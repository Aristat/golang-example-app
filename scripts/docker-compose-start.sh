#!/bin/bash

docker-compose rm
REMOVE_CONTAINERS=on DOCKER_IMAGE=golang-example-app TAG=development make docker-image
docker-compose up
