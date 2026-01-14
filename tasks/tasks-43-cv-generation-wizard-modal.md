---
created: 2026-01-14T12:00
modified: 2026-01-14T12:00
---
# Task 43: CV Generation - Wizard Modal Architecture

## Overview
- **Goal**: Refactor CV Generation from 17-state screen architecture to wizard modal + 2 screens
- **Time Estimate**: 9-10 hours
- **Prerequisites**: Task 42 Phase 4.2 complete (BrowseTimeline reference implementation)
- **Reference**: BrowseTimeline modal-heavy pattern (5 modals + 2 screens)

## ⚠️ CRITICAL REQUIREMENTS

### State Matrix Integration (MANDATORY)
After EVERY component creation or deletion:
```bash
make generate-diagrams
```
Verify STATE_MATRIX.md updated correctly before committing.

### TUI Standards Compliance (MANDATORY)
All components MUST follow: `docs/TUI_STANDARDS.md`
- ✅ Universal keyboard shortcuts (esc, ↑/↓, j/k, enter)
- ✅ Escape key behavior matches state type
- ✅ StandardView component with logo and breadcrumbs (screens only)
- ✅ Help text visible in footer
- ✅ Theme system integration

**Quick Check Commands**:
```bash
make check-compliance    # Run before EVERY commit
make generate-diagrams   # Update state matrix after component changes
make test                # Verify all tests pass
```

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

---

## Architecture Decision Summary

### Current State (Before)
- **States**: 17 separate states
- **Architecture**: Screen per state approach
- **Intent size**: 2,019 lines
- **Screens exist**: 4 (profile_select, audience_select, generating, preview)
- **Problem**: Only 24% coverage, 13 missing screens blocking migration

### Target State (After)
- **States**: 5 states (70% reduction)
- **Architecture**: Wizard modal + progress modals + preview screen
- **Intent size**: ~400 lines (80% reduction)
- **Components**: 4 total (1 wizard + 1 progress + 1 export + 1 screen)
- **Benefit**: Modal-heavy pattern proven by BrowseTimeline

---

## Design Rationale

### Modal vs Screen Decision Matrix

**User's Rules Applied**:
1. **Error** → Modal
2. **Exporting** → Modal (async call)
3. **Select*** → Modal (anything with "Select" in name)
4. **Generating** → Modal (async call)
5. **Preview** → Screen (viewport)
6. **Review** → Screen (viewport) - **ELIMINATED (not required)**

**Result**: 85% modal-based workflow

### Wizard Modal Approach (Option A - User's Choice)

**Why Wizard?**
- ✅ **Single modal** for all configuration steps
- ✅ **Guided workflow** - users progress through logical sections
- ✅ **Skippable steps** - sensible defaults for power users
- ✅ **State consolidation** - 8 selection states → 1 wizard state
- ✅ **Better UX** - no context switching between selections

**Alternative Rejected**: Separate modal per selection (too many modal opens/closes)

---

## Component Architecture

### Final Component List: 4 Components

| Component | Type | Purpose | Lines | Time |
|-----------|------|---------|-------|------|
| `CVConfigWizardModal` | Modal | Multi-step config (3 steps) | 300 | 2.5h |
| `CVProgressModal` | Modal | Generic async progress | 100 | 1h |
| `ExportOptionsModal` | Modal | Export format + location | 120 | 45m |
| `CVPreviewScreen` | Screen | Viewport for CV preview | 150 | 1.5h |

**Total New Code**: ~670 lines  
**Code Removed**: ~1,200 lines (legacy state handlers)  
**Net Change**: **-530 lines** (26% reduction)

---

## Wizard Modal Design

### CVConfigWizardModal Structure

**3 Steps** (all skippable with defaults):

```
Step 1: WHO          Step 2: TECH           Step 3: FORMAT
┌──────────────┐     ┌──────────────┐       ┌──────────────┐
│ • Profile    │ ──► │ • Tech Focus │ ──►   │ • Skills Fmt │
│ • Audience   │     │ • Technologies│       │ • CV Length  │
└──────────────┘     │ • Focus Area │       └──────────────┘
                     └──────────────┘
                     
                     (Step 2 conditional:
                      only if techs extracted)
```

