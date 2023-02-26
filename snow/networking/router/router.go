// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package router

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dioneprotocol/dionego/api/health"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/message"
	"github.com/dioneprotocol/dionego/snow/networking/benchlist"
	"github.com/dioneprotocol/dionego/snow/networking/handler"
	"github.com/dioneprotocol/dionego/snow/networking/timeout"
	"github.com/dioneprotocol/dionego/utils/logging"
	"github.com/dioneprotocol/dionego/utils/set"
)

// Router routes consensus messages to the Handler of the consensus
// engine that the messages are intended for
type Router interface {
	ExternalHandler
	InternalHandler

	Initialize(
		nodeID ids.NodeID,
		log logging.Logger,
		timeouts timeout.Manager,
		shutdownTimeout time.Duration,
		criticalChains set.Set[ids.ID],
		stakingEnabled bool,
		trackedSubnets set.Set[ids.ID],
		onFatal func(exitCode int),
		healthConfig HealthConfig,
		metricsNamespace string,
		metricsRegisterer prometheus.Registerer,
	) error
	Shutdown(context.Context)
	AddChain(ctx context.Context, chain handler.Handler)
	health.Checker
}

// InternalHandler deals with messages internal to this node
type InternalHandler interface {
	benchlist.Benchable

	RegisterRequest(
		ctx context.Context,
		nodeID ids.NodeID,
		sourceChainID ids.ID,
		destinationChainID ids.ID,
		requestID uint32,
		op message.Op,
		failedMsg message.InboundMessage,
	)
}
