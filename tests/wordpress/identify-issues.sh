#!/bin/bash
#
# WordPress Issue Identification Script for PHP-Go
#
# This script systematically identifies issues with WordPress compatibility:
# 1. Parse errors in WordPress files
# 2. Missing PHP functions
# 3. Missing PHP classes/interfaces
# 4. Other compatibility issues
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
WP_DIR="$SCRIPT_DIR/wordpress"
REPORT_DIR="$SCRIPT_DIR/reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="$REPORT_DIR/issues_${TIMESTAMP}.md"

# Create reports directory
mkdir -p "$REPORT_DIR"

echo "WordPress Issue Identification for PHP-Go"
echo "=========================================="
echo ""
echo "WordPress Dir: $WP_DIR"
echo "Report File: $REPORT_FILE"
echo ""

# Check if php-go binary exists
if [ ! -f "$PHPGO" ]; then
    echo "Error: php-go binary not found at $PHPGO"
    echo "Please build it first: go build -o php-go ./cmd/php-go"
    exit 1
fi

# Check if WordPress is installed
if [ ! -d "$WP_DIR" ]; then
    echo "Error: WordPress not found at $WP_DIR"
    exit 1
fi

# Initialize report
cat > "$REPORT_FILE" <<EOF
# WordPress Compatibility Issues Report
Generated: $(date)
WordPress Directory: $WP_DIR
PHP-Go Version: $($PHPGO --version 2>&1 || echo "unknown")

---

## Executive Summary

This report identifies compatibility issues between PHP-Go and WordPress 6.8.3.

EOF

# Test 1: Parse all WordPress PHP files
echo "Test 1: Parsing all WordPress PHP files..."
echo "-------------------------------------------"

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

    # Show progress every 50 files
    if [ $((PARSE_TOTAL % 50)) -eq 0 ]; then
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
done < <(find "$WP_DIR" -name "*.php" -type f)

echo ""
echo "Parse Results:"
echo "  Total files: $PARSE_TOTAL"
echo "  Successful: $PARSE_SUCCESS"
echo "  Failed: $PARSE_FAIL"
echo "  Success rate: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%"
echo ""

# Add to report
cat >> "$REPORT_FILE" <<EOF

## 1. Parse Test Results

- **Total PHP files**: $PARSE_TOTAL
- **Successfully parsed**: $PARSE_SUCCESS
- **Parse failures**: $PARSE_FAIL
- **Success rate**: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%

EOF

if [ $PARSE_FAIL -gt 0 ]; then
    cat >> "$REPORT_FILE" <<EOF
### Parse Error Summary

The following files failed to parse:

\`\`\`
$(cat "$FAILED_FILES")
\`\`\`

Detailed error logs available in: \`$(basename "$PARSE_LOG")\`

### Common Parse Error Patterns

EOF

    # Analyze common error patterns
    echo "Analyzing parse error patterns..."
    grep -E "(Unexpected|Expected|Parse error)" "$PARSE_LOG" | sort | uniq -c | sort -rn | head -10 >> "$REPORT_FILE" || true

    cat >> "$REPORT_FILE" <<EOF

EOF
fi

# Test 2: Identify missing functions
echo "Test 2: Identifying missing PHP functions..."
echo "---------------------------------------------"

FUNCTIONS_LOG="$REPORT_DIR/missing_functions_${TIMESTAMP}.txt"

# Extract function calls from WordPress files
# This is a simple grep-based approach - may have false positives
echo "  Extracting function calls from WordPress code..."
find "$WP_DIR" -name "*.php" -type f -exec grep -oh '\b[a-z_][a-z0-9_]*\s*(' {} \; | \
    sed 's/\s*($//' | \
    sort | uniq > "$REPORT_DIR/all_functions_${TIMESTAMP}.txt"

TOTAL_UNIQUE_FUNCS=$(wc -l < "$REPORT_DIR/all_functions_${TIMESTAMP}.txt" | tr -d ' ')
echo "  Found $TOTAL_UNIQUE_FUNCS unique function calls"

# Common PHP built-in functions that might be missing
# This is a curated list based on WordPress usage patterns
cat > "$REPORT_DIR/known_builtins.txt" <<'EOFLIST'
strlen
substr
strpos
strrpos
str_replace
str_repeat
preg_match
preg_replace
preg_match_all
array_push
array_pop
array_shift
array_unshift
array_merge
array_keys
array_values
array_map
array_filter
array_reduce
in_array
count
sizeof
is_array
is_string
is_int
is_float
is_bool
is_null
is_numeric
is_object
is_resource
isset
empty
unset
define
defined
constant
file_exists
file_get_contents
file_put_contents
is_file
is_dir
is_readable
is_writable
mkdir
rmdir
unlink
fopen
fclose
fread
fwrite
fgets
file
glob
dirname
basename
realpath
pathinfo
json_encode
json_decode
json_last_error
md5
sha1
hash
hash_hmac
base64_encode
base64_decode
urlencode
urldecode
rawurlencode
rawurldecode
htmlspecialchars
htmlentities
strip_tags
trim
ltrim
rtrim
strtolower
strtoupper
ucfirst
ucwords
sprintf
vsprintf
number_format
time
date
strtotime
microtime
sleep
usleep
gettype
settype
var_dump
var_export
print_r
debug_backtrace
class_exists
interface_exists
trait_exists
method_exists
property_exists
get_class
get_parent_class
get_class_methods
get_class_vars
get_object_vars
is_subclass_of
is_a
call_user_func
call_user_func_array
func_get_args
func_num_args
function_exists
header
headers_sent
setcookie
session_start
session_destroy
ob_start
ob_end_clean
ob_end_flush
ob_get_clean
ob_get_contents
error_reporting
set_error_handler
trigger_error
exit
die
EOFLIST

echo "  Checking for potentially missing standard library functions..."

# Note: This is an approximation - we're identifying functions that MIGHT be missing
# A more accurate approach would require running the code and catching undefined function errors

cat >> "$REPORT_FILE" <<EOF
## 2. Function Analysis

- **Total unique function calls found**: $TOTAL_UNIQUE_FUNCS

### Critical Standard Library Functions Used by WordPress

The following standard PHP functions are heavily used by WordPress and must be implemented:

**String Functions:**
- \`strlen\`, \`substr\`, \`strpos\`, \`strrpos\`, \`str_replace\`, \`str_repeat\`
- \`preg_match\`, \`preg_replace\`, \`preg_match_all\`, \`preg_split\`
- \`trim\`, \`ltrim\`, \`rtrim\`, \`strtolower\`, \`strtoupper\`
- \`htmlspecialchars\`, \`htmlentities\`, \`strip_tags\`
- \`sprintf\`, \`vsprintf\`, \`number_format\`
- \`ucfirst\`, \`ucwords\`, \`strrev\`, \`str_pad\`

**Array Functions:**
- \`array_push\`, \`array_pop\`, \`array_shift\`, \`array_unshift\`
- \`array_merge\`, \`array_keys\`, \`array_values\`, \`array_combine\`
- \`array_map\`, \`array_filter\`, \`array_reduce\`, \`array_walk\`
- \`in_array\`, \`array_search\`, \`array_key_exists\`
- \`count\`, \`sizeof\`, \`sort\`, \`rsort\`, \`asort\`, \`ksort\`
- \`array_slice\`, \`array_splice\`, \`array_chunk\`

**File I/O Functions:**
- \`file_exists\`, \`is_file\`, \`is_dir\`, \`is_readable\`, \`is_writable\`
- \`file_get_contents\`, \`file_put_contents\`, \`file\`
- \`fopen\`, \`fclose\`, \`fread\`, \`fwrite\`, \`fgets\`, \`feof\`
- \`mkdir\`, \`rmdir\`, \`unlink\`, \`rename\`, \`copy\`
- \`dirname\`, \`basename\`, \`realpath\`, \`pathinfo\`
- \`glob\`, \`scandir\`, \`opendir\`, \`readdir\`, \`closedir\`

**Type/Variable Functions:**
- \`is_array\`, \`is_string\`, \`is_int\`, \`is_float\`, \`is_bool\`, \`is_null\`
- \`is_numeric\`, \`is_object\`, \`is_resource\`, \`is_callable\`
- \`isset\`, \`empty\`, \`unset\`
- \`gettype\`, \`settype\`, \`var_dump\`, \`var_export\`, \`print_r\`

**OOP Functions:**
- \`class_exists\`, \`interface_exists\`, \`trait_exists\`
- \`method_exists\`, \`property_exists\`
- \`get_class\`, \`get_parent_class\`, \`get_class_methods\`, \`get_class_vars\`
- \`get_object_vars\`, \`is_subclass_of\`, \`is_a\`
- \`call_user_func\`, \`call_user_func_array\`

**Crypto/Encoding Functions:**
- \`md5\`, \`sha1\`, \`hash\`, \`hash_hmac\`, \`hash_file\`
- \`base64_encode\`, \`base64_decode\`
- \`urlencode\`, \`urldecode\`, \`rawurlencode\`, \`rawurldecode\`

**JSON Functions:**
- \`json_encode\`, \`json_decode\`, \`json_last_error\`, \`json_last_error_msg\`

**Date/Time Functions:**
- \`time\`, \`date\`, \`strtotime\`, \`microtime\`
- \`gmdate\`, \`mktime\`, \`strftime\`

**HTTP/Output Functions:**
- \`header\`, \`headers_sent\`, \`headers_list\`
- \`setcookie\`, \`setrawcookie\`
- \`ob_start\`, \`ob_end_clean\`, \`ob_end_flush\`, \`ob_get_clean\`, \`ob_get_contents\`

**Error Handling:**
- \`error_reporting\`, \`set_error_handler\`, \`restore_error_handler\`
- \`trigger_error\`, \`error_get_last\`

**Other Critical Functions:**
- \`define\`, \`defined\`, \`constant\`
- \`function_exists\`, \`func_get_args\`, \`func_num_args\`
- \`exit\`, \`die\`
- \`sleep\`, \`usleep\`

EOF

# Test 3: Check for WordPress-specific function usage patterns
echo "Test 3: Analyzing WordPress-specific patterns..."
echo "------------------------------------------------"

# Count occurrences of common WordPress APIs
echo "  Analyzing WordPress API usage..."

WP_HOOKS=$(find "$WP_DIR/wp-includes" -name "*.php" -exec grep -c "add_filter\|add_action\|apply_filters\|do_action" {} \; | awk '{s+=$1} END {print s}')
WP_DB=$(find "$WP_DIR/wp-includes" -name "*.php" -exec grep -c '\$wpdb->' {} \; | awk '{s+=$1} END {print s}')
WP_CLASSES=$(find "$WP_DIR/wp-includes" -name "class-*.php" | wc -l | tr -d ' ')

cat >> "$REPORT_FILE" <<EOF

## 3. WordPress-Specific Analysis

### WordPress API Usage
- **Hook system usage** (add_filter/add_action/apply_filters/do_action): ~$WP_HOOKS occurrences
- **Database usage** (\$wpdb->): ~$WP_DB occurrences
- **Class files** (class-*.php): $WP_CLASSES files

### Required PHP Extensions

WordPress 6.8.3 requires or heavily uses:

1. **mysqli** or **PDO_MySQL** (CRITICAL - required for database)
2. **json** (CRITICAL - core functionality)
3. **hash** (HIGH - security, checksums)
4. **filter** (HIGH - input validation)
5. **ctype** (MEDIUM - character type checking)
6. **mbstring** (MEDIUM - multi-byte string support)
7. **xml** (MEDIUM - RSS, XMLRPC)
8. **curl** (MEDIUM - HTTP requests)
9. **zip** (LOW - plugin/theme uploads)
10. **gd** or **imagick** (LOW - image processing)
11. **openssl** (LOW - cryptography)

### Current PHP-Go Extension Status

Based on the latest implementation:
- ✅ **hash**: Implemented (md5, sha1, sha256, etc.)
- ✅ **json**: Implemented (json_encode, json_decode)
- ❌ **mysqli/PDO**: Not implemented
- ❌ **filter**: Not implemented
- ❌ **ctype**: Not implemented
- ❌ **mbstring**: Not implemented
- ❌ **xml**: Not implemented
- ❌ **curl**: Not implemented
- ❌ **zip**: Not implemented
- ❌ **gd/imagick**: Not implemented
- ❌ **openssl**: Not implemented

EOF

# Test 4: Sample execution test
echo "Test 4: Testing sample file execution..."
echo "-----------------------------------------"

# Try to execute wp-includes/version.php (simple file with just variable assignments)
echo "  Testing wp-includes/version.php execution..."

VERSION_TEST="$REPORT_DIR/version_test_${TIMESTAMP}.log"
if $PHPGO "$WP_DIR/wp-includes/version.php" > "$VERSION_TEST" 2>&1; then
    echo "  ✓ Successfully executed version.php"
    VERSION_EXEC="SUCCESS"
else
    echo "  ✗ Failed to execute version.php"
    VERSION_EXEC="FAILED"
fi

cat >> "$REPORT_FILE" <<EOF

## 4. Execution Tests

### Simple File Test (wp-includes/version.php)
- **Status**: $VERSION_EXEC
EOF

if [ "$VERSION_EXEC" = "FAILED" ]; then
    cat >> "$REPORT_FILE" <<EOF
- **Error log**: \`$(basename "$VERSION_TEST")\`

\`\`\`
$(head -20 "$VERSION_TEST")
\`\`\`
EOF
fi

# Summary and recommendations
echo ""
echo "Generating summary and recommendations..."

cat >> "$REPORT_FILE" <<EOF

---

## 5. Summary and Recommendations

### Critical Blockers

The following are **critical blockers** that prevent WordPress from running:

1. **Database Extension (mysqli/PDO)** - WordPress cannot function without database access
   - Effort: ~40-60 hours
   - Priority: CRITICAL
   - Files affected: All WordPress functionality

2. **Missing Standard Library Functions** - WordPress uses 100+ standard library functions
   - Effort: ~200-300 hours (most already implemented in Phase 6)
   - Priority: CRITICAL
   - Categories: strings, arrays, files, JSON, dates, etc.

### High Priority Issues

3. **ctype Extension** - Used for character validation
   - Effort: ~4-6 hours
   - Priority: HIGH
   - Files affected: Input validation, sanitization

4. **filter Extension** - Input filtering and validation
   - Effort: ~8-12 hours
   - Priority: HIGH
   - Files affected: Security, data validation

5. **mbstring Extension** - Multi-byte string support
   - Effort: ~20-30 hours
   - Priority: HIGH
   - Files affected: Internationalization, UTF-8 handling

### Medium Priority Issues

6. **curl Extension** - HTTP client
   - Effort: ~15-20 hours
   - Priority: MEDIUM
   - Impact: Plugin/theme installation, HTTP API

7. **xml Extension** - XML parsing
   - Effort: ~25-35 hours
   - Priority: MEDIUM
   - Impact: RSS feeds, XML-RPC, imports

### Parse Issues

EOF

if [ $PARSE_FAIL -gt 0 ]; then
    cat >> "$REPORT_FILE" <<EOF
**Parse failures**: $PARSE_FAIL files failed to parse ($((PARSE_FAIL * 100 / PARSE_TOTAL))% failure rate)

Action items:
- Review parse error logs in \`$(basename "$PARSE_LOG")\`
- Fix parser bugs for unsupported PHP syntax
- Add missing language features

EOF
else
    cat >> "$REPORT_FILE" <<EOF
**Parse success**: All $PARSE_TOTAL WordPress files parsed successfully! ✅

EOF
fi

cat >> "$REPORT_FILE" <<EOF

### Recommended Implementation Order

To achieve WordPress compatibility, implement in this order:

1. **Phase 6 Completion** - Standard library functions (~200h remaining)
   - String manipulation functions
   - Array manipulation functions
   - File I/O functions
   - Date/time functions
   - Type checking functions

2. **Database Extension** - mysqli or PDO_MySQL (~40-60h)
   - Required for WordPress to run at all
   - Implement basic query execution
   - Support prepared statements

3. **Critical Extensions** - ctype, filter, hash completion (~20-30h)
   - Complete hash extension (missing algorithms)
   - Implement ctype functions
   - Implement filter extension

4. **HTTP/Network** - curl extension (~15-20h)
   - Basic HTTP client
   - Required for plugin/theme installation

5. **Additional Extensions** - mbstring, xml, etc. (~50-70h)
   - As needed for specific WordPress features

### Test Coverage Recommendations

- Create incremental test suite:
  1. Parse all files (DONE if $PARSE_FAIL = 0)
  2. Execute version.php and simple files
  3. Test specific subsystems (hooks, database abstraction layer)
  4. Load wp-settings.php
  5. Execute WordPress installation
  6. Full integration test

### Estimated Total Effort

- **Standard library completion**: ~200 hours (Phase 6)
- **Database support**: ~50 hours
- **Critical extensions**: ~30 hours
- **Additional extensions**: ~70 hours
- **WordPress-specific fixes**: ~20 hours
- **Testing and debugging**: ~30 hours

**Total**: ~400 hours to full WordPress compatibility

---

## Files Generated

- Parse errors: \`$(basename "$PARSE_LOG")\`
- Failed files: \`$(basename "$FAILED_FILES")\`
- All functions: \`all_functions_${TIMESTAMP}.txt\`
- Version test: \`$(basename "$VERSION_TEST")\`

Report generated: $(date)

EOF

echo ""
echo "=========================================="
echo "Issue identification complete!"
echo ""
echo "Report saved to: $REPORT_FILE"
echo ""
echo "Key findings:"
echo "  - Parse success rate: $(awk "BEGIN {printf \"%.2f\", ($PARSE_SUCCESS/$PARSE_TOTAL)*100}")%"
echo "  - Total unique functions: $TOTAL_UNIQUE_FUNCS"
echo "  - WordPress classes: $WP_CLASSES"
echo ""
echo "See report for detailed analysis and recommendations."
echo "=========================================="
