# GORM Repository Skill

You are an expert in GORM ORM, SQLite, and KaRiya's repository patterns.

## Overview

KaRiya uses GORM with SQLite for persistence, following a clean repository pattern with interface-based design.

## Repository Architecture

### Layer Structure

```
Repository Interface (internal/repository/career/)
         │
    ┌────┴────┐
    │         │
    ▼         ▼
SQL Impl   Memory Impl
(sql/)     (memory/)
    │
    ▼
GORM Models (internal/repository/models/)
    │
    ▼
SQLite Database
```

### Interface Definition

```go
// internal/repository/career/event_repository.go
package career

//go:generate mockgen -destination=../testutil/mocks/repository/event_repository_mock.go -package=mocks . EventRepository

// EventRepository defines the contract for event persistence.
//
// Expected: Valid domain entities for write operations
// Returns: Domain entities, never GORM models
// Side effects: Database operations
type EventRepository interface {
    Save(ctx context.Context, event *career.Event) error
    FindByID(ctx context.Context, id string) (*career.Event, error)
    FindAll(ctx context.Context) ([]*career.Event, error)
    FindWithFilters(ctx context.Context, filters EventListFilters) ([]*career.Event, error)
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}

// Sentinel errors
var (
    ErrEventNotFound   = errors.New("event not found")
    ErrDuplicateEvent  = errors.New("duplicate event")
)
```

---

## SQL Implementation

### Structure

```go
// internal/repository/career/sql/event_repository.go
package sql

import (
    "context"
    "gorm.io/gorm"
    
    "github.com/baphled/kariya/internal/domain/career"
    repo "github.com/baphled/kariya/internal/repository/career"
    "github.com/baphled/kariya/internal/repository/models"
)

// Verify interface implementation
var _ repo.EventRepository = (*EventRepository)(nil)

// EventRepository implements EventRepository using GORM/SQLite.
type EventRepository struct {
    db *gorm.DB
}

// NewEventRepository creates a new SQL-backed event repository.
//
// Expected: db (required, must be initialized)
// Returns: Repository instance ready for use
// Side effects: None
func NewEventRepository(db *gorm.DB) *EventRepository {
    return &EventRepository{db: db}
}
```

### CRUD Operations

#### Save (Create/Update)

```go
func (r *EventRepository) Save(ctx context.Context, event *career.Event) error {
    // Convert domain to model
    model := models.EventFromDomain(event)
    
    // Use GORM's Save (upsert behavior)
    result := r.db.WithContext(ctx).Save(model)
    if result.Error != nil {
        return fmt.Errorf("failed to save event: %w", result.Error)
    }
    
    return nil
}
```

#### FindByID

```go
func (r *EventRepository) FindByID(ctx context.Context, id string) (*career.Event, error) {
    var model models.Event
    
    result := r.db.WithContext(ctx).
        Preload("Skills").  // Eager load relationships
        First(&model, "id = ?", id)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, repo.ErrEventNotFound
        }
        return nil, fmt.Errorf("failed to find event: %w", result.Error)
    }
    
    return model.ToDomain(), nil
}
```

#### FindAll with Preloading

```go
func (r *EventRepository) FindAll(ctx context.Context) ([]*career.Event, error) {
    var models []models.Event
    
    result := r.db.WithContext(ctx).
        Preload("Skills").
        Order("date DESC").
        Find(&models)
    
    if result.Error != nil {
        return nil, fmt.Errorf("failed to find events: %w", result.Error)
    }
    
    events := make([]*career.Event, len(models))
    for i, m := range models {
        events[i] = m.ToDomain()
    }
    
    return events, nil
}
```

#### FindWithFilters

