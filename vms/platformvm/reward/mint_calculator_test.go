package reward

import (
	"fmt"
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
	c := NewMintCalculator(defaultMintConfig)
	tests := []struct {
		validatorsAmount   uint64
		lastSyncTime       uint64
		newChainTime       uint64
		expectedMintAmount uint64
	}{
		{
			validatorsAmount:   1,
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintUntil * 2,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount,
		},
		{
			validatorsAmount:   2,
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			validatorsAmount:   2,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       0,
			newChainTime:       defaultMintConfig.MintSince,
			expectedMintAmount: 0,
		},
		{
			validatorsAmount:   0,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: 0,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       (defaultMintConfig.MintSince + defaultMintConfig.MintUntil) / 2,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       (defaultMintConfig.MintSince + defaultMintConfig.MintUntil) / 2,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)/4,
			newChainTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)*3/4,
			expectedMintAmount: defaultMintConfig.MintAmount / 2,
		},
		{
			validatorsAmount:   2,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       (defaultMintConfig.MintSince + defaultMintConfig.MintUntil) / 2,
			expectedMintAmount: defaultMintConfig.MintAmount / 4,
		},
		{
			validatorsAmount:   2,
			lastSyncTime:       (defaultMintConfig.MintSince + defaultMintConfig.MintUntil) / 2,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount / 4,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       defaultMintConfig.MintSince,
			newChainTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)/4,
			expectedMintAmount: defaultMintConfig.MintAmount / 4,
		},
		{
			validatorsAmount:   1,
			lastSyncTime:       defaultMintConfig.MintSince + (defaultMintConfig.MintUntil-defaultMintConfig.MintSince)*3/4,
			newChainTime:       defaultMintConfig.MintUntil,
			expectedMintAmount: defaultMintConfig.MintAmount / 4,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("mint(%d,%d,%d)==%d",
			test.validatorsAmount,
			test.lastSyncTime,
			test.newChainTime,
			test.expectedMintAmount,
		)
		t.Run(name, func(t *testing.T) {
			reward := c.CalculateMintRate(
				test.validatorsAmount,
				time.Unix(int64(test.lastSyncTime), 0),
				time.Unix(int64(test.newChainTime), 0),
			)
			require.Equal(t, test.expectedMintAmount, reward)
		})
	}
}
