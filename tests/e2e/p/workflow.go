// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package p

import (
	"context"
	"errors"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/DioneProtocol/odysseygo/api/info"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/alpha"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/status"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

// OChainWorkflow is an integration test for normal O-Chain operations
// - Issues an Add Validator using the funding address
// - Exports DIONE from the O-Chain funding address to the A-Chain created address
// - Exports DIONE from the A-Chain created address to the O-Chain created address
// - Checks the expected value of the funding address

var _ = e2e.DescribeOChain("[Workflow]", func() {
	ginkgo.It("O-chain main operations",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"require-network-runner",
			"xp",
			"workflow",
		),
		ginkgo.FlakeAttempts(2),
		func() {
			rpcEps := e2e.Env.GetURIs()
			gomega.Expect(rpcEps).ShouldNot(gomega.BeEmpty())
			nodeURI := rpcEps[0]

			tests.Outf("{{blue}} setting up keys {{/}}\n")
			_, testKeyAddrs, keyChain := e2e.Env.GetTestKeys()

			tests.Outf("{{blue}} setting up wallet {{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
			baseWallet, err := primary.NewWalletFromURI(ctx, nodeURI, keyChain)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			pWallet := baseWallet.O()
			dioneAssetID := baseWallet.O().DIONEAssetID()
			xWallet := baseWallet.A()
			oChainClient := omegavm.NewClient(nodeURI)
			aChainClient := alpha.NewClient(nodeURI, xWallet.BlockchainID().String())

			tests.Outf("{{blue}} fetching minimal stake amounts {{/}}\n")
			ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
			minValStake, err := oChainClient.GetMinStake(ctx, constants.OmegaChainID)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
			tests.Outf("{{green}} minimal validator stake: %d {{/}}\n", minValStake)

			tests.Outf("{{blue}} fetching tx fee {{/}}\n")
			infoClient := info.NewClient(nodeURI)
			ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultWalletCreationTimeout)
			fees, err := infoClient.GetTxFee(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
			txFees := uint64(fees.TxFee)
			tests.Outf("{{green}} txFee: %d {{/}}\n", txFees)

			// amount to transfer from O to A chain
			toTransfer := 1 * units.Dione

			pShortAddr := testKeyAddrs[0]
			xTargetAddr := testKeyAddrs[1]
			ginkgo.By("check selected keys have sufficient funds", func() {
				pBalances, err := pWallet.Builder().GetBalance()
				pBalance := pBalances[dioneAssetID]
				minBalance := minValStake + txFees + txFees + toTransfer + txFees
				gomega.Expect(pBalance, err).To(gomega.BeNumerically(">=", minBalance))
			})
			// create validator data
			validatorStartTimeDiff := 30 * time.Second
			vdrStartTime := time.Now().Add(validatorStartTimeDiff)

			vdr := &txs.Validator{
				NodeID: ids.GenerateTestNodeID(),
				Start:  uint64(vdrStartTime.Unix()),
				End:    uint64(vdrStartTime.Add(72 * time.Hour).Unix()),
				Wght:   minValStake,
			}
			rewardOwner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{pShortAddr},
			}

			ginkgo.By("issue add validator tx", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				addValidatorTxID, err := pWallet.IssueAddValidatorTx(
					vdr,
					rewardOwner,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				txStatus, err := oChainClient.GetTxStatus(ctx, addValidatorTxID)
				cancel()
				gomega.Expect(txStatus.Status, err).To(gomega.Equal(status.Committed))
			})

			// retrieve initial balances
			pBalances, err := pWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			pStartBalance := pBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance before O->A export: %d {{/}}\n", pStartBalance)

			xBalances, err := xWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			xStartBalance := xBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance before O->A export: %d {{/}}\n", xStartBalance)

			outputOwner := secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					xTargetAddr,
				},
			}
			output := &secp256k1fx.TransferOutput{
				Amt:          toTransfer,
				OutputOwners: outputOwner,
			}

			ginkgo.By("export dione from O to A chain", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				exportTxID, err := pWallet.IssueExportTx(
					xWallet.BlockchainID(),
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

				ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				txStatus, err := oChainClient.GetTxStatus(ctx, exportTxID)
				cancel()
				gomega.Expect(txStatus.Status, err).To(gomega.Equal(status.Committed))
			})

			// check balances post export
			pBalances, err = pWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			pPreImportBalance := pBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance after O->A export: %d {{/}}\n", pPreImportBalance)

			xBalances, err = xWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			xPreImportBalance := xBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance after O->A export: %d {{/}}\n", xPreImportBalance)

			gomega.Expect(xPreImportBalance).To(gomega.Equal(xStartBalance)) // import not performed yet
			gomega.Expect(pPreImportBalance).To(gomega.Equal(pStartBalance - toTransfer - txFees))

			ginkgo.By("import dione from O into A chain", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				importTxID, err := xWallet.IssueImportTx(
					constants.OmegaChainID,
					&outputOwner,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil(), "is context.DeadlineExceeded: %v", errors.Is(err, context.DeadlineExceeded))

				ctx, cancel = context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				txStatus, err := aChainClient.GetTxStatus(ctx, importTxID)
				cancel()
				gomega.Expect(txStatus, err).To(gomega.Equal(choices.Accepted))
			})

			// check balances post import
			pBalances, err = pWallet.Builder().GetBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			pFinalBalance := pBalances[dioneAssetID]
			tests.Outf("{{blue}} O-chain balance after O->A import: %d {{/}}\n", pFinalBalance)

			xBalances, err = xWallet.Builder().GetFTBalance()
			gomega.Expect(err).Should(gomega.BeNil())
			xFinalBalance := xBalances[dioneAssetID]
			tests.Outf("{{blue}} A-chain balance after O->A import: %d {{/}}\n", xFinalBalance)

			gomega.Expect(xFinalBalance).To(gomega.Equal(xPreImportBalance + toTransfer - txFees)) // import not performed yet
			gomega.Expect(pFinalBalance).To(gomega.Equal(pPreImportBalance))
		})
})
