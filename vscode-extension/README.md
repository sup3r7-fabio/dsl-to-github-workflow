# GitHub Actions DSL - VSCode Extension

A comprehensive Visual Studio Code extension that provides advanced language support and debugging capabilities for GitHub Actions DSL files.

## Features

### 🎨 Syntax Highlighting
- Full syntax highlighting for DSL files with `.dsl` extension
- Support for GitHub Actions workflow constructs
- Expression interpolation highlighting (`${{ }}`)
- Action reference highlighting (`owner/repo@version`)
- Event and runner label highlighting

### 🐛 Debugging Support
- Full Debug Adapter Protocol (DAP) integration
- Set breakpoints in DSL files
- Step through workflow execution
- Variable inspection and watch
- Call stack visualization
- Debug console with expression evaluation

### 🧠 IntelliSense & Language Features
- Auto-completion for workflow keywords
- Hover information for DSL constructs
- Real-time syntax validation
- Error highlighting and diagnostics

### ⚡ Commands & Tools
- **Start/Stop/Restart DAP Server**: Manage the debug server
- **Validate DSL Syntax**: Check syntax of current file
- **Show AST Visualization**: View abstract syntax tree
- **Open Debug Console**: Quick access to debug tools

## Installation

### From Source
1. Clone the repository
2. Navigate to the `vscode-extension` directory
3. Run `npm install`
4. Run `npm run compile`
5. Open VSCode and press `F5` to launch extension development host

### From VSIX Package
1. Build the extension: `npm run package`
2. Install in VSCode: `code --install-extension github-actions-dsl-1.0.0.vsix`

## Usage

### Basic Language Support
1. Create or open a file with `.dsl` extension
2. The extension automatically provides syntax highlighting
3. Use `Ctrl+Space` for auto-completion
4. Hover over keywords for documentation

### Debugging Workflows
1. Open a DSL file in VSCode
2. Set breakpoints by clicking in the editor gutter
3. Press `F5` or use "Run > Start Debugging"
4. Select "Debug DSL File" configuration
5. The extension will automatically start the DAP server if needed

### Debug Configuration
Add to your `.vscode/launch.json`:
```json
{
  "type": "dsl",
  "request": "launch",
  "name": "Debug DSL File",
  "program": "${file}",
  "stopOnEntry": true,
  "dapServerPort": 4711,
  "autoStartDapServer": true
}
```

## Commands

Access these commands via `Ctrl+Shift+P`:

- `GitHub Actions DSL: Start DAP Server` - Start the debug server
- `GitHub Actions DSL: Stop DAP Server` - Stop the debug server
- `GitHub Actions DSL: Restart DAP Server` - Restart the debug server
- `GitHub Actions DSL: Validate DSL Syntax` - Validate current file
- `GitHub Actions DSL: Show AST Visualization` - Show syntax tree
- `GitHub Actions DSL: Open Debug Console` - Open debug console

## Configuration

### Extension Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `github-actions-dsl.dapServer.autoStart` | `true` | Auto-start DAP server when debugging |
| `github-actions-dsl.dapServer.port` | `4711` | Default DAP server port |
| `github-actions-dsl.dapServer.host` | `localhost` | Default DAP server host |
| `github-actions-dsl.dapServer.executable` | `./bin/dsl-dap` | Path to DAP server executable |
| `github-actions-dsl.syntax.validation.enabled` | `true` | Enable real-time syntax validation |
| `github-actions-dsl.debug.verbose` | `false` | Enable verbose debug logging |
| `github-actions-dsl.ast.autoShow` | `false` | Auto-show AST when opening DSL files |

### File Associations
The extension automatically associates `.dsl` files with the GitHub Actions DSL language. You can also manually set the language mode using `Ctrl+K M`.

## Debugging Features

### Breakpoints
- **Line Breakpoints**: Click in the gutter to set/unset
- **Conditional Breakpoints**: Right-click breakpoint for conditions
- **Logpoints**: Add logging without stopping execution

