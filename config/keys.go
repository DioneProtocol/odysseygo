// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package config

// #nosec G101
const (
	ConfigFileKey                                      = "config-file"
	ConfigContentKey                                   = "config-file-content"
	ConfigContentTypeKey                               = "config-file-content-type"
	VersionKey                                         = "version"
	GenesisConfigFileKey                               = "genesis"
	GenesisConfigContentKey                            = "genesis-content"
	NetworkNameKey                                     = "network-id"
	TxFeeKey                                           = "tx-fee"
	CreateAssetTxFeeKey                                = "create-asset-tx-fee"
	CreateSubnetTxFeeKey                               = "create-subnet-tx-fee"
	CreateBlockchainTxFeeKey                           = "create-blockchain-tx-fee"
	UptimeRequirementKey                               = "uptime-requirement"
	MinValidatorStakeKey                               = "min-validator-stake"
	MaxValidatorStakeKey                               = "max-validator-stake"
	MinDelegatorStakeKey                               = "min-delegator-stake"
	MinDelegatorFeeKey                                 = "min-delegation-fee"
	MinStakeDurationKey                                = "min-stake-duration"
	MaxStakeDurationKey                                = "max-stake-duration"
	StakeMaxConsumptionRateKey                         = "stake-max-consumption-rate"
	StakeMinConsumptionRateKey                         = "stake-min-consumption-rate"
	StakeMintingPeriodKey                              = "stake-minting-period"
	StakeSupplyCapKey                                  = "stake-supply-cap"
	AssertionsEnabledKey                               = "assertions-enabled"
	SignatureVerificationEnabledKey                    = "signature-verification-enabled"
	DBTypeKey                                          = "db-type"
	DBPathKey                                          = "db-dir"
	DBConfigFileKey                                    = "db-config-file"
	DBConfigContentKey                                 = "db-config-file-content"
	PublicIPKey                                        = "public-ip"
	DynamicUpdateDurationKey                           = "dynamic-update-duration"
	DynamicPublicIPResolverKey                         = "dynamic-public-ip"
	InboundConnUpgradeThrottlerCooldownKey             = "inbound-connection-throttling-cooldown"
	InboundThrottlerMaxConnsPerSecKey                  = "inbound-connection-throttling-max-conns-per-sec"
	OutboundConnectionThrottlingRps                    = "outbound-connection-throttling-rps"
	OutboundConnectionTimeout                          = "outbound-connection-timeout"
	HTTPHostKey                                        = "http-host"
	HTTPPortKey                                        = "http-port"
	HTTPSEnabledKey                                    = "http-tls-enabled"
	HTTPSKeyFileKey                                    = "http-tls-key-file"
	HTTPSKeyContentKey                                 = "http-tls-key-file-content"
	HTTPSCertFileKey                                   = "http-tls-cert-file"
	HTTPSCertContentKey                                = "http-tls-cert-file-content"
	HTTPAllowedOrigins                                 = "http-allowed-origins"
	HTTPShutdownTimeoutKey                             = "http-shutdown-timeout"
	HTTPShutdownWaitKey                                = "http-shutdown-wait"
	APIAuthRequiredKey                                 = "api-auth-required"
	APIAuthPasswordKey                                 = "api-auth-password"
	APIAuthPasswordFileKey                             = "api-auth-password-file"
	BootstrapIPsKey                                    = "bootstrap-ips"
	BootstrapIDsKey                                    = "bootstrap-ids"
	StakingPortKey                                     = "staking-port"
	StakingEnabledKey                                  = "staking-enabled"
	StakingEphemeralCertEnabledKey                     = "staking-ephemeral-cert-enabled"
	StakingKeyPathKey                                  = "staking-tls-key-file"
	StakingKeyContentKey                               = "staking-tls-key-file-content"
	StakingCertPathKey                                 = "staking-tls-cert-file"
	StakingCertContentKey                              = "staking-tls-cert-file-content"
	StakingDisabledWeightKey                           = "staking-disabled-weight"
	NetworkInitialTimeoutKey                           = "network-initial-timeout"
	NetworkMinimumTimeoutKey                           = "network-minimum-timeout"
	NetworkMaximumTimeoutKey                           = "network-maximum-timeout"
	NetworkMaximumInboundTimeoutKey                    = "network-maximum-inbound-timeout"
	NetworkTimeoutHalflifeKey                          = "network-timeout-halflife"
	NetworkTimeoutCoefficientKey                       = "network-timeout-coefficient"
	NetworkHealthMinPeersKey                           = "network-health-min-conn-peers"
	NetworkHealthMaxTimeSinceMsgReceivedKey            = "network-health-max-time-since-msg-received"
	NetworkHealthMaxTimeSinceMsgSentKey                = "network-health-max-time-since-msg-sent"
	NetworkHealthMaxPortionSendQueueFillKey            = "network-health-max-portion-send-queue-full"
	NetworkHealthMaxSendFailRateKey                    = "network-health-max-send-fail-rate"
	NetworkHealthMaxOutstandingDurationKey             = "network-health-max-outstanding-request-duration"
	NetworkPeerListNumValidatorIPsKey                  = "network-peer-list-num-validator-ips"
	NetworkPeerListValidatorGossipSizeKey              = "network-peer-list-validator-gossip-size"
	NetworkPeerListNonValidatorGossipSizeKey           = "network-peer-list-non-validator-gossip-size"
	NetworkPeerListPeersGossipSizeKey                  = "network-peer-list-peers-gossip-size"
	NetworkPeerListGossipFreqKey                       = "network-peer-list-gossip-frequency"
	NetworkInitialReconnectDelayKey                    = "network-initial-reconnect-delay"
	NetworkReadHandshakeTimeoutKey                     = "network-read-handshake-timeout"
	NetworkPingTimeoutKey                              = "network-ping-timeout"
	NetworkPingFrequencyKey                            = "network-ping-frequency"
	NetworkMaxReconnectDelayKey                        = "network-max-reconnect-delay"
	NetworkCompressionEnabledKey                       = "network-compression-enabled"
	NetworkMaxClockDifferenceKey                       = "network-max-clock-difference"
	NetworkAllowPrivateIPsKey                          = "network-allow-private-ips"
	NetworkRequireValidatorToConnectKey                = "network-require-validator-to-connect"
	NetworkPeerReadBufferSizeKey                       = "network-peer-read-buffer-size"
	NetworkPeerWriteBufferSizeKey                      = "network-peer-write-buffer-size"
	BenchlistFailThresholdKey                          = "benchlist-fail-threshold"
	BenchlistDurationKey                               = "benchlist-duration"
	BenchlistMinFailingDurationKey                     = "benchlist-min-failing-duration"
	BuildDirKey                                        = "build-dir"
	LogsDirKey                                         = "log-dir"
	LogLevelKey                                        = "log-level"
	LogDisplayLevelKey                                 = "log-display-level"
	LogDisplayHighlightKey                             = "log-display-highlight"
	LogDisableDisplayPluginLogsKey                     = "log-disable-display-plugin-logs"
	SnowSampleSizeKey                                  = "snow-sample-size"
	SnowQuorumSizeKey                                  = "snow-quorum-size"
	SnowVirtuousCommitThresholdKey                     = "snow-virtuous-commit-threshold"
	SnowRogueCommitThresholdKey                        = "snow-rogue-commit-threshold"
	SnowAvalancheNumParentsKey                         = "snow-avalanche-num-parents"
	SnowAvalancheBatchSizeKey                          = "snow-avalanche-batch-size"
	SnowConcurrentRepollsKey                           = "snow-concurrent-repolls"
	SnowOptimalProcessingKey                           = "snow-optimal-processing"
	SnowMaxProcessingKey                               = "snow-max-processing"
	SnowMaxTimeProcessingKey                           = "snow-max-time-processing"
	WhitelistedSubnetsKey                              = "whitelisted-subnets"
	AdminAPIEnabledKey                                 = "api-admin-enabled"
	InfoAPIEnabledKey                                  = "api-info-enabled"
	KeystoreAPIEnabledKey                              = "api-keystore-enabled"
	MetricsAPIEnabledKey                               = "api-metrics-enabled"
	HealthAPIEnabledKey                                = "api-health-enabled"
	IpcAPIEnabledKey                                   = "api-ipcs-enabled"
	IpcsChainIDsKey                                    = "ipcs-chain-ids"
	IpcsPathKey                                        = "ipcs-path"
	MeterVMsEnabledKey                                 = "meter-vms-enabled"
	ConsensusGossipFrequencyKey                        = "consensus-gossip-frequency"
	ConsensusGossipAcceptedFrontierValidatorSizeKey    = "consensus-accepted-frontier-gossip-validator-size"
	ConsensusGossipAcceptedFrontierNonValidatorSizeKey = "consensus-accepted-frontier-gossip-non-validator-size"
	ConsensusGossipAcceptedFrontierPeerSizeKey         = "consensus-accepted-frontier-gossip-peer-size"
	ConsensusGossipOnAcceptValidatorSizeKey            = "consensus-on-accept-gossip-validator-size"
	ConsensusGossipOnAcceptNonValidatorSizeKey         = "consensus-on-accept-gossip-non-validator-size"
	ConsensusGossipOnAcceptPeerSizeKey                 = "consensus-on-accept-gossip-peer-size"
	AppGossipValidatorSizeKey                          = "consensus-app-gossip-validator-size"
	AppGossipNonValidatorSizeKey                       = "consensus-app-gossip-non-validator-size"
	AppGossipPeerSizeKey                               = "consensus-app-gossip-peer-size"
	ConsensusShutdownTimeoutKey                        = "consensus-shutdown-timeout"
	FdLimitKey                                         = "fd-limit"
	IndexEnabledKey                                    = "index-enabled"
	IndexAllowIncompleteKey                            = "index-allow-incomplete"
	ResetProposerVMHeightIndexKey                      = "reset-proposervm-height-index"
	RouterHealthMaxDropRateKey                         = "router-health-max-drop-rate"
	RouterHealthMaxOutstandingRequestsKey              = "router-health-max-outstanding-requests"
	HealthCheckFreqKey                                 = "health-check-frequency"
	HealthCheckAveragerHalflifeKey                     = "health-check-averager-halflife"
	RetryBootstrapKey                                  = "bootstrap-retry-enabled"
	RetryBootstrapWarnFrequencyKey                     = "bootstrap-retry-warn-frequency"
	PluginModeKey                                      = "plugin-mode-enabled"
	BootstrapBeaconConnectionTimeoutKey                = "bootstrap-beacon-connection-timeout"
	BootstrapMaxTimeGetAncestorsKey                    = "boostrap-max-time-get-ancestors"
	BootstrapAncestorsMaxContainersSentKey             = "bootstrap-ancestors-max-containers-sent"
	BootstrapAncestorsMaxContainersReceivedKey         = "bootstrap-ancestors-max-containers-received"
	ChainConfigDirKey                                  = "chain-config-dir"
	ChainConfigContentKey                              = "chain-config-content"
	SubnetConfigDirKey                                 = "subnet-config-dir"
	SubnetConfigContentKey                             = "subnet-config-content"
	ProfileDirKey                                      = "profile-dir"
	ProfileContinuousEnabledKey                        = "profile-continuous-enabled"
	ProfileContinuousFreqKey                           = "profile-continuous-freq"
	ProfileContinuousMaxFilesKey                       = "profile-continuous-max-files"
	InboundThrottlerAtLargeAllocSizeKey                = "throttler-inbound-at-large-alloc-size"
	InboundThrottlerVdrAllocSizeKey                    = "throttler-inbound-validator-alloc-size"
	InboundThrottlerNodeMaxAtLargeBytesKey             = "throttler-inbound-node-max-at-large-bytes"
	InboundThrottlerMaxProcessingMsgsPerNodeKey        = "throttler-inbound-node-max-processing-msgs"
	InboundThrottlerBandwidthRefillRateKey             = "throttler-inbound-bandwidth-refill-rate"
	InboundThrottlerBandwidthMaxBurstSizeKey           = "throttler-inbound-bandwidth-max-burst-size"
	OutboundThrottlerAtLargeAllocSizeKey               = "throttler-outbound-at-large-alloc-size"
	OutboundThrottlerVdrAllocSizeKey                   = "throttler-outbound-validator-alloc-size"
	OutboundThrottlerNodeMaxAtLargeBytesKey            = "throttler-outbound-node-max-at-large-bytes"
	UptimeMetricFreqKey                                = "uptime-metric-freq"
	VMAliasesFileKey                                   = "vm-aliases-file"
	VMAliasesContentKey                                = "vm-aliases-file-content"
)
