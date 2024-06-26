// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements tests for the banff network upgrade.
package banff

import (
	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var _ = ginkgo.Describe("[Banff]", func() {
	ginkgo.It("can send custom assets A->O and O->A",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"xp",
			"banff",
		),
		func() {
			keychain := e2e.Env.NewKeychain(1)
			wallet := e2e.Env.NewWallet(keychain, e2e.Env.GetRandomNodeURI())

			// Get the O-chain and the A-chain wallets
			oWallet := wallet.O()
			aWallet := wallet.A()

			// Pull out useful constants to use when issuing transactions.
			aChainID := aWallet.BlockchainID()
			owner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					keychain.Keys[0].Address(),
				},
			}

			var assetID ids.ID
			ginkgo.By("create new A-chain asset", func() {
				assetTx, err := aWallet.IssueCreateAssetTx(
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
				assetID = assetTx.ID()

				tests.Outf("{{green}}created new A-chain asset{{/}}: %s\n", assetID)
			})

			ginkgo.By("export new A-chain asset to O-chain", func() {
				tx, err := aWallet.IssueExportTx(
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

				tests.Outf("{{green}}issued A-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import new asset from A-chain on the O-chain", func() {
				tx, err := oWallet.IssueImportTx(aChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued O-chain import{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("export asset from O-chain to the A-chain", func() {
				tx, err := oWallet.IssueExportTx(
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

				tests.Outf("{{green}}issued O-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import asset from O-chain on the A-chain", func() {
				tx, err := aWallet.IssueImportTx(constants.OmegaChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued A-chain import{{/}}: %s\n", tx.ID())
			})
		})
})
