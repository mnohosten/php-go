#!/bin/bash
#
# WordPress Testing Suite for PHP-Go
#
# This script runs various WordPress compatibility tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
WP_DIR="$SCRIPT_DIR/wordpress"

echo "WordPress Testing Suite for PHP-Go"
echo "===================================="
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
    echo "Please run the installation first"
    exit 1
fi

echo "PHP-Go Binary: $PHPGO"
echo "WordPress Dir: $WP_DIR"
echo ""

# Test 1: Parse Test
echo "Test 1: Parse WordPress Files"
echo "------------------------------"
if [ -f "$SCRIPT_DIR/test-parse.php" ]; then
    $PHPGO "$SCRIPT_DIR/test-parse.php"
    echo ""
else
    echo "SKIP: test-parse.php not found"
    echo ""
fi

# Test 2: Version Detection
echo "Test 2: WordPress Version Detection"
echo "------------------------------------"
$PHPGO -c "<?php require_once '$WP_DIR/wp-includes/version.php'; echo 'WordPress Version: ' . \$wp_version . PHP_EOL;" 2>&1 || echo "FAIL: Could not detect WordPress version"
echo ""

# Test 3: Load Index (without execution)
echo "Test 3: Parse index.php"
echo "-----------------------"
$PHPGO parse "$WP_DIR/index.php" > /dev/null 2>&1 && echo "PASS: index.php parses successfully" || echo "FAIL: index.php parse failed"
echo ""

# Test 4: Load wp-blog-header (without execution)
echo "Test 4: Parse wp-blog-header.php"
echo "---------------------------------"
$PHPGO parse "$WP_DIR/wp-blog-header.php" > /dev/null 2>&1 && echo "PASS: wp-blog-header.php parses successfully" || echo "FAIL: wp-blog-header.php parse failed"
echo ""

echo "===================================="
echo "WordPress testing complete"
echo "See README.md for next steps"
