#!/usr/bin/env bash
set -e

if [[ -z "${GIT_COMMIT}" ]]; then
  GIT_COMMIT=$(git rev-parse --short HEAD || echo "GitNotFound")
fi

if [[ -z "${RELEASE_VERSION}" ]]; then
  RELEASE_VERSION=v0.0.1
fi

if [[ -z "${BUILD_TIME}" ]]; then
  BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
fi

CGO_ENABLED=0 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -v \
  -ldflags "-s -w \
  -X github.com/gyuho/hello-world/version.GitCommit=${GIT_COMMIT} \
  -X github.com/gyuho/hello-world/version.ReleaseVersion=${RELEASE_VERSION} \
  -X github.com/gyuho/hello-world/version.BuildTime=${BUILD_TIME}" \
  -o ./hello-world
