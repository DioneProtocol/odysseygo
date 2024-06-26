// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
)

type StakingConfig struct {
	// Staking uptime requirements
	UptimeRequirement float64 `json:"uptimeRequirement"`
	// Minimum stake, in nDIONE, required to validate the primary network
	MinValidatorStake uint64 `json:"minValidatorStake"`
	// Maximum stake, in nDIONE, allowed to be placed on a single validator in
	// the primary network
	MaxValidatorStake uint64 `json:"maxValidatorStake"`
	// Minimum stake, in nDIONE, that can be delegated on the primary network
	MinDelegatorStake uint64 `json:"minDelegatorStake"`
	// Minimum delegation fee, in the range [0, 1000000], that can be charged
	// for delegation on the primary network.
	MinDelegationFee uint32 `json:"minDelegationFee"`
	// MinValidatorStakeDuration is the minimum amount of time a validator can validate
	// for in a single period.
	MinValidatorStakeDuration time.Duration `json:"minValidatorStakeDuration"`
	// MaxValidatorStakeDuration is the maximum amount of time a validator can validate
	// for in a single period.
	MaxValidatorStakeDuration time.Duration `json:"maxValidatorStakeDuration"`
	// MinDelegatorStakeDuration is the minimum amount of time a delegator can delegate
	// for in a single period.
	MinDelegatorStakeDuration time.Duration `json:"minDelekatorStakeDuration"`
	// MaxDelegatorStakeDuration is the maximum amount of time a delegator can delegate
	// for in a single period.
	MaxDelegatorStakeDuration time.Duration `json:"maxDelegatorStakeDuration"`
	// RewardConfig is the config for the reward function.
	RewardConfig reward.Config `json:"rewardConfig"`
	// Config for the minting function
	MintConfig reward.MintConfig `json:"mintConfig"`
}

type TxFeeConfig struct {
	// Transaction fee
	TxFee uint64 `json:"txFee"`
	// Transaction fee for create asset transactions
	CreateAssetTxFee uint64 `json:"createAssetTxFee"`
	// Transaction fee for create subnet transactions
	CreateSubnetTxFee uint64 `json:"createSubnetTxFee"`
	// Transaction fee for transform subnet transactions
	TransformSubnetTxFee uint64 `json:"transformSubnetTxFee"`
	// Transaction fee for create blockchain transactions
	CreateBlockchainTxFee uint64 `json:"createBlockchainTxFee"`
	// Transaction fee for adding a primary network validator
	AddPrimaryNetworkValidatorFee uint64 `json:"addPrimaryNetworkValidatorFee"`
	// Transaction fee for adding a primary network delegator
	AddPrimaryNetworkDelegatorFee uint64 `json:"addPrimaryNetworkDelegatorFee"`
	// Transaction fee for adding a subnet validator
	AddSubnetValidatorFee uint64 `json:"addSubnetValidatorFee"`
	// Transaction fee for adding a subnet delegator
	AddSubnetDelegatorFee uint64 `json:"addSubnetDelegatorFee"`
}

type Params struct {
	StakingConfig
	TxFeeConfig
}

func GetTxFeeConfig(networkID uint32) TxFeeConfig {
	switch networkID {
	case constants.MainnetID:
		return MainnetParams.TxFeeConfig
	case constants.TestnetID:
		return TestnetParams.TxFeeConfig
	case constants.LocalID:
		return LocalParams.TxFeeConfig
	default:
		return LocalParams.TxFeeConfig
	}
}

func GetStakingConfig(networkID uint32) StakingConfig {
	switch networkID {
	case constants.MainnetID:
		return MainnetParams.StakingConfig
	case constants.TestnetID:
		return TestnetParams.StakingConfig
	case constants.LocalID:
		return LocalParams.StakingConfig
	default:
		return LocalParams.StakingConfig
	}
}
