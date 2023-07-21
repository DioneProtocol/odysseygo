// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/txs"
)

var (
	_ BanffBlock = (*BanffAbortBlock)(nil)
	_ Block      = (*OdysseyAbortBlock)(nil)
)

type BanffAbortBlock struct {
	Time              uint64 `serialize:"true" json:"time"`
	OdysseyAbortBlock `serialize:"true"`
}

func (b *BanffAbortBlock) Timestamp() time.Time {
	return time.Unix(int64(b.Time), 0)
}

func (b *BanffAbortBlock) Visit(v Visitor) error {
	return v.BanffAbortBlock(b)
}

func NewBanffAbortBlock(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
) (*BanffAbortBlock, error) {
	blk := &BanffAbortBlock{
		Time: uint64(timestamp.Unix()),
		OdysseyAbortBlock: OdysseyAbortBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
		},
	}
	return blk, initialize(blk)
}

type OdysseyAbortBlock struct {
	CommonBlock `serialize:"true"`
}

func (b *OdysseyAbortBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	return nil
}

func (*OdysseyAbortBlock) InitCtx(*snow.Context) {}

func (*OdysseyAbortBlock) Txs() []*txs.Tx {
	return nil
}

func (b *OdysseyAbortBlock) Visit(v Visitor) error {
	return v.OdysseyAbortBlock(b)
}

// NewOdysseyAbortBlock is kept for testing purposes only.
// Following Banff activation and subsequent code cleanup, Odyssey Abort blocks
// should be only verified (upon bootstrap), never created anymore
func NewOdysseyAbortBlock(
	parentID ids.ID,
	height uint64,
) (*OdysseyAbortBlock, error) {
	blk := &OdysseyAbortBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
	}
	return blk, initialize(blk)
}
