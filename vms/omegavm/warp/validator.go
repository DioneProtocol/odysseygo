// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/validators"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/utils/math"
	"github.com/DioneProtocol/odysseygo/utils/set"
)

var (
	_ utils.Sortable[*Validator] = (*Validator)(nil)

	ErrUnknownValidator = errors.New("unknown validator")
	ErrWeightOverflow   = errors.New("weight overflowed")
)

// ValidatorState defines the functions that must be implemented to get
// the canonical validator set for warp message validation.
type ValidatorState interface {
	GetValidatorSet(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error)
}

type Validator struct {
	PublicKey      *bls.PublicKey
	PublicKeyBytes []byte
	Weight         uint64
	NodeIDs        []ids.NodeID
}

func (v *Validator) Less(o *Validator) bool {
	return bytes.Compare(v.PublicKeyBytes, o.PublicKeyBytes) < 0
}

// GetCanonicalValidatorSet returns the validator set of [subnetID] at
// [oChainHeight] in a canonical ordering. Also returns the total weight on
// [subnetID].
func GetCanonicalValidatorSet(
	ctx context.Context,
	oChainState ValidatorState,
	oChainHeight uint64,
	subnetID ids.ID,
) ([]*Validator, uint64, error) {
	// Get the validator set at the given height.
	vdrSet, err := oChainState.GetValidatorSet(ctx, oChainHeight, subnetID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch validator set (O-Chain Height: %d, SubnetID: %s): %w", oChainHeight, subnetID, err)
	}

	var (
		vdrs        = make(map[string]*Validator, len(vdrSet))
		totalWeight uint64
	)
	for _, vdr := range vdrSet {
		totalWeight, err = math.Add64(totalWeight, vdr.Weight)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %w", ErrWeightOverflow, err)
		}

		if vdr.PublicKey == nil {
			continue
		}

		pkBytes := vdr.PublicKey.Serialize()
		uniqueVdr, ok := vdrs[string(pkBytes)]
		if !ok {
			uniqueVdr = &Validator{
				PublicKey:      vdr.PublicKey,
				PublicKeyBytes: pkBytes,
			}
			vdrs[string(pkBytes)] = uniqueVdr
		}

		uniqueVdr.Weight += vdr.Weight // Impossible to overflow here
		uniqueVdr.NodeIDs = append(uniqueVdr.NodeIDs, vdr.NodeID)
	}

	// Sort validators by public key
	vdrList := maps.Values(vdrs)
	utils.Sort(vdrList)
	return vdrList, totalWeight, nil
}

// FilterValidators returns the validators in [vdrs] whose bit is set to 1 in
// [indices].
//
// Returns an error if [indices] references an unknown validator.
func FilterValidators(
	indices set.Bits,
	vdrs []*Validator,
) ([]*Validator, error) {
	// Verify that all alleged signers exist
	if indices.BitLen() > len(vdrs) {
		return nil, fmt.Errorf(
			"%w: NumIndices (%d) >= NumFilteredValidators (%d)",
			ErrUnknownValidator,
			indices.BitLen()-1, // -1 to convert from length to index
			len(vdrs),
		)
	}

	filteredVdrs := make([]*Validator, 0, len(vdrs))
	for i, vdr := range vdrs {
		if !indices.Contains(i) {
			continue
		}

		filteredVdrs = append(filteredVdrs, vdr)
	}
	return filteredVdrs, nil
}

// SumWeight returns the total weight of the provided validators.
func SumWeight(vdrs []*Validator) (uint64, error) {
	var (
		weight uint64
		err    error
	)
	for _, vdr := range vdrs {
		weight, err = math.Add64(weight, vdr.Weight)
		if err != nil {
			return 0, fmt.Errorf("%w: %w", ErrWeightOverflow, err)
		}
	}
	return weight, nil
}

// AggregatePublicKeys returns the public key of the provided validators.
//
// Invariant: All of the public keys in [vdrs] are valid.
func AggregatePublicKeys(vdrs []*Validator) (*bls.PublicKey, error) {
	pks := make([]*bls.PublicKey, len(vdrs))
	for i, vdr := range vdrs {
		pks[i] = vdr.PublicKey
	}
	return bls.AggregatePublicKeys(pks)
}
