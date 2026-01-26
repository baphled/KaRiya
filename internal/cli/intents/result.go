package intents

import "fmt"

// ResultStatus represents the outcome of an intent execution.
type ResultStatus string

const (
	// Completed indicates the intent finished successfully.
	Completed ResultStatus = "completed"

	// Cancelled indicates the user explicitly cancelled the intent.
	Cancelled ResultStatus = "cancelled"

	// Failed indicates the intent encountered an error.
	Failed ResultStatus = "failed"

	// Partial indicates the intent succeeded partially (e.g., some data accepted, some rejected).
	Partial ResultStatus = "partial"
)

// IntentError provides debug and logging information without violating type safety.
type IntentError struct {
	Code    string // Machine-readable error code
	Message string // Human-readable error message
	Cause   error  // Underlying error for debugging
}

func (e *IntentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithCause adds or updates the cause error.
func (e *IntentError) WithCause(cause error) *IntentError {
	e.Cause = cause
	return e
}

// WithMessage updates the human-readable message.
func (e *IntentError) WithMessage(message string) *IntentError {
	e.Message = message
	return e
}

// IntentResult[T] is the type-safe boundary contract for intent communication.
// Each intent MUST return a strongly-typed IntentResult.
// Results are the ONLY mechanism for intent communication.
type IntentResult[T any] struct {
	// Status indicates how the intent completed.
	Status ResultStatus

	// Data contains the intent's output, only valid when Status is Completed or Partial.
	Data T

	// Error provides debug/logging information without violating type safety.
	Error *IntentError

	// Metadata is for additional context (e.g., breadcrumbs, timestamps).
	Metadata map[string]interface{}
}

// NewCompletedResult creates a successful IntentResult with output data.
func NewCompletedResult[T any](data T) *IntentResult[T] {
	return &IntentResult[T]{
		Status:   Completed,
		Data:     data,
		Metadata: make(map[string]interface{}),
	}
}

// NewCancelledResult creates an IntentResult indicating user cancellation.
func NewCancelledResult[T any]() *IntentResult[T] {
	return &IntentResult[T]{
		Status:   Cancelled,
		Metadata: make(map[string]interface{}),
	}
}

// NewFailedResult creates an IntentResult indicating an error.
func NewFailedResult[T any](code, message string, cause error) *IntentResult[T] {
	return &IntentResult[T]{
		Status: Failed,
		Error: &IntentError{
			Code:    code,
			Message: message,
			Cause:   cause,
		},
		Metadata: make(map[string]interface{}),
	}
}

// NewPartialResult creates an IntentResult indicating partial success.
// Used for flows where some data is accepted and some is rejected (e.g., CaptureEvent).
func NewPartialResult[T any](data T, code, message string) *IntentResult[T] {
	return &IntentResult[T]{
		Status: Partial,
		Data:   data,
		Error: &IntentError{
			Code:    code,
			Message: message,
		},
		Metadata: make(map[string]interface{}),
	}
}

// IsSuccessful returns true if the intent completed successfully (Completed or Partial).
func (r *IntentResult[T]) IsSuccessful() bool {
	return r.Status == Completed || r.Status == Partial
}

// IsCancelled returns true if the user cancelled the intent.
func (r *IntentResult[T]) IsCancelled() bool {
	return r.Status == Cancelled
}

// IsFailed returns true if the intent encountered an error.
func (r *IntentResult[T]) IsFailed() bool {
	return r.Status == Failed
}

// IsTerminal returns true if the result is in a terminal state (cannot transition further).
func (r *IntentResult[T]) IsTerminal() bool {
	return r.IsCancelled() || r.IsFailed() || r.IsSuccessful()
}

// WithMetadata adds or updates metadata on the result.
// Returns the result for method chaining.
func (r *IntentResult[T]) WithMetadata(key string, value interface{}) *IntentResult[T] {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

// GetMetadata retrieves metadata by key.
// Returns the value and a boolean indicating if the key exists.
func (r *IntentResult[T]) GetMetadata(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	val, ok := r.Metadata[key]
	return val, ok
}

// GetAllMetadata returns a copy of all metadata.
func (r *IntentResult[T]) GetAllMetadata() map[string]interface{} {
	if r.Metadata == nil {
		return make(map[string]interface{})
	}
	// Create a copy to prevent external mutation.
	result := make(map[string]interface{})
	for k, v := range r.Metadata {
		result[k] = v
	}
	return result
}

// WithError sets the error on the result.
// Returns the result for method chaining.
func (r *IntentResult[T]) WithError(err *IntentError) *IntentResult[T] {
	r.Error = err
	return r
}

// WithStatus sets the status on the result.
// Returns the result for method chaining.
func (r *IntentResult[T]) WithStatus(status ResultStatus) *IntentResult[T] {
	r.Status = status
	return r
}

// WithData sets the data on the result.
// Returns the result for method chaining.
func (r *IntentResult[T]) WithData(data T) *IntentResult[T] {
	r.Data = data
	return r
}

// IsValid checks if the result is in a valid state.
// A result is valid if:
// - Completed results have data
// - Failed results have an error
// - Cancelled results have no data or error
// - Partial results have data and an error
func (r *IntentResult[T]) IsValid() error {
	switch r.Status {
	case Completed:
		// Completed results should have data (but we can't check if T is zero value)
		// So we just verify the status is set
		return nil

	case Cancelled:
		// Cancelled results should not have error
		if r.Error != nil {
			return fmt.Errorf("cancelled result should not have error")
		}
		return nil

	case Failed:
		// Failed results must have error
		if r.Error == nil {
			return fmt.Errorf("failed result must have error")
		}
		return nil

	case Partial:
		// Partial results must have error
		if r.Error == nil {
			return fmt.Errorf("partial result must have error")
		}
		return nil

	default:
		return fmt.Errorf("unknown result status: %s", r.Status)
	}
}
