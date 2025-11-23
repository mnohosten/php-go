package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionParameter_BasicInfo tests basic parameter information
func TestReflectionParameter_BasicInfo(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		NumParams:  1,
		Parameters: []*types.ParameterDef{
			{
				Name: "testParam",
				Type: "string",
			},
		},
	}

	rm := NewReflectionMethod(rc, "testMethod", methodDef)
	paramDef := methodDef.Parameters[0]

	rp := NewReflectionParameter(rm, 0, paramDef)

	// Test GetName
	if name := rp.GetName(); name != "testParam" {
		t.Errorf("GetName() = %s, want testParam", name)
	}

	// Test GetPosition
	if pos := rp.GetPosition(); pos != 0 {
		t.Errorf("GetPosition() = %d, want 0", pos)
	}

	// Test GetDeclaringFunction
	if declaringFunc := rp.GetDeclaringFunction(); declaringFunc != rm {
		t.Error("GetDeclaringFunction() returned wrong method")
	}
}

// TestReflectionParameter_Type tests type information
func TestReflectionParameter_Type(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	// Parameter with type
	paramDef1 := &types.ParameterDef{
		Name: "typedParam",
		Type: "int",
	}
	rp1 := NewReflectionParameter(rm, 0, paramDef1)

	if !rp1.HasType() {
		t.Error("HasType() = false for typed parameter, want true")
	}
	if typ := rp1.GetType(); typ != "int" {
		t.Errorf("GetType() = %s, want int", typ)
	}

	// Parameter without type
	paramDef2 := &types.ParameterDef{
		Name: "untypedParam",
		Type: "",
	}
	rp2 := NewReflectionParameter(rm, 0, paramDef2)

	if rp2.HasType() {
		t.Error("HasType() = true for untyped parameter, want false")
	}
}

// TestReflectionParameter_Optional tests optional parameter checks
func TestReflectionParameter_Optional(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	// Required parameter
	paramDef1 := &types.ParameterDef{
		Name:       "requiredParam",
		Type:       "string",
		HasDefault: false,
	}
	rp1 := NewReflectionParameter(rm, 0, paramDef1)

	if rp1.IsOptional() {
		t.Error("IsOptional() = true for required parameter, want false")
	}

	// Optional parameter
	paramDef2 := &types.ParameterDef{
		Name:       "optionalParam",
		Type:       "string",
		HasDefault: true,
		Default:    types.NewString("default"),
	}
	rp2 := NewReflectionParameter(rm, 0, paramDef2)

	if !rp2.IsOptional() {
		t.Error("IsOptional() = false for optional parameter, want true")
	}
}

// TestReflectionParameter_DefaultValue tests default value handling
func TestReflectionParameter_DefaultValue(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	// Parameter with default value
	defaultVal := types.NewInt(42)
	paramDef1 := &types.ParameterDef{
		Name:       "paramWithDefault",
		Type:       "int",
		HasDefault: true,
		Default:    defaultVal,
	}
	rp1 := NewReflectionParameter(rm, 0, paramDef1)

	val := rp1.GetDefaultValue()
	if val.ToInt() != 42 {
		t.Errorf("GetDefaultValue().ToInt() = %d, want 42", val.ToInt())
	}

	// Parameter without default value
	paramDef2 := &types.ParameterDef{
		Name:       "paramNoDefault",
		Type:       "int",
		HasDefault: false,
	}
	rp2 := NewReflectionParameter(rm, 0, paramDef2)

	val2 := rp2.GetDefaultValue()
	if !val2.IsNull() {
		t.Error("GetDefaultValue() should return null for parameter without default")
	}
}

// TestReflectionParameter_PassedByReference tests reference passing check
func TestReflectionParameter_PassedByReference(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	// Parameter passed by value
	paramDef1 := &types.ParameterDef{
		Name:        "byValue",
		Type:        "int",
		PassedByRef: false,
	}
	rp1 := NewReflectionParameter(rm, 0, paramDef1)

	if rp1.IsPassedByReference() {
		t.Error("IsPassedByReference() = true for by-value parameter, want false")
	}

	// Parameter passed by reference
	paramDef2 := &types.ParameterDef{
		Name:        "byRef",
		Type:        "int",
		PassedByRef: true,
	}
	rp2 := NewReflectionParameter(rm, 0, paramDef2)

	if !rp2.IsPassedByReference() {
		t.Error("IsPassedByReference() = false for by-reference parameter, want true")
	}
}

