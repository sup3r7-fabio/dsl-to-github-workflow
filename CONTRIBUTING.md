# Contributing to golang-vibe-coding

Thank you for your interest in contributing! This document provides guidelines and information for contributing to the golang-vibe-coding project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://golang.org/conduct). Please read and follow it in all your interactions.

## Getting Started

### Prerequisites

- **Go 1.24+**: Latest stable Go version
- **PowerShell 7+**: For build scripts (cross-platform)
- **VS Code**: Recommended for development
- **Git**: Version control
- **Node.js 18+**: Required for VS Code extension development

### Development Environment Setup

1. **Fork and Clone**:
   ```powershell
   git clone https://github.com/YOUR_USERNAME/golang-vibe-coding.git
   cd golang-vibe-coding
   ```

2. **Install Dependencies**:
   ```powershell
   go mod download
   go mod tidy
   ```

3. **Build Project**:
   ```powershell
   .\scripts\build.ps1
   ```

4. **Run Tests**:
   ```powershell
   .\scripts\test.ps1
   ```

## Development Workflow

### Branch Naming

- `feature/description` - New features
- `bugfix/description` - Bug fixes  
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring
- `test/description` - Test improvements

### Development Process

1. **Create Branch**:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Make Changes**: Follow coding standards and write tests

3. **Test Locally**:
   ```powershell
   .\scripts\test.ps1 -Coverage -Race
   .\scripts\build.ps1
   ```

4. **Commit Changes**:
   ```bash
   git add .
   git commit -m "Add awesome feature with comprehensive tests"
   ```

5. **Push and PR**:
   ```bash
   git push origin feature/my-awesome-feature
   ```

## Project Structure

### Directory Organization

```
├── cmd/                    # Command-line applications (main packages)
│   ├── golang-vibe-coding/ # Main CLI application
│   ├── debugger/          # Interactive debugger
│   └── dap/               # Debug Adapter Protocol server
├── internal/              # Private application code
│   ├── ast/               # Abstract Syntax Tree definitions
│   ├── lexer/             # Tokenization logic
│   ├── parser/            # DSL parsing logic
│   └── generator/         # YAML generation logic
├── pkg/                   # Public API packages
│   └── debug/             # Debug engine and utilities
├── docs/                  # Documentation
├── examples/              # Example DSL files and usage
├── scripts/               # Build and utility scripts
├── vscode-extension/      # VS Code extension
└── .vscode/               # VS Code workspace configuration
```

### Package Guidelines

- **`cmd/`**: Main applications, keep minimal (mainly argument parsing and setup)
- **`internal/`**: Private code, not importable by external projects
- **`pkg/`**: Public APIs that external projects can import
- **`docs/`**: All documentation beyond README
- **`examples/`**: Working examples and test DSL files
- **`scripts/`**: PowerShell scripts for build, test, and CI/CD

## Coding Standards

### Go Code Style

1. **Follow Go Standards**:
   - Use `go fmt`
   - Follow [Effective Go](https://golang.org/doc/effective_go.html)
   - Use `golangci-lint` for linting

2. **Naming Conventions**:
   ```go
   // Good
   type WorkflowParser struct {}
   func (p *WorkflowParser) ParseWorkflow() *ast.Workflow {}
   
   // Package names: short, lowercase, no underscores
   package lexer  // not github_lexer or githubLexer
   ```

3. **Error Handling**:
   ```go
   // Always handle errors explicitly
   result, err := someFunction()
   if err != nil {
       return fmt.Errorf("failed to do something: %w", err)
   }
   
   // Use wrapped errors for context
   return fmt.Errorf("parsing workflow %s: %w", filename, err)
   ```

4. **Comments**:
   ```go
   // Package comments should be complete sentences
   // Package lexer provides tokenization for GitHub Actions DSL.
   package lexer
   
   // Public functions need comments
   // ParseWorkflow parses a DSL string into an AST workflow.
   func ParseWorkflow(input string) (*ast.Workflow, error) {}
   ```

### PowerShell Scripts

1. **Use Approved Verbs**: Follow PowerShell verb guidelines
2. **Parameter Validation**: Use proper parameter types and validation
3. **Error Handling**: Use `$ErrorActionPreference = "Stop"`
4. **Documentation**: Include help comments and examples

### VS Code Extension

1. **TypeScript**: Use strict type checking
2. **Follow VS Code Guidelines**: Use official VS Code extension guidelines
3. **Async/Await**: Prefer async/await over promises where possible

## Testing

### Test Structure

```powershell
# Run all tests
.\scripts\test.ps1

# Run with coverage  
.\scripts\test.ps1 -Coverage

# Run with race detection
.\scripts\test.ps1 -Race -Verbose

# Run specific package tests
go test ./internal/parser -v
```

### Test Guidelines

1. **Test Files**: Use `*_test.go` naming
2. **Test Functions**: Use `TestFunctionName` pattern
3. **Table Tests**: Use table-driven tests for multiple scenarios
4. **Coverage**: Aim for >80% coverage on new code
5. **Benchmarks**: Add benchmarks for performance-critical code

Example test structure:
```go
func TestParseWorkflow(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *ast.Workflow
        wantErr  bool
    }{
        {
            name:  "simple workflow",
            input: `workflow "test" { on = "push" }`,
            expected: &ast.Workflow{Name: "test"},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseWorkflow(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseWorkflow() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Additional assertions...
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. **Pre-submission Checklist**:
   - [ ] Tests pass: `.\scripts\test.ps1`
   - [ ] Builds successfully: `.\scripts\build.ps1`  
   - [ ] Linting passes: `golangci-lint run`
   - [ ] Documentation updated
   - [ ] Examples work with changes

2. **PR Description Template**:
   ```markdown
   ## Description
   Brief description of changes and motivation.
   
   ## Type of Change
   - [ ] Bug fix (non-breaking change)
   - [ ] New feature (non-breaking change) 
   - [ ] Breaking change (fix or feature causing existing functionality to change)
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests added/updated
   - [ ] Integration tests pass
   - [ ] Manual testing completed
   
   ## Checklist
   - [ ] Code follows project style guidelines
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] No new warnings introduced
   ```

3. **Review Process**:
   - At least one maintainer approval required
   - All CI checks must pass
   - Address review feedback promptly
   - Keep PR focused and reasonably sized

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code refactoring
- `test`: Adding/updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(parser): add support for matrix workflows

Add parsing logic for matrix workflow syntax including
strategy definitions and variable substitution.

Fixes #123
```

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)
- Pre-releases: `1.2.3-alpha.1`, `1.2.3-beta.1`, `1.2.3-rc.1`

### Release Steps

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md with new features and fixes
3. **Tag**: Create git tag with version
4. **Build**: Generate release binaries for all platforms
5. **VS Code Extension**: Update and package extension
6. **GitHub Release**: Create GitHub release with artifacts

## Getting Help

- **Issues**: Create GitHub issue for bugs/features
- **Discussions**: Use GitHub Discussions for questions
- **Code Review**: Request review from maintainers
- **Documentation**: Check `docs/` directory for detailed guides

## Recognition

Contributors are recognized in:
- GitHub contributors list
- Release notes for significant contributions  
- Special recognition for first-time contributors

Thank you for contributing to golang-vibe-coding! 🚀
