// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"bytes"
	"testing"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/codec"
	"github.com/ava-labs/gecko/vms/components/ava"
	"github.com/ava-labs/gecko/vms/components/verify"
	"github.com/ava-labs/gecko/vms/secp256k1fx"
)

func TestCreateAssetTxSerialization(t *testing.T) {
	expected := []byte{
		// Codec version:
		0x00, 0x00,
		// txID:
		0x00, 0x00, 0x00, 0x01,
		// networkID:
		0x00, 0x00, 0x00, 0x02,
		// blockchainID:
		0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
		0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
		0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
		0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
		// number of outs:
		0x00, 0x00, 0x00, 0x01,
		// output[0]:
		// assetID:
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		// output:
		0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
		// number of inputs:
		0x00, 0x00, 0x00, 0x01,
		// txID:
		0xf1, 0xe1, 0xd1, 0xc1, 0xb1, 0xa1, 0x91, 0x81,
		0x71, 0x61, 0x51, 0x41, 0x31, 0x21, 0x11, 0x01,
		0xf0, 0xe0, 0xd0, 0xc0, 0xb0, 0xa0, 0x90, 0x80,
		0x70, 0x60, 0x50, 0x40, 0x30, 0x20, 0x10, 0x00,
		// utxoIndex:
		0x00, 0x00, 0x00, 0x05,
		// assetID:
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		// input:
		0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00,
		0x07, 0x5b, 0xcd, 0x15, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x07,
		// Memo length:
		0x00, 0x00, 0x00, 0x04,
		// Memo:
		0x00, 0x01, 0x02, 0x03,
		// name:
		0x00, 0x10, 0x56, 0x6f, 0x6c, 0x61, 0x74, 0x69,
		0x6c, 0x69, 0x74, 0x79, 0x20, 0x49, 0x6e, 0x64,
		0x65, 0x78,
		// symbol:
		0x00, 0x03, 0x56, 0x49, 0x58,
		// denomination:
		0x02,
		// number of InitialStates:
		0x00, 0x00, 0x00, 0x01,
		// InitialStates[0]:
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
	}

	tx := &Tx{UnsignedTx: &CreateAssetTx{
		BaseTx: BaseTx{
			NetID: 2,
			BCID: ids.NewID([32]byte{
				0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
				0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
				0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
				0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
			}),
			Memo: []byte{0x00, 0x01, 0x02, 0x03},
			Outs: []*ava.TransferableOutput{{
				Asset: ava.Asset{
					ID: ids.NewID([32]byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
						0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					}),
				},
				Out: &secp256k1fx.TransferOutput{
					Amt: 12345,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  54321,
						Threshold: 1,
						Addrs: []ids.ShortID{
							ids.NewShortID([20]byte{
								0x51, 0x02, 0x5c, 0x61, 0xfb, 0xcf, 0xc0, 0x78,
								0xf6, 0x93, 0x34, 0xf8, 0x34, 0xbe, 0x6d, 0xd2,
								0x6d, 0x55, 0xa9, 0x55,
							}),
							ids.NewShortID([20]byte{
								0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
								0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
								0x43, 0xab, 0x08, 0x59,
							}),
						},
					},
				},
			}},
			Ins: []*ava.TransferableInput{{
				UTXOID: ava.UTXOID{
					TxID: ids.NewID([32]byte{
						0xf1, 0xe1, 0xd1, 0xc1, 0xb1, 0xa1, 0x91, 0x81,
						0x71, 0x61, 0x51, 0x41, 0x31, 0x21, 0x11, 0x01,
						0xf0, 0xe0, 0xd0, 0xc0, 0xb0, 0xa0, 0x90, 0x80,
						0x70, 0x60, 0x50, 0x40, 0x30, 0x20, 0x10, 0x00,
					}),
					OutputIndex: 5,
				},
				Asset: ava.Asset{
					ID: ids.NewID([32]byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
						0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					}),
				},
				In: &secp256k1fx.TransferInput{
					Amt: 123456789,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{3, 7},
					},
				},
			}},
		},
		Name:         "Volatility Index",
		Symbol:       "VIX",
		Denomination: 2,
		States: []*InitialState{
			{
				FxID: 0,
				Outs: []verify.State{
					&secp256k1fx.TransferOutput{
						Amt: 12345,
						OutputOwners: secp256k1fx.OutputOwners{
							Locktime:  54321,
							Threshold: 1,
							Addrs: []ids.ShortID{
								ids.NewShortID([20]byte{
									0x51, 0x02, 0x5c, 0x61, 0xfb, 0xcf, 0xc0, 0x78,
									0xf6, 0x93, 0x34, 0xf8, 0x34, 0xbe, 0x6d, 0xd2,
									0x6d, 0x55, 0xa9, 0x55,
								}),
								ids.NewShortID([20]byte{
									0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
									0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
									0x43, 0xab, 0x08, 0x59,
								}),
							},
						},
					},
				},
			},
		},
	}}

	c := codec.NewDefault()
	c.RegisterType(&BaseTx{})
	c.RegisterType(&CreateAssetTx{})
	c.RegisterType(&OperationTx{})
	c.RegisterType(&ImportTx{})
	c.RegisterType(&ExportTx{})
	c.RegisterType(&secp256k1fx.TransferInput{})
	c.RegisterType(&secp256k1fx.MintOutput{})
	c.RegisterType(&secp256k1fx.TransferOutput{})
	c.RegisterType(&secp256k1fx.MintOperation{})
	c.RegisterType(&secp256k1fx.Credential{})

	b, err := c.Marshal(&tx.UnsignedTx)
	if err != nil {
		t.Fatal(err)
	}
	tx.Initialize(b)

	result := tx.Bytes()
	if !bytes.Equal(expected, result) {
		t.Fatalf("\nExpected: 0x%x\nResult:   0x%x", expected, result)
	}
}

func TestCreateAssetTxNotState(t *testing.T) {
	intf := interface{}(&CreateAssetTx{})
	if _, ok := intf.(verify.State); ok {
		t.Fatalf("shouldn't be marked as state")
	}
}
