package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// TestSimpleClassWithMethods tests parsing a realistic class with methods
func TestSimpleClassWithMethods(t *testing.T) {
	input := `<?php

class Calculator {
	private $result;

	public function add($a, $b) {
		$this->result = $a + $b;
		return $this->result;
	}

	public function subtract($a, $b) {
		$this->result = $a - $b;
		return $this->result;
	}

	public function multiply($a, $b) {
		$this->result = $a * $b;
		return $this->result;
	}

	public function divide($a, $b) {
		if ($b == 0) {
			throw new Exception("Division by zero");
		}
		$this->result = $a / $b;
		return $this->result;
	}

	public function getResult() {
		return $this->result;
	}
}`

	l := lexer.New(input, "Calculator.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) == 0 {
		t.Fatal("parser returned empty program")
	}

	classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("statement not ClassDeclaration. got=%T", program.Statements[0])
	}

	if classDecl.Name.Value != "Calculator" {
		t.Errorf("class name wrong. expected=Calculator, got=%s", classDecl.Name.Value)
	}

	methodCount := 0
	for _, stmt := range classDecl.Body {
		if _, ok := stmt.(*ast.MethodDeclaration); ok {
			methodCount++
		}
	}

	if methodCount != 5 {
		t.Errorf("expected 5 methods, got %d", methodCount)
	}
}

func TestComplexControlFlow(t *testing.T) {
	input := `<?php

function validatePassword($password) {
	if ($password == null) {
		return false;
	}

	$length = strlen($password);
	if ($length < 8) {
		return false;
	}

	if ($length > 100) {
		return false;
	}

	return true;
}

function processItems($items) {
	$results = [];

	for ($i = 0; $i < count($items); $i++) {
		$item = $items[$i];

		if ($item == null) {
			continue;
		}

		if ($item["skip"] == true) {
			continue;
		}

		if ($item["stop"] == true) {
			break;
		}

		$results[$i] = $item;
	}

	return $results;
}

function getGrade($score) {
	if ($score >= 90) {
		return "A";
	} elseif ($score >= 80) {
		return "B";
	} elseif ($score >= 70) {
		return "C";
	} elseif ($score >= 60) {
		return "D";
	} else {
		return "F";
	}
}`

	l := lexer.New(input, "control.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Errorf("expected 3 functions, got %d", len(program.Statements))
	}
}

func TestNestedStructures(t *testing.T) {
	input := `<?php

function processMatrix($matrix) {
	$result = [];

	for ($i = 0; $i < count($matrix); $i++) {
		$row = [];

		for ($j = 0; $j < count($matrix[$i]); $j++) {
			$value = $matrix[$i][$j];

			if ($value > 0) {
				$row[$j] = $value * 2;
			} else {
				$row[$j] = 0;
			}
		}

		$result[$i] = $row;
	}

	return $result;
}

class TreeNode {
	public $value;
	public $left;
	public $right;

	public function __construct($value) {
		$this->value = $value;
		$this->left = null;
		$this->right = null;
	}

	public function insert($value) {
		if ($value < $this->value) {
			if ($this->left == null) {
				$this->left = new TreeNode($value);
			} else {
				$this->left->insert($value);
			}
		} else {
			if ($this->right == null) {
				$this->right = new TreeNode($value);
			} else {
				$this->right->insert($value);
			}
		}
	}
}`

	l := lexer.New(input, "nested.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 2 {
		t.Errorf("expected at least 2 top-level statements, got %d", len(program.Statements))
	}
}

func TestTryCatchFinally(t *testing.T) {
	input := `<?php

class ValidationError extends Exception {}

function processUser($data) {
	try {
		if ($data["email"] == "") {
			throw new ValidationError("Email is required");
		}

		if ($data["age"] < 18) {
			throw new ValidationError("Must be 18 or older");
		}

		return true;
	} catch (ValidationError $e) {
		echo "Validation failed: " . $e->getMessage();
		return false;
	} finally {
		echo "Validation complete";
	}
}`

	l := lexer.New(input, "error.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 2 {
		t.Errorf("expected at least 2 statements, got %d", len(program.Statements))
	}
}

