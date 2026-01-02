# PRD: TUI Intent Architecture Refactoring

**Document Status**: Ready for Implementation
**Last Updated**: 2026-01-02
**Target Audience**: Development Team
**Confidence Level**: Very High

---

## 1. Introduction & Overview

### Problem Statement

The KaRiya TUI currently has inconsistent workflow patterns, state management issues, and scattered architectural approaches that make it difficult to maintain, test, and extend. The existing code contains:

- Mixed state management patterns (some local, some global)
- Unclear navigation boundaries between screens/intents
- Inconsistent error handling
- Difficulty testing individual workflows in isolation
- Challenges with back navigation and context restoration
- Type-unsafe intent communication

### Solution Overview

Refactor the entire TUI to strictly implement the **Intent-Driven Architecture** as defined in `docs/TUI_INTENT_DIAGRAM.md`. This refactoring will:

1. Establish **type-safe intent communication** via `IntentResult[T]`
2. Create **clear intent boundaries** that prevent state pollution
3. Implement **predictable state machines** for all workflows
4. Enable **reliable back navigation** with full context restoration
5. Support **modal sub-flows** for inline editing
6. Establish **consistent async operation patterns**
7. Make **illegal states unrepresentable** through the type system

### Goal

Transform the KaRiya TUI into a production-grade, type-safe system where:
- Each workflow is a self-contained intent with explicit state transitions
- All intents communicate via strongly-typed `IntentResult[T]`
- No intent mutates global state or assumes prior context
- Back navigation is reliable and preserves complete context
- Testing is straightforward with isolated, mockable intents
- The codebase is maintainable and extensible

### Success Definition

A **consistent, reliable workflow system** where:
- ✅ All intents follow the same architectural pattern
- ✅ No broken workflows or navigation issues
- ✅ All state transitions are type-safe and explicit
- ✅ Complete test coverage (>90%) with unit and E2E tests
- ✅ Clear, documented architecture that developers can follow
- ✅ Extensible pattern for adding new intents/workflows

---

## 2. Goals

### Primary Goals

1. **Establish Type-Safe Intent Communication**
   - Implement `IntentResult[T]` as the only mechanism for intent communication
   - Support four explicit states: Completed, Cancelled, Failed, Partial
   - Include `IntentError` for debug/logging without type pollution
   - Enable strongly-typed metadata for context passing

2. **Refactor All Intents to Follow Consistent Pattern**
   - Each intent is a self-contained state machine
   - All intents implement the `Intent` interface
   - All intents own only local state and navigation
   - All intents return typed results
   - All intents can be tested in isolation

3. **Implement Clear Intent Boundaries**
   - No cross-intent state mutation
   - No direct navigation between intents
   - Explicit ownership rules (MAY/MAY NOT)
   - Compile-time type safety prevents illegal states

4. **Enable Reliable Back Navigation**
   - Full context restoration (scroll position, filters, sort order, selection)
   - Type-safe metadata preservation
   - Automatic restoration via IntentRouter
   - No context leakage between intents

5. **Support Modal Sub-Flows Pattern**
   - `ModalEditResult[T]` for typed diffs
   - Context preservation on cancel
   - No global state mutation from modals
   - Consistent sub-flow handling across all intents

6. **Establish Async Operations Pattern**
   - Ephemeral `InProgress` state for non-blocking operations
   - Strongly-typed completion triggers
   - Consistent error handling and retry logic
   - Progress feedback for long operations

### Secondary Goals

1. **Add Recommended Enhancements**
   - GlobalContext pattern for non-mutable shared metadata
   - Async feedback in TUI (progress indicators, spinners)
   - CI/CD integration with automated checks
   - Performance benchmarking and optimization

2. **Improve Developer Experience**
   - Clear documentation and examples
   - Consistent naming conventions
   - Reusable test utilities and harnesses
   - Implementation guides for new intents

3. **Ensure Code Quality**
   - >90% test coverage on all code
   - All linting checks pass
   - No type safety violations
   - Comprehensive error handling

---

## 3. User Stories

### For End Users

**As a user**, I want:
- Consistent navigation that works reliably across all workflows
- Clear feedback when operations succeed or fail
- Ability to go back and return to where I was
- Smooth, responsive UI that doesn't lose my context
- Clear error messages when something goes wrong

