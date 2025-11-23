<?php
// Test: Function definitions and calls

function greet($name) {
    echo "Hello, ";
    echo $name;
    echo "!\n";
}

function add($a, $b) {
    return $a + $b;
}

function multiply($x, $y) {
    $result = $x * $y;
    return $result;
}

// Test function calls
greet("World");
greet("PHP-Go");

// Test return values
$sum = add(5, 3);
echo "5 + 3 = ";
echo $sum;
echo "\n";

$product = multiply(4, 7);
echo "4 * 7 = ";
echo $product;
echo "\n";

// Test nested calls
$result = add(multiply(2, 3), 4);
echo "2 * 3 + 4 = ";
echo $result;
echo "\n";
