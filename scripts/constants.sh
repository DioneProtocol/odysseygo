#!/usr/bin/env bash
#
# Use lower_case variables in the scripts and UPPER_CASE variables for override
# Use the constants.sh for env overrides

ODYSSEY_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd ) # Directory above this script

# Where OdysseyGo binary goes
odysseygo_path="$ODYSSEY_PATH/build/odysseygo"

# Avalabs docker hub
# odyplatform/odysseygo - defaults to local as to avoid unintentional pushes
# You should probably set it - export DOCKER_REPO='odyplatform/odysseygo'
odysseygo_dockerhub_repo=${DOCKER_REPO:-"odysseygo"}

# Current branch
# TODO: fix "fatal: No names found, cannot describe anything" in github CI
current_branch=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match || true)

git_commit=${ODYSSEYGO_COMMIT:-$( git rev-list -1 HEAD )}

# Static compilation
static_ld_flags=''
if [ "${STATIC_COMPILATION:-}" = 1 ]
then
    export CC=musl-gcc
    which $CC > /dev/null || ( echo $CC must be available for static compilation && exit 1 )
    static_ld_flags=' -extldflags "-static" -linkmode external '
fi

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"
# While CGO_ENABLED doesn't need to be explicitly set, it produces a much more
# clear error due to the default value change in go1.20.
export CGO_ENABLED=1
