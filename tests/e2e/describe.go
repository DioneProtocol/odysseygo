// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e

import (
	ginkgo "github.com/onsi/ginkgo/v2"
)

// DescribeAChain annotates the tests for A-Chain.
// Can run with any type of cluster (e.g., local, testnet, mainnet).
func DescribeAChain(text string, body func()) bool {
	return ginkgo.Describe("[A-Chain] "+text, body)
}

// DescribeOChain annotates the tests for O-Chain.
// Can run with any type of cluster (e.g., local, testnet, mainnet).
func DescribeOChain(text string, body func()) bool {
	return ginkgo.Describe("[O-Chain] "+text, body)
}
