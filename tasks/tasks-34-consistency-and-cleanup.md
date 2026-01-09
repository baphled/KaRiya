# Task 34: Consistency & Cleanup (MERGED WITH DEAD CODE AUDIT)

**Created**: 2026-01-08
**Updated**: 2026-01-09
**Status**: PARTIALLY COMPLETE (Dead Code Removal Done)
**Priority**: MEDIUM
**Estimated Time**: 4-5 hours (expanded from 2-3 hours)
**Related**: Codebase Audit (2026-01-08), Dead Code Audit (2026-01-09)

---

## Overview

This task was merged with a comprehensive dead code audit on 2026-01-09. The dead code removal is complete. Remaining items (help screen, error display, etc.) still need implementation.

**Original Issues**:
- Help screen ('?' key) does nothing
- CV generation errors never shown to user
- Unused loadingRotator components in 3 intents
- 12+ orphaned message types with no handlers
- 3 .bak files (3,039 total lines of dead code)
- Stub functions that should be removed or implemented
- Unused CLI flags

**Added from Dead Code Audit**:
- Superseded packages (context, layout, workflow, validation)
- Superseded navigation files (handlers, key_handler, registry)
- 17 unused component files
- Unused service files (data_quality, config_initializer, burst_fact/workflow)
- Unused examples directory
- Unused testing.go utilities

---

## Completed: Dead Code Removal (2026-01-09)

### Summary
- **Files removed**: 56
- **Lines removed**: 14,367
- **Dead functions reduced**: 613 → 236 (377 removed, 61% reduction)
- **All tests passing**: 20 packages
- **Staticcheck warnings**: 0

### Packages Removed (Entire Directories)
- [x] `internal/cli/context/` - GlobalContext (Phase 5.1 never started)
- [x] `internal/cli/layout/` - Superseded by StandardView
- [x] `internal/cli/workflow/` - Superseded by Intent architecture
- [x] `internal/cli/validation/` - Never wired up, legacy Task 03 code

### Navigation Files Removed
- [x] `internal/cli/navigation/handlers.go`
- [x] `internal/cli/navigation/key_handler.go`
- [x] `internal/cli/navigation/key_handler_test.go`
- [x] `internal/cli/navigation/registry.go`
- [x] `internal/cli/navigation/registry_test.go`
- [x] `internal/cli/navigation/registry_integration_test.go`

**Kept** (actively used):
- `list_navigator.go`, `constants.go`, `help.go` and their tests
- Added `suite_test.go` for Ginkgo test runner

### Components Removed (17 files)
- [x] `card_container.go` + test
- [x] `centered_container.go`
- [x] `form.go` + test
- [x] `form_container.go` + test
- [x] `form_field_container.go` + test
- [x] `list_container.go` + test
- [x] `list_item.go` + test
- [x] `navigation_menu.go` + test
- [x] `pagination.go` + test
- [x] `progress.go` + test
- [x] `screen_container.go` + test
- [x] `section_container.go` + test
- [x] `smart_container.go` + test
- [x] `spinner.go` + test
- [x] `text_utils.go` + test
- [x] `intent_header.go`

**Kept** (actively used or needed for Task 31):
- `modal_container.go` - Used by `modals.go` (Task 31)
- `audience_relevance_selector.go` - Decision pending (stub function)

### Services Removed
- [x] `internal/service/career/data_quality.go` + test
- [x] `internal/service/career/cv/config_initializer.go` + test
- [x] `internal/service/career/burst_fact/workflow.go` + test
- [x] `internal/service/career/burst_fact/integration_test.go`

**Kept** (PRD requirements):
- `traceability_service.go` - PRD Section 5.5 explainability
- `classification/classifier.go` - Used by category_selector

### Other Files Removed
- [x] `examples/enhanced_capture_example.go` - Orphaned example
- [x] `internal/cli/intents/testing.go` + test - YAGNI
- [x] `internal/cli/app/app.go.bak`
- [x] `internal/cli/models/fact_list.go.bak`
- [x] `internal/cli/models/burst_list.go.bak`

**Kept** (used by tests):
- `internal/repository/career/mocks/mock_helper.go` - Used by service_test.go

---

## Remaining: Original Task 34 Items

### Files to Modify

- [ ] `internal/cli/app/app.go` - Help screen implementation
- [ ] `internal/cli/app/messages.go` - Remove orphaned message types
- [ ] `internal/cli/intents/capture_event_intent.go` - LoadingRotator decision
- [ ] `internal/cli/intents/browse_timeline_intent.go` - LoadingRotator decision
- [ ] `internal/cli/intents/generate_cv_intent.go` - Error display, LoadingRotator
- [ ] `internal/cli/components/audience_relevance_selector.go` - ToggleSelected stub
- [ ] `cmd/cli/main.go` - Unused CLI flags

