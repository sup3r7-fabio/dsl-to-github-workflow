// Package debug provides debugging utilities for the GitHub Actions DSL
package debug

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	DebugEnabled = false
	debugLogger  *log.Logger
	debugger     *Debugger
)

// DebugLevel represents different levels of debugging output
type DebugLevel int

const (
	LevelNone DebugLevel = iota
	LevelBasic
	LevelVerbose
	LevelTrace
)

// Breakpoint represents a debugging breakpoint
type Breakpoint struct {
	File     string
	Line     int
	Column   int
	Function string
	Enabled  bool
}

// ExecutionContext represents the current execution state
type ExecutionContext struct {
	CurrentFile     string
	CurrentLine     int
	CurrentColumn   int
	CurrentFunction string
	CallStack       []string
	Variables       map[string]interface{}
}

// Debugger provides interactive debugging capabilities
type Debugger struct {
	breakpoints    map[string]*Breakpoint
	executionStack []ExecutionContext
	currentContext *ExecutionContext
	stepMode       bool
	stepIntoMode   bool
	continueMode   bool
	level          DebugLevel
	mutex          sync.RWMutex
	reader         *bufio.Reader
}

// Init initializes the debug system
func Init() {
	debugMode := os.Getenv("DSL_DEBUG")
	if debugMode != "" {
		DebugEnabled = true
		debugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)

		// Initialize debugger
		debugger = &Debugger{
			breakpoints: make(map[string]*Breakpoint),
			level:       parseDebugLevel(debugMode),
			reader:      bufio.NewReader(os.Stdin),
		}

		fmt.Println("🐛 DSL Debugger initialized")
		fmt.Println("Available commands: help, break, step, continue, vars, stack, quit")
	}
}

// GetGlobalDebugger returns the global debugger instance
func GetGlobalDebugger() *Debugger {
	return debugger
}

// parseDebugLevel converts string debug level to DebugLevel
func parseDebugLevel(level string) DebugLevel {
	switch strings.ToLower(level) {
	case "1", "basic":
		return LevelBasic
	case "2", "verbose":
		return LevelVerbose
	case "3", "trace":
		return LevelTrace
	default:
		return LevelBasic
	}
}

// SetBreakpoint sets a breakpoint at the specified location
func SetBreakpoint(file string, line int, function string) {
	if debugger == nil {
		return
	}

	debugger.mutex.Lock()
	defer debugger.mutex.Unlock()

	key := fmt.Sprintf("%s:%d", file, line)
	debugger.breakpoints[key] = &Breakpoint{
		File:     file,
		Line:     line,
		Function: function,
		Enabled:  true,
	}

	Log("Breakpoint set at %s:%d in %s", file, line, function)
}

// CheckBreakpoint checks if execution should pause at current location
func CheckBreakpoint(file string, line int, function string, variables map[string]interface{}) {
	if debugger == nil {
		return
	}

	debugger.mutex.Lock()
	defer debugger.mutex.Unlock()

	key := fmt.Sprintf("%s:%d", file, line)

	// Update current context
	debugger.currentContext = &ExecutionContext{
		CurrentFile:     file,
		CurrentLine:     line,
		CurrentFunction: function,
		Variables:       variables,
	}

	// Check if we should pause
	shouldPause := false

	if breakpoint, exists := debugger.breakpoints[key]; exists && breakpoint.Enabled {
		shouldPause = true
		fmt.Printf("🔴 Breakpoint hit at %s:%d in %s\n", file, line, function)
	} else if debugger.stepMode {
		shouldPause = true
		fmt.Printf("👣 Step: %s:%d in %s\n", file, line, function)
	}

	if shouldPause {
		debugger.interactiveSession()
	}
}

// interactiveSession starts an interactive debugging session
func (d *Debugger) interactiveSession() {
	d.stepMode = false // Reset step mode

	fmt.Println("\n--- Interactive Debug Session ---")
	fmt.Println("Type 'help' for available commands")

	for {
		fmt.Print("(dsl-debug) ")

		input, err := d.reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if d.handleCommand(input) {
			break // Continue execution
		}
	}
}

// handleCommand processes debug commands and returns true if execution should continue
func (d *Debugger) handleCommand(command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "help", "h":
		d.showHelp()

	case "continue", "c":
		fmt.Println("Continuing execution...")
		return true

	case "step", "s":
		fmt.Println("Stepping to next line...")
		d.stepMode = true
		return true

	case "vars", "v":
		d.showVariables()

	case "stack", "st":
		d.showStack()

	case "break", "b":
		d.handleBreakCommand(args)

	case "list", "l":
		d.listBreakpoints()

	case "enable":
		d.toggleBreakpoint(args)

	case "disable":
		d.toggleBreakpoint(args)

	case "quit", "q":
		fmt.Println("Exiting debugger...")
		os.Exit(0)

	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", cmd)
	}

	return false
}

