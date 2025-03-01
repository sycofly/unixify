#!/bin/bash

# Build the application using podman
podman build -f Dockerfile.build -t unixify-build .

# Extract the compiled binary from the container
podman create --name unixify-extract unixify-build
podman cp unixify-extract:/app/server ./server
podman rm unixify-extract

echo "Build completed. The binary is available at: ./server"