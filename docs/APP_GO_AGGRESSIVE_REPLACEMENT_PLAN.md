# app.go Aggressive Replacement Plan - Full Intent-Based Architecture

**Document Version**: 1.0
**Date Created**: 2026-01-03
**Status**: READY FOR IMMEDIATE IMPLEMENTATION
**Timeline**: 2-3 weeks (aggressive, no transition period)
**Approach**: Complete replacement, zero legacy code
**User**: Solo developer (no backward compatibility needed)

---

## Executive Summary

This document outlines an **aggressive, no-compromise replacement** of the entire screen-based architecture with the intent-based system. Since you are the sole user and want to move fast, we will:

✅ **Remove all 31 screens entirely**
✅ **Remove all 37+ model fields**
✅ **Rebuild app.go from scratch** (~200 lines vs 1766)
✅ **Replace all legacy models** with intent-based equivalents
✅ **Complete migration in 2-3 weeks** (full-time)
✅ **Zero backward compatibility** (no transition period)

### Key Differences from Phased Plan
- ❌ **NO dual-mode operation** - Only intent mode
- ❌ **NO legacy screen code** - All removed
- ❌ **NO gradual migration** - All at once
- ✅ **MUCH simpler** - Clean slate
- ✅ **MUCH faster** - No compatibility concerns
- ✅ **MUCH cleaner** - No legacy code pollution

### Expected Result
- **app.go**: 1766 lines → ~200 lines (88% reduction)
- **Update() method**: 906 lines → ~30 lines (97% reduction)
- **Model fields**: 37+ → ~5 (86% reduction)
- **Code quality**: Dramatically improved
- **Maintainability**: Significantly better

---

## Current State vs Target State

### Current State (Screen-Based)
```
app.go (1766 lines)
├── 31 Screen constants
├── 37+ model fields (one per screen)
├── NewModel() - 200+ lines of initialization
├── Update() - 906 lines of message handling
├── View() - 180 lines of screen rendering
├── 20+ helper methods for screen-specific logic
└── Breadcrumb/workflow/navigation complexity
```

### Target State (Intent-Based)
```
app.go (~200 lines)
├── Minimal constants (HomeScreen, HelpScreen only)
├── 5-7 core model fields (services, router, state)
├── NewModel() - 50 lines of initialization
├── Update() - 30 lines (delegates to router)
├── View() - 20 lines (delegates to router)
├── 3-4 helper methods (menu, shortcuts)
└── Clean, focused, maintainable
```

---

## Implementation Strategy

### Phase 1: Preparation (2-3 days)

#### 1.1 Audit All Screen Models
**Goal**: Understand what each screen model does

**Current Screen Models**:
1. FormModel (CaptureScreen)
2. ListModel (ListScreen)
3. DetailsModel (ViewScreen)
4. SuccessModel (SuccessScreen)
5. ActionMenuModel
6. FactActionMenuModel
7. ConfirmationDialog
8. ImportReviewModel
9. ImportProgressModel
10. MetadataReviewModel
11. MetadataEditorModel
12. BulkOperationsModel
13. BurstSuggestionModel
14. BurstListModel
15. BurstDetailsModel
16. BurstEditorModel
17. FactListModel
18. FactsResultsModel
19. FactEditorModel
20. FactDetailsModel
21. MenuModel
22. HelpModel
23. CVConfigManagerModel
24. CVGeneratorModel
25. CVPreviewModel
26. CVListModel
27. CVExportDialogModel
28. CVExportSuccessModel
29. CVExportProgressModel

**Tasks**:
1. Document what each model does
2. Identify which intent should handle it
3. Note any special logic that needs preservation
4. Check for cross-model dependencies

**Effort**: 4 hours

#### 1.2 Identify Intent Gaps
**Goal**: Map legacy screens to intents, identify missing intents

