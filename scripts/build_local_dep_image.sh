#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

echo "Building docker image based off of most recent local commits of dionego and coreth"

DIONE_REMOTE="git@github.com:dioneprotocol/dionego.git"
CORETH_REMOTE="git@github.com:dioneprotocol/coreth.git"
DOCKERHUB_REPO="dioneprotocol/dionego"

DOCKER="${DOCKER:-docker}"
SCRIPT_DIRPATH=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

DIONE_LABS_RELATIVE_PATH="src/github.com/dioneprotocol"
EXISTING_GOPATH="$GOPATH"

export GOPATH="$SCRIPT_DIRPATH/.build_image_gopath"
WORKPREFIX="$GOPATH/src/github.com/dioneprotocol"

# Clone the remotes and checkout the desired branch/commits
DIONE_CLONE="$WORKPREFIX/dionego"
CORETH_CLONE="$WORKPREFIX/coreth"

# Replace the WORKPREFIX directory
rm -rf "$WORKPREFIX"
mkdir -p "$WORKPREFIX"


DIONE_COMMIT_HASH="$(git -C "$EXISTING_GOPATH/$DIONE_LABS_RELATIVE_PATH/dionego" rev-parse --short HEAD)"
CORETH_COMMIT_HASH="$(git -C "$EXISTING_GOPATH/$DIONE_LABS_RELATIVE_PATH/coreth" rev-parse --short HEAD)"

git config --global credential.helper cache

git clone "$DIONE_REMOTE" "$DIONE_CLONE"
git -C "$DIONE_CLONE" checkout "$DIONE_COMMIT_HASH"

git clone "$CORETH_REMOTE" "$CORETH_CLONE"
git -C "$CORETH_CLONE" checkout "$CORETH_COMMIT_HASH"

CONCATENATED_HASHES="$DIONE_COMMIT_HASH-$CORETH_COMMIT_HASH"

"$DOCKER" build -t "$DOCKERHUB_REPO:$CONCATENATED_HASHES" "$WORKPREFIX" -f "$SCRIPT_DIRPATH/local.Dockerfile"
