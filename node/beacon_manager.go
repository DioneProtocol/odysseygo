// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"sync/atomic"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow/networking/router"
	"github.com/dioneprotocol/dionego/snow/validators"
	"github.com/dioneprotocol/dionego/utils/constants"
	"github.com/dioneprotocol/dionego/utils/timer"
	"github.com/dioneprotocol/dionego/version"
)

var _ router.Router = (*beaconManager)(nil)

type beaconManager struct {
	router.Router
	timer         *timer.Timer
	beacons       validators.Set
	requiredConns int64
	numConns      int64
}

func (b *beaconManager) Connected(nodeID ids.NodeID, nodeVersion *version.Application, subnetID ids.ID) {
	if constants.PrimaryNetworkID == subnetID &&
		b.beacons.Contains(nodeID) &&
		atomic.AddInt64(&b.numConns, 1) >= b.requiredConns {
		b.timer.Cancel()
	}
	b.Router.Connected(nodeID, nodeVersion, subnetID)
}

func (b *beaconManager) Disconnected(nodeID ids.NodeID) {
	if b.beacons.Contains(nodeID) {
		atomic.AddInt64(&b.numConns, -1)
	}
	b.Router.Disconnected(nodeID)
}
