// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package alpha

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	stdjson "encoding/json"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/api/keystore"
	"github.com/DioneProtocol/odysseygo/chains/atomic"
	"github.com/DioneProtocol/odysseygo/database/manager"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/snow/engine/common"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/utils/cb58"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/formatting"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/utils/json"
	"github.com/DioneProtocol/odysseygo/utils/linkedhashmap"
	"github.com/DioneProtocol/odysseygo/utils/sampler"
	"github.com/DioneProtocol/odysseygo/version"
	"github.com/DioneProtocol/odysseygo/vms/alpha/block/executor"
	"github.com/DioneProtocol/odysseygo/vms/alpha/config"
	"github.com/DioneProtocol/odysseygo/vms/alpha/fxs"
	"github.com/DioneProtocol/odysseygo/vms/alpha/txs"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/feecollector"
	"github.com/DioneProtocol/odysseygo/vms/nftfx"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"

	keystoreutils "github.com/DioneProtocol/odysseygo/vms/components/keystore"
)

const (
	testTxFee    uint64 = 1000
	startBalance uint64 = 50000

	username       = "bobby"
	password       = "StrnasfqewiurPasswdn56d" //#nosec G101
	feeAssetName   = "TEST"
	otherAssetName = "OTHER"
)

var (
	testChangeAddr = ids.GenerateTestShortID()
	testCases      = []struct {
		name       string
		dioneAsset bool
	}{
		{
			name:       "genesis asset is DIONE",
			dioneAsset: true,
		},
		{
			name:       "genesis asset is TEST",
			dioneAsset: false,
		},
	}

	chainID = ids.ID{5, 4, 3, 2, 1}
	assetID = ids.ID{1, 2, 3}

	keys  []*secp256k1.PrivateKey
	addrs []ids.ShortID // addrs[i] corresponds to keys[i]

	errMissing = errors.New("missing")
)

func init() {
	factory := secp256k1.Factory{}

	for _, key := range []string{
		"24jUJ9vZexUM6expyMcT48LBx27k1m7xpraoV62oSQAHdziao5",
		"2MMvUMsxx6zsHSNXJdFD8yc5XkancvwyKPwpw4xUK3TCGDuNBY",
		"cxb7KpGWhDMALTjNNSJ7UQkkomPesyWAPUaWRGdyeBNzR6f35",
	} {
		keyBytes, _ := cb58.Decode(key)
		pk, _ := factory.ToPrivateKey(keyBytes)
		keys = append(keys, pk)
		addrs = append(addrs, pk.PublicKey().Address())
	}
}

type user struct {
	username    string
	password    string
	initialKeys []*secp256k1.PrivateKey
}

type envConfig struct {
	isCustomFeeAsset bool
	keystoreUsers    []*user
	vmStaticConfig   *config.Config
	vmDynamicConfig  *Config
	additionalFxs    []*common.Fx
	notLinearized    bool
	notBootstrapped  bool
}

type environment struct {
	genesisBytes  []byte
	genesisTx     *txs.Tx
	sharedMemory  *atomic.Memory
	issuer        chan common.Message
	vm            *VM
	service       *Service
	walletService *WalletService
}

