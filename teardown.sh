#!/bin/bash
docker rm -f $(docker ps -aq)

docker rmi $(docker images | grep dev | awk '{print $3}')

docker volume prune

docker network prune
