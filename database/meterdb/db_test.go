// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package meterdb

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/memdb"
)

func TestInterface(t *testing.T) {
	for _, test := range database.Tests {
		baseDB := memdb.New()
		db, err := New("", prometheus.NewRegistry(), baseDB)
		require.NoError(t, err)

		test(t, db)
	}
}

func FuzzKeyValue(f *testing.F) {
	baseDB := memdb.New()
	db, err := New("", prometheus.NewRegistry(), baseDB)
	require.NoError(f, err)
	database.FuzzKeyValue(f, db)
}

func FuzzNewIteratorWithPrefix(f *testing.F) {
	baseDB := memdb.New()
	db, err := New("", prometheus.NewRegistry(), baseDB)
	require.NoError(f, err)
	database.FuzzNewIteratorWithPrefix(f, db)
}

func BenchmarkInterface(b *testing.B) {
	for _, size := range database.BenchmarkSizes {
		keys, values := database.SetupBenchmark(b, size[0], size[1], size[2])
		for _, bench := range database.Benchmarks {
			baseDB := memdb.New()
			db, err := New("", prometheus.NewRegistry(), baseDB)
			require.NoError(b, err)
			bench(b, db, "meterdb", keys, values)
		}
	}
}
