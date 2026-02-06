# Domain Modeling Skill

You are an expert in Domain-Driven Design (DDD) and KaRiya's domain modeling patterns.

## Overview

KaRiya follows DDD principles with a clear separation between domain entities (pure business logic) and infrastructure concerns (persistence, UI).

## Domain Layer Structure

```
internal/domain/
├── career/
│   ├── event.go      # Event entity
│   ├── fact.go       # Fact entity
│   ├── burst.go      # Burst entity (event grouping)
│   ├── skill.go      # Skill entity
│   └── cv.go         # CV aggregate
└── constants/
    ├── role_fit.go   # RoleFit enum
    └── capture_mode.go
```

## Entity Pattern

### Structure

```go
// internal/domain/career/event.go
package career

import (
    "errors"
    "time"
)

// Event represents a career event (accomplishment, task, milestone).
//
// An Event is the core unit of career tracking. It captures what happened,
// when, and in what context (company, project, role).
type Event struct {
    ID        string
    Text      string
    Date      time.Time
    Company   string
    Project   string
    Role      string
    Skills    []string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Validation Method

```go
// Validate checks the event for business rule compliance.
//
// Expected: Populated event fields
// Returns: Error if validation fails, nil if valid
// Side effects: None
func (e *Event) Validate() error {
    if err := e.validateText(); err != nil {
        return err
    }
    if err := e.validateDate(); err != nil {
        return err
    }
    return nil
}

func (e *Event) validateText() error {
    if e.Text == "" {
        return errors.New("event text is required")
    }
    if len(e.Text) > 2000 {
        return errors.New("event text exceeds maximum length of 2000 characters")
    }
    return nil
}

