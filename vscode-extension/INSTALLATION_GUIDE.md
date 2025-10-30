# VSCode Extension Installation & Testing Guide

This guide explains how to install and test the GitHub Actions DSL VSCode extension.

## 📦 Installation Methods

### Method 1: Install from VSIX Package (Recommended for Testing)

1. **Locate the VSIX file**:
   ```
   s:\development\github.com\golang-vibe-coding\vscode-extension\github-actions-dsl-1.0.0.vsix
   ```

2. **Install using VSCode UI**:
   - Open VSCode
   - Press `Ctrl+Shift+P` to open command palette
   - Type "Extensions: Install from VSIX"
   - Select the command and choose the `.vsix` file
   - Restart VSCode when prompted

3. **Install using Command Line**:
   ```powershell
   cd "s:\development\github.com\golang-vibe-coding\vscode-extension"
   code --install-extension github-actions-dsl-1.0.0.vsix
   ```

### Method 2: Development Mode (For Development & Testing)

1. **Open Extension Development Host**:
   ```powershell
   cd "s:\development\github.com\golang-vibe-coding\vscode-extension"
   code .
   # Press F5 to launch Extension Development Host
   ```

2. **Test in Development Host**:
   - A new VSCode window opens with the extension loaded
   - Open or create `.dsl` files to test functionality
   - Make changes to extension code and reload to test

## 🧪 Testing the Extension

### Prerequisites
Ensure the DAP server is built before testing debugging features:
```powershell
cd "s:\development\github.com\golang-vibe-coding"
go build -o bin/dsl-dap ./cmd/dap
go build -o bin/golang-vibe-coding .
```

### Test 1: Syntax Highlighting
1. Open the example DSL file:
   ```
   s:\development\github.com\golang-vibe-coding\vscode-extension\example.dsl
   ```
2. Verify syntax highlighting is working:
   - Keywords like `name`, `on`, `jobs`, `steps` should be highlighted
   - Strings in quotes should be colored differently
   - Comments (lines starting with `#`) should be styled as comments
   - Expression interpolation (`${{ }}`) should be highlighted

### Test 2: Auto-completion
1. Create a new file with `.dsl` extension
2. Start typing workflow keywords:
   - Type `na` and press `Ctrl+Space` → should suggest `name`
   - Type `o` and press `Ctrl+Space` → should suggest `on`
   - Type `jo` and press `Ctrl+Space` → should suggest `jobs`
3. Verify completion works in different contexts

### Test 3: Hover Information
1. Open a DSL file
2. Hover over keywords like:
   - `name` → Should show "The name of the workflow"
   - `runs-on` → Should show "The type of machine to run the job on"
   - `uses` → Should show "Selects an action to run as part of a step"

### Test 4: Commands
Test extension commands via `Ctrl+Shift+P`:

1. **Start DAP Server**:
   - Command: "GitHub Actions DSL: Start DAP Server"
   - Should show success message and start server on port 4711

2. **Validate Syntax**:
   - Command: "GitHub Actions DSL: Validate DSL Syntax"
   - Open a DSL file first, then run command
   - Should show validation result

3. **Show AST Visualization**:
   - Command: "GitHub Actions DSL: Show AST Visualization" 
   - Should open a new panel with AST information

### Test 5: Debugging (Full DAP Integration)

