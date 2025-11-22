package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Object Property Opcode Handlers
// ============================================================================

// opFetchObjR handles fetching object property for read: result = $obj->prop
// OpFetchObjR - Fetch object property for read
func (vm *VM) opFetchObjR(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Check if it's an object
	if objVal.Type() != types.TypeObject {
		// PHP allows property access on non-objects (returns null with warning)
		return vm.setOperandValue(frame, instr.Result, types.NewNull())
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	// Get current class context for visibility checking
	// For now, assume public access (nil context)
	// TODO: Track current class context in frame
	var accessContext *types.ClassEntry = nil

	// Get property value
	value, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		// Property doesn't exist or is not accessible
		// Check for __get magic method
		if obj.ClassEntry != nil {
			if magicGet, hasMagic := obj.ClassEntry.MagicMethods["__get"]; hasMagic {
				// TODO: Call __get($name) magic method
				_ = magicGet
				// For now, return null
				return vm.setOperandValue(frame, instr.Result, types.NewNull())
			}
		}
		// No magic method, return null (PHP behavior for undefined property)
		return vm.setOperandValue(frame, instr.Result, types.NewNull())
	}

	return vm.setOperandValue(frame, instr.Result, value)
}

// opFetchObjW handles fetching object property for write: $obj->prop = ...
// OpFetchObjW - Fetch object property for write
func (vm *VM) opFetchObjW(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Auto-vivify to object if needed
	if objVal.Type() != types.TypeObject {
		// In PHP, writing to a property of a non-object converts it to stdClass
		// For now, create a simple object
		newObj := types.NewObjectInstance("stdClass")
		objVal = types.NewObject(newObj)
		// Update the original variable
		vm.setOperandValue(frame, instr.Op1, objVal)
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	// Get current class context for visibility checking
	var accessContext *types.ClassEntry = nil

	// Check if property exists
	value, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		// Property doesn't exist - create it or call __set magic method
		if obj.ClassEntry != nil {
			if magicSet, hasMagic := obj.ClassEntry.MagicMethods["__set"]; hasMagic {
				// TODO: Call __set($name, $value) magic method
				_ = magicSet
			}
		}
		// Create new property with null value
		value = types.NewNull()
		obj.SetProperty(propNameStr, value, accessContext)
	}

	// Return the property value (for write fetch, we return reference to the property)
	return vm.setOperandValue(frame, instr.Result, value)
}

// opFetchObjRW handles fetching object property for read-write: $obj->prop += 1
// OpFetchObjRW - Fetch object property for read-write
func (vm *VM) opFetchObjRW(frame *Frame, instr Instruction) error {
	// For read-write operations, we need both read and write access
	return vm.opFetchObjW(frame, instr)
}

// opFetchObjIs handles fetching object property for isset/empty check
// OpFetchObjIs - Fetch object property for isset/empty check
func (vm *VM) opFetchObjIs(frame *Frame, instr Instruction) error {
	// Similar to FetchObjR but doesn't generate notices
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return vm.setOperandValue(frame, instr.Result, types.NewNull())
	}

	obj := objVal.ToObject()

	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	// Check for __isset magic method first
	if obj.ClassEntry != nil {
		if magicIsset, hasMagic := obj.ClassEntry.MagicMethods["__isset"]; hasMagic {
			// TODO: Call __isset($name) magic method
			_ = magicIsset
		}
	}

	value, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		return vm.setOperandValue(frame, instr.Result, types.NewNull())
	}

	return vm.setOperandValue(frame, instr.Result, value)
}

// opFetchObjFuncArg handles fetching object property as function argument
// OpFetchObjFuncArg - Fetch object property as function argument
func (vm *VM) opFetchObjFuncArg(frame *Frame, instr Instruction) error {
	// For function arguments, use read semantics
	return vm.opFetchObjR(frame, instr)
}

