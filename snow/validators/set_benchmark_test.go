// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"fmt"
	"testing"
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
)

// setupValidators creates a validator set with the specified number of validators
func setupValidators(b *testing.B, totalValidators int, orionCount int) (*vdrSet, []ids.NodeID) {
	b.Helper()

	s := NewSet().(*vdrSet)
	orionNodes := make([]ids.NodeID, 0, orionCount)

	// Create validators
	for i := 0; i < totalValidators; i++ {
		nodeID := ids.GenerateTestNodeID()
		weight := uint64(1000 + i) // Varying weights

		if err := s.Add(nodeID, nil, ids.Empty, weight); err != nil {
			b.Fatalf("Failed to add validator: %v", err)
		}

		// Mark first orionCount validators as orions
		if i < orionCount {
			orionNodes = append(orionNodes, nodeID)
		}
	}

	return s, orionNodes
}

// BenchmarkSampleApricotPhase7_SmallSet benchmarks sampling with a small validator set
func BenchmarkSampleApricotPhase7_SmallSet(b *testing.B) {
	// Activate Apricot Phase 7
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	s, orionNodes := setupValidators(b, 10, 2)
	SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
	defer SetOrionChecker(nil)

	sampleSize := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Sample(sampleSize)
		if err != nil {
			b.Fatalf("Sampling failed: %v", err)
		}
	}
}

// BenchmarkSampleApricotPhase7_MediumSet benchmarks sampling with a medium validator set
func BenchmarkSampleApricotPhase7_MediumSet(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	s, orionNodes := setupValidators(b, 100, 20)
	SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
	defer SetOrionChecker(nil)

	sampleSize := 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Sample(sampleSize)
		if err != nil {
			b.Fatalf("Sampling failed: %v", err)
		}
	}
}

// BenchmarkSampleApricotPhase7_LargeSet benchmarks sampling with a large validator set
func BenchmarkSampleApricotPhase7_LargeSet(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	s, orionNodes := setupValidators(b, 1000, 200)
	SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
	defer SetOrionChecker(nil)

	sampleSize := 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Sample(sampleSize)
		if err != nil {
			b.Fatalf("Sampling failed: %v", err)
		}
	}
}

// BenchmarkSampleApricotPhase7_AllOrions benchmarks when all validators are orions
func BenchmarkSampleApricotPhase7_AllOrions(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	totalValidators := 100
	s, orionNodes := setupValidators(b, totalValidators, totalValidators)
	SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
	defer SetOrionChecker(nil)

	sampleSize := 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Sample(sampleSize)
		if err != nil {
			b.Fatalf("Sampling failed: %v", err)
		}
	}
}

// BenchmarkSampleApricotPhase7_NoOrions benchmarks when no validators are orions
func BenchmarkSampleApricotPhase7_NoOrions(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	s, _ := setupValidators(b, 100, 0)
	SetOrionChecker(&mockOrionChecker{nodes: []ids.NodeID{}})
	defer SetOrionChecker(nil)

	sampleSize := 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Sample(sampleSize)
		if err != nil {
			b.Fatalf("Sampling failed: %v", err)
		}
	}
}

// BenchmarkSampleApricotPhase7_VaryingSampleSizes benchmarks different sample sizes
func BenchmarkSampleApricotPhase7_VaryingSampleSizes(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	sampleSizes := []int{1, 5, 10, 20, 50}

	for _, sampleSize := range sampleSizes {
		b.Run(fmt.Sprintf("sampleSize_%d", sampleSize), func(b *testing.B) {
			s, orionNodes := setupValidators(b, 100, 20)
			SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
			defer SetOrionChecker(nil)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := s.Sample(sampleSize)
				if err != nil {
					b.Fatalf("Sampling failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkSampleApricotPhase7_VaryingOrionRatios benchmarks different orion ratios
func BenchmarkSampleApricotPhase7_VaryingOrionRatios(b *testing.B) {
	SetApricotPhase7ActivationTime(time.Now().Add(time.Hour))
	defer SetApricotPhase7ActivationTime(time.Time{})

	totalValidators := 100
	orionRatios := []struct {
		name       string
		orionCount int
	}{
		{"0%", 0},
		{"10%", 10},
		{"20%", 20},
		{"50%", 50},
		{"100%", 100},
	}

	for _, ratio := range orionRatios {
		b.Run(ratio.name, func(b *testing.B) {
			s, orionNodes := setupValidators(b, totalValidators, ratio.orionCount)
			SetOrionChecker(&mockOrionChecker{nodes: orionNodes})
			defer SetOrionChecker(nil)

			sampleSize := 20

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := s.Sample(sampleSize)
				if err != nil {
					b.Fatalf("Sampling failed: %v", err)
				}
			}
		})
	}
}
