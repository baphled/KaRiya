# Task 35: Test Coverage and Quality Improvements

**Status**: In Progress (Phases 1-6 ✅ complete, Phases 7-11 pending)
**Priority**: MEDIUM-HIGH
**Estimated Time**: 4-5 days total
**Current Coverage**: 80.34%
**Target Coverage**: 85%+

---

## Overview

This task consolidates test coverage improvements, architectural refinements, and code quality enhancements identified across multiple reviews:

1. **PR #74 Analysis**: Test patterns, mocks, edge cases (COMPLETE)
2. **PR #72 Review**: Architectural improvements for intent routing and modals
3. **Coverage Gaps**: Critical workflows with 0% coverage

---

## Current Situation (Updated 2026-01-11)

### Coverage Stats
| Package | Coverage | Status |
|---------|----------|--------|
| **Overall** | 80.34% | ✅ Meets 80% threshold |
| Intent Framework | 88.1% | ✅ Good |
| GlobalContext | 100% | ✅ Excellent |
| Domain Models | >95% | ✅ Excellent |
| Repository | 52.6% | ⚠️ Needs improvement |
| Service | >85% | ✅ Good |
| CV Service | 100% | ✅ Excellent |
| Intents | 61.6% | ⚠️ Needs improvement |

### Key Gaps Identified
1. **Flaky timer test** in `contract_test.go` (uses `time.Sleep`)
2. **CaptureEvent submit workflow** - 0% coverage on core functions
3. **Delete event workflow** - 0% coverage
4. **Modal lifecycle** - Update/View/Result methods untested
5. **Export formats** - Marshal functions untested
6. **Intent registration** - Temporary patterns need cleanup
7. **Delete confirmation** - Uses inline view instead of modal

---

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

---

## PART A: Test Coverage Improvements (from PR #74 analysis)

### Phase 1: Fix Flaky Timer Test ✅ COMPLETE

**Issue**: Test uses `time.Sleep(3500ms)` - slow and potentially flaky.
**File**: `internal/cli/intents/contract_test.go:147`
**Status**: ✅ Fixed - now uses `Eventually` pattern

---

### Phase 2: CaptureEvent Submit Workflow Tests ✅ COMPLETE

**Issue**: Core user journey (event submission) has 0% coverage.
**File**: `internal/cli/intents/capture_event_intent.go`
**Status**: ✅ Fixed - created `capture_event_submit_test.go`

**Subtasks**:
- [x] Add submit workflow tests (TDD)
- [x] Add enrichment workflow tests (TDD)
- [x] Add accept/reject tests (TDD)

**Files Created**: `internal/cli/intents/capture_event_submit_test.go`

---

### Phase 3: Delete Event Workflow Tests ✅ COMPLETE

**Issue**: Delete event has 0% coverage.
**File**: `internal/cli/intents/browse_timeline_intent.go`
**Status**: ✅ Fixed - created `browse_timeline_delete_test.go`

**Subtasks**:
- [x] Add delete confirmation tests (TDD)
- [x] Add event removal tests (TDD)

**Files Created**: `internal/cli/intents/browse_timeline_delete_test.go`

---

### Phase 4: Modal Lifecycle Tests ✅ COMPLETE

**Issue**: Modal Update/View/Result methods have 0% coverage.
**File**: `internal/cli/intents/modals.go`
**Status**: ✅ Fixed - created `modals_lifecycle_test.go`

**Subtasks**:
- [x] Add EditMetadataModal lifecycle tests (TDD)
- [x] Add EditBurstModal lifecycle tests (TDD)
- [x] Add EditFactModal lifecycle tests (TDD)

**Files Created**: `internal/cli/intents/modals_lifecycle_test.go`

---

### Phase 5: Export Format Tests ✅ COMPLETE

**Issue**: Export marshal functions have 0% coverage.
**File**: `internal/cli/intents/export_artifact.go`
**Status**: ✅ Fixed - created `export_formats_test.go` with 66 tests

