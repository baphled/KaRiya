# Task 14: Migrate CaptureEvent and BrowseTimeline Intents to StandardView

## Overview
- **Goal**: Migrate CaptureEventIntent and BrowseTimelineIntent to use StandardView with logo, modals, and standardized layout
- **Time Estimate**: 5 hours
- **Prerequisites**: Tasks 12 and 13 completed (StandardView components and infrastructure ready)
- **Status**: Not Started

## Motivation
Migrate the two most complex and frequently used intents to the new standardized view system. These intents serve as the primary patterns that other intents will follow.

## Files to Modify
- [ ] `internal/cli/intents/capture_event.go` (data structures)
- [ ] `internal/cli/intents/capture_event_intent.go` (implementation)
- [ ] `internal/cli/intents/capture_event_test.go` (tests)
- [ ] `internal/cli/intents/browse_timeline.go` (data structures)
- [ ] `internal/cli/intents/browse_timeline_intent.go` (implementation)
- [ ] `internal/cli/intents/browse_timeline_test.go` (tests)

## Implementation Checklist

### Phase 1: Preparation (15 min)
- [ ] Create task file
- [ ] Run compliance check (baseline): `make check-compliance`
- [ ] Verify Tasks 12 and 13 components exist and tests pass
- [ ] Review current CaptureEventIntent implementation
- [ ] Review current BrowseTimelineIntent implementation
- [ ] Run existing tests: `go test ./internal/cli/intents/...`
- [ ] Document current test coverage baseline

### Phase 2: Migrate CaptureEventIntent (150 min)

#### 2.1 Update CaptureEventIntent structure (20 min)
- [ ] Open `internal/cli/intents/capture_event.go`
- [ ] Verify BaseIntent is embedded
- [ ] Add loading state fields if not inherited:
  - Verify `isLoading`, `loadingMessage` from BaseIntent
  - Verify `errorState` from BaseIntent
  - Verify `successMessage`, `successTime` from BaseIntent
- [ ] Add loading message rotator:
  - `loadingRotator *components.LoadingMessageRotator`
- [ ] Initialize in constructor:
  - Call `InitializeLogo()` in NewCaptureEventIntent
  - Create loading rotator with capture-specific messages
  - Set logo spacing to 2
- [ ] Commit: `feat(capture): update CaptureEventIntent for StandardView`

#### 2.2 Refactor View method (30 min)
- [ ] Open `internal/cli/intents/capture_event_intent.go`
- [ ] Locate main `View() string` method (around line 609)
- [ ] Refactor to use StandardView pattern:
```go
func (i *CaptureEventIntent) View() string {
    if !i.active {
        return "CaptureEvent intent is not active"
    }
    
    // Create standard view with breadcrumbs
    view := i.CreateViewWithBreadcrumbs("Main Menu", "Capture Event", i.getStateName())
    
    // Handle modal states
    if i.isLoading {
        ShowLoadingModal(view, i.loadingRotator.GetCurrent(), true)
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    if i.ShouldShowSuccess() {
        ShowSuccessModal(view, i.successMessage)
    }
    
    // Get content for current state
    content := i.getStateContent()
    view.WithContent(content)
    
    // Get context-aware help
    help := i.getContextHelp()
    view.WithHelp(help).WithFooterSeparator(true)
    
    return view.Render()
}
```
- [ ] Commit: `refactor(capture): update View method to use StandardView`

#### 2.3 Create getStateName helper (10 min)
- [ ] Add `getStateName() string` method:
```go
func (i *CaptureEventIntent) getStateName() string {
    switch i.state.currentState {
    case CaptureStateChooseStrategy:
        return "Choose Strategy"
    case CaptureStateForm:
        return "Enter Details"
    case CaptureStateReview:
        return "Review"
    case CaptureStateSubmit:
        return "Submit"
    default:
        return string(i.state.currentState)
    }
}
```
- [ ] Commit: `feat(capture): add state name helper`

#### 2.4 Create getStateContent method (30 min)
- [ ] Add `getStateContent() string` method that delegates to state-specific renderers
- [ ] Rename existing view methods to render methods:
  - `viewChooseStrategy()` → `renderChooseStrategy()`
  - `viewCaptureForm()` → `renderCaptureForm()`
  - `viewReviewInferredEvent()` → `renderReviewInferredEvent()`
  - `viewSubmit()` → `renderSubmit()`
  - `viewError()` → `renderError()` (deprecated - now use modal)
- [ ] Update render methods to return content only (no logo, no footer):
  - Remove header text
  - Focus on core content
  - Use full width
  - Return clean content string
