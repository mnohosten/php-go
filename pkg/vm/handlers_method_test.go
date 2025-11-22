package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// OpNew Tests - Object Instantiation
// ============================================================================

func TestOpNew_BasicObjectCreation(t *testing.T) {
	vm := New()

	// Register a simple class
	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	// Create main function with OpNew instruction
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			// OpNew: class name in Op1 (const), result in Result
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 0}, // "TestClass"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	vm.constants = []interface{}{"TestClass"}
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
		t.Errorf("Expected object, got %v", objVal.Type())
	}

	obj := objVal.ToObject()
	if obj == nil {
		t.Fatal("Expected object, got nil")
	}

	if obj.ClassName != "TestClass" {
		t.Errorf("Expected class name 'TestClass', got '%s'", obj.ClassName)
	}

	if obj.ClassEntry != classEntry {
		t.Error("Object should reference the class entry")
	}
}

func TestOpNew_AbstractClass(t *testing.T) {
	vm := New()

	// Register an abstract class
	classEntry := types.NewClassEntry("AbstractClass")
	classEntry.IsAbstract = true
	vm.classes["AbstractClass"] = classEntry

	vm.constants = []interface{}{"AbstractClass"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpNew - should fail for abstract class
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when instantiating abstract class")
	}

	if err.Error() != "Cannot instantiate abstract class 'AbstractClass'" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestOpNew_InterfaceClass(t *testing.T) {
	vm := New()

	// Register an interface
	classEntry := types.NewClassEntry("TestInterface")
	classEntry.IsInterface = true
	vm.classes["TestInterface"] = classEntry

	vm.constants = []interface{}{"TestInterface"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpNew - should fail for interface
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when instantiating interface")
	}

	if err.Error() != "Cannot instantiate interface 'TestInterface'" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestOpNew_ClassNotFound(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"NonExistentClass"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpNew - should fail for non-existent class
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when class not found")
	}

	if err.Error() != "Class 'NonExistentClass' not found" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// ============================================================================
// OpInitMethodCall Tests - Instance Method Calls
// ============================================================================

func TestOpInitMethodCall_BasicMethod(t *testing.T) {
	vm := New()

	// Create class with a method
	classEntry := types.NewClassEntry("TestClass")
	classEntry.Methods["testMethod"] = &types.MethodDef{
		Name:       "testMethod",
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}
	vm.classes["TestClass"] = classEntry

	// Create object
	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"testMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitMethodCall,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // object
				Op2:    Operand{Type: OpConst, Value: 0},  // "testMethod"
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInitMethodCall
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInitMethodCall failed: %v", err)
	}

	// Verify pending method was set
	if frame.pendingMethod == nil {
		t.Fatal("Expected pending method to be set")
	}

	if frame.pendingMethod.Name != "testMethod" {
		t.Errorf("Expected method name 'testMethod', got '%s'", frame.pendingMethod.Name)
	}

	// Verify pending object was set
	if frame.pendingObject != obj {
		t.Error("Expected pending object to be set to the object")
	}
}

func TestOpInitMethodCall_MethodNotFound(t *testing.T) {
	vm := New()

	// Create class without the method
	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"nonExistentMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitMethodCall,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInitMethodCall - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when method not found")
	}
}

func TestOpInitMethodCall_NonObject(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"testMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitMethodCall,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, types.NewInt(42)) // Not an object
	vm.pushFrame(frame)

	// Execute OpInitMethodCall - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when calling method on non-object")
	}
}

// ============================================================================
// OpInitStaticMethodCall Tests - Static Method Calls
// ============================================================================

func TestOpInitStaticMethodCall_BasicStaticMethod(t *testing.T) {
	vm := New()

	// Create class with a static method
	classEntry := types.NewClassEntry("TestClass")
	classEntry.Methods["staticMethod"] = &types.MethodDef{
		Name:       "staticMethod",
		Visibility: types.VisibilityPublic,
		IsStatic:   true,
	}
	vm.classes["TestClass"] = classEntry

	vm.constants = []interface{}{"TestClass", "staticMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitStaticMethodCall,
				Op1:    Operand{Type: OpConst, Value: 0}, // "TestClass"
				Op2:    Operand{Type: OpConst, Value: 1}, // "staticMethod"
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpInitStaticMethodCall
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInitStaticMethodCall failed: %v", err)
	}

	// Verify pending method was set
	if frame.pendingMethod == nil {
		t.Fatal("Expected pending method to be set")
	}

	if frame.pendingMethod.Name != "staticMethod" {
		t.Errorf("Expected method name 'staticMethod', got '%s'", frame.pendingMethod.Name)
	}

	// Verify pending object is nil (static call)
	if frame.pendingObject != nil {
		t.Error("Expected pending object to be nil for static call")
	}
}