### For Developers

**As a developer**, I want:
- Clear architectural patterns I can follow for all intents
- Type-safe communication between intents
- Easy-to-test, isolated intent implementations
- Reliable back navigation that preserves state
- Clear guidelines on what each intent can and cannot do

**As a developer implementing a new intent**, I want:
- A template (CaptureEvent) showing the exact pattern to follow
- Clear state machine design guidelines
- Examples of modal sub-flows and async operations
- Test utilities for isolated testing
- Documentation of the Intent interface contract

**As a QA engineer**, I want:
- Predictable state transitions that are easy to test
- Clear error states and recovery paths
- Ability to test workflows in isolation
- Back navigation that works consistently
- E2E tests that validate complete workflows

### For Architects/Reviewers

**As an architect**, I want:
- Strict adherence to the TUI_INTENT_DIAGRAM specification
- Type-safe communication preventing runtime errors
- Clear separation of concerns
- Extensible patterns for future intents
- Comprehensive documentation of design decisions

---

## 4. Functional Requirements

### 4.1 Core Architecture Requirements

#### Requirement 1: Intent Boundary Contract Types
The system must implement the following types in `internal/cli/intents/contract.go`:

1. **IntentStatus** - Enumeration with four values:
   - `StatusCompleted` - Intent finished successfully
   - `StatusCancelled` - User explicitly cancelled
   - `StatusFailed` - Intent encountered an error
   - `StatusPartial` - Intent succeeded partially

2. **IntentError** - Error type with:
   - `Code` field (string) - Machine-readable error code
   - `Message` field (string) - User-facing error message
   - `Cause` field (error) - Underlying error for logging
   - Error() method for error interface compliance

3. **IntentResult[T]** - Generic result type with:
   - `Status` field (IntentStatus)
   - `Data` field (T)
   - `Error` field (*IntentError)
   - `Metadata` field (map[string]interface{})
   - Helper methods: `WithMetadata()`, `GetMetadata()`, `IsSuccess()`, `IsError()`

4. **Intent** - Interface that all intents must implement:
   - `Init(ctx context.Context) tea.Cmd` - Initialize intent
   - `Update(msg tea.Msg) (Intent, tea.Cmd)` - Handle updates
   - `View() string` - Render the intent
   - `Result() *IntentResult[interface{}]` - Return typed result

5. **ModalEditResult[T]** - Modal sub-flow result type with:
   - `Original` field (T) - Original value before editing
   - `Modified` field (T) - Modified value after editing
   - `Accepted` field (bool) - Whether user confirmed
   - `Changes` field (map[string]interface{}) - Tracked changes

**Acceptance Criteria**:
- All types compile without errors
- All helper methods work correctly
- No runtime type assertions needed
- >95% test coverage on all types
- Types are properly documented

#### Requirement 2: IntentRouter Implementation
The system must implement `IntentRouter` in `internal/cli/intents/router.go`:

1. **Core Functionality**:
   - Register intents by name
   - Activate intents by name with context
   - Handle intent updates and results
   - Navigate back through history
   - Restore previous intent with metadata

2. **Router Interface**:
   - `NewIntentRouter(ctx context.Context) *IntentRouter`
   - `Register(name string, intent Intent) error`
   - `Activate(name string) tea.Cmd`
   - `Update(msg tea.Msg) (tea.Cmd, error)`
   - `View() string`
   - `Back() tea.Cmd`
   - `GetCurrentIntent() Intent`
   - `GetHistory() []Intent`

3. **Behavior**:
   - Thread-safe with proper mutex protection
   - Validates intent exists before activation
   - Pushes current intent to history on activation
   - Calls `Init()` on newly activated intents
   - Delegates all updates to current intent
   - Handles result callbacks when intents complete
   - Restores previous intent on back navigation

**Acceptance Criteria**:
- Router can activate and switch intents correctly
- Back navigation restores previous intent
- History is properly maintained
- Result handler is called when intent completes
- >90% test coverage
- No goroutine leaks or race conditions

#### Requirement 3: Root Model Refactoring
The root Bubble Tea model (`internal/cli/app/app.go`) must:

1. **Replace inline state with IntentRouter**:
   - Remove all inline intent-specific state
   - Add `router *IntentRouter` field
   - Keep only global UI state (theme, preferences)

2. **Delegate to Router**:
   - `Update()` delegates all intent messages to router
   - `View()` calls router's view method
   - `Init()` initializes router and registers intents
   - All intent navigation goes through router

3. **Handle Global Shortcuts**:
   - Quit/exit from any intent
   - Help system accessible from any intent
   - Main menu shortcut (Ctrl+Home)
   - Back navigation (Esc)

4. **Preserve Existing Functionality**:
   - No breaking changes to existing CLI
   - All existing tests pass
   - No regressions in user workflows

**Acceptance Criteria**:
- Root model successfully delegates to router
- All intents activate and return results correctly
- Back navigation works from all intents
- No breaking changes to existing CLI
- >90% test coverage

### 4.2 Intent Implementation Requirements

#### Requirement 4: CaptureEvent Intent (Template Pattern)
Implement `CaptureEvent` intent in `internal/cli/intents/capture/` as the template for all other intents:

1. **State Machine**:
   - `StateChooseStrategy` - Choose capture method
   - `StateCaptureForm` - Capture event details
   - `StateReviewInferred` - Review inferred metadata/bursts/facts
   - `StateSubmit` - Confirm and submit event

2. **Form Validation**:
   - Required field validation
   - Date/time validation
   - Format validation
   - Clear error messages

3. **Domain Integration**:
   - Call enrichment service for inference
   - Call domain service to save event
   - Handle service errors gracefully
   - Return typed result with event data

4. **Modal Sub-Flows**:
   - EditMetadataModal returning `ModalEditResult[Metadata]`
   - EditBurstModal returning `ModalEditResult[Burst]`
   - EditFactModal returning `ModalEditResult[Fact]`
   - Context preservation on cancel

5. **Views**:
   - Strategy selection view
   - Form capture view with validation errors
   - Review/inline editing view
   - Result confirmation view

6. **Result Handling**:
   - Return `IntentResult[CaptureEventResult]`
   - Include event data and timestamp
   - Include source strategy in metadata

**Acceptance Criteria**:
- Implements full Intent interface
- All state transitions work correctly
- Form validation prevents invalid data
- Modal sub-flows preserve context
- Error handling is consistent and user-friendly
- >90% test coverage
- Back navigation restores complete state
- Serves as template for other intents

#### Requirement 5: BrowseTimeline Intent
Implement `BrowseTimeline` intent following the CaptureEvent pattern:

1. **State Machine**:
   - `StateTimelineView` - Display timeline with filters/sort/select
   - `StateEventDetail` - Show event details

2. **Filtering & Sorting**:
   - Apply filters without state mutation
   - Support multiple sort orders
   - Preserve filter/sort state in metadata

3. **Navigation**:
   - Select event to view details
   - Return to timeline with context preserved
   - Back navigation restores scroll position, filters, sort, selection

4. **Large Dataset Handling**:
   - Efficient rendering with 1000+ events
   - Pagination or lazy loading
   - Performance benchmarks <100ms per render

**Acceptance Criteria**:
- Follows CaptureEvent pattern
- Filtering and sorting work correctly
- Back navigation preserves all context
- Performance acceptable with large datasets
- >90% test coverage

#### Requirement 6: GenerateCV Intent
Implement `GenerateCV` intent following the CaptureEvent pattern:

1. **State Machine**:
   - `StateSelectProfile` - Choose profile
   - `StateValidateProfile` - Validate profile completeness
   - `StateSelectAudience` - Choose target audience
   - `StateValidateAudience` - Validate audience compatibility
   - `StateGeneratePreview` - Generate CV preview
   - `StateReviewCV` - Review with inline edit mode
   - `StateConfirmCV` - Confirm generation
   - `StateArtifactReady` - Show result

2. **Validation**:
   - Profile completeness checks
   - Audience compatibility checks
   - Clear error messages

3. **Preview Generation**:
   - Call CV generation service
   - Display formatted preview
   - Handle generation errors

4. **Inline Editing**:
   - Toggle edit mode in review state
   - Edit CV content
   - Validate changes
   - Apply or cancel changes

**Acceptance Criteria**:
- Follows CaptureEvent pattern
- All validations work correctly
- Preview generation is performant
- Inline editing works without breaking state
- >90% test coverage

