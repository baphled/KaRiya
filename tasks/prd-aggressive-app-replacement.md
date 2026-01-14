---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Product Requirements Document: Aggressive app.go Replacement with Intent-Based Architecture

**Document Version**: 1.0
**Date Created**: 2026-01-03
**Status**: APPROVED FOR IMPLEMENTATION
**Timeline**: 3 weeks (121 hours full-time)
**Approach**: Complete replacement with zero legacy code
**Target User**: Solo developer (KaRiya project maintainer)

---

## 1. Introduction & Overview

### Problem Statement
The current KaRiya application uses a legacy **screen-based architecture** with:
- 1,766 lines in `app.go` (bloated and hard to maintain)
- 31 screen constants
- 37+ model fields (one per screen)
- 906-line `Update()` method (complex message routing)
- 180-line `View()` method (screen rendering)
- 20+ helper methods for screen-specific logic
- Heavy breadcrumb and workflow state management
- Tight coupling between screens and application state

This architecture has become difficult to maintain, test, and extend. The codebase has already implemented a modern **intent-based system** for 5 core intents (CaptureEvent, BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem), but the root application still uses the legacy screen-based approach.

### Solution Overview
**Completely replace** the screen-based architecture in `app.go` with the proven intent-based system, creating a unified, clean, and maintainable application structure. This involves:
- Rebuilding `app.go` from scratch (~200 lines vs 1,766)
- Creating 5 new missing intents for functionality not yet migrated
- Removing all 31 legacy screen constants and 37+ model fields
- Consolidating message routing and state management
- Achieving 88% code reduction while preserving 100% functionality

### Goals
1. **Reduce complexity**: Shrink `app.go` from 1,766 to ~200 lines (88% reduction)
2. **Improve maintainability**: Eliminate legacy code pollution and screen-specific logic
3. **Preserve functionality**: Ensure all 29+ legacy screens' features work through intents
4. **Enhance testability**: Increase test coverage and ensure 164+ tests pass
5. **Maintain performance**: Meet all existing performance benchmarks
6. **Complete in 3 weeks**: Deliver aggressive timeline without sacrificing quality

---

## 2. Goals & Objectives

### Primary Goals
- [ ] **G1**: Replace all 31 legacy screens with 10 intent-based workflows
- [ ] **G2**: Reduce `app.go` to under 250 lines of clean, focused code
- [ ] **G3**: Implement 5 new missing intents (BurstManagement, FactManagement, ImportWizard, MetadataEditor, BulkOperations)
- [ ] **G4**: Achieve 164+ passing tests with 87%+ code coverage
- [ ] **G5**: Eliminate all 37+ screen-specific model fields from root model
- [ ] **G6**: Complete migration in 3 weeks (121 hours) with zero technical debt

### Secondary Goals
- [ ] **G7**: Improve developer experience with cleaner architecture
- [ ] **G8**: Enable easier addition of future intents
- [ ] **G9**: Document intent patterns for future development
- [ ] **G10**: Establish benchmarks for intent performance

---

## 3. User Stories

### User Story 1: Solo Developer Wants Cleaner Codebase
**As a** solo developer maintaining KaRiya,
**I want to** replace the bloated screen-based architecture with the modern intent-based system,
**So that** I can maintain the codebase more easily and add new features faster.

**Acceptance Criteria:**
- `app.go` is reduced to ~200 lines
- All 29+ legacy screens are replaced by intents
- Code is easier to understand and modify
- No functionality is lost

### User Story 2: Developer Wants Type-Safe Intent Communication
**As a** developer working with intents,
**I want to** use type-safe `IntentResult[T]` for all intent communication,
**So that** I can catch errors at compile-time and avoid runtime type assertions.

**Acceptance Criteria:**
- All 10 intents use `IntentResult[T]` for type safety
- No runtime type assertions needed
- Clear intent boundaries prevent state pollution
- Back navigation preserves complete context

### User Story 3: Developer Wants Comprehensive Test Coverage
**As a** developer adding new features,
**I want to** have 164+ tests covering all intents and workflows,
**So that** I can confidently refactor and extend the application.

**Acceptance Criteria:**
- 164+ Ginkgo specs passing
- 87%+ code coverage across all modules
- Zero race conditions detected
- All performance benchmarks met

