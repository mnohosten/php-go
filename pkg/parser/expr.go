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
	p.prefixParseFns[lexer.ARRAY] = p.parseArrayConstructor
	p.prefixParseFns[lexer.NEW] = p.parseNewExpression
	p.prefixParseFns[lexer.MATCH] = p.parseMatchExpression
	p.prefixParseFns[lexer.FUNCTION] = p.parseClosureExpression
	p.prefixParseFns[lexer.FN] = p.parseArrowFunctionExpression
	p.prefixParseFns[lexer.STATIC] = p.parseStaticClosureOrProperty
	p.prefixParseFns[lexer.EXIT] = p.parseExitExpression
	p.prefixParseFns[lexer.ISSET] = p.parseIssetExpression
	p.prefixParseFns[lexer.EMPTY] = p.parseEmptyExpression
	p.prefixParseFns[lexer.LIST] = p.parseListExpression
	p.prefixParseFns[lexer.INCLUDE] = p.parseIncludeExpression
	p.prefixParseFns[lexer.INCLUDE_ONCE] = p.parseIncludeExpression
	p.prefixParseFns[lexer.REQUIRE] = p.parseIncludeExpression
	p.prefixParseFns[lexer.REQUIRE_ONCE] = p.parseIncludeExpression

	// Magic constants
	p.prefixParseFns[lexer.LINE_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.FILE_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.DIR_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.FUNCTION_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.CLASS_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.TRAIT_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.METHOD_CONST] = p.parseMagicConstant
	p.prefixParseFns[lexer.NAMESPACE_CONST] = p.parseMagicConstant

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

	// Postfix ++ and --
	p.infixParseFns[lexer.INC] = p.parsePostfixExpression
	p.infixParseFns[lexer.DEC] = p.parsePostfixExpression
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
	strValue := p.curToken.Literal
	token := p.curToken

	// Check if string contains interpolation (simple $variable detection)
	if !p.hasInterpolation(strValue) {
		// No interpolation - return simple string literal
		return &ast.StringLiteral{
			Token: token,
			Value: strValue,
		}
	}

	// String has interpolation - parse it into parts
	return p.parseInterpolatedString(token, strValue)
}

// hasInterpolation checks if a string contains variable interpolation
func (p *Parser) hasInterpolation(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == '\\' && i+1 < len(str) {
			// Skip escaped characters
			i++
			continue
		}
		if str[i] == '$' && i+1 < len(str) {
			// Check if next char is valid variable start (letter or underscore)
			next := str[i+1]
			if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_' {
				return true
			}
		}
	}
	return false
}

