// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"

	"github.com/ava-labs/avalanchego/nat"
	"github.com/ava-labs/avalanchego/node"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/dynamicip"
	"github.com/ava-labs/avalanchego/utils/logging"
)

const (
	header = "" +
		`     _____               .__                       .__` + "\n" +
		`    /  _  \___  _______  |  | _____    ____   ____ |  |__   ____    ,_ o` + "\n" +
		`   /  /_\  \  \/ /\__  \ |  | \__  \  /    \_/ ___\|  |  \_/ __ \   / //\,` + "\n" +
		`  /    |    \   /  / __ \|  |__/ __ \|   |  \  \___|   Y  \  ___/    \>> |` + "\n" +
		`  \____|__  /\_/  (____  /____(____  /___|  /\___  >___|  /\___  >    \\` + "\n" +
		`          \/           \/          \/     \/     \/     \/     \/`
)

var (
	stakingPortName = fmt.Sprintf("%s-staking", constants.AppName)
	httpPortName    = fmt.Sprintf("%s-http", constants.AppName)
)

// main is the primary entry point to Avalanche.
func main() {
	// Err is set based on the CLI arguments
	if Err != nil {
		fmt.Printf("parsing parameters returned with error %s\n", Err)
		return
	}

	factory := logging.NewFactory(Config.LoggingConfig)
	defer factory.Close()

	log, err := factory.Make()
	if err != nil {
		fmt.Printf("starting logger failed with: %s\n", err)
		return
	}
	fmt.Println(header)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic from", r)
		}
	}()

	defer func() {
		if err := Config.DB.Close(); err != nil {
			log.Warn("failed to close the node's DB: %s", err)
		}
		log.StopOnPanic()
		log.Stop()
	}()

	// Track if sybil control is enforced
	if !Config.EnableStaking && Config.EnableP2PTLS {
		log.Warn("Staking is disabled. Sybil control is not enforced.")
	}
	if !Config.EnableStaking && !Config.EnableP2PTLS {
		log.Warn("Staking and p2p encryption are disabled. Packet spoofing is possible.")
	}

	// Check if transaction signatures should be checked
	if !Config.EnableCrypto {
		log.Warn("transaction signatures are not being checked")
	}
	crypto.EnableCrypto = Config.EnableCrypto

	if err := Config.ConsensusParams.Valid(); err != nil {
		log.Error("consensus parameters are invalid: %s", err)
		return
	}

	// Track if assertions should be executed
	if Config.LoggingConfig.Assertions {
		log.Debug("assertions are enabled. This may slow down execution")
	}

	// SupportsNAT() for NoRouter is false.
	// Which means we tried to perform a NAT activity but we were not successful.
	if Config.AttemptedNATTraversal && !Config.Nat.SupportsNAT() {
		log.Error("UPnP or NAT-PMP router attach failed, you may not be listening publicly," +
			" please confirm the settings in your router")
	}

	mapper := nat.NewPortMapper(log, Config.Nat)
	defer mapper.UnmapAllPorts()

	// Open staking port we want for NAT Traversal to have the external port
	// (Config.StakingIP.Port) to connect to our internal listening port
	// (Config.InternalStakingPort) which should be the same in most cases.
	mapper.Map(
		"TCP",
		Config.StakingIP.IP().Port,
		Config.StakingIP.IP().Port,
		stakingPortName,
		&Config.StakingIP,
		Config.DynamicUpdateDuration,
	)

	// Open the HTTP port iff the HTTP server is not listening on localhost
	if Config.HTTPHost != "127.0.0.1" && Config.HTTPHost != "localhost" {
		// For NAT Traversal we want to route from the external port
		// (Config.ExternalHTTPPort) to our internal port (Config.HTTPPort)
		mapper.Map(
			"TCP",
			Config.HTTPPort,
			Config.HTTPPort,
			httpPortName,
			nil,
			Config.DynamicUpdateDuration,
		)
	}

	// Regularly updates our public IP (or does nothing, if configured that way)
	externalIPUpdater := dynamicip.NewDynamicIPManager(
		Config.DynamicPublicIPResolver,
		Config.DynamicUpdateDuration,
		log,
		&Config.StakingIP,
	)
	defer externalIPUpdater.Stop()

	log.Info("Node started with IP: %s", Config.StakingIP.IP().String())

	log.Debug("initializing node state")
	node := node.Node{}
	if err := node.Initialize(&Config, log, factory); err != nil {
		log.Error("error initializing node state: %s", err)
		return
	}

	defer node.Shutdown()

	log.Debug("dispatching node handlers")
	err = node.Dispatch()
	log.Debug("node dispatching returned with %s", err)
}
