package lexer

import (
	"testing"
)

func TestHeredoc(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{
			name: "simple heredoc",
			input: `<<<EOT
Hello World
EOT;`,
			expectedType:    HEREDOC,
			expectedLiteral: "Hello World\n",
		},
		{
			name: "heredoc with multiple lines",
			input: `<<<EOT
Line 1
Line 2
Line 3
EOT;`,
			expectedType:    HEREDOC,
			expectedLiteral: "Line 1\nLine 2\nLine 3\n",
		},
		{
			name: "heredoc with variables",
			input: `<<<EOT
Hello $name
EOT;`,
			expectedType:    HEREDOC,
			expectedLiteral: "Hello $name\n",
		},
		{
			name: "heredoc without semicolon",
			input: `<<<EOT
Hello
EOT
`,
			expectedType:    HEREDOC,
			expectedLiteral: "Hello\n",
		},
		{
			name: "heredoc with different label",
			input: `<<<HTML
<div>Content</div>
HTML;`,
			expectedType:    HEREDOC,
			expectedLiteral: "<div>Content</div>\n",
		},
		{
			name: "heredoc with indented content",
			input: `<<<EOT
    Indented line
        More indented
EOT;`,
			expectedType:    HEREDOC,
			expectedLiteral: "    Indented line\n        More indented\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("tokentype wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong.\nexpected=%q\ngot=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNowdoc(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{
			name: "simple nowdoc single quotes",
			input: `<<<'EOT'
Hello World
EOT;`,
			expectedType:    NOWDOC,
			expectedLiteral: "Hello World\n",
		},
		{
			name: "nowdoc double quotes",
			input: `<<<"EOT"
Hello World
EOT;`,
			expectedType:    NOWDOC,
			expectedLiteral: "Hello World\n",
		},
		{
			name: "nowdoc with variables (not interpolated)",
			input: `<<<'EOT'
Hello $name
Price: $100
EOT;`,
			expectedType:    NOWDOC,
			expectedLiteral: "Hello $name\nPrice: $100\n",
		},
		{
			name: "nowdoc with backslashes",
			input: `<<<'EOT'
Path: C:\Users\Name
EOT;`,
			expectedType:    NOWDOC,
			expectedLiteral: "Path: C:\\Users\\Name\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("tokentype wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong.\nexpected=%q\ngot=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestHeredocEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType TokenType
		expectError  bool
	}{
		{
			name: "heredoc empty",
			input: `<<<EOT
EOT;`,
			expectedType: HEREDOC,
			expectError:  false,
		},
		{
			name: "heredoc with blank lines",
			input: `<<<EOT

Line with content

EOT;`,
			expectedType: HEREDOC,
			expectError:  false,
		},
		{
			name: "unterminated heredoc",
			input: `<<<EOT
Hello World
`,
			expectedType: ILLEGAL,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tt.expectError {
				if tok.Type != ILLEGAL {
					t.Errorf("expected ILLEGAL token, got %q", tok.Type)
				}
			} else {
				if tok.Type != tt.expectedType {
					t.Errorf("tokentype wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
				}
			}
		})
	}
}

func TestHeredocIndented(t *testing.T) {
	// PHP 7.3+ allows indented closing tags
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name: "indented closing tag with spaces",
			input: `<<<EOT
Hello
    EOT;`,
			expectedLiteral: "Hello\n",
		},
		{
			name: "indented closing tag with tabs",
			input: "<<<EOT\nHello\n\tEOT;",
			expectedLiteral: "Hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != HEREDOC {
				t.Errorf("tokentype wrong. expected=HEREDOC, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong.\nexpected=%q\ngot=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []StringInterpolation
	}{
		{
			name:  "simple variable",
			input: "Hello $name",
			expected: []StringInterpolation{
				{IsVariable: false, Value: "Hello "},
				{IsVariable: true, Value: "$name"},
			},
		},
		{
			name:  "multiple variables",
			input: "$greeting $name",
			expected: []StringInterpolation{
				{IsVariable: true, Value: "$greeting"},
				{IsVariable: false, Value: " "},
				{IsVariable: true, Value: "$name"},
			},
		},
		{
			name:  "variable at end",
			input: "Hello $name",
			expected: []StringInterpolation{
				{IsVariable: false, Value: "Hello "},
				{IsVariable: true, Value: "$name"},
			},
		},
		{
			name:  "no interpolation",
			input: "Hello World",
			expected: []StringInterpolation{
				{IsVariable: false, Value: "Hello World"},
			},
		},
		{
			name:  "dollar sign not followed by variable",
			input: "Price: $100",
			expected: []StringInterpolation{
				{IsVariable: false, Value: "Price: $100"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStringInterpolation(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("wrong number of parts. expected=%d, got=%d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i].IsVariable != tt.expected[i].IsVariable {
					t.Errorf("part[%d] IsVariable wrong. expected=%v, got=%v",
						i, tt.expected[i].IsVariable, result[i].IsVariable)
				}
				if result[i].Value != tt.expected[i].Value {
					t.Errorf("part[%d] Value wrong. expected=%q, got=%q",
						i, tt.expected[i].Value, result[i].Value)
				}
			}
		})
	}
}

func TestHasInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello $name", true},
		{"Hello World", false},
		{"$var", true},
		{"Price: $100", false}, // $ not followed by valid identifier
		{"Test ${expr}", true},
		{"Test {$expr}", true},
		{"$", false},
		{"$$var", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasInterpolation(tt.input)
			if result != tt.expected {
				t.Errorf("hasInterpolation(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidHeredocLabel(t *testing.T) {
	tests := []struct {
		label    string
		expected bool
	}{
		{"EOT", true},
		{"HTML", true},
		{"_label", true},
		{"label123", true},
		{"MY_LABEL", true},
		{"123invalid", false},
		{"", false},
		{"label-dash", false},
		{"label.dot", false},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			result := isValidHeredocLabel(tt.label)
			if result != tt.expected {
				t.Errorf("isValidHeredocLabel(%q) = %v, expected %v", tt.label, result, tt.expected)
			}
		})
	}
}

func TestDoubleQuotedStringEscapes(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name:            "hex escape",
			input:           `"Hello\x41"`,
			expectedLiteral: "HelloA",
		},
		{
			name:            "multiple hex escapes",
			input:           `"\x48\x65\x6C\x6C\x6F"`,
			expectedLiteral: "Hello",
		},
		{
			name:            "null byte",
			input:           `"test\0"`,
			expectedLiteral: "test\x00",
		},
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

func TestShiftOperatorVsHeredoc(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType TokenType
	}{
		{
			name:         "left shift operator",
			input:        "$x << 2",
			expectedType: VARIABLE, // First token
		},
		{
			name:         "valid heredoc",
			input:        "<<<EOT\nHello\nEOT;",
			expectedType: HEREDOC,
		},
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

func TestComplexHeredocScenarios(t *testing.T) {
	input := `<?php
$text = <<<EOT
This is a heredoc string
with multiple lines
and some $variables
EOT;
echo $text;
`

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
	}{
		{OPEN_TAG, "<?php"},
		{VARIABLE, "$text"},
		{ASSIGN, "="},
		{HEREDOC, "This is a heredoc string\nwith multiple lines\nand some $variables\n"},
		{SEMICOLON, ";"},
		{ECHO, "echo"},
		{VARIABLE, "$text"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input, "test.php")

	for i, expected := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] type wrong. expected=%q, got=%q",
				i, expected.tokenType, tok.Type)
		}

		if expected.literal != "" && tok.Literal != expected.literal {
			t.Errorf("token[%d] literal wrong.\nexpected=%q\ngot=%q",
				i, expected.literal, tok.Literal)
		}
	}
}

