// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/dioneprotocol/dionego/genesis"
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/utils/formatting/address"
	"github.com/dioneprotocol/dionego/utils/units"
	"github.com/dioneprotocol/dionego/vms/components/dione"
	"github.com/dioneprotocol/dionego/vms/platformvm/stakeable"
	"github.com/dioneprotocol/dionego/vms/secp256k1fx"
	"github.com/dioneprotocol/dionego/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	amount := 500 * units.MilliDione
	locktime := uint64(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	destAddrStr := "P-local18jma8ppw3nhx5r4ap8clazz0dps7rv5u00z96u"

	destAddr, err := address.ParseToID(destAddrStr)
	if err != nil {
		log.Fatalf("failed to parse address: %s\n", err)
	}

	ctx := context.Background()

	// NewWalletFromURI fetches the available UTXOs owned by [kc] on the network
	// that [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.NewWalletFromURI(ctx, uri, kc)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()
	dioneAssetID := pWallet.DIONEAssetID()

	issueTxStartTime := time.Now()
	txID, err := pWallet.IssueBaseTx([]*dione.TransferableOutput{
		{
			Asset: dione.Asset{
				ID: dioneAssetID,
			},
			Out: &stakeable.LockOut{
				Locktime: locktime,
				TransferableOut: &secp256k1fx.TransferOutput{
					Amt: amount,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs: []ids.ShortID{
							destAddr,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to issue transaction: %s\n", err)
	}
	log.Printf("issued %s in %s\n", txID, time.Since(issueTxStartTime))
}
