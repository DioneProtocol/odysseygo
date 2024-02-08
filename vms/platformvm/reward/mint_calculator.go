// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
	"time"
)

var epochTime = time.Unix(0, 0)

type mintCalculator struct {
	mintFrom   time.Time
	mintUntil  time.Time
	mintPeriod *big.Int
	mintAmount *big.Int
}

func NewMintCalculator(config MintConfig) *mintCalculator {
	return &mintCalculator{
		mintFrom:   time.Unix(config.MintSince, 0),
		mintUntil:  time.Unix(config.MintUntil, 0),
		mintPeriod: new(big.Int).SetInt64(config.MintUntil - config.MintSince),
		mintAmount: new(big.Int).SetUint64(config.MintAmount),
	}
}

func (c *mintCalculator) CalculateMintRate(validatorsAmount uint64, lastSyncTime, newChainTime time.Time) uint64 {
	if validatorsAmount == 0 {
		return 0
	}

	if lastSyncTime.Compare(c.mintFrom) < 0 {
		lastSyncTime = c.mintFrom
	}

	if newChainTime.Compare(c.mintUntil) > 0 {
		newChainTime = c.mintUntil
	}

	lastSyncTimeUnix := lastSyncTime.Unix()
	newChainTimeUnix := newChainTime.Unix()
	if newChainTimeUnix <= lastSyncTimeUnix {
		return 0
	}

	elapsed := new(big.Int).SetInt64(newChainTimeUnix - lastSyncTimeUnix)
	validatorsAmountBigInt := new(big.Int).SetUint64(validatorsAmount)

	result := elapsed
	result.Mul(result, c.mintAmount)
	result.Div(result, c.mintPeriod)
	result.Div(result, validatorsAmountBigInt)

	return result.Uint64()
}
