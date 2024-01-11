// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/fx"
)

// ValidatorTx defines the interface for a validator transaction that supports
// delegation.
type ValidatorTx interface {
	UnsignedTx
	PermissionlessStaker

	ValidationRewardsOwner() fx.Owner
	DelegationRewardsOwner() fx.Owner
	Shares() uint32
}

type DelegatorTx interface {
	UnsignedTx
	PermissionlessStaker

	RewardsOwner() fx.Owner
}

type StakerTx interface {
	UnsignedTx
	Staker
}

type PermissionlessStaker interface {
	Staker

	Outputs() []*dione.TransferableOutput
	Stake() []*dione.TransferableOutput
}

type Staker interface {
	SubnetID() ids.ID
	NodeID() ids.NodeID
	// PublicKey returns the BLS public key registered by this transaction. If
	// there was no key registered by this transaction, it will return false.
	PublicKey() (*bls.PublicKey, bool, error)
	StartTime() time.Time
	EndTime() time.Time
	Weight() uint64
	PendingPriority() Priority
	CurrentPriority() Priority
}
