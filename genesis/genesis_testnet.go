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
	//go:embed genesis_testnet.json
	testnetGenesisConfigJSON []byte

	// TestnetParams are the params used for the testnet testnet
	TestnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         5 * units.Dione,
			CreateAssetTxFee:              10 * units.Dione,
			CreateSubnetTxFee:             300 * units.Dione,
			TransformSubnetTxFee:          5 * units.Dione,
			CreateBlockchainTxFee:         300 * units.Dione,
			AddPrimaryNetworkValidatorFee: 0,
			AddPrimaryNetworkDelegatorFee: 0,
			AddSubnetValidatorFee:         5 * units.Dione,
			AddSubnetDelegatorFee:         5 * units.Dione,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement:         .8, // 80%
			MinValidatorStake:         500 * units.KiloDione,
			MinDelegatorStake:         500 * units.Dione,
			MaxValidatorStake:         60 * units.MegaDione,
			MinDelegationFee:          20000, // 2%
			MinValidatorStakeDuration: 365 * 24 * time.Hour,
			MaxValidatorStakeDuration: 6 * 365 * 24 * time.Hour,
			MinDelegatorStakeDuration: 30 * 24 * time.Hour,
			MaxDelegatorStakeDuration: 6 * 365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          720 * units.MegaDione,
			},
			MintConfig: reward.MintConfig{
				MintingPeriod: 365 * 24 * time.Hour,
				MaxMintAmount: 500 * units.MegaDione,
				MintRate:      40_000, // 4%
			},
		},
	}
)
