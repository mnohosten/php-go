package vm

import (
	"testing"
)

// TestOperandTypeString tests the String() method for operand types
func TestOperandTypeString(t *testing.T) {
	tests := []struct {
		typ      OperandType
		expected string
	}{
		{OpUnused, "UNUSED"},
		{OpConst, "CONST"},
		{OpTmpVar, "TMPVAR"},
		{OpVar, "VAR"},
		{OpCV, "CV"},
		{OperandType(99), "UNKNOWN(99)"},
	}

	for _, tt := range tests {
		got := tt.typ.String()
		if got != tt.expected {
			t.Errorf("OperandType(%d).String() = %q, want %q", tt.typ, got, tt.expected)
		}
	}
}

// TestOperandCreation tests operand constructor functions
func TestOperandCreation(t *testing.T) {
	tests := []struct {
		name     string
		operand  Operand
		wantType OperandType
		wantVal  uint32
	}{
		{"unused", UnusedOperand(), OpUnused, 0},
		{"const", ConstOperand(5), OpConst, 5},
		{"tmpvar", TmpVarOperand(10), OpTmpVar, 10},
		{"var", VarOperand(15), OpVar, 15},
		{"cv", CVOperand(20), OpCV, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.operand.Type != tt.wantType {
				t.Errorf("%s: Type = %v, want %v", tt.name, tt.operand.Type, tt.wantType)
			}
			if tt.operand.Value != tt.wantVal {
				t.Errorf("%s: Value = %d, want %d", tt.name, tt.operand.Value, tt.wantVal)
			}
		})
	}
}

// TestOperandIsChecks tests the Is*() helper methods
func TestOperandIsChecks(t *testing.T) {
	tests := []struct {
		name     string
		operand  Operand
		isUnused bool
		isConst  bool
		isTmpVar bool
		isVar    bool
		isCV     bool
	}{
		{"unused", UnusedOperand(), true, false, false, false, false},
		{"const", ConstOperand(1), false, true, false, false, false},
		{"tmpvar", TmpVarOperand(1), false, false, true, false, false},
		{"var", VarOperand(1), false, false, false, true, false},
		{"cv", CVOperand(1), false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.operand.IsUnused() != tt.isUnused {
				t.Errorf("%s.IsUnused() = %v, want %v", tt.name, tt.operand.IsUnused(), tt.isUnused)
			}
			if tt.operand.IsConst() != tt.isConst {
				t.Errorf("%s.IsConst() = %v, want %v", tt.name, tt.operand.IsConst(), tt.isConst)
			}
			if tt.operand.IsTmpVar() != tt.isTmpVar {
				t.Errorf("%s.IsTmpVar() = %v, want %v", tt.name, tt.operand.IsTmpVar(), tt.isTmpVar)
			}
			if tt.operand.IsVar() != tt.isVar {
				t.Errorf("%s.IsVar() = %v, want %v", tt.name, tt.operand.IsVar(), tt.isVar)
			}
			if tt.operand.IsCV() != tt.isCV {
				t.Errorf("%s.IsCV() = %v, want %v", tt.name, tt.operand.IsCV(), tt.isCV)
			}
		})
	}
}

// TestOperandString tests the String() method for operands
func TestOperandString(t *testing.T) {
	tests := []struct {
		operand  Operand
		expected string
	}{
		{UnusedOperand(), "<unused>"},
		{ConstOperand(5), "CONST(5)"},
		{TmpVarOperand(10), "TMPVAR(10)"},
		{VarOperand(15), "VAR(15)"},
		{CVOperand(20), "CV(20)"},
	}

	for _, tt := range tests {
		got := tt.operand.String()
		if got != tt.expected {
			t.Errorf("Operand.String() = %q, want %q", got, tt.expected)
		}
	}
}

