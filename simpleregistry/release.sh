#!/bin/bash
set -euo pipefail

if [ ! -e "$( which github-release )" ]; then
  echo "You need github-release installed."
  echo "go get github.com/aktau/github-release"
  exit 2
fi

TAG="v0.0.1"
NAME="Actual Actuator"
DESCRIPTION="Prototypal release of the binder template registry"
USER="binder-project"
REPO="binder-registry"

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

docker push 

github-release upload \
  --user "$USER" \
  --repo "$REPO" \
  --tag "$TAG" \
  --name "darwin-amd64-simpleregistry" \
  --file bin/darwin-amd64-simpleregistry

