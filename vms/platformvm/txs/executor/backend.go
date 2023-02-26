// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/snow/uptime"
	"github.com/dioneprotocol/dionego/utils"
	"github.com/dioneprotocol/dionego/utils/timer/mockable"
	"github.com/dioneprotocol/dionego/vms/platformvm/config"
	"github.com/dioneprotocol/dionego/vms/platformvm/fx"
	"github.com/dioneprotocol/dionego/vms/platformvm/reward"
	"github.com/dioneprotocol/dionego/vms/platformvm/utxo"
)

type Backend struct {
	Config       *config.Config
	Ctx          *snow.Context
	Clk          *mockable.Clock
	Fx           fx.Fx
	FlowChecker  utxo.Verifier
	Uptimes      uptime.Manager
	Rewards      reward.Calculator
	Bootstrapped *utils.Atomic[bool]
}
