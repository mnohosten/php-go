package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestAttributeParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedAttrs  int
		expectedArgs   int
		expectedNamed  int
		attrName       string
	}{
		{
			name:          "simple attribute on class",
			input:         `<?php #[Attribute] class Foo {}`,
			expectedAttrs: 1,
			attrName:      "Attribute",
		},
		{
			name:          "attribute with single argument",
			input:         `<?php #[Route("/api")] class Foo {}`,
			expectedAttrs: 1,
			expectedArgs:  1,
			attrName:      "Route",
		},
		{
			name:          "attribute with multiple arguments",
			input:         `<?php #[Route("/api", "GET")] class Foo {}`,
			expectedAttrs: 1,
			expectedArgs:  2,
			attrName:      "Route",
		},
		{
			name:           "attribute with named argument",
			input:          `<?php #[Route(path: "/api")] class Foo {}`,
			expectedAttrs:  1,
			expectedNamed:  1,
			attrName:       "Route",
		},
		{
			name:           "attribute with mixed arguments",
			input:          `<?php #[Route("/api", methods: ["GET", "POST"])] class Foo {}`,
			expectedAttrs:  1,
			expectedArgs:   1,
			expectedNamed:  1,
			attrName:       "Route",
		},
		{
			name:          "multiple attributes on class",
			input:         `<?php #[Attribute] #[Route("/api")] class Foo {}`,
			expectedAttrs: 2,
		},
		{
			name:          "multiple attributes in single group",
			input:         `<?php #[Attr1, Attr2] class Foo {}`,
			expectedAttrs: 1, // one group with 2 attributes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
			if !ok {
				t.Fatalf("expected ClassDeclaration, got %T", program.Statements[0])
			}

			if len(classDecl.Attributes) != tt.expectedAttrs {
				t.Errorf("expected %d attribute groups, got %d", tt.expectedAttrs, len(classDecl.Attributes))
			}

			if tt.attrName != "" && len(classDecl.Attributes) > 0 {
				attr := classDecl.Attributes[0].Attributes[0]
				name, ok := attr.Name.(*ast.Identifier)
				if !ok {
					t.Fatalf("expected Identifier for attribute name, got %T", attr.Name)
				}
				if name.Value != tt.attrName {
					t.Errorf("expected attribute name %q, got %q", tt.attrName, name.Value)
				}

				if tt.expectedArgs > 0 && len(attr.Arguments) != tt.expectedArgs {
					t.Errorf("expected %d arguments, got %d", tt.expectedArgs, len(attr.Arguments))
				}

				if tt.expectedNamed > 0 && len(attr.Named) != tt.expectedNamed {
					t.Errorf("expected %d named arguments, got %d", tt.expectedNamed, len(attr.Named))
				}
			}
		})
	}
}

func TestAttributeOnProperty(t *testing.T) {
	input := `<?php
class User {
    #[Column("id")]
    public int $id;
}`
	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("expected ClassDeclaration, got %T", program.Statements[0])
	}

	if len(classDecl.Body) != 1 {
		t.Fatalf("expected 1 class member, got %d", len(classDecl.Body))
	}

	propDecl, ok := classDecl.Body[0].(*ast.PropertyDeclaration)
	if !ok {
		t.Fatalf("expected PropertyDeclaration, got %T", classDecl.Body[0])
	}

	if len(propDecl.Attributes) != 1 {
		t.Errorf("expected 1 attribute group on property, got %d", len(propDecl.Attributes))
	}
}

func TestAttributeOnMethod(t *testing.T) {
	input := `<?php
class Controller {
    #[Get("/users")]
    #[Authenticated]
    public function getUsers(): array {
        return [];
    }
}`
	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("expected ClassDeclaration, got %T", program.Statements[0])
	}

	if len(classDecl.Body) != 1 {
		t.Fatalf("expected 1 class member, got %d", len(classDecl.Body))
	}

	methodDecl, ok := classDecl.Body[0].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("expected MethodDeclaration, got %T", classDecl.Body[0])
	}

	if len(methodDecl.Attributes) != 2 {
		t.Errorf("expected 2 attribute groups on method, got %d", len(methodDecl.Attributes))
	}
}

func TestAttributeOnFunction(t *testing.T) {
	input := `<?php
#[Pure]
function calculate(int $x): int {
    return $x * 2;
}`
	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	funcDecl, ok := program.Statements[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expected FunctionDeclaration, got %T", program.Statements[0])
	}

	if len(funcDecl.Attributes) != 1 {
		t.Errorf("expected 1 attribute group on function, got %d", len(funcDecl.Attributes))
	}
}

func TestNamespacedAttribute(t *testing.T) {
	input := `<?php
#[\Doctrine\ORM\Mapping\Entity]
class User {}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	classDecl, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("expected ClassDeclaration, got %T", program.Statements[0])
	}

	if len(classDecl.Attributes) != 1 {
		t.Errorf("expected 1 attribute group, got %d", len(classDecl.Attributes))
	}

	attr := classDecl.Attributes[0].Attributes[0]
	name, ok := attr.Name.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected Identifier for namespaced attribute, got %T", attr.Name)
	}

	expectedName := "\\Doctrine\\ORM\\Mapping\\Entity"
	if name.Value != expectedName {
		t.Errorf("expected attribute name %q, got %q", expectedName, name.Value)
	}
}
