// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package dione

import (
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/snow/consensus/dione"
	"github.com/dioneprotocol/dionego/snow/engine/dione/vertex"
	"github.com/dioneprotocol/dionego/snow/engine/common"
	"github.com/dioneprotocol/dionego/snow/validators"
)

// Config wraps all the parameters needed for an dione engine
type Config struct {
	Ctx *snow.ConsensusContext
	common.AllGetsServer
	VM         vertex.DAGVM
	Manager    vertex.Manager
	Sender     common.Sender
	Validators validators.Set

	Params    dione.Parameters
	Consensus dione.Consensus
}
