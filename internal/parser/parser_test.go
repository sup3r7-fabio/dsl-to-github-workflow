package parser

import (
	"strings"
	"testing"

	"golang-vibe-coding/internal/ast"
	"golang-vibe-coding/internal/lexer"
)

// TestParseSimpleWorkflow tests parsing a basic workflow
func TestParseSimpleWorkflow(t *testing.T) {
	input := `workflow "CI" {
		on = "push"
		
		job "build" {
			runs-on = "ubuntu-latest"
			step "checkout" uses = "actions/checkout@v4"
			step "setup" run = "echo 'Setting up'"
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if workflow.Name != "CI" {
		t.Errorf("workflow.Name not 'CI'. got=%q", workflow.Name)
	}

	if len(workflow.On) != 1 {
		t.Errorf("workflow.On should have 1 trigger. got=%d", len(workflow.On))
	}

	if workflow.On[0].Event != "push" {
		t.Errorf("workflow.On[0].Event not 'push'. got=%q", workflow.On[0].Event)
	}

	if len(workflow.Jobs) != 1 {
		t.Errorf("workflow should have 1 job. got=%d", len(workflow.Jobs))
	}

	job, exists := workflow.Jobs["build"]
	if !exists {
		t.Errorf("job 'build' not found in workflow")
	}

	if job.RunsOn != "ubuntu-latest" {
		t.Errorf("job.RunsOn not 'ubuntu-latest'. got=%q", job.RunsOn)
	}

	if len(job.Steps) != 2 {
		t.Errorf("job should have 2 steps. got=%d", len(job.Steps))
	}
}

// TestParseForLoop tests parsing a for loop
func TestParseForLoop(t *testing.T) {
	input := `workflow "Matrix Build" {
		on = "push"
		
		for i in range(1..3) {
			job "build-${i}" {
				runs-on = "ubuntu-latest"
				step "echo" run = "echo 'Build ${i}'"
			}
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Should have 3 jobs generated from the loop
	expectedJobs := []string{"build-1_1", "build-2_2", "build-3_3"}
	if len(workflow.Jobs) != 3 {
		t.Errorf("workflow should have 3 jobs from loop. got=%d", len(workflow.Jobs))
	}

	for _, jobName := range expectedJobs {
		if _, exists := workflow.Jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in workflow", jobName)
		}
	}
}

// TestParseForeachLoop tests parsing a foreach loop
func TestParseForeachLoop(t *testing.T) {
	input := `workflow "Multi Environment" {
		on = "push"
		
		foreach env in ["dev", "staging", "prod"] {
			job "deploy-${env}" {
				runs-on = "ubuntu-latest"
				step "deploy" run = "deploy to ${env}"
			}
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Should have 3 jobs generated from the foreach loop
	expectedJobs := []string{"deploy-dev_dev", "deploy-staging_staging", "deploy-prod_prod"}
	if len(workflow.Jobs) != 3 {
		t.Errorf("workflow should have 3 jobs from foreach loop. got=%d", len(workflow.Jobs))
	}

	for _, jobName := range expectedJobs {
		if _, exists := workflow.Jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in workflow", jobName)
		}
	}
}

// TestParseRepeatLoop tests parsing a repeat loop
func TestParseRepeatLoop(t *testing.T) {
	input := `workflow "Repeat Test" {
		on = "push"
		
		repeat 2 {
			job "test-iteration" {
				runs-on = "ubuntu-latest"
				step "test" run = "run tests"
			}
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Should have 2 jobs generated from the repeat loop
	if len(workflow.Jobs) != 2 {
		t.Errorf("workflow should have 2 jobs from repeat loop. got=%d", len(workflow.Jobs))
	}

	expectedJobs := []string{"test-iteration_1", "test-iteration_2"}
	for _, jobName := range expectedJobs {
		if _, exists := workflow.Jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in workflow", jobName)
		}
	}
}

// TestParseJobWithComplexSteps tests parsing jobs with complex step configurations
func TestParseJobWithComplexSteps(t *testing.T) {
	input := `workflow "Complex Job" {
		on = "push"
		
		job "complex" {
			runs-on = "ubuntu-latest"
			needs = "setup"
			
			step "checkout" {
				uses = "actions/checkout@v4"
				with {
					fetch-depth = "0"
					token = "${{ secrets.GITHUB_TOKEN }}"
				}
			}
			
			step "setup-node" {
				uses = "actions/setup-node@v4"
				with {
					node-version = "18"
				}
				env {
					NODE_ENV = "production"
				}
				if = "success()"
			}
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	job, exists := workflow.Jobs["complex"]
	if !exists {
		t.Fatalf("job 'complex' not found")
	}

	if len(job.Needs) != 1 || job.Needs[0] != "setup" {
		t.Errorf("job.Needs should be ['setup']. got=%v", job.Needs)
	}

	if len(job.Steps) != 2 {
		t.Fatalf("job should have 2 steps. got=%d", len(job.Steps))
	}

	// Test first step (checkout)
	checkoutStep := job.Steps[0]
	if checkoutStep.Name != "checkout" {
		t.Errorf("step name should be 'checkout'. got=%q", checkoutStep.Name)
	}

	if checkoutStep.Uses != "actions/checkout@v4" {
		t.Errorf("step uses should be 'actions/checkout@v4'. got=%q", checkoutStep.Uses)
	}

	if len(checkoutStep.With) != 2 {
		t.Errorf("step should have 2 with parameters. got=%d", len(checkoutStep.With))
	}

	// Test second step (setup-node)
	setupStep := job.Steps[1]
	if setupStep.If != "success()" {
		t.Errorf("step if condition should be 'success()'. got=%q", setupStep.If)
	}

	if len(setupStep.Env) != 1 {
		t.Errorf("step should have 1 env variable. got=%d", len(setupStep.Env))
	}

	if setupStep.Env["NODE_ENV"] != "production" {
		t.Errorf("NODE_ENV should be 'production'. got=%q", setupStep.Env["NODE_ENV"])
	}
}

// TestParseWorkflowWithEnv tests parsing workflow with environment variables
func TestParseWorkflowWithEnv(t *testing.T) {
	input := `workflow "Env Test" {
		on = "push"
		
		env {
			GLOBAL_VAR = "global_value"
			API_URL = "https://api.example.com"
		}
		
		job "test" {
			runs-on = "ubuntu-latest"
			env {
				JOB_VAR = "job_value"
			}
			step "test" run = "echo 'Testing'"
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()

	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	} // Test global env
	if len(workflow.Env) != 2 {
		t.Errorf("workflow should have 2 global env vars. got=%d", len(workflow.Env))
	}

	if workflow.Env["GLOBAL_VAR"] != "global_value" {
		t.Errorf("GLOBAL_VAR should be 'global_value'. got=%q", workflow.Env["GLOBAL_VAR"])
	}

	// Test job env
	job, exists := workflow.Jobs["test"]
	if !exists {
		t.Fatalf("job 'test' not found in workflow")
	}

	if len(job.Env) != 1 {
		t.Errorf("job should have 1 env var. got=%d", len(job.Env))
	}

	if job.Env["JOB_VAR"] != "job_value" {
		t.Errorf("JOB_VAR should be 'job_value'. got=%q", job.Env["JOB_VAR"])
	}
}

// TestParseErrors tests various error conditions
func TestParseErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "Missing workflow keyword",
			input:         `job "test" {}`,
			expectedError: "program must start with 'workflow' keyword",
		},
		{
			name:          "Missing workflow name",
			input:         `workflow {}`,
			expectedError: "expected next token to be STRING",
		},
		{
			name:          "Missing opening brace",
			input:         `workflow "test"`,
			expectedError: "expected next token to be LBRACE",
		},
		{
			name: "Invalid token in workflow",
			input: `workflow "test" {
				invalid_keyword = "value"
			}`,
			expectedError: "unexpected token in workflow",
		},
		{
			name: "Missing job name",
			input: `workflow "test" {
				job {}
			}`,
			expectedError: "expected next token to be STRING",
		},
		{
			name: "Invalid for loop syntax",
			input: `workflow "test" {
				for {}
			}`,
			expectedError: "expected identifier after for",
		},
		{
			name: "Missing array in foreach",
			input: `workflow "test" {
				foreach env in "not_array" {}
			}`,
			expectedError: "expected next token to be LBRACKET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			workflow := p.ParseProgram()

			if len(p.Errors()) == 0 {
				t.Errorf("expected parsing errors, but got none. workflow=%v", workflow)
			}

			found := false
			for _, err := range p.Errors() {
				if strings.Contains(err, tt.expectedError) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error containing '%s', got errors: %v", tt.expectedError, p.Errors())
			}
		})
	}
}

