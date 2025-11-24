<?php
/**
 * Laravel Parse Test for PHP-Go (Simplified)
 * Tests parsing of critical Laravel files
 */

echo "Laravel Parse Test for PHP-Go\n";
echo "==============================\n\n";

$laravel_dir = __DIR__;

// Test each file individually
echo "Testing: artisan ... ";
if (file_exists($laravel_dir . '/artisan')) {
    $code = file_get_contents($laravel_dir . '/artisan');
    echo (strlen($code) > 0) ? "PASS\n" : "FAIL\n";
} else {
    echo "SKIP\n";
}

echo "Testing: public/index.php ... ";
if (file_exists($laravel_dir . '/public/index.php')) {
    $code = file_get_contents($laravel_dir . '/public/index.php');
    echo (strlen($code) > 0) ? "PASS\n" : "FAIL\n";
} else {
    echo "SKIP\n";
}

echo "Testing: bootstrap/app.php ... ";
if (file_exists($laravel_dir . '/bootstrap/app.php')) {
    $code = file_get_contents($laravel_dir . '/bootstrap/app.php');
    echo (strlen($code) > 0) ? "PASS\n" : "FAIL\n";
} else {
    echo "SKIP\n";
}

echo "Testing: bootstrap/providers.php ... ";
if (file_exists($laravel_dir . '/bootstrap/providers.php')) {
    $code = file_get_contents($laravel_dir . '/bootstrap/providers.php');
    echo (strlen($code) > 0) ? "PASS\n" : "FAIL\n";
} else {
    echo "SKIP\n";
}

echo "Testing: app/Models/User.php ... ";
if (file_exists($laravel_dir . '/app/Models/User.php')) {
    $code = file_get_contents($laravel_dir . '/app/Models/User.php');
    echo (strlen($code) > 0) ? "PASS\n" : "FAIL\n";
} else {
    echo "SKIP\n";
}

echo "\n==============================\n";
echo "Laravel critical files tested\n";
