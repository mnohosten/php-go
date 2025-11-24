<?php
$test_files = [
    'artisan',
    'public/index.php',
    'bootstrap/app.php',
];

$total = 0;
foreach ($test_files as $file) {
    $total++;
    echo "Testing: $file\n";
}

echo "Total: $total\n";
