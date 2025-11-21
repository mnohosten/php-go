package compiler

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/vm"
)

// ========================================
// Helper Functions
// ========================================

// parseAndCompile parses PHP code and compiles it
func parseAndCompile(t *testing.T, input string) *Bytecode {
	l := lexer.New(input, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors:\n%v", p.Errors())
	}

	c := New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	return c.Bytecode()
}

// ========================================
// Constant Table Tests
// ========================================

func TestAddConstant(t *testing.T) {
	c := New()

	// Add first constant
	idx1 := c.AddConstant(int64(42))
	if idx1 != 0 {
		t.Errorf("First constant should have index 0, got %d", idx1)
	}

	// Add second constant
	idx2 := c.AddConstant("hello")
	if idx2 != 1 {
		t.Errorf("Second constant should have index 1, got %d", idx2)
	}

	// Add duplicate constant (should reuse index)
	idx3 := c.AddConstant(int64(42))
	if idx3 != 0 {
		t.Errorf("Duplicate constant should reuse index 0, got %d", idx3)
	}

	// Verify constant count
	if len(c.constants) != 2 {
		t.Errorf("Expected 2 constants, got %d", len(c.constants))
	}
}

func TestGetConstant(t *testing.T) {
	c := New()

	// Add constants
	c.AddConstant(int64(42))
	c.AddConstant("hello")
	c.AddConstant(true)

	// Get valid constants
	val1, err := c.GetConstant(0)
	if err != nil || val1 != int64(42) {
		t.Errorf("GetConstant(0) = %v, %v; want 42, nil", val1, err)
	}

	val2, err := c.GetConstant(1)
	if err != nil || val2 != "hello" {
		t.Errorf("GetConstant(1) = %v, %v; want 'hello', nil", val2, err)
	}

	// Get invalid index
	_, err = c.GetConstant(10)
	if err == nil {
		t.Error("Expected error for invalid constant index")
	}
}

func TestConstants(t *testing.T) {
	c := New()

	// Add constants
	c.AddConstant(int64(1))
	c.AddConstant(int64(2))
	c.AddConstant(int64(3))

	// Get constants
	constants := c.Constants()

	// Verify it's a copy
	if len(constants) != 3 {
		t.Errorf("Expected 3 constants, got %d", len(constants))
	}

	// Modify the copy (should not affect original)
	constants[0] = int64(999)
	if c.constants[0] == int64(999) {
		t.Error("Constants() should return a copy, not original slice")
	}
}

// ========================================
// Opcode Emission Tests
// ========================================

func TestEmit(t *testing.T) {
	c := New()

	pos1 := c.Emit(vm.OpNop)
	if pos1 != 0 {
		t.Errorf("First instruction should be at position 0, got %d", pos1)
	}

	pos2 := c.Emit(vm.OpAdd,
		vm.TmpVarOperand(0),
		vm.TmpVarOperand(1),
		vm.TmpVarOperand(2))
	if pos2 != 1 {
		t.Errorf("Second instruction should be at position 1, got %d", pos2)
	}

	if len(c.instructions) != 2 {
		t.Errorf("Expected 2 instructions, got %d", len(c.instructions))
	}
}

func TestEmitWithLine(t *testing.T) {
	c := New()

	c.EmitWithLine(vm.OpAdd, 42,
		vm.TmpVarOperand(0),
		vm.TmpVarOperand(1),
		vm.TmpVarOperand(2))

	instr := c.instructions[0]
	if instr.Lineno != 42 {
		t.Errorf("Expected line number 42, got %d", instr.Lineno)
	}
	if instr.Opcode != vm.OpAdd {
		t.Errorf("Expected ADD opcode, got %v", instr.Opcode)
	}
}

func TestEmitWithExtended(t *testing.T) {
	c := New()

	c.EmitWithExtended(vm.OpCast, 10, 4, // Cast to int
		vm.TmpVarOperand(0),
		vm.UnusedOperand(),
		vm.TmpVarOperand(1))

	instr := c.instructions[0]
	if instr.ExtendedValue != 4 {
		t.Errorf("Expected extended value 4, got %d", instr.ExtendedValue)
	}
	if instr.Lineno != 10 {
		t.Errorf("Expected line number 10, got %d", instr.Lineno)
	}
}

// ========================================
// Instruction Manipulation Tests
// ========================================

func TestReplaceInstruction(t *testing.T) {
	c := New()

	c.Emit(vm.OpNop)
	c.Emit(vm.OpAdd)

	// Replace first instruction
	newInstr := vm.Instruction{
		Opcode: vm.OpSub,
		Lineno: 99,
	}
	err := c.ReplaceInstruction(0, newInstr)
	if err != nil {
		t.Fatalf("ReplaceInstruction failed: %v", err)
	}

	if c.instructions[0].Opcode != vm.OpSub {
		t.Errorf("Expected SUB opcode, got %v", c.instructions[0].Opcode)
	}

	// Try invalid index
	err = c.ReplaceInstruction(10, newInstr)
	if err == nil {
		t.Error("Expected error for invalid instruction index")
	}
}

