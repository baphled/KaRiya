package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel defines the severity of log messages
type LogLevel int

const (
	// DebugLevel is the lowest logging level, used for detailed debugging information
	DebugLevel LogLevel = iota
	// InfoLevel is used for general information about system operations
	InfoLevel
	// WarnLevel is used for potentially harmful situations
	WarnLevel
	// ErrorLevel is used for error events that might still allow the application to continue running
	ErrorLevel
	// FatalLevel is used for severe errors that cause the application to terminate
	FatalLevel
)

// Logger provides a structured and configurable logging mechanism
type Logger struct {
	logger     *log.Logger
	level      LogLevel
	output     io.Writer
	mu         sync.Mutex
	contextMap map[string]string
}

// New creates a new Logger with default configuration
func New(output io.Writer, level LogLevel) *Logger {
	return &Logger{
		logger:     log.New(output, "", log.LstdFlags|log.Lmicroseconds),
		level:      level,
		output:     output,
		contextMap: make(map[string]string),
	}
}

// DefaultLogger creates a standard logger writing to stdout
func DefaultLogger() *Logger {
	return New(os.Stdout, InfoLevel)
}

// SetContext adds a key-value pair to the logger's context
func (l *Logger) SetContext(key, value string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.contextMap[key] = value
	return l
}

// ClearContext removes all context from the logger
func (l *Logger) ClearContext() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.contextMap = make(map[string]string)
	return l
}

// formatContext converts the context map to a formatted string
func (l *Logger) formatContext() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	var contextParts []string
	for k, v := range l.contextMap {
		contextParts = append(contextParts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(contextParts, " ")
}

// getCallerInfo retrieves the file, line, and function name of the caller
func getCallerInfo(skip int) (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}
	function = runtime.FuncForPC(pc).Name()
	return file, line, function
}

// log writes a log message with the specified level and format
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	file, line, function := getCallerInfo(2)

	// Construct log message
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf(
		"%s [%s] %s:%d (%s) %s %s",
		level.String(),
		time.Now().Format(time.RFC3339),
		file,
		line,
		function,
		l.formatContext(),
		message,
	)

	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Println(logEntry)
}

// Debug logs a message at Debug level
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs a message at Info level
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a message at Warn level
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs a message at Error level
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a message at Fatal level and terminates the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
	os.Exit(1)
}

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// WithFields creates a new logger with additional context
func (l *Logger) WithFields(fields map[string]string) *Logger {
	newLogger := &Logger{
		logger:     l.logger,
		level:      l.level,
		output:     l.output,
		contextMap: make(map[string]string),
	}

	// Copy existing context
	for k, v := range l.contextMap {
		newLogger.contextMap[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.contextMap[k] = v
	}

	return newLogger
}

// FileLogger creates a logger that writes to a file in the user's home directory
func FileLogger() *Logger {
	// Create logs directory in home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to stdout if we can't get home directory
		return DefaultLogger()
	}

	logDir := homeDir + "/.kariya/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Fallback to stdout if we can't create directory
		return DefaultLogger()
	}

	// Create log file with timestamp
	logFile := logDir + "/kariya.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback to stdout if we can't open file
		return DefaultLogger()
	}

	return New(file, InfoLevel)
}
