// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package alpha

import (
	"cmp"

	"github.com/DioneProtocol/odysseygo/utils"
	"github.com/DioneProtocol/odysseygo/vms/alpha/txs"
)

var _ utils.Sortable[*GenesisAsset] = (*GenesisAsset)(nil)

type Genesis struct {
	Txs []*GenesisAsset `serialize:"true"`
}

type GenesisAsset struct {
	Alias             string `serialize:"true"`
	txs.CreateAssetTx `serialize:"true"`
}

func (g *GenesisAsset) Compare(other *GenesisAsset) int {
	return cmp.Compare(g.Alias, other.Alias)
}
