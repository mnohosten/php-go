package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
)

func TestParserNew(t *testing.T) {
	input := `<?php`
	l := lexer.New(input, "test.php")
	p := New(l)

	if p == nil {
		t.Fatal("New() returned nil")
	}

	if p.l == nil {
		t.Error("parser lexer is nil")
	}

	if len(p.errors) != 0 {
		t.Errorf("parser has errors on initialization: %v", p.errors)
	}
}

func TestTokenManagement(t *testing.T) {
	input := `<?php $x = 5;`
	l := lexer.New(input, "test.php")
	p := New(l)

	// Check initial tokens
	if p.curToken.Type != lexer.OPEN_TAG {
		t.Errorf("curToken type wrong. expected=OPEN_TAG, got=%s", p.curToken.Type)
	}

	if p.peekToken.Type != lexer.VARIABLE {
		t.Errorf("peekToken type wrong. expected=VARIABLE, got=%s", p.peekToken.Type)
	}

	// Test nextToken
	p.nextToken()
	if p.curToken.Type != lexer.VARIABLE {
		t.Errorf("after nextToken, curToken type wrong. expected=VARIABLE, got=%s", p.curToken.Type)
	}

	if p.peekToken.Type != lexer.ASSIGN {
		t.Errorf("after nextToken, peekToken type wrong. expected=ASSIGN, got=%s", p.peekToken.Type)
	}
}

func TestCurTokenIs(t *testing.T) {
	input := `<?php $x`
	l := lexer.New(input, "test.php")
	p := New(l)

	if !p.curTokenIs(lexer.OPEN_TAG) {
		t.Error("curTokenIs(OPEN_TAG) should be true")
	}

	if p.curTokenIs(lexer.VARIABLE) {
		t.Error("curTokenIs(VARIABLE) should be false")
	}
}

func TestPeekTokenIs(t *testing.T) {
	input := `<?php $x`
	l := lexer.New(input, "test.php")
	p := New(l)

	if !p.peekTokenIs(lexer.VARIABLE) {
		t.Error("peekTokenIs(VARIABLE) should be true")
	}

	if p.peekTokenIs(lexer.OPEN_TAG) {
		t.Error("peekTokenIs(OPEN_TAG) should be false")
	}
}

func TestExpectPeek(t *testing.T) {
	input := `<?php $x = 5;`
	l := lexer.New(input, "test.php")
	p := New(l)

	// Skip opening tag
	p.nextToken()

	// Current: $x, Peek: =
	if !p.expectPeek(lexer.ASSIGN) {
		t.Error("expectPeek(ASSIGN) should succeed")
	}

	// Now current should be =
	if !p.curTokenIs(lexer.ASSIGN) {
		t.Error("after expectPeek, curToken should be ASSIGN")
	}

	// Try to expect wrong token
	if p.expectPeek(lexer.LPAREN) {
		t.Error("expectPeek(LPAREN) should fail")
	}

	// Should have an error now
	if !p.HasErrors() {
		t.Error("parser should have errors after failed expectPeek")
	}

	if len(p.Errors()) != 1 {
		t.Errorf("parser should have exactly 1 error, got %d", len(p.Errors()))
	}
}

func TestExpectCurrent(t *testing.T) {
	input := `<?php $x`
	l := lexer.New(input, "test.php")
	p := New(l)

	if !p.expectCurrent(lexer.OPEN_TAG) {
		t.Error("expectCurrent(OPEN_TAG) should succeed")
	}

	if p.expectCurrent(lexer.VARIABLE) {
		t.Error("expectCurrent(VARIABLE) should fail")
	}

	if !p.HasErrors() {
		t.Error("parser should have errors after failed expectCurrent")
	}
}

func TestErrorHandling(t *testing.T) {
	input := `<?php`
	l := lexer.New(input, "test.php")
	p := New(l)

	// Test error method
	p.error("test error")

	if !p.HasErrors() {
		t.Error("HasErrors() should return true after error()")
	}

	errors := p.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}

	if errors[0] == "" {
		t.Error("error message should not be empty")
	}
}

func TestPeekError(t *testing.T) {
	input := `<?php $x`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.peekError(lexer.SEMICOLON)

	if !p.HasErrors() {
		t.Error("HasErrors() should return true after peekError()")
	}

	errors := p.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}
}

func TestSynchronize(t *testing.T) {
	input := `<?php $x = invalid syntax here; function test() {}`
	l := lexer.New(input, "test.php")
	p := New(l)

	// Skip to somewhere in the middle
	p.nextToken() // $x
	p.nextToken() // =
	p.nextToken() // invalid

	// Synchronize should skip to next statement
	p.synchronize()

	// Should now be at or near 'function'
	found := false
	for i := 0; i < 5; i++ {
		if p.curTokenIs(lexer.FUNCTION) || p.peekTokenIs(lexer.FUNCTION) {
			found = true
			break
		}
		p.nextToken()
	}

	if !found {
		t.Error("synchronize() should have moved to next statement boundary")
	}
}

