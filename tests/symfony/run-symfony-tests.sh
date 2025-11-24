#!/bin/bash

# Symfony Framework Test Runner using PHP-Go
# This script runs a sample of Symfony framework tests using PHP-Go

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SYMFONY_DIR="$SCRIPT_DIR/symfony-framework"
PHP_GO="$SCRIPT_DIR/../../php-go"
REPORT_DIR="$SCRIPT_DIR/reports"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Symfony Framework Test Suite Runner ==="
echo "Date: $(date)"
echo ""

# Build php-go if needed
if [ ! -f "$PHP_GO" ]; then
    echo "Building php-go..."
    cd "$SCRIPT_DIR/../.."
    go build -o php-go ./cmd/php-go
    cd "$SCRIPT_DIR"
fi

# Create reports directory
mkdir -p "$REPORT_DIR"

# Select sample tests from different components
echo "=== Selecting Sample Tests ==="
echo "Testing a representative sample from major components..."
echo ""

COMPONENTS=(
    "Console"
    "DependencyInjection"
    "HttpFoundation"
    "HttpKernel"
    "Cache"
    "EventDispatcher"
    "Routing"
    "Config"
)

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="$REPORT_DIR/test_results_$TIMESTAMP.txt"
SUMMARY_FILE="$REPORT_DIR/test_summary_$TIMESTAMP.md"

total_tests=0
tests_run=0
tests_passed=0
tests_failed=0
tests_error=0

echo "# Symfony Framework Test Results" > "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "**Date**: $(date)" >> "$SUMMARY_FILE"
echo "**PHP-Go Version**: Custom Build" >> "$SUMMARY_FILE"
echo "**Symfony Version**: 7.3" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

for component in "${COMPONENTS[@]}"; do
    comp_dir="$SYMFONY_DIR/src/Symfony/Component/$component/Tests"
    
    if [ ! -d "$comp_dir" ]; then
        echo "Component $component not found, skipping..."
        continue
    fi
    
    echo -e "${YELLOW}Testing: $component${NC}"
    echo "" >> "$SUMMARY_FILE"
    echo "## $component Component" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    
    # Find first 5 test files for this component
    test_files=($(find "$comp_dir" -name "*Test.php" | head -5))
    
    comp_total=0
    comp_passed=0
    comp_failed=0
    comp_error=0
    
    for test_file in "${test_files[@]}"; do
        ((comp_total++))
        ((total_tests++))
        
        test_name=$(basename "$test_file")
        rel_path="${test_file#$SYMFONY_DIR/}"
        
        echo -n "  Testing: $test_name ... "
        
        # Try to parse the test file
        if "$PHP_GO" parse "$test_file" > /dev/null 2>&1; then
            echo -e "${GREEN}PASS (parse)${NC}"
            echo "- ✅ $test_name (parse success)" >> "$SUMMARY_FILE"
            ((comp_passed++))
            ((tests_passed++))
        else
            error_msg=$("$PHP_GO" parse "$test_file" 2>&1 | head -5)
            echo -e "${RED}FAIL (parse)${NC}"
            echo "- ❌ $test_name (parse failed)" >> "$SUMMARY_FILE"
            echo "  \`\`\`" >> "$SUMMARY_FILE"
            echo "  $error_msg" | head -2 >> "$SUMMARY_FILE"
            echo "  \`\`\`" >> "$SUMMARY_FILE"
            ((comp_failed++))
            ((tests_failed++))
            
            # Log detailed error
            echo "FAILED: $rel_path" >> "$RESULTS_FILE"
            echo "$error_msg" >> "$RESULTS_FILE"
            echo "---" >> "$RESULTS_FILE"
        fi
    done
    
    echo ""
    echo "**Results**: $comp_passed/$comp_total passed" >> "$SUMMARY_FILE"
    echo "Component: $component - $comp_passed/$comp_total passed"
    echo ""
done

# Generate summary
echo "" >> "$SUMMARY_FILE"
echo "## Overall Summary" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "| Metric | Value |" >> "$SUMMARY_FILE"
echo "|--------|-------|" >> "$SUMMARY_FILE"
echo "| Total Tests Sampled | $total_tests |" >> "$SUMMARY_FILE"
echo "| Passed (Parse) | $tests_passed |" >> "$SUMMARY_FILE"
echo "| Failed (Parse) | $tests_failed |" >> "$SUMMARY_FILE"
echo "| Success Rate | $(awk "BEGIN {printf \"%.2f\", ($tests_passed/$total_tests)*100}")% |" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

echo "=== Test Results Summary ==="
echo "Total Tests Sampled: $total_tests"
echo -e "${GREEN}Passed (Parse): $tests_passed${NC}"
echo -e "${RED}Failed (Parse): $tests_failed${NC}"
echo "Success Rate: $(awk "BEGIN {printf \"%.2f\", ($tests_passed/$total_tests)*100}")%"
echo ""
echo "Detailed results saved to:"
echo "  - $SUMMARY_FILE"
echo "  - $RESULTS_FILE"

