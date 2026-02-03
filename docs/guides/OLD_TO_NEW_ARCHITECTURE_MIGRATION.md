# Old to New Architecture Migration Guide

**Purpose**: Complete guide for migrating from legacy flat structure to new standardized architecture  
**Status**: MANDATORY for all legacy intent migrations  
**Target**: Legacy intents using old patterns (flat structure, components/, raw lipgloss, models/)

---

## Table of Contents

1. [Overview](#overview)
2. [Old vs New Architecture](#old-vs-new-architecture)
3. [What Changed](#what-changed)
4. [Migration Phases](#migration-phases)
5. [Phase 1: Assessment](#phase-1-assessment)
6. [Phase 2: Structure Migration](#phase-2-structure-migration)
7. [Phase 3: Component Migration](#phase-3-component-migration)
8. [Phase 4: UIKit Migration](#phase-4-uikit-migration)
9. [Phase 5: Forms Migration](#phase-5-forms-migration)
10. [Phase 6: Validation](#phase-6-validation)
11. [Migration Examples](#migration-examples)
12. [Troubleshooting](#troubleshooting)

---

## Overview

### What Is "Old Architecture"?

**Old Architecture** (Legacy Pattern):
```
intents/
├── my_intent.go (2,000+ lines)          # Single flat file
│   ├── Intent struct
│   ├── All types (Context, Result, State, Msg)
│   ├── All business logic
│   ├── All render methods (8+)
│   ├── Custom modals
│   └── models.*Form usage
│
components/                               # DEPRECATED
├── key_badge.go                         # DEPRECATED
├── standard_view.go                     # DEPRECATED
└── custom_modals.go                     # DEPRECATED
│
models/                                   # DEPRECATED for forms
└── capture_form.go                      # DEPRECATED

# Raw lipgloss everywhere
lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
```

**Problems**:
- ❌ Single 2,000+ line file (unmaintainable)
- ❌ No file organization
- ❌ Rendering in intent (violates SRP)
- ❌ Using deprecated components/
- ❌ Using deprecated models/ for forms
- ❌ Raw lipgloss (hardcoded colors)
- ❌ No type safety
- ❌ Business logic mixed with UI

---

### What Is "New Architecture"?

**New Architecture** (Standardized Pattern):
```
intents/events/                           # Subdirectory structure
├── context.go (120 lines)               # Business logic only
├── result.go (20 lines)                 # Output type
├── constants.go (60 lines)              # State enum, errors
├── messages.go (80 lines)               # Message types
└── intent.go (300 lines)                # Broker only

screens/events/                           # Screen components
├── list_screen.go                       # List view
├── detail_screen.go                     # Detail view
└── form_screen.go                       # Form view (NOT models/)

# UIKit components
uikit/
├── primitives/                          # Text, buttons, badges
├── layout/                              # Screen layouts
├── feedback/                            # Centralized modals
└── theme/                               # Theme system

# Theme-based colors
theme.Primary(), theme.Error(), theme.Success()
```

**Benefits**:
- ✅ 5-file subdirectory structure (organized)
- ✅ Intent ≤600 lines (maintainable)
- ✅ Screens handle rendering (SRP)
- ✅ UIKit components (consistent)
- ✅ Centralized modals (reusable)
- ✅ Theme system (no hardcoded colors)
- ✅ Type-safe (typed enums)
- ✅ Business logic isolated (testable)

---

## Old vs New Architecture

### Visual Comparison

**OLD (Flat Structure)**:
```
┌─────────────────────────────────────────────────┐
│ my_intent.go (2,746 lines)                     │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ Types (Context, Result, State, Msgs)   │   │
│  │ Business Logic (validate, filter, etc) │   │
│  │ Rendering (8+ render methods)          │   │
│  │ Custom Modals                           │   │
│  │ Raw Lipgloss Styling                    │   │
│  │ models.*Form usage                      │   │
│  │ Intent Implementation                   │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  Everything in one giant file                  │
└─────────────────────────────────────────────────┘
```

**NEW (Subdirectory Structure)**:
```
┌─────────────────────────────────────────────────┐
│ intents/events/                                 │
│  ├── context.go (120 lines)                    │
│  │   └── Business logic only                   │
│  ├── result.go (20 lines)                      │
│  │   └── Output type                           │
│  ├── constants.go (60 lines)                   │
│  │   └── State enum, errors                    │
│  ├── messages.go (80 lines)                    │
│  │   └── Message types                         │
│  └── intent.go (300 lines)                     │
│      └── Broker/orchestration only             │
│                                                 │
│ screens/events/                                 │
│  ├── list_screen.go                            │
│  │   └── List view with TableBehavior          │
│  ├── detail_screen.go                          │
│  │   └── Detail view with UIKit                │
│  └── form_screen.go                            │
│      └── Form view (NOT models/)               │
│                                                 │
│ uikit/feedback/                                 │
│  └── Centralized modals                        │
│                                                 │
│  Clean separation, organized, maintainable     │
└─────────────────────────────────────────────────┘
```

---

## What Changed

### 1. File Structure

| Old | New | Why |
|-----|-----|-----|
| Single file (2,000+ lines) | 5 files (300-600 lines total) | Maintainability |
| Flat structure | Subdirectory structure | Organization |
| All types in one file | Types in separate files | Clarity |
| Everything mixed | Clear separation | SRP compliance |

**Enforcement**: Check #17 (subdirectory), Check #18 (file size)

---

### 2. Intent Role

| Old | New | Why |
|-----|-----|-----|
| Intent does everything | Intent orchestrates only | Separation of concerns |
| Business logic in intent | Business logic in context.go | Testability |
| Rendering in intent | Rendering in screens/ | SRP compliance |
| 8+ render methods | 1 View() method | Simplicity |

**Enforcement**: Check #19 (≤2 render methods)

---

### 3. Components

| Old | New | Why |
|-----|-----|-----|
| `components.KeyBadge()` | `primitives.HelpKeyBadge()` | Consistent naming |
| `components.StandardView` | `layout.NewScreenLayout()` | UIKit architecture |
| Custom modals | `feedback.*Modal` | Centralized, reusable |
| Raw lipgloss | UIKit primitives | Consistency |

**Enforcement**: Check #21 (UIKit usage)

---

### 4. Styling

| Old | New | Why |
|-----|-----|-----|
| `lipgloss.Color("#fff")` | `theme.Primary()` | Theme support |
| `lipgloss.NewStyle()...` | `primitives.Title()` | Consistency |
| Hardcoded colors | Theme system | Customizable |
| Manual layout | `layout.NewScreenLayout()` | Reusable patterns |

**Enforcement**: Check #21 (warns on excessive lipgloss)

---

### 5. Forms

| Old | New | Why |
|-----|-----|-----|
| `models.CaptureForm` | `base.FormScreen[T]` + `screens/events/FormScreen` | Consistent pattern |
| `models/` package | `screens/` package with `base.FormScreen[T]` | Architecture alignment |
| Wrapper model (manual) | Generic base (automatic resize, StandardView) | Less boilerplate |
| Intent → models → forms | Intent → screens → base.FormScreen → forms | Clean separation |

**Enforcement**: Check #22 (blocks models/ for forms)

---

### 6. Type Safety

| Old | New | Why |
|-----|-----|-----|
| `state string` | `state EventsState` (typed) | Type safety |
| Types scattered | Types in correct files | Organization |
| No enum | Typed enum constants | Compile-time checks |
| Magic strings | Typed constants | Refactoring safety |

**Enforcement**: Check #1 (typed state), Check #20 (type location)

---

## Migration Phases

### Complete Migration Workflow

```
Phase 1: Assessment
    ↓
Phase 2: Structure Migration (create subdirectory)
    ↓
Phase 3: Component Migration (deprecated → UIKit)
    ↓
Phase 4: UIKit Migration (raw lipgloss → theme)
    ↓
Phase 5: Forms Migration (models/ → screens/)
    ↓
Phase 6: Validation (automated checks)
```

**Estimated Time**: 2-4 hours per intent (depending on size)

**Expected Reduction**: 48-87% code reduction

---

## Phase 1: Assessment

### Step 1: Analyze Current State

```bash
# 1. Check file size
wc -l internal/cli/intents/my_intent.go

# 2. Count render methods
grep -c "^func.*render\|^func.*Render\|^func.*View" internal/cli/intents/my_intent.go

# 3. Check for deprecated patterns
grep "components\.KeyBadge\|components\.StandardView" internal/cli/intents/my_intent.go
grep "models\.\w*Form" internal/cli/intents/my_intent.go
grep "lipgloss\.Color" internal/cli/intents/my_intent.go | wc -l

# 4. Identify types to extract
grep "^type.*struct\|^type.*string\|^const" internal/cli/intents/my_intent.go

# 5. Run checks to see violations
make check-intent-architecture
```

**Output Analysis**:
```
File size: 2,746 lines             # Violation (Check #18)
Render methods: 8                  # Violation (Check #19)
Deprecated components: 5           # Violation (Check #21)
models.Form usage: 2               # Violation (Check #22)
Raw lipgloss: 147                  # Warning (Check #21)
```

**Decision**: MUST migrate (multiple blocking violations)

---

### Step 2: Create Migration Plan

**Template**:
```markdown
# Migration Plan: {intent_name}

## Current State
- File: intents/{intent_name}.go
- Lines: {count}
- Render methods: {count}
- Violations: {list}

## Target State
- Directory: intents/{feature}/
- Files: context.go, result.go, constants.go, messages.go, intent.go
- Screens: screens/{feature}/*.go
- Target lines: ~300-600 total

## Migration Steps
1. [ ] Create subdirectory structure
2. [ ] Extract constants (state enum, errors)
3. [ ] Extract messages (all *Msg types)
4. [ ] Extract result (output type)
5. [ ] Extract context (business logic)
6. [ ] Create screens (list, detail, form)
7. [ ] Migrate to UIKit (deprecated → UIKit)
8. [ ] Migrate forms (models/ → screens/)
9. [ ] Refactor intent (broker only)
10. [ ] Validate (all checks pass)

## Expected Reduction
- From: {original_lines} lines
- To: ~{target_lines} lines
- Reduction: ~{percentage}%
```

---

## Phase 2: Structure Migration

### Step 1: Create Directory Structure

```bash
# 1. Create subdirectory
mkdir -p internal/cli/intents/{feature}

# 2. Create screen directory
mkdir -p internal/cli/screens/{feature}

# 3. Copy templates
cp examples/intent_subdirectory_template/context.go.template internal/cli/intents/{feature}/context.go
cp examples/intent_subdirectory_template/result.go.template internal/cli/intents/{feature}/result.go
cp examples/intent_subdirectory_template/constants.go.template internal/cli/intents/{feature}/constants.go
cp examples/intent_subdirectory_template/messages.go.template internal/cli/intents/{feature}/messages.go
cp examples/intent_subdirectory_template/intent.go.template internal/cli/intents/{feature}/intent.go
cp examples/intent_subdirectory_template/screens/list_screen.go.template internal/cli/screens/{feature}/list_screen.go
```

---

### Step 2: Extract Constants

**From old file**:
```go
// File: my_intent.go (OLD)
type MyIntentState string

const (
    StateInitial MyIntentState = "initial"
    StateList    MyIntentState = "list"
    StateDetail  MyIntentState = "detail"
    StateForm    MyIntentState = "form"
)

var (
    ErrNoItems = errors.New("no items found")
    ErrInvalid = errors.New("invalid input")
)
```

**To new file**:
```go
// File: intents/{feature}/constants.go (NEW)
package {feature}

import "errors"

// {Feature}State defines the state machine states.
type {Feature}State string

const (
    // StateInitial is the initial state.
    StateInitial {Feature}State = "initial"
    
    // StateList is the list view state.
    StateList {Feature}State = "list"
    
    // StateDetail is the detail view state.
    StateDetail {Feature}State = "detail"
    
    // StateForm is the form view state.
    StateForm {Feature}State = "form"
)

// Error constants.
var (
    ErrNoItems = errors.New("no items found")
    ErrInvalid = errors.New("invalid input")
)
```

**Checklist**:
- [ ] Copy state enum type definition
- [ ] Copy all state constants
- [ ] Copy all error constants
- [ ] Add godoc comments
- [ ] Update package name
- [ ] Remove from old file

---

### Step 3: Extract Messages

**From old file**:
```go
// File: my_intent.go (OLD)
type ItemSelectedMsg struct {
    Item *Item
}

type ItemCreatedMsg struct {
    Item *Item
}

type ItemDeletedMsg struct {
    ID string
}
```

**To new file**:
```go
// File: intents/{feature}/messages.go (NEW)
package {feature}

import "github.com/baphled/kariya/internal/domain"

// ItemSelectedMsg is sent when an item is selected.
type ItemSelectedMsg struct {
    Item *domain.Item
}

// ItemCreatedMsg is sent when an item is created.
type ItemCreatedMsg struct {
    Item *domain.Item
}

// ItemDeletedMsg is sent when an item is deleted.
type ItemDeletedMsg struct {
    ID string
}
```

**Checklist**:
- [ ] Copy ALL *Msg types
- [ ] Add godoc comments
- [ ] Update package name
- [ ] Update imports
- [ ] Remove from old file

---

### Step 4: Extract Result

**From old file**:
```go
// File: my_intent.go (OLD)
type MyIntentResult struct {
    Items   []*Item
    Created *Item
    Updated *Item
    Deleted string
}
```

**To new file**:
```go
// File: intents/{feature}/result.go (NEW)
package {feature}

import "github.com/baphled/kariya/internal/domain"

// {Feature}Result is the output of the {feature} intent.
type {Feature}Result struct {
    // Items is the list of items (for list mode).
    Items []*domain.Item
    
    // Created is the newly created item.
    Created *domain.Item
    
    // Updated is the updated item.
    Updated *domain.Item
    
    // Deleted is the ID of the deleted item.
    Deleted string
}
```

**Checklist**:
- [ ] Copy result type
- [ ] Add godoc comments
- [ ] Update package name
- [ ] Update imports
- [ ] Remove from old file

---

### Step 5: Extract Context

**From old file**:
```go
// File: my_intent.go (OLD)
type MyIntent struct {
    *BaseIntent
    service *service.MyService
    items   []*Item
    // ... other fields
}

func (i *MyIntent) validateItem(item *Item) error { ... }
func (i *MyIntent) filterItems(criteria *Filter) []*Item { ... }
func (i *MyIntent) loadItems() error { ... }
```

**To new file**:
```go
// File: intents/{feature}/context.go (NEW)
package {feature}

import (
    "context"
    "github.com/baphled/kariya/internal/domain"
    "github.com/baphled/kariya/internal/service"
)

// {Feature}Context holds input parameters and business logic.
type {Feature}Context struct {
    // Input parameters
    Service *service.MyService
    
    // Data
    items         []*domain.Item
    selectedItem  *domain.Item
}

// New{Feature}Context creates a new context.
func New{Feature}Context(service *service.MyService) *{Feature}Context {
    return &{Feature}Context{
        Service: service,
    }
}

// ValidateItem validates an item.
func (c *{Feature}Context) ValidateItem(item *domain.Item) error {
    // ... validation logic
}

// FilterItems filters items by criteria.
func (c *{Feature}Context) FilterItems(criteria *Filter) []*domain.Item {
    // ... filter logic
}

// LoadItems loads all items from the service.
func (c *{Feature}Context) LoadItems(ctx context.Context) error {
    items, err := c.Service.GetAll(ctx)
    if err != nil {
        return err
    }
    c.items = items
    return nil
}

// GetItems returns all items.
func (c *{Feature}Context) GetItems() []*domain.Item {
    return c.items
}

// GetSelectedItem returns the selected item.
func (c *{Feature}Context) GetSelectedItem() *domain.Item {
    return c.selectedItem
}

// SetSelectedItem sets the selected item.
func (c *{Feature}Context) SetSelectedItem(item *domain.Item) {
    c.selectedItem = item
}
```

**Checklist**:
- [ ] Copy input parameters to Context struct
- [ ] Copy data fields to Context struct
- [ ] Move business logic methods to Context
- [ ] Add getter/setter methods
- [ ] Add godoc comments
- [ ] Update package name
- [ ] Remove from old file

---

## Phase 3: Component Migration

### Deprecated Components → UIKit

**Migration Map**:

| Old (DEPRECATED) | New (UIKit) | Location |
|-----------------|-------------|----------|
| `components.KeyBadge("k", "label")` | `primitives.HelpKeyBadge("k", "label", theme)` | `uikit/primitives/` |
| `components.StandardView{...}` | `layout.NewScreenLayout(theme)` | `uikit/layout/` |
| Custom `DeleteModal` | `feedback.NewConfirmModal()` | `uikit/feedback/` |
| Custom `ErrorModal` | `feedback.NewErrorModal()` | `uikit/feedback/` |
| Custom `SuccessModal` | `feedback.NewSuccessModal()` | `uikit/feedback/` |

---

### Example: KeyBadge Migration

**OLD**:
```go
// File: my_intent.go
import "github.com/baphled/kariya/internal/cli/components"

func (i *MyIntent) renderHelp() string {
    help := components.KeyBadge("enter", "select") + " " +
            components.KeyBadge("esc", "cancel")
    return help
}
```

**NEW**:
```go
// File: screens/{feature}/list_screen.go
import "github.com/baphled/kariya/internal/cli/uikit/primitives"

func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("esc", "cancel", theme)
    
    return layout.NewScreenLayout(theme).
        WithFooter(help).
        Render()
}
```

---

### Example: StandardView Migration

**OLD**:
```go
import "github.com/baphled/kariya/internal/cli/components"

func (i *MyIntent) View() string {
    return components.StandardView{
        Title:   "My Feature",
        Content: bodyContent,
        Help:    helpText,
    }.Render()
}
```

**NEW**:
```go
import "github.com/baphled/kariya/internal/cli/uikit/layout"
import "github.com/baphled/kariya/internal/cli/uikit/primitives"

func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    header := primitives.Title("My Feature", theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(bodyContent).
        WithFooter(helpText).
        Render()
}
```

---

### Example: Custom Modal Migration

**OLD**:
```go
// File: my_intent.go
type DeleteConfirmModal struct {
    visible bool
    message string
    // ... 100 lines of custom modal code
}

func (i *MyIntent) showDeleteConfirm() {
    i.deleteModal = &DeleteConfirmModal{
        visible: true,
        message: "Delete this item?",
    }
}

func (i *MyIntent) View() string {
    if i.deleteModal != nil && i.deleteModal.visible {
        // 50 lines of custom modal rendering
    }
    return baseView
}
```

**NEW**:
```go
// File: intents/{feature}/intent.go
import "github.com/baphled/kariya/internal/cli/uikit/feedback"
import "github.com/baphled/kariya/internal/cli/behaviors"

type {Feature}Intent struct {
    *BaseIntent
    deleteModal *feedback.ConfirmModal
}

func (i *{Feature}Intent) showDeleteConfirm(item *Item) {
    i.deleteModal = feedback.NewConfirmModal(
        fmt.Sprintf("Delete '%s'?", item.Name),
        feedback.WithDanger(),
    )
    i.deleteModal.Show()
}

func (i *{Feature}Intent) View() string {
    baseView := i.activeScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**Reduction**: 150 lines → 10 lines (93% reduction!)

---

## Phase 4: UIKit Migration

### Raw Lipgloss → Theme System

**Migration Map**:

| Old (Hardcoded) | New (Theme) |
|----------------|-------------|
| `lipgloss.Color("#ffffff")` | `theme.Primary()` |
| `lipgloss.Color("#ff0000")` | `theme.Error()` |
| `lipgloss.Color("#00ff00")` | `theme.Success()` |
| `lipgloss.Color("#888888")` | `theme.Muted()` |
| `lipgloss.Color("#5c5cff")` | `theme.Accent()` |
| `lipgloss.Color("#ffaa00")` | `theme.Warning()` |

---

### Raw Styling → UIKit Primitives

**Migration Map**:

| Old (Raw Lipgloss) | New (UIKit) |
|-------------------|-------------|
| `lipgloss.NewStyle().Bold(true).Render("Title")` | `primitives.Title("Title", theme)` |
| `lipgloss.NewStyle().Render("Body")` | `primitives.Body("Body", theme)` |
| `lipgloss.NewStyle().Foreground(red).Render("Error")` | `primitives.ErrorText("Error", theme)` |
| `lipgloss.NewStyle().Foreground(green).Render("OK")` | `primitives.SuccessText("OK", theme)` |
| `lipgloss.NewStyle().Foreground(gray).Render("Muted")` | `primitives.Muted("Muted", theme)` |

---

### Example: Complete Styling Migration

**OLD** (Raw Lipgloss):
```go
// 150 lines of styling code
func (i *MyIntent) renderItem(item *Item) string {
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        Background(lipgloss.Color("#5c5cff")).
        Bold(true).
        Padding(1, 2)
    
    bodyStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#dddddd")).
        Padding(1).
        Width(60)
    
    dateStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#888888")).
        Italic(true)
    
    title := titleStyle.Render(item.Title)
    body := bodyStyle.Render(item.Description)
    date := dateStyle.Render(item.Date.Format("2006-01-02"))
    
    return lipgloss.JoinVertical(lipgloss.Left, title, body, date)
}
```

**NEW** (UIKit):
```go
// 15 lines of styling code (90% reduction!)
func (s *DetailScreen) View() string {
    theme := s.GetTheme()
    
    title := primitives.Title(s.item.Title, theme)
    body := primitives.Body(s.item.Description, theme)
    date := primitives.Muted(s.item.Date.Format("2006-01-02"), theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(title).
        WithContent(body).
        WithFooter(date).
        Render()
}
```

---

## Phase 5: Forms Migration

### models/ Package → screens/ Package

**OLD Pattern** (DEPRECATED):
```go
// File: my_intent.go
import "github.com/baphled/kariya/internal/cli/models"

type MyIntent struct {
    *BaseIntent
    form *models.CaptureForm  // DEPRECATED
}

func (i *MyIntent) Init() tea.Cmd {
    i.form = models.NewCaptureForm(i.service)
    return i.form.Init()
}

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    i.form, cmd = i.form.Update(msg)
    
    if i.form.Submitted() {
        // handle submission
    }
    
    return cmd
}

func (i *MyIntent) View() string {
    return i.form.View()
}
```

**NEW Pattern** (CORRECT):
```go
// File: intents/{feature}/intent.go
import "github.com/baphled/kariya/internal/cli/screens/{feature}"

type {Feature}Intent struct {
    *BaseIntent
    formScreen *{feature}.FormScreen  // CORRECT
}

func (i *{Feature}Intent) Init() tea.Cmd {
    i.formScreen = {feature}.NewFormScreen()
    i.formScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.formScreen.SetTheme(i.Theme())
    i.activeScreen = i.formScreen
    i.state = StateForm
    return nil
}

func (i *{Feature}Intent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.formScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

func (i *{Feature}Intent) View() string {
    return i.activeScreen.View()
}
```

---

### Create FormScreen

**NEW File**: `screens/{feature}/form_screen.go`

Uses `base.FormScreen[T]` which handles `Update()`, `View()`, window resize, and escape
cancellation automatically. `FormData` types and form builders live in `internal/cli/forms/`.

```go
package {feature}

import (
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/cli/screens/base"
)

// FormScreen handles form input.
type FormScreen struct {
    *base.FormScreen[*forms.{Feature}FormData]
}

// NewFormScreen creates a new form screen.
func NewFormScreen() *FormScreen {
    formData := &forms.{Feature}FormData{}
    baseScreen := base.NewBaseFormScreen(
        []string{"Main Menu", "{Feature}", "Create"},
        forms.New{Feature}FormWithDataAndDimensions, // FormBuilder[T]
        formData,
    )
    return &FormScreen{FormScreen: baseScreen}
}
```

The `base.FormScreen[T]` provides:
- Automatic `Update()` that delegates to `forms.Update()` and checks `forms.IsCompleted()`
- Automatic `View()` using `CreateView()` with breadcrumbs and footer
- Window resize handling via `SetTerminalInfo()` that rebuilds the form
- Returns `*screens.SubmitResult` (with `FormData` field) or `*screens.CancelResult`

**Real example**: See `internal/cli/screens/skills/skill_form.go` (93 lines).

---

## Phase 6: Validation

### Step 1: Run Automated Checks

```bash
# 1. Architecture checks (REQUIRED)
make check-intent-architecture

# Expected: All checks pass
# - Check #17: Subdirectory structure ✅
# - Check #18: File size ≤600 lines ✅
# - Check #19: ≤2 render methods ✅
# - Check #20: Types in correct files ✅
# - Check #21: UIKit usage ✅
# - Check #22: No models/ for forms ✅

# 2. Full compliance
make check-compliance

# 3. Tests
make test
make coverage

# Expected: All tests pass, coverage ≥95%
```

---

### Step 2: Manual Validation

**Checklist**:
- [ ] Subdirectory structure created
- [ ] 5 files present (context, result, constants, messages, intent)
- [ ] Screen files created
- [ ] Intent.go ≤600 lines
- [ ] Context.go ≤500 lines
- [ ] Types in correct files
- [ ] No raw lipgloss in intent
- [ ] UIKit components used
- [ ] No deprecated components (KeyBadge, StandardView)
- [ ] No models/ for forms
- [ ] Centralized modals used
- [ ] All tests pass
- [ ] Old file deleted

---

### Step 3: Delete Old File

```bash
# Only after all checks pass!
git rm internal/cli/intents/my_intent.go

# Commit migration
git add internal/cli/intents/{feature}/
git add internal/cli/screens/{feature}/
git commit -m "refactor({feature}): migrate to new architecture

- Migrate from flat structure to subdirectory
- Extract types to separate files
- Migrate components/ to UIKit
- Migrate raw lipgloss to theme system
- Migrate models/ to screens/ for forms
- Reduce from 2,746 lines to 580 lines (79% reduction)

Closes #XXX"
```

---

## Migration Examples

### Example 1: Complete Migration (capture_event_intent.go)

**BEFORE**:
```
internal/cli/intents/capture_event_intent.go (1,884 lines)
├── CaptureEventIntent struct
├── CaptureEventState string
├── const (StateList, StateForm, ...)
├── CaptureEventResult struct
├── ItemSelectedMsg struct
├── Business logic (validate, filter, etc.)
├── 8 render methods
├── Custom modals
├── models.CaptureForm usage
└── Raw lipgloss styling
```

**AFTER**:
```
internal/cli/intents/events/
├── context.go (150 lines)
│   └── Business logic
├── result.go (25 lines)
│   └── CaptureEventResult struct
├── constants.go (70 lines)
│   └── CaptureEventState enum + errors
├── messages.go (90 lines)
│   └── ItemSelectedMsg, etc.
└── intent.go (320 lines)
    └── Broker/orchestration only

internal/cli/screens/events/
├── list_screen.go (100 lines)
├── detail_screen.go (90 lines)
└── form_screen.go (110 lines)

Total: 955 lines (from 1,884)
Reduction: 929 lines (49% reduction!)
```

---

### Example 2: UIKit Migration (manage_skills_intent.go)

**BEFORE** (Raw Lipgloss):
```go
// 147 instances of raw lipgloss.Color()
titleStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#ffffff")).
    Background(lipgloss.Color("#5c5cff")).
    Bold(true)

errorStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#ff0000"))

successStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#00ff00"))
```

**AFTER** (Theme + UIKit):
```go
// 0 instances of raw lipgloss.Color()
theme := s.GetTheme()

title := primitives.Title("Manage Skills", theme)
error := primitives.ErrorText("Failed to load", theme)
success := primitives.SuccessText("Saved successfully", theme)
```

**Result**: 147 color instances → 0 (100% migrated to theme)

---

### Example 3: Forms Migration (capture_event_intent.go)

**BEFORE** (models/ Package):
```go
import "github.com/baphled/kariya/internal/cli/models"

type CaptureEventIntent struct {
    form *models.CaptureForm  // DEPRECATED
}

func (i *CaptureEventIntent) Init() tea.Cmd {
    i.form = models.NewCaptureForm(i.service)
    return i.form.Init()
}
```

**AFTER** (screens/ Package):
```go
import "github.com/baphled/kariya/internal/cli/screens/events"

type CaptureEventIntent struct {
    formScreen *events.FormScreen  // CORRECT
}

func (i *CaptureEventIntent) Init() tea.Cmd {
    i.formScreen = events.NewFormScreen()
    i.formScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.formScreen.SetTheme(i.Theme())
    i.activeScreen = i.formScreen
    return nil
}
```

**Result**: Check #22 violation fixed ✅

---

## Troubleshooting

### Issue 1: Import Cycle

**Error**:
```
import cycle not allowed
package github.com/baphled/kariya/internal/cli/screens/events
    imports github.com/baphled/kariya/internal/cli/intents/events
    imports github.com/baphled/kariya/internal/cli/screens/events
```

**Cause**: Screen importing intents package

**Fix**: Screens NEVER import intents
```go
// ❌ WRONG
import "github.com/baphled/kariya/internal/cli/intents/events"

// ✅ CORRECT
// Don't import intents, use ScreenResult to communicate
return nil, screens.NewNavigateResult("detail", item)
```

---

### Issue 2: Tests Failing After Migration

**Error**:
```
undefined: MyIntent
undefined: NewMyIntent
```

**Cause**: Old imports still referencing flat structure

**Fix**: Update test imports
```go
// OLD
import "github.com/baphled/kariya/internal/cli/intents"
intent := intents.NewMyIntent()

// NEW
import myfeature "github.com/baphled/kariya/internal/cli/intents/myfeature"
intent := myfeature.NewMyFeatureIntent()
```

---

### Issue 3: Type Not Found

**Error**:
```
undefined: MyIntentState
undefined: StateList
```

**Cause**: Types not imported from correct file

**Fix**: Import from subdirectory package
```go
// OLD (flat)
import "github.com/baphled/kariya/internal/cli/intents"

// NEW (subdirectory)
import myfeature "github.com/baphled/kariya/internal/cli/intents/myfeature"

state := myfeature.StateList
```

---

### Issue 4: Check #22 Still Failing

**Error**:
```
❌ VIOLATION: Deprecated models/ package for forms
File: internal/cli/intents/events/intent.go
```

**Cause**: Still using models.*Form

**Fix**: Change to FormScreen
```go
// Remove this
import "github.com/baphled/kariya/internal/cli/models"
form *models.CaptureForm

// Add this
import "github.com/baphled/kariya/internal/cli/screens/events"
formScreen *events.FormScreen
```

---

### Issue 5: Check #21 Warning (Excessive Lipgloss)

**Warning**:
```
⚠️  WARNING: Excessive raw lipgloss usage
File: internal/cli/screens/events/list_screen.go (15 instances)
```

**Fix**: Replace with UIKit
```go
// Remove raw lipgloss
lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

// Add UIKit/theme
primitives.Title("text", theme)
theme.Primary()
```

---

## Summary

### Migration Workflow

```
Old Architecture (2,000+ lines)
    ↓
Phase 1: Assessment
    ↓
Phase 2: Structure Migration (subdirectory)
    ↓
Phase 3: Component Migration (deprecated → UIKit)
    ↓
Phase 4: UIKit Migration (lipgloss → theme)
    ↓
Phase 5: Forms Migration (models/ → screens/)
    ↓
Phase 6: Validation (all checks pass)
    ↓
New Architecture (300-600 lines, 48-87% reduction)
```

### Key Changes Summary

| Aspect | Old | New |
|--------|-----|-----|
| Structure | Flat (1 file) | Subdirectory (5 files) |
| Intent role | Does everything | Orchestrates only |
| Rendering | In intent (8+ methods) | In screens (1 method) |
| Components | Deprecated components/ | UIKit components |
| Styling | Raw lipgloss | Theme system |
| Forms | models/ package | screens/ package |
| File size | 2,000+ lines | 300-600 lines |
| Organization | Mixed | Separated |
| Maintainability | Low | High |

### Expected Results

**Code Reduction**: 48-87%  
**Violations**: Many → Zero  
**Maintainability**: Low → High  
**Consistency**: None → Full  
**Type Safety**: Weak → Strong  

---

## Resources

### Documentation
- [Migration Rules](MIGRATION_RULES.md) - Complete migration reference
- [View Extraction Rules](../rules/VIEW_EXTRACTION_RULES.md) - View extraction
- [Component Usage Rules](../rules/COMPONENT_USAGE_RULES.md) - Component selection
- [Intent Development Checklist](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - All checks
- [UIKit Guide](../UIKIT_GUIDE.md) - UIKit components
- [Forms Guide](../FORMS_GUIDE.md) - Forms architecture

### Templates
- `examples/intent_subdirectory_template/` - Complete templates

### Tools
- `make check-intent-architecture` - Validate migration
- `make what-to-use NEED="keyword"` - Component lookup

### Enforcement
- Check #17: Subdirectory structure
- Check #18: File size ≤600 lines
- Check #19: ≤2 render methods
- Check #20: Type location
- Check #21: UIKit usage
- Check #22: No models/ for forms

---

**Last Updated**: 2026-01-22  
**Status**: MANDATORY for all legacy intent migrations  
**Enforcement**: Automated by 22 architecture checks
