# Screen Extraction Guide

**Purpose**: How to extract rendering logic from intents to screens  
**Status**: REQUIRED for Check #19 compliance (no rendering in intent)  
**Enforcement**: `check-intent-architecture.sh` (Check #19)

---

## Overview

### Why Extract Screens?

**Problem**: Intents with inline rendering violate single responsibility principle

```go
// ❌ BAD: Rendering in intent
func (i *MyIntent) View() string {
    // 200+ lines of lipgloss styling
    title := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        Background(lipgloss.Color("#5c5cff")).
        Bold(true).
        Padding(1, 2).
        Render("My Feature")
    
    body := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#dddddd")).
        Padding(2).
        Render(i.getBodyContent())
    
    // ... 150 more lines of styling
    
    return lipgloss.JoinVertical(lipgloss.Left, title, body, footer)
}

func (i *MyIntent) renderHeader() string { /* 50 lines */ }
func (i *MyIntent) renderBody() string { /* 80 lines */ }
func (i *MyIntent) renderFooter() string { /* 40 lines */ }
func (i *MyIntent) renderModal() string { /* 60 lines */ }
// 8+ render methods = VIOLATION
```

**Solution**: Extract to dedicated screen components

```go
// ✅ GOOD: Delegate to screen
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    
    return baseView
}
// Only 1 method, delegates rendering
```

**Benefits**:
- Intent stays lean (broker only)
- Screens are reusable
- Easier to test rendering
- Clear separation of concerns
- UIKit component usage

---

## Extraction Process

### Step 1: Identify Render Methods

**Find all rendering** in intent file:
```bash
# Search for render methods
grep -n "func.*render\|func.*Render\|func.*View" my_intent.go

# Count render methods (should be ≤2)
grep -c "^func.*render\|^func.*Render\|^func.*View" my_intent.go
```

**Example output**:
```
456:func (i *MyIntent) View() string {
512:func (i *MyIntent) renderHeader() string {
568:func (i *MyIntent) renderBody() string {
624:func (i *MyIntent) renderFooter() string {
680:func (i *MyIntent) renderList() string {
736:func (i *MyIntent) renderDetail() string {
792:func (i *MyIntent) renderForm() string {
848:func (i *MyIntent) renderModal() string {
```

**Analysis**: 8 render methods = **VIOLATION** (max 2 allowed)

---

### Step 2: Group by Screen

**Identify screens** based on state machine:

```go
// From intent state enum
type MyIntentState string

const (
    StateList   MyIntentState = "list"    // → ListScreen
    StateDetail MyIntentState = "detail"  // → DetailScreen
    StateForm   MyIntentState = "form"    // → FormScreen
)
```

**Mapping**:
| State | Screen | Render Methods |
|-------|--------|----------------|
| `StateList` | `ListScreen` | `renderList()` |
| `StateDetail` | `DetailScreen` | `renderDetail()` |
| `StateForm` | `FormScreen` | `renderForm()` |
| All | - | `renderHeader()`, `renderFooter()` (shared, use UIKit) |

---

### Step 3: Create Screen Files

**File structure**:
```
screens/{feature}/
├── list_screen.go
├── detail_screen.go
└── form_screen.go
```

**Template** for each screen:
```go
package {feature}

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
    // Setup table behavior
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
        }
    }
    
    // Delegate to table behavior
    cmd := s.table.Update(msg)
    return cmd, nil
}

// View renders the screen.
func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    header := primitives.Title("My Feature", theme)
    tableView := s.table.View()
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("esc", "back", theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

---

### Step 4: Move Rendering Logic

**Copy rendering code** from intent to screen:

#### BEFORE (Intent)
```go
// File: my_intent.go
func (i *MyIntent) renderList() string {
    // 100+ lines of lipgloss styling
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        Background(lipgloss.Color("#5c5cff")).
        Bold(true).
        Padding(1, 2)
    
    title := titleStyle.Render("My Feature")
    
    bodyStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#dddddd")).
        Padding(2)
    
    // Manual table rendering
    rows := []string{}
    for i, item := range i.items {
        selected := ""
        if i == i.selectedIndex {
            selected = "> "
        }
        rows = append(rows, fmt.Sprintf("%s%s - %s", selected, item.Name, item.Date))
    }
    
    body := bodyStyle.Render(strings.Join(rows, "\n"))
    footer := "Press enter to select, esc to cancel"
    
    return lipgloss.JoinVertical(lipgloss.Left, title, body, footer)
}
```

#### AFTER (Screen)
```go
// File: screens/myfeature/list_screen.go
func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    // Use UIKit primitives (not raw lipgloss)
    header := primitives.Title("My Feature", theme)
    
    // Use TableBehavior (not manual rendering)
    tableView := s.table.View()
    
    // Use UIKit help badges (not raw strings)
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("esc", "cancel", theme)
    
    // Use layout system (not manual JoinVertical)
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

**Key Changes**:
1. Raw lipgloss → UIKit primitives
2. Manual table rendering → `TableBehavior`
3. Raw strings → `HelpKeyBadge()`
4. Manual layout → `ScreenLayout`
5. Result: **15 lines vs 100 lines** (85% reduction!)

---

### Step 5: Update Intent

**Replace rendering** with delegation:

#### BEFORE (Intent with inline rendering)
```go
func (i *MyIntent) View() string {
    switch i.state {
    case StateList:
        return i.renderList()    // 100 lines
    case StateDetail:
        return i.renderDetail()  // 80 lines
    case StateForm:
        return i.renderForm()    // 120 lines
    }
    return ""
}
```

#### AFTER (Intent delegates to screens)
```go
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    // Overlay modals (centralized)
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**Result**: Intent rendering reduced from **300+ lines to ~10 lines** (97% reduction!)

---

### Step 6: Update Intent Fields

**Add typed screen fields**:

#### BEFORE
```go
type MyIntent struct {
    *BaseIntent
    // No screen fields
}
```

#### AFTER
```go
type MyIntent struct {
    *BaseIntent
    
    // Explicit typed screen fields (REQUIRED)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    formScreen   *myfeature.FormScreen
    
    // Active screen (generic pointer)
    activeScreen screens.Screen
}
```

---

### Step 7: Initialize Screens

**In `Init()` method**:

```go
func (i *MyIntent) Init() tea.Cmd {
    // Create list screen
    i.listScreen = myfeature.NewListScreen(i.context.GetItems())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    
    // Set as active
    i.activeScreen = i.listScreen
    i.state = StateList
    
    return nil
}
```

---

### Step 8: Handle Screen Results

**Update method** to handle screen results:

```go
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Delegate to active screen
    cmd, result := i.activeScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

func (i *MyIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    switch r := result.(type) {
    case *screens.CancelResult:
        return i.HandleCancel(r)
    case *screens.NavigateResult:
        return i.HandleNavigate(r)
    case *screens.SubmitResult:
        return i.HandleSubmit(r)
    case *screens.ErrorResult:
        return i.HandleError(r)
    }
    return nil
}
```

---

## UIKit Components Reference

**Use these instead of raw lipgloss**:

| Need | UIKit Component | Not |
|------|-----------------|-----|
| Title | `primitives.Title("text", theme)` | `lipgloss.NewStyle().Bold(true).Render()` |
| Body text | `primitives.Body("text", theme)` | `lipgloss.NewStyle().Render()` |
| Error text | `primitives.ErrorText("text", theme)` | `lipgloss.NewStyle().Foreground(red)` |
| Success text | `primitives.SuccessText("text", theme)` | `lipgloss.NewStyle().Foreground(green)` |
| Help badge | `primitives.HelpKeyBadge("key", "label", theme)` | `components.KeyBadge()` |
| Layout | `layout.NewScreenLayout(theme).WithContent()` | `lipgloss.JoinVertical()` |
| Table | `behaviors.TableBehavior[T]` | Manual table rendering |
| Modal | `feedback.ConfirmModal`, `feedback.ErrorModal` | Custom modal code |

**See**: [UIKit Guide](../UIKIT_GUIDE.md)

---

## Extraction Checklist

For each screen extraction:

### Discovery
- [ ] Identify all render methods in intent
- [ ] Map render methods to states
- [ ] Determine screens needed (list, detail, form, etc.)

### Implementation
- [ ] Create `screens/{feature}/` directory
- [ ] Create screen file for each state
- [ ] Implement `BaseScreen` embedding
- [ ] Implement `Update()` method (handle input, return `ScreenResult`)
- [ ] Implement `View()` method (use UIKit, not raw lipgloss)
- [ ] Use behaviors (`TableBehavior`, `CRUDBehavior`, etc.)

### Integration
- [ ] Add typed screen fields to intent
- [ ] Initialize screens in `Init()`
- [ ] Update intent `Update()` to delegate to screens
- [ ] Update intent `View()` to delegate to screens
- [ ] Implement `ScreenResultHandler` interface

### Validation
- [ ] Run `make check-intent-architecture` (Check #19)
- [ ] Verify intent has ≤2 render methods
- [ ] Verify screens use UIKit components (Check #21)
- [ ] Run unit tests
- [ ] Manual testing

---

## Common Pitfalls

### 1. Forgetting ScreenResult

```go
// ❌ WRONG - No return value
func (s *ListScreen) Update(msg tea.Msg) tea.Cmd {
    // ...
    return nil
}

// ✅ CORRECT - Return ScreenResult
func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg.String() {
    case "esc":
        return nil, screens.NewCancelResult("list")
    }
    return nil, nil
}
```

### 2. Raw Lipgloss Instead of UIKit

```go
// ❌ WRONG - Raw lipgloss
title := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#ffffff")).
    Bold(true).
    Render("Title")

