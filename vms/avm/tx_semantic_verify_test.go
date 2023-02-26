// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/dioneprotocol/dionego/api/keystore"
	"github.com/dioneprotocol/dionego/chains/atomic"
	"github.com/dioneprotocol/dionego/database/manager"
	"github.com/dioneprotocol/dionego/database/prefixdb"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/snow"
	"github.com/dioneprotocol/dionego/snow/engine/common"
	"github.com/dioneprotocol/dionego/utils/constants"
	"github.com/dioneprotocol/dionego/utils/crypto/secp256k1"
	"github.com/dioneprotocol/dionego/version"
	"github.com/dioneprotocol/dionego/vms/avm/fxs"
	"github.com/dioneprotocol/dionego/vms/avm/txs"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

func TestBaseTxSemanticVerify(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
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
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBaseTxSemanticVerifyUnknownFx(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVMWithArgs(
		t,
		[]*common.Fx{{
			ID: ids.GenerateTestID(),
			Fx: &FxTest{
				InitializeF: func(vmIntf interface{}) error {
					vm := vmIntf.(secp256k1fx.VM)
					return vm.CodecRegistry().RegisterType(&dione.TestVerifiable{})
				},
			},
		}},
		nil,
	)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{
		Unsigned: &txs.BaseTx{
			BaseTx: dione.BaseTx{
				NetworkID:    networkID,
				BlockchainID: chainID,
				Ins: []*dione.TransferableInput{
					{
						UTXOID: dione.UTXOID{
							TxID:        genesisTx.ID(),
							OutputIndex: 1,
						},
						Asset: dione.Asset{ID: genesisTx.ID()},
						In: &secp256k1fx.TransferInput{
							Amt: startBalance,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{
									0,
								},
							},
						},
					},
				},
			},
		},
		Creds: []*fxs.FxCredential{{
			Verifiable: &dione.TestVerifiable{},
		}},
	}
	if err := vm.parser.InitializeTx(tx); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an unknown feature extension")
	}
}

func TestBaseTxSemanticVerifyWrongAssetID(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: assetID},
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
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an asset ID mismatch")
	}
}

func TestBaseTxSemanticVerifyUnauthorizedFx(t *testing.T) {
	ctx := NewContext(t)
	vm := &VM{}
	ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	fx := &FxTest{}
	fx.InitializeF = func(vmIntf interface{}) error {
		vm := vmIntf.(secp256k1fx.VM)
		return vm.CodecRegistry().RegisterType(&dione.TestTransferable{})
	}

	genesisBytes := BuildGenesisTest(t)
	issuer := make(chan common.Message, 1)
	err := vm.Initialize(
		context.Background(),
		ctx,
		manager.NewMemDB(version.Semantic1_0_0),
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*common.Fx{
			{
				ID: ids.Empty,
				Fx: &secp256k1fx.Fx{},
			},
			{
				ID: ids.ID{1},
				Fx: fx,
			},
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	if err := vm.SetState(context.Background(), snow.Bootstrapping); err != nil {
		t.Fatal(err)
	}

	err = vm.SetState(context.Background(), snow.NormalOp)
	if err != nil {
		t.Fatal(err)
	}

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In:    &dione.TestTransferable{},
			}},
		},
	}}
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err = tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an unsupported fx")
	}
}

func TestBaseTxSemanticVerifyInvalidSignature(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{
		Unsigned: &txs.BaseTx{
			BaseTx: dione.BaseTx{
				NetworkID:    networkID,
				BlockchainID: chainID,
				Ins: []*dione.TransferableInput{
					{
						UTXOID: dione.UTXOID{
							TxID:        genesisTx.ID(),
							OutputIndex: 2,
						},
						Asset: dione.Asset{ID: genesisTx.ID()},
						In: &secp256k1fx.TransferInput{
							Amt: startBalance,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{
									0,
								},
							},
						},
					},
				},
			},
		},
		Creds: []*fxs.FxCredential{
			{
				Verifiable: &secp256k1fx.Credential{
					Sigs: [][secp256k1.SignatureLen]byte{{}},
				},
			},
		},
	}
	if err := vm.parser.InitializeTx(tx); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Invalid credential should have failed verification")
	}
}

