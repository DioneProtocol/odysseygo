// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowball

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/utils/sampler"
)

func TestSnowballOptimized(t *testing.T) {
	require := require.New(t)

	numColors := 10
	numNodes := 100
	params := Parameters{
		K: 20, Alpha: 15, BetaVirtuous: 20, BetaRogue: 30,
	}
	seed := int64(0)

	nBitwise := Network{}
	nBitwise.Initialize(params, numColors)

	nNaive := nBitwise

	sampler.Seed(seed)
	for i := 0; i < numNodes; i++ {
		nBitwise.AddNode(&Tree{})
	}

	sampler.Seed(seed)
	for i := 0; i < numNodes; i++ {
		nNaive.AddNode(&Flat{})
	}

	numRounds := 0
	for !nBitwise.Finalized() && !nBitwise.Disagreement() && !nNaive.Finalized() && !nNaive.Disagreement() {
		sampler.Seed(int64(numRounds) + seed)
		nBitwise.Round()

		sampler.Seed(int64(numRounds) + seed)
		nNaive.Round()
		numRounds++
	}

	require.False(nBitwise.Disagreement())
	require.False(nNaive.Disagreement())

	// Although this can theoretically fail with a correct implementation, it
	// shouldn't in practice
	require.True(nBitwise.Finalized())
	require.True(nBitwise.Agreement())
}
