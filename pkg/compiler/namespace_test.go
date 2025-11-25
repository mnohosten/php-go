package compiler

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
)

// TestNamespaceTracking verifies that namespace context is properly tracked
func TestNamespaceTracking(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedNamespace string
	}{
		{
			name:              "global namespace (no declaration)",
			input:             `<?php echo "hello";`,
			expectedNamespace: "",
		},
		{
			name:              "simple namespace",
			input:             `<?php namespace MyNamespace;`,
			expectedNamespace: "MyNamespace",
		},
		{
			name:              "nested namespace",
			input:             `<?php namespace Vendor\Package\SubPackage;`,
			expectedNamespace: `Vendor\Package\SubPackage`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			c := New()
			err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			if c.GetCurrentNamespace() != tt.expectedNamespace {
				t.Errorf("Expected namespace %q, got %q", tt.expectedNamespace, c.GetCurrentNamespace())
			}
		})
	}
}

// TestUseImportTracking verifies that use imports are properly registered
func TestUseImportTracking(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedImports map[string]string
	}{
		{
			name: "simple use",
			input: `<?php
			use MyNamespace\MyClass;
			`,
			expectedImports: map[string]string{
				"MyClass": `MyNamespace\MyClass`,
			},
		},
		{
			name: "use with alias",
			input: `<?php
			use MyNamespace\MyClass as Alias;
			`,
			expectedImports: map[string]string{
				"Alias": `MyNamespace\MyClass`,
			},
		},
		{
			name: "multiple use",
			input: `<?php
			use Namespace\ClassA;
			use Namespace\ClassB;
			`,
			expectedImports: map[string]string{
				"ClassA": `Namespace\ClassA`,
				"ClassB": `Namespace\ClassB`,
			},
		},
		{
			name: "grouped use",
			input: `<?php
			use Vendor\Package\{ClassA, ClassB, ClassC};
			`,
			expectedImports: map[string]string{
				"ClassA": `Vendor\Package\ClassA`,
				"ClassB": `Vendor\Package\ClassB`,
				"ClassC": `Vendor\Package\ClassC`,
			},
		},
		{
			name: "grouped use with alias",
			input: `<?php
			use Vendor\Package\{ClassA as A, ClassB};
			`,
			expectedImports: map[string]string{
				"A":      `Vendor\Package\ClassA`,
				"ClassB": `Vendor\Package\ClassB`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			c := New()
			err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			imports := c.GetUseImports()
			if len(imports) != len(tt.expectedImports) {
				t.Errorf("Expected %d imports, got %d: %v", len(tt.expectedImports), len(imports), imports)
			}

			for alias, expectedFQN := range tt.expectedImports {
				if fqn, ok := imports[alias]; !ok {
					t.Errorf("Missing import for alias %q", alias)
				} else if fqn != expectedFQN {
					t.Errorf("Import %q: expected FQN %q, got %q", alias, expectedFQN, fqn)
				}
			}
		})
	}
}

// TestFunctionUseImports verifies function use imports
func TestFunctionUseImports(t *testing.T) {
	input := `<?php
	use function Vendor\Package\myFunction;
	use function Vendor\Other\otherFunc as aliasFunc;
	`

	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	funcImports := c.GetUseFunctionImports()
	expected := map[string]string{
		"myFunction": `Vendor\Package\myFunction`,
		"aliasFunc":  `Vendor\Other\otherFunc`,
	}

	for alias, expectedFQN := range expected {
		if fqn, ok := funcImports[alias]; !ok {
			t.Errorf("Missing function import for alias %q", alias)
		} else if fqn != expectedFQN {
			t.Errorf("Function import %q: expected FQN %q, got %q", alias, expectedFQN, fqn)
		}
	}
}

