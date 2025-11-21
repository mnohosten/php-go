package lexer

import (
	"strings"
)

// scanHeredoc scans a heredoc string (<<<EOT ... EOT)
func (l *Lexer) scanHeredoc(label string, isNowdoc bool) Token {
	pos := l.currentPosition()
	var result strings.Builder

	// Skip the newline after <<<LABEL
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.lineStart = l.pos + 1
		l.readChar()
	} else if l.ch == '\r' && l.peekChar() == '\n' {
		l.readChar() // skip \r
		l.line++
		l.column = 0
		l.lineStart = l.pos + 1
		l.readChar() // skip \n
	}

	// Scan until we find the closing label on its own line
	for l.ch != 0 {
		// Check if this line starts with the closing label
		if l.checkHeredocEnd(label) {
			break
		}

		// Read the entire line
		for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
			if !isNowdoc && l.ch == '$' {
				// TODO: Handle interpolation in heredoc
				// For now, just include the $ in the string
				result.WriteByte(l.ch)
				l.readChar()
			} else if !isNowdoc && l.ch == '\\' && (l.peekChar() == '\\' || l.peekChar() == '$') {
				// Escape sequences in heredoc
				l.readChar()
				result.WriteByte(l.ch)
				l.readChar()
			} else {
				result.WriteByte(l.ch)
				l.readChar()
			}
		}

		// Include newline in heredoc content (except for the last line)
		if l.ch == '\n' {
			result.WriteByte('\n')
			l.line++
			l.column = 0
			l.lineStart = l.pos + 1
			l.readChar()
		} else if l.ch == '\r' {
			if l.peekChar() == '\n' {
				result.WriteByte('\n')
				l.readChar() // skip \r
				l.readChar() // skip \n
			} else {
				result.WriteByte('\n')
				l.readChar()
			}
			l.line++
			l.column = 0
			l.lineStart = l.pos + 1
		}
	}

	if l.ch == 0 {
		return Token{
			Type:    ILLEGAL,
			Literal: "unterminated heredoc",
			Pos:     pos,
		}
	}

	// Consume the closing label
	for l.ch != '\n' && l.ch != '\r' && l.ch != ';' && l.ch != 0 {
		l.readChar()
	}

	// Don't consume semicolon - it's part of the statement syntax
	// Just consume trailing newline if present
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.lineStart = l.pos + 1
		l.readChar()
	} else if l.ch == '\r' {
		l.readChar()
		if l.ch == '\n' {
			l.readChar()
		}
		l.line++
		l.column = 0
		l.lineStart = l.pos + 1
	}

	tokenType := HEREDOC
	if isNowdoc {
		tokenType = NOWDOC
	}

	return Token{
		Type:    tokenType,
		Literal: result.String(),
		Pos:     pos,
	}
}

// checkHeredocEnd checks if the current position is at the heredoc closing label
func (l *Lexer) checkHeredocEnd(label string) bool {
	// Save current position
	savedPos := l.pos
	savedReadPos := l.readPos
	savedCh := l.ch
	savedLine := l.line
	savedColumn := l.column

	// Skip leading whitespace (PHP 7.3+ allows indented closing tags)
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}

	// Check if the label matches
	matches := true
	for i := 0; i < len(label); i++ {
		if l.ch != label[i] {
			matches = false
			break
		}
		l.readChar()
	}

	// After the label, we should see ; or newline or EOF
	if matches {
		if l.ch == ';' || l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
			// Restore position (we'll consume it in scanHeredoc)
			l.pos = savedPos
			l.readPos = savedReadPos
			l.ch = savedCh
			l.line = savedLine
			l.column = savedColumn
			return true
		}
	}

	// Restore position
	l.pos = savedPos
	l.readPos = savedReadPos
	l.ch = savedCh
	l.line = savedLine
	l.column = savedColumn
	return false
}

// scanHeredocLabel scans the heredoc/nowdoc label after <<<
func (l *Lexer) scanHeredocLabel() (string, bool) {
	isNowdoc := false

	// Skip whitespace after <<<
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}

	// Check for nowdoc (quoted label)
	if l.ch == '\'' || l.ch == '"' {
		isNowdoc = true
		quote := l.ch
		l.readChar()

		start := l.pos
		for l.ch != quote && l.ch != 0 && l.ch != '\n' {
			l.readChar()
		}

		label := l.input[start:l.pos]

		if l.ch != quote {
			return "", isNowdoc
		}
		l.readChar() // consume closing quote

		return label, isNowdoc
	}

	// Unquoted label (heredoc)
	start := l.pos
	if !isLetter(l.ch) && l.ch != '_' {
		return "", isNowdoc
	}

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	label := l.input[start:l.pos]
	return label, isNowdoc
}

