// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
)

// UTXO adds messages to UTXOs
type UTXO struct {
	dione.UTXO `serialize:"true"`
	Message    []byte `serialize:"true" json:"message"`
}

// Genesis represents a genesis state of the omega chain
type Genesis struct {
	UTXOs         []*UTXO   `serialize:"true"`
	Validators    []*txs.Tx `serialize:"true"`
	Chains        []*txs.Tx `serialize:"true"`
	Timestamp     uint64    `serialize:"true"`
	InitialSupply uint64    `serialize:"true"`
	Message       string    `serialize:"true"`
}

func Parse(genesisBytes []byte) (*Genesis, error) {
	gen := &Genesis{}
	if _, err := Codec.Unmarshal(genesisBytes, gen); err != nil {
		return nil, err
	}
	for _, tx := range gen.Validators {
		if err := tx.Initialize(txs.GenesisCodec); err != nil {
			return nil, err
		}
	}
	for _, tx := range gen.Chains {
		if err := tx.Initialize(txs.GenesisCodec); err != nil {
			return nil, err
		}
	}
	return gen, nil
}

// State represents the genesis state of the omega chain
type State struct {
	UTXOs         []*dione.UTXO
	Validators    []*txs.Tx
	Chains        []*txs.Tx
	Timestamp     uint64
	InitialSupply uint64
}

func ParseState(genesisBytes []byte) (*State, error) {
	genesis, err := Parse(genesisBytes)
	if err != nil {
		return nil, err
	}

	utxos := make([]*dione.UTXO, 0, len(genesis.UTXOs))
	for _, utxo := range genesis.UTXOs {
		utxos = append(utxos, &utxo.UTXO)
	}

	return &State{
		UTXOs:         utxos,
		Validators:    genesis.Validators,
		Chains:        genesis.Chains,
		Timestamp:     genesis.Timestamp,
		InitialSupply: genesis.InitialSupply,
	}, nil
}
