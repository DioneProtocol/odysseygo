// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/state"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

var (
	ErrChildBlockAfterStakerChangeTime = errors.New("proposed timestamp later than next staker change time")
	ErrChildBlockBeyondSyncBound       = errors.New("proposed timestamp is too far in the future relative to local time")
)

// VerifyNewChainTime returns nil if the [newChainTime] is a valid chain time
// given the wall clock time ([now]) and when the next staking set change occurs
// ([nextStakerChangeTime]).
// Requires:
//   - [newChainTime] <= [nextStakerChangeTime]: so that no staking set changes
//     are skipped.
//   - [newChainTime] <= [now] + [SyncBound]: to ensure chain time approximates
//     "real" time.
func VerifyNewChainTime(
	newChainTime,
	nextStakerChangeTime,
	now time.Time,
) error {
	// Only allow timestamp to move as far forward as the time of the next
	// staker set change
	if newChainTime.After(nextStakerChangeTime) {
		return fmt.Errorf(
			"%w, proposed timestamp (%s), next staker change time (%s)",
			ErrChildBlockAfterStakerChangeTime,
			newChainTime,
			nextStakerChangeTime,
		)
	}

	// Only allow timestamp to reasonably far forward
	maxNewChainTime := now.Add(SyncBound)
	if newChainTime.After(maxNewChainTime) {
		return fmt.Errorf(
			"%w, proposed time (%s), local time (%s)",
			ErrChildBlockBeyondSyncBound,
			newChainTime,
			now,
		)
	}
	return nil
}

type StateChanges interface {
	Apply(onAccept state.Diff)
	Len() int
}

type stateChanges struct {
	updatedSupplies           map[ids.ID]uint64
	currentValidatorsToAdd    []*state.Staker
	currentDelegatorsToAdd    []*state.Staker
	pendingValidatorsToRemove []*state.Staker
	pendingDelegatorsToRemove []*state.Staker
	currentValidatorsToRemove []*state.Staker
	stakeSyncTimestamp        time.Time
	accumulatedMintRate       *big.Int
	feePerWeightStored        *big.Int
	lastAccumulatedFee        uint64
}

func (s *stateChanges) Apply(stateDiff state.Diff) {
	for subnetID, supply := range s.updatedSupplies {
		stateDiff.SetCurrentSupply(subnetID, supply)
	}
	for _, currentValidatorToAdd := range s.currentValidatorsToAdd {
		stateDiff.PutCurrentValidator(currentValidatorToAdd)
	}
	for _, pendingValidatorToRemove := range s.pendingValidatorsToRemove {
		stateDiff.DeletePendingValidator(pendingValidatorToRemove)
	}
	for _, currentDelegatorToAdd := range s.currentDelegatorsToAdd {
		stateDiff.PutCurrentDelegator(currentDelegatorToAdd)
	}
	for _, pendingDelegatorToRemove := range s.pendingDelegatorsToRemove {
		stateDiff.DeletePendingDelegator(pendingDelegatorToRemove)
	}
	for _, currentValidatorToRemove := range s.currentValidatorsToRemove {
		stateDiff.DeleteCurrentValidator(currentValidatorToRemove)
	}
	if s.stakeSyncTimestamp.Compare(time.Time{}) != 0 {
		stateDiff.SetStakeSyncTimestamp(s.stakeSyncTimestamp)
	}
	if s.lastAccumulatedFee != 0 {
		stateDiff.SetLastAccumulatedFee(s.lastAccumulatedFee)
	}
	if s.accumulatedMintRate != nil {
		stateDiff.SetStakerAccumulatedMintRate(s.accumulatedMintRate)
	}
	if s.feePerWeightStored != nil {
		stateDiff.SetFeePerWeightStored(s.feePerWeightStored)
	}
}

func (s *stateChanges) Len() int {
	return len(s.currentValidatorsToAdd) + len(s.currentDelegatorsToAdd) +
		len(s.pendingValidatorsToRemove) + len(s.pendingDelegatorsToRemove) +
		len(s.currentValidatorsToRemove)
}

func (s *stateChanges) updateAccumulatedMintRate(backend *Backend, parentState state.Chain, newChainTime time.Time) error {
	if s.stakeSyncTimestamp.Compare(newChainTime) == 0 {
		return nil
	}

	validators, ok := backend.Config.Validators.Get(constants.PrimaryNetworkID)
	if !ok {
		return fmt.Errorf("couldn't get primary validators")
	}

	totalWeight := validators.Weight()
	mintRate, err := parentState.GetStakerAccumulatedMintRate()
	if err != nil {
		return err
	}

	lastSyncTime, err := parentState.GetStakeSyncTimestamp()
	if err != nil {
		return err
	}

	mintConfig := backend.Config.MintConfig

	// Config is not set
	if mintConfig.MintingPeriod == 0 {
		s.accumulatedMintRate = new(big.Int).SetUint64(0)
		return nil
	}

	if totalWeight != 0 {
		s.accumulatedMintRate = backend.Mint.CalculateMintRate(totalWeight, lastSyncTime, newChainTime)
		s.accumulatedMintRate.Add(s.accumulatedMintRate, mintRate)
	} else {
		s.accumulatedMintRate = new(big.Int).SetUint64(0)
	}
	s.stakeSyncTimestamp = newChainTime

	return nil
}

