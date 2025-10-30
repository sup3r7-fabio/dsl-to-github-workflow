// Package debug - Debug Adapter Protocol (DAP) implementation for VS Code integration
package debug

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

// DAP Protocol Messages - based on Microsoft Debug Adapter Protocol specification

// DAPMessage represents the base structure for all DAP messages
type DAPMessage struct {
	Seq  int    `json:"seq"`
	Type string `json:"type"` // "request", "response", "event"
}

// DAPRequest represents a DAP request message
type DAPRequest struct {
	DAPMessage
	Command   string      `json:"command"`
	Arguments interface{} `json:"arguments,omitempty"`
}

// DAPResponse represents a DAP response message
type DAPResponse struct {
	DAPMessage
	RequestSeq int         `json:"request_seq"`
	Success    bool        `json:"success"`
	Command    string      `json:"command"`
	Message    string      `json:"message,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

// DAPEvent represents a DAP event message
type DAPEvent struct {
	DAPMessage
	Event string      `json:"event"`
	Body  interface{} `json:"body,omitempty"`
}

// DAP Request Arguments
type InitializeRequestArguments struct {
	ClientID                     string `json:"clientID,omitempty"`
	ClientName                   string `json:"clientName,omitempty"`
	AdapterID                    string `json:"adapterID"`
	Locale                       string `json:"locale,omitempty"`
	LinesStartAt1                bool   `json:"linesStartAt1,omitempty"`
	ColumnsStartAt1              bool   `json:"columnsStartAt1,omitempty"`
	PathFormat                   string `json:"pathFormat,omitempty"`
	SupportsVariableType         bool   `json:"supportsVariableType,omitempty"`
	SupportsVariablePaging       bool   `json:"supportsVariablePaging,omitempty"`
	SupportsRunInTerminalRequest bool   `json:"supportsRunInTerminalRequest,omitempty"`
}

type LaunchRequestArguments struct {
	NoDebug bool     `json:"noDebug,omitempty"`
	Program string   `json:"program"`
	Args    []string `json:"args,omitempty"`
}

type SetBreakpointsArguments struct {
	Source      Source             `json:"source"`
	Breakpoints []SourceBreakpoint `json:"breakpoints,omitempty"`
}

type Source struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type SourceBreakpoint struct {
	Line      int    `json:"line"`
	Column    int    `json:"column,omitempty"`
	Condition string `json:"condition,omitempty"`
}

// DAP Response Bodies
type Capabilities struct {
	SupportsConfigurationDoneRequest  bool `json:"supportsConfigurationDoneRequest,omitempty"`
	SupportsFunctionBreakpoints       bool `json:"supportsFunctionBreakpoints,omitempty"`
	SupportsConditionalBreakpoints    bool `json:"supportsConditionalBreakpoints,omitempty"`
	SupportsHitConditionalBreakpoints bool `json:"supportsHitConditionalBreakpoints,omitempty"`
	SupportsEvaluateForHovers         bool `json:"supportsEvaluateForHovers,omitempty"`
	SupportsStepBack                  bool `json:"supportsStepBack,omitempty"`
	SupportsSetVariable               bool `json:"supportsSetVariable,omitempty"`
	SupportsRestartFrame              bool `json:"supportsRestartFrame,omitempty"`
	SupportsGotoTargetsRequest        bool `json:"supportsGotoTargetsRequest,omitempty"`
	SupportsStepInTargetsRequest      bool `json:"supportsStepInTargetsRequest,omitempty"`
	SupportsCompletionsRequest        bool `json:"supportsCompletionsRequest,omitempty"`
}

type SetBreakpointsResponseBody struct {
	Breakpoints []DAPBreakpoint `json:"breakpoints"`
}

type DAPBreakpoint struct {
	ID       int    `json:"id,omitempty"`
	Verified bool   `json:"verified"`
	Message  string `json:"message,omitempty"`
	Source   Source `json:"source,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
}

type StackTraceResponseBody struct {
	StackFrames []StackFrame `json:"stackFrames"`
	TotalFrames int          `json:"totalFrames,omitempty"`
}

