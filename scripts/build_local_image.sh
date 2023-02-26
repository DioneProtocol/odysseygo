#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Directory above this script
DIONE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Load the constants
source "$DIONE_PATH"/scripts/constants.sh

# WARNING: this will use the most recent commit even if there are un-committed changes present
full_commit_hash="$(git --git-dir="$DIONE_PATH/.git" rev-parse HEAD)"
commit_hash="${full_commit_hash::8}"

echo "Building Docker Image with tags: $dionego_dockerhub_repo:$commit_hash , $dionego_dockerhub_repo:$current_branch"
docker build -t "$dionego_dockerhub_repo:$commit_hash" \
        -t "$dionego_dockerhub_repo:$current_branch" "$DIONE_PATH" -f "$DIONE_PATH/Dockerfile"
