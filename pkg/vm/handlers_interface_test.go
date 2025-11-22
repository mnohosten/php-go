package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// instanceof with Interfaces Tests
// ============================================================================

func TestInstanceof_Interface(t *testing.T) {
	vm := New()

	// Create interface
	iface := types.NewInterfaceEntry("Countable")
	iface.Methods["count"] = &types.MethodDef{
		Name:       "count",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
		NumParams:  0,
	}

	// Create class implementing interface
	class := types.NewClassEntry("MyCollection")
	class.Interfaces = []*types.InterfaceEntry{iface}
	class.Methods["count"] = &types.MethodDef{
		Name:       "count",
		Visibility: types.VisibilityPublic,
		NumParams:  0,
		Instructions: []interface{}{
			Instruction{
				Opcode: OpReturn,
				Op1:    Operand{Type: OpConst, Value: 0}, // return 5
				Result: Operand{Type: OpUnused},
			},
		},
		NumLocals: 10,
	}

	// Validate interface implementation
	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Fatalf("Interface validation failed: %v", err)
	}

	vm.classes["MyCollection"] = class
	vm.constants = []interface{}{int64(5), "Countable"}

	// Create object
	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	// Test: $obj instanceof Countable
	mainFunc := &CompiledFunction{
		Name: "main",
		Instructions: Instructions{
			{
				Opcode: OpInstanceof,
				Op1:    Operand{Type: OpTmpVar, Value: 0}, // object
				Op2:    Operand{Type: OpConst, Value: 1},  // "Countable"
				Result: Operand{Type: OpTmpVar, Value: 1}, // result
			},
		},
		NumLocals: 10,
		NumParams: 0,
	}

	frame := NewFrame(mainFunc)
	frame.setLocal(0, objVal)
	vm.pushFrame(frame)

	frame.ip++
	err = vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("instanceof failed: %v", err)
	}

	// Verify result is true
	result := frame.getLocal(1)
	if !result.ToBool() {
		t.Error("Expected instanceof to return true for implemented interface")
	}
}

