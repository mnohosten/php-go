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

	isOnce := (typeVal == 1 || typeVal == 3) // include_once or require_once
	isRequire := (typeVal == 2 || typeVal == 3) // require or require_once

	// TODO: Full implementation requires breaking circular dependency
	// between VM and compiler packages. For now, we implement basic
	// path resolution and tracking for _once variants.
	//
	// Full implementation should:
	// 1. Use runtime.GetGlobalIncludeManager() to resolve path
	// 2. Read and parse the file
	// 3. Compile it to bytecode
	// 4. Execute in a new frame
	// 5. Return the result value
	//
	// This requires either:
	// - Moving compilation logic out of compiler package
	// - Using dependency injection for the compiler
	// - Creating an execution engine layer above VM and compiler

	// For now, just do basic path validation and tracking
	_ = path
	_ = isOnce
	_ = isRequire

	// Return 1 (success) as placeholder
	// In a full implementation, this would be the return value from the included file
	if err := vm.setOperandValue(frame, instr.Result, types.NewInt(1)); err != nil {
		return err
	}

	return nil
}
