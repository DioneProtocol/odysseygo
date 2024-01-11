// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/engine/common"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/avm/txs"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

func TestSetsAndGets(t *testing.T) {
	require := require.New(t)

	env := setup(t, &envConfig{
		additionalFxs: []*common.Fx{{
			ID: ids.GenerateTestID(),
			Fx: &FxTest{
				InitializeF: func(vmIntf interface{}) error {
					vm := vmIntf.(secp256k1fx.VM)
					return vm.CodecRegistry().RegisterType(&dione.TestState{})
				},
			},
		}},
	})
	defer func() {
		require.NoError(env.vm.Shutdown(context.Background()))
		env.vm.ctx.Lock.Unlock()
	}()

	utxo := &dione.UTXO{
		UTXOID: dione.UTXOID{
			TxID:        ids.Empty,
			OutputIndex: 1,
		},
		Asset: dione.Asset{ID: ids.Empty},
		Out:   &dione.TestState{},
	}
	utxoID := utxo.InputID()

	tx := &txs.Tx{Unsigned: &txs.BaseTx{BaseTx: dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: chainID,
		Ins: []*dione.TransferableInput{{
			UTXOID: dione.UTXOID{
				TxID:        ids.Empty,
				OutputIndex: 0,
			},
			Asset: dione.Asset{ID: assetID},
			In: &secp256k1fx.TransferInput{
				Amt: 20 * units.KiloDione,
				Input: secp256k1fx.Input{
					SigIndices: []uint32{
						0,
					},
				},
			},
		}},
	}}}
	require.NoError(tx.SignSECP256K1Fx(env.vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}))

	txID := tx.ID()

	env.vm.state.AddUTXO(utxo)
	env.vm.state.AddTx(tx)

	resultUTXO, err := env.vm.state.GetUTXO(utxoID)
	require.NoError(err)
	resultTx, err := env.vm.state.GetTx(txID)
	require.NoError(err)

	require.Equal(uint32(1), resultUTXO.OutputIndex)
	require.Equal(tx.ID(), resultTx.ID())
}

func TestFundingNoAddresses(t *testing.T) {
	env := setup(t, &envConfig{
		additionalFxs: []*common.Fx{{
			ID: ids.GenerateTestID(),
			Fx: &FxTest{
				InitializeF: func(vmIntf interface{}) error {
					vm := vmIntf.(secp256k1fx.VM)
					return vm.CodecRegistry().RegisterType(&dione.TestState{})
				},
			},
		}},
	})
	defer func() {
		require.NoError(t, env.vm.Shutdown(context.Background()))
		env.vm.ctx.Lock.Unlock()
	}()

	utxo := &dione.UTXO{
		UTXOID: dione.UTXOID{
			TxID:        ids.Empty,
			OutputIndex: 1,
		},
		Asset: dione.Asset{ID: ids.Empty},
		Out:   &dione.TestState{},
	}

	env.vm.state.AddUTXO(utxo)
	env.vm.state.DeleteUTXO(utxo.InputID())
}

func TestFundingAddresses(t *testing.T) {
	require := require.New(t)

	env := setup(t, &envConfig{
		additionalFxs: []*common.Fx{{
			ID: ids.GenerateTestID(),
			Fx: &FxTest{
				InitializeF: func(vmIntf interface{}) error {
					vm := vmIntf.(secp256k1fx.VM)
					return vm.CodecRegistry().RegisterType(&dione.TestAddressable{})
				},
			},
		}},
	})
	defer func() {
		require.NoError(env.vm.Shutdown(context.Background()))
		env.vm.ctx.Lock.Unlock()
	}()

	utxo := &dione.UTXO{
		UTXOID: dione.UTXOID{
			TxID:        ids.Empty,
			OutputIndex: 1,
		},
		Asset: dione.Asset{ID: ids.Empty},
		Out: &dione.TestAddressable{
			Addrs: [][]byte{{0}},
		},
	}

	env.vm.state.AddUTXO(utxo)
	require.NoError(env.vm.state.Commit())

	utxos, err := env.vm.state.UTXOIDs([]byte{0}, ids.Empty, math.MaxInt32)
	require.NoError(err)
	require.Len(utxos, 1)
	require.Equal(utxo.InputID(), utxos[0])

	env.vm.state.DeleteUTXO(utxo.InputID())
	require.NoError(env.vm.state.Commit())

	utxos, err = env.vm.state.UTXOIDs([]byte{0}, ids.Empty, math.MaxInt32)
	require.NoError(err)
	require.Empty(utxos)
}