**Tests Added**:
- marshalToJSON (3 tests)
- marshalToYAML (3 tests)
- marshalEventsToCSV (6 tests)
- marshalEventsToText (14 tests)
- marshalFactsToCSV (5 tests)
- marshalFactsToText (12 tests)
- marshalBurstsToCSV (5 tests)
- marshalBurstsToText (12 tests)
- escapeCSV (6 tests)

**Files Created**: `internal/cli/intents/export_formats_test.go`

---

### Phase 6: Consolidate Test Fixtures ✅ COMPLETE

**Issue**: Test data creation duplicated across files.
**Status**: ✅ Fixed - created fixtures package with factory-go + gofakeit

**Libraries Added**:
- `github.com/bluele/factory-go` - Factory pattern (like FactoryBot)
- `github.com/brianvoe/gofakeit/v7` - Realistic fake data (300+ generators)

**Features Implemented**:
- `EventFactory` - Creates events with realistic text, companies, projects
- `BurstFactory` - Creates bursts with sequential IDs, linked events
- `FactFactory` - Creates facts with valid role fits, audiences, signals
- Quick helpers: `Event()`, `Burst()`, `Fact()`, `BurstConfirmed()`, `FactFromBurst()`
- Batch helpers: `Events(n)`, `Bursts(n, events)`, `Facts(n, events)`
- `SetSeed()` for reproducible tests

**Files Created**:
- `internal/testutil/fixtures/factories.go` - Factory implementations
- `internal/testutil/fixtures/factories_test.go` - 26 tests

**Usage Examples**:
```go
// Factory with fake data
event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

// With overrides
event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
    "Company": "TechCorp",
}).(*career.CareerEvent)

// Quick minimal objects
event := fixtures.Event("my-id")
burst := fixtures.Burst("burst-id", "evt-1", "evt-2")

// Batch creation
events := fixtures.Events(10)
facts := fixtures.Facts(5, events)
```

---

## PART B: Architectural Improvements (from PR #72 review)

### Phase 7: Intent Registration Refactoring (Priority: MEDIUM)

**Issue**: Temporary intent registration creates new intent names dynamically.
**File**: `internal/cli/app/app.go:200`

**Options to Evaluate**:
1. Edit flag approach: Reuse `capture_event` with `isEditMode` flag
2. Dedicated edit intent: Register `capture_event_edit` once at startup
3. Context-based routing: Router handles edit context automatically

**Subtasks**:
- [ ] Document chosen approach with pros/cons
- [ ] Implement new registration pattern (TDD)
- [ ] Remove temporary registration code

**Files to Modify**:
- `internal/cli/app/app.go`
- `internal/cli/intents/capture_event_intent.go`
- `internal/cli/intents/router.go`

**Estimated Time**: 4-6 hours

---

### Phase 8: Delete Confirmation Modal (Priority: HIGH)

**Issue**: Delete confirmation uses inline view instead of modal.
**Impact**: Less prominent warning, risk of accidental deletions.

**Subtasks**:
- [ ] Create DeleteConfirmationModal using existing modal system (TDD)
- [ ] Integrate modal into BrowseTimeline (TDD)
- [ ] Apply pattern to BurstManagement and FactManagement (TDD)

**Files to Modify**:
- `internal/cli/intents/modals.go` - Add DeleteConfirmationModal
- `internal/cli/intents/browse_timeline_intent.go` - Use modal
- `internal/cli/intents/burst_management_intent.go` - Add delete modal
- `internal/cli/intents/fact_management_intent.go` - Add delete modal

**Estimated Time**: 6-8 hours

---

### Phase 9: Error Modal for Delete Failures (Priority: MEDIUM)

**Issue**: Delete errors stored in state but not prominently displayed.
**Goal**: Immediate, clear error feedback using modal system.

**Subtasks**:
- [ ] Show error modal when delete fails (TDD)
- [ ] Add "Retry" and "Cancel" options
- [ ] Include error details and suggested actions

**Files to Modify**:
- `internal/cli/intents/browse_timeline_intent.go`
- `internal/cli/components/modal.go` (if enhancement needed)

**Estimated Time**: 3-4 hours