// TestReflectionParameter_Variadic tests variadic parameter check
func TestReflectionParameter_Variadic(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	// Regular parameter
	paramDef1 := &types.ParameterDef{
		Name:       "regularParam",
		Type:       "string",
		IsVariadic: false,
	}
	rp1 := NewReflectionParameter(rm, 0, paramDef1)

	if rp1.IsVariadic() {
		t.Error("IsVariadic() = true for regular parameter, want false")
	}

	// Variadic parameter
	paramDef2 := &types.ParameterDef{
		Name:       "variadicParam",
		Type:       "string",
		IsVariadic: true,
	}
	rp2 := NewReflectionParameter(rm, 1, paramDef2)

	if !rp2.IsVariadic() {
		t.Error("IsVariadic() = false for variadic parameter, want true")
	}
}

// TestReflectionParameter_AllowsNull tests null allowance check
func TestReflectionParameter_AllowsNull(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	tests := []struct {
		name       string
		paramDef   *types.ParameterDef
		allowsNull bool
	}{
		{
			name: "no type allows null",
			paramDef: &types.ParameterDef{
				Name: "param",
				Type: "",
			},
			allowsNull: true,
		},
		{
			name: "non-nullable type",
			paramDef: &types.ParameterDef{
				Name: "param",
				Type: "string",
			},
			allowsNull: false,
		},
		{
			name: "nullable type",
			paramDef: &types.ParameterDef{
				Name: "param",
				Type: "?string",
			},
			allowsNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := NewReflectionParameter(rm, 0, tt.paramDef)
			if allowsNull := rp.AllowsNull(); allowsNull != tt.allowsNull {
				t.Errorf("AllowsNull() = %v, want %v", allowsNull, tt.allowsNull)
			}
		})
	}
}

// TestReflectionParameter_String tests string representation
func TestReflectionParameter_String(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	tests := []struct {
		name     string
		paramDef *types.ParameterDef
		position int
		contains []string
	}{
		{
			name: "simple parameter",
			paramDef: &types.ParameterDef{
				Name: "param",
				Type: "string",
			},
			position: 0,
			contains: []string{"#0", "string", "$param"},
		},
		{
			name: "parameter with default",
			paramDef: &types.ParameterDef{
				Name:       "param",
				Type:       "int",
				HasDefault: true,
				Default:    types.NewInt(42),
			},
			position: 1,
			contains: []string{"#1", "int", "$param", "="},
		},
		{
			name: "reference parameter",
			paramDef: &types.ParameterDef{
				Name:        "param",
				Type:        "array",
				PassedByRef: true,
			},
			position: 0,
			contains: []string{"#0", "array", "&", "$param"},
		},
		{
			name: "variadic parameter",
			paramDef: &types.ParameterDef{
				Name:       "params",
				Type:       "string",
				IsVariadic: true,
			},
			position: 0,
			contains: []string{"#0", "string", "...", "$params"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := NewReflectionParameter(rm, tt.position, tt.paramDef)
			str := rp.String()

			for _, substr := range tt.contains {
				if !hasSubstring(str, substr) {
					t.Errorf("String() = %s, should contain %s", str, substr)
				}
			}
		})
	}
}

// TestReflectionParameter_ComplexExample tests a complex parameter scenario
func TestReflectionParameter_ComplexExample(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		NumParams:  4,
		Parameters: []*types.ParameterDef{
			{
				Name: "required",
				Type: "string",
			},
			{
				Name:       "optional",
				Type:       "?int",
				HasDefault: true,
				Default:    types.NewNull(),
			},
			{
				Name:        "byRef",
				Type:        "array",
				PassedByRef: true,
			},
			{
				Name:       "variadic",
				Type:       "mixed",
				IsVariadic: true,
			},
		},
	}

	rm := NewReflectionMethod(rc, "complexMethod", methodDef)
	params := rm.GetParameters()

	if len(params) != 4 {
		t.Fatalf("Expected 4 parameters, got %d", len(params))
	}

	// Check required parameter
	if params[0].IsOptional() {
		t.Error("First parameter should be required")
	}
	if !params[0].HasType() {
		t.Error("First parameter should have type")
	}

	// Check optional parameter
	if !params[1].IsOptional() {
		t.Error("Second parameter should be optional")
	}
	if !params[1].AllowsNull() {
		t.Error("Second parameter should allow null")
	}

	// Check reference parameter
	if !params[2].IsPassedByReference() {
		t.Error("Third parameter should be passed by reference")
	}

	// Check variadic parameter
	if !params[3].IsVariadic() {
		t.Error("Fourth parameter should be variadic")
	}
}
