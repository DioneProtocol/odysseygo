// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utxo

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/utils/crypto/secp256k1"
	"github.com/dioneprotocol/dionego/utils/timer/mockable"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/components/verify"
	"github.com/dioneprotocol/dionego/vms/platformvm/stakeable"
	"github.com/dioneprotocol/dionego/vms/platformvm/txs"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

var _ txs.UnsignedTx = (*dummyUnsignedTx)(nil)

type dummyUnsignedTx struct {
	txs.BaseTx
}

func (*dummyUnsignedTx) Visit(txs.Visitor) error {
	return nil
}

func TestVerifySpendUTXOs(t *testing.T) {
	fx := &secp256k1fx.Fx{}

	require.NoError(t, fx.InitializeVM(&secp256k1fx.TestVM{}))
	require.NoError(t, fx.Bootstrapped())

	h := &handler{
		ctx: snow.DefaultContextTest(),
		clk: &mockable.Clock{},
		fx:  fx,
	}

	// The handler time during a test, unless [chainTimestamp] is set
	now := time.Unix(1607133207, 0)

	unsignedTx := dummyUnsignedTx{
		BaseTx: txs.BaseTx{},
	}
	unsignedTx.SetBytes([]byte{0})

	customAssetID := ids.GenerateTestID()

	// Note that setting [chainTimestamp] also set's the handler's clock.
	// Adjust input/output locktimes accordingly.
	tests := []struct {
		description     string
		utxos           []*dione.UTXO
		ins             []*dione.TransferableInput
		outs            []*dione.TransferableOutput
		creds           []verify.Verifiable
		producedAmounts map[ids.ID]uint64
		shouldErr       bool
	}{
		{
			description:     "no inputs, no outputs, no fee",
			utxos:           []*dione.UTXO{},
			ins:             []*dione.TransferableInput{},
			outs:            []*dione.TransferableOutput{},
			creds:           []verify.Verifiable{},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       false,
		},
		{
			description: "no inputs, no outputs, positive fee",
			utxos:       []*dione.UTXO{},
			ins:         []*dione.TransferableInput{},
			outs:        []*dione.TransferableOutput{},
			creds:       []verify.Verifiable{},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "wrong utxo assetID, one input, no outputs, no fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: customAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "one wrong assetID input, no outputs, no fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: customAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "one input, one wrong assetID output, no fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempt to consume locked output as unlocked",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Add(time.Second).Unix()),
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempt to modify locktime",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Add(time.Second).Unix()),
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &stakeable.LockIn{
					Locktime: uint64(now.Unix()),
					TransferableIn: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "one input, no outputs, positive fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "wrong number of credentials",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs:  []*dione.TransferableOutput{},
			creds: []verify.Verifiable{},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "wrong number of UTXOs",
			utxos:       []*dione.UTXO{},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "invalid credential",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				(*secp256k1fx.Credential)(nil),
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "invalid signature",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs: []ids.ShortID{
							ids.GenerateTestShortID(),
						},
					},
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{0},
					},
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{
					Sigs: [][secp256k1.SignatureLen]byte{
						{},
					},
				},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "one input, no outputs, positive fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "locked one input, no outputs, no fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Unix()) + 1,
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &stakeable.LockIn{
					Locktime: uint64(now.Unix()) + 1,
					TransferableIn: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       false,
		},
		{
			description: "locked one input, no outputs, positive fee",
			utxos: []*dione.UTXO{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Unix()) + 1,
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*dione.TransferableInput{{
				Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
				In: &stakeable.LockIn{
					Locktime: uint64(now.Unix()) + 1,
					TransferableIn: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			}},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "one locked and one unlocked input, one locked output, positive fee",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &stakeable.LockIn{
						Locktime: uint64(now.Unix()) + 1,
						TransferableIn: &secp256k1fx.TransferInput{
							Amt: 1,
						},
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "one locked and one unlocked input, one locked output, positive fee, partially locked",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 2,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &stakeable.LockIn{
						Locktime: uint64(now.Unix()) + 1,
						TransferableIn: &secp256k1fx.TransferInput{
							Amt: 1,
						},
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 2,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 2,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "one unlocked input, one locked output, zero fee",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) - 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       false,
		},
		{
			description: "attempted overflow",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 2,
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: math.MaxUint64,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempted mint",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 2,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempted mint through locking",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 2,
						},
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: math.MaxUint64,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempted mint through mixed locking (low then high)",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 2,
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: math.MaxUint64,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempted mint through mixed locking (high then low)",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: math.MaxUint64,
					},
				},
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &stakeable.LockOut{
						Locktime: 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 2,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "transfer non-dione asset",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       false,
		},
		{
			description: "lock non-dione asset",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Add(time.Second).Unix()),
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       false,
		},
		{
			description: "attempted asset conversion",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			shouldErr:       true,
		},
		{
			description: "attempted asset conversion with burn",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "two inputs, one output with custom asset, with fee",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "one input, fee, custom asset",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "one input, custom fee",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				customAssetID: 1,
			},
			shouldErr: false,
		},
		{
			description: "one input, custom fee, wrong burn",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				customAssetID: 1,
			},
			shouldErr: true,
		},
		{
			description: "two inputs, multiple fee",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: h.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{
				h.ctx.DIONEAssetID: 1,
				customAssetID:     1,
			},
			shouldErr: false,
		},
		{
			description: "one unlock input, one locked output, zero fee, unlocked, custom asset",
			utxos: []*dione.UTXO{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) - 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			ins: []*dione.TransferableInput{
				{
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*dione.TransferableOutput{
				{
					Asset: dione.Asset{ID: customAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: make(map[ids.ID]uint64),
			shouldErr:       false,
		},
	}

	for _, test := range tests {
		h.clk.Set(now)

		t.Run(test.description, func(t *testing.T) {
			require := require.New(t)
			err := h.VerifySpendUTXOs(
				&unsignedTx,
				test.utxos,
				test.ins,
				test.outs,
				test.creds,
				test.producedAmounts,
			)

			if test.shouldErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
