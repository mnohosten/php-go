package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// GeneratorState represents the current state of a generator
type GeneratorState int

const (
	GenStateStart GeneratorState = iota // Generator has not been started
	GenStateRunning                      // Generator is currently running
	GenStateYielded                      // Generator has yielded a value
	GenStateDone                         // Generator has finished execution
	GenStateClosed                       // Generator has been explicitly closed
)

// Generator represents a PHP generator object
// Generators are created by functions that use the yield keyword
type Generator struct {
	// VM instance for execution
	vm *VM

	// Function to execute
	function *CompiledFunction

	// Execution frame - saved when yielding, restored when resuming
	frame *Frame

	// Current state
	state GeneratorState

	// Current key/value pair
	key   *types.Value
	value *types.Value

	// Return value (set when generator finishes)
	returnValue *types.Value

	// Exception to throw into generator
	exception *types.Value

	// Instruction pointer when yielded
	savedIP int

	// Whether this is the first iteration
	firstIteration bool

	// Auto-increment key (for yield without key)
	autoKey int64
}

// NewGenerator creates a new generator
func NewGenerator(v *VM, fn *CompiledFunction) *Generator {
	return &Generator{
		vm:             v,
		function:       fn,
		state:          GenStateStart,
		key:            types.NewInt(0),
		value:          types.NewNull(),
		returnValue:    types.NewNull(),
		firstIteration: true,
		autoKey:        0,
	}
}

// Current returns the current value (Iterator interface)
func (g *Generator) Current() *types.Value {
	if g.state == GenStateYielded {
		return g.value
	}
	return types.NewNull()
}

// Key returns the current key (Iterator interface)
func (g *Generator) Key() *types.Value {
	if g.state == GenStateYielded {
		return g.key
	}
	return types.NewNull()
}

// Valid returns whether the iterator is valid (Iterator interface)
func (g *Generator) Valid() bool {
	return g.state == GenStateYielded
}

// Next advances the generator (Iterator interface)
func (g *Generator) Next() error {
	if g.state == GenStateDone || g.state == GenStateClosed {
		return nil
	}

	// Resume execution
	return g.Resume(types.NewNull())
}

// Rewind restarts the generator (Iterator interface)
// Note: In PHP, calling rewind() on an already-started generator throws an exception
func (g *Generator) Rewind() error {
	if g.state != GenStateStart {
		return fmt.Errorf("Cannot rewind a generator that was already run")
	}

	// Initialize and run to first yield
	return g.Resume(types.NewNull())
}

// Resume resumes generator execution
// The actual execution happens in the VM's execution loop
// This method prepares the generator state for resumption
func (g *Generator) Resume(sendValue *types.Value) error {
	if g.state == GenStateDone || g.state == GenStateClosed {
		return nil
	}

	// If this is the first iteration, we need to create the frame and start execution
	if g.state == GenStateStart {
		// Create a new frame for the generator function
		g.frame = NewFrame(g.function)
		g.firstIteration = false
	}

	// Mark as running
	g.state = GenStateRunning

	// Note: The actual execution is driven by the VM
	// The VM will call g.Yield() or g.Return() when appropriate
	return nil
}

// GetFrame returns the execution frame (for VM integration)
func (g *Generator) GetFrame() *Frame {
	return g.frame
}

// SetSendValue sets the value to be sent into the generator (for VM integration)
// This is used when resuming from a yield
func (g *Generator) SetSendValue(value *types.Value) {
	// The VM will push this value onto the stack before resuming
	if g.frame != nil && value != nil {
		// This will be handled by the VM execution loop
	}
}

// Send sends a value to the generator
func (g *Generator) Send(value *types.Value) (*types.Value, error) {
	if g.state == GenStateDone || g.state == GenStateClosed {
		return types.NewNull(), fmt.Errorf("Cannot send to a closed generator")
	}

	if g.firstIteration && value != nil && value.Type() != types.TypeNull {
		return nil, fmt.Errorf("Cannot send a value to a generator on the first iteration")
	}

	err := g.Resume(value)
	if err != nil {
		return nil, err
	}

	return g.value, nil
}

// Throw throws an exception into the generator
func (g *Generator) Throw(exception *types.Value) (*types.Value, error) {
	if g.state == GenStateDone || g.state == GenStateClosed {
		return nil, fmt.Errorf("Cannot throw into a closed generator")
	}

	g.exception = exception
	err := g.Resume(types.NewNull())
	if err != nil {
		return nil, err
	}

	if g.exception != nil {
		// Exception was not caught, re-throw
		return nil, fmt.Errorf("Uncaught exception in generator")
	}

	return g.value, nil
}

// GetReturn returns the return value of the generator
func (g *Generator) GetReturn() (*types.Value, error) {
	if g.state != GenStateDone {
		return nil, fmt.Errorf("Cannot get return value of a generator that hasn't finished")
	}
	return g.returnValue, nil
}

// Close closes the generator
func (g *Generator) Close() {
	g.state = GenStateClosed
	g.frame = nil
}

// Yield is called when the generator yields a value
// This is called by the OpYield opcode handler
// If key is nil, an auto-incrementing integer key is used
func (g *Generator) Yield(key, value *types.Value, ip int) {
	if key == nil || key.Type() == types.TypeNull {
		// Auto-increment key
		g.key = types.NewInt(g.autoKey)
		g.autoKey++
	} else {
		g.key = key
	}
	g.value = value
	g.savedIP = ip
	g.state = GenStateYielded
}

// Return is called when the generator returns
// This is called by the OpGeneratorReturn opcode handler
func (g *Generator) Return(value *types.Value) {
	g.returnValue = value
	g.state = GenStateDone
}

// State returns the current state
func (g *Generator) State() GeneratorState {
	return g.state
}

// SetState sets the generator state
func (g *Generator) SetState(state GeneratorState) {
	g.state = state
}