func TestBaseTxSemanticVerifyMissingUTXO(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        ids.Empty,
					OutputIndex: 1,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
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
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Unknown UTXO should have failed verification")
	}
}

func TestBaseTxSemanticVerifyInvalidUTXO(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: math.MaxUint32,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
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
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Invalid UTXO should have failed verification")
	}
}

func TestBaseTxSemanticVerifyPendingInvalidUTXO(t *testing.T) {
	genesisBytes, issuer, vm, _ := GenesisVM(t)
	ctx := vm.ctx

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	pendingTx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: genesisTx.ID()},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance - vm.TxFee,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
		},
	}}
	if err := pendingTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	txID, err := vm.IssueTx(pendingTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	<-issuer

	ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	_ = vm.PendingTxs(context.Background())

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        txID,
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
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
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err = tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Invalid UTXO should have failed verification")
	}
}

func TestBaseTxSemanticVerifyPendingWrongAssetID(t *testing.T) {
	genesisBytes, issuer, vm, _ := GenesisVM(t)
	ctx := vm.ctx

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	pendingTx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: genesisTx.ID()},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance - vm.TxFee,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
		},
	}}
	if err := pendingTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	txID, err := vm.IssueTx(pendingTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	<-issuer

	ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	_ = vm.PendingTxs(context.Background())

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        txID,
					OutputIndex: 0,
				},
				Asset: dione.Asset{ID: assetID},
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

	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err = tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Wrong asset ID should have failed verification")
	}
}

func TestBaseTxSemanticVerifyPendingUnauthorizedFx(t *testing.T) {
	genesisBytes := BuildGenesisTest(t)
	ctx := NewContext(t)

	issuer := make(chan common.Message, 1)

	ctx.Lock.Lock()

	vm := &VM{}

	fx := &FxTest{}
	fx.InitializeF = func(vmIntf interface{}) error {
		vm := vmIntf.(secp256k1fx.VM)
		return vm.CodecRegistry().RegisterType(&dione.TestVerifiable{})
	}

	err := vm.Initialize(
		context.Background(),
		ctx,
		manager.NewMemDB(version.Semantic1_0_0),
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*common.Fx{
			{
				ID: ids.ID{1},
				Fx: &secp256k1fx.Fx{},
			},
			{
				ID: ids.Empty,
				Fx: fx,
			},
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	if err := vm.SetState(context.Background(), snow.Bootstrapping); err != nil {
		t.Fatal(err)
	}

	err = vm.SetState(context.Background(), snow.NormalOp)
	if err != nil {
		t.Fatal(err)
	}

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	pendingTx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: genesisTx.ID()},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance - vm.TxFee,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
		},
	}}
	if err := pendingTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	txID, err := vm.IssueTx(pendingTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	<-issuer

	ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	_ = vm.PendingTxs(context.Background())

	tx := &txs.Tx{
		Unsigned: &txs.BaseTx{
			BaseTx: dione.BaseTx{
				NetworkID:    networkID,
				BlockchainID: chainID,
				Ins: []*dione.TransferableInput{
					{
						UTXOID: dione.UTXOID{
							TxID:        txID,
							OutputIndex: 0,
						},
						Asset: dione.Asset{ID: genesisTx.ID()},
						In: &secp256k1fx.TransferInput{
							Amt: startBalance,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{
									0,
								},
							},
						},
					},
				},
			},
		},
		Creds: []*fxs.FxCredential{{
			Verifiable: &dione.TestVerifiable{},
		}},
	}
	if err := vm.parser.InitializeTx(tx); err != nil {
		t.Fatal(err)
	}

	err = tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Unsupported feature extension should have failed verification")
	}
}

