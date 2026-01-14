---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# CV Generation Bug Fix - Complete Analysis and Solution

**Date**: 2026-01-04
**Status**: ✅ **FIXED AND VERIFIED**
**Impact**: Critical - CV generation with filters now works correctly

---

## Executive Summary

A critical bug in CV generation was preventing CVs from being generated with event filters. When users tried to generate a CV with category filters (e.g., "technical", "achievement"), the system would return **0 events and 0 facts**, resulting in empty CVs.

**Root Cause**: The SQLite database schema was missing the `categories` column, and the filter logic didn't have a fallback mechanism to handle missing categories.

**Solution**:
1. Added `categories` column to the SQLite schema
2. Updated SQLiteRepository to persist and retrieve categories
3. Modified filter logic to fall back to tags when categories are missing
4. Created migration script for existing databases

---

## Bug Investigation

### Initial Observation
Integration tests revealed:
```
CV Generated Successfully:
  Source Event Count: 0  ❌ (Expected > 0)
  Source Fact Count: 0   ❌ (Expected > 0)
```

### Root Cause Analysis

#### Issue 1: Missing Database Column
```sql
-- Actual schema (missing categories):
CREATE TABLE career_events (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    date DATETIME NOT NULL,
    tags TEXT,
    company TEXT,
    project TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Expected schema (with categories):
CREATE TABLE career_events (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    date DATETIME NOT NULL,
    tags TEXT,
    categories TEXT,  -- ← MISSING!
    company TEXT,
    project TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

#### Issue 2: Filter Logic Bug
In `cv_generation_service.go`, the `eventMatchesFilters()` function was checking:
```go
// Check categories filter
if categories, ok := filters["categories"].([]string); ok {
    found := false
    for _, category := range categories {
        for _, eventCategory := range event.Categories {
            if category == eventCategory {
                found = true
                break
            }
        }
        if found {
            break
        }
    }
    if !found && len(categories) > 0 {
        return false  // ← Rejects ALL events since event.Categories is always empty!
    }
}
```

**Problem**: Since no events had categories populated (column didn't exist), `event.Categories` was always empty, causing all events to be filtered out.

#### Issue 3: Database Schema Mismatch
The code tried to read categories in `sqlite_repository.go`:
```go
var tagString, categoriesString string  // ← Declared but never populated!

// SELECT query didn't include categories:
query := "SELECT id, text, date, tags, company, project, created_at, updated_at FROM career_events"
// Missing: categories column

// Later tried to use the unpopulated variable:
if categoriesString != "" {
    event.Categories = parseCategories(categoriesString)
}
```

---

## Solution Implementation

### 1. Updated SQLite Schema
**File**: `internal/repository/career/sqlite_repository.go`

- Added `categories TEXT` column to CREATE TABLE statement
- Added automatic migration function `migrateAddCategoriesColumn()` to handle existing databases
- Updated all SQL queries to include the categories column

### 2. Fixed SQLiteRepository
**Changes**:

#### Create Method
```go
// Convert categories to comma-separated strings
categoriesString := ""
if len(event.Categories) > 0 {
    categoriesString = event.Categories[0]
    for _, category := range event.Categories[1:] {
        categoriesString += "," + category
    }
}

