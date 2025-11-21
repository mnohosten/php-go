# PHP 8.4 Source Code Analysis

This document contains a comprehensive analysis of the PHP 8.4 interpreter source code from the `php-src` repository.

## Executive Summary

- **Version**: PHP 8.6.0-dev (development branch, targeting 8.4 features)
- **Total C Files**: 1,034 files
- **Total Header Files**: 1,073 files
- **Total Lines of C Code**: ~1,203,034 lines
- **Extensions**: 72 total (70 real + 2 test/skeleton)
- **Opcodes**: 210 opcodes in the Zend VM

## Core Architecture

PHP consists of several major components:

1. **Zend Engine** - The core VM and runtime
2. **Main** - SAPI abstraction and PHP core functions
3. **SAPI** - Server API implementations (CLI, FPM, CGI, etc.)
4. **Extensions** - Standard library and optional functionality
5. **TSRM** - Thread-Safe Resource Manager

## Zend Engine (`/Zend/`)

The heart of PHP - implements the virtual machine, compiler, and runtime.

**Statistics**: 79 C files, 108 header files, 187 total files

### Lexer & Parser

**Files**:
- `/Zend/zend_language_scanner.l` (3,231 lines) - Lexical analyzer (Flex/re2c)
- `/Zend/zend_language_parser.y` (1,908 lines) - Grammar definition (Bison)
- `/Zend/zend_ini_scanner.l` - INI file lexer
- `/Zend/zend_ini_parser.y` - INI parser

**Key Features**:
- Hand-written lexer using re2c (generates C code)
- LALR parser using Bison
- Produces AST (not opcodes directly)
- Supports PHP's complex syntax (heredoc, nowdoc, string interpolation, etc.)

### Compiler

**Files**:
- `/Zend/zend_compile.c` (12,479 lines, 378KB) - Main compiler
- `/Zend/zend_compile.h` (55,425 lines) - Compilation structures
- `/Zend/zend_ast.c` (86KB) - Abstract Syntax Tree implementation
- `/Zend/zend_ast.h` - AST node definitions

**Compilation Pipeline**:
1. Source code → Tokens (lexer)
2. Tokens → AST (parser)
3. AST → Opcodes (compiler)
4. Opcodes → Optimized opcodes (optimizer, optional)

**Key Responsibilities**:
- AST traversal and opcode generation
- Symbol table management
- Compile-time constant folding
- Early binding optimization
- Error checking and diagnostics

### Virtual Machine

**Files**:
- `/Zend/zend_vm_def.h` (10,573 lines) - VM opcode handler definitions
- `/Zend/zend_vm_execute.h` (129,697 lines, 4.2MB) - **GENERATED** VM executor
- `/Zend/zend_vm_opcodes.h` - Opcode constants (210 opcodes)
- `/Zend/zend_vm_gen.php` - Generator script (creates executor from def)
- `/Zend/zend_execute.c` (5,949 lines, 182KB) - Execution API
- `/Zend/zend_execute_API.c` (55KB) - Execution helpers

**VM Architecture**:
- **Specialized VM**: Each opcode has specialized handlers for different operand types
  - CONST: Compile-time constant
  - VAR: Variable (temporary)
  - TMP: Temporary value
  - CV: Compiled variable (faster variable access)
- **VM Kinds**: Multiple threading models
  - CALL: Function call-based dispatch
  - SWITCH: Switch-based dispatch
  - GOTO: Computed goto (fastest on GCC)
  - HYBRID: Mix of goto and switch
  - TAILCALL: Tail-call optimization
- **Generated Code**: The `zend_vm_gen.php` script generates optimized C code
- **210 Opcodes**: Everything from `ZEND_NOP` to complex operations

