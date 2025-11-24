# Plugin System Guide

This guide explains how to create and load dynamic Go plugins for PHP-Go extensions.

## Overview

The plugin system allows you to compile Go extensions as shared libraries (.so on Linux/macOS, .dll on Windows) and dynamically load them into PHP-Go at runtime. This enables:

- **Distribution**: Ship extensions as binary plugins without source code
- **Modularity**: Load only the extensions you need
- **Hot updates**: Potentially replace extensions without recompiling PHP-Go (with limitations)
- **Third-party extensions**: Allow others to create extensions for your PHP-Go installation

## Creating a Plugin

### Plugin Structure

A plugin must be a `main` package that exports an `Extension` symbol:

```go
// File: myextension/plugin.go
package main

import (
	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// MyExtension implements the Extension interface
type MyExtension struct {
	*goext.BaseExtension
}

// NewMyExtension creates a new instance of MyExtension
func NewMyExtension() *MyExtension {
	ext := &MyExtension{
		BaseExtension: goext.NewBaseExtension("myextension", "1.0.0"),
	}

	// Register functions
	ext.AddFunction("my_hello", myHello)
	ext.AddFunction("my_add", myAdd)

	// Register constants
	ext.AddConstant("MY_CONSTANT", types.NewString("Hello from plugin!"))

	return ext
}

// Init is called when the extension is loaded
func (e *MyExtension) Init(v *vm.VM) error {
	// Perform any initialization here
	return nil
}

// Function implementations
func myHello(v *vm.VM, args []*types.Value) (*types.Value, error) {
	name := "World"
	if len(args) > 0 {
		name = args[0].String()
	}
	return types.NewString("Hello, " + name + "!"), nil
}

func myAdd(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("my_add requires 2 arguments")
	}

	a := args[0].ToInt()
	b := args[1].ToInt()

	return types.NewInt(a + b), nil
}

// Extension is the exported symbol that PHP-Go will look for
var Extension goext.Extension = NewMyExtension()
```

### Building the Plugin

Build the plugin as a shared library:

```bash
# Linux/macOS
go build -buildmode=plugin -o myextension.so myextension/plugin.go

# Note: Windows requires different build flags
# go build -buildmode=plugin -o myextension.dll myextension/plugin.go
```

**Important Build Considerations:**

