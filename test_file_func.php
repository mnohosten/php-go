<?php
$path = "/Users/krizos/code/mnohosten/php-go/test_foreach.php";
echo "Testing file: $path\n";

if (file_exists($path)) {
    echo "File exists\n";
    $content = file_get_contents($path);
    echo "Got content\n";
    $len = strlen($content);
    echo "Length: $len\n";
} else {
    echo "File not found\n";
}
