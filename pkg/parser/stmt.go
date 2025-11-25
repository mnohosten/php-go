package parser

import (
	"fmt"

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

	// Check for alternative syntax (colon), regular syntax (brace), or single statement
	p.nextToken()
	useAlternativeSyntax := false

	if p.curTokenIs(lexer.COLON) {
		// Alternative syntax: if (...): ... endif;
		useAlternativeSyntax = true
		stmt.Consequence = p.parseAlternativeBlockStatement(lexer.ENDIF, lexer.ELSEIF, lexer.ELSE)
		// After parseAlternativeBlockStatement, curToken is at the terminator (endif/elseif/else)
	} else if p.curTokenIs(lexer.LBRACE) {
		// Regular syntax: if (...) { ... }
		stmt.Consequence = p.parseBlockStatement()
	} else {
		// Single statement without braces: if (...) statement;
		stmt.Consequence = p.parseSingleStatement()
	}

	// Parse elseif clauses
	// In alternative syntax, curToken is already at elseif/else/endif after parseAlternativeBlockStatement
	// In regular syntax, we need to peek for the next token
	for (useAlternativeSyntax && p.curTokenIs(lexer.ELSEIF)) || (!useAlternativeSyntax && p.peekTokenIs(lexer.ELSEIF)) {
		if !useAlternativeSyntax {
			p.nextToken() // move to elseif in regular syntax
		}

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

		p.nextToken()
		if useAlternativeSyntax {
			if !p.curTokenIs(lexer.COLON) {
				p.peekError(lexer.COLON)
				return nil
			}
			elseIfClause.Consequence = p.parseAlternativeBlockStatement(lexer.ENDIF, lexer.ELSEIF, lexer.ELSE)
			// After parseAlternativeBlockStatement, curToken is at the terminator
		} else if p.curTokenIs(lexer.LBRACE) {
			// Regular syntax with braces: elseif (...) { ... }
			elseIfClause.Consequence = p.parseBlockStatement()
		} else {
			// Single statement without braces: elseif (...) statement;
			elseIfClause.Consequence = p.parseSingleStatement()
		}

		stmt.ElseIfs = append(stmt.ElseIfs, elseIfClause)
	}

	// Parse else clause
	if (useAlternativeSyntax && p.curTokenIs(lexer.ELSE)) || (!useAlternativeSyntax && p.peekTokenIs(lexer.ELSE)) {
		if !useAlternativeSyntax {
			p.nextToken() // move to else in regular syntax
		}

		if useAlternativeSyntax {
			if !p.expectPeek(lexer.COLON) {
				return nil
			}
			stmt.Alternative = p.parseAlternativeBlockStatement(lexer.ENDIF)
			// After parseAlternativeBlockStatement, curToken is at endif
		} else {
			p.nextToken() // move past else
			if p.curTokenIs(lexer.LBRACE) {
				// Regular syntax with braces: else { ... }
				stmt.Alternative = p.parseBlockStatement()
			} else {
				// Single statement without braces: else statement;
				stmt.Alternative = p.parseSingleStatement()
			}
		}
	}

	// Consume endif token for alternative syntax
	if useAlternativeSyntax {
		if !p.curTokenIs(lexer.ENDIF) {
			p.peekError(lexer.ENDIF)
			return nil
		}
		// Optional semicolon after endif
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	return stmt
}

// parseWhileStatement parses while loop (both regular and alternative syntax)
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

	// Check for alternative syntax (colon), regular syntax (brace), or single statement
	p.nextToken()

	if p.curTokenIs(lexer.COLON) {
		// Alternative syntax: while (...): ... endwhile;
		stmt.Body = p.parseAlternativeBlockStatement(lexer.ENDWHILE)
		if !p.curTokenIs(lexer.ENDWHILE) {
			p.peekError(lexer.ENDWHILE)
			return nil
		}
		// Optional semicolon after endwhile
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.LBRACE) {
		// Regular syntax: while (...) { ... }
		stmt.Body = p.parseBlockStatement()
	} else {
		// Single statement without braces: while (...) statement;
		stmt.Body = p.parseSingleStatement()
	}

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

	// Check for alternative syntax (colon), regular syntax (brace), or single statement
	p.nextToken()

	if p.curTokenIs(lexer.COLON) {
		// Alternative syntax: for (...): ... endfor;
		stmt.Body = p.parseAlternativeBlockStatement(lexer.ENDFOR)
		if !p.curTokenIs(lexer.ENDFOR) {
			p.peekError(lexer.ENDFOR)
			return nil
		}
		// Optional semicolon after endfor
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.LBRACE) {
		// Regular syntax: for (...) { ... }
		stmt.Body = p.parseBlockStatement()
	} else {
		// Single statement without braces: for (...) statement;
		stmt.Body = p.parseSingleStatement()
	}

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

	// Check for alternative syntax (colon), regular syntax (brace), or single statement
	p.nextToken()

	if p.curTokenIs(lexer.COLON) {
		// Alternative syntax: foreach (...): ... endforeach;
		stmt.Body = p.parseAlternativeBlockStatement(lexer.ENDFOREACH)
		if !p.curTokenIs(lexer.ENDFOREACH) {
			p.peekError(lexer.ENDFOREACH)
			return nil
		}
		// Optional semicolon after endforeach
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.LBRACE) {
		// Regular syntax: foreach (...) { ... }
		stmt.Body = p.parseBlockStatement()
	} else {
		// Single statement without braces: foreach (...) statement;
		stmt.Body = p.parseSingleStatement()
	}

	return stmt
}

