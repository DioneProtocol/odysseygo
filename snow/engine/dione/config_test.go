// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package dione

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dioneprotocol/dionego/database/memdb"
	"github.com/dioneprotocol/dionego/snow/consensus/dione"
	"github.com/dioneprotocol/dionego/snow/consensus/snowball"
	"github.com/dioneprotocol/dionego/snow/engine/dione/bootstrap"
	"github.com/dioneprotocol/dionego/snow/engine/dione/vertex"
	"github.com/dioneprotocol/dionego/snow/engine/common"
	"github.com/dioneprotocol/dionego/snow/engine/common/queue"
	"github.com/dioneprotocol/dionego/snow/validators"
)

func DefaultConfig() (common.Config, bootstrap.Config, Config) {
	vtxBlocked, _ := queue.NewWithMissing(memdb.New(), "", prometheus.NewRegistry())
	txBlocked, _ := queue.New(memdb.New(), "", prometheus.NewRegistry())

	commonCfg := common.DefaultConfigTest()

	bootstrapConfig := bootstrap.Config{
		Config:     commonCfg,
		VtxBlocked: vtxBlocked,
		TxBlocked:  txBlocked,
		Manager:    &vertex.TestManager{},
		VM:         &vertex.TestVM{},
	}

	engineConfig := Config{
		Ctx:        bootstrapConfig.Ctx,
		VM:         bootstrapConfig.VM,
		Manager:    bootstrapConfig.Manager,
		Sender:     bootstrapConfig.Sender,
		Validators: validators.NewSet(),
		Params: dione.Parameters{
			Parameters: snowball.Parameters{
				K:                       1,
				Alpha:                   1,
				BetaVirtuous:            1,
				BetaRogue:               2,
				ConcurrentRepolls:       1,
				OptimalProcessing:       100,
				MaxOutstandingItems:     1,
				MaxItemProcessingTime:   1,
				MixedQueryNumPushVdr:    1,
				MixedQueryNumPushNonVdr: 1,
			},
			Parents:   2,
			BatchSize: 1,
		},
		Consensus: &dione.Topological{},
	}

	return commonCfg, bootstrapConfig, engineConfig
}
