package compiler

import "fmt"

// SymbolScope represents the scope of a symbol
type SymbolScope string

const (
	// GlobalScope - global variables (accessible everywhere with 'global' keyword)
	GlobalScope SymbolScope = "GLOBAL"

	// LocalScope - function-local variables
	LocalScope SymbolScope = "LOCAL"

	// BuiltinScope - built-in functions and constants
	BuiltinScope SymbolScope = "BUILTIN"

	// FreeScope - free variables (closure variables)
	FreeScope SymbolScope = "FREE"
)

// Symbol represents a variable or function in the symbol table
type Symbol struct {
	// Name of the symbol (variable name without $)
	Name string

	// Scope of the symbol (GLOBAL, LOCAL, BUILTIN, FREE)
	Scope SymbolScope

	// Index in the compiled variable array (for CV operands)
	Index int
}

// SymbolTable manages symbols in a scope
type SymbolTable struct {
	// outer is the parent scope (nil for global scope)
	outer *SymbolTable

	// store maps symbol names to Symbol structs
	store map[string]Symbol

	// numDefinitions tracks the number of symbols defined in this scope
	numDefinitions int

	// freeSymbols tracks free variables (for closures)
	freeSymbols []Symbol
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:       make(map[string]Symbol),
		freeSymbols: []Symbol{},
	}
}

// NewEnclosedSymbolTable creates a new symbol table with an outer scope
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.outer = outer
	return s
}

// Define adds a new symbol to the symbol table
// Returns the defined symbol
func (s *SymbolTable) Define(name string) Symbol {
	// Determine scope based on whether we have an outer scope
	scope := GlobalScope
	if s.outer != nil {
		scope = LocalScope
	}

	symbol := Symbol{
		Name:  name,
		Scope: scope,
		Index: s.numDefinitions,
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

// DefineBuiltin adds a built-in symbol (function or constant)
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: BuiltinScope,
		Index: index,
	}
	s.store[name] = symbol
	return symbol
}

// DefineFree adds a free variable (closure variable)
func (s *SymbolTable) DefineFree(original Symbol) Symbol {
	s.freeSymbols = append(s.freeSymbols, original)

	symbol := Symbol{
		Name:  original.Name,
		Scope: FreeScope,
		Index: len(s.freeSymbols) - 1,
	}

	s.store[symbol.Name] = symbol
	return symbol
}

// Resolve looks up a symbol by name
// Searches current scope and outer scopes
// Returns the symbol and true if found, empty symbol and false otherwise
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	// Look in current scope
	symbol, ok := s.store[name]
	if ok {
		return symbol, true
	}

	// If not found and we have an outer scope, search there
	if s.outer != nil {
		symbol, ok := s.outer.Resolve(name)
		if !ok {
			return symbol, false
		}

		// If found in outer scope and it's not global or builtin,
		// make it a free variable in current scope (for closures)
		if symbol.Scope == GlobalScope || symbol.Scope == BuiltinScope {
			return symbol, true
		}

		// Make it a free variable
		free := s.DefineFree(symbol)
		return free, true
	}

	return Symbol{}, false
}

// IsDefined checks if a symbol is defined in the current scope (not outer scopes)
func (s *SymbolTable) IsDefined(name string) bool {
	_, ok := s.store[name]
	return ok
}

// NumDefinitions returns the number of symbols defined in this scope
func (s *SymbolTable) NumDefinitions() int {
	return s.numDefinitions
}

// FreeSymbols returns the free variables in this scope
func (s *SymbolTable) FreeSymbols() []Symbol {
	return s.freeSymbols
}

// Outer returns the parent scope
func (s *SymbolTable) Outer() *SymbolTable {
	return s.outer
}

// IsGlobalScope returns true if this is the global scope
func (s *SymbolTable) IsGlobalScope() bool {
	return s.outer == nil
}

// ========================================
// Helper Methods
// ========================================

// String returns a string representation of the symbol table for debugging
func (s *SymbolTable) String() string {
	result := fmt.Sprintf("SymbolTable (definitions=%d):\n", s.numDefinitions)
	for name, symbol := range s.store {
		result += fmt.Sprintf("  %s: %s[%d]\n", name, symbol.Scope, symbol.Index)
	}
	if len(s.freeSymbols) > 0 {
		result += "  Free variables:\n"
		for i, symbol := range s.freeSymbols {
			result += fmt.Sprintf("    [%d] %s\n", i, symbol.Name)
		}
	}
	if s.outer != nil {
		result += "  (has outer scope)\n"
	}
	return result
}

// String returns a string representation of a symbol
func (sym Symbol) String() string {
	return fmt.Sprintf("%s:%s[%d]", sym.Name, sym.Scope, sym.Index)
}

// ========================================
// Compiler Integration
// ========================================

// InitSymbolTable initializes the compiler's symbol table with built-ins
func (c *Compiler) InitSymbolTable() {
	c.symbolTable = NewSymbolTable()

	// Define built-in functions
	// These will be expanded as we implement the standard library
	builtins := []string{
		"echo",
		"print",
		"var_dump",
		"isset",
		"empty",
		"count",
		"strlen",
		// More built-ins will be added in Phase 6
	}

	for i, name := range builtins {
		c.symbolTable.DefineBuiltin(i, name)
	}
}

// EnterScope creates a new nested scope (for functions, blocks, etc.)
func (c *Compiler) EnterScope() {
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

// ExitScope returns to the parent scope
func (c *Compiler) ExitScope() {
	if c.symbolTable.outer != nil {
		c.symbolTable = c.symbolTable.outer
	}
}

// DefineVariable defines a new variable in the current scope
func (c *Compiler) DefineVariable(name string) Symbol {
	return c.symbolTable.Define(name)
}

// ResolveVariable looks up a variable by name
func (c *Compiler) ResolveVariable(name string) (Symbol, bool) {
	return c.symbolTable.Resolve(name)
}

// IsVariableDefined checks if a variable is defined in the current scope
func (c *Compiler) IsVariableDefined(name string) bool {
	return c.symbolTable.IsDefined(name)
}
