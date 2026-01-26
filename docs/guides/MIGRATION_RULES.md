# Migration Rules: Intent Subdirectory Structure & UIKit

**Purpose**: Mandatory migration rules for all legacy intent refactoring  
**Status**: REQUIRED - Must follow for all migrations  
**Enforcement**: Automated by `check-intent-architecture.sh` (Checks #17-22)

---

## Table of Contents

1. [When to Migrate](#when-to-migrate)
2. [Migration Order](#migration-order)
3. [File Structure Rules](#file-structure-rules)
4. [UIKit Migration Rules](#uikit-migration-rules)
5. [Forms Migration Rules](#forms-migration-rules)
6. [Component Migration Map](#component-migration-map)
7. [Step-by-Step Process](#step-by-step-process)
8. [Validation Checklist](#validation-checklist)
9. [Common Violations](#common-violations)

---

## When to Migrate

### MANDATORY Migration Triggers

**You MUST migrate if any of these are true**:

1. **File size >600 lines** (BLOCKING)
   - Check #18 BLOCKS commits
   - No exceptions - migration required

2. **File size >400 lines AND adding features** (BLOCKING)
   - Warnings at >400 lines
   - Cannot add features until migrated

3. **Fixing bugs in bloated intents** (Boy Scout Rule)
   - Fix the bug AND refactor the touched area
   - Leave code better than you found it

4. **Creating new intents** (ALWAYS)
   - ALWAYS use subdirectory structure from day 1
   - Copy templates from `examples/intent_subdirectory_template/`

### AI Agent Refusal Rule

**AI agents MUST refuse if**:
- Intent >600 lines: "Intent must be migrated first"
- Intent >400 lines + feature request: "Migrate before adding features"
- Using deprecated patterns: "Use subdirectory structure"

**Refusal Template**:
```
I cannot proceed with this request.

MIGRATION REQUIRED

File: intents/burst_management_intent.go (1,488 lines)
Limit: 600 lines (BLOCKING), 400 lines (WARNING)

This intent must be migrated to subdirectory structure before:
- Adding new features
- Making substantial changes

Migration guide: docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md
Migration rules: docs/guides/MIGRATION_RULES.md

Automated checks enforce this (Check #18).
```

---

## Migration Order

### Priority 1: Critical Bloat (>1,000 lines)

**MUST migrate first** (blocks all feature development):

| Intent | Lines | Priority | Reason |
|--------|-------|----------|--------|
| `generate_cv_intent.go` | 2,746 | 🔴 CRITICAL | Largest file, blocks features |
| `manage_skills_intent.go` | 2,415 | 🔴 CRITICAL | Second largest, active development |
| `capture_event_intent.go` | 1,884 | 🔴 CRITICAL | High touch file, uses deprecated models/ |
| `burst_management_intent.go` | 1,488 | 🔴 CRITICAL | Complex state machine |
| `browse_timeline_intent.go` | 1,195 | 🔴 CRITICAL | Core feature |

**Estimate**: 2-4 hours per migration

---

### Priority 2: Deprecated Patterns

**Migrate to fix architectural violations**:

| Intent | Violation | Check |
|--------|-----------|-------|
| `capture_event_intent.go` | Uses `models.CaptureForm` | #22 |
| `manage_skills_intent.go` | Uses `models.SkillForm` | #22 |
| `burst_management_intent.go` | Excessive raw lipgloss | #21 |
| `manage_skills_intent.go` | Excessive raw lipgloss | #21 |

---

### Priority 3: Flat Structure (Warnings)

**Migrate for consistency** (7 legacy intents):

| Intent | Status | Migration Urgency |
|--------|--------|-------------------|
| `configure_system` | Flat | Medium |
| `fact_management` | Flat | Medium |
| *(others from Check #17)* | Flat | Low |

---

## File Structure Rules

### Rule 1: Subdirectory Organization (MANDATORY)

**ALWAYS use this structure**:

```
intents/{feature}/
├── context.go     # Business logic, data (100-200 lines)
├── result.go      # Output type (20-50 lines)
├── constants.go   # State enum, error constants (50-80 lines)
├── messages.go    # ALL *Msg types (50-100 lines)
└── intent.go      # Broker ONLY (200-400 lines, MAX 600)

screens/{feature}/
├── list_screen.go    # List view (if applicable)
├── detail_screen.go  # Detail view (if applicable)
└── form_screen.go    # Form view (if applicable)
```

**Enforcement**: Check #17, #18, #19, #20

---

### Rule 2: File Size Limits (STRICTLY ENFORCED)

| File | Target | Max | Enforcement |
|------|--------|-----|-------------|
| `intent.go` | 200-400 | **600** | Check #18 (BLOCKS) |
| `context.go` | 100-200 | 500 | Warning only |
| `constants.go` | 50-80 | 150 | Warning only |
| `messages.go` | 50-100 | 200 | Warning only |
| `result.go` | 20-50 | 100 | Warning only |

**If intent.go >600 lines**: BLOCKED by Check #18

---

### Rule 3: Type Location (STRICTLY ENFORCED)

**Types MUST be in correct files**:

| Type | MUST be in | NEVER in |
|------|------------|----------|
| `*Context struct` | `context.go` | `intent.go` |
| `*Result struct` | `result.go` | `intent.go` |
| `*State string` | `constants.go` | `intent.go` |
| `const (...)` | `constants.go` | `intent.go` |
| `var Err...` | `constants.go` | `intent.go` |
| `*Msg struct` | `messages.go` | `intent.go`, `constants.go` |
| `*Intent struct` | `intent.go` | ✅ Correct |

**Enforcement**: Check #20 (BLOCKS violations)

---

### Rule 4: No Rendering in Intent (STRICTLY ENFORCED)

**Intent MUST only have View() method that delegates**:

```go
// ✅ CORRECT - Only View() that delegates
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    
    return baseView
}
// Total: 1 method, ~10 lines

// ❌ WRONG - Multiple render methods
func (i *MyIntent) View() string { /* ... */ }
func (i *MyIntent) renderHeader() string { /* ... */ }
func (i *MyIntent) renderBody() string { /* ... */ }
func (i *MyIntent) renderFooter() string { /* ... */ }
func (i *MyIntent) renderList() string { /* ... */ }
func (i *MyIntent) renderDetail() string { /* ... */ }
// Total: 6+ methods, 300+ lines (VIOLATION)
```

**Max render methods**: 2 (View() + optional 1 helper)  
**Enforcement**: Check #19 (BLOCKS >2 render methods)

---

## UIKit Migration Rules

### Rule 5: No Raw Lipgloss (STRICTLY ENFORCED)

**Replace raw lipgloss with UIKit primitives**:

#### Text Styling

| Old (Raw Lipgloss) | New (UIKit) |
|-------------------|-------------|
| `lipgloss.NewStyle().Bold(true).Render("Title")` | `primitives.Title("Title", theme)` |
| `lipgloss.NewStyle().Render("Body text")` | `primitives.Body("Body text", theme)` |
| `lipgloss.NewStyle().Foreground(red).Render("Error")` | `primitives.ErrorText("Error", theme)` |
| `lipgloss.NewStyle().Foreground(green).Render("Success")` | `primitives.SuccessText("Success", theme)` |
| `lipgloss.NewStyle().Foreground(gray).Render("Muted")` | `primitives.Muted("Muted", theme)` |

#### Colors

| Old (Hardcoded) | New (Theme) |
|----------------|-------------|
| `lipgloss.Color("#ffffff")` | `theme.Primary()` |
| `lipgloss.Color("#ff0000")` | `theme.Error()` |
| `lipgloss.Color("#00ff00")` | `theme.Success()` |
| `lipgloss.Color("#888888")` | `theme.Muted()` |
| `lipgloss.Color("#5c5cff")` | `theme.Accent()` |

#### Layout

| Old (Manual) | New (UIKit) |
|-------------|-------------|
| `lipgloss.JoinVertical(lipgloss.Left, ...)` | `layout.NewScreenLayout(theme).WithContent(...)` |
| `lipgloss.Place(w, h, ...)` | `containers.NewBox(theme).Center()` |
| Manual padding/margin | `containers.NewBox(theme).Padding(...)` |

#### Components

| Old (Deprecated) | New (UIKit) |
|-----------------|-------------|
| `components.KeyBadge("key", "label")` | `primitives.HelpKeyBadge("key", "label", theme)` |
| `components.StandardView{...}` | `layout.NewScreenLayout(theme)` |
| Custom modal rendering | `behaviors.RenderModalOverlay(modal, baseView)` |

**Enforcement**: Check #21 (BLOCKS deprecated components, WARNS excessive lipgloss)

---

### Rule 6: UIKit Primitives Reference

**Use these in screens** (NOT in intents):

#### Text Primitives
```go
import "github.com/baphled/kariya/internal/cli/uikit/primitives"

// Headings and titles
primitives.Title("Main Title", theme)
primitives.Heading("Section Heading", theme)
primitives.Subheading("Subsection", theme)

// Body text
primitives.Body("Regular text", theme)
primitives.Muted("Less important text", theme)

// Status text
primitives.ErrorText("Error message", theme)
primitives.SuccessText("Success message", theme)
primitives.WarningText("Warning message", theme)
primitives.InfoText("Info message", theme)

// Help text
primitives.HelpKeyBadge("key", "description", theme)
```

#### Layout Components
```go
import "github.com/baphled/kariya/internal/cli/uikit/layout"

// Screen layout (replaces StandardView)
layout.NewScreenLayout(theme).
    WithHeader(header).
    WithContent(content).
    WithFooter(footer).
    Render()
```

#### Containers
```go
import "github.com/baphled/kariya/internal/cli/uikit/containers"

// Box (padding, borders, background)
containers.NewBox(theme).
    Content(content).
    Padding(2).
    Border(true).
    Background(theme.BackgroundColor()).
    Render()

// Overlay (for modals - MUST have solid background)
containers.NewBox(theme).
    Content(modalContent).
    Background(theme.BackgroundColor()).  // REQUIRED
    Render()
```

#### Modals (Centralized)
```go
import "github.com/baphled/kariya/internal/cli/uikit/feedback"

// Confirm modal (replaces custom delete confirmation)
deleteModal := feedback.NewConfirmModal("Delete this item?")
deleteModal.Show()

// Error modal
errorModal := feedback.NewErrorModal("Something went wrong")
errorModal.Show()

// Success modal
successModal := feedback.NewSuccessModal("Item saved successfully")
successModal.Show()

// Loading modal
loadingModal := feedback.NewLoadingModal("Loading...")
loadingModal.Show()

// Render over base view
if modal != nil && modal.IsVisible() {
    return behaviors.RenderModalOverlay(modal, baseView)
}
```

**See**: [UIKit Guide](../UIKIT_GUIDE.md)

---

## Forms Migration Rules

### Rule 7: No models/ Package for Forms (STRICTLY ENFORCED)

**DEPRECATED (DO NOT USE)**:
```go
// ❌ WRONG - models/ package
import "github.com/baphled/kariya/internal/cli/models"

type MyIntent struct {
    form *models.CaptureForm  // DEPRECATED
}

func (i *MyIntent) Init() tea.Cmd {
    i.form = models.NewCaptureForm(i.service)
    return i.form.Init()
}

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    i.form, cmd = i.form.Update(msg)
    return cmd
}

func (i *MyIntent) View() string {
    return i.form.View()
}
```

**CORRECT (USE THIS)**:
```go
// ✅ CORRECT - screens/ package
import "github.com/baphled/kariya/internal/cli/screens/myfeature"

type MyIntent struct {
    formScreen *myfeature.FormScreen  // CORRECT
}

func (i *MyIntent) Init() tea.Cmd {
    i.formScreen = myfeature.NewFormScreen(i.service)
    i.formScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.formScreen.SetTheme(i.Theme())
    i.activeScreen = i.formScreen
    return nil
}

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.formScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}

func (i *MyIntent) View() string {
    return i.activeScreen.View()
}
```

**Enforcement**: Check #22 (BLOCKS models/ usage for forms)

---

### Rule 8: Form Screen Pattern

**Form screens MUST follow this pattern**:

```go
// File: screens/myfeature/form_screen.go
package myfeature

import (
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/base"
    tea "github.com/charmbracelet/bubbletea"
)

// FormScreen handles form input.
type FormScreen struct {
    *base.BaseScreen
    form     forms.Form      // Direct use of forms package
    formData *MyFormData
}

// NewFormScreen creates a new form screen.
func NewFormScreen() *FormScreen {
    // Build form using forms package
    form := forms.NewForm(
        forms.NewGroup(
            forms.NewInput(forms.FieldConfig{
                Key:         "name",
                Title:       "Name",
                Placeholder: "Enter name",
                Required:    true,
            }),
            forms.NewText(forms.FieldConfig{
                Key:         "description",
                Title:       "Description",
                Placeholder: "Enter description",
            }),
        ),
    )
    
    return &FormScreen{
        BaseScreen: base.NewBaseScreen(),
        form:       form,
        formData:   &MyFormData{},
    }
}

// Update handles form updates.
func (s *FormScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    // Check form state via forms package
    if forms.IsCompleted(s.form) {
        // Extract data
        s.formData.Name = s.form.GetString("name")
        s.formData.Description = s.form.GetString("description")
        return nil, screens.NewSubmitResult(s.formData)
    }
    
    if forms.IsAborted(s.form) {
        return nil, screens.NewCancelResult("form")
    }
    
    // Delegate to form
    var cmd tea.Cmd
    s.form, cmd = s.form.Update(msg)
    return cmd, nil
}

// View renders the form.
func (s *FormScreen) View() string {
    return s.form.View()
}
```

**Key Points**:
- Embed `*base.BaseScreen`
- Use `forms.Form` directly (NOT `*huh.Form`)
- Check state with `forms.IsCompleted()`, `forms.IsAborted()`
- Return `ScreenResult` (Submit, Cancel, Error)

**See**: [Forms Guide](../FORMS_GUIDE.md), [Forms Workflow](../rules/FORMS_WORKFLOW_GUIDE.md)

---

## Component Migration Map

### Complete Migration Reference

| Component Type | Old Location | New Location | Pattern |
|---------------|--------------|--------------|---------|
| **Intent Logic** | `my_intent.go` (2,000 lines) | `intents/myfeature/intent.go` (300 lines) | Broker only |
| **Business Logic** | Intent file | `intents/myfeature/context.go` | Methods on Context |
| **State Enum** | Intent file | `intents/myfeature/constants.go` | Typed enum |
| **Error Constants** | Intent file | `intents/myfeature/constants.go` | `var Err...` |
| **Message Types** | Intent file | `intents/myfeature/messages.go` | All `*Msg` types |
| **Result Type** | Intent file | `intents/myfeature/result.go` | Output struct |
| **List View** | `renderList()` in intent | `screens/myfeature/list_screen.go` | TableBehavior |
| **Detail View** | `renderDetail()` in intent | `screens/myfeature/detail_screen.go` | UIKit layout |
| **Form View** | `models.*Form` | `screens/myfeature/form_screen.go` | FormScreen |
| **Delete Modal** | Custom modal in intent | `feedback.ConfirmModal` | Centralized |
| **Error Modal** | Custom modal in intent | `feedback.ErrorModal` | Centralized |
| **Success Modal** | Custom modal in intent | `feedback.SuccessModal` | Centralized |
| **Colors** | `lipgloss.Color("#xxx")` | `theme.Primary()`, etc. | Theme system |
| **Text Styling** | `lipgloss.NewStyle()` | `primitives.Title()`, etc. | UIKit primitives |
| **Layout** | `lipgloss.JoinVertical()` | `layout.NewScreenLayout()` | UIKit layout |
| **Help Keys** | `components.KeyBadge()` | `primitives.HelpKeyBadge()` | UIKit primitive |
| **Table** | Manual rendering | `behaviors.TableBehavior[T]` | Behavior pattern |

---

## Step-by-Step Process

### Phase 1: Preparation

**Before starting migration**:

1. **Read original intent completely**
   ```bash
   # Get line count
   wc -l internal/cli/intents/my_intent.go
   
   # Count render methods
   grep -c "^func.*render\|^func.*Render\|^func.*View" internal/cli/intents/my_intent.go
   
   # Identify deprecated patterns
   grep "models\.\w*Form\|components\.KeyBadge\|components\.StandardView" internal/cli/intents/my_intent.go
   ```

2. **Create migration task**
   ```bash
   make new-feature TASK="Migrate my_intent to subdirectory structure"
   ```

3. **Create feature branch**
   ```bash
   git checkout -b refactor/migrate-my-intent
   ```

---

### Phase 2: Create Directory Structure

```bash
# Create subdirectory
mkdir -p internal/cli/intents/myfeature

# Create screen directory
mkdir -p internal/cli/screens/myfeature

# Copy templates
cp examples/intent_subdirectory_template/context.go.template internal/cli/intents/myfeature/context.go
cp examples/intent_subdirectory_template/result.go.template internal/cli/intents/myfeature/result.go
cp examples/intent_subdirectory_template/constants.go.template internal/cli/intents/myfeature/constants.go
cp examples/intent_subdirectory_template/messages.go.template internal/cli/intents/myfeature/messages.go
cp examples/intent_subdirectory_template/intent.go.template internal/cli/intents/myfeature/intent.go
```

---

### Phase 3: Extract Types (Order Matters)

**3.1. Extract Constants** (`constants.go`)
- State enum
- Error constants
- Magic numbers

**3.2. Extract Messages** (`messages.go`)
- ALL `*Msg` types
- Move from intent.go or constants.go

**3.3. Extract Result** (`result.go`)
- Output type
- Return value struct

**3.4. Extract Context** (`context.go`)
- Input parameters
- Business logic methods
- Data management

**See**: [Intent Migration Guide](INTENT_MIGRATION_TO_SUBDIRECTORY.md) Step 2-5

---

### Phase 4: Extract Screens

**4.1. Identify States**
```go
// From constants.go
type MyFeatureState string

const (
    StateList   MyFeatureState = "list"    // → ListScreen
    StateDetail MyFeatureState = "detail"  // → DetailScreen
    StateForm   MyFeatureState = "form"    // → FormScreen
)
```

**4.2. Create Screens**
- `screens/myfeature/list_screen.go` - List view with TableBehavior
- `screens/myfeature/detail_screen.go` - Detail view with UIKit
- `screens/myfeature/form_screen.go` - Form view (NOT models/)

**4.3. Migrate Rendering**
- Move `renderList()` → `ListScreen.View()`
- Move `renderDetail()` → `DetailScreen.View()`
- Replace raw lipgloss with UIKit

**See**: [Screen Extraction Guide](SCREEN_EXTRACTION_GUIDE.md)

---

### Phase 5: Migrate to UIKit

**5.1. Replace Raw Lipgloss**
```bash
# Find raw lipgloss usage
grep "lipgloss.NewStyle()" screens/myfeature/*.go
```

**5.2. Apply UIKit Replacements**
- Hardcoded colors → `theme.*()` methods
- Manual styling → `primitives.*()` functions
- Manual layout → `layout.NewScreenLayout()`
- Custom modals → `feedback.*Modal`

**5.3. Update Component Imports**
```go
// Remove deprecated
- import "github.com/baphled/kariya/internal/cli/components"
- import "github.com/charmbracelet/lipgloss"

// Add UIKit
+ import "github.com/baphled/kariya/internal/cli/uikit/primitives"
+ import "github.com/baphled/kariya/internal/cli/uikit/layout"
+ import "github.com/baphled/kariya/internal/cli/uikit/feedback"
```

---

### Phase 6: Refactor Intent

**6.1. Update Intent Struct**
```go
type MyFeatureIntent struct {
    *BaseIntent
    
    // Context (input parameters)
    context *MyFeatureContext
    
    // State machine
    state  MyFeatureState
    active bool
    result *IntentResult[*MyFeatureResult]
    
    // Screens (explicit typed fields)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    formScreen   *myfeature.FormScreen
    
    // Active screen (generic pointer)
    activeScreen screens.Screen
    
    // Modals (centralized)
    deleteModal  *feedback.ConfirmModal
    errorModal   *feedback.ErrorModal
    successModal *feedback.SuccessModal
}
```

**6.2. Implement Broker Pattern**
```go
// Update delegates to screens
func (i *MyFeatureIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Modal handling
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        model, cmd := i.deleteModal.Update(msg)
        i.deleteModal = model.(*feedback.ConfirmModal)
        
        if i.deleteModal.Confirmed() {
            return i.handleDeleteConfirm()
        }
        
        return cmd
    }
    
    // Screen delegation
    cmd, result := i.activeScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

// View delegates to screens + modals
func (i *MyFeatureIntent) View() string {
    baseView := i.activeScreen.View()
    
    // Overlay modals
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**6.3. Extract Helpers**
- Business logic → `context.go`
- Rendering → `screens/`
- Keep ≤5 orchestration helpers in intent
- Utilities → `internal/util/`

**See**: [Helper Extraction Guide](HELPER_EXTRACTION_GUIDE.md)

---

### Phase 7: Delete Old File

```bash
# After migration complete and tested
git rm internal/cli/intents/my_intent.go
```

---

## Validation Checklist

### Pre-Commit Validation

**Run ALL these commands**:

```bash
# 1. Architecture checks (REQUIRED)
make check-intent-architecture

# 2. Full compliance (REQUIRED)
make check-compliance

# 3. Code formatting
make fmt

# 4. Static analysis
make vet
make staticcheck

# 5. Comprehensive linting
make golangci-lint

# 6. Tests
make test
make coverage
```

### Manual Validation Checklist

**File Structure**:
- [ ] Subdirectory created: `intents/{feature}/`
- [ ] 5 files present: context.go, result.go, constants.go, messages.go, intent.go
- [ ] Screen directory created: `screens/{feature}/`
- [ ] Screen files created (list, detail, form as needed)

**File Sizes**:
- [ ] intent.go ≤600 lines (target 200-400)
- [ ] context.go ≤500 lines (target 100-200)
- [ ] constants.go ≤150 lines (target 50-80)
- [ ] messages.go ≤200 lines (target 50-100)
- [ ] result.go ≤100 lines (target 20-50)

**Type Locations**:
- [ ] Context struct in context.go (NOT intent.go)
- [ ] Result struct in result.go (NOT intent.go)
- [ ] State enum in constants.go (NOT intent.go)
- [ ] Error constants in constants.go (NOT intent.go)
- [ ] Message types in messages.go (NOT intent.go or constants.go)
- [ ] Intent struct in intent.go

**Rendering**:
- [ ] Intent has ≤2 render methods (ideally only View())
- [ ] Rendering delegated to screens
- [ ] No raw lipgloss in intent
- [ ] UIKit components used in screens

**UIKit Migration**:
- [ ] No `lipgloss.Color("#xxx")` (use theme)
- [ ] No `lipgloss.NewStyle()` for text (use primitives)
- [ ] No `components.KeyBadge()` (use primitives.HelpKeyBadge())
- [ ] No `components.StandardView` (use layout.NewScreenLayout())
- [ ] Centralized modals (feedback.*)

**Forms**:
- [ ] No `models.*Form` usage
- [ ] Form screens in `screens/{feature}/form_screen.go`
- [ ] Forms use `forms.Form` directly (NOT `*huh.Form`)

**Architecture**:
- [ ] Intent embeds `*BaseIntent`
- [ ] Context field present
- [ ] State field present (typed enum)
- [ ] Explicit screen fields present
- [ ] Implements `ScreenResultHandler` interface
- [ ] Uses `behaviors.RenderModalOverlay()` for modals

**Tests**:
- [ ] All existing tests pass
- [ ] E2E tests updated (if needed)
- [ ] Coverage ≥95%

**Cleanup**:
- [ ] Old flat file deleted
- [ ] Imports updated in dependent files
- [ ] No compilation errors
- [ ] No linter warnings

---

## Common Violations

### Violation 1: Wrong Package Name

```go
// ❌ WRONG - Old package name
package intents

// ✅ CORRECT - Subdirectory package name
package myfeature
```

---

### Violation 2: Types in Wrong Files

```go
// ❌ WRONG - Context in intent.go
// File: intents/myfeature/intent.go
type MyFeatureContext struct { ... }  // WRONG FILE

// ✅ CORRECT - Context in context.go
// File: intents/myfeature/context.go
type MyFeatureContext struct { ... }  // CORRECT FILE
```

---

### Violation 3: Business Logic in Intent

```go
// ❌ WRONG - Business logic in intent
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // Direct DB query
    rows, err := db.Query("SELECT * FROM items")
    
    // Complex calculations
    for _, item := range items {
        // Heavy processing
    }
}

// ✅ CORRECT - Delegate to context
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    ctx := i.GetContext()
    err := i.context.LoadItems(ctx)
    // ...
}
```

---

### Violation 4: Rendering in Intent

```go
// ❌ WRONG - Rendering in intent
func (i *MyIntent) View() string {
    title := lipgloss.NewStyle().Bold(true).Render("Title")
    // 200 lines of rendering
}

// ✅ CORRECT - Delegate to screen
func (i *MyIntent) View() string {
    return i.activeScreen.View()
}
```

---

### Violation 5: Using models/ for Forms

```go
// ❌ WRONG - models/ package
import "github.com/baphled/kariya/internal/cli/models"
type MyIntent struct {
    form *models.CaptureForm
}

// ✅ CORRECT - screens/ package
import "github.com/baphled/kariya/internal/cli/screens/myfeature"
type MyIntent struct {
    formScreen *myfeature.FormScreen
}
```

---

### Violation 6: Raw Lipgloss

```go
// ❌ WRONG - Hardcoded colors
lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render("Error")

// ✅ CORRECT - Theme colors
primitives.ErrorText("Error", theme)
```

---

### Violation 7: Deprecated Components

```go
// ❌ WRONG - Deprecated components
components.KeyBadge("esc", "cancel")
components.StandardView{...}

// ✅ CORRECT - UIKit components
primitives.HelpKeyBadge("esc", "cancel", theme)
layout.NewScreenLayout(theme).WithContent(...)
```

---

## Migration Support

### Documentation
- [Intent Migration Guide](INTENT_MIGRATION_TO_SUBDIRECTORY.md) - Full migration process
- [Screen Extraction Guide](SCREEN_EXTRACTION_GUIDE.md) - Extracting rendering
- [Helper Extraction Guide](HELPER_EXTRACTION_GUIDE.md) - Extracting helpers
- [UIKit Guide](../UIKIT_GUIDE.md) - UIKit component reference
- [Forms Guide](../FORMS_GUIDE.md) - Forms architecture
- [Intent Development Checklist](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - All checks

### Templates
- `examples/intent_subdirectory_template/` - Complete template system

### Tools
- `make check-intent-architecture` - Validate migration
- `make check-compliance` - Full compliance check
- `make what-to-use NEED="keyword"` - Component lookup

### Enforcement
- **Check #17**: Subdirectory structure
- **Check #18**: File size limits (blocks >600 lines)
- **Check #19**: No rendering in intent (blocks >2 render methods)
- **Check #20**: Type location (blocks wrong files)
- **Check #21**: UIKit usage (blocks deprecated components)
- **Check #22**: No models/ for forms (blocks deprecated pattern)

---

## Summary

### Migration Must-Dos

1. ✅ Use subdirectory structure (5 files)
2. ✅ Keep intent.go ≤600 lines (target 200-400)
3. ✅ Types in correct files (Check #20)
4. ✅ No rendering in intent (delegate to screens)
5. ✅ Use UIKit components (not raw lipgloss)
6. ✅ Use screens/ for forms (not models/)
7. ✅ Centralized modals (feedback.*)
8. ✅ Run `make check-intent-architecture` before commit

### Migration Must-Nots

1. ❌ Don't keep flat structure
2. ❌ Don't exceed 600 lines in intent.go
3. ❌ Don't put types in wrong files
4. ❌ Don't render in intent (use screens)
5. ❌ Don't use raw lipgloss (use UIKit)
6. ❌ Don't use models/ for forms (use screens/)
7. ❌ Don't use deprecated components (KeyBadge, StandardView)
8. ❌ Don't skip validation checks

---

**Last Updated**: 2026-01-22  
**Status**: REQUIRED for all intent migrations  
**Enforcement**: Automated by 22 architecture checks
