# Task List: Burst & Fact CLI Integration

**Based on**: `BURST_FACTS_INTEGRATION_REPORT.md`
**Status**: ✅ Algorithms working, ❌ Integration missing
**Impact**: High - Users can't access burst suggestions or facts after import/capture

---

## Executive Summary

The burst detection and fact extraction algorithms are fully functional with 100% test coverage, but the CLI doesn't initialize the required repositories or trigger processing. This task list prioritizes fixing integration gaps to deliver user-visible functionality.

**Current State - What's Working ✅**:
- ✅ Domain models (Burst, Fact) with validation
- ✅ Repository interfaces (BurstRepository, FactRepository)
- ✅ Memory repositories (MemoryBurstRepository, MemoryFactRepository)
- ✅ SQLite repositories (SQLiteBurstRepository, SQLiteFactRepository) with schema creation
- ✅ Service methods (SuggestBursts, ExtractFactsFromEvent, SaveFact, etc.)
- ✅ Burst detection algorithm (BurstDetector with 0.6+ confidence)
- ✅ Fact extraction algorithm (Extractor)
- ✅ CLI models (BurstSuggestionModel, FactListModel, FactEditorModel, FactCard, BurstList)
- ✅ App screens defined (BurstSuggestionScreen in app.go)
- ✅ Integration tests written (burst_integration_test.go, view_event_with_facts_integration_test.go)
- ✅ 248 events imported successfully

**Current State - What's Missing ❌**:
- ❌ Service.burstRepo field doesn't exist (only factRepo exists)
- ❌ SetBurstRepository() method doesn't exist
- ❌ Repositories not initialized in cmd/cli/main.go
- ❌ Database connection not shared/reused for burst/fact repos
- ❌ No post-import triggers in internal/cli/importer/service.go
- ❌ ImportResult struct doesn't include burst/fact results
- ❌ No display of burst/fact results after import

**Required Fixes**:
1. Add burstRepo field and SetBurstRepository() to Service (P1)
2. Initialize burst/fact repositories in CLI with shared DB connection (P1)
3. Update ImportResult to include burst/fact results (P2)
4. Trigger burst detection after import (P2)
5. Trigger fact extraction after import (P3)
6. Display results to user (P2-P3)
7. Add CLI commands for re-running extraction (P4)
8. Add interactive UI navigation (already partially done) (P5)

---

## Relevant Files

**Modified (Integration)**:
- `cmd/cli/main.go` - Initialize repositories (P1)
- `internal/cli/importer/service.go` - Add post-import triggers (P2-P3)
- `internal/cli/models/import_review.go` - Display results (P2-P3)

**Created (New Features)**:
- `internal/cli/models/burst_results_screen.go` - Display burst detection results (P2)
- `internal/cli/models/facts_results_screen.go` - Display fact extraction results (P3)
- `internal/cli/cli_burst_fact_service.go` - CLI service wrapper for burst/fact operations (P4)

**Supporting Files**:
- `internal/service/career/service.go` - Verify SetBurstRepository method exists (P1)
- `internal/repository/career/burst_repository.go` - Verify schema creation (P1)
- `internal/repository/career/fact_repository.go` - Verify schema creation (P1)

### Notes

- Test files should follow existing patterns: `*_test.go` alongside implementation
- SQLite tables must be created when repositories are initialized
- All burst/fact operations should log with structured logging
- Error handling should be graceful (allow degradation if burst/fact unavailable)

---

## Tasks

### Priority 1: Repository Initialization & Database Setup

