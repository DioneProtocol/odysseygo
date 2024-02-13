// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
)

var _ DistrubuteCalculator = (*distrubuteCalculator)(nil)

type DistrubuteCalculator interface {
	Calculate(weight uint64, feePerWeightPaid *big.Int) uint64
}

type distrubuteCalculator struct {
	feePerWeightStored *big.Int
}

func NewDistrubuteCalculator(feePerWeightStored *big.Int) DistrubuteCalculator {
	return &distrubuteCalculator{
		feePerWeightStored: feePerWeightStored,
	}
}

func (dc *distrubuteCalculator) Calculate(weight uint64, feePerWeightPaid *big.Int) uint64 {
	feePerWeightDiff := new(big.Int).Set(dc.feePerWeightStored)
	feePerWeightDiff.Sub(feePerWeightDiff, feePerWeightPaid)
	potentialReward := new(big.Int).SetUint64(weight)
	potentialReward.Mul(potentialReward, feePerWeightDiff)
	potentialReward.Div(potentialReward, BigFeePerWeightDenominator)

	return potentialReward.Uint64()
}