// setup the testing environment
func setup(tb testing.TB, c *envConfig) *environment {
	require := require.New(tb)

	var (
		genesisArgs *BuildGenesisArgs
		assetName   = "DIONE"
	)
	if c.isCustomFeeAsset {
		genesisArgs = makeCustomAssetGenesis(tb)
		assetName = feeAssetName
	} else {
		genesisArgs = makeDefaultGenesis(tb)
	}

	genesisBytes := buildGenesisTestWithArgs(tb, genesisArgs)
	ctx := newContext(tb)

	baseDBManager := manager.NewMemDB(version.Semantic1_0_0)

	m := atomic.NewMemory(prefixdb.New([]byte{0}, baseDBManager.Current().Database))
	ctx.SharedMemory = m.NewSharedMemory(ctx.ChainID)

	f, err := feecollector.New(prefixdb.New([]byte{1}, baseDBManager.Current().Database))
	require.NoError(err)
	ctx.FeeCollector = f

	// NB: this lock is intentionally left locked when this function returns.
	// The caller of this function is responsible for unlocking.
	ctx.Lock.Lock()

	userKeystore, err := keystore.CreateTestKeystore()
	require.NoError(err)
	ctx.Keystore = userKeystore.NewBlockchainKeyStore(ctx.ChainID)

	for _, user := range c.keystoreUsers {
		require.NoError(userKeystore.CreateUser(user.username, user.password))

		// Import the initially funded private keys
		keystoreUser, err := keystoreutils.NewUserFromKeystore(ctx.Keystore, user.username, user.password)
		require.NoError(err)

		require.NoError(keystoreUser.PutKeys(user.initialKeys...))
		require.NoError(keystoreUser.Close())
	}

	vmStaticConfig := config.Config{
		TxFee:            testTxFee,
		CreateAssetTxFee: testTxFee,
	}
	if c.vmStaticConfig != nil {
		vmStaticConfig = *c.vmStaticConfig
	}

	vm := &VM{
		Config: vmStaticConfig,
	}

	vmDynamicConfig := Config{
		IndexTransactions: true,
	}
	if c.vmDynamicConfig != nil {
		vmDynamicConfig = *c.vmDynamicConfig
	}
	configBytes, err := stdjson.Marshal(vmDynamicConfig)
	require.NoError(err)

	require.NoError(vm.Initialize(
		context.Background(),
		ctx,
		baseDBManager.NewPrefixDBManager([]byte{1}),
		genesisBytes,
		nil,
		configBytes,
		nil,
		append(
			[]*common.Fx{
				{
					ID: secp256k1fx.ID,
					Fx: &secp256k1fx.Fx{},
				},
				{
					ID: nftfx.ID,
					Fx: &nftfx.Fx{},
				},
			},
			c.additionalFxs...,
		),
		&common.SenderTest{},
	))

	stopVertexID := ids.GenerateTestID()
	issuer := make(chan common.Message, 1)

	env := &environment{
		genesisBytes: genesisBytes,
		genesisTx:    getCreateTxFromGenesisTest(tb, genesisBytes, assetName),
		sharedMemory: m,
		issuer:       issuer,
		vm:           vm,
		service: &Service{
			vm: vm,
		},
		walletService: &WalletService{
			vm:         vm,
			pendingTxs: linkedhashmap.New[ids.ID, *txs.Tx](),
		},
	}

	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	if c.notLinearized {
		return env
	}

	require.NoError(vm.Linearize(context.Background(), stopVertexID, issuer))
	if c.notBootstrapped {
		return env
	}

	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	return env
}

func newContext(tb testing.TB) *snow.Context {
	require := require.New(tb)

	genesisBytes := buildGenesisTest(tb)
	tx := getCreateTxFromGenesisTest(tb, genesisBytes, "DIONE")

	ctx := snow.DefaultContextTest()
	ctx.NetworkID = constants.UnitTestID
	ctx.ChainID = chainID
	ctx.DIONEAssetID = tx.ID()
	ctx.AChainID = ids.Empty.Prefix(0)
	ctx.DChainID = ids.Empty.Prefix(1)
	aliaser := ctx.BCLookup.(ids.Aliaser)

	require.NoError(aliaser.Alias(chainID, "A"))
	require.NoError(aliaser.Alias(chainID, chainID.String()))
	require.NoError(aliaser.Alias(constants.OmegaChainID, "O"))
	require.NoError(aliaser.Alias(constants.OmegaChainID, constants.OmegaChainID.String()))

	ctx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				constants.OmegaChainID: ctx.SubnetID,
				chainID:                ctx.SubnetID,
			}[chainID]
			if !ok {
				return ids.Empty, errMissing
			}
			return subnetID, nil
		},
	}
	return ctx
}

// Returns:
//
//  1. tx in genesis that creates asset
//  2. the index of the output
func getCreateTxFromGenesisTest(tb testing.TB, genesisBytes []byte, assetName string) *txs.Tx {
	require := require.New(tb)

	parser, err := txs.NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
	})
	require.NoError(err)

	cm := parser.GenesisCodec()
	genesis := Genesis{}
	_, err = cm.Unmarshal(genesisBytes, &genesis)
	require.NoError(err)
	require.NotEmpty(genesis.Txs)

	var assetTx *GenesisAsset
	for _, tx := range genesis.Txs {
		if tx.Name == assetName {
			assetTx = tx
			break
		}
	}
	require.NotNil(assetTx)

	tx := &txs.Tx{
		Unsigned: &assetTx.CreateAssetTx,
	}
	require.NoError(parser.InitializeGenesisTx(tx))
	return tx
}

// buildGenesisTest is the common Genesis builder for most tests
func buildGenesisTest(tb testing.TB) []byte {
	defaultArgs := makeDefaultGenesis(tb)
	return buildGenesisTestWithArgs(tb, defaultArgs)
}

// buildGenesisTestWithArgs allows building the genesis while injecting different starting points (args)
func buildGenesisTestWithArgs(tb testing.TB, args *BuildGenesisArgs) []byte {
	require := require.New(tb)

	ss := CreateStaticService()

	reply := BuildGenesisReply{}
	require.NoError(ss.BuildGenesis(nil, args, &reply))

	b, err := formatting.Decode(reply.Encoding, reply.Bytes)
	require.NoError(err)
	return b
}

func newTx(tb testing.TB, genesisBytes []byte, vm *VM, assetName string) *txs.Tx {
	require := require.New(tb)

	createTx := getCreateTxFromGenesisTest(tb, genesisBytes, assetName)
	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    constants.UnitTestID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        createTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: createTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
		},
	}}
	require.NoError(tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}))
	return tx
}

