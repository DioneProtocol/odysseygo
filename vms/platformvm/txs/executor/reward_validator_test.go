// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/reward"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/state"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/status"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

func TestRewardValidatorTxExecuteOnCommit(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(false /*=postBanff*/, false /*=postCortina*/)
	defer func() {
		require.NoError(shutdownEnvironment(env))
	}()
	dummyHeight := uint64(1)

	currentStakerIterator, err := env.state.GetCurrentStakerIterator()
	require.NoError(err)
	require.True(currentStakerIterator.Next())

	stakerToRemove := currentStakerIterator.Value()
	currentStakerIterator.Release()

	stakerToRemoveTxIntf, _, err := env.state.GetTx(stakerToRemove.TxID)
	require.NoError(err)
	stakerToRemoveTx := stakerToRemoveTxIntf.Unsigned.(*txs.AddValidatorTx)

	// Case 1: Chain timestamp is wrong
	tx, err := env.txBuilder.NewRewardValidatorTx(stakerToRemove.TxID)
	require.NoError(err)

	onCommitState, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	onAbortState, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	txExecutor := ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	err = tx.Unsigned.Visit(&txExecutor)
	require.ErrorIs(err, ErrRemoveStakerTooEarly)

	// Advance chain timestamp to time that next validator leaves
	env.state.SetTimestamp(stakerToRemove.EndTime)

	// Case 2: Wrong validator
	tx, err = env.txBuilder.NewRewardValidatorTx(ids.GenerateTestID())
	require.NoError(err)

	onCommitState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	onAbortState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	txExecutor = ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	err = tx.Unsigned.Visit(&txExecutor)
	require.ErrorIs(err, ErrRemoveWrongStaker)

	// Case 3: Happy path
	tx, err = env.txBuilder.NewRewardValidatorTx(stakerToRemove.TxID)
	require.NoError(err)

	onCommitState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	onAbortState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	txExecutor = ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	require.NoError(tx.Unsigned.Visit(&txExecutor))

	onCommitStakerIterator, err := txExecutor.OnCommitState.GetCurrentStakerIterator()
	require.NoError(err)
	require.True(onCommitStakerIterator.Next())

	nextToRemove := onCommitStakerIterator.Value()
	onCommitStakerIterator.Release()
	require.NotEqual(stakerToRemove.TxID, nextToRemove.TxID)

	// check that stake/reward is given back
	stakeOwners := stakerToRemoveTx.StakeOuts[0].Out.(*secp256k1fx.TransferOutput).AddressesSet()

	// Get old balances
	oldBalance, err := dione.GetBalance(env.state, stakeOwners)
	require.NoError(err)

	require.NoError(txExecutor.OnCommitState.Apply(env.state))

	env.state.SetHeight(dummyHeight)
	require.NoError(env.state.Commit())

	onCommitBalance, err := dione.GetBalance(env.state, stakeOwners)
	require.NoError(err)
	require.Equal(oldBalance+stakerToRemove.Weight+27697, onCommitBalance)
}

func TestRewardValidatorTxExecuteOnAbort(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(false /*=postBanff*/, false /*=postCortina*/)
	defer func() {
		require.NoError(shutdownEnvironment(env))
	}()
	dummyHeight := uint64(1)

	currentStakerIterator, err := env.state.GetCurrentStakerIterator()
	require.NoError(err)
	require.True(currentStakerIterator.Next())

	stakerToRemove := currentStakerIterator.Value()
	currentStakerIterator.Release()

	stakerToRemoveTxIntf, _, err := env.state.GetTx(stakerToRemove.TxID)
	require.NoError(err)
	stakerToRemoveTx := stakerToRemoveTxIntf.Unsigned.(*txs.AddValidatorTx)

	// Case 1: Chain timestamp is wrong
	tx, err := env.txBuilder.NewRewardValidatorTx(stakerToRemove.TxID)
	require.NoError(err)

	onCommitState, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	onAbortState, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	txExecutor := ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	err = tx.Unsigned.Visit(&txExecutor)
	require.ErrorIs(err, ErrRemoveStakerTooEarly)

	// Advance chain timestamp to time that next validator leaves
	env.state.SetTimestamp(stakerToRemove.EndTime)

	// Case 2: Wrong validator
	tx, err = env.txBuilder.NewRewardValidatorTx(ids.GenerateTestID())
	require.NoError(err)

	txExecutor = ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	err = tx.Unsigned.Visit(&txExecutor)
	require.ErrorIs(err, ErrRemoveWrongStaker)

	// Case 3: Happy path
	tx, err = env.txBuilder.NewRewardValidatorTx(stakerToRemove.TxID)
	require.NoError(err)

	onCommitState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	onAbortState, err = state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	txExecutor = ProposalTxExecutor{
		OnCommitState: onCommitState,
		OnAbortState:  onAbortState,
		Backend:       &env.backend,
		Tx:            tx,
	}
	require.NoError(tx.Unsigned.Visit(&txExecutor))

	onAbortStakerIterator, err := txExecutor.OnAbortState.GetCurrentStakerIterator()
	require.NoError(err)
	require.True(onAbortStakerIterator.Next())

	nextToRemove := onAbortStakerIterator.Value()
	onAbortStakerIterator.Release()
	require.NotEqual(stakerToRemove.TxID, nextToRemove.TxID)

	// check that stake/reward isn't given back
	stakeOwners := stakerToRemoveTx.StakeOuts[0].Out.(*secp256k1fx.TransferOutput).AddressesSet()

	// Get old balances
	oldBalance, err := dione.GetBalance(env.state, stakeOwners)
	require.NoError(err)

	require.NoError(txExecutor.OnAbortState.Apply(env.state))

	env.state.SetHeight(dummyHeight)
	require.NoError(env.state.Commit())

	onAbortBalance, err := dione.GetBalance(env.state, stakeOwners)
	require.NoError(err)
	require.Equal(oldBalance+stakerToRemove.Weight, onAbortBalance)
}
