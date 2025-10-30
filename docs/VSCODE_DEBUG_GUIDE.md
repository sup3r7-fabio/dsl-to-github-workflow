# VSCode Debug Adapter Protocol (DAP) Guide

This guide explains how to use the GitHub Actions DSL debugger with Visual Studio Code through the Debug Adapter Protocol (DAP).

## Overview

The DSL debugger now supports the Debug Adapter Protocol (DAP), which enables full integration with Visual Studio Code's debugging interface. This provides a professional IDE debugging experience with breakpoints, variable inspection, stack traces, and step-through debugging.

## Prerequisites

1. **Visual Studio Code** installed on your system
2. **Go environment** properly configured
3. **DSL debugger** compiled and ready to use

## Setup Instructions

### 1. Build the DAP Server

First, build the DAP adapter server:

```powershell
# Build the main DSL binary
go build -o bin/golang-vibe-coding .

# Build the DAP adapter (if separate)
go build -o bin/dsl-dap ./cmd/dap
```

### 2. VSCode Launch Configuration

Create or update `.vscode/launch.json` in your workspace:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug DSL File",
            "type": "dsl",
            "request": "launch",
            "program": "${file}",
            "stopOnEntry": false,
            "console": "integratedTerminal",
            "cwd": "${workspaceFolder}"
        },
        {
            "name": "Debug DSL with Arguments",
            "type": "dsl",
            "request": "launch",
            "program": "${workspaceFolder}/examples/simple.dsl",
            "args": ["--debug", "--verbose"],
            "stopOnEntry": true,
            "console": "integratedTerminal",
            "cwd": "${workspaceFolder}"
        }
    ]
}
```

### 3. VSCode Settings

Add the following to your `.vscode/settings.json`:

```json
{
    "debug.allowBreakpointsEverywhere": true,
    "debug.showBreakpointsInOverviewRuler": true,
    "files.associations": {
        "*.dsl": "yaml"
    }
}
```

## DAP Server Usage

### Starting the DAP Server

The DAP server can be started in two ways:

#### Method 1: Standalone Server
```powershell
# Start the DAP server on default port (4711)
.\bin\golang-vibe-coding --dap-server

# Start on custom port
.\bin\golang-vibe-coding --dap-server --dap-port 5000
```

#### Method 2: Launch Mode (Recommended)
The DAP adapter will automatically start when VSCode initiates a debug session.

### Debug Session Workflow

1. **Set Breakpoints**: Click in the gutter next to line numbers in your `.dsl` files
2. **Start Debugging**: Press F5 or use "Run > Start Debugging"
3. **Debug Controls**:
   - **Continue** (F5): Resume execution
   - **Step Over** (F10): Execute next line
   - **Step Into** (F11): Step into function calls
   - **Step Out** (Shift+F11): Step out of current function
   - **Stop** (Shift+F5): Terminate debug session

## Debugging Features

### 1. Breakpoints

Set breakpoints in your DSL files by clicking in the editor gutter:

```yaml
# example.dsl
name: Test Workflow
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout          # <- Set breakpoint here
        uses: actions/checkout@v3
      
      - name: Setup Node        # <- Or here
        uses: actions/setup-node@v3
        with:
          node-version: '18'
```

### 2. Variable Inspection

During debugging, you can inspect variables in the **Variables** panel:

- **Workflow Context**: Current workflow state
- **Job Context**: Active job information  
- **Step Context**: Current step details
- **Environment Variables**: Available env vars
- **Inputs/Outputs**: Step inputs and outputs

### 3. Call Stack

The **Call Stack** panel shows:
- Current execution position
- Workflow → Job → Step hierarchy
- Function call stack (if applicable)

### 4. Debug Console

Use the **Debug Console** to:
- Evaluate expressions
- Inspect variable values
- Execute debug commands

Example commands:
```
> workflow.name
"Test Workflow"

> job.id
"build"

> step.name
"Checkout"

> env.NODE_VERSION
"18"
```

## Advanced Features

### Conditional Breakpoints

Right-click on a breakpoint to add conditions:

```javascript
// Example conditions
step.name === "Setup Node"
job.runs-on === "ubuntu-latest"
env.NODE_VERSION !== undefined
```

### Logpoints

Add logpoints (breakpoints that log without stopping):
1. Right-click in gutter
2. Select "Add Logpoint"
3. Enter message: `Step: {step.name}, Status: {step.status}`

### Exception Breakpoints

Configure the debugger to break on:
- Parse errors
- Runtime errors
- Validation failures
- Timeout conditions

## Troubleshooting

### Common Issues

#### 1. DAP Server Not Starting
```powershell
# Check if port is available
netstat -an | findstr :4711

