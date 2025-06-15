# Performance Test Runner Script for E-commerce Microservices (PowerShell)
# This script runs k6 performance tests with proper setup and result handling

param(
    [string]$Environment = "local",
    [string]$TestType = "all",
    [switch]$SkipSetup,
    [switch]$Help
)

# Configuration
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ResultsDir = Join-Path $ScriptDir "results"
$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"

# Function to show usage
function Show-Usage {
    Write-Host "Usage: .\run-tests.ps1 [OPTIONS]" -ForegroundColor White
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Environment ENV     Set environment (local|staging|production) [default: local]"
    Write-Host "  -TestType TYPE       Test type (smoke|load|stress|pipeline|all) [default: all]"
    Write-Host "  -SkipSetup          Skip environment setup checks"
    Write-Host "  -Help               Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Green
    Write-Host "  .\run-tests.ps1                          # Run all tests on local environment"
    Write-Host "  .\run-tests.ps1 -TestType smoke          # Run only smoke test"
    Write-Host "  .\run-tests.ps1 -Environment staging -TestType load  # Run load test on staging"
    Write-Host "  .\run-tests.ps1 -TestType pipeline       # Run pipeline test for CI/CD"
}

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Show help if requested
if ($Help) {
    Show-Usage
    exit 0
}

# Validate environment
if ($Environment -notin @("local", "staging", "production")) {
    Write-Error "Invalid environment: $Environment"
    Write-Error "Valid environments: local, staging, production"
    exit 1
}

# Validate test type
if ($TestType -notin @("smoke", "load", "stress", "pipeline", "all")) {
    Write-Error "Invalid test type: $TestType"
    Write-Error "Valid test types: smoke, load, stress, pipeline, all"
    exit 1
}

# Function to check prerequisites
function Test-Prerequisites {
    Write-Status "Checking prerequisites..."
    
    # Check if k6 is installed
    try {
        $k6Version = & k6 version 2>$null
        if ($LASTEXITCODE -ne 0) {
            throw "k6 not found"
        }
        Write-Status "k6 version: $($k6Version[0])"
    }
    catch {
        Write-Error "k6 is not installed. Please install k6 first."
        Write-Status "Installation instructions: https://k6.io/docs/getting-started/installation/"
        exit 1
    }
    
    # Create results directory
    if (!(Test-Path $ResultsDir)) {
        New-Item -ItemType Directory -Path $ResultsDir -Force | Out-Null
    }
    
    Write-Success "Prerequisites check completed"
}

# Function to check environment connectivity
function Test-Environment {
    if ($SkipSetup) {
        Write-Warning "Skipping environment setup checks"
        return
    }
    
    Write-Status "Checking environment connectivity..."
    
    switch ($Environment) {
        "local" { $BaseUrl = "http://localhost:58080" }
        "staging" { $BaseUrl = "https://staging-api.ecommerce.com" }
        "production" { $BaseUrl = "https://api.ecommerce.com" }
    }
    
    # Check if the API is accessible
    try {
        $response = Invoke-WebRequest -Uri "$BaseUrl/api/v1/health" -TimeoutSec 10 -UseBasicParsing -ErrorAction Stop
        Write-Success "Environment $Environment is accessible at $BaseUrl"
    }
    catch {
        Write-Warning "Cannot reach $BaseUrl - tests may fail"
        Write-Warning "Make sure the services are running for $Environment environment"
        
        if ($Environment -eq "local") {
            Write-Status "For local environment, run: docker-compose up -d"
        }
        
        $continue = Read-Host "Continue anyway? (y/N)"
        if ($continue -ne "y" -and $continue -ne "Y") {
            exit 1
        }
    }
}

# Function to run a specific test
function Invoke-Test {
    param([string]$TestName)
    
    $TestFile = Join-Path $ScriptDir "tests\$TestName-test.js"
    $ResultFile = Join-Path $ResultsDir "$TestName`_$Environment`_$Timestamp.json"
    
    if (!(Test-Path $TestFile)) {
        Write-Error "Test file not found: $TestFile"
        return $false
    }
    
    Write-Status "Running $TestName test..."
    Write-Status "Environment: $Environment"
    Write-Status "Results will be saved to: $ResultFile"
    
    # Set environment variable and run k6 test
    $env:ENVIRONMENT = $Environment
    
    try {
        & k6 run --out "json=$ResultFile" --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" $TestFile
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "$TestName test completed successfully"
            New-SummaryReport -TestName $TestName -ResultFile $ResultFile
            return $true
        }
        else {
            Write-Error "$TestName test failed"
            return $false
        }
    }
    catch {
        Write-Error "$TestName test failed with exception: $_"
        return $false
    }
}

# Function to generate summary report
function New-SummaryReport {
    param(
        [string]$TestName,
        [string]$ResultFile
    )
    
    $SummaryFile = Join-Path $ResultsDir "$TestName`_$Environment`_$Timestamp`_summary.txt"
    
    Write-Status "Generating summary report..."
    
    $summaryContent = @"
==========================================
Performance Test Summary Report
==========================================
Test Type: $TestName
Environment: $Environment
Timestamp: $(Get-Date)
==========================================

Key Metrics:
- Check the detailed JSON results in: $ResultFile
- Use k6 dashboard or analysis tools for detailed insights

Next Steps:
1. Review the detailed results in the JSON file
2. Compare with previous test runs
3. Investigate any performance regressions
4. Update performance baselines if needed
"@
    
    $summaryContent | Out-File -FilePath $SummaryFile -Encoding UTF8
    Write-Success "Summary report saved to: $SummaryFile"
}

# Function to run all tests
function Invoke-AllTests {
    $failedTests = @()
    
    Write-Status "Running all performance tests..."
    
    # Run tests in order of increasing load
    $tests = @("smoke", "load", "stress")
    
    foreach ($test in $tests) {
        if (Invoke-Test -TestName $test) {
            Write-Success "$test test passed"
        }
        else {
            Write-Error "$test test failed"
            $failedTests += $test
        }
        
        # Wait between tests to allow system recovery
        if ($test -ne "stress") {
            Write-Status "Waiting 30 seconds for system recovery..."
            Start-Sleep -Seconds 30
        }
    }
    
    # Report results
    if ($failedTests.Count -eq 0) {
        Write-Success "All performance tests completed successfully!"
        return $true
    }
    else {
        Write-Error "The following tests failed: $($failedTests -join ', ')"
        return $false
    }
}

# Main execution
function Main {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "E-commerce Microservices Performance Tests" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""
    
    Test-Prerequisites
    Test-Environment
    
    $success = $false
    
    switch ($TestType) {
        { $_ -in @("smoke", "load", "stress", "pipeline") } {
            $success = Invoke-Test -TestName $TestType
        }
        "all" {
            $success = Invoke-AllTests
        }
    }
    
    Write-Host ""
    if ($success) {
        Write-Success "Performance testing completed successfully!"
    }
    else {
        Write-Error "Performance testing completed with failures!"
    }
    
    Write-Status "Results are available in: $ResultsDir"
    
    # Show recent results
    Write-Host ""
    Write-Status "Recent test results:"
    Get-ChildItem $ResultsDir | Sort-Object LastWriteTime -Descending | Select-Object -First 5 | Format-Table Name, LastWriteTime, Length
}

# Run main function
Main 