#!/bin/bash

docker-compose rm
GOOS=linux GOARCH=amd64 make build
REMOVE_CONTAINERS=on DOCKER_IMAGE=golang-example-app TAG=development make docker-image
docker-compose up
