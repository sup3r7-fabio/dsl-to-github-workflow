// Package parser provides parsing functionality for GitHub Actions DSL
package parser

import (
	"fmt"
	"strconv"

	"golang-vibe-coding/internal/ast"
	"golang-vibe-coding/internal/lexer"
	"golang-vibe-coding/pkg/debug"
)

// Parser represents the parser state
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// addError adds a parsing error
func (p *Parser) addError(msg string) {
	errorMsg := fmt.Sprintf("line %d, column %d: %s", p.curToken.Line, p.curToken.Column, msg)
	p.errors = append(p.errors, errorMsg)
}

// expectPeek checks the next token type and advances if it matches
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
		return false
	}
}

// curTokenIs checks if current token matches type
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if peek token matches type
func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// skipNewlines skips any newline tokens
func (p *Parser) skipNewlines() {
	for p.curTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}
}

// ParseProgram parses the entire DSL program
func (p *Parser) ParseProgram() *ast.Workflow {
	p.skipNewlines()

	if !p.curTokenIs(lexer.WORKFLOW) {
		p.addError("program must start with 'workflow' keyword")
		return nil
	}

	workflow := p.parseWorkflow()
	return workflow
}

// parseWorkflow parses a workflow block
func (p *Parser) parseWorkflow() *ast.Workflow {
	workflow := &ast.Workflow{
		Jobs: make(map[string]*ast.Job),
		Env:  make(map[string]string),
	}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}
	workflow.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch p.curToken.Type {
		case lexer.ON:
			triggers := p.parseOnClause()
			workflow.On = triggers
		case lexer.JOB:
			job := p.parseJob()
			if job != nil {
				workflow.Jobs[job.Name] = job
			} else {
				// If parseJob failed, advance past tokens to avoid infinite loop
				p.nextToken()
			}
		case lexer.FOR, lexer.FOREACH, lexer.REPEAT:
			loop := p.parseLoop()
			if loop != nil {
				// Expand loop and add generated jobs to workflow
				jobs := p.expandLoop(loop)
				for name, job := range jobs {
					workflow.Jobs[name] = job
				}
			} else {
				// If parseLoop failed, advance past tokens to avoid infinite loop
				p.nextToken()
			}
		case lexer.ENV:
			env := p.parseEnvBlock()
			if env != nil {
				workflow.Env = env
			}
		default:
			p.addError(fmt.Sprintf("unexpected token in workflow: %s", p.curToken.Type))
			p.nextToken()
		}
		p.skipNewlines()
	}

	if !p.curTokenIs(lexer.RBRACE) {
		p.addError("expected '}' at end of workflow")
		return nil
	}

	return workflow
}

// parseOnClause parses the 'on' clause for workflow triggers
func (p *Parser) parseOnClause() []ast.Trigger {
	var triggers []ast.Trigger

	if !p.expectPeek(lexer.ASSIGN) {
		return triggers
	}

	if !p.expectPeek(lexer.STRING) {
		return triggers
	}

	// For now, support simple string triggers like "push", "pull_request"
	event := p.curToken.Literal
	trigger := ast.Trigger{
		Event:   event,
		Filters: make(map[string]string),
	}
	triggers = append(triggers, trigger)

	p.nextToken()
	return triggers
}

// parseJob parses a job block
func (p *Parser) parseJob() *ast.Job {
	job := &ast.Job{
		Steps: []*ast.Step{},
		Env:   make(map[string]string),
		Needs: []string{},
	}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}
	job.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch p.curToken.Type {
		case lexer.IDENT:
			// Handle assignments like runs-on = "ubuntu-latest"
			switch p.curToken.Literal {
			case "runs-on":
				if !p.expectPeek(lexer.ASSIGN) {
					return nil
				}
				if !p.expectPeek(lexer.STRING) {
					return nil
				}
				job.RunsOn = p.curToken.Literal
				p.nextToken()
			case "timeout-minutes":
				// Handle timeout-minutes assignment (simplified)
				p.nextToken() // skip assignment and value for now
				p.nextToken()
			default:
				p.addError(fmt.Sprintf("unexpected identifier in job: %s", p.curToken.Literal))
				p.nextToken()
			}
		case lexer.STEP:
			step := p.parseStep()
			if step != nil {
				job.Steps = append(job.Steps, step)
			}
		case lexer.NEEDS:
			needs := p.parseNeeds()
			job.Needs = needs
		case lexer.ENV:
			env := p.parseEnvBlock()
			if env != nil {
				job.Env = env
			}
		case lexer.IF:
			condition := p.parseIfCondition()
			job.If = condition
			p.nextToken()
		default:
			p.addError(fmt.Sprintf("unexpected token in job: %s", p.curToken.Type))
			p.nextToken()
		}
		p.skipNewlines()
	}

	if !p.curTokenIs(lexer.RBRACE) {
		p.addError("expected '}' at end of job")
		return nil
	}

	// Advance past the closing brace
	p.nextToken()

	return job
}

