<?php
// Test: exit() language construct
// Tests that exit() outputs message and stops execution
// The "After exit" line should NOT be printed
// Expected output: Before exit\nExiting now

echo "Before exit\n";
exit("Exiting now");
echo "After exit\n";
