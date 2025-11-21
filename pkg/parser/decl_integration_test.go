package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// TestCompleteClassParsing tests parsing a complete PHP class with all features
func TestCompleteClassParsing(t *testing.T) {
	input := `<?php
abstract class Database {
	private string $host;
	protected static int $connections = 0;
	public const MAX_CONNECTIONS = 100;

	public function __construct(string $host, int $port = 3306) {
		$this->host = $host;
		self::$connections++;
	}

	abstract public function connect(): bool;

	final public function getHost(): string {
		return $this->host;
	}

	public static function getConnections(): int {
		return self::$connections;
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

	// Check class modifiers
	if len(classDecl.Modifiers) != 1 || classDecl.Modifiers[0] != "abstract" {
		t.Errorf("class should have 'abstract' modifier")
	}

	// Check class name
	if classDecl.Name.Value != "Database" {
		t.Errorf("class name not 'Database'. got=%s", classDecl.Name.Value)
	}

	// Check members count (2 properties, 1 constant, 4 methods)
	expectedMembers := 7
	if len(classDecl.Body) != expectedMembers {
		t.Errorf("expected %d members. got=%d", expectedMembers, len(classDecl.Body))
	}

	// Check first property (private string $host)
	prop1, ok := classDecl.Body[0].(*ast.PropertyDeclaration)
	if !ok {
		t.Fatalf("first member is not *ast.PropertyDeclaration. got=%T", classDecl.Body[0])
	}
	if prop1.Visibility != "private" {
		t.Errorf("first property visibility not 'private'. got=%s", prop1.Visibility)
	}
	if prop1.Type == nil {
		t.Error("first property should have type hint")
	}

	// Check second property (protected static int $connections = 0)
	prop2, ok := classDecl.Body[1].(*ast.PropertyDeclaration)
	if !ok {
		t.Fatalf("second member is not *ast.PropertyDeclaration. got=%T", classDecl.Body[1])
	}
	if prop2.Visibility != "protected" || !prop2.Static {
		t.Errorf("second property should be 'protected static'")
	}

	// Check constant
	constDecl, ok := classDecl.Body[2].(*ast.ClassConstantDeclaration)
	if !ok {
		t.Fatalf("third member is not *ast.ClassConstantDeclaration. got=%T", classDecl.Body[2])
	}
	if constDecl.Constants[0].Name.Value != "MAX_CONNECTIONS" {
		t.Errorf("constant name not 'MAX_CONNECTIONS'. got=%s", constDecl.Constants[0].Name.Value)
	}

	// Check constructor
	method1, ok := classDecl.Body[3].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("fourth member is not *ast.MethodDeclaration. got=%T", classDecl.Body[3])
	}
	if method1.Name.Value != "__construct" {
		t.Errorf("constructor name not '__construct'. got=%s", method1.Name.Value)
	}
	if len(method1.Parameters) != 2 {
		t.Errorf("constructor should have 2 parameters. got=%d", len(method1.Parameters))
	}

	// Check abstract method
	method2, ok := classDecl.Body[4].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("fifth member is not *ast.MethodDeclaration. got=%T", classDecl.Body[4])
	}
	if !method2.Abstract {
		t.Error("connect method should be abstract")
	}
	if method2.Body != nil {
		t.Error("abstract method should not have body")
	}

	// Check final method
	method3, ok := classDecl.Body[5].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("sixth member is not *ast.MethodDeclaration. got=%T", classDecl.Body[5])
	}
	if !method3.Final {
		t.Error("getHost method should be final")
	}
	if method3.ReturnType == nil {
		t.Error("getHost should have return type")
	}

	// Check static method
	method4, ok := classDecl.Body[6].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("seventh member is not *ast.MethodDeclaration. got=%T", classDecl.Body[6])
	}
	if !method4.Static {
		t.Error("getConnections method should be static")
	}
}

// TestCompleteInterfaceParsing tests parsing a complete interface
func TestCompleteInterfaceParsing(t *testing.T) {
	input := `<?php
interface Cacheable extends Serializable {
	public function cache(): void;
	public function invalidate(): void;
	public function getCacheKey(): string;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	interfaceDecl := program.Statements[0].(*ast.InterfaceDeclaration)

	if interfaceDecl.Name.Value != "Cacheable" {
		t.Errorf("interface name not 'Cacheable'. got=%s", interfaceDecl.Name.Value)
	}

	if len(interfaceDecl.Extends) != 1 {
		t.Errorf("interface should extend 1 interface. got=%d", len(interfaceDecl.Extends))
	}

	if len(interfaceDecl.Body) != 3 {
		t.Errorf("interface should have 3 methods. got=%d", len(interfaceDecl.Body))
	}

	// Check all methods have signatures but no bodies
	for i, method := range interfaceDecl.Body {
		if method.Name == nil {
			t.Errorf("method %d has no name", i)
		}
	}
}

// TestCompleteTraitParsing tests parsing a complete trait
func TestCompleteTraitParsing(t *testing.T) {
	input := `<?php
trait Logger {
	private $logFile;

	public function log(string $message): void {
		file_put_contents($this->logFile, $message);
	}

	protected function getLogFile(): string {
		return $this->logFile;
	}
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	traitDecl := program.Statements[0].(*ast.TraitDeclaration)

	if traitDecl.Name.Value != "Logger" {
		t.Errorf("trait name not 'Logger'. got=%s", traitDecl.Name.Value)
	}

	if len(traitDecl.Body) != 3 {
		t.Errorf("trait should have 3 members. got=%d", len(traitDecl.Body))
	}
}

// TestComplexFunctionSignature tests parsing complex function signatures
func TestComplexFunctionSignature(t *testing.T) {
	input := `<?php
function process(
	string $input,
	int $count = 10,
	bool $verbose = false,
	array &$results = [],
	...$options
): array {
	return [];
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FunctionDeclaration)

	if len(funcDecl.Parameters) != 5 {
		t.Fatalf("expected 5 parameters. got=%d", len(funcDecl.Parameters))
	}

	// Check string $input
	if funcDecl.Parameters[0].Type == nil {
		t.Error("first parameter should have type")
	}
	if funcDecl.Parameters[0].DefaultValue != nil {
		t.Error("first parameter should not have default")
	}

	// Check int $count = 10
	if funcDecl.Parameters[1].DefaultValue == nil {
		t.Error("second parameter should have default")
	}

	// Check array &$results = []
	if !funcDecl.Parameters[3].ByRef {
		t.Error("fourth parameter should be by reference")
	}

	// Check ...$options
	if !funcDecl.Parameters[4].Variadic {
		t.Error("fifth parameter should be variadic")
	}

	// Check return type
	if funcDecl.ReturnType == nil {
		t.Error("function should have return type")
	}
}

// TestClassWithTraitUse tests parsing class with trait usage
func TestClassWithTraitUse(t *testing.T) {
	input := `<?php
class Post {
	use Timestampable, Loggable;

	public $title;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	classDecl := program.Statements[0].(*ast.ClassDeclaration)

	if len(classDecl.Body) != 2 {
		t.Errorf("expected 2 members. got=%d", len(classDecl.Body))
	}

	traitUse, ok := classDecl.Body[0].(*ast.TraitUse)
	if !ok {
		t.Fatalf("first member is not *ast.TraitUse. got=%T", classDecl.Body[0])
	}

	if len(traitUse.Traits) != 2 {
		t.Errorf("expected 2 traits. got=%d", len(traitUse.Traits))
	}
}

// TestMultipleFunctionsAndClasses tests parsing multiple declarations
func TestMultipleFunctionsAndClasses(t *testing.T) {
	input := `<?php
function helper(): void {
}

class User {
	public $name;
}

interface Serializable {
	public function serialize(): string;
}

trait Timestampable {
	public $timestamp;
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements. got=%d", len(program.Statements))
	}

	// Check types
	_, ok1 := program.Statements[0].(*ast.FunctionDeclaration)
	_, ok2 := program.Statements[1].(*ast.ClassDeclaration)
	_, ok3 := program.Statements[2].(*ast.InterfaceDeclaration)
	_, ok4 := program.Statements[3].(*ast.TraitDeclaration)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		t.Error("statements are not of expected types")
	}
}
