package vm

import "fmt"

// ========================================
// Operand Types
// ========================================

// OperandType represents the type of an instruction operand
type OperandType uint8

const (
	// OpUnused - Operand is not used
	OpUnused OperandType = 0

	// OpConst - Operand is a constant from the constant table
	OpConst OperandType = 1 << 0 // 1

	// OpTmpVar - Operand is a temporary variable (intermediate computation result)
	OpTmpVar OperandType = 1 << 1 // 2

	// OpVar - Operand is a variable (may require indirect fetch)
	OpVar OperandType = 1 << 2 // 4

	// OpCV - Operand is a compiled variable (direct stack access, optimized)
	OpCV OperandType = 1 << 3 // 8
)

// String returns the name of the operand type
func (t OperandType) String() string {
	switch t {
	case OpUnused:
		return "UNUSED"
	case OpConst:
		return "CONST"
	case OpTmpVar:
		return "TMPVAR"
	case OpVar:
		return "VAR"
	case OpCV:
		return "CV"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

// ========================================
// Operand
// ========================================

// Operand represents a single operand in an instruction
// Each operand has a type and a value (index/number depending on type)
type Operand struct {
	// Type of the operand (unused, const, tmp, var, cv)
	Type OperandType

	// Value depends on the type:
	// - OpConst: index into constant table
	// - OpTmpVar: temporary variable number
	// - OpVar: variable number
	// - OpCV: compiled variable index
	// - OpUnused: ignored
	Value uint32
}

// String returns a human-readable representation of the operand
func (op Operand) String() string {
	if op.Type == OpUnused {
		return "<unused>"
	}
	return fmt.Sprintf("%s(%d)", op.Type, op.Value)
}

// IsUnused returns true if the operand is not used
func (op Operand) IsUnused() bool {
	return op.Type == OpUnused
}

// IsConst returns true if the operand is a constant
func (op Operand) IsConst() bool {
	return op.Type == OpConst
}

// IsTmpVar returns true if the operand is a temporary variable
func (op Operand) IsTmpVar() bool {
	return op.Type == OpTmpVar
}

// IsVar returns true if the operand is a variable
func (op Operand) IsVar() bool {
	return op.Type == OpVar
}

// IsCV returns true if the operand is a compiled variable
func (op Operand) IsCV() bool {
	return op.Type == OpCV
}

// ========================================
// Instruction
// ========================================

// Instruction represents a single VM instruction (analogous to zend_op)
// This is the core unit of execution in the VM
type Instruction struct {
	// Opcode - the operation to perform
	Opcode Opcode

	// Op1 - first operand
	Op1 Operand

	// Op2 - second operand
	Op2 Operand

	// Result - where to store the result
	Result Operand

	// ExtendedValue - additional data for certain opcodes
	// Used for:
	// - Type cast target type (OpCast)
	// - Number of arguments (function calls)
	// - Fetch mode flags
	// - Assignment operation type (+=, -=, etc.)
	// - And other opcode-specific data
	ExtendedValue uint32

	// Lineno - source code line number (for debugging and error messages)
	Lineno uint32
}

// String returns a human-readable representation of the instruction
func (instr Instruction) String() string {
	result := fmt.Sprintf("%-20s", instr.Opcode.String())

	// Add result if used
	if !instr.Result.IsUnused() {
		result += fmt.Sprintf(" %s = ", instr.Result)
	} else {
		result += " "
	}

	// Add op1 if used
	if !instr.Op1.IsUnused() {
		result += fmt.Sprintf("%s", instr.Op1)
	}

	// Add op2 if used
	if !instr.Op2.IsUnused() {
		result += fmt.Sprintf(", %s", instr.Op2)
	}

	// Add extended value if non-zero
	if instr.ExtendedValue != 0 {
		result += fmt.Sprintf(" [ext=%d]", instr.ExtendedValue)
	}

	// Add line number
	result += fmt.Sprintf(" (line %d)", instr.Lineno)

	return result
}

// ========================================
// Instruction Builder Helpers
// ========================================

// NewInstruction creates a new instruction with the given opcode
func NewInstruction(opcode Opcode, lineno uint32) *Instruction {
	return &Instruction{
		Opcode: opcode,
		Lineno: lineno,
	}
}

// WithResult sets the result operand and returns the instruction (for chaining)
func (instr *Instruction) WithResult(typ OperandType, value uint32) *Instruction {
	instr.Result = Operand{Type: typ, Value: value}
	return instr
}

// WithOp1 sets the first operand and returns the instruction (for chaining)
func (instr *Instruction) WithOp1(typ OperandType, value uint32) *Instruction {
	instr.Op1 = Operand{Type: typ, Value: value}
	return instr
}

// WithOp2 sets the second operand and returns the instruction (for chaining)
func (instr *Instruction) WithOp2(typ OperandType, value uint32) *Instruction {
	instr.Op2 = Operand{Type: typ, Value: value}
	return instr
}

// WithExtended sets the extended value and returns the instruction (for chaining)
func (instr *Instruction) WithExtended(value uint32) *Instruction {
	instr.ExtendedValue = value
	return instr
}

// ========================================
// Operand Constructors
// ========================================

// UnusedOperand creates an unused operand
func UnusedOperand() Operand {
	return Operand{Type: OpUnused, Value: 0}
}

// ConstOperand creates a constant operand
func ConstOperand(index uint32) Operand {
	return Operand{Type: OpConst, Value: index}
}

// TmpVarOperand creates a temporary variable operand
func TmpVarOperand(number uint32) Operand {
	return Operand{Type: OpTmpVar, Value: number}
}

// VarOperand creates a variable operand
func VarOperand(number uint32) Operand {
	return Operand{Type: OpVar, Value: number}
}

// CVOperand creates a compiled variable operand
func CVOperand(index uint32) Operand {
	return Operand{Type: OpCV, Value: index}
}

// ========================================
// Encoding/Decoding
// ========================================

// InstructionSize is the size of an encoded instruction in bytes
// Layout: [opcode:1][op1_type:1][op2_type:1][result_type:1][op1_value:4][op2_value:4][result_value:4][extended:4][lineno:4]
const InstructionSize = 24

// Encode serializes the instruction into a byte slice
func (instr *Instruction) Encode() []byte {
	buf := make([]byte, InstructionSize)

	// Byte 0: opcode
	buf[0] = uint8(instr.Opcode)

	// Byte 1: op1 type
	buf[1] = uint8(instr.Op1.Type)

	// Byte 2: op2 type
	buf[2] = uint8(instr.Op2.Type)

	// Byte 3: result type
	buf[3] = uint8(instr.Result.Type)

	// Bytes 4-7: op1 value (little-endian)
	buf[4] = uint8(instr.Op1.Value)
	buf[5] = uint8(instr.Op1.Value >> 8)
	buf[6] = uint8(instr.Op1.Value >> 16)
	buf[7] = uint8(instr.Op1.Value >> 24)

	// Bytes 8-11: op2 value (little-endian)
	buf[8] = uint8(instr.Op2.Value)
	buf[9] = uint8(instr.Op2.Value >> 8)
	buf[10] = uint8(instr.Op2.Value >> 16)
	buf[11] = uint8(instr.Op2.Value >> 24)

	// Bytes 12-15: result value (little-endian)
	buf[12] = uint8(instr.Result.Value)
	buf[13] = uint8(instr.Result.Value >> 8)
	buf[14] = uint8(instr.Result.Value >> 16)
	buf[15] = uint8(instr.Result.Value >> 24)

	// Bytes 16-19: extended value (little-endian)
	buf[16] = uint8(instr.ExtendedValue)
	buf[17] = uint8(instr.ExtendedValue >> 8)
	buf[18] = uint8(instr.ExtendedValue >> 16)
	buf[19] = uint8(instr.ExtendedValue >> 24)

	// Bytes 20-23: line number (little-endian)
	buf[20] = uint8(instr.Lineno)
	buf[21] = uint8(instr.Lineno >> 8)
	buf[22] = uint8(instr.Lineno >> 16)
	buf[23] = uint8(instr.Lineno >> 24)

	return buf
}

// DecodeInstruction deserializes an instruction from a byte slice
func DecodeInstruction(buf []byte) (*Instruction, error) {
	if len(buf) < InstructionSize {
		return nil, fmt.Errorf("buffer too small: expected %d bytes, got %d", InstructionSize, len(buf))
	}

	instr := &Instruction{}

	// Byte 0: opcode
	instr.Opcode = Opcode(buf[0])

	// Byte 1: op1 type
	instr.Op1.Type = OperandType(buf[1])

	// Byte 2: op2 type
	instr.Op2.Type = OperandType(buf[2])

	// Byte 3: result type
	instr.Result.Type = OperandType(buf[3])

	// Bytes 4-7: op1 value (little-endian)
	instr.Op1.Value = uint32(buf[4]) |
		uint32(buf[5])<<8 |
		uint32(buf[6])<<16 |
		uint32(buf[7])<<24

	// Bytes 8-11: op2 value (little-endian)
	instr.Op2.Value = uint32(buf[8]) |
		uint32(buf[9])<<8 |
		uint32(buf[10])<<16 |
		uint32(buf[11])<<24

	// Bytes 12-15: result value (little-endian)
	instr.Result.Value = uint32(buf[12]) |
		uint32(buf[13])<<8 |
		uint32(buf[14])<<16 |
		uint32(buf[15])<<24

	// Bytes 16-19: extended value (little-endian)
	instr.ExtendedValue = uint32(buf[16]) |
		uint32(buf[17])<<8 |
		uint32(buf[18])<<16 |
		uint32(buf[19])<<24

	// Bytes 20-23: line number (little-endian)
	instr.Lineno = uint32(buf[20]) |
		uint32(buf[21])<<8 |
		uint32(buf[22])<<16 |
		uint32(buf[23])<<24

	return instr, nil
}

// ========================================
// Instruction Sequence
// ========================================

// Instructions is a sequence of instructions (a compiled program or function)
type Instructions []Instruction

// String returns a human-readable representation of the instruction sequence
func (instrs Instructions) String() string {
	result := ""
	for i, instr := range instrs {
		result += fmt.Sprintf("%04d: %s\n", i, instr.String())
	}
	return result
}

// Encode serializes all instructions into a byte slice
func (instrs Instructions) Encode() []byte {
	buf := make([]byte, 0, len(instrs)*InstructionSize)
	for i := range instrs {
		buf = append(buf, instrs[i].Encode()...)
	}
	return buf
}

// DecodeInstructions deserializes a sequence of instructions from a byte slice
func DecodeInstructions(buf []byte) (Instructions, error) {
	if len(buf)%InstructionSize != 0 {
		return nil, fmt.Errorf("buffer size not a multiple of instruction size: %d", len(buf))
	}

	count := len(buf) / InstructionSize
	instrs := make(Instructions, count)

	for i := 0; i < count; i++ {
		offset := i * InstructionSize
		instr, err := DecodeInstruction(buf[offset : offset+InstructionSize])
		if err != nil {
			return nil, fmt.Errorf("failed to decode instruction %d: %w", i, err)
		}
		instrs[i] = *instr
	}

	return instrs, nil
}
