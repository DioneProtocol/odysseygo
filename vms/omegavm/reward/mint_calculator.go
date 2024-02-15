// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
	"time"
)

var (
	_ MintCalculator = (*mintCalculator)(nil)

	// 32 bits for unix time + 64 bits for a weight
	mintShift uint = 96
)

type MintCalculator interface {
	CalculateMintRate(totalWeight uint64, lastSyncTime, newChainTime time.Time) *big.Int
}

type mintCalculator struct {
	mintSince  time.Time
	mintUntil  time.Time
	mintPeriod *big.Int
	mintAmount *big.Int
}

func NewMintCalculator(config MintConfig) *mintCalculator {
	mintSince := time.Unix(config.MintSince, 0)
	mintUntil := mintSince.Add(config.MintingPeriod)
	mintPeriod := uint64(config.MintingPeriod.Seconds())
	return &mintCalculator{
		mintSince:  mintSince,
		mintUntil:  mintUntil,
		mintPeriod: new(big.Int).SetUint64(mintPeriod),
		mintAmount: new(big.Int).SetUint64(config.MintAmount),
	}
}

func CalculateMintReward(weight uint64, stakerMintRate, accumulatedMintRate *big.Int) uint64 {
	weightBigInt := new(big.Int).SetUint64(weight)
	result := new(big.Int).Set(accumulatedMintRate)
	result.Sub(result, stakerMintRate)
	result.Mul(result, weightBigInt)
	result.Rsh(result, mintShift)
	return result.Uint64()
}

func (c *mintCalculator) CalculateMintRate(totalWeight uint64, lastSyncTime, newChainTime time.Time) *big.Int {
	if lastSyncTime.Compare(c.mintSince) < 0 {
		lastSyncTime = c.mintSince
	}

	if newChainTime.Compare(c.mintUntil) > 0 {
		newChainTime = c.mintUntil
	}

	lastSyncTimeUnix := lastSyncTime.Unix()
	newChainTimeUnix := newChainTime.Unix()
	if newChainTimeUnix <= lastSyncTimeUnix {
		return new(big.Int)
	}

	elapsed := new(big.Int).SetInt64(newChainTimeUnix - lastSyncTimeUnix)
	totalWeightBigInt := new(big.Int).SetUint64(totalWeight)

	result := elapsed
	result.Mul(result, c.mintAmount)
	result.Lsh(result, mintShift)
	result.Div(result, c.mintPeriod)
	result.Div(result, totalWeightBigInt)

	return result
}
