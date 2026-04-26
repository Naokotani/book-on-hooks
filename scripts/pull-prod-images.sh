#!/usr/bin/env sh
set -eu

if [ "$#" -ne 2 ]; then
  echo "Usage: ./scripts/pull-prod-images.sh <github-username-or-org> <tag>"
  echo "Example: ./scripts/pull-prod-images.sh your-user 2026-04-26"
  exit 2
fi

IMAGE_NAMESPACE="ghcr.io/$1"
IMAGE_TAG="$2"

API_IMAGE="$IMAGE_NAMESPACE/books-on-hooks-api:$IMAGE_TAG"
CLIENT_IMAGE="$IMAGE_NAMESPACE/books-on-hooks-client:$IMAGE_TAG"

docker pull "$API_IMAGE"
docker pull "$CLIENT_IMAGE"

echo "Pulled:"
echo "  $API_IMAGE"
echo "  $CLIENT_IMAGE"
