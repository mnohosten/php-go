package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// VM represents the PHP Virtual Machine
type VM struct {
	// Constants pool from compilation
	constants []interface{}

	// Global variables ($_GET, $_POST, user globals, etc.)
	globals map[string]*types.Value

	// Function registry (user functions and built-ins)
	functions map[string]*CompiledFunction

	// Class registry
	classes map[string]*CompiledClass

	// Call stack (frames)
	frames []*Frame
	// Current frame index
	frameIndex int

	// Output buffer
	output []byte

	// Maximum stack depth (default 1000)
	maxStackDepth int
}

// CompiledFunction represents a compiled PHP function
type CompiledFunction struct {
	Name         string
	Instructions Instructions
	NumLocals    int // Number of local variables
	NumParams    int // Number of parameters
}

// Closure represents a PHP closure/anonymous function with captured variables
type Closure struct {
	Function        *CompiledFunction
	CapturedVars    map[string]*types.Value // Variables from use clause
	Static          bool                     // static closure (no $this access)
	ReturnByRef     bool                     // Returns by reference
}

// CompiledClass represents a compiled PHP class
type CompiledClass struct {
	Name       string
	ParentName string
	Properties map[string]*types.Value
	Methods    map[string]*CompiledFunction
}

// New creates a new virtual machine
func New() *VM {
	return &VM{
		constants:     make([]interface{}, 0),
		globals:       make(map[string]*types.Value),
		functions:     make(map[string]*CompiledFunction),
		classes:       make(map[string]*CompiledClass),
		frames:        make([]*Frame, 1024), // Pre-allocate frame stack
		frameIndex:    -1,                   // -1 means no frames on stack
		output:        make([]byte, 0),
		maxStackDepth: 1000,
	}
}

// NewWithBytecode creates a new VM and loads the bytecode
func NewWithBytecode(instructions Instructions, constants []interface{}) *VM {
	vm := New()
	vm.constants = constants

	// Create main function frame
	mainFunc := &CompiledFunction{
		Name:         "main",
		Instructions: instructions,
		NumLocals:    100, // Allocate space for locals
		NumParams:    0,
	}

	// Push main frame
	vm.pushFrame(NewFrame(mainFunc))

	return vm
}

// LoadConstants loads constants from compiled bytecode
func (vm *VM) LoadConstants(constants []interface{}) {
	vm.constants = constants
}

// Execute executes the bytecode starting from the main program
func (vm *VM) Execute(instructions Instructions) error {
	// Create main function
	mainFunc := &CompiledFunction{
		Name:         "main",
		Instructions: instructions,
		NumLocals:    100,
		NumParams:    0,
	}

	// Push main frame
	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Run the execution loop
	return vm.run()
}

// run executes the main VM loop
func (vm *VM) run() error {
	for vm.frameIndex >= 0 {
		frame := vm.currentFrame()

		// Check if we've finished this frame
		if frame.ip >= len(frame.fn.Instructions) {
			// Pop frame and return
			vm.popFrame()
			continue
		}

		// Fetch next instruction
		instr := frame.fn.Instructions[frame.ip]
		frame.ip++

		// Dispatch instruction
		if err := vm.dispatch(frame, instr); err != nil {
			return err
		}
	}

	return nil
}