#### Requirement 7: ExportArtifact Intent
Implement `ExportArtifact` intent with async operations pattern:

1. **State Machine**:
   - `StateSelectArtifactType` - Choose what to export
   - `StateConfigureExport` - Configure export settings
   - `StatePreviewExport` - Show preview
   - `StateConfirmExport` - Confirm export
   - `StateExportInProgress` - Async export operation
   - Handle completion and errors

2. **Async Operations**:
   - Use ephemeral `InProgress` state
   - Non-blocking operation with progress feedback
   - Handle errors with retry logic
   - Support cancellation

3. **Progress Feedback**:
   - Show progress bar or spinner
   - Display progress percentage
   - Allow user to cancel

**Acceptance Criteria**:
- Follows CaptureEvent pattern with async extension
- Async operations don't block UI
- Progress feedback is clear
- Error handling allows retry
- >90% test coverage

#### Requirement 8: ConfigureSystem Intent
Implement `ConfigureSystem` intent with staged changes pattern:

1. **State Machine**:
   - `StateSelectDomain` - Choose config domain
   - `StateEditSettings` - Edit settings
   - `StateStageChanges` - Stage changes for review
   - `StateSaveConfiguration` - Save or discard

2. **Staged Changes**:
   - Track all changes locally
   - No partial writes
   - Discard all on cancel
   - Validate before save

3. **Configuration Domains**:
   - Support multiple configuration domains
   - Each domain has its own settings
   - Changes are domain-scoped

**Acceptance Criteria**:
- Follows CaptureEvent pattern with staged changes
- All changes are local until save
- Validation prevents invalid saves
- Discard works without side effects
- >90% test coverage

### 4.3 Cross-Intent Requirements

#### Requirement 9: Back Navigation with Metadata
All intents must support back navigation with full context restoration:

1. **Metadata Preservation**:
   - Store scroll position in metadata
   - Store filter state in metadata
   - Store sort order in metadata
   - Store current selection in metadata

2. **Automatic Restoration**:
   - Router restores metadata on back navigation
   - Intents restore their state from metadata
   - No manual context passing needed

3. **Type Safety**:
   - All metadata is explicitly typed
   - No arbitrary object passing
   - Clear documentation of metadata keys

**Acceptance Criteria**:
- Back navigation works from all intents
- All context is restored correctly
- No data loss on back navigation
- Metadata is properly typed

#### Requirement 10: Modal Sub-Flows Pattern
All intents with inline editing must use modal sub-flows:

1. **Modal Sub-Flow Contract**:
   - Modal returns `ModalEditResult[T]`
   - Original value preserved on cancel
   - Changes tracked in Changes map
   - Accepted flag indicates user confirmation

2. **Context Preservation**:
   - Parent intent state unchanged on cancel
   - No side effects if user cancels
   - Clear undo/cancel behavior

3. **Consistent Handling**:
   - All modals follow same pattern
   - All modals return typed results
   - All modals preserve context

**Acceptance Criteria**:
- All modals follow pattern consistently
- Context preserved on cancel
- No global state mutation
- >90% test coverage

#### Requirement 11: GlobalContext Pattern (Enhancement)
Implement lightweight global context for non-mutable shared metadata:

1. **GlobalContext Structure**:
   - Read-only user preferences
   - Transient UI state (last selected items)
   - Read-only application config
   - Thread-safe access with mutex

2. **Usage**:
   - Pass via context to all intents
   - Intents can read but not write
   - Session-only, no persistence

3. **Benefits**:
   - Reduces metadata boilerplate
   - Maintains intent independence
   - Improves UX with remembered preferences

**Acceptance Criteria**:
- GlobalContext properly defined and thread-safe
- All intents can access GlobalContext
- No persistence from GlobalContext
- >90% test coverage

#### Requirement 12: Async Feedback Pattern (Enhancement)
Implement progress feedback for long-running operations:

1. **Progress Indicators**:
   - Progress bar for operations with known duration
   - Spinner for indeterminate operations
   - Progress percentage display
   - Cancellation support

2. **Async Operation Flow**:
   - InProgress ephemeral state
   - Non-blocking operation execution
   - Progress updates via messages
   - Completion trigger with typed result

