#!/bin/bash
#
# Symfony Issue Identification Script for PHP-Go
#
# This script systematically identifies issues with Symfony compatibility:
# 1. Parse errors in Symfony files
# 2. Missing PHP functions
# 3. Missing PHP classes/interfaces
# 4. Other compatibility issues
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
SYMFONY_DIR="$SCRIPT_DIR/symfony-app"
REPORT_DIR="$SCRIPT_DIR/reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="$REPORT_DIR/issues_${TIMESTAMP}.md"

# Create reports directory
mkdir -p "$REPORT_DIR"

echo "Symfony Issue Identification for PHP-Go"
echo "========================================"
echo ""
echo "Symfony Dir: $SYMFONY_DIR"
echo "Report File: $REPORT_FILE"
echo ""

# Check if php-go binary exists
if [ ! -f "$PHPGO" ]; then
    echo "Error: php-go binary not found at $PHPGO"
    echo "Please build it first: go build -o php-go ./cmd/php-go"
    exit 1
fi

# Check if Symfony is installed
if [ ! -f "$SYMFONY_DIR/bin/console" ]; then
    echo "Error: Symfony not found at $SYMFONY_DIR"
    exit 1
fi

# Get Symfony version
SYMFONY_VERSION=$($SYMFONY_DIR/bin/console --version 2>/dev/null | head -1 || echo "Unknown")

# Initialize report
cat > "$REPORT_FILE" <<EOF
# Symfony Compatibility Issues Report
Generated: $(date)
Symfony Directory: $SYMFONY_DIR
Symfony Version: $SYMFONY_VERSION
PHP-Go Version: $($PHPGO --version 2>&1 || echo "unknown")

---

## Executive Summary

This report identifies compatibility issues between PHP-Go and Symfony 7.3.

EOF

# Test 1: Parse all Symfony PHP files
echo "Test 1: Parsing all Symfony PHP files..."
echo "-----------------------------------------"

PARSE_LOG="$REPORT_DIR/parse_errors_${TIMESTAMP}.log"
PARSE_SUCCESS=0
PARSE_FAIL=0
PARSE_TOTAL=0

# Create temporary file for failed files
FAILED_FILES="$REPORT_DIR/failed_files_${TIMESTAMP}.txt"
> "$FAILED_FILES"

# Find and parse all PHP files
while IFS= read -r php_file; do
    PARSE_TOTAL=$((PARSE_TOTAL + 1))

    # Show progress every 100 files
    if [ $((PARSE_TOTAL % 100)) -eq 0 ]; then
        echo "  Parsed $PARSE_TOTAL files... ($PARSE_SUCCESS successful, $PARSE_FAIL failed)"
    fi

    # Try to parse the file
    if $PHPGO parse "$php_file" > /dev/null 2>&1; then
        PARSE_SUCCESS=$((PARSE_SUCCESS + 1))
    else
        PARSE_FAIL=$((PARSE_FAIL + 1))
        echo "$php_file" >> "$FAILED_FILES"
        # Capture error details
        echo "=== $php_file ===" >> "$PARSE_LOG"
        $PHPGO parse "$php_file" 2>&1 >> "$PARSE_LOG" || true
        echo "" >> "$PARSE_LOG"
    fi
done < <(find "$SYMFONY_DIR" -name "*.php" -type f)

echo ""
echo "Parse Results:"
echo "  Total files: $PARSE_TOTAL"
echo "  Successful: $PARSE_SUCCESS"
echo "  Failed: $PARSE_FAIL"
echo "  Success rate: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%"
echo ""

# Add parse results to report
cat >> "$REPORT_FILE" <<EOF

## Parse Test Results

Total PHP files tested: $PARSE_TOTAL
Successfully parsed: $PARSE_SUCCESS
Parse failures: $PARSE_FAIL
Success rate: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%

### Parse Error Log
See: \`$(basename "$PARSE_LOG")\` for detailed parse errors.

### Failed Files
See: \`$(basename "$FAILED_FILES")\` for list of files that failed to parse.

EOF

# Test 2: Analyze parse errors to identify common issues
echo "Test 2: Analyzing parse errors..."
echo "-----------------------------------"

if [ -f "$PARSE_LOG" ] && [ -s "$PARSE_LOG" ]; then
    echo "Analyzing error patterns..."

    # Extract unique error messages
    ERROR_SUMMARY="$REPORT_DIR/error_summary_${TIMESTAMP}.txt"

    # Count error types
    grep -h "Parse error:" "$PARSE_LOG" 2>/dev/null | sort | uniq -c | sort -rn > "$ERROR_SUMMARY" || true

    if [ -s "$ERROR_SUMMARY" ]; then
        echo "  Found $(wc -l < "$ERROR_SUMMARY" | tr -d ' ') unique error types"

        cat >> "$REPORT_FILE" <<EOF

## Error Analysis

### Top Parse Errors

\`\`\`
$(head -20 "$ERROR_SUMMARY")
\`\`\`

See: \`$(basename "$ERROR_SUMMARY")\` for full error breakdown.

EOF
    fi
fi

# Finalize report
cat >> "$REPORT_FILE" <<EOF

---

## Next Steps

1. Review parse errors to identify missing language features
2. Prioritize fixes based on frequency of errors
3. Implement missing features in PHP-Go parser
4. Re-run tests to measure improvement

EOF

echo ""
echo "========================================"
echo "Report generated: $REPORT_FILE"
echo "Parse log: $PARSE_LOG"
echo "Failed files: $FAILED_FILES"
if [ -f "$ERROR_SUMMARY" ]; then
    echo "Error summary: $ERROR_SUMMARY"
fi
echo "========================================"
