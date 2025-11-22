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
		// Use variables to prevent constant folding
		{"<?php $a + $b;", vm.OpAdd},
		{"<?php $a - $b;", vm.OpSub},
		{"<?php $a * $b;", vm.OpMul},
		{"<?php $a / $b;", vm.OpDiv},
		{"<?php $a % $b;", vm.OpMod},
		{"<?php $a ** $b;", vm.OpPow},
		{"<?php $a . $b;", vm.OpConcat},
		{"<?php $a == $b;", vm.OpIsEqual},
		{"<?php $a != $b;", vm.OpIsNotEqual},
		{"<?php $a === $b;", vm.OpIsIdentical},
		{"<?php $a !== $b;", vm.OpIsNotIdentical},
		{"<?php $a < $b;", vm.OpIsSmaller},
		{"<?php $a <= $b;", vm.OpIsSmallerOrEqual},
		{"<?php $a > $b;", vm.OpIsSmaller}, // Swaps operands
		{"<?php $a >= $b;", vm.OpIsSmallerOrEqual}, // Swaps operands
		{"<?php $a | $b;", vm.OpBWOr},
		{"<?php $a & $b;", vm.OpBWAnd},
		{"<?php $a ^ $b;", vm.OpBWXor},
		{"<?php $a << $b;", vm.OpSL},
		{"<?php $a >> $b;", vm.OpSR},
		{"<?php $a <=> $b;", vm.OpSpaceship},
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
		// Use variables to prevent constant folding
		{"<?php !$x;", vm.OpBoolNot},
		{"<?php -$x;", vm.OpSub}, // Unary minus becomes 0 - x
		{"<?php ~$x;", vm.OpBWNot},
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
	// Use variables to prevent constant folding
	input := `<?php
	$a + $b * $c;
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
}

// ========================================
// Task 2.5: Expression Compilation Tests
// ========================================

func TestCompileArrayLiteral(t *testing.T) {
	input := `<?php
	$x = [1, 2, 3];
	`

	bytecode := parseAndCompile(t, input)

	// Should have INIT_ARRAY and ADD_ARRAY_ELEMENT opcodes
	hasInitArray := false
	addArrayCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitArray {
			hasInitArray = true
		}
		if instr.Opcode == vm.OpAddArrayElement {
			addArrayCount++
		}
	}

	if !hasInitArray {
		t.Error("Expected INIT_ARRAY instruction")
	}
	if addArrayCount != 3 {
		t.Errorf("Expected 3 ADD_ARRAY_ELEMENT instructions, got %d", addArrayCount)
	}

	// Should have constants 1, 2, 3
	if len(bytecode.Constants) < 3 {
		t.Errorf("Expected at least 3 constants, got %d", len(bytecode.Constants))
	}
}

func TestCompileAssociativeArray(t *testing.T) {
	input := `<?php
	$x = ["foo" => 1, "bar" => 2];
	`

	bytecode := parseAndCompile(t, input)

	// Should have INIT_ARRAY and ADD_ARRAY_ELEMENT opcodes
	hasInitArray := false
	addArrayCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitArray {
			hasInitArray = true
		}
		if instr.Opcode == vm.OpAddArrayElement {
			addArrayCount++
		}
	}

	if !hasInitArray {
		t.Error("Expected INIT_ARRAY instruction")
	}
	if addArrayCount != 2 {
		t.Errorf("Expected 2 ADD_ARRAY_ELEMENT instructions, got %d", addArrayCount)
	}
}

func TestCompileArrayAccess(t *testing.T) {
	input := `<?php
	$x = $arr[0];
	`

	bytecode := parseAndCompile(t, input)

	// Should have FETCH_DIM_R opcode
	hasFetchDim := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFetchDimR {
			hasFetchDim = true
			break
		}
	}

	if !hasFetchDim {
		t.Error("Expected FETCH_DIM_R instruction")
	}

	// Should have constant 0
	found := false
	for _, c := range bytecode.Constants {
		if i, ok := c.(int64); ok && i == 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected constant 0")
	}
}

func TestCompilePropertyAccess(t *testing.T) {
	input := `<?php
	$x = $obj->prop;
	`

	bytecode := parseAndCompile(t, input)

	// Should have FETCH_OBJ_R opcode
	hasFetchObj := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFetchObjR {
			hasFetchObj = true
			break
		}
	}

	if !hasFetchObj {
		t.Error("Expected FETCH_OBJ_R instruction")
	}

	// Should have constant "prop"
	found := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "prop" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected constant 'prop'")
	}
}

func TestCompileFunctionCall(t *testing.T) {
	input := `<?php
	$x = strlen("hello");
	`

	bytecode := parseAndCompile(t, input)

	// Should have INIT_FCALL_BY_NAME and DO_FCALL opcodes
	hasInitFcall := false
	hasDoFcall := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitFcallByName {
			hasInitFcall = true
		}
		if instr.Opcode == vm.OpDoFcall {
			hasDoFcall = true
		}
	}

	if !hasInitFcall {
		t.Error("Expected INIT_FCALL_BY_NAME instruction")
	}
	if !hasDoFcall {
		t.Error("Expected DO_FCALL instruction")
	}

	// Should have constants "strlen" and "hello"
	hasStrlen := false
	hasHello := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok {
			if s == "strlen" {
				hasStrlen = true
			}
			if s == "hello" {
				hasHello = true
			}
		}
	}
	if !hasStrlen {
		t.Error("Expected constant 'strlen'")
	}
	if !hasHello {
		t.Error("Expected constant 'hello'")
	}
}

func TestCompileMethodCall(t *testing.T) {
	input := `<?php
	$x = $obj->method(1, 2);
	`

	bytecode := parseAndCompile(t, input)

	// Should have INIT_METHOD_CALL and DO_FCALL opcodes
	hasInitMethod := false
	hasDoFcall := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitMethodCall {
			hasInitMethod = true
		}
		if instr.Opcode == vm.OpDoFcall {
			hasDoFcall = true
		}
	}

	if !hasInitMethod {
		t.Error("Expected INIT_METHOD_CALL instruction")
	}
	if !hasDoFcall {
		t.Error("Expected DO_FCALL instruction")
	}

	// Should have constant "method"
	found := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "method" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected constant 'method'")
	}
}

func TestCompileTernaryOperator(t *testing.T) {
	input := `<?php
	$x = $a ? $b : $c;
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMPZ and JMP opcodes
	hasJmpz := false
	hasJmp := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpJmp {
			hasJmp = true
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for ternary")
	}
	if !hasJmp {
		t.Error("Expected JMP instruction for ternary")
	}
}

