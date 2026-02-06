# Error Handling Skill

You are an expert in Go error handling patterns and KaRiya's error conventions.

## Overview

KaRiya uses standard Go error patterns: sentinel errors, error wrapping, and structured validation errors.

---

## Sentinel Errors

### Definition

```go
// internal/repository/career/errors.go
package career

import "errors"

// Repository sentinel errors.
var (
    ErrEventNotFound   = errors.New("event not found")
    ErrDuplicateEvent  = errors.New("duplicate event")
    ErrInvalidEventID  = errors.New("invalid event ID")
)
```

```go
// internal/service/career/errors.go
package career

import "errors"

// Service sentinel errors.
var (
    ErrEventNotFound              = errors.New("event not found")
    ErrFactRepositoryNotConfigured = errors.New("fact repository not configured")
    ErrInvalidEventData           = errors.New("invalid event data")
)
```

### Usage

```go
// Returning sentinel error
func (r *Repository) FindByID(ctx context.Context, id string) (*Event, error) {
    var event Event
    result := r.db.First(&event, "id = ?", id)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrEventNotFound
    }
    return &event, result.Error
}

// Checking sentinel error
event, err := repo.FindByID(ctx, id)
if errors.Is(err, ErrEventNotFound) {
    // Handle not found case
    return nil, ErrEventNotFound
}
```

---

## Error Wrapping

### Adding Context

```go
import "fmt"

func (s *Service) CreateEvent(ctx context.Context, data EventData) (*Event, error) {
    event, err := career.NewEvent(data.Text, data.Date)
    if err != nil {
        return nil, fmt.Errorf("invalid event data: %w", err)
    }
    
    if err := s.repo.Save(ctx, event); err != nil {
        return nil, fmt.Errorf("failed to save event: %w", err)
    }
    
    return event, nil
}
```

### Preserving Error Chain

```go
// Error chain: service -> repository -> gorm
// "failed to create event: failed to save: UNIQUE constraint failed"

func processEvent(ctx context.Context, data EventData) error {
    _, err := service.CreateEvent(ctx, data)
    if err != nil {
        return fmt.Errorf("failed to process event: %w", err)
    }
    return nil
}
```

### Unwrapping Errors

```go
// Check for specific error in chain
if errors.Is(err, ErrEventNotFound) {
    // Handle not found
}

// Get underlying error type
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error specifically
    fmt.Println("Validation failed:", validationErr.Field)
}
```

---

## Validation Errors

### Domain Validation

```go
// internal/domain/career/event.go
func (e *Event) Validate() error {
    if e.Text == "" {
        return errors.New("event text is required")
    }
    if len(e.Text) > 2000 {
        return errors.New("event text exceeds maximum length of 2000 characters")
    }
    if e.Date.IsZero() {
        return errors.New("event date is required")
    }
    if e.Date.After(time.Now().Add(24 * time.Hour)) {
        return errors.New("event date cannot be in the future")
    }
    return nil
}
```

### Structured Validation Error

```go
// ValidationError contains details about validation failures.
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e *Event) Validate() error {
    if e.Text == "" {
        return &ValidationError{
            Field:   "text",
            Message: "is required",
        }
    }
    return nil
}

// Usage
err := event.Validate()
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    fmt.Printf("Field %s %s\n", validationErr.Field, validationErr.Message)
}
```

### Multiple Validation Errors

```go
// MultiValidationError collects multiple validation failures.
type MultiValidationError struct {
    Errors []ValidationError
}

func (e *MultiValidationError) Error() string {
    var msgs []string
    for _, err := range e.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (e *Event) ValidateAll() error {
    var errs []ValidationError
    
    if e.Text == "" {
        errs = append(errs, ValidationError{Field: "text", Message: "is required"})
    }
    if e.Date.IsZero() {
        errs = append(errs, ValidationError{Field: "date", Message: "is required"})
    }
    
    if len(errs) > 0 {
        return &MultiValidationError{Errors: errs}
    }
    return nil
}
```

---

## Error Handling in Layers

### Repository Layer

```go
func (r *Repository) Save(ctx context.Context, event *Event) error {
    result := r.db.WithContext(ctx).Save(event)
    if result.Error != nil {
        // Translate DB errors to domain errors
        if isDuplicateKeyError(result.Error) {
            return ErrDuplicateEvent
        }
        return fmt.Errorf("database error: %w", result.Error)
    }
    return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Event, error) {
    var event Event
    result := r.db.WithContext(ctx).First(&event, "id = ?", id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, ErrEventNotFound
        }
        return nil, fmt.Errorf("database error: %w", result.Error)
    }
    return &event, nil
}
```

### Service Layer

```go
func (s *Service) GetEvent(ctx context.Context, id string) (*Event, error) {
    if id == "" {
        return nil, errors.New("event ID is required")
    }
    
    event, err := s.repo.FindByID(ctx, id)
    if err != nil {
        // Translate repository errors to service errors
        if errors.Is(err, repo.ErrEventNotFound) {
            return nil, ErrEventNotFound  // Service-level error
        }
        // Log unexpected errors
        s.logger.Error("failed to find event", "error", err)
        return nil, fmt.Errorf("failed to get event: %w", err)
    }
    
    return event, nil
}
```

### Intent/UI Layer

```go
func (i *Intent) handleSaveResult(msg SaveResultMsg) tea.Cmd {
    if msg.Error != nil {
        // Translate to user-friendly message
        var userMsg string
        
        switch {
        case errors.Is(msg.Error, service.ErrEventNotFound):
            userMsg = "The event could not be found. It may have been deleted."
        case errors.Is(msg.Error, service.ErrDuplicateEvent):
            userMsg = "An event with this information already exists."
        default:
            userMsg = "An unexpected error occurred. Please try again."
            // Log the actual error
            i.logger.Error("save failed", "error", msg.Error)
        }
        
        i.showErrorModal(userMsg)
        return nil
    }
    
    return i.handleSuccess(msg.Event)
}
```