// parseSwitchStatement parses switch statement (both regular and alternative syntax)
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

	// Check for alternative syntax (colon) or regular syntax (brace)
	p.nextToken()
	useAlternativeSyntax := false

	if p.curTokenIs(lexer.COLON) {
		// Alternative syntax: switch (...): ... endswitch;
		useAlternativeSyntax = true
		p.nextToken()
	} else if p.curTokenIs(lexer.LBRACE) {
		// Regular syntax: switch (...) { ... }
		p.nextToken()
	} else {
		p.peekError(lexer.LBRACE)
		return nil
	}

	// Parse cases
	endToken := lexer.RBRACE
	if useAlternativeSyntax {
		endToken = lexer.ENDSWITCH
	}

	for !p.curTokenIs(endToken) && !p.curTokenIs(lexer.EOF) {
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

			// Parse case body until next case/default/closing token
			for !p.curTokenIs(lexer.CASE) && !p.curTokenIs(lexer.DEFAULT) &&
				!p.curTokenIs(endToken) && !p.curTokenIs(lexer.EOF) {
				caseStmt := p.parseStatement()
				if caseStmt != nil {
					caseClause.Body = append(caseClause.Body, caseStmt)
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
				!p.curTokenIs(endToken) && !p.curTokenIs(lexer.EOF) {
				defaultStmt := p.parseStatement()
				if defaultStmt != nil {
					defaultClause.Body = append(defaultClause.Body, defaultStmt)
				}
				p.nextToken()
			}

			stmt.Cases = append(stmt.Cases, defaultClause)
		} else {
			p.nextToken()
		}
	}

	// Handle endswitch for alternative syntax
	if useAlternativeSyntax {
		if !p.curTokenIs(lexer.ENDSWITCH) {
			p.peekError(lexer.ENDSWITCH)
			return nil
		}
		// Optional semicolon after endswitch
		if p.peekTokenIs(lexer.SEMICOLON) {
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

				// Check if next is => (comma was just before the arrow)
				if p.peekTokenIs(lexer.DOUBLE_ARROW) {
					break
				}

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
			if p.peekTokenIs(lexer.RBRACE) {
				break // Trailing comma before closing brace
			}
			p.nextToken() // move to next arm
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
// parseSingleStatement parses a single statement and wraps it in a BlockStatement
// This is used for single-line if/else/while/for bodies without braces
func (p *Parser) parseSingleStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Stmt{},
	}

	// Parse one statement
	stmt := p.parseStatement()
	if stmt != nil {
		block.Statements = append(block.Statements, stmt)
	}

	return block
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Stmt{},
	}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Handle inline HTML (convert to echo statement)
		if p.curTokenIs(lexer.INLINE_HTML) {
			html := p.curToken.Literal
			echoStmt := &ast.EchoStatement{
				Token:       p.curToken,
				Expressions: []ast.Expr{&ast.StringLiteral{Token: p.curToken, Value: html}},
			}
			block.Statements = append(block.Statements, echoStmt)
			p.nextToken()
			continue
		}

		// Handle PHP close tag - just skip it
		if p.curTokenIs(lexer.CLOSE_TAG) {
			p.nextToken()
			continue
		}

		// Handle PHP open tags
		if p.curTokenIs(lexer.OPEN_TAG) || p.curTokenIs(lexer.OPEN_TAG_ECHO) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseAlternativeBlockStatement parses statements in alternative syntax (until terminator tokens)
// Used for if/endif, while/endwhile, foreach/endforeach, etc.
func (p *Parser) parseAlternativeBlockStatement(terminators ...lexer.TokenType) *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Stmt{},
	}

	p.nextToken()

	// Continue until we hit one of the terminator tokens or EOF
	for !p.curTokenIs(lexer.EOF) {
		// Check if current token is one of the terminators
		isTerminator := false
		for _, term := range terminators {
			if p.curTokenIs(term) {
				isTerminator = true
				break
			}
		}
		if isTerminator {
			break
		}

		// Handle inline HTML (convert to echo statement)
		if p.curTokenIs(lexer.INLINE_HTML) {
			html := p.curToken.Literal
			echoStmt := &ast.EchoStatement{
				Token:       p.curToken,
				Expressions: []ast.Expr{&ast.StringLiteral{Token: p.curToken, Value: html}},
			}
			block.Statements = append(block.Statements, echoStmt)
			p.nextToken()
			continue
		}

		// Handle PHP close tag - just skip it
		if p.curTokenIs(lexer.CLOSE_TAG) {
			p.nextToken()
			continue
		}

		// Handle PHP open tags
		if p.curTokenIs(lexer.OPEN_TAG) || p.curTokenIs(lexer.OPEN_TAG_ECHO) {
			p.nextToken()
			continue
		}

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
// parseDeclareStatement parses declare() statement
// Syntax: declare(directive=value);  declare(directive=value) { ... }
// Examples: declare(strict_types=1); declare(ticks=1) { ... }
func (p *Parser) parseDeclareStatement() *ast.DeclareStatement {
	stmt := &ast.DeclareStatement{
		Token:      p.curToken,
		Directives: make(map[string]interface{}),
	}

	// Expect opening parenthesis
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse directives
	for {
		p.nextToken() // move to directive name

		// Get directive name
		if !p.curTokenIs(lexer.IDENT) {
			p.error("expected directive name in declare()")
			return nil
		}
		directiveName := p.curToken.Literal

		// Expect equals sign
		if !p.expectPeek(lexer.ASSIGN) {
			return nil
		}

		p.nextToken() // move to value

		// Parse directive value (must be a literal)
		var value interface{}
		switch p.curToken.Type {
		case lexer.INTEGER:
			// Parse as integer
			intVal := 0
			fmt.Sscanf(p.curToken.Literal, "%d", &intVal)
			value = intVal
		case lexer.STRING:
			// Parse as string (remove quotes)
			value = p.curToken.Literal
		default:
			p.error("directive value must be a literal (integer or string)")
			return nil
		}

		stmt.Directives[directiveName] = value

		// Check for comma (more directives) or closing parenthesis
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			continue
		} else if p.peekTokenIs(lexer.RPAREN) {
			p.nextToken() // consume closing parenthesis
			break
		} else {
			p.error("expected ',' or ')' in declare()")
			return nil
		}
	}

	// Check for statement body or semicolon
	if p.peekTokenIs(lexer.LBRACE) {
		p.nextToken() // move to opening brace
		stmt.Body = p.parseBlockStatement()
	} else if p.peekTokenIs(lexer.COLON) {
		// Alternative syntax: declare(...): ... enddeclare;
		p.nextToken() // consume colon
		stmt.Body = p.parseAlternativeBlockStatement(lexer.ENDDECLARE)

		// Expect enddeclare token
		if !p.curTokenIs(lexer.ENDDECLARE) {
			p.error("expected 'enddeclare' after declare block")
			return nil
		}

		// Optional semicolon after enddeclare
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	} else {
		// No body, just semicolon
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	return stmt
}

// parseAttributeGroup parses a PHP 8 attribute group: #[Attr1, Attr2(arg)]
// Returns a list of attribute groups (allows multiple #[] blocks)
func (p *Parser) parseAttributeGroups() []*ast.AttributeGroup {
	var groups []*ast.AttributeGroup

	for p.curTokenIs(lexer.ATTRIBUTE_START) {
		group := p.parseAttributeGroup()
		if group != nil {
			groups = append(groups, group)
		}
		// Move past the closing bracket if we haven't already
		if p.curTokenIs(lexer.RBRACKET) {
			p.nextToken()
		}
	}

	return groups
}

// parseAttributeGroup parses a single attribute group: #[Attr1, Attr2(arg)]
func (p *Parser) parseAttributeGroup() *ast.AttributeGroup {
	group := &ast.AttributeGroup{
		Token:      p.curToken,
		Attributes: []*ast.Attribute{},
	}

	p.nextToken() // consume #[

	// Parse first attribute
	attr := p.parseAttribute()
	if attr != nil {
		group.Attributes = append(group.Attributes, attr)
	}

	// Parse additional attributes separated by commas
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		attr := p.parseAttribute()
		if attr != nil {
			group.Attributes = append(group.Attributes, attr)
		}
	}

	// Expect closing bracket
	if !p.curTokenIs(lexer.RBRACKET) {
		p.error(fmt.Sprintf("expected ']' after attribute, got %s", p.curToken.Type))
		return nil
	}

	return group
}

