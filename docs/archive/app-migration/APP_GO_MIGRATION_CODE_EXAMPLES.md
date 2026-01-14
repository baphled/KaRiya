---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# app.go Migration - Detailed Code Examples

**Version**: 1.0
**Date**: 2026-01-03
**Purpose**: Complete code examples for all migration changes

---

## Phase 1: Foundation Changes

### Change 1: Update() Method - Complete Example

**File**: `internal/cli/app/app.go`
**Location**: Line 368 (func (m *Model) Update)
**Effort**: 15 minutes

#### Before (Current):
```go
// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit and back messages first
	switch msg.(type) {
	case models.BackMsg:
		// Special handling for ListScreen - always go back to HomeScreen
		if m.currentScreen == ListScreen {
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
			return m, nil
		}
		// ... 30+ more screen-specific back handling cases
		return m, nil
	case models.QuitMsg:
		return m, tea.Quit
	case models.HelpMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	case models.MainMenuMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = MainMenuScreen
		return m, nil
	}

	// Handle models.ViewEventMsg
	if viewMsg, ok := msg.(models.ViewEventMsg); ok {
		// ... 800+ lines of screen-specific logic
	}
	// ... more message handling
}
```

#### After (New):
```go
// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// PRIORITY 1: Check if we're in intent mode and delegate to router
	// This must be FIRST before any other logic
	if m.inIntentMode {
		return m.handleIntentMessage(msg)
	}

	// Handle global shortcuts (these work in both screen and intent mode)
	switch msg.(type) {
	case models.QuitMsg:
		return m, tea.Quit
	case models.HelpMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	case models.MainMenuMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = MainMenuScreen
		return m, nil
	case models.BackMsg:
		// Screen-based back navigation (only when NOT in intent mode)
		// Special handling for ListScreen - always go back to HomeScreen
		if m.currentScreen == ListScreen {
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
			return m, nil
		}
		// ... rest of existing back handling
		return m, nil
	}

	// Handle models.ViewEventMsg
	if viewMsg, ok := msg.(models.ViewEventMsg); ok {
		// ... existing screen-based logic continues unchanged
	}
	// ... rest of existing logic continues unchanged
}
```

**Key Points**:
- ✅ Intent check is **FIRST** before any other logic
- ✅ All existing screen-based code remains **unchanged**
- ✅ Global shortcuts (Quit, Help, MainMenu) work in both modes
- ✅ BackMsg now checks intent mode before screen logic

---

### Change 2: View() Method - Complete Example

**File**: `internal/cli/app/app.go`
**Location**: Line 1273 (func (m *Model) View)
**Effort**: 10 minutes

#### Before (Current):
```go
func (m *Model) View() string {
	switch m.currentScreen {
	case MainMenuScreen:
		if m.menuModel != nil {
			return m.menuModel.View()
		}
		return "Error: Menu model not initialized\n"
	case HomeScreen:
		return m.renderHome()
	case CaptureScreen:
		if m.formModel != nil {
			m.formModel.SetBreadcrumbs(m.breadcrumbs)
			return m.formModel.View()
		}
		return "Error: Form model not initialized\n"
	case ListScreen:
		if m.listModel != nil {
			m.listModel.SetBreadcrumbs(m.breadcrumbs)
			return m.listModel.View()
		}
		return "Error: List model not initialized\n"
	// ... 28 more case statements
	default:
		return m.renderHome()
	}
}
```

#### After (New):
```go
func (m *Model) View() string {
	// PRIORITY 1: Check if we're in intent mode and delegate to router
	// This must be FIRST before any other logic
	if m.inIntentMode && m.intentRouter != nil {
		return m.intentRouter.View()
	}

	// Screen-based rendering (all existing code continues unchanged)
	switch m.currentScreen {
	case MainMenuScreen:
		if m.menuModel != nil {
			return m.menuModel.View()
		}
		return "Error: Menu model not initialized\n"
	case HomeScreen:
		return m.renderHome()
	case CaptureScreen:
		if m.formModel != nil {
			m.formModel.SetBreadcrumbs(m.breadcrumbs)
			return m.formModel.View()
		}
		return "Error: Form model not initialized\n"
	case ListScreen:
		if m.listModel != nil {
			m.listModel.SetBreadcrumbs(m.breadcrumbs)
			return m.listModel.View()
		}
		return "Error: List model not initialized\n"
	// ... rest of existing cases unchanged
	default:
		return m.renderHome()
	}
}
```

