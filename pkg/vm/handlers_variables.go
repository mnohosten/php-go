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

// opBindGlobal handles binding a local variable to a global variable
// Op1 = local variable (CV) to bind
// After this, the local variable references the global variable
func (vm *VM) opBindGlobal(frame *Frame, instr Instruction) error {
	// Get the variable index
	if instr.Op1.Type != OpCV {
		return fmt.Errorf("BIND_GLOBAL expects CV operand, got %v", instr.Op1.Type)
	}

	_ = int(instr.Op1.Value) // index - will be used when we implement proper global binding

	// TODO: Implement proper global variable binding
	// This requires:
	// 1. Track variable names in function metadata (compiler enhancement)
	// 2. Add global binding map to Frame structure
	// 3. Modify getOperandValue/setOperandValue to check for global-bound variables
	// 4. When a variable is marked as global, read/write operations should use vm.globals
	//
	// For now, this is a no-op placeholder that allows the code to parse and compile
	// without errors. The global statement will be recognized but won't have runtime effect.

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