// TestConstUseImports verifies constant use imports
func TestConstUseImports(t *testing.T) {
	input := `<?php
	use const Vendor\Package\MY_CONST;
	use const Vendor\Other\OTHER_CONST as ALIAS_CONST;
	`

	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	constImports := c.GetUseConstImports()
	expected := map[string]string{
		"MY_CONST":    `Vendor\Package\MY_CONST`,
		"ALIAS_CONST": `Vendor\Other\OTHER_CONST`,
	}

	for alias, expectedFQN := range expected {
		if fqn, ok := constImports[alias]; !ok {
			t.Errorf("Missing const import for alias %q", alias)
		} else if fqn != expectedFQN {
			t.Errorf("Const import %q: expected FQN %q, got %q", alias, expectedFQN, fqn)
		}
	}
}

// TestResolveClassName verifies class name resolution
func TestResolveClassName(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(c *Compiler)
		className        string
		expectedResolved string
	}{
		{
			name:             "unqualified in global namespace",
			setup:            func(c *Compiler) {},
			className:        "MyClass",
			expectedResolved: "MyClass",
		},
		{
			name: "unqualified in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
			},
			className:        "MyClass",
			expectedResolved: `App\MyClass`,
		},
		{
			name: "unqualified with use import",
			setup: func(c *Compiler) {
				c.useImports["MyClass"] = `Vendor\Package\MyClass`
			},
			className:        "MyClass",
			expectedResolved: `Vendor\Package\MyClass`,
		},
		{
			name: "unqualified with use import in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
				c.useImports["MyClass"] = `Vendor\Package\MyClass`
			},
			className:        "MyClass",
			expectedResolved: `Vendor\Package\MyClass`,
		},
		{
			name:             "fully qualified (starts with \\)",
			setup:            func(c *Compiler) {},
			className:        `\Vendor\Package\MyClass`,
			expectedResolved: `Vendor\Package\MyClass`,
		},
		{
			name: "fully qualified in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
			},
			className:        `\Vendor\Package\MyClass`,
			expectedResolved: `Vendor\Package\MyClass`,
		},
		{
			name: "qualified name in global namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = ""
			},
			className:        `Sub\MyClass`,
			expectedResolved: `Sub\MyClass`,
		},
		{
			name: "qualified name in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
			},
			className:        `Sub\MyClass`,
			expectedResolved: `App\Sub\MyClass`,
		},
		{
			name: "qualified name with first part imported",
			setup: func(c *Compiler) {
				c.useImports["Sub"] = `Vendor\Package\Sub`
			},
			className:        `Sub\MyClass`,
			expectedResolved: `Vendor\Package\Sub\MyClass`,
		},
		{
			name: "qualified name with first part imported in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
				c.useImports["Sub"] = `Vendor\Package\Sub`
			},
			className:        `Sub\MyClass`,
			expectedResolved: `Vendor\Package\Sub\MyClass`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			tt.setup(c)

			resolved := c.ResolveClassName(tt.className)
			if resolved != tt.expectedResolved {
				t.Errorf("ResolveClassName(%q) = %q, expected %q", tt.className, resolved, tt.expectedResolved)
			}
		})
	}
}

// TestResolveFunctionName verifies function name resolution
func TestResolveFunctionName(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(c *Compiler)
		funcName         string
		expectedResolved string
	}{
		{
			name:             "unqualified in global namespace",
			setup:            func(c *Compiler) {},
			funcName:         "myFunc",
			expectedResolved: "myFunc",
		},
		{
			name: "unqualified in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
			},
			funcName:         "myFunc",
			expectedResolved: `App\myFunc`,
		},
		{
			name: "with function use import",
			setup: func(c *Compiler) {
				c.useFunctionImports["myFunc"] = `Vendor\Package\myFunc`
			},
			funcName:         "myFunc",
			expectedResolved: `Vendor\Package\myFunc`,
		},
		{
			name:             "fully qualified",
			setup:            func(c *Compiler) {},
			funcName:         `\strlen`,
			expectedResolved: "strlen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			tt.setup(c)

			resolved := c.ResolveFunctionName(tt.funcName)
			if resolved != tt.expectedResolved {
				t.Errorf("ResolveFunctionName(%q) = %q, expected %q", tt.funcName, resolved, tt.expectedResolved)
			}
		})
	}
}

