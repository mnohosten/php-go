package main

import (
	"fmt"
	"log"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/runtime"
	"github.com/krizos/php-go/pkg/vm"
)

func main() {
	code := `<?php
$i = 0;
while ($i < 2) {
    echo $i;
    $i = $i + 1;
}
`

	// Lex
	l := lexer.New(code, "test.php")
	tokens, err := l.ScanAll()
	if err != nil {
		log.Fatal(err)
	}

	// Parse
	p := parser.New(tokens, "test.php")
	ast, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Compile
	c := compiler.New()
	bytecode, err := c.Compile(ast)
	if err != nil {
		log.Fatal(err)
	}

	// Print bytecode details
	fmt.Printf("Constants count: %d\n", len(bytecode.Constants))
	fmt.Printf("Constants: %v\n", bytecode.Constants)
	fmt.Printf("\nInstructions count: %d\n", len(bytecode.Instructions))

	for i, instr := range bytecode.Instructions {
		fmt.Printf("[%d] Op=%s, Op1=%d, Op2=%d, Op3=%d, Line=%d\n",
			i, vm.OpcodeName(instr.Opcode), instr.Operand1, instr.Operand2, instr.Operand3, instr.Line)
	}

	// Run
	fmt.Println("\n=== Execution ===")
	rt := runtime.New()
	v := vm.New(bytecode, rt)
	result, err := v.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %v\n", result)
	fmt.Printf("Output: %s\n", rt.GetOutput())
}
