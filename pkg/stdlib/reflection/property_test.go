package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionProperty_BasicInfo tests basic property information
func TestReflectionProperty_BasicInfo(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
	}
	rc := NewReflectionClass("TestClass", classEntry)

	propDef := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		Type:       "string",
		HasDefault: true,
		Default:    types.NewString("test"),
	}

	rp := NewReflectionProperty(rc, "testProp", propDef)

	// Test GetName
	if name := rp.GetName(); name != "testProp" {
		t.Errorf("GetName() = %s, want testProp", name)
	}

	// Test GetDeclaringClass
	if declaringClass := rp.GetDeclaringClass(); declaringClass != rc {
		t.Error("GetDeclaringClass() returned wrong class")
	}
}

// TestReflectionProperty_Visibility tests visibility checks
func TestReflectionProperty_Visibility(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	tests := []struct {
		name       string
		visibility types.PropertyVisibility
		isPublic   bool
		isProtected bool
		isPrivate  bool
	}{
		{
			name:       "public property",
			visibility: types.VisibilityPublic,
			isPublic:   true,
			isProtected: false,
			isPrivate:  false,
		},
		{
			name:       "protected property",
			visibility: types.VisibilityProtected,
			isPublic:   false,
			isProtected: true,
			isPrivate:  false,
		},
		{
			name:       "private property",
			visibility: types.VisibilityPrivate,
			isPublic:   false,
			isProtected: false,
			isPrivate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propDef := &types.PropertyDef{
				Visibility: tt.visibility,
			}
			rp := NewReflectionProperty(rc, "prop", propDef)

			if isPublic := rp.IsPublic(); isPublic != tt.isPublic {
				t.Errorf("IsPublic() = %v, want %v", isPublic, tt.isPublic)
			}
			if isProtected := rp.IsProtected(); isProtected != tt.isProtected {
				t.Errorf("IsProtected() = %v, want %v", isProtected, tt.isProtected)
			}
			if isPrivate := rp.IsPrivate(); isPrivate != tt.isPrivate {
				t.Errorf("IsPrivate() = %v, want %v", isPrivate, tt.isPrivate)
			}
		})
	}
}

// TestReflectionProperty_Static tests static property checks
func TestReflectionProperty_Static(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Non-static property
	propDef1 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}
	rp1 := NewReflectionProperty(rc, "instanceProp", propDef1)

	if rp1.IsStatic() {
		t.Error("IsStatic() = true for instance property, want false")
	}

	// Static property
	propDef2 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		IsStatic:   true,
	}
	rp2 := NewReflectionProperty(rc, "staticProp", propDef2)

	if !rp2.IsStatic() {
		t.Error("IsStatic() = false for static property, want true")
	}
}

// TestReflectionProperty_ReadOnly tests readonly property checks
func TestReflectionProperty_ReadOnly(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Regular property
	propDef1 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		IsReadOnly: false,
	}
	rp1 := NewReflectionProperty(rc, "regularProp", propDef1)

	if rp1.IsReadOnly() {
		t.Error("IsReadOnly() = true for regular property, want false")
	}

	// Readonly property
	propDef2 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		IsReadOnly: true,
	}
	rp2 := NewReflectionProperty(rc, "readonlyProp", propDef2)

	if !rp2.IsReadOnly() {
		t.Error("IsReadOnly() = false for readonly property, want true")
	}
}

// TestReflectionProperty_Type tests type information
func TestReflectionProperty_Type(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Property with type
	propDef1 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		Type:       "string",
	}
	rp1 := NewReflectionProperty(rc, "typedProp", propDef1)

	if !rp1.HasType() {
		t.Error("HasType() = false for typed property, want true")
	}
	if typ := rp1.GetType(); typ != "string" {
		t.Errorf("GetType() = %s, want string", typ)
	}

	// Property without type
	propDef2 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		Type:       "",
	}
	rp2 := NewReflectionProperty(rc, "untypedProp", propDef2)

	if rp2.HasType() {
		t.Error("HasType() = true for untyped property, want false")
	}
}

// TestReflectionProperty_DefaultValue tests default value handling
func TestReflectionProperty_DefaultValue(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Property with default value
	defaultVal := types.NewString("default")
	propDef1 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		HasDefault: true,
		Default:    defaultVal,
	}
	rp1 := NewReflectionProperty(rc, "propWithDefault", propDef1)

	if !rp1.HasDefaultValue() {
		t.Error("HasDefaultValue() = false, want true")
	}
	val := rp1.GetDefaultValue()
	if val.ToString() != "default" {
		t.Errorf("GetDefaultValue().ToString() = %s, want default", val.ToString())
	}

	// Property without default value
	propDef2 := &types.PropertyDef{
		Visibility: types.VisibilityPublic,
		HasDefault: false,
	}
	rp2 := NewReflectionProperty(rc, "propNoDefault", propDef2)

	if rp2.HasDefaultValue() {
		t.Error("HasDefaultValue() = true, want false")
	}
	val2 := rp2.GetDefaultValue()
	if !val2.IsNull() {
		t.Error("GetDefaultValue() should return null for property without default")
	}
}

