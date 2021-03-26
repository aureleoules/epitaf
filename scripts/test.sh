#!/bin/sh
docker-compose -f ./docker/test/docker-compose.yml up --build --abort-on-container-exit
docker-compose -f ./docker/test/docker-compose.yml down