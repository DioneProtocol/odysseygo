// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"context"

	"github.com/DioneProtocol/odysseygo/utils/rpc"
)

var _ StaticClient = (*staticClient)(nil)

// StaticClient for interacting with the omegavm static api
type StaticClient interface {
	BuildGenesis(
		ctx context.Context,
		args *BuildGenesisArgs,
		options ...rpc.Option,
	) (*BuildGenesisReply, error)
}

// staticClient is an implementation of a omegavm client for interacting with
// the omegavm static api
type staticClient struct {
	requester rpc.EndpointRequester
}

// NewClient returns a omegavm client for interacting with the omegavm static api
func NewStaticClient(uri string) StaticClient {
	return &staticClient{requester: rpc.NewEndpointRequester(
		uri + "/ext/vm/omega",
	)}
}

func (c *staticClient) BuildGenesis(
	ctx context.Context,
	args *BuildGenesisArgs,
	options ...rpc.Option,
) (resp *BuildGenesisReply, err error) {
	resp = &BuildGenesisReply{}
	err = c.requester.SendRequest(ctx, "omega.buildGenesis", args, resp, options...)
	return resp, err
}
