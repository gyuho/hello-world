#!/usr/bin/env bash
set -e

if [[ -z "${RELEASE_VERSION}" ]]; then
  RELEASE_VERSION=v0.0.1
fi

docker run \
  --name hello-world-${RELEASE_VERSION} \
  -p 32001:32001 \
  --rm \
  -e "PORT=32001" \
  docker.io/gyuho/hello-world:${RELEASE_VERSION} \
  /hello-world