---

## Error Handling in Async Operations

### Message Pattern

```go
// Define result message with error field
type LoadEventsMsg struct {
    Events []*Event
    Error  error
}

// Create command
func (i *Intent) loadEventsAsync() tea.Cmd {
    return func() tea.Msg {
        events, err := i.service.ListEvents(i.ctx)
        return LoadEventsMsg{
            Events: events,
            Error:  err,
        }
    }
}

// Handle in Update
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case LoadEventsMsg:
        if msg.Error != nil {
            return i.handleLoadError(msg.Error)
        }
        i.events = msg.Events
        return nil
    }
    return nil
}
```

---

## Error Display

### Error Modal Pattern

```go
func (i *Intent) showErrorModal(message string) {
    i.errorModal = feedback.NewErrorModal(message)
    i.errorModal.Show()
}

func (i *Intent) handleError(err error) tea.Cmd {
    // Map errors to user messages
    userMsg := mapErrorToUserMessage(err)
    i.showErrorModal(userMsg)
    return nil
}

func mapErrorToUserMessage(err error) string {
    switch {
    case errors.Is(err, ErrEventNotFound):
        return "Event not found"
    case errors.Is(err, ErrInvalidEventData):
        return "Please check your input and try again"
    case errors.Is(err, context.DeadlineExceeded):
        return "Operation timed out. Please try again."
    case errors.Is(err, context.Canceled):
        return "Operation was cancelled"
    default:
        return "An unexpected error occurred"
    }
}
```

---

## Logging Errors

```go
import "github.com/baphled/kariya/internal/logger"

func (s *Service) Operation(ctx context.Context) error {
    result, err := s.repo.Find(ctx)
    if err != nil {
        // Log with context
        s.logger.WithFields(map[string]string{
            "operation": "find",
            "error":     err.Error(),
        }).Error("Repository operation failed")
        
        return fmt.Errorf("failed to find: %w", err)
    }
    return nil
}
```

---

## Common Patterns

### Early Return

```go
func (s *Service) Process(ctx context.Context, data Data) error {
    // Validate early
    if err := data.Validate(); err != nil {
        return fmt.Errorf("invalid data: %w", err)
    }
    
    // Check preconditions
    if s.repo == nil {
        return ErrRepositoryNotConfigured
    }
    
    // Main operation
    if err := s.repo.Save(ctx, data); err != nil {
        return fmt.Errorf("failed to save: %w", err)
    }
    
    return nil
}
```

### Defer with Error

```go
func (r *Repository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
    tx := r.db.Begin()
    if tx.Error != nil {
        return fmt.Errorf("failed to begin transaction: %w", tx.Error)
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // Re-panic after rollback
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit: %w", err)
    }
    
    return nil
}
```

---

## Testing Errors

```go
var _ = Describe("EventService", func() {
    Describe("GetEvent", func() {
        Context("when event not found", func() {
            BeforeEach(func() {
                repo.EXPECT().
                    FindByID(ctx, "not-found").
                    Return(nil, repo.ErrEventNotFound)
            })
            
            It("returns service not found error", func() {
                _, err := service.GetEvent(ctx, "not-found")
                
                Expect(err).To(MatchError(service.ErrEventNotFound))
            })
        })
        
        Context("when repository fails", func() {
            BeforeEach(func() {
                repo.EXPECT().
                    FindByID(ctx, gomock.Any()).
                    Return(nil, errors.New("database error"))
            })
            
            It("wraps the error", func() {
                _, err := service.GetEvent(ctx, "test-id")
                
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("failed to get event"))
                Expect(err.Error()).To(ContainSubstring("database error"))
            })
        })
    })
})
```

---

## Anti-Patterns

### DON'T: Ignore Errors

```go
// WRONG
result, _ := service.GetEvent(ctx, id)

// CORRECT
result, err := service.GetEvent(ctx, id)
if err != nil {
    return err
}
```

### DON'T: Panic on Recoverable Errors

```go
// WRONG
func (s *Service) GetEvent(ctx context.Context, id string) *Event {
    event, err := s.repo.FindByID(ctx, id)
    if err != nil {
        panic(err)  // Never panic for business errors!
    }
    return event
}

// CORRECT
func (s *Service) GetEvent(ctx context.Context, id string) (*Event, error) {
    return s.repo.FindByID(ctx, id)
}
```

### DON'T: Lose Error Context

```go
// WRONG - Original error lost
if err != nil {
    return errors.New("operation failed")
}

// CORRECT - Wrap to preserve chain
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### DON'T: Double Wrap

```go
// WRONG - Redundant wrapping
if err != nil {
    return fmt.Errorf("failed: %w", fmt.Errorf("error: %w", err))
}

// CORRECT - Single wrap with context
if err != nil {
    return fmt.Errorf("failed to process event: %w", err)
}
```

### DON'T: Use String Comparison

```go
// WRONG
if err.Error() == "event not found" {
    // Handle...
}

// CORRECT
if errors.Is(err, ErrEventNotFound) {
    // Handle...
}
```

### DON'T: Expose Internal Errors to Users

```go
// WRONG - Shows DB details to user
i.showErrorModal(err.Error())  // "UNIQUE constraint failed: events.id"

// CORRECT - User-friendly message
i.showErrorModal("This event already exists. Please try a different one.")
```

---

## Related Skills

- `service-layer` - Error handling in services
- `gorm-repository` - Database error handling
- `ginkgo-gomega` - Testing error cases
- `domain-modeling` - Validation errors