**Key Points**:
- ✅ Intent check is **FIRST** before any other logic
- ✅ Checks both `inIntentMode` flag AND router is not nil
- ✅ All existing screen-based rendering remains **unchanged**
- ✅ Clean delegation to router when in intent mode

---

### Change 3: Global Shortcuts - BackMsg Handler

**File**: `internal/cli/app/app.go`
**Location**: Line 400+ (in Update method, BackMsg case)
**Effort**: 15 minutes

#### Before (Current):
```go
case models.BackMsg:
	// Special handling for ListScreen - always go back to HomeScreen
	if m.currentScreen == ListScreen {
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	}
	// Special handling for BurstListScreen - always go back to HomeScreen
	if m.currentScreen == BurstListScreen {
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	}
	// ... 30+ more screen-specific back handling
	// Special handling for CVListScreen - go back to home
	if m.currentScreen == CVListScreen {
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	}
	// ... more handling
	return m, nil
```

#### After (New):
```go
case models.BackMsg:
	// PRIORITY 1: Check if in intent mode and delegate to router
	if m.inIntentMode {
		// Let the intent handle back navigation
		cmd, err := m.intentRouter.Back()
		if err != nil {
			// If back navigation fails in intent, exit intent mode
			m.deactivateIntent()
			return m, nil
		}
		return m, cmd
	}

	// Screen-based back navigation (only when NOT in intent mode)
	// Special handling for ListScreen - always go back to HomeScreen
	if m.currentScreen == ListScreen {
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	}
	// Special handling for BurstListScreen - always go back to HomeScreen
	if m.currentScreen == BurstListScreen {
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		return m, nil
	}
	// ... rest of existing back handling continues unchanged
	return m, nil
```

