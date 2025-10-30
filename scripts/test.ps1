#!/usr/bin/env pwsh
# Test script for golang-vibe-coding project

param(
    [string]$Target = "all",
    [switch]$Coverage = $false,
    [switch]$Verbose = $false,
    [switch]$Race = $false,
    [string]$Pattern = ""
)

$ErrorActionPreference = "Stop"

# Colors for output
function Write-Success { param($Message) Write-Host "✅ $Message" -ForegroundColor Green }
function Write-Info { param($Message) Write-Host "ℹ️ $Message" -ForegroundColor Blue }
function Write-Warning { param($Message) Write-Host "⚠️ $Message" -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host "❌ $Message" -ForegroundColor Red }

function Show-Help {
    Write-Host @"
Test script for golang-vibe-coding

Usage: .\scripts\test.ps1 [OPTIONS]

OPTIONS:
    -Target <target>     Test target: all, unit, integration (default: all)
    -Coverage           Generate coverage report
    -Verbose            Enable verbose test output
    -Race              Enable race condition detection
    -Pattern <pattern>  Run tests matching pattern
    -Help              Show this help message

EXAMPLES:
    .\scripts\test.ps1                        # Run all tests
    .\scripts\test.ps1 -Coverage             # Run tests with coverage
    .\scripts\test.ps1 -Pattern "*Parser*"   # Run parser tests only
    .\scripts\test.ps1 -Race -Verbose       # Run with race detection
"@
}

function Invoke-Tests {
    param(
        [string]$TestPattern,
        [bool]$EnableCoverage,
        [bool]$EnableVerbose,
        [bool]$EnableRace
    )
    
    Write-Info "Running Go tests..."
    
    $testArgs = @("test")
    
    # Add test pattern if specified
    if ($TestPattern) {
        $testArgs += "./..."
        $env:TESTPATTERN = $TestPattern
    } else {
        $testArgs += "./..."
    }
    
    # Add flags
    if ($EnableVerbose) {
        $testArgs += "-v"
    }
    
    if ($EnableRace) {
        $testArgs += "-race"
    }
    
    if ($EnableCoverage) {
        $testArgs += "-coverprofile=coverage.out"
        $testArgs += "-covermode=atomic"
    }
    
    Write-Info "Command: go $($testArgs -join ' ')"
    
    try {
        & go @testArgs
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "All tests passed!"
            
            if ($EnableCoverage -and (Test-Path "coverage.out")) {
                Write-Info "Generating coverage report..."
                & go tool cover -html=coverage.out -o coverage.html
                
                Write-Info "Getting coverage summary..."
                $coverageOutput = & go tool cover -func=coverage.out
                Write-Host $coverageOutput
                
                Write-Success "Coverage report generated: coverage.html"
            }
        } else {
            Write-Error "Tests failed with exit code $LASTEXITCODE"
            exit 1
        }
    } catch {
        Write-Error "Test execution failed: $($_.Exception.Message)"
        exit 1
    }
}

function Test-CodeLinting {
    Write-Info "Running code linting..."
    
    # Check if golangci-lint is available
    $lintCmd = Get-Command golangci-lint -ErrorAction SilentlyContinue
    
    if ($lintCmd) {
        try {
            & golangci-lint run
            if ($LASTEXITCODE -eq 0) {
                Write-Success "Linting passed!"
            } else {
                Write-Warning "Linting found issues (exit code $LASTEXITCODE)"
            }
        } catch {
            Write-Warning "Linting failed: $($_.Exception.Message)"
        }
    } else {
        Write-Warning "golangci-lint not found, skipping linting"
        Write-Info "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    }
}

function Invoke-ModTidy {
    Write-Info "Running go mod tidy..."
    
    try {
        & go mod tidy
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Go modules tidied"
        } else {
            Write-Error "go mod tidy failed"
            exit 1
        }
    } catch {
        Write-Error "go mod tidy failed: $($_.Exception.Message)"
        exit 1
    }
}

# Main execution
Write-Info "Starting test process for golang-vibe-coding"

# Show help if requested
if ($args -contains "-Help" -or $args -contains "--help" -or $args -contains "-h") {
    Show-Help
    exit 0
}

# Ensure we're in the project root
if (-not (Test-Path "go.mod")) {
    Write-Error "go.mod not found. Please run from project root."
    exit 1
}

# Run go mod tidy first
Invoke-ModTidy

# Execute test target
switch ($Target.ToLower()) {
    "unit" {
        Invoke-Tests -TestPattern $Pattern -EnableCoverage $Coverage -EnableVerbose $Verbose -EnableRace $Race
    }
    "integration" {
        Write-Warning "Integration tests not yet implemented"
    }
    "lint" {
        Test-CodeLinting
    }
    "all" {
        Invoke-Tests -TestPattern $Pattern -EnableCoverage $Coverage -EnableVerbose $Verbose -EnableRace $Race
        Test-CodeLinting
    }
    default {
        Write-Error "Unknown target: $Target"
        Write-Info "Valid targets: all, unit, integration, lint"
        exit 1
    }
}

Write-Success "Test process completed successfully!"
