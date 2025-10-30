# golang-vibe-coding

A comprehensive Go toolkit for GitHub Actions DSL with advanced debugging capabilities, IDE integration, and Visual Studio Code extension support.

## 🚀 Features

- **DSL Processing**: Parse custom DSL files and convert them to GitHub Actions YAML
- **Advanced Debugging**: Full-featured debugger with breakpoints, step execution, and AST visualization
- **IDE Integration**: Debug Adapter Protocol (DAP) server for VS Code and other editors
- **VS Code Extension**: Complete syntax highlighting and debugging support
- **Multiple Interfaces**: CLI, HTTP server, interactive debugger, and web UI
- **Cross-Platform**: Windows, Linux, and macOS support

## 📦 Installation

### From Source

1. Clone the repository:
```powershell
git clone <repository-url>
cd golang-vibe-coding
```

2. Build all binaries:
```powershell
.\scripts\build.ps1
# Or for specific components:
.\scripts\build.ps1 -Target main
```

### Pre-built Binaries

Download from the [releases page](../../releases) for your platform.

## 🛠 Usage

### Main CLI Application

```powershell
# Show help
.\bin\golang-vibe-coding.exe --help

# Compile DSL to GitHub Actions YAML
.\bin\golang-vibe-coding.exe compile input.dsl -o .github\workflows\ci.yml

# Generate example DSL file
.\bin\golang-vibe-coding.exe example

# Greet command (example)
.\bin\golang-vibe-coding.exe greet "World" -m "Hello"
```

### Debugging Tools

#### CLI Debugger
```powershell
# Interactive debugging session
.\bin\debugger.exe --interactive examples\simple.dsl

# Batch analysis with visualization
.\bin\debugger.exe --visualize --analyze examples\production.dsl

# Start web UI debugger
.\bin\debugger.exe --web --port 8080 examples\simple.dsl
```

#### DAP Server (for VS Code)
```powershell
# Start DAP server for VS Code integration
.\bin\dap.exe --dap-server --dap-port 4711
```

### VS Code Extension

1. Install the extension from `vscode-extension/github-actions-dsl-1.0.0.vsix`
2. Open a `.dsl` file to get syntax highlighting
3. Use F5 to start debugging with full breakpoint support

## 🏗 Project Structure

```
├── cmd/                    # Command-line applications
│   ├── golang-vibe-coding/ # Main CLI application
│   ├── debugger/          # Interactive debugger
│   └── dap/               # Debug Adapter Protocol server
├── internal/              # Private application code
│   ├── ast/               # Abstract Syntax Tree definitions
│   ├── lexer/             # Tokenization
│   ├── parser/            # DSL parsing
│   └── generator/         # YAML generation logic
├── pkg/                   # Public API packages
│   └── debug/             # Debug engine and utilities
├── docs/                  # Documentation
├── examples/              # Example DSL files
├── scripts/               # Build and utility scripts
├── vscode-extension/      # VS Code extension
├── bin/                   # Compiled binaries
└── .vscode/               # VS Code configuration
```

## 🔧 Development

### Prerequisites

- Go 1.24+ 
- PowerShell 7+ (for scripts)
- VS Code (for extension development)
- Node.js 18+ (for VS Code extension)

### Building

```powershell
# Build all components
.\scripts\build.ps1

# Build specific targets
.\scripts\build.ps1 -Target main      # Main CLI only
.\scripts\build.ps1 -Target debugger  # Debugger only
.\scripts\build.ps1 -Target dap       # DAP server only

# Cross-compilation
.\scripts\build.ps1 -OS linux -Arch amd64
.\scripts\build.ps1 -OS darwin -Arch arm64

# Clean build
.\scripts\build.ps1 -Clean
```

### Testing

```powershell
# Run all tests
.\scripts\test.ps1

# Run with coverage
.\scripts\test.ps1 -Coverage

# Run with race detection
.\scripts\test.ps1 -Race -Verbose

# Run specific tests
.\scripts\test.ps1 -Pattern "*Parser*"
```

### VS Code Tasks

The project includes pre-configured VS Code tasks:

- **Build Go CLI**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Build Go CLI"
- **Run Go CLI**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Run Go CLI"  
- **Test Go CLI**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Test Go CLI"

## 📚 Documentation

- [Debug Guide](docs/DEBUG_GUIDE.md) - Comprehensive debugging documentation
- [VS Code Debug Guide](docs/VSCODE_DEBUG_GUIDE.md) - VS Code integration setup
- [Debugger Implementation Summary](docs/DEBUGGER_IMPLEMENTATION_SUMMARY.md) - Technical details

## 🎯 DSL Syntax

The DSL supports GitHub Actions workflow definitions with simplified syntax:

```dsl
workflow "CI Pipeline" {
    on: [push, pull_request]
    
    job "test" {
        runs-on: ubuntu-latest
        
        step "checkout" {
            uses: actions/checkout@v4
        }
        
        step "setup" {
            uses: actions/setup-go@v4
            with: {
                go-version: "1.24"
            }
        }
        
        step "test" {
            run: "go test ./..."
        }
    }
}
```

See `examples/` directory for more complex examples including:
- Matrix builds
- Conditional execution  
- Environment variables
- Multi-job workflows
- Dependency management

## 🐛 Debugging Features

### Interactive Debugger Commands

- `parse` - Parse DSL with error reporting
- `lex` - Tokenize DSL with loop protection
- `visualize` - Show AST tree visualization
- `analyze` - Display AST statistics
- `breakpoint` - Set execution breakpoints
- `step` - Step through parsing
- `vars` - Show current variables
- `reload` - Reload DSL file

### Web UI Debugger

Access at `http://localhost:8080/debug/ui/` when running:
```powershell
.\bin\debugger.exe --web examples\simple.dsl
```

Features:
- Real-time AST visualization
- Interactive breakpoints
- Step-by-step execution
- Variable inspection
- Syntax highlighting

## 🔌 VS Code Integration

The included VS Code extension provides:

- **Syntax Highlighting**: Full TextMate grammar for `.dsl` files
- **Debug Support**: Native DAP integration with breakpoints
- **IntelliSense**: Basic completion and error reporting
- **Commands**: Compile, debug, and analyze DSL files
- **Themes**: Support for all VS Code color themes

### Extension Installation

```powershell
# Install from VSIX
code --install-extension vscode-extension\github-actions-dsl-1.0.0.vsix

# Or in VS Code: Extensions → Install from VSIX
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following Go best practices
4. Add tests for new functionality
5. Run tests: `.\scripts\test.ps1`
6. Build all targets: `.\scripts\build.ps1`
7. Commit changes: `git commit -m 'Add amazing feature'`
8. Push to branch: `git push origin feature/amazing-feature`
9. Open a Pull Request

### Code Standards

- Follow Go naming conventions
- Add comprehensive tests
- Update documentation
- Use PowerShell for scripts (Windows-first approach)
- Maintain backward compatibility in public APIs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Microsoft Debug Adapter Protocol specification
- VS Code extension development guidance
- Go community best practices
- GitHub Actions DSL inspiration
