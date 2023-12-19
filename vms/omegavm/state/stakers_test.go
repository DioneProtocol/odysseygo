// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

func TestBaseStakersPruning(t *testing.T) {
	require := require.New(t)
	staker := newTestStaker()

	v := newBaseStakers()

	v.PutValidator(staker)

	_, err := v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)

	v.DeleteValidator(staker)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.ErrorIs(err, database.ErrNotFound)

	require.Empty(v.validators)

	v.PutValidator(staker)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)

	v.DeleteValidator(staker)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.ErrorIs(err, database.ErrNotFound)

	require.Empty(v.validators)
}

func TestBaseStakersValidator(t *testing.T) {
	require := require.New(t)
	staker := newTestStaker()

	v := newBaseStakers()

	_, err := v.GetValidator(ids.GenerateTestID(), staker.NodeID)
	require.ErrorIs(err, database.ErrNotFound)

	_, err = v.GetValidator(staker.SubnetID, ids.GenerateTestNodeID())
	require.ErrorIs(err, database.ErrNotFound)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.ErrorIs(err, database.ErrNotFound)

	stakerIterator := v.GetStakerIterator()

	v.PutValidator(staker)

	returnedStaker, err := v.GetValidator(staker.SubnetID, staker.NodeID)
	require.NoError(err)
	require.Equal(staker, returnedStaker)

	stakerIterator = v.GetStakerIterator()
	assertIteratorsEqual(t, NewSliceIterator(staker), stakerIterator)

	v.DeleteValidator(staker)

	_, err = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.ErrorIs(err, database.ErrNotFound)

	stakerIterator = v.GetStakerIterator()
	assertIteratorsEqual(t, EmptyIterator, stakerIterator)
}

func TestDiffStakersValidator(t *testing.T) {
	require := require.New(t)
	staker := newTestStaker()

	v := diffStakers{}

	// validators not available in the diff are marked as unmodified
	_, status := v.GetValidator(ids.GenerateTestID(), staker.NodeID)
	require.Equal(unmodified, status)

	_, status = v.GetValidator(staker.SubnetID, ids.GenerateTestNodeID())
	require.Equal(unmodified, status)

	stakerIterator := v.GetStakerIterator(EmptyIterator)

	v.PutValidator(staker)

	returnedStaker, status := v.GetValidator(staker.SubnetID, staker.NodeID)
	require.Equal(added, status)
	require.Equal(staker, returnedStaker)

	v.DeleteValidator(staker)

	// Validators created and deleted in the same diff are marked as unmodified.
	// This means they won't be pushed to baseState if diff.Apply(baseState) is
	// called.
	_, status = v.GetValidator(staker.SubnetID, staker.NodeID)
	require.Equal(unmodified, status)
}

func TestDiffStakersDeleteValidator(t *testing.T) {
	require := require.New(t)
	staker := newTestStaker()

	v := diffStakers{}

	_, status := v.GetValidator(ids.GenerateTestID(), staker.NodeID)
	require.Equal(unmodified, status)

	v.DeleteValidator(staker)

	returnedStaker, status := v.GetValidator(staker.SubnetID, staker.NodeID)
	require.Equal(deleted, status)
	require.Nil(returnedStaker)
}

func newTestStaker() *Staker {
	startTime := time.Now().Round(time.Second)
	endTime := startTime.Add(28 * 24 * time.Hour)
	return &Staker{
		TxID:            ids.GenerateTestID(),
		NodeID:          ids.GenerateTestNodeID(),
		SubnetID:        ids.GenerateTestID(),
		Weight:          1,
		StartTime:       startTime,
		EndTime:         endTime,
		PotentialReward: 1,

		NextTime: endTime,
		Priority: txs.SubnetPermissionlessValidatorCurrentPriority,
	}
}

func assertIteratorsEqual(t *testing.T, expected, actual StakerIterator) {
	t.Helper()

	for expected.Next() {
		require.True(t, actual.Next())

		expectedStaker := expected.Value()
		actualStaker := actual.Value()

		require.Equal(t, expectedStaker, actualStaker)
	}
	require.False(t, actual.Next())

	expected.Release()
	actual.Release()
}
