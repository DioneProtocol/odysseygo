// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"context"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
)

type ChainUTXOs interface {
	AddUTXO(ctx context.Context, destinationChainID ids.ID, utxo *dione.UTXO) error
	RemoveUTXO(ctx context.Context, sourceChainID, utxoID ids.ID) error

	UTXOs(ctx context.Context, sourceChainID ids.ID) ([]*dione.UTXO, error)
	GetUTXO(ctx context.Context, sourceChainID, utxoID ids.ID) (*dione.UTXO, error)
}
