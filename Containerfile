# Build stage
FROM docker.io/library/golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o unixify ./cmd/unixify

# Final stage
FROM docker.io/library/alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/unixify /app/unixify

# Copy web templates and static files
COPY ./web-build /app/web

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Command to run
CMD ["/app/unixify"]