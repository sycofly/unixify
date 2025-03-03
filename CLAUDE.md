# Unixify - UNIX Account/Group Registry

Unixify is a Go application that serves as a registry for UNIX account UIDs and Group GIDs.

## Project Overview

The application provides a web interface for managing UNIX accounts and groups with the following features:

1. PostgreSQL database backend
2. Web interface with four sections: People, System, Database, and Service
3. Complete audit log system for all operations
4. Full RESTful API for all operations
5. JWT-based authentication with optional TOTP 2FA
6. Light/dark mode theme switching with auto-detection
7. Read-only guest mode with visual indicators
8. Gradient text and consistent button styling

## Account/Group Types and UID/GID Ranges

| Type     | Account UID Range | Group GID Range |
|----------|-------------------|-----------------|
| People   | 5000-6000         | 1000-3000       |
| System   | 1000-2000         | 3000-5000       |
| Database | 2000-7999         | 2000-7500       |
| Service  | 8000-8999         | 4000-5000       |

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
- Automatic guest mode with clear visual indicators
- Proper separation between regular users and guest accounts

## Theming System

The application supports light and dark themes:
- CSS variables for comprehensive theme support
- Theme toggle button integrated in the navigation bar
- Theme preference stored in localStorage for persistence
- System preference detection via `prefers-color-scheme`
- Dark mode for all UI components including forms, tables, and alerts
- Light grey text in tables for better dark mode readability
- Gradient text effects for headings and descriptions
- Consistent color palette for buttons and interactive elements
- Custom colored badges with theme-appropriate styling

## Access Control

The application implements a role-based access control system:
- Guests have automatic read-only access to view data
- "Guest Account (Read-Only)" indicator clearly shows status
- Registration is required to request edit permissions
- New registrations require admin approval
- Authenticated users can perform edits based on their role
- UI dynamically adapts to show/hide edit controls based on permissions
- Clear visual indicators show current access mode
- Proper user dropdown menu visibility control for guest vs regular users

## UI Enhancements

The application features a modern and user-friendly interface:
- Clean, responsive layout with Bootstrap 5
- Soft purple and green accent colors for key UI elements
- Blue-to-purple gradient text for headings and descriptions
- Consistent button styling across the application
- Card-based UI with subtle shadows and animations
- Interactive hover effects for better user feedback
- Custom document viewer for application documentation
- UID/GID Range Reference card with optimized dark mode display
