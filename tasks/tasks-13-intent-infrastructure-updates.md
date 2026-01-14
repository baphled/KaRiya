---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 13: Intent Infrastructure Updates for Standardized Views

## Overview
- **Goal**: Update BaseIntent and router infrastructure to support terminal awareness and StandardView creation across all intents
- **Time Estimate**: 2-3 hours
- **Prerequisites**: Task 12 completed (StandardView, Modal, LoadingMessages components created)
- **Status**: Not Started

## Motivation
Enable all intents to easily create standardized views with logo, modals, and context-aware help. This requires enhancing BaseIntent with terminal info tracking, logo management, and StandardView creation helpers.

## Files to Modify
- [ ] `internal/cli/intents/contract.go` (enhance BaseIntent)
- [ ] `internal/cli/intents/router.go` (ensure terminal info propagation)
- [ ] `internal/cli/app/app.go` (update to pass terminal info)

## Files to Create
- [ ] `internal/cli/intents/view_helpers.go` (helper functions for view creation)
- [ ] `internal/cli/intents/view_helpers_test.go` (tests)

## Implementation Checklist

### Phase 1: Preparation (10 min)
- [ ] Create task file
- [ ] Run compliance check (baseline): `make check-compliance`
- [ ] Verify Task 12 components exist and tests pass
- [ ] Review current BaseIntent implementation
- [ ] Review TerminalAwareIntent interface
- [ ] Run existing tests: `go test ./internal/cli/intents/...`

### Phase 2: Enhance BaseIntent (45 min)

#### 2.1 Add terminal and logo fields (15 min)
- [ ] Open `internal/cli/intents/contract.go`
- [ ] Add imports for `components` package
- [ ] Add to `BaseIntent` struct:
  - `terminalInfo *terminal.Info`
  - `logo *components.ASCIILogo`
  - `logoSpacing int` (default: 2)
- [ ] Update `NewBaseIntent()` constructor:
  - Initialize `terminalInfo` with `terminal.NewInfo()`
  - Set `logoSpacing` to 2 (default)
  - Don't create logo yet (lazy initialization)
- [ ] Commit: `feat(intents): add terminal info and logo fields to BaseIntent`

#### 2.2 Implement TerminalAwareIntent (15 min)
- [ ] Add `UpdateTerminalInfo(info *terminal.Info)` method to BaseIntent:
  - Store info in `terminalInfo` field
  - Mark as valid if width/height > 0
- [ ] Add `GetTerminalInfo() *terminal.Info` method to BaseIntent
- [ ] Verify BaseIntent implements TerminalAwareIntent interface
- [ ] Add documentation comments explaining terminal awareness
- [ ] Commit: `feat(intents): implement TerminalAwareIntent in BaseIntent`

#### 2.3 Add logo management methods (15 min)
- [ ] Add `InitializeLogo()` method to BaseIntent:
  - Create new ASCIILogo instance
  - Set static mode (no animation)
  - Set external centering to true
  - Store in `logo` field
- [ ] Add `GetLogo() *components.ASCIILogo` method:
  - Lazy initialize if nil
  - Return logo instance
- [ ] Add `SetLogoSpacing(spacing int)` method
- [ ] Add `GetLogoSpacing() int` method
- [ ] Commit: `feat(intents): add logo management to BaseIntent`

### Phase 3: Create View Helper Functions (60 min)

#### 3.1 Create view_helpers.go structure (20 min)
- [ ] Create `internal/cli/intents/view_helpers.go`
- [ ] Add package documentation
- [ ] Add `CreateStandardView(b *BaseIntent) *components.StandardView`:
  - Create new StandardView with terminal info
  - Attach logo with configured spacing
  - Set UseFullWidth to true
  - Return configured view
- [ ] Add `CreateStandardViewWithBreadcrumbs(b *BaseIntent, crumbs ...string) *components.StandardView`:
  - Call CreateStandardView
  - Add breadcrumbs
  - Return view
- [ ] Commit: `feat(intents): add basic view helper functions`

#### 3.2 Add modal helper functions (20 min)
- [ ] Add `ShowErrorModal(view *components.StandardView, err error) *components.StandardView`:
  - Create error modal from error message
  - Set bell to true
  - Attach to view
  - Return view
- [ ] Add `ShowLoadingModal(view *components.StandardView, message string, cancellable bool) *components.StandardView`:
  - Create loading modal
  - Set cancellable flag
  - Attach to view
  - Return view
- [ ] Add `ShowProgressModal(view *components.StandardView, title, message string, progress float64) *components.StandardView`:
  - Create progress modal
  - Set progress value (0.0 - 1.0)
  - Attach to view
  - Return view
- [ ] Add `ShowSuccessModal(view *components.StandardView, message string) *components.StandardView`:
  - Create success modal
  - Set auto-dismiss to 3 seconds
  - Attach to view
  - Return view
- [ ] Commit: `feat(intents): add modal helper functions`

#### 3.3 Add footer helper functions (20 min)
- [ ] Add `StandardHelpFooter(shortcuts map[string]string) string`:
  - Format shortcuts as "key1 Action1  key2 Action2"
  - Return formatted string