func TestChangeOperand(t *testing.T) {
	c := New()

	c.Emit(vm.OpAdd,
		vm.TmpVarOperand(0),
		vm.TmpVarOperand(1),
		vm.TmpVarOperand(2))

	// Change Op1
	err := c.ChangeOperand(0, 1, vm.CVOperand(5))
	if err != nil {
		t.Fatalf("ChangeOperand failed: %v", err)
	}

	if !c.instructions[0].Op1.IsCV() || c.instructions[0].Op1.Value != 5 {
		t.Errorf("Op1 not changed correctly: %v", c.instructions[0].Op1)
	}

	// Try invalid operand number
	err = c.ChangeOperand(0, 4, vm.CVOperand(5))
	if err == nil {
		t.Error("Expected error for invalid operand number")
	}
}

func TestLastInstructionIs(t *testing.T) {
	c := New()

	c.Emit(vm.OpAdd)
	if !c.LastInstructionIs(vm.OpAdd) {
		t.Error("LastInstructionIs(OpAdd) should be true")
	}

	c.Emit(vm.OpSub)
	if c.LastInstructionIs(vm.OpAdd) {
		t.Error("LastInstructionIs(OpAdd) should be false after emitting SUB")
	}
	if !c.LastInstructionIs(vm.OpSub) {
		t.Error("LastInstructionIs(OpSub) should be true")
	}
}

func TestRemoveLastInstruction(t *testing.T) {
	c := New()

	c.Emit(vm.OpAdd)
	c.Emit(vm.OpSub)
	c.Emit(vm.OpMul)

	if len(c.instructions) != 3 {
		t.Fatalf("Expected 3 instructions, got %d", len(c.instructions))
	}

	c.RemoveLastInstruction()

	if len(c.instructions) != 2 {
		t.Errorf("Expected 2 instructions after remove, got %d", len(c.instructions))
	}

	if !c.LastInstructionIs(vm.OpSub) {
		t.Error("Last instruction should be SUB after removing MUL")
	}
}

// ========================================
// Compilation Tests - Literals
// ========================================

func TestCompileIntegerLiteral(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php 42;")

	// Should have one constant (42)
	if len(bytecode.Constants) != 1 {
		t.Fatalf("Expected 1 constant, got %d", len(bytecode.Constants))
	}

	if bytecode.Constants[0] != int64(42) {
		t.Errorf("Expected constant 42, got %v", bytecode.Constants[0])
	}

	// Should have instructions: QM_ASSIGN, FREE
	if len(bytecode.Instructions) != 2 {
		t.Fatalf("Expected 2 instructions, got %d", len(bytecode.Instructions))
	}

	if bytecode.Instructions[0].Opcode != vm.OpQMAssign {
		t.Errorf("First instruction should be QM_ASSIGN, got %v", bytecode.Instructions[0].Opcode)
	}
}

func TestCompileStringLiteral(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php \"hello\";")

	if len(bytecode.Constants) != 1 {
		t.Fatalf("Expected 1 constant, got %d", len(bytecode.Constants))
	}

	if bytecode.Constants[0] != "hello" {
		t.Errorf("Expected constant 'hello', got %v", bytecode.Constants[0])
	}
}

func TestCompileBooleanLiteral(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php true; false;")

	if len(bytecode.Constants) != 2 {
		t.Fatalf("Expected 2 constants, got %d", len(bytecode.Constants))
	}

	if bytecode.Constants[0] != true || bytecode.Constants[1] != false {
		t.Error("Boolean constants not correct")
	}
}

func TestCompileNullLiteral(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php null;")

	if len(bytecode.Constants) != 1 {
		t.Fatalf("Expected 1 constant, got %d", len(bytecode.Constants))
	}

	if bytecode.Constants[0] != nil {
		t.Errorf("Expected constant nil, got %v", bytecode.Constants[0])
	}
}

// ========================================
// Compilation Tests - Expressions
// ========================================

func TestCompileInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator vm.Opcode
	}{
		{"<?php 1 + 2;", vm.OpAdd},
		{"<?php 1 - 2;", vm.OpSub},
		{"<?php 1 * 2;", vm.OpMul},
		{"<?php 1 / 2;", vm.OpDiv},
		{"<?php 1 % 2;", vm.OpMod},
		{"<?php 1 ** 2;", vm.OpPow},
		{"<?php \"a\" . \"b\";", vm.OpConcat},
		{"<?php 1 == 2;", vm.OpIsEqual},
		{"<?php 1 != 2;", vm.OpIsNotEqual},
		{"<?php 1 === 2;", vm.OpIsIdentical},
		{"<?php 1 !== 2;", vm.OpIsNotIdentical},
		{"<?php 1 < 2;", vm.OpIsSmaller},
		{"<?php 1 <= 2;", vm.OpIsSmallerOrEqual},
		{"<?php 1 > 2;", vm.OpIsSmaller}, // Swaps operands
		{"<?php 1 >= 2;", vm.OpIsSmallerOrEqual}, // Swaps operands
		{"<?php 1 | 2;", vm.OpBWOr},
		{"<?php 1 & 2;", vm.OpBWAnd},
		{"<?php 1 ^ 2;", vm.OpBWXor},
		{"<?php 1 << 2;", vm.OpSL},
		{"<?php 1 >> 2;", vm.OpSR},
		{"<?php 1 <=> 2;", vm.OpSpaceship},
	}

	for _, tt := range tests {
		bytecode := parseAndCompile(t, tt.input)

		// Find the operator instruction (should be the third instruction: 2 QM_ASSIGNs, then operator)
		found := false
		for _, instr := range bytecode.Instructions {
			if instr.Opcode == tt.operator {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Input %q: expected opcode %v not found in instructions", tt.input, tt.operator)
		}
	}
}

func TestCompilePrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator vm.Opcode
	}{
		{"<?php !true;", vm.OpBoolNot},
		{"<?php -5;", vm.OpSub}, // Unary minus becomes 0 - x
		{"<?php ~5;", vm.OpBWNot},
	}

	for _, tt := range tests {
		bytecode := parseAndCompile(t, tt.input)

		// Find the operator instruction
		found := false
		for _, instr := range bytecode.Instructions {
			if instr.Opcode == tt.operator {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Input %q: expected opcode %v not found", tt.input, tt.operator)
		}
	}
}

// ========================================
// Compilation Tests - Statements
// ========================================

func TestCompileEchoStatement(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php echo \"hello\";")

	// Should have ECHO instruction
	found := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpEcho {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected ECHO instruction not found")
	}
}

func TestCompileMultipleEcho(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php echo 1, 2, 3;")

	// Should have 3 ECHO instructions
	echoCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpEcho {
			echoCount++
		}
	}

	if echoCount != 3 {
		t.Errorf("Expected 3 ECHO instructions, got %d", echoCount)
	}
}

func TestCompileReturnStatement(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php return 42;")

	// Should have RETURN instruction
	found := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected RETURN instruction not found")
	}
}

func TestCompileEmptyReturn(t *testing.T) {
	bytecode := parseAndCompile(t, "<?php return;")

	// Should have RETURN instruction with no operand
	found := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected RETURN instruction not found")
	}
}

// ========================================
// Program Assembly Tests
// ========================================

func TestBytecode(t *testing.T) {
	c := New()

	c.AddConstant(int64(42))
	c.Emit(vm.OpAdd)

	bytecode := c.Bytecode()

	if len(bytecode.Instructions) != 1 {
		t.Errorf("Expected 1 instruction, got %d", len(bytecode.Instructions))
	}

	if len(bytecode.Constants) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(bytecode.Constants))
	}
}

// ========================================
// Reset Tests
// ========================================

func TestReset(t *testing.T) {
	c := New()

	// Add some data
	c.AddConstant(int64(42))
	c.Emit(vm.OpAdd)

	// Reset
	c.Reset()

	// Verify everything is cleared
	if len(c.instructions) != 0 {
		t.Errorf("Expected 0 instructions after reset, got %d", len(c.instructions))
	}

	if len(c.constants) != 0 {
		t.Errorf("Expected 0 constants after reset, got %d", len(c.constants))
	}

	if len(c.constantMap) != 0 {
		t.Errorf("Expected empty constant map after reset, got %d entries", len(c.constantMap))
	}
}

// ========================================
// Real-World Example Tests
// ========================================

func TestCompileSimpleProgram(t *testing.T) {
	input := `<?php
	echo "Hello, World!";
	return 0;
	`

	bytecode := parseAndCompile(t, input)

	// Should have both ECHO and RETURN
	hasEcho := false
	hasReturn := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpEcho {
			hasEcho = true
		}
		if instr.Opcode == vm.OpReturn {
			hasReturn = true
		}
	}

	if !hasEcho {
		t.Error("Expected ECHO instruction")
	}
	if !hasReturn {
		t.Error("Expected RETURN instruction")
	}
}

func TestCompileArithmeticExpression(t *testing.T) {
	input := `<?php
	2 + 3 * 4;
	`

	bytecode := parseAndCompile(t, input)

	// Should have MUL and ADD opcodes
	hasMul := false
	hasAdd := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpMul {
			hasMul = true
		}
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
		}
	}

	if !hasMul {
		t.Error("Expected MUL instruction")
	}
	if !hasAdd {
		t.Error("Expected ADD instruction")
	}

	// Should have constants 2, 3, 4
	if len(bytecode.Constants) != 3 {
		t.Errorf("Expected 3 constants, got %d", len(bytecode.Constants))
	}
}
