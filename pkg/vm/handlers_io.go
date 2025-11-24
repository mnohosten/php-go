package vm

import "github.com/krizos/php-go/pkg/types"

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

// opIncludeOrEval handles include, include_once, require, require_once
func (vm *VM) opIncludeOrEval(frame *Frame, instr Instruction) error {
	// Get the path from Op1
	pathValue, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the include type from Op2
	// 0 = include, 1 = include_once, 2 = require, 3 = require_once
	includeType, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	path := pathValue.ToString()
	typeVal := includeType.ToInt()

	// TODO: Implement actual file inclusion logic
	// For now, we'll emit a warning and return true
	// This allows parsing to succeed without full implementation

	var typeStr string
	switch typeVal {
	case 0:
		typeStr = "include"
	case 1:
		typeStr = "include_once"
	case 2:
		typeStr = "require"
	case 3:
		typeStr = "require_once"
	default:
		typeStr = "unknown"
	}

	// Emit a warning (not implemented yet)
	// For now, just log to stderr
	_ = typeStr
	_ = path

	// Return true (1) to indicate success
	// In PHP, include/require return 1 on success
	if err := vm.setOperandValue(frame, instr.Result, types.NewInt(1)); err != nil {
		return err
	}

	return nil
}
