#!/bin/sh

set -e

TIMEOUT=${TIMEOUT:-30}

docker compose -f $PWD/tests/docker-compose.yml up -d

trap "docker compose -f $PWD/tests/docker-compose.yml down" EXIT

echo "Waiting for Nebraska to be up"
i=0
while true; do
  if curl -sSf -o /dev/null http://localhost:8000/api/apps 2>/dev/null; then
    printf "\n"
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
TF_ACC=true NEBRASKA_ENDPOINT=http://localhost:8000 NEBRASKA_USERNAME=user NEBRASKA_PASSWORD=pass go test ./nebraska -v -timeout 120m
TF_ACC=true NEBRASKA_ENDPOINT=http://localhost:8000 NEBRASKA_BEARER_TOKEN=token go test ./nebraska -v -timeout 120m