func TestBaseTxSemanticVerifyPendingInvalidSignature(t *testing.T) {
	genesisBytes := BuildGenesisTest(t)
	ctx := NewContext(t)

	issuer := make(chan common.Message, 1)

	ctx.Lock.Lock()

	vm := &VM{}

	fx := &FxTest{}
	fx.InitializeF = func(vmIntf interface{}) error {
		vm := vmIntf.(secp256k1fx.VM)
		return vm.CodecRegistry().RegisterType(&dione.TestVerifiable{})
	}

	err := vm.Initialize(
		context.Background(),
		ctx,
		manager.NewMemDB(version.Semantic1_0_0),
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*common.Fx{
			{
				ID: ids.ID{1},
				Fx: &secp256k1fx.Fx{},
			},
			{
				ID: ids.Empty,
				Fx: fx,
			},
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	if err := vm.SetState(context.Background(), snow.Bootstrapping); err != nil {
		t.Fatal(err)
	}

	err = vm.SetState(context.Background(), snow.NormalOp)
	if err != nil {
		t.Fatal(err)
	}

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	pendingTx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: genesisTx.ID()},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance - vm.TxFee,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
		},
	}}
	if err := pendingTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	txID, err := vm.IssueTx(pendingTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	ctx.Lock.Unlock()

	<-issuer

	ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	_ = vm.PendingTxs(context.Background())

	tx := &txs.Tx{
		Unsigned: &txs.BaseTx{
			BaseTx: dione.BaseTx{
				NetworkID:    networkID,
				BlockchainID: chainID,
				Ins: []*dione.TransferableInput{
					{
						UTXOID: dione.UTXOID{
							TxID:        txID,
							OutputIndex: 0,
						},
						Asset: dione.Asset{ID: genesisTx.ID()},
						In: &secp256k1fx.TransferInput{
							Amt: startBalance,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{
									0,
								},
							},
						},
					},
				},
			},
		},
		Creds: []*fxs.FxCredential{{
			Verifiable: &secp256k1fx.Credential{
				Sigs: [][secp256k1.SignatureLen]byte{{}},
			},
		}},
	}

	if err := vm.parser.InitializeTx(tx); err != nil {
		t.Fatal(err)
	}

	err = tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("Invalid signature should have failed verification")
	}
}

func TestBaseTxSemanticVerifyMalformedOutput(t *testing.T) {
	_, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	txBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x66,
		0x05, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x70, 0xae, 0x33, 0xb5,
		0x60, 0x9c, 0xd8, 0x9a, 0x72, 0x92, 0x4f, 0xa2,
		0x88, 0x3f, 0x9b, 0xf1, 0xc6, 0xd8, 0x9f, 0x07,
		0x09, 0x9b, 0x2a, 0xd7, 0x1b, 0xe1, 0x7c, 0x5d,
		0x44, 0x93, 0x23, 0xdb, 0x00, 0x00, 0x00, 0x05,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc3, 0x50,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x70, 0xae, 0x33, 0xb5, 0x60, 0x9c, 0xd8, 0x9a,
		0x72, 0x92, 0x4f, 0xa2, 0x88, 0x3f, 0x9b, 0xf1,
		0xc6, 0xd8, 0x9f, 0x07, 0x09, 0x9b, 0x2a, 0xd7,
		0x1b, 0xe1, 0x7c, 0x5d, 0x44, 0x93, 0x23, 0xdb,
		0x00, 0x00, 0x00, 0x01, 0x70, 0xae, 0x33, 0xb5,
		0x60, 0x9c, 0xd8, 0x9a, 0x72, 0x92, 0x4f, 0xa2,
		0x88, 0x3f, 0x9b, 0xf1, 0xc6, 0xd8, 0x9f, 0x07,
		0x09, 0x9b, 0x2a, 0xd7, 0x1b, 0xe1, 0x7c, 0x5d,
		0x44, 0x93, 0x23, 0xdb, 0x00, 0x00, 0x00, 0x05,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc3, 0x50,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x01, 0x50, 0x6b, 0xd9, 0x2d,
		0xe5, 0xeb, 0xc2, 0xbf, 0x8f, 0xaa, 0xf1, 0x7d,
		0xbb, 0xae, 0xb3, 0xf3, 0x13, 0x9e, 0xae, 0xb4,
		0xad, 0x32, 0x95, 0x6e, 0x92, 0x74, 0xf9, 0x53,
		0x0e, 0xcc, 0x03, 0xd8, 0x02, 0xab, 0x1c, 0x16,
		0x52, 0xd0, 0xe3, 0xfc, 0xe5, 0x93, 0xa9, 0x8e,
		0x96, 0x1e, 0x83, 0xf0, 0x12, 0x27, 0x66, 0x9f,
		0x03, 0x56, 0x9f, 0x17, 0x1b, 0xd1, 0x22, 0x90,
		0xfd, 0x64, 0xf5, 0x73, 0x01,
	}

	tx := &txs.Tx{}
	if _, err := vm.parser.Codec().Unmarshal(txBytes, tx); err == nil {
		t.Fatalf("should have failed to unmarshal the tx")
	}
}