#### 1.0 Initialize Burst & Fact Repositories in CLI
- [x] 1.1 Open `cmd/cli/main.go` and locate the service initialization (~line 113)
- [x] 1.2 After `svc := careerservice.NewService(repo)`, add burst repository initialization
- [x] 1.3 Get the *sql.DB connection from the event repository (added GetDB() method to SQLiteRepository)
- [x] 1.4 Create SQLiteBurstRepository with the shared DB connection
- [x] 1.5 Create SQLiteFactRepository with the shared DB connection
- [x] 1.6 Call `svc.SetFactRepository(factRepo)` to wire fact repository (method already exists)
- [x] 1.7 Add `burstRepo` field to Service struct (in internal/service/career/service.go)
- [x] 1.8 Add `SetBurstRepository(burstRepo)` method to Service (follow SetFactRepository pattern)
- [x] 1.9 Call `svc.SetBurstRepository(burstRepo)` to wire burst repository
- [x] 1.10 Add error handling for repository initialization failures (log and continue gracefully)
- [x] 1.11 Write tests in `cmd/cli/main_test.go` to verify repositories are initialized (tests pass)
- [x] 1.12 Verify database tables are created by checking SQLite schema after initialization

**Implementation Details**:
```go
// STEP 1: First, we need to refactor repository initialization to share DB connection
// Current approach: SQLiteRepository creates DB internally
// New approach: Extract DB creation, then pass to all repos

var db *sql.DB
if !inMemory {
    // Open database connection (extract from SQLiteRepository constructor)
    db, err = sql.Open("sqlite", dbPath)
    if err != nil {
        fmt.Fprintf(errOut, "Error opening database: %v\n", err)
        return 1
    }

    // Create event repository with DB
    repo, err = career.NewSQLiteRepositoryWithDB(db)
    if err != nil {
        fmt.Fprintf(errOut, "Error initializing event repository: %v\n", err)
        return 1
    }
} else {
    repo = career.NewMemoryRepository()
}

svc := careerservice.NewService(repo)

// STEP 2: Initialize burst/fact repositories if we have a DB
if db != nil {
    // Initialize fact repository
    factRepo, err := career.NewSQLiteFactRepository(db)
    if err != nil {
        // Log warning but continue - facts are optional enhancement
        fmt.Fprintf(errOut, "Warning: Failed to initialize fact repository: %v\n", err)
    } else {
        svc.SetFactRepository(factRepo)
    }

    // Initialize burst repository
    burstRepo, err := career.NewSQLiteBurstRepository(db)
    if err != nil {
        fmt.Fprintf(errOut, "Warning: Failed to initialize burst repository: %v\n", err)
    } else {
        svc.SetBurstRepository(burstRepo)
    }
}
```

