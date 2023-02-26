#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

print_usage() {
  printf "Usage: build [OPTIONS]

  Build dionego

  Options:

    -r  Build with race detector
"
}

race=''
while getopts 'r' flag; do
  case "${flag}" in
    r) race='-r' ;;
    *) print_usage
      exit 1 ;;
  esac
done

# Dionego root folder
DIONE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$DIONE_PATH"/scripts/constants.sh

# Download dependencies
echo "Downloading dependencies..."
go mod download

build_args="$race"

# Build dionego
"$DIONE_PATH"/scripts/build_dione.sh $build_args

# Exit build successfully if the DioneGo binary is created successfully
if [[ -f "$dionego_path" ]]; then
        echo "Build Successful"
        exit 0
else
        echo "Build failure" >&2
        exit 1
fi
