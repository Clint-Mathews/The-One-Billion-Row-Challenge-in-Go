#!/bin/bash

# Build the Docker image
docker build -t create_data . > /dev/null

echo "Docker image created successfully."

# Run the Docker container with a name
docker run -d --name create_data create_data > /dev/null
echo "Docker container running!"

# Remove the docker image
docker rmi -f create_data > /dev/null
echo "Image cleaned up successfully"

# Copy data from the container to the current directory
docker cp create_data:/1brc/data/measurements.txt . 

echo "Data copied successfully."

# Stop and remove the container
docker stop create_data > /dev/null
docker rm create_data > /dev/null

echo "Container cleaned up successfully"