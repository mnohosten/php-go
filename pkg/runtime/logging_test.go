package runtime

import (
	"bytes"
	"strings"
	"testing"
)

// TestLogLevel tests log level string representation
func TestLogLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelFatal, "FATAL"},
		{LogLevelError, "ERROR"},
		{LogLevelWarn, "WARN"},
		{LogLevelInfo, "INFO"},
		{LogLevelDebug, "DEBUG"},
		{LogLevelTrace, "TRACE"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestTextFormatter tests the text formatter
func TestTextFormatter(t *testing.T) {
	formatter := NewTextFormatter()
	formatter.TimestampFormat = "2006-01-02"
	formatter.IncludeLocation = true
	formatter.IncludeFields = true

	entry := &LogEntry{
		Level:     LogLevelInfo,
		Message:   "test message",
		Component: "test",
		File:      "test.go",
		Line:      42,
		Function:  "TestFunc",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	output := formatter.Format(entry)

	// Check that output contains expected components
	if !strings.Contains(output, "INFO") {
		t.Error("Output missing log level")
	}
	if !strings.Contains(output, "test message") {
		t.Error("Output missing message")
	}
	if !strings.Contains(output, "[test]") {
		t.Error("Output missing component")
	}
	if !strings.Contains(output, "test.go:42") {
		t.Error("Output missing file location")
	}
	if !strings.Contains(output, "TestFunc") {
		t.Error("Output missing function name")
	}
	if !strings.Contains(output, "key1=value1") {
		t.Error("Output missing field key1")
	}
	if !strings.Contains(output, "key2=123") {
		t.Error("Output missing field key2")
	}
}

// TestTextFormatterWithoutLocation tests text formatter without location
func TestTextFormatterWithoutLocation(t *testing.T) {
	formatter := NewTextFormatter()
	formatter.IncludeLocation = false

	entry := &LogEntry{
		Level:   LogLevelInfo,
		Message: "test message",
		File:    "test.go",
		Line:    42,
	}

	output := formatter.Format(entry)

	if strings.Contains(output, "test.go") {
		t.Error("Output should not contain file location when IncludeLocation=false")
	}
}

// TestTextFormatterWithoutFields tests text formatter without fields
func TestTextFormatterWithoutFields(t *testing.T) {
	formatter := NewTextFormatter()
	formatter.IncludeFields = false

	entry := &LogEntry{
		Level:   LogLevelInfo,
		Message: "test message",
		Fields: map[string]interface{}{
			"key": "value",
		},
	}

	output := formatter.Format(entry)

	if strings.Contains(output, "key=value") {
		t.Error("Output should not contain fields when IncludeFields=false")
	}
}

// TestTextFormatterColor tests colored output
func TestTextFormatterColor(t *testing.T) {
	formatter := NewTextFormatter()
	formatter.ColorEnabled = true

	tests := []struct {
		level    LogLevel
		contains string
	}{
		{LogLevelFatal, "\033[31m"}, // Red
		{LogLevelError, "\033[31m"}, // Red
		{LogLevelWarn, "\033[33m"},  // Yellow
		{LogLevelInfo, "\033[34m"},  // Blue
		{LogLevelDebug, "\033[97m"}, // White
		{LogLevelTrace, "\033[37m"}, // Gray
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			entry := &LogEntry{
				Level:   tt.level,
				Message: "test",
			}
			output := formatter.Format(entry)
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected color code %s in output", tt.contains)
			}
		})
	}
}

// TestJSONFormatter tests the JSON formatter
func TestJSONFormatter(t *testing.T) {
	formatter := NewJSONFormatter()

	entry := &LogEntry{
		Level:     LogLevelInfo,
		Message:   "test message",
		Component: "test",
		File:      "test.go",
		Line:      42,
		Function:  "TestFunc",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	output := formatter.Format(entry)

	// Check that output contains expected JSON fields
	expectedFields := []string{
		`"level":"INFO"`,
		`"message":"test message"`,
		`"component":"test"`,
		`"file":"test.go"`,
		`"line":42`,
		`"function":"TestFunc"`,
		`"key1":"value1"`,
		`"key2":123`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("JSON output missing field: %s\nGot: %s", field, output)
		}
	}
}

