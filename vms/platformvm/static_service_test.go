// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"testing"

	"github.com/ava-labs/gecko/utils/constants"

	"github.com/ava-labs/gecko/utils/formatting"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/json"
)

func TestBuildGenesisInvalidUTXOBalance(t *testing.T) {
	id := ids.NewShortID([20]byte{1, 2, 3})
	nodeID := id.PrefixedString(constants.NodeIDPrefix)
	hrp := constants.NetworkIDToHRP[testNetworkID]
	addr, err := formatting.FormatBech32(hrp, id.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	utxo := APIUTXO{
		Address: addr,
		Amount:  0,
	}
	weight := json.Uint64(987654321)
	validator := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			EndTime: 15,
			Weight:  &weight,
			ID:      nodeID,
		},
		Destination: addr,
	}

	args := BuildGenesisArgs{
		UTXOs: []APIUTXO{
			utxo,
		},
		Validators: []FormattedAPIDefaultSubnetValidator{
			validator,
		},
		Time: 5,
	}
	reply := BuildGenesisReply{}

	ss := StaticService{}
	if err := ss.BuildGenesis(nil, &args, &reply); err == nil {
		t.Fatalf("Should have errored due to an invalid balance")
	}
}

func TestBuildGenesisInvalidAmount(t *testing.T) {
	id := ids.NewShortID([20]byte{1, 2, 3})
	nodeID := id.PrefixedString(constants.NodeIDPrefix)
	hrp := constants.NetworkIDToHRP[testNetworkID]
	addr, err := formatting.FormatBech32(hrp, id.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	utxo := APIUTXO{
		Address: addr,
		Amount:  123456789,
	}
	weight := json.Uint64(0)
	validator := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			StartTime: 0,
			EndTime:   15,
			Weight:    &weight,
			ID:        nodeID,
		},
		Destination: addr,
	}

	args := BuildGenesisArgs{
		UTXOs: []APIUTXO{
			utxo,
		},
		Validators: []FormattedAPIDefaultSubnetValidator{
			validator,
		},
		Time: 5,
	}
	reply := BuildGenesisReply{}

	ss := StaticService{}
	if err := ss.BuildGenesis(nil, &args, &reply); err == nil {
		t.Fatalf("Should have errored due to an invalid amount")
	}
}

func TestBuildGenesisInvalidEndtime(t *testing.T) {
	id := ids.NewShortID([20]byte{1, 2, 3})
	nodeID := id.PrefixedString(constants.NodeIDPrefix)
	hrp := constants.NetworkIDToHRP[testNetworkID]
	addr, err := formatting.FormatBech32(hrp, id.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	utxo := APIUTXO{
		Address: addr,
		Amount:  123456789,
	}

	weight := json.Uint64(987654321)
	validator := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			StartTime: 0,
			EndTime:   5,
			Weight:    &weight,
			ID:        nodeID,
		},
		Destination: addr,
	}

	args := BuildGenesisArgs{
		UTXOs: []APIUTXO{
			utxo,
		},
		Validators: []FormattedAPIDefaultSubnetValidator{
			validator,
		},
		Time: 5,
	}
	reply := BuildGenesisReply{}

	ss := StaticService{}
	if err := ss.BuildGenesis(nil, &args, &reply); err == nil {
		t.Fatalf("Should have errored due to an invalid end time")
	}
}

func TestBuildGenesisReturnsSortedValidators(t *testing.T) {
	id := ids.NewShortID([20]byte{1})
	nodeID := id.PrefixedString(constants.NodeIDPrefix)
	hrp := constants.NetworkIDToHRP[testNetworkID]
	addr, err := formatting.FormatBech32(hrp, id.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	utxo := APIUTXO{
		Address: addr,
		Amount:  123456789,
	}

	weight := json.Uint64(987654321)
	validator1 := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			StartTime: 0,
			EndTime:   20,
			Weight:    &weight,
			ID:        nodeID,
		},
		Destination: addr,
	}

	validator2 := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			StartTime: 3,
			EndTime:   15,
			Weight:    &weight,
			ID:        nodeID,
		},
		Destination: addr,
	}

	validator3 := FormattedAPIDefaultSubnetValidator{
		FormattedAPIValidator: FormattedAPIValidator{
			StartTime: 1,
			EndTime:   10,
			Weight:    &weight,
			ID:        nodeID,
		},
		Destination: addr,
	}

	args := BuildGenesisArgs{
		AvaxAssetID: ids.NewID([32]byte{'d', 'u', 'm', 'm', 'y', ' ', 'I', 'D'}),
		UTXOs: []APIUTXO{
			utxo,
		},
		Validators: []FormattedAPIDefaultSubnetValidator{
			validator1,
			validator2,
			validator3,
		},
		Time: 5,
	}
	reply := BuildGenesisReply{}

	ss := StaticService{}
	if err := ss.BuildGenesis(nil, &args, &reply); err != nil {
		t.Fatalf("BuildGenesis should not have errored but got error: %s", err)
	}

	genesis := &Genesis{}
	if err := Codec.Unmarshal(reply.Bytes.Bytes, genesis); err != nil {
		t.Fatal(err)
	}
	validators := genesis.Validators
	if validators.Len() == 0 {
		t.Fatal("Validators should contain 3 validators")
	}
	currentValidator := validators.Remove().UnsignedProposalTx.(*UnsignedAddDefaultSubnetValidatorTx)
	for validators.Len() > 0 {
		nextValidator := validators.Remove().UnsignedProposalTx.(*UnsignedAddDefaultSubnetValidatorTx)
		if currentValidator.EndTime().Unix() > nextValidator.EndTime().Unix() {
			t.Fatalf("Validators returned by genesis should be a min heap sorted by end time")
		}
		currentValidator = nextValidator
	}
}
