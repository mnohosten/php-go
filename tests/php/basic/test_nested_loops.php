<?php
// Test: Nested loops

$result = 0;
$i = 0;
while ($i < 3) {
    $j = 0;
    while ($j < 2) {
        $result = $result + 1;
        $j = $j + 1;
    }
    $i = $i + 1;
}

echo "Result: ";
echo $result;
echo "\n";

// Expected: 3 * 2 = 6
