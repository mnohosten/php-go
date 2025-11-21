package parser

import (
	"fmt"
	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// Parser parses PHP source code into an Abstract Syntax Tree (AST)
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	// For error recovery
	panicMode bool
}

// New creates a new Parser from a lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// ParseProgram parses the entire PHP program and returns the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Stmt{},
	}

	// Skip PHP opening tag if present
	if p.curTokenIs(lexer.OPEN_TAG) || p.curTokenIs(lexer.OPEN_TAG_ECHO) {
		p.nextToken()
	}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a single statement
// This is a placeholder - will be implemented in Task 1.7
func (p *Parser) parseStatement() ast.Stmt {
	// TODO: Implement statement parsing in Task 1.7
	return nil
}

// Token management methods

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Skip comments automatically
	for p.peekToken.Type == lexer.COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

// curTokenIs checks if the current token is of the given type
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if the peek token is of the given type
func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the peek token is of the expected type and advances if so
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// expectCurrent checks if the current token is of the expected type
func (p *Parser) expectCurrent(t lexer.TokenType) bool {
	if p.curTokenIs(t) {
		return true
	}
	p.error(fmt.Sprintf("expected token %s, got %s", t, p.curToken.Type))
	return false
}

// Error handling methods

// error adds an error message to the parser's error list
func (p *Parser) error(msg string) {
	errMsg := fmt.Sprintf("[%s] Parse error: %s", p.curToken.Pos, msg)
	p.errors = append(p.errors, errMsg)
}

// peekError adds an error about an unexpected peek token
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.error(msg)
}

// Errors returns all parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// HasErrors returns true if there are any parsing errors
func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

// Error recovery methods

// synchronize attempts to recover from a parse error by skipping tokens
// until we reach a statement boundary
func (p *Parser) synchronize() {
	p.panicMode = false

	for !p.curTokenIs(lexer.EOF) {
		// If we just passed a semicolon, we're at a statement boundary
		if p.curToken.Type == lexer.SEMICOLON {
			return
		}

		// Check for statement-starting keywords
		switch p.peekToken.Type {
		case lexer.CLASS, lexer.FUNCTION, lexer.INTERFACE, lexer.TRAIT,
			lexer.NAMESPACE, lexer.USE, lexer.CONST,
			lexer.IF, lexer.WHILE, lexer.FOR, lexer.FOREACH,
			lexer.SWITCH, lexer.RETURN, lexer.BREAK, lexer.CONTINUE,
			lexer.ECHO, lexer.TRY, lexer.THROW:
			return
		}

		p.nextToken()
	}
}

// skipTo skips tokens until the specified token type is encountered
func (p *Parser) skipTo(t lexer.TokenType) {
	for !p.curTokenIs(t) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()
	}
}

// skipToStatementEnd skips to the end of the current statement (semicolon or })
func (p *Parser) skipToStatementEnd() {
	for !p.curTokenIs(lexer.SEMICOLON) && !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()
	}
}

// Helper methods for parsing

// parseIdentifier parses an identifier token
func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// currentTokenPrecedence returns the precedence of the current token
// This will be used for expression parsing (Task 1.6)
func (p *Parser) currentTokenPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// peekTokenPrecedence returns the precedence of the peek token
func (p *Parser) peekTokenPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// Operator precedence levels (for expression parsing)
const (
	_ int = iota
	LOWEST
	LOGICAL_OR       // or, ||
	LOGICAL_XOR      // xor
	LOGICAL_AND      // and, &&
	ASSIGNMENT       // =, +=, -=, etc.
	TERNARY          // ? :
	COALESCE         // ??
	BITWISE_OR       // |
	BITWISE_XOR      // ^
	BITWISE_AND      // &
	EQUALITY         // ==, ===, !=, !==
	COMPARISON       // <, >, <=, >=, <=>
	SHIFT            // <<, >>
	CONCAT           // .
	SUM              // +, -
	PRODUCT          // *, /, %
	POWER            // **
	UNARY            // !, ~, ++, --, @, cast
	POSTFIX          // [], ->, ::, ()
	NEW_CLONE        // new, clone
)

// precedences maps token types to their precedence levels
var precedences = map[lexer.TokenType]int{
	lexer.OR:                LOGICAL_OR,
	lexer.LOGICAL_OR:        LOGICAL_OR,
	lexer.XOR:               LOGICAL_XOR,
	lexer.AND:               LOGICAL_AND,
	lexer.LOGICAL_AND:       LOGICAL_AND,
	lexer.ASSIGN:            ASSIGNMENT,
	lexer.PLUS_ASSIGN:       ASSIGNMENT,
	lexer.MINUS_ASSIGN:      ASSIGNMENT,
	lexer.MUL_ASSIGN:        ASSIGNMENT,
	lexer.DIV_ASSIGN:        ASSIGNMENT,
	lexer.MOD_ASSIGN:        ASSIGNMENT,
	lexer.CONCAT_ASSIGN:     ASSIGNMENT,
	lexer.POWER_ASSIGN:      ASSIGNMENT,
	lexer.AND_ASSIGN:        ASSIGNMENT,
	lexer.OR_ASSIGN:         ASSIGNMENT,
	lexer.XOR_ASSIGN:        ASSIGNMENT,
	lexer.SL_ASSIGN:         ASSIGNMENT,
	lexer.SR_ASSIGN:         ASSIGNMENT,
	lexer.COALESCE_ASSIGN:   ASSIGNMENT,
	lexer.QUESTION:          TERNARY,
	lexer.COALESCE:          COALESCE,
	lexer.BITWISE_OR:        BITWISE_OR,
	lexer.BITWISE_XOR:       BITWISE_XOR,
	lexer.BITWISE_AND:       BITWISE_AND,
	lexer.EQ:                EQUALITY,
	lexer.IDENTICAL:         EQUALITY,
	lexer.NE:                EQUALITY,
	lexer.NOT_IDENTICAL:     EQUALITY,
	lexer.LT:                COMPARISON,
	lexer.LE:                COMPARISON,
	lexer.GT:                COMPARISON,
	lexer.GE:                COMPARISON,
	lexer.SPACESHIP:         COMPARISON,
	lexer.INSTANCEOF:        COMPARISON,
	lexer.SL:                SHIFT,
	lexer.SR:                SHIFT,
	lexer.CONCAT:            CONCAT,
	lexer.PLUS:              SUM,
	lexer.MINUS:             SUM,
	lexer.ASTERISK:          PRODUCT,
	lexer.SLASH:             PRODUCT,
	lexer.PERCENT:           PRODUCT,
	lexer.POWER:             POWER,
	lexer.LBRACKET:          POSTFIX,
	lexer.OBJECT_OPERATOR:   POSTFIX,
	lexer.NULLSAFE_OPERATOR: POSTFIX,
	lexer.PAAMAYIM_NEKUDOTAYIM: POSTFIX,
	lexer.LPAREN:            POSTFIX,
	lexer.NEW:               NEW_CLONE,
	lexer.CLONE:             NEW_CLONE,
}

// ParseFile is a convenience function to parse a PHP file
func ParseFile(filename string) (*ast.Program, error) {
	// This will read the file and parse it
	// For now, return nil as we need to implement file reading
	// TODO: Implement in Task 1.5
	return nil, fmt.Errorf("ParseFile not yet implemented")
}

// ParseString is a convenience function to parse PHP source code from a string
func ParseString(input string) (*ast.Program, []string) {
	l := lexer.New(input, "<string>")
	p := New(l)
	program := p.ParseProgram()
	return program, p.Errors()
}
