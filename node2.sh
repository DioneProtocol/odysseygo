./build/odysseygo \
  --public-ip-resolution-service opendns \
  --network-id=123456 \
  --http-port=9652 \
  --staking-port=9653 \
  --db-dir=./data/node2 \
  --log-dir=./data/node2/logs \
  --data-dir=./data/node2/.odysseygo \
  --genesis-file=./genesis/genesis_custom.json \
  --bootstrap-beacon-connection-timeout=5m \
  --staking-signer-key-file=data/node2/.odysseygo/staking/signer.key \
  --staking-tls-cert-file=data/node2/.odysseygo/staking/staker.crt \
  --staking-tls-key-file=data/node2/.odysseygo/staking/staker.key \
  --bootstrap-ips=127.0.0.1:9651,127.0.0.1:9655,127.0.0.1:9657,127.0.0.1:9659 \
  --bootstrap-ids=NodeID-DFRpSjRPPByYm2jtX74aCpdrcstAfpAwR,NodeID-2u2AdhwmTZDWfD4C3ppu52uwzcqfgtz9X,NodeID-C9mbTTgDLwbuYJnKdQAJDwrbSvHKLwctc,NodeID-7oowrokMRzbUoDRzytLsVhSny3mPTCrJi \
  --chain-config-dir=./scripts/configs/validators/D \
  # --log-level=debug \
  # --consensus-shutdown-timeout=5s \
  # --bootstrap-retry-enabled=true \