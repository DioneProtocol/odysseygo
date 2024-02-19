// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	stdmath "math"

	"github.com/google/btree"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/DioneProtocol/odysseygo/cache"
	"github.com/DioneProtocol/odysseygo/cache/metercacher"
	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/linkeddb"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/database/versiondb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/snow/uptime"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/utils/hashing"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/utils/timer"
	"github.com/DioneProtocol/odysseygo/utils/wrappers"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/config"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/fx"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/genesis"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/metrics"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/status"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

const (
	pruneCommitLimit           = 1024
	pruneCommitSleepMultiplier = 5
	pruneCommitSleepCap        = 10 * time.Second
	pruneUpdateFrequency       = 30 * time.Second
)

var (
	_ State = (*state)(nil)

	ErrDelegatorSubset              = errors.New("delegator's time range must be a subset of the validator's time range")
	ErrCantFindSubnet               = errors.New("couldn't find subnet")
	errMissingValidatorSet          = errors.New("missing validator set")
	errValidatorSetAlreadyPopulated = errors.New("validator set already populated")
	errDuplicateValidatorSet        = errors.New("duplicate validator set")
	errIsNotSubnet                  = errors.New("is not a subnet")

	blockIDPrefix                       = []byte("blockID")
	blockPrefix                         = []byte("block")
	validatorsPrefix                    = []byte("validators")
	currentPrefix                       = []byte("current")
	pendingPrefix                       = []byte("pending")
	validatorPrefix                     = []byte("validator")
	delegatorPrefix                     = []byte("delegator")
	subnetValidatorPrefix               = []byte("subnetValidator")
	subnetDelegatorPrefix               = []byte("subnetDelegator")
	nestedValidatorWeightDiffsPrefix    = []byte("validatorDiffs")
	nestedValidatorPublicKeyDiffsPrefix = []byte("publicKeyDiffs")
	flatValidatorWeightDiffsPrefix      = []byte("flatValidatorDiffs")
	flatValidatorPublicKeyDiffsPrefix   = []byte("flatPublicKeyDiffs")
	txPrefix                            = []byte("tx")
	rewardUTXOsPrefix                   = []byte("rewardUTXOs")
	utxoPrefix                          = []byte("utxo")
	subnetPrefix                        = []byte("subnet")
	transformedSubnetPrefix             = []byte("transformedSubnet")
	supplyPrefix                        = []byte("supply")
	chainPrefix                         = []byte("chain")
	singletonPrefix                     = []byte("singleton")

	timestampKey      = []byte("timestamp")
	currentSupplyKey  = []byte("current supply")
	lastAcceptedKey   = []byte("last accepted")
	heightsIndexedKey = []byte("heights indexed")
	initializedKey    = []byte("initialized")
	prunedKey         = []byte("pruned")

	stakeSyncTimestampKey    = []byte("stake sync timestamp")
	stakerMintRateKey        = []byte("staker mint rate")
	feePerWeightStoredKey    = []byte("fee per staker stored")
	lastAccumulatedFeeKey    = []byte("last accumulated fee")
	currentAccumulatedFeeKey = []byte("current accumulated fee")
)

// Chain collects all methods to manage the state of the chain for block
// execution.
type Chain interface {
	Stakers
	dione.UTXOAdder
	dione.UTXOGetter
	dione.UTXODeleter

	GetTimestamp() time.Time
	SetTimestamp(tm time.Time)

	SetStakeSyncTimestamp(tm time.Time)
	GetStakeSyncTimestamp() (time.Time, error)

	SetStakerAccumulatedMintRate(mr *big.Int)
	GetStakerAccumulatedMintRate() (*big.Int, error)

	GetCurrentSupply(subnetID ids.ID) (uint64, error)
	SetCurrentSupply(subnetID ids.ID, cs uint64)

	GetRewardUTXOs(txID ids.ID) ([]*dione.UTXO, error)
	AddRewardUTXO(txID ids.ID, utxo *dione.UTXO)

	GetSubnets() ([]*txs.Tx, error)
	AddSubnet(createSubnetTx *txs.Tx)

	GetSubnetOwner(subnetID ids.ID) (fx.Owner, error)

	GetSubnetTransformation(subnetID ids.ID) (*txs.Tx, error)
	AddSubnetTransformation(transformSubnetTx *txs.Tx)

	GetChains(subnetID ids.ID) ([]*txs.Tx, error)
	AddChain(createChainTx *txs.Tx)

	GetTx(txID ids.ID) (*txs.Tx, status.Status, error)
	AddTx(tx *txs.Tx, status status.Status)

	GetFeePerWeightStored() (*big.Int, error)
	SetFeePerWeightStored(*big.Int)

	GetCurrentAccumulatedFee() (uint64, error)
	AddCurrentAccumulatedFee(uint64)

	GetLastAccumulatedFee() (uint64, error)
	SetLastAccumulatedFee(uint64)
}

type State interface {
	Chain
	uptime.State
	dione.UTXOReader

	GetLastAccepted() ids.ID
	SetLastAccepted(blkID ids.ID)

	GetStatelessBlock(blockID ids.ID) (blocks.Block, error)

	// Invariant: [block] is an accepted block.
	AddStatelessBlock(block blocks.Block)

	GetBlockIDAtHeight(height uint64) (ids.ID, error)

	// ValidatorSet adds all the validators and delegators of [subnetID] into
	// [vdrs].
	ValidatorSet(subnetID ids.ID, vdrs validators.Set) error

	// ApplyValidatorWeightDiffs iterates from [startHeight] towards the genesis
	// block until it has applied all of the diffs up to and including
	// [endHeight]. Applying the diffs modifies [validators].
	//
	// Invariant: If attempting to generate the validator set for
	// [endHeight - 1], [validators] must initially contain the validator
	// weights for [startHeight].
	//
	// Note: Because this function iterates towards the genesis, [startHeight]
	// will typically be greater than or equal to [endHeight]. If [startHeight]
	// is less than [endHeight], no diffs will be applied.
	ApplyValidatorWeightDiffs(
		ctx context.Context,
		validators map[ids.NodeID]*validators.GetValidatorOutput,
		startHeight uint64,
		endHeight uint64,
		subnetID ids.ID,
	) error

	// ApplyValidatorPublicKeyDiffs iterates from [startHeight] towards the
	// genesis block until it has applied all of the diffs up to and including
	// [endHeight]. Applying the diffs modifies [validators].
	//
	// Invariant: If attempting to generate the validator set for
	// [endHeight - 1], [validators] must initially contain the validator
	// weights for [startHeight].
	//
	// Note: Because this function iterates towards the genesis, [startHeight]
	// will typically be greater than or equal to [endHeight]. If [startHeight]
	// is less than [endHeight], no diffs will be applied.
	ApplyValidatorPublicKeyDiffs(
		ctx context.Context,
		validators map[ids.NodeID]*validators.GetValidatorOutput,
		startHeight uint64,
		endHeight uint64,
	) error

	SetHeight(height uint64)

	// Discard uncommitted changes to the database.
	Abort()

	// Returns if the state should be pruned and indexed to remove rejected
	// blocks and generate the block height index.
	//
	// TODO: Remove after v1.11.x is activated
	ShouldPrune() (bool, error)

	// Removes rejected blocks from disk and indexes accepted blocks by height. This
	// function supports being (and is recommended to be) called asynchronously.
	//
	// TODO: Remove after v1.11.x is activated
	PruneAndIndex(sync.Locker, logging.Logger) error

	// Commit changes to the base database.
	Commit() error

	// Returns a batch of unwritten changes that, when written, will commit all
	// pending changes to the base database.
	CommitBatch() (database.Batch, error)

	Checksum() ids.ID

	Close() error
}

// TODO: Remove after v1.11.x is activated
type stateBlk struct {
	Blk    blocks.Block
	Bytes  []byte         `serialize:"true"`
	Status choices.Status `serialize:"true"`
}

/*
 * VMDB
 * |-. validators
 * | |-. current
 * | | |-. validator
 * | | | '-. list
 * | | |   '-- txID -> uptime + potential reward + potential delegatee reward
 * | | |-. delegator
 * | | | '-. list
 * | | |   '-- txID -> potential reward
 * | | |-. subnetValidator
 * | | | '-. list
 * | | |   '-- txID -> uptime + potential reward + potential delegatee reward
 * | | '-. subnetDelegator
 * | |   '-. list
 * | |     '-- txID -> potential reward
 * | |-. pending
 * | | |-. validator
 * | | | '-. list
 * | | |   '-- txID -> nil
 * | | |-. delegator
 * | | | '-. list
 * | | |   '-- txID -> nil
 * | | |-. subnetValidator
 * | | | '-. list
 * | | |   '-- txID -> nil
 * | | '-. subnetDelegator
 * | |   '-. list
 * | |     '-- txID -> nil
 * | |-. nested weight diffs TODO: Remove once only the flat db is needed
 * | | '-. height+subnet
 * | |   '-. list
 * | |     '-- nodeID -> weightChange
 * | |-. nested pub key diffs TODO: Remove once only the flat db is needed
 * | | '-. height
 * | |   '-. list
 * | |     '-- nodeID -> compressed public key
 * | |-. flat weight diffs
 * | | '-- subnet+height+nodeID -> weightChange
 * | '-. flat pub key diffs
 * |   '-- subnet+height+nodeID -> uncompressed public key or nil
 * |-. blockIDs
 * | '-- height -> blockID
 * |-. blocks
 * | '-- blockID -> block bytes
 * |-. txs
 * | '-- txID -> tx bytes + tx status
 * |- rewardUTXOs
 * | '-. txID
 * |   '-. list
 * |     '-- utxoID -> utxo bytes
 * |- utxos
 * | '-- utxoDB
 * |-. subnets
 * | '-. list
 * |   '-- txID -> nil
 * |-. chains
 * | '-. subnetID
 * |   '-. list
 * |     '-- txID -> nil
 * '-. singletons
 *   |-- initializedKey -> nil
 *   |-- prunedKey -> nil
 *   |-- timestampKey -> timestamp
 *   |-- currentSupplyKey -> currentSupply
 *   |-- lastAcceptedKey -> lastAccepted
 *   '-- heightsIndexKey -> startIndexHeight + endIndexHeight
 */
