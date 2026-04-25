#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/compose.yaml"
HEALTH_URL="http://localhost:8082/api/healthz"

podman compose -f "$COMPOSE_FILE" up -d --build test-db api client-nginx-test

DB_READY=0
for _ in $(seq 1 60); do
  if podman compose -f "$COMPOSE_FILE" exec -T test-db pg_isready -U postgres -d books >/dev/null 2>&1; then
    DB_READY=1
    break
  fi
  sleep 1
done

if [ "$DB_READY" -ne 1 ]; then
  echo "test-db did not become ready in time" >&2
  podman compose -f "$COMPOSE_FILE" down
  exit 1
fi

HTTP_READY=0
for _ in $(seq 1 60); do
  if curl -fsS "$HEALTH_URL" >/dev/null 2>&1; then
    HTTP_READY=1
    break
  fi
  sleep 1
done

if [ "$HTTP_READY" -ne 1 ]; then
  echo "api/nginx health endpoint did not become ready in time: $HEALTH_URL" >&2
  podman compose -f "$COMPOSE_FILE" down
  exit 1
fi

podman compose -f "$COMPOSE_FILE" build playwright

set +e
podman compose -f "$COMPOSE_FILE" run --rm playwright npx playwright test /work/app-smoke.spec.js --reporter=line
SMOKE_EXIT_CODE=$?
if [ "$SMOKE_EXIT_CODE" -eq 0 ]; then
  podman compose -f "$COMPOSE_FILE" run --rm playwright
  EXIT_CODE=$?
else
  EXIT_CODE=$SMOKE_EXIT_CODE
fi
set -e

podman compose -f "$COMPOSE_FILE" down

exit "$EXIT_CODE"
