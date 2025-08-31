#!/bin/bash

COMPANY_NAME=webdevelop-pro
SERVICE_NAME=go-common

case $1 in

run)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  BUILD_DATE=$(date "+%Y%m%d")
  build && ./app
  ;;


*)
  # This is the default case for building and pushing multi-arch images.
  BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  echo "Building and pushing multi-arch image for branch: $BRANCH_NAME, commit: $GIT_COMMIT"

  # Define full image names for each architecture
  AMD64_IMG="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-amd64"
  ARM64_IMG="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-arm64"

  # Define manifest list names
  MANIFEST_COMMIT="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT"
  MANIFEST_LATEST="cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev"

  # Step 1: Build each architecture-specific image.
  echo "Building linux/amd64 image..."
  podman build --platform linux/amd64 -t $AMD64_IMG .
  echo "Building linux/arm64 image..."
  podman build --platform linux/arm64 -t $ARM64_IMG .

  # Step 2: Push the individual architecture-specific images to the registry.
  echo "Pushing individual architecture images..."
  podman push $AMD64_IMG
  podman push $ARM64_IMG

  # Step 3: Create manifest lists and add the images.
  echo "Creating manifest for commit tag: $MANIFEST_COMMIT"
  podman manifest create $MANIFEST_COMMIT
  podman manifest add $MANIFEST_COMMIT $AMD64_IMG
  podman manifest add $MANIFEST_COMMIT $ARM64_IMG

  echo "Creating manifest for latest-dev tag: $MANIFEST_LATEST"
  podman manifest create $MANIFEST_LATEST
  podman manifest add $MANIFEST_LATEST $AMD64_IMG
  podman manifest add $MANIFEST_LATEST $ARM64_IMG

  # Step 4: Push the manifest lists to the registry.
  echo "Pushing manifest lists..."
  podman manifest push --all $MANIFEST_COMMIT
  podman manifest push --all $MANIFEST_LATEST

  # docker buildx build --platform linux/amd64,linux/arm64 -t docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT -t docker.io/webdeveloppro/$SERVICE_NAME:latest-dev -t cr.webdevelop.biz/$COMPANY_NAME/$SERVICE_NAME:latest-dev --platform=linux/amd64 .
  # docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT
  # docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  # docker push docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT
  # docker push webdeveloppro/$SERVICE_NAME:latest-dev

  echo "Push complete for tags $GIT_COMMIT and latest-dev."
  ;;
esac

esac

