// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
)

func (b *Backend) GetAccumulatedFees(key []byte) (map[ids.ID][]byte, error) {
	var accumulatedFees map[ids.ID][]byte
	for _, chain := range b.AccumulatedFeeChainIDs {
		value, err := b.Ctx.SharedMemory.GetBigInt(chain, key)
		if err != nil {
			return nil, err
		}
		switch value.Sign() {
		case -1:
			return nil, fmt.Errorf("negative accumulated fees")
		case 0:
			continue
		}
		if len(accumulatedFees) == 0 {
			accumulatedFees = make(map[ids.ID][]byte)
		}
		accumulatedFees[chain] = value.Bytes()
	}
	return accumulatedFees, nil
}
