// Package debug - Debug server for IDE integration
package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// DebugServer provides HTTP endpoints for IDE integration
type DebugServer struct {
	port       int
	running    bool
	mutex      sync.RWMutex
	sessions   map[string]*DebugSession
	nextSessID int
}

// DebugSession represents an active debug session
type DebugSession struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"` // running, paused, stopped
	Breakpoints map[string]*Breakpoint `json:"breakpoints"`
	Context     *ExecutionContext      `json:"context"`
	Source      string                 `json:"source"`
}

// DebugResponse represents API response structure
type DebugResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewDebugServer creates a new debug server
func NewDebugServer(port int) *DebugServer {
	return &DebugServer{
		port:     port,
		sessions: make(map[string]*DebugSession),
	}
}

// Start starts the debug server
func (ds *DebugServer) Start() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if ds.running {
		return fmt.Errorf("debug server already running")
	}

	// Setup HTTP endpoints
	http.HandleFunc("/debug/session/create", ds.handleCreateSession)
	http.HandleFunc("/debug/session/", ds.handleSession)
	http.HandleFunc("/debug/breakpoints", ds.handleBreakpoints)
	http.HandleFunc("/debug/step", ds.handleStep)
	http.HandleFunc("/debug/continue", ds.handleContinue)
	http.HandleFunc("/debug/variables", ds.handleVariables)
	http.HandleFunc("/debug/stack", ds.handleStack)
	http.HandleFunc("/debug/evaluate", ds.handleEvaluate)

	// Static file server for debug UI
	http.Handle("/debug/ui/", http.StripPrefix("/debug/ui/",
		http.FileServer(http.Dir("debug/ui/"))))

	ds.running = true

	fmt.Printf("🌐 Debug server starting on port %d\n", ds.port)
	fmt.Printf("   Debug UI: http://localhost:%d/debug/ui/\n", ds.port)
	fmt.Printf("   API Base: http://localhost:%d/debug/\n", ds.port)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", ds.port), nil); err != nil {
			fmt.Printf("Debug server error: %v\n", err)
		}
	}()

	return nil
}

// handleCreateSession creates a new debug session
func (ds *DebugServer) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ds.mutex.Lock()
	ds.nextSessID++
	sessionID := fmt.Sprintf("session_%d", ds.nextSessID)

	session := &DebugSession{
		ID:          sessionID,
		Status:      "stopped",
		Breakpoints: make(map[string]*Breakpoint),
	}

	ds.sessions[sessionID] = session
	ds.mutex.Unlock()

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Message: "Session created",
		Data:    session,
	})
}

// handleSession handles session-specific operations
func (ds *DebugServer) handleSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	path := r.URL.Path[len("/debug/session/"):]
	parts := splitPath(path)

	if len(parts) < 1 {
		ds.sendError(w, "Session ID required", http.StatusBadRequest)
		return
	}

	sessionID := parts[0]

	ds.mutex.RLock()
	session, exists := ds.sessions[sessionID]
	ds.mutex.RUnlock()

	if !exists {
		ds.sendError(w, "Session not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ds.sendResponse(w, DebugResponse{
			Success: true,
			Data:    session,
		})
	case http.MethodDelete:
		ds.mutex.Lock()
		delete(ds.sessions, sessionID)
		ds.mutex.Unlock()

		ds.sendResponse(w, DebugResponse{
			Success: true,
			Message: "Session deleted",
		})
	default:
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleBreakpoints handles breakpoint operations
func (ds *DebugServer) handleBreakpoints(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		ds.sendError(w, "Session ID required", http.StatusBadRequest)
		return
	}

	ds.mutex.RLock()
	session, exists := ds.sessions[sessionID]
	ds.mutex.RUnlock()

	if !exists {
		ds.sendError(w, "Session not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// List breakpoints
		ds.sendResponse(w, DebugResponse{
			Success: true,
			Data:    session.Breakpoints,
		})

	case http.MethodPost:
		// Add breakpoint
		var bp Breakpoint
		if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
			ds.sendError(w, "Invalid breakpoint data", http.StatusBadRequest)
			return
		}

		key := fmt.Sprintf("%s:%d", bp.File, bp.Line)
		bp.Enabled = true
		session.Breakpoints[key] = &bp

		// Also set in global debugger if available
		if debugger != nil {
			SetBreakpoint(bp.File, bp.Line, bp.Function)
		}

		ds.sendResponse(w, DebugResponse{
			Success: true,
			Message: "Breakpoint added",
			Data:    &bp,
		})

	default:
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleStep handles step operations
func (ds *DebugServer) handleStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if debugger != nil {
		debugger.stepMode = true
	}

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Message: "Step command sent",
	})
}

// handleContinue handles continue operations
func (ds *DebugServer) handleContinue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if debugger != nil {
		debugger.continueMode = true
		debugger.stepMode = false
	}

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Message: "Continue command sent",
	})
}

// handleVariables returns current variables
func (ds *DebugServer) handleVariables(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	variables := make(map[string]interface{})
	if debugger != nil && debugger.currentContext != nil {
		variables = debugger.currentContext.Variables
	}

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Data:    variables,
	})
}

// handleStack returns current call stack
func (ds *DebugServer) handleStack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stack := []string{}
	if debugger != nil && debugger.currentContext != nil {
		stack = debugger.currentContext.CallStack
	}

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Data:    stack,
	})
}

// handleEvaluate evaluates expressions in current context
func (ds *DebugServer) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ds.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ds.sendError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// For now, just echo the expression - in a full implementation,
	// this would evaluate the expression in the current context
	result := fmt.Sprintf("Expression: %s (evaluation not yet implemented)", request.Expression)

	ds.sendResponse(w, DebugResponse{
		Success: true,
		Data: map[string]interface{}{
			"expression": request.Expression,
			"result":     result,
		},
	})
}

// sendResponse sends a JSON response
func (ds *DebugServer) sendResponse(w http.ResponseWriter, response DebugResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(response)
}

// sendError sends an error response
func (ds *DebugServer) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	response := DebugResponse{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// splitPath splits URL path into components
func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	// Remove leading/trailing slashes and split
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

// StartDebugServer starts the debug server if debugging is enabled
func StartDebugServer(port int) {
	if !DebugEnabled {
		return
	}

	server := NewDebugServer(port)
	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start debug server: %v\n", err)
	}
}
