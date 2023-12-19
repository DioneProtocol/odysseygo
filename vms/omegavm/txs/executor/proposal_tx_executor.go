// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"
	"fmt"
	"time"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

const (
	// Maximum future start time for staking
	MaxFutureStartTime = 24 * 7 * 2 * time.Hour

	// SyncBound is the synchrony bound used for safe decision making
	SyncBound = 10 * time.Second

	MaxValidatorWeightFactor = 5
)

var (
	_ txs.Visitor = (*ProposalTxExecutor)(nil)

	ErrRemoveStakerTooEarly          = errors.New("attempting to remove staker before their end time")
	ErrRemoveWrongStaker             = errors.New("attempting to remove wrong staker")
	ErrChildBlockNotAfterParent      = errors.New("proposed timestamp not after current chain time")
	ErrInvalidState                  = errors.New("generated output isn't valid state")
	ErrShouldBePermissionlessStaker  = errors.New("expected permissionless staker")
	ErrWrongTxType                   = errors.New("wrong transaction type")
	ErrInvalidID                     = errors.New("invalid ID")
	ErrProposedAddStakerTxAfterBanff = errors.New("staker transaction proposed after Banff")
	ErrAdvanceTimeTxIssuedAfterBanff = errors.New("AdvanceTimeTx issued after Banff")
)

type ProposalTxExecutor struct {
	// inputs, to be filled before visitor methods are called
	*Backend
	Tx *txs.Tx
	// [OnCommitState] is the state used for validation.
	// [OnCommitState] is modified by this struct's methods to
	// reflect changes made to the state if the proposal is committed.
	//
	// Invariant: Both [OnCommitState] and [OnAbortState] represent the same
	//            state when provided to this struct.
	OnCommitState state.Diff
	// [OnAbortState] is modified by this struct's methods to
	// reflect changes made to the state if the proposal is aborted.
	OnAbortState state.Diff

	// outputs populated by this struct's methods:
	//
	// [PrefersCommit] is true iff this node initially prefers to
	// commit this block transaction.
	PrefersCommit bool
}

func (*ProposalTxExecutor) CreateChainTx(*txs.CreateChainTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) CreateSubnetTx(*txs.CreateSubnetTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) ImportTx(*txs.ImportTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) ExportTx(*txs.ExportTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) RemoveSubnetValidatorTx(*txs.RemoveSubnetValidatorTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) TransformSubnetTx(*txs.TransformSubnetTx) error {
	return ErrWrongTxType
}

func (*ProposalTxExecutor) AddPermissionlessValidatorTx(*txs.AddPermissionlessValidatorTx) error {
	return ErrWrongTxType
}

func (e *ProposalTxExecutor) AddValidatorTx(tx *txs.AddValidatorTx) error {
	// AddValidatorTx is a proposal transaction until the Banff fork
	// activation. Following the activation, AddValidatorTxs must be issued into
	// StandardBlocks.
	currentTimestamp := e.OnCommitState.GetTimestamp()
	if e.Config.IsBanffActivated(currentTimestamp) {
		return fmt.Errorf(
			"%w: timestamp (%s) >= Banff fork time (%s)",
			ErrProposedAddStakerTxAfterBanff,
			currentTimestamp,
			e.Config.BanffTime,
		)
	}

	onAbortOuts, err := verifyAddValidatorTx(
		e.Backend,
		e.OnCommitState,
		e.Tx,
		tx,
	)
	if err != nil {
		return err
	}

	txID := e.Tx.ID()

	// Set up the state if this tx is committed
	// Consume the UTXOs
	dione.Consume(e.OnCommitState, tx.Ins)
	// Produce the UTXOs
	dione.Produce(e.OnCommitState, txID, tx.Outs)

	newStaker, err := state.NewPendingStaker(txID, tx)
	if err != nil {
		return err
	}

	e.OnCommitState.PutPendingValidator(newStaker)

	// Set up the state if this tx is aborted
	// Consume the UTXOs
	dione.Consume(e.OnAbortState, tx.Ins)
	// Produce the UTXOs
	dione.Produce(e.OnAbortState, txID, onAbortOuts)

	e.PrefersCommit = tx.StartTime().After(e.Clk.Time())
	return nil
}