**Acceptance Criteria**:
- Progress feedback is clear and responsive
- Operations don't block UI
- Cancellation works correctly
- >90% test coverage

#### Requirement 13: CI/CD Integration (Enhancement)
Set up automated quality checks:

1. **Automated Tests**:
   - Run all tests on every PR
   - Enforce >90% coverage requirement
   - Run property-based tests
   - Check for race conditions

2. **Code Quality**:
   - Linting with golangci-lint
   - Type checking with go vet
   - Format checking with gofmt
   - No unused imports

3. **Performance**:
   - Baseline rendering benchmarks
   - Form validation benchmarks
   - CV generation benchmarks
   - Track performance over time

**Acceptance Criteria**:
- CI/CD pipeline configured and working
- All checks pass for all code
- Coverage reports available
- Performance benchmarks tracked

#### Requirement 14: Performance Benchmarking (Enhancement)
Establish performance baselines and optimize:

1. **Benchmark Tests**:
   - Rendering performance (<100ms per frame)
   - Form validation (<1ms)
   - List filtering (<50ms)
   - CV preview generation (<500ms)

2. **Performance Optimization**:
   - Profile rendering bottlenecks
   - Optimize hot paths
   - Add caching where appropriate
   - Test with large datasets

**Acceptance Criteria**:
- All benchmarks meet targets
- Performance tracked over time
- Large datasets handled efficiently
- No memory leaks

### 4.4 Testing Requirements

#### Requirement 15: Unit Tests
All code must have comprehensive unit tests:

1. **Test Coverage**:
   - >90% code coverage required
   - All state transitions tested
   - All views tested
   - All validations tested
   - All error cases tested

2. **Test Organization**:
   - Use Ginkgo for test organization
   - Use Gomega for assertions
   - Create test utilities and helpers
   - Mock external dependencies

3. **Test Utilities**:
   - `IntentTestHarness` for isolated testing
   - `IntentRouterTestHelper` for router testing
   - Mock domain services
   - Test data generators

**Acceptance Criteria**:
- >90% test coverage on all code
- All tests pass
- Tests are isolated and repeatable
- Mocking works correctly

#### Requirement 16: Integration Tests
Test complete workflows from entry to exit:

1. **Workflow Tests**:
   - Complete intent workflows
   - Intent switching and navigation
   - Back navigation with state restoration
   - Result propagation

2. **Error Handling Tests**:
   - Error recovery paths
   - Validation error handling
   - Service error handling
   - Timeout handling

3. **Edge Cases**:
   - Empty data sets
   - Large data sets
   - Rapid user input
   - Network failures

**Acceptance Criteria**:
- All workflows tested end-to-end
- Error paths tested
- Edge cases covered
- >90% test coverage

#### Requirement 17: E2E Tests
Test complete user journeys:

1. **User Journey Tests**:
   - Complete capture workflow
   - Complete CV generation workflow
   - Complete export workflow
   - Complete configuration workflow

2. **Navigation Tests**:
   - Forward navigation works
   - Back navigation works
   - Context preservation works
   - Menu navigation works

3. **Error Scenarios**:
   - User cancellation
   - Service errors
   - Validation errors
   - Recovery paths

**Acceptance Criteria**:
- All major user journeys tested
- Navigation tested thoroughly
- Error scenarios covered
- Tests are maintainable

---

## 5. Non-Goals (Out of Scope)

### Explicitly Out of Scope

1. **Phase 5 Secondary Intents** - Skill Tracking, Career Goals, Mentor Matching, Continuous Learning are post-release (not included in this refactoring)

2. **Data Migration** - No migration of existing data required; refactoring is architectural only

3. **UI/UX Redesign** - No visual redesign; keep existing visual style and components

4. **Performance Optimization Beyond Baselines** - Set baselines; detailed optimization is future work

5. **Advanced Analytics** - Tracking and analytics are out of scope

6. **Mobile-Specific Optimization** - Desktop-first focus

7. **Localization/Internationalization** - Not included in this refactoring

8. **Custom Theming System** - Keep existing color scheme and styling

### Intentionally Excluded

- Changes to domain logic or data models
- Changes to database schema or persistence layer
- Changes to CLI command structure
- Changes to CSV import functionality
- Changes to existing test infrastructure beyond TUI tests

