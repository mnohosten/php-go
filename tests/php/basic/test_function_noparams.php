<?php
// Test: Functions without parameters

function sayHello() {
    echo "Hello from function!\n";
}

function getNumber() {
    return 42;
}

// Test function calls
sayHello();
sayHello();

$num = getNumber();
echo "Number: ";
echo $num;
echo "\n";
