package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// VM represents the PHP Virtual Machine
type VM struct {
	// Constants pool from compilation
	constants []interface{}

	// User-defined constants (from define())
	userConstants map[string]*types.Value

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

	// Exit flag (set by exit() or die())
	exited bool
	exitCode int
}

// CompiledFunction represents a compiled PHP function
type CompiledFunction struct {
	Name         string
	Instructions Instructions
	NumLocals    int // Number of local variables (CVs + temps)
	NumParams    int // Number of parameters
	NumCVs       int // Number of compiled variables (user variables)
}

// Closure represents a PHP closure/anonymous function with captured variables
type Closure struct {
	Function        *CompiledFunction
	CapturedVars    map[string]*types.Value // Variables from use clause
	Static          bool                     // static closure (no $this access)
	ReturnByRef     bool                     // Returns by reference
}

// CompiledClass is deprecated - use types.ClassEntry instead
// Kept for backwards compatibility
type CompiledClass = types.ClassEntry

// New creates a new virtual machine
func New() *VM {
	vm := &VM{
		constants:     make([]interface{}, 0),
		userConstants: make(map[string]*types.Value),
		globals:       make(map[string]*types.Value),
		functions:     make(map[string]*CompiledFunction),
		classes:       make(map[string]*CompiledClass),
		frames:        make([]*Frame, 1024), // Pre-allocate frame stack
		frameIndex:    -1,                   // -1 means no frames on stack
		output:        make([]byte, 0),
		maxStackDepth: 1000,
	}
	// Register built-in exception classes
	vm.registerBuiltinClasses()
	// Initialize built-in constants
	vm.initBuiltinConstants()
	return vm
}

