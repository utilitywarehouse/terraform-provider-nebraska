#!/bin/sh

set -e

TIMEOUT=${TIMEOUT:-30}

command -v docker-compose >/dev/null \
  || { echo "ERROR: tests require docker-compose"; exit 1; }

docker-compose -f "$(pwd)"/tests/docker-compose.yml up -d

trap "docker-compose -f $(pwd)/tests/docker-compose.yml down"  EXIT

echo "Waiting for Nebraska to be up"
i=0
while true; do
  if curl -sSf -o /dev/null http://localhost:8000/api/apps 2>/dev/null; then
    break
  fi
  i=$((i+1))
  if [ $i -ge $TIMEOUT ]; then
    echo "Couldn't reach Nebraska within the timeout: ${TIMEOUT}s" >&2
    exit 1
  fi
  printf "."
  sleep 1
done

TF_ACC=true NEBRASKA_ENDPOINT=http://localhost:8000 go test ./... -v -timeout 120m
