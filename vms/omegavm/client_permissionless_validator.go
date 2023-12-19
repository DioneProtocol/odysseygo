// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package omegavm

import (
	"encoding/json"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/api"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/signer"
)

// ClientStaker is the representation of a staker sent via client.
type ClientStaker struct {
	// the txID of the transaction that added this staker.
	TxID ids.ID
	// the Unix time when they start staking
	StartTime uint64
	// the Unix time when they are done staking
	EndTime uint64
	// the validator weight when sampling validators
	Weight uint64
	// the amount of tokens being staked.
	StakeAmount *uint64
	// the node ID of the staker
	NodeID ids.NodeID
}

// ClientOwner is the repr. of a reward owner sent over client
type ClientOwner struct {
	Locktime  uint64
	Threshold uint32
	Addresses []ids.ShortID
}

// ClientPermissionlessValidator is the repr. of a permissionless validator sent
// over client
type ClientPermissionlessValidator struct {
	ClientStaker
	ValidationRewardOwner *ClientOwner
	PotentialReward       *uint64
	Uptime                *float32
	Connected             *bool
	Signer                *signer.ProofOfPossession
}

func apiStakerToClientStaker(validator api.Staker) ClientStaker {
	return ClientStaker{
		TxID:        validator.TxID,
		StartTime:   uint64(validator.StartTime),
		EndTime:     uint64(validator.EndTime),
		Weight:      uint64(validator.Weight),
		StakeAmount: (*uint64)(validator.StakeAmount),
		NodeID:      validator.NodeID,
	}
}

func apiOwnerToClientOwner(rewardOwner *api.Owner) (*ClientOwner, error) {
	if rewardOwner == nil {
		return nil, nil
	}

	addrs, err := address.ParseToIDs(rewardOwner.Addresses)
	return &ClientOwner{
		Locktime:  uint64(rewardOwner.Locktime),
		Threshold: uint32(rewardOwner.Threshold),
		Addresses: addrs,
	}, err
}

func getClientPermissionlessValidators(validatorsSliceIntf []interface{}) ([]ClientPermissionlessValidator, error) {
	clientValidators := make([]ClientPermissionlessValidator, len(validatorsSliceIntf))
	for i, validatorMapIntf := range validatorsSliceIntf {
		validatorMapJSON, err := json.Marshal(validatorMapIntf)
		if err != nil {
			return nil, err
		}

		var apiValidator api.PermissionlessValidator
		err = json.Unmarshal(validatorMapJSON, &apiValidator)
		if err != nil {
			return nil, err
		}

		validationRewardOwner, err := apiOwnerToClientOwner(apiValidator.ValidationRewardOwner)
		if err != nil {
			return nil, err
		}

		clientValidators[i] = ClientPermissionlessValidator{
			ClientStaker:          apiStakerToClientStaker(apiValidator.Staker),
			ValidationRewardOwner: validationRewardOwner,
			PotentialReward:       (*uint64)(apiValidator.PotentialReward),
			Uptime:                (*float32)(apiValidator.Uptime),
			Connected:             &apiValidator.Connected,
			Signer:                apiValidator.Signer,
		}
	}
	return clientValidators, nil
}
