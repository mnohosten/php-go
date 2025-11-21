package parser

import (
	"fmt"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// parseFunctionDeclaration parses a function declaration
// function [&]name(params): returnType { body }
func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	funcDecl := &ast.FunctionDeclaration{
		Token: p.curToken,
	}

	// Check for reference return (&function)
	if p.peekTokenIs(lexer.BITWISE_AND) {
		p.nextToken()
		funcDecl.ByRef = true
	}

	// Expect function name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	funcDecl.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse parameter list
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	funcDecl.Parameters = p.parseFunctionParameters()

	// Parse return type hint
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		funcDecl.ReturnType = p.parseTypeHint()
	}

	// Parse function body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	funcDecl.Body = p.parseBlockStatement()

	return funcDecl
}

// parseFunctionParameters parses function parameter list
// ($param1, $param2 = default, Type $param3, ...$rest)
func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	params := []*ast.Parameter{}

	// Empty parameter list
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken() // move to first parameter

	for {
		param := p.parseParameter()
		if param != nil {
			params = append(params, param)
		}

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken() // move to next parameter
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return params
}

// parseParameter parses a single function parameter
// [Type] [&][$]name [= default]
func (p *Parser) parseParameter() *ast.Parameter {
	param := &ast.Parameter{}

	// Check for variadic (...)
	if p.curTokenIs(lexer.ELLIPSIS) {
		param.Variadic = true
		p.nextToken()
	}

	// Check for reference before anything else
	if p.curTokenIs(lexer.BITWISE_AND) {
		param.ByRef = true
		p.nextToken()
	}

	// Check for type hint (not a variable)
	if !p.curTokenIs(lexer.VARIABLE) && !p.curTokenIs(lexer.BITWISE_AND) {
		// This is a type hint
		param.Type = p.parseTypeHint()
		// At this point, parseTypeHint has parsed the type but hasn't advanced past it
		// We're still on the last token of the type

		// Don't advance yet - check what's next
		// If next is &, it's part of the parameter syntax (by-reference), not the type
		// If next is $, we're done with the type

		// Advance to see what's after the type
		if p.peekTokenIs(lexer.BITWISE_AND) {
			// Next is &, which means by-reference parameter
			p.nextToken() // move to &
			param.ByRef = true
			// Now advance to the variable
			if !p.expectPeek(lexer.VARIABLE) {
				return nil
			}
		} else if p.peekTokenIs(lexer.VARIABLE) {
			// Next is the variable name
			p.nextToken() // move to variable
		} else {
			// Unexpected token after type
			p.error(fmt.Sprintf("expected variable name or & after type, got %s instead", p.peekToken.Type))
			return nil
		}
	}

	// Parse variable name
	if !p.curTokenIs(lexer.VARIABLE) {
		p.error("expected variable name in parameter")
		return nil
	}

	param.Name = &ast.Variable{
		Token: p.curToken,
		Name:  p.curToken.Literal[1:], // Remove $
	}

	// Check for default value
	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken() // consume =
		p.nextToken() // move to value
		param.DefaultValue = p.parseExpression(LOWEST)
	}

	return param
}

// parseTypeHint parses a type hint (now uses comprehensive type parser from types.go)
// Supports: scalar types, nullable (?Type), union (A|B), intersection (A&B)
func (p *Parser) parseTypeHint() ast.Expr {
	return p.parseType()
}

// parseClassDeclaration parses a class declaration
// [abstract|final] class Name [extends Parent] [implements Interface1, Interface2] { body }
func (p *Parser) parseClassDeclaration() *ast.ClassDeclaration {
	classDecl := &ast.ClassDeclaration{
		Token:     p.curToken,
		Modifiers: []string{},
		Body:      []ast.Stmt{},
	}

	// Expect class name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	classDecl.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse extends clause
	if p.peekTokenIs(lexer.EXTENDS) {
		p.nextToken() // consume extends
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		classDecl.Extends = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}

	// Parse implements clause
	if p.peekTokenIs(lexer.IMPLEMENTS) {
		p.nextToken() // consume implements
		classDecl.Implements = p.parseInterfaceList()
	}

	// Parse class body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken() // move into body

	// Parse class members
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		member := p.parseClassMember()
		if member != nil {
			classDecl.Body = append(classDecl.Body, member)
		}
		p.nextToken()
	}

	return classDecl
}

