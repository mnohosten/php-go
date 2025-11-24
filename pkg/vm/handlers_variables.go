package vm

import "github.com/krizos/php-go/pkg/types"

// ============================================================================
// Variable Opcode Handlers
// ============================================================================

// opConst loads a constant value
func (vm *VM) opConst(frame *Frame, instr Instruction) error {
	// Op1 contains the constant index
	value, err := vm.GetConstant(int(instr.Op1.Value))
	if err != nil {
		return err
	}

	return vm.setOperandValue(frame, instr.Result, value)
}

// opAssign handles variable assignment
func (vm *VM) opAssign(frame *Frame, instr Instruction) error {
	// Get the value to assign (from Op1 - this is where compiler puts it)
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Assign to result
	return vm.setOperandValue(frame, instr.Result, value)
}

// opQMAssign handles quick assign (no side effects): result = op1
// This is used for simple value assignments without modifications
func (vm *VM) opQMAssign(frame *Frame, instr Instruction) error {
	// Get the value from Op1
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Assign to result
	return vm.setOperandValue(frame, instr.Result, value)
}

// opFetch handles variable fetch (read)
func (vm *VM) opFetch(frame *Frame, instr Instruction) error {
	// Get the variable value
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Store in result
	return vm.setOperandValue(frame, instr.Result, value)
}

// opFree handles freeing temporary variables
// In Go with garbage collection, this is a no-op
// PHP's VM uses this to clean up reference counts, but we don't need it
func (vm *VM) opFree(frame *Frame, instr Instruction) error {
	// No-op: Go's garbage collector handles memory management
	// Attempting to clear the variable here causes issues with the compiler
	// which sometimes emits FREE for compiled variables
	return nil
}

// opUnset handles unsetting a variable
func (vm *VM) opUnset(frame *Frame, instr Instruction) error {
	// Set variable to null/undef
	return vm.setOperandValue(frame, instr.Op1, types.NewUndef())
}

// opIsset handles isset() check
func (vm *VM) opIsset(frame *Frame, instr Instruction) error {
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// isset returns false for null and undef
	result := !value.IsNull() && !value.IsUndef()

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// opEmpty handles empty() check
func (vm *VM) opEmpty(frame *Frame, instr Instruction) error {
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// empty() returns true for falsy values
	result := value.IsFalse()

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// opIssetIsemptyVar handles isset/empty check on a variable
// Op1 = variable to check (CV or temp)
// Op2 = mode: 0 for isset, 1 for empty
// Result = boolean result
func (vm *VM) opIssetIsemptyVar(frame *Frame, instr Instruction) error {
	// For isset/empty, we need to directly check the local variable
	// because getOperandValue converts nil to NewNull(), which hides
	// the difference between uninitialized and null variables

	var value *types.Value
	var index int

	// Calculate the actual index in the locals array
	switch instr.Op1.Type {
	case OpCV, OpVar:
		index = int(instr.Op1.Value)
	case OpTmpVar:
		offset := frame.fn.NumParams
		if frame.fn.NumCVs > 0 {
			offset = frame.fn.NumCVs
		}
		index = int(instr.Op1.Value) + offset
	default:
		// For other operand types, use normal value retrieval
		var err error
		value, err = vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}
		index = -1
	}

	// If we have a valid index, check the raw local value
	if index >= 0 {
		if index < len(frame.locals) {
			value = frame.locals[index]
		} else {
			value = nil
		}
	}

	// Get the mode from Op2 (0 = isset, 1 = empty)
	mode := instr.Op2.Value
	var result bool

	if mode == 0 {
		// isset mode: returns false for nil (uninitialized), null, and undef
		if value == nil {
			result = false
		} else {
			result = !value.IsNull() && !value.IsUndef()
		}
	} else {
		// empty mode: returns true for nil (uninitialized), falsy values, null, or undef
		if value == nil {
			result = true
		} else {
			result = value.IsFalse() || value.IsNull() || value.IsUndef()
		}
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}