// ✅ CORRECT - UIKit primitive
title := primitives.Title("Title", theme)
```

### 3. Manual Table Rendering

```go
// ❌ WRONG - Manual rendering
for i, item := range items {
    selected := ""
    if i == selectedIndex {
        selected = "> "
    }
    fmt.Printf("%s%s\n", selected, item.Name)
}

// ✅ CORRECT - TableBehavior
table := behaviors.NewTableBehavior(items, columns, renderFunc)
tableView := table.View()
```

### 4. Screens Importing Intents

```go
// ❌ FORBIDDEN - Circular dependency
package myfeature  // screens/myfeature/

import "github.com/baphled/kariya/internal/cli/intents"

// ✅ CORRECT - Screens NEVER import intents
package myfeature

import "github.com/baphled/kariya/internal/cli/screens"
```

---

## Before/After Comparison

### BEFORE (All rendering in intent)

```
File: my_intent.go (1,200 lines)

- Intent struct (50 lines)
- Init/Update methods (100 lines)
- renderHeader() (50 lines)
- renderBody() (80 lines)
- renderFooter() (40 lines)
- renderList() (200 lines)
- renderDetail() (180 lines)
- renderForm() (250 lines)
- renderModal() (100 lines)
- Helper methods (150 lines)

TOTAL: 1,200 lines in ONE file
VIOLATIONS: 8 render methods (>2 limit)
```

### AFTER (Screens extracted)

```
File: intents/myfeature/intent.go (290 lines)
- Intent struct (50 lines)
- Init/Update methods (100 lines)
- View() delegation (10 lines)
- ScreenResultHandler (80 lines)
- Helpers (50 lines)

File: screens/myfeature/list_screen.go (90 lines)
File: screens/myfeature/detail_screen.go (85 lines)
File: screens/myfeature/form_screen.go (110 lines)

TOTAL: 575 lines across 4 files (48% reduction!)
VIOLATIONS: 0 (intent has 1 render method)
```

---

## Resources

### Templates
- `examples/intent_subdirectory_template/intent.go.template` - Intent template
- `examples/intent_subdirectory_template/screens/list_screen.go.template` - Screen template

### Documentation
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md) - Architecture patterns
- [UIKit Guide](../UIKIT_GUIDE.md) - UIKit component reference
- [Intent Migration Guide](INTENT_MIGRATION_TO_SUBDIRECTORY.md) - Full migration process
- [AGENTS.md](../../AGENTS.md) - AI agent rules

### Tools
- `make check-intent-architecture` - Validate extraction (Check #19)
- `make what-to-use NEED="keyword"` - Component lookup

---

**Last Updated**: 2026-01-22  
**Status**: Required for Check #19 compliance  
**Enforcement**: Automated by `check-intent-architecture.sh` (Check #19)
