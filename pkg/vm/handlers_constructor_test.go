package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Constructor Tests
// ============================================================================

func TestConstructor_BasicInvocation(t *testing.T) {
	vm := New()

	// Create a class with a constructor
	classEntry := types.NewClassEntry("TestClass")
	classEntry.Constructor = &types.MethodDef{
		Name:       "__construct",
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
		// Simple constructor that sets a property: $this->status = "initialized";
		Instructions: []interface{}{
			// Fetch $this
			Instruction{
				Opcode: OpFetchThis,
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			// Fetch constant "initialized"
			Instruction{
				Opcode: OpFetchConstant,
				Op1:    Operand{Type: OpConst, Value: 0}, // "initialized"
				Result: Operand{Type: OpTmpVar, Value: 2},
			},
			// $this->status = "initialized"
			Instruction{
				Opcode: OpAssignObj,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // $this
				Op2:    Operand{Type: OpConst, Value: 1},  // property name "status"
				Result: Operand{Type: OpTmpVar, Value: 2}, // value "initialized"
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpConst, Value: 2}, // null
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals:    10,
		NumParams:    0,
		IsConstructor: true,
	}
	classEntry.Methods["__construct"] = classEntry.Constructor

	// Add property definition
	classEntry.Properties["status"] = &types.PropertyDef{
		Name:       "status",
		Visibility: types.VisibilityPublic,
		Default:    types.NewString("uninitialized"),
	}

	vm.classes["TestClass"] = classEntry
	vm.constants = []interface{}{"initialized", "status", nil}

	// Simulate: $obj = new TestClass(); which compiles to:
	// 1. OpNew - create object
	// 2. OpInitMethodCall - init __construct
	// 3. OpDoFcall - execute __construct
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			// $obj = new TestClass()
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 3}, // "TestClass"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			// $obj->__construct()
			{
				Opcode:        OpInitMethodCall,
				Op1:           Operand{Type: OpTmpVar, Value: 0}, // object
				Op2:           Operand{Type: OpConst, Value: 4},  // "__construct"
				ExtendedValue: 0,                                  // 0 args
			},
			{
				Opcode: OpDoFcall,
				Result: Operand{Type: OpUnused}, // constructor returns void
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	vm.constants = append(vm.constants, "TestClass", "__construct")

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpNew
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpNew failed: %v", err)
	}

	// Verify object was created
	objVal := frame.getLocal(0)
	if objVal.Type() != types.TypeObject {
		t.Fatalf("Expected object, got %v", objVal.Type())
	}

	obj := objVal.ToObject()

	// Execute OpInitMethodCall
	err = vm.dispatch(frame, mainFunc.Instructions[1])
	if err != nil {
		t.Fatalf("OpInitMethodCall failed: %v", err)
	}

	// Verify pending method was set
	if frame.pendingMethod == nil {
		t.Fatal("Expected pending method to be set")
	}

	if frame.pendingMethod.Name != "__construct" {
		t.Errorf("Expected method name '__construct', got '%s'", frame.pendingMethod.Name)
	}

	// Execute OpDoFcall (this will call the constructor)
	frame.ip++ // Advance IP as dispatch does
	err = vm.dispatch(frame, mainFunc.Instructions[2])
	if err != nil {
		t.Fatalf("OpDoFcall failed: %v", err)
	}

	// Verify constructor was executed and property was set
	prop, exists := obj.Properties["status"]
	if !exists {
		t.Fatal("Expected 'status' property to exist")
	}

	if prop.Value.ToString() != "initialized" {
		t.Errorf("Expected property value 'initialized', got '%s'", prop.Value.ToString())
	}
}

func TestConstructor_WithParameters(t *testing.T) {
	vm := New()

	// Create a class with a constructor that takes parameters
	classEntry := types.NewClassEntry("Person")
	classEntry.Constructor = &types.MethodDef{
		Name:       "__construct",
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
		// Constructor: function __construct($name) { $this->name = $name; }
		Instructions: []interface{}{
			// Fetch $this
			Instruction{
				Opcode: OpFetchThis,
				Result: Operand{Type: OpTmpVar, Value: 5},
			},
			// Fetch parameter $name (local variable 0)
			Instruction{
				Opcode: OpFetchR,
				Op1:    Operand{Type: OpCV, Value: 0}, // param 0
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
			// Assign to $this->name
			Instruction{
				Opcode: OpAssignObj,
				Op1:    Operand{Type: OpTmpVar, Value: 5}, // $this
				Op2:    Operand{Type: OpConst, Value: 0},  // "name"
				Result: Operand{Type: OpTmpVar, Value: 1}, // value
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpConst, Value: 1}, // null
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals:     10,
		NumParams:     1,
		IsConstructor: true,
	}
	classEntry.Methods["__construct"] = classEntry.Constructor

	classEntry.Properties["name"] = &types.PropertyDef{
		Name:       "name",
		Visibility: types.VisibilityPublic,
		Default:    types.NewString(""),
	}

	vm.classes["Person"] = classEntry
	vm.constants = []interface{}{"name", nil, "John", "Person", "__construct"}

	// Simulate: $person = new Person("John");
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			// new Person()
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 3}, // "Person"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			// $person->__construct("John")
			{
				Opcode:        OpInitMethodCall,
				Op1:           Operand{Type: OpTmpVar, Value: 0}, // object
				Op2:           Operand{Type: OpConst, Value: 4},  // "__construct"
				ExtendedValue: 1,                                  // 1 arg
			},
			// Send parameter "John"
			{
				Opcode: OpSendVal,
				Op1:    Operand{Type: OpConst, Value: 2}, // "John"
			},
			// Execute constructor
			{
				Opcode: OpDoFcall,
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute all instructions
	for _, instr := range mainFunc.Instructions {
		frame.ip++
		err := vm.dispatch(frame, instr)
		if err != nil {
			t.Fatalf("Instruction %v failed: %v", instr.Opcode, err)
		}
	}

	// Verify property was set to parameter value
	objVal := frame.getLocal(0)
	obj := objVal.ToObject()

	prop, exists := obj.Properties["name"]
	if !exists {
		t.Fatal("Expected 'name' property to exist")
	}

	if prop.Value.ToString() != "John" {
		t.Errorf("Expected property value 'John', got '%s'", prop.Value.ToString())
	}
}

func TestConstructor_MultipleParameters(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("User")
	classEntry.Constructor = &types.MethodDef{
		Name:       "__construct",
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
		// Constructor: function __construct($name, $age) { $this->name = $name; $this->age = $age; }
		Instructions: []interface{}{
			// Fetch $this
			Instruction{
				Opcode: OpFetchThis,
				Result: Operand{Type: OpTmpVar, Value: 5},
			},
			// $this->name = $name (param 0)
			Instruction{
				Opcode: OpFetchR,
				Op1:    Operand{Type: OpCV, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
			Instruction{
				Opcode: OpAssignObj,
				Op1:    Operand{Type: OpTmpVar, Value: 5}, // $this
				Op2:    Operand{Type: OpConst, Value: 0},  // "name"
				Result: Operand{Type: OpTmpVar, Value: 1}, // value
			},
			// $this->age = $age (param 1)
			Instruction{
				Opcode: OpFetchR,
				Op1:    Operand{Type: OpCV, Value: 1},
				Result: Operand{Type: OpTmpVar, Value: 2},
			},
			Instruction{
				Opcode: OpAssignObj,
				Op1:    Operand{Type: OpTmpVar, Value: 5}, // $this
				Op2:    Operand{Type: OpConst, Value: 1},  // "age"
				Result: Operand{Type: OpTmpVar, Value: 2}, // value
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpConst, Value: 2}, // null
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals:     10,
		NumParams:     2,
		IsConstructor: true,
	}
	classEntry.Methods["__construct"] = classEntry.Constructor

	classEntry.Properties["name"] = &types.PropertyDef{
		Name:       "name",
		Visibility: types.VisibilityPublic,
		Default:    types.NewString(""),
	}
	classEntry.Properties["age"] = &types.PropertyDef{
		Name:       "age",
		Visibility: types.VisibilityPublic,
		Default:    types.NewInt(0),
	}

	vm.classes["User"] = classEntry
	vm.constants = []interface{}{"name", "age", nil, "Alice", int64(30), "User", "__construct"}

	// Simulate: $user = new User("Alice", 30);
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 5},
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			{
				Opcode:        OpInitMethodCall,
				Op1:           Operand{Type: OpTmpVar, Value: 0},
				Op2:           Operand{Type: OpConst, Value: 6},
				ExtendedValue: 2, // 2 args
			},
			{
				Opcode: OpSendVal,
				Op1:    Operand{Type: OpConst, Value: 3}, // "Alice"
			},
			{
				Opcode: OpSendVal,
				Op1:    Operand{Type: OpConst, Value: 4}, // 30
			},
			{
				Opcode: OpDoFcall,
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute all instructions
	for _, instr := range mainFunc.Instructions {
		frame.ip++
		err := vm.dispatch(frame, instr)
		if err != nil {
			t.Fatalf("Instruction %v failed: %v", instr.Opcode, err)
		}
	}

	// Verify both properties were set
	objVal := frame.getLocal(0)
	obj := objVal.ToObject()

	nameProp, exists := obj.Properties["name"]
	if !exists {
		t.Fatal("Property 'name' does not exist")
	}
	if nameProp.Value.ToString() != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", nameProp.Value.ToString())
	}

	ageProp, exists := obj.Properties["age"]
	if !exists {
		t.Fatal("Property 'age' does not exist")
	}
	if ageProp.Value.ToInt() != 30 {
		t.Errorf("Expected age 30, got %d", ageProp.Value.ToInt())
	}
}

// ============================================================================
// Regular Function Call Tests
// ============================================================================

func TestFunctionCall_Basic(t *testing.T) {
	vm := New()

	// Create a simple function: function greet() { return "Hello"; }
	greetFunc := &CompiledFunction{
		Name: "greet",
		Instructions: Instructions{
			{
				Opcode: OpFetchConstant,
				Op1:    Operand{Type: OpConst, Value: 0}, // "Hello"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	vm.RegisterFunction("greet", greetFunc)
	vm.constants = []interface{}{"Hello", "greet"}

	// Simulate: $result = greet();
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode:        OpInitFcall,
				Op2:           Operand{Type: OpConst, Value: 1}, // "greet"
				ExtendedValue: 0,
			},
			{
				Opcode: OpDoFcall,
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute instructions
	for _, instr := range mainFunc.Instructions {
		frame.ip++
		err := vm.dispatch(frame, instr)
		if err != nil {
			t.Fatalf("Instruction %v failed: %v", instr.Opcode, err)
		}
	}

	// Verify return value
	result := frame.getLocal(0)
	if result.ToString() != "Hello" {
		t.Errorf("Expected return value 'Hello', got '%s'", result.ToString())
	}
}

func TestFunctionCall_WithParameters(t *testing.T) {
	vm := New()

	// Create function: function add($a, $b) { return $a + $b; }
	addFunc := &CompiledFunction{
		Name: "add",
		Instructions: Instructions{
			{
				Opcode: OpAdd,
				Op1:    Operand{Type: OpCV, Value: 0}, // param $a
				Op2:    Operand{Type: OpCV, Value: 1}, // param $b
				Result: Operand{Type: OpTmpVar, Value: 2},
			},
			{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpTmpVar, Value: 2},
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 2,
	}

	vm.RegisterFunction("add", addFunc)
	vm.constants = []interface{}{int64(5), int64(3), "add"}

	// Simulate: $result = add(5, 3);
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode:        OpInitFcall,
				Op2:           Operand{Type: OpConst, Value: 2}, // "add"
				ExtendedValue: 2,
			},
			{
				Opcode: OpSendVal,
				Op1:    Operand{Type: OpConst, Value: 0}, // 5
			},
			{
				Opcode: OpSendVal,
				Op1:    Operand{Type: OpConst, Value: 1}, // 3
			},
			{
				Opcode: OpDoFcall,
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute instructions
	for _, instr := range mainFunc.Instructions {
		frame.ip++
		err := vm.dispatch(frame, instr)
		if err != nil {
			t.Fatalf("Instruction %v failed: %v", instr.Opcode, err)
		}
	}

	// Verify return value
	result := frame.getLocal(0)
	if result.ToInt() != 8 {
		t.Errorf("Expected return value 8, got %d", result.ToInt())
	}
}

func TestFunctionCall_UndefinedFunction(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"nonExistent"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode:        OpInitFcall,
				Op2:           Operand{Type: OpConst, Value: 0},
				ExtendedValue: 0,
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	frame.ip++
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error for undefined function")
	}

	expectedMsg := "Call to undefined function nonExistent()"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}
