#!/usr/bin/env bash
set -e

ENDPOINT=http://localhost:32001
echo ENDPOINT: ${ENDPOINT}

curl -L ${ENDPOINT}/hello-world
curl -L ${ENDPOINT}/hello-world-readiness
curl -L ${ENDPOINT}/hello-world-liveness
curl -L ${ENDPOINT}/hello-world-status
