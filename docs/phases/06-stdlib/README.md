# Phase 6: Standard Library

## Overview

Phase 6 implements PHP's standard library - the massive collection of built-in functions that make PHP practical. This includes array functions, string functions, file I/O, and essential extensions.

## Scope

**Target**: Core + Standard Library from Phase 1 plan

**Essential Extensions**:
1. **standard** - The big one! (~66 files in PHP source)
2. **pcre** - Regular expressions
3. **date** - Date/time functions
4. **spl** - Standard PHP Library (data structures, iterators)
5. **json** - JSON encoding/decoding
6. **hash** - Hashing functions
7. **filter** - Input validation
8. **ctype** - Character type checking

## Goals

1. Implement ~300+ built-in functions
2. Full array function suite
3. Full string function suite
4. File I/O functions
5. Math functions
6. JSON support
7. Regular expressions (PCRE)
8. Date/time support

## Success Criteria

- [ ] All critical array functions implemented
- [ ] All critical string functions implemented
- [ ] File I/O working
- [ ] JSON encode/decode working
- [ ] Regular expressions working
- [ ] Date/time functions working
- [ ] Can run real PHP applications
- [ ] Test coverage >80%

## Function Categories

### 1. Array Functions (~100 functions)

**File**: `pkg/stdlib/array/`

**Priority 1 (Essential)**:
- count(), sizeof()
- array_keys(), array_values()
- array_push(), array_pop()
- array_shift(), array_unshift()
- array_merge(), array_merge_recursive()
- array_slice(), array_splice()
- array_search(), in_array()
- array_key_exists(), isset()
- array_unique()
- array_reverse()
- sort(), rsort(), asort(), arsort(), ksort(), krsort()
- usort(), uasort(), uksort()
- array_map(), array_filter(), array_reduce()
- array_walk(), array_walk_recursive()
- current(), next(), prev(), reset(), end(), key()
- each() (deprecated but still used)

**Priority 2 (Common)**:
- array_chunk(), array_column()
- array_combine(), array_fill(), array_fill_keys()
- array_flip(), array_intersect(), array_diff()
- array_pad(), array_product(), array_sum()
- array_rand(), shuffle()
- range(), array_replace()

**Priority 3 (Less Common)**:
- array_change_key_case(), array_count_values()
- array_diff_assoc(), array_diff_key()
- array_intersect_assoc(), array_intersect_key()
- array_multisort()
- array_udiff(), array_uintersect()
- compact(), extract()

**Reference**: `php-src/ext/standard/array.c` (198KB!)

### 2. String Functions (~100 functions)

**File**: `pkg/stdlib/string/`

**Priority 1 (Essential)**:
- strlen()
- substr(), mb_substr()
- strpos(), strrpos(), stripos(), strripos()
- str_replace(), str_ireplace()
- strtolower(), strtoupper()
- ucfirst(), ucwords()
- trim(), ltrim(), rtrim()
- explode(), implode(), join()
- str_split(), chunk_split()
- sprintf(), vsprintf()
- str_pad(), str_repeat()
- strcmp(), strcasecmp(), strncmp(), strncasecmp()
- strstr(), stristr(), strrchr()
- htmlspecialchars(), htmlentities()
- addslashes(), stripslashes()
- nl2br(), wordwrap()

**Priority 2 (Common)**:
- substr_count(), substr_replace()
- str_contains(), str_starts_with(), str_ends_with() (PHP 8.0+)
- str_word_count()
- strrev(), strtr()
- parse_str(), http_build_query()
- base64_encode(), base64_decode()
- urlencode(), urldecode(), rawurlencode(), rawurldecode()
- quoted_printable_encode(), quoted_printable_decode()
- md5(), sha1() (basic hashing)
- crypt(), password_hash(), password_verify()

**Reference**: `php-src/ext/standard/string.c` (154KB!)

### 3. File I/O Functions (~50 functions)

**File**: `pkg/stdlib/file/`

**Priority 1 (Essential)**:
- fopen(), fclose()
- fread(), fwrite()
- fgets(), fgetss(), fgetc()
- fputs(), fputc()
- file_get_contents(), file_put_contents()
- file(), readfile()
- feof(), ftell(), fseek(), rewind()
- is_file(), is_dir(), is_readable(), is_writable(), is_executable()
- file_exists(), filesize(), filetype()
- mkdir(), rmdir(), unlink(), rename(), copy()
- dirname(), basename(), pathinfo()
- realpath(), getcwd(), chdir()

**Priority 2 (Common)**:
- glob(), scandir(), opendir(), readdir(), closedir()
- fstat(), stat(), lstat()
- filemtime(), fileatime(), filectime()
- chmod(), chown(), chgrp()
- touch(), umask()
- flock(), ftruncate()
- tempnam(), tmpfile()

