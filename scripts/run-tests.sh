#!/bin/bash

# URL Database Test Runner
# Comprehensive testing script for the URL Database system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_TIMEOUT=300 # 5 minutes
COVERAGE_THRESHOLD=80
OUTPUT_DIR="test-output"

# Helper functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Check dependencies
check_dependencies() {
    print_header "Checking Dependencies"
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    print_success "Go is installed: $(go version)"
    
    if ! command -v golangci-lint &> /dev/null; then
        print_warning "golangci-lint not found, installing..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    print_success "golangci-lint is available"
    
    # Check if we're in the correct directory
    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi
    print_success "Project root directory confirmed"
}

# Download dependencies
download_dependencies() {
    print_header "Downloading Dependencies"
    
    print_info "Running go mod download..."
    go mod download
    
    print_info "Running go mod tidy..."
    go mod tidy
    
    print_success "Dependencies updated"
}

# Format code
format_code() {
    print_header "Formatting Code"
    
    print_info "Running go fmt..."
    go fmt ./...
    
    print_success "Code formatted"
}

# Lint code
lint_code() {
    print_header "Linting Code"
    
    print_info "Running golangci-lint..."
    if golangci-lint run --timeout=5m --out-format=colored-line-number 2>&1 | tee "$OUTPUT_DIR/lint-report.txt"; then
        print_success "Linting passed"
        return 0
    else
        print_error "Linting failed"
        return 1
    fi
}

# Build project
build_project() {
    print_header "Building Project"
    
    print_info "Building main application..."
    if go build -o "$OUTPUT_DIR/url-db" ./cmd/server; then
        print_success "Build successful"
        return 0
    else
        print_error "Build failed"
        return 1
    fi
}

# Generate documentation
generate_docs() {
    print_header "Generating Documentation"
    
    if command -v swag &> /dev/null; then
        print_info "Generating Swagger documentation..."
        swag init -g cmd/server/main.go -o docs
        print_success "Swagger documentation generated"
    else
        print_warning "swag not found, skipping Swagger documentation generation"
    fi
}

# Run unit tests
run_unit_tests() {
    print_header "Running Unit Tests"
    
    print_info "Running all unit tests with coverage..."
    
    local test_args=(
        "./..."
        "-v"
        "-race"
        "-timeout=${TEST_TIMEOUT}s"
        "-coverprofile=$OUTPUT_DIR/coverage.out"
        "-covermode=atomic"
        "-count=1"
    )
    
    if go test "${test_args[@]}" 2>&1 | tee "$OUTPUT_DIR/test-results.txt"; then
        print_success "Unit tests passed"
        return 0
    else
        print_error "Unit tests failed"
        return 1
    fi
}

# Run specific package tests
run_package_tests() {
    local package="$1"
    print_header "Running Tests for Package: $package"
    
    print_info "Testing package $package..."
    
    if go test "./$package/..." -v -race -timeout="${TEST_TIMEOUT}s" 2>&1 | tee "$OUTPUT_DIR/test-${package//\//-}.txt"; then
        print_success "Package $package tests passed"
        return 0
    else
        print_error "Package $package tests failed"
        return 1
    fi
}

