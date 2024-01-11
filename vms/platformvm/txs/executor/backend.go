// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/uptime"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/timer/mockable"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/config"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/fx"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/reward"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/utxo"
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
