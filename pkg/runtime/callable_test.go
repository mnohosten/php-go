package runtime

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestNewFunctionCallable tests creating function callable
func TestNewFunctionCallable(t *testing.T) {
	callable := NewFunctionCallable("strlen")

	if callable == nil {
		t.Fatal("NewFunctionCallable() returned nil")
	}

	if callable.Type != CallableTypeFunction {
		t.Errorf("Type = %v, want CallableTypeFunction", callable.Type)
	}

	if callable.FunctionName != "strlen" {
		t.Errorf("FunctionName = %s, want strlen", callable.FunctionName)
	}
}

// TestNewMethodCallable tests creating method callable
func TestNewMethodCallable(t *testing.T) {
	obj := &types.Object{
		ClassName: "MyClass",
	}

	callable := NewMethodCallable(obj, "myMethod")

	if callable == nil {
		t.Fatal("NewMethodCallable() returned nil")
	}

	if callable.Type != CallableTypeMethod {
		t.Errorf("Type = %v, want CallableTypeMethod", callable.Type)
	}

	if callable.Object != obj {
		t.Error("Object not set correctly")
	}

	if callable.MethodName != "myMethod" {
		t.Errorf("MethodName = %s, want myMethod", callable.MethodName)
	}
}

// TestNewStaticMethodCallable tests creating static method callable
func TestNewStaticMethodCallable(t *testing.T) {
	callable := NewStaticMethodCallable("MyClass", "staticMethod")

	if callable == nil {
		t.Fatal("NewStaticMethodCallable() returned nil")
	}

	if callable.Type != CallableTypeStaticMethod {
		t.Errorf("Type = %v, want CallableTypeStaticMethod", callable.Type)
	}

	if callable.ClassName != "MyClass" {
		t.Errorf("ClassName = %s, want MyClass", callable.ClassName)
	}

	if callable.MethodName != "staticMethod" {
		t.Errorf("MethodName = %s, want staticMethod", callable.MethodName)
	}
}

// TestNewClosureCallable tests creating closure callable
func TestNewClosureCallable(t *testing.T) {
	closure := types.NewNull() // Placeholder

	callable := NewClosureCallable(closure)

	if callable == nil {
		t.Fatal("NewClosureCallable() returned nil")
	}

	if callable.Type != CallableTypeClosure {
		t.Errorf("Type = %v, want CallableTypeClosure", callable.Type)
	}

	if callable.Closure != closure {
		t.Error("Closure not set correctly")
	}
}

