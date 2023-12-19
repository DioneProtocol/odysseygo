// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

const (
	// then primary network validators,
	PrimaryNetworkValidatorPendingPriority Priority = iota + 1
	// then permissionless subnet validators,
	SubnetPermissionlessValidatorPendingPriority
	// then permissioned subnet validators,
	SubnetPermissionedValidatorPendingPriority

	// First permissioned subnet validators are removed from the current
	// validator set,
	// Invariant: All permissioned stakers must be removed first because they
	//            are removed by the advancement of time. Permissionless stakers
	//            are removed with a RewardValidatorTx after time has advanced.
	SubnetPermissionedValidatorCurrentPriority
	// then permissionless subnet validators,
	SubnetPermissionlessValidatorCurrentPriority
	// then primary network validators.
	PrimaryNetworkValidatorCurrentPriority
)

var PendingToCurrentPriorities = []Priority{
	PrimaryNetworkValidatorPendingPriority:       PrimaryNetworkValidatorCurrentPriority,
	SubnetPermissionlessValidatorPendingPriority: SubnetPermissionlessValidatorCurrentPriority,
	SubnetPermissionedValidatorPendingPriority:   SubnetPermissionedValidatorCurrentPriority,
}

type Priority byte

func (p Priority) IsCurrent() bool {
	return p.IsCurrentValidator()
}

func (p Priority) IsPending() bool {
	return p.IsPendingValidator()
}

func (p Priority) IsValidator() bool {
	return p.IsCurrentValidator() || p.IsPendingValidator()
}

func (p Priority) IsPermissionedValidator() bool {
	return p == SubnetPermissionedValidatorCurrentPriority ||
		p == SubnetPermissionedValidatorPendingPriority
}

func (p Priority) IsCurrentValidator() bool {
	return p == PrimaryNetworkValidatorCurrentPriority ||
		p == SubnetPermissionedValidatorCurrentPriority ||
		p == SubnetPermissionlessValidatorCurrentPriority
}

func (p Priority) IsPendingValidator() bool {
	return p == PrimaryNetworkValidatorPendingPriority ||
		p == SubnetPermissionedValidatorPendingPriority ||
		p == SubnetPermissionlessValidatorPendingPriority
}
