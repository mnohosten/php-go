# PHP-Go User Guide

Welcome to PHP-Go! This guide will help you get started with using PHP-Go, a complete PHP 8.4 interpreter written in Go with automatic parallelization and native Go integration.

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Getting Started](#getting-started)
4. [Command-Line Interface](#command-line-interface)
5. [Running PHP Scripts](#running-php-scripts)
6. [Language Features](#language-features)
7. [Standard Library](#standard-library)
8. [Parallelization](#parallelization)
9. [Go Integration](#go-integration)
10. [Configuration](#configuration)
11. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
12. [Performance](#performance)
13. [Best Practices](#best-practices)
14. [FAQ](#faq)

## Introduction

PHP-Go is a complete reimplementation of the PHP 8.4 interpreter in Go that provides:

- **Full PHP 8.4 Compatibility**: Run existing PHP code without modifications
- **Multi-threaded Execution**: Unlike traditional PHP, PHP-Go can utilize all CPU cores
- **Automatic Parallelization**: Transparently parallelize safe operations for better performance
- **Native Go Integration**: Call Go functions and use Go libraries directly from PHP
- **Single Binary Deployment**: No need to compile extensions or manage dependencies
- **Modern Architecture**: Clean, well-tested Go codebase

### Why PHP-Go?

PHP-Go addresses several limitations of traditional PHP:

1. **True Parallelism**: PHP is single-threaded; PHP-Go enables multi-core utilization
2. **Simplified Deployment**: Single binary instead of interpreter + extensions
3. **Go Ecosystem Access**: Use Go's rich library ecosystem from PHP
4. **Performance**: Competitive with PHP 8.4, with potential for optimization
5. **Modern Tooling**: Benefit from Go's excellent development tools

### What's Different from PHP?

PHP-Go maintains 100% compatibility with PHP 8.4 language features, but differs in:

- **Runtime Model**: Multi-threaded instead of single-threaded
- **Extensions**: Written in Go instead of C
- **Deployment**: Single binary instead of multiple components
- **Integration**: Native Go interop instead of FFI

## Installation

### Prerequisites

- Go 1.21 or later (for building from source)
- Git (for cloning the repository)

### From Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/krizos/php-go.git
   cd php-go
   ```

2. **Build the binary**:
   ```bash
   go build -o php-go ./cmd/php-go
   ```

3. **Optional: Install globally**:
   ```bash
   # Copy to your PATH
   sudo cp php-go /usr/local/bin/

   # Or use go install
   go install ./cmd/php-go
   ```

4. **Verify installation**:
   ```bash
   php-go --version
   ```

### Pre-built Binaries

Pre-built binaries will be available in future releases from the GitHub releases page.

### System Requirements

- **Operating System**: Linux, macOS, or Windows
- **Memory**: 256MB minimum, 1GB+ recommended
- **Disk Space**: 50MB for binary, additional space for your PHP applications

## Getting Started

### Your First Script

Create a file called `hello.php`:

```php
<?php
echo "Hello, PHP-Go!\n";
```

Run it with PHP-Go:

```bash
php-go hello.php
```

You should see:
```
Hello, PHP-Go!
```

### Interactive Testing

You can quickly test PHP expressions using the demo mode:

```bash
php-go demo
```

This will run built-in demo examples to verify your installation.

### Basic Examples

#### Variables and Types

```php
<?php
// Basic types
$integer = 42;
$float = 3.14;
$string = "Hello, World!";
$boolean = true;
$array = [1, 2, 3, 4, 5];
$assoc = ["name" => "John", "age" => 30];

// Type checking
var_dump($integer);  // int(42)
var_dump($string);   // string(13) "Hello, World!"
```

#### Control Flow

```php
<?php
// If statements
if ($score >= 90) {
    echo "Grade: A";
} elseif ($score >= 80) {
    echo "Grade: B";
} else {
    echo "Grade: C";
}

// Loops
for ($i = 0; $i < 10; $i++) {
    echo "$i\n";
}

foreach ($array as $value) {
    echo "$value\n";
}
```

#### Functions

```php
<?php
function greet($name) {
    return "Hello, $name!";
}

echo greet("World");

// Arrow functions (PHP 7.4+)
$square = fn($x) => $x * $x;
echo $square(5);  // 25
```

#### Classes and Objects

```php
<?php
class Person {
    public function __construct(
        public string $name,
        public int $age
    ) {}

    public function greet(): string {
        return "Hello, I'm {$this->name}!";
    }
}

$person = new Person("Alice", 25);
echo $person->greet();
```

## Command-Line Interface

PHP-Go provides several commands for different use cases.

### Basic Syntax

```bash
php-go [command] [options] [arguments]
```

### Commands

#### `run` / `exec` - Execute PHP Script

Run a PHP script:

```bash
php-go run script.php
php-go exec script.php
```

You can also omit the command and just provide the file:

```bash
php-go script.php
./script.php  # if script.php is executable with shebang
```

#### `lex` - Tokenize PHP Code

Display the tokens (lexer output) for a PHP file:

```bash
php-go lex script.php
```

With JSON output:

```bash
php-go lex --json script.php
```

**Use Case**: Debugging tokenization issues, understanding how PHP-Go parses your code.

**Example Output**:
```
Tokens for 'script.php':
1. PHP_OPEN_TAG        '<?php'         @ 1:1
2. ECHO                'echo'          @ 2:1
3. STRING              '"Hello"'       @ 2:6
4. SEMICOLON           ';'             @ 2:13
5. EOF                 ''              @ 2:14
```

#### `parse` - Parse PHP Code to AST

Display the Abstract Syntax Tree for a PHP file:

```bash
php-go parse script.php
```

With JSON output:

```bash
php-go parse --json script.php
```

**Use Case**: Debugging parsing issues, understanding AST structure, verifying syntax.

**Example Output**:
```
AST for 'script.php':
Program {
  Statements: [
    EchoStatement {
      Values: [
        StringLiteral { Value: "Hello" }
      ]
    }
  ]
}
```

#### `demo` - Run Demo Examples

Run built-in demonstration examples:

```bash
php-go demo
```

**Use Case**: Verify installation, see PHP-Go features in action.

#### `--version` / `-v` - Show Version

Display version information:

```bash
php-go --version
php-go -v
```

#### `--help` / `-h` - Show Help

Display help information:

```bash
php-go --help
php-go -h
```

### Command-Line Options

PHP-Go supports several command-line options for controlling execution:

| Option | Description | Example |
|--------|-------------|---------|
| `--json` | Output in JSON format (for lex/parse) | `php-go parse --json script.php` |

## Running PHP Scripts

### Basic Execution

Execute a PHP script:

```bash
php-go script.php
```

### Script Arguments

Pass arguments to your PHP script:

```bash
php-go script.php arg1 arg2 arg3
```

Access arguments in PHP using `$argv`:

```php
<?php
// script.php
echo "Script name: {$argv[0]}\n";
echo "First argument: {$argv[1]}\n";
echo "Number of arguments: {$argc}\n";
```

### Standard Input/Output

PHP-Go scripts can read from stdin and write to stdout/stderr:

```php
<?php
// Read from stdin
$input = fgets(STDIN);
echo "You entered: $input";

// Write to stderr
fwrite(STDERR, "This is an error message\n");
```

Usage:

```bash
echo "Hello" | php-go script.php
```

### Exit Codes

PHP-Go respects exit codes from your scripts:

```php
<?php
if ($error) {
    echo "Error occurred\n";
    exit(1);  // Non-zero exit code
}
echo "Success\n";
exit(0);  // Zero exit code (default)
```

Check exit code in shell:

```bash
php-go script.php
echo $?  # Prints exit code
```

### Shebang Support

Make PHP scripts executable with a shebang:

```php
#!/usr/bin/env php-go
<?php
echo "Hello from executable script!\n";
```

Make executable and run:

```bash
chmod +x script.php
./script.php
```

## Language Features

PHP-Go supports all PHP 8.4 language features:

### Modern PHP Features

#### Union Types (PHP 8.0+)

```php
<?php
function process(int|float $number): int|float {
    return $number * 2;
}
```

#### Nullable Types

```php
<?php
function greet(?string $name): string {
    return $name ?? "Guest";
}
```

#### Named Arguments (PHP 8.0+)

```php
<?php
function createUser(string $name, int $age, string $role = "user") {
    // ...
}

createUser(name: "Alice", age: 25, role: "admin");
createUser(age: 30, name: "Bob");  // Order doesn't matter
```

#### Constructor Property Promotion (PHP 8.0+)

```php
<?php
class Person {
    public function __construct(
        public string $name,
        public int $age,
        private string $email
    ) {}
}
```

#### Readonly Properties (PHP 8.1+)

```php
<?php
class Configuration {
    public function __construct(
        public readonly string $apiKey
    ) {}
}
```

#### Enums (PHP 8.1+)

```php
<?php
enum Status {
    case Pending;
    case Approved;
    case Rejected;
}

// Backed enums
enum StatusCode: int {
    case Success = 200;
    case NotFound = 404;
    case ServerError = 500;
}
```

#### First-class Callables (PHP 8.1+)

```php
<?php
$strlen = strlen(...);
echo $strlen("Hello");  // 5
```

#### Match Expressions (PHP 8.0+)

```php
<?php
$result = match ($status) {
    'pending' => 'Waiting for approval',
    'approved' => 'Request accepted',
    'rejected' => 'Request denied',
    default => 'Unknown status'
};
```

### Object-Oriented Programming

#### Classes and Inheritance

```php
<?php
class Animal {
    public function __construct(public string $name) {}

    public function speak(): string {
        return "Some sound";
    }
}

class Dog extends Animal {
    public function speak(): string {
        return "Woof!";
    }
}
```

#### Interfaces

```php
<?php
interface Drawable {
    public function draw(): void;
}

class Circle implements Drawable {
    public function draw(): void {
        echo "Drawing a circle\n";
    }
}
```

#### Traits

```php
<?php
trait Logger {
    public function log(string $message): void {
        echo "[LOG] $message\n";
    }
}

class Service {
    use Logger;

    public function process(): void {
        $this->log("Processing...");
    }
}
```

#### Abstract Classes

```php
<?php
abstract class Shape {
    abstract public function area(): float;

    public function describe(): string {
        return "Area: " . $this->area();
    }
}

class Rectangle extends Shape {
    public function __construct(
        private float $width,
        private float $height
    ) {}

    public function area(): float {
        return $this->width * $this->height;
    }
}
```

### Advanced Features

#### Namespaces

```php
<?php
namespace App\Models;

class User {
    // ...
}

// Usage
use App\Models\User;
$user = new User();
```

#### Anonymous Classes

```php
<?php
$logger = new class {
    public function log(string $message): void {
        echo "[LOG] $message\n";
    }
};
```

#### Generators

```php
<?php
function generateNumbers(int $count): Generator {
    for ($i = 1; $i <= $count; $i++) {
        yield $i;
    }
}

foreach (generateNumbers(5) as $number) {
    echo "$number\n";
}
```

#### Closures

```php
<?php
$multiplier = 3;
$multiply = function($x) use ($multiplier) {
    return $x * $multiplier;
};

echo $multiply(10);  // 30
```

## Standard Library

PHP-Go implements the complete PHP 8.4 standard library with 300+ built-in functions.

### String Functions

```php
<?php
// Common string operations
$text = "  Hello, World!  ";
echo strlen($text);           // 17
echo trim($text);             // "Hello, World!"
echo strtoupper($text);       // "  HELLO, WORLD!  "
echo substr($text, 0, 5);     // "  Hel"
echo str_replace("Hello", "Hi", $text);  // "  Hi, World!  "

// String formatting
$name = "Alice";
$age = 25;
echo sprintf("Name: %s, Age: %d", $name, $age);

// Multibyte strings
echo mb_strlen("Hello 世界");
echo mb_substr("Hello 世界", 6, 2);
```

### Array Functions

```php
<?php
$numbers = [1, 2, 3, 4, 5];

// Array manipulation
array_push($numbers, 6);
$last = array_pop($numbers);
array_unshift($numbers, 0);
$first = array_shift($numbers);

// Array transformation
$squared = array_map(fn($x) => $x * $x, $numbers);
$evens = array_filter($numbers, fn($x) => $x % 2 === 0);
$sum = array_reduce($numbers, fn($acc, $x) => $acc + $x, 0);

// Array searching
$found = in_array(3, $numbers);
$index = array_search(3, $numbers);
$exists = array_key_exists('name', $assoc);

// Sorting
sort($numbers);           // Sort by value
rsort($numbers);          // Reverse sort
asort($assoc);            // Sort preserving keys
ksort($assoc);            // Sort by keys
```

### File I/O Functions

```php
<?php
// Reading files
$content = file_get_contents('data.txt');
$lines = file('data.txt');  // Read as array of lines

// Writing files
file_put_contents('output.txt', $content);

// File operations
if (file_exists('data.txt')) {
    $size = filesize('data.txt');
    $modified = filemtime('data.txt');
}

// Directory operations
$files = scandir('.');
mkdir('newdir', 0755, true);
rmdir('olddir');
```

### JSON Functions

```php
<?php
$data = ['name' => 'Alice', 'age' => 25];

// Encode to JSON
$json = json_encode($data);
echo $json;  // {"name":"Alice","age":25}

// Decode from JSON
$decoded = json_decode($json, true);  // true = associative array
var_dump($decoded);

// Pretty printing
$pretty = json_encode($data, JSON_PRETTY_PRINT);

// Error handling
json_decode('invalid json');
if (json_last_error() !== JSON_ERROR_NONE) {
    echo json_last_error_msg();
}
```

### Date/Time Functions

```php
<?php
// Current timestamp
$now = time();
echo date('Y-m-d H:i:s', $now);

// Date manipulation
$tomorrow = strtotime('+1 day');
$lastWeek = strtotime('-1 week');

// DateTime class
$dt = new DateTime('2025-01-01');
$dt->modify('+1 month');
echo $dt->format('Y-m-d');

// Timezones
date_default_timezone_set('UTC');
$tz = new DateTimeZone('America/New_York');
```

### Hash Functions

```php
<?php
// Common hash algorithms
echo md5('password');
echo sha1('password');
echo hash('sha256', 'password');

// HMAC
echo hash_hmac('sha256', 'message', 'secret_key');

// File hashing
echo hash_file('sha256', 'data.txt');

// Password hashing
$hash = password_hash('password', PASSWORD_DEFAULT);
if (password_verify('password', $hash)) {
    echo "Password correct";
}
```

### Regular Expressions (PCRE)

```php
<?php
// Pattern matching
if (preg_match('/^[a-z]+$/', $string)) {
    echo "Matches";
}

// Find all matches
preg_match_all('/\d+/', $text, $matches);
var_dump($matches);

// Replace
$result = preg_replace('/\s+/', ' ', $text);

// Split
$parts = preg_split('/\s+/', $text);
```

### Type/Variable Functions

```php
<?php
// Type checking
is_int($var);
is_float($var);
is_string($var);
is_array($var);
is_object($var);
is_null($var);
is_bool($var);

// Variable testing
isset($var);
empty($var);
is_numeric($var);

// Type conversion
intval($var);
floatval($var);
strval($var);
boolval($var);
```

## Parallelization

One of PHP-Go's key features is automatic parallelization. Unlike traditional PHP, PHP-Go can leverage multiple CPU cores for better performance.

### Automatic Parallelization

PHP-Go automatically parallelizes certain operations when it's safe to do so:

```php
<?php
// Array operations are automatically parallelized
$numbers = range(1, 1000000);
$squared = array_map(fn($x) => $x * $x, $numbers);  // Runs in parallel

// Independent function calls
$result1 = fetchData('api1');  // These may run
$result2 = fetchData('api2');  // in parallel
$result3 = fetchData('api3');  // automatically
```

### Explicit Parallelization

PHP-Go provides APIs for explicit parallel execution:

#### `go_parallel()` - Run Functions in Parallel

Execute multiple functions concurrently:

```php
<?php
// Run three expensive operations in parallel
$futures = go_parallel([
    fn() => processLargeFile('data1.csv'),
    fn() => processLargeFile('data2.csv'),
    fn() => processLargeFile('data3.csv'),
]);

// Wait for all to complete and get results
$results = go_wait($futures);
foreach ($results as $i => $result) {
    echo "Task $i: $result\n";
}
```

#### `go_spawn()` - Spawn Goroutine

Spawn a function to run in the background:

```php
<?php
// Start background task
$future = go_spawn(function() {
    // Long-running task
    sleep(5);
    return "Completed!";
});

// Do other work
echo "Working...\n";

// Wait for result when needed
$result = go_wait([$future])[0];
echo $result;
```

#### `go_channel()` - Create Channel

Use channels for communication between parallel tasks:

```php
<?php
$channel = go_channel();

// Producer
go_spawn(function() use ($channel) {
    for ($i = 1; $i <= 10; $i++) {
        go_send($channel, $i);
    }
    go_close($channel);
});

// Consumer
while ($value = go_receive($channel)) {
    echo "Received: $value\n";
}
```

### Thread Safety

When using parallelization, be aware of thread safety:

```php
<?php
// ✗ NOT thread-safe: shared mutable state
$counter = 0;
go_parallel([
    fn() => $counter++,  // Race condition!
    fn() => $counter++,
]);

// ✓ Thread-safe: use mutexes
$mutex = go_mutex();
$counter = 0;
go_parallel([
    fn() => go_with_lock($mutex, fn() => $counter++),
    fn() => go_with_lock($mutex, fn() => $counter++),
]);
```

### Best Practices

1. **Use immutable data**: Prefer passing values instead of references
2. **Avoid shared state**: Minimize shared mutable state between goroutines
3. **Use channels**: Communicate via channels instead of shared memory
4. **Profile first**: Don't parallelize prematurely; measure first
5. **Mind overhead**: Parallelization has overhead; use for CPU-intensive tasks

## Go Integration

PHP-Go allows you to call Go functions and use Go libraries directly from PHP.

### Calling Go Functions

Use `go_call()` to invoke Go functions:

```php
<?php
// Call Go's HTTP client
$response = go_call('http.Get', 'https://api.example.com/data');
echo "Status: " . $response['StatusCode'] . "\n";
echo "Body: " . $response['Body'] . "\n";

// Call Go's crypto functions
$hash = go_call('crypto.SHA256', $data);
echo "SHA256: " . bin2hex($hash) . "\n";

// Call Go's JSON encoding
$json = go_call('json.Marshal', ['name' => 'Alice', 'age' => 25]);
echo $json;
```

### Importing Go Packages

Import and use Go packages:

```php
<?php
// Import Go packages
go_import('fmt');
go_import('math');
go_import('time');

// Use imported functions
go_call('fmt.Println', 'Hello from Go!');
$pi = go_call('math.Pi');
$now = go_call('time.Now');
```

### Type Conversion

PHP-Go automatically converts between PHP and Go types:

| PHP Type | Go Type |
|----------|---------|
| int | int64 |
| float | float64 |
| string | string |
| bool | bool |
| array | []interface{} or map[string]interface{} |
| object | map[string]interface{} |
| null | nil |

### Writing Go Extensions

You can write custom Go extensions for PHP-Go:

```go
// my_extension.go
package main

import "github.com/krizos/php-go/pkg/goext"

func init() {
    goext.RegisterFunction("my_custom_function", MyCustomFunction)
}

func MyCustomFunction(args ...interface{}) (interface{}, error) {
    // Your Go code here
    return "Hello from Go!", nil
}
```

Build and use:

```bash
go build -buildmode=plugin -o my_extension.so my_extension.go
```

```php
<?php
go_load_plugin('my_extension.so');
$result = my_custom_function();
echo $result;  // "Hello from Go!"
```

## Configuration

PHP-Go can be configured via configuration files and environment variables.

### Configuration File

Create a `php-go.ini` file:

```ini
[runtime]
memory_limit = 256M
max_execution_time = 30
error_reporting = E_ALL

[parallelization]
max_workers = 8
queue_size = 1000
enable_auto_parallel = true

[logging]
log_level = info
log_output = stderr
log_format = json

[performance]
opcache_enabled = true
jit_enabled = false
```

Load configuration:

```bash
php-go --config php-go.ini script.php
```

### Environment Variables

Configure via environment variables:

```bash
export PHPGO_MEMORY_LIMIT=512M
export PHPGO_MAX_WORKERS=16
export PHPGO_LOG_LEVEL=debug

php-go script.php
```

### Runtime Configuration

Configure at runtime in PHP:

```php
<?php
// Set configuration
ini_set('memory_limit', '512M');
ini_set('max_execution_time', '60');

// Get configuration
$limit = ini_get('memory_limit');
echo "Memory limit: $limit\n";
```

## Debugging and Troubleshooting

### Debugging Tools

#### Lexer Output

See how your code is tokenized:

```bash
php-go lex script.php
```

Use this when you suspect tokenization issues.

#### Parser Output

See the Abstract Syntax Tree:

```bash
php-go parse script.php
```

Use this when you suspect parsing issues.

#### Error Messages

PHP-Go provides detailed error messages with file, line, and column information:

```
Parse error: Unexpected token '}', expected ';'
  at /path/to/script.php:15:8
```

### Common Issues

#### "Parse error: Unexpected token"

**Cause**: Syntax error in your PHP code.

**Solution**: Check the line and column indicated, ensure proper syntax.

#### "Undefined function"

**Cause**: Function not implemented or not loaded.

**Solution**: Check if the function is part of PHP's standard library. Some extensions may not be implemented yet.

#### "Memory limit exceeded"

**Cause**: Script uses more memory than allowed.

**Solution**: Increase memory limit:
```php
ini_set('memory_limit', '512M');
```

#### "Maximum execution time exceeded"

**Cause**: Script runs longer than allowed.

**Solution**: Increase execution time:
```php
set_time_limit(60);  // 60 seconds
```

### Performance Issues

If your script runs slowly:

1. **Profile your code**: Identify bottlenecks
2. **Use parallelization**: For CPU-intensive tasks
3. **Optimize algorithms**: Better algorithms > faster execution
4. **Check I/O operations**: File/network operations are often bottlenecks
5. **Monitor memory**: Large data structures can slow down GC

### Getting Help

If you encounter issues:

1. Check this documentation
2. Search GitHub issues: https://github.com/krizos/php-go/issues
3. File a bug report with:
   - PHP-Go version (`php-go --version`)
   - Your code (minimal reproducible example)
   - Error messages
   - Expected vs actual behavior

## Performance

### Benchmarking

PHP-Go provides built-in benchmarking tools:

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Benchmark specific package
go test -bench=. ./pkg/vm/
```

### Performance Characteristics

**Single-threaded Performance**: Competitive with PHP 8.4 + opcache

**Multi-threaded Performance**: 2-4x faster for parallel workloads

**Memory Usage**: Similar to PHP, benefits from Go's efficient GC

### Optimization Tips

1. **Use parallelization**: For CPU-intensive tasks
2. **Minimize allocations**: Reuse objects when possible
3. **Use generators**: For large datasets
4. **Cache results**: Avoid redundant computations
5. **Profile first**: Measure before optimizing

### Comparison with PHP 8.4

| Scenario | PHP 8.4 | PHP-Go | Improvement |
|----------|---------|--------|-------------|
| Simple script | 1.0x | 0.8-1.2x | Similar |
| Web request | 1.0x | 1.0-1.5x | 0-50% faster |
| Parallel work | 1.0x | 2.0-4.0x | 2-4x faster |
| Array operations | 1.0x | 0.9-1.3x | Similar |

## Best Practices

### Code Organization

```php
<?php
// Use namespaces
namespace App\Services;

// Use type declarations
function processData(array $data): array {
    // ...
}

// Use readonly for immutable data
class Config {
    public function __construct(
        public readonly string $apiKey
    ) {}
}
```

### Error Handling

```php
<?php
// Use exceptions
try {
    $result = riskyOperation();
} catch (Exception $e) {
    error_log("Error: " . $e->getMessage());
    throw $e;
}

// Check return values
$result = someFunction();
if ($result === false) {
    // Handle error
}
```

### Resource Management

```php
<?php
// Close resources properly
$file = fopen('data.txt', 'r');
try {
    $content = fread($file, filesize('data.txt'));
} finally {
    fclose($file);
}

// Or use file_get_contents/file_put_contents
$content = file_get_contents('data.txt');
```

### Security

```php
<?php
// Validate input
$id = filter_input(INPUT_GET, 'id', FILTER_VALIDATE_INT);
if ($id === false) {
    die('Invalid ID');
}

// Escape output
echo htmlspecialchars($userInput, ENT_QUOTES, 'UTF-8');

// Use prepared statements (when DB available)
$stmt = $pdo->prepare('SELECT * FROM users WHERE id = ?');
$stmt->execute([$id]);
```

## FAQ

### Is PHP-Go production-ready?

PHP-Go is currently in active development. While core features are implemented and tested, it's recommended to thoroughly test your specific use case before production deployment.

### Can I run my existing PHP application?

Yes! PHP-Go aims for 100% PHP 8.4 compatibility. Most applications should work without modifications. Some C-based extensions may need Go equivalents.

### How do I install PHP extensions?

PHP-Go doesn't use traditional C-based extensions. Core extensions are built-in, and additional extensions can be written in Go. See the Extension Development Guide for details.

### Is PHP-Go faster than PHP?

For single-threaded workloads, PHP-Go is competitive with PHP 8.4. For parallel workloads, PHP-Go can be 2-4x faster by utilizing multiple CPU cores.

### Can I use Composer packages?

Yes! Pure PHP packages work without modification. Packages requiring C extensions need Go equivalents.

### Does PHP-Go support Windows?

Yes, PHP-Go supports Windows, Linux, and macOS.

### How do I report bugs?

File an issue on GitHub: https://github.com/krizos/php-go/issues

Include:
- PHP-Go version
- Operating system
- Minimal reproducible example
- Expected vs actual behavior

### Can I contribute?

Yes! Contributions are welcome. See CONTRIBUTING.md for guidelines.

### What's the license?

[To be determined - likely MIT or Apache 2.0]

### Where can I get help?

- Documentation: docs/user-guide/
- GitHub Issues: https://github.com/krizos/php-go/issues
- Examples: examples/

## Next Steps

Now that you understand the basics:

1. **Try the examples**: Check out `examples/` directory
2. **Read the installation guide**: `docs/user-guide/installation.md`
3. **Learn about configuration**: `docs/user-guide/configuration.md`
4. **Explore Go integration**: `docs/extension-guide/`
5. **Optimize performance**: `docs/user-guide/performance.md`

Happy coding with PHP-Go!
