// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package builder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/chains"
	"github.com/DioneProtocol/odysseygo/chains/atomic"
	"github.com/DioneProtocol/odysseygo/codec"
	"github.com/DioneProtocol/odysseygo/codec/linearcodec"
	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/manager"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/database/versiondb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/engine/common"
	"github.com/DioneProtocol/odysseygo/snow/uptime"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/formatting"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/utils/json"
	"github.com/DioneProtocol/odysseygo/utils/logging"
	"github.com/DioneProtocol/odysseygo/utils/timer/mockable"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/utils/wrappers"
	"github.com/DioneProtocol/odysseygo/version"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/feecollector"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/api"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/config"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/fx"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/metrics"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/status"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs/mempool"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/utxo"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"

	blockexecutor "github.com/DioneProtocol/odysseygo/vms/omegavm/blocks/executor"
	txbuilder "github.com/DioneProtocol/odysseygo/vms/omegavm/txs/builder"
	txexecutor "github.com/DioneProtocol/odysseygo/vms/omegavm/txs/executor"
	ovalidators "github.com/DioneProtocol/odysseygo/vms/omegavm/validators"
)

const (
	defaultWeight = 10000
	trackChecksum = false
)

var (
	defaultMinValidatorStakingDuration = 24 * time.Hour
	defaultMaxValidatorStakingDuration = 365 * 24 * time.Hour
	defaultMinDelegatorStakingDuration = 24 * time.Hour
	defaultMaxDelegatorStakingDuration = 365 * 24 * time.Hour
	defaultGenesisTime                 = time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultValidateStartTime           = defaultGenesisTime
	defaultValidateEndTime             = defaultValidateStartTime.Add(10 * defaultMinValidatorStakingDuration)
	defaultMinValidatorStake           = 5 * units.MilliDione
	defaultBalance                     = 100 * defaultMinValidatorStake
	preFundedKeys                      = secp256k1.TestKeys()
	dioneAssetID                       = ids.ID{'y', 'e', 'e', 't'}
	defaultTxFee                       = uint64(100)
	aChainID                           = ids.Empty.Prefix(0)
	dChainID                           = ids.Empty.Prefix(1)

	testSubnet1            *txs.Tx
	testSubnet1ControlKeys = preFundedKeys[0:3]

	errMissing = errors.New("missing")
)

type mutableSharedMemory struct {
	atomic.SharedMemory
}

type environment struct {
	Builder
	blkManager blockexecutor.Manager
	mempool    mempool.Mempool
	sender     *common.SenderTest

	isBootstrapped *utils.Atomic[bool]
	config         *config.Config
	clk            *mockable.Clock
	baseDB         *versiondb.Database
	ctx            *snow.Context
	msm            *mutableSharedMemory
	fx             fx.Fx
	state          state.State
	atomicUTXOs    dione.AtomicUTXOManager
	uptimes        uptime.Manager
	utxosHandler   utxo.Handler
	txBuilder      txbuilder.Builder
	backend        txexecutor.Backend
}

func newEnvironment(t *testing.T) *environment {
	require := require.New(t)

	res := &environment{
		isBootstrapped: &utils.Atomic[bool]{},
		config:         defaultConfig(),
		clk:            defaultClock(),
	}
	res.isBootstrapped.Set(true)

	baseDBManager := manager.NewMemDB(version.Semantic1_0_0)
	res.baseDB = versiondb.New(baseDBManager.Current().Database)
	res.ctx, res.msm = defaultCtx(res.baseDB)

	res.ctx.Lock.Lock()
	defer res.ctx.Lock.Unlock()

	res.fx = defaultFx(t, res.clk, res.ctx.Log, res.isBootstrapped.Get())

	rewardsCalc := reward.NewCalculator(res.config.RewardConfig)
	res.state = defaultState(t, res.config, res.ctx, res.baseDB, rewardsCalc)

	res.atomicUTXOs = dione.NewAtomicUTXOManager(res.ctx.SharedMemory, txs.Codec)
	res.uptimes = uptime.NewManager(res.state)
	res.utxosHandler = utxo.NewHandler(res.ctx, res.clk, res.fx)

	res.txBuilder = txbuilder.New(
		res.ctx,
		res.config,
		res.clk,
		res.fx,
		res.state,
		res.atomicUTXOs,
		res.utxosHandler,
	)

	genesisID := res.state.GetLastAccepted()
	res.backend = txexecutor.Backend{
		Config:       res.config,
		Ctx:          res.ctx,
		Clk:          res.clk,
		Bootstrapped: res.isBootstrapped,
		Fx:           res.fx,
		FlowChecker:  res.utxosHandler,
		Uptimes:      res.uptimes,
		Rewards:      rewardsCalc,
	}

	registerer := prometheus.NewRegistry()
	res.sender = &common.SenderTest{T: t}

	metrics, err := metrics.New("", registerer)
	require.NoError(err)

	res.mempool, err = mempool.NewMempool("mempool", registerer, res)
	require.NoError(err)

	res.blkManager = blockexecutor.NewManager(
		res.mempool,
		metrics,
		res.state,
		&res.backend,
		ovalidators.TestManager,
	)

	res.Builder = New(
		res.mempool,
		res.txBuilder,
		&res.backend,
		res.blkManager,
		nil, // toEngine,
		res.sender,
	)

	res.Builder.SetPreference(genesisID)
	addSubnet(t, res)

	return res
}

