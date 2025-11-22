package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

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

// opInitFcall initializes a regular function call
// Op2: function name (constant or variable)
// ExtendedValue: number of arguments
func (vm *VM) opInitFcall(frame *Frame, instr Instruction) error {
	// Get function name
	funcName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	funcNameStr := funcName.ToString()

	// Look up the function in VM's function registry
	fn, exists := vm.GetFunction(funcNameStr)
	if !exists {
		return fmt.Errorf("Call to undefined function %s()", funcNameStr)
	}

	// Store pending function call info in frame
	frame.pendingFunction = fn
	frame.pendingParams = &CallParams{
		params: make([]*types.Value, 0, int(instr.ExtendedValue)),
	}

	return nil
}

// opSendVal sends a parameter value for the pending function/method call
// Op1: parameter value
func (vm *VM) opSendVal(frame *Frame, instr Instruction) error {
	// Get the parameter value
	paramValue, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Add to pending parameters
	if frame.pendingParams == nil {
		frame.pendingParams = &CallParams{
			params: make([]*types.Value, 0, 8),
		}
	}
	frame.pendingParams.params = append(frame.pendingParams.params, paramValue)

	return nil
}

// opDoFcall executes a function or method call
// This handles both regular function calls (from OpInitFcall) and method calls (from OpInitMethodCall)
// Result: return value
func (vm *VM) opDoFcall(frame *Frame, instr Instruction) error {
	var fn *CompiledFunction
	var thisObj *types.Object
	var currentClass *types.ClassEntry
	var calledClass *types.ClassEntry

	// Check if this is a method call or regular function call
	if frame.pendingMethod != nil {
		// Method call - convert MethodDef to CompiledFunction
		fn = &CompiledFunction{
			Name:         frame.pendingMethod.Name,
			Instructions: convertInstructions(frame.pendingMethod.Instructions),
			NumLocals:    frame.pendingMethod.NumLocals,
			NumParams:    frame.pendingMethod.NumParams,
		}

		thisObj = frame.pendingObject
		if thisObj != nil && thisObj.ClassEntry != nil {
			currentClass = thisObj.ClassEntry
			calledClass = thisObj.ClassEntry
		}

		// Clear pending method
		frame.pendingMethod = nil
		frame.pendingObject = nil
	} else if frame.pendingFunction != nil {
		// Regular function call
		fn = frame.pendingFunction
		frame.pendingFunction = nil
	} else {
		return fmt.Errorf("DO_FCALL: no pending function or method call")
	}

	// Get parameters
	params := make([]*types.Value, 0)
	if frame.pendingParams != nil {
		params = frame.pendingParams.params
		frame.pendingParams = nil
	}

	// Create new frame for the function/method
	newFrame := NewFrame(fn)

	// Set object/class context for methods
	newFrame.thisObject = thisObj
	newFrame.currentClass = currentClass
	newFrame.calledClass = calledClass

	// Copy parameters to the new frame's local variables
	for i, param := range params {
		if i < fn.NumParams {
			newFrame.setParam(i, param)
			// Debug: log parameter assignment (disabled)
			// fmt.Printf("DEBUG: setParam(%d, %v) - type=%v\n", i, param, param.Type())
		}
	}

	// Push the new frame onto the call stack
	if err := vm.pushFrame(newFrame); err != nil {
		return err
	}

	// Execute the function immediately in this context
	// The function will run until it returns or hits an error
	err := vm.runFrame(newFrame)
	if err != nil {
		return err
	}

	// Pop the completed frame
	completedFrame := vm.popFrame()

	// Store the return value in the result operand
	returnValue := completedFrame.getReturnValue()
	if instr.Result.Type != OpUnused {
		return vm.setOperandValue(frame, instr.Result, returnValue)
	}

	return nil
}

// opDoUcall executes a user-defined function call (same as OpDoFcall)
func (vm *VM) opDoUcall(frame *Frame, instr Instruction) error {
	return vm.opDoFcall(frame, instr)
}

// opDoIcall executes an internal (built-in) function call
func (vm *VM) opDoIcall(frame *Frame, instr Instruction) error {
	// For now, treat the same as regular function call
	// In the future, this could be optimized for built-in functions
	return vm.opDoFcall(frame, instr)
}

// ============================================================================
// Helper Functions
// ============================================================================

// convertInstructions converts []interface{} to Instructions
func convertInstructions(instrArray []interface{}) Instructions {
	if instrArray == nil {
		return Instructions{}
	}

	result := make(Instructions, 0, len(instrArray))
	for _, instr := range instrArray {
		if i, ok := instr.(Instruction); ok {
			result = append(result, i)
		}
	}
	return result
}