---

### Phase 10: Form State Management Documentation (Priority: LOW)

**Issue**: Form cancel behavior not explicitly documented or tested.

**Subtasks**:
- [ ] Document HuhCaptureForm lifecycle in code comments
- [ ] Add state diagram to `docs/HUH_FORMS_GUIDE.md`
- [ ] Add form cancel tests (TDD)

**Files to Modify**:
- `internal/cli/models/huh_capture_form.go`
- `docs/HUH_FORMS_GUIDE.md`
- Create `internal/cli/models/huh_capture_form_test.go`

**Estimated Time**: 2-3 hours

---

### Phase 11: Intent Router Enhancement (Priority: LOW)

**Issue**: Message-based routing is ad-hoc, pattern not formalized.
**Goal**: Formalize message-based routing as reusable pattern.

**Subtasks**:
- [ ] Document `RequestEditEventMsg` pattern in `TUI_INTENT_DIAGRAM.md`
- [ ] Create `RouteToIntent()` helper in router (TDD)
- [ ] Use helper for `RequestEditEventMsg`

**Files to Modify**:
- `internal/cli/intents/router.go`
- `internal/cli/app/app.go`
- `docs/TUI_INTENT_DIAGRAM.md`

**Estimated Time**: 3-4 hours

---

## Summary

### Effort Breakdown

| Phase | Description | Priority | Time |
|-------|-------------|----------|------|
| 1 | Fix Flaky Timer | HIGH | ✅ Done |
| 2 | CaptureEvent Submit Tests | HIGH | 2-3 hrs |
| 3 | Delete Event Tests | HIGH | 1-2 hrs |
| 4 | Modal Lifecycle Tests | MEDIUM | 3-4 hrs |
| 5 | Export Format Tests | MEDIUM | 2-3 hrs |
| 6 | Consolidate Fixtures | MEDIUM | 4-5 hrs |
| 7 | Intent Registration | MEDIUM | 4-6 hrs |
| 8 | Delete Confirmation Modal | HIGH | 6-8 hrs |
| 9 | Error Modal for Delete | MEDIUM | 3-4 hrs |
| 10 | Form Documentation | LOW | 2-3 hrs |
| 11 | Router Enhancement | LOW | 3-4 hrs |
| **Total** | | | **32-45 hours** |

### Suggested Order
1. Phase 2-3: High priority test coverage (4-5 hrs)
2. Phase 8: Delete modal - highest user impact (6-8 hrs)
3. Phase 4-5: Medium priority tests (5-7 hrs)
4. Phase 9: Error modal - builds on Phase 8 (3-4 hrs)
5. Phase 7: Intent registration - careful design needed (4-6 hrs)
6. Phase 6: Fixture consolidation - maintenance improvement (4-5 hrs)
7. Phase 10-11: Documentation and router (5-7 hrs)

### Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Overall Coverage | 80.34% | 85%+ |
| Intents Coverage | 61.6% | 75%+ |
| Repository Coverage | 52.6% | 75%+ |
| Flaky Tests | 0 | 0 |
| Delete Uses Modal | No | Yes |
| Intent Registration | Temporary | Clean |

---

## Rollback Plan

- **Test additions (Phases 2-6)**: Delete new test files if issues
- **Modal changes (Phases 8-9)**: Revert to inline view if modal causes issues
- **Intent registration (Phase 7)**: Revert to temporary pattern
- **Documentation (Phases 10-11)**: No rollback needed

---

## Related Documentation
- `docs/TESTING_PATTERNS.md` - Testing patterns guide
- `docs/TUI_INTENT_DIAGRAM.md` - Intent architecture
- `docs/MODAL_PATTERNS.md` - Modal usage patterns
- `docs/HUH_FORMS_GUIDE.md` - Huh forms developer guide
- `docs/rules/master-task-prompt.md` - Development workflow

---

**Document Version**: 2.0 (Consolidated)
**Created**: 2026-01-11
**Last Updated**: 2026-01-11
**Status**: IN PROGRESS
**Related PRs**: #72, #74
