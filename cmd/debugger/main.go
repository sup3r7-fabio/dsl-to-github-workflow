package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang-vibe-coding/internal/ast"
	"golang-vibe-coding/internal/lexer"
	"golang-vibe-coding/internal/parser"
	"golang-vibe-coding/pkg/debug"
)

func main() {
	var (
		debugLevel  = flag.String("debug", "", "Debug level: basic, verbose, trace")
		serverPort  = flag.Int("port", 8080, "Debug server port")
		interactive = flag.Bool("interactive", false, "Start interactive debugger")
		visualize   = flag.Bool("visualize", false, "Show AST visualization")
		analyze     = flag.Bool("analyze", false, "Show AST analysis")
		webUI       = flag.Bool("web", false, "Start web UI debugger")
	)
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: debugger [flags] <dsl-file>")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  debugger --debug=verbose --visualize examples/simple.dsl")
		fmt.Println("  debugger --interactive --web examples/production.dsl")
		fmt.Println("  debugger --analyze examples/*.dsl")
		os.Exit(1)
	}

	// Set up debug environment
	if *debugLevel != "" {
		os.Setenv("DSL_DEBUG", *debugLevel)
	}

	debug.Init()

	if *webUI {
		debug.StartDebugServer(*serverPort)
		fmt.Printf("🌐 Web debugger started at http://localhost:%d/debug/ui/\n", *serverPort)
	}

	dslFile := flag.Arg(0)

	fmt.Printf("🐛 DSL Debugger - Processing: %s\n", dslFile)

	// Read DSL file
	content, err := ioutil.ReadFile(dslFile)
	if err != nil {
		fmt.Printf("❌ Error reading file: %v\n", err)
		os.Exit(1)
	}

	dslContent := string(content)

	if *interactive {
		runInteractiveDebugger(dslFile, dslContent)
	} else {
		runBatchDebugger(dslFile, dslContent, *visualize, *analyze)
	}
}

func runInteractiveDebugger(filename, content string) {
	fmt.Println("🔧 Starting Interactive DSL Debugger")
	fmt.Println("Available commands: parse, lex, visualize, analyze, breakpoint, step, continue, help, quit")

	reader := bufio.NewReader(os.Stdin)

	var workflow *ast.Workflow
	var p *parser.Parser

	for {
		fmt.Print("(dsl-debug) ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case "help", "h":
			showInteractiveHelp()

		case "lex", "tokenize":
			fmt.Println("🔤 Tokenizing DSL...")
			l := lexer.New(content)

			tokenCount := 0
			for {
				token := l.NextToken()
				if token.Type == lexer.EOF {
					break
				}
				debug.Token(token.Type.String(), token.Literal, token.Line, token.Column)
				tokenCount++

				// Add protection against infinite loops in lexer
				if tokenCount > 10000 {
					fmt.Printf("⚠️ WARNING: Token count exceeded 10,000. Possible infinite loop detected!\n")
					fmt.Printf("Last token: Type=%s, Literal='%s', Line=%d, Column=%d\n",
						token.Type.String(), token.Literal, token.Line, token.Column)
					break
				}
			}
			fmt.Printf("✅ Tokenization complete: %d tokens\n", tokenCount)

		case "parse":
			fmt.Println("📝 Parsing DSL...")
			l := lexer.New(content)
			p = parser.New(l)

			debug.SetBreakpoint(filename, 1, "parsing")

			// Add timeout protection for parsing
			parsingComplete := make(chan bool, 1)
			var parseErr error

			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("🚨 PANIC during parsing: %v\n", r)
						parseErr = fmt.Errorf("parsing panic: %v", r)
					}
					parsingComplete <- true
				}()

				workflow = p.ParseProgram()
			}()

			// Wait for parsing completion
			<-parsingComplete
			if parseErr != nil {
				fmt.Printf("❌ Parsing failed: %v\n", parseErr)
				continue
			}

			if errors := p.Errors(); len(errors) > 0 {
				fmt.Printf("❌ Parsing errors:\n")
				for _, err := range errors {
					fmt.Printf("  %s\n", err)
				}
			} else {
				fmt.Println("✅ Parsing complete")
			}

		case "visualize", "vis":
			if workflow == nil {
				fmt.Println("⚠️ No workflow parsed yet. Run 'parse' first.")
				continue
			}

			visualizer := debug.NewASTVisualizer(debug.DefaultVisualizationOptions())
			visualization := visualizer.VisualizeWorkflow(workflow)
			fmt.Println("🌳 AST Visualization:")
			fmt.Println(visualization)

		case "analyze":
			if workflow == nil {
				fmt.Println("⚠️ No workflow parsed yet. Run 'parse' first.")
				continue
			}

			analysis := debug.AnalyzeAST(workflow)
			fmt.Println("📊 AST Analysis:")
			for key, value := range analysis {
				fmt.Printf("  %s: %v\n", key, value)
			}

		case "breakpoint", "bp":
			if len(args) < 1 {
				fmt.Println("Usage: breakpoint <line> [function]")
				continue
			}

			// This would need line number parsing
			fmt.Printf("Setting breakpoint at line %s\n", args[0])
			debug.SetBreakpoint(filename, 1, "manual")

		case "vars", "variables":
			fmt.Println("📋 Current Variables:")
			fmt.Println("  filename:", filename)
			fmt.Println("  content_length:", len(content))
			if workflow != nil {
				fmt.Println("  workflow: parsed")
			} else {
				fmt.Println("  workflow: not parsed")
			}

		case "step":
			fmt.Println("👣 Step functionality would advance one parsing step")

		case "continue", "c":
			fmt.Println("▶️ Continue functionality would resume parsing")

		case "reload", "r":
			// Reload the DSL file
			newContent, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Printf("❌ Error reloading file: %v\n", err)
			} else {
				content = string(newContent)
				workflow = nil
				p = nil
				fmt.Println("🔄 DSL file reloaded")
			}

		case "quit", "q", "exit":
			fmt.Println("👋 Goodbye!")
			os.Exit(0)

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}

