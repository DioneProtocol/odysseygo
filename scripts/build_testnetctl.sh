#!/usr/bin/env bash

set -euo pipefail

# Odysseygo root folder
AVALANCHE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$AVALANCHE_PATH"/scripts/constants.sh

echo "Building testnetctl..."
go build -ldflags\
   "-X github.com/DioneProtocol/odysseygo/version.GitCommit=$git_commit $static_ld_flags"\
   -o "$AVALANCHE_PATH/build/testnetctl"\
   "$AVALANCHE_PATH/tests/fixture/testnet/cmd/"*.go
