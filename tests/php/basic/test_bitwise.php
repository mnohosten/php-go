<?php
// Test: Bitwise operations

// Test AND
$a = 12;  // 1100 in binary
$b = 10;  // 1010 in binary
echo "12 & 10 = ";
echo $a & $b;  // Should be 8 (1000)
echo "\n";

// Test OR
echo "12 | 10 = ";
echo $a | $b;  // Should be 14 (1110)
echo "\n";

// Test XOR
echo "12 ^ 10 = ";
echo $a ^ $b;  // Should be 6 (0110)
echo "\n";

// Test NOT
echo "~12 = ";
echo ~$a;  // Should be -13
echo "\n";

// Test left shift
echo "5 << 2 = ";
echo 5 << 2;  // Should be 20
echo "\n";

// Test right shift
echo "20 >> 2 = ";
echo 20 >> 2;  // Should be 5
echo "\n";

// Combined bitwise operations
$result = (15 & 7) | (3 << 2);
echo "(15 & 7) | (3 << 2) = ";
echo $result;
echo "\n";