**Mapping**:
- CaptureScreen → CaptureEvent Intent ✅ (already exists)
- ListScreen → BrowseTimeline Intent ✅ (already exists)
- ViewScreen → Part of BrowseTimeline Intent ✅
- CVConfigManagerScreen → ConfigureSystem Intent ✅ (already exists)
- CVGeneratorScreen → GenerateCV Intent ✅ (already exists)
- CVPreviewScreen → Part of GenerateCV Intent ✅
- CVExportDialogScreen → ExportArtifact Intent ✅ (already exists)
- CVListScreen → BrowseTimeline Intent (extend) ⚠️
- BurstListScreen → **New: BurstManagement Intent** ❌
- BurstDetailsScreen → Part of BurstManagement Intent ❌
- BurstEditorScreen → Part of BurstManagement Intent ❌
- BurstSuggestionScreen → Part of BurstManagement Intent ❌
- FactListScreen → **New: FactManagement Intent** ❌
- FactDetailsScreen → Part of FactManagement Intent ❌
- FactEditorScreen → Part of FactManagement Intent ❌
- FactActionMenuScreen → Part of FactManagement Intent ❌
- FactsResultsScreen → Part of FactManagement Intent ❌
- ImportReviewScreen → **New: ImportWizard Intent** ❌
- ImportProgressScreen → Part of ImportWizard Intent ❌
- MetadataReviewScreen → **New: MetadataEditor Intent** ❌
- MetadataEditorScreen → Part of MetadataEditor Intent ❌
- BulkOperationsScreen → **New: BulkOperations Intent** ❌
- SuccessScreen → Generic success message (not an intent)
- ActionMenuScreen → Generic menu component (not an intent)
- MainMenuScreen → Keep as startup screen
- HomeScreen → Keep as home screen
- HelpScreen → Keep as help screen
- ConfirmationScreen → Modal dialog (not an intent)
- QuitScreen → Quit handler (not an intent)

**Missing Intents to Create**:
1. BurstManagement Intent
2. FactManagement Intent
3. ImportWizard Intent
4. MetadataEditor Intent
5. BulkOperations Intent

**Effort**: 3 hours

#### 1.3 Identify Shared Logic
**Goal**: Find logic used by multiple screens

**Common Patterns**:
- Breadcrumb management → Move to intent router or global state
- Workflow state tracking → Move to GlobalContext
- Window size handling → Pass to intents via context
- Event operations (CRUD) → Already in services
- Confirmation dialogs → Create modal component
- Success messages → Return from intent results

**Effort**: 2 hours

**Total Phase 1 Effort**: 9 hours

---

### Phase 2: Create Missing Intents (4-5 days)

All 5 core intents already exist. We need to create 5 new intents:

#### 2.1 BurstManagement Intent

**Screens to Replace**:
- BurstListScreen
- BurstDetailsScreen
- BurstEditorScreen
- BurstSuggestionScreen

**Implementation**:
```
internal/cli/intents/burst_management.go
  ├── BurstManagementState (enum)
  ├── BurstManagementContext (input data)
  ├── BurstManagementResult (output data)
  └── BurstManagementModel (state machine)

internal/cli/intents/burst_management_intent.go
  ├── Init() - Load bursts
  ├── Update() - Handle all states
  ├── View() - Render current state
  └── Result() - Return result

internal/cli/intents/burst_management_test.go
  ├── State transition tests
  ├── View rendering tests
  ├── Result handling tests
  └── Integration tests
```

**States**:
1. StateListBursts - Show list of bursts
2. StateViewBurst - Show burst details
3. StateEditBurst - Edit burst
4. StateSuggestBursts - Show burst suggestions
5. StateConfirm - Confirm changes

**Effort**: 16 hours (8 hours implementation, 8 hours testing)

#### 2.2 FactManagement Intent

**Screens to Replace**:
- FactListScreen
- FactDetailsScreen
- FactEditorScreen
- FactActionMenuScreen
- FactsResultsScreen

**Implementation**:
```
internal/cli/intents/fact_management.go
  ├── FactManagementState (enum)
  ├── FactManagementContext (input data)
  ├── FactManagementResult (output data)
  └── FactManagementModel (state machine)

internal/cli/intents/fact_management_intent.go
  ├── Init() - Load facts
  ├── Update() - Handle all states
  ├── View() - Render current state
  └── Result() - Return result

internal/cli/intents/fact_management_test.go
  ├── State transition tests
  ├── View rendering tests
  ├── Result handling tests
  └── Integration tests
```

**States**:
1. StateListFacts - Show list of facts
2. StateViewFact - Show fact details
3. StateEditFact - Edit fact
4. StateReviewFacts - Review extracted facts
5. StateConfirm - Confirm changes

