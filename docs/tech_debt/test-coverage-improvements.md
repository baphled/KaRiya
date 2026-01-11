# Task 42: Test Coverage and Quality Improvements

## Overview
- **Goal**: Address test coverage gaps, fix flaky tests, and improve test maintainability
- **Time Estimate**: 3-5 days (20-30 hours total)
- **Prerequisites**: PR #74 technical debt complete
- **Priority**: Medium-High

## Context

Analysis of the KaRiya test suite identified several areas for improvement:
- Overall test count: 2,078+ tests (100% pass rate)
- Coverage maintained at 80.34% overall
- Several critical workflows have 0% coverage
- One flaky test identified with timing dependency
- Test fixture duplication across packages

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

---

## Phase 1: Fix Flaky Timer Test (Priority: HIGH)

**Issue**: Test uses `time.Sleep(3500ms)` making it slow and potentially flaky.

**File**: `internal/cli/intents/contract_test.go:147`

```go
It("should expire success after 3 seconds", func() {
    base.SetSuccess("Test message")
    Expect(base.ShouldShowSuccess()).To(BeTrue())
    time.Sleep(3500 * time.Millisecond)  // FLAKY
    Expect(base.ShouldShowSuccess()).To(BeFalse())
})
```

### Subtask 1.1: Refactor Timer Test (TDD)
- [ ] **RED**: Write test using `Eventually` pattern
- [ ] **GREEN**: Update BaseIntent to support clock injection (optional) or use Eventually
- [ ] **REFACTOR**: Remove time.Sleep from test

**Recommended Fix**:
```go
It("should expire success after 3 seconds", func() {
    base.SetSuccess("Test message")
    Expect(base.ShouldShowSuccess()).To(BeTrue())
    Eventually(base.ShouldShowSuccess, 5*time.Second, 100*time.Millisecond).Should(BeFalse())
})
```

**Acceptance Criteria**:
- [ ] Test no longer uses time.Sleep
- [ ] Test still validates 3-second expiry behavior
- [ ] Test execution time reduced
- [ ] No flakiness in 10 consecutive runs

**Estimated Time**: 30 minutes

---

## Phase 2: CaptureEvent Submit Workflow Tests (Priority: HIGH)

**Issue**: Core user journey (event submission) has 0% test coverage.

**File**: `internal/cli/intents/capture_event_intent.go`

**Uncovered Functions**:
| Line | Function | Description |
|------|----------|-------------|
| 168 | `initializeFormForEdit()` | Edit form initialization |
| 595 | `performSubmit()` | Core submit logic |
| 696 | `performEnrichment()` | Enrichment workflow |
| 1101 | `viewError()` | Error view |
| 1131 | `acceptCurrentItem()` | Accept enrichment item |
| 1165 | `rejectCurrentItem()` | Reject enrichment item |

### Subtask 2.1: Add Submit Workflow Tests (TDD)
- [ ] **RED**: Write failing tests for `performSubmit()`
- [ ] **GREEN**: Verify existing implementation passes
- [ ] **REFACTOR**: Ensure test isolation

### Subtask 2.2: Add Enrichment Workflow Tests (TDD)
- [ ] **RED**: Write failing tests for `performEnrichment()`
- [ ] **GREEN**: Verify existing implementation passes
- [ ] **REFACTOR**: Add edge cases

### Subtask 2.3: Add Accept/Reject Tests (TDD)
- [ ] **RED**: Write failing tests for accept/reject item flows
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

**Test Cases to Add**:
```go
Describe("Submit Workflow", func() {
    It("should submit event with valid data", func() {...})
    It("should handle submission errors gracefully", func() {...})
    It("should validate required fields before submit", func() {...})
})

Describe("Enrichment Workflow", func() {
    It("should start enrichment after manual entry", func() {...})
    It("should allow accepting enrichment suggestions", func() {...})
    It("should allow rejecting enrichment suggestions", func() {...})
    It("should handle enrichment service errors", func() {...})
})
```

**Files to Create**:
- `internal/cli/intents/capture_event_submit_test.go`

**Acceptance Criteria**:
- [ ] Submit workflow fully tested
- [ ] Enrichment accept/reject tested
- [ ] Error handling tested
- [ ] Coverage increased for capture_event_intent.go

**Estimated Time**: 2-3 hours

---

## Phase 3: Delete Event Workflow Tests (Priority: HIGH)

**Issue**: Delete event is a critical operation with 0% test coverage.

**File**: `internal/cli/intents/browse_timeline_intent.go`

**Uncovered Functions**:
| Line | Function | Description |
|------|----------|-------------|
| 324 | `updateDeleteConfirm()` | Delete confirmation handler |
| 367 | `removeEventFromList()` | Event removal logic |
| 389 | `getContext()` | Context retrieval |
| 616 | `viewDeleteConfirm()` | Delete confirmation view |

### Subtask 3.1: Add Delete Confirmation Tests (TDD)
- [ ] **RED**: Write failing tests for delete confirmation
- [ ] **GREEN**: Verify existing implementation
- [ ] **REFACTOR**: Cover edge cases

