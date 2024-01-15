// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/components/avax"
)

var (
	errOutsGreaterThanInputs = errors.New("outputs are greater than inputs")

	_ Visitor = (*burned)(nil)
)

type burned struct {
	tx      *Tx
	assetId ids.ID
	burned  uint64
}

func (b *burned) calculateBurned(tx *avax.BaseTx) error {
	var totalInputs, totalOutputs uint64
	for _, i := range tx.Ins {
		if i.AssetID() == b.assetId {
			totalInputs += i.In.Amount()
		}
	}
	for _, o := range tx.Outs {
		if o.AssetID() == b.assetId {
			totalOutputs += o.Out.Amount()
		}
	}
	if totalOutputs < totalInputs {
		b.burned = totalInputs - totalOutputs
	} else {
		b.burned = 0
	}
	return nil
}

func (b *burned) BaseTx(tx *BaseTx) error {
	return b.calculateBurned(&tx.BaseTx)
}

func (b *burned) CreateAssetTx(tx *CreateAssetTx) error {
	return b.calculateBurned(&tx.BaseTx.BaseTx)
}

func (b *burned) ExportTx(tx *ExportTx) error {
	return b.calculateBurned(&tx.BaseTx.BaseTx)
}

func (b *burned) ImportTx(tx *ImportTx) error {
	return b.calculateBurned(&tx.BaseTx.BaseTx)
}

func (b *burned) OperationTx(tx *OperationTx) error {
	return b.calculateBurned(&tx.BaseTx.BaseTx)
}
