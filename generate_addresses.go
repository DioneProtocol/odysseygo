package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/DioneProtocol/odysseygo/utils/cb58"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("Generating addresses for kaito network...")
	fmt.Println("========================================")

	// Generate 3 new private keys for validators
	for i := 0; i < 3; i++ {
		// Generate random private key
		privateKeyBytes := make([]byte, 32)
		_, err := rand.Read(privateKeyBytes)
		if err != nil {
			log.Fatal("Failed to generate random bytes:", err)
		}

		// Create secp256k1 private key
		factory := &secp256k1.Factory{}
		privateKey, err := factory.ToPrivateKey(privateKeyBytes)
		if err != nil {
			log.Fatal("Failed to create private key:", err)
		}

		// Get public key
		publicKey := privateKey.PublicKey()

		// Generate address for kaito prefix
		addr, err := address.FormatBech32("kaito", publicKey.Address().Bytes())
		if err != nil {
			log.Fatal("Failed to format address:", err)
		}

		// Format with chain prefix
		fullAddr := fmt.Sprintf("A-%s", addr)

		// Derive EVM address and hex key from the same private key
		evmPriv := privateKey.ToECDSA()
		evmAddr := crypto.PubkeyToAddress(evmPriv.PublicKey).Hex()

		fmt.Printf("\nValidator %d:\n", i+1)
		privateKeyStr, _ := cb58.Encode(privateKeyBytes)
		fmt.Printf("Private Key (CB58): %s\n", privateKeyStr)
		fmt.Printf("Private Key (hex, no 0x): %064x\n", privateKeyBytes)
		publicKeyStr, _ := cb58.Encode(publicKey.Bytes())
		fmt.Printf("Public Key: %s\n", publicKeyStr)
		fmt.Printf("Address: %s\n", fullAddr)
		fmt.Printf("EVM Address: %s\n", evmAddr)
		fmt.Println("---")
	}

	fmt.Println("\nCopy these addresses to your genesis_custom.json file!")
}
