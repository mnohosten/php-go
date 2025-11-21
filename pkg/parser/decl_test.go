package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// Test function declarations

func TestFunctionDeclaration(t *testing.T) {
	input := `<?php
function add($a, $b) {
	return $a + $b;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	funcDecl, ok := program.Statements[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("statement is not *ast.FunctionDeclaration. got=%T", program.Statements[0])
	}

	if funcDecl.Name.Value != "add" {
		t.Errorf("function name not 'add'. got=%s", funcDecl.Name.Value)
	}

	if len(funcDecl.Parameters) != 2 {
		t.Errorf("expected 2 parameters. got=%d", len(funcDecl.Parameters))
	}
}

func TestFunctionWithReturnType(t *testing.T) {
	input := `<?php
function divide($a, $b): float {
	return $a / $b;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	if funcDecl.ReturnType == nil {
		t.Fatal("return type is nil")
	}

	returnType, ok := funcDecl.ReturnType.(*ast.Identifier)
	if !ok || returnType.Value != "float" {
		t.Errorf("return type not 'float'. got=%v", funcDecl.ReturnType)
	}
}

func TestFunctionWithReferenceReturn(t *testing.T) {
	input := `<?php
function &getValue() {
	return $value;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	if !funcDecl.ByRef {
		t.Error("function should return by reference")
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input       string
		paramCount  int
		checkParam  func(*testing.T, *ast.Parameter)
	}{
		{
			input:      `<?php function test($a) {}`,
			paramCount: 1,
			checkParam: func(t *testing.T, p *ast.Parameter) {
				if p.Name.Name != "a" {
					t.Errorf("param name not 'a'. got=%s", p.Name.Name)
				}
			},
		},
		{
			input:      `<?php function test(int $x) {}`,
			paramCount: 1,
			checkParam: func(t *testing.T, p *ast.Parameter) {
				if p.Type == nil {
					t.Fatal("param type is nil")
				}
				typeIdent, ok := p.Type.(*ast.Identifier)
				if !ok || typeIdent.Value != "int" {
					t.Errorf("param type not 'int'. got=%v", p.Type)
				}
			},
		},
		{
			input:      `<?php function test($x = 10) {}`,
			paramCount: 1,
			checkParam: func(t *testing.T, p *ast.Parameter) {
				if p.DefaultValue == nil {
					t.Fatal("param default value is nil")
				}
			},
		},
		{
			input:      `<?php function test(&$ref) {}`,
			paramCount: 1,
			checkParam: func(t *testing.T, p *ast.Parameter) {
				if !p.ByRef {
					t.Error("param should be by reference")
				}
			},
		},
		{
			input:      `<?php function test(...$rest) {}`,
			paramCount: 1,
			checkParam: func(t *testing.T, p *ast.Parameter) {
				if !p.Variadic {
					t.Error("param should be variadic")
				}
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		funcDecl := program.Statements[0].(*ast.FunctionDeclaration)
		if len(funcDecl.Parameters) != tt.paramCount {
			t.Errorf("expected %d parameters. got=%d", tt.paramCount, len(funcDecl.Parameters))
			continue
		}

		tt.checkParam(t, funcDecl.Parameters[0])
	}
}

// Test class declarations

func TestClassDeclaration(t *testing.T) {
	input := `<?php
class User {
	public $name;

	public function getName() {
		return $this->name;
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got=%d", len(program.Statements))
	}

	classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("statement is not *ast.ClassDeclaration. got=%T", program.Statements[0])
	}

	if classDecl.Name.Value != "User" {
		t.Errorf("class name not 'User'. got=%s", classDecl.Name.Value)
	}

	if len(classDecl.Body) != 2 {
		t.Errorf("expected 2 members. got=%d", len(classDecl.Body))
	}
}

func TestClassWithInheritance(t *testing.T) {
	input := `<?php
class Admin extends User {
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	if classDecl.Extends == nil {
		t.Fatal("extends is nil")
	}

	if classDecl.Extends.Value != "User" {
		t.Errorf("extends not 'User'. got=%s", classDecl.Extends.Value)
	}
}

func TestClassWithInterfaces(t *testing.T) {
	input := `<?php
class MyClass implements Serializable, JsonSerializable {
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	if len(classDecl.Implements) != 2 {
		t.Errorf("expected 2 interfaces. got=%d", len(classDecl.Implements))
	}

	if classDecl.Implements[0].Value != "Serializable" {
		t.Errorf("first interface not 'Serializable'. got=%s", classDecl.Implements[0].Value)
	}
}

func TestAbstractClass(t *testing.T) {
	input := `<?php
abstract class Shape {
	abstract public function area(): float;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	if len(classDecl.Modifiers) != 1 || classDecl.Modifiers[0] != "abstract" {
		t.Errorf("class should have 'abstract' modifier. got=%v", classDecl.Modifiers)
	}

	method := classDecl.Body[0].(*ast.MethodDeclaration)
	if !method.Abstract {
		t.Error("method should be abstract")
	}

	if method.Body != nil {
		t.Error("abstract method should not have body")
	}
}

func TestFinalClass(t *testing.T) {
	input := `<?php
final class Config {
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	if len(classDecl.Modifiers) != 1 || classDecl.Modifiers[0] != "final" {
		t.Errorf("class should have 'final' modifier. got=%v", classDecl.Modifiers)
	}
}

// Test class properties

func TestPropertyDeclaration(t *testing.T) {
	tests := []struct {
		input      string
		visibility string
		propName   string
		hasType    bool
		hasDefault bool
	}{
		{
			input:      `<?php class Test { public $name; }`,
			visibility: "public",
			propName:   "name",
			hasType:    false,
			hasDefault: false,
		},
		{
			input:      `<?php class Test { private string $email; }`,
			visibility: "private",
			propName:   "email",
			hasType:    true,
			hasDefault: false,
		},
		{
			input:      `<?php class Test { protected $count = 0; }`,
			visibility: "protected",
			propName:   "count",
			hasType:    false,
			hasDefault: true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		classDecl := program.Statements[0].(*ast.ClassDeclaration)
		propDecl, ok := classDecl.Body[0].(*ast.PropertyDeclaration)
		if !ok {
			t.Fatalf("member is not *ast.PropertyDeclaration. got=%T", classDecl.Body[0])
		}

		if propDecl.Visibility != tt.visibility {
			t.Errorf("visibility not '%s'. got=%s", tt.visibility, propDecl.Visibility)
		}

		if propDecl.Properties[0].Name.Name != tt.propName {
			t.Errorf("property name not '%s'. got=%s", tt.propName, propDecl.Properties[0].Name.Name)
		}

		if (propDecl.Type != nil) != tt.hasType {
			t.Errorf("hasType mismatch. expected=%v, got=%v", tt.hasType, propDecl.Type != nil)
		}

		if (propDecl.Properties[0].DefaultValue != nil) != tt.hasDefault {
			t.Errorf("hasDefault mismatch. expected=%v, got=%v", tt.hasDefault, propDecl.Properties[0].DefaultValue != nil)
		}
	}
}

func TestStaticProperty(t *testing.T) {
	input := `<?php
class Counter {
	public static $count = 0;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	propDecl := classDecl.Body[0].(*ast.PropertyDeclaration)

	if !propDecl.Static {
		t.Error("property should be static")
	}
}

func TestReadonlyProperty(t *testing.T) {
	input := `<?php
class User {
	public readonly string $id;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	propDecl := classDecl.Body[0].(*ast.PropertyDeclaration)

	if !propDecl.Readonly {
		t.Error("property should be readonly")
	}
}

func TestMultiplePropertiesOneLine(t *testing.T) {
	input := `<?php
class Test {
	public $x, $y, $z;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	propDecl := classDecl.Body[0].(*ast.PropertyDeclaration)

	if len(propDecl.Properties) != 3 {
		t.Errorf("expected 3 properties. got=%d", len(propDecl.Properties))
	}
}

// Test class methods

func TestMethodDeclaration(t *testing.T) {
	input := `<?php
class Calculator {
	public function add($a, $b): int {
		return $a + $b;
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	method, ok := classDecl.Body[0].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("member is not *ast.MethodDeclaration. got=%T", classDecl.Body[0])
	}

	if method.Name.Value != "add" {
		t.Errorf("method name not 'add'. got=%s", method.Name.Value)
	}

	if method.Visibility != "public" {
		t.Errorf("visibility not 'public'. got=%s", method.Visibility)
	}
}

func TestStaticMethod(t *testing.T) {
	input := `<?php
class Helper {
	public static function format($str): string {
		return $str;
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	method := classDecl.Body[0].(*ast.MethodDeclaration)

	if !method.Static {
		t.Error("method should be static")
	}
}

func TestPrivateMethod(t *testing.T) {
	input := `<?php
class Secret {
	private function encrypt(): void {
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	method := classDecl.Body[0].(*ast.MethodDeclaration)

	if method.Visibility != "private" {
		t.Errorf("visibility not 'private'. got=%s", method.Visibility)
	}
}

func TestFinalMethod(t *testing.T) {
	input := `<?php
class Parent {
	final public function cannotOverride() {
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	method := classDecl.Body[0].(*ast.MethodDeclaration)

	if !method.Final {
		t.Error("method should be final")
	}
}

// Test interfaces

func TestInterfaceDeclaration(t *testing.T) {
	input := `<?php
interface Renderable {
	public function render(): string;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	interfaceDecl, ok := program.Statements[0].(*ast.InterfaceDeclaration)
	if !ok {
		t.Fatalf("statement is not *ast.InterfaceDeclaration. got=%T", program.Statements[0])
	}

	if interfaceDecl.Name.Value != "Renderable" {
		t.Errorf("interface name not 'Renderable'. got=%s", interfaceDecl.Name.Value)
	}

	if len(interfaceDecl.Body) != 1 {
		t.Errorf("expected 1 method. got=%d", len(interfaceDecl.Body))
	}
}

func TestInterfaceExtends(t *testing.T) {
	input := `<?php
interface JsonRenderable extends Renderable, Serializable {
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	interfaceDecl := program.Statements[0].(*ast.InterfaceDeclaration)

	if len(interfaceDecl.Extends) != 2 {
		t.Errorf("expected 2 parent interfaces. got=%d", len(interfaceDecl.Extends))
	}
}

// Test traits

func TestTraitDeclaration(t *testing.T) {
	input := `<?php
trait Timestampable {
	public $created_at;

	public function updateTimestamp() {
		$this->created_at = time();
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	traitDecl, ok := program.Statements[0].(*ast.TraitDeclaration)
	if !ok {
		t.Fatalf("statement is not *ast.TraitDeclaration. got=%T", program.Statements[0])
	}

	if traitDecl.Name.Value != "Timestampable" {
		t.Errorf("trait name not 'Timestampable'. got=%s", traitDecl.Name.Value)
	}

	if len(traitDecl.Body) != 2 {
		t.Errorf("expected 2 members. got=%d", len(traitDecl.Body))
	}
}

func TestTraitUse(t *testing.T) {
	input := `<?php
class Post {
	use Timestampable;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	traitUse, ok := classDecl.Body[0].(*ast.TraitUse)
	if !ok {
		t.Fatalf("member is not *ast.TraitUse. got=%T", classDecl.Body[0])
	}

	if len(traitUse.Traits) != 1 {
		t.Errorf("expected 1 trait. got=%d", len(traitUse.Traits))
	}

	if traitUse.Traits[0].Value != "Timestampable" {
		t.Errorf("trait name not 'Timestampable'. got=%s", traitUse.Traits[0].Value)
	}
}

func TestMultipleTraitUse(t *testing.T) {
	input := `<?php
class Post {
	use Timestampable, Loggable;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	traitUse := classDecl.Body[0].(*ast.TraitUse)

	if len(traitUse.Traits) != 2 {
		t.Errorf("expected 2 traits. got=%d", len(traitUse.Traits))
	}
}

// Test class constants

func TestClassConstant(t *testing.T) {
	input := `<?php
class Status {
	const ACTIVE = 1;
	const INACTIVE = 0;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	// First constant declaration
	constDecl, ok := classDecl.Body[0].(*ast.ClassConstantDeclaration)
	if !ok {
		t.Fatalf("member is not *ast.ClassConstantDeclaration. got=%T", classDecl.Body[0])
	}

	if len(constDecl.Constants) != 1 {
		t.Errorf("expected 1 constant in first declaration. got=%d", len(constDecl.Constants))
	}

	if constDecl.Constants[0].Name.Value != "ACTIVE" {
		t.Errorf("constant name not 'ACTIVE'. got=%s", constDecl.Constants[0].Name.Value)
	}
}

func TestClassConstantWithVisibility(t *testing.T) {
	input := `<?php
class Config {
	private const SECRET = "hidden";
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)
	constDecl := classDecl.Body[0].(*ast.ClassConstantDeclaration)

	if constDecl.Visibility != "private" {
		t.Errorf("visibility not 'private'. got=%s", constDecl.Visibility)
	}
}
