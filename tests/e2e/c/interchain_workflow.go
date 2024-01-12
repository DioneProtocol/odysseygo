// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package c

import (
	"math/big"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/coreth/core/types"
	"github.com/DioneProtocol/coreth/plugin/evm"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var _ = e2e.DescribeCChain("[Interchain Workflow]", func() {
	require := require.New(ginkgo.GinkgoT())

	const (
		txAmount = 10 * units.Dione // Arbitrary amount to send and transfer
		gasLimit = uint64(21000)    // Standard gas limit
	)

	ginkgo.It("should ensure that funds can be transferred from the C-Chain to the X-Chain and the O-Chain", func() {
		ginkgo.By("initializing a new eth client")
		// Select a random node URI to use for both the eth client and
		// the wallet to avoid having to verify that all nodes are at
		// the same height before initializing the wallet.
		nodeURI := e2e.Env.GetRandomNodeURI()
		ethClient := e2e.Env.NewEthClient(nodeURI)

		ginkgo.By("allocating a pre-funded key to send from and a recipient key to deliver to")
		senderKey := e2e.Env.AllocateFundedKey()
		senderEthAddress := evm.GetEthAddress(senderKey)
		factory := secp256k1.Factory{}
		recipientKey, err := factory.NewPrivateKey()
		require.NoError(err)
		recipientEthAddress := evm.GetEthAddress(recipientKey)

		ginkgo.By("sending funds from one address to another on the C-Chain", func() {
			// Create transaction
			acceptedNonce, err := ethClient.AcceptedNonceAt(e2e.DefaultContext(), senderEthAddress)
			require.NoError(err)
			gasPrice := e2e.SuggestGasPrice(ethClient)
			tx := types.NewTransaction(
				acceptedNonce,
				recipientEthAddress,
				big.NewInt(int64(txAmount)),
				gasLimit,
				gasPrice,
				nil,
			)

			// Sign transaction
			cChainID, err := ethClient.ChainID(e2e.DefaultContext())
			require.NoError(err)
			signer := types.NewEIP155Signer(cChainID)
			signedTx, err := types.SignTx(tx, signer, senderKey.ToECDSA())
			require.NoError(err)

			_ = e2e.SendEthTransaction(ethClient, signedTx)

			ginkgo.By("waiting for the C-Chain recipient address to have received the sent funds")
			e2e.Eventually(func() bool {
				balance, err := ethClient.BalanceAt(e2e.DefaultContext(), recipientEthAddress, nil)
				require.NoError(err)
				return balance.Cmp(big.NewInt(0)) > 0
			}, e2e.DefaultTimeout, e2e.DefaultPollingInterval, "failed to see funds delivered before timeout")
		})

		// Wallet must be initialized after sending funds on the
		// C-Chain with the same node URI to ensure wallet state
		// matches on-chain state.
		ginkgo.By("initializing a keychain and associated wallet")
		keychain := secp256k1fx.NewKeychain(senderKey, recipientKey)
		baseWallet := e2e.Env.NewWallet(keychain, nodeURI)
		xWallet := baseWallet.X()
		cWallet := baseWallet.C()
		oWallet := baseWallet.O()

		ginkgo.By("defining common configuration")
		dioneAssetID := xWallet.DIONEAssetID()
		// Use the same owner for import funds to X-Chain and O-Chain
		recipientOwner := secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				recipientKey.Address(),
			},
		}
		// Use the same outputs for both X-Chain and O-Chain exports
		exportOutputs := []*secp256k1fx.TransferOutput{
			{
				Amt: txAmount,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs: []ids.ShortID{
						keychain.Keys[0].Address(),
					},
				},
			},
		}

		ginkgo.By("exporting DIONE from the C-Chain to the X-Chain", func() {
			_, err := cWallet.IssueExportTx(
				xWallet.BlockchainID(),
				exportOutputs,
				e2e.WithDefaultContext(),
				e2e.WithSuggestedGasPrice(ethClient),
			)
			require.NoError(err)
		})

		ginkgo.By("importing DIONE from the C-Chain to the X-Chain", func() {
			_, err := xWallet.IssueImportTx(
				cWallet.BlockchainID(),
				&recipientOwner,
				e2e.WithDefaultContext(),
			)
			require.NoError(err)
		})

		ginkgo.By("checking that the recipient address has received imported funds on the X-Chain", func() {
			balances, err := xWallet.Builder().GetFTBalance(common.WithCustomAddresses(set.Of(
				recipientKey.Address(),
			)))
			require.NoError(err)
			require.Positive(balances[dioneAssetID])
		})

		ginkgo.By("exporting DIONE from the C-Chain to the O-Chain", func() {
			_, err := cWallet.IssueExportTx(
				constants.OmegaChainID,
				exportOutputs,
				e2e.WithDefaultContext(),
				e2e.WithSuggestedGasPrice(ethClient),
			)
			require.NoError(err)
		})

		ginkgo.By("importing DIONE from the C-Chain to the O-Chain", func() {
			_, err = oWallet.IssueImportTx(
				cWallet.BlockchainID(),
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