// Insert with categories column
_, err = r.db.ExecContext(ctx, `
    INSERT INTO career_events
    (id, text, date, tags, categories, company, project, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, event.ID, event.Text, event.Date, tagString, categoriesString, ...)
```

#### List Method
```go
// Updated SELECT to include categories
query := "SELECT id, text, date, tags, categories, company, project, created_at, updated_at FROM career_events WHERE 1=1"

// Updated Scan to populate categoriesString
err := rows.Scan(
    &event.ID,
    &event.Text,
    &event.Date,
    &tagString,
    &categoriesString,  // ← Now properly scanned from query
    &event.Company,
    &event.Project,
    &createdAt,
    &updatedAt,
)

// Parse categories correctly
if categoriesString.Valid && categoriesString.String != "" {
    event.Categories = parseCategories(categoriesString.String)
}
```

#### GetByID and Update Methods
Similar updates applied to all methods that interact with categories.

### 3. Smart Filter Fallback Logic
**File**: `internal/service/career/cv/cv_generation_service.go`

```go
// Check categories filter with fallback to tags
if categories, ok := filters["categories"].([]string); ok && len(categories) > 0 {
    found := false

    // If event has categories, match against them
    if len(event.Categories) > 0 {
        for _, category := range categories {
            for _, eventCategory := range event.Categories {
                if category == eventCategory {
                    found = true
                    break
                }
            }
            if found {
                break
            }
        }
    } else {
        // Fall back to tags if no categories are set
        for _, category := range categories {
            for _, eventTag := range event.Tags {
                if category == eventTag {
                    found = true
                    break
                }
            }
            if found {
                break
            }
        }
    }

    if !found {
        return false
    }
}
```

**Why This Works**:
- Events in the database have `tags` (e.g., "technical", "achievement")
- Categories filter requests match against these same values
- The fallback ensures backward compatibility with existing data

### 4. Migration Script
**File**: `scripts/migrate_add_categories_column.sh`

Automatically:
- Checks if categories column exists
- Adds it if missing
- Verifies the migration
- Provides clear feedback to users

---

## Test Results

### Before Fix
```
CV Generated Successfully:
  Source Event Count: 0
  Source Fact Count: 0

CV with Multiple Audiences Generated:
  Source Event Count: 0
  Source Fact Count: 0

CV with Technical Category Filter:
  Events Used: 0
  Facts Used: 0
```

### After Fix
```
CV Generated Successfully:
  Source Event Count: 28
  Source Fact Count: 596

CV with Multiple Audiences Generated:
  Source Event Count: 28
  Source Fact Count: 596

CV with Technical Category Filter:
  Events Used: 28
  Facts Used: 596
```

### Test Suite Results
- ✅ **210/210** CV Service tests passing
- ✅ **84/84** Repository tests passing
- ✅ **164/164** Career Service tests passing
- ✅ **132/132** Burst Fact tests passing
- ✅ **8/8** Classification tests passing
- ✅ **All integration tests** passing

**Total**: 598 tests passing, 0 failures

---

## Files Modified

1. **`internal/repository/career/sqlite_repository.go`**
   - Added categories column to schema
   - Added migration function
   - Updated all CRUD operations to handle categories

2. **`internal/service/career/cv/cv_generation_service.go`**
   - Fixed filter logic to fall back to tags when categories are missing
   - Improved robustness of category filtering

3. **`internal/service/career/cv/cv_generation_integration_test.go`** (new)
   - Added comprehensive integration tests using real production database
   - Tests for various filtering scenarios
   - Tests for metadata preservation

4. **`scripts/migrate_add_categories_column.sh`** (new)
   - Migration script for existing databases

---

## Migration Guide for Users

### Automatic Migration
The SQLiteRepository now automatically handles the migration:
```go
// This runs automatically when opening the database
err = migrateAddCategoriesColumn(db)
```

### Manual Migration (if needed)
```bash
./scripts/migrate_add_categories_column.sh ~/.kariya/events.db
```

### Backward Compatibility
- ✅ Existing events with no categories work fine
- ✅ Filter logic falls back to tags automatically
- ✅ No data loss
- ✅ No breaking changes to the API

---

## Impact Assessment

### What Was Broken
- CV generation with category filters (returns 0 events)
- Users couldn't generate CVs with "technical", "achievement", "leadership" categories
- The Senior IC, Principal Engineer, and other role-based CV configs didn't work

### What's Fixed
- ✅ CV generation with all filter types now works
- ✅ Events are correctly retrieved based on category/tag filters
- ✅ Facts are included in generated CVs
- ✅ Fallback logic handles legacy data gracefully

### Performance Impact
- Negligible: One additional column in the database
- Migration is automatic and non-blocking
- No impact on existing queries

---

## Lessons Learned

1. **Schema-Code Mismatch**: Always ensure database schema matches code expectations
2. **Graceful Degradation**: Implement fallback logic for missing data
3. **Integration Testing**: Real database tests catch issues that unit tests miss
4. **Migration Planning**: Always provide migration scripts for schema changes

---

## Verification Steps

To verify the fix is working:

```bash
# Run integration tests
ginkgo -v ./internal/service/career/cv/

# Check for events in filtered CV
sqlite3 ~/.kariya/events.db "SELECT COUNT(*) FROM career_events WHERE tags LIKE '%technical%';"
# Should show: 28

# Generate a CV and verify event count
# The CV should now show: Source Event Count: 28
```

---

## Future Improvements

1. Consider populating categories from tags during migration
2. Add UI/CLI tool to bulk-assign categories to existing events
3. Add database schema versioning system
4. Create comprehensive database health check tool

---

**Bug Status**: ✅ **CLOSED - FIXED AND VERIFIED**

All tests passing, production database working correctly, backward compatible.