**Navigation**:
- `Tab`: Next field
- `Enter`: Submit field / Next step
- `Esc`: Previous step (or cancel if step 1)
- `Ctrl+Enter` or `s`: Skip to generate (use defaults)

**Implementation Pattern**:
- 3 `huh.Group`s (one per step)
- huh.Form handles group progression automatically
- Back navigation via Esc
- Step 2 conditionally shown based on tech extraction results

### Data Structure

```go
type CVConfigWizardModal struct {
    form        *huh.Form
    data        *CVConfigData
    currentStep int      // 0, 1, 2
    visible     bool
    completed   bool
    skipped     bool     // User pressed skip shortcut
    
    // Extracted technologies (from async call)
    extractedTechs []*ExtractedTechnology
    techsAvailable bool
}

type CVConfigData struct {
    // Step 1: WHO
    ProfileID string   // Required
    Audience  string   // "hiring_manager" | "recruiter" | "peer"
    
    // Step 2: TECH (optional - only if techs extracted)
    TechFocus    string   // "language_agnostic" | "generalist" | "specialist"
    Technologies []string // MultiSelect (if generalist/specialist)
    FocusArea    string   // "backend" | "frontend" | "fullstack" | "devops"
    
    // Step 3: FORMAT
    SkillsFormat string   // "grouped" | "flat" | "categorized"
    CVLength     string   // "1_page" | "2_page" | "detailed"
}
```

---

## Workflow Design

### Complete User Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  1. [CVConfigWizardModal] - Configuration                        │
│      Step 1: WHO (Profile + Audience)                           │
│      Step 2: TECH (Focus + Techs + Area) - conditional          │
│      Step 3: FORMAT (Skills + Length)                           │
│          ↓ (on complete or skip)                                │
│                                                                 │
│  2. [CVProgressModal] - "Extracting technologies..."            │
│      • Background async call                                    │
│      • Esc: Cancel (return to wizard)                           │
│          ↓                                                      │
│                                                                 │
│  3. [CVProgressModal] - "Generating CV..."                      │
│      • Uses wizard config + extracted techs                     │
│      • Esc: Cancel (return to wizard)                           │
│          ↓                                                      │
│                                                                 │
│  4. [CVPreviewScreen] - Scrollable CV Preview                   │
│      • ↑/↓/j/k: Scroll through CV                               │
│      • Enter: Complete (return to menu)                         │
│      • x: Export (opens export modal)                           │
│      • Esc: Back to wizard (re-configure)                       │
│          ↓ (if x pressed)                                       │
│                                                                 │
│  5. [ExportOptionsModal] - Export Configuration                 │
│      • Format: Text / Markdown / YAML                           │
│      • Location: File / Clipboard                               │
│      • Enter: Export                                            │
│      • Esc: Cancel export                                       │
│          ↓                                                      │
│                                                                 │
│  6. [CVProgressModal] - "Exporting..."                          │
│      • Writes file or copies to clipboard                       │
│          ↓                                                      │
│                                                                 │
│  7. [Success Toast] - "Exported to /path/to/file.md"           │
│      • Auto-dismiss or Enter/Esc                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## State Machine Simplification

### Before: 17 States

```go
const (
    GenerateCVStateSelectProfile            // Screen
    GenerateCVStateSelectAudience           // Screen
    GenerateCVStateExtractingTechnologies   // Progress
    GenerateCVStateSelectTechnologyFocus    // Selection
    GenerateCVStateSelectTechnologies       // Multi-select
    GenerateCVStateSelectFocusArea          // Selection
    GenerateCVStateSelectSkillsConfig       // Form
    GenerateCVStateSelectLengthFormat       // Selection
    GenerateCVStateGenerating               // Progress
    GenerateCVStatePreview                  // Screen
    GenerateCVStateReview                   // Screen
    GenerateCVStateConfirm                  // Confirm
    GenerateCVStateExportSelectFormat       // Selection
    GenerateCVStateExportSelectLocation     // Selection
    GenerateCVStateExporting                // Progress
    GenerateCVStateExportComplete           // Success
    // + Error states
)
```

### After: 5 States (70% reduction!)

```go
const (
    CVStateConfiguring    // Wizard modal active (replaces 8 select states)
    CVStateExtracting     // Progress modal (tech extraction)
    CVStateGenerating     // Progress modal (CV generation)
    CVStatePreview        // Preview screen
    CVStateExporting      // Progress modal (export) - optional
)
```

