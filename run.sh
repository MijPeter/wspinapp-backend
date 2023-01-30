#!/bin/bash

if [ "$1" == "--test" ]; then
  docker compose -f docker-compose.test.yml up --build --remove-orphans -d
  go test tests/highlevel_test.go -v
  docker compose -f docker-compose.test.yml down
else
  docker compose up --build --remove-orphans -d
fi
