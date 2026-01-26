# Component Usage Rules: When to Use What

**Purpose**: Decision tree for choosing the right component type  
**Status**: MANDATORY - Must follow for all TUI development  
**Enforcement**: Automated checks + code review

---

## Table of Contents

1. [Quick Decision Tree](#quick-decision-tree)
2. [Modals vs Screens](#modals-vs-screens)
3. [When to Use Modals](#when-to-use-modals)
4. [When to Use Screens](#when-to-use-screens)
5. [When to Use Forms](#when-to-use-forms)
6. [When to Use Behaviors](#when-to-use-behaviors)
7. [Component Selection Matrix](#component-selection-matrix)
8. [Common Scenarios](#common-scenarios)
9. [Anti-Patterns](#anti-patterns)

---

## Quick Decision Tree

```
Need to show UI?
├─ Temporary overlay (2-5 seconds)?
│  └─ YES → Use Modal (feedback/*)
│     └─ Confirmation? → ConfirmModal
│     └─ Error message? → ErrorModal
│     └─ Success message? → SuccessModal
│     └─ Loading? → LoadingModal
│     └─ Info? → InfoModal
│
├─ Capture data (form inputs)?
│  └─ YES → Use FormScreen (screens/{feature}/form_screen.go)
│     └─ Embed forms.Form
│     └─ Return ScreenResult on submit/cancel
│
├─ Display list of items?
│  └─ YES → Use ListScreen with TableBehavior
│     └─ screens/{feature}/list_screen.go
│     └─ Embed behaviors.TableBehavior[T]
│
├─ Display single item details?
│  └─ YES → Use DetailScreen
│     └─ screens/{feature}/detail_screen.go
│     └─ Use UIKit layout
│
└─ Complex interaction (CRUD, filters, etc.)?
   └─ YES → Use Multiple Screens + Behaviors
      └─ State machine in intent
      └─ Each state = different screen
```

---

## Modals vs Screens

### Key Differences

| Aspect | Modal | Screen |
|--------|-------|--------|
| **Purpose** | Temporary feedback/confirmation | Primary content/interaction |
| **Duration** | 2-5 seconds (auto-dismiss or quick action) | Until user navigates away |
| **Interaction** | Minimal (Y/N, OK, dismiss) | Full (forms, lists, navigation) |
| **Rendering** | Overlay on top of screen | Full screen content |
| **State** | Transient | Persistent (part of state machine) |
| **Location** | `uikit/feedback/*Modal` (centralized) | `screens/{feature}/*_screen.go` |
| **Example** | "Are you sure?", "Success!", "Loading..." | List view, detail view, form |

### Visual Distinction

**Modal (Overlay)**:
```
┌────────────────────────────────┐
│ Base Screen Content            │
│                                │
│    ┌──────────────────┐       │
│    │  Modal Overlay   │       │  ← Appears on top
│    │  [Yes] [No]      │       │
│    └──────────────────┘       │
│                                │
│                                │
└────────────────────────────────┘
```

**Screen (Full Content)**:
```
┌────────────────────────────────┐
│ Screen Header                  │
├────────────────────────────────┤
│                                │
│ Screen Content                 │  ← Takes full area
│ (List, Form, Detail, etc.)     │
│                                │
├────────────────────────────────┤
│ Footer (help keys)             │
└────────────────────────────────┘
```

---

## When to Use Modals

### Use Modals When

1. **Confirmation Required** (2-3 seconds)
   - Delete confirmation
   - Destructive action warning
   - Overwrite confirmation
   - Discard changes prompt

2. **Feedback Message** (2-3 seconds)
   - Success notification
   - Error message
   - Warning alert
   - Info notification

3. **Temporary Status** (while processing)
   - Loading indicator
   - Processing status
   - Waiting for async operation

4. **Quick Input** (single field, <5 seconds)
   - Enter password
   - Confirm name
   - Quick yes/no decision

### Modal Types Available

**Location**: `internal/cli/uikit/feedback/`

| Modal | Use Case | Duration | Input |
|-------|----------|----------|-------|
| **ConfirmModal** | Yes/No questions | Until user responds | Y/N |
| **ErrorModal** | Show errors | Until dismissed | None (auto-dismiss or key) |
| **SuccessModal** | Show success | Until dismissed | None (auto-dismiss or key) |
| **WarningModal** | Show warnings | Until dismissed | None (auto-dismiss or key) |
| **InfoModal** | Show information | Until dismissed | None (auto-dismiss or key) |
| **LoadingModal** | Show loading state | Until operation completes | None |

### Modal Pattern

```go
// 1. Create modal
i.deleteModal = feedback.NewConfirmModal(
    "Delete this item?",
    feedback.WithDanger(),  // Optional styling
)
i.deleteModal.Show()

// 2. Handle in Update()
if i.deleteModal != nil && i.deleteModal.IsVisible() {
    model, cmd := i.deleteModal.Update(msg)
    i.deleteModal = model.(*feedback.ConfirmModal)
    
    if i.deleteModal.Confirmed() {
        return i.handleDeleteConfirm()
    }
    if i.deleteModal.Cancelled() {
        i.deleteModal.Hide()
    }
    
    return cmd
}

// 3. Render in View()
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

### Modal Rules

**DO**:
- ✅ Use centralized modals from `uikit/feedback/`
- ✅ Use `behaviors.RenderModalOverlay()` for rendering
- ✅ Keep modal interactions simple (Y/N, OK, dismiss)
- ✅ Auto-dismiss after action (if appropriate)
- ✅ Show modal on top of current screen

**DON'T**:
- ❌ Create custom modal types (use centralized)
- ❌ Use modals for complex forms (use FormScreen)
- ❌ Use modals for list navigation (use ListScreen)
- ❌ Keep modals visible for >10 seconds (use Screen)
- ❌ Put business logic in modals (delegate to context)

---

## When to Use Screens

### Use Screens When

1. **Primary Content Display**
   - List of items (with navigation)
   - Detail view of single item
   - Dashboard/overview
   - Settings panel

2. **Complex Forms** (multiple fields)
   - Create/edit forms
   - Multi-step forms
   - Forms with validation
   - Forms with dependencies

3. **Data Interaction**
   - Browsing/searching data
   - Filtering/sorting
   - Pagination
   - Selection from list

4. **State Machine States**
   - Each state in workflow = Screen
   - Example: List → Detail → Edit → Confirm

### Screen Types

**Location**: `internal/cli/screens/{feature}/`

| Screen Type | Use Case | Behavior | Example |
|-------------|----------|----------|---------|
| **ListScreen** | Display list of items | TableBehavior | Browse events |
| **DetailScreen** | Show single item | UIKit layout | Event details |
| **FormScreen** | Capture/edit data | forms.Form | Create event |
| **FilterScreen** | Filter/search options | Form-like | Filter by date |
| **DashboardScreen** | Overview/summary | Mixed components | Stats dashboard |

### Screen Pattern

```go
// File: screens/myfeature/list_screen.go
package myfeature

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/base"
    "github.com/baphled/kariya/internal/cli/uikit/layout"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
    tea "github.com/charmbracelet/bubbletea"
)

// ListScreen displays a list of items.
type ListScreen struct {
    *base.BaseScreen
    table *behaviors.TableBehavior[*Item]
}

// NewListScreen creates a new list screen.
func NewListScreen(items []*Item) *ListScreen {
    columns := []behaviors.Column{
        {Title: "Name", Width: 30},
        {Title: "Date", Width: 20},
    }
    
    table := behaviors.NewTableBehavior(
        items,
        columns,
        func(item *Item) []string {
            return []string{item.Name, item.Date}
        },
    )
    
    return &ListScreen{
        BaseScreen: base.NewBaseScreen(),
        table:      table,
    }
}

// Update handles screen updates.
func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            if selected := s.table.SelectedItem(); selected != nil {
                return nil, screens.NewNavigateResult("detail", selected)
            }
        case "esc":
            return nil, screens.NewCancelResult("list")
        case "d":
            if selected := s.table.SelectedItem(); selected != nil {
                return nil, screens.NewNavigateResult("delete-confirm", selected)
            }
        }
    }
    
    cmd := s.table.Update(msg)
    return cmd, nil
}

// View renders the screen.
func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    header := primitives.Title("My Feature", theme)
    tableView := s.table.View()
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("d", "delete", theme) + " " +
            primitives.HelpKeyBadge("esc", "back", theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

### Screen Rules

**DO**:
- ✅ Embed `*base.BaseScreen`
- ✅ Implement `Update(msg) (tea.Cmd, screens.ScreenResult)`
- ✅ Implement `View() string`
- ✅ Use UIKit components (not raw lipgloss)
- ✅ Use Behaviors for complex UI (TableBehavior, etc.)
- ✅ Return ScreenResult to communicate with intent
- ✅ Keep screens stateless from intent perspective

**DON'T**:
- ❌ Import intents package (circular dependency)
- ❌ Call services directly (return ScreenResult instead)
- ❌ Mutate intent state (use ScreenResult)
- ❌ Put business logic in screens (delegate to context)
- ❌ Use raw lipgloss (use UIKit)

---

## When to Use Forms

### Use Forms When

1. **Data Capture Required**
   - Create new record
   - Edit existing record
   - Multi-field input
   - Structured data entry

2. **Validation Needed**
   - Required fields
   - Format validation (email, date, etc.)
   - Business rule validation
   - Cross-field validation

3. **Multi-Step Input**
   - Wizard-style forms
   - Progressive disclosure
   - Conditional fields
   - Grouped inputs

### Form Pattern

**CORRECT Pattern** (Use FormScreen):
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
    form     forms.Form
    formData *MyFormData
}

// NewFormScreen creates a new form screen.
func NewFormScreen() *FormScreen {
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
            forms.NewSelect(
                "category",
                "Category",
                []forms.Option{
                    {Key: "work", Label: "Work"},
                    {Key: "personal", Label: "Personal"},
                },
            ),
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
    // Check form state
    if forms.IsCompleted(s.form) {
        // Extract data
        s.formData.Name = s.form.GetString("name")
        s.formData.Description = s.form.GetString("description")
        s.formData.Category = s.form.GetString("category")
        
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

**DEPRECATED Pattern** (Don't Use):
```go
// ❌ WRONG - models/ package (DEPRECATED)
import "github.com/baphled/kariya/internal/cli/models"

type MyIntent struct {
    form *models.CaptureForm  // DEPRECATED
}
```

### Form Rules

**DO**:
- ✅ Create FormScreen in `screens/{feature}/form_screen.go`
- ✅ Use `forms.Form` directly (NOT `*huh.Form`)
- ✅ Build forms with `forms.NewInput()`, `forms.NewSelect()`, etc.
- ✅ Check state with `forms.IsCompleted()`, `forms.IsAborted()`
- ✅ Return ScreenResult (Submit, Cancel, Error)
- ✅ Validate in intent (not in screen)

**DON'T**:
- ❌ Use `models/` package for forms (DEPRECATED - Check #22)
- ❌ Import `huh` directly (use `forms/` package)
- ❌ Put business logic in FormScreen
- ❌ Call services from FormScreen
- ❌ Create custom form components (use forms builders)

---

## When to Use Behaviors

### Use Behaviors When

1. **Complex Reusable UI**
   - Tables with sorting/filtering
   - CRUD operations pattern
   - Pagination
   - Complex selection logic

2. **Shared Interaction Patterns**
   - Same behavior across multiple screens
   - Standard keyboard shortcuts
   - Common state management
   - Reusable business logic

### Behavior Types

**Location**: `internal/cli/behaviors/`

| Behavior | Use Case | Features |
|----------|----------|----------|
| **TableBehavior[T]** | Display lists | Sorting, selection, navigation, formatting |
| **CRUDBehavior[T]** | CRUD operations | Create, read, update, delete patterns |
| **PaginationBehavior** | Paginate lists | Page navigation, size control |
| **FilterBehavior** | Filter data | Multiple criteria, reset, apply |

### Behavior Pattern

```go
// In screen
type ListScreen struct {
    *base.BaseScreen
    table *behaviors.TableBehavior[*Item]
}

func NewListScreen(items []*Item) *ListScreen {
    columns := []behaviors.Column{
        {Title: "Name", Width: 30},
        {Title: "Date", Width: 20},
        {Title: "Status", Width: 15},
    }
    
    table := behaviors.NewTableBehavior(
        items,
        columns,
        func(item *Item) []string {
            return []string{
                item.Name,
                item.Date.Format("2006-01-02"),
                item.Status,
            }
        },
    )
    
    return &ListScreen{
        BaseScreen: base.NewBaseScreen(),
        table:      table,
    }
}

func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    // Handle screen-specific keys
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "enter":
            if selected := s.table.SelectedItem(); selected != nil {
                return nil, screens.NewNavigateResult("detail", selected)
            }
        }
    }
    
    // Delegate to behavior
    cmd := s.table.Update(msg)
    return cmd, nil
}

func (s *ListScreen) View() string {
    return s.table.View()
}
```

### Behavior Rules

**DO**:
- ✅ Use behaviors for complex reusable patterns
- ✅ Embed behaviors in screens (not intents)
- ✅ Delegate Update() to behavior
- ✅ Let behavior handle rendering
- ✅ Use generic types for type safety

**DON'T**:
- ❌ Put behaviors directly in intents (use screens)
- ❌ Duplicate behavior logic (extract to behavior)
- ❌ Create behavior for single-use case (inline instead)

---

## Component Selection Matrix

### By Use Case

| Use Case | Component | Location | Pattern |
|----------|-----------|----------|---------|
| **List of items** | ListScreen + TableBehavior | `screens/{feature}/list_screen.go` | Screen with embedded behavior |
| **Item details** | DetailScreen | `screens/{feature}/detail_screen.go` | Screen with UIKit layout |
| **Create/edit form** | FormScreen | `screens/{feature}/form_screen.go` | Screen with embedded forms.Form |
| **Delete confirmation** | ConfirmModal | Intent field | Centralized modal |
| **Error message** | ErrorModal | Intent field | Centralized modal |
| **Success message** | SuccessModal | Intent field | Centralized modal |
| **Loading indicator** | LoadingModal | Intent field | Centralized modal |
| **Filter/search** | FilterScreen | `screens/{feature}/filter_screen.go` | Screen with form-like inputs |
| **Wizard/multi-step** | Multiple FormScreens | `screens/{feature}/step*_screen.go` | State machine with screens |
| **Dashboard** | DashboardScreen | `screens/{feature}/dashboard_screen.go` | Screen with mixed components |

### By Interaction Time

| Duration | Component | Example |
|----------|-----------|---------|
| **2-3 seconds** | Modal (feedback) | "Success!", "Are you sure?" |
| **5-30 seconds** | FormScreen (simple) | Name + description form |
| **30+ seconds** | FormScreen (complex) | Multi-field wizard |
| **1-5 minutes** | ListScreen + DetailScreen | Browse and view items |
| **Ongoing** | Screen (primary) | Main application view |

### By Complexity

| Complexity | Component | When to Use |
|------------|-----------|-------------|
| **Trivial** | Modal | Single question/message |
| **Simple** | FormScreen (1-3 fields) | Quick data entry |
| **Medium** | ListScreen + DetailScreen | Browse/view pattern |
| **Complex** | Multiple Screens + Behaviors | Multi-state workflows |
| **Very Complex** | Multiple Screens + Behaviors + Modals | Full CRUD with confirmations |

---

## Common Scenarios

### Scenario 1: Browse and View Items

**Pattern**: ListScreen → DetailScreen

```go
// State machine
type MyFeatureState string

const (
    StateList   MyFeatureState = "list"    // ListScreen
    StateDetail MyFeatureState = "detail"  // DetailScreen
)

// Intent fields
type MyIntent struct {
    *BaseIntent
    state        MyFeatureState
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    activeScreen screens.Screen
}

// Init: Start with list
func (i *MyIntent) Init() tea.Cmd {
    i.listScreen = myfeature.NewListScreen(i.context.GetItems())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.activeScreen = i.listScreen
    i.state = StateList
    return nil
}

// HandleNavigate: Switch to detail
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    if result.Target == "detail" {
        item := result.Data.(*Item)
        i.detailScreen = myfeature.NewDetailScreen(item)
        i.detailScreen.SetTerminalInfo(i.GetTerminalInfo())
        i.detailScreen.SetTheme(i.Theme())
        i.activeScreen = i.detailScreen
        i.state = StateDetail
    }
    return nil
}
```

**No modals needed** (unless delete confirmation)

---

### Scenario 2: Create New Item

**Pattern**: FormScreen → SuccessModal

```go
// State machine
const (
    StateList MyFeatureState = "list"  // ListScreen
    StateForm MyFeatureState = "form"  // FormScreen
)

// Intent fields
type MyIntent struct {
    *BaseIntent
    state         MyFeatureState
    listScreen    *myfeature.ListScreen
    formScreen    *myfeature.FormScreen
    activeScreen  screens.Screen
    successModal  *feedback.SuccessModal
}

// HandleNavigate: Show form
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    if result.Target == "create" {
        i.formScreen = myfeature.NewFormScreen()
        i.formScreen.SetTerminalInfo(i.GetTerminalInfo())
        i.formScreen.SetTheme(i.Theme())
        i.activeScreen = i.formScreen
        i.state = StateForm
    }
    return nil
}

// HandleSubmit: Save and show success
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    formData := result.Data.(*myfeature.FormData)
    
    // Save (delegate to context)
    ctx := i.GetContext()
    err := i.context.SaveItem(ctx, formData)
    if err != nil {
        i.errorModal = feedback.NewErrorModal(err.Error())
        i.errorModal.Show()
        return nil
    }
    
    // Show success modal
    i.successModal = feedback.NewSuccessModal("Item created successfully")
    i.successModal.Show()
    
    // Return to list after modal dismissed
    i.activeScreen = i.listScreen
    i.state = StateList
    
    return nil
}
```

**Uses**: FormScreen + SuccessModal

---

### Scenario 3: Delete Item with Confirmation

**Pattern**: ListScreen → ConfirmModal → SuccessModal

```go
// Intent fields
type MyIntent struct {
    *BaseIntent
    listScreen    *myfeature.ListScreen
    activeScreen  screens.Screen
    deleteModal   *feedback.ConfirmModal
    successModal  *feedback.SuccessModal
    itemToDelete  *Item
}

// In ListScreen Update: Trigger delete
case "d":
    if selected := s.table.SelectedItem(); selected != nil {
        return nil, screens.NewNavigateResult("delete-confirm", selected)
    }

// In Intent HandleNavigate: Show confirm modal
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    if result.Target == "delete-confirm" {
        i.itemToDelete = result.Data.(*Item)
        i.deleteModal = feedback.NewConfirmModal(
            fmt.Sprintf("Delete '%s'?", i.itemToDelete.Name),
            feedback.WithDanger(),
        )
        i.deleteModal.Show()
    }
    return nil
}

// In Intent Update: Handle modal response
if i.deleteModal != nil && i.deleteModal.IsVisible() {
    model, cmd := i.deleteModal.Update(msg)
    i.deleteModal = model.(*feedback.ConfirmModal)
    
    if i.deleteModal.Confirmed() {
        return i.handleDeleteConfirm()
    }
    if i.deleteModal.Cancelled() {
        i.deleteModal.Hide()
    }
    
    return cmd
}

// handleDeleteConfirm: Delete and show success
func (i *MyIntent) handleDeleteConfirm() tea.Cmd {
    ctx := i.GetContext()
    err := i.context.DeleteItem(ctx, i.itemToDelete.ID)
    if err != nil {
        i.errorModal = feedback.NewErrorModal(err.Error())
        i.errorModal.Show()
        return nil
    }
    
    i.deleteModal.Hide()
    i.successModal = feedback.NewSuccessModal("Item deleted")
    i.successModal.Show()
    
    // Refresh list
    i.listScreen = myfeature.NewListScreen(i.context.GetItems())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.activeScreen = i.listScreen
    
    return nil
}
```

**Uses**: ListScreen + ConfirmModal + SuccessModal

---

### Scenario 4: Multi-Step Wizard

**Pattern**: Multiple FormScreens (step1, step2, step3)

```go
// State machine
const (
    StateStep1 MyFeatureState = "step1"
    StateStep2 MyFeatureState = "step2"
    StateStep3 MyFeatureState = "step3"
    StateReview MyFeatureState = "review"
)

// Intent fields
type MyIntent struct {
    *BaseIntent
    state        MyFeatureState
    step1Screen  *myfeature.Step1Screen
    step2Screen  *myfeature.Step2Screen
    step3Screen  *myfeature.Step3Screen
    reviewScreen *myfeature.ReviewScreen
    activeScreen screens.Screen
    wizardData   *WizardData
}

// HandleSubmit: Advance through wizard
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    switch i.state {
    case StateStep1:
        // Save step 1 data
        i.wizardData.Step1 = result.Data.(*Step1Data)
        
        // Advance to step 2
        i.step2Screen = myfeature.NewStep2Screen()
        i.step2Screen.SetTerminalInfo(i.GetTerminalInfo())
        i.step2Screen.SetTheme(i.Theme())
        i.activeScreen = i.step2Screen
        i.state = StateStep2
        
    case StateStep2:
        i.wizardData.Step2 = result.Data.(*Step2Data)
        i.activeScreen = i.step3Screen
        i.state = StateStep3
        
    case StateStep3:
        i.wizardData.Step3 = result.Data.(*Step3Data)
        i.activeScreen = i.reviewScreen
        i.state = StateReview
        
    case StateReview:
        // Final submission
        ctx := i.GetContext()
        err := i.context.SaveWizardData(ctx, i.wizardData)
        if err != nil {
            i.errorModal = feedback.NewErrorModal(err.Error())
            i.errorModal.Show()
            return nil
        }
        
        i.successModal = feedback.NewSuccessModal("Wizard completed!")
        i.successModal.Show()
        i.setCompleted(i.wizardData)
    }
    
    return nil
}
```

**Uses**: Multiple FormScreens + ReviewScreen + SuccessModal

---

### Scenario 5: Filter and Browse

**Pattern**: ListScreen ⇄ FilterScreen

```go
// State machine
const (
    StateList   MyFeatureState = "list"
    StateFilter MyFeatureState = "filter"
)

// Intent fields
type MyIntent struct {
    *BaseIntent
    state        MyFeatureState
    listScreen   *myfeature.ListScreen
    filterScreen *myfeature.FilterScreen
    activeScreen screens.Screen
    filters      *FilterConfig
}

// In ListScreen: Navigate to filter
case "f":
    return nil, screens.NewNavigateResult("filter", nil)

// HandleNavigate: Show filter screen
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    if result.Target == "filter" {
        i.filterScreen = myfeature.NewFilterScreen(i.filters)
        i.filterScreen.SetTerminalInfo(i.GetTerminalInfo())
        i.filterScreen.SetTheme(i.Theme())
        i.activeScreen = i.filterScreen
        i.state = StateFilter
    }
    return nil
}

// HandleSubmit: Apply filters
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    if i.state == StateFilter {
        i.filters = result.Data.(*FilterConfig)
        
        // Apply filters (delegate to context)
        filteredItems := i.context.FilterItems(i.filters)
        
        // Update list screen
        i.listScreen = myfeature.NewListScreen(filteredItems)
        i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
        i.listScreen.SetTheme(i.Theme())
        i.activeScreen = i.listScreen
        i.state = StateList
    }
    return nil
}
```

**Uses**: ListScreen + FilterScreen (no modals needed)

---

## Anti-Patterns

### Anti-Pattern 1: Using Modal for Form

```go
// ❌ WRONG - Complex form in modal
i.formModal = feedback.NewFormModal(
    "Create Item",
    []Field{
        {Name: "name", ...},
        {Name: "description", ...},
        {Name: "category", ...},
        {Name: "tags", ...},
    },
)

// ✅ CORRECT - Use FormScreen
i.formScreen = myfeature.NewFormScreen()
i.activeScreen = i.formScreen
i.state = StateForm
```

**Why**: Modals are for quick interactions, forms need full screen space

---

### Anti-Pattern 2: Using Modal for List

```go
// ❌ WRONG - List in modal
i.selectModal = feedback.NewSelectModal("Choose item", items)

// ✅ CORRECT - Use ListScreen
i.listScreen = myfeature.NewListScreen(items)
i.activeScreen = i.listScreen
i.state = StateList
```

**Why**: Lists need navigation, selection, and detail views

---

### Anti-Pattern 3: Custom Modal Instead of Centralized

```go
// ❌ WRONG - Custom modal type
type DeleteConfirmModal struct {
    visible bool
    message string
    // ... 100 lines of custom modal code
}

// ✅ CORRECT - Use centralized modal
i.deleteModal = feedback.NewConfirmModal(
    "Delete this item?",
    feedback.WithDanger(),
)
```

**Why**: Centralized modals are tested, consistent, and maintained

---

### Anti-Pattern 4: Screen for Quick Confirmation

```go
// ❌ WRONG - Dedicated screen for yes/no
type ConfirmDeleteScreen struct {
    *base.BaseScreen
    message string
    // ... rendering code
}

// ✅ CORRECT - Use ConfirmModal
i.deleteModal = feedback.NewConfirmModal("Delete?")
```

**Why**: Confirmations are temporary, don't need full screen

---

### Anti-Pattern 5: Multiple Modals Stacked

```go
// ❌ WRONG - Stacking modals
if i.confirmModal != nil && i.confirmModal.IsVisible() {
    if i.warningModal != nil && i.warningModal.IsVisible() {
        // Two modals visible at once
    }
}

// ✅ CORRECT - One modal at a time
if i.confirmModal != nil && i.confirmModal.IsVisible() {
    return behaviors.RenderModalOverlay(i.confirmModal, baseView)
}
if i.warningModal != nil && i.warningModal.IsVisible() {
    return behaviors.RenderModalOverlay(i.warningModal, baseView)
}
```

**Why**: Only one modal should be visible at a time

---

### Anti-Pattern 6: Business Logic in Modal

```go
// ❌ WRONG - Business logic in modal handling
if i.deleteModal.Confirmed() {
    // Direct DB call
    db.Query("DELETE FROM items WHERE id = ?", id)
}

// ✅ CORRECT - Delegate to context
if i.deleteModal.Confirmed() {
    ctx := i.GetContext()
    err := i.context.DeleteItem(ctx, id)
}
```

**Why**: Business logic belongs in context, not UI layer

---

### Anti-Pattern 7: Using models/ Package for Forms

```go
// ❌ WRONG - models/ package (DEPRECATED)
import "github.com/baphled/kariya/internal/cli/models"
type MyIntent struct {
    form *models.CaptureForm
}

// ✅ CORRECT - FormScreen
import "github.com/baphled/kariya/internal/cli/screens/myfeature"
type MyIntent struct {
    formScreen *myfeature.FormScreen
}
```

**Why**: Check #22 blocks models/ usage for forms

---

## Summary

### Quick Reference

**Use Modal when**:
- Confirmation (Y/N)
- Feedback message (success/error)
- Loading indicator
- Duration: 2-5 seconds
- Minimal interaction

**Use Screen when**:
- Primary content
- Complex forms (>3 fields)
- Lists/tables
- Detail views
- Duration: >10 seconds
- Full interaction

**Use FormScreen when**:
- Data capture
- Validation needed
- Multi-field input
- Create/edit operations

**Use Behavior when**:
- Complex reusable pattern
- Tables with features
- Shared interaction logic
- CRUD operations

### Component Hierarchy

```
Intent (orchestration)
├── Screens (primary content)
│   ├── ListScreen (with TableBehavior)
│   ├── DetailScreen (with UIKit layout)
│   └── FormScreen (with forms.Form)
└── Modals (temporary overlay)
    ├── ConfirmModal
    ├── ErrorModal
    ├── SuccessModal
    └── LoadingModal
```

### Decision Rules

1. **Duration < 5 seconds** → Modal
2. **Duration > 5 seconds** → Screen
3. **Simple Y/N** → ConfirmModal
4. **Form (1-2 fields)** → Modal or FormScreen
5. **Form (3+ fields)** → FormScreen (always)
6. **List navigation** → ListScreen (always)
7. **Feedback message** → Modal (always)
8. **Primary content** → Screen (always)

---

## Resources

### Documentation
- [Modal Patterns](../MODAL_PATTERNS.md)
- [Forms Guide](../FORMS_GUIDE.md)
- [UIKit Guide](../UIKIT_GUIDE.md)
- [Screen Extraction Guide](../guides/SCREEN_EXTRACTION_GUIDE.md)
- [Migration Rules](MIGRATION_RULES.md)

### Examples
- `examples/intent_subdirectory_template/`
- `internal/cli/uikit/feedback/` (modal examples)
- `internal/cli/screens/` (screen examples)

### Enforcement
- Check #19: No rendering in intent
- Check #21: UIKit component usage
- Check #22: No models/ for forms

---

**Last Updated**: 2026-01-22  
**Status**: MANDATORY - Must follow for all component selection  
**Enforcement**: Automated checks + code review
