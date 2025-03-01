# CLAUDE.md - Development Guidelines

## Build & Test Commands
- Build: `go build -o bin/server ./cmd/server`
- Run: `go run ./cmd/server`
- Test all: `go test ./...`
- Test specific: `go test ./internal/package-name`
- Test with coverage: `go test -cover ./...`
- Lint: `golangci-lint run`

## Code Style Guidelines
- Follow standard Go conventions (gofmt)
- Use meaningful variable/function names (camelCase)
- Group imports: standard lib first, then third-party, then internal
- Error handling: check errors immediately, avoid nested error handling
- Use context for request cancellation and timeouts
- Dependencies: use dependency injection where appropriate
- Comments: follow godoc conventions for exported items
- Testing: aim for >80% coverage of business logic
- Keep functions focused and under 50 lines where possible