// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"context"
	"crypto"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/message"
	"github.com/DioneProtocol/odysseygo/network/throttling"
	"github.com/DioneProtocol/odysseygo/snow/networking/router"
	"github.com/DioneProtocol/odysseygo/snow/networking/tracker"
	"github.com/DioneProtocol/odysseygo/snow/uptime"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/staking"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/ips"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/utils/math/meter"
	"github.com/DioneProtocol/odysseygo/utils/resource"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/version"
)

const maxMessageToSend = 1024

// StartTestPeer provides a simple interface to create a peer that has finished
// the p2p handshake.
//
// This function will generate a new TLS key to use when connecting to the peer.
//
// The returned peer will not throttle inbound or outbound messages.
//
//   - [ctx] provides a way of canceling the connection request.
//   - [ip] is the remote that will be dialed to create the connection.
//   - [networkID] will be sent to the peer during the handshake. If the peer is
//     expecting a different [networkID], the handshake will fail and an error
//     will be returned.
//   - [router] will be called with all non-handshake messages received by the
//     peer.
func StartTestPeer(
	ctx context.Context,
	ip ips.IPPort,
	networkID uint32,
	router router.InboundHandler,
) (Peer, error) {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, constants.NetworkType, ip.String())
	if err != nil {
		return nil, err
	}

	tlsCert, err := staking.NewTLSCert()
	if err != nil {
		return nil, err
	}

	tlsConfg := TLSConfig(*tlsCert, nil)
	clientUpgrader := NewTLSClientUpgrader(
		tlsConfg,
		prometheus.NewCounter(prometheus.CounterOpts{}),
	)

	peerID, conn, cert, err := clientUpgrader.Upgrade(conn)
	if err != nil {
		return nil, err
	}

	mc, err := message.NewCreator(
		logging.NoLog{},
		prometheus.NewRegistry(),
		"",
		constants.DefaultNetworkCompressionType,
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	metrics, err := NewMetrics(
		logging.NoLog{},
		"",
		prometheus.NewRegistry(),
	)
	if err != nil {
		return nil, err
	}

	resourceTracker, err := tracker.NewResourceTracker(
		prometheus.NewRegistry(),
		resource.NoUsage,
		meter.ContinuousFactory{},
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	signerIP := ips.NewDynamicIPPort(net.IPv6zero, 0)
	tls := tlsCert.PrivateKey.(crypto.Signer)

	peer := Start(
		&Config{
			Metrics:              metrics,
			MessageCreator:       mc,
			Log:                  logging.NoLog{},
			InboundMsgThrottler:  throttling.NewNoInboundThrottler(),
			Network:              TestNetwork,
			Router:               router,
			VersionCompatibility: version.GetCompatibility(networkID),
			MySubnets:            set.Set[ids.ID]{},
			Beacons:              validators.NewSet(),
			NetworkID:            networkID,
			PingFrequency:        constants.DefaultPingFrequency,
			PongTimeout:          constants.DefaultPingPongTimeout,
			MaxClockDifference:   time.Minute,
			ResourceTracker:      resourceTracker,
			UptimeCalculator:     uptime.NoOpCalculator,
			IPSigner:             NewIPSigner(signerIP, tls),
		},
		conn,
		cert,
		peerID,
		NewBlockingMessageQueue(
			metrics,
			logging.NoLog{},
			maxMessageToSend,
		),
	)
	return peer, peer.AwaitReady(ctx)
}