// TestReflectionProperty_GetValue tests getting property value from object
func TestReflectionProperty_GetValue(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"prop": {
				Visibility: types.VisibilityPublic,
			},
		},
	}
	rc := NewReflectionClass("TestClass", classEntry)

	// Create object with property
	obj := &types.Object{
		ClassName:  "TestClass",
		ClassEntry: classEntry,
		Properties: map[string]*types.Property{
			"prop": {
				Value:      types.NewString("test value"),
				Visibility: types.VisibilityPublic,
			},
		},
	}

	propDef := classEntry.Properties["prop"]
	rp := NewReflectionProperty(rc, "prop", propDef)

	// Test GetValue
	val, err := rp.GetValue(obj)
	if err != nil {
		t.Fatalf("GetValue() error: %v", err)
	}
	if val.ToString() != "test value" {
		t.Errorf("GetValue().ToString() = %s, want test value", val.ToString())
	}

	// Test GetValue with nil object
	_, err = rp.GetValue(nil)
	if err == nil {
		t.Error("GetValue(nil) should return error")
	}
}

// TestReflectionProperty_SetValue tests setting property value on object
func TestReflectionProperty_SetValue(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"prop": {
				Visibility: types.VisibilityPublic,
			},
		},
	}
	rc := NewReflectionClass("TestClass", classEntry)

	// Create object
	obj := &types.Object{
		ClassName:  "TestClass",
		ClassEntry: classEntry,
		Properties: map[string]*types.Property{
			"prop": {
				Value:      types.NewString("old value"),
				Visibility: types.VisibilityPublic,
			},
		},
	}

	propDef := classEntry.Properties["prop"]
	rp := NewReflectionProperty(rc, "prop", propDef)

	// Test SetValue
	newVal := types.NewString("new value")
	err := rp.SetValue(obj, newVal)
	if err != nil {
		t.Fatalf("SetValue() error: %v", err)
	}

	// Verify value was set
	val, _ := rp.GetValue(obj)
	if val.ToString() != "new value" {
		t.Errorf("After SetValue, GetValue() = %s, want new value", val.ToString())
	}

	// Test SetValue with nil object
	err = rp.SetValue(nil, newVal)
	if err == nil {
		t.Error("SetValue(nil, ...) should return error")
	}
}

// TestReflectionProperty_SetValue_ReadOnly tests setting readonly property fails
func TestReflectionProperty_SetValue_ReadOnly(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"readonlyProp": {
				Visibility: types.VisibilityPublic,
				IsReadOnly: true,
			},
		},
	}
	rc := NewReflectionClass("TestClass", classEntry)

	obj := &types.Object{
		ClassName:  "TestClass",
		ClassEntry: classEntry,
		Properties: map[string]*types.Property{
			"readonlyProp": {
				Value:      types.NewString("readonly"),
				Visibility: types.VisibilityPublic,
			},
		},
	}

	propDef := classEntry.Properties["readonlyProp"]
	rp := NewReflectionProperty(rc, "readonlyProp", propDef)

	// Try to set readonly property
	err := rp.SetValue(obj, types.NewString("modified"))
	if err == nil {
		t.Error("SetValue() on readonly property should return error")
	}
}

// TestReflectionProperty_StaticValue tests getting/setting static property value
func TestReflectionProperty_StaticValue(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"staticProp": {
				Visibility: types.VisibilityPublic,
				IsStatic:   true,
			},
		},
		StaticProperties: map[string]*types.Value{
			"staticProp": types.NewString("static value"),
		},
	}
	rc := NewReflectionClass("TestClass", classEntry)

	obj := &types.Object{
		ClassName:  "TestClass",
		ClassEntry: classEntry,
		Properties: make(map[string]*types.Property),
	}

	propDef := classEntry.Properties["staticProp"]
	rp := NewReflectionProperty(rc, "staticProp", propDef)

	// Test GetValue for static property
	val, err := rp.GetValue(obj)
	if err != nil {
		t.Fatalf("GetValue() error: %v", err)
	}
	if val.ToString() != "static value" {
		t.Errorf("GetValue().ToString() = %s, want static value", val.ToString())
	}
}

// TestReflectionProperty_String tests string representation
func TestReflectionProperty_String(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	tests := []struct {
		name     string
		propDef  *types.PropertyDef
		propName string
		contains []string
	}{
		{
			name: "public property",
			propDef: &types.PropertyDef{
				Visibility: types.VisibilityPublic,
			},
			propName: "prop",
			contains: []string{"public", "$prop"},
		},
		{
			name: "static property",
			propDef: &types.PropertyDef{
				Visibility: types.VisibilityPublic,
				IsStatic:   true,
			},
			propName: "staticProp",
			contains: []string{"static", "$staticProp"},
		},
		{
			name: "readonly property",
			propDef: &types.PropertyDef{
				Visibility: types.VisibilityPublic,
				IsReadOnly: true,
			},
			propName: "readonlyProp",
			contains: []string{"readonly", "$readonlyProp"},
		},
		{
			name: "typed property",
			propDef: &types.PropertyDef{
				Visibility: types.VisibilityPublic,
				Type:       "string",
			},
			propName: "typedProp",
			contains: []string{"string", "$typedProp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := NewReflectionProperty(rc, tt.propName, tt.propDef)
			str := rp.String()

			for _, substr := range tt.contains {
				if !containsString(str, substr) {
					t.Errorf("String() = %s, should contain %s", str, substr)
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
