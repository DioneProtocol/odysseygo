// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
	"time"
)

// PercentDenominator is the denominator used to calculate percentages
const PercentDenominator = 1_000_000

// consumptionRateDenominator is the magnitude offset used to emulate
// floating point fractions.
var consumptionRateDenominator = new(big.Int).SetUint64(PercentDenominator)

type Config struct {
	// MaxConsumptionRate is the rate to allocate funds if the validator's stake
	// duration is equal to [MintingPeriod]
	MaxConsumptionRate uint64 `json:"maxConsumptionRate"`

	// MinConsumptionRate is the rate to allocate funds if the validator's stake
	// duration is 0.
	MinConsumptionRate uint64 `json:"minConsumptionRate"`

	// MintingPeriod is period that the staking calculator runs on. It is
	// not valid for a validator's stake duration to be larger than this.
	MintingPeriod time.Duration `json:"mintingPeriod"`

	// SupplyCap is the target value that the reward calculation should be
	// asymptotic to.
	SupplyCap uint64 `json:"supplyCap"`
}

type MintConfig struct {
	// MintSince is the Unix Epoch timestamp since which the reward will be
	// minted
	MintSince uint64 `json:"mintSince"`

	// MintUntil is the Unix Epoch timestamp until which the reward will be
	// minted
	MintUntil uint64 `json:"mintUntil"`

	// SupplyCap is the target value that the reward calculation should be
	// asymptotic to.
	MintAmount uint64 `json:"mintAmount"`
}