**Alternative Approach** (if we don't want to refactor SQLiteRepository):
```go
// Let each repository open its own connection to the same DB file
if !inMemory {
    // Initialize fact repository with same dbPath
    factRepo, err := career.NewSQLiteFactRepositoryWithPath(dbPath)
    if err != nil {
        fmt.Fprintf(errOut, "Warning: Failed to initialize fact repository: %v\n", err)
    } else {
        svc.SetFactRepository(factRepo)
    }

    // Initialize burst repository with same dbPath
    burstRepo, err := career.NewSQLiteBurstRepositoryWithPath(dbPath)
    if err != nil {
        fmt.Fprintf(errOut, "Warning: Failed to initialize burst repository: %v\n", err)
    } else {
        svc.SetBurstRepository(burstRepo)
    }
}
```
```

**Success Criteria**:
- [x] No compilation errors
- [x] Repositories initialize without error
- [x] SQLite tables created (verify with `sqlite3 ~/.kariya/events.db ".schema"`)
- [x] Tests pass: `go test ./cmd/cli -v`

#### 1.1 Verify SetBurstRepository Method Exists
- [x] 2.1 Open `internal/service/career/service.go`
- [x] 2.2 Search for `SetBurstRepository` method
- [x] 2.3 If missing, add method following same pattern as `SetFactRepository`:
  ```go
  func (s *Service) SetBurstRepository(burstRepo repo.BurstRepository) {
      s.burstRepo = burstRepo
  }
  ```
- [x] 2.4 Verify field `burstRepo repo.BurstRepository` exists in Service struct
- [x] 2.5 If missing, add field to Service struct
- [x] 2.6 Write test verifying method sets repository correctly

**Test Example**:
```go
It("sets burst repository", func() {
    burstRepo := &mockBurstRepository{}
    svc.SetBurstRepository(burstRepo)
    // Verify it's set (may need to expose via getter or check internal state)
})
```

**Success Criteria**:
- [x] SetBurstRepository method exists
- [x] burstRepo field exists in Service struct
- [x] Test passes: `go test ./internal/service/career -v -run SetBurstRepository`

---

### Priority 2: Post-Import Burst Detection

#### 2.0 Add Burst Detection Trigger After Import
- [x] 3.1 Open `internal/cli/importer/service.go` and locate `ImportRows` method
- [x] 3.2 After successful event creation loop, add burst detection trigger
- [x] 3.3 Extract event IDs from successfully created events
- [x] 3.4 Call `is.careerService.SuggestBursts(ctx, eventIDs)` to generate suggestions
- [x] 3.5 Store burst suggestions in importer state (add field to ImporterService or return in result)
- [x] 3.6 Log burst detection results with structured logging
- [x] 3.7 Add error handling (log warning if burst detection fails, continue)
- [x] 3.8 Update ImportResult struct to include burst suggestions if not already present
- [x] 3.9 Write tests for burst detection trigger
- [x] 3.10 Test with imported CSV data to verify burst suggestions generated

**Implementation Details**:
```go
// After successful event creation in ImportRows:

// Detect bursts from new events
if len(result.CreatedEvents) > 0 {
    eventIDs := make([]string, len(result.CreatedEvents))
    for i, event := range result.CreatedEvents {
        eventIDs[i] = event.ID
    }

    suggestions, err := is.careerService.SuggestBursts(ctx, eventIDs)
    if err != nil {
        is.logger.Warnf("Failed to suggest bursts: %v", err)
    } else {
        result.BurstSuggestions = suggestions
        is.logger.Infof("Detected %d burst suggestions from %d events",
            len(suggestions), len(result.CreatedEvents))
    }
}
```

**Success Criteria**:
- [ ] No compilation errors
- [ ] Burst suggestions generated after import
- [ ] Suggestions stored in ImportResult
- [ ] Tests pass: `go test ./internal/cli/importer -v`

#### 2.1 Display Burst Detection Results to User
- [ ] 4.1 Modify `handleNonInteractiveImport` in `cmd/cli/main.go` to display burst results
- [ ] 4.2 After import summary, print burst detection summary:
  ```
  === Burst Suggestions ===
  Detected 5 potential bursts:
  - Performance optimization (3 events, confidence: 0.85)
  - Cloud migration (2 events, confidence: 0.72)
  ...
  ```
- [ ] 4.3 If using interactive mode, navigate to burst review screen instead of exiting
- [ ] 4.4 Add option to skip burst review (--skip-burst-detection flag)
- [ ] 4.5 Add option to save burst suggestions for later review (--review-bursts-later flag)
- [ ] 4.6 Write tests for result display

**Implementation Details**:
```go
// In handleNonInteractiveImport, after import completes:

if result.BurstSuggestions != nil && len(result.BurstSuggestions) > 0 {
    fmt.Fprintf(out, "\n=== Burst Suggestions ===\n")
    fmt.Fprintf(out, "Detected %d potential bursts:\n", len(result.BurstSuggestions))
    for i, burst := range result.BurstSuggestions {
        fmt.Fprintf(out, "%d. %s (%d events, confidence: %.2f)\n",
            i+1, burst.Name, len(burst.EventIDs), burst.ConfidenceScore)
    }
    fmt.Fprintf(out, "\nPress 'Enter' to review burst suggestions, or 'Ctrl+C' to exit\n")
}
```

**Success Criteria**:
- [ ] Burst results displayed to user
- [ ] Confidence scores shown
- [ ] Event counts correct
- [ ] Tests pass: `go test ./cmd/cli -v`

---

### Priority 3: Post-Import Fact Extraction

#### 3.0 Add Fact Extraction Trigger After Import
- [ ] 5.1 Open `internal/cli/importer/service.go` and locate ImportRows after burst detection
- [ ] 5.2 After burst suggestions, add fact extraction trigger
- [ ] 5.3 Extract facts from each created event using `is.careerService.ExtractFactsFromEvent`
- [ ] 5.4 Store extracted facts (keyed by event ID or in ImportResult)
- [ ] 5.5 Call `is.careerService.SaveFact(ctx, fact)` to persist facts
- [ ] 5.6 Handle extraction failures gracefully (log warning, continue)
- [ ] 5.7 Count successfully extracted facts
- [ ] 5.8 Update ImportResult to include fact extraction summary
- [ ] 5.9 Write tests for fact extraction trigger
- [ ] 5.10 Test with imported CSV data to verify facts extracted

**Implementation Details**:
```go
// After burst detection in ImportRows:

// Extract facts from new events
factCount := 0
for _, event := range result.CreatedEvents {
    facts, err := is.careerService.ExtractFactsFromEvent(ctx, event)
    if err != nil {
        is.logger.Warnf("Failed to extract facts from event %s: %v", event.ID, err)
        continue
    }

    for _, fact := range facts {
        if err := is.careerService.SaveFact(ctx, fact); err != nil {
            is.logger.Warnf("Failed to save fact: %v", err)
        } else {
            factCount++
        }
    }
}
result.ExtractedFactsCount = factCount
result.FactsByEventID = make(map[string][]*career.Fact) // Store if needed
```

**Success Criteria**:
- [ ] No compilation errors
- [ ] Facts extracted from each event
- [ ] Facts persisted to database
- [ ] Counts tracked correctly
- [ ] Tests pass: `go test ./internal/cli/importer -v`

#### 3.1 Display Fact Extraction Results to User
- [ ] 6.1 Modify `handleNonInteractiveImport` to display fact extraction summary
- [ ] 6.2 After burst summary, print fact extraction summary:
  ```
  === Fact Extraction ===
  Extracted 47 facts from 248 events
  - Technical facts: 22
  - Leadership facts: 15
  - Mentoring facts: 10
  ```
- [ ] 6.3 Show competency breakdown if available
- [ ] 6.4 Add option to review facts interactively (--review-facts flag)
- [ ] 6.5 Add option to save facts for later review (--review-facts-later flag)
- [ ] 6.6 Write tests for result display

**Implementation Details**:
```go
// In handleNonInteractiveImport, after burst results:

if result.ExtractedFactsCount > 0 {
    fmt.Fprintf(out, "\n=== Fact Extraction ===\n")
    fmt.Fprintf(out, "Extracted %d facts from %d events\n",
        result.ExtractedFactsCount, len(result.CreatedEvents))
    fmt.Fprintf(out, "These facts highlight key competencies and achievements\n")
}
```

**Success Criteria**:
- [ ] Fact extraction summary displayed
- [ ] Counts accurate
- [ ] User can see what was extracted
- [ ] Tests pass: `go test ./cmd/cli -v`

---

### Priority 4: CLI Commands for Re-running Extraction

#### 4.0 Add CLI Flags for Burst/Fact Operations
- [ ] 7.1 Open `cmd/cli/main.go` and locate flag definitions
- [ ] 7.2 Add `--detect-bursts` flag to re-run burst detection on all events
- [ ] 7.3 Add `--extract-facts` flag to re-run fact extraction on all events
- [ ] 7.4 Add `--show-bursts` flag to display all existing bursts
- [ ] 7.5 Add `--show-facts` flag to display all existing facts
- [ ] 7.6 Implement handlers for each flag
- [ ] 7.7 Test flag parsing with various combinations
- [ ] 7.8 Write comprehensive tests for all flags

**Flag Definitions**:
```go
detectBursts := flag.Bool("detect-bursts", false, "Re-run burst detection on all events")
extractFacts := flag.Bool("extract-facts", false, "Re-run fact extraction on all events")
showBursts := flag.Bool("show-bursts", false, "Display all existing bursts")
showFacts := flag.Bool("show-facts", false, "Display all existing facts")
```

**Success Criteria**:
- [ ] Flags parse without error
- [ ] Handlers execute correct operations
- [ ] Results display correctly
- [ ] Tests pass: `go test ./cmd/cli -v`

#### 4.1 Implement Burst Detection Handler
- [ ] 8.1 Create handler function `handleDetectBursts(ctx, svc, out, logger)`
- [ ] 8.2 Retrieve all event IDs from repository
- [ ] 8.3 Call `svc.SuggestBursts(ctx, eventIDs)` with all events
- [ ] 8.4 Display suggestions to user with option to confirm/reject each
- [ ] 8.5 Persist confirmed bursts to database
- [ ] 8.6 Handle errors gracefully
- [ ] 8.7 Write tests

**Success Criteria**:
- [ ] Handler executes without error
- [ ] All bursts detected
- [ ] Results displayed
- [ ] Tests pass

#### 4.2 Implement Fact Extraction Handler
- [ ] 9.1 Create handler function `handleExtractFacts(ctx, svc, out, logger)`
- [ ] 9.2 Retrieve all events from repository
- [ ] 9.3 Extract facts from each event
- [ ] 9.4 Count facts by competency
- [ ] 9.5 Display extraction summary
- [ ] 9.6 Offer to save/confirm facts
- [ ] 9.7 Write tests

**Success Criteria**:
- [ ] Handler executes without error
- [ ] Facts extracted from all events
- [ ] Summary displayed
- [ ] Tests pass

#### 4.3 Implement Show Bursts Handler
- [ ] 10.1 Create handler function `handleShowBursts(ctx, svc, out)`
- [ ] 10.2 List all bursts from database
- [ ] 10.3 Display burst details (name, events, competency focus)
- [ ] 10.4 Sort by creation date or event count
- [ ] 10.5 Add option for detailed view (with event text)
- [ ] 10.6 Write tests

**Success Criteria**:
- [ ] All bursts displayed
- [ ] Formatting clear
- [ ] Tests pass

#### 4.4 Implement Show Facts Handler
- [ ] 11.1 Create handler function `handleShowFacts(ctx, svc, out)`
- [ ] 11.2 List all facts from database
- [ ] 11.3 Display fact details (text, competencies, role fit, audience)
- [ ] 11.4 Sort by source event/burst or creation date
- [ ] 11.5 Filter by competency, role fit, or audience (optional)
- [ ] 11.6 Write tests

**Success Criteria**:
- [ ] All facts displayed
- [ ] Details shown clearly
- [ ] Tests pass

---

### Priority 5: Interactive UI for Reviewing Bursts/Facts (Optional Enhancement)

#### 5.0 Create Burst Results Screen
- [ ] 12.1 Create `internal/cli/models/burst_results_screen.go`
- [ ] 12.2 Implement BubbleTea Model for displaying burst detection results
- [ ] 12.3 Show list of detected bursts with confidence scores
- [ ] 12.4 Allow user to expand each burst to see related events
- [ ] 12.5 Implement keyboard navigation (↑/↓, Enter, Space, 'y'/'n')
- [ ] 12.6 Allow user to confirm or reject bursts
- [ ] 12.7 Display confirmation summary at end
- [ ] 12.8 Write comprehensive tests

**Model Structure**:
```go
type BurstResultsModel struct {
    suggestions    []BurstSuggestion
    focusIndex     int
    selectedBursts map[int]bool // Track accepted bursts
    submitted      bool
    service        *careerservice.Service
    ctx            context.Context
}
```

**Success Criteria**:
- [ ] Model compiles and integrates with BubbleTea
- [ ] All bursts displayed with confidence
- [ ] User can accept/reject interactively
- [ ] Tests pass

#### 5.1 Create Facts Results Screen
- [ ] 13.1 Create `internal/cli/models/facts_results_screen.go`
- [ ] 13.2 Implement BubbleTea Model for displaying extracted facts
- [ ] 13.3 Show list of extracted facts with competencies, role fit, audience
- [ ] 13.4 Implement keyboard navigation
- [ ] 13.5 Allow user to confirm, edit, or reject facts
- [ ] 13.6 Display confirmation summary at end
- [ ] 13.7 Write comprehensive tests

**Model Structure**:
```go
type FactsResultsModel struct {
    facts           []*career.Fact
    focusIndex      int
    selectedFacts   map[int]bool // Track accepted facts
    submitted       bool
    service         *careerservice.Service
    ctx             context.Context
}
```

**Success Criteria**:
- [ ] Model compiles and integrates with BubbleTea
- [ ] All facts displayed with details
- [ ] User can accept/reject/edit interactively
- [ ] Tests pass

#### 5.2 Integrate Results Screens into Import Workflow
- [ ] 14.1 Modify app.go to add BurstResultsScreen and FactsResultsScreen
- [ ] 14.2 After import, navigate to burst results screen if bursts detected
- [ ] 14.3 After burst confirmation, navigate to facts results screen if facts extracted
- [ ] 14.4 After facts confirmation, return to home or import complete screen
- [ ] 14.5 Allow user to skip screens (--skip-review flags)
- [ ] 14.6 Write integration tests

**Navigation Flow**:
```
Import → [Validation] → BurstResultsScreen → FactsResultsScreen → Home/Complete
```

**Success Criteria**:
- [ ] Navigation works smoothly
- [ ] Screens integrate properly
- [ ] User can skip if desired
- [ ] Tests pass

---

### Priority 6: Verification & Documentation

#### 6.0 Integration Verification
- [ ] 15.1 Run full integration test with CSV import
- [ ] 15.2 Verify 248 events imported
- [ ] 15.3 Verify burst suggestions generated and displayed
- [ ] 15.4 Verify fact extraction completed and displayed
- [ ] 15.5 Verify database tables created and populated
- [ ] 15.6 Verify results persisted to database
- [ ] 15.7 Test with --skip-burst-detection flag
- [ ] 15.8 Test with --skip-fact-extraction flag
- [ ] 15.9 Test with interactive mode (BubbleTea screens)
- [ ] 15.10 Run all tests: `go test -race ./...`

**Verification Commands**:
```bash
# Import and process
./kariya --import career_entries.csv

# Verify bursts created
sqlite3 ~/.kariya/events.db "SELECT COUNT(*) FROM bursts;"

# Verify facts created
sqlite3 ~/.kariya/events.db "SELECT COUNT(*) FROM facts;"

# List bursts
./kariya --show-bursts

# List facts
./kariya --show-facts
```

**Success Criteria**:
- [ ] 248 events imported
- [ ] ≥5 bursts detected
- [ ] ≥40 facts extracted
- [ ] Database tables populated
- [ ] All tests pass
- [ ] No race conditions

#### 6.1 Update Documentation
- [ ] 16.1 Update README.md with burst/fact feature section
- [ ] 16.2 Update CLI_GUIDE.md with new flags and workflow
- [ ] 16.3 Update BURST_FACT_EXTRACTION_GUIDE.md with import workflow
- [ ] 16.4 Add troubleshooting for common issues
- [ ] 16.5 Update CHANGELOG.md with integration details
- [ ] 16.6 Document expected performance metrics

**Success Criteria**:
- [ ] All docs updated
- [ ] Workflow explained clearly
- [ ] Flags documented
- [ ] Examples provided

---

## Success Metrics

### Before Integration
- ❌ Burst suggestions not generated after import
- ❌ Facts not extracted after import
- ❌ No database tables created
- ❌ No user-visible results

### After Integration (Target)
- ✅ ≥5 bursts detected from 248 imported events
- ✅ ≥40 facts extracted from 248 events
- ✅ Database tables created and populated
- ✅ Results displayed to user
- ✅ User can review/confirm bursts and facts
- ✅ All tests passing (0 race conditions)
- ✅ Documentation complete

---

## Testing Strategy

### Unit Tests Required
- [ ] Repository initialization (1.0-1.1)
- [ ] Burst detection trigger (2.0)
- [ ] Burst result display (2.1)
- [ ] Fact extraction trigger (3.0)
- [ ] Fact result display (3.1)
- [ ] CLI flag parsing (4.0)
- [ ] Flag handlers (4.1-4.4)
- [ ] Screen models (5.0-5.2)

### Integration Tests Required
- [ ] Full import → burst detection → fact extraction flow
- [ ] Skip flags work correctly
- [ ] Interactive screens integrate smoothly
- [ ] Database persistence works
- [ ] Error recovery graceful

### Performance Tests Required
- [ ] Burst detection ≤2s for ≤500 events (verify with 248)
- [ ] Fact extraction ≤1s per event (verify with 248 events)
- [ ] Database queries ≤100ms for 1000+ items

---

## Implementation Order

**Recommended Sequence**:
1. **Start with P1**: Repository initialization (Tasks 1.0-1.1)
   - Unblocks all subsequent work
   - Required for persistence
   - Can test standalone

2. **Then P2**: Burst detection trigger + display (Tasks 2.0-2.1)
   - Delivers user-visible value
   - Can test with imported CSV

3. **Then P3**: Fact extraction + display (Tasks 3.0-3.1)
   - Builds on burst implementation
   - Completes core integration

4. **Optional P4**: CLI commands (Tasks 4.0-4.4)
   - Adds power-user features
   - Can skip if time-constrained

5. **Optional P5**: Interactive UI (Tasks 5.0-5.2)
   - Polish and enhancement
   - Can skip if time-constrained

**Time Estimates**:
- P1: 1-2 hours
- P2: 2-3 hours
- P3: 2-3 hours
- P4: 3-4 hours
- P5: 4-6 hours

**MVP Scope** (P1-P3): 5-8 hours
**Full Scope** (P1-P5): 12-18 hours

---

## Reference Materials

**Related Files**:
- `BURST_FACTS_INTEGRATION_REPORT.md` - Root cause analysis
- `tasks-05-burst-fact-extraction.md` - Feature implementation (Phases 1-5)
- `internal/service/career/burst_fact/detector.go` - Burst detection algorithm
- `internal/service/career/burst_fact/extractor.go` - Fact extraction algorithm

**Similar Patterns in Codebase**:
- `cmd/cli/main.go` - Entry point initialization pattern
- `internal/cli/importer/service.go` - Post-action trigger pattern
- `internal/cli/models/metadata_review.go` - Interactive review screen pattern
- `internal/service/career/service.go` - Service method pattern

---

## Progress Tracking

- [ ] P1: Repository Initialization (2 tasks, 20 subtasks)
- [ ] P2: Burst Detection Integration (2 tasks, 20 subtasks)
- [ ] P3: Fact Extraction Integration (2 tasks, 20 subtasks)
- [ ] P4: CLI Commands (4 tasks, 30 subtasks)
- [ ] P5: Interactive UI (3 tasks, 25 subtasks)
- [ ] P6: Verification & Documentation (2 tasks, 20 subtasks)

**Total**: 15 parent tasks, 135+ subtasks

---

**Document Version**: 1.0
**Created**: 2025-12-31
**Status**: Ready for implementation
**Next Step**: Start with Task 1.0 (Repository Initialization)
**Estimated Total Effort**: 12-18 hours (MVP: 5-8 hours)
**Process Guide**: docs/rules/master-task-prompt.md
