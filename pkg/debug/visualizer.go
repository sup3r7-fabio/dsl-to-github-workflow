// Package debug - AST visualization utilities
package debug

import (
	"fmt"
	"strings"

	"golang-vibe-coding/internal/ast"
)

// ASTVisualizer provides functionality to visualize AST structures
type ASTVisualizer struct {
	indent  int
	output  strings.Builder
	options VisualizationOptions
}

// VisualizationOptions controls AST visualization
type VisualizationOptions struct {
	ShowTypes      bool
	ShowPositions  bool
	ShowMetadata   bool
	MaxDepth       int
	HighlightNodes []string
}

// DefaultVisualizationOptions returns sensible defaults
func DefaultVisualizationOptions() VisualizationOptions {
	return VisualizationOptions{
		ShowTypes:     true,
		ShowPositions: false,
		ShowMetadata:  false,
		MaxDepth:      10,
	}
}

// NewASTVisualizer creates a new AST visualizer
func NewASTVisualizer(opts VisualizationOptions) *ASTVisualizer {
	return &ASTVisualizer{
		options: opts,
	}
}

// VisualizeWorkflow renders a workflow AST as a tree structure
func (v *ASTVisualizer) VisualizeWorkflow(workflow *ast.Workflow) string {
	v.output.Reset()
	v.indent = 0

	if workflow == nil {
		return "❌ No workflow to visualize"
	}

	v.writeNode("🔄 Workflow", workflow.Name, "workflow")
	v.indent++

	// Show triggers
	if len(workflow.On) > 0 {
		v.writeNode("🎯 Triggers", "", "triggers")
		v.indent++
		for _, trigger := range workflow.On {
			v.writeTrigger(&trigger)
		}
		v.indent--
	}

	// Show global environment
	if len(workflow.Env) > 0 {
		v.writeNode("🌍 Global Environment", "", "env")
		v.indent++
		for key, value := range workflow.Env {
			v.writeKeyValue(key, value)
		}
		v.indent--
	}

	// Show jobs
	if len(workflow.Jobs) > 0 {
		v.writeNode("⚙️ Jobs", fmt.Sprintf("(%d)", len(workflow.Jobs)), "jobs")
		v.indent++
		for name, job := range workflow.Jobs {
			v.writeJob(name, job)
		}
		v.indent--
	}

	return v.output.String()
}

// writeTrigger renders a trigger node
func (v *ASTVisualizer) writeTrigger(trigger *ast.Trigger) {
	v.writeNode("📡", trigger.Event, "trigger")

	if len(trigger.Filters) > 0 {
		v.indent++
		for key, value := range trigger.Filters {
			v.writeKeyValue(key, value)
		}
		v.indent--
	}
}

// writeJob renders a job node
func (v *ASTVisualizer) writeJob(name string, job *ast.Job) {
	v.writeNode("🏗️ Job", name, "job")
	v.indent++

	// Basic job properties
	if job.RunsOn != "" {
		v.writeProperty("runs-on", job.RunsOn)
	}

	if job.If != "" {
		v.writeProperty("if", job.If)
	}

	if len(job.Needs) > 0 {
		v.writeProperty("needs", strings.Join(job.Needs, ", "))
	}

	if job.TimeoutMin > 0 {
		v.writeProperty("timeout-minutes", fmt.Sprintf("%d", job.TimeoutMin))
	}

	// Job environment
	if len(job.Env) > 0 {
		v.writeNode("🔧 Environment", "", "env")
		v.indent++
		for key, value := range job.Env {
			v.writeKeyValue(key, value)
		}
		v.indent--
	}

	// Job strategy (matrix)
	if job.Strategy != nil {
		v.writeStrategy(job.Strategy)
	}

	// Job steps
	if len(job.Steps) > 0 {
		v.writeNode("📋 Steps", fmt.Sprintf("(%d)", len(job.Steps)), "steps")
		v.indent++
		for i, step := range job.Steps {
			v.writeStep(i+1, step)
		}
		v.indent--
	}

	v.indent--
}

