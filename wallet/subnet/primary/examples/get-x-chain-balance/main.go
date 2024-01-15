// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/wallet/chain/a"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

func main() {
	uri := primary.LocalAPIURI
	addrStr := "A-local18jma8ppw3nhx5r4ap8clazz0dps7rv5u00z96u"

	addr, err := address.ParseToID(addrStr)
	if err != nil {
		log.Fatalf("failed to parse address: %s\n", err)
	}

	addresses := set.Set[ids.ShortID]{}
	addresses.Add(addr)

	ctx := context.Background()

	fetchStartTime := time.Now()
	state, err := primary.FetchState(ctx, uri, addresses)
	if err != nil {
		log.Fatalf("failed to fetch state: %s\n", err)
	}
	log.Printf("fetched state of %s in %s\n", addrStr, time.Since(fetchStartTime))

	aChainID := state.ACTX.BlockchainID()

	aUTXOs := primary.NewChainUTXOs(aChainID, state.UTXOs)
	aBackend := a.NewBackend(state.ACTX, aUTXOs)
	aBuilder := a.NewBuilder(addresses, aBackend)

	currentBalances, err := aBuilder.GetFTBalance()
	if err != nil {
		log.Fatalf("failed to get the balance: %s\n", err)
	}

	dioneID := state.ACTX.DIONEAssetID()
	dioneBalance := currentBalances[dioneID]
	log.Printf("current DIONE balance of %s is %d nDIONE\n", addrStr, dioneBalance)
}
