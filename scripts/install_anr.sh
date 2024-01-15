#!/usr/bin/env bash

set -euo pipefail

# Odyssey root directory
ODYSSEY_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

#################################
# download odyssey-network-runner
# https://github.com/DioneProtocol/odyssey-network-runner
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
NETWORK_RUNNER_VERSION=1.7.0
anr_workdir=${ANR_WORKDIR:-"/tmp"}
DOWNLOAD_PATH=${anr_workdir}/odyssey-network-runner-v${NETWORK_RUNNER_VERSION}.tar.gz
DOWNLOAD_URL="https://github.com/DioneProtocol/odyssey-network-runner/releases/download/v${NETWORK_RUNNER_VERSION}/odyssey-network-runner_${NETWORK_RUNNER_VERSION}_${GOOS}_${GOARCH}.tar.gz"
echo "Installing odyssey-network-runner ${NETWORK_RUNNER_VERSION} to ${anr_workdir}/odyssey-network-runner"

# download only if not already downloaded
if [ ! -f "$DOWNLOAD_PATH" ]; then
  echo "downloading odyssey-network-runner ${NETWORK_RUNNER_VERSION} at ${DOWNLOAD_URL} to ${DOWNLOAD_PATH}"
  curl --fail -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}
else
  echo "odyssey-network-runner ${NETWORK_RUNNER_VERSION} already downloaded at ${DOWNLOAD_PATH}"
fi

rm -f ${anr_workdir}/odyssey-network-runner

echo "extracting downloaded odyssey-network-runner"
tar xzvf ${DOWNLOAD_PATH} -C ${anr_workdir}
${anr_workdir}/odyssey-network-runner -h
