// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/consensus/odyssey"
)

var (
	errBuild = errors.New("unexpectedly called Build")

	_ Builder = (*TestBuilder)(nil)
)

type TestBuilder struct {
	T             *testing.T
	CantBuildVtx  bool
	BuildStopVtxF func(ctx context.Context, parentIDs []ids.ID) (odyssey.Vertex, error)
}

func (b *TestBuilder) Default(cant bool) {
	b.CantBuildVtx = cant
}

func (b *TestBuilder) BuildStopVtx(ctx context.Context, parentIDs []ids.ID) (odyssey.Vertex, error) {
	if b.BuildStopVtxF != nil {
		return b.BuildStopVtxF(ctx, parentIDs)
	}
	if b.CantBuildVtx && b.T != nil {
		require.FailNow(b.T, errBuild.Error())
	}
	return nil, errBuild
}