# Try different port
.\bin\golang-vibe-coding --dap-server --dap-port 5000
```

#### 2. Breakpoints Not Hit
- Ensure the DSL file is valid
- Check that debug symbols are enabled
- Verify the file path matches the workspace

#### 3. Variables Not Showing
- Make sure the debugger is properly initialized
- Check that the execution context is available
- Verify the DSL parser is running in debug mode

### Debug Logs

Enable detailed logging:

```powershell
# Set debug environment variables
$env:DSL_DEBUG = "3"
$env:DAP_DEBUG = "true"

# Run with verbose output
.\bin\golang-vibe-coding --dap-server --verbose
```

### Manual DAP Testing

Test the DAP server manually:

```powershell
# Start server
.\bin\golang-vibe-coding --dap-server --dap-port 4711

# In another terminal, connect via telnet
telnet localhost 4711
```

Send DAP initialize request:
```json
Content-Length: 87

{"seq":1,"type":"request","command":"initialize","arguments":{"clientID":"vscode"}}
```

## Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DSL_DEBUG` | Debug level (1-3) | `""` |
| `DAP_PORT` | DAP server port | `4711` |
| `DAP_HOST` | DAP server host | `localhost` |
| `DAP_TIMEOUT` | Connection timeout (seconds) | `30` |

### Command Line Options

| Option | Description |
|--------|-------------|
| `--dap-server` | Start DAP server mode |
| `--dap-port <port>` | Set DAP server port |
| `--dap-host <host>` | Set DAP server host |
| `--debug` | Enable debug mode |
| `--verbose` | Verbose output |

## Integration Examples

### Example 1: Debug Simple Workflow

```yaml
# simple-debug.dsl
name: Debug Example
on: [push]

jobs:
  debug-job:
    runs-on: ubuntu-latest
    steps:
      - name: Debug Step 1    # Breakpoint here
        run: echo "Starting debug"
      
      - name: Debug Step 2    # Breakpoint here
        run: echo "Variable: ${{ env.TEST_VAR }}"
        env:
          TEST_VAR: "debug-value"
```

### Example 2: Debug Complex Pipeline

```yaml
# complex-debug.dsl
name: Complex Debug
on: [push, pull_request]

jobs:
  matrix-debug:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest]
        node: ['16', '18', '20']
    runs-on: ${{ matrix.os }}
    steps:
      - name: Matrix Debug     # Breakpoint with condition: matrix.node === '18'
        run: |
          echo "OS: ${{ matrix.os }}"
          echo "Node: ${{ matrix.node }}"
```

## VSCode Extension Integration

For even better integration, consider creating a VSCode extension:

### Extension Structure
```
dsl-vscode-extension/
├── package.json
├── src/
│   ├── extension.ts
│   └── debugAdapter.ts
└── syntaxes/
    └── dsl.tmGrammar.json
```

### Extension Features
- Syntax highlighting for DSL files
- Integrated debugger activation
- Custom debug views
- IntelliSense support
- Error highlighting

## Best Practices

### 1. Effective Breakpoint Usage
- Set breakpoints at decision points
- Use conditional breakpoints for specific scenarios
- Place logpoints for monitoring without stopping

### 2. Variable Inspection Strategy
- Focus on workflow context first
- Check job-level variables for job-specific issues
- Inspect step inputs/outputs for data flow problems

### 3. Debug Session Management
- Start with simple workflows
- Use step-over for high-level flow understanding
- Use step-into for detailed execution analysis
- Monitor the call stack for context

### 4. Performance Considerations
- Limit breakpoints in loops
- Use conditional breakpoints judiciously
- Consider using logpoints instead of frequent stops

## Conclusion

The DAP integration provides a powerful debugging experience for GitHub Actions DSL development. With proper setup and usage of breakpoints, variable inspection, and debug controls, you can efficiently troubleshoot and develop complex workflows.

For additional help or advanced configuration, refer to the main [DEBUG_GUIDE.md](./DEBUG_GUIDE.md) or the [README.md](./README.md) documentation.
