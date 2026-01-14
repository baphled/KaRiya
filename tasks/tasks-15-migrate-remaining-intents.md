---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 15: Migrate GenerateCV, ExportArtifact, and ConfigureSystem Intents to StandardView

## Overview
- **Goal**: Migrate the remaining primary intents (GenerateCV, ExportArtifact, ConfigureSystem) to use StandardView with logo, modals, and standardized layout
- **Time Estimate**: 4.5 hours
- **Prerequisites**: Task 14 completed (CaptureEvent and BrowseTimeline migrated successfully)
- **Status**: Not Started

## Motivation
Complete the standardization of all primary user-facing intents. Following the patterns established in Task 14, ensure consistency across the entire application.

## Files to Modify
- [ ] `internal/cli/intents/generate_cv.go` (data structures)
- [ ] `internal/cli/intents/generate_cv_intent.go` (implementation)
- [ ] `internal/cli/intents/generate_cv_test.go` (tests)
- [ ] `internal/cli/intents/export_artifact.go` (data structures)
- [ ] `internal/cli/intents/export_artifact_intent.go` (implementation)
- [ ] `internal/cli/intents/export_artifact_test.go` (tests)
- [ ] `internal/cli/intents/configure_system.go` (data structures)
- [ ] `internal/cli/intents/configure_system_intent.go` (implementation)
- [ ] `internal/cli/intents/configure_system_test.go` (tests)

## Implementation Checklist

### Phase 1: Preparation (10 min)
- [ ] Create task file
- [ ] Run compliance check (baseline): `make check-compliance`
- [ ] Verify Task 14 completed successfully
- [ ] Review patterns from CaptureEvent and BrowseTimeline migrations
- [ ] Run existing tests: `go test ./internal/cli/intents/...`
- [ ] Document current test coverage baseline

### Phase 2: Migrate GenerateCVIntent (110 min)

#### 2.1 Update GenerateCVIntent structure (15 min)
- [ ] Open `internal/cli/intents/generate_cv.go`
- [ ] Verify BaseIntent is embedded
- [ ] Add loading rotator field
- [ ] Add progress tracking for multi-step CV generation
- [ ] Initialize in constructor:
  - Call `InitializeLogo()`
  - Create loading rotator with CV-specific messages:
    - "🔍 Analyzing career events..."
    - "📊 Calculating impact metrics..."
    - "✨ Generating professional bullets..."
    - "📝 Formatting final document..."
    - "✅ CV ready!"
- [ ] Commit: `feat(cv): update GenerateCVIntent for StandardView`

#### 2.2 Refactor View method (25 min)
- [ ] Open `internal/cli/intents/generate_cv_intent.go`
- [ ] Locate main `View() string` method (around line 324)
- [ ] Refactor to StandardView pattern:
```go
func (i *GenerateCVIntent) View() string {
    if !i.active {
        return "GenerateCV intent is not active"
    }
    
    view := i.CreateViewWithBreadcrumbs("Main Menu", "Generate CV", i.getStateName())
    
    // Handle modal states with progress for generation
    if i.isLoading {
        if i.generationProgress > 0 {
            ShowProgressModal(view, "Generating CV", 
                i.loadingRotator.GetCurrent(), i.generationProgress)
        } else {
            ShowLoadingModal(view, i.loadingRotator.GetCurrent(), false)
        }
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    if i.ShouldShowSuccess() {
        ShowSuccessModal(view, i.successMessage)
    }
    
    content := i.getStateContent()
    view.WithContent(content)
    
    help := i.getContextHelp()
    view.WithHelp(help).WithFooterSeparator(true)
    
    return view.Render()
}
```
- [ ] Commit: `refactor(cv): update View method to use StandardView`

#### 2.3 Create state helpers (15 min)
- [ ] Add `getStateName() string` method:
```go
func (i *GenerateCVIntent) getStateName() string {
    switch i.state {
    case CVStateSelectProfile:
        return "Select Profile"
    case CVStateSelectAudience:
        return "Select Audience"
    case CVStateGenerating:
        return "Generating"
    case CVStatePreview:
        return "Preview"
    case CVStateReview:
        return "Review"
    case CVStateConfirm:
        return "Confirm"
    case CVStateExportFormat:
        return "Export Format"
    case CVStateExportLocation:
        return "Export Location"
    case CVStateExporting:
        return "Exporting"
    case CVStateComplete:
        return "Complete"
    default:
        return string(i.state)
    }
}
```
- [ ] Add `getStateContent() string` method delegating to render methods
- [ ] Commit: `feat(cv): add state helper methods`

