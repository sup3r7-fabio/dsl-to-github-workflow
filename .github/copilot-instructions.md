# Go CLI Project Instructions

This workspace contains a Go command-line interface (CLI) application.

## Project Structure
- Entry point: `main.go`
- Go module: `go.mod`
- Build output: `bin/golang-vibe-coding`
- Commands: Simple flag-based CLI (ready for Cobra if needed)

## Development Guidelines
- Follow Go naming conventions and best practices
- Use proper error handling with wrapped errors
- Include proper help text and usage examples
- Write unit tests for core functionality
- Use VS Code tasks for building and running

## Build and Run
- Build: `go build -o bin/golang-vibe-coding .`
- Run: `go run . [args]`
- Test: `go test ./...`

## Available Commands
- `--help` - Show help information
- `--version` - Show version information  
- `greet [name]` - Greet someone

## Dependencies
- Manage dependencies using `go mod`
- Keep go.sum file committed