// TestNewInstruction tests instruction creation
func TestNewInstruction(t *testing.T) {
	instr := NewInstruction(OpAdd, 42)

	if instr.Opcode != OpAdd {
		t.Errorf("Opcode = %v, want %v", instr.Opcode, OpAdd)
	}
	if instr.Lineno != 42 {
		t.Errorf("Lineno = %d, want %d", instr.Lineno, 42)
	}
	if !instr.Op1.IsUnused() || !instr.Op2.IsUnused() || !instr.Result.IsUnused() {
		t.Error("New instruction should have all operands as unused")
	}
	if instr.ExtendedValue != 0 {
		t.Errorf("ExtendedValue = %d, want 0", instr.ExtendedValue)
	}
}

// TestInstructionBuilder tests the builder pattern methods
func TestInstructionBuilder(t *testing.T) {
	instr := NewInstruction(OpAdd, 10).
		WithResult(OpTmpVar, 1).
		WithOp1(OpCV, 2).
		WithOp2(OpConst, 3).
		WithExtended(100)

	if instr.Opcode != OpAdd {
		t.Errorf("Opcode = %v, want %v", instr.Opcode, OpAdd)
	}
	if instr.Result != TmpVarOperand(1) {
		t.Errorf("Result = %v, want TmpVar(1)", instr.Result)
	}
	if instr.Op1 != CVOperand(2) {
		t.Errorf("Op1 = %v, want CV(2)", instr.Op1)
	}
	if instr.Op2 != ConstOperand(3) {
		t.Errorf("Op2 = %v, want Const(3)", instr.Op2)
	}
	if instr.ExtendedValue != 100 {
		t.Errorf("ExtendedValue = %d, want 100", instr.ExtendedValue)
	}
	if instr.Lineno != 10 {
		t.Errorf("Lineno = %d, want 10", instr.Lineno)
	}
}

// TestInstructionString tests the String() method for instructions
func TestInstructionString(t *testing.T) {
	// Simple instruction with no operands
	instr1 := NewInstruction(OpNop, 1)
	str1 := instr1.String()
	if str1 == "" {
		t.Error("String() returned empty for NOP instruction")
	}

	// Instruction with all operands
	instr2 := NewInstruction(OpAdd, 10).
		WithResult(OpTmpVar, 1).
		WithOp1(OpCV, 2).
		WithOp2(OpConst, 3)
	str2 := instr2.String()
	if str2 == "" {
		t.Error("String() returned empty for ADD instruction")
	}

	// Instruction with extended value
	instr3 := NewInstruction(OpCast, 20).
		WithResult(OpTmpVar, 1).
		WithOp1(OpCV, 2).
		WithExtended(4) // Cast to int
	str3 := instr3.String()
	if str3 == "" {
		t.Error("String() returned empty for CAST instruction")
	}
}

// TestInstructionEncodeDecode tests encoding and decoding of instructions
func TestInstructionEncodeDecode(t *testing.T) {
	original := NewInstruction(OpAdd, 42).
		WithResult(OpTmpVar, 1).
		WithOp1(OpCV, 2).
		WithOp2(OpConst, 3).
		WithExtended(100)

	// Encode
	encoded := original.Encode()

	// Check size
	if len(encoded) != InstructionSize {
		t.Fatalf("Encoded size = %d, want %d", len(encoded), InstructionSize)
	}

	// Decode
	decoded, err := DecodeInstruction(encoded)
	if err != nil {
		t.Fatalf("DecodeInstruction failed: %v", err)
	}

	// Compare all fields
	if decoded.Opcode != original.Opcode {
		t.Errorf("Opcode: got %v, want %v", decoded.Opcode, original.Opcode)
	}
	if decoded.Op1 != original.Op1 {
		t.Errorf("Op1: got %v, want %v", decoded.Op1, original.Op1)
	}
	if decoded.Op2 != original.Op2 {
		t.Errorf("Op2: got %v, want %v", decoded.Op2, original.Op2)
	}
	if decoded.Result != original.Result {
		t.Errorf("Result: got %v, want %v", decoded.Result, original.Result)
	}
	if decoded.ExtendedValue != original.ExtendedValue {
		t.Errorf("ExtendedValue: got %d, want %d", decoded.ExtendedValue, original.ExtendedValue)
	}
	if decoded.Lineno != original.Lineno {
		t.Errorf("Lineno: got %d, want %d", decoded.Lineno, original.Lineno)
	}
}

