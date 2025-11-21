package lexer

import "fmt"

// TokenType represents the type of a lexical token
type TokenType int

// Token represents a lexical token in PHP source code
type Token struct {
	Type    TokenType // Token type
	Literal string    // Actual text of the token
	Pos     Position  // Position in source code
}

// Position represents a location in the source code
type Position struct {
	Filename string // Source file name
	Offset   int    // Byte offset in source (0-based)
	Line     int    // Line number (1-based)
	Column   int    // Column number (1-based)
}

// Token type constants - PHP 8.4 compatible
const (
	// Special tokens
	ILLEGAL TokenType = iota // Illegal/unknown token
	EOF                      // End of file
	COMMENT                  // Comment (single-line or multi-line)

	// Literals
	INTEGER        // 123, 0x1A, 0b1010, 0o777
	FLOAT          // 123.45, 1.23e4
	STRING         // "string", 'string'
	HEREDOC        // <<<EOT ... EOT
	NOWDOC         // <<<'EOT' ... EOT
	ENCAPSED_START // Start of string interpolation
	ENCAPSED_END   // End of string interpolation

	// Identifiers and variables
	IDENT        // function_name, class_name
	VARIABLE     // $variable
	VAR_IDENT    // Variable name without $
	NS_SEPARATOR // \ in namespace paths

	// Keywords
	ABSTRACT
	AND
	ARRAY
	AS
	BREAK
	CALLABLE
	CASE
	CATCH
	CLASS
	CLONE
	CONST
	CONTINUE
	DECLARE
	DEFAULT
	DO
	ECHO
	ELSE
	ELSEIF
	EMPTY
	ENDDECLARE
	ENDFOR
	ENDFOREACH
	ENDIF
	ENDSWITCH
	ENDWHILE
	ENUM
	EVAL
	EXIT
	EXTENDS
	FINAL
	FINALLY
	FN          // Arrow function
	FOR
	FOREACH
	FUNCTION
	GLOBAL
	GOTO
	IF
	IMPLEMENTS
	INCLUDE
	INCLUDE_ONCE
	INSTANCEOF
	INSTEADOF
	INTERFACE
	ISSET
	LIST
	MATCH
	NAMESPACE
	NEW
	OR
	PRINT
	PRIVATE
	PROTECTED
	PUBLIC
	READONLY
	REQUIRE
	REQUIRE_ONCE
	RETURN
	STATIC
	SWITCH
	THROW
	TRAIT
	TRY
	UNSET
	USE
	VAR
	WHILE
	XOR
	YIELD
	YIELD_FROM

	// Type keywords
	INT
	FLOAT_TYPE
	BOOL
	STRING_TYPE
	TRUE
	FALSE
	NULL
	VOID
	NEVER
	MIXED
	OBJECT
	ITERABLE

	// Magic constants
	LINE_CONST      // __LINE__
	FILE_CONST      // __FILE__
	DIR_CONST       // __DIR__
	FUNCTION_CONST  // __FUNCTION__
	CLASS_CONST     // __CLASS__
	TRAIT_CONST     // __TRAIT__
	METHOD_CONST    // __METHOD__
	NAMESPACE_CONST // __NAMESPACE__
	PROPERTY_CONST  // __PROPERTY__

	// Operators and delimiters
	PLUS              // +
	MINUS             // -
	ASTERISK          // *
	SLASH             // /
	PERCENT           // %
	POWER             // **

	ASSIGN            // =
	PLUS_ASSIGN       // +=
	MINUS_ASSIGN      // -=
	MUL_ASSIGN        // *=
	DIV_ASSIGN        // /=
	MOD_ASSIGN        // %=
	CONCAT_ASSIGN     // .=
	POWER_ASSIGN      // **=
	AND_ASSIGN        // &=
	OR_ASSIGN         // |=
	XOR_ASSIGN        // ^=
	SL_ASSIGN         // <<=
	SR_ASSIGN         // >>=
	COALESCE_ASSIGN   // ??=

	EQ                // ==
	IDENTICAL         // ===
	NE                // !=
	NOT_IDENTICAL     // !==
	LT                // <
	LE                // <=
	GT                // >
	GE                // >=
	SPACESHIP         // <=>

	INC               // ++
	DEC               // --

	LOGICAL_AND       // &&
	LOGICAL_OR        // ||
	LOGICAL_NOT       // !

	BITWISE_AND       // &
	BITWISE_OR        // |
	BITWISE_XOR       // ^
	BITWISE_NOT       // ~
	SL                // <<
	SR                // >>

	CONCAT            // .

	QUESTION          // ?
	COLON             // :
	SEMICOLON         // ;
	COMMA             // ,
	DOT               // . (also concat)

	DOUBLE_ARROW      // =>
	OBJECT_OPERATOR   // ->
	PAAMAYIM_NEKUDOTAYIM // :: (scope resolution)
	NULLSAFE_OPERATOR // ?->
	ELLIPSIS          // ...
	COALESCE          // ??

	AT                // @ (error suppression)
	DOLLAR            // $
	AMPERSAND         // &
	BACKTICK          // ` (shell execution)

	// Delimiters
	LPAREN            // (
	RPAREN            // )
	LBRACE            // {
	RBRACE            // }
	LBRACKET          // [
	RBRACKET          // ]

	// PHP tags
	OPEN_TAG          // <?php or <?
	OPEN_TAG_ECHO     // <?=
	CLOSE_TAG         // ?>

	// Attributes (PHP 8.0+)
	ATTRIBUTE_START   // #[

	// Whitespace
	WHITESPACE        // Space, tab, newline
)

