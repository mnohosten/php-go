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

// opRecv receives a function parameter
// Op1: parameter index (constant)
// Result: CV where to store the parameter
func (vm *VM) opRecv(frame *Frame, instr Instruction) error {
	// Get parameter index
	paramIndexVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	paramIndex := int(paramIndexVal.ToInt())

	// Get parameter value from frame's locals
	// Parameters were already placed in locals by opDoFcall via setParam
	paramValue := frame.getLocal(paramIndex)
	if paramValue == nil {
		return fmt.Errorf("Missing required parameter at position %d", paramIndex)
	}

	// Store in result CV
	if instr.Result.Type != OpUnused {
		return vm.setOperandValue(frame, instr.Result, paramValue)
	}

	return nil
}

// opRecvInit receives a function parameter with a default value
// Op1: parameter index (constant)
// Op2: default value (from temp variable)
// Result: CV where to store the parameter
func (vm *VM) opRecvInit(frame *Frame, instr Instruction) error {
	// Get parameter index
	paramIndexVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	paramIndex := int(paramIndexVal.ToInt())

	// Try to get parameter value from frame's locals
	paramValue := frame.getLocal(paramIndex)

	// If parameter was not provided, use default value
	if paramValue == nil || paramValue.IsUndef() {
		// Get default value from Op2
		defaultValue, err := vm.getOperandValue(frame, instr.Op2)
		if err != nil {
			return err
		}
		paramValue = defaultValue
	}

	// Store in result CV
	if instr.Result.Type != OpUnused {
		return vm.setOperandValue(frame, instr.Result, paramValue)
	}

	return nil
}

// opInitFcallByName initializes a function call by name
// Op1: function name operand
// ExtendedValue: argument count
func (vm *VM) opInitFcallByName(frame *Frame, instr Instruction) error {
	// Get function name or callable
	funcOrName, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get argument count from ExtendedValue
	argCount := int64(instr.ExtendedValue)

	// Check if this is a closure (Resource type)
	if funcOrName.Type() == types.TypeResource {
		resource := funcOrName.ToResource()
		if closure, ok := resource.Data().(*Closure); ok {
			// It's a closure! Store it as pending closure
			frame.pendingClosure = closure
			frame.pendingFunction = closure.Function
			frame.pendingParams = &CallParams{
				params: make([]*types.Value, 0, int(argCount)),
			}
			return nil
		}
	}

	// Otherwise, treat as function name string
	funcNameStr := funcOrName.ToString()

	// Check if it's a built-in function
	if IsBuiltin(funcNameStr) {
		// Create a placeholder function for built-ins
		frame.pendingFunction = &CompiledFunction{
			Name: funcNameStr,
		}
		frame.pendingParams = &CallParams{
			params: make([]*types.Value, 0, int(argCount)),
		}
		return nil
	}

	// Look up the function in VM's function registry
	fn, exists := vm.GetFunction(funcNameStr)
	if !exists {
		return fmt.Errorf("Call to undefined function %s()", funcNameStr)
	}

	// Store pending function call info in frame
	frame.pendingFunction = fn
	frame.pendingParams = &CallParams{
		params: make([]*types.Value, 0, int(argCount)),
	}

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

	// Check if it's a built-in function
	if IsBuiltin(funcNameStr) {
		// Create a placeholder function for built-ins
		// We'll handle the actual call in opDoFcall
		frame.pendingFunction = &CompiledFunction{
			Name: funcNameStr,
		}
		frame.pendingParams = &CallParams{
			params: make([]*types.Value, 0, int(instr.ExtendedValue)),
		}
		return nil
	}

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
			byRef:  make([]bool, 0, 8),
		}
	}
	frame.pendingParams.params = append(frame.pendingParams.params, paramValue)
	frame.pendingParams.byRef = append(frame.pendingParams.byRef, false)

	return nil
}

