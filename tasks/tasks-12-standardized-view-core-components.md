---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 12: Core Components for Standardized Views

## Overview
- **Goal**: Create foundational components for standardized centered views with logo, modal system, and loading messages
- **Time Estimate**: 4-5 hours
- **Prerequisites**: Understanding of BubbleTea, Lipgloss, SmartContainer, ASCIILogo components
- **Status**: Not Started

## Motivation
Standardize all TUI views to use a consistent layout with logo always visible at the top, content centered in the middle, and context-aware help at the bottom. This provides a professional, consistent user experience across all screens.

## Files to Create
- [ ] `internal/cli/components/standard_view.go` (main component)
- [ ] `internal/cli/components/standard_view_test.go` (tests)
- [ ] `internal/cli/components/modal.go` (modal system)
- [ ] `internal/cli/components/modal_test.go` (modal tests)
- [ ] `internal/cli/components/loading_messages.go` (loading message rotator)
- [ ] `internal/cli/components/loading_messages_test.go` (loading tests)

## Files to Modify
- [ ] `internal/cli/components/smart_container.go` (add dimming support for modals)

## Implementation Checklist

### Phase 1: Preparation (10 min)
- [ ] Create task file
- [ ] Run compliance check (baseline): `make check-compliance`
- [ ] Review existing SmartContainer implementation
- [ ] Review existing ASCIILogo component
- [ ] Verify test suite passes: `go test ./internal/cli/components/...`

### Phase 2: Create StandardView Component (60 min)

#### 2.1 Create basic structure (20 min)
- [ ] Create `internal/cli/components/standard_view.go`
- [ ] Add package documentation
- [ ] Define `StandardView` struct with fields:
  - `ShowLogo bool` (default: true)
  - `Logo *ASCIILogo`
  - `LogoSpacing int` (default: 2)
  - `ShowHeader bool`
  - `Breadcrumbs []string`
  - `Title string`
  - `Subtitle string`
  - `Content string`
  - `ContentStyle lipgloss.Style`
  - `ShowModal bool`
  - `Modal *ModalContent`
  - `HelpText string`
  - `ShowFooter bool`
  - `ShowFooterSeparator bool`
  - `TerminalInfo *terminal.Info`
  - `UseFullWidth bool` (default: true)
- [ ] Add necessary imports (lipgloss, strings, terminal package)
- [ ] Commit: `feat(components): add StandardView struct and basic structure`

#### 2.2 Add constructor and builder methods (20 min)
- [ ] Add `NewStandardView(info *terminal.Info) *StandardView` constructor
- [ ] Add `WithLogo(logo *ASCIILogo, spacing int) *StandardView`
- [ ] Add `WithBreadcrumbs(crumbs ...string) *StandardView`
- [ ] Add `WithTitle(title, subtitle string) *StandardView`
- [ ] Add `WithContent(content string) *StandardView`
- [ ] Add `WithContentStyle(style lipgloss.Style) *StandardView`
- [ ] Add `WithHelp(helpText string) *StandardView`
- [ ] Add `WithFooterSeparator(show bool) *StandardView`
- [ ] Add `ShowModalOverlay(modal *ModalContent) *StandardView`
- [ ] Add `UseFullWidth(full bool) *StandardView`
- [ ] Commit: `feat(components): add StandardView builder methods`

#### 2.3 Implement rendering logic (20 min)
- [ ] Add `Render() string` method
- [ ] Implement logo rendering with spacing:
  - Add `LogoSpacing` blank lines before logo
  - Render logo using `logo.ViewStatic()`
  - Add 1 blank line after logo
- [ ] Implement breadcrumb/header rendering:
  - Join breadcrumbs with " > "
  - Style with subtle color
  - Center horizontally
- [ ] Implement content rendering:
  - Apply ContentStyle if provided
  - Use full width if `UseFullWidth` is true
- [ ] Implement footer rendering:
  - Add visual separator if `ShowFooterSeparator` is true
  - Style separator as horizontal line (repeated "─")
  - Add blank line before help text
  - Render help text centered