### User Story 4: Developer Wants Fast Feature Development
**As a** developer adding a new intent,
**I want to** follow clear patterns and templates,
**So that** I can implement new intents in hours instead of days.

**Acceptance Criteria:**
- Clear state machine pattern established
- Reusable component library available
- Intent template documentation provided
- Example implementations available

### User Story 5: Developer Wants Reliable Navigation
**As a** a user navigating the application,
**I want to** move between intents seamlessly with back navigation,
**So that** I can explore features without losing my context.

**Acceptance Criteria:**
- Back navigation preserves complete state via metadata
- All intents support returning to previous screen
- Context is restored when returning
- No state loss during navigation

---

## 4. Functional Requirements

### 4.1 Application Architecture

**FR-1.1**: The application must use an intent-based routing system where:
- Each major workflow is implemented as a separate intent
- All intents implement the `Intent` interface with `Init()`, `Update()`, `View()`, and `Result()` methods
- The root application (`app.go`) delegates all message handling to the intent router
- Only 5-7 application-level fields exist in the root model (services, router, state, width, height)

**FR-1.2**: The application must maintain a minimal root model with:
- Services (CLIEventService, CareerService, Logger)
- IntentRouter instance
- Boolean flag for intent mode (`inIntentMode`)
- Current and previous screen state
- Window dimensions (width, height)
- Reusable components (MenuModel, HelpModel)

**FR-1.3**: The application must have only 4 root screens:
- `HomeScreen` - Main application home
- `MainMenuScreen` - Menu before home
- `HelpScreen` - Help documentation
- `QuitScreen` - Quit confirmation

### 4.2 Intent System

**FR-2.1**: The application must implement 10 intents total:
- **Already implemented (5)**:
  1. CaptureEvent - Capture new career events
  2. BrowseTimeline - View career timeline and events
  3. GenerateCV - Generate CVs from career data
  4. ExportArtifact - Export CVs and artifacts
  5. ConfigureSystem - System configuration

- **To be implemented (5)**:
  6. BurstManagement - Manage career bursts
  7. FactManagement - Manage extracted facts
  8. ImportWizard - Import career data from files
  9. MetadataEditor - Edit career event metadata
  10. BulkOperations - Perform bulk operations on events

**FR-2.2**: Each intent must implement a state machine with:
- Clear state constants (e.g., `StateListBursts`, `StateViewBurst`, `StateEditBurst`)
- Explicit state transitions in the `Update()` method
- State-specific view rendering in the `View()` method
- Typed context and result structures

**FR-2.3**: Each intent must provide:
- Input context structure (e.g., `BurstManagementContext`)
- Output result structure (e.g., `BurstManagementResult`)
- Type-safe `IntentResult[T]` for communication
- Metadata preservation for back navigation

### 4.3 Message Routing

**FR-3.1**: The application's `Update()` method must:
- Check if currently in intent mode (priority 1)
- If in intent mode, delegate to intent router's `HandleMessage()`
- If not in intent mode, handle global shortcuts and screen-specific messages
- Support menu selection to activate intents
- Support back navigation from intents

**FR-3.2**: The application must support global shortcuts:
- `Ctrl+C` or `q` - Quit application
- `?` - Show help screen
- `Esc` - Back navigation or return to home
- `Home` - Return to home screen
- Intent activation keys from menu (e.g., `c` for capture)

**FR-3.3**: The intent router must:
- Support registering intents with factory functions
- Support activating intents by name with context
- Support handling messages and returning results
- Support result handling callbacks

### 4.4 Code Reduction Requirements

**FR-4.1**: The `app.go` file must be reduced to approximately 200 lines, achieving:
- Removal of all 31 screen constants
- Removal of all 37+ screen-specific model fields
- Reduction of `Update()` from 906 to ~30 lines
- Reduction of `View()` from 180 to ~20 lines
- Elimination of 20+ screen-specific helper methods

**FR-4.2**: Legacy model files must be removed:
- All 30+ screen-specific model files in `internal/cli/models/`
- All screen-specific service methods
- All breadcrumb and workflow management code
- All screen-specific navigation logic

**FR-4.3**: Code quality must be maintained or improved:
- All code must pass `gofmt` formatting
- All code must pass `golangci-lint` checks
- All code must have no race conditions (verified with `-race` flag)
- All code must follow existing project conventions

### 4.5 Testing Requirements