type StackFrame struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Source Source `json:"source,omitempty"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type ScopesResponseBody struct {
	Scopes []Scope `json:"scopes"`
}

type Scope struct {
	Name               string `json:"name"`
	VariablesReference int    `json:"variablesReference"`
	Expensive          bool   `json:"expensive,omitempty"`
}

type VariablesResponseBody struct {
	Variables []Variable `json:"variables"`
}

type Variable struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	Type               string `json:"type,omitempty"`
	VariablesReference int    `json:"variablesReference,omitempty"`
}

// DAP Event Bodies
type InitializedEventBody struct{}

type StoppedEventBody struct {
	Reason            string `json:"reason"`
	Description       string `json:"description,omitempty"`
	ThreadID          int    `json:"threadId,omitempty"`
	PreserveFocusHint bool   `json:"preserveFocusHint,omitempty"`
	Text              string `json:"text,omitempty"`
	AllThreadsStopped bool   `json:"allThreadsStopped,omitempty"`
}

type ExitedEventBody struct {
	ExitCode int `json:"exitCode"`
}

// DAPAdapter represents the Debug Adapter Protocol server
type DAPAdapter struct {
	conn             net.Conn
	reader           *bufio.Reader
	writer           io.Writer
	seq              int
	mutex            sync.Mutex
	debugger         *Debugger
	running          bool
	breakpoints      map[string][]DAPBreakpoint
	nextBreakpointID int
	variableRefs     map[int]interface{}
	nextVarRef       int
}

// NewDAPAdapter creates a new DAP adapter
func NewDAPAdapter(conn net.Conn) *DAPAdapter {
	return &DAPAdapter{
		conn:         conn,
		reader:       bufio.NewReader(conn),
		writer:       conn,
		breakpoints:  make(map[string][]DAPBreakpoint),
		variableRefs: make(map[int]interface{}),
		nextVarRef:   1000,
	}
}

// Start starts the DAP adapter
func (dap *DAPAdapter) Start() {
	defer dap.conn.Close()
	dap.running = true

	for dap.running {
		message, err := dap.readMessage()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("DAP: Error reading message: %v\n", err)
			}
			break
		}

		dap.handleMessage(message)
	}
}

// readMessage reads a DAP message from the connection
func (dap *DAPAdapter) readMessage() ([]byte, error) {
	// Read Content-Length header
	contentLength := 0
	for {
		line, err := dap.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}

		if strings.HasPrefix(line, "Content-Length: ") {
			lengthStr := strings.TrimPrefix(line, "Content-Length: ")
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid content length: %s", lengthStr)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("no content length specified")
	}

	// Read the message content
	content := make([]byte, contentLength)
	_, err := io.ReadFull(dap.reader, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// writeMessage writes a DAP message to the connection
func (dap *DAPAdapter) writeMessage(data []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	dap.mutex.Lock()
	defer dap.mutex.Unlock()

	_, err := dap.writer.Write([]byte(header))
	if err != nil {
		return err
	}

	_, err = dap.writer.Write(data)
	return err
}

// handleMessage processes a DAP message
func (dap *DAPAdapter) handleMessage(messageData []byte) {
	var message map[string]interface{}
	if err := json.Unmarshal(messageData, &message); err != nil {
		fmt.Printf("DAP: Error unmarshaling message: %v\n", err)
		return
	}

	msgType, ok := message["type"].(string)
	if !ok {
		fmt.Printf("DAP: Invalid message type\n")
		return
	}

	switch msgType {
	case "request":
		dap.handleRequest(message)
	default:
		fmt.Printf("DAP: Unknown message type: %s\n", msgType)
	}
}

// handleRequest processes DAP requests
func (dap *DAPAdapter) handleRequest(message map[string]interface{}) {
	seq, _ := message["seq"].(float64)
	command, _ := message["command"].(string)
	arguments := message["arguments"]

	switch command {
	case "initialize":
		dap.handleInitialize(int(seq))
	case "launch":
		dap.handleLaunch(int(seq))
	case "setBreakpoints":
		dap.handleSetBreakpoints(int(seq), arguments)
	case "configurationDone":
		dap.handleConfigurationDone(int(seq))
	case "continue":
		dap.handleContinueRequest(int(seq), arguments)
	case "next":
		dap.handleNext(int(seq), arguments)
	case "stepIn":
		dap.handleStepIn(int(seq), arguments)
	case "stepOut":
		dap.handleStepOut(int(seq), arguments)
	case "stackTrace":
		dap.handleStackTrace(int(seq), arguments)
	case "scopes":
		dap.handleScopes(int(seq), arguments)
	case "variables":
		dap.handleVariables(int(seq), arguments)
	case "evaluate":
		dap.handleEvaluate(int(seq), arguments)
	case "disconnect":
		dap.handleDisconnect(int(seq))
	default:
		dap.sendErrorResponse(int(seq), command, fmt.Sprintf("Unknown command: %s", command))
	}
}

// DAP Request Handlers

func (dap *DAPAdapter) handleInitialize(seq int) {
	capabilities := Capabilities{
		SupportsConfigurationDoneRequest:  true,
		SupportsFunctionBreakpoints:       false,
		SupportsConditionalBreakpoints:    true,
		SupportsHitConditionalBreakpoints: false,
		SupportsEvaluateForHovers:         true,
		SupportsStepBack:                  false,
		SupportsSetVariable:               false,
		SupportsRestartFrame:              false,
		SupportsGotoTargetsRequest:        false,
		SupportsStepInTargetsRequest:      false,
		SupportsCompletionsRequest:        false,
	}

	dap.sendResponse(seq, "initialize", true, "", capabilities)
	dap.sendEvent("initialized", InitializedEventBody{})
}

func (dap *DAPAdapter) handleLaunch(seq int) {
	// Initialize debugger if not already done
	if dap.debugger == nil {
		dap.debugger = GetGlobalDebugger()
		if dap.debugger == nil {
			Init()
			dap.debugger = GetGlobalDebugger()
		}
	}

	dap.sendResponse(seq, "launch", true, "", nil)
}

func (dap *DAPAdapter) handleSetBreakpoints(seq int, args interface{}) {
	var arguments SetBreakpointsArguments
	if argsBytes, err := json.Marshal(args); err == nil {
		json.Unmarshal(argsBytes, &arguments)
	}

	sourcePath := arguments.Source.Path
	breakpoints := []DAPBreakpoint{}

	// Clear existing breakpoints for this source
	delete(dap.breakpoints, sourcePath)

	// Set new breakpoints
	for _, bp := range arguments.Breakpoints {
		dap.nextBreakpointID++

		dapBp := DAPBreakpoint{
			ID:       dap.nextBreakpointID,
			Verified: true,
			Line:     bp.Line,
			Source:   arguments.Source,
		}

		// Set breakpoint in debugger
		if dap.debugger != nil {
			SetBreakpoint(sourcePath, bp.Line, "")
		}

		breakpoints = append(breakpoints, dapBp)
	}

	dap.breakpoints[sourcePath] = breakpoints

	body := SetBreakpointsResponseBody{
		Breakpoints: breakpoints,
	}

	dap.sendResponse(seq, "setBreakpoints", true, "", body)
}

func (dap *DAPAdapter) handleConfigurationDone(seq int) {
	dap.sendResponse(seq, "configurationDone", true, "", nil)
}

func (dap *DAPAdapter) handleContinueRequest(seq int, _ interface{}) {
	if dap.debugger != nil {
		dap.debugger.continueMode = true
		dap.debugger.stepMode = false
	}

	dap.sendResponse(seq, "continue", true, "", map[string]bool{
		"allThreadsContinued": true,
	})
}

func (dap *DAPAdapter) handleNext(seq int, _ interface{}) {
	if dap.debugger != nil {
		dap.debugger.stepMode = true
	}

	dap.sendResponse(seq, "next", true, "", nil)
}

func (dap *DAPAdapter) handleStepIn(seq int, _ interface{}) {
	if dap.debugger != nil {
		dap.debugger.stepMode = true
	}

	dap.sendResponse(seq, "stepIn", true, "", nil)
}

func (dap *DAPAdapter) handleStepOut(seq int, _ interface{}) {
	if dap.debugger != nil {
		dap.debugger.stepMode = true
	}

	dap.sendResponse(seq, "stepOut", true, "", nil)
}

func (dap *DAPAdapter) handleStackTrace(seq int, _ interface{}) {
	frames := []StackFrame{}

	if dap.debugger != nil && dap.debugger.currentContext != nil {
		for i, frame := range dap.debugger.currentContext.CallStack {
			frames = append(frames, StackFrame{
				ID:   i + 1,
				Name: frame,
				Line: 1, // Would need actual line tracking
			})
		}
	}

	body := StackTraceResponseBody{
		StackFrames: frames,
		TotalFrames: len(frames),
	}

	dap.sendResponse(seq, "stackTrace", true, "", body)
}

func (dap *DAPAdapter) handleScopes(seq int, _ interface{}) {
	scopes := []Scope{
		{
			Name:               "Local",
			VariablesReference: dap.getVariableReference("local"),
		},
	}

	body := ScopesResponseBody{
		Scopes: scopes,
	}

	dap.sendResponse(seq, "scopes", true, "", body)
}

func (dap *DAPAdapter) handleVariables(seq int, _ interface{}) {
	variables := []Variable{}

	if dap.debugger != nil && dap.debugger.currentContext != nil {
		for name, value := range dap.debugger.currentContext.Variables {
			variables = append(variables, Variable{
				Name:  name,
				Value: fmt.Sprintf("%v", value),
				Type:  fmt.Sprintf("%T", value),
			})
		}
	}

	body := VariablesResponseBody{
		Variables: variables,
	}

	dap.sendResponse(seq, "variables", true, "", body)
}

func (dap *DAPAdapter) handleEvaluate(seq int, _ interface{}) {
	// Simple evaluation - in a full implementation this would evaluate expressions
	result := "Evaluation not yet implemented"

	body := map[string]interface{}{
		"result":             result,
		"variablesReference": 0,
	}

	dap.sendResponse(seq, "evaluate", true, "", body)
}

func (dap *DAPAdapter) handleDisconnect(seq int) {
	dap.sendResponse(seq, "disconnect", true, "", nil)
	dap.running = false
}

// Utility methods

func (dap *DAPAdapter) sendResponse(requestSeq int, command string, success bool, message string, body interface{}) {
	dap.seq++
	response := DAPResponse{
		DAPMessage: DAPMessage{
			Seq:  dap.seq,
			Type: "response",
		},
		RequestSeq: requestSeq,
		Success:    success,
		Command:    command,
		Message:    message,
		Body:       body,
	}

	data, _ := json.Marshal(response)
	dap.writeMessage(data)
}

func (dap *DAPAdapter) sendErrorResponse(requestSeq int, command string, message string) {
	dap.sendResponse(requestSeq, command, false, message, nil)
}

func (dap *DAPAdapter) sendEvent(event string, body interface{}) {
	dap.seq++
	eventMsg := DAPEvent{
		DAPMessage: DAPMessage{
			Seq:  dap.seq,
			Type: "event",
		},
		Event: event,
		Body:  body,
	}

	data, _ := json.Marshal(eventMsg)
	dap.writeMessage(data)
}

func (dap *DAPAdapter) getVariableReference(scope string) int {
	dap.nextVarRef++
	dap.variableRefs[dap.nextVarRef] = scope
	return dap.nextVarRef
}

// SendStoppedEvent sends a stopped event to the client
func (dap *DAPAdapter) SendStoppedEvent(reason string, threadID int) {
	body := StoppedEventBody{
		Reason:            reason,
		ThreadID:          threadID,
		AllThreadsStopped: true,
	}
	dap.sendEvent("stopped", body)
}

// SendExitedEvent sends an exited event to the client
func (dap *DAPAdapter) SendExitedEvent(exitCode int) {
	body := ExitedEventBody{
		ExitCode: exitCode,
	}
	dap.sendEvent("exited", body)
}

// StartDAPServer starts a DAP server on the specified port
func StartDAPServer(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start DAP server: %w", err)
	}

	fmt.Printf("🔌 DAP Server started on port %d\n", port)
	fmt.Printf("   VS Code can connect to: localhost:%d\n", port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("DAP: Error accepting connection: %v\n", err)
				continue
			}

			fmt.Println("🔗 DAP Client connected")

			adapter := NewDAPAdapter(conn)
			go adapter.Start()
		}
	}()

	return nil
}
