# Task 38: TUI Visual Overhaul - btop-Inspired Professional Polish

## Overview
- **Goal**: Comprehensive visual overhaul with theme system, btop-inspired aesthetics, and professional polish
- **Time Estimate**: 10 days (leveraging existing bubbles components)
- **Prerequisites**: 
  - Task 20 (Coverage Improvement) - ✅ Completed
  - Understanding of Lipgloss, Bubble Tea, and KaRiya's current TUI architecture
- **Reference**: `docs/TUI_VISUAL_OVERHAUL_SPEC.md` for complete design specification

## Key Principle: Leverage Existing Libraries

**DO NOT reinvent existing components.** We already have these in `go.mod`:

| Package | Purpose | Action |
|---------|---------|--------|
| `bubbles/list` | Lists with selection | Apply theme styling |
| `bubbles/table` | Tables | Apply theme styling |
| `bubbles/progress` | Progress bars with gradients | Apply theme styling |
| `bubbles/spinner` | Loading spinners | Apply theme styling |
| `bubbles/help` | Help key bindings | Apply theme styling |
| `bubbles/viewport` | Scrollable content | Apply theme styling |
| `harmonica` | Spring animations | Already available |

**New component needed:** Only `KeyBadge` - a thin wrapper for styled keyboard shortcuts.

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - [ ] `internal/cli/styles/styles.go` (current styling)
  - [ ] `internal/cli/components/` (current components)
  - [ ] `internal/cli/intents/` (current view implementations)
  - [ ] `docs/TUI_VISUAL_OVERHAUL_SPEC.md` (design specification)
- [ ] Confirmed this is ONE atomic task (visual overhaul system)
- [ ] Identified test files that will be created/modified

## Files to Modify/Create

### New Files (Theme System)
- [ ] `internal/cli/themes/theme.go` - Theme interface and types
- [ ] `internal/cli/themes/theme_test.go` - Theme tests
- [ ] `internal/cli/themes/manager.go` - Theme manager
- [ ] `internal/cli/themes/manager_test.go` - Manager tests
- [ ] `internal/cli/themes/styles.go` - StyleSet and GenerateStyles
- [ ] `internal/cli/themes/styles_test.go` - Styles tests
- [ ] `internal/cli/themes/default.go` - Default theme (current colors)
- [ ] `internal/cli/themes/detector.go` - Terminal theme detection
- [ ] `internal/cli/themes/detector_test.go` - Detector tests

### New Files (Theme Integration & KeyBadge)
- [ ] `internal/cli/themes/bubbles.go` - Theme integration helpers for bubbles components
- [ ] `internal/cli/themes/bubbles_test.go` - Tests for theme integration helpers
- [ ] `internal/cli/components/key_badge.go` - Styled key badges (thin wrapper)
- [ ] `internal/cli/components/key_badge_test.go`

**Note:** We use existing `bubbles/*` components (list, table, progress, spinner, help)
rather than building custom versions. The `themes/bubbles.go` provides helper functions
to apply theme-aware styling to these existing components.

### Modified Files
- [ ] `internal/cli/styles/styles.go` - Refactor to use theme system
- [ ] `internal/cli/context/global.go` - Add ThemeManager to context
- [ ] `internal/cli/app/app.go` - Initialize theme system
- [ ] `internal/cli/intents/capture_event_intent.go` - Use new components
- [ ] `internal/cli/intents/browse_timeline_intent.go` - Use new components
- [ ] `internal/cli/intents/generate_cv_intent.go` - Use new components
- [ ] `internal/cli/intents/export_artifact_intent.go` - Use new components
- [ ] `internal/cli/intents/configure_system_intent.go` - Use new components
- [ ] `internal/cli/components/standard_view.go` - Theme-aware rendering
- [ ] `internal/cli/components/modal.go` - Theme-aware modals
- [ ] `internal/cli/forms/forms.go` - Unify with main theme
- [ ] `internal/config/config.go` - Theme configuration

### Documentation
- [ ] `docs/TUI_VISUAL_OVERHAUL_SPEC.md` - Already created (design spec)
- [ ] `docs/THEME_CUSTOMIZATION_GUIDE.md` - User guide for themes
- [ ] Update `docs/TUI_DEVELOPER_GUIDE.md` - Add theme system info
- [ ] Update `docs/TUI_STANDARDS.md` - New styling standards

---

## Implementation Plan

