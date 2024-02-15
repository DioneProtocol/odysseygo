// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import (
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/gorilla/rpc/v2"

	"go.uber.org/zap"

	"github.com/DioneProtocol/odysseygo/api/server"
	"github.com/DioneProtocol/odysseygo/chains"
	"github.com/DioneProtocol/odysseygo/codec"
	"github.com/DioneProtocol/odysseygo/codec/linearcodec"
	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/engine/avalanche/vertex"
	"github.com/DioneProtocol/odysseygo/snow/engine/common"
	"github.com/DioneProtocol/odysseygo/snow/engine/snowman/block"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/json"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/utils/timer/mockable"
	"github.com/DioneProtocol/odysseygo/utils/wrappers"
)

const (
	indexNamePrefix = "index-"
	codecVersion    = uint16(0)
	// Max size, in bytes, of something serialized by this indexer
	// Assumes no containers are larger than math.MaxUint32
	// wrappers.IntLen accounts for the size of the container bytes
	// wrappers.LongLen accounts for the timestamp of the container
	// ids.IDLen accounts for the container ID
	// wrappers.ShortLen accounts for the codec version
	codecMaxSize = int(constants.DefaultMaxMessageSize) + wrappers.IntLen + wrappers.LongLen + ids.IDLen + wrappers.ShortLen
)

var (
	txPrefix                = byte(0x01)
	vtxPrefix               = byte(0x02)
	blockPrefix             = byte(0x03)
	isIncompletePrefix      = byte(0x04)
	previouslyIndexedPrefix = byte(0x05)
	hasRunKey               = []byte{0x07}

	_ Indexer = (*indexer)(nil)
)

// Config for an indexer
type Config struct {
	DB                   database.Database
	Log                  logging.Logger
	IndexingEnabled      bool
	AllowIncompleteIndex bool
	BlockAcceptorGroup   snow.AcceptorGroup
	TxAcceptorGroup      snow.AcceptorGroup
	VertexAcceptorGroup  snow.AcceptorGroup
	APIServer            server.PathAdder
	ShutdownF            func()
}

// Indexer causes accepted containers for a given chain
// to be indexed by their ID and by the order in which
// they were accepted by this node.
// Indexer is threadsafe.
type Indexer interface {
	chains.Registrant
	// Close will do nothing and return nil after the first call
	io.Closer
}

// NewIndexer returns a new Indexer and registers a new endpoint on the given API server.
func NewIndexer(config Config) (Indexer, error) {
	indexer := &indexer{
		codec:                codec.NewManager(codecMaxSize),
		log:                  config.Log,
		db:                   config.DB,
		allowIncompleteIndex: config.AllowIncompleteIndex,
		indexingEnabled:      config.IndexingEnabled,
		blockAcceptorGroup:   config.BlockAcceptorGroup,
		txAcceptorGroup:      config.TxAcceptorGroup,
		vertexAcceptorGroup:  config.VertexAcceptorGroup,
		txIndices:            map[ids.ID]Index{},
		vtxIndices:           map[ids.ID]Index{},
		blockIndices:         map[ids.ID]Index{},
		pathAdder:            config.APIServer,
		shutdownF:            config.ShutdownF,
	}

	if err := indexer.codec.RegisterCodec(
		codecVersion,
		linearcodec.NewCustomMaxLength(math.MaxUint32),
	); err != nil {
		return nil, fmt.Errorf("couldn't register codec: %w", err)
	}
	hasRun, err := indexer.hasRun()
	if err != nil {
		return nil, err
	}
	indexer.hasRunBefore = hasRun
	return indexer, indexer.markHasRun()
}

type indexer struct {
	codec  codec.Manager
	clock  mockable.Clock
	lock   sync.RWMutex
	log    logging.Logger
	db     database.Database
	closed bool

	// Called in a goroutine on shutdown
	shutdownF func()

	// true if this is not the first run using this database
	hasRunBefore bool

	// Used to add API endpoint for new indices
	pathAdder server.PathAdder

	// If true, allow running in such a way that could allow the creation
	// of an index which could be missing accepted containers.
	allowIncompleteIndex bool

	// If false, don't create index for a chain when RegisterChain is called
	indexingEnabled bool

	// Chain ID --> index of blocks of that chain (if applicable)
	blockIndices map[ids.ID]Index
	// Chain ID --> index of vertices of that chain (if applicable)
	vtxIndices map[ids.ID]Index
	// Chain ID --> index of txs of that chain (if applicable)
	txIndices map[ids.ID]Index

	// Notifies of newly accepted blocks
	blockAcceptorGroup snow.AcceptorGroup
	// Notifies of newly accepted transactions
	txAcceptorGroup snow.AcceptorGroup
	// Notifies of newly accepted vertices
	vertexAcceptorGroup snow.AcceptorGroup
}

