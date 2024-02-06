// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package reward

import (
	"math/big"
	"time"
)

type mintCalculator struct {
	mintFrom   uint64
	mintUntil  uint64
	mintPeriod *big.Int
	mintAmount *big.Int
}

func NewMintCalculator(config MintConfig) *mintCalculator {
	return &mintCalculator{
		mintFrom:   config.MintSince,
		mintUntil:  config.MintUntil,
		mintPeriod: new(big.Int).SetUint64(config.MintUntil - config.MintSince),
		mintAmount: new(big.Int).SetUint64(config.MintAmount),
	}
}

func (c *mintCalculator) CalculateMintRate(validatorsAmount uint64, lastSyncTime, newChainTime time.Time) uint64 {
	if validatorsAmount == 0 {
		return 0
	}

	lastSyncTimeUnix := uint64(lastSyncTime.Unix())
	if lastSyncTimeUnix < c.mintFrom {
		lastSyncTimeUnix = c.mintFrom
	}

	newChainTimeUnix := uint64(newChainTime.Unix())
	if newChainTimeUnix > c.mintUntil {
		newChainTimeUnix = c.mintUntil
	}

	if newChainTimeUnix <= lastSyncTimeUnix {
		return 0
	}

	if lastSyncTimeUnix >= c.mintUntil {
		return 0
	}

	elapsed := new(big.Int).SetUint64(newChainTimeUnix - lastSyncTimeUnix)
	validatorsAmountBigInt := new(big.Int).SetUint64(validatorsAmount)

	result := elapsed
	result.Mul(result, c.mintAmount)
	result.Div(result, c.mintPeriod)
	result.Div(result, validatorsAmountBigInt)

	return result.Uint64()
}