// TestVariableSubstitution tests the variable substitution functionality
func TestVariableSubstitution(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		context  *ast.LoopContext
		expected string
	}{
		{
			name: "Simple variable substitution",
			text: "deploy-${env}",
			context: &ast.LoopContext{
				Variable: "env",
				Value:    "prod",
				Index:    0,
			},
			expected: "deploy-prod",
		},
		{
			name: "Multiple variable substitutions",
			text: "${env}-${env}-deployment",
			context: &ast.LoopContext{
				Variable: "env",
				Value:    "staging",
				Index:    1,
			},
			expected: "staging-staging-deployment",
		},
		{
			name: "No substitution needed",
			text: "static-text",
			context: &ast.LoopContext{
				Variable: "env",
				Value:    "dev",
				Index:    0,
			},
			expected: "static-text",
		},
		{
			name: "Numeric variable substitution",
			text: "build-${i}",
			context: &ast.LoopContext{
				Variable: "i",
				Value:    5,
				Index:    4,
			},
			expected: "build-5",
		},
	}

	p := &Parser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.substituteVariables(tt.text, tt.context)
			if result != tt.expected {
				t.Errorf("substituteVariables() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestExpandLoops tests the loop expansion functionality
func TestExpandLoops(t *testing.T) {
	// Test for loop expansion
	forLoop := &ast.Loop{
		Type:     ast.LoopTypeFor,
		Variable: "i",
		Range: &ast.LoopRange{
			Start: 1,
			End:   3,
			Step:  1,
		},
		Body: []ast.Node{
			&ast.Job{
				Name:   "build-${i}",
				RunsOn: "ubuntu-latest",
				Steps: []*ast.Step{
					{
						Name: "step-${i}",
						Run:  "echo ${i}",
					},
				},
			},
		},
	}

	p := &Parser{}
	jobs := p.expandForLoop(forLoop)

	if len(jobs) != 3 {
		t.Errorf("expandForLoop should generate 3 jobs. got=%d", len(jobs))
	}

	expectedJobNames := []string{"build-1_1", "build-2_2", "build-3_3"}
	for _, jobName := range expectedJobNames {
		if _, exists := jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in expanded jobs", jobName)
		}
	}

	// Test foreach loop expansion
	foreachLoop := &ast.Loop{
		Type:     ast.LoopTypeForeach,
		Variable: "env",
		Items:    []string{"dev", "prod"},
		Body: []ast.Node{
			&ast.Job{
				Name:   "deploy-${env}",
				RunsOn: "ubuntu-latest",
			},
		},
	}

	jobs = p.expandForeachLoop(foreachLoop)

	if len(jobs) != 2 {
		t.Errorf("expandForeachLoop should generate 2 jobs. got=%d", len(jobs))
	}

	expectedJobNames = []string{"deploy-dev_dev", "deploy-prod_prod"}
	for _, jobName := range expectedJobNames {
		if _, exists := jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in expanded jobs", jobName)
		}
	}
}

