package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Array Opcode Handlers
// ============================================================================

// opInitArray handles array initialization: result = []
// OpInitArray - Initialize array
func (vm *VM) opInitArray(frame *Frame, instr Instruction) error {
	// Create a new empty array
	arr := types.NewEmptyArray()
	result := types.NewArray(arr)

	return vm.setOperandValue(frame, instr.Result, result)
}

// opAddArrayElement handles adding an element to an array during initialization
// OpAddArrayElement - Add element to array: result[] = op1 or result[op2] = op1
func (vm *VM) opAddArrayElement(frame *Frame, instr Instruction) error {
	// Get the array (should be in result)
	arrayVal, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	if arrayVal.Type() != types.TypeArray {
		return fmt.Errorf("ADD_ARRAY_ELEMENT: result is not an array")
	}

	arr := arrayVal.ToArray()

	// Get the value to add
	value, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Check if we have a key (op2)
	if instr.Op2.Type != OpUnused {
		// Add with specific key: result[op2] = op1
		key, err := vm.getOperandValue(frame, instr.Op2)
		if err != nil {
			return err
		}
		arr.Set(key, value)
	} else {
		// Append: result[] = op1
		arr.Append(value)
	}

	return nil
}

// ============================================================================
// Fetch Operations (Read)
// ============================================================================