- [ ] Add predefined footer helpers:
  - `NavigationFooter() string` - Up/Down/Enter/Esc shortcuts
  - `FormFooter() string` - Tab/Shift+Tab/Enter/Esc shortcuts
  - `ListFooter() string` - Navigation + search shortcuts
  - `DetailViewFooter() string` - Scroll + back shortcuts
  - `ModalFooter(actions []string) string` - Modal-specific actions
- [ ] Add `CombineFooters(footers ...string) string`:
  - Combine multiple footer strings
  - Separate with "  |  "
  - Return combined string
- [ ] Commit: `feat(intents): add footer helper functions`

### Phase 4: Update Router for Terminal Info (30 min)

#### 4.1 Review current router implementation (10 min)
- [ ] Open `internal/cli/intents/router.go`
- [ ] Verify `terminalInfo` field exists
- [ ] Verify `UpdateTerminalInfo()` is called in `Update()`
- [ ] Check terminal info propagation in `ActivateIntent()`

#### 4.2 Enhance terminal info propagation (20 min)
- [ ] In `ActivateIntent()`:
  - After creating intent, check if it implements TerminalAwareIntent
  - If yes, immediately call `UpdateTerminalInfo(r.terminalInfo)`
  - Ensure terminal info is passed before Init() is called
- [ ] In `Update()` method for WindowSizeMsg:
  - Update router's terminal info
  - Propagate to active intent if it's TerminalAwareIntent
  - Return cmd to trigger re-render
- [ ] Add debug logging (using existing logger):
  - Log when terminal info is updated
  - Log terminal dimensions
  - Use Debug level to avoid noise
- [ ] Commit: `feat(intents): enhance terminal info propagation in router`

### Phase 5: Update App Model (20 min)

#### 5.1 Ensure terminal info flows from app to router (20 min)
- [ ] Open `internal/cli/app/app.go`
- [ ] In `New()` constructor:
  - Verify intentRouter gets terminal info
  - Ensure initial window size triggers terminal info update
- [ ] In `Update()` method for WindowSizeMsg:
  - Update model's width and height
  - Update terminalInfo
  - Pass to intentRouter via UpdateTerminalInfo()
  - Ensure StateIntent forwards to active intent
- [ ] Verify main menu already uses terminal info correctly
- [ ] Commit: `feat(app): ensure terminal info flows to all intents`

### Phase 6: Add Intent Helper Methods (30 min)

#### 6.1 Add convenience methods to BaseIntent (30 min)
- [ ] Open `internal/cli/intents/contract.go`
- [ ] Add `CreateView() *components.StandardView`:
  - Wrapper for CreateStandardView(b)
  - Makes it easier for intents to create views
- [ ] Add `CreateViewWithBreadcrumbs(crumbs ...string) *components.StandardView`:
  - Wrapper for CreateStandardViewWithBreadcrumbs
- [ ] Add state helpers (to be used by intents):
  - `isLoading bool` field
  - `loadingMessage string` field
  - `errorState error` field
  - `successMessage string` field
  - `successTime time.Time` field
  - `SetLoading(message string)` method
  - `ClearLoading()` method
  - `SetError(err error)` method
  - `ClearError()` method
  - `SetSuccess(message string)` method
  - `ClearSuccess()` method
  - `ShouldShowSuccess() bool` (checks if within 3s window)
- [ ] Commit: `feat(intents): add state management helpers to BaseIntent`

### Phase 7: Write Tests (45 min)

#### 7.1 Test BaseIntent enhancements (20 min)
- [ ] Open or create `internal/cli/intents/contract_test.go`
- [ ] Add test: BaseIntent initializes with terminal info
- [ ] Add test: UpdateTerminalInfo stores info correctly
- [ ] Add test: GetTerminalInfo returns stored info
- [ ] Add test: InitializeLogo creates logo in static mode
- [ ] Add test: GetLogo lazy initializes logo
- [ ] Add test: SetLogoSpacing/GetLogoSpacing work correctly
- [ ] Add test: State helpers (SetLoading, SetError, SetSuccess) work
- [ ] Add test: ShouldShowSuccess respects 3s window
- [ ] Commit: `test(intents): add BaseIntent enhancement tests`

#### 7.2 Test view helpers (25 min)
- [ ] Create `internal/cli/intents/view_helpers_test.go`
- [ ] Add test: CreateStandardView creates view with logo and terminal info
- [ ] Add test: CreateStandardViewWithBreadcrumbs includes breadcrumbs
- [ ] Add test: ShowErrorModal attaches error modal with bell
- [ ] Add test: ShowLoadingModal attaches loading modal
- [ ] Add test: ShowProgressModal attaches progress modal with value
- [ ] Add test: ShowSuccessModal attaches success modal with auto-dismiss
- [ ] Add test: StandardHelpFooter formats shortcuts correctly
- [ ] Add test: Predefined footer helpers return expected strings
- [ ] Add test: CombineFooters joins with separator
- [ ] Commit: `test(intents): add view helper tests`