// opFetchObjUnset handles fetching object property for unset
// OpFetchObjUnset - Fetch object property for unset
func (vm *VM) opFetchObjUnset(frame *Frame, instr Instruction) error {
	// For unset, we just need to identify the property
	return vm.opFetchObjR(frame, instr)
}

// ============================================================================
// Object Property Assignment Operations
// ============================================================================

// opAssignObj handles object property assignment: $obj->prop = value
// OpAssignObj - Object property assignment
func (vm *VM) opAssignObj(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Auto-vivify to object if needed
	if objVal.Type() != types.TypeObject {
		newObj := types.NewObjectInstance("stdClass")
		objVal = types.NewObject(newObj)
		vm.setOperandValue(frame, instr.Op1, objVal)
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	// Get the value to assign
	value, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	// Get current class context for visibility checking
	var accessContext *types.ClassEntry = nil

	// Check for __set magic method
	if obj.ClassEntry != nil {
		if magicSet, hasMagic := obj.ClassEntry.MagicMethods["__set"]; hasMagic {
			// If property doesn't exist or is not accessible, use __set
			if _, exists := obj.GetProperty(propNameStr, accessContext); !exists {
				// TODO: Call __set($name, $value) magic method
				_ = magicSet
				// For now, fall through to direct assignment
			}
		}
	}

	// Set the property
	obj.SetProperty(propNameStr, value, accessContext)

	return nil
}

// opAssignObjOp handles compound object assignment: $obj->prop += value
// OpAssignObjOp - Compound object assignment
func (vm *VM) opAssignObjOp(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("ASSIGN_OBJ_OP: not an object")
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	// Get the current value
	currentVal, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		currentVal = types.NewNull()
	}

	// Get the operand value
	operand, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	// Perform the operation (this would come from extended opcode info)
	// For simplicity, assume += operation
	// TODO: Get actual operation type from instruction
	newVal := types.NewInt(currentVal.ToInt() + operand.ToInt())

	obj.SetProperty(propNameStr, newVal, accessContext)
	return nil
}

// opAssignObjRef handles object property assignment by reference: $obj->prop =& $var
// OpAssignObjRef - Object property assignment by reference
func (vm *VM) opAssignObjRef(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("ASSIGN_OBJ_REF: not an object")
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	// Get the value to assign (should be a reference)
	value, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	var accessContext *types.ClassEntry = nil

	// Set the property (references handled by value system)
	obj.SetProperty(propNameStr, value, accessContext)

	return nil
}

// ============================================================================
// Object Property Unset Operations
// ============================================================================

// opUnsetObj handles unsetting object property: unset($obj->prop)
// OpUnsetObj - Unset object property
func (vm *VM) opUnsetObj(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		// Unset on non-object is a no-op
		return nil
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	// Check for __unset magic method
	if obj.ClassEntry != nil {
		if magicUnset, hasMagic := obj.ClassEntry.MagicMethods["__unset"]; hasMagic {
			// TODO: Call __unset($name) magic method
			_ = magicUnset
			// For now, fall through to direct unset
		}
	}

	// Remove the property
	delete(obj.Properties, propNameStr)

	return nil
}

// ============================================================================
// Object Property Isset/Empty Operations
// ============================================================================

// opIssetIsemptyPropObj handles isset/empty check on object property
// OpIssetIsemptyPropObj - Check isset/empty on object property
func (vm *VM) opIssetIsemptyPropObj(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		// Non-object is considered not set
		return vm.setOperandValue(frame, instr.Result, types.NewBool(false))
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	// Check for __isset magic method
	if obj.ClassEntry != nil {
		if magicIsset, hasMagic := obj.ClassEntry.MagicMethods["__isset"]; hasMagic {
			// TODO: Call __isset($name) magic method and return result
			_ = magicIsset
			// For now, fall through to direct check
		}
	}

	value, exists := obj.GetProperty(propNameStr, accessContext)

	var result bool
	// For isset: check if exists and not null
	// For empty: check if exists and is "empty" (falsy)
	// TODO: Determine from instruction if this is isset or empty check
	// For now, implement isset semantics
	result = exists && !value.IsNull()

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// ============================================================================
// Object Increment/Decrement Operations
// ============================================================================

// opPreIncObj handles pre-increment object property: result = ++$obj->prop
// OpPreIncObj - Pre-increment object property
func (vm *VM) opPreIncObj(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("PRE_INC_OBJ: not an object")
	}

	obj := objVal.ToObject()

	// Get the property name
	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	// Get current value
	currentVal, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		currentVal = types.NewInt(0)
	}

	// Increment
	newVal := types.NewInt(currentVal.ToInt() + 1)

	// Set back
	obj.SetProperty(propNameStr, newVal, accessContext)

	// Return new value
	return vm.setOperandValue(frame, instr.Result, newVal)
}

