#!/bin/bash

# Coverage Analysis Script for URL-DB
# This script provides comprehensive test coverage analysis

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_FILE="coverage.out"
HTML_FILE="coverage.html"
MIN_COVERAGE=20
TARGET_COVERAGE=50

echo -e "${BLUE}🧪 URL-DB Test Coverage Analysis${NC}"
echo "=================================="

# Clean up previous coverage files
rm -f ${COVERAGE_FILE} ${HTML_FILE}

# Step 1: Run tests with coverage
echo -e "${BLUE}📊 Running tests with coverage...${NC}"
go test -coverprofile=${COVERAGE_FILE} ./... 2>/dev/null || {
    echo -e "${RED}❌ Tests failed. Please fix failing tests first.${NC}"
    exit 1
}

# Step 2: Generate overall coverage stats
echo -e "${BLUE}📈 Overall Coverage Statistics${NC}"
echo "------------------------------"
TOTAL_COVERAGE=$(go tool cover -func=${COVERAGE_FILE} | tail -1 | awk '{print $3}' | sed 's/%//')
echo -e "Total Coverage: ${GREEN}${TOTAL_COVERAGE}%${NC}"

if (( $(echo "$TOTAL_COVERAGE >= $TARGET_COVERAGE" | bc -l) )); then
    echo -e "Status: ${GREEN}✅ Above target (${TARGET_COVERAGE}%)${NC}"
elif (( $(echo "$TOTAL_COVERAGE >= $MIN_COVERAGE" | bc -l) )); then
    echo -e "Status: ${YELLOW}⚠️  Above minimum (${MIN_COVERAGE}%) but below target${NC}"
else
    echo -e "Status: ${RED}❌ Below minimum threshold${NC}"
fi

# Step 3: Package-level coverage analysis
echo -e "\n${BLUE}📦 Package-level Coverage${NC}"
echo "-------------------------"

echo "🏆 HIGH COVERAGE (60%+):"
go test -cover ./... 2>/dev/null | grep -E "coverage: [6-9][0-9]\.[0-9]%|coverage: 100\.0%" | sort -k4 -nr | while read line; do
    package=$(echo "$line" | awk '{print $2}')
    coverage=$(echo "$line" | awk '{print $4}')
    echo -e "  ${GREEN}${package}${NC}: ${coverage}"
done

echo -e "\n📊 MEDIUM COVERAGE (30-59%):"
go test -cover ./... 2>/dev/null | grep -E "coverage: [3-5][0-9]\.[0-9]%" | sort -k4 -nr | while read line; do
    package=$(echo "$line" | awk '{print $2}')
    coverage=$(echo "$line" | awk '{print $4}')
    echo -e "  ${YELLOW}${package}${NC}: ${coverage}"
done

echo -e "\n🔴 LOW COVERAGE (0-29%):"
go test -cover ./... 2>/dev/null | grep -E "coverage: [0-2][0-9]\.[0-9]%|coverage: 0\.0%" | sort -k4 -nr | while read line; do
    package=$(echo "$line" | awk '{print $2}')
    coverage=$(echo "$line" | awk '{print $4}')
    echo -e "  ${RED}${package}${NC}: ${coverage}"
done

echo -e "\n⚪ NO STATEMENTS:"
go test -cover ./... 2>/dev/null | grep "coverage: \[no statements\]" | while read line; do
    package=$(echo "$line" | awk '{print $2}')
    echo -e "  ${package}"
done

# Step 4: Function-level analysis - Most Critical (0% coverage)
echo -e "\n${BLUE}🔍 Critical Functions (0% Coverage)${NC}"
echo "-----------------------------------"
go tool cover -func=${COVERAGE_FILE} | grep "0.0%" | head -20 | while read line; do
    func_info=$(echo "$line" | awk '{print $1":"$2}')
    echo -e "  ${RED}${func_info}${NC}"
done

# Step 5: High potential improvement functions (75-95%)
echo -e "\n${BLUE}⚡ High Potential Functions (75-95% Coverage)${NC}"
echo "---------------------------------------------"
go tool cover -func=${COVERAGE_FILE} | grep -E "[7-9][0-9]\.[0-9]%" | grep -v "100.0%" | head -10 | while read line; do
    func_info=$(echo "$line" | awk '{print $1":"$2}')
    coverage=$(echo "$line" | awk '{print $3}')
    echo -e "  ${YELLOW}${func_info}${NC}: ${coverage}"
done

# Step 6: Identify untested files/packages
echo -e "\n${BLUE}📁 Packages Without Tests${NC}"
echo "-------------------------"
find ./internal -name "*.go" -not -name "*_test.go" | while read file; do
    dir=$(dirname "$file")
    if [ ! -f "${dir}"/*_test.go ] 2>/dev/null; then
        if [ "$dir" != "$last_dir" ]; then
            echo -e "  ${RED}${dir}${NC}"
            last_dir="$dir"
        fi
    fi
done | sort | uniq

# Step 7: Test file statistics  
echo -e "\n${BLUE}📊 Test File Statistics${NC}"
echo "-----------------------"
TOTAL_GO_FILES=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
TOTAL_TEST_FILES=$(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)
TEST_RATIO=$(echo "scale=2; $TOTAL_TEST_FILES * 100 / $TOTAL_GO_FILES" | bc)

echo "Total Go files: $TOTAL_GO_FILES"
echo "Total test files: $TOTAL_TEST_FILES"
echo -e "Test file ratio: ${YELLOW}${TEST_RATIO}%${NC}"

# Step 8: Generate HTML report
echo -e "\n${BLUE}📄 Generating HTML Coverage Report${NC}"
echo "-----------------------------------"
go tool cover -html=${COVERAGE_FILE} -o ${HTML_FILE}
echo -e "HTML report generated: ${GREEN}${HTML_FILE}${NC}"

# Step 9: Quick improvement suggestions
echo -e "\n${BLUE}💡 Quick Improvement Suggestions${NC}"
echo "--------------------------------"

# Count 0% coverage functions by package
echo "Top packages for immediate improvement:"
go tool cover -func=${COVERAGE_FILE} | grep "0.0%" | cut -d'/' -f1-3 | sort | uniq -c | sort -nr | head -5 | while read count package; do
    echo -e "  ${RED}${package}${NC}: ${count} untested functions"
done

echo -e "\n${GREEN}🎯 Recommended Next Steps:${NC}"
echo "1. Focus on packages with 0% coverage first"
echo "2. Improve functions with 75-95% coverage to 100%"
echo "3. Add integration tests for main.go and server initialization"
echo "4. Create missing test files for untested packages"

echo -e "\n${BLUE}📊 Coverage Analysis Complete${NC}"
echo "To view detailed HTML report: open ${HTML_FILE}"