package parser

import (
	"fmt"
	"strconv"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// Pratt parsing function types
type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

// Register prefix and infix parsing functions
func (p *Parser) registerExpressionParsers() {
	// Prefix parsers (expressions that start with these tokens)
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.prefixParseFns[lexer.IDENT] = p.parseIdentifier
	p.prefixParseFns[lexer.VARIABLE] = p.parseVariable
	p.prefixParseFns[lexer.INTEGER] = p.parseIntegerLiteral
	p.prefixParseFns[lexer.FLOAT] = p.parseFloatLiteral
	p.prefixParseFns[lexer.STRING] = p.parseStringLiteral
	p.prefixParseFns[lexer.HEREDOC] = p.parseStringLiteral
	p.prefixParseFns[lexer.NOWDOC] = p.parseStringLiteral
	p.prefixParseFns[lexer.TRUE] = p.parseBooleanLiteral
	p.prefixParseFns[lexer.FALSE] = p.parseBooleanLiteral
	p.prefixParseFns[lexer.NULL] = p.parseNullLiteral
	p.prefixParseFns[lexer.LOGICAL_NOT] = p.parsePrefixExpression
	p.prefixParseFns[lexer.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[lexer.PLUS] = p.parsePrefixExpression
	p.prefixParseFns[lexer.BITWISE_NOT] = p.parsePrefixExpression
	p.prefixParseFns[lexer.INC] = p.parsePrefixExpression
	p.prefixParseFns[lexer.DEC] = p.parsePrefixExpression
	p.prefixParseFns[lexer.AT] = p.parsePrefixExpression
	p.prefixParseFns[lexer.LPAREN] = p.parseGroupedOrCastExpression
	p.prefixParseFns[lexer.LBRACKET] = p.parseArrayExpression
	p.prefixParseFns[lexer.NEW] = p.parseNewExpression

	// Infix parsers (operators that appear between expressions)
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)

	// Arithmetic operators
	p.infixParseFns[lexer.PLUS] = p.parseInfixExpression
	p.infixParseFns[lexer.MINUS] = p.parseInfixExpression
	p.infixParseFns[lexer.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[lexer.SLASH] = p.parseInfixExpression
	p.infixParseFns[lexer.PERCENT] = p.parseInfixExpression
	p.infixParseFns[lexer.POWER] = p.parseInfixExpression

	// Comparison operators
	p.infixParseFns[lexer.EQ] = p.parseInfixExpression
	p.infixParseFns[lexer.IDENTICAL] = p.parseInfixExpression
	p.infixParseFns[lexer.NE] = p.parseInfixExpression
	p.infixParseFns[lexer.NOT_IDENTICAL] = p.parseInfixExpression
	p.infixParseFns[lexer.LT] = p.parseInfixExpression
	p.infixParseFns[lexer.LE] = p.parseInfixExpression
	p.infixParseFns[lexer.GT] = p.parseInfixExpression
	p.infixParseFns[lexer.GE] = p.parseInfixExpression
	p.infixParseFns[lexer.SPACESHIP] = p.parseInfixExpression

	// Logical operators
	p.infixParseFns[lexer.LOGICAL_AND] = p.parseInfixExpression
	p.infixParseFns[lexer.LOGICAL_OR] = p.parseInfixExpression
	p.infixParseFns[lexer.AND] = p.parseInfixExpression
	p.infixParseFns[lexer.OR] = p.parseInfixExpression
	p.infixParseFns[lexer.XOR] = p.parseInfixExpression

	// Bitwise operators
	p.infixParseFns[lexer.BITWISE_AND] = p.parseInfixExpression
	p.infixParseFns[lexer.BITWISE_OR] = p.parseInfixExpression
	p.infixParseFns[lexer.BITWISE_XOR] = p.parseInfixExpression
	p.infixParseFns[lexer.SL] = p.parseInfixExpression
	p.infixParseFns[lexer.SR] = p.parseInfixExpression

	// String concatenation
	p.infixParseFns[lexer.CONCAT] = p.parseInfixExpression

	// Assignment operators
	p.infixParseFns[lexer.ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.PLUS_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.MINUS_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.MUL_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.DIV_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.MOD_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.CONCAT_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.POWER_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.AND_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.OR_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.XOR_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.SL_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.SR_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[lexer.COALESCE_ASSIGN] = p.parseAssignmentExpression

	// Ternary and coalescing
	p.infixParseFns[lexer.QUESTION] = p.parseTernaryExpression
	p.infixParseFns[lexer.COALESCE] = p.parseInfixExpression

	// Postfix operators (array access, property access, method calls)
	p.infixParseFns[lexer.LBRACKET] = p.parseIndexExpression
	p.infixParseFns[lexer.OBJECT_OPERATOR] = p.parsePropertyOrMethodCall
	p.infixParseFns[lexer.NULLSAFE_OPERATOR] = p.parseNullsafePropertyOrMethodCall
	p.infixParseFns[lexer.PAAMAYIM_NEKUDOTAYIM] = p.parseStaticAccessOrCall
	p.infixParseFns[lexer.LPAREN] = p.parseCallExpression

	// instanceof
	p.infixParseFns[lexer.INSTANCEOF] = p.parseInstanceofExpression
}

// parseExpression is the main entry point for parsing expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.error(fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}

	leftExp := prefix()

	// Pratt parsing: continue parsing while the next operator has higher precedence
	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekTokenPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix parsing functions

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseVariable() ast.Expr {
	// Remove the $ prefix for the name
	name := p.curToken.Literal
	if len(name) > 0 && name[0] == '$' {
		name = name[1:]
	}

	return &ast.Variable{
		Token: p.curToken,
		Name:  name,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expr {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.error(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expr {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.error(fmt.Sprintf("could not parse %q as float", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expr {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(lexer.TRUE),
	}
}

func (p *Parser) parseNullLiteral() ast.Expr {
	return &ast.NullLiteral{
		Token: p.curToken,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expr {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(UNARY)

	return expression
}

func (p *Parser) parseGroupedOrCastExpression() ast.Expr {
	// Look ahead to determine if this is a cast or grouped expression
	// Cast: (int), (string), (bool), (float), (array), (object)
	// Grouped: any other expression in parentheses

	if p.peekTokenIs(lexer.INT) || p.peekTokenIs(lexer.STRING_TYPE) ||
	   p.peekTokenIs(lexer.BOOL) || p.peekTokenIs(lexer.FLOAT) ||
	   p.peekTokenIs(lexer.ARRAY) || p.peekTokenIs(lexer.OBJECT) {
		return p.parseCastExpression()
	}

	// Grouped expression
	token := p.curToken
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return &ast.GroupedExpression{
		Token: token,
		Expr:  exp,
	}
}

func (p *Parser) parseCastExpression() ast.Expr {
	token := p.curToken

	p.nextToken() // move to type token
	typeToken := p.curToken

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	p.nextToken() // move past )

	return &ast.CastExpression{
		Token: token,
		Type:  typeToken.Literal,
		Expr:  p.parseExpression(UNARY),
	}
}

func (p *Parser) parseArrayExpression() ast.Expr {
	array := &ast.ArrayExpression{
		Token:    p.curToken,
		Elements: []ast.ArrayElement{},
	}

	if p.peekTokenIs(lexer.RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()

	// Parse first element
	array.Elements = append(array.Elements, p.parseArrayElement())

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next element

		// Allow trailing comma
		if p.curTokenIs(lexer.RBRACKET) {
			break
		}

		array.Elements = append(array.Elements, p.parseArrayElement())
	}

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return array
}

func (p *Parser) parseArrayElement() ast.ArrayElement {
	// Parse first expression
	expr := p.parseExpression(LOWEST)

	// Check for => (associative array)
	if p.peekTokenIs(lexer.DOUBLE_ARROW) {
		p.nextToken() // consume =>
		p.nextToken() // move to value

		value := p.parseExpression(LOWEST)
		return ast.ArrayElement{
			Key:   expr,
			Value: value,
		}
	}

	// Non-associative element
	return ast.ArrayElement{
		Key:   nil,
		Value: expr,
	}
}

func (p *Parser) parseNewExpression() ast.Expr {
	expression := &ast.NewExpression{
		Token: p.curToken,
	}

	p.nextToken()

	// Parse class name or expression
	expression.Class = p.parseExpression(NEW_CLONE)

	// Optional arguments
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken() // move to (
		expression.Arguments = p.parseCallArguments()
	}

	return expression
}

// Infix parsing functions

func (p *Parser) parseInfixExpression(left ast.Expr) ast.Expr {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.currentTokenPrecedence()

	// Power operator is right-associative
	if p.curTokenIs(lexer.POWER) {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left ast.Expr) ast.Expr {
	expression := &ast.AssignmentExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expression.Right = p.parseExpression(ASSIGNMENT - 1) // Right-associative

	return expression
}

func (p *Parser) parseTernaryExpression(left ast.Expr) ast.Expr {
	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: left,
	}

	p.nextToken()

	// Check for short ternary (?:)
	if p.curTokenIs(lexer.COLON) {
		expression.Consequence = nil
		p.nextToken()
		expression.Alternative = p.parseExpression(TERNARY)
		return expression
	}

	// Full ternary
	expression.Consequence = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	expression.Alternative = p.parseExpression(TERNARY)

	return expression
}

func (p *Parser) parseIndexExpression(left ast.Expr) ast.Expr {
	expression := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return expression
}

func (p *Parser) parsePropertyOrMethodCall(left ast.Expr) ast.Expr {
	token := p.curToken
	p.nextToken()

	// Parse property name (can be identifier or dynamic)
	property := p.parseExpression(POSTFIX)

	// Check if this is a method call
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken() // move to (

		return &ast.MethodCallExpression{
			Token:     token,
			Object:    left,
			Method:    property,
			Arguments: p.parseCallArguments(),
		}
	}

	// Property access
	return &ast.PropertyExpression{
		Token:    token,
		Object:   left,
		Property: property,
	}
}

func (p *Parser) parseNullsafePropertyOrMethodCall(left ast.Expr) ast.Expr {
	token := p.curToken
	p.nextToken()

	property := p.parseExpression(POSTFIX)

	// Check if this is a method call
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()

		// For nullsafe method calls, we wrap in NullsafePropertyExpression
		// The actual method call handling will be in the VM
		return &ast.MethodCallExpression{
			Token:     token,
			Object:    left,
			Method:    property,
			Arguments: p.parseCallArguments(),
		}
	}

	// Nullsafe property access
	return &ast.NullsafePropertyExpression{
		Token:    token,
		Object:   left,
		Property: property,
	}
}

func (p *Parser) parseStaticAccessOrCall(left ast.Expr) ast.Expr {
	token := p.curToken
	p.nextToken()

	// Parse member (method, property, or constant)
	member := p.parseExpression(POSTFIX)

	// Check if this is a method call
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()

		return &ast.StaticCallExpression{
			Token:     token,
			Class:     left,
			Method:    member,
			Arguments: p.parseCallArguments(),
		}
	}

	// Static property or constant access
	return &ast.StaticPropertyExpression{
		Token:    token,
		Class:    left,
		Property: member,
	}
}

func (p *Parser) parseCallExpression(left ast.Expr) ast.Expr {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  left,
		Arguments: p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []ast.Expr {
	args := []ast.Expr{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next argument
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseInstanceofExpression(left ast.Expr) ast.Expr {
	expression := &ast.InstanceofExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Right = p.parseExpression(COMPARISON)

	return expression
}