```go
func (r *EventRepository) FindWithFilters(ctx context.Context, filters repo.EventListFilters) ([]*career.Event, error) {
    query := r.db.WithContext(ctx).Model(&models.Event{})
    
    // Apply filters
    if filters.Company != "" {
        query = query.Where("company = ?", filters.Company)
    }
    if filters.Project != "" {
        query = query.Where("project = ?", filters.Project)
    }
    if !filters.StartDate.IsZero() {
        query = query.Where("date >= ?", filters.StartDate)
    }
    if !filters.EndDate.IsZero() {
        query = query.Where("date <= ?", filters.EndDate)
    }
    if filters.SearchText != "" {
        query = query.Where("text LIKE ?", "%"+filters.SearchText+"%")
    }
    
    // Apply ordering
    if filters.OrderBy != "" {
        direction := "ASC"
        if filters.OrderDesc {
            direction = "DESC"
        }
        query = query.Order(fmt.Sprintf("%s %s", filters.OrderBy, direction))
    }
    
    // Apply pagination
    if filters.Limit > 0 {
        query = query.Limit(filters.Limit)
    }
    if filters.Offset > 0 {
        query = query.Offset(filters.Offset)
    }
    
    var models []models.Event
    if err := query.Preload("Skills").Find(&models).Error; err != nil {
        return nil, fmt.Errorf("failed to find events: %w", err)
    }
    
    return r.toDomainSlice(models), nil
}
```

#### Delete

```go
func (r *EventRepository) Delete(ctx context.Context, id string) error {
    result := r.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", id)
    
    if result.Error != nil {
        return fmt.Errorf("failed to delete event: %w", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return repo.ErrEventNotFound
    }
    
    return nil
}
```

#### Count

```go
func (r *EventRepository) Count(ctx context.Context) (int64, error) {
    var count int64
    
    result := r.db.WithContext(ctx).Model(&models.Event{}).Count(&count)
    if result.Error != nil {
        return 0, fmt.Errorf("failed to count events: %w", result.Error)
    }
    
    return count, nil
}
```

---

## GORM Models

### Model Definition

```go
// internal/repository/models/career.go
package models

import (
    "time"
    "github.com/baphled/kariya/internal/domain/career"
)

// Event is the GORM model for career events.
type Event struct {
    ID        string    `gorm:"primaryKey;type:text"`
    Text      string    `gorm:"type:text;not null"`
    Date      time.Time `gorm:"type:datetime;not null"`
    Company   string    `gorm:"type:text"`
    Project   string    `gorm:"type:text"`
    Role      string    `gorm:"type:text"`
    Skills    []Skill   `gorm:"many2many:event_skills"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name.
