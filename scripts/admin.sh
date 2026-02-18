#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

if [ "$#" -lt 1 ]; then
  echo "Usage: ./scripts/admin.sh <counts|reset-database> [args...]"
  exit 2
fi

podman compose -f "$REPO_ROOT/compose.yaml" exec -T api /usr/local/bin/admin "$@"
