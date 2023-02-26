// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"context"
	"errors"
	"testing"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow/consensus/dione"
	"github.com/dioneprotocol/dionego/snow/consensus/snowstorm"
)

var (
	errBuild = errors.New("unexpectedly called Build")

	_ Builder = (*TestBuilder)(nil)
)

type TestBuilder struct {
	T             *testing.T
	CantBuildVtx  bool
	BuildVtxF     func(ctx context.Context, parentIDs []ids.ID, txs []snowstorm.Tx) (dione.Vertex, error)
	BuildStopVtxF func(ctx context.Context, parentIDs []ids.ID) (dione.Vertex, error)
}

func (b *TestBuilder) Default(cant bool) {
	b.CantBuildVtx = cant
}

func (b *TestBuilder) BuildVtx(ctx context.Context, parentIDs []ids.ID, txs []snowstorm.Tx) (dione.Vertex, error) {
	if b.BuildVtxF != nil {
		return b.BuildVtxF(ctx, parentIDs, txs)
	}
	if b.CantBuildVtx && b.T != nil {
		b.T.Fatal(errBuild)
	}
	return nil, errBuild
}

func (b *TestBuilder) BuildStopVtx(ctx context.Context, parentIDs []ids.ID) (dione.Vertex, error) {
	if b.BuildStopVtxF != nil {
		return b.BuildStopVtxF(ctx, parentIDs)
	}
	if b.CantBuildVtx && b.T != nil {
		b.T.Fatal(errBuild)
	}
	return nil, errBuild
}