- [ ] Add modal overlay rendering (if `ShowModal` is true)
- [ ] Use SmartContainer for final centering with `CenterBoth`
- [ ] Commit: `feat(components): implement StandardView render method`

### Phase 3: Create Modal System (60 min)

#### 3.1 Define modal types and structures (15 min)
- [ ] Create `internal/cli/components/modal.go`
- [ ] Define `ModalType` enum:
  - `ModalError`
  - `ModalLoading`
  - `ModalProgress`
  - `ModalSuccess`
  - `ModalWarning`
- [ ] Define `ModalContent` struct with fields:
  - `Type ModalType`
  - `Title string`
  - `Message string`
  - `Progress float64` (0.0 - 1.0)
  - `Actions []string`
  - `FadeInDuration time.Duration` (default: 150ms)
  - `AutoDismiss time.Duration`
  - `Bell bool` (sound alert)
  - `Cancellable bool`
  - `fadeStartTime time.Time` (internal)
- [ ] Add necessary imports
- [ ] Commit: `feat(components): add modal types and structures`

#### 3.2 Implement modal constructors (10 min)
- [ ] Add `NewErrorModal(title, message string) *ModalContent`
- [ ] Add `NewLoadingModal(message string, cancellable bool) *ModalContent`
- [ ] Add `NewProgressModal(title, message string, progress float64) *ModalContent`
- [ ] Add `NewSuccessModal(message string) *ModalContent` (auto-dismiss 3s)
- [ ] Add `NewWarningModal(title, message string) *ModalContent`
- [ ] Commit: `feat(components): add modal constructor functions`

#### 3.3 Implement modal rendering (35 min)
- [ ] Add `Render(terminalWidth, terminalHeight int) string` method
- [ ] Implement fade-in animation:
  - Calculate opacity based on elapsed time since `fadeStartTime`
  - Use lipgloss alpha channel for dimming
  - Fade in over `FadeInDuration` (150ms default)
- [ ] Implement modal box styling by type:
  - Error: Red border, "⚠️" icon
  - Loading: Blue border, rotating spinner
  - Progress: Blue border with progress bar
  - Success: Green border, "✅" icon
  - Warning: Yellow border, "⚠️" icon
- [ ] Add adaptive sizing:
  - Calculate box width based on content (min: 40, max: 80)
  - Calculate box height based on lines (min: 8, max: 20)
  - Center modal in available space
- [ ] Add progress bar rendering for `ModalProgress`:
  - Show percentage (e.g., "65%")
  - Visual bar with filled/unfilled sections
- [ ] Add action buttons if `Actions` is not empty:
  - Render as horizontal list at bottom
  - Style first action as primary (highlighted)
- [ ] Add dismissal hint:
  - "Press Esc to dismiss" for errors
  - "Press Esc to cancel" for cancellable loading
  - Auto-dismiss countdown for success modals
- [ ] Commit: `feat(components): implement modal rendering with fade-in`

### Phase 4: Create Loading Message Rotator (45 min)

#### 4.1 Define loading message types (10 min)
- [ ] Create `internal/cli/components/loading_messages.go`
- [ ] Define `LoadingMessageRotator` struct:
  - `messages []string`
  - `currentIndex int`
  - `rotateInterval time.Duration` (default: 2s)
  - `lastRotation time.Time`
- [ ] Define predefined message sets:
  - `LoadingMessagesCV` for CV generation (5 messages)
  - `LoadingMessagesExport` for export operations (4 messages)
  - `LoadingMessagesGeneric` for general operations (5 messages)
  - `LoadingMessagesFetch` for data fetching (4 messages)
- [ ] Add necessary imports (time)
- [ ] Commit: `feat(components): add loading message rotator structure`

