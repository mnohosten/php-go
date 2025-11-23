<?php
// Example: Arithmetic and Logical Expressions
// Demonstrates expression evaluation

// Arithmetic expressions
echo "Arithmetic:\n";
$a = 10;
$b = 3;

echo "Addition: ";
echo $a + $b;
echo "\n";

echo "Subtraction: ";
echo $a - $b;
echo "\n";

echo "Multiplication: ";
echo $a * $b;
echo "\n";

echo "Division: ";
echo $a / $b;
echo "\n";

echo "Modulo: ";
echo $a % $b;
echo "\n";

// Compound expressions
echo "\nCompound expression:\n";
$result = 5 + 3 * 2;
echo "5 + 3 * 2 = ";
echo $result;
echo "\n";

// Assignment with expression
echo "\nAssignment with expression:\n";
$x = 5;
$y = $x + 10;
echo "x = 5, y = x + 10: ";
echo $y;
echo "\n";

// Comparison (returns boolean)
echo "\nComparisons:\n";
echo "10 == 10: ";
echo 10 == 10;  // true (1)
echo "\n";

echo "10 > 5: ";
echo 10 > 5;    // true (1)
echo "\n";

echo "3 < 2: ";
echo 3 < 2;     // false (empty)
echo "\n";
