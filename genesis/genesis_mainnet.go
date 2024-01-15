// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
)

var (
	//go:embed genesis_mainnet.json
	mainnetGenesisConfigJSON []byte

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         50 * units.Dione,
			CreateAssetTxFee:              100 * units.Dione,
			CreateSubnetTxFee:             3 * units.KiloDione,
			TransformSubnetTxFee:          50 * units.Dione,
			CreateBlockchainTxFee:         3 * units.KiloDione,
			AddPrimaryNetworkValidatorFee: 0,
			AddPrimaryNetworkDelegatorFee: 0,
			AddSubnetValidatorFee:         50 * units.Dione,
			AddSubnetDelegatorFee:         50 * units.Dione,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 500 * units.KiloDione,
			MaxValidatorStake: 3 * units.MegaDione,
			MinDelegatorStake: 25 * units.Dione,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  30 * 24 * time.Hour,
			MaxStakeDuration:  6 * 365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          720 * units.MegaDione,
			},
		},
	}
)
