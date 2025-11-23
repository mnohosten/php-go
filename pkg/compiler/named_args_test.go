package compiler

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/types"
)

// Helper to create a test argument
func createTestArgument(name string, value ast.Expr) *ast.Argument {
	return &ast.Argument{
		Token: lexer.Token{Type: lexer.IDENT},
		Name:  name,
		Value: value,
	}
}

// Helper to create a test identifier expression
func createTestIdentifier(name string) *ast.Identifier {
	return &ast.Identifier{
		Token: lexer.Token{Type: lexer.IDENT, Literal: name},
		Value: name,
	}
}

// Helper to create a test parameter
func createTestParam(name string, hasDefault bool, isVariadic bool) *types.ParameterDef {
	return &types.ParameterDef{
		Name:       name,
		HasDefault: hasDefault,
		IsVariadic: isVariadic,
	}
}

func TestNamedArgumentResolver_Creation(t *testing.T) {
	resolver := NewNamedArgumentResolver()
	if resolver == nil {
		t.Fatal("NewNamedArgumentResolver() returned nil")
	}
	if resolver.parameterCache == nil {
		t.Error("parameterCache not initialized")
	}
}

func TestNamedArgumentResolver_RegisterFunction(t *testing.T) {
	resolver := NewNamedArgumentResolver()
	params := []string{"x", "y", "z"}
	resolver.RegisterFunction("testFunc", params)

	cached, ok := resolver.parameterCache["testFunc"]
	if !ok {
		t.Fatal("Function not registered")
	}
	if len(cached) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(cached))
	}
	if cached[0] != "x" || cached[1] != "y" || cached[2] != "z" {
		t.Errorf("Parameters not stored correctly: %v", cached)
	}
}

func TestNamedArgumentResolver_RegisterMethod(t *testing.T) {
	resolver := NewNamedArgumentResolver()
	params := []string{"a", "b"}
	resolver.RegisterMethod("MyClass", "myMethod", params)

	cached, ok := resolver.parameterCache["MyClass::myMethod"]
	if !ok {
		t.Fatal("Method not registered")
	}
	if len(cached) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(cached))
	}
}

func TestResolveArguments_AllPositional(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
		createTestParam("z", false, false),
	}

	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")),
		createTestArgument("", createTestIdentifier("arg2")),
		createTestArgument("", createTestIdentifier("arg3")),
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	if len(resolved) != 3 {
		t.Errorf("Expected 3 resolved arguments, got %d", len(resolved))
	}

	// Should be in same order
	for i, arg := range resolved {
		if arg != args[i] {
			t.Errorf("Argument at position %d changed", i)
		}
	}
}

func TestResolveArguments_AllNamed(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
		createTestParam("z", false, false),
	}

	// Arguments in reverse order
	args := []*ast.Argument{
		createTestArgument("z", createTestIdentifier("arg3")),
		createTestArgument("y", createTestIdentifier("arg2")),
		createTestArgument("x", createTestIdentifier("arg1")),
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	if len(resolved) != 3 {
		t.Errorf("Expected 3 resolved arguments, got %d", len(resolved))
	}

	// Should be reordered to match parameters
	if resolved[0].Value.(*ast.Identifier).Value != "arg1" {
		t.Errorf("First argument should be arg1, got %v", resolved[0].Value)
	}
	if resolved[1].Value.(*ast.Identifier).Value != "arg2" {
		t.Errorf("Second argument should be arg2, got %v", resolved[1].Value)
	}
	if resolved[2].Value.(*ast.Identifier).Value != "arg3" {
		t.Errorf("Third argument should be arg3, got %v", resolved[2].Value)
	}
}