// TestJSONFormatterEscaping tests JSON special character escaping
func TestJSONFormatterEscaping(t *testing.T) {
	formatter := NewJSONFormatter()

	entry := &LogEntry{
		Level:   LogLevelInfo,
		Message: `test "quoted" message\nwith\tspecial chars`,
	}

	output := formatter.Format(entry)

	// Check that special characters are escaped
	if !strings.Contains(output, `\"quoted\"`) {
		t.Error("Quotes not properly escaped")
	}
	if !strings.Contains(output, `\\n`) {
		t.Error("Newlines not properly escaped")
	}
	if !strings.Contains(output, `\\t`) {
		t.Error("Tabs not properly escaped")
	}
}

// TestNewLogger tests logger creation with defaults
func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if logger.GetLevel() != LogLevelInfo {
		t.Errorf("Default level = %v, want %v", logger.GetLevel(), LogLevelInfo)
	}

	if logger.formatter == nil {
		t.Error("Logger formatter is nil")
	}

	if logger.output == nil {
		t.Error("Logger output is nil")
	}
}

// TestLoggerSetLevel tests setting log level
func TestLoggerSetLevel(t *testing.T) {
	logger := NewLogger()

	levels := []LogLevel{
		LogLevelFatal,
		LogLevelError,
		LogLevelWarn,
		LogLevelInfo,
		LogLevelDebug,
		LogLevelTrace,
	}

	for _, level := range levels {
		logger.SetLevel(level)
		if got := logger.GetLevel(); got != level {
			t.Errorf("GetLevel() = %v, want %v", got, level)
		}
	}
}

// TestLoggerSetFormatter tests setting formatter
func TestLoggerSetFormatter(t *testing.T) {
	logger := NewLogger()
	formatter := NewJSONFormatter()

	logger.SetFormatter(formatter)

	// Log a message and check it's JSON formatted
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Error("Formatter not applied correctly")
	}
}

// TestLoggerSetOutput tests setting output writer
func TestLoggerSetOutput(t *testing.T) {
	logger := NewLogger()
	buf := &bytes.Buffer{}

	logger.SetOutput(buf)
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Message not written to custom output: %s", output)
	}
}

// TestLoggerSetComponent tests setting component name
func TestLoggerSetComponent(t *testing.T) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	logger.SetComponent("TestComponent")
	logger.Info("test")

	output := buf.String()
	if !strings.Contains(output, "[TestComponent]") {
		t.Errorf("Component not included in output: %s", output)
	}
}

// TestLoggerWithField tests adding fields to logger
func TestLoggerWithField(t *testing.T) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	loggerWithField := logger.WithField("request_id", "12345")
	loggerWithField.Info("test")

	output := buf.String()
	if !strings.Contains(output, "request_id=12345") {
		t.Errorf("Field not included in output: %s", output)
	}

	// Original logger should not have the field
	buf.Reset()
	logger.Info("test")
	output = buf.String()
	if strings.Contains(output, "request_id") {
		t.Error("Field should not be in original logger")
	}
}

// TestLoggerWithFields tests adding multiple fields
func TestLoggerWithFields(t *testing.T) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	fields := map[string]interface{}{
		"user_id":    123,
		"session_id": "abc",
	}

	loggerWithFields := logger.WithFields(fields)
	loggerWithFields.Info("test")

	output := buf.String()
	if !strings.Contains(output, "user_id=123") {
		t.Error("user_id field not included")
	}
	if !strings.Contains(output, "session_id=abc") {
		t.Error("session_id field not included")
	}
}

