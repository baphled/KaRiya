# KaRiya TUI Intent Architecture Implementation Roadmap

## Overview

This document provides a detailed implementation roadmap for the production-ready TUI intent architecture defined in [TUI_INTENT_DIAGRAM.md](TUI_INTENT_DIAGRAM.md). The architecture is **ready for immediate implementation** with no blockers.

## Architecture Status: ✅ Production-Ready

### Key Achievements
- ✅ Type-safe intent communication via `IntentResult[T]`
- ✅ Clear intent ownership rules preventing state pollution
- ✅ Predictable state machines with explicit transitions
- ✅ Back navigation preserving complete context
- ✅ Consistent async operation patterns
- ✅ Modal edit pattern with typed diffs
- ✅ Comprehensive testing strategy
- ✅ Future-proof secondary intent design

---

## Phase 1: Foundation & Core Infrastructure (Weeks 1-2)

### Goal
Establish the foundational types, interfaces, and root-level orchestration that all intents depend on.

### Tasks

#### 1.1 Implement Intent Boundary Contract Types
**Location**: `internal/cli/intents/contract.go`

```go
// Core types to implement:
type IntentStatus string
const (
    StatusCompleted IntentStatus = "completed"
    StatusCancelled  IntentStatus = "cancelled"
    StatusFailed     IntentStatus = "failed"
    StatusPartial    IntentStatus = "partial"
)

type IntentError struct {
    Code    string        // Machine-readable error code
    Message string        // User-facing message
    Cause   error         // Underlying error for logging
}

type IntentResult[T any] struct {
    Status   IntentStatus
    Data     T
    Error    *IntentError
    Metadata map[string]interface{}
}

type Intent interface {
    // Initialize the intent with context
    Init(ctx context.Context) tea.Cmd

    // Handle updates
    Update(msg tea.Msg) (Intent, tea.Cmd)

    // Render the intent
    View() string

    // Return the result when complete
    Result() *IntentResult[interface{}]
}

type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}
```

**Checklist**:
- [ ] Define `IntentStatus` constants
- [ ] Implement `IntentError` with proper error wrapping
- [ ] Create generic `IntentResult[T]` with metadata support
- [ ] Define `Intent` interface contract
- [ ] Implement `ModalEditResult[T]` for sub-flows
- [ ] Add helper methods: `WithMetadata()`, `GetMetadata()`, `IsSuccess()`
- [ ] Write unit tests for all types

**Estimated Time**: 2-3 days

---

#### 1.2 Create IntentRouter
**Location**: `internal/cli/intents/router.go`

```go
type IntentRouter struct {
    intents       map[string]Intent
    current       Intent
    history       []Intent
    resultHandler func(IntentResult[interface{}])
}

func (r *IntentRouter) Activate(intentName string, ctx context.Context) tea.Cmd
func (r *IntentRouter) HandleResult(result IntentResult[interface{}]) tea.Cmd
func (r *IntentRouter) Back() tea.Cmd
func (r *IntentRouter) Update(msg tea.Msg) (IntentRouter, tea.Cmd)
func (r *IntentRouter) View() string
```

**Checklist**:
- [ ] Implement intent activation with context
- [ ] Implement result propagation
- [ ] Implement back navigation with history
- [ ] Add intent registration mechanism
- [ ] Write integration tests for router
- [ ] Test intent switching and back navigation

**Estimated Time**: 2-3 days

---

#### 1.3 Refactor Root Model to Use IntentRouter
**Location**: `internal/cli/app/app.go`

Update the root Bubble Tea model to:
- Use `IntentRouter` instead of inline state management
- Delegate all updates to the active intent
- Render the active intent's view
- Handle intent results properly

**Checklist**:
- [ ] Replace inline state with `IntentRouter`
- [ ] Implement `Update()` to delegate to router
- [ ] Implement `View()` to render active intent
- [ ] Add intent result handling
- [ ] Update all existing message types
- [ ] Write integration tests

**Estimated Time**: 2-3 days

