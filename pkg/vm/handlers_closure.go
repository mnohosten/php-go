package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Closure Opcode Handlers
// ============================================================================

// opDeclareLambdaFunction handles OpDeclareLambdaFunction
// Creates an anonymous function/closure
// ExtendedValue: number of parameters
// Op1: closure name (unique identifier from constants)
// Op2: closure start position
// Result: closure end position (also stores the closure object)
func (vm *VM) opDeclareLambdaFunction(frame *Frame, instr Instruction) error {
	// Get closure name
	funcNameVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	funcName := funcNameVal.ToString()

	// Get closure start and end positions
	funcStartVal, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	funcStart := int(funcStartVal.ToInt())

	funcEndVal, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}
	funcEnd := int(funcEndVal.ToInt())

	// Extract closure instructions from the main instruction stream
	if funcStart < 0 || funcEnd > len(frame.fn.Instructions) || funcStart >= funcEnd {
		return fmt.Errorf("Invalid closure bounds: start=%d, end=%d", funcStart, funcEnd)
	}

	closureInstructions := make(Instructions, funcEnd-funcStart)
	copy(closureInstructions, frame.fn.Instructions[funcStart:funcEnd])

	// CRITICAL: Adjust jump targets to be relative to closure start
	for i := range closureInstructions {
		instr := &closureInstructions[i]

		// Check if this is a jump instruction
		if instr.Opcode == OpJmp || instr.Opcode == OpJmpZ || instr.Opcode == OpJmpNZ {
			var targetOperand *Operand
			if instr.Opcode == OpJmp {
				targetOperand = &instr.Op1
			} else {
				targetOperand = &instr.Op2
			}

			// If it's a constant operand with absolute position
			if targetOperand.Type == OpConst {
				if int(targetOperand.Value) < len(vm.constants) {
					if targetVal, ok := vm.constants[targetOperand.Value].(int64); ok {
						absoluteTarget := int(targetVal)
						// Convert to relative position
						relativeTarget := absoluteTarget - funcStart
						// Store relative target back in constants
						vm.constants[targetOperand.Value] = int64(relativeTarget)
					}
				}
			}
		}
	}

	// Create CompiledFunction for the closure
	numParams := int(instr.ExtendedValue)
	compiledFunc := &CompiledFunction{
		Name:         funcName,
		Instructions: closureInstructions,
		NumParams:    numParams,
		NumLocals:    numParams + 10, // Parameters + space for locals/temps
	}

	// Register the closure function
	vm.RegisterFunction(funcName, compiledFunc)

	// Check if this should be a static closure
	// We can infer from the context or add flags if needed
	isStatic := false // For now, closures default to non-static

	// Create closure object
	var closure *Closure
	if isStatic {
		closure = NewStaticClosure(compiledFunc)
	} else {
		closure = NewClosure(compiledFunc)

		// If not static, inherit $this from current frame
		if frame.thisObject != nil {
			// The closure captures the current $this
			closure.BindVariable("this", types.NewObject(frame.thisObject), false)
		}
	}

	// Wrap closure in a resource value
	resource := types.NewResourceHandle("Closure", closure)
	closureVal := types.NewResource(resource)

	// Store the closure object in temp 0 for OpBindLexical to use
	return vm.setOperandValue(frame, Operand{Type: OpTmpVar, Value: 0}, closureVal)
}

// opBindLexical handles OpBindLexical
// Binds a lexical variable to a closure (use clause)
// Op1: variable name
// Op2: by-reference flag
// Result: closure object
func (vm *VM) opBindLexical(frame *Frame, instr Instruction) error {
	// Get variable name
	varNameVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	varName := varNameVal.ToString()

	// Get by-reference flag
	byRefVal, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	isByRef := byRefVal.ToInt() != 0

	// Get closure object
	closureVal, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	// Extract closure from resource
	if closureVal.Type() != types.TypeResource {
		return fmt.Errorf("OpBindLexical expects closure resource, got %s", closureVal.TypeString())
	}

	resource := closureVal.ToResource()
	data := resource.Data()
	closure, ok := data.(*Closure)
	if !ok {
		return fmt.Errorf("OpBindLexical expects Closure resource, got %T", data)
	}

	// Look up the variable in the current scope
	// Try locals first, then globals
	var varValue *types.Value
	found := false

	// Check frame locals by iterating through them
	// The compiler should have mapped variable names to local indices
	// For now, we'll look in globals
	if val, ok := vm.globals[varName]; ok {
		varValue = val
		found = true
	}

	if !found {
		// Variable doesn't exist - use null (PHP behavior)
		varValue = types.NewNull()
	}

	// Bind the variable to the closure
	closure.BindVariable(varName, varValue, isByRef)

	return nil
}

