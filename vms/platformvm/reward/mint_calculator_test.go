package reward

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var defaultMintConfig = MintConfig{
	MintSince:  100,
	MintUntil:  200,
	MintAmount: 1_000_000,
}

func TestMint(t *testing.T) {
	c, _ := NewMintCalculator(defaultMintConfig)
	tests := []struct {
		lastSyncTime       int64
		newChainTime       int64
		stakerWeight       uint64
		totalWeight        uint64
		expectedMintAmount uint64
	}{
		{
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintUntil * 2,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintSince,
			expectedMintAmount: 0,
		},
		{
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       (defaultMintConfig.MintSince + defaultMintConfig.MintUntil) / 2,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			lastSyncTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)/4,
			newChainTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)*3/4,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
	}

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
					require.True(t, expectedReward-reward <= 1)
				})
			}
		}
	}
}
