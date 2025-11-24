#!/bin/bash
# Symfony Installation Inventory Script
# Generates detailed statistics about the Symfony installation

cd "$(dirname "$0")/symfony-app" || exit 1

echo "============================================"
echo "Symfony Installation Inventory"
echo "============================================"
echo ""
echo "Date: $(date)"
echo ""

# Version information
echo "## Version Information"
echo -n "Symfony Framework: "
./bin/console --version 2>/dev/null | head -1 || echo "Unknown"
echo "PHP Required: $(grep '"php"' composer.json | cut -d'"' -f4)"
echo ""

# File statistics
echo "## File Statistics"
echo -n "Total PHP files: "
find . -name "*.php" -type f | wc -l | tr -d ' '

echo -n "Total PHP lines: "
find . -name "*.php" -type f -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}'

echo -n "Total directories: "
find . -type d | wc -l | tr -d ' '

echo ""

# Breakdown by directory
echo "## PHP Files by Directory"
for dir in bin config public src var vendor; do
    if [ -d "$dir" ]; then
        count=$(find "$dir" -name "*.php" -type f 2>/dev/null | wc -l | tr -d ' ')
        lines=$(find "$dir" -name "*.php" -type f -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}')
        echo "$dir: $count files, $lines lines"
    fi
done

echo ""

# Symfony packages
echo "## Symfony Packages"
echo -n "Symfony packages: "
grep -c '"name": "symfony/' composer.lock 2>/dev/null || echo "0"

echo ""
echo "Main Symfony components:"
grep '"name": "symfony/' composer.lock | cut -d'"' -f4 | head -15

echo ""

# Composer packages
echo "## All Installed Packages"
echo -n "Total packages: "
grep -c '"name":' composer.lock 2>/dev/null || echo "0"

echo ""

# Configuration
echo "## Configuration"
echo -n "Config files: "
find config -name "*.php" -o -name "*.yaml" -o -name "*.yml" 2>/dev/null | wc -l | tr -d ' '

echo -n "Routes: "
find config/routes -name "*.yaml" -o -name "*.yml" 2>/dev/null | wc -l | tr -d ' '

echo ""

# Key application files
echo "## Key Files"
echo "bin/console: $(if [ -f bin/console ]; then echo "✓"; else echo "✗"; fi)"
echo ".env: $(if [ -f .env ]; then echo "✓"; else echo "✗"; fi)"
echo "composer.json: $(if [ -f composer.json ]; then echo "✓"; else echo "✗"; fi)"
echo "composer.lock: $(if [ -f composer.lock ]; then echo "✓"; else echo "✗"; fi)"
echo "public/index.php: $(if [ -f public/index.php ]; then echo "✓"; else echo "✗"; fi)"

echo ""
echo "============================================"
echo "Inventory complete!"
echo "============================================"
