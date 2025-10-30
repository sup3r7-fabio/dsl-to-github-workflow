"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = __importStar(require("vscode"));
const path = __importStar(require("path"));
const fs = __importStar(require("fs"));
const child_process_1 = require("child_process");
const debugAdapter_1 = require("./debugAdapter");
let dapServerProcess = null;
let outputChannel;
function activate(context) {
    console.log('GitHub Actions DSL extension is now active');
    // Create output channel for logging
    outputChannel = vscode.window.createOutputChannel('GitHub Actions DSL');
    context.subscriptions.push(outputChannel);
    // Register debug adapter factory
    const factory = new debugAdapter_1.DslDebugAdapterDescriptorFactory();
    context.subscriptions.push(vscode.debug.registerDebugAdapterDescriptorFactory('dsl', factory));
    // Register commands
    registerCommands(context);
    // Register language features
    registerLanguageFeatures(context);
    // Auto-start DAP server if configured
    const config = vscode.workspace.getConfiguration('github-actions-dsl');
    if (config.get('dapServer.autoStart')) {
        startDapServer();
    }
    outputChannel.appendLine('GitHub Actions DSL extension activated successfully');
}
exports.activate = activate;
function registerCommands(context) {
    // Start DAP Server command
    const startDapServerCommand = vscode.commands.registerCommand('github-actions-dsl.startDapServer', startDapServer);
    // Stop DAP Server command
    const stopDapServerCommand = vscode.commands.registerCommand('github-actions-dsl.stopDapServer', stopDapServer);
    // Restart DAP Server command
    const restartDapServerCommand = vscode.commands.registerCommand('github-actions-dsl.restartDapServer', async () => {
        await stopDapServer();
        await new Promise(resolve => setTimeout(resolve, 1000));
        await startDapServer();
    });
    // Open Debug Console command
    const openDebugConsoleCommand = vscode.commands.registerCommand('github-actions-dsl.openDebugConsole', () => {
        vscode.commands.executeCommand('workbench.debug.action.toggleRepl');
    });
    // Validate Syntax command
    const validateSyntaxCommand = vscode.commands.registerCommand('github-actions-dsl.validateSyntax', validateSyntax);
    // Show AST command
    const showAstCommand = vscode.commands.registerCommand('github-actions-dsl.showAst', showAst);
    // Get program name for debugger
    const getProgramNameCommand = vscode.commands.registerCommand('github-actions-dsl.getProgramName', async () => {
        const editor = vscode.window.activeTextEditor;
        if (editor && editor.document.languageId === 'github-actions-dsl') {
            return editor.document.uri.fsPath;
        }
        const files = await vscode.workspace.findFiles('**/*.dsl', '**/node_modules/**');
        if (files.length > 0) {
            return files[0].fsPath;
        }
        return vscode.window.showInputBox({
            placeHolder: 'Enter path to DSL file',
            value: 'examples/simple.dsl'
        });
    });
    context.subscriptions.push(startDapServerCommand, stopDapServerCommand, restartDapServerCommand, openDebugConsoleCommand, validateSyntaxCommand, showAstCommand, getProgramNameCommand);
}
function registerLanguageFeatures(context) {
    // Register completion provider
    const completionProvider = vscode.languages.registerCompletionItemProvider('github-actions-dsl', new DslCompletionProvider(), ':', '-', ' ');
    // Register hover provider
    const hoverProvider = vscode.languages.registerHoverProvider('github-actions-dsl', new DslHoverProvider());
    // Register diagnostic collection
    const diagnosticCollection = vscode.languages.createDiagnosticCollection('github-actions-dsl');
    context.subscriptions.push(diagnosticCollection);
    // Register document validation
    const validateDocument = (document) => {
        if (document.languageId === 'github-actions-dsl') {
            validateDocumentSyntax(document, diagnosticCollection);
        }
    };
    vscode.workspace.onDidOpenTextDocument(validateDocument);
    vscode.workspace.onDidChangeTextDocument(event => validateDocument(event.document));
    context.subscriptions.push(completionProvider, hoverProvider, diagnosticCollection);
}
async function startDapServer() {
    if (dapServerProcess) {
        outputChannel.appendLine('DAP server is already running');
        return;
    }
    const config = vscode.workspace.getConfiguration('github-actions-dsl');
    const executable = config.get('dapServer.executable', './bin/dap');
    const port = config.get('dapServer.port', 4711);
    const verbose = config.get('debug.verbose', false);
    // Find the executable path
    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        vscode.window.showErrorMessage('No workspace folder found');
        return;
    }
    const executablePath = path.resolve(workspaceFolders[0].uri.fsPath, executable);
    // Check if executable exists
    if (!fs.existsSync(executablePath)) {
        const buildFirst = await vscode.window.showWarningMessage(`DAP server executable not found at ${executablePath}. Would you like to build it first?`, 'Build', 'Cancel');
        if (buildFirst === 'Build') {
            await buildDapServer();
        }
        else {
            return;
        }
    }
    try {
        const args = ['--dap-server', '--dap-port', port.toString()];
        if (verbose) {
            args.push('--verbose');
        }
        outputChannel.appendLine(`Starting DAP server: ${executablePath} ${args.join(' ')}`);
        dapServerProcess = (0, child_process_1.spawn)(executablePath, args, {
            cwd: workspaceFolders[0].uri.fsPath
        });
        dapServerProcess.stdout?.on('data', (data) => {
            outputChannel.appendLine(`DAP Server: ${data.toString().trim()}`);
        });
        dapServerProcess.stderr?.on('data', (data) => {
            outputChannel.appendLine(`DAP Server Error: ${data.toString().trim()}`);
        });
        dapServerProcess.on('close', (code) => {
            outputChannel.appendLine(`DAP server exited with code ${code}`);
            dapServerProcess = null;
        });
        outputChannel.appendLine(`DAP server started on port ${port}`);
        vscode.window.showInformationMessage(`DAP server started on port ${port}`);
    }
    catch (error) {
        outputChannel.appendLine(`Error starting DAP server: ${error}`);
        vscode.window.showErrorMessage(`Failed to start DAP server: ${error}`);
    }
}
async function stopDapServer() {
    if (dapServerProcess) {
        outputChannel.appendLine('Stopping DAP server...');
        dapServerProcess.kill();
        dapServerProcess = null;
        vscode.window.showInformationMessage('DAP server stopped');
    }
}
async function buildDapServer() {
    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        return;
    }
    return new Promise((resolve, reject) => {
        outputChannel.appendLine('Building DAP server...');
        (0, child_process_1.exec)('go build -o bin/dap ./cmd/dap', {
            cwd: workspaceFolders[0].uri.fsPath
        }, (error, stdout, stderr) => {
            if (error) {
                outputChannel.appendLine(`Build error: ${error.message}`);
                vscode.window.showErrorMessage(`Failed to build DAP server: ${error.message}`);
                reject(error);
            }
            else {
                outputChannel.appendLine('DAP server built successfully');
                if (stdout)
                    outputChannel.appendLine(`Build output: ${stdout}`);
                if (stderr)
                    outputChannel.appendLine(`Build stderr: ${stderr}`);
                resolve();
            }
        });
    });
}
async function validateSyntax() {
    const editor = vscode.window.activeTextEditor;
    if (!editor || editor.document.languageId !== 'github-actions-dsl') {
        vscode.window.showWarningMessage('Please open a DSL file first');
        return;
    }
    const document = editor.document;
    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        return;
    }
    // Save document if modified
    if (document.isDirty) {
        await document.save();
    }
    const filePath = document.uri.fsPath;
    const mainExecutable = path.resolve(workspaceFolders[0].uri.fsPath, 'bin/golang-vibe-coding');
    (0, child_process_1.exec)(`"${mainExecutable}" compile "${filePath}" --validate`, {
        cwd: workspaceFolders[0].uri.fsPath
    }, (error, stdout, stderr) => {
        if (error) {
            vscode.window.showErrorMessage(`Validation failed: ${error.message}`);
            outputChannel.appendLine(`Validation error: ${error.message}`);
        }
        else {
            vscode.window.showInformationMessage('DSL syntax is valid');
            outputChannel.appendLine('DSL validation passed');
        }
        if (stdout)
            outputChannel.appendLine(`Validation output: ${stdout}`);
        if (stderr)
            outputChannel.appendLine(`Validation stderr: ${stderr}`);
    });
}
async function showAst() {
    const editor = vscode.window.activeTextEditor;
    if (!editor || editor.document.languageId !== 'github-actions-dsl') {
        vscode.window.showWarningMessage('Please open a DSL file first');
        return;
    }
    // Create and show AST visualization panel
    const panel = vscode.window.createWebviewPanel('dslAst', 'DSL AST Visualization', vscode.ViewColumn.Beside, {
        enableScripts: true,
        retainContextWhenHidden: true
    });
    panel.webview.html = getAstVisualizationHtml();
    // Send AST data to webview (would need to call our Go parser)
    // For now, show a placeholder
    panel.webview.postMessage({
        command: 'showAst',
        ast: { type: 'placeholder', message: 'AST visualization coming soon' }
    });
}
function getAstVisualizationHtml() {
    return `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>DSL AST Visualization</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                .ast-node { margin: 10px 0; padding: 5px; border-left: 2px solid #007ACC; }
                .ast-children { margin-left: 20px; }
            </style>
        </head>
        <body>
            <h1>DSL AST Visualization</h1>
            <div id="ast-container">Loading...</div>
            <script>
                window.addEventListener('message', event => {
                    const message = event.data;
                    if (message.command === 'showAst') {
                        document.getElementById('ast-container').innerHTML = 
                            '<pre>' + JSON.stringify(message.ast, null, 2) + '</pre>';
                    }
                });
            </script>
        </body>
        </html>
    `;
}
function validateDocumentSyntax(document, diagnostics) {
    const config = vscode.workspace.getConfiguration('github-actions-dsl');
    if (!config.get('syntax.validation.enabled')) {
        return;
    }
    // Basic syntax validation - this would be enhanced with actual parser integration
    const text = document.getText();
    const lines = text.split('\n');
    const diagnosticItems = [];
    lines.forEach((line, index) => {
        // Example validation rules
        if (line.trim().startsWith('name:') && line.indexOf(':') === line.lastIndexOf(':')) {
            // Check for missing value after name:
            const nameValue = line.substring(line.indexOf(':') + 1).trim();
            if (!nameValue) {
                const diagnostic = new vscode.Diagnostic(new vscode.Range(index, line.indexOf(':') + 1, index, line.length), 'Name field cannot be empty', vscode.DiagnosticSeverity.Error);
                diagnosticItems.push(diagnostic);
            }
        }
    });
    diagnostics.set(document.uri, diagnosticItems);
}
// Completion provider for DSL constructs
class DslCompletionProvider {
    provideCompletionItems(document, position) {
        const items = [];
        // Add common DSL keywords
        items.push(new vscode.CompletionItem('name', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('on', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('jobs', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('steps', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('runs-on', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('uses', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('with', vscode.CompletionItemKind.Keyword), new vscode.CompletionItem('env', vscode.CompletionItemKind.Keyword));
        return items;
    }
}
// Hover provider for DSL constructs
class DslHoverProvider {
    provideHover(document, position) {
        const range = document.getWordRangeAtPosition(position);
        const word = document.getText(range);
        const hoverInfo = {
            'name': 'The name of the workflow',
            'on': 'The event that triggers the workflow',
            'jobs': 'A workflow run is made up of one or more jobs',
            'steps': 'A job contains a sequence of tasks called steps',
            'runs-on': 'The type of machine to run the job on',
            'uses': 'Selects an action to run as part of a step',
            'with': 'A map of the input parameters defined by the action'
        };
        if (hoverInfo[word]) {
            return new vscode.Hover(hoverInfo[word]);
        }
        return undefined;
    }
}
function deactivate() {
    if (dapServerProcess) {
        dapServerProcess.kill();
    }
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map