// TestLoggerLevels tests logging at different levels
func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel LogLevel
		logFunc  func(*Logger, *bytes.Buffer)
		contains string
	}{
		{
			name:     "Error",
			logLevel: LogLevelError,
			logFunc: func(l *Logger, buf *bytes.Buffer) {
				l.SetOutput(buf)
				l.Error("error message")
			},
			contains: "ERROR: error message",
		},
		{
			name:     "Warn",
			logLevel: LogLevelWarn,
			logFunc: func(l *Logger, buf *bytes.Buffer) {
				l.SetOutput(buf)
				l.Warn("warn message")
			},
			contains: "WARN: warn message",
		},
		{
			name:     "Info",
			logLevel: LogLevelInfo,
			logFunc: func(l *Logger, buf *bytes.Buffer) {
				l.SetOutput(buf)
				l.Info("info message")
			},
			contains: "INFO: info message",
		},
		{
			name:     "Debug",
			logLevel: LogLevelDebug,
			logFunc: func(l *Logger, buf *bytes.Buffer) {
				l.SetOutput(buf)
				l.Debug("debug message")
			},
			contains: "DEBUG: debug message",
		},
		{
			name:     "Trace",
			logLevel: LogLevelTrace,
			logFunc: func(l *Logger, buf *bytes.Buffer) {
				l.SetOutput(buf)
				l.Trace("trace message")
			},
			contains: "TRACE: trace message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			logger.SetLevel(tt.logLevel)
			buf := &bytes.Buffer{}
			tt.logFunc(logger, buf)

			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected %q in output, got: %s", tt.contains, output)
			}
		})
	}
}

// TestLoggerFormattedMessages tests formatted logging methods
func TestLoggerFormattedMessages(t *testing.T) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Test Errorf
	buf.Reset()
	logger.Errorf("error %d: %s", 404, "not found")
	if !strings.Contains(buf.String(), "error 404: not found") {
		t.Error("Errorf not formatting correctly")
	}

	// Test Warnf
	buf.Reset()
	logger.Warnf("warning %d", 123)
	if !strings.Contains(buf.String(), "warning 123") {
		t.Error("Warnf not formatting correctly")
	}

	// Test Infof
	buf.Reset()
	logger.Infof("info %s", "test")
	if !strings.Contains(buf.String(), "info test") {
		t.Error("Infof not formatting correctly")
	}

	// Test Debugf
	buf.Reset()
	logger.SetLevel(LogLevelDebug)
	logger.Debugf("debug %v", true)
	if !strings.Contains(buf.String(), "debug true") {
		t.Error("Debugf not formatting correctly")
	}

	// Test Tracef
	buf.Reset()
	logger.SetLevel(LogLevelTrace)
	logger.Tracef("trace %d", 42)
	if !strings.Contains(buf.String(), "trace 42") {
		t.Error("Tracef not formatting correctly")
	}
}

// TestLoggerFiltering tests that messages below log level are filtered
func TestLoggerFiltering(t *testing.T) {
	logger := NewLogger()
	logger.SetLevel(LogLevelWarn)
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// These should be logged (Warn level and above)
	logger.Error("error")
	logger.Warn("warn")

	// These should NOT be logged (below Warn level)
	logger.Info("info")
	logger.Debug("debug")
	logger.Trace("trace")

	output := buf.String()

	if !strings.Contains(output, "error") {
		t.Error("Error message should be logged")
	}
	if !strings.Contains(output, "warn") {
		t.Error("Warn message should be logged")
	}
	if strings.Contains(output, "info") {
		t.Error("Info message should be filtered")
	}
	if strings.Contains(output, "debug") {
		t.Error("Debug message should be filtered")
	}
	if strings.Contains(output, "trace") {
		t.Error("Trace message should be filtered")
	}
}

// TestGlobalLogger tests the global logger instance
func TestGlobalLogger(t *testing.T) {
	logger := GetGlobalLogger()
	if logger == nil {
		t.Fatal("GetGlobalLogger returned nil")
	}

	// Create a new logger and set it as global
	newLogger := NewLogger()
	newLogger.SetLevel(LogLevelDebug)
	SetGlobalLogger(newLogger)

	// Verify the global logger was updated
	if GetGlobalLogger() != newLogger {
		t.Error("SetGlobalLogger did not update global logger")
	}

	if GetGlobalLogger().GetLevel() != LogLevelDebug {
		t.Error("Global logger level not updated")
	}
}

