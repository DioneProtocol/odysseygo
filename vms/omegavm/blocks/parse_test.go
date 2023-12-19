// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/codec"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var preFundedKeys = secp256k1.TestKeys()

func TestStandardBlocks(t *testing.T) {
	// check Odyssey standard block can be built and parsed
	require := require.New(t)
	blkTimestamp := time.Now()
	parentID := ids.ID{'p', 'a', 'r', 'e', 'n', 't', 'I', 'D'}
	height := uint64(2022)
	txs, err := testDecisionTxs()
	require.NoError(err)

	for _, cdc := range []codec.Manager{Codec, GenesisCodec} {
		// build block
		odysseyStandardBlk, err := NewOdysseyStandardBlock(parentID, height, txs)
		require.NoError(err)

		// parse block
		parsed, err := Parse(cdc, odysseyStandardBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(odysseyStandardBlk.ID(), parsed.ID())
		require.Equal(odysseyStandardBlk.Bytes(), parsed.Bytes())
		require.Equal(odysseyStandardBlk.Parent(), parsed.Parent())
		require.Equal(odysseyStandardBlk.Height(), parsed.Height())

		require.IsType(&OdysseyStandardBlock{}, parsed)
		require.Equal(txs, parsed.Txs())

		// check that banff standard block can be built and parsed
		banffStandardBlk, err := NewBanffStandardBlock(blkTimestamp, parentID, height, txs)
		require.NoError(err)

		// parse block
		parsed, err = Parse(cdc, banffStandardBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(banffStandardBlk.ID(), parsed.ID())
		require.Equal(banffStandardBlk.Bytes(), parsed.Bytes())
		require.Equal(banffStandardBlk.Parent(), parsed.Parent())
		require.Equal(banffStandardBlk.Height(), parsed.Height())
		require.IsType(&BanffStandardBlock{}, parsed)
		parsedBanffStandardBlk := parsed.(*BanffStandardBlock)
		require.Equal(txs, parsedBanffStandardBlk.Txs())

		// timestamp check for banff blocks only
		require.Equal(banffStandardBlk.Timestamp(), parsedBanffStandardBlk.Timestamp())

		// backward compatibility check
		require.Equal(parsed.Txs(), parsedBanffStandardBlk.Txs())
	}
}

func TestProposalBlocks(t *testing.T) {
	// check Odyssey proposal block can be built and parsed
	require := require.New(t)
	blkTimestamp := time.Now()
	parentID := ids.ID{'p', 'a', 'r', 'e', 'n', 't', 'I', 'D'}
	height := uint64(2022)
	tx, err := testProposalTx()
	require.NoError(err)

	for _, cdc := range []codec.Manager{Codec, GenesisCodec} {
		// build block
		odysseyProposalBlk, err := NewOdysseyProposalBlock(
			parentID,
			height,
			tx,
		)
		require.NoError(err)

		// parse block
		parsed, err := Parse(cdc, odysseyProposalBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(odysseyProposalBlk.ID(), parsed.ID())
		require.Equal(odysseyProposalBlk.Bytes(), parsed.Bytes())
		require.Equal(odysseyProposalBlk.Parent(), parsed.Parent())
		require.Equal(odysseyProposalBlk.Height(), parsed.Height())

		require.IsType(&OdysseyProposalBlock{}, parsed)
		parsedOdysseyProposalBlk := parsed.(*OdysseyProposalBlock)
		require.Equal([]*txs.Tx{tx}, parsedOdysseyProposalBlk.Txs())

		// check that banff proposal block can be built and parsed
		banffProposalBlk, err := NewBanffProposalBlock(
			blkTimestamp,
			parentID,
			height,
			tx,
		)
		require.NoError(err)

		// parse block
		parsed, err = Parse(cdc, banffProposalBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(banffProposalBlk.ID(), parsed.ID())
		require.Equal(banffProposalBlk.Bytes(), parsed.Bytes())
		require.Equal(banffProposalBlk.Parent(), banffProposalBlk.Parent())
		require.Equal(banffProposalBlk.Height(), parsed.Height())
		require.IsType(&BanffProposalBlock{}, parsed)
		parsedBanffProposalBlk := parsed.(*BanffProposalBlock)
		require.Equal([]*txs.Tx{tx}, parsedBanffProposalBlk.Txs())

		// timestamp check for banff blocks only
		require.Equal(banffProposalBlk.Timestamp(), parsedBanffProposalBlk.Timestamp())

		// backward compatibility check
		require.Equal(parsedOdysseyProposalBlk.Txs(), parsedBanffProposalBlk.Txs())
	}
}

func TestCommitBlock(t *testing.T) {
	// check Odyssey commit block can be built and parsed
	require := require.New(t)
	blkTimestamp := time.Now()
	parentID := ids.ID{'p', 'a', 'r', 'e', 'n', 't', 'I', 'D'}
	height := uint64(2022)

	for _, cdc := range []codec.Manager{Codec, GenesisCodec} {
		// build block
		odysseyCommitBlk, err := NewOdysseyCommitBlock(parentID, height)
		require.NoError(err)

		// parse block
		parsed, err := Parse(cdc, odysseyCommitBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(odysseyCommitBlk.ID(), parsed.ID())
		require.Equal(odysseyCommitBlk.Bytes(), parsed.Bytes())
		require.Equal(odysseyCommitBlk.Parent(), parsed.Parent())
		require.Equal(odysseyCommitBlk.Height(), parsed.Height())

		// check that banff commit block can be built and parsed
		banffCommitBlk, err := NewBanffCommitBlock(blkTimestamp, parentID, height)
		require.NoError(err)

		// parse block
		parsed, err = Parse(cdc, banffCommitBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(banffCommitBlk.ID(), parsed.ID())
		require.Equal(banffCommitBlk.Bytes(), parsed.Bytes())
		require.Equal(banffCommitBlk.Parent(), banffCommitBlk.Parent())
		require.Equal(banffCommitBlk.Height(), parsed.Height())

		// timestamp check for banff blocks only
		require.IsType(&BanffCommitBlock{}, parsed)
		parsedBanffCommitBlk := parsed.(*BanffCommitBlock)
		require.Equal(banffCommitBlk.Timestamp(), parsedBanffCommitBlk.Timestamp())
	}
}

func TestAbortBlock(t *testing.T) {
	// check Odyssey abort block can be built and parsed
	require := require.New(t)
	blkTimestamp := time.Now()
	parentID := ids.ID{'p', 'a', 'r', 'e', 'n', 't', 'I', 'D'}
	height := uint64(2022)

	for _, cdc := range []codec.Manager{Codec, GenesisCodec} {
		// build block
		odysseyAbortBlk, err := NewOdysseyAbortBlock(parentID, height)
		require.NoError(err)

		// parse block
		parsed, err := Parse(cdc, odysseyAbortBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(odysseyAbortBlk.ID(), parsed.ID())
		require.Equal(odysseyAbortBlk.Bytes(), parsed.Bytes())
		require.Equal(odysseyAbortBlk.Parent(), parsed.Parent())
		require.Equal(odysseyAbortBlk.Height(), parsed.Height())

		// check that banff abort block can be built and parsed
		banffAbortBlk, err := NewBanffAbortBlock(blkTimestamp, parentID, height)
		require.NoError(err)

		// parse block
		parsed, err = Parse(cdc, banffAbortBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(banffAbortBlk.ID(), parsed.ID())
		require.Equal(banffAbortBlk.Bytes(), parsed.Bytes())
		require.Equal(banffAbortBlk.Parent(), banffAbortBlk.Parent())
		require.Equal(banffAbortBlk.Height(), parsed.Height())

		// timestamp check for banff blocks only
		require.IsType(&BanffAbortBlock{}, parsed)
		parsedBanffAbortBlk := parsed.(*BanffAbortBlock)
		require.Equal(banffAbortBlk.Timestamp(), parsedBanffAbortBlk.Timestamp())
	}
}

func TestAtomicBlock(t *testing.T) {
	// check atomic block can be built and parsed
	require := require.New(t)
	parentID := ids.ID{'p', 'a', 'r', 'e', 'n', 't', 'I', 'D'}
	height := uint64(2022)
	tx, err := testAtomicTx()
	require.NoError(err)

	for _, cdc := range []codec.Manager{Codec, GenesisCodec} {
		// build block
		atomicBlk, err := NewOdysseyAtomicBlock(
			parentID,
			height,
			tx,
		)
		require.NoError(err)

		// parse block
		parsed, err := Parse(cdc, atomicBlk.Bytes())
		require.NoError(err)

		// compare content
		require.Equal(atomicBlk.ID(), parsed.ID())
		require.Equal(atomicBlk.Bytes(), parsed.Bytes())
		require.Equal(atomicBlk.Parent(), parsed.Parent())
		require.Equal(atomicBlk.Height(), parsed.Height())

		require.IsType(&OdysseyAtomicBlock{}, parsed)
		parsedAtomicBlk := parsed.(*OdysseyAtomicBlock)
		require.Equal([]*txs.Tx{tx}, parsedAtomicBlk.Txs())
	}
}

func testAtomicTx() (*txs.Tx, error) {
	utx := &txs.ImportTx{
		BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
			NetworkID:    10,
			BlockchainID: ids.ID{'c', 'h', 'a', 'i', 'n', 'I', 'D'},
			Outs: []*dione.TransferableOutput{{
				Asset: dione.Asset{ID: ids.ID{'a', 's', 's', 'e', 'r', 't'}},
				Out: &secp256k1fx.TransferOutput{
					Amt: uint64(1234),
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
					},
				},
			}},
			Ins: []*dione.TransferableInput{{
				UTXOID: dione.UTXOID{
					TxID:        ids.ID{'t', 'x', 'I', 'D'},
					OutputIndex: 2,
				},
				Asset: dione.Asset{ID: ids.ID{'a', 's', 's', 'e', 'r', 't'}},
				In: &secp256k1fx.TransferInput{
					Amt:   uint64(5678),
					Input: secp256k1fx.Input{SigIndices: []uint32{0}},
				},
			}},
			Memo: []byte{1, 2, 3, 4, 5, 6, 7, 8},
		}},
		SourceChain: ids.ID{'c', 'h', 'a', 'i', 'n'},
		ImportedInputs: []*dione.TransferableInput{{
			UTXOID: dione.UTXOID{
				TxID:        ids.Empty.Prefix(1),
				OutputIndex: 1,
			},
			Asset: dione.Asset{ID: ids.ID{'a', 's', 's', 'e', 'r', 't'}},
			In: &secp256k1fx.TransferInput{
				Amt:   50000,
				Input: secp256k1fx.Input{SigIndices: []uint32{0}},
			},
		}},
	}
	signers := [][]*secp256k1.PrivateKey{{preFundedKeys[0]}}
	return txs.NewSigned(utx, txs.Codec, signers)
}

func testDecisionTxs() ([]*txs.Tx, error) {
	countTxs := 2
	decisionTxs := make([]*txs.Tx, 0, countTxs)
	for i := 0; i < countTxs; i++ {
		// Create the tx
		utx := &txs.CreateChainTx{
			BaseTx: txs.BaseTx{BaseTx: dione.BaseTx{
				NetworkID:    10,
				BlockchainID: ids.ID{'c', 'h', 'a', 'i', 'n', 'I', 'D'},
				Outs: []*dione.TransferableOutput{{
					Asset: dione.Asset{ID: ids.ID{'a', 's', 's', 'e', 'r', 't'}},
					Out: &secp256k1fx.TransferOutput{
						Amt: uint64(1234),
						OutputOwners: secp256k1fx.OutputOwners{
							Threshold: 1,
							Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
						},
					},
				}},
				Ins: []*dione.TransferableInput{{
					UTXOID: dione.UTXOID{
						TxID:        ids.ID{'t', 'x', 'I', 'D'},
						OutputIndex: 2,
					},
					Asset: dione.Asset{ID: ids.ID{'a', 's', 's', 'e', 'r', 't'}},
					In: &secp256k1fx.TransferInput{
						Amt:   uint64(5678),
						Input: secp256k1fx.Input{SigIndices: []uint32{0}},
					},
				}},
				Memo: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			}},
			SubnetID:    ids.ID{'s', 'u', 'b', 'n', 'e', 't', 'I', 'D'},
			ChainName:   "a chain",
			VMID:        ids.GenerateTestID(),
			FxIDs:       []ids.ID{ids.GenerateTestID()},
			GenesisData: []byte{'g', 'e', 'n', 'D', 'a', 't', 'a'},
			SubnetAuth:  &secp256k1fx.Input{SigIndices: []uint32{1}},
		}

		signers := [][]*secp256k1.PrivateKey{{preFundedKeys[0]}}
		tx, err := txs.NewSigned(utx, txs.Codec, signers)
		if err != nil {
			return nil, err
		}
		decisionTxs = append(decisionTxs, tx)
	}
	return decisionTxs, nil
}

func testProposalTx() (*txs.Tx, error) {
	utx := &txs.RewardValidatorTx{
		TxID: ids.ID{'r', 'e', 'w', 'a', 'r', 'd', 'I', 'D'},
	}

	signers := [][]*secp256k1.PrivateKey{{preFundedKeys[0]}}
	return txs.NewSigned(utx, txs.Codec, signers)
}