// parseInterpolatedString parses a string with interpolated variables
// Example: "Hello $name" becomes ["Hello ", $name]
func (p *Parser) parseInterpolatedString(token lexer.Token, str string) ast.Expr {
	parts := []ast.Expr{}
	var currentPart []byte
	i := 0

	for i < len(str) {
		if str[i] == '\\' && i+1 < len(str) {
			// Handle escape sequences
			i++
			switch str[i] {
			case 'n':
				currentPart = append(currentPart, '\n')
			case 't':
				currentPart = append(currentPart, '\t')
			case 'r':
				currentPart = append(currentPart, '\r')
			case '\\':
				currentPart = append(currentPart, '\\')
			case '$':
				currentPart = append(currentPart, '$')
			default:
				currentPart = append(currentPart, str[i])
			}
			i++
			continue
		}

		if str[i] == '$' && i+1 < len(str) {
			next := str[i+1]
			if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_' {
				// Found variable interpolation
				// Add current string part if non-empty
				if len(currentPart) > 0 {
					parts = append(parts, &ast.StringLiteral{
						Token: token,
						Value: string(currentPart),
					})
					currentPart = nil
				}

				// Extract variable name
				i++ // skip $
				varStart := i
				for i < len(str) && (isLetter(str[i]) || isDigit(str[i]) || str[i] == '_') {
					i++
				}
				varName := str[varStart:i]

				// Add variable to parts
				parts = append(parts, &ast.Variable{
					Token: token,
					Name:  varName,
				})
				continue
			}
		}

		// Regular character
		currentPart = append(currentPart, str[i])
		i++
	}

	// Add final string part if non-empty
	if len(currentPart) > 0 {
		parts = append(parts, &ast.StringLiteral{
			Token: token,
			Value: string(currentPart),
		})
	}

	// If only one part, return it directly
	if len(parts) == 1 {
		return parts[0]
	}

	return &ast.InterpolatedStringExpression{
		Token: token,
		Parts: parts,
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

func (p *Parser) parseMagicConstant() ast.Expr {
	return &ast.MagicConstant{
		Token: p.curToken,
		Kind:  p.curToken.Type,
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
	   p.peekTokenIs(lexer.BOOL) || p.peekTokenIs(lexer.FLOAT_TYPE) ||
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

		// Allow trailing comma - check if next token is closing bracket
		if p.peekTokenIs(lexer.RBRACKET) {
			break
		}

		p.nextToken() // move to next element
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

// parseArrayConstructor parses array() constructor syntax (legacy PHP array syntax)
// Example: array(1, 2, 3) or array('key' => 'value')
func (p *Parser) parseArrayConstructor() ast.Expr {
	array := &ast.ArrayExpression{
		Token:    p.curToken, // 'array' token
		Elements: []ast.ArrayElement{},
	}

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Check for empty array: array()
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return array
	}

	p.nextToken() // move to first element

	// Parse first element
	array.Elements = append(array.Elements, p.parseArrayElement())

	// Parse remaining elements
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma

		// Allow trailing comma - check if next token is closing parenthesis
		if p.peekTokenIs(lexer.RPAREN) {
			break
		}

		p.nextToken() // move to next element
		array.Elements = append(array.Elements, p.parseArrayElement())
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return array
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

func (p *Parser) parseExitExpression() ast.Expr {
	// exit and die can be used as:
	// 1. exit; or die;  (no parentheses, no argument)
	// 2. exit() or die() (empty parentheses)
	// 3. exit("message") or die("message") (with argument)
	// 4. exit(0) or die(0) (with exit code)

	// We'll parse it as a CallExpression with "exit" as the function name
	token := p.curToken

	// Create identifier for "exit"
	funcIdent := &ast.Identifier{
		Token: token,
		Value: "exit", // Normalize both exit and die to "exit"
	}

	expression := &ast.CallExpression{
		Token:    token,
		Function: funcIdent,
	}

	// Check for optional parentheses
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken() // move to (
		expression.Arguments = p.parseCallArguments()
	} else {
		// No parentheses, no arguments
		expression.Arguments = []*ast.Argument{}
	}

	return expression
}

func (p *Parser) parseIssetExpression() ast.Expr {
	// isset() checks if one or more variables are set and not null
	// Examples:
	// isset($var)
	// isset($a, $b, $c)
	// isset($arr['key'])

	token := p.curToken // The ISSET token

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	expression := &ast.IssetExpression{
		Token:     token,
		Variables: []ast.Expr{},
	}

	// Check for empty isset()
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // move to )
		p.error("isset() expects at least one argument")
		return expression
	}

	p.nextToken() // move past (

	// Parse first variable
	variable := p.parseExpression(LOWEST)
	if variable == nil {
		return nil
	}
	expression.Variables = append(expression.Variables, variable)

	// Parse additional variables (comma-separated)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // move to ,
		p.nextToken() // move past ,

		variable := p.parseExpression(LOWEST)
		if variable == nil {
			return nil
		}
		expression.Variables = append(expression.Variables, variable)
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseEmptyExpression() ast.Expr {
	// empty() checks if a variable is empty (falsy or not set)
	// Examples:
	// empty($var)
	// empty($arr['key'])

	token := p.curToken // The EMPTY token

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Check for empty empty()
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // move to )
		p.error("empty() expects exactly one argument")
		return nil
	}

	p.nextToken() // move past (

	// Parse the variable
	variable := p.parseExpression(LOWEST)
	if variable == nil {
		return nil
	}

	expression := &ast.EmptyExpression{
		Token:    token,
		Variable: variable,
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseListExpression() ast.Expr {
	// list() is used for array destructuring assignment
	// Examples:
	// list($a, $b) = array(1, 2)                  // Basic list
	// list($a, , $c) = array(1, 2, 3)             // Skip middle element
	// list($x[0], $y->prop) = array(1, 2)         // Assign to array/object
	// list("a" => $a, "b" => $b) = ["a"=>1,"b"=>2] // Keyed list (PHP 7.1+)

	token := p.curToken // The LIST token

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	expression := &ast.ListExpression{
		Token:    token,
		Elements: []*ast.ListElement{},
	}

	// Check for empty list()
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // move to )
		p.error("list() cannot be empty")
		return expression
	}

	p.nextToken() // move past (

	// Parse list elements
	for !p.curTokenIs(lexer.RPAREN) {
		var element *ast.ListElement

		// Check for skipped element at start or after comma: list(, $a) or list($a, , $c)
		if p.curTokenIs(lexer.COMMA) {
			// Skipped element
			element = &ast.ListElement{
				Key:   nil,
				Value: nil,
			}
			expression.Elements = append(expression.Elements, element)
			p.nextToken() // move past ,
			continue
		}

		// Parse a normal element
		// Check for keyed list (PHP 7.1+): list("key" => $var)
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}

		// Check if this is a keyed element
		if p.peekTokenIs(lexer.DOUBLE_ARROW) {
			p.nextToken() // move to =>
			p.nextToken() // move past =>

			// Parse the value
			value := p.parseExpression(LOWEST)
			if value == nil {
				return nil
			}

			element = &ast.ListElement{
				Key:   expr,
				Value: value,
			}
		} else {
			// Non-keyed element
			element = &ast.ListElement{
				Key:   nil,
				Value: expr,
			}
		}

		expression.Elements = append(expression.Elements, element)

		// Check for comma
		if !p.peekTokenIs(lexer.COMMA) {
			// No more elements
			break
		}

		p.nextToken() // move to ,

		// Check if next is closing paren (trailing comma)
		if p.peekTokenIs(lexer.RPAREN) {
			break
		}

		p.nextToken() // move past ,
		// Continue loop - if next is comma, it's a skipped element
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIncludeExpression() ast.Expr {
	// include, include_once, require, require_once
	// These can be used with or without parentheses:
	// require 'file.php'
	// require('file.php')
	// $result = include 'file.php'

	token := p.curToken // The INCLUDE/INCLUDE_ONCE/REQUIRE/REQUIRE_ONCE token

	// Determine the type
	includeType := ""
	switch token.Type {
	case lexer.INCLUDE:
		includeType = "include"
	case lexer.INCLUDE_ONCE:
		includeType = "include_once"
	case lexer.REQUIRE:
		includeType = "require"
	case lexer.REQUIRE_ONCE:
		includeType = "require_once"
	}

	p.nextToken() // move to the path expression

	// Parse the path expression
	// Note: We parse with LOWEST precedence to handle expressions like:
	// require DIR . '/file.php'
	path := p.parseExpression(LOWEST)
	if path == nil {
		return nil
	}

	return &ast.IncludeExpression{
		Token: token,
		Path:  path,
		Type:  includeType,
	}
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
	token := p.curToken

	// Check for first-class callable syntax: foo(...)
	// In PHP 8.1+, if the only thing in parentheses is ..., it creates a callable
	if p.peekTokenIs(lexer.ELLIPSIS) {
		p.nextToken() // move to ...
		if p.peekTokenIs(lexer.RPAREN) {
			p.nextToken() // consume )
			return &ast.FirstClassCallableExpression{
				Token:    token,
				Callable: left,
			}
		}
		// Otherwise it's unpacking - go back and parse normally
		// This is a bit hacky but works for our use case
	}

	return &ast.CallExpression{
		Token:     token,
		Function:  left,
		Arguments: p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []*ast.Argument {
	args := []*ast.Argument{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	// Parse first argument (may be named or positional)
	arg := p.parseArgument()
	if arg != nil {
		args = append(args, arg)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next argument
		arg := p.parseArgument()
		if arg != nil {
			args = append(args, arg)
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return args
}

// parseArgument parses a function call argument (positional or named)
// Named arguments have the syntax: name: value
// Unpacked arguments have the syntax: ...$expr
// Positional arguments are just expressions
func (p *Parser) parseArgument() *ast.Argument {
	arg := &ast.Argument{
		Token: p.curToken,
	}

	// Check for unpacking operator (...)
	if p.curTokenIs(lexer.ELLIPSIS) {
		arg.Unpack = true
		p.nextToken() // consume ...
		arg.Value = p.parseExpression(LOWEST)
		return arg
	}

	// Check if this is a named argument: identifier followed by colon
	if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.COLON) {
		// Named argument
		arg.Name = p.curToken.Literal
		p.nextToken() // consume identifier
		arg.Token = p.curToken // update token to colon
		p.nextToken() // consume colon, move to value expression
		arg.Value = p.parseExpression(LOWEST)
	} else {
		// Positional argument
		arg.Value = p.parseExpression(LOWEST)
	}

	return arg
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

func (p *Parser) parsePostfixExpression(left ast.Expr) ast.Expr {
	// Postfix ++ and -- are unary operators applied after the operand
	return &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal + "(postfix)",
		Right:    left,
	}
}

// parseClosureExpression parses a closure (anonymous function)
// Example: function($x, $y) use ($z) { return $x + $y + $z; }
func (p *Parser) parseClosureExpression() ast.Expr {
	expr := &ast.ClosureExpression{
		Token: p.curToken, // FUNCTION token
	}

	// Check for reference return
	if p.peekTokenIs(lexer.BITWISE_AND) {
		p.nextToken()
		expr.ByRef = true
	}

	// Expect '('
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse parameters
	expr.Parameters = p.parseFunctionParameters()

	// Check for 'use' clause
	if p.peekTokenIs(lexer.USE) {
		p.nextToken() // consume USE

		if !p.expectPeek(lexer.LPAREN) {
			return nil
		}
		p.nextToken() // move to first variable

		expr.Use = p.parseUseClause()

		if !p.curTokenIs(lexer.RPAREN) {
			return nil
		}
	}

	// Check for return type hint
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume COLON
		p.nextToken() // move to type
		expr.ReturnType = p.parseType()
	}

	// Expect body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

// parseUseClause parses the use clause of a closure
// Example: ($x, &$y, $z)
func (p *Parser) parseUseClause() []*ast.UseClause {
	useClauses := []*ast.UseClause{}

	// Handle empty use clause
	if p.curTokenIs(lexer.RPAREN) {
		return useClauses
	}

	for {
		useClause := &ast.UseClause{}

		// Check for reference capture
		if p.curTokenIs(lexer.BITWISE_AND) {
			useClause.ByRef = true
			p.nextToken()
		}

		// Expect variable
		if !p.curTokenIs(lexer.VARIABLE) {
			p.errors = append(p.errors, "expected variable in use clause, got "+p.curToken.Literal)
			return nil
		}

		useClause.Variable = &ast.Variable{
			Token: p.curToken,
			Name:  p.curToken.Literal,
		}

		useClauses = append(useClauses, useClause)

		// Check for more variables
		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume COMMA
		p.nextToken() // move to next variable
	}

	// Expect closing )
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return useClauses
}

// parseArrowFunctionExpression parses an arrow function (PHP 7.4+)
// Example: fn($x): int => $x * 2
func (p *Parser) parseArrowFunctionExpression() ast.Expr {
	expr := &ast.ArrowFunctionExpression{
		Token: p.curToken, // FN token
	}

	// Check for reference return
	if p.peekTokenIs(lexer.BITWISE_AND) {
		p.nextToken()
		expr.ByRef = true
	}

	// Expect '('
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse parameters
	expr.Parameters = p.parseFunctionParameters()

	// Check for return type hint
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume COLON
		p.nextToken() // move to type
		expr.ReturnType = p.parseType()
	}

	// Expect '=>'
	if !p.expectPeek(lexer.DOUBLE_ARROW) {
		return nil
	}

	// Parse body expression
	p.nextToken()
	expr.Body = p.parseExpression(LOWEST)

	return expr
}

// parseStaticClosureOrProperty handles 'static' keyword which can be:
// - static function() { ... }  (static closure)
// - static fn() => ...  (static arrow function)
// - static::$property  (static property access - handled elsewhere)
func (p *Parser) parseStaticClosureOrProperty() ast.Expr {
	staticToken := p.curToken

	// Look ahead to determine what follows
	if p.peekTokenIs(lexer.FUNCTION) {
		p.nextToken() // consume FUNCTION

		expr := p.parseClosureExpression()
		if closure, ok := expr.(*ast.ClosureExpression); ok {
			closure.Static = true
		}
		return expr
	}

	if p.peekTokenIs(lexer.FN) {
		p.nextToken() // consume FN

		expr := p.parseArrowFunctionExpression()
		if arrowFn, ok := expr.(*ast.ArrowFunctionExpression); ok {
			arrowFn.Static = true
		}
		return expr
	}

	// If not followed by function or fn, it's likely static::$property
	// Treat 'static' as an identifier for static property/method access
	return &ast.Identifier{
		Token: staticToken,
		Value: staticToken.Literal,
	}
}

// Helper functions for string interpolation

// isLetter checks if a byte is a letter (a-z, A-Z)
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit checks if a byte is a digit (0-9)
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
