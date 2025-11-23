package vm

import "github.com/krizos/php-go/pkg/types"

// ============================================================================
// Control Flow Opcode Handlers
// ============================================================================

// opJmp handles unconditional jump
func (vm *VM) opJmp(frame *Frame, instr Instruction) error {
	// Op1 contains the jump target (may be CONST or direct value)
	var target int
	if instr.Op1.Type == OpConst {
		// Resolve constant
		targetVal, err := vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}
		target = int(targetVal.ToInt())
	} else {
		// Direct value
		target = int(instr.Op1.Value)
	}
	frame.ip = target
	return nil
}

// opJmpZ handles jump if zero (false)
func (vm *VM) opJmpZ(frame *Frame, instr Instruction) error {
	// Op1 contains the condition
	condition, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// If condition is false, jump
	if !condition.ToBool() {
		// Op2 contains the jump target
		// Due to historical reasons, this may be a CONST operand where the value is actually
		// a direct instruction index, not a constant pool index
		// We try to resolve it as a constant first, and if that fails, use it directly
		var target int
		if instr.Op2.Type == OpConst {
			targetVal, err := vm.getOperandValue(frame, instr.Op2)
			if err == nil {
				// Successfully resolved as constant
				target = int(targetVal.ToInt())
			} else {
				// Not in constant pool - use as direct value
				target = int(instr.Op2.Value)
			}
		} else {
			// Direct value
			target = int(instr.Op2.Value)
		}
		frame.ip = target
	}

	return nil
}

// opJmpNZ handles jump if not zero (true)
func (vm *VM) opJmpNZ(frame *Frame, instr Instruction) error {
	// Op1 contains the condition
	condition, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// If condition is true, jump
	if condition.ToBool() {
		// Op2 contains the jump target
		// Due to historical reasons, this may be a CONST operand where the value is actually
		// a direct instruction index, not a constant pool index
		// We try to resolve it as a constant first, and if that fails, use it directly
		var target int
		if instr.Op2.Type == OpConst {
			targetVal, err := vm.getOperandValue(frame, instr.Op2)
			if err == nil {
				// Successfully resolved as constant
				target = int(targetVal.ToInt())
			} else {
				// Not in constant pool - use as direct value
				target = int(instr.Op2.Value)
			}
		} else {
			// Direct value
			target = int(instr.Op2.Value)
		}
		frame.ip = target
	}

	return nil
}

// opExit handles exit() and die() language constructs
// Op1: optional exit message/status code
func (vm *VM) opExit(frame *Frame, instr Instruction) error {
	// Get the exit argument (if any)
	if instr.Op1.Type != OpUnused {
		value, err := vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}

		// If it's a string, output it
		// If it's an integer, set exit code
		if value.Type() == types.TypeString {
			vm.output = append(vm.output, []byte(value.ToString())...)
			vm.exitCode = 0
		} else if value.Type() == types.TypeInt {
			vm.exitCode = int(value.ToInt())
		} else {
			// Convert to string and output
			vm.output = append(vm.output, []byte(value.ToString())...)
			vm.exitCode = 0
		}
	} else {
		// No argument, just exit with code 0
		vm.exitCode = 0
	}

	// Set the exit flag to stop execution
	vm.exited = true
	return nil
}
