<?php
// Test: Control flow - if/else, while, for, foreach, switch

// Test if/else
$x = 10;
if ($x > 5) {
    echo "x is greater than 5\n";
} else {
    echo "x is not greater than 5\n";
}

// Test elseif
$score = 75;
if ($score >= 90) {
    echo "Grade: A\n";
} elseif ($score >= 80) {
    echo "Grade: B\n";
} elseif ($score >= 70) {
    echo "Grade: C\n";
} else {
    echo "Grade: F\n";
}

// Test while loop
$i = 0;
echo "While loop: ";
while ($i < 3) {
    echo $i;
    echo " ";
    $i++;
}
echo "\n";

// Test do-while loop
$j = 0;
echo "Do-while loop: ";
do {
    echo $j;
    echo " ";
    $j++;
} while ($j < 3);
echo "\n";

// Test for loop
echo "For loop: ";
for ($k = 0; $k < 5; $k++) {
    echo $k;
    echo " ";
}
echo "\n";

// Test foreach with array
$colors = ["red", "green", "blue"];
echo "Colors: ";
foreach ($colors as $color) {
    echo $color;
    echo " ";
}
echo "\n";

// Test foreach with associative array
$person = ["name" => "John", "age" => 30];
echo "Person: ";
foreach ($person as $key => $value) {
    echo $key;
    echo "=";
    echo $value;
    echo " ";
}
echo "\n";

// Test switch
$day = 2;
echo "Day ";
echo $day;
echo ": ";
switch ($day) {
    case 1:
        echo "Monday";
        break;
    case 2:
        echo "Tuesday";
        break;
    case 3:
        echo "Wednesday";
        break;
    default:
        echo "Unknown";
        break;
}
echo "\n";

// Test break in loop
echo "Break test: ";
for ($m = 0; $m < 10; $m++) {
    if ($m == 5) {
        break;
    }
    echo $m;
    echo " ";
}
echo "\n";

// Test continue in loop
echo "Continue test: ";
for ($n = 0; $n < 6; $n++) {
    if ($n == 3) {
        continue;
    }
    echo $n;
    echo " ";
}
echo "\n";

// Test ternary operator
$age = 20;
$status = $age >= 18 ? "adult" : "minor";
echo "Status: ";
echo $status;
echo "\n";
