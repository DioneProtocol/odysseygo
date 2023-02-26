// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package dione

import (
	"context"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow/consensus/dione"
	"github.com/dioneprotocol/dionego/snow/engine/common"
)

// Engine describes the events that can occur on a consensus instance
type Engine interface {
	common.Engine

	// GetVtx returns a vertex by its ID.
	// Returns an error if unknown.
	GetVtx(ctx context.Context, vtxID ids.ID) (dione.Vertex, error)
}
