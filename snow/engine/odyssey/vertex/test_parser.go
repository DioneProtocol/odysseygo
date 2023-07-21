// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"context"
	"errors"
	"testing"

	"github.com/DioneProtocol/odysseygo/snow/consensus/odyssey"
)

var (
	errParse = errors.New("unexpectedly called Parse")

	_ Parser = (*TestParser)(nil)
)

type TestParser struct {
	T            *testing.T
	CantParseVtx bool
	ParseVtxF    func(context.Context, []byte) (odyssey.Vertex, error)
}

func (p *TestParser) Default(cant bool) {
	p.CantParseVtx = cant
}

func (p *TestParser) ParseVtx(ctx context.Context, b []byte) (odyssey.Vertex, error) {
	if p.ParseVtxF != nil {
		return p.ParseVtxF(ctx, b)
	}
	if p.CantParseVtx && p.T != nil {
		p.T.Fatal(errParse)
	}
	return nil, errParse
}
