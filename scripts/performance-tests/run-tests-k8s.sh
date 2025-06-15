#!/bin/bash

# K8s Performance Test Runner for E-commerce Microservices
# Specialized script for K8s environments with enhanced monitoring and staging validation

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Default values
ENVIRONMENT="staging"
TEST_TYPE="staging"
NAMESPACE="default"
MONITOR_RESOURCES=true
GENERATE_REPORT=true

# Function to print colored output
print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_k8s() { echo -e "${PURPLE}[K8S]${NC} $1"; }
print_monitor() { echo -e "${CYAN}[MONITOR]${NC} $1"; }

# Function to show usage
show_usage() {
    echo "K8s Performance Test Runner for E-commerce Microservices"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --environment ENV    Set environment (staging|production) [default: staging]"
    echo "  -t, --test-type TYPE     Test type (smoke|load|stress|staging|pipeline) [default: staging]"
    echo "  -n, --namespace NS       K8s namespace [default: default]"
    echo "  -m, --monitor           Enable resource monitoring [default: true]"
    echo "  -r, --report            Generate detailed K8s report [default: true]"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "K8s Optimized Test Types:"
    echo "  staging                 Comprehensive staging validation (200+ users, 20min)"
    echo "  stress                  K8s stress testing (300+ users with spikes)"
    echo "  load                    K8s load testing (up to 200 users)"
    echo "  smoke                   Quick K8s validation (5 users)"
    echo "  pipeline                CI/CD pipeline validation"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Full staging validation"
    echo "  $0 -t stress -e staging              # Stress test on staging"
    echo "  $0 -t load -n ecommerce              # Load test in ecommerce namespace"
    echo "  $0 -t pipeline -e staging            # Pipeline validation"
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
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -m|--monitor)
            MONITOR_RESOURCES=true
            shift
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
if [[ ! "$ENVIRONMENT" =~ ^(staging|production)$ ]]; then
    print_error "Invalid environment: $ENVIRONMENT"
    print_error "Valid environments for K8s: staging, production"
    exit 1
fi

# Validate test type
if [[ ! "$TEST_TYPE" =~ ^(smoke|load|stress|staging|pipeline)$ ]]; then
    print_error "Invalid test type: $TEST_TYPE"
    print_error "Valid test types: smoke, load, stress, staging, pipeline"
    exit 1
fi

# Function to check K8s prerequisites
check_k8s_prerequisites() {
    print_k8s "Checking K8s prerequisites..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is required for K8s performance testing"
        exit 1
    fi
    
    # Check k6
    if ! command -v k6 &> /dev/null; then
        print_error "k6 is required for performance testing"
        exit 1
    fi
    
    # Check cluster connectivity
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to K8s cluster"
        exit 1
    fi
    
    # Get cluster info
    KUBECTL_VERSION=$(kubectl version --client --short 2>/dev/null | cut -d' ' -f3)
    K6_VERSION=$(k6 version | head -n1 | cut -d' ' -f2)
    
    print_k8s "kubectl version: $KUBECTL_VERSION"
    print_k8s "k6 version: $K6_VERSION"
    
    # Check namespace
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        print_warning "Namespace '$NAMESPACE' not found, using default"
        NAMESPACE="default"
    fi
    
    print_success "K8s prerequisites check completed"
}

# Function to get cluster status
get_cluster_status() {
    print_k8s "Getting cluster status..."
    
    # Node information
    NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
    READY_NODES=$(kubectl get nodes --no-headers | grep -c " Ready ")
    
    print_k8s "📊 Cluster: $READY_NODES/$NODE_COUNT nodes ready"
    
    # Pod information
    TOTAL_PODS=$(kubectl get pods -n "$NAMESPACE" --no-headers | wc -l)
    RUNNING_PODS=$(kubectl get pods -n "$NAMESPACE" --field-selector=status.phase=Running --no-headers | wc -l)
    
    print_k8s "🚀 Namespace '$NAMESPACE': $RUNNING_PODS/$TOTAL_PODS pods running"
    
    # Check for HPA
    if kubectl get hpa -n "$NAMESPACE" &> /dev/null; then
        HPA_COUNT=$(kubectl get hpa -n "$NAMESPACE" --no-headers | wc -l)
        print_k8s "📈 $HPA_COUNT HPA configurations found"
    fi
    
    # Check for resource quotas
    if kubectl get resourcequota -n "$NAMESPACE" &> /dev/null; then
        QUOTA_COUNT=$(kubectl get resourcequota -n "$NAMESPACE" --no-headers | wc -l)
        print_k8s "💾 $QUOTA_COUNT resource quotas configured"
    fi
}

