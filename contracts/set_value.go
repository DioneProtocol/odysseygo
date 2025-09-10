package main

import (
	"context"
	"log"
	"math/big"

	BasicContract "github.com/DioneProtocol/odysseygo/contracts/basicContract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcURL := "http://localhost:9652/ext/bc/sCo6dHJvJsrGGHUfcdMVnTWBm6LFEV6pvucfNhBTvpCLuMFwF/rpc"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(43112)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = uint64(5000000)

	contract, _ := BasicContract.NewBasicContract(common.HexToAddress("0x000000000000000000000000000000000000012a"), client)
	// Set value
	setTx, err := contract.Set(auth, big.NewInt(987654))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Set tx sent: %s", setTx.Hash().Hex())

	// Read value
	val, err := contract.Get(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read value: %s", val.String())
}
