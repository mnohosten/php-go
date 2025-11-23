package goext_test

import (
	"fmt"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/goext/bindings"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Example demonstrates basic type marshaling between PHP and Go
func Example_basicMarshaling() {
	m := goext.NewMarshaler()

	// Go → PHP
	phpInt, _ := m.ToPHP(int64(42))
	fmt.Println("PHP int:", phpInt.ToInt())

	phpStr, _ := m.ToPHP("hello")
	fmt.Println("PHP string:", phpStr.ToString())

	// PHP → Go
	goVal, _ := m.ToGo(types.NewInt(100))
	fmt.Println("Go value:", goVal.(int64))

	// Output:
	// PHP int: 42
	// PHP string: hello
	// Go value: 100
}

// Example demonstrates registering and calling Go functions
func Example_functionRegistry() {
	registry := goext.NewFunctionRegistry()
	testVM := vm.New()

	// Register a Go function
	addFunc := func(a, b int64) int64 {
		return a + b
	}
	registry.Register("math.Add", addFunc)

	// Call it using the registry
	regFunc, _ := registry.Get("math.Add")
	result, _ := regFunc.Handler(testVM, []*types.Value{
		types.NewInt(10),
		types.NewInt(20),
	})

	fmt.Println("Result:", result.ToInt())

	// Output:
	// Result: 30
}

// Example demonstrates using the FFI system
func Example_ffiManager() {
	registry := goext.NewFunctionRegistry()
	ffi := goext.NewFFIManager(registry)
	testVM := vm.New()

	// Register functions
	registry.Register("string.Upper", func(s string) string {
		return "HELLO"
	})

	// Call via FFI (simulates PHP: go_call('string.Upper', 'hello'))
	result, _ := ffi.GoCall(testVM, []*types.Value{
		types.NewString("string.Upper"),
		types.NewString("hello"),
	})

	fmt.Println("Uppercased:", result.ToString())

	// Output:
	// Uppercased: HELLO
}

// Example demonstrates creating a custom extension
func Example_customExtension() {
	// Create custom extension
	ext := goext.NewBaseExtension("my_extension", "1.0.0")

	// Add functions
	ext.AddFunction("greet", func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		name := args[0].ToString()
		return types.NewString("Hello, " + name + "!"), nil
	})

	// Add constants
	ext.AddConstant("MY_CONST", types.NewInt(42))

	fmt.Println("Extension:", ext.Name())
	fmt.Println("Version:", ext.Version())
	fmt.Println("Functions:", len(ext.Functions()))
	fmt.Println("Constants:", len(ext.Constants()))

	// Output:
	// Extension: my_extension
	// Version: 1.0.0
	// Functions: 1
	// Constants: 1
}

// Example demonstrates using HTTP extension
func Example_httpExtension() {
	ext := bindings.NewHTTPExtension()

	fmt.Println("Extension:", ext.Name())
	fmt.Println("Functions:", len(ext.Functions()))

	// Check function names
	funcs := ext.Functions()
	_, hasGet := funcs["go_http_get"]
	_, hasPost := funcs["go_http_post"]

	fmt.Println("Has go_http_get:", hasGet)
	fmt.Println("Has go_http_post:", hasPost)

	// Output:
	// Extension: go_http
	// Functions: 5
	// Has go_http_get: true
	// Has go_http_post: true
}

// Example demonstrates using JSON extension
func Example_jsonExtension() {
	ext := bindings.NewJSONExtension()
	testVM := vm.New()

	// Encode PHP array to JSON
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))

	jsonFunc := ext.Functions()["go_json_encode"]
	result, _ := jsonFunc(testVM, []*types.Value{types.NewArray(arr)})

	fmt.Println("JSON:", result.ToString())

	// Output:
	// JSON: {"age":30,"name":"John"}
}

// Example demonstrates using crypto extension
func Example_cryptoExtension() {
	ext := bindings.NewCryptoExtension()
	testVM := vm.New()

	// Hash a string
	hashFunc := ext.Functions()["go_hash_sha256"]
	result, _ := hashFunc(testVM, []*types.Value{types.NewString("hello")})

	fmt.Println("SHA256 length:", len(result.ToString()))

	// Base64 encode
	encodeFunc := ext.Functions()["go_base64_encode"]
	encoded, _ := encodeFunc(testVM, []*types.Value{types.NewString("test")})

	fmt.Println("Base64:", encoded.ToString())

	// Output:
	// SHA256 length: 64
	// Base64: dGVzdA==
}

// Example demonstrates advanced marshaling with custom types
func Example_advancedMarshaler() {
	m := goext.NewAdvancedMarshaler()

	// Register custom converter for a specific type
	type Point struct {
		X, Y int
	}

	m.RegisterConverter(
		nil, // Would use reflect.TypeOf(Point{}) in real code
		&goext.TypeConverter{
			ToPHP: func(v interface{}) (*types.Value, error) {
				p := v.(Point)
				arr := types.NewEmptyArray()
				arr.Set(types.NewString("x"), types.NewInt(int64(p.X)))
				arr.Set(types.NewString("y"), types.NewInt(int64(p.Y)))
				return types.NewArray(arr), nil
			},
		},
	)

	fmt.Println("Custom converter registered")

	// Output:
	// Custom converter registered
}

// Example demonstrates complete integration - registering and using an extension
func Example_completeIntegration() {
	// Create extension manager
	registry := goext.NewFunctionRegistry()
	extManager := goext.NewExtensionManager(registry)

	// Create and register an extension
	mathExt := goext.NewBaseExtension("math", "1.0.0")
	mathExt.AddFunction("multiply", func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		a := args[0].ToInt()
		b := args[1].ToInt()
		return types.NewInt(a * b), nil
	})
	mathExt.AddConstant("PI", types.NewFloat(3.14159))

	extManager.Register(mathExt)

	// Load into VM
	testVM := vm.New()
	extManager.LoadIntoVM("math", testVM)

	// Check constant was loaded
	pi, ok := testVM.GetGlobal("PI")
	if ok {
		fmt.Printf("PI constant: %.2f\n", pi.ToFloat())
	}

	// Extension is now loaded and ready to use
	info, _ := extManager.Get("math")
	fmt.Println("Extension loaded:", info.Loaded)

	// Output:
	// PI constant: 3.14
	// Extension loaded: true
}
