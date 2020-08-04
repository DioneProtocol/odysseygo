// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"bytes"
	"math"
	"testing"

	"github.com/ava-labs/gecko/chains/atomic"
	"github.com/ava-labs/gecko/database/memdb"
	"github.com/ava-labs/gecko/database/prefixdb"
	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/snow/engine/common"
	"github.com/ava-labs/gecko/utils/codec"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/utils/logging"
	"github.com/ava-labs/gecko/vms/components/ava"
	"github.com/ava-labs/gecko/vms/components/verify"
	"github.com/ava-labs/gecko/vms/secp256k1fx"
)

func TestExportTxSyntacticVerify(t *testing.T) {
	ctx := NewContext()
	c := setupCodec()

	tx := &ExportTx{
		BaseTx: BaseTx{
			NetID: networkID,
			BCID:  chainID,
			Ins: []*ava.TransferableInput{{
				UTXOID: ava.UTXOID{
					TxID: ids.NewID([32]byte{
						0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf9, 0xf8,
						0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1, 0xf0,
						0xef, 0xee, 0xed, 0xec, 0xeb, 0xea, 0xe9, 0xe8,
						0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
					}),
					OutputIndex: 0,
				},
				Asset: ava.Asset{ID: asset},
				In: &secp256k1fx.TransferInput{
					Amt: 54321,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{2},
					},
				},
			}},
		},
		Outs: []*ava.TransferableOutput{{
			Asset: ava.Asset{ID: asset},
			Out: &secp256k1fx.TransferOutput{
				Amt: 12345,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}
	tx.Initialize([]byte{})

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0); err != nil {
		t.Fatal(err)
	}
}

func TestExportTxSyntacticVerifyInvalidMemo(t *testing.T) {
	ctx := NewContext()
	c := setupCodec()

	tx := &ExportTx{
		BaseTx: BaseTx{
			NetID: networkID,
			BCID:  chainID,
			Ins: []*ava.TransferableInput{{
				UTXOID: ava.UTXOID{
					TxID: ids.NewID([32]byte{
						0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf9, 0xf8,
						0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1, 0xf0,
						0xef, 0xee, 0xed, 0xec, 0xeb, 0xea, 0xe9, 0xe8,
						0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
					}),
					OutputIndex: 0,
				},
				Asset: ava.Asset{ID: asset},
				In: &secp256k1fx.TransferInput{
					Amt: 54321,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{2},
					},
				},
			}},
			Memo: make([]byte, maxMemoSize+1),
		},
		Outs: []*ava.TransferableOutput{{
			Asset: ava.Asset{ID: asset},
			Out: &secp256k1fx.TransferOutput{
				Amt: 12345,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}
	tx.Initialize([]byte{})

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0); err == nil {
		t.Fatalf("should have errored due to memo field being too long")
	}
}

func TestExportTxSerialization(t *testing.T) {
	expected := []byte{
		// Codec version:
		0x00, 0x00,
		// txID:
		0x00, 0x00, 0x00, 0x04,
		// networkID:
		0x00, 0x00, 0x00, 0x02,
		// blockchainID:
		0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
		0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
		0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
		0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
		// number of outs:
		0x00, 0x00, 0x00, 0x00,
		// number of inputs:
		0x00, 0x00, 0x00, 0x01,
		// utxoID:
		0x0f, 0x2f, 0x4f, 0x6f, 0x8e, 0xae, 0xce, 0xee,
		0x0d, 0x2d, 0x4d, 0x6d, 0x8c, 0xac, 0xcc, 0xec,
		0x0b, 0x2b, 0x4b, 0x6b, 0x8a, 0xaa, 0xca, 0xea,
		0x09, 0x29, 0x49, 0x69, 0x88, 0xa8, 0xc8, 0xe8,
		// output index
		0x00, 0x00, 0x00, 0x00,
		// assetID:
		0x1f, 0x3f, 0x5f, 0x7f, 0x9e, 0xbe, 0xde, 0xfe,
		0x1d, 0x3d, 0x5d, 0x7d, 0x9c, 0xbc, 0xdc, 0xfc,
		0x1b, 0x3b, 0x5b, 0x7b, 0x9a, 0xba, 0xda, 0xfa,
		0x19, 0x39, 0x59, 0x79, 0x98, 0xb8, 0xd8, 0xf8,
		// input:
		// input ID:
		0x00, 0x00, 0x00, 0x05,
		// amount:
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8,
		// num sig indices:
		0x00, 0x00, 0x00, 0x01,
		// sig index[0]:
		0x00, 0x00, 0x00, 0x00,
		// Memo length:
		0x00, 0x00, 0x00, 0x04,
		// Memo:
		0x00, 0x01, 0x02, 0x03,
		// number of exported outs:
		0x00, 0x00, 0x00, 0x00,
	}

	tx := &Tx{UnsignedTx: &ExportTx{BaseTx: BaseTx{
		NetID: 2,
		BCID: ids.NewID([32]byte{
			0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
			0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
			0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
			0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
		}),
		Ins: []*ava.TransferableInput{{
			UTXOID: ava.UTXOID{TxID: ids.NewID([32]byte{
				0x0f, 0x2f, 0x4f, 0x6f, 0x8e, 0xae, 0xce, 0xee,
				0x0d, 0x2d, 0x4d, 0x6d, 0x8c, 0xac, 0xcc, 0xec,
				0x0b, 0x2b, 0x4b, 0x6b, 0x8a, 0xaa, 0xca, 0xea,
				0x09, 0x29, 0x49, 0x69, 0x88, 0xa8, 0xc8, 0xe8,
			})},
			Asset: ava.Asset{ID: ids.NewID([32]byte{
				0x1f, 0x3f, 0x5f, 0x7f, 0x9e, 0xbe, 0xde, 0xfe,
				0x1d, 0x3d, 0x5d, 0x7d, 0x9c, 0xbc, 0xdc, 0xfc,
				0x1b, 0x3b, 0x5b, 0x7b, 0x9a, 0xba, 0xda, 0xfa,
				0x19, 0x39, 0x59, 0x79, 0x98, 0xb8, 0xd8, 0xf8,
			})},
			In: &secp256k1fx.TransferInput{
				Amt:   1000,
				Input: secp256k1fx.Input{SigIndices: []uint32{0}},
			},
		}},
		Memo: []byte{0x00, 0x01, 0x02, 0x03},
	}}}

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

// Test issuing an import transaction.
func TestIssueExportTx(t *testing.T) {
	genesisBytes := BuildGenesisTest(t)

	issuer := make(chan common.Message, 1)
	baseDB := memdb.New()

	sm := &atomic.SharedMemory{}
	sm.Initialize(logging.NoLog{}, prefixdb.New([]byte{0}, baseDB))

	ctx := snow.DefaultContextTest()
	ctx.NetworkID = networkID
	ctx.ChainID = chainID
	ctx.SharedMemory = sm.NewBlockchainSharedMemory(chainID)

	genesisTx := GetFirstTxFromGenesisTest(genesisBytes, t)

	avaID := genesisTx.ID()
	platformID := ids.Empty.Prefix(0)

	ctx.Lock.Lock()
	vm := &VM{
		ava:      avaID,
		platform: platformID,
	}
	err := vm.Initialize(
		ctx,
		prefixdb.New([]byte{1}, baseDB),
		genesisBytes,
		issuer,
		[]*common.Fx{{
			ID: ids.Empty,
			Fx: &secp256k1fx.Fx{},
		}},
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	err = vm.Bootstrapping()
	if err != nil {
		t.Fatal(err)
	}

	err = vm.Bootstrapped()
	if err != nil {
		t.Fatal(err)
	}

	key := keys[0]

	tx := &Tx{UnsignedTx: &ExportTx{
		BaseTx: BaseTx{
			NetID: networkID,
			BCID:  chainID,
			Ins: []*ava.TransferableInput{{
				UTXOID: ava.UTXOID{
					TxID:        avaID,
					OutputIndex: 1,
				},
				Asset: ava.Asset{ID: avaID},
				In: &secp256k1fx.TransferInput{
					Amt:   50000,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		},
		Outs: []*ava.TransferableOutput{{
			Asset: ava.Asset{ID: avaID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 50000,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{key.PublicKey().Address()},
				},
			},
		}},
	}}

	unsignedBytes, err := vm.codec.Marshal(&tx.UnsignedTx)
	if err != nil {
		t.Fatal(err)
	}

	sig, err := key.Sign(unsignedBytes)
	if err != nil {
		t.Fatal(err)
	}
	fixedSig := [crypto.SECP256K1RSigLen]byte{}
	copy(fixedSig[:], sig)

	tx.Creds = append(tx.Creds, &secp256k1fx.Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			fixedSig,
		},
	})

	b, err := vm.codec.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	tx.Initialize(b)

	if _, err := vm.IssueTx(tx.Bytes(), nil); err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	msg := <-issuer
	if msg != common.PendingTxs {
		t.Fatalf("Wrong message")
	}

	ctx.Lock.Lock()
	defer func() {
		vm.Shutdown()
		ctx.Lock.Unlock()
	}()

	txs := vm.PendingTxs()
	if len(txs) != 1 {
		t.Fatalf("Should have returned %d tx(s)", 1)
	}

	parsedTx := txs[0]
	if err := parsedTx.Verify(); err != nil {
		t.Fatal(err)
	}
	parsedTx.Accept()

	smDB := vm.ctx.SharedMemory.GetDatabase(platformID)
	defer vm.ctx.SharedMemory.ReleaseDatabase(platformID)

	state := ava.NewPrefixedState(smDB, vm.codec)

	utxo := ava.UTXOID{
		TxID:        tx.ID(),
		OutputIndex: 0,
	}
	utxoID := utxo.InputID()
	if _, err := state.AVMUTXO(utxoID); err != nil {
		t.Fatal(err)
	}

	utxoIDs, err := state.AVMFunds(key.PublicKey().Address().LongID(), ids.Empty, math.MaxInt32)
	if err != nil {
		t.Fatal(err)
	}
	if len(utxoIDs) != 1 {
		t.Fatalf("wrong number of utxoIDs %d", len(utxoIDs))
	}
}

// Test force accepting an import transaction.
func TestClearForceAcceptedExportTx(t *testing.T) {
	genesisBytes := BuildGenesisTest(t)

	issuer := make(chan common.Message, 1)
	baseDB := memdb.New()

	sm := &atomic.SharedMemory{}
	sm.Initialize(logging.NoLog{}, prefixdb.New([]byte{0}, baseDB))

	ctx := snow.DefaultContextTest()
	ctx.NetworkID = networkID
	ctx.ChainID = chainID
	ctx.SharedMemory = sm.NewBlockchainSharedMemory(chainID)

	genesisTx := GetFirstTxFromGenesisTest(genesisBytes, t)

	avaID := genesisTx.ID()
	platformID := ids.Empty.Prefix(0)

	ctx.Lock.Lock()
	vm := &VM{
		ava:      avaID,
		platform: platformID,
	}
	err := vm.Initialize(
		ctx,
		prefixdb.New([]byte{1}, baseDB),
		genesisBytes,
		issuer,
		[]*common.Fx{{
			ID: ids.Empty,
			Fx: &secp256k1fx.Fx{},
		}},
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	err = vm.Bootstrapping()
	if err != nil {
		t.Fatal(err)
	}

	err = vm.Bootstrapped()
	if err != nil {
		t.Fatal(err)
	}

	key := keys[0]

	tx := &Tx{UnsignedTx: &ExportTx{
		BaseTx: BaseTx{
			NetID: networkID,
			BCID:  chainID,
			Ins: []*ava.TransferableInput{{
				UTXOID: ava.UTXOID{
					TxID:        avaID,
					OutputIndex: 1,
				},
				Asset: ava.Asset{ID: avaID},
				In: &secp256k1fx.TransferInput{
					Amt:   50000,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		},
		Outs: []*ava.TransferableOutput{{
			Asset: ava.Asset{ID: avaID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 50000,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{key.PublicKey().Address()},
				},
			},
		}},
	}}

	unsignedBytes, err := vm.codec.Marshal(&tx.UnsignedTx)
	if err != nil {
		t.Fatal(err)
	}

	sig, err := key.Sign(unsignedBytes)
	if err != nil {
		t.Fatal(err)
	}
	fixedSig := [crypto.SECP256K1RSigLen]byte{}
	copy(fixedSig[:], sig)

	tx.Creds = append(tx.Creds, &secp256k1fx.Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			fixedSig,
		},
	})

	b, err := vm.codec.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	tx.Initialize(b)

	if _, err := vm.IssueTx(tx.Bytes(), nil); err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	msg := <-issuer
	if msg != common.PendingTxs {
		t.Fatalf("Wrong message")
	}

	ctx.Lock.Lock()
	defer func() {
		vm.Shutdown()
		ctx.Lock.Unlock()
	}()

	txs := vm.PendingTxs()
	if len(txs) != 1 {
		t.Fatalf("Should have returned %d tx(s)", 1)
	}

	parsedTx := txs[0]
	if err := parsedTx.Verify(); err != nil {
		t.Fatal(err)
	}

	smDB := vm.ctx.SharedMemory.GetDatabase(platformID)

	state := ava.NewPrefixedState(smDB, vm.codec)

	utxo := ava.UTXOID{
		TxID:        tx.ID(),
		OutputIndex: 0,
	}
	utxoID := utxo.InputID()
	if err := state.SpendAVMUTXO(utxoID); err != nil {
		t.Fatal(err)
	}

	vm.ctx.SharedMemory.ReleaseDatabase(platformID)

	parsedTx.Accept()

	smDB = vm.ctx.SharedMemory.GetDatabase(platformID)
	defer vm.ctx.SharedMemory.ReleaseDatabase(platformID)

	state = ava.NewPrefixedState(smDB, vm.codec)

	if _, err := state.AVMUTXO(utxoID); err == nil {
		t.Fatalf("should have failed to read the utxo")
	}
}

func TestExportTxNotState(t *testing.T) {
	intf := interface{}(&ExportTx{})
	if _, ok := intf.(verify.State); ok {
		t.Fatalf("shouldn't be marked as state")
	}
}
