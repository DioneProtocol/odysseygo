package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Contract address (0x9 in the genesis)
	contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000009")

	// RPC endpoint
	rpcURL := "http://localhost:9650/ext/bc/m3Fo7ibkiC3311FGwE8qKZpv6EmA4KmwWCLb4etPEtPhRcpNN/rpc"

	// Connect to the client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Get storage value directly from storage slot 0
	storageSlot := common.Hash{} // Storage slot 0
	storageValue, err := client.StorageAt(context.Background(), contractAddress, storageSlot, nil)
	if err != nil {
		log.Fatalf("Failed to get storage value: %v", err)
	}

	// Convert bytes to big.Int
	value := new(big.Int)
	value.SetBytes(storageValue)

	fmt.Printf("Contract value: %s\n", value.String())
}
