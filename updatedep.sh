#!/usr/bin/env bash
set -e

if [[ $(go version) != "go version go1.11"* ]]; then
  echo "expect Go 1.11+, got:" $(go version)
  exit 255
fi

<<COMMENT
GO111MODULE=on go mod init
COMMENT

GO111MODULE=on go mod tidy -v
GO111MODULE=on go mod vendor -v