func TestOpInitStaticMethodCall_ClassNotFound(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"NonExistentClass", "staticMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitStaticMethodCall,
				Op1:    Operand{Type: OpConst, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpInitStaticMethodCall - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when class not found")
	}
}

func TestOpInitStaticMethodCall_MethodNotFound(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	vm.constants = []interface{}{"TestClass", "nonExistentMethod"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInitStaticMethodCall,
				Op1:    Operand{Type: OpConst, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Execute OpInitStaticMethodCall - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when method not found")
	}
}

// ============================================================================
// OpClone Tests - Object Cloning
// ============================================================================

func TestOpClone_BasicCloning(t *testing.T) {
	vm := New()

	// Create class with a property
	classEntry := types.NewClassEntry("TestClass")
	classEntry.Properties["prop"] = &types.PropertyDef{
		Name:       "prop",
		Visibility: types.VisibilityPublic,
		Default:    types.NewInt(42),
	}
	vm.classes["TestClass"] = classEntry

	// Create object with a property value
	obj := types.NewObjectFromClass(classEntry)
	obj.Properties["prop"].Value = types.NewInt(100)
	objVal := types.NewObject(obj)

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpClone,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // object to clone
				Result: Operand{Type: OpTmpVar, Value: 1}, // cloned object
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpClone
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpClone failed: %v", err)
	}

	// Verify cloned object was created
	clonedVal := frame.getLocal(1)
	if clonedVal.Type() != types.TypeObject {
		t.Fatalf("Expected object, got %v", clonedVal.Type())
	}

	cloned := clonedVal.ToObject()
	if cloned == nil {
		t.Fatal("Expected cloned object, got nil")
	}

	// Verify class name is the same
	if cloned.ClassName != "TestClass" {
		t.Errorf("Expected class name 'TestClass', got '%s'", cloned.ClassName)
	}

	// Verify object IDs are different
	if cloned.ObjectID == obj.ObjectID {
		t.Error("Cloned object should have a different ObjectID")
	}

	// Verify property was cloned
	if prop, exists := cloned.Properties["prop"]; !exists {
		t.Error("Expected property 'prop' in cloned object")
	} else {
		if prop.Value.ToInt() != 100 {
			t.Errorf("Expected property value 100, got %d", prop.Value.ToInt())
		}
	}

	// Verify modifying cloned property doesn't affect original
	cloned.Properties["prop"].Value = types.NewInt(200)
	if obj.Properties["prop"].Value.ToInt() != 100 {
		t.Error("Modifying cloned property should not affect original")
	}
}

func TestOpClone_NonObject(t *testing.T) {
	vm := New()

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpClone,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, types.NewInt(42)) // Not an object
	vm.pushFrame(frame)

	// Execute OpClone - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when cloning non-object")
	}
}

// ============================================================================
// OpInstanceof Tests - Type Checking
// ============================================================================

func TestOpInstanceof_SameClass(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"TestClass"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // object
				Op2:    Operand{Type: OpConst, Value: 0},  // "TestClass"
				Result: Operand{Type: OpTmpVar, Value: 1}, // result
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInstanceof
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInstanceof failed: %v", err)
	}

	// Verify result is true
	result := frame.getLocal(1)
	if !result.ToBool() {
		t.Error("Expected instanceof to return true for same class")
	}
}

func TestOpInstanceof_ParentClass(t *testing.T) {
	vm := New()

	// Create parent and child classes
	parentClass := types.NewClassEntry("ParentClass")
	childClass := types.NewClassEntry("ChildClass")
	childClass.ParentClass = parentClass

	vm.classes["ParentClass"] = parentClass
	vm.classes["ChildClass"] = childClass

	// Create object of child class
	obj := types.NewObjectFromClass(childClass)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"ParentClass"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInstanceof
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInstanceof failed: %v", err)
	}

	// Verify result is true (child is instance of parent)
	result := frame.getLocal(1)
	if !result.ToBool() {
		t.Error("Expected instanceof to return true for parent class")
	}
}