**Eliminated States** (12):
- All "Select*" states → wizard steps (not separate states)
- Review → eliminated (not required)
- Confirm → handled by wizard completion
- ExportSelectFormat/Location → export modal (not separate states)
- ExportComplete → success toast

---

## Implementation Plan

### Phase 1: Create CVConfigWizardModal (2.5 hours)

**Files to Create**:
- `internal/cli/components/cv_config_wizard_modal.go` (300 lines)
- `internal/cli/components/cv_config_wizard_modal_test.go` (150 lines)

**TDD Checklist**:
- [ ] RED: Test wizard creation with 3 steps
- [ ] RED: Test step navigation (Tab, Enter, Esc)
- [ ] RED: Test skip shortcut (Ctrl+Enter)
- [ ] RED: Test conditional Step 2 (tech step)
- [ ] RED: Test data extraction after completion
- [ ] GREEN: Implement CVConfigWizardModal
- [ ] REFACTOR: Extract step builders

**Acceptance Criteria**:
- [ ] 3 huh.Groups (WHO, TECH, FORMAT)
- [ ] Tab/Enter navigation works
- [ ] Esc goes back a step
- [ ] Skip shortcut uses defaults
- [ ] Step 2 conditionally shown
- [ ] Validates required fields (Profile)
- [ ] Theme integration (Catppuccin)
- [ ] KeyBadge footer

**Key Features**:
```go
// Step 1: WHO
huh.NewSelect[string]().
    Key("profile").
    Title("Select CV Profile").
    Options(...).
    Value(&data.ProfileID)

huh.NewSelect[string]().
    Key("audience").
    Title("Target Audience").
    Options(
        huh.NewOption("Hiring Manager", "hiring_manager"),
        huh.NewOption("Recruiter", "recruiter"),
        huh.NewOption("Peer", "peer"),
    ).
    Value(&data.Audience)
    
// Step 2: TECH (conditional)
huh.NewSelect[string]().
    Key("tech_focus").
    Title("Technology Focus").
    Options(
        huh.NewOption("Language Agnostic", "language_agnostic"),
        huh.NewOption("Generalist", "generalist"),
        huh.NewOption("Specialist", "specialist"),
    ).
    Value(&data.TechFocus)

huh.NewMultiSelect[string]().
    Key("technologies").
    Title("Select Technologies").
    Options(...).  // From extractedTechs
    Value(&data.Technologies)

// Step 3: FORMAT
huh.NewSelect[string]().
    Key("cv_length").
    Title("CV Length").
    Options(
        huh.NewOption("1 Page", "1_page"),
        huh.NewOption("2 Pages", "2_page"),
        huh.NewOption("Detailed", "detailed"),
    ).
    Value(&data.CVLength)
```

**Commit**:
- Test: `test(components): add CVConfigWizardModal tests (RED phase)`
- Impl: `feat(components): add CVConfigWizardModal with 3-step workflow (GREEN phase)`

---

### Phase 2: Create CVProgressModal (1 hour)

**Files to Create**:
- `internal/cli/components/cv_progress_modal.go` (100 lines)
- `internal/cli/components/cv_progress_modal_test.go` (80 lines)

**TDD Checklist**:
- [ ] RED: Test modal creation with title/subtitle
- [ ] RED: Test spinner animation
- [ ] RED: Test cancellable vs non-cancellable
- [ ] RED: Test completion handling
- [ ] RED: Test error handling
- [ ] GREEN: Implement CVProgressModal
- [ ] REFACTOR: Extract spinner animation

**Acceptance Criteria**:
- [ ] Animated spinner (10-frame animation)
- [ ] Configurable title/subtitle
- [ ] Cancellable flag (Esc enabled/disabled)
- [ ] Completion/error state handling
- [ ] Theme integration
- [ ] KeyBadge footer
- [ ] Solid background (no transparency)

