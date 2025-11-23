<?php
// Test: String concatenation with variable
// KNOWN ISSUE: Currently outputs only "World" instead of "Hello World"
// The concat operator needs investigation
// Expected output: Hello World

$name = "World";
echo "Hello " . $name;