func (s *stateChanges) updateFeePerWeight(backend *Backend, parentState state.Chain) error {
	curAccumFee, err := parentState.GetCurrentAccumulatedFee()
	if err != nil {
		return err
	}

	lastAccumulatedFee, err := parentState.GetLastAccumulatedFee()
	if err != nil {
		return err
	}

	feePerWeightStored, err := parentState.GetFeePerWeightStored()
	if err != nil {
		return err
	}

	if lastAccumulatedFee == curAccumFee {
		return nil
	}

	if s.feePerWeightStored == nil {
		s.feePerWeightStored = new(big.Int)
	}

	vdrs, exists := backend.Config.Validators.Get(constants.OmegaChainID)
	if !exists {
		return fmt.Errorf("primary network vdrs not exists")
	}
	totalWeight := vdrs.Weight()
	if totalWeight == 0 {
		return nil
	}

	accumFeeDiff := curAccumFee - lastAccumulatedFee
	bigTotalWeight := new(big.Int).SetUint64(totalWeight)

	feePerWeightIncrement := new(big.Int).SetUint64(accumFeeDiff)
	feePerWeightIncrement.Lsh(feePerWeightIncrement, reward.BitShift)
	feePerWeightIncrement.Div(feePerWeightIncrement, bigTotalWeight)

	s.lastAccumulatedFee = curAccumFee
	s.feePerWeightStored.Set(feePerWeightStored)
	s.feePerWeightStored.Add(s.feePerWeightStored, feePerWeightIncrement)

	return nil
}