// TestInstructionEncodeDecodeMultiple tests encoding/decoding multiple instructions
func TestInstructionEncodeDecodeMultiple(t *testing.T) {
	instrs := Instructions{
		*NewInstruction(OpAdd, 1).WithResult(OpTmpVar, 1).WithOp1(OpCV, 0).WithOp2(OpConst, 0),
		*NewInstruction(OpSub, 2).WithResult(OpTmpVar, 2).WithOp1(OpTmpVar, 1).WithOp2(OpConst, 1),
		*NewInstruction(OpEcho, 3).WithOp1(OpTmpVar, 2),
		*NewInstruction(OpReturn, 4),
	}

	// Encode
	encoded := instrs.Encode()

	// Check size
	expectedSize := len(instrs) * InstructionSize
	if len(encoded) != expectedSize {
		t.Fatalf("Encoded size = %d, want %d", len(encoded), expectedSize)
	}

	// Decode
	decoded, err := DecodeInstructions(encoded)
	if err != nil {
		t.Fatalf("DecodeInstructions failed: %v", err)
	}

	// Check count
	if len(decoded) != len(instrs) {
		t.Fatalf("Decoded count = %d, want %d", len(decoded), len(instrs))
	}

	// Compare each instruction
	for i := range instrs {
		if decoded[i].Opcode != instrs[i].Opcode {
			t.Errorf("Instruction %d: Opcode = %v, want %v", i, decoded[i].Opcode, instrs[i].Opcode)
		}
		if decoded[i].Lineno != instrs[i].Lineno {
			t.Errorf("Instruction %d: Lineno = %d, want %d", i, decoded[i].Lineno, instrs[i].Lineno)
		}
	}
}

// TestDecodeInstructionError tests error handling in DecodeInstruction
func TestDecodeInstructionError(t *testing.T) {
	// Buffer too small
	_, err := DecodeInstruction(make([]byte, InstructionSize-1))
	if err == nil {
		t.Error("Expected error for buffer too small, got nil")
	}

	// Empty buffer
	_, err = DecodeInstruction([]byte{})
	if err == nil {
		t.Error("Expected error for empty buffer, got nil")
	}
}

// TestDecodeInstructionsError tests error handling in DecodeInstructions
func TestDecodeInstructionsError(t *testing.T) {
	// Buffer not a multiple of instruction size
	_, err := DecodeInstructions(make([]byte, InstructionSize+1))
	if err == nil {
		t.Error("Expected error for invalid buffer size, got nil")
	}
}

// TestInstructionSize verifies the instruction size constant
func TestInstructionSize(t *testing.T) {
	// 1 byte opcode + 3 bytes operand types + 12 bytes operand values + 4 bytes extended + 4 bytes lineno = 24 bytes
	expectedSize := 1 + 3 + 12 + 4 + 4
	if InstructionSize != expectedSize {
		t.Errorf("InstructionSize = %d, want %d", InstructionSize, expectedSize)
	}
}

// TestInstructionsString tests the String() method for instruction sequences
func TestInstructionsString(t *testing.T) {
	instrs := Instructions{
		*NewInstruction(OpAdd, 1).WithResult(OpTmpVar, 1).WithOp1(OpCV, 0).WithOp2(OpConst, 0),
		*NewInstruction(OpEcho, 2).WithOp1(OpTmpVar, 1),
	}

	str := instrs.String()
	if str == "" {
		t.Error("Instructions.String() returned empty string")
	}

	// Should contain instruction numbers
	if len(str) < 10 {
		t.Errorf("Instructions.String() too short: %q", str)
	}
}