---

## 6. Design Considerations

### Architectural Constraints

**Strict Adherence to TUI_INTENT_DIAGRAM.md**:
- All intents MUST implement the Intent interface exactly as specified
- All intents MUST return `IntentResult[T]` for communication
- All intents MUST own only local state
- All intents MUST NOT mutate global state
- All navigation MUST go through IntentRouter
- All back navigation MUST restore metadata

**Type Safety Requirements**:
- No runtime type assertions
- No unsafe casts
- All communication must be strongly typed
- Illegal states must be unrepresentable

**State Machine Requirements**:
- All states must be explicit constants
- All transitions must be validated
- No implicit behavior
- Clear entry and exit points

### Component Refactoring Strategy

1. **Keep Existing Components**:
   - Reuse existing UI components (forms, lists, modals, etc.)
   - Refactor to fit intent pattern where needed
   - Maintain visual consistency

2. **Refactor Models**:
   - Move from monolithic models to intent-based models
   - Each intent has its own model file
   - Clear separation of concerns

3. **Update Navigation**:
   - Replace direct model-to-model navigation with IntentRouter
   - Update message types to work with router
   - Implement proper result handling

### Code Organization

**File Structure**:
```
internal/cli/
  /intents
    contract.go           # Intent interface and result types
    router.go             # IntentRouter implementation
    /capture
      model.go            # CaptureEvent model and states
      update.go           # State transition logic
      view.go             # Rendering logic
      /modals             # Modal sub-flows
        edit_metadata.go
        edit_burst.go
        edit_fact.go
      capture_test.go
    /browse
      model.go
      update.go
      view.go
      browse_test.go
    /generate_cv
      model.go
      update.go
      view.go
      generate_cv_test.go
    /export
      model.go
      update.go
      view.go
      export_test.go
    /configure
      model.go
      update.go
      view.go
      configure_test.go
  /app
    app.go                # Root model
    update.go
    view.go
  /components
    # Existing components, refactored as needed
  /models
    # Existing models, refactored for intents
```

### Naming Conventions

**Intent Names**:
- `CaptureEventIntent` (not `CaptureScreen` or `CaptureModel`)
- `BrowseTimelineIntent`
- `GenerateCVIntent`
- `ExportArtifactIntent`
- `ConfigureSystemIntent`

**State Constants**:
- `StateChooseCaptureStrategy` (descriptive, workflow-specific)
- `StateReviewInferredEvent` (not `State1` or `ReviewState`)
- `StateGeneratePreview` (not `GeneratingCV`)

**Component Names**:
- `CaptureEventForm` (reflects workflow position)
- `CVAudienceSelection`
- `ExportDestinationConfirm`

**Message Types**:
- `ChooseCaptureStrategyMsg`
- `CaptureFormInputMsg`
- `ReviewInferredEventMsg`

---

## 7. Technical Considerations

### Dependencies

**Required**:
- Go 1.21+
- Bubble Tea (TUI framework)
- Gomega/Ginkgo (testing)
- gopter (property-based testing)

**Existing**:
- Domain services for enrichment, CV generation, etc.
- Database layer for persistence
- CLI infrastructure

### Integration Points

1. **Domain Services**:
   - Enrichment service for burst/fact inference
   - CV generation service
   - Configuration management service
   - All integrated via intent handlers

2. **Existing Models**:
   - Reuse existing components where possible
   - Refactor models to fit intent pattern
   - Maintain data flow to domain layer

3. **Database Layer**:
   - No changes to persistence layer
   - Intents call domain services which handle persistence
   - All data flow through service layer

### Known Technical Constraints

1. **Bubble Tea Limitations**:
   - Single-threaded message processing
   - Commands must be non-blocking
   - Long operations need async pattern

2. **State Machine Complexity**:
   - Some workflows have many states
   - Must balance clarity with complexity
   - Documentation is critical

3. **Testing Complexity**:
   - Mocking Bubble Tea messages is tricky
   - Need good test utilities
   - Property-based tests help verify invariants

### Suggested Approaches

1. **For Async Operations**:
   - Use ephemeral `InProgress` state
   - Spawn goroutine in command
   - Send completion message when done
   - Handle errors with retry logic

