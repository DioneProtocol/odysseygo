// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"crypto"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/staking"
)

func TestBuild(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	oChainHeight := uint64(2)
	innerBlockBytes := []byte{3}
	chainID := ids.ID{4}

	tlsCert, err := staking.NewTLSCert()
	require.NoError(err)

	cert := staking.CertificateFromX509(tlsCert.Leaf)
	key := tlsCert.PrivateKey.(crypto.Signer)

	builtBlock, err := Build(
		parentID,
		timestamp,
		oChainHeight,
		cert,
		innerBlockBytes,
		chainID,
		key,
	)
	require.NoError(err)

	require.Equal(parentID, builtBlock.ParentID())
	require.Equal(oChainHeight, builtBlock.OChainHeight())
	require.Equal(timestamp, builtBlock.Timestamp())
	require.Equal(innerBlockBytes, builtBlock.Block())

	require.NoError(builtBlock.Verify(true, chainID))

	err = builtBlock.Verify(false, chainID)
	require.ErrorIs(err, errUnexpectedProposer)
}

func TestBuildUnsigned(t *testing.T) {
	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	oChainHeight := uint64(2)
	innerBlockBytes := []byte{3}

	require := require.New(t)

	builtBlock, err := BuildUnsigned(parentID, timestamp, oChainHeight, innerBlockBytes)
	require.NoError(err)

	require.Equal(parentID, builtBlock.ParentID())
	require.Equal(oChainHeight, builtBlock.OChainHeight())
	require.Equal(timestamp, builtBlock.Timestamp())
	require.Equal(innerBlockBytes, builtBlock.Block())
	require.Equal(ids.EmptyNodeID, builtBlock.Proposer())

	require.NoError(builtBlock.Verify(false, ids.Empty))

	err = builtBlock.Verify(true, ids.Empty)
	require.ErrorIs(err, errMissingProposer)
}

func TestBuildHeader(t *testing.T) {
	require := require.New(t)

	chainID := ids.ID{1}
	parentID := ids.ID{2}
	bodyID := ids.ID{3}

	builtHeader, err := BuildHeader(
		chainID,
		parentID,
		bodyID,
	)
	require.NoError(err)

	require.Equal(chainID, builtHeader.ChainID())
	require.Equal(parentID, builtHeader.ParentID())
	require.Equal(bodyID, builtHeader.BodyID())
}

func TestBuildOption(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	innerBlockBytes := []byte{3}

	builtOption, err := BuildOption(parentID, innerBlockBytes)
	require.NoError(err)

	require.Equal(parentID, builtOption.ParentID())
	require.Equal(innerBlockBytes, builtOption.Block())
}