// parseAttribute parses a single attribute: AttributeName or AttributeName(args)
func (p *Parser) parseAttribute() *ast.Attribute {
	attr := &ast.Attribute{
		Token:     p.curToken,
		Arguments: []ast.Expr{},
		Named:     make(map[string]ast.Expr),
	}

	// Parse attribute name (can be namespaced like \Namespace\Attribute or Attribute)
	attr.Name = p.parseNamespacedName()

	// Check for arguments
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken() // consume (

		// Parse arguments
		if !p.curTokenIs(lexer.RPAREN) {
			p.parseAttributeArguments(attr)
		}

		// Expect closing paren
		if !p.curTokenIs(lexer.RPAREN) {
			p.error(fmt.Sprintf("expected ')' after attribute arguments, got %s", p.curToken.Type))
			return nil
		}
		p.nextToken() // consume )
	}

	return attr
}

// parseAttributeArguments parses attribute arguments (positional and named)
func (p *Parser) parseAttributeArguments(attr *ast.Attribute) {
	for {
		// Check for named argument: name: value
		if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.COLON) {
			name := p.curToken.Literal
			p.nextToken() // consume name
			p.nextToken() // consume :
			value := p.parseExpression(LOWEST)
			attr.Named[name] = value
		} else {
			// Positional argument
			value := p.parseExpression(LOWEST)
			attr.Arguments = append(attr.Arguments, value)
		}

		// After parseExpression, curToken might be:
		// 1. On the last token of a simple expression (literal)
		// 2. On a closing bracket/paren from complex expressions (array, call, etc.)
		// We need to check peek for comma to continue, or advance to reach )

		// If peek is comma, we have more arguments
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // move past current token
			p.nextToken() // consume comma
			continue
		}

		// If peek is ), we're done - advance to it
		if p.peekTokenIs(lexer.RPAREN) {
			p.nextToken() // move to )
			break
		}

		// If cur is already ), we're done
		if p.curTokenIs(lexer.RPAREN) {
			break
		}

		// Otherwise advance and check again
		p.nextToken()
		if p.curTokenIs(lexer.RPAREN) || p.curTokenIs(lexer.COMMA) {
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken() // consume comma
				continue
			}
			break
		}
	}
}

