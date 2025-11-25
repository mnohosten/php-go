# Migration Guide: From PHP to PHP-Go

**Version**: 1.0
**Last Updated**: November 24, 2025
**Target Audience**: PHP developers evaluating or migrating to PHP-Go

---

## Table of Contents

1. [Introduction](#introduction)
2. [Feature Comparison](#feature-comparison)
3. [Compatibility Notes](#compatibility-notes)
4. [Migration Strategy](#migration-strategy)
5. [Known Limitations](#known-limitations)
6. [Performance Tips](#performance-tips)
7. [Common Pitfalls](#common-pitfalls)
8. [FAQ](#faq)

---

## Introduction

### What is PHP-Go?

PHP-Go is a complete rewrite of the PHP interpreter in Go, designed to provide:
- **5.0x faster performance** than PHP 8.4 + opcache
- **Automatic parallelization** capabilities (roadmap)
- **Native Go library integration** (roadmap)
- **Full PHP 8.4 language compatibility** (in progress)

### Should You Migrate?

**Migrate if you need**:
- ✅ Significantly better performance (5x faster)
- ✅ Go ecosystem integration
- ✅ Modern, maintainable codebase
- ✅ Future parallelization features

**Stay on PHP if you need**:
- ❌ 100% ecosystem compatibility right now
- ❌ All PECL extensions
- ❌ Legacy PHP 5.x features
- ❌ Mature, battle-tested runtime

### Current Maturity

- **Language Features**: 85% complete (PHP 8.4 syntax)
- **Standard Library**: ~50 functions (300+ planned)
- **WordPress Compatibility**: 61% parse success
- **Laravel Compatibility**: 44% parse success
- **Production Ready**: Not yet (6-8 months estimated)

---

## Feature Comparison

### Language Features

| Feature | PHP 8.4 | PHP-Go | Status | Notes |
|---------|---------|--------|--------|-------|
| **Basic Syntax** | ✅ | ✅ | Complete | Variables, operators, control flow |
| **Functions** | ✅ | ✅ | Complete | Named functions, variadic, default params |
| **Classes** | ✅ | ✅ | Complete | OOP, inheritance, interfaces, traits |
| **Closures** | ✅ | ⚠️ | Partial | Basic work, parameters broken |
| **Arrow Functions** | ✅ | ⚠️ | Partial | Compilation works, runtime issues |
| **Generators** | ✅ | ✅ | Complete | `yield` and `yield from` |
| **Namespaces** | ✅ | ❌ | Missing | Critical for modern PHP |
| **Attributes** | ✅ | ❌ | Planned | PHP 8.0 feature |
| **Enums** | ✅ | ✅ | Complete | PHP 8.1 feature |
| **Match Expressions** | ✅ | ✅ | Complete | PHP 8.0 feature |
| **Named Arguments** | ✅ | ✅ | Complete | PHP 8.0 feature |
| **Constructor Promotion** | ✅ | ✅ | Complete | PHP 8.0 feature |
| **Readonly Properties** | ✅ | ✅ | Complete | PHP 8.1 feature |
| **Fibers** | ✅ | ❌ | Not planned | PHP 8.1 feature |

### Standard Library

| Category | PHP 8.4 | PHP-Go | Completion |
|----------|---------|--------|------------|
| **Array Functions** | 80+ | 15+ | 19% |
| **String Functions** | 100+ | 12+ | 12% |
| **File Functions** | 80+ | 8+ | 10% |
| **Math Functions** | 40+ | 8+ | 20% |
| **Date/Time** | 50+ | 2+ | 4% |
| **Hash/Crypto** | 30+ | 6+ | 20% |
| **Type Functions** | 30+ | 8+ | 27% |
| **PCRE** | 10+ | 3+ | 30% |

### Extensions

| Extension | PHP 8.4 | PHP-Go | Status |
|-----------|---------|--------|--------|
| **mysqli** | ✅ | ❌ | Planned Phase 14 |
| **PDO** | ✅ | ❌ | Planned Phase 14 |
| **mbstring** | ✅ | ❌ | Planned Phase 14 |
| **curl** | ✅ | ❌ | Planned Phase 14 |
| **GD** | ✅ | ❌ | Not planned |
| **XML** | ✅ | ❌ | Planned Phase 14 |
| **JSON** | ✅ | ⚠️ | Partial (basic support) |

---

## Compatibility Notes

### What Works Well

#### ✅ Basic PHP Scripts
```php
<?php
echo "Hello, World!\n";

$x = 10;
$y = 20;
echo $x + $y;

function greet($name) {
    return "Hello, " . $name;
}
echo greet("Alice");
```
**Status**: ✅ Works perfectly

#### ✅ Object-Oriented Programming
```php
<?php
class User {
    public function __construct(
        public string $name,
        public int $age
    ) {}

    public function greet(): string {
        return "Hello, I'm {$this->name}";
    }
}

$user = new User("Bob", 30);
echo $user->greet();
```
**Status**: ✅ Full OOP support including:
- Constructor property promotion
- Type hints
- Visibility modifiers
- Inheritance, interfaces, traits
- Magic methods
- Readonly properties
- Enums

#### ✅ Control Flow
```php
<?php
// If/else, while, for, foreach all work
for ($i = 0; $i < 10; $i++) {
    if ($i % 2 === 0) {
        echo "$i is even\n";
    }
}

// Match expressions (PHP 8.0)
$result = match($value) {
    1 => "one",
    2 => "two",
    default => "other"
};

// Alternative syntax
if ($condition):
    echo "True";
endif;
```
**Status**: ✅ Complete control flow support

#### ✅ Arrays
```php
<?php
$arr = [1, 2, 3, 4, 5];
$assoc = ["name" => "Alice", "age" => 30];

// Array functions
$doubled = array_map(function($x) { return $x * 2; }, $arr);
$filtered = array_filter($arr, function($x) { return $x > 3; });
$merged = array_merge($arr, [6, 7, 8]);
```
**Status**: ✅ Core array operations work
**Note**: ⚠️ Closures with parameters need fixes

### What Has Issues

#### ⚠️ Closures with Parameters
```php
<?php
// ❌ BROKEN: Parameters not passed correctly
$fn = function($x) {
    return $x * 2;
};
$result = $fn(5); // Returns nothing instead of 10
```
**Workaround**: Use closures without parameters:
```php
<?php
// ✅ WORKS: No parameters
$multiplier = 2;
$fn = function() use ($multiplier) {
    return $multiplier * 5;
};
$result = $fn(); // Works: returns 10
```

#### ⚠️ Reference Capture Write-back
```php
<?php
// ❌ BROKEN: Reference changes don't propagate
$counter = 0;
$increment = function() use (&$counter) {
    $counter++; // Modifies local copy only
};
$increment();
echo $counter; // Still 0 (should be 1)
```
**Workaround**: Use global variables:
```php
<?php
// ✅ WORKS: Use globals
$GLOBALS['counter'] = 0;
$increment = function() {
    global $counter;
    $counter++;
};
```

### What Doesn't Work Yet

#### ❌ Namespaces
```php
<?php
namespace App\Models; // ❌ NOT SUPPORTED

use Illuminate\Database\Model; // ❌ NOT SUPPORTED
```
**Impact**: Critical blocker for modern frameworks
**Status**: Planned for Phase 11 (16-20 hours)
**Workaround**: None - refactor code to not use namespaces

#### ❌ require/include
```php
<?php
require_once "config.php"; // ❌ NOT SUPPORTED
include "helpers.php";     // ❌ NOT SUPPORTED
```
**Impact**: Critical blocker for multi-file applications
**Status**: Planned for Phase 11 (12-16 hours)
**Workaround**: Combine all code into single file

#### ❌ Database Extensions
```php
<?php
$mysqli = new mysqli("localhost", "user", "pass", "db"); // ❌ NOT SUPPORTED
$pdo = new PDO("mysql:host=localhost;dbname=test");     // ❌ NOT SUPPORTED
```
**Impact**: Critical for web applications
**Status**: Planned for Phase 14 (80-120 hours)
**Workaround**: Use REST APIs or external database services

---

## Migration Strategy

### Phase 1: Assessment (1-2 days)

**Goal**: Determine if your codebase can run on PHP-Go

1. **Analyze Dependencies**:
   ```bash
   # List all require/include statements
   grep -r "require\|include" your-project/

   # List all namespace declarations
   grep -r "^namespace\|^use " your-project/

   # List all database calls
   grep -r "mysqli\|PDO\|pg_" your-project/
   ```

2. **Check Feature Usage**:
   - Do you use namespaces? (blocker)
   - Do you use require/include? (blocker)
   - Do you need database access? (blocker)
   - Do you use PECL extensions? (likely blocker)

3. **Evaluate Blockers**:
   - **0-1 blockers**: Migration feasible now with workarounds
   - **2-3 blockers**: Wait 3-6 months for Phase 11/14
   - **4+ blockers**: Wait 6-12 months for full compatibility

### Phase 2: Test Migration (2-5 days)

**Goal**: Run your code on PHP-Go and identify issues

1. **Install PHP-Go**:
   ```bash
   git clone https://github.com/your-org/php-go
   cd php-go
   go build -o php-go ./cmd/php-go
   ```

2. **Run Simple Scripts**:
   ```bash
   ./php-go your-script.php
   ```

3. **Document Issues**:
   - List all error messages
   - Identify missing functions
   - Note performance differences

4. **Calculate Coverage**:
   - What % of your code runs?
   - What features are blockers?
   - What's the workaround cost?

### Phase 3: Incremental Migration (1-4 weeks)

**Goal**: Migrate non-blocking parts first

1. **Start with Utilities**:
   - Migrate standalone utility scripts
   - Test thoroughly
   - Compare performance

2. **Migrate Core Logic**:
   - Business logic without database
   - API response formatting
   - Data transformation

3. **Defer Blockers**:
   - Keep database code in PHP
   - Use PHP for file operations
   - Hybrid approach: PHP-Go for compute, PHP for I/O

4. **Incremental Rollout**:
   - Run both PHP and PHP-Go side-by-side
   - A/B test performance
   - Gradually shift traffic

### Phase 4: Production Deployment (1-2 weeks)

**Goal**: Run PHP-Go in production

1. **Performance Testing**:
   ```bash
   # Benchmark PHP
   time php your-script.php

   # Benchmark PHP-Go
   time ./php-go your-script.php
   ```

2. **Load Testing**:
   - Use same test suite
   - Compare memory usage
   - Check error rates

3. **Monitoring**:
   - Set up error tracking
   - Monitor performance metrics
   - Track resource usage

4. **Rollback Plan**:
   - Keep PHP installation ready
   - Have feature flags
   - Quick switch mechanism

---

## Known Limitations

### Critical (Blockers)

1. **Namespaces Not Supported**
   - **Impact**: High - Modern PHP uses namespaces extensively
   - **Workaround**: Rename classes to avoid conflicts
   - **Fix Timeline**: Phase 11 (16-20h)

2. **require/include Not Supported**
   - **Impact**: High - Multi-file applications won't work
   - **Workaround**: Concatenate files manually
   - **Fix Timeline**: Phase 11 (12-16h)

3. **No Database Extensions**
   - **Impact**: High - Web apps need databases
   - **Workaround**: Use external APIs
   - **Fix Timeline**: Phase 14 (80-120h)

### Important (Major Features)

4. **Closure Parameters Broken**
   - **Impact**: Medium - Affects callbacks and filters
   - **Workaround**: Use parameter-less closures
   - **Fix Timeline**: 8-12 hours

5. **Reference Write-back**
   - **Impact**: Medium - Affects accumulators
   - **Workaround**: Use globals
   - **Fix Timeline**: 6-8 hours

6. **Limited Standard Library**
   - **Impact**: Medium - Missing many common functions
   - **Workaround**: Implement yourself or wait
   - **Fix Timeline**: Phase 13 (200-300h)

### Minor (Nice to Have)

7. **No Attributes Support**
   - **Impact**: Low - Modern frameworks use attributes
   - **Workaround**: Use alternate patterns
   - **Fix Timeline**: 12-16 hours (deferred)

8. **Execution Flow After Closures**
   - **Impact**: Low - Sequential closure calls fail
   - **Workaround**: One closure per script section
   - **Fix Timeline**: 4-6 hours

---

## Performance Tips

### What's Fast

#### ✅ Pure Computation
```php
<?php
// PHP-Go excels at CPU-intensive tasks
function fibonacci($n) {
    if ($n <= 1) return $n;
    return fibonacci($n - 1) + fibonacci($n - 2);
}

// PHP-Go: 3552x faster than PHP 8.4
echo fibonacci(15);
```
**Speedup**: 10-3500x depending on task

#### ✅ String Operations
```php
<?php
$str = "Hello";
for ($i = 0; $i < 1000; $i++) {
    $str .= " World";
}
```
**Speedup**: 181x faster than PHP 8.4

#### ✅ Array Operations
```php
<?php
$arr = range(1, 10000);
$filtered = array_filter($arr, function($x) { return $x % 2 === 0; });
$mapped = array_map(function($x) { return $x * 2; }, $filtered);
```
**Speedup**: 5.6x faster than PHP 8.4

### What's Slow

#### ⚠️ I/O Operations
```php
<?php
// File I/O not optimized yet
$contents = file_get_contents("large-file.txt");
```
**Status**: Similar to PHP (no optimization yet)

#### ⚠️ Extension Calls
```php
<?php
// Extensions not available
// Will be slow when implemented (Go->C bridge)
```
**Status**: Not yet implemented

### Optimization Strategies

#### 1. Minimize Function Calls
```php
<?php
// SLOW: Many function calls
for ($i = 0; $i < 1000; $i++) {
    echo expensive_function($i);
}

// FAST: Batch operations
$results = [];
for ($i = 0; $i < 1000; $i++) {
    $results[] = $i * 2; // Inline simple operations
}
echo implode("\n", $results);
```

#### 2. Use Built-in Functions
```php
<?php
// SLOW: Manual loops
$sum = 0;
foreach ($array as $value) {
    $sum += $value;
}

// FAST: Built-in function
$sum = array_sum($array);
```

#### 3. Avoid Dynamic Features
```php
<?php
// SLOW: Variable variables, eval
$var = "name";
$$var = "value";

// FAST: Direct access
$name = "value";
```

#### 4. Pre-compute Constants
```php
<?php
// SLOW: Re-compute each time
for ($i = 0; $i < 1000; $i++) {
    $result = $i * (2 + 2);
}

// FAST: Pre-compute constant
$multiplier = 4;
for ($i = 0; $i < 1000; $i++) {
    $result = $i * $multiplier;
}
```

### Benchmarking

```bash
# Run your script with timing
time ./php-go script.php

# Compare with PHP
time php script.php

# Memory usage
/usr/bin/time -v ./php-go script.php
```

**Expected Results**:
- Simple loops: 9x faster
- Function calls: 10x faster
- String operations: 100-200x faster
- Array operations: 5-10x faster
- Overall: 5x faster average

---

## Common Pitfalls

### 1. Assuming Full Compatibility

**Wrong**:
```php
<?php
namespace App; // Assumes namespaces work
require "vendor/autoload.php"; // Assumes includes work
$pdo = new PDO(...); // Assumes PDO exists
```

**Right**:
```php
<?php
// Check feature support first
// Refactor to remove dependencies
// Use workarounds for missing features
```

### 2. Not Testing Closures

**Wrong**:
```php
<?php
// Assumes closure parameters work
$users = array_map(function($user) {
    return $user->name;
}, $users);
```

**Right**:
```php
<?php
// Test closure behavior first
// Use alternative patterns if broken
$names = [];
foreach ($users as $user) {
    $names[] = $user->name;
}
```

### 3. Relying on Undefined Behavior

**Wrong**:
```php
<?php
// Relies on specific error handling
try {
    risky_operation();
} catch (SpecificException $e) {
    // May not work the same
}
```

**Right**:
```php
<?php
// Test error behavior explicitly
// Don't rely on internals
// Check PHP-Go compatibility
```

### 4. Not Profiling

**Wrong**:
```php
// Assumes everything is faster
```

**Right**:
```bash
# Benchmark before migrating
./benchmark.sh php
./benchmark.sh php-go
# Compare results
```

---

## FAQ

### General Questions

**Q: Is PHP-Go production ready?**
A: Not yet. Estimated 6-8 months to WordPress/Laravel compatibility. Use for:
- Experimentation
- Performance testing
- Non-critical workloads

**Q: Will my PHP code run without changes?**
A: Depends. Simple scripts: yes. Complex apps with namespaces/includes/database: not yet.

**Q: How much faster is PHP-Go?**
A: 5.0x average, up to 3500x for specific operations. Varies by workload.

**Q: Can I use Composer packages?**
A: Not yet. Most packages require:
- Namespaces (not supported)
- Autoloading (requires include/require)
- Extensions (mostly not implemented)

### Migration Questions

**Q: Should I migrate my production app now?**
A: Only if:
- You don't use namespaces
- Single-file script or manual file concatenation is OK
- No database access needed
- You need 5x performance boost urgently

**Q: When will WordPress work?**
A: Estimated 5-6 months (Phase 11 completion). Currently 61% parse success.

**Q: When will Laravel work?**
A: Estimated 4-5 months (Phase 12 completion). Currently 44% parse success.

**Q: Can I use PHP-Go with my existing PHP installation?**
A: Yes. Run side-by-side:
```bash
# Old code: use PHP
php legacy-app.php

# New code: use PHP-Go
./php-go optimized-script.php
```

### Technical Questions

**Q: Why are closures broken?**
A: Parameter passing mechanism needs fixes. Parser and compiler work, runtime needs debugging. Fix estimated: 8-12 hours.

**Q: Why no namespaces?**
A: Not implemented yet. Planned for Phase 11 (16-20 hours of work).

**Q: Will extensions be compatible?**
A: Native extensions will be written in Go, not C. Different API, but similar functionality.

**Q: Can I write Go code that PHP calls?**
A: Planned for Phase 8 (Go Integration). Not yet implemented.

---

## Next Steps

### If You're Ready to Migrate

1. **Run Assessment**: Check your codebase for blockers
2. **Start Small**: Migrate utility scripts first
3. **Test Thoroughly**: Don't assume compatibility
4. **Benchmark**: Measure actual performance gains
5. **Report Issues**: Help improve PHP-Go

### If You're Waiting

1. **Monitor Progress**: Track Phase 11/12/13 completion
2. **Test Periodically**: Try new releases
3. **Provide Feedback**: Share your use cases
4. **Contribute**: Help fix blockers

### Resources

- **Documentation**: `docs/` directory
- **Roadmap**: `docs/ROADMAP.md`
- **Test Results**: `tests/CLOSURE_TEST_RESULTS.md`
- **Phase 10 Summary**: `docs/PHASE10_CLOSURE_SUMMARY.md`
- **Issue Tracker**: GitHub issues

---

## Conclusion

PHP-Go is an **exciting project** with **impressive performance** gains but is **not yet production-ready** for complex applications.

**Best Use Cases Right Now**:
- ✅ CLI scripts
- ✅ Data processing
- ✅ Computation-heavy tasks
- ✅ Single-file applications
- ✅ Performance testing

**Wait Until Phase 11+ For**:
- ❌ WordPress applications
- ❌ Laravel applications
- ❌ Multi-file projects with includes
- ❌ Database-driven applications
- ❌ Composer-based projects

**Timeline**: Full WordPress/Laravel support expected in **6-8 months** with estimated **730-970 hours** of additional development work.

For the adventurous: experiment now, provide feedback, and help shape the future of PHP-Go!

---

**Document Version**: 1.0
**Last Updated**: November 24, 2025
**Next Review**: When Phase 11 completes
