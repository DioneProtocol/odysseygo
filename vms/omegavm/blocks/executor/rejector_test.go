// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs/mempool"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

func TestRejectBlock(t *testing.T) {
	type test struct {
		name         string
		newBlockFunc func() (blocks.Block, error)
		rejectFunc   func(*rejector, blocks.Block) error
	}

	tests := []test{
		{
			name: "proposal block",
			newBlockFunc: func() (blocks.Block, error) {
				return blocks.NewBanffProposalBlock(
					time.Now(),
					ids.GenerateTestID(),
					1,
				)
			},
			rejectFunc: func(r *rejector, b blocks.Block) error {
				return r.BanffProposalBlock(b.(*blocks.BanffProposalBlock))
			},
		},
		{
			name: "atomic block",
			newBlockFunc: func() (blocks.Block, error) {
				return blocks.NewOdysseyAtomicBlock(
					ids.GenerateTestID(),
					1,
				)
			},
			rejectFunc: func(r *rejector, b blocks.Block) error {
				return r.OdysseyAtomicBlock(b.(*blocks.OdysseyAtomicBlock))
			},
		},
		{
			name: "standard block",
			newBlockFunc: func() (blocks.Block, error) {
				return blocks.NewBanffStandardBlock(
					time.Now(),
					ids.GenerateTestID(),
					1,
				)
			},
			rejectFunc: func(r *rejector, b blocks.Block) error {
				return r.BanffStandardBlock(b.(*blocks.BanffStandardBlock))
			},
		},
		{
			name: "commit",
			newBlockFunc: func() (blocks.Block, error) {
				return blocks.NewBanffCommitBlock(time.Now(), ids.GenerateTestID() /*parent*/, 1 /*height*/)
			},
			rejectFunc: func(r *rejector, blk blocks.Block) error {
				return r.BanffCommitBlock(blk.(*blocks.BanffCommitBlock))
			},
		},
		{
			name: "abort",
			newBlockFunc: func() (blocks.Block, error) {
				return blocks.NewBanffAbortBlock(time.Now(), ids.GenerateTestID() /*parent*/, 1 /*height*/)
			},
			rejectFunc: func(r *rejector, blk blocks.Block) error {
				return r.BanffAbortBlock(blk.(*blocks.BanffAbortBlock))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blk, err := tt.newBlockFunc()
			require.NoError(err)

			mempool := mempool.NewMockMempool(ctrl)
			state := state.NewMockState(ctrl)
			blkIDToState := map[ids.ID]*blockState{
				blk.Parent(): nil,
				blk.ID():     nil,
			}
			rejector := &rejector{
				backend: &backend{
					ctx: &snow.Context{
						Log: logging.NoLog{},
					},
					blkIDToState: blkIDToState,
					Mempool:      mempool,
					state:        state,
				},
			}

			// Set expected calls on dependencies.
			for _, tx := range blk.Txs() {
				mempool.EXPECT().Add(tx).Return(nil).Times(1)
			}
			gomock.InOrder(
				state.EXPECT().AddStatelessBlock(blk, choices.Rejected).Times(1),
				state.EXPECT().Commit().Return(nil).Times(1),
			)

			err = tt.rejectFunc(rejector, blk)
			require.NoError(err)
			// Make sure block and its parent are removed from the state map.
			require.NotContains(rejector.blkIDToState, blk.ID())
		})
	}
}
