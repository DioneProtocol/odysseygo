// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package testnet

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/DioneProtocol/coreth/core"
	"github.com/DioneProtocol/coreth/params"
	"github.com/DioneProtocol/coreth/plugin/delta"

	"github.com/DioneProtocol/odysseygo/config"
	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/staking"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/utils/perms"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/reward"
)

const (
	DefaultNodeCount      = 2 // Minimum required to ensure connectivity-based health checks will pass
	DefaultFundedKeyCount = 50

	DefaultGasLimit = uint64(100_000_000) // Gas limit is arbitrary

	// Arbitrarily large amount of DIONE to fund keys on the A-Chain for testing
	DefaultFundedKeyAChainAmount = 30 * units.MegaDione
)

var (
	// Arbitrarily large amount of DIONE (10^12) to fund keys on the D-Chain for testing
	DefaultFundedKeyDChainAmount = new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)

	errEmptyValidatorsForGenesis   = errors.New("failed to generate genesis: empty validator IDs")
	errNoKeysForGenesis            = errors.New("failed to generate genesis: no keys to fund")
	errInvalidNetworkIDForGenesis  = errors.New("network ID can't be mainnet, testnet or local network ID")
	errMissingValidatorsForGenesis = errors.New("no genesis validators provided")
	errMissingBalancesForGenesis   = errors.New("no genesis balances given")
	errMissingTLSKeyForNodeID      = fmt.Errorf("failed to ensure node ID: missing value for %q", config.StakingTLSKeyContentKey)
	errMissingCertForNodeID        = fmt.Errorf("failed to ensure node ID: missing value for %q", config.StakingCertContentKey)
	errInvalidKeypair              = fmt.Errorf("%q and %q must be provided together or not at all", config.StakingTLSKeyContentKey, config.StakingCertContentKey)
)

// Defines a mapping of flag keys to values intended to be supplied to
// an invocation of an OdysseyGo node.
type FlagsMap map[string]interface{}

// SetDefaults ensures the effectiveness of flag overrides by only
// setting values supplied in the defaults map that are not already
// explicitly set.
func (f FlagsMap) SetDefaults(defaults FlagsMap) {
	for key, value := range defaults {
		if _, ok := f[key]; !ok {
			f[key] = value
		}
	}
}

// GetStringVal simplifies retrieving a map value as a string.
func (f FlagsMap) GetStringVal(key string) (string, error) {
	rawVal, ok := f[key]
	if !ok {
		return "", nil
	}

	val, err := cast.ToStringE(rawVal)
	if err != nil {
		return "", fmt.Errorf("failed to cast value for %q: %w", key, err)
	}
	return val, nil
}

// Write simplifies writing a FlagsMap to the provided path. The
// description is used in error messages.
func (f FlagsMap) Write(path string, description string) error {
	bytes, err := DefaultJSONMarshal(f)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", description, err)
	}
	if err := os.WriteFile(path, bytes, perms.ReadWrite); err != nil {
		return fmt.Errorf("failed to write %s: %w", description, err)
	}
	return nil
}

// Utility function simplifying construction of a FlagsMap from a file.
func ReadFlagsMap(path string, description string) (*FlagsMap, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", description, err)
	}
	flagsMap := &FlagsMap{}
	if err := json.Unmarshal(bytes, flagsMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", description, err)
	}
	return flagsMap, nil
}

// Marshal to json with default prefix and indent.
func DefaultJSONMarshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NetworkConfig defines configuration shared or
// common to all nodes in a given network.
type NetworkConfig struct {
	Genesis      *genesis.UnparsedConfig
	DChainConfig FlagsMap
	DefaultFlags FlagsMap
	FundedKeys   []*secp256k1.PrivateKey
}

// Ensure genesis is generated if not already present.
func (c *NetworkConfig) EnsureGenesis(networkID uint32, validatorIDs []ids.NodeID) error {
	if c.Genesis != nil {
		return nil
	}

	if len(validatorIDs) == 0 {
		return errEmptyValidatorsForGenesis
	}
	if len(c.FundedKeys) == 0 {
		return errNoKeysForGenesis
	}

	// Ensure pre-funded keys have arbitrary large balances on both chains to support testing
	aChainBalances := make(AChainBalanceMap, len(c.FundedKeys))
	dChainBalances := make(core.GenesisAlloc, len(c.FundedKeys))
	for _, key := range c.FundedKeys {
		aChainBalances[key.Address()] = DefaultFundedKeyAChainAmount
		dChainBalances[delta.GetEthAddress(key)] = core.GenesisAccount{
			Balance: DefaultFundedKeyDChainAmount,
		}
	}

	genesis, err := NewTestGenesis(networkID, aChainBalances, dChainBalances, validatorIDs)
	if err != nil {
		return err
	}

	c.Genesis = genesis
	return nil
}