#### Setup Debug Configuration
1. Create `.vscode/launch.json` in your workspace:
   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "type": "dsl",
         "request": "launch",
         "name": "Debug DSL File",
         "program": "${file}",
         "stopOnEntry": true,
         "autoStartDapServer": true
       }
     ]
   }
   ```

#### Test Debugging Flow
1. **Set Breakpoints**:
   - Open a DSL file (`example.dsl`)
   - Click in the gutter next to line numbers to set breakpoints
   - Breakpoints should appear as red dots

2. **Start Debug Session**:
   - Press `F5` or go to "Run > Start Debugging"
   - Select "Debug DSL File" configuration
   - Extension should automatically start DAP server
   - Debug session should begin

3. **Debug Controls**:
   - Test Continue (F5)
   - Test Step Over (F10)
   - Test Step Into (F11)
   - Test Step Out (Shift+F11)

4. **Variable Inspection**:
   - Check Variables panel for workflow context
   - Use Debug Console to evaluate expressions
   - Test hover evaluation over variables

### Test 6: Error Handling

1. **Invalid DSL Syntax**:
   - Create DSL file with syntax errors
   - Should see error highlighting
   - Validation command should report errors

2. **Missing DAP Server**:
   - Try debugging without DAP server built
   - Should prompt to build server first
   - Should handle connection failures gracefully

3. **Invalid Debug Configuration**:
   - Try debugging with malformed launch.json
   - Should show helpful error messages

## 📊 Expected Results

### ✅ Success Indicators
- Syntax highlighting works for all DSL constructs
- Auto-completion provides relevant suggestions
- Hover shows helpful documentation
- Commands execute without errors
- DAP server starts and accepts connections
- Breakpoints can be set and hit during debugging
- Variables and call stack are visible during debug sessions
- Debug controls (step, continue) work properly

### ❌ Common Issues and Solutions

#### Extension Not Loading
**Problem**: Extension doesn't appear to be active
**Solutions**:
- Check if `.dsl` files are present in workspace
- Reload window: `Ctrl+Shift+P` → "Developer: Reload Window"
- Check Output panel for errors: View → Output → "GitHub Actions DSL"

#### Syntax Highlighting Missing
**Problem**: DSL files show as plain text
**Solutions**:
- Ensure file has `.dsl` extension
- Set language manually: `Ctrl+K M` → Select "GitHub Actions DSL"
- Verify extension is enabled: Extensions → Search "GitHub Actions DSL"

#### DAP Server Won't Start
**Problem**: Debug server fails to start
**Solutions**:
```powershell
# Check if binary exists
ls bin/dsl-dap*

# Rebuild if missing
go build -o bin/dsl-dap ./cmd/dap

# Test manual start
.\bin\dsl-dap --dap-server --verbose

# Check port availability
netstat -an | findstr 4711
```

#### Debugging Not Working
**Problem**: Breakpoints not hit or debug session fails
**Solutions**:
- Ensure DAP server is running
- Check launch configuration is correct
- Verify file paths are accurate
- Enable verbose debug logging in settings
- Check Debug Console for error messages

#### Auto-completion Not Working
**Problem**: No suggestions appear
**Solutions**:
- Try triggering manually with `Ctrl+Space`
- Ensure cursor is in valid context
- Check that file is recognized as DSL language
- Restart VSCode if language features aren't loading

## 🔧 Development Testing

### Hot Reload Development
1. Open extension project in VSCode
2. Make changes to TypeScript files
3. Press `Ctrl+R` in Extension Development Host to reload
4. Test changes immediately

### Debug Extension Code
1. Set breakpoints in extension TypeScript files
2. Launch Extension Development Host (F5)
3. Extension code breakpoints will be hit when features are used
4. Use Debug Console to inspect extension state

### Build and Test Cycle
```powershell
# Make changes to extension
# Compile TypeScript
npm run compile

# Run tests (if available)
npm test

# Package for testing
npm run package

# Install and test
code --install-extension github-actions-dsl-1.0.0.vsix
```

## 📋 Testing Checklist

- [ ] Extension installs successfully
- [ ] Syntax highlighting works for DSL files
- [ ] Auto-completion provides relevant suggestions
- [ ] Hover information displays correctly  
- [ ] Extension commands are available and work
- [ ] DAP server starts without errors
- [ ] Breakpoints can be set in DSL files
- [ ] Debug sessions start successfully
- [ ] Debug controls (step, continue) work
- [ ] Variables are visible during debugging
- [ ] Debug console accepts expressions
- [ ] Error handling works gracefully
- [ ] Extension can be uninstalled cleanly

## 🚀 Distribution

Once testing is complete, the extension can be:
1. **Published to VSCode Marketplace** (requires publisher account)
2. **Distributed as VSIX file** for internal use
3. **Integrated into CI/CD** for automatic builds and testing

The packaged extension (`github-actions-dsl-1.0.0.vsix`) is ready for distribution and provides a complete debugging and development environment for GitHub Actions DSL files.