**Reference**: `php-src/ext/standard/file.c` (57KB)

### 4. Variable Functions (~20 functions)

**File**: `pkg/stdlib/var/`

**Priority 1**:
- var_dump(), print_r()
- var_export()
- serialize(), unserialize()
- gettype(), settype()
- is_null(), is_bool(), is_int(), is_float(), is_string(), is_array(), is_object(), is_resource()
- is_numeric(), is_scalar()
- isset(), empty()
- unset()

**Reference**: `php-src/ext/standard/var.c` (46KB)

### 5. Math Functions (~30 functions)

**File**: `pkg/stdlib/math/`

**Priority 1**:
- abs(), ceil(), floor(), round()
- min(), max()
- pow(), sqrt(), exp(), log(), log10()
- sin(), cos(), tan(), asin(), acos(), atan(), atan2()
- rand(), mt_rand(), random_int()
- srand(), mt_srand()
- number_format()
- deg2rad(), rad2deg()
- pi(), M_PI constants

### 6. JSON Extension

**File**: `pkg/stdlib/json/`

**Functions**:
- json_encode()
- json_decode()
- json_last_error()
- json_last_error_msg()

**Reference**: `php-src/ext/json/` (~10K lines)

### 7. PCRE Extension (Regular Expressions)

**File**: `pkg/stdlib/pcre/`

**Functions**:
- preg_match()
- preg_match_all()
- preg_replace()
- preg_replace_callback()
- preg_split()
- preg_grep()
- preg_quote()
- preg_last_error()

**Implementation**: Use Go's regexp package + custom PCRE compatibility layer

**Reference**: `php-src/ext/pcre/` (~20K lines)

### 8. Date/Time Extension

**File**: `pkg/stdlib/date/`

**Functions**:
- date(), gmdate()
- time(), microtime()
- strtotime()
- mktime(), gmmktime()
- date_create(), date_format()
- DateTime class
- DateTimeImmutable class
- DateInterval class
- DateTimeZone class

**Reference**: `php-src/ext/date/` (~30K lines)

### 9. SPL Extension

**File**: `pkg/stdlib/spl/`

**Data Structures**:
- SplStack, SplQueue
- SplHeap, SplMaxHeap, SplMinHeap
- SplPriorityQueue
- SplFixedArray
- SplDoublyLinkedList

**Iterators**:
- DirectoryIterator, RecursiveDirectoryIterator
- ArrayIterator, RecursiveArrayIterator
- Iterator, IteratorAggregate

**Functions**:
- spl_autoload_register()
- class_implements(), class_parents(), class_uses()
- iterator_to_array(), iterator_count()

**Reference**: `php-src/ext/spl/` (~40K lines)

### 10. Hash Extension

**File**: `pkg/stdlib/hash/`

**Functions**:
- hash(), hash_file()
- hash_hmac(), hash_hmac_file()
- hash_algos()
- md5(), md5_file()
- sha1(), sha1_file()
- hash_equals()

**Algorithms**: md5, sha1, sha256, sha512, etc.

**Reference**: `php-src/ext/hash/` (~15K lines)

### 11. Filter Extension

**File**: `pkg/stdlib/filter/`

**Functions**:
- filter_var()
- filter_var_array()
- filter_input()
- filter_input_array()
- filter_list()

**Filters**: FILTER_VALIDATE_EMAIL, FILTER_VALIDATE_INT, FILTER_SANITIZE_STRING, etc.

**Reference**: `php-src/ext/filter/` (~10K lines)

### 12. Ctype Extension

**File**: `pkg/stdlib/ctype/`

**Functions**:
- ctype_alnum(), ctype_alpha(), ctype_cntrl()
- ctype_digit(), ctype_graph()
- ctype_lower(), ctype_upper()
- ctype_print(), ctype_punct()
- ctype_space(), ctype_xdigit()

**Reference**: `php-src/ext/ctype/` (~5K lines)

## Implementation Strategy

### Phase 6A: Array & String Functions (3-4 weeks)
- All essential array functions
- All essential string functions
- File I/O basics
- Variable functions

### Phase 6B: Extensions (2-3 weeks)
- JSON
- PCRE (regex)
- Date/Time basics
- Math functions

### Phase 6C: SPL & Advanced (2 weeks)
- SPL data structures
- SPL iterators
- Hash functions
- Filter/Ctype

## Implementation Tasks

### Task 6.1: Array Functions - Basic
**Effort**: 16 hours

- [ ] count(), sizeof()
- [ ] array_keys(), array_values()
- [ ] array_push(), array_pop(), array_shift(), array_unshift()
- [ ] array_merge()
- [ ] in_array(), array_search()
- [ ] array_slice(), array_splice()

