# Unixify - UNIX Account Management

![Unixify Logo](https://via.placeholder.com/150x150.png?text=Unixify)

Unixify is a comprehensive Go application that serves as a registry for UNIX account UIDs and Group GIDs, providing a modern web interface and API for managing these resources.

## 🌟 Features

- **PostgreSQL Database**: Robust data persistence with proper relations and constraints
- **Web Interface**: Clean, responsive UI built with Bootstrap 5 and Gin web framework
- **Multiple Account Types**: Dedicated sections for People, System, Database, and Service accounts
- **Group Management**: Create and manage groups with proper GID ranges
- **Membership Management**: Assign users to appropriate groups with validation
- **Search Functionality**: Find accounts and groups by UID, username, GID, or groupname
- **Audit Logging**: Track all changes with detailed audit entries
- **Soft Deletion**: Preserve data integrity with soft delete for accounts and groups
- **UID/GID Validation**: Enforce proper UID/GID ranges for different account types
- **RESTful API**: Comprehensive API for programmatic access
- **Containerization**: Docker and Podman support for easy deployment

## 📋 Prerequisites

- Go 1.22 or higher
- PostgreSQL 16 or higher
- Podman or Docker (optional, for containerized deployment)

## 🚀 Quick Start

### Using Docker/Podman

The fastest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/home/unixify.git
cd unixify

# Start the application and database
docker-compose up -d
```

Then access the web interface at http://localhost:8080

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/home/unixify.git
   cd unixify
   ```

2. Create and configure the environment:
   ```bash
   cp .env.example .env
   # Edit .env with your database settings
   ```

3. Set up the database:
   ```bash
   make migrate-up
   ```

4. Build and run:
   ```bash
   make build
   make run
   ```

5. Access the web interface at http://localhost:8080

## 📊 Database Schema

The application uses the following database schema:

- **accounts**: User accounts with UIDs, usernames, and types
- **groups**: Groups with GIDs, groupnames, and types
- **account_groups**: Many-to-many relationship between accounts and groups
- **audit_entries**: Comprehensive audit log for all system actions

## 🔍 UID and GID Ranges

The application enforces the following UID and GID ranges:

| Type     | UID/GID Range | Description                                |
|----------|---------------|--------------------------------------------|
| People   | 1000-60000    | Regular user accounts and groups           |
| System   | 1-999         | System accounts and groups                 |
| Service  | 60001-65535   | Service accounts and application services  |
| Database | 70000-79999   | Database users and related groups          |

## 📚 Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- [Installation Guide](docs/INSTALL.md): Detailed installation instructions
- [Usage Guide](docs/USAGE.md): How to use the application
- API Documentation: Available at `/api` endpoint when running the application

## 🔧 Development

### Makefile Commands

```bash
# Build the application
make build

# Run the application
make run

# Run database migrations
make migrate-up

# Run database migrations down
make migrate-down

# Run tests
make test

# Run linting
make lint

# Build Docker image
make docker-build

# Run with Docker
make docker-run

# Run with Docker Compose
make docker-compose
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request