### PHASE 1: Theme Infrastructure (Days 1-3)

#### Day 1: Core Theme System

##### TDD: Theme Interface Tests
- [ ] Create test file: `internal/cli/themes/theme_test.go`
- [ ] Test Theme interface contract
- [ ] Test ColorPalette structure
- [ ] Test all required colors are defined
- [ ] Run tests and confirm they FAIL
- [ ] Commit: `test(themes): add failing tests for theme interface`

##### Implement Theme Interface
- [ ] Create `internal/cli/themes/theme.go`
- [ ] Implement `Theme` interface with all methods
- [ ] Implement `ColorPalette` struct
- [ ] Add semantic color helper methods
- [ ] Run tests and confirm they PASS
- [ ] Commit: `feat(themes): implement theme interface and types`

##### TDD: StyleSet Tests
- [ ] Add tests for StyleSet generation
- [ ] Test all style categories are generated
- [ ] Run tests and confirm they FAIL
- [ ] Commit: `test(themes): add tests for style set generation`

##### Implement StyleSet
- [ ] Create `internal/cli/themes/styles.go`
- [ ] Implement `StyleSet` struct
- [ ] Implement `GenerateStyles(palette)` function
- [ ] Run tests and confirm they PASS
- [ ] Commit: `feat(themes): implement style set generation`

#### Day 2: Default Theme & Manager

##### TDD: Default Theme Tests
- [ ] Test default theme has all required colors
- [ ] Test default theme matches current KaRiya colors
- [ ] Test default theme generates valid styles
- [ ] Commit: `test(themes): add tests for default theme`

##### Implement Default Theme
- [ ] Create `internal/cli/themes/default.go`
- [ ] Define DefaultPalette with current colors
- [ ] Implement DefaultTheme struct
- [ ] Ensure backwards compatibility
- [ ] Commit: `feat(themes): implement default theme with current colors`

##### TDD: Theme Manager Tests
- [ ] Test file: `internal/cli/themes/manager_test.go`
- [ ] Test theme registration
- [ ] Test theme switching
- [ ] Test default theme fallback
- [ ] Test OnChange callbacks
- [ ] Commit: `test(themes): add tests for theme manager`

##### Implement Theme Manager
- [ ] Create `internal/cli/themes/manager.go`
- [ ] Implement `ThemeManager` struct
- [ ] Implement `Register()`, `SetActive()`, `Active()`, `List()`
- [ ] Implement `OnChange()` callback system
- [ ] Commit: `feat(themes): implement theme manager`

#### Day 3: Terminal Detection & App Integration

##### TDD: Terminal Detection Tests
- [ ] Test file: `internal/cli/themes/detector_test.go`
- [ ] Test color depth detection
- [ ] Test dark/light mode detection
- [ ] Test with various TERM values
- [ ] Commit: `test(themes): add tests for terminal detection`

##### Implement Terminal Detection
- [ ] Create `internal/cli/themes/detector.go`
- [ ] Implement `DetectColorDepth()`
- [ ] Implement `DetectDarkMode()`
- [ ] Implement `AutoSelect()` for ThemeManager
- [ ] Commit: `feat(themes): implement terminal theme detection`

##### Wire Theme System into App
- [ ] Update `internal/cli/context/global.go` - Add ThemeManager field
- [ ] Update `internal/cli/app/app.go` - Initialize ThemeManager
- [ ] Update `internal/config/config.go` - Add theme config field
- [ ] Run all existing tests to verify no regressions
- [ ] Commit: `feat(app): wire theme system into application`

##### Compliance Check
- [ ] Run `make check-compliance`
- [ ] Fix any issues
- [ ] Commit any fixes

---

### PHASE 2: Theme Integration & KeyBadge (Days 4-5)

**Note:** We use existing `bubbles/*` components instead of building custom ones.
This phase focuses on creating theme integration helpers and the only new component (KeyBadge).

#### Day 4: Bubbles Theme Integration Helpers

##### TDD: Theme Integration Helper Tests
- [ ] Test file: `internal/cli/themes/bubbles_test.go`
- [ ] Test `NewThemedList()` applies theme to bubbles/list
- [ ] Test `NewThemedTable()` applies theme to bubbles/table
- [ ] Test `NewThemedProgress()` applies theme to bubbles/progress
- [ ] Test `NewThemedSpinner()` applies theme to bubbles/spinner
- [ ] Test `NewThemedHelp()` applies theme to bubbles/help
- [ ] Commit: `test(themes): add tests for bubbles theme integration`

