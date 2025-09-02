./build/odysseygo \
  --public-ip-resolution-service opendns \
  --network-id=123456 \
  --http-port=9650 \
  --staking-port=9651 \
  --db-dir=./data/node1 \
  --log-dir=./data/node1/logs \
  --data-dir=./data/node1/.odysseygo \
  --genesis-file=./genesis/genesis_custom.json \
  --bootstrap-beacon-connection-timeout=5m \
  --staking-signer-key-file=data/node1/.odysseygo/staking/signer.key \
  --staking-tls-cert-file=data/node1/.odysseygo/staking/staker.crt \
  --staking-tls-key-file=data/node1/.odysseygo/staking/staker.key \
  --bootstrap-ips=127.0.0.1:9653,127.0.0.1:9655,127.0.0.1:9657,127.0.0.1:9659 \
  --bootstrap-ids=NodeID-4tGmJnmgLYDq5aL4eYTD94E3Dr637kdtD,NodeID-2u2AdhwmTZDWfD4C3ppu52uwzcqfgtz9X,NodeID-C9mbTTgDLwbuYJnKdQAJDwrbSvHKLwctc,NodeID-7oowrokMRzbUoDRzytLsVhSny3mPTCrJi \
  --chain-config-dir=./scripts/configs/validators/D \
  # --log-level=debug \
  # --consensus-shutdown-timeout=5s \
  # --bootstrap-retry-enabled=true \