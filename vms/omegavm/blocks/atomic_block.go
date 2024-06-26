// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"fmt"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

var _ Block = (*ApricotAtomicBlock)(nil)

// ApricotAtomicBlock being accepted results in the atomic transaction contained
// in the block to be accepted and committed to the chain.
type ApricotAtomicBlock struct {
	CommonBlock `serialize:"true"`
	Tx          *txs.Tx `serialize:"true" json:"tx"`
}

func (b *ApricotAtomicBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	if err := b.Tx.Initialize(txs.Codec); err != nil {
		return fmt.Errorf("failed to initialize tx: %w", err)
	}
	return nil
}

func (b *ApricotAtomicBlock) InitCtx(ctx *snow.Context) {
	b.Tx.Unsigned.InitCtx(ctx)
}

func (b *ApricotAtomicBlock) Txs() []*txs.Tx {
	return []*txs.Tx{b.Tx}
}

func (b *ApricotAtomicBlock) Visit(v Visitor) error {
	return v.ApricotAtomicBlock(b)
}

func (b *ApricotAtomicBlock) FeeFromAChain() uint64 {
	return 0
}

func (b *ApricotAtomicBlock) FeeFromDChain() uint64 {
	return 0
}

func (b *ApricotAtomicBlock) FeeFromOChain(ids.ID) uint64 {
	return 0
}

func (b *ApricotAtomicBlock) AccumulatedFee(ids.ID) uint64 {
	return 0
}

func NewApricotAtomicBlock(
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*ApricotAtomicBlock, error) {
	blk := &ApricotAtomicBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		Tx: tx,
	}
	return blk, initialize(blk)
}
