package vm

import "testing"

// TestOpcodeValues ensures all opcodes have correct numeric values
func TestOpcodeValues(t *testing.T) {
	tests := []struct {
		op       Opcode
		expected uint8
		name     string
	}{
		{OpNop, 0, "NOP"},
		{OpAdd, 1, "ADD"},
		{OpSub, 2, "SUB"},
		{OpMul, 3, "MUL"},
		{OpDiv, 4, "DIV"},
		{OpMod, 5, "MOD"},
		{OpEcho, 136, "ECHO"},
		{OpDeclareAttributedConst, 210, "DECLARE_ATTRIBUTED_CONST"},
	}

	for _, tt := range tests {
		if uint8(tt.op) != tt.expected {
			t.Errorf("%s: expected value %d, got %d", tt.name, tt.expected, uint8(tt.op))
		}
	}
}

// TestOpcodeString verifies the String() method returns correct names
func TestOpcodeString(t *testing.T) {
	tests := []struct {
		op       Opcode
		expected string
	}{
		{OpNop, "NOP"},
		{OpAdd, "ADD"},
		{OpSub, "SUB"},
		{OpMul, "MUL"},
		{OpDiv, "DIV"},
		{OpMod, "MOD"},
		{OpConcat, "CONCAT"},
		{OpEcho, "ECHO"},
		{OpAssign, "ASSIGN"},
		{OpJmp, "JMP"},
		{OpJmpZ, "JMPZ"},
		{OpJmpNZ, "JMPNZ"},
		{OpInitFcall, "INIT_FCALL"},
		{OpDoFcall, "DO_FCALL"},
		{OpReturn, "RETURN"},
		{OpNew, "NEW"},
		{OpClone, "CLONE"},
		{OpThrow, "THROW"},
		{OpCatch, "CATCH"},
		{OpYield, "YIELD"},
		{OpMatch, "MATCH"},
		{OpDeclareAttributedConst, "DECLARE_ATTRIBUTED_CONST"},
	}

	for _, tt := range tests {
		got := tt.op.String()
		if got != tt.expected {
			t.Errorf("%s.String(): expected %q, got %q", tt.expected, tt.expected, got)
		}
	}
}

// TestOpcodeStringUnknown tests String() for invalid opcodes
func TestOpcodeStringUnknown(t *testing.T) {
	invalidOp := Opcode(255)
	if invalidOp.String() != "UNKNOWN" {
		t.Errorf("Invalid opcode should return 'UNKNOWN', got %q", invalidOp.String())
	}
}

// TestOpcodeCount verifies we have exactly 211 opcodes (0-210)
func TestOpcodeCount(t *testing.T) {
	if OpcodeLast != 210 {
		t.Errorf("Expected OpcodeLast to be 210, got %d", OpcodeLast)
	}

	// Verify opcodeNames array has correct length
	expectedLength := OpcodeLast + 1
	if len(opcodeNames) != expectedLength {
		t.Errorf("Expected opcodeNames length %d, got %d", expectedLength, len(opcodeNames))
	}
}

