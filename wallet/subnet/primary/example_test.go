// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package primary

import (
	"context"
	"log"
	"time"

	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/signer"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

func ExampleWallet() {
	ctx := context.Background()
	kc := secp256k1fx.NewKeychain(genesis.EWOQKey)

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [LocalAPIURI] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := MakeWallet(ctx, &WalletConfig{
		URI:           LocalAPIURI,
		DIONEKeychain: kc,
		EthKeychain:   kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet with: %s\n", err)
		return
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the O-chain and the X-chain wallets
	oWallet := wallet.O()
	xWallet := wallet.X()

	// Pull out useful constants to use when issuing transactions.
	xChainID := xWallet.BlockchainID()
	owner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			genesis.EWOQKey.PublicKey().Address(),
		},
	}

	// Create a custom asset to send to the O-chain.
	createAssetStartTime := time.Now()
	createAssetTx, err := xWallet.IssueCreateAssetTx(
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
	)
	if err != nil {
		log.Fatalf("failed to create new X-chain asset with: %s\n", err)
		return
	}
	createAssetTxID := createAssetTx.ID()
	log.Printf("created X-chain asset %s in %s\n", createAssetTxID, time.Since(createAssetStartTime))

	// Send 100 MegaDione to the O-chain.
	exportStartTime := time.Now()
	exportTx, err := xWallet.IssueExportTx(
		constants.OmegaChainID,
		[]*dione.TransferableOutput{
			{
				Asset: dione.Asset{
					ID: createAssetTxID,
				},
				Out: &secp256k1fx.TransferOutput{
					Amt:          100 * units.MegaDione,
					OutputOwners: *owner,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to issue X->O export transaction with: %s\n", err)
		return
	}
	exportTxID := exportTx.ID()
	log.Printf("issued X->O export %s in %s\n", exportTxID, time.Since(exportStartTime))

	// Import the 100 MegaDione from the X-chain into the P-chain.
	importStartTime := time.Now()
	importTx, err := oWallet.IssueImportTx(xChainID, owner)
	if err != nil {
		log.Fatalf("failed to issue X->O import transaction with: %s\n", err)
		return
	}
	importTxID := importTx.ID()
	log.Printf("issued X->O import %s in %s\n", importTxID, time.Since(importStartTime))

	createSubnetStartTime := time.Now()
	createSubnetTx, err := oWallet.IssueCreateSubnetTx(owner)
	if err != nil {
		log.Fatalf("failed to issue create subnet transaction with: %s\n", err)
		return
	}
	createSubnetTxID := createSubnetTx.ID()
	log.Printf("issued create subnet transaction %s in %s\n", createSubnetTxID, time.Since(createSubnetStartTime))

	transformSubnetStartTime := time.Now()
	transformSubnetTx, err := oWallet.IssueTransformSubnetTx(
		createSubnetTxID,
		createAssetTxID,
		50*units.MegaDione,
		100*units.MegaDione,
		reward.PercentDenominator,
		reward.PercentDenominator,
		1,
		100*units.MegaDione,
		time.Second,
		365*24*time.Hour,
		0,
		1,
		5,
		.80*reward.PercentDenominator,
	)
	if err != nil {
		log.Fatalf("failed to issue transform subnet transaction with: %s\n", err)
		return
	}
	transformSubnetTxID := transformSubnetTx.ID()
	log.Printf("issued transform subnet transaction %s in %s\n", transformSubnetTxID, time.Since(transformSubnetStartTime))

	addPermissionlessValidatorStartTime := time.Now()
	startTime := time.Now().Add(time.Minute)
	addSubnetValidatorTx, err := oWallet.IssueAddPermissionlessValidatorTx(
		&txs.SubnetValidator{
			Validator: txs.Validator{
				NodeID: genesis.LocalConfig.InitialStakers[0].NodeID,
				Start:  uint64(startTime.Unix()),
				End:    uint64(startTime.Add(5 * time.Second).Unix()),
				Wght:   25 * units.MegaDione,
			},
			Subnet: createSubnetTxID,
		},
		&signer.Empty{},
		createAssetTx.ID(),
		&secp256k1fx.OutputOwners{},
		&secp256k1fx.OutputOwners{},
		reward.PercentDenominator,
	)
	if err != nil {
		log.Fatalf("failed to issue add subnet validator with: %s\n", err)
		return
	}
	addSubnetValidatorTxID := addSubnetValidatorTx.ID()
	log.Printf("issued add subnet validator transaction %s in %s\n", addSubnetValidatorTxID, time.Since(addPermissionlessValidatorStartTime))

	addPermissionlessDelegatorStartTime := time.Now()
	addSubnetDelegatorTx, err := oWallet.IssueAddPermissionlessDelegatorTx(
		&txs.SubnetValidator{
			Validator: txs.Validator{
				NodeID: genesis.LocalConfig.InitialStakers[0].NodeID,
				Start:  uint64(startTime.Unix()),
				End:    uint64(startTime.Add(5 * time.Second).Unix()),
				Wght:   25 * units.MegaDione,
			},
			Subnet: createSubnetTxID,
		},
		createAssetTxID,
		&secp256k1fx.OutputOwners{},
	)
	if err != nil {
		log.Fatalf("failed to issue add subnet delegator with: %s\n", err)
		return
	}
	addSubnetDelegatorTxID := addSubnetDelegatorTx.ID()
	log.Printf("issued add subnet validator delegator %s in %s\n", addSubnetDelegatorTxID, time.Since(addPermissionlessDelegatorStartTime))
}