**Opcode Categories**:
- Arithmetic: ADD, SUB, MUL, DIV, MOD, POW
- Bitwise: BW_OR, BW_AND, BW_XOR, SL, SR
- Comparison: IS_EQUAL, IS_IDENTICAL, IS_SMALLER, etc.
- Logic: BOOL_NOT, BOOL_AND, BOOL_OR
- Control Flow: JMP, JMPZ, JMPNZ, SWITCH, MATCH
- Variables: ASSIGN, FETCH, UNSET, ISSET
- Arrays: FETCH_DIM, ASSIGN_DIM, INIT_ARRAY
- Objects: NEW, FETCH_OBJ, ASSIGN_OBJ, INIT_METHOD_CALL
- Functions: INIT_FCALL, DO_FCALL, RETURN, YIELD
- And many more...

### Memory Management

**Files**:
- `/Zend/zend_alloc.c` (107KB) - Zend Memory Manager
- `/Zend/zend_alloc.h` (20,807 lines)
- `/Zend/zend_gc.c` (62KB) - Garbage collector

**Key Features**:
- **Custom Allocator**: Optimized for PHP's allocation patterns
- **Request-based Pools**: Fast allocation/deallocation per request
- **Reference Counting**: Manual ref counting for most objects
- **Cycle Detection GC**: Handles circular references
- **Arena Allocation**: For compiler temporaries

**Memory Model**:
- Reference counting as primary mechanism
- Cycle-detecting GC runs periodically
- Copy-on-write for strings and arrays
- Interned strings (deduplicated)
- Small-string optimization

### Type System

**Files**:
- `/Zend/zend_types.h` (52,919 lines) - Core type definitions
- `/Zend/zend_type_info.h` - Type information
- `/Zend/zend_operators.c` (101KB) - Type conversions/operations

**Core Types**:

```c
// zval - Universal value container (16 bytes on 64-bit)
struct _zval_struct {
    zend_value value;        // 8 bytes - union of all types
    union {
        uint32_t type_info;  // Type tag + flags
        struct { ... }       // GC info for references
    } u1;
    union { ... } u2;        // Extra info (cache slot, etc.)
};

// Type tags
#define IS_UNDEF     0
#define IS_NULL      1
#define IS_FALSE     2
#define IS_TRUE      3
#define IS_LONG      4
#define IS_DOUBLE    5
#define IS_STRING    6
#define IS_ARRAY     7
#define IS_OBJECT    8
#define IS_RESOURCE  9
#define IS_REFERENCE 10
```

**Complex Types** (PHP 8.0+):
- Union types: `int|string`
- Intersection types: `A&B`
- Mixed type
- Never type
- Void type
- Static type
- DNF types (Disjunctive Normal Form)

### Data Structures

**Hash Tables** (`/Zend/zend_hash.c`, 87KB):
- **PHP arrays are hash tables** - Ordered associative arrays
- Maintains insertion order
- Optimized for sequential integer keys ("packed arrays")
- Separate chaining for collisions
- Efficient iteration
- Copy-on-write

**Strings** (`/Zend/zend_string.c`, 16KB):
- Reference-counted strings
- Interned strings (deduplicated, immortal)
- Copy-on-write
- Length-prefixed (O(1) strlen)

**Linked Lists** (`/Zend/zend_llist.c`):
- Simple doubly-linked lists
- Used for various internal purposes

### Object System

**Files**:
- `/Zend/zend_object_handlers.c` (77KB) - Property/method access
- `/Zend/zend_objects.c` (12KB) - Object allocation
- `/Zend/zend_objects_API.c` (7KB) - Object API
- `/Zend/zend_inheritance.c` (142KB) - Class inheritance (COMPLEX!)
- `/Zend/zend_API.c` (159KB, 5,301 lines) - Internal API

**Features**:
- Single inheritance + interfaces
- Traits (horizontal code reuse)
- Magic methods (__get, __set, __call, etc.)
- Property visibility (public, protected, private)
- Static properties and methods
- Late static binding
- Type hints for properties
- Readonly properties (PHP 8.1+)
- Readonly classes (PHP 8.2+)

### Advanced Features

**Generators** (`/Zend/zend_generators.c`, 38KB):
- Implemented as special objects
- Maintain execution state between yields
- Support yield from

**Fibers** (`/Zend/zend_fibers.c`, 34KB):
- Lightweight coroutines (PHP 8.1+)
- Cooperative multitasking
- Separate stack for each fiber

