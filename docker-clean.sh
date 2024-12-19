#!/bin/bash
# Stop all running containers
docker stop $(docker ps -q)

# Remove all stopped containers
docker rm $(docker ps -a -q)

# Remove all unused images, containers, volumes, and networks
docker system prune -a --volumes --force

# Remove Docker build cache
docker builder prune --force

echo "Docker cleanup completed."