// String returns a human-readable representation of the token
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %s", t.Type.String(), t.Literal, t.Pos)
}

// String returns the name of the token type
func (tt TokenType) String() string {
	return tokenNames[tt]
}

// tokenNames maps token types to their string representations
var tokenNames = map[TokenType]string{
	ILLEGAL:           "ILLEGAL",
	EOF:               "EOF",
	COMMENT:           "COMMENT",

	INTEGER:           "INTEGER",
	FLOAT:             "FLOAT",
	STRING:            "STRING",
	HEREDOC:           "HEREDOC",
	NOWDOC:            "NOWDOC",
	ENCAPSED_START:    "ENCAPSED_START",
	ENCAPSED_END:      "ENCAPSED_END",

	IDENT:             "IDENT",
	VARIABLE:          "VARIABLE",
	VAR_IDENT:         "VAR_IDENT",
	NS_SEPARATOR:      "NS_SEPARATOR",

	ABSTRACT:          "ABSTRACT",
	AND:               "AND",
	ARRAY:             "ARRAY",
	AS:                "AS",
	BREAK:             "BREAK",
	CALLABLE:          "CALLABLE",
	CASE:              "CASE",
	CATCH:             "CATCH",
	CLASS:             "CLASS",
	CLONE:             "CLONE",
	CONST:             "CONST",
	CONTINUE:          "CONTINUE",
	DECLARE:           "DECLARE",
	DEFAULT:           "DEFAULT",
	DO:                "DO",
	ECHO:              "ECHO",
	ELSE:              "ELSE",
	ELSEIF:            "ELSEIF",
	EMPTY:             "EMPTY",
	ENDDECLARE:        "ENDDECLARE",
	ENDFOR:            "ENDFOR",
	ENDFOREACH:        "ENDFOREACH",
	ENDIF:             "ENDIF",
	ENDSWITCH:         "ENDSWITCH",
	ENDWHILE:          "ENDWHILE",
	ENUM:              "ENUM",
	EVAL:              "EVAL",
	EXIT:              "EXIT",
	EXTENDS:           "EXTENDS",
	FINAL:             "FINAL",
	FINALLY:           "FINALLY",
	FN:                "FN",
	FOR:               "FOR",
	FOREACH:           "FOREACH",
	FUNCTION:          "FUNCTION",
	GLOBAL:            "GLOBAL",
	GOTO:              "GOTO",
	IF:                "IF",
	IMPLEMENTS:        "IMPLEMENTS",
	INCLUDE:           "INCLUDE",
	INCLUDE_ONCE:      "INCLUDE_ONCE",
	INSTANCEOF:        "INSTANCEOF",
	INSTEADOF:         "INSTEADOF",
	INTERFACE:         "INTERFACE",
	ISSET:             "ISSET",
	LIST:              "LIST",
	MATCH:             "MATCH",
	NAMESPACE:         "NAMESPACE",
	NEW:               "NEW",
	OR:                "OR",
	PRINT:             "PRINT",
	PRIVATE:           "PRIVATE",
	PROTECTED:         "PROTECTED",
	PUBLIC:            "PUBLIC",
	READONLY:          "READONLY",
	REQUIRE:           "REQUIRE",
	REQUIRE_ONCE:      "REQUIRE_ONCE",
	RETURN:            "RETURN",
	STATIC:            "STATIC",
	SWITCH:            "SWITCH",
	THROW:             "THROW",
	TRAIT:             "TRAIT",
	TRY:               "TRY",
	UNSET:             "UNSET",
	USE:               "USE",
	VAR:               "VAR",
	WHILE:             "WHILE",
	XOR:               "XOR",
	YIELD:             "YIELD",
	YIELD_FROM:        "YIELD_FROM",

	INT:               "INT",
	FLOAT_TYPE:        "FLOAT_TYPE",
	BOOL:              "BOOL",
	STRING_TYPE:       "STRING_TYPE",
	TRUE:              "TRUE",
	FALSE:             "FALSE",
	NULL:              "NULL",
	VOID:              "VOID",
	NEVER:             "NEVER",
	MIXED:             "MIXED",
	OBJECT:            "OBJECT",
	ITERABLE:          "ITERABLE",

	LINE_CONST:        "__LINE__",
	FILE_CONST:        "__FILE__",
	DIR_CONST:         "__DIR__",
	FUNCTION_CONST:    "__FUNCTION__",
	CLASS_CONST:       "__CLASS__",
	TRAIT_CONST:       "__TRAIT__",
	METHOD_CONST:      "__METHOD__",
	NAMESPACE_CONST:   "__NAMESPACE__",
	PROPERTY_CONST:    "__PROPERTY__",

	PLUS:              "+",
	MINUS:             "-",
	ASTERISK:          "*",
	SLASH:             "/",
	PERCENT:           "%",
	POWER:             "**",

	ASSIGN:            "=",
	PLUS_ASSIGN:       "+=",
	MINUS_ASSIGN:      "-=",
	MUL_ASSIGN:        "*=",
	DIV_ASSIGN:        "/=",
	MOD_ASSIGN:        "%=",
	CONCAT_ASSIGN:     ".=",
	POWER_ASSIGN:      "**=",
	AND_ASSIGN:        "&=",
	OR_ASSIGN:         "|=",
	XOR_ASSIGN:        "^=",
	SL_ASSIGN:         "<<=",
	SR_ASSIGN:         ">>=",
	COALESCE_ASSIGN:   "??=",

	EQ:                "==",
	IDENTICAL:         "===",
	NE:                "!=",
	NOT_IDENTICAL:     "!==",
	LT:                "<",
	LE:                "<=",
	GT:                ">",
	GE:                ">=",
	SPACESHIP:         "<=>",

	INC:               "++",
	DEC:               "--",

	LOGICAL_AND:       "&&",
	LOGICAL_OR:        "||",
	LOGICAL_NOT:       "!",

	BITWISE_AND:       "&",
	BITWISE_OR:        "|",
	BITWISE_XOR:       "^",
	BITWISE_NOT:       "~",
	SL:                "<<",
	SR:                ">>",

	CONCAT:            ".",

	QUESTION:          "?",
	COLON:             ":",
	SEMICOLON:         ";",
	COMMA:             ",",
	DOT:               ".",

	DOUBLE_ARROW:      "=>",
	OBJECT_OPERATOR:   "->",
	PAAMAYIM_NEKUDOTAYIM: "::",
	NULLSAFE_OPERATOR: "?->",
	ELLIPSIS:          "...",
	COALESCE:          "??",

	AT:                "@",
	DOLLAR:            "$",
	AMPERSAND:         "&",
	BACKTICK:          "`",

	LPAREN:            "(",
	RPAREN:            ")",
	LBRACE:            "{",
	RBRACE:            "}",
	LBRACKET:          "[",
	RBRACKET:          "]",

	OPEN_TAG:          "<?php",
	OPEN_TAG_ECHO:     "<?=",
	CLOSE_TAG:         "?>",

	ATTRIBUTE_START:   "#[",

	WHITESPACE:        "WHITESPACE",
}

