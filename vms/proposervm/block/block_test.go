// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/units"
)

func equal(require *require.Assertions, chainID ids.ID, want, have SignedBlock) {
	require.Equal(want.ID(), have.ID())
	require.Equal(want.ParentID(), have.ParentID())
	require.Equal(want.OChainHeight(), have.OChainHeight())
	require.Equal(want.Timestamp(), have.Timestamp())
	require.Equal(want.Block(), have.Block())
	require.Equal(want.Proposer(), have.Proposer())
	require.Equal(want.Bytes(), have.Bytes())
	require.Equal(want.Verify(false, chainID), have.Verify(false, chainID))
	require.Equal(want.Verify(true, chainID), have.Verify(true, chainID))
}

func TestVerifyNoCertWithSignature(t *testing.T) {
	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	oChainHeight := uint64(2)
	innerBlockBytes := []byte{3}

	require := require.New(t)

	builtBlockIntf, err := BuildUnsigned(parentID, timestamp, oChainHeight, innerBlockBytes)
	require.NoError(err)

	builtBlock := builtBlockIntf.(*statelessBlock)
	builtBlock.Signature = []byte{0}

	err = builtBlock.Verify(false, ids.Empty)
	require.ErrorIs(err, errUnexpectedProposer)

	err = builtBlock.Verify(true, ids.Empty)
	require.ErrorIs(err, errMissingProposer)
}

func TestBlockSizeLimit(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	oChainHeight := uint64(2)
	innerBlockBytes := bytes.Repeat([]byte{0}, 270*units.KiB)

	// with the large limit, it should be able to build large blocks
	_, err := BuildUnsigned(parentID, timestamp, oChainHeight, innerBlockBytes)
	require.NoError(err)
}