2. **For Modal Sub-Flows**:
   - Create modal as separate component
   - Return `ModalEditResult[T]` on completion
   - Parent handles result in state transition
   - Context preserved automatically

3. **For Back Navigation**:
   - Store metadata in `IntentResult`
   - Router restores on back navigation
   - Intent restores state from metadata
   - No manual context passing

4. **For Global State**:
   - Minimal global state only (theme, preferences)
   - Use GlobalContext pattern for shared reads
   - Never write to global state from intents
   - Pass context explicitly where needed

---

## 8. Success Metrics

### Functional Success

1. **Workflow Consistency** (Primary)
   - All intents follow the same architectural pattern
   - No broken workflows or navigation issues
   - Back navigation works reliably
   - Error recovery is consistent

2. **Type Safety** (Primary)
   - No runtime type assertions in intent code
   - All communication via `IntentResult[T]`
   - Illegal states unrepresentable
   - Type checker catches errors at compile time

3. **State Machine Correctness** (Primary)
   - All state transitions are explicit
   - No implicit behavior
   - All transitions validated
   - Clear entry and exit points

### Code Quality Success

1. **Test Coverage** (Required)
   - >90% test coverage on all code
   - All intents have unit tests
   - All workflows have integration tests
   - All user journeys have E2E tests

2. **Code Standards** (Required)
   - All linting checks pass
   - All code properly formatted
   - No type safety violations
   - Clear, documented code

3. **Performance** (Required)
   - Rendering <100ms per frame
   - Form validation <1ms
   - List filtering <50ms
   - CV preview generation <500ms

### Developer Experience Success

1. **Documentation** (Required)
   - Clear implementation guide
   - Code examples for all patterns
   - Testing guide
   - Troubleshooting guide

2. **Extensibility** (Required)
   - Clear pattern for new intents
   - Template (CaptureEvent) for reference
   - Test utilities available
   - Documentation complete

3. **Maintainability** (Required)
   - Code is easy to understand
   - State machines are clear
   - Naming conventions consistent
   - No hidden dependencies

### User Experience Success

1. **Workflow Reliability**
   - No broken workflows
   - Consistent navigation
   - Back navigation works
   - Error recovery is clear

2. **Performance**
   - No UI lag or freezing
   - Responsive to user input
   - Fast feedback for operations
   - Smooth transitions

3. **Error Handling**
   - Clear error messages
   - Suggested recovery actions
   - No silent failures
   - Helpful logging

---

## 9. Implementation Phases

### Phase 1: Foundation & Core Infrastructure (1.5 weeks)

**Deliverables**:
- Intent boundary contract types
- IntentRouter implementation
- Root model refactoring
- All Phase 1 tests pass with >90% coverage

**Validation Gate**: All Phase 1 acceptance criteria met

### Phase 2: CaptureEvent Intent Template (2 weeks)

**Deliverables**:
- Complete CaptureEvent intent
- All state transitions
- Modal sub-flows
- Comprehensive tests (>90% coverage)

**Validation Gate**: CaptureEvent serves as template for other intents

### Phase 3: Remaining Core Intents (4 weeks)

**Deliverables**:
- BrowseTimeline intent
- GenerateCV intent
- ExportArtifact intent
- ConfigureSystem intent
- All tests pass with >90% coverage

**Validation Gate**: All intents follow CaptureEvent pattern

### Phase 4: Integration & Polish (2 weeks)

**Deliverables**:
- All intents integrated
- Global shortcuts working
- Comprehensive logging
- Performance optimized
- Documentation complete

**Validation Gate**: Ready for production release

### Phase 5: Enhancements (Concurrent with Phase 4)

**Deliverables**:
- GlobalContext implementation
- Async feedback pattern
- CI/CD pipeline
- Performance benchmarks

**Validation Gate**: All enhancements working and tested

---

## 10. Open Questions

### Architecture Questions

1. **Existing Component Reuse**: Should we keep all existing components in `internal/cli/components/` and `internal/cli/models/`, or refactor/remove ones that don't fit the intent pattern?

2. **Message Types**: Should message types be defined per-intent (e.g., `capture/messages.go`) or globally?

3. **Service Layer**: Should domain services be passed to intents via context or dependency injection?

