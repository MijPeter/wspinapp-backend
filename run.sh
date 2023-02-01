#!/bin/bash


# TODO make this script to work with man
# TODO create a flag --generate-goldens that makes it so that instead of failing when goldens don't match -> generate them again

if [ "$1" == "--test" ]; then
  docker compose -f docker-compose.test.yml up --build --remove-orphans --abort-on-container-exit
elif [ "$1" == "--fast" ]; then
  docker compose -f docker-compose.test.yml up -d test_db --build --remove-orphans
  go test tests/highlevel_test.go -v
  docker compose -f docker-compose.test.yml down
else
  docker compose up --build --remove-orphans -d
fi
