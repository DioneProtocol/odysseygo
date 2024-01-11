// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package secp256k1fx

import (
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/vms"
)

var (
	_ vms.Factory = (*Factory)(nil)

	// ID that this Fx uses when labeled
	ID = ids.ID{'s', 'e', 'c', 'p', '2', '5', '6', 'k', '1', 'f', 'x'}
)

type Factory struct{}

func (*Factory) New(logging.Logger) (interface{}, error) {
	return &Fx{}, nil
}
