#!/bin/bash

# Create network if it doesn't exist
echo "Ensuring network exists..."
podman network exists unixify-network || podman network create unixify-network

# Start database container if not running
echo "Ensuring database is running..."
podman container exists unixify-db || podman run -d --name unixify-db \
  --network unixify-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=unixify \
  -p 5432:5432 \
  docker.io/library/postgres:14-alpine

# Wait for the database to be ready
echo "Waiting for database to be ready..."
sleep 10

# Create a temporary file with the migration SQL
echo "Preparing migration SQL..."
MIGRATION_FILE=$(pwd)/internal/database/migrations/001_initial_schema.sql
TEMP_FILE=$(mktemp)
cat $MIGRATION_FILE > $TEMP_FILE

# Run migrations using a temporary container
echo "Running database migrations..."
podman run --rm --network unixify-network \
  -v $TEMP_FILE:/migrations.sql:Z \
  -e PGPASSWORD=postgres \
  docker.io/library/postgres:14-alpine \
  psql -h unixify-db -U postgres -d unixify -f /migrations.sql

# Clean up
rm $TEMP_FILE

echo "Migration completed."