**Key Points**:
- ✅ Checks `inIntentMode` **FIRST**
- ✅ Delegates to router's Back() method when in intent mode
- ✅ Handles error case (intent can't go back further)
- ✅ All existing screen-based logic remains **unchanged**

---

### Change 4: Menu Item Activation

**File**: `internal/cli/app/app.go`
**Location**: Line 1611+ (func (m *Model) handleMenuItemSelection)
**Effort**: 1 hour

#### Before (Current):
```go
func (m *Model) handleMenuItemSelection(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "c":
		// Capture Career Event
		m.previousScreen = m.currentScreen
		m.currentScreen = CaptureScreen
		m.formModel = models.NewFormModel(m.cliService)
		return m, nil
	case "l":
		// List Events
		ctx := context.Background()
		m.listModel = models.NewListModel(m.service, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = ListScreen
		return m, nil
	case "b":
		// View Bursts
		ctx := context.Background()
		m.burstListModel = models.NewBurstListModel(m.service, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = BurstListScreen
		return m, m.burstListModel.Init()
	case "g":
		// Generate CV
		var cmd tea.Cmd
		if m.cvConfigManagerModel == nil {
			m.cvConfigManagerModel = models.NewCVConfigManagerModel(models.NewBaseStandardModel(), m.configManager)
			cmd = m.cvConfigManagerModel.Init()
		} else {
			cmd = m.cvConfigManagerModel.RefreshConfigs()
		}
		m.previousScreen = m.currentScreen
		m.currentScreen = CVConfigManagerScreen
		return m, cmd
	case "?":
		// Help
		if m.helpModel == nil {
			m.helpModel = models.NewHelpModel()
		}
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	}

	return m, nil
}
```

#### After (New):
```go
func (m *Model) handleMenuItemSelection(key string) (tea.Model, tea.Cmd) {
	// PRIORITY 1: Activate intents instead of switching screens
	switch key {
	case "c":
		// Capture Career Event (Intent-based)
		return m, m.activateIntent("capture_event", make(map[string]interface{}))
	case "l":
		// Browse Timeline (Intent-based)
		return m, m.activateIntent("browse_timeline", make(map[string]interface{}))
	case "b":
		// Burst Management (Intent-based, to be created in Phase 3)
		// For now, keep legacy implementation
		ctx := context.Background()
		m.burstListModel = models.NewBurstListModel(m.service, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = BurstListScreen
		return m, m.burstListModel.Init()
	case "g":
		// Generate CV (Intent-based)
		return m, m.activateIntent("generate_cv", make(map[string]interface{}))
	case "x":
		// Export Artifact (Intent-based)
		return m, m.activateIntent("export_artifact", make(map[string]interface{}))
	case "s":
		// Configure System (Intent-based)
		return m, m.activateIntent("configure_system", make(map[string]interface{}))
	case "?":
		// Help (still screen-based for now)
		if m.helpModel == nil {
			m.helpModel = models.NewHelpModel()
		}
		m.previousScreen = m.currentScreen
		m.currentScreen = HelpScreen
		return m, nil
	case "t":
		// View Facts (to be migrated in Phase 3)
		if m.factListModel == nil {
			ctx := context.Background()
			m.factListModel = models.NewFactListModel(m.service, ctx)
		} else {
			m.factListModel.Refresh()
		}
		m.previousScreen = m.currentScreen
		m.currentScreen = FactListScreen
		return m, nil
	}

	return m, nil
}
```

**Key Points**:
- ✅ High-value screens now use intent activation
- ✅ Remaining screens kept as-is for Phase 3 migration
- ✅ Simple call to `activateIntent()` with intent name and empty context
- ✅ No model initialization needed (router handles it)

---

## Phase 2: Result Handler Enhancement

### Enhance Result Handlers

**File**: `internal/cli/app/app.go`
**Location**: Lines 238-304 (in NewModel function)
**Effort**: 2 hours

#### Before (Current):
```go
// Register result handlers for all intents
router.RegisterResultHandler("capture_event", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		if captureResult, ok := result.Data.(*intents.CaptureEventResult); ok {
			log.Info("CaptureEvent intent completed with event: %s", captureResult.Event.ID)
			return func() tea.Msg {
				return FormSubmittedMsg{
					Event: captureResult.Event,
				}
			}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("CaptureEvent intent cancelled by user")
		return func() tea.Msg {
			return models.BackMsg{}
		}
	} else if result.Status == intents.Failed {
		log.Error("CaptureEvent intent failed: %v", result.Error)
		return func() tea.Msg {
			return models.BackMsg{}
		}
	}
	return nil
})

router.RegisterResultHandler("browse_timeline", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		log.Info("BrowseTimeline intent completed")
	} else if result.Status == intents.Cancelled {
		log.Info("BrowseTimeline intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("BrowseTimeline intent failed: %v", result.Error)
	}
	return func() tea.Msg {
		return models.BackMsg{}
	}
})

router.RegisterResultHandler("generate_cv", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		log.Info("GenerateCV intent completed")
	} else if result.Status == intents.Cancelled {
		log.Info("GenerateCV intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("GenerateCV intent failed: %v", result.Error)
	}
	return func() tea.Msg {
		return models.BackMsg{}
	}
})
// ... similar for ExportArtifact and ConfigureSystem
```

#### After (Enhanced):
```go
// Register result handlers for all intents
router.RegisterResultHandler("capture_event", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		if captureResult, ok := result.Data.(*intents.CaptureEventResult); ok {
			log.Info("CaptureEvent intent completed with event: %s", captureResult.Event.ID)
			// Return FormSubmittedMsg to trigger app state update
			return func() tea.Msg {
				return FormSubmittedMsg{
					Event: captureResult.Event,
				}
			}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("CaptureEvent intent cancelled by user")
		// Exit intent mode and return to screen mode
		return nil
	} else if result.Status == intents.Failed {
		log.Error("CaptureEvent intent failed: %v", result.Error)
		// Could return error message or retry prompt
		return nil
	}
	return nil
})

router.RegisterResultHandler("browse_timeline", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		if timelineResult, ok := result.Data.(*intents.BrowseTimelineResult); ok {
			log.Info("BrowseTimeline intent completed with event: %s", timelineResult.SelectedEventID)
			// Could trigger viewing the selected event
			if timelineResult.SelectedEventID != "" {
				return func() tea.Msg {
					return models.ViewEventMsg{
						EventID: timelineResult.SelectedEventID,
					}
				}
			}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("BrowseTimeline intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("BrowseTimeline intent failed: %v", result.Error)
	}
	// Return to screen mode
	return nil
})

router.RegisterResultHandler("generate_cv", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		if cvResult, ok := result.Data.(*intents.GenerateCVResult); ok {
			log.Info("GenerateCV intent completed with CV: %s", cvResult.CV.ID)
			// Could trigger export dialog or success screen
			return func() tea.Msg {
				return models.CVGeneratedMsg{
					CV: cvResult.CV,
				}
			}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("GenerateCV intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("GenerateCV intent failed: %v", result.Error)
	}
	// Return to screen mode
	return nil
})

router.RegisterResultHandler("export_artifact", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		if exportResult, ok := result.Data.(*intents.ExportArtifactResult); ok {
			log.Info("ExportArtifact intent completed with file: %s", exportResult.FilePath)
			// Show success message with file path
			return func() tea.Msg {
				return models.ExportSuccessMsg{
					FilePath: exportResult.FilePath,
				}
			}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("ExportArtifact intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("ExportArtifact intent failed: %v", result.Error)
	}
	// Return to screen mode
	return nil
})

router.RegisterResultHandler("configure_system", func(result *intents.IntentResult[interface{}]) tea.Cmd {
	if result.Status == intents.Completed {
		log.Info("ConfigureSystem intent completed")
		// Could trigger config reload or success message
		return func() tea.Msg {
			return models.ConfigSavedMsg{}
		}
	} else if result.Status == intents.Cancelled {
		log.Info("ConfigureSystem intent cancelled by user")
	} else if result.Status == intents.Failed {
		log.Error("ConfigureSystem intent failed: %v", result.Error)
	}
	// Return to screen mode
	return nil
})
```

**Key Points**:
- ✅ Extract typed result data from IntentResult[interface{}]
- ✅ Return meaningful messages (not just BackMsg)
- ✅ Log important information
- ✅ Handle all three status cases (Completed, Cancelled, Failed)
- ✅ Update app state based on results

---

## Complete Before/After Comparison

### Lines of Code Impact

| Component | Before | After | Change |
|-----------|--------|-------|--------|
| Update() | 906 | 920 | +14 lines |
| View() | 180 | 190 | +10 lines |
| handleMenuItemSelection() | 50 | 65 | +15 lines |
| Result Handlers | 67 | 150 | +83 lines |
| **Total app.go** | 1766 | 1850 | +84 lines |

**Note**: These are temporary increases. After Phase 4 cleanup, app.go will be **reduced to ~400 lines**.

---

## Testing Checklist

### Unit Tests
```bash
# Test Update() delegation
go test -v -run TestUpdateIntentMode ./internal/cli/app/...

# Test View() delegation
go test -v -run TestViewIntentMode ./internal/cli/app/...

# Test intent activation
go test -v -run TestActivateIntent ./internal/cli/app/...

# Test result handlers
go test -v -run TestResultHandlers ./internal/cli/app/...
```

### Integration Tests
```bash
# Test end-to-end capture workflow
go test -v -run TestCaptureEventWorkflow ./internal/cli/app/...

# Test end-to-end timeline workflow
go test -v -run TestBrowseTimelineWorkflow ./internal/cli/app/...

# Test menu activation
go test -v -run TestMenuActivation ./internal/cli/app/...
```

### Manual Testing
1. Run `go run ./cmd/kariya`
2. Activate "Capture Event" from menu
3. Complete a capture workflow
4. Press Esc to go back
5. Activate "Browse Timeline" from menu
6. Browse events
7. Press Esc to go back
8. Verify no crashes or panics

---

## Common Pitfalls to Avoid

### ❌ Wrong: Forgetting to check inIntentMode

```go
// BAD - doesn't check inIntentMode
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case models.BackMsg:
        // ... screen logic immediately
    }
}
```

### ✅ Right: Check inIntentMode FIRST

```go
// GOOD - checks inIntentMode first
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.inIntentMode {
        return m.handleIntentMessage(msg)
    }

    switch msg.(type) {
    case models.BackMsg:
        // ... screen logic
    }
}
```

### ❌ Wrong: Nil pointer dereference

```go
// BAD - doesn't check if router is nil
if m.inIntentMode {
    return m.intentRouter.View()  // Could panic!
}
```

### ✅ Right: Check for nil

```go
// GOOD - checks both flags and nil
if m.inIntentMode && m.intentRouter != nil {
    return m.intentRouter.View()
}
```

### ❌ Wrong: Modifying screen state in intent mode

```go
// BAD - modifies screen state while in intent mode
if m.inIntentMode {
    m.currentScreen = SomeScreen  // Confusing!
    return m.handleIntentMessage(msg)
}
```

### ✅ Right: Keep modes separate

```go
// GOOD - intent mode doesn't touch screen state
if m.inIntentMode {
    return m.handleIntentMessage(msg)
}
// Screen state only modified when not in intent mode
m.currentScreen = SomeScreen
```

---

## Summary

**Phase 1 Implementation**:
- 4 changes to app.go
- ~84 lines added (temporary)
- 0 lines removed
- 2-3 hours total effort
- 0 breaking changes
- Full backward compatibility

**Result**: Intent system becomes fully functional while maintaining all existing screen-based functionality.