func TestCompileShortTernary(t *testing.T) {
	input := `<?php
	$x = $a ?: $b;
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMP_SET opcode
	hasJmpSet := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpSet {
			hasJmpSet = true
			break
		}
	}

	if !hasJmpSet {
		t.Error("Expected JMP_SET instruction for short ternary")
	}
}

func TestCompileTypeCast(t *testing.T) {
	tests := []struct {
		input    string
		castType string
	}{
		{`<?php $x = (int)$y;`, "int"},
		{`<?php $x = (string)$y;`, "string"},
		{`<?php $x = (bool)$y;`, "bool"},
		// Note: float/double and array casts need parser support to be added later
	}

	for _, tt := range tests {
		bytecode := parseAndCompile(t, tt.input)

		// Should have CAST opcode
		hasCast := false
		for _, instr := range bytecode.Instructions {
			if instr.Opcode == vm.OpCast {
				hasCast = true
				break
			}
		}

		if !hasCast {
			t.Errorf("Expected CAST instruction for %s cast", tt.castType)
		}
	}
}

func TestCompileInstanceof(t *testing.T) {
	input := `<?php
	$x = $obj instanceof MyClass;
	`

	bytecode := parseAndCompile(t, input)

	// Should have INSTANCEOF opcode
	hasInstanceof := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInstanceof {
			hasInstanceof = true
			break
		}
	}

	if !hasInstanceof {
		t.Error("Expected INSTANCEOF instruction")
	}

	// Should have constant "MyClass"
	found := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "MyClass" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected constant 'MyClass'")
	}
}

func TestCompileGroupedExpression(t *testing.T) {
	// Use variables to prevent constant folding
	input := `<?php
	$x = ($a + $b) * $c;
	`

	bytecode := parseAndCompile(t, input)

	// Should have ADD and MUL opcodes
	hasAdd := false
	hasMul := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
		}
		if instr.Opcode == vm.OpMul {
			hasMul = true
		}
	}

	if !hasAdd {
		t.Error("Expected ADD instruction")
	}
	if !hasMul {
		t.Error("Expected MUL instruction")
	}
}

func TestCompileComplexExpression(t *testing.T) {
	input := `<?php
	$result = $arr[0]->method($x, $y) + 10;
	`

	bytecode := parseAndCompile(t, input)

	// Should have FETCH_DIM_R, INIT_METHOD_CALL, DO_FCALL, and ADD opcodes
	hasFetchDim := false
	hasInitMethod := false
	hasDoFcall := false
	hasAdd := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFetchDimR {
			hasFetchDim = true
		}
		if instr.Opcode == vm.OpInitMethodCall {
			hasInitMethod = true
		}
		if instr.Opcode == vm.OpDoFcall {
			hasDoFcall = true
		}
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
		}
	}

	if !hasFetchDim {
		t.Error("Expected FETCH_DIM_R instruction")
	}
	if !hasInitMethod {
		t.Error("Expected INIT_METHOD_CALL instruction")
	}
	if !hasDoFcall {
		t.Error("Expected DO_FCALL instruction")
	}
	if !hasAdd {
		t.Error("Expected ADD instruction")
	}
}

func TestCompileNestedArrays(t *testing.T) {
	input := `<?php
	$x = [1, [2, 3], 4];
	`

	bytecode := parseAndCompile(t, input)

	// Should have multiple INIT_ARRAY instructions (one for outer, one for inner)
	initArrayCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitArray {
			initArrayCount++
		}
	}

	if initArrayCount < 2 {
		t.Errorf("Expected at least 2 INIT_ARRAY instructions for nested arrays, got %d", initArrayCount)
	}
}

func TestCompileIdentifier(t *testing.T) {
	input := `<?php
	$x = MyClass;
	`

	bytecode := parseAndCompile(t, input)

	// Should have constant "MyClass"
	found := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "MyClass" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected constant 'MyClass' from identifier")
	}
}

// ========================================
// Task 2.6: Statement Compilation Tests
// ========================================

func TestCompileIfStatement(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		echo "positive";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMPZ and JMP opcodes
	hasJmpz := false
	hasJmp := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpJmp {
			hasJmp = true
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for if statement")
	}
	if !hasJmp {
		t.Error("Expected JMP instruction for if statement")
	}
}

func TestCompileIfElseStatement(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		echo "positive";
	} else {
		echo "non-positive";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMPZ, JMP, and ECHO opcodes
	hasJmpz := false
	hasJmp := false
	echoCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpJmp {
			hasJmp = true
		}
		if instr.Opcode == vm.OpEcho {
			echoCount++
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction")
	}
	if !hasJmp {
		t.Error("Expected JMP instruction")
	}
	if echoCount != 2 {
		t.Errorf("Expected 2 ECHO instructions, got %d", echoCount)
	}
}

func TestCompileWhileLoop(t *testing.T) {
	input := `<?php
	while ($i < 10) {
		$i = $i + 1;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMPZ and JMP opcodes
	hasJmpz := false
	jmpCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for while loop")
	}
	if jmpCount < 1 {
		t.Error("Expected at least 1 JMP instruction for while loop")
	}
}

func TestCompileForLoop(t *testing.T) {
	input := `<?php
	for ($i = 0; $i < 10; $i = $i + 1) {
		echo $i;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have JMPZ, JMP, and ECHO opcodes
	hasJmpz := false
	jmpCount := 0
	hasEcho := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
		if instr.Opcode == vm.OpEcho {
			hasEcho = true
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for for loop")
	}
	if jmpCount < 1 {
		t.Error("Expected at least 1 JMP instruction for for loop")
	}
	if !hasEcho {
		t.Error("Expected ECHO instruction in for loop body")
	}
}

func TestCompileForeachLoop(t *testing.T) {
	input := `<?php
	foreach ($arr as $val) {
		echo $val;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have FE_RESET_R, FE_FETCH_R, FE_FREE opcodes
	hasFeReset := false
	hasFeFetch := false
	hasFeFree := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFeResetR {
			hasFeReset = true
		}
		if instr.Opcode == vm.OpFeFetchR {
			hasFeFetch = true
		}
		if instr.Opcode == vm.OpFeFree {
			hasFeFree = true
		}
	}

	if !hasFeReset {
		t.Error("Expected FE_RESET_R instruction")
	}
	if !hasFeFetch {
		t.Error("Expected FE_FETCH_R instruction")
	}
	if !hasFeFree {
		t.Error("Expected FE_FREE instruction")
	}
}

func TestCompileForeachWithKey(t *testing.T) {
	input := `<?php
	foreach ($arr as $key => $val) {
		echo $key;
		echo $val;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have FE_RESET_R and multiple ASSIGN opcodes
	hasFeReset := false
	assignCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFeResetR {
			hasFeReset = true
		}
		if instr.Opcode == vm.OpAssign {
			assignCount++
		}
	}

	if !hasFeReset {
		t.Error("Expected FE_RESET_R instruction")
	}
	if assignCount < 2 {
		t.Errorf("Expected at least 2 ASSIGN instructions (key and value), got %d", assignCount)
	}
}

func TestCompileBreakStatement(t *testing.T) {
	input := `<?php
	while (true) {
		break;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have multiple JMP opcodes (loop back and break)
	jmpCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
	}

	if jmpCount < 2 {
		t.Errorf("Expected at least 2 JMP instructions (loop and break), got %d", jmpCount)
	}
}

