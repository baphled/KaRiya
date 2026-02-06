---
name: db-operations
description: Database operations following KaRiya repository patterns with GORM and SQLite
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide database operations following KaRiya's repository pattern with GORM and SQLite.

## When to use me

Use this skill when:
- Creating or modifying database schemas
- Implementing repository methods
- Writing database migrations
- Troubleshooting database issues

## Repository Architecture

```
internal/repository/
├── career/
│   ├── event_repository.go      # Interface definition
│   ├── skill_repository.go      # Interface definition
│   ├── fact_repository.go       # Interface definition
│   ├── burst_repository.go      # Interface definition
│   ├── migrator.go              # Migration runner
│   ├── migrations/              # SQL migration files
│   │   ├── 001_create_career_events.sql
│   │   ├── 002_add_categories_column.sql
│   │   └── ...
│   ├── sql/                     # GORM implementations
│   │   ├── event_repository.go
│   │   ├── skill_repository.go
│   │   └── ...
│   └── memory/                  # In-memory implementations (testing)
│       ├── event_repository.go
│       └── ...
└── models/                      # GORM models
    └── career.go
```

## Creating a New Repository

### 1. Define the Interface

```go
// internal/repository/career/thing_repository.go

//go:generate mockgen -destination=../../testutil/mocks/repository/thing_repository_mock.go -package=mockrepo github.com/baphled/kariya/internal/repository/career ThingRepository

package career

import (
    "context"
    "errors"
    
    "github.com/baphled/kariya/internal/domain/career"
)

var (
    ErrThingNotFound = errors.New("thing not found")
)

// ThingRepository defines the interface for thing persistence.
type ThingRepository interface {
    Create(ctx context.Context, thing *career.Thing) error
    GetByID(ctx context.Context, id string) (*career.Thing, error)
    Update(ctx context.Context, thing *career.Thing) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filters ThingListFilters) ([]*career.Thing, error)
}

type ThingListFilters struct {
    Offset int
    Limit  int
}
```

### 2. Create the GORM Model

```go
// internal/repository/models/career.go

type Thing struct {
    ID        string    `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (Thing) TableName() string {
    return "things"
}

func ThingFromDomain(t *career.Thing) *Thing {
    return &Thing{
        ID:        t.ID,
        Name:      t.Name,
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
    }
}

func (m *Thing) ToDomain() *career.Thing {
    return &career.Thing{
        ID:        m.ID,
        Name:      m.Name,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}
```

### 3. Create Migration

```sql
-- internal/repository/career/migrations/007_create_things.sql

CREATE TABLE IF NOT EXISTS things (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_things_name ON things(name);
```

### 4. Implement SQL Repository

```go
// internal/repository/career/sql/thing_repository.go

package sql

import (
    "context"
    "errors"
    
    "github.com/baphled/kariya/internal/domain/career"
    career_repo "github.com/baphled/kariya/internal/repository/career"
    "github.com/baphled/kariya/internal/repository/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

var _ career_repo.ThingRepository = (*ThingRepository)(nil)

type ThingRepository struct {
    db *gorm.DB
}

func NewThingRepository(db *gorm.DB) *ThingRepository {
    return &ThingRepository{db: db}
}

func (r *ThingRepository) Create(ctx context.Context, thing *career.Thing) error {
    if thing.ID == "" {
        thing.ID = uuid.New().String()
    }
    
    now := time.Now()
    thing.CreatedAt = now
    thing.UpdatedAt = now
    
    model := models.ThingFromDomain(thing)
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *ThingRepository) GetByID(ctx context.Context, id string) (*career.Thing, error) {
    var model models.Thing
    err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, career_repo.ErrThingNotFound
    }
    if err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}
```

### 5. Implement Memory Repository (for testing)

```go
// internal/repository/career/memory/thing_repository.go

package memory

type ThingRepository struct {
    mu     sync.RWMutex
    things map[string]*career.Thing
}

func NewThingRepository() *ThingRepository {
    return &ThingRepository{
        things: make(map[string]*career.Thing),
    }
}
```

### 6. Generate Mocks

```bash
make generate-mocks
```

## Migration Best Practices

1. **Always use `IF NOT EXISTS`** for tables and indexes
2. **Never modify existing migrations** - create new ones
3. **Use descriptive names**: `007_add_status_to_things.sql`
4. **Include indexes** for frequently queried columns
5. **Keep migrations small** - one logical change per file

## Query Patterns

### Filtering with GORM

```go
func (r *ThingRepository) List(ctx context.Context, filters ThingListFilters) ([]*career.Thing, error) {
    query := r.db.WithContext(ctx).Model(&models.Thing{})
    
    if filters.Name != "" {
        query = query.Where("name LIKE ?", "%"+filters.Name+"%")
    }
    
    if filters.Limit > 0 {
        query = query.Limit(filters.Limit)
    }
    
    if filters.Offset > 0 {
        query = query.Offset(filters.Offset)
    }
    
    var models []models.Thing
    if err := query.Find(&models).Error; err != nil {
        return nil, err
    }
    
    result := make([]*career.Thing, len(models))
    for i, m := range models {
        result[i] = m.ToDomain()
    }
    return result, nil
}
```

### Transactions

```go
func (r *Repository) CreateWithRelations(ctx context.Context, thing *Thing) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(thing).Error; err != nil {
            return err
        }
        for _, rel := range thing.Relations {
            if err := tx.Create(rel).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## Testing Repositories

```go
var _ = Describe("ThingRepository", func() {
    var repo *sql.ThingRepository
    var db *gorm.DB
    
    BeforeEach(func() {
        db = testutil.SetupTestDB()
        repo = sql.NewThingRepository(db)
    })
    
    Describe("Create", func() {
        It("creates a thing with generated ID", func() {
            thing := &career.Thing{Name: "Test"}
            
            err := repo.Create(context.Background(), thing)
            
            Expect(err).NotTo(HaveOccurred())
            Expect(thing.ID).NotTo(BeEmpty())
        })
    })
})
```

## Common Issues

| Issue | Solution |
|-------|----------|
| `UNIQUE constraint failed` | Check for duplicates, use upsert pattern |
| `no such table` | Run migrations: migrator.Migrate() |
| `database is locked` | Use single connection for SQLite |
| Slow queries | Add indexes, check N+1 issues |

## Related skills

- `architecture` - Repository layer in architecture
- `security` - SQL injection prevention
- `code-reviewer` - Review DB code
