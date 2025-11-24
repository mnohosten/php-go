#!/bin/bash

# Symfony Framework Test Suite Inventory
# Generates statistics about Symfony test files

BASE_DIR="/Users/krizos/code/mnohosten/php-go/tests/symfony/symfony-framework"

echo "=== Symfony Framework Test Suite Inventory ==="
echo ""
echo "Date: $(date)"
echo "Repository: symfony/symfony v7.3"
echo ""

# Count total test files
total_tests=$(find "$BASE_DIR" -name "*Test.php" | wc -l | tr -d ' ')
echo "Total Test Files: $total_tests"
echo ""

# Count tests by category
echo "=== Tests by Category ==="
bridges=$(find "$BASE_DIR/src/Symfony/Bridge" -name "*Test.php" 2>/dev/null | wc -l | tr -d ' ')
components=$(find "$BASE_DIR/src/Symfony/Component" -name "*Test.php" 2>/dev/null | wc -l | tr -d ' ')
contracts=$(find "$BASE_DIR/src/Symfony/Contracts" -name "*Test.php" 2>/dev/null | wc -l | tr -d ' ')
bundles=$(find "$BASE_DIR/src/Symfony/Bundle" -name "*Test.php" 2>/dev/null | wc -l | tr -d ' ')

echo "Bridge Tests: $bridges"
echo "Component Tests: $components"
echo "Contract Tests: $contracts"
echo "Bundle Tests: $bundles"
echo ""

# Top 10 components by test count
echo "=== Top 10 Components by Test Count ==="
find "$BASE_DIR/src/Symfony/Component" -type d -name "Tests" | while read dir; do
    component=$(echo "$dir" | sed 's|.*/Component/\([^/]*\)/.*|\1|')
    count=$(find "$dir" -name "*Test.php" | wc -l | tr -d ' ')
    echo "$count $component"
done | sort -rn | head -10
echo ""

# Sample a few components for detailed analysis
echo "=== Sample Components (detailed) ==="
for component in Console DependencyInjection HttpFoundation HttpKernel; do
    comp_dir="$BASE_DIR/src/Symfony/Component/$component"
    if [ -d "$comp_dir/Tests" ]; then
        test_count=$(find "$comp_dir/Tests" -name "*Test.php" | wc -l | tr -d ' ')
        php_files=$(find "$comp_dir" -name "*.php" | wc -l | tr -d ' ')
        lines=$(find "$comp_dir" -name "*.php" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}')
        echo "$component: $test_count tests, $php_files total files, ~$lines LOC"
    fi
done
echo ""

echo "=== Test Suite Structure ==="
echo "Tests are organized by:"
echo "- Bridge: Integration with 3rd-party libraries"
echo "- Component: Core Symfony components (57 total)"
echo "- Contract: Interface definitions and contracts"
echo "- Bundle: Full-stack bundles"
