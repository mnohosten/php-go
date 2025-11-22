package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Parent Method Call Tests
// ============================================================================

func TestParentMethodCall_Basic(t *testing.T) {
	vm := New()

	// Create parent class with a method
	parent := types.NewClassEntry("ParentClass")
	parent.Methods["greet"] = &types.MethodDef{
		Name:       "greet",
		Visibility: types.VisibilityPublic,
		// return "Hello from Parent"
		Instructions: []interface{}{
			Instruction{
				Opcode: OpFetchConstant,
				Op1:    Operand{Type: OpConst, Value: 0}, // "Hello from Parent"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	// Create child class that extends parent
	child := types.NewClassEntry("ChildClass")
	child.ParentClass = parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	vm.classes["ParentClass"] = parent
	vm.classes["ChildClass"] = child
	vm.constants = []interface{}{"Hello from Parent", "parent", "greet"}

	// Create object of child class
	obj := types.NewObjectFromClass(child)

	// Simulate: parent::greet()
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			// parent::greet()
			{
				Opcode:        OpInitStaticMethodCall,
				Op1:           Operand{Type: OpConst, Value: 1}, // "parent"
				Op2:           Operand{Type: OpConst, Value: 2}, // "greet"
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
	frame.thisObject = obj
	frame.currentClass = child
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
	if result.ToString() != "Hello from Parent" {
		t.Errorf("Expected 'Hello from Parent', got '%s'", result.ToString())
	}
}

func TestParentMethodCall_OverriddenMethod(t *testing.T) {
	vm := New()

	// Create parent class with a method
	parent := types.NewClassEntry("ParentClass")
	parent.Methods["getValue"] = &types.MethodDef{
		Name:       "getValue",
		Visibility: types.VisibilityPublic,
		// return 100
		Instructions: []interface{}{
			Instruction{
				Opcode: OpFetchConstant,
				Op1:    Operand{Type: OpConst, Value: 0}, // 100
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	// Create child class that overrides the method
	child := types.NewClassEntry("ChildClass")
	child.ParentClass = parent
	child.Methods["getValue"] = &types.MethodDef{
		Name:       "getValue",
		Visibility: types.VisibilityPublic,
		// return 200
		Instructions: []interface{}{
			Instruction{
				Opcode: OpFetchConstant,
				Op1:    Operand{Type: OpConst, Value: 1}, // 200
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpTmpVar, Value: 0},
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	vm.classes["ParentClass"] = parent
	vm.classes["ChildClass"] = child
	vm.constants = []interface{}{int64(100), int64(200), "parent", "getValue"}

	// Create object of child class
	obj := types.NewObjectFromClass(child)

	// Test calling parent::getValue() should return 100
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode:        OpInitStaticMethodCall,
				Op1:           Operand{Type: OpConst, Value: 2}, // "parent"
				Op2:           Operand{Type: OpConst, Value: 3}, // "getValue"
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
	frame.thisObject = obj
	frame.currentClass = child
	vm.pushFrame(frame)

	// Execute instructions
	for _, instr := range mainFunc.Instructions {
		frame.ip++
		err := vm.dispatch(frame, instr)
		if err != nil {
			t.Fatalf("Instruction %v failed: %v", instr.Opcode, err)
		}
	}

	// Verify we got parent's version (100, not 200)
	result := frame.getLocal(0)
	if result.ToInt() != 100 {
		t.Errorf("Expected parent's value 100, got %d", result.ToInt())
	}
}

// ============================================================================
// Abstract Class Tests
// ============================================================================

func TestAbstractClass_CannotInstantiate(t *testing.T) {
	vm := New()

	// Create abstract class
	abstractClass := types.NewClassEntry("AbstractClass")
	abstractClass.IsAbstract = true
	abstractClass.Methods["abstractMethod"] = &types.MethodDef{
		Name:       "abstractMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	vm.classes["AbstractClass"] = abstractClass
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

	// Should fail to instantiate abstract class
	frame.ip++
	err := vm.dispatch(frame, mainFunc.Instructions[0])
	if err == nil {
		t.Fatal("Expected error when instantiating abstract class")
	}

	if err.Error() != "Cannot instantiate abstract class 'AbstractClass'" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAbstractClass_ConcreteChildCanInstantiate(t *testing.T) {
	vm := New()

	// Create abstract parent
	parent := types.NewClassEntry("AbstractParent")
	parent.IsAbstract = true
	parent.Methods["abstractMethod"] = &types.MethodDef{
		Name:       "abstractMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Create concrete child that implements abstract method
	child := types.NewClassEntry("ConcreteChild")
	child.ParentClass = parent
	child.Methods["abstractMethod"] = &types.MethodDef{
		Name:       "abstractMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: false, // Implemented
		Instructions: []interface{}{
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpConst, Value: 0}, // null
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	vm.classes["AbstractParent"] = parent
	vm.classes["ConcreteChild"] = child
	vm.constants = []interface{}{nil, "ConcreteChild"}

	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpNew,
				Op1:    Operand{Type: OpConst, Value: 1}, // "ConcreteChild"
				Result: Operand{Type: OpTmpVar, Value: 0},
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	vm.pushFrame(frame)

	// Should succeed - concrete child can be instantiated
	frame.ip++
	err = vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("Should be able to instantiate concrete child: %v", err)
	}

	objVal := frame.getLocal(0)
	if objVal.Type() != types.TypeObject {
		t.Error("Expected object to be created")
	}
}

// ============================================================================
// Final Class/Method Tests
// ============================================================================

func TestFinalClass_CannotExtend(t *testing.T) {
	// Create final parent
	parent := types.NewClassEntry("FinalParent")
	parent.IsFinal = true

	// Try to create child
	child := types.NewClassEntry("Child")
	child.ParentClass = parent

	err := child.InheritFrom(parent)
	if err == nil {
		t.Fatal("Expected error when extending final class")
	}

	expectedMsg := "Class Child cannot extend final class FinalParent"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestFinalMethod_CannotOverride(t *testing.T) {
	// Create parent with final method
	parent := types.NewClassEntry("Parent")
	parent.Methods["finalMethod"] = &types.MethodDef{
		Name:       "finalMethod",
		Visibility: types.VisibilityPublic,
		IsFinal:    true,
	}

	// Try to override in child
	child := types.NewClassEntry("Child")
	child.ParentClass = parent
	child.Methods["finalMethod"] = &types.MethodDef{
		Name:       "finalMethod",
		Visibility: types.VisibilityPublic,
	}

	err := child.InheritFrom(parent)
	if err == nil {
		t.Fatal("Expected error when overriding final method")
	}

	expectedMsg := "Cannot override final method Parent::finalMethod() in Child"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}
