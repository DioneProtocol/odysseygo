// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/proto/pb/p2p"
)

func TestEngineManager_Get(t *testing.T) {
	type args struct {
		engineType p2p.EngineType
	}

	odyssey := &Engine{}
	snowman := &Engine{}

	type expected struct {
		engine *Engine
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "request unspecified engine",
			args: args{
				engineType: p2p.EngineType_ENGINE_TYPE_UNSPECIFIED,
			},
			expected: expected{
				engine: nil,
			},
		},
		{
			name: "request odyssey engine",
			args: args{
				engineType: p2p.EngineType_ENGINE_TYPE_ODYSSEY,
			},
			expected: expected{
				engine: odyssey,
			},
		},
		{
			name: "request snowman engine",
			args: args{
				engineType: p2p.EngineType_ENGINE_TYPE_SNOWMAN,
			},
			expected: expected{
				engine: snowman,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			e := EngineManager{
				Odyssey: odyssey,
				Snowman:   snowman,
			}

			r.Equal(test.expected.engine, e.Get(test.args.engineType))
		})
	}
}
