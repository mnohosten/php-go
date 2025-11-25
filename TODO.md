# PHP-Go TODO - Roadmap to WordPress & Laravel Support

**Ultimate Goal**: Full support for WordPress 6.8.3 and Laravel v12.10.1
**Current Status**: Phase 10 (93.6% complete, 1340/1430 hours)
**See also**: `docs/ROADMAP.md` for comprehensive roadmap documentation

**Status Legend**:
- ⬜ Not Started
- 🔄 In Progress
- ✅ Complete
- ⏸️ Blocked
- ⏭️ Deferred

---

## Current Progress Summary

**Phases 0-9**: ✅ Complete (1340 hours)
- ✅ Foundation, compiler, VM, object system
- ✅ 5.0x faster than PHP 8.4 on benchmarks
- ✅ 50+ standard library functions
- ✅ Security hardening (ReDoS, overflow protection)

**Phase 10**: 🔄 In Progress (90h remaining)
- ✅ Named arguments, generators, match expressions
- ✅ Array spread, first-class callables, throw expressions
- ✅ Constructor property promotion, alternative syntax
- ❌ Closures/anonymous functions (15-25h) - **LAST MAJOR BLOCKER**
- ❌ Documentation and polish (20h)

**Framework Compatibility**:
- WordPress 6.8.3: 61.20% parse success (768/1255 files)
- Laravel v12.10.1: 43.57% parse success (3260/7483 files)
- Symfony 7.3: 57.5% parse success

**Estimated Path to Production**: 730-970h (6-8 months)

---

## Phase 10: Testing & Production Readiness (Remaining: 90h)

**Status**: 🔄 In Progress (43.3% complete, 103/240 hours)

### 10.16 Closures and Anonymous Functions (15-25h) ⬜ NOT STARTED

**Priority**: P0 - CRITICAL (Last major language blocker)
**Impact**: Blocks 40% of WordPress code (hooks, callbacks), Laravel functional programming
**Blocks**: WordPress hooks (`add_action`, `add_filter`), `array_map`, `array_filter` with callbacks

