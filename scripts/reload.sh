#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

touch "$REPO_ROOT/frontend/src/main.jsx"
touch "$REPO_ROOT/api/cmd/web/main.go"

echo "Triggered reload markers in frontend and api."
