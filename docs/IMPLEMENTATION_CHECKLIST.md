---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya TUI Intent Architecture - Implementation Checklist

This document provides a detailed, phase-by-phase checklist for implementing the production-ready TUI intent architecture defined in [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md).

---

## Phase 1: Foundation & Core Infrastructure

### 1.1 Intent Boundary Contract Types

**File**: `internal/cli/intents/contract.go`

#### Core Types Implementation
- [ ] Define `IntentStatus` string type with constants:
  - [ ] `StatusCompleted = "completed"`
  - [ ] `StatusCancelled = "cancelled"`
  - [ ] `StatusFailed = "failed"`
  - [ ] `StatusPartial = "partial"`

- [ ] Implement `IntentError` struct:
  - [ ] `Code string` - machine-readable error code
  - [ ] `Message string` - user-facing error message
  - [ ] `Cause error` - underlying error for logging
  - [ ] Constructor: `NewIntentError(code, message string, cause error) *IntentError`
  - [ ] Error() method for `error` interface compliance

- [ ] Implement generic `IntentResult[T]` struct:
  - [ ] `Status IntentStatus`
  - [ ] `Data T`
  - [ ] `Error *IntentError`
  - [ ] `Metadata map[string]interface{}`
  - [ ] Helper method: `WithMetadata(key string, value interface{}) *IntentResult[T]`
  - [ ] Helper method: `GetMetadata(key string) (interface{}, bool)`
  - [ ] Helper method: `IsSuccess() bool` (returns Status == StatusCompleted)
  - [ ] Helper method: `IsPartialSuccess() bool` (returns Status == StatusPartial)
  - [ ] Helper method: `IsFailed() bool` (returns Status == StatusFailed)
  - [ ] Helper method: `IsCancelled() bool` (returns Status == StatusCancelled)

- [ ] Implement `Intent` interface:
  - [ ] `Init(ctx context.Context) tea.Cmd`
  - [ ] `Update(msg tea.Msg) (Intent, tea.Cmd)`
  - [ ] `View() string`
  - [ ] `Result() *IntentResult[interface{}]`

- [ ] Implement generic `ModalEditResult[T]` struct:
  - [ ] `Original T`
  - [ ] `Modified T`
  - [ ] `Accepted bool`
  - [ ] `Changes map[string]interface{}`
  - [ ] Constructor: `NewModalEditResult(original, modified T) *ModalEditResult[T]`
  - [ ] Helper method: `WithChange(field string, value interface{}) *ModalEditResult[T]`
  - [ ] Helper method: `Accept() *ModalEditResult[T]`
  - [ ] Helper method: `Reject() *ModalEditResult[T]`

#### Testing
- [ ] Unit tests for all types
- [ ] Test `IntentError` wrapping and logging
- [ ] Test `IntentResult[T]` metadata operations
- [ ] Test `ModalEditResult[T]` state transitions
- [ ] Verify all helper methods work correctly
- [ ] Test with multiple generic types (string, struct, slice)

#### Acceptance Criteria
- [ ] All types compile without errors
- [ ] All helper methods have >90% test coverage
- [ ] No runtime type assertions needed
- [ ] Clear error messages on type mismatches

---

### 1.2 IntentRouter Implementation

**File**: `internal/cli/intents/router.go`

#### Core Router Structure
- [ ] Define `IntentRouter` struct:
  - [ ] `intents map[string]Intent` - registered intents by name
  - [ ] `current Intent` - currently active intent
  - [ ] `history []Intent` - navigation history for back navigation
  - [ ] `ctx context.Context` - shared context
  - [ ] `resultHandler func(IntentResult[interface{}])` - result callback

#### Router Methods
- [ ] `NewIntentRouter(ctx context.Context) *IntentRouter`
  - [ ] Initialize empty intents map
  - [ ] Initialize empty history
  - [ ] Store context

- [ ] `Register(name string, intent Intent) error`
  - [ ] Validate intent implements `Intent` interface
  - [ ] Prevent duplicate registrations
  - [ ] Return error on failure

- [ ] `Activate(intentName string) tea.Cmd`
  - [ ] Find intent by name
  - [ ] Call `Init()` on intent
  - [ ] Store previous intent in history
  - [ ] Set as current intent
  - [ ] Return init command

