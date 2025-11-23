package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/runtime"
)

// ============================================================================
// Exception Opcode Handlers
// ============================================================================

// opThrow handles OpThrow
// Throws an exception and unwinds the stack
// Op1: exception value (object or resource)
func (vm *VM) opThrow(frame *Frame, instr Instruction) error {
	// Get exception value
	excVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Extract exception from value
	exc, err := runtime.GetExceptionFromValue(excVal)
	if err != nil {
		// If it's not a proper exception, create a generic one
		exc = runtime.NewException(excVal.ToString(), 0, nil)
	}

	// Set file and line information
	exc.File = frame.fn.Name
	exc.Line = int(instr.Lineno)

	// Generate stack trace
	trace := vm.generateStackTrace()
	exc.Trace = trace

	// Start exception handling - unwind stack
	return vm.handleException(exc)
}

// opCatch handles OpCatch
// Catches an exception if it matches the specified type
// ExtendedValue: exception type index (constant)
// Result: caught exception (temp var)
func (vm *VM) opCatch(frame *Frame, instr Instruction) error {
	// Check if there's a pending exception
	if frame.exception == nil {
		// No exception to catch
		return nil
	}

	// Get exception type if specified
	var exceptionType string
	if instr.ExtendedValue != 0 {
		typeConst, err := vm.GetConstant(int(instr.ExtendedValue))
		if err == nil {
			exceptionType = typeConst.ToString()
		}
	}

	// Check if exception matches type (if specified)
	// Need to assert frame.exception to *runtime.Exception
	exc, ok := frame.exception.(*runtime.Exception)
	if !ok || exc == nil {
		// No exception to catch
		return nil
	}

	if exceptionType != "" && exc.ClassName != exceptionType {
		// Type doesn't match, continue unwinding
		return nil
	}

	// Exception matches - catch it
	// Store exception in result operand
	excVal := runtime.WrapExceptionAsValue(exc)
	if err := vm.setOperandValue(frame, instr.Result, excVal); err != nil {
		return err
	}

	// Clear exception from frame
	frame.exception = nil

	return nil
}

// ============================================================================
// Exception Handling Helpers
// ============================================================================

// handleException handles an exception by unwinding the stack
func (vm *VM) handleException(exc *runtime.Exception) error {
	// Unwind stack looking for a catch handler
	for vm.frameIndex >= 0 {
		frame := vm.currentFrame()

		// Set exception in frame
		frame.exception = exc

		// Look for CATCH instruction in current frame
		// Scan forward from current IP to find OpCatch
		foundCatch := vm.findCatchHandler(frame)
		if foundCatch {
			// Found a catch handler in this frame
			return nil
		}

		// No catch handler in this frame, pop and continue unwinding
		vm.popFrame()
	}

	// No catch handler found - uncaught exception
	return fmt.Errorf("Uncaught %s: %s in %s:%d\n%s",
		exc.ClassName,
		exc.Message,
		exc.File,
		exc.Line,
		exc.GetTraceAsString())
}

// findCatchHandler looks for a CATCH instruction in the current frame
func (vm *VM) findCatchHandler(frame *Frame) bool {
	// Scan instructions from current IP
	for i := frame.ip; i < len(frame.fn.Instructions); i++ {
		instr := frame.fn.Instructions[i]
		if instr.Opcode == OpCatch {
			// Found catch handler
			// Set IP to catch instruction
			frame.ip = i
			return true
		}

		// If we hit certain opcodes, we've gone past the try-catch block
		if instr.Opcode == OpReturn ||
			instr.Opcode == OpReturnByRef ||
			instr.Opcode == OpGeneratorReturn {
			break
		}
	}

	return false
}

// generateStackTrace generates a stack trace from the current call stack
func (vm *VM) generateStackTrace() *runtime.StackTrace {
	trace := runtime.NewStackTrace()

	// Walk the frame stack
	for i := vm.frameIndex; i >= 0; i-- {
		frame := vm.frames[i]
		if frame == nil {
			continue
		}

		stackFrame := &runtime.StackFrame{
			File:     frame.fn.Name,
			Line:     int(frame.lastLine),
			Function: frame.fn.Name,
		}

		// Add class context if available
		if frame.classEntry != nil {
			stackFrame.Class = frame.classEntry.Name
			stackFrame.Type = "->"
		}

		trace.AddFrame(stackFrame)
	}

	return trace
}

// ============================================================================
// Frame Exception Support
// ============================================================================

// SetException sets the exception for a frame
func (f *Frame) SetException(exc *runtime.Exception) {
	f.exception = exc
}

// GetException gets the exception from a frame
func (f *Frame) GetException() *runtime.Exception {
	if f.exception == nil {
		return nil
	}
	exc, _ := f.exception.(*runtime.Exception)
	return exc
}

// HasException checks if a frame has an exception
func (f *Frame) HasException() bool {
	return f.exception != nil
}

// ClearException clears the exception from a frame
func (f *Frame) ClearException() {
	f.exception = nil
}