// TestParseComplexWorkflow tests parsing a complete workflow with multiple features
func TestParseComplexWorkflow(t *testing.T) {
	input := `workflow "Complete CI/CD" {
		on = "push"
		
		env {
			GLOBAL_VAR = "global"
		}
		
		job "setup" {
			runs-on = "ubuntu-latest"
			step "checkout" uses = "actions/checkout@v4"
		}
		
		foreach env in ["dev", "staging"] {
			job "test-${env}" {
				runs-on = "ubuntu-latest"
				needs = "setup"
				env {
					ENVIRONMENT = "${env}"
				}
				step "test" run = "npm test"
				step "deploy" {
					run = "deploy to ${env}"
					if = "success()"
				}
			}
		}
		
		job "notify" {
			runs-on = "ubuntu-latest"
			needs = ["test-dev_dev", "test-staging_staging"]
			step "notify" run = "send notification"
		}
	}`

	l := lexer.New(input)
	p := New(l)

	workflow := p.ParseProgram()
	checkParserErrors(t, p)

	if workflow == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Should have 4 jobs: setup + 2 from foreach + notify
	if len(workflow.Jobs) != 4 {
		t.Errorf("workflow should have 4 jobs. got=%d", len(workflow.Jobs))
	}

	// Check that jobs exist
	expectedJobs := []string{"setup", "test-dev_dev", "test-staging_staging", "notify"}
	for _, jobName := range expectedJobs {
		if _, exists := workflow.Jobs[jobName]; !exists {
			t.Errorf("job '%s' not found in workflow", jobName)
		}
	}

	// Check notify job has correct dependencies
	notifyJob := workflow.Jobs["notify"]
	if len(notifyJob.Needs) != 2 {
		t.Errorf("notify job should have 2 dependencies. got=%d", len(notifyJob.Needs))
	}
}

