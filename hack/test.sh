#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=0

TARGETS=$(for d in "$@"; do echo ./$d/...; done)

echo "Building tests..."
go test -i -installsuffix "static" ${TARGETS}
echo "Running tests..."
go test -v -installsuffix "static" ${TARGETS}