// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/dioneprotocol/dionego/codec"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/utils/set"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

var (
	_ UnsignedTx             = (*BaseTx)(nil)
	_ secp256k1fx.UnsignedTx = (*BaseTx)(nil)
)

// BaseTx is the basis of all transactions.
type BaseTx struct {
	dione.BaseTx `serialize:"true"`

	bytes []byte
}

func (t *BaseTx) InitCtx(ctx *snow.Context) {
	for _, out := range t.Outs {
		out.InitCtx(ctx)
	}
}

func (t *BaseTx) SetBytes(bytes []byte) {
	t.bytes = bytes
}

func (t *BaseTx) Bytes() []byte {
	return t.bytes
}

func (t *BaseTx) InputIDs() set.Set[ids.ID] {
	inputIDs := set.NewSet[ids.ID](len(t.Ins))
	for _, in := range t.Ins {
		inputIDs.Add(in.InputID())
	}
	return inputIDs
}

func (t *BaseTx) SyntacticVerify(
	ctx *snow.Context,
	c codec.Manager,
	txFeeAssetID ids.ID,
	txFee uint64,
	_ uint64,
	_ int,
) error {
	if t == nil {
		return errNilTx
	}

	if err := t.BaseTx.Verify(ctx); err != nil {
		return err
	}

	return dione.VerifyTx(
		txFee,
		txFeeAssetID,
		[][]*dione.TransferableInput{t.Ins},
		[][]*dione.TransferableOutput{t.Outs},
		c,
	)
}

func (t *BaseTx) Visit(v Visitor) error {
	return v.BaseTx(t)
}