**FR-5.1**: All 5 new intents must have comprehensive test coverage:
- BurstManagement: 30+ tests
- FactManagement: 40+ tests
- ImportWizard: 25+ tests
- MetadataEditor: 20+ tests
- BulkOperations: 20+ tests

**FR-5.2**: All tests must:
- Use Ginkgo v2 + Gomega framework
- Achieve 87%+ code coverage for each intent
- Test all state transitions
- Test all view rendering
- Test result handling
- Pass with zero race conditions

**FR-5.3**: Integration tests must verify:
- All intents activate from menu
- All workflows complete end-to-end
- Navigation between intents works correctly
- Results are properly handled
- Back navigation preserves context

**FR-5.4**: Performance benchmarks must:
- Intent Init: < 50ms
- View Render: < 100ms
- State Transition: < 10ms
- Router Operations: < 1ms
- Full test suite: < 5s with race detector

### 4.6 Feature Preservation

**FR-6.1**: All legacy screen functionality must be preserved:
- Event capture workflow (CaptureEvent intent)
- Timeline browsing (BrowseTimeline intent)
- Event details viewing (part of BrowseTimeline)
- CV generation (GenerateCV intent)
- CV export (ExportArtifact intent)
- System configuration (ConfigureSystem intent)
- Burst management (new BurstManagement intent)
- Fact management (new FactManagement intent)
- Data import (new ImportWizard intent)
- Metadata editing (new MetadataEditor intent)
- Bulk operations (new BulkOperations intent)

**FR-6.2**: All data operations must work identically:
- Create/Read/Update/Delete operations unchanged
- Service layer methods unchanged
- Repository queries unchanged
- Data validation unchanged
- Error handling unchanged

**FR-6.3**: All user-facing features must be accessible:
- All menu items available
- All workflows completable
- All data operations functional
- All configuration options available

---

## 5. Non-Goals (Out of Scope)

### What This PRD Does NOT Include

- ❌ **New features** - This is purely a refactoring, no new functionality added
- ❌ **Database schema changes** - All data structures remain the same
- ❌ **API changes** - Service interfaces remain compatible
- ❌ **Performance optimization** - No algorithmic improvements, only code cleanup
- ❌ **Additional intents** - Only the 5 missing intents, not future ones
- ❌ **Mobile or web versions** - TUI-only, no platform expansion
- ❌ **Backward compatibility layer** - Complete replacement, no dual-mode operation
- ❌ **Gradual migration** - All-or-nothing aggressive replacement
- ❌ **Documentation of removed code** - Legacy code is deleted, not archived

---

## 6. Design Considerations

### 6.1 Architecture Pattern

The application follows the **Type-Safe Intent Architecture** pattern:

```
┌─────────────────────────────────────────────┐
│         Root Model (app.go)                 │
│  - Services, Router, State (5-7 fields)     │
│  - Minimal Update() and View()               │
└────────────┬────────────────────────────────┘
             │
             ├─────────────────────────────────────────┐
             │                                         │
             ▼                                         ▼
    ┌─────────────────────┐          ┌──────────────────────────┐
    │  Home/Help Screens  │          │   Intent Router          │
    │  (Simple views)     │          │  - 10 registered intents │
    │                     │          │  - Message delegation    │
    └─────────────────────┘          │  - Result handling       │
                                     └──────────────────────────┘
                                               │
                                ┌──────────────┴──────────────┐
                                │                             │
                    ┌───────────────────────┐    ┌──────────────────────┐
                    │  Core Intents (5)     │    │  New Intents (5)     │
                    ├───────────────────────┤    ├──────────────────────┤
                    │ 1. CaptureEvent       │    │ 6. BurstManagement   │
                    │ 2. BrowseTimeline     │    │ 7. FactManagement    │
                    │ 3. GenerateCV         │    │ 8. ImportWizard      │
                    │ 4. ExportArtifact     │    │ 9. MetadataEditor    │
                    │ 5. ConfigureSystem    │    │10. BulkOperations    │
                    └───────────────────────┘    └──────────────────────┘
```

### 6.2 Intent State Machine Pattern

Each intent follows this pattern:

