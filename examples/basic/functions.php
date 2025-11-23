<?php
// Example: Function Definitions and Calls
// Demonstrates user-defined functions (Phase 3 complete)

// Note: This example shows the syntax for when function definition is complete
// Currently, the VM supports function calls but the parser/compiler for
// user-defined functions may be in progress

echo "Function examples:\n";
echo "User-defined functions coming in Phase 3/6\n";

// Built-in function example that works now
echo "\nBuilt-in function - strlen:\n";
$text = "PHP-Go";
$length = strlen($text);
echo "Length: ";
echo $length;
echo "\n";

// Built-in function - exit/die
echo "\nBuilt-in constructs - exit/die:\n";
echo "These work for PHPT SKIPIF tests\n";
// die("This would stop execution");