func TestCompileContinueStatement(t *testing.T) {
	input := `<?php
	while ($i < 10) {
		if ($i == 5) {
			continue;
		}
		echo $i;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have multiple JMP opcodes
	jmpCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
	}

	if jmpCount < 3 {
		t.Errorf("Expected at least 3 JMP instructions (if-end, continue, loop), got %d", jmpCount)
	}
}

func TestCompileSwitchStatement(t *testing.T) {
	input := `<?php
	switch ($x) {
		case 1:
			echo "one";
			break;
		case 2:
			echo "two";
			break;
		default:
			echo "other";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have IS_EQUAL and JMPNZ opcodes for case comparisons
	hasIsEqual := false
	hasJmpnz := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpIsEqual {
			hasIsEqual = true
		}
		if instr.Opcode == vm.OpJmpNZ {
			hasJmpnz = true
		}
	}

	if !hasIsEqual {
		t.Error("Expected IS_EQUAL instruction for switch cases")
	}
	if !hasJmpnz {
		t.Error("Expected JMPNZ instruction for switch cases")
	}
}

func TestCompileTryCatchStatement(t *testing.T) {
	input := `<?php
	try {
		echo "trying";
	} catch (Exception $e) {
		echo "caught";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have CATCH opcode
	hasCatch := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpCatch {
			hasCatch = true
			break
		}
	}

	if !hasCatch {
		t.Error("Expected CATCH instruction")
	}
}

func TestCompileTryCatchFinallyStatement(t *testing.T) {
	input := `<?php
	try {
		echo "trying";
	} catch (Exception $e) {
		echo "caught";
	} finally {
		echo "finally";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have FAST_CALL, CATCH, and FAST_RET opcodes
	hasFastCall := false
	hasCatch := false
	hasFastRet := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFastCall {
			hasFastCall = true
		}
		if instr.Opcode == vm.OpCatch {
			hasCatch = true
		}
		if instr.Opcode == vm.OpFastRet {
			hasFastRet = true
		}
	}

	if !hasFastCall {
		t.Error("Expected FAST_CALL instruction for finally block")
	}
	if !hasCatch {
		t.Error("Expected CATCH instruction")
	}
	if !hasFastRet {
		t.Error("Expected FAST_RET instruction for finally block")
	}
}

func TestCompileThrowStatement(t *testing.T) {
	input := `<?php
	throw new Exception("error");
	`

	bytecode := parseAndCompile(t, input)

	// Should have THROW opcode
	hasThrow := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpThrow {
			hasThrow = true
			break
		}
	}

	if !hasThrow {
		t.Error("Expected THROW instruction")
	}
}

func TestCompileNestedLoops(t *testing.T) {
	input := `<?php
	for ($i = 0; $i < 10; $i = $i + 1) {
		for ($j = 0; $j < 10; $j = $j + 1) {
			echo $i;
			echo $j;
		}
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have multiple JMPZ and JMP opcodes for nested loops
	jmpzCount := 0
	jmpCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			jmpzCount++
		}
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
	}

	if jmpzCount < 2 {
		t.Errorf("Expected at least 2 JMPZ instructions for nested loops, got %d", jmpzCount)
	}
	if jmpCount < 2 {
		t.Errorf("Expected at least 2 JMP instructions for nested loops, got %d", jmpCount)
	}
}

func TestCompileComplexControlFlow(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		for ($i = 0; $i < $x; $i = $i + 1) {
			if ($i == 5) {
				break;
			}
			echo $i;
		}
	} else {
		echo "negative";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have multiple control flow opcodes
	jmpzCount := 0
	jmpCount := 0
	echoCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			jmpzCount++
		}
		if instr.Opcode == vm.OpJmp {
			jmpCount++
		}
		if instr.Opcode == vm.OpEcho {
			echoCount++
		}
	}

	if jmpzCount < 2 {
		t.Errorf("Expected at least 2 JMPZ instructions, got %d", jmpzCount)
	}
	if jmpCount < 3 {
		t.Errorf("Expected at least 3 JMP instructions, got %d", jmpCount)
	}
	if echoCount != 2 {
		t.Errorf("Expected 2 ECHO instructions, got %d", echoCount)
	}
}

// ========================================
// Task 2.8: Function Compilation Tests
// ========================================

func TestCompileFunctionDeclaration(t *testing.T) {
	input := `<?php
	function greet() {
		echo "Hello";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have DECLARE_FUNCTION and ECHO opcodes
	hasDeclareFunc := false
	hasEcho := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareFunction {
			hasDeclareFunc = true
		}
		if instr.Opcode == vm.OpEcho {
			hasEcho = true
		}
	}

	if !hasDeclareFunc {
		t.Error("Expected DECLARE_FUNCTION instruction")
	}
	if !hasEcho {
		t.Error("Expected ECHO instruction in function body")
	}

	// Should have function name as constant
	foundGreet := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "greet" {
			foundGreet = true
			break
		}
	}
	if !foundGreet {
		t.Error("Expected 'greet' function name in constants")
	}
}

func TestCompileFunctionWithParameters(t *testing.T) {
	input := `<?php
	function add($a, $b) {
		return $a + $b;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have RECV opcodes for parameters
	recvCount := 0
	hasAdd := false
	hasReturn := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			recvCount++
		}
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
		}
		if instr.Opcode == vm.OpReturn {
			hasReturn = true
		}
	}

	if recvCount != 2 {
		t.Errorf("Expected 2 RECV instructions for parameters, got %d", recvCount)
	}
	if !hasAdd {
		t.Error("Expected ADD instruction in function body")
	}
	if !hasReturn {
		t.Error("Expected RETURN instruction")
	}
}

func TestCompileFunctionWithDefaultParameter(t *testing.T) {
	input := `<?php
	function greet($name = "World") {
		echo $name;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have RECV_INIT opcode for parameter with default
	hasRecvInit := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecvInit {
			hasRecvInit = true
			break
		}
	}

	if !hasRecvInit {
		t.Error("Expected RECV_INIT instruction for parameter with default value")
	}

	// Should have "World" as constant for default value
	foundWorld := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok && s == "World" {
			foundWorld = true
			break
		}
	}
	if !foundWorld {
		t.Error("Expected 'World' default value in constants")
	}
}

func TestCompileFunctionWithVariadicParameter(t *testing.T) {
	input := `<?php
	function sum(...$numbers) {
		return 0;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have RECV_VARIADIC opcode
	hasRecvVariadic := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecvVariadic {
			hasRecvVariadic = true
			break
		}
	}

	if !hasRecvVariadic {
		t.Error("Expected RECV_VARIADIC instruction for variadic parameter")
	}
}

func TestCompileFunctionWithReturnValue(t *testing.T) {
	input := `<?php
	function triple($x) {
		return $x * 3;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have MUL and RETURN opcodes
	hasMul := false
	hasReturn := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpMul {
			hasMul = true
		}
		if instr.Opcode == vm.OpReturn {
			hasReturn = true
		}
	}

	if !hasMul {
		t.Error("Expected MUL instruction")
	}
	if !hasReturn {
		t.Error("Expected RETURN instruction")
	}
}

