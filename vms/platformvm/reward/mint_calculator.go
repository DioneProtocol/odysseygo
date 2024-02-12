// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"fmt"
	"math/big"
	"time"
)

var (
	epochTime      = time.Unix(0, 0)
	mintShift uint = 64

	errInvalidMintPeriod = fmt.Errorf("mintFrom must be less than mintUntil")
)

type mintCalculator struct {
	mintFrom   time.Time
	mintUntil  time.Time
	mintPeriod *big.Int
	mintAmount *big.Int
}

func NewMintCalculator(config MintConfig) (*mintCalculator, error) {
	if config.MintUntil <= config.MintSince {
		return nil, errInvalidMintPeriod
	}
	return &mintCalculator{
		mintFrom:   time.Unix(config.MintSince, 0),
		mintUntil:  time.Unix(config.MintUntil, 0),
		mintPeriod: new(big.Int).SetInt64(config.MintUntil - config.MintSince),
		mintAmount: new(big.Int).SetUint64(config.MintAmount),
	}, nil
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
	if lastSyncTime.Compare(c.mintFrom) < 0 {
		lastSyncTime = c.mintFrom
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