// NodeURI associates a node ID with its API URI.
type NodeURI struct {
	NodeID ids.NodeID
	URI    string
}

// NodeConfig defines configuration for an OdysseyGo node.
type NodeConfig struct {
	NodeID ids.NodeID
	Flags  FlagsMap
}

func NewNodeConfig() *NodeConfig {
	return &NodeConfig{
		Flags: FlagsMap{},
	}
}

// Convenience method for setting networking flags.
func (nc *NodeConfig) SetNetworkingConfigDefaults(
	httpPort uint16,
	stakingPort uint16,
	bootstrapIDs []string,
	bootstrapIPs []string,
) {
	nc.Flags.SetDefaults(FlagsMap{
		config.HTTPPortKey:     httpPort,
		config.StakingPortKey:  stakingPort,
		config.BootstrapIDsKey: strings.Join(bootstrapIDs, ","),
		config.BootstrapIPsKey: strings.Join(bootstrapIPs, ","),
	})
}

// Ensures staking and signing keys are generated if not already present and
// that the node ID (derived from the staking keypair) is set.
func (nc *NodeConfig) EnsureKeys() error {
	if err := nc.EnsureBLSSigningKey(); err != nil {
		return err
	}
	if err := nc.EnsureStakingKeypair(); err != nil {
		return err
	}
	// Once a staking keypair is guaranteed it is safe to derive the node ID
	return nc.EnsureNodeID()
}

// Ensures a BLS signing key is generated if not already present.
func (nc *NodeConfig) EnsureBLSSigningKey() error {
	// Attempt to retrieve an existing key
	existingKey, err := nc.Flags.GetStringVal(config.StakingSignerKeyContentKey)
	if err != nil {
		return err
	}
	if len(existingKey) > 0 {
		// Nothing to do
		return nil
	}

	// Generate a new signing key
	newKey, err := bls.NewSecretKey()
	if err != nil {
		return fmt.Errorf("failed to generate staking signer key: %w", err)
	}
	nc.Flags[config.StakingSignerKeyContentKey] = base64.StdEncoding.EncodeToString(newKey.Serialize())
	return nil
}

// Ensures a staking keypair is generated if not already present.
func (nc *NodeConfig) EnsureStakingKeypair() error {
	keyKey := config.StakingTLSKeyContentKey
	certKey := config.StakingCertContentKey

	key, err := nc.Flags.GetStringVal(keyKey)
	if err != nil {
		return err
	}

	cert, err := nc.Flags.GetStringVal(certKey)
	if err != nil {
		return err
	}

	if len(key) == 0 && len(cert) == 0 {
		// Generate new keypair
		tlsCertBytes, tlsKeyBytes, err := staking.NewCertAndKeyBytes()
		if err != nil {
			return fmt.Errorf("failed to generate staking keypair: %w", err)
		}
		nc.Flags[keyKey] = base64.StdEncoding.EncodeToString(tlsKeyBytes)
		nc.Flags[certKey] = base64.StdEncoding.EncodeToString(tlsCertBytes)
	} else if len(key) == 0 || len(cert) == 0 {
		// Only one of key and cert was provided
		return errInvalidKeypair
	}

	err = nc.EnsureNodeID()
	if err != nil {
		return fmt.Errorf("failed to derive a node ID: %w", err)
	}

	return nil
}

// Attempt to derive the node ID from the node configuration.
func (nc *NodeConfig) EnsureNodeID() error {
	keyKey := config.StakingTLSKeyContentKey
	certKey := config.StakingCertContentKey

	key, err := nc.Flags.GetStringVal(keyKey)
	if err != nil {
		return err
	}
	if len(key) == 0 {
		return errMissingTLSKeyForNodeID
	}
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return fmt.Errorf("failed to ensure node ID: failed to base64 decode value for %q: %w", keyKey, err)
	}

	cert, err := nc.Flags.GetStringVal(certKey)
	if err != nil {
		return err
	}
	if len(cert) == 0 {
		return errMissingCertForNodeID
	}
	certBytes, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return fmt.Errorf("failed to ensure node ID: failed to base64 decode value for %q: %w", certKey, err)
	}

	tlsCert, err := staking.LoadTLSCertFromBytes(keyBytes, certBytes)
	if err != nil {
		return fmt.Errorf("failed to ensure node ID: failed to load tls cert: %w", err)
	}
	stakingCert := staking.CertificateFromX509(tlsCert.Leaf)
	nc.NodeID = ids.NodeIDFromCert(stakingCert)

	return nil
}