// parseNamespacedName parses a potentially namespaced identifier like \Foo\Bar or Foo\Bar or Foo
func (p *Parser) parseNamespacedName() ast.Expr {
	// Check for leading backslash (fully qualified name)
	prefix := ""
	if p.curTokenIs(lexer.NS_SEPARATOR) {
		prefix = "\\"
		p.nextToken()
	}

	// Parse first identifier
	if !p.curTokenIs(lexer.IDENT) {
		p.error(fmt.Sprintf("expected identifier in attribute name, got %s", p.curToken.Type))
		return nil
	}

	name := &ast.Identifier{
		Token: p.curToken,
		Value: prefix + p.curToken.Literal,
	}
	p.nextToken()

	// Check for namespace separator
	for p.curTokenIs(lexer.NS_SEPARATOR) {
		p.nextToken() // consume \
		if !p.curTokenIs(lexer.IDENT) {
			p.error("expected identifier after namespace separator")
			return name
		}
		// Build qualified name
		name.Value = name.Value + "\\" + p.curToken.Literal
		p.nextToken()
	}

	return name
}

// parseStatementWithAttributes parses a statement that may be preceded by attributes
// and returns the appropriate declaration with attributes attached
func (p *Parser) parseStatementWithAttributes() ast.Stmt {
	// Parse attribute groups
	attrs := p.parseAttributeGroups()

	// Now parse the actual declaration
	switch p.curToken.Type {
	case lexer.FUNCTION:
		fn := p.parseFunctionDeclaration()
		if fn != nil {
			fn.Attributes = attrs
		}
		return fn
	case lexer.CLASS:
		cls := p.parseClassDeclaration()
		if cls != nil {
			cls.Attributes = attrs
		}
		return cls
	case lexer.INTERFACE:
		iface := p.parseInterfaceDeclaration()
		// Interface declaration should support attributes
		return iface
	case lexer.TRAIT:
		trait := p.parseTraitDeclaration()
		// Trait declaration should support attributes
		return trait
	case lexer.ABSTRACT, lexer.FINAL:
		cls := p.parseClassDeclarationWithModifiers()
		if cls != nil {
			cls.Attributes = attrs
		}
		return cls
	case lexer.READONLY:
		// Could be readonly class
		cls := p.parseClassDeclarationWithModifiers()
		if cls != nil {
			cls.Attributes = attrs
		}
		return cls
	default:
		p.error(fmt.Sprintf("unexpected token %s after attributes", p.curToken.Type))
		return nil
	}
}
