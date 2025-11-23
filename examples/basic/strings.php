<?php
// Example: String Operations
// Demonstrates string handling in PHP-Go

// String literals
$greeting = "Hello";
$name = "World";

echo "Basic strings:\n";
echo $greeting;
echo "\n";
echo $name;
echo "\n";

// String with escape sequences
$message = "Line 1\nLine 2\nLine 3";
echo "\nMulti-line string:\n";
echo $message;
echo "\n";

// Concatenation (Note: may have current limitations - see Issue #xxx)
// $combined = $greeting . " " . $name;
// echo "\nConcatenated: ";
// echo $combined;

// String length (built-in function)
echo "\nString length:\n";
echo "Length of 'Hello': ";
echo strlen("Hello");
echo "\n";
echo "Length of variable: ";
echo strlen($name);
echo "\n";

// Empty string
$empty = "";
echo "\nEmpty string length: ";
echo strlen($empty);
echo "\n";