```go
type IntentState string

const (
    StateInitial IntentState = "initial"
    StateWorking IntentState = "working"
    StateFinal   IntentState = "final"
)

type IntentModel struct {
    state  IntentState
    data   *IntentContext
    result *IntentResult[*IntentResult]
}

func (i *IntentModel) Update(msg tea.Msg) tea.Cmd {
    switch i.state {
    case StateInitial:
        return i.handleInitial(msg)
    case StateWorking:
        return i.handleWorking(msg)
    case StateFinal:
        return i.handleFinal(msg)
    }
    return nil
}

func (i *IntentModel) View() string {
    switch i.state {
    case StateInitial:
        return i.viewInitial()
    case StateWorking:
        return i.viewWorking()
    case StateFinal:
        return i.viewFinal()
    }
    return ""
}
```

### 6.3 Component Reuse

Reusable UI components in `internal/cli/components/`:
- Card - Display data in a card format
- List - Render scrollable lists
- Form - Handle form input
- Modal - Display modal dialogs
- Progress - Show progress indicators
- Menu - Display menu options

### 6.4 Styling

All styling uses Lipgloss and is centralized in `internal/cli/styles/styles.go`:
- Color scheme consistent across all intents
- Border styles uniform
- Spacing and padding standardized
- Responsive to terminal width changes

---

## 7. Technical Considerations

### 7.1 Dependencies

**Core Dependencies** (already in use):
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/onsi/ginkgo/v2` - Testing framework
- `github.com/onsi/gomega` - Assertion library
- `modernc.org/sqlite` - Database

**No new dependencies required** - All infrastructure already exists.

### 7.2 Compatibility

**Breaking Changes**: None for external users (solo developer)
- This is an internal refactoring
- All service interfaces remain the same
- All data structures remain the same
- All functionality preserved

**Internal Changes**: Significant
- Screen-based models replaced with intents
- Message routing consolidated
- State management simplified

### 7.3 Performance Implications

**Expected Performance**: No degradation
- Intent system already benchmarked
- Simpler code should be faster
- Message routing more efficient
- State management cleaner

**Verification**: All benchmarks must pass
- Intent Init: < 50ms
- View Render: < 100ms
- State Transition: < 10ms
- Router Operations: < 1ms

### 7.4 Data Integrity

**No changes to**:
- Database schema
- Data access layer (repository)
- Business logic (service layer)
- Domain models

**All CRUD operations** remain identical.

### 7.5 Testing Strategy

**Test Framework**: Ginkgo v2 + Gomega (already in use)
**Test Execution**:
```bash
# All tests
go test -v ./...

# With race detector
go test -race ./...

# Specific intent
go test -v ./internal/cli/intents/burst_management_test.go

