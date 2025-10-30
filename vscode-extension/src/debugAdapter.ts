import * as vscode from 'vscode';
import * as net from 'net';

export class DslDebugAdapterDescriptorFactory implements vscode.DebugAdapterDescriptorFactory {
    createDebugAdapterDescriptor(
        session: vscode.DebugSession,
        executable: vscode.DebugAdapterExecutable | undefined
    ): vscode.ProviderResult<vscode.DebugAdapterDescriptor> {
        
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

export class DslDebugAdapter {
    private socket: net.Socket | null = null;

    constructor(private port: number, private host: string = 'localhost') {}

    public async connect(): Promise<void> {
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

    public disconnect(): void {
        if (this.socket) {
            this.socket.destroy();
            this.socket = null;
        }
    }

    public send(message: any): void {
        if (this.socket) {
            const messageStr = JSON.stringify(message);
            const contentLength = Buffer.byteLength(messageStr, 'utf8');
            const header = `Content-Length: ${contentLength}\r\n\r\n`;
            this.socket.write(header + messageStr);
        }
    }

    public onMessage(callback: (message: any) => void): void {
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
                    } catch (error) {
                        console.error('Error parsing DAP message:', error);
                    }
                }
            });
        }
    }
}