##### Implement Theme Integration Helpers
- [ ] Create `internal/cli/themes/bubbles.go`
- [ ] Implement `NewThemedListStyles(theme)` - returns `list.Styles`
- [ ] Implement `NewThemedTableStyles(theme)` - returns `table.Styles`
- [ ] Implement `NewThemedProgress(theme)` - returns configured `progress.Model`
- [ ] Implement `NewThemedSpinner(theme)` - returns configured `spinner.Model`
- [ ] Implement `NewThemedHelpStyles(theme)` - returns `help.Styles`
- [ ] Commit: `feat(themes): implement bubbles theme integration helpers`

#### Day 5: KeyBadge Component

##### TDD: KeyBadge Tests
- [ ] Test file: `internal/cli/components/key_badge_test.go`
- [ ] Test single badge rendering
- [ ] Test badge styling with theme
- [ ] Test help footer composition
- [ ] Test multiple badges in footer
- [ ] Commit: `test(components): add tests for key badges`

##### Implement KeyBadge
- [ ] Create `internal/cli/components/key_badge.go`
- [ ] Implement `KeyBadge` struct (simple, not a full Bubble Tea model)
- [ ] Implement `Render(theme)` method
- [ ] Implement `RenderHelpFooter()` helper
- [ ] Commit: `feat(components): implement styled key badges`

##### Update StandardView and Intents
- [ ] Update StandardView footer to use key badges
- [ ] Update all intent footers
- [ ] Commit: `refactor(intents): use key badges in all footers`

##### Compliance Check
- [ ] Run `make check-compliance`
- [ ] Fix any issues

---

### PHASE 3: Glamour Integration & Polish (Days 6-7)

**Note:** `harmonica` is already available in `go.mod`. We only need to add `glamour` for markdown rendering.

#### Day 6: Add Glamour & Update Loading States

##### Add Glamour Dependency
- [ ] Run: `go get github.com/charmbracelet/glamour`
- [ ] Update `go.mod` and `go.sum`
- [ ] Commit: `deps: add glamour for markdown rendering`

##### Update Loading States with Themed Spinners
- [ ] Update modal loading spinner to use `NewThemedSpinner()`
- [ ] Verify existing `LoadingMessageRotator` works with theme
- [ ] Apply theme colors throughout
- [ ] Commit: `refactor(components): use themed spinners in loading states`

#### Day 7: Markdown Rendering (Glamour)

##### TDD: CV Preview Tests
- [ ] Test markdown rendering
- [ ] Test theme application to glamour
- [ ] Commit: `test(generate-cv): add tests for markdown preview`

##### Implement CV Markdown Preview
- [ ] Create `GenerateGlamourStyle(theme)` helper
- [ ] Update GenerateCV intent to use Glamour for CV preview
- [ ] Apply theme-matched styling
- [ ] Commit: `feat(generate-cv): use glamour for cv preview`

##### Compliance Check
- [ ] Run `make check-compliance`
- [ ] Fix any issues

---

### PHASE 4: Intent Updates (Days 8-9)

#### Day 8: CaptureEvent & BrowseTimeline

##### Update CaptureEvent Intent
- [ ] Use theme from context throughout
- [ ] Apply new components (SelectionList, KeyBadge, etc.)
- [ ] Remove all hardcoded `styles.Color*` references
- [ ] Run existing tests
- [ ] Commit: `refactor(capture-event): migrate to theme system`

##### Update BrowseTimeline Intent
- [ ] Use theme from context throughout
- [ ] Apply themed `bubbles/table` styles via `NewThemedTableStyles()`
- [ ] Apply key badges
- [ ] Remove all hardcoded styles
- [ ] Run existing tests
- [ ] Commit: `refactor(browse-timeline): migrate to theme system`

#### Day 9: GenerateCV, ExportArtifact & ConfigureSystem

##### Update GenerateCV Intent
- [ ] Use theme from context
- [ ] Apply themed `bubbles/list` for profile/audience via `NewThemedListStyles()`
- [ ] Apply themed `bubbles/progress` for generation via `NewThemedProgress()`
- [ ] Apply Glamour for preview
- [ ] Commit: `refactor(generate-cv): migrate to theme system`

