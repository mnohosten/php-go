package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Foreach Iterator Support
// ============================================================================

// ArrayIterator represents an iterator for foreach loops
type ArrayIterator struct {
	array     *types.Array
	position  int
	keys      []interface{} // Cached keys for iteration order
	byRef     bool          // Whether values should be by reference
	exhausted bool          // Whether iteration is complete
}

// NewArrayIterator creates a new array iterator
func NewArrayIterator(arr *types.Array, byRef bool) *ArrayIterator {
	if arr == nil {
		return &ArrayIterator{
			array:     nil,
			position:  0,
			keys:      nil,
			byRef:     byRef,
			exhausted: true,
		}
	}

	// Get all keys in iteration order
	keys := arr.GetKeysSlice()

	return &ArrayIterator{
		array:     arr,
		position:  0,
		keys:      keys,
		byRef:     byRef,
		exhausted: len(keys) == 0,
	}
}

// Valid returns whether the iterator is at a valid position
func (it *ArrayIterator) Valid() bool {
	if it == nil || it.array == nil {
		return false
	}
	return it.position >= 0 && it.position < len(it.keys) && !it.exhausted
}

// Current returns the current value
func (it *ArrayIterator) Current() (*types.Value, error) {
	if !it.Valid() {
		return nil, fmt.Errorf("invalid iterator position")
	}
	key := it.keys[it.position]

	// Convert key to Value for array lookup
	var keyVal *types.Value
	switch k := key.(type) {
	case int64:
		keyVal = types.NewInt(k)
	case string:
		keyVal = types.NewString(k)
	default:
		keyVal = types.NewString(fmt.Sprintf("%v", k))
	}

	value, exists := it.array.Get(keyVal)
	if !exists {
		return types.NewNull(), nil
	}
	return value, nil
}

// Key returns the current key
func (it *ArrayIterator) Key() (*types.Value, error) {
	if !it.Valid() {
		return nil, fmt.Errorf("invalid iterator position")
	}
	key := it.keys[it.position]

	// Convert key to Value
	switch k := key.(type) {
	case int64:
		return types.NewInt(k), nil
	case string:
		return types.NewString(k), nil
	default:
		return types.NewString(fmt.Sprintf("%v", k)), nil
	}
}

// Next advances the iterator to the next position
func (it *ArrayIterator) Next() {
	if it.Valid() {
		it.position++
		if it.position >= len(it.keys) {
			it.exhausted = true
		}
	}
}

// ============================================================================
// Foreach Opcode Handlers
// ============================================================================

// opFeResetR handles FE_RESET_R opcode - Initialize foreach iterator for read
// Op1: Array value (TmpVar 0)
// Op2: Unused
// Result: Iterator stored in result operand (TmpVar 1)
func (vm *VM) opFeResetR(frame *Frame, instr Instruction) error {
	// Get the array value
	arrayVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Convert to array
	arr := arrayVal.ToArray()

	// Create iterator
	iterator := NewArrayIterator(arr, false)

	// Store iterator as a resource in the result temporary variable
	resource := types.NewResourceHandle("foreach_iterator", iterator)
	return vm.setOperandValue(frame, instr.Result, types.NewResource(resource))
}

// opFeResetRW handles FE_RESET_RW opcode - Initialize foreach iterator for read-write
// Op1: Array value (TmpVar 0)
// Op2: Unused
// Result: Iterator stored in result operand (TmpVar 1)
func (vm *VM) opFeResetRW(frame *Frame, instr Instruction) error {
	// Get the array value
	arrayVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Convert to array
	arr := arrayVal.ToArray()

	// Create iterator with by-ref flag
	iterator := NewArrayIterator(arr, true)

	// Store iterator as a resource in the result temporary variable
	resource := types.NewResourceHandle("foreach_iterator", iterator)
	return vm.setOperandValue(frame, instr.Result, types.NewResource(resource))
}