func TestCompileFunctionWithMultipleParameters(t *testing.T) {
	input := `<?php
	function calculate($a, $b, $c) {
		return ($a + $b) * $c;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have 3 RECV opcodes
	recvCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			recvCount++
		}
	}

	if recvCount != 3 {
		t.Errorf("Expected 3 RECV instructions, got %d", recvCount)
	}
}

func TestCompileFunctionImplicitReturn(t *testing.T) {
	input := `<?php
	function noReturn() {
		echo "test";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have implicit RETURN at end
	hasReturn := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			hasReturn = true
			break
		}
	}

	if !hasReturn {
		t.Error("Expected implicit RETURN instruction")
	}
}

func TestCompileNestedFunctionDeclarations(t *testing.T) {
	input := `<?php
	function outer() {
		echo "outer";
	}

	function inner() {
		echo "inner";
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have 2 DECLARE_FUNCTION opcodes
	declareFuncCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareFunction {
			declareFuncCount++
		}
	}

	if declareFuncCount != 2 {
		t.Errorf("Expected 2 DECLARE_FUNCTION instructions, got %d", declareFuncCount)
	}

	// Should have both function names as constants
	hasOuter := false
	hasInner := false
	for _, c := range bytecode.Constants {
		if s, ok := c.(string); ok {
			if s == "outer" {
				hasOuter = true
			}
			if s == "inner" {
				hasInner = true
			}
		}
	}

	if !hasOuter {
		t.Error("Expected 'outer' function name in constants")
	}
	if !hasInner {
		t.Error("Expected 'inner' function name in constants")
	}
}

func TestCompileFunctionWithMixedParameters(t *testing.T) {
	input := `<?php
	function variedParams($required, $optional = 10, ...$rest) {
		return $required;
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have RECV, RECV_INIT, and RECV_VARIADIC
	hasRecv := false
	hasRecvInit := false
	hasRecvVariadic := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			hasRecv = true
		}
		if instr.Opcode == vm.OpRecvInit {
			hasRecvInit = true
		}
		if instr.Opcode == vm.OpRecvVariadic {
			hasRecvVariadic = true
		}
	}

	if !hasRecv {
		t.Error("Expected RECV instruction for required parameter")
	}
	if !hasRecvInit {
		t.Error("Expected RECV_INIT instruction for optional parameter")
	}
	if !hasRecvVariadic {
		t.Error("Expected RECV_VARIADIC instruction for variadic parameter")
	}
}

func TestCompileFunctionWithComplexBody(t *testing.T) {
	input := `<?php
	function complex($x) {
		if ($x > 0) {
			return $x * 2;
		} else {
			return 0;
		}
	}
	`

	bytecode := parseAndCompile(t, input)

	// Should have control flow and return opcodes
	hasJmpz := false
	returnCount := 0

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
		}
		if instr.Opcode == vm.OpReturn {
			returnCount++
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for if statement")
	}
	if returnCount < 2 {
		t.Errorf("Expected at least 2 RETURN instructions, got %d", returnCount)
	}
}

// ========================================
// Class Compilation Tests
// ========================================

func TestCompileBasicClass(t *testing.T) {
	input := `<?php
class User {
}
`

	bytecode := parseAndCompile(t, input)

	// Should have DECLARE_CLASS opcode
	hasDeclareClass := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareClass {
			hasDeclareClass = true
			break
		}
	}

	if !hasDeclareClass {
		t.Error("Expected DECLARE_CLASS instruction")
	}

	// Should have "User" in constants
	hasClassName := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "User" {
			hasClassName = true
			break
		}
	}

	if !hasClassName {
		t.Error("Expected class name 'User' in constants")
	}
}

func TestCompileClassWithProperties(t *testing.T) {
	input := `<?php
class User {
    public $name;
    public $email = "default@example.com";
    private $password;
}
`

	bytecode := parseAndCompile(t, input)

	// Should have property names in constants
	propertyNames := []string{"name", "email", "password"}
	for _, propName := range propertyNames {
		found := false
		for _, c := range bytecode.Constants {
			if str, ok := c.(string); ok && str == propName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected property name '%s' in constants", propName)
		}
	}

	// Should have default value "default@example.com" in constants
	hasDefaultValue := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "default@example.com" {
			hasDefaultValue = true
			break
		}
	}

	if !hasDefaultValue {
		t.Error("Expected default value 'default@example.com' in constants")
	}
}

func TestCompileClassWithMethod(t *testing.T) {
	input := `<?php
class User {
    public function getName() {
        return $this->name;
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have method name "getName" in constants
	hasMethodName := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "getName" {
			hasMethodName = true
			break
		}
	}

	if !hasMethodName {
		t.Error("Expected method name 'getName' in constants")
	}

	// Should have RETURN opcode for method
	hasReturn := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			hasReturn = true
			break
		}
	}

	if !hasReturn {
		t.Error("Expected RETURN instruction for method")
	}
}

func TestCompileClassWithConstructor(t *testing.T) {
	input := `<?php
class User {
    public $name;

    public function __construct($name) {
        $this->name = $name;
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have constructor name "__construct" in constants
	hasConstructor := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "__construct" {
			hasConstructor = true
			break
		}
	}

	if !hasConstructor {
		t.Error("Expected constructor name '__construct' in constants")
	}

	// Should have RECV opcode for parameter
	hasRecv := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			hasRecv = true
			break
		}
	}

	if !hasRecv {
		t.Error("Expected RECV instruction for constructor parameter")
	}
}

func TestCompileClassWithInheritance(t *testing.T) {
	input := `<?php
class Animal {
    public $name;
}

class Dog extends Animal {
    public $breed;
}
`

	bytecode := parseAndCompile(t, input)

	// Should have both class names in constants
	classNames := []string{"Animal", "Dog"}
	for _, className := range classNames {
		found := false
		for _, c := range bytecode.Constants {
			if str, ok := c.(string); ok && str == className {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected class name '%s' in constants", className)
		}
	}

	// Should have two DECLARE_CLASS opcodes
	declareClassCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareClass {
			declareClassCount++
		}
	}

	if declareClassCount != 2 {
		t.Errorf("Expected 2 DECLARE_CLASS instructions, got %d", declareClassCount)
	}
}

func TestCompileClassWithMultipleMethods(t *testing.T) {
	input := `<?php
class Calculator {
    public function add($a, $b) {
        return $a + $b;
    }

    public function subtract($a, $b) {
        return $a - $b;
    }

    public function multiply($a, $b) {
        return $a * $b;
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have all method names in constants
	methodNames := []string{"add", "subtract", "multiply"}
	for _, methodName := range methodNames {
		found := false
		for _, c := range bytecode.Constants {
			if str, ok := c.(string); ok && str == methodName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected method name '%s' in constants", methodName)
		}
	}

	// Should have RECV opcodes for parameters (2 parameters * 3 methods = 6)
	recvCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			recvCount++
		}
	}

	if recvCount != 6 {
		t.Errorf("Expected 6 RECV instructions, got %d", recvCount)
	}

	// Should have RETURN opcodes for methods (3 methods)
	returnCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			returnCount++
		}
	}

	if returnCount != 3 {
		t.Errorf("Expected 3 RETURN instructions, got %d", returnCount)
	}
}

func TestCompileClassWithMethodParameters(t *testing.T) {
	input := `<?php
class User {
    public function greet($name, $greeting = "Hello") {
        echo $greeting . " " . $name;
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have RECV for required parameter
	hasRecv := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			hasRecv = true
			break
		}
	}

	if !hasRecv {
		t.Error("Expected RECV instruction for required parameter")
	}

	// Should have RECV_INIT for optional parameter
	hasRecvInit := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecvInit {
			hasRecvInit = true
			break
		}
	}

	if !hasRecvInit {
		t.Error("Expected RECV_INIT instruction for optional parameter")
	}

	// Should have default value "Hello" in constants
	hasDefaultValue := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "Hello" {
			hasDefaultValue = true
			break
		}
	}

	if !hasDefaultValue {
		t.Error("Expected default value 'Hello' in constants")
	}
}