// showHelp displays available debug commands
func (d *Debugger) showHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  help (h)           - Show this help")
	fmt.Println("  continue (c)       - Continue execution")
	fmt.Println("  step (s)           - Step to next line")
	fmt.Println("  vars (v)           - Show variables")
	fmt.Println("  stack (st)         - Show call stack")
	fmt.Println("  break (b) <line>   - Set breakpoint at line")
	fmt.Println("  list (l)           - List all breakpoints")
	fmt.Println("  enable <id>        - Enable breakpoint")
	fmt.Println("  disable <id>       - Disable breakpoint")
	fmt.Println("  quit (q)           - Exit debugger")
	fmt.Println()
}

// showVariables displays current variables
func (d *Debugger) showVariables() {
	if d.currentContext == nil || len(d.currentContext.Variables) == 0 {
		fmt.Println("No variables in current scope")
		return
	}

	fmt.Println("\nCurrent Variables:")
	for name, value := range d.currentContext.Variables {
		fmt.Printf("  %s = %v\n", name, value)
	}
	fmt.Println()
}

// showStack displays the current call stack
func (d *Debugger) showStack() {
	if d.currentContext == nil {
		fmt.Println("No execution context available")
		return
	}

	fmt.Println("\nCall Stack:")
	fmt.Printf("  -> %s:%d in %s\n",
		d.currentContext.CurrentFile,
		d.currentContext.CurrentLine,
		d.currentContext.CurrentFunction)

	for i, frame := range d.currentContext.CallStack {
		fmt.Printf("  %d. %s\n", i+1, frame)
	}
	fmt.Println()
}

// handleBreakCommand handles breakpoint setting
func (d *Debugger) handleBreakCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: break <line> [file]")
		return
	}

	line, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid line number: %s\n", args[0])
		return
	}

	file := "current"
	if len(args) > 1 {
		file = args[1]
	} else if d.currentContext != nil {
		file = d.currentContext.CurrentFile
	}

	SetBreakpoint(file, line, "")
	fmt.Printf("Breakpoint set at %s:%d\n", file, line)
}

// listBreakpoints shows all breakpoints
func (d *Debugger) listBreakpoints() {
	if len(d.breakpoints) == 0 {
		fmt.Println("No breakpoints set")
		return
	}

	fmt.Println("\nBreakpoints:")
	i := 1
	for _, bp := range d.breakpoints {
		status := "enabled"
		if !bp.Enabled {
			status = "disabled"
		}
		fmt.Printf("  %d. %s:%d (%s)\n", i, bp.File, bp.Line, status)
		i++
	}
	fmt.Println()
}

// toggleBreakpoint enables/disables a breakpoint
func (d *Debugger) toggleBreakpoint(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: enable/disable <breakpoint_id>")
		return
	}

	// For simplicity, this could be enhanced to match by ID
	fmt.Printf("Breakpoint toggle functionality needs breakpoint ID mapping\n")
}

// Log prints a debug message if debugging is enabled
func Log(format string, args ...interface{}) {
	if DebugEnabled && debugLogger != nil {
		debugLogger.Printf(format, args...)
	}
}

// Trace prints a trace message with function entry/exit info
func Trace(funcName string, action string, info ...interface{}) {
	if !DebugEnabled || debugger == nil {
		return
	}

	if debugger.level >= LevelTrace {
		if len(info) > 0 {
			fmt.Printf("[TRACE] %s %s: %v\n", funcName, action, info)
		} else {
			fmt.Printf("[TRACE] %s %s\n", funcName, action)
		}
	}

	// Update call stack if entering/exiting functions
	if debugger.currentContext != nil {
		if action == "enter" {
			debugger.currentContext.CallStack = append(debugger.currentContext.CallStack, funcName)
		} else if action == "exit" && len(debugger.currentContext.CallStack) > 0 {
			debugger.currentContext.CallStack = debugger.currentContext.CallStack[:len(debugger.currentContext.CallStack)-1]
		}
	}
}

// Token prints token information
func Token(tokenType, literal string, line, col int) {
	if DebugEnabled && debugger != nil && debugger.level >= LevelVerbose {
		fmt.Printf("[TOKEN] %s='%s' at %d:%d\n", tokenType, literal, line, col)
	}
}
