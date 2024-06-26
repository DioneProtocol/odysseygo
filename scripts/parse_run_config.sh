#!/usr/bin/env bash

RAW_BOOTSTRAP_IDS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].id' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_IDS=""
BOOTSTRAP_NODE_IDS=()
for line in $RAW_BOOTSTRAP_IDS; do
  BOOTSTRAP_IDS+="${line},"
  BOOTSTRAP_NODE_IDS+=($line)
done
BOOTSTRAP_IDS=$(echo $BOOTSTRAP_IDS | sed 's/.$//')

RAW_BOOTSTRAP_STAKING_PORTS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].stakingPort' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_STAKING_PORTS=()
for line in $RAW_BOOTSTRAP_STAKING_PORTS; do
  BOOTSTRAP_STAKING_PORTS+=("$line")
done

RAW_STAKING_PORTS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes[].stakingPort' | sed -e 's/^"//' -e 's/"$//')
STAKING_PORTS=()
for line in $RAW_STAKING_PORTS; do
  STAKING_PORTS+=("$line")
done

RAW_BOOTSTRAP_HTTP_PORTS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].httpPort' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_HTTP_PORTS=()
for line in $RAW_BOOTSTRAP_HTTP_PORTS; do
  BOOTSTRAP_HTTP_PORTS+=("$line")
done

RAW_HTTP_PORTS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes[].httpPort' | sed -e 's/^"//' -e 's/"$//')
HTTP_PORTS=()
for line in $RAW_HTTP_PORTS; do
  HTTP_PORTS+=("$line")
done

RAW_BOOTSTRAP_IPS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].ip' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_IPS=""
BOOTSTRAP_PUBLIC_IPS=()
i=0
for line in $RAW_BOOTSTRAP_IPS; do
  BOOTSTRAP_IPS+="${line}:${BOOTSTRAP_STAKING_PORTS[$i]},"
  BOOTSTRAP_PUBLIC_IPS+=($line)
  i=$(($i + 1))
done
BOOTSTRAP_IPS=$(echo $BOOTSTRAP_IPS | sed 's/.$//')

RAW_BOOTSTRAP_HOST_IPS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].hostIp' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_HOST_IPS=""
i=0
for line in $RAW_BOOTSTRAP_HOST_IPS; do
  BOOTSTRAP_HOST_IPS+="${line}:${BOOTSTRAP_STAKING_PORTS[$i]},"
  i=$(($i + 1))
done
BOOTSTRAP_HOST_IPS=$(echo $BOOTSTRAP_HOST_IPS | sed 's/.$//')

RAW_IPS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes[].ip' | sed -e 's/^"//' -e 's/"$//')
PUBLIC_IPS=()
for line in $RAW_BOOTSTRAP_IPS; do
  PUBLIC_IPS+=($line)
done

RAW_BOOTSTRAP_DB_DIRS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].dbDir' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_DB_DIRS=()
for line in $RAW_BOOTSTRAP_DB_DIRS; do
  BOOTSTRAP_DB_DIRS+=("$ODYSSEY_PATH$line")
done

RAW_DB_DIRS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes[].dbDir' | sed -e 's/^"//' -e 's/"$//')
DB_DIRS=()
for line in $RAW_DB_DIRS; do
  DB_DIRS+=("$ODYSSEY_PATH$line")
done

RAW_BOOTSTRAP_NODE_NAMES=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes | keys[]' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_NODE_NAMES=()
for line in $RAW_BOOTSTRAP_NODE_NAMES; do
  BOOTSTRAP_NODE_NAMES+=("$line")
done

RAW_NODE_NAMES=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes | keys[]' | sed -e 's/^"//' -e 's/"$//')
NODE_NAMES=()
for line in $RAW_NODE_NAMES; do
  NODE_NAMES+=("$line")
done

RAW_BOOTSTRAP_NETWORK_IDS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].networkId' | sed -e 's/^"//' -e 's/"$//')
BOOTSTRAP_NETWORK_IDS=()
for line in $RAW_BOOTSTRAP_NETWORK_IDS; do
  BOOTSTRAP_NETWORK_IDS+=("$line")
done

RAW_NETWORK_IDS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.nodes[].networkId' | sed -e 's/^"//' -e 's/"$//')
NETWORK_IDS=()
for line in $RAW_NETWORK_IDS; do
  NETWORK_IDS+=("$line")
done

RAW_TLS_KEY_PATHS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].tlsKeyFilePath' | sed -e 's/^"//' -e 's/"$//')
TLS_KEY_PATHS=()
for line in $RAW_TLS_KEY_PATHS; do
  TLS_KEY_PATHS+=("$ODYSSEY_PATH$line")
done

RAW_TLS_CERT_PATHS=$(cat "$ODYSSEY_PATH"/scripts/run_config.json | jq '.bootstrapNodes[].tlsCertFilePath' | sed -e 's/^"//' -e 's/"$//')
TLS_CERT_PATHS=()
for line in $RAW_TLS_CERT_PATHS; do
  TLS_CERT_PATHS+=("$ODYSSEY_PATH$line")
done
