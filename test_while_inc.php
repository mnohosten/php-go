<?php
$arr = ["a", "b", "c"];
$i = 0;
$total = 0;

while ($i < 3) {
    echo "Item: " . $arr[$i] . "\n";
    $total++;
    $i++;
}

echo "Total: $total\n";