**Closures** (`/Zend/zend_closures.c`, 32KB):
- Anonymous functions
- Variable capture (by value or reference)
- Arrow functions (PHP 7.4+)

**Exceptions** (`/Zend/zend_exceptions.c`, 34KB):
- Exception hierarchy
- Stack trace generation
- Try-catch-finally

**Other Advanced Features**:
- `/Zend/zend_weakrefs.c` (28KB) - Weak references
- `/Zend/zend_lazy_objects.c` (26KB) - Lazy initialization
- `/Zend/zend_enum.c` (22KB) - Enums (PHP 8.1+)
- `/Zend/zend_attributes.c` (18KB) - Attributes/annotations (PHP 8.0+)

### Optimizer (`/Zend/Optimizer/`)

**Files**: 35 files implementing SSA-based optimization

**Key Components**:
- `/Zend/Optimizer/zend_optimizer.c` (58KB) - Main optimizer
- `/Zend/Optimizer/zend_inference.c` (166KB) - Type inference
- `/Zend/Optimizer/sccp.c` (72KB) - Sparse Conditional Constant Propagation
- `/Zend/Optimizer/dfa_pass.c` (57KB) - Data flow analysis
- `/Zend/Optimizer/zend_ssa.c` (57KB) - Static Single Assignment form

**Optimizations**:
- Constant folding and propagation
- Dead code elimination
- Type narrowing via inference
- Function inlining (limited)
- Opcode specialization

### Built-in Functions

**Files**:
- `/Zend/zend_builtin_functions.c` (63KB) - Core PHP functions
- `/Zend/zend_constants.c` (17KB) - Constant management
- `/Zend/zend_ini.c` (28KB) - INI configuration
- `/Zend/zend_interfaces.c` (22KB) - Standard interfaces

## Main (`/main/`)

SAPI abstraction and core PHP runtime.

**Files**: 75 total files

**Key Components**:
- `/main/main.c` (84KB) - PHP engine startup/shutdown
- `/main/php.h` - Main PHP header
- `/main/SAPI.c` (32KB) - Server API abstraction
- `/main/php_ini.c` (27KB) - INI configuration
- `/main/php_variables.c` (28KB) - Request variables ($_GET, $_POST, etc.)
- `/main/output.c` (50KB) - Output buffering
- `/main/network.c` (35KB) - Network functions
- `/main/fopen_wrappers.c` (23KB) - File opening wrappers
- `/main/rfc1867.c` (35KB) - Multipart form-data (file uploads)
- `/main/fastcgi.c` (41KB) - FastCGI protocol

**Streams** (`/main/streams/`):
- `/main/streams/streams.c` (72KB) - Stream abstraction
- `/main/streams/plain_wrapper.c` (47KB) - File stream wrapper
- `/main/streams/userspace.c` (43KB) - Userspace streams
- `/main/streams/transports.c` (13KB) - Network transports
- Plus filters, memory streams, mmap support

## SAPI (`/sapi/`)

Server API implementations - how PHP runs in different environments.

**Available SAPIs**:
1. **cli** - Command-line interface (essential)
2. **fpm** - FastCGI Process Manager (production web server)
3. **cgi** - CGI interface
4. **apache2handler** - Apache module
5. **embed** - Embedded PHP library
6. **phpdbg** - PHP debugger
7. **litespeed** - LiteSpeed server
8. **fuzzer** - Fuzzing interface

**For php-go**: We'll start with CLI SAPI only.

## Extensions (`/ext/`)

**Total**: 72 extensions (70 real + 2 test)

### Essential Extensions (Must Implement)

#### 1. standard
**The most important extension!**

- **Files**: 66 C files
- **Size**: Massive (~300KB total)
- **Key Files**:
  - `array.c` (198KB) - Array functions (sort, map, filter, etc.)
  - `string.c` (154KB) - String functions (strlen, substr, str_replace, etc.)
  - `file.c` (57KB) - File I/O (fopen, fread, fwrite, etc.)
  - `var.c` (46KB) - Variable handling (serialize, unserialize, var_dump, etc.)
  - `math.c` - Mathematical functions
  - `dir.c` - Directory functions
  - `exec.c` - Process execution
  - `basic_functions.c` - Misc core functions
  - And ~55 more files