**Key Features**:
```go
type CVProgressModal struct {
    title       string
    subtitle    string
    spinner     int
    visible     bool
    cancellable bool
    completed   bool
    error       error
    theme       themes.Theme
}

// Factory functions
func NewExtractingTechsProgress() *CVProgressModal {
    return &CVProgressModal{
        title:       "Extracting Technologies",
        subtitle:    "Analyzing your skills...",
        cancellable: true,
        visible:     true,
    }
}

func NewGeneratingCVProgress(profile, audience string) *CVProgressModal {
    return &CVProgressModal{
        title:       "Generating CV",
        subtitle:    fmt.Sprintf("Profile: %s | Audience: %s", profile, audience),
        cancellable: true,
        visible:     true,
    }
}

func NewExportingProgress(format string) *CVProgressModal {
    return &CVProgressModal{
        title:       "Exporting CV",
        subtitle:    fmt.Sprintf("Format: %s", format),
        cancellable: false,
        visible:     true,
    }
}
```

**Commit**:
- Test: `test(components): add CVProgressModal tests (RED phase)`
- Impl: `feat(components): add CVProgressModal for async operations (GREEN phase)`

---

### Phase 3: Create ExportOptionsModal (45 minutes)

**Files to Create**:
- `internal/cli/components/export_options_modal.go` (120 lines)
- `internal/cli/components/export_options_modal_test.go` (60 lines)

**TDD Checklist**:
- [ ] RED: Test modal creation
- [ ] RED: Test format selection
- [ ] RED: Test location selection
- [ ] RED: Test submission
- [ ] RED: Test cancellation
- [ ] GREEN: Implement ExportOptionsModal
- [ ] REFACTOR: Extract form builders

**Acceptance Criteria**:
- [ ] 2 fields: Format + Location
- [ ] Format options: Text, Markdown, YAML
- [ ] Location options: File, Clipboard
- [ ] Enter to submit
- [ ] Esc to cancel
- [ ] Theme integration
- [ ] KeyBadge footer

**Key Features**:
```go
type ExportOptionsModal struct {
    form     *huh.Form
    data     *ExportData
    visible  bool
}

type ExportData struct {
    Format   string  // "text" | "markdown" | "yaml"
    Location string  // "file" | "clipboard"
}

// Simple 2-field form
huh.NewSelect[string]().
    Key("format").
    Title("Export Format").
    Options(
        huh.NewOption("Plain Text", "text"),
        huh.NewOption("Markdown", "markdown"),
        huh.NewOption("YAML", "yaml"),
    ).
    Value(&data.Format)

huh.NewSelect[string]().
    Key("location").
    Title("Save To").
    Options(
        huh.NewOption("File", "file"),
        huh.NewOption("Clipboard", "clipboard"),
    ).
    Value(&data.Location)
```

**Commit**:
- Test: `test(components): add ExportOptionsModal tests (RED phase)`
- Impl: `feat(components): add ExportOptionsModal for export config (GREEN phase)`

---

### Phase 4: Refactor CVPreviewScreen (1.5 hours)

**Files to Modify**:
- `internal/cli/screens/cv/preview.go` (refactor to 150 lines)
- `internal/cli/screens/cv/preview_test.go` (create, 100 lines)

**TDD Checklist**:
- [ ] RED: Test viewport scrolling
- [ ] RED: Test keyboard navigation (↑/↓/j/k/g/G)
- [ ] RED: Test actions (Enter, x, Esc)
- [ ] RED: Test CV rendering
- [ ] GREEN: Refactor CVPreviewScreen with viewport
- [ ] REFACTOR: Extract CV formatter

**Current Issues** (from analysis):
1. ❌ No proper scrolling - scrollOffset not used with viewport
2. ❌ Plain text footer - not using KeyBadge components
3. ❌ Multiple actions mixed - needs clear separation

**Acceptance Criteria**:
- [ ] Uses bubbles `viewport.Model` for scrolling
- [ ] Keyboard navigation works (↑/↓/j/k/g/G)
- [ ] Enter: Complete workflow
- [ ] x: Open export modal
- [ ] Esc: Back to wizard
- [ ] KeyBadge footer
- [ ] Theme integration
- [ ] CV content formatted properly

