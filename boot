#!/bin/bash

# Docker image tarball file path
IMAGE_PATH=$1

if [ -z "$IMAGE_PATH" ]
then
    echo "Please provide docker image path as an argument."
    exit 1
fi

# Change directory to the path of Docker image
cd "$(dirname "$IMAGE_PATH")" || exit

# Load the docker image file into Docker's local image store
docker load -i "$(basename "$IMAGE_PATH")"

# Container Name
CONTAINER_NAME=$(basename "$IMAGE_PATH" .tar.gz)

# Check if the container is running
CONTAINER_ID=$(docker ps -a -q -f name="$CONTAINER_NAME")

# If container found
if [ -n "$CONTAINER_ID" ]
then
    echo "Stopping and removing container named '$CONTAINER_NAME'."

    docker stop "$CONTAINER_ID"
    docker rm "$CONTAINER_ID"
fi

# Build and run with docker compose
docker-compose up --build -d

# Remove docker image tarball file
rm -f "$IMAGE_PATH"

echo "Docker image file '$IMAGE_PATH' has been removed."