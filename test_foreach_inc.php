<?php
$arr = ["a", "b"];
$total = 0;
foreach ($arr as $item) {
    echo "Item: $item, Total before: $total\n";
    $total = $total + 1;  // Use explicit addition instead of ++
    echo "Total after: $total\n";
}
echo "Final total: $total\n";