type state struct {
	validatorState

	cfg          *config.Config
	ctx          *snow.Context
	metrics      metrics.Metrics
	rewards      reward.Calculator
	bootstrapped *utils.Atomic[bool]

	baseDB *versiondb.Database

	currentStakers *baseStakers
	pendingStakers *baseStakers

	currentHeight uint64

	addedBlockIDs map[uint64]ids.ID            // map of height -> blockID
	blockIDCache  cache.Cacher[uint64, ids.ID] // cache of height -> blockID. If the entry is ids.Empty, it is not in the database
	blockIDDB     database.Database

	addedBlocks map[ids.ID]blocks.Block            // map of blockID -> Block
	blockCache  cache.Cacher[ids.ID, blocks.Block] // cache of blockID -> Block. If the entry is nil, it is not in the database
	blockDB     database.Database

	validatorsDB                 database.Database
	currentValidatorsDB          database.Database
	currentValidatorBaseDB       database.Database
	currentValidatorList         linkeddb.LinkedDB
	currentDelegatorBaseDB       database.Database
	currentDelegatorList         linkeddb.LinkedDB
	currentSubnetValidatorBaseDB database.Database
	currentSubnetValidatorList   linkeddb.LinkedDB
	currentSubnetDelegatorBaseDB database.Database
	currentSubnetDelegatorList   linkeddb.LinkedDB
	pendingValidatorsDB          database.Database
	pendingValidatorBaseDB       database.Database
	pendingValidatorList         linkeddb.LinkedDB
	pendingDelegatorBaseDB       database.Database
	pendingDelegatorList         linkeddb.LinkedDB
	pendingSubnetValidatorBaseDB database.Database
	pendingSubnetValidatorList   linkeddb.LinkedDB
	pendingSubnetDelegatorBaseDB database.Database
	pendingSubnetDelegatorList   linkeddb.LinkedDB

	nestedValidatorWeightDiffsDB    database.Database
	nestedValidatorPublicKeyDiffsDB database.Database
	flatValidatorWeightDiffsDB      database.Database
	flatValidatorPublicKeyDiffsDB   database.Database

	addedTxs map[ids.ID]*txAndStatus            // map of txID -> {*txs.Tx, Status}
	txCache  cache.Cacher[ids.ID, *txAndStatus] // txID -> {*txs.Tx, Status}. If the entry is nil, it isn't in the database
	txDB     database.Database

	addedRewardUTXOs map[ids.ID][]*dione.UTXO            // map of txID -> []*UTXO
	rewardUTXOsCache cache.Cacher[ids.ID, []*dione.UTXO] // txID -> []*UTXO
	rewardUTXODB     database.Database

	modifiedUTXOs map[ids.ID]*dione.UTXO // map of modified UTXOID -> *UTXO if the UTXO is nil, it has been removed
	utxoDB        database.Database
	utxoState     dione.UTXOState

	cachedSubnets []*txs.Tx // nil if the subnets haven't been loaded
	addedSubnets  []*txs.Tx
	subnetBaseDB  database.Database
	subnetDB      linkeddb.LinkedDB

	transformedSubnets     map[ids.ID]*txs.Tx            // map of subnetID -> transformSubnetTx
	transformedSubnetCache cache.Cacher[ids.ID, *txs.Tx] // cache of subnetID -> transformSubnetTx if the entry is nil, it is not in the database
	transformedSubnetDB    database.Database

	modifiedSupplies map[ids.ID]uint64             // map of subnetID -> current supply
	supplyCache      cache.Cacher[ids.ID, *uint64] // cache of subnetID -> current supply if the entry is nil, it is not in the database
	supplyDB         database.Database

	addedChains  map[ids.ID][]*txs.Tx                    // maps subnetID -> the newly added chains to the subnet
	chainCache   cache.Cacher[ids.ID, []*txs.Tx]         // cache of subnetID -> the chains after all local modifications []*txs.Tx
	chainDBCache cache.Cacher[ids.ID, linkeddb.LinkedDB] // cache of subnetID -> linkedDB
	chainDB      database.Database

	// The persisted fields represent the current database value
	timestamp, persistedTimestamp         time.Time
	currentSupply, persistedCurrentSupply uint64
	// [lastAccepted] is the most recently accepted block.
	lastAccepted, persistedLastAccepted ids.ID
	indexedHeights                      *heightRange
	singletonDB                         database.Database

	lastAccumulatedFee, persistedLastAccumulatedFee               uint64
	currentAccumulatedFee, persistedCurrentAccumulatedFee         uint64
	feePerWeightStored, persistedFeePerWeightStored               *big.Int
	stakerAccumulatedMintRate, persistedStakerAccumulatedMintRate *big.Int
	stakeSyncTimestamp, persistedStakeSyncTimestamp               time.Time
}

// heightRange is used to track which heights are safe to use the native DB
// iterator for querying validator diffs.
//
// TODO: Remove once we are guaranteed nodes can not rollback to not support the
// new indexing mechanism.
type heightRange struct {
	LowerBound uint64 `serialize:"true"`
	UpperBound uint64 `serialize:"true"`
}

type ValidatorWeightDiff struct {
	Decrease bool   `serialize:"true"`
	Amount   uint64 `serialize:"true"`
}

func (v *ValidatorWeightDiff) Add(negative bool, amount uint64) error {
	if v.Decrease == negative {
		var err error
		v.Amount, err = math.Add64(v.Amount, amount)
		return err
	}

	if v.Amount > amount {
		v.Amount -= amount
	} else {
		v.Amount = math.AbsDiff(v.Amount, amount)
		v.Decrease = negative
	}
	return nil
}

type heightWithSubnet struct {
	Height   uint64 `serialize:"true"`
	SubnetID ids.ID `serialize:"true"`
}

type txBytesAndStatus struct {
	Tx     []byte        `serialize:"true"`
	Status status.Status `serialize:"true"`
}

type txAndStatus struct {
	tx     *txs.Tx
	status status.Status
}

func txSize(_ ids.ID, tx *txs.Tx) int {
	if tx == nil {
		return ids.IDLen + constants.PointerOverhead
	}
	return ids.IDLen + len(tx.Bytes()) + constants.PointerOverhead
}

func txAndStatusSize(_ ids.ID, t *txAndStatus) int {
	if t == nil {
		return ids.IDLen + constants.PointerOverhead
	}
	return ids.IDLen + len(t.tx.Bytes()) + wrappers.IntLen + 2*constants.PointerOverhead
}

func blockSize(_ ids.ID, blk blocks.Block) int {
	if blk == nil {
		return ids.IDLen + constants.PointerOverhead
	}
	return ids.IDLen + len(blk.Bytes()) + constants.PointerOverhead
}

func New(
	db database.Database,
	genesisBytes []byte,
	metricsReg prometheus.Registerer,
	cfg *config.Config,
	execCfg *config.ExecutionConfig,
	ctx *snow.Context,
	metrics metrics.Metrics,
	rewards reward.Calculator,
	bootstrapped *utils.Atomic[bool],
) (State, error) {
	s, err := newState(
		db,
		metrics,
		cfg,
		execCfg,
		ctx,
		metricsReg,
		rewards,
		bootstrapped,
	)
	if err != nil {
		return nil, err
	}

	if err := s.sync(genesisBytes); err != nil {
		// Drop any errors on close to return the first error
		_ = s.Close()

		return nil, err
	}

	// Before we start accepting new blocks, we check if the pruning process needs
	// to be run.
	//
	// TODO: Cleanup after v1.11.x is activated
	shouldPrune, err := s.ShouldPrune()
	if err != nil {
		return nil, err
	}
	if shouldPrune {
		// If the pruned key is on disk, we must delete it to ensure our disk
		// can't get into a partially pruned state if the node restarts mid-way
		// through pruning.
		if err := s.singletonDB.Delete(prunedKey); err != nil {
			return nil, fmt.Errorf("failed to remove prunedKey from singletonDB: %w", err)
		}

		if err := s.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit to baseDB: %w", err)
		}
	}

	return s, nil
}

func newState(
	db database.Database,
	metrics metrics.Metrics,
	cfg *config.Config,
	execCfg *config.ExecutionConfig,
	ctx *snow.Context,
	metricsReg prometheus.Registerer,
	rewards reward.Calculator,
	bootstrapped *utils.Atomic[bool],
) (*state, error) {
	blockIDCache, err := metercacher.New[uint64, ids.ID](
		"block_id_cache",
		metricsReg,
		&cache.LRU[uint64, ids.ID]{Size: execCfg.BlockIDCacheSize},
	)
	if err != nil {
		return nil, err
	}

	blockCache, err := metercacher.New[ids.ID, blocks.Block](
		"block_cache",
		metricsReg,
		cache.NewSizedLRU[ids.ID, blocks.Block](execCfg.BlockCacheSize, blockSize),
	)
	if err != nil {
		return nil, err
	}

	baseDB := versiondb.New(db)

	validatorsDB := prefixdb.New(validatorsPrefix, baseDB)

	currentValidatorsDB := prefixdb.New(currentPrefix, validatorsDB)
	currentValidatorBaseDB := prefixdb.New(validatorPrefix, currentValidatorsDB)
	currentDelegatorBaseDB := prefixdb.New(delegatorPrefix, currentValidatorsDB)
	currentSubnetValidatorBaseDB := prefixdb.New(subnetValidatorPrefix, currentValidatorsDB)
	currentSubnetDelegatorBaseDB := prefixdb.New(subnetDelegatorPrefix, currentValidatorsDB)

	pendingValidatorsDB := prefixdb.New(pendingPrefix, validatorsDB)
	pendingValidatorBaseDB := prefixdb.New(validatorPrefix, pendingValidatorsDB)
	pendingDelegatorBaseDB := prefixdb.New(delegatorPrefix, pendingValidatorsDB)
	pendingSubnetValidatorBaseDB := prefixdb.New(subnetValidatorPrefix, pendingValidatorsDB)
	pendingSubnetDelegatorBaseDB := prefixdb.New(subnetDelegatorPrefix, pendingValidatorsDB)

	nestedValidatorWeightDiffsDB := prefixdb.New(nestedValidatorWeightDiffsPrefix, validatorsDB)
	nestedValidatorPublicKeyDiffsDB := prefixdb.New(nestedValidatorPublicKeyDiffsPrefix, validatorsDB)
	flatValidatorWeightDiffsDB := prefixdb.New(flatValidatorWeightDiffsPrefix, validatorsDB)
	flatValidatorPublicKeyDiffsDB := prefixdb.New(flatValidatorPublicKeyDiffsPrefix, validatorsDB)

	txCache, err := metercacher.New(
		"tx_cache",
		metricsReg,
		cache.NewSizedLRU[ids.ID, *txAndStatus](execCfg.TxCacheSize, txAndStatusSize),
	)
	if err != nil {
		return nil, err
	}

	rewardUTXODB := prefixdb.New(rewardUTXOsPrefix, baseDB)
	rewardUTXOsCache, err := metercacher.New[ids.ID, []*dione.UTXO](
		"reward_utxos_cache",
		metricsReg,
		&cache.LRU[ids.ID, []*dione.UTXO]{Size: execCfg.RewardUTXOsCacheSize},
	)
	if err != nil {
		return nil, err
	}

	utxoDB := prefixdb.New(utxoPrefix, baseDB)
	utxoState, err := dione.NewMeteredUTXOState(utxoDB, txs.GenesisCodec, metricsReg, execCfg.ChecksumsEnabled)
	if err != nil {
		return nil, err
	}

	subnetBaseDB := prefixdb.New(subnetPrefix, baseDB)

	transformedSubnetCache, err := metercacher.New(
		"transformed_subnet_cache",
		metricsReg,
		cache.NewSizedLRU[ids.ID, *txs.Tx](execCfg.TransformedSubnetTxCacheSize, txSize),
	)
	if err != nil {
		return nil, err
	}

	supplyCache, err := metercacher.New[ids.ID, *uint64](
		"supply_cache",
		metricsReg,
		&cache.LRU[ids.ID, *uint64]{Size: execCfg.ChainCacheSize},
	)
	if err != nil {
		return nil, err
	}

	chainCache, err := metercacher.New[ids.ID, []*txs.Tx](
		"chain_cache",
		metricsReg,
		&cache.LRU[ids.ID, []*txs.Tx]{Size: execCfg.ChainCacheSize},
	)
	if err != nil {
		return nil, err
	}

	chainDBCache, err := metercacher.New[ids.ID, linkeddb.LinkedDB](
		"chain_db_cache",
		metricsReg,
		&cache.LRU[ids.ID, linkeddb.LinkedDB]{Size: execCfg.ChainDBCacheSize},
	)
	if err != nil {
		return nil, err
	}

	return &state{
		validatorState: newValidatorState(),

		cfg:          cfg,
		ctx:          ctx,
		metrics:      metrics,
		rewards:      rewards,
		bootstrapped: bootstrapped,
		baseDB:       baseDB,

		addedBlockIDs: make(map[uint64]ids.ID),
		blockIDCache:  blockIDCache,
		blockIDDB:     prefixdb.New(blockIDPrefix, baseDB),

		addedBlocks: make(map[ids.ID]blocks.Block),
		blockCache:  blockCache,
		blockDB:     prefixdb.New(blockPrefix, baseDB),

		currentStakers: newBaseStakers(),
		pendingStakers: newBaseStakers(),

		validatorsDB:                    validatorsDB,
		currentValidatorsDB:             currentValidatorsDB,
		currentValidatorBaseDB:          currentValidatorBaseDB,
		currentValidatorList:            linkeddb.NewDefault(currentValidatorBaseDB),
		currentDelegatorBaseDB:          currentDelegatorBaseDB,
		currentDelegatorList:            linkeddb.NewDefault(currentDelegatorBaseDB),
		currentSubnetValidatorBaseDB:    currentSubnetValidatorBaseDB,
		currentSubnetValidatorList:      linkeddb.NewDefault(currentSubnetValidatorBaseDB),
		currentSubnetDelegatorBaseDB:    currentSubnetDelegatorBaseDB,
		currentSubnetDelegatorList:      linkeddb.NewDefault(currentSubnetDelegatorBaseDB),
		pendingValidatorsDB:             pendingValidatorsDB,
		pendingValidatorBaseDB:          pendingValidatorBaseDB,
		pendingValidatorList:            linkeddb.NewDefault(pendingValidatorBaseDB),
		pendingDelegatorBaseDB:          pendingDelegatorBaseDB,
		pendingDelegatorList:            linkeddb.NewDefault(pendingDelegatorBaseDB),
		pendingSubnetValidatorBaseDB:    pendingSubnetValidatorBaseDB,
		pendingSubnetValidatorList:      linkeddb.NewDefault(pendingSubnetValidatorBaseDB),
		pendingSubnetDelegatorBaseDB:    pendingSubnetDelegatorBaseDB,
		pendingSubnetDelegatorList:      linkeddb.NewDefault(pendingSubnetDelegatorBaseDB),
		nestedValidatorWeightDiffsDB:    nestedValidatorWeightDiffsDB,
		nestedValidatorPublicKeyDiffsDB: nestedValidatorPublicKeyDiffsDB,
		flatValidatorWeightDiffsDB:      flatValidatorWeightDiffsDB,
		flatValidatorPublicKeyDiffsDB:   flatValidatorPublicKeyDiffsDB,

		addedTxs: make(map[ids.ID]*txAndStatus),
		txDB:     prefixdb.New(txPrefix, baseDB),
		txCache:  txCache,

		addedRewardUTXOs: make(map[ids.ID][]*dione.UTXO),
		rewardUTXODB:     rewardUTXODB,
		rewardUTXOsCache: rewardUTXOsCache,

		modifiedUTXOs: make(map[ids.ID]*dione.UTXO),
		utxoDB:        utxoDB,
		utxoState:     utxoState,

		subnetBaseDB: subnetBaseDB,
		subnetDB:     linkeddb.NewDefault(subnetBaseDB),

		transformedSubnets:     make(map[ids.ID]*txs.Tx),
		transformedSubnetCache: transformedSubnetCache,
		transformedSubnetDB:    prefixdb.New(transformedSubnetPrefix, baseDB),

		modifiedSupplies: make(map[ids.ID]uint64),
		supplyCache:      supplyCache,
		supplyDB:         prefixdb.New(supplyPrefix, baseDB),

		addedChains:  make(map[ids.ID][]*txs.Tx),
		chainDB:      prefixdb.New(chainPrefix, baseDB),
		chainCache:   chainCache,
		chainDBCache: chainDBCache,

		singletonDB: prefixdb.New(singletonPrefix, baseDB),

		stakerAccumulatedMintRate:          new(big.Int),
		persistedStakerAccumulatedMintRate: new(big.Int),
		feePerWeightStored:                 new(big.Int),
		persistedFeePerWeightStored:        new(big.Int),
	}, nil
}

