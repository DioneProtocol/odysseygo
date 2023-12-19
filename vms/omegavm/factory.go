// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package omegavm

import (
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/vms"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/config"
)

var _ vms.Factory = (*Factory)(nil)

// Factory can create new instances of the Omega Chain
type Factory struct {
	config.Config
}

// New returns a new instance of the Omega Chain
func (f *Factory) New(logging.Logger) (interface{}, error) {
	return &VM{Config: f.Config}, nil
}
