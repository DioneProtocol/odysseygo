// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package odyssey

import (
	"context"

	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowstorm"
)

// Vertex is a collection of multiple transactions tied to other vertices
//
// Note: Verify is not part of this interface because bootstrapping uses IDs to
// verify the vertex is valid.
type Vertex interface {
	choices.Decidable

	// Returns the vertices this vertex depends on
	Parents() ([]Vertex, error)

	// Returns the height of this vertex. A vertex's height is defined by one
	// greater than the maximum height of the parents.
	Height() (uint64, error)

	// Returns a series of state transitions to be performed on acceptance
	Txs(context.Context) ([]snowstorm.Tx, error)

	// Returns the binary representation of this vertex
	Bytes() []byte
}