// opFeFetchR handles FE_FETCH_R opcode - Fetch next foreach element for read
// Op1: Iterator (TmpVar 1)
// Op2: Jump target if done (instruction index)
// Result: Fetched value (TmpVar 2)
// Side effect: Key stored in TmpVar 3
func (vm *VM) opFeFetchR(frame *Frame, instr Instruction) error {
	// Get the iterator
	iteratorVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if iteratorVal.Type() != types.TypeResource {
		return fmt.Errorf("expected iterator resource, got %v", iteratorVal.Type())
	}

	resource := iteratorVal.ToResource()
	if resource == nil {
		return fmt.Errorf("invalid resource")
	}

	iterator, ok := resource.Data().(*ArrayIterator)
	if !ok {
		return fmt.Errorf("expected ArrayIterator, got %T", resource.Data())
	}

	// Check if iterator is valid
	if !iterator.Valid() {
		// Iterator exhausted, jump to end of loop
		var target int
		if instr.Op2.Type == OpConst {
			targetVal, err := vm.getOperandValue(frame, instr.Op2)
			if err == nil {
				target = int(targetVal.ToInt())
			} else {
				target = int(instr.Op2.Value)
			}
		} else {
			target = int(instr.Op2.Value)
		}
		frame.ip = target
		return nil
	}

	// Fetch current key and value
	key, err := iterator.Key()
	if err != nil {
		return err
	}

	value, err := iterator.Current()
	if err != nil {
		return err
	}

	// Store value in result
	if err := vm.setOperandValue(frame, instr.Result, value); err != nil {
		return err
	}

	// Store key in Result+1 (next tmpvar after result)
	// This assumes Result is a TmpVar operand
	if instr.Result.Type == OpTmpVar {
		keyOp := TmpVarOperand(instr.Result.Value + 1)
		if err := vm.setOperandValue(frame, keyOp, key); err != nil {
			return err
		}
	}

	// Advance iterator
	iterator.Next()

	return nil
}

// opFeFetchRW handles FE_FETCH_RW opcode - Fetch next foreach element for read-write
// Op1: Iterator (TmpVar 1)
// Op2: Jump target if done (instruction index)
// Result: Fetched value (TmpVar 2)
// Side effect: Key stored in TmpVar 3
func (vm *VM) opFeFetchRW(frame *Frame, instr Instruction) error {
	// Get the iterator
	iteratorVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	if iteratorVal.Type() != types.TypeResource {
		return fmt.Errorf("expected iterator resource, got %v", iteratorVal.Type())
	}

	resource := iteratorVal.ToResource()
	if resource == nil {
		return fmt.Errorf("invalid resource")
	}

	iterator, ok := resource.Data().(*ArrayIterator)
	if !ok {
		return fmt.Errorf("expected ArrayIterator, got %T", resource.Data())
	}

	// Check if iterator is valid
	if !iterator.Valid() {
		// Iterator exhausted, jump to end of loop
		var target int
		if instr.Op2.Type == OpConst {
			targetVal, err := vm.getOperandValue(frame, instr.Op2)
			if err == nil {
				target = int(targetVal.ToInt())
			} else {
				target = int(instr.Op2.Value)
			}
		} else {
			target = int(instr.Op2.Value)
		}
		frame.ip = target
		return nil
	}

	// Fetch current key and value
	key, err := iterator.Key()
	if err != nil {
		return err
	}

	value, err := iterator.Current()
	if err != nil {
		return err
	}

	// For by-ref iteration, we should store a reference
	// For now, we'll just store the value directly
	// TODO: Implement proper reference semantics

	// Store value in result
	if err := vm.setOperandValue(frame, instr.Result, value); err != nil {
		return err
	}

	// Store key in Result+1 (next tmpvar after result)
	// This assumes Result is a TmpVar operand
	if instr.Result.Type == OpTmpVar {
		keyOp := TmpVarOperand(instr.Result.Value + 1)
		if err := vm.setOperandValue(frame, keyOp, key); err != nil {
			return err
		}
	}

	// Advance iterator
	iterator.Next()

	return nil
}

// opFeFree handles FE_FREE opcode - Free foreach iterator
// Op1: Iterator (TmpVar 1)
// Op2: Unused
// Result: Unused
func (vm *VM) opFeFree(frame *Frame, instr Instruction) error {
	// Get the iterator
	iteratorVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		// If we can't get the iterator, it might already be freed
		// Just ignore the error
		return nil
	}

	if iteratorVal.Type() != types.TypeResource {
		// Not a resource, nothing to free
		return nil
	}

	// Set the iterator to null to allow garbage collection
	return vm.setOperandValue(frame, instr.Op1, types.NewNull())
}
