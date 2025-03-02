# Unixify - UNIX Account/Group Registry

Unixify is a Go application that serves as a registry for UNIX account UIDs and Group GIDs.

## Project Overview

The application provides a web interface for managing UNIX accounts and groups with the following features:

1. PostgreSQL database backend
2. Web interface with four sections: People, System, Database, and Service
3. Complete audit log system for all operations
4. Full RESTful API for all operations

## Account/Group Types and UID/GID Ranges

| Type     | Account UID Range | Group GID Range |
|----------|-------------------|-----------------|
| People   | 1000-6000         | 1000-6000       |
| System   | 9000-9100         | 9000-9100       |
| Database | 7000-7999         | 7000-7999       |
| Service  | 8000-8999         | 8000-8999       |

## Key Operations

- Add/edit/delete accounts and groups
- Assign/remove users from groups
- View detailed audit logs of all system events
- Search by UID, GID, username, or groupname

## Development Commands

```bash
# Run the application locally
go run cmd/unixify/main.go

# Build the application
go build -o unixify cmd/unixify/main.go

# Run database migrations
go run cmd/migrate/main.go
```

## Deployment Commands

```bash
# Start the application with Podman
podman-compose up -d

# Stop the application
podman-compose down

# View application logs
podman logs uno-861-acc-man_unixify_1
```

## Container Notes

- The application runs on port 8080
- Templates are stored in `/app/web/templates/`
- Static assets are in `/app/web/static/`
- Do NOT use volume mounts that override the container's web directory

## Tech Stack

- Go with Gin web framework
- PostgreSQL database
- HTML/CSS/JavaScript frontend
- RESTful API backend
- Audit logging for all operations
