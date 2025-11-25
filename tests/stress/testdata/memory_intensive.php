<?php
// Memory-intensive PHP script for memory leak detection
// Allocates and uses memory, then should release it

// Create arrays
$data = [];
for ($i = 0; $i < 1000; $i++) {
    $data[] = [
        'id' => $i,
        'name' => 'Item ' . $i,
        'value' => $i * 2,
        'data' => str_repeat('x', 100)
    ];
}

// Process arrays
$sum = 0;
foreach ($data as $item) {
    $sum += $item['value'];
}

// String operations
$text = '';
for ($i = 0; $i < 100; $i++) {
    $text .= 'Line ' . $i . "\n";
}

// Array operations
$filtered = array_filter($data, function($item) {
    return $item['id'] % 2 == 0;
});

$mapped = array_map(function($item) {
    return $item['value'] * 2;
}, $filtered);

echo "Processed " . count($data) . " items, sum: " . $sum . "\n";
echo "Filtered: " . count($filtered) . ", Mapped: " . count($mapped) . "\n";
