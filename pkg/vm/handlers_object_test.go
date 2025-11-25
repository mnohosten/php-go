package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Object Property Fetch Tests
// ============================================================================

func TestOpFetchObjR_PublicProperty(t *testing.T) {
	vm := New()

	// Create class and object
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["name"] = &types.Property{
		Value:      types.NewString("John"),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	// Create frame with object and property name
	fn := &CompiledFunction{
		Instructions: Instructions{},
		NumLocals:    10,
	}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)                        // Object
	frame.setLocal(1, types.NewString("name"))       // Property name

	// Execute OpFetchObjR
	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0}, // Object
		Op2:    Operand{Type: OpTmpVar, Value: 1}, // Property name
		Result: Operand{Type: OpTmpVar, Value: 2}, // Result
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("opFetchObjR failed: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToString() != "John" {
		t.Errorf("Expected 'John', got '%s'", result.ToString())
	}
}

func TestOpFetchObjR_NonexistentProperty(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("nonexistent"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("opFetchObjR failed: %v", err)
	}

	result := frame.getLocal(2)
	if !result.IsNull() {
		t.Errorf("Expected NULL for nonexistent property, got %v", result.Type())
	}
}

func TestOpFetchObjR_NonObject(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, types.NewInt(42))              // Not an object
	frame.setLocal(1, types.NewString("property"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("opFetchObjR failed: %v", err)
	}

	result := frame.getLocal(2)
	if !result.IsNull() {
		t.Errorf("Expected NULL for non-object, got %v", result.Type())
	}
}

// ============================================================================
// Object Property Assignment Tests
// ============================================================================

func TestOpAssignObj_BasicAssignment(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)                         // Object
	frame.setLocal(1, types.NewString("name"))        // Property name
	frame.setLocal(2, types.NewString("Alice"))       // Value to assign

	instr := Instruction{
		Opcode: OpAssignObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0}, // Object
		Op2:    Operand{Type: OpTmpVar, Value: 1}, // Property name
		Result: Operand{Type: OpTmpVar, Value: 2}, // Value
	}

	err := vm.opAssignObj(frame, instr)
	if err != nil {
		t.Fatalf("opAssignObj failed: %v", err)
	}

	// Verify property was set
	prop, exists := obj.Properties["name"]
	if !exists {
		t.Fatal("Property 'name' was not set")
	}
	if prop.Value.ToString() != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", prop.Value.ToString())
	}
}

func TestOpAssignObj_AutoVivification(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, types.NewInt(42))               // Not an object
	frame.setLocal(1, types.NewString("prop"))
	frame.setLocal(2, types.NewString("value"))

	instr := Instruction{
		Opcode: OpAssignObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opAssignObj(frame, instr)
	if err != nil {
		t.Fatalf("opAssignObj failed: %v", err)
	}

	// Check that the value was converted to an object
	result := frame.getLocal(0)
	if result.Type() != types.TypeObject {
		t.Errorf("Expected auto-vivification to object, got %v", result.Type())
	}
}

// ============================================================================
// Object Property Increment/Decrement Tests
// ============================================================================

func TestOpPreIncObj(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["count"] = &types.Property{
		Value:      types.NewInt(5),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("count"))

	instr := Instruction{
		Opcode: OpPreIncObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0}, // Object
		Op2:    Operand{Type: OpTmpVar, Value: 1}, // Property name
		Result: Operand{Type: OpTmpVar, Value: 2}, // Result
	}

	err := vm.opPreIncObj(frame, instr)
	if err != nil {
		t.Fatalf("opPreIncObj failed: %v", err)
	}

	// Check result is new value (6)
	result := frame.getLocal(2)
	if result.ToInt() != 6 {
		t.Errorf("Expected result 6, got %d", result.ToInt())
	}

	// Check property was updated
	prop, _ := obj.Properties["count"]
	if prop.Value.ToInt() != 6 {
		t.Errorf("Expected property value 6, got %d", prop.Value.ToInt())
	}
}