**Effort**: 16 hours (8 hours implementation, 8 hours testing)

#### 2.3 ImportWizard Intent

**Screens to Replace**:
- ImportReviewScreen
- ImportProgressScreen

**Implementation**:
```
internal/cli/intents/import_wizard.go
  ├── ImportWizardState (enum)
  ├── ImportWizardContext (input data)
  ├── ImportWizardResult (output data)
  └── ImportWizardModel (state machine)

internal/cli/intents/import_wizard_intent.go
  ├── Init() - Start with file
  ├── Update() - Handle import progress
  ├── View() - Show progress/review
  └── Result() - Return import result

internal/cli/intents/import_wizard_test.go
  ├── State transition tests
  ├── Progress tracking tests
  ├── Result handling tests
  └── Integration tests
```

**States**:
1. StateSelectFile - File selection
2. StateReviewImport - Review before import
3. StateImportInProgress - Show progress
4. StateImportComplete - Show results
5. StateConfirm - Confirm completion

**Effort**: 12 hours (6 hours implementation, 6 hours testing)

#### 2.4 MetadataEditor Intent

**Screens to Replace**:
- MetadataReviewScreen
- MetadataEditorScreen

**Implementation**:
```
internal/cli/intents/metadata_editor.go
  ├── MetadataEditorState (enum)
  ├── MetadataEditorContext (input data)
  ├── MetadataEditorResult (output data)
  └── MetadataEditorModel (state machine)

internal/cli/intents/metadata_editor_intent.go
  ├── Init() - Load metadata
  ├── Update() - Handle editing
  ├── View() - Show editor
  └── Result() - Return changes

internal/cli/intents/metadata_editor_test.go
  ├── State transition tests
  ├── Validation tests
  ├── Result handling tests
  └── Integration tests
```

**States**:
1. StateReviewMetadata - Show metadata
2. StateEditMetadata - Edit fields
3. StateConfirm - Confirm changes

**Effort**: 10 hours (5 hours implementation, 5 hours testing)

#### 2.5 BulkOperations Intent

**Screens to Replace**:
- BulkOperationsScreen

**Implementation**:
```
internal/cli/intents/bulk_operations.go
  ├── BulkOperationsState (enum)
  ├── BulkOperationsContext (input data)
  ├── BulkOperationsResult (output data)
  └── BulkOperationsModel (state machine)

internal/cli/intents/bulk_operations_intent.go
  ├── Init() - Load events
  ├── Update() - Handle bulk operations
  ├── View() - Show operations menu
  └── Result() - Return results

internal/cli/intents/bulk_operations_test.go
  ├── State transition tests
  ├── Operation tests
  ├── Result handling tests
  └── Integration tests
```

**States**:
1. StateSelectOperation - Choose operation
2. StateConfigureOperation - Set parameters
3. StateExecuteOperation - Run operation
4. StateConfirm - Confirm results

**Effort**: 10 hours (5 hours implementation, 5 hours testing)

**Total Phase 2 Effort**: 64 hours (8 days)

---

### Phase 3: Rebuild app.go (3-4 days)

#### 3.1 Create New Minimal app.go