func (s *state) GetCurrentValidator(subnetID ids.ID, nodeID ids.NodeID) (*Staker, error) {
	return s.currentStakers.GetValidator(subnetID, nodeID)
}

func (s *state) PutCurrentValidator(staker *Staker) {
	s.currentStakers.PutValidator(staker)
}

func (s *state) DeleteCurrentValidator(staker *Staker) {
	s.currentStakers.DeleteValidator(staker)
}

func (s *state) GetCurrentDelegatorIterator(subnetID ids.ID, nodeID ids.NodeID) (StakerIterator, error) {
	return s.currentStakers.GetDelegatorIterator(subnetID, nodeID), nil
}

func (s *state) PutCurrentDelegator(staker *Staker) {
	s.currentStakers.PutDelegator(staker)
}

func (s *state) DeleteCurrentDelegator(staker *Staker) {
	s.currentStakers.DeleteDelegator(staker)
}

func (s *state) GetCurrentStakerIterator() (StakerIterator, error) {
	return s.currentStakers.GetStakerIterator(), nil
}

func (s *state) GetCurrentStakersLen() (uint64, error) {
	return uint64(s.currentStakers.stakers.Len()), nil
}

func (s *state) GetPendingValidator(subnetID ids.ID, nodeID ids.NodeID) (*Staker, error) {
	return s.pendingStakers.GetValidator(subnetID, nodeID)
}

func (s *state) PutPendingValidator(staker *Staker) {
	s.pendingStakers.PutValidator(staker)
}

func (s *state) DeletePendingValidator(staker *Staker) {
	s.pendingStakers.DeleteValidator(staker)
}

func (s *state) GetPendingDelegatorIterator(subnetID ids.ID, nodeID ids.NodeID) (StakerIterator, error) {
	return s.pendingStakers.GetDelegatorIterator(subnetID, nodeID), nil
}

func (s *state) PutPendingDelegator(staker *Staker) {
	s.pendingStakers.PutDelegator(staker)
}

func (s *state) DeletePendingDelegator(staker *Staker) {
	s.pendingStakers.DeleteDelegator(staker)
}

func (s *state) GetPendingStakerIterator() (StakerIterator, error) {
	return s.pendingStakers.GetStakerIterator(), nil
}

func (s *state) GetPendingStakersLen() (uint64, error) {
	return uint64(s.pendingStakers.stakers.Len()), nil
}

func (s *state) shouldInit() (bool, error) {
	has, err := s.singletonDB.Has(initializedKey)
	return !has, err
}

func (s *state) doneInit() error {
	return s.singletonDB.Put(initializedKey, nil)
}

func (s *state) ShouldPrune() (bool, error) {
	has, err := s.singletonDB.Has(prunedKey)
	if err != nil {
		return true, err
	}

	// If [prunedKey] is not in [singletonDB], [PruneAndIndex()] did not finish
	// execution.
	if !has {
		return true, nil
	}

	// To ensure the db was not modified since we last ran [PruneAndIndex()], we
	// must verify that [s.lastAccepted] is height indexed.
	blk, err := s.GetStatelessBlock(s.lastAccepted)
	if err != nil {
		return true, err
	}

	_, err = s.GetBlockIDAtHeight(blk.Height())
	if err == database.ErrNotFound {
		return true, nil
	}

	return false, err
}

func (s *state) donePrune() error {
	return s.singletonDB.Put(prunedKey, nil)
}

func (s *state) GetSubnets() ([]*txs.Tx, error) {
	if s.cachedSubnets != nil {
		return s.cachedSubnets, nil
	}

	subnetDBIt := s.subnetDB.NewIterator()
	defer subnetDBIt.Release()

	txs := []*txs.Tx(nil)
	for subnetDBIt.Next() {
		subnetIDBytes := subnetDBIt.Key()
		subnetID, err := ids.ToID(subnetIDBytes)
		if err != nil {
			return nil, err
		}
		subnetTx, _, err := s.GetTx(subnetID)
		if err != nil {
			return nil, err
		}
		txs = append(txs, subnetTx)
	}
	if err := subnetDBIt.Error(); err != nil {
		return nil, err
	}
	txs = append(txs, s.addedSubnets...)
	s.cachedSubnets = txs
	return txs, nil
}

func (s *state) AddSubnet(createSubnetTx *txs.Tx) {
	s.addedSubnets = append(s.addedSubnets, createSubnetTx)
	if s.cachedSubnets != nil {
		s.cachedSubnets = append(s.cachedSubnets, createSubnetTx)
	}
}

func (s *state) GetSubnetOwner(subnetID ids.ID) (fx.Owner, error) {
	subnetIntf, _, err := s.GetTx(subnetID)
	if err != nil {
		return nil, fmt.Errorf(
			"%w %q: %w",
			ErrCantFindSubnet,
			subnetID,
			err,
		)
	}

	subnet, ok := subnetIntf.Unsigned.(*txs.CreateSubnetTx)
	if !ok {
		return nil, fmt.Errorf("%q %w", subnetID, errIsNotSubnet)
	}

	return subnet.Owner, nil
}

