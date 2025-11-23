<?php
// Test: Nested if statements

$x = 10;
$y = 5;

if ($x > 5) {
    echo "x > 5\n";
    if ($y > 3) {
        echo "y > 3\n";
        if ($x > $y) {
            echo "x > y\n";
        }
    }
}

// Test if-else nesting
$a = 8;
if ($a > 10) {
    echo "a > 10\n";
} else {
    if ($a > 5) {
        echo "5 < a <= 10\n";
    } else {
        echo "a <= 5\n";
    }
}
