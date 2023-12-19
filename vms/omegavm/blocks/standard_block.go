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
	_ BanffBlock = (*BanffStandardBlock)(nil)
	_ Block      = (*OdysseyStandardBlock)(nil)
)

type BanffStandardBlock struct {
	Time                 uint64 `serialize:"true" json:"time"`
	OdysseyStandardBlock `serialize:"true"`
}

func (b *BanffStandardBlock) Timestamp() time.Time {
	return time.Unix(int64(b.Time), 0)
}

func (b *BanffStandardBlock) Visit(v Visitor) error {
	return v.BanffStandardBlock(b)
}

func NewBanffStandardBlock(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
	txs []*txs.Tx,
) (*BanffStandardBlock, error) {
	blk := &BanffStandardBlock{
		Time: uint64(timestamp.Unix()),
		OdysseyStandardBlock: OdysseyStandardBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			Transactions: txs,
		},
	}
	return blk, initialize(blk)
}

type OdysseyStandardBlock struct {
	CommonBlock  `serialize:"true"`
	Transactions []*txs.Tx `serialize:"true" json:"txs"`
}

func (b *OdysseyStandardBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	for _, tx := range b.Transactions {
		if err := tx.Initialize(txs.Codec); err != nil {
			return fmt.Errorf("failed to sign block: %w", err)
		}
	}
	return nil
}

func (b *OdysseyStandardBlock) InitCtx(ctx *snow.Context) {
	for _, tx := range b.Transactions {
		tx.Unsigned.InitCtx(ctx)
	}
}

func (b *OdysseyStandardBlock) Txs() []*txs.Tx {
	return b.Transactions
}

func (b *OdysseyStandardBlock) Visit(v Visitor) error {
	return v.OdysseyStandardBlock(b)
}

// NewOdysseyStandardBlock is kept for testing purposes only.
// Following Banff activation and subsequent code cleanup, Odyssey Standard blocks
// should be only verified (upon bootstrap), never created anymore
func NewOdysseyStandardBlock(
	parentID ids.ID,
	height uint64,
	txs []*txs.Tx,
) (*OdysseyStandardBlock, error) {
	blk := &OdysseyStandardBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		Transactions: txs,
	}
	return blk, initialize(blk)
}