// writeStrategy renders a strategy node
func (v *ASTVisualizer) writeStrategy(strategy *ast.Strategy) {
	v.writeNode("🎯 Strategy", "", "strategy")
	v.indent++

	if len(strategy.Matrix) > 0 {
		v.writeNode("📊 Matrix", "", "matrix")
		v.indent++
		for key, values := range strategy.Matrix {
			valueStr := "["
			for i, val := range values {
				if i > 0 {
					valueStr += ", "
				}
				valueStr += fmt.Sprintf("%v", val)
			}
			valueStr += "]"
			v.writeKeyValue(key, valueStr)
		}
		v.indent--
	}

	if strategy.FailFast != nil {
		v.writeProperty("fail-fast", fmt.Sprintf("%t", *strategy.FailFast))
	}

	if strategy.MaxParallel > 0 {
		v.writeProperty("max-parallel", fmt.Sprintf("%d", strategy.MaxParallel))
	}

	v.indent--
}

// writeStep renders a step node
func (v *ASTVisualizer) writeStep(index int, step *ast.Step) {
	stepTitle := fmt.Sprintf("Step %d", index)
	if step.Name != "" {
		stepTitle = fmt.Sprintf("Step %d: %s", index, step.Name)
	}

	v.writeNode("📝", stepTitle, "step")
	v.indent++

	if step.Uses != "" {
		v.writeProperty("uses", step.Uses)
	}

	if step.Run != "" {
		v.writeProperty("run", v.truncateString(step.Run, 50))
	}

	if step.If != "" {
		v.writeProperty("if", step.If)
	}

	// Note: ContinueOnError not available in current Step struct

	// Step with parameters
	if len(step.With) > 0 {
		v.writeNode("⚙️ With", "", "with")
		v.indent++
		for key, value := range step.With {
			v.writeKeyValue(key, value)
		}
		v.indent--
	}

	// Step environment
	if len(step.Env) > 0 {
		v.writeNode("🔧 Environment", "", "env")
		v.indent++
		for key, value := range step.Env {
			v.writeKeyValue(key, value)
		}
		v.indent--
	}

	v.indent--
}

// writeNode writes a node with icon and type information
func (v *ASTVisualizer) writeNode(icon, name, nodeType string) {
	prefix := v.getIndentPrefix()

	line := fmt.Sprintf("%s%s %s", prefix, icon, name)

	if v.options.ShowTypes {
		line += fmt.Sprintf(" <%s>", nodeType)
	}

	// Highlight specific nodes if requested
	for _, highlight := range v.options.HighlightNodes {
		if strings.Contains(strings.ToLower(name), strings.ToLower(highlight)) ||
			strings.Contains(strings.ToLower(nodeType), strings.ToLower(highlight)) {
			line = ">>> " + line + " <<<"
			break
		}
	}

	v.output.WriteString(line + "\n")
}

// writeProperty writes a property line
func (v *ASTVisualizer) writeProperty(key, value string) {
	prefix := v.getIndentPrefix()
	v.output.WriteString(fmt.Sprintf("%s  %s: %s\n", prefix, key, value))
}

// writeKeyValue writes a key-value pair
func (v *ASTVisualizer) writeKeyValue(key, value string) {
	prefix := v.getIndentPrefix()
	v.output.WriteString(fmt.Sprintf("%s  %s = %s\n", prefix, key, value))
}

// getIndentPrefix returns the current indentation prefix
func (v *ASTVisualizer) getIndentPrefix() string {
	if v.indent == 0 {
		return ""
	}

	prefix := ""
	for i := 0; i < v.indent-1; i++ {
		prefix += "│   "
	}
	prefix += "├── "

	return prefix
}

