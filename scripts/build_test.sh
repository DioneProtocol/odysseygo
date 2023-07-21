#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Directory above this script
ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$ODYSSEY_PATH"/scripts/constants.sh

go test -race -timeout="120s" -coverprofile="coverage.out" -covermode="atomic" $(go list ./... | grep -v /mocks | grep -v proto | grep -v tests)
