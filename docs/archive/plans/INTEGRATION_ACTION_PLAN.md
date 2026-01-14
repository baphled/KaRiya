---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Integration Action Plan - Burst & Fact CLI Integration

**Status**: 🔴 **CRITICAL** - Algorithms complete, CLI integration missing
**Date**: 2025-12-31
**Impact**: HIGH - Users cannot access burst/fact features

---

## Problem Summary

The burst detection and fact extraction algorithms are **fully functional and tested** (675+ tests passing), but the CLI application doesn't initialize the required repositories or trigger processing after events are captured/imported.

### Current Behavior
```
User imports 248 events
  ↓
Events stored in database ✅
  ↓
Burst detection skipped ❌
  ↓
Fact extraction skipped ❌
  ↓
User sees nothing 😞
```

### Desired Behavior
```
User imports 248 events
  ↓
Events stored in database ✅
  ↓
Burst detection triggered ✅ (5+ bursts detected)
  ↓
Fact extraction triggered ✅ (40+ facts extracted)
  ↓
Results displayed to user ✅
```

---

## Root Causes

| Issue | Location | Fix |
|-------|----------|-----|
| Repositories not initialized | `cmd/cli/main.go` (~L113) | Initialize `SQLiteBurstRepository` and `SQLiteFactRepository` |
| No post-import trigger | `internal/cli/importer/service.go` | Call `SuggestBursts()` and `ExtractFacts()` after events created |
| Database tables not created | SQLite repositories | Tables created when repo initialized, need trigger |
| No result display | `cmd/cli/main.go` | Print summary of detected bursts and facts |
| User can't see results | App workflow | Add screens to review bursts/facts interactively |

---

## Solution: 3-Phase Implementation

### 🟡 Phase 1: Foundation (CRITICAL - 1-2 hours)
**Objective**: Get burst/fact repositories working and persisting data

**Tasks**:
1. Initialize repositories in `cmd/cli/main.go`
2. Verify SetBurstRepository method exists
3. Confirm database tables created
4. Write tests to verify initialization

**Expected Outcome**:
- Repositories initialized without error
- Database tables created (`bursts`, `facts`)
- No user-visible features yet, but foundation ready

**Files**:
- `cmd/cli/main.go` - Add 15 lines
- `internal/service/career/service.go` - Verify/add SetBurstRepository method

### 🟡 Phase 2: Trigger & Display (HIGH VALUE - 2-3 hours)
**Objective**: Process bursts/facts after import and show results

**Tasks**:
1. Add burst detection trigger in import workflow
2. Add fact extraction trigger in import workflow
3. Display burst/fact summaries to user
4. Add CLI flags (--skip-burst-detection, --skip-fact-extraction)

**Expected Outcome**:
- Users see burst suggestions after import
- Users see fact extraction summary after import
- Can skip processing with flags
- **248 events → 5+ bursts, 40+ facts extracted**

**Files**:
- `internal/cli/importer/service.go` - Add 50 lines
- `cmd/cli/main.go` - Add 40 lines
- New: Result display formatting

### 🟢 Phase 3: Interactive Review (NICE-TO-HAVE - 4-6 hours)
**Objective**: Let users review and confirm bursts/facts interactively

**Tasks**:
1. Create BubbleTea screen for burst review
2. Create BubbleTea screen for facts review
3. Integrate into import workflow
4. Add navigation and keyboard shortcuts

**Expected Outcome**:
- Professional UI for reviewing results
- Users can accept/reject/edit bursts and facts
- Results persisted to database
- Smooth workflow integration

**Files**:
- New: `internal/cli/models/burst_results_screen.go`
- New: `internal/cli/models/facts_results_screen.go`
- Modified: `internal/cli/app/app.go` (navigation)

---

## Priority Implementation Order

