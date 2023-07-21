// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/reward"
)

var (
	//go:embed genesis_testnet.json
	testnetGenesisConfigJSON []byte

	// TestnetParams are the params used for the testnet
	TestnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         50 * units.MilliDione,
			CreateAssetTxFee:              100 * units.MilliDione,
			CreateSubnetTxFee:             1 * units.Dione,
			TransformSubnetTxFee:          10 * units.Dione,
			CreateBlockchainTxFee:         1 * units.Dione,
			AddPrimaryNetworkValidatorFee: 0,
			AddSubnetValidatorFee:         units.MilliDione,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 1 * units.Dione,
			MaxValidatorStake: 3 * units.MegaDione,
			MinStakeDuration:  24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          720 * units.MegaDione,
			},
		},
	}
)
