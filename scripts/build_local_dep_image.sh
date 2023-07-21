#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

echo "Building docker image based off of most recent local commits of odysseygo and coreth"

ODYSSEY_REMOTE="git@github.com:DioneProtocol/odysseygo.git"
CORETH_REMOTE="git@github.com:DioneProtocol/coreth.git"
DOCKERHUB_REPO="dioneprotocol/odysseygo"

DOCKER="${DOCKER:-docker}"
SCRIPT_DIRPATH=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

DIONE_PROTOCOL_RELATIVE_PATH="src/github.com/DioneProtocol"
EXISTING_GOPATH="$GOPATH"

export GOPATH="$SCRIPT_DIRPATH/.build_image_gopath"
WORKPREFIX="$GOPATH/src/github.com/DioneProtocol"

# Clone the remotes and checkout the desired branch/commits
ODYSSEY_CLONE="$WORKPREFIX/odysseygo"
CORETH_CLONE="$WORKPREFIX/coreth"

# Replace the WORKPREFIX directory
rm -rf "$WORKPREFIX"
mkdir -p "$WORKPREFIX"


ODYSSEY_COMMIT_HASH="$(git -C "$EXISTING_GOPATH/$DIONE_PROTOCOL_RELATIVE_PATH/odysseygo" rev-parse --short HEAD)"
CORETH_COMMIT_HASH="$(git -C "$EXISTING_GOPATH/$DIONE_PROTOCOL_RELATIVE_PATH/coreth" rev-parse --short HEAD)"

git config --global credential.helper cache

git clone "$ODYSSEY_REMOTE" "$ODYSSEY_CLONE"
git -C "$ODYSSEY_CLONE" checkout "$ODYSSEY_COMMIT_HASH"

git clone "$CORETH_REMOTE" "$CORETH_CLONE"
git -C "$CORETH_CLONE" checkout "$CORETH_COMMIT_HASH"

CONCATENATED_HASHES="$ODYSSEY_COMMIT_HASH-$CORETH_COMMIT_HASH"

"$DOCKER" build -t "$DOCKERHUB_REPO:$CONCATENATED_HASHES" "$WORKPREFIX" -f "$SCRIPT_DIRPATH/local.Dockerfile"