**New Structure**:
```go
package app

import (
	// Minimal imports
	"context"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen constants - only what we actually need
type Screen string

const (
	HomeScreen       Screen = "home"       // Main menu
	HelpScreen       Screen = "help"       // Help screen
	MainMenuScreen   Screen = "main_menu"  // Menu before home
	QuitScreen       Screen = "quit"       // Quit confirmation
)

// Model - minimal, focused on app-level state
type Model struct {
	// Services
	cliService    *service.CLIEventService
	careerService *career.Service
	logger        *logger.Logger

	// Intent system
	intentRouter intents.IntentRouter
	inIntentMode bool

	// App-level state (only what's truly app-level)
	currentScreen  Screen
	previousScreen Screen
	width          int
	height         int

	// Reusable components
	menuModel *models.MenuModel
	helpModel *models.HelpModel
}

// NewModel creates a minimal app model
func NewModel(cliService *service.CLIEventService, careerService *career.Service) *Model {
	// Initialize router with all intents
	router := intents.NewDefaultIntentRouter()

	// Register all 10 intents
	registerIntents(router, careerService)

	return &Model{
		cliService:    cliService,
		careerService: careerService,
		logger:        logger.DefaultLogger(),
		intentRouter:  router,
		inIntentMode:  false,
		currentScreen: MainMenuScreen,
		width:         80,
		height:        24,
	}
}

// Init initializes the app
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// PRIORITY: Check intent mode
	if m.inIntentMode {
		return m.handleIntentMessage(msg)
	}

	// Handle global shortcuts
	switch msg := msg.(type) {
	case tea.QuitMsg:
		return m, tea.Quit
	case models.BackMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	case models.HelpMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	case models.MainMenuMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = MainMenuScreen
		return m, nil
	}

	// Handle menu selection
	if menuMsg, ok := msg.(models.MenuSelectionMsg); ok {
		return m.handleMenuSelection(menuMsg.Key)
	}

	// Update current screen models
	switch m.currentScreen {
	case HomeScreen:
		return m, nil
	case HelpScreen:
		if m.helpModel != nil {
			_, cmd := m.helpModel.Update(msg)
			return m, cmd
		}
	case MainMenuScreen:
		if m.menuModel != nil {
			_, cmd := m.menuModel.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View renders the app
func (m *Model) View() string {
	// PRIORITY: Check intent mode
	if m.inIntentMode && m.intentRouter != nil {
		return m.intentRouter.View()
	}

	// Render current screen
	switch m.currentScreen {
	case MainMenuScreen:
		if m.menuModel == nil {
			m.menuModel = models.NewMenuModel(false, "")
		}
		return m.menuModel.View()
	case HomeScreen:
		return renderHome()
	case HelpScreen:
		if m.helpModel == nil {
			m.helpModel = models.NewHelpModel()
		}
		return m.helpModel.View()
	case QuitScreen:
		return "Goodbye!\n"
	}

	return renderHome()
}

// Helper methods
func (m *Model) handleIntentMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.intentRouter == nil {
		return m, nil
	}

	cmd, result := m.intentRouter.HandleMessage(msg)
	if result != nil {
		m.inIntentMode = false
		m.currentScreen = HomeScreen
	}

	return m, cmd
}

func (m *Model) handleMenuSelection(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "c":
		return m, m.activateIntent("capture_event", make(map[string]interface{}))
	case "l":
		return m, m.activateIntent("browse_timeline", make(map[string]interface{}))
	case "g":
		return m, m.activateIntent("generate_cv", make(map[string]interface{}))
	case "x":
		return m, m.activateIntent("export_artifact", make(map[string]interface{}))
	case "s":
		return m, m.activateIntent("configure_system", make(map[string]interface{}))
	case "b":
		return m, m.activateIntent("burst_management", make(map[string]interface{}))
	case "f":
		return m, m.activateIntent("fact_management", make(map[string]interface{}))
	case "i":
		return m, m.activateIntent("import_wizard", make(map[string]interface{}))
	case "m":
		return m, m.activateIntent("metadata_editor", make(map[string]interface{}))
	case "o":
		return m, m.activateIntent("bulk_operations", make(map[string]interface{}))
	case "?":
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	}
	return m, nil
}

func (m *Model) activateIntent(intentName string, ctx map[string]interface{}) tea.Cmd {
	ctx["previousScreen"] = m.currentScreen
	cmd, err := m.intentRouter.ActivateIntent(intentName, ctx)
	if err != nil {
		m.logger.Error("Failed to activate intent %s: %v", intentName, err)
		return nil
	}
	m.inIntentMode = true
	return cmd
}

func renderHome() string {
	// Simple home screen
	return "Welcome to KaRiya\n"
}

// Register all intents with the router
func registerIntents(router intents.IntentRouter, svc *career.Service) {
	// Register all 10 intents here
	// (implementation details omitted for brevity)
}
```

**Size Comparison**:
- Old: 1766 lines
- New: ~200 lines
- Reduction: 88%

**Effort**: 8 hours (4 hours implementation, 4 hours testing)

#### 3.2 Remove All Legacy Files

**Files to Delete**:
- All screen-specific model files in `internal/cli/models/` that are now handled by intents
- Unused helper functions
- Unused components

**Effort**: 2 hours

#### 3.3 Update Service Layer

**Changes**:
- Remove screen-specific service methods
- Consolidate service interfaces
- Ensure all intent dependencies are met

**Effort**: 4 hours

