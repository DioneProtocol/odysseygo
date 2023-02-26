// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/vms"
)

var _ vms.Factory = (*Factory)(nil)

type Factory struct {
	TxFee            uint64
	CreateAssetTxFee uint64
}

func (f *Factory) New(*snow.Context) (interface{}, error) {
	return &VM{Factory: *f}, nil
}