func addSubnet(t *testing.T, env *environment) {
	require := require.New(t)

	// Create a subnet
	var err error
	testSubnet1, err = env.txBuilder.NewCreateSubnetTx(
		2, // threshold; 2 sigs from keys[0], keys[1], keys[2] needed to add validator to this subnet
		[]ids.ShortID{ // control keys
			preFundedKeys[0].PublicKey().Address(),
			preFundedKeys[1].PublicKey().Address(),
			preFundedKeys[2].PublicKey().Address(),
		},
		[]*secp256k1.PrivateKey{preFundedKeys[0]},
		preFundedKeys[0].PublicKey().Address(),
	)
	require.NoError(err)

	// store it
	genesisID := env.state.GetLastAccepted()
	stateDiff, err := state.NewDiff(genesisID, env.blkManager)
	require.NoError(err)

	executor := txexecutor.StandardTxExecutor{
		Backend: &env.backend,
		State:   stateDiff,
		Tx:      testSubnet1,
	}
	require.NoError(testSubnet1.Unsigned.Visit(&executor))

	stateDiff.AddTx(testSubnet1, status.Committed)
	require.NoError(stateDiff.Apply(env.state))
}

func defaultState(
	t *testing.T,
	cfg *config.Config,
	ctx *snow.Context,
	db database.Database,
	rewards reward.Calculator,
) state.State {
	require := require.New(t)

	execCfg, _ := config.GetExecutionConfig([]byte(`{}`))
	genesisBytes := buildGenesisTest(t, ctx)
	state, err := state.New(
		db,
		genesisBytes,
		prometheus.NewRegistry(),
		cfg,
		execCfg,
		ctx,
		metrics.Noop,
		rewards,
		&utils.Atomic[bool]{},
	)
	require.NoError(err)

	// persist and reload to init a bunch of in-memory stuff
	state.SetHeight(0)
	require.NoError(state.Commit())
	return state
}

func defaultCtx(db database.Database) (*snow.Context, *mutableSharedMemory) {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = 10
	ctx.AChainID = aChainID
	ctx.DChainID = dChainID
	ctx.DIONEAssetID = dioneAssetID

	atomicDB := prefixdb.New([]byte{1}, db)
	m := atomic.NewMemory(atomicDB)

	msm := &mutableSharedMemory{
		SharedMemory: m.NewSharedMemory(ctx.ChainID),
	}
	ctx.SharedMemory = msm

	feeDb := prefixdb.New([]byte{2}, db)
	f, err := feecollector.New(feeDb)
	if err != nil {
		panic(err)
	}
	ctx.FeeCollector = f

	ctx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				constants.OmegaChainID: constants.PrimaryNetworkID,
				aChainID:               constants.PrimaryNetworkID,
				dChainID:               constants.PrimaryNetworkID,
			}[chainID]
			if !ok {
				return ids.Empty, errMissing
			}
			return subnetID, nil
		},
	}

	return ctx, msm
}

func defaultConfig() *config.Config {
	vdrs := validators.NewManager()
	primaryVdrs := validators.NewSet()
	_ = vdrs.Add(constants.PrimaryNetworkID, primaryVdrs)
	return &config.Config{
		Chains:                    chains.TestManager,
		UptimeLockedCalculator:    uptime.NewLockedCalculator(),
		Validators:                vdrs,
		TxFee:                     defaultTxFee,
		CreateSubnetTxFee:         100 * defaultTxFee,
		CreateBlockchainTxFee:     100 * defaultTxFee,
		MinValidatorStake:         5 * units.MilliDione,
		MaxValidatorStake:         500 * units.MilliDione,
		MinDelegatorStake:         1 * units.MilliDione,
		MinValidatorStakeDuration: defaultMinValidatorStakingDuration,
		MaxValidatorStakeDuration: defaultMaxValidatorStakingDuration,
		MinDelegatorStakeDuration: defaultMinDelegatorStakingDuration,
		MaxDelegatorStakeDuration: defaultMaxDelegatorStakingDuration,
		RewardConfig: reward.Config{
			MaxConsumptionRate: .12 * reward.PercentDenominator,
			MinConsumptionRate: .10 * reward.PercentDenominator,
			MintingPeriod:      365 * 24 * time.Hour,
			SupplyCap:          720 * units.MegaDione,
		},
		ApricotPhase3Time: defaultValidateEndTime,
		ApricotPhase5Time: defaultValidateEndTime,
		BanffTime:         time.Time{}, // neglecting fork ordering this for package tests
	}
}

