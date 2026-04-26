#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

if [ "$#" -ne 2 ]; then
  echo "Usage: ./scripts/push-prod-images.sh <github-username-or-org> <tag>"
  echo "Example: ./scripts/push-prod-images.sh your-user 2026-04-26"
  exit 2
fi

IMAGE_NAMESPACE="ghcr.io/$1"
IMAGE_TAG="$2"

API_IMAGE="$IMAGE_NAMESPACE/books-on-hooks-api:$IMAGE_TAG"
CLIENT_IMAGE="$IMAGE_NAMESPACE/books-on-hooks-client:$IMAGE_TAG"

podman build -t "$API_IMAGE" "$REPO_ROOT/api"
podman build -t "$CLIENT_IMAGE" -f "$REPO_ROOT/frontend/Containerfile.nginx" "$REPO_ROOT/frontend"

podman push "$API_IMAGE"
podman push "$CLIENT_IMAGE"

echo "Pushed:"
echo "  $API_IMAGE"
echo "  $CLIENT_IMAGE"