// opSendRef sends a parameter by reference for the pending function/method call
// Op1: variable reference (the variable to pass by reference)
func (vm *VM) opSendRef(frame *Frame, instr Instruction) error {
	// For by-reference parameters, we need to pass a reference to the variable
	// so the called function can modify the original variable

	// Add to pending parameters
	if frame.pendingParams == nil {
		frame.pendingParams = &CallParams{
			params: make([]*types.Value, 0, 8),
			byRef:  make([]bool, 0, 8),
		}
	}

	// Get the operand - for references, we need to get a pointer to the actual
	// variable storage location, not just its value
	var refValue *types.Value

	switch instr.Op1.Type {
	case OpVar, OpCV:
		// Local variable - get a reference to it
		idx := int(instr.Op1.Value)
		if idx < len(frame.locals) {
			if frame.locals[idx] == nil {
				frame.locals[idx] = types.NewNull()
			}
			// Create a reference to the variable's value slot
			refValue = types.NewReference(frame.locals[idx])
			// Store the reference back so assignments go through it
			frame.locals[idx] = refValue
		} else {
			refValue = types.NewNull()
		}
	case OpTmpVar:
		// Temporary - get its value (less common for refs)
		val, err := vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}
		refValue = types.NewReference(val)
	default:
		// For other operand types, just get the value and wrap it
		val, err := vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}
		refValue = types.NewReference(val)
	}

	frame.pendingParams.params = append(frame.pendingParams.params, refValue)
	frame.pendingParams.byRef = append(frame.pendingParams.byRef, true)

	return nil
}

