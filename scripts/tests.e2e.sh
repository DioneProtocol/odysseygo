#!/usr/bin/env bash
set -e
set -o nounset
set -o pipefail

# e.g.,
# ./scripts/build.sh
# ./scripts/tests.e2e.sh ./build/odysseygo
if ! [[ "$0" =~ scripts/tests.e2e.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

ODYSSEYGO_PATH="${1-}"
if [[ -z "${ODYSSEYGO_PATH}" ]]; then
  echo "Missing ODYSSEYGO_PATH argument!"
  echo "Usage: ${0} [ODYSSEYGO_PATH]" >> /dev/stderr
  exit 255
fi

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"
# While CGO_ENABLED doesn't need to be explicitly set, it produces a much more
# clear error due to the default value change in go1.20.
export CGO_ENABLED=1

#################################
# download odyssey-network-runner
# https://github.com/DioneProtocol/odyssey-network-runner
# TODO: migrate to upstream odyssey-network-runner
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
NETWORK_RUNNER_VERSION=1.3.5-rc.0
DOWNLOAD_PATH=/tmp/odyssey-network-runner.tar.gz
DOWNLOAD_URL="https://github.com/DioneProtocol/odyssey-network-runner/releases/download/v${NETWORK_RUNNER_VERSION}/odyssey-network-runner_${NETWORK_RUNNER_VERSION}_${GOOS}_${GOARCH}.tar.gz"

rm -f ${DOWNLOAD_PATH}
rm -f /tmp/odyssey-network-runner

echo "downloading odyssey-network-runner ${NETWORK_RUNNER_VERSION} at ${DOWNLOAD_URL} to ${DOWNLOAD_PATH}"
curl --fail -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

echo "extracting downloaded odyssey-network-runner"
tar xzvf ${DOWNLOAD_PATH} -C /tmp
/tmp/odyssey-network-runner -h

GOPATH="$(go env GOPATH)"
PATH="${GOPATH}/bin:${PATH}"

#################################
echo "building e2e.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.1.4
ACK_GINKGO_RC=true ginkgo build ./tests/e2e
./tests/e2e/e2e.test --help

#################################
# run "odyssey-network-runner" server
echo "launch odyssey-network-runner in the background"
/tmp/odyssey-network-runner \
server \
--log-level debug \
--port=":12342" \
--disable-grpc-gateway &
PID=${!}

#################################
echo "running e2e tests against the local cluster with ${ODYSSEYGO_PATH}"
./tests/e2e/e2e.test \
--ginkgo.v \
--log-level debug \
--network-runner-grpc-endpoint="0.0.0.0:12342" \
--network-runner-odysseygo-path=${ODYSSEYGO_PATH} \
--network-runner-odysseygo-log-level="WARN" \
--test-keys-file=tests/test.insecure.secp256k1.keys \
&& EXIT_CODE=$? || EXIT_CODE=$?

kill ${PID}

if [[ ${EXIT_CODE} -gt 0 ]]; then
  echo "FAILURE with exit code ${EXIT_CODE}"
  exit ${EXIT_CODE}
else
  echo "ALL SUCCESS!"
fi
