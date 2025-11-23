<?php
// Example: Control Flow Structures
// Demonstrates if/else, loops, and control structures

// If statement
$age = 25;
echo "Age check:\n";
if ($age >= 18) {
    echo "Adult\n";
} else {
    echo "Minor\n";
}

// For loop
echo "\nCounting with for loop:\n";
for ($i = 1; $i <= 5; $i = $i + 1) {
    echo $i;
    echo "\n";
}

// While loop
echo "\nCounting with while loop:\n";
$count = 1;
while ($count <= 3) {
    echo "Count: ";
    echo $count;
    echo "\n";
    $count = $count + 1;
}

// Nested if
$score = 85;
echo "\nGrade calculation:\n";
if ($score >= 90) {
    echo "Grade: A\n";
} else {
    if ($score >= 80) {
        echo "Grade: B\n";
    } else {
        if ($score >= 70) {
            echo "Grade: C\n";
        } else {
            echo "Grade: F\n";
        }
    }
}
