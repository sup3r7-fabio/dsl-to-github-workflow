# Testing Documentation

## Parser Testing

### Test Coverage
The parser has **78.5% test coverage** with comprehensive unit tests covering:

- **Basic workflow parsing**: Simple workflows with jobs, steps, and triggers
- **Complex workflows**: Nested structures with environment variables
- **Loop constructs**: For loops, foreach loops, and repeat loops with variable expansion
- **Error handling**: Table-driven tests for various error conditions
- **Variable substitution**: Template variable replacement in loop contexts
- **Performance**: Benchmark tests showing ~1μs for simple workflows

### Test Structure

#### Core Functionality Tests
- `TestParseSimpleWorkflow`: Basic workflow with jobs and steps
- `TestParseForLoop`: For loop with range expansion (1..3)
- `TestParseForeachLoop`: Foreach loop with array iteration 
- `TestParseRepeatLoop`: Repeat loop with count-based iteration
- `TestParseJobWithComplexSteps`: Jobs with complex step configurations
- `TestParseWorkflowWithEnv`: Environment variable parsing at workflow and job levels
- `TestParseComplexWorkflow`: Complete workflow combining multiple features

#### Error Handling Tests
- `TestParseErrors`: Table-driven tests covering:
  - Missing workflow keyword
  - Missing workflow name
  - Missing opening brace
  - Invalid tokens
  - Missing job names
  - Invalid loop syntax
  - Missing array elements

#### Advanced Features Tests
- `TestVariableSubstitution`: Variable replacement in loop contexts
- `TestExpandLoops`: Loop expansion functionality
- `TestParserNew`: Constructor validation
- `TestParserTokenAdvancement`: Token management

### Performance Benchmarks
- **BenchmarkParseSimpleWorkflow**: ~1,032 ns/op
- **BenchmarkParseComplexWorkflow**: ~10,405 ns/op

### Coverage by Function
| Function | Coverage | Notes |
|----------|----------|-------|
| Basic parser methods | 100% | New, nextToken, Errors, etc. |
| parseWorkflow | 93.8% | Core workflow parsing |
| parseNeeds | 90.5% | Job dependency parsing |
| expandForLoop | 90.9% | Loop expansion logic |
| parseLoop | 88.2% | Loop type detection |
| Loop expansion methods | 100% | Foreach, repeat, variable substitution |
| parseJob | 65.0% | Job parsing with error paths |
| parseStep | 77.8% | Step parsing |
| parseWithBlock | 68.4% | Action input parsing |
| parseEnvBlock | 71.4% | Environment variable parsing |

### Running Tests

```bash
# Run all parser tests
go test ./internal/parser -v

# Run tests with coverage
go test ./internal/parser -cover

# Run benchmarks
go test ./internal/parser -bench=.

# Generate coverage report
go test ./internal/parser -coverprofile=coverage.out
go tool cover -func coverage.out
```

### Bug Fixes During Testing

1. **Environment Block Parsing**: Fixed parseEnvBlock not advancing past closing brace
2. **Infinite Loop Prevention**: Added token advancement in error recovery
3. **Parser State Management**: Improved error handling in parseWorkflow

These tests ensure the parser correctly handles the DSL syntax and provides reliable error reporting for invalid input.
