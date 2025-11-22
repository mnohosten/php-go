package vm

// ============================================================================
// Control Flow Opcode Handlers
// ============================================================================

// opJmp handles unconditional jump
func (vm *VM) opJmp(frame *Frame, instr Instruction) error {
	// Op1 contains the jump target
	target := int(instr.Op1.Value)
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
		target := int(instr.Op2.Value)
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
		target := int(instr.Op2.Value)
		frame.ip = target
	}

	return nil
}
