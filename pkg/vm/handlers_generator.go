package vm

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Generator Opcode Handlers
// ============================================================================

// opGeneratorCreate handles OpGeneratorCreate
// Creates a generator object from a function
// Result operand receives the generator
func (vm *VM) opGeneratorCreate(frame *Frame, instr Instruction) error {
	// Get the function from operand
	funcVal, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Get the function
	var fn *CompiledFunction
	switch funcVal.Type() {
	case types.TypeString:
		// Function name
		fnName := funcVal.ToString()
		var ok bool
		fn, ok = vm.functions[fnName]
		if !ok {
			return fmt.Errorf("Function '%s' not found", fnName)
		}
	default:
		return fmt.Errorf("OpGeneratorCreate expects function name, got %s", funcVal.TypeString())
	}

	// Create generator
	generator := NewGenerator(vm, fn)

	// Wrap in resource
	// Store the generator pointer as the resource data
	resource := types.NewResourceHandle("Generator", generator)
	generatorVal := types.NewResource(resource)

	return vm.setOperandValue(frame, instr.Result, generatorVal)
}

// opYield handles OpYield
// Yields a value from the generator
// Op1: value to yield
// Op2: key to yield (optional, if Arg1 != 0)
// Arg1: flags (bit 0: has key)
func (vm *VM) opYield(frame *Frame, instr Instruction) error {
	// Check if current frame is executing a generator
	if !frame.IsGenerator() {
		return fmt.Errorf("yield can only be used in a generator function")
	}

	// Get the generator
	gen, ok := frame.GetGenerator().(*Generator)
	if !ok {
		return fmt.Errorf("Invalid generator object")
	}

	// Check if we have a key (flag in ExtendedValue)
	hasKey := (instr.ExtendedValue & 1) != 0

	var key, value *types.Value
	var err error

	if hasKey {
		// Read key and value
		key, err = vm.getOperandValue(frame, instr.Op2)
		if err != nil {
			return err
		}
		value, err = vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}

		// Yield with key
		gen.Yield(key, value, frame.GetIP())
	} else {
		// Read value only
		value, err = vm.getOperandValue(frame, instr.Op1)
		if err != nil {
			return err
		}

		// Yield with auto-increment key
		gen.Yield(nil, value, frame.GetIP())
	}

	// Suspend execution - return control to the generator caller
	// The VM will need to handle this specially
	// For now, we'll pop the frame to exit the generator execution
	// TODO: Implement proper generator suspension/resumption
	return fmt.Errorf("Generator suspension not yet fully implemented")
}

// opGeneratorReturn handles OpGeneratorReturn
// Returns from a generator with a value (accessible via getReturn())
// Op1: return value
func (vm *VM) opGeneratorReturn(frame *Frame, instr Instruction) error {
	// Check if current frame is executing a generator
	if !frame.IsGenerator() {
		return fmt.Errorf("generator return can only be used in a generator function")
	}

	// Get the generator
	gen, ok := frame.GetGenerator().(*Generator)
	if !ok {
		return fmt.Errorf("Invalid generator object")
	}

	// Get return value
	returnValue, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Set generator return value
	gen.Return(returnValue)

	// Exit the generator
	vm.popFrame()
	return nil
}

// opYieldFrom handles OpYieldFrom
// Yields all values from another generator or iterable
// Op1: generator/iterable to yield from
func (vm *VM) opYieldFrom(frame *Frame, instr Instruction) error {
	// Check if current frame is executing a generator
	if !frame.IsGenerator() {
		return fmt.Errorf("yield from can only be used in a generator function")
	}

	// Get the iterable
	iterable, err := vm.getOperandValue(frame, instr.Op1)
	if err != nil {
		return err
	}

	// Check if it's a generator
	if iterable.Type() == types.TypeResource {
		resource := iterable.ToResource()
		// Get the data from the resource
		data := resource.Data()
		if subGen, ok := data.(*Generator); ok {
			// Yield from this generator
			// This is complex - we need to:
			// 1. Execute the sub-generator until it's done
			// 2. Yield each of its values
			// 3. Return its final return value

			// For now, this is a simplified implementation
			// TODO: Implement full yield from semantics

			// Get current generator
			gen, ok := frame.GetGenerator().(*Generator)
			if !ok {
				return fmt.Errorf("Invalid generator object")
			}

			// Execute sub-generator to completion
			// Yielding each value from the current generator
			for {
				err := subGen.Next()
				if err != nil {
					return err
				}

				if !subGen.Valid() {
					// Sub-generator is done
					break
				}

				// Yield the value from sub-generator
				gen.Yield(subGen.Key(), subGen.Current(), frame.GetIP())

				// Suspend current generator
				// TODO: Proper suspension
			}

			// Get return value from sub-generator
			returnVal, err := subGen.GetReturn()
			if err == nil && returnVal != nil {
				// Push return value onto stack for use
				// TODO: Determine correct behavior
			}

			return nil
		}
	}

	return fmt.Errorf("yield from expects generator or iterable, got %s", iterable.TypeString())
}
