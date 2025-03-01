# Unixify

A web application for managing user accounts and groups, built with Go, HTML, and HTMX.

## Features

- User account management (add, edit, search, soft delete)
- Group management (add, edit, search, soft delete)
- Audit tracking (created_by, updated_by, deleted_by)
- PostgreSQL database with indexed fields and proper schema
- Modern UI with HTMX for interactive features
- Containerized with Podman for easy deployment

## Technology Stack

- **Backend**: Go (Golang) with Chi router
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, HTMX, Alpine.js
- **Containerization**: Podman and Podman Compose
- **Platform**: Red Hat family Linux

## Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Podman and Podman Compose

### Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/home/unixify.git
   cd unixify
   ```

2. Build the application:
   ```
   ./build.sh
   ```

3. Set up the database and run migrations:
   ```
   ./migrate.sh
   ```

4. Start the application:
   ```
   ./run.sh
   ```

5. Access the application at http://localhost:8080

6. To stop the application:
   ```
   podman-compose down
   ```

### Project Structure

- `/cmd/server`: Main application entry point
- `/internal`: Internal packages
  - `/config`: Application configuration
  - `/database`: Database connection and migrations
  - `/handlers`: HTTP handlers
  - `/middleware`: HTTP middleware
  - `/models`: Data models
  - `/services`: Business logic services
  - `/utils`: Utility functions
- `/web`: Web assets
  - `/static`: Static files (JS, CSS, images)
  - `/templates`: HTML templates

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.