- [ ] `Update(msg tea.Msg) (tea.Cmd, error)`
  - [ ] Delegate to current intent
  - [ ] Check if result is ready
  - [ ] Call result handler if ready
  - [ ] Return command

- [ ] `View() string`
  - [ ] Return current intent's view
  - [ ] Handle nil current intent gracefully

- [ ] `Back() tea.Cmd`
  - [ ] Check if history is empty
  - [ ] Pop previous intent from history
  - [ ] Call `Init()` on restored intent
  - [ ] Restore metadata from previous result
  - [ ] Return command

- [ ] `GetCurrentIntent() Intent`
  - [ ] Return currently active intent
  - [ ] Used for testing

- [ ] `GetHistory() []Intent`
  - [ ] Return navigation history
  - [ ] Used for testing and debugging

#### Testing
- [ ] Unit tests for all router methods
- [ ] Test intent registration
- [ ] Test intent activation and initialization
- [ ] Test intent switching
- [ ] Test back navigation with history
- [ ] Test result propagation
- [ ] Test error handling for missing intents
- [ ] Test concurrent intent operations

#### Acceptance Criteria
- [ ] Router successfully activates and switches intents
- [ ] Back navigation restores previous intent
- [ ] Result handler is called when intent completes
- [ ] All tests pass with >90% coverage
- [ ] No goroutine leaks or race conditions

---

### 1.3 Root Model Refactoring

**File**: `internal/cli/app/app.go`

#### Current State Analysis
- [ ] Document current root model structure
- [ ] Identify all inline state that should move to intents
- [ ] Document all message types
- [ ] Identify global state mutations

#### Refactoring Steps
- [ ] Replace inline state with `IntentRouter`
  - [ ] Remove all intent-specific state from root model
  - [ ] Keep only global UI state (e.g., theme, preferences)
  - [ ] Add `router *IntentRouter` field

- [ ] Update `Update()` method:
  - [ ] Delegate all intent messages to router
  - [ ] Handle router result callbacks
  - [ ] Handle global shortcuts (quit, help, main menu)
  - [ ] Return appropriate commands

- [ ] Update `View()` method:
  - [ ] Call `router.View()`
  - [ ] Add global UI chrome (status bar, help text)
  - [ ] Preserve existing visual style

- [ ] Update `Init()` method:
  - [ ] Initialize router
  - [ ] Register all intents
  - [ ] Activate main menu or initial intent
  - [ ] Return appropriate commands

#### Testing
- [ ] Integration tests for root model
- [ ] Test intent activation from root
- [ ] Test result handling from intents
- [ ] Test back navigation from intents
- [ ] Test global shortcuts
- [ ] Verify no regressions in existing functionality

#### Acceptance Criteria
- [ ] Root model successfully delegates to router
- [ ] All intents activate and return results correctly
- [ ] Back navigation works from all intents
- [ ] No breaking changes to existing CLI functionality
- [ ] All tests pass

---

### Phase 1 Validation Checklist

**Before Moving to Phase 2**:
- [ ] All intent boundary contract types implemented and tested
- [ ] `IntentRouter` fully functional with comprehensive tests
- [ ] Root model successfully refactored to use router
- [ ] All Phase 1 tests pass with >90% coverage
- [ ] No type safety violations
- [ ] Code review completed
- [ ] Documentation updated

---

## Phase 2: CaptureEvent Intent Implementation

### 2.1 Intent Model Definition

**File**: `internal/cli/intents/capture/model.go`

#### State Definition
- [ ] Define `CaptureState` type with constants:
  - [ ] `StateChooseStrategy = "choose_strategy"`
  - [ ] `StateCaptureForm = "capture_form"`
  - [ ] `StateReviewInferred = "review_inferred"`
  - [ ] `StateSubmit = "submit"`

- [ ] Define `CaptureEventIntent` struct:
  - [ ] `state CaptureState`
  - [ ] `form *CaptureForm`
  - [ ] `review *ReviewState`
  - [ ] `result *IntentResult[CaptureEventResult]`
  - [ ] `domainService *domain.CareerEventService`