---

### Phase 1 Validation

**Acceptance Criteria**:
- [ ] All intent boundary contract types compile and have full test coverage
- [ ] `IntentRouter` can activate, switch, and navigate back through intents
- [ ] Root model successfully delegates to `IntentRouter`
- [ ] No breaking changes to existing CLI functionality
- [ ] All tests pass (unit + integration)

**Estimated Phase 1 Duration**: 1.5 weeks

---

## Phase 2: CaptureEvent Intent Implementation (Weeks 3-4)

### Goal
Implement the complete `CaptureEvent` intent with all sub-flows, establishing the pattern for other intents.

### Tasks

#### 2.1 Define CaptureEvent Intent Model
**Location**: `internal/cli/intents/capture/model.go`

```go
type CaptureEventIntent struct {
    state    CaptureState
    form     *CaptureForm
    review   *ReviewState
    result   *IntentResult[CaptureEventResult]
}

type CaptureState string
const (
    StateChooseStrategy    CaptureState = "choose_strategy"
    StateCaptureForm      CaptureState = "capture_form"
    StateReviewInferred   CaptureState = "review_inferred"
    StateSubmit           CaptureState = "submit"
)

type CaptureEventResult struct {
    Event     *domain.CareerEvent
    Timestamp time.Time
}
```

**Checklist**:
- [ ] Define all state constants
- [ ] Create `CaptureForm` model with validation
- [ ] Create `ReviewState` model for inline editing
- [ ] Implement state transition logic
- [ ] Add form validation helpers
- [ ] Write unit tests for model

**Estimated Time**: 3-4 days

---

#### 2.2 Implement CaptureEvent State Transitions
**Location**: `internal/cli/intents/capture/update.go`

Implement transitions for:
- `StateChooseStrategy` → `StateCaptureForm`
- `StateCaptureForm` → `StateReviewInferred`
- `StateReviewInferred` → `StateSubmit`
- Sub-flow handling: `EditMetadata`, `EditBursts`, `EditFacts`
- Error handling and validation

**Checklist**:
- [ ] Implement all state transitions
- [ ] Handle form input and validation
- [ ] Implement modal sub-flows for editing
- [ ] Handle user cancellation
- [ ] Handle submission and domain service calls
- [ ] Write comprehensive unit tests for each transition
- [ ] Test error cases and recovery

**Estimated Time**: 4-5 days

---

#### 2.3 Implement CaptureEvent Views
**Location**: `internal/cli/intents/capture/view.go`

Create views for each state:
- Strategy selection view
- Form capture view
- Review/inline editing view
- Result confirmation view

**Checklist**:
- [ ] Implement view for each state
- [ ] Add proper styling and formatting
- [ ] Implement help text and instructions
- [ ] Add error message display
- [ ] Test view rendering for all states
- [ ] Ensure keyboard navigation works

**Estimated Time**: 3-4 days

---

#### 2.4 Implement Modal Sub-Flows
**Location**: `internal/cli/intents/capture/modals/`

Create modal components for:
- `EditMetadataModal` → `ModalEditResult[Metadata]`
- `EditBurstModal` → `ModalEditResult[Burst]`
- `EditFactModal` → `ModalEditResult[Fact]`

**Checklist**:
- [ ] Create modal base component
- [ ] Implement metadata edit modal
- [ ] Implement burst edit modal
- [ ] Implement fact edit modal
- [ ] Ensure context preservation on cancel
- [ ] Test modal lifecycle and result propagation

**Estimated Time**: 3-4 days

---

#### 2.5 Add Comprehensive Tests
**Location**: `internal/cli/intents/capture/` (all `*_test.go` files)

- Unit tests for model state transitions
- Unit tests for view rendering
- Integration tests for complete workflows
- Property-based tests for invariants
- Modal sub-flow tests

**Checklist**:
- [ ] Write unit tests for all state transitions
- [ ] Write integration tests for full workflows
- [ ] Write property-based tests for invariants
- [ ] Achieve >90% code coverage
- [ ] Test error cases and edge cases
- [ ] Test back navigation and state restoration

