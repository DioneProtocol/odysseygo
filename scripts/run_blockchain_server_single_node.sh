#!/usr/bin/env bash

# Odysseygo root folder
ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
cd $ODYSSEY_PATH
# Load the constants
source "$ODYSSEY_PATH"/scripts/constants.sh

if [ -z $1 ]; then
    echo "Specify node index"
    exit -1
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

# Parse run_config and load nodes info
source "$ODYSSEY_PATH"/scripts/parse_run_config.sh

if [ -z "${BOOTSTRAP_PUBLIC_IPS[$1]}" ]; then
    echo "Node with index $1 not found"
    exit -1
fi

run_node() {
echo -ne "[Unit]\nDescription=Service for odyssey network's validator\nAfter=network.target\nAfter=syslog.target\n\n[Service]\nType=simple\nUser=root\n
ExecStart='$odysseygo_path' --public-ip=${BOOTSTRAP_PUBLIC_IPS[$1]} --http-host=${BOOTSTRAP_PUBLIC_IPS[$1]} --http-port=${BOOTSTRAP_HTTP_PORTS[$1]} --staking-port=${BOOTSTRAP_STAKING_PORTS[$1]} --db-dir=${BOOTSTRAP_DB_DIRS[$1]}/${BOOTSTRAP_NODE_NAMES[$1]} --log-dir=${BOOTSTRAP_DB_DIRS[$1]}/${BOOTSTRAP_NODE_NAMES[$1]}/logs --chain-config-dir="$ODYSSEY_PATH"/scripts/configs/validators --network-id=${BOOTSTRAP_NETWORK_IDS[$1]} --http-allowed-hosts=* --bootstrap-ips=$BOOTSTRAP_HOST_IPS --bootstrap-ids=$BOOTSTRAP_IDS --staking-tls-cert-file=${TLS_CERT_PATHS[$1]} --staking-tls-key-file=${TLS_KEY_PATHS[$1]}
Restart=on-failure\nRestartSec=5\nPIDFile=/tmp/odyssey$1.pid\n\n[Install]\nWantedBy=default.target" > /etc/systemd/system/node-validator-$1.service

systemctl daemon-reload
systemctl start node-validator-$1.service
systemctl enable node-validator-$1.service
}

run_node $1 &

wait
