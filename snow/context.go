// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snow

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/DioneProtocol/odysseygo/api/keystore"
	"github.com/DioneProtocol/odysseygo/api/metrics"
	"github.com/DioneProtocol/odysseygo/chains/atomic"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/warp"
)

// ContextInitializable represents an object that can be initialized
// given a *Context object
type ContextInitializable interface {
	// InitCtx initializes an object provided a *Context object
	InitCtx(ctx *Context)
}

// Context is information about the current execution.
// [NetworkID] is the ID of the network this context exists within.
// [ChainID] is the ID of the chain this context exists within.
// [NodeID] is the ID of this node
type Context struct {
	NetworkID uint32
	SubnetID  ids.ID
	ChainID   ids.ID
	NodeID    ids.NodeID
	PublicKey *bls.PublicKey

	XChainID     ids.ID
	DChainID     ids.ID
	DIONEAssetID ids.ID

	Log          logging.Logger
	Lock         sync.RWMutex
	Keystore     keystore.BlockchainKeystore
	SharedMemory atomic.SharedMemory
	BCLookup     ids.AliaserReader
	Metrics      metrics.OptionalGatherer

	WarpSigner warp.Signer

	// snowman++ attributes
	ValidatorState validators.State // interface for O-Chain validators
	// Chain-specific directory where arbitrary data can be written
	ChainDataDir string
}

// Expose gatherer interface for unit testing.
type Registerer interface {
	prometheus.Registerer
	prometheus.Gatherer
}

type ConsensusContext struct {
	*Context

	// Registers all common and snowman consensus metrics. Unlike the odyssey
	// consensus engine metrics, we do not prefix the name with the engine name,
	// as snowman is used for all chains by default.
	Registerer Registerer
	// Only used to register Odyssey consensus metrics. Previously, all
	// metrics were prefixed with "odyssey_{chainID}_". Now we add odyssey
	// to the prefix, "odyssey_{chainID}_odyssey_", to differentiate
	// consensus operations after the DAG linearization.
	OdysseyRegisterer Registerer

	// BlockAcceptor is the callback that will be fired whenever a VM is
	// notified that their block was accepted.
	BlockAcceptor Acceptor

	// TxAcceptor is the callback that will be fired whenever a VM is notified
	// that their transaction was accepted.
	TxAcceptor Acceptor

	// VertexAcceptor is the callback that will be fired whenever a vertex was
	// accepted.
	VertexAcceptor Acceptor

	// State indicates the current state of this consensus instance.
	State utils.Atomic[EngineState]

	// True iff this chain is executing transactions as part of bootstrapping.
	Executing utils.Atomic[bool]

	// True iff this chain is currently state-syncing
	StateSyncing utils.Atomic[bool]
}

func DefaultContextTest() *Context {
	sk, err := bls.NewSecretKey()
	if err != nil {
		panic(err)
	}
	pk := bls.PublicFromSecretKey(sk)
	return &Context{
		NetworkID:    0,
		SubnetID:     ids.Empty,
		ChainID:      ids.Empty,
		NodeID:       ids.EmptyNodeID,
		PublicKey:    pk,
		Log:          logging.NoLog{},
		BCLookup:     ids.NewAliaser(),
		Metrics:      metrics.NewOptionalGatherer(),
		ChainDataDir: "",
	}
}

func DefaultConsensusContextTest() *ConsensusContext {
	return &ConsensusContext{
		Context:           DefaultContextTest(),
		Registerer:        prometheus.NewRegistry(),
		OdysseyRegisterer: prometheus.NewRegistry(),
		BlockAcceptor:     noOpAcceptor{},
		TxAcceptor:        noOpAcceptor{},
		VertexAcceptor:    noOpAcceptor{},
	}
}
