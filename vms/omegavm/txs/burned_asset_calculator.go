// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
)

var (
	_ Visitor = (*BurnedAssetCalculator)(nil)
)

type BurnedAssetCalculator struct {
	tx      *Tx
	assetId ids.ID
	burned  uint64
}

type stakeGetter interface {
	Stake() []*dione.TransferableOutput
}

func calculateInputs(ins []*dione.TransferableInput, assetId ids.ID) uint64 {
	var totalInputs uint64
	for _, i := range ins {
		if i.AssetID() == assetId {
			totalInputs += i.In.Amount()
		}
	}
	return totalInputs
}

func calculateOutputs(outs []*dione.TransferableOutput, assetId ids.ID) uint64 {
	var totalOutputs uint64
	for _, o := range outs {
		if o.AssetID() == assetId {
			totalOutputs += o.Out.Amount()
		}
	}
	return totalOutputs
}

func (b *BurnedAssetCalculator) setDifference(tx *dione.BaseTx) error {
	ins := calculateInputs(tx.Ins, b.assetId)
	outs := calculateOutputs(tx.Outs, b.assetId)
	if ins > outs {
		b.burned = ins - outs
	}
	return nil
}

func (b *BurnedAssetCalculator) setDifferenceWithStake(tx *dione.BaseTx, s stakeGetter) error {
	ins := calculateInputs(tx.Ins, b.assetId)
	baseOuts := calculateOutputs(tx.Outs, b.assetId)
	stakeOuts := calculateOutputs(s.Stake(), b.assetId)
	outs := baseOuts + stakeOuts
	if ins > outs {
		b.burned = ins - outs
	}
	return nil
}

func (b *BurnedAssetCalculator) AddDelegatorTx(tx *AddDelegatorTx) error {
	return b.setDifferenceWithStake(&tx.BaseTx.BaseTx, tx)
}

func (b *BurnedAssetCalculator) AddPermissionlessDelegatorTx(tx *AddPermissionlessDelegatorTx) error {
	return b.setDifferenceWithStake(&tx.BaseTx.BaseTx, tx)
}

func (b *BurnedAssetCalculator) AddPermissionlessValidatorTx(tx *AddPermissionlessValidatorTx) error {
	return b.setDifferenceWithStake(&tx.BaseTx.BaseTx, tx)
}

func (b *BurnedAssetCalculator) AddSubnetValidatorTx(tx *AddSubnetValidatorTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}

func (b *BurnedAssetCalculator) AddValidatorTx(tx *AddValidatorTx) error {
	return b.setDifferenceWithStake(&tx.BaseTx.BaseTx, tx)
}

func (*BurnedAssetCalculator) AdvanceTimeTx(*AdvanceTimeTx) error {
	return nil
}

func (b *BurnedAssetCalculator) CreateChainTx(tx *CreateChainTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}

func (b *BurnedAssetCalculator) CreateSubnetTx(tx *CreateSubnetTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}

func (b *BurnedAssetCalculator) ExportTx(tx *ExportTx) error {
	ins := calculateInputs(tx.Ins, b.assetId)
	baseOuts := calculateOutputs(tx.Outs, b.assetId)
	exportedOuts := calculateOutputs(tx.ExportedOutputs, b.assetId)
	outs := baseOuts + exportedOuts
	if ins > outs {
		b.burned = ins - outs
	}
	return nil
}

func (b *BurnedAssetCalculator) ImportTx(tx *ImportTx) error {
	baseIns := calculateInputs(tx.Ins, b.assetId)
	importedIns := calculateInputs(tx.ImportedInputs, b.assetId)
	outs := calculateOutputs(tx.Outs, b.assetId)
	ins := baseIns + importedIns
	if ins > outs {
		b.burned = ins - outs
	}
	return nil
}

func (b *BurnedAssetCalculator) RemoveSubnetValidatorTx(tx *RemoveSubnetValidatorTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}

func (*BurnedAssetCalculator) RewardValidatorTx(*RewardValidatorTx) error {
	return nil
}

func (b *BurnedAssetCalculator) TransformSubnetTx(tx *TransformSubnetTx) error {
	return b.setDifference(&tx.BaseTx.BaseTx)
}
