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
	_ Block      = (*OdysseyProposalBlock)(nil)
)

type BanffProposalBlock struct {
	Time uint64 `serialize:"true" json:"time"`
	// Transactions is currently unused. This is populated so that introducing
	// them in the future will not require a codec change.
	//
	// TODO: when Transactions is used, we must correctly verify and apply their
	//       changes.
	Transactions         []*txs.Tx `serialize:"true" json:"-"`
	OdysseyProposalBlock `serialize:"true"`
}

func (b *BanffProposalBlock) InitCtx(ctx *snow.Context) {
	for _, tx := range b.Transactions {
		tx.Unsigned.InitCtx(ctx)
	}
	b.OdysseyProposalBlock.InitCtx(ctx)
}

func (b *BanffProposalBlock) Timestamp() time.Time {
	return time.Unix(int64(b.Time), 0)
}

func (b *BanffProposalBlock) Visit(v Visitor) error {
	return v.BanffProposalBlock(b)
}

func NewBanffProposalBlock(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*BanffProposalBlock, error) {
	blk := &BanffProposalBlock{
		Time: uint64(timestamp.Unix()),
		OdysseyProposalBlock: OdysseyProposalBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			Tx: tx,
		},
	}
	return blk, initialize(blk)
}

type OdysseyProposalBlock struct {
	CommonBlock `serialize:"true"`
	Tx          *txs.Tx `serialize:"true" json:"tx"`
}

func (b *OdysseyProposalBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	if err := b.Tx.Initialize(txs.Codec); err != nil {
		return fmt.Errorf("failed to initialize tx: %w", err)
	}
	return nil
}

func (b *OdysseyProposalBlock) InitCtx(ctx *snow.Context) {
	b.Tx.Unsigned.InitCtx(ctx)
}

func (b *OdysseyProposalBlock) Txs() []*txs.Tx {
	return []*txs.Tx{b.Tx}
}

func (b *OdysseyProposalBlock) Visit(v Visitor) error {
	return v.OdysseyProposalBlock(b)
}

// NewOdysseyProposalBlock is kept for testing purposes only.
// Following Banff activation and subsequent code cleanup, Odyssey Proposal blocks
// should be only verified (upon bootstrap), never created anymore
func NewOdysseyProposalBlock(
	parentID ids.ID,
	height uint64,
	tx *txs.Tx,
) (*OdysseyProposalBlock, error) {
	blk := &OdysseyProposalBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		Tx: tx,
	}
	return blk, initialize(blk)
}
