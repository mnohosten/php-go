# Standard Library API Reference

Package: `github.com/krizos/php-go/pkg/stdlib`

## Overview

The stdlib package will implement PHP's standard library functions (~300+ functions). This package is planned for Phase 6 and is currently a placeholder.

**Status**: Not yet implemented (Phase 6 - Planned)

**Implementation Timeline**: Phase 6 of the project roadmap

**Estimated Scope**: ~300+ PHP built-in functions across multiple categories

## Planned Function Categories

### Array Functions (`pkg/stdlib/array/`)
Functions for manipulating PHP arrays.

**Planned Functions (50+ functions):**
- `array_push`, `array_pop`, `array_shift`, `array_unshift`
- `array_merge`, `array_combine`, `array_slice`, `array_splice`
- `array_map`, `array_filter`, `array_reduce`, `array_walk`
- `array_keys`, `array_values`, `array_flip`, `array_reverse`
- `array_search`, `in_array`, `array_key_exists`
- `array_sum`, `array_product`, `array_unique`
- `sort`, `asort`, `ksort`, `usort`, `uasort`, `uksort`
- `rsort`, `arsort`, `krsort`, `shuffle`, `array_rand`
- `array_chunk`, `array_pad`, `array_fill`, `array_fill_keys`
- `array_intersect`, `array_diff`, `array_column`
- And many more...

### String Functions (`pkg/stdlib/string/`)
Functions for string manipulation and processing.

**Planned Functions (80+ functions):**
- `strlen`, `substr`, `str_replace`, `str_ireplace`
- `strpos`, `strrpos`, `stripos`, `strripos`
- `strtolower`, `strtoupper`, `ucfirst`, `ucwords`
- `trim`, `ltrim`, `rtrim`, `str_pad`
- `str_repeat`, `str_split`, `chunk_split`
- `explode`, `implode`, `join`
- `strcmp`, `strcasecmp`, `strncmp`, `strncasecmp`
- `sprintf`, `vsprintf`, `sscanf`
- `str_contains`, `str_starts_with`, `str_ends_with`
- `htmlspecialchars`, `htmlentities`, `html_entity_decode`
- `addslashes`, `stripslashes`, `quotemeta`
- `nl2br`, `wordwrap`, `str_word_count`
- And many more...

### Math Functions (`pkg/stdlib/math/`)
Mathematical operations and functions.

**Planned Functions (30+ functions):**
- `abs`, `ceil`, `floor`, `round`
- `min`, `max`, `pow`, `sqrt`, `exp`, `log`, `log10`
- `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`
- `sinh`, `cosh`, `tanh`
- `deg2rad`, `rad2deg`
- `pi`, `rand`, `mt_rand`, `random_int`
- `is_nan`, `is_finite`, `is_infinite`
- `intdiv`, `fmod`, `hypot`
- And more...

### File I/O Functions (`pkg/stdlib/file/`)
Functions for file system operations.

**Planned Functions (50+ functions):**
- `fopen`, `fclose`, `fread`, `fwrite`, `fgets`, `fputs`
- `file_get_contents`, `file_put_contents`
- `file`, `readfile`, `file_exists`
- `is_file`, `is_dir`, `is_readable`, `is_writable`
- `mkdir`, `rmdir`, `unlink`, `rename`, `copy`
- `scandir`, `glob`, `dirname`, `basename`, `pathinfo`
- `realpath`, `chmod`, `chown`, `chgrp`
- `filesize`, `filetype`, `filemtime`, `filectime`, `fileatime`
- `flock`, `fseek`, `ftell`, `rewind`, `feof`, `fflush`
- `tempnam`, `tmpfile`
- And more...

### Date/Time Functions (`pkg/stdlib/date/`)
Functions for date and time manipulation.

**Planned Functions (20+ functions):**
- `time`, `microtime`, `date`, `gmdate`, `idate`
- `strtotime`, `mktime`, `gmmktime`
- `date_create`, `date_format`, `date_parse`
- `date_diff`, `date_add`, `date_sub`
- `getdate`, `localtime`, `gmdate`
- `checkdate`, `date_default_timezone_set`
- And more...

### JSON Functions (`pkg/stdlib/json/`)
JSON encoding and decoding.

**Planned Functions:**
- `json_encode` - Encode value to JSON
- `json_decode` - Decode JSON to PHP value
- `json_last_error` - Returns the last JSON error
- `json_last_error_msg` - Returns error message

### PCRE Functions (`pkg/stdlib/pcre/`)
Perl-compatible regular expression functions.

**Planned Functions (10+ functions):**
- `preg_match` - Perform a regular expression match
- `preg_match_all` - Perform a global regular expression match
- `preg_replace` - Perform a regular expression search and replace
- `preg_replace_callback` - Replace using callback
- `preg_split` - Split string by a regular expression
- `preg_grep` - Return array entries that match pattern
- `preg_quote` - Quote regular expression characters
- And more...

### Variable Functions (`pkg/stdlib/var/`)
Functions for variable inspection and manipulation.