### Phase 8: Integration Testing (30 min)
- [ ] Run all intent tests: `go test -v ./internal/cli/intents/...`
- [ ] Run with race detector: `go test -race ./internal/cli/intents/...`
- [ ] Run with coverage: `go test -cover ./internal/cli/intents/...`
- [ ] Verify coverage maintained or improved
- [ ] Fix any failing tests
- [ ] Create manual test in `cmd/test_intent_view/main.go`:
  - Create mock intent using BaseIntent
  - Call CreateView() and CreateViewWithBreadcrumbs()
  - Render and verify output
  - Test with different terminal sizes
- [ ] Commit if fixes needed: `fix(intents): resolve test failures`

### Phase 9: Documentation (15 min)
- [ ] Add godoc comments to all new methods in BaseIntent
- [ ] Add godoc comments to all view helper functions
- [ ] Add usage examples in comments:
  - How to create a standard view from an intent
  - How to add modals to views
  - How to use footer helpers
- [ ] Add code example showing complete pattern:
```go
// Example intent View() method
func (i *MyIntent) View() string {
    view := i.CreateViewWithBreadcrumbs("Main Menu", "My Intent", i.stateName)
    
    if i.isLoading {
        ShowLoadingModal(view, i.loadingMessage, true)
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    if i.ShouldShowSuccess() {
        ShowSuccessModal(view, i.successMessage)
    }
    
    view.WithContent(i.renderContent())
    view.WithHelp(CombineFooters(NavigationFooter(), "q Quit"))
    view.WithFooterSeparator(true)
    
    return view.Render()
}
```
- [ ] Commit: `docs(intents): add documentation for view helpers`

### Phase 10: Final Verification (15 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./internal/cli/intents/...`
- [ ] Format code: `go fmt ./internal/cli/intents/...`
- [ ] Review all commits for atomicity
- [ ] Verify commit messages follow conventional format
- [ ] Run full test suite: `go test ./...`
- [ ] Verify no regressions in app or components
- [ ] Test terminal resize handling manually

## Testing Instructions

### Automated Tests
```bash
# Run all intent tests
go test -v ./internal/cli/intents/

# Run with coverage
go test -cover ./internal/cli/intents/

# Run with race detector
go test -race ./internal/cli/intents/

# Run specific test
go test -v ./internal/cli/intents/ -run TestBaseIntent
```

### Manual Testing
Create a test program to verify infrastructure:

```go
// cmd/test_intent_view/main.go
package main

import (
    "fmt"
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/cli/terminal"
)

func main() {
    // Create base intent
    base := intents.NewBaseIntent()
    
    // Update terminal info
    termInfo := terminal.NewInfo()
    termInfo.Width = 120
    termInfo.Height = 40
    termInfo.IsValid = true
    base.UpdateTerminalInfo(termInfo)
    
    // Initialize logo
    base.InitializeLogo()
    
    // Create standard view
    view := intents.CreateStandardViewWithBreadcrumbs(base, "Main Menu", "Test Intent")
    view.WithContent("Test content here")
    view.WithHelp(intents.CombineFooters(intents.NavigationFooter(), "q Quit"))
    view.WithFooterSeparator(true)
    
    fmt.Println(view.Render())
}
```

Test scenarios:
1. BaseIntent creates logo correctly
2. Terminal info propagates to view
3. View helpers create properly configured views
4. Modal helpers attach modals correctly
5. Footer helpers generate correct shortcut text
6. Different terminal sizes handled properly

## Acceptance Criteria
- [ ] BaseIntent has terminalInfo and logo fields
- [ ] BaseIntent implements TerminalAwareIntent interface
- [ ] Logo initializes in static mode with external centering
- [ ] Terminal info propagates from app → router → intent
- [ ] Window resize updates propagate to active intent
- [ ] CreateStandardView() helper creates configured view
- [ ] Modal helpers (Error, Loading, Progress, Success) work correctly
- [ ] Footer helpers generate properly formatted shortcuts
- [ ] Predefined footers (Navigation, Form, List, DetailView) available
- [ ] State management helpers (SetLoading, SetError, SetSuccess) work
- [ ] ShouldShowSuccess() respects 3-second window
- [ ] All tests pass (100%)
- [ ] Code coverage maintained or improved
- [ ] Code passes linting and formatting
- [ ] Compliance check passes
- [ ] No regressions in existing functionality
- [ ] Terminal resize handled gracefully

## Rollback Plan
If issues are discovered:
1. Identify problematic commit(s)
2. Run: `git revert <commit-hash>`
3. Alternative: `git reset --hard <previous-working-commit>`
4. Re-run tests to verify working state
5. Review and fix issues before re-implementing

## Notes
- Logo should be lazy-initialized to avoid overhead
- Terminal info should default to valid state even if not updated
- State helpers are optional but recommended for consistency
- Footer helpers make it easy to maintain consistent shortcuts
- All changes should be backward compatible with existing intents

## Dependencies
- Requires Task 12 (StandardView, Modal, LoadingMessages components)
- Requires existing terminal.Info infrastructure
- Requires existing TerminalAwareIntent interface
- Uses existing BaseIntent as foundation