func TestOpInstanceof_UnrelatedClass(t *testing.T) {
	vm := New()

	class1 := types.NewClassEntry("Class1")
	class2 := types.NewClassEntry("Class2")

	vm.classes["Class1"] = class1
	vm.classes["Class2"] = class2

	obj := types.NewObjectFromClass(class1)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"Class2"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInstanceof
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInstanceof failed: %v", err)
	}

	// Verify result is false
	result := frame.getLocal(1)
	if result.ToBool() {
		t.Error("Expected instanceof to return false for unrelated class")
	}
}

func TestOpInstanceof_NonObject(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	vm.constants = []interface{}{"TestClass"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, types.NewInt(42)) // Not an object
	vm.pushFrame(frame)

	// Execute OpInstanceof
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInstanceof failed: %v", err)
	}

	// Verify result is false
	result := frame.getLocal(1)
	if result.ToBool() {
		t.Error("Expected instanceof to return false for non-object")
	}
}

func TestOpInstanceof_Interface(t *testing.T) {
	vm := New()

	// Create interface and implementing class
	iface := types.NewInterfaceEntry("TestInterface")
	classEntry := types.NewClassEntry("TestClass")
	classEntry.Interfaces = []*types.InterfaceEntry{iface}

	vm.classes["TestInterface"] = (*types.ClassEntry)(nil) // Interfaces stored separately in real implementation
	vm.classes["TestClass"] = classEntry

	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	vm.constants = []interface{}{"TestInterface"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Op2:    Operand{Type: OpConst, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpInstanceof
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpInstanceof failed: %v", err)
	}

	// Verify result is true (class implements interface)
	result := frame.getLocal(1)
	if !result.ToBool() {
		t.Error("Expected instanceof to return true for implemented interface")
	}
}

// ============================================================================
// OpGetClass Tests - Get Class Name
// ============================================================================

func TestOpGetClass_BasicObject(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpGetClass,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // object
				Result: Operand{Type: OpTmpVar, Value: 1}, // class name
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	// Execute OpGetClass
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpGetClass failed: %v", err)
	}

	// Verify result is the class name
	result := frame.getLocal(1)
	if result.Type() != types.TypeString {
		t.Fatalf("Expected string, got %v", result.Type())
	}

	if result.ToString() != "TestClass" {
		t.Errorf("Expected class name 'TestClass', got '%s'", result.ToString())
	}
}

func TestOpGetClass_NonObject(t *testing.T) {
	vm := New()

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpGetClass,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpTmpVar, Value: 1},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, types.NewInt(42)) // Not an object
	vm.pushFrame(frame)

	// Execute OpGetClass - should return false for non-object
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpGetClass failed: %v", err)
	}

	// Verify result is false (PHP behavior)
	result := frame.getLocal(1)
	if result.Type() != types.TypeBool {
		t.Fatalf("Expected bool, got %v", result.Type())
	}

	if result.ToBool() {
		t.Error("Expected get_class() to return false for non-object")
	}
}

// ============================================================================
// OpFetchThis Tests - Fetch $this
// ============================================================================

func TestOpFetchThis_WithThisObject(t *testing.T) {
	vm := New()

	classEntry := types.NewClassEntry("TestClass")
	vm.classes["TestClass"] = classEntry

	obj := types.NewObjectFromClass(classEntry)

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpFetchThis,
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.thisObject = obj // Set $this
	vm.pushFrame(frame)

	// Execute OpFetchThis
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("OpFetchThis failed: %v", err)
	}

	// Verify result is the object
	result := frame.getLocal(0)
	if result.Type() != types.TypeObject {
		t.Fatalf("Expected object, got %v", result.Type())
	}

	resultObj := result.ToObject()
	if resultObj != obj {
		t.Error("Expected $this to be the frame's thisObject")
	}
}

func TestOpFetchThis_WithoutThisObject(t *testing.T) {
	vm := New()

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpFetchThis,
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	// Don't set thisObject - simulates static context
	vm.pushFrame(frame)

	// Execute OpFetchThis - should fail
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when fetching $this in static context")
	}
}
