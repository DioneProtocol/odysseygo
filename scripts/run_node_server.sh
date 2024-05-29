#!/usr/bin/env bash

print_usage() {
  printf "Usage: run_node [OPTIONS]

  Run odysseygo node

  Options:

    -n Node name to run from run_config file 'node' section
"
}

while getopts n: flag; do
  case "${flag}" in
    n) node_name=${OPTARG} ;;
    *) print_usage
      exit 1 ;;
  esac
done

# Odysseygo root folder
ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
cd $ODYSSEY_PATH
# Load the constants
source "$ODYSSEY_PATH"/scripts/constants.sh
# Parse run_config and load nodes info
source "$ODYSSEY_PATH"/scripts/parse_run_config.sh

if [ -n "$node_name" ]; then
  node_index=""
  for i in ${!NODE_NAMES[@]}; do
    if [ ${NODE_NAMES[$i]} == $node_name ]; then
      node_index=$i
    fi
  done

  if [ -z $node_index ]; then
    echo error: Node with name: $node_name not found
    exit 1
  fi
else
  print_usage
  exit 1
fi
# Download dependencies
echo "Downloading dependencies..."
go mod download

# Build odysseygo
"$ODYSSEY_PATH"/scripts/build_odyssey.sh

# Exit if the OdysseyGo binary is not created successfully
if [[ -f "$odysseygo_path" ]]; then
        echo "Build Successful"
else
        echo "Build failure" >&2
        exit 1
fi

run_node() {
echo -ne "[Unit]\nDescription=Service for odyssey network's node\nAfter=network.target\nAfter=syslog.target\n\n[Service]\nType=simple\nUser=root\n
ExecStart='$odysseygo_path' --public-ip=${PUBLIC_IPS[$1]} --http-host=${PUBLIC_IPS[$1]} --http-port=${HTTP_PORTS[$1]} --staking-port=${STAKING_PORTS[$1]} --db-dir=${DB_DIRS[$1]}/${NODE_NAMES[$1]} --log-dir=${DB_DIRS[$1]}/${NODE_NAMES[$1]}/logs --chain-config-dir="$ODYSSEY_PATH"/scripts/configs/archive --network-id=${NETWORK_IDS[$1]} --http-allowed-hosts=* --bootstrap-ips=$BOOTSTRAP_HOST_IPS --bootstrap-ids=$BOOTSTRAP_IDS
Restart=on-failure\nRestartSec=5\nPIDFile=/tmp/node-$node_name.pid\n\n[Install]\nWantedBy=default.target" > /etc/systemd/system/node-$node_name.service
systemctl daemon-reload
systemctl start node-$node_name.service
systemctl enable node-$node_name.service
}

run_node $node_index

wait
