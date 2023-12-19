// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/DioneProtocol/odysseygo/api/info"
	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	startTime := time.Now().Add(time.Minute)
	duration := 3 * 7 * 24 * time.Hour // 3 weeks
	weight := 2_000 * units.Dione
	validatorRewardAddr := key.Address()

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	nodeInfoStartTime := time.Now()
	nodeID, nodePOP, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

	// NewWalletFromURI fetches the available UTXOs owned by [kc] on the network
	// that [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.NewWalletFromURI(ctx, uri, kc)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the O-chain wallet
	pWallet := wallet.O()
	dioneAssetID := pWallet.DIONEAssetID()

	addValidatorStartTime := time.Now()
	addValidatorTxID, err := pWallet.IssueAddPermissionlessValidatorTx(
		&txs.SubnetValidator{Validator: txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(startTime.Add(duration).Unix()),
			Wght:   weight,
		}},
		nodePOP,
		dioneAssetID,
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{validatorRewardAddr},
		},
	)
	if err != nil {
		log.Fatalf("failed to issue add permissionless validator transaction: %s\n", err)
	}
	log.Printf("added new primary network validator %s with %s in %s\n", nodeID, addValidatorTxID, time.Since(addValidatorStartTime))
}
