<?php
if (file_exists("test_foreach.php")) {
    echo "File exists!\n";
} else {
    echo "File not found\n";
}

$content = file_get_contents("test_foreach.php");
echo "Length: " . strlen($content) . "\n";
