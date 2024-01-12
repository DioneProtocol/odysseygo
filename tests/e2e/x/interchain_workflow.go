// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"math/big"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/coreth/plugin/delta"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var _ = e2e.DescribeXChain("[Interchain Workflow]", func() {
	require := require.New(ginkgo.GinkgoT())

	const transferAmount = 10 * units.Dione

	ginkgo.It("should ensure that funds can be transferred from the X-Chain to the C-Chain and the O-Chain", func() {
		nodeURI := e2e.Env.GetRandomNodeURI()

		ginkgo.By("creating wallet with a funded key to send from and recipient key to deliver to")
		factory := secp256k1.Factory{}
		recipientKey, err := factory.NewPrivateKey()
		require.NoError(err)
		keychain := e2e.Env.NewKeychain(1)
		keychain.Add(recipientKey)
		baseWallet := e2e.Env.NewWallet(keychain, nodeURI)
		xWallet := baseWallet.X()
		cWallet := baseWallet.C()
		oWallet := baseWallet.O()

		ginkgo.By("defining common configuration")
		recipientEthAddress := delta.GetEthAddress(recipientKey)
		dioneAssetID := xWallet.DIONEAssetID()
		// Use the same owner for sending to X-Chain and importing funds to O-Chain
		recipientOwner := secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				recipientKey.Address(),
			},
		}
		// Use the same outputs for both C-Chain and O-Chain exports
		exportOutputs := []*dione.TransferableOutput{
			{
				Asset: dione.Asset{
					ID: dioneAssetID,
				},
				Out: &secp256k1fx.TransferOutput{
					Amt: transferAmount,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs: []ids.ShortID{
							keychain.Keys[0].Address(),
						},
					},
				},
			},
		}

		ginkgo.By("sending funds from one address to another on the X-Chain", func() {
			_, err = xWallet.IssueBaseTx(
				[]*dione.TransferableOutput{{
					Asset: dione.Asset{
						ID: dioneAssetID,
					},
					Out: &secp256k1fx.TransferOutput{
						Amt:          transferAmount,
						OutputOwners: recipientOwner,
					},
				}},
				e2e.WithDefaultContext(),
			)
			require.NoError(err)
		})

		ginkgo.By("checking that the X-Chain recipient address has received the sent funds", func() {
			balances, err := xWallet.Builder().GetFTBalance(common.WithCustomAddresses(set.Of(
				recipientKey.Address(),
			)))
			require.NoError(err)
			require.Greater(balances[dioneAssetID], uint64(0))
		})

		ginkgo.By("exporting DIONE from the X-Chain to the C-Chain", func() {
			_, err := xWallet.IssueExportTx(
				cWallet.BlockchainID(),
				exportOutputs,
				e2e.WithDefaultContext(),
			)
			require.NoError(err)
		})

		ginkgo.By("initializing a new eth client")
		ethClient := e2e.Env.NewEthClient(nodeURI)

		ginkgo.By("importing DIONE from the X-Chain to the C-Chain", func() {
			_, err := cWallet.IssueImportTx(
				xWallet.BlockchainID(),
				recipientEthAddress,
				e2e.WithDefaultContext(),
				e2e.WithSuggestedGasPrice(ethClient),
			)
			require.NoError(err)
		})

		ginkgo.By("checking that the recipient address has received imported funds on the C-Chain")
		e2e.Eventually(func() bool {
			balance, err := ethClient.BalanceAt(e2e.DefaultContext(), recipientEthAddress, nil)
			require.NoError(err)
			return balance.Cmp(big.NewInt(0)) > 0
		}, e2e.DefaultTimeout, e2e.DefaultPollingInterval, "failed to see recipient address funded before timeout")

		ginkgo.By("exporting DIONE from the X-Chain to the O-Chain", func() {
			_, err := xWallet.IssueExportTx(
				constants.OmegaChainID,
				exportOutputs,
				e2e.WithDefaultContext(),
			)
			require.NoError(err)
		})

		ginkgo.By("importing DIONE from the X-Chain to the O-Chain", func() {
			_, err := oWallet.IssueImportTx(
				xWallet.BlockchainID(),
				&recipientOwner,
				e2e.WithDefaultContext(),
			)
			require.NoError(err)
		})

		ginkgo.By("checking that the recipient address has received imported funds on the O-Chain", func() {
			balances, err := oWallet.Builder().GetBalance(common.WithCustomAddresses(set.Of(
				recipientKey.Address(),
			)))
			require.NoError(err)
			require.Greater(balances[dioneAssetID], uint64(0))
		})
	})
})