**Categories**:
- Arrays: ~100 functions
- Strings: ~100 functions
- Files: ~50 functions
- Math: ~30 functions
- Variables: ~20 functions
- Process: ~15 functions
- Misc: ~100+ functions

#### 2. pcre (Regular Expressions)
- Perl-Compatible Regular Expressions
- Wrapper around PCRE2 library
- Essential for text processing
- ~20K lines

#### 3. date (Date/Time)
- Date and time functions
- Timezone support
- Date formatting and parsing
- ~30K lines

#### 4. spl (Standard PHP Library)
- Data structures (SplStack, SplQueue, SplHeap, etc.)
- Iterators (DirectoryIterator, RecursiveIterator, etc.)
- Interfaces (Countable, ArrayAccess, etc.)
- Exceptions
- ~40K lines

#### 5. json
- JSON encoding/decoding
- Essential for modern web apps
- ~10K lines

#### 6. hash
- Hashing functions (md5, sha1, sha256, etc.)
- HMAC support
- ~15K lines

#### 7. filter
- Input validation and sanitization
- Essential for security
- ~10K lines

#### 8. ctype
- Character type checking
- Simple but commonly used
- ~5K lines

#### 9. tokenizer
- PHP tokenizer
- Used by tools and IDEs
- ~5K lines

#### 10. reflection
- Runtime introspection
- Essential for frameworks
- ~20K lines

### Database Extensions (Important)

- **pdo** - PHP Data Objects (abstraction layer)
- **pdo_mysql**, **pdo_sqlite**, **pdo_pgsql** - PDO drivers
- **mysqli** - MySQL improved extension
- **mysqlnd** - MySQL native driver
- **sqlite3** - SQLite 3
- **pgsql** - PostgreSQL

### Network Extensions

- **curl** - HTTP client (~30K lines)
- **ftp** - FTP client
- **sockets** - Low-level sockets

### XML Extensions

- **libxml** - XML base support
- **dom** - DOM API
- **simplexml** - SimpleXML API
- **xml** - XML parser
- **xmlreader** / **xmlwriter**

### Compression

- **zlib** - Gzip compression (~10K lines)
- **bz2** - Bzip2
- **zip** - ZIP archives
- **phar** - PHP Archives

### Cryptography

- **openssl** - OpenSSL crypto (~50K lines)
- **sodium** - Libsodium (modern crypto)
- **random** - CSPRNG (PHP 8.2+)
- **password** - Password hashing

### String/Encoding

- **mbstring** - Multi-byte strings
- **iconv** - Character encoding conversion
- **intl** - Internationalization (ICU)

### Image Processing

- **gd** - GD graphics
- **exif** - Image metadata

### System/IPC

- **posix** - POSIX functions
- **pcntl** - Process control
- **shmop** - Shared memory
- **sysvmsg**, **sysvsem**, **sysvshm** - System V IPC

### Other Notable Extensions

- **opcache** - Opcode cache + JIT compiler (very complex, can defer)
- **ffi** - Foreign Function Interface
- **calendar** - Calendar conversions
- **bcmath** - Arbitrary precision math
- **gmp** - GNU MP math

## JIT Compiler (`/ext/opcache/jit/`)

**Extremely complex** - can skip for initial implementation.

**Components**:
- Tracing JIT (traces hot code paths)
- IR library (Intermediate Representation)
- DynASM (runtime assembler)
- x86-64 and ARM64 backends
- 6 major files (~1.2MB of code)

## Scope Estimates

### Minimal Core (Basic PHP Execution)
**~250K lines of C code to port**:
- Lexer/Parser (~5K lines)
- Compiler (~25K lines)
- VM Executor (~20K lines)
- Type System (~30K lines)
- Memory Manager (~35K lines)
- Hash Tables (~10K lines)
- Strings (~5K lines)
- Core Functions (~5K lines)
- SAPI/CLI (~10K lines)
- Main Runtime (~20K lines)