**Estimated Time**: 3-4 days

---

### Phase 2 Validation

**Acceptance Criteria**:
- [ ] `CaptureEvent` intent fully implements `Intent` interface
- [ ] All state transitions work correctly
- [ ] Form validation prevents invalid data
- [ ] Modal sub-flows preserve context and return typed results
- [ ] Error handling is consistent and user-friendly
- [ ] All tests pass with >90% coverage
- [ ] Back navigation restores complete state

**Estimated Phase 2 Duration**: 2 weeks

---

## Phase 3: Remaining Core Intents (Weeks 5-8)

### Goal
Implement the remaining four core intents, replicating the CaptureEvent pattern.

### 3.1 BrowseTimeline Intent
**Location**: `internal/cli/intents/browse/`

**States**:
- `StateTimelineView` (with Filter, Sort, Select sub-states)
- `StateEventDetail`

**Pattern Replication**:
- [ ] Define model with states
- [ ] Implement state transitions
- [ ] Implement views
- [ ] Add metadata preservation for scroll position, filters, sort order
- [ ] Add comprehensive tests

**Estimated Time**: 1.5 weeks

---

### 3.2 GenerateCV Intent
**Location**: `internal/cli/intents/generate_cv/`

**States**:
- `StateSelectProfile`
- `StateValidateProfile`
- `StateSelectAudience`
- `StateValidateAudience`
- `StateGeneratePreview`
- `StateReviewCV`
- `StateConfirmCV`
- `StateArtifactReady`

**Pattern Replication**:
- [ ] Define model with states and validation
- [ ] Implement all state transitions
- [ ] Implement views for each state
- [ ] Add inline edit mode for review state
- [ ] Add validation with clear error messages
- [ ] Add comprehensive tests

**Estimated Time**: 2 weeks

---

### 3.3 ExportArtifact Intent
**Location**: `internal/cli/intents/export/`

**States**:
- `StateSelectArtifactType`
- `StateConfigureExport`
- `StatePreviewExport`
- `StateConfirmExport`
- `StateExportInProgress`

**Pattern Replication**:
- [ ] Define model with states
- [ ] Implement state transitions
- [ ] Implement async operation pattern
- [ ] Add retry logic for failures
- [ ] Add progress indication
- [ ] Add comprehensive tests including async handling

**Estimated Time**: 1.5 weeks

---

### 3.4 ConfigureSystem Intent
**Location**: `internal/cli/intents/configure/`

**States**:
- `StateSelectDomain`
- `StateEditSettings`
- `StateStageChanges`
- `StateSaveConfiguration`

**Pattern Replication**:
- [ ] Define model with states
- [ ] Implement staged changes pattern
- [ ] Implement validation
- [ ] Implement state transitions
- [ ] Add discard changes functionality
- [ ] Add comprehensive tests

**Estimated Time**: 1 week

---

### Phase 3 Validation

**Acceptance Criteria**:
- [ ] All four intents fully implement `Intent` interface
- [ ] All intents follow the established CaptureEvent pattern
- [ ] All intents have consistent error handling
- [ ] All intents have comprehensive test coverage (>90%)
- [ ] Navigation between intents works correctly
- [ ] Back navigation works for all intents
- [ ] No global state mutations occur

**Estimated Phase 3 Duration**: 4 weeks

---

## Phase 4: Integration & Polish (Weeks 9-10)

### Goal
Integrate all intents, add cross-intent features, and polish the UX.

### Tasks

#### 4.1 Integrate All Intents
- [ ] Register all intents with `IntentRouter`
- [ ] Update main menu to activate intents correctly
- [ ] Test navigation between all intents
- [ ] Verify back navigation works across all intents
- [ ] Test result propagation from all intents

**Estimated Time**: 2-3 days

---