func TestInstanceof_MultipleInterfaces(t *testing.T) {
	vm := New()

	// Create two interfaces
	iface1 := types.NewInterfaceEntry("Countable")
	iface1.Methods["count"] = &types.MethodDef{
		Name:       "count",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	iface2 := types.NewInterfaceEntry("Serializable")
	iface2.Methods["serialize"] = &types.MethodDef{
		Name:       "serialize",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Create class implementing both
	class := types.NewClassEntry("MyClass")
	class.Interfaces = []*types.InterfaceEntry{iface1, iface2}
	class.Methods["count"] = &types.MethodDef{
		Name:       "count",
		Visibility: types.VisibilityPublic,
	}
	class.Methods["serialize"] = &types.MethodDef{
		Name:       "serialize",
		Visibility: types.VisibilityPublic,
	}

	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Fatalf("Interface validation failed: %v", err)
	}

	vm.classes["MyClass"] = class
	vm.constants = []interface{}{"Countable", "Serializable"}

	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	// Test instanceof for both interfaces
	for i, ifaceName := range []string{"Countable", "Serializable"} {
		mainFunc := &CompiledFunction{
			Name: "main",
			Instructions: Instructions{
				{
					Opcode: OpInstanceof,
					Op1:    Operand{Type: OpTmpVar, Value: 0},
					Op2:    Operand{Type: OpConst, Value: uint32(i)},
					Result: Operand{Type: OpTmpVar, Value: 1},
				},
			},
			NumLocals: 10,
			NumParams: 0,
		}

		frame := NewFrame(mainFunc)
		frame.setLocal(0, objVal)
		vm.pushFrame(frame)

		frame.ip++
		err = vm.dispatch(frame, mainFunc.Instructions[0])
		if err != nil {
			t.Fatalf("instanceof %s failed: %v", ifaceName, err)
		}

		result := frame.getLocal(1)
		if !result.ToBool() {
			t.Errorf("Expected instanceof %s to return true", ifaceName)
		}

		vm.popFrame()
	}
}

func TestInstanceof_InheritedInterface(t *testing.T) {
	vm := New()

	// Create interface
	iface := types.NewInterfaceEntry("TestInterface")
	iface.Methods["test"] = &types.MethodDef{
		Name:       "test",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Parent class implements interface
	parent := types.NewClassEntry("ParentClass")
	parent.Interfaces = []*types.InterfaceEntry{iface}
	parent.Methods["test"] = &types.MethodDef{
		Name:       "test",
		Visibility: types.VisibilityPublic,
	}

	// Child class extends parent
	child := types.NewClassEntry("ChildClass")
	child.ParentClass = parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	vm.classes["ChildClass"] = child
	vm.constants = []interface{}{"TestInterface"}

	// Create object of child class
	obj := types.NewObjectFromClass(child)
	objVal := types.NewObject(obj)

	// Child should implement interface through parent
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

	frame.ip++
	err = vm.dispatch(frame, mainFunc.Instructions[0])
	if err != nil {
		t.Fatalf("instanceof failed: %v", err)
	}

	result := frame.getLocal(1)
	if !result.ToBool() {
		t.Error("Expected child to implement interface through parent")
	}
}

func TestInstanceof_ExtendedInterface(t *testing.T) {
	vm := New()

	// Create parent interface
	parentIface := types.NewInterfaceEntry("ParentInterface")
	parentIface.Methods["parentMethod"] = &types.MethodDef{
		Name:       "parentMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Create child interface extending parent
	childIface := types.NewInterfaceEntry("ChildInterface")
	childIface.ParentInterfaces = []*types.InterfaceEntry{parentIface}
	childIface.Methods["childMethod"] = &types.MethodDef{
		Name:       "childMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Create class implementing child interface
	class := types.NewClassEntry("TestClass")
	class.Interfaces = []*types.InterfaceEntry{childIface}
	class.Methods["parentMethod"] = &types.MethodDef{
		Name:       "parentMethod",
		Visibility: types.VisibilityPublic,
	}
	class.Methods["childMethod"] = &types.MethodDef{
		Name:       "childMethod",
		Visibility: types.VisibilityPublic,
	}

	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Fatalf("Interface validation failed: %v", err)
	}

	vm.classes["TestClass"] = class
	vm.constants = []interface{}{"ParentInterface", "ChildInterface"}

	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	// Test instanceof for both parent and child interfaces
	for i, ifaceName := range []string{"ParentInterface", "ChildInterface"} {
		mainFunc := &CompiledFunction{
			Name: "main",
			Instructions: Instructions{
				{
					Opcode: OpInstanceof,
					Op1:    Operand{Type: OpTmpVar, Value: 0},
					Op2:    Operand{Type: OpConst, Value: uint32(i)},
					Result: Operand{Type: OpTmpVar, Value: 1},
				},
			},
			NumLocals: 10,
			NumParams: 0,
		}

		frame := NewFrame(mainFunc)
		frame.setLocal(0, objVal)
		vm.pushFrame(frame)

		frame.ip++
		err = vm.dispatch(frame, mainFunc.Instructions[0])
		if err != nil {
			t.Fatalf("instanceof %s failed: %v", ifaceName, err)
		}

		result := frame.getLocal(1)
		if !result.ToBool() {
			t.Errorf("Expected instanceof %s to return true", ifaceName)
		}

		vm.popFrame()
	}
}

// ============================================================================
// Interface Validation Tests
// ============================================================================

func TestInterface_ValidationOnInstantiation(t *testing.T) {
	// This test ensures validation happens at the right time
	// In PHP, validation happens at class declaration time, not instantiation

	// Create interface
	iface := types.NewInterfaceEntry("TestInterface")
	iface.Methods["requiredMethod"] = &types.MethodDef{
		Name:       "requiredMethod",
		Visibility: types.VisibilityPublic,
		IsAbstract: true,
	}

	// Create class claiming to implement interface but missing method
	class := types.NewClassEntry("BadClass")
	class.Interfaces = []*types.InterfaceEntry{iface}
	// Intentionally not implementing requiredMethod

	// Validation should fail
	err := class.ValidateInterfaceImplementation()
	if err == nil {
		t.Fatal("Expected validation error for missing interface method")
	}

	expectedMsg := "Class BadClass must implement method requiredMethod() from interface TestInterface"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInterface_ConstantAccess(t *testing.T) {
	// Create interface with constants
	iface := types.NewInterfaceEntry("ConfigInterface")
	iface.Constants["MAX_SIZE"] = &types.ClassConstant{
		Name:       "MAX_SIZE",
		Value:      types.NewInt(100),
		Visibility: types.VisibilityPublic,
	}

	// In PHP, interface constants can be accessed as InterfaceName::CONSTANT
	// This is tested through the class system
	if len(iface.Constants) != 1 {
		t.Error("Interface should have 1 constant")
	}

	constant := iface.Constants["MAX_SIZE"]
	if constant.Value.ToInt() != 100 {
		t.Errorf("Expected constant value 100, got %d", constant.Value.ToInt())
	}
}
