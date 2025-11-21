package compiler

import (
	"testing"

	"github.com/krizos/php-go/pkg/vm"
)

// ========================================
// Symbol Table Tests
// ========================================

func TestDefine(t *testing.T) {
	global := NewSymbolTable()

	// Define first variable
	a := global.Define("a")
	if a.Name != "a" {
		t.Errorf("a.Name = %q, want 'a'", a.Name)
	}
	if a.Scope != GlobalScope {
		t.Errorf("a.Scope = %v, want %v", a.Scope, GlobalScope)
	}
	if a.Index != 0 {
		t.Errorf("a.Index = %d, want 0", a.Index)
	}

	// Define second variable
	b := global.Define("b")
	if b.Name != "b" {
		t.Errorf("b.Name = %q, want 'b'", b.Name)
	}
	if b.Index != 1 {
		t.Errorf("b.Index = %d, want 1", b.Index)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()

	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table    *SymbolTable
		expected []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expected {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		{Name: "a", Scope: BuiltinScope, Index: 0},
		{Name: "c", Scope: BuiltinScope, Index: 1},
		{Name: "e", Scope: BuiltinScope, Index: 2},
		{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table               *SymbolTable
		expectedSymbols     []Symbol
		expectedFreeSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
			[]Symbol{},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
			[]Symbol{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}

		if len(tt.table.FreeSymbols()) != len(tt.expectedFreeSymbols) {
			t.Errorf("wrong number of free symbols. got=%d, want=%d",
				len(tt.table.FreeSymbols()), len(tt.expectedFreeSymbols))
			continue
		}

		for i, sym := range tt.expectedFreeSymbols {
			result := tt.table.FreeSymbols()[i]
			if result != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v", result, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 0},
		{Name: "f", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}

	expectedUnresolvable := []string{
		"b",
		"d",
	}

	for _, name := range expectedUnresolvable {
		_, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got %+v", expected.Name, expected, result)
	}
}

func TestShadowing(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("a") // Shadow global 'a'
	local.Define("c")

	// Resolve 'a' in local scope - should get local, not global
	result, ok := local.Resolve("a")
	if !ok {
		t.Fatal("could not resolve 'a'")
	}

	if result.Scope != LocalScope {
		t.Errorf("expected LocalScope, got %v", result.Scope)
	}
	if result.Index != 0 {
		t.Errorf("expected index 0, got %d", result.Index)
	}
}

func TestIsDefined(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	// Should be defined in current scope
	if !global.IsDefined("a") {
		t.Error("'a' should be defined in global scope")
	}

	// Should not be defined
	if global.IsDefined("b") {
		t.Error("'b' should not be defined in global scope")
	}

	// Test with nested scope
	local := NewEnclosedSymbolTable(global)
	local.Define("b")

	// 'b' is defined in local scope
	if !local.IsDefined("b") {
		t.Error("'b' should be defined in local scope")
	}

	// 'a' is NOT defined in local scope (it's in outer scope)
	if local.IsDefined("a") {
		t.Error("'a' should not be defined in local scope (only in outer)")
	}
}

func TestIsGlobalScope(t *testing.T) {
	global := NewSymbolTable()
	if !global.IsGlobalScope() {
		t.Error("global table should return true for IsGlobalScope()")
	}

	local := NewEnclosedSymbolTable(global)
	if local.IsGlobalScope() {
		t.Error("local table should return false for IsGlobalScope()")
	}
}

func TestNumDefinitions(t *testing.T) {
	table := NewSymbolTable()

	if table.NumDefinitions() != 0 {
		t.Errorf("expected 0 definitions, got %d", table.NumDefinitions())
	}

	table.Define("a")
	if table.NumDefinitions() != 1 {
		t.Errorf("expected 1 definition, got %d", table.NumDefinitions())
	}

	table.Define("b")
	table.Define("c")
	if table.NumDefinitions() != 3 {
		t.Errorf("expected 3 definitions, got %d", table.NumDefinitions())
	}
}

func TestOuter(t *testing.T) {
	global := NewSymbolTable()
	if global.Outer() != nil {
		t.Error("global scope should have no outer scope")
	}

	local := NewEnclosedSymbolTable(global)
	if local.Outer() != global {
		t.Error("local scope should have global as outer scope")
	}
}

// ========================================
// Compiler Integration Tests
// ========================================

func TestCompilerInitSymbolTable(t *testing.T) {
	c := New()

	// Should have builtins defined
	builtins := []string{"echo", "print", "var_dump", "isset", "empty", "count", "strlen"}

	for _, name := range builtins {
		symbol, ok := c.symbolTable.Resolve(name)
		if !ok {
			t.Errorf("builtin %q not found in symbol table", name)
			continue
		}

		if symbol.Scope != BuiltinScope {
			t.Errorf("%q should have BuiltinScope, got %v", name, symbol.Scope)
		}
	}
}

func TestCompilerEnterExitScope(t *testing.T) {
	c := New()

	// Should start in global scope
	if !c.symbolTable.IsGlobalScope() {
		t.Error("compiler should start in global scope")
	}

	// Enter a scope
	c.EnterScope()
	if c.symbolTable.IsGlobalScope() {
		t.Error("after EnterScope, should not be in global scope")
	}

	// Exit scope
	c.ExitScope()
	if !c.symbolTable.IsGlobalScope() {
		t.Error("after ExitScope, should be back in global scope")
	}
}

func TestCompilerDefineResolveVariable(t *testing.T) {
	c := New()

	// Define a variable
	sym1 := c.DefineVariable("x")
	if sym1.Name != "x" {
		t.Errorf("expected name 'x', got %q", sym1.Name)
	}
	if sym1.Scope != GlobalScope {
		t.Errorf("expected GlobalScope, got %v", sym1.Scope)
	}

	// Resolve it
	sym2, ok := c.ResolveVariable("x")
	if !ok {
		t.Fatal("could not resolve variable 'x'")
	}

	if sym2 != sym1 {
		t.Errorf("resolved symbol doesn't match defined symbol")
	}
}

func TestCompilerVariableScopes(t *testing.T) {
	c := New()

	// Define in global scope
	c.DefineVariable("global")

	// Enter local scope
	c.EnterScope()
	c.DefineVariable("local")

	// Resolve both
	globalSym, ok := c.ResolveVariable("global")
	if !ok {
		t.Fatal("could not resolve 'global'")
	}
	if globalSym.Scope != GlobalScope {
		t.Errorf("'global' should be GlobalScope, got %v", globalSym.Scope)
	}

	localSym, ok := c.ResolveVariable("local")
	if !ok {
		t.Fatal("could not resolve 'local'")
	}
	if localSym.Scope != LocalScope {
		t.Errorf("'local' should be LocalScope, got %v", localSym.Scope)
	}

	// Exit scope
	c.ExitScope()

	// 'local' should no longer be resolvable
	_, ok = c.ResolveVariable("local")
	if ok {
		t.Error("'local' should not be resolvable after exiting scope")
	}
}

func TestCompileVariableStatement(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php $x = 42; echo $x;")

	// Should have constant 42
	if len(bytecode.Constants) != 1 {
		t.Fatalf("expected 1 constant, got %d", len(bytecode.Constants))
	}

	// Should have ASSIGN and ECHO instructions
	hasAssign := false
	hasEcho := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAssign {
			hasAssign = true
		}
		if instr.Opcode == vm.OpEcho {
			hasEcho = true
		}
	}

	if !hasAssign {
		t.Error("expected ASSIGN instruction")
	}
	if !hasEcho {
		t.Error("expected ECHO instruction")
	}
}

func TestCompileMultipleVariables(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php $x = 1; $y = 2; $z = 3;")

	// Should have 3 constants
	if len(bytecode.Constants) != 3 {
		t.Fatalf("expected 3 constants, got %d", len(bytecode.Constants))
	}

	// Should have 3 ASSIGN instructions
	assignCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAssign {
			assignCount++
		}
	}

	if assignCount != 3 {
		t.Errorf("expected 3 ASSIGN instructions, got %d", assignCount)
	}
}

func TestCompileVariableArithmetic(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php $x = 1; $y = 2; $z = $x + $y;")

	// Should have constants 1 and 2
	if len(bytecode.Constants) != 2 {
		t.Fatalf("expected 2 constants, got %d", len(bytecode.Constants))
	}

	// Should have ADD instruction
	hasAdd := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
			break
		}
	}

	if !hasAdd {
		t.Error("expected ADD instruction")
	}
}