### Subtask 3.2: Add Event Removal Tests (TDD)
- [ ] **RED**: Write failing tests for event removal
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Test list state after removal

**Test Cases to Add**:
```go
Describe("Delete Event Workflow", func() {
    It("should show delete confirmation when pressing 'd'", func() {...})
    It("should delete event when confirming with 'y'", func() {...})
    It("should cancel delete when pressing 'n'", func() {...})
    It("should cancel delete when pressing Escape", func() {...})
    It("should update list after successful deletion", func() {...})
    It("should handle delete errors gracefully", func() {...})
})
```

**Files to Create**:
- `internal/cli/intents/browse_timeline_delete_test.go`

**Acceptance Criteria**:
- [ ] Delete confirmation fully tested
- [ ] Delete cancellation tested
- [ ] List state verified after deletion
- [ ] Error handling tested

**Estimated Time**: 1-2 hours

---

## Phase 4: Modal Lifecycle Tests (Priority: MEDIUM)

**Issue**: Modal Update/View/Result methods have 0% coverage.

**File**: `internal/cli/intents/modals.go`

**Uncovered Functions**:
| Lines | Modal | Methods |
|-------|-------|---------|
| 156-285 | EditMetadataModal | Update, View, Result, computeChanges |
| 398-502 | EditBurstModal | View, IsComplete, syncModified, computeChanges |
| 575-710 | EditFactModal | Update, View, Result, computeChanges |

### Subtask 4.1: Add EditMetadataModal Tests (TDD)
- [ ] **RED**: Write failing tests for modal lifecycle
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

### Subtask 4.2: Add EditBurstModal Tests (TDD)
- [ ] **RED**: Write failing tests for modal lifecycle
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

### Subtask 4.3: Add EditFactModal Tests (TDD)
- [ ] **RED**: Write failing tests for modal lifecycle
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

**Test Cases for Each Modal**:
```go
Describe("EditXxxModal", func() {
    Describe("Update", func() {
        It("should handle field navigation", func() {...})
        It("should handle value input", func() {...})
        It("should handle Enter to confirm", func() {...})
        It("should handle Escape to cancel", func() {...})
    })
    
    Describe("View", func() {
        It("should render all fields", func() {...})
        It("should highlight focused field", func() {...})
        It("should show current values", func() {...})
    })
    
    Describe("Result", func() {
        It("should return modified data on confirm", func() {...})
        It("should return original data on cancel", func() {...})
        It("should track changes correctly", func() {...})
    })
})
```

**Files to Modify**:
- `internal/cli/intents/modals_test.go` (expand existing)

**Acceptance Criteria**:
- [ ] All 3 modals have Update/View/Result tests
- [ ] Change tracking tested
- [ ] Focus navigation tested
- [ ] Cancel/confirm flows tested

**Estimated Time**: 3-4 hours

---

## Phase 5: Export Format Tests (Priority: MEDIUM)

**Issue**: Export marshal functions (CSV, YAML, Text) have 0% coverage.

**File**: `internal/cli/intents/export_artifact.go`

**Uncovered Functions**:
| Line | Function | Description |
|------|----------|-------------|
| 656 | `exportFacts()` | Fact export |
| 698 | `exportBursts()` | Burst export |
| 1023 | `generateBurstsPreview()` | Burst preview |
| 1073 | `generateProfilePreview()` | Profile preview |
| 1216 | `marshalEventsYAML()` | Events to YAML |
| 1243 | `marshalEventsCSV()` | Events to CSV |
| 1284 | `marshalEventsText()` | Events to Text |
| 1308 | `marshalFactsYAML()` | Facts to YAML |
| 1329 | `marshalFactsCSV()` | Facts to CSV |
| 1366 | `marshalFactsText()` | Facts to Text |

### Subtask 5.1: Add Export Function Tests (TDD)
- [ ] **RED**: Write failing tests for each export type
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

### Subtask 5.2: Add Marshal Function Tests (TDD)
- [ ] **RED**: Write failing tests for each marshal function
- [ ] **GREEN**: Verify output format
- [ ] **REFACTOR**: Test with various data sizes

**Test Cases**:
```go
Describe("Export Functions", func() {
    DescribeTable("marshal formats",
        func(dataType, format string, expected interface{}) {...},
        Entry("events to YAML", "events", "yaml", ...),
        Entry("events to CSV", "events", "csv", ...),
        Entry("events to Text", "events", "text", ...),
        Entry("facts to YAML", "facts", "yaml", ...),
        Entry("facts to CSV", "facts", "csv", ...),
        Entry("facts to Text", "facts", "text", ...),
        Entry("bursts to YAML", "bursts", "yaml", ...),
    )
})
```

**Files to Create**:
- `internal/cli/intents/export_formats_test.go`

**Acceptance Criteria**:
- [ ] All marshal functions tested
- [ ] Output format validated
- [ ] Edge cases (empty data, special characters) tested

**Estimated Time**: 2-3 hours

---

## Phase 6: Consolidate Test Fixtures (Priority: MEDIUM)

**Issue**: Test data creation duplicated across multiple files.