---

## Implementation Plan

### Phase 1: Help Screen (30 min)

**Location**: `internal/cli/app/app.go:147-148`

**Current (STUB)**:
```go
case "?":
    // Help screen could be implemented here
    return m, nil
```

**Implementation**:
```go
case "?", "h":
    m.showingHelp = !m.showingHelp
    return m, nil

func (m *Model) View() string {
    if m.showingHelp {
        return m.renderHelpScreen()
    }
    // ... normal view
}

func (m *Model) renderHelpScreen() string {
    help := `
╔════════════════════════════════════════════════════════════════════╗
║                    KaRiya Keyboard Reference                       ║
╠════════════════════════════════════════════════════════════════════╣
║ Navigation              Global Shortcuts                           ║
║ ──────────────────────  ──────────────────────────────────────     ║
║ ↑/k   Previous          ?/h    Toggle this help                    ║
║ ↓/j   Next              q      Quit application                    ║
║ ←     Back              c      Capture event                       ║
║ →     Forward           l      List events                         ║
║ Esc   Cancel/Back       m      Main menu                           ║
║                         g      Generate CV                         ║
║                                                                    ║
║ Context Actions                                                    ║
║ ──────────────────────────────────────────────────────────────     ║
║ Enter   Confirm/Select                                             ║
║ Space   Toggle checkbox                                            ║
║ Tab     Next field                                                 ║
║ /       Search                                                     ║
║ f       Filter                                                     ║
║ e       Edit                                                       ║
║ d       Delete                                                     ║
║                                                                    ║
║ Press '?' or 'h' to close this help                                ║
╚════════════════════════════════════════════════════════════════════╝
`
    return help
}
```

**Tasks**:
- [ ] Add showingHelp boolean to Model
- [ ] Implement renderHelpScreen()
- [ ] Toggle help on '?' or 'h' key
- [ ] Show help overlay in View()
- [ ] Test help screen displays and hides

---

### Phase 2: Display CV Generation Errors (30 min)

**Location**: `internal/cli/intents/generate_cv_intent.go:506-558`

**Issue**: `generationError` is stored (line 248) but never displayed

**Fix viewSelectAudience()**:
```go
func (i *GenerateCVIntent) viewSelectAudience() string {
    // Check for generation error
    if i.state.generationError != nil {
        errorMsg := fmt.Sprintf("\n❌ Error: %s\n\n", i.state.generationError.Error())
        errorMsg += "Please try selecting a different audience or check your data.\n\n"
        
        // Show error at top of audience view
        return errorMsg + i.renderAudienceSelection()
    }
    
    return i.renderAudienceSelection()
}

func (i *GenerateCVIntent) renderAudienceSelection() string {
    // ... existing view code
}
```

**Tasks**:
- [ ] Check for generationError in viewSelectAudience()
- [ ] Display error message prominently
- [ ] Extract view rendering to separate method
- [ ] Clear error on successful generation
- [ ] Test error display

---

### Phase 3: Remove Dead Code (1 hour)

#### 3.1 Delete Backup Files
```bash
rm internal/cli/app/app.go.bak
rm internal/cli/models/fact_list.go.bak
rm internal/cli/models/burst_list.go.bak
```

**Tasks**:
- [ ] Delete app.go.bak (1767 lines)
- [ ] Delete fact_list.go.bak (668 lines)
- [ ] Delete burst_list.go.bak (604 lines)
- [ ] Verify no references to these files
- [ ] Commit deletion

#### 3.2 Remove/Use LoadingRotators

**Option A: Use Them** (Recommended)
```go
// In capture_event_intent.go viewSubmit()
message := i.loadingRotator.GetMessage()
s := message + "\n\n"

// Add tick handling
case LoadingTickMsg:
    if i.state == CaptureStateSubmit {
        i.loadingRotator.Rotate()
        return i, i.tickLoadingRotator()
    }
```

**Option B: Remove Them**
```go
// Remove from struct
// loadingRotator *components.LoadingMessageRotator

// Remove initialization
// loadingRotator := components.NewLoadingMessageRotator(...)
```

**Files to update**:
- `internal/cli/intents/capture_event_intent.go:84-85`
- `internal/cli/intents/browse_timeline_intent.go:61-62`
- `internal/cli/intents/generate_cv_intent.go:32-33`

**Tasks**:
- [ ] Decide: use or remove
- [ ] If using: add tick handling and view usage
- [ ] If removing: delete fields and initialization
- [ ] Test no regressions

#### 3.3 Remove Orphaned Message Types

