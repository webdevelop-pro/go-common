#!/bin/bash

set -e

if [ ! -d ".git" ]; then
  echo "Not found .git directory. Please run the script in the root of the repository."
  exit 1
fi

MODULE_DIRS=$(find . -maxdepth 2 -name "go.mod" -exec dirname {} \;)

TAG_VERSION="v1.0.9"

for DIR in $MODULE_DIRS; do
  # skip dir if it doesn't have go.mod file
  if [ ! -f "$DIR/go.mod" ]; then
    continue
  fi

  MODULE_NAME=$(basename "$DIR")

  TAG_NAME="$MODULE_NAME/$TAG_VERSION"

  echo "Creating tag: $TAG_NAME for module $MODULE_NAME"

   git tag -m $TAG_NAME"" "$TAG_NAME"  || true

done

echo "Pushing tags to the remote repository..."
git push --tags

echo "All tags have been created and pushed to the remote repository."
