// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/vms/components/dione"
)

type UTXOGetter interface {
	GetUTXO(utxoID ids.ID) (*dione.UTXO, error)
}

type UTXOAdder interface {
	AddUTXO(utxo *dione.UTXO)
}

type UTXODeleter interface {
	DeleteUTXO(utxoID ids.ID)
}
