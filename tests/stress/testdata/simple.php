<?php
// Simple PHP script for load testing
// Performs basic operations quickly

$sum = 0;
for ($i = 0; $i < 100; $i++) {
    $sum += $i;
}

echo "Sum: " . $sum . "\n";