// opFetchDimR handles fetching array element for read: result = $arr[$key]
// OpFetchDimR - Fetch array element for read
func (vm *VM) opFetchDimR(frame *Frame, instr Instruction) error {
	// Get the array/string
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the key/index
	key, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result *types.Value

	switch container.Type() {
	case types.TypeArray:
		// Array access
		arr := container.ToArray()
		val, exists := arr.Get(key)
		if !exists {
			// PHP returns NULL for undefined array keys (with notice)
			result = types.NewNull()
		} else {
			result = val
		}

	case types.TypeString:
		// String offset access: $str[$index]
		str := container.ToString()
		index := int(key.ToInt())

		if index < 0 || index >= len(str) {
			// Out of bounds
			result = types.NewString("")
		} else {
			result = types.NewString(string(str[index]))
		}

	default:
		// Non-array, non-string - PHP returns NULL (with warning)
		result = types.NewNull()
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opFetchDimW handles fetching array element for write: $arr[$key] = ...
// OpFetchDimW - Fetch array element for write
func (vm *VM) opFetchDimW(frame *Frame, instr Instruction) error {
	// Get the array
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// For write operations, we need to ensure we have an array
	if container.Type() != types.TypeArray {
		// PHP auto-vivifies to an array
		newArr := types.NewEmptyArray()
		container = types.NewArray(newArr)
		// Update the original variable
		vm.setOperandValue(frame, instr.Op1, container)
	}

	arr := container.ToArray()

	// Get the key (might be unspecified for append operation)
	var key *types.Value
	if instr.Op2.Type != OpUnused {
		key, err = vm.getOperandValue(frame, instr.Op2)
		if err != nil {
			return err
		}
	} else {
		// No key specified - this is for append operation $arr[] = ...
		// We'll use the next auto-increment index
		key = nil
	}

	// For write fetch, we return a reference to the array element
	// In a simplified VM, we'll create/get the element
	if key != nil {
		// Check if element exists
		val, exists := arr.Get(key)
		if !exists {
			// Create new element
			val = types.NewNull()
			arr.Set(key, val)
		}
		return vm.setOperandValue(frame, instr.Result, val)
	} else {
		// Append operation - create a placeholder
		// The actual append will happen in OpAssignDim
		return vm.setOperandValue(frame, instr.Result, types.NewNull())
	}
}

// opFetchDimRW handles fetching array element for read-write: $arr[$key] += 1
// OpFetchDimRW - Fetch array element for read-write
func (vm *VM) opFetchDimRW(frame *Frame, instr Instruction) error {
	// Similar to FetchDimW, but we need the current value
	return vm.opFetchDimW(frame, instr)
}

// opFetchDimIs handles fetching array element for isset/empty check
// OpFetchDimIs - Fetch array element for isset/empty check
func (vm *VM) opFetchDimIs(frame *Frame, instr Instruction) error {
	// Similar to FetchDimR but doesn't generate notices
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	key, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result *types.Value

	if container.Type() == types.TypeArray {
		arr := container.ToArray()
		val, exists := arr.Get(key)
		if !exists {
			result = types.NewNull()
		} else {
			result = val
		}
	} else {
		result = types.NewNull()
	}

	return vm.setOperandValue(frame, instr.Result, result)
}

// opFetchDimFuncArg handles fetching array element as function argument
// OpFetchDimFuncArg - Fetch array element as function argument
func (vm *VM) opFetchDimFuncArg(frame *Frame, instr Instruction) error {
	// For function arguments, use read semantics
	return vm.opFetchDimR(frame, instr)
}

// opFetchDimUnset handles fetching array element for unset
// OpFetchDimUnset - Fetch array element for unset
func (vm *VM) opFetchDimUnset(frame *Frame, instr Instruction) error {
	// For unset, we just need to identify the element
	return vm.opFetchDimR(frame, instr)
}

// ============================================================================
// Assignment Operations
// ============================================================================

// opAssignDim handles array element assignment: op1[op2] = result
// OpAssignDim - Array element assignment
func (vm *VM) opAssignDim(frame *Frame, instr Instruction) error {
	// Get the array
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Auto-vivify to array if needed
	if container.Type() != types.TypeArray {
		newArr := types.NewEmptyArray()
		container = types.NewArray(newArr)
		vm.setOperandValue(frame, instr.Op1, container)
	}

	arr := container.ToArray()

	// Get the value to assign
	value, err := vm.getOperandValue(frame, instr.Result)
	if err != nil {
		return err
	}

	// Get the key (might be unspecified for append)
	if instr.Op2.Type != OpUnused {
		// Specific key: $arr[$key] = $value
		key, err := vm.getOperandValue(frame, instr.Op2)
		if err != nil {
			return err
		}
		arr.Set(key, value)
	} else {
		// Append: $arr[] = $value
		arr.Append(value)
	}

	return nil
}

// opAssignDimOp handles compound array assignment: op1[op2] += result
// OpAssignDimOp - Compound array assignment
func (vm *VM) opAssignDimOp(frame *Frame, instr Instruction) error {
	// Get the array
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if container.Type() != types.TypeArray {
		return fmt.Errorf("ASSIGN_DIM_OP: container is not an array")
	}

	arr := container.ToArray()

	// Get the key
	key, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Get the current value
	currentVal, exists := arr.Get(key)
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

	arr.Set(key, newVal)
	return nil
}

// ============================================================================
// Unset Operations
// ============================================================================

// opUnsetDim handles unsetting array element: unset($arr[$key])
// OpUnsetDim - Unset array element
func (vm *VM) opUnsetDim(frame *Frame, instr Instruction) error {
	// Get the array
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if container.Type() != types.TypeArray {
		// Unset on non-array is a no-op in PHP
		return nil
	}

	arr := container.ToArray()

	// Get the key
	key, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	// Remove the element
	arr.Unset(key)
	return nil
}

// ============================================================================
// Isset/Empty Operations
// ============================================================================

// opIssetIsemptyDimObj handles isset/empty check on array element
// OpIssetIsemptyDimObj - Check isset/empty on array element or object property
func (vm *VM) opIssetIsemptyDimObj(frame *Frame, instr Instruction) error {
	// Get the container (array or object)
	container, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the key
	key, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var result bool

	switch container.Type() {
	case types.TypeArray:
		arr := container.ToArray()
		val, exists := arr.Get(key)

		// For isset: check if exists and not null
		// For empty: check if exists and is "empty" (falsy)
		// TODO: Determine from instruction if this is isset or empty check
		// For now, implement isset semantics
		result = exists && !val.IsNull()

	case types.TypeObject:
		// TODO: Implement object property isset check in Phase 5
		result = false

	default:
		result = false
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(result))
}

// ============================================================================
// Array Built-in Functions
// ============================================================================

// opCount handles count() function: result = count($arr)
// OpCount - Count array elements
func (vm *VM) opCount(frame *Frame, instr Instruction) error {
	// Get the array
	arrayVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	var count int64

	switch arrayVal.Type() {
	case types.TypeArray:
		arr := arrayVal.ToArray()
		count = int64(arr.Len())

	case types.TypeObject:
		// TODO: For objects, call Countable interface or count properties
		count = 1

	case types.TypeNull:
		count = 0

	default:
		// For scalars, count is 1
		count = 1
	}

	return vm.setOperandValue(frame, instr.Result, types.NewInt(count))
}

// opInArray handles in_array() function: result = in_array($needle, $haystack)
// OpInArray - Check if value is in array
func (vm *VM) opInArray(frame *Frame, instr Instruction) error {
	// Get the needle (value to search for)
	needle, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the haystack (array to search in)
	haystackVal, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var found bool

	if haystackVal.Type() == types.TypeArray {
		haystack := haystackVal.ToArray()
		found = haystack.Contains(needle)
	} else {
		found = false
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(found))
}

// opArrayKeyExists handles array_key_exists() function
// OpArrayKeyExists - Check if array key exists
func (vm *VM) opArrayKeyExists(frame *Frame, instr Instruction) error {
	// Get the key
	key, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the array
	arrayVal, err := vm.getOperandValue(frame, instr.Op2)
	if err != nil {
		return err
	}

	var exists bool

	if arrayVal.Type() == types.TypeArray {
		arr := arrayVal.ToArray()
		exists = arr.HasKey(key)
	} else {
		exists = false
	}

	return vm.setOperandValue(frame, instr.Result, types.NewBool(exists))
}