// TestResolveConstantName verifies constant name resolution
func TestResolveConstantName(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(c *Compiler)
		constName        string
		expectedResolved string
	}{
		{
			name:             "unqualified in global namespace",
			setup:            func(c *Compiler) {},
			constName:        "MY_CONST",
			expectedResolved: "MY_CONST",
		},
		{
			name: "unqualified in namespace",
			setup: func(c *Compiler) {
				c.currentNamespace = "App"
			},
			constName:        "MY_CONST",
			expectedResolved: `App\MY_CONST`,
		},
		{
			name: "with const use import",
			setup: func(c *Compiler) {
				c.useConstImports["MY_CONST"] = `Vendor\Package\MY_CONST`
			},
			constName:        "MY_CONST",
			expectedResolved: `Vendor\Package\MY_CONST`,
		},
		{
			name:             "fully qualified",
			setup:            func(c *Compiler) {},
			constName:        `\PHP_VERSION`,
			expectedResolved: "PHP_VERSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			tt.setup(c)

			resolved := c.ResolveConstantName(tt.constName)
			if resolved != tt.expectedResolved {
				t.Errorf("ResolveConstantName(%q) = %q, expected %q", tt.constName, resolved, tt.expectedResolved)
			}
		})
	}
}

// TestBracketedNamespaceRestoresContext verifies that bracketed namespaces restore context
func TestBracketedNamespaceRestoresContext(t *testing.T) {
	input := `<?php
	namespace App {
		// Inside App namespace
	}
	// After bracketed namespace, should be back to global
	`

	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	// After a bracketed namespace, we should be back to global
	if c.GetCurrentNamespace() != "" {
		t.Errorf("Expected empty namespace after bracketed namespace, got %q", c.GetCurrentNamespace())
	}
}

// TestNamespaceClearsUseImports verifies that entering a namespace clears use imports
func TestNamespaceClearsUseImports(t *testing.T) {
	input := `<?php
	use Vendor\OldClass;

	namespace App;
	// use imports should be cleared here
	`

	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	// After entering a new namespace, use imports should be cleared
	imports := c.GetUseImports()
	if len(imports) != 0 {
		t.Errorf("Expected no use imports after entering namespace, got %v", imports)
	}
}

// TestClassNameExpressionWithNamespace tests ::class with namespace resolution
func TestClassNameExpressionWithNamespace(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedClassName string
	}{
		{
			name: "simple class in global namespace",
			input: `<?php
			echo MyClass::class;
			`,
			expectedClassName: "MyClass",
		},
		{
			name: "simple class in namespace",
			input: `<?php
			namespace App;
			echo MyClass::class;
			`,
			expectedClassName: `App\MyClass`,
		},
		{
			name: "imported class",
			input: `<?php
			namespace App;
			use Vendor\Package\MyClass;
			echo MyClass::class;
			`,
			expectedClassName: `Vendor\Package\MyClass`,
		},
		{
			name: "fully qualified class",
			input: `<?php
			namespace App;
			echo \Vendor\Package\MyClass::class;
			`,
			expectedClassName: `Vendor\Package\MyClass`,
		},
		{
			name: "qualified class in namespace",
			input: `<?php
			namespace App;
			echo Sub\MyClass::class;
			`,
			expectedClassName: `App\Sub\MyClass`,
		},
		{
			name: "qualified class with imported prefix",
			input: `<?php
			use Vendor\Package as VP;
			echo VP\MyClass::class;
			`,
			expectedClassName: `Vendor\Package\MyClass`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			c := New()
			err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			// Find the class name constant in the compiled constants
			found := false
			for _, constant := range c.constants {
				if strVal, ok := constant.(string); ok {
					if strVal == tt.expectedClassName {
						found = true
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected class name %q not found in constants. Constants: %v", tt.expectedClassName, c.constants)
			}
		})
	}
}