1. **Go version**: The plugin and PHP-Go must be built with the same Go version
2. **Dependencies**: The plugin must use the exact same dependency versions as PHP-Go
3. **Platform**: Plugins are platform-specific (Linux plugins won't work on macOS)

## Loading Plugins

### Using the Plugin Manager

```go
import "github.com/krizos/php-go/pkg/goext"

// Load a plugin
err := goext.LoadPlugin("/path/to/myextension.so")
if err != nil {
    log.Fatalf("Failed to load plugin: %v", err)
}

// The extension is now registered and can be loaded into a VM
```

### Loading into a VM

After loading the plugin, you can load its extension into a VM:

```go
import (
	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/vm"
)

// Create a VM
v := vm.New()

// Load the extension into the VM
extMgr := goext.GetGlobalExtensionManager()
err := extMgr.LoadIntoVM("myextension", v)
if err != nil {
    log.Fatalf("Failed to load extension into VM: %v", err)
}

// Now you can use the extension's functions in PHP code
```

### Using in PHP

Once loaded, the extension's functions are available in PHP:

```php
<?php
// Call plugin functions
echo my_hello("Alice");  // Output: Hello, Alice!
echo my_add(10, 20);     // Output: 30

// Use plugin constants
echo MY_CONSTANT;        // Output: Hello from plugin!
```

## Plugin Manager API

### Loading Plugins

```go
// Load a single plugin
err := pm.LoadPlugin("/path/to/plugin.so")

// Load and get extension name
extName, err := pm.LoadPluginByName("/path/to/plugin.so")
```

### Querying Plugins

```go
// Check if plugin is loaded
if pm.IsLoaded("/path/to/plugin.so") {
    fmt.Println("Plugin is loaded")
}

// Get plugin info
info, ok := pm.GetPlugin("/path/to/plugin.so")
if ok {
    fmt.Printf("Plugin: %s, Extension: %s\n",
        info.Path, info.Extension.Name())
}

// Get plugin by extension name
info, ok := pm.GetPluginByExtension("myextension")

// List all loaded plugins
for _, path := range pm.List() {
    fmt.Println(path)
}

// Count loaded plugins
count := pm.Count()
```

### Unloading Plugins

```go
// Unload a plugin
err := pm.UnloadPlugin("/path/to/plugin.so")

// Clear all plugins
pm.Clear()
```

**Note:** Go's plugin system has limitations - plugins cannot be truly unloaded from memory. The `UnloadPlugin` method only removes the plugin from tracking and unregisters its extension.

## Complete Example

### 1. Create the Plugin

```go
// File: examples/plugin/mathext/plugin.go
package main

import (
	"fmt"
	"math"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

type MathExtension struct {
	*goext.BaseExtension
}

func NewMathExtension() *MathExtension {
	ext := &MathExtension{
		BaseExtension: goext.NewBaseExtension("mathext", "1.0.0"),
	}

	ext.AddFunction("math_sqrt", mathSqrt)
	ext.AddFunction("math_pow", mathPow)
	ext.AddConstant("MATH_PI", types.NewFloat(math.Pi))

	return ext
}

func mathSqrt(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("math_sqrt requires 1 argument")
	}

	x := args[0].ToFloat()
	result := math.Sqrt(x)

	return types.NewFloat(result), nil
}

func mathPow(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("math_pow requires 2 arguments")
	}

	base := args[0].ToFloat()
	exp := args[1].ToFloat()
	result := math.Pow(base, exp)

	return types.NewFloat(result), nil
}

var Extension goext.Extension = NewMathExtension()
```

### 2. Build the Plugin

```bash
cd examples/plugin/mathext
go build -buildmode=plugin -o mathext.so plugin.go
```

### 3. Use the Plugin

```go
// File: examples/plugin/main.go
package main

import (
	"fmt"
	"log"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/vm"
)

func main() {
	// Load the plugin
	err := goext.LoadPlugin("./mathext/mathext.so")
	if err != nil {
		log.Fatalf("Failed to load plugin: %v", err)
	}

	// Create VM
	v := vm.New()

	// Load extension into VM
	extMgr := goext.GetGlobalExtensionManager()
	err = extMgr.LoadIntoVM("mathext", v)
	if err != nil {
		log.Fatalf("Failed to load extension: %v", err)
	}

	// Now run PHP code that uses the extension
	code := `<?php
		echo "Square root of 16: " . math_sqrt(16) . "\n";
		echo "2 to the power of 8: " . math_pow(2, 8) . "\n";
		echo "Value of PI: " . MATH_PI . "\n";
	?>`

	// Execute the code
	// ... (execute code with VM)
}
```

## Best Practices

### 1. Version Compatibility

Always document the required PHP-Go version:

```go
const RequiredPHPGoVersion = "0.1.0"

func (e *MyExtension) Init(v *vm.VM) error {
	// Check version compatibility
	// if !isCompatible(v.Version(), RequiredPHPGoVersion) {
	//     return fmt.Errorf("requires PHP-Go %s or later", RequiredPHPGoVersion)
	// }
	return nil
}
```

### 2. Error Handling

Always provide clear error messages:

```go
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("myFunc() expects at least 2 arguments, %d given", len(args))
	}

	// Validate argument types
	if args[0].Type != types.TypeInt {
		return nil, fmt.Errorf("myFunc() argument 1 must be an integer, %s given", args[0].TypeName())
	}

	// ... function logic
}
```

### 3. Thread Safety

Ensure your plugin is thread-safe if used with parallel features:

```go
type MyExtension struct {
	*goext.BaseExtension
	mu    sync.RWMutex
	cache map[string]interface{}
}

func (e *MyExtension) Get(key string) interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.cache[key]
}
```

### 4. Resource Cleanup

Clean up resources in your extension:

```go
func (e *MyExtension) Init(v *vm.VM) error {
	// Initialize resources
	e.connection = openConnection()

	// Register cleanup handler
	// v.RegisterShutdownHandler(func() {
	//     e.connection.Close()
	// })

	return nil
}
```

### 5. Documentation

Document your functions using PHP-style docblocks:

```go
// myFunction implements the PHP function my_function()
//
// PHP Signature: string my_function(string $input, int $flags = 0)
//
// Returns the processed input string with optional flags.
func myFunction(v *vm.VM, args []*types.Value) (*types.Value, error) {
	// ...
}
```

## Limitations

### Go Plugin System Limitations

1. **No unloading**: Plugins cannot be truly unloaded from memory
2. **Version sensitivity**: Plugin and main program must use same Go version
3. **Dependency sensitivity**: All dependencies must match exactly
4. **Platform-specific**: Plugins are not cross-platform
5. **No Windows support**: Limited or no support on Windows (check Go version)

### Workarounds

For hot reloading or dynamic updates:

1. **Process isolation**: Run plugins in separate processes and communicate via IPC
2. **Embedded extensions**: For critical extensions, compile them directly into PHP-Go
3. **Configuration-based loading**: Use config files to control which extensions to load at startup

## Troubleshooting

### Plugin fails to load

```
Error: plugin.Open: plugin was built with a different version of package X
```

**Solution**: Rebuild the plugin with the same Go version and dependency versions as PHP-Go.

### Symbol not found

```
Error: plugin 'myext.so' does not export 'Extension' symbol
```

**Solution**: Ensure your plugin exports the `Extension` variable:

```go
var Extension goext.Extension = NewMyExtension()
```

### Type assertion failed

```
Error: symbol 'Extension' is not of type goext.Extension
```

**Solution**: Make sure your extension implements the `goext.Extension` interface:

```go
func (e *MyExtension) Name() string
func (e *MyExtension) Version() string
func (e *MyExtension) Init(*vm.VM) error
func (e *MyExtension) Functions() map[string]FunctionHandler
func (e *MyExtension) Constants() map[string]*types.Value
```

## Next Steps

- Learn about [Advanced Marshaling](03-advanced-marshaling.md) for complex type conversions
- Explore [Extension Development](01-extension-development.md) for in-depth extension creation
- Review [Type Marshaling](02-type-marshaling.md) for PHP ↔ Go conversions