func TestSkipTo(t *testing.T) {
	input := `<?php $x = 5;`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.skipTo(lexer.SEMICOLON)

	if !p.curTokenIs(lexer.SEMICOLON) {
		t.Errorf("skipTo(SEMICOLON) failed. current token is %s", p.curToken.Type)
	}
}

func TestSkipToStatementEnd(t *testing.T) {
	input := `<?php $x = 5 + 3;`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.nextToken() // $x

	p.skipToStatementEnd()

	if !p.curTokenIs(lexer.SEMICOLON) {
		t.Errorf("skipToStatementEnd() failed. current token is %s", p.curToken.Type)
	}
}

func TestParseIdentifier(t *testing.T) {
	input := `<?php test`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.nextToken() // Move to 'test'

	ident := p.parseIdentifier()

	if ident == nil {
		t.Fatal("parseIdentifier() returned nil")
	}

	if ident.Value != "test" {
		t.Errorf("identifier value wrong. expected=test, got=%s", ident.Value)
	}

	if ident.TokenLiteral() != "test" {
		t.Errorf("identifier TokenLiteral() wrong. expected=test, got=%s", ident.TokenLiteral())
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		tokenType  lexer.TokenType
		precedence int
	}{
		{lexer.PLUS, SUM},
		{lexer.MINUS, SUM},
		{lexer.ASTERISK, PRODUCT},
		{lexer.SLASH, PRODUCT},
		{lexer.POWER, POWER},
		{lexer.EQ, EQUALITY},
		{lexer.IDENTICAL, EQUALITY},
		{lexer.LT, COMPARISON},
		{lexer.LOGICAL_AND, LOGICAL_AND},
		{lexer.LOGICAL_OR, LOGICAL_OR},
		{lexer.ASSIGN, ASSIGNMENT},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType.String(), func(t *testing.T) {
			if precedences[tt.tokenType] != tt.precedence {
				t.Errorf("precedence for %s wrong. expected=%d, got=%d",
					tt.tokenType, tt.precedence, precedences[tt.tokenType])
			}
		})
	}
}

func TestCurrentTokenPrecedence(t *testing.T) {
	input := `<?php + - *`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.nextToken() // +

	if p.currentTokenPrecedence() != SUM {
		t.Errorf("precedence wrong. expected=%d, got=%d", SUM, p.currentTokenPrecedence())
	}

	p.nextToken() // -

	if p.currentTokenPrecedence() != SUM {
		t.Errorf("precedence wrong. expected=%d, got=%d", SUM, p.currentTokenPrecedence())
	}

	p.nextToken() // *

	if p.currentTokenPrecedence() != PRODUCT {
		t.Errorf("precedence wrong. expected=%d, got=%d", PRODUCT, p.currentTokenPrecedence())
	}
}

func TestPeekTokenPrecedence(t *testing.T) {
	input := `<?php + *`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.nextToken() // +, peek is *

	if p.peekTokenPrecedence() != PRODUCT {
		t.Errorf("peek precedence wrong. expected=%d, got=%d", PRODUCT, p.peekTokenPrecedence())
	}
}

func TestParseProgram(t *testing.T) {
	input := `<?php`
	l := lexer.New(input, "test.php")
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	if program.Statements == nil {
		t.Error("program.Statements is nil")
	}

	// With no statements, should have empty slice
	if len(program.Statements) != 0 {
		t.Errorf("program has wrong number of statements. expected=0, got=%d", len(program.Statements))
	}
}

func TestParseString(t *testing.T) {
	input := `<?php`

	program, errors := ParseString(input)

	if program == nil {
		t.Fatal("ParseString() returned nil program")
	}

	if errors == nil {
		t.Error("ParseString() returned nil errors slice")
	}

	if len(errors) != 0 {
		t.Errorf("ParseString() returned unexpected errors: %v", errors)
	}
}

func TestCommentSkipping(t *testing.T) {
	input := `<?php
	// This is a comment
	/* This is also a comment */
	$x
	`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.nextToken() // Should skip comments and land on $x

	if !p.curTokenIs(lexer.VARIABLE) {
		t.Errorf("comments not skipped properly. current token is %s", p.curToken.Type)
	}

	if p.curToken.Literal != "$x" {
		t.Errorf("variable literal wrong. expected=$x, got=%s", p.curToken.Literal)
	}
}

func TestErrorMessageFormat(t *testing.T) {
	input := `<?php $x`
	l := lexer.New(input, "test.php")
	p := New(l)

	p.error("test error message")

	errors := p.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	// Error message should include position and message
	if errors[0] == "" {
		t.Error("error message is empty")
	}

	// Should contain "Parse error"
	if len(errors[0]) < 10 {
		t.Errorf("error message too short: %s", errors[0])
	}
}

// Benchmark tests
func BenchmarkParserNew(b *testing.B) {
	input := `<?php $x = 5;`

	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		New(l)
	}
}

func BenchmarkParseProgram(b *testing.B) {
	input := `<?php
	$x = 5;
	$y = 10;
	$sum = $x + $y;
	`

	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}
