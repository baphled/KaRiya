---
name: bubble-tea-expert
description: Expert in Charm's Bubble Tea TUI framework and KaRiya implementation patterns
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Bubble Tea Expert Skill

You are an expert in Charm's Bubble Tea TUI framework and KaRiya's specific implementation patterns.

## Core Concepts

### The Elm Architecture (TEA)

Bubble Tea implements The Elm Architecture:
- **Model**: Application state
- **Update**: State transitions based on messages
- **View**: Render state to string

```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

### KaRiya's Modified Pattern

KaRiya modifies the pattern for its layer architecture:

| Layer | Init | Update | View |
|-------|------|--------|------|
| App | `tea.Cmd` | `(tea.Model, tea.Cmd)` | `string` |
| Intent | `tea.Cmd` | `tea.Cmd` (no Model) | `string` |
| Screen | N/A | `(tea.Cmd, ScreenResult)` | `string` |
| Modal | `tea.Cmd` | `(tea.Model, tea.Cmd)` | `string` |

**Critical**: Intents return only `tea.Cmd` from Update, NOT `(tea.Model, tea.Cmd)`.

---

## Message Flow

```
User Input (tea.KeyMsg)
         │
         ▼
    App.Update()
         │
         ▼
  Intent.Update()
         │
    ┌────┴────┐
    │         │
    ▼         ▼
 Modal    Screen.Update()
.Update()      │
    │          ▼
    │     ScreenResult
    │          │
    └────┬─────┘
         │
         ▼
  Intent handles result
         │
         ▼
    State transition
```

---

## Intent Pattern

### Structure

```go
type Intent struct {
    *intents.BaseIntent  // REQUIRED: Provides terminal, theme, logo
    
    context *IntentContext   // Input parameters (context.go)
    state   State            // Typed enum (constants.go)
    active  bool             // Is intent active
    result  *intents.IntentResult[*Result]
    
    // Explicit typed screen fields
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Generic pointer for delegation
    activeScreen screens.Screen
    
    // Modal fields
    detailModal  *modals.DetailModal
    deleteModal  *feedback.ConfirmModal
}
```

### Init Pattern

```go
func (i *Intent) Init() tea.Cmd {
    // 1. Validate context
    if err := i.context.Validate(); err != nil {
        i.SetError(err)
        return nil
    }
    
    // 2. Load initial data
    items, err := i.context.Repository.FindAll()
    if err != nil {
        i.SetError(err)
        return nil
    }
    
    // 3. Initialize state
    i.items = items
    i.state = StateList
    i.active = true
    
    // 4. Create and transition to first screen
    i.listScreen = NewListScreen(items)
    i.transitionToScreen(i.listScreen)
    
    return nil
}
```

### Update Pattern

```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // 1. Guard: Check if active
    if !i.active {
        return nil
    }
    
    // 2. Handle custom messages (from async operations)
    switch msg := msg.(type) {
    case DataLoadedMsg:
        return i.handleDataLoaded(msg)
    case OperationCompleteMsg:
        return i.handleOperationComplete(msg)
    }
    
    // 3. Handle modals (highest UI priority)
    if cmd := i.handleModalUpdates(msg); cmd != nil || i.HasActiveModal() {
        return cmd
    }
    
    // 4. Handle keyboard shortcuts
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if cmd := i.handleKeyShortcuts(keyMsg); cmd != nil {
            return cmd
        }
    }
    
    // 5. Delegate to active screen
    if i.activeScreen != nil {
        cmd, result := i.activeScreen.Update(msg)
        if result != nil {
            return i.handleScreenResult(result)
        }
        return cmd
    }
    
    return nil
}
```

### View Pattern

```go
func (i *Intent) View() string {
    // 1. Guard
    if !i.active {
        return "Intent not active"
    }
    
    // 2. Render screen
    var baseView string
    if i.activeScreen != nil {
        baseView = i.activeScreen.View()
    } else {
        view := i.CreateViewWithBreadcrumbs("Main", "Feature", i.getStateName())
        view.WithContent("Loading...")
        view.WithHelp(i.getContextHelp())
        baseView = view.Render()
    }
    
    // 3. Apply modal overlay (REQUIRED pattern)
    if i.detailModal != nil && i.detailModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.detailModal, baseView)
    }
    
    return baseView
}
```

---

## Screen Pattern

### Structure

```go
type ListScreen struct {
    *base.Screen  // REQUIRED: Provides terminal, theme, logo
    
    items         []*domain.Item
    tableBehavior *behaviors.TableBehavior[*domain.Item]
}
```

### Constructor

```go
func NewListScreen(items []*domain.Item, theme themes.Theme) *ListScreen {
    columns := []behaviors.ColumnDef{
        {Title: "Name", Width: 30},
        {Title: "Status", Width: 15},
        {Title: "Date", Width: 12},
    }
    
    formatter := func(item *domain.Item, _ int) []string {
        return []string{
            item.Name,
            item.Status,
            item.CreatedAt.Format("2006-01-02"),
        }
    }
    
    table := behaviors.NewTableBehavior[*domain.Item](theme, columns, formatter).
        PageSize(15).
        PaginationPrefix("Items").
        EmptyMessage("No items found.")
    
    table.SetItems(items)
    
    return &ListScreen{
        Screen:        base.NewScreen(),
        items:         items,
        tableBehavior: table,
    }
}
```

### Update (Returns ScreenResult)

```go
func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        s.SetTerminalInfo(msg.Width, msg.Height)
        return nil, nil
        
    case tea.KeyMsg:
        // Special keys: ALWAYS use tea.Key* constants
        switch msg.Type {
        case tea.KeyEsc:
            return nil, &screens.CancelResult{}
        case tea.KeyEnter:
            if item := s.tableBehavior.GetSelectedItem(); item != nil {
                return nil, &screens.NavigateResult{ResultData: *item}
            }
        case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown, tea.KeyHome, tea.KeyEnd:
            s.tableBehavior.HandleNavigation(s.keyToNav(msg.Type))
            return nil, nil
        }
        
        // Rune keys: Use string comparison
        switch msg.String() {
        case "j":
            s.tableBehavior.HandleNavigation("down")
        case "k":
            s.tableBehavior.HandleNavigation("up")
        case "g":
            s.tableBehavior.HandleNavigation("home")
        case "G":
            s.tableBehavior.HandleNavigation("end")
        case "a":
            return nil, &screens.NavigateResult{
                ResultData: map[string]interface{}{"action": "add"},
            }
        case "e":
            if item := s.tableBehavior.GetSelectedItem(); item != nil {
                return nil, &screens.NavigateResult{
                    ResultData: map[string]interface{}{
                        "action": "edit",
                        "item":   *item,
                    },
                }
            }
        case "d":
            if item := s.tableBehavior.GetSelectedItem(); item != nil {
                return nil, &screens.NavigateResult{
                    ResultData: map[string]interface{}{
                        "action": "delete",
                        "item":   *item,
                    },
                }
            }
        }
    }
    
    return nil, nil
}