// TestRealWorldExample tests a realistic instruction sequence
func TestRealWorldExample(t *testing.T) {
	// Simulate: $result = $a + $b;
	instrs := Instructions{
		// $tmp1 = $a + $b
		*NewInstruction(OpAdd, 1).
			WithResult(OpTmpVar, 0).
			WithOp1(OpCV, 0). // $a (compiled variable 0)
			WithOp2(OpCV, 1), // $b (compiled variable 1)

		// $result = $tmp1
		*NewInstruction(OpAssign, 1).
			WithResult(OpCV, 2).     // $result (compiled variable 2)
			WithOp1(OpTmpVar, 0).    // $tmp1
			WithOp2(OpUnused, 0),    // unused

		// Free $tmp1
		*NewInstruction(OpFree, 1).
			WithOp1(OpTmpVar, 0),
	}

	// Encode and decode
	encoded := instrs.Encode()
	decoded, err := DecodeInstructions(encoded)
	if err != nil {
		t.Fatalf("Failed to encode/decode: %v", err)
	}

	// Verify instruction count
	if len(decoded) != 3 {
		t.Fatalf("Expected 3 instructions, got %d", len(decoded))
	}

	// Verify first instruction (ADD)
	if decoded[0].Opcode != OpAdd {
		t.Errorf("First instruction should be ADD, got %v", decoded[0].Opcode)
	}
	if !decoded[0].Result.IsTmpVar() || decoded[0].Result.Value != 0 {
		t.Errorf("ADD result should be TmpVar(0), got %v", decoded[0].Result)
	}

	// Verify second instruction (ASSIGN)
	if decoded[1].Opcode != OpAssign {
		t.Errorf("Second instruction should be ASSIGN, got %v", decoded[1].Opcode)
	}
	if !decoded[1].Result.IsCV() || decoded[1].Result.Value != 2 {
		t.Errorf("ASSIGN result should be CV(2), got %v", decoded[1].Result)
	}

	// Verify third instruction (FREE)
	if decoded[2].Opcode != OpFree {
		t.Errorf("Third instruction should be FREE, got %v", decoded[2].Opcode)
	}
}

// TestOperandTypeValues verifies the operand type bit flags
func TestOperandTypeValues(t *testing.T) {
	if OpUnused != 0 {
		t.Errorf("OpUnused = %d, want 0", OpUnused)
	}
	if OpConst != 1 {
		t.Errorf("OpConst = %d, want 1", OpConst)
	}
	if OpTmpVar != 2 {
		t.Errorf("OpTmpVar = %d, want 2", OpTmpVar)
	}
	if OpVar != 4 {
		t.Errorf("OpVar = %d, want 4", OpVar)
	}
	if OpCV != 8 {
		t.Errorf("OpCV = %d, want 8", OpCV)
	}
}

// TestLargeValues tests encoding/decoding with large values
func TestLargeValues(t *testing.T) {
	// Test with maximum uint32 values
	maxUint32 := uint32(0xFFFFFFFF)

	instr := NewInstruction(OpAdd, maxUint32).
		WithResult(OpTmpVar, maxUint32).
		WithOp1(OpCV, maxUint32).
		WithOp2(OpConst, maxUint32).
		WithExtended(maxUint32)

	// Encode and decode
	encoded := instr.Encode()
	decoded, err := DecodeInstruction(encoded)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	// Verify all fields preserved
	if decoded.Result.Value != maxUint32 {
		t.Errorf("Result value = %d, want %d", decoded.Result.Value, maxUint32)
	}
	if decoded.Op1.Value != maxUint32 {
		t.Errorf("Op1 value = %d, want %d", decoded.Op1.Value, maxUint32)
	}
	if decoded.Op2.Value != maxUint32 {
		t.Errorf("Op2 value = %d, want %d", decoded.Op2.Value, maxUint32)
	}
	if decoded.ExtendedValue != maxUint32 {
		t.Errorf("ExtendedValue = %d, want %d", decoded.ExtendedValue, maxUint32)
	}
	if decoded.Lineno != maxUint32 {
		t.Errorf("Lineno = %d, want %d", decoded.Lineno, maxUint32)
	}
}
