// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/tests/fixture"
	"github.com/DioneProtocol/odysseygo/tests/fixture/testnet"
	"github.com/DioneProtocol/odysseygo/tests/fixture/testnet/local"

	// ensure test packages are scanned by ginkgo
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/banff"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/c"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/faultinjection"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/p"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/static-handlers"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/x"
	_ "github.com/DioneProtocol/odysseygo/tests/e2e/x/transfer"
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "e2e test suites")
}

var (
	odysseyGoExecPath    string
	persistentNetworkDir string
	usePersistentNetwork bool
)

func init() {
	flag.StringVar(
		&odysseyGoExecPath,
		"odysseygo-path",
		os.Getenv(local.OdysseyGoPathEnvName),
		fmt.Sprintf("odysseygo executable path (required if not using a persistent network). Also possible to configure via the %s env variable.", local.OdysseyGoPathEnvName),
	)
	flag.StringVar(
		&persistentNetworkDir,
		"network-dir",
		"",
		fmt.Sprintf("[optional] the dir containing the configuration of a persistent network to target for testing. Useful for speeding up test development. Also possible to configure via the %s env variable.", local.NetworkDirEnvName),
	)
	flag.BoolVar(
		&usePersistentNetwork,
		"use-persistent-network",
		false,
		"[optional] whether to target the persistent network identified by --network-dir.",
	)
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Run only once in the first ginkgo process

	require := require.New(ginkgo.GinkgoT())

	if usePersistentNetwork && len(persistentNetworkDir) == 0 {
		persistentNetworkDir = os.Getenv(local.NetworkDirEnvName)
	}

	// Load or create a test network
	var network *local.LocalNetwork
	if len(persistentNetworkDir) > 0 {
		tests.Outf("{{yellow}}Using a pre-existing network configured at %s{{/}}\n", persistentNetworkDir)

		var err error
		network, err = local.ReadNetwork(persistentNetworkDir)
		require.NoError(err)
	} else {
		tests.Outf("{{magenta}}Starting network with %q{{/}}\n", odysseyGoExecPath)

		ctx, cancel := context.WithTimeout(context.Background(), local.DefaultNetworkStartTimeout)
		defer cancel()
		var err error
		network, err = local.StartNetwork(
			ctx,
			ginkgo.GinkgoWriter,
			"", // Use the default path to ensure a predictable target for github's upload-artifact action
			&local.LocalNetwork{
				LocalConfig: local.LocalConfig{
					ExecPath: odysseyGoExecPath,
				},
			},
			testnet.DefaultNodeCount,
			testnet.DefaultFundedKeyCount,
		)
		require.NoError(err)
		ginkgo.DeferCleanup(func() {
			tests.Outf("Shutting down network\n")
			require.NoError(network.Stop())
		})

		tests.Outf("{{green}}Successfully started network{{/}}\n")
	}

	uris := network.GetURIs()
	require.NotEmpty(uris, "network contains no nodes")
	tests.Outf("{{green}}network URIs: {{/}} %+v\n", uris)

	testDataServerURI, err := fixture.ServeTestData(fixture.TestData{
		FundedKeys: network.FundedKeys,
	})
	tests.Outf("{{green}}test data server URI: {{/}} %+v\n", testDataServerURI)
	require.NoError(err)

	env := &e2e.TestEnvironment{
		NetworkDir:        network.Dir,
		URIs:              uris,
		TestDataServerURI: testDataServerURI,
	}
	bytes, err := json.Marshal(env)
	require.NoError(err)
	return bytes
}, func(envBytes []byte) {
	// Run in every ginkgo process

	// Initialize the local test environment from the global state
	e2e.InitTestEnvironment(envBytes)
})
