// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/utils/crypto/bls"
)

func TestSigner(t *testing.T) {
	for _, test := range SignerTests {
		sk, err := bls.NewSecretKey()
		require.NoError(t, err)

		chainID := ids.GenerateTestID()
		s := NewSigner(sk, chainID)

		test(t, s, sk, chainID)
	}
}