- [ ] Implement getStateContent():
```go
func (i *CaptureEventIntent) getStateContent() string {
    // Show error content if there's an error and no modal
    if i.state.error != nil && i.errorState == nil {
        return i.renderError()
    }
    
    switch i.state.currentState {
    case CaptureStateChooseStrategy:
        return i.renderChooseStrategy()
    case CaptureStateForm:
        return i.renderCaptureForm()
    case CaptureStateReview:
        return i.renderReviewInferredEvent()
    case CaptureStateSubmit:
        return i.renderSubmit()
    default:
        return fmt.Sprintf("Unknown state: %s", i.state.currentState)
    }
}
```
- [ ] Commit: `refactor(capture): extract state content rendering`

#### 2.5 Create getContextHelp method (20 min)
- [ ] Add `getContextHelp() string` method:
```go
func (i *CaptureEventIntent) getContextHelp() string {
    base := "q Quit  m Main Menu"
    
    switch i.state.currentState {
    case CaptureStateChooseStrategy:
        return CombineFooters("1-3 Select  Esc Cancel", base)
    case CaptureStateForm:
        return CombineFooters(FormFooter(), base)
    case CaptureStateReview:
        return CombineFooters("Enter Confirm  e Edit  Esc Back", base)
    case CaptureStateSubmit:
        return base
    default:
        return base
    }
}
```
- [ ] Commit: `feat(capture): add context-aware help text`

#### 2.6 Update loading states (20 min)
- [ ] In form submission logic, add loading state:
```go
i.SetLoading("Saving event...")
i.loadingRotator.SetMessages([]string{
    "⏳ Validating input...",
    "💾 Saving to database...",
    "📊 Creating indexes...",
    "✨ Finalizing...",
})
```
- [ ] After successful save:
```go
i.ClearLoading()
i.SetSuccess("✅ Event captured successfully!")
```
- [ ] On error:
```go
i.ClearLoading()
i.SetError(err)
```
- [ ] Commit: `feat(capture): add loading and success states`

#### 2.7 Update render methods for full width (20 min)
- [ ] Update `renderChooseStrategy()`:
  - Remove manual box drawing
  - Use lipgloss for styling
  - Return content that uses full width
- [ ] Update `renderCaptureForm()`:
  - Use form model's view directly
  - Apply full-width styling
- [ ] Update `renderReviewInferredEvent()`:
  - Use lipgloss card styles
  - Full width content
- [ ] Update `renderSubmit()`:
  - Simple centered message
  - Progress indicator if needed
- [ ] Commit: `refactor(capture): update render methods for full width`

### Phase 3: Migrate BrowseTimelineIntent (120 min)

#### 3.1 Update BrowseTimelineIntent structure (15 min)
- [ ] Open `internal/cli/intents/browse_timeline.go`
- [ ] Verify BaseIntent is embedded
- [ ] Add loading rotator field
- [ ] Initialize in constructor:
  - Call `InitializeLogo()`
  - Create loading rotator with timeline-specific messages
- [ ] Commit: `feat(browse): update BrowseTimelineIntent for StandardView`

#### 3.2 Refactor View method (25 min)
- [ ] Open `internal/cli/intents/browse_timeline_intent.go`
- [ ] Locate main `View() string` method (around line 364)
- [ ] Refactor to StandardView pattern:
```go
func (i *BrowseTimelineIntent) View() string {
    if !i.active {
        return "BrowseTimeline intent is not active"
    }
    
    view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
    
    if i.isLoading {
        ShowLoadingModal(view, i.loadingRotator.GetCurrent(), false)
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    
    content := i.getStateContent()
    view.WithContent(content)
    
    help := i.getContextHelp()
    view.WithHelp(help).WithFooterSeparator(true)
    
    return view.Render()
}
```
- [ ] Commit: `refactor(browse): update View method to use StandardView`

#### 3.3 Create state helpers (15 min)
- [ ] Add `getStateName() string` method:
```go
func (i *BrowseTimelineIntent) getStateName() string {
    switch i.state {
    case BrowseStateTimeline:
        return "Timeline"
    case BrowseStateEventDetail:
        return "Event Detail"
    default:
        return string(i.state)
    }
}
```
- [ ] Add `getStateContent() string` method:
```go
func (i *BrowseTimelineIntent) getStateContent() string {
    switch i.state {
    case BrowseStateTimeline:
        return i.renderTimeline()
    case BrowseStateEventDetail:
        return i.renderEventDetail()
    default:
        return "Unknown state"
    }
}
```
- [ ] Commit: `feat(browse): add state helper methods`