# With coverage
go test -coverprofile=coverage.out ./...
```

---

## 8. Success Metrics

### 8.1 Code Metrics

| Metric | Target | Verification |
|--------|--------|--------------|
| `app.go` lines | < 250 | `wc -l internal/cli/app/app.go` |
| `Update()` lines | < 50 | Manual review |
| `View()` lines | < 30 | Manual review |
| Model fields | < 10 | Manual review |
| Screen constants | 4 only | Grep for "const (" |
| Legacy files removed | All 30+ | Check `internal/cli/models/` |

### 8.2 Quality Metrics

| Metric | Target | Verification |
|--------|--------|--------------|
| Tests passing | 164+ | `go test -v ./...` |
| Code coverage | 87%+ | `go test -cover ./...` |
| Race conditions | 0 | `go test -race ./...` |
| Lint issues | 0 | `golangci-lint run ./...` |
| Format issues | 0 | `gofmt -l ./...` |

### 8.3 Functional Metrics

| Metric | Target | Verification |
|--------|--------|--------------|
| Intents implemented | 10 | Count intent files |
| Workflows functional | 100% | Manual testing |
| Features preserved | 100% | Feature checklist |
| Navigation working | 100% | Manual testing |
| Data operations | 100% | Integration tests |

### 8.4 Performance Metrics

| Metric | Target | Verification |
|--------|--------|--------------|
| Intent Init | < 50ms | `go test -bench=BenchmarkInit` |
| View Render | < 100ms | `go test -bench=BenchmarkView` |
| State Transition | < 10ms | `go test -bench=BenchmarkUpdate` |
| Router Operations | < 1ms | `go test -bench=BenchmarkRouter` |
| Full test suite | < 5s | `go test -race ./...` (with timing) |

### 8.5 Timeline Metrics

| Phase | Duration | Verification |
|-------|----------|--------------|
| Phase 1: Preparation | 2-3 days | Audit checklist complete |
| Phase 2: New Intents | 4-5 days | 5 intents, 135+ tests |
| Phase 3: Rebuild app.go | 3-4 days | ~200 lines, all compiled |
| Phase 4: Testing | 2-3 days | 164+ tests passing |
| **Total** | **3 weeks** | All milestones met |

---

## 9. Implementation Phases

### Phase 1: Preparation (9 hours)

**Goal**: Audit legacy code and plan implementation

**Tasks**:
1. Document all 29+ legacy screen models
2. Identify which intent handles each screen
3. Note special logic to preserve
4. Check for cross-model dependencies
5. Identify missing intents (5 total)
6. Plan new intent implementations

**Deliverables**:
- Audit document with screen-to-intent mapping
- Feature gap analysis
- Implementation plan for 5 new intents

**Success Criteria**:
- All screens mapped to intents
- No missing features identified
- Clear implementation plan established

### Phase 2: Create Missing Intents (64 hours)

**Goal**: Implement 5 new intents with full test coverage

**Intents to Create**:
1. **BurstManagement** (16 hours)
   - States: List, View, Edit, Suggest, Confirm
   - Tests: 30+
   - Coverage: >90%

2. **FactManagement** (16 hours)
   - States: List, View, Edit, Review, Confirm
   - Tests: 40+
   - Coverage: >90%

3. **ImportWizard** (12 hours)
   - States: SelectFile, Review, Progress, Complete, Confirm
   - Tests: 25+
   - Coverage: >90%

4. **MetadataEditor** (10 hours)
   - States: Review, Edit, Confirm
   - Tests: 20+
   - Coverage: >90%

5. **BulkOperations** (10 hours)
   - States: SelectOp, Configure, Execute, Confirm
   - Tests: 20+
   - Coverage: >90%

**Deliverables**:
- 5 intent implementations
- 135+ tests (all passing)
- 100% functionality coverage

**Success Criteria**:
- All intents compile
- All tests pass
- Coverage > 90% per intent
- No race conditions

### Phase 3: Rebuild app.go (14 hours)

**Goal**: Replace legacy app.go with clean intent-based version

**Tasks**:
1. Create new minimal `app.go` (~200 lines)
2. Register all 10 intents with router
3. Implement menu selection
4. Implement intent activation
5. Remove all legacy code
6. Verify compilation

**Deliverables**:
- New `app.go` (~200 lines)
- All legacy model files removed
- All legacy service methods removed
- Clean project structure

**Success Criteria**:
- Code compiles without errors
- `app.go` < 250 lines
- `Update()` < 50 lines
- `View()` < 30 lines
- All 30+ legacy files removed

### Phase 4: Testing & Validation (34 hours)

**Goal**: Comprehensive testing and quality assurance

**Tasks**:
1. Unit tests for new intents (20 hours)
2. Integration tests (20 hours)
3. End-to-end workflow testing (8 hours)
4. Performance benchmarking (6 hours)
5. Code quality review (4 hours)
6. Documentation updates (4 hours)

**Deliverables**:
- 164+ passing tests
- Performance benchmarks verified
- Code quality report
- Updated documentation

**Success Criteria**:
- 164+ tests passing
- 87%+ code coverage
- 0 race conditions
- All benchmarks met
- All lint checks passing

---

## 10. Risk Assessment & Mitigation

### Risk 1: Breaking Changes
**Severity**: HIGH
**Probability**: LOW (solo developer, fully tested)

**Mitigation**:
- ✅ You are the sole user
- ✅ All intents are fully tested
- ✅ Can test locally before deploying
- ✅ Git history allows rollback
- ✅ Feature branch for isolated work

### Risk 2: Missing Features
**Severity**: HIGH
**Probability**: LOW (comprehensive audit planned)

**Mitigation**:
- ✅ Phase 1 includes feature audit
- ✅ Each intent maps to existing screens
- ✅ No new features, just reorganization
- ✅ Comprehensive feature checklist
- ✅ Integration tests verify all features

### Risk 3: Performance Degradation
**Severity**: MEDIUM
**Probability**: LOW (intent system already optimized)

**Mitigation**:
- ✅ Intent system already benchmarked
- ✅ Phase 4 includes performance testing
- ✅ Simpler code should be faster
- ✅ Profile and optimize as needed
- ✅ Rollback if issues found

### Risk 4: Timeline Slippage
**Severity**: MEDIUM
**Probability**: MEDIUM (aggressive timeline)

**Mitigation**:
- ✅ Conservative estimates per phase
- ✅ Clear checkpoints at each phase
- ✅ Buffer time in schedule
- ✅ Parallel implementation possible
- ✅ Prioritize core intents if needed

### Risk 5: Code Quality Issues
**Severity**: MEDIUM
**Probability**: LOW (comprehensive testing)

**Mitigation**:
- ✅ All code must pass linting
- ✅ 87%+ code coverage required
- ✅ Race detector verification
- ✅ Code review before merge
- ✅ Automated CI/CD checks

### Risk 6: Integration Issues
**Severity**: MEDIUM
**Probability**: LOW (existing patterns)

**Mitigation**:
- ✅ Follow established intent patterns
- ✅ Reuse existing components
- ✅ Integration tests verify workflows
- ✅ Manual testing of all intents
- ✅ Gradual activation from menu

---

## 11. Open Questions

1. **Prioritization**: If timeline becomes constrained, should we prioritize completing all 5 new intents, or can some be deferred?

2. **Documentation**: Should we create a detailed "Intent Implementation Guide" for future developers, or is the current documentation sufficient?

3. **Backward Compatibility**: Should we maintain any backward compatibility with the old screen-based architecture, or proceed with complete replacement?

4. **Deployment Strategy**: Should we deploy immediately after completion, or perform a staged rollout with testing period?

5. **Future Intents**: What other intents might be needed in the future? Should we design the framework with specific extensibility in mind?

6. **Performance Profiling**: Should we include detailed performance profiling in Phase 4, or focus on basic benchmarks?

7. **Error Handling**: Should we implement enhanced error handling for edge cases in the new intents?

8. **Logging**: Should we add structured logging to track intent activation and state transitions?

---

## 12. Acceptance Criteria (Master Checklist)

### Code Reduction
- [ ] `app.go` reduced to < 250 lines (from 1,766)
- [ ] `Update()` method < 50 lines (from 906)
- [ ] `View()` method < 30 lines (from 180)
- [ ] Model fields < 10 (from 37+)
- [ ] All 31 screen constants removed
- [ ] All 30+ legacy model files deleted

### Intent Implementation
- [ ] 5 new intents fully implemented
- [ ] 10 total intents registered with router
- [ ] All intents follow state machine pattern
- [ ] All intents use type-safe `IntentResult[T]`
- [ ] All intents support back navigation

### Testing & Quality
- [ ] 164+ tests passing (100% pass rate)
- [ ] 87%+ code coverage
- [ ] 0 race conditions detected
- [ ] All lint checks passing
- [ ] All performance benchmarks met
- [ ] Code formatted with gofmt

### Feature Preservation
- [ ] All 29+ legacy screen features preserved
- [ ] All CRUD operations working
- [ ] All menu items functional
- [ ] All workflows completable
- [ ] All data operations unchanged
- [ ] All error handling preserved

### Integration
- [ ] All intents activate from menu
- [ ] Navigation between intents works
- [ ] Results properly handled
- [ ] Context preserved on back navigation
- [ ] Services properly initialized
- [ ] Global state properly managed

### Documentation
- [ ] Implementation guide updated
- [ ] Intent patterns documented
- [ ] Code comments clear and complete
- [ ] README reflects new architecture
- [ ] Phase completion report written
- [ ] Known issues documented

### Deployment
- [ ] Feature branch tested locally
- [ ] All tests passing on CI/CD
- [ ] Code review completed
- [ ] Merge to main approved
- [ ] Deployment checklist verified
- [ ] Post-deployment validation done

---

## 13. Definition of Done

This PRD is considered complete when:

1. ✅ All 10 intents are fully implemented and tested
2. ✅ `app.go` is rebuilt with < 250 lines of clean code
3. ✅ All 30+ legacy model files are removed
4. ✅ 164+ tests are passing with 87%+ coverage
5. ✅ Zero race conditions detected
6. ✅ All performance benchmarks met
7. ✅ All lint and format checks passing
8. ✅ All features from legacy screens preserved and working
9. ✅ Complete navigation between all intents functional
10. ✅ Code merged to main branch
11. ✅ Phase completion report written
12. ✅ Documentation updated

---

## 14. Timeline Summary

| Week | Phase | Duration | Key Deliverables |
|------|-------|----------|------------------|
| **Week 1** | Preparation + Phase 2 Start | 5 days | Audit complete, 2-3 intents started |
| **Week 2** | Phase 2 Completion + Phase 3 | 5 days | All 5 intents complete, app.go rebuilt |
| **Week 3** | Phase 4 + Polish | 5 days | All tests passing, ready for deployment |

**Total**: 3 weeks (121 hours full-time)

---

## 15. Approval & Sign-Off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Developer | [Your Name] | 2026-01-03 | ________________ |
| Reviewer | (Self-review) | 2026-01-03 | ________________ |

---

## Appendix A: Screen-to-Intent Mapping

| Legacy Screen | Intent | Status | Tests |
|---------------|--------|--------|-------|
| CaptureScreen | CaptureEvent | ✅ Exists | 30+ |
| ListScreen | BrowseTimeline | ✅ Exists | 37 |
| ViewScreen | BrowseTimeline | ✅ Exists | 37 |
| CVConfigManagerScreen | ConfigureSystem | ✅ Exists | 400+ |
| CVGeneratorScreen | GenerateCV | ✅ Exists | 41 |
| CVPreviewScreen | GenerateCV | ✅ Exists | 41 |
| CVExportDialogScreen | ExportArtifact | ✅ Exists | 411 |
| CVListScreen | BrowseTimeline | ✅ Exists | 37 |
| BurstListScreen | BurstManagement | ✨ NEW | 30+ |
| BurstDetailsScreen | BurstManagement | ✨ NEW | 30+ |
| BurstEditorScreen | BurstManagement | ✨ NEW | 30+ |
| BurstSuggestionScreen | BurstManagement | ✨ NEW | 30+ |
| FactListScreen | FactManagement | ✨ NEW | 40+ |
| FactDetailsScreen | FactManagement | ✨ NEW | 40+ |
| FactEditorScreen | FactManagement | ✨ NEW | 40+ |
| FactActionMenuScreen | FactManagement | ✨ NEW | 40+ |
| FactsResultsScreen | FactManagement | ✨ NEW | 40+ |
| ImportReviewScreen | ImportWizard | ✨ NEW | 25+ |
| ImportProgressScreen | ImportWizard | ✨ NEW | 25+ |
| MetadataReviewScreen | MetadataEditor | ✨ NEW | 20+ |
| MetadataEditorScreen | MetadataEditor | ✨ NEW | 20+ |
| BulkOperationsScreen | BulkOperations | ✨ NEW | 20+ |
| SuccessScreen | Generic (not intent) | ✅ Keep | - |
| ActionMenuScreen | Generic (not intent) | ✅ Keep | - |
| MainMenuScreen | Root app | ✅ Keep | - |
| HomeScreen | Root app | ✅ Keep | - |
| HelpScreen | Root app | ✅ Keep | - |
| ConfirmationScreen | Modal dialog | ✅ Keep | - |
| QuitScreen | Root app | ✅ Keep | - |

---

## Appendix B: Implementation Patterns

### Pattern 1: Intent State Machine
```go
type YourState string

const (
    StateInitial YourState = "initial"
    StateWorking YourState = "working"
    StateFinal   YourState = "final"
)

func (y *YourIntentModel) Update(msg tea.Msg) tea.Cmd {
    switch y.state {
    case StateInitial:
        return y.handleInitial(msg)
    case StateWorking:
        return y.handleWorking(msg)
    case StateFinal:
        return y.handleFinal(msg)
    }
    return nil
}
```

### Pattern 2: Back Navigation with Context
```go
result := &IntentResult[*YourIntentResult]{
    Status: StatusCompleted,
    Data:   y.result.Data,
}
result.WithMetadata("scroll_position", 42)
result.WithMetadata("selection", selectedID)
return result
```

### Pattern 3: Modal Sub-Flow
```go
type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}

if editModal.WasAccepted() {
    y.data.Field = editModal.Modified
} else {
    y.data.Field = editModal.Original
}
```

---

**Document Version**: 1.0
**Status**: READY FOR IMPLEMENTATION
**Last Updated**: 2026-01-03

