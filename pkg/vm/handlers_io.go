package vm

// ============================================================================
// I/O Opcode Handlers
// ============================================================================

// opEcho handles echo statement
func (vm *VM) opEcho(frame *Frame, instr Instruction) error {
	// Get the value to echo
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Convert to string and write to output
	output := value.ToString()
	vm.writeOutput([]byte(output))

	return nil
}

// opPrint handles print statement (same as echo but returns 1)
func (vm *VM) opPrint(frame *Frame, instr Instruction) error {
	// Echo the value
	if err := vm.opEcho(frame, instr); err != nil {
		return err
	}

	// Print returns 1
	// TODO: Set result to 1 if needed

	return nil
}
