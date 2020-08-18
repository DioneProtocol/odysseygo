// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"errors"
	"fmt"

	"github.com/ava-labs/gecko/chains/atomic"
	"github.com/ava-labs/gecko/database"
	"github.com/ava-labs/gecko/database/versiondb"
	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/utils/codec"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/vms/components/avax"
	"github.com/ava-labs/gecko/vms/secp256k1fx"

	safemath "github.com/ava-labs/gecko/utils/math"
)

var (
	errNoExportOutputs  = errors.New("no export outputs")
	errOutputsNotSorted = errors.New("outputs not sorted")
	errOverflowExport   = errors.New("overflow when computing export amount + txFee")
	errWrongChainID     = errors.New("tx has wrong chain ID")

	_ UnsignedAtomicTx = &UnsignedExportTx{}
)

// UnsignedExportTx is an unsigned ExportTx
type UnsignedExportTx struct {
	BaseTx `serialize:"true"`

	// Which chain to send the funds to
	DestinationChain ids.ID `serialize:"true" json:"destinationChain"`

	// Outputs that are exported to the chain
	ExportedOutputs []*avax.TransferableOutput `serialize:"true" json:"exportedOutputs"`
}

// InputUTXOs returns an empty set
func (tx *UnsignedExportTx) InputUTXOs() ids.Set { return ids.Set{} }

// Verify this transaction is well-formed
func (tx *UnsignedExportTx) Verify(
	avmID ids.ID,
	ctx *snow.Context,
	c codec.Codec,
	feeAmount uint64,
	feeAssetID ids.ID,
) error {
	switch {
	case tx == nil:
		return errNilTx
	case tx.syntacticallyVerified: // already passed syntactic verification
		return nil
	case tx.DestinationChain.IsZero():
		return errWrongChainID
	case !tx.DestinationChain.Equals(avmID):
		// TODO: remove this check if we allow for P->C swaps
		return errWrongChainID
	case len(tx.ExportedOutputs) == 0:
		return errNoExportOutputs
	}

	if err := tx.BaseTx.Verify(ctx, c); err != nil {
		return err
	}

	for _, out := range tx.ExportedOutputs {
		if err := out.Verify(); err != nil {
			return err
		}
		if _, ok := out.Output().(*StakeableLockOut); ok {
			return errWrongLocktime
		}
	}
	if !avax.IsSortedTransferableOutputs(tx.ExportedOutputs, Codec) {
		return errOutputsNotSorted
	}

	allOuts := make([]*avax.TransferableOutput, len(tx.Outs)+len(tx.ExportedOutputs))
	copy(allOuts, tx.Outs)
	copy(allOuts[len(tx.Outs):], tx.ExportedOutputs)
	if err := syntacticVerifySpend(tx.Ins, allOuts, nil, 0, feeAmount, feeAssetID); err != nil {
		return err
	}

	tx.syntacticallyVerified = true
	return nil
}

// SemanticVerify this transaction is valid.
func (tx *UnsignedExportTx) SemanticVerify(
	vm *VM,
	db database.Database,
	stx *Tx,
) TxError {
	if err := tx.Verify(vm.avm, vm.Ctx, vm.codec, vm.txFee, vm.avaxAssetID); err != nil {
		return permError{err}
	}

	outs := make([]*avax.TransferableOutput, len(tx.Outs)+len(tx.ExportedOutputs))
	copy(outs, tx.Outs)
	copy(outs[len(tx.Outs):], tx.ExportedOutputs)

	// Verify the flowcheck
	if err := vm.semanticVerifySpend(db, tx, tx.Ins, outs, stx.Creds); err != nil {
		return err
	}

	txID := tx.ID()

	// Consume the UTXOS
	if err := vm.consumeInputs(db, tx.Ins); err != nil {
		return tempError{err}
	}
	// Produce the UTXOS
	if err := vm.produceOutputs(db, txID, tx.Outs); err != nil {
		return tempError{err}
	}
	return nil
}

// Accept this transaction.
func (tx *UnsignedExportTx) Accept(ctx *snow.Context, batch database.Batch) error {
	// Produce exported UTXOs
	smDB := ctx.SharedMemory.GetDatabase(tx.DestinationChain)
	defer ctx.SharedMemory.ReleaseDatabase(tx.DestinationChain)

	vsmDB := versiondb.New(smDB)
	state := avax.NewPrefixedState(vsmDB, Codec, ctx.ChainID, tx.DestinationChain)

	txID := tx.ID()
	for i, out := range tx.ExportedOutputs {
		if err := state.FundUTXO(&avax.UTXO{
			UTXOID: avax.UTXOID{
				TxID:        txID,
				OutputIndex: uint32(len(tx.BaseTx.Outs) + i),
			},
			Asset: out.Asset,
			Out:   out.Output(),
		}); err != nil {
			return err
		}
	}

	sharedBatch, err := vsmDB.CommitBatch()
	if err != nil {
		return err
	}
	return atomic.WriteAll(batch, sharedBatch)
}

// Create a new transaction
func (vm *VM) newExportTx(
	amount uint64, // Amount of tokens to export
	chainID ids.ID, // Chain to send the UTXOs to
	to ids.ShortID, // Address of chain recipient
	keys []*crypto.PrivateKeySECP256K1R, // Pay the fee and provide the tokens
) (*Tx, error) {
	if !vm.avm.Equals(chainID) {
		return nil, errWrongChainID
	}

	toBurn, err := safemath.Add64(amount, vm.txFee)
	if err != nil {
		return nil, errOverflowExport
	}
	ins, outs, _, signers, err := vm.stake(vm.DB, keys, 0, toBurn)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate tx inputs/outputs: %w", err)
	}

	// Create the transaction
	utx := &UnsignedExportTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    vm.Ctx.NetworkID,
			BlockchainID: vm.Ctx.ChainID,
			Ins:          ins,
			Outs:         outs, // Non-exported outputs
		}},
		DestinationChain: chainID,
		ExportedOutputs: []*avax.TransferableOutput{{ // Exported to X-Chain
			Asset: avax.Asset{ID: vm.avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: amount,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Threshold: 1,
					Addrs:     []ids.ShortID{to},
				},
			},
		}},
	}
	tx := &Tx{UnsignedTx: utx}
	if err := tx.Sign(vm.codec, signers); err != nil {
		return nil, err
	}
	return tx, utx.Verify(vm.avm, vm.Ctx, vm.codec, vm.txFee, vm.avaxAssetID)
}