func (s *state) GetSubnetTransformation(subnetID ids.ID) (*txs.Tx, error) {
	if tx, exists := s.transformedSubnets[subnetID]; exists {
		return tx, nil
	}

	if tx, cached := s.transformedSubnetCache.Get(subnetID); cached {
		if tx == nil {
			return nil, database.ErrNotFound
		}
		return tx, nil
	}

	transformSubnetTxID, err := database.GetID(s.transformedSubnetDB, subnetID[:])
	if err == database.ErrNotFound {
		s.transformedSubnetCache.Put(subnetID, nil)
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	transformSubnetTx, _, err := s.GetTx(transformSubnetTxID)
	if err != nil {
		return nil, err
	}
	s.transformedSubnetCache.Put(subnetID, transformSubnetTx)
	return transformSubnetTx, nil
}

func (s *state) AddSubnetTransformation(transformSubnetTxIntf *txs.Tx) {
	transformSubnetTx := transformSubnetTxIntf.Unsigned.(*txs.TransformSubnetTx)
	s.transformedSubnets[transformSubnetTx.Subnet] = transformSubnetTxIntf
}

func (s *state) GetChains(subnetID ids.ID) ([]*txs.Tx, error) {
	if chains, cached := s.chainCache.Get(subnetID); cached {
		return chains, nil
	}
	chainDB := s.getChainDB(subnetID)
	chainDBIt := chainDB.NewIterator()
	defer chainDBIt.Release()

	txs := []*txs.Tx(nil)
	for chainDBIt.Next() {
		chainIDBytes := chainDBIt.Key()
		chainID, err := ids.ToID(chainIDBytes)
		if err != nil {
			return nil, err
		}
		chainTx, _, err := s.GetTx(chainID)
		if err != nil {
			return nil, err
		}
		txs = append(txs, chainTx)
	}
	if err := chainDBIt.Error(); err != nil {
		return nil, err
	}
	txs = append(txs, s.addedChains[subnetID]...)
	s.chainCache.Put(subnetID, txs)
	return txs, nil
}

func (s *state) AddChain(createChainTxIntf *txs.Tx) {
	createChainTx := createChainTxIntf.Unsigned.(*txs.CreateChainTx)
	subnetID := createChainTx.SubnetID
	s.addedChains[subnetID] = append(s.addedChains[subnetID], createChainTxIntf)
	if chains, cached := s.chainCache.Get(subnetID); cached {
		chains = append(chains, createChainTxIntf)
		s.chainCache.Put(subnetID, chains)
	}
}

func (s *state) getChainDB(subnetID ids.ID) linkeddb.LinkedDB {
	if chainDB, cached := s.chainDBCache.Get(subnetID); cached {
		return chainDB
	}
	rawChainDB := prefixdb.New(subnetID[:], s.chainDB)
	chainDB := linkeddb.NewDefault(rawChainDB)
	s.chainDBCache.Put(subnetID, chainDB)
	return chainDB
}

func (s *state) GetTx(txID ids.ID) (*txs.Tx, status.Status, error) {
	if tx, exists := s.addedTxs[txID]; exists {
		return tx.tx, tx.status, nil
	}
	if tx, cached := s.txCache.Get(txID); cached {
		if tx == nil {
			return nil, status.Unknown, database.ErrNotFound
		}
		return tx.tx, tx.status, nil
	}
	txBytes, err := s.txDB.Get(txID[:])
	if err == database.ErrNotFound {
		s.txCache.Put(txID, nil)
		return nil, status.Unknown, database.ErrNotFound
	} else if err != nil {
		return nil, status.Unknown, err
	}

	stx := txBytesAndStatus{}
	if _, err := txs.GenesisCodec.Unmarshal(txBytes, &stx); err != nil {
		return nil, status.Unknown, err
	}

	tx, err := txs.Parse(txs.GenesisCodec, stx.Tx)
	if err != nil {
		return nil, status.Unknown, err
	}

	ptx := &txAndStatus{
		tx:     tx,
		status: stx.Status,
	}

	s.txCache.Put(txID, ptx)
	return ptx.tx, ptx.status, nil
}

func (s *state) AddTx(tx *txs.Tx, status status.Status) {
	s.addedTxs[tx.ID()] = &txAndStatus{
		tx:     tx,
		status: status,
	}
}

func (s *state) GetRewardUTXOs(txID ids.ID) ([]*dione.UTXO, error) {
	if utxos, exists := s.addedRewardUTXOs[txID]; exists {
		return utxos, nil
	}
	if utxos, exists := s.rewardUTXOsCache.Get(txID); exists {
		return utxos, nil
	}

	rawTxDB := prefixdb.New(txID[:], s.rewardUTXODB)
	txDB := linkeddb.NewDefault(rawTxDB)
	it := txDB.NewIterator()
	defer it.Release()

	utxos := []*dione.UTXO(nil)
	for it.Next() {
		utxo := &dione.UTXO{}
		if _, err := txs.Codec.Unmarshal(it.Value(), utxo); err != nil {
			return nil, err
		}
		utxos = append(utxos, utxo)
	}
	if err := it.Error(); err != nil {
		return nil, err
	}

	s.rewardUTXOsCache.Put(txID, utxos)
	return utxos, nil
}

func (s *state) AddRewardUTXO(txID ids.ID, utxo *dione.UTXO) {
	s.addedRewardUTXOs[txID] = append(s.addedRewardUTXOs[txID], utxo)
}

func (s *state) GetUTXO(utxoID ids.ID) (*dione.UTXO, error) {
	if utxo, exists := s.modifiedUTXOs[utxoID]; exists {
		if utxo == nil {
			return nil, database.ErrNotFound
		}
		return utxo, nil
	}
	return s.utxoState.GetUTXO(utxoID)
}

func (s *state) UTXOIDs(addr []byte, start ids.ID, limit int) ([]ids.ID, error) {
	return s.utxoState.UTXOIDs(addr, start, limit)
}

func (s *state) AddUTXO(utxo *dione.UTXO) {
	s.modifiedUTXOs[utxo.InputID()] = utxo
}

func (s *state) DeleteUTXO(utxoID ids.ID) {
	s.modifiedUTXOs[utxoID] = nil
}

func (s *state) GetStartTime(nodeID ids.NodeID, subnetID ids.ID) (time.Time, error) {
	staker, err := s.currentStakers.GetValidator(subnetID, nodeID)
	if err != nil {
		return time.Time{}, err
	}
	return staker.StartTime, nil
}

func (s *state) GetTimestamp() time.Time {
	return s.timestamp
}

func (s *state) SetTimestamp(tm time.Time) {
	s.timestamp = tm
}

func (s *state) GetStakeSyncTimestamp() (time.Time, error) {
	return s.stakeSyncTimestamp, nil
}

func (s *state) SetStakeSyncTimestamp(tm time.Time) {
	s.stakeSyncTimestamp = tm
}

func (s *state) GetStakerAccumulatedMintRate() (*big.Int, error) {
	return new(big.Int).Set(s.stakerAccumulatedMintRate), nil
}

func (s *state) SetStakerAccumulatedMintRate(mr *big.Int) {
	s.stakerAccumulatedMintRate.Set(mr)
}

func (s *state) GetLastAccepted() ids.ID {
	return s.lastAccepted
}

func (s *state) SetLastAccepted(lastAccepted ids.ID) {
	s.lastAccepted = lastAccepted
}

func (s *state) GetCurrentSupply(subnetID ids.ID) (uint64, error) {
	if subnetID == constants.PrimaryNetworkID {
		return s.currentSupply, nil
	}

	supply, ok := s.modifiedSupplies[subnetID]
	if ok {
		return supply, nil
	}

	cachedSupply, ok := s.supplyCache.Get(subnetID)
	if ok {
		if cachedSupply == nil {
			return 0, database.ErrNotFound
		}
		return *cachedSupply, nil
	}

	supply, err := database.GetUInt64(s.supplyDB, subnetID[:])
	if err == database.ErrNotFound {
		s.supplyCache.Put(subnetID, nil)
		return 0, database.ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	s.supplyCache.Put(subnetID, &supply)
	return supply, nil
}

func (s *state) SetCurrentSupply(subnetID ids.ID, cs uint64) {
	if subnetID == constants.PrimaryNetworkID {
		s.currentSupply = cs
	} else {
		s.modifiedSupplies[subnetID] = cs
	}
}

func (s *state) ValidatorSet(subnetID ids.ID, vdrs validators.Set) error {
	for nodeID, validator := range s.currentStakers.validators[subnetID] {
		staker := validator.validator
		if err := vdrs.Add(nodeID, staker.PublicKey, staker.TxID, staker.Weight); err != nil {
			return err
		}

		delegatorIterator := NewTreeIterator(validator.delegators)
		for delegatorIterator.Next() {
			staker := delegatorIterator.Value()
			if err := vdrs.AddWeight(nodeID, staker.Weight); err != nil {
				delegatorIterator.Release()
				return err
			}
		}
		delegatorIterator.Release()
	}
	return nil
}

func (s *state) ApplyValidatorWeightDiffs(
	ctx context.Context,
	validators map[ids.NodeID]*validators.GetValidatorOutput,
	startHeight uint64,
	endHeight uint64,
	subnetID ids.ID,
) error {
	diffIter := s.flatValidatorWeightDiffsDB.NewIteratorWithStartAndPrefix(
		marshalStartDiffKey(subnetID, startHeight),
		subnetID[:],
	)
	defer diffIter.Release()

	prevHeight := startHeight + 1
	// TODO: Remove the index continuity checks once we are guaranteed nodes can
	// not rollback to not support the new indexing mechanism.
	for diffIter.Next() && s.indexedHeights != nil && s.indexedHeights.LowerBound <= endHeight {
		if err := ctx.Err(); err != nil {
			return err
		}

		_, parsedHeight, nodeID, err := unmarshalDiffKey(diffIter.Key())
		if err != nil {
			return err
		}
		// If the parsedHeight is less than our target endHeight, then we have
		// fully processed the diffs from startHeight through endHeight.
		if parsedHeight < endHeight {
			return diffIter.Error()
		}

		prevHeight = parsedHeight

		weightDiff, err := unmarshalWeightDiff(diffIter.Value())
		if err != nil {
			return err
		}

		if err := applyWeightDiff(validators, nodeID, weightDiff); err != nil {
			return err
		}
	}
	if err := diffIter.Error(); err != nil {
		return err
	}

	// TODO: Remove this once it is assumed that all subnet validators have
	// adopted the new indexing.
	for height := prevHeight - 1; height >= endHeight; height-- {
		if err := ctx.Err(); err != nil {
			return err
		}

		prefixStruct := heightWithSubnet{
			Height:   height,
			SubnetID: subnetID,
		}
		prefixBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, prefixStruct)
		if err != nil {
			return err
		}

		rawDiffDB := prefixdb.New(prefixBytes, s.nestedValidatorWeightDiffsDB)
		diffDB := linkeddb.NewDefault(rawDiffDB)
		diffIter := diffDB.NewIterator()
		defer diffIter.Release()

		for diffIter.Next() {
			nodeID, err := ids.ToNodeID(diffIter.Key())
			if err != nil {
				return err
			}

			weightDiff := ValidatorWeightDiff{}
			_, err = blocks.GenesisCodec.Unmarshal(diffIter.Value(), &weightDiff)
			if err != nil {
				return err
			}

			if err := applyWeightDiff(validators, nodeID, &weightDiff); err != nil {
				return err
			}
		}
	}

	return nil
}

func applyWeightDiff(
	vdrs map[ids.NodeID]*validators.GetValidatorOutput,
	nodeID ids.NodeID,
	weightDiff *ValidatorWeightDiff,
) error {
	vdr, ok := vdrs[nodeID]
	if !ok {
		// This node isn't in the current validator set.
		vdr = &validators.GetValidatorOutput{
			NodeID: nodeID,
		}
		vdrs[nodeID] = vdr
	}

	// The weight of this node changed at this block.
	var err error
	if weightDiff.Decrease {
		// The validator's weight was decreased at this block, so in the
		// prior block it was higher.
		vdr.Weight, err = math.Add64(vdr.Weight, weightDiff.Amount)
	} else {
		// The validator's weight was increased at this block, so in the
		// prior block it was lower.
		vdr.Weight, err = math.Sub(vdr.Weight, weightDiff.Amount)
	}
	if err != nil {
		return err
	}

	if vdr.Weight == 0 {
		// The validator's weight was 0 before this block so they weren't in the
		// validator set.
		delete(vdrs, nodeID)
	}
	return nil
}

func (s *state) ApplyValidatorPublicKeyDiffs(
	ctx context.Context,
	validators map[ids.NodeID]*validators.GetValidatorOutput,
	startHeight uint64,
	endHeight uint64,
) error {
	diffIter := s.flatValidatorPublicKeyDiffsDB.NewIteratorWithStartAndPrefix(
		marshalStartDiffKey(constants.PrimaryNetworkID, startHeight),
		constants.PrimaryNetworkID[:],
	)
	defer diffIter.Release()

	for diffIter.Next() {
		if err := ctx.Err(); err != nil {
			return err
		}

		_, parsedHeight, nodeID, err := unmarshalDiffKey(diffIter.Key())
		if err != nil {
			return err
		}
		// If the parsedHeight is less than our target endHeight, then we have
		// fully processed the diffs from startHeight through endHeight.
		if parsedHeight < endHeight {
			break
		}

		vdr, ok := validators[nodeID]
		if !ok {
			continue
		}

		pkBytes := diffIter.Value()
		if len(pkBytes) == 0 {
			vdr.PublicKey = nil
			continue
		}

		vdr.PublicKey = new(bls.PublicKey).Deserialize(pkBytes)
	}

	// Note: this does not fallback to the linkeddb index because the linkeddb
	// index does not contain entries for when to remove the public key.
	//
	// Nodes may see inconsistent public keys for heights before the new public
	// key index was populated.
	return diffIter.Error()
}

func (s *state) syncGenesis(genesisBlk blocks.Block, genesis *genesis.State) error {
	genesisBlkID := genesisBlk.ID()
	s.SetLastAccepted(genesisBlkID)
	s.SetTimestamp(time.Unix(int64(genesis.Timestamp), 0))
	s.SetStakeSyncTimestamp(s.GetTimestamp())
	s.SetCurrentSupply(constants.PrimaryNetworkID, genesis.InitialSupply)
	s.AddStatelessBlock(genesisBlk)

	// Persist UTXOs that exist at genesis
	for _, utxo := range genesis.UTXOs {
		s.AddUTXO(utxo)
	}

	// Persist primary network validator set at genesis
	for _, vdrTx := range genesis.Validators {
		tx, ok := vdrTx.Unsigned.(*txs.AddValidatorTx)
		if !ok {
			return fmt.Errorf("expected tx type *txs.AddValidatorTx but got %T", vdrTx.Unsigned)
		}

		staker, err := NewCurrentStaker(vdrTx.ID(), tx, 0)
		if err != nil {
			return err
		}

		s.PutCurrentValidator(staker)
		s.AddTx(vdrTx, status.Committed)
	}

	for _, chain := range genesis.Chains {
		unsignedChain, ok := chain.Unsigned.(*txs.CreateChainTx)
		if !ok {
			return fmt.Errorf("expected tx type *txs.CreateChainTx but got %T", chain.Unsigned)
		}

		// Ensure all chains that the genesis bytes say to create have the right
		// network ID
		if unsignedChain.NetworkID != s.ctx.NetworkID {
			return dione.ErrWrongNetworkID
		}

		s.AddChain(chain)
		s.AddTx(chain, status.Committed)
	}

	// updateValidators is set to false here to maintain the invariant that the
	// primary network's validator set is empty before the validator sets are
	// initialized.
	return s.write(false /*=updateValidators*/, 0)
}

// Load pulls data previously stored on disk that is expected to be in memory.
func (s *state) load() error {
	errs := wrappers.Errs{}
	errs.Add(
		s.loadMetadata(),
		s.loadCurrentValidators(),
		s.loadPendingValidators(),
		s.initValidatorSets(),
	)
	return errs.Err
}

func (s *state) loadMetadata() error {
	timestamp, err := database.GetTimestamp(s.singletonDB, timestampKey)
	if err != nil {
		return err
	}
	s.persistedTimestamp = timestamp
	s.SetTimestamp(timestamp)

	stakeSyncTimestamp, err := database.GetTimestamp(s.singletonDB, stakeSyncTimestampKey)
	if err != nil {
		return err
	}
	s.persistedStakeSyncTimestamp = stakeSyncTimestamp
	s.SetStakeSyncTimestamp(stakeSyncTimestamp)

	stakerMintRate, err := database.GetBigInt(s.singletonDB, stakerMintRateKey)
	if err != nil && err != database.ErrNotFound {
		return err
	}
	s.persistedStakerAccumulatedMintRate.Set(stakerMintRate)
	s.SetStakerAccumulatedMintRate(stakerMintRate)

	currentSupply, err := database.GetUInt64(s.singletonDB, currentSupplyKey)
	if err != nil {
		return err
	}
	s.persistedCurrentSupply = currentSupply
	s.SetCurrentSupply(constants.PrimaryNetworkID, currentSupply)

	lastAccepted, err := database.GetID(s.singletonDB, lastAcceptedKey)
	if err != nil {
		return err
	}
	s.persistedLastAccepted = lastAccepted
	s.lastAccepted = lastAccepted

	// Lookup the most recently indexed range on disk. If we haven't started
	// indexing the weights, then we keep the indexed heights as nil.
	indexedHeightsBytes, err := s.singletonDB.Get(heightsIndexedKey)
	if err == database.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	indexedHeights := &heightRange{}
	_, err = blocks.GenesisCodec.Unmarshal(indexedHeightsBytes, indexedHeights)
	if err != nil {
		return err
	}

	// If the indexed range is not up to date, then we will act as if the range
	// doesn't exist.
	lastAcceptedBlock, err := s.GetStatelessBlock(lastAccepted)
	if err != nil {
		return err
	}
	if indexedHeights.UpperBound != lastAcceptedBlock.Height() {
		return nil
	}
	s.indexedHeights = indexedHeights

	feePerStakerBytes, err := s.singletonDB.Get(feePerWeightStoredKey)
	if err != nil && err != database.ErrNotFound {
		return err
	}
	s.feePerWeightStored = new(big.Int).SetBytes(feePerStakerBytes)
	s.persistedFeePerWeightStored = new(big.Int).Set(s.feePerWeightStored)

	lastAccumulatedFee, err := database.GetUInt64(s.singletonDB, lastAccumulatedFeeKey)
	if err != nil && err != database.ErrNotFound {
		return err
	}
	s.lastAccumulatedFee = lastAccumulatedFee
	s.persistedLastAccumulatedFee = lastAccumulatedFee

	currentAccumulatedFee, err := database.GetUInt64(s.singletonDB, currentAccumulatedFeeKey)
	if err != nil && err != database.ErrNotFound {
		return err
	}
	s.currentAccumulatedFee = currentAccumulatedFee
	s.persistedCurrentAccumulatedFee = currentAccumulatedFee

	return nil
}

func (s *state) loadCurrentValidators() error {
	s.currentStakers = newBaseStakers()

	validatorIt := s.currentValidatorList.NewIterator()
	defer validatorIt.Release()
	for validatorIt.Next() {
		txIDBytes := validatorIt.Key()
		txID, err := ids.ToID(txIDBytes)
		if err != nil {
			return err
		}
		tx, _, err := s.GetTx(txID)
		if err != nil {
			return err
		}

		metadataBytes := validatorIt.Value()
		metadata := &validatorMetadata{
			txID: txID,
			// Note: we don't provide [LastUpdated] here because we expect it to
			// always be present on disk.
		}
		if err := parseValidatorMetadata(metadataBytes, metadata); err != nil {
			return err
		}

		stakerTx, ok := tx.Unsigned.(txs.Staker)
		if !ok {
			return fmt.Errorf("expected tx type txs.Staker but got %T", tx.Unsigned)
		}

		mintRate := new(big.Int).SetBytes(metadata.MintRateBytes)
		feePerWeightPaid := new(big.Int).SetBytes(metadata.FeePerWeightPaidBytes)
		staker, err := NewCurrentStakerWithRewardRate(txID, stakerTx, metadata.PotentialReward, mintRate, feePerWeightPaid)
		if err != nil {
			return err
		}

		validator := s.currentStakers.getOrCreateValidator(staker.SubnetID, staker.NodeID)
		validator.validator = staker

		s.currentStakers.stakers.ReplaceOrInsert(staker)

		s.validatorState.LoadValidatorMetadata(staker.NodeID, staker.SubnetID, metadata)
	}

	subnetValidatorIt := s.currentSubnetValidatorList.NewIterator()
	defer subnetValidatorIt.Release()
	for subnetValidatorIt.Next() {
		txIDBytes := subnetValidatorIt.Key()
		txID, err := ids.ToID(txIDBytes)
		if err != nil {
			return err
		}
		tx, _, err := s.GetTx(txID)
		if err != nil {
			return err
		}

		stakerTx, ok := tx.Unsigned.(txs.Staker)
		if !ok {
			return fmt.Errorf("expected tx type txs.Staker but got %T", tx.Unsigned)
		}

		metadataBytes := subnetValidatorIt.Value()
		metadata := &validatorMetadata{
			txID: txID,
			// use the start time as the fallback value
			// in case it's not stored in the database
			LastUpdated: uint64(stakerTx.StartTime().Unix()),
		}
		if err := parseValidatorMetadata(metadataBytes, metadata); err != nil {
			return err
		}

		mintRate := new(big.Int).SetBytes(metadata.MintRateBytes)
		feePerWeightPaid := new(big.Int).SetBytes(metadata.FeePerWeightPaidBytes)
		staker, err := NewCurrentStakerWithRewardRate(txID, stakerTx, metadata.PotentialReward, mintRate, feePerWeightPaid)
		if err != nil {
			return err
		}
		validator := s.currentStakers.getOrCreateValidator(staker.SubnetID, staker.NodeID)
		validator.validator = staker

		s.currentStakers.stakers.ReplaceOrInsert(staker)

		s.validatorState.LoadValidatorMetadata(staker.NodeID, staker.SubnetID, metadata)
	}

	delegatorIt := s.currentDelegatorList.NewIterator()
	defer delegatorIt.Release()

	subnetDelegatorIt := s.currentSubnetDelegatorList.NewIterator()
	defer subnetDelegatorIt.Release()

	for _, delegatorIt := range []database.Iterator{delegatorIt, subnetDelegatorIt} {
		for delegatorIt.Next() {
			txIDBytes := delegatorIt.Key()
			txID, err := ids.ToID(txIDBytes)
			if err != nil {
				return err
			}
			tx, _, err := s.GetTx(txID)
			if err != nil {
				return err
			}

			metadataBytes := delegatorIt.Value()
			metadata := delegatorMetadata{}
			if err := parseDelegatorMetadata(metadataBytes, &metadata); err != nil {
				return err
			}

			stakerTx, ok := tx.Unsigned.(txs.Staker)
			if !ok {
				return fmt.Errorf("expected tx type txs.Staker but got %T", tx.Unsigned)
			}

			mintRate := new(big.Int).SetBytes(metadata.MintRateBytes)
			feePerWeightPaid := new(big.Int).SetBytes(metadata.FeePerWeightPaidBytes)
			staker, err := NewCurrentStakerWithRewardRate(txID, stakerTx, metadata.PotentialReward, mintRate, feePerWeightPaid)
			if err != nil {
				return err
			}

			validator := s.currentStakers.getOrCreateValidator(staker.SubnetID, staker.NodeID)
			if validator.delegators == nil {
				validator.delegators = btree.NewG(defaultTreeDegree, (*Staker).Less)
			}
			validator.delegators.ReplaceOrInsert(staker)

			s.currentStakers.stakers.ReplaceOrInsert(staker)
		}
	}

	errs := wrappers.Errs{}
	errs.Add(
		validatorIt.Error(),
		subnetValidatorIt.Error(),
		delegatorIt.Error(),
		subnetDelegatorIt.Error(),
	)
	return errs.Err
}

func (s *state) loadPendingValidators() error {
	s.pendingStakers = newBaseStakers()

	validatorIt := s.pendingValidatorList.NewIterator()
	defer validatorIt.Release()

	subnetValidatorIt := s.pendingSubnetValidatorList.NewIterator()
	defer subnetValidatorIt.Release()

	for _, validatorIt := range []database.Iterator{validatorIt, subnetValidatorIt} {
		for validatorIt.Next() {
			txIDBytes := validatorIt.Key()
			txID, err := ids.ToID(txIDBytes)
			if err != nil {
				return err
			}
			tx, _, err := s.GetTx(txID)
			if err != nil {
				return err
			}

			stakerTx, ok := tx.Unsigned.(txs.Staker)
			if !ok {
				return fmt.Errorf("expected tx type txs.Staker but got %T", tx.Unsigned)
			}

			staker, err := NewPendingStaker(txID, stakerTx)
			if err != nil {
				return err
			}

			validator := s.pendingStakers.getOrCreateValidator(staker.SubnetID, staker.NodeID)
			validator.validator = staker

			s.pendingStakers.stakers.ReplaceOrInsert(staker)
		}
	}

	delegatorIt := s.pendingDelegatorList.NewIterator()
	defer delegatorIt.Release()

	subnetDelegatorIt := s.pendingSubnetDelegatorList.NewIterator()
	defer subnetDelegatorIt.Release()

	for _, delegatorIt := range []database.Iterator{delegatorIt, subnetDelegatorIt} {
		for delegatorIt.Next() {
			txIDBytes := delegatorIt.Key()
			txID, err := ids.ToID(txIDBytes)
			if err != nil {
				return err
			}
			tx, _, err := s.GetTx(txID)
			if err != nil {
				return err
			}

			stakerTx, ok := tx.Unsigned.(txs.Staker)
			if !ok {
				return fmt.Errorf("expected tx type txs.Staker but got %T", tx.Unsigned)
			}

			staker, err := NewPendingStaker(txID, stakerTx)
			if err != nil {
				return err
			}

			validator := s.pendingStakers.getOrCreateValidator(staker.SubnetID, staker.NodeID)
			if validator.delegators == nil {
				validator.delegators = btree.NewG(defaultTreeDegree, (*Staker).Less)
			}
			validator.delegators.ReplaceOrInsert(staker)

			s.pendingStakers.stakers.ReplaceOrInsert(staker)
		}
	}

	errs := wrappers.Errs{}
	errs.Add(
		validatorIt.Error(),
		subnetValidatorIt.Error(),
		delegatorIt.Error(),
		subnetDelegatorIt.Error(),
	)
	return errs.Err
}

// Invariant: initValidatorSets requires loadCurrentValidators to have already
// been called.
func (s *state) initValidatorSets() error {
	primaryValidators, ok := s.cfg.Validators.Get(constants.PrimaryNetworkID)
	if !ok {
		return errMissingValidatorSet
	}
	if primaryValidators.Len() != 0 {
		// Enforce the invariant that the validator set is empty here.
		return errValidatorSetAlreadyPopulated
	}
	err := s.ValidatorSet(constants.PrimaryNetworkID, primaryValidators)
	if err != nil {
		return err
	}

	vl := validators.NewLogger(s.ctx.Log, s.bootstrapped, constants.PrimaryNetworkID, s.ctx.NodeID)
	primaryValidators.RegisterCallbackListener(vl)

	s.metrics.SetLocalStake(primaryValidators.GetWeight(s.ctx.NodeID))
	s.metrics.SetTotalStake(primaryValidators.Weight())

	for subnetID := range s.cfg.TrackedSubnets {
		subnetValidators := validators.NewSet()
		err := s.ValidatorSet(subnetID, subnetValidators)
		if err != nil {
			return err
		}

		if !s.cfg.Validators.Add(subnetID, subnetValidators) {
			return fmt.Errorf("%w: %s", errDuplicateValidatorSet, subnetID)
		}

		vl := validators.NewLogger(s.ctx.Log, s.bootstrapped, subnetID, s.ctx.NodeID)
		subnetValidators.RegisterCallbackListener(vl)
	}
	return nil
}

func (s *state) write(updateValidators bool, height uint64) error {
	errs := wrappers.Errs{}
	errs.Add(
		s.writeBlocks(),
		s.writeCurrentStakers(updateValidators, height),
		s.writePendingStakers(),
		s.WriteValidatorMetadata(s.currentValidatorList, s.currentSubnetValidatorList), // Must be called after writeCurrentStakers
		s.writeTXs(),
		s.writeRewardUTXOs(),
		s.writeUTXOs(),
		s.writeSubnets(),
		s.writeTransformedSubnets(),
		s.writeSubnetSupplies(),
		s.writeChains(),
		s.writeMetadata(),
	)
	return errs.Err
}

func (s *state) Close() error {
	errs := wrappers.Errs{}
	errs.Add(
		s.pendingSubnetValidatorBaseDB.Close(),
		s.pendingSubnetDelegatorBaseDB.Close(),
		s.pendingDelegatorBaseDB.Close(),
		s.pendingValidatorBaseDB.Close(),
		s.pendingValidatorsDB.Close(),
		s.currentSubnetValidatorBaseDB.Close(),
		s.currentSubnetDelegatorBaseDB.Close(),
		s.currentDelegatorBaseDB.Close(),
		s.currentValidatorBaseDB.Close(),
		s.currentValidatorsDB.Close(),
		s.validatorsDB.Close(),
		s.txDB.Close(),
		s.rewardUTXODB.Close(),
		s.utxoDB.Close(),
		s.subnetBaseDB.Close(),
		s.transformedSubnetDB.Close(),
		s.supplyDB.Close(),
		s.chainDB.Close(),
		s.singletonDB.Close(),
		s.blockDB.Close(),
		s.blockIDDB.Close(),
	)
	return errs.Err
}

func (s *state) sync(genesis []byte) error {
	shouldInit, err := s.shouldInit()
	if err != nil {
		return fmt.Errorf(
			"failed to check if the database is initialized: %w",
			err,
		)
	}

	// If the database is empty, create the omega chain anew using the
	// provided genesis state
	if shouldInit {
		if err := s.init(genesis); err != nil {
			return fmt.Errorf(
				"failed to initialize the database: %w",
				err,
			)
		}
	}

	if err := s.load(); err != nil {
		return fmt.Errorf(
			"failed to load the database state: %w",
			err,
		)
	}
	return nil
}

func (s *state) init(genesisBytes []byte) error {
	// Create the genesis block and save it as being accepted (We don't do
	// genesisBlock.Accept() because then it'd look for genesisBlock's
	// non-existent parent)
	genesisID := hashing.ComputeHash256Array(genesisBytes)
	genesisBlock, err := blocks.NewApricotCommitBlock(genesisID, 0 /*height*/)
	if err != nil {
		return err
	}

	genesisState, err := genesis.ParseState(genesisBytes)
	if err != nil {
		return err
	}
	if err := s.syncGenesis(genesisBlock, genesisState); err != nil {
		return err
	}

	if err := s.doneInit(); err != nil {
		return err
	}

	return s.Commit()
}

func (s *state) AddStatelessBlock(block blocks.Block) {
	blkID := block.ID()
	s.addedBlockIDs[block.Height()] = blkID
	s.addedBlocks[blkID] = block
}

func (s *state) SetHeight(height uint64) {
	if s.indexedHeights == nil {
		// If indexedHeights hasn't been created yet, then we are newly tracking
		// the range. This means we should initialize the LowerBound to the
		// current height.
		s.indexedHeights = &heightRange{
			LowerBound: height,
		}
	}

	s.indexedHeights.UpperBound = height
	s.currentHeight = height
}

func (s *state) GetFeePerWeightStored() (*big.Int, error) {
	return new(big.Int).Set(s.feePerWeightStored), nil
}

func (s *state) SetFeePerWeightStored(f *big.Int) {
	s.feePerWeightStored.Set(f)
}

func (s *state) GetCurrentAccumulatedFee() (uint64, error) {
	return s.currentAccumulatedFee, nil
}

func (s *state) AddCurrentAccumulatedFee(f uint64) {
	s.currentAccumulatedFee += f
}

func (s *state) GetLastAccumulatedFee() (uint64, error) {
	return s.lastAccumulatedFee, nil
}

func (s *state) SetLastAccumulatedFee(f uint64) {
	s.lastAccumulatedFee = f
}

func (s *state) Commit() error {
	defer s.Abort()
	batch, err := s.CommitBatch()
	if err != nil {
		return err
	}
	return batch.Write()
}

func (s *state) Abort() {
	s.baseDB.Abort()
}

func (s *state) Checksum() ids.ID {
	return s.utxoState.Checksum()
}

func (s *state) CommitBatch() (database.Batch, error) {
	// updateValidators is set to true here so that the validator manager is
	// kept up to date with the last accepted state.
	if err := s.write(true /*=updateValidators*/, s.currentHeight); err != nil {
		return nil, err
	}
	return s.baseDB.CommitBatch()
}

func (s *state) writeBlocks() error {
	for blkID, blk := range s.addedBlocks {
		blkID := blkID
		blkBytes := blk.Bytes()
		blkHeight := blk.Height()
		heightKey := database.PackUInt64(blkHeight)

		delete(s.addedBlockIDs, blkHeight)
		s.blockIDCache.Put(blkHeight, blkID)
		if err := database.PutID(s.blockIDDB, heightKey, blkID); err != nil {
			return fmt.Errorf("failed to add blockID: %w", err)
		}

		delete(s.addedBlocks, blkID)
		// Note: Evict is used rather than Put here because blk may end up
		// referencing additional data (because of shared byte slices) that
		// would not be properly accounted for in the cache sizing.
		s.blockCache.Evict(blkID)
		if err := s.blockDB.Put(blkID[:], blkBytes); err != nil {
			return fmt.Errorf("failed to write block %s: %w", blkID, err)
		}
	}
	return nil
}

func (s *state) GetStatelessBlock(blockID ids.ID) (blocks.Block, error) {
	if blk, exists := s.addedBlocks[blockID]; exists {
		return blk, nil
	}
	if blk, cached := s.blockCache.Get(blockID); cached {
		if blk == nil {
			return nil, database.ErrNotFound
		}

		return blk, nil
	}

	blkBytes, err := s.blockDB.Get(blockID[:])
	if err == database.ErrNotFound {
		s.blockCache.Put(blockID, nil)
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	blk, status, _, err := parseStoredBlock(blkBytes)
	if err != nil {
		return nil, err
	}

	if status != choices.Accepted {
		s.blockCache.Put(blockID, nil)
		return nil, database.ErrNotFound
	}

	s.blockCache.Put(blockID, blk)
	return blk, nil
}

func (s *state) GetBlockIDAtHeight(height uint64) (ids.ID, error) {
	if blkID, exists := s.addedBlockIDs[height]; exists {
		return blkID, nil
	}
	if blkID, cached := s.blockIDCache.Get(height); cached {
		if blkID == ids.Empty {
			return ids.Empty, database.ErrNotFound
		}

		return blkID, nil
	}

	heightKey := database.PackUInt64(height)

	blkID, err := database.GetID(s.blockIDDB, heightKey)
	if err == database.ErrNotFound {
		s.blockIDCache.Put(height, ids.Empty)
		return ids.Empty, database.ErrNotFound
	}
	if err != nil {
		return ids.Empty, err
	}

	s.blockIDCache.Put(height, blkID)
	return blkID, nil
}

func (s *state) writeCurrentStakers(updateValidators bool, height uint64) error {
	heightBytes := database.PackUInt64(height)
	rawNestedPublicKeyDiffDB := prefixdb.New(heightBytes, s.nestedValidatorPublicKeyDiffsDB)
	nestedPKDiffDB := linkeddb.NewDefault(rawNestedPublicKeyDiffDB)

	for subnetID, validatorDiffs := range s.currentStakers.validatorDiffs {
		delete(s.currentStakers.validatorDiffs, subnetID)

		// Select db to write to
		validatorDB := s.currentSubnetValidatorList
		delegatorDB := s.currentSubnetDelegatorList
		if subnetID == constants.PrimaryNetworkID {
			validatorDB = s.currentValidatorList
			delegatorDB = s.currentDelegatorList
		}

		prefixStruct := heightWithSubnet{
			Height:   height,
			SubnetID: subnetID,
		}
		prefixBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, prefixStruct)
		if err != nil {
			return fmt.Errorf("failed to create prefix bytes: %w", err)
		}
		rawNestedWeightDiffDB := prefixdb.New(prefixBytes, s.nestedValidatorWeightDiffsDB)
		nestedWeightDiffDB := linkeddb.NewDefault(rawNestedWeightDiffDB)

		// Record the change in weight and/or public key for each validator.
		for nodeID, validatorDiff := range validatorDiffs {
			// Copy [nodeID] so it doesn't get overwritten next iteration.
			nodeID := nodeID

			weightDiff := &ValidatorWeightDiff{
				Decrease: validatorDiff.validatorStatus == deleted,
			}
			switch validatorDiff.validatorStatus {
			case added:
				staker := validatorDiff.validator
				weightDiff.Amount = staker.Weight

				// Invariant: Only the Primary Network contains non-nil public
				// keys.
				if staker.PublicKey != nil {
					// Record that the public key for the validator is being
					// added. This means the prior value for the public key was
					// nil.
					err := s.flatValidatorPublicKeyDiffsDB.Put(
						marshalDiffKey(constants.PrimaryNetworkID, height, nodeID),
						nil,
					)
					if err != nil {
						return err
					}
				}

				mintRate := staker.MintRate
				if mintRate == nil {
					mintRate = new(big.Int)
				}

				feePerWeightPaid := staker.FeePerWeightPaid
				if feePerWeightPaid == nil {
					feePerWeightPaid = new(big.Int)
				}

				// The validator is being added.
				//
				// Invariant: It's impossible for a delegator to have been
				// rewarded in the same block that the validator was added.
				metadata := &validatorMetadata{
					txID:        staker.TxID,
					lastUpdated: staker.StartTime,

					UpDuration:               0,
					LastUpdated:              uint64(staker.StartTime.Unix()),
					PotentialReward:          staker.PotentialReward,
					PotentialDelegateeReward: 0,
					MintRateBytes:            mintRate.Bytes(),
					FeePerWeightPaidBytes:    feePerWeightPaid.Bytes(),
				}

				metadataBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, metadata)
				if err != nil {
					return fmt.Errorf("failed to serialize current validator: %w", err)
				}

				if err = validatorDB.Put(staker.TxID[:], metadataBytes); err != nil {
					return fmt.Errorf("failed to write current validator to list: %w", err)
				}

				s.validatorState.LoadValidatorMetadata(nodeID, subnetID, metadata)
			case deleted:
				staker := validatorDiff.validator
				weightDiff.Amount = staker.Weight

				// Invariant: Only the Primary Network contains non-nil public
				// keys.
				if staker.PublicKey != nil {
					// Record that the public key for the validator is being
					// removed. This means we must record the prior value of the
					// public key.
					//
					// Note: We store the uncompressed public key here as it is
					// significantly more efficient to parse when applying
					// diffs.
					err := s.flatValidatorPublicKeyDiffsDB.Put(
						marshalDiffKey(constants.PrimaryNetworkID, height, nodeID),
						staker.PublicKey.Serialize(),
					)
					if err != nil {
						return err
					}

					// TODO: Remove this once we no longer support version
					// rollbacks.
					//
					// Note: We store the compressed public key here.
					pkBytes := bls.PublicKeyToBytes(staker.PublicKey)
					if err := nestedPKDiffDB.Put(nodeID[:], pkBytes); err != nil {
						return err
					}
				}

				if err := validatorDB.Delete(staker.TxID[:]); err != nil {
					return fmt.Errorf("failed to delete current staker: %w", err)
				}

				s.validatorState.DeleteValidatorMetadata(nodeID, subnetID)
			}

			err := writeCurrentDelegatorDiff(
				delegatorDB,
				weightDiff,
				validatorDiff,
			)
			if err != nil {
				return err
			}

			if weightDiff.Amount == 0 {
				// No weight change to record; go to next validator.
				continue
			}

			err = s.flatValidatorWeightDiffsDB.Put(
				marshalDiffKey(subnetID, height, nodeID),
				marshalWeightDiff(weightDiff),
			)
			if err != nil {
				return err
			}

			// TODO: Remove this once we no longer support version rollbacks.
			weightDiffBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, weightDiff)
			if err != nil {
				return fmt.Errorf("failed to serialize validator weight diff: %w", err)
			}
			if err := nestedWeightDiffDB.Put(nodeID[:], weightDiffBytes); err != nil {
				return err
			}

			// TODO: Move the validator set management out of the state package
			if !updateValidators {
				continue
			}

			// We only track the current validator set of tracked subnets.
			if subnetID != constants.PrimaryNetworkID && !s.cfg.TrackedSubnets.Contains(subnetID) {
				continue
			}

			if weightDiff.Decrease {
				err = validators.RemoveWeight(s.cfg.Validators, subnetID, nodeID, weightDiff.Amount)
			} else {
				if validatorDiff.validatorStatus == added {
					staker := validatorDiff.validator
					err = validators.Add(
						s.cfg.Validators,
						subnetID,
						nodeID,
						staker.PublicKey,
						staker.TxID,
						weightDiff.Amount,
					)
				} else {
					err = validators.AddWeight(s.cfg.Validators, subnetID, nodeID, weightDiff.Amount)
				}
			}
			if err != nil {
				return fmt.Errorf("failed to update validator weight: %w", err)
			}
		}
	}

	// TODO: Move validator set management out of the state package
	//
	// Attempt to update the stake metrics
	if !updateValidators {
		return nil
	}
	primaryValidators, ok := s.cfg.Validators.Get(constants.PrimaryNetworkID)
	if !ok {
		return nil
	}
	s.metrics.SetLocalStake(primaryValidators.GetWeight(s.ctx.NodeID))
	s.metrics.SetTotalStake(primaryValidators.Weight())
	return nil
}

