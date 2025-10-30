#!/usr/bin/env pwsh
# Build script for golang-vibe-coding project

param(
    [string]$Target = "all",
    [string]$OS = "",
    [string]$Arch = "",
    [switch]$Clean = $false,
    [switch]$Verbose = $false
)

$ErrorActionPreference = "Stop"

# Project configuration
$ProjectName = "golang-vibe-coding"
$BinDir = "bin"
$Commands = @("golang-vibe-coding", "debugger", "dap")

# Colors for output
function Write-Success { param($Message) Write-Host "✅ $Message" -ForegroundColor Green }
function Write-Info { param($Message) Write-Host "ℹ️ $Message" -ForegroundColor Blue }
function Write-Warning { param($Message) Write-Host "⚠️ $Message" -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host "❌ $Message" -ForegroundColor Red }

function Show-Help {
    Write-Host @"
Build script for $ProjectName

Usage: .\scripts\build.ps1 [OPTIONS]

OPTIONS:
    -Target <target>     Build target: all, main, debugger, dap, clean (default: all)
    -OS <os>            Target OS: windows, linux, darwin (default: current)
    -Arch <arch>        Target architecture: amd64, arm64 (default: current)
    -Clean              Clean build artifacts before building
    -Verbose            Enable verbose output
    -Help               Show this help message

EXAMPLES:
    .\scripts\build.ps1                          # Build all targets for current platform
    .\scripts\build.ps1 -Target main             # Build only main CLI
    .\scripts\build.ps1 -OS linux -Arch amd64    # Cross-compile for Linux
    .\scripts\build.ps1 -Clean                   # Clean and rebuild all
"@
}

function Clean-Artifacts {
    Write-Info "Cleaning build artifacts..."
    
    if (Test-Path $BinDir) {
        Remove-Item "$BinDir/*" -Force -Recurse -ErrorAction SilentlyContinue
        Write-Success "Cleaned $BinDir directory"
    }
}

function Build-Binary {
    param(
        [string]$Command,
        [string]$OutputName,
        [string]$TargetOS,
        [string]$TargetArch
    )
    
    $env:GOOS = $TargetOS
    $env:GOARCH = $TargetArch
    
    $suffix = ""
    if ($TargetOS -eq "windows") {
        $suffix = ".exe"
    }
    
    $outputPath = "$BinDir/$OutputName$suffix"
    $sourceDir = "./cmd/$Command"
    
    Write-Info "Building $Command for $TargetOS/$TargetArch..."
    
    if ($Verbose) {
        Write-Host "Command: go build -o $outputPath $sourceDir"
    }
    
    try {
        go build -o $outputPath $sourceDir
        if ($LASTEXITCODE -eq 0) {
            $fileSize = (Get-Item $outputPath).Length
            $fileSizeKB = [math]::Round($fileSize / 1024, 2)
            Write-Success "Built $outputPath ($fileSizeKB KB)"
        } else {
            Write-Error "Failed to build $Command"
            exit 1
        }
    } catch {
        Write-Error "Build failed for $Command`: $($_.Exception.Message)"
        exit 1
    }
}

function Build-All {
    param([string]$TargetOS, [string]$TargetArch)
    
    # Ensure bin directory exists
    if (-not (Test-Path $BinDir)) {
        New-Item -ItemType Directory -Path $BinDir | Out-Null
    }
    
    # Build each command
    foreach ($cmd in $Commands) {
        $outputName = $cmd
        if ($cmd -eq "golang-vibe-coding") {
            $outputName = $ProjectName
        }
        
        Build-Binary -Command $cmd -OutputName $outputName -TargetOS $TargetOS -TargetArch $TargetArch
    }
}

# Main execution
Write-Info "Starting build process for $ProjectName"

# Show help if requested
if ($args -contains "-Help" -or $args -contains "--help" -or $args -contains "-h") {
    Show-Help
    exit 0
}

# Determine target platform
$currentOS = if ($IsWindows) { "windows" } elseif ($IsLinux) { "linux" } elseif ($IsMacOS) { "darwin" } else { "windows" }
$currentArch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

$targetOS = if ($OS) { $OS } else { $currentOS }
$targetArch = if ($Arch) { $Arch } else { $currentArch }

Write-Info "Target platform: $targetOS/$targetArch"

# Clean if requested
if ($Clean) {
    Clean-Artifacts
}

# Execute build target
switch ($Target.ToLower()) {
    "clean" {
        Clean-Artifacts
        Write-Success "Clean completed"
    }
    "main" {
        if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir | Out-Null }
        Build-Binary -Command "golang-vibe-coding" -OutputName $ProjectName -TargetOS $targetOS -TargetArch $targetArch
    }
    "debugger" {
        if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir | Out-Null }
        Build-Binary -Command "debugger" -OutputName "debugger" -TargetOS $targetOS -TargetArch $targetArch
    }
    "dap" {
        if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir | Out-Null }
        Build-Binary -Command "dap" -OutputName "dap" -TargetOS $targetOS -TargetArch $targetArch
    }
    "all" {
        Build-All -TargetOS $targetOS -TargetArch $targetArch
    }
    default {
        Write-Error "Unknown target: $Target"
        Write-Info "Valid targets: all, main, debugger, dap, clean"
        exit 1
    }
}

Write-Success "Build process completed successfully!"

# Show built artifacts
Write-Info "Built artifacts:"
if (Test-Path $BinDir) {
    Get-ChildItem $BinDir | ForEach-Object {
        $size = [math]::Round($_.Length / 1024, 2)
        Write-Host "  $($_.Name) ($size KB)" -ForegroundColor Cyan
    }
}
