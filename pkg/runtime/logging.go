package runtime

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// Log levels from most to least severe
	LogLevelFatal LogLevel = iota // Fatal errors that cause program termination
	LogLevelError                  // Errors that prevent normal operation
	LogLevelWarn                   // Warning messages for potentially harmful situations
	LogLevelInfo                   // Informational messages about program execution
	LogLevelDebug                  // Detailed debugging information
	LogLevelTrace                  // Very detailed tracing information
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelFatal:
		return "FATAL"
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// LogFormatter defines the interface for log formatters
type LogFormatter interface {
	Format(entry *LogEntry) string
}

// LogEntry represents a single log entry
type LogEntry struct {
	Time      time.Time
	Level     LogLevel
	Message   string
	Fields    map[string]interface{}
	File      string
	Line      int
	Function  string
	Component string
}

// TextFormatter formats log entries as plain text
type TextFormatter struct {
	TimestampFormat string
	IncludeLocation bool
	IncludeFields   bool
	ColorEnabled    bool
}

// NewTextFormatter creates a new text formatter with default settings
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		IncludeLocation: true,
		IncludeFields:   true,
		ColorEnabled:    false,
	}
}

// Format formats a log entry as plain text
func (f *TextFormatter) Format(entry *LogEntry) string {
	// Timestamp
	timestamp := entry.Time.Format(f.TimestampFormat)

	// Level with optional color
	level := entry.Level.String()
	if f.ColorEnabled {
		level = f.colorizeLevel(entry.Level, level)
	}

	// Build message
	msg := fmt.Sprintf("[%s] %s: %s", timestamp, level, entry.Message)

	// Add component if present
	if entry.Component != "" {
		msg = fmt.Sprintf("[%s] %s [%s]: %s", timestamp, level, entry.Component, entry.Message)
	}

	// Add location if enabled
	if f.IncludeLocation && entry.File != "" {
		msg += fmt.Sprintf(" (%s:%d", entry.File, entry.Line)
		if entry.Function != "" {
			msg += fmt.Sprintf(" in %s", entry.Function)
		}
		msg += ")"
	}

	// Add fields if enabled and present
	if f.IncludeFields && len(entry.Fields) > 0 {
		msg += " "
		for k, v := range entry.Fields {
			msg += fmt.Sprintf("%s=%v ", k, v)
		}
	}

	return msg
}

// colorizeLevel adds ANSI color codes to log levels
func (f *TextFormatter) colorizeLevel(level LogLevel, text string) string {
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
		colorGray   = "\033[37m"
		colorWhite  = "\033[97m"
	)

	switch level {
	case LogLevelFatal, LogLevelError:
		return colorRed + text + colorReset
	case LogLevelWarn:
		return colorYellow + text + colorReset
	case LogLevelInfo:
		return colorBlue + text + colorReset
	case LogLevelDebug:
		return colorWhite + text + colorReset
	case LogLevelTrace:
		return colorGray + text + colorReset
	default:
		return text
	}
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct {
	PrettyPrint bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		PrettyPrint: false,
	}
}

// Format formats a log entry as JSON
func (f *JSONFormatter) Format(entry *LogEntry) string {
	// Build JSON manually for better control and performance
	json := "{"
	json += fmt.Sprintf(`"time":"%s",`, entry.Time.Format(time.RFC3339))
	json += fmt.Sprintf(`"level":"%s",`, entry.Level.String())
	json += fmt.Sprintf(`"message":"%s"`, escapeJSON(entry.Message))

	if entry.Component != "" {
		json += fmt.Sprintf(`,"component":"%s"`, escapeJSON(entry.Component))
	}

	if entry.File != "" {
		json += fmt.Sprintf(`,"file":"%s"`, escapeJSON(entry.File))
		json += fmt.Sprintf(`,"line":%d`, entry.Line)
	}

	if entry.Function != "" {
		json += fmt.Sprintf(`,"function":"%s"`, escapeJSON(entry.Function))
	}

	if len(entry.Fields) > 0 {
		json += `,"fields":{`
		first := true
		for k, v := range entry.Fields {
			if !first {
				json += ","
			}
			json += fmt.Sprintf(`"%s":`, escapeJSON(k))
			switch v := v.(type) {
			case string:
				json += fmt.Sprintf(`"%s"`, escapeJSON(v))
			case int, int64, float64, bool:
				json += fmt.Sprintf(`%v`, v)
			default:
				json += fmt.Sprintf(`"%v"`, v)
			}
			first = false
		}
		json += "}"
	}

	json += "}"
	return json
}

