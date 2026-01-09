# Task 37: E2E Integration Test Enhancement

## Overview
- **Goal**: Create comprehensive end-to-end integration tests to ensure complete workflow coverage, prevent regressions, and verify core functionality across all 10 intents
- **Time Estimate**: 6-8 days
- **Prerequisites**: All 10 intents functional, SQLite persistence working, existing test infrastructure

## Session Contract Acknowledgment
- [x] Ran `make session-start` and reviewed results
- [x] Acknowledge and commit to following all workflow rules
- [x] Note: Coverage check ignored (64% vs 80%) - this task will improve it
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` reviewed (coverage exception noted)
- [x] Reviewed existing integration test patterns in:
  - `internal/cli/app/app_integration_test.go`
  - `internal/cli/intents/generate_cv_integration_test.go`
  - `internal/cli/intents/*_escape_test.go`
  - `internal/service/career/integration_test.go`
  - `cmd/cli/persistence_test.go`
- [x] Confirmed test utilities available in `internal/testutil/db.go`
- [x] Completed implementation audit of all 10 intents

---

## Implementation Audit Summary

### Intent Implementation Status Matrix

| Intent | Category | Implementation | Service Integration | E2E Testable |
|--------|----------|----------------|---------------------|--------------|
| **CaptureEvent** | Core | Fully functional | Full CRUD | Yes |
| **BrowseTimeline** | Core | Partial (filters missing) | Read-only | Yes |
| **GenerateCV** | Core | Partial (edit stub) | Full generation | Yes |
| **ExportArtifact** | Core | Partial (PDF/Email stubs) | File export | Yes |
| **ConfigureSystem** | Core | Partial (select inputs) | Config persistence | Yes |
| **BurstManagement** | Secondary | Fully functional | Full CRUD | Yes |
| **FactManagement** | Secondary | Fully functional | Full CRUD | Yes |
| **ImportWizard** | Secondary | UI shell only | None | Limited |
| **MetadataEditor** | Secondary | UI shell only | None | Limited |
| **BulkOperations** | Secondary | UI shell only | None | Limited |

### Key Gaps Identified (GitHub Issues to Create)

#### BrowseTimeline (5 issues)
- [ ] 'f' Filter shortcut not implemented
- [ ] Date range filter not implemented
- [ ] Company filter not implemented
- [ ] Categories filter not implemented
- [ ] Fact selection not connected

#### GenerateCV (2 issues)
- [ ] CV editing in review state is a stub
- [ ] Profile creation not available

#### ExportArtifact (3 issues)
- [ ] PDF export defined but not implemented
- [ ] Email export defined but not implemented
- [ ] Profile export returns hardcoded data

#### ConfigureSystem (3 issues)
- [ ] Select type shows text input instead of dropdown
- [ ] Type assertion panic risk
- [ ] Config not applied at runtime

#### Secondary Intents (3 umbrella issues)
- [ ] ImportWizard: No file browser, CSV parsing, or DB import
- [ ] MetadataEditor: No text inputs, entity loading, or persistence
- [ ] BulkOperations: No item selection or actual operations

---

## Files to Create

### Phase 1: Test Infrastructure & Audit
- [ ] `internal/cli/intents/e2e_test_helpers.go` - Reusable E2E test utilities
- [ ] `internal/cli/intents/e2e_test_fixtures.go` - Sample data fixtures
- [ ] `docs/E2E_IMPLEMENTATION_GAPS.md` - Gap documentation

### Phase 2: Core 5 Intents
- [ ] `internal/cli/intents/e2e_capture_workflow_test.go`
- [ ] `internal/cli/intents/e2e_browse_workflow_test.go`
- [ ] `internal/cli/intents/e2e_generate_cv_workflow_test.go`
- [ ] `internal/cli/intents/e2e_export_workflow_test.go`
- [ ] `internal/cli/intents/e2e_configure_workflow_test.go`

### Phase 3: Secondary 5 Intents
- [ ] `internal/cli/intents/e2e_burst_management_test.go`
- [ ] `internal/cli/intents/e2e_fact_management_test.go`
- [ ] `internal/cli/intents/e2e_import_wizard_test.go` (navigation only)
- [ ] `internal/cli/intents/e2e_metadata_editor_test.go` (navigation only)
- [ ] `internal/cli/intents/e2e_bulk_operations_test.go` (navigation only)

### Phase 4: Cross-Intent & Error Recovery
- [ ] `internal/cli/intents/e2e_chained_workflows_test.go`
- [ ] `internal/cli/intents/e2e_error_recovery_test.go`

---

## Phase 1: Test Infrastructure

### 1.1 Create E2E Test Helpers (`e2e_test_helpers.go`)

**Functions to implement:**

```go
// E2ETestEnv holds all test dependencies
type E2ETestEnv struct {
    Model      *app.Model
    DB         *sql.DB
    Repo       career.EventRepository
    BurstRepo  career.BurstRepository
    FactRepo   career.FactRepository
    Service    *careerservice.Service
    CLIService *service.CLIEventService
    Cleanup    func()
}

// SetupE2ETest creates full test environment with SQLite
func SetupE2ETest(t *testing.T) *E2ETestEnv

// Navigation helpers
func (e *E2ETestEnv) SelectIntent(index int) *E2ETestEnv
func (e *E2ETestEnv) PressKey(key tea.KeyType) *E2ETestEnv
func (e *E2ETestEnv) PressKeyRune(r rune) *E2ETestEnv
func (e *E2ETestEnv) PressKeys(keys ...interface{}) *E2ETestEnv
func (e *E2ETestEnv) TypeText(text string) *E2ETestEnv

// Assertion helpers
func (e *E2ETestEnv) AssertViewContains(substr string) *E2ETestEnv
func (e *E2ETestEnv) AssertViewNotContains(substr string) *E2ETestEnv
func (e *E2ETestEnv) GetView() string

// Data verification
func (e *E2ETestEnv) AssertEventCount(count int) *E2ETestEnv
func (e *E2ETestEnv) AssertBurstCount(count int) *E2ETestEnv
func (e *E2ETestEnv) AssertFactCount(count int) *E2ETestEnv

// Session simulation
func (e *E2ETestEnv) SimulateRestart() *E2ETestEnv
```

**TDD Checklist - Phase 1.1**:
- [ ] RED: Write test for `SetupE2ETest()` returning valid environment
- [ ] GREEN: Implement `SetupE2ETest()`
- [ ] RED: Write test for `SelectIntent()` navigating correctly
- [ ] GREEN: Implement `SelectIntent()`
- [ ] RED: Write test for `PressKey()` updating model
- [ ] GREEN: Implement `PressKey()`
- [ ] RED: Write test for assertion helpers
- [ ] GREEN: Implement assertion helpers
- [ ] REFACTOR: Extract common patterns

### 1.2 Create Test Fixtures (`e2e_test_fixtures.go`)

**Functions to implement:**

```go
// CreateSampleEvents generates test career events
func CreateSampleEvents(count int) []*career.CareerEvent

// CreateSampleBursts generates test bursts linked to events
func CreateSampleBursts(count int, events []*career.CareerEvent) []*career.Burst

// CreateSampleFacts generates test facts linked to events
func CreateSampleFacts(count int, events []*career.CareerEvent) []*career.Fact

// CreateSampleProfiles generates CV profiles
func CreateSampleProfiles() []*intents.CVProfile

// PopulateTestDatabase adds fixtures to repositories
func PopulateTestDatabase(env *E2ETestEnv, events []*career.CareerEvent, bursts []*career.Burst, facts []*career.Fact) error
```

**TDD Checklist - Phase 1.2**:
- [ ] RED: Write test for `CreateSampleEvents()` returning valid events
- [ ] GREEN: Implement `CreateSampleEvents()`
- [ ] RED: Write test for `CreateSampleBursts()` with event references
- [ ] GREEN: Implement `CreateSampleBursts()`
- [ ] RED: Write test for `CreateSampleFacts()` with event references
- [ ] GREEN: Implement `CreateSampleFacts()`
- [ ] RED: Write test for `PopulateTestDatabase()`
- [ ] GREEN: Implement `PopulateTestDatabase()`
- [ ] REFACTOR: Ensure fixtures follow domain validation

### 1.3 Create Gap Documentation

- [ ] Create `docs/E2E_IMPLEMENTATION_GAPS.md`
- [ ] Create GitHub issues for each gap (~15 issues)

---

## Phase 2: Core 5 Intents E2E Tests

### 2.1 CaptureEvent E2E (`e2e_capture_workflow_test.go`)

| Scenario | States Covered | Verification |
|----------|----------------|--------------|
| Complete Quick Capture | ChooseStrategy -> Form -> Review -> Submit | Event persisted with auto-date |
| Complete Manual Capture | All states with all fields | All fields stored correctly |
| Edit Burst Modal | Review -> EditBurst -> Review | Burst changes tracked |
| Edit Fact Modal | Review -> EditFact -> Review | Fact changes tracked |
| Cancel at Each State | All states -> Cancel | No data persisted |
| Back Navigation | All states -> Back | Previous state restored |

**Target: 25+ test specs**

### 2.2 BrowseTimeline E2E (`e2e_browse_workflow_test.go`)

| Scenario | Verification |
|----------|--------------|
| View Timeline with Events | Events displayed chronologically |
| Navigate to Event Detail | Detail view shows event info |
| Back from Detail | Scroll position preserved |
| Empty Timeline | Graceful empty message |
| Context Preservation | State preserved on back navigation |

**Note**: Filter tests skipped (not implemented)
**Target: 15+ test specs**

### 2.3 GenerateCV E2E (`e2e_generate_cv_workflow_test.go`)

| Scenario | Verification |
|----------|--------------|
| Complete CV Generation | CV generated with bullets |
| Export to File (Text) | File written with content |
| Export to File (Markdown) | Proper markdown formatting |
| Export to File (YAML) | Valid YAML structure |
| Export to Clipboard | Content copied (mock) |
| Error During Generation | Error displayed |
| CV with No Facts | Graceful handling |
| Cancel at Each State | Returns to appropriate state |

**Target: 25+ test specs**

### 2.4 ExportArtifact E2E (`e2e_export_workflow_test.go`)

| Scenario | Verification |
|----------|--------------|
| Export Events to JSON | Valid JSON created |
| Export Events to CSV | Valid CSV with headers |
| Export Events to YAML | Valid YAML structure |
| Export Events to TXT | Plain text format |
| Export Facts to JSON | All facts exported |
| Export Bursts to JSON | Bursts with event refs |
| Export to Clipboard | Content copied |
| Preview Before Export | Content shown |
| Cancel Export | No partial exports |

**Note**: PDF/Email tests skipped (not implemented)
**Target: 20+ test specs**

### 2.5 ConfigureSystem E2E (`e2e_configure_workflow_test.go`)

| Scenario | Verification |
|----------|--------------|
| Configure Profile Settings | Settings saved to file |
| Configure System Settings | System settings saved |
| Configure Export Settings | Export settings saved |
| Configure UI Settings | UI settings saved |
| Cancel Without Saving | Original values preserved |
| Persist After Restart | Settings loaded from file |
| Review Changes Diff | Diff shown correctly |

**Target: 20+ test specs**

---

## Phase 3: Secondary 5 Intents E2E Tests

### 3.1 BurstManagement E2E (`e2e_burst_management_test.go`)

| Scenario | Verification |
|----------|--------------|
| View Burst List | Bursts displayed with event counts |
| View Burst Detail | Associated events shown |
| Confirm Burst | Marked confirmed in DB |
| Delete Burst | Removed from DB |
| View Facts for Burst | Facts displayed |
| Extract Facts | Facts created in DB |

**Target: 20+ test specs**

### 3.2 FactManagement E2E (`e2e_fact_management_test.go`)

| Scenario | Verification |
|----------|--------------|
| View Fact List | Facts displayed with sources |
| View Fact Detail | Full details shown |
| Create New Fact | Fact persisted |
| Delete Fact | Removed from DB |
| Pagination | Pages work correctly |

**Target: 15+ test specs**

### 3.3-3.5 ImportWizard/MetadataEditor/BulkOperations (Navigation Only)

| Test Type | Verification |
|-----------|--------------|
| State Navigation | All states reachable |
| Escape Navigation | Back works at each state |
| View Rendering | No panics, content displayed |
| Cancel Workflow | Returns to menu |

**Target: 10+ test specs each (30 total)**

---

## Phase 4: Cross-Intent & Error Recovery

### 4.1 Chained Workflows (`e2e_chained_workflows_test.go`)

| Scenario | Workflow |
|----------|----------|
| Full Career to CV | Capture 5 events -> Browse -> GenerateCV -> Export |
| Capture to Burst | Capture events -> BurstManagement -> View bursts |
| Multi-Session | Session1: Capture -> Session2: Browse |
| Browse then CV | Browse timeline -> Select events -> Generate CV |
| Config then Export | Configure defaults -> Export with defaults |

**Target: 15+ test specs**

### 4.2 Error Recovery (`e2e_error_recovery_test.go`)

| Scenario | Verification |
|----------|--------------|
| Empty Database | All intents handle gracefully |
| Invalid Input | Validation errors shown |
| Service Error | Error displayed, can retry |
| Partial Data | Graceful degradation |

**Target: 15+ test specs**

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make review-commit` passes
- [ ] AI attribution included (if AI-generated)
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] All tests pass with race detector: `go test -race ./internal/cli/intents/...`
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] GitHub issues created for all gaps

## Acceptance Criteria

### Coverage Requirements
- [ ] All 10 intents have E2E tests
- [ ] At least 5 chained workflow scenarios tested
- [ ] Error recovery tests cover 4+ error types
- [ ] Empty state handling for all intents

### Quality Requirements
- [ ] All tests pass with race detector enabled
- [ ] 100% pass rate on new E2E tests
- [ ] No flaky tests (run 3x to verify)
- [ ] Test suite completes in <60 seconds

### Documentation Requirements
- [ ] Gap documentation complete
- [ ] ~15 GitHub issues created
- [ ] Task file updated with completion status

### Verification Commands
```bash
# Run all E2E tests
go test -v -race ./internal/cli/intents/... -run "E2E"

# Run with coverage
go test -coverprofile=coverage.out ./internal/cli/intents/... -run "E2E"
go tool cover -func=coverage.out | grep -E "e2e_"

# Run specific phase tests
go test -v ./internal/cli/intents/... -run "CaptureWorkflow"
go test -v ./internal/cli/intents/... -run "ChainedWorkflows"
go test -v ./internal/cli/intents/... -run "ErrorRecovery"
```

## Success Metrics

| Metric | Target |
|--------|--------|
| E2E Test Files | 14 |
| Total E2E Test Specs | 200+ |
| Intents with E2E Coverage | 10/10 |
| Chained Workflow Tests | 5+ |
| Error Recovery Tests | 7+ |
| GitHub Issues Created | ~15 |
| Test Execution Time | <60s |
| Pass Rate | 100% |

## Rollback Plan

If E2E tests introduce issues:

1. **Revert commits**: Each phase is independently committable
2. **Keep infrastructure**: Test helpers are reusable even if some tests removed
3. **Disable flaky tests**: Use `Skip()` with explanation rather than delete
4. **Isolate issues**: Each test file is independent, can be disabled individually