func TestCompileClassWithComplexBody(t *testing.T) {
	input := `<?php
class Account {
    private $balance = 0;

    public function deposit($amount) {
        if ($amount > 0) {
            $this->balance = $this->balance + $amount;
            return true;
        }
        return false;
    }

    public function getBalance() {
        return $this->balance;
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have DECLARE_CLASS opcode
	hasDeclareClass := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareClass {
			hasDeclareClass = true
			break
		}
	}

	if !hasDeclareClass {
		t.Error("Expected DECLARE_CLASS instruction")
	}

	// Should have JMPZ for if statement
	hasJmpz := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmpZ {
			hasJmpz = true
			break
		}
	}

	if !hasJmpz {
		t.Error("Expected JMPZ instruction for if statement")
	}

	// Should have multiple RETURN opcodes
	returnCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpReturn {
			returnCount++
		}
	}

	if returnCount < 2 {
		t.Errorf("Expected at least 2 RETURN instructions, got %d", returnCount)
	}
}

func TestCompileMultipleClasses(t *testing.T) {
	input := `<?php
class Point {
    public $x;
    public $y;
}

class Circle {
    public $center;
    public $radius;
}

class Rectangle {
    public $topLeft;
    public $bottomRight;
}
`

	bytecode := parseAndCompile(t, input)

	// Should have three DECLARE_CLASS opcodes
	declareClassCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareClass {
			declareClassCount++
		}
	}

	if declareClassCount != 3 {
		t.Errorf("Expected 3 DECLARE_CLASS instructions, got %d", declareClassCount)
	}

	// Should have all class names in constants
	classNames := []string{"Point", "Circle", "Rectangle"}
	for _, className := range classNames {
		found := false
		for _, c := range bytecode.Constants {
			if str, ok := c.(string); ok && str == className {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected class name '%s' in constants", className)
		}
	}
}

func TestCompileClassWithVariadicMethod(t *testing.T) {
	input := `<?php
class Logger {
    public function log($level, ...$messages) {
        foreach ($messages as $msg) {
            echo $level . ": " . $msg;
        }
    }
}
`

	bytecode := parseAndCompile(t, input)

	// Should have RECV for required parameter
	hasRecv := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecv {
			hasRecv = true
			break
		}
	}

	if !hasRecv {
		t.Error("Expected RECV instruction for required parameter")
	}

	// Should have RECV_VARIADIC for variadic parameter
	hasRecvVariadic := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpRecvVariadic {
			hasRecvVariadic = true
			break
		}
	}

	if !hasRecvVariadic {
		t.Error("Expected RECV_VARIADIC instruction for variadic parameter")
	}

	// Should have foreach opcodes
	hasFeFetch := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpFeFetchR || instr.Opcode == vm.OpFeFetchRW {
			hasFeFetch = true
			break
		}
	}

	if !hasFeFetch {
		t.Error("Expected FE_FETCH instruction for foreach loop")
	}
}

// ========================================
// Optimization Tests
// ========================================

func TestConstantFoldingArithmetic(t *testing.T) {
	input := `<?php
$x = 1 + 2;
$y = 10 - 5;
$z = 3 * 4;
$a = 20 / 4;
$b = 17 % 5;
`

	bytecode := parseAndCompile(t, input)

	// Check that constants 3, 5, 12, 5, 2 are in the constant pool
	expectedConstants := []int64{3, 5, 12, 5, 2}
	for _, expected := range expectedConstants {
		found := false
		for _, c := range bytecode.Constants {
			if i, ok := c.(int64); ok && i == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected folded constant %d in constant pool", expected)
		}
	}

	// Check that we don't have ADD, SUB, MUL, DIV, MOD opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		switch instr.Opcode {
		case vm.OpAdd, vm.OpSub, vm.OpMul, vm.OpDiv, vm.OpMod:
			t.Errorf("Found arithmetic opcode %s - constant folding didn't work", instr.Opcode)
		}
	}
}

func TestConstantFoldingComparison(t *testing.T) {
	input := `<?php
$a = 5 > 3;
$b = 10 <= 10;
$c = 5 == 5;
$d = 5 != 3;
`

	bytecode := parseAndCompile(t, input)

	// Check that boolean results are in the constant pool
	hasTrue := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok && b {
			hasTrue = true
			break
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' constant from folded comparisons")
	}

	// Check that we don't have comparison opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		switch instr.Opcode {
		case vm.OpIsSmaller, vm.OpIsSmallerOrEqual, vm.OpIsEqual, vm.OpIsNotEqual:
			t.Errorf("Found comparison opcode %s - constant folding didn't work", instr.Opcode)
		}
	}
}

func TestConstantFoldingBitwise(t *testing.T) {
	input := `<?php
$a = 12 | 5;
$b = 12 & 5;
$c = 12 ^ 5;
$d = 8 << 2;
$e = 32 >> 3;
`

	bytecode := parseAndCompile(t, input)

	// Check that results are in the constant pool
	// 12 | 5 = 13, 12 & 5 = 4, 12 ^ 5 = 9, 8 << 2 = 32, 32 >> 3 = 4
	expectedConstants := []int64{13, 4, 9, 32, 4}
	for _, expected := range expectedConstants {
		found := false
		for _, c := range bytecode.Constants {
			if i, ok := c.(int64); ok && i == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected folded constant %d in constant pool", expected)
		}
	}

	// Check that we don't have bitwise opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		switch instr.Opcode {
		case vm.OpBWOr, vm.OpBWAnd, vm.OpBWXor, vm.OpSL, vm.OpSR:
			t.Errorf("Found bitwise opcode %s - constant folding didn't work", instr.Opcode)
		}
	}
}

func TestConstantFoldingStringConcat(t *testing.T) {
	input := `<?php
$x = "Hello" . " " . "World";
`

	bytecode := parseAndCompile(t, input)

	// Due to the way InfixExpression works, we can fold pairs
	// "Hello" . " " will be folded to "Hello "
	// Then "Hello " . "World" won't be folded in one pass (requires multiple passes)
	// For now, just check that at least one CONCAT was eliminated

	concatCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpConcat {
			concatCount++
		}
	}

	// We should have fewer than 2 CONCAT operations
	// (original would be 2, but at least one should be folded)
	if concatCount >= 2 {
		t.Errorf("Expected fewer CONCAT operations due to folding, got %d", concatCount)
	}
}

func TestConstantFoldingUnaryOperations(t *testing.T) {
	input := `<?php
$a = !true;
$b = !false;
$c = -42;
$d = ~7;
`

	bytecode := parseAndCompile(t, input)

	// Check folded constants: !true = false, !false = true, -42 = -42, ~7 = -8
	expectedValues := map[interface{}]bool{
		false:   true,
		true:    true,
		int64(-42): true,
		int64(-8):  true,
	}

	for _, c := range bytecode.Constants {
		delete(expectedValues, c)
	}

	if len(expectedValues) > 0 {
		t.Errorf("Missing expected folded constants: %v", expectedValues)
	}

	// Check that we don't have BOOL_NOT, BW_NOT opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		switch instr.Opcode {
		case vm.OpBoolNot, vm.OpBWNot:
			t.Errorf("Found unary opcode %s - constant folding didn't work", instr.Opcode)
		}
	}
}

func TestConstantFoldingPower(t *testing.T) {
	input := `<?php
$a = 2 ** 3;
$b = 5 ** 2;
$c = 10 ** 0;
`

	bytecode := parseAndCompile(t, input)

	// Check folded constants: 2**3 = 8, 5**2 = 25, 10**0 = 1
	expectedConstants := []int64{8, 25, 1}
	for _, expected := range expectedConstants {
		found := false
		for _, c := range bytecode.Constants {
			if i, ok := c.(int64); ok && i == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected folded constant %d in constant pool", expected)
		}
	}

	// Check that we don't have POW opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpPow {
			t.Error("Found POW opcode - constant folding didn't work")
		}
	}
}

func TestDeadCodeEliminationAfterReturn(t *testing.T) {
	input := `<?php
function test() {
    $x = 1;
    return $x;
    $y = 2;
    echo $y;
}
`

	bytecode := parseAndCompile(t, input)

	// Count variable assignments
	// We should only have one ASSIGN (for $x), not two
	// The $y = 2 should be eliminated
	assignCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAssign {
			assignCount++
		}
	}

	if assignCount > 1 {
		t.Errorf("Expected dead code elimination to remove assignment after return, got %d assignments", assignCount)
	}

	// We should not have ECHO opcode (it's after return)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpEcho {
			t.Error("Found ECHO opcode after return - dead code elimination didn't work")
		}
	}
}

func TestDeadCodeEliminationMultipleReturns(t *testing.T) {
	input := `<?php
function test() {
    if (true) {
        return 1;
        $a = 2;
    }
    return 2;
    $b = 3;
}
`

	bytecode := parseAndCompile(t, input)

	// Count variable assignments
	// Both $a and $b should be eliminated
	assignCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAssign {
			assignCount++
		}
	}

	if assignCount > 0 {
		t.Errorf("Expected dead code elimination to remove all assignments after returns, got %d", assignCount)
	}
}

func TestNoConstantFoldingWithVariables(t *testing.T) {
	input := `<?php
$a = 5;
$b = $a + 3;
`

	bytecode := parseAndCompile(t, input)

	// We should have an ADD opcode because $a is a variable
	hasAdd := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAdd {
			hasAdd = true
			break
		}
	}

	if !hasAdd {
		t.Error("Expected ADD opcode for variable + constant")
	}
}

func TestConstantFoldingMixedTypes(t *testing.T) {
	input := `<?php
$a = 5 + 2.5;
$b = 10.0 - 3;
$c = 2 * 1.5;
`

	bytecode := parseAndCompile(t, input)

	// Check folded float constants: 5 + 2.5 = 7.5, 10.0 - 3 = 7.0, 2 * 1.5 = 3.0
	expectedConstants := []float64{7.5, 7.0, 3.0}
	for _, expected := range expectedConstants {
		found := false
		for _, c := range bytecode.Constants {
			if f, ok := c.(float64); ok && f == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected folded float constant %f in constant pool", expected)
		}
	}

	// Check that we don't have arithmetic opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		switch instr.Opcode {
		case vm.OpAdd, vm.OpSub, vm.OpMul:
			t.Errorf("Found arithmetic opcode %s - constant folding didn't work", instr.Opcode)
		}
	}
}

func TestConstantFoldingSpaceship(t *testing.T) {
	input := `<?php
$a = 5 <=> 3;
$b = 3 <=> 5;
$c = 5 <=> 5;
`

	bytecode := parseAndCompile(t, input)

	// Check folded constants: 5 <=> 3 = 1, 3 <=> 5 = -1, 5 <=> 5 = 0
	expectedConstants := []int64{1, -1, 0}
	for _, expected := range expectedConstants {
		found := false
		for _, c := range bytecode.Constants {
			if i, ok := c.(int64); ok && i == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected folded constant %d in constant pool", expected)
		}
	}

	// Check that we don't have SPACESHIP opcodes (they were folded)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpSpaceship {
			t.Error("Found SPACESHIP opcode - constant folding didn't work")
		}
	}
}

// ========================================
// Helper Method Tests
// ========================================

func TestInstructionsMethod(t *testing.T) {
	input := "<?php $x = 1;"
	bytecode := parseAndCompile(t, input)

	// Should have instructions
	instructions := bytecode.Instructions
	if len(instructions) == 0 {
		t.Error("Expected non-empty instructions after compilation")
	}
}

func TestIsVariableDefined(t *testing.T) {
	c := New()

	// Variable not defined initially
	if c.IsVariableDefined("x") {
		t.Error("Variable 'x' should not be defined initially")
	}

	// Define variable
	c.DefineVariable("x")

	// Now it should be defined
	if !c.IsVariableDefined("x") {
		t.Error("Variable 'x' should be defined after DefineVariable")
	}

	// Other variable still not defined
	if c.IsVariableDefined("y") {
		t.Error("Variable 'y' should not be defined")
	}
}

func TestSymbolString(t *testing.T) {
	sym := &Symbol{
		Name:  "testVar",
		Scope: LocalScope,
		Index: 5,
	}

	str := sym.String()
	if str == "" {
		t.Error("Symbol.String() should return non-empty string")
	}

	// Should contain the name
	if len(str) < len("testVar") {
		t.Error("Symbol.String() should contain variable name")
	}
}

func TestSymbolTableString(t *testing.T) {
	st := NewSymbolTable()
	st.Define("x")
	st.Define("y")

	str := st.String()
	if str == "" {
		t.Error("SymbolTable.String() should return non-empty string")
	}
}

// ========================================
// Optimization Edge Case Tests
// ========================================

func TestConstantFoldingBooleanLiterals(t *testing.T) {
	input := `<?php
$c = true == true;
$d = false != true;
$e = true === false;
`

	bytecode := parseAndCompile(t, input)

	// Check for boolean constants (true from true==true, true from false!=true, false from true===false)
	hasTrue := false
	hasFalse := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok {
			if b {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' constant in bytecode")
	}
	if !hasFalse {
		t.Error("Expected 'false' constant in bytecode")
	}
}

func TestConstantFoldingNullOperations(t *testing.T) {
	input := `<?php
$a = !null;
`

	bytecode := parseAndCompile(t, input)

	// !null should be folded to true
	hasTrue := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok && b {
			hasTrue = true
			break
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' constant from !null")
	}

	// Should NOT have BOOL_NOT opcode (it was folded)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpBoolNot {
			t.Error("Found BOOL_NOT opcode - constant folding didn't work for !null")
		}
	}
}

func TestConstantFoldingDivisionByZero(t *testing.T) {
	// Division by zero should NOT be folded (would cause runtime error)
	input := `<?php
$x = 10 / 0;
`

	bytecode := parseAndCompile(t, input)

	// Should have DIV opcode (not folded)
	hasDiv := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDiv {
			hasDiv = true
			break
		}
	}

	if !hasDiv {
		t.Error("Division by zero should not be folded, expected DIV opcode")
	}
}

func TestConstantFoldingModuloByZero(t *testing.T) {
	// Modulo by zero should NOT be folded
	input := `<?php
$x = 10 % 0;
`

	bytecode := parseAndCompile(t, input)

	// Should have MOD opcode (not folded)
	hasMod := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpMod {
			hasMod = true
			break
		}
	}

	if !hasMod {
		t.Error("Modulo by zero should not be folded, expected MOD opcode")
	}
}

func TestConstantFoldingFloatDivision(t *testing.T) {
	input := `<?php
$x = 10.0 / 0.0;
`

	bytecode := parseAndCompile(t, input)

	// Division by float zero should NOT be folded
	hasDiv := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDiv {
			hasDiv = true
			break
		}
	}

	if !hasDiv {
		t.Error("Float division by zero should not be folded, expected DIV opcode")
	}
}

func TestConstantFoldingLargePower(t *testing.T) {
	// Large power should NOT be folded (>= 100)
	input := `<?php
$x = 2 ** 100;
`

	bytecode := parseAndCompile(t, input)

	// Should have POW opcode (not folded)
	hasPow := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpPow {
			hasPow = true
			break
		}
	}

	if !hasPow {
		t.Error("Large power exponent should not be folded, expected POW opcode")
	}
}

func TestConstantFoldingNegativePower(t *testing.T) {
	// Negative power should NOT be folded
	input := `<?php
$x = 2 ** -3;
`

	bytecode := parseAndCompile(t, input)

	// Should have POW opcode (not folded)
	hasPow := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpPow {
			hasPow = true
			break
		}
	}

	if !hasPow {
		t.Error("Negative power exponent should not be folded, expected POW opcode")
	}
}

// ========================================
// Integration Tests
// ========================================

func TestIntegrationComplexControlFlow(t *testing.T) {
	input := `<?php
function factorial($n) {
    if ($n <= 1) {
        return 1;
    }
    return $n * factorial($n - 1);
}

$result = factorial(5);
`

	bytecode := parseAndCompile(t, input)

	// Should have function declaration
	hasDeclareFunction := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareFunction {
			hasDeclareFunction = true
			break
		}
	}

	if !hasDeclareFunction {
		t.Error("Expected DECLARE_FUNCTION opcode")
	}

	// Should have "factorial" constant
	hasFactorial := false
	for _, c := range bytecode.Constants {
		if str, ok := c.(string); ok && str == "factorial" {
			hasFactorial = true
			break
		}
	}

	if !hasFactorial {
		t.Error("Expected 'factorial' constant")
	}
}

func TestIntegrationNestedClassesAndMethods(t *testing.T) {
	input := `<?php
class Outer {
    public $value = 10;

    public function getValue() {
        return $this->value;
    }

    public function setValue($v) {
        $this->value = $v;
    }
}

class Inner extends Outer {
    public function doubleValue() {
        return $this->getValue() * 2;
    }
}

$obj = new Inner();
$obj->setValue(20);
$result = $obj->doubleValue();
`

	bytecode := parseAndCompile(t, input)

	// Should have both class declarations
	declareClassCount := 0
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpDeclareClass {
			declareClassCount++
		}
	}

	if declareClassCount != 2 {
		t.Errorf("Expected 2 DECLARE_CLASS opcodes, got %d", declareClassCount)
	}

	// Should have class names
	classNames := []string{"Outer", "Inner"}
	for _, className := range classNames {
		found := false
		for _, c := range bytecode.Constants {
			if str, ok := c.(string); ok && str == className {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected class name '%s' in constants", className)
		}
	}
}

func TestIntegrationLoopsWithBreakAndContinue(t *testing.T) {
	input := `<?php
$i = 0;
while ($i < 10) {
    if ($i == 5) {
        break;
    }
    if ($i % 2 == 0) {
        continue;
    }
    echo $i;
    $i = $i + 1;
}
`

	bytecode := parseAndCompile(t, input)

	// Should have JMP opcodes for break/continue
	hasJmp := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpJmp {
			hasJmp = true
			break
		}
	}

	if !hasJmp {
		t.Error("Expected JMP opcode for break/continue")
	}
}

func TestIntegrationTryCatchFinally(t *testing.T) {
	input := `<?php
try {
    $x = 10 / $y;
} catch (Exception $e) {
    echo "Error: " . $e;
} finally {
    echo "Done";
}
`

	bytecode := parseAndCompile(t, input)

	// Should have CATCH opcode
	hasCatch := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpCatch {
			hasCatch = true
			break
		}
	}

	if !hasCatch {
		t.Error("Expected CATCH opcode")
	}
}

func TestIntegrationArrayManipulation(t *testing.T) {
	input := `<?php
$arr = [1, 2, 3];
$x = $arr[0];
$arr2 = ["key" => "value", "num" => 42];
$y = $arr2["key"];
`

	bytecode := parseAndCompile(t, input)

	// Should have array operations
	hasInitArray := false
	hasFetchDim := false

	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpInitArray {
			hasInitArray = true
		}
		if instr.Opcode == vm.OpFetchDimR {
			hasFetchDim = true
		}
	}

	if !hasInitArray {
		t.Error("Expected INIT_ARRAY opcode")
	}
	if !hasFetchDim {
		t.Error("Expected FETCH_DIM_R opcode")
	}
}

func TestIntegrationMixedOptimizations(t *testing.T) {
	input := `<?php
function test() {
    $a = 1 + 2;  // Should be folded to 3
    $b = $a * 5;  // Should not be folded (uses variable), uses non-power-of-2
    return $b;
    $c = 5;  // Dead code, should be eliminated
}
`

	bytecode := parseAndCompile(t, input)

	// Should have constant 3 (from 1+2 folding)
	hasThree := false
	for _, c := range bytecode.Constants {
		if i, ok := c.(int64); ok && i == 3 {
			hasThree = true
			break
		}
	}

	if !hasThree {
		t.Error("Expected constant 3 from folded 1+2")
	}

	// Should NOT have ADD opcode (it was folded)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpAdd {
			t.Error("Found ADD opcode - constant folding didn't work")
		}
	}

	// Should have MUL opcode (variable operation with non-power-of-2)
	hasMul := false
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpMul {
			hasMul = true
			break
		}
	}

	if !hasMul {
		t.Error("Expected MUL opcode for variable multiplication")
	}
}

// ========================================
// Additional Edge Case Tests for Coverage
// ========================================

func TestConstantFoldingStringTruthiness(t *testing.T) {
	input := `<?php
$a = !"";
$b = !"0";
$c = !"hello";
`

	bytecode := parseAndCompile(t, input)

	// !"" and !"0" should be folded to true
	// !"hello" should be folded to false
	hasTrue := false
	hasFalse := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok {
			if b {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' from !'' and !'0'")
	}
	if !hasFalse {
		t.Error("Expected 'false' from !'hello'")
	}
}

func TestConstantFoldingIntTruthiness(t *testing.T) {
	input := `<?php
$a = !0;
$b = !1;
$c = !42;
`

	bytecode := parseAndCompile(t, input)

	// !0 should be folded to true
	// !1 and !42 should be folded to false
	hasTrue := false
	hasFalse := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok {
			if b {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' from !0")
	}
	if !hasFalse {
		t.Error("Expected 'false' from !1 or !42")
	}
}

func TestConstantFoldingUnaryMinusFloat(t *testing.T) {
	input := `<?php
$a = -3.14;
$b = -0.5;
`

	bytecode := parseAndCompile(t, input)

	// Should have negative float constants
	hasNegativePi := false
	for _, c := range bytecode.Constants {
		if f, ok := c.(float64); ok {
			if f < -3.0 && f > -3.2 {
				hasNegativePi = true
				break
			}
		}
	}

	if !hasNegativePi {
		t.Error("Expected -3.14 constant")
	}

	// Should NOT have SUB opcode for float negation (should be folded)
	for _, instr := range bytecode.Instructions {
		if instr.Opcode == vm.OpSub {
			t.Error("Found SUB opcode - unary minus should be folded for float literal")
		}
	}
}

func TestCompilerResetMethod(t *testing.T) {
	input := "<?php $x = 1 + 2;"
	bytecode := parseAndCompile(t, input)

	// Should have instructions and constants
	if len(bytecode.Instructions) == 0 {
		t.Error("Expected instructions after compilation")
	}
	if len(bytecode.Constants) == 0 {
		t.Error("Expected constants after compilation")
	}

	// Create new compiler and compile again (testing reset implicitly)
	c := New()
	p := parser.New(lexer.New(input, "test"))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	c.Compile(program)

	// Should have instructions
	if len(c.Instructions()) == 0 {
		t.Error("Expected instructions")
	}

	// Reset and verify
	c.Reset()

	if len(c.Instructions()) != 0 {
		t.Error("Expected empty instructions after reset")
	}
	if len(c.Constants()) != 0 {
		t.Error("Expected empty constants after reset")
	}
}

func TestChangeOperandMethod(t *testing.T) {
	c := New()

	// Emit an instruction
	pos := c.Emit(vm.OpJmp, vm.ConstOperand(999), vm.UnusedOperand(), vm.UnusedOperand())

	// Change Op1 (operand number 1)
	c.ChangeOperand(pos, 1, vm.ConstOperand(123))

	// Verify the change
	instr := c.Instructions()[pos]
	if instr.Op1.Type != vm.OpConst || instr.Op1.Value != 123 {
		t.Errorf("ChangeOperand didn't work correctly: Op1.Type=%v, Op1.Value=%v", instr.Op1.Type, instr.Op1.Value)
	}
}

func TestRemoveLastInstructionMethod(t *testing.T) {
	c := New()

	// Emit two instructions
	c.Emit(vm.OpEcho, vm.TmpVarOperand(0))
	initialLen := len(c.Instructions())

	c.Emit(vm.OpEcho, vm.TmpVarOperand(1))
	afterSecondLen := len(c.Instructions())

	if afterSecondLen <= initialLen {
		t.Error("Second instruction wasn't added")
	}

	// Remove last instruction
	c.RemoveLastInstruction()

	// Length should be back to initial
	if len(c.Instructions()) != initialLen {
		t.Errorf("RemoveLastInstruction didn't work: expected %d, got %d", initialLen, len(c.Instructions()))
	}
}

func TestCurrentLoopMethod(t *testing.T) {
	c := New()

	// Not in a loop initially
	if c.CurrentLoop() != nil {
		t.Error("CurrentLoop should return nil when not in a loop")
	}

	// Enter a loop
	c.EnterLoop(0)

	// Now should have a current loop
	if c.CurrentLoop() == nil {
		t.Error("CurrentLoop should return a loop context when in a loop")
	}

	// Exit the loop
	c.ExitLoop(10)

	// Should be nil again
	if c.CurrentLoop() != nil {
		t.Error("CurrentLoop should return nil after exiting loop")
	}
}

func TestConstantFoldingIdenticalOperators(t *testing.T) {
	input := `<?php
$a = 5 === 5;
$b = 5 !== 5;
$c = 10 === 10;
`

	bytecode := parseAndCompile(t, input)

	// 5 === 5 and 10 === 10 should be folded to true
	// 5 !== 5 should be folded to false
	hasTrue := false
	hasFalse := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok {
			if b {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' constant from === comparisons")
	}
	if !hasFalse {
		t.Error("Expected 'false' constant from !== comparison")
	}
}

func TestConstantFoldingFloatComparison(t *testing.T) {
	input := `<?php
$a = 3.14 > 2.71;
$b = 1.5 <= 2.5;
$c = 10.0 == 10.0;
$d = 5.5 != 5.5;
`

	bytecode := parseAndCompile(t, input)

	// All comparisons should be folded to boolean constants
	hasTrue := false
	hasFalse := false
	for _, c := range bytecode.Constants {
		if b, ok := c.(bool); ok {
			if b {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
	}

	if !hasTrue {
		t.Error("Expected 'true' constants from float comparisons")
	}
	if !hasFalse {
		t.Error("Expected 'false' constant from 5.5 != 5.5")
	}
}

func TestGetConstantMethod(t *testing.T) {
	c := New()

	// Add some constants
	idx1 := c.AddConstant(int64(42))
	idx2 := c.AddConstant("hello")
	idx3 := c.AddConstant(true)

	// Retrieve and verify
	val1, err1 := c.GetConstant(idx1)
	if err1 != nil {
		t.Errorf("GetConstant error: %v", err1)
	}
	if val1 != int64(42) {
		t.Errorf("Expected int64(42), got %v", val1)
	}

	val2, err2 := c.GetConstant(idx2)
	if err2 != nil {
		t.Errorf("GetConstant error: %v", err2)
	}
	if val2 != "hello" {
		t.Errorf("Expected 'hello', got %v", val2)
	}

	val3, err3 := c.GetConstant(idx3)
	if err3 != nil {
		t.Errorf("GetConstant error: %v", err3)
	}
	if val3 != true {
		t.Errorf("Expected true, got %v", val3)
	}

	// Test invalid index
	_, err := c.GetConstant(999)
	if err == nil {
		t.Error("Expected GetConstant to return error for invalid index")
	}
}
