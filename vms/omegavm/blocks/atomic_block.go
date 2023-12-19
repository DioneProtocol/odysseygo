// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"fmt"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

var _ Block = (*OdysseyAtomicBlock)(nil)

// OdysseyAtomicBlock being accepted results in the atomic transaction contained
// in the block to be accepted and committed to the chain.
type OdysseyAtomicBlock struct {
	CommonBlock `serialize:"true"`
	Tx          *txs.Tx `serialize:"true" json:"tx"`
}

func (b *OdysseyAtomicBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	if err := b.Tx.Initialize(txs.Codec); err != nil {
		return fmt.Errorf("failed to initialize tx: %w", err)
	}
	return nil
}

func (b *OdysseyAtomicBlock) InitCtx(ctx *snow.Context) {
	b.Tx.Unsigned.InitCtx(ctx)
}

func (b *OdysseyAtomicBlock) Txs() []*txs.Tx {
	return []*txs.Tx{b.Tx}
}

func (b *OdysseyAtomicBlock) Visit(v Visitor) error {
	return v.OdysseyAtomicBlock(b)
}

func NewOdysseyAtomicBlock(
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*OdysseyAtomicBlock, error) {
	blk := &OdysseyAtomicBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		Tx: tx,
	}
	return blk, initialize(blk)
}