// Benchmark tests
func BenchmarkHeredoc(b *testing.B) {
	input := `<<<EOT
Line 1
Line 2
Line 3
Line 4
Line 5
EOT;`

	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		l.NextToken()
	}
}

func BenchmarkStringInterpolation(b *testing.B) {
	input := "Hello $name, you have $count messages"

	for i := 0; i < b.N; i++ {
		parseStringInterpolation(input)
	}
}

// Additional edge case tests for better coverage

func TestStringEscapeSequences(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name:            "newline escape",
			input:           `"line1\nline2"`,
			expectedLiteral: "line1\nline2",
		},
		{
			name:            "tab escape",
			input:           `"col1\tcol2"`,
			expectedLiteral: "col1\tcol2",
		},
		{
			name:            "carriage return",
			input:           `"before\rafter"`,
			expectedLiteral: "before\rafter",
		},
		{
			name:            "backslash escape",
			input:           `"path\\to\\file"`,
			expectedLiteral: "path\\to\\file",
		},
		{
			name:            "dollar sign escape",
			input:           `"price is \$10"`,
			expectedLiteral: "price is $10",
		},
		{
			name:            "quote escape",
			input:           `"He said \"Hello\""`,
			expectedLiteral: "He said \"Hello\"",
		},
		{
			name:            "null byte escape",
			input:           `"null\0byte"`,
			expectedLiteral: "null\x00byte",
		},
		{
			name:            "unknown escape keeps backslash",
			input:           `"test\q"`,
			expectedLiteral: "test\\q",
		},
		{
			name:            "hex uppercase letters",
			input:           `"\x41\x42\x43"`,
			expectedLiteral: "ABC",
		},
		{
			name:            "hex lowercase letters",
			input:           `"\x61\x62\x63"`,
			expectedLiteral: "abc",
		},
		{
			name:            "hex mixed case",
			input:           `"\x4A\x6b"`,
			expectedLiteral: "Jk",
		},
		{
			name:            "invalid hex escape - non-hex chars",
			input:           `"\xGH"`,
			expectedLiteral: "\\xGH",
		},
		{
			name:            "invalid hex escape - one digit",
			input:           `"\x4Z"`,
			expectedLiteral: "\\x4Z",
		},
		{
			name:            "mixed escapes",
			input:           `"line1\nline2\ttab\x41"`,
			expectedLiteral: "line1\nline2\ttabA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != STRING {
				t.Errorf("token type wrong. expected=STRING, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestUnterminatedString(t *testing.T) {
	input := `"unterminated string`
	l := New(input, "test.php")
	tok := l.NextToken()

	if tok.Type != ILLEGAL {
		t.Errorf("expected ILLEGAL token for unterminated string, got %q", tok.Type)
	}

	if tok.Literal != "unterminated string" {
		t.Errorf("expected error message 'unterminated string', got %q", tok.Literal)
	}
}

func TestStringWithNewlines(t *testing.T) {
	input := "\"line1\nline2\nline3\""
	l := New(input, "test.php")
	tok := l.NextToken()

	if tok.Type != STRING {
		t.Errorf("token type wrong. expected=STRING, got=%q", tok.Type)
	}

	expected := "line1\nline2\nline3"
	if tok.Literal != expected {
		t.Errorf("literal wrong. expected=%q, got=%q", expected, tok.Literal)
	}

	// Check that line number was incremented
	if tok.Pos.Line != 1 {
		t.Errorf("expected string to start at line 1, got line %d", tok.Pos.Line)
	}
}

func TestStringWithDollarSign(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name:            "dollar sign alone",
			input:           `"price $"`,
			expectedLiteral: "price $",
		},
		{
			name:            "dollar sign with text",
			input:           `"total $amount"`,
			expectedLiteral: "total $amount",
		},
		{
			name:            "multiple dollar signs",
			input:           `"$a $b $c"`,
			expectedLiteral: "$a $b $c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != STRING {
				t.Errorf("token type wrong. expected=STRING, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestHeredocVariations(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name: "heredoc with blank lines",
			input: `<<<EOT
line1

line3
EOT;`,
			expectedLiteral: "line1\n\nline3\n",
		},
		{
			name: "heredoc with only newlines",
			input: `<<<EOT


EOT;`,
			expectedLiteral: "\n\n",
		},
		{
			name: "heredoc with trailing spaces",
			input: `<<<EOT
text
EOT;`,
			expectedLiteral: "text\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != HEREDOC {
				t.Errorf("token type wrong. expected=HEREDOC, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNowdocVariations(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedLiteral string
	}{
		{
			name: "nowdoc with variables not interpolated",
			input: `<<<'EOT'
Hello $name
Total $count
EOT;`,
			expectedLiteral: "Hello $name\nTotal $count\n",
		},
		{
			name: "nowdoc with backslashes",
			input: `<<<'EOT'
C:\path\to\file
EOT;`,
			expectedLiteral: "C:\\path\\to\\file\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input, "test.php")
			tok := l.NextToken()

			if tok.Type != NOWDOC {
				t.Errorf("token type wrong. expected=NOWDOC, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}
