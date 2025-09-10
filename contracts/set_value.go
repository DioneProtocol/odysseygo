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
	rpcURL := "http://localhost:9652/ext/bc/2fUpUve6v2Hp54e6iCU26wQrbkCjYLf7q7Gmv8Mq9qh4TkRsL3/rpc"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("68dcbb8de4f2d15d7820f33548e857e3f957aa98d428cf487c0c5f553693f7a6")
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
