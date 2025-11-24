#!/bin/bash
#
# WordPress Installation Inventory
#
# This script provides statistics about the WordPress installation

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WP_DIR="$SCRIPT_DIR/wordpress"

echo "WordPress Installation Inventory"
echo "================================="
echo ""

# Check WordPress version
if [ -f "$WP_DIR/wp-includes/version.php" ]; then
    VERSION=$(grep '$wp_version =' "$WP_DIR/wp-includes/version.php" | sed "s/.*'\(.*\)'.*/\1/")
    echo "WordPress Version: $VERSION"
fi

echo ""
echo "File Statistics:"
echo "----------------"

# Count PHP files
PHP_FILES=$(find "$WP_DIR" -name "*.php" | wc -l | tr -d ' ')
echo "Total PHP files: $PHP_FILES"

# Count lines of PHP code
PHP_LINES=$(find "$WP_DIR" -name "*.php" -exec cat {} \; | wc -l | tr -d ' ')
echo "Total PHP lines: $PHP_LINES"

# Major directories
echo ""
echo "Directory Structure:"
echo "--------------------"
echo "wp-admin:    $(find "$WP_DIR/wp-admin" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') PHP files"
echo "wp-includes: $(find "$WP_DIR/wp-includes" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') PHP files"
echo "wp-content:  $(find "$WP_DIR/wp-content" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') PHP files"

# Top-level files
echo ""
echo "Core Files:"
echo "-----------"
ls -lh "$WP_DIR"/*.php 2>/dev/null | awk '{print $9, "(" $5 ")"}'

# Classes
echo ""
echo "Class Files:"
echo "------------"
CLASS_FILES=$(find "$WP_DIR/wp-includes" -name "class-*.php" | wc -l | tr -d ' ')
echo "Class files in wp-includes: $CLASS_FILES"

echo ""
echo "Key Subsystems:"
echo "---------------"
echo "REST API:    $(find "$WP_DIR/wp-includes/rest-api" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') files"
echo "Blocks:      $(find "$WP_DIR/wp-includes/blocks" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') files"
echo "Widgets:     $(find "$WP_DIR/wp-includes/widgets" -name "*.php" 2>/dev/null | wc -l | tr -d ' ') files"

echo ""
echo "================================="
