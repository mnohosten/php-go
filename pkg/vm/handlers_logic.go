package vm

import "github.com/krizos/php-go/pkg/types"

// ============================================================================
// Logic & Bitwise Opcode Handlers
// ============================================================================

// opBoolNot handles boolean NOT (!)
func (vm *VM) opBoolNot(frame *Frame, instr Instruction) error {
	operand, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	result := types.NewBool(!operand.ToBool())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opBWNot handles bitwise NOT (~)
func (vm *VM) opBWNot(frame *Frame, instr Instruction) error {
	operand, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Bitwise NOT on integer
	result := types.NewInt(^operand.ToInt())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opBWAnd handles bitwise AND (&)
func (vm *VM) opBWAnd(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewInt(left.ToInt() & right.ToInt())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opBWOr handles bitwise OR (|)
func (vm *VM) opBWOr(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewInt(left.ToInt() | right.ToInt())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opBWXor handles bitwise XOR (^)
func (vm *VM) opBWXor(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewInt(left.ToInt() ^ right.ToInt())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opShiftLeft handles left shift (<<)
func (vm *VM) opShiftLeft(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewInt(left.ToInt() << uint(right.ToInt()))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opShiftRight handles right shift (>>)
func (vm *VM) opShiftRight(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewInt(left.ToInt() >> uint(right.ToInt()))

	return vm.setOperandValue(frame, instr.Result, result)
}