// parseStep parses a step block or inline step
func (p *Parser) parseStep() *ast.Step {
	step := &ast.Step{
		With: make(map[string]string),
		Env:  make(map[string]string),
	}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}
	step.Name = p.curToken.Literal

	// Check for inline step (uses = "action" or run = "command")
	if p.peekTokenIs(lexer.USES) {
		p.nextToken() // move to "uses"
		if !p.expectPeek(lexer.ASSIGN) {
			return nil
		}
		if !p.expectPeek(lexer.STRING) {
			return nil
		}
		step.Uses = p.curToken.Literal
		p.nextToken()
		return step
	} else if p.peekTokenIs(lexer.RUN) {
		p.nextToken() // move to "run"
		if !p.expectPeek(lexer.ASSIGN) {
			return nil
		}
		if !p.expectPeek(lexer.STRING) {
			return nil
		}
		step.Run = p.curToken.Literal
		p.nextToken()
		return step
	}

	// Handle block-style step
	if p.peekTokenIs(lexer.LBRACE) {
		p.nextToken() // move to "{"
		p.nextToken()
		p.skipNewlines()

		for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
			switch p.curToken.Type {
			case lexer.USES:
				if !p.expectPeek(lexer.ASSIGN) {
					return nil
				}
				if !p.expectPeek(lexer.STRING) {
					return nil
				}
				step.Uses = p.curToken.Literal
			case lexer.RUN:
				if !p.expectPeek(lexer.ASSIGN) {
					return nil
				}
				if !p.expectPeek(lexer.STRING) {
					return nil
				}
				step.Run = p.curToken.Literal
			case lexer.WITH:
				with := p.parseWithBlock()
				if with != nil {
					step.With = with
				}
			case lexer.ENV:
				env := p.parseEnvBlock()
				if env != nil {
					step.Env = env
				}
			case lexer.IF:
				condition := p.parseIfCondition()
				step.If = condition
			default:
				p.addError(fmt.Sprintf("unexpected token in step: %s", p.curToken.Type))
			}
			p.nextToken()
			p.skipNewlines()
		}

		if !p.curTokenIs(lexer.RBRACE) {
			p.addError("expected '}' at end of step")
			return nil
		}
	}

	p.nextToken()
	return step
}

// parseWithBlock parses a 'with' block for action inputs
func (p *Parser) parseWithBlock() map[string]string {
	with := make(map[string]string)

	if !p.expectPeek(lexer.LBRACE) {
		return with
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		if !p.curTokenIs(lexer.IDENT) && !p.curTokenIs(lexer.STRING) {
			p.addError("expected identifier or string in with block")
			p.nextToken()
			continue
		}

		key := p.curToken.Literal
		if !p.expectPeek(lexer.ASSIGN) {
			return with
		}
		if !p.expectPeek(lexer.STRING) {
			return with
		}

		with[key] = p.curToken.Literal
		p.nextToken()
		p.skipNewlines()
	}

	return with
}

