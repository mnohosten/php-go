package vm

import (
	"fmt"
	"math"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Arithmetic Opcode Handlers
// ============================================================================

// opAdd handles addition (result = op1 + op2)
func (vm *VM) opAdd(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// PHP addition rules:
	// - If both are ints, result is int
	// - If either is float, result is float
	// - Otherwise convert to numeric

	var result *types.Value

	if left.IsInt() && right.IsInt() {
		result = types.NewInt(left.ToInt() + right.ToInt())
	} else {
		// At least one is float, or needs conversion
		result = types.NewFloat(left.ToFloat() + right.ToFloat())
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opSub handles subtraction (result = op1 - op2)
func (vm *VM) opSub(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result *types.Value

	if left.IsInt() && right.IsInt() {
		result = types.NewInt(left.ToInt() - right.ToInt())
	} else {
		result = types.NewFloat(left.ToFloat() - right.ToFloat())
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opMul handles multiplication (result = op1 * op2)
func (vm *VM) opMul(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result *types.Value

	if left.IsInt() && right.IsInt() {
		result = types.NewInt(left.ToInt() * right.ToInt())
	} else {
		result = types.NewFloat(left.ToFloat() * right.ToFloat())
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opDiv handles division (result = op1 / op2)
func (vm *VM) opDiv(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Division by zero check
	if right.ToFloat() == 0.0 {
		return fmt.Errorf("division by zero")
	}

	// Division always returns float in PHP
	result := types.NewFloat(left.ToFloat() / right.ToFloat())

	return vm.setOperandValue(frame, instr.Result, result)
}

// opMod handles modulo (result = op1 % op2)
func (vm *VM) opMod(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Convert to integers for modulo
	leftInt := left.ToInt()
	rightInt := right.ToInt()

	// Modulo by zero check
	if rightInt == 0 {
		return fmt.Errorf("modulo by zero")
	}

	// Modulo returns int
	result := types.NewInt(leftInt % rightInt)

	return vm.setOperandValue(frame, instr.Result, result)
}

// opPow handles exponentiation (result = op1 ** op2)
func (vm *VM) opPow(frame *Frame, instr Instruction) error {
	left, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	right, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Convert to floats for power
	base := left.ToFloat()
	exp := right.ToFloat()

	// Calculate power
	result := types.NewFloat(math.Pow(base, exp))

	return vm.setOperandValue(frame, instr.Result, result)
}

// opNegate handles unary negation (result = -op1)
// This might be called OpNegate or mapped to OpSub with zero
func (vm *VM) opNegate(frame *Frame, instr Instruction) error {
	operand, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	var result *types.Value

	if operand.IsInt() {
		result = types.NewInt(-operand.ToInt())
	} else {
		result = types.NewFloat(-operand.ToFloat())
	}

	return vm.setOperandValue(frame, instr.Result, result)
}
