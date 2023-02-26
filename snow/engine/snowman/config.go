// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/snow/consensus/snowball"
	"github.com/dioneprotocol/dionego/snow/consensus/snowman"
	"github.com/dioneprotocol/dionego/snow/engine/common"
	"github.com/dioneprotocol/dionego/snow/engine/snowman/block"
	"github.com/dioneprotocol/dionego/snow/validators"
)

// Config wraps all the parameters needed for a snowman engine
type Config struct {
	common.AllGetsServer

	Ctx        *snow.ConsensusContext
	VM         block.ChainVM
	Sender     common.Sender
	Validators validators.Set
	Params     snowball.Parameters
	Consensus  snowman.Consensus
}
