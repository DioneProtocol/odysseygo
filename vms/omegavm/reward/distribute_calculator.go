// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
)

func CalculateFeeReward(feePerWeightStored *big.Int, weight uint64, feePerWeightPaid *big.Int) uint64 {
	feePerWeightDiff := new(big.Int).Set(feePerWeightStored)
	feePerWeightDiff.Sub(feePerWeightDiff, feePerWeightPaid)
	potentialReward := new(big.Int).SetUint64(weight)
	potentialReward.Mul(potentialReward, feePerWeightDiff)
	potentialReward.Rsh(potentialReward, BitShift)

	return potentialReward.Uint64()
}