**Planned Functions (20+ functions):**
- `var_dump`, `var_export`, `print_r`
- `isset`, `empty`, `unset`
- `is_null`, `is_bool`, `is_int`, `is_float`, `is_string`
- `is_array`, `is_object`, `is_resource`, `is_callable`
- `is_numeric`, `is_scalar`
- `gettype`, `settype`
- `serialize`, `unserialize`
- `get_defined_vars`, `get_defined_functions`, `get_defined_constants`
- And more...

### Character Type Functions (`pkg/stdlib/ctype/`)
Character type checking functions.

**Planned Functions (10+ functions):**
- `ctype_alnum`, `ctype_alpha`, `ctype_digit`
- `ctype_lower`, `ctype_upper`, `ctype_space`
- `ctype_punct`, `ctype_xdigit`, `ctype_cntrl`
- And more...

### Filter Functions (`pkg/stdlib/filter/`)
Data filtering and validation.

**Planned Functions:**
- `filter_var` - Filter a variable with a specified filter
- `filter_input` - Get input from external source and filter
- `filter_var_array` - Get multiple variables and filter them
- `filter_list` - Returns a list of all supported filters
- And more...

### Hash Functions (`pkg/stdlib/hash/`)
Hashing functions.

**Planned Functions:**
- `md5`, `sha1`, `hash`, `hash_file`
- `hash_hmac`, `hash_pbkdf2`
- `password_hash`, `password_verify`, `password_needs_rehash`
- And more...

### SPL Functions (`pkg/stdlib/spl/`)
Standard PHP Library data structures and iterators.

**Planned Classes:**
- `SplStack`, `SplQueue`, `SplHeap`, `SplPriorityQueue`
- `SplDoublyLinkedList`, `SplFixedArray`
- `ArrayIterator`, `DirectoryIterator`, `RecursiveDirectoryIterator`
- And more...

### Reflection (`pkg/stdlib/reflection/`)
Reflection API for introspection.

**Planned Classes:**
- `ReflectionClass`, `ReflectionMethod`, `ReflectionProperty`
- `ReflectionFunction`, `ReflectionParameter`
- And more...

## Current Implementation Status

As of Phase 5 completion:
- **Core VM functions**: Basic functions like `echo`, `print`, `var_dump`, `isset`, `empty`, `count`, `strlen` are implemented in the VM's builtins
- **Standard library**: Not yet implemented
- **Phase 6 roadmap**: ~300 functions to be implemented

## Usage Example (Future)

Once implemented, stdlib functions will be used like this:

```go
package main

import (
    "github.com/krizos/php-go/pkg/stdlib/array"
    "github.com/krizos/php-go/pkg/stdlib/string"
    "github.com/krizos/php-go/pkg/types"
)

func main() {
    // Array functions
    arr := types.NewEmptyArray()
    arr.Append(types.NewInt(1))
    arr.Append(types.NewInt(2))

    count := array.Count(arr)  // Will return 2

    // String functions
    str := types.NewString("hello world")
    upper := string.Strtoupper(str)  // Will return "HELLO WORLD"
}
```

## Implementation Plan (Phase 6)

### Task 6.1: Array Functions (50 hours)
Implement core array manipulation functions.

### Task 6.2: String Functions (60 hours)
Implement string processing and manipulation.

### Task 6.3: Math Functions (20 hours)
Implement mathematical operations.

### Task 6.4: File I/O (40 hours)
Implement file system operations.

### Task 6.5: Date/Time (30 hours)
Implement date and time functions.

### Task 6.6: JSON Extension (20 hours)
Implement JSON encode/decode.

### Task 6.7: PCRE (Regular Expressions) (40 hours)
Implement regex matching and replacement.

### Task 6.8: Additional Functions (40 hours)
Implement remaining commonly-used functions.

**Total Estimated**: 300 hours for Phase 6

## Design Considerations

When implementing stdlib functions:

1. **PHP Compatibility**: Follow PHP 8.4 behavior exactly
2. **Error Handling**: Use PHP error reporting (warnings, notices, errors)
3. **Type Juggling**: Apply PHP type conversion rules
4. **Performance**: Optimize hot paths, use Go standard library where possible
5. **Testing**: Comprehensive test coverage (85%+ target)
6. **Documentation**: Document each function's parameters, return values, and behavior

## References

- PHP Manual: https://www.php.net/manual/en/
- PHP Source Code: `php-src/ext/standard/`
- Phase 6 Documentation: `/Users/krizos/code/mnohosten/php-go/docs/phases/06-stdlib/`

## Contributing

When Phase 6 begins, contributions to stdlib implementation are welcome. Each function should:
- Match PHP 8.4 behavior
- Include comprehensive tests
- Have performance benchmarks
- Be documented with examples

## Related Packages

- `pkg/types` - Provides Value, Array, Object types used by stdlib
- `pkg/vm` - VM provides built-in functions that stdlib extends
- `pkg/runtime` - Runtime context for error handling and configuration

## Future Enhancements (Beyond Phase 6)

- Native Go implementations of performance-critical functions
- Parallel versions of array functions (Phase 7)
- FFI extensions (Phase 8)
- Custom extensions API (Phase 10)
