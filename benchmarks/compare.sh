#!/bin/bash

# compare.sh - Compare php-go performance against standard PHP
# Usage: ./compare.sh [script.php]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if PHP is installed
if ! command -v php &> /dev/null; then
    echo -e "${RED}Error: PHP not found in PATH${NC}"
    echo "Please install PHP to run comparisons"
    exit 1
fi

# Check if php-go is built
if [ ! -f "../php-go" ]; then
    echo -e "${YELLOW}Building php-go...${NC}"
    (cd .. && go build -o php-go ./cmd/php-go)
fi

PHP_GO="../php-go"

# Default test script
SCRIPT="${1:-scripts/simple_bench.php}"

if [ ! -f "$SCRIPT" ]; then
    echo -e "${RED}Error: Script not found: $SCRIPT${NC}"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}PHP-Go vs PHP Performance Comparison${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Script: $SCRIPT"
echo "PHP Version: $(php -v | head -1)"
echo ""

# Function to run benchmark
run_benchmark() {
    local engine=$1
    local command=$2
    local iterations=5

    echo -e "${YELLOW}Running $engine (${iterations} iterations)...${NC}"

    local times=()
    local outputs=()

    for i in $(seq 1 $iterations); do
        # Measure time using bash's built-in time
        START=$(date +%s%N)
        OUTPUT=$(eval "$command" 2>&1)
        END=$(date +%s%N)

        # Calculate duration in milliseconds
        DURATION=$(echo "scale=3; ($END - $START) / 1000000" | bc)
        times+=($DURATION)

        if [ $i -eq 1 ]; then
            outputs+=("$OUTPUT")
        fi

        echo "  Iteration $i: ${DURATION}ms"
    done

    # Calculate average
    local sum=0
    for t in "${times[@]}"; do
        sum=$(echo "$sum + $t" | bc)
    done
    local avg=$(echo "scale=3; $sum / $iterations" | bc)

    # Calculate min and max
    local min=${times[0]}
    local max=${times[0]}
    for t in "${times[@]}"; do
        if (( $(echo "$t < $min" | bc -l) )); then
            min=$t
        fi
        if (( $(echo "$t > $max" | bc -l) )); then
            max=$t
        fi
    done

    echo ""
    echo "$avg|$min|$max|${outputs[0]}"
}

# Run PHP benchmark
echo ""
php_result=$(run_benchmark "PHP" "php $SCRIPT")
php_avg=$(echo "$php_result" | cut -d'|' -f1)
php_min=$(echo "$php_result" | cut -d'|' -f2)
php_max=$(echo "$php_result" | cut -d'|' -f3)
php_output=$(echo "$php_result" | cut -d'|' -f4-)

echo -e "${GREEN}PHP Results:${NC}"
echo "  Average: ${php_avg}ms"
echo "  Min: ${php_min}ms"
echo "  Max: ${php_max}ms"
echo ""

# Run php-go benchmark
phpgo_result=$(run_benchmark "php-go" "$PHP_GO run $SCRIPT")
phpgo_avg=$(echo "$phpgo_result" | cut -d'|' -f1)
phpgo_min=$(echo "$phpgo_result" | cut -d'|' -f2)
phpgo_max=$(echo "$phpgo_result" | cut -d'|' -f3)
phpgo_output=$(echo "$phpgo_result" | cut -d'|' -f4-)

echo -e "${GREEN}php-go Results:${NC}"
echo "  Average: ${phpgo_avg}ms"
echo "  Min: ${phpgo_min}ms"
echo "  Max: ${phpgo_max}ms"
echo ""

# Calculate ratio
ratio=$(echo "scale=2; $phpgo_avg / $php_avg" | bc)

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Comparison${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Performance Ratio (php-go/PHP): ${YELLOW}${ratio}x${NC}"

if (( $(echo "$ratio < 1.0" | bc -l) )); then
    speedup=$(echo "scale=1; (1 - $ratio) * 100" | bc)
    echo -e "${GREEN}✓ php-go is FASTER by ${speedup}%${NC}"
elif (( $(echo "$ratio < 2.0" | bc -l) )); then
    slowdown=$(echo "scale=1; ($ratio - 1) * 100" | bc)
    echo -e "${YELLOW}○ php-go is ${slowdown}% slower (within 2x target)${NC}"
else
    slowdown=$(echo "scale=1; ($ratio - 1) * 100" | bc)
    echo -e "${RED}✗ php-go is ${slowdown}% slower (exceeds 2x target)${NC}"
fi

echo ""

# Compare outputs
if [ "$php_output" = "$phpgo_output" ]; then
    echo -e "${GREEN}✓ Output matches exactly${NC}"
else
    echo -e "${RED}✗ Output differs${NC}"
    echo ""
    echo "PHP output (first 200 chars):"
    echo "$php_output" | head -c 200
    echo ""
    echo ""
    echo "php-go output (first 200 chars):"
    echo "$phpgo_output" | head -c 200
    echo ""
fi

echo ""
echo -e "${BLUE}========================================${NC}"