##### Update ExportArtifact Intent
- [ ] Use theme from context
- [ ] Apply themed `bubbles/progress` for export
- [ ] Update all views to use theme
- [ ] Commit: `refactor(export-artifact): migrate to theme system`

##### Update ConfigureSystem Intent
- [ ] Use theme from context
- [ ] Add theme selection option in settings
- [ ] Commit: `refactor(configure-system): migrate to theme system`

##### Unify Huh Forms Theme
- [ ] Update `internal/cli/forms/forms.go`
- [ ] Create `GenerateHuhTheme(theme Theme)` function
- [ ] Ensure Huh forms match main UI
- [ ] Commit: `feat(forms): unify huh theme with main ui theme`

##### Compliance Check
- [ ] Run `make check-compliance`
- [ ] Fix any issues

---

### PHASE 5: Testing & Documentation (Day 10)

#### Day 10: Comprehensive Testing & Documentation

##### Visual Testing
- [ ] Test default theme in all intents
- [ ] Verify consistent styling
- [ ] Check all screens render correctly

##### Responsive Testing
- [ ] Test 60 column terminal
- [ ] Test 80 column terminal
- [ ] Test 120 column terminal
- [ ] Test 200 column terminal
- [ ] Test minimum height (24 lines)

##### Terminal Mode Testing
- [ ] Test with TERM=xterm (basic)
- [ ] Test with TERM=xterm-256color
- [ ] Test with COLORTERM=truecolor

##### Integration Testing
- [ ] Test complete capture workflow
- [ ] Test complete browse workflow
- [ ] Test complete CV generation workflow
- [ ] Test complete export workflow

##### Create Theme Customization Guide
- [ ] Create `docs/THEME_CUSTOMIZATION_GUIDE.md`
- [ ] Document how to switch themes
- [ ] Document theme configuration
- [ ] Commit: `docs: add theme customization guide`

##### Update Developer Documentation
- [ ] Update `docs/TUI_DEVELOPER_GUIDE.md`
- [ ] Add section on theme system
- [ ] Add examples of using themes in components
- [ ] Commit: `docs: update tui developer guide for themes`

##### Update TUI Standards
- [ ] Update `docs/TUI_STANDARDS.md`
- [ ] Add new styling standards
- [ ] Add component usage guidelines
- [ ] Commit: `docs: update tui standards for new components`

##### Final Compliance Check
- [ ] Run `make check-compliance`
- [ ] Run full test suite: `make test`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Verify all tests pass
- [ ] Fix any remaining issues

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make review-commit` passes
- [ ] AI attribution included (AI-generated code)
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)
- [ ] Tests pass for committed code

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] Theme system fully functional with runtime switching
- [ ] Terminal theme auto-detection working
- [ ] All hardcoded `styles.Color*` references replaced with theme calls
- [ ] Themed bubbles/list used for selections (via `NewThemedListStyles()`)
- [ ] Themed bubbles/table used for data display (via `NewThemedTableStyles()`)
- [ ] Key badges used in all footers
- [ ] Themed bubbles/progress bars in loading states (via `NewThemedProgress()`)
- [ ] Glamour used for CV preview rendering
- [ ] Huh forms unified with main theme
- [ ] All existing tests pass
- [ ] Code coverage maintained ≥ 80%
- [ ] Documentation complete
- [ ] Visual quality matches professional TUI apps

## Rollback Plan
1. All work is in feature branches
2. Each phase can be rolled back independently
3. Theme system is additive (doesn't break existing code)
4. Default theme preserves current appearance
5. If needed, revert commits in reverse order
6. Run `make check-compliance` after rollback

## Notes
- **Leverages existing bubbles/* components** - no need to build custom list/table/progress/spinner
- Theme system is backwards compatible (default theme = current colors)
- Only new component is KeyBadge (thin wrapper for styled keys)
- Each phase builds on previous phases
- Can pause after any phase if needed
- Reference `docs/TUI_VISUAL_OVERHAUL_SPEC.md` for detailed specifications
- Additional themes can be added incrementally after Phase 1
- Time reduced from 15 to 10 days by reusing existing components

## Success Metrics
- [ ] KaRiya looks as polished as btop
- [ ] Users can switch themes at runtime
- [ ] Terminal settings auto-detected
- [ ] No performance regression
- [ ] All tests pass
- [ ] Zero staticcheck warnings
- [ ] Developer experience improved (semantic styles vs hardcoded colors)