// parseClassMember parses a class member (property, method, constant, trait use)
func (p *Parser) parseClassMember() ast.Stmt {
	// Check for use statement (traits)
	if p.curTokenIs(lexer.USE) {
		return p.parseTraitUse()
	}

	// Check for const
	if p.curTokenIs(lexer.CONST) {
		return p.parseClassConstant("public")
	}

	// Collect modifiers
	var modifiers []string
	visibility := "public" // default visibility

	for {
		switch p.curToken.Type {
		case lexer.PUBLIC:
			visibility = "public"
			p.nextToken()
		case lexer.PROTECTED:
			visibility = "protected"
			p.nextToken()
		case lexer.PRIVATE:
			visibility = "private"
			p.nextToken()
		case lexer.STATIC:
			modifiers = append(modifiers, "static")
			p.nextToken()
		case lexer.ABSTRACT:
			modifiers = append(modifiers, "abstract")
			p.nextToken()
		case lexer.FINAL:
			modifiers = append(modifiers, "final")
			p.nextToken()
		case lexer.READONLY:
			modifiers = append(modifiers, "readonly")
			p.nextToken()
		default:
			goto endModifiers
		}
	}

endModifiers:
	// Check for const after visibility modifiers
	if p.curTokenIs(lexer.CONST) {
		return p.parseClassConstant(visibility)
	}

	// Now we should have either 'function' or 'var' or a type hint or variable
	if p.curTokenIs(lexer.FUNCTION) {
		return p.parseMethodDeclaration(visibility, modifiers)
	}

	// Check for VAR keyword (old style)
	if p.curTokenIs(lexer.VAR) {
		p.nextToken() // consume 'var'
		return p.parsePropertyDeclaration("public", []string{})
	}

	// Otherwise, it's a property declaration (with or without type hint)
	return p.parsePropertyDeclaration(visibility, modifiers)
}

// parseMethodDeclaration parses a method declaration
func (p *Parser) parseMethodDeclaration(visibility string, modifiers []string) *ast.MethodDeclaration {
	method := &ast.MethodDeclaration{
		Token:      p.curToken,
		Visibility: visibility,
	}

	// Process modifiers
	for _, mod := range modifiers {
		switch mod {
		case "static":
			method.Static = true
		case "abstract":
			method.Abstract = true
		case "final":
			method.Final = true
		}
	}

	// Check for reference return
	if p.peekTokenIs(lexer.BITWISE_AND) {
		p.nextToken()
		method.ByRef = true
	}

	// Parse method name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	method.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse parameters
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	method.Parameters = p.parseFunctionParameters()

	// Parse return type
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		method.ReturnType = p.parseTypeHint()
	}

	// Parse body (abstract methods have no body)
	if method.Abstract {
		// Expect semicolon
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
		return method
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	method.Body = p.parseBlockStatement()

	return method
}

// parsePropertyDeclaration parses a property declaration
func (p *Parser) parsePropertyDeclaration(visibility string, modifiers []string) *ast.PropertyDeclaration {
	prop := &ast.PropertyDeclaration{
		Token:      p.curToken,
		Visibility: visibility,
		Properties: []*ast.PropertyItem{},
	}

	// Process modifiers
	for _, mod := range modifiers {
		switch mod {
		case "static":
			prop.Static = true
		case "readonly":
			prop.Readonly = true
		}
	}

	// Check if we have a type hint (not a variable)
	if !p.curTokenIs(lexer.VARIABLE) {
		prop.Type = p.parseTypeHint()
		if !p.expectPeek(lexer.VARIABLE) && !p.curTokenIs(lexer.VARIABLE) {
			return nil
		}
	}

	// Parse first property
	for {
		propItem := &ast.PropertyItem{}

		if !p.curTokenIs(lexer.VARIABLE) {
			if !p.expectPeek(lexer.VARIABLE) {
				return nil
			}
		}

		propItem.Name = &ast.Variable{
			Token: p.curToken,
			Name:  p.curToken.Literal[1:], // Remove $
		}

		// Check for default value
		if p.peekTokenIs(lexer.ASSIGN) {
			p.nextToken() // consume =
			p.nextToken() // move to value
			propItem.DefaultValue = p.parseExpression(LOWEST)
		}

		prop.Properties = append(prop.Properties, propItem)

		// Check for multiple properties on same line
		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken() // move to next variable
	}

	// Expect semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return prop
}

// parseInterfaceDeclaration parses an interface declaration
func (p *Parser) parseInterfaceDeclaration() *ast.InterfaceDeclaration {
	interfaceDecl := &ast.InterfaceDeclaration{
		Token: p.curToken,
		Body:  []*ast.MethodSignature{},
	}

	// Expect interface name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	interfaceDecl.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse extends clause (interfaces can extend multiple interfaces)
	if p.peekTokenIs(lexer.EXTENDS) {
		p.nextToken() // consume extends
		interfaceDecl.Extends = p.parseInterfaceList()
	}

	// Parse interface body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken() // move into body

	// Parse method signatures
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Skip visibility modifiers (public is implicit in interfaces)
		if p.curTokenIs(lexer.PUBLIC) {
			p.nextToken()
		}

		if p.curTokenIs(lexer.FUNCTION) {
			signature := p.parseMethodSignature()
			if signature != nil {
				interfaceDecl.Body = append(interfaceDecl.Body, signature)
			}
		}
		p.nextToken()
	}

	return interfaceDecl
}