4. **Logging**: Should logging be handled by each intent or centrally by the router?

### Implementation Questions

1. **Async Operations**: For ExportArtifact, should we use goroutines or channels for async work?

2. **Modal Lifecycle**: Should modals be separate intent implementations or helper components within intents?

3. **Error Recovery**: Should failed intents return to previous state or stay in error state for retry?

4. **Metadata Storage**: Should metadata be stored in `IntentResult` or separately in router?

### Testing Questions

1. **Property-Based Tests**: Which invariants should we focus on for property-based testing?

2. **Mock Services**: Should we create a mock implementation of all domain services?

3. **E2E Test Framework**: Should we use existing test infrastructure or create new E2E test utilities?

### Performance Questions

1. **Large Datasets**: What's the expected maximum number of events in timeline view?

2. **Rendering Strategy**: Should we implement pagination or lazy loading for large lists?

3. **Caching**: What data should be cached to improve performance?

---

## 11. Success Criteria Checklist

### Phase 1 Acceptance

- [ ] All intent boundary contract types implemented
- [ ] IntentRouter fully functional
- [ ] Root model successfully refactored
- [ ] All Phase 1 tests pass with >90% coverage
- [ ] No breaking changes to existing CLI
- [ ] Code review approved

### Phase 2 Acceptance

- [ ] CaptureEvent intent fully implemented
- [ ] All state transitions working
- [ ] Modal sub-flows functional
- [ ] All Phase 2 tests pass with >90% coverage
- [ ] Serves as template for other intents
- [ ] Code review approved

### Phase 3 Acceptance

- [ ] All four core intents implemented
- [ ] All intents follow CaptureEvent pattern
- [ ] All intents have >90% test coverage
- [ ] Navigation working between all intents
- [ ] Back navigation works for all intents
- [ ] Code review approved

### Phase 4 Acceptance

- [ ] All intents integrated
- [ ] Global shortcuts working
- [ ] Logging comprehensive
- [ ] Performance acceptable
- [ ] Documentation complete
- [ ] Code review approved
- [ ] **Ready for production release**

### Enhancement Acceptance

- [ ] GlobalContext implemented and tested
- [ ] Async feedback pattern working
- [ ] CI/CD pipeline configured
- [ ] Performance benchmarks established
- [ ] All enhancements tested and documented

---

## 12. References & Resources

### Architecture Documentation
- `docs/TUI_INTENT_DIAGRAM.md` - Complete architectural specification (MUST READ)
- `docs/IMPLEMENTATION_ROADMAP.md` - Detailed phase-by-phase plan
- `docs/IMPLEMENTATION_CHECKLIST.md` - Detailed developer checklist
- `docs/IMPLEMENTATION_ENHANCEMENTS.md` - Enhancement recommendations

### Guidelines & Standards
- `docs/TUI_STANDARDS.md` - UI/UX standards
- `docs/TUI_DEVELOPER_GUIDE.md` - General development guidelines
- `docs/KEYBOARD_REFERENCE.md` - Keyboard shortcuts
- `AGENTS.md` - Agent guidelines

### Workflow Documentation
- `docs/WORKFLOW_DIAGRAM.md` - High-level workflows
- `docs/PRD_MASTER.md` - Product requirements
- `docs/PRD_USER_STORIES.md` - User stories

---

## 13. Approval & Sign-Off

### Document Approval

- **Author**: AI Assistant
- **Date**: 2026-01-02
- **Status**: Ready for Implementation
- **Confidence Level**: Very High

### Required Approvals Before Implementation

- [ ] Architecture Lead Review
- [ ] Development Team Review
- [ ] Product Owner Sign-Off

### Implementation Readiness

- [x] Architecture fully specified
- [x] Requirements clearly defined
- [x] Acceptance criteria established
- [x] Test strategy defined
- [x] Implementation roadmap created
- [x] Documentation prepared

**Status**: ✅ **READY FOR IMMEDIATE IMPLEMENTATION**

---

## 14. Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-02 | AI Assistant | Initial PRD based on TUI_INTENT_DIAGRAM.md |

---

*This PRD is based strictly on the specifications in `docs/TUI_INTENT_DIAGRAM.md` and supporting documentation. All requirements must be implemented exactly as specified to achieve the architectural goals.*

