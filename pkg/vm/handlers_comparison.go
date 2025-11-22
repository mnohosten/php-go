package vm

import "github.com/krizos/php-go/pkg/types"

// ============================================================================
// Comparison Opcode Handlers
// ============================================================================

// opIsEqual handles loose equality (==)
func (vm *VM) opIsEqual(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewBool(left.Equals(right))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opIsNotEqual handles loose inequality (!=)
func (vm *VM) opIsNotEqual(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewBool(!left.Equals(right))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opIsIdentical handles strict equality (===)
func (vm *VM) opIsIdentical(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewBool(left.Identical(right))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opIsNotIdentical handles strict inequality (!==)
func (vm *VM) opIsNotIdentical(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	result := types.NewBool(!left.Identical(right))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opIsSmaller handles less than (<)
func (vm *VM) opIsSmaller(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// PHP comparison: convert to numbers if needed
	var result bool

	if left.IsString() && right.IsString() {
		// String comparison
		result = left.ToString() < right.ToString()
	} else {
		// Numeric comparison
		result = left.ToFloat() < right.ToFloat()
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// opIsSmallerOrEqual handles less than or equal (<=)
func (vm *VM) opIsSmallerOrEqual(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result bool

	if left.IsString() && right.IsString() {
		result = left.ToString() <= right.ToString()
	} else {
		result = left.ToFloat() <= right.ToFloat()
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// opSpaceship handles spaceship operator (<=>)
// Returns -1 if left < right, 0 if equal, 1 if left > right
func (vm *VM) opSpaceship(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result int64

	if left.IsString() && right.IsString() {
		leftStr := left.ToString()
		rightStr := right.ToString()
		if leftStr < rightStr {
			result = -1
		} else if leftStr > rightStr {
			result = 1
		} else {
			result = 0
		}
	} else {
		leftNum := left.ToFloat()
		rightNum := right.ToFloat()
		if leftNum < rightNum {
			result = -1
		} else if leftNum > rightNum {
			result = 1
		} else {
			result = 0
		}
	}

	return vm.setOperandValue(frame, instr.Result, types.NewInt(result))
}
