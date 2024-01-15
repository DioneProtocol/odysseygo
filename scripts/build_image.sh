#!/usr/bin/env bash

set -euo pipefail

# Directory above this script
ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Load the constants
source "$ODYSSEY_PATH"/scripts/constants.sh

# WARNING: this will use the most recent commit even if there are un-committed changes present
full_commit_hash="$(git --git-dir="$ODYSSEY_PATH/.git" rev-parse HEAD)"
commit_hash="${full_commit_hash::8}"

echo "Building Docker Image with tags: $odysseygo_dockerhub_repo:$commit_hash , $odysseygo_dockerhub_repo:$current_branch"
docker build -t "$odysseygo_dockerhub_repo:$commit_hash" \
        -t "$odysseygo_dockerhub_repo:$current_branch" "$ODYSSEY_PATH" -f "$ODYSSEY_PATH/Dockerfile"