func TestComplexExpressions(t *testing.T) {
	input := `<?php

function calculate($a, $b, $c) {
	$result1 = ($a + $b) * $c;
	$result2 = $a / ($b - $c);
	$result3 = $a % $b + $c ** 2;
	$result4 = ($a << 2) | ($b >> 1);
	$result5 = $a & $b | $c ^ 0xFF;

	return [$result1, $result2, $result3, $result4, $result5];
}

function compareValues($x, $y) {
	$eq = ($x == $y);
	$identical = ($x === $y);
	$notEq = ($x != $y);
	$notIdentical = ($x !== $y);
	$less = ($x < $y);
	$lessEq = ($x <= $y);
	$greater = ($x > $y);
	$greaterEq = ($x >= $y);
	$spaceship = ($x <=> $y);

	return [
		$eq,
		$identical,
		$notEq,
		$notIdentical,
		$less,
		$lessEq,
		$greater,
		$greaterEq,
		$spaceship
	];
}

function stringOperations($str1, $str2) {
	$concat = $str1 . $str2;
	$repeat = $str1 . $str1 . $str1;
	$interpolation = "Hello" . $str1 . "World";

	return $concat . $repeat . $interpolation;
}`

	l := lexer.New(input, "expressions.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Errorf("expected 3 functions, got %d", len(program.Statements))
	}
}

func TestArrayOperations(t *testing.T) {
	input := `<?php

function arrayDemo() {
	$simple = [1, 2, 3, 4, 5];
	$assoc = ["name" => "John", "age" => 30, "city" => "New York"];
	$mixed = [1, "two", 3.0, true, null];

	$nested = [
		[1, 2, 3],
		[4, 5, 6],
		[7, 8, 9]
	];

	$value1 = $simple[0];
	$value2 = $assoc["name"];
	$value3 = $nested[1][2];

	$simple[0] = 10;
	$assoc["country"] = "USA";

	return [$simple, $assoc, $mixed, $nested];
}`

	l := lexer.New(input, "arrays.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) == 0 {
		t.Fatal("parser returned empty program")
	}
}

func TestWhileLoops(t *testing.T) {
	input := `<?php

function countdown($n) {
	while ($n > 0) {
		echo $n;
		$n--;
	}
}

function factorial($n) {
	$result = 1;
	$i = 1;

	do {
		$result = $result * $i;
		$i++;
	} while ($i <= $n);

	return $result;
}`

	l := lexer.New(input, "loops.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Errorf("expected 2 functions, got %d", len(program.Statements))
	}
}

func TestSwitchStatements(t *testing.T) {
	input := `<?php

function getDayName($day) {
	switch ($day) {
		case 1:
			return "Monday";
		case 2:
			return "Tuesday";
		case 3:
			return "Wednesday";
		case 4:
			return "Thursday";
		case 5:
			return "Friday";
		case 6:
			return "Saturday";
		case 7:
			return "Sunday";
		default:
			return "Invalid day";
	}
}

function getHttpStatus($code) {
	switch ($code) {
		case 200:
		case 201:
		case 204:
			return "Success";
		case 400:
		case 404:
			return "Client Error";
		case 500:
		case 503:
			return "Server Error";
		default:
			return "Unknown";
	}
}`

	l := lexer.New(input, "switch.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Errorf("expected 2 functions, got %d", len(program.Statements))
	}
}

func TestInterfacesAndTraits(t *testing.T) {
	input := `<?php

interface Logger {
	public function log($message);
	public function error($message);
}

trait Timestamped {
	public function getTimestamp() {
		return time();
	}
}

class FileLogger implements Logger {
	use Timestamped;

	private $filename;

	public function __construct($filename) {
		$this->filename = $filename;
	}

	public function log($message) {
		$timestamp = $this->getTimestamp();
		$line = "[" . $timestamp . "] " . $message;
		file_put_contents($this->filename, $line, FILE_APPEND);
	}

	public function error($message) {
		$this->log("ERROR: " . $message);
	}
}`

	l := lexer.New(input, "oop.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Count interfaces, traits, and classes
	interfaceCount := 0
	traitCount := 0
	classCount := 0

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.InterfaceDeclaration:
			interfaceCount++
		case *ast.TraitDeclaration:
			traitCount++
		case *ast.ClassDeclaration:
			classCount++
		}
	}

	if interfaceCount != 1 {
		t.Errorf("expected 1 interface, got %d", interfaceCount)
	}

	if traitCount != 1 {
		t.Errorf("expected 1 trait, got %d", traitCount)
	}

	if classCount != 1 {
		t.Errorf("expected 1 class, got %d", classCount)
	}
}
