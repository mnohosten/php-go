package lexer

import (
	"testing"
)

func TestLexerBasicTokens(t *testing.T) {
	input := `<?php
$x = 5;
$y = 10;
$sum = $x + $y;
echo $sum;
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{OPEN_TAG, "<?php"},
		{VARIABLE, "$x"},
		{ASSIGN, "="},
		{INTEGER, "5"},
		{SEMICOLON, ";"},
		{VARIABLE, "$y"},
		{ASSIGN, "="},
		{INTEGER, "10"},
		{SEMICOLON, ";"},
		{VARIABLE, "$sum"},
		{ASSIGN, "="},
		{VARIABLE, "$x"},
		{PLUS, "+"},
		{VARIABLE, "$y"},
		{SEMICOLON, ";"},
		{ECHO, "echo"},
		{VARIABLE, "$sum"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input, "test.php")

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerKeywords(t *testing.T) {
	input := `if else while for foreach function class interface return`

	tests := []TokenType{
		IF, ELSE, WHILE, FOR, FOREACH, FUNCTION, CLASS, INTERFACE, RETURN, EOF,
	}

	l := New(input, "test.php")

	for i, expectedType := range tests {
		tok := l.NextToken()

		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"123", INTEGER, "123"},
		{"0x1A", INTEGER, "0x1A"},
		{"0b1010", INTEGER, "0b1010"},
		{"0o777", INTEGER, "0o777"},
		{"123.45", FLOAT, "123.45"},
		{"1.23e4", FLOAT, "1.23e4"},
		{"1.23E-4", FLOAT, "1.23E-4"},
		{".5", FLOAT, ".5"},
		{"123_456", INTEGER, "123_456"},
		{"1_234.567_89", FLOAT, "1_234.567_89"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("tokentype wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{"single quoted", "'hello world'", "hello world"},
		{"double quoted", `"hello world"`, "hello world"},
		{"escaped single quote", `'it\'s'`, "it's"},
		{"escaped backslash", `'back\\slash'`, "back\\slash"},
		{"escaped newline", `"hello\nworld"`, "hello\nworld"},
		{"escaped tab", `"hello\tworld"`, "hello\tworld"},
		{"escaped quote", `"say \"hi\""`, "say \"hi\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != STRING {
				t.Errorf("tokentype wrong. expected=STRING, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestLexerOperators(t *testing.T) {
	input := `+ - * / % **
== === != !== < <= > >= <=>
++ -- += -= *= /= %=
&& || !
& | ^ ~ << >>
-> :: ?-> => ...
?? ??=
`

	tests := []TokenType{
		PLUS, MINUS, ASTERISK, SLASH, PERCENT, POWER,
		EQ, IDENTICAL, NE, NOT_IDENTICAL, LT, LE, GT, GE, SPACESHIP,
		INC, DEC, PLUS_ASSIGN, MINUS_ASSIGN, MUL_ASSIGN, DIV_ASSIGN, MOD_ASSIGN,
		LOGICAL_AND, LOGICAL_OR, LOGICAL_NOT,
		BITWISE_AND, BITWISE_OR, BITWISE_XOR, BITWISE_NOT, SL, SR,
		OBJECT_OPERATOR, PAAMAYIM_NEKUDOTAYIM, NULLSAFE_OPERATOR, DOUBLE_ARROW, ELLIPSIS,
		COALESCE, COALESCE_ASSIGN,
		EOF,
	}

	l := New(input, "test.php")

	for i, expectedType := range tests {
		tok := l.NextToken()

		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestLexerComments(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{"single line //", "// this is a comment", "// this is a comment"},
		{"single line #", "# this is a comment", "# this is a comment"},
		{"multi line", "/* multi\nline\ncomment */", "/* multi\nline\ncomment */"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != COMMENT {
				t.Errorf("tokentype wrong. expected=COMMENT, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestLexerVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"$x", "$x"},
		{"$name", "$name"},
		{"$_private", "$_private"},
		{"$camelCase", "$camelCase"},
		{"$under_score", "$under_score"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != VARIABLE {
				t.Errorf("tokentype wrong. expected=VARIABLE, got=%q", tok.Type)
			}

			if tok.Literal != tt.expected {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expected, tok.Literal)
			}
		})
	}
}

func TestLexerPHPTags(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType TokenType
	}{
		{"full tag", "<?php", OPEN_TAG},
		{"echo tag", "<?=", OPEN_TAG_ECHO},
		{"short tag", "<?", OPEN_TAG},
		{"close tag", "?>", CLOSE_TAG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("tokentype wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}
		})
	}
}

func TestLexerAttributes(t *testing.T) {
	input := "#[Route('/path')]"

	l := New(input, "test.php")
	tok := l.NextToken()

	if tok.Type != ATTRIBUTE_START {
		t.Errorf("tokentype wrong. expected=ATTRIBUTE_START, got=%q", tok.Type)
	}

	if tok.Literal != "#[" {
		t.Errorf("literal wrong. expected=%q, got=%q", "#[", tok.Literal)
	}
}

func TestLexerDelimiters(t *testing.T) {
	input := "( ) { } [ ] ; , : ::"

	tests := []TokenType{
		LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET,
		SEMICOLON, COMMA, COLON, PAAMAYIM_NEKUDOTAYIM, EOF,
	}

	l := New(input, "test.php")

	for i, expectedType := range tests {
		tok := l.NextToken()

		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expectedType, tok.Type)
		}
	}
}

func TestLexerPositionTracking(t *testing.T) {
	input := `<?php
$x = 5;`

	l := New(input, "test.php")

	tok := l.NextToken() // <?php
	if tok.Pos.Line != 1 || tok.Pos.Column != 1 {
		t.Errorf("wrong position for <?php. expected line=1, col=1, got line=%d, col=%d",
			tok.Pos.Line, tok.Pos.Column)
	}

	tok = l.NextToken() // $x
	if tok.Pos.Line != 2 {
		t.Errorf("wrong line for $x. expected=2, got=%d", tok.Pos.Line)
	}
}

func TestLexerComplexExpression(t *testing.T) {
	input := `<?php
if ($x >= 10 && $y !== null) {
    echo "Valid";
}`

	expectedTokens := []TokenType{
		OPEN_TAG,
		IF, LPAREN, VARIABLE, GE, INTEGER, LOGICAL_AND, VARIABLE,
		NOT_IDENTICAL, NULL, RPAREN, LBRACE,
		ECHO, STRING, SEMICOLON,
		RBRACE, EOF,
	}

	l := New(input, "test.php")

	for i, expectedType := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestLexerErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unterminated string double", `"hello`},
		{"unterminated string single", `'hello`},
		{"unterminated comment", `/* comment`},
		{"invalid variable", `$`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != ILLEGAL {
				t.Errorf("expected ILLEGAL token for %q, got %q", tt.input, tok.Type)
			}
		})
	}
}

// Benchmark tests
func BenchmarkLexer(b *testing.B) {
	input := `<?php
$x = 5;
$y = 10;
$sum = $x + $y;
echo $sum;
`

	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

func BenchmarkLexerLargeFile(b *testing.B) {
	// Simulate a larger PHP file
	input := `<?php
namespace App\\Controllers;

class UserController {
    private $db;

    public function __construct($db) {
        $this->db = $db;
    }

    public function getUser($id) {
        if ($id <= 0) {
            return null;
        }

        $sql = "SELECT * FROM users WHERE id = ?";
        $result = $this->db->query($sql, [$id]);

        return $result->fetch();
    }
}
`

	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}
