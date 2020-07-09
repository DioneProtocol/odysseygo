// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"errors"
	"net/http"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/utils/formatting"
	"github.com/ava-labs/gecko/utils/json"
	"github.com/ava-labs/gecko/vms/components/ava"
	"github.com/ava-labs/gecko/vms/secp256k1fx"
)

// Note that since an AVA network has exactly one Platform Chain,
// and the Platform Chain defines the genesis state of the network
// (who is staking, which chains exist, etc.), defining the genesis
// state of the Platform Chain is the same as defining the genesis
// state of the network.

var (
	errUTXOHasNoValue       = errors.New("genesis UTXO has no value")
	errValidatorAddsNoValue = errors.New("validator would have already unstaked")
)

// StaticService defines the static API methods exposed by the platform VM
type StaticService struct{}

// APIUTXO is a UTXO on the Platform Chain that exists at the chain's genesis.
type APIUTXO struct {
	Amount  json.Uint64 `json:"amount"`
	Address ids.ShortID `json:"address"`
}

// FormattedAPIAccount is an APIAccount but allows for a formatted Address
type FormattedAPIAccount struct {
	Address string      `json:"address"`
	Nonce   json.Uint64 `json:"nonce"`
	Balance json.Uint64 `json:"balance"`
}

// APIValidator is a validator.
// [Amount] is the amount of $AVA being staked.
// [Endtime] is the Unix time repr. of when they are done staking
// [ID] is the node ID of the staker
// [Address] is the address where the staked AVA (and, if applicable, reward)
// is sent when this staker is done staking.
type APIValidator struct {
	StartTime   json.Uint64  `json:"startTime"`
	EndTime     json.Uint64  `json:"endTime"`
	Weight      *json.Uint64 `json:"weight,omitempty"`
	StakeAmount *json.Uint64 `json:"stakeAmount,omitempty"`
	Address     *ids.ShortID `json:"address,omitempty"`
	ID          ids.ShortID  `json:"id"`
}

func (v *APIValidator) weight() uint64 {
	switch {
	case v.Weight != nil:
		return uint64(*v.Weight)
	case v.StakeAmount != nil:
		return uint64(*v.StakeAmount)
	default:
		return 0
	}
}

// APIDefaultSubnetValidator is a validator of the default subnet
type APIDefaultSubnetValidator struct {
	APIValidator

	Destination       ids.ShortID `json:"destination"`
	DelegationFeeRate json.Uint32 `json:"delegationFeeRate"`
}

// FormattedAPIValidator allows for a formatted address
type FormattedAPIValidator struct {
	StartTime   json.Uint64  `json:"startTime"`
	EndTime     json.Uint64  `json:"endTime"`
	Weight      *json.Uint64 `json:"weight,omitempty"`
	StakeAmount *json.Uint64 `json:"stakeAmount,omitempty"`
	Address     string       `json:"address,omitempty"`
	ID          ids.ShortID  `json:"id"`
}

func (v *FormattedAPIValidator) weight() uint64 {
	switch {
	case v.Weight != nil:
		return uint64(*v.Weight)
	case v.StakeAmount != nil:
		return uint64(*v.StakeAmount)
	default:
		return 0
	}
}

// FormattedAPIDefaultSubnetValidator is a formatted validator of the default subnet
type FormattedAPIDefaultSubnetValidator struct {
	FormattedAPIValidator

	Destination       string      `json:"destination"`
	DelegationFeeRate json.Uint32 `json:"delegationFeeRate"`
}

// APIChain defines a chain that exists
// at the network's genesis.
// [GenesisData] is the initial state of the chain.
// [VMID] is the ID of the VM this chain runs.
// [FxIDs] are the IDs of the Fxs the chain supports.
// [Name] is a human-readable, non-unique name for the chain.
// [SubnetID] is the ID of the subnet that validates the chain
type APIChain struct {
	GenesisData formatting.CB58 `json:"genesisData"`
	VMID        ids.ID          `json:"vmID"`
	FxIDs       []ids.ID        `json:"fxIDs"`
	Name        string          `json:"name"`
	SubnetID    ids.ID          `json:"subnetID"`
}

