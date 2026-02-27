#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/compose.yaml"

set +e
podman compose -f "$COMPOSE_FILE" up --build --abort-on-container-exit --exit-code-from playwright playwright
EXIT_CODE=$?
set -e

podman compose -f "$COMPOSE_FILE" down

exit "$EXIT_CODE"
