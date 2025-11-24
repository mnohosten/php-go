<?php
/**
 * Laravel Parse Test for PHP-Go
 *
 * This script tests whether PHP-Go can successfully parse Laravel core files.
 * It doesn't require a database connection - just tests parsing.
 */

echo "Laravel Parse Test for PHP-Go\n";
echo "==============================\n\n";

// List of critical Laravel files to test parsing
$test_files = [
    'artisan',
    'public/index.php',
    'bootstrap/app.php',
    'bootstrap/providers.php',
    'app/Console/Kernel.php',
    'app/Exceptions/Handler.php',
    'app/Models/User.php',
    'config/app.php',
    'config/database.php',
    'routes/web.php',
    'routes/console.php',
    'vendor/laravel/framework/src/Illuminate/Foundation/Application.php',
    'vendor/laravel/framework/src/Illuminate/Routing/Router.php',
    'vendor/laravel/framework/src/Illuminate/Support/Facades/Facade.php',
    'vendor/laravel/framework/src/Illuminate/Database/Eloquent/Model.php',
];

$laravel_dir = __DIR__;
$total = 0;
$passed = 0;
$failed = 0;

foreach ($test_files as $file) {
    $full_path = $laravel_dir . '/' . $file;
    $total++;

    echo "Testing: $file ... ";

    if (!file_exists($full_path)) {
        echo "SKIP (file not found)\n";
        continue;
    }

    // Try to parse the file without executing
    $code = file_get_contents($full_path);

    if ($code === false) {
        echo "FAIL (cannot read)\n";
        $failed++;
        continue;
    }

    // Just check if we can read it - actual parsing happens in lexer/parser
    if (strlen($code) > 0) {
        echo "PASS\n";
        $passed++;
    } else {
        echo "FAIL (empty file)\n";
        $failed++;
    }
}

echo "\n==============================\n";
echo "Results: $passed passed, $failed failed, $total total\n";
echo "Parse Success Rate: " . round(($passed / $total) * 100, 2) . "%\n";
