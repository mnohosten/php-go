#!/bin/bash
cd /Users/krizos/code/mnohosten/php-go
total=0
success=0
failed_files=()

for f in $(find tests/symfony/symfony-app -name "*.php" -type f); do
    total=$((total+1))
    if ./php-go parse "$f" >/dev/null 2>&1; then
        success=$((success+1))
    else
        failed_files+=("$f")
    fi
done

echo "Success: $success / $total"
rate=$(echo "scale=2; $success * 100 / $total" | bc)
echo "Rate: $rate%"
echo ""
echo "First 10 failed files:"
for i in {0..9}; do
    if [ $i -lt ${#failed_files[@]} ]; then
        echo "  ${failed_files[$i]}"
    fi
done
