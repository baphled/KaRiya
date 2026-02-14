# Task 3: BDD Performance — Shared :memory: DB with Transaction Rollback

**Date**: 2026-02-14
**Status**: IMPLEMENTED — DB optimisation delivered; 120s target partially met

## Summary

Implemented shared `:memory:` SQLite database with per-scenario transaction rollback for BDD test isolation. Eliminates per-scenario DB creation, migration, and cleanup overhead.

## Changes Made

### `features/support/hooks.go`
- Added shared `:memory:` SQLite DB initialised once in `SetTestingT()`
- Migrations run once at suite start (not per scenario)
- Each scenario gets a GORM `Begin()` transaction
- `afterScenario` calls `tx.Rollback()` (~8µs cleanup)
- Removed per-scenario `env.Cleanup()` (which included 100ms Windows sleep)

### `features/support/env.go`
- Added `NewAppEnvFromGormDB(t, gormDB)` factory function
- Creates `TestEnv` from a GORM connection (transaction)
- Builds repos, services, and model without DB creation overhead

## Timing Results

| Metric | Before | After | Delta |
|--------|--------|-------|-------|
| Total wall-clock | 145.8s | 139.1s | **-6.7s (4.6%)** |
| Fast tests (50+) | ~0.10-0.11s each | ~0.00s each | **~100% reduction** |
| Slow form tests (8) | 6-25s each | 6-25s each | No change |
| Pre-existing failure | `Edit_event_metadata_during_review` (13s) | Same | Pre-existing |

### Scenario Counts
- Total scenarios: 65
- Passing: 64
- Failing: 1 (pre-existing: `Edit_event_metadata_during_review`)
- Zero regressions

### `make test` Result
- Ginkgo ran 58 suites in 18.3s — **all pass**

## Analysis: Why Not Under 120s

The 120s target assumed DB overhead was the primary bottleneck. In reality:

1. **DB overhead eliminated**: ~6.7s saved (50+ fast tests × ~0.11s each + 100ms cleanup sleep × scenarios)
2. **TUI timeout bottleneck**: 8 slow tests account for ~130s total due to cursor blink timeouts (500-600ms each, multiple per test)
3. **Remaining gap**: 139s - 120s = 19s, entirely from TUI form interaction timeouts

The slow tests (e.g., `Quick_capture_with_minimal_input` at 24.6s, `Accept_inferred_skill_during_review` at 25.6s) have form processing paths that block on cursor blink timers. These cannot be optimised via DB changes.

### Breakdown of Slow Tests
| Test | Time | Bottleneck |
|------|------|-----------|
| Accept_inferred_skill_during_review | 25.6s | TUI form timeouts |
| Quick_capture_with_minimal_input | 24.6s | TUI form timeouts |
| Edit_suggested_burst_before_accepting | 17.1s | TUI form timeouts |
| Accept_suggested_burst_during_review | 17.0s | TUI form timeouts |
| Edit_event_metadata_during_review | 13.1s | TUI form timeouts (FAILING) |
| Add_event_from_timeline | 11.5s | TUI form timeouts |
| Edit_event_metadata_and_save | 11.1s | TUI form timeouts |
| Edit_burst_description_and_save | 11.0s | TUI form timeouts |

## Architecture

```
SetTestingT (once):
  └── sql.Open("sqlite", ":memory:")
  └── SetMaxOpenConns(1)
  └── RunMigrationsForTests(sqlDB)
  └── NewGormDB(sqlDB) → sharedGormDB

beforeScenario (per scenario):
  └── sharedGormDB.Begin() → tx
  └── NewAppEnvFromGormDB(t, tx) → TestEnv with repos on tx
  └── Store tx and env in context

afterScenario (per scenario):
  └── tx.Rollback() (~8µs)
```

## Recommendations for Hitting 120s

To reach 120s, address TUI cursor blink timeouts:
1. Reduce `executeCmd` timeout from 500ms → 100ms for BDD tests
2. Or use a test-mode flag to disable cursor blink processing
3. Or extract form submission to skip UI interaction entirely

These changes would require modifying step definitions or the `e2e.TestEnv` helpers, which are outside the scope of this task.
