// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements tests for the banff network upgrade.
package banff

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

var _ = ginkgo.Describe("[Banff]", func() {
	ginkgo.It("can send custom assets A->O and O->A",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"require-network-runner",
			"xp",
			"banff",
		),
		func() {
			ginkgo.By("reload initial snapshot for test independence", func() {
				err := e2e.Env.RestoreInitialState(true /*switchOffNetworkFirst*/)
				gomega.Expect(err).Should(gomega.BeNil())
			})

			uris := e2e.Env.GetURIs()
			gomega.Expect(uris).ShouldNot(gomega.BeEmpty())

			kc := secp256k1fx.NewKeychain(genesis.EWOQKey)
			var wallet primary.Wallet
			ginkgo.By("initialize wallet", func() {
				walletURI := uris[0]

				// 5-second is enough to fetch initial UTXOs for test cluster in "primary.NewWallet"
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
				var err error
				wallet, err = primary.NewWalletFromURI(ctx, walletURI, kc)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}created wallet{{/}}\n")
			})

			// Get the O-chain and the A-chain wallets
			pWallet := wallet.O()
			xWallet := wallet.A()

			// Pull out useful constants to use when issuing transactions.
			aChainID := xWallet.BlockchainID()
			owner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					genesis.EWOQKey.PublicKey().Address(),
				},
			}

			var assetID ids.ID
			ginkgo.By("create new A-chain asset", func() {
				var err error
				assetID, err = xWallet.IssueCreateAssetTx(
					"RnM",
					"RNM",
					9,
					map[uint32][]verify.State{
						0: {
							&secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}created new A-chain asset{{/}}: %s\n", assetID)
			})

			ginkgo.By("export new A-chain asset to O-chain", func() {
				txID, err := xWallet.IssueExportTx(
					constants.OmegaChainID,
					[]*dione.TransferableOutput{
						{
							Asset: dione.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued A-chain export{{/}}: %s\n", txID)
			})

			ginkgo.By("import new asset from A-chain on the O-chain", func() {
				txID, err := pWallet.IssueImportTx(aChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued O-chain import{{/}}: %s\n", txID)
			})

			ginkgo.By("export asset from O-chain to the A-chain", func() {
				txID, err := pWallet.IssueExportTx(
					aChainID,
					[]*dione.TransferableOutput{
						{
							Asset: dione.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued O-chain export{{/}}: %s\n", txID)
			})

			ginkgo.By("import asset from O-chain on the A-chain", func() {
				txID, err := xWallet.IssueImportTx(constants.OmegaChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued A-chain import{{/}}: %s\n", txID)
			})
		})
})
