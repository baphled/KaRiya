# Spike: :memory: SQLite + GORM Transaction Rollback

**Date**: 2026-02-13
**Status**: VALIDATED - Approach works
**Spike file**: `internal/repository/career/sql/spike_memory_txn_test.go` (deleted after documentation)

## Summary

`:memory:` SQLite with GORM transaction rollback is fully viable for BDD test isolation. All critical behaviours validated.

## Test Results

### 1. Basic Transaction Rollback (PASS)

- Open `:memory:` SQLite with GORM
- Run `RunMigrationsForTests()` (6 goose migrations)
- `gormDB.Begin()` → insert career event via `EventRepository.Create()` → verify exists → `Rollback()`
- Event gone after rollback
- Tables (migrations) remain intact

**Verdict**: Transaction rollback on `:memory:` works identically to on-disk SQLite.

### 2. Nested Transactions / Savepoints (PASS with caveat)

**Caveat**: `outerTx.Begin()` returns `"invalid transaction"` error. This is a GORM limitation — nested `Begin()` calls are not supported.

**Working pattern**: Use `outerTx.Transaction(func(savepointTx *gorm.DB) error { ... })` instead. This creates a proper SQLite SAVEPOINT under the hood.

Validated behaviour:
- Outer transaction inserts event1
- Inner `Transaction()` (savepoint) inserts event2, then returns error (triggers savepoint rollback)
- event1 survives savepoint rollback
- event2 is gone after savepoint rollback
- Outer commit preserves event1, event2 remains gone

**Implication for BDD**: If any repository method internally calls `db.Transaction()`, it will create a savepoint within our outer test transaction. When the outer transaction rolls back, everything (including savepoint-committed data) is rolled back. This is exactly what we want.

### 3. Repository Internal `db.Transaction()` (PASS)

Tested the exact pattern we care about for BDD:
- `gormDB.Begin()` (outer test transaction)
- Inside: `outerTx.Transaction(func(innerTx) { repo.Create(ctx, event) })` (simulating repo method that uses `db.Transaction()`)
- `outerTx.Rollback()`
- All events gone — outer rollback clears everything including nested `Transaction()` commits

**This is the critical result**: Even when repository code internally uses `db.Transaction()` for its own atomicity, wrapping everything in an outer test transaction and rolling back still cleans up all data.

### 4. Timing Comparison

| Setup Method | Average Time | Total (10 iterations) |
|---|---|---|
| `:memory:` + migrations | ~1.8-2.3ms | ~18-23ms |
| On-disk + WAL + migrations | ~2.7-2.8ms | ~27-28ms |
| **Speedup** | **1.2-1.5x** | |
| Txn begin + rollback | **~8µs** | ~84µs |

**Key insight**: The migration cost dominates (~1.8ms for `:memory:`, ~2.8ms for disk). The actual transaction begin+rollback overhead is negligible at ~8µs.

**For BDD optimisation**: The real win is NOT `:memory:` vs disk speed (only 1.2-1.5x). The real win is:
- **Current**: Create DB + migrate per scenario = ~2.8ms × N scenarios
- **Proposed**: Create DB + migrate ONCE in BeforeSuite, then begin/rollback per scenario = ~2.8ms + (8µs × N scenarios)
- With 50 BDD scenarios: 140ms → 2.8ms + 0.4ms = 3.2ms for DB setup

The transaction wrapping pattern eliminates per-scenario migration cost entirely.

## Recommended Architecture for BDD

```go
// BeforeSuite: once
sqlDB, _ := sql.Open("sqlite", ":memory:")
RunMigrationsForTests(sqlDB)
gormDB, _ := gorm.Open(sqlite.New(sqlite.Config{Conn: sqlDB}), &gorm.Config{})

// BeforeEach: per scenario
tx := gormDB.Begin()
repos := NewRepositoriesFromDB(tx) // pass transaction as DB
// ... run scenario against repos ...

// AfterEach: per scenario
tx.Rollback() // instant cleanup, ~8µs
```

## Gotchas

1. **GORM nested Begin() fails**: Use `tx.Transaction(func(inner *gorm.DB) error { ... })` for savepoints, NOT `tx.Begin()`
2. **Single connection required**: `:memory:` SQLite with `SetMaxOpenConns(1)` is required — multiple connections would get separate memory databases
3. **Migration tables persist through rollback**: Schema changes (DDL) in SQLite are auto-committed and not transactional. Our migrations run before the test transaction starts, so this is fine.

## Conclusion

**GO AHEAD with `:memory:` + transaction rollback for BDD test isolation.**

The pattern is proven, fast (~8µs per test cleanup vs ~2.8ms for full DB recreation), and handles all edge cases including repository methods that internally use `db.Transaction()`.