#### 4.2 Add Global Shortcuts
- [ ] Implement quit/exit from any intent
- [ ] Implement help system accessible from any intent
- [ ] Implement main menu shortcut (Ctrl+Home)
- [ ] Test keyboard shortcuts across all intents

**Estimated Time**: 1-2 days

---

#### 4.3 Add Comprehensive Logging
- [ ] Add intent activation logging
- [ ] Add state transition logging
- [ ] Add error logging with context
- [ ] Add result propagation logging
- [ ] Add performance metrics

**Estimated Time**: 2-3 days

---

#### 4.4 Performance Optimization
- [ ] Profile rendering performance
- [ ] Optimize view rendering
- [ ] Add caching where appropriate
- [ ] Test with large datasets

**Estimated Time**: 2-3 days

---

#### 4.5 Documentation & Examples
- [ ] Create intent implementation guide
- [ ] Create testing guide for intents
- [ ] Create troubleshooting guide
- [ ] Add code examples for common patterns
- [ ] Update README with new architecture

**Estimated Time**: 2-3 days

---

### Phase 4 Validation

**Acceptance Criteria**:
- [ ] All intents are registered and accessible
- [ ] Navigation is smooth and predictable
- [ ] All keyboard shortcuts work
- [ ] Error handling is consistent across all intents
- [ ] Performance is acceptable
- [ ] Documentation is complete and clear
- [ ] No regressions in existing functionality

**Estimated Phase 4 Duration**: 2 weeks

---

## Phase 5: Secondary Intents & Future Expansion (Post-Release)

### Goal
Establish patterns for secondary/contextual intents and prepare for future expansion.

### Secondary Intents to Consider

#### 5.1 Skill Tracking (Contextual Intent)
- Accessed from: CaptureEvent (ReviewInferredEvent), BrowseTimeline (EventDetail)
- Pattern: Modal sub-flow returning `ModalEditResult[Skill]`
- Implementation: Similar to burst/fact editing

#### 5.2 Career Goal Setting (Top-Level or Contextual)
- Accessed from: Main menu or GenerateCV
- Pattern: Full intent with multiple states
- Implementation: Similar to GenerateCV intent

#### 5.3 Mentor Matching (Contextual Intent)
- Accessed from: GenerateCV (ReviewCV), BrowseTimeline
- Pattern: Modal sub-flow with recommendations
- Implementation: Similar to burst/fact editing

#### 5.4 Continuous Learning Tracking (Top-Level or Contextual)
- Accessed from: Main menu or CaptureEvent
- Pattern: Full intent with timeline view
- Implementation: Similar to BrowseTimeline intent

### Promotion Criteria
- Promote secondary intent to top-level only when:
  - UX is validated with users
  - Demand is clear from usage metrics
  - Integration with core workflows is complete
  - Performance impact is acceptable

---

## Implementation Checklist

### Pre-Implementation
- [ ] Review [TUI_INTENT_DIAGRAM.md](TUI_INTENT_DIAGRAM.md) thoroughly
- [ ] Review this roadmap with team
- [ ] Set up development environment
- [ ] Create feature branches for each phase

### Phase 1: Foundation
- [ ] Implement intent boundary contract types
- [ ] Implement `IntentRouter`
- [ ] Refactor root model
- [ ] All Phase 1 tests pass

### Phase 2: CaptureEvent
- [ ] Implement CaptureEvent intent model
- [ ] Implement state transitions
- [ ] Implement views
- [ ] Implement modal sub-flows
- [ ] All Phase 2 tests pass

### Phase 3: Remaining Intents
- [ ] Implement BrowseTimeline intent
- [ ] Implement GenerateCV intent
- [ ] Implement ExportArtifact intent
- [ ] Implement ConfigureSystem intent
- [ ] All Phase 3 tests pass

### Phase 4: Integration & Polish
- [ ] Integrate all intents
- [ ] Add global shortcuts
- [ ] Add comprehensive logging
- [ ] Optimize performance
- [ ] Complete documentation
- [ ] All Phase 4 tests pass