**Implementation Pattern** (use BaseDetailScreen pattern):
```go
type CVPreviewScreen struct {
    *base.BaseScreen
    cv       *career.CVView
    viewport viewport.Model
}

// Keyboard handling
case "enter", "y":
    return nil, &screens.SubmitResult{
        ResultData: "complete",
    }
case "x":
    return nil, &screens.NavigateResult{
        Action: "export",
    }
case "esc":
    return nil, &screens.CancelResult{}

// Footer using KeyBadge
components.ScrollBadge(),           // ↑↓/jk: Scroll
components.NewKeyBadge("g/G", "Top/Bottom"),
components.NewKeyBadge("Enter", "Complete"),
components.NewKeyBadge("x", "Export"),
components.BackBadge(),             // Esc: Back
```

**Files to Delete**:
- `internal/cli/screens/cv/audience_select.go` (replaced by wizard)
- `internal/cli/screens/cv/profile_select.go` (replaced by wizard)
- `internal/cli/screens/cv/profile_select_test.go`

**Commit**:
- Test: `test(screens): add CVPreviewScreen viewport tests (RED phase)`
- Impl: `refactor(screens): use viewport in CVPreviewScreen with KeyBadge footer (GREEN phase)`
- Cleanup: `refactor(screens): remove profile/audience screens (replaced by wizard)`

---

### Phase 5: Intent Integration (2 hours)

**Files to Modify**:
- `internal/cli/intents/generate_cv_intent.go` (massive simplification)
- `internal/cli/intents/generate_cv.go` (update states)
- `internal/cli/intents/generate_cv_test.go` (update tests)

**Changes Required**:

#### 5.1 Update State Constants (10 min)

```go
// Before: 17 states
const (
    GenerateCVStateSelectProfile
    GenerateCVStateSelectAudience
    GenerateCVStateExtractingTechnologies
    GenerateCVStateSelectTechnologyFocus
    GenerateCVStateSelectTechnologies
    GenerateCVStateSelectFocusArea
    GenerateCVStateSelectSkillsConfig
    GenerateCVStateSelectLengthFormat
    GenerateCVStateGenerating
    GenerateCVStatePreview
    GenerateCVStateReview
    GenerateCVStateConfirm
    GenerateCVStateExportSelectFormat
    GenerateCVStateExportSelectLocation
    GenerateCVStateExporting
    GenerateCVStateExportComplete
)

// After: 5 states (70% reduction!)
const (
    CVStateConfiguring    GenerateCVState = "configuring"
    CVStateExtracting     GenerateCVState = "extracting"
    CVStateGenerating     GenerateCVState = "generating"
    CVStatePreview        GenerateCVState = "preview"
    CVStateExporting      GenerateCVState = "exporting"
)
```

#### 5.2 Simplify Intent Structure (20 min)

```go
type GenerateCVIntent struct {
    *BaseIntent
    
    context *GenerateCVContext
    state   *GenerateCVModel
    active  bool
    result  *IntentResult[*GenerateCVResult]
    
    // Modals
    wizardModal       *components.CVConfigWizardModal
    progressModal     *components.CVProgressModal
    exportModal       *components.ExportOptionsModal
    
    // Screens
    previewScreen     *cv.CVPreviewScreen
    
    // State
    currentState      GenerateCVState
    cvGenerated       *career.CVView
}
```

#### 5.3 Implement 3-Tier Update() (30 min)

```go
func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Tier 1: Modal updates (HIGHEST PRIORITY)
        if i.wizardModal != nil && i.wizardModal.IsVisible() {
            return i.handleWizardModalUpdate(msg)
        }
        if i.exportModal != nil && i.exportModal.IsVisible() {
            return i.handleExportModalUpdate(msg)
        }
        if i.progressModal != nil && i.progressModal.IsVisible() {
            return i.handleProgressModalUpdate(msg)
        }
        
        // Tier 2: Global keys
        switch msg.String() {
        case "ctrl+c", "q":
            return i.handleQuit()
        case "?", "h":
            return i.handleHelp()
        case "m":
            return i.handleMainMenu()
        }
        
        // Tier 3: Screen delegation
        if i.previewScreen != nil {
            return i.handlePreviewScreenUpdate(msg)
        }
        
    case tea.WindowSizeMsg:
        i.SetTerminalInfo(msg.Width, msg.Height)
        // Update modal/screen dimensions
        
    case WizardCompleteMsg:
        return i.handleWizardComplete(msg)
        
    case TechnologiesExtractedMsg:
        return i.handleTechsExtracted(msg)
        
    case CVGenerationCompleteMsg:
        return i.handleCVGenerated(msg)
        
    case ExportCompleteMsg:
        return i.handleExportComplete(msg)
    }
    return nil
}
```