// dispatch executes a single instruction
func (vm *VM) dispatch(frame *Frame, instr Instruction) error {
	switch instr.Opcode {
	// Arithmetic operations
	case OpAdd:
		return vm.opAdd(frame, instr)
	case OpSub:
		return vm.opSub(frame, instr)
	case OpMul:
		return vm.opMul(frame, instr)
	case OpDiv:
		return vm.opDiv(frame, instr)
	case OpMod:
		return vm.opMod(frame, instr)
	case OpPow:
		return vm.opPow(frame, instr)

	// Comparison operations
	case OpIsEqual:
		return vm.opIsEqual(frame, instr)
	case OpIsNotEqual:
		return vm.opIsNotEqual(frame, instr)
	case OpIsIdentical:
		return vm.opIsIdentical(frame, instr)
	case OpIsNotIdentical:
		return vm.opIsNotIdentical(frame, instr)
	case OpIsSmaller:
		return vm.opIsSmaller(frame, instr)
	case OpIsSmallerOrEqual:
		return vm.opIsSmallerOrEqual(frame, instr)
	case OpSpaceship:
		return vm.opSpaceship(frame, instr)

	// Bitwise operations
	case OpBWAnd:
		return vm.opBWAnd(frame, instr)
	case OpBWOr:
		return vm.opBWOr(frame, instr)
	case OpBWXor:
		return vm.opBWXor(frame, instr)
	case OpBWNot:
		return vm.opBWNot(frame, instr)
	case OpSL:
		return vm.opShiftLeft(frame, instr)
	case OpSR:
		return vm.opShiftRight(frame, instr)

	// Logical operations
	case OpBoolNot:
		return vm.opBoolNot(frame, instr)

	// Constants
	case OpFetchConstant:
		return vm.opConst(frame, instr)

	// Variables
	case OpAssign:
		return vm.opAssign(frame, instr)
	case OpFetchR:
		return vm.opFetch(frame, instr)

	// Control flow
	case OpJmp:
		return vm.opJmp(frame, instr)
	case OpJmpZ:
		return vm.opJmpZ(frame, instr)
	case OpJmpNZ:
		return vm.opJmpNZ(frame, instr)

	// Functions
	case OpReturn:
		return vm.opReturn(frame, instr)

	// I/O
	case OpEcho:
		return vm.opEcho(frame, instr)

	// String operations
	case OpConcat:
		return vm.opConcat(frame, instr)

	// Closure operations
	case OpDeclareLambdaFunction:
		return vm.opDeclareLambdaFunction(frame, instr)
	case OpBindLexical:
		return vm.opBindLexical(frame, instr)

	default:
		return fmt.Errorf("unknown opcode: %s", instr.Opcode)
	}
}

// ============================================================================
// Frame Management
// ============================================================================

// currentFrame returns the current execution frame
func (vm *VM) currentFrame() *Frame {
	if vm.frameIndex < 0 {
		return nil
	}
	return vm.frames[vm.frameIndex]
}

// pushFrame pushes a new frame onto the call stack
func (vm *VM) pushFrame(frame *Frame) error {
	if vm.frameIndex+1 >= vm.maxStackDepth {
		return fmt.Errorf("stack overflow: maximum depth %d exceeded", vm.maxStackDepth)
	}

	vm.frameIndex++
	vm.frames[vm.frameIndex] = frame
	return nil
}

// popFrame pops the current frame from the call stack
func (vm *VM) popFrame() *Frame {
	if vm.frameIndex < 0 {
		return nil
	}

	frame := vm.frames[vm.frameIndex]
	vm.frames[vm.frameIndex] = nil // Clear reference
	vm.frameIndex--
	return frame
}

// ============================================================================
// Global Variables
// ============================================================================

// SetGlobal sets a global variable
func (vm *VM) SetGlobal(name string, value *types.Value) {
	vm.globals[name] = value
}

// GetGlobal gets a global variable
func (vm *VM) GetGlobal(name string) (*types.Value, bool) {
	val, ok := vm.globals[name]
	return val, ok
}

// ============================================================================
// Functions
// ============================================================================

// RegisterFunction registers a compiled function
func (vm *VM) RegisterFunction(name string, fn *CompiledFunction) {
	vm.functions[name] = fn
}

// GetFunction gets a compiled function
func (vm *VM) GetFunction(name string) (*CompiledFunction, bool) {
	fn, ok := vm.functions[name]
	return fn, ok
}

// ============================================================================
// Constants
// ============================================================================

// GetConstant retrieves a constant from the constant pool
func (vm *VM) GetConstant(index int) (*types.Value, error) {
	if index < 0 || index >= len(vm.constants) {
		return nil, fmt.Errorf("constant index out of range: %d", index)
	}

	c := vm.constants[index]

	// Convert to Value
	switch v := c.(type) {
	case int64:
		return types.NewInt(v), nil
	case float64:
		return types.NewFloat(v), nil
	case string:
		return types.NewString(v), nil
	case bool:
		return types.NewBool(v), nil
	case nil:
		return types.NewNull(), nil
	default:
		return nil, fmt.Errorf("unsupported constant type: %T", c)
	}
}