#### 2.4 Update render methods (30 min)
- [ ] Rename view methods to render methods:
  - `viewSelectProfile()` → `renderSelectProfile()`
  - `viewSelectAudience()` → `renderSelectAudience()`
  - `viewGenerating()` → `renderGenerating()` (now just returns empty - modal shows progress)
  - `viewPreview()` → `renderPreview()`
  - `viewReview()` → `renderReview()`
  - `viewConfirm()` → `renderConfirm()`
- [ ] Update each render method:
  - Remove header/footer elements
  - Use full width
  - Apply lipgloss styling
  - Return clean content
- [ ] For renderGenerating(), return simple placeholder:
  - Modal will show progress, content can be minimal
- [ ] Commit: `refactor(cv): update render methods for full width`

#### 2.5 Create getContextHelp method (10 min)
- [ ] Add `getContextHelp() string` method:
```go
func (i *GenerateCVIntent) getContextHelp() string {
    base := "q Quit  m Main Menu"
    
    switch i.state {
    case CVStateSelectProfile:
        return CombineFooters(NavigationFooter(), "Esc Cancel", base)
    case CVStateSelectAudience:
        return CombineFooters(NavigationFooter(), "Esc Back", base)
    case CVStateGenerating:
        return base // No actions during generation
    case CVStatePreview:
        return CombineFooters("e Edit  c Confirm  Esc Back", base)
    case CVStateReview:
        return CombineFooters("Enter Continue  Esc Back", base)
    case CVStateConfirm:
        return CombineFooters("y Confirm  e/x Export  Esc Back", base)
    default:
        return base
    }
}
```
- [ ] Commit: `feat(cv): add context-aware help text`

#### 2.6 Add progress tracking (15 min)
- [ ] Add `generationProgress float64` field
- [ ] In CV generation logic, update progress:
```go
i.SetLoading("Generating CV...")
i.generationProgress = 0.0

// Step 1: Analyze events
i.generationProgress = 0.25
i.loadingRotator.Rotate()

// Step 2: Calculate metrics
i.generationProgress = 0.50
i.loadingRotator.Rotate()

// Step 3: Generate bullets
i.generationProgress = 0.75
i.loadingRotator.Rotate()

// Step 4: Format document
i.generationProgress = 1.0
i.loadingRotator.Rotate()

i.ClearLoading()
i.SetSuccess("✅ CV generated successfully!")
```
- [ ] Commit: `feat(cv): add progress tracking for CV generation`

### Phase 3: Migrate ExportArtifactIntent (80 min)

#### 3.1 Update ExportArtifactIntent structure (10 min)
- [ ] Open `internal/cli/intents/export_artifact.go`
- [ ] Verify BaseIntent is embedded
- [ ] Add loading rotator field
- [ ] Add progress tracking for export
- [ ] Initialize in constructor:
  - Call `InitializeLogo()`
  - Create loading rotator with export messages:
    - "📝 Preparing export..."
    - "💾 Writing file..."
    - "✨ Finalizing..."
- [ ] Commit: `feat(export): update ExportArtifactIntent for StandardView`

#### 3.2 Refactor View method (20 min)
- [ ] Open `internal/cli/intents/export_artifact_intent.go`
- [ ] Locate main `View() string` method (around line 40)
- [ ] Refactor to StandardView pattern:
```go
func (e *ExportArtifactIntent) View() string {
    if !e.active {
        return "ExportArtifact intent is not active"
    }
    
    view := e.CreateViewWithBreadcrumbs("Main Menu", "Export Artifact", e.getStateName())
    
    if e.isLoading {
        if e.exportProgress > 0 {
            ShowProgressModal(view, "Exporting",
                e.loadingRotator.GetCurrent(), e.exportProgress)
        } else {
            ShowLoadingModal(view, e.loadingRotator.GetCurrent(), false)
        }
    }
    if e.errorState != nil {
        ShowErrorModal(view, e.errorState)
    }
    if e.ShouldShowSuccess() {
        ShowSuccessModal(view, e.successMessage)
    }
    
    content := e.getStateContent()
    view.WithContent(content)
    
    help := e.getContextHelp()
    view.WithHelp(help).WithFooterSeparator(true)
    
    return view.Render()
}
```
- [ ] Commit: `refactor(export): update View method to use StandardView`

#### 3.3 Create state helpers (15 min)
- [ ] Add `getStateName() string` method
- [ ] Add `getStateContent() string` method
- [ ] Add `getContextHelp() string` method with state-specific shortcuts
- [ ] Commit: `feat(export): add state helper methods`

#### 3.4 Update render methods (20 min)
- [ ] Rename view methods to render methods
- [ ] Update for full width and clean content
- [ ] Remove manual headers/footers
- [ ] Apply lipgloss styling
- [ ] Commit: `refactor(export): update render methods for full width`