// keywords maps PHP keywords to their token types
var keywords = map[string]TokenType{
	"abstract":      ABSTRACT,
	"and":           AND,
	"array":         ARRAY,
	"as":            AS,
	"break":         BREAK,
	"callable":      CALLABLE,
	"case":          CASE,
	"catch":         CATCH,
	"class":         CLASS,
	"clone":         CLONE,
	"const":         CONST,
	"continue":      CONTINUE,
	"declare":       DECLARE,
	"default":       DEFAULT,
	"do":            DO,
	"echo":          ECHO,
	"else":          ELSE,
	"elseif":        ELSEIF,
	"empty":         EMPTY,
	"enddeclare":    ENDDECLARE,
	"endfor":        ENDFOR,
	"endforeach":    ENDFOREACH,
	"endif":         ENDIF,
	"endswitch":     ENDSWITCH,
	"endwhile":      ENDWHILE,
	"enum":          ENUM,
	"eval":          EVAL,
	"exit":          EXIT,
	"extends":       EXTENDS,
	"final":         FINAL,
	"finally":       FINALLY,
	"fn":            FN,
	"for":           FOR,
	"foreach":       FOREACH,
	"function":      FUNCTION,
	"global":        GLOBAL,
	"goto":          GOTO,
	"if":            IF,
	"implements":    IMPLEMENTS,
	"include":       INCLUDE,
	"include_once":  INCLUDE_ONCE,
	"instanceof":    INSTANCEOF,
	"insteadof":     INSTEADOF,
	"interface":     INTERFACE,
	"isset":         ISSET,
	"list":          LIST,
	"match":         MATCH,
	"namespace":     NAMESPACE,
	"new":           NEW,
	"or":            OR,
	"print":         PRINT,
	"private":       PRIVATE,
	"protected":     PROTECTED,
	"public":        PUBLIC,
	"readonly":      READONLY,
	"require":       REQUIRE,
	"require_once":  REQUIRE_ONCE,
	"return":        RETURN,
	"static":        STATIC,
	"switch":        SWITCH,
	"throw":         THROW,
	"trait":         TRAIT,
	"try":           TRY,
	"unset":         UNSET,
	"use":           USE,
	"var":           VAR,
	"while":         WHILE,
	"xor":           XOR,
	"yield":         YIELD,

	// Type keywords
	"int":           INT,
	"float":         FLOAT_TYPE,
	"bool":          BOOL,
	"string":        STRING_TYPE,
	"true":          TRUE,
	"false":         FALSE,
	"null":          NULL,
	"void":          VOID,
	"never":         NEVER,
	"mixed":         MIXED,
	"object":        OBJECT,
	"iterable":      ITERABLE,

	// Magic constants
	"__LINE__":      LINE_CONST,
	"__FILE__":      FILE_CONST,
	"__DIR__":       DIR_CONST,
	"__FUNCTION__":  FUNCTION_CONST,
	"__CLASS__":     CLASS_CONST,
	"__TRAIT__":     TRAIT_CONST,
	"__METHOD__":    METHOD_CONST,
	"__NAMESPACE__": NAMESPACE_CONST,
	"__PROPERTY__":  PROPERTY_CONST,
}

