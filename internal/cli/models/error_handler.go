package models

import (
	"fmt"
	"log"
	"os"
	"time"
)

// ErrorSeverity defines the severity level of an error
type ErrorSeverity int

const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// String returns the string representation of ErrorSeverity
func (es ErrorSeverity) String() string {
	switch es {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ErrorLog represents a logged error with metadata
type ErrorLog struct {
	Severity    ErrorSeverity
	Message     string
	Error       error
	Timestamp   time.Time
	Context     map[string]interface{}
	Suggestion  string
	ScreenID    string
}

// ErrorHandler manages error logging and recovery suggestions
type ErrorHandler struct {
	errors      []ErrorLog
	logger      *log.Logger
	maxErrors   int
	suggestions map[string]string
}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		errors:      make([]ErrorLog, 0),
		logger:      log.New(os.Stderr, "[KaRiya] ", log.LstdFlags),
		maxErrors:   100,
		suggestions: make(map[string]string),
	}
}

// LogError logs an error with the specified severity
func (eh *ErrorHandler) LogError(severity ErrorSeverity, message string, err error, screenID string) {
	errorLog := ErrorLog{
		Severity:   severity,
		Message:    message,
		Error:      err,
		Timestamp:  time.Now(),
		Context:    make(map[string]interface{}),
		ScreenID:   screenID,
	}

	// Add suggestion if available
	if err != nil {
		if suggestion, exists := eh.suggestions[err.Error()]; exists {
			errorLog.Suggestion = suggestion
		}
	}

	eh.errors = append(eh.errors, errorLog)

	// Keep only the latest maxErrors
	if len(eh.errors) > eh.maxErrors {
		eh.errors = eh.errors[1:]
	}

	// Log to file/stderr
	eh.logToOutput(errorLog)
}

// logToOutput writes error information to the output
func (eh *ErrorHandler) logToOutput(errorLog ErrorLog) {
	logMessage := fmt.Sprintf("%s: %s", errorLog.Severity, errorLog.Message)
	if errorLog.Error != nil {
		logMessage += fmt.Sprintf(" - %v", errorLog.Error)
	}
	eh.logger.Println(logMessage)
}

// GetLastError returns the most recent error
func (eh *ErrorHandler) GetLastError() *ErrorLog {
	if len(eh.errors) == 0 {
		return nil
	}
	return &eh.errors[len(eh.errors)-1]
}

// GetErrorsByScreenID returns all errors for a specific screen
func (eh *ErrorHandler) GetErrorsByScreenID(screenID string) []ErrorLog {
	var results []ErrorLog
	for _, errLog := range eh.errors {
		if errLog.ScreenID == screenID {
			results = append(results, errLog)
		}
	}
	return results
}

// GetErrorsBySeverity returns all errors with the specified severity
func (eh *ErrorHandler) GetErrorsBySeverity(severity ErrorSeverity) []ErrorLog {
	var results []ErrorLog
	for _, errLog := range eh.errors {
		if errLog.Severity == severity {
			results = append(results, errLog)
		}
	}
	return results
}

// ClearErrors removes all logged errors
func (eh *ErrorHandler) ClearErrors() {
	eh.errors = make([]ErrorLog, 0)
}

// ClearErrorsByScreenID removes errors for a specific screen
func (eh *ErrorHandler) ClearErrorsByScreenID(screenID string) {
	filtered := make([]ErrorLog, 0)
	for _, errLog := range eh.errors {
		if errLog.ScreenID != screenID {
			filtered = append(filtered, errLog)
		}
	}
	eh.errors = filtered
}

// RegisterSuggestion registers a recovery suggestion for a specific error
func (eh *ErrorHandler) RegisterSuggestion(errorMessage string, suggestion string) {
	eh.suggestions[errorMessage] = suggestion
}

// GetSuggestion retrieves the recovery suggestion for an error
func (eh *ErrorHandler) GetSuggestion(err error) string {
	if err == nil {
		return ""
	}
	if suggestion, exists := eh.suggestions[err.Error()]; exists {
		return suggestion
	}
	return ""
}

// GetErrorHistory returns all logged errors
func (eh *ErrorHandler) GetErrorHistory() []ErrorLog {
	return eh.errors
}

// FormatErrorMessage returns a user-friendly error message
func FormatErrorMessage(severity ErrorSeverity, message string, suggestion string) string {
	formatted := fmt.Sprintf("[%s] %s", severity, message)
	if suggestion != "" {
		formatted += fmt.Sprintf("\n💡 Suggestion: %s", suggestion)
	}
	return formatted
}
