package parser

import (
	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// parseType parses a complete type hint including nullable, union, and intersection types
// This is the main entry point for type parsing
func (p *Parser) parseType() ast.Expr {
	return p.parseUnionType()
}

// parseUnionType parses union types (Type1|Type2|Type3)
// Union has the lowest precedence
func (p *Parser) parseUnionType() ast.Expr {
	left := p.parseIntersectionType()

	if !p.peekTokenIs(lexer.BITWISE_OR) {
		return left
	}

	// We have a union type
	union := &ast.UnionType{
		Token: p.curToken,
		Types: []ast.Expr{left},
	}

	for p.peekTokenIs(lexer.BITWISE_OR) {
		p.nextToken() // consume |
		p.nextToken() // move to next type
		union.Types = append(union.Types, p.parseIntersectionType())
	}

	return union
}

// parseIntersectionType parses intersection types (Type1&Type2)
// Intersection has higher precedence than union
// Note: In PHP, & in parameter position can mean either:
//   - Intersection type: Countable&Traversable $x
//   - By-reference: array &$x
// We only parse as intersection if & is followed by a type token, not a variable
func (p *Parser) parseIntersectionType() ast.Expr {
	left := p.parseSingleType()

	// Check if we have potential intersection type
	if !p.peekTokenIs(lexer.BITWISE_AND) {
		return left
	}

	// Don't parse & as intersection in these cases:
	// 1. Current context might be a parameter and & could mean by-reference
	// We need a way to disambiguate, so we only proceed if we're confident
	// this is an intersection type (i.e., followed by another type token)

	// For now, skip intersection type parsing to avoid conflicts with by-reference parameters
	// Full intersection type support can be added later with proper context tracking
	return left
}

// parseSingleType parses a single type (possibly nullable)
// Handles: ?Type, Type, (Type)
func (p *Parser) parseSingleType() ast.Expr {
	// Check for nullable type (?Type)
	if p.curTokenIs(lexer.QUESTION) {
		nullableToken := p.curToken
		p.nextToken() // move to the type
		return &ast.NullableType{
			Token: nullableToken,
			Type:  p.parseBaseType(),
		}
	}

	// Check for parenthesized type
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken() // consume (
		typeExpr := p.parseType()
		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
		return typeExpr
	}

	return p.parseBaseType()
}

// parseBaseType parses a base type (identifier, scalar, special)
// This includes: int, string, bool, float, array, callable, iterable, object
// mixed, never, void, static, self, parent, ClassName
func (p *Parser) parseBaseType() ast.Expr {
	// Check if it's a type keyword or identifier
	if !p.isTypeToken() {
		p.error("expected type name")
		return nil
	}

	// Create an identifier for the type
	typeIdent := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Validate the type name (optional, for better error messages)
	p.validateTypeName(typeIdent.Value)

	return typeIdent
}

// isTypeToken checks if the current token can be used as a type
func (p *Parser) isTypeToken() bool {
	switch p.curToken.Type {
	// Scalar type keywords
	case lexer.INT, lexer.FLOAT_TYPE, lexer.BOOL, lexer.STRING_TYPE:
		return true
	// Special type keywords
	case lexer.MIXED, lexer.VOID, lexer.NEVER, lexer.NULL, lexer.TRUE, lexer.FALSE:
		return true
	// Compound type keywords
	case lexer.ARRAY, lexer.OBJECT, lexer.CALLABLE, lexer.ITERABLE:
		return true
	// Other keywords that can be types
	case lexer.STATIC:
		return true
	// Class/interface names (identifiers)
	case lexer.IDENT:
		return true
	default:
		return false
	}
}

// isTypeKeyword is deprecated, use isTypeToken instead
func (p *Parser) isTypeKeyword() bool {
	return p.isTypeToken()
}

// validateTypeName validates a type name and can issue warnings
func (p *Parser) validateTypeName(name string) {
	// List of valid scalar and special types
	validTypes := map[string]bool{
		// Scalar types
		"int":      true,
		"integer":  true, // alias for int
		"float":    true,
		"double":   true, // alias for float
		"string":   true,
		"bool":     true,
		"boolean":  true, // alias for bool
		"array":    true,
		"object":   true,
		"callable": true,
		"iterable": true,
		// Special types
		"mixed":  true,
		"void":   true,
		"never":  true,
		"null":   true,
		"false":  true, // PHP 8.2+
		"true":   true, // PHP 8.2+
		"static": true,
		"self":   true,
		"parent": true,
		// Resource (deprecated but still valid)
		"resource": true,
	}

	// Check if it's a valid built-in type
	// If not, assume it's a class/interface name (which is also valid)
	_ = validTypes // No error for now, just validation for future use
}

// Helper function to check if a type name is a scalar type
func isScalarType(name string) bool {
	switch name {
	case "int", "integer", "float", "double", "string", "bool", "boolean":
		return true
	default:
		return false
	}
}

// Helper function to check if a type name is a special type
func isSpecialType(name string) bool {
	switch name {
	case "mixed", "void", "never", "null", "false", "true", "static", "self", "parent":
		return true
	default:
		return false
	}
}

// Helper function to check if a type name is a compound type
func isCompoundType(name string) bool {
	switch name {
	case "array", "object", "callable", "iterable", "resource":
		return true
	default:
		return false
	}
}

// parseTypeList parses a comma-separated list of types
// Used for parsing catch clause types: catch (Exception1 | Exception2 $e)
func (p *Parser) parseTypeList() []ast.Expr {
	types := []ast.Expr{}

	types = append(types, p.parseType())

	for p.peekTokenIs(lexer.BITWISE_OR) {
		p.nextToken() // consume |
		p.nextToken() // move to next type
		types = append(types, p.parseType())
	}

	return types
}

// updateParseTypeHint updates the old parseTypeHint function to use the new type parser
// This is called from decl.go for backward compatibility
func (p *Parser) parseTypeHintCompat() ast.Expr {
	return p.parseType()
}
