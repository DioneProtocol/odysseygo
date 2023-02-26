// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/components/verify"
	"github.com/dioneprotocol/dionego/vms/platformvm/txs"
	"github.com/dioneprotocol/dionego/vms/platformvm/validator"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
)

func TestNewBanffStandardBlock(t *testing.T) {
	require := require.New(t)

	timestamp := time.Now().Truncate(time.Second)
	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: dione.BaseTx{
					Ins:  []*dione.TransferableInput{},
					Outs: []*dione.TransferableOutput{},
				},
			},
			StakeOuts: []*dione.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Initialize(txs.Codec))

	blk, err := NewBanffStandardBlock(
		timestamp,
		parentID,
		height,
		[]*txs.Tx{tx},
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotEmpty(blk.Bytes())
	require.NotEmpty(blk.Transactions[0].Bytes())
	require.NotEqual(ids.Empty, blk.Transactions[0].ID())
	require.Equal(tx.Bytes(), blk.Transactions[0].Bytes())
	require.Equal(timestamp, blk.Timestamp())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}

func TestNewApricotStandardBlock(t *testing.T) {
	require := require.New(t)

	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: dione.BaseTx{
					Ins:  []*dione.TransferableInput{},
					Outs: []*dione.TransferableOutput{},
				},
			},
			StakeOuts: []*dione.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Initialize(txs.Codec))

	blk, err := NewApricotStandardBlock(
		parentID,
		height,
		[]*txs.Tx{tx},
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotEmpty(blk.Bytes())
	require.NotEmpty(blk.Transactions[0].Bytes())
	require.NotEqual(ids.Empty, blk.Transactions[0].ID())
	require.Equal(tx.Bytes(), blk.Transactions[0].Bytes())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}
