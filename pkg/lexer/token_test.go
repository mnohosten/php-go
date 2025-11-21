package lexer

import (
	"testing"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name: "Integer token",
			token: Token{
				Type:    INTEGER,
				Literal: "123",
				Pos:     Position{Filename: "test.php", Line: 1, Column: 1},
			},
			expected: "INTEGER(\"123\") at test.php:1:1",
		},
		{
			name: "Variable token",
			token: Token{
				Type:    VARIABLE,
				Literal: "$name",
				Pos:     Position{Filename: "test.php", Line: 2, Column: 5},
			},
			expected: "VARIABLE(\"$name\") at test.php:2:5",
		},
		{
			name: "Keyword token",
			token: Token{
				Type:    IF,
				Literal: "if",
				Pos:     Position{Line: 3, Column: 1},
			},
			expected: "IF(\"if\") at 3:1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
			if result != tt.expected {
				t.Errorf("Token.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{INTEGER, "INTEGER"},
		{FLOAT, "FLOAT"},
		{STRING, "STRING"},
		{VARIABLE, "VARIABLE"},
		{IF, "IF"},
		{ELSE, "ELSE"},
		{FUNCTION, "FUNCTION"},
		{CLASS, "CLASS"},
		{ECHO, "ECHO"},
		{RETURN, "RETURN"},
		{PLUS, "+"},
		{MINUS, "-"},
		{ASTERISK, "*"},
		{SLASH, "/"},
		{EQ, "=="},
		{IDENTICAL, "==="},
		{DOUBLE_ARROW, "=>"},
		{OBJECT_OPERATOR, "->"},
		{PAAMAYIM_NEKUDOTAYIM, "::"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tokenType.String()
			if result != tt.expected {
				t.Errorf("TokenType.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident    string
		expected TokenType
	}{
		// Keywords
		{"if", IF},
		{"else", ELSE},
		{"elseif", ELSEIF},
		{"while", WHILE},
		{"for", FOR},
		{"foreach", FOREACH},
		{"function", FUNCTION},
		{"class", CLASS},
		{"interface", INTERFACE},
		{"trait", TRAIT},
		{"enum", ENUM},
		{"namespace", NAMESPACE},
		{"use", USE},
		{"public", PUBLIC},
		{"private", PRIVATE},
		{"protected", PROTECTED},
		{"static", STATIC},
		{"final", FINAL},
		{"abstract", ABSTRACT},
		{"readonly", READONLY},
		{"return", RETURN},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"echo", ECHO},
		{"print", PRINT},
		{"new", NEW},
		{"clone", CLONE},
		{"instanceof", INSTANCEOF},
		{"throw", THROW},
		{"try", TRY},
		{"catch", CATCH},
		{"finally", FINALLY},
		{"match", MATCH},
		{"yield", YIELD},
		{"fn", FN},

		// Type keywords
		{"int", INT},
		{"float", FLOAT_TYPE},
		{"bool", BOOL},
		{"string", STRING_TYPE},
		{"true", TRUE},
		{"false", FALSE},
		{"null", NULL},
		{"void", VOID},
		{"never", NEVER},
		{"mixed", MIXED},
		{"object", OBJECT},
		{"iterable", ITERABLE},

		// Magic constants
		{"__LINE__", LINE_CONST},
		{"__FILE__", FILE_CONST},
		{"__DIR__", DIR_CONST},
		{"__FUNCTION__", FUNCTION_CONST},
		{"__CLASS__", CLASS_CONST},
		{"__TRAIT__", TRAIT_CONST},
		{"__METHOD__", METHOD_CONST},
		{"__NAMESPACE__", NAMESPACE_CONST},

		// Non-keywords (should return IDENT)
		{"myFunction", IDENT},
		{"MyClass", IDENT},
		{"$variable", IDENT},
		{"foo", IDENT},
		{"bar", IDENT},
		{"test123", IDENT},
		{"_underscore", IDENT},
	}

	for _, tt := range tests {
		t.Run(tt.ident, func(t *testing.T) {
			result := LookupIdent(tt.ident)
			if result != tt.expected {
				t.Errorf("LookupIdent(%q) = %v, want %v", tt.ident, result, tt.expected)
			}
		})
	}
}

func TestIsKeyword(t *testing.T) {
	// Test that keywords are identified correctly
	keywords := []TokenType{
		IF, ELSE, ELSEIF, WHILE, FOR, FOREACH, FUNCTION, CLASS,
		INTERFACE, TRAIT, ENUM, NAMESPACE, USE, PUBLIC, PRIVATE,
		PROTECTED, STATIC, FINAL, ABSTRACT, RETURN, BREAK,
		CONTINUE, ECHO, PRINT, NEW, CLONE, INSTANCEOF,
	}

	for _, tok := range keywords {
		t.Run(tok.String(), func(t *testing.T) {
			if !tok.IsKeyword() {
				t.Errorf("%v.IsKeyword() = false, want true", tok)
			}
		})
	}

	// Test that non-keywords return false
	nonKeywords := []TokenType{
		INTEGER, FLOAT, STRING, VARIABLE, IDENT, PLUS, MINUS,
		LPAREN, RPAREN, SEMICOLON, EOF, ILLEGAL,
	}

	for _, tok := range nonKeywords {
		t.Run(tok.String(), func(t *testing.T) {
			if tok.IsKeyword() {
				t.Errorf("%v.IsKeyword() = true, want false", tok)
			}
		})
	}
}

func TestIsLiteral(t *testing.T) {
	// Test that literals are identified correctly
	literals := []TokenType{
		INTEGER, FLOAT, STRING, HEREDOC, NOWDOC,
	}

	for _, tok := range literals {
		t.Run(tok.String(), func(t *testing.T) {
			if !tok.IsLiteral() {
				t.Errorf("%v.IsLiteral() = false, want true", tok)
			}
		})
	}

	// Test that non-literals return false
	nonLiterals := []TokenType{
		VARIABLE, IDENT, IF, PLUS, LPAREN, SEMICOLON, EOF,
	}

	for _, tok := range nonLiterals {
		t.Run(tok.String(), func(t *testing.T) {
			if tok.IsLiteral() {
				t.Errorf("%v.IsLiteral() = true, want false", tok)
			}
		})
	}
}

func TestIsOperator(t *testing.T) {
	// Test that operators are identified correctly
	operators := []TokenType{
		PLUS, MINUS, ASTERISK, SLASH, PERCENT, POWER,
		EQ, IDENTICAL, NE, NOT_IDENTICAL, LT, LE, GT, GE,
		LOGICAL_AND, LOGICAL_OR, LOGICAL_NOT,
		BITWISE_AND, BITWISE_OR, BITWISE_XOR, BITWISE_NOT,
		ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN,
		DOUBLE_ARROW, OBJECT_OPERATOR, PAAMAYIM_NEKUDOTAYIM,
	}

	for _, tok := range operators {
		t.Run(tok.String(), func(t *testing.T) {
			if !tok.IsOperator() {
				t.Errorf("%v.IsOperator() = false, want true", tok)
			}
		})
	}

	// Test that non-operators return false
	nonOperators := []TokenType{
		INTEGER, VARIABLE, IF, LPAREN, SEMICOLON, EOF,
	}

	for _, tok := range nonOperators {
		t.Run(tok.String(), func(t *testing.T) {
			if tok.IsOperator() {
				t.Errorf("%v.IsOperator() = true, want false", tok)
			}
		})
	}
}

func TestPositionString(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		expected string
	}{
		{
			name:     "Position with filename",
			pos:      Position{Filename: "test.php", Line: 10, Column: 5},
			expected: "test.php:10:5",
		},
		{
			name:     "Position without filename",
			pos:      Position{Line: 20, Column: 15},
			expected: "20:15",
		},
		{
			name:     "Position at start",
			pos:      Position{Filename: "main.php", Line: 1, Column: 1},
			expected: "main.php:1:1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.String()
			if result != tt.expected {
				t.Errorf("Position.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPositionIsValid(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		expected bool
	}{
		{
			name:     "Valid position",
			pos:      Position{Line: 1, Column: 1},
			expected: true,
		},
		{
			name:     "Valid position with offset",
			pos:      Position{Line: 10, Column: 5, Offset: 100},
			expected: true,
		},
		{
			name:     "Invalid position (line 0)",
			pos:      Position{Line: 0, Column: 1},
			expected: false,
		},
		{
			name:     "Invalid position (negative line)",
			pos:      Position{Line: -1, Column: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsValid()
			if result != tt.expected {
				t.Errorf("Position.IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPositionComparison(t *testing.T) {
	tests := []struct {
		name     string
		pos1     Position
		pos2     Position
		before   bool
		after    bool
	}{
		{
			name:   "pos1 before pos2",
			pos1:   Position{Filename: "test.php", Offset: 10, Line: 1, Column: 10},
			pos2:   Position{Filename: "test.php", Offset: 20, Line: 2, Column: 5},
			before: true,
			after:  false,
		},
		{
			name:   "pos1 after pos2",
			pos1:   Position{Filename: "test.php", Offset: 30, Line: 3, Column: 1},
			pos2:   Position{Filename: "test.php", Offset: 15, Line: 2, Column: 1},
			before: false,
			after:  true,
		},
		{
			name:   "pos1 equals pos2",
			pos1:   Position{Filename: "test.php", Offset: 10, Line: 1, Column: 10},
			pos2:   Position{Filename: "test.php", Offset: 10, Line: 1, Column: 10},
			before: false,
			after:  false,
		},
		{
			name:   "different files (not comparable)",
			pos1:   Position{Filename: "file1.php", Offset: 10},
			pos2:   Position{Filename: "file2.php", Offset: 5},
			before: false,
			after:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeResult := tt.pos1.Before(tt.pos2)
			if beforeResult != tt.before {
				t.Errorf("Before() = %v, want %v", beforeResult, tt.before)
			}

			afterResult := tt.pos1.After(tt.pos2)
			if afterResult != tt.after {
				t.Errorf("After() = %v, want %v", afterResult, tt.after)
			}
		})
	}
}

func TestAllTokenTypesHaveNames(t *testing.T) {
	// Test that all token types have corresponding names
	// This helps catch missing entries in the tokenNames map
	for i := TokenType(0); i <= WHITESPACE; i++ {
		name := i.String()
		if name == "" {
			t.Errorf("TokenType %d is missing a name in tokenNames map", i)
		}
	}
}

func TestAllKeywordsAreInMap(t *testing.T) {
	// Verify that all keyword token types can be found through LookupIdent
	testKeywords := []struct {
		keyword string
		token   TokenType
	}{
		{"if", IF},
		{"class", CLASS},
		{"function", FUNCTION},
		{"namespace", NAMESPACE},
		{"readonly", READONLY},
		{"enum", ENUM},
		{"match", MATCH},
		{"fn", FN},
	}

	for _, tk := range testKeywords {
		t.Run(tk.keyword, func(t *testing.T) {
			result := LookupIdent(tk.keyword)
			if result != tk.token {
				t.Errorf("LookupIdent(%q) = %v, want %v", tk.keyword, result, tk.token)
			}
		})
	}
}