// parseEnvBlock parses an 'env' block for environment variables
func (p *Parser) parseEnvBlock() map[string]string {
	env := make(map[string]string)

	if !p.expectPeek(lexer.LBRACE) {
		return env
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		if !p.curTokenIs(lexer.IDENT) && !p.curTokenIs(lexer.STRING) {
			p.addError("expected identifier or string in env block")
			p.nextToken()
			continue
		}

		key := p.curToken.Literal
		if !p.expectPeek(lexer.ASSIGN) {
			return env
		}
		if !p.expectPeek(lexer.STRING) {
			return env
		}

		env[key] = p.curToken.Literal
		p.nextToken()
		p.skipNewlines()
	}

	// Advance past the closing brace
	if p.curTokenIs(lexer.RBRACE) {
		p.nextToken()
	}

	return env
}

// parseNeeds parses a 'needs' clause for job dependencies
func (p *Parser) parseNeeds() []string {
	var needs []string

	if !p.expectPeek(lexer.ASSIGN) {
		return needs
	}

	// Handle single dependency
	if p.peekTokenIs(lexer.STRING) {
		p.nextToken()
		needs = append(needs, p.curToken.Literal)
		p.nextToken()
		return needs
	}

	// Handle array of dependencies (simplified - expect comma-separated strings)
	if p.peekTokenIs(lexer.LBRACKET) {
		p.nextToken() // move to "["
		p.nextToken()

		for !p.curTokenIs(lexer.RBRACKET) && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.STRING) {
				needs = append(needs, p.curToken.Literal)
			}
			p.nextToken()
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}

		if !p.curTokenIs(lexer.RBRACKET) {
			p.addError("expected ']' at end of needs array")
		}
		p.nextToken()
	}

	return needs
}

// parseIfCondition parses an 'if' condition
func (p *Parser) parseIfCondition() string {
	if !p.expectPeek(lexer.ASSIGN) {
		return ""
	}
	if !p.expectPeek(lexer.STRING) {
		return ""
	}

	condition := p.curToken.Literal
	return condition
}

// parseLoop parses different types of loops (for, foreach, repeat)
func (p *Parser) parseLoop() *ast.Loop {
	debug.Trace("parseLoop", "ENTER", fmt.Sprintf("token=%s", p.curToken.Type))
	loop := &ast.Loop{}

	switch p.curToken.Type {
	case lexer.FOR:
		debug.Log("Parsing FOR loop")
		loop.Type = ast.LoopTypeFor
		p.nextToken() // consume FOR token
		return p.parseForLoop(loop)
	case lexer.FOREACH:
		debug.Log("Parsing FOREACH loop")
		loop.Type = ast.LoopTypeForeach
		p.nextToken() // consume FOREACH token
		return p.parseForeachLoop(loop)
	case lexer.REPEAT:
		debug.Log("Parsing REPEAT loop")
		loop.Type = ast.LoopTypeRepeat
		p.nextToken() // consume REPEAT token
		return p.parseRepeatLoop(loop)
	default:
		p.addError(fmt.Sprintf("unexpected loop type: %s", p.curToken.Type))
		return nil
	}
}