// Helper type to simplify configuring A-Chain genesis balances
type AChainBalanceMap map[ids.ShortID]uint64

// Create a genesis struct valid for bootstrapping a test
// network. Note that many of the genesis fields (e.g. reward
// addresses) are randomly generated or hard-coded.
func NewTestGenesis(
	networkID uint32,
	aChainBalances AChainBalanceMap,
	dChainBalances core.GenesisAlloc,
	validatorIDs []ids.NodeID,
) (*genesis.UnparsedConfig, error) {
	// Validate inputs
	switch networkID {
	case constants.TestnetID, constants.MainnetID, constants.LocalID:
		return nil, errInvalidNetworkIDForGenesis
	}
	if len(validatorIDs) == 0 {
		return nil, errMissingValidatorsForGenesis
	}
	if len(aChainBalances) == 0 || len(dChainBalances) == 0 {
		return nil, errMissingBalancesForGenesis
	}

	// Address that controls stake doesn't matter -- generate it randomly
	stakeAddress, err := address.Format(
		"A",
		constants.GetHRP(networkID),
		ids.GenerateTestShortID().Bytes(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to format stake address: %w", err)
	}

	// Ensure the total stake allows a MegaDione per validator
	totalStake := uint64(len(validatorIDs)) * units.MegaDione

	// The eth address is only needed to link pre-mainnet assets. Until that capability
	// becomes necessary for testing, use a bogus address.
	//
	// Reference: https://github.com/DioneProtocol/odysseygo/issues/1365#issuecomment-1511508767
	ethAddress := "0x0000000000000000000000000000000000000000"

	now := time.Now()

	config := &genesis.UnparsedConfig{
		NetworkID: networkID,
		Allocations: []genesis.UnparsedAllocation{
			{
				ETHAddr:       ethAddress,
				DIONEAddr:     stakeAddress,
				InitialAmount: 0,
				UnlockSchedule: []genesis.LockedAmount{ // Provides stake to validators
					{
						Amount:   totalStake,
						Locktime: uint64(now.Add(7 * 24 * time.Hour).Unix()), // 1 Week
					},
				},
			},
		},
		StartTime:                  uint64(now.Unix()),
		InitialStakedFunds:         []string{stakeAddress},
		InitialStakeDuration:       365 * 24 * 60 * 60, // 1 year
		InitialStakeDurationOffset: 90 * 60,            // 90 minutes
		Message:                    "hello odyssey!",
	}

	// Set A-Chain balances
	for aChainAddress, balance := range aChainBalances {
		dioneAddr, err := address.Format("A", constants.GetHRP(networkID), aChainAddress[:])
		if err != nil {
			return nil, fmt.Errorf("failed to format A-Chain address: %w", err)
		}
		config.Allocations = append(
			config.Allocations,
			genesis.UnparsedAllocation{
				ETHAddr:       ethAddress,
				DIONEAddr:     dioneAddr,
				InitialAmount: balance,
				UnlockSchedule: []genesis.LockedAmount{
					{
						Amount: 20 * units.MegaDione,
					},
					{
						Amount:   totalStake,
						Locktime: uint64(now.Add(7 * 24 * time.Hour).Unix()), // 1 Week
					},
				},
			},
		)
	}

	// Define D-Chain genesis
	dChainGenesis := &core.Genesis{
		Config: &params.ChainConfig{
			ChainID: big.NewInt(131312), // Arbitrary chain ID is arbitrary
		},
		Difficulty: big.NewInt(0), // Difficulty is a mandatory field
		GasLimit:   DefaultGasLimit,
		Alloc:      dChainBalances,
	}
	dChainGenesisBytes, err := json.Marshal(dChainGenesis)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal D-Chain genesis: %w", err)
	}
	config.DChainGenesis = string(dChainGenesisBytes)

	// Give staking rewards for initial validators to a random address. Any testing of staking rewards
	// will be easier to perform with nodes other than the initial validators since the timing of
	// staking can be more easily controlled.
	rewardAddr, err := address.Format("A", constants.GetHRP(networkID), ids.GenerateTestShortID().Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format reward address: %w", err)
	}

	// Configure provided validator node IDs as initial stakers
	for _, validatorID := range validatorIDs {
		config.InitialStakers = append(
			config.InitialStakers,
			genesis.UnparsedStaker{
				NodeID:        validatorID,
				RewardAddress: rewardAddr,
				DelegationFee: .01 * reward.PercentDenominator,
			},
		)
	}

	return config, nil
}
