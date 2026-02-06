# Test Fixtures Skill

You are an expert in creating test data using factory-go and gofakeit for Go tests.

## Overview

KaRiya uses factory-go for structured test data creation and gofakeit for realistic fake data generation.

## Factory Location

```
internal/testutil/fixtures/
├── event_factory.go   # Event test data
├── fact_factory.go    # Fact test data
├── burst_factory.go   # Burst test data
└── skill_factory.go   # Skill test data
```

---

## Basic Factory Pattern

### Factory Definition

```go
// internal/testutil/fixtures/event_factory.go
package fixtures

import (
    "time"
    
    "github.com/bluele/factory-go/factory"
    "github.com/brianvoe/gofakeit/v7"
    "github.com/baphled/kariya/internal/domain/career"
)

var eventFactory = factory.NewFactory(
    &career.Event{},
).SeqInt("ID", func(n int) (interface{}, error) {
    return fmt.Sprintf("event-%d", n), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
    return gofakeit.Sentence(8), nil
}).Attr("Date", func(args factory.Args) (interface{}, error) {
    return gofakeit.DateRange(
        time.Now().AddDate(-2, 0, 0),
        time.Now(),
    ), nil
}).Attr("Company", func(args factory.Args) (interface{}, error) {
    return gofakeit.Company(), nil
}).Attr("Project", func(args factory.Args) (interface{}, error) {
    return gofakeit.BuzzWord(), nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
    return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
    return time.Now(), nil
})
```

---

## Helper Functions

### Single Entity

```go
// Event returns a pointer to a test event with the given ID suffix.
//
// Expected: n (ID suffix number)
// Returns: *Event with sequential ID and random data
// Side effects: None
func Event(n int) *career.Event {
    event := eventFactory.MustCreateWithOption(map[string]interface{}{
        "ID": fmt.Sprintf("event-%d", n),
    }).(*career.Event)
    return event
}

// EventVal returns a value (not pointer) for inline struct comparisons.
func EventVal(n int) career.Event {
    return *Event(n)
}
```

### Multiple Entities

```go
// Events returns a slice of n test events.
//
// Expected: n (count of events to create)
// Returns: Slice of event pointers with sequential IDs
// Side effects: None
func Events(n int) []*career.Event {
    events := make([]*career.Event, n)
    for i := 0; i < n; i++ {
        events[i] = Event(i + 1)
    }
    return events
}
```

### With Custom Attributes

```go
// EventOption is a function that modifies an event.
type EventOption func(*career.Event)

// WithText sets the event text.
func WithText(text string) EventOption {
    return func(e *career.Event) {
        e.Text = text
    }
}

// WithCompany sets the company.
func WithCompany(company string) EventOption {
    return func(e *career.Event) {
        e.Company = company
    }
}

// WithDate sets the date.
func WithDate(date time.Time) EventOption {
    return func(e *career.Event) {
        e.Date = date
    }
}

// WithProject sets the project.
func WithProject(project string) EventOption {
    return func(e *career.Event) {
        e.Project = project
    }
}

// WithSkills sets the skills.
func WithSkills(skills ...string) EventOption {
    return func(e *career.Event) {
        e.Skills = skills
    }
}

// EventWith creates an event with custom options.
//
// Expected: opts (variadic EventOptions)
// Returns: *Event with applied options
// Side effects: None
func EventWith(opts ...EventOption) *career.Event {
    event := Event(1)
    for _, opt := range opts {
        opt(event)
    }
    return event
}
```

---

## Using Gofakeit

### Common Generators