// TestClassNameExpressionSpecialKeywords tests self::class, parent::class, static::class
func TestClassNameExpressionSpecialKeywords(t *testing.T) {
	// These keywords should NOT be resolved at compile time
	// They must be passed to the VM for runtime resolution
	tests := []struct {
		name           string
		input          string
		expectedInCode string
	}{
		{
			name: "self::class",
			input: `<?php
			namespace App;
			echo self::class;
			`,
			expectedInCode: "self", // Should be "self", not resolved
		},
		{
			name: "parent::class",
			input: `<?php
			namespace App;
			echo parent::class;
			`,
			expectedInCode: "parent", // Should be "parent", not resolved
		},
		{
			name: "static::class",
			input: `<?php
			namespace App;
			echo static::class;
			`,
			expectedInCode: "static", // Should be "static", not resolved
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			c := New()
			err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			// Special keywords should be in constants as-is (not namespace-resolved)
			found := false
			for _, constant := range c.constants {
				if strVal, ok := constant.(string); ok {
					if strVal == tt.expectedInCode {
						found = true
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected keyword %q not found in constants. Constants: %v", tt.expectedInCode, c.constants)
			}

			// Also verify it was NOT prepended with namespace
			for _, constant := range c.constants {
				if strVal, ok := constant.(string); ok {
					if strVal == "App\\"+tt.expectedInCode {
						t.Errorf("Special keyword %q was incorrectly namespace-resolved to %q", tt.expectedInCode, strVal)
					}
				}
			}
		})
	}
}

// TestWordPressStyleNamespaces tests WordPress-like namespace patterns
func TestWordPressStyleNamespaces(t *testing.T) {
	input := `<?php
	namespace WP\REST\V1;

	use WP\Core\Application;
	use WP\Database\{Connection, Query};
	use function WP\Helpers\sanitize_text;
	use const WP\Config\VERSION;
	`

	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	// Check namespace
	if c.GetCurrentNamespace() != `WP\REST\V1` {
		t.Errorf("Expected namespace 'WP\\REST\\V1', got %q", c.GetCurrentNamespace())
	}

	// Check class imports
	classImports := c.GetUseImports()
	expectedClasses := map[string]string{
		"Application": `WP\Core\Application`,
		"Connection":  `WP\Database\Connection`,
		"Query":       `WP\Database\Query`,
	}
	for alias, fqn := range expectedClasses {
		if classImports[alias] != fqn {
			t.Errorf("Class import %q: expected %q, got %q", alias, fqn, classImports[alias])
		}
	}

	// Check function import
	funcImports := c.GetUseFunctionImports()
	if funcImports["sanitize_text"] != `WP\Helpers\sanitize_text` {
		t.Errorf("Function import: expected 'WP\\Helpers\\sanitize_text', got %q", funcImports["sanitize_text"])
	}

	// Check const import
	constImports := c.GetUseConstImports()
	if constImports["VERSION"] != `WP\Config\VERSION` {
		t.Errorf("Const import: expected 'WP\\Config\\VERSION', got %q", constImports["VERSION"])
	}

	// Test class resolution
	resolved := c.ResolveClassName("Application")
	if resolved != `WP\Core\Application` {
		t.Errorf("ResolveClassName(Application) = %q, expected 'WP\\Core\\Application'", resolved)
	}

	// Test unresolved class (should prepend current namespace)
	resolved = c.ResolveClassName("Controller")
	if resolved != `WP\REST\V1\Controller` {
		t.Errorf("ResolveClassName(Controller) = %q, expected 'WP\\REST\\V1\\Controller'", resolved)
	}
}
