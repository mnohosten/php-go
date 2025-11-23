<?php
// Test: Simple class definition and instantiation

class Person {
    public $name;
    public $age;
}

$person = new Person();
$person->name = "John";
$person->age = 30;

echo "Name: ";
echo $person->name;
echo "\n";

echo "Age: ";
echo $person->age;
echo "\n";