### Task 6.2: Array Functions - Advanced
**Effort**: 20 hours

- [ ] Sorting functions (sort, rsort, asort, etc.)
- [ ] array_map(), array_filter(), array_reduce()
- [ ] array_walk(), array_walk_recursive()
- [ ] array_diff(), array_intersect()
- [ ] Array pointer functions (current, next, etc.)

### Task 6.3: String Functions - Basic
**Effort**: 16 hours

- [ ] strlen(), substr()
- [ ] strpos(), strrpos()
- [ ] str_replace()
- [ ] strtolower(), strtoupper()
- [ ] trim(), ltrim(), rtrim()
- [ ] explode(), implode()

### Task 6.4: String Functions - Advanced
**Effort**: 20 hours

- [ ] sprintf() family
- [ ] str_pad(), str_repeat(), str_split()
- [ ] strcmp() family
- [ ] strstr(), strrchr()
- [ ] htmlspecialchars(), htmlentities()
- [ ] URL encoding functions
- [ ] Base64 functions

### Task 6.5: File I/O
**Effort**: 20 hours

- [ ] fopen/fclose/fread/fwrite
- [ ] file_get_contents/file_put_contents
- [ ] file(), readfile()
- [ ] File info functions
- [ ] Directory functions
- [ ] Path functions

### Task 6.6: Variable Functions
**Effort**: 8 hours

- [ ] var_dump(), print_r()
- [ ] serialize(), unserialize()
- [ ] Type checking functions
- [ ] gettype(), settype()

### Task 6.7: Math Functions
**Effort**: 8 hours

- [ ] Basic math (abs, ceil, floor, round, etc.)
- [ ] Trigonometric functions
- [ ] Random number generation
- [ ] number_format()

### Task 6.8: JSON Extension
**Effort**: 12 hours

- [ ] json_encode() implementation
- [ ] json_decode() implementation
- [ ] Options and flags
- [ ] Error handling

**Use Go's encoding/json as base**

### Task 6.9: PCRE Extension
**Effort**: 20 hours

- [ ] preg_match()
- [ ] preg_match_all()
- [ ] preg_replace()
- [ ] preg_split()
- [ ] Pattern compilation
- [ ] PCRE compatibility layer

**Challenging - Go's regexp is different from PCRE!**

### Task 6.10: Date/Time Extension
**Effort**: 16 hours

- [ ] date(), gmdate()
- [ ] time(), microtime()
- [ ] strtotime() (complex!)
- [ ] DateTime class
- [ ] Timezone support

### Task 6.11: SPL Data Structures
**Effort**: 16 hours

- [ ] SplStack, SplQueue
- [ ] SplHeap
- [ ] SPL classes
- [ ] Iterator interfaces

### Task 6.12: Hash/Filter/Ctype
**Effort**: 12 hours

- [ ] Hash functions
- [ ] Filter functions
- [ ] Ctype functions

### Task 6.13: Testing
**Effort**: 24 hours

- [ ] Test every function
- [ ] Edge cases
- [ ] PHP compatibility tests
- [ ] Real-world usage tests

## Estimated Timeline

**Total Effort**: ~210 hours (10-12 weeks)

This is the longest phase!

**Breakdown**:
- Array functions: ~36 hours
- String functions: ~36 hours
- File I/O: ~20 hours
- Other stdlib: ~16 hours
- JSON: ~12 hours
- PCRE: ~20 hours
- Date/Time: ~16 hours
- SPL: ~16 hours
- Hash/Filter/Ctype: ~12 hours
- Testing: ~24 hours

## Progress Tracking

- [ ] Task 6.1: Array Functions - Basic (0%)
- [ ] Task 6.2: Array Functions - Advanced (0%)
- [ ] Task 6.3: String Functions - Basic (0%)
- [ ] Task 6.4: String Functions - Advanced (0%)
- [ ] Task 6.5: File I/O (0%)
- [ ] Task 6.6: Variable Functions (0%)
- [ ] Task 6.7: Math Functions (0%)
- [ ] Task 6.8: JSON Extension (0%)
- [ ] Task 6.9: PCRE Extension (0%)
- [ ] Task 6.10: Date/Time Extension (0%)
- [ ] Task 6.11: SPL Data Structures (0%)
- [ ] Task 6.12: Hash/Filter/Ctype (0%)
- [ ] Task 6.13: Testing (0%)

**Overall Phase 6 Progress**: 0%

**Note**: This is the largest and most time-consuming phase. After this, you'll be able to run real PHP applications!

---

**Status**: Not Started
**Last Updated**: 2025-11-21
