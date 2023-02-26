// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"

	"github.com/dioneprotocol/dionego/codec"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/utils/set"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

var (
	errNoImportInputs = errors.New("no import inputs")

	_ UnsignedTx             = (*ImportTx)(nil)
	_ secp256k1fx.UnsignedTx = (*ImportTx)(nil)
)

// ImportTx is a transaction that imports an asset from another blockchain.
type ImportTx struct {
	BaseTx `serialize:"true"`

	// Which chain to consume the funds from
	SourceChain ids.ID `serialize:"true" json:"sourceChain"`

	// The inputs to this transaction
	ImportedIns []*dione.TransferableInput `serialize:"true" json:"importedInputs"`
}

// InputUTXOs track which UTXOs this transaction is consuming.
func (t *ImportTx) InputUTXOs() []*dione.UTXOID {
	utxos := t.BaseTx.InputUTXOs()
	for _, in := range t.ImportedIns {
		in.Symbol = true
		utxos = append(utxos, &in.UTXOID)
	}
	return utxos
}

func (t *ImportTx) InputIDs() set.Set[ids.ID] {
	inputs := t.BaseTx.InputIDs()
	for _, in := range t.ImportedIns {
		inputs.Add(in.InputID())
	}
	return inputs
}

// ConsumedAssetIDs returns the IDs of the assets this transaction consumes
func (t *ImportTx) ConsumedAssetIDs() set.Set[ids.ID] {
	assets := t.BaseTx.AssetIDs()
	for _, in := range t.ImportedIns {
		assets.Add(in.AssetID())
	}
	return assets
}

// AssetIDs returns the IDs of the assets this transaction depends on
func (t *ImportTx) AssetIDs() set.Set[ids.ID] {
	assets := t.BaseTx.AssetIDs()
	for _, in := range t.ImportedIns {
		assets.Add(in.AssetID())
	}
	return assets
}

// NumCredentials returns the number of expected credentials
func (t *ImportTx) NumCredentials() int {
	return t.BaseTx.NumCredentials() + len(t.ImportedIns)
}

// SyntacticVerify that this import transaction is well-formed.
func (t *ImportTx) SyntacticVerify(
	ctx *snow.Context,
	c codec.Manager,
	txFeeAssetID ids.ID,
	txFee uint64,
	_ uint64,
	_ int,
) error {
	switch {
	case t == nil:
		return errNilTx
	case len(t.ImportedIns) == 0:
		return errNoImportInputs
	}

	// We don't call [t.BaseTx.SyntacticVerify] because the flow check performed
	// here is less strict than the flow check performed in the [BaseTx].
	if err := t.BaseTx.BaseTx.Verify(ctx); err != nil {
		return err
	}

	return dione.VerifyTx(
		txFee,
		txFeeAssetID,
		[][]*dione.TransferableInput{
			t.Ins,
			t.ImportedIns,
		},
		[][]*dione.TransferableOutput{t.Outs},
		c,
	)
}

func (t *ImportTx) Visit(v Visitor) error {
	return v.ImportTx(t)
}
