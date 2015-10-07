#!/bin/bash
set -euo pipefail

if [ ! -e "$( which github-release )" ]; then
  echo "You need github-release installed."
  echo "go get github.com/aktau/github-release"
  exit 2
fi

declare -xr USER="binder-project"
declare -xr REPO="binder-registry"

TAG=${1:-}
NAME=${2:-}
DESCRIPTION="Prototypal release of the binder template registry"

echo "Releasing '$TAG' - $NAME: $DESCRIPTION"

make
make docker-upload

github-release release \
  --user "$USER" \
  --repo "$REPO" \
  --tag "$TAG" \
  --pre-release \
  --name "$NAME" \
  --description "$DESCRIPTION"

github-release upload \
  --user "$USER" \
  --repo "$REPO" \
  --tag "$TAG" \
  --name "linux-amd64-simpleregistry" \
  --file bin/linux-amd64-simpleregistry

github-release upload \
  --user "$USER" \
  --repo "$REPO" \
  --tag "$TAG" \
  --name "darwin-amd64-simpleregistry" \
  --file bin/darwin-amd64-simpleregistry