### Core + Standard Library
**~600K lines of C code to port**:
- All of Minimal Core
- ext/standard (~150K lines) - THE BIG ONE
- ext/pcre (~20K lines)
- ext/date (~30K lines)
- ext/spl (~40K lines)
- ext/json (~10K lines)
- ext/hash (~15K lines)
- ext/filter (~10K lines)
- Object System (~60K lines)
- Exceptions (~10K lines)
- Closures (~10K lines)
- Generators (~10K lines)
- Reflection (~20K lines)

### Production-Ready
**~800K-1M lines**:
- Everything above
- Database extensions (~50K lines)
- XML extensions (~40K lines)
- Compression (~30K lines)
- Network extensions (~50K lines)

### Full PHP 8.4
**~1.2M lines**:
- Everything above
- All remaining extensions
- Full compatibility

## Key Insights for Go Implementation

### Simplifications Possible

1. **Memory Management**: Go GC eliminates need for:
   - Reference counting
   - Circular reference detection
   - Manual memory pools
   - Arena allocation

2. **Thread Safety**: Go's goroutines different from TSRM:
   - No need for complex thread-local storage
   - Simpler concurrency model
   - Built-in synchronization primitives

3. **String Handling**: Go strings simplify:
   - No manual interning needed (Go does it)
   - Immutable strings by default
   - Unicode support built-in

4. **Type System**: Go's interfaces simplify:
   - Type-safe implementation
   - Clear polymorphism
   - Compile-time checking where possible

### Complexities to Handle

1. **Inheritance System**: 142KB of complex C code
   - Property visibility rules
   - Method override checking
   - Trait resolution
   - Interface compliance

2. **Type Juggling**: PHP's loose typing
   - Automatic type conversions
   - Comparison rules
   - String to number coercion
   - Array to scalar conversions

3. **Variable Variables**: Dynamic names
   - `$$var`, `$obj->$prop`, `$class::$static`
   - Requires runtime symbol lookup

4. **References**: PHP references
   - Not same as Go pointers
   - Reference semantics for assignment
   - Reference parameters

5. **Error Handling**: Multiple error modes
   - Exceptions
   - Errors (PHP 7+ Error class)
   - Warnings/Notices
   - @ error suppression

## Test Suite

PHP has comprehensive tests:
- **Location**: `/tests/` and per-extension `tests/` directories
- **Format**: `.phpt` files (PHP test format)
- **Coverage**: Thousands of tests
- **Categories**: Language, functions, extensions, edge cases

**For php-go**: We can reuse these tests for compatibility verification.

## Build System

- **Unix**: Autoconf-based (`./configure`)
- **Windows**: CMake/NMake
- **Extension System**: Each extension has config.m4/config.w32

**For php-go**: Go's build system much simpler.

## Documentation

- `/docs/` - Sphinx-based internal documentation
- `/README.md` - Main readme
- `/CODING_STANDARDS.md` - Coding style
- `/UPGRADING` - Version upgrade notes
- `/UPGRADING.INTERNALS` - API changes

## Recommended Reading Order

For implementing php-go, study these files first:

1. `/Zend/zend_types.h` - Understand core types
2. `/Zend/zend_vm_opcodes.h` - All 210 opcodes
3. `/Zend/zend_language_parser.y` - PHP grammar
4. `/Zend/zend_compile.h` - Compilation structures
5. `/Zend/zend_vm_def.h` - VM handler definitions
6. `/ext/standard/array.c` - Example of extension functions
7. `/Zend/zend_hash.c` - Hash table implementation

## Conclusion

PHP 8.4 is a mature, complex interpreter with:
- Clean separation of concerns (lexer → parser → compiler → VM)
- Well-optimized core data structures
- Extensive standard library
- Complex features (JIT, fibers, optimizer)

The total scope is large (~1.2M lines), but core + standard library (~600K lines) is achievable for a Go rewrite.

---

**Source**: Analyzed from php-src repository (PHP 8.6.0-dev branch)
**Last Updated**: 2025-11-21
