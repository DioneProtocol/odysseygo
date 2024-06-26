// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"fmt"
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

var (
	_ BanffBlock = (*BanffProposalBlock)(nil)
	_ Block      = (*ApricotProposalBlock)(nil)
)

type BanffProposalBlock struct {
	Time uint64 `serialize:"true" json:"time"`
	// Transactions is currently unused. This is populated so that introducing
	// them in the future will not require a codec change.
	//
	// TODO: when Transactions is used, we must correctly verify and apply their
	//       changes.
	Transactions         []*txs.Tx `serialize:"true" json:"-"`
	ApricotProposalBlock `serialize:"true"`
}

func (b *BanffProposalBlock) InitCtx(ctx *snow.Context) {
	for _, tx := range b.Transactions {
		tx.Unsigned.InitCtx(ctx)
	}
	b.ApricotProposalBlock.InitCtx(ctx)
}

func (b *BanffProposalBlock) Timestamp() time.Time {
	return time.Unix(int64(b.Time), 0)
}

func (b *BanffProposalBlock) Visit(v Visitor) error {
	return v.BanffProposalBlock(b)
}

func (b *BanffProposalBlock) FeeFromOChain(assetID ids.ID) uint64 {
	feeOChain := b.ApricotProposalBlock.FeeFromOChain(assetID)
	for _, tx := range b.Transactions {
		feeOChain += tx.Burned(assetID)
	}
	return feeOChain
}

func NewBanffProposalBlock(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*BanffProposalBlock, error) {
	blk := &BanffProposalBlock{
		Time: uint64(timestamp.Unix()),
		ApricotProposalBlock: ApricotProposalBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			Tx: tx,
		},
	}
	return blk, initialize(blk)
}

func NewBanffProposalBlockWithFee(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
	feeFromAChain uint64,
	feeFromDChain uint64,
) (*BanffProposalBlock, error) {
	blk := &BanffProposalBlock{
		Time: uint64(timestamp.Unix()),
		ApricotProposalBlock: ApricotProposalBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			FeeAChain: feeFromAChain,
			FeeDChain: feeFromDChain,
			Tx:        tx,
		},
	}
	return blk, initialize(blk)
}

type ApricotProposalBlock struct {
	CommonBlock `serialize:"true"`
	Tx          *txs.Tx `serialize:"true" json:"tx"`
	FeeAChain   uint64  `serialize:"true" json:"feeFromAChain"`
	FeeDChain   uint64  `serialize:"true" json:"feeFromDChain"`
}

func (b *ApricotProposalBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	if err := b.Tx.Initialize(txs.Codec); err != nil {
		return fmt.Errorf("failed to initialize tx: %w", err)
	}
	return nil
}

func (b *ApricotProposalBlock) InitCtx(ctx *snow.Context) {
	b.Tx.Unsigned.InitCtx(ctx)
}

func (b *ApricotProposalBlock) Txs() []*txs.Tx {
	return []*txs.Tx{b.Tx}
}

func (b *ApricotProposalBlock) Visit(v Visitor) error {
	return v.ApricotProposalBlock(b)
}

func (b *ApricotProposalBlock) FeeFromAChain() uint64 {
	return b.FeeAChain
}

func (b *ApricotProposalBlock) FeeFromDChain() uint64 {
	return b.FeeDChain
}

func (b *ApricotProposalBlock) FeeFromOChain(assetID ids.ID) uint64 {
	return b.Tx.Burned(assetID)
}

func (b *ApricotProposalBlock) AccumulatedFee(assetID ids.ID) uint64 {
	accumulatedFee := b.FeeFromOChain(assetID)
	accumulatedFee += b.FeeFromAChain()
	accumulatedFee += b.FeeFromDChain()
	return accumulatedFee
}

// NewApricotProposalBlock is kept for testing purposes only.
// Following Banff activation and subsequent code cleanup, Apricot Proposal blocks
// should be only verified (upon bootstrap), never created anymore
func NewApricotProposalBlock(
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*ApricotProposalBlock, error) {
	blk := &ApricotProposalBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		Tx: tx,
	}
	return blk, initialize(blk)
}