- [ ] Define `CaptureEventResult` struct:
  - [ ] `Event *domain.CareerEvent`
  - [ ] `Timestamp time.Time`
  - [ ] `SourceStrategy string` (manual, inference, etc.)

- [ ] Define `CaptureForm` struct:
  - [ ] Form fields for event capture
  - [ ] Validation state
  - [ ] Error messages
  - [ ] Helper methods for field access

- [ ] Define `ReviewState` struct:
  - [ ] Inferred event data
  - [ ] Metadata for editing
  - [ ] Bursts for editing
  - [ ] Facts for editing
  - [ ] Current selection
  - [ ] Edit mode state

#### Initialization
- [ ] `NewCaptureEventIntent(service *domain.CareerEventService) *CaptureEventIntent`
  - [ ] Initialize with clean state
  - [ ] Set initial state to `StateChooseStrategy`

#### Testing
- [ ] Unit tests for model initialization
- [ ] Test state transitions
- [ ] Test form validation
- [ ] Test review state management

---

### 2.2 State Transitions

**File**: `internal/cli/intents/capture/update.go`

#### Message Types
- [ ] Define message types for capture workflow:
  - [ ] `ChooseCaptureStrategyMsg`
  - [ ] `CaptureFormInputMsg`
  - [ ] `ValidateCaptureFormMsg`
  - [ ] `ReviewInferredEventMsg`
  - [ ] `EditMetadataMsg`
  - [ ] `EditBurstMsg`
  - [ ] `EditFactMsg`
  - [ ] `SubmitCaptureMsg`
  - [ ] `CancelCaptureMsg`

#### Update Logic
- [ ] Implement `Update(msg tea.Msg) (Intent, tea.Cmd)`:
  - [ ] Route messages based on current state
  - [ ] Call appropriate state handler
  - [ ] Return updated intent and command

- [ ] Implement state handlers:
  - [ ] `handleChooseStrategy(msg tea.Msg) (Intent, tea.Cmd)`
  - [ ] `handleCaptureForm(msg tea.Msg) (Intent, tea.Cmd)`
  - [ ] `handleReviewInferred(msg tea.Msg) (Intent, tea.Cmd)`
  - [ ] `handleSubmit(msg tea.Msg) (Intent, tea.Cmd)`

#### Form Validation
- [ ] Implement form field validation:
  - [ ] Required field validation
  - [ ] Date/time validation
  - [ ] Format validation
  - [ ] Clear error messages

- [ ] Implement `ValidateForm() error`:
  - [ ] Check all required fields
  - [ ] Return first error encountered
  - [ ] Store error in form state

#### Domain Service Integration
- [ ] Implement event inference:
  - [ ] Call domain service to infer bursts/facts
  - [ ] Handle inference errors gracefully
  - [ ] Store inferred data in review state

- [ ] Implement event submission:
  - [ ] Call domain service to save event
  - [ ] Handle submission errors
  - [ ] Return typed result

#### Modal Sub-Flow Handling
- [ ] Implement metadata edit handling:
  - [ ] Activate metadata edit modal
  - [ ] Handle modal result
  - [ ] Update form state

- [ ] Implement burst edit handling:
  - [ ] Activate burst edit modal
  - [ ] Handle modal result
  - [ ] Update review state

- [ ] Implement fact edit handling:
  - [ ] Activate fact edit modal
  - [ ] Handle modal result
  - [ ] Update review state

#### Testing
- [ ] Unit tests for each state handler
- [ ] Test form validation
- [ ] Test domain service integration
- [ ] Test modal sub-flow handling
- [ ] Test error cases and recovery
- [ ] Test cancellation at each state

---

### 2.3 Intent Views

**File**: `internal/cli/intents/capture/view.go`

#### View Methods
- [ ] Implement `View() string`:
  - [ ] Route to state-specific view
  - [ ] Return formatted string

- [ ] Implement `viewChooseStrategy() string`:
  - [ ] Display strategy options
  - [ ] Show selection indicator
  - [ ] Show help text

- [ ] Implement `viewCaptureForm() string`:
  - [ ] Display form fields
  - [ ] Show current input
  - [ ] Show validation errors
  - [ ] Show help text and navigation hints

