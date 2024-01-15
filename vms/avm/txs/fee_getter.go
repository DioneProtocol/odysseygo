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

func calculateInputs(ins []*avax.TransferableInput, assetId ids.ID) uint64 {
	var totalInputs uint64
	for _, i := range ins {
		if i.AssetID() == assetId {
			totalInputs += i.In.Amount()
		}
	}
	return totalInputs
}

func calculateOutputs(outs []*avax.TransferableOutput, assetId ids.ID) uint64 {
	var totalOutputs uint64
	for _, o := range outs {
		if o.AssetID() == assetId {
			totalOutputs += o.Out.Amount()
		}
	}
	return totalOutputs
}

func (b *burned) setDifference(tx *avax.BaseTx) error {
	ins := calculateInputs(tx.Ins, b.assetId)
	outs := calculateOutputs(tx.Outs, b.assetId)
	if ins > outs {
		b.burned = ins - outs
	} else {
		b.burned = 0
	}
	return nil
}

func (b *burned) BaseTx(tx *BaseTx) error {
	return b.setDifference(&tx.BaseTx)
}

func (b *burned) CreateAssetTx(tx *CreateAssetTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}

func (b *burned) ExportTx(tx *ExportTx) error {
	baseTx := &tx.BaseTx.BaseTx
	ins := calculateInputs(baseTx.Ins, b.assetId)
	baseOuts := calculateOutputs(baseTx.Outs, b.assetId)
	exportedOuts := calculateOutputs(tx.ExportedOuts, b.assetId)
	outs := baseOuts + exportedOuts
	if ins > outs {
		b.burned = ins - outs
	} else {
		b.burned = 0
	}
	return nil
}

func (b *burned) ImportTx(tx *ImportTx) error {
	baseTx := &tx.BaseTx.BaseTx
	outs := calculateOutputs(baseTx.Outs, b.assetId)
	baseIns := calculateInputs(baseTx.Ins, b.assetId)
	importedIns := calculateInputs(tx.ImportedIns, b.assetId)
	ins := baseIns + importedIns
	if ins > outs {
		b.burned = ins - outs
	} else {
		b.burned = 0
	}
	return nil
}

func (b *burned) OperationTx(tx *OperationTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}