```go
import "github.com/brianvoe/gofakeit/v7"

// Text
gofakeit.Sentence(wordCount)      // "The quick brown fox."
gofakeit.Paragraph(sentenceCount) // Multiple sentences
gofakeit.BuzzWord()               // "synergize"
gofakeit.HackerPhrase()           // "We need to hack the neural API!"

// Names & Companies
gofakeit.Name()           // "John Smith"
gofakeit.Company()        // "Acme Corp"
gofakeit.JobTitle()       // "Senior Developer"
gofakeit.Email()          // "john@example.com"

// Dates
gofakeit.Date()                           // Random date
gofakeit.DateRange(start, end)            // Date within range
gofakeit.FutureDate()                     // Future date
gofakeit.PastDate()                       // Past date

// Numbers
gofakeit.Number(min, max)                 // Random int
gofakeit.Float64Range(min, max)           // Random float
gofakeit.UUID()                           // UUID string

// Collections
gofakeit.RandomString([]string{"a", "b"}) // Random from slice
```

### Seeded Randomness (For Reproducibility)

```go
import "github.com/brianvoe/gofakeit/v7"

// Set seed for reproducible tests
gofakeit.Seed(12345)

// Now all generated values will be consistent
name := gofakeit.Name()  // Always same name with this seed
```

---

## Factory Patterns

### Fact Factory

```go
var factFactory = factory.NewFactory(
    &career.Fact{},
).SeqInt("ID", func(n int) (interface{}, error) {
    return fmt.Sprintf("fact-%d", n), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
    return fmt.Sprintf("Demonstrated %s skills", gofakeit.BuzzWord()), nil
}).Attr("Category", func(args factory.Args) (interface{}, error) {
    categories := []string{"technical", "leadership", "communication"}
    return gofakeit.RandomString(categories), nil
}).Attr("RoleFit", func(args factory.Args) (interface{}, error) {
    fits := []career.RoleFit{career.RoleFitHigh, career.RoleFitMedium, career.RoleFitLow}
    return fits[gofakeit.Number(0, 2)], nil
}).Attr("Confidence", func(args factory.Args) (interface{}, error) {
    return gofakeit.Float64Range(0.5, 1.0), nil
})

func Fact(n int) *career.Fact {
    return factFactory.MustCreateWithOption(map[string]interface{}{
        "ID": fmt.Sprintf("fact-%d", n),
    }).(*career.Fact)
}

func Facts(n int) []*career.Fact {
    facts := make([]*career.Fact, n)
    for i := 0; i < n; i++ {
        facts[i] = Fact(i + 1)
    }
    return facts
}

func FactWith(opts ...FactOption) *career.Fact {
    fact := Fact(1)
    for _, opt := range opts {
        opt(fact)
    }
    return fact
}
```

### Burst Factory

```go
var burstFactory = factory.NewFactory(
    &career.Burst{},
).SeqInt("ID", func(n int) (interface{}, error) {
    return fmt.Sprintf("burst-%d", n), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
    return fmt.Sprintf("%s Project", gofakeit.BuzzWord()), nil
}).Attr("StartDate", func(args factory.Args) (interface{}, error) {
    return gofakeit.DateRange(
        time.Now().AddDate(-1, 0, 0),
        time.Now().AddDate(0, -1, 0),
    ), nil
}).Attr("EndDate", func(args factory.Args) (interface{}, error) {
    return gofakeit.DateRange(
        time.Now().AddDate(0, -1, 0),
        time.Now(),
    ), nil
})

func Burst(n int) *career.Burst {
    return burstFactory.MustCreateWithOption(map[string]interface{}{
        "ID": fmt.Sprintf("burst-%d", n),
    }).(*career.Burst)
}

// BurstWithEvents creates a burst with attached events.
func BurstWithEvents(eventCount int) *career.Burst {
    burst := Burst(1)
    burst.Events = Events(eventCount)
    return burst
}
```

---

## Using Fixtures in Tests

### Basic Usage

```go
var _ = Describe("EventService", func() {
    It("saves an event", func() {
        event := fixtures.Event(1)
        
        err := service.Save(ctx, event)
        
        Expect(err).ToNot(HaveOccurred())
    })
    
    It("processes multiple events", func() {
        events := fixtures.Events(5)
        
        for _, e := range events {
            err := service.Save(ctx, e)
            Expect(err).ToNot(HaveOccurred())
        }
    })
})
```

### With Custom Data