// StringInterpolation represents a part of an interpolated string
type StringInterpolation struct {
	IsVariable bool   // true if this is a variable, false if it's text
	Value      string // the variable name (with $) or the text content
}

// parseStringInterpolation parses string interpolation in double-quoted strings
// Returns a list of parts (text and variables)
func parseStringInterpolation(input string) []StringInterpolation {
	var parts []StringInterpolation
	var current strings.Builder
	i := 0

	for i < len(input) {
		ch := input[i]

		if ch == '$' && i+1 < len(input) {
			// Check if this is a variable
			next := input[i+1]
			if isLetter(byte(next)) || next == '_' {
				// Save any pending text
				if current.Len() > 0 {
					parts = append(parts, StringInterpolation{
						IsVariable: false,
						Value:      current.String(),
					})
					current.Reset()
				}

				// Extract variable name
				varStart := i
				i++ // skip $
				for i < len(input) && (isLetter(input[i]) || isDigit(input[i]) || input[i] == '_') {
					i++
				}

				parts = append(parts, StringInterpolation{
					IsVariable: true,
					Value:      input[varStart:i],
				})
				continue
			} else if next == '{' {
				// Complex interpolation like ${expr} or {$expr}
				// For now, treat as regular text (will be enhanced later)
				current.WriteByte(ch)
				i++
				continue
			}
		}

		current.WriteByte(ch)
		i++
	}

	// Add any remaining text
	if current.Len() > 0 {
		parts = append(parts, StringInterpolation{
			IsVariable: false,
			Value:      current.String(),
		})
	}

	return parts
}

// hasInterpolation checks if a string contains variable interpolation
func hasInterpolation(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '$' && i+1 < len(s) {
			next := s[i+1]
			if isLetter(byte(next)) || next == '_' || next == '{' {
				return true
			}
		}
	}
	return false
}

// isValidHeredocLabel checks if a string is a valid heredoc/nowdoc label
func isValidHeredocLabel(label string) bool {
	if len(label) == 0 {
		return false
	}

	// First character must be a letter or underscore
	first := label[0]
	if !isLetter(first) && first != '_' {
		return false
	}

	// Rest must be letters, digits, or underscores
	for i := 1; i < len(label); i++ {
		ch := label[i]
		if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			return false
		}
	}

	return true
}

// scanStringWithInterpolation scans a double-quoted string and handles interpolation
// Returns ENCAPSED_START token if interpolation is found, otherwise returns STRING
func (l *Lexer) scanStringWithInterpolation() Token {
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
			case '0':
				result.WriteByte(0)
			case 'x': // Hex escape \xHH
				l.readChar()
				hex1 := l.ch
				l.readChar()
				hex2 := l.ch
				if isHexDigit(hex1) && isHexDigit(hex2) {
					val := hexValue(hex1)*16 + hexValue(hex2)
					result.WriteByte(byte(val))
				} else {
					result.WriteString("\\x")
					result.WriteByte(hex1)
					result.WriteByte(hex2)
				}
				l.readChar()
				continue
			default:
				result.WriteByte('\\')
				result.WriteByte(l.ch)
			}
			l.readChar()
		} else if l.ch == '$' && l.peekChar() != 0 {
			// Check if this is a variable interpolation
			// For now, just include the $ in the string
			// Full interpolation support will be added later
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

	// Consume closing quote
	l.readChar()

	// For now, return STRING even if it has interpolation
	// Full interpolation support will tokenize variables separately
	return Token{
		Type:    STRING,
		Literal: result.String(),
		Pos:     pos,
	}
}

// hexValue converts a hex digit character to its numeric value
func hexValue(ch byte) byte {
	if ch >= '0' && ch <= '9' {
		return ch - '0'
	} else if ch >= 'a' && ch <= 'f' {
		return ch - 'a' + 10
	} else if ch >= 'A' && ch <= 'F' {
		return ch - 'A' + 10
	}
	return 0
}
