// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package d

import (
	"errors"
	"math/big"

	stdcontext "context"

	"github.com/DioneProtocol/coreth/plugin/delta"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

const dioneConversionRateInt = 1_000_000_000

var (
	_ Builder = (*builder)(nil)

	errInsufficientFunds = errors.New("insufficient funds")

	// dioneConversionRate is the conversion rate between the smallest
	// denomination on the X-Chain and O-chain, 1 nDIONE, and the smallest
	// denomination on the D-Chain 1 wei. Where 1 nDIONE = 1 gWei.
	//
	// This is only required for DIONE because the denomination of 1 DIONE is 9
	// decimal places on the X and O chains, but is 18 decimal places within the
	// DELTA.
	dioneConversionRate = big.NewInt(dioneConversionRateInt)
)

// Builder provides a convenient interface for building unsigned D-chain
// transactions.
type Builder interface {
	// GetBalance calculates the amount of DIONE that this builder has control
	// over.
	GetBalance(
		options ...common.Option,
	) (*big.Int, error)

	// GetImportableBalance calculates the amount of DIONE that this builder
	// could import from the provided chain.
	//
	// - [chainID] specifies the chain the funds are from.
	GetImportableBalance(
		chainID ids.ID,
		options ...common.Option,
	) (uint64, error)

	// NewImportTx creates an import transaction that attempts to consume all
	// the available UTXOs and import the funds to [to].
	//
	// - [chainID] specifies the chain to be importing funds from.
	// - [to] specifies where to send the imported funds to.
	// - [baseFee] specifies the fee price willing to be paid by this tx.
	NewImportTx(
		chainID ids.ID,
		to ethcommon.Address,
		baseFee *big.Int,
		options ...common.Option,
	) (*delta.UnsignedImportTx, error)

	// NewExportTx creates an export transaction that attempts to send all the
	// provided [outputs] to the requested [chainID].
	//
	// - [chainID] specifies the chain to be exporting the funds to.
	// - [outputs] specifies the outputs to send to the [chainID].
	// - [baseFee] specifies the fee price willing to be paid by this tx.
	NewExportTx(
		chainID ids.ID,
		outputs []*secp256k1fx.TransferOutput,
		baseFee *big.Int,
		options ...common.Option,
	) (*delta.UnsignedExportTx, error)
}

// BuilderBackend specifies the required information needed to build unsigned
// D-chain transactions.
type BuilderBackend interface {
	Context

	UTXOs(ctx stdcontext.Context, sourceChainID ids.ID) ([]*dione.UTXO, error)
	Balance(ctx stdcontext.Context, addr ethcommon.Address) (*big.Int, error)
	Nonce(ctx stdcontext.Context, addr ethcommon.Address) (uint64, error)
}

type builder struct {
	dioneAddrs set.Set[ids.ShortID]
	ethAddrs   set.Set[ethcommon.Address]
	backend    BuilderBackend
}

// NewBuilder returns a new transaction builder.
//
//   - [dioneAddrs] is the set of addresses in the DIONE format that the builder
//     assumes can be used when signing the transactions in the future.
//   - [ethAddrs] is the set of addresses in the Eth format that the builder
//     assumes can be used when signing the transactions in the future.
//   - [backend] provides the required access to the chain's context and state
//     to build out the transactions.
func NewBuilder(
	dioneAddrs set.Set[ids.ShortID],
	ethAddrs set.Set[ethcommon.Address],
	backend BuilderBackend,
) Builder {
	return &builder{
		dioneAddrs: dioneAddrs,
		ethAddrs:   ethAddrs,
		backend:    backend,
	}
}

func (b *builder) GetBalance(
	options ...common.Option,
) (*big.Int, error) {
	var (
		ops          = common.NewOptions(options)
		ctx          = ops.Context()
		addrs        = ops.EthAddresses(b.ethAddrs)
		totalBalance = new(big.Int)
	)
	for addr := range addrs {
		balance, err := b.backend.Balance(ctx, addr)
		if err != nil {
			return nil, err
		}
		totalBalance.Add(totalBalance, balance)
	}

	return totalBalance, nil
}

func (b *builder) GetImportableBalance(
	chainID ids.ID,
	options ...common.Option,
) (uint64, error) {
	ops := common.NewOptions(options)
	utxos, err := b.backend.UTXOs(ops.Context(), chainID)
	if err != nil {
		return 0, err
	}

	var (
		addrs           = ops.Addresses(b.dioneAddrs)
		minIssuanceTime = ops.MinIssuanceTime()
		dioneAssetID    = b.backend.DIONEAssetID()
		balance         uint64
	)
	for _, utxo := range utxos {
		amount, _, ok := getSpendableAmount(utxo, addrs, minIssuanceTime, dioneAssetID)
		if !ok {
			continue
		}

		newBalance, err := math.Add64(balance, amount)
		if err != nil {
			return 0, err
		}
		balance = newBalance
	}

	return balance, nil
}

func (b *builder) NewImportTx(
	chainID ids.ID,
	to ethcommon.Address,
	baseFee *big.Int,
	options ...common.Option,
) (*delta.UnsignedImportTx, error) {
	ops := common.NewOptions(options)
	utxos, err := b.backend.UTXOs(ops.Context(), chainID)
	if err != nil {
		return nil, err
	}

	var (
		addrs           = ops.Addresses(b.dioneAddrs)
		minIssuanceTime = ops.MinIssuanceTime()
		dioneAssetID    = b.backend.DIONEAssetID()

		importedInputs = make([]*dione.TransferableInput, 0, len(utxos))
		importedAmount uint64
	)
	for _, utxo := range utxos {
		amount, inputSigIndices, ok := getSpendableAmount(utxo, addrs, minIssuanceTime, dioneAssetID)
		if !ok {
			continue
		}

		importedInputs = append(importedInputs, &dione.TransferableInput{
			UTXOID: utxo.UTXOID,
			Asset:  utxo.Asset,
			In: &secp256k1fx.TransferInput{
				Amt: amount,
				Input: secp256k1fx.Input{
					SigIndices: inputSigIndices,
				},
			},
		})

		newImportedAmount, err := math.Add64(importedAmount, amount)
		if err != nil {
			return nil, err
		}
		importedAmount = newImportedAmount
	}

	utils.Sort(importedInputs)
	tx := &delta.UnsignedImportTx{
		NetworkID:      b.backend.NetworkID(),
		BlockchainID:   b.backend.BlockchainID(),
		SourceChain:    chainID,
		ImportedInputs: importedInputs,
	}

	// We must initialize the bytes of the tx to calculate the initial cost
	wrappedTx := &delta.Tx{UnsignedAtomicTx: tx}
	if err := wrappedTx.Sign(delta.Codec, nil); err != nil {
		return nil, err
	}

	gasUsedWithoutOutput, err := tx.GasUsed(true /*=IsApricotPhase5*/)
	if err != nil {
		return nil, err
	}
	gasUsedWithOutput := gasUsedWithoutOutput + delta.DELTAOutputGas

	txFee, err := delta.CalculateDynamicFee(gasUsedWithOutput, baseFee)
	if err != nil {
		return nil, err
	}

	if importedAmount <= txFee {
		return nil, errInsufficientFunds
	}

	tx.Outs = []delta.DELTAOutput{{
		Address: to,
		Amount:  importedAmount - txFee,
		AssetID: dioneAssetID,
	}}
	return tx, nil
}

func (b *builder) NewExportTx(
	chainID ids.ID,
	outputs []*secp256k1fx.TransferOutput,
	baseFee *big.Int,
	options ...common.Option,
) (*delta.UnsignedExportTx, error) {
	var (
		dioneAssetID    = b.backend.DIONEAssetID()
		exportedOutputs = make([]*dione.TransferableOutput, len(outputs))
		exportedAmount  uint64
	)
	for i, output := range outputs {
		exportedOutputs[i] = &dione.TransferableOutput{
			Asset: dione.Asset{ID: dioneAssetID},
			Out:   output,
		}

		newExportedAmount, err := math.Add64(exportedAmount, output.Amt)
		if err != nil {
			return nil, err
		}
		exportedAmount = newExportedAmount
	}

	dione.SortTransferableOutputs(exportedOutputs, delta.Codec)
	tx := &delta.UnsignedExportTx{
		NetworkID:        b.backend.NetworkID(),
		BlockchainID:     b.backend.BlockchainID(),
		DestinationChain: chainID,
		ExportedOutputs:  exportedOutputs,
	}

	// We must initialize the bytes of the tx to calculate the initial cost
	wrappedTx := &delta.Tx{UnsignedAtomicTx: tx}
	if err := wrappedTx.Sign(delta.Codec, nil); err != nil {
		return nil, err
	}

	cost, err := tx.GasUsed(true /*=IsApricotPhase5*/)
	if err != nil {
		return nil, err
	}

	initialFee, err := delta.CalculateDynamicFee(cost, baseFee)
	if err != nil {
		return nil, err
	}

	amountToConsume, err := math.Add64(exportedAmount, initialFee)
	if err != nil {
		return nil, err
	}

	var (
		ops    = common.NewOptions(options)
		ctx    = ops.Context()
		addrs  = ops.EthAddresses(b.ethAddrs)
		inputs = make([]delta.DELTAInput, 0, addrs.Len())
	)
	for addr := range addrs {
		if amountToConsume == 0 {
			break
		}

		prevFee, err := delta.CalculateDynamicFee(cost, baseFee)
		if err != nil {
			return nil, err
		}

		newCost := cost + delta.DELTAInputGas
		newFee, err := delta.CalculateDynamicFee(newCost, baseFee)
		if err != nil {
			return nil, err
		}

		additionalFee := newFee - prevFee

		balance, err := b.backend.Balance(ctx, addr)
		if err != nil {
			return nil, err
		}

		// Since the asset is DIONE, we divide by the dioneConversionRate to
		// convert back to the correct denomination of DIONE that can be
		// exported.
		dioneBalance := new(big.Int).Div(balance, dioneConversionRate).Uint64()

		// If the balance for [addr] is insufficient to cover the additional
		// cost of adding an input to the transaction, skip adding the input
		// altogether.
		if dioneBalance <= additionalFee {
			continue
		}

		// Update the cost for the next iteration
		cost = newCost

		amountToConsume, err = math.Add64(amountToConsume, additionalFee)
		if err != nil {
			return nil, err
		}

		nonce, err := b.backend.Nonce(ctx, addr)
		if err != nil {
			return nil, err
		}

		inputAmount := math.Min(amountToConsume, dioneBalance)
		inputs = append(inputs, delta.DELTAInput{
			Address: addr,
			Amount:  inputAmount,
			AssetID: dioneAssetID,
			Nonce:   nonce,
		})
		amountToConsume -= inputAmount
	}

	if amountToConsume > 0 {
		return nil, errInsufficientFunds
	}

	utils.Sort(inputs)
	tx.Ins = inputs
	return tx, nil
}

func getSpendableAmount(
	utxo *dione.UTXO,
	addrs set.Set[ids.ShortID],
	minIssuanceTime uint64,
	dioneAssetID ids.ID,
) (uint64, []uint32, bool) {
	if utxo.Asset.ID != dioneAssetID {
		// Only DIONE can be imported
		return 0, nil, false
	}

	out, ok := utxo.Out.(*secp256k1fx.TransferOutput)
	if !ok {
		// Can't import an unknown transfer output type
		return 0, nil, false
	}

	inputSigIndices, ok := common.MatchOwners(&out.OutputOwners, addrs, minIssuanceTime)
	return out.Amt, inputSigIndices, ok
}