// registerBuiltinClasses registers PHP's built-in exception classes
func (vm *VM) registerBuiltinClasses() {
	// Exception - base exception class
	exceptionClass := &types.ClassEntry{
		Name:       "Exception",
		Properties: make(map[string]*types.PropertyDef),
		Methods:    make(map[string]*types.MethodDef),
	}
	// Add exception properties
	exceptionClass.Properties["message"] = &types.PropertyDef{
		Name:       "message",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewString(""),
	}
	exceptionClass.Properties["code"] = &types.PropertyDef{
		Name:       "code",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewInt(0),
	}
	exceptionClass.Properties["file"] = &types.PropertyDef{
		Name:       "file",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewString(""),
	}
	exceptionClass.Properties["line"] = &types.PropertyDef{
		Name:       "line",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewInt(0),
	}
	exceptionClass.Properties["previous"] = &types.PropertyDef{
		Name:       "previous",
		Visibility: types.VisibilityPrivate,
		HasDefault: true,
		Default:    types.NewNull(),
	}
	vm.classes["Exception"] = exceptionClass

	// Error - base error class (PHP 7+)
	errorClass := &types.ClassEntry{
		Name:       "Error",
		Properties: make(map[string]*types.PropertyDef),
		Methods:    make(map[string]*types.MethodDef),
	}
	errorClass.Properties["message"] = &types.PropertyDef{
		Name:       "message",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewString(""),
	}
	errorClass.Properties["code"] = &types.PropertyDef{
		Name:       "code",
		Visibility: types.VisibilityProtected,
		HasDefault: true,
		Default:    types.NewInt(0),
	}
	vm.classes["Error"] = errorClass

	// RuntimeException extends Exception
	runtimeExceptionClass := &types.ClassEntry{
		Name:        "RuntimeException",
		ParentClass: exceptionClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["RuntimeException"] = runtimeExceptionClass

	// LogicException extends Exception
	logicExceptionClass := &types.ClassEntry{
		Name:        "LogicException",
		ParentClass: exceptionClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["LogicException"] = logicExceptionClass

	// InvalidArgumentException extends LogicException
	invalidArgExceptionClass := &types.ClassEntry{
		Name:        "InvalidArgumentException",
		ParentClass: logicExceptionClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["InvalidArgumentException"] = invalidArgExceptionClass

	// OutOfBoundsException extends LogicException
	outOfBoundsExceptionClass := &types.ClassEntry{
		Name:        "OutOfBoundsException",
		ParentClass: logicExceptionClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["OutOfBoundsException"] = outOfBoundsExceptionClass

	// TypeError extends Error
	typeErrorClass := &types.ClassEntry{
		Name:        "TypeError",
		ParentClass: errorClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["TypeError"] = typeErrorClass

	// ArgumentCountError extends TypeError
	argumentCountErrorClass := &types.ClassEntry{
		Name:        "ArgumentCountError",
		ParentClass: typeErrorClass,
		Properties:  make(map[string]*types.PropertyDef),
		Methods:     make(map[string]*types.MethodDef),
	}
	vm.classes["ArgumentCountError"] = argumentCountErrorClass
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
	// Safety: This is the first frame (frameIndex == -1), so it will never exceed maxStackDepth
	// If this somehow fails, it's a programming error and we should panic
	if err := vm.pushFrame(NewFrame(mainFunc)); err != nil {
		panic(fmt.Sprintf("failed to push initial frame: %v", err))
	}

	return vm
}

// LoadConstants loads constants from compiled bytecode
func (vm *VM) LoadConstants(constants []interface{}) {
	vm.constants = constants
}

// Execute executes the bytecode starting from the main program
func (vm *VM) Execute(instructions Instructions) error {
	return vm.ExecuteWithCVs(instructions, 0)
}

// ExecuteWithCVs executes bytecode with a specified number of compiled variables
func (vm *VM) ExecuteWithCVs(instructions Instructions, numCVs int) error {
	// Create main function
	mainFunc := &CompiledFunction{
		Name:         "main",
		Instructions: instructions,
		NumLocals:    100,
		NumParams:    0,
		NumCVs:       numCVs,
	}

	// Push main frame
	frame := NewFrame(mainFunc)
	if err := vm.pushFrame(frame); err != nil {
		return err
	}

	// Run the execution loop
	return vm.run()
}

// run executes the main VM loop
func (vm *VM) run() error {
	for vm.frameIndex >= 0 {
		// Check if exit() or die() was called
		if vm.exited {
			return nil
		}

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

		// Track line number for stack traces
		frame.lastLine = instr.Lineno

		// Dispatch instruction
		if err := vm.dispatch(frame, instr); err != nil {
			return err
		}
	}

	return nil
}

// runFrame executes a single frame until completion
func (vm *VM) runFrame(frame *Frame) error {
	for frame.ip < len(frame.fn.Instructions) {
		// Fetch next instruction
		instr := frame.fn.Instructions[frame.ip]
		frame.ip++

		// Track line number for stack traces
		frame.lastLine = instr.Lineno

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
	case OpQMAssign:
		return vm.opQMAssign(frame, instr)
	case OpFetchR:
		return vm.opFetch(frame, instr)
	case OpFree:
		return vm.opFree(frame, instr)
	case OpIssetIsemptyVar:
		return vm.opIssetIsemptyVar(frame, instr)
	case OpBindGlobal:
		return vm.opBindGlobal(frame, instr)
	case OpUnsetVar:
		return vm.opUnsetVar(frame, instr)
	case OpIncludeOrEval:
		return vm.opIncludeOrEval(frame, instr)

	// Control flow
	case OpJmp:
		return vm.opJmp(frame, instr)
	case OpJmpZ:
		return vm.opJmpZ(frame, instr)
	case OpJmpNZ:
		return vm.opJmpNZ(frame, instr)

	// Foreach
	case OpFeResetR:
		return vm.opFeResetR(frame, instr)
	case OpFeResetRW:
		return vm.opFeResetRW(frame, instr)
	case OpFeFetchR:
		return vm.opFeFetchR(frame, instr)
	case OpFeFetchRW:
		return vm.opFeFetchRW(frame, instr)
	case OpFeFree:
		return vm.opFeFree(frame, instr)

	// Functions
	case OpReturn:
		return vm.opReturn(frame, instr)
	case OpRecv:
		return vm.opRecv(frame, instr)
	case OpRecvInit:
		return vm.opRecvInit(frame, instr)
	case OpInitFcall:
		return vm.opInitFcall(frame, instr)
	case OpInitFcallByName:
		return vm.opInitFcallByName(frame, instr)
	case OpSendVal:
		return vm.opSendVal(frame, instr)
	case OpSendRef:
		return vm.opSendRef(frame, instr)
	case OpDoFcall:
		return vm.opDoFcall(frame, instr)
	case OpDoUcall:
		return vm.opDoUcall(frame, instr)
	case OpDoIcall:
		return vm.opDoIcall(frame, instr)

	// I/O
	case OpEcho:
		return vm.opEcho(frame, instr)

	// String operations
	case OpConcat:
		return vm.opConcat(frame, instr)
	case OpFastConcat:
		return vm.opFastConcat(frame, instr)

	// Array operations
	case OpInitArray:
		return vm.opInitArray(frame, instr)
	case OpAddArrayElement:
		return vm.opAddArrayElement(frame, instr)
	case OpFetchDimR:
		return vm.opFetchDimR(frame, instr)
	case OpFetchDimW:
		return vm.opFetchDimW(frame, instr)
	case OpFetchDimRW:
		return vm.opFetchDimRW(frame, instr)
	case OpFetchDimIs:
		return vm.opFetchDimIs(frame, instr)
	case OpFetchDimFuncArg:
		return vm.opFetchDimFuncArg(frame, instr)
	case OpFetchDimUnset:
		return vm.opFetchDimUnset(frame, instr)
	case OpAssignDim:
		return vm.opAssignDim(frame, instr)
	case OpAssignDimOp:
		return vm.opAssignDimOp(frame, instr)
	case OpUnsetDim:
		return vm.opUnsetDim(frame, instr)
	case OpIssetIsemptyDimObj:
		return vm.opIssetIsemptyDimObj(frame, instr)
	case OpCount:
		return vm.opCount(frame, instr)
	case OpInArray:
		return vm.opInArray(frame, instr)
	case OpArrayKeyExists:
		return vm.opArrayKeyExists(frame, instr)

	// Closure operations
	case OpDeclareFunction:
		return vm.opDeclareFunction(frame, instr)
	case OpDeclareLambdaFunction:
		return vm.opDeclareLambdaFunction(frame, instr)
	case OpBindLexical:
		return vm.opBindLexical(frame, instr)

	// Object property operations - Fetch
	case OpFetchObjR:
		return vm.opFetchObjR(frame, instr)
	case OpFetchObjW:
		return vm.opFetchObjW(frame, instr)
	case OpFetchObjRW:
		return vm.opFetchObjRW(frame, instr)
	case OpFetchObjIs:
		return vm.opFetchObjIs(frame, instr)
	case OpFetchObjFuncArg:
		return vm.opFetchObjFuncArg(frame, instr)
	case OpFetchObjUnset:
		return vm.opFetchObjUnset(frame, instr)

	// Object property operations - Assignment
	case OpAssignObj:
		return vm.opAssignObj(frame, instr)
	case OpAssignObjOp:
		return vm.opAssignObjOp(frame, instr)
	case OpAssignObjRef:
		return vm.opAssignObjRef(frame, instr)

	// Object property operations - Unset/Isset
	case OpUnsetObj:
		return vm.opUnsetObj(frame, instr)
	case OpIssetIsemptyPropObj:
		return vm.opIssetIsemptyPropObj(frame, instr)

	// Object property operations - Increment/Decrement
	case OpPreIncObj:
		return vm.opPreIncObj(frame, instr)
	case OpPreDecObj:
		return vm.opPreDecObj(frame, instr)
	case OpPostIncObj:
		return vm.opPostIncObj(frame, instr)
	case OpPostDecObj:
		return vm.opPostDecObj(frame, instr)

	// Object creation and method calls
	case OpNew:
		return vm.opNew(frame, instr)
	case OpInitMethodCall:
		return vm.opInitMethodCall(frame, instr)
	case OpInitStaticMethodCall:
		return vm.opInitStaticMethodCall(frame, instr)
	case OpClone:
		return vm.opClone(frame, instr)
	case OpInstanceof:
		return vm.opInstanceof(frame, instr)
	case OpGetClass:
		return vm.opGetClass(frame, instr)
	case OpFetchClassName:
		return vm.opFetchClassName(frame, instr)
	case OpFetchThis:
		return vm.opFetchThis(frame, instr)
	case OpDeclareClass:
		return vm.opDeclareClass(frame, instr)
	case OpFetchClassConstant:
		return vm.opFetchClassConstant(frame, instr)

	// Generator operations
	case OpGeneratorCreate:
		return vm.opGeneratorCreate(frame, instr)
	case OpYield:
		return vm.opYield(frame, instr)
	case OpGeneratorReturn:
		return vm.opGeneratorReturn(frame, instr)
	case OpYieldFrom:
		return vm.opYieldFrom(frame, instr)

	// Exception operations
	case OpThrow:
		return vm.opThrow(frame, instr)
	case OpCatch:
		return vm.opCatch(frame, instr)
	case OpFastCall:
		return vm.opFastCall(frame, instr)
	case OpFastRet:
		return vm.opFastRet(frame, instr)

	// Exit operations
	case OpExit:
		return vm.opExit(frame, instr)

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
	case *types.ClassEntry:
		// ClassEntry constants are stored as-is for use by DECLARE_CLASS
		// They are not converted to Values - just kept as metadata
		return nil, nil // Return nil Value, will be accessed directly from constants
	default:
		return nil, fmt.Errorf("unsupported constant type: %T", c)
	}
}

// ============================================================================
// User-Defined Constants
// ============================================================================

// DefineUserConstant defines a user constant (from define())
func (vm *VM) DefineUserConstant(name string, value *types.Value) bool {
	// Check if already defined (case-sensitive)
	if _, exists := vm.userConstants[name]; exists {
		return false // Already defined
	}
	vm.userConstants[name] = value
	return true
}

// GetUserConstant retrieves a user-defined constant
func (vm *VM) GetUserConstant(name string) (*types.Value, bool) {
	val, ok := vm.userConstants[name]
	return val, ok
}

// UserConstantExists checks if a user constant exists
func (vm *VM) UserConstantExists(name string) bool {
	_, exists := vm.userConstants[name]
	return exists
}

// initBuiltinConstants initializes PHP built-in constants
func (vm *VM) initBuiltinConstants() {
	// Boolean constants
	vm.userConstants["true"] = types.NewBool(true)
	vm.userConstants["TRUE"] = types.NewBool(true)
	vm.userConstants["false"] = types.NewBool(false)
	vm.userConstants["FALSE"] = types.NewBool(false)
	vm.userConstants["null"] = types.NewNull()
	vm.userConstants["NULL"] = types.NewNull()

	// PHP version constants
	vm.userConstants["PHP_VERSION"] = types.NewString("8.4.0")
	vm.userConstants["PHP_MAJOR_VERSION"] = types.NewInt(8)
	vm.userConstants["PHP_MINOR_VERSION"] = types.NewInt(4)
	vm.userConstants["PHP_RELEASE_VERSION"] = types.NewInt(0)

	// System constants
	vm.userConstants["PHP_EOL"] = types.NewString("\n")
	vm.userConstants["PHP_INT_MAX"] = types.NewInt(9223372036854775807)
	vm.userConstants["PHP_INT_MIN"] = types.NewInt(-9223372036854775808)

	// Common constants
	vm.userConstants["DIRECTORY_SEPARATOR"] = types.NewString("/")
	vm.userConstants["PATH_SEPARATOR"] = types.NewString(":")

	// Sort constants for array functions
	vm.userConstants["SORT_REGULAR"] = types.NewInt(0)
	vm.userConstants["SORT_NUMERIC"] = types.NewInt(1)
	vm.userConstants["SORT_STRING"] = types.NewInt(2)
	vm.userConstants["SORT_LOCALE_STRING"] = types.NewInt(5)
	vm.userConstants["SORT_NATURAL"] = types.NewInt(6)
	vm.userConstants["SORT_FLAG_CASE"] = types.NewInt(8)

	// Case constants
	vm.userConstants["CASE_LOWER"] = types.NewInt(0)
	vm.userConstants["CASE_UPPER"] = types.NewInt(1)

	// JSON constants
	vm.userConstants["JSON_ERROR_NONE"] = types.NewInt(0)
	vm.userConstants["JSON_ERROR_DEPTH"] = types.NewInt(1)
	vm.userConstants["JSON_ERROR_STATE_MISMATCH"] = types.NewInt(2)
	vm.userConstants["JSON_ERROR_CTRL_CHAR"] = types.NewInt(3)
	vm.userConstants["JSON_ERROR_SYNTAX"] = types.NewInt(4)
	vm.userConstants["JSON_ERROR_UTF8"] = types.NewInt(5)
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
// Exit Status
// ============================================================================

// GetExitCode returns the exit code set by exit() or die()
func (vm *VM) GetExitCode() int {
	return vm.exitCode
}

// HasExited returns true if exit() or die() was called
func (vm *VM) HasExited() bool {
	return vm.exited
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
		// Compiled variable (parameters and user variables)
		// CVs start at index 0 (or after params in functions)
		idx := int(op.Value)
		// Check if this variable is bound to a global
		if frame.globalBindings != nil {
			if globalName, bound := frame.globalBindings[idx]; bound {
				// Return value from globals instead of locals
				if val, exists := vm.globals[globalName]; exists {
					return val, nil
				}
				return types.NewNull(), nil
			}
		}
		return frame.getLocal(idx), nil
	case OpTmpVar:
		// Temporary variable (offset to avoid conflicts with CVs)
		// TMPVARs start after all compiled variables
		// If NumCVs is explicitly set and > 0, use it as offset
		// Otherwise, use NumParams for backward compatibility
		offset := frame.fn.NumParams
		if frame.fn.NumCVs > 0 {
			offset = frame.fn.NumCVs
		}
		return frame.getLocal(int(op.Value) + offset), nil
	case OpUnused:
		return types.NewNull(), nil
	default:
		return nil, fmt.Errorf("unknown operand type: %v", op.Type)
	}
}

// setOperandValue sets the value of an operand
func (vm *VM) setOperandValue(frame *Frame, op Operand, value *types.Value) error {
	switch op.Type {
	case OpVar, OpCV:
		// Compiled variable (parameters and user variables)
		idx := int(op.Value)
		// Check if this variable is bound to a global
		if frame.globalBindings != nil {
			if globalName, bound := frame.globalBindings[idx]; bound {
				// Write to globals instead of locals
				vm.globals[globalName] = value
				return nil
			}
		}
		frame.setLocal(idx, value)
		return nil
	case OpTmpVar:
		// Temporary variable (offset to avoid conflicts with CVs)
		// Must match getOperandValue offset calculation
		offset := frame.fn.NumParams
		if frame.fn.NumCVs > 0 {
			offset = frame.fn.NumCVs
		}
		frame.setLocal(int(op.Value)+offset, value)
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

