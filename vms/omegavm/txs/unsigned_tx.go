// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

// UnsignedTx is an unsigned transaction
type UnsignedTx interface {
	// TODO: Remove this initialization pattern from both the omegavm and the
	// avm.
	snow.ContextInitializable
	secp256k1fx.UnsignedTx
	SetBytes(unsignedBytes []byte)

	// InputIDs returns the set of inputs this transaction consumes
	InputIDs() set.Set[ids.ID]

	Outputs() []*dione.TransferableOutput

	// Attempts to verify this transaction without any provided state.
	SyntacticVerify(ctx *snow.Context) error

	// Visit calls [visitor] with this transaction's concrete type
	Visit(visitor Visitor) error
}