```go
It("filters by company", func() {
    // Create events with specific companies
    acmeEvents := []*career.Event{
        fixtures.EventWith(fixtures.WithCompany("Acme Corp")),
        fixtures.EventWith(fixtures.WithCompany("Acme Corp")),
    }
    otherEvent := fixtures.EventWith(fixtures.WithCompany("Other Inc"))
    
    for _, e := range append(acmeEvents, otherEvent) {
        repo.Save(ctx, e)
    }
    
    results, err := repo.FindWithFilters(ctx, Filters{Company: "Acme Corp"})
    
    Expect(err).ToNot(HaveOccurred())
    Expect(results).To(HaveLen(2))
})
```

### With Date Ranges

```go
It("finds events in date range", func() {
    lastWeek := time.Now().AddDate(0, 0, -7)
    yesterday := time.Now().AddDate(0, 0, -1)
    lastMonth := time.Now().AddDate(0, -1, 0)
    
    recentEvent := fixtures.EventWith(fixtures.WithDate(yesterday))
    oldEvent := fixtures.EventWith(fixtures.WithDate(lastMonth))
    
    // Save events...
    
    results, err := repo.FindWithFilters(ctx, Filters{
        StartDate: lastWeek,
        EndDate:   time.Now(),
    })
    
    Expect(results).To(HaveLen(1))
    Expect(results[0].ID).To(Equal(recentEvent.ID))
})
```

---

## Mock Data for UI Tests

### Timeline Data

```go
// TimelineData creates realistic timeline test data.
func TimelineData() []*career.Event {
    events := []*career.Event{
        EventWith(
            WithText("Led cross-functional team to deliver microservices migration"),
            WithCompany("TechCorp"),
            WithProject("Platform Modernization"),
            WithDate(time.Now().AddDate(0, -1, 0)),
        ),
        EventWith(
            WithText("Implemented CI/CD pipeline reducing deployment time by 60%"),
            WithCompany("TechCorp"),
            WithProject("DevOps Initiative"),
            WithDate(time.Now().AddDate(0, -2, 0)),
        ),
        EventWith(
            WithText("Mentored 3 junior developers on Go best practices"),
            WithCompany("StartupXYZ"),
            WithProject("Team Growth"),
            WithDate(time.Now().AddDate(0, -6, 0)),
        ),
    }
    return events
}
```

---

## Anti-Patterns

### DON'T: Use Real Data

```go
// WRONG - Real data in tests
event := &career.Event{
    Text:    "Led project at my actual company",  // PII risk!
    Company: "My Real Company",
}

// CORRECT - Use fixtures
event := fixtures.Event(1)
```

### DON'T: Hardcode IDs

```go
// WRONG - Hardcoded IDs may conflict
event1 := &Event{ID: "test-1"}
event2 := &Event{ID: "test-1"}  // Duplicate!

// CORRECT - Use sequential factory
event1 := fixtures.Event(1)  // "event-1"
event2 := fixtures.Event(2)  // "event-2"
```

### DON'T: Create in Test Body

```go
// WRONG - Verbose, repetitive
It("test", func() {
    event := &career.Event{
        ID:        "test-1",
        Text:      "Some text",
        Date:      time.Now(),
        Company:   "Test Co",
        Project:   "Test Project",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    // ... 
})

// CORRECT - Use factory
It("test", func() {
    event := fixtures.Event(1)
    // ...
})
```

### DON'T: Share Mutable Fixtures

```go
// WRONG - Shared mutable state
var sharedEvent = fixtures.Event(1)

It("test 1", func() {
    sharedEvent.Text = "Modified"  // Affects other tests!
})

It("test 2", func() {
    Expect(sharedEvent.Text).To(Equal("original"))  // Fails!
})

// CORRECT - Create fresh for each test
BeforeEach(func() {
    event = fixtures.Event(1)
})
```

---

## Related Skills

- `ginkgo-gomega` - Testing framework using fixtures
- `gomock` - Mocking with fixture data
- `e2e-testing` - E2E tests with fixture data