// opDoFcall executes a function or method call
// This handles both regular function calls (from OpInitFcall) and method calls (from OpInitMethodCall)
// Result: return value
func (vm *VM) opDoFcall(frame *Frame, instr Instruction) error {
	// Check for native constructor call (built-in exception classes)
	if frame.pendingNativeConstructor != "" {
		// Get parameters
		params := make([]*types.Value, 0)
		if frame.pendingParams != nil {
			params = frame.pendingParams.params
			frame.pendingParams = nil
		}

		// Handle built-in exception constructor
		obj := frame.pendingObject
		handleBuiltinExceptionConstructor(obj, params)

		// Clear pending state
		frame.pendingNativeConstructor = ""
		frame.pendingObject = nil
		return nil
	}

	var fn *CompiledFunction
	var closure *Closure
	var thisObj *types.Object
	var currentClass *types.ClassEntry
	var calledClass *types.ClassEntry
	var funcName string

	// Check if this is a method call, closure call, or regular function call
	if frame.pendingMethod != nil {
		// Method call - convert MethodDef to CompiledFunction
		fn = &CompiledFunction{
			Name:         frame.pendingMethod.Name,
			Instructions: convertInstructions(frame.pendingMethod.Instructions),
			NumLocals:    frame.pendingMethod.NumLocals,
			NumParams:    frame.pendingMethod.NumParams,
		}
		funcName = frame.pendingMethod.Name

		thisObj = frame.pendingObject
		if thisObj != nil && thisObj.ClassEntry != nil {
			currentClass = thisObj.ClassEntry
			calledClass = thisObj.ClassEntry
		}

		// Clear pending method
		frame.pendingMethod = nil
		frame.pendingObject = nil
	} else if frame.pendingClosure != nil {
		// Closure call
		closure = frame.pendingClosure
		fn = closure.Function
		funcName = fn.Name

		// Get $this from closure if available
		if !closure.Static {
			if thisVal, ok := closure.GetCapturedVariable("this"); ok {
				if thisVal.Type() == types.TypeObject {
					thisObj = thisVal.ToObject()
				}
			}
		}

		frame.pendingClosure = nil
		frame.pendingFunction = nil
	} else if frame.pendingFunction != nil {
		// Regular function call
		fn = frame.pendingFunction
		funcName = fn.Name
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

	// Check for special functions that need VM access
	if handled, returnValue, err := vm.handleSpecialFunction(funcName, params, frame); handled {
		if err != nil {
			return err
		}
		if instr.Result.Type != OpUnused {
			return vm.setOperandValue(frame, instr.Result, returnValue)
		}
		return nil
	}

	// Check if this is a built-in function
	if IsBuiltin(funcName) {
		returnValue, err := CallBuiltin(funcName, params)
		if err != nil {
			return err
		}
		if instr.Result.Type != OpUnused {
			return vm.setOperandValue(frame, instr.Result, returnValue)
		}
		return nil
	}

	// Create new frame for the function/method
	newFrame := NewFrame(fn)

	// Set object/class context for methods
	newFrame.thisObject = thisObj
	newFrame.currentClass = currentClass
	newFrame.calledClass = calledClass

	// If this is a closure, set up captured variables
	if closure != nil {
		// Copy captured variables to the new frame
		// We need to make them available as local variables
		// For now, we'll add them to the globals from the closure's perspective
		// A better approach would be to reserve CV slots, but this works for initial support
		for varName, varValue := range closure.CapturedVars {
			// Store in globals so the closure can access them
			// This is a simplified approach - a full implementation would use proper CV management
			vm.globals[varName] = varValue
		}
	}

	// Copy parameters to the new frame's local variables
	for i, param := range params {
		if i < fn.NumParams {
			newFrame.setParam(i, param)
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

// ============================================================================
// Special Functions (require VM access)
// ============================================================================

// handleSpecialFunction handles functions that require VM access
// Returns (handled bool, result *types.Value, error)
func (vm *VM) handleSpecialFunction(funcName string, params []*types.Value, frame *Frame) (bool, *types.Value, error) {
	switch funcName {
	case "define":
		return true, vm.funcDefine(params), nil
	case "defined":
		return true, vm.funcDefined(params), nil
	case "constant":
		return true, vm.funcConstant(params), nil
	case "call_user_func":
		return vm.funcCallUserFunc(params, frame)
	case "call_user_func_array":
		return vm.funcCallUserFuncArray(params, frame)
	case "func_get_args":
		return true, vm.funcGetArgs(frame), nil
	case "func_num_args":
		return true, vm.funcNumArgs(frame), nil
	case "function_exists":
		return true, vm.funcFunctionExists(params), nil
	default:
		return false, nil, nil
	}
}

// funcDefine implements define(string $name, mixed $value): bool
func (vm *VM) funcDefine(params []*types.Value) *types.Value {
	if len(params) < 2 {
		return types.NewBool(false)
	}
	name := params[0].ToString()
	value := params[1]
	success := vm.DefineUserConstant(name, value)
	return types.NewBool(success)
}

// funcDefined implements defined(string $name): bool
func (vm *VM) funcDefined(params []*types.Value) *types.Value {
	if len(params) < 1 {
		return types.NewBool(false)
	}
	name := params[0].ToString()
	exists := vm.UserConstantExists(name)
	return types.NewBool(exists)
}

// funcConstant implements constant(string $name): mixed
func (vm *VM) funcConstant(params []*types.Value) *types.Value {
	if len(params) < 1 {
		return types.NewNull()
	}
	name := params[0].ToString()
	if val, ok := vm.GetUserConstant(name); ok {
		return val
	}
	// Constant not found - in PHP this throws a warning
	return types.NewNull()
}

// funcCallUserFunc implements call_user_func(callable $callback, mixed ...$args): mixed
func (vm *VM) funcCallUserFunc(params []*types.Value, frame *Frame) (bool, *types.Value, error) {
	if len(params) < 1 {
		return true, types.NewNull(), nil
	}

	callback := params[0]
	args := params[1:]

	// Handle string callback (function name)
	if callback.Type() == types.TypeString {
		funcName := callback.ToString()

		// Check if it's a builtin
		if IsBuiltin(funcName) {
			result, err := CallBuiltin(funcName, args)
			if err != nil {
				return true, types.NewNull(), err
			}
			return true, result, nil
		}

		// Check if it's a user function
		if fn, ok := vm.GetFunction(funcName); ok {
			// Create a new frame and execute
			newFrame := NewFrame(fn)
			for i, arg := range args {
				if i < fn.NumParams {
					newFrame.setParam(i, arg)
				}
			}
			if err := vm.pushFrame(newFrame); err != nil {
				return true, types.NewNull(), err
			}
			if err := vm.runFrame(newFrame); err != nil {
				return true, types.NewNull(), err
			}
			completedFrame := vm.popFrame()
			return true, completedFrame.getReturnValue(), nil
		}

		return true, types.NewNull(), fmt.Errorf("call_user_func: undefined function %s", funcName)
	}

	// Handle closure/resource callback
	if callback.Type() == types.TypeResource {
		res := callback.ToResource()
		if res != nil {
			if closure, ok := res.Data().(*Closure); ok {
				newFrame := NewFrame(closure.Function)
				for i, arg := range args {
					if i < closure.Function.NumParams {
						newFrame.setParam(i, arg)
					}
				}
				// Set up captured variables
				for varName, varValue := range closure.CapturedVars {
					vm.globals[varName] = varValue
				}
				if err := vm.pushFrame(newFrame); err != nil {
					return true, types.NewNull(), err
				}
				if err := vm.runFrame(newFrame); err != nil {
					return true, types.NewNull(), err
				}
				completedFrame := vm.popFrame()
				return true, completedFrame.getReturnValue(), nil
			}
		}
	}

	return true, types.NewNull(), fmt.Errorf("call_user_func: invalid callback")
}

// funcCallUserFuncArray implements call_user_func_array(callable $callback, array $args): mixed
func (vm *VM) funcCallUserFuncArray(params []*types.Value, frame *Frame) (bool, *types.Value, error) {
	if len(params) < 2 {
		return true, types.NewNull(), nil
	}

	callback := params[0]
	argsArray := params[1]

	if argsArray.Type() != types.TypeArray {
		return true, types.NewNull(), nil
	}

	// Convert array to slice of values
	arr := argsArray.ToArray()
	args := make([]*types.Value, 0)
	for _, key := range arr.GetKeysSlice() {
		var keyVal *types.Value
		switch k := key.(type) {
		case int64:
			keyVal = types.NewInt(k)
		case string:
			keyVal = types.NewString(k)
		default:
			continue
		}
		val, _ := arr.Get(keyVal)
		if val != nil {
			args = append(args, val)
		}
	}

	// Reuse call_user_func logic
	newParams := make([]*types.Value, 0, len(args)+1)
	newParams = append(newParams, callback)
	newParams = append(newParams, args...)
	return vm.funcCallUserFunc(newParams, frame)
}

// funcGetArgs implements func_get_args(): array
func (vm *VM) funcGetArgs(frame *Frame) *types.Value {
	arr := types.NewEmptyArray()

	// Get all parameters passed to the current function
	if frame != nil && frame.fn != nil {
		for i := 0; i < frame.fn.NumParams; i++ {
			val := frame.getLocal(i)
			if val != nil {
				arr.Append(val)
			}
		}
	}

	return types.NewArray(arr)
}

// funcNumArgs implements func_num_args(): int
func (vm *VM) funcNumArgs(frame *Frame) *types.Value {
	if frame != nil && frame.fn != nil {
		return types.NewInt(int64(frame.fn.NumParams))
	}
	return types.NewInt(0)
}

// funcFunctionExists implements function_exists(string $function_name): bool
// This version checks both builtins and user-defined functions
func (vm *VM) funcFunctionExists(params []*types.Value) *types.Value {
	if len(params) < 1 {
		return types.NewBool(false)
	}
	funcName := params[0].ToString()

	// Check builtins
	if IsBuiltin(funcName) {
		return types.NewBool(true)
	}

	// Check user-defined functions
	if _, ok := vm.GetFunction(funcName); ok {
		return types.NewBool(true)
	}

	return types.NewBool(false)
}
