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
		"print_r",
		"isset",
		"empty",
		// Date/Time functions
		"microtime",
		"date",
		// Array functions
		"count",
		// String functions
		"strlen",
		"substr",
		"str_replace",
		"explode",
		"implode",
		"join",
		"array_push",
		"array_pop",
		"array_map",
		"array_filter",
		"array_merge",
		// Type checking functions
		"is_null",
		"is_bool",
		"is_int",
		"is_long",
		"is_integer",
		"is_float",
		"is_double",
		"is_real",
		"is_string",
		"is_array",
		"is_object",
		"is_resource",
		"is_numeric",
		"is_scalar",
		"is_callable",
		"is_iterable",
		"is_countable",
		"gettype",
		// Math functions
		"abs",
		"ceil",
		"floor",
		"round",
		"min",
		"max",
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

// ========================================
// Namespace Resolution
// ========================================

// ResolveClassName resolves a class name to its fully qualified name (FQN)
// based on the current namespace context and use imports.
//
// Resolution rules (PHP semantics):
// 1. If name starts with \, it's already fully qualified (remove leading \)
// 2. If name is in use imports, use the imported FQN
// 3. If name is unqualified (no \), prepend current namespace
// 4. If name is qualified (has \ but doesn't start with \), prepend current namespace
func (c *Compiler) ResolveClassName(name string) string {
	if name == "" {
		return ""
	}

	// Rule 1: Fully qualified name (starts with \)
	if name[0] == '\\' {
		return name[1:] // Remove leading backslash
	}

	// Check if name contains a backslash (qualified name)
	hasBackslash := false
	for i := 0; i < len(name); i++ {
		if name[i] == '\\' {
			hasBackslash = true
			break
		}
	}

	if !hasBackslash {
		// Rule 2: Unqualified name - check use imports first
		if fqn, ok := c.useImports[name]; ok {
			return fqn
		}

		// Rule 3: Prepend current namespace (if any)
		if c.currentNamespace != "" {
			return c.currentNamespace + "\\" + name
		}
		return name
	}

	// Qualified name (has \ but doesn't start with \)
	// Get the first part before the first backslash
	firstPart := ""
	restPart := ""
	for i := 0; i < len(name); i++ {
		if name[i] == '\\' {
			firstPart = name[:i]
			restPart = name[i+1:]
			break
		}
	}

	// Check if first part is in use imports
	if fqn, ok := c.useImports[firstPart]; ok {
		return fqn + "\\" + restPart
	}

	// Rule 4: Prepend current namespace
	if c.currentNamespace != "" {
		return c.currentNamespace + "\\" + name
	}
	return name
}

// ResolveFunctionName resolves a function name to its FQN
// Functions follow different resolution rules than classes:
// 1. Check use function imports
// 2. Check current namespace
// 3. Fall back to global namespace
func (c *Compiler) ResolveFunctionName(name string) string {
	if name == "" {
		return ""
	}

	// Fully qualified name
	if name[0] == '\\' {
		return name[1:]
	}

	// Check use function imports
	if fqn, ok := c.useFunctionImports[name]; ok {
		return fqn
	}

	// For unqualified names, PHP looks in current namespace first,
	// then falls back to global namespace (for functions)
	// We return the namespaced version; runtime handles fallback
	if c.currentNamespace != "" {
		// Check if it contains backslash (qualified)
		for i := 0; i < len(name); i++ {
			if name[i] == '\\' {
				return c.currentNamespace + "\\" + name
			}
		}
		return c.currentNamespace + "\\" + name
	}
	return name
}

// ResolveConstantName resolves a constant name to its FQN
// Constants follow similar rules to functions
func (c *Compiler) ResolveConstantName(name string) string {
	if name == "" {
		return ""
	}

	// Fully qualified name
	if name[0] == '\\' {
		return name[1:]
	}

	// Check use const imports
	if fqn, ok := c.useConstImports[name]; ok {
		return fqn
	}

	// For unqualified names in a namespace, prepend namespace
	if c.currentNamespace != "" {
		return c.currentNamespace + "\\" + name
	}
	return name
}

// GetCurrentNamespace returns the current namespace context
func (c *Compiler) GetCurrentNamespace() string {
	return c.currentNamespace
}

// GetUseImports returns the current class use imports
func (c *Compiler) GetUseImports() map[string]string {
	return c.useImports
}

// GetUseFunctionImports returns the current function use imports
func (c *Compiler) GetUseFunctionImports() map[string]string {
	return c.useFunctionImports
}

// GetUseConstImports returns the current constant use imports
func (c *Compiler) GetUseConstImports() map[string]string {
	return c.useConstImports
}
