---
name: service-layer
description: KaRiya service layer patterns for business logic orchestration
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Service Layer Skill

You are an expert in KaRiya's service layer patterns for business logic orchestration.

## Overview

Services in KaRiya orchestrate business operations, coordinating between domain entities and repositories while keeping domain logic pure.

## Service Location

```
internal/service/
├── career/
│   ├── service.go          # Main career service
│   ├── errors.go           # Service-specific errors
│   ├── cv/                  # CV generation service
│   ├── burstfact/          # Burst detection service
│   └── skillinference/     # Skill inference service
```

---

## Service Pattern

### Structure

```go
// internal/service/career/service.go
package career

import (
    "context"
    
    "github.com/baphled/kariya/internal/domain/career"
    repo "github.com/baphled/kariya/internal/repository/career"
    "github.com/baphled/kariya/internal/logger"
)

// Service orchestrates career-related business operations.
type Service struct {
    eventRepo repo.EventRepository
    factRepo  repo.FactRepository  // Optional
    burstRepo repo.BurstRepository // Optional
    logger    *logger.Logger
}

// NewService creates a new career service.
//
// Expected: eventRepo (required)
// Returns: Initialized service
// Side effects: None
func NewService(eventRepo repo.EventRepository) *Service {
    return &Service{
        eventRepo: eventRepo,
        logger:    logger.NewLogger("career-service"),
    }
}
```

### Optional Dependencies

```go
// SetFactRepository configures the optional fact repository.
//
// Expected: repo (can be nil to disable)
// Returns: None
// Side effects: Enables fact-related operations
func (s *Service) SetFactRepository(repo repo.FactRepository) {
    s.factRepo = repo
}

// SetBurstRepository configures the optional burst repository.
func (s *Service) SetBurstRepository(repo repo.BurstRepository) {
    s.burstRepo = repo
}
```

---

## Service Methods

### Create Operation

```go
// CreateEvent creates and persists a new career event.
//
// Expected: ctx, text (non-empty), date (not future)
// Returns: Created event or error
// Side effects: Persists to database, logs operation
func (s *Service) CreateEvent(ctx context.Context, text string, date time.Time) (*career.Event, error) {
    s.logger.WithFields(map[string]string{
        "operation": "create_event",
    }).Info("Creating event")
    
    // Create domain entity (validates internally)
    event, err := career.NewEvent(text, date)
    if err != nil {
        return nil, fmt.Errorf("invalid event: %w", err)
    }
    
    // Persist
    if err := s.eventRepo.Save(ctx, event); err != nil {
        s.logger.WithFields(map[string]string{
            "error": err.Error(),
        }).Error("Failed to save event")
        return nil, fmt.Errorf("failed to save event: %w", err)
    }
    
    s.logger.WithFields(map[string]string{
        "event_id": event.ID,
    }).Info("Event created successfully")
    
    return event, nil
}
```

### Read Operation

```go
// GetEvent retrieves an event by ID.
//
// Expected: ctx, id (non-empty)
// Returns: Event or ErrEventNotFound
// Side effects: None
func (s *Service) GetEvent(ctx context.Context, id string) (*career.Event, error) {
    if id == "" {
        return nil, errors.New("event ID is required")
    }
    
    event, err := s.eventRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, repo.ErrEventNotFound) {
            return nil, ErrEventNotFound
        }
        return nil, fmt.Errorf("failed to get event: %w", err)
    }
    
    return event, nil
}

// ListEvents retrieves events with optional filters.
//
// Expected: ctx, filters (can be empty struct for all events)
// Returns: Slice of events (may be empty), error
// Side effects: None
func (s *Service) ListEvents(ctx context.Context, filters EventFilters) ([]*career.Event, error) {
    events, err := s.eventRepo.FindWithFilters(ctx, repo.EventListFilters{
        Company:    filters.Company,
        Project:    filters.Project,
        StartDate:  filters.StartDate,
        EndDate:    filters.EndDate,
        SearchText: filters.SearchText,
        OrderBy:    filters.OrderBy,
        OrderDesc:  filters.OrderDesc,
        Limit:      filters.Limit,
        Offset:     filters.Offset,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list events: %w", err)
    }
    
    return events, nil
}
```

### Update Operation

