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
exports.DslDebugAdapter = exports.DslDebugAdapterDescriptorFactory = void 0;
const vscode = __importStar(require("vscode"));
const net = __importStar(require("net"));
class DslDebugAdapterDescriptorFactory {
    createDebugAdapterDescriptor(session, executable) {
        const config = session.configuration;
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) {
            throw new Error('No workspace folder found');
        }
        // Get DAP server configuration
        const dslConfig = vscode.workspace.getConfiguration('github-actions-dsl');
        const port = config.dapServerPort || dslConfig.get('dapServer.port', 4711);
        const host = config.dapServerHost || dslConfig.get('dapServer.host', 'localhost');
        // Return a socket connection to our DAP server
        return new vscode.DebugAdapterServer(port, host);
    }
}
exports.DslDebugAdapterDescriptorFactory = DslDebugAdapterDescriptorFactory;
class DslDebugAdapter {
    constructor(port, host = 'localhost') {
        this.port = port;
        this.host = host;
        this.socket = null;
    }
    async connect() {
        return new Promise((resolve, reject) => {
            this.socket = new net.Socket();
            this.socket.connect(this.port, this.host, () => {
                resolve();
            });
            this.socket.on('error', (error) => {
                reject(error);
            });
        });
    }
    disconnect() {
        if (this.socket) {
            this.socket.destroy();
            this.socket = null;
        }
    }
    send(message) {
        if (this.socket) {
            const messageStr = JSON.stringify(message);
            const contentLength = Buffer.byteLength(messageStr, 'utf8');
            const header = `Content-Length: ${contentLength}\r\n\r\n`;
            this.socket.write(header + messageStr);
        }
    }
    onMessage(callback) {
        if (this.socket) {
            let buffer = '';
            this.socket.on('data', (data) => {
                buffer += data.toString();
                while (buffer.length > 0) {
                    const headerMatch = buffer.match(/^Content-Length: (\d+)\r\n\r\n/);
                    if (!headerMatch) {
                        break;
                    }
                    const contentLength = parseInt(headerMatch[1]);
                    const headerLength = headerMatch[0].length;
                    if (buffer.length < headerLength + contentLength) {
                        break;
                    }
                    const messageStr = buffer.substring(headerLength, headerLength + contentLength);
                    buffer = buffer.substring(headerLength + contentLength);
                    try {
                        const message = JSON.parse(messageStr);
                        callback(message);
                    }
                    catch (error) {
                        console.error('Error parsing DAP message:', error);
                    }
                }
            });
        }
    }
}
exports.DslDebugAdapter = DslDebugAdapter;
//# sourceMappingURL=debugAdapter.js.map