# Generate coverage report
generate_coverage_report() {
    print_header "Generating Coverage Report"
    
    if [[ ! -f "$OUTPUT_DIR/coverage.out" ]]; then
        print_warning "Coverage file not found, skipping coverage report"
        return 1
    fi
    
    print_info "Generating HTML coverage report..."
    go tool cover -html="$OUTPUT_DIR/coverage.out" -o "$OUTPUT_DIR/coverage.html"
    print_success "Coverage report generated: $OUTPUT_DIR/coverage.html"
    
    print_info "Calculating coverage percentage..."
    local coverage=$(go tool cover -func="$OUTPUT_DIR/coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
    
    echo "Coverage: $coverage%" | tee "$OUTPUT_DIR/coverage-summary.txt"
    
    if (( $(echo "$coverage >= $COVERAGE_THRESHOLD" | bc -l) )); then
        print_success "Coverage threshold met: $coverage% >= $COVERAGE_THRESHOLD%"
        return 0
    else
        print_error "Coverage threshold not met: $coverage% < $COVERAGE_THRESHOLD%"
        return 1
    fi
}

# Run benchmarks
run_benchmarks() {
    print_header "Running Benchmarks"
    
    print_info "Running benchmark tests..."
    
    if go test ./... -bench=. -benchmem -run=^$ > "$OUTPUT_DIR/benchmark-results.txt" 2>&1; then
        print_success "Benchmarks completed"
        cat "$OUTPUT_DIR/benchmark-results.txt"
        return 0
    else
        print_error "Benchmarks failed"
        return 1
    fi
}

# Check for race conditions
check_race_conditions() {
    print_header "Checking for Race Conditions"
    
    print_info "Running tests with race detector..."
    
    if go test ./... -race -short > "$OUTPUT_DIR/race-detection.txt" 2>&1; then
        print_success "No race conditions detected"
        return 0
    else
        print_error "Race conditions detected"
        cat "$OUTPUT_DIR/race-detection.txt"
        return 1
    fi
}

# Validate test coverage for critical packages
validate_critical_coverage() {
    print_header "Validating Critical Package Coverage"
    
    local critical_packages=(
        "internal/domains"
        "internal/nodes"
        "internal/attributes"
        "internal/mcp"
        "internal/nodeattributes"
    )
    
    local failed_packages=()
    
    for package in "${critical_packages[@]}"; do
        print_info "Checking coverage for $package..."
        
        # Run tests for specific package with coverage
        if go test "./$package/..." -coverprofile="$OUTPUT_DIR/coverage-${package//\//-}.out" > /dev/null 2>&1; then
            local coverage=$(go tool cover -func="$OUTPUT_DIR/coverage-${package//\//-}.out" | grep total | awk '{print $3}' | sed 's/%//')
            
            if [[ -n "$coverage" ]]; then
                if (( $(echo "$coverage >= 85" | bc -l) )); then
                    print_success "$package: $coverage%"
                else
                    print_error "$package: $coverage% (below 85% threshold)"
                    failed_packages+=("$package")
                fi
            else
                print_warning "$package: No coverage data"
            fi
        else
            print_error "$package: Tests failed"
            failed_packages+=("$package")
        fi
    done
    
    if [[ ${#failed_packages[@]} -eq 0 ]]; then
        print_success "All critical packages meet coverage requirements"
        return 0
    else
        print_error "The following packages failed coverage requirements: ${failed_packages[*]}"
        return 1
    fi
}

# Test MCP integration
test_mcp_integration() {
    print_header "Testing MCP Integration"
    
    print_info "Testing MCP stdio mode..."
    
    # Build the application
    if ! go build -o "$OUTPUT_DIR/url-db-test" ./cmd/server; then
        print_error "Failed to build application for MCP testing"
        return 1
    fi
    
    # Test help command
    if "$OUTPUT_DIR/url-db-test" -help > "$OUTPUT_DIR/mcp-help.txt" 2>&1; then
        print_success "MCP help command works"
    else
        print_error "MCP help command failed"
        return 1
    fi
    
    # Test version command
    if "$OUTPUT_DIR/url-db-test" -version > "$OUTPUT_DIR/mcp-version.txt" 2>&1; then
        print_success "MCP version command works"
    else
        print_error "MCP version command failed"
        return 1
    fi
    
    print_success "MCP integration tests completed"
    return 0
}

# Generate test summary
generate_test_summary() {
    print_header "Test Summary"
    
    local summary_file="$OUTPUT_DIR/test-summary.txt"
    
    {
        echo "URL Database Test Summary"
        echo "========================"
        echo "Timestamp: $(date)"
        echo "Go Version: $(go version)"
        echo ""
        
        if [[ -f "$OUTPUT_DIR/coverage-summary.txt" ]]; then
            echo "Coverage Summary:"
            cat "$OUTPUT_DIR/coverage-summary.txt"
            echo ""
        fi
        
        echo "Test Files Generated:"
        ls -la "$OUTPUT_DIR"
        echo ""
        
        echo "Test Commands Used:"
        echo "- Unit Tests: go test ./... -v -race -timeout=300s -coverprofile=coverage.out"
        echo "- Linting: golangci-lint run --timeout=5m"
        echo "- Benchmarks: go test ./... -bench=. -benchmem"
        echo "- Race Detection: go test ./... -race -short"
        
    } > "$summary_file"
    
    cat "$summary_file"
    print_success "Test summary generated: $summary_file"
}

# Cleanup function
cleanup() {
    print_info "Cleaning up temporary files..."
    # Remove test binaries
    rm -f "$OUTPUT_DIR/url-db-test"
    # Keep important reports
}

# Main execution
main() {
    local exit_code=0
    local start_time=$(date +%s)
    
    print_header "URL Database Test Suite"
    print_info "Starting comprehensive test execution..."
    
    # Parse command line arguments
    local run_all=true
    local run_lint=false
    local run_build=false
    local run_tests=false
    local run_coverage=false
    local run_benchmarks=false
    local run_mcp=false
    local package_filter=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --lint-only)
                run_all=false
                run_lint=true
                shift
                ;;
            --build-only)
                run_all=false
                run_build=true
                shift
                ;;
            --tests-only)
                run_all=false
                run_tests=true
                shift
                ;;
            --coverage-only)
                run_all=false
                run_coverage=true
                shift
                ;;
            --benchmarks-only)
                run_all=false
                run_benchmarks=true
                shift
                ;;
            --mcp-only)
                run_all=false
                run_mcp=true
                shift
                ;;
            --package)
                package_filter="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --lint-only      Run only linting"
                echo "  --build-only     Run only build"
                echo "  --tests-only     Run only unit tests"
                echo "  --coverage-only  Run only coverage analysis"
                echo "  --benchmarks-only Run only benchmarks"
                echo "  --mcp-only       Run only MCP integration tests"
                echo "  --package DIR    Run tests for specific package"
                echo "  --help           Show this help"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Execute tests based on arguments
    if [[ "$run_all" == true ]]; then
        check_dependencies || exit_code=1
        download_dependencies || exit_code=1
        format_code || exit_code=1
        generate_docs || exit_code=1
        lint_code || exit_code=1
        build_project || exit_code=1
        run_unit_tests || exit_code=1
        generate_coverage_report || exit_code=1
        validate_critical_coverage || exit_code=1
        check_race_conditions || exit_code=1
        run_benchmarks || exit_code=1
        test_mcp_integration || exit_code=1
    else
        [[ "$run_lint" == true ]] && (lint_code || exit_code=1)
        [[ "$run_build" == true ]] && (build_project || exit_code=1)
        [[ "$run_tests" == true ]] && (run_unit_tests || exit_code=1)
        [[ "$run_coverage" == true ]] && (generate_coverage_report || exit_code=1)
        [[ "$run_benchmarks" == true ]] && (run_benchmarks || exit_code=1)
        [[ "$run_mcp" == true ]] && (test_mcp_integration || exit_code=1)
    fi
    
    # Run package-specific tests if specified
    if [[ -n "$package_filter" ]]; then
        run_package_tests "$package_filter" || exit_code=1
    fi
    
    generate_test_summary
    cleanup
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "All tests completed successfully in ${duration}s"
    else
        print_error "Some tests failed. Check $OUTPUT_DIR for details. Total time: ${duration}s"
    fi
    
    exit $exit_code
}

# Set trap for cleanup
trap cleanup EXIT

# Run main function with all arguments
main "$@"