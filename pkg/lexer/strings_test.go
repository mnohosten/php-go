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
