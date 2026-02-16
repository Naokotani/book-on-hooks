#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/api"

go test ./... "$@"
go test -tags=integration ./integration -v "$@"