func TestBaseTxSemanticVerifyInvalidFxOutput(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVMWithArgs(
		t,
		[]*common.Fx{{
			ID: ids.GenerateTestID(),
			Fx: &FxTest{
				InitializeF: func(vmIntf interface{}) error {
					vm := vmIntf.(secp256k1fx.VM)
					return vm.CodecRegistry().RegisterType(&dione.TestTransferable{})
				},
			},
		}},
		nil,
	)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	tx := &txs.Tx{Unsigned: &txs.BaseTx{
		BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        genesisTx.ID(),
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: genesisTx.ID()},
				In: &secp256k1fx.TransferInput{
					Amt: startBalance,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: genesisTx.ID()},
				Out: &dione.TestTransferable{
					Val: 1,
				},
			}},
		},
	}}
	if err := tx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	err := tx.Unsigned.Visit(&txSemanticVerify{
		tx: tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to sending funds to an un-authorized fx")
	}
}

func TestExportTxSemanticVerify(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisTx.ID()
	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        dioneID,
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: dioneID},
				In: &secp256k1fx.TransferInput{
					Amt:   startBalance,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: dioneID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}

	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportTxSemanticVerifyUnknownCredFx(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisTx.ID()
	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        dioneID,
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: dioneID},
				In: &secp256k1fx.TransferInput{
					Amt:   startBalance,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: dioneID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}
	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	utx.Tx.Creds[0].Verifiable = nil
	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an unknown credential fx")
	}
}

func TestExportTxSemanticVerifyMissingUTXO(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisTx.ID()
	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        dioneID,
					OutputIndex: 1000,
				},
				Asset: dione.Asset{ID: dioneID},
				In: &secp256k1fx.TransferInput{
					Amt:   startBalance,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: dioneID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}

	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an unknown utxo")
	}
}

// Test that we can't create an output of by consuming a UTXO that doesn't exist
func TestExportTxSemanticVerifyInvalidAssetID(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisTx.ID()
	assetID := dioneID
	// so the inputs below are sorted
	copy(assetID[len(assetID)-5:], []byte{255, 255, 255, 255})
	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{
				{
					UTXOID: dione.UTXOID{
						TxID:        dioneID,
						OutputIndex: 0,
					},
					Asset: dione.Asset{ID: vm.ctx.DIONEAssetID},
					In: &secp256k1fx.TransferInput{
						Amt:   startBalance,
						Input: secp256k1fx.Input{SigIndices: []uint32{0}},
					},
				},
				{
					UTXOID: dione.UTXOID{
						TxID:        assetID, // This tx doesn't exist
						OutputIndex: 0,
					},
					Asset: dione.Asset{ID: assetID}, // This asset doesn't exist
					In: &secp256k1fx.TransferInput{
						Amt:   startBalance,
						Input: secp256k1fx.Input{SigIndices: []uint32{0}},
					},
				},
			},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: assetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}
	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{
		{
			keys[0],
		},
		{
			keys[0],
		},
	}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an invalid asset ID")
	}
}

