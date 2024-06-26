// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/metrics"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs/executor"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs/mempool"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/validators"
)

var _ Manager = (*manager)(nil)

type Manager interface {
	state.Versions

	// Returns the ID of the most recently accepted block.
	LastAccepted() ids.ID
	GetBlock(blkID ids.ID) (snowman.Block, error)
	GetStatelessBlock(blkID ids.ID) (blocks.Block, error)
	NewBlock(blocks.Block) snowman.Block
}

func NewManager(
	mempool mempool.Mempool,
	metrics metrics.Metrics,
	s state.State,
	txExecutorBackend *executor.Backend,
	validatorManager validators.Manager,
) Manager {
	backend := &backend{
		Mempool:      mempool,
		lastAccepted: s.GetLastAccepted(),
		state:        s,
		ctx:          txExecutorBackend.Ctx,
		blkIDToState: map[ids.ID]*blockState{},
	}

	return &manager{
		backend: backend,
		verifier: &verifier{
			backend:           backend,
			txExecutorBackend: txExecutorBackend,
		},
		acceptor: &acceptor{
			backend:      backend,
			metrics:      metrics,
			validators:   validatorManager,
			bootstrapped: txExecutorBackend.Bootstrapped,
		},
		rejector: &rejector{
			backend:         backend,
			addTxsToMempool: !txExecutorBackend.Config.PartialSyncPrimaryNetwork,
		},
	}
}

type manager struct {
	*backend
	verifier blocks.Visitor
	acceptor blocks.Visitor
	rejector blocks.Visitor
}

func (m *manager) GetBlock(blkID ids.ID) (snowman.Block, error) {
	blk, err := m.backend.GetBlock(blkID)
	if err != nil {
		return nil, err
	}
	return m.NewBlock(blk), nil
}

func (m *manager) GetStatelessBlock(blkID ids.ID) (blocks.Block, error) {
	return m.backend.GetBlock(blkID)
}

func (m *manager) NewBlock(blk blocks.Block) snowman.Block {
	return &Block{
		manager: m,
		Block:   blk,
	}
}
