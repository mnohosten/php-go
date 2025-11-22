package vm

import "github.com/krizos/php-go/pkg/types"

// ============================================================================
// String Opcode Handlers
// ============================================================================

// opConcat handles string concatenation
func (vm *VM) opConcat(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Convert both to strings and concatenate
	result := types.NewString(left.ToString() + right.ToString())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opFastConcat handles optimized string concatenation (same as opConcat for now)
func (vm *VM) opFastConcat(frame *Frame, instr Instruction) error {
	return vm.opConcat(frame, instr)
}
