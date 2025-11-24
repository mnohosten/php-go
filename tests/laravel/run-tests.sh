#!/bin/bash
#
# Laravel Testing Suite for PHP-Go
#
# This script runs various Laravel compatibility tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHPGO="$SCRIPT_DIR/../../php-go"
LARAVEL_DIR="$SCRIPT_DIR"

echo "Laravel Testing Suite for PHP-Go"
echo "=================================="
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
    echo "Please run the installation first"
    exit 1
fi

echo "PHP-Go Binary: $PHPGO"
echo "Laravel Dir: $LARAVEL_DIR"
echo ""

# Test 1: Parse Test
echo "Test 1: Parse Laravel Files"
echo "----------------------------"
if [ -f "$SCRIPT_DIR/test-parse.php" ]; then
    $PHPGO "$SCRIPT_DIR/test-parse.php"
    echo ""
else
    echo "SKIP: test-parse.php not found"
    echo ""
fi

# Test 2: Version Detection
echo "Test 2: Laravel Version Detection"
echo "----------------------------------"
$PHPGO -c "<?php require_once '$LARAVEL_DIR/vendor/autoload.php'; echo 'Laravel Version: ' . \Illuminate\Foundation\Application::VERSION . PHP_EOL;" 2>&1 || echo "FAIL: Could not detect Laravel version"
echo ""

# Test 3: Parse artisan
echo "Test 3: Parse artisan"
echo "---------------------"
$PHPGO parse "$LARAVEL_DIR/artisan" > /dev/null 2>&1 && echo "PASS: artisan parses successfully" || echo "FAIL: artisan parse failed"
echo ""

# Test 4: Parse public/index.php
echo "Test 4: Parse public/index.php"
echo "-------------------------------"
$PHPGO parse "$LARAVEL_DIR/public/index.php" > /dev/null 2>&1 && echo "PASS: public/index.php parses successfully" || echo "FAIL: public/index.php parse failed"
echo ""

# Test 5: Parse bootstrap/app.php
echo "Test 5: Parse bootstrap/app.php"
echo "--------------------------------"
$PHPGO parse "$LARAVEL_DIR/bootstrap/app.php" > /dev/null 2>&1 && echo "PASS: bootstrap/app.php parses successfully" || echo "FAIL: bootstrap/app.php parse failed"
echo ""

echo "=================================="
echo "Laravel testing complete"
echo "See LARAVEL_README.md for next steps"
