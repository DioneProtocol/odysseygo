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
	mintSince          int64
	mintPeriod         int64
	mintPeriodBigInt   *big.Int
	maxMintAmount      *big.Int
	mintRate           *big.Int
	initialSupply      *big.Int
	percentDenominator *big.Int
}

func NewMintCalculator(config MintConfig, initialSupply uint64) *mintCalculator {
	mintPeriod := uint64(config.MintingPeriod.Seconds())
	return &mintCalculator{
		mintPeriod:         int64(mintPeriod),
		mintSince:          config.MintSince,
		mintPeriodBigInt:   new(big.Int).SetUint64(mintPeriod),
		maxMintAmount:      new(big.Int).SetUint64(config.MaxMintAmount),
		mintRate:           new(big.Int).SetUint64(config.MintRate),
		initialSupply:      new(big.Int).SetUint64(initialSupply),
		percentDenominator: new(big.Int).SetUint64(PercentDenominator),
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
	lastSyncTimeUnix := lastSyncTime.Unix()
	if lastSyncTimeUnix < c.mintSince {
		lastSyncTimeUnix = c.mintSince
	}

	newChainTimeUnix := newChainTime.Unix()
	if newChainTimeUnix <= lastSyncTimeUnix {
		return new(big.Int)
	}

	lastSyncTimePeriod := (lastSyncTimeUnix - c.mintSince) / c.mintPeriod
	newChainTimePeriod := (newChainTimeUnix - c.mintSince) / c.mintPeriod

	if lastSyncTimePeriod != newChainTimePeriod {
		newPeriodTimestamp := c.mintSince + c.mintPeriod*lastSyncTimePeriod + c.mintPeriod
		if newPeriodTimestamp != newChainTimeUnix {
			mintRateBeforeNewPeriod := c.CalculateMintRate(totalWeight, lastSyncTime, time.Unix(newPeriodTimestamp, 0))
			mintRateAfterNewPeriod := c.CalculateMintRate(totalWeight, time.Unix(newPeriodTimestamp, 0), newChainTime)
			return new(big.Int).Add(mintRateBeforeNewPeriod, mintRateAfterNewPeriod)
		}
	}

	elapsed := new(big.Int).SetInt64(newChainTimeUnix - lastSyncTimeUnix)
	totalWeightBigInt := new(big.Int).SetUint64(totalWeight)

	period := lastSyncTimePeriod

	periodBigInt := new(big.Int).SetInt64(period)
	startPeriodSupply := new(big.Int).Set(c.percentDenominator)
	startPeriodSupply.Add(startPeriodSupply, c.mintRate)
	startPeriodSupply.Exp(startPeriodSupply, periodBigInt, nil)
	startPeriodSupply.Mul(startPeriodSupply, c.initialSupply)

	supplyDenominator := new(big.Int).Set(c.percentDenominator)
	supplyDenominator.Exp(supplyDenominator, periodBigInt, nil)

	startPeriodSupply.Div(startPeriodSupply, supplyDenominator)

	mintAmount := new(big.Int).Set(startPeriodSupply)
	mintAmount.Mul(mintAmount, c.mintRate)
	mintAmount.Div(mintAmount, c.percentDenominator)

	if mintAmount.Cmp(c.maxMintAmount) > 0 {
		mintAmount.Set(c.maxMintAmount)
	}

	result := elapsed
	result.Mul(result, mintAmount)
	result.Lsh(result, mintShift)
	result.Div(result, c.mintPeriodBigInt)
	result.Div(result, totalWeightBigInt)

	return result
}
