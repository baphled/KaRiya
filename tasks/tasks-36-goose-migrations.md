# Task 36: Goose Migration Integration

## Overview
- **Goal**: Integrate goose migration framework to provide Rails-like database migrations for KaRiya
- **Time Estimate**: 3-5 hours
- **Prerequisites**: Understanding of current SQLite repository pattern, goose migration framework

## Context

Currently, KaRiya manages database schema through embedded `CREATE TABLE IF NOT EXISTS` statements in repository constructors and manual migration functions (`migrateAddCategoriesColumn`). This approach has limitations:
- No migration versioning or tracking
- No rollback capability
- Manual migration code mixed with repository logic
- No centralized migration history

This task integrates the **goose** migration framework to provide:
- Versioned SQL migration files
- Automatic migration tracking
- Baseline detection for existing databases
- Embedded migrations in binary (single-file distribution)
- Database indexes for improved query performance

## Implementation Strategy

### Migration-Only Approach
After discussion, we're implementing **goose for migrations only** (not a full ORM):
- Keep existing repository pattern (raw SQL queries)
- Add versioned migration files
- Maintain current codebase structure
- Low-risk, incremental improvement

### Key Design Decisions
1. **Migration storage**: Embedded in binary via `//go:embed`
2. **Existing DB handling**: Auto-detect and mark baseline migrations as applied
3. **Rollback support**: Only for additive migrations (SQLite limitation)
4. **CLI commands**: Skip for now (not frequently needed)
5. **Add indexes**: Yes, include performance indexes

## Files to Create

### Migration Files (SQL)
- [ ] `internal/repository/career/migrations/001_create_career_events.sql`
- [ ] `internal/repository/career/migrations/002_add_categories_column.sql`
- [ ] `internal/repository/career/migrations/003_create_bursts.sql`
- [ ] `internal/repository/career/migrations/004_create_facts.sql`

### Go Implementation Files
- [ ] `internal/repository/career/migrator.go` - Goose wrapper with baseline detection
- [ ] `internal/repository/career/migrator_test.go` - Migration tests
- [ ] `internal/testutil/db.go` - Reusable test database helpers
- [ ] `internal/testutil/db_test.go` - Test helper tests

## Files to Modify

### Production Code
- [ ] `go.mod` - Add goose dependency
- [ ] `cmd/cli/main.go` - Run migrations before creating repositories
- [ ] `internal/repository/career/sqlite_repository.go` - Remove embedded schema, add `NewSQLiteRepositoryWithDB`
- [ ] `internal/repository/career/sqlite_burst_repository.go` - Remove `initSchema`, add `NewSQLiteBurstRepositoryWithDB`
- [ ] `internal/repository/career/sqlite_fact_repository.go` - Remove `initSchema`, add `NewSQLiteFactRepositoryWithDB`

### Test Files (using testutil)
- [ ] `internal/repository/career/sqlite_repository_test.go`
- [ ] `internal/repository/career/suite_test.go`
- [ ] `internal/repository/career/sqlite_project_test.go`
- [ ] `internal/repository/career/integration_test.go`
- [ ] `cmd/cli/persistence_test.go`
- [ ] `internal/service/career/integration_test.go`
- [ ] `internal/service/career/cv/cv_generation_integration_test.go`

## Files to Delete
- [ ] `scripts/migrate_add_categories_column.sh` - Obsolete manual migration script

## Implementation Checklist

### Phase 1: Add Goose Infrastructure
- [ ] Add goose dependency: `go get github.com/pressly/goose/v3`
- [ ] Create `internal/repository/career/migrations/` directory
- [ ] Create migration 001: career_events table with indexes (date, company)
- [ ] Create migration 002: add categories column (non-reversible)
- [ ] Create migration 003: bursts table with index (confirmed)
- [ ] Create migration 004: facts table with indexes (source_event_id, source_burst_id)
- [ ] Create `migrator.go` with:
  - `RunMigrations(db *sql.DB) error` - Main migration runner
  - `handleBaselineMigration(db *sql.DB) error` - Auto-detect existing databases
  - `MigrationStatus(db *sql.DB) (int64, error)` - Get current version
  - Helper functions: `hasTable`, `hasColumn`