// opPreDecObj handles pre-decrement object property: result = --$obj->prop
// OpPreDecObj - Pre-decrement object property
func (vm *VM) opPreDecObj(frame *Frame, instr Instruction) error {
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("PRE_DEC_OBJ: not an object")
	}

	obj := objVal.ToObject()

	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	currentVal, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		currentVal = types.NewInt(0)
	}

	newVal := types.NewInt(currentVal.ToInt() - 1)
	obj.SetProperty(propNameStr, newVal, accessContext)

	return vm.setOperandValue(frame, instr.Result, newVal)
}

// opPostIncObj handles post-increment object property: result = $obj->prop++
// OpPostIncObj - Post-increment object property
func (vm *VM) opPostIncObj(frame *Frame, instr Instruction) error {
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("POST_INC_OBJ: not an object")
	}

	obj := objVal.ToObject()

	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	currentVal, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		currentVal = types.NewInt(0)
	}

	// Return old value
	oldVal := currentVal.Copy()

	// Increment and set back
	newVal := types.NewInt(currentVal.ToInt() + 1)
	obj.SetProperty(propNameStr, newVal, accessContext)

	return vm.setOperandValue(frame, instr.Result, oldVal)
}

// opPostDecObj handles post-decrement object property: result = $obj->prop--
// OpPostDecObj - Post-decrement object property
func (vm *VM) opPostDecObj(frame *Frame, instr Instruction) error {
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("POST_DEC_OBJ: not an object")
	}

	obj := objVal.ToObject()

	propName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	propNameStr := propName.ToString()

	var accessContext *types.ClassEntry = nil

	currentVal, exists := obj.GetProperty(propNameStr, accessContext)
	if !exists {
		currentVal = types.NewInt(0)
	}

	// Return old value
	oldVal := currentVal.Copy()

	// Decrement and set back
	newVal := types.NewInt(currentVal.ToInt() - 1)
	obj.SetProperty(propNameStr, newVal, accessContext)

	return vm.setOperandValue(frame, instr.Result, oldVal)
}

// ============================================================================
// Object Creation and Method Call Operations
// ============================================================================

// opNew handles object instantiation: result = new Class()
// OpNew - Create new object instance
func (vm *VM) opNew(frame *Frame, instr Instruction) error {
	// Get the class name
	className, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	classNameStr := className.ToString()

	// Look up the class in the VM's class registry
	classEntry, exists := vm.classes[classNameStr]
	if !exists {
		// Class not found - in PHP this is a fatal error
		return fmt.Errorf("Class '%s' not found", classNameStr)
	}

	// Check if class is abstract or interface
	if classEntry.IsAbstract {
		return fmt.Errorf("Cannot instantiate abstract class '%s'", classNameStr)
	}
	if classEntry.IsInterface {
		return fmt.Errorf("Cannot instantiate interface '%s'", classNameStr)
	}

	// Create new object instance
	obj := types.NewObjectFromClass(classEntry)
	objVal := types.NewObject(obj)

	// Store the object in the result operand
	// The constructor will be called separately via OpInitMethodCall + OpDoFcall
	return vm.setOperandValue(frame, instr.Result, objVal)
}

