// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow/networking/router"
	"github.com/dioneprotocol/dionego/version"
)

var _ router.ExternalHandler = (*testHandler)(nil)

type testHandler struct {
	router.InboundHandler
	ConnectedF    func(nodeID ids.NodeID, nodeVersion *version.Application, subnetID ids.ID)
	DisconnectedF func(nodeID ids.NodeID)
}

func (h *testHandler) Connected(id ids.NodeID, nodeVersion *version.Application, subnetID ids.ID) {
	if h.ConnectedF != nil {
		h.ConnectedF(id, nodeVersion, subnetID)
	}
}

func (h *testHandler) Disconnected(id ids.NodeID) {
	if h.DisconnectedF != nil {
		h.DisconnectedF(id)
	}
}