// parseMethodSignature parses a method signature (no body)
func (p *Parser) parseMethodSignature() *ast.MethodSignature {
	signature := &ast.MethodSignature{
		Token: p.curToken,
	}

	// Check for reference return
	if p.peekTokenIs(lexer.BITWISE_AND) {
		p.nextToken()
		signature.ByRef = true
	}

	// Parse method name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	signature.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse parameters
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	signature.Parameters = p.parseFunctionParameters()

	// Parse return type
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		signature.ReturnType = p.parseTypeHint()
	}

	// Expect semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return signature
}

// parseTraitDeclaration parses a trait declaration
func (p *Parser) parseTraitDeclaration() *ast.TraitDeclaration {
	traitDecl := &ast.TraitDeclaration{
		Token: p.curToken,
		Body:  []ast.Stmt{},
	}

	// Expect trait name
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	traitDecl.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Parse trait body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken() // move into body

	// Parse trait members (properties and methods)
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		member := p.parseTraitMember()
		if member != nil {
			traitDecl.Body = append(traitDecl.Body, member)
		}
		p.nextToken()
	}

	return traitDecl
}

// parseTraitMember parses a trait member (property or method)
func (p *Parser) parseTraitMember() ast.Stmt {
	// Collect modifiers
	var modifiers []string
	visibility := "public"

	for {
		switch p.curToken.Type {
		case lexer.PUBLIC:
			visibility = "public"
			p.nextToken()
		case lexer.PROTECTED:
			visibility = "protected"
			p.nextToken()
		case lexer.PRIVATE:
			visibility = "private"
			p.nextToken()
		case lexer.STATIC:
			modifiers = append(modifiers, "static")
			p.nextToken()
		case lexer.ABSTRACT:
			modifiers = append(modifiers, "abstract")
			p.nextToken()
		case lexer.FINAL:
			modifiers = append(modifiers, "final")
			p.nextToken()
		default:
			goto endModifiers
		}
	}

endModifiers:
	if p.curTokenIs(lexer.FUNCTION) {
		return p.parseMethodDeclaration(visibility, modifiers)
	}

	// Property declaration
	return p.parsePropertyDeclaration(visibility, modifiers)
}

// parseTraitUse parses trait usage in a class
func (p *Parser) parseTraitUse() *ast.TraitUse {
	traitUse := &ast.TraitUse{
		Token:       p.curToken,
		Traits:      []*ast.Identifier{},
		Adaptations: []ast.TraitAdaptation{},
	}

	// Parse trait list
	p.nextToken() // move to first trait name

	for {
		if !p.curTokenIs(lexer.IDENT) {
			p.error("expected trait name")
			return nil
		}

		traitUse.Traits = append(traitUse.Traits, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken() // move to next trait
	}

	// Check for trait adaptations block
	if p.peekTokenIs(lexer.LBRACE) {
		p.nextToken() // consume {
		p.nextToken() // move into body

		// Parse adaptations (insteadof, as)
		for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
			// Simplified trait adaptation parsing
			// Full implementation would handle complex cases
			// For now, just skip to semicolon
			p.skipToStatementEnd()
			p.nextToken()
		}
	} else {
		// Simple trait use without adaptations - expect semicolon
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	return traitUse
}

// parseClassConstant parses a class constant declaration
func (p *Parser) parseClassConstant(visibility string) *ast.ClassConstantDeclaration {
	constDecl := &ast.ClassConstantDeclaration{
		Token:      p.curToken,
		Visibility: visibility,
		Constants:  []*ast.ConstantItem{},
	}

	p.nextToken() // move to first constant name

	for {
		if !p.curTokenIs(lexer.IDENT) {
			p.error("expected constant name")
			return nil
		}

		constItem := &ast.ConstantItem{
			Name: &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			},
		}

		// Expect = value
		if !p.expectPeek(lexer.ASSIGN) {
			return nil
		}

		p.nextToken() // move to value
		constItem.Value = p.parseExpression(LOWEST)

		constDecl.Constants = append(constDecl.Constants, constItem)

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken() // move to next constant
	}

	// Expect semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return constDecl
}

// parseInterfaceList parses a comma-separated list of interface names
func (p *Parser) parseInterfaceList() []*ast.Identifier {
	interfaces := []*ast.Identifier{}

	p.nextToken() // move to first interface

	for {
		if !p.curTokenIs(lexer.IDENT) {
			p.error("expected interface name")
			return nil
		}

		interfaces = append(interfaces, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken() // move to next interface
	}

	return interfaces
}