func (s *ListScreen) keyToNav(keyType tea.KeyType) string {
    switch keyType {
    case tea.KeyUp:
        return "up"
    case tea.KeyDown:
        return "down"
    case tea.KeyPgUp:
        return "pgup"
    case tea.KeyPgDown:
        return "pgdn"
    case tea.KeyHome:
        return "home"
    case tea.KeyEnd:
        return "end"
    default:
        return ""
    }
}
```

---

## ScreenResult Pattern

### Result Types

```go
// Navigate to new state/screen
&screens.NavigateResult{
    ResultData: selectedItem,            // Single item
    // OR
    ResultData: map[string]interface{}{  // Action with context
        "action": "edit",
        "item":   item,
    },
}

// Cancel current operation
&screens.CancelResult{}

// Submit form data
&screens.SubmitResult{
    FormData: &MyFormData{
        Name: name,
        Description: desc,
    },
}

// Report error
&screens.ErrorResult{
    Err:     err,
    Message: "Failed to save",
}
```

### Handler Interface

```go
// Intent MUST implement this interface
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    // Handle action maps
    if actionData, ok := result.ResultData.(map[string]interface{}); ok {
        return i.handleAction(actionData)
    }
    
    // Handle direct item navigation
    if item, ok := result.ResultData.(*domain.Item); ok {
        i.selectedItem = item
        i.state = StateDetail
        i.transitionToScreen(NewDetailScreen(item))
    }
    return nil
}

func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    switch i.state {
    case StateList:
        i.SetCancelled()
    case StateDetail, StateEdit:
        i.state = StateList
        i.transitionToScreen(i.listScreen)
    }
    return nil
}

func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    formData := result.FormData.(*MyFormData)
    return i.saveItem(formData)
}

