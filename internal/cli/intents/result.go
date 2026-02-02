package intents

import (
	"errors"
	"fmt"
)

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
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface, producing a formatted diagnostic string
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (e *IntentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithCause attaches the root-cause error to the IntentError for debugging.
//
// Expected:
//   - cause is the originating lower-level error, or nil to clear the cause.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Cause field on the receiver.
func (e *IntentError) WithCause(cause error) *IntentError {
	e.Cause = cause
	return e
}

// WithMessage attaches a human-readable description that can be surfaced to the user or written to logs.
//
// Expected:
//   - message is a non-empty string describing what happened in plain language.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Message field on the receiver.
func (e *IntentError) WithMessage(message string) *IntentError {
	e.Message = message
	return e
}

// IntentResult wraps the outcome of an intent execution into a type-safe
// envelope that the intent router uses to determine what happened.
//
// The generic type parameter T specifies the concrete payload type the
// intent produces on success (e.g., *captureevent.Result or
// *GenerateCVResult). T must satisfy the "any" constraint.
//
// The Status field is one of the following ResultStatus values:
//
//   - Completed: the intent finished successfully. The Data field of type T
//     holds the intent output. The Error field is nil.
//   - Partial: the intent succeeded for some inputs but not all. The Data
//     field holds accepted output and the Error field contains an
//     IntentError describing what was rejected.
//   - Failed: the intent encountered an unrecoverable error. The Error
//     field contains an IntentError with a machine-readable Code, a
//     human-readable Message, and an optional Cause. The Data field is
//     zero-valued.
//   - Cancelled: the user explicitly aborted the intent. Both Data and
//     Error are zero-valued.
//
// The Metadata field is a map[string]interface{} carrying auxiliary
// key-value pairs such as breadcrumbs, timestamps, or diagnostic hints.
//
// Every intent must return an IntentResult as the sole communication
// channel back to the intent router; direct state mutation across intent
// boundaries is forbidden.
type IntentResult[T any] struct {
	Status   ResultStatus
	Data     T
	Error    *IntentError
	Metadata map[string]interface{}
}

// NewCompletedResult creates a successful IntentResult with output data.
//
// Expected:
//   - data is the intent output payload.
//
// Returns:
//   - A new IntentResult with Completed status and initialized metadata.
//
// Side effects:
//   - None.
func NewCompletedResult[T any](data T) *IntentResult[T] {
	return &IntentResult[T]{
		Status:   Completed,
		Data:     data,
		Metadata: make(map[string]interface{}),
	}
}

// NewCancelledResult creates an IntentResult indicating user cancellation.
//
// Returns:
//   - A new IntentResult with Cancelled status and initialized metadata.
//
// Side effects:
//   - None.
func NewCancelledResult[T any]() *IntentResult[T] {
	return &IntentResult[T]{
		Status:   Cancelled,
		Metadata: make(map[string]interface{}),
	}
}

// NewFailedResult creates an IntentResult indicating an error.
//
// Expected:
//   - code is a machine-readable error code.
//   - message is a human-readable error description.
//   - cause is the underlying error, or nil if none.
//
// Returns:
//   - A new IntentResult with Failed status, an IntentError, and initialized metadata.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - data is the accepted output payload.
//   - code is a machine-readable error code for the rejected portion.
//   - message is a human-readable description of what was rejected.
//
// Returns:
//   - A new IntentResult with Partial status, data, an IntentError, and initialized metadata.
//
// Side effects:
//   - None.
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

// IsSuccessful indicates the intent produced usable output, either fully (Completed) or
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) IsSuccessful() bool {
	return r.Status == Completed || r.Status == Partial
}

// IsCancelled indicates the user explicitly aborted the intent before it could
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) IsCancelled() bool {
	return r.Status == Cancelled
}

// IsFailed indicates the intent encountered an unrecoverable error and produced
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) IsFailed() bool {
	return r.Status == Failed
}

// IsTerminal indicates the result has reached a final state and no further processing
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) IsTerminal() bool {
	return r.IsCancelled() || r.IsFailed() || r.IsSuccessful()
}

// WithMetadata attaches an auxiliary key-value pair (breadcrumbs, timestamps, diagnostic
// hints) to the result for downstream consumers that need context beyond status and data.
//
// Expected:
//   - key is a non-empty, dot-namespaced identifier (e.g. "timing.duration").
//   - value is any serialisable payload; callers must agree on the concrete type per key.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Metadata map on the receiver, initializing it if nil.
func (r *IntentResult[T]) WithMetadata(key string, value interface{}) *IntentResult[T] {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

// GetMetadata looks up a single auxiliary value (breadcrumb, timestamp, diagnostic hint)
// previously attached via WithMetadata.
//
// Expected:
//   - key is the exact identifier used when the metadata was stored.
//
// Returns:
//   - The value associated with the key and true, or nil and false if not found.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) GetMetadata(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	val, ok := r.Metadata[key]
	return val, ok
}

// GetAllMetadata produces a shallow copy of every auxiliary key-value pair attached to the
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) GetAllMetadata() map[string]interface{} {
	if r.Metadata == nil {
		return make(map[string]interface{})
	}
	result := make(map[string]interface{})
	for k, v := range r.Metadata {
		result[k] = v
	}
	return result
}

// WithError attaches structured error information to the result for diagnostics
// and downstream error handling.
//
// Expected:
//   - err is a fully populated IntentError with Code and Message, or nil to clear.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Error field on the receiver.
func (r *IntentResult[T]) WithError(err *IntentError) *IntentResult[T] {
	r.Error = err
	return r
}

// WithStatus overrides the result's outcome classification, useful when a result
// must be reclassified after construction (e.g., downgrading Completed to Partial).
//
// Expected:
//   - status is one of Completed, Cancelled, Failed, or Partial.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Status field on the receiver.
func (r *IntentResult[T]) WithStatus(status ResultStatus) *IntentResult[T] {
	r.Status = status
	return r
}

// WithData attaches or replaces the intent's output payload, allowing post-construction
// enrichment of the result.
//
// Expected:
//   - data is a fully populated payload of type T appropriate for the current status.
//
// Returns:
//   - The receiver for method chaining.
//
// Side effects:
//   - Mutates the Data field on the receiver.
func (r *IntentResult[T]) WithData(data T) *IntentResult[T] {
	r.Data = data
	return r
}

// IsValid enforces the structural invariants each status requires: Completed results
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *IntentResult[T]) IsValid() error {
	switch r.Status {
	case Completed:
		return nil

	case Cancelled:
		if r.Error != nil {
			return errors.New("cancelled result should not have error")
		}
		return nil

	case Failed:
		if r.Error == nil {
			return errors.New("failed result must have error")
		}
		return nil

	case Partial:
		if r.Error == nil {
			return errors.New("partial result must have error")
		}
		return nil

	default:
		return fmt.Errorf("unknown result status: %s", r.Status)
	}
}
