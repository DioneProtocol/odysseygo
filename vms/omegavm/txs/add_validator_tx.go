// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"fmt"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/fx"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var (
	_ ValidatorTx = (*AddValidatorTx)(nil)
)

// AddValidatorTx is an unsigned addValidatorTx
type AddValidatorTx struct {
	// Metadata, inputs and outputs
	BaseTx `serialize:"true"`
	// Describes the validator
	Validator `serialize:"true" json:"validator"`
	// Where to send staked tokens when done validating
	StakeOuts []*dione.TransferableOutput `serialize:"true" json:"stake"`
	// Where to send staking rewards when done validating
	RewardsOwner fx.Owner `serialize:"true" json:"rewardsOwner"`
}

// InitCtx sets the FxID fields in the inputs and outputs of this
// [AddValidatorTx]. Also sets the [ctx] to the given [vm.ctx] so that
// the addresses can be json marshalled into human readable format
func (tx *AddValidatorTx) InitCtx(ctx *snow.Context) {
	tx.BaseTx.InitCtx(ctx)
	for _, out := range tx.StakeOuts {
		out.FxID = secp256k1fx.ID
		out.InitCtx(ctx)
	}
	tx.RewardsOwner.InitCtx(ctx)
}

func (*AddValidatorTx) SubnetID() ids.ID {
	return constants.PrimaryNetworkID
}

func (tx *AddValidatorTx) NodeID() ids.NodeID {
	return tx.Validator.NodeID
}

func (*AddValidatorTx) PublicKey() (*bls.PublicKey, bool, error) {
	return nil, false, nil
}

func (*AddValidatorTx) PendingPriority() Priority {
	return PrimaryNetworkValidatorPendingPriority
}

func (*AddValidatorTx) CurrentPriority() Priority {
	return PrimaryNetworkValidatorCurrentPriority
}

func (tx *AddValidatorTx) Stake() []*dione.TransferableOutput {
	return tx.StakeOuts
}

func (tx *AddValidatorTx) ValidationRewardsOwner() fx.Owner {
	return tx.RewardsOwner
}

// SyntacticVerify returns nil iff [tx] is valid
func (tx *AddValidatorTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return ErrNilTx
	case tx.SyntacticallyVerified: // already passed syntactic verification
		return nil
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return fmt.Errorf("failed to verify BaseTx: %w", err)
	}
	if err := verify.All(&tx.Validator, tx.RewardsOwner); err != nil {
		return fmt.Errorf("failed to verify validator or rewards owner: %w", err)
	}

	// cache that this is valid
	tx.SyntacticallyVerified = true
	return nil
}

func (tx *AddValidatorTx) Visit(visitor Visitor) error {
	return visitor.AddValidatorTx(tx)
}
