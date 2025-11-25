<?php
// Long-running PHP script for stability testing
// Runs continuously performing various operations

$startTime = microtime(true);
$iteration = 0;

// Run for a long time (will be killed by test)
while (true) {
    $iteration++;

    // Various operations
    $sum = 0;
    for ($i = 0; $i < 100; $i++) {
        $sum += $i;
    }

    // String operations
    $text = "Iteration $iteration";
    $len = strlen($text);

    // Array operations
    $arr = [];
    for ($i = 0; $i < 50; $i++) {
        $arr[] = $i;
    }
    $count = count($arr);

    // Every 1000 iterations, output status
    if ($iteration % 1000 == 0) {
        $elapsed = microtime(true) - $startTime;
        echo "Iteration $iteration after " . round($elapsed, 2) . " seconds\n";
    }
}
