// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package poll

import (
	"fmt"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/utils/bag"
	"github.com/dioneprotocol/dionego/utils/formatting"
)

// Set is a collection of polls
type Set interface {
	fmt.Stringer

	Add(requestID uint32, vdrs bag.Bag[ids.NodeID]) bool
	Vote(requestID uint32, vdr ids.NodeID, votes []ids.ID) []bag.UniqueBag[ids.ID]
	Len() int
}

// Poll is an outstanding poll
type Poll interface {
	formatting.PrefixedStringer

	Vote(vdr ids.NodeID, votes []ids.ID)
	Finished() bool
	Result() bag.UniqueBag[ids.ID]
}

// Factory creates a new Poll
type Factory interface {
	New(vdrs bag.Bag[ids.NodeID]) Poll
}
