# Unixify Installation Guide

This document provides instructions on how to install and configure the Unixify UNIX account management system.

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 16 or higher
- Git

## Installation Methods

There are several ways to install and run Unixify:

1. [Local Development Setup](#local-development-setup)
2. [Using Docker](#using-docker)
3. [Using Podman](#using-podman)
4. [Production Deployment](#production-deployment)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/home/unixify.git
cd unixify
```

### 2. Set Up Environment Variables

Copy the example environment file and modify as needed:

```bash
cp .env.example .env
```

Edit the `.env` file to configure your database connection and other settings.

### 3. Create the Database

Create a PostgreSQL database for Unixify:

```bash
createdb unixify
```

Or using SQL:

```sql
CREATE DATABASE unixify;
```

### 4. Run Database Migrations

```bash
make migrate-up
```

### 5. Build and Run

```bash
make build
make run
```

Or directly:

```bash
go build -o unixify ./cmd/unixify
./unixify
```

### 6. Access the Application

Open your browser and navigate to: http://localhost:8080

## Using Docker

### 1. Clone the Repository

```bash
git clone https://github.com/home/unixify.git
cd unixify
```

### 2. Build and Run using Docker Compose

```bash
docker-compose up -d
```

This will:
- Start a PostgreSQL container
- Build and start the Unixify container
- Set up networking between containers
- Set environment variables

### 3. Access the Application

Open your browser and navigate to: http://localhost:8080

## Using Podman

### 1. Clone the Repository

```bash
git clone https://github.com/home/unixify.git
cd unixify
```

### 2. Build and Run using Podman

```bash
# Build the image
podman build -t unixify -f Containerfile .

# Run PostgreSQL container
podman run -d --name postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_DB=unixify \
    -p 5432:5432 \
    docker.io/library/postgres:16-alpine

# Run Unixify container
podman run -d --name unixify \
    -e DB_HOST=host.containers.internal \
    -e DB_PORT=5432 \
    -e DB_USER=postgres \
    -e DB_PASSWORD=postgres \
    -e DB_NAME=unixify \
    -p 8080:8080 \
    unixify
```

### 3. Or using Podman Compose

```bash
podman-compose up -d
```

### 4. Access the Application

Open your browser and navigate to: http://localhost:8080

## Production Deployment

For production deployment, consider the following security and stability enhancements:

### 1. Environment Configuration

Edit your `.env` file or set environment variables:

```
# Set to production mode
GIN_MODE=release

# Use a strong secret for JWT
JWT_SECRET=your_very_strong_secret_key

# Use SSL for database connection
DB_SSLMODE=require
```

### 2. Database Security

- Create a dedicated database user with appropriate permissions
- Enable SSL connections to the database
- Set up database backups

### 3. Reverse Proxy

Set up a reverse proxy like Nginx to:
- Handle SSL termination
- Provide HTTPS support
- Handle request buffering
- Optional: Basic authentication

Example Nginx configuration:

```nginx
server {
    listen 80;
    server_name unixify.example.com;
    
    # Redirect to HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name unixify.example.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. Process Management

For robust process management, consider:

- Systemd service
- Docker with auto-restart
- Container orchestration (Kubernetes, etc.)

Example systemd service file (`/etc/systemd/system/unixify.service`):

```ini
[Unit]
Description=Unixify UNIX Account Management
After=network.target postgresql.service

[Service]
User=unixify
WorkingDirectory=/opt/unixify
ExecStart=/opt/unixify/unixify
Restart=on-failure
RestartSec=5
Environment=GIN_MODE=release
EnvironmentFile=/opt/unixify/.env

[Install]
WantedBy=multi-user.target
```

Activate and start the service:

```bash
systemctl daemon-reload
systemctl enable unixify
systemctl start unixify
```

## Troubleshooting

### Database Connection Issues

- Verify PostgreSQL is running: `pg_isready -h <host> -p <port>`
- Check connectivity: `psql -h <host> -p <port> -U <user> -d <dbname> -c "SELECT 1;"`
- Verify environment variables match your PostgreSQL configuration

### Application Not Starting

- Check logs: `journalctl -u unixify` (if using systemd)
- Verify permissions on the application directory
- Check for port conflicts: `netstat -tuln | grep 8080`

### Database Migration Failures

- Check database schema manually: `psql -U <user> -d <dbname> -c "\\dt"`
- Try running migrations manually: `./migrate -direction up`
- Check migration logs for specific errors

## Next Steps

After successful installation, refer to the [Usage Guide](USAGE.md) for information on how to use the application.