- [ ] Implement `viewReviewInferred() string`:
  - [ ] Display inferred event data
  - [ ] Show metadata with edit indicator
  - [ ] Show bursts list with edit indicators
  - [ ] Show facts list with edit indicators
  - [ ] Show navigation options

- [ ] Implement `viewSubmit() string`:
  - [ ] Display confirmation message
  - [ ] Show saved event summary
  - [ ] Show next steps

#### Styling & Formatting
- [ ] Use consistent color scheme (from TUI_STANDARDS.md)
- [ ] Use consistent typography
- [ ] Implement proper indentation and spacing
- [ ] Add visual indicators for errors
- [ ] Add visual indicators for success

#### Testing
- [ ] Unit tests for view rendering
- [ ] Test view output for each state
- [ ] Test view with various data sizes
- [ ] Test view with error states
- [ ] Visual regression testing

---

### 2.4 Modal Sub-Flows

**File**: `internal/cli/intents/capture/modals/`

#### Metadata Edit Modal
**File**: `edit_metadata.go`

- [ ] Define `EditMetadataModal` struct:
  - [ ] `original *domain.Metadata`
  - [ ] `modified *domain.Metadata`
  - [ ] `state ModalState`
  - [ ] Form fields for editing

- [ ] Implement `Intent` interface:
  - [ ] `Init(ctx context.Context) tea.Cmd`
  - [ ] `Update(msg tea.Msg) (Intent, tea.Cmd)`
  - [ ] `View() string`
  - [ ] `Result() *IntentResult[interface{}]`

- [ ] Implement result generation:
  - [ ] Return `ModalEditResult[Metadata]`
  - [ ] Track changes in `Changes` map
  - [ ] Preserve original on cancel

#### Burst Edit Modal
**File**: `edit_burst.go`

- [ ] Define `EditBurstModal` struct similar to metadata
- [ ] Implement all required methods
- [ ] Return `ModalEditResult[Burst]`

#### Fact Edit Modal
**File**: `edit_fact.go`

- [ ] Define `EditFactModal` struct similar to metadata
- [ ] Implement all required methods
- [ ] Return `ModalEditResult[Fact]`

#### Testing
- [ ] Unit tests for each modal
- [ ] Test modal lifecycle
- [ ] Test form validation in modals
- [ ] Test result generation
- [ ] Test cancellation preserves original state

---

### 2.5 Comprehensive Tests

**File**: `internal/cli/intents/capture/*_test.go`

#### Unit Tests
- [ ] Test model initialization
- [ ] Test all state transitions
- [ ] Test form validation
- [ ] Test view rendering for each state
- [ ] Test modal sub-flows
- [ ] Test error handling
- [ ] Achieve >90% code coverage

#### Integration Tests
- [ ] Test complete capture workflow
- [ ] Test with valid input
- [ ] Test with invalid input
- [ ] Test with inference errors
- [ ] Test modal sub-flow integration
- [ ] Test cancellation at each state
- [ ] Test back navigation

#### Property-Based Tests
- [ ] Verify no global state mutation
- [ ] Verify no illegal transitions
- [ ] Verify all results are strongly typed
- [ ] Verify context preservation in modals

#### Test Utilities
- [ ] Create `CaptureEventTestHarness`:
  - [ ] Isolated intent testing
  - [ ] Mock domain service
  - [ ] Helper methods for state setup

#### Performance Tests
- [ ] Test rendering performance with large datasets
- [ ] Test form validation performance
- [ ] Test domain service call performance

---

### Phase 2 Validation Checklist

**Before Moving to Phase 3**:
- [ ] `CaptureEvent` intent fully implements `Intent` interface
- [ ] All state transitions work correctly
- [ ] Form validation prevents invalid data
- [ ] Modal sub-flows preserve context and return typed results
- [ ] Error handling is consistent and user-friendly
- [ ] All tests pass with >90% coverage
- [ ] Back navigation restores complete state
- [ ] Code review completed
- [ ] Documentation updated

---

## Phase 3: Remaining Core Intents

### 3.1 BrowseTimeline Intent

**File**: `internal/cli/intents/browse/`

#### Checklist
- [ ] Define state machine with states:
  - [ ] `StateTimelineView` (with Filter, Sort, Select sub-states)
  - [ ] `StateEventDetail`