#### 3.4 Create getContextHelp method (15 min)
- [ ] Add `getContextHelp() string` method:
```go
func (i *BrowseTimelineIntent) getContextHelp() string {
    base := "q Quit  m Main Menu"
    
    switch i.state {
    case BrowseStateTimeline:
        if len(i.events) == 0 {
            return CombineFooters("Esc Back", base)
        }
        return CombineFooters(NavigationFooter(), "Enter View Details", base)
    case BrowseStateEventDetail:
        return CombineFooters("Esc Back to Timeline", base)
    default:
        return base
    }
}
```
- [ ] Commit: `feat(browse): add context-aware help text`

#### 3.5 Update render methods (30 min)
- [ ] Rename and refactor `viewTimeline()` → `renderTimeline()`:
  - Remove header/footer elements
  - Focus on timeline list content
  - Use full width for table/list
  - Handle empty state gracefully
- [ ] Rename and refactor `viewEventDetail()` → `renderEventDetail()`:
  - Remove header/footer elements
  - Use lipgloss for detail card
  - Full width content
  - Show metadata, bursts, facts
- [ ] Add empty state message in renderTimeline():
```go
if len(i.events) == 0 {
    return lipgloss.NewStyle().
        Align(lipgloss.Center).
        Render("No events found. Create your first event!")
}
```
- [ ] Commit: `refactor(browse): update render methods for full width`

#### 3.6 Add loading states (20 min)
- [ ] If events are loaded asynchronously, add loading:
```go
i.SetLoading("Loading events...")
i.loadingRotator.SetMessages([]string{
    "🔍 Loading events...",
    "📅 Organizing timeline...",
    "✨ Preparing display...",
})
```
- [ ] After loading:
```go
i.ClearLoading()
```
- [ ] On error:
```go
i.ClearLoading()
i.SetError(err)
```
- [ ] Commit: `feat(browse): add loading states for event fetching`

### Phase 4: Update Tests (90 min)

#### 4.1 Update CaptureEventIntent tests (45 min)
- [ ] Open `internal/cli/intents/capture_event_test.go`
- [ ] Update tests that check view output:
  - Look for logo presence in rendered view
  - Look for breadcrumbs in rendered view
  - Look for footer separator
  - Update assertions for new structure
- [ ] Add new tests:
  - Test getStateName() returns correct names
  - Test getContextHelp() returns appropriate shortcuts
  - Test modal states (loading, error, success)
  - Test loading message rotation
- [ ] Update integration tests:
  - Verify StandardView is used
  - Verify terminal info is respected
  - Verify full width rendering
- [ ] Run tests: `go test -v ./internal/cli/intents/ -run TestCaptureEvent`
- [ ] Fix any failing tests
- [ ] Commit: `test(capture): update tests for StandardView migration`

#### 4.2 Update BrowseTimelineIntent tests (45 min)
- [ ] Open `internal/cli/intents/browse_timeline_test.go`
- [ ] Update tests that check view output:
  - Look for logo presence
  - Look for breadcrumbs
  - Look for footer separator
  - Update assertions for new structure
- [ ] Add new tests:
  - Test getStateName() returns correct names
  - Test getContextHelp() returns appropriate shortcuts
  - Test empty state rendering
  - Test error modal display
- [ ] Update integration tests:
  - Verify StandardView usage
  - Verify terminal info handling
- [ ] Run tests: `go test -v ./internal/cli/intents/ -run TestBrowseTimeline`
- [ ] Fix any failing tests
- [ ] Commit: `test(browse): update tests for StandardView migration`

### Phase 5: Integration Testing (45 min)

#### 5.1 Run full test suite (15 min)
- [ ] Run all intent tests: `go test -v ./internal/cli/intents/...`
- [ ] Run with race detector: `go test -race ./internal/cli/intents/...`
- [ ] Run with coverage: `go test -cover ./internal/cli/intents/...`
- [ ] Verify coverage maintained or improved
- [ ] Fix any failing tests
- [ ] Commit if fixes needed: `fix(intents): resolve integration test failures`

#### 5.2 Manual testing - CaptureEvent (15 min)
- [ ] Build: `go build -o kariya ./cmd/kariya`
- [ ] Run: `./kariya`
- [ ] Navigate to Capture Event (c key)
- [ ] Verify logo appears with spacing at top
- [ ] Verify breadcrumbs show "Main Menu > Capture Event > Choose Strategy"
- [ ] Verify footer separator is visible
- [ ] Verify context help shows appropriate shortcuts
- [ ] Test all states:
  - Choose Strategy: Select option, verify view
  - Form Entry: Enter data, verify full width
  - Review: Verify preview, edit capabilities
  - Submit: Verify success modal appears and auto-dismisses
- [ ] Test error handling: Trigger validation error, verify error modal
- [ ] Test cancellation: Press Esc, verify returns to menu