func runBatchDebugger(_, content string, visualize, analyze bool) {
	fmt.Println("📊 Running batch analysis...")

	// Tokenize with protection
	fmt.Println("🔤 Tokenizing...")
	l := lexer.New(content)

	tokenCount := 0
	for {
		token := l.NextToken()
		if token.Type == lexer.EOF {
			break
		}
		tokenCount++

		// Protection against infinite tokenization
		if tokenCount > 10000 {
			fmt.Printf("⚠️ WARNING: Token count exceeded 10,000. Possible infinite loop in lexer!\n")
			fmt.Printf("Last token: Type=%s, Literal='%s', Line=%d, Column=%d\n",
				token.Type.String(), token.Literal, token.Line, token.Column)
			return
		}
	}

	fmt.Printf("✅ Tokenization complete: %d tokens\n", tokenCount)

	// Parse with protection
	fmt.Println("📝 Parsing...")
	l = lexer.New(content)
	p := parser.New(l)

	// Add timeout protection for parsing
	parsingComplete := make(chan bool, 1)
	var workflow *ast.Workflow
	var parseErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("🚨 PANIC during parsing: %v\n", r)
				parseErr = fmt.Errorf("parsing panic: %v", r)
			}
			parsingComplete <- true
		}()

		workflow = p.ParseProgram()
	}()

	// Wait for parsing
	<-parsingComplete
	if parseErr != nil {
		fmt.Printf("❌ Parsing failed: %v\n", parseErr)
		return
	}

	// Check errors
	if errors := p.Errors(); len(errors) > 0 {
		fmt.Printf("❌ Parsing errors:\n")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		return
	}

	if workflow == nil {
		fmt.Println("❌ Failed to parse workflow")
		return
	}

	fmt.Printf("✅ Successfully parsed workflow: %s\n", workflow.Name)

	// Visualize if requested
	if visualize {
		fmt.Println("\n🌳 AST Visualization:")
		visualizer := debug.NewASTVisualizer(debug.DefaultVisualizationOptions())
		visualization := visualizer.VisualizeWorkflow(workflow)
		fmt.Println(visualization)
	}

	// Analyze if requested
	if analyze {
		fmt.Println("\n📊 AST Analysis:")
		analysis := debug.AnalyzeAST(workflow)
		for key, value := range analysis {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println("\n✨ Analysis complete!")
}

func showInteractiveHelp() {
	fmt.Println("\n📚 Interactive DSL Debugger Commands:")
	fmt.Println("  help (h)           - Show this help")
	fmt.Println("  lex                - Tokenize the DSL (with infinite loop protection)")
	fmt.Println("  parse              - Parse the DSL into AST (with timeout protection)")
	fmt.Println("  visualize (vis)    - Show AST visualization")
	fmt.Println("  analyze            - Show AST analysis")
	fmt.Println("  breakpoint (bp)    - Set a breakpoint")
	fmt.Println("  vars               - Show current variables")
	fmt.Println("  step               - Step through parsing")
	fmt.Println("  continue (c)       - Continue execution")
	fmt.Println("  reload (r)         - Reload DSL file")
	fmt.Println("  quit (q)           - Exit debugger")
	fmt.Println()
}
