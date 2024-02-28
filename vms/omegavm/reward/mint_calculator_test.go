package reward

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMintWithMaxMintAmount(t *testing.T) {
	maxMintAmount := uint64(1_000_000)
	initialSupply := uint64(100_000_000)
	mintSince := int64(50)
	mintingPeriod := int64(100)

	tests := []struct {
		lastSyncTime       int64
		newChainTime       int64
		expectedMintAmount uint64
	}{
		{
			lastSyncTime:       0,
			newChainTime:       mintSince,
			expectedMintAmount: 0,
		},
		{
			lastSyncTime:       0,
			newChainTime:       mintSince + mintingPeriod,
			expectedMintAmount: maxMintAmount,
		},
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod,
			expectedMintAmount: maxMintAmount,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod,
			newChainTime:       mintSince + mintingPeriod*2,
			expectedMintAmount: maxMintAmount,
		},
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod*2,
			expectedMintAmount: maxMintAmount * 2,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod*3/4,
			newChainTime:       mintSince + mintingPeriod*7/4,
			expectedMintAmount: maxMintAmount,
		},
	}

	mintConfig := MintConfig{
		MintSince:     int64(mintSince),
		MintingPeriod: time.Duration(mintingPeriod * int64(time.Second)),
		MintRate:      PercentDenominator,
		MaxMintAmount: maxMintAmount,
	}

	c := NewMintCalculator(mintConfig, initialSupply)
	for totalWeight := uint64(1); totalWeight < 10; totalWeight++ {
		for weight := uint64(0); weight <= totalWeight; weight++ {
			for _, test := range tests {
				expectedReward := test.expectedMintAmount * weight / totalWeight
				name := fmt.Sprintf("mint(%d,%d,%d,%d)==%d",
					weight,
					totalWeight,
					test.lastSyncTime,
					test.newChainTime,
					expectedReward,
				)
				t.Run(name, func(t *testing.T) {
					mintRate := c.CalculateMintRate(
						totalWeight,
						time.Unix(int64(test.lastSyncTime), 0),
						time.Unix(int64(test.newChainTime), 0),
					)
					reward := CalculateMintReward(weight, new(big.Int), mintRate)

					// might happen roundoff error
					require.True(t, expectedReward-reward <= 1, "%d != %d", expectedReward, reward)
				})
			}
		}
	}

}

func TestMintWithInflationRate(t *testing.T) {
	maxMintAmount := uint64(1_000_000_000_000_000)
	initialSupply := uint64(1_000_000)
	mintRate := uint64(PercentDenominator * 0.1)

	mintSince := int64(50)
	mintingPeriod := int64(100)

	tests := []struct {
		lastSyncTime       int64
		newChainTime       int64
		expectedMintAmount uint64
	}{
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod,
			expectedMintAmount: initialSupply / 10,
		},
		{
			lastSyncTime:       0,
			newChainTime:       mintSince + mintingPeriod,
			expectedMintAmount: initialSupply / 10,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod,
			newChainTime:       mintSince + mintingPeriod*2,
			expectedMintAmount: initialSupply * 11 / 100,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod*2,
			newChainTime:       mintSince + mintingPeriod*3,
			expectedMintAmount: initialSupply * 121 / 1000,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod*3,
			newChainTime:       mintSince + mintingPeriod*4,
			expectedMintAmount: initialSupply * 1331 / 10000,
		},
		{
			lastSyncTime:       mintSince + mintingPeriod*99,
			newChainTime:       mintSince + mintingPeriod*100,
			expectedMintAmount: 1252782939,
		},
		{
			lastSyncTime:       0,
			newChainTime:       mintSince + mintingPeriod*2,
			expectedMintAmount: initialSupply * 21 / 100,
		},
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod*2,
			expectedMintAmount: initialSupply * 21 / 100,
		},
		{
			lastSyncTime:       0,
			newChainTime:       mintSince + mintingPeriod*3,
			expectedMintAmount: initialSupply * 331 / 1000,
		},
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod*3,
			expectedMintAmount: initialSupply * 331 / 1000,
		},
		{
			lastSyncTime:       0,
			newChainTime:       mintSince + mintingPeriod*4,
			expectedMintAmount: initialSupply * 4641 / 10000,
		},
		{
			lastSyncTime:       mintSince,
			newChainTime:       mintSince + mintingPeriod*4,
			expectedMintAmount: initialSupply * 4641 / 10000,
		},
	}

	mintConfig := MintConfig{
		MintSince:     int64(mintSince),
		MintingPeriod: time.Duration(mintingPeriod * int64(time.Second)),
		MintRate:      mintRate,
		MaxMintAmount: maxMintAmount,
	}
	c := NewMintCalculator(mintConfig, initialSupply)
	for totalWeight := uint64(1); totalWeight < 10; totalWeight++ {
		for weight := uint64(0); weight <= totalWeight; weight++ {
			for _, test := range tests {
				expectedReward := test.expectedMintAmount * weight / totalWeight
				name := fmt.Sprintf("mint(%d,%d,%d,%d)==%d",
					weight,
					totalWeight,
					test.lastSyncTime,
					test.newChainTime,
					expectedReward,
				)
				t.Run(name, func(t *testing.T) {
					mintRate := c.CalculateMintRate(
						totalWeight,
						time.Unix(int64(test.lastSyncTime), 0),
						time.Unix(int64(test.newChainTime), 0),
					)
					reward := CalculateMintReward(weight, new(big.Int), mintRate)

					// might happen roundoff error
					require.True(t, expectedReward-reward <= 1, "%d != %d", expectedReward, reward)
				})
			}
		}
	}
}