- [ ] Implement model with:
  - [ ] Timeline data structure
  - [ ] Filter state
  - [ ] Sort state
  - [ ] Selection state
  - [ ] Scroll position (for metadata)

- [ ] Implement state transitions:
  - [ ] Filter application
  - [ ] Sort application
  - [ ] Event selection
  - [ ] Event detail view
  - [ ] Back to timeline

- [ ] Implement views:
  - [ ] Timeline view with filters/sort
  - [ ] Event detail view
  - [ ] Proper scrolling and pagination

- [ ] Implement metadata preservation:
  - [ ] Store scroll position
  - [ ] Store filter state
  - [ ] Store sort order
  - [ ] Store selection
  - [ ] Restore on back navigation

- [ ] Implement tests:
  - [ ] Unit tests for all transitions
  - [ ] Integration tests for workflows
  - [ ] Property-based tests for invariants
  - [ ] Performance tests for large datasets

---

### 3.2 GenerateCV Intent

**File**: `internal/cli/intents/generate_cv/`

#### Checklist
- [ ] Define state machine with states:
  - [ ] `StateSelectProfile`
  - [ ] `StateValidateProfile`
  - [ ] `StateSelectAudience`
  - [ ] `StateValidateAudience`
  - [ ] `StateGeneratePreview`
  - [ ] `StateReviewCV`
  - [ ] `StateConfirmCV`
  - [ ] `StateArtifactReady`

- [ ] Implement model with:
  - [ ] Profile selection
  - [ ] Audience selection
  - [ ] CV preview data
  - [ ] Edit mode state (for review)
  - [ ] Validation results

- [ ] Implement state transitions:
  - [ ] Profile selection and validation
  - [ ] Audience selection and validation
  - [ ] Preview generation
  - [ ] Review with edit mode toggle
  - [ ] Confirmation
  - [ ] Artifact generation

- [ ] Implement validation:
  - [ ] Profile completeness check
  - [ ] Audience compatibility check
  - [ ] Required data validation
  - [ ] Clear error messages

- [ ] Implement views:
  - [ ] Profile selector view
  - [ ] Audience selector view
  - [ ] CV preview view
  - [ ] Edit mode toggle in review
  - [ ] Confirmation view
  - [ ] Success view

- [ ] Implement inline editing:
  - [ ] Toggle edit mode in review
  - [ ] Edit CV content
  - [ ] Validate changes
  - [ ] Apply changes or cancel

- [ ] Implement tests:
  - [ ] Unit tests for all transitions
  - [ ] Validation tests
  - [ ] Integration tests for workflows
  - [ ] Property-based tests
  - [ ] Performance tests for large profiles

---

### 3.3 ExportArtifact Intent

**File**: `internal/cli/intents/export/`

#### Checklist
- [ ] Define state machine with states:
  - [ ] `StateSelectArtifactType`
  - [ ] `StateConfigureExport`
  - [ ] `StatePreviewExport`
  - [ ] `StateConfirmExport`
  - [ ] `StateExportInProgress`

- [ ] Implement model with:
  - [ ] Artifact type selection
  - [ ] Export configuration
  - [ ] Preview data
  - [ ] Progress state (for async)
  - [ ] Error state (for retry)

- [ ] Implement state transitions:
  - [ ] Artifact type selection
  - [ ] Export configuration
  - [ ] Preview generation
  - [ ] Confirmation
  - [ ] Async export operation
  - [ ] Completion or error

- [ ] Implement async operations:
  - [ ] Use ephemeral `InProgress` state
  - [ ] Trigger domain service call
  - [ ] Handle completion
  - [ ] Handle errors with retry logic
  - [ ] Handle timeouts with backoff

- [ ] Implement views:
  - [ ] Artifact type selector
  - [ ] Configuration form
  - [ ] Preview view
  - [ ] Confirmation view
  - [ ] Progress indicator
  - [ ] Success/error view

- [ ] Implement tests:
  - [ ] Unit tests for all transitions
  - [ ] Async operation tests
  - [ ] Error handling tests
  - [ ] Retry logic tests
  - [ ] Integration tests
  - [ ] Property-based tests

---

### 3.4 ConfigureSystem Intent

