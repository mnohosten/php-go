package lexer

import (
	"fmt"
	"strings"
)

// Lexer tokenizes PHP source code into a stream of tokens
type Lexer struct {
	input     string // The input string being scanned
	filename  string // Source filename (for error messages)
	pos       int    // Current position in input (byte offset)
	readPos   int    // Current reading position (pos + 1)
	ch        byte   // Current character under examination
	line      int    // Current line number (1-based)
	column    int    // Current column number (1-based)
	lineStart int    // Byte offset of the start of the current line
	inPHP     bool   // Whether we're currently inside PHP mode (between <?php and ?>)
}

// New creates a new Lexer for the given input
func New(input, filename string) *Lexer {
	l := &Lexer{
		input:    input,
		filename: filename,
		line:     1,
		column:   0, // Start at 0, readChar will increment to 1
		inPHP:    true, // Default to PHP mode for backward compatibility
	}
	l.readChar() // Initialize first character

	// Determine initial mode based on file content
	// If file starts with HTML (< but not <?), start in HTML mode
	// Otherwise, start in PHP mode for backward compatibility
	if len(input) >= 2 && input[0] == '<' && input[1] != '?' {
		// Starts with HTML tag (e.g., <html, <div, etc.)
		l.inPHP = false
	}
	// All other cases: start in PHP mode
	// - Pure PHP code without tags
	// - <?php or <? tags
	// - Code starting with PHP constructs (<<<, $, etc.)

	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII NUL - represents EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// peekCharN returns the character n positions ahead without advancing
func (l *Lexer) peekCharN(n int) byte {
	pos := l.readPos + n - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// currentPosition returns the current source position
func (l *Lexer) currentPosition() Position {
	return Position{
		Filename: l.filename,
		Offset:   l.pos,
		Line:     l.line,
		Column:   l.column,
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	// If we're not in PHP mode, scan inline HTML
	if !l.inPHP && l.ch != 0 {
		// Check if this is a PHP open tag
		if l.ch == '<' && l.peekChar() == '?' {
			// This is a PHP open tag, let it be processed below
			// Don't scan as inline HTML
		} else {
			// We're in HTML mode - scan inline HTML until we find <?php or <?=
			return l.scanInlineHTML()
		}
	}

	l.skipWhitespace()

	tok.Pos = l.currentPosition()

	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Literal = ""

	// Single character tokens
	case ';':
		tok = l.makeToken(SEMICOLON, string(l.ch))
	case ',':
		tok = l.makeToken(COMMA, string(l.ch))
	case '(':
		tok = l.makeToken(LPAREN, string(l.ch))
	case ')':
		tok = l.makeToken(RPAREN, string(l.ch))
	case '{':
		tok = l.makeToken(LBRACE, string(l.ch))
	case '}':
		tok = l.makeToken(RBRACE, string(l.ch))
	case '[':
		tok = l.makeToken(LBRACKET, string(l.ch))
	case ']':
		tok = l.makeToken(RBRACKET, string(l.ch))
	case '~':
		tok = l.makeToken(BITWISE_NOT, string(l.ch))
	case '@':
		tok = l.makeToken(AT, string(l.ch))
	case '`':
		tok = l.scanBacktickString()
		return tok

	// Operators (potentially multi-character)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(INC, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(PLUS_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(PLUS, string(l.ch))
		}

	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(DEC, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(MINUS_ASSIGN, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(OBJECT_OPERATOR, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(MINUS, string(l.ch))
		}

	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(POWER_ASSIGN, "**=")
			} else {
				tok = l.makeToken(POWER, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(MUL_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(ASTERISK, string(l.ch))
		}

	case '/':
		if l.peekChar() == '/' {
			tok = l.scanSingleLineComment()
			return tok
		} else if l.peekChar() == '*' {
			tok = l.scanMultiLineComment()
			return tok
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(DIV_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(SLASH, string(l.ch))
		}

	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(MOD_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(PERCENT, string(l.ch))
		}

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(IDENTICAL, "===")
			} else {
				tok = l.makeToken(EQ, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(DOUBLE_ARROW, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(ASSIGN, string(l.ch))
		}

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(NOT_IDENTICAL, "!==")
			} else {
				tok = l.makeToken(NE, string(ch)+string(l.ch))
			}
		} else {
			tok = l.makeToken(LOGICAL_NOT, string(l.ch))
		}

	case '<':
		if l.peekChar() == '<' {
			ch := l.ch
			l.readChar() // Now on second <
			if l.peekChar() == '<' {
				// Could be heredoc/nowdoc (<<<)
				l.readChar() // Consume second <, now on third <
				tok = l.scanHeredocOrShift()
				return tok
			} else if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(SL_ASSIGN, "<<=")
			} else {
				tok = l.makeToken(SL, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '>' {
				l.readChar()
				tok = l.makeToken(SPACESHIP, "<=>")
			} else {
				tok = l.makeToken(LE, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '?' {
			// Check for <?php or <?=
			tok = l.scanPHPTag()
			return tok
		} else {
			tok = l.makeToken(LT, string(l.ch))
		}

	case '>':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(SR_ASSIGN, ">>=")
			} else {
				tok = l.makeToken(SR, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(GE, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(GT, string(l.ch))
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(LOGICAL_AND, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(AND_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(BITWISE_AND, string(l.ch))
		}

	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(LOGICAL_OR, string(ch)+string(l.ch))
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(OR_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(BITWISE_OR, string(l.ch))
		}

	case '^':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(XOR_ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(BITWISE_XOR, string(l.ch))
		}

	case '.':
		if l.peekChar() == '.' && l.peekCharN(2) == '.' {
			l.readChar()
			l.readChar()
			tok = l.makeToken(ELLIPSIS, "...")
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(CONCAT_ASSIGN, string(ch)+string(l.ch))
		} else if isDigit(l.peekChar()) {
			// Could be a float like .5
			tok = l.scanNumber()
			return tok
		} else {
			tok = l.makeToken(CONCAT, string(l.ch))
		}

	case '?':
		if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.makeToken(COALESCE_ASSIGN, "??=")
			} else {
				tok = l.makeToken(COALESCE, string(ch)+string(l.ch))
			}
		} else if l.peekChar() == '-' && l.peekCharN(2) == '>' {
			l.readChar()
			l.readChar()
			tok = l.makeToken(NULLSAFE_OPERATOR, "?->")
		} else if l.peekChar() == '>' {
			// PHP close tag ?>
			l.readChar()
			l.inPHP = false // Exit PHP mode
			tok = l.makeToken(CLOSE_TAG, "?>")
		} else {
			tok = l.makeToken(QUESTION, string(l.ch))
		}

	case ':':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(PAAMAYIM_NEKUDOTAYIM, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(COLON, string(l.ch))
		}

	case '$':
		tok = l.scanVariable()
		return tok

	case '"':
		tok = l.scanStringWithInterpolation()
		return tok

	case '\'':
		tok = l.scanSingleQuotedString()
		return tok

	case '#':
		// Check for attribute #[
		if l.peekChar() == '[' {
			l.readChar()
			tok = l.makeToken(ATTRIBUTE_START, "#[")
		} else {
			// Single-line comment
			tok = l.scanSingleLineComment()
			return tok
		}

	default:
		if isLetter(l.ch) || l.ch == '_' || l.ch == '\\' {
			tok = l.scanIdentifier()
			return tok
		} else if isDigit(l.ch) {
			tok = l.scanNumber()
			return tok
		} else {
			tok = l.makeToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// makeToken creates a token with the given type and literal
func (l *Lexer) makeToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Pos:     l.currentPosition(),
	}
}

// skipWhitespace skips whitespace characters (space, tab, newline, carriage return)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
			l.lineStart = l.pos + 1
		}
		l.readChar()
	}
}

// scanIdentifier scans an identifier or keyword
func (l *Lexer) scanIdentifier() Token {
	pos := l.currentPosition()
	start := l.pos

	// PHP identifiers can start with a letter, underscore, or backslash (for namespaces)
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '\\' {
		l.readChar()
	}

	literal := l.input[start:l.pos]

	// Check if it's a keyword
	tokenType := LookupIdent(literal)

	return Token{
		Type:    tokenType,
		Literal: literal,
		Pos:     pos,
	}
}

// scanVariable scans a PHP variable ($var)
func (l *Lexer) scanVariable() Token {
	pos := l.currentPosition()
	start := l.pos

	l.readChar() // consume '$'

	if !isLetter(l.ch) && l.ch != '_' {
		return Token{
			Type:    ILLEGAL,
			Literal: "$",
			Pos:     pos,
		}
	}

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	literal := l.input[start:l.pos]

	return Token{
		Type:    VARIABLE,
		Literal: literal,
		Pos:     pos,
	}
}

// scanNumber scans an integer or float literal
func (l *Lexer) scanNumber() Token {
	pos := l.currentPosition()
	start := l.pos
	tokenType := INTEGER

	// Handle hex (0x), octal (0o), binary (0b)
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar()
		l.readChar()
		for isHexDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	} else if l.ch == '0' && (l.peekChar() == 'b' || l.peekChar() == 'B') {
		l.readChar()
		l.readChar()
		for l.ch == '0' || l.ch == '1' || l.ch == '_' {
			l.readChar()
		}
	} else if l.ch == '0' && (l.peekChar() == 'o' || l.peekChar() == 'O') {
		l.readChar()
		l.readChar()
		for isOctalDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	} else {
		// Decimal number or float starting with .
		if l.ch == '.' {
			tokenType = FLOAT
			l.readChar() // consume '.'
		}

		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}

		// Check for decimal point (if not already processed)
		if l.ch == '.' && tokenType == INTEGER && isDigit(l.peekChar()) {
			tokenType = FLOAT
			l.readChar() // consume '.'
			for isDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
		}

		// Check for exponent (e or E)
		if l.ch == 'e' || l.ch == 'E' {
			tokenType = FLOAT
			l.readChar()
			if l.ch == '+' || l.ch == '-' {
				l.readChar()
			}
			for isDigit(l.ch) || l.ch == '_' {
				l.readChar()
			}
		}
	}

	literal := l.input[start:l.pos]

	return Token{
		Type:    tokenType,
		Literal: literal,
		Pos:     pos,
	}
}

// scanSingleQuotedString scans a single-quoted string (no interpolation)
func (l *Lexer) scanSingleQuotedString() Token {
	pos := l.currentPosition()
	var result strings.Builder

	l.readChar() // consume opening '

	for l.ch != '\'' && l.ch != 0 {
		if l.ch == '\\' && (l.peekChar() == '\'' || l.peekChar() == '\\') {
			l.readChar()
			result.WriteByte(l.ch)
			l.readChar()
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			result.WriteByte(l.ch)
			l.readChar()
		}
	}

	if l.ch == 0 {
		return Token{
			Type:    ILLEGAL,
			Literal: "unterminated string",
			Pos:     pos,
		}
	}

	// l.ch is now on the closing quote; consume it
	l.readChar()

	return Token{
		Type:    STRING,
		Literal: result.String(),
		Pos:     pos,
	}
}

// scanDoubleQuotedString scans a double-quoted string (with interpolation support)
func (l *Lexer) scanDoubleQuotedString() Token {
	pos := l.currentPosition()
	var result strings.Builder

	l.readChar() // consume opening "

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '$':
				result.WriteByte('$')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte('\\')
				result.WriteByte(l.ch)
			}
			l.readChar()
		} else if l.ch == '$' {
			// TODO: Handle string interpolation in Phase 1, Task 1.4
			// For now, just include it in the string
			result.WriteByte(l.ch)
			l.readChar()
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			result.WriteByte(l.ch)
			l.readChar()
		}
	}

	if l.ch == 0 {
		return Token{
			Type:    ILLEGAL,
			Literal: "unterminated string",
			Pos:     pos,
		}
	}

	// l.ch is now on the closing quote; consume it
	l.readChar()

	return Token{
		Type:    STRING,
		Literal: result.String(),
		Pos:     pos,
	}
}

// scanBacktickString scans a backtick string (shell execution)
func (l *Lexer) scanBacktickString() Token {
	pos := l.currentPosition()
	var result strings.Builder

	l.readChar() // consume opening `

	for l.ch != '`' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		result.WriteByte(l.ch)
		l.readChar()
	}

	if l.ch == 0 {
		return Token{
			Type:    ILLEGAL,
			Literal: "unterminated shell execution string",
			Pos:     pos,
		}
	}

	// l.ch is now on the closing backtick; consume it
	l.readChar()

	return Token{
		Type:    STRING, // For now, treat as string; will be handled specially in parser
		Literal: result.String(),
		Pos:     pos,
	}
}

// scanSingleLineComment scans a single-line comment (// or #)
func (l *Lexer) scanSingleLineComment() Token {
	pos := l.currentPosition()
	start := l.pos

	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	literal := l.input[start:l.pos]

	return Token{
		Type:    COMMENT,
		Literal: literal,
		Pos:     pos,
	}
}

// scanMultiLineComment scans a multi-line comment (/* ... */)
func (l *Lexer) scanMultiLineComment() Token {
	pos := l.currentPosition()
	start := l.pos

	l.readChar() // consume '/'
	l.readChar() // consume '*'

	for {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		if l.ch == 0 {
			return Token{
				Type:    ILLEGAL,
				Literal: "unterminated comment",
				Pos:     pos,
			}
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}

	literal := l.input[start:l.pos]

	return Token{
		Type:    COMMENT,
		Literal: literal,
		Pos:     pos,
	}
}

// scanPHPTag scans PHP opening tags (<?php, <?, <?=)
func (l *Lexer) scanPHPTag() Token {
	pos := l.currentPosition()

	l.readChar() // consume '<'
	l.readChar() // consume '?'

	if l.ch == '=' {
		l.readChar()
		l.inPHP = true // Enter PHP mode
		return Token{
			Type:    OPEN_TAG_ECHO,
			Literal: "<?=",
			Pos:     pos,
		}
	}

	if l.ch == 'p' && l.peekChar() == 'h' && l.peekCharN(2) == 'p' {
		l.readChar()
		l.readChar()
		l.readChar()
		// Optionally consume following whitespace
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.skipWhitespace()
		}
		l.inPHP = true // Enter PHP mode
		return Token{
			Type:    OPEN_TAG,
			Literal: "<?php",
			Pos:     pos,
		}
	}

	// Short open tag <?
	l.inPHP = true // Enter PHP mode
	return Token{
		Type:    OPEN_TAG,
		Literal: "<?",
		Pos:     pos,
	}
}

// scanHeredocOrShift determines if <<< is heredoc/nowdoc or left shift
func (l *Lexer) scanHeredocOrShift() Token {
	// We're already on the third '<'
	// Consume it and check for heredoc label
	l.readChar() // Move past third '<'

	// Try to scan heredoc/nowdoc label
	label, isNowdoc := l.scanHeredocLabel()

	if label == "" || !isValidHeredocLabel(label) {
		// Not a valid heredoc, this is an error case
		// We already consumed <<<, so we can't easily recover
		// Return illegal token
		return l.makeToken(ILLEGAL, "invalid heredoc/nowdoc syntax")
	}

	// Valid heredoc/nowdoc found
	return l.scanHeredoc(label, isNowdoc)
}

// Helper functions

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch >= 128
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isOctalDigit(ch byte) bool {
	return ch >= '0' && ch <= '7'
}

// scanInlineHTML scans HTML content between ?> and <?php
func (l *Lexer) scanInlineHTML() Token {
	pos := l.currentPosition()
	startPos := l.pos

	// Scan until we find <?php, <?=, <?, or EOF
	for {
		if l.ch == 0 {
			// EOF - return all remaining content as inline HTML
			break
		}

		if l.ch == '<' && l.peekChar() == '?' {
			// Check if it's <?php, <?=, or <?
			if l.peekCharN(2) == 'p' && l.peekCharN(3) == 'h' && l.peekCharN(4) == 'p' {
				// Found <?php - stop here
				break
			} else if l.peekCharN(2) == '=' {
				// Found <?= - stop here
				break
			} else {
				// Found <? (short tag) - stop here
				break
			}
		}

		// Track line numbers in HTML content
		if l.ch == '\n' {
			l.line++
			l.column = 0
			l.lineStart = l.readPos
		}

		l.readChar()
	}

	// Extract the HTML content
	html := l.input[startPos:l.pos]

	return Token{
		Type:    INLINE_HTML,
		Literal: html,
		Pos:     pos,
	}
}

// Error creates an error token
func (l *Lexer) Error(message string) Token {
	return Token{
		Type:    ILLEGAL,
		Literal: message,
		Pos:     l.currentPosition(),
	}
}

// String returns a string representation of the lexer state (for debugging)
func (l *Lexer) String() string {
	return fmt.Sprintf("Lexer{pos:%d, line:%d, col:%d, ch:%q}", l.pos, l.line, l.column, l.ch)
}
