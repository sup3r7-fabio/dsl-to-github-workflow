# DSL Debugger Guide

## Overview

This guide shows you how to create and use a language-specific debugger for your GitHub Actions DSL. The debugger provides multiple interfaces and debugging capabilities to help you understand, troubleshoot, and develop your DSL effectively.

## Features

### 🐛 Core Debugging Features
- **Breakpoints**: Set breakpoints at specific lines or functions
- **Step-through debugging**: Step through parsing and execution
- **Variable inspection**: View current variables and context
- **Call stack tracking**: Monitor function call hierarchy
- **AST visualization**: Visualize the parsed syntax tree
- **Interactive REPL**: Debug interactively with command-line interface
- **Web UI**: Browser-based debugging interface
- **IDE integration**: VS Code debug configurations

### 📊 Analysis Tools
- **Token analysis**: Inspect lexical tokens
- **Parse tree visualization**: See how your DSL is structured
- **Complexity analysis**: Get metrics about your workflows
- **Error reporting**: Detailed parsing error messages

## Quick Start

### 1. Enable Debug Mode

Set the `DSL_DEBUG` environment variable to enable debugging:

```powershell
# Basic debugging
$env:DSL_DEBUG = "1"

# Verbose output  
$env:DSL_DEBUG = "verbose"

# Full tracing
$env:DSL_DEBUG = "trace"
```

### 2. Run with Debugging

```powershell
# Test the debug system
go run test-debug.go

# Debug a specific DSL file
go run main.go examples/simple.dsl
```

### 3. Use VS Code Integration

Open VS Code and use the pre-configured debug configurations:
- **Debug DSL Parser**: Debug the test file with verbose output
- **Debug DSL with Trace**: Debug with full tracing enabled
- **Debug DSL Main Program**: Debug the main CLI with a sample file
- **Debug DSL with Custom File**: Debug with a file you specify

## Debugging Interfaces

### Command Line Interface

```powershell
# Basic usage
$env:DSL_DEBUG = "verbose"; go run test-debug.go

# With analysis and visualization
$env:DSL_DEBUG = "trace"; go run test-debug.go
```

### Interactive Debugger

The interactive debugger provides a REPL-like interface for step-by-step debugging:

```powershell
# Start interactive mode (when implemented)
go run dsl-debugger.go --interactive examples/simple.dsl
```

Available commands:
- `help` - Show available commands
- `lex` - Tokenize the DSL and show tokens
- `parse` - Parse the DSL into AST
- `visualize` - Show AST tree visualization
- `analyze` - Show AST analysis and metrics
- `breakpoint <line>` - Set a breakpoint
- `step` - Step through execution
- `continue` - Continue execution
- `vars` - Show current variables
- `quit` - Exit debugger

### Web Interface

Start the debug server for browser-based debugging:

```powershell
# Start with web UI
$env:DSL_DEBUG = "verbose"; go run test-debug.go
```

Then visit `http://localhost:8080/debug/ui/` for the web interface.

The web UI provides:
- **Breakpoint management**: Add, remove, enable/disable breakpoints
- **Variable inspection**: Real-time variable monitoring
- **Call stack viewer**: See function call hierarchy
- **Console**: Execute debug commands and see output
- **Step controls**: Step, continue, restart, stop

## Setting Breakpoints

### Programmatic Breakpoints

```go
import "golang-vibe-coding/pkg/debug"

// Set a breakpoint at a specific location
debug.SetBreakpoint("parser.go", 45, "parseWorkflow")

// Check if execution should pause
variables := map[string]interface{}{
    "current_token": token.Literal,
    "line": token.Line,
}
debug.CheckBreakpoint("parser.go", 45, "parseWorkflow", variables)
```

### Interactive Breakpoints

In the interactive debugger:
```
(dsl-debug) breakpoint 25
(dsl-debug) breakpoint parser.go:45
```

## AST Visualization

### Text-based Visualization

```go
visualizer := debug.NewASTVisualizer(debug.DefaultVisualizationOptions())
visualization := visualizer.VisualizeWorkflow(workflow)
fmt.Println(visualization)
```

Example output:
```
🔄 Workflow Test Workflow <workflow>
├── 🎯 Triggers <triggers>
│   ├── 📡 push <trigger>
├── 🌍 Global Environment <env>
│   ├──   NODE_ENV = development
│   ├──   API_URL = https://api.example.com
├── ⚙️ Jobs (2) <jobs>
│   ├── 🏗️ Job build <job>
│   │   ├──   runs-on: ubuntu-latest
│   │   ├── 📋 Steps (5) <steps>
│   │   │   ├── 📝 Step 1: Checkout <step>
│   │   │   │   ├──   uses: actions/checkout@v3
│   │   │   ├── 📝 Step 2: Setup Node <step>
│   │   │   │   ├──   uses: actions/setup-node@v3
│   │   │   │   ├── ⚙️ With <with>
│   │   │   │   │   ├──   node-version = 18
│   │   │   │   │   ├──   cache = npm
```