```go
// UpdateEvent updates an existing event.
//
// Expected: ctx, event (with valid ID and data)
// Returns: Updated event or error
// Side effects: Updates database, logs operation
func (s *Service) UpdateEvent(ctx context.Context, event *career.Event) (*career.Event, error) {
    // Verify exists
    existing, err := s.eventRepo.FindByID(ctx, event.ID)
    if err != nil {
        if errors.Is(err, repo.ErrEventNotFound) {
            return nil, ErrEventNotFound
        }
        return nil, fmt.Errorf("failed to find event: %w", err)
    }
    
    // Validate updated data
    if err := event.Validate(); err != nil {
        return nil, fmt.Errorf("invalid event: %w", err)
    }
    
    // Preserve timestamps
    event.CreatedAt = existing.CreatedAt
    event.UpdatedAt = time.Now()
    
    // Save
    if err := s.eventRepo.Save(ctx, event); err != nil {
        return nil, fmt.Errorf("failed to update event: %w", err)
    }
    
    return event, nil
}
```

### Delete Operation

```go
// DeleteEvent removes an event by ID.
//
// Expected: ctx, id (non-empty)
// Returns: Error if not found or deletion fails
// Side effects: Removes from database, may cascade to related entities
func (s *Service) DeleteEvent(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("event ID is required")
    }
    
    // Check exists
    _, err := s.eventRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, repo.ErrEventNotFound) {
            return ErrEventNotFound
        }
        return fmt.Errorf("failed to find event: %w", err)
    }
    
    // Delete related facts if configured
    if s.factRepo != nil {
        if err := s.factRepo.DeleteByEventID(ctx, id); err != nil {
            s.logger.WithFields(map[string]string{
                "event_id": id,
                "error":    err.Error(),
            }).Warn("Failed to delete related facts")
        }
    }
    
    // Delete event
    if err := s.eventRepo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete event: %w", err)
    }
    
    return nil
}
```

---

## Optional Feature Pattern

```go
// ExtractFacts extracts facts from an event.
//
// Expected: ctx, eventID
// Returns: Extracted facts or ErrFactRepositoryNotConfigured
// Side effects: Persists facts to database
func (s *Service) ExtractFacts(ctx context.Context, eventID string) ([]*career.Fact, error) {
    // Check optional dependency
    if s.factRepo == nil {
        return nil, ErrFactRepositoryNotConfigured
    }
    
    // Get event
    event, err := s.eventRepo.FindByID(ctx, eventID)
    if err != nil {
        return nil, err
    }
    
    // Extract facts (business logic)
    facts := s.extractFactsFromText(event.Text)
    
    // Link to event
    for _, fact := range facts {
        fact.EventID = eventID
    }
    
    // Persist
    if err := s.factRepo.SaveAll(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to save facts: %w", err)
    }
    
    return facts, nil
}
```

---

## Service Errors

```go
// internal/service/career/errors.go
package career

import "errors"

// Service-level sentinel errors.
var (
    ErrEventNotFound              = errors.New("event not found")
    ErrFactRepositoryNotConfigured = errors.New("fact repository not configured")
    ErrBurstRepositoryNotConfigured = errors.New("burst repository not configured")
    ErrInvalidEventData           = errors.New("invalid event data")
)
```

---

## Logging Pattern

```go
import "github.com/baphled/kariya/internal/logger"

func (s *Service) SomeOperation(ctx context.Context) error {
    // Log entry
    s.logger.WithFields(map[string]string{
        "operation": "some_operation",
    }).Info("Starting operation")
    
    // ... do work ...
    
    if err != nil {
        // Log error with context
        s.logger.WithFields(map[string]string{
            "operation": "some_operation",
            "error":     err.Error(),
        }).Error("Operation failed")
        return err
    }
    
    // Log success
    s.logger.WithFields(map[string]string{
        "operation": "some_operation",
        "result":    "success",
    }).Info("Operation completed")
    
    return nil
}
```

---

## Context Propagation

```go
// Always accept and pass context
func (s *Service) Operation(ctx context.Context, ...) error {
    // Pass to repository
    result, err := s.repo.Find(ctx, id)
    
    // Pass to other services
    facts, err := s.factService.Extract(ctx, eventID)
    
    // Use for cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Continue processing
    }
}
```

---

## Testing Services

