# Build stage
FROM docker.io/library/golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM docker.io/library/alpine:3.18

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/server /app/server

# Copy assets and templates
COPY --from=builder /app/web /app/web

# Create a non-root user and set ownership
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser

# Set environment variables
ENV SERVER_PORT=8080

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/server"]