**Total Phase 3 Effort**: 14 hours (2 days)

---

### Phase 4: Testing & Validation (2-3 days)

#### 4.1 Unit Tests for New Intents
- BurstManagement Intent tests (8 hours) - already counted in Phase 2
- FactManagement Intent tests (8 hours) - already counted in Phase 2
- ImportWizard Intent tests (6 hours) - already counted in Phase 2
- MetadataEditor Intent tests (5 hours) - already counted in Phase 2
- BulkOperations Intent tests (5 hours) - already counted in Phase 2

#### 4.2 Integration Tests
- Test all intents activate from menu (4 hours)
- Test all workflows end-to-end (8 hours)
- Test result handling (4 hours)
- Test navigation (4 hours)

**Effort**: 20 hours (3 days)

#### 4.3 Performance & Quality
- Run benchmarks (2 hours)
- Performance profiling (4 hours)
- Code review & cleanup (4 hours)
- Documentation (4 hours)

**Effort**: 14 hours (2 days)

**Total Phase 4 Effort**: 34 hours (4-5 days)

---

## Complete Timeline

### Week 1: Preparation & Missing Intents
```
Mon-Tue: Phase 1 - Audit and planning (9 hours)
Wed-Sun: Phase 2 - Create 5 missing intents (64 hours)
```

### Week 2: Rebuild & Testing
```
Mon-Tue: Phase 3 - Rebuild app.go (14 hours)
Wed-Fri: Phase 4 - Testing & validation (34 hours)
Sat-Sun: Buffer & final polish
```

**Total Effort**: 121 hours (~3 weeks full-time)

---

## Detailed Breakdown by Intent

### Already Implemented (5 intents) ✅

| Intent | Status | Tests | Coverage | Work Needed |
|--------|--------|-------|----------|-------------|
| CaptureEvent | ✅ COMPLETE | 30+ | 88.3% | Register & test |
| BrowseTimeline | ✅ COMPLETE | 37 | >90% | Register & test |
| GenerateCV | ✅ COMPLETE | 41 | >90% | Register & test |
| ExportArtifact | ✅ COMPLETE | 411 | >90% | Register & test |
| ConfigureSystem | ✅ COMPLETE | 400+ | >90% | Register & test |

### To Create (5 intents) ❌

| Intent | Effort | Tests Needed | Complexity |
|--------|--------|--------------|------------|
| BurstManagement | 16 hours | 30+ | MEDIUM |
| FactManagement | 16 hours | 40+ | MEDIUM-HIGH |
| ImportWizard | 12 hours | 25+ | MEDIUM |
| MetadataEditor | 10 hours | 20+ | LOW-MEDIUM |
| BulkOperations | 10 hours | 20+ | LOW |

**Total New Intent Creation**: 64 hours

---

## Key Implementation Patterns

### For Each New Intent

#### 1. Create State Machine
```go
type BurstManagementState string

const (
	StateListBursts    BurstManagementState = "list"
	StateViewBurst     BurstManagementState = "view"
	StateEditBurst     BurstManagementState = "edit"
	StateSuggestBursts BurstManagementState = "suggest"
	StateConfirm       BurstManagementState = "confirm"
)
```

#### 2. Define Context & Result
```go
type BurstManagementContext struct {
	Bursts []*career.Burst
	// ... other data
}

type BurstManagementResult struct {
	Changes map[string]interface{}
	// ... result data
}
```

#### 3. Implement Intent Interface
```go
type BurstManagementModel struct {
	state BurstManagementState
	data  *BurstManagementContext
	// ... other fields
}

func (b *BurstManagementModel) Init(ctx context.Context) tea.Cmd { ... }
func (b *BurstManagementModel) Update(msg tea.Msg) tea.Cmd { ... }
func (b *BurstManagementModel) View() string { ... }
func (b *BurstManagementModel) Result() *IntentResult[interface{}] { ... }
```

#### 4. Write Tests
```go
var _ = Describe("BurstManagement Intent", func() {
	It("should list bursts", func() { ... })
	It("should view burst details", func() { ... })
	It("should edit burst", func() { ... })
	// ... more tests
})
```

---

## What Gets Deleted

