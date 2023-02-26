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
	//go:embed genesis_fuji.json
	fujiGenesisConfigJSON []byte

	// FujiParams are the params used for the fuji testnet
	FujiParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         units.MilliDione,
			CreateAssetTxFee:              10 * units.MilliDione,
			CreateSubnetTxFee:             100 * units.MilliDione,
			TransformSubnetTxFee:          1 * units.Dione,
			CreateBlockchainTxFee:         100 * units.MilliDione,
			AddPrimaryNetworkValidatorFee: 0,
			AddPrimaryNetworkDelegatorFee: 0,
			AddSubnetValidatorFee:         units.MilliDione,
			AddSubnetDelegatorFee:         units.MilliDione,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 1 * units.Dione,
			MaxValidatorStake: 3 * units.MegaDione,
			MinDelegatorStake: 1 * units.Dione,
			MinDelegationFee:  20000, // 2%
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