func TestResolveArguments_MixedPositionalAndNamed(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
		createTestParam("z", false, false),
	}

	// First argument positional, then named
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")), // x
		createTestArgument("z", createTestIdentifier("arg3")), // z (named)
		createTestArgument("y", createTestIdentifier("arg2")), // y (named)
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	if len(resolved) != 3 {
		t.Errorf("Expected 3 resolved arguments, got %d", len(resolved))
	}

	if resolved[0].Value.(*ast.Identifier).Value != "arg1" {
		t.Errorf("First argument (x) should be arg1")
	}
	if resolved[1].Value.(*ast.Identifier).Value != "arg2" {
		t.Errorf("Second argument (y) should be arg2")
	}
	if resolved[2].Value.(*ast.Identifier).Value != "arg3" {
		t.Errorf("Third argument (z) should be arg3")
	}
}

func TestResolveArguments_WithDefaults(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", true, false), // has default
		createTestParam("z", true, false), // has default
	}

	// Only provide first argument
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")),
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	// Should only have the provided argument
	if len(resolved) != 1 {
		t.Errorf("Expected 1 resolved argument, got %d", len(resolved))
	}
}

func TestResolveArguments_SkipMiddleParameter(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", true, false), // has default
		createTestParam("z", false, false),
	}

	// Skip middle parameter (y) using named argument
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")),  // x
		createTestArgument("z", createTestIdentifier("arg3")), // z (named)
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	if len(resolved) != 3 {
		t.Errorf("Expected 3 resolved arguments, got %d", len(resolved))
	}

	if resolved[0].Value.(*ast.Identifier).Value != "arg1" {
		t.Errorf("First argument should be arg1")
	}
	if resolved[1] != nil {
		t.Errorf("Middle argument should be nil (uses default)")
	}
	if resolved[2].Value.(*ast.Identifier).Value != "arg3" {
		t.Errorf("Third argument should be arg3")
	}
}

func TestResolveArguments_Error_PositionalAfterNamed(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
	}

	// Positional after named - should error
	args := []*ast.Argument{
		createTestArgument("y", createTestIdentifier("arg2")),
		createTestArgument("", createTestIdentifier("arg1")), // ERROR
	}

	_, err := resolver.ResolveArguments("test", args, params)
	if err == nil {
		t.Error("Expected error for positional argument after named argument")
	}
}

func TestResolveArguments_Error_UnknownParameter(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
	}

	args := []*ast.Argument{
		createTestArgument("unknown", createTestIdentifier("arg1")),
	}

	_, err := resolver.ResolveArguments("test", args, params)
	if err == nil {
		t.Error("Expected error for unknown parameter name")
	}
}

func TestResolveArguments_Error_MissingRequiredParameter(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false), // required, no default
	}

	// Only provide first argument
	args := []*ast.Argument{
		createTestArgument("x", createTestIdentifier("arg1")),
	}

	_, err := resolver.ResolveArguments("test", args, params)
	if err == nil {
		t.Error("Expected error for missing required parameter")
	}
}

func TestResolveArguments_Error_DuplicateParameter(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("y", false, false),
	}

	// Provide x positionally, then again by name
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")), // x positional
		createTestArgument("x", createTestIdentifier("arg2")), // x named - ERROR
	}

	_, err := resolver.ResolveArguments("test", args, params)
	if err == nil {
		t.Error("Expected error for duplicate parameter")
	}
}

func TestValidateArguments_Valid(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")),
		createTestArgument("y", createTestIdentifier("arg2")),
		createTestArgument("z", createTestIdentifier("arg3")),
	}

	err := resolver.ValidateArguments("test", args)
	if err != nil {
		t.Errorf("ValidateArguments() error = %v", err)
	}
}

func TestValidateArguments_Invalid(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	args := []*ast.Argument{
		createTestArgument("y", createTestIdentifier("arg2")),
		createTestArgument("", createTestIdentifier("arg1")), // ERROR
	}

	err := resolver.ValidateArguments("test", args)
	if err == nil {
		t.Error("Expected validation error")
	}
}

func TestGetParameterNames(t *testing.T) {
	params := []*types.ParameterDef{
		createTestParam("foo", false, false),
		createTestParam("bar", false, false),
		createTestParam("baz", false, false),
	}

	names := GetParameterNames(params)
	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}
	if names[0] != "foo" || names[1] != "bar" || names[2] != "baz" {
		t.Errorf("Parameter names not extracted correctly: %v", names)
	}
}

