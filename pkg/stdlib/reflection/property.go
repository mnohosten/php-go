package reflection

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ReflectionProperty represents a property reflection
type ReflectionProperty struct {
	// The class this property belongs to
	class *ReflectionClass

	// Property name
	name string

	// Property definition
	property *types.PropertyDef
}

// NewReflectionProperty creates a new ReflectionProperty
func NewReflectionProperty(class *ReflectionClass, name string, property *types.PropertyDef) *ReflectionProperty {
	return &ReflectionProperty{
		class:    class,
		name:     name,
		property: property,
	}
}

// GetName returns the property name
func (rp *ReflectionProperty) GetName() string {
	return rp.name
}

// GetDeclaringClass returns the class that declared this property
func (rp *ReflectionProperty) GetDeclaringClass() *ReflectionClass {
	return rp.class
}

// IsPublic checks if the property is public
func (rp *ReflectionProperty) IsPublic() bool {
	return rp.property.Visibility == types.VisibilityPublic
}

// IsProtected checks if the property is protected
func (rp *ReflectionProperty) IsProtected() bool {
	return rp.property.Visibility == types.VisibilityProtected
}

// IsPrivate checks if the property is private
func (rp *ReflectionProperty) IsPrivate() bool {
	return rp.property.Visibility == types.VisibilityPrivate
}

// IsStatic checks if the property is static
func (rp *ReflectionProperty) IsStatic() bool {
	return rp.property.IsStatic
}

// IsReadOnly checks if the property is readonly
func (rp *ReflectionProperty) IsReadOnly() bool {
	return rp.property.IsReadOnly
}

// GetType returns the property type declaration
func (rp *ReflectionProperty) GetType() string {
	return rp.property.Type
}

// HasType checks if the property has a type declaration
func (rp *ReflectionProperty) HasType() bool {
	return rp.property.Type != ""
}

// GetDefaultValue returns the default value of the property
func (rp *ReflectionProperty) GetDefaultValue() *types.Value {
	if rp.property.HasDefault {
		return rp.property.Default
	}
	return types.NewNull()
}

// HasDefaultValue checks if the property has a default value
func (rp *ReflectionProperty) HasDefaultValue() bool {
	return rp.property.HasDefault
}

// GetValue gets the value of the property from an object
func (rp *ReflectionProperty) GetValue(obj *types.Object) (*types.Value, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot get property value on null object")
	}

	if prop, ok := obj.Properties[rp.name]; ok {
		return prop.Value, nil
	}

	// Check for static property
	if rp.property.IsStatic && rp.class.class != nil {
		if val, ok := rp.class.class.StaticProperties[rp.name]; ok {
			return val, nil
		}
	}

	return types.NewNull(), nil
}

// SetValue sets the value of the property on an object
func (rp *ReflectionProperty) SetValue(obj *types.Object, value *types.Value) error {
	if obj == nil {
		return fmt.Errorf("Cannot set property value on null object")
	}

	if rp.property.IsReadOnly {
		return fmt.Errorf("Cannot modify readonly property %s", rp.name)
	}

	if prop, ok := obj.Properties[rp.name]; ok {
		prop.Value = value
		return nil
	}

	// Create property if it doesn't exist
	obj.Properties[rp.name] = &types.Property{
		Value:      value,
		Visibility: rp.property.Visibility,
		IsStatic:   false,
	}

	return nil
}

// SetAccessible makes the property accessible (bypasses visibility)
func (rp *ReflectionProperty) SetAccessible(accessible bool) {
	// In Go implementation, we can always access properties
	// This is mainly for API compatibility with PHP
	// In a full implementation, we'd track accessibility state
}

// String returns a string representation of the property
func (rp *ReflectionProperty) String() string {
	visibility := "public"
	if rp.IsProtected() {
		visibility = "protected"
	} else if rp.IsPrivate() {
		visibility = "private"
	}

	static := ""
	if rp.IsStatic() {
		static = "static "
	}

	readonly := ""
	if rp.IsReadOnly() {
		readonly = "readonly "
	}

	typeStr := ""
	if rp.HasType() {
		typeStr = rp.GetType() + " "
	}

	return fmt.Sprintf("Property [ %s%s%s%s$%s ]", visibility, " ", static, readonly+typeStr, rp.name)
}
