#!/bin/bash

# URL-DB Test Runner Script
# Comprehensive test execution with different modes and options

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
MODE="all"
VERBOSE=false
COVERAGE=false
PACKAGE=""
TIMEOUT="10m"
CLEAN=false

# Function to show usage
show_usage() {
    cat << EOF
${BLUE}URL-DB Test Runner${NC}

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -m, --mode MODE        Test mode: all, unit, integration, mcp, coverage (default: all)
    -p, --package PKG      Run tests for specific package (e.g., internal/mcp)
    -v, --verbose          Enable verbose output
    -c, --coverage         Run with coverage analysis
    -t, --timeout TIME     Set test timeout (default: 10m)
    --clean               Clean coverage files before running
    -h, --help            Show this help message

MODES:
    all           Run all tests
    unit          Run unit tests only
    integration   Run integration tests only  
    mcp           Run MCP-specific tests only
    coverage      Run tests with detailed coverage analysis

EXAMPLES:
    $0                                    # Run all tests
    $0 -m unit -v                        # Run unit tests with verbose output
    $0 -m coverage                       # Run with coverage analysis
    $0 -p internal/mcp -v                # Test specific package
    $0 -m integration --timeout 15m      # Integration tests with custom timeout

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--mode)
            MODE="$2"
            shift 2
            ;;
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_usage
            exit 1
            ;;
    esac
done

# Clean coverage files if requested
if [ "$CLEAN" = true ]; then
    echo -e "${YELLOW}🧹 Cleaning coverage files...${NC}"
    rm -f coverage.out coverage.html
fi

# Build test command based on options
build_test_cmd() {
    local cmd="go test"
    
    if [ "$VERBOSE" = true ]; then
        cmd="$cmd -v"
    fi
    
    if [ "$COVERAGE" = true ]; then
        cmd="$cmd -coverprofile=coverage.out"
    fi
    
    cmd="$cmd -timeout $TIMEOUT"
    
    echo "$cmd"
}

# Run specific package tests
run_package_tests() {
    local pkg="$1"
    local cmd=$(build_test_cmd)
    
    echo -e "${BLUE}🧪 Running tests for package: ${pkg}${NC}"
    echo "Command: $cmd ./$pkg/..."
    echo "----------------------------------------"
    
    $cmd ./$pkg/...
}

# Run all tests
run_all_tests() {
    local cmd=$(build_test_cmd)
    
    echo -e "${BLUE}🧪 Running all tests${NC}"
    echo "Command: $cmd ./..."
    echo "----------------------------------------"
    
    $cmd ./...
}

# Run unit tests only (exclude integration tags)
run_unit_tests() {
    local cmd=$(build_test_cmd)
    cmd="$cmd -short"
    
    echo -e "${BLUE}🧪 Running unit tests only${NC}"
    echo "Command: $cmd ./..."
    echo "----------------------------------------"
    
    $cmd ./...
}

# Run integration tests only
run_integration_tests() {
    local cmd=$(build_test_cmd)
    cmd="$cmd -tags=integration"
    
    echo -e "${BLUE}🧪 Running integration tests${NC}"
    echo "Command: $cmd ./..."
    echo "----------------------------------------"
    
    $cmd ./...
}

# Run MCP-specific tests
run_mcp_tests() {
    local cmd=$(build_test_cmd)
    
    echo -e "${BLUE}🧪 Running MCP-specific tests${NC}"
    echo "Command: $cmd ./internal/mcp/... ./internal/interfaces/mcp/..."
    echo "----------------------------------------"
    
    $cmd ./internal/mcp/... ./internal/interfaces/mcp/...
}

# Run coverage analysis
run_coverage_analysis() {
    echo -e "${BLUE}📊 Running comprehensive coverage analysis${NC}"
    echo "============================================"
    
    # Run tests with coverage
    local cmd="go test -coverprofile=coverage.out -timeout $TIMEOUT"
    if [ "$VERBOSE" = true ]; then
        cmd="$cmd -v"
    fi
    
    echo "Running: $cmd ./..."
    $cmd ./...
    
    # Run coverage analysis script
    if [ -f "./scripts/coverage_analysis.sh" ]; then
        echo -e "\n${BLUE}📊 Running detailed coverage analysis...${NC}"
        ./scripts/coverage_analysis.sh
    else
        echo -e "${YELLOW}⚠️  Coverage analysis script not found${NC}"
        echo "Basic coverage summary:"
        go tool cover -func=coverage.out | tail -5
    fi
}

# Main execution
echo -e "${GREEN}🚀 URL-DB Test Runner Starting${NC}"
echo "================================"
echo "Mode: $MODE"
echo "Package: ${PACKAGE:-"all"}"
echo "Verbose: $VERBOSE"
echo "Coverage: $COVERAGE"
echo "Timeout: $TIMEOUT"
echo ""

# Execute based on mode
case $MODE in
    "all")
        if [ -n "$PACKAGE" ]; then
            run_package_tests "$PACKAGE"
        else
            run_all_tests
        fi
        ;;
    "unit")
        run_unit_tests
        ;;
    "integration")
        run_integration_tests
        ;;
    "mcp")
        run_mcp_tests
        ;;
    "coverage")
        run_coverage_analysis
        ;;
    *)
        echo -e "${RED}❌ Unknown mode: $MODE${NC}"
        show_usage
        exit 1
        ;;
esac

# Show results
if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}✅ Tests completed successfully${NC}"
    
    # Show coverage summary if coverage was enabled
    if [ "$COVERAGE" = true ] && [ -f "coverage.out" ]; then
        echo -e "\n${BLUE}📊 Coverage Summary${NC}"
        echo "-------------------"
        go tool cover -func=coverage.out | tail -1
        echo -e "📄 HTML report: ${GREEN}coverage.html${NC} (if generated)"
    fi
else
    echo -e "\n${RED}❌ Tests failed${NC}"
    exit 1
fi

# Additional information
echo -e "\n${BLUE}ℹ️  Additional Information${NC}"
echo "-------------------------"
echo "• Use -h or --help for more options"
echo "• Coverage HTML report: open coverage.html"
echo "• For detailed analysis: ./scripts/coverage_analysis.sh"