// Sample from a set of addresses and return them raw and formatted as strings.
// The size of the sample is between 1 and len(addrs)
// If len(addrs) == 0, returns nil
func sampleAddrs(tb testing.TB, vm *VM, addrs []ids.ShortID) ([]ids.ShortID, []string) {
	require := require.New(tb)

	sampledAddrs := []ids.ShortID{}
	sampledAddrsStr := []string{}

	sampler := sampler.NewUniform()
	sampler.Initialize(uint64(len(addrs)))

	numAddrs := 1 + rand.Intn(len(addrs)) // #nosec G404
	indices, err := sampler.Sample(numAddrs)
	require.NoError(err)
	for _, index := range indices {
		addr := addrs[index]
		addrStr, err := vm.FormatLocalAddress(addr)
		require.NoError(err)

		sampledAddrs = append(sampledAddrs, addr)
		sampledAddrsStr = append(sampledAddrsStr, addrStr)
	}
	return sampledAddrs, sampledAddrsStr
}

func makeDefaultGenesis(tb testing.TB) *BuildGenesisArgs {
	require := require.New(tb)

	addr0Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[0].Bytes())
	require.NoError(err)

	addr1Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[1].Bytes())
	require.NoError(err)

	addr2Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[2].Bytes())
	require.NoError(err)

	return &BuildGenesisArgs{
		Encoding: formatting.Hex,
		GenesisData: map[string]AssetDefinition{
			"asset1": {
				Name:   "DIONE",
				Symbol: "SYMB",
				InitialState: map[string][]interface{}{
					"fixedCap": {
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr0Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr1Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr2Str,
						},
					},
				},
			},
			"asset2": {
				Name:   "myVarCapAsset",
				Symbol: "MVCA",
				InitialState: map[string][]interface{}{
					"variableCap": {
						Owners{
							Threshold: 1,
							Minters: []string{
								addr0Str,
								addr1Str,
							},
						},
						Owners{
							Threshold: 2,
							Minters: []string{
								addr0Str,
								addr1Str,
								addr2Str,
							},
						},
					},
				},
			},
			"asset3": {
				Name: "myOtherVarCapAsset",
				InitialState: map[string][]interface{}{
					"variableCap": {
						Owners{
							Threshold: 1,
							Minters: []string{
								addr0Str,
							},
						},
					},
				},
			},
			"asset4": {
				Name: "myFixedCapAsset",
				InitialState: map[string][]interface{}{
					"fixedCap": {
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr0Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr1Str,
						},
					},
				},
			},
		},
	}
}

func makeCustomAssetGenesis(tb testing.TB) *BuildGenesisArgs {
	require := require.New(tb)

	addr0Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[0].Bytes())
	require.NoError(err)

	addr1Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[1].Bytes())
	require.NoError(err)

	addr2Str, err := address.FormatBech32(constants.UnitTestHRP, addrs[2].Bytes())
	require.NoError(err)

	return &BuildGenesisArgs{
		Encoding: formatting.Hex,
		GenesisData: map[string]AssetDefinition{
			"asset1": {
				Name:   feeAssetName,
				Symbol: "TST",
				InitialState: map[string][]interface{}{
					"fixedCap": {
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr0Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr1Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr2Str,
						},
					},
				},
			},
			"asset2": {
				Name:   otherAssetName,
				Symbol: "OTH",
				InitialState: map[string][]interface{}{
					"fixedCap": {
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr0Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr1Str,
						},
						Holder{
							Amount:  json.Uint64(startBalance),
							Address: addr2Str,
						},
					},
				},
			},
		},
	}
}

// issueAndAccept expects the context lock to be held
func issueAndAccept(
	require *require.Assertions,
	vm *VM,
	issuer <-chan common.Message,
	tx *txs.Tx,
) {
	txID, err := vm.IssueTx(tx.Bytes())
	require.NoError(err)
	require.Equal(tx.ID(), txID)

	buildAndAccept(require, vm, issuer, txID)
}

// buildAndAccept expects the context lock to be held
func buildAndAccept(
	require *require.Assertions,
	vm *VM,
	issuer <-chan common.Message,
	txID ids.ID,
) {
	vm.ctx.Lock.Unlock()
	require.Equal(common.PendingTxs, <-issuer)
	vm.ctx.Lock.Lock()

	blkIntf, err := vm.BuildBlock(context.Background())
	require.NoError(err)
	require.IsType(&executor.Block{}, blkIntf)

	blk := blkIntf.(*executor.Block)
	txs := blk.Txs()
	require.Len(txs, 1)

	issuedTx := txs[0]
	require.Equal(txID, issuedTx.ID())
	require.NoError(blk.Verify(context.Background()))
	require.NoError(vm.SetPreference(context.Background(), blk.ID()))
	require.NoError(blk.Accept(context.Background()))
}