### Phase 5: Secondary Intents
- [ ] Design secondary intent patterns
- [ ] Implement skill tracking
- [ ] Implement career goal setting
- [ ] Implement mentor matching
- [ ] Implement continuous learning
- [ ] All Phase 5 tests pass

---

## Testing Strategy Summary

### Unit Tests
- Test each state transition in isolation
- Test view rendering for each state
- Test input validation
- Test error handling

### Integration Tests
- Test complete workflows from entry to exit
- Test intent switching and navigation
- Test back navigation with state restoration
- Test result propagation

### Property-Based Tests
- Verify no intent mutates global state
- Verify no illegal transitions occur
- Verify all results are strongly typed
- Verify back navigation restores complete context

### Test Tools
- `IntentTestHarness`: Isolated intent testing
- `IntentRouterTestHelper`: Router testing
- Mock domain services for testing
- Property-based testing framework (e.g., gopter)

---

## Success Metrics

### Code Quality
- [ ] All code has >90% test coverage
- [ ] All linting checks pass
- [ ] No type safety violations
- [ ] Clear, documented code

### Architecture
- [ ] No global state mutations
- [ ] All intent boundaries enforced
- [ ] All state transitions are type-safe
- [ ] Clear separation of concerns

### User Experience
- [ ] Smooth, predictable navigation
- [ ] Clear error messages
- [ ] Responsive UI
- [ ] Intuitive keyboard shortcuts

### Performance
- [ ] Rendering completes in <100ms
- [ ] No memory leaks
- [ ] Handles large datasets efficiently
- [ ] Smooth scrolling and transitions

---

## Risk Mitigation

### Risk: Scope Creep
- **Mitigation**: Stick to the defined phases. New features go to Phase 5.
- **Monitor**: Weekly progress reviews against this roadmap.

### Risk: Test Coverage Gaps
- **Mitigation**: Require >90% coverage for all code. Use property-based tests for invariants.
- **Monitor**: Coverage reports in CI/CD pipeline.

### Risk: Performance Degradation
- **Mitigation**: Profile regularly. Add performance benchmarks.
- **Monitor**: Performance metrics in CI/CD pipeline.

### Risk: Breaking Changes
- **Mitigation**: Use feature flags for new intents. Maintain backward compatibility.
- **Monitor**: Integration tests catch regressions early.

---

## Timeline Summary

| Phase | Duration | Key Deliverable |
|-------|----------|-----------------|
| 1: Foundation | 1.5 weeks | IntentRouter + Root Model Refactor |
| 2: CaptureEvent | 2 weeks | Complete CaptureEvent Intent |
| 3: Remaining Intents | 4 weeks | BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem |
| 4: Integration & Polish | 2 weeks | Integrated, Tested, Documented System |
| 5: Secondary Intents | Post-Release | Skill Tracking, Career Goals, Mentor Matching, Continuous Learning |

**Total Estimated Timeline**: 9.5 weeks for core implementation

---

## References

- [TUI Intent & Flow State Diagram](TUI_INTENT_DIAGRAM.md) - Complete architectural specification
- [WORKFLOW_DIAGRAM.md](WORKFLOW_DIAGRAM.md) - High-level workflow overview
- [TUI_DEVELOPER_GUIDE.md](TUI_DEVELOPER_GUIDE.md) - General TUI development guidelines
- [TUI_STANDARDS.md](TUI_STANDARDS.md) - UI/UX standards and conventions

---

## Next Steps

1. **Week 1**: Review this roadmap with the team
2. **Week 1**: Set up development environment and feature branches
3. **Week 1-2**: Begin Phase 1 implementation (IntentRouter + Root Model)
4. **Week 3-4**: Begin Phase 2 implementation (CaptureEvent)
5. **Week 5-8**: Begin Phase 3 implementation (Remaining Intents)
6. **Week 9-10**: Begin Phase 4 implementation (Integration & Polish)

---

*Last Updated: 2026-01-02*
*Status: Ready for Implementation*

