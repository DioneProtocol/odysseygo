// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"strings"
	"testing"

	stdmath "math"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/vms/alpha/config"
	"github.com/DioneProtocol/odysseygo/vms/alpha/fxs"
	"github.com/DioneProtocol/odysseygo/vms/alpha/txs"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var (
	keys      = secp256k1.TestKeys()
	feeConfig = config.Config{
		TxFee:            2,
		CreateAssetTxFee: 3,
	}
)

func newContext(t testing.TB) *snow.Context {
	require := require.New(t)

	ctx := snow.DefaultContextTest()
	ctx.NetworkID = constants.UnitTestID
	ctx.ChainID = ids.GenerateTestID()
	ctx.AChainID = ctx.ChainID
	ctx.DChainID = ids.GenerateTestID()

	aliaser := ctx.BCLookup.(ids.Aliaser)
	require.NoError(aliaser.Alias(ctx.AChainID, "A"))
	require.NoError(aliaser.Alias(ctx.AChainID, ctx.AChainID.String()))
	require.NoError(aliaser.Alias(constants.OmegaChainID, "O"))
	require.NoError(aliaser.Alias(constants.OmegaChainID, constants.OmegaChainID.String()))
	return ctx
}

func TestSyntacticVerifierBaseTx(t *testing.T) {
	ctx := newContext(t)

	fx := &secp256k1fx.Fx{}
	parser, err := txs.NewParser([]fxs.Fx{
		fx,
	})
	require.NoError(t, err)

	feeAssetID := ids.GenerateTestID()
	asset := dione.Asset{
		ID: feeAssetID,
	}
	outputOwners := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
	}
	fxOutput := secp256k1fx.TransferOutput{
		Amt:          12345,
		OutputOwners: outputOwners,
	}
	output := dione.TransferableOutput{
		Asset: asset,
		Out:   &fxOutput,
	}
	inputTxID := ids.GenerateTestID()
	utxoID := dione.UTXOID{
		TxID:        inputTxID,
		OutputIndex: 0,
	}
	inputSigners := secp256k1fx.Input{
		SigIndices: []uint32{2},
	}
	fxInput := secp256k1fx.TransferInput{
		Amt:   54321,
		Input: inputSigners,
	}
	input := dione.TransferableInput{
		UTXOID: utxoID,
		Asset:  asset,
		In:     &fxInput,
	}
	baseTx := dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: ctx.ChainID,
		Outs: []*dione.TransferableOutput{
			&output,
		},
		Ins: []*dione.TransferableInput{
			&input,
		},
	}
	cred := fxs.FxCredential{
		Verifiable: &secp256k1fx.Credential{},
	}
	creds := []*fxs.FxCredential{
		&cred,
	}

	codec := parser.Codec()
	backend := &Backend{
		Ctx:    ctx,
		Config: &feeConfig,
		Fxs: []*fxs.ParsedFx{
			{
				ID: secp256k1fx.ID,
				Fx: fx,
			},
		},
		Codec:      codec,
		FeeAssetID: feeAssetID,
	}

	tests := []struct {
		name   string
		txFunc func() *txs.Tx
		err    error
	}{
		{
			name: "valid",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "wrong networkID",
			txFunc: func() *txs.Tx {
				baseTx := baseTx
				baseTx.NetworkID++
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrWrongNetworkID,
		},
		{
			name: "wrong chainID",
			txFunc: func() *txs.Tx {
				baseTx := baseTx
				baseTx.BlockchainID = ids.GenerateTestID()
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrWrongChainID,
		},
		{
			name: "memo too large",
			txFunc: func() *txs.Tx {
				baseTx := baseTx
				baseTx.Memo = make([]byte, dione.MaxMemoSize+1)
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrMemoTooLarge,
		},
		{
			name: "invalid output",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          0,
					OutputOwners: outputOwners,
				}

				baseTx := baseTx
				baseTx.Outs = []*dione.TransferableOutput{
					&output,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				baseTx := baseTx
				baseTx.Outs = outputs
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "invalid input",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   0,
					Input: inputSigners,
				}

				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueInput,
		},
		{
			name: "duplicate inputs",
			txFunc: func() *txs.Tx {
				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "input overflow",
			txFunc: func() *txs.Tx {
				input0 := input
				input0.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				input1 := input
				input1.UTXOID.OutputIndex++
				input1.In = &secp256k1fx.TransferInput{
					Amt:   stdmath.MaxUint64,
					Input: inputSigners,
				}

				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input0,
					&input1,
				}
				dione.SortTransferableInputsWithSigners(baseTx.Ins, make([][]*secp256k1.PrivateKey, 2))
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "output overflow",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          stdmath.MaxUint64,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)

				baseTx := baseTx
				baseTx.Outs = outputs
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
		{
			name: "invalid credential",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds: []*fxs.FxCredential{{
						Verifiable: (*secp256k1fx.Credential)(nil),
					}},
				}
			},
			err: secp256k1fx.ErrNilCredential,
		},
		{
			name: "wrong number of credentials",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
				}
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "barely sufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee,
					Input: inputSigners,
				}

				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "barely insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee - 1,
					Input: inputSigners,
				}

				baseTx := baseTx
				baseTx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &txs.BaseTx{BaseTx: baseTx},
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := test.txFunc()
			verifier := &SyntacticVerifier{
				Backend: backend,
				Tx:      tx,
			}
			err := tx.Unsigned.Visit(verifier)
			require.ErrorIs(t, err, test.err)
		})
	}
}