- [x] Parse closure expression syntax (2h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/expr.go` (lines 1171-1221)
  - Parse `function($params) { body }` as expression ✅
  - Parse `function($params) use ($vars) { body }` with variable capture ✅
  - Handle both value and reference capture (`use ($x, &$y)`) ✅
  - **Status**: Fully implemented in `parseClosureExpression()` with support for:
    - Reference return (`&function`)
    - Parameters with type hints
    - Use clause with value and reference capture
    - Return type hints
    - Block statement body
  - **Tested**: Successfully parses all closure patterns

- [x] Implement variable capture analysis (5h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/compiler/compiler.go` (lines 3341-3464)
  - Analyze closure body to identify captured variables ✅
  - Distinguish between captured and local variables ✅
  - Build capture list with reference flags ✅
  - **Status**: Fully implemented with two functions:
    - `findReferencedVariables(expr)` - Main entry point, returns list of variable names
    - `findVarsRecursive(node, vars)` - Recursive visitor pattern traversing all AST nodes
  - **Features**:
    - Handles all expression types (infix, prefix, ternary, calls, arrays, etc.)
    - Properly handles nested closures (doesn't traverse into their bodies)
    - Tracks variables in use clauses
    - Used by arrow function auto-capture (lines 1127-1146)
    - Used by closure use clause binding (lines 1105-1118)
  - **Tested**: Successfully analyzes variable references in closures

- [x] Create closure compilation (4h) ⚠️ PARTIALLY IMPLEMENTED - NEEDS FIXES
  - File: `pkg/compiler/compiler.go` (lines 1023-1120)
  - Compile closure as anonymous function ✅ (structure exists)
  - Generate unique closure name ❌ **BLOCKER**: Not generating function names
  - Emit OpClosure instruction with capture list ⚠️ Uses `OpDeclareLambdaFunction` but parameters mismatch
  - Store function metadata ❌ **BLOCKER**: Not storing CompiledFunction in vm.functions map
  - **Status**: Compilation structure exists but has critical issues:
    1. **Compiler emits**: flags, start_pos, end_pos (lines 1099-1103)
    2. **VM handler expects**: function_name in Op1 (handlers_closure.go:19-29)
    3. **Missing**: Function registration in vm.functions map
    4. **Missing**: Unique name generation for closures (e.g., "{closure}#1")
  - **Testing**: Closures parse and compile but don't execute correctly
  - **Next**: Need to fix the compiler-VM interface mismatch

- [x] Implement closure runtime support (4h) ✅ IMPLEMENTED - Basic closures working!
  - File: `pkg/vm/closure.go`, `pkg/vm/handlers_closure.go`, `pkg/vm/handlers_functions.go`
  - Create Closure value type with captured variables ✅
  - Implement OpClosure handler to create closure objects ✅
  - Support closure invocation via DO_FCALL ✅
  - Handle variable capture (value and reference) ✅ (basic support)
  - **Implementation**:
    1. Fixed `opDeclareLambdaFunction` to extract and register closure instructions (handlers_closure.go:19-112)
    2. Updated `opInitFcallByName` to detect and handle closure resources (handlers_functions.go:92-151)
    3. Added `pendingClosure` field to Frame struct (frame.go:38)
    4. Updated `opDoFcall` to handle closure calls with captured variables (handlers_functions.go:215-338)
    5. Added `closureCounter` to Compiler for unique closure names (compiler.go:44)
    6. Fixed closure compilation to use proper operands (compiler.go:1026-1124)
  - **Testing**:
    - ✅ Basic closures work: `$fn = function() { echo "Works!"; }; $fn();` outputs "Works!"
    - ⚠️ Parameters have issue: `$fn = function($x) { echo $x; }; $fn(42);` outputs nothing
    - ⚠️ Use clause not fully tested yet
  - **Known Issues**:
    - Parameters not being received correctly in closures (RECV opcodes may need adjustment)
    - Need to debug parameter passing mechanism
  - **Next**: Fix parameter passing, then test use clause

- [x] Add by-reference capture support (3h) ✅ IMPLEMENTED - Infrastructure complete
  - File: `pkg/compiler/compiler.go`, `pkg/vm/closure.go`, `pkg/vm/handlers_closure.go`
  - Parse `use (&$var)` syntax ✅ (already implemented in parser)
  - Implement reference capture in closure creation ✅
  - Ensure references update outer scope ⚠️ (partial)
  - **Implementation**:
    1. Parser already supports `use (&$var)` syntax (parseClosureExpression)
    2. Compiler emits `OpBindLexical` with byRef flag (compiler.go:1109-1122)
    3. `opBindLexical` handler processes reference flag (handlers_closure.go:119-174)
    4. `Closure.BindVariable()` creates Reference type for by-ref captures (closure.go:35-50)
    5. `GetCapturedVariable()` properly dereferences References (closure.go:52-65)
  - **Testing**:
    - ✅ Value capture works: `use ($x)` - reads correctly
    - ✅ Reference capture read works: `use (&$x)` - reads correctly
    - ⚠️ Reference capture write partial: `use (&$x)` with `$x++` doesn't update outer scope
  - **Known Limitation**:
    - Captured variables are currently copied to globals (opDoFcall:302-306)
    - Reference updates within closure don't propagate back to outer scope
    - Full fix requires proper CV (Compiled Variable) slot management
    - This is a known simplification documented in code comments
  - **Status**: Core infrastructure complete, usable for read-only scenarios

- [x] Test with WordPress patterns (2h) ✅ COMPLETE - Comprehensive testing done
  - File: `tests/CLOSURE_TEST_RESULTS.md`, `tests/closures_wordpress_test.php`
  - Test: `add_action('init', function() { echo "Hello"; });` ✅ Works (no params)
  - Test: `array_map(function($x) { return $x * 2; }, $arr);` ❌ Blocked by params
  - Test: Variable capture with `use ($x)` ✅ Works perfectly
  - Test: Reference capture with `use (&$count)` ⚠️ Read works, write doesn't
  - Test: Nested closures ❌ Not tested (blocked by param issue)
  - **Test Results Summary**:
    - ✅ Basic closures (no parameters): PASS
    - ✅ Value capture `use ($var)`: PASS
    - ✅ Multiple captures `use ($a, $b)`: PASS
    - ✅ Reference capture read `use (&$var)`: PASS
    - ❌ Closure parameters: FAIL (critical bug)
    - ❌ Reference capture write: FAIL (known limitation)
    - ❌ Execution flow after closure: FAIL (bug)
  - **Documentation**:
    - Created comprehensive test results: `tests/CLOSURE_TEST_RESULTS.md`
    - 10 test cases with detailed pass/fail analysis
    - WordPress/Laravel use case analysis
    - Identified 3 critical issues with severity and fix estimates
  - **Impact Assessment**:
    - 60-70% of real-world closure patterns blocked by parameter issue
    - Simple hooks and context capture work perfectly
    - Array functions with closures blocked
  - **Status**: Testing complete, blockers documented for Phase 11 work

**Example Code**:
```php
// Simple closure
$greet = function($name) {
    echo "Hello, $name!";
};
$greet("World");

// Closure with capture
$multiplier = 2;
$double = function($x) use ($multiplier) {
    return $x * $multiplier;
};

// Reference capture
$count = 0;
$increment = function() use (&$count) {
    $count++;
};

// WordPress hook
add_action('init', function() {
    // initialization code
});
```

### 10.17 Arrow Functions (6-8h) ✅ IMPLEMENTED - Compiler fixed, same issues as closures

**Priority**: P1 - HIGH
**Depends on**: Closures (10.16)
**Impact**: Modern PHP shorthand, cleaner syntax

- [x] Add ARROW token (`=>`) to lexer (0.5h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/lexer/lexer.go`, `pkg/lexer/token.go`
  - Token `FN` exists (token.go:77)
  - Token `DOUBLE_ARROW` (=>) exists (token.go:197)
  - Both already properly implemented

- [x] Parse arrow function syntax (2h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/expr.go` (lines 1271-1309)
  - Parse `fn($params) => expression` ✅
  - Parse parameter list ✅
  - Parse single expression (not block) ✅
  - Create ArrowFunction AST node ✅
  - **Features**:
    - Reference return `&fn()`
    - Parameter type hints
    - Return type hints with `:`
    - Static arrow functions
  - **Tests**: Comprehensive tests in `pkg/parser/closure_test.go`

- [x] Compile arrow function to closure (3h) ✅ FIXED - Now uses proper operands
  - File: `pkg/compiler/compiler.go` (lines 1127-1243)
  - Compile as closure with automatic capture ✅
  - Capture all variables used in expression ✅
  - **Fixed Implementation**:
    1. Generate unique arrow function name `{arrow}#N` (line 1130)
    2. Emit OpDeclareLambdaFunction with name in Op1 (lines 1155-1159)
    3. Patch start/end positions (lines 1223-1224)
    4. Auto-capture variables with `findReferencedVariables()` (line 1135)
    5. Filter out parameters from capture list (lines 1138-1150)
    6. Emit OpBindLexical for each captured var (lines 1229-1241)
  - **Key Feature**: Automatic variable capture (no `use` clause needed)
  - **Status**: Compiler fixed to match closure implementation
  - **Known Issues**: Same as closures:
    - Parameters not being passed correctly
    - Execution flow stops after declaration
    - Need parameter passing fix (shared with closures)

- [x] Test arrow functions (1h) ✅ TESTED - Same issues as closures
  - File: Tested manually
  - Test: `$fn = fn() => 42;` ✅ Compiles successfully
  - Test: Auto-capture detection works (variables analyzed)
  - Test: With parameters: ❌ Blocked by parameter issue
  - **Status**: Compilation works, runtime blocked by closure issues
  - **Note**: Arrow functions share the same VM execution path as closures
  - Once closure parameters are fixed, arrow functions will work automatically

### 10.18 Documentation & User Guides (20h) 🔄 IN PROGRESS

**Priority**: P2 - MEDIUM

- [x] Complete user guide documentation (8h) ⚠️ PARTIAL - Phase 10 summary created
  - File: `docs/PHASE10_CLOSURE_SUMMARY.md` (comprehensive 450+ line document)
  - **What was created**:
    - Executive summary of Phase 10 closure/arrow function implementation
    - Detailed implementation breakdown for parser, compiler, and runtime
    - Complete test results with 11 test cases analyzed
    - WordPress/Laravel pattern coverage analysis (40% working, 60% blocked)
    - 3 critical issues documented with root cause and fix estimates
    - Performance analysis and recommendations
    - Architecture quality assessment
    - Code locations and file references
  - **What's missing** (for full user guide):
    - Installation guide (needs separate doc)
    - Getting started tutorial (needs separate doc)
    - Configuration reference (needs separate doc)
    - CLI commands documentation (needs separate doc)
  - **Status**: Phase 10 technical summary complete, general user docs pending

- [x] Create migration guide from PHP (6h) ✅ COMPLETE
  - File: `docs/MIGRATION_GUIDE.md` (comprehensive 600+ line document)
  - Compatibility notes ✅
  - Feature comparison ✅ (tables for language features, stdlib, extensions)
  - Known limitations ✅ (8 documented with severity, workarounds, fix timelines)
  - Performance tips ✅ (benchmarking commands, optimization strategies)
  - Migration strategy ✅ (4-phase approach)
  - Common pitfalls ✅ (code examples showing wrong vs right)
  - FAQ section ✅ (general, migration, technical questions)
  - Timeline estimates ✅ (6-8 months for full WordPress/Laravel support)

- [x] Write extension development guide (4h) ✅ ALREADY COMPLETE
  - Directory: `docs/extension-guide/` (5 files, 2,401 total lines)
  - **Main guide** (`README.md`, 509 lines):
    - Go integration basics ✅
    - Creating native extensions ✅
    - Quick start examples ✅
    - Extension interface documentation ✅
    - Best practices (error handling, validation, thread safety, resource cleanup) ✅
    - Complete examples (Math, HTTP, Cache extensions) ✅
  - **FFI guide** (`ffi-guide.md`, 502 lines):
    - Registering Go functions ✅
    - Calling from PHP with `go_call()` ✅
    - FFI helper functions ✅
  - **Type bridging** (`marshaling-guide.md`, 568 lines):
    - PHP ↔ Go type conversions ✅
    - Custom type marshaling ✅
    - Advanced marshaling patterns ✅
  - **Plugin system** (`04-plugin-system.md`, 543 lines):
    - Dynamic plugin loading ✅
    - Building shared libraries ✅
    - Plugin structure and distribution ✅
  - **Quick reference** (`quick-reference.md`, 279 lines):
    - Common patterns ✅
    - Type mapping tables ✅
    - Testing examples ✅
  - **Working examples** in `pkg/goext/bindings/`:
    - HTTP client extension (http.go) ✅
    - JSON extension (json.go) ✅
    - Time extension (time.go) ✅
    - Crypto extension (crypto.go) ✅
    - Filesystem extension (filesystem.go) ✅
  - **Status**: Comprehensive documentation already exists, no action needed

- [x] Performance tuning guide (2h) ✅ ALREADY COMPLETE
  - File: `docs/user-guide/performance-tuning.md` (1,202 lines, comprehensive)
  - Optimization strategies ✅
    - Quick start checklist ✅
    - Value pooling (automatic) ✅
    - Integer caching ✅
    - Resource limits ✅
  - Benchmarking tools ✅
    - Go test benchmarks ✅
    - benchstat comparison ✅
    - Custom benchmark creation ✅
  - Memory profiling ✅
    - Built-in profiling scripts ✅
    - Analysis tools and workflows ✅
    - Memory optimization best practices ✅
  - Production deployment ✅
    - Resource limits configuration ✅
    - Logging configuration (minimal overhead) ✅
    - Metrics collection (near-zero overhead) ✅
    - Health checks (Kubernetes-compatible) ✅
    - Graceful shutdown ✅
    - Kubernetes deployment examples ✅
  - Performance characteristics ✅
    - Micro-benchmarks with real numbers ✅
    - Macro-benchmarks (real-world workloads) ✅
    - Optimization impact analysis ✅
  - Common performance issues ✅
    - High memory usage diagnosis and solutions ✅
    - Slow execution analysis ✅
    - Memory leaks detection ✅
    - High allocation rate optimization ✅
  - Advanced topics ✅
    - Go extension performance tips ✅
    - Memory management strategies ✅
    - Large-scale deployment patterns ✅
  - **Status**: Comprehensive guide already exists, no action needed

---

## Phase 11: WordPress Compatibility (180-220h)

**Goal**: Achieve 95%+ parse success on WordPress 6.8.3 (1255 files)
**Current**: 61.20% parse success (768/1255 files)
**Priority**: HIGH - Real-world framework support

### 11.1 Critical Parser Features - P0 (100-120h)

#### 11.1.1 require/include System (12-16h) ⚠️ MOSTLY COMPLETE (Architectural blocker)

**Priority**: P0 - CRITICAL
**Occurrences**: 1,211 in WordPress (780 require_once, 398 require, 14 include_once, 19 include)

- [x] Parse require/require_once/include/include_once as expressions (2h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/expr.go` (lines 791-828)
  - Tokens: INCLUDE, INCLUDE_ONCE, REQUIRE, REQUIRE_ONCE defined in `pkg/lexer/token.go` ✅
  - AST node: IncludeExpression defined in `pkg/ast/ast.go` (lines 588-597) ✅
  - Parser function: `parseIncludeExpression()` fully implemented ✅
  - Registered in prefix parsers (lines 51-54) ✅
  - Handles both with and without parentheses:
    - `require 'file.php'` ✅
    - `require('file.php')` ✅
  - Supports all 4 variants: include, include_once, require, require_once ✅
  - Parses path as expression (supports concatenation like `require __DIR__ . '/file.php'`) ✅
  - Works as expression value (e.g., `$result = require 'file.php'`) ✅
  - **Tests**: Created comprehensive test suite in `pkg/parser/include_test.go` ✅
    - 4 test functions with 11 test cases total
    - All tests passing
  - **Status**: Parser fully implemented and tested, ready for compiler/runtime work

- [x] Implement file path resolution (3h) ✅ COMPLETE
  - File: `pkg/runtime/include.go` (created, 216 lines)
  - **IncludeManager** implementation:
    - Resolve relative paths ✅ (handles `./` and `../` prefixes)
    - Handle absolute paths ✅ (canonicalized with symlink resolution)
    - Resolve `.` and `..` in paths ✅ (via filepath.Clean and EvalSymlinks)
  - **Path resolution logic**:
    - Absolute paths: returned as-is after canonicalization
    - Relative paths with `./ or ../`: resolved relative to current script directory
    - Other paths: searched in include paths, fallback to script directory
  - **Include path support**:
    - SetIncludePaths/GetIncludePaths for managing include directories ✅
    - AddIncludePath for appending to include path ✅
    - Searches multiple directories in order ✅
  - **Module caching for _once variants**:
    - MarkIncluded/IsIncluded for tracking included files ✅
    - ClearIncluded for resetting cache ✅
    - GetIncludedFiles for listing all included files ✅
    - Thread-safe with sync.RWMutex ✅
  - **Additional features**:
    - SetCurrentScriptDir for context-aware resolution ✅
    - SetBaseDir for base directory fallback ✅
    - Global singleton via GetGlobalIncludeManager() ✅
    - Canonicalization handles symlinks (e.g., /var -> /private/var on macOS) ✅
  - **Tests**: Created comprehensive test suite in `pkg/runtime/include_test.go` ✅
    - 16 test functions covering all functionality
    - Tests for absolute, relative, dotdot paths
    - Tests for include path searching
    - Tests for module caching (_once variants)
    - Concurrent access testing for thread safety
    - All tests passing ✅
  - **Status**: Fully implemented and tested, ready for compiler/VM integration

- [x] Add include path support (2h) ✅ ALREADY IMPLEMENTED (part of previous task)
  - File: `pkg/runtime/include.go` (same as above)
  - Support `include_path` configuration ✅
    - SetIncludePaths() sets full list of paths
    - AddIncludePath() appends single path
    - GetIncludePaths() retrieves current paths
  - Search multiple directories ✅
    - ResolvePath() searches all include paths in order
    - First match wins (PHP-compatible behavior)
  - Match PHP include path behavior ✅
    - Searches include paths for non-relative files
    - Falls back to script directory if not found
    - Absolute and ./ ../ paths skip include path search
  - **Status**: Fully implemented in IncludeManager

- [x] Create module cache for _once variants (3h) ✅ ALREADY IMPLEMENTED (part of previous task)
  - File: `pkg/runtime/include.go` (same as above)
  - Track included files ✅
    - MarkIncluded() marks file as included, returns true if already included
    - IsIncluded() checks if file has been included
  - Prevent duplicate includes for require_once/include_once ✅
    - MarkIncluded() returns true on second call for same file
    - VM/compiler can check before executing file
  - Store by canonicalized path ✅
    - All paths canonicalized via EvalSymlinks before storage
    - Handles symlinks, . and .. correctly
    - Thread-safe with sync.RWMutex
  - **Additional**: ClearIncluded(), GetIncludedFiles() for testing/debugging
  - **Status**: Fully implemented in IncludeManager

- [x] Handle return values from included files (2h) ⚠️ PARTIALLY COMPLETE - Architecture documented
  - **Infrastructure in place**:
    - Parser: IncludeExpression AST node fully functional ✅
    - Compiler: Emits OpIncludeOrEval opcode correctly ✅
    - VM: Has opIncludeOrEval handler stub ✅
    - Runtime: IncludeManager for path resolution and _once tracking ✅
  - **Tests created**: `pkg/runtime/include_integration_test.go` ✅
    - TestIncludeExpressionParsing: Verifies parsing ✅
    - TestIncludeExpressionCompilation: Verifies compilation to opcodes ✅
    - TestIncludeManagerIntegration: Verifies file resolution and _once tracking ✅
    - TestIncludePathResolution: Verifies include path searching ✅
    - All tests passing (4/4) ✅
  - **Architecture blocker documented**:
    - Circular dependency between pkg/vm and pkg/compiler
    - VM needs compiler to compile included files
    - Compiler imports VM for opcodes
    - **Solution options**:
      1. Create separate execution engine package above VM+compiler
      2. Use dependency injection (pass compiler to VM)
      3. Move opcode definitions to separate package
      4. Use interfaces to break dependency
    - **Current implementation**: VM handler returns success (1) as stub
    - **Full implementation needs**: ~4-6h additional work to break circular dependency
  - File: `pkg/vm/handlers_io.go` (lines 41-90, detailed TODO comment)
  - **Status**: Infrastructure complete, execution blocked by architectural constraint

- [x] Add tests with WordPress-like includes (2h) ✅ COMPLETE
  - File: `tests/wordpress_includes_test.go` (470 lines, 7 comprehensive test functions)
  - **Test 1: TestWordPressStyleRequireOnce** ✅
    - Tests require_once with relative paths
    - Mimics WordPress pattern: `require_once __DIR__ . '/wp-load.php'`
    - Verifies path resolution with __DIR__ concatenation
    - Creates temporary file structure and validates resolution
  - **Test 2: TestMultipleIncludesWithOnceVariants** ✅
    - Tests _once variants prevent duplicate includes
    - WordPress pattern to prevent redefinition errors
    - Verifies MarkIncluded() behavior (returns false first time, true on duplicate)
    - Tests IsIncluded() and ClearIncluded() methods
  - **Test 3: TestNestedIncludes** ✅
    - Tests deep include chains (WordPress: index.php → wp-blog-header.php → wp-load.php → wp-config.php)
    - Creates 4-level nested file structure
    - Verifies all files parse and compile correctly
    - Tests include manager tracks all 4 files in chain
  - **Test 4: TestIncludeWithRelativePaths** ✅
    - Tests parent directory references (../)
    - WordPress pattern: wp-admin/admin.php includes ../wp-load.php
    - Creates directory structure: wp-admin/, wp-includes/, root files
    - Verifies resolution of ../wp-load.php and ../wp-includes/functions.php
  - **Test 5: TestReturnValuesFromIncludes** ✅
    - Tests include expressions as values: `$config = require 'config.php'`
    - Valid PHP pattern used in WordPress
    - Verifies parsing and compilation of 3 assignment variations
    - Tests require, include, and require_once as expression values
  - **Test 6: TestIncludePathSearching** ✅
    - Tests multiple include paths (WordPress plugins/themes)
    - Creates wp-content/plugins/ and wp-content/themes/ structure
    - Verifies SetIncludePaths() with multiple directories
    - Tests that files are found in correct order (first path wins)
  - **Test 7: TestConditionalIncludes** ✅
    - Tests WordPress pattern: `if (file_exists(...)) require_once ...`
    - Verifies conditional includes parse and compile
    - Tests `if (!defined('WP_LOADED')) require ...` pattern
  - **All tests passing**: 7/7 (100%) ✅
  - **Status**: Complete WordPress compatibility testing for include infrastructure

#### 11.1.2 Namespace Support (16-20h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Occurrences**: 333 namespace + 210 use in WordPress

- [x] Parse namespace declarations (2h) ✅
  - File: `pkg/parser/stmt.go` (lines 933-1009)
  - Already implemented: parseNamespaceStatement(), parseNamespaceName()
  - Test file: `pkg/parser/namespace_test.go` (420+ lines, 10 test functions)
  - Supports all syntax variants: simple, nested, bracketed, global, unbracketed
  - All tests passing (10/10)

- [x] Parse use statements (3h) ✅
  - File: `pkg/parser/stmt.go` (lines 1013-1131)
  - Already implemented: parseUseStatement()
  - Test file: `pkg/parser/use_test.go` (501 lines, 10 test functions)
  - Supports all use variants:
    - Simple use: `use Namespace\Class;`
    - Aliased use: `use Namespace\Class as Alias;`
    - Multiple use: `use A, B, C;`
    - Grouped use: `use Namespace\{ClassA, ClassB};`
    - Function use: `use function Namespace\func;`
    - Const use: `use const Namespace\CONST;`
  - All tests passing (10/10 test functions, includes WordPress-style patterns)

- [x] Implement namespace resolution in compiler (4h) ✅
  - File: `pkg/compiler/compiler.go`, `pkg/compiler/symbols.go`
  - Added currentNamespace tracking in Compiler struct
  - Added useImports, useFunctionImports, useConstImports maps
  - Implemented ResolveClassName(), ResolveFunctionName(), ResolveConstantName()
  - Test file: `pkg/compiler/namespace_test.go` (548 lines, 10 test functions)
  - All tests passing (35+ subtests covering all resolution scenarios)

- [x] Add FQN (Fully Qualified Name) support (3h) ✅
  - Implemented in ResolveClassName() in `pkg/compiler/symbols.go`
  - Handles relative class names (prepends current namespace)
  - Handles qualified names with first part in use imports
  - Resolves use aliases automatically

- [x] Handle global namespace access (2h) ✅
  - Implemented in ResolveClassName(), ResolveFunctionName(), ResolveConstantName()
  - Names starting with `\` are treated as fully qualified
  - Leading backslash is stripped to get the global FQN

- [x] Update ::class to work with namespaces (2h) ✅
  - Added ClassNameExpression compilation in `pkg/compiler/compiler.go` (lines 1565-1612)
  - Resolves class names to FQN using ResolveClassName() at compile time
  - Handles special keywords (self, parent, static) - passed to VM for runtime resolution
  - Test file updated: `pkg/compiler/namespace_test.go` (717 lines, 12 test functions)
  - Added 2 new test functions: TestClassNameExpressionWithNamespace (6 subtests),
    TestClassNameExpressionSpecialKeywords (3 subtests)
  - All tests passing

- [x] Add comprehensive tests (4h) ✅
  - Test file: `pkg/compiler/namespace_test.go` (717 lines, 12 test functions)
  - Covered: Basic namespace declaration, nested namespaces, use statements with aliases,
    grouped use statements, class resolution in namespaces, function/constant namespaces
  - WordPress-style patterns tested
  - All tests passing (40+ subtests)

#### 11.1.3 global Keyword (4-6h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Occurrences**: 1,061 in WordPress

- [x] Parse global statement (1h) ✅
  - File: `pkg/parser/stmt.go` (lines 97-130)
  - Already implemented: parseGlobalStatement()
  - Compiler: `pkg/compiler/compiler.go` (lines 394-407) emits OpBindGlobal
  - Tests passing: TestGlobalStatement (3 subtests), TestCompileGlobalStatement (2 subtests)

- [x] Implement global variable importing in VM (3h) ✅
  - File: `pkg/compiler/compiler.go` (lines 394-411) - passes variable name as constant
  - File: `pkg/vm/handlers_variables.go` (lines 167-212) - opBindGlobal implementation
  - File: `pkg/vm/frame.go` (line 29) - added globalBindings map
  - File: `pkg/vm/vm.go` (lines 561-629) - modified getOperandValue/setOperandValue for global bindings
  - Implementation details:
    - OpBindGlobal takes CV index (Op1) and variable name as constant (Op2)
    - Frame tracks global bindings: local index -> global variable name
    - Get/set operations check for global bindings and redirect to vm.globals
  - All existing tests passing

- [x] Add tests (1h) ✅
  - Test file: `pkg/vm/vm_test.go` (lines 970-1250)
  - Added 6 test functions:
    - TestGlobalVariableBinding - verifies binding is recorded in frame
    - TestGlobalVariableRead - reads global value through binding
    - TestGlobalVariableWrite - writes to global through binding
    - TestMultipleGlobalBindings - tests 3 globals simultaneously
    - TestGlobalVariableNewCreation - creates non-existent global
    - TestGlobalVariableUnboundLocalAccess - verifies unbound locals work normally
  - All tests passing (6/6)

#### 11.1.4 unset() Construct (3-4h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Occurrences**: 1,042 in WordPress

- [x] Parse unset() with multiple arguments (1h) ✅
  - File: `pkg/parser/stmt.go` (lines 134-175)
  - Already implemented: parseUnsetStatement()
  - Supports: single variable, multiple variables, array elements, object properties, nested
  - Tests passing: TestUnsetStatement (7 subtests)

- [x] Implement OpUnset for variables, array elements, properties (2h) ✅
  - Already implemented:
    - Compiler: `pkg/compiler/compiler.go` (lines 500-568) handles UnsetStatement
    - VM handlers:
      - `opUnsetVar` in handlers_variables.go (lines 214-231) - sets variable to Undef
      - `opUnsetDim` in handlers_array.go (lines 306-331) - removes array element
      - `opUnsetObj` in handlers_object.go (lines 313-348) - deletes object property
    - Opcodes: OpUnsetVar (74), OpUnsetDim (75), OpUnsetObj (76)
  - Tests passing: TestUnsetStatement (3 subtests), TestOpUnsetObj (2 subtests)

- [x] Add tests (1h) ✅
  - Test: unset variable - TestOpUnsetVar (4 subtests: basic, string, array, already unset)
  - Test: unset array element - TestOpUnsetDim (4 subtests: string key, int key, non-existent, non-array)
  - Test: unset object property - TestOpUnsetObj, TestOpUnsetObj_NonObject
  - Test: Multiple unset targets - TestMultipleUnsets
  - All tests in pkg/vm/vm_test.go (11 new tests total)

#### 11.1.5 Visibility Modifiers (4-6h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Occurrences**: 1,825 in WordPress (1,406 public, 339 protected, 80 private)

- [x] Parse visibility modifiers in method declarations (1h) ✅
  - File: `pkg/parser/decl.go` (lines 240-322)
  - Already implemented: parseClassMember() collects visibility modifiers
  - parseMethodDeclaration() receives and stores visibility
  - AST: MethodDeclaration.Visibility field

- [x] Parse visibility modifiers in property declarations (1h) ✅
  - File: `pkg/parser/decl.go` (lines 372-445)
  - Already implemented: parsePropertyDeclaration() receives and stores visibility
  - AST: PropertyDeclaration.Visibility field

- [x] Implement runtime visibility checks (2h) ✅
  - File: `pkg/vm/handlers_object.go`
  - Property access: Updated all 11 accessContext usages to use frame.currentClass
  - Method calls: Added canAccessMethod() helper at line 988
  - opInitMethodCall: Added visibility check at line 648
  - opInitStaticMethodCall: Added visibility check at line 732
  - Throws error for invalid access with message "cannot access {visibility} method/property"

- [x] Add tests (1h) ✅
  - Test: public/protected/private properties - 5 tests (TestPropertyVisibility_*)
  - Test: public/protected/private methods - 5 tests (TestMethodVisibility_*)
  - Test: Inheritance visibility - included in protected tests
  - Test: Visibility violations - 4 tests (block private from outside, protected from unrelated)
  - Test: Static method visibility - TestStaticMethodVisibility_PrivateFromOutside
  - All 11 tests in pkg/vm/handlers_object_test.go

#### 11.1.6 isset() and empty() Enhancements (6-8h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Occurrences**: 5,565 isset + 5,149 empty in WordPress

- [x] Enhance isset() for all cases (3-4h) ✅
  - File: `pkg/compiler/compiler.go` (lines 1632-1754)
  - Support `isset($var)` with undefined variables - uses OpIssetIsemptyVar
  - Support `isset($arr[$key])` for array access - uses OpIssetIsemptyDimObj
  - Support `isset($obj->prop)` for properties - uses OpIssetIsemptyPropObj
  - Support multiple arguments: `isset($a, $b, $c)` - short-circuit evaluation
  - Added compileIssetCheck() helper function for type-specific handling

- [x] Enhance empty() for PHP truthiness (3-4h) ✅
  - File: `pkg/compiler/compiler.go` (lines 1756-1835)
  - Support `empty($var)` - uses OpIssetIsemptyVar with mode=1
  - Support `empty($arr[$key])` - uses OpIssetIsemptyDimObj
  - Support `empty($obj->prop)` - uses OpIssetIsemptyPropObj
  - Truthiness rules already implemented in IsFalse() in pkg/types/value.go
  - Edge cases ("0", 0, 0.0, [], null, false) handled by IsFalse()

- [x] Add comprehensive tests (1h) ✅
  - Test file: `pkg/compiler/compiler_test.go`
  - Test: isset with undefined variables (4 tests) ✅
  - Test: isset with array access (5 tests) ✅
  - Test: isset with multiple arguments (4 tests) ✅
  - Test: empty with all falsy values (9 tests) ✅
  - Test: empty with array access (4 tests) ✅
  - Test: isset/empty in conditionals (5 tests) ✅
  - Total: 31 new tests for isset/empty functionality
  - **Bug fixes during testing**:
    - Fixed `OpInitFcallByName` to use ExtendedValue for arg count (was using CONST operand)
    - Fixed `OpIssetIsemptyVar` to use ExtendedValue for mode (0=isset, 1=empty)
    - Fixed `OpIssetIsemptyDimObj` and `OpIssetIsemptyPropObj` to use ExtendedValue for mode
    - Fixed multi-argument isset() short-circuit evaluation with proper jump targets

#### 11.1.7 list() Array Destructuring (6-8h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Occurrences**: 265 in WordPress
**Note**: `[$a, $b] = ...` already works in foreach, need assignment context

- [x] Parse list() as assignment target (2h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/expr.go` (parseListExpression)
  - Parse `list($a, $b) = expr;` ✅
  - ListExpression AST node exists ✅

- [x] Compile to same opcodes as [] destructuring (2h) ✅ COMPLETE
  - File: `pkg/compiler/compiler.go` (lines 1079-1087, 3163-3363)
  - Added `compileListDestructuring()` helper function
  - Added `compileArrayDestructuring()` helper function
  - Emit FETCH_DIM_R for each element
  - Support for keyed destructuring: `list("x" => $a) = $arr`

- [x] Support nested list() (2h) ✅ COMPLETE
  - Parse `list($a, list($b, $c)) = $arr;` ✅
  - Compile nested destructuring ✅
  - Mixed nesting: `list($a, [$b, $c]) = $arr` ✅

- [x] Support list() in foreach (1h) ✅ COMPLETE
  - File: `pkg/compiler/compiler.go` (lines 2400-2422)
  - Added support for `foreach ($arr as list($a, $b))`
  - Added support for `foreach ($arr as [$a, $b])`

- [x] Add tests (1h) ✅ COMPLETE
  - File: `pkg/compiler/regression_test.go` (TestRegression_Phase10_ListDestructuring)
  - Test: Basic list assignment ✅
  - Test: Nested list assignment ✅
  - Test: list with array keys ✅
  - Test: list in foreach (verify compatibility) ✅
  - 12 comprehensive test cases all passing

### 11.2 High Priority Features - P1 (40-50h)

#### 11.2.1 Class Constants (6-8h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Occurrences**: 156 in WordPress

- [x] Parse const declarations in class body (2h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/decl.go` (parseClassConstant)
  - Parse `const NAME = value;` ✅
  - ClassConstantDeclaration AST node exists ✅

- [x] Store constants in ClassEntry (1h) ✅ COMPLETE
  - File: `pkg/types/object.go` - ClassConstant struct already exists
  - File: `pkg/compiler/compiler.go` (lines 2995-3014)
  - Added constant evaluation during class compilation
  - Implemented `evaluateConstantExpression()` for compile-time constant evaluation
  - Supports: int, float, string, bool, null, arrays, arithmetic, string concat

- [x] Implement constant access via :: (2h) ✅ COMPLETE
  - File: `pkg/compiler/compiler.go` (lines 1530-1542)
  - Compile `Class::CONSTANT` via OpFetchClassConstant
  - File: `pkg/vm/vm.go` (line 400) - Added dispatch case
  - File: `pkg/vm/handlers_object.go` (lines 1004-1058) - Implemented opFetchClassConstant

- [x] Support const visibility modifiers (1h) ✅ COMPLETE
  - Parse public/protected/private const ✅
  - Visibility checking in opFetchClassConstant ✅

- [x] Add tests (1h) ✅ COMPLETE
  - File: `pkg/compiler/regression_test.go` (TestRegression_Phase10_ClassConstants)
  - 12 comprehensive test cases all passing
  - Test: Class constants ✅
  - Test: Constant access ✅
  - Test: Multiple constants ✅
  - Test: Constant expressions ✅
  - Test: Visibility modifiers ✅

#### 11.2.2 Switch/Case Statements (8-10h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Occurrences**: 14 case + 29 default in WordPress

- [x] Complete switch statement parsing (2h)
  - File: `pkg/parser/stmt.go`
  - Verify switch/case/default parsing
  - Note: Parser already implemented (lines 526-635)
  - Supports both regular `switch() {}` and alternative `switch(): endswitch;` syntax

- [x] Implement OpSwitch and case comparison (4h)
  - File: `pkg/compiler/compiler.go`
  - Uses OpIsEqual for loose comparison and OpJmpNZ for conditional jumps
  - Fixed bug: OpJmpNZ was patching Op1 (condition) instead of Op2 (jump target)
  - Proper temp variable management with saved/restored stack levels

- [x] Handle fall-through behavior (2h)
  - PHP fall-through semantics implemented correctly
  - Break statements properly compile to jumps to switch end
  - Cases without break fall through to next case

- [x] Add tests (2h)
  - Test: Basic switch statement ✅
  - Test: Fall-through cases ✅
  - Test: default case ✅
  - Test: String case values ✅
  - Test: Alternative syntax ✅
  - Test: Expression subject ✅
  - Added 9 regression tests in `pkg/compiler/regression_test.go`

#### 11.2.3 Try/Catch Exception Handling (12-16h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Occurrences**: 1 try + 3 catch in WordPress core (more in plugins)
**Status**: Completed in Phase 10

- [x] Complete try/catch parsing (2h) ✅
  - File: `pkg/parser/stmt.go` (lines 723-793)
  - Parser supports try/catch/finally with multiple catch clauses ✅

- [x] Implement exception handling in VM (6h) ✅
  - File: `pkg/vm/vm.go`, `pkg/vm/handlers_exception.go`
  - Registered built-in exception classes (Exception, Error, RuntimeException, etc.) ✅
  - Implemented OpThrow, OpCatch handlers ✅
  - Implemented exception unwinding with findCatchHandler() ✅
  - Match exception type with inheritance checking ✅
  - Native constructor for Exception class ✅

- [x] Support multiple catch blocks (2h) ✅
  - Compiler stores catch type in OpCatch instruction ✅
  - VM scans catch blocks for matching exception type ✅
  - Supports exception inheritance (RuntimeException caught by Exception) ✅

- [x] Implement finally block (3h) ✅ BASIC
  - Simplified implementation: finally executes after try or catch ✅
  - Compiler patches JMPs to finally block ✅
  - **Note**: Finally with uncaught exception propagation not yet supported
  - **Note**: FAST_CALL/FAST_RET opcodes implemented but not used

- [x] Add tests (2h) ✅
  - Test: Basic try/catch ✅
  - Test: Multiple catch blocks ✅
  - Test: Exception message access ✅
  - Test: Exception inheritance ✅
  - Test: Simple finally block ✅
  - Test: Uncaught exception error ✅
  - Unit tests in `pkg/vm/exception_test.go` ✅

**Known Limitations**:
- Finally block only runs when try/catch completes normally; uncaught exceptions skip finally
- User-defined exception subclasses need `parent::__construct()` for message/code (not yet supported)
- Union types in catch (`catch (A|B $e)`) not yet supported

### 11.3 Medium Priority Features - P2 (20-30h)

#### 11.3.1 clone Keyword (4-6h) ✅ COMPLETE

**Priority**: P2 - MEDIUM
**Occurrences**: 80 in WordPress

- [x] Parse clone expression (1h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/expr.go`
  - Parse `clone $object`
  - Implemented in `parseCloneExpression()` at line 547

- [x] Implement OpClone (2h) ✅ COMPLETE
  - File: `pkg/compiler/compiler.go`, `pkg/vm/handlers_object.go`
  - Add OpClone instruction (opcode 110)
  - Shallow copy object with proper property copying
  - Fixed compiler temp variable allocation issue

- [x] Support __clone() magic method (1h) ✅ COMPLETE
  - Call __clone() after cloning on the cloned object
  - Implemented `callMagicClone()` helper in handlers_object.go
  - Checks both MagicMethods and Methods maps for __clone

- [x] Add tests (1h) ✅ COMPLETE
  - Test: `TestOpClone_BasicCloning` - Basic clone
  - Test: `TestOpClone_NonObject` - Error handling
  - Test: `TestOpClone_MagicCloneMethod` - Clone with __clone()

#### 11.3.2 var Property Declaration (2-3h) ✅ COMPLETE

**Priority**: P2 - MEDIUM
**Occurrences**: 63 in WordPress

- [x] Parse var as public property (1h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/decl.go` (lines 296-299)
  - Treat var as public
  - Parser already handles `var` keyword correctly

- [x] Add tests (1h) ✅ COMPLETE
  - Test: `TestPropertyDeclaration` - Added var keyword test cases
  - Test: `TestVarPropertyDeclaration` - Comprehensive var property tests
  - Test: `TestVarMultipleProperties` - Multiple var properties on one line

**Note**: Property default value initialization is a separate VM issue (not related to var parsing)

#### 11.3.3 final Modifier (3-4h) ✅ COMPLETE

**Priority**: P2 - MEDIUM
**Occurrences**: 23 in WordPress

- [x] Parse final on classes and methods (1h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/parser/parser.go` (lines 143-164)
  - `parseClassDeclarationWithModifiers()` handles `final` class modifier
  - File: `pkg/parser/decl.go` - handles `final` method modifier

- [x] Set IsFinal on ClassEntry and resolve parent at runtime
  - File: `pkg/compiler/compiler.go` (lines 2994-3002) - Set `classEntry.IsFinal`
  - File: `pkg/types/object.go` - Added `ParentClassName` field for deferred resolution
  - File: `pkg/compiler/compiler.go` (line 3006) - Set `ParentClassName` from `node.Extends`

- [x] Enforce final in inheritance (2h) ✅ COMPLETE
  - File: `pkg/vm/handlers_object.go` (opDeclareClass, lines 1095-1106)
  - Resolve parent class at runtime and call `InheritFrom()`
  - File: `pkg/types/object.go` (InheritFrom, lines 508-511) - Prevent extending final classes
  - File: `pkg/types/object.go` (validateMethodOverride, lines 636-638) - Prevent overriding final methods

- [x] Add tests (1h) ✅ COMPLETE
  - Test: `TestFinalClass` (parser/decl_test.go) - Parsing final class
  - Test: `TestOpDeclareClass_FinalClassCannotBeExtended` (vm/handlers_inheritance_test.go)
  - Test: `TestOpDeclareClass_FinalMethodCannotBeOverridden` (vm/handlers_inheritance_test.go)
  - Test: `TestOpDeclareClass_ValidInheritance` (vm/handlers_inheritance_test.go)

### 11.4 HTML/PHP Mixed Mode (12-16h) ✅ COMPLETE

**Priority**: P0 - CRITICAL
**Impact**: All WordPress templates fail without this

- [x] Add HTML mode to lexer (4h) ✅ COMPLETE
  - File: `pkg/lexer/lexer.go`
  - HTML lexing mode already implemented (`inPHP` field, `scanInlineHTML()` function)
  - Emits `INLINE_HTML` tokens for HTML content

- [x] Switch between HTML and PHP mode on tags (4h) ✅ COMPLETE
  - `scanPHPTag()` handles `<?php`, `<?=`, `<?` tags (sets `inPHP = true`)
  - `?>` close tag sets `inPHP = false`
  - Multiple PHP blocks work correctly

- [x] Handle `<?=` short echo tags (2h) ✅ COMPLETE
  - File: `pkg/parser/parser.go` (ParseProgram, lines 79-105)
  - Parse `<?= expr ?>` as echo statement
  - Supports multiple expressions separated by commas

- [x] Fix PHP newline suppression after ?> (1h) ✅ COMPLETE
  - File: `pkg/lexer/lexer.go` (scanInlineHTML, lines 826-838)
  - Single newline immediately following `?>` is consumed (PHP behavior)
  - Handles both Unix (`\n`) and Windows (`\r\n`) line endings

- [x] Test with WordPress templates (4h) ✅ COMPLETE
  - Test: HTML with embedded PHP - PASS
  - Test: Multiple PHP blocks - PASS
  - Test: Short echo tags (`<?= ?>`) - PASS
  - Test: WordPress-style templates - PASS

### 11.5 WordPress Standard Library (60-80h) ⬜ NOT STARTED

#### String Functions (20-25h)

- [x] HTML Encoding (4-5h) ✅ COMPLETE
  - `htmlspecialchars` - Full implementation with ENT_QUOTES/ENT_COMPAT/ENT_NOQUOTES flag support
  - `htmlentities` - Full implementation with 100+ named HTML entities
  - `html_entity_decode` - Decodes named entities, decimal (&#60;) and hex (&#x3C;) numeric entities
  - `htmlspecialchars_decode` - Reverses htmlspecialchars with flag support
  - `strip_tags` - Removes HTML/PHP tags, supports allowed tags parameter
  - `addslashes` - Escapes single/double quotes, backslashes, NUL bytes
  - `stripslashes` - Reverses addslashes

- [x] String Manipulation (8-10h) ✓ COMPLETE
  - `str_repeat`, `str_pad`, `str_split`, `str_word_count`, `str_shuffle` ✓
  - `substr_count`, `substr_replace`, `strrev` ✓
  - `strcmp`, `strcasecmp`, `strncmp`, `strncasecmp` ✓
  - `chunk_split`, `wordwrap`, `nl2br` ✓

- [x] String Search (4-5h) ✓ COMPLETE
  - `strstr`, `stristr`, `strrchr`, `strpbrk` ✓
  - `strspn`, `strcspn` ✓
  - Also registered: `strchr` (alias of strstr) ✓

- [x] Formatting (4-5h) ✓ COMPLETE
  - `sprintf`, `vsprintf`, `printf`, `vprintf` ✓
  - `number_format`, `sscanf`, `str_getcsv` ✓

#### Array Functions (15-20h)

- [x] Array Manipulation (8-10h) ✅ COMPLETE
  - `array_pad`, `array_sum`, `array_product` ✅
  - `array_column`, `array_change_key_case` ✅
  - `array_replace`, `array_replace_recursive` ✅
  - Also added: `array_intersect`, `current`, `key`, `reset`, `end`, `next`, `prev`

- [x] Array Search/Filter (4-5h) ✅ COMPLETE
  - `array_key_exists`, `key_exists` ✅
  - `array_diff_key`, `array_intersect_key` ✅
  - `array_diff_assoc`, `array_intersect_assoc` ✅
  - `array_key_first`, `array_key_last`, `array_count_values` ✅

- [x] Array Iteration (3-5h) ✅ COMPLETE
  - `array_walk_recursive` ✅ (stub - requires callable support)
  - `current`, `next`, `prev`, `reset`, `end`, `key` ✅ (already implemented)

#### File System Functions (12-16h)

- [x] File Operations (6-8h) ✅ COMPLETE
  - `is_dir`, `is_file`, `is_readable`, `is_writable` ✅
  - `mkdir`, `rmdir`, `unlink`, `rename`, `copy` ✅
  - `filesize`, `filemtime`, `touch`, `chmod` ✅
  - Also added: `file_put_contents`, `file`, `readfile`, `filetype` ✅
  - Also added: `fopen`, `fclose`, `fread`, `fwrite`, `fgets`, `fgetc` ✅

- [x] Directory Operations (4-6h) ✅ COMPLETE
  - `scandir`, `glob` ✅ (already implemented)
  - Note: `opendir`, `readdir`, `closedir` - low priority (scandir covers most use cases)
  - Note: `fnmatch` - deferred (use glob instead)

- [x] Path Operations (2-3h) ✅ COMPLETE
  - `basename`, `dirname`, `pathinfo`, `realpath` ✅

#### Other Functions (13-19h)

- [x] URL Functions (3-4h) ✅ COMPLETE
  - `parse_url`, `http_build_query` ✅
  - `urlencode`, `urldecode`, `rawurlencode`, `rawurldecode` ✅
  - `base64_encode`, `base64_decode` ✅

- [x] Type Functions (2-3h) ✅ COMPLETE
  - `function_exists`, `class_exists`, `method_exists` ✅
  - `serialize`, `unserialize` ✅ (full PHP format implementation)
  - `json_encode`, `json_decode`, `json_last_error`, `json_last_error_msg` ✅

- [x] Math Functions (2-3h) ✅ COMPLETE
  - `abs`, `ceil`, `floor`, `round`, `min`, `max` ✅
  - `pow`, `sqrt`, `log`, `log10`, `exp` ✅
  - `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2` ✅
  - `deg2rad`, `rad2deg`, `pi` ✅
  - `rand`, `mt_rand`, `random_int`, `getrandmax`, `mt_getrandmax` ✅
  - `is_nan`, `is_infinite`, `is_finite`, `hypot`, `fmod`, `intdiv`, `fdiv` ✅
  - 37 total math functions registered in builtins.go ✅

- [x] Misc Functions (3-4h) ✅ COMPLETE
  - `define`, `defined`, `constant` ✅ (with built-in constants: PHP_VERSION, PHP_EOL, etc.)
  - `call_user_func`, `call_user_func_array` ✅ (supports string callbacks and closures)
  - `func_get_args`, `func_num_args` ✅
  - `function_exists` ✅ (improved to check both builtins and user-defined functions)
  - Added userConstants map to VM with built-in PHP constants ✅
  - 10 new tests added ✅

---

## Phase 12: Laravel Compatibility (120-160h)

**Goal**: Achieve 95%+ parse success on Laravel v12.10.1 (7,483 files)
**Current**: 43.57% parse success (3,260/7,483 files)
**Priority**: HIGH - Modern framework support

### 12.1 Laravel-Specific Features (50-70h)

#### 12.1.1 Array Destructuring in Assignment (2-3h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Note**: Already fully implemented - works in both foreach and assignment statements

- [x] Support `[$a, $b] = expr` at statement level (1h)
  - File: `pkg/compiler/compiler.go` - compileArrayDestructuring at line 3364
  - Implemented: parseArrayExpression parses `[...]`, assignment infix handles it

- [x] Support nested destructuring (1h)
  - Parse `[[$a, $b], $c] = $arr;`
  - Implemented: Recursive handling in compileArrayDestructuring

- [x] Add tests (0.5h)
  - Test: 13 tests in TestRegression_Phase10_ListDestructuring
  - Test: 8 tests in TestForeachArrayDestructuring
  - Test: 4 tests in TestRegression_Phase6B_ForeachArrayDestructuring

#### 12.1.2 Reference Parameters (8-10h) ✅ COMPLETE

**Priority**: P1 - HIGH
**Occurrences**: 50+ in Laravel

- [x] Parse & in parameter declarations (1h)
  - File: `pkg/parser/stmt.go`
  - Parse `function foo(&$param)`
  - Already implemented: `param.ByRef` field in AST

- [x] Parse & in function call arguments (1h)
  - Parse `foo(&$var)`
  - Already implemented in parser

- [x] Implement reference passing in compiler (3h)
  - File: `pkg/compiler/compiler.go`
  - Added `FunctionSignature` struct to track which params are by-ref
  - Added `functionSignatures` map to Compiler
  - Emit OpSendRef instead of OpSendVal for by-ref params with variable operand

- [x] Implement reference passing in VM (3h)
  - File: `pkg/vm/handlers_functions.go`
  - Added `opSendRef` handler with reference creation
  - File: `pkg/vm/frame.go`
  - Modified `setLocal` to update through references
  - Modified `getLocal` to auto-dereference
  - File: `pkg/types/value.go`
  - Added `SetReferenceTarget` and `GetReferenceTarget` methods

- [x] Add tests (1h)
  - Test: 6 tests in TestRegression_Phase12_ReferenceParameters
  - Test: basic_reference_parameter, multiple_calls_with_reference
  - Test: reference_parameter_string_modification, reference_parameter_with_array_assignment
  - Test: nested_function_calls_with_reference, reference_parameter_set_to_null

#### 12.1.3 Attributes #[...] (12-16h) 🔄 IN PROGRESS (Parsing complete, Reflection pending)

**Priority**: P2 - MEDIUM (Complex, can defer)
**Occurrences**: 239 in Laravel

- [x] Add # token to lexer (0.5h) ✅ ALREADY IMPLEMENTED
  - File: `pkg/lexer/lexer.go` (line 388-395)
  - Token `ATTRIBUTE_START` defined in `pkg/lexer/token.go` (line 224)
  - Lexer recognizes `#[` and returns `ATTRIBUTE_START` token
  - Tested: `#[Attribute]` correctly tokenizes

- [x] Parse attribute syntax (3h) ✅ COMPLETE
  - File: `pkg/parser/stmt.go` (parseAttributeGroups, parseAttribute, parseAttributeArguments, parseNamespacedName)
  - File: `pkg/parser/parser.go` (added ATTRIBUTE_START case to parseStatement)
  - File: `pkg/parser/decl.go` (added attribute support to parseClassMember)
  - Parse `#[AttributeName]` ✅
  - Parse `#[Attr(arg1, arg2)]` with positional arguments ✅
  - Parse `#[Attr(name: value)]` with named arguments ✅
  - Parse `#[Attr1, Attr2]` multiple attributes in group ✅
  - Parse namespaced attributes `#[\Foo\Bar\Attr]` ✅
  - Support attributes on classes, functions, methods, properties ✅
  - Test file: `pkg/parser/attribute_test.go` (12 tests, all passing)

- [x] Store attributes in AST nodes (2h) ✅ ALREADY IMPLEMENTED
  - Attributes field already exists in ClassDeclaration, FunctionDeclaration,
    MethodDeclaration, PropertyDeclaration, Parameter AST nodes
  - AttributeGroup and Attribute types defined in `pkg/ast/ast.go`

- [ ] Implement attribute access via reflection (4h)
  - File: `pkg/stdlib/reflection/` (new)
  - Store attributes in metadata
  - Create reflection API

- [x] Support attribute arguments (named args) (2h) ✅ COMPLETE
  - Parse attribute with named arguments ✅ (parseAttributeArguments)
  - Store argument values ✅ (attr.Named map[string]Expr)
  - Tested with `#[Route(path: "/api", methods: ["GET"])]`

- [x] Add tests (2h) ✅ COMPLETE
  - Test file: `pkg/parser/attribute_test.go`
  - Test: Basic attributes ✅ (TestAttributeParsing/simple_attribute_on_class)
  - Test: Attribute arguments ✅ (multiple subtests)
  - Test: Multiple attributes ✅ (TestAttributeParsing/multiple_attributes_on_class)
  - Test: Named arguments ✅ (TestAttributeParsing/attribute_with_named_argument)
  - Test: Attributes on properties ✅ (TestAttributeOnProperty)
  - Test: Attributes on methods ✅ (TestAttributeOnMethod)
  - Test: Attributes on functions ✅ (TestAttributeOnFunction)
  - Test: Namespaced attributes ✅ (TestNamespacedAttribute)
  - Total: 12 tests, all passing

#### 12.1.4 Anonymous Classes (8-12h) 🔄 IN PROGRESS (Parsing complete)

**Priority**: P2 - MEDIUM
**Impact**: Used in tests and mocks

- [x] Parse `new class` syntax (2h) ✅ COMPLETE
  - File: `pkg/parser/expr.go` (parseNewExpression, parseAnonymousClass)
  - File: `pkg/ast/ast.go` (AnonymousClassExpression struct)
  - Parse `new class { ... }` ✅
  - Parse `new class(args) { ... }` with constructor args ✅
  - Parse `new class extends Base { ... }` ✅
  - Parse `new class implements Iface { ... }` ✅
  - Parse `new class implements Iface1, Iface2 { ... }` multiple interfaces ✅
  - Parse properties, methods in anonymous class body ✅
  - Test file: `pkg/parser/anonymous_class_test.go` (10 tests, all passing)

- [ ] Create anonymous class compilation (4h)
  - File: `pkg/compiler/compiler.go`
  - Generate unique class name
  - Compile class definition
  - Instantiate immediately

- [ ] Generate unique class names (1h)
  - Use counter or hash for uniqueness

- [x] Support inheritance and interfaces (2h) ✅ COMPLETE (Parsing)
  - Parse `new class extends Base implements Iface` ✅
  - Compile inheritance - pending

- [x] Add tests (2h) ✅ COMPLETE
  - Test: Basic anonymous class ✅
  - Test: With constructor ✅
  - Test: With inheritance ✅
  - Test: With interfaces ✅
  - Test: As function argument ✅
  - Test: With properties ✅

### 12.2 Laravel Standard Library (40-60h)

#### Reflection API (15-20h) ⬜ NOT STARTED

- [ ] Class Reflection (8-10h)
  - `ReflectionClass`, `ReflectionMethod`, `ReflectionProperty`
  - Methods: getName, getAttributes, isPublic, etc.

- [ ] Attribute Reflection (4-6h)
  - `ReflectionAttribute`
  - Methods: getName, getArguments, newInstance

- [ ] Advanced Reflection (3-4h)
  - `ReflectionParameter`, `ReflectionType`
  - `get_class_methods`, `get_class_vars`

#### SPL (Standard PHP Library) (12-16h) ⬜ NOT STARTED

- [ ] Iterators (6-8h)
  - `Iterator`, `IteratorAggregate`
  - `ArrayIterator`, `DirectoryIterator`

- [ ] Data Structures (4-6h)
  - `SplFixedArray`, `SplQueue`, `SplStack`

- [ ] Exceptions (2-3h)
  - `LogicException`, `RuntimeException`
  - `InvalidArgumentException`

#### JSON Functions (2-3h) ⬜ NOT STARTED

- [ ] JSON Encoding/Decoding (2-3h)
  - `json_encode`, `json_decode` with all options
  - `json_last_error`, `json_last_error_msg`
  - Support JSON_THROW_ON_ERROR, JSON_PRETTY_PRINT

#### DateTime Classes (5-7h) ⬜ NOT STARTED

- [ ] DateTime Class (5-7h)
  - `DateTime`, `DateTimeImmutable`
  - `DateInterval`, `DatePeriod`, `DateTimeZone`
  - Methods: format, modify, add, sub, diff

### 12.3 Service Container Support (10-15h)

#### Type Resolution for DI (6-8h) ⬜ NOT STARTED

- [ ] Extract parameter types from function signatures (2h)
  - File: `pkg/compiler/compiler.go`
  - Store type information

- [ ] Store type information in CompiledFunction (2h)
  - Add ParameterTypes field

- [ ] Provide runtime type access (2h)
  - Reflection API for parameter types

- [ ] Add tests (1h)

#### Reflection Enhancements for DI (4-7h) ⬜ NOT STARTED

- [ ] Expose constructor parameters with types (2h)
  - ReflectionClass::getConstructor()
  - ReflectionMethod::getParameters()

- [ ] Support automatic dependency resolution (3h)
  - Resolve dependencies by type
  - Create instances with dependencies

- [ ] Add tests (1h)

---

## Phase 13: Standard Library Expansion (200-300h)

**Goal**: Complete ~300 remaining PHP standard library functions
**Priority**: MEDIUM - Needed for full WordPress/Laravel support

### 13.1 String Functions (40-60h) ⬜ NOT STARTED

- [ ] Multi-byte functions (mb_*) (20-30h)
- [ ] Encoding functions (10-15h)
- [ ] Locale-aware comparisons (5-8h)
- [ ] Advanced parsing (5-7h)

### 13.2 Array Functions (30-50h) ⬜ NOT STARTED

- [ ] Advanced sorting (10-15h)
- [ ] Set operations (10-15h)
- [ ] Array cursors (5-10h)
- [ ] Custom comparators (5-10h)

### 13.3 File System Functions (30-40h) ⬜ NOT STARTED

- [ ] Stream functions (10-15h)
- [ ] Directory traversal (10-15h)
- [ ] File uploads (5-10h)

### 13.4 Network Functions (20-30h) ⬜ NOT STARTED

- [ ] Socket support (10-15h)
- [ ] HTTP functions (5-10h)
- [ ] DNS functions (5-8h)

### 13.5 Math Functions (10-15h) ⬜ NOT STARTED

- [ ] Trigonometry (4-6h)
- [ ] Random functions (3-5h)
- [ ] Number theory (3-4h)

### 13.6 Date/Time Functions (20-30h) ⬜ NOT STARTED

- [ ] Formatting (8-12h)
- [ ] Parsing (6-10h)
- [ ] Timezone support (6-8h)

### 13.7 Regular Expressions (15-20h) ⬜ NOT STARTED

- [ ] Complete preg_* functions (10-15h)
- [ ] Unicode support (5-8h)

### 13.8 Cryptography (15-20h) ⬜ NOT STARTED

- [ ] Password hashing (6-8h)
- [ ] Random bytes (4-6h)
- [ ] OpenSSL basics (5-6h)

---

## Phase 14: Extension System (80-120h)

**Goal**: Implement critical PHP extensions for WordPress/Laravel
**Priority**: HIGH - Required for database and network functionality

### 14.1 MySQLi Extension (30-40h) ⬜ NOT STARTED

**Priority**: CRITICAL - WordPress requires database

- [ ] Connection Management (10-12h)
  - mysqli_connect, mysqli_close, mysqli_select_db
  - mysqli_ping, mysqli_set_charset

- [ ] Query Execution (12-16h)
  - mysqli_query, mysqli_multi_query
  - mysqli_prepare, mysqli_stmt_bind_param, mysqli_stmt_execute

- [ ] Result Handling (8-12h)
  - mysqli_fetch_assoc, mysqli_fetch_array
  - mysqli_num_rows, mysqli_affected_rows

### 14.2 PDO Extension (25-35h) ⬜ NOT STARTED

**Priority**: HIGH - Laravel uses PDO

- [ ] PDO Class (15-20h)
  - Constructor with DSN parsing
  - query, exec, prepare
  - Transaction support

- [ ] PDOStatement Class (10-15h)
  - execute, fetch, fetchAll
  - Parameter binding

### 14.3 mbstring Extension (12-16h) ⬜ NOT STARTED

**Priority**: HIGH - UTF-8 support

- [ ] Multi-byte String Functions (12-16h)
  - mb_strlen, mb_substr, mb_strpos
  - mb_strtolower, mb_strtoupper
  - mb_detect_encoding, mb_convert_encoding

### 14.4 XML Extension (8-12h) ⬜ NOT STARTED

**Priority**: MEDIUM - WordPress imports/exports

- [ ] SimpleXML (5-7h)
  - simplexml_load_string, simplexml_load_file
  - XML tree navigation

- [ ] XML Parser (3-5h)
  - xml_parser_create, xml_parse
  - Event-based parsing

### 14.5 cURL Extension (10-15h) ⬜ NOT STARTED

**Priority**: HIGH - HTTP requests

- [ ] cURL Functions (10-15h)
  - curl_init, curl_setopt, curl_exec
  - HTTP methods (GET, POST, PUT, DELETE)

### 14.6 Other Extensions (5-10h) ⬜ NOT STARTED

- [ ] ctype (2-3h): Character type functions
- [ ] filter (3-5h): Input validation/filtering

---

## Phase 15: Production Optimization (60-80h)

**Goal**: Optimize for production WordPress/Laravel deployments
**Priority**: MEDIUM - After basic compatibility achieved

### 15.1 Performance Optimization (25-35h) ⬜ NOT STARTED

- [ ] Memory Pooling (10-15h)
  - Implement value pooling
  - Target: 50% memory reduction

- [ ] JIT Compilation (15-20h)
  - Compile hot paths to native code
  - Target: 2-3x additional speedup

### 15.2 Production Features (20-25h) ⬜ NOT STARTED

- [ ] Opcode Caching (8-10h)
  - Cache compiled bytecode on disk
  - opcache-like functionality

- [ ] Built-in Web Server (8-10h)
  - Implement `php-go -S localhost:8000`
  - FastCGI support

- [ ] Debugger Integration (4-5h)
  - Xdebug protocol support
  - Breakpoint debugging

### 15.3 Framework-Specific Optimization (15-25h) ⬜ NOT STARTED

- [ ] WordPress Profiling (5-7h)
  - Profile WordPress page loads
  - Identify bottlenecks

- [ ] WordPress Plugin Testing (5-8h)
  - Test top 20 plugins
  - Fix compatibility issues

- [ ] Laravel Profiling (3-5h)
  - Profile Laravel requests
  - Optimize service container

- [ ] Artisan Support (2-5h)
  - Ensure Artisan commands work
  - CLI performance

---

## Summary & Metrics

### Current Status

| Metric | Value |
|--------|-------|
| Total Hours Completed | 1340h |
| Project Completion | 93.6% |
| Benchmark Success Rate | 91% (10/11) |
| Performance vs PHP 8.4 | 5.0x faster |
| Standard Library Functions | 50+ |
| WordPress Parse Success | 61.20% (768/1255) |
| Laravel Parse Success | 43.57% (3260/7483) |

### Estimated Remaining Work

| Phase | Hours | Duration | Priority |
|-------|-------|----------|----------|
| Phase 10 (remaining) | 90 | 2-3 weeks | HIGH |
| Phase 11 (WordPress) | 180-220 | 5-6 weeks | HIGH |
| Phase 12 (Laravel) | 120-160 | 4-5 weeks | HIGH |
| Phase 13 (stdlib) | 200-300 | 8-10 weeks | MEDIUM |
| Phase 14 (extensions) | 80-120 | 3-4 weeks | HIGH |
| Phase 15 (optimization) | 60-80 | 2-3 weeks | MEDIUM |
| **TOTAL** | **730-970h** | **24-31 weeks** | **(6-8 months)** |

### Success Metrics

**WordPress Compatibility**:
- [ ] 95%+ parse success rate (currently 61.20%)
- [ ] WordPress installs and runs
- [ ] Admin panel fully functional
- [ ] Top 20 plugins work
- [ ] Performance: 3-5x faster than PHP 8.4

**Laravel Compatibility**:
- [ ] 95%+ parse success rate (currently 43.57%)
- [ ] Laravel app boots successfully
- [ ] All Artisan commands work
- [ ] Database operations work
- [ ] Performance: 3-5x faster than PHP 8.4

**Standard Library**:
- [ ] 300+ functions implemented (currently 50+)
- [ ] Core extensions (mysqli, PDO, mbstring, XML, curl)
- [ ] 90%+ test coverage

**Performance**:
- [ ] Maintain 5.0x speedup over PHP 8.4
- [ ] Memory usage: 50% of current (with pooling)
- [ ] Startup time: < 10ms
- [ ] 99th percentile latency: < PHP 8.4

---

## Prioritization

### Immediate Focus (Next 4 weeks)

1. ✅ Complete Phase 10 closures (15-25h) - **CRITICAL BLOCKER**
2. ⬜ Arrow functions (6-8h)
3. ⬜ Documentation (20h)
4. ⬜ WordPress require/include system (12-16h)
5. ⬜ Namespaces (16-20h)

### Short-term (Weeks 5-12)

1. ⬜ WordPress P0 features (global, unset, visibility, isset/empty)
2. ⬜ WordPress P1 features (const, switch, try/catch)
3. ⬜ Laravel-specific (references, destructuring)
4. ⬜ Begin standard library expansion

### Medium-term (Weeks 13-24)

1. ⬜ Extension system (mysqli, PDO, mbstring)
2. ⬜ Advanced features (attributes, anonymous classes)
3. ⬜ Complete standard library

### Long-term (Weeks 25+)

1. ⬜ Production optimization
2. ⬜ Framework testing and tuning
3. ⬜ Ecosystem development

---

## Notes

- **Previous Work**: Phases 0-9 details in `docs/CHANGELOG.md`
- **Detailed Roadmap**: See `docs/ROADMAP.md` for comprehensive documentation
- **Testing Strategy**: Maintain 85%+ code coverage, add regression tests for all new features
- **Commit Format**: `feat(phase<N>): <description>` or `fix(phase<N>): <description>`

---

**Last Updated**: November 24, 2025
**Next Milestone**: Complete Phase 10 closures (15-25h)
**Path to Production**: 730-970 hours (6-8 months) to full WordPress/Laravel support