// AdvanceTimeTo does not modify [parentState].
// Instead it returns all the StateChanges caused by advancing the chain time to
// the [newChainTime].
func AdvanceTimeTo(
	backend *Backend,
	parentState state.Chain,
	newChainTime time.Time,
) (StateChanges, error) {
	pendingStakerIterator, err := parentState.GetPendingStakerIterator()
	if err != nil {
		return nil, err
	}
	defer pendingStakerIterator.Release()

	changes := &stateChanges{
		updatedSupplies: make(map[ids.ID]uint64),
	}

	// Add to the staker set any pending stakers whose start time is at or
	// before the new timestamp

	// Note: we process pending stakers ready to be promoted to current ones and
	// then we process current stakers to be demoted out of stakers set. It is
	// guaranteed that no promoted stakers would be demoted immediately. A
	// failure of this invariant would cause a staker to be added to
	// StateChanges and be persisted among current stakers even if it already
	// expired. The following invariants ensure this does not happens:
	// Invariant: minimum stake duration is > 0, so staker.StartTime != staker.EndTime.
	// Invariant: [newChainTime] does not skip stakers set change times.

	for pendingStakerIterator.Next() {
		stakerToRemove := pendingStakerIterator.Value()
		if stakerToRemove.StartTime.After(newChainTime) {
			break
		}

		stakerToAdd := *stakerToRemove
		stakerToAdd.NextTime = stakerToRemove.EndTime
		stakerToAdd.Priority = txs.PendingToCurrentPriorities[stakerToRemove.Priority]

		if stakerToRemove.Priority == txs.SubnetPermissionedValidatorPendingPriority {
			changes.currentValidatorsToAdd = append(changes.currentValidatorsToAdd, &stakerToAdd)
			changes.pendingValidatorsToRemove = append(changes.pendingValidatorsToRemove, stakerToRemove)
			continue
		}

		switch stakerToRemove.Priority {
		case txs.PrimaryNetworkValidatorPendingPriority, txs.PrimaryNetworkDelegatorApricotPendingPriority, txs.PrimaryNetworkDelegatorBanffPendingPriority:
			if err := changes.updateFeePerWeight(backend, parentState); err != nil {
				return nil, err
			}
			feePerWeightStored := changes.feePerWeightStored
			if feePerWeightStored == nil {
				feePerWeightStored, err = parentState.GetFeePerWeightStored()
				if err != nil {
					return nil, err
				}
			}
			stakerToAdd.FeePerWeightPaid = feePerWeightStored
			if err := changes.updateAccumulatedMintRate(backend, parentState, newChainTime); err != nil {
				return nil, err
			}
			stakerToAdd.MintRate = changes.accumulatedMintRate
		default:
			supply, ok := changes.updatedSupplies[stakerToRemove.SubnetID]
			if !ok {
				supply, err = parentState.GetCurrentSupply(stakerToRemove.SubnetID)
				if err != nil {
					return nil, err
				}
			}

			rewards, err := GetRewardsCalculator(backend, parentState, stakerToRemove.SubnetID)
			if err != nil {
				return nil, err
			}

			potentialReward := rewards.Calculate(
				stakerToRemove.EndTime.Sub(stakerToRemove.StartTime),
				stakerToRemove.Weight,
				supply,
			)
			stakerToAdd.PotentialReward = potentialReward

			// Invariant: [rewards.Calculate] can never return a [potentialReward]
			//            such that [supply + potentialReward > maximumSupply].
			changes.updatedSupplies[stakerToRemove.SubnetID] = supply + potentialReward
		}

		switch stakerToRemove.Priority {
		case txs.PrimaryNetworkValidatorPendingPriority, txs.SubnetPermissionlessValidatorPendingPriority:
			changes.currentValidatorsToAdd = append(changes.currentValidatorsToAdd, &stakerToAdd)
			changes.pendingValidatorsToRemove = append(changes.pendingValidatorsToRemove, stakerToRemove)

		case txs.PrimaryNetworkDelegatorApricotPendingPriority, txs.PrimaryNetworkDelegatorBanffPendingPriority, txs.SubnetPermissionlessDelegatorPendingPriority:
			changes.currentDelegatorsToAdd = append(changes.currentDelegatorsToAdd, &stakerToAdd)
			changes.pendingDelegatorsToRemove = append(changes.pendingDelegatorsToRemove, stakerToRemove)

		default:
			return nil, fmt.Errorf("expected staker priority got %d", stakerToRemove.Priority)
		}
	}

	currentStakerIterator, err := parentState.GetCurrentStakerIterator()
	if err != nil {
		return nil, err
	}
	defer currentStakerIterator.Release()

	for currentStakerIterator.Next() {
		stakerToRemove := currentStakerIterator.Value()
		if stakerToRemove.EndTime.After(newChainTime) {
			break
		}

		switch stakerToRemove.Priority {
		case txs.PrimaryNetworkValidatorCurrentPriority, txs.PrimaryNetworkDelegatorCurrentPriority:
			supply, ok := changes.updatedSupplies[stakerToRemove.SubnetID]
			if !ok {
				supply, err = parentState.GetCurrentSupply(stakerToRemove.SubnetID)
				if err != nil {
					return nil, err
				}
			}
			if err := changes.updateAccumulatedMintRate(backend, parentState, newChainTime); err != nil {
				return nil, err
			}
			if err := changes.updateFeePerWeight(backend, parentState); err != nil {
				return nil, err
			}

			feePerWeightStored := changes.feePerWeightStored
			if feePerWeightStored == nil {
				feePerWeightStored, err = parentState.GetFeePerWeightStored()
				if err != nil {
					return nil, err
				}
			}

			mint := reward.CalculateMintReward(stakerToRemove.Weight, stakerToRemove.MintRate, changes.accumulatedMintRate)
			fee := reward.CalculateFeeReward(feePerWeightStored, stakerToRemove.Weight, stakerToRemove.FeePerWeightPaid)

			stakerToRemove.PotentialReward = mint + fee
			changes.updatedSupplies[stakerToRemove.SubnetID] = supply + stakerToRemove.PotentialReward
		}

		// Invariant: Permissioned stakers are encountered first for a given
		//            timestamp because their priority is the smallest.
		if stakerToRemove.Priority != txs.SubnetPermissionedValidatorCurrentPriority {
			// Permissionless stakers are removed by the RewardValidatorTx, not
			// an AdvanceTimeTx.
			break
		}

		changes.currentValidatorsToRemove = append(changes.currentValidatorsToRemove, stakerToRemove)
	}
	return changes, nil
}

func GetRewardsCalculator(
	backend *Backend,
	parentState state.Chain,
	subnetID ids.ID,
) (reward.Calculator, error) {
	if subnetID == constants.PrimaryNetworkID {
		return backend.Rewards, nil
	}

	transformSubnetIntf, err := parentState.GetSubnetTransformation(subnetID)
	if err != nil {
		return nil, err
	}
	transformSubnet, ok := transformSubnetIntf.Unsigned.(*txs.TransformSubnetTx)
	if !ok {
		return nil, ErrIsNotTransformSubnetTx
	}

	return reward.NewCalculator(reward.Config{
		MaxConsumptionRate: transformSubnet.MaxConsumptionRate,
		MinConsumptionRate: transformSubnet.MinConsumptionRate,
		MintingPeriod:      backend.Config.RewardConfig.MintingPeriod,
		SupplyCap:          transformSubnet.MaximumSupply,
	}), nil
}
