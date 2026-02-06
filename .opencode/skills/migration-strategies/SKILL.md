---
name: migration-strategies
description: Execute migrations safely - database schema changes, data transformations, zero-downtime deployments
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide safe execution of migrations - database schema changes, data transformations, and code migrations with minimal risk and zero downtime where possible.

## When to use me

- Adding/modifying database columns
- Renaming fields or tables
- Data format transformations
- Large-scale refactoring
- Breaking changes to APIs

## Core Principles

1. **Backwards compatible first** - Old code must work with new schema
2. **Forward compatible second** - New code must work with old schema (during rollout)
3. **Small steps** - Break large migrations into smaller, safer steps
4. **Reversible** - Always have a rollback plan
5. **Test thoroughly** - Migrations are risky, test more than usual

## Migration Patterns

### Expand-Contract Pattern

The safest approach for most changes:

```
PHASE 1: EXPAND
- Add new column/table
- Old code continues working
- New code can use either

PHASE 2: MIGRATE
- Backfill data to new structure
- Update code to use new structure
- Both structures work

PHASE 3: CONTRACT
- Remove old column/table
- All code uses new structure
```

### Example: Renaming a Column

```sql
-- WRONG: Direct rename (breaks running code)
ALTER TABLE events RENAME COLUMN name TO title;

-- RIGHT: Expand-Contract

-- Phase 1: Add new column
ALTER TABLE events ADD COLUMN title VARCHAR(255);

-- Phase 2: Backfill
UPDATE events SET title = name WHERE title IS NULL;

-- Phase 2: Code changes (use title, write to both)
func (e *Event) Save() {
    e.Title = e.Name  // Write to both during transition
}

-- Phase 3: Remove old column (after all code updated)
ALTER TABLE events DROP COLUMN name;
```

## Database Migrations

### Migration File Structure

```
migrations/
├── 001_create_events.sql
├── 002_add_event_type.sql
├── 003_add_event_title.sql       # Phase 1: Expand
├── 004_backfill_event_title.sql  # Phase 2: Migrate
└── 005_drop_event_name.sql       # Phase 3: Contract
```

### Safe Migration Practices

```sql
-- Always use IF NOT EXISTS / IF EXISTS
CREATE TABLE IF NOT EXISTS events (...);
ALTER TABLE events ADD COLUMN IF NOT EXISTS title VARCHAR(255);
DROP TABLE IF EXISTS old_events;

-- Add columns as nullable first
ALTER TABLE events ADD COLUMN title VARCHAR(255);  -- Nullable
-- Later, after backfill:
ALTER TABLE events ALTER COLUMN title SET NOT NULL;

-- Add indexes concurrently (PostgreSQL)
CREATE INDEX CONCURRENTLY idx_events_title ON events(title);

-- Never lock tables for long
-- BAD: UPDATE events SET x = y;  (locks entire table)
-- GOOD: Batch updates
UPDATE events SET title = name 
WHERE id IN (SELECT id FROM events WHERE title IS NULL LIMIT 1000);
```

### Goose Migration Example

```go
// migrations/003_add_event_title.go

//go:build goose

package migrations

import (
    "database/sql"
    "github.com/pressly/goose/v3"
)

func init() {
    goose.AddMigration(upAddEventTitle, downAddEventTitle)
}

func upAddEventTitle(tx *sql.Tx) error {
    _, err := tx.Exec(`
        ALTER TABLE events 
        ADD COLUMN IF NOT EXISTS title VARCHAR(255)
    `)
    return err
}

func downAddEventTitle(tx *sql.Tx) error {
    _, err := tx.Exec(`
        ALTER TABLE events 
        DROP COLUMN IF EXISTS title
    `)
    return err
}
```

## Data Migrations

### Backfill Strategies

```go
// Small dataset: Single transaction
func backfillTitles(db *sql.DB) error {
    _, err := db.Exec(`UPDATE events SET title = name WHERE title IS NULL`)
    return err
}

// Large dataset: Batched
func backfillTitlesBatched(db *sql.DB) error {
    batchSize := 1000
    for {
        result, err := db.Exec(`
            UPDATE events 
            SET title = name 
            WHERE id IN (
                SELECT id FROM events 
                WHERE title IS NULL 
                LIMIT $1
            )
        `, batchSize)
        if err != nil {
            return err
        }
        
        affected, _ := result.RowsAffected()
        if affected == 0 {
            break  // Done
        }
        
        log.Printf("Backfilled %d events", affected)
        time.Sleep(100 * time.Millisecond)  // Reduce load
    }
    return nil
}

// Very large dataset: Background job
func scheduleBackfill() {
    // Run in background, track progress
    // Can be paused/resumed
    // Doesn't block deployment
}
```