- [ ] Create `migrator_test.go` with tests for:
  - Fresh database migration
  - Existing database (pre-goose) detection
  - Baseline version detection (categories column exists)
  - Migration status reporting

### Phase 2: Create Test Utilities
- [ ] Create `internal/testutil/` directory
- [ ] Create `db.go` with:
  - `SetupTestDB(t testing.TB) (*sql.DB, func())` - Standard test DB setup
  - `SetupTestDBWithPath(t testing.TB) (string, *sql.DB, func())` - With path for legacy constructors
- [ ] Create `db_test.go` with tests for test helpers

### Phase 3: Integrate Migrations into Application
- [ ] Add `NewSQLiteRepositoryWithDB(db *sql.DB) *SQLiteRepository` to sqlite_repository.go
- [ ] Add `NewSQLiteBurstRepositoryWithDB(db *sql.DB) *SQLiteBurstRepository` to sqlite_burst_repository.go
- [ ] Add `NewSQLiteFactRepositoryWithDB(db *sql.DB) *SQLiteFactRepository` to sqlite_fact_repository.go
- [ ] Update `cmd/cli/main.go`:
  - Open database connection with `sql.Open`
  - Call `career.RunMigrations(db)` before creating repositories
  - Use `NewSQLiteRepositoryWithDB(db)` instead of `NewSQLiteRepository(dbPath)`
  - Same for fact and burst repositories
  - Keep in-memory mode unchanged

### Phase 4: Remove Embedded Schema
- [ ] Remove `CREATE TABLE` statement from `NewSQLiteRepository` (lines 30-45)
- [ ] Remove `migrateAddCategoriesColumn` function (lines 56-94)
- [ ] Remove `migrateAddCategoriesColumn()` call (line 48)
- [ ] Remove `initSchema()` function from sqlite_burst_repository.go (lines 32-49)
- [ ] Remove `initSchema()` call from `NewSQLiteBurstRepository` (line 25)
- [ ] Remove `initSchema()` function from sqlite_fact_repository.go (lines 32-50)
- [ ] Remove `initSchema()` call from `NewSQLiteFactRepository` (line 25)
- [ ] Keep constructors minimal (just assign db connection)

### Phase 5: Update Tests
- [ ] Update `sqlite_repository_test.go`:
  - Import testutil
  - Replace temp dir creation with `testutil.SetupTestDBWithPath`
  - Use `DeferCleanup` for cleanup
- [ ] Update `suite_test.go`: Same pattern
- [ ] Update `sqlite_project_test.go`: Same pattern
- [ ] Update `integration_test.go`: Same pattern
- [ ] Update `cmd/cli/persistence_test.go`: Same pattern
- [ ] Update `internal/service/career/integration_test.go`: Same pattern
- [ ] Update `internal/service/career/cv/cv_generation_integration_test.go`: Same pattern
- [ ] Run tests after each file update to catch issues early

### Phase 6: Cleanup and Documentation
- [ ] Delete `scripts/migrate_add_categories_column.sh`
- [ ] Update `AGENTS.md` "Database Setup" section:
  - Document new migration system
  - Document goose integration
  - Update "How Migrations Are Currently Handled" section
  - Add migration file examples
- [ ] Run `go mod tidy`
- [ ] Run full test suite: `make test`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Run staticcheck: `staticcheck ./...`

## Testing Instructions

### Unit Tests
```bash
# Test migrator
go test -v ./internal/repository/career/migrator_test.go

# Test testutil
go test -v ./internal/testutil/

# Test repositories with migrations
go test -v ./internal/repository/career/
```

### Integration Tests
```bash
# Test full application
go test -v ./cmd/cli/

# Test service layer
go test -v ./internal/service/career/
```

### Manual Testing
```bash
# Test with fresh database
rm -rf ~/.kariya/events.db
go run ./cmd/kariya

# Test with existing database
# (should auto-detect baseline and not recreate tables)
go run ./cmd/kariya
```

