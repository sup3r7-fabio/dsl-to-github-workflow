// Package ast defines the Abstract Syntax Tree structures for GitHub Actions DSL
package ast

import "fmt"

// Node represents a node in the AST
type Node interface {
	String() string
}

// Workflow represents a complete GitHub Actions workflow
type Workflow struct {
	Name     string            `json:"name"`
	On       []Trigger         `json:"on"`
	Jobs     map[string]*Job   `json:"jobs"`
	Env      map[string]string `json:"env,omitempty"`      // Global environment variables
	Defaults map[string]string `json:"defaults,omitempty"` // Global defaults
}

func (w *Workflow) String() string {
	return fmt.Sprintf("Workflow{Name: %s, Jobs: %d}", w.Name, len(w.Jobs))
}

// Trigger represents workflow triggers (push, pull_request, etc.)
type Trigger struct {
	Event   string            `json:"event"`             // push, pull_request, schedule, etc.
	Filters map[string]string `json:"filters,omitempty"` // branches, paths, etc.
}

func (t *Trigger) String() string {
	return fmt.Sprintf("Trigger{Event: %s}", t.Event)
}

// Job represents a job within a workflow
type Job struct {
	Name          string            `json:"name"`
	RunsOn        string            `json:"runs-on"`
	Steps         []*Step           `json:"steps"`
	Needs         []string          `json:"needs,omitempty"`    // Job dependencies
	If            string            `json:"if,omitempty"`       // Conditional execution
	Env           map[string]string `json:"env,omitempty"`      // Job-level environment variables
	Strategy      *Strategy         `json:"strategy,omitempty"` // Matrix strategy
	TimeoutMin    int               `json:"timeout-minutes,omitempty"`
	ContinueOnErr bool              `json:"continue-on-error,omitempty"`
}

func (j *Job) String() string {
	return fmt.Sprintf("Job{Name: %s, RunsOn: %s, Steps: %d}", j.Name, j.RunsOn, len(j.Steps))
}

// Strategy represents job strategy (matrix builds, etc.)
type Strategy struct {
	Matrix      map[string][]string `json:"matrix,omitempty"`
	FailFast    *bool               `json:"fail-fast,omitempty"`
	MaxParallel int                 `json:"max-parallel,omitempty"`
}

func (s *Strategy) String() string {
	return fmt.Sprintf("Strategy{Matrix: %v}", s.Matrix)
}

// Step represents a single step within a job
type Step struct {
	Name string            `json:"name,omitempty"`
	ID   string            `json:"id,omitempty"`
	Uses string            `json:"uses,omitempty"` // For action steps
	Run  string            `json:"run,omitempty"`  // For run steps
	With map[string]string `json:"with,omitempty"` // Action inputs
	Env  map[string]string `json:"env,omitempty"`  // Step environment variables
	If   string            `json:"if,omitempty"`   // Conditional execution
}

func (s *Step) String() string {
	if s.Uses != "" {
		return fmt.Sprintf("Step{Name: %s, Uses: %s}", s.Name, s.Uses)
	}
	return fmt.Sprintf("Step{Name: %s, Run: %s}", s.Name, s.Run)
}

// StepType represents the type of step (action or run)
type StepType int

const (
	StepTypeAction StepType = iota
	StepTypeRun
)

// Assignment represents a key-value assignment in the DSL
type Assignment struct {
	Key   string
	Value string
}

func (a *Assignment) String() string {
	return fmt.Sprintf("Assignment{%s = %s}", a.Key, a.Value)
}

// Block represents a block structure in the DSL (workflow, job, step)
type Block struct {
	Type        string
	Name        string
	Assignments []*Assignment
	Children    []*Block
}

func (b *Block) String() string {
	return fmt.Sprintf("Block{Type: %s, Name: %s, Assignments: %d, Children: %d}",
		b.Type, b.Name, len(b.Assignments), len(b.Children))
}

// Loop represents different types of loops in the DSL
type Loop struct {
	Type     LoopType   `json:"type"`               // for, foreach, repeat
	Variable string     `json:"variable,omitempty"` // loop variable name
	Range    *LoopRange `json:"range,omitempty"`    // for range loops
	Items    []string   `json:"items,omitempty"`    // foreach items
	Count    int        `json:"count,omitempty"`    // repeat count
	Body     []Node     `json:"body"`               // loop body (jobs/steps)
}

func (l *Loop) String() string {
	return fmt.Sprintf("Loop{Type: %s, Variable: %s}", l.Type.String(), l.Variable)
}

// LoopType represents different loop types
type LoopType int

const (
	LoopTypeFor LoopType = iota
	LoopTypeForeach
	LoopTypeRepeat
)

func (lt LoopType) String() string {
	switch lt {
	case LoopTypeFor:
		return "for"
	case LoopTypeForeach:
		return "foreach"
	case LoopTypeRepeat:
		return "repeat"
	default:
		return "unknown"
	}
}

// LoopRange represents a range for 'for' loops
type LoopRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Step  int `json:"step,omitempty"` // default 1
}

func (lr *LoopRange) String() string {
	if lr.Step > 0 && lr.Step != 1 {
		return fmt.Sprintf("Range{%d..%d step %d}", lr.Start, lr.End, lr.Step)
	}
	return fmt.Sprintf("Range{%d..%d}", lr.Start, lr.End)
}

// LoopContext provides context for loop expansion
type LoopContext struct {
	Variable string      // current loop variable
	Value    interface{} // current loop value (int, string, etc.)
	Index    int         // current iteration index
}

func (lc *LoopContext) String() string {
	return fmt.Sprintf("LoopContext{%s=%v, index=%d}", lc.Variable, lc.Value, lc.Index)
}