### Model Files (30+ files)
```
internal/cli/models/
├── form_model.go (replaced by CaptureEvent Intent)
├── list_model.go (replaced by BrowseTimeline Intent)
├── details_model.go (replaced by BrowseTimeline Intent)
├── burst_list_model.go (replaced by BurstManagement Intent)
├── burst_details_model.go (replaced by BurstManagement Intent)
├── burst_editor_model.go (replaced by BurstManagement Intent)
├── fact_list_model.go (replaced by FactManagement Intent)
├── fact_details_model.go (replaced by FactManagement Intent)
├── fact_editor_model.go (replaced by FactManagement Intent)
├── import_review_model.go (replaced by ImportWizard Intent)
├── import_progress_model.go (replaced by ImportWizard Intent)
├── metadata_review_model.go (replaced by MetadataEditor Intent)
├── metadata_editor_model.go (replaced by MetadataEditor Intent)
├── cv_config_manager_model.go (replaced by ConfigureSystem Intent)
├── cv_generator_model.go (replaced by GenerateCV Intent)
├── cv_preview_model.go (replaced by GenerateCV Intent)
├── cv_list_model.go (replaced by BrowseTimeline Intent)
├── cv_export_dialog_model.go (replaced by ExportArtifact Intent)
├── cv_export_success_model.go (handled by result)
├── cv_export_progress_model.go (handled by intent progress)
├── bulk_operations_model.go (replaced by BulkOperations Intent)
├── action_menu_model.go (handled by intents)
├── confirmation_dialog.go (replaced by modal component)
└── ... and more
```

### app.go Sections (1500+ lines deleted)
```
❌ All 31 screen constants
❌ All 37+ model fields
❌ Entire NewModel() initialization (150+ lines)
❌ Entire Update() method (906 lines)
❌ Entire View() method (180 lines)
❌ 20+ helper methods
❌ Breadcrumb management logic
❌ Screen-specific navigation logic
❌ Workflow state integration (moved to GlobalContext)
```

---

## Risks & Mitigation

### Risk: Breaking Changes
**Mitigation**:
- ✅ You're the sole user
- ✅ All intents are fully tested
- ✅ Can test locally before deploying
- ✅ Git history allows rollback

### Risk: Missing Features
**Mitigation**:
- ✅ Audit phase identifies all features
- ✅ Each intent maps to existing screens
- ✅ No new features, just reorganization
- ✅ Comprehensive testing

### Risk: Performance Issues
**Mitigation**:
- ✅ Intent system is already benchmarked
- ✅ Phase 4 includes performance testing
- ✅ Simpler code should be faster
- ✅ Profile and optimize as needed

### Risk: Timeline Slippage
**Mitigation**:
- ✅ Aggressive timeline built in
- ✅ Parallel implementation possible
- ✅ Clear checkpoints at each phase
- ✅ Buffer time in schedule

---

## Success Criteria

### Phase 1 Complete
- ✅ All screens audited
- ✅ All intents mapped
- ✅ No missing features identified

### Phase 2 Complete
- ✅ 5 new intents created
- ✅ All intents have >90% test coverage
- ✅ All tests passing

### Phase 3 Complete
- ✅ app.go rebuilt (~200 lines)
- ✅ All legacy files removed
- ✅ Code compiles without errors

### Phase 4 Complete
- ✅ All intents activate from menu
- ✅ All workflows end-to-end functional
- ✅ All tests passing (164+ total)
- ✅ Performance benchmarks met
- ✅ Code quality verified

---

## File Structure After Replacement

```
internal/
├── cli/
│   ├── app/
│   │   ├── app.go (200 lines, clean and focused)
│   │   ├── app_test.go (comprehensive tests)
│   │   └── messages.go (message types)
│   ├── intents/
│   │   ├── contract.go ✅
│   │   ├── router.go ✅
│   │   ├── testing.go ✅
│   │   ├── capture_event_intent.go ✅
│   │   ├── browse_timeline_intent.go ✅
│   │   ├── generate_cv_intent.go ✅
│   │   ├── export_artifact_intent.go ✅
│   │   ├── configure_system_intent.go ✅
│   │   ├── burst_management_intent.go ✨ NEW
│   │   ├── fact_management_intent.go ✨ NEW
│   │   ├── import_wizard_intent.go ✨ NEW
│   │   ├── metadata_editor_intent.go ✨ NEW
│   │   ├── bulk_operations_intent.go ✨ NEW
│   │   └── *_test.go (comprehensive tests)
│   ├── models/
│   │   ├── base_model.go (shared components)
│   │   ├── menu_model.go (main menu)
│   │   ├── help_model.go (help screen)
│   │   └── messages.go (message types)
│   ├── components/
│   │   ├── card.go
│   │   ├── list.go
│   │   ├── form.go
│   │   └── ... (reusable UI components)
│   ├── styles/
│   │   └── styles.go (all styling)
│   └── context/
│       └── global.go (global context)
├── domain/
├── service/
├── repository/
└── logger/
```