func writeCurrentDelegatorDiff(
	currentDelegatorList linkeddb.LinkedDB,
	weightDiff *ValidatorWeightDiff,
	validatorDiff *diffValidator,
) error {
	addedDelegatorIterator := NewTreeIterator(validatorDiff.addedDelegators)
	defer addedDelegatorIterator.Release()
	for addedDelegatorIterator.Next() {
		staker := addedDelegatorIterator.Value()

		if err := weightDiff.Add(false, staker.Weight); err != nil {
			return fmt.Errorf("failed to increase node weight diff: %w", err)
		}

		metadata := delegatorMetadata{
			PotentialReward:       staker.PotentialReward,
			MintRateBytes:         staker.MintRate.Bytes(),
			FeePerWeightPaidBytes: staker.FeePerWeightPaid.Bytes(),
		}

		metadataBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, &metadata)
		if err != nil {
			return fmt.Errorf("failed to serialize current delegator: %w", err)
		}

		if err := currentDelegatorList.Put(staker.TxID[:], metadataBytes); err != nil {
			return fmt.Errorf("failed to write current delegator to list: %w", err)
		}
	}

	for _, staker := range validatorDiff.deletedDelegators {
		if err := weightDiff.Add(true, staker.Weight); err != nil {
			return fmt.Errorf("failed to decrease node weight diff: %w", err)
		}

		if err := currentDelegatorList.Delete(staker.TxID[:]); err != nil {
			return fmt.Errorf("failed to delete current staker: %w", err)
		}
	}
	return nil
}

