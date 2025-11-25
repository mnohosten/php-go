#!/bin/bash

# WordPress/Laravel Performance Testing Script
# Compares php-go vs PHP 8.4 on framework-specific workloads

set -e

PHPGO="./php-go"
PHP="/opt/homebrew/bin/php"
OUTPUT_DIR="benchmarks/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "==================================================================="
echo "WordPress/Laravel Performance Benchmarking"
echo "==================================================================="
echo "Date: $(date)"
echo "PHP Version: $($PHP --version | head -1)"
echo "PHP-Go: Built from source"
echo "==================================================================="
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to time execution
benchmark() {
    local name=$1
    local script=$2
    local iterations=${3:-10}

    echo -e "${YELLOW}Benchmarking: $name${NC}"

    # Create temporary file for php-go (it doesn't support stdin)
    local tmpfile=$(mktemp /tmp/phpgo_bench.XXXXXX.php)
    echo "$script" > "$tmpfile"

    # Time PHP
    local php_total=0
    for i in $(seq 1 $iterations); do
        local start=$(perl -MTime::HiRes -e 'print Time::HiRes::time()')
        $PHP -r "$script" > /dev/null 2>&1 || true
        local end=$(perl -MTime::HiRes -e 'print Time::HiRes::time()')
        local duration=$(echo "$end - $start" | bc)
        php_total=$(echo "$php_total + $duration" | bc)
    done
    local php_avg=$(echo "scale=6; $php_total / $iterations" | bc)

    # Time php-go
    local phpgo_total=0
    local phpgo_success=0
    for i in $(seq 1 $iterations); do
        local start=$(perl -MTime::HiRes -e 'print Time::HiRes::time()')
        $PHPGO "$tmpfile" > /dev/null 2>&1 && phpgo_success=$((phpgo_success + 1)) || true
        local end=$(perl -MTime::HiRes -e 'print Time::HiRes::time()')
        local duration=$(echo "$end - $start" | bc)
        phpgo_total=$(echo "$phpgo_total + $duration" | bc)
    done
    local phpgo_avg=$(echo "scale=6; $phpgo_total / $iterations" | bc)

    # Clean up
    rm -f "$tmpfile"

    # Calculate ratio
    local ratio=$(echo "scale=2; $phpgo_avg / $php_avg" | bc)

    # Determine status
    local status="✗"
    if [ $phpgo_success -eq $iterations ]; then
        status="✓"
    fi

    # Output results
    printf "  %-30s PHP: %8.6fs | PHP-Go: %8.6fs | Ratio: %6.2fx | Status: %s (%d/%d)\n" \
        "$name" "$php_avg" "$phpgo_avg" "$ratio" "$status" "$phpgo_success" "$iterations"

    # Log to file
    echo "$TIMESTAMP,$name,$php_avg,$phpgo_avg,$ratio,$phpgo_success,$iterations" >> "$OUTPUT_DIR/performance_log.csv"
}

# Initialize CSV log
if [ ! -f "$OUTPUT_DIR/performance_log.csv" ]; then
    echo "timestamp,test_name,php_avg,phpgo_avg,ratio,phpgo_success,iterations" > "$OUTPUT_DIR/performance_log.csv"
fi

echo ""
echo "=== Basic Operations ==="
echo ""

benchmark "Simple Loop (10k iterations)" \
'<?php
$a = 0;
for ($i = 0; $i < 10000; $i++) {
    $a++;
}
echo $a;
' 20

benchmark "Function Calls (5k calls)" \
'<?php
function test($x) { return $x; }
for ($i = 0; $i < 5000; $i++) {
    test($i);
}
echo "OK";
' 20

benchmark "String Operations" \
'<?php
$s = "";
for ($i = 0; $i < 100; $i++) {
    $s .= "test";
}
echo strlen($s);
' 20

benchmark "Array Operations" \
'<?php
$arr = [];
for ($i = 0; $i < 1000; $i++) {
    $arr[$i] = $i * 2;
}
echo "OK";
' 20

echo ""
echo "=== WordPress-like Patterns ==="
echo ""

benchmark "Nested Conditionals" \
'<?php
$result = 0;
for ($i = 0; $i < 1000; $i++) {
    if ($i % 2 == 0) {
        if ($i % 3 == 0) {
            $result++;
        } else {
            $result--;
        }
    } else {
        $result += 2;
    }
}
echo $result;
' 15

benchmark "WordPress Hook Pattern" \
'<?php
class HookSystem {
    private $hooks = [];

    public function add($name, $callback) {
        $this->hooks[$name][] = $callback;
    }

    public function run($name, $value) {
        if (isset($this->hooks[$name])) {
            foreach ($this->hooks[$name] as $callback) {
                $value = $callback($value);
            }
        }
        return $value;
    }
}

$hooks = new HookSystem();
$hooks->add("filter", function($x) { return $x + 1; });
$hooks->add("filter", function($x) { return $x * 2; });

$result = 0;
for ($i = 0; $i < 500; $i++) {
    $result = $hooks->run("filter", $i);
}
echo $result;
' 10

echo ""
echo "=== Laravel-like Patterns ==="
echo ""

benchmark "Class Instantiation" \
'<?php
class User {
    private $name;
    private $email;

    public function __construct($name, $email) {
        $this->name = $name;
        $this->email = $email;
    }

    public function getName() {
        return $this->name;
    }
}

for ($i = 0; $i < 1000; $i++) {
    $user = new User("Test", "test@example.com");
    $user->getName();
}
echo "OK";
' 10

benchmark "Array Map/Filter Pattern" \
'<?php
$data = [];
for ($i = 0; $i < 100; $i++) {
    $data[] = $i;
}

// Simple map
$result = [];
foreach ($data as $item) {
    $result[] = $item * 2;
}

// Simple filter
$filtered = [];
foreach ($result as $item) {
    if ($item % 4 == 0) {
        $filtered[] = $item;
    }
}

echo "OK";
' 15

benchmark "Method Chaining" \
'<?php
class Builder {
    private $value = 0;

    public function add($n) {
        $this->value += $n;
        return $this;
    }

    public function multiply($n) {
        $this->value *= $n;
        return $this;
    }

    public function get() {
        return $this->value;
    }
}

$result = 0;
for ($i = 0; $i < 500; $i++) {
    $builder = new Builder();
    $result = $builder->add(5)->multiply(3)->add(2)->get();
}
echo $result;
' 10

echo ""
echo "=== Complex Operations ==="
echo ""

benchmark "Recursive Function (Fibonacci)" \
'<?php
function fib($n) {
    if ($n < 2) return 1;
    return fib($n - 1) + fib($n - 2);
}
echo fib(15);
' 10

benchmark "OOP Inheritance" \
'<?php
class Animal {
    protected $name;

    public function __construct($name) {
        $this->name = $name;
    }

    public function speak() {
        return "Animal speaks";
    }
}

class Dog extends Animal {
    public function speak() {
        return "Woof!";
    }
}

for ($i = 0; $i < 500; $i++) {
    $dog = new Dog("Buddy");
    $dog->speak();
}
echo "OK";
' 10

echo ""
echo "==================================================================="
echo "Benchmark Summary"
echo "==================================================================="
echo ""
echo "Results saved to: $OUTPUT_DIR/performance_log.csv"
echo "View results with: cat $OUTPUT_DIR/performance_log.csv | tail -20"
echo ""
echo "==================================================================="