// TestArithmeticOpcodes tests all arithmetic opcodes
func TestArithmeticOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpAdd, "ADD"},
		{OpSub, "SUB"},
		{OpMul, "MUL"},
		{OpDiv, "DIV"},
		{OpMod, "MOD"},
		{OpPow, "POW"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Arithmetic opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestComparisonOpcodes tests all comparison opcodes
func TestComparisonOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpIsIdentical, "IS_IDENTICAL"},
		{OpIsNotIdentical, "IS_NOT_IDENTICAL"},
		{OpIsEqual, "IS_EQUAL"},
		{OpIsNotEqual, "IS_NOT_EQUAL"},
		{OpIsSmaller, "IS_SMALLER"},
		{OpIsSmallerOrEqual, "IS_SMALLER_OR_EQUAL"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Comparison opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestControlFlowOpcodes tests all control flow opcodes
func TestControlFlowOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpJmp, "JMP"},
		{OpJmpZ, "JMPZ"},
		{OpJmpNZ, "JMPNZ"},
		{OpJmpZEx, "JMPZ_EX"},
		{OpJmpNZEx, "JMPNZ_EX"},
		{OpCase, "CASE"},
		{OpJmpSet, "JMP_SET"},
		{OpJmpNull, "JMP_NULL"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Control flow opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestFunctionCallOpcodes tests function call opcodes
func TestFunctionCallOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpInitFcall, "INIT_FCALL"},
		{OpInitFcallByName, "INIT_FCALL_BY_NAME"},
		{OpDoFcall, "DO_FCALL"},
		{OpDoIcall, "DO_ICALL"},
		{OpDoUcall, "DO_UCALL"},
		{OpReturn, "RETURN"},
		{OpReturnByRef, "RETURN_BY_REF"},
		{OpSendVal, "SEND_VAL"},
		{OpSendVar, "SEND_VAR"},
		{OpSendRef, "SEND_REF"},
		{OpRecv, "RECV"},
		{OpRecvInit, "RECV_INIT"},
		{OpRecvVariadic, "RECV_VARIADIC"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Function call opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestObjectOpcodes tests object-related opcodes
func TestObjectOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpNew, "NEW"},
		{OpClone, "CLONE"},
		{OpInstanceof, "INSTANCEOF"},
		{OpFetchClass, "FETCH_CLASS"},
		{OpFetchClassName, "FETCH_CLASS_NAME"},
		{OpGetClass, "GET_CLASS"},
		{OpGetCalledClass, "GET_CALLED_CLASS"},
		{OpInitMethodCall, "INIT_METHOD_CALL"},
		{OpInitStaticMethodCall, "INIT_STATIC_METHOD_CALL"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Object opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestArrayOpcodes tests array-related opcodes
func TestArrayOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpInitArray, "INIT_ARRAY"},
		{OpAddArrayElement, "ADD_ARRAY_ELEMENT"},
		{OpAddArrayUnpack, "ADD_ARRAY_UNPACK"},
		{OpFetchDimR, "FETCH_DIM_R"},
		{OpFetchDimW, "FETCH_DIM_W"},
		{OpAssignDim, "ASSIGN_DIM"},
		{OpUnsetDim, "UNSET_DIM"},
		{OpInArray, "IN_ARRAY"},
		{OpCount, "COUNT"},
		{OpArrayKeyExists, "ARRAY_KEY_EXISTS"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Array opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestGeneratorOpcodes tests generator-related opcodes
func TestGeneratorOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpGeneratorCreate, "GENERATOR_CREATE"},
		{OpYield, "YIELD"},
		{OpYieldFrom, "YIELD_FROM"},
		{OpGeneratorReturn, "GENERATOR_RETURN"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Generator opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestExceptionOpcodes tests exception handling opcodes
func TestExceptionOpcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpThrow, "THROW"},
		{OpCatch, "CATCH"},
		{OpHandleException, "HANDLE_EXCEPTION"},
		{OpDiscardException, "DISCARD_EXCEPTION"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("Exception opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestPHP8Opcodes tests PHP 8.x specific opcodes
func TestPHP8Opcodes(t *testing.T) {
	ops := []struct {
		op   Opcode
		name string
	}{
		{OpMatch, "MATCH"},
		{OpCaseStrict, "CASE_STRICT"},
		{OpMatchError, "MATCH_ERROR"},
		{OpCallableConvert, "CALLABLE_CONVERT"},
		{OpVerifyNeverType, "VERIFY_NEVER_TYPE"},
		{OpDeclareAttributedConst, "DECLARE_ATTRIBUTED_CONST"},
		{OpInitParentPropertyHookCall, "INIT_PARENT_PROPERTY_HOOK_CALL"},
	}

	for _, tt := range ops {
		if tt.op.String() != tt.name {
			t.Errorf("PHP 8 opcode %s has wrong name: %s", tt.name, tt.op.String())
		}
	}
}

// TestOpcodeSequence verifies there are no gaps in the opcode sequence
func TestOpcodeSequence(t *testing.T) {
	// These opcodes are intentionally missing in PHP's implementation
	missingOpcodes := map[uint8]bool{
		45: true, // Gap between JMPNZ and JMPZ_EX
		79: true, // Gap between FE_FETCH_R and FETCH_R
	}

	for i := uint8(0); i <= OpcodeLast; i++ {
		if missingOpcodes[i] {
			continue
		}

		op := Opcode(i)
		name := op.String()

		if name == "UNKNOWN" {
			t.Errorf("Opcode %d returns UNKNOWN, should have a name", i)
		}

		if opcodeNames[i] == "" {
			t.Errorf("Opcode %d has empty name in opcodeNames array", i)
		}
	}
}