// opInitMethodCall handles initialization of instance method call: $obj->method()
// OpInitMethodCall - Initialize method call
// Op1: object
// Op2: method name (constant or variable)
// ExtendedValue: number of arguments
func (vm *VM) opInitMethodCall(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("INIT_METHOD_CALL: not an object")
	}

	obj := objVal.ToObject()

	// Get the method name
	methodName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	methodNameStr := methodName.ToString()

	// Check for class entry
	if obj.ClassEntry == nil {
		// Check for __call magic method
		return fmt.Errorf("INIT_METHOD_CALL: object has no class entry")
	}

	// Look up the method in the class hierarchy
	method, exists := obj.ClassEntry.GetMethod(methodNameStr)
	if !exists {
		// Check for __call magic method
		if magicCall, hasMagic := obj.ClassEntry.MagicMethods["__call"]; hasMagic {
			// TODO: Set up __call($method, $args) invocation
			_ = magicCall
			return fmt.Errorf("INIT_METHOD_CALL: method '%s' not found (magic method __call not yet implemented)", methodNameStr)
		}
		return fmt.Errorf("INIT_METHOD_CALL: method '%s' not found in class '%s'", methodNameStr, obj.ClassEntry.Name)
	}

	// Check if method is static (cannot call static method as instance method in strict mode)
	// In PHP, you can call static methods on instances, but we'll allow it for now

	// TODO: Check method visibility based on current context
	// For now, we'll assume all methods are accessible

	// Store method information for OpDoFcall
	// We'll use frame locals to pass this information
	// This is a simplified approach - a real implementation would use a call stack
	frame.pendingMethod = method
	frame.pendingObject = obj

	return nil
}

// opInitStaticMethodCall handles initialization of static method call: Class::method()
// OpInitStaticMethodCall - Initialize static method call
// Op1: class name (constant or variable)
// Op2: method name (constant or variable)
// ExtendedValue: number of arguments
func (vm *VM) opInitStaticMethodCall(frame *Frame, instr Instruction) error {
	// Get the class name
	className, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}
	classNameStr := className.ToString()

	// Special handling for 'self', 'parent', 'static'
	// TODO: Implement late static binding for 'static'
	switch classNameStr {
	case "self":
		// Use current class
		if frame.currentClass == nil {
			return fmt.Errorf("INIT_STATIC_METHOD_CALL: 'self' used outside class context")
		}
		classNameStr = frame.currentClass.Name
	case "parent":
		// Use parent class
		if frame.currentClass == nil || frame.currentClass.ParentClass == nil {
			return fmt.Errorf("INIT_STATIC_METHOD_CALL: 'parent' used without parent class")
		}
		classNameStr = frame.currentClass.ParentClass.Name
	case "static":
		// Late static binding - use the called class
		// For now, use current class (will implement LSB in later task)
		if frame.currentClass == nil {
			return fmt.Errorf("INIT_STATIC_METHOD_CALL: 'static' used outside class context")
		}
		classNameStr = frame.currentClass.Name
	}

	// Look up the class
	classEntry, exists := vm.classes[classNameStr]
	if !exists {
		return fmt.Errorf("INIT_STATIC_METHOD_CALL: class '%s' not found", classNameStr)
	}

	// Get the method name
	methodName, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	methodNameStr := methodName.ToString()

	// Look up the method
	method, exists := classEntry.GetMethod(methodNameStr)
	if !exists {
		// Check for __callStatic magic method
		if magicCallStatic, hasMagic := classEntry.MagicMethods["__callStatic"]; hasMagic {
			// TODO: Set up __callStatic($method, $args) invocation
			_ = magicCallStatic
			return fmt.Errorf("INIT_STATIC_METHOD_CALL: method '%s' not found (magic method __callStatic not yet implemented)", methodNameStr)
		}
		return fmt.Errorf("INIT_STATIC_METHOD_CALL: static method '%s' not found in class '%s'", methodNameStr, classNameStr)
	}

	// Check if method is actually static
	if !method.IsStatic {
		// PHP allows calling instance methods statically in some cases
		// We'll allow it for now with a warning
		// TODO: Add warning/notice system
	}

	// Store method information for OpDoFcall
	frame.pendingMethod = method
	frame.pendingObject = nil // No object for static calls

	return nil
}

