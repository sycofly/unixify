# Unixify - UNIX Account Management

Unixify is a Go application that serves as a registry for UNIX account UIDs and Group GIDs, providing a web interface for managing these resources.

## Features

- PostgreSQL database backend for data persistence
- Web interface for managing UIDs and GIDs using Gin web framework
- Four sections: People, System, Database, and Service
- Support for account and group operations:
  - Add users and groups
  - Assign users to groups
  - Remove users from groups
  - List users and groups
  - Search functionality
  - Audit logging

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 16 or higher
- Podman or Docker (for containerized deployment)

## Getting Started

### Local Development

1. Clone this repository:
   ```bash
   git clone https://github.com/home/unixify.git
   cd unixify
   ```

2. Set up environment variables or create a `.env` file:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=unixify
   DB_SSLMODE=disable
   SERVER_PORT=8080
   GIN_MODE=debug
   JWT_SECRET=your_secret_key
   ```

3. Build and run the application:
   ```bash
   go build -o unixify ./cmd/unixify
   ./unixify
   ```

4. Access the web interface at http://localhost:8080

### Using Podman/Docker

1. Build and run using docker-compose:
   ```bash
   podman-compose up -d
   # or with Docker
   docker-compose up -d
   ```

2. Access the web interface at http://localhost:8080

## Database Schema

The application uses the following database schema:

- `accounts`: Stores user accounts with UIDs
- `groups`: Stores groups with GIDs
- `account_groups`: Many-to-many relationship between accounts and groups
- `audit_entries`: Audit log for all actions

## API Endpoints

### Account Endpoints

- `GET /api/accounts` - Get all accounts
- `GET /api/accounts/:id` - Get account by ID
- `POST /api/accounts` - Create a new account
- `PUT /api/accounts/:id` - Update an account
- `DELETE /api/accounts/:id` - Delete an account
- `GET /api/accounts/uid/:uid` - Get account by UID
- `GET /api/accounts/username/:username` - Get account by username
- `GET /api/accounts/:id/groups` - Get groups for an account

### Group Endpoints

- `GET /api/groups` - Get all groups
- `GET /api/groups/:id` - Get group by ID
- `POST /api/groups` - Create a new group
- `PUT /api/groups/:id` - Update a group
- `DELETE /api/groups/:id` - Delete a group
- `GET /api/groups/gid/:gid` - Get group by GID
- `GET /api/groups/groupname/:groupname` - Get group by groupname
- `GET /api/groups/:id/accounts` - Get accounts in a group

### Membership Endpoints

- `POST /api/memberships` - Assign account to group
- `DELETE /api/memberships` - Remove account from group

### Search Endpoints

- `GET /api/search/accounts?q=query` - Search accounts
- `GET /api/search/groups?q=query` - Search groups

### Audit Endpoints

- `GET /api/audit` - Get audit entries
- `GET /api/audit/:id` - Get specific audit entry

## UID and GID Ranges

The application enforces the following UID and GID ranges:

- People: 1000-60000
- System: 1-999
- Service: 60001-65535
- Database: 70000-79999

## License

This project is licensed under the MIT License - see the LICENSE file for details.