// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package o

import (
	"context"
	"errors"
	"time"

	"github.com/DioneProtocol/odysseygo/vms/omegavm"
	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/DioneProtocol/odysseygo/api/info"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

// OChainWorkflow is an integration test for normal O-Chain operations
// - Issues an Add Validator and an Add Delegator using the funding address
// - Exports DIONE from the O-Chain funding address to the A-Chain created address
// - Exports DIONE from the A-Chain created address to the O-Chain created address
// - Checks the expected value of the funding address

var _ = e2e.DescribeOChain("[Workflow]", func() {
	ginkgo.It("O-chain main operations",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"xp",
			"workflow",
		),
		ginkgo.FlakeAttempts(2),
		func() {
			nodeURI := e2e.Env.GetRandomNodeURI()
			keychain := e2e.Env.NewKeychain(2)
			baseWallet := e2e.Env.NewWallet(keychain, nodeURI)

			oWallet := baseWallet.O()
			dioneAssetID := baseWallet.O().DIONEAssetID()
			aWallet := baseWallet.A()
			oChainClient := omegavm.NewClient(nodeURI.URI)

			tests.Outf("{{blue}} fetching minimal stake amounts {{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
			minValStake, minDelStake, err := oChainClient.GetMinStake(ctx, constants.OmegaChainID)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
			tests.Outf("{{green}} minimal validator stake: %d {{/}}\n", minValStake)
			tests.Outf("{{green}} minimal delegator stake: %d {{/}}\n", minDelStake)

			tests.Outf("{{blue}} fetching tx fee {{/}}\n")
			infoClient := info.NewClient(nodeURI.URI)
			ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
			fees, err := infoClient.GetTxFee(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
			txFees := uint64(fees.TxFee)
			tests.Outf("{{green}} txFee: %d {{/}}\n", txFees)

			// amount to transfer from O to A chain
			toTransfer := 1 * units.Dione

			oShortAddr := keychain.Keys[0].Address()
			aTargetAddr := keychain.Keys[1].Address()
			ginkgo.By("check selected keys have sufficient funds", func() {
				oBalances, err := oWallet.Builder().GetBalance()
				oBalance := oBalances[dioneAssetID]
				minBalance := minValStake + txFees + minDelStake + txFees + toTransfer + txFees
				gomega.Expect(oBalance, err).To(gomega.BeNumerically(">=", minBalance))
			})
			// create validator data
			validatorStartTimeDiff := 30 * time.Second
			vdrStartTime := time.Now().Add(validatorStartTimeDiff)

			// Use a random node ID to ensure that repeated test runs
			// will succeed against a persistent network.
			validatorID, err := ids.ToNodeID(utils.RandomBytes(ids.NodeIDLen))
			gomega.Expect(err).Should(gomega.BeNil())

			vdr := &txs.Validator{
				NodeID: validatorID,
				Start:  uint64(vdrStartTime.Unix()),
				End:    uint64(vdrStartTime.Add(72 * time.Hour).Unix()),
				Wght:   minValStake,
			}
			rewardOwner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{oShortAddr},
			}
			shares := uint32(20000) // TODO: retrieve programmatically

			ginkgo.By("issue add validator tx", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				_, err := oWallet.IssueAddValidatorTx(
					vdr,
					rewardOwner,
					shares,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			ginkgo.By("issue add delegator tx", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				_, err := oWallet.IssueAddDelegatorTx(
					vdr,
					rewardOwner,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			// retrieve initial balances
			oBalances, err := oWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			oStartBalance := oBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance before O->A export: %d {{/}}\n", oStartBalance)

			aBalances, err := aWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			aStartBalance := aBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance before O->A export: %d {{/}}\n", aStartBalance)

			outputOwner := secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					aTargetAddr,
				},
			}
			output := &secp256k1fx.TransferOutput{
				Amt:          toTransfer,
				OutputOwners: outputOwner,
			}

			ginkgo.By("export dione from O to A chain", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				_, err := oWallet.IssueExportTx(
					aWallet.BlockchainID(),
					[]*dione.TransferableOutput{
						{
							Asset: dione.Asset{
								ID: dioneAssetID,
							},
							Out: output,
						},
					},
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			// check balances post export
			oBalances, err = oWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			oPreImportBalance := oBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance after O->A export: %d {{/}}\n", oPreImportBalance)

			aBalances, err = aWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			aPreImportBalance := aBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance after O->A export: %d {{/}}\n", aPreImportBalance)

			gomega.Expect(aPreImportBalance).To(gomega.Equal(aStartBalance)) // import not performed yet
			gomega.Expect(oPreImportBalance).To(gomega.Equal(oStartBalance - toTransfer - txFees))

			ginkgo.By("import dione from O into A chain", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				_, err := aWallet.IssueImportTx(
					constants.OmegaChainID,
					&outputOwner,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil(), "is context.DeadlineExceeded: %v", errors.Is(err, context.DeadlineExceeded))
			})

			// check balances post import
			oBalances, err = oWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			oFinalBalance := oBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance after O->A import: %d {{/}}\n", oFinalBalance)

			aBalances, err = aWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			aFinalBalance := aBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance after O->A import: %d {{/}}\n", aFinalBalance)

			gomega.Expect(aFinalBalance).To(gomega.Equal(aPreImportBalance + toTransfer - txFees)) // import not performed yet
			gomega.Expect(oFinalBalance).To(gomega.Equal(oPreImportBalance))
		})
})
