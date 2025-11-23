package main

import (
	"fmt"
	"log"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
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

	// Parse
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for errors
	errors := p.Errors()
	if len(errors) > 0 {
		log.Fatal(errors)
	}

	// Compile
	c := compiler.New()
	err := c.Compile(program)
	if err != nil {
		log.Fatal(err)
	}

	// Get bytecode
	bytecode := c.Bytecode()

	// Print bytecode details
	fmt.Printf("Constants count: %d\n", len(bytecode.Constants))
	for i, c := range bytecode.Constants {
		fmt.Printf("  [%d] = %v (%T)\n", i, c, c)
	}
	fmt.Printf("\nInstructions count: %d\n", len(bytecode.Instructions))
	fmt.Printf("NumCVs: %d\n\n", bytecode.NumCVs)

	for i, instr := range bytecode.Instructions {
		fmt.Printf("[%3d] %s (Line %d)\n",
			i, instr.String(), instr.Lineno)
	}

	// Run
	fmt.Println("\n=== Execution ===")
	machine := vm.New()
	machine.LoadConstants(bytecode.Constants)
	err = machine.ExecuteWithCVs(bytecode.Instructions, bytecode.NumCVs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nOutput: %s\n", machine.GetOutput())
}
