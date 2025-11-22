package vm

import "github.com/krizos/php-go/pkg/types"

// CallParams holds parameters being collected for a function call
type CallParams struct {
	params []*types.Value
}

// Frame represents a single execution frame (function call)
type Frame struct {
	// The function being executed
	fn *CompiledFunction

	// Instruction pointer (current instruction index)
	ip int

	// Local variables (including parameters and temporaries)
	locals []*types.Value

	// Return value (set by return statement)
	returnValue *types.Value

	// Base pointer (for stack-based operations)
	bp int

	// Object/class context (for method calls)
	thisObject    *types.Object     // $this for instance methods
	currentClass  *types.ClassEntry // Current class context for self/parent
	calledClass   *types.ClassEntry // Called class for late static binding (static::)

	// Pending method call information (set by OpInitMethodCall)
	pendingMethod *types.MethodDef // Method to be called
	pendingObject *types.Object    // Object for instance method calls (nil for static)

	// Pending function call information (set by OpInitFcall)
	pendingFunction *CompiledFunction // Function to be called
	pendingParams   *CallParams       // Parameters being collected
}

// NewFrame creates a new execution frame for a function
func NewFrame(fn *CompiledFunction) *Frame {
	// Allocate local variable storage
	numLocals := fn.NumLocals
	if numLocals < 10 {
		numLocals = 10 // Minimum allocation
	}

	return &Frame{
		fn:          fn,
		ip:          0,
		locals:      make([]*types.Value, numLocals),
		returnValue: types.NewNull(),
		bp:          0,
	}
}

// ============================================================================
// Local Variable Access
// ============================================================================

// getLocal retrieves a local variable by index
func (f *Frame) getLocal(index int) *types.Value {
	if index < 0 || index >= len(f.locals) {
		return types.NewNull()
	}

	val := f.locals[index]
	if val == nil {
		return types.NewNull()
	}

	// Debug: log local variable access (disabled)
	// fmt.Printf("DEBUG getLocal [%s]: index=%d, value=%v, type=%v\n", f.fn.Name, index, val, val.Type())

	return val
}

// setLocal sets a local variable by index
func (f *Frame) setLocal(index int, value *types.Value) {
	// Expand locals if needed
	if index >= len(f.locals) {
		newSize := index + 1
		if newSize < len(f.locals)*2 {
			newSize = len(f.locals) * 2
		}

		newLocals := make([]*types.Value, newSize)
		copy(newLocals, f.locals)
		f.locals = newLocals
	}

	f.locals[index] = value
}

// ============================================================================
// Stack Operations (for temporaries)
// ============================================================================

// In PHP VM, temporaries are just local variables
// The stack operations are helper methods for managing temporaries

// pushTemp pushes a temporary value onto the stack
// Returns the index where the value was stored
func (f *Frame) pushTemp(value *types.Value) int {
	// Find next available temp slot
	for i, v := range f.locals {
		if v == nil || v.IsUndef() {
			f.locals[i] = value
			return i
		}
	}

	// Need to allocate more space
	index := len(f.locals)
	f.setLocal(index, value)
	return index
}

// popTemp pops a temporary value from the stack
func (f *Frame) popTemp(index int) *types.Value {
	val := f.getLocal(index)
	f.locals[index] = nil // Clear the slot
	return val
}

// ============================================================================
// Parameter Handling
// ============================================================================

// setParam sets a parameter value
func (f *Frame) setParam(index int, value *types.Value) {
	if index < 0 || index >= f.fn.NumParams {
		return
	}

	f.setLocal(index, value)
}

// getParam gets a parameter value
func (f *Frame) getParam(index int) *types.Value {
	if index < 0 || index >= f.fn.NumParams {
		return types.NewNull()
	}

	return f.getLocal(index)
}

// ============================================================================
// Return Value
// ============================================================================

// setReturnValue sets the return value for this frame
func (f *Frame) setReturnValue(value *types.Value) {
	f.returnValue = value
}

// getReturnValue gets the return value
func (f *Frame) getReturnValue() *types.Value {
	if f.returnValue == nil {
		return types.NewNull()
	}
	return f.returnValue
}

// ============================================================================
// Debugging
// ============================================================================

// String returns a string representation of the frame for debugging
func (f *Frame) String() string {
	if f == nil {
		return "<nil frame>"
	}

	return f.fn.Name + " @ " + string(rune(f.ip))
}
