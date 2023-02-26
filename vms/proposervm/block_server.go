// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package proposervm

import (
	"context"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow/consensus/snowman"
	"github.com/dioneprotocol/dionego/vms/proposervm/indexer"
)

var _ indexer.BlockServer = (*VM)(nil)

// Note: this is a contention heavy call that should be avoided
// for frequent/repeated indexer ops
func (vm *VM) GetFullPostForkBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error) {
	vm.ctx.Lock.Lock()
	defer vm.ctx.Lock.Unlock()

	return vm.getPostForkBlock(ctx, blkID)
}

func (vm *VM) Commit() error {
	vm.ctx.Lock.Lock()
	defer vm.ctx.Lock.Unlock()

	return vm.db.Commit()
}
