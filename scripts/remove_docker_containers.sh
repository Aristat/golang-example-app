#!/bin/bash

echo "Start removing containers"

if [ "$REMOVE_CONTAINERS" == ON ]; then
  docker_containers=$(docker ps -a -q --filter "name=$DOCKER_IMAGE")

  if [ -z "$docker_containers" ]; then
    echo "No containers for remove"
  else
    docker rm $docker_containers
  fi
else
    echo "Containers removal disabled"
fi