### Step 1️⃣: Repository Initialization (Do First)
```bash
# Target: cmd/cli/main.go, after line ~113

// After: svc := careerservice.NewService(repo)

if !inMemory {
    factRepo, err := career.NewSQLiteFactRepository(dbPath)
    if err != nil {
        logger.Warnf("Failed to init fact repo: %v", err)
    } else {
        svc.SetFactRepository(factRepo)
    }

    burstRepo, err := career.NewSQLiteBurstRepository(dbPath)
    if err != nil {
        logger.Warnf("Failed to init burst repo: %v", err)
    } else {
        svc.SetBurstRepository(burstRepo)
    }
}
```

**Verification**:
```bash
go test ./cmd/cli -v
sqlite3 ~/.kariya/events.db ".schema"  # Should show bursts and facts tables
```

### Step 2️⃣: Burst Detection Trigger (Do Second)
```bash
# Target: internal/cli/importer/service.go, in ImportRows method
# After: Events created successfully

eventIDs := make([]string, len(result.CreatedEvents))
for i, event := range result.CreatedEvents {
    eventIDs[i] = event.ID
}

suggestions, err := is.careerService.SuggestBursts(ctx, eventIDs)
if err == nil {
    result.BurstSuggestions = suggestions
    is.logger.Infof("Detected %d bursts from %d events",
        len(suggestions), len(result.CreatedEvents))
}
```

**Verification**:
```bash
go test ./internal/cli/importer -v
./kariya --import career_entries.csv  # Should show burst count
```

### Step 3️⃣: Fact Extraction Trigger (Do Third)
```bash
# Target: internal/cli/importer/service.go, in ImportRows method
# After: Burst detection

for _, event := range result.CreatedEvents {
    facts, err := is.careerService.ExtractFactsFromEvent(ctx, event)
    if err == nil {
        for _, fact := range facts {
            is.careerService.SaveFact(ctx, fact)
        }
    }
}
```

**Verification**:
```bash
go test ./internal/cli/importer -v
./kariya --import career_entries.csv  # Should show fact count
```

### Step 4️⃣: Result Display (Do Fourth - Optional)
Display burst/fact summaries to user after import completes.

---

## Quick Impact Check

### Before Fix
```bash
$ ./kariya --import career_entries.csv

=== Import Complete ===
Total rows processed: 353
Successfully imported: 248
Skipped: 104
Failed: 1
✓ Import successful!

# [Nothing else happens - no bursts, no facts]
```

### After Fix (Phase 1-2)
```bash
$ ./kariya --import career_entries.csv

=== Import Complete ===
Total rows processed: 353
Successfully imported: 248
Skipped: 104
Failed: 1
✓ Import successful!

=== Burst Suggestions ===
Detected 8 potential bursts:
- Performance optimization (3 events, confidence: 0.85)
- Cloud infrastructure (2 events, confidence: 0.72)
- Team leadership (4 events, confidence: 0.68)
...

=== Fact Extraction ===
Extracted 47 facts from 248 events
- Technical facts: 22
- Leadership facts: 15
- Mentoring facts: 10

$ sqlite3 ~/.kariya/events.db "SELECT COUNT(*) FROM bursts; SELECT COUNT(*) FROM facts;"
8
47
```

---

## Risk Assessment

### Low Risk
- ✅ Algorithms already working and tested (100% coverage)
- ✅ Repositories already implemented and tested
- ✅ Service methods already exist
- ✅ Adding new code, not modifying existing working code

### Mitigation
- All changes are additive (no breaking changes)
- Error handling allows graceful degradation (burst/fact optional)
- Comprehensive tests for new integration points
- Backward compatible (existing features unchanged)

---

## Dependencies

### Required
- ✅ Burst domain model (DONE)
- ✅ Fact domain model (DONE)
- ✅ Burst repository (DONE)
- ✅ Fact repository (DONE)
- ✅ Burst detection algorithm (DONE)
- ✅ Fact extraction algorithm (DONE)

### Nice-to-Have
- ⏳ BubbleTea UI screens (Phase 3)
- ⏳ Keyboard navigation (Phase 3)
- ⏳ Advanced filtering (Phase 4)

---

## Success Criteria