func (e *Event) validateDate() error {
    if e.Date.IsZero() {
        return errors.New("event date is required")
    }
    if e.Date.After(time.Now().Add(24 * time.Hour)) {
        return errors.New("event date cannot be in the future")
    }
    return nil
}
```

### Constructor Function

```go
// NewEvent creates a new Event with required fields.
//
// Expected: text (required, non-empty), date (required, not future)
// Returns: New Event with generated ID, or error if validation fails
// Side effects: Generates UUID for ID
func NewEvent(text string, date time.Time) (*Event, error) {
    event := &Event{
        ID:        uuid.New().String(),
        Text:      text,
        Date:      date,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := event.Validate(); err != nil {
        return nil, err
    }
    
    return event, nil
}
```

### Clone Method (For Immutability)

```go
// Clone creates a deep copy of the event.
//
// Expected: None
// Returns: New Event with copied values
// Side effects: None
func (e *Event) Clone() *Event {
    skills := make([]string, len(e.Skills))
    copy(skills, e.Skills)
    
    return &Event{
        ID:        e.ID,
        Text:      e.Text,
        Date:      e.Date,
        Company:   e.Company,
        Project:   e.Project,
        Role:      e.Role,
        Skills:    skills,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}
```

---

## Type Aliases for Enums

### Defining Enums

```go
// internal/domain/constants/role_fit.go
package constants

// RoleFit represents how well a fact aligns with a target role.
type RoleFit string

const (
    RoleFitHigh   RoleFit = "high"
    RoleFitMedium RoleFit = "medium"
    RoleFitLow    RoleFit = "low"
)

// Valid returns true if the RoleFit is a known value.
func (r RoleFit) Valid() bool {
    switch r {
    case RoleFitHigh, RoleFitMedium, RoleFitLow:
        return true
    default:
        return false
    }
}

// AllRoleFits returns all valid RoleFit values.
func AllRoleFits() []RoleFit {
    return []RoleFit{RoleFitHigh, RoleFitMedium, RoleFitLow}
}
```

### Using Type Aliases in Domain

```go
// internal/domain/career/fact.go
package career

import "github.com/baphled/kariya/internal/domain/constants"

// Type alias for cleaner usage
type RoleFit = constants.RoleFit

// Use constants
const (
    RoleFitHigh   = constants.RoleFitHigh
    RoleFitMedium = constants.RoleFitMedium
    RoleFitLow    = constants.RoleFitLow
)

type Fact struct {
    ID      string
    Text    string
    RoleFit RoleFit  // Uses type alias
    // ...
}
```

---

## Aggregate Pattern

### CV as Aggregate Root

```go
// internal/domain/career/cv.go
package career

// CV is an aggregate root containing all CV-related data.
type CV struct {
    Profile      Profile
    Summary      string
    Experience   []Experience
    Skills       []SkillCategory
    Education    []Education
    Achievements []string
}

// Profile contains personal information.
type Profile struct {
    Name     string
    Title    string
    Email    string
    Phone    string
    Location string
    LinkedIn string
    GitHub   string
}

// Experience represents a work experience entry.
type Experience struct {
    Company     string
    Role        string
    StartDate   time.Time
    EndDate     *time.Time  // Nil if current
    Description string
    Highlights  []string
}
```

### Aggregate Validation

```go
func (cv *CV) Validate() error {
    if err := cv.Profile.Validate(); err != nil {
        return fmt.Errorf("profile: %w", err)
    }
    
    for i, exp := range cv.Experience {
        if err := exp.Validate(); err != nil {
            return fmt.Errorf("experience[%d]: %w", i, err)
        }
    }
    
    return nil
}

func (p *Profile) Validate() error {
    if p.Name == "" {
        return errors.New("name is required")
    }
    if p.Email != "" && !isValidEmail(p.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

---

## Value Objects

### Immutable Value Object

```go
// DateRange represents a period of time.
type DateRange struct {
    start time.Time
    end   time.Time
}

// NewDateRange creates a validated date range.
func NewDateRange(start, end time.Time) (DateRange, error) {
    if end.Before(start) {
        return DateRange{}, errors.New("end date cannot be before start date")
    }
    return DateRange{start: start, end: end}, nil
}

func (d DateRange) Start() time.Time { return d.start }
func (d DateRange) End() time.Time   { return d.end }

func (d DateRange) Contains(t time.Time) bool {
    return !t.Before(d.start) && !t.After(d.end)
}

func (d DateRange) Duration() time.Duration {
    return d.end.Sub(d.start)
}
```

---

## Domain Rules

### Business Rules in Domain

```go
// Fact represents an extracted competency from an event.
type Fact struct {
    ID          string
    EventID     string
    Text        string
    Category    string
    RoleFit     RoleFit
    Confidence  float64
    Skills      []string
}

func (f *Fact) Validate() error {
    if f.Text == "" {
        return errors.New("fact text is required")
    }
    if f.EventID == "" {
        return errors.New("fact must be linked to an event")
    }
    if !f.RoleFit.Valid() {
        return fmt.Errorf("invalid role fit: %s", f.RoleFit)
    }
    if f.Confidence < 0 || f.Confidence > 1 {
        return errors.New("confidence must be between 0 and 1")
    }
    return nil
}

// IsHighValue determines if the fact is valuable for CV inclusion.
func (f *Fact) IsHighValue() bool {
    return f.RoleFit == RoleFitHigh && f.Confidence >= 0.8
}
```

---

## Domain Entity Guidelines

### MUST Have

1. **ID field** - Unique identifier (usually UUID string)
2. **Validate() method** - Business rule validation
3. **No persistence logic** - Pure domain, no DB concerns
4. **No UI logic** - No rendering, formatting for display

### SHOULD Have

1. **Constructor function** - `NewXxx()` that validates
2. **Clone() method** - For immutability when needed
3. **Timestamps** - CreatedAt, UpdatedAt for audit

### MUST NOT Have

1. **GORM tags** - Those belong in repository models
2. **JSON tags** - Unless specifically needed for API
3. **Database operations** - No `Save()`, `Load()` methods
4. **External dependencies** - No imports from infrastructure

---

## Testing Domain Entities

```go
var _ = Describe("Event", func() {
    Describe("Validate", func() {
        Context("when text is empty", func() {
            It("returns error", func() {
                event := &Event{Text: "", Date: time.Now()}
                
                err := event.Validate()
                
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("text"))
            })
        })
        
        Context("when date is in future", func() {
            It("returns error", func() {
                event := &Event{
                    Text: "Valid text",
                    Date: time.Now().Add(7 * 24 * time.Hour),
                }
                
                err := event.Validate()
                
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("future"))
            })
        })
        
        Context("when all fields valid", func() {
            It("returns nil", func() {
                event := &Event{
                    Text: "Valid event",
                    Date: time.Now(),
                }
                
                err := event.Validate()
                
                Expect(err).ToNot(HaveOccurred())
            })
        })
    })
    
    Describe("NewEvent", func() {
        It("creates event with generated ID", func() {
            event, err := NewEvent("Test event", time.Now())
            
            Expect(err).ToNot(HaveOccurred())
            Expect(event.ID).ToNot(BeEmpty())
            Expect(event.Text).To(Equal("Test event"))
        })
        
        It("rejects invalid input", func() {
            event, err := NewEvent("", time.Now())
            
            Expect(err).To(HaveOccurred())
            Expect(event).To(BeNil())
        })
    })
})
```

---

## Domain vs Model Separation

### Domain Entity (Pure)

```go
// internal/domain/career/event.go
type Event struct {
    ID      string
    Text    string
    Date    time.Time
    Skills  []string  // Simple slice
}
```

### Repository Model (GORM)

```go
// internal/repository/models/career.go
type Event struct {
    ID     string  `gorm:"primaryKey;type:text"`
    Text   string  `gorm:"type:text;not null"`
    Date   time.Time `gorm:"type:datetime"`
    Skills []Skill `gorm:"many2many:event_skills"`  // Relationship
}

// Conversion methods
func (e *Event) ToDomain() *career.Event { ... }
func EventFromDomain(d *career.Event) *Event { ... }
```

---

## Anti-Patterns

### DON'T: Add Persistence to Domain

```go
// WRONG - Domain knows about database
type Event struct {
    ID string `gorm:"primaryKey"`  // NO GORM tags!
}

func (e *Event) Save(db *gorm.DB) error {  // NO DB operations!
    return db.Save(e).Error
}
```

### DON'T: Add UI Logic to Domain

```go
// WRONG - Domain knows about rendering
func (e *Event) Format() string {
    return lipgloss.NewStyle().Render(e.Text)  // NO UI!
}
```

### DON'T: Skip Validation

```go
// WRONG - Accepting invalid state
func NewEvent(text string) *Event {
    return &Event{Text: text}  // No validation!
}

// CORRECT
func NewEvent(text string) (*Event, error) {
    e := &Event{Text: text}
    if err := e.Validate(); err != nil {
        return nil, err
    }
    return e, nil
}
```

### DON'T: Use Raw Strings for Enums

```go
// WRONG - String without type safety
type Event struct {
    Status string  // Could be anything!
}

// CORRECT - Typed enum
type EventStatus string
const (
    EventStatusDraft     EventStatus = "draft"
    EventStatusPublished EventStatus = "published"
)

type Event struct {
    Status EventStatus
}
```

---

## Related Skills

- `gorm-repository` - Persistence layer patterns
- `service-layer` - Business logic orchestration
- `error-handling` - Domain validation errors
