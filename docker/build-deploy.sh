#!/bin/sh

set -eu

COMPANY_NAME="${COMPANY_NAME:-global-torque}"
SERVICE_NAME="${SERVICE_NAME:-go-common}"
REGISTRY="${REGISTRY:-cr.webdevelop.pro}"
CONTAINER_ENGINE="${CONTAINER_ENGINE:-docker}"
PLATFORMS="linux/amd64,linux/arm64"

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPOSITORY_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)

GIT_COMMIT=$(git -C "$REPOSITORY_ROOT" rev-parse --short=8 HEAD)
BUILD_DATE=$(date "+%Y%m%d")
REPOSITORY="$COMPANY_NAME/$SERVICE_NAME"
IMAGE="$REGISTRY/$REPOSITORY"

build_with_docker() {
  echo "Building and pushing $IMAGE for $PLATFORMS with Docker Buildx"
  "$CONTAINER_ENGINE" buildx build \
    --file "$SCRIPT_DIR/Dockerfile" \
    --platform "$PLATFORMS" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_DATE="$BUILD_DATE" \
    --build-arg SERVICE_NAME="$SERVICE_NAME" \
    --build-arg REPOSITORY="$REPOSITORY" \
    --build-arg VERSION="$BUILD_DATE:$GIT_COMMIT" \
    --tag "$IMAGE:$GIT_COMMIT" \
    --tag "$IMAGE:latest-dev" \
    --push \
    "$REPOSITORY_ROOT"
}

build_with_podman() {
  local_manifest="$SERVICE_NAME-$GIT_COMMIT-multiarch"

  echo "Building $local_manifest for $PLATFORMS with Podman"
  "$CONTAINER_ENGINE" manifest rm "$local_manifest" >/dev/null 2>&1 || true
  "$CONTAINER_ENGINE" build \
    --file "$SCRIPT_DIR/Dockerfile" \
    --platform "$PLATFORMS" \
    --manifest "$local_manifest" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_DATE="$BUILD_DATE" \
    --build-arg SERVICE_NAME="$SERVICE_NAME" \
    --build-arg REPOSITORY="$REPOSITORY" \
    --build-arg VERSION="$BUILD_DATE:$GIT_COMMIT" \
    "$REPOSITORY_ROOT"

  echo "Pushing $IMAGE:$GIT_COMMIT"
  "$CONTAINER_ENGINE" manifest push --all \
    "$local_manifest" "docker://$IMAGE:$GIT_COMMIT"

  echo "Pushing $IMAGE:latest-dev"
  "$CONTAINER_ENGINE" manifest push --all \
    "$local_manifest" "docker://$IMAGE:latest-dev"
}

BUILDX_VERSION=$("$CONTAINER_ENGINE" buildx version 2>/dev/null || true)
case "$BUILDX_VERSION" in
  *github.com/docker/buildx*)
    build_with_docker
    ;;
  *buildah* | *Buildah*)
    build_with_podman
    ;;
  *)
    echo "A Docker Buildx or Podman-compatible container engine is required." >&2
    exit 1
    ;;
esac
