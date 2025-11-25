package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Closure Methods
// ============================================================================

// NewClosure creates a new closure
func NewClosure(fn *CompiledFunction) *Closure {
	return &Closure{
		Function:     fn,
		CapturedVars: make(map[string]*types.Value),
		Static:       false,
		ReturnByRef:  false,
	}
}

// NewStaticClosure creates a new static closure (no $this access)
func NewStaticClosure(fn *CompiledFunction) *Closure {
	return &Closure{
		Function:     fn,
		CapturedVars: make(map[string]*types.Value),
		Static:       true,
		ReturnByRef:  false,
	}
}

// BindVariable binds a variable to the closure
// If byRef is true, the variable is captured by reference
func (c *Closure) BindVariable(name string, value *types.Value, byRef bool) {
	if byRef {
		// Capture by reference - store the reference itself
		if value.Type() == types.TypeReference {
			c.CapturedVars[name] = value
		} else {
			// Create a reference to the value
			c.CapturedVars[name] = types.NewReference(value)
		}
	} else {
		// Capture by value - create a copy
		// For now, we store the value directly
		// PHP actually does a copy-on-write optimization
		c.CapturedVars[name] = value
	}
}

// GetCapturedVariable retrieves a captured variable
func (c *Closure) GetCapturedVariable(name string) (*types.Value, bool) {
	val, ok := c.CapturedVars[name]
	if !ok {
		return types.NewNull(), false
	}

	// If it's a reference, dereference it
	if val.Type() == types.TypeReference {
		return val.Deref(), true
	}

	return val, true
}

// HasCapturedVariable checks if a variable is captured
func (c *Closure) HasCapturedVariable(name string) bool {
	_, ok := c.CapturedVars[name]
	return ok
}

// Call invokes the closure
// This is a higher-level method that the VM can use
func (c *Closure) Call(vm *VM, args []*types.Value, thisObj *types.Object) (*types.Value, error) {
	// Create a new frame for the closure
	frame := NewFrame(c.Function)

	// Set $this if not a static closure
	if !c.Static && thisObj != nil {
		frame.thisObject = thisObj
	}

	// Set parameters
	for i, arg := range args {
		if i < c.Function.NumParams {
			frame.setParam(i, arg)
		}
	}

	// Copy captured variables into frame locals
	// The compiler should have allocated locals for captured variables
	// For now, we'll assume they're accessible as variables
	// TODO: Integrate with frame variable management

	// Push frame and execute
	if err := vm.pushFrame(frame); err != nil {
		return nil, err
	}
	err := vm.runFrame(frame)
	if err != nil {
		return nil, err
	}

	// Get return value
	result := frame.getReturnValue()
	vm.popFrame()

	return result, nil
}

// Bind creates a new closure with a different $this and/or scope
// This implements Closure::bindTo()
func (c *Closure) Bind(newThis *types.Object, newScope *types.ClassEntry) (*Closure, error) {
	// Static closures cannot be rebound
	if c.Static {
		return nil, fmt.Errorf("Cannot bind an instance to a static closure")
	}

	// Create a new closure with the same function and captured variables
	newClosure := &Closure{
		Function:     c.Function,
		CapturedVars: make(map[string]*types.Value),
		Static:       false,
		ReturnByRef:  c.ReturnByRef,
	}

	// Copy captured variables (shallow copy)
	for k, v := range c.CapturedVars {
		newClosure.CapturedVars[k] = v
	}

	// The new $this and scope will be handled by the caller
	// This is just the data structure binding

	return newClosure, nil
}

// BindStatic creates a static version of the closure
func (c *Closure) BindStatic(newScope *types.ClassEntry) (*Closure, error) {
	// Create a new static closure
	newClosure := &Closure{
		Function:     c.Function,
		CapturedVars: make(map[string]*types.Value),
		Static:       true,
		ReturnByRef:  c.ReturnByRef,
	}

	// Copy captured variables
	for k, v := range c.CapturedVars {
		newClosure.CapturedVars[k] = v
	}

	return newClosure, nil
}

// IsStatic returns true if this is a static closure
func (c *Closure) IsStatic() bool {
	return c.Static
}

// GetFunction returns the compiled function
func (c *Closure) GetFunction() *CompiledFunction {
	return c.Function
}

// GetCapturedVariables returns all captured variables
func (c *Closure) GetCapturedVariables() map[string]*types.Value {
	return c.CapturedVars
}

// Clone creates a deep copy of the closure
func (c *Closure) Clone() *Closure {
	newClosure := &Closure{
		Function:     c.Function, // Function is immutable, so we can share it
		CapturedVars: make(map[string]*types.Value),
		Static:       c.Static,
		ReturnByRef:  c.ReturnByRef,
	}

	// Deep copy captured variables
	for k, v := range c.CapturedVars {
		// For references, we want to keep the reference
		// For values, we create a new value
		if v.Type() == types.TypeReference {
			newClosure.CapturedVars[k] = v // Keep the same reference
		} else {
			// Clone the value
			// TODO: Implement proper value cloning for complex types
			newClosure.CapturedVars[k] = v
		}
	}

	return newClosure
}
