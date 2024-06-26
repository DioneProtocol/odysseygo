// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package primary

import (
	"context"
	"fmt"

	"github.com/DioneProtocol/coreth/ethclient"
	"github.com/DioneProtocol/coreth/plugin/delta"

	"github.com/ethereum/go-ethereum/common"

	"github.com/DioneProtocol/odysseygo/api/info"
	"github.com/DioneProtocol/odysseygo/codec"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/rpc"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/alpha"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/wallet/chain/a"
	"github.com/DioneProtocol/odysseygo/wallet/chain/d"
	"github.com/DioneProtocol/odysseygo/wallet/chain/o"
)

const (
	MainnetAPIURI = "https://api.dione.network"
	TestnetAPIURI = "https://api.dione-test.network"
	LocalAPIURI   = "http://localhost:9650"

	fetchLimit = 1024
)

// TODO: Refactor UTXOClient definition to allow the client implementations to
// perform their own assertions.
var (
	_ UTXOClient = omegavm.Client(nil)
	_ UTXOClient = alpha.Client(nil)
)

type UTXOClient interface {
	GetAtomicUTXOs(
		ctx context.Context,
		addrs []ids.ShortID,
		sourceChain string,
		limit uint32,
		startAddress ids.ShortID,
		startUTXOID ids.ID,
		options ...rpc.Option,
	) ([][]byte, ids.ShortID, ids.ID, error)
}

type DIONEState struct {
	OClient omegavm.Client
	OCTX    o.Context
	AClient alpha.Client
	ACTX    a.Context
	DClient delta.Client
	DCTX    d.Context
	UTXOs   UTXOs
}

func FetchState(
	ctx context.Context,
	uri string,
	addrs set.Set[ids.ShortID],
) (
	*DIONEState,
	error,
) {
	infoClient := info.NewClient(uri)
	oClient := omegavm.NewClient(uri)
	aClient := alpha.NewClient(uri, "A")
	dClient := delta.NewDChainClient(uri)

	oCTX, err := o.NewContextFromClients(ctx, infoClient, aClient)
	if err != nil {
		return nil, err
	}

	aCTX, err := a.NewContextFromClients(ctx, infoClient, aClient)
	if err != nil {
		return nil, err
	}

	dCTX, err := d.NewContextFromClients(ctx, infoClient, aClient)
	if err != nil {
		return nil, err
	}

	utxos := NewUTXOs()
	addrList := addrs.List()
	chains := []struct {
		id     ids.ID
		client UTXOClient
		codec  codec.Manager
	}{
		{
			id:     constants.OmegaChainID,
			client: oClient,
			codec:  txs.Codec,
		},
		{
			id:     aCTX.BlockchainID(),
			client: aClient,
			codec:  a.Parser.Codec(),
		},
		{
			id:     dCTX.BlockchainID(),
			client: dClient,
			codec:  delta.Codec,
		},
	}
	for _, destinationChain := range chains {
		for _, sourceChain := range chains {
			err = AddAllUTXOs(
				ctx,
				utxos,
				destinationChain.client,
				destinationChain.codec,
				sourceChain.id,
				destinationChain.id,
				addrList,
			)
			if err != nil {
				return nil, err
			}
		}
	}
	return &DIONEState{
		OClient: oClient,
		OCTX:    oCTX,
		AClient: aClient,
		ACTX:    aCTX,
		DClient: dClient,
		DCTX:    dCTX,
		UTXOs:   utxos,
	}, nil
}

type EthState struct {
	Client   ethclient.Client
	Accounts map[common.Address]*d.Account
}

func FetchEthState(
	ctx context.Context,
	uri string,
	addrs set.Set[common.Address],
) (*EthState, error) {
	path := fmt.Sprintf(
		"%s/ext/%s/D/rpc",
		uri,
		constants.ChainAliasPrefix,
	)
	client, err := ethclient.Dial(path)
	if err != nil {
		return nil, err
	}

	accounts := make(map[common.Address]*d.Account, addrs.Len())
	for addr := range addrs {
		balance, err := client.BalanceAt(ctx, addr, nil)
		if err != nil {
			return nil, err
		}
		nonce, err := client.NonceAt(ctx, addr, nil)
		if err != nil {
			return nil, err
		}
		accounts[addr] = &d.Account{
			Balance: balance,
			Nonce:   nonce,
		}
	}
	return &EthState{
		Client:   client,
		Accounts: accounts,
	}, nil
}

// AddAllUTXOs fetches all the UTXOs referenced by [addresses] that were sent
// from [sourceChainID] to [destinationChainID] from the [client]. It then uses
// [codec] to parse the returned UTXOs and it adds them into [utxos]. If [ctx]
// expires, then the returned error will be immediately reported.
func AddAllUTXOs(
	ctx context.Context,
	utxos UTXOs,
	client UTXOClient,
	codec codec.Manager,
	sourceChainID ids.ID,
	destinationChainID ids.ID,
	addrs []ids.ShortID,
) error {
	var (
		sourceChainIDStr = sourceChainID.String()
		startAddr        ids.ShortID
		startUTXO        ids.ID
	)
	for {
		utxosBytes, endAddr, endUTXO, err := client.GetAtomicUTXOs(
			ctx,
			addrs,
			sourceChainIDStr,
			fetchLimit,
			startAddr,
			startUTXO,
		)
		if err != nil {
			return err
		}

		for _, utxoBytes := range utxosBytes {
			var utxo dione.UTXO
			_, err := codec.Unmarshal(utxoBytes, &utxo)
			if err != nil {
				return err
			}

			if err := utxos.AddUTXO(ctx, sourceChainID, destinationChainID, &utxo); err != nil {
				return err
			}
		}

		if len(utxosBytes) < fetchLimit {
			break
		}

		// Update the vars to query the next page of UTXOs.
		startAddr = endAddr
		startUTXO = endUTXO
	}
	return nil
}
