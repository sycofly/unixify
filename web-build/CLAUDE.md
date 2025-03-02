# Unixify - UNIX Account/Group Registry

Unixify is a Go application that serves as a registry for UNIX account UIDs and Group GIDs.

## Project Overview

The application provides a web interface for managing UNIX accounts and groups with the following features:

1. PostgreSQL database backend
2. Web interface with four sections: People, System, Database, and Service
3. Complete audit log system for all operations
4. Full RESTful API for all operations
5. JWT-based authentication with optional TOTP 2FA
6. Light/dark mode theme switching
7. Read-only guest access with registration for edit permissions

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
- User authentication with optional TOTP 2FA
- Theme switching (light/dark mode)
- Guest read-only access with registration for edit permissions

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

# View specific log entries
podman logs uno-861-acc-man_unixify_1 | grep "ERROR"

# Run a standalone frontend with Caddy (for testing UI changes)
podman build -t unixify-caddy -f Dockerfile.caddy .
podman run -d -p 8081:80 --name unixify-caddy unixify-caddy
```

## Container Notes

- The main application runs on port 8080
- The frontend-only container runs on port 8081
- Templates are stored in `/app/web/templates/`
- Static assets are in `/app/web/static/`
- Do NOT use volume mounts that override the container's web directory
- Use custom frontend container with Caddy for UI-only changes

## Tech Stack

- Go with Gin web framework
- PostgreSQL database
- HTML/CSS/JavaScript frontend
- RESTful API backend
- JWT-based authentication
- Google Authenticator TOTP support
- Theme switching with CSS variables
- Audit logging for all operations

## Authentication System

The application includes a comprehensive authentication system:
- JWT token-based authentication
- Password hashing with bcrypt
- Optional TOTP second factor with Google Authenticator
- Protected API routes with middleware
- User profiles and password management
- Self-registration with email verification and admin approval
- Guest read-only access for unauthenticated users

## Theming System

The application supports light and dark themes:
- CSS variables for theming
- Theme toggle button integrated in the navigation bar
- Theme preference stored in localStorage
- System preference detection via `prefers-color-scheme`
- Dark mode for all UI components including forms, tables, and alerts

## Access Control

The application implements a role-based access control system:
- Guests (unauthenticated users) have read-only access to view data
- Registration is required to request edit permissions
- New registrations require admin approval
- Authenticated users can perform edits based on their role
- UI dynamically adapts to show/hide edit controls based on permissions
- Clear visual indicators show current access mode (read-only vs. edit mode)