func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
    i.showErrorModal(result.Message)
    return nil
}
```

---

## Keyboard Handling

### CRITICAL: Use tea.Key* for Special Keys

```go
// CORRECT: Use constants for special keys
switch msg.Type {
case tea.KeyEsc:      // Escape
case tea.KeyEnter:    // Enter
case tea.KeyUp:       // Arrow up
case tea.KeyDown:     // Arrow down
case tea.KeyLeft:     // Arrow left
case tea.KeyRight:    // Arrow right
case tea.KeyPgUp:     // Page up
case tea.KeyPgDown:   // Page down
case tea.KeyHome:     // Home
case tea.KeyEnd:      // End
case tea.KeyTab:      // Tab
case tea.KeyBackspace:// Backspace
case tea.KeyDelete:   // Delete
case tea.KeySpace:    // Space
case tea.KeyCtrlC:    // Ctrl+C
case tea.KeyCtrlD:    // Ctrl+D
case tea.KeyCtrlU:    // Ctrl+U
}

// CORRECT: Use string for rune/vim keys
switch msg.String() {
case "j", "k", "h", "l":  // vim navigation
case "g", "G":             // vim go to
case "a", "e", "d":        // actions
case "q":                  // quit
case "/":                  // search
case "?":                  // help
}

// WRONG: Never use string for special keys
switch msg.String() {
case "esc":     // WRONG - inconsistent across terminals
case "up":      // WRONG - use tea.KeyUp
case "pgdown":  // WRONG - use tea.KeyPgDown
}
```

### Navigation Mapping

```go
// TableBehavior navigation strings
"up"     // Move up one row
"down"   // Move down one row
"pgup"   // Page up
"pgdn"   // Page down (note: NOT "pgdown")
"home"   // First row
"end"    // Last row
"ctrl+u" // Half page up
"ctrl+d" // Half page down
```

---

## Async Commands

### Pattern

```go
// Define message type (in messages.go)
type DataLoadedMsg struct {
    Items []*domain.Item
    Error error
}

// Create command (in helpers.go)
func (i *Intent) loadDataAsync() tea.Cmd {
    return func() tea.Msg {
        items, err := i.context.Repository.FindAll()
        return DataLoadedMsg{Items: items, Error: err}
    }
}

// Handle in Update (in intent.go or handlers.go)
func (i *Intent) handleDataLoaded(msg DataLoadedMsg) tea.Cmd {
    if msg.Error != nil {
        i.showErrorModal(msg.Error.Error())
        return nil
    }
    
    i.items = msg.Items
    i.listScreen.SetItems(msg.Items)
    return nil
}
```

### With Cancellation

```go
func (i *Intent) startLongOperation() tea.Cmd {
    ctx, cancel := context.WithCancel(context.Background())
    i.cancelFunc = cancel
    i.loading = true
    
    return tea.Batch(
        i.loadingModal.Init(),
        func() tea.Msg {
            result, err := i.context.Service.Process(ctx)
            return OperationCompleteMsg{Result: result, Error: err}
        },
    )
}

func (i *Intent) cancelOperation() {
    if i.cancelFunc != nil {
        i.cancelFunc()
        i.cancelFunc = nil
    }
    i.loading = false
}
```

---

## Modal Overlay

### REQUIRED: Use RenderModalOverlay

```go
// In intent View()
func (i *Intent) View() string {
    baseView := i.activeScreen.View()
    
    // CORRECT: Use RenderModalOverlay
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    
    return baseView
}

// WRONG: Direct modal rendering (loses background)
func (i *Intent) View() string {
    if i.modal != nil && i.modal.IsVisible() {
        return i.modal.View()  // WRONG!
    }
    return i.activeScreen.View()
}
```

### Modal Registry Pattern

```go
// For multiple modals
func (i *Intent) rebuildModalRegistry() {
    i.modalRegistry = behaviors.NewModalRegistry()
    
    if i.detailModal != nil {
        i.modalRegistry.Register("detail", i.detailModal)
    }
    if i.deleteModal != nil {
        i.modalRegistry.Register("delete", i.deleteModal)
    }
    if i.errorModal != nil {
        i.modalRegistry.Register("error", i.errorModal)
    }
}

func (i *Intent) View() string {
    baseView := i.activeScreen.View()
    i.rebuildModalRegistry()
    return i.modalRegistry.RenderOverlay(baseView)
}
```

---

## State Transitions

### Pattern

```go
func (i *Intent) transitionToScreen(screen screens.Screen) {
    // 1. Pass theme
    if i.Theme() != nil {
        screen.SetTheme(i.Theme())
    }
    
    // 2. Pass terminal dimensions
    if info := i.GetTerminalInfo(); info != nil {
        screen.SetTerminalInfo(info.Width, info.Height)
    }
    
    // 3. Pass logo
    if i.GetLogo() != nil {
        screen.SetLogo(i.GetLogo(), i.GetLogoSpacing())
    }
    
    // 4. Update active screen
    i.activeScreen = screen
}
```

### State Machine

```go
// constants.go
type State string

