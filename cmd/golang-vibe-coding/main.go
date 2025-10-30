package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"golang-vibe-coding/internal/generator"
	"golang-vibe-coding/internal/lexer"
	"golang-vibe-coding/internal/parser"
	"golang-vibe-coding/pkg/debug"
)

const version = "1.0.0"

func main() {
	// Initialize debug system
	debug.Init()

	app := &cli.App{
		Name:    "golang-vibe-coding",
		Usage:   "A Go CLI application",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:      "greet",
				Aliases:   []string{"g"},
				Usage:     "Greet someone",
				ArgsUsage: "[name]",
				Action:    greetCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "message",
						Aliases: []string{"m"},
						Value:   "Hello",
						Usage:   "Custom greeting message",
					},
				},
			},
			{
				Name:      "compile",
				Aliases:   []string{"c"},
				Usage:     "Compile DSL file to GitHub Actions YAML",
				ArgsUsage: "[input.dsl]",
				Action:    compileCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "",
						Usage:   "Output YAML file path",
					},
					&cli.BoolFlag{
						Name:    "validate",
						Aliases: []string{"v"},
						Value:   true,
						Usage:   "Validate workflow before generating",
					},
				},
			},
			{
				Name:    "example",
				Aliases: []string{"ex"},
				Usage:   "Generate example DSL file",
				Action:  exampleCommand,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("Hello, World! Welcome to your Go CLI application.")
			fmt.Println("Use --help for more options.")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func greetCommand(ctx *cli.Context) error {
	name := "World"
	if ctx.NArg() > 0 {
		name = ctx.Args().Get(0)
	}

	message := ctx.String("message")
	fmt.Printf("%s, %s!\n", message, name)
	return nil
}

func compileCommand(ctx *cli.Context) error {
	// Get input file
	if ctx.NArg() < 1 {
		return fmt.Errorf("input DSL file is required")
	}

	inputFile := ctx.Args().Get(0)

	// Read input file
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Tokenize
	l := lexer.New(string(input))

	// Parse
	p := parser.New(l)
	workflow := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parsing errors:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "  %s\n", err)
		}
		return fmt.Errorf("failed to parse DSL file")
	}

	if workflow == nil {
		return fmt.Errorf("failed to parse workflow")
	}

	// Generate YAML
	gen := generator.New(workflow)

	// Validate if requested
	if ctx.Bool("validate") {
		if validationErrors := gen.ValidateWorkflow(); len(validationErrors) > 0 {
			fmt.Fprintf(os.Stderr, "Validation errors:\n")
			for _, err := range validationErrors {
				fmt.Fprintf(os.Stderr, "  %s\n", err)
			}
			return fmt.Errorf("workflow validation failed")
		}
	}

	yamlContent, err := gen.GenerateYAML()
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// Determine output file
	outputFile := ctx.String("output")
	if outputFile == "" {
		// Default output file
		baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile = baseName + ".yml"
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputFile, []byte(yamlContent), 0o644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully compiled %s to %s\n", inputFile, outputFile)
	return nil
}

func exampleCommand(ctx *cli.Context) error {
	exampleContent := generator.GenerateExample()

	// Write to example.dsl file
	outputFile := "example.dsl"
	if err := os.WriteFile(outputFile, []byte(exampleContent), 0o644); err != nil {
		return fmt.Errorf("failed to write example file: %w", err)
	}

	fmt.Printf("Generated example DSL file: %s\n", outputFile)
	fmt.Println("You can compile it with: go run . compile example.dsl -o .github/workflows/example.yml")
	return nil
}
