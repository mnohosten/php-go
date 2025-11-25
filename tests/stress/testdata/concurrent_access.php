<?php
// Script for testing concurrent access
// Simulates operations that might have race conditions

// Simulate operations
$counter = 0;
for ($i = 0; $i < 100; $i++) {
    $counter++;
}

// Array operations
$arr = [];
for ($i = 0; $i < 50; $i++) {
    $arr[$i] = $i * 2;
}

$result = 0;
foreach ($arr as $value) {
    $result += $value;
}

echo "Counter: $counter, Array sum: $result\n";