func TestExportTxSemanticVerifyInvalidFx(t *testing.T) {
	genesisBytes := BuildGenesisTest(t)
	ctx := NewContext(t)

	baseDBManager := manager.NewMemDB(version.Semantic1_0_0)

	m := atomic.NewMemory(prefixdb.New([]byte{0}, baseDBManager.Current().Database))
	ctx.SharedMemory = m.NewSharedMemory(ctx.ChainID)

	ctx.Lock.Lock()

	userKeystore, err := keystore.CreateTestKeystore()
	if err != nil {
		t.Fatal(err)
	}
	if err := userKeystore.CreateUser(username, password); err != nil {
		t.Fatal(err)
	}
	ctx.Keystore = userKeystore.NewBlockchainKeyStore(ctx.ChainID)

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)

	dioneID := genesisTx.ID()

	issuer := make(chan common.Message, 1)
	vm := &VM{}
	err = vm.Initialize(
		context.Background(),
		ctx,
		baseDBManager.NewPrefixDBManager([]byte{1}),
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*common.Fx{
			{
				ID: ids.Empty,
				Fx: &secp256k1fx.Fx{},
			},
			{
				ID: ids.Empty.Prefix(0),
				Fx: &FxTest{
					InitializeF: func(vmIntf interface{}) error {
						vm := vmIntf.(secp256k1fx.VM)
						return vm.CodecRegistry().RegisterType(&dione.TestVerifiable{})
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	vm.batchTimeout = 0

	if err := vm.SetState(context.Background(), snow.Bootstrapping); err != nil {
		t.Fatal(err)
	}

	if err := vm.SetState(context.Background(), snow.NormalOp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        dioneID,
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: dioneID},
				In: &secp256k1fx.TransferInput{
					Amt:   startBalance,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: dioneID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}
	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[0]}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	utx.Tx.Creds[0].Verifiable = &dione.TestVerifiable{}
	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to using an invalid fxID")
	}
}

func TestExportTxSemanticVerifyInvalidTransfer(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	genesisTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisTx.ID()
	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        dioneID,
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: dioneID},
				In: &secp256k1fx.TransferInput{
					Amt:   startBalance,
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{{
			Asset: dione.Asset{ID: dioneID},
			Out: &secp256k1fx.TransferOutput{
				Amt: startBalance - vm.TxFee,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
				},
			},
		}},
	}}
	if err := rawTx.SignSECP256K1Fx(vm.parser.Codec(), [][]*secp256k1.PrivateKey{{keys[1]}}); err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err == nil {
		t.Fatalf("should have erred due to an invalid credential")
	}
}

func TestExportTxSemanticVerifyTransferCustomAsset(t *testing.T) {
	genesisBytes, _, vm, _ := GenesisVM(t)
	ctx := vm.ctx
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		ctx.Lock.Unlock()
	}()

	vm.clock.Set(testBanffTime.Add(time.Second))

	genesisDioneTx := GetDIONETxFromGenesisTest(genesisBytes, t)
	dioneID := genesisDioneTx.ID()

	genesisCustomAssetTx := GetCreateTxFromGenesisTest(t, genesisBytes, "myFixedCapAsset")
	customAssetID := genesisCustomAssetTx.ID()

	rawTx := &txs.Tx{Unsigned: &txs.ExportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*dione.TransferableInput{
				{
					UTXOID: dione.UTXOID{
						TxID:        customAssetID,
						OutputIndex: 1,
					},
					Asset: dione.Asset{ID: customAssetID},
					In: &secp256k1fx.TransferInput{
						Amt:   startBalance,
						Input: secp256k1fx.Input{SigIndices: []uint32{0}},
					},
				},
				{
					UTXOID: dione.UTXOID{
						TxID:        dioneID,
						OutputIndex: 2,
					},
					Asset: dione.Asset{ID: dioneID},
					In: &secp256k1fx.TransferInput{
						Amt:   startBalance,
						Input: secp256k1fx.Input{SigIndices: []uint32{0}},
					},
				},
			},
		}},
		DestinationChain: constants.PlatformChainID,
		ExportedOuts: []*dione.TransferableOutput{
			{
				Asset: dione.Asset{ID: customAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			},
			{
				Asset: dione.Asset{ID: dioneID},
				Out: &secp256k1fx.TransferOutput{
					Amt: startBalance - vm.TxFee,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			},
		},
	}}

	err := rawTx.SignSECP256K1Fx(
		vm.parser.Codec(),
		[][]*secp256k1.PrivateKey{
			{keys[0]},
			{keys[0]},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := vm.ParseTx(context.Background(), rawTx.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	utx, ok := tx.(*UniqueTx)
	if !ok {
		t.Fatalf("wrong tx type")
	}

	err = rawTx.Unsigned.Visit(&txSemanticVerify{
		tx: utx.Tx,
		vm: vm,
	})
	if err != nil {
		t.Fatal(err)
	}
}
