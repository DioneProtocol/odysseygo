// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPriorityIsCurrent(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsCurrent())
		})
	}
}

func TestPriorityIsPending(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: false,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsPending())
		})
	}
}

func TestPriorityIsValidator(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsValidator())
		})
	}
}

func TestPriorityIsPermissionedValidator(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: false,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsPermissionedValidator())
		})
	}
}

func TestPriorityIsCurrentValidator(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: true,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsCurrentValidator())
		})
	}
}

func TestPriorityIsPendingValidator(t *testing.T) {
	tests := []struct {
		priority Priority
		expected bool
	}{
		{
			priority: PrimaryNetworkValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionlessValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorPendingPriority,
			expected: true,
		},
		{
			priority: SubnetPermissionedValidatorCurrentPriority,
			expected: false,
		},
		{
			priority: SubnetPermissionlessValidatorCurrentPriority,
			expected: false,
		},
		{
			priority: PrimaryNetworkValidatorCurrentPriority,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.priority), func(t *testing.T) {
			require.Equal(t, test.expected, test.priority.IsPendingValidator())
		})
	}
}