func (s *state) writePendingStakers() error {
	for subnetID, subnetValidatorDiffs := range s.pendingStakers.validatorDiffs {
		delete(s.pendingStakers.validatorDiffs, subnetID)

		validatorDB := s.pendingSubnetValidatorList
		delegatorDB := s.pendingSubnetDelegatorList
		if subnetID == constants.PrimaryNetworkID {
			validatorDB = s.pendingValidatorList
			delegatorDB = s.pendingDelegatorList
		}

		for _, validatorDiff := range subnetValidatorDiffs {
			err := writePendingDiff(
				validatorDB,
				delegatorDB,
				validatorDiff,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writePendingDiff(
	pendingValidatorList linkeddb.LinkedDB,
	pendingDelegatorList linkeddb.LinkedDB,
	validatorDiff *diffValidator,
) error {
	switch validatorDiff.validatorStatus {
	case added:
		err := pendingValidatorList.Put(validatorDiff.validator.TxID[:], nil)
		if err != nil {
			return fmt.Errorf("failed to add pending validator: %w", err)
		}
	case deleted:
		err := pendingValidatorList.Delete(validatorDiff.validator.TxID[:])
		if err != nil {
			return fmt.Errorf("failed to delete pending validator: %w", err)
		}
	}

	addedDelegatorIterator := NewTreeIterator(validatorDiff.addedDelegators)
	defer addedDelegatorIterator.Release()
	for addedDelegatorIterator.Next() {
		staker := addedDelegatorIterator.Value()

		if err := pendingDelegatorList.Put(staker.TxID[:], nil); err != nil {
			return fmt.Errorf("failed to write pending delegator to list: %w", err)
		}
	}

	for _, staker := range validatorDiff.deletedDelegators {
		if err := pendingDelegatorList.Delete(staker.TxID[:]); err != nil {
			return fmt.Errorf("failed to delete pending delegator: %w", err)
		}
	}
	return nil
}

func (s *state) writeTXs() error {
	for txID, txStatus := range s.addedTxs {
		txID := txID

		stx := txBytesAndStatus{
			Tx:     txStatus.tx.Bytes(),
			Status: txStatus.status,
		}

		// Note that we're serializing a [txBytesAndStatus] here, not a
		// *txs.Tx, so we don't use [txs.Codec].
		txBytes, err := txs.GenesisCodec.Marshal(txs.Version, &stx)
		if err != nil {
			return fmt.Errorf("failed to serialize tx: %w", err)
		}

		delete(s.addedTxs, txID)
		// Note: Evict is used rather than Put here because stx may end up
		// referencing additional data (because of shared byte slices) that
		// would not be properly accounted for in the cache sizing.
		s.txCache.Evict(txID)
		if err := s.txDB.Put(txID[:], txBytes); err != nil {
			return fmt.Errorf("failed to add tx: %w", err)
		}
	}
	return nil
}

func (s *state) writeRewardUTXOs() error {
	for txID, utxos := range s.addedRewardUTXOs {
		delete(s.addedRewardUTXOs, txID)
		s.rewardUTXOsCache.Put(txID, utxos)
		rawTxDB := prefixdb.New(txID[:], s.rewardUTXODB)
		txDB := linkeddb.NewDefault(rawTxDB)

		for _, utxo := range utxos {
			utxoBytes, err := txs.GenesisCodec.Marshal(txs.Version, utxo)
			if err != nil {
				return fmt.Errorf("failed to serialize reward UTXO: %w", err)
			}
			utxoID := utxo.InputID()
			if err := txDB.Put(utxoID[:], utxoBytes); err != nil {
				return fmt.Errorf("failed to add reward UTXO: %w", err)
			}
		}
	}
	return nil
}

func (s *state) writeUTXOs() error {
	for utxoID, utxo := range s.modifiedUTXOs {
		delete(s.modifiedUTXOs, utxoID)

		if utxo == nil {
			if err := s.utxoState.DeleteUTXO(utxoID); err != nil {
				return fmt.Errorf("failed to delete UTXO: %w", err)
			}
			continue
		}
		if err := s.utxoState.PutUTXO(utxo); err != nil {
			return fmt.Errorf("failed to add UTXO: %w", err)
		}
	}
	return nil
}

func (s *state) writeSubnets() error {
	for _, subnet := range s.addedSubnets {
		subnetID := subnet.ID()

		if err := s.subnetDB.Put(subnetID[:], nil); err != nil {
			return fmt.Errorf("failed to write subnet: %w", err)
		}
	}
	s.addedSubnets = nil
	return nil
}

func (s *state) writeTransformedSubnets() error {
	for subnetID, tx := range s.transformedSubnets {
		txID := tx.ID()

		delete(s.transformedSubnets, subnetID)
		// Note: Evict is used rather than Put here because tx may end up
		// referencing additional data (because of shared byte slices) that
		// would not be properly accounted for in the cache sizing.
		s.transformedSubnetCache.Evict(subnetID)
		if err := database.PutID(s.transformedSubnetDB, subnetID[:], txID); err != nil {
			return fmt.Errorf("failed to write transformed subnet: %w", err)
		}
	}
	return nil
}

func (s *state) writeSubnetSupplies() error {
	for subnetID, supply := range s.modifiedSupplies {
		supply := supply
		delete(s.modifiedSupplies, subnetID)
		s.supplyCache.Put(subnetID, &supply)
		if err := database.PutUInt64(s.supplyDB, subnetID[:], supply); err != nil {
			return fmt.Errorf("failed to write subnet supply: %w", err)
		}
	}
	return nil
}

func (s *state) writeChains() error {
	for subnetID, chains := range s.addedChains {
		for _, chain := range chains {
			chainDB := s.getChainDB(subnetID)

			chainID := chain.ID()
			if err := chainDB.Put(chainID[:], nil); err != nil {
				return fmt.Errorf("failed to write chain: %w", err)
			}
		}
		delete(s.addedChains, subnetID)
	}
	return nil
}

func (s *state) writeMetadata() error {
	if !s.persistedTimestamp.Equal(s.timestamp) {
		if err := database.PutTimestamp(s.singletonDB, timestampKey, s.timestamp); err != nil {
			return fmt.Errorf("failed to write timestamp: %w", err)
		}
		s.persistedTimestamp = s.timestamp
	}
	if !s.persistedStakeSyncTimestamp.Equal(s.stakeSyncTimestamp) {
		if err := database.PutTimestamp(s.singletonDB, stakeSyncTimestampKey, s.stakeSyncTimestamp); err != nil {
			return fmt.Errorf("failed to write timestamp: %w", err)
		}
		s.persistedStakeSyncTimestamp = s.stakeSyncTimestamp
	}
	if s.persistedStakerAccumulatedMintRate.Cmp(s.stakerAccumulatedMintRate) != 0 {
		if err := database.PutBigInt(s.singletonDB, stakerMintRateKey, s.stakerAccumulatedMintRate); err != nil {
			return fmt.Errorf("failed to write staker mint rate: %w", err)
		}
		s.persistedStakerAccumulatedMintRate = s.stakerAccumulatedMintRate
	}
	if s.persistedCurrentSupply != s.currentSupply {
		if err := database.PutUInt64(s.singletonDB, currentSupplyKey, s.currentSupply); err != nil {
			return fmt.Errorf("failed to write current supply: %w", err)
		}
		s.persistedCurrentSupply = s.currentSupply
	}
	if s.persistedLastAccepted != s.lastAccepted {
		if err := database.PutID(s.singletonDB, lastAcceptedKey, s.lastAccepted); err != nil {
			return fmt.Errorf("failed to write last accepted: %w", err)
		}
		s.persistedLastAccepted = s.lastAccepted
	}

	if s.indexedHeights != nil {
		indexedHeightsBytes, err := blocks.GenesisCodec.Marshal(blocks.Version, s.indexedHeights)
		if err != nil {
			return err
		}
		if err := s.singletonDB.Put(heightsIndexedKey, indexedHeightsBytes); err != nil {
			return fmt.Errorf("failed to write indexed range: %w", err)
		}
	}

	if s.persistedFeePerWeightStored.Cmp(s.feePerWeightStored) != 0 {
		if err := s.singletonDB.Put(feePerWeightStoredKey, s.feePerWeightStored.Bytes()); err != nil {
			return fmt.Errorf("failed to write fee per token stored: %w", err)
		}
		s.persistedFeePerWeightStored.Set(s.feePerWeightStored)
	}

	if s.persistedLastAccumulatedFee != s.lastAccumulatedFee {
		if err := database.PutUInt64(s.singletonDB, lastAccumulatedFeeKey, s.lastAccumulatedFee); err != nil {
			return fmt.Errorf("failed to write last accumulated fee: %w", err)
		}
		s.persistedLastAccumulatedFee = s.lastAccumulatedFee
	}

	if s.persistedCurrentAccumulatedFee != s.currentAccumulatedFee {
		if err := database.PutUInt64(s.singletonDB, currentAccumulatedFeeKey, s.currentAccumulatedFee); err != nil {
			return fmt.Errorf("failed to write current accumulated fee: %w", err)
		}
		s.persistedCurrentAccumulatedFee = s.currentAccumulatedFee
	}

	return nil
}

// Returns the block, status of the block, and whether it is a [stateBlk].
// Invariant: blkBytes is safe to parse with blocks.GenesisCodec
//
// TODO: Remove after v1.11.x is activated
func parseStoredBlock(blkBytes []byte) (blocks.Block, choices.Status, bool, error) {
	// Attempt to parse as blocks.Block
	blk, err := blocks.Parse(blocks.GenesisCodec, blkBytes)
	if err == nil {
		return blk, choices.Accepted, false, nil
	}

	// Fallback to [stateBlk]
	blkState := stateBlk{}
	if _, err := blocks.GenesisCodec.Unmarshal(blkBytes, &blkState); err != nil {
		return nil, choices.Processing, false, err
	}

	blkState.Blk, err = blocks.Parse(blocks.GenesisCodec, blkState.Bytes)
	if err != nil {
		return nil, choices.Processing, false, err
	}

	return blkState.Blk, blkState.Status, true, nil
}

func (s *state) PruneAndIndex(lock sync.Locker, log logging.Logger) error {
	lock.Lock()
	// It is possible that new blocks are added after grabbing this iterator. New
	// blocks are guaranteed to be accepted and height-indexed, so we don't need to
	// check them.
	blockIterator := s.blockDB.NewIterator()
	// Releasing is done using a closure to ensure that updating blockIterator will
	// result in having the most recent iterator released when executing the
	// deferred function.
	defer func() {
		blockIterator.Release()
	}()

	// While we are pruning the disk, we disable caching of the data we are
	// modifying. Caching is re-enabled when pruning finishes.
	//
	// Note: If an unexpected error occurs the caches are never re-enabled.
	// That's fine as the node is going to be in an unhealthy state regardless.
	oldBlockIDCache := s.blockIDCache
	s.blockIDCache = &cache.Empty[uint64, ids.ID]{}
	lock.Unlock()

	log.Info("starting state pruning and indexing")

	var (
		startTime  = time.Now()
		lastCommit = startTime
		lastUpdate = startTime
		numPruned  = 0
		numIndexed = 0
	)

	for blockIterator.Next() {
		blkBytes := blockIterator.Value()

		blk, status, isStateBlk, err := parseStoredBlock(blkBytes)
		if err != nil {
			return err
		}

		if status != choices.Accepted {
			// Remove non-accepted blocks from disk.
			if err := s.blockDB.Delete(blockIterator.Key()); err != nil {
				return fmt.Errorf("failed to delete block: %w", err)
			}

			numPruned++

			// We don't index the height of non-accepted blocks.
			continue
		}

		blkHeight := blk.Height()
		blkID := blk.ID()

		// Populate the map of height -> blockID.
		heightKey := database.PackUInt64(blkHeight)
		if err := database.PutID(s.blockIDDB, heightKey, blkID); err != nil {
			return fmt.Errorf("failed to add blockID: %w", err)
		}

		// Since we only store accepted blocks on disk, we only need to store a map of
		// ids.ID to Block.
		if isStateBlk {
			if err := s.blockDB.Put(blkID[:], blkBytes); err != nil {
				return fmt.Errorf("failed to write block: %w", err)
			}
		}

		numIndexed++

		if numIndexed%pruneCommitLimit == 0 {
			// We must hold the lock during committing to make sure we don't
			// attempt to commit to disk while a block is concurrently being
			// accepted.
			lock.Lock()
			errs := wrappers.Errs{}
			errs.Add(
				s.Commit(),
				blockIterator.Error(),
			)
			lock.Unlock()
			if errs.Errored() {
				return errs.Err
			}

			// We release the iterator here to allow the underlying database to
			// clean up deleted state.
			blockIterator.Release()

			now := time.Now()
			if now.Sub(lastUpdate) > pruneUpdateFrequency {
				lastUpdate = now

				progress := timer.ProgressFromHash(blkID[:])
				eta := timer.EstimateETA(
					startTime,
					progress,
					stdmath.MaxUint64,
				)

				log.Info("committing state pruning and indexing",
					zap.Int("numPruned", numPruned),
					zap.Int("numIndexed", numIndexed),
					zap.Duration("eta", eta),
				)
			}

			// We take the minimum here because it's possible that the node is
			// currently bootstrapping. This would mean that grabbing the lock
			// could take an extremely long period of time; which we should not
			// delay processing for.
			pruneDuration := now.Sub(lastCommit)
			sleepDuration := math.Min(
				pruneCommitSleepMultiplier*pruneDuration,
				pruneCommitSleepCap,
			)
			time.Sleep(sleepDuration)

			// Make sure not to include the sleep duration into the next prune
			// duration.
			lastCommit = time.Now()

			blockIterator = s.blockDB.NewIteratorWithStart(blkID[:])
		}
	}

	// Ensure we fully iterated over all blocks before writing that pruning has
	// finished.
	//
	// Note: This is needed because a transient read error could cause the
	// iterator to stop early.
	if err := blockIterator.Error(); err != nil {
		return err
	}

	if err := s.donePrune(); err != nil {
		return err
	}

	// We must hold the lock during committing to make sure we don't
	// attempt to commit to disk while a block is concurrently being
	// accepted.
	lock.Lock()
	defer lock.Unlock()

	// Make sure we flush the original cache before re-enabling it to prevent
	// surfacing any stale data.
	oldBlockIDCache.Flush()
	s.blockIDCache = oldBlockIDCache

	log.Info("finished state pruning and indexing",
		zap.Int("numPruned", numPruned),
		zap.Int("numIndexed", numIndexed),
		zap.Duration("duration", time.Since(startTime)),
	)

	return s.Commit()
}