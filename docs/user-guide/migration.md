# PHP-Go Migration Guide

This guide helps you migrate existing PHP applications to PHP-Go, covering compatibility considerations, common migration patterns, and troubleshooting strategies.

## Table of Contents

1. [Introduction](#introduction)
2. [Pre-Migration Assessment](#pre-migration-assessment)
3. [Compatibility Overview](#compatibility-overview)
4. [Migration Strategy](#migration-strategy)
5. [Code Migration](#code-migration)
6. [Extension Migration](#extension-migration)
7. [Configuration Migration](#configuration-migration)
8. [Testing and Validation](#testing-and-validation)
9. [Performance Optimization](#performance-optimization)
10. [Common Issues and Solutions](#common-issues-and-solutions)
11. [Rollback Strategy](#rollback-strategy)
12. [Migration Checklist](#migration-checklist)

## Introduction

PHP-Go is designed to be a drop-in replacement for PHP 8.4, providing full language compatibility while offering enhanced performance through multi-threaded execution and native Go integration.

### Who Should Migrate?

PHP-Go is ideal for:

- **Performance-Critical Applications**: Applications that need better CPU utilization
- **Microservices**: Modern architectures benefiting from single-binary deployment
- **Go Ecosystem Integration**: Projects that want to leverage Go libraries
- **Concurrent Workloads**: Applications that can benefit from parallelization

### When to Migrate

Consider migrating when:

- Your application is on PHP 8.2+ (easier migration path)
- You need better performance without scaling infrastructure
- You want to integrate with Go libraries
- You're starting a greenfield project with modern PHP
- Your team is comfortable with both PHP and Go

### When NOT to Migrate

Avoid migrating if:

- Your application heavily relies on PHP extensions not yet available in PHP-Go
- You need 100% compatibility with legacy PHP 5.x/7.x code
- Your team lacks Go expertise for troubleshooting
- You require specific PHP SAPI modules (mod_php, php-fpm specific features)

## Pre-Migration Assessment

### Step 1: Analyze Your Codebase

Before migrating, understand your application's requirements:

#### Check PHP Version Compatibility

```bash
# Check your current PHP version
php -v

# Scan for PHP version-specific features
grep -r "declare(strict_types" .
grep -r "readonly" .
grep -r "enum " .
```

**Action Items:**
- Applications on PHP 8.2+ are easiest to migrate
- PHP 7.4-8.1 applications may need minor adjustments
- PHP 5.x/7.0-7.3 applications will require code modernization first

#### Inventory Extensions

```bash
# List installed PHP extensions
php -m

# Check extension usage in code
grep -r "extension_loaded" .
grep -r "function_exists" .
```

**Currently Supported Extensions (PHP-Go):**
- Core language features (arrays, strings, variables)
- JSON (full support)
- Hash (12 core functions: md5, sha1, sha256, etc.)
- Date/Time (partial support)
- SPL (partial support)
- PCRE - Regular expressions (partial support)

**Extensions Not Yet Available:**
- Database: mysqli, PDO, pgsql
- XML: dom, simplexml, xml, xmlreader, xmlwriter
- Compression: zlib, zip, bzip2
- Image: gd, imagick
- Network: curl, sockets
- Caching: opcache, apcu
- Multi-byte: mbstring
- Cryptography: openssl (beyond basic hash functions)

**Action Items:**
- List all extensions your application uses
- Check compatibility with PHP-Go supported extensions
- Plan workarounds or wait for extension implementation
- Consider using Go libraries as alternatives for missing extensions

#### Analyze Standard Library Usage

```bash
# Find commonly used functions
grep -roh '\b[a-z_]\+(' . | sort | uniq -c | sort -rn | head -50

# Check for specific function categories
grep -r "mysqli_" .          # Database
grep -r "curl_" .            # HTTP client
grep -r "file_get_contents" . # File I/O
grep -r "preg_" .            # Regex
grep -r "mb_" .              # Multi-byte strings
```

**Action Items:**
- Identify critical functions used in your application
- Verify these functions are implemented in PHP-Go (see [API Reference](../api/stdlib.md))
- Prepare fallback implementations for unsupported functions

#### Check Framework Compatibility

**WordPress:**
- Parse success rate: ~61% (after inline HTML support)
- Most core files parse correctly
- Missing: mysqli/PDO extensions (critical blocker)
- Estimated completion: 2-3 months for full support

**Laravel:**
- Parse success rate: ~43%
- Requires modern PHP 8.2+ features
- Missing: Advanced type declarations, attributes
- Estimated completion: 3-4 months for full support

**Symfony:**
- Parse success rate: ~50%
- Heavy attribute usage
- Missing: Generators, first-class callables
- Estimated completion: 3-4 months for full support

**Action Items:**
- If using a major framework, check compatibility status above
- Consider waiting for framework-specific support
- Test critical framework features in PHP-Go environment

### Step 2: Set Up Testing Environment

Create a parallel environment for testing:

```bash
# Create test directory
mkdir php-go-test
cd php-go-test

# Copy your application
cp -r /path/to/your/app .

# Install PHP-Go
wget https://github.com/krizos/php-go/releases/latest/download/php-go-linux-amd64.tar.gz
tar -xzf php-go-linux-amd64.tar.gz
chmod +x php-go
```

### Step 3: Create Migration Plan

Document your migration approach:

1. **Phase 1: Assessment** (1-2 weeks)
   - Inventory dependencies
   - Test parse compatibility
   - Identify blockers

2. **Phase 2: Environment Setup** (1 week)
   - Set up PHP-Go testing environment
   - Configure CI/CD for dual testing
   - Train team on PHP-Go basics

3. **Phase 3: Incremental Migration** (2-8 weeks)
   - Migrate non-critical components first
   - Update code for compatibility
   - Run parallel testing

4. **Phase 4: Production Rollout** (2-4 weeks)
   - Gradual traffic shift
   - Monitor performance and errors
   - Optimize based on metrics

5. **Phase 5: Cleanup** (1-2 weeks)
   - Remove compatibility shims
   - Optimize for PHP-Go features
   - Document learnings

## Compatibility Overview

### Language Features

#### Fully Supported

PHP-Go supports all PHP 8.4 language features:

- **Variables and Types**: All PHP types (int, float, string, array, object, resource, null, bool)
- **Operators**: Arithmetic, comparison, logical, bitwise, assignment, ternary, null coalescing
- **Control Structures**: if/else, switch, while, do-while, for, foreach, break, continue, goto
- **Functions**: Function declaration, parameters, return values, variadic functions, type declarations
- **Classes**: Classes, inheritance, interfaces, traits, abstract classes, final classes
- **OOP Features**: Properties, methods, constructors, destructors, static members, constants
- **Visibility**: public, protected, private
- **Magic Methods**: __construct, __destruct, __get, __set, __call, __callStatic, __toString, etc.
- **Enums**: Pure enums and backed enums (PHP 8.1+)
- **Readonly**: Readonly properties (PHP 8.1+) and readonly classes (PHP 8.2+)
- **Type System**: Union types, nullable types, mixed type, void type, never type
- **Attributes**: PHP 8.0+ attributes (parsing complete, runtime access in progress)
- **Named Arguments**: PHP 8.0+ named arguments
- **Match Expressions**: PHP 8.0+ match expressions
- **Null Safe Operator**: ?-> operator
- **Array Spread**: Array unpacking with ... operator
- **Arrow Functions**: fn() => short closures
- **First-class Callables**: $fn(...) syntax (in progress)

#### Partially Supported

- **Generators**: yield/yield from syntax (parser complete, runtime in progress)
- **Closures**: Anonymous functions (basic support, variable capture in progress)
- **Reflection**: Basic reflection API (full API in Phase 9)

#### Differences from Standard PHP

| Feature | PHP 8.4 | PHP-Go | Notes |
|---------|---------|--------|-------|
| Threading Model | Single-threaded | Multi-threaded | PHP-Go can use multiple cores |
| Memory Management | Reference counting | Garbage collection | Less predictable memory usage |
| Extension API | C language | Go language | Extensions written in Go, not C |
| SAPI | Multiple (CLI, FPM, etc.) | CLI + embedded | No mod_php or FPM |
| Request Handling | One per process | Concurrent | Multiple requests per process |
| Resource Cleanup | Deterministic | Non-deterministic | GC runs when needed, not immediately |

### Standard Library

#### Core Functions (Fully Supported)

```php
// Variable functions
var_dump(), print_r(), isset(), empty(), unset()

// Array functions (partial)
count(), array_keys(), array_values(), array_merge(),
array_push(), array_pop(), array_shift(), array_unshift()

// String functions (partial)
strlen(), substr(), str_repeat(), strpos(), strtolower(), strtoupper()

// Type functions
is_int(), is_float(), is_string(), is_array(), is_object(),
is_null(), is_bool(), gettype(), settype()

// Output functions
echo, print, printf(), sprintf()
```

#### Extension Functions

**JSON Extension (Full Support):**
```php
json_encode(), json_decode(), json_last_error(), json_last_error_msg()
```

**Hash Extension (Partial Support):**
```php
hash(), hash_file(), hash_hmac(), hash_hmac_file(),
hash_init(), hash_update(), hash_final(),
hash_equals(), hash_algos()

// Supported algorithms:
// md5, sha1, sha224, sha256, sha384, sha512,
// sha3-224, sha3-256, sha3-384, sha3-512
```

**Date/Time Extension (Partial Support):**
```php
time(), date(), strtotime(), mktime()
// DateTime class support in progress
```

**SPL Extension (Partial Support):**
```php
// Iterators, data structures in progress
```

#### Missing Functions

See the [Standard Library Status](#standard-library-status) section for complete lists.

Common missing functions that may affect migration:

**Database:**
- mysqli_*, PDO::*, pg_*

**HTTP:**
- curl_*, file_get_contents() with HTTP context

**File System:**
- Many advanced file operations

**Multi-byte Strings:**
- mb_* functions

**XML:**
- dom_*, simplexml_*, xml_*

**Action Items:**
- Audit your code for missing functions
- Implement wrapper functions or polyfills where possible
- Use Go integration to access Go standard library equivalents

## Migration Strategy

### Strategy 1: Big Bang Migration (Not Recommended)

Replace entire PHP installation with PHP-Go at once.

**Pros:**
- Quick migration
- Single deployment
- Clean cutover

**Cons:**
- High risk
- Difficult rollback
- All-or-nothing approach

**Use When:**
- Small, simple applications
- Complete test coverage
- Easy rollback mechanism

### Strategy 2: Incremental Migration (Recommended)

Gradually migrate components one at a time.

**Pros:**
- Lower risk
- Easier troubleshooting
- Gradual validation

**Cons:**
- Longer migration period
- Dual maintenance
- Integration complexity

**Use When:**
- Large applications
- Critical production systems
- Limited PHP-Go experience

### Strategy 3: Parallel Run

Run both PHP and PHP-Go side-by-side, comparing results.

**Pros:**
- Validates correctness
- Confidence building
- Performance comparison

**Cons:**
- Double infrastructure
- Higher complexity
- Doubled costs

**Use When:**
- Critical financial systems
- Need 100% compatibility validation
- Risk-averse organization

### Strategy 4: Greenfield Projects

Use PHP-Go for new projects, keep existing on PHP.

**Pros:**
- Zero migration risk
- Learn PHP-Go gradually
- Gain experience before migration

**Cons:**
- Longer time to production use
- Dual technology stack
- Team split

**Use When:**
- Starting new microservices
- Building new features
- Experimental adoption

## Code Migration

### Step 1: Parse Compatibility Test

Test if your code parses correctly:

```bash
# Test individual file
php-go parse yourfile.php

# Test entire codebase
find . -name "*.php" -exec php-go parse {} \;

# Automated parse test script
cat > test-parse.php << 'EOF'
<?php
$dir = $argv[1] ?? '.';
$files = new RecursiveIteratorIterator(
    new RecursiveDirectoryIterator($dir)
);

$total = 0;
$success = 0;
$failed = [];

foreach ($files as $file) {
    if ($file->getExtension() !== 'php') continue;
    $total++;

    exec("php-go parse " . escapeshellarg($file->getPathname()), $output, $code);
    if ($code === 0) {
        $success++;
    } else {
        $failed[] = $file->getPathname();
    }
}

echo "Parsed: $success / $total (" . round($success/$total*100, 2) . "%)\n";
if ($failed) {
    echo "Failed files:\n";
    foreach ($failed as $f) echo "  - $f\n";
}
EOF

php-go test-parse.php /path/to/your/code
```

### Step 2: Fix Parse Errors

Common parse errors and fixes:

#### Alternative Syntax (Already Fixed)

If using older code with alternative control structure syntax:

```php
// This now works in PHP-Go
<?php if ($condition): ?>
    <html content>
<?php endif; ?>
```

#### Inline HTML Support (Already Fixed)

PHP-Go now supports mixing PHP and HTML:

```php
<?php
if ($user) {
    ?>
    <div>Welcome <?php echo $user; ?></div>
    <?php
}
?>
```

#### Reserved Keywords as Method Names

PHP-Go may be stricter about reserved keywords:

```php
// Problematic in PHP-Go
class Container {
    public function do() { }  // 'do' is reserved
}

// Solution: Rename method
class Container {
    public function execute() { }
}
```

### Step 3: Update Extensions

Replace PHP extensions with PHP-Go equivalents or Go libraries.

#### Database Extensions

**Standard PHP:**
```php
$mysqli = new mysqli('localhost', 'user', 'pass', 'db');
$result = $mysqli->query('SELECT * FROM users');
```

**Migration Options:**

Option 1: Wait for mysqli/PDO implementation (recommended)

Option 2: Use Go database libraries via PHP-Go integration:

```php
// Using PHP-Go's Go integration (future feature)
$db = go\sql\Open('mysql', 'user:pass@tcp(localhost:3306)/db');
$rows = $db->Query('SELECT * FROM users');
```

Option 3: Use REST API wrapper:

```php
// Create a Go microservice for database access
function query($sql) {
    $ch = curl_init('http://localhost:8080/query');
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode(['sql' => $sql]));
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    $result = curl_exec($ch);
    return json_decode($result, true);
}
```

#### HTTP Client

**Standard PHP:**
```php
$ch = curl_init('https://api.example.com/data');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$result = curl_exec($ch);
```

**PHP-Go (using file_get_contents when available):**
```php
// For simple GET requests
$context = stream_context_create([
    'http' => [
        'method' => 'GET',
        'header' => 'Accept: application/json'
    ]
]);
$result = file_get_contents('https://api.example.com/data', false, $context);
```

#### Image Processing

**Standard PHP:**
```php
$img = imagecreatefrompng('input.png');
imagefilter($img, IMG_FILTER_GRAYSCALE);
imagepng($img, 'output.png');
```

**Migration Options:**

Option 1: Use Go image libraries (future feature):
```php
// Via PHP-Go Go integration
$img = go\image\LoadPNG('input.png');
$gray = go\image\Grayscale($img);
go\image\SavePNG($gray, 'output.png');
```

Option 2: External service:
```bash
# Use ImageMagick CLI
exec('convert input.png -colorspace Gray output.png');
```

### Step 4: Update Configuration

Migrate php.ini settings to PHP-Go configuration:

```bash
# Create php-go.ini or php-go.yaml
cat > php-go.yaml << 'EOF'
# PHP-Go Configuration

# Error handling
error_reporting: E_ALL
display_errors: true
log_errors: true
error_log: /var/log/php-go/error.log

# Resource limits
memory_limit: 256M
max_execution_time: 30
max_input_time: 60

# Output
output_buffering: 4096
implicit_flush: false

# File uploads
upload_max_filesize: 64M
post_max_size: 64M
max_file_uploads: 20

# Date/Time
date.timezone: UTC

# Extensions
extensions:
  - json
  - hash

# Parallelization (PHP-Go specific)
parallel:
  enabled: true
  max_workers: 8
  auto_parallelize_arrays: true

# Go integration (PHP-Go specific)
go:
  enabled: true
  plugin_dir: /usr/local/lib/php-go/plugins
EOF
```

See [Configuration Guide](configuration.md) for complete options.

### Step 5: Handle Behavioral Differences

#### Floating Point Precision

PHP-Go uses Go's float64, which may have slight differences:

```php
// Results may differ in edge cases
$a = 0.1 + 0.2;
// PHP: 0.30000000000000004
// PHP-Go: Same, but internal representation differs
```

**Solution:** Use proper floating point comparisons:

```php
function floatEquals($a, $b, $epsilon = 0.00001) {
    return abs($a - $b) < $epsilon;
}
```

#### Resource Cleanup

PHP has deterministic resource cleanup; PHP-Go uses garbage collection:

```php
// PHP: File closed when $handle goes out of scope or on unset()
$handle = fopen('file.txt', 'r');
// ... use handle ...
// Automatically closed here

// PHP-Go: File closed when GC runs (non-deterministic)
$handle = fopen('file.txt', 'r');
// ... use handle ...
fclose($handle);  // Explicit close recommended
```

**Solution:** Always explicitly close resources:

```php
// Pattern 1: Explicit close
$handle = fopen('file.txt', 'r');
try {
    // Use handle
} finally {
    fclose($handle);
}

// Pattern 2: Use temporary variables
function processFile($path) {
    $handle = fopen($path, 'r');
    $data = fread($handle, filesize($path));
    fclose($handle);
    return $data;
}
```

#### Error Handling

Error handling is compatible but may have subtle differences:

```php
// PHP: set_error_handler affects all errors in current request
set_error_handler(function($errno, $errstr) {
    // Handle error
});

// PHP-Go: Same behavior in single-threaded mode
// In parallel mode: Each goroutine has its own error handler context
```

**Solution:** Use try/catch for critical code:

```php
try {
    // Critical operation
} catch (Exception $e) {
    // Handle exception
}
```

## Extension Migration

### Identify Extension Usage

```bash
# Find extension-specific functions
grep -roh '\b[a-z_]\+_[a-z_]\+(' . | grep -E '^(mysqli|pdo|curl|gd|xml)_' | sort -u
```

### Extension Migration Strategies

#### 1. Use PHP-Go Native Extensions (Preferred)

If the extension is available in PHP-Go:

```php
// Works in both PHP and PHP-Go
$json = json_encode($data);
$hash = hash('sha256', $data);
```

#### 2. Polyfill Missing Functions

Create PHP implementations for missing functions:

```php
// polyfills.php
if (!function_exists('mb_strlen')) {
    function mb_strlen($str, $encoding = 'UTF-8') {
        // Basic implementation
        return strlen(utf8_decode($str));
    }
}

if (!function_exists('array_is_list')) {
    function array_is_list(array $array) {
        $i = 0;
        foreach ($array as $k => $v) {
            if ($k !== $i++) return false;
        }
        return true;
    }
}
```

#### 3. Create Go Extensions

Write custom extensions in Go:

```go
// myextension/extension.go
package main

import "github.com/krizos/php-go/pkg/goext"

func init() {
    goext.Register(&goext.Extension{
        Name: "myextension",
        Version: "1.0.0",
        Functions: map[string]goext.Function{
            "my_function": {
                Name: "my_function",
                Handler: myFunction,
            },
        },
    })
}

func myFunction(args []interface{}) (interface{}, error) {
    // Your Go implementation
    return result, nil
}
```

Build and load:

```bash
go build -buildmode=plugin -o myextension.so
php-go -d extension=/path/to/myextension.so script.php
```

See [Extension Development Guide](extension-development.md) for details.

#### 4. Use External Services

For complex extensions, use microservices:

```php
// Redis example: Use REST wrapper instead of phpredis
class RedisClient {
    private $baseUrl = 'http://localhost:6379';

    public function get($key) {
        return file_get_contents("{$this->baseUrl}/get/$key");
    }

    public function set($key, $value, $ttl = 0) {
        $context = stream_context_create([
            'http' => [
                'method' => 'POST',
                'content' => json_encode(['value' => $value, 'ttl' => $ttl])
            ]
        ]);
        return file_get_contents("{$this->baseUrl}/set/$key", false, $context);
    }
}
```

## Configuration Migration

### Convert php.ini to php-go Configuration

PHP-Go supports two configuration formats: INI (PHP-compatible) and YAML (extended features).

#### INI Format (Compatible)

```ini
; php-go.ini - PHP-compatible configuration

[PHP]
error_reporting = E_ALL
display_errors = On
log_errors = On
error_log = /var/log/php-go/error.log

memory_limit = 256M
max_execution_time = 30
post_max_size = 64M
upload_max_filesize = 64M

date.timezone = UTC

[extensions]
extension=json
extension=hash
```

#### YAML Format (Extended)

```yaml
# php-go.yaml - Extended configuration with PHP-Go features

# PHP-compatible settings
php:
  error_reporting: E_ALL
  display_errors: true
  log_errors: true
  error_log: /var/log/php-go/error.log

  memory_limit: 256M
  max_execution_time: 30
  post_max_size: 64M
  upload_max_filesize: 64M

  date:
    timezone: UTC

# Extensions
extensions:
  enabled:
    - json
    - hash

  paths:
    - /usr/local/lib/php-go/extensions
    - /opt/php-go/extensions

# PHP-Go specific features
parallelization:
  enabled: true
  max_workers: 8
  array_operations:
    enabled: true
    min_size: 1000

logging:
  level: info
  format: json
  output: /var/log/php-go/app.log

metrics:
  enabled: true
  port: 9090
  path: /metrics

health:
  enabled: true
  port: 8080
  path: /health

go_integration:
  enabled: true
  plugin_directory: /usr/local/lib/php-go/plugins
  import_path: /usr/local/lib/php-go/imports
```

### Configuration Mapping

| PHP Setting | PHP-Go Equivalent | Notes |
|-------------|-------------------|-------|
| memory_limit | memory_limit | Same syntax (128M, 1G, etc.) |
| max_execution_time | max_execution_time | Seconds |
| error_reporting | error_reporting | Same constants |
| display_errors | display_errors | On/Off or true/false |
| log_errors | log_errors | On/Off or true/false |
| error_log | error_log | File path |
| date.timezone | date.timezone | IANA timezone |
| post_max_size | post_max_size | Size limit |
| upload_max_filesize | upload_max_filesize | Size limit |
| extension=name | extensions: [name] | YAML array format |

### Environment-Specific Configuration

Use environment variables for deployment:

```bash
# Production
export PHPGO_CONFIG=/etc/php-go/production.yaml
export PHPGO_ENV=production
php-go /var/www/app.php

# Development
export PHPGO_CONFIG=/etc/php-go/development.yaml
export PHPGO_ENV=development
php-go /var/www/app.php
```

Configuration precedence:
1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

## Testing and Validation

### Unit Testing

Run existing PHPUnit tests:

```bash
# If PHPUnit works with PHP-Go
php-go vendor/bin/phpunit

# Generate coverage report
php-go -d extension=xdebug vendor/bin/phpunit --coverage-html coverage/
```

### Integration Testing

Test critical workflows:

```php
// tests/integration/DatabaseTest.php
class DatabaseTest extends TestCase {
    public function testUserCreation() {
        $user = User::create(['name' => 'Test', 'email' => 'test@example.com']);
        $this->assertNotNull($user->id);
        $this->assertEquals('Test', $user->name);
    }

    public function testUserQuery() {
        $users = User::where('active', true)->get();
        $this->assertGreaterThan(0, $users->count());
    }
}
```

### Functional Testing

Test user-facing features:

```bash
# Selenium, Behat, or custom scripts
./vendor/bin/behat
```

### Performance Testing

Compare performance between PHP and PHP-Go:

```php
// benchmark.php
$iterations = 10000;

// Test 1: Array operations
$start = microtime(true);
for ($i = 0; $i < $iterations; $i++) {
    $arr = range(1, 1000);
    array_map(fn($x) => $x * 2, $arr);
}
$time1 = microtime(true) - $start;

// Test 2: String operations
$start = microtime(true);
for ($i = 0; $i < $iterations; $i++) {
    $str = str_repeat('test', 100);
    strlen($str);
    strtoupper($str);
}
$time2 = microtime(true) - $start;

echo "Array ops: {$time1}s\n";
echo "String ops: {$time2}s\n";
```

Run with both:

```bash
php benchmark.php
php-go benchmark.php
```

### Load Testing

Use Apache Bench, wrk, or similar tools:

```bash
# PHP-FPM
ab -n 10000 -c 100 http://localhost:8000/

# PHP-Go
ab -n 10000 -c 100 http://localhost:8000/
```

### Regression Testing

Automated testing for behavior changes:

```bash
# Test suite comparing outputs
for test in tests/*.php; do
    php "$test" > php.out
    php-go "$test" > phpgo.out
    if ! diff -q php.out phpgo.out > /dev/null; then
        echo "FAIL: $test"
        diff php.out phpgo.out
    else
        echo "PASS: $test"
    fi
done
```

## Performance Optimization

### Leverage Parallelization

PHP-Go can automatically parallelize certain operations:

```php
// Enable parallelization in config
// php-go.yaml:
// parallel:
//   enabled: true
//   array_operations: true
//   min_size: 1000

// This will be automatically parallelized if array is large enough
$result = array_map(function($x) {
    return expensiveOperation($x);
}, $largeArray);

// Explicitly request parallel execution (future feature)
$result = parallel_map(function($x) {
    return expensiveOperation($x);
}, $largeArray);
```

### Optimize Resource Usage

```php
// Bad: Creates many small allocations
$str = '';
for ($i = 0; $i < 10000; $i++) {
    $str .= 'x';
}

// Good: Single allocation
$str = str_repeat('x', 10000);

// Bad: Repeated array copies
$arr = [];
for ($i = 0; $i < 10000; $i++) {
    $arr[] = $i;
}

// Good: Pre-allocate size (future optimization)
$arr = array_fill(0, 10000, 0);
for ($i = 0; $i < 10000; $i++) {
    $arr[$i] = $i;
}
```

### Use PHP-Go Profiling

```bash
# Enable profiling
php-go -d phpgo.profile=cpu -d phpgo.profile_output=cpu.prof script.php

# Analyze profile (using Go tools)
go tool pprof -http=:8080 cpu.prof
```

### Monitor Performance Metrics

```yaml
# Enable metrics collection
metrics:
  enabled: true
  port: 9090

# Metrics available at http://localhost:9090/metrics:
# - phpgo_requests_total
# - phpgo_request_duration_seconds
# - phpgo_memory_usage_bytes
# - phpgo_goroutines
```

## Common Issues and Solutions

### Issue 1: Parse Errors

**Symptom:**
```
Parse error: Unexpected token 'IDENTIFIER', expected ';'
  at file.php:42:15
```

**Solution:**
- Check for syntax incompatible with PHP 8.4
- Verify all quotes and brackets are balanced
- Look for reserved keywords used incorrectly
- Check for malformed inline HTML/PHP mixing

### Issue 2: Missing Extension Functions

**Symptom:**
```
Fatal error: Call to undefined function mysqli_connect()
```

**Solution:**
```php
// Option 1: Check if function exists
if (!function_exists('mysqli_connect')) {
    die('mysqli extension not available in PHP-Go');
}

// Option 2: Polyfill
require_once 'polyfills/mysqli.php';

// Option 3: Use alternative
// Create database wrapper using Go integration
```

### Issue 3: Different Output

**Symptom:**
Results differ between PHP and PHP-Go for the same code.

**Common Causes:**
1. Floating point precision differences
2. Different random number generators
3. Timezone handling
4. Hash function implementations
5. Array iteration order (in rare cases)

**Solution:**
```php
// Floating point: Use epsilon comparison
function floatEquals($a, $b, $epsilon = 0.00001) {
    return abs($a - $b) < $epsilon;
}

// Random: Seed explicitly
mt_srand(12345);

// Timezone: Set explicitly
date_default_timezone_set('UTC');

// Hash: Verify algorithm support
if (!in_array('sha256', hash_algos())) {
    die('sha256 not supported');
}
```

### Issue 4: Performance Regression

**Symptom:**
PHP-Go is slower than PHP for your workload.

**Common Causes:**
1. Small datasets (parallelization overhead)
2. Excessive allocations
3. Missing optimizations
4. Non-optimized code paths

**Solution:**
```yaml
# Adjust parallelization threshold
parallel:
  enabled: true
  min_size: 10000  # Increase threshold

# Enable optimizations
optimize:
  inline: true
  constant_folding: true
```

```bash
# Profile to find bottleneck
php-go -d phpgo.profile=cpu script.php
```

### Issue 5: Memory Usage

**Symptom:**
Higher memory usage in PHP-Go.

**Cause:**
Go's garbage collector vs. PHP's reference counting.

**Solution:**
```yaml
# Configure GC
runtime:
  gc_percent: 100  # Default
  max_memory: 512M

# Monitor memory
metrics:
  enabled: true
```

```php
// Explicit cleanup for large data
$largeData = loadLargeDataset();
processData($largeData);
unset($largeData);  // Hint to GC
gc_collect_cycles();  // Force GC (if available)
```

### Issue 6: Extension Incompatibility

**Symptom:**
Code relies on PHP extension not available in PHP-Go.

**Solution:**

**Step 1:** Check if extension is planned:
- See TODO.md for extension roadmap
- Check [Extension Status](#extension-status)

**Step 2:** Find alternative:
```php
// Example: mbstring replacement
if (!function_exists('mb_strlen')) {
    // Use grapheme functions (if available)
    if (function_exists('grapheme_strlen')) {
        function mb_strlen($str, $encoding = 'UTF-8') {
            return grapheme_strlen($str);
        }
    } else {
        // Basic UTF-8 length
        function mb_strlen($str, $encoding = 'UTF-8') {
            return count(preg_split('//u', $str, -1, PREG_SPLIT_NO_EMPTY));
        }
    }
}
```

**Step 3:** Use Go library via integration:
```php
// Future feature: Call Go directly
$length = go\unicode\utf8\RuneCountInString($str);
```

**Step 4:** Request implementation:
- File issue on GitHub
- Contribute Go implementation
- Use external service meanwhile

## Rollback Strategy

### Pre-Migration Backup

Before migration:

```bash
# Backup entire application
tar -czf app-backup-$(date +%Y%m%d).tar.gz /var/www/app

# Backup database
mysqldump -u root -p database > database-backup-$(date +%Y%m%d).sql

# Backup configuration
cp -r /etc/php /etc/php-backup-$(date +%Y%m%d)
```

### Rollback Plan

If migration fails:

**Step 1: Stop PHP-Go**
```bash
systemctl stop php-go
# or
killall php-go
```

**Step 2: Restore PHP**
```bash
systemctl start php-fpm
systemctl start nginx  # or apache2
```

**Step 3: Restore Code** (if modified)
```bash
cd /var/www
rm -rf app
tar -xzf app-backup-20240101.tar.gz
```

**Step 4: Restore Database** (if modified)
```bash
mysql -u root -p database < database-backup-20240101.sql
```

**Step 5: Verify**
```bash
curl http://localhost
# Check application works correctly
```

### Blue-Green Deployment

Recommended for production:

```bash
# Keep both versions running
# PHP on port 8000
php-fpm -p 8000

# PHP-Go on port 8001
php-go serve -p 8001

# Route traffic via load balancer
# Gradually shift: 0% -> 10% -> 50% -> 100%
```

Rollback by changing load balancer configuration.

### Feature Flags

Use feature flags for gradual migration:

```php
// config.php
define('USE_PHPGO_FEATURES', getenv('PHPGO_FEATURES') === '1');

// In code
if (USE_PHPGO_FEATURES) {
    // Use PHP-Go specific features
    $result = parallel_map($func, $data);
} else {
    // Standard PHP code
    $result = array_map($func, $data);
}
```

## Migration Checklist

### Pre-Migration (1-2 weeks)

- [ ] Review this migration guide completely
- [ ] Assess current PHP version (target PHP 8.2+)
- [ ] Inventory all PHP extensions used
- [ ] List all critical standard library functions used
- [ ] Check framework compatibility (WordPress/Laravel/Symfony)
- [ ] Review [Extension Status](#extension-status) section
- [ ] Identify blocking dependencies
- [ ] Set up test environment with PHP-Go
- [ ] Install PHP-Go build tools
- [ ] Create migration timeline and plan
- [ ] Get team buy-in and training scheduled
- [ ] Set up monitoring and logging
- [ ] Create rollback procedure document
- [ ] Backup production environment

### Parse Testing (1 week)

- [ ] Run parse test on entire codebase
- [ ] Document parse success rate
- [ ] Identify and categorize parse errors
- [ ] Fix critical parse errors
- [ ] Update code for PHP 8.4 compatibility
- [ ] Test inline HTML/PHP mixing
- [ ] Verify alternative syntax support
- [ ] Test with all application entry points

### Functional Testing (2-4 weeks)

- [ ] Set up parallel test environment (PHP and PHP-Go)
- [ ] Run unit test suite on PHP-Go
- [ ] Document test failures and incompatibilities
- [ ] Create polyfills for missing functions
- [ ] Implement workarounds for missing extensions
- [ ] Run integration tests
- [ ] Test critical user workflows
- [ ] Verify data consistency
- [ ] Test error handling and edge cases
- [ ] Validate API responses match exactly
- [ ] Check database interactions
- [ ] Test file uploads and downloads
- [ ] Verify session management
- [ ] Test authentication and authorization
- [ ] Validate email sending
- [ ] Check cron jobs and background tasks

### Configuration (1 week)

- [ ] Migrate php.ini to php-go.yaml
- [ ] Configure resource limits
- [ ] Set up logging and error reporting
- [ ] Configure extensions
- [ ] Set up environment variables
- [ ] Document configuration differences
- [ ] Test configuration in dev/staging/prod environments
- [ ] Configure monitoring and metrics
- [ ] Set up health checks
- [ ] Configure graceful shutdown

### Performance Testing (1-2 weeks)

- [ ] Run baseline performance tests on PHP
- [ ] Run same tests on PHP-Go
- [ ] Compare response times
- [ ] Compare memory usage
- [ ] Compare CPU usage
- [ ] Test concurrent request handling
- [ ] Run load tests
- [ ] Profile hot code paths
- [ ] Optimize based on profiling results
- [ ] Enable parallelization where beneficial
- [ ] Tune configuration for workload
- [ ] Verify performance meets requirements

### Deployment Preparation (1 week)

- [ ] Create deployment scripts
- [ ] Set up blue-green or canary deployment
- [ ] Prepare rollback procedure
- [ ] Document deployment steps
- [ ] Create monitoring dashboards
- [ ] Set up alerts for errors and performance
- [ ] Train operations team
- [ ] Schedule deployment window
- [ ] Notify stakeholders
- [ ] Prepare communication plan

### Production Rollout (2-4 weeks)

- [ ] Deploy to staging environment
- [ ] Run full test suite in staging
- [ ] Perform smoke tests
- [ ] Start with 1-5% traffic
- [ ] Monitor errors and performance
- [ ] Gradually increase to 10%
- [ ] Continue monitoring
- [ ] Increase to 25%
- [ ] Increase to 50%
- [ ] Monitor for 24-48 hours
- [ ] Increase to 100%
- [ ] Monitor for 1 week
- [ ] Document any issues encountered
- [ ] Optimize based on production metrics

### Post-Migration (Ongoing)

- [ ] Remove compatibility shims
- [ ] Optimize for PHP-Go features
- [ ] Enable parallelization
- [ ] Implement Go integration where beneficial
- [ ] Update documentation
- [ ] Share learnings with team
- [ ] Contribute improvements to PHP-Go
- [ ] Monitor long-term stability
- [ ] Plan next optimization phase

## Extension Status

### Fully Implemented

| Extension | Functions | Status | Notes |
|-----------|-----------|--------|-------|
| JSON | 4 core | ✅ Complete | json_encode, json_decode, error functions |
| Hash | 12 core | ✅ Partial | md5, sha*, hmac functions; missing 4 algorithms |
| Core | 50+ | ✅ Partial | Variables, arrays, strings, output, types |

### Partially Implemented

| Extension | Status | Missing | ETA |
|-----------|--------|---------|-----|
| Date/Time | 🟡 30% | DateTime class, timezone functions | Phase 6 |
| SPL | 🟡 20% | Most iterators, data structures | Phase 6 |
| PCRE | 🟡 40% | Some preg_* functions | Phase 6 |
| Strings | 🟡 60% | Advanced functions | Phase 6 |
| Arrays | 🟡 50% | Advanced manipulation | Phase 6 |

### Not Yet Implemented

| Extension | Priority | Estimated Effort | ETA |
|-----------|----------|------------------|-----|
| mysqli | 🔴 Critical | 40-60h | Phase 6-7 |
| PDO | 🔴 Critical | 50-70h | Phase 6-7 |
| ctype | 🟡 High | 4-6h | Phase 6 |
| filter | 🟡 High | 8-12h | Phase 6 |
| mbstring | 🟡 High | 20-30h | Phase 6-7 |
| curl | 🟡 High | 15-20h | Phase 6-7 |
| xml | 🟢 Medium | 25-35h | Phase 7 |
| dom | 🟢 Medium | 30-40h | Phase 7 |
| simplexml | 🟢 Medium | 15-20h | Phase 7 |
| gd | 🟢 Medium | 40-60h | Phase 7-8 |
| openssl | 🟢 Medium | 30-50h | Phase 7 |
| zlib | 🟢 Low | 10-15h | Phase 8 |
| zip | 🟢 Low | 15-20h | Phase 8 |
| pgsql | 🟢 Low | 30-40h | Phase 8 |
| redis | 🟢 Low | 20-30h | Phase 8 |
| memcached | 🟢 Low | 15-25h | Phase 8 |

Priority:
- 🔴 Critical: Required for most PHP applications
- 🟡 High: Commonly used, significant impact
- 🟢 Medium/Low: Nice to have, specific use cases

## Standard Library Status

### String Functions

**Implemented:**
- strlen(), substr(), str_repeat(), strpos(), strtolower(), strtoupper()

**Missing (Phase 6):**
- str_replace(), str_pad(), trim(), ltrim(), rtrim()
- explode(), implode(), str_split()
- sprintf(), vsprintf()
- wordwrap(), str_word_count()
- strcmp(), strcasecmp(), strncmp()
- ucfirst(), ucwords(), lcfirst()
- And ~40 more functions

### Array Functions

**Implemented:**
- count(), array_keys(), array_values(), array_merge()
- array_push(), array_pop(), array_shift(), array_unshift()

**Missing (Phase 6):**
- array_map(), array_filter(), array_reduce()
- array_slice(), array_splice(), array_chunk()
- array_diff(), array_intersect(), array_unique()
- sort(), rsort(), asort(), arsort(), ksort(), krsort()
- in_array(), array_search(), array_key_exists()
- And ~50 more functions

### File I/O Functions

**Implemented:**
- Basic file operations

**Missing (Phase 6):**
- file_get_contents(), file_put_contents()
- fopen(), fclose(), fread(), fwrite()
- file_exists(), is_file(), is_dir()
- mkdir(), rmdir(), unlink()
- glob(), scandir()
- And ~30 more functions

## Additional Resources

### Documentation

- [Installation Guide](installation.md) - Detailed installation instructions
- [Configuration Guide](configuration.md) - Complete configuration reference
- [Extension Development Guide](extension-development.md) - Writing Go extensions
- [Performance Tuning Guide](performance-tuning.md) - Optimization strategies
- [API Reference](../api/README.md) - Complete API documentation

### Community

- **GitHub Repository**: https://github.com/krizos/php-go
- **Issue Tracker**: https://github.com/krizos/php-go/issues
- **Discussions**: https://github.com/krizos/php-go/discussions

### Getting Help

1. Check this migration guide
2. Review [FAQ](#faq) section in User Guide
3. Search existing GitHub issues
4. Ask in GitHub Discussions
5. File a bug report with reproduction steps

## FAQ

### Q: When should I migrate to PHP-Go?

**A:** Migrate when:
- Your application is on PHP 8.2+
- You need better multi-core utilization
- You want single-binary deployment
- Critical extensions are available in PHP-Go
- You have time for testing and validation

Wait if:
- You rely on unavailable extensions (mysqli, PDO, curl, etc.)
- Your application is on PHP 7.x or below
- You need 100% compatibility immediately

### Q: Is PHP-Go production-ready?

**A:** PHP-Go is currently at ~93% completion (Phase 10 in progress). It's suitable for:
- ✅ Greenfield projects
- ✅ Internal tools
- ✅ Non-critical services
- ✅ Development/testing environments

Wait for v1.0 for:
- ⏳ Critical production systems
- ⏳ Applications requiring database extensions
- ⏳ Applications requiring XML/image processing
- ⏳ High-traffic public websites

### Q: How compatible is PHP-Go with standard PHP?

**A:** PHP-Go aims for 100% PHP 8.4 language compatibility:
- ✅ Language features: ~100% (all syntax, operators, control structures)
- ✅ Core functions: ~40-60% (basic operations implemented)
- 🟡 Extensions: ~10-20% (JSON, Hash, partial Date/SPL)
- ⏳ Standard library: ~30-40% (growing rapidly)

See [Extension Status](#extension-status) for details.

### Q: Will my existing code work without changes?

**A:** Most modern PHP 8.2+ code will work with minimal changes:
- Pure PHP logic: Usually works as-is
- Standard library: May need polyfills for missing functions
- Extensions: May need alternatives for unavailable extensions
- Configuration: May need migration to php-go.yaml

Use the [Parse Compatibility Test](#step-1-parse-compatibility-test) to check your code.

### Q: How do I handle missing extensions?

**A:** Multiple options:
1. Wait for extension implementation (check roadmap in TODO.md)
2. Write Go extension (see [Extension Development Guide](extension-development.md))
3. Use polyfills for simple functions
4. Use external services (microservices, REST APIs)
5. Use Go libraries via PHP-Go integration (future)

### Q: What's the performance difference?

**A:** Performance varies by workload:
- **Parser/Compiler**: Similar to PHP 8.4
- **Type operations**: Excellent (sub-ns to 40ns)
- **VM execution**: Currently 2-5x slower (optimizations ongoing)
- **Parallel workloads**: Can be faster than PHP due to multi-threading
- **Expected v1.0**: Within 2x of PHP 8.4 + opcache

See [Performance Benchmarks](../TODO.md#106-performance-benchmarks-20h) for details.

### Q: Can I run WordPress/Laravel/Symfony on PHP-Go?

**A:** Current status:

**WordPress**: 61% parse rate
- ✅ Most core files parse correctly
- ❌ Blocked by missing mysqli/PDO
- ⏳ Estimated 2-3 months for full support

**Laravel**: 43% parse rate
- ✅ Critical files parse (with features)
- ❌ Needs more advanced PHP 8.2+ features
- ⏳ Estimated 3-4 months for full support

**Symfony**: 50% parse rate
- ✅ Core application files parse
- ❌ Needs generators, first-class callables
- ⏳ Estimated 3-4 months for full support

### Q: How do I contribute?

**A:** Contributions welcome!
1. Check TODO.md for open tasks
2. Read CONTRIBUTING.md (if available)
3. Fork repository
4. Implement feature or fix
5. Write tests (85%+ coverage required)
6. Submit pull request

Focus areas:
- Standard library functions
- Extensions (especially mysqli, PDO)
- Performance optimizations
- Bug fixes
- Documentation

### Q: What's the long-term roadmap?

**A:** See TODO.md for complete roadmap:

- **Phase 10** (Current, 93% complete): Testing, optimization, documentation
- **v1.0** (Target: Q2 2025): Production-ready release
  - 95%+ PHP test suite pass rate
  - Critical extensions (mysqli, PDO, curl)
  - Performance parity with PHP 8.4
  - Complete standard library

- **v1.x** (2025-2026): Additional features
  - All common extensions
  - Advanced parallelization
  - JIT compiler for hot paths
  - Performance exceeding PHP

### Q: How do I report bugs or request features?

**A:** Use GitHub:
1. Search existing issues first
2. For bugs: Provide minimal reproduction case
3. For features: Explain use case and benefits
4. Include PHP-Go version and environment info

**Bug Report Template:**
```markdown
## Description
Brief description of the bug

## Reproduction
1. Create file `test.php` with:
   <?php
   // Minimal code to reproduce
   ?>
2. Run: php-go test.php
3. Expected: ...
4. Actual: ...

## Environment
- PHP-Go version: vX.Y.Z
- OS: Linux/macOS/Windows
- PHP version compared: 8.4.0
```

## Conclusion

Migrating to PHP-Go offers significant benefits for performance, deployment simplicity, and Go ecosystem integration. While PHP-Go is still maturing, it's rapidly approaching production readiness.

**Success Factors:**
1. **Thorough Assessment**: Understand your dependencies
2. **Incremental Approach**: Migrate gradually, not all at once
3. **Comprehensive Testing**: Test early and often
4. **Community Engagement**: Report issues, contribute fixes
5. **Patience**: Some features are still in development

**Next Steps:**
1. Complete [Pre-Migration Assessment](#pre-migration-assessment)
2. Set up test environment
3. Run parse compatibility tests
4. Identify and document blockers
5. Create detailed migration plan
6. Start with non-critical components
7. Monitor and optimize

For questions or help, visit the [GitHub Discussions](https://github.com/krizos/php-go/discussions) or file an issue.

Happy migrating! 🚀
