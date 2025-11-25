package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

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
// ExtendedValue = mode: 0 for isset, 1 for empty
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

	// Get the mode from ExtendedValue (0 = isset, 1 = empty)
	mode := instr.ExtendedValue
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

// opBindGlobal handles binding a local variable to a global variable
// Op1 = local variable index (CV)
// Op2 = constant index containing the variable name
// After this, the local variable references the global variable
func (vm *VM) opBindGlobal(frame *Frame, instr Instruction) error {
	// Get the variable index
	if instr.Op1.Type != OpCV {
		return fmt.Errorf("BIND_GLOBAL expects CV operand for Op1, got %v", instr.Op1.Type)
	}
	localIdx := int(instr.Op1.Value)

	// Get the variable name from constants
	if instr.Op2.Type != OpConst {
		return fmt.Errorf("BIND_GLOBAL expects Const operand for Op2, got %v", instr.Op2.Type)
	}
	nameIdx := int(instr.Op2.Value)
	if nameIdx < 0 || nameIdx >= len(vm.constants) {
		return fmt.Errorf("BIND_GLOBAL: invalid constant index %d", nameIdx)
	}
	varName, ok := vm.constants[nameIdx].(string)
	if !ok {
		return fmt.Errorf("BIND_GLOBAL: constant at index %d is not a string", nameIdx)
	}

	// Initialize globalBindings map if needed
	if frame.globalBindings == nil {
		frame.globalBindings = make(map[int]string)
	}

	// Record the binding: local index -> global variable name
	frame.globalBindings[localIdx] = varName

	// Ensure the global variable exists in vm.globals
	// If it doesn't exist, create it with null value
	if _, exists := vm.globals[varName]; !exists {
		vm.globals[varName] = types.NewNull()
	}

	// Copy current global value to local slot for consistency
	// (This makes sure the local slot has the current global value)
	if localIdx >= 0 && localIdx < len(frame.locals) {
		frame.locals[localIdx] = vm.globals[varName]
	}

	return nil
}

// opUnsetVar handles UNSET_VAR opcode - unset($var)
func (vm *VM) opUnsetVar(frame *Frame, instr Instruction) error {
	// Get the variable index
	if instr.Op1.Type != OpCV {
		return fmt.Errorf("UNSET_VAR expects CV operand, got %v", instr.Op1.Type)
	}

	varIdx := int(instr.Op1.Value)

	// Set the variable to undefined
	// In PHP, unset() removes the variable from the symbol table
	// We simulate this by setting it to Undef
	if varIdx >= 0 && varIdx < len(frame.locals) {
		frame.locals[varIdx] = types.NewUndef()
	}

	return nil
}
