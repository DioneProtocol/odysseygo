./build/odysseygo \
  --public-ip-resolution-service opendns \
  --network-id=123456 \
  --http-port=9660 \
  --staking-port=9661 \
  --db-dir=./data/node6 \
  --log-dir=./data/node6/logs \
  --data-dir=./data/node6/.odysseygo \
  --genesis-file=./genesis/genesis_custom.json \
  --bootstrap-beacon-connection-timeout=5m \
  --bootstrap-ips=127.0.0.1:9651,127.0.0.1:9653,127.0.0.1:9655,127.0.0.1:9657 \
  --bootstrap-ids=NodeID-DFRpSjRPPByYm2jtX74aCpdrcstAfpAwR,NodeID-4tGmJnmgLYDq5aL4eYTD94E3Dr637kdtD,NodeID-2u2AdhwmTZDWfD4C3ppu52uwzcqfgtz9X,NodeID-C9mbTTgDLwbuYJnKdQAJDwrbSvHKLwctc \
  --chain-config-dir=./scripts/configs/archive/D \
#   --log-level=debug \
#   --consensus-shutdown-timeout=5s \
#   --bootstrap-retry-enabled=true \