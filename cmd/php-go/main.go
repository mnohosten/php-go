package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/vm"
)

const version = "0.0.1-dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run", "exec":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: run command requires a file argument")
			fmt.Fprintln(os.Stderr, "Usage: php-go run <file>")
			os.Exit(1)
		}
		handleRun(os.Args[2])

	case "lex":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: lex command requires a file argument")
			fmt.Fprintln(os.Stderr, "Usage: php-go lex [--json] <file>")
			os.Exit(1)
		}
		handleLex(os.Args[2:])

	case "parse":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: parse command requires a file argument")
			fmt.Fprintln(os.Stderr, "Usage: php-go parse [--json] <file>")
			os.Exit(1)
		}
		handleParse(os.Args[2:])

	case "demo":
		handleDemo()

	case "--version", "-v":
		fmt.Printf("PHP-Go v%s\n", version)
		fmt.Println("PHP 8.4 Interpreter in Go with Automatic Parallelization")

	case "--help", "-h":
		printUsage()

	default:
		// Check if it looks like a file path (contains / or ends with .php)
		if len(command) > 0 && (command[0] == '/' || command[0] == '.' || len(command) > 4 && command[len(command)-4:] == ".php") {
			handleRun(command)
		} else {
			fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func handleLex(args []string) {
	jsonOutput := false
	var filePath string

	// Parse flags
	for _, arg := range args {
		if arg == "--json" {
			jsonOutput = true
		} else if filePath == "" {
			filePath = arg
		}
	}

	if filePath == "" {
		fmt.Fprintln(os.Stderr, "Error: no file specified")
		os.Exit(1)
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	// Tokenize
	l := lexer.New(string(content), filePath)
	tokens := []lexer.Token{}

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.EOF {
			break
		}
	}

	// Output
	if jsonOutput {
		outputJSON(tokens)
	} else {
		outputTokensHuman(tokens, filePath)
	}
}

func handleParse(args []string) {
	jsonOutput := false
	var filePath string

	// Parse flags
	for _, arg := range args {
		if arg == "--json" {
			jsonOutput = true
		} else if filePath == "" {
			filePath = arg
		}
	}

	if filePath == "" {
		fmt.Fprintln(os.Stderr, "Error: no file specified")
		os.Exit(1)
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	// Parse
	l := lexer.New(string(content), filePath)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for errors
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parser encountered %d error(s):\n", len(errors))
		for i, msg := range errors {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, msg)
		}
		os.Exit(1)
	}

	// Output
	if jsonOutput {
		outputJSON(program)
	} else {
		outputASTHuman(program, filePath)
	}
}

func handleRun(filePath string) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	// Lex
	l := lexer.New(string(content), filePath)

	// Parse
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parse errors
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parse error in '%s':\n", filePath)
		for i, msg := range errors {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, msg)
		}
		os.Exit(1)
	}

	// Compile
	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error: %v\n", err)
		os.Exit(1)
	}

	// Get bytecode
	bytecode := comp.Bytecode()

	// Execute
	machine := vm.New()
	machine.LoadConstants(bytecode.Constants)
	err = machine.Execute(bytecode.Instructions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}

	// Print output
	output := machine.GetOutput()
	if output != "" {
		fmt.Print(output)
	}
}

func outputTokensHuman(tokens []lexer.Token, filePath string) {
	fmt.Printf("Tokens for: %s\n", filePath)
	fmt.Printf("Total: %d tokens\n\n", len(tokens))
	fmt.Println("┌────────┬─────────────────────┬──────────────────────────────────────┐")
	fmt.Println("│ Line   │ Type                │ Literal                              │")
	fmt.Println("├────────┼─────────────────────┼──────────────────────────────────────┤")

	for _, tok := range tokens {
		literal := tok.Literal
		if len(literal) > 36 {
			literal = literal[:33] + "..."
		}
		// Escape newlines and tabs for display
		literal = escapeString(literal)
		fmt.Printf("│ %-6d │ %-19s │ %-36s │\n", tok.Pos.Line, tok.Type, literal)
	}

	fmt.Println("└────────┴─────────────────────┴──────────────────────────────────────┘")
}

func outputASTHuman(program interface{}, filePath string) {
	fmt.Printf("AST for: %s\n\n", filePath)
	fmt.Printf("%s\n", formatAST(program))
}

func outputJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func formatAST(node interface{}) string {
	// Simple string representation using JSON formatting
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting AST: %v", err)
	}
	return string(data)
}

func escapeString(s string) string {
	result := ""
	for _, ch := range s {
		switch ch {
		case '\n':
			result += "\\n"
		case '\t':
			result += "\\t"
		case '\r':
			result += "\\r"
		default:
			result += string(ch)
		}
	}
	return result
}

func handleDemo() {
	fmt.Println("=== PHP-Go Parallelization Demo ===\n")
	runParallelDemo()
}

func printUsage() {
	fmt.Printf("PHP-Go v%s\n", version)
	fmt.Println("PHP 8.4 Interpreter in Go with Automatic Parallelization")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  php-go <file>              Execute PHP file")
	fmt.Println("  php-go run <file>          Execute PHP file")
	fmt.Println("  php-go -a                  Interactive mode (coming soon)")
	fmt.Println("  php-go -S host:port        Built-in web server (coming soon)")
	fmt.Println("  php-go --version, -v       Show version")
	fmt.Println("  php-go --help, -h          Show this help")
	fmt.Println()
	fmt.Println("Development commands:")
	fmt.Println("  php-go lex [--json] <file>     Tokenize file and show tokens")
	fmt.Println("  php-go parse [--json] <file>   Parse file and show AST")
	fmt.Println("  php-go demo                    Show parallelization demo")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --json                     Output in JSON format")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  php-go test.php            Execute test.php")
	fmt.Println("  php-go run test.php        Execute test.php")
	fmt.Println("  php-go lex test.php        Show tokens from test.php")
	fmt.Println("  php-go parse test.php      Show AST from test.php")
	fmt.Println("  php-go demo                Show parallel features demo")
	fmt.Println()
}
