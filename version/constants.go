// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package version

import (
	"encoding/json"
	"time"

	_ "embed"

	"github.com/DioneProtocol/odysseygo/utils/constants"
)

// RPCChainVMProtocol should be bumped anytime changes are made which require
// the plugin vm to upgrade to latest odysseygo release to be compatible.
const RPCChainVMProtocol uint = 26

// These are globals that describe network upgrades and node versions
var (
	Current = &Semantic{
		Major: 1,
		Minor: 10,
		Patch: 2,
	}
	CurrentApp = &Application{
		Major: Current.Major,
		Minor: Current.Minor,
		Patch: Current.Patch,
	}
	MinimumCompatibleVersion = &Application{
		Major: 1,
		Minor: 10,
		Patch: 0,
	}
	PrevMinimumCompatibleVersion = &Application{
		Major: 1,
		Minor: 9,
		Patch: 0,
	}

	CurrentDatabase = DatabaseVersion1_4_5
	PrevDatabase    = DatabaseVersion1_0_0

	DatabaseVersion1_4_5 = &Semantic{
		Major: 1,
		Minor: 4,
		Patch: 5,
	}
	DatabaseVersion1_0_0 = &Semantic{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}

	//go:embed compatibility.json
	rpcChainVMProtocolCompatibilityBytes []byte
	// RPCChainVMProtocolCompatibility maps RPCChainVMProtocol versions to the
	// set of odysseygo versions that supported that version. This is not used
	// by odysseygo, but is useful for downstream libraries.
	RPCChainVMProtocolCompatibility map[uint][]*Semantic

	OdysseyPhase1Times = map[uint32]time.Time{
		constants.MainnetID: time.Date(2023, time.October, 26, 11, 46, 0, 0, time.UTC),
		constants.TestnetID: time.Date(2021, time.September, 16, 21, 0, 0, 0, time.UTC),
	}
	OdysseyPhase1DefaultTime     = time.Date(2020, time.December, 5, 5, 0, 0, 0, time.UTC)
	OdysseyPhase1MinPChainHeight = map[uint32]uint64{
		constants.MainnetID: 0,
		constants.TestnetID: 0,
	}
	OdysseyPhase1DefaultMinPChainHeight uint64

	BanffTimes = map[uint32]time.Time{
		constants.MainnetID: time.Date(2023, time.October, 26, 11, 46, 0, 0, time.UTC),
		constants.TestnetID: time.Date(2022, time.October, 3, 14, 0, 0, 0, time.UTC),
	}
	BanffDefaultTime = time.Date(2020, time.December, 5, 5, 0, 0, 0, time.UTC)

	CortinaTimes = map[uint32]time.Time{
		// constants.MainnetID: time.Date(2023, time.October, 26, 11, 54, 0, 0, time.UTC),
		// constants.TestnetID: time.Date(2023, time.April, 6, 15, 0, 0, 0, time.UTC),
	}
	CortinaDefaultTime = time.Date(2020, time.December, 5, 5, 0, 0, 0, time.UTC)
)

func init() {
	var parsedRPCChainVMCompatibility map[uint][]string
	err := json.Unmarshal(rpcChainVMProtocolCompatibilityBytes, &parsedRPCChainVMCompatibility)
	if err != nil {
		panic(err)
	}

	RPCChainVMProtocolCompatibility = make(map[uint][]*Semantic)
	for rpcChainVMProtocol, versionStrings := range parsedRPCChainVMCompatibility {
		versions := make([]*Semantic, len(versionStrings))
		for i, versionString := range versionStrings {
			version, err := Parse(versionString)
			if err != nil {
				panic(err)
			}
			versions[i] = version
		}
		RPCChainVMProtocolCompatibility[rpcChainVMProtocol] = versions
	}
}

func GetOdysseyPhase1Time(networkID uint32) time.Time {
	if upgradeTime, exists := OdysseyPhase1Times[networkID]; exists {
		return upgradeTime
	}
	return OdysseyPhase1DefaultTime
}

func GetOdysseyPhase1MinPChainHeight(networkID uint32) uint64 {
	if minHeight, exists := OdysseyPhase1MinPChainHeight[networkID]; exists {
		return minHeight
	}
	return OdysseyPhase1DefaultMinPChainHeight
}

func GetBanffTime(networkID uint32) time.Time {
	if upgradeTime, exists := BanffTimes[networkID]; exists {
		return upgradeTime
	}
	return BanffDefaultTime
}

func GetCortinaTime(networkID uint32) time.Time {
	if upgradeTime, exists := CortinaTimes[networkID]; exists {
		return upgradeTime
	}
	return CortinaDefaultTime
}

func GetCompatibility(networkID uint32) Compatibility {
	return NewCompatibility(
		CurrentApp,
		MinimumCompatibleVersion,
		GetCortinaTime(networkID),
		PrevMinimumCompatibleVersion,
	)
}
