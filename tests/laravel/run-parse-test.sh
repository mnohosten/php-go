#!/bin/bash
# Laravel Parse Test - Test critical files

PHPGO="../../php-go"
FILES=(
    "artisan"
    "public/index.php"
    "bootstrap/app.php"
    "bootstrap/providers.php"
    "app/Console/Kernel.php"
    "app/Exceptions/Handler.php"
    "app/Models/User.php"
    "config/app.php"
    "config/database.php"
    "routes/web.php"
    "routes/console.php"
)

echo "Laravel Critical Files Parse Test"
echo "=================================="
echo ""

total=0
passed=0
failed=0

for file in "${FILES[@]}"; do
    total=$((total + 1))
    echo -n "Testing: $file ... "
    
    if [ ! -f "$file" ]; then
        echo "SKIP (not found)"
        continue
    fi
    
    if $PHPGO parse "$file" > /dev/null 2>&1; then
        echo "PASS"
        passed=$((passed + 1))
    else
        echo "FAIL"
        failed=$((failed + 1))
    fi
done

echo ""
echo "=================================="
echo "Results: $passed passed, $failed failed, $total total"
if [ $total -gt 0 ]; then
    success_rate=$(echo "scale=2; ($passed / $total) * 100" | bc)
    echo "Success Rate: $success_rate%"
fi
