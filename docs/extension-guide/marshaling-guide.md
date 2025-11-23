# Type Marshaling Guide

This guide explains how to convert data between PHP and Go types in the PHP-Go system.

## Table of Contents

1. [Overview](#overview)
2. [Basic Type Mapping](#basic-type-mapping)
3. [Basic Marshaling](#basic-marshaling)
4. [Advanced Marshaling](#advanced-marshaling)
5. [Custom Types](#custom-types)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

## Overview

Type marshaling converts values between PHP's dynamic type system and Go's static type system. The PHP-Go marshaler provides:

- Automatic conversion for common types
- Custom type converters for user-defined types
- Struct ↔ Object conversion with field mapping
- Array detection (list vs associative)
- Error handling for unsupported types

## Basic Type Mapping

### Primitive Types

| Go Type | PHP Type | Example |
|---------|----------|---------|
| `nil` | `null` | `null` |
| `bool` | `bool` | `true`, `false` |
| `int`, `int64` | `int` | `42` |
| `float64` | `float` | `3.14` |
| `string` | `string` | `"hello"` |

### Complex Types

| Go Type | PHP Type | Notes |
|---------|----------|-------|
| `[]interface{}` | `array` (list) | Sequential keys [0, 1, 2, ...] |
| `map[string]interface{}` | `array` (assoc) | String keys ["key" => value] |
| `struct` | `object` or `array` | Depends on context |
| `*types.Value` | Any PHP type | Native PHP value |

## Basic Marshaling

### Creating a Marshaler

```go
import "github.com/krizos/php-go/pkg/goext"

m := goext.NewMarshaler()
```

### PHP → Go Conversion

```go
// Convert any PHP value to Go
phpVal := types.NewInt(42)
goVal, err := m.ToGo(phpVal)
// goVal is int64(42)

// Type-specific conversions
phpStr := types.NewString("hello")
str, err := m.ToGoString(phpStr)  // string

phpInt := types.NewInt(100)
num, err := m.ToGoInt(phpInt)     // int64

phpFloat := types.NewFloat(3.14)
f, err := m.ToGoFloat(phpFloat)   // float64

phpBool := types.NewBool(true)
b, err := m.ToGoBool(phpBool)     // bool
```

### Go → PHP Conversion

```go
// Convert Go values to PHP
phpInt, err := m.ToPHP(int64(42))
// phpInt is *types.Value with TypeInt

phpStr, err := m.ToPHP("hello")
// phpStr is *types.Value with TypeString

phpBool, err := m.ToPHP(true)
// phpBool is *types.Value with TypeBool
```

### Arrays and Slices

```go
// Go slice → PHP array (list)
goSlice := []int64{1, 2, 3, 4, 5}
phpArr, err := m.ToPHP(goSlice)
// PHP: [1, 2, 3, 4, 5]

// Go map → PHP array (associative)
goMap := map[string]interface{}{
    "name": "John",
    "age":  30,
}
phpArr, err := m.ToPHP(goMap)
// PHP: ["name" => "John", "age" => 30]

// PHP array → Go (automatic detection)
phpList := types.NewEmptyArray()
phpList.Append(types.NewInt(1))
phpList.Append(types.NewInt(2))

goVal, err := m.ToGo(types.NewArray(phpList))
// goVal is []interface{}{int64(1), int64(2)}

// PHP associative array → Go map
phpMap := types.NewEmptyArray()
phpMap.Set(types.NewString("key"), types.NewString("value"))

goVal, err := m.ToGo(types.NewArray(phpMap))
// goVal is map[string]interface{}{"key": "value"}
```

## Advanced Marshaling

### Creating an Advanced Marshaler

```go
import "github.com/krizos/php-go/pkg/goext"

am := goext.NewAdvancedMarshaler()
```

### Struct ↔ Object Conversion

```go
// Define a Go struct
type Person struct {
    Name  string
    Age   int
    Email string `json:"email"`  // Use json tag for field naming
}

// Struct → PHP Object
person := Person{
    Name:  "John Doe",
    Age:   30,
    Email: "john@example.com",
}

phpObj, err := am.StructToObject(person, "Person")
// Creates PHP object with class name "Person"
// Properties: Name, Age, email (note: uses json tag)

// PHP Object/Array → Struct
phpArr := types.NewEmptyArray()
phpArr.Set(types.NewString("Name"), types.NewString("Jane"))
phpArr.Set(types.NewString("Age"), types.NewInt(25))
phpArr.Set(types.NewString("email"), types.NewString("jane@example.com"))

var result Person
goVal, err := am.phpToStruct(types.NewArray(phpArr), reflect.TypeOf(Person{}))
result = goVal.(Person)
// result.Name = "Jane"
// result.Age = 25
// result.Email = "jane@example.com"
```

### Nested Structures

```go
type Address struct {
    Street  string
    City    string
    ZipCode string `json:"zip_code"`
}

type PersonWithAddress struct {
    Name    string
    Address Address
}

// Convert nested struct
person := PersonWithAddress{
    Name: "John",
    Address: Address{
        Street:  "123 Main St",
        City:    "Springfield",
        ZipCode: "12345",
    },
}

phpObj, err := am.StructToObject(person, "PersonWithAddress")
// PHP object with nested Address object
```

### Custom Type Converters

```go
// Register a custom converter for a specific type
import "time"
import "reflect"

am := goext.NewAdvancedMarshaler()

// Register time.Time converter
am.RegisterConverter(reflect.TypeOf(time.Time{}), &goext.TypeConverter{
    ToPHP: func(v interface{}) (*types.Value, error) {
        t := v.(time.Time)
        return types.NewInt(t.Unix()), nil
    },
    ToGo: func(v *types.Value) (interface{}, error) {
        return time.Unix(v.ToInt(), 0), nil
    },
})

// Now time.Time converts automatically
now := time.Now()
phpTimestamp, err := am.ToPHPAdvanced(now)
// phpTimestamp is Unix timestamp

// And back
goTime, err := am.ToGoAdvanced(phpTimestamp, reflect.TypeOf(time.Time{}))
// goTime is time.Time
```

## Custom Types

### Implementing CustomMarshaler

```go
type MyType struct {
    Value int
    Label string
}

// Implement MarshalPHP
func (m *MyType) MarshalPHP() (*types.Value, error) {
    arr := types.NewEmptyArray()
    arr.Set(types.NewString("value"), types.NewInt(int64(m.Value)))
    arr.Set(types.NewString("label"), types.NewString(m.Label))
    return types.NewArray(arr), nil
}

// Now MyType marshals automatically
mt := &MyType{Value: 42, Label: "answer"}
phpVal, err := am.ToPHPAdvanced(mt)
// PHP: ["value" => 42, "label" => "answer"]
```

### Implementing CustomUnmarshaler

```go
func (m *MyType) UnmarshalPHP(v *types.Value) error {
    if v.Type() != types.TypeArray {
        return fmt.Errorf("expected array")
    }

    arr := v.ToArray()

    if val, ok := arr.Get(types.NewString("value")); ok {
        m.Value = int(val.ToInt())
    }

    if lbl, ok := arr.Get(types.NewString("label")); ok {
        m.Label = lbl.ToString()
    }

    return nil
}

// Unmarshal from PHP
phpArr := types.NewEmptyArray()
phpArr.Set(types.NewString("value"), types.NewInt(42))
phpArr.Set(types.NewString("label"), types.NewString("answer"))

var mt MyType
err := mt.UnmarshalPHP(types.NewArray(phpArr))
// mt.Value = 42, mt.Label = "answer"
```

### Type Converter vs Custom Marshaler

Use **TypeConverter** when:
- You don't control the type definition
- You want centralized conversion logic
- You need different conversions in different contexts

Use **CustomMarshaler interface** when:
- You control the type definition
- Conversion is intrinsic to the type
- You want the simplest approach

## Best Practices

### 1. Use Type-Specific Conversions

```go
// Prefer type-specific methods
str, err := m.ToGoString(phpVal)

// Over generic conversion with type assertion
val, err := m.ToGo(phpVal)
str := val.(string)  // Can panic if wrong type
```

### 2. Handle Errors

```go
goVal, err := m.ToGo(phpVal)
if err != nil {
    return fmt.Errorf("marshaling failed: %w", err)
}
```

### 3. Use JSON Tags for Struct Fields

```go
type User struct {
    UserID    int    `json:"user_id"`     // Maps to "user_id" in PHP
    FirstName string `json:"first_name"`  // Maps to "first_name" in PHP
    LastName  string `json:"last_name"`   // Maps to "last_name" in PHP
}
```

### 4. Validate Types

```go
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf("expected string, got %s",
            args[0].TypeString())
    }

    str := args[0].ToString()
    // Use str safely
}
```

### 5. Batch Conversions

```go
// Batch convert arrays
phpArray := types.NewEmptyArray()
for _, item := range goSlice {
    phpVal, _ := m.ToPHP(item)
    phpArray.Append(phpVal)
}
result := types.NewArray(phpArray)
```

### 6. Cache Marshalers

```go
// Create once, reuse many times
var globalMarshaler = goext.NewMarshaler()

func convert(phpVal *types.Value) (interface{}, error) {
    return globalMarshaler.ToGo(phpVal)
}
```

## Type Conversion Examples

### Complex Structures

```go
type Config struct {
    Database struct {
        Host     string
        Port     int
        Username string
        Password string
    }
    Cache struct {
        Enabled bool
        TTL     int
    }
    Features []string
}

// Convert to PHP
config := Config{ /* ... */ }
phpVal, err := am.StructToObject(config, "Config")

// PHP receives:
// object(Config) {
//     ["Database"] => array(
//         ["Host"] => "localhost",
//         ["Port"] => 5432,
//         ...
//     ),
//     ["Cache"] => array(...),
//     ["Features"] => array("feature1", "feature2")
// }
```

### Working with interfaces{}

```go
// Mixed type slice
data := []interface{}{
    int64(42),
    "hello",
    true,
    []int64{1, 2, 3},
}

phpVal, err := m.ToPHP(data)
// PHP: [42, "hello", true, [1, 2, 3]]

// Maps with mixed values
mixed := map[string]interface{}{
    "int":    int64(42),
    "string": "hello",
    "array":  []int64{1, 2, 3},
}

phpVal, err := m.ToPHP(mixed)
// PHP: ["int" => 42, "string" => "hello", "array" => [1, 2, 3]]
```

## Troubleshooting

### Type Mismatch Errors

```
Error: cannot convert string to int64
```

**Solution**: Validate PHP types before conversion
```go
if args[0].Type() != types.TypeInt {
    return nil, fmt.Errorf("expected int")
}
```

### Unsupported Type Errors

```
Error: unsupported type: chan int
```

**Solution**: Register a custom converter or convert manually
```go
// Not all Go types can be marshaled
// Channels, functions, unsafe pointers are not supported
```

### Array vs Slice Confusion

```go
// PHP arrays with sequential keys become slices
phpArr := [0 => 1, 1 => 2, 2 => 3]  // []interface{}{1, 2, 3}

// PHP arrays with non-sequential or string keys become maps
phpArr := ["a" => 1, "b" => 2]      // map[string]interface{}
phpArr := [0 => 1, 5 => 2]          // map[string]interface{} (gap at 1-4)
```

### Nil vs Null

```go
// Go nil → PHP null
var ptr *MyType  // nil
phpVal, _ := m.ToPHP(ptr)  // PHP null

// PHP null → Go nil
phpNull := types.NewNull()
goVal, _ := m.ToGo(phpNull)  // nil
```

### Performance Issues

For high-frequency conversions:

1. **Cache type converters**
```go
var converter *goext.TypeConverter
converter = &goext.TypeConverter{...}
am.RegisterConverter(myType, converter)
```

2. **Avoid unnecessary conversions**
```go
// Bad: Convert repeatedly
for i := 0; i < 1000; i++ {
    phpVal, _ := m.ToPHP(data)
    process(phpVal)
}

// Good: Convert once
phpVal, _ := m.ToPHP(data)
for i := 0; i < 1000; i++ {
    process(phpVal)
}
```

3. **Use native types when possible**
```go
// Prefer working with *types.Value directly
func process(v *types.Value) {
    // No marshaling needed
}
```

## Complete Example

```go
// Define types
type Book struct {
    Title     string
    Author    string
    Year      int
    Available bool
    Tags      []string
}

type Library struct {
    Name  string
    Books []Book
}

// Create marshaler
am := goext.NewAdvancedMarshaler()

// Create data
library := Library{
    Name: "City Library",
    Books: []Book{
        {
            Title:     "Go Programming",
            Author:    "John Doe",
            Year:      2023,
            Available: true,
            Tags:      []string{"programming", "go"},
        },
        {
            Title:     "PHP Mastery",
            Author:    "Jane Smith",
            Year:      2022,
            Available: false,
            Tags:      []string{"programming", "php"},
        },
    },
}

// Convert to PHP
phpObj, err := am.StructToObject(library, "Library")
if err != nil {
    log.Fatal(err)
}

// In PHP, this becomes:
// object(Library) {
//     ["Name"] => "City Library",
//     ["Books"] => [
//         ["Title" => "Go Programming", "Author" => "John Doe", ...],
//         ["Title" => "PHP Mastery", "Author" => "Jane Smith", ...]
//     ]
// }
```

## Next Steps

- Review [Extension Development Guide](README.md) for using marshalers in extensions
- Check [FFI Usage Guide](ffi-guide.md) for automatic marshaling in FFI calls
- See [examples](examples/) for more complex marshaling scenarios
