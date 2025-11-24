#!/bin/bash
#
# Laravel Issue Identification Script for PHP-Go
#
# This script systematically identifies issues with Laravel compatibility:
# 1. Parse errors in Laravel files
# 2. Missing PHP functions
# 3. Missing PHP classes/interfaces
# 4. Other compatibility issues
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
LARAVEL_DIR="$SCRIPT_DIR"
REPORT_DIR="$SCRIPT_DIR/reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="$REPORT_DIR/issues_${TIMESTAMP}.md"

# Create reports directory
mkdir -p "$REPORT_DIR"

echo "Laravel Issue Identification for PHP-Go"
echo "========================================"
echo ""
echo "Laravel Dir: $LARAVEL_DIR"
echo "Report File: $REPORT_FILE"
echo ""

# Check if php-go binary exists
if [ ! -f "$PHPGO" ]; then
    echo "Error: php-go binary not found at $PHPGO"
    echo "Please build it first: go build -o php-go ./cmd/php-go"
    exit 1
fi

# Check if Laravel is installed
if [ ! -f "$LARAVEL_DIR/artisan" ]; then
    echo "Error: Laravel not found at $LARAVEL_DIR"
    exit 1
fi

# Initialize report
cat > "$REPORT_FILE" <<EOF
# Laravel Compatibility Issues Report
Generated: $(date)
Laravel Directory: $LARAVEL_DIR
Laravel Version: v12.10.1 (Framework v12.39.0)
PHP-Go Version: $($PHPGO --version 2>&1 || echo "unknown")

---

## Executive Summary

This report identifies compatibility issues between PHP-Go and Laravel v12.10.1.

EOF

# Test 1: Parse all Laravel PHP files
echo "Test 1: Parsing all Laravel PHP files..."
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

    # Show progress every 100 files (Laravel is much larger than WordPress)
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
done < <(find "$LARAVEL_DIR" -name "*.php" -type f)

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

- **Total PHP files**: $PARSE_TOTAL
- **Successfully parsed**: $PARSE_SUCCESS
- **Parse failures**: $PARSE_FAIL
- **Success rate**: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%

### Parse Error Analysis

EOF

# Test 2: Analyze parse error patterns
echo "Test 2: Analyzing parse error patterns..."
echo "------------------------------------------"

# Extract unique error patterns from parse log
if [ -f "$PARSE_LOG" ]; then
    echo "Extracting error patterns..."

    # Count error types
    ERROR_PATTERNS="$REPORT_DIR/error_patterns_${TIMESTAMP}.txt"

    # Extract error messages and count them
    grep -h "Parse error:" "$PARSE_LOG" 2>/dev/null | \
        sed 's/at .*$//' | \
        sed 's/line [0-9]*/line N/' | \
        sed 's/column [0-9]*/column N/' | \
        sort | uniq -c | sort -rn > "$ERROR_PATTERNS" || true

    # Display top 20 error patterns
    echo ""
    echo "Top 20 Parse Error Patterns:"
    head -20 "$ERROR_PATTERNS" 2>/dev/null || echo "No error patterns found"
    echo ""

    # Add to report
    echo "#### Top Parse Error Patterns" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "\`\`\`" >> "$REPORT_FILE"
    head -30 "$ERROR_PATTERNS" 2>/dev/null >> "$REPORT_FILE" || echo "No error patterns found" >> "$REPORT_FILE"
    echo "\`\`\`" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
fi

# Test 3: Categorize errors by directory
echo "Test 3: Categorizing errors by directory..."
echo "--------------------------------------------"

cat >> "$REPORT_FILE" <<EOF

### Parse Errors by Directory

EOF

for dir in app bootstrap config database public resources routes tests vendor; do
    if [ -d "$LARAVEL_DIR/$dir" ]; then
        total=$(find "$LARAVEL_DIR/$dir" -name "*.php" -type f 2>/dev/null | wc -l | tr -d ' ')
        failed=$(grep -c "^$LARAVEL_DIR/$dir/" "$FAILED_FILES" 2>/dev/null || echo "0")
        success=$((total - failed))
        if [ $total -gt 0 ]; then
            success_rate=$(awk "BEGIN {printf \"%.2f\", ($success/$total)*100}")
            echo "  $dir: $success/$total ($success_rate%)"
            echo "- **$dir**: $success/$total files parsed successfully ($success_rate%)" >> "$REPORT_FILE"
        fi
    fi
done
echo ""
echo "" >> "$REPORT_FILE"

# Test 4: Check specific critical files
echo "Test 4: Checking critical Laravel files..."
echo "-------------------------------------------"

cat >> "$REPORT_FILE" <<EOF

### Critical File Analysis

EOF

