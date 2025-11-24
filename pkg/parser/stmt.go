package parser

import (
	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// parseEchoStatement parses echo statement
func (p *Parser) parseEchoStatement() *ast.EchoStatement {
	stmt := &ast.EchoStatement{
		Token:       p.curToken,
		Expressions: []ast.Expr{},
	}

	p.nextToken()

	// Parse first expression
	stmt.Expressions = append(stmt.Expressions, p.parseExpression(LOWEST))

	// Parse remaining expressions separated by commas
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to expression
		stmt.Expressions = append(stmt.Expressions, p.parseExpression(LOWEST))
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	// Check if there's a return value
	if !p.peekTokenIs(lexer.SEMICOLON) && !p.peekTokenIs(lexer.EOF) {
		p.nextToken()
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBreakStatement parses break statement
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{
		Token: p.curToken,
	}

	// Optional depth
	if p.peekTokenIs(lexer.INTEGER) {
		p.nextToken()
		stmt.Depth = p.parseExpression(LOWEST)
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseContinueStatement parses continue statement
func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{
		Token: p.curToken,
	}

	// Optional depth
	if p.peekTokenIs(lexer.INTEGER) {
		p.nextToken()
		stmt.Depth = p.parseExpression(LOWEST)
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseGlobalStatement parses global variable declaration
// Syntax: global $var1, $var2, ...;
func (p *Parser) parseGlobalStatement() *ast.GlobalStatement {
	stmt := &ast.GlobalStatement{
		Token:     p.curToken,
		Variables: []*ast.Identifier{},
	}

	// Parse variable list
	for {
		// Expect a variable
		if !p.expectPeek(lexer.VARIABLE) {
			return nil
		}

		// Create identifier for the variable
		variable := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		stmt.Variables = append(stmt.Variables, variable)

		// Check for comma (more variables)
		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseUnsetStatement parses unset() statement
func (p *Parser) parseUnsetStatement() *ast.UnsetStatement {
	stmt := &ast.UnsetStatement{
		Token:     p.curToken,
		Variables: []ast.Expr{},
	}

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse variable list
	for {
		p.nextToken() // move to next expression

		// Parse the expression (variable, array access, object property, etc.)
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}
		stmt.Variables = append(stmt.Variables, expr)

		// Check for comma (more variables)
		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
	}

	// Expect closing parenthesis
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseIfStatement parses if/elseif/else statement
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{
		Token:   p.curToken,
		ElseIfs: []*ast.ElseIfClause{},
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	// Parse elseif clauses
	for p.peekTokenIs(lexer.ELSEIF) {
		p.nextToken() // move to elseif

		elseIfClause := &ast.ElseIfClause{
			Token: p.curToken,
		}

		if !p.expectPeek(lexer.LPAREN) {
			return nil
		}

		p.nextToken()
		elseIfClause.Condition = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		elseIfClause.Consequence = p.parseBlockStatement()
		stmt.ElseIfs = append(stmt.ElseIfs, elseIfClause)
	}

	// Parse else clause
	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken() // move to else

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseWhileStatement parses while loop
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseDoWhileStatement parses do-while loop
func (p *Parser) parseDoWhileStatement() *ast.DoWhileStatement {
	stmt := &ast.DoWhileStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	if !p.expectPeek(lexer.WHILE) {
		return nil
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseForStatement parses for loop
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{
		Token:     p.curToken,
		Init:      []ast.Expr{},
		Condition: []ast.Expr{},
		Increment: []ast.Expr{},
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse init expressions
	p.nextToken()
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Init = append(stmt.Init, p.parseExpression(LOWEST))

		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to expression
			stmt.Init = append(stmt.Init, p.parseExpression(LOWEST))
		}
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Parse condition expressions
	p.nextToken()
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Condition = append(stmt.Condition, p.parseExpression(LOWEST))

		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to expression
			stmt.Condition = append(stmt.Condition, p.parseExpression(LOWEST))
		}
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Parse increment expressions
	p.nextToken()
	if !p.curTokenIs(lexer.RPAREN) {
		stmt.Increment = append(stmt.Increment, p.parseExpression(LOWEST))

		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to expression
			stmt.Increment = append(stmt.Increment, p.parseExpression(LOWEST))
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseForeachStatement parses foreach loop
func (p *Parser) parseForeachStatement() *ast.ForeachStatement {
	stmt := &ast.ForeachStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Array = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.AS) {
		return nil
	}

	p.nextToken()

	// Check for key => value syntax
	firstExpr := p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.DOUBLE_ARROW) {
		// Has key
		stmt.Key = firstExpr

		p.nextToken() // consume =>
		p.nextToken() // move to value

		// Check for reference (&$value)
		if p.curTokenIs(lexer.BITWISE_AND) {
			stmt.ByRef = true
			p.nextToken()
		}

		stmt.Value = p.parseExpression(LOWEST)
	} else {
		// No key, just value
		stmt.Value = firstExpr
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseSwitchStatement parses switch statement
func (p *Parser) parseSwitchStatement() *ast.SwitchStatement {
	stmt := &ast.SwitchStatement{
		Token: p.curToken,
		Cases: []*ast.SwitchCase{},
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Subject = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()

	// Parse cases
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		if p.curTokenIs(lexer.CASE) {
			caseClause := &ast.SwitchCase{
				Token: p.curToken,
				Body:  []ast.Stmt{},
			}

			p.nextToken()
			caseClause.Value = p.parseExpression(LOWEST)

			if !p.expectPeek(lexer.COLON) {
				return nil
			}

			p.nextToken()

			// Parse case body until next case/default/closing brace
			for !p.curTokenIs(lexer.CASE) && !p.curTokenIs(lexer.DEFAULT) &&
				!p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
				stmt := p.parseStatement()
				if stmt != nil {
					caseClause.Body = append(caseClause.Body, stmt)
				}
				p.nextToken()
			}

			stmt.Cases = append(stmt.Cases, caseClause)
		} else if p.curTokenIs(lexer.DEFAULT) {
			defaultClause := &ast.SwitchCase{
				Token: p.curToken,
				Value: nil, // nil indicates default case
				Body:  []ast.Stmt{},
			}

			if !p.expectPeek(lexer.COLON) {
				return nil
			}

			p.nextToken()

			// Parse default body
			for !p.curTokenIs(lexer.CASE) && !p.curTokenIs(lexer.DEFAULT) &&
				!p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
				stmt := p.parseStatement()
				if stmt != nil {
					defaultClause.Body = append(defaultClause.Body, stmt)
				}
				p.nextToken()
			}

			stmt.Cases = append(stmt.Cases, defaultClause)
		} else {
			p.nextToken()
		}
	}

	return stmt
}

// parseMatchExpression parses match expression (PHP 8.0+)
func (p *Parser) parseMatchExpression() ast.Expr {
	matchExpr := &ast.MatchExpression{
		Token: p.curToken,
		Arms:  []*ast.MatchArm{},
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	matchExpr.Subject = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()

	// Parse match arms
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		arm := &ast.MatchArm{
			Conditions: []ast.Expr{},
		}

		// Check for default
		if p.curTokenIs(lexer.DEFAULT) {
			arm.IsDefault = true
		} else {
			// Parse conditions
			arm.Conditions = append(arm.Conditions, p.parseExpression(LOWEST))

			// Handle multiple conditions separated by comma
			for p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // consume comma

				// Check if this is the end (trailing comma) or another condition
				if p.peekTokenIs(lexer.RBRACE) || p.peekTokenIs(lexer.DEFAULT) {
					break
				}

				p.nextToken() // move to expression
				arm.Conditions = append(arm.Conditions, p.parseExpression(LOWEST))
			}
		}

		if !p.expectPeek(lexer.DOUBLE_ARROW) {
			return nil
		}

		p.nextToken()
		arm.Body = p.parseExpression(LOWEST)

		matchExpr.Arms = append(matchExpr.Arms, arm)

		// Check for comma (more arms) or closing brace
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma

			// Check if there's another arm or just trailing comma
			if !p.peekTokenIs(lexer.RBRACE) {
				p.nextToken() // move to next arm
			}
		} else if p.peekTokenIs(lexer.RBRACE) {
			break
		}
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return matchExpr
}

// parseTryStatement parses try-catch-finally statement
func (p *Parser) parseTryStatement() *ast.TryStatement {
	stmt := &ast.TryStatement{
		Token:        p.curToken,
		CatchClauses: []*ast.CatchClause{},
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	// Parse catch clauses
	for p.peekTokenIs(lexer.CATCH) {
		p.nextToken() // move to catch

		catchClause := &ast.CatchClause{
			Token: p.curToken,
			Types: []ast.Expr{},
		}

		if !p.expectPeek(lexer.LPAREN) {
			return nil
		}

		p.nextToken()

		// Parse exception types (can be multiple with |)
		catchClause.Types = append(catchClause.Types, p.parseExpression(POSTFIX))

		for p.peekTokenIs(lexer.BITWISE_OR) {
			p.nextToken() // consume |
			p.nextToken() // move to next type
			catchClause.Types = append(catchClause.Types, p.parseExpression(POSTFIX))
		}

		// Optional variable
		if p.peekTokenIs(lexer.VARIABLE) {
			p.nextToken()
			variable, ok := p.parseExpression(POSTFIX).(*ast.Variable)
			if ok {
				catchClause.Variable = variable
			}
		}

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		catchClause.Body = p.parseBlockStatement()
		stmt.CatchClauses = append(stmt.CatchClauses, catchClause)
	}

	// Parse finally clause
	if p.peekTokenIs(lexer.FINALLY) {
		p.nextToken() // move to finally

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		stmt.Finally = p.parseBlockStatement()
	}

	return stmt
}

// parseThrowStatement parses throw statement
func (p *Parser) parseThrowStatement() *ast.ThrowStatement {
	stmt := &ast.ThrowStatement{
		Token: p.curToken,
	}

	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parses a block of statements { ... }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Stmt{},
	}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseNamespaceStatement parses namespace declaration
// Syntax: namespace Name\Space;  or  namespace Name\Space { ... }  or  namespace { ... }
func (p *Parser) parseNamespaceStatement() *ast.NamespaceStatement {
	stmt := &ast.NamespaceStatement{
		Token:      p.curToken,
		Statements: []ast.Stmt{},
	}

	p.nextToken()

	// Check for global namespace: namespace { ... }
	if p.curTokenIs(lexer.LBRACE) {
		stmt.Body = p.parseBlockStatement()
		return stmt
	}

	// Parse namespace name
	stmt.Name = p.parseNamespaceName()

	// Check for bracketed or unbracketed syntax
	if p.peekTokenIs(lexer.LBRACE) {
		// Bracketed syntax: namespace Name { ... }
		p.nextToken()
		stmt.Body = p.parseBlockStatement()
	} else {
		// Unbracketed syntax: namespace Name;
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
		// Parse all following statements until EOF or next namespace
		for !p.peekTokenIs(lexer.EOF) && !p.peekTokenIs(lexer.NAMESPACE) {
			p.nextToken()
			if s := p.parseStatement(); s != nil {
				stmt.Statements = append(stmt.Statements, s)
			}
		}
	}

	return stmt
}

// parseNamespaceName parses a namespace or class name with backslash separators
// e.g., Foo\Bar\Baz or \Foo\Bar (leading backslash for global namespace)
// Note: The lexer scans Foo\Bar\Baz as a single IDENT token including backslashes
func (p *Parser) parseNamespaceName() *ast.NamespaceName {
	name := &ast.NamespaceName{
		Token: p.curToken,
		Parts: []string{},
	}

	if !p.curTokenIs(lexer.IDENT) {
		p.error("expected identifier in namespace name")
		return name
	}

	// Split the identifier by backslashes to get parts
	// The lexer includes backslashes in the identifier literal
	literal := p.curToken.Literal

	// Split on backslash
	parts := []string{}
	current := ""
	for _, ch := range literal {
		if ch == '\\' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	name.Parts = parts
	return name
}

// parseUseStatement parses use declaration(s)
// Syntax: use Name\Space;  use Name\Space as Alias;  use Name\{A, B};  use function Name;  use const Name;
func (p *Parser) parseUseStatement() *ast.UseStatement {
	stmt := &ast.UseStatement{
		Token: p.curToken,
		Uses:  []*ast.UseImport{},
	}

	p.nextToken()

	// Check for function or const type
	if p.curTokenIs(lexer.FUNCTION) {
		stmt.Type = "function"
		p.nextToken()
	} else if p.curTokenIs(lexer.CONST) {
		stmt.Type = "const"
		p.nextToken()
	}

	// Parse the namespace name
	startName := p.parseNamespaceName()

	// Check for group use syntax: use Name\{A, B, C};
	if p.peekTokenIs(lexer.LBRACE) {
		p.nextToken() // consume name
		p.nextToken() // consume {

		stmt.Prefix = startName.String()

		// Parse group use clauses
		for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
			clause := &ast.UseImport{}

			// Check for function/const prefix in group
			if p.curTokenIs(lexer.FUNCTION) || p.curTokenIs(lexer.CONST) {
				p.nextToken()
			}

			// Parse the name
			clause.Name = p.parseNamespaceName()

			// Check for alias
			if p.peekTokenIs(lexer.AS) {
				p.nextToken() // consume name
				p.nextToken() // consume 'as'

				if !p.curTokenIs(lexer.IDENT) {
					p.error("expected identifier after 'as'")
				} else {
					clause.Alias = p.curToken.Literal
				}
			}

			stmt.Uses = append(stmt.Uses, clause)

			// Check for comma or end of group
			if p.peekTokenIs(lexer.COMMA) {
				p.nextToken() // consume comma
				p.nextToken() // move to next name
			} else if p.peekTokenIs(lexer.RBRACE) {
				p.nextToken() // move to }
				break
			}
		}

		// Expect closing brace and semicolon
		if !p.curTokenIs(lexer.RBRACE) {
			p.error("expected } after group use")
		}
	} else {
		// Simple use: use Name\Space [as Alias];
		clause := &ast.UseImport{
			Name: startName,
		}

		// Check for alias
		if p.peekTokenIs(lexer.AS) {
			p.nextToken() // consume name
			p.nextToken() // consume 'as'

			if !p.curTokenIs(lexer.IDENT) {
				p.error("expected identifier after 'as'")
			} else {
				clause.Alias = p.curToken.Literal
			}
		}

		stmt.Uses = append(stmt.Uses, clause)

		// Handle multiple use declarations separated by commas
		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next name

			clause := &ast.UseImport{
				Name: p.parseNamespaceName(),
			}

			// Check for alias
			if p.peekTokenIs(lexer.AS) {
				p.nextToken() // consume name
				p.nextToken() // consume 'as'

				if !p.curTokenIs(lexer.IDENT) {
					p.error("expected identifier after 'as'")
				} else {
					clause.Alias = p.curToken.Literal
				}
			}

			stmt.Uses = append(stmt.Uses, clause)
		}
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