#### 4.2 Implement rotator logic (20 min)
- [ ] Add `NewLoadingMessageRotator(messages []string, interval time.Duration) *LoadingMessageRotator`
- [ ] Add `GetCurrent() string` - Returns current message
- [ ] Add `Rotate() string` - Advances to next message if interval elapsed
- [ ] Add `Reset()` - Resets to first message
- [ ] Add `SetMessages(messages []string)` - Updates message set
- [ ] Implement rotation logic:
  - Check if `rotateInterval` has elapsed since `lastRotation`
  - If yes, increment `currentIndex` (wrap around)
  - Update `lastRotation`
  - Return current message
- [ ] Commit: `feat(components): implement loading message rotation logic`

#### 4.3 Add spinner integration (15 min)
- [ ] Add `Spinner` struct:
  - `frames []string` (default: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
  - `currentFrame int`
  - `lastUpdate time.Time`
  - `frameInterval time.Duration` (default: 80ms)
- [ ] Add `NewSpinner() *Spinner` constructor
- [ ] Add `GetFrame() string` - Returns current spinner frame
- [ ] Add `Advance()` - Advances to next frame
- [ ] Integrate spinner with loading messages:
  - Combine spinner frame with message
  - Format: "{spinner} {message}"
- [ ] Commit: `feat(components): add spinner for loading animations`

### Phase 5: Add SmartContainer Dimming Support (20 min)
- [ ] Open `internal/cli/components/smart_container.go`
- [ ] Add `dimContent bool` field to `SmartContainer`
- [ ] Add `SetDimContent(dim bool) *SmartContainer` method
- [ ] Update `Render()` to apply dimming if enabled:
  - Reduce color intensity using lipgloss Foreground with gray
  - Apply subtle alpha/opacity reduction
- [ ] Commit: `feat(components): add content dimming support to SmartContainer`

### Phase 6: Write Tests (60 min)

#### 6.1 StandardView tests (25 min)
- [ ] Create `internal/cli/components/standard_view_test.go`
- [ ] Add test: StandardView renders logo with spacing
- [ ] Add test: StandardView renders breadcrumbs
- [ ] Add test: StandardView renders content with full width
- [ ] Add test: StandardView renders footer with separator
- [ ] Add test: StandardView renders help text
- [ ] Add test: StandardView uses SmartContainer for centering
- [ ] Add test: Builder methods chain correctly
- [ ] Add test: Default values are set correctly
- [ ] Commit: `test(components): add StandardView tests`

#### 6.2 Modal tests (20 min)
- [ ] Create `internal/cli/components/modal_test.go`
- [ ] Add test: ErrorModal renders with red border and icon
- [ ] Add test: LoadingModal renders with spinner
- [ ] Add test: ProgressModal renders with progress bar
- [ ] Add test: SuccessModal renders with green border
- [ ] Add test: WarningModal renders with yellow border
- [ ] Add test: Modal fade-in animation works
- [ ] Add test: Modal adaptive sizing works
- [ ] Add test: Modal actions render correctly
- [ ] Add test: Bell flag is respected
- [ ] Commit: `test(components): add modal system tests`

#### 6.3 LoadingMessageRotator tests (15 min)
- [ ] Create `internal/cli/components/loading_messages_test.go`
- [ ] Add test: Rotator returns first message initially
- [ ] Add test: Rotator advances after interval
- [ ] Add test: Rotator wraps around to first message
- [ ] Add test: Reset returns to first message
- [ ] Add test: SetMessages updates message set
- [ ] Add test: Spinner advances frames correctly
- [ ] Add test: Spinner integrates with messages
- [ ] Commit: `test(components): add loading message rotator tests`

### Phase 7: Integration Testing (30 min)
- [ ] Run all component tests: `go test -v ./internal/cli/components/...`
- [ ] Run with race detector: `go test -race ./internal/cli/components/...`
- [ ] Run with coverage: `go test -cover ./internal/cli/components/...`
- [ ] Verify coverage > 85% for new components
- [ ] Fix any failing tests
- [ ] Test StandardView rendering manually:
  - Create test program in `cmd/test_standard_view/main.go`
  - Render StandardView with various configurations
  - Verify logo spacing, separator, centering
  - Test modal overlays
  - Test loading messages
- [ ] Commit if fixes needed: `fix(components): resolve test failures`

### Phase 8: Documentation (20 min)
- [ ] Add godoc comments to all public types and methods
- [ ] Add usage examples in comments:
  - StandardView basic usage
  - Modal creation and display
  - Loading message rotation
- [ ] Create code example in comments showing complete flow
- [ ] Commit: `docs(components): add documentation for standardized views`

### Phase 9: Final Verification (15 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./internal/cli/components/...`
- [ ] Format code: `go fmt ./internal/cli/components/...`
- [ ] Review all commits for atomicity
- [ ] Verify commit messages follow conventional format
- [ ] Run full test suite: `go test ./...`
- [ ] Verify no regressions in other packages

## Testing Instructions

### Automated Tests
```bash
# Run all component tests
go test -v ./internal/cli/components/

# Run with coverage
go test -cover ./internal/cli/components/

# Run with race detector
go test -race ./internal/cli/components/

# Run specific test
go test -v ./internal/cli/components/ -run TestStandardView
```

### Manual Testing
Create a test program to visualize the components:

```go
// cmd/test_standard_view/main.go
package main

import (
    "fmt"
    "github.com/baphled/kariya/internal/cli/components"
    "github.com/baphled/kariya/internal/cli/terminal"
)

func main() {
    termInfo := terminal.NewInfo()
    termInfo.Width = 120
    termInfo.Height = 40
    
    logo := components.NewASCIILogo()
    
    view := components.NewStandardView(termInfo).
        WithLogo(logo, 2).
        WithBreadcrumbs("Main Menu", "Test View").
        WithContent("This is test content\nLine 2\nLine 3").
        WithHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back").
        WithFooterSeparator(true).
        UseFullWidth(true)
    
    fmt.Println(view.Render())
}
```

Test scenarios:
1. Basic view without modal
2. View with error modal
3. View with loading modal and spinner
4. View with progress modal at various percentages
5. View with success modal
6. Different terminal sizes (80x24, 120x40, 200x60)

## Acceptance Criteria
- [ ] StandardView component renders logo with configurable spacing (default: 2 lines)
- [ ] Breadcrumbs render as "Crumb1 > Crumb2 > Crumb3"
- [ ] Content uses full terminal width when UseFullWidth is true
- [ ] Footer separator renders as horizontal line of "─" characters
- [ ] Help text is centered and visible
- [ ] SmartContainer centers entire composition
- [ ] Modal system supports 5 types (Error, Loading, Progress, Success, Warning)
- [ ] Modals fade in over 150ms
- [ ] Error modals have red border and "⚠️" icon
- [ ] Loading modals show spinner animation
- [ ] Progress modals show percentage and progress bar
- [ ] Success modals have green border and auto-dismiss after 3s
- [ ] Modal size adapts to content (min 40x8, max 80x20)
- [ ] Loading messages rotate every 2 seconds
- [ ] Spinner animates with 10 frames at 80ms intervals
- [ ] All tests pass (100%)
- [ ] Code coverage > 85% for new components
- [ ] Code passes linting and formatting
- [ ] Compliance check passes
- [ ] No regressions in existing components

## Rollback Plan
If issues are discovered:
1. Identify problematic commit(s)
2. Run: `git revert <commit-hash>`
3. Alternative: `git reset --hard <previous-working-commit>`
4. Re-run tests to verify working state
5. Review and fix issues before re-implementing

## Notes
- Logo should use static mode (no animation) for consistency
- Modal fade-in should be subtle and quick (150ms)
- Loading message rotation should feel natural (2s interval)
- Footer separator should span full content width
- All components should be terminal-size aware
- Consider accessibility: terminal bell for errors, clear visual hierarchy

## Dependencies
- Requires existing SmartContainer component
- Requires existing ASCIILogo component
- Requires terminal.Info for size awareness
- Uses lipgloss for styling