# Function to start resource monitoring
start_monitoring() {
    if [[ "$MONITOR_RESOURCES" != true ]]; then
        return
    fi
    
    print_monitor "Starting resource monitoring..."
    
    local monitor_file="$RESULTS_DIR/k8s_monitoring_${TIMESTAMP}.log"
    
    # Start background monitoring
    {
        echo "=== K8s Resource Monitoring Started at $(date) ==="
        echo ""
        
        while true; do
            echo "--- $(date) ---"
            
            # Node resources
            echo "Node Resources:"
            kubectl top nodes --no-headers 2>/dev/null || echo "Metrics server not available"
            echo ""
            
            # Pod resources in namespace
            echo "Pod Resources (namespace: $NAMESPACE):"
            kubectl top pods -n "$NAMESPACE" --no-headers 2>/dev/null || echo "Pod metrics not available"
            echo ""
            
            # HPA status
            if kubectl get hpa -n "$NAMESPACE" &> /dev/null; then
                echo "HPA Status:"
                kubectl get hpa -n "$NAMESPACE" --no-headers
                echo ""
            fi
            
            # Pod status
            echo "Pod Status:"
            kubectl get pods -n "$NAMESPACE" --no-headers | awk '{print $1, $3, $4}'
            echo ""
            
            echo "=================================="
            echo ""
            
            sleep 30
        done
    } > "$monitor_file" 2>&1 &
    
    MONITOR_PID=$!
    echo $MONITOR_PID > "$RESULTS_DIR/monitor.pid"
    
    print_monitor "Resource monitoring started (PID: $MONITOR_PID)"
    print_monitor "Monitoring log: $monitor_file"
}

# Function to stop resource monitoring
stop_monitoring() {
    if [[ "$MONITOR_RESOURCES" != true ]]; then
        return
    fi
    
    local pid_file="$RESULTS_DIR/monitor.pid"
    
    if [[ -f "$pid_file" ]]; then
        local monitor_pid=$(cat "$pid_file")
        if kill -0 "$monitor_pid" 2>/dev/null; then
            kill "$monitor_pid"
            print_monitor "Resource monitoring stopped"
        fi
        rm -f "$pid_file"
    fi
}

