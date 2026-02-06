---
name: create-screen
description: Create a new screen component following KaRiya naming conventions and architecture
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Create a new screen component with proper naming, structure, and architecture compliance.

## When to use me

Use this skill when adding a new screen (list, detail, form, etc.) to an intent.

## Screen Naming Conventions

| Aspect | Convention | Example |
|--------|------------|---------|
| Package | `screens/{feature_name}/` (underscore) | `screens/fact_management/` |
| File | `{type}.go` | `list.go`, `detail.go`, `form.go` |
| Struct | `{Entity}{Type}Screen` | `FactListScreen`, `FactDetailScreen` |
| Constructor | `New{StructName}()` | `NewFactListScreen()` |

**Screen Types:** `List`, `Detail`, `Form`, `Delete`, `Select`, `Confirm`, `Preview`, `Review`

## Screen Structure

```go
// Package screens/{feature} provides UI screens for {feature} management.
package myfeature

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/screens/base"
    tea "github.com/charmbracelet/bubbletea"
)

// MyListScreen displays a list of items with filtering and selection.
//
// Expected: theme (required), items (can be empty slice)
// Returns: ScreenResult on user actions
// Side effects: None
type MyListScreen struct {
    *base.BaseScreen
    table *behaviors.TableBehavior[*MyItem]
}

// NewMyListScreen creates a new list screen.
//
// Expected: theme (required), items (can be empty)
// Returns: initialized screen ready for display
// Side effects: None
func NewMyListScreen(theme themes.Theme, items []*MyItem) *MyListScreen {
    s := &MyListScreen{
        BaseScreen: base.NewBaseScreen(theme),
        table:      behaviors.NewTableBehavior(items, nil),
    }
    return s
}

// Update handles input and returns screen results.
//
// Expected: tea.Msg from Bubble Tea runtime
// Returns: tea.Cmd and optional ScreenResult
// Side effects: May update internal selection state
func (s *MyListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEsc:
            return nil, screens.NewCancelResult("")
        case tea.KeyEnter:
            if item := s.table.SelectedItem(); item != nil {
                return nil, screens.NewNavigateResult("detail", item)
            }
        }
    }
    
    cmd := s.table.Update(msg)
    return cmd, nil
}

// View renders the screen.
//
// Expected: None
// Returns: rendered string for display
// Side effects: None
func (s *MyListScreen) View() string {
    return s.table.View()
}
```

## Keyboard Handling Rules

**ALWAYS use `tea.Key*` constants for special keys:**

```go
// GOOD
switch msg.Type {
case tea.KeyEsc:
    return nil, &screens.CancelResult{}
case tea.KeyUp:
    s.table.HandleNavigation("up")
case tea.KeyDown:
    s.table.HandleNavigation("down")
case tea.KeyEnter:
    // Handle enter
}

// String comparison ONLY for vim keys
switch msg.String() {
case "j":
    s.table.HandleNavigation("down")
case "k":
    s.table.HandleNavigation("up")
case "q":
    // Quit
}
```

## Screen Result Pattern

Screens communicate with intents via ScreenResult:

```go
// Return results, don't mutate intent state
return nil, screens.NewCancelResult("reason")
return nil, screens.NewNavigateResult("detail", selectedItem)
return nil, screens.NewSubmitResult(formData)
return nil, screens.NewErrorResult(err)
```

## Modal Naming

| Aspect | Convention | Example |
|--------|------------|---------|
| Package | `screens/{feature}/modals/` | `screens/fact_management/modals/` |
| File | `{action}_modal.go` | `edit_modal.go`, `filter_modal.go` |
| Struct | `{Action}Modal` | `EditModal`, `FilterModal` |

**Allowed Modal Locations:**
- `internal/cli/uikit/feedback/` (reusable modals)
- `internal/cli/screens/{feature}/modals/` (feature-specific)

**FORBIDDEN Modal Locations:**
- `internal/cli/intents/` - NEVER
- `internal/cli/models/` - NEVER

## Components to Use

| Need | Use | Not |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Colors | `theme.Primary()` etc | `lipgloss.Color("#xxx")` |
| Layout | `layout.ScreenLayout` | Manual composition |
| Badges | `primitives.HelpKeyBadge()` | `components.KeyBadge` |

## Related skills

- `create-intent` - Create the parent intent
- `bubble-tea-expert` - Bubble Tea TUI patterns (Update returns `(tea.Cmd, ScreenResult)`)
- `component-lookup` - Find the right component to use
- `fix-architecture` - Fix violations
