// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package primary

import (
	"context"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/keychain"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/wallet/chain/a"
	"github.com/DioneProtocol/odysseygo/wallet/chain/d"
	"github.com/DioneProtocol/odysseygo/wallet/chain/o"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var _ Wallet = (*wallet)(nil)

// Wallet provides chain wallets for the primary network.
type Wallet interface {
	O() o.Wallet
	A() a.Wallet
	D() d.Wallet
}

type wallet struct {
	o o.Wallet
	a a.Wallet
	d d.Wallet
}

func (w *wallet) O() o.Wallet {
	return w.o
}

func (w *wallet) A() a.Wallet {
	return w.a
}

func (w *wallet) D() d.Wallet {
	return w.d
}

// Creates a new default wallet
func NewWallet(o o.Wallet, a a.Wallet, d d.Wallet) Wallet {
	return &wallet{
		o: o,
		a: a,
		d: d,
	}
}

// Creates a Wallet with the given set of options
func NewWalletWithOptions(w Wallet, options ...common.Option) Wallet {
	return NewWallet(
		o.NewWalletWithOptions(w.O(), options...),
		a.NewWalletWithOptions(w.A(), options...),
		d.NewWalletWithOptions(w.D(), options...),
	)
}

type WalletConfig struct {
	// Base URI to use for all node requests.
	URI string // required
	// Keys to use for signing all transactions.
	DIONEKeychain keychain.Keychain // required
	EthKeychain   d.EthKeychain     // required
	// Set of O-chain transactions that the wallet should know about to be able
	// to generate transactions.
	OChainTxs map[ids.ID]*txs.Tx // optional
	// Set of O-chain transactions that the wallet should fetch to be able to
	// generate transactions.
	OChainTxsToFetch set.Set[ids.ID] // optional
}

// MakeWallet returns a wallet that supports issuing transactions to the chains
// living in the primary network.
//
// On creation, the wallet attaches to the provided uri and fetches all UTXOs
// that reference any of the provided keys. If the UTXOs are modified through an
// external issuance process, such as another instance of the wallet, the UTXOs
// may become out of sync. The wallet will also fetch all requested O-chain
// transactions.
//
// The wallet manages all state locally, and performs all tx signing locally.
func MakeWallet(ctx context.Context, config *WalletConfig) (Wallet, error) {
	dioneAddrs := config.DIONEKeychain.Addresses()
	dioneState, err := FetchState(ctx, config.URI, dioneAddrs)
	if err != nil {
		return nil, err
	}

	ethAddrs := config.EthKeychain.EthAddresses()
	ethState, err := FetchEthState(ctx, config.URI, ethAddrs)
	if err != nil {
		return nil, err
	}

	oChainTxs := config.OChainTxs
	if oChainTxs == nil {
		oChainTxs = make(map[ids.ID]*txs.Tx)
	}

	for txID := range config.OChainTxsToFetch {
		txBytes, err := dioneState.OClient.GetTx(ctx, txID)
		if err != nil {
			return nil, err
		}
		tx, err := txs.Parse(txs.Codec, txBytes)
		if err != nil {
			return nil, err
		}
		oChainTxs[txID] = tx
	}

	oUTXOs := NewChainUTXOs(constants.OmegaChainID, dioneState.UTXOs)
	oBackend := o.NewBackend(dioneState.OCTX, oUTXOs, oChainTxs)
	oBuilder := o.NewBuilder(dioneAddrs, oBackend)
	oSigner := o.NewSigner(config.DIONEKeychain, oBackend)

	aChainID := dioneState.ACTX.BlockchainID()
	aUTXOs := NewChainUTXOs(aChainID, dioneState.UTXOs)
	aBackend := a.NewBackend(dioneState.ACTX, aUTXOs)
	aBuilder := a.NewBuilder(dioneAddrs, aBackend)
	aSigner := a.NewSigner(config.DIONEKeychain, aBackend)

	dChainID := dioneState.DCTX.BlockchainID()
	dUTXOs := NewChainUTXOs(dChainID, dioneState.UTXOs)
	dBackend := d.NewBackend(dioneState.DCTX, dUTXOs, ethState.Accounts)
	dBuilder := d.NewBuilder(dioneAddrs, ethAddrs, dBackend)
	dSigner := d.NewSigner(config.DIONEKeychain, config.EthKeychain, dBackend)

	return NewWallet(
		o.NewWallet(oBuilder, oSigner, dioneState.OClient, oBackend),
		a.NewWallet(aBuilder, aSigner, dioneState.AClient, aBackend),
		d.NewWallet(dBuilder, dSigner, dioneState.DClient, ethState.Client, dBackend),
	), nil
}
