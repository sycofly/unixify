#!/bin/bash

# Ensure network exists
podman network exists unixify-network || podman network create unixify-network

# Start database if not running
echo "Ensuring database is running..."
podman container exists unixify-db || podman run -d --name unixify-db \
  --network unixify-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=unixify \
  -p 5432:5432 \
  docker.io/library/postgres:14-alpine

# Start the application
echo "Starting the application..."
podman container exists unixify-app && podman rm -f unixify-app
podman run -d --name unixify-app \
  --network unixify-network \
  -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e DB_HOST=unixify-db \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=unixify \
  -e DB_SSLMODE=disable \
  -v $(pwd)/server:/app/server:Z \
  -v $(pwd)/web:/web:Z \
  -w /app \
  --entrypoint "/app/server" \
  docker.io/library/alpine:3.18

echo "Application started. Access it at http://localhost:8080"
echo "To stop the application, run: podman rm -f unixify-app"