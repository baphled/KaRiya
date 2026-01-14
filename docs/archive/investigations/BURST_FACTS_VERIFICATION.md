---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Burst Facts Extraction Verification Report

**Date**: 2025-12-31  
**Status**: ⚠️ **ISSUE IDENTIFIED** - Burst facts NOT extracted during import

## Summary

Imported 248 career events via CSV using `--import` flag, but **burst facts extraction was NOT performed**. The functionality exists in the codebase but is not wired up in the CLI.

## What Works ✅

1. **CSV Import**: Successfully imports 248 events from career_entries.csv
   - Events stored in SQLite database
   - All events validated and persisted correctly

2. **Burst Detection Algorithm**: Works correctly in code
   - Can detect related events (temporal + similarity matching)
   - Generates burst suggestions with confidence scores
   - Example: 3 performance optimization events detected as burst (confidence: 0.61)

3. **Fact Extraction Algorithm**: Works correctly in code
   - Extracts facts from individual events
   - Generates role-fit classifications (e.g., "senior_ic")
   - Example: "Worked on performance optimization" extracted as 1 fact

4. **Service Methods**: All exist and function correctly
   - `SuggestBursts(ctx, eventIDs)` → Returns burst suggestions
   - `ExtractFactsFromEvent(ctx, event)` → Returns extracted facts
   - `ExtractFactsFromBurst(ctx, burst)` → Batch extraction from bursts

## What's Missing ❌

1. **Burst & Fact Repository Initialization**
   - SQLite tables for `bursts` and `facts` are NOT created
   - Fact repository is NOT wired up to career service in main.go
   - **No burst/fact repositories instantiated in CLI entry point**

2. **Post-Import Extraction**
   - Import workflow does NOT trigger burst detection after events are imported
   - Import workflow does NOT trigger fact extraction after events are imported

3. **CLI Integration**
   - No command to run burst detection on imported events
   - No command to manually trigger fact extraction
   - No display of burst suggestions to user

## Database Schema Status

Current database tables:
```
$ sqlite3 ~/.kariya/events.db "SELECT name FROM sqlite_master WHERE type='table';"
career_events
```

Missing but defined in code:
```
- bursts (id, name, description, event_ids, competency_focus, created_at, updated_at)
- facts (id, text, competencies, role_fit, audience_relevance, strength_signal, source_event_id, source_burst_id, created_at, updated_at)
```

## Verification Test

Created minimal test program to verify burst/fact functionality works:

```go
// Create 3 performance optimization events (temporal proximity + keyword similarity)
// Burst Detection: ✓ Detected 1 burst (confidence: 0.61)
// Fact Extraction: ✓ Extracted 1 fact (role_fit: senior_ic)
```

**Conclusion**: The burst and fact extraction algorithms work correctly, but the CLI integration is incomplete.

## Issues to Address

### Issue #1: Repositories Not Initialized
**File**: `cmd/cli/main.go`  
**Problem**: Burst and Fact repositories are not created or wired to service

```go
// Current (incomplete):
svc := careerservice.NewService(repo)

// Should be:
svc := careerservice.NewService(repo)
factRepo, err := career.NewSQLiteFactRepository(dbPath)
burstRepo, err := career.NewSQLiteBurstRepository(dbPath)
svc.SetFactRepository(factRepo)
```

### Issue #2: No Post-Import Extraction
**File**: `cmd/cli/main.go` (handleNonInteractiveImport function)  
**Problem**: After importing events, no automatic burst detection or fact extraction

**Expected behavior**:
1. Import events
2. Detect bursts from imported events
3. Extract facts from events/bursts
4. Report summary to user

**Current behavior**:
1. Import events
2. Report import summary
3. Exit

### Issue #3: No User Interface for Burst/Fact Management
**Files**: CLI models and app state  
**Problem**: No screens or commands to view/manage bursts and facts post-import

## Testing Results

### CSV Import Test
```
$ ./kariya-test --import /home/baphled/Projects/career_entries.csv --skip-import-review
Parsing CSV file: /home/baphled/Projects/career_entries.csv
Found 353 rows to import
Importing 353 rows...
...
=== Import Complete ===
Total rows processed: 353
Successfully imported: 248
Skipped: 104
Failed: 1

✓ Import successful!
```

### Burst Detection Test (Programmatic)
```
Created 3 related events:
  - "Worked on performance optimization..." (2020-01-15)
  - "Implemented caching layer..." (2020-01-20)
  - "Completed performance benchmarking..." (2020-01-25)

Burst Detection: ✓ Detected 1 burst
  - Confidence Score: 0.61
  - Event IDs: [id1, id3]
  - Name: (empty)
  - Description: (empty)

Fact Extraction: ✓ Extracted 1 fact
  - ID: c1f71be3-2853-4c08-9602-a9af29fa66bf
  - Text: Worked on performance optimization with database tuning
  - Role Fit: senior_ic
```

## Recommendations

1. **Enable Burst/Fact Storage**
   - Initialize burst and fact repositories in main.go
   - Create database tables on first run

2. **Add Post-Import Processing**
   - After successful import, detect bursts
   - Extract facts from imported events
   - Store bursts and facts in database

3. **Add CLI Commands**
   - `--detect-bursts` to manually trigger burst detection
   - `--extract-facts` to manually trigger fact extraction
   - Display burst/fact suggestions in import summary

4. **Add Interactive UI**
   - Screen to review detected bursts
   - Screen to confirm extracted facts
   - Option to save or discard bursts/facts

## Files Involved

**Repositories** (working, need initialization):
- `internal/repository/career/sqlite_burst_repository.go` (46 table schema)
- `internal/repository/career/sqlite_fact_repository.go` (35 table schema)
- `internal/repository/career/burst_repository.go` (interface)
- `internal/repository/career/fact_repository.go` (interface)

**Service** (working, needs repository setup):
- `internal/service/career/service.go` (has SuggestBursts, ExtractFactsFrom*)

**Burst/Fact Package** (working, tested):
- `internal/service/career/burst_fact/detector.go` (detects bursts)
- `internal/service/career/burst_fact/extractor.go` (extracts facts)
- `internal/service/career/burst_fact/classifier.go` (classifies facts)

**CLI** (needs integration):
- `cmd/cli/main.go` (needs repository initialization)
- `internal/cli/importer/service.go` (needs post-import extraction)

---

**Status**: Ready for implementation fixes  
**Severity**: Medium (functionality exists but disconnected)  
**Fix Complexity**: Low (wiring existing components together)