**File**: `internal/cli/intents/configure/`

#### Checklist
- [ ] Define state machine with states:
  - [ ] `StateSelectDomain`
  - [ ] `StateEditSettings`
  - [ ] `StateStageChanges`
  - [ ] `StateSaveConfiguration`

- [ ] Implement model with:
  - [ ] Configuration domain selection
  - [ ] Original settings
  - [ ] Staged changes
  - [ ] Validation state

- [ ] Implement state transitions:
  - [ ] Domain selection
  - [ ] Settings editing
  - [ ] Change staging
  - [ ] Validation
  - [ ] Save or discard

- [ ] Implement staged changes:
  - [ ] Track all changes locally
  - [ ] No partial writes
  - [ ] Discard all on cancel
  - [ ] Validate before save

- [ ] Implement views:
  - [ ] Domain selector
  - [ ] Settings editor
  - [ ] Change preview
  - [ ] Confirmation view
  - [ ] Success/error view

- [ ] Implement tests:
  - [ ] Unit tests for all transitions
  - [ ] Staged changes tests
  - [ ] Validation tests
  - [ ] Integration tests
  - [ ] Property-based tests

---

### Phase 3 Validation Checklist

**Before Moving to Phase 4**:
- [ ] All four intents fully implement `Intent` interface
- [ ] All intents follow the established CaptureEvent pattern
- [ ] All intents have consistent error handling
- [ ] All intents have comprehensive test coverage (>90%)
- [ ] Navigation between intents works correctly
- [ ] Back navigation works for all intents
- [ ] No global state mutations occur
- [ ] Code review completed
- [ ] Documentation updated

---

## Phase 4: Integration & Polish

### 4.1 Integration Checklist
- [ ] Register all intents with `IntentRouter`
- [ ] Update main menu to activate intents
- [ ] Test navigation between all intents
- [ ] Verify back navigation across intents
- [ ] Test result propagation from all intents
- [ ] Verify no state leakage between intents

### 4.2 Global Shortcuts Checklist
- [ ] Implement quit/exit from any intent
- [ ] Implement help system
- [ ] Implement main menu shortcut
- [ ] Test shortcuts across all intents
- [ ] Document keyboard shortcuts

### 4.3 Logging Checklist
- [ ] Add intent activation logging
- [ ] Add state transition logging
- [ ] Add error logging with context
- [ ] Add result propagation logging
- [ ] Add performance metrics logging

### 4.4 Performance Optimization Checklist
- [ ] Profile rendering performance
- [ ] Optimize view rendering
- [ ] Add caching where appropriate
- [ ] Test with large datasets
- [ ] Document performance characteristics

### 4.5 Documentation Checklist
- [ ] Create intent implementation guide
- [ ] Create testing guide
- [ ] Create troubleshooting guide
- [ ] Add code examples
- [ ] Update README
- [ ] Update AGENTS.md

---

## Phase 5: Secondary Intents (Post-Release)

### Secondary Intent Patterns
- [ ] Design skill tracking intent
- [ ] Design career goal setting intent
- [ ] Design mentor matching intent
- [ ] Design continuous learning intent
- [ ] Establish promotion criteria

---

## Enhancements & Recommendations

### Cross-Intent Metadata (Recommended for Phase 1.5)

**Consideration**: Lightweight `GlobalContext` for non-mutable shared metadata

```go
type GlobalContext struct {
    UserPreferences map[string]interface{}
    LastSelectedProfile string
    LastSelectedAudience string
    // Non-mutable, read-only shared state
}
```

**Benefits**:
- Simplifies repeated back-and-forth data without violating intent boundaries
- Reduces metadata passing between intents
- Improves UX for common workflows

**Implementation**:
- [ ] Define `GlobalContext` struct in `internal/cli/intents/context.go`
- [ ] Pass via router to all intents
- [ ] Document usage patterns
- [ ] Add tests for context sharing

---

### Async Feedback in TUI (Recommended for Phase 4)

**Consideration**: Non-blocking spinner/progress bar for in-progress operations

**Implementation**:
- [ ] Create progress component in `internal/cli/components/progress/`
- [ ] Implement spinner animation
- [ ] Integrate with async operations
- [ ] Add progress percentage for long operations
- [ ] Add cancellation support