func TestSyntacticVerifierCreateAssetTx(t *testing.T) {
	ctx := newContext(t)

	fx := &secp256k1fx.Fx{}
	parser, err := txs.NewParser([]fxs.Fx{
		fx,
	})
	require.NoError(t, err)

	feeAssetID := ids.GenerateTestID()
	asset := dione.Asset{
		ID: feeAssetID,
	}
	outputOwners := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
	}
	fxOutput := secp256k1fx.TransferOutput{
		Amt:          12345,
		OutputOwners: outputOwners,
	}
	output := dione.TransferableOutput{
		Asset: asset,
		Out:   &fxOutput,
	}
	inputTxID := ids.GenerateTestID()
	utxoID := dione.UTXOID{
		TxID:        inputTxID,
		OutputIndex: 0,
	}
	inputSigners := secp256k1fx.Input{
		SigIndices: []uint32{2},
	}
	fxInput := secp256k1fx.TransferInput{
		Amt:   54321,
		Input: inputSigners,
	}
	input := dione.TransferableInput{
		UTXOID: utxoID,
		Asset:  asset,
		In:     &fxInput,
	}
	baseTx := dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: ctx.ChainID,
		Outs: []*dione.TransferableOutput{
			&output,
		},
		Ins: []*dione.TransferableInput{
			&input,
		},
	}
	initialState := txs.InitialState{
		FxIndex: 0,
		Outs: []verify.State{
			&fxOutput,
		},
	}
	tx := txs.CreateAssetTx{
		BaseTx:       txs.BaseTx{BaseTx: baseTx},
		Name:         "NormalName",
		Symbol:       "TICK",
		Denomination: byte(2),
		States: []*txs.InitialState{
			&initialState,
		},
	}
	cred := fxs.FxCredential{
		Verifiable: &secp256k1fx.Credential{},
	}
	creds := []*fxs.FxCredential{
		&cred,
	}

	codec := parser.Codec()
	backend := &Backend{
		Ctx:    ctx,
		Config: &feeConfig,
		Fxs: []*fxs.ParsedFx{
			{
				ID: secp256k1fx.ID,
				Fx: fx,
			},
		},
		Codec:      codec,
		FeeAssetID: feeAssetID,
	}

	tests := []struct {
		name   string
		txFunc func() *txs.Tx
		err    error
	}{
		{
			name: "valid",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "name too short",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Name = ""
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNameTooShort,
		},
		{
			name: "name too long",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Name = strings.Repeat("A", maxNameLen+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNameTooLong,
		},
		{
			name: "symbol too short",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Symbol = ""
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errSymbolTooShort,
		},
		{
			name: "symbol too long",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Symbol = strings.Repeat("A", maxSymbolLen+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errSymbolTooLong,
		},
		{
			name: "no feature extensions",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.States = nil
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNoFxs,
		},
		{
			name: "denomination too large",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Denomination = maxDenomination + 1
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errDenominationTooLarge,
		},
		{
			name: "bounding whitespace in name",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Name = " DIONE"
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errUnexpectedWhitespace,
		},
		{
			name: "illegal character in name",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Name = "h8*32"
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errIllegalNameCharacter,
		},
		{
			name: "illegal character in ticker",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Symbol = "H I"
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errIllegalSymbolCharacter,
		},
		{
			name: "wrong networkID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.NetworkID++
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongNetworkID,
		},
		{
			name: "wrong chainID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.BlockchainID = ids.GenerateTestID()
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongChainID,
		},
		{
			name: "memo too large",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Memo = make([]byte, dione.MaxMemoSize+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrMemoTooLarge,
		},
		{
			name: "invalid output",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          0,
					OutputOwners: outputOwners,
				}

				tx := tx
				tx.Outs = []*dione.TransferableOutput{
					&output,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "invalid input",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   0,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueInput,
		},
		{
			name: "duplicate inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "input overflow",
			txFunc: func() *txs.Tx {
				input0 := input
				input0.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				input1 := input
				input1.UTXOID.OutputIndex++
				input1.In = &secp256k1fx.TransferInput{
					Amt:   stdmath.MaxUint64,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input0,
					&input1,
				}
				dione.SortTransferableInputsWithSigners(baseTx.Ins, make([][]*secp256k1.PrivateKey, 2))
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "output overflow",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          stdmath.MaxUint64,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
		{
			name: "invalid nil state",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.States = []*txs.InitialState{
					nil,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrNilInitialState,
		},
		{
			name: "invalid fx",
			txFunc: func() *txs.Tx {
				initialState := initialState
				initialState.FxIndex = 1

				tx := tx
				tx.States = []*txs.InitialState{
					&initialState,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrUnknownFx,
		},
		{
			name: "invalid nil state output",
			txFunc: func() *txs.Tx {
				initialState := initialState
				initialState.Outs = []verify.State{
					nil,
				}

				tx := tx
				tx.States = []*txs.InitialState{
					&initialState,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrNilFxOutput,
		},
		{
			name: "invalid state output",
			txFunc: func() *txs.Tx {
				fxOutput := fxOutput
				fxOutput.Amt = 0

				initialState := initialState
				initialState.Outs = []verify.State{
					&fxOutput,
				}

				tx := tx
				tx.States = []*txs.InitialState{
					&initialState,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted initial state",
			txFunc: func() *txs.Tx {
				fxOutput0 := fxOutput

				fxOutput1 := fxOutput
				fxOutput1.Amt++

				initialState := initialState
				initialState.Outs = []verify.State{
					&fxOutput0,
					&fxOutput1,
				}
				initialState.Sort(codec)
				initialState.Outs[0], initialState.Outs[1] = initialState.Outs[1], initialState.Outs[0]

				tx := tx
				tx.States = []*txs.InitialState{
					&initialState,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrOutputsNotSorted,
		},
		{
			name: "non-unique initial states",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.States = []*txs.InitialState{
					&initialState,
					&initialState,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errInitialStatesNotSortedUnique,
		},
		{
			name: "invalid credential",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{{
						Verifiable: (*secp256k1fx.Credential)(nil),
					}},
				}
			},
			err: secp256k1fx.ErrNilCredential,
		},
		{
			name: "wrong number of credentials",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
				}
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "barely sufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.CreateAssetTxFee,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "barely insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.CreateAssetTxFee - 1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := test.txFunc()
			verifier := &SyntacticVerifier{
				Backend: backend,
				Tx:      tx,
			}
			err := tx.Unsigned.Visit(verifier)
			require.ErrorIs(t, err, test.err)
		})
	}
}

func TestSyntacticVerifierOperationTx(t *testing.T) {
	ctx := newContext(t)

	fx := &secp256k1fx.Fx{}
	parser, err := txs.NewParser([]fxs.Fx{
		fx,
	})
	require.NoError(t, err)

	feeAssetID := ids.GenerateTestID()
	asset := dione.Asset{
		ID: feeAssetID,
	}
	outputOwners := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
	}
	fxOutput := secp256k1fx.TransferOutput{
		Amt:          12345,
		OutputOwners: outputOwners,
	}
	output := dione.TransferableOutput{
		Asset: asset,
		Out:   &fxOutput,
	}
	inputTxID := ids.GenerateTestID()
	utxoID := dione.UTXOID{
		TxID:        inputTxID,
		OutputIndex: 0,
	}
	inputSigners := secp256k1fx.Input{
		SigIndices: []uint32{2},
	}
	fxInput := secp256k1fx.TransferInput{
		Amt:   54321,
		Input: inputSigners,
	}
	input := dione.TransferableInput{
		UTXOID: utxoID,
		Asset:  asset,
		In:     &fxInput,
	}
	baseTx := dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: ctx.ChainID,
		Ins: []*dione.TransferableInput{
			&input,
		},
		Outs: []*dione.TransferableOutput{
			&output,
		},
	}
	opUTXOID := utxoID
	opUTXOID.OutputIndex++
	fxOp := secp256k1fx.MintOperation{
		MintInput: inputSigners,
		MintOutput: secp256k1fx.MintOutput{
			OutputOwners: outputOwners,
		},
		TransferOutput: fxOutput,
	}
	op := txs.Operation{
		Asset: asset,
		UTXOIDs: []*dione.UTXOID{
			&opUTXOID,
		},
		Op: &fxOp,
	}
	tx := txs.OperationTx{
		BaseTx: txs.BaseTx{BaseTx: baseTx},
		Ops: []*txs.Operation{
			&op,
		},
	}
	cred := fxs.FxCredential{
		Verifiable: &secp256k1fx.Credential{},
	}
	creds := []*fxs.FxCredential{
		&cred,
		&cred,
	}

	codec := parser.Codec()
	backend := &Backend{
		Ctx:    ctx,
		Config: &feeConfig,
		Fxs: []*fxs.ParsedFx{
			{
				ID: secp256k1fx.ID,
				Fx: fx,
			},
		},
		Codec:      codec,
		FeeAssetID: feeAssetID,
	}

	tests := []struct {
		name   string
		txFunc func() *txs.Tx
		err    error
	}{
		{
			name: "valid",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "no operation",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ops = nil
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNoOperations,
		},
		{
			name: "wrong networkID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.NetworkID++
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongNetworkID,
		},
		{
			name: "wrong chainID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.BlockchainID = ids.GenerateTestID()
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongChainID,
		},
		{
			name: "memo too large",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Memo = make([]byte, dione.MaxMemoSize+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrMemoTooLarge,
		},
		{
			name: "invalid output",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          0,
					OutputOwners: outputOwners,
				}

				tx := tx
				tx.Outs = []*dione.TransferableOutput{
					&output,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "invalid input",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   0,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueInput,
		},
		{
			name: "duplicate inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "input overflow",
			txFunc: func() *txs.Tx {
				input0 := input
				input0.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				input1 := input
				input1.UTXOID.OutputIndex++
				input1.In = &secp256k1fx.TransferInput{
					Amt:   stdmath.MaxUint64,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input0,
					&input1,
				}
				dione.SortTransferableInputsWithSigners(tx.Ins, make([][]*secp256k1.PrivateKey, 2))
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "output overflow",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          stdmath.MaxUint64,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output,
				}
				dione.SortTransferableOutputs(outputs, codec)

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
		{
			name: "invalid nil op",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ops = []*txs.Operation{
					nil,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrNilOperation,
		},
		{
			name: "invalid nil fx op",
			txFunc: func() *txs.Tx {
				op := op
				op.Op = nil

				tx := tx
				tx.Ops = []*txs.Operation{
					&op,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrNilFxOperation,
		},
		{
			name: "invalid duplicated op UTXOs",
			txFunc: func() *txs.Tx {
				op := op
				op.UTXOIDs = []*dione.UTXOID{
					&opUTXOID,
					&opUTXOID,
				}

				tx := tx
				tx.Ops = []*txs.Operation{
					&op,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: txs.ErrNotSortedAndUniqueUTXOIDs,
		},
		{
			name: "invalid duplicated UTXOs across ops",
			txFunc: func() *txs.Tx {
				newOp := op
				op.Asset.ID = ids.GenerateTestID()

				tx := tx
				tx.Ops = []*txs.Operation{
					&op,
					&newOp,
				}
				txs.SortOperations(tx.Ops, codec)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errDoubleSpend,
		},
		{
			name: "invalid duplicated op",
			txFunc: func() *txs.Tx {
				op := op
				op.UTXOIDs = nil

				tx := tx
				tx.Ops = []*txs.Operation{
					&op,
					&op,
				}
				txs.SortOperations(tx.Ops, codec)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errOperationsNotSortedUnique,
		},
		{
			name: "invalid credential",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{{
						Verifiable: (*secp256k1fx.Credential)(nil),
					}},
				}
			},
			err: secp256k1fx.ErrNilCredential,
		},
		{
			name: "wrong number of credentials",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
				}
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "barely sufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "barely insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee - 1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := test.txFunc()
			verifier := &SyntacticVerifier{
				Backend: backend,
				Tx:      tx,
			}
			err := tx.Unsigned.Visit(verifier)
			require.ErrorIs(t, err, test.err)
		})
	}
}

func TestSyntacticVerifierImportTx(t *testing.T) {
	ctx := newContext(t)

	fx := &secp256k1fx.Fx{}
	parser, err := txs.NewParser([]fxs.Fx{
		fx,
	})
	require.NoError(t, err)

	feeAssetID := ids.GenerateTestID()
	asset := dione.Asset{
		ID: feeAssetID,
	}
	outputOwners := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
	}
	fxOutput := secp256k1fx.TransferOutput{
		Amt:          12345,
		OutputOwners: outputOwners,
	}
	output := dione.TransferableOutput{
		Asset: asset,
		Out:   &fxOutput,
	}
	inputTxID := ids.GenerateTestID()
	utxoID := dione.UTXOID{
		TxID:        inputTxID,
		OutputIndex: 0,
	}
	inputSigners := secp256k1fx.Input{
		SigIndices: []uint32{2},
	}
	fxInput := secp256k1fx.TransferInput{
		Amt:   54321,
		Input: inputSigners,
	}
	input := dione.TransferableInput{
		UTXOID: utxoID,
		Asset:  asset,
		In:     &fxInput,
	}
	baseTx := dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: ctx.ChainID,
		Outs: []*dione.TransferableOutput{
			&output,
		},
	}
	tx := txs.ImportTx{
		BaseTx:      txs.BaseTx{BaseTx: baseTx},
		SourceChain: ctx.DChainID,
		ImportedIns: []*dione.TransferableInput{
			&input,
		},
	}
	cred := fxs.FxCredential{
		Verifiable: &secp256k1fx.Credential{},
	}
	creds := []*fxs.FxCredential{
		&cred,
	}

	codec := parser.Codec()
	backend := &Backend{
		Ctx:    ctx,
		Config: &feeConfig,
		Fxs: []*fxs.ParsedFx{
			{
				ID: secp256k1fx.ID,
				Fx: fx,
			},
		},
		Codec:      codec,
		FeeAssetID: feeAssetID,
	}

	tests := []struct {
		name   string
		txFunc func() *txs.Tx
		err    error
	}{
		{
			name: "valid",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "no imported inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.ImportedIns = nil
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNoImportInputs,
		},
		{
			name: "wrong networkID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.NetworkID++
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongNetworkID,
		},
		{
			name: "wrong chainID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.BlockchainID = ids.GenerateTestID()
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongChainID,
		},
		{
			name: "memo too large",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Memo = make([]byte, dione.MaxMemoSize+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrMemoTooLarge,
		},
		{
			name: "invalid output",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          0,
					OutputOwners: outputOwners,
				}

				tx := tx
				tx.Outs = []*dione.TransferableOutput{
					&output,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "invalid input",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   0,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueInput,
		},
		{
			name: "duplicate inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "duplicate imported inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.ImportedIns = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "input overflow",
			txFunc: func() *txs.Tx {
				input0 := input
				input0.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				input1 := input
				input1.UTXOID.OutputIndex++
				input1.In = &secp256k1fx.TransferInput{
					Amt:   stdmath.MaxUint64,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input0,
					&input1,
				}
				dione.SortTransferableInputsWithSigners(tx.Ins, make([][]*secp256k1.PrivateKey, 2))
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "output overflow",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          stdmath.MaxUint64,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output,
				}
				dione.SortTransferableOutputs(outputs, codec)

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				tx := tx
				tx.ImportedIns = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
		{
			name: "invalid credential",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{{
						Verifiable: (*secp256k1fx.Credential)(nil),
					}},
				}
			},
			err: secp256k1fx.ErrNilCredential,
		},
		{
			name: "wrong number of credentials",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
				}
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "barely sufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee,
					Input: inputSigners,
				}

				tx := tx
				tx.ImportedIns = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "barely insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee - 1,
					Input: inputSigners,
				}

				tx := tx
				tx.ImportedIns = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := test.txFunc()
			verifier := &SyntacticVerifier{
				Backend: backend,
				Tx:      tx,
			}
			err := tx.Unsigned.Visit(verifier)
			require.ErrorIs(t, err, test.err)
		})
	}
}

func TestSyntacticVerifierExportTx(t *testing.T) {
	ctx := newContext(t)

	fx := &secp256k1fx.Fx{}
	parser, err := txs.NewParser([]fxs.Fx{
		fx,
	})
	require.NoError(t, err)

	feeAssetID := ids.GenerateTestID()
	asset := dione.Asset{
		ID: feeAssetID,
	}
	outputOwners := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
	}
	fxOutput := secp256k1fx.TransferOutput{
		Amt:          12345,
		OutputOwners: outputOwners,
	}
	output := dione.TransferableOutput{
		Asset: asset,
		Out:   &fxOutput,
	}
	inputTxID := ids.GenerateTestID()
	utxoID := dione.UTXOID{
		TxID:        inputTxID,
		OutputIndex: 0,
	}
	inputSigners := secp256k1fx.Input{
		SigIndices: []uint32{2},
	}
	fxInput := secp256k1fx.TransferInput{
		Amt:   54321,
		Input: inputSigners,
	}
	input := dione.TransferableInput{
		UTXOID: utxoID,
		Asset:  asset,
		In:     &fxInput,
	}
	baseTx := dione.BaseTx{
		NetworkID:    constants.UnitTestID,
		BlockchainID: ctx.ChainID,
		Ins: []*dione.TransferableInput{
			&input,
		},
	}
	tx := txs.ExportTx{
		BaseTx:           txs.BaseTx{BaseTx: baseTx},
		DestinationChain: ctx.DChainID,
		ExportedOuts: []*dione.TransferableOutput{
			&output,
		},
	}
	cred := fxs.FxCredential{
		Verifiable: &secp256k1fx.Credential{},
	}
	creds := []*fxs.FxCredential{
		&cred,
	}

	codec := parser.Codec()
	backend := &Backend{
		Ctx:    ctx,
		Config: &feeConfig,
		Fxs: []*fxs.ParsedFx{
			{
				ID: secp256k1fx.ID,
				Fx: fx,
			},
		},
		Codec:      codec,
		FeeAssetID: feeAssetID,
	}

	tests := []struct {
		name   string
		txFunc func() *txs.Tx
		err    error
	}{
		{
			name: "valid",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "no exported outputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.ExportedOuts = nil
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: errNoExportOutputs,
		},
		{
			name: "wrong networkID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.NetworkID++
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongNetworkID,
		},
		{
			name: "wrong chainID",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.BlockchainID = ids.GenerateTestID()
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrWrongChainID,
		},
		{
			name: "memo too large",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Memo = make([]byte, dione.MaxMemoSize+1)
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrMemoTooLarge,
		},
		{
			name: "invalid output",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          0,
					OutputOwners: outputOwners,
				}

				tx := tx
				tx.Outs = []*dione.TransferableOutput{
					&output,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueOutput,
		},
		{
			name: "unsorted outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "unsorted exported outputs",
			txFunc: func() *txs.Tx {
				output0 := output
				output0.Out = &secp256k1fx.TransferOutput{
					Amt:          1,
					OutputOwners: outputOwners,
				}

				output1 := output
				output1.Out = &secp256k1fx.TransferOutput{
					Amt:          2,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output0,
					&output1,
				}
				dione.SortTransferableOutputs(outputs, codec)
				outputs[0], outputs[1] = outputs[1], outputs[0]

				tx := tx
				tx.ExportedOuts = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrOutputsNotSorted,
		},
		{
			name: "invalid input",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   0,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: secp256k1fx.ErrNoValueInput,
		},
		{
			name: "duplicate inputs",
			txFunc: func() *txs.Tx {
				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: dione.ErrInputsNotSortedUnique,
		},
		{
			name: "input overflow",
			txFunc: func() *txs.Tx {
				input0 := input
				input0.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				input1 := input
				input1.UTXOID.OutputIndex++
				input1.In = &secp256k1fx.TransferInput{
					Amt:   stdmath.MaxUint64,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input0,
					&input1,
				}
				dione.SortTransferableInputsWithSigners(tx.Ins, make([][]*secp256k1.PrivateKey, 2))
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{
						&cred,
						&cred,
					},
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "output overflow",
			txFunc: func() *txs.Tx {
				output := output
				output.Out = &secp256k1fx.TransferOutput{
					Amt:          stdmath.MaxUint64,
					OutputOwners: outputOwners,
				}

				outputs := []*dione.TransferableOutput{
					&output,
				}
				dione.SortTransferableOutputs(outputs, codec)

				tx := tx
				tx.Outs = outputs
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: math.ErrOverflow,
		},
		{
			name: "insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
		{
			name: "invalid credential",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
					Creds: []*fxs.FxCredential{{
						Verifiable: (*secp256k1fx.Credential)(nil),
					}},
				}
			},
			err: secp256k1fx.ErrNilCredential,
		},
		{
			name: "wrong number of credentials",
			txFunc: func() *txs.Tx {
				return &txs.Tx{
					Unsigned: &tx,
				}
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "barely sufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: nil,
		},
		{
			name: "barely insufficient funds",
			txFunc: func() *txs.Tx {
				input := input
				input.In = &secp256k1fx.TransferInput{
					Amt:   fxOutput.Amt + feeConfig.TxFee - 1,
					Input: inputSigners,
				}

				tx := tx
				tx.Ins = []*dione.TransferableInput{
					&input,
				}
				return &txs.Tx{
					Unsigned: &tx,
					Creds:    creds,
				}
			},
			err: dione.ErrInsufficientFunds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := test.txFunc()
			verifier := &SyntacticVerifier{
				Backend: backend,
				Tx:      tx,
			}
			err := tx.Unsigned.Visit(verifier)
			require.ErrorIs(t, err, test.err)
		})
	}
}
