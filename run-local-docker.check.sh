#!/usr/bin/env bash
set -e

ENDPOINT=http://localhost:32001
echo ENDPOINT: ${ENDPOINT}

curl -L ${ENDPOINT}/hello-world
curl -L ${ENDPOINT}/readiness
curl -L ${ENDPOINT}/liveness
curl -L ${ENDPOINT}/status
