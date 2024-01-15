#!/usr/bin/env bash

set -euo pipefail

# Odysseygo root folder
ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$ODYSSEY_PATH"/scripts/constants.sh

echo "Building testnetctl..."
go build -ldflags\
   "-X github.com/DioneProtocol/odysseygo/version.GitCommit=$git_commit $static_ld_flags"\
   -o "$ODYSSEY_PATH/build/testnetctl"\
   "$ODYSSEY_PATH/tests/fixture/testnet/cmd/"*.go