// LookupIdent checks if an identifier is a keyword and returns the appropriate token type
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// IsKeyword returns true if the token type is a keyword
func (tt TokenType) IsKeyword() bool {
	return tt >= ABSTRACT && tt <= YIELD_FROM
}

// IsLiteral returns true if the token type is a literal
func (tt TokenType) IsLiteral() bool {
	return tt >= INTEGER && tt <= ENCAPSED_END
}

// IsOperator returns true if the token type is an operator
func (tt TokenType) IsOperator() bool {
	// Arithmetic, comparison, logical, bitwise, and assignment operators
	if tt >= PLUS && tt <= CONCAT {
		return true
	}
	// PHP-specific operators (=>, ->, ::, ?->, ..., ??)
	if tt >= DOUBLE_ARROW && tt <= COALESCE {
		return true
	}
	// Special operators
	return tt == AT || tt == AMPERSAND || tt == BACKTICK
}

// Position methods

// String returns a formatted string representation of the position
func (p Position) String() string {
	if p.Filename != "" {
		return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
	}
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// IsValid returns true if the position is valid (line > 0)
func (p Position) IsValid() bool {
	return p.Line > 0
}

// Before returns true if this position is before the other position
func (p Position) Before(other Position) bool {
	if p.Filename != other.Filename {
		return false
	}
	return p.Offset < other.Offset
}

// After returns true if this position is after the other position
func (p Position) After(other Position) bool {
	if p.Filename != other.Filename {
		return false
	}
	return p.Offset > other.Offset
}
