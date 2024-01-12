// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utxo

import (
	"testing"
	"time"

	stdmath "math"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/utils/timer/mockable"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/stakeable"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
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
		expectedErr     error
	}{
		{
			description:     "no inputs, no outputs, no fee",
			utxos:           []*dione.UTXO{},
			ins:             []*dione.TransferableInput{},
			outs:            []*dione.TransferableOutput{},
			creds:           []verify.Verifiable{},
			producedAmounts: map[ids.ID]uint64{},
			expectedErr:     nil,
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
			expectedErr: ErrInsufficientUnlockedFunds,
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
			expectedErr:     errAssetIDMismatch,
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
			expectedErr:     errAssetIDMismatch,
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
			expectedErr:     ErrInsufficientUnlockedFunds,
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
			expectedErr:     errLockedFundsNotMarkedAsLocked,
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
			expectedErr:     errLocktimeMismatch,
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
			expectedErr: nil,
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
			expectedErr: errWrongNumberCredentials,
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
			expectedErr: errWrongNumberUTXOs,
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
			expectedErr: secp256k1fx.ErrNilCredential,
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
			expectedErr: secp256k1.ErrInvalidSig,
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
			expectedErr: nil,
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
			expectedErr:     nil,
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
			expectedErr: ErrInsufficientUnlockedFunds,
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
			expectedErr: nil,
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
			expectedErr: nil,
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
			expectedErr:     nil,
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
						Amt: stdmath.MaxUint64,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			expectedErr:     math.ErrOverflow,
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
			expectedErr:     ErrInsufficientLockedFunds,
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
							Amt: stdmath.MaxUint64,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			expectedErr:     math.ErrOverflow,
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
							Amt: stdmath.MaxUint64,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			producedAmounts: map[ids.ID]uint64{},
			expectedErr:     ErrInsufficientLockedFunds,
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
						Amt: stdmath.MaxUint64,
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
			expectedErr:     ErrInsufficientLockedFunds,
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
			expectedErr:     nil,
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
			expectedErr:     nil,
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
			expectedErr:     ErrInsufficientUnlockedFunds,
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
			expectedErr: ErrInsufficientUnlockedFunds,
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
			expectedErr: nil,
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
			expectedErr: ErrInsufficientUnlockedFunds,
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
			expectedErr: nil,
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
			expectedErr: ErrInsufficientUnlockedFunds,
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
				customAssetID:      1,
			},
			expectedErr: nil,
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
			expectedErr:     nil,
		},
	}

	for _, test := range tests {
		h.clk.Set(now)

		t.Run(test.description, func(t *testing.T) {
			err := h.VerifySpendUTXOs(
				&unsignedTx,
				test.utxos,
				test.ins,
				test.outs,
				test.creds,
				test.producedAmounts,
			)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}
