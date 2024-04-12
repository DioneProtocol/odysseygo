// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package o

import (
	"context"
	"fmt"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/tests/e2e"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/signer"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var _ = e2e.DescribeOChain("[Permissionless Subnets]", func() {
	ginkgo.It("subnets operations",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"xp",
			"permissionless-subnets",
		),
		func() {
			nodeURI := e2e.Env.GetRandomNodeURI()

			keychain := e2e.Env.NewKeychain(1)
			baseWallet := e2e.Env.NewWallet(keychain, nodeURI)

			oWallet := baseWallet.O()
			aWallet := baseWallet.A()
			aChainID := aWallet.BlockchainID()

			var validatorID ids.NodeID
			ginkgo.By("retrieving the node ID of a primary network validator", func() {
				oChainClient := omegavm.NewClient(nodeURI.URI)
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				validatorIDs, err := oChainClient.SampleValidators(ctx, constants.PrimaryNetworkID, 1)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(validatorIDs).Should(gomega.HaveLen(1))
				validatorID = validatorIDs[0]
			})

			owner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					keychain.Keys[0].Address(),
				},
			}

			var subnetID ids.ID
			ginkgo.By("create a permissioned subnet", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				subnetTx, err := oWallet.IssueCreateSubnetTx(
					owner,
					common.WithContext(ctx),
				)
				cancel()

				subnetID = subnetTx.ID()
				gomega.Expect(subnetID, err).Should(gomega.Not(gomega.Equal(constants.PrimaryNetworkID)))
			})

			var subnetAssetID ids.ID
			ginkgo.By("create a custom asset for the permissionless subnet", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				subnetAssetTx, err := aWallet.IssueCreateAssetTx(
					"RnM",
					"RNM",
					9,
					map[uint32][]verify.State{
						0: {
							&secp256k1fx.TransferOutput{
								Amt:          100 * units.MegaDione,
								OutputOwners: *owner,
							},
						},
					},
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
				subnetAssetID = subnetAssetTx.ID()
			})

			ginkgo.By(fmt.Sprintf("Send 100 MegaDione of asset %s to the O-chain", subnetAssetID), func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				_, err := aWallet.IssueExportTx(
					constants.OmegaChainID,
					[]*dione.TransferableOutput{
						{
							Asset: dione.Asset{
								ID: subnetAssetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.MegaDione,
								OutputOwners: *owner,
							},
						},
					},
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			ginkgo.By(fmt.Sprintf("Import the 100 MegaDione of asset %s from the A-chain into the O-chain", subnetAssetID), func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				_, err := oWallet.IssueImportTx(
					aChainID,
					owner,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			ginkgo.By("make subnet permissionless", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				_, err := oWallet.IssueTransformSubnetTx(
					subnetID,
					subnetAssetID,
					50*units.MegaDione,
					100*units.MegaDione,
					reward.PercentDenominator,
					reward.PercentDenominator,
					1,
					100*units.MegaDione,
					time.Second,
					365*24*time.Hour,
					time.Second,
					365*24*time.Hour,
					0,
					1,
					5,
					.80*reward.PercentDenominator,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			validatorStartTime := time.Now().Add(time.Minute)
			ginkgo.By("add permissionless validator", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				_, err := oWallet.IssueAddPermissionlessValidatorTx(
					&txs.SubnetValidator{
						Validator: txs.Validator{
							NodeID: validatorID,
							Start:  uint64(validatorStartTime.Unix()),
							End:    uint64(validatorStartTime.Add(5 * time.Second).Unix()),
							Wght:   25 * units.MegaDione,
						},
						Subnet: subnetID,
					},
					&signer.Empty{},
					subnetAssetID,
					&secp256k1fx.OutputOwners{},
					&secp256k1fx.OutputOwners{},
					reward.PercentDenominator,
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})

			delegatorStartTime := validatorStartTime
			ginkgo.By("add permissionless delegator", func() {
				ctx, cancel := context.WithTimeout(context.Background(), e2e.DefaultTimeout)
				_, err := oWallet.IssueAddPermissionlessDelegatorTx(
					&txs.SubnetValidator{
						Validator: txs.Validator{
							NodeID: validatorID,
							Start:  uint64(delegatorStartTime.Unix()),
							End:    uint64(delegatorStartTime.Add(5 * time.Second).Unix()),
							Wght:   25 * units.MegaDione,
						},
						Subnet: subnetID,
					},
					subnetAssetID,
					&secp256k1fx.OutputOwners{},
					common.WithContext(ctx),
				)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			})
		})
})
