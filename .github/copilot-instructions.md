# DSL-to-GitHub-Workflow CLI Project

This workspace contains a Go CLI application that converts DSL files to GitHub Actions workflows.

## Project Structure
- Entry point: `cmd/golang-vibe-coding/main.go`
- Go module: `go.mod`
- Internal packages: `internal/` (lexer, parser, ast, generator)
- Public packages: `pkg/` (debug utilities)
- VS Code extension: `vscode-extension/`
- Examples: `examples/` (DSL sample files)
- Documentation: `docs/`

## CI/CD Workflows
- **CI**: Build, test, and lint on push/PR
- **Release**: Cross-platform builds and GitHub releases  
- **Security**: Dependency scanning and vulnerability checks
- **VS Code Extension**: Extension testing and marketplace publishing
- **Documentation**: Example validation and link checking

## Development Guidelines
- Follow Go naming conventions and best practices
- Use proper error handling with wrapped errors
- Write unit tests for core functionality
- Validate DSL examples in CI
- Keep documentation up-to-date

## Build and Run
- Build: `go build -o bin/dsl-cli ./cmd/golang-vibe-coding`
- Run: `go run ./cmd/golang-vibe-coding [args]`
- Test: `go test ./...`
- Lint: `golangci-lint run`

## Available Commands
- `compile [input.dsl]` - Convert DSL to GitHub Actions YAML
- `example [type]` - Show example DSL files
- `greet [name]` - Greet someone (demo command)
- `--help` - Show help information
- `--version` - Show version information

## DSL Features
- Loop constructs: `for`, `foreach`, `repeat`
- Variable interpolation: `${{ vars.name }}`
- GitHub Actions workflow generation
- Debug support with VS Code extension

## Release Process
1. Create a tag: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. GitHub Actions will build cross-platform binaries and create a release

## Dependencies
- CLI framework: `github.com/urfave/cli/v2`
- YAML processing: `gopkg.in/yaml.v3`
- Manage dependencies using `go mod`
- Dependabot keeps dependencies updated automatically