// TestGlobalLogFunctions tests package-level logging functions
func TestGlobalLogFunctions(t *testing.T) {
	// Save original global logger
	originalLogger := GetGlobalLogger()
	defer SetGlobalLogger(originalLogger)

	// Create test logger
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	SetGlobalLogger(logger)

	// Test Error
	buf.Reset()
	Error("error message")
	if !strings.Contains(buf.String(), "ERROR: error message") {
		t.Error("Error() not working")
	}

	// Test Errorf
	buf.Reset()
	Errorf("error %d", 123)
	if !strings.Contains(buf.String(), "error 123") {
		t.Error("Errorf() not working")
	}

	// Test Warn
	buf.Reset()
	Warn("warn message")
	if !strings.Contains(buf.String(), "WARN: warn message") {
		t.Error("Warn() not working")
	}

	// Test Warnf
	buf.Reset()
	Warnf("warn %s", "test")
	if !strings.Contains(buf.String(), "warn test") {
		t.Error("Warnf() not working")
	}

	// Test Info
	buf.Reset()
	Info("info message")
	if !strings.Contains(buf.String(), "INFO: info message") {
		t.Error("Info() not working")
	}

	// Test Infof
	buf.Reset()
	Infof("info %d", 456)
	if !strings.Contains(buf.String(), "info 456") {
		t.Error("Infof() not working")
	}

	// Test Debug
	buf.Reset()
	logger.SetLevel(LogLevelDebug)
	Debug("debug message")
	if !strings.Contains(buf.String(), "DEBUG: debug message") {
		t.Error("Debug() not working")
	}

	// Test Debugf
	buf.Reset()
	Debugf("debug %v", true)
	if !strings.Contains(buf.String(), "debug true") {
		t.Error("Debugf() not working")
	}

	// Test Trace
	buf.Reset()
	logger.SetLevel(LogLevelTrace)
	Trace("trace message")
	if !strings.Contains(buf.String(), "TRACE: trace message") {
		t.Error("Trace() not working")
	}

	// Test Tracef
	buf.Reset()
	Tracef("trace %d", 789)
	if !strings.Contains(buf.String(), "trace 789") {
		t.Error("Tracef() not working")
	}
}

// TestEscapeJSON tests JSON escaping function
func TestEscapeJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`hello`, `hello`},
		{`"quoted"`, `\"quoted\"`},
		{`back\slash`, `back\\slash`},
		{"new\nline", `new\nline`},
		{"carriage\rreturn", `carriage\rreturn`},
		{"tab\there", `tab\there`},
		{`all "special"\chars\n\r\t`, `all \"special\"\\chars\\n\\r\\t`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := escapeJSON(tt.input); got != tt.expected {
				t.Errorf("escapeJSON(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// BenchmarkTextFormatter benchmarks text formatting
func BenchmarkTextFormatter(b *testing.B) {
	formatter := NewTextFormatter()
	entry := &LogEntry{
		Level:     LogLevelInfo,
		Message:   "benchmark message",
		Component: "bench",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.Format(entry)
	}
}

// BenchmarkJSONFormatter benchmarks JSON formatting
func BenchmarkJSONFormatter(b *testing.B) {
	formatter := NewJSONFormatter()
	entry := &LogEntry{
		Level:     LogLevelInfo,
		Message:   "benchmark message",
		Component: "bench",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.Format(entry)
	}
}

// BenchmarkLoggerInfo benchmarks Info logging
func BenchmarkLoggerInfo(b *testing.B) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

// BenchmarkLoggerWithFields benchmarks logging with fields
func BenchmarkLoggerWithFields(b *testing.B) {
	logger := NewLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loggerWithFields := logger.WithFields(fields)
		loggerWithFields.Info("benchmark message")
	}
}