// parseForLoop parses 'for i in range(1..5)' or 'for i in 1..5'
func (p *Parser) parseForLoop(loop *ast.Loop) *ast.Loop {
	// for variable in range(...) { ... }
	if !p.curTokenIs(lexer.IDENT) {
		p.addError("expected identifier after for")
		return nil
	}
	loop.Variable = p.curToken.Literal

	if !p.expectPeek(lexer.IN) {
		return nil
	}

	// Parse range: either range(1..5) or just 1..5
	if p.peekTokenIs(lexer.RANGE) {
		p.nextToken() // move to "range"
		if !p.expectPeek(lexer.LPAREN) {
			return nil
		}
		loop.Range = p.parseRange()
		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
	} else {
		// Direct range: 1..5
		loop.Range = p.parseRange()
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	loop.Body = p.parseLoopBody()
	// parseLoopBody advances past the closing brace

	return loop
}

// parseForeachLoop parses 'foreach env in ["dev", "staging", "prod"]'
func (p *Parser) parseForeachLoop(loop *ast.Loop) *ast.Loop {
	// foreach variable in [items] { ... }
	// Allow both IDENT and certain keywords as variable names
	if !p.curTokenIs(lexer.IDENT) && !p.curTokenIs(lexer.ENV) && !p.curTokenIs(lexer.STRATEGY) {
		p.addError("expected identifier after foreach")
		return nil
	}
	loop.Variable = p.curToken.Literal

	if !p.expectPeek(lexer.IN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACKET) {
		return nil
	}

	// Parse array items carefully to avoid infinite loops
	var items []string

	// Move to first element after '['
	p.nextToken()

	// Skip any newlines
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACKET) && !p.curTokenIs(lexer.EOF) {
		if p.curTokenIs(lexer.STRING) {
			items = append(items, p.curToken.Literal)
			p.nextToken()
			p.skipNewlines()

			// Skip comma if present
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
				p.skipNewlines()
			}
		} else if p.curTokenIs(lexer.NEWLINE) {
			p.nextToken()
		} else {
			// Unexpected token - report error and try to recover
			p.addError(fmt.Sprintf("unexpected token in array: %s, expected string or ']'", p.curToken.Type))
			p.nextToken()
		}
	}

	loop.Items = items

	// We should be at RBRACKET now
	if !p.curTokenIs(lexer.RBRACKET) {
		p.addError("expected ']' to close array")
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	loop.Body = p.parseLoopBody()
	// parseLoopBody advances past the closing brace

	return loop
} // parseRepeatLoop parses 'repeat 3 { ... }'
func (p *Parser) parseRepeatLoop(loop *ast.Loop) *ast.Loop {
	if !p.curTokenIs(lexer.NUMBER) {
		p.addError("expected number after repeat")
		return nil
	}

	count, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid number: %s", p.curToken.Literal))
		return nil
	}
	loop.Count = count

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	loop.Body = p.parseLoopBody()
	// parseLoopBody advances past the closing brace

	return loop
}

// parseRange parses range expressions like 1..5 or 0..10
func (p *Parser) parseRange() *ast.LoopRange {
	if !p.expectPeek(lexer.NUMBER) {
		return nil
	}

	start, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid start number: %s", p.curToken.Literal))
		return nil
	}

	if !p.expectPeek(lexer.DOTDOT) {
		return nil
	}

	if !p.expectPeek(lexer.NUMBER) {
		return nil
	}

	end, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid end number: %s", p.curToken.Literal))
		return nil
	}

	return &ast.LoopRange{
		Start: start,
		End:   end,
		Step:  1, // default step
	}
}

// parseLoopBody parses the body of a loop (jobs or steps)
func (p *Parser) parseLoopBody() []ast.Node {
	var body []ast.Node
	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch p.curToken.Type {
		case lexer.JOB:
			job := p.parseJob()
			if job != nil {
				body = append(body, job)
			}
		case lexer.STEP:
			step := p.parseStep()
			if step != nil {
				body = append(body, step)
			}
		default:
			p.addError(fmt.Sprintf("unexpected token in loop body: %s", p.curToken.Type))
			p.nextToken()
		}
		p.skipNewlines()
	}

	// Advance past the closing brace
	if p.curTokenIs(lexer.RBRACE) {
		p.nextToken()
	}

	return body
}

// expandLoop expands a loop into multiple jobs
func (p *Parser) expandLoop(loop *ast.Loop) map[string]*ast.Job {
	jobs := make(map[string]*ast.Job)

	switch loop.Type {
	case ast.LoopTypeFor:
		return p.expandForLoop(loop)
	case ast.LoopTypeForeach:
		return p.expandForeachLoop(loop)
	case ast.LoopTypeRepeat:
		return p.expandRepeatLoop(loop)
	}

	return jobs
}

// expandForLoop expands a for loop into multiple jobs
func (p *Parser) expandForLoop(loop *ast.Loop) map[string]*ast.Job {
	jobs := make(map[string]*ast.Job)

	if loop.Range == nil {
		return jobs
	}

	for i := loop.Range.Start; i <= loop.Range.End; i += loop.Range.Step {
		context := &ast.LoopContext{
			Variable: loop.Variable,
			Value:    i,
			Index:    i - loop.Range.Start,
		}

		for _, node := range loop.Body {
			if job, ok := node.(*ast.Job); ok {
				expandedJob := p.expandJobWithContext(job, context)
				// Use the substituted name for the job key
				jobName := p.substituteVariables(job.Name, context) + fmt.Sprintf("_%d", i)
				jobs[jobName] = expandedJob
			}
		}
	}

	return jobs
}