#### 5.3 Manual testing - BrowseTimeline (15 min)
- [ ] From main menu, select Browse Timeline (b key)
- [ ] Verify logo appears with spacing
- [ ] Verify breadcrumbs show "Main Menu > Browse Timeline > Timeline"
- [ ] Verify footer separator and context help
- [ ] Test timeline view:
  - Navigate with arrow keys/vim keys
  - Select event with Enter
  - Verify full width table/list
- [ ] Test event detail view:
  - Verify breadcrumbs update to "Event Detail"
  - Verify full width detail display
  - Verify all metadata, bursts, facts visible
  - Press Esc to return to timeline
- [ ] Test empty state: Clear all events, verify message displays
- [ ] Test different terminal sizes: Resize terminal, verify responsive

### Phase 6: Documentation (15 min)
- [ ] Add comments to new methods explaining StandardView usage
- [ ] Update any relevant documentation files
- [ ] Add migration notes for other developers:
  - Document the pattern used
  - Show before/after examples
  - Explain modal state management
- [ ] Commit: `docs(intents): document StandardView migration pattern`

### Phase 7: Final Verification (15 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./internal/cli/intents/...`
- [ ] Format code: `go fmt ./internal/cli/intents/...`
- [ ] Review all commits for atomicity
- [ ] Verify commit messages follow conventional format
- [ ] Run full test suite: `go test ./...`
- [ ] Verify no regressions in other intents
- [ ] Test on different terminal sizes (80x24, 120x40, 200x60)

## Testing Instructions

### Automated Tests
```bash
# Run intent tests
go test -v ./internal/cli/intents/

# Run specific intent tests
go test -v ./internal/cli/intents/ -run TestCaptureEvent
go test -v ./internal/cli/intents/ -run TestBrowseTimeline

# Run with coverage
go test -cover ./internal/cli/intents/

# Run with race detector
go test -race ./internal/cli/intents/
```

### Manual Testing Checklist

#### CaptureEvent Intent
- [ ] Logo visible with 2-line spacing from top
- [ ] Breadcrumbs show navigation path
- [ ] Footer separator visible as horizontal line
- [ ] Context help shows relevant shortcuts for each state
- [ ] Choose Strategy: 3 options visible, selection works
- [ ] Form Entry: Full width form, tab navigation works
- [ ] Review: Event preview shown, edit button works
- [ ] Submit: Success modal appears and auto-dismisses after 3s
- [ ] Error handling: Error modal shows with bell, Esc dismisses
- [ ] Loading: Loading modal shows during save, can cancel

#### BrowseTimeline Intent
- [ ] Logo visible with 2-line spacing
- [ ] Breadcrumbs show "Main Menu > Browse Timeline > Timeline"
- [ ] Footer separator visible
- [ ] Timeline: Events listed, navigation works
- [ ] Event Detail: Full event shown, Esc returns to timeline
- [ ] Empty State: Message shown when no events
- [ ] Error modal: Shows errors properly
- [ ] Responsive: Works on different terminal sizes

## Acceptance Criteria
- [ ] CaptureEventIntent uses StandardView throughout
- [ ] CaptureEventIntent shows logo on all screens
- [ ] CaptureEventIntent breadcrumbs update per state
- [ ] CaptureEventIntent context help is state-aware
- [ ] CaptureEventIntent success modal auto-dismisses
- [ ] CaptureEventIntent error modal shows with bell
- [ ] CaptureEventIntent loading modal shows during operations
- [ ] BrowseTimelineIntent uses StandardView throughout
- [ ] BrowseTimelineIntent shows logo on all screens
- [ ] BrowseTimelineIntent breadcrumbs update per state
- [ ] BrowseTimelineIntent context help is state-aware
- [ ] BrowseTimelineIntent handles empty state gracefully
- [ ] Both intents use full terminal width
- [ ] Both intents show footer separator
- [ ] All existing tests pass (100%)
- [ ] New tests added for StandardView features
- [ ] Code coverage maintained or improved
- [ ] Code passes linting and formatting
- [ ] Compliance check passes
- [ ] No regressions in other intents
- [ ] Manual testing successful on different terminal sizes

## Rollback Plan
If issues are discovered:
1. Identify problematic commit(s)
2. Run: `git revert <commit-hash>`
3. Alternative: `git reset --hard <previous-working-commit>`
4. Re-run tests to verify working state
5. Review and fix issues before re-implementing

## Notes
- These are the most complex intents, so take time to get them right
- The patterns established here will be used for other intents
- Pay special attention to state management and modal timing
- Ensure loading messages rotate smoothly
- Success modal auto-dismiss should feel natural (3s is good)
- Error modals should be dismissible but not auto-dismiss
- Terminal resize should trigger re-render with updated dimensions

## Dependencies
- Requires Task 12 (StandardView components)
- Requires Task 13 (Infrastructure updates)
- Requires existing intent implementations
- Uses existing form models and components