// truncateString truncates a string to the specified length
func (v *ASTVisualizer) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// VisualizationToHTML converts AST visualization to HTML
func VisualizationToHTML(visualization string) string {
	html := `<div style="font-family: monospace; white-space: pre-wrap; color: #d4d4d4; background-color: #1e1e1e; padding: 16px; border-radius: 4px;">`

	lines := strings.Split(visualization, "\n")
	for _, line := range lines {
		// Color-code different elements
		styledLine := line

		// Highlight different node types
		if strings.Contains(line, "🔄 Workflow") {
			styledLine = fmt.Sprintf(`<span style="color: #4fc3f7; font-weight: bold;">%s</span>`, line)
		} else if strings.Contains(line, "🏗️ Job") {
			styledLine = fmt.Sprintf(`<span style="color: #81c784; font-weight: bold;">%s</span>`, line)
		} else if strings.Contains(line, "📝") && strings.Contains(line, "Step") {
			styledLine = fmt.Sprintf(`<span style="color: #ffb74d;">%s</span>`, line)
		} else if strings.Contains(line, "🎯 Triggers") || strings.Contains(line, "📡") {
			styledLine = fmt.Sprintf(`<span style="color: #f06292;">%s</span>`, line)
		} else if strings.Contains(line, "🌍") || strings.Contains(line, "🔧") {
			styledLine = fmt.Sprintf(`<span style="color: #ba68c8;">%s</span>`, line)
		} else if strings.Contains(line, ">>>") && strings.Contains(line, "<<<") {
			styledLine = fmt.Sprintf(`<span style="background-color: #ffeb3b; color: #000; font-weight: bold;">%s</span>`, line)
		}

		html += styledLine + "\n"
	}

	html += "</div>"
	return html
}

// AnalyzeAST provides analysis of the AST structure
func AnalyzeAST(workflow *ast.Workflow) map[string]interface{} {
	if workflow == nil {
		return map[string]interface{}{
			"error": "No workflow provided",
		}
	}

	analysis := map[string]interface{}{
		"workflow_name":    workflow.Name,
		"total_jobs":       len(workflow.Jobs),
		"total_triggers":   len(workflow.On),
		"has_global_env":   len(workflow.Env) > 0,
		"has_defaults":     len(workflow.Defaults) > 0,
		"jobs_analysis":    make(map[string]interface{}),
		"complexity_score": 0,
	}

	totalSteps := 0
	matrixJobs := 0
	conditionalJobs := 0

	for jobName, job := range workflow.Jobs {
		jobAnalysis := map[string]interface{}{
			"steps_count":       len(job.Steps),
			"has_strategy":      job.Strategy != nil,
			"has_conditions":    job.If != "",
			"has_dependencies":  len(job.Needs) > 0,
			"has_environment":   len(job.Env) > 0,
			"timeout_minutes":   job.TimeoutMin,
			"continue_on_error": job.ContinueOnErr,
		}

		if job.Strategy != nil && len(job.Strategy.Matrix) > 0 {
			matrixJobs++
			matrixSize := 1
			for _, values := range job.Strategy.Matrix {
				matrixSize *= len(values)
			}
			jobAnalysis["matrix_size"] = matrixSize
		}

		if job.If != "" {
			conditionalJobs++
		}

		totalSteps += len(job.Steps)
		analysis["jobs_analysis"].(map[string]interface{})[jobName] = jobAnalysis
	}

	// Calculate complexity score
	complexityScore := len(workflow.Jobs) * 2
	complexityScore += totalSteps
	complexityScore += matrixJobs * 5
	complexityScore += conditionalJobs * 2

	analysis["total_steps"] = totalSteps
	analysis["matrix_jobs"] = matrixJobs
	analysis["conditional_jobs"] = conditionalJobs
	analysis["complexity_score"] = complexityScore

	return analysis
}

// LogAST logs the AST visualization to the debug output
func LogAST(workflow *ast.Workflow) {
	if !DebugEnabled {
		return
	}

	visualizer := NewASTVisualizer(DefaultVisualizationOptions())
	visualization := visualizer.VisualizeWorkflow(workflow)

	Log("AST Visualization:\n%s", visualization)

	if debugger != nil && debugger.level >= LevelVerbose {
		analysis := AnalyzeAST(workflow)
		Log("AST Analysis: %+v", analysis)
	}
}
