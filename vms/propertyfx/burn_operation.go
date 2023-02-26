// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/vms/components/verify"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

type BurnOperation struct {
	secp256k1fx.Input `serialize:"true"`
}

func (*BurnOperation) InitCtx(*snow.Context) {}

func (*BurnOperation) Outs() []verify.State {
	return nil
}