// TestCallable_GetName tests getting callable name
func TestCallable_GetName(t *testing.T) {
	tests := []struct {
		name     string
		callable *Callable
		want     string
	}{
		{
			name:     "function",
			callable: NewFunctionCallable("strlen"),
			want:     "strlen",
		},
		{
			name: "method",
			callable: NewMethodCallable(&types.Object{
				ClassName: "MyClass",
			}, "myMethod"),
			want: "MyClass->myMethod",
		},
		{
			name:     "static method",
			callable: NewStaticMethodCallable("MyClass", "staticMethod"),
			want:     "MyClass::staticMethod",
		},
		{
			name:     "closure",
			callable: NewClosureCallable(types.NewNull()),
			want:     "{closure}",
		},
		{
			name:     "nil",
			callable: nil,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.callable.GetName()
			if got != tt.want {
				t.Errorf("GetName() = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestCallable_GetName_MethodWithNilObject tests method callable with nil object
func TestCallable_GetName_MethodWithNilObject(t *testing.T) {
	callable := &Callable{
		Type:       CallableTypeMethod,
		MethodName: "test",
	}

	name := callable.GetName()
	if name != "test" {
		t.Errorf("GetName() = %s, want test", name)
	}
}

// TestCallable_ToValue tests converting to PHP value
func TestCallable_ToValue(t *testing.T) {
	callable := NewFunctionCallable("strlen")
	value := callable.ToValue()

	if value.Type() != types.TypeResource {
		t.Errorf("ToValue() type = %v, want TypeResource", value.Type())
	}
}

// TestCallable_ToValue_Nil tests nil callable
func TestCallable_ToValue_Nil(t *testing.T) {
	var callable *Callable
	value := callable.ToValue()

	if value.Type() != types.TypeNull {
		t.Errorf("ToValue() on nil type = %v, want TypeNull", value.Type())
	}
}

// TestCallable_ToValue_Closure tests closure callable
func TestCallable_ToValue_Closure(t *testing.T) {
	closureValue := types.NewInt(42) // Placeholder
	callable := NewClosureCallable(closureValue)

	value := callable.ToValue()

	// For closures, should return the closure value directly
	if value != closureValue {
		t.Error("ToValue() for closure should return closure value directly")
	}
}

// TestIsCallable_String tests string callable
func TestIsCallable_String(t *testing.T) {
	value := types.NewString("strlen")

	if !IsCallable(value) {
		t.Error("IsCallable(string) = false, want true")
	}
}

// TestIsCallable_Array tests array callable
func TestIsCallable_Array(t *testing.T) {
	arr := types.NewArrayWithCapacity(2)
	arr.Set(types.NewInt(0), types.NewString("MyClass"))
	arr.Set(types.NewInt(1), types.NewString("method"))

	value := types.NewArray(arr)

	if !IsCallable(value) {
		t.Error("IsCallable(array) = false, want true for 2-element array")
	}
}

// TestIsCallable_Array_WrongSize tests array with wrong size
func TestIsCallable_Array_WrongSize(t *testing.T) {
	arr := types.NewArrayWithCapacity(1)
	arr.Set(types.NewInt(0), types.NewString("MyClass"))

	value := types.NewArray(arr)

	if IsCallable(value) {
		t.Error("IsCallable(1-element array) = true, want false")
	}
}

// TestIsCallable_Object tests object with __invoke
func TestIsCallable_Object(t *testing.T) {
	obj := &types.Object{
		ClassName: "InvokableClass",
		ClassEntry: &types.ClassEntry{
			Name: "InvokableClass",
			Methods: map[string]*types.MethodDef{
				"__invoke": {
					Name: "__invoke",
				},
			},
		},
	}

	value := types.NewObject(obj)

	if !IsCallable(value) {
		t.Error("IsCallable(object with __invoke) = false, want true")
	}
}

// TestIsCallable_Object_NoInvoke tests object without __invoke
func TestIsCallable_Object_NoInvoke(t *testing.T) {
	obj := &types.Object{
		ClassName: "RegularClass",
		ClassEntry: &types.ClassEntry{
			Name:    "RegularClass",
			Methods: map[string]*types.MethodDef{},
		},
	}

	value := types.NewObject(obj)

	if IsCallable(value) {
		t.Error("IsCallable(object without __invoke) = true, want false")
	}
}

// TestIsCallable_Resource tests callable resource
func TestIsCallable_Resource(t *testing.T) {
	callable := NewFunctionCallable("strlen")
	value := callable.ToValue()

	if !IsCallable(value) {
		t.Error("IsCallable(Callable resource) = false, want true")
	}
}

// TestIsCallable_Nil tests nil value
func TestIsCallable_Nil(t *testing.T) {
	if IsCallable(nil) {
		t.Error("IsCallable(nil) = true, want false")
	}
}

// TestIsCallable_Int tests non-callable value
func TestIsCallable_Int(t *testing.T) {
	value := types.NewInt(42)

	if IsCallable(value) {
		t.Error("IsCallable(int) = true, want false")
	}
}

// TestCallableFromValue tests extracting callable from value
func TestCallableFromValue(t *testing.T) {
	original := NewFunctionCallable("strlen")
	value := original.ToValue()

	extracted, err := CallableFromValue(value)
	if err != nil {
		t.Fatalf("CallableFromValue() error = %v", err)
	}

	if extracted.FunctionName != "strlen" {
		t.Errorf("Extracted callable function name = %s, want strlen", extracted.FunctionName)
	}
}

// TestCallableFromValue_Error_NotResource tests non-resource value
func TestCallableFromValue_Error_NotResource(t *testing.T) {
	value := types.NewInt(42)

	_, err := CallableFromValue(value)
	if err == nil {
		t.Error("CallableFromValue(int) should return error")
	}
}

// TestCallableFromValue_Error_InvalidResource tests invalid resource
func TestCallableFromValue_Error_InvalidResource(t *testing.T) {
	// Create a different type of resource
	resource := types.NewResourceHandle("NotCallable", "some data")
	value := types.NewResource(resource)

	_, err := CallableFromValue(value)
	if err == nil {
		t.Error("CallableFromValue(wrong resource type) should return error")
	}
}

// TestValidateCallable_Function tests function callable validation
func TestValidateCallable_Function(t *testing.T) {
	callable := NewFunctionCallable("strlen")

	err := ValidateCallable(callable)
	if err != nil {
		t.Errorf("ValidateCallable(function) error = %v", err)
	}
}

// TestValidateCallable_Function_Error_EmptyName tests empty function name
func TestValidateCallable_Function_Error_EmptyName(t *testing.T) {
	callable := NewFunctionCallable("")

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(empty function name) should return error")
	}
}

// TestValidateCallable_Method tests method callable validation
func TestValidateCallable_Method(t *testing.T) {
	obj := &types.Object{ClassName: "Test"}
	callable := NewMethodCallable(obj, "method")

	err := ValidateCallable(callable)
	if err != nil {
		t.Errorf("ValidateCallable(method) error = %v", err)
	}
}

// TestValidateCallable_Method_Error_NilObject tests nil object
func TestValidateCallable_Method_Error_NilObject(t *testing.T) {
	callable := NewMethodCallable(nil, "method")

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(nil object) should return error")
	}
}

// TestValidateCallable_Method_Error_EmptyName tests empty method name
func TestValidateCallable_Method_Error_EmptyName(t *testing.T) {
	obj := &types.Object{ClassName: "Test"}
	callable := NewMethodCallable(obj, "")

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(empty method name) should return error")
	}
}

// TestValidateCallable_StaticMethod tests static method callable validation
func TestValidateCallable_StaticMethod(t *testing.T) {
	callable := NewStaticMethodCallable("MyClass", "staticMethod")

	err := ValidateCallable(callable)
	if err != nil {
		t.Errorf("ValidateCallable(static method) error = %v", err)
	}
}

// TestValidateCallable_StaticMethod_Error_EmptyClass tests empty class name
func TestValidateCallable_StaticMethod_Error_EmptyClass(t *testing.T) {
	callable := NewStaticMethodCallable("", "method")

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(empty class name) should return error")
	}
}

// TestValidateCallable_StaticMethod_Error_EmptyMethod tests empty method name
func TestValidateCallable_StaticMethod_Error_EmptyMethod(t *testing.T) {
	callable := NewStaticMethodCallable("MyClass", "")

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(empty method name) should return error")
	}
}