func (e *ProposalTxExecutor) AddSubnetValidatorTx(tx *txs.AddSubnetValidatorTx) error {
	// AddSubnetValidatorTx is a proposal transaction until the Banff fork
	// activation. Following the activation, AddSubnetValidatorTxs must be
	// issued into StandardBlocks.
	currentTimestamp := e.OnCommitState.GetTimestamp()
	if e.Config.IsBanffActivated(currentTimestamp) {
		return fmt.Errorf(
			"%w: timestamp (%s) >= Banff fork time (%s)",
			ErrProposedAddStakerTxAfterBanff,
			currentTimestamp,
			e.Config.BanffTime,
		)
	}

	if err := verifyAddSubnetValidatorTx(
		e.Backend,
		e.OnCommitState,
		e.Tx,
		tx,
	); err != nil {
		return err
	}

	txID := e.Tx.ID()

	// Set up the state if this tx is committed
	// Consume the UTXOs
	dione.Consume(e.OnCommitState, tx.Ins)
	// Produce the UTXOs
	dione.Produce(e.OnCommitState, txID, tx.Outs)

	newStaker, err := state.NewPendingStaker(txID, tx)
	if err != nil {
		return err
	}

	e.OnCommitState.PutPendingValidator(newStaker)

	// Set up the state if this tx is aborted
	// Consume the UTXOs
	dione.Consume(e.OnAbortState, tx.Ins)
	// Produce the UTXOs
	dione.Produce(e.OnAbortState, txID, tx.Outs)

	e.PrefersCommit = tx.StartTime().After(e.Clk.Time())
	return nil
}

func (e *ProposalTxExecutor) AdvanceTimeTx(tx *txs.AdvanceTimeTx) error {
	switch {
	case tx == nil:
		return txs.ErrNilTx
	case len(e.Tx.Creds) != 0:
		return errWrongNumberOfCredentials
	}

	// Validate [newChainTime]
	newChainTime := tx.Timestamp()
	if e.Config.IsBanffActivated(newChainTime) {
		return fmt.Errorf(
			"%w: proposed timestamp (%s) >= Banff fork time (%s)",
			ErrAdvanceTimeTxIssuedAfterBanff,
			newChainTime,
			e.Config.BanffTime,
		)
	}

	parentChainTime := e.OnCommitState.GetTimestamp()
	if !newChainTime.After(parentChainTime) {
		return fmt.Errorf(
			"%w, proposed timestamp (%s), chain time (%s)",
			ErrChildBlockNotAfterParent,
			parentChainTime,
			parentChainTime,
		)
	}

	// Only allow timestamp to move forward as far as the time of next staker
	// set change time
	nextStakerChangeTime, err := GetNextStakerChangeTime(e.OnCommitState)
	if err != nil {
		return err
	}

	now := e.Clk.Time()
	if err := VerifyNewChainTime(
		newChainTime,
		nextStakerChangeTime,
		now,
	); err != nil {
		return err
	}

	changes, err := AdvanceTimeTo(e.Backend, e.OnCommitState, newChainTime)
	if err != nil {
		return err
	}

	// Update the state if this tx is committed
	e.OnCommitState.SetTimestamp(newChainTime)
	changes.Apply(e.OnCommitState)

	e.PrefersCommit = !newChainTime.After(now.Add(SyncBound))

	// Note that state doesn't change if this proposal is aborted
	return nil
}