// Assumes [ctx.Lock] is not held
func (i *indexer) RegisterChain(chainName string, ctx *snow.ConsensusContext, vm common.VM) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.closed {
		i.log.Debug("not registering chain to indexer",
			zap.String("reason", "indexer is closed"),
			zap.String("chainName", chainName),
		)
		return
	} else if ctx.SubnetID != constants.PrimaryNetworkID {
		i.log.Debug("not registering chain to indexer",
			zap.String("reason", "not in the primary network"),
			zap.String("chainName", chainName),
		)
		return
	}

	chainID := ctx.ChainID
	if i.blockIndices[chainID] != nil || i.txIndices[chainID] != nil || i.vtxIndices[chainID] != nil {
		i.log.Warn("chain is already being indexed",
			zap.Stringer("chainID", chainID),
		)
		return
	}

	// If the index is incomplete, make sure that's OK. Otherwise, cause node to die.
	isIncomplete, err := i.isIncomplete(chainID)
	if err != nil {
		i.log.Error("couldn't get whether chain is incomplete",
			zap.String("chainName", chainName),
			zap.Error(err),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}

	// See if this chain was indexed in a previous run
	previouslyIndexed, err := i.previouslyIndexed(chainID)
	if err != nil {
		i.log.Error("couldn't get whether chain was previously indexed",
			zap.String("chainName", chainName),
			zap.Error(err),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}

	if !i.indexingEnabled { // Indexing is disabled
		if previouslyIndexed && !i.allowIncompleteIndex {
			// We indexed this chain in a previous run but not in this run.
			// This would create an incomplete index, which is not allowed, so exit.
			i.log.Fatal("running would cause index to become incomplete but incomplete indices are disabled",
				zap.String("chainName", chainName),
			)
			if err := i.close(); err != nil {
				i.log.Error("failed to close indexer",
					zap.Error(err),
				)
			}
			return
		}

		// Creating an incomplete index is allowed. Mark index as incomplete.
		err := i.markIncomplete(chainID)
		if err == nil {
			return
		}
		i.log.Fatal("couldn't mark chain as incomplete",
			zap.String("chainName", chainName),
			zap.Error(err),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}

	if !i.allowIncompleteIndex && isIncomplete && (previouslyIndexed || i.hasRunBefore) {
		i.log.Fatal("index is incomplete but incomplete indices are disabled. Shutting down",
			zap.String("chainName", chainName),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}

	// Mark that in this run, this chain was indexed
	if err := i.markPreviouslyIndexed(chainID); err != nil {
		i.log.Error("couldn't mark chain as indexed",
			zap.String("chainName", chainName),
			zap.Error(err),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}

	index, err := i.registerChainHelper(chainID, blockPrefix, chainName, "block", i.blockAcceptorGroup)
	if err != nil {
		i.log.Fatal("failed to create index",
			zap.String("chainName", chainName),
			zap.String("endpoint", "block"),
			zap.Error(err),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
		return
	}
	i.blockIndices[chainID] = index

	switch vm.(type) {
	case vertex.DAGVM:
		vtxIndex, err := i.registerChainHelper(chainID, vtxPrefix, chainName, "vtx", i.vertexAcceptorGroup)
		if err != nil {
			i.log.Fatal("couldn't create index",
				zap.String("chainName", chainName),
				zap.String("endpoint", "vtx"),
				zap.Error(err),
			)
			if err := i.close(); err != nil {
				i.log.Error("failed to close indexer",
					zap.Error(err),
				)
			}
			return
		}
		i.vtxIndices[chainID] = vtxIndex

		txIndex, err := i.registerChainHelper(chainID, txPrefix, chainName, "tx", i.txAcceptorGroup)
		if err != nil {
			i.log.Fatal("couldn't create index",
				zap.String("chainName", chainName),
				zap.String("endpoint", "tx"),
				zap.Error(err),
			)
			if err := i.close(); err != nil {
				i.log.Error("failed to close indexer",
					zap.Error(err),
				)
			}
			return
		}
		i.txIndices[chainID] = txIndex
	case block.ChainVM:
	default:
		vmType := fmt.Sprintf("%T", vm)
		i.log.Error("got unexpected vm type",
			zap.String("vmType", vmType),
		)
		if err := i.close(); err != nil {
			i.log.Error("failed to close indexer",
				zap.Error(err),
			)
		}
	}
}

func (i *indexer) registerChainHelper(
	chainID ids.ID,
	prefixEnd byte,
	name, endpoint string,
	acceptorGroup snow.AcceptorGroup,
) (Index, error) {
	prefix := make([]byte, ids.IDLen+wrappers.ByteLen)
	copy(prefix, chainID[:])
	prefix[ids.IDLen] = prefixEnd
	indexDB := prefixdb.New(prefix, i.db)
	index, err := newIndex(indexDB, i.log, i.codec, i.clock)
	if err != nil {
		_ = indexDB.Close()
		return nil, err
	}

	// Register index to learn about new accepted vertices
	if err := acceptorGroup.RegisterAcceptor(chainID, fmt.Sprintf("%s%s", indexNamePrefix, chainID), index, true); err != nil {
		_ = index.Close()
		return nil, err
	}

	// Create an API endpoint for this index
	apiServer := rpc.NewServer()
	codec := json.NewCodec()
	apiServer.RegisterCodec(codec, "application/json")
	apiServer.RegisterCodec(codec, "application/json;charset=UTF-8")
	if err := apiServer.RegisterService(&service{Index: index}, "index"); err != nil {
		_ = index.Close()
		return nil, err
	}
	handler := &common.HTTPHandler{LockOptions: common.NoLock, Handler: apiServer}
	if err := i.pathAdder.AddRoute(handler, &sync.RWMutex{}, "index/"+name, "/"+endpoint); err != nil {
		_ = index.Close()
		return nil, err
	}
	return index, nil
}

// Close this indexer. Stops indexing all chains.
// Closes [i.db]. Assumes Close is only called after
// the node is done making decisions.
// Calling Close after it has been called does nothing.
func (i *indexer) Close() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.close()
}

func (i *indexer) close() error {
	if i.closed {
		return nil
	}
	i.closed = true

	errs := &wrappers.Errs{}
	for chainID, txIndex := range i.txIndices {
		errs.Add(
			txIndex.Close(),
			i.txAcceptorGroup.DeregisterAcceptor(chainID, fmt.Sprintf("%s%s", indexNamePrefix, chainID)),
		)
	}
	for chainID, vtxIndex := range i.vtxIndices {
		errs.Add(
			vtxIndex.Close(),
			i.vertexAcceptorGroup.DeregisterAcceptor(chainID, fmt.Sprintf("%s%s", indexNamePrefix, chainID)),
		)
	}
	for chainID, blockIndex := range i.blockIndices {
		errs.Add(
			blockIndex.Close(),
			i.blockAcceptorGroup.DeregisterAcceptor(chainID, fmt.Sprintf("%s%s", indexNamePrefix, chainID)),
		)
	}
	errs.Add(i.db.Close())

	go i.shutdownF()
	return errs.Err
}

func (i *indexer) markIncomplete(chainID ids.ID) error {
	key := make([]byte, ids.IDLen+wrappers.ByteLen)
	copy(key, chainID[:])
	key[ids.IDLen] = isIncompletePrefix
	return i.db.Put(key, nil)
}

// Returns true if this chain is incomplete
func (i *indexer) isIncomplete(chainID ids.ID) (bool, error) {
	key := make([]byte, ids.IDLen+wrappers.ByteLen)
	copy(key, chainID[:])
	key[ids.IDLen] = isIncompletePrefix
	return i.db.Has(key)
}

func (i *indexer) markPreviouslyIndexed(chainID ids.ID) error {
	key := make([]byte, ids.IDLen+wrappers.ByteLen)
	copy(key, chainID[:])
	key[ids.IDLen] = previouslyIndexedPrefix
	return i.db.Put(key, nil)
}

// Returns true if this chain is incomplete
func (i *indexer) previouslyIndexed(chainID ids.ID) (bool, error) {
	key := make([]byte, ids.IDLen+wrappers.ByteLen)
	copy(key, chainID[:])
	key[ids.IDLen] = previouslyIndexedPrefix
	return i.db.Has(key)
}

// Mark that the node has run at least once
func (i *indexer) markHasRun() error {
	return i.db.Put(hasRunKey, nil)
}

// Returns true if the node has run before
func (i *indexer) hasRun() (bool, error) {
	return i.db.Has(hasRunKey)
}