#### 3.5 Add export progress (15 min)
- [ ] Add `exportProgress float64` field
- [ ] Track progress during export:
  - 0.33 - Preparing data
  - 0.66 - Writing file/clipboard
  - 1.0 - Complete
- [ ] Show success modal with file path on completion
- [ ] Commit: `feat(export): add progress tracking for export operations`

### Phase 4: Migrate ConfigureSystemIntent (80 min)

#### 4.1 Update ConfigureSystemIntent structure (10 min)
- [ ] Open `internal/cli/intents/configure_system.go`
- [ ] Verify BaseIntent is embedded
- [ ] Add loading rotator field
- [ ] Initialize in constructor:
  - Call `InitializeLogo()`
  - Create loading rotator with config messages:
    - "⚙️ Validating configuration..."
    - "💾 Saving settings..."
    - "🔄 Applying changes..."
- [ ] Commit: `feat(config): update ConfigureSystemIntent for StandardView`

#### 4.2 Refactor View method (20 min)
- [ ] Open `internal/cli/intents/configure_system_intent.go`
- [ ] Locate main `View() string` method (around line 40)
- [ ] Refactor to StandardView pattern (similar to above intents)
- [ ] Commit: `refactor(config): update View method to use StandardView`

#### 4.3 Create state helpers (15 min)
- [ ] Add `getStateName() string` method
- [ ] Add `getStateContent() string` method
- [ ] Add `getContextHelp() string` method
- [ ] Commit: `feat(config): add state helper methods`

#### 4.4 Update render methods (20 min)
- [ ] Rename view methods to render methods
- [ ] Update for full width
- [ ] Remove manual headers/footers
- [ ] Apply lipgloss styling for settings display
- [ ] Commit: `refactor(config): update render methods for full width`

#### 4.5 Add configuration states (15 min)
- [ ] Add loading during configuration apply
- [ ] Add success modal: "✅ Configuration updated!"
- [ ] Add error handling with modal
- [ ] Commit: `feat(config): add loading and success states`

### Phase 5: Update Tests (70 min)

#### 5.1 Update GenerateCVIntent tests (25 min)
- [ ] Open `internal/cli/intents/generate_cv_test.go`
- [ ] Update view output tests for StandardView structure
- [ ] Add tests for:
  - Logo presence
  - Breadcrumbs
  - Footer separator
  - Progress modal during generation
  - Success modal on completion
- [ ] Fix any failing tests
- [ ] Run: `go test -v ./internal/cli/intents/ -run TestGenerateCV`
- [ ] Commit: `test(cv): update tests for StandardView migration`

#### 5.2 Update ExportArtifactIntent tests (20 min)
- [ ] Open `internal/cli/intents/export_artifact_test.go`
- [ ] Update view output tests
- [ ] Add tests for progress tracking and success modal
- [ ] Fix any failing tests
- [ ] Run: `go test -v ./internal/cli/intents/ -run TestExportArtifact`
- [ ] Commit: `test(export): update tests for StandardView migration`

#### 5.3 Update ConfigureSystemIntent tests (25 min)
- [ ] Open `internal/cli/intents/configure_system_test.go`
- [ ] Update view output tests
- [ ] Add tests for configuration states and modals
- [ ] Fix any failing tests
- [ ] Run: `go test -v ./internal/cli/intents/ -run TestConfigureSystem`
- [ ] Commit: `test(config): update tests for StandardView migration`

### Phase 6: Integration Testing (45 min)

#### 6.1 Run full test suite (15 min)
- [ ] Run all intent tests: `go test -v ./internal/cli/intents/...`
- [ ] Run with race detector: `go test -race ./internal/cli/intents/...`
- [ ] Run with coverage: `go test -cover ./internal/cli/intents/...`
- [ ] Verify coverage maintained or improved
- [ ] Fix any failing tests
- [ ] Commit if fixes needed: `fix(intents): resolve integration test failures`

#### 6.2 Manual testing - GenerateCV (10 min)
- [ ] Build and run application
- [ ] Navigate to Generate CV (g key)
- [ ] Verify logo, breadcrumbs, footer separator
- [ ] Select profile, verify view updates
- [ ] Select audience, verify CV generation starts
- [ ] Verify progress modal shows with percentage
- [ ] Verify loading messages rotate
- [ ] Verify success modal appears and auto-dismisses
- [ ] Test export workflow from CV generation

