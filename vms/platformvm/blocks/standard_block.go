// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
)

var (
	_ BanffBlock = (*BanffStandardBlock)(nil)
	_ Block      = (*ApricotStandardBlock)(nil)
)

type BanffStandardBlock struct {
	Time                 uint64 `serialize:"true" json:"time"`
	ApricotStandardBlock `serialize:"true"`
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
		ApricotStandardBlock: ApricotStandardBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			AccFee:       []byte{},
			FeeXChain:    []byte{},
			FeeCChain:    []byte{},
			Transactions: txs,
		},
	}
	return blk, initialize(blk)
}

func NewBanffStandardBlockWithFee(
	timestamp time.Time,
	parentID ids.ID,
	height uint64,
	txs []*txs.Tx,
	accumulatedFee *big.Int,
	feeFromXChain *big.Int,
	feeFromCChain *big.Int,
) (*BanffStandardBlock, error) {
	blk := &BanffStandardBlock{
		Time: uint64(timestamp.Unix()),
		ApricotStandardBlock: ApricotStandardBlock{
			CommonBlock: CommonBlock{
				PrntID: parentID,
				Hght:   height,
			},
			AccFee:       accumulatedFee.Bytes(),
			FeeXChain:    feeFromXChain.Bytes(),
			FeeCChain:    feeFromCChain.Bytes(),
			Transactions: txs,
		},
	}
	return blk, initialize(blk)
}

type ApricotStandardBlock struct {
	CommonBlock  `serialize:"true"`
	Transactions []*txs.Tx `serialize:"true" json:"txs"`
	AccFee       []byte    `serialize:"true" json:"accumulatedFee"`
	FeeXChain    []byte    `serialize:"true" json:"feeFromXChain"`
	FeeCChain    []byte    `serialize:"true" json:"feeFromCChain"`
}

func (b *ApricotStandardBlock) initialize(bytes []byte) error {
	b.CommonBlock.initialize(bytes)
	for _, tx := range b.Transactions {
		if err := tx.Initialize(txs.Codec); err != nil {
			return fmt.Errorf("failed to sign block: %w", err)
		}
	}
	return nil
}

func (b *ApricotStandardBlock) InitCtx(ctx *snow.Context) {
	for _, tx := range b.Transactions {
		tx.Unsigned.InitCtx(ctx)
	}
}

func (b *ApricotStandardBlock) Txs() []*txs.Tx {
	return b.Transactions
}

func (b *ApricotStandardBlock) Visit(v Visitor) error {
	return v.ApricotStandardBlock(b)
}

func (b *ApricotStandardBlock) FeeFromXChain() *big.Int {
	return new(big.Int).SetBytes(b.FeeXChain)
}

func (b *ApricotStandardBlock) FeeFromCChain() *big.Int {
	return new(big.Int).SetBytes(b.FeeCChain)
}

func (b *ApricotStandardBlock) AccumulatedFee() *big.Int {
	return new(big.Int).SetBytes(b.AccFee)
}

// NewApricotStandardBlock is kept for testing purposes only.
// Following Banff activation and subsequent code cleanup, Apricot Standard blocks
// should be only verified (upon bootstrap), never created anymore
func NewApricotStandardBlock(
	parentID ids.ID,
	height uint64,
	txs []*txs.Tx,
) (*ApricotStandardBlock, error) {
	blk := &ApricotStandardBlock{
		CommonBlock: CommonBlock{
			PrntID: parentID,
			Hght:   height,
		},
		AccFee:       []byte{},
		FeeXChain:    []byte{},
		FeeCChain:    []byte{},
		Transactions: txs,
	}
	return blk, initialize(blk)
}
