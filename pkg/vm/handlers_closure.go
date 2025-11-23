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
// Op1: function name (or unique identifier)
// Result: closure object
func (vm *VM) opDeclareLambdaFunction(frame *Frame, instr Instruction) error {
	// Get function name/identifier
	funcNameVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	funcName := funcNameVal.ToString()

	// Look up the compiled function
	fn, ok := vm.functions[funcName]
	if !ok {
		return fmt.Errorf("Lambda function '%s' not found", funcName)
	}

	// Check if this should be a static closure (from ExtendedValue flag)
	isStatic := (instr.ExtendedValue & 1) != 0

	// Create closure
	var closure *Closure
	if isStatic {
		closure = NewStaticClosure(fn)
	} else {
		closure = NewClosure(fn)

		// If not static, inherit $this from current frame
		if frame.thisObject != nil {
			// The closure captures the current $this
			// This will be available when the closure is called
			// We store it as a special captured variable
			closure.BindVariable("this", types.NewObject(frame.thisObject), false)
		}
	}

	// Wrap closure in a value
	// For now, we'll use a resource to store the closure
	resource := types.NewResourceHandle("Closure", closure)
	closureVal := types.NewResource(resource)

	return vm.setOperandValue(frame, instr.Result, closureVal)
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
// Declares a named function (not a closure, but similar handling)
// Op1: function name
// Result: (optional) function value
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

	// The function should have been added to vm.functions during compilation
	// This opcode is mainly for runtime declaration in certain contexts
	// For now, this is a no-op

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
	vm.pushFrame(frame)
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
