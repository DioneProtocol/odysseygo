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
	subnetIDStr := "29uVeLPJB1eQJkzRemU8g8wZDw5uJRqpab5U2mX9euieVwiEbL"
	startTime := time.Now().Add(time.Minute)
	duration := 2 * 7 * 24 * time.Hour // 2 weeks
	weight := units.Schmeckle

	subnetID, err := ids.FromString(subnetIDStr)
	if err != nil {
		log.Fatalf("failed to parse subnet ID: %s\n", err)
	}

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	nodeInfoStartTime := time.Now()
	nodeID, _, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

	// NewWalletWithTxs fetches the available UTXOs owned by [kc] on the network
	// that [uri] is hosting and registers [subnetID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.NewWalletWithTxs(ctx, uri, kc, subnetID)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the O-chain wallet
	pWallet := wallet.O()

	addValidatorStartTime := time.Now()
	addValidatorTxID, err := pWallet.IssueAddSubnetValidatorTx(&txs.SubnetValidator{
		Validator: txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(startTime.Add(duration).Unix()),
			Wght:   weight,
		},
		Subnet: subnetID,
	})
	if err != nil {
		log.Fatalf("failed to issue add subnet validator transaction: %s\n", err)
	}
	log.Printf("added new subnet validator %s to %s with %s in %s\n", nodeID, subnetID, addValidatorTxID, time.Since(addValidatorStartTime))
}
