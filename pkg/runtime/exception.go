package runtime

import (
	"fmt"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// Exception represents a PHP exception
// Implements PHP's Throwable interface
type Exception struct {
	// Message - exception message
	Message string

	// Code - exception code
	Code int64

	// File - filename where exception was thrown
	File string

	// Line - line number where exception was thrown
	Line int

	// Previous - previous exception (for chaining)
	Previous *Exception

	// Trace - stack trace
	Trace *StackTrace

	// ClassName - name of the exception class
	ClassName string
}

// NewException creates a new exception
func NewException(message string, code int64, previous *Exception) *Exception {
	return &Exception{
		Message:   message,
		Code:      code,
		Previous:  previous,
		Trace:     NewStackTrace(),
		ClassName: "Exception",
	}
}

// NewExceptionWithTrace creates a new exception with a stack trace
func NewExceptionWithTrace(message string, code int64, previous *Exception, trace *StackTrace) *Exception {
	return &Exception{
		Message:   message,
		Code:      code,
		Previous:  previous,
		Trace:     trace,
		ClassName: "Exception",
	}
}

// GetMessage returns the exception message
func (e *Exception) GetMessage() string {
	return e.Message
}

// GetCode returns the exception code
func (e *Exception) GetCode() int64 {
	return e.Code
}

// GetFile returns the filename where exception was thrown
func (e *Exception) GetFile() string {
	return e.File
}

// GetLine returns the line number where exception was thrown
func (e *Exception) GetLine() int {
	return e.Line
}

// GetPrevious returns the previous exception in the chain
func (e *Exception) GetPrevious() *Exception {
	return e.Previous
}

// GetTrace returns the stack trace
func (e *Exception) GetTrace() *StackTrace {
	return e.Trace
}

// GetTraceAsString returns the stack trace as a string
func (e *Exception) GetTraceAsString() string {
	if e.Trace == nil {
		return ""
	}
	return e.Trace.String()
}

// String returns a string representation of the exception
func (e *Exception) String() string {
	var sb strings.Builder

	// Class name and message
	sb.WriteString(fmt.Sprintf("%s: %s", e.ClassName, e.Message))

	// File and line if available
	if e.File != "" {
		sb.WriteString(fmt.Sprintf(" in %s:%d", e.File, e.Line))
	}

	// Stack trace
	if e.Trace != nil && len(e.Trace.Frames) > 0 {
		sb.WriteString("\n")
		sb.WriteString(e.GetTraceAsString())
	}

	return sb.String()
}

// ToObject converts the exception to a PHP object
func (e *Exception) ToObject() *types.Object {
	obj := &types.Object{
		ClassName:  e.ClassName,
		Properties: make(map[string]*types.Property),
	}

	// Set properties
	obj.Properties["message"] = &types.Property{
		Value:      types.NewString(e.Message),
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}
	obj.Properties["code"] = &types.Property{
		Value:      types.NewInt(e.Code),
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}
	obj.Properties["file"] = &types.Property{
		Value:      types.NewString(e.File),
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}
	obj.Properties["line"] = &types.Property{
		Value:      types.NewInt(int64(e.Line)),
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}

	// Previous exception
	if e.Previous != nil {
		obj.Properties["previous"] = &types.Property{
			Value:      types.NewObject(e.Previous.ToObject()),
			Visibility: types.VisibilityPublic,
			IsStatic:   false,
		}
	} else {
		obj.Properties["previous"] = &types.Property{
			Value:      types.NewNull(),
			Visibility: types.VisibilityPublic,
			IsStatic:   false,
		}
	}

	return obj
}

// FromObject creates an exception from a PHP object
func FromObject(obj *types.Object) (*Exception, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot create exception from null object")
	}

	exc := &Exception{
		ClassName: obj.ClassName,
		Trace:     NewStackTrace(),
	}

	// Extract message
	if msgProp, ok := obj.Properties["message"]; ok {
		exc.Message = msgProp.Value.ToString()
	}

	// Extract code
	if codeProp, ok := obj.Properties["code"]; ok {
		exc.Code = codeProp.Value.ToInt()
	}

	// Extract file
	if fileProp, ok := obj.Properties["file"]; ok {
		exc.File = fileProp.Value.ToString()
	}

	// Extract line
	if lineProp, ok := obj.Properties["line"]; ok {
		exc.Line = int(lineProp.Value.ToInt())
	}

	// Extract previous
	if prevProp, ok := obj.Properties["previous"]; ok {
		if prevProp.Value.Type() == types.TypeObject {
			prevExc, err := FromObject(prevProp.Value.ToObject())
			if err == nil {
				exc.Previous = prevExc
			}
		}
	}

	return exc, nil
}

