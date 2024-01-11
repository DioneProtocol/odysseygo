// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowball"
	"github.com/DioneProtocol/odysseygo/snow/consensus/snowman"
	"github.com/DioneProtocol/odysseygo/snow/engine/common"
	"github.com/DioneProtocol/odysseygo/snow/engine/snowman/block"
	"github.com/DioneProtocol/odysseygo/snow/validators"
)

// Config wraps all the parameters needed for a snowman engine
type Config struct {
	common.AllGetsServer

	Ctx         *snow.ConsensusContext
	VM          block.ChainVM
	Sender      common.Sender
	Validators  validators.Set
	Params      snowball.Parameters
	Consensus   snowman.Consensus
	PartialSync bool
}