### Data Transformation

```go
// Transform data format
func migrateEventDates(db *sql.DB) error {
    // Old: date stored as string "2024-03-15"
    // New: date stored as timestamp
    
    rows, err := db.Query(`
        SELECT id, date_string 
        FROM events 
        WHERE date_timestamp IS NULL
    `)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    for rows.Next() {
        var id string
        var dateStr string
        rows.Scan(&id, &dateStr)
        
        parsed, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            log.Printf("Failed to parse date for event %s: %v", id, err)
            continue  // Don't fail entire migration
        }
        
        db.Exec(`
            UPDATE events 
            SET date_timestamp = $1 
            WHERE id = $2
        `, parsed, id)
    }
    
    return nil
}
```

## Code Migrations

### Feature Flags for Gradual Rollout

```go
// Use feature flag during migration
func (s *Service) GetEvent(id string) (*Event, error) {
    event, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    if s.features.UseNewEventFormat {
        return s.transformToNewFormat(event), nil
    }
    
    return event, nil
}
```

### Parallel Running

```go
// Run both old and new code, compare results
func (s *Service) GetEvent(id string) (*Event, error) {
    oldResult, oldErr := s.oldImplementation(id)
    newResult, newErr := s.newImplementation(id)
    
    // Log differences
    if !reflect.DeepEqual(oldResult, newResult) {
        s.logger.Warn("migration mismatch",
            "event_id", id,
            "old", oldResult,
            "new", newResult,
        )
    }
    
    // Return old result during verification period
    return oldResult, oldErr
}
```

## Zero-Downtime Deployment

### Rolling Update Requirements

```markdown
## Compatibility Requirements

For zero-downtime deploys, at any moment:
- Old code may run with new schema
- New code may run with old schema
- Both versions may run simultaneously

## Checklist
- [ ] New code handles missing new columns
- [ ] Old code ignores new columns
- [ ] No breaking API changes
- [ ] Database migrations are backwards compatible
```

### Deployment Order

```
1. Deploy database migration (expand)
2. Deploy new code (handles both old and new)
3. Run data migration (backfill)
4. Verify all data migrated
5. Deploy code that requires new structure
6. Remove old structure (contract)
```

## Rollback Strategies

### Database Rollback

```bash
# Check current version
goose status

# Rollback last migration
goose down

# Rollback to specific version
goose down-to 003
```

### Code Rollback

```bash
# Kubernetes rollback
kubectl rollout undo deployment/kariya

# Or redeploy previous version
kubectl set image deployment/kariya kariya=kariya:v1.2.2
```

### When NOT to Rollback

- Data has been transformed (can't un-transform easily)
- New data created with new schema
- External systems updated

In these cases: **roll forward** with a fix.

## Testing Migrations

### Test on Copy of Production Data

```bash
# Create test database from production backup
pg_restore -d kariya_migration_test production_backup.dump

# Run migration
goose -dir migrations up

# Verify
psql kariya_migration_test -c "SELECT COUNT(*) FROM events WHERE title IS NULL"
```

### Migration Test Cases

```go
func TestMigration003(t *testing.T) {
    // Setup: Create database with old schema
    db := setupTestDB(t)
    
    // Insert test data
    db.Exec(`INSERT INTO events (id, name) VALUES ('1', 'Test Event')`)
    
    // Run migration
    goose.Up(db, "migrations")
    
    // Verify
    var title sql.NullString
    db.QueryRow(`SELECT title FROM events WHERE id = '1'`).Scan(&title)
    
    assert.False(t, title.Valid)  // Should be null until backfill
}
```

## Common Migration Mistakes

| Mistake | Problem | Fix |
|---------|---------|-----|
| Direct rename | Breaks running code | Expand-contract |
| Locking table | Downtime | Batched updates |
| No rollback plan | Stuck if fails | Always plan rollback |
| Untested migration | Failures in prod | Test on prod-like data |
| Big bang migration | High risk | Small incremental steps |
| Forgetting indexes | Slow queries | Add indexes (concurrently) |

## Related Skills

- `gorm-repository` - Database operations
- `devops` - Deployment strategies
- `incident-response` - When migrations fail