// expandForeachLoop expands a foreach loop into multiple jobs
func (p *Parser) expandForeachLoop(loop *ast.Loop) map[string]*ast.Job {
	jobs := make(map[string]*ast.Job)

	for index, item := range loop.Items {
		context := &ast.LoopContext{
			Variable: loop.Variable,
			Value:    item,
			Index:    index,
		}

		for _, node := range loop.Body {
			if job, ok := node.(*ast.Job); ok {
				expandedJob := p.expandJobWithContext(job, context)
				// Use the substituted name for the job key
				jobName := p.substituteVariables(job.Name, context) + fmt.Sprintf("_%s", item)
				jobs[jobName] = expandedJob
			}
		}
	}

	return jobs
}

// expandRepeatLoop expands a repeat loop into multiple jobs
func (p *Parser) expandRepeatLoop(loop *ast.Loop) map[string]*ast.Job {
	jobs := make(map[string]*ast.Job)

	for i := 0; i < loop.Count; i++ {
		context := &ast.LoopContext{
			Variable: fmt.Sprintf("iteration_%d", i),
			Value:    i + 1, // 1-based counting
			Index:    i,
		}

		for _, node := range loop.Body {
			if job, ok := node.(*ast.Job); ok {
				expandedJob := p.expandJobWithContext(job, context)
				// Use the substituted name for the job key
				jobName := p.substituteVariables(job.Name, context) + fmt.Sprintf("_%d", i+1)
				jobs[jobName] = expandedJob
			}
		}
	}

	return jobs
}

// expandJobWithContext expands a job template with loop context
func (p *Parser) expandJobWithContext(job *ast.Job, context *ast.LoopContext) *ast.Job {
	expanded := &ast.Job{
		Name:          p.substituteVariables(job.Name, context),
		RunsOn:        p.substituteVariables(job.RunsOn, context),
		Steps:         make([]*ast.Step, len(job.Steps)),
		Needs:         make([]string, len(job.Needs)),
		If:            p.substituteVariables(job.If, context),
		Env:           make(map[string]string),
		TimeoutMin:    job.TimeoutMin,
		ContinueOnErr: job.ContinueOnErr,
	}

	// Copy and substitute steps
	for i, step := range job.Steps {
		expanded.Steps[i] = &ast.Step{
			Name: p.substituteVariables(step.Name, context),
			ID:   p.substituteVariables(step.ID, context),
			Uses: p.substituteVariables(step.Uses, context),
			Run:  p.substituteVariables(step.Run, context),
			With: make(map[string]string),
			Env:  make(map[string]string),
			If:   p.substituteVariables(step.If, context),
		}

		// Copy with parameters
		for k, v := range step.With {
			expanded.Steps[i].With[k] = p.substituteVariables(v, context)
		}

		// Copy env variables
		for k, v := range step.Env {
			expanded.Steps[i].Env[k] = p.substituteVariables(v, context)
		}
	}

	// Copy needs (job dependencies)
	for i, need := range job.Needs {
		expanded.Needs[i] = p.substituteVariables(need, context)
	}

	// Copy env variables
	for k, v := range job.Env {
		expanded.Env[k] = p.substituteVariables(v, context)
	}

	return expanded
}

// substituteVariables substitutes loop variables in strings
func (p *Parser) substituteVariables(text string, context *ast.LoopContext) string {
	if text == "" {
		return text
	}

	// Simple variable substitution - replace ${variable} with value
	result := text

	// Replace loop variable
	varPlaceholder := fmt.Sprintf("${%s}", context.Variable)
	valueStr := fmt.Sprintf("%v", context.Value)

	// Replace all occurrences
	for i := 0; i < len(result); i++ {
		if i+len(varPlaceholder) <= len(result) && result[i:i+len(varPlaceholder)] == varPlaceholder {
			result = result[:i] + valueStr + result[i+len(varPlaceholder):]
			i += len(valueStr) - 1 // adjust index
		}
	}

	return result
}

// Error represents a parser error
type Error struct {
	Line    int
	Column  int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("parser error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