// BuildGenesisArgs are the arguments used to create
// the genesis data of the Platform Chain.
// [NetworkID] is the ID of the network
// [UTXOs] are the UTXOs on the Platform Chain that exist at genesis.
// [Validators] are the validators of the default subnet at genesis.
// [Chains] are the chains that exist at genesis.
// [Time] is the Platform Chain's time at network genesis.
// TODO replace Accounts with UTXOs
type BuildGenesisArgs struct {
	AvaxAssetID ids.ID                      `json:"avaxAssetID"`
	NetworkID   json.Uint32                 `json:"address"`
	UTXOs       []APIUTXO                   `json:"utxos"`
	Validators  []APIDefaultSubnetValidator `json:"defaultSubnetValidators"`
	Chains      []APIChain                  `json:"chains"`
	Time        json.Uint64                 `json:"time"`
}

// BuildGenesisReply is the reply from BuildGenesis
type BuildGenesisReply struct {
	Bytes formatting.CB58 `json:"bytes"`
}

// Genesis represents a genesis state of the platform chain
type Genesis struct {
	UTXOs      []*ava.UTXO      `serialize:"true"`
	Validators *EventHeap       `serialize:"true"`
	Chains     []*CreateChainTx `serialize:"true"`
	Timestamp  uint64           `serialize:"true"`
}

// Initialize ...
func (g *Genesis) Initialize() error {
	for _, tx := range g.Validators.Txs {
		if err := tx.initialize(nil); err != nil {
			return err
		}
	}
	for _, chain := range g.Chains {
		if err := chain.initialize(nil); err != nil {
			return err
		}
	}
	return nil
}

// BuildGenesis build the genesis state of the Platform Chain (and thereby the AVA network.)
func (ss *StaticService) BuildGenesis(_ *http.Request, args *BuildGenesisArgs, reply *BuildGenesisReply) error {
	// Specify the accounts on the Platform chain that exist at genesis.
	utxos := make([]*ava.UTXO, 0, len(args.UTXOs))
	for i, utxo := range args.UTXOs {
		if utxo.Amount == 0 {
			return errUTXOHasNoValue
		}
		utxos = append(utxos, &ava.UTXO{
			UTXOID: ava.UTXOID{
				TxID:        ids.Empty,
				OutputIndex: uint32(i),
			},
			Asset: ava.Asset{ID: args.AvaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt:      uint64(utxo.Amount),
				Locktime: 0,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxo.Address},
				},
			},
		})
	}

	// Specify the validators that are validating the default subnet at genesis.
	validators := &EventHeap{}
	for _, validator := range args.Validators {
		weight := validator.weight()
		if weight == 0 {
			return errValidatorAddsNoValue
		}
		if uint64(validator.EndTime) <= uint64(args.Time) {
			return errValidatorAddsNoValue
		}

		tx := &addDefaultSubnetValidatorTx{
			UnsignedAddDefaultSubnetValidatorTx: UnsignedAddDefaultSubnetValidatorTx{
				BaseTx: BaseTx{
					NetworkID:    uint32(args.NetworkID),
					BlockchainID: ids.Empty,
				},
				DurationValidator: DurationValidator{
					Validator: Validator{
						NodeID: validator.ID,
						Wght:   weight,
					},
					Start: uint64(args.Time),
					End:   uint64(validator.EndTime),
				},
				Destination: validator.Destination,
			},
		}
		if err := tx.initialize(nil); err != nil {
			return err
		}

		validators.Add(tx)
	}

	// Specify the chains that exist at genesis.
	chains := []*CreateChainTx{}
	for _, chain := range args.Chains {
		// Ordinarily we sign a createChainTx. For genesis, there is no key.
		// We generate the ID of this tx by hashing the bytes of the unsigned transaction
		// TODO: Should we just sign this tx with a private key that we share publicly?
		tx := &CreateChainTx{
			UnsignedCreateChainTx: UnsignedCreateChainTx{
				BaseTx: BaseTx{
					NetworkID:    uint32(args.NetworkID),
					BlockchainID: ids.Empty,
				},
				SubnetID:    chain.SubnetID,
				ChainName:   chain.Name,
				VMID:        chain.VMID,
				FxIDs:       chain.FxIDs,
				GenesisData: chain.GenesisData.Bytes,
			},
			ControlSigs: [][crypto.SECP256K1RSigLen]byte{},
		}
		if err := tx.initialize(nil); err != nil {
			return err
		}

		chains = append(chains, tx)
	}

	// genesis holds the genesis state
	genesis := Genesis{
		UTXOs:      utxos,
		Validators: validators,
		Chains:     chains,
		Timestamp:  uint64(args.Time),
	}

	// Marshal genesis to bytes
	bytes, err := Codec.Marshal(genesis)
	reply.Bytes.Bytes = bytes
	return err
}
