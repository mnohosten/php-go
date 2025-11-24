#!/bin/bash
#
# Symfony Test Runner for PHP-Go
# Quick test script for critical Symfony files
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
SYMFONY_DIR="$SCRIPT_DIR/symfony-app"

echo "Symfony Test Runner for PHP-Go"
echo "==============================="
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
    echo "Please run: composer create-project symfony/skeleton symfony-app"
    exit 1
fi

# Get Symfony version
echo "Symfony Version: $($SYMFONY_DIR/bin/console --version 2>/dev/null | head -1 || echo "Unknown")"
echo "PHP-Go Version: $($PHPGO --version 2>&1 || echo "unknown")"
echo ""

# Critical files to test
declare -a CRITICAL_FILES=(
    "public/index.php"
    "bin/console"
    "src/Kernel.php"
    "config/bundles.php"
)

PASSED=0
FAILED=0

echo "Testing critical Symfony files..."
echo "-----------------------------------"

for file in "${CRITICAL_FILES[@]}"; do
    full_path="$SYMFONY_DIR/$file"

    if [ ! -f "$full_path" ]; then
        echo "⏭  $file (not found)"
        continue
    fi

    echo -n "Testing: $file ... "

    if $PHPGO parse "$full_path" > /dev/null 2>&1; then
        echo "✓ PASS"
        PASSED=$((PASSED + 1))
    else
        echo "✗ FAIL"
        FAILED=$((FAILED + 1))
        # Show error details
        echo "  Error details:"
        $PHPGO parse "$full_path" 2>&1 | sed 's/^/    /'
    fi
done

echo ""
echo "-----------------------------------"
echo "Results:"
echo "  Passed: $PASSED"
echo "  Failed: $FAILED"
echo "  Total:  $((PASSED + FAILED))"
echo "==============================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
