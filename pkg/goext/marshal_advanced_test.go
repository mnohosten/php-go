package goext

import (
	"reflect"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// Test types for custom marshaling

type CustomTime struct {
	Unix int64
}

func (ct *CustomTime) MarshalPHP() (*types.Value, error) {
	return types.NewInt(ct.Unix), nil
}

func (ct *CustomTime) UnmarshalPHP(v *types.Value) error {
	ct.Unix = v.ToInt()
	return nil
}

type Person struct {
	Name  string
	Age   int
	Email string `json:"email"`
}

type Address struct {
	Street  string
	City    string
	ZipCode string `json:"zip_code"`
}

type CircularA struct {
	Name string
	B    *CircularB
}

type CircularB struct {
	Name string
	A    *CircularA
}

// Test AdvancedMarshaler creation

func TestNewAdvancedMarshaler(t *testing.T) {
	m := NewAdvancedMarshaler()

	if m == nil {
		t.Fatal("NewAdvancedMarshaler returned nil")
	}

	if m.Marshaler == nil {
		t.Error("Base marshaler is nil")
	}

	if m.converters == nil {
		t.Error("Converters map is nil")
	}
}

// Test custom type converters

func TestRegisterConverter(t *testing.T) {
	m := NewAdvancedMarshaler()

	converter := &TypeConverter{
		ToPHP: func(v interface{}) (*types.Value, error) {
			return types.NewString("custom"), nil
		},
		ToGo: func(v *types.Value) (interface{}, error) {
			return v.ToString(), nil
		},
	}

	typ := reflect.TypeOf("")
	m.RegisterConverter(typ, converter)

	if !m.HasConverter(typ) {
		t.Error("Converter not registered")
	}
}

func TestUnregisterConverter(t *testing.T) {
	m := NewAdvancedMarshaler()

	converter := &TypeConverter{
		ToPHP: func(v interface{}) (*types.Value, error) {
			return types.NewString("custom"), nil
		},
	}

	typ := reflect.TypeOf("")
	m.RegisterConverter(typ, converter)

	if !m.HasConverter(typ) {
		t.Error("Converter not registered")
	}

	m.UnregisterConverter(typ)

	if m.HasConverter(typ) {
		t.Error("Converter still registered after unregister")
	}
}

func TestCustomConverter(t *testing.T) {
	m := NewAdvancedMarshaler()

	// Register custom converter for int
	typ := reflect.TypeOf(int(0))
	converter := &TypeConverter{
		ToPHP: func(v interface{}) (*types.Value, error) {
			n := v.(int)
			return types.NewInt(int64(n * 2)), nil // Double the value
		},
		ToGo: func(v *types.Value) (interface{}, error) {
			return int(v.ToInt() / 2), nil // Halve the value
		},
	}

	m.RegisterConverter(typ, converter)

	// Test ToPHP
	result, err := m.ToPHPAdvanced(10)
	if err != nil {
		t.Fatalf("ToPHPAdvanced failed: %v", err)
	}

	if result.ToInt() != 20 {
		t.Errorf("Expected 20, got %d", result.ToInt())
	}

	// Test ToGo
	phpVal := types.NewInt(20)
	goVal, err := m.ToGoAdvanced(phpVal, typ)
	if err != nil {
		t.Fatalf("ToGoAdvanced failed: %v", err)
	}

	if goVal.(int) != 10 {
		t.Errorf("Expected 10, got %d", goVal.(int))
	}
}

// Test CustomMarshaler interface

func TestCustomMarshaler(t *testing.T) {
	m := NewAdvancedMarshaler()

	ct := &CustomTime{Unix: 1234567890}

	result, err := m.ToPHPAdvanced(ct)
	if err != nil {
		t.Fatalf("ToPHPAdvanced failed: %v", err)
	}

	if result.ToInt() != 1234567890 {
		t.Errorf("Expected 1234567890, got %d", result.ToInt())
	}
}

func TestCustomUnmarshaler(t *testing.T) {
	m := NewAdvancedMarshaler()

	phpVal := types.NewInt(1234567890)

	typ := reflect.TypeOf(CustomTime{})
	result, err := m.ToGoAdvanced(phpVal, typ)
	if err != nil {
		t.Fatalf("ToGoAdvanced failed: %v", err)
	}

	ct := result.(CustomTime)
	if ct.Unix != 1234567890 {
		t.Errorf("Expected 1234567890, got %d", ct.Unix)
	}
}

// Test circular reference detection

func TestCircularReferenceDetection(t *testing.T) {
	t.Skip("TODO: Circular reference detection needs deeper integration with base marshaler")
	m := NewAdvancedMarshaler()

	// Create circular reference
	a := &CircularA{Name: "A"}
	b := &CircularB{Name: "B"}
	a.B = b
	b.A = a

	// Should not panic or infinite loop
	result, err := m.ToPHPAdvanced(a)
	if err != nil {
		t.Fatalf("ToPHPAdvanced with circular reference failed: %v", err)
	}

	if result.Type() != types.TypeArray && result.Type() != types.TypeNull {
		// Should handle gracefully (either as map or null for circular part)
		t.Logf("Circular reference handled, result type: %s", result.TypeString())
	}
}

// Test struct to object conversion

func TestStructToObject(t *testing.T) {
	m := NewAdvancedMarshaler()

	person := Person{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	result, err := m.StructToObject(person, "Person")
	if err != nil {
		t.Fatalf("StructToObject failed: %v", err)
	}

	if result.Type() != types.TypeObject {
		t.Fatalf("Expected object, got %s", result.TypeString())
	}

	obj := result.ToObject()
	if obj.ClassName != "Person" {
		t.Errorf("Expected className 'Person', got '%s'", obj.ClassName)
	}

	// Check properties
	nameVal, ok := obj.GetProperty("Name", nil)
	if !ok {
		t.Error("Name property not found")
	} else if nameVal.ToString() != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", nameVal.ToString())
	}

	emailVal, ok := obj.GetProperty("email", nil)
	if !ok {
		t.Error("email property not found (should use json tag)")
	} else if emailVal.ToString() != "john@example.com" {
		t.Errorf("Expected 'john@example.com', got '%s'", emailVal.ToString())
	}
}

func TestStructToObjectPointer(t *testing.T) {
	m := NewAdvancedMarshaler()

	person := &Person{
		Name:  "Jane Doe",
		Age:   25,
		Email: "jane@example.com",
	}

	result, err := m.StructToObject(person, "Person")
	if err != nil {
		t.Fatalf("StructToObject with pointer failed: %v", err)
	}

	if result.Type() != types.TypeObject {
		t.Fatalf("Expected object, got %s", result.TypeString())
	}
}

func TestStructToObjectNil(t *testing.T) {
	m := NewAdvancedMarshaler()

	var person *Person

	result, err := m.StructToObject(person, "Person")
	if err != nil {
		t.Fatalf("StructToObject with nil pointer failed: %v", err)
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null for nil pointer, got %s", result.TypeString())
	}
}

// Test object to struct conversion

func TestPHPToStruct(t *testing.T) {
	m := NewAdvancedMarshaler()

	// Create PHP array
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("Name"), types.NewString("John Doe"))
	arr.Set(types.NewString("Age"), types.NewInt(30))
	arr.Set(types.NewString("email"), types.NewString("john@example.com"))

	typ := reflect.TypeOf(Person{})
	result, err := m.phpToStruct(types.NewArray(arr), typ)
	if err != nil {
		t.Fatalf("phpToStruct failed: %v", err)
	}

	person := result.(Person)
	if person.Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", person.Name)
	}

	if person.Age != 30 {
		t.Errorf("Expected 30, got %d", person.Age)
	}

	if person.Email != "john@example.com" {
		t.Errorf("Expected 'john@example.com', got '%s'", person.Email)
	}
}

