// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import (
	"context"

	"github.com/DioneProtocol/odysseygo/database/versiondb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
)

// BlockServer represents all requests heightIndexer can issue
// against ProposerVM. All methods must be thread-safe.
type BlockServer interface {
	versiondb.Commitable

	// Note: this is a contention heavy call that should be avoided
	// for frequent/repeated indexer ops
	GetFullPostForkBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error)
}
