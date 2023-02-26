// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/dioneprotocol/dionego/utils/units"
	"github.com/dioneprotocol/dionego/vms/platformvm/reward"
)

var (
	//go:embed genesis_mainnet.json
	mainnetGenesisConfigJSON []byte

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         units.MilliDione,
			CreateAssetTxFee:              10 * units.MilliDione,
			CreateSubnetTxFee:             1 * units.Dione,
			TransformSubnetTxFee:          10 * units.Dione,
			CreateBlockchainTxFee:         1 * units.Dione,
			AddPrimaryNetworkValidatorFee: 0,
			AddPrimaryNetworkDelegatorFee: 0,
			AddSubnetValidatorFee:         units.MilliDione,
			AddSubnetDelegatorFee:         units.MilliDione,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloDione,
			MaxValidatorStake: 3 * units.MegaDione,
			MinDelegatorStake: 25 * units.Dione,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  2 * 7 * 24 * time.Hour,
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
