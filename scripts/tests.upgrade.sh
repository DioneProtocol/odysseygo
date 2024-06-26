#!/usr/bin/env bash

set -euo pipefail

# e.g.,
# ./scripts/build.sh
# ./scripts/tests.upgrade.sh 1.7.16 ./build/odysseygo
if ! [[ "$0" =~ scripts/tests.upgrade.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
  echo "Missing version argument!"
  echo "Usage: ${0} [VERSION] [NEW-BINARY]" >>/dev/stderr
  exit 255
fi

NEW_BINARY="${2:-}"
if [[ -z "${NEW_BINARY}" ]]; then
  echo "Missing new binary path argument!"
  echo "Usage: ${0} [VERSION] [NEW-BINARY]" >>/dev/stderr
  exit 255
fi

#################################
# download odysseygo
# https://github.com/DioneProtocol/odysseygo/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
DOWNLOAD_URL=https://github.com/DioneProtocol/odysseygo/releases/download/v${VERSION}/odysseygo-linux-${GOARCH}-v${VERSION}.tar.gz
DOWNLOAD_PATH=/tmp/odysseygo.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  DOWNLOAD_URL=https://github.com/DioneProtocol/odysseygo/releases/download/v${VERSION}/odysseygo-macos-v${VERSION}.zip
  DOWNLOAD_PATH=/tmp/odysseygo.zip
fi

rm -f ${DOWNLOAD_PATH}
rm -rf /tmp/odysseygo-v${VERSION}
rm -rf /tmp/odysseygo-build

echo "downloading odysseygo ${VERSION} at ${DOWNLOAD_URL}"
curl -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

echo "extracting downloaded odysseygo"
if [[ ${GOOS} == "linux" ]]; then
  tar xzvf ${DOWNLOAD_PATH} -C /tmp
elif [[ ${GOOS} == "darwin" ]]; then
  unzip ${DOWNLOAD_PATH} -d /tmp/odysseygo-build
  mv /tmp/odysseygo-build/build /tmp/odysseygo-v${VERSION}
fi
find /tmp/odysseygo-v${VERSION}

#################################
echo "installing odyssey-network-runner"
ANR_WORKDIR="/tmp"
./scripts/install_anr.sh

# Sourcing constants.sh ensures that the necessary CGO flags are set to
# build the portable version of BLST. Without this, ginkgo may fail to
# build the test binary if run on a host (e.g. github worker) that lacks
# the instructions to build non-portable BLST.
source ./scripts/constants.sh

#################################
echo "building upgrade.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.1.4
ACK_GINKGO_RC=true ginkgo build ./tests/upgrade
./tests/upgrade/upgrade.test --help

#################################
# run "odyssey-network-runner" server
echo "launch odyssey-network-runner in the background"
$ANR_WORKDIR/odyssey-network-runner \
  server \
  --log-level debug \
  --port=":12340" \
  --disable-grpc-gateway &
PID=${!}

#################################
# By default, it runs all upgrade test cases!
echo "running upgrade tests against the local cluster with ${NEW_BINARY}"
./tests/upgrade/upgrade.test \
  --ginkgo.v \
  --log-level debug \
  --network-runner-grpc-endpoint="0.0.0.0:12340" \
  --network-runner-odysseygo-path=/tmp/odysseygo-v${VERSION}/odysseygo \
  --network-runner-odysseygo-path-to-upgrade=${NEW_BINARY} \
  --network-runner-odysseygo-log-level="WARN" || EXIT_CODE=$?

# "e2e.test" already terminates the cluster
# just in case tests are aborted, manually terminate them again
pkill -P ${PID} || true
kill -2 ${PID}

if [[ "${EXIT_CODE:-}" -gt 0 ]]; then
  echo "FAILURE with exit code ${EXIT_CODE}"
  exit ${EXIT_CODE}
else
  echo "ALL SUCCESS!"
fi