CRITICAL_FILES=(
    "artisan"
    "public/index.php"
    "bootstrap/app.php"
    "bootstrap/providers.php"
    "app/Console/Kernel.php"
    "app/Exceptions/Handler.php"
    "app/Models/User.php"
    "vendor/laravel/framework/src/Illuminate/Foundation/Application.php"
    "vendor/laravel/framework/src/Illuminate/Routing/Router.php"
    "vendor/laravel/framework/src/Illuminate/Support/ServiceProvider.php"
)

for critical_file in "${CRITICAL_FILES[@]}"; do
    full_path="$LARAVEL_DIR/$critical_file"
    if [ -f "$full_path" ]; then
        if grep -q "^$full_path$" "$FAILED_FILES" 2>/dev/null; then
            echo "  ❌ $critical_file - FAILED"
            echo "- ❌ **$critical_file** - Parse failed" >> "$REPORT_FILE"
        else
            echo "  ✅ $critical_file - PASSED"
            echo "- ✅ **$critical_file** - Parse successful" >> "$REPORT_FILE"
        fi
    else
        echo "  ⚠️  $critical_file - NOT FOUND"
        echo "- ⚠️  **$critical_file** - File not found" >> "$REPORT_FILE"
    fi
done
echo ""
echo "" >> "$REPORT_FILE"

# Test 5: Analyze specific syntax features
echo "Test 5: Analyzing PHP 8+ syntax features..."
echo "--------------------------------------------"

cat >> "$REPORT_FILE" <<EOF

### PHP 8+ Feature Analysis

Analyzing usage of PHP 8+ features in Laravel codebase:

EOF

# Search for PHP 8+ features
echo "Searching for PHP 8+ features in failed files..."

# Attributes
ATTR_COUNT=$(cat "$FAILED_FILES" 2>/dev/null | xargs grep -l "^#\[" 2>/dev/null | wc -l | tr -d ' ')
echo "  Attributes (#[...]): $ATTR_COUNT files"
echo "- **Attributes** (\`#[...]\`): Found in $ATTR_COUNT failed files" >> "$REPORT_FILE"

# Named arguments
NAMED_COUNT=$(cat "$FAILED_FILES" 2>/dev/null | xargs grep -l ": " 2>/dev/null | grep -v "case " | wc -l | tr -d ' ')
echo "  Named arguments (potential): $NAMED_COUNT files"
echo "- **Named arguments**: Potentially used in $NAMED_COUNT failed files" >> "$REPORT_FILE"

# Match expressions
MATCH_COUNT=$(cat "$FAILED_FILES" 2>/dev/null | xargs grep -l "match\s*(" 2>/dev/null | wc -l | tr -d ' ')
echo "  Match expressions: $MATCH_COUNT files"
echo "- **Match expressions**: Found in $MATCH_COUNT failed files" >> "$REPORT_FILE"

# Constructor property promotion
PROMOTED_COUNT=$(cat "$FAILED_FILES" 2>/dev/null | xargs grep -l "public\s*function\s*__construct.*public\|private\|protected" 2>/dev/null | wc -l | tr -d ' ')
echo "  Constructor property promotion: $PROMOTED_COUNT files"
echo "- **Constructor property promotion**: Found in $PROMOTED_COUNT failed files" >> "$REPORT_FILE"

echo ""
echo "" >> "$REPORT_FILE"

# Finalize report
cat >> "$REPORT_FILE" <<EOF

## Recommendations

Based on this analysis, the following areas need attention:

1. **Parse Error Fixes**: Address the top parse error patterns identified above
2. **Critical File Support**: Ensure all critical Laravel files parse successfully
3. **PHP 8+ Features**: Implement missing PHP 8+ syntax features
4. **Testing**: Run comprehensive tests on fixed features

## Files Generated

- Parse error log: \`$(basename "$PARSE_LOG")\`
- Failed files list: \`$(basename "$FAILED_FILES")\`
- Error patterns: \`$(basename "$ERROR_PATTERNS")\`
- This report: \`$(basename "$REPORT_FILE")\`

---

*Report generated by identify-issues.sh on $(date)*
EOF

echo "========================================"
echo "Issue identification complete!"
echo ""
echo "Summary:"
echo "  Total files: $PARSE_TOTAL"
echo "  Parse success: $PARSE_SUCCESS ($((PARSE_SUCCESS * 100 / PARSE_TOTAL))%)"
echo "  Parse failed: $PARSE_FAIL ($((PARSE_FAIL * 100 / PARSE_TOTAL))%)"
echo ""
echo "Report saved to: $REPORT_FILE"
echo "Parse errors logged to: $PARSE_LOG"
echo "Failed files listed in: $FAILED_FILES"
echo ""
echo "Next steps:"
echo "1. Review the error patterns in $REPORT_FILE"
echo "2. Prioritize fixing the most common parse errors"
echo "3. Re-run this script after fixes to measure progress"
echo "========================================"
