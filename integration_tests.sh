#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/api"
GOCACHE=/tmp/go-build-cache go test -tags=integration ./integration -v "$@"
