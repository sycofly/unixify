.PHONY: build run clean test lint migrate-up migrate-down docker-build docker-run docker-compose

# Build variables
BINARY_NAME=unixify
MIGRATE_BINARY=migrate

# Go build flags
LDFLAGS=-ldflags "-s -w"

# Build the application
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/unixify
	go build $(LDFLAGS) -o $(MIGRATE_BINARY) ./cmd/migrate

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f $(MIGRATE_BINARY)

# Run tests
test:
	go test -v ./...

# Run linting
lint:
	go vet ./...
	go fmt ./...

# Run database migrations up
migrate-up: build
	./$(MIGRATE_BINARY) -direction up

# Run database migrations down
migrate-down: build
	./$(MIGRATE_BINARY) -direction down

# Build Docker image
docker-build:
	podman build -t unixify -f Containerfile .

# Run Docker container
docker-run: docker-build
	podman run --name unixify -p 8080:8080 -e DB_HOST=host.containers.internal unixify

# Run with Docker Compose
docker-compose:
	podman-compose up -d

# Initialize project
init:
	cp .env.example .env
	mkdir -p web/static/js web/static/css web/static/img
	@echo "Project initialized. Edit .env file with your configuration."

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linters"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Run database migrations down"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make docker-compose - Run with Docker Compose"
	@echo "  make init           - Initialize project"
	@echo "  make help           - Show this help"