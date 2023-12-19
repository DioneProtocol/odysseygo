// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
)

func TestStatus(t *testing.T) {
	type test struct {
		name           string
		blockF         func(*gomock.Controller) *Block
		expectedStatus choices.Status
	}

	tests := []test{
		{
			name: "last accepted",
			blockF: func(ctrl *gomock.Controller) *Block {
				blkID := ids.GenerateTestID()
				statelessBlk := blocks.NewMockBlock(ctrl)
				statelessBlk.EXPECT().ID().Return(blkID)

				manager := &manager{
					backend: &backend{
						lastAccepted: blkID,
					},
				}

				return &Block{
					Block:   statelessBlk,
					manager: manager,
				}
			},
			expectedStatus: choices.Accepted,
		},
		{
			name: "processing",
			blockF: func(ctrl *gomock.Controller) *Block {
				blkID := ids.GenerateTestID()
				statelessBlk := blocks.NewMockBlock(ctrl)
				statelessBlk.EXPECT().ID().Return(blkID)

				manager := &manager{
					backend: &backend{
						blkIDToState: map[ids.ID]*blockState{
							blkID: {},
						},
					},
				}
				return &Block{
					Block:   statelessBlk,
					manager: manager,
				}
			},
			expectedStatus: choices.Processing,
		},
		{
			name: "in database",
			blockF: func(ctrl *gomock.Controller) *Block {
				blkID := ids.GenerateTestID()
				statelessBlk := blocks.NewMockBlock(ctrl)
				statelessBlk.EXPECT().ID().Return(blkID)

				state := state.NewMockState(ctrl)
				state.EXPECT().GetStatelessBlock(blkID).Return(statelessBlk, choices.Accepted, nil)

				manager := &manager{
					backend: &backend{
						state: state,
					},
				}
				return &Block{
					Block:   statelessBlk,
					manager: manager,
				}
			},
			expectedStatus: choices.Accepted,
		},
		{
			name: "not in map or database",
			blockF: func(ctrl *gomock.Controller) *Block {
				blkID := ids.GenerateTestID()
				statelessBlk := blocks.NewMockBlock(ctrl)
				statelessBlk.EXPECT().ID().Return(blkID)

				state := state.NewMockState(ctrl)
				state.EXPECT().GetStatelessBlock(blkID).Return(nil, choices.Unknown, database.ErrNotFound)

				manager := &manager{
					backend: &backend{
						state: state,
					},
				}
				return &Block{
					Block:   statelessBlk,
					manager: manager,
				}
			},
			expectedStatus: choices.Processing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blk := tt.blockF(ctrl)
			require.Equal(t, tt.expectedStatus, blk.Status())
		})
	}
}

func TestBlockOptions(t *testing.T) {
	type test struct {
		name                   string
		blkF                   func() *Block
		expectedPreferenceType blocks.Block
		expectedErr            error
	}

	tests := []test{
		{
			name: "odyssey proposal block; commit preferred",
			blkF: func() *Block {
				innerBlk := &blocks.OdysseyProposalBlock{}
				blkID := innerBlk.ID()

				manager := &manager{
					backend: &backend{
						blkIDToState: map[ids.ID]*blockState{
							blkID: {
								proposalBlockState: proposalBlockState{
									initiallyPreferCommit: true,
								},
							},
						},
					},
				}

				return &Block{
					Block:   innerBlk,
					manager: manager,
				}
			},
			expectedPreferenceType: &blocks.OdysseyCommitBlock{},
		},
		{
			name: "odyssey proposal block; abort preferred",
			blkF: func() *Block {
				innerBlk := &blocks.OdysseyProposalBlock{}
				blkID := innerBlk.ID()

				manager := &manager{
					backend: &backend{
						blkIDToState: map[ids.ID]*blockState{
							blkID: {},
						},
					},
				}

				return &Block{
					Block:   innerBlk,
					manager: manager,
				}
			},
			expectedPreferenceType: &blocks.OdysseyAbortBlock{},
		},
		{
			name: "banff proposal block; commit preferred",
			blkF: func() *Block {
				innerBlk := &blocks.BanffProposalBlock{}
				blkID := innerBlk.ID()

				manager := &manager{
					backend: &backend{
						blkIDToState: map[ids.ID]*blockState{
							blkID: {
								proposalBlockState: proposalBlockState{
									initiallyPreferCommit: true,
								},
							},
						},
					},
				}

				return &Block{
					Block:   innerBlk,
					manager: manager,
				}
			},
			expectedPreferenceType: &blocks.BanffCommitBlock{},
		},
		{
			name: "banff proposal block; abort preferred",
			blkF: func() *Block {
				innerBlk := &blocks.BanffProposalBlock{}
				blkID := innerBlk.ID()

				manager := &manager{
					backend: &backend{
						blkIDToState: map[ids.ID]*blockState{
							blkID: {},
						},
					},
				}

				return &Block{
					Block:   innerBlk,
					manager: manager,
				}
			},
			expectedPreferenceType: &blocks.BanffAbortBlock{},
		},
		{
			name: "non oracle block",
			blkF: func() *Block {
				return &Block{
					Block:   &blocks.BanffStandardBlock{},
					manager: &manager{},
				}
			},
			expectedErr: snowman.ErrNotOracle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blk := tt.blkF()
			options, err := blk.Options(context.Background())
			if tt.expectedErr != nil {
				require.ErrorIs(err, tt.expectedErr)
				return
			}
			require.IsType(tt.expectedPreferenceType, options[0].(*Block).Block)
		})
	}
}