// ============================================================================
// Output
// ============================================================================

// GetOutput returns the captured output
func (vm *VM) GetOutput() string {
	return string(vm.output)
}

// ClearOutput clears the output buffer
func (vm *VM) ClearOutput() {
	vm.output = vm.output[:0]
}

// writeOutput writes to the output buffer
func (vm *VM) writeOutput(data []byte) {
	vm.output = append(vm.output, data...)
}

// ============================================================================
// Helper Methods
// ============================================================================

// getOperandValue retrieves the value of an operand
func (vm *VM) getOperandValue(frame *Frame, op Operand) (*types.Value, error) {
	switch op.Type {
	case OpConst:
		return vm.GetConstant(int(op.Value))
	case OpVar, OpCV:
		// Local variable
		return frame.getLocal(int(op.Value)), nil
	case OpTmpVar:
		// Temporary variable (on stack)
		return frame.getLocal(int(op.Value)), nil
	case OpUnused:
		return types.NewNull(), nil
	default:
		return nil, fmt.Errorf("unknown operand type: %v", op.Type)
	}
}

// setOperandValue sets the value of an operand
func (vm *VM) setOperandValue(frame *Frame, op Operand, value *types.Value) error {
	switch op.Type {
	case OpVar, OpCV, OpTmpVar:
		frame.setLocal(int(op.Value), value)
		return nil
	case OpUnused:
		// Do nothing
		return nil
	default:
		return fmt.Errorf("cannot assign to operand type: %v", op.Type)
	}
}

// ============================================================================
// Closure Operations
// ============================================================================

// opDeclareLambdaFunction creates a closure object
// ExtendedValue: number of parameters
// Op1: flags (static, byref)
// Op2: closure start position
// Result: closure end position, closure object placed in temp var 0
func (vm *VM) opDeclareLambdaFunction(frame *Frame, instr Instruction) error {
	numParams := int(instr.ExtendedValue)
	flags, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	startPos, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	endPos, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	flagsInt := int64(flags.ToInt())
	isStatic := (flagsInt & 1) != 0
	isByRef := (flagsInt & 2) != 0

	// Extract instructions for the closure body
	start := int(startPos.ToInt())
	end := int(endPos.ToInt())
	closureInstructions := frame.fn.Instructions[start:end]

	// Create compiled function for the closure
	compiledFunc := &CompiledFunction{
		Name:         "<closure>",
		Instructions: closureInstructions,
		NumLocals:    100, // TODO: Calculate actual number of locals
		NumParams:    numParams,
	}

	// Create closure object
	_ = &Closure{
		Function:     compiledFunc,
		CapturedVars: make(map[string]*types.Value),
		Static:       isStatic,
		ReturnByRef:  isByRef,
	}

	// TODO: Store closure properly once we have a proper closure value type
	// For now, we'll store it as a PHP object with the closure embedded
	obj := &types.Object{
		ClassName:  "Closure",
		Properties: map[string]*types.Value{
			"__closure__": {}, // Store actual closure here (hack for now)
		},
	}
	closureValue := types.NewObject(obj)
	frame.setLocal(0, closureValue) // Store in temp var 0

	return nil
}

// opBindLexical binds a captured variable to a closure
// Op1: variable name (constant index)
// Op2: by-reference flag (0 = by value, 1 = by reference)
// Result: closure object (in temp var)
func (vm *VM) opBindLexical(frame *Frame, instr Instruction) error {
	// Get variable name from constants
	varNameConst, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	varName := varNameConst.ToString()

	// Get by-reference flag
	byRefFlag, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	isByRef := byRefFlag.ToInt() != 0

	// Get closure object from Result operand
	closureValue, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	// For now, we'll just store a placeholder
	// TODO: Implement proper closure variable binding
	// This requires accessing the parent frame's variables by name
	// and storing them in the closure's captured variables map

	// Create a placeholder captured variable
	varValue := types.NewNull()

	// Log binding for debugging (could be removed in production)
	_ = varName
	_ = isByRef
	_ = closureValue
	_ = varValue

	return nil
}
