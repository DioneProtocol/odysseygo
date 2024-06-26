// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/DioneProtocol/coreth/plugin/delta"

	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	dioneAddr := key.Address()
	ethAddr := delta.PublicKeyToEthAddress(key.PublicKey())

	ctx := context.Background()

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:           uri,
		DIONEKeychain: kc,
		EthKeychain:   kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the O-chain wallet
	oWallet := wallet.O()
	dWallet := wallet.D()

	// Pull out useful constants to use when issuing transactions.
	dChainID := dWallet.BlockchainID()
	dioneAssetID := dWallet.DIONEAssetID()
	owner := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			dioneAddr,
		},
	}

	exportStartTime := time.Now()
	exportTx, err := oWallet.IssueExportTx(dChainID, []*dione.TransferableOutput{{
		Asset: dione.Asset{ID: dioneAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt:          units.Dione,
			OutputOwners: owner,
		},
	}})
	if err != nil {
		log.Fatalf("failed to issue export transaction: %s\n", err)
	}
	log.Printf("issued export %s in %s\n", exportTx.ID(), time.Since(exportStartTime))

	importStartTime := time.Now()
	importTx, err := dWallet.IssueImportTx(constants.OmegaChainID, ethAddr)
	if err != nil {
		log.Fatalf("failed to issue import transaction: %s\n", err)
	}
	log.Printf("issued import %s to %s in %s\n", importTx.ID(), ethAddr.Hex(), time.Since(importStartTime))
}
