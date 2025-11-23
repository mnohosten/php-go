<?php
// Test: Arithmetic operations including negate and modulo

// Test negation
$x = 5;
$y = -$x;
echo "Negate 5: ";
echo $y;
echo "\n";

// Test modulo
$a = 10;
$b = 3;
echo "10 % 3 = ";
echo $a % $b;
echo "\n";

// Test power
$p = 2 ** 8;
echo "2 ** 8 = ";
echo $p;
echo "\n";

// Test negative modulo
$c = -10;
$d = 3;
echo "-10 % 3 = ";
echo $c % $d;
echo "\n";

// Test float modulo
$e = 10.5;
$f = 3.2;
echo "10.5 % 3.2 = ";
echo $e % $f;
echo "\n";

// Test division by zero protection
$g = 10;
$h = 2;
echo "10 / 2 = ";
echo $g / $h;
echo "\n";

// Combined operations
$result = 5 + 3 * 2 - 8 / 4;
echo "5 + 3 * 2 - 8 / 4 = ";
echo $result;
echo "\n";

// Test increment and decrement
$i = 10;
$i = $i + 1;
echo "10 + 1 = ";
echo $i;
echo "\n";

// Test with assignment
$k = 5;
$k = $k - 1;
echo "5 - 1 = ";
echo $k;
echo "\n";
