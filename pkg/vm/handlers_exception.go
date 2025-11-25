package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/runtime"
)

// Note: fmt is still used for error formatting in handleException

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
// Called when findCatchHandler has already determined this is the right catch block
// Op1: exception type constant index (used by findCatchHandler)
// Result: caught exception (temp var)
func (vm *VM) opCatch(frame *Frame, instr Instruction) error {
	// Check if there's a pending exception
	if frame.exception == nil {
		// No exception to catch - this happens when catch block is reached
		// via normal control flow (not exception), skip the catch body
		return nil
	}

	// Get the exception
	exc, ok := frame.exception.(*runtime.Exception)
	if !ok || exc == nil {
		// No valid exception to catch
		return nil
	}

	// Exception matches (findCatchHandler already verified this)
	// Store exception in result operand
	excVal := runtime.WrapExceptionAsValue(exc)
	if err := vm.setOperandValue(frame, instr.Result, excVal); err != nil {
		return err
	}

	// Clear exception from frame
	frame.exception = nil

	return nil
}

// opFastCall implements the FAST_CALL opcode for finally block setup
// Op1: target address (patched by compiler to point to finally block)
func (vm *VM) opFastCall(frame *Frame, instr Instruction) error {
	// FAST_CALL is used to set up a finally block call
	// Op1 contains the target address of the finally block
	// We need to save the return address so FAST_RET can return here

	// Get finally block address from Op1
	finallyAddr := int(instr.Op1.Value)

	// Save current IP as return address (so we can continue after finally)
	// We push the return address onto a simple stack
	if frame.finallyStack == nil {
		frame.finallyStack = make([]int, 0, 4)
	}
	frame.finallyStack = append(frame.finallyStack, frame.ip)

	// Jump to finally block
	frame.ip = finallyAddr

	return nil
}

// opFastRet implements the FAST_RET opcode for returning from finally block
func (vm *VM) opFastRet(frame *Frame, instr Instruction) error {
	// FAST_RET returns from the finally block to where we came from

	// Pop return address from finally stack
	if frame.finallyStack == nil || len(frame.finallyStack) == 0 {
		// No return address - just continue normally
		// This happens when finally is reached via normal control flow
		return nil
	}

	// Pop return address
	n := len(frame.finallyStack)
	returnAddr := frame.finallyStack[n-1]
	frame.finallyStack = frame.finallyStack[:n-1]

	// Jump back to the return address
	frame.ip = returnAddr

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

// findCatchHandler looks for a CATCH instruction in the current frame that matches the exception type
func (vm *VM) findCatchHandler(frame *Frame) bool {
	exc, ok := frame.exception.(*runtime.Exception)
	if !ok || exc == nil {
		return false
	}

	// Scan instructions from current IP
	for i := frame.ip; i < len(frame.fn.Instructions); i++ {
		instr := frame.fn.Instructions[i]
		if instr.Opcode == OpCatch {
			// Get catch type from Op1 (constant index)
			catchTypeName := ""
			if instr.Op1.Type == OpConst {
				typeConst, err := vm.GetConstant(int(instr.Op1.Value))
				if err == nil && typeConst != nil {
					catchTypeName = typeConst.ToString()
				}
			}

			// Check if exception matches catch type
			if vm.exceptionMatchesCatchType(exc, catchTypeName) {
				// Found matching catch handler
				frame.ip = i
				return true
			}
			// Type doesn't match, continue searching for next catch
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

// exceptionMatchesCatchType checks if an exception matches a catch type
// Handles inheritance: Exception is caught by catch(Exception) or catch(Throwable)
func (vm *VM) exceptionMatchesCatchType(exc *runtime.Exception, catchTypeName string) bool {
	// Empty catch type catches everything
	if catchTypeName == "" {
		return true
	}

	// Exact match
	if exc.ClassName == catchTypeName {
		return true
	}

	// Check if exception class inherits from catch type
	// Look up the exception's class definition
	exceptionClass, exists := vm.classes[exc.ClassName]
	if !exists {
		// Unknown class - check basic hierarchy
		// Exception and Error are both Throwable
		if catchTypeName == "Throwable" {
			return true
		}
		if catchTypeName == "Exception" && (exc.ClassName == "Exception" ||
			exc.ClassName == "RuntimeException" ||
			exc.ClassName == "LogicException" ||
			exc.ClassName == "InvalidArgumentException" ||
			exc.ClassName == "OutOfBoundsException") {
			return true
		}
		return false
	}

	// Walk up the class hierarchy
	currentClass := exceptionClass
	for currentClass != nil {
		if currentClass.Name == catchTypeName {
			return true
		}
		currentClass = currentClass.ParentClass
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
