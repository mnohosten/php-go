package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// Test scalar types

func TestScalarTypes(t *testing.T) {
	tests := []struct {
		input    string
		typeName string
	}{
		{`<?php function test(int $x) {}`, "int"},
		{`<?php function test(string $x) {}`, "string"},
		{`<?php function test(bool $x) {}`, "bool"},
		{`<?php function test(float $x) {}`, "float"},
		{`<?php function test(array $x) {}`, "array"},
		{`<?php function test(object $x) {}`, "object"},
		{`<?php function test(callable $x) {}`, "callable"},
		{`<?php function test(iterable $x) {}`, "iterable"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
		param := funcDecl.Parameters[0]

		if param.Type == nil {
			t.Errorf("parameter type is nil for input: %s", tt.input)
			continue
		}

		typeIdent, ok := param.Type.(*ast.Identifier)
		if !ok {
			t.Errorf("type is not *ast.Identifier. got=%T", param.Type)
			continue
		}

		if typeIdent.Value != tt.typeName {
			t.Errorf("type name not '%s'. got=%s", tt.typeName, typeIdent.Value)
		}
	}
}

func TestScalarTypeAliases(t *testing.T) {
	tests := []struct {
		input    string
		typeName string
	}{
		{`<?php function test(integer $x) {}`, "integer"}, // alias for int
		{`<?php function test(boolean $x) {}`, "boolean"}, // alias for bool
		{`<?php function test(double $x) {}`, "double"},   // alias for float
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
		param := funcDecl.Parameters[0]

		typeIdent, ok := param.Type.(*ast.Identifier)
		if !ok {
			t.Errorf("type is not *ast.Identifier. got=%T", param.Type)
			continue
		}

		if typeIdent.Value != tt.typeName {
			t.Errorf("type name not '%s'. got=%s", tt.typeName, typeIdent.Value)
		}
	}
}

// Test special types

func TestSpecialTypes(t *testing.T) {
	tests := []struct {
		input    string
		typeName string
	}{
		{`<?php function test(mixed $x) {}`, "mixed"},
		{`<?php function test(): void {}`, "void"},
		{`<?php function test(): never {}`, "never"},
		{`<?php function test(null $x) {}`, "null"},
		{`<?php function test(): static {}`, "static"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

		var typeExpr ast.Expr
		if tt.typeName == "void" || tt.typeName == "never" || tt.typeName == "static" {
			// These are return types
			typeExpr = funcDecl.ReturnType
		} else {
			// These are parameter types
			typeExpr = funcDecl.Parameters[0].Type
		}

		if typeExpr == nil {
			t.Errorf("type is nil for input: %s", tt.input)
			continue
		}

		typeIdent, ok := typeExpr.(*ast.Identifier)
		if !ok {
			t.Errorf("type is not *ast.Identifier. got=%T", typeExpr)
			continue
		}

		if typeIdent.Value != tt.typeName {
			t.Errorf("type name not '%s'. got=%s", tt.typeName, typeIdent.Value)
		}
	}
}

func TestSelfParentTypes(t *testing.T) {
	input := `<?php
class Test {
	public function getSelf(): self {
		return $this;
	}

	public function getParent(): parent {
		return parent::create();
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	// Check first method (self)
	method1 := classDecl.Body[0].(*ast.MethodDeclaration)
	typeIdent1, ok := method1.ReturnType.(*ast.Identifier)
	if !ok || typeIdent1.Value != "self" {
		t.Errorf("first method return type not 'self'. got=%v", method1.ReturnType)
	}

	// Check second method (parent)
	method2 := classDecl.Body[1].(*ast.MethodDeclaration)
	typeIdent2, ok := method2.ReturnType.(*ast.Identifier)
	if !ok || typeIdent2.Value != "parent" {
		t.Errorf("second method return type not 'parent'. got=%v", method2.ReturnType)
	}
}

// Test nullable types

func TestNullableType(t *testing.T) {
	input := `<?php function test(?int $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	nullableType, ok := param.Type.(*ast.NullableType)
	if !ok {
		t.Fatalf("type is not *ast.NullableType. got=%T", param.Type)
	}

	innerType, ok := nullableType.Type.(*ast.Identifier)
	if !ok {
		t.Fatalf("inner type is not *ast.Identifier. got=%T", nullableType.Type)
	}

	if innerType.Value != "int" {
		t.Errorf("inner type not 'int'. got=%s", innerType.Value)
	}
}

func TestNullableReturnType(t *testing.T) {
	input := `<?php function test(): ?string {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	nullableType, ok := funcDecl.ReturnType.(*ast.NullableType)
	if !ok {
		t.Fatalf("return type is not *ast.NullableType. got=%T", funcDecl.ReturnType)
	}

	innerType, ok := nullableType.Type.(*ast.Identifier)
	if !ok {
		t.Fatalf("inner type is not *ast.Identifier. got=%T", nullableType.Type)
	}

	if innerType.Value != "string" {
		t.Errorf("inner type not 'string'. got=%s", innerType.Value)
	}
}

func TestNullableClassType(t *testing.T) {
	input := `<?php function test(?User $user) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	nullableType, ok := param.Type.(*ast.NullableType)
	if !ok {
		t.Fatalf("type is not *ast.NullableType. got=%T", param.Type)
	}

	innerType, ok := nullableType.Type.(*ast.Identifier)
	if !ok {
		t.Fatalf("inner type is not *ast.Identifier. got=%T", nullableType.Type)
	}

	if innerType.Value != "User" {
		t.Errorf("inner type not 'User'. got=%s", innerType.Value)
	}
}

// Test union types (PHP 8.0+)

func TestUnionType(t *testing.T) {
	input := `<?php function test(int|string $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	unionType, ok := param.Type.(*ast.UnionType)
	if !ok {
		t.Fatalf("type is not *ast.UnionType. got=%T", param.Type)
	}

	if len(unionType.Types) != 2 {
		t.Errorf("expected 2 types in union. got=%d", len(unionType.Types))
	}

	// Check first type (int)
	type1, ok := unionType.Types[0].(*ast.Identifier)
	if !ok || type1.Value != "int" {
		t.Errorf("first type not 'int'. got=%v", unionType.Types[0])
	}

	// Check second type (string)
	type2, ok := unionType.Types[1].(*ast.Identifier)
	if !ok || type2.Value != "string" {
		t.Errorf("second type not 'string'. got=%v", unionType.Types[1])
	}
}

func TestUnionTypeMultiple(t *testing.T) {
	input := `<?php function test(int|string|float $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	unionType, ok := param.Type.(*ast.UnionType)
	if !ok {
		t.Fatalf("type is not *ast.UnionType. got=%T", param.Type)
	}

	if len(unionType.Types) != 3 {
		t.Errorf("expected 3 types in union. got=%d", len(unionType.Types))
	}

	expectedTypes := []string{"int", "string", "float"}
	for i, expected := range expectedTypes {
		typeIdent, ok := unionType.Types[i].(*ast.Identifier)
		if !ok || typeIdent.Value != expected {
			t.Errorf("type %d not '%s'. got=%v", i, expected, unionType.Types[i])
		}
	}
}

func TestUnionTypeWithNull(t *testing.T) {
	input := `<?php function test(int|null $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	unionType, ok := param.Type.(*ast.UnionType)
	if !ok {
		t.Fatalf("type is not *ast.UnionType. got=%T", param.Type)
	}

	if len(unionType.Types) != 2 {
		t.Errorf("expected 2 types in union. got=%d", len(unionType.Types))
	}
}

func TestUnionReturnType(t *testing.T) {
	input := `<?php function test(): int|string {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	unionType, ok := funcDecl.ReturnType.(*ast.UnionType)
	if !ok {
		t.Fatalf("return type is not *ast.UnionType. got=%T", funcDecl.ReturnType)
	}

	if len(unionType.Types) != 2 {
		t.Errorf("expected 2 types in union. got=%d", len(unionType.Types))
	}
}

// Test intersection types (PHP 8.1+)
// NOTE: Intersection type parsing is currently disabled in parameter position
// to avoid conflicts with by-reference parameters (array &$x)
// This will be re-enabled with proper context tracking

func TestIntersectionType(t *testing.T) {
	t.Skip("Intersection type parsing temporarily disabled - conflicts with by-reference syntax")

	input := `<?php function test(Countable&Traversable $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	intersectionType, ok := param.Type.(*ast.IntersectionType)
	if !ok {
		t.Fatalf("type is not *ast.IntersectionType. got=%T", param.Type)
	}

	if len(intersectionType.Types) != 2 {
		t.Errorf("expected 2 types in intersection. got=%d", len(intersectionType.Types))
	}

	// Check first type (Countable)
	type1, ok := intersectionType.Types[0].(*ast.Identifier)
	if !ok || type1.Value != "Countable" {
		t.Errorf("first type not 'Countable'. got=%v", intersectionType.Types[0])
	}

	// Check second type (Traversable)
	type2, ok := intersectionType.Types[1].(*ast.Identifier)
	if !ok || type2.Value != "Traversable" {
		t.Errorf("second type not 'Traversable'. got=%v", intersectionType.Types[1])
	}
}

func TestIntersectionTypeMultiple(t *testing.T) {
	t.Skip("Intersection type parsing temporarily disabled - conflicts with by-reference syntax")

	input := `<?php function test(A&B&C $x) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	intersectionType, ok := param.Type.(*ast.IntersectionType)
	if !ok {
		t.Fatalf("type is not *ast.IntersectionType. got=%T", param.Type)
	}

	if len(intersectionType.Types) != 3 {
		t.Errorf("expected 3 types in intersection. got=%d", len(intersectionType.Types))
	}

	expectedTypes := []string{"A", "B", "C"}
	for i, expected := range expectedTypes {
		typeIdent, ok := intersectionType.Types[i].(*ast.Identifier)
		if !ok || typeIdent.Value != expected {
			t.Errorf("type %d not '%s'. got=%v", i, expected, intersectionType.Types[i])
		}
	}
}

// Test class/interface type names

func TestClassTypeName(t *testing.T) {
	input := `<?php function test(User $user) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	typeIdent, ok := param.Type.(*ast.Identifier)
	if !ok {
		t.Fatalf("type is not *ast.Identifier. got=%T", param.Type)
	}

	if typeIdent.Value != "User" {
		t.Errorf("type name not 'User'. got=%s", typeIdent.Value)
	}
}

func TestFullyQualifiedClassName(t *testing.T) {
	input := `<?php function test(App\Models\User $user) {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
	param := funcDecl.Parameters[0]

	// The lexer includes the full namespace in the identifier literal
	typeIdent, ok := param.Type.(*ast.Identifier)
	if !ok {
		t.Fatalf("type is not *ast.Identifier. got=%T", param.Type)
	}

	// The full namespace path is in the identifier value
	if typeIdent.Value != `App\Models\User` {
		t.Errorf("type name not 'App\\Models\\User'. got=%s", typeIdent.Value)
	}
}

// Test complex type combinations

func TestComplexPropertyTypes(t *testing.T) {
	input := `<?php
class Example {
	private int $count;
	protected ?string $name;
	public array|object $data;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	// Check first property (int)
	prop1 := classDecl.Body[0].(*ast.PropertyDeclaration)
	if prop1.Type == nil {
		t.Error("first property type is nil")
	}

	// Check second property (?string)
	prop2 := classDecl.Body[1].(*ast.PropertyDeclaration)
	if _, ok := prop2.Type.(*ast.NullableType); !ok {
		t.Errorf("second property type is not nullable. got=%T", prop2.Type)
	}

	// Check third property (array|object)
	prop3 := classDecl.Body[2].(*ast.PropertyDeclaration)
	if _, ok := prop3.Type.(*ast.UnionType); !ok {
		t.Errorf("third property type is not union. got=%T", prop3.Type)
	}
}

func TestComplexMethodSignatures(t *testing.T) {
	input := `<?php
interface Repository {
	public function find(int|string $id): ?object;
	public function findAll(int $limit): array|Collection;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	interfaceDecl := program.Statements[0].(*ast.InterfaceDeclaration)

	// Check first method parameter (int|string)
	method1 := interfaceDecl.Body[0]
	if _, ok := method1.Parameters[0].Type.(*ast.UnionType); !ok {
		t.Errorf("first method parameter is not union type. got=%T", method1.Parameters[0].Type)
	}

	// Check first method return type (?object)
	if _, ok := method1.ReturnType.(*ast.NullableType); !ok {
		t.Errorf("first method return type is not nullable. got=%T", method1.ReturnType)
	}

	// Check second method return type (array|Collection)
	method2 := interfaceDecl.Body[1]
	if _, ok := method2.ReturnType.(*ast.UnionType); !ok {
		t.Errorf("second method return type is not union. got=%T", method2.ReturnType)
	}
}

// Test edge cases

func TestMixedType(t *testing.T) {
	input := `<?php function test(mixed $x): mixed {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	// Check parameter type
	paramType, ok := funcDecl.Parameters[0].Type.(*ast.Identifier)
	if !ok || paramType.Value != "mixed" {
		t.Errorf("parameter type not 'mixed'. got=%v", funcDecl.Parameters[0].Type)
	}

	// Check return type
	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "mixed" {
		t.Errorf("return type not 'mixed'. got=%v", funcDecl.ReturnType)
	}
}

func TestVoidReturnType(t *testing.T) {
	input := `<?php function test(): void {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "void" {
		t.Errorf("return type not 'void'. got=%v", funcDecl.ReturnType)
	}
}

func TestNeverReturnType(t *testing.T) {
	input := `<?php function terminate(): never {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "never" {
		t.Errorf("return type not 'never'. got=%v", funcDecl.ReturnType)
	}
}

// Additional edge case tests for better coverage

func TestTypeHelperFunctions(t *testing.T) {
	// Test isScalarType
	scalarTests := []struct {
		name     string
		expected bool
	}{
		{"int", true},
		{"integer", true},
		{"float", true},
		{"double", true},
		{"string", true},
		{"bool", true},
		{"boolean", true},
		{"User", false},
		{"mixed", false},
		{"array", false},
	}

	for _, tt := range scalarTests {
		result := isScalarType(tt.name)
		if result != tt.expected {
			t.Errorf("isScalarType(%q) = %v, expected %v", tt.name, result, tt.expected)
		}
	}

	// Test isSpecialType
	specialTests := []struct {
		name     string
		expected bool
	}{
		{"mixed", true},
		{"void", true},
		{"never", true},
		{"null", true},
		{"false", true},
		{"true", true},
		{"static", true},
		{"self", true},
		{"parent", true},
		{"int", false},
		{"User", false},
	}

	for _, tt := range specialTests {
		result := isSpecialType(tt.name)
		if result != tt.expected {
			t.Errorf("isSpecialType(%q) = %v, expected %v", tt.name, result, tt.expected)
		}
	}

	// Test isCompoundType
	compoundTests := []struct {
		name     string
		expected bool
	}{
		{"array", true},
		{"object", true},
		{"callable", true},
		{"iterable", true},
		{"resource", true},
		{"int", false},
		{"mixed", false},
		{"User", false},
	}

	for _, tt := range compoundTests {
		result := isCompoundType(tt.name)
		if result != tt.expected {
			t.Errorf("isCompoundType(%q) = %v, expected %v", tt.name, result, tt.expected)
		}
	}
}

func TestStaticReturnType(t *testing.T) {
	input := `<?php
	class Foo {
		public function getInstance(): static {}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	methodDecl := classDecl.Body[0].(*ast.MethodDeclaration)

	returnType, ok := methodDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "static" {
		t.Errorf("return type not 'static'. got=%v", methodDecl.ReturnType)
	}
}

func TestSelfReturnType(t *testing.T) {
	input := `<?php
	class Foo {
		public function getSelf(): self {}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	methodDecl := classDecl.Body[0].(*ast.MethodDeclaration)

	returnType, ok := methodDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "self" {
		t.Errorf("return type not 'self'. got=%v", methodDecl.ReturnType)
	}
}

func TestParentReturnType(t *testing.T) {
	input := `<?php
	class Bar extends Foo {
		public function getParent(): parent {}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	methodDecl := classDecl.Body[0].(*ast.MethodDeclaration)

	returnType, ok := methodDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "parent" {
		t.Errorf("return type not 'parent'. got=%v", methodDecl.ReturnType)
	}
}

func TestFalseReturnType(t *testing.T) {
	input := `<?php function alwaysFalse(): false {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "false" {
		t.Errorf("return type not 'false'. got=%v", funcDecl.ReturnType)
	}
}

func TestTrueReturnType(t *testing.T) {
	input := `<?php function alwaysTrue(): true {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "true" {
		t.Errorf("return type not 'true'. got=%v", funcDecl.ReturnType)
	}
}

func TestResourceType(t *testing.T) {
	input := `<?php function getResource(): resource {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "resource" {
		t.Errorf("return type not 'resource'. got=%v", funcDecl.ReturnType)
	}
}

func TestNamespacedTypeHint(t *testing.T) {
	input := `<?php function getUser(): \App\Models\User {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok {
		t.Fatalf("return type not Identifier. got=%T", funcDecl.ReturnType)
	}

	if returnType.Value != "\\App\\Models\\User" {
		t.Errorf("return type not '\\App\\Models\\User'. got=%q", returnType.Value)
	}
}

func TestComplexUnionWithNullable(t *testing.T) {
	input := `<?php function getValue(): int|string|null {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	unionType, ok := funcDecl.ReturnType.(*ast.UnionType)
	if !ok {
		t.Fatalf("return type not UnionType. got=%T", funcDecl.ReturnType)
	}

	if len(unionType.Types) != 3 {
		t.Errorf("union should have 3 types. got=%d", len(unionType.Types))
	}
}

func TestTypeAliases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "integer alias",
			input:    `<?php function test(): integer {}`,
			expected: "integer",
		},
		{
			name:     "double alias",
			input:    `<?php function test(): double {}`,
			expected: "double",
		},
		{
			name:     "boolean alias",
			input:    `<?php function test(): boolean {}`,
			expected: "boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
			returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
			if !ok || returnType.Value != tt.expected {
				t.Errorf("return type not %q. got=%v", tt.expected, funcDecl.ReturnType)
			}
		})
	}
}