// opDeclareFunction handles OpDeclareFunction
// Declares a named function
// ExtendedValue: number of parameters
// Op1: function name (constant index)
// Op2: function start position (constant index)
// Result: function end position (constant index)
func (vm *VM) opDeclareFunction(frame *Frame, instr Instruction) error {
	// Get function name
	funcNameVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	funcName := funcNameVal.ToString()

	// Check if function is already registered
	if _, ok := vm.functions[funcName]; ok {
		// Function already exists
		// In PHP, this would be a fatal error in most contexts
		// but we'll just skip re-registration
		return nil
	}

	// Get function start and end positions
	funcStartVal, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	funcStart := int(funcStartVal.ToInt())

	funcEndVal, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}
	funcEnd := int(funcEndVal.ToInt())

	// Extract function instructions from the main instruction stream
	if funcStart < 0 || funcEnd > len(frame.fn.Instructions) || funcStart >= funcEnd {
		return fmt.Errorf("Invalid function bounds: start=%d, end=%d", funcStart, funcEnd)
	}

	funcInstructions := make(Instructions, funcEnd-funcStart)
	copy(funcInstructions, frame.fn.Instructions[funcStart:funcEnd])

	// CRITICAL FIX: Adjust jump targets to be relative to function start
	// When instructions are extracted, absolute positions need to be converted to relative positions
	for i := range funcInstructions {
		instr := &funcInstructions[i]

		// Check if this is a jump instruction (JMP, JMPZ, JMPNZ)
		if instr.Opcode == OpJmp || instr.Opcode == OpJmpZ || instr.Opcode == OpJmpNZ {
			// Get the jump target from the appropriate operand
			var targetOperand *Operand
			if instr.Opcode == OpJmp {
				targetOperand = &instr.Op1 // JMP uses Op1 for target
			} else {
				targetOperand = &instr.Op2 // JMPZ/JMPNZ use Op2 for target
			}

			// If it's a constant operand, it contains an absolute position
			if targetOperand.Type == OpConst {
				// Get the absolute target from constants
				if int(targetOperand.Value) < len(vm.constants) {
					if targetVal, ok := vm.constants[targetOperand.Value].(int64); ok {
						absoluteTarget := int(targetVal)
						// Convert to relative position
						relativeTarget := absoluteTarget - funcStart
						// Store the relative target back in constants
						vm.constants[targetOperand.Value] = int64(relativeTarget)
					}
				}
			}
		}
	}

	// Create CompiledFunction
	numParams := int(instr.ExtendedValue)
	compiledFunc := &CompiledFunction{
		Name:         funcName,
		Instructions: funcInstructions,
		NumParams:    numParams,
		NumLocals:    numParams + 10, // Parameters + space for local vars and temps
	}

	// Register the function
	vm.RegisterFunction(funcName, compiledFunc)

	return nil
}

// Helper function to invoke a closure
// This is called when a closure is used as a callable
func (vm *VM) invokeClosure(closure *Closure, args []*types.Value, thisObj *types.Object) (*types.Value, error) {
	// Create new frame for the closure
	frame := NewFrame(closure.Function)

	// Set $this if applicable
	if !closure.Static {
		if thisObj != nil {
			frame.thisObject = thisObj
		} else {
			// Check if closure captured $this
			if thisVal, ok := closure.GetCapturedVariable("this"); ok {
				if thisVal.Type() == types.TypeObject {
					frame.thisObject = thisVal.ToObject()
				}
			}
		}
	}

	// Set parameters
	for i, arg := range args {
		if i < closure.Function.NumParams {
			frame.setParam(i, arg)
		}
	}

	// Set up captured variables as locals
	// The compiler should have allocated local slots for captured variables
	// For a proper implementation, we need to know which local index
	// corresponds to which captured variable
	// For now, we'll store them in a way that the closure can access them
	// TODO: Integrate with compiler's variable allocation

	// Push frame and execute
	if err := vm.pushFrame(frame); err != nil {
		return nil, err
	}
	err := vm.runFrame(frame)
	if err != nil {
		vm.popFrame()
		return nil, err
	}

	// Get return value
	result := frame.getReturnValue()
	vm.popFrame()

	return result, nil
}

// GetClosureFromValue extracts a Closure from a Value
// Helper function for other parts of the VM
func GetClosureFromValue(val *types.Value) (*Closure, error) {
	if val.Type() != types.TypeResource {
		return nil, fmt.Errorf("Expected closure resource, got %s", val.TypeString())
	}

	resource := val.ToResource()
	data := resource.Data()
	closure, ok := data.(*Closure)
	if !ok {
		return nil, fmt.Errorf("Expected Closure resource, got %T", data)
	}

	return closure, nil
}
