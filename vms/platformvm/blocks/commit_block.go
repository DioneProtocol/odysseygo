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
	_ BanffBlock = (*BanffCommitBlock)(nil)
	_ Block      = (*OdysseyCommitBlock)(nil)
)

type BanffCommitBlock struct {
	Time               uint64 `serialize:"true" json:"time"`
	OdysseyCommitBlock `serialize:"true"`
}

func (b *BanffCommitBlock) Timestamp() time.Time {
	return time.Unix(int64(b.Time), 0)
}

func (b *BanffCommitBlock) Visit(v Visitor) error {
	return v.BanffCommitBlock(b)
}

func NewBanffCommitBlock(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
) (*BanffCommitBlock, error) {
	blk := &BanffCommitBlock{
		Time: uint64(timestamp.Unix()),
		OdysseyCommitBlock: OdysseyCommitBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
		},
	}
	return blk, initialize(blk)
}

type OdysseyCommitBlock struct {
	CommonBlock `serialize:"true"`
}

func (b *OdysseyCommitBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	return nil
}

func (*OdysseyCommitBlock) InitCtx(*snow.Context) {}

func (*OdysseyCommitBlock) Txs() []*txs.Tx {
	return nil
}

func (b *OdysseyCommitBlock) Visit(v Visitor) error {
	return v.OdysseyCommitBlock(b)
}

func NewOdysseyCommitBlock(
	parentID ids.ID,
	height uint64,
) (*OdysseyCommitBlock, error) {
	blk := &OdysseyCommitBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
	}
	return blk, initialize(blk)
}