// ============================================================================
// Exception Subclasses (PHP standard exception hierarchy)
// ============================================================================

// ErrorException represents a PHP ErrorException
type ErrorException struct {
	*Exception
	Severity ErrorType
}

// NewErrorException creates a new error exception
func NewErrorException(message string, code int64, severity ErrorType, previous *Exception) *ErrorException {
	return &ErrorException{
		Exception: &Exception{
			Message:   message,
			Code:      code,
			Previous:  previous,
			Trace:     NewStackTrace(),
			ClassName: "ErrorException",
		},
		Severity: severity,
	}
}

// GetSeverity returns the error severity
func (e *ErrorException) GetSeverity() ErrorType {
	return e.Severity
}

// LogicException represents a PHP LogicException
type LogicException struct {
	*Exception
}

// NewLogicException creates a new logic exception
func NewLogicException(message string, code int64, previous *Exception) *LogicException {
	return &LogicException{
		Exception: &Exception{
			Message:   message,
			Code:      code,
			Previous:  previous,
			Trace:     NewStackTrace(),
			ClassName: "LogicException",
		},
	}
}

// RuntimeException represents a PHP RuntimeException
type RuntimeException struct {
	*Exception
}

// NewRuntimeException creates a new runtime exception
func NewRuntimeException(message string, code int64, previous *Exception) *RuntimeException {
	return &RuntimeException{
		Exception: &Exception{
			Message:   message,
			Code:      code,
			Previous:  previous,
			Trace:     NewStackTrace(),
			ClassName: "RuntimeException",
		},
	}
}

// InvalidArgumentException represents a PHP InvalidArgumentException
type InvalidArgumentException struct {
	*LogicException
}

// NewInvalidArgumentException creates a new invalid argument exception
func NewInvalidArgumentException(message string, code int64, previous *Exception) *InvalidArgumentException {
	return &InvalidArgumentException{
		LogicException: &LogicException{
			Exception: &Exception{
				Message:   message,
				Code:      code,
				Previous:  previous,
				Trace:     NewStackTrace(),
				ClassName: "InvalidArgumentException",
			},
		},
	}
}

// OutOfBoundsException represents a PHP OutOfBoundsException
type OutOfBoundsException struct {
	*LogicException
}

// NewOutOfBoundsException creates a new out of bounds exception
func NewOutOfBoundsException(message string, code int64, previous *Exception) *OutOfBoundsException {
	return &OutOfBoundsException{
		LogicException: &LogicException{
			Exception: &Exception{
				Message:   message,
				Code:      code,
				Previous:  previous,
				Trace:     NewStackTrace(),
				ClassName: "OutOfBoundsException",
			},
		},
	}
}

// OutOfRangeException represents a PHP OutOfRangeException
type OutOfRangeException struct {
	*LogicException
}

// NewOutOfRangeException creates a new out of range exception
func NewOutOfRangeException(message string, code int64, previous *Exception) *OutOfRangeException {
	return &OutOfRangeException{
		LogicException: &LogicException{
			Exception: &Exception{
				Message:   message,
				Code:      code,
				Previous:  previous,
				Trace:     NewStackTrace(),
				ClassName: "OutOfRangeException",
			},
		},
	}
}

// ============================================================================
// Exception Helpers
// ============================================================================

// IsException checks if a value is an exception
func IsException(val *types.Value) bool {
	if val.Type() != types.TypeObject {
		return false
	}

	obj := val.ToObject()
	// Check if class name ends with "Exception" or is "Error" or "Throwable"
	className := obj.ClassName
	return strings.HasSuffix(className, "Exception") ||
		className == "Error" ||
		className == "Throwable"
}

// GetExceptionFromValue extracts an exception from a value
func GetExceptionFromValue(val *types.Value) (*Exception, error) {
	if val.Type() == types.TypeResource {
		// Check if it's an exception resource
		resource := val.ToResource()
		data := resource.Data()
		if exc, ok := data.(*Exception); ok {
			return exc, nil
		}
	}

	if val.Type() == types.TypeObject {
		// Try to convert object to exception
		return FromObject(val.ToObject())
	}

	return nil, fmt.Errorf("Value is not an exception")
}

// WrapExceptionAsValue wraps an exception as a PHP value (object)
func WrapExceptionAsValue(exc *Exception) *types.Value {
	// Convert exception to object so PHP code can access properties
	return types.NewObject(exc.ToObject())
}