---

## No Transition Period

### Aggressive Approach Advantages
✅ **Simpler** - One clear path, no compatibility concerns
✅ **Faster** - No need to maintain dual systems
✅ **Cleaner** - No legacy code pollution
✅ **Safer** - Less surface area for bugs
✅ **Better** - Forces clean implementation

### Aggressive Approach Challenges
⚠️ **Riskier** - All or nothing approach
⚠️ **Demanding** - Requires sustained focus
⚠️ **Intensive** - 3 weeks of solid work
⚠️ **Thorough** - Needs comprehensive testing

### Mitigation
- ✅ Work on feature branch
- ✅ Test locally before merging
- ✅ Git allows rollback
- ✅ You're the sole user
- ✅ Can iterate quickly

---

## Implementation Checklist

### Phase 1: Preparation
- [ ] Audit all 29 legacy screen models
- [ ] Document feature mapping to intents
- [ ] Identify any missing functionality
- [ ] Plan new intent implementations
- [ ] Estimate effort for each intent

### Phase 2: Create Missing Intents
- [ ] Create BurstManagement intent (16 hours)
- [ ] Create FactManagement intent (16 hours)
- [ ] Create ImportWizard intent (12 hours)
- [ ] Create MetadataEditor intent (10 hours)
- [ ] Create BulkOperations intent (10 hours)
- [ ] All tests passing (30+ tests per intent)

### Phase 3: Rebuild app.go
- [ ] Write new minimal app.go (~200 lines)
- [ ] Register all 10 intents
- [ ] Implement menu selection
- [ ] Implement intent activation
- [ ] Remove all legacy code
- [ ] Verify compilation

### Phase 4: Testing & Validation
- [ ] Unit tests for new intents (20+ hours)
- [ ] Integration tests (20 hours)
- [ ] End-to-end workflow testing (8 hours)
- [ ] Performance benchmarking (6 hours)
- [ ] Code quality review (4 hours)
- [ ] Documentation update (4 hours)

---

## Next Steps

### Immediate
1. ✅ Review this plan
2. ✅ Confirm 3-week timeline works
3. ✅ Prepare development environment
4. ✅ Create feature branch

### This Week
1. Start Phase 1 (Preparation)
2. Audit all legacy models
3. Document feature mapping
4. Plan Phase 2

### Next Week
1. Complete Phase 1
2. Begin Phase 2 (Create missing intents)
3. Start with BurstManagement intent
4. Move to FactManagement intent

### Week 3
1. Complete Phase 2
2. Begin Phase 3 (Rebuild app.go)
3. Start Phase 4 (Testing)

### Week 4
1. Complete Phase 4
2. Final validation
3. Merge to main
4. Deploy

---

## Conclusion

This aggressive replacement plan provides a **clear, direct path** to completely replace the screen-based architecture with the intent-based system in **3 weeks**.

### Key Advantages
✅ **Fast** - 3 weeks vs 6+ weeks
✅ **Clean** - No legacy code pollution
✅ **Simple** - One clear implementation path
✅ **Comprehensive** - All features preserved
✅ **Well-Tested** - 164+ tests included

### Expected Result
- **app.go**: 1766 → ~200 lines (88% reduction)
- **Update()**: 906 → ~30 lines (97% reduction)
- **Model fields**: 37+ → ~5 (86% reduction)
- **Code quality**: Dramatically improved
- **Maintainability**: Significantly better
- **Test coverage**: 164+ tests passing

### Start Date
Ready to begin immediately upon your approval.

---

**Document Version**: 1.0
**Status**: READY FOR IMPLEMENTATION
**Timeline**: 3 weeks (121 hours)
**Approach**: Aggressive, complete replacement, no transition period