---

### CI/CD Integration (Recommended for Phase 1)

**Consideration**: Automatic verification of implementation phases

**Implementation**:
- [ ] Add linting checks to CI
- [ ] Add coverage reports to CI
- [ ] Add property-based tests to CI
- [ ] Add performance benchmarks to CI
- [ ] Create phase validation scripts

**Checklist**:
- [ ] Configure linter (golangci-lint)
- [ ] Configure coverage threshold (>90%)
- [ ] Configure property-based test runner
- [ ] Configure performance benchmarks
- [ ] Create phase validation CI job

---

### Performance Benchmarking (Recommended for Phase 1-2)

**Consideration**: Catch slow rendering early, especially for large datasets

**Implementation**:
- [ ] Create benchmark suite in `internal/cli/intents/benchmarks/`
- [ ] Add baseline rendering benchmarks
- [ ] Add form validation benchmarks
- [ ] Add domain service call benchmarks
- [ ] Track performance across phases

**Benchmarks to Include**:
- [ ] `BenchmarkCaptureEventView` - rendering performance
- [ ] `BenchmarkFormValidation` - validation performance
- [ ] `BenchmarkTimelineRendering` - large dataset rendering
- [ ] `BenchmarkCVGeneration` - CV generation performance

---

## Validation Gates

### Gate 1: Phase 1 Complete
- [ ] All intent boundary contract types implemented
- [ ] `IntentRouter` fully functional
- [ ] Root model refactored
- [ ] All Phase 1 tests pass with >90% coverage
- [ ] Code review approved
- [ ] Ready for Phase 2

### Gate 2: Phase 2 Complete
- [ ] `CaptureEvent` intent fully implemented
- [ ] All state transitions working
- [ ] Modal sub-flows functional
- [ ] All Phase 2 tests pass with >90% coverage
- [ ] Code review approved
- [ ] Ready for Phase 3

### Gate 3: Phase 3 Complete
- [ ] All four core intents implemented
- [ ] All intents follow CaptureEvent pattern
- [ ] All intents have >90% test coverage
- [ ] Navigation working between all intents
- [ ] Code review approved
- [ ] Ready for Phase 4

### Gate 4: Phase 4 Complete
- [ ] All intents integrated
- [ ] Global shortcuts working
- [ ] Logging comprehensive
- [ ] Performance acceptable
- [ ] Documentation complete
- [ ] Code review approved
- [ ] Ready for release

---

## Developer Quick Reference

### Creating a New Intent

1. **Create directory**: `internal/cli/intents/your_intent/`
2. **Create model**: `model.go` with state machine
3. **Create update**: `update.go` with state transitions
4. **Create view**: `view.go` with rendering logic
5. **Create tests**: `*_test.go` with comprehensive tests
6. **Register**: Add to `IntentRouter` in `app.go`

### Implementing State Transitions

```go
func (y *YourIntent) Update(msg tea.Msg) (Intent, tea.Cmd) {
    switch y.state {
    case StateInitial:
        return y.handleInitial(msg)
    case StateWorking:
        return y.handleWorking(msg)
    case StateFinal:
        return y.handleFinal(msg)
    default:
        return y, nil
    }
}
```

### Returning Results

```go
func (y *YourIntent) Result() *IntentResult[interface{}] {
    if y.result != nil {
        return y.result
    }
    return &IntentResult[interface{}]{
        Status: StatusCompleted,
        Data:   y.data,
    }
}
```

### Testing Intents

```go
func TestYourIntentWorkflow(t *testing.T) {
    harness := NewIntentTestHarness(NewYourIntent())
    harness.Update(yourMessage)
    assert.Equal(t, StateWorking, harness.State())
}
```

---

## Resources

- [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md) - Phase details
- [TUI_INTENT_DIAGRAM.md](../TUI_INTENT_DIAGRAM.md) - Architecture specification
- [TUI_STANDARDS.md](TUI_STANDARDS.md) - UI/UX standards
- [TUI_DEVELOPER_GUIDE.md](TUI_DEVELOPER_GUIDE.md) - General TUI guidelines

---

*Last Updated: 2026-01-02*
*Status: Ready for Implementation*

