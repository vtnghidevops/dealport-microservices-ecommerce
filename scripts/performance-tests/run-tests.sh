#!/bin/bash

# Performance Test Runner Script for E-commerce Microservices
# This script runs k6 performance tests with proper setup and result handling
# Updated for K8s deployment with enhanced staging validation

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
ENVIRONMENT="local"
TEST_TYPE="all"
SKIP_SETUP=false
K8S_MODE=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_k8s() {
    echo -e "${PURPLE}[K8S]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --environment ENV    Set environment (local|staging|production) [default: local]"
    echo "  -t, --test-type TYPE     Test type (smoke|load|stress|pipeline|all) [default: all]"
    echo "  -s, --skip-setup         Skip environment setup checks"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                       # Run all tests on local environment"
    echo "  $0 -t smoke              # Run only smoke test"
    echo "  $0 -e staging -t load    # Run load test on staging"
    echo "  $0 -t pipeline           # Run pipeline test for CI/CD"
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
        -s|--skip-setup)
            SKIP_SETUP=true
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
if [[ ! "$TEST_TYPE" =~ ^(smoke|load|stress|pipeline|all)$ ]]; then
    print_error "Invalid test type: $TEST_TYPE"
    print_error "Valid test types: smoke, load, stress, pipeline, all"
    exit 1
fi

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if k6 is installed
    if ! command -v k6 &> /dev/null; then
        print_error "k6 is not installed. Please install k6 first."
        print_status "Installation instructions: https://k6.io/docs/getting-started/installation/"
        exit 1
    fi
    
    # Check k6 version
    K6_VERSION=$(k6 version | head -n1 | cut -d' ' -f2)
    print_status "k6 version: $K6_VERSION"
    
    # Create results directory
    mkdir -p "$RESULTS_DIR"
    
    print_success "Prerequisites check completed"
}

# Function to check environment connectivity
check_environment() {
    if [[ "$SKIP_SETUP" == true ]]; then
        print_warning "Skipping environment setup checks"
        return
    fi
    
    print_status "Checking environment connectivity..."
    
    case $ENVIRONMENT in
        "local")
            BASE_URL="http://localhost:58080"
            ;;
        "staging")
            BASE_URL="https://staging-api.ecommerce.com"
            ;;
        "production")
            BASE_URL="https://api.ecommerce.com"
            ;;
    esac
    
    # Check if the API is accessible
    if curl -s --max-time 10 "$BASE_URL/api/v1/health" > /dev/null 2>&1; then
        print_success "Environment $ENVIRONMENT is accessible at $BASE_URL"
    else
        print_warning "Cannot reach $BASE_URL - tests may fail"
        print_warning "Make sure the services are running for $ENVIRONMENT environment"
        
        if [[ "$ENVIRONMENT" == "local" ]]; then
            print_status "For local environment, run: docker-compose up -d"
        fi
        
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Function to run a specific test
run_test() {
    local test_name=$1
    local test_file="$SCRIPT_DIR/tests/${test_name}-test.js"
    local result_file="$RESULTS_DIR/${test_name}_${ENVIRONMENT}_${TIMESTAMP}.json"
    
    if [[ ! -f "$test_file" ]]; then
        print_error "Test file not found: $test_file"
        return 1
    fi
    
    print_status "Running $test_name test..."
    print_status "Environment: $ENVIRONMENT"
    print_status "Results will be saved to: $result_file"
    
    # Run k6 test with environment variable and JSON output
    if ENVIRONMENT="$ENVIRONMENT" k6 run \
        --out json="$result_file" \
        --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" \
        "$test_file"; then
        
        print_success "$test_name test completed successfully"
        
        # Generate summary report
        generate_summary_report "$test_name" "$result_file"
        
        return 0
    else
        print_error "$test_name test failed"
        return 1
    fi
}

# Function to generate summary report
generate_summary_report() {
    local test_name=$1
    local result_file=$2
    local summary_file="$RESULTS_DIR/${test_name}_${ENVIRONMENT}_${TIMESTAMP}_summary.txt"
    
    print_status "Generating summary report..."
    
    {
        echo "=========================================="
        echo "Performance Test Summary Report"
        echo "=========================================="
        echo "Test Type: $test_name"
        echo "Environment: $ENVIRONMENT"
        echo "Timestamp: $(date)"
        echo "=========================================="
        echo ""
        
        # Extract key metrics from JSON (basic parsing)
        if [[ -f "$result_file" ]]; then
            echo "Key Metrics:"
            echo "- Check the detailed JSON results in: $result_file"
            echo "- Use k6 dashboard or analysis tools for detailed insights"
        fi
        
        echo ""
        echo "Next Steps:"
        echo "1. Review the detailed results in the JSON file"
        echo "2. Compare with previous test runs"
        echo "3. Investigate any performance regressions"
        echo "4. Update performance baselines if needed"
        
    } > "$summary_file"
    
    print_success "Summary report saved to: $summary_file"
}

# Function to run all tests
run_all_tests() {
    local failed_tests=()
    
    print_status "Running all performance tests..."
    
    # Run tests in order of increasing load
    for test in smoke load stress; do
        if run_test "$test"; then
            print_success "$test test passed"
        else
            print_error "$test test failed"
            failed_tests+=("$test")
        fi
        
        # Wait between tests to allow system recovery
        if [[ "$test" != "stress" ]]; then
            print_status "Waiting 30 seconds for system recovery..."
            sleep 30
        fi
    done
    
    # Report results
    if [[ ${#failed_tests[@]} -eq 0 ]]; then
        print_success "All performance tests completed successfully!"
    else
        print_error "The following tests failed: ${failed_tests[*]}"
        return 1
    fi
}

# Main execution
main() {
    echo "=========================================="
    echo "E-commerce Microservices Performance Tests"
    echo "=========================================="
    echo ""
    
    check_prerequisites
    check_environment
    
    case $TEST_TYPE in
        "smoke"|"load"|"stress"|"pipeline")
            run_test "$TEST_TYPE"
            ;;
        "all")
            run_all_tests
            ;;
    esac
    
    echo ""
    print_success "Performance testing completed!"
    print_status "Results are available in: $RESULTS_DIR"
    
    # Show recent results
    echo ""
    print_status "Recent test results:"
    ls -la "$RESULTS_DIR" | tail -5
}

# Run main function
main "$@" 