**Location**: `internal/cli/app/messages.go:38-156`

**Orphaned messages** (no handlers exist):
- `ConfirmDeleteMsg`
- `CancelDeleteMsg`
- `EventDeletedMsg`
- `EventUpdatedMsg`
- `ApplyBulkOperationsMsg`
- `CancelBulkOperationsMsg`
- `SkipStepMsg`
- `ReviewLaterMsg`
- `FactExtractionTriggeredMsg`
- `FactsReadyMsg`
- `ConfirmFactMsg`
- `RejectFactMsg`
- `FactProcessingCompleteMsg`

**Tasks**:
- [ ] Search for usage of each message type
- [ ] Remove messages with zero usage
- [ ] Keep messages if future implementation planned
- [ ] Document decision for kept messages
- [ ] Test no compilation errors

---

### Phase 4: Stub Function Cleanup (30 min)

#### 4.1 SetInitialScreen and SetInitialCaptureMode

**Location**: `internal/cli/app/app.go:629-632, 635-638`

**Option A: Implement Them**
```go
func (m *Model) SetInitialScreen(screen Screen) {
    // Navigate to specific intent based on screen
    switch screen {
    case ListScreen:
        m.router.ActivateIntent(context.Background(), "browse_timeline")
    case "capture":
        m.router.ActivateIntent(context.Background(), "capture_event")
    // ... other screens
    }
}

func (m *Model) SetInitialCaptureMode(mode string) {
    // Store in global context for capture intent
    m.globalContext.SetPreference("initial_capture_strategy", mode)
}
```

**Option B: Remove Them** (if not used)

**Tasks**:
- [ ] Check if methods are called anywhere
- [ ] If called: implement them
- [ ] If not called: remove them
- [ ] Update CLI flags if removing

#### 4.2 AudienceRelevanceSelector.ToggleSelected()

**Location**: `internal/cli/components/audience_relevance_selector.go:39-41`

**Current (STUB)**:
```go
func (a *AudienceRelevanceSelector) ToggleSelected() {
    // This is a placeholder - implementation would depend on current focus
}
```

**Option A: Implement It**
```go
func (a *AudienceRelevanceSelector) ToggleSelected() {
    if a.focusedIndex >= 0 && a.focusedIndex < len(a.options) {
        a.selected[a.focusedIndex] = !a.selected[a.focusedIndex]
    }
}
```

**Option B: Remove It** (if not used)

**Tasks**:
- [ ] Check if method is called
- [ ] If called: implement it
- [ ] If not called: remove it

#### 4.3 CLI Flags

**Location**: `cmd/cli/main.go`

**Unused flags**:
- `--review-facts` (line 39, 82, 437)
- `--skip-import-review` (line 79-80)

**Tasks**:
- [ ] Check if flags are actually used
- [ ] If used: implement functionality
- [ ] If not used: remove flag parsing
- [ ] Update help text

---

## Acceptance Criteria

### Must Have
- [ ] '?' key shows/hides help screen
- [ ] CV generation errors are displayed to user
- [ ] All .bak files deleted
- [ ] LoadingRotators either used or removed
- [ ] Orphaned message types removed (or documented why kept)
- [ ] Stub functions either implemented or removed
- [ ] All tests passing (2,078/2,078)
- [ ] Zero staticcheck warnings
- [ ] Build successful

### Should Have
- [ ] Help screen shows all major shortcuts
- [ ] Error messages are clear and actionable
- [ ] No dead code remains
- [ ] Code is more maintainable

---

## Testing Strategy

```go
It("should show help screen on ? key", func() {
    model.Update(tea.KeyMsg{String: "?"})
    view := model.View()
    Expect(view).To(ContainSubstring("Keyboard Reference"))
})

It("should display generation errors", func() {
    intent.state.generationError = errors.New("test error")
    view := intent.View()
    Expect(view).To(ContainSubstring("test error"))
})
```

### Verification Commands
```bash
# Check for .bak files (should be none)
find . -name "*.bak"

# Check for unused loadingRotator
rg "loadingRotator" internal/cli/intents/

# Run tests
go test ./... -v

# Staticcheck
staticcheck ./...

# Compliance
make check-compliance
```

---

## Risk Assessment

### Low Risk
- All changes are cleanup or simple additions
- Easy to revert if issues found
- Well-tested existing functionality

### Mitigation
- Test after each phase
- Commit each phase separately
- Keep changes atomic

---

## References

- `internal/cli/app/app.go` - Help screen implementation
- `internal/cli/intents/generate_cv_intent.go` - Error display
- `docs/KEYBOARD_REFERENCE.md` - Help content source

---

**Last Updated**: 2026-01-08
**Status**: Ready for implementation