// opClone handles object cloning: result = clone $obj
// OpClone - Clone object
func (vm *VM) opClone(frame *Frame, instr Instruction) error {
	// Get the object to clone
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if objVal.Type() != types.TypeObject {
		return fmt.Errorf("CLONE: not an object")
	}

	obj := objVal.ToObject()

	// Create a shallow copy of the object
	newObj := &types.Object{
		ClassName:   obj.ClassName,
		ClassEntry:  obj.ClassEntry,
		Properties:  make(map[string]*types.Property),
		ObjectID:    types.NextObjectID(),
		IsDestroyed: false,
	}

	// Copy properties
	for name, prop := range obj.Properties {
		// Create new property with copied value
		newProp := &types.Property{
			Value:      prop.Value.Copy(), // Shallow copy of the value
			Visibility: prop.Visibility,
			IsStatic:   prop.IsStatic,
			Type:       prop.Type,
			HasDefault: prop.HasDefault,
			Default:    prop.Default,
			IsReadOnly: prop.IsReadOnly,
			Hooks:      prop.Hooks,
		}
		newObj.Properties[name] = newProp
	}

	newObjVal := types.NewObject(newObj)

	// Check for __clone magic method
	if obj.ClassEntry != nil {
		if magicClone, hasMagic := obj.ClassEntry.MagicMethods["__clone"]; hasMagic {
			// TODO: Call __clone() on the new object
			// The __clone method is called on the copy, not the original
			_ = magicClone
		}
	}

	return vm.setOperandValue(frame, instr.Result, newObjVal)
}

// opInstanceof handles instanceof check: result = $obj instanceof Class
// OpInstanceof - Check if object is instance of class
func (vm *VM) opInstanceof(frame *Frame, instr Instruction) error {
	// Get the object/value
	val, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the class name
	className, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}
	classNameStr := className.ToString()

	var result bool

	// Only objects can be instanceof a class
	if val.Type() != types.TypeObject {
		result = false
	} else {
		obj := val.ToObject()

		if obj.ClassEntry == nil {
			result = false
		} else {
			// Check if object's class matches or is a subclass of the target class
			result = vm.isInstanceOf(obj.ClassEntry, classNameStr)
		}
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// isInstanceOf checks if a class is an instance of a target class (including inheritance and interfaces)
func (vm *VM) isInstanceOf(class *types.ClassEntry, targetClassName string) bool {
	if class == nil {
		return false
	}

	// Check direct match
	if class.Name == targetClassName {
		return true
	}

	// Check parent classes
	current := class.ParentClass
	for current != nil {
		if current.Name == targetClassName {
			return true
		}
		current = current.ParentClass
	}

	// Check implemented interfaces
	if class.ImplementsInterface(targetClassName) {
		return true
	}

	return false
}

// opGetClass handles get_class() function: result = get_class($obj)
// OpGetClass - Get class name
func (vm *VM) opGetClass(frame *Frame, instr Instruction) error {
	// Get the object
	objVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	var result *types.Value

	if objVal.Type() != types.TypeObject {
		// get_class() on non-object returns false
		result = types.NewBool(false)
	} else {
		obj := objVal.ToObject()
		if obj.ClassEntry != nil {
			result = types.NewString(obj.ClassEntry.Name)
		} else {
			result = types.NewString(obj.ClassName)
		}
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opFetchThis handles fetching $this variable
// OpFetchThis - Fetch $this variable
func (vm *VM) opFetchThis(frame *Frame, instr Instruction) error {
	// Get $this from frame context
	if frame.thisObject == nil {
		// $this is not available (static context or function context)
		return fmt.Errorf("FETCH_THIS: $this not available in this context")
	}

	return vm.setOperandValue(frame, instr.Result, types.NewObject(frame.thisObject))
}