### Verification Checklist
- [ ] Fresh database: All 4 tables created with indexes
- [ ] Existing database (with categories): Baseline detected, migrations skipped
- [ ] Existing database (without categories): Migration 002 applied
- [ ] All 2,078+ tests pass
- [ ] No race conditions: `go test -race ./...`
- [ ] Zero staticcheck warnings
- [ ] In-memory mode still works (MemoryRepository unchanged)
- [ ] `MigrationStatus()` returns correct version (1-4)

## Acceptance Criteria

### Functionality
- [ ] Goose migration framework integrated
- [ ] 4 migration files created (001-004)
- [ ] All tables have appropriate indexes
- [ ] Baseline detection works for existing databases
- [ ] Migration status reporting works
- [ ] Test utilities created in `internal/testutil/`

### Code Quality
- [ ] All repository schema code removed (~140 lines)
- [ ] New constructors added (`*WithDB` variants)
- [ ] All tests updated to use testutil
- [ ] All tests pass (100% pass rate maintained)
- [ ] No race conditions detected
- [ ] Zero staticcheck warnings
- [ ] Code coverage maintained (>87%)

### Documentation
- [ ] AGENTS.md updated with new migration system
- [ ] Migration files include up/down comments
- [ ] Migrator code includes comprehensive comments
- [ ] Test helpers documented

### Backward Compatibility
- [ ] Existing databases work without manual intervention
- [ ] No data loss during migration
- [ ] In-memory mode unaffected
- [ ] CLI flags work as before (`--db`, `--in-memory`)

## Rollback Plan

If issues arise during implementation:

### Phase 1-2 Issues
- Delete new files
- Run `go mod tidy`
- No production code affected yet

### Phase 3-4 Issues
- Revert `cmd/cli/main.go` changes
- Keep old constructors alongside new ones
- Repository code still has schema creation (safety net)

### Phase 5 Issues
- Revert test file changes one by one
- Each test file is independent

### Complete Rollback
```bash
git checkout next
git branch -D feat/goose-migrations
```

## Migration File Details

### 001_create_career_events.sql
- CREATE TABLE career_events (all existing columns)
- CREATE INDEX idx_career_events_date
- CREATE INDEX idx_career_events_company
- Uses `IF NOT EXISTS` for idempotency

### 002_add_categories_column.sql
- ALTER TABLE career_events ADD COLUMN categories TEXT
- Non-reversible rollback (SQLite limitation)

### 003_create_bursts.sql
- CREATE TABLE bursts (all existing columns)
- CREATE INDEX idx_bursts_confirmed
- Uses `IF NOT EXISTS` for idempotency

### 004_create_facts.sql
- CREATE TABLE facts (all existing columns)
- CREATE INDEX idx_facts_source_event_id
- CREATE INDEX idx_facts_source_burst_id
- Uses `IF NOT EXISTS` for idempotency

## Notes

### Why Goose Over GORM AutoMigrate?
- Current raw SQL approach working well
- No need for full ORM overhead
- Better separation of concerns
- Explicit migration files (reviewable, trackable)
- Lower learning curve

### Why Embedded Migrations?
- Single binary distribution
- No external files needed for deployment
- Migrations always in sync with code version

### Why Separate testutil Package?
- Reusable across multiple test packages
- Clear separation of test utilities from production code
- Can be expanded for other test helpers later

### Performance Expectations
- Index addition should improve:
  - Timeline queries (date index)
  - Company filtering (company index)
  - Fact-to-event joins (source_event_id index)
- No performance degradation expected from goose itself

## Expected Metrics After Completion

| Metric | Before | After | Change |
|--------|--------|-------|--------|
| Total files | N/A | +9 new | +9 |
| Repository schema code | ~140 lines | 0 lines | -140 |
| Migration tracking | None | goose_db_version table | ✅ |
| Database indexes | 0 | 5 | +5 |
| Test pass rate | 100% | 100% | Maintained |
| Race conditions | 0 | 0 | Maintained |
| Staticcheck warnings | 0 | 0 | Maintained |

## References

- [Goose GitHub](https://github.com/pressly/goose)
- [AGENTS.md - Database Setup](../AGENTS.md#database-setup)
- [docs/rules/master-task-prompt.md](../docs/rules/master-task-prompt.md)
- [docs/rules/go-guidelines.md](../docs/rules/go-guidelines.md)