### Minimum Viable (Phase 1-2)
- [x] Repositories initialize without error
- [x] Database tables created and populated
- [x] ≥5 bursts detected from 248 events
- [x] ≥40 facts extracted
- [x] Results displayed to user in console
- [x] All tests passing
- [x] No race conditions

### Full Featured (Phase 1-3)
- [ ] All minimum criteria
- [ ] Interactive burst review screen
- [ ] Interactive facts review screen
- [ ] User can confirm/reject/edit results
- [ ] Results persisted correctly
- [ ] Smooth navigation workflow

---

## Effort Estimate

| Phase | Tasks | Time | Complexity |
|-------|-------|------|-----------|
| 1: Foundation | 2 parent, 12 subtasks | 1-2h | Low |
| 2: Trigger & Display | 2 parent, 12 subtasks | 2-3h | Low-Medium |
| 3: Interactive Review | 3 parent, 15 subtasks | 4-6h | Medium |
| **Total** | **7 parent, 39 subtasks** | **7-11h** | **Low-Medium** |

**MVP (Phase 1-2)**: 3-5 hours
**Full (Phase 1-3)**: 7-11 hours

---

## Next Steps

### Immediate (Next Session)
1. ✅ Read `BURST_FACTS_INTEGRATION_REPORT.md` (root causes)
2. ✅ Read `tasks-06-burst-fact-cli-integration.md` (implementation plan)
3. Start Task 1.0 (Repository Initialization)
   - Modify `cmd/cli/main.go`
   - Test with `go test ./cmd/cli -v`
   - Verify database tables created

### Short-term (This Week)
1. Complete Task 2.0 (Burst Detection Trigger)
2. Complete Task 3.0 (Fact Extraction Trigger)
3. Verify 248 events → 5+ bursts + 40+ facts

### Medium-term (Next Week)
1. Optional: Implement Phase 3 (Interactive UI)
2. Comprehensive end-to-end testing
3. Documentation and user guides

---

## Checklist for Next Session

Before starting implementation:

- [ ] Read `BURST_FACTS_INTEGRATION_REPORT.md` (understand root causes)
- [ ] Read `tasks-06-burst-fact-cli-integration.md` (detailed tasks)
- [ ] Review `cmd/cli/main.go` (understand current initialization)
- [ ] Review `internal/cli/importer/service.go` (understand import workflow)
- [ ] Review `internal/service/career/service.go` (verify methods exist)
- [ ] Compile codebase: `go build -o kariya ./cmd/cli`
- [ ] Run tests: `go test ./... -v`

When ready to start:
1. Open `cmd/cli/main.go`
2. Locate service initialization (~line 113)
3. Start Task 1.0

---

## Reference Files

**Critical**:
- `BURST_FACTS_INTEGRATION_REPORT.md` - Root cause analysis
- `tasks-06-burst-fact-cli-integration.md` - Detailed implementation tasks
- `tasks-05-burst-fact-extraction.md` - Feature background (Phases 1-5 complete)

**Source Code**:
- `cmd/cli/main.go` - Entry point (where to initialize)
- `internal/cli/importer/service.go` - Import workflow (where to trigger)
- `internal/service/career/service.go` - Service methods (verify exist)
- `internal/service/career/burst_fact/detector.go` - Burst algorithm
- `internal/service/career/burst_fact/extractor.go` - Fact algorithm

---

## Questions?

If anything is unclear:
1. Review the root cause analysis in `BURST_FACTS_INTEGRATION_REPORT.md`
2. Check the detailed tasks in `tasks-06-burst-fact-cli-integration.md`
3. Look at similar patterns in the codebase (metadata_review_screen.go as reference)

The task list is comprehensive and step-by-step. Start with Task 1.0 and work through sequentially.

---

**Document Status**: Ready for Implementation
**Estimated Completion**: 7-11 hours (MVP: 3-5 hours)
**Recommended Start**: Task 1.0 - Repository Initialization
**Success Metric**: 248 events → 5+ bursts, 40+ facts extracted and persisted