```go
var _ = Describe("EventService", func() {
    var (
        ctrl      *gomock.Controller
        eventRepo *mocks.MockEventRepository
        factRepo  *mocks.MockFactRepository
        service   *career.Service
        ctx       context.Context
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        ctrl = gomock.NewController(GinkgoT())
        eventRepo = mocks.NewMockEventRepository(ctrl)
        factRepo = mocks.NewMockFactRepository(ctrl)
        
        service = career.NewService(eventRepo)
        service.SetFactRepository(factRepo)
    })
    
    AfterEach(func() {
        ctrl.Finish()
    })
    
    Describe("CreateEvent", func() {
        It("creates and saves event", func() {
            eventRepo.EXPECT().
                Save(ctx, gomock.Any()).
                Return(nil)
            
            event, err := service.CreateEvent(ctx, "Test event", time.Now())
            
            Expect(err).ToNot(HaveOccurred())
            Expect(event).ToNot(BeNil())
            Expect(event.Text).To(Equal("Test event"))
        })
        
        It("returns error for invalid data", func() {
            // No repo expectations - validation fails first
            
            _, err := service.CreateEvent(ctx, "", time.Now())
            
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("invalid"))
        })
    })
    
    Describe("ExtractFacts", func() {
        Context("when fact repo not configured", func() {
            BeforeEach(func() {
                service = career.NewService(eventRepo)
                // Don't set fact repo
            })
            
            It("returns configuration error", func() {
                _, err := service.ExtractFacts(ctx, "event-1")
                
                Expect(err).To(MatchError(career.ErrFactRepositoryNotConfigured))
            })
        })
    })
})
```

---

## Specialized Services

### CV Generation Service

```go
// internal/service/career/cv/cv_generation_service.go
package cv

type GenerationService struct {
    eventRepo repo.EventRepository
    factRepo  repo.FactRepository
    config    *GenerationConfig
}

func NewGenerationService(eventRepo repo.EventRepository, factRepo repo.FactRepository) *GenerationService {
    return &GenerationService{
        eventRepo: eventRepo,
        factRepo:  factRepo,
        config:    DefaultConfig(),
    }
}

func (s *GenerationService) Generate(ctx context.Context, opts GenerateOptions) (*career.CV, error) {
    // Gather events
    events, err := s.eventRepo.FindWithFilters(ctx, opts.EventFilters)
    if err != nil {
        return nil, err
    }
    
    // Gather facts
    facts, err := s.factRepo.FindByEventIDs(ctx, getEventIDs(events))
    if err != nil {
        return nil, err
    }
    
    // Generate CV
    cv := s.buildCV(events, facts, opts)
    
    return cv, nil
}
```

---

## Anti-Patterns

### DON'T: Put Domain Logic in Service

```go
// WRONG - Validation belongs in domain
func (s *Service) CreateEvent(ctx context.Context, text string) (*Event, error) {
    if text == "" {
        return nil, errors.New("text required")  // Move to domain!
    }
    if len(text) > 2000 {
        return nil, errors.New("text too long")  // Move to domain!
    }
}

// CORRECT - Delegate to domain
func (s *Service) CreateEvent(ctx context.Context, text string) (*Event, error) {
    event, err := career.NewEvent(text, time.Now())  // Domain validates
    if err != nil {
        return nil, fmt.Errorf("invalid event: %w", err)
    }
}
```

### DON'T: Return Repository Errors Directly

```go
// WRONG - Leaks repository implementation
func (s *Service) GetEvent(ctx context.Context, id string) (*Event, error) {
    return s.repo.FindByID(ctx, id)  // Returns repo.ErrEventNotFound
}

// CORRECT - Translate to service errors
func (s *Service) GetEvent(ctx context.Context, id string) (*Event, error) {
    event, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, repo.ErrEventNotFound) {
            return nil, ErrEventNotFound  // Service error
        }
        return nil, fmt.Errorf("failed to get event: %w", err)
    }
    return event, nil
}
```

### DON'T: Skip Context

```go
// WRONG - No context
func (s *Service) Process() error {
    s.repo.FindAll()  // No cancellation support!
}

// CORRECT - Always use context
func (s *Service) Process(ctx context.Context) error {
    s.repo.FindAll(ctx)
}
```

### DON'T: Panic on Optional Dependencies

```go
// WRONG - Panic if not configured
func (s *Service) ExtractFacts(ctx context.Context, id string) ([]*Fact, error) {
    return s.factRepo.Find(ctx, id)  // Panics if nil!
}

// CORRECT - Return error
func (s *Service) ExtractFacts(ctx context.Context, id string) ([]*Fact, error) {
    if s.factRepo == nil {
        return nil, ErrFactRepositoryNotConfigured
    }
    return s.factRepo.Find(ctx, id)
}
```

---

## Related Skills

- `domain-modeling` - Domain entities used by services
- `gorm-repository` - Repository layer that services use
- `error-handling` - Error patterns in services
- `gomock` - Mocking repositories for service tests
