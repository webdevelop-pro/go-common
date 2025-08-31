#!/bin/bash
set -euo pipefail

COMPANY_NAME=webdevelop-pro
SERVICE_NAME=go-common

case ${1:-} in
run)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  BUILD_DATE=$(date "+%Y%m%d")
  build && ./app
  ;;
*)
  BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  echo "Building and pushing multi-arch image for branch: $BRANCH_NAME, commit: $GIT_COMMIT"

  AMD64_IMG="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:${GIT_COMMIT}-amd64"
  ARM64_IMG="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:${GIT_COMMIT}-arm64"

  # Final remote tags
  REMOTE_COMMIT="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:${GIT_COMMIT}"
  REMOTE_LATEST="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev"

  # Use a local-only manifest name to avoid clashing with existing image tags
  LOCAL_MANIFEST_COMMIT="${SERVICE_NAME}-manifest-${GIT_COMMIT}"
  LOCAL_MANIFEST_LATEST="${SERVICE_NAME}-manifest-latest-dev"

  echo "Building linux/amd64 image..."
  podman build --platform linux/amd64 -t "$AMD64_IMG" .
  echo "Building linux/arm64 image..."
  podman build --platform linux/arm64 -t "$ARM64_IMG" .

  echo "Pushing individual architecture images..."
  podman push "$AMD64_IMG"
  podman push "$ARM64_IMG"

  echo "Creating local manifest lists..."
  podman manifest create "$LOCAL_MANIFEST_COMMIT"
  podman manifest create "$LOCAL_MANIFEST_LATEST"

  # Add the remote images by reference, avoids needing them locally again
  podman manifest add "$LOCAL_MANIFEST_COMMIT" "docker://$AMD64_IMG"
  podman manifest add "$LOCAL_MANIFEST_COMMIT" "docker://$ARM64_IMG"

  podman manifest add "$LOCAL_MANIFEST_LATEST" "docker://$AMD64_IMG"
  podman manifest add "$LOCAL_MANIFEST_LATEST" "docker://$ARM64_IMG"

  echo "Pushing manifest lists to the registry…"
  podman manifest push --all "$LOCAL_MANIFEST_COMMIT" "docker://$REMOTE_COMMIT"
  podman manifest push --all "$LOCAL_MANIFEST_LATEST" "docker://$REMOTE_LATEST"

  echo "Push complete for tags ${GIT_COMMIT} and latest-dev."
  ;;
esac

