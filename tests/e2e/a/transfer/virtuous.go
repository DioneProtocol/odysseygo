// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements A-chain transfer tests.
package transfer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/onsi/gomega"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/snow/choices"
	"github.com/DioneProtocol/odysseygo/tests"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/alpha"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

const (
	totalRounds = 50

	metricBlksProcessing = "odyssey_A_blks_processing"
	metricBlksAccepted   = "odyssey_A_blks_accepted_count"
)

// This test requires that the network not have ongoing blocks and
// cannot reliably be run in parallel.
var _ = e2e.DescribeAChainSerial("[Virtuous Transfer Tx DIONE]", func() {
	ginkgo.It("can issue a virtuous transfer tx for DIONE asset",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"a",
			"virtuous-transfer-tx-dione",
		),
		func() {
			rpcEps := make([]string, len(e2e.Env.URIs))
			for i, nodeURI := range e2e.Env.URIs {
				rpcEps[i] = nodeURI.URI
			}

			// Waiting for ongoing blocks to have completed before starting this
			// test avoids the case of a previous test having initiated block
			// processing but not having completed it.
			gomega.Eventually(func() bool {
				allNodeMetrics, err := tests.GetNodesMetrics(rpcEps, metricBlksProcessing)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, metrics := range allNodeMetrics {
					if metrics[metricBlksProcessing] > 0 {
						return false
					}
				}
				return true
			}).
				WithTimeout(e2e.DefaultTimeout).
				WithPolling(e2e.DefaultPollingInterval).
				Should(gomega.BeTrue(), "The cluster is generating ongoing blocks. Is this test being run in parallel?")

			allMetrics := []string{
				metricBlksProcessing,
				metricBlksAccepted,
			}

			// Ensure the same set of 10 keys is used for all tests
			// by retrieving them outside of runFunc.
			testKeys := e2e.Env.AllocateFundedKeys(10)

			runFunc := func(round int) {
				tests.Outf("{{green}}\n\n\n\n\n\n---\n[ROUND #%02d]:{{/}}\n", round)

				needPermute := round > 3
				if needPermute {
					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(testKeys), func(i, j int) {
						testKeys[i], testKeys[j] = testKeys[j], testKeys[i]
					})
				}

				keychain := secp256k1fx.NewKeychain(testKeys...)
				baseWallet := e2e.Env.NewWallet(keychain, e2e.Env.GetRandomNodeURI())
				dioneAssetID := baseWallet.A().DIONEAssetID()

				wallets := make([]primary.Wallet, len(testKeys))
				shortAddrs := make([]ids.ShortID, len(testKeys))
				for i := range wallets {
					shortAddrs[i] = testKeys[i].PublicKey().Address()

					wallets[i] = primary.NewWalletWithOptions(
						baseWallet,
						common.WithCustomAddresses(set.Of(
							testKeys[i].PublicKey().Address(),
						)),
					)
				}

				metricsBeforeTx, err := tests.GetNodesMetrics(rpcEps, allMetrics...)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, uri := range rpcEps {
					tests.Outf("{{green}}metrics at %q:{{/}} %v\n", uri, metricsBeforeTx[uri])
				}

				testBalances := make([]uint64, 0)
				for i, w := range wallets {
					balances, err := w.A().Builder().GetFTBalance()
					gomega.Expect(err).Should(gomega.BeNil())

					bal := balances[dioneAssetID]
					testBalances = append(testBalances, bal)

					fmt.Printf(`CURRENT BALANCE %21d DIONE (SHORT ADDRESS %q)
`,
						bal,
						testKeys[i].PublicKey().Address(),
					)
				}
				fromIdx := -1
				for i := range testBalances {
					if fromIdx < 0 && testBalances[i] > 0 {
						fromIdx = i
						break
					}
				}
				if fromIdx < 0 {
					gomega.Expect(fromIdx).Should(gomega.BeNumerically(">", 0), "no address found with non-zero balance")
				}

				toIdx := -1
				for i := range testBalances {
					// prioritize the address with zero balance
					if toIdx < 0 && i != fromIdx && testBalances[i] == 0 {
						toIdx = i
						break
					}
				}
				if toIdx < 0 {
					// no zero balance address, so just transfer between any two addresses
					toIdx = (fromIdx + 1) % len(testBalances)
				}

				senderOrigBal := testBalances[fromIdx]
				receiverOrigBal := testBalances[toIdx]

				amountToTransfer := senderOrigBal / 10

				senderNewBal := senderOrigBal - amountToTransfer - baseWallet.A().BaseTxFee()
				receiverNewBal := receiverOrigBal + amountToTransfer

				ginkgo.By("A-Chain transfer with wrong amount must fail", func() {
					ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
					_, err := wallets[fromIdx].A().IssueBaseTx(
						[]*dione.TransferableOutput{{
							Asset: dione.Asset{
								ID: dioneAssetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt: senderOrigBal + 1,
								OutputOwners: secp256k1fx.OutputOwners{
									Threshold: 1,
									Addrs:     []ids.ShortID{shortAddrs[toIdx]},
								},
							},
						}},
						common.WithContext(ctx),
					)
					cancel()
					gomega.Expect(err.Error()).Should(gomega.ContainSubstring("insufficient funds"))
				})

				fmt.Printf(`===
TRANSFERRING

FROM [%q]
SENDER    CURRENT BALANCE     : %21d DIONE
SENDER    NEW BALANCE (AFTER) : %21d DIONE

TRANSFER AMOUNT FROM SENDER   : %21d DIONE

TO [%q]
RECEIVER  CURRENT BALANCE     : %21d DIONE
RECEIVER  NEW BALANCE (AFTER) : %21d DIONE
===
`,
					shortAddrs[fromIdx],
					senderOrigBal,
					senderNewBal,
					amountToTransfer,
					shortAddrs[toIdx],
					receiverOrigBal,
					receiverNewBal,
				)

				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
				tx, err := wallets[fromIdx].A().IssueBaseTx(
					[]*dione.TransferableOutput{{
						Asset: dione.Asset{
							ID: dioneAssetID,
						},
						Out: &secp256k1fx.TransferOutput{
							Amt: amountToTransfer,
							OutputOwners: secp256k1fx.OutputOwners{
								Threshold: 1,
								Addrs:     []ids.ShortID{shortAddrs[toIdx]},
							},
						},
					}},
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				balances, err := wallets[fromIdx].A().Builder().GetFTBalance()
				gomega.Expect(err).Should(gomega.BeNil())
				senderCurBalX := balances[dioneAssetID]
				tests.Outf("{{green}}first wallet balance:{{/}}  %d\n", senderCurBalX)

				balances, err = wallets[toIdx].A().Builder().GetFTBalance()
				gomega.Expect(err).Should(gomega.BeNil())
				receiverCurBalX := balances[dioneAssetID]
				tests.Outf("{{green}}second wallet balance:{{/}} %d\n", receiverCurBalX)

				gomega.Expect(senderCurBalX).Should(gomega.Equal(senderNewBal))
				gomega.Expect(receiverCurBalX).Should(gomega.Equal(receiverNewBal))

				txID := tx.ID()
				for _, u := range rpcEps {
					xc := alpha.NewClient(u, "A")
					ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
					status, err := xc.ConfirmTx(ctx, txID, 2*time.Second)
					cancel()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(status).Should(gomega.Equal(choices.Accepted))
				}

				for _, u := range rpcEps {
					xc := alpha.NewClient(u, "A")
					ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultConfirmTxTimeout)
					status, err := xc.ConfirmTx(ctx, txID, 2*time.Second)
					cancel()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(status).Should(gomega.Equal(choices.Accepted))

					mm, err := tests.GetNodeMetrics(u, allMetrics...)
					gomega.Expect(err).Should(gomega.BeNil())

					prev := metricsBeforeTx[u]

					// +0 since A-chain tx must have been processed and accepted
					// by now
					gomega.Expect(mm[metricBlksProcessing]).Should(gomega.Equal(prev[metricBlksProcessing]))

					// +1 since A-chain tx must have been accepted by now
					gomega.Expect(mm[metricBlksAccepted]).Should(gomega.Equal(prev[metricBlksAccepted] + 1))

					metricsBeforeTx[u] = mm
				}
			}

			for i := 0; i < totalRounds; i++ {
				runFunc(i)
				time.Sleep(time.Second)
			}
		})
})
