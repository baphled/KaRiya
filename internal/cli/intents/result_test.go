package intents

import (
	"errors"
	"testing"
)

func TestIntentResult_NewCompletedResult(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "creates completed result with data",
			data: "test data",
		},
		{
			name: "creates completed result with empty string",
			data: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewCompletedResult(tt.data)

			if result.Status != Completed {
				t.Errorf("expected Completed status, got %s", result.Status)
			}

			if result.Data != tt.data {
				t.Errorf("expected data %q, got %q", tt.data, result.Data)
			}

			if result.Error != nil {
				t.Errorf("expected no error, got %v", result.Error)
			}

			if !result.IsSuccessful() {
				t.Errorf("expected IsSuccessful() to return true")
			}
		})
	}
}

func TestIntentResult_NewCancelledResult(t *testing.T) {
	result := NewCancelledResult[string]()

	if result.Status != Cancelled {
		t.Errorf("expected Cancelled status, got %s", result.Status)
	}

	if result.IsCancelled() == false {
		t.Errorf("expected IsCancelled() to return true")
	}

	if result.IsSuccessful() {
		t.Errorf("expected IsSuccessful() to return false")
	}
}

func TestIntentResult_NewFailedResult(t *testing.T) {
	cause := errors.New("test error")
	result := NewFailedResult[string]("test_code", "test message", cause)

	if result.Status != Failed {
		t.Errorf("expected Failed status, got %s", result.Status)
	}

	if result.IsFailed() == false {
		t.Errorf("expected IsFailed() to return true")
	}

	if result.Error == nil {
		t.Errorf("expected error to be set")
	}

	if result.Error.Code != "test_code" {
		t.Errorf("expected error code test_code, got %s", result.Error.Code)
	}

	if result.Error.Cause != cause {
		t.Errorf("expected error cause to be set")
	}

	if result.IsSuccessful() {
		t.Errorf("expected IsSuccessful() to return false")
	}
}

func TestIntentResult_NewPartialResult(t *testing.T) {
	data := "partial data"
	result := NewPartialResult(data, "partial_code", "partial message")

	if result.Status != Partial {
		t.Errorf("expected Partial status, got %s", result.Status)
	}

	if result.Data != data {
		t.Errorf("expected data %q, got %q", data, result.Data)
	}

	if result.Error == nil {
		t.Errorf("expected error to be set for partial result")
	}

	if result.IsSuccessful() == false {
		t.Errorf("expected IsSuccessful() to return true for partial result")
	}
}

func TestIntentResult_WithMetadata(t *testing.T) {
	result := NewCompletedResult("test")

	result.WithMetadata("key1", "value1")
	result.WithMetadata("key2", 42)

	val, ok := result.GetMetadata("key1")
	if !ok || val != "value1" {
		t.Errorf("expected metadata key1 to be value1, got %v", val)
	}

	val, ok = result.GetMetadata("key2")
	if !ok || val != 42 {
		t.Errorf("expected metadata key2 to be 42, got %v", val)
	}

	_, ok = result.GetMetadata("nonexistent")
	if ok {
		t.Errorf("expected nonexistent metadata to return false")
	}
}

func TestIntentResult_GetAllMetadata(t *testing.T) {
	result := NewCompletedResult("test")
	result.WithMetadata("key1", "value1")
	result.WithMetadata("key2", 42)

	all := result.GetAllMetadata()

	if len(all) != 2 {
		t.Errorf("expected 2 metadata entries, got %d", len(all))
	}

	if all["key1"] != "value1" {
		t.Errorf("expected key1 to be value1")
	}

	if all["key2"] != 42 {
		t.Errorf("expected key2 to be 42")
	}
}

func TestIntentResult_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		result   *IntentResult[string]
		expected bool
	}{
		{
			name:     "completed is terminal",
			result:   NewCompletedResult("data"),
			expected: true,
		},
		{
			name:     "cancelled is terminal",
			result:   NewCancelledResult[string](),
			expected: true,
		},
		{
			name:     "failed is terminal",
			result:   NewFailedResult[string]("code", "message", nil),
			expected: true,
		},
		{
			name:     "partial is terminal",
			result:   NewPartialResult("data", "code", "message"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.IsTerminal() != tt.expected {
				t.Errorf("expected IsTerminal() to return %v", tt.expected)
			}
		})
	}
}

func TestIntentResult_MethodChaining(t *testing.T) {
	result := NewCompletedResult("data").
		WithMetadata("key1", "value1").
		WithMetadata("key2", 42)

	val, ok := result.GetMetadata("key1")
	if !ok || val != "value1" {
		t.Errorf("expected chained metadata to be set")
	}
}

func TestIntentResult_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		result      *IntentResult[string]
		shouldError bool
	}{
		{
			name:        "completed result is valid",
			result:      NewCompletedResult("data"),
			shouldError: false,
		},
		{
			name:        "cancelled result is valid",
			result:      NewCancelledResult[string](),
			shouldError: false,
		},
		{
			name:        "failed result without error is invalid",
			result:      &IntentResult[string]{Status: Failed},
			shouldError: true,
		},
		{
			name:        "failed result with error is valid",
			result:      NewFailedResult[string]("code", "message", nil),
			shouldError: false,
		},
		{
			name:        "partial result without error is invalid",
			result:      &IntentResult[string]{Status: Partial, Data: "data"},
			shouldError: true,
		},
		{
			name:        "partial result with error is valid",
			result:      NewPartialResult("data", "code", "message"),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.IsValid()
			if (err != nil) != tt.shouldError {
				t.Errorf("expected error=%v, got error=%v", tt.shouldError, err != nil)
			}
		})
	}
}

func TestIntentError_WithCause(t *testing.T) {
	err := &IntentError{
		Code:    "code1",
		Message: "message1",
	}

	cause := errors.New("cause error")
	err.WithCause(cause) // nolint: errcheck

	if err.Cause != cause {
		t.Errorf("expected cause to be set")
	}
}

func TestIntentError_WithMessage(t *testing.T) {
	err := &IntentError{
		Code:    "code1",
		Message: "original message",
	}

	err.WithMessage("updated message") // nolint: errcheck

	if err.Message != "updated message" {
		t.Errorf("expected message to be updated")
	}
}

func TestIntentError_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "error with cause",
			code:     "code1",
			message:  "message1",
			cause:    errors.New("cause error"),
			expected: "code1: message1 (cause: cause error)",
		},
		{
			name:     "error without cause",
			code:     "code2",
			message:  "message2",
			cause:    nil,
			expected: "code2: message2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &IntentError{
				Code:    tt.code,
				Message: tt.message,
				Cause:   tt.cause,
			}

			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}
		})
	}
}
