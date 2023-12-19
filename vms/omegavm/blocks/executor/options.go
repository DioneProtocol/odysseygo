// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"fmt"

	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
)

var _ blocks.Visitor = (*verifier)(nil)

// options supports build new option blocks
type options struct {
	// outputs populated by this struct's methods:
	commitBlock blocks.Block
	abortBlock  blocks.Block
}

func (*options) BanffAbortBlock(*blocks.BanffAbortBlock) error {
	return snowman.ErrNotOracle
}

func (*options) BanffCommitBlock(*blocks.BanffCommitBlock) error {
	return snowman.ErrNotOracle
}

func (o *options) BanffProposalBlock(b *blocks.BanffProposalBlock) error {
	timestamp := b.Timestamp()
	blkID := b.ID()
	nextHeight := b.Height() + 1

	var err error
	o.commitBlock, err = blocks.NewBanffCommitBlock(timestamp, blkID, nextHeight)
	if err != nil {
		return fmt.Errorf(
			"failed to create commit block: %w",
			err,
		)
	}

	o.abortBlock, err = blocks.NewBanffAbortBlock(timestamp, blkID, nextHeight)
	if err != nil {
		return fmt.Errorf(
			"failed to create abort block: %w",
			err,
		)
	}
	return nil
}

func (*options) BanffStandardBlock(*blocks.BanffStandardBlock) error {
	return snowman.ErrNotOracle
}

func (*options) OdysseyAbortBlock(*blocks.OdysseyAbortBlock) error {
	return snowman.ErrNotOracle
}

func (*options) OdysseyCommitBlock(*blocks.OdysseyCommitBlock) error {
	return snowman.ErrNotOracle
}

func (o *options) OdysseyProposalBlock(b *blocks.OdysseyProposalBlock) error {
	blkID := b.ID()
	nextHeight := b.Height() + 1

	var err error
	o.commitBlock, err = blocks.NewOdysseyCommitBlock(blkID, nextHeight)
	if err != nil {
		return fmt.Errorf(
			"failed to create commit block: %w",
			err,
		)
	}

	o.abortBlock, err = blocks.NewOdysseyAbortBlock(blkID, nextHeight)
	if err != nil {
		return fmt.Errorf(
			"failed to create abort block: %w",
			err,
		)
	}
	return nil
}

func (*options) OdysseyStandardBlock(*blocks.OdysseyStandardBlock) error {
	return snowman.ErrNotOracle
}

func (*options) OdysseyAtomicBlock(*blocks.OdysseyAtomicBlock) error {
	return snowman.ErrNotOracle
}
