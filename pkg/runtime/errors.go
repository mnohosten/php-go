package runtime

// ErrorType represents PHP error types
type ErrorType int

const (
	// Error levels (from PHP)
	E_ERROR             ErrorType = 1      // Fatal run-time errors
	E_WARNING           ErrorType = 2      // Run-time warnings (non-fatal errors)
	E_PARSE             ErrorType = 4      // Compile-time parse errors
	E_NOTICE            ErrorType = 8      // Run-time notices
	E_CORE_ERROR        ErrorType = 16     // Fatal errors during PHP's initial startup
	E_CORE_WARNING      ErrorType = 32     // Warnings during PHP's initial startup
	E_COMPILE_ERROR     ErrorType = 64     // Fatal compile-time errors
	E_COMPILE_WARNING   ErrorType = 128    // Compile-time warnings
	E_USER_ERROR        ErrorType = 256    // User-generated error message
	E_USER_WARNING      ErrorType = 512    // User-generated warning message
	E_USER_NOTICE       ErrorType = 1024   // User-generated notice message
	E_STRICT            ErrorType = 2048   // Enable to have PHP suggest changes
	E_RECOVERABLE_ERROR ErrorType = 4096   // Catchable fatal error
	E_DEPRECATED        ErrorType = 8192   // Run-time notices
	E_USER_DEPRECATED   ErrorType = 16384  // User-generated warning message
	E_ALL               ErrorType = 32767  // All errors and warnings
)

// String returns the string representation of an error type
func (et ErrorType) String() string {
	switch et {
	case E_ERROR:
		return "Fatal error"
	case E_WARNING:
		return "Warning"
	case E_PARSE:
		return "Parse error"
	case E_NOTICE:
		return "Notice"
	case E_CORE_ERROR:
		return "Core error"
	case E_CORE_WARNING:
		return "Core warning"
	case E_COMPILE_ERROR:
		return "Compile error"
	case E_COMPILE_WARNING:
		return "Compile warning"
	case E_USER_ERROR:
		return "User error"
	case E_USER_WARNING:
		return "User warning"
	case E_USER_NOTICE:
		return "User notice"
	case E_STRICT:
		return "Strict standards"
	case E_RECOVERABLE_ERROR:
		return "Recoverable error"
	case E_DEPRECATED:
		return "Deprecated"
	case E_USER_DEPRECATED:
		return "User deprecated"
	default:
		return "Unknown error"
	}
}

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	File     string
	Line     int
	Function string
	Class    string
	Type     string // "->" or "::"
	Args     []interface{}
}

// StackTrace represents a stack trace
type StackTrace struct {
	Frames []*StackFrame
}

// NewStackTrace creates a new stack trace
func NewStackTrace() *StackTrace {
	return &StackTrace{
		Frames: make([]*StackFrame, 0),
	}
}

// AddFrame adds a frame to the stack trace
func (st *StackTrace) AddFrame(frame *StackFrame) {
	st.Frames = append(st.Frames, frame)
}

// String returns a string representation of the stack trace
func (st *StackTrace) String() string {
	result := "Stack trace:\n"
	for i, frame := range st.Frames {
		result += formatFrame(i, frame) + "\n"
	}
	return result
}

// formatFrame formats a single stack frame
func formatFrame(index int, frame *StackFrame) string {
	location := ""
	if frame.File != "" {
		location = frame.File
		if frame.Line > 0 {
			location += ":" + string(rune(frame.Line))
		}
	}

	function := frame.Function
	if frame.Class != "" {
		function = frame.Class + frame.Type + frame.Function
	}

	return "#" + string(rune(index)) + " " + location + " " + function + "()"
}