func (Event) TableName() string {
    return "events"
}
```

### Domain Conversion

```go
// ToDomain converts GORM model to domain entity.
func (e *Event) ToDomain() *career.Event {
    skills := make([]string, len(e.Skills))
    for i, s := range e.Skills {
        skills[i] = s.Name
    }
    
    return &career.Event{
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

// EventFromDomain converts domain entity to GORM model.
func EventFromDomain(e *career.Event) *Event {
    skills := make([]Skill, len(e.Skills))
    for i, name := range e.Skills {
        skills[i] = Skill{Name: name}
    }
    
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

### Custom Types (StringSlice)

```go
// StringSlice is a custom type for storing string slices in SQLite.
type StringSlice []string

// Scan implements sql.Scanner for database reads.
func (s *StringSlice) Scan(value interface{}) error {
    if value == nil {
        *s = nil
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("failed to scan StringSlice: expected []byte, got %T", value)
    }
    
    if len(bytes) == 0 {
        *s = nil
        return nil
    }
    
    return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer for database writes.
func (s StringSlice) Value() (driver.Value, error) {
    if s == nil || len(s) == 0 {
        return nil, nil
    }
    return json.Marshal(s)
}
```

---

## Relationships

### Many-to-Many

```go
// Model with relationship
type Event struct {
    ID     string  `gorm:"primaryKey"`
    Skills []Skill `gorm:"many2many:event_skills"`
}

type Skill struct {
    ID   string `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex"`
}

// Query with preload
r.db.Preload("Skills").Find(&events)

// Associate skills
r.db.Model(&event).Association("Skills").Append(&skills)

// Replace all skills
r.db.Model(&event).Association("Skills").Replace(&skills)

// Remove association
r.db.Model(&event).Association("Skills").Delete(&skill)

// Clear all
r.db.Model(&event).Association("Skills").Clear()
```

### One-to-Many

```go
type Burst struct {
    ID     string  `gorm:"primaryKey"`
    Events []Event `gorm:"foreignKey:BurstID"`
}

type Event struct {
    ID      string `gorm:"primaryKey"`
    BurstID string `gorm:"index"`
}

// Preload
r.db.Preload("Events").Find(&bursts)
```

---

## Database Setup

### Connection

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewDatabase(path string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Enable foreign keys for SQLite
    if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }
    
    return db, nil
}
```

### Migrations with Goose

```go
// internal/repository/career/migrator.go
package career

import (
    "database/sql"
    "embed"
    
    "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
    goose.SetBaseFS(embedMigrations)
    
    if err := goose.SetDialect("sqlite3"); err != nil {
        return fmt.Errorf("failed to set dialect: %w", err)
    }
    
    if err := goose.Up(db, "migrations"); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }
    
    return nil
}
```

### Migration Files

```sql
-- migrations/001_create_events.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    date DATETIME NOT NULL,
    company TEXT,
    project TEXT,
    role TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_company ON events(company);

-- +goose Down
DROP TABLE IF EXISTS events;
```

---

## Memory Implementation

```go
// internal/repository/career/memory/event_repository.go
package memory

type EventRepository struct {
    mu     sync.RWMutex
    events map[string]*career.Event
}

func NewEventRepository() *EventRepository {
    return &EventRepository{
        events: make(map[string]*career.Event),
    }
}

func (r *EventRepository) Save(ctx context.Context, event *career.Event) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Deep copy to avoid external mutations
    r.events[event.ID] = event.Clone()
    return nil
}

func (r *EventRepository) FindByID(ctx context.Context, id string) (*career.Event, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    event, exists := r.events[id]
    if !exists {
        return nil, repo.ErrEventNotFound
    }
    
    return event.Clone(), nil
}
```

---

## Filter Pattern

```go
// internal/repository/career/filters.go
package career

import "time"

// EventListFilters defines query parameters for event listing.
type EventListFilters struct {
    // Search
    SearchText string
    
    // Filters
    Company   string
    Project   string
    Role      string
    StartDate time.Time
    EndDate   time.Time
    
    // Sorting
    OrderBy   string
    OrderDesc bool
    
    // Pagination
    Limit  int
    Offset int
}

// HasFilters returns true if any filter is set.
func (f EventListFilters) HasFilters() bool {
    return f.SearchText != "" ||
        f.Company != "" ||
        f.Project != "" ||
        f.Role != "" ||
        !f.StartDate.IsZero() ||
        !f.EndDate.IsZero()
}
```

---

## Testing Repositories

### With In-Memory SQLite

```go
var _ = Describe("EventRepository", func() {
    var (
        db   *gorm.DB
        repo *sql.EventRepository
        ctx  context.Context
    )
    
    BeforeEach(func() {
        var err error
        db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
        Expect(err).ToNot(HaveOccurred())
        
        // Run migrations
        sqlDB, _ := db.DB()
        err = career.RunMigrations(sqlDB)
        Expect(err).ToNot(HaveOccurred())
        
        repo = sql.NewEventRepository(db)
        ctx = context.Background()
    })
    
    AfterEach(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })
    
    Describe("Save", func() {
        It("persists the event", func() {
            event := fixtures.Event(1)
            
            err := repo.Save(ctx, event)
            
            Expect(err).ToNot(HaveOccurred())
            
            found, err := repo.FindByID(ctx, event.ID)
            Expect(err).ToNot(HaveOccurred())
            Expect(found.Text).To(Equal(event.Text))
        })
    })
})
```

---

## Anti-Patterns

### DON'T: Return GORM Models

```go
// WRONG - Leaks implementation details
func (r *Repository) FindByID(ctx context.Context, id string) (*models.Event, error)

// CORRECT - Return domain entities
func (r *Repository) FindByID(ctx context.Context, id string) (*career.Event, error)
```

### DON'T: Skip Context

```go
// WRONG - No context propagation
r.db.First(&model, id)

// CORRECT - Always use context
r.db.WithContext(ctx).First(&model, id)
```

### DON'T: Ignore Errors

```go
// WRONG - Silent failure
r.db.Save(model)

// CORRECT - Check and wrap errors
if err := r.db.Save(model).Error; err != nil {
    return fmt.Errorf("failed to save: %w", err)
}
```

### DON'T: Hardcode SQL

```go
// WRONG - Raw SQL for simple queries
r.db.Raw("SELECT * FROM events WHERE id = ?", id)

// CORRECT - Use GORM query builder
r.db.Where("id = ?", id).First(&model)
```

---

## Related Skills

- `domain-modeling` - Domain entity patterns
- `error-handling` - Sentinel errors and wrapping
- `ginkgo-gomega` - Testing repositories
- `db-operations` - Broader database patterns
