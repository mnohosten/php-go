<?php
// Test: Array operations - creation, access, and modification

// Create empty array
$arr = [];
echo "Empty array created\n";

// Add elements
$arr[0] = "first";
$arr[1] = "second";
$arr[2] = "third";
echo "Element 0: ";
echo $arr[0];
echo "\n";

// Associative array
$person = [];
$person["name"] = "John";
$person["age"] = 30;
$person["city"] = "Prague";
echo "Name: ";
echo $person["name"];
echo "\n";
echo "Age: ";
echo $person["age"];
echo "\n";

// Array literal
$colors = ["red", "green", "blue"];
echo "First color: ";
echo $colors[0];
echo "\n";
echo "Second color: ";
echo $colors[1];
echo "\n";

// Nested arrays
$matrix = [];
$matrix[0] = [1, 2, 3];
$matrix[1] = [4, 5, 6];
echo "Matrix[0][1]: ";
echo $matrix[0][1];
echo "\n";
echo "Matrix[1][2]: ";
echo $matrix[1][2];
echo "\n";

// Modify array element
$arr[0] = "modified";
echo "Modified element 0: ";
echo $arr[0];
echo "\n";

// Array with mixed keys
$mixed = [];
$mixed[0] = "zero";
$mixed["key"] = "value";
$mixed[10] = "ten";
echo "Mixed[0]: ";
echo $mixed[0];
echo "\n";
echo "Mixed[key]: ";
echo $mixed["key"];
echo "\n";
echo "Mixed[10]: ";
echo $mixed[10];
echo "\n";

// Check array element exists
if ($arr[0]) {
    echo "arr[0] is set\n";
}

// Array with explicit keys
$nums = [1, 2, 3];
$nums[3] = 4;
$nums[4] = 5;
echo "Last element: ";
echo $nums[4];
echo "\n";
