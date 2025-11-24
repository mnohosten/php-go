#!/bin/bash
# Laravel Installation Inventory Script
# Generates detailed statistics about the Laravel installation

echo "============================================"
echo "Laravel Installation Inventory"
echo "============================================"
echo ""
echo "Date: $(date)"
echo ""

# Version information
echo "## Version Information"
echo "Laravel Framework: $(grep -A 2 '"laravel/framework"' composer.lock | grep version | cut -d'"' -f4)"
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
for dir in app bootstrap config database public resources routes storage tests vendor; do
    if [ -d "$dir" ]; then
        count=$(find "$dir" -name "*.php" -type f 2>/dev/null | wc -l | tr -d ' ')
        lines=$(find "$dir" -name "*.php" -type f -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}')
        echo "$dir: $count files, $lines lines"
    fi
done

echo ""

# Composer packages
echo "## Installed Packages"
echo -n "Total packages: "
grep -c '"name":' composer.lock

echo ""
echo "Main dependencies:"
grep -A 1 '"name": "laravel/' composer.lock | grep '"name"' | cut -d'"' -f4 | head -10

echo ""

# Database
echo "## Database"
if [ -f "database/database.sqlite" ]; then
    size=$(ls -lh database/database.sqlite | awk '{print $5}')
    echo "SQLite database: database/database.sqlite ($size)"
else
    echo "SQLite database: Not found"
fi

echo ""

# Configuration
echo "## Configuration"
echo -n "Config files: "
find config -name "*.php" -type f | wc -l | tr -d ' '

echo -n "Routes: "
find routes -name "*.php" -type f | wc -l | tr -d ' '

echo -n "Migrations: "
find database/migrations -name "*.php" -type f 2>/dev/null | wc -l | tr -d ' '

echo ""

# Key application files
echo "## Key Files"
echo "Artisan: $(if [ -f artisan ]; then echo "✓"; else echo "✗"; fi)"
echo ".env: $(if [ -f .env ]; then echo "✓"; else echo "✗"; fi)"
echo "composer.json: $(if [ -f composer.json ]; then echo "✓"; else echo "✗"; fi)"
echo "composer.lock: $(if [ -f composer.lock ]; then echo "✓"; else echo "✗"; fi)"

echo ""
echo "============================================"
echo "Inventory complete!"
echo "============================================"
