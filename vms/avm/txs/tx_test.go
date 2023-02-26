// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"testing"

	"github.com/dioneprotocol/dionego/codec"
	"github.com/dioneprotocol/dionego/codec/linearcodec"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/utils/crypto/secp256k1"
	"github.com/dioneprotocol/dionego/utils/units"
	"github.com/dioneprotocol/dionego/utils/wrappers"
	"github.com/dioneprotocol/dionego/vms/avm/fxs"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

var (
	networkID       uint32 = 10
	chainID                = ids.ID{5, 4, 3, 2, 1}
	platformChainID        = ids.Empty.Prefix(0)

	keys = secp256k1.TestKeys()

	assetID = ids.ID{1, 2, 3}
)

func setupCodec() codec.Manager {
	parser, err := NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
	})
	if err != nil {
		panic(err)
	}
	return parser.Codec()
}

func NewContext(tb testing.TB) *snow.Context {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = networkID
	ctx.ChainID = chainID
	dioneAssetID, err := ids.FromString("2XGxUr7VF7j1iwUp2aiGe4b6Ue2yyNghNS1SuNTNmZ77dPpXFZ")
	if err != nil {
		tb.Fatal(err)
	}
	ctx.DIONEAssetID = dioneAssetID
	ctx.XChainID = ids.Empty.Prefix(0)
	ctx.CChainID = ids.Empty.Prefix(1)
	aliaser := ctx.BCLookup.(ids.Aliaser)

	errs := wrappers.Errs{}
	errs.Add(
		aliaser.Alias(chainID, "X"),
		aliaser.Alias(chainID, chainID.String()),
		aliaser.Alias(platformChainID, "P"),
		aliaser.Alias(platformChainID, platformChainID.String()),
	)
	if errs.Errored() {
		tb.Fatal(errs.Err)
	}
	return ctx
}

func TestTxNil(t *testing.T) {
	ctx := NewContext(t)
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	tx := (*Tx)(nil)
	if err := tx.SyntacticVerify(ctx, m, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Should have erred due to nil tx")
	}
}

func TestTxEmpty(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()
	tx := &Tx{}
	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Should have erred due to nil tx")
	}
}

func TestTxInvalidCredential(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
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
		}},
		Creds: []*fxs.FxCredential{{Verifiable: &dione.TestVerifiable{Err: errTest}}},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid credential")
	}
}

func TestTxInvalidUnsignedTx(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{
				{
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
				},
				{
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
				},
			},
		}},
		Creds: []*fxs.FxCredential{
			{Verifiable: &dione.TestVerifiable{}},
			{Verifiable: &dione.TestVerifiable{}},
		},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid unsigned tx")
	}
}

func TestTxInvalidNumberOfCredentials(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{
				{
					UTXOID: dione.UTXOID{TxID: ids.Empty, OutputIndex: 0},
					Asset:  dione.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloDione,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
				{
					UTXOID: dione.UTXOID{TxID: ids.Empty, OutputIndex: 1},
					Asset:  dione.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloDione,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
			},
		}},
		Creds: []*fxs.FxCredential{{Verifiable: &dione.TestVerifiable{}}},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid number of credentials")
	}
}