func (e *ProposalTxExecutor) RewardValidatorTx(tx *txs.RewardValidatorTx) error {
	switch {
	case tx == nil:
		return txs.ErrNilTx
	case tx.TxID == ids.Empty:
		return ErrInvalidID
	case len(e.Tx.Creds) != 0:
		return errWrongNumberOfCredentials
	}

	currentStakerIterator, err := e.OnCommitState.GetCurrentStakerIterator()
	if err != nil {
		return err
	}
	if !currentStakerIterator.Next() {
		return fmt.Errorf("failed to get next staker to remove: %w", database.ErrNotFound)
	}
	stakerToRemove := currentStakerIterator.Value()
	currentStakerIterator.Release()

	if stakerToRemove.TxID != tx.TxID {
		return fmt.Errorf(
			"%w: %s != %s",
			ErrRemoveWrongStaker,
			stakerToRemove.TxID,
			tx.TxID,
		)
	}

	// Verify that the chain's timestamp is the validator's end time
	currentChainTime := e.OnCommitState.GetTimestamp()
	if !stakerToRemove.EndTime.Equal(currentChainTime) {
		return fmt.Errorf(
			"%w: TxID = %s with %s < %s",
			ErrRemoveStakerTooEarly,
			tx.TxID,
			currentChainTime,
			stakerToRemove.EndTime,
		)
	}

	primaryNetworkValidator, err := e.OnCommitState.GetCurrentValidator(
		constants.PrimaryNetworkID,
		stakerToRemove.NodeID,
	)
	if err != nil {
		// This should never error because the staker set is in memory and
		// primary network validators are removed last.
		return err
	}

	stakerTx, _, err := e.OnCommitState.GetTx(stakerToRemove.TxID)
	if err != nil {
		return fmt.Errorf("failed to get next removed staker tx: %w", err)
	}

	switch uStakerTx := stakerTx.Unsigned.(type) {
	case txs.ValidatorTx:
		e.OnCommitState.DeleteCurrentValidator(stakerToRemove)
		e.OnAbortState.DeleteCurrentValidator(stakerToRemove)

		stake := uStakerTx.Stake()
		outputs := uStakerTx.Outputs()
		// Invariant: The staked asset must be equal to the reward asset.
		stakeAsset := stake[0].Asset

		// Refund the stake here
		for i, out := range stake {
			utxo := &dione.UTXO{
				UTXOID: dione.UTXOID{
					TxID:        tx.TxID,
					OutputIndex: uint32(len(outputs) + i),
				},
				Asset: out.Asset,
				Out:   out.Output(),
			}
			e.OnCommitState.AddUTXO(utxo)
			e.OnAbortState.AddUTXO(utxo)
		}

		offset := 0

		// Provide the reward here
		if stakerToRemove.PotentialReward > 0 {
			validationRewardsOwner := uStakerTx.ValidationRewardsOwner()
			outIntf, err := e.Fx.CreateOutput(stakerToRemove.PotentialReward, validationRewardsOwner)
			if err != nil {
				return fmt.Errorf("failed to create output: %w", err)
			}
			out, ok := outIntf.(verify.State)
			if !ok {
				return ErrInvalidState
			}

			utxo := &dione.UTXO{
				UTXOID: dione.UTXOID{
					TxID:        tx.TxID,
					OutputIndex: uint32(len(outputs) + len(stake)),
				},
				Asset: stakeAsset,
				Out:   out,
			}

			e.OnCommitState.AddUTXO(utxo)
			e.OnCommitState.AddRewardUTXO(tx.TxID, utxo)

			offset++
		}
	default:
		// Invariant: Permissioned stakers are removed by the advancement of
		//            time and the current chain timestamp is == this staker's
		//            EndTime. This means only permissionless stakers should be
		//            left in the staker set.
		return ErrShouldBePermissionlessStaker
	}

	// If the reward is aborted, then the current supply should be decreased.
	currentSupply, err := e.OnAbortState.GetCurrentSupply(stakerToRemove.SubnetID)
	if err != nil {
		return err
	}
	newSupply, err := math.Sub(currentSupply, stakerToRemove.PotentialReward)
	if err != nil {
		return err
	}
	e.OnAbortState.SetCurrentSupply(stakerToRemove.SubnetID, newSupply)

	var expectedUptimePercentage float64
	if stakerToRemove.SubnetID != constants.PrimaryNetworkID {
		transformSubnetIntf, err := e.OnCommitState.GetSubnetTransformation(stakerToRemove.SubnetID)
		if err != nil {
			return err
		}
		transformSubnet, ok := transformSubnetIntf.Unsigned.(*txs.TransformSubnetTx)
		if !ok {
			return ErrIsNotTransformSubnetTx
		}

		expectedUptimePercentage = float64(transformSubnet.UptimeRequirement) / reward.PercentDenominator
	} else {
		expectedUptimePercentage = e.Config.UptimePercentage
	}

	// TODO: calculate subnet uptimes
	uptime, err := e.Uptimes.CalculateUptimePercentFrom(
		primaryNetworkValidator.NodeID,
		constants.PrimaryNetworkID,
		primaryNetworkValidator.StartTime,
	)
	if err != nil {
		return fmt.Errorf("failed to calculate uptime: %w", err)
	}

	e.PrefersCommit = uptime >= expectedUptimePercentage
	return nil
}

// GetNextStakerChangeTime returns the next time a staker will be either added
// or removed to/from the current validator set.
func GetNextStakerChangeTime(state state.Chain) (time.Time, error) {
	currentStakerIterator, err := state.GetCurrentStakerIterator()
	if err != nil {
		return time.Time{}, err
	}
	defer currentStakerIterator.Release()

	pendingStakerIterator, err := state.GetPendingStakerIterator()
	if err != nil {
		return time.Time{}, err
	}
	defer pendingStakerIterator.Release()

	hasCurrentStaker := currentStakerIterator.Next()
	hasPendingStaker := pendingStakerIterator.Next()
	switch {
	case hasCurrentStaker && hasPendingStaker:
		nextCurrentTime := currentStakerIterator.Value().NextTime
		nextPendingTime := pendingStakerIterator.Value().NextTime
		if nextCurrentTime.Before(nextPendingTime) {
			return nextCurrentTime, nil
		}
		return nextPendingTime, nil
	case hasCurrentStaker:
		return currentStakerIterator.Value().NextTime, nil
	case hasPendingStaker:
		return pendingStakerIterator.Value().NextTime, nil
	default:
		return time.Time{}, database.ErrNotFound
	}
}

// GetValidator returns information about the given validator, which may be a
// current validator or pending validator.
func GetValidator(state state.Chain, subnetID ids.ID, nodeID ids.NodeID) (*state.Staker, error) {
	validator, err := state.GetCurrentValidator(subnetID, nodeID)
	if err == nil {
		// This node is currently validating the subnet.
		return validator, nil
	}
	if err != database.ErrNotFound {
		// Unexpected error occurred.
		return nil, err
	}
	return state.GetPendingValidator(subnetID, nodeID)
}