**Current Duplication**:
- `internal/repository/career/sqlite_repository_test.go:49-56` - creates test events
- `internal/testutil/e2e/fixtures.go:252-261` - creates minimal events
- Various intent tests create their own test data

### Subtask 6.1: Create Shared Fixture Package
- [ ] Create `internal/testutil/fixtures/` package
- [ ] Create `event_fixtures.go` with event builders
- [ ] Create `burst_fixtures.go` with burst builders
- [ ] Create `fact_fixtures.go` with fact builders

### Subtask 6.2: Refactor Existing Tests
- [ ] Update repository tests to use shared fixtures
- [ ] Update e2e tests to use shared fixtures
- [ ] Update intent tests to use shared fixtures

**Fixture Builder Pattern**:
```go
// internal/testutil/fixtures/event_fixtures.go
package fixtures

type EventBuilder struct {
    event *career.CareerEvent
}

func NewEventBuilder() *EventBuilder {
    return &EventBuilder{
        event: &career.CareerEvent{
            ID:        uuid.New().String(),
            Text:      "Test event",
            Date:      time.Now(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }
}

func (b *EventBuilder) WithText(text string) *EventBuilder {
    b.event.Text = text
    return b
}

func (b *EventBuilder) WithCompany(company string) *EventBuilder {
    b.event.Company = company
    return b
}

func (b *EventBuilder) Build() *career.CareerEvent {
    return b.event
}

// Usage:
// event := fixtures.NewEventBuilder().WithCompany("Acme").WithText("Did stuff").Build()
```

**Files to Create**:
- `internal/testutil/fixtures/event_fixtures.go`
- `internal/testutil/fixtures/burst_fixtures.go`
- `internal/testutil/fixtures/fact_fixtures.go`
- `internal/testutil/fixtures/fixtures_test.go`

**Acceptance Criteria**:
- [ ] Builder pattern implemented for all domain types
- [ ] Existing tests refactored to use builders
- [ ] No duplicate test data creation
- [ ] Builders have tests

**Estimated Time**: 4-5 hours

---

## Phase 7: Burst Management Context Tests (Priority: LOW)

**Issue**: Most BurstManagementContext methods have 0% coverage.

**File**: `internal/cli/intents/burst_management.go`

**Uncovered Methods** (lines 165-342):
- `GetPageBursts()`
- `SelectBurst()`
- `GetSelectedBurst()`
- `CreateBurst()`
- `StartNewBurst()`
- `StartEditBurst()`
- `CancelEdit()`
- `SaveEdit()`
- `ToggleRowExpansion()`

### Subtask 7.1: Add Context Method Tests (TDD)
- [ ] **RED**: Write failing tests for context methods
- [ ] **GREEN**: Verify implementation
- [ ] **REFACTOR**: Cover edge cases

**Files to Create**:
- `internal/cli/intents/burst_management_context_test.go`

**Acceptance Criteria**:
- [ ] All context methods tested
- [ ] State transitions verified
- [ ] Edge cases covered

**Estimated Time**: 2-3 hours

---

## Phase 8: Repository Coverage Improvement (Priority: LOW)

**Issue**: Repository layer at 52.6% coverage.

**Target**: Increase to 75%+

### Subtask 8.1: Review Uncovered Repository Code
- [ ] Identify uncovered functions
- [ ] Prioritize by usage frequency

### Subtask 8.2: Add Missing Repository Tests (TDD)
- [ ] Add tests for uncovered functions
- [ ] Include error handling scenarios

**Estimated Time**: 3-4 hours

---

## Summary

### Effort Breakdown

| Phase | Priority | Time Estimate |
|-------|----------|---------------|
| 1. Fix Flaky Timer | HIGH | 30 min |
| 2. CaptureEvent Submit | HIGH | 2-3 hours |
| 3. Delete Event Workflow | HIGH | 1-2 hours |
| 4. Modal Lifecycle | MEDIUM | 3-4 hours |
| 5. Export Formats | MEDIUM | 2-3 hours |
| 6. Consolidate Fixtures | MEDIUM | 4-5 hours |
| 7. Burst Context | LOW | 2-3 hours |
| 8. Repository Coverage | LOW | 3-4 hours |
| **Total** | | **18-25 hours** |

### Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Overall Coverage | 80.34% | 85%+ |
| Intents Coverage | 61.6% | 75%+ |
| Repository Coverage | 52.6% | 75%+ |
| Flaky Tests | 1 | 0 |
| Test Fixtures Duplication | High | Eliminated |

---

## Rollback Plan

All phases are additive (adding tests). No production code changes except Phase 1 (timer fix).

- **Phase 1**: Revert `contract_test.go` changes if issues arise
- **Phases 2-8**: Delete new test files if they cause problems

---

## Related Documentation
- `docs/TESTING_PATTERNS.md` - Testing patterns guide
- `docs/rules/go-guidelines.md` - Go testing standards
- `docs/rules/senior-engineer-guidelines.md` - TDD protocol

---

**Document Version**: 1.0
**Created**: 2026-01-11
**Status**: READY FOR IMPLEMENTATION
**Process Guide**: docs/rules/master-task-prompt.md
