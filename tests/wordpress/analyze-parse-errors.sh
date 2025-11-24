#!/bin/bash
#
# Analyze parse errors from WordPress files
#

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
WP_DIR="$SCRIPT_DIR/wordpress"
FAILED_FILES="$SCRIPT_DIR/reports/failed_files_20251124_024134.txt"
OUTPUT="$SCRIPT_DIR/reports/parse_error_analysis.txt"

echo "Analyzing parse errors from WordPress files..."
echo ""

# Extract all unique "no prefix parse function for X" errors
echo "=== Missing Prefix Parse Functions ===" > "$OUTPUT"
while IFS= read -r file; do
    $PHPGO parse "$file" 2>&1 | grep -oE "no prefix parse function for [A-Z_]+"
done < "$FAILED_FILES" | sort | uniq -c | sort -rn >> "$OUTPUT"

echo "" >> "$OUTPUT"

# Extract all unique "expected X got Y" errors
echo "=== Unexpected Token Errors ===" >> "$OUTPUT"
while IFS= read -r file; do
    $PHPGO parse "$file" 2>&1 | grep -oE "expected next token to be [^,]+, got [^ ]+ instead"
done < "$FAILED_FILES" | sort | uniq -c | sort -rn >> "$OUTPUT"

echo "" >> "$OUTPUT"

# Extract other error patterns
echo "=== Other Parse Errors ===" >> "$OUTPUT"
while IFS= read -r file; do
    $PHPGO parse "$file" 2>&1 | grep "Parse error:" | grep -v "no prefix parse function" | grep -v "expected next token"
done < "$FAILED_FILES" | sed 's/.*Parse error: //' | sort | uniq -c | sort -rn >> "$OUTPUT"

echo "Analysis complete!"
echo ""
cat "$OUTPUT"