func TestHasNamedArguments(t *testing.T) {
	tests := []struct {
		name string
		args []*ast.Argument
		want bool
	}{
		{
			name: "no named arguments",
			args: []*ast.Argument{
				createTestArgument("", createTestIdentifier("arg1")),
				createTestArgument("", createTestIdentifier("arg2")),
			},
			want: false,
		},
		{
			name: "has named arguments",
			args: []*ast.Argument{
				createTestArgument("", createTestIdentifier("arg1")),
				createTestArgument("y", createTestIdentifier("arg2")),
			},
			want: true,
		},
		{
			name: "all named arguments",
			args: []*ast.Argument{
				createTestArgument("x", createTestIdentifier("arg1")),
				createTestArgument("y", createTestIdentifier("arg2")),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasNamedArguments(tt.args)
			if got != tt.want {
				t.Errorf("HasNamedArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountPositionalArguments(t *testing.T) {
	tests := []struct {
		name string
		args []*ast.Argument
		want int
	}{
		{
			name: "all positional",
			args: []*ast.Argument{
				createTestArgument("", createTestIdentifier("arg1")),
				createTestArgument("", createTestIdentifier("arg2")),
				createTestArgument("", createTestIdentifier("arg3")),
			},
			want: 3,
		},
		{
			name: "mixed",
			args: []*ast.Argument{
				createTestArgument("", createTestIdentifier("arg1")),
				createTestArgument("", createTestIdentifier("arg2")),
				createTestArgument("z", createTestIdentifier("arg3")),
			},
			want: 2,
		},
		{
			name: "all named",
			args: []*ast.Argument{
				createTestArgument("x", createTestIdentifier("arg1")),
				createTestArgument("y", createTestIdentifier("arg2")),
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountPositionalArguments(tt.args)
			if got != tt.want {
				t.Errorf("CountPositionalArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveArguments_VariadicParameter(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	params := []*types.ParameterDef{
		createTestParam("x", false, false),
		createTestParam("args", false, true), // variadic
	}

	// Only provide required parameter
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("arg1")),
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	if len(resolved) != 1 {
		t.Errorf("Expected 1 resolved argument, got %d", len(resolved))
	}
}

func TestResolveArguments_ComplexScenario(t *testing.T) {
	resolver := NewNamedArgumentResolver()

	// function test($a, $b = 10, $c = 20, $d, $e = 50)
	params := []*types.ParameterDef{
		createTestParam("a", false, false), // required
		createTestParam("b", true, false),  // has default
		createTestParam("c", true, false),  // has default
		createTestParam("d", false, false), // required
		createTestParam("e", true, false),  // has default
	}

	// test(1, d: 4, c: 3)
	// Should provide a=1, skip b (default), c=3, d=4, skip e (default)
	args := []*ast.Argument{
		createTestArgument("", createTestIdentifier("val_a")),
		createTestArgument("d", createTestIdentifier("val_d")),
		createTestArgument("c", createTestIdentifier("val_c")),
	}

	resolved, err := resolver.ResolveArguments("test", args, params)
	if err != nil {
		t.Fatalf("ResolveArguments() error = %v", err)
	}

	// Should have a, nil, c, d, nil (but trailing nils removed)
	if len(resolved) != 4 {
		t.Errorf("Expected 4 resolved arguments, got %d", len(resolved))
	}

	if resolved[0].Value.(*ast.Identifier).Value != "val_a" {
		t.Errorf("Parameter a should be val_a")
	}
	if resolved[1] != nil {
		t.Errorf("Parameter b should be nil (uses default)")
	}
	if resolved[2].Value.(*ast.Identifier).Value != "val_c" {
		t.Errorf("Parameter c should be val_c")
	}
	if resolved[3].Value.(*ast.Identifier).Value != "val_d" {
		t.Errorf("Parameter d should be val_d")
	}
}