### HTML Visualization

Convert to HTML for web display:
```go
html := debug.VisualizationToHTML(visualization)
```

## Analysis and Metrics

### AST Analysis

```go
analysis := debug.AnalyzeAST(workflow)
```

Returns metrics like:
```go
{
    "workflow_name": "Test Workflow",
    "total_jobs": 2,
    "total_steps": 5,
    "total_triggers": 1,
    "matrix_jobs": 0,
    "conditional_jobs": 1,
    "complexity_score": 12,
    "has_global_env": true
}
```

### Complexity Scoring

The complexity score is calculated based on:
- Number of jobs (×2 points each)
- Number of steps (×1 point each)  
- Matrix jobs (×5 points each)
- Conditional jobs (×2 points each)

## Debug Output Levels

### Basic (`DSL_DEBUG=1` or `DSL_DEBUG=basic`)
- Basic parsing information
- Error messages
- Completion status

### Verbose (`DSL_DEBUG=verbose`)  
- Detailed parsing steps
- Token information
- Function entry/exit
- Variable values

### Trace (`DSL_DEBUG=trace`)
- Full execution tracing
- Complete call stack
- All debug output
- Performance timing

## Integration Examples

### Adding Debug Points to Your Parser

```go
func (p *Parser) parseWorkflow() *ast.Workflow {
    debug.Trace("parseWorkflow", "enter")
    defer debug.Trace("parseWorkflow", "exit")
    
    variables := map[string]interface{}{
        "current_token": p.curToken.Literal,
        "position": fmt.Sprintf("%d:%d", p.curToken.Line, p.curToken.Column),
    }
    
    debug.CheckBreakpoint("parser.go", 100, "parseWorkflow", variables)
    
    // ... parsing logic ...
    
    debug.Log("Parsed workflow: %s", workflow.Name)
    return workflow
}
```

### Custom Visualization Options

```go
options := debug.VisualizationOptions{
    ShowTypes:      true,
    ShowPositions:  true,
    ShowMetadata:   false,
    MaxDepth:       5,
    HighlightNodes: []string{"deploy", "test"},
}

visualizer := debug.NewASTVisualizer(options)
```

## Best Practices

### 1. Strategic Breakpoint Placement
- Place breakpoints at key parsing decision points
- Add breakpoints before complex operations
- Use conditional breakpoints for specific scenarios

### 2. Meaningful Variable Context
- Include relevant parsing state in variable maps
- Add position information (line, column)
- Include current token information

### 3. Structured Debug Output
- Use consistent debug.Log() formatting
- Group related debug output
- Include context in trace messages

### 4. Performance Considerations
- Debug overhead is only active when `DSL_DEBUG` is set
- Use appropriate debug levels for your needs
- Consider disabling verbose output in production

## Troubleshooting

### Common Issues

**Problem**: Debugger not starting
- **Solution**: Ensure `DSL_DEBUG` environment variable is set

**Problem**: Breakpoints not triggering  
- **Solution**: Check that `debug.CheckBreakpoint()` is called in your code

**Problem**: Web UI not accessible
- **Solution**: Verify the debug server started on the correct port

**Problem**: Variables not showing
- **Solution**: Ensure variables map is populated before calling `CheckBreakpoint()`

### Debug the Debugger

If the debugger itself has issues:

```powershell
# Enable Go debugging
$env:GODEBUG = "schedtrace=1000"

# Run with verbose Go output
go run -x test-debug.go
```

## Advanced Usage

### Custom Debug Commands

Extend the interactive debugger with custom commands by modifying the `handleCommand()` function in the debugger interface.

### API Integration

The debug server provides REST API endpoints:
- `GET /debug/session/` - List debug sessions
- `POST /debug/breakpoints` - Add breakpoints  
- `GET /debug/variables` - Get current variables
- `POST /debug/step` - Step through execution
- `POST /debug/continue` - Continue execution

### VS Code Extension Integration

The debug server is compatible with VS Code's Debug Adapter Protocol (DAP) for potential future VS Code extension integration.

## Conclusion

This debugger provides comprehensive debugging capabilities for your DSL, from simple logging to interactive step-through debugging with visual representations. Use the appropriate level of debugging for your current needs, and leverage the different interfaces (CLI, interactive, web) based on your workflow preferences.
