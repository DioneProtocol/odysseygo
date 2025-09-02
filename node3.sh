./build/odysseygo \
  --public-ip-resolution-service opendns \
  --network-id=123456 \
  --http-port=9654 \
  --staking-port=9655 \
  --db-dir=./data/node3 \
  --log-dir=./data/node3/logs \
  --data-dir=./data/node3/.odysseygo \
  --genesis-file=./genesis/genesis_custom.json \
  --bootstrap-beacon-connection-timeout=5m \
  --staking-signer-key-file=data/node3/.odysseygo/staking/signer.key \
  --staking-tls-cert-file=data/node3/.odysseygo/staking/staker.crt \
  --staking-tls-key-file=data/node3/.odysseygo/staking/staker.key \
  --bootstrap-ips=127.0.0.1:9651,127.0.0.1:9653,127.0.0.1:9657,127.0.0.1:9659 \
  --bootstrap-ids=NodeID-DFRpSjRPPByYm2jtX74aCpdrcstAfpAwR,NodeID-4tGmJnmgLYDq5aL4eYTD94E3Dr637kdtD,NodeID-C9mbTTgDLwbuYJnKdQAJDwrbSvHKLwctc,NodeID-7oowrokMRzbUoDRzytLsVhSny3mPTCrJi \
  --chain-config-dir=./scripts/configs/validators/D \
  # --log-level=debug \
  # --consensus-shutdown-timeout=5s \
  # --bootstrap-retry-enabled=true \