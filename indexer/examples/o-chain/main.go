// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DioneProtocol/odysseygo/indexer"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/blocks"
	"github.com/DioneProtocol/odysseygo/vms/proposervm/block"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

// This example program continuously polls for the next O-Chain block
// and prints the ID of the block and its transactions.
func main() {
	var (
		uri       = fmt.Sprintf("%s/ext/index/O/block", primary.LocalAPIURI)
		client    = indexer.NewClient(uri)
		ctx       = context.Background()
		nextIndex uint64
	)
	for {
		container, err := client.GetContainerByIndex(ctx, nextIndex)
		if err != nil {
			time.Sleep(time.Second)
			log.Printf("polling for next accepted block\n")
			continue
		}

		omegavmBlockBytes := container.Bytes
		proposerVMBlock, err := block.Parse(container.Bytes)
		if err == nil {
			omegavmBlockBytes = proposerVMBlock.Block()
		}

		omegavmBlock, err := blocks.Parse(blocks.Codec, omegavmBlockBytes)
		if err != nil {
			log.Fatalf("failed to parse omegavm block: %s\n", err)
		}

		acceptedTxs := omegavmBlock.Txs()
		log.Printf("accepted block %s with %d transactions\n", omegavmBlock.ID(), len(acceptedTxs))

		for _, tx := range acceptedTxs {
			log.Printf("accepted transaction %s\n", tx.ID())
		}

		nextIndex++
	}
}
