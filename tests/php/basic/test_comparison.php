<?php
// Test: Comparison operations

// Test equality (==)
$a = 5;
$b = 5;
$c = "5";
echo "5 == 5: ";
echo $a == $b ? "true" : "false";
echo "\n";

echo "5 == '5': ";
echo $a == $c ? "true" : "false";
echo "\n";

// Test identity (===)
echo "5 === 5: ";
echo $a === $b ? "true" : "false";
echo "\n";

echo "5 === '5': ";
echo $a === $c ? "true" : "false";
echo "\n";

// Test inequality (!=)
echo "5 != 3: ";
echo $a != 3 ? "true" : "false";
echo "\n";

// Test non-identity (!==)
echo "5 !== '5': ";
echo $a !== $c ? "true" : "false";
echo "\n";

// Test greater than
echo "10 > 5: ";
echo 10 > 5 ? "true" : "false";
echo "\n";

// Test less than
echo "3 < 7: ";
echo 3 < 7 ? "true" : "false";
echo "\n";

// Test greater than or equal
echo "5 >= 5: ";
echo 5 >= 5 ? "true" : "false";
echo "\n";

// Test less than or equal
echo "4 <= 9: ";
echo 4 <= 9 ? "true" : "false";
echo "\n";

// Test spaceship operator (<=>)
echo "5 <=> 3: ";
echo 5 <=> 3;
echo "\n";

echo "3 <=> 5: ";
echo 3 <=> 5;
echo "\n";

echo "5 <=> 5: ";
echo 5 <=> 5;
echo "\n";

// Test null coalescing
$x = null;
$y = $x ?? "default";
echo "null ?? 'default': ";
echo $y;
echo "\n";

$z = "value";
$w = $z ?? "default";
echo "'value' ?? 'default': ";
echo $w;
echo "\n";