# Function to run performance test
run_k8s_test() {
    local test_name=$1
    local test_file="$SCRIPT_DIR/tests/${test_name}-test.js"
    local result_file="$RESULTS_DIR/${test_name}_k8s_${ENVIRONMENT}_${TIMESTAMP}.json"
    
    # Handle special case for staging deployment test
    if [[ "$test_name" == "staging" ]]; then
        test_file="$SCRIPT_DIR/tests/staging-deployment-test.js"
    fi
    
    if [[ ! -f "$test_file" ]]; then
        print_error "Test file not found: $test_file"
        return 1
    fi
    
    print_k8s "🚀 Running K8s $test_name test..."
    print_k8s "Environment: $ENVIRONMENT"
    print_k8s "Namespace: $NAMESPACE"
    print_k8s "Results: $result_file"
    
    # Set environment variables
    export ENVIRONMENT="$ENVIRONMENT"
    export K8S_MODE="true"
    export K8S_NAMESPACE="$NAMESPACE"
    
    # Start monitoring
    start_monitoring
    
    # Record pre-test state
    local pre_test_file="$RESULTS_DIR/pre_test_state_${TIMESTAMP}.log"
    {
        echo "=== Pre-test K8s State ==="
        echo "Timestamp: $(date)"
        echo "Test: $test_name"
        echo "Environment: $ENVIRONMENT"
        echo "Namespace: $NAMESPACE"
        echo ""
        
        echo "Nodes:"
        kubectl get nodes
        echo ""
        
        echo "Pods in namespace $NAMESPACE:"
        kubectl get pods -n "$NAMESPACE"
        echo ""
        
        if kubectl get hpa -n "$NAMESPACE" &> /dev/null; then
            echo "HPA Status:"
            kubectl get hpa -n "$NAMESPACE"
            echo ""
        fi
        
    } > "$pre_test_file"
    
    # Run the test
    local success=false
    if ENVIRONMENT="$ENVIRONMENT" k6 run \
        --out json="$result_file" \
        --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" \
        "$test_file"; then
        
        print_success "K8s $test_name test completed successfully"
        success=true
    else
        print_error "K8s $test_name test failed"
    fi
    
    # Stop monitoring
    stop_monitoring
    
    # Record post-test state
    local post_test_file="$RESULTS_DIR/post_test_state_${TIMESTAMP}.log"
    {
        echo "=== Post-test K8s State ==="
        echo "Timestamp: $(date)"
        echo "Test: $test_name"
        echo "Success: $success"
        echo ""
        
        echo "Nodes:"
        kubectl get nodes
        echo ""
        
        echo "Pods in namespace $NAMESPACE:"
        kubectl get pods -n "$NAMESPACE"
        echo ""
        
        echo "Pod Restart Counts:"
        kubectl get pods -n "$NAMESPACE" -o custom-columns=NAME:.metadata.name,RESTARTS:.status.containerStatuses[0].restartCount
        echo ""
        
        if kubectl get hpa -n "$NAMESPACE" &> /dev/null; then
            echo "HPA Status:"
            kubectl get hpa -n "$NAMESPACE"
            echo ""
        fi
        
        echo "Events (last 10):"
        kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -10
        echo ""
        
    } > "$post_test_file"
    
    # Generate K8s report
    if [[ "$GENERATE_REPORT" == true ]]; then
        generate_k8s_report "$test_name" "$result_file" "$pre_test_file" "$post_test_file"
    fi
    
    if [[ "$success" == true ]]; then
        return 0
    else
        return 1
    fi
}

