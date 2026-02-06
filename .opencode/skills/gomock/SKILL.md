# GoMock Skill

You are an expert in GoMock for generating and using mock implementations of Go interfaces.

## Overview

KaRiya uses GoMock to generate mock implementations of repository and service interfaces for unit testing.

## Mock Generation

### go:generate Directive

```go
// internal/repository/career/event_repository.go
package career

//go:generate mockgen -destination=../../testutil/mocks/repository/event_repository_mock.go -package=mocks . EventRepository

type EventRepository interface {
    Save(ctx context.Context, event *career.Event) error
    FindByID(ctx context.Context, id string) (*career.Event, error)
    FindAll(ctx context.Context) ([]*career.Event, error)
    Delete(ctx context.Context, id string) error
}
```

### Running Generation

```bash
# Generate all mocks
go generate ./...

# Generate for specific package
go generate ./internal/repository/career/...

# Verify mocks are up to date
make generate  # If available in Makefile
```

### Mock File Location

```
internal/testutil/mocks/
├── repository/
│   ├── event_repository_mock.go
│   ├── fact_repository_mock.go
│   └── burst_repository_mock.go
└── service/
    └── career_service_mock.go
```

---

## Using Mocks in Tests

### Setup Pattern

```go
package service_test

import (
    "context"
    "testing"
    
    "github.com/golang/mock/gomock"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "github.com/baphled/kariya/internal/service/career"
    "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("EventService", func() {
    var (
        ctrl    *gomock.Controller
        repo    *mocks.MockEventRepository
        service *career.EventService
        ctx     context.Context
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        ctrl = gomock.NewController(GinkgoT())
        repo = mocks.NewMockEventRepository(ctrl)
        service = career.NewEventService(repo)
    })
    
    AfterEach(func() {
        ctrl.Finish()
    })
    
    // Tests here...
})
```

---

## EXPECT Patterns

### Basic Expectation

```go
It("saves the event", func() {
    event := fixtures.Event(1)
    
    // Expect Save to be called with any context and return nil
    repo.EXPECT().
        Save(ctx, event).
        Return(nil)
    
    err := service.Create(ctx, event)
    
    Expect(err).ToNot(HaveOccurred())
})
```

### With gomock.Any()

```go
It("saves any event", func() {
    // Match any event argument
    repo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Return(nil)
    
    err := service.Create(ctx, &Event{Text: "test"})
    Expect(err).ToNot(HaveOccurred())
})
```

### Returning Values

```go
It("returns event by ID", func() {
    expected := fixtures.Event(1)
    
    repo.EXPECT().
        FindByID(ctx, "event-1").
        Return(expected, nil)
    
    event, err := service.GetByID(ctx, "event-1")
    
    Expect(err).ToNot(HaveOccurred())
    Expect(event).To(Equal(expected))
})
```

### Returning Errors

```go
It("handles repository errors", func() {
    repo.EXPECT().
        FindByID(ctx, "not-found").
        Return(nil, repository.ErrEventNotFound)
    
    event, err := service.GetByID(ctx, "not-found")
    
    Expect(err).To(MatchError(repository.ErrEventNotFound))
    Expect(event).To(BeNil())
})
```

---

## Call Counts

### Times()

```go
// Called exactly once (default)
repo.EXPECT().
    Save(gomock.Any(), gomock.Any()).
    Times(1)

// Called exactly N times
repo.EXPECT().
    FindByID(gomock.Any(), gomock.Any()).
    Times(3)

// Never called
repo.EXPECT().
    Delete(gomock.Any(), gomock.Any()).
    Times(0)
```

### MinTimes / MaxTimes

```go
// At least once
repo.EXPECT().
    Save(gomock.Any(), gomock.Any()).
    MinTimes(1)

// At most twice
repo.EXPECT().
    FindByID(gomock.Any(), gomock.Any()).
    MaxTimes(2)
```

### AnyTimes

```go
// Called any number of times (including zero)
repo.EXPECT().
    Count(gomock.Any()).
    Return(int64(5), nil).
    AnyTimes()
```

---

## Argument Matchers

### gomock.Any()

```go
// Match any value
repo.EXPECT().
    Save(gomock.Any(), gomock.Any())
```

### gomock.Eq()

```go
// Match exact value
repo.EXPECT().
    FindByID(gomock.Any(), gomock.Eq("event-123"))
```

### Custom Matcher

```go
// Create custom matcher
type eventWithTextMatcher struct {
    expectedText string
}

func (m eventWithTextMatcher) Matches(x interface{}) bool {
    event, ok := x.(*career.Event)
    if !ok {
        return false
    }
    return event.Text == m.expectedText
}

func (m eventWithTextMatcher) String() string {
    return fmt.Sprintf("has text %q", m.expectedText)
}

func EventWithText(text string) gomock.Matcher {
    return eventWithTextMatcher{expectedText: text}
}

// Usage
repo.EXPECT().
    Save(gomock.Any(), EventWithText("expected text"))
```

### gomock.Not()

```go
// Match anything except this value
repo.EXPECT().
    FindByID(gomock.Any(), gomock.Not("invalid-id"))
```

---

## Call Ordering

### InOrder

```go
It("validates then saves", func() {
    gomock.InOrder(
        repo.EXPECT().
            FindByID(ctx, "event-1").
            Return(nil, repository.ErrEventNotFound),
        repo.EXPECT().
            Save(ctx, gomock.Any()).
            Return(nil),
    )
    
    err := service.CreateIfNotExists(ctx, fixtures.Event(1))
    Expect(err).ToNot(HaveOccurred())
})
```

### After()

```go
call1 := repo.EXPECT().
    FindAll(ctx).
    Return([]*career.Event{}, nil)

repo.EXPECT().
    Save(ctx, gomock.Any()).
    Return(nil).
    After(call1)  // Must happen after call1
```