func defaultClock() *mockable.Clock {
	// set time after Banff fork (and before default nextStakerTime)
	clk := &mockable.Clock{}
	clk.Set(defaultGenesisTime)
	return clk
}

type fxVMInt struct {
	registry codec.Registry
	clk      *mockable.Clock
	log      logging.Logger
}

func (fvi *fxVMInt) CodecRegistry() codec.Registry {
	return fvi.registry
}

func (fvi *fxVMInt) Clock() *mockable.Clock {
	return fvi.clk
}

func (fvi *fxVMInt) Logger() logging.Logger {
	return fvi.log
}

func defaultFx(t *testing.T, clk *mockable.Clock, log logging.Logger, isBootstrapped bool) fx.Fx {
	require := require.New(t)

	fxVMInt := &fxVMInt{
		registry: linearcodec.NewDefault(),
		clk:      clk,
		log:      log,
	}
	res := &secp256k1fx.Fx{}
	require.NoError(res.Initialize(fxVMInt))
	if isBootstrapped {
		require.NoError(res.Bootstrapped())
	}
	return res
}

func buildGenesisTest(t *testing.T, ctx *snow.Context) []byte {
	require := require.New(t)

	genesisUTXOs := make([]api.UTXO, len(preFundedKeys))
	for i, key := range preFundedKeys {
		id := key.PublicKey().Address()
		addr, err := address.FormatBech32(constants.UnitTestHRP, id.Bytes())
		require.NoError(err)
		genesisUTXOs[i] = api.UTXO{
			Amount:  json.Uint64(defaultBalance),
			Address: addr,
		}
	}

	genesisValidators := make([]api.PermissionlessValidator, len(preFundedKeys))
	for i, key := range preFundedKeys {
		nodeID := ids.NodeID(key.PublicKey().Address())
		addr, err := address.FormatBech32(constants.UnitTestHRP, nodeID.Bytes())
		require.NoError(err)
		genesisValidators[i] = api.PermissionlessValidator{
			Staker: api.Staker{
				StartTime: json.Uint64(defaultValidateStartTime.Unix()),
				EndTime:   json.Uint64(defaultValidateEndTime.Unix()),
				NodeID:    nodeID,
			},
			RewardOwner: &api.Owner{
				Threshold: 1,
				Addresses: []string{addr},
			},
			Staked: []api.UTXO{{
				Amount:  json.Uint64(defaultWeight),
				Address: addr,
			}},
			DelegationFee: reward.PercentDenominator,
		}
	}

	buildGenesisArgs := api.BuildGenesisArgs{
		NetworkID:     json.Uint32(constants.UnitTestID),
		DioneAssetID:  ctx.DIONEAssetID,
		UTXOs:         genesisUTXOs,
		Validators:    genesisValidators,
		Chains:        nil,
		Time:          json.Uint64(defaultGenesisTime.Unix()),
		InitialSupply: json.Uint64(360 * units.MegaDione),
		Encoding:      formatting.Hex,
	}

	buildGenesisResponse := api.BuildGenesisReply{}
	omegavmSS := api.StaticService{}
	require.NoError(omegavmSS.BuildGenesis(nil, &buildGenesisArgs, &buildGenesisResponse))

	genesisBytes, err := formatting.Decode(buildGenesisResponse.Encoding, buildGenesisResponse.Bytes)
	require.NoError(err)

	return genesisBytes
}

func shutdownEnvironment(env *environment) error {
	if env.isBootstrapped.Get() {
		validatorIDs, err := validators.NodeIDs(env.config.Validators, constants.PrimaryNetworkID)
		if err != nil {
			return err
		}

		if err := env.uptimes.StopTracking(validatorIDs, constants.PrimaryNetworkID); err != nil {
			return err
		}
		if err := env.state.Commit(); err != nil {
			return err
		}
	}

	errs := wrappers.Errs{}
	errs.Add(
		env.state.Close(),
		env.baseDB.Close(),
	)
	return errs.Err
}
