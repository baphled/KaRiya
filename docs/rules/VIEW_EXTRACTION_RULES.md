# View Extraction Rules

**Purpose**: Strict rules for extracting views/rendering from intents to screens  
**Status**: MANDATORY - Must follow for all view extractions  
**Enforcement**: Check #19 (blocks >2 render methods in intent)

---

## Table of Contents

1. [Core Principle](#core-principle)
2. [The Rule of One](#the-rule-of-one)
3. [Extraction Triggers](#extraction-triggers)
4. [What to Extract](#what-to-extract)
5. [Where to Extract](#where-to-extract)
6. [Extraction Process](#extraction-process)
7. [View Method Limits](#view-method-limits)
8. [UIKit Requirements](#uikit-requirements)
9. [Screen Responsibilities](#screen-responsibilities)
10. [Common Violations](#common-violations)

---

## Core Principle

**INTENTS ORCHESTRATE, SCREENS RENDER**

```
Intent:  Coordinates state machine, delegates to screens
Screen:  Handles rendering, user input, returns results
```

**One sentence rule**:
> An intent's View() method should ONLY delegate to screens and overlay modals.

---

## The Rule of One

**Intent View() Method Rules**:

1. **One responsibility**: Delegate to active screen
2. **One method**: Only `View()` (no `render*()` helpers)
3. **One pattern**: Screen delegation + modal overlay
4. **One exception**: May have 1 tiny helper (≤5 lines) for modal overlay logic

**Maximum**: 2 render methods total (View() + optional 1 helper ≤5 lines)

**Enforcement**: Check #19 blocks commits with >2 render methods

---

## Extraction Triggers

### MANDATORY Extraction (Must Extract)

Extract views if **ANY** of these are true:

1. **Multiple render methods** (>2 total)
   ```go
   // VIOLATION - 6 render methods
   func (i *MyIntent) View() string { ... }
   func (i *MyIntent) renderHeader() string { ... }
   func (i *MyIntent) renderBody() string { ... }
   func (i *MyIntent) renderFooter() string { ... }
   func (i *MyIntent) renderList() string { ... }
   func (i *MyIntent) renderDetail() string { ... }
   ```

2. **View() method >20 lines**
   ```go
   // VIOLATION - 80 lines of rendering
   func (i *MyIntent) View() string {
       // 80 lines of lipgloss styling
   }
   ```

3. **Any lipgloss styling in intent**
   ```go
   // VIOLATION - Direct styling in intent
   title := lipgloss.NewStyle().Bold(true).Render("Title")
   ```

4. **Manual layout in intent**
   ```go
   // VIOLATION - Manual composition in intent
   return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
   ```

5. **Switch statement for rendering**
   ```go
   // VIOLATION - Rendering logic in intent
   func (i *MyIntent) View() string {
       switch i.state {
       case StateList:
           return i.renderList()      // Extract to ListScreen
       case StateDetail:
           return i.renderDetail()    // Extract to DetailScreen
       case StateForm:
           return i.renderForm()      // Extract to FormScreen
       }
   }
   ```

**If any trigger is true**: Extract views to screens immediately.

---

## What to Extract

### ALWAYS Extract These

1. **All render*() methods**
   - `renderHeader()` → Screen method or UIKit primitive
   - `renderBody()` → Screen.View()
   - `renderFooter()` → Screen method or UIKit primitive
   - `renderList()` → ListScreen.View()
   - `renderDetail()` → DetailScreen.View()
   - `renderForm()` → FormScreen.View()

2. **All lipgloss styling**
   ```go
   // Extract this
   titleStyle := lipgloss.NewStyle().Bold(true).Foreground(...)
   bodyStyle := lipgloss.NewStyle().Padding(2).Margin(1)
   
   // To
   primitives.Title("text", theme)
   layout.NewScreenLayout(theme).WithContent(...)
   ```

3. **All layout composition**
   ```go
   // Extract this
   lipgloss.JoinVertical(lipgloss.Left, parts...)
   lipgloss.Place(width, height, ...)
   
   // To
   layout.NewScreenLayout(theme).WithHeader(...).WithContent(...)
   ```

4. **All formatting logic**
   ```go
   // Extract this
   func (i *MyIntent) formatItem(item *Item) string {
       return fmt.Sprintf("%s - %s", item.Name, item.Date)
   }
   
   // To (in screen)
   func (s *ListScreen) formatItem(item *Item) string {
       return fmt.Sprintf("%s - %s", item.Name, item.Date)
   }
   ```

5. **All view state**
   ```go
   // Extract this (from intent)
   selectedIndex int
   scrollOffset  int
   viewWidth     int
   
   // To (in screen)
   type ListScreen struct {
       *base.BaseScreen
       selectedIndex int
       scrollOffset  int
       // viewWidth from terminal info
   }
   ```

### NEVER Extract These (Keep in Intent)

1. **Business logic** (goes to context.go, not screen)
   ```go
   // Keep in context.go (NOT screen)
   func (c *Context) ValidateItem(item *Item) error { ... }
   func (c *Context) FilterItems(items []*Item) []*Item { ... }
   ```

2. **State machine logic** (stays in intent)
   ```go
   // Keep in intent
   func (i *MyIntent) transitionTo(state MyState) tea.Cmd { ... }
   ```

3. **Service calls** (stays in intent, delegates to context)
   ```go
   // Keep in intent
   func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
       ctx := i.GetContext()
       err := i.context.SaveItem(ctx, result.Data)
       // ...
   }
   ```

---

## Where to Extract

### Extraction Target Map

| What | Where | File |
|------|-------|------|
| List view rendering | ListScreen | `screens/{feature}/list_screen.go` |
| Detail view rendering | DetailScreen | `screens/{feature}/detail_screen.go` |
| Form view rendering | FormScreen | `screens/{feature}/form_screen.go` |
| Filter view rendering | FilterScreen | `screens/{feature}/filter_screen.go` |
| Header rendering | UIKit or Screen method | `primitives.Title()` or helper |
| Footer rendering | UIKit or Screen method | `primitives.HelpKeyBadge()` or helper |
| Modal rendering | Centralized modals | `uikit/feedback/*Modal` |
| Layout composition | UIKit | `layout.NewScreenLayout()` |
| Text styling | UIKit | `primitives.*()` functions |
| Colors | Theme | `theme.*()` methods |

### File Structure After Extraction

```
intents/{feature}/
└── intent.go              # Only View() that delegates (≤10 lines)

screens/{feature}/
├── list_screen.go         # Extracted from renderList()
├── detail_screen.go       # Extracted from renderDetail()
└── form_screen.go         # Extracted from renderForm()
```

---

## Extraction Process

### Step 1: Identify Render Methods

**Command**:
```bash
# Find all render methods
grep -n "^func.*render\|^func.*Render\|^func.*View" internal/cli/intents/my_intent.go

# Count render methods
grep -c "^func.*render\|^func.*Render\|^func.*View" internal/cli/intents/my_intent.go
```

**Example output**:
```
456:func (i *MyIntent) View() string {
512:func (i *MyIntent) renderHeader() string {
568:func (i *MyIntent) renderBody() string {
624:func (i *MyIntent) renderList() string {
680:func (i *MyIntent) renderDetail() string {
```

**Analysis**: 5 methods = VIOLATION (max 2)

---

### Step 2: Map Methods to Screens

**Create mapping table**:

| Method | Lines | Screen | Reason |
|--------|-------|--------|--------|
| `View()` | 15 | - | Keep (delegates) |
| `renderHeader()` | 25 | Remove | Use UIKit primitives |
| `renderBody()` | 40 | Remove | Use UIKit layout |
| `renderList()` | 120 | ListScreen | List view |
| `renderDetail()` | 95 | DetailScreen | Detail view |

---

### Step 3: Create Screen Files

**For each screen needed**:

```bash
# Create screen directory
mkdir -p screens/myfeature

# Create screen files
touch screens/myfeature/list_screen.go
touch screens/myfeature/detail_screen.go
```

**Screen template**:
```go
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
            primitives.HelpKeyBadge("esc", "back", theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

---

### Step 4: Move Rendering Code

**BEFORE (Intent)**:
```go
// File: my_intent.go (BEFORE extraction)
func (i *MyIntent) renderList() string {
    // 120 lines of lipgloss styling
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        Background(lipgloss.Color("#5c5cff")).
        Bold(true).
        Padding(1, 2)
    
    title := titleStyle.Render("My Feature")
    
    // Manual table rendering (40 lines)
    rows := []string{}
    for idx, item := range i.items {
        selected := ""
        if idx == i.selectedIndex {
            selected = "> "
        }
        rows = append(rows, fmt.Sprintf("%s%s - %s", selected, item.Name, item.Date))
    }
    
    body := strings.Join(rows, "\n")
    
    // Manual footer (20 lines)
    footer := "Press enter to select, esc to cancel"
    
    return lipgloss.JoinVertical(lipgloss.Left, title, body, footer)
}
```

**AFTER (Screen)**:
```go
// File: screens/myfeature/list_screen.go (AFTER extraction)
func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    // UIKit primitive (not raw lipgloss)
    header := primitives.Title("My Feature", theme)
    
    // TableBehavior (not manual rendering)
    tableView := s.table.View()
    
    // UIKit help badges (not raw strings)
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("esc", "cancel", theme)
    
    // UIKit layout (not manual JoinVertical)
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

**Result**: 120 lines → 15 lines (87% reduction!)

---

### Step 5: Update Intent View()

**BEFORE (Intent with rendering)**:
```go
// File: my_intent.go (BEFORE extraction)
func (i *MyIntent) View() string {
    switch i.state {
    case StateList:
        return i.renderList()      // 120 lines
    case StateDetail:
        return i.renderDetail()    // 95 lines
    case StateForm:
        return i.renderForm()      // 110 lines
    }
    return ""
}

// Total: 4 methods, 325+ lines of rendering
```

**AFTER (Intent delegates)**:
```go
// File: intents/myfeature/intent.go (AFTER extraction)
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    // Overlay modals (centralized)
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}

// Total: 1 method, ~10 lines (97% reduction!)
```

---

### Step 6: Update Intent Fields

**Add typed screen fields**:

```go
type MyIntent struct {
    *BaseIntent
    
    // Context
    context *MyContext
    
    // State machine
    state  MyState
    active bool
    
    // Screens (explicit typed fields - REQUIRED)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    formScreen   *myfeature.FormScreen
    
    // Active screen (generic pointer)
    activeScreen screens.Screen
    
    // Modals (centralized)
    deleteModal *feedback.ConfirmModal
}
```

---

### Step 7: Initialize Screens

**In Init() method**:

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

**In Update() method**:

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
    case *screens.NavigateResult:
        return i.HandleNavigate(r)
    case *screens.SubmitResult:
        return i.HandleSubmit(r)
    case *screens.CancelResult:
        return i.HandleCancel(r)
    case *screens.ErrorResult:
        return i.HandleError(r)
    }
    return nil
}
```

---

## View Method Limits

### Rule: Max 2 Render Methods

**Allowed configurations**:

1. **Configuration A (IDEAL)**:
   ```go
   // Only View() that delegates
   func (i *MyIntent) View() string {
       baseView := i.activeScreen.View()
       
       if i.modal != nil && i.modal.IsVisible() {
           return behaviors.RenderModalOverlay(i.modal, baseView)
       }
       
       return baseView
   }
   // Total: 1 method ✅
   ```

2. **Configuration B (ACCEPTABLE)**:
   ```go
   // View() + 1 tiny helper (≤5 lines)
   func (i *MyIntent) View() string {
       return i.renderWithModalOverlay(i.activeScreen.View())
   }
   
   func (i *MyIntent) renderWithModalOverlay(baseView string) string {
       if i.modal != nil && i.modal.IsVisible() {
           return behaviors.RenderModalOverlay(i.modal, baseView)
       }
       return baseView
   }
   // Total: 2 methods (helper ≤5 lines) ✅
   ```

3. **Configuration C (VIOLATION)**:
   ```go
   // View() + multiple helpers
   func (i *MyIntent) View() string { ... }
   func (i *MyIntent) renderHeader() string { ... }
   func (i *MyIntent) renderBody() string { ... }
   func (i *MyIntent) renderFooter() string { ... }
   // Total: 4 methods ❌ BLOCKED by Check #19
   ```

**Enforcement**: Check #19 blocks commits with >2 render methods

---

## UIKit Requirements

### Rule: No Raw Lipgloss in Intents

**NEVER in intents**:
```go
// ❌ VIOLATION - Raw lipgloss in intent
titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fff"))
title := titleStyle.Render("Title")
```

**ALWAYS in screens**:
```go
// ✅ CORRECT - UIKit in screen
header := primitives.Title("Title", theme)
```

### UIKit Component Map

**Use these in screens** (NOT intents):

| Need | UIKit Component | Never |
|------|-----------------|-------|
| Title | `primitives.Title("text", theme)` | `lipgloss.NewStyle().Bold()` |
| Body | `primitives.Body("text", theme)` | `lipgloss.NewStyle().Render()` |
| Error | `primitives.ErrorText("msg", theme)` | `lipgloss.Color("#ff0000")` |
| Success | `primitives.SuccessText("msg", theme)` | `lipgloss.Color("#00ff00")` |
| Help | `primitives.HelpKeyBadge("k", "l", theme)` | `components.KeyBadge()` |
| Layout | `layout.NewScreenLayout(theme)` | `lipgloss.JoinVertical()` |
| Modal | `behaviors.RenderModalOverlay()` | Custom modal rendering |

**See**: [UIKit Guide](../UIKIT_GUIDE.md)

---

## Screen Responsibilities

### What Screens MUST Do

1. **Embed BaseScreen**
   ```go
   type ListScreen struct {
       *base.BaseScreen  // REQUIRED
       // ... other fields
   }
   ```

2. **Implement Update()**
   ```go
   func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
       // Handle input
       // Return ScreenResult (Navigate, Submit, Cancel, Error)
   }
   ```

3. **Implement View()**
   ```go
   func (s *ListScreen) View() string {
       // Render using UIKit
       // Return string
   }
   ```

4. **Use UIKit components**
   - Primitives for text
   - Layout for composition
   - Theme for colors
   - Behaviors for complex UI (tables, forms)

5. **Return ScreenResult**
   ```go
   // Navigation
   return nil, screens.NewNavigateResult("detail", item)
   
   // Submission
   return nil, screens.NewSubmitResult(formData)
   
   // Cancellation
   return nil, screens.NewCancelResult("screen-name")
   
   // Error
   return nil, screens.NewErrorResult(err)
   ```

### What Screens MUST NOT Do

1. **Import intents package** (FORBIDDEN - circular dependency)
   ```go
   // ❌ NEVER do this
   import "github.com/baphled/kariya/internal/cli/intents"
   ```

2. **Call services directly** (use ScreenResult to communicate back)
   ```go
   // ❌ WRONG - Direct service call
   func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
       err := s.service.SaveItem(item)  // NO!
   }
   
   // ✅ CORRECT - Return result to intent
   func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
       return nil, screens.NewSubmitResult(item)  // Intent handles service call
   }
   ```

3. **Mutate intent state** (screens are stateless from intent perspective)
   ```go
   // ❌ WRONG - Mutating intent state
   func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
       i.state = StateDetail  // NO! Can't access intent
   }
   
   // ✅ CORRECT - Return navigation result
   func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
       return nil, screens.NewNavigateResult("detail", item)
   }
   ```

4. **Have business logic** (business logic goes in context.go)
   ```go
   // ❌ WRONG - Business logic in screen
   func (s *ListScreen) validateItem(item *Item) error { ... }
   
   // ✅ CORRECT - Business logic in context
   // File: context.go
   func (c *MyContext) ValidateItem(item *Item) error { ... }
   ```

---

## Common Violations

### Violation 1: Multiple Render Methods

```go
// ❌ VIOLATION - 6 render methods
func (i *MyIntent) View() string { ... }
func (i *MyIntent) renderHeader() string { ... }
func (i *MyIntent) renderBody() string { ... }
func (i *MyIntent) renderFooter() string { ... }
func (i *MyIntent) renderList() string { ... }
func (i *MyIntent) renderDetail() string { ... }

// ✅ CORRECT - 1 method that delegates
func (i *MyIntent) View() string {
    return i.activeScreen.View()
}
```

**Fix**: Extract all render*() methods to screens

---

### Violation 2: Lipgloss Styling in Intent

```go
// ❌ VIOLATION - Raw lipgloss in intent
func (i *MyIntent) View() string {
    titleStyle := lipgloss.NewStyle().Bold(true).Foreground(...)
    title := titleStyle.Render("Title")
    // ...
}

// ✅ CORRECT - UIKit in screen
func (s *ListScreen) View() string {
    title := primitives.Title("Title", theme)
    // ...
}
```

**Fix**: Move styling to screens, use UIKit primitives

---

### Violation 3: Switch for Rendering

```go
// ❌ VIOLATION - Rendering switch in intent
func (i *MyIntent) View() string {
    switch i.state {
    case StateList:
        return i.renderList()
    case StateDetail:
        return i.renderDetail()
    }
}

// ✅ CORRECT - Delegate to active screen
func (i *MyIntent) View() string {
    return i.activeScreen.View()
}
```

**Fix**: Extract each case to a screen, update active screen on state transitions

---

### Violation 4: Manual Layout in Intent

```go
// ❌ VIOLATION - Manual layout in intent
func (i *MyIntent) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        i.renderHeader(),
        i.renderBody(),
        i.renderFooter(),
    )
}

// ✅ CORRECT - UIKit layout in screen
func (s *ListScreen) View() string {
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(body).
        WithFooter(footer).
        Render()
}
```

**Fix**: Use UIKit layout in screens

---

### Violation 5: View State in Intent

```go
// ❌ VIOLATION - View state in intent
type MyIntent struct {
    *BaseIntent
    selectedIndex int      // View state
    scrollOffset  int      // View state
    cursorPos     int      // View state
}

// ✅ CORRECT - View state in screen
type ListScreen struct {
    *base.BaseScreen
    selectedIndex int      // Belongs in screen
    scrollOffset  int      // Belongs in screen
}
```

**Fix**: Move view state to screens

---

### Violation 6: Custom Modal Rendering

```go
// ❌ VIOLATION - Custom modal rendering in intent
func (i *MyIntent) View() string {
    baseView := i.screen.View()
    
    if i.showDeleteModal {
        // Custom modal rendering (50 lines)
        modalContent := lipgloss.NewStyle()...
        overlay := lipgloss.Place(...)
        return overlay
    }
    
    return baseView
}

// ✅ CORRECT - Centralized modal with RenderModalOverlay
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**Fix**: Use centralized modals from `uikit/feedback/` with `behaviors.RenderModalOverlay()`

---

### Violation 7: Formatting in Intent

```go
// ❌ VIOLATION - Formatting in intent
func (i *MyIntent) formatItem(item *Item) string {
    return fmt.Sprintf("%s - %s", item.Name, item.Date)
}

func (i *MyIntent) View() string {
    items := []string{}
    for _, item := range i.items {
        items = append(items, i.formatItem(item))
    }
    // ...
}

// ✅ CORRECT - Formatting in screen
func (s *ListScreen) formatItem(item *Item) string {
    return fmt.Sprintf("%s - %s", item.Name, item.Date)
}

func (s *ListScreen) View() string {
    // Use TableBehavior which handles formatting
    return s.table.View()
}
```

**Fix**: Move formatting to screens, prefer TableBehavior for structured data

---

## Validation

### Pre-Extraction Validation

**Before extracting views**:
```bash
# 1. Count render methods
grep -c "^func.*render\|^func.*Render\|^func.*View" my_intent.go

# Should be: >2 (extraction needed)

# 2. Find lipgloss usage
grep -c "lipgloss\." my_intent.go

# Should be: >0 (extraction needed)

# 3. Check file size
wc -l my_intent.go

# If >600 lines: MANDATORY extraction
# If >400 lines: RECOMMENDED extraction
```

### Post-Extraction Validation

**After extracting views**:
```bash
# 1. Architecture checks (REQUIRED)
make check-intent-architecture

# Check #19 should pass (≤2 render methods)

# 2. Count render methods (should be ≤2)
grep -c "^func.*render\|^func.*Render\|^func.*View" intents/myfeature/intent.go

# Expected: 1 (only View())

# 3. Check lipgloss usage (should be 0 in intent)
grep "lipgloss\." intents/myfeature/intent.go

# Expected: No matches

# 4. Verify screens created
ls screens/myfeature/

# Expected: list_screen.go, detail_screen.go, etc.

# 5. Run tests
make test

# Expected: All tests pass
```

### Manual Checklist

**After extraction**:
- [ ] Intent has ≤2 render methods (ideally 1)
- [ ] Intent View() ≤20 lines
- [ ] No lipgloss styling in intent
- [ ] No manual layout in intent
- [ ] Screen files created in `screens/{feature}/`
- [ ] Screens use UIKit components
- [ ] Screens embed `*base.BaseScreen`
- [ ] Screens implement Update() returning ScreenResult
- [ ] Screens implement View() using UIKit
- [ ] Intent delegates to activeScreen
- [ ] Modals use `behaviors.RenderModalOverlay()`
- [ ] All tests pass
- [ ] Check #19 passes

---

## Summary

### The Rules

1. **Max 2 render methods** in intent (ideally 1)
2. **No lipgloss styling** in intent
3. **No manual layout** in intent
4. **Delegate to screens** - intent orchestrates, screens render
5. **Use UIKit components** in screens
6. **Return ScreenResult** from screens
7. **Use centralized modals** with RenderModalOverlay
8. **Extract immediately** if any trigger is met

### The Process

1. Identify render methods (grep)
2. Map methods to screens
3. Create screen files
4. Move rendering code to screens
5. Update intent View() to delegate
6. Add typed screen fields
7. Initialize screens in Init()
8. Handle ScreenResults in Update()
9. Validate (checks + tests)

### The Outcome

**Before extraction**:
- Intent: 2,000+ lines with 8+ render methods
- Violations: Multiple (Check #19 fails)
- Maintainability: Low

**After extraction**:
- Intent: ~300 lines with 1 render method
- Screens: 3-5 screens, ~100 lines each
- Violations: Zero (Check #19 passes)
- Maintainability: High
- Code reduction: 48-87%

---

## Resources

### Documentation
- [Screen Extraction Guide](../guides/SCREEN_EXTRACTION_GUIDE.md) - Complete extraction guide
- [Migration Rules](MIGRATION_RULES.md) - Full migration rules
- [UIKit Guide](../UIKIT_GUIDE.md) - UIKit component reference
- [Intent Development Checklist](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - All checks

### Templates
- `examples/intent_subdirectory_template/screens/list_screen.go.template`

### Tools
- `make check-intent-architecture` - Validate extraction (Check #19)
- `grep "^func.*render" file.go` - Find render methods

### Enforcement
- **Check #19**: No rendering in intent (blocks >2 render methods)
- **Check #21**: UIKit component usage (blocks deprecated components)

---

**Last Updated**: 2026-01-22  
**Status**: MANDATORY - Must follow for all view extractions  
**Enforcement**: Check #19 (automated blocking)