// TestValidateCallable_Closure tests closure callable validation
func TestValidateCallable_Closure(t *testing.T) {
	closure := types.NewNull()
	callable := NewClosureCallable(closure)

	err := ValidateCallable(callable)
	if err != nil {
		t.Errorf("ValidateCallable(closure) error = %v", err)
	}
}

// TestValidateCallable_Closure_Error_Nil tests nil closure
func TestValidateCallable_Closure_Error_Nil(t *testing.T) {
	callable := NewClosureCallable(nil)

	err := ValidateCallable(callable)
	if err == nil {
		t.Error("ValidateCallable(nil closure) should return error")
	}
}

// TestValidateCallable_Error_Nil tests nil callable
func TestValidateCallable_Error_Nil(t *testing.T) {
	err := ValidateCallable(nil)
	if err == nil {
		t.Error("ValidateCallable(nil) should return error")
	}
}

// TestInvoke_NotImplemented tests that Invoke returns error (not yet implemented)
func TestInvoke_NotImplemented(t *testing.T) {
	callable := NewFunctionCallable("strlen")
	args := []*types.Value{types.NewString("test")}

	_, err := callable.Invoke(args)
	if err == nil {
		t.Error("Invoke() should return error (not yet implemented)")
	}
}

// TestInvoke_Nil tests invoking nil callable
func TestInvoke_Nil(t *testing.T) {
	var callable *Callable
	args := []*types.Value{}

	_, err := callable.Invoke(args)
	if err == nil {
		t.Error("Invoke() on nil should return error")
	}
}
