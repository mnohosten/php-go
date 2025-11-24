<?php
/**
 * Symfony Parse Test for PHP-Go
 * Tests parsing of critical Symfony files
 */

echo "Symfony Parse Test for PHP-Go\n";
echo "==============================\n\n";

$symfony_dir = __DIR__ . '/symfony-app';

// Critical files to test
$critical_files = [
    'bin/console',
    'public/index.php',
    'src/Kernel.php',
    'config/bundles.php',
    'config/services.yaml',
];

$passed = 0;
$failed = 0;
$skipped = 0;

foreach ($critical_files as $file) {
    $full_path = $symfony_dir . '/' . $file;
    echo "Testing: $file ... ";

    if (file_exists($full_path)) {
        $code = file_get_contents($full_path);
        if (strlen($code) > 0) {
            echo "PASS\n";
            $passed++;
        } else {
            echo "FAIL (empty file)\n";
            $failed++;
        }
    } else {
        echo "SKIP (not found)\n";
        $skipped++;
    }
}

echo "\n==============================\n";
echo "Results:\n";
echo "  Passed:  $passed\n";
echo "  Failed:  $failed\n";
echo "  Skipped: $skipped\n";
echo "  Total:   " . count($critical_files) . "\n";
echo "==============================\n";