const (
    StateList          State = "list"
    StateDetail        State = "detail"
    StateEdit          State = "edit"
    StateDeleteConfirm State = "delete_confirm"
)

// Transition in handler
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    switch i.state {
    case StateList:
        if item, ok := result.ResultData.(*domain.Item); ok {
            i.selectedItem = item
            i.state = StateDetail
            i.detailScreen = NewDetailScreen(item)
            i.transitionToScreen(i.detailScreen)
        }
    case StateDetail:
        if actionData, ok := result.ResultData.(map[string]interface{}); ok {
            switch actionData["action"] {
            case "edit":
                i.state = StateEdit
                i.editScreen = NewEditScreen(i.selectedItem)
                i.transitionToScreen(i.editScreen)
            }
        }
    }
    return nil
}
```

---

## Forms (ONLY in forms/ package)

### Creating Forms

```go
// ONLY in internal/cli/forms/ package
import "github.com/charmbracelet/huh"

form := forms.NewForm(
    forms.NewGroup(
        forms.NewInput(forms.FieldConfig{
            Key:         "name",
            Title:       "Name",
            Placeholder: "Enter name...",
            Required:    true,
            Validate:    validateName,
        }),
        forms.NewText(forms.FieldConfig{
            Key:         "description",
            Title:       "Description",
            CharLimit:   500,
        }),
        forms.NewSelect("status", "Status",
            []string{"Active", "Inactive", "Pending"},
            "Active",
        ),
    ),
)
```

### Checking State

```go
// In screen Update
func (s *FormScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    cmd := s.form.Update(msg)
    
    if forms.IsCompleted(s.form) {
        data := &MyFormData{
            Name:        forms.GetString(s.form, "name"),
            Description: forms.GetString(s.form, "description"),
            Status:      forms.GetString(s.form, "status"),
        }
        return cmd, &screens.SubmitResult{FormData: data}
    }
    
    if forms.IsAborted(s.form) {
        return cmd, &screens.CancelResult{}
    }
    
    return cmd, nil
}
```

---

## Anti-Patterns (NEVER DO)

### 1. Import huh outside forms/

```go
// WRONG
import "github.com/charmbracelet/huh"  // Only allowed in forms/

// CORRECT
import "github.com/baphled/kariya/internal/cli/forms"
```

### 2. Screen importing intents

```go
// WRONG
import "github.com/baphled/kariya/internal/cli/intents"  // FORBIDDEN in screens/

// CORRECT: Use ScreenResult to communicate
return nil, &screens.NavigateResult{ResultData: item}
```

### 3. Intent returning (tea.Model, tea.Cmd)

```go
// WRONG
func (i *Intent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return i, nil
}

// CORRECT
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    return nil
}
```

### 4. Direct modal rendering

```go
// WRONG
return i.modal.View()

// CORRECT
return behaviors.RenderModalOverlay(i.modal, baseView)
```

### 5. Untyped state

```go
// WRONG
state string

// CORRECT
state State  // Typed enum
```

### 6. Missing BaseIntent/BaseScreen

```go
// WRONG
type Intent struct {
    // No base
}

// CORRECT
type Intent struct {
    *intents.BaseIntent
}
```

### 7. String comparison for special keys

```go
// WRONG
if msg.String() == "esc" { }

// CORRECT
if msg.Type == tea.KeyEsc { }
```

---

## Quick Reference

### Key Packages

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/base"
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/cli/themes"
    "github.com/baphled/kariya/internal/cli/uikit/feedback"
    "github.com/baphled/kariya/internal/cli/uikit/layout"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
)
```

### Update Signature by Layer

| Layer | Signature |
|-------|-----------|
| App | `Update(msg tea.Msg) (tea.Model, tea.Cmd)` |
| Intent | `Update(msg tea.Msg) tea.Cmd` |
| Screen | `Update(msg tea.Msg) (tea.Cmd, ScreenResult)` |
| Modal | `Update(msg tea.Msg) (tea.Model, tea.Cmd)` |

### ScreenResult Types

| Type | Purpose |
|------|---------|
| `NavigateResult` | Navigate to new state/screen |
| `CancelResult` | Cancel current operation |
| `SubmitResult` | Submit form data |
| `ErrorResult` | Report error |
