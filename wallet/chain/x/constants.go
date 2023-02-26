// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"github.com/dioneprotocol/dionego/vms/avm/fxs"
	"github.com/dioneprotocol/dionego/vms/avm/txs"
	"github.com/dioneprotocol/dionego/vms/nftfx"
	"github.com/dioneprotocol/dionego/vms/propertyfx"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

const (
	SECP256K1FxIndex = 0
	NFTFxIndex       = 1
	PropertyFxIndex  = 2
)

// Parser to support serialization and deserialization
var Parser txs.Parser

func init() {
	var err error
	Parser, err = txs.NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
		&nftfx.Fx{},
		&propertyfx.Fx{},
	})
	if err != nil {
		panic(err)
	}
}
