#!/bin/bash

# Simple Performance Test Runner for E-commerce Microservices
# Only requires k6 and target URL - no kubernetes dependencies

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Default values
ENVIRONMENT="staging"
TEST_TYPE="load"
TARGET_URL=""
GENERATE_REPORT=true

# Function to print colored output
print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to show usage
show_usage() {
    echo "Performance Test Runner for E-commerce Microservices"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --environment ENV    Set environment (local|staging|production) [default: staging]"
    echo "  -t, --test-type TYPE     Test type (smoke|load|stress) [default: load]"
    echo "  -u, --url URL           Target URL [required]"
    echo "  -r, --report            Generate detailed report [default: true]"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Available Test Types:"
    echo "  smoke                   Quick validation (5 users, 2min)"
    echo "  load                    Load testing (up to 200 users, 10min)"
    echo "  stress                  Stress testing (300+ users, 15min)"
    echo ""
    echo "Examples:"
    echo "  $0 -u https://api.example.com                    # Default load test"
    echo "  $0 -t smoke -u https://staging.example.com      # Quick smoke test"
    echo "  $0 -t stress -e production -u https://prod.com  # Production stress test"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -t|--test-type)
            TEST_TYPE="$2"
            shift 2
            ;;
        -u|--url)
            TARGET_URL="$2"
            shift 2
            ;;
        -r|--report)
            GENERATE_REPORT=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(local|staging|production)$ ]]; then
    print_error "Invalid environment: $ENVIRONMENT"
    print_error "Valid environments: local, staging, production"
    exit 1
fi

# Validate test type
if [[ ! "$TEST_TYPE" =~ ^(smoke|load|stress)$ ]]; then
    print_error "Invalid test type: $TEST_TYPE"
    print_error "Valid test types: smoke, load, stress"
    exit 1
fi

# Validate target URL
if [[ -z "$TARGET_URL" ]]; then
    print_error "Target URL is required. Use -u or --url option."
    show_usage
    exit 1
fi

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check k6
    if ! command -v k6 &> /dev/null; then
        print_error "k6 is required for performance testing"
        print_error "Install k6: https://k6.io/docs/getting-started/installation/"
        exit 1
    fi
    
    # Get k6 version
    K6_VERSION=$(k6 version | head -n1 | cut -d' ' -f2)
    print_status "k6 version: $K6_VERSION"
    
    # Test target URL connectivity
    print_status "Testing connectivity to $TARGET_URL..."
    if curl -f -s --max-time 10 "$TARGET_URL/health" > /dev/null 2>&1; then
        print_success "Target URL is reachable"
    else
        print_warning "Health endpoint not reachable, continuing anyway..."
    fi
    
    print_success "Prerequisites check completed"
}

# Function to prepare results directory
prepare_results() {
    print_status "Preparing results directory..."
    
    mkdir -p "$RESULTS_DIR"
    
    # Clean old results (keep last 10)
    if ls "$RESULTS_DIR"/*.json 1> /dev/null 2>&1; then
        ls -t "$RESULTS_DIR"/*.json | tail -n +11 | xargs -r rm --
    fi
    
    print_success "Results directory ready: $RESULTS_DIR"
}

# Function to run performance test
run_test() {
    local test_file="tests/${TEST_TYPE}-test.js"
    local output_file="$RESULTS_DIR/${TEST_TYPE}_${ENVIRONMENT}_${TIMESTAMP}.json"
    
    print_status "Running $TEST_TYPE test..."
    print_status "Test file: $test_file"
    print_status "Target: $TARGET_URL"
    print_status "Output: $output_file"
    
    # Check if test file exists
    if [[ ! -f "$test_file" ]]; then
        print_error "Test file not found: $test_file"
        exit 1
    fi
    
    # Set environment variables for k6
    export ENVIRONMENT="$ENVIRONMENT"
    export TARGET_URL="$TARGET_URL"
    
    # Run k6 test
    print_status "Starting k6 test execution..."
    
    if k6 run \
        --out json="$output_file" \
        --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" \
        "$test_file"; then
        
        print_success "$TEST_TYPE test completed successfully"
        print_status "Results saved to: $output_file"
        
        # Show quick summary
        if [[ -f "$output_file" ]]; then
            print_status "Quick Summary:"
            echo "  - Total requests: $(grep '"type":"Point"' "$output_file" | grep '"metric":"http_reqs"' | wc -l)"
            echo "  - Test duration: $(grep '"type":"Point"' "$output_file" | grep '"metric":"iteration_duration"' | tail -1 | jq -r '.data.value' 2>/dev/null || echo "N/A") ms"
        fi
        
        return 0
    else
        print_error "$TEST_TYPE test failed"
        return 1
    fi
}

# Function to generate report
generate_report() {
    if [[ "$GENERATE_REPORT" != true ]]; then
        return
    fi
    
    print_status "Generating test report..."
    
    local output_file="$RESULTS_DIR/${TEST_TYPE}_${ENVIRONMENT}_${TIMESTAMP}.json"
    local report_file="$RESULTS_DIR/${TEST_TYPE}_${ENVIRONMENT}_${TIMESTAMP}_report.txt"
    
    if [[ ! -f "$output_file" ]]; then
        print_warning "No results file found for report generation"
        return
    fi
    
    # Generate simple text report
    {
        echo "=========================="
        echo "Performance Test Report"
        echo "=========================="
        echo "Test Type: $TEST_TYPE"
        echo "Environment: $ENVIRONMENT"
        echo "Target URL: $TARGET_URL"
        echo "Timestamp: $(date)"
        echo "=========================="
        echo ""
        echo "Results file: $output_file"
        echo ""
        echo "To analyze detailed results, use k6 analysis tools or"
        echo "import the JSON file into your preferred analysis tool."
    } > "$report_file"
    
    print_success "Report generated: $report_file"
}

# Main execution
main() {
    print_status "Starting E-commerce Performance Test"
    print_status "Environment: $ENVIRONMENT"
    print_status "Test Type: $TEST_TYPE"
    print_status "Target URL: $TARGET_URL"
    
    check_prerequisites
    prepare_results
    
    if run_test; then
        generate_report
        print_success "Performance test completed successfully!"
        exit 0
    else
        print_error "Performance test failed!"
        exit 1
    fi
}

# Run main function
main "$@" 