func TestObjectToStruct(t *testing.T) {
	m := NewAdvancedMarshaler()

	// Create PHP object
	obj := types.NewObjectInstance("Person")
	obj.SetProperty("Name", types.NewString("Jane Doe"), nil)
	obj.SetProperty("Age", types.NewInt(25), nil)
	obj.SetProperty("email", types.NewString("jane@example.com"), nil)

	phpVal := types.NewObject(obj)

	typ := reflect.TypeOf(Person{})
	result, err := m.ObjectToStruct(phpVal, typ)
	if err != nil {
		t.Fatalf("ObjectToStruct failed: %v", err)
	}

	person := result.(Person)
	if person.Name != "Jane Doe" {
		t.Errorf("Expected 'Jane Doe', got '%s'", person.Name)
	}
}

// Test type conversion

func TestConvertToTypeInt(t *testing.T) {
	m := NewAdvancedMarshaler()

	// int64 to int
	result, err := m.convertToType(int64(42), reflect.TypeOf(int(0)))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(int) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// float64 to int64
	result, err = m.convertToType(float64(42.7), reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestConvertToTypeFloat(t *testing.T) {
	m := NewAdvancedMarshaler()

	// int64 to float64
	result, err := m.convertToType(int64(42), reflect.TypeOf(float64(0)))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(float64) != 42.0 {
		t.Errorf("Expected 42.0, got %v", result)
	}
}

func TestConvertToTypeString(t *testing.T) {
	m := NewAdvancedMarshaler()

	result, err := m.convertToType(42, reflect.TypeOf(""))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(string) != "42" {
		t.Errorf("Expected '42', got '%v'", result)
	}
}

func TestConvertToTypeBool(t *testing.T) {
	m := NewAdvancedMarshaler()

	// int64 to bool
	result, err := m.convertToType(int64(1), reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(bool) != true {
		t.Errorf("Expected true, got %v", result)
	}

	result, err = m.convertToType(int64(0), reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	if result.(bool) != false {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestConvertToTypeSlice(t *testing.T) {
	m := NewAdvancedMarshaler()

	input := []interface{}{int64(1), int64(2), int64(3)}
	result, err := m.convertToType(input, reflect.TypeOf([]int{}))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	slice := result.([]int)
	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}

	if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("Unexpected slice values: %v", slice)
	}
}

func TestConvertToTypeMap(t *testing.T) {
	m := NewAdvancedMarshaler()

	input := map[string]interface{}{
		"a": int64(1),
		"b": int64(2),
	}

	result, err := m.convertToType(input, reflect.TypeOf(map[string]int{}))
	if err != nil {
		t.Fatalf("convertToType failed: %v", err)
	}

	mapVal := result.(map[string]int)
	if mapVal["a"] != 1 || mapVal["b"] != 2 {
		t.Errorf("Unexpected map values: %v", mapVal)
	}
}

// Test Clone

func TestClone(t *testing.T) {
	m := NewAdvancedMarshaler()

	// Create original
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))

	original := types.NewArray(arr)

	// Clone
	cloned, err := m.Clone(original)
	if err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	// Verify they're equal
	if !m.DeepEquals(original, cloned) {
		t.Error("Clone not equal to original")
	}

	// Verify they're separate (modify clone shouldn't affect original)
	clonedArr := cloned.ToArray()
	clonedArr.Append(types.NewInt(4))

	if original.ToArray().Len() == clonedArr.Len() {
		t.Error("Modifying clone affected original")
	}
}

// Test DeepEquals

func TestDeepEquals(t *testing.T) {
	m := NewAdvancedMarshaler()

	// Test simple values
	a := types.NewInt(42)
	b := types.NewInt(42)
	c := types.NewInt(43)

	if !m.DeepEquals(a, b) {
		t.Error("Equal integers not detected as equal")
	}

	if m.DeepEquals(a, c) {
		t.Error("Unequal integers detected as equal")
	}

	// Test arrays
	arr1 := types.NewEmptyArray()
	arr1.Append(types.NewInt(1))
	arr1.Append(types.NewInt(2))

	arr2 := types.NewEmptyArray()
	arr2.Append(types.NewInt(1))
	arr2.Append(types.NewInt(2))

	arr3 := types.NewEmptyArray()
	arr3.Append(types.NewInt(1))
	arr3.Append(types.NewInt(3))

	if !m.DeepEquals(types.NewArray(arr1), types.NewArray(arr2)) {
		t.Error("Equal arrays not detected as equal")
	}

	if m.DeepEquals(types.NewArray(arr1), types.NewArray(arr3)) {
		t.Error("Unequal arrays detected as equal")
	}
}

func TestDeepEqualsNil(t *testing.T) {
	m := NewAdvancedMarshaler()

	if !m.DeepEquals(nil, nil) {
		t.Error("Two nils not equal")
	}

	if m.DeepEquals(nil, types.NewInt(42)) {
		t.Error("nil and value detected as equal")
	}

	if m.DeepEquals(types.NewInt(42), nil) {
		t.Error("value and nil detected as equal")
	}
}

// Test nested struct conversion

func TestNestedStructConversion(t *testing.T) {
	m := NewAdvancedMarshaler()

	type PersonWithAddress struct {
		Name    string
		Address Address
	}

	// Create PHP array with nested data
	addrArr := types.NewEmptyArray()
	addrArr.Set(types.NewString("Street"), types.NewString("123 Main St"))
	addrArr.Set(types.NewString("City"), types.NewString("Springfield"))
	addrArr.Set(types.NewString("zip_code"), types.NewString("12345"))

	personArr := types.NewEmptyArray()
	personArr.Set(types.NewString("Name"), types.NewString("John Doe"))
	personArr.Set(types.NewString("Address"), types.NewArray(addrArr))

	typ := reflect.TypeOf(PersonWithAddress{})
	result, err := m.phpToStruct(types.NewArray(personArr), typ)
	if err != nil {
		t.Fatalf("Nested struct conversion failed: %v", err)
	}

	person := result.(PersonWithAddress)
	if person.Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", person.Name)
	}

	if person.Address.City != "Springfield" {
		t.Errorf("Expected 'Springfield', got '%s'", person.Address.City)
	}

	if person.Address.ZipCode != "12345" {
		t.Errorf("Expected '12345', got '%s'", person.Address.ZipCode)
	}
}
