package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
)

func TestArrayAppendSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "Simple array append",
			input:   `<?php $arr[] = 1;`,
			wantErr: false,
		},
		{
			name:    "Multiple array appends",
			input:   `<?php $arr[] = 1; $arr[] = 2; $arr[] = 3;`,
			wantErr: false,
		},
		{
			name:    "Array append in loop",
			input:   `<?php for ($i = 0; $i < 10; $i++) { $arr[] = $i; }`,
			wantErr: false,
		},
		{
			name:    "Array append with expression",
			input:   `<?php $arr[] = $x + $y * 2;`,
			wantErr: false,
		},
		{
			name:    "Nested array append",
			input:   `<?php $arr[][] = "value";`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()

			errors := p.Errors()
			if tt.wantErr {
				if len(errors) == 0 {
					t.Error("Expected parse errors, got none")
				}
				return
			}

			if len(errors) > 0 {
				t.Errorf("Unexpected parse errors: %v", errors)
				return
			}

			if program == nil {
				t.Fatal("ParseProgram returned nil")
			}

			if len(program.Statements) == 0 {
				t.Fatal("ParseProgram returned no statements")
			}

			// Just verify it parsed successfully
			t.Logf("Successfully parsed: %s", tt.input)
		})
	}
}