// checkParserErrors is a helper function to check for parser errors
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// TestParserNew tests the parser constructor
func TestParserNew(t *testing.T) {
	input := "workflow"
	l := lexer.New(input)
	p := New(l)

	if p == nil {
		t.Fatalf("New() returned nil")
	}

	if p.l != l {
		t.Errorf("parser lexer not set correctly")
	}

	if len(p.errors) != 0 {
		t.Errorf("new parser should have no errors. got=%d", len(p.errors))
	}
}

// TestParserTokenAdvancement tests that tokens are advanced correctly
func TestParserTokenAdvancement(t *testing.T) {
	input := "workflow test"
	l := lexer.New(input)
	p := New(l)

	if p.curToken.Type != lexer.WORKFLOW {
		t.Errorf("current token should be WORKFLOW. got=%s", p.curToken.Type)
	}

	if p.peekToken.Type != lexer.IDENT {
		t.Errorf("peek token should be IDENT. got=%s", p.peekToken.Type)
	}

	p.nextToken()

	if p.curToken.Type != lexer.IDENT {
		t.Errorf("after nextToken(), current should be IDENT. got=%s", p.curToken.Type)
	}
}

// BenchmarkParseSimpleWorkflow benchmarks parsing a simple workflow
func BenchmarkParseSimpleWorkflow(b *testing.B) {
	input := `workflow "CI" {
		on = "push"
		job "build" {
			runs-on = "ubuntu-latest"
			step "checkout" uses = "actions/checkout@v4"
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkParseComplexWorkflow benchmarks parsing a complex workflow with loops
func BenchmarkParseComplexWorkflow(b *testing.B) {
	input := `workflow "Complex" {
		foreach env in ["dev", "staging", "prod"] {
			job "deploy-${env}" {
				runs-on = "ubuntu-latest"
				step "deploy" run = "deploy to ${env}"
			}
		}
		
		for i in range(1..5) {
			job "test-${i}" {
				runs-on = "ubuntu-latest"
				step "test" run = "test ${i}"
			}
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}
