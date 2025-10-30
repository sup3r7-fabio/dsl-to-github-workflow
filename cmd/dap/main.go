package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang-vibe-coding/pkg/debug"
)

func main() {
	var (
		dapServer = flag.Bool("dap-server", false, "Start DAP (Debug Adapter Protocol) server")
		dapPort   = flag.Int("dap-port", 4711, "DAP server port")
		dapHost   = flag.String("dap-host", "localhost", "DAP server host")
		verbose   = flag.Bool("verbose", false, "Enable verbose logging")
		help      = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("🔍 Verbose logging enabled")
	}

	if *dapServer {
		startDAPServer(*dapHost, *dapPort, *verbose)
		return
	}

	// Default behavior - show help if no specific command given
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	fmt.Println("Use --help to see available options")
}

func startDAPServer(host string, port int, verbose bool) {
	fmt.Printf("🚀 Starting DAP server on %s:%d\n", host, port)
	fmt.Println("📡 Waiting for VSCode connection...")

	if verbose {
		fmt.Println("🐛 Debug mode enabled")
		os.Setenv("DSL_DEBUG", "3")
		os.Setenv("DAP_DEBUG", "true")
	}

	// Initialize debug system
	debug.Init()

	// Start DAP server (this blocks and handles connections)
	if err := debug.StartDAPServer(port); err != nil {
		log.Fatalf("❌ Failed to start DAP server: %v", err)
	}
}

func printHelp() {
	fmt.Println("🐛 DSL Debug Adapter Protocol (DAP) Server")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  dsl-dap [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --dap-server        Start DAP server for VSCode integration")
	fmt.Println("  --dap-port <port>   DAP server port (default: 4711)")
	fmt.Println("  --dap-host <host>   DAP server host (default: localhost)")
	fmt.Println("  --verbose           Enable verbose logging")
	fmt.Println("  --help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dsl-dap --dap-server")
	fmt.Println("  dsl-dap --dap-server --dap-port 5000")
	fmt.Println("  dsl-dap --dap-server --verbose")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DSL_DEBUG     Debug level (1-3)")
	fmt.Println("  DAP_DEBUG     Enable DAP debug logging")
	fmt.Println()
	fmt.Println("For VSCode setup instructions, see VSCODE_DEBUG_GUIDE.md")
}