# Function to generate K8s report
generate_k8s_report() {
    local test_name=$1
    local result_file=$2
    local pre_test_file=$3
    local post_test_file=$4
    local report_file="$RESULTS_DIR/k8s_report_${test_name}_${TIMESTAMP}.md"
    
    print_k8s "📊 Generating K8s performance report..."
    
    {
        echo "# K8s Performance Test Report"
        echo ""
        echo "## Test Information"
        echo "- **Test Type**: $test_name"
        echo "- **Environment**: $ENVIRONMENT"
        echo "- **Namespace**: $NAMESPACE"
        echo "- **Timestamp**: $(date)"
        echo "- **Duration**: $(get_test_duration "$result_file")"
        echo ""
        
        echo "## Cluster Information"
        echo "- **Nodes**: $NODE_COUNT total, $READY_NODES ready"
        echo "- **Pods**: $RUNNING_PODS/$TOTAL_PODS running in namespace"
        if [[ -n "$HPA_COUNT" ]]; then
            echo "- **HPA**: $HPA_COUNT configurations"
        fi
        echo ""
        
        echo "## Test Results Summary"
        echo "- **Result File**: \`$result_file\`"
        echo "- **Pre-test State**: \`$pre_test_file\`"
        echo "- **Post-test State**: \`$post_test_file\`"
        if [[ "$MONITOR_RESOURCES" == true ]]; then
            echo "- **Monitoring Log**: \`$RESULTS_DIR/k8s_monitoring_${TIMESTAMP}.log\`"
        fi
        echo ""
        
        echo "## Performance Metrics"
        echo "Detailed performance metrics are available in the JSON result file."
        echo "Key areas to analyze:"
        echo "- Response time percentiles (p95, p99)"
        echo "- Error rates and patterns"
        echo "- Throughput and concurrent user handling"
        echo "- Resource utilization during load"
        echo ""
        
        echo "## K8s Specific Observations"
        echo ""
        echo "### Pod Stability"
        local restart_count=$(kubectl get pods -n "$NAMESPACE" -o jsonpath='{.items[*].status.containerStatuses[0].restartCount}' | tr ' ' '\n' | awk '{sum+=$1} END {print sum+0}')
        echo "- Total pod restarts during test: $restart_count"
        
        echo ""
        echo "### Auto-scaling Behavior"
        if kubectl get hpa -n "$NAMESPACE" &> /dev/null; then
            echo "HPA configurations were active during the test."
            echo "Review the monitoring logs for scaling events."
        else
            echo "No HPA configurations detected."
        fi
        
        echo ""
        echo "### Resource Utilization"
        echo "Check the monitoring logs for detailed resource usage patterns:"
        echo "- CPU utilization across nodes"
        echo "- Memory usage patterns"
        echo "- Network throughput"
        echo "- Storage I/O if applicable"
        
        echo ""
        echo "## Recommendations"
        echo ""
        
        case $test_name in
            "staging")
                echo "### Staging Deployment Validation"
                echo "- ✅ Review response time thresholds"
                echo "- ✅ Validate error rates under load"
                echo "- ✅ Check auto-scaling effectiveness"
                echo "- ✅ Verify resource limits and requests"
                echo "- ✅ Confirm system recovery after spikes"
                echo ""
                echo "**Production Readiness**: If all metrics pass thresholds, the system is ready for production deployment."
                ;;
            "stress")
                echo "### Stress Test Analysis"
                echo "- 🔍 Identify breaking points and bottlenecks"
                echo "- 🔍 Review system behavior under extreme load"
                echo "- 🔍 Validate circuit breaker and failover mechanisms"
                echo "- 🔍 Plan capacity and scaling strategies"
                ;;
            "load")
                echo "### Load Test Analysis"
                echo "- 📊 Confirm system handles expected traffic"
                echo "- 📊 Validate performance under normal conditions"
                echo "- 📊 Review resource efficiency"
                ;;
        esac
        
        echo ""
        echo "## Next Steps"
        echo "1. Analyze detailed performance metrics in JSON results"
        echo "2. Review K8s monitoring logs for resource patterns"
        echo "3. Compare results with previous test runs"
        echo "4. Update performance baselines if needed"
        echo "5. Plan infrastructure optimizations based on findings"
        
        if [[ "$test_name" == "staging" ]] && [[ "$ENVIRONMENT" == "staging" ]]; then
            echo "6. **Proceed with production deployment if all criteria are met**"
        fi
        
    } > "$report_file"
    
    print_success "K8s report generated: $report_file"
}

# Function to get test duration (placeholder)
get_test_duration() {
    local result_file=$1
    # This would parse the JSON file to get actual duration
    # For now, return a placeholder
    echo "Check JSON results for duration"
}

# Function to cleanup
cleanup() {
    print_k8s "Cleaning up..."
    stop_monitoring
    
    # Remove any temporary files
    rm -f "$RESULTS_DIR/monitor.pid"
}

# Trap cleanup on exit
trap cleanup EXIT

# Main execution
main() {
    echo "=========================================="
    echo "🚀 K8s Performance Testing Suite"
    echo "E-commerce Microservices"
    echo "=========================================="
    echo ""
    
    # Create results directory
    mkdir -p "$RESULTS_DIR"
    
    check_k8s_prerequisites
    get_cluster_status
    
    echo ""
    print_k8s "🎯 Test Configuration:"
    print_k8s "Environment: $ENVIRONMENT"
    print_k8s "Test Type: $TEST_TYPE"
    print_k8s "Namespace: $NAMESPACE"
    print_k8s "Monitoring: $MONITOR_RESOURCES"
    print_k8s "Reporting: $GENERATE_REPORT"
    echo ""
    
    # Run the test
    if run_k8s_test "$TEST_TYPE"; then
        print_success "🎉 K8s performance test completed successfully!"
        
        if [[ "$TEST_TYPE" == "staging" ]] && [[ "$ENVIRONMENT" == "staging" ]]; then
            echo ""
            print_k8s "🚀 Staging Validation Complete!"
            print_k8s "Review the generated report for production deployment readiness"
        fi
    else
        print_error "❌ K8s performance test failed"
        exit 1
    fi
    
    echo ""
    print_k8s "📁 Results available in: $RESULTS_DIR"
    print_k8s "📊 Review the K8s report for detailed analysis"
}

# Run main function
main "$@" 