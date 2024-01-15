// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"path"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/vms/nftfx"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/genesis"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/propertyfx"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

// Aliases returns the default aliases based on the network ID
func Aliases(genesisBytes []byte) (map[string][]string, map[ids.ID][]string, error) {
	apiAliases := map[string][]string{
		path.Join(constants.ChainAliasPrefix, constants.OmegaChainID.String()): {
			"O",
			"omega",
			path.Join(constants.ChainAliasPrefix, "O"),
			path.Join(constants.ChainAliasPrefix, "omega"),
		},
	}
	chainAliases := map[ids.ID][]string{
		constants.OmegaChainID: {"O", "omega"},
	}

	genesis, err := genesis.Parse(genesisBytes) // TODO let's not re-create genesis to do aliasing
	if err != nil {
		return nil, nil, err
	}
	for _, chain := range genesis.Chains {
		uChain := chain.Unsigned.(*txs.CreateChainTx)
		chainID := chain.ID()
		endpoint := path.Join(constants.ChainAliasPrefix, chainID.String())
		switch uChain.VMID {
		case constants.AlphaID:
			apiAliases[endpoint] = []string{
				"A",
				"alpha",
				path.Join(constants.ChainAliasPrefix, "A"),
				path.Join(constants.ChainAliasPrefix, "alpha"),
			}
			chainAliases[chainID] = GetAChainAliases()
		case constants.DeltaID:
			apiAliases[endpoint] = []string{
				"D",
				"delta",
				path.Join(constants.ChainAliasPrefix, "D"),
				path.Join(constants.ChainAliasPrefix, "delta"),
			}
			chainAliases[chainID] = GetDChainAliases()
		}
	}
	return apiAliases, chainAliases, nil
}

func GetDChainAliases() []string {
	return []string{"D", "delta"}
}

func GetAChainAliases() []string {
	return []string{"A", "alpha"}
}

func GetVMAliases() map[ids.ID][]string {
	return map[ids.ID][]string{
		constants.OmegaVMID: {"omega"},
		constants.AlphaID:   {"alpha"},
		constants.DeltaID:   {"delta"},
		secp256k1fx.ID:      {"secp256k1fx"},
		nftfx.ID:            {"nftfx"},
		propertyfx.ID:       {"propertyfx"},
	}
}