// escapeJSON escapes special characters in JSON strings
func escapeJSON(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '"':
			result += `\"`
		case '\\':
			result += `\\`
		case '\n':
			result += `\n`
		case '\r':
			result += `\r`
		case '\t':
			result += `\t`
		default:
			result += string(c)
		}
	}
	return result
}

// Logger is the main logging interface
type Logger struct {
	level     LogLevel
	formatter LogFormatter
	output    io.Writer
	component string
	fields    map[string]interface{}
	mu        sync.Mutex
}

// NewLogger creates a new logger with default settings
func NewLogger() *Logger {
	return &Logger{
		level:     LogLevelInfo,
		formatter: NewTextFormatter(),
		output:    os.Stderr,
		fields:    make(map[string]interface{}),
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetFormatter sets the log formatter
func (l *Logger) SetFormatter(formatter LogFormatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = formatter
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// SetComponent sets the component name for all log entries
func (l *Logger) SetComponent(component string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.component = component
}

// WithField creates a new logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:     l.level,
		formatter: l.formatter,
		output:    l.output,
		component: l.component,
		fields:    make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value

	return newLogger
}

// WithFields creates a new logger with multiple additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:     l.level,
		formatter: l.formatter,
		output:    l.output,
		component: l.component,
		fields:    make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// log writes a log entry at the specified level
func (l *Logger) log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Skip if below minimum level
	if level > l.level {
		return
	}

	// Create log entry
	entry := &LogEntry{
		Time:      time.Now(),
		Level:     level,
		Message:   message,
		Fields:    l.fields,
		Component: l.component,
	}

	// Format and write
	formatted := l.formatter.Format(entry)
	fmt.Fprintln(l.output, formatted)

	// Flush if we're writing to a file
	if f, ok := l.output.(*os.File); ok {
		f.Sync()
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string) {
	l.log(LogLevelFatal, message)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log(LogLevelError, message)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	l.log(LogLevelWarn, message)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Info logs an informational message
func (l *Logger) Info(message string) {
	l.log(LogLevelInfo, message)
}

// Infof logs a formatted informational message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.log(LogLevelDebug, message)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Trace logs a trace message
func (l *Logger) Trace(message string) {
	l.log(LogLevelTrace, message)
}

// Tracef logs a formatted trace message
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Trace(fmt.Sprintf(format, args...))
}

// Global logger instance
var globalLogger = NewLogger()

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return globalLogger
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// Package-level convenience functions that use the global logger

// Fatal logs a fatal message using the global logger
func Fatal(message string) {
	globalLogger.Fatal(message)
}

// Fatalf logs a formatted fatal message using the global logger
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// Error logs an error message using the global logger
func Error(message string) {
	globalLogger.Error(message)
}

// Errorf logs a formatted error message using the global logger
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Warn logs a warning message using the global logger
func Warn(message string) {
	globalLogger.Warn(message)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Info logs an informational message using the global logger
func Info(message string) {
	globalLogger.Info(message)
}

// Infof logs a formatted informational message using the global logger
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Debug logs a debug message using the global logger
func Debug(message string) {
	globalLogger.Debug(message)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Trace logs a trace message using the global logger
func Trace(message string) {
	globalLogger.Trace(message)
}

// Tracef logs a formatted trace message using the global logger
func Tracef(format string, args ...interface{}) {
	globalLogger.Tracef(format, args...)
}