#### 6.3 Manual testing - ExportArtifact (10 min)
- [ ] Navigate to Export (from CV or other artifact)
- [ ] Verify logo, breadcrumbs, footer separator
- [ ] Select format, verify options
- [ ] Select location, verify options
- [ ] Watch export progress modal
- [ ] Verify success modal with file path
- [ ] Test different export formats (Text, Markdown, YAML)
- [ ] Test different locations (File, Clipboard)

#### 6.4 Manual testing - ConfigureSystem (10 min)
- [ ] Navigate to Configure System
- [ ] Verify logo, breadcrumbs, footer separator
- [ ] Select domain, verify settings display
- [ ] Modify settings, verify changes tracked
- [ ] Apply configuration, verify loading modal
- [ ] Verify success modal appears
- [ ] Test error handling (invalid config)

### Phase 7: Documentation (10 min)
- [ ] Add comments explaining StandardView usage in each intent
- [ ] Document progress tracking patterns
- [ ] Add notes about modal state management
- [ ] Commit: `docs(intents): document StandardView migration for remaining intents`

### Phase 8: Final Verification (15 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./internal/cli/intents/...`
- [ ] Format code: `go fmt ./internal/cli/intents/...`
- [ ] Review all commits for atomicity
- [ ] Verify commit messages follow conventional format
- [ ] Run full test suite: `go test ./...`
- [ ] Verify no regressions
- [ ] Test on different terminal sizes

## Testing Instructions

### Automated Tests
```bash
# Run all intent tests
go test -v ./internal/cli/intents/

# Run specific tests
go test -v ./internal/cli/intents/ -run TestGenerateCV
go test -v ./internal/cli/intents/ -run TestExportArtifact
go test -v ./internal/cli/intents/ -run TestConfigureSystem

# Run with coverage
go test -cover ./internal/cli/intents/

# Run with race detector
go test -race ./internal/cli/intents/
```

### Manual Testing Checklist

#### GenerateCV Intent
- [ ] Logo visible with 2-line spacing
- [ ] Breadcrumbs update through workflow
- [ ] Footer separator visible
- [ ] Context help state-aware
- [ ] Profile selection works
- [ ] Audience selection works
- [ ] Progress modal shows during generation
- [ ] Loading messages rotate (5 messages)
- [ ] Success modal auto-dismisses
- [ ] Export workflow accessible from confirmation

#### ExportArtifact Intent
- [ ] Logo visible
- [ ] Breadcrumbs show path
- [ ] Footer separator visible
- [ ] Format selection (Text, Markdown, YAML)
- [ ] Location selection (File, Clipboard)
- [ ] Progress modal during export
- [ ] Success modal shows file path
- [ ] Different formats work correctly

#### ConfigureSystem Intent
- [ ] Logo visible
- [ ] Breadcrumbs show navigation
- [ ] Domain selection works
- [ ] Settings display correctly
- [ ] Changes tracked
- [ ] Apply shows loading modal
- [ ] Success modal on completion
- [ ] Error handling works

## Acceptance Criteria
- [ ] GenerateCVIntent uses StandardView throughout
- [ ] GenerateCVIntent shows progress modal with percentage
- [ ] GenerateCVIntent loading messages rotate (5 messages)
- [ ] GenerateCVIntent success modal auto-dismisses
- [ ] ExportArtifactIntent uses StandardView throughout
- [ ] ExportArtifactIntent shows export progress
- [ ] ExportArtifactIntent success modal shows file path
- [ ] ConfigureSystemIntent uses StandardView throughout
- [ ] ConfigureSystemIntent shows loading during apply
- [ ] ConfigureSystemIntent success modal on completion
- [ ] All intents show logo with 2-line spacing
- [ ] All intents use full terminal width
- [ ] All intents show footer separator
- [ ] All intents have context-aware help
- [ ] All existing tests pass (100%)
- [ ] New tests added for StandardView features
- [ ] Code coverage maintained or improved
- [ ] Code passes linting and formatting
- [ ] Compliance check passes
- [ ] No regressions
- [ ] Manual testing successful

## Rollback Plan
If issues are discovered:
1. Identify problematic commit(s)
2. Run: `git revert <commit-hash>`
3. Alternative: `git reset --hard <previous-working-commit>`
4. Re-run tests to verify working state
5. Review and fix issues before re-implementing

## Notes
- Follow the patterns established in Task 14
- Progress modals are especially important for CV generation and export
- Loading message rotation adds polish to long operations
- Success modals should include useful information (file paths, etc.)
- Ensure all modals have appropriate timing (fade-in, auto-dismiss)
- Terminal resize should work smoothly

## Dependencies
- Requires Task 14 (CaptureEvent and BrowseTimeline migrated)
- Uses StandardView components from Task 12
- Uses infrastructure from Task 13
- Requires existing intent implementations