### Variable Inspection
- **Variables Panel**: Shows workflow context, job info, step details
- **Watch Expressions**: Monitor specific values
- **Hover Evaluation**: Hover over variables to see values

### Debug Console
- Evaluate expressions in the current context
- Access workflow variables: `workflow.name`, `job.id`, `step.name`
- Check environment variables: `env.NODE_VERSION`

### Step Controls
- **Continue** (F5): Resume execution
- **Step Over** (F10): Execute next line
- **Step Into** (F11): Step into function calls  
- **Step Out** (Shift+F11): Step out of current function

## Requirements

### Prerequisites
1. **Go Environment**: Go 1.19+ installed
2. **DSL Parser**: The main DSL parser project built
3. **DAP Server**: The `dsl-dap` executable built

### Building Dependencies
```bash
# Build the main DSL parser
go build -o bin/golang-vibe-coding .

# Build the DAP server
go build -o bin/dsl-dap ./cmd/dap

# Verify builds
ls bin/
```

## Architecture

### Extension Components
- **Main Extension** (`extension.ts`): Core functionality and command registration
- **Debug Adapter** (`debugAdapter.ts`): DAP integration and communication
- **Language Features**: Completion, hover, and validation providers
- **Syntax Grammar** (`dsl.tmGrammar.json`): TextMate grammar for highlighting

### Integration Flow
```
VSCode Extension → DAP Adapter → TCP Connection → Go DAP Server → DSL Parser
```

## Development

### Extension Development
1. Clone the repository
2. Navigate to `vscode-extension/`
3. Run `npm install`
4. Press `F5` to launch extension development host
5. Open a `.dsl` file to test features

### Debugging the Extension
1. Set breakpoints in TypeScript files
2. Use the Extension Development Host
3. Check the Debug Console for logs
4. Use "Developer: Reload Window" to reload changes

### Building VSIX Package
```bash
npm run package
# Creates: github-actions-dsl-1.0.0.vsix
```

## Example DSL File

```yaml
# example.dsl
name: CI/CD Pipeline
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
      
      - name: Run Tests
        run: npm test
        env:
          CI: true
```

## Troubleshooting

### Common Issues

#### Extension Not Activating
- Check that `.dsl` files are present in workspace
- Verify extension is enabled in VSCode
- Check Output panel for error messages

#### Syntax Highlighting Missing
- Ensure file has `.dsl` extension
- Try setting language mode manually: `Ctrl+K M` → "GitHub Actions DSL"
- Reload window: `Ctrl+Shift+P` → "Developer: Reload Window"

#### Debugging Not Working
- Verify DAP server is built: `ls bin/dsl-dap*`
- Check DAP server port is available: `netstat -an | findstr 4711`
- Enable verbose logging in settings
- Check Debug Console for error messages

#### Auto-completion Not Working
- Ensure cursor is in valid context (after keywords)
- Try triggering manually with `Ctrl+Space`
- Check that language server features are enabled

### Debug Information
Enable verbose logging:
1. Open Settings: `Ctrl+,`
2. Search for "github-actions-dsl"
3. Enable "Debug: Verbose"
4. Check Output panel → "GitHub Actions DSL"

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly with sample DSL files
5. Submit a pull request

## License

This extension is part of the golang-vibe-coding DSL project and follows the same license terms.

## Support

- 📖 Documentation: See project README and guides
- 🐛 Issues: Report bugs in the project repository
- 💬 Discussions: Use GitHub Discussions for questions

## Related

- [GitHub Actions DSL Parser](https://github.com/golang-vibe-coding/golang-vibe-coding/blob/main/README.md)
- [Debug Guide](https://github.com/golang-vibe-coding/golang-vibe-coding/blob/main/DEBUG_GUIDE.md)  
- [VSCode Debug Guide](https://github.com/golang-vibe-coding/golang-vibe-coding/blob/main/VSCODE_DEBUG_GUIDE.md)
- [Implementation Summary](https://github.com/golang-vibe-coding/golang-vibe-coding/blob/main/DEBUGGER_IMPLEMENTATION_SUMMARY.md)
