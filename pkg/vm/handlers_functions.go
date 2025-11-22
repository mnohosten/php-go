package vm

// ============================================================================
// Function Opcode Handlers
// ============================================================================

// opReturn handles function return
func (vm *VM) opReturn(frame *Frame, instr Instruction) error {
	// Get return value (Op1)
	returnValue, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Set frame's return value
	frame.setReturnValue(returnValue)

	// Set IP to end of instructions to exit frame
	frame.ip = len(frame.fn.Instructions)

	return nil
}

// opInitFcall initializes a function call (placeholder for Phase 3)
func (vm *VM) opInitFcall(frame *Frame, instr Instruction) error {
	// TODO: Implement in full function call support
	return nil
}

// opSendVal sends a value parameter (placeholder)
func (vm *VM) opSendVal(frame *Frame, instr Instruction) error {
	// TODO: Implement parameter passing
	return nil
}

// opDoFcall executes a function call (placeholder)
func (vm *VM) opDoFcall(frame *Frame, instr Instruction) error {
	// TODO: Implement function calling
	return nil
}