---

## Do() for Side Effects

### Capture Arguments

```go
It("captures saved event", func() {
    var savedEvent *career.Event
    
    repo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Do(func(_ context.Context, event *career.Event) {
            savedEvent = event
        }).
        Return(nil)
    
    service.Create(ctx, &Event{Text: "test"})
    
    Expect(savedEvent).ToNot(BeNil())
    Expect(savedEvent.Text).To(Equal("test"))
})
```

### Conditional Returns

```go
It("returns different values based on input", func() {
    repo.EXPECT().
        FindByID(gomock.Any(), gomock.Any()).
        DoAndReturn(func(_ context.Context, id string) (*career.Event, error) {
            if id == "exists" {
                return fixtures.Event(1), nil
            }
            return nil, repository.ErrEventNotFound
        }).
        AnyTimes()
})
```

---

## Multiple Mock Setup

### Testing Service with Multiple Repos

```go
var _ = Describe("CareerService", func() {
    var (
        ctrl        *gomock.Controller
        eventRepo   *mocks.MockEventRepository
        factRepo    *mocks.MockFactRepository
        burstRepo   *mocks.MockBurstRepository
        service     *career.Service
        ctx         context.Context
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        ctrl = gomock.NewController(GinkgoT())
        
        eventRepo = mocks.NewMockEventRepository(ctrl)
        factRepo = mocks.NewMockFactRepository(ctrl)
        burstRepo = mocks.NewMockBurstRepository(ctrl)
        
        service = career.NewService(eventRepo)
        service.SetFactRepository(factRepo)
        service.SetBurstRepository(burstRepo)
    })
    
    AfterEach(func() {
        ctrl.Finish()
    })
    
    Describe("CreateEventWithFacts", func() {
        It("saves event and extracts facts", func() {
            event := fixtures.Event(1)
            facts := fixtures.Facts(3)
            
            eventRepo.EXPECT().Save(ctx, event).Return(nil)
            factRepo.EXPECT().SaveAll(ctx, gomock.Any()).Return(nil)
            
            err := service.CreateEventWithFacts(ctx, event)
            Expect(err).ToNot(HaveOccurred())
        })
    })
})
```

---

## Common Patterns

### Repository Not Found

```go
Context("when event not found", func() {
    BeforeEach(func() {
        repo.EXPECT().
            FindByID(ctx, "not-found").
            Return(nil, repository.ErrEventNotFound)
    })
    
    It("returns not found error", func() {
        _, err := service.GetByID(ctx, "not-found")
        Expect(err).To(MatchError(repository.ErrEventNotFound))
    })
})
```

### Empty Results

```go
Context("when no events exist", func() {
    BeforeEach(func() {
        repo.EXPECT().
            FindAll(ctx).
            Return([]*career.Event{}, nil)
    })
    
    It("returns empty slice", func() {
        events, err := service.GetAll(ctx)
        Expect(err).ToNot(HaveOccurred())
        Expect(events).To(BeEmpty())
    })
})
```

### Verifying No Calls

```go
It("does not save invalid event", func() {
    // Expect Save to never be called
    repo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Times(0)
    
    invalidEvent := &Event{Text: ""}  // Invalid: empty text
    err := service.Create(ctx, invalidEvent)
    
    Expect(err).To(HaveOccurred())
})
```

---

## Anti-Patterns

### DON'T: Forget ctrl.Finish()

```go
// WRONG - Expectations not verified
It("test without finish", func() {
    ctrl := gomock.NewController(GinkgoT())
    repo := mocks.NewMockEventRepository(ctrl)
    // Missing: ctrl.Finish()
})

// CORRECT - Use AfterEach
AfterEach(func() {
    ctrl.Finish()  // Verifies all expectations were met
})
```

### DON'T: Over-Specify

```go
// WRONG - Too specific, brittle test
repo.EXPECT().
    Save(
        gomock.Eq(context.Background()),
        gomock.Eq(&Event{ID: "exactly-this", Text: "exactly-this"}),
    )

// CORRECT - Specify only what matters
repo.EXPECT().
    Save(gomock.Any(), EventWithText("exactly-this"))
```

### DON'T: Ignore Context

```go
// WRONG - Hard-coded context
repo.EXPECT().
    Save(context.Background(), gomock.Any())

// CORRECT - Use test context variable
repo.EXPECT().
    Save(ctx, gomock.Any())

// OR use gomock.Any() for context
repo.EXPECT().
    Save(gomock.Any(), gomock.Any())
```

### DON'T: Mock Everything

```go
// WRONG - Mocking domain logic
func (m *MockEvent) Validate() error {
    return nil
}

// CORRECT - Only mock external dependencies (repos, services)
// Use real domain entities
event := &career.Event{Text: "test"}
Expect(event.Validate()).ToNot(HaveOccurred())
```

---

## Debugging Mocks

### Unexpected Call Error

```
unexpected call to *mocks.MockEventRepository.FindByID
expected call at /path/to/test.go:42 was not matched
```

**Fix**: Add missing EXPECT or use AnyTimes()

### Missing Call Error

```
missing call(s) to *mocks.MockEventRepository.Save
expected call at /path/to/test.go:42 doesn't have matching call
```

**Fix**: Ensure the code under test actually makes the expected call

### Argument Mismatch

```
doesn't match expected arguments for *mocks.MockEventRepository.FindByID
Got: ctx, "actual-id"
Expected: ctx, "expected-id"
```

**Fix**: Use correct expected values or gomock.Any()

---

## Related Skills

- `ginkgo-gomega` - Testing framework for using mocks
- `test-fixtures` - Creating test data for mock returns
- `service-layer` - Services that use mocked dependencies
