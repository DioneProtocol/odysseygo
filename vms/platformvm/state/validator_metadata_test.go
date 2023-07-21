// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/codec"
	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/database/memdb"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/wrappers"
)

func TestValidatorUptimes(t *testing.T) {
	require := require.New(t)
	state := newValidatorState()

	// get non-existent uptime
	nodeID := ids.GenerateTestNodeID()
	subnetID := ids.GenerateTestID()
	_, _, err := state.GetUptime(nodeID, subnetID)
	require.ErrorIs(err, database.ErrNotFound)

	// set non-existent uptime
	err = state.SetUptime(nodeID, subnetID, 1, time.Now())
	require.ErrorIs(err, database.ErrNotFound)

	testMetadata := &validatorMetadata{
		UpDuration:  time.Hour,
		lastUpdated: time.Now(),
	}
	// load uptime
	state.LoadValidatorMetadata(nodeID, subnetID, testMetadata)

	// get uptime
	upDuration, lastUpdated, err := state.GetUptime(nodeID, subnetID)
	require.NoError(err)
	require.Equal(testMetadata.UpDuration, upDuration)
	require.Equal(testMetadata.lastUpdated, lastUpdated)

	// set uptime
	newUpDuration := testMetadata.UpDuration + 1
	newLastUpdated := testMetadata.lastUpdated.Add(time.Hour)
	err = state.SetUptime(nodeID, subnetID, newUpDuration, newLastUpdated)
	require.NoError(err)

	// get new uptime
	upDuration, lastUpdated, err = state.GetUptime(nodeID, subnetID)
	require.NoError(err)
	require.Equal(newUpDuration, upDuration)
	require.Equal(newLastUpdated, lastUpdated)

	// load uptime changes uptimes
	newTestMetadata := &validatorMetadata{
		UpDuration:  testMetadata.UpDuration + time.Hour,
		lastUpdated: testMetadata.lastUpdated.Add(time.Hour),
	}
	state.LoadValidatorMetadata(nodeID, subnetID, newTestMetadata)

	// get new uptime
	upDuration, lastUpdated, err = state.GetUptime(nodeID, subnetID)
	require.NoError(err)
	require.Equal(newTestMetadata.UpDuration, upDuration)
	require.Equal(newTestMetadata.lastUpdated, lastUpdated)

	// delete uptime
	state.DeleteValidatorMetadata(nodeID, subnetID)

	// get deleted uptime
	_, _, err = state.GetUptime(nodeID, subnetID)
	require.ErrorIs(err, database.ErrNotFound)
}

func TestWriteValidatorMetadata(t *testing.T) {
	require := require.New(t)
	state := newValidatorState()

	primaryDB := memdb.New()
	subnetDB := memdb.New()
	// write empty uptimes
	err := state.WriteValidatorMetadata(primaryDB, subnetDB)
	require.NoError(err)

	// load uptime
	nodeID := ids.GenerateTestNodeID()
	subnetID := ids.GenerateTestID()
	testUptimeReward := &validatorMetadata{
		UpDuration:      time.Hour,
		lastUpdated:     time.Now(),
		PotentialReward: 100,
		txID:            ids.GenerateTestID(),
	}
	state.LoadValidatorMetadata(nodeID, subnetID, testUptimeReward)

	// write state, should not reflect to DB yet
	err = state.WriteValidatorMetadata(primaryDB, subnetDB)
	require.NoError(err)
	require.False(primaryDB.Has(testUptimeReward.txID[:]))
	require.False(subnetDB.Has(testUptimeReward.txID[:]))

	// get uptime should still return the loaded value
	upDuration, lastUpdated, err := state.GetUptime(nodeID, subnetID)
	require.NoError(err)
	require.Equal(testUptimeReward.UpDuration, upDuration)
	require.Equal(testUptimeReward.lastUpdated, lastUpdated)

	// update uptimes
	newUpDuration := testUptimeReward.UpDuration + 1
	newLastUpdated := testUptimeReward.lastUpdated.Add(time.Hour)
	err = state.SetUptime(nodeID, subnetID, newUpDuration, newLastUpdated)
	require.NoError(err)

	// write uptimes, should reflect to subnet DB
	err = state.WriteValidatorMetadata(primaryDB, subnetDB)
	require.NoError(err)
	require.False(primaryDB.Has(testUptimeReward.txID[:]))
	require.True(subnetDB.Has(testUptimeReward.txID[:]))
}

func TestParseValidatorMetadata(t *testing.T) {
	type test struct {
		name        string
		bytes       []byte
		expected    *validatorMetadata
		expectedErr error
	}
	tests := []test{
		{
			name:  "nil",
			bytes: nil,
			expected: &validatorMetadata{
				lastUpdated: time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name:  "nil",
			bytes: []byte{},
			expected: &validatorMetadata{
				lastUpdated: time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name: "potential reward only",
			bytes: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0,
			},
			expected: &validatorMetadata{
				PotentialReward: 100000,
				lastUpdated:     time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name: "uptime + potential reward",
			bytes: []byte{
				// codec version
				0x00, 0x00,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
				// potential reward
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0,
			},
			expected: &validatorMetadata{
				UpDuration:      time.Duration(6000000),
				LastUpdated:     900000,
				PotentialReward: 100000,
				lastUpdated:     time.Unix(900000, 0),
			},
			expectedErr: nil,
		},
		{
			name: "uptime + potential reward",
			bytes: []byte{
				// codec version
				0x00, 0x00,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
				// potential reward
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0,
			},
			expected: &validatorMetadata{
				UpDuration:               time.Duration(6000000),
				LastUpdated:              900000,
				PotentialReward:          100000,
				lastUpdated:              time.Unix(900000, 0),
			},
			expectedErr: nil,
		},
		{
			name: "invalid codec version",
			bytes: []byte{
				// codec version
				0x00, 0x01,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
				// potential reward
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0,
			},
			expected:    nil,
			expectedErr: codec.ErrUnknownVersion,
		},
		{
			name: "short byte len",
			bytes: []byte{
				// codec version
				0x00, 0x00,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
				// potential reward
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0xA0,
			},
			expected:    nil,
			expectedErr: wrappers.ErrInsufficientLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			var metadata validatorMetadata
			err := parseValidatorMetadata(tt.bytes, &metadata)
			require.ErrorIs(err, tt.expectedErr)
			if err == nil {
				require.Equal(tt.expected, &metadata)
			}
		})
	}
}