func TestOpPostIncObj(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["count"] = &types.Property{
		Value:      types.NewInt(5),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("count"))

	instr := Instruction{
		Opcode: OpPostIncObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opPostIncObj(frame, instr)
	if err != nil {
		t.Fatalf("opPostIncObj failed: %v", err)
	}

	// Check result is old value (5)
	result := frame.getLocal(2)
	if result.ToInt() != 5 {
		t.Errorf("Expected result 5 (old value), got %d", result.ToInt())
	}

	// Check property was updated to 6
	prop, _ := obj.Properties["count"]
	if prop.Value.ToInt() != 6 {
		t.Errorf("Expected property value 6, got %d", prop.Value.ToInt())
	}
}

func TestOpPreDecObj(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["count"] = &types.Property{
		Value:      types.NewInt(10),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("count"))

	instr := Instruction{
		Opcode: OpPreDecObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opPreDecObj(frame, instr)
	if err != nil {
		t.Fatalf("opPreDecObj failed: %v", err)
	}

	// Check result is new value (9)
	result := frame.getLocal(2)
	if result.ToInt() != 9 {
		t.Errorf("Expected result 9, got %d", result.ToInt())
	}

	// Check property was updated
	prop, _ := obj.Properties["count"]
	if prop.Value.ToInt() != 9 {
		t.Errorf("Expected property value 9, got %d", prop.Value.ToInt())
	}
}

func TestOpPostDecObj(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["count"] = &types.Property{
		Value:      types.NewInt(10),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("count"))

	instr := Instruction{
		Opcode: OpPostDecObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opPostDecObj(frame, instr)
	if err != nil {
		t.Fatalf("opPostDecObj failed: %v", err)
	}

	// Check result is old value (10)
	result := frame.getLocal(2)
	if result.ToInt() != 10 {
		t.Errorf("Expected result 10 (old value), got %d", result.ToInt())
	}

	// Check property was updated to 9
	prop, _ := obj.Properties["count"]
	if prop.Value.ToInt() != 9 {
		t.Errorf("Expected property value 9, got %d", prop.Value.ToInt())
	}
}

// ============================================================================
// Object Property Unset Tests
// ============================================================================

func TestOpUnsetObj(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["temp"] = &types.Property{
		Value:      types.NewString("delete me"),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("temp"))

	instr := Instruction{
		Opcode: OpUnsetObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opUnsetObj(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetObj failed: %v", err)
	}

	// Verify property was deleted
	if _, exists := obj.Properties["temp"]; exists {
		t.Error("Property 'temp' should have been deleted")
	}
}

func TestOpUnsetObj_NonObject(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, types.NewInt(42))
	frame.setLocal(1, types.NewString("prop"))

	instr := Instruction{
		Opcode: OpUnsetObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	// Should not error on non-object (PHP behavior)
	err := vm.opUnsetObj(frame, instr)
	if err != nil {
		t.Errorf("opUnsetObj should not error on non-object: %v", err)
	}
}

// ============================================================================
// Object Property Isset/Empty Tests
// ============================================================================

func TestOpIssetIsemptyPropObj_PropertyExists(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["name"] = &types.Property{
		Value:      types.NewString("John"),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("name"))

	instr := Instruction{
		Opcode: OpIssetIsemptyPropObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opIssetIsemptyPropObj(frame, instr)
	if err != nil {
		t.Fatalf("opIssetIsemptyPropObj failed: %v", err)
	}

	result := frame.getLocal(2)
	if !result.ToBool() {
		t.Error("Expected isset to return true for existing non-null property")
	}
}

func TestOpIssetIsemptyPropObj_PropertyNull(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["nullProp"] = &types.Property{
		Value:      types.NewNull(),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("nullProp"))

	instr := Instruction{
		Opcode: OpIssetIsemptyPropObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opIssetIsemptyPropObj(frame, instr)
	if err != nil {
		t.Fatalf("opIssetIsemptyPropObj failed: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToBool() {
		t.Error("Expected isset to return false for null property")
	}
}

func TestOpIssetIsemptyPropObj_PropertyNotExists(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("nonexistent"))

	instr := Instruction{
		Opcode: OpIssetIsemptyPropObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opIssetIsemptyPropObj(frame, instr)
	if err != nil {
		t.Fatalf("opIssetIsemptyPropObj failed: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToBool() {
		t.Error("Expected isset to return false for nonexistent property")
	}
}

func TestOpIssetIsemptyPropObj_NonObject(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, types.NewInt(42))
	frame.setLocal(1, types.NewString("prop"))

	instr := Instruction{
		Opcode: OpIssetIsemptyPropObj,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opIssetIsemptyPropObj(frame, instr)
	if err != nil {
		t.Fatalf("opIssetIsemptyPropObj failed: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToBool() {
		t.Error("Expected isset to return false for non-object")
	}
}

// ============================================================================
// Object Property Compound Assignment Tests
// ============================================================================

func TestOpAssignObjOp_Addition(t *testing.T) {
	vm := New()

	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["count"] = &types.Property{
		Value:      types.NewInt(10),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("count"))
	frame.setLocal(2, types.NewInt(5)) // Add 5

	instr := Instruction{
		Opcode: OpAssignObjOp,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opAssignObjOp(frame, instr)
	if err != nil {
		t.Fatalf("opAssignObjOp failed: %v", err)
	}

	// Check property was updated (10 + 5 = 15)
	prop, _ := obj.Properties["count"]
	if prop.Value.ToInt() != 15 {
		t.Errorf("Expected property value 15, got %d", prop.Value.ToInt())
	}
}

// ============================================================================
// Visibility Tests
// ============================================================================

// TestPropertyVisibility_PublicFromOutside tests public property access from outside class
func TestPropertyVisibility_PublicFromOutside(t *testing.T) {
	vm := New()

	// Create class with public property
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["name"] = &types.Property{
		Value:      types.NewString("John"),
		Visibility: types.VisibilityPublic,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = nil // Outside class context
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("name"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("Public property should be accessible from outside: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToString() != "John" {
		t.Errorf("Expected 'John', got '%s'", result.ToString())
	}
}

// TestPropertyVisibility_PrivateFromOutside tests private property blocked from outside
func TestPropertyVisibility_PrivateFromOutside(t *testing.T) {
	vm := New()

	// Create class with private property
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["secret"] = &types.Property{
		Value:      types.NewString("hidden"),
		Visibility: types.VisibilityPrivate,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = nil // Outside class context
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("secret"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	// Should not error but should return null (PHP behavior for inaccessible property)
	if err != nil {
		t.Fatalf("opFetchObjR should not error: %v", err)
	}

	// Private property should not be accessible - returns null
	result := frame.getLocal(2)
	if result.Type() != types.TypeNull {
		t.Errorf("Private property should return null from outside, got %v", result.Type())
	}
}

// TestPropertyVisibility_PrivateFromSameClass tests private property access from same class
func TestPropertyVisibility_PrivateFromSameClass(t *testing.T) {
	vm := New()

	// Create class with private property
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	obj.Properties["secret"] = &types.Property{
		Value:      types.NewString("hidden"),
		Visibility: types.VisibilityPrivate,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = class // Inside same class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("secret"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("Private property should be accessible from same class: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToString() != "hidden" {
		t.Errorf("Expected 'hidden', got '%s'", result.ToString())
	}
}

// TestPropertyVisibility_ProtectedFromChild tests protected property access from child class
func TestPropertyVisibility_ProtectedFromChild(t *testing.T) {
	vm := New()

	// Create parent and child class
	parent := types.NewClassEntry("Parent")
	child := types.NewClassEntry("Child")
	child.ParentClass = parent

	// Create object from parent
	obj := types.NewObjectFromClass(parent)
	obj.Properties["protected_prop"] = &types.Property{
		Value:      types.NewString("protected value"),
		Visibility: types.VisibilityProtected,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = child // Inside child class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("protected_prop"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("Protected property should be accessible from child class: %v", err)
	}

	result := frame.getLocal(2)
	if result.ToString() != "protected value" {
		t.Errorf("Expected 'protected value', got '%s'", result.ToString())
	}
}

// TestPropertyVisibility_ProtectedFromUnrelated tests protected property blocked from unrelated class
func TestPropertyVisibility_ProtectedFromUnrelated(t *testing.T) {
	vm := New()

	// Create unrelated classes
	testClass := types.NewClassEntry("TestClass")
	unrelated := types.NewClassEntry("Unrelated")

	obj := types.NewObjectFromClass(testClass)
	obj.Properties["protected_prop"] = &types.Property{
		Value:      types.NewString("protected value"),
		Visibility: types.VisibilityProtected,
	}
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = unrelated // Inside unrelated class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("protected_prop"))

	instr := Instruction{
		Opcode: OpFetchObjR,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
		Result: Operand{Type: OpTmpVar, Value: 2},
	}

	err := vm.opFetchObjR(frame, instr)
	if err != nil {
		t.Fatalf("opFetchObjR should not error: %v", err)
	}

	// Protected property should not be accessible from unrelated class
	result := frame.getLocal(2)
	if result.Type() != types.TypeNull {
		t.Errorf("Protected property should return null from unrelated class, got %v", result.Type())
	}
}

// TestMethodVisibility_PublicFromOutside tests public method access from outside
func TestMethodVisibility_PublicFromOutside(t *testing.T) {
	vm := New()

	// Create class with public method
	class := types.NewClassEntry("TestClass")
	class.Methods = map[string]*types.MethodDef{
		"publicMethod": {
			Name:         "publicMethod",
			Visibility:   types.VisibilityPublic,
			Instructions: []interface{}{},
		},
	}
	vm.classes["TestClass"] = class

	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = nil // Outside class context
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("publicMethod"))

	instr := Instruction{
		Opcode: OpInitMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitMethodCall(frame, instr)
	if err != nil {
		t.Fatalf("Public method should be callable from outside: %v", err)
	}

	if frame.pendingMethod == nil {
		t.Error("Expected pending method to be set")
	}
}

// TestMethodVisibility_PrivateFromOutside tests private method blocked from outside
func TestMethodVisibility_PrivateFromOutside(t *testing.T) {
	vm := New()

	// Create class with private method
	class := types.NewClassEntry("TestClass")
	class.Methods = map[string]*types.MethodDef{
		"privateMethod": {
			Name:         "privateMethod",
			Visibility:   types.VisibilityPrivate,
			Instructions: []interface{}{},
		},
	}
	vm.classes["TestClass"] = class

	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = nil // Outside class context
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("privateMethod"))

	instr := Instruction{
		Opcode: OpInitMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitMethodCall(frame, instr)
	if err == nil {
		t.Fatal("Expected error when calling private method from outside")
	}

	// Check error message mentions "private"
	if err.Error() == "" || !errorContainsWord(err.Error(), "private") {
		t.Errorf("Expected error message to mention 'private', got: %v", err)
	}
}

// TestMethodVisibility_PrivateFromSameClass tests private method access from same class
func TestMethodVisibility_PrivateFromSameClass(t *testing.T) {
	vm := New()

	// Create class with private method
	class := types.NewClassEntry("TestClass")
	class.Methods = map[string]*types.MethodDef{
		"privateMethod": {
			Name:         "privateMethod",
			Visibility:   types.VisibilityPrivate,
			Instructions: []interface{}{},
		},
	}
	vm.classes["TestClass"] = class

	obj := types.NewObjectFromClass(class)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = class // Inside same class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("privateMethod"))

	instr := Instruction{
		Opcode: OpInitMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitMethodCall(frame, instr)
	if err != nil {
		t.Fatalf("Private method should be callable from same class: %v", err)
	}
}

// TestMethodVisibility_ProtectedFromChild tests protected method access from child
func TestMethodVisibility_ProtectedFromChild(t *testing.T) {
	vm := New()

	// Create parent and child class
	parent := types.NewClassEntry("Parent")
	parent.Methods = map[string]*types.MethodDef{
		"protectedMethod": {
			Name:         "protectedMethod",
			Visibility:   types.VisibilityProtected,
			Instructions: []interface{}{},
		},
	}
	vm.classes["Parent"] = parent

	child := types.NewClassEntry("Child")
	child.ParentClass = parent
	vm.classes["Child"] = child

	obj := types.NewObjectFromClass(parent)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = child // Inside child class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("protectedMethod"))

	instr := Instruction{
		Opcode: OpInitMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitMethodCall(frame, instr)
	if err != nil {
		t.Fatalf("Protected method should be callable from child class: %v", err)
	}
}

// TestMethodVisibility_ProtectedFromUnrelated tests protected method blocked from unrelated
func TestMethodVisibility_ProtectedFromUnrelated(t *testing.T) {
	vm := New()

	// Create unrelated classes
	testClass := types.NewClassEntry("TestClass")
	testClass.Methods = map[string]*types.MethodDef{
		"protectedMethod": {
			Name:         "protectedMethod",
			Visibility:   types.VisibilityProtected,
			Instructions: []interface{}{},
		},
	}
	vm.classes["TestClass"] = testClass

	unrelated := types.NewClassEntry("Unrelated")
	vm.classes["Unrelated"] = unrelated

	obj := types.NewObjectFromClass(testClass)
	objVal := types.NewObject(obj)

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = unrelated // Inside unrelated class
	frame.setLocal(0, objVal)
	frame.setLocal(1, types.NewString("protectedMethod"))

	instr := Instruction{
		Opcode: OpInitMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitMethodCall(frame, instr)
	if err == nil {
		t.Fatal("Expected error when calling protected method from unrelated class")
	}

	// Check error message mentions "protected"
	if err.Error() == "" || !errorContainsWord(err.Error(), "protected") {
		t.Errorf("Expected error message to mention 'protected', got: %v", err)
	}
}

// TestStaticMethodVisibility_PrivateFromOutside tests static private method blocked
func TestStaticMethodVisibility_PrivateFromOutside(t *testing.T) {
	vm := New()

	// Create class with private static method
	class := types.NewClassEntry("TestClass")
	class.Methods = map[string]*types.MethodDef{
		"privateStatic": {
			Name:         "privateStatic",
			Visibility:   types.VisibilityPrivate,
			IsStatic:     true,
			Instructions: []interface{}{},
		},
	}
	vm.classes["TestClass"] = class

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)
	frame.currentClass = nil // Outside class context
	frame.setLocal(0, types.NewString("TestClass"))
	frame.setLocal(1, types.NewString("privateStatic"))

	instr := Instruction{
		Opcode: OpInitStaticMethodCall,
		Op1:    Operand{Type: OpTmpVar, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opInitStaticMethodCall(frame, instr)
	if err == nil {
		t.Fatal("Expected error when calling private static method from outside")
	}
}

// Helper to check if error message contains a word
func errorContainsWord(s, word string) bool {
	// Simple substring check for visibility keywords
	for i := 0; i <= len(s)-len(word); i++ {
		if s[i:i+len(word)] == word {
			return true
		}
	}
	return false
}
