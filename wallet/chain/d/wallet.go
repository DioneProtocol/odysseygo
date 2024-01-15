// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package d

import (
	"errors"
	"math/big"
	"time"

	"github.com/DioneProtocol/coreth/ethclient"
	"github.com/DioneProtocol/coreth/plugin/delta"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var (
	_ Wallet = (*wallet)(nil)

	errNotCommitted = errors.New("not committed")
)

type Wallet interface {
	Context

	// Builder returns the builder that will be used to create the transactions.
	Builder() Builder

	// Signer returns the signer that will be used to sign the transactions.
	Signer() Signer

	// IssueImportTx creates, signs, and issues an import transaction that
	// attempts to consume all the available UTXOs and import the funds to [to].
	//
	// - [chainID] specifies the chain to be importing funds from.
	// - [to] specifies where to send the imported funds to.
	IssueImportTx(
		chainID ids.ID,
		to ethcommon.Address,
		options ...common.Option,
	) (*delta.Tx, error)

	// IssueExportTx creates, signs, and issues an export transaction that
	// attempts to send all the provided [outputs] to the requested [chainID].
	//
	// - [chainID] specifies the chain to be exporting the funds to.
	// - [outputs] specifies the outputs to send to the [chainID].
	IssueExportTx(
		chainID ids.ID,
		outputs []*secp256k1fx.TransferOutput,
		options ...common.Option,
	) (*delta.Tx, error)

	// IssueUnsignedTx signs and issues the unsigned tx.
	IssueUnsignedAtomicTx(
		utx delta.UnsignedAtomicTx,
		options ...common.Option,
	) (*delta.Tx, error)

	// IssueAtomicTx issues the signed tx.
	IssueAtomicTx(
		tx *delta.Tx,
		options ...common.Option,
	) error
}

func NewWallet(
	builder Builder,
	signer Signer,
	dioneClient delta.Client,
	ethClient ethclient.Client,
	backend Backend,
) Wallet {
	return &wallet{
		Backend:     backend,
		builder:     builder,
		signer:      signer,
		dioneClient: dioneClient,
		ethClient:   ethClient,
	}
}

type wallet struct {
	Backend
	builder     Builder
	signer      Signer
	dioneClient delta.Client
	ethClient   ethclient.Client
}

func (w *wallet) Builder() Builder {
	return w.builder
}

func (w *wallet) Signer() Signer {
	return w.signer
}

func (w *wallet) IssueImportTx(
	chainID ids.ID,
	to ethcommon.Address,
	options ...common.Option,
) (*delta.Tx, error) {
	baseFee, err := w.baseFee(options)
	if err != nil {
		return nil, err
	}

	utx, err := w.builder.NewImportTx(chainID, to, baseFee, options...)
	if err != nil {
		return nil, err
	}
	return w.IssueUnsignedAtomicTx(utx, options...)
}

func (w *wallet) IssueExportTx(
	chainID ids.ID,
	outputs []*secp256k1fx.TransferOutput,
	options ...common.Option,
) (*delta.Tx, error) {
	baseFee, err := w.baseFee(options)
	if err != nil {
		return nil, err
	}

	utx, err := w.builder.NewExportTx(chainID, outputs, baseFee, options...)
	if err != nil {
		return nil, err
	}
	return w.IssueUnsignedAtomicTx(utx, options...)
}

func (w *wallet) IssueUnsignedAtomicTx(
	utx delta.UnsignedAtomicTx,
	options ...common.Option,
) (*delta.Tx, error) {
	ops := common.NewOptions(options)
	ctx := ops.Context()
	tx, err := w.signer.SignUnsignedAtomic(ctx, utx)
	if err != nil {
		return nil, err
	}

	return tx, w.IssueAtomicTx(tx, options...)
}

func (w *wallet) IssueAtomicTx(
	tx *delta.Tx,
	options ...common.Option,
) error {
	ops := common.NewOptions(options)
	ctx := ops.Context()
	txID, err := w.dioneClient.IssueTx(ctx, tx.SignedBytes())
	if err != nil {
		return err
	}

	if f := ops.PostIssuanceFunc(); f != nil {
		f(txID)
	}

	if ops.AssumeDecided() {
		return w.Backend.AcceptAtomicTx(ctx, tx)
	}

	pollFrequency := ops.PollFrequency()
	ticker := time.NewTicker(pollFrequency)
	defer ticker.Stop()

	for {
		status, err := w.dioneClient.GetAtomicTxStatus(ctx, txID)
		if err != nil {
			return err
		}

		switch status {
		case delta.Accepted:
			return w.Backend.AcceptAtomicTx(ctx, tx)
		case delta.Dropped, delta.Unknown:
			return errNotCommitted
		}

		// The tx is Processing.

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *wallet) baseFee(options []common.Option) (*big.Int, error) {
	ops := common.NewOptions(options)
	baseFee := ops.BaseFee(nil)
	if baseFee != nil {
		return baseFee, nil
	}

	ctx := ops.Context()
	return w.ethClient.EstimateBaseFee(ctx)
}
