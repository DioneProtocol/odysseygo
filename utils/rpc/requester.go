// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"context"
	"net/url"
)

var _ EndpointRequester = (*odysseyEndpointRequester)(nil)

type EndpointRequester interface {
	SendRequest(ctx context.Context, method string, params interface{}, reply interface{}, options ...Option) error
}

type odysseyEndpointRequester struct {
	uri string
}

func NewEndpointRequester(uri string) EndpointRequester {
	return &odysseyEndpointRequester{
		uri: uri,
	}
}

func (e *odysseyEndpointRequester) SendRequest(
	ctx context.Context,
	method string,
	params interface{},
	reply interface{},
	options ...Option,
) error {
	uri, err := url.Parse(e.uri)
	if err != nil {
		return err
	}

	return SendJSONRequest(
		ctx,
		uri,
		method,
		params,
		reply,
		options...,
	)
}