#### 5.4 Implement View() with Modal Overlays (20 min)

```go
func (i *GenerateCVIntent) View() string {
    // Create StandardView
    view := i.CreateViewWithBreadcrumbs(
        i.getStateName(),
        i.terminal.Width,
        i.terminal.Height,
    )
    
    // Render content based on state
    var content string
    switch i.currentState {
    case CVStatePreview:
        if i.previewScreen != nil {
            content = i.previewScreen.RenderContent()
        }
    default:
        content = "Processing..."
    }
    
    view.WithContent(content)
    view.WithHelp(i.getContextHelp())
    baseView := view.Render()
    
    // Overlay visible modal (LAST STEP)
    if i.wizardModal != nil && i.wizardModal.IsVisible() {
        return i.renderWizardModalOverlay(baseView)
    }
    if i.progressModal != nil && i.progressModal.IsVisible() {
        return i.renderProgressModalOverlay(baseView)
    }
    if i.exportModal != nil && i.exportModal.IsVisible() {
        return i.renderExportModalOverlay(baseView)
    }
    
    return baseView
}
```

#### 5.5 Remove Legacy Code (~1,200 lines)

**Delete These Methods**:
- `updateSelectProfile()`
- `updateSelectAudience()`
- `updateExtractingTechnologies()`
- `updateSelectTechnologyFocus()`
- `updateSelectTechnologies()`
- `updateSelectFocusArea()`
- `updateSelectSkillsConfig()`
- `updateSelectLengthFormat()`
- `updateGenerating()`
- `updatePreview()`
- `updateReview()`
- `updateConfirm()`
- `updateExportSelectFormat()`
- `updateExportSelectLocation()`
- `updateExporting()`
- `updateExportComplete()`
- `viewSelectProfile()`
- `viewSelectAudience()`
- `viewExtractingTechnologies()`
- `viewSelectTechnologyFocus()`
- `viewSelectTechnologies()`
- `viewSelectFocusArea()`
- `viewSelectSkillsConfig()`
- `viewSelectLengthFormat()`
- `viewGenerating()`
- `viewPreview()`
- `viewReview()`
- `viewConfirm()`
- `viewExportSelectFormat()`
- `viewExportSelectLocation()`
- `viewExporting()`
- `viewExportComplete()`

**Estimated Removal**: ~1,200 lines

**Commit**:
- `refactor(intents): simplify GenerateCV to 5 states with wizard modal`
- `refactor(intents): implement 3-tier key handling in GenerateCV`
- `refactor(intents): remove legacy state handlers (1,200 lines)`

---

### Phase 6: Tests and Documentation (1.5 hours)

#### 6.1 Update Intent Tests (1 hour)

**Files to Modify**:
- `internal/cli/intents/generate_cv_test.go`
- `internal/cli/intents/generate_cv_technology_test.go`
- `internal/cli/intents/generate_cv_escape_test.go`

**Test Updates Required**:
- [ ] Update state expectations (17 → 5 states)
- [ ] Add wizard modal tests
- [ ] Add progress modal tests
- [ ] Update screen tests (preview only)
- [ ] Remove tests for eliminated states
- [ ] Add E2E wizard workflow test

**E2E Test Example**:
```go
It("should complete full CV generation workflow", func() {
    // 1. Init opens wizard
    Expect(intent.wizardModal.IsVisible()).To(BeTrue())
    
    // 2. Complete wizard
    intent.wizardModal.data.ProfileID = "profile-1"
    intent.wizardModal.data.Audience = "hiring_manager"
    intent.wizardModal.Complete()
    
    // 3. Tech extraction starts
    Expect(intent.currentState).To(Equal(CVStateExtracting))
    
    // 4. Simulate extraction complete
    intent.Update(TechnologiesExtractedMsg{Technologies: techs})
    
    // 5. CV generation starts
    Expect(intent.currentState).To(Equal(CVStateGenerating))
    
    // 6. Simulate generation complete
    intent.Update(CVGenerationCompleteMsg{CV: cv})
    
    // 7. Preview screen shown
    Expect(intent.currentState).To(Equal(CVStatePreview))
    Expect(intent.previewScreen).NotTo(BeNil())
})
```

