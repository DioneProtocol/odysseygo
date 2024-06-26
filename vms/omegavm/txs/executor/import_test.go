// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/chains/atomic"
	"github.com/DioneProtocol/odysseygo/database/prefixdb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/utxo"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

func TestNewImportTx(t *testing.T) {
	env := newEnvironment(t, false /*=postBanff*/, false /*=postCortina*/)
	defer func() {
		require.NoError(t, shutdownEnvironment(env))
	}()

	type test struct {
		description   string
		sourceChainID ids.ID
		sharedMemory  atomic.SharedMemory
		sourceKeys    []*secp256k1.PrivateKey
		timestamp     time.Time
		expectedErr   error
	}

	factory := secp256k1.Factory{}
	sourceKey, err := factory.NewPrivateKey()
	require.NoError(t, err)

	cnt := new(byte)

	// Returns a shared memory where GetDatabase returns a database
	// where [recipientKey] has a balance of [amt]
	fundedSharedMemory := func(peerChain ids.ID, assets map[ids.ID]uint64) atomic.SharedMemory {
		*cnt++
		m := atomic.NewMemory(prefixdb.New([]byte{*cnt}, env.baseDB))

		sm := m.NewSharedMemory(env.ctx.ChainID)
		peerSharedMemory := m.NewSharedMemory(peerChain)

		for assetID, amt := range assets {
			// #nosec G404
			utxo := &dione.UTXO{
				UTXOID: dione.UTXOID{
					TxID:        ids.GenerateTestID(),
					OutputIndex: rand.Uint32(),
				},
				Asset: dione.Asset{ID: assetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: amt,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Addrs:     []ids.ShortID{sourceKey.PublicKey().Address()},
						Threshold: 1,
					},
				},
			}
			utxoBytes, err := txs.Codec.Marshal(txs.Version, utxo)
			require.NoError(t, err)

			inputID := utxo.InputID()
			require.NoError(t, peerSharedMemory.Apply(map[ids.ID]*atomic.Requests{
				env.ctx.ChainID: {
					PutRequests: []*atomic.Element{
						{
							Key:   inputID[:],
							Value: utxoBytes,
							Traits: [][]byte{
								sourceKey.PublicKey().Address().Bytes(),
							},
						},
					},
				},
			}))
		}

		return sm
	}

	customAssetID := ids.GenerateTestID()

	tests := []test{
		{
			description:   "can't pay fee",
			sourceChainID: env.ctx.AChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.AChainID,
				map[ids.ID]uint64{
					env.ctx.DIONEAssetID: env.config.TxFee - 1,
				},
			),
			sourceKeys:  []*secp256k1.PrivateKey{sourceKey},
			expectedErr: utxo.ErrInsufficientFunds,
		},
		{
			description:   "can barely pay fee",
			sourceChainID: env.ctx.AChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.AChainID,
				map[ids.ID]uint64{
					env.ctx.DIONEAssetID: env.config.TxFee,
				},
			),
			sourceKeys:  []*secp256k1.PrivateKey{sourceKey},
			expectedErr: nil,
		},
		{
			description:   "attempting to import from D-chain",
			sourceChainID: dChainID,
			sharedMemory: fundedSharedMemory(
				dChainID,
				map[ids.ID]uint64{
					env.ctx.DIONEAssetID: env.config.TxFee,
				},
			),
			sourceKeys:  []*secp256k1.PrivateKey{sourceKey},
			timestamp:   env.config.ApricotPhase5Time,
			expectedErr: nil,
		},
		{
			description:   "attempting to import non-dione from A-chain",
			sourceChainID: env.ctx.AChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.AChainID,
				map[ids.ID]uint64{
					env.ctx.DIONEAssetID: env.config.TxFee,
					customAssetID:        1,
				},
			),
			sourceKeys:  []*secp256k1.PrivateKey{sourceKey},
			timestamp:   env.config.BanffTime,
			expectedErr: nil,
		},
	}

	to := ids.GenerateTestShortID()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require := require.New(t)

			env.msm.SharedMemory = tt.sharedMemory
			tx, err := env.txBuilder.NewImportTx(
				tt.sourceChainID,
				to,
				tt.sourceKeys,
				ids.ShortEmpty,
			)
			require.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}
			require.NoError(err)

			unsignedTx := tx.Unsigned.(*txs.ImportTx)
			require.NotEmpty(unsignedTx.ImportedInputs)
			numInputs := len(unsignedTx.Ins) + len(unsignedTx.ImportedInputs)
			require.Equal(len(tx.Creds), numInputs, "should have the same number of credentials as inputs")

			totalIn := uint64(0)
			for _, in := range unsignedTx.Ins {
				totalIn += in.Input().Amount()
			}
			for _, in := range unsignedTx.ImportedInputs {
				totalIn += in.Input().Amount()
			}
			totalOut := uint64(0)
			for _, out := range unsignedTx.Outs {
				totalOut += out.Out.Amount()
			}

			require.Equal(env.config.TxFee, totalIn-totalOut)

			fakedState, err := state.NewDiff(lastAcceptedID, env)
			require.NoError(err)

			fakedState.SetTimestamp(tt.timestamp)

			fakedParent := ids.GenerateTestID()
			env.SetState(fakedParent, fakedState)

			verifier := MempoolTxVerifier{
				Backend:       &env.backend,
				ParentID:      fakedParent,
				StateVersions: env,
				Tx:            tx,
			}
			require.NoError(tx.Unsigned.Visit(&verifier))
		})
	}
}