#### 6.2 Update Documentation (30 min)

**Files to Update**:
- [ ] `docs/workflows/CV_GENERATION_WORKFLOW.md` - Update state diagram and steps
- [ ] `docs/KEYBOARD_SHORTCUTS_GUIDE.md` - Update CV generation shortcuts
- [ ] `AGENTS.md` - Update Task 43 completion status
- [ ] `CV_GENERATION_SCREEN_ANALYSIS.md` - Mark as obsolete, link to Task 43

**Documentation Updates**:
- Update state machine diagram (17 → 5 states)
- Document wizard modal steps
- Update keyboard shortcuts
- Add modal overlay patterns used
- Update workflow timing estimates

**Commit**:
- `test(intents): update GenerateCV tests for wizard modal architecture`
- `docs: update CV generation workflow documentation`

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)

- [ ] `make check-compliance` passes (REQUIRED)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)
- [ ] All tests pass in affected areas

---

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)

- [ ] `make check-compliance` passes
- [ ] All 6 phases complete
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

---

## Acceptance Criteria

### Functionality
- [ ] Wizard modal with 3 steps works
- [ ] Step skipping with defaults functional
- [ ] Tech extraction progress shown
- [ ] CV generation progress shown
- [ ] Preview screen with viewport scrolling
- [ ] Export from preview works
- [ ] All keyboard shortcuts functional

### Code Quality
- [ ] Intent reduced from 2,019 → ~400 lines (80% reduction)
- [ ] States reduced from 17 → 5 (70% reduction)
- [ ] All tests pass (>95% pass rate)
- [ ] Coverage maintained ≥ 80%
- [ ] Zero race conditions
- [ ] Zero staticcheck warnings

### Patterns Compliance
- [ ] Pattern 1: Modal Overlay Rendering (StandardView first, modal last)
- [ ] Pattern 2: Themed Footer Building (KeyBadge components only)
- [ ] Pattern 3: View Rendering with Modal Overlay
- [ ] Pattern 4: Global Key Interception (modal → global → screen)
- [ ] Pattern 5: Context-Aware Footer Generation
- [ ] Pattern 12: Form Modal with Immediate Init

### Documentation
- [ ] STATE_MATRIX.md updated (`make generate-diagrams`)
- [ ] Workflow guide updated
- [ ] Keyboard shortcuts documented
- [ ] Task marked complete

---

## Rollback Plan

If implementation fails or introduces regressions:

1. **Revert to hybrid approach** - Keep `useScreens = false`
2. **Incremental rollback** - Remove phases in reverse order
3. **Safety net** - All legacy code kept in commits until final cleanup

**Recovery Steps**:
```bash
# Revert specific phase
git revert <commit-hash>

# Or full task revert
git revert <first-commit>..<last-commit>

# Verify tests pass
make test
```

---

## Success Metrics

### Primary Goals
- ✅ **80% code reduction** - 2,019 → ~400 lines
- ✅ **70% state reduction** - 17 → 5 states
- ✅ **Modal-heavy pattern** - 3 modals + 1 screen
- ✅ **Wizard UX** - Single modal for all config

### Secondary Goals
- ✅ **Test pass rate** - Maintain >95%
- ✅ **Zero regressions** - All workflows functional
- ✅ **Pattern compliance** - Follow BrowseTimeline patterns
- ✅ **Performance** - No degradation vs legacy

### Timeline Goal
- ✅ **Complete in 9-10 hours** - All phases done

---

## References

- **Task 42**: TUI Architecture Refactoring (BrowseTimeline pattern)
- **BrowseTimeline**: Reference implementation (`internal/cli/intents/browse_timeline_intent.go`)
- **Modal Patterns**: `docs/MODAL_PATTERNS.md`
- **Forms Guide**: `docs/FORMS_GUIDE.md`
- **Intent Patterns**: `docs/development/INTENT_PATTERNS_LIBRARY.md`
- **TUI Standards**: `docs/TUI_STANDARDS.md`

---

**Status**: Ready for implementation  
**Estimated Completion**: 9-10 hours  
**Next Action**: Begin Phase 1 - Create CVConfigWizardModal
