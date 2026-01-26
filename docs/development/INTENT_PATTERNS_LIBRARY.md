# Intent Patterns Library

**Last Updated**: 2026-01-20
**Status**: ✅ **PRODUCTION STANDARDS - COMPREHENSIVE PATTERN CATALOG**

> **UIKit Migration Note**: New code should use UIKit components from `internal/cli/uikit/`:
> - `layout.ScreenLayout` for view layouts (replaces `CreateStandardView`)
> - `feedback.Modal` + `behaviors.RenderModalOverlay()` for modals
> - `primitives.HelpKeyBadge()` for footer badges (replaces `components.KeyBadge`)
> - See [UIKIT_GUIDE.md](../UIKIT_GUIDE.md) for complete reference

---

## Table of Contents

1. [Overview](#overview)
2. [Pattern Index](#pattern-index)
3. [Core Rendering Patterns](#core-rendering-patterns)
4. [State Management Patterns](#state-management-patterns)
5. [Navigation Patterns](#navigation-patterns)
6. [Modal Patterns](#modal-patterns)
7. [Footer Patterns](#footer-patterns)
8. [Action Routing Patterns](#action-routing-patterns)
9. [Data Management Patterns](#data-management-patterns)
10. [Testing Patterns](#testing-patterns)
11. [Complete Intent Template](#complete-intent-template)
12. [Migration Checklist](#migration-checklist)

---

## Overview

This document catalogs **12 standardized patterns** discovered during the BrowseTimeline intent refactoring. All intents MUST follow these patterns for consistency, maintainability, and correctness.

**Reference Implementation**: `internal/cli/intents/browse_timeline/` (subdirectory structure)

**Applies To**: All intent implementations in KaRiya TUI

### File Organization (NEW - REQUIRED)

All new intents MUST use the subdirectory structure:

```
intents/{feature}/
├── context.go     # IntentContext + Validate() + domain types
├── result.go      # Result struct
├── constants.go   # State enum ONLY
├── messages.go    # ALL *Msg types
├── intent.go      # NewIntent, Init, Update, View, Result
├── types.go       # (optional) Intent struct if intent.go > 300 lines
├── handlers.go    # (optional) ScreenResultHandler methods
├── helpers.go     # (optional) Helper methods
└── interfaces.go  # (optional) Service interfaces

screens/{feature}/
├── list_screen.go
├── detail_screen.go
└── modals/        # Feature-specific modals
    ├── filter_modal.go
    └── helpers.go
```

**Reference**: See `intents/browse_timeline/` and `screens/timeline/` for complete example.

---

## Pattern Index

| # | Pattern Name | Priority | Applies To | Documentation |
|---|--------------|----------|------------|---------------|
| 1 | **Modal Overlay Rendering** | CRITICAL | All intents with modals | [MODAL_OVERLAY_PATTERN.md](MODAL_OVERLAY_PATTERN.md) |
| 2 | **Themed Footer Building** | HIGH | All intents | [THEMED_FOOTER_GUIDE.md](THEMED_FOOTER_GUIDE.md) |
| 3 | **View Rendering with Modal Overlay** | CRITICAL | All intents with modals | This doc + MODAL_OVERLAY_PATTERN |
| 4 | **Global Key Interception** | HIGH | All intents | This doc |
| 5 | **Context-Aware Footer Generation** | HIGH | All intents | This doc + THEMED_FOOTER_GUIDE |
| 6 | **State-to-Breadcrumb Mapping** | MEDIUM | All intents | This doc |
| 7 | **Screen Transition Helper** | MEDIUM | All intents | This doc |
| 8 | **Screen Result Handling** | HIGH | All intents with screens | This doc |
| 9 | **Filter/Sort Application** | MEDIUM | List-based intents | This doc |
| 10 | **Action Routing** | HIGH | List-based intents | This doc |
| 11 | **Delete Confirmation Flow** | MEDIUM | Intents with delete | This doc |
| 12 | **Form Modal with Immediate Init** | CRITICAL | Form-based modals | This doc + MODAL_OVERLAY_PATTERN |

---

## Core Rendering Patterns

### Pattern 1: Modal Overlay Rendering (CRITICAL)

**Purpose**: Render modals over complete views without misalignment

**Rule**: StandardView FIRST, modal overlay LAST

```go
func (i *YourIntent) View() string {
    // 1. Create StandardView with complete content
    view := i.CreateViewWithBreadcrumbs(
        i.GetState(),
        i.terminal.Width,
        i.terminal.Height,
    )
    view.WithContent(screen.RenderContent())
    view.WithHelp(i.getContextHelp())
    
    // 2. Render COMPLETE view
    baseView := view.Render()
    
    // 3. Overlay modal as FINAL step
    if i.modal != nil && i.modal.IsVisible() {
        return i.renderModalOverlay(baseView)
    }
    
    return baseView
}
```

**Why**: StandardView uses `lipgloss.Place()` to center content. Overlaying modal before StandardView causes double-centering and misalignment.

**See**: [MODAL_OVERLAY_PATTERN.md](MODAL_OVERLAY_PATTERN.md) for complete documentation

---

### Pattern 3: View Rendering with Modal Overlay

**Purpose**: Complete view rendering workflow

**Template**:

```go
func (i *YourIntent) View() string {
    // Handle special states first
    if i.state == StateError {
        return i.renderError()
    }
    
    // Render based on screen type
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        return i.renderListScreenView(screen)
    case *YourDetailScreen:
        return i.renderDetailScreenView(screen)
    case *YourConfirmScreen:
        return i.renderConfirmScreenView(screen)
    default:
        return "Unknown screen"
    }
}

func (i *YourIntent) renderListScreenView(screen *YourListScreen) string {
    // 1. Create StandardView
    view := i.CreateViewWithBreadcrumbs(
        i.GetState(),
        i.terminal.Width,
        i.terminal.Height,
    )
    
    // 2. Add content and help
    view.WithContent(screen.RenderContent())
    view.WithHelp(i.getContextHelp())
    
    // 3. Render complete view
    baseView := view.Render()
    
    // 4. Overlay modal if visible
    if i.filterModal != nil && i.filterModal.IsVisible() {
        return i.renderFilterModalOverlay(baseView)
    }
    
    return baseView
}
```

**Key Points**:
- One method per screen type for clarity
- StandardView always created with complete content
- Modal overlay always last step
- Help footer is context-aware (calls `getContextHelp()`)

---

## State Management Patterns

### Pattern 6: State-to-Breadcrumb Mapping

**Purpose**: Dynamic breadcrumbs based on current state

**Implementation**:

```go
func (i *YourIntent) GetState() string {
    switch i.currentScreen.(type) {
    case *YourListScreen:
        return "Browse"
    case *YourDetailScreen:
        return "Browse > Details"
    case *YourEditScreen:
        return "Browse > Edit"
    case *YourConfirmScreen:
        return "Browse > Confirm Delete"
    default:
        return "Browse"
    }
}
```

**Why**: Users see their location in the workflow. Breadcrumbs update automatically with state changes.

**Pattern**:
- Root state: Single word (e.g., "Browse")
- Nested states: Parent > Child (e.g., "Browse > Details")
- Action states: Parent > Action (e.g., "Browse > Edit")

---

### Pattern 7: Screen Transition Helper

**Purpose**: Consistent screen transitions with proper initialization

**Implementation**:

```go
func (i *YourIntent) transitionToScreen(screen screens.Screen) tea.Cmd {
    // Set screen properties
    screen.SetTerminalInfo(i.terminal)
    screen.SetTheme(i.theme)
    screen.SetLogo(i.logo)
    
    // Update current screen
    i.currentScreen = screen
    
    // Initialize screen
    return screen.Init()
}
```

**Usage**:

```go
// Transition to list screen
listScreen := NewYourListScreen(items)
return i.transitionToScreen(listScreen)

// Transition to detail screen
detailScreen := NewYourDetailScreen(selectedItem)
return i.transitionToScreen(detailScreen)
```

**Why**: Ensures all screens have required context (terminal, theme, logo) before initialization.

---

## Navigation Patterns

### Pattern 4: Global Key Interception

**Purpose**: Handle global keys before screen delegation

**Implementation**:

```go
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // 1. HIGHEST PRIORITY: Modal updates (if modal visible)
        if i.modal != nil && i.modal.IsVisible() {
            newModal, cmd := i.modal.Update(msg)
            i.modal = newModal.(*components.YourModal)
            
            // Handle modal result
            if !i.modal.IsVisible() {
                return i.handleModalClosed()
            }
            
            return cmd
        }
        
        // 2. MEDIUM PRIORITY: Global keys
        switch msg.String() {
        case "ctrl+c", "q":
            return i.handleQuit()
        case "?", "h":
            return i.handleHelp()
        case "m":
            return i.handleMainMenu()
        }
        
        // 3. LOWEST PRIORITY: Delegate to screen
        return i.handleScreenUpdate(msg)
        
    case tea.WindowSizeMsg:
        return i.handleWindowResize(msg)
        
    case screens.ScreenResultMsg:
        return i.handleScreenResult(msg)
        
    default:
        return nil
    }
}
```

**Priority Order**:
1. **Modal updates** - Highest (modals consume keys first)
2. **Global keys** - Medium (q, ?, m always work)
3. **Screen delegation** - Lowest (screen handles remaining keys)

**Why**: Prevents modals or screens from consuming global keys. Users can always quit, get help, or return to main menu.

---

### Pattern 8: Screen Result Handling

**Purpose**: Route screen results to appropriate handlers

**Implementation**:

```go
func (i *YourIntent) handleScreenResult(msg screens.ScreenResultMsg) tea.Cmd {
    switch result := msg.Result.(type) {
    case screens.ListActionResult:
        return i.handleListAction(result)
    case screens.NavigateBackResult:
        return i.handleNavigateBack()
    case screens.SelectItemResult:
        return i.handleItemSelected(result)
    case screens.DeleteConfirmedResult:
        return i.handleDeleteConfirmed(result)
    default:
        return nil
    }
}

func (i *YourIntent) handleListAction(result screens.ListActionResult) tea.Cmd {
    switch result.Action {
    case "add":
        return i.handleAddNew()
    case "edit":
        return i.handleEdit(result.ItemID)
    case "delete":
        return i.handleDelete(result.ItemID)
    case "filter":
        return i.handleOpenFilter()
    case "sort":
        return i.handleOpenSort()
    default:
        return nil
    }
}
```

**Why**: Type-safe routing based on result type. Each result type has dedicated handler. Extensible for new result types.

---

## Modal Patterns

### Pattern 12: Form Modal with Immediate Init

**Purpose**: Form modals render immediately on display

**Implementation**:

```go
func (i *YourIntent) handleOpenFilter() tea.Cmd {
    // 1. Create filter modal
    i.filterModal = components.NewFilterModal(
        i.theme,
        i.currentFilters,
        i.terminal.Width,
        i.terminal.Height,
    )
    
    // 2. CRITICAL: Call Init() for immediate rendering
    return i.filterModal.Init()
}
```

**Why**: Form lifecycle requires `Init()` to set up internal state. Without it, form appears blank until user input.

**When to Use**:
- ✅ All form-based modals (filter, edit, add)
- ✅ Any modal with huh.Form inside
- ❌ Simple confirmation modals (no form)

---

### Pattern 11: Delete Confirmation Flow

**Purpose**: Safe delete workflow with confirmation

**Implementation**:

```go
// Step 1: Show confirmation screen
func (i *YourIntent) handleDelete(itemID string) tea.Cmd {
    item, _ := i.findItemByID(itemID)
    
    confirmScreen := screens.NewDeleteConfirmScreen(
        "Delete Item",
        fmt.Sprintf("Are you sure you want to delete '%s'?", item.Name),
        itemID,
    )
    
    return i.transitionToScreen(confirmScreen)
}

// Step 2: Handle confirmation result
func (i *YourIntent) handleDeleteConfirmed(result screens.DeleteConfirmedResult) tea.Cmd {
    if !result.Confirmed {
        // User cancelled - go back to list
        return i.handleNavigateBack()
    }
    
    // User confirmed - delete item
    if err := i.service.Delete(result.ItemID); err != nil {
        return i.showError(err)
    }
    
    // Refresh list
    return i.refreshList()
}

// Step 3: Refresh list after delete
func (i *YourIntent) refreshList() tea.Cmd {
    items, err := i.service.GetAll()
    if err != nil {
        return i.showError(err)
    }
    
    listScreen := NewYourListScreen(items)
    return i.transitionToScreen(listScreen)
}
```

**Flow**:
1. User triggers delete (keyboard shortcut or action)
2. Confirmation screen displays with item details
3. User confirms or cancels
4. If confirmed: delete + refresh list
5. If cancelled: return to list unchanged

---

## Footer Patterns

### Pattern 2: Themed Footer Building

**Purpose**: Consistent, themed footers using KeyBadge components

**Implementation**:

```go
func (i *YourIntent) getContextHelp() string {
    // Different footer based on screen type
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        return i.getListScreenFooter()
    case *YourDetailScreen:
        return i.getDetailScreenFooter()
    default:
        return ""
    }
}

func (i *YourIntent) getListScreenFooter() string {
    badges := []components.KeyBadge{
        components.NavigateBadge(),         // ↑↓/j/k: Navigate
        components.NewKeyBadge("Enter", "View Details"),
        components.NewKeyBadge("a", "Add"),
        components.NewKeyBadge("e", "Edit"),
        components.NewKeyBadge("d", "Delete"),
        components.NewKeyBadge("f", "Filter"),
        components.BackBadge(),             // Esc: Back
        components.QuitBadge(),             // q: Quit
    }
    return components.RenderHelpFooter(i.theme, badges...)
}

func (i *YourIntent) getDetailScreenFooter() string {
    badges := []components.KeyBadge{
        components.NewKeyBadge("e", "Edit"),
        components.NewKeyBadge("d", "Delete"),
        components.BackBadge(),
        components.QuitBadge(),
    }
    return components.RenderHelpFooter(i.theme, badges...)
}
```

**Why**:
- Consistent styling via theme
- Reusable badge components
- Context-aware (different per screen)
- Professional appearance

**See**: [THEMED_FOOTER_GUIDE.md](THEMED_FOOTER_GUIDE.md) for complete documentation

---

### Pattern 5: Context-Aware Footer Generation

**Purpose**: Footer adapts to current screen and state

**States to Consider**:

| State | Footer Keys | Example |
|-------|-------------|---------|
| List Screen | Navigate, View, Add, Edit, Delete, Filter, Back, Quit | Browse items |
| Detail Screen | Edit, Delete, Back, Quit | View single item |
| Edit Screen | Save, Cancel | Edit form |
| Confirm Screen | Yes, No | Confirm delete |
| Loading Screen | None | Async operation |
| Error Screen | Retry, Back | Error occurred |

**Implementation**:

```go
func (i *YourIntent) getContextHelp() string {
    // Handle special states
    if i.state == StateLoading {
        return ""  // No footer during loading
    }
    
    if i.state == StateError {
        return i.getErrorFooter()
    }
    
    // Screen-specific footer
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        // Different footer if filtering is active
        if i.activeFilters != nil {
            return i.getFilteredListFooter()
        }
        return i.getListScreenFooter()
        
    case *YourDetailScreen:
        return i.getDetailScreenFooter()
        
    case *YourEditScreen:
        return i.getEditScreenFooter()
        
    case *YourConfirmScreen:
        return i.getConfirmScreenFooter()
        
    default:
        return ""
    }
}
```

---

## Action Routing Patterns

### Pattern 10: Action Routing

**Purpose**: Consistent action handling from list screens

**Implementation**:

```go
// Step 1: Screen emits action result
type ListActionResult struct {
    Action string
    ItemID string
}

// In YourListScreen Update()
case "a":
    return screens.ScreenResultMsg{
        Result: screens.ListActionResult{Action: "add"},
    }
case "e":
    return screens.ScreenResultMsg{
        Result: screens.ListActionResult{
            Action: "edit",
            ItemID: i.getSelectedItemID(),
        },
    }
case "d":
    return screens.ScreenResultMsg{
        Result: screens.ListActionResult{
            Action: "delete",
            ItemID: i.getSelectedItemID(),
        },
    }
case "f":
    return screens.ScreenResultMsg{
        Result: screens.ListActionResult{Action: "filter"},
    }

// Step 2: Intent routes action
func (i *YourIntent) handleListAction(result screens.ListActionResult) tea.Cmd {
    switch result.Action {
    case "add":
        return i.handleAddNew()
    case "edit":
        return i.handleEdit(result.ItemID)
    case "delete":
        return i.handleDelete(result.ItemID)
    case "filter":
        return i.handleOpenFilter()
    case "sort":
        return i.handleOpenSort()
    default:
        return nil
    }
}

// Step 3: Action handlers
func (i *YourIntent) handleAddNew() tea.Cmd {
    // Show add form/screen
}

func (i *YourIntent) handleEdit(itemID string) tea.Cmd {
    // Show edit form/screen for itemID
}

func (i *YourIntent) handleDelete(itemID string) tea.Cmd {
    // Show delete confirmation for itemID
}

func (i *YourIntent) handleOpenFilter() tea.Cmd {
    // Show filter modal
    i.filterModal = components.NewFilterModal(...)
    return i.filterModal.Init()
}
```

**Why**: Type-safe action routing. Screens emit actions, intent handles business logic. Easy to add new actions.

---

## Data Management Patterns

### Pattern 9: Filter/Sort Application

**Purpose**: Apply filters and sorting to data

**Implementation**:

```go
func (i *YourIntent) applyFilters(items []YourType) []YourType {
    if i.activeFilters == nil {
        return items
    }
    
    filtered := items
    
    // Apply each filter criterion
    if i.activeFilters.DateRange != nil {
        filtered = i.filterByDateRange(filtered, i.activeFilters.DateRange)
    }
    
    if len(i.activeFilters.Categories) > 0 {
        filtered = i.filterByCategories(filtered, i.activeFilters.Categories)
    }
    
    if i.activeFilters.SearchText != "" {
        filtered = i.filterBySearchText(filtered, i.activeFilters.SearchText)
    }
    
    // Apply sorting
    if i.activeSort != nil {
        filtered = i.applySorting(filtered, i.activeSort)
    }
    
    return filtered
}

func (i *YourIntent) filterByDateRange(items []YourType, dateRange *DateRange) []YourType {
    result := []YourType{}
    for _, item := range items {
        if item.Date.After(dateRange.Start) && item.Date.Before(dateRange.End) {
            result = append(result, item)
        }
    }
    return result
}

func (i *YourIntent) filterByCategories(items []YourType, categories []string) []YourType {
    categorySet := make(map[string]bool)
    for _, cat := range categories {
        categorySet[cat] = true
    }
    
    result := []YourType{}
    for _, item := range items {
        for _, cat := range item.Categories {
            if categorySet[cat] {
                result = append(result, item)
                break
            }
        }
    }
    return result
}

func (i *YourIntent) applySorting(items []YourType, sortBy string) []YourType {
    sorted := make([]YourType, len(items))
    copy(sorted, items)
    
    switch sortBy {
    case "date_asc":
        sort.Slice(sorted, func(i, j int) bool {
            return sorted[i].Date.Before(sorted[j].Date)
        })
    case "date_desc":
        sort.Slice(sorted, func(i, j int) bool {
            return sorted[i].Date.After(sorted[j].Date)
        })
    case "name_asc":
        sort.Slice(sorted, func(i, j int) bool {
            return sorted[i].Name < sorted[j].Name
        })
    case "name_desc":
        sort.Slice(sorted, func(i, j int) bool {
            return sorted[i].Name > sorted[j].Name
        })
    }
    
    return sorted
}
```

**Usage**:

```go
// After loading data
items, err := i.service.GetAll()
if err != nil {
    return i.showError(err)
}

// Apply filters/sorting
filtered := i.applyFilters(items)

// Create screen with filtered data
listScreen := NewYourListScreen(filtered)
return i.transitionToScreen(listScreen)
```

---

## Testing Patterns

### Unit Test Template

```go
var _ = Describe("YourIntent", func() {
    var (
        ctx     context.Context
        intent  *YourIntent
        service *mocks.MockService
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        service = mocks.NewMockService()
        intent = NewYourIntent(ctx, service)
        intent.terminal = &cli.TerminalInfo{Width: 120, Height: 40}
        intent.theme = themes.CatppuccinMocha()
        intent.logo = "KaRiya"
    })
    
    Describe("View()", func() {
        Context("with list screen", func() {
            It("should render StandardView with content", func() {
                view := intent.View()
                
                Expect(view).To(ContainSubstring("KaRiya"))  // Logo
                Expect(view).To(ContainSubstring("Browse"))   // Breadcrumb
                Expect(view).NotTo(BeEmpty())
            })
            
            It("should overlay modal when visible", func() {
                intent.filterModal = components.NewFilterModal(...)
                intent.filterModal.Show()
                
                view := intent.View()
                
                Expect(view).To(ContainSubstring("Filter"))
                Expect(view).To(ContainSubstring("Expected Content"))
            })
        })
    })
    
    Describe("Update()", func() {
        Context("global keys", func() {
            It("should handle quit", func() {
                msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
                
                cmd := intent.Update(msg)
                
                Expect(cmd).NotTo(BeNil())
                // Verify quit command
            })
            
            It("should handle main menu", func() {
                msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
                
                cmd := intent.Update(msg)
                
                Expect(cmd).NotTo(BeNil())
                // Verify main menu navigation
            })
        })
        
        Context("modal visible", func() {
            BeforeEach(func() {
                intent.filterModal = components.NewFilterModal(...)
                intent.filterModal.Show()
            })
            
            It("should route keys to modal first", func() {
                msg := tea.KeyMsg{Type: tea.KeyTab}
                
                cmd := intent.Update(msg)
                
                // Modal should consume tab, not screen
                Expect(cmd).NotTo(BeNil())
            })
        })
    })
    
    Describe("Action Routing", func() {
        It("should handle add action", func() {
            result := screens.ListActionResult{Action: "add"}
            
            cmd := intent.handleListAction(result)
            
            Expect(cmd).NotTo(BeNil())
            // Verify add screen shown
        })
        
        It("should handle delete action", func() {
            result := screens.ListActionResult{
                Action: "delete",
                ItemID: "test-id",
            }
            
            cmd := intent.handleListAction(result)
            
            Expect(cmd).NotTo(BeNil())
            // Verify delete confirmation shown
        })
    })
})
```

---

## Complete Intent Template (Subdirectory Structure)

Use this template for new intents. All types are in separate files per the subdirectory pattern.

### File: intents/myfeature/constants.go
```go
package myfeature

type State string

const (
    StateList    State = "list"
    StateDetail  State = "detail"
)
```

### File: intents/myfeature/context.go
```go
package myfeature

type IntentContext struct {
    Items   []*domain.Item
    Service ItemService
}

func (c *IntentContext) Validate() error {
    if c.Items == nil {
        c.Items = make([]*domain.Item, 0)
    }
    return nil
}
```

### File: intents/myfeature/result.go
```go
package myfeature

type Result struct {
    SelectedItem *domain.Item
}
```

### File: intents/myfeature/messages.go
```go
package myfeature

type ItemSelectedMsg struct {
    Item *domain.Item
}
```

### File: intents/myfeature/types.go
```go
package myfeature

type Intent struct {
    *intents.BaseIntent
    
    context      *IntentContext
    state        State
    active       bool
    result       *intents.IntentResult[*Result]
    activeScreen screens.Screen
    
    // Modals (from screens/myfeature/modals/ or uikit/feedback/)
    filterModal *modals.FilterModal
    errorModal  *feedback.Modal
    
    // Modal registry for unified handling
    modalRegistry *intents.ModalRegistry
}
```

### File: intents/myfeature/intent.go
```go
package myfeature

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/cli/screens"
)

var _ intents.ScreenResultHandler = (*Intent)(nil)

// NewIntent creates a new myfeature intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
    if err := ctx.Validate(); err != nil {
        return nil, err
    }
    return &Intent{
        BaseIntent:    intents.NewBaseIntent(),
        context:       ctx,
        state:         StateList,
        active:        true,
        modalRegistry: intents.NewModalRegistry(),
    }, nil
}

// Init implements Intent.Init
func (i *YourIntent) Init(ctx context.Context) tea.Cmd {
    return i.loadInitialData()
}

// Update implements Intent.Update
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // PATTERN 4: Global Key Interception
        
        // 1. HIGHEST PRIORITY: Modal updates
        if i.filterModal != nil && i.filterModal.IsVisible() {
            return i.handleModalUpdate(msg)
        }
        
        // 2. MEDIUM PRIORITY: Global keys
        switch msg.String() {
        case "ctrl+c", "q":
            return i.handleQuit()
        case "?", "h":
            return i.handleHelp()
        case "m":
            return i.handleMainMenu()
        }
        
        // 3. LOWEST PRIORITY: Screen delegation
        return i.handleScreenUpdate(msg)
        
    case tea.WindowSizeMsg:
        return i.handleWindowResize(msg)
        
    case screens.ScreenResultMsg:
        // PATTERN 8: Screen Result Handling
        return i.handleScreenResult(msg)
        
    default:
        return nil
    }
}

// View implements Intent.View
func (i *YourIntent) View() string {
    // PATTERN 3: View Rendering with Modal Overlay
    
    // Handle special states
    if i.state == StateLoading {
        return i.renderLoading()
    }
    
    if i.state == StateError {
        return i.renderError()
    }
    
    // Render based on screen type
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        return i.renderListScreenView(screen)
    case *YourDetailScreen:
        return i.renderDetailScreenView(screen)
    case *YourEditScreen:
        return i.renderEditScreenView(screen)
    default:
        return "Unknown screen"
    }
}

// PATTERN 1 & 3: Modal Overlay Rendering
func (i *YourIntent) renderListScreenView(screen *YourListScreen) string {
    // 1. Create StandardView
    view := i.CreateViewWithBreadcrumbs(
        i.GetState(),
        i.terminal.Width,
        i.terminal.Height,
    )
    
    // 2. Add content and help
    view.WithContent(screen.RenderContent())
    view.WithHelp(i.getContextHelp())
    
    // 3. Render complete view
    baseView := view.Render()
    
    // 4. Overlay modal (if visible)
    if i.filterModal != nil && i.filterModal.IsVisible() {
        return i.renderFilterModalOverlay(baseView)
    }
    
    return baseView
}

func (i *YourIntent) renderFilterModalOverlay(baseView string) string {
    if i.filterModal == nil || !i.filterModal.IsVisible() {
        return baseView
    }
    
    modalContent := i.filterModal.View()
    
    return components.RenderOverlay(
        baseView,
        modalContent,
        i.terminal.Width,
        i.terminal.Height,
    )
}

// PATTERN 6: State-to-Breadcrumb Mapping
func (i *YourIntent) GetState() string {
    switch i.currentScreen.(type) {
    case *YourListScreen:
        return "Browse"
    case *YourDetailScreen:
        return "Browse > Details"
    case *YourEditScreen:
        return "Browse > Edit"
    default:
        return "Browse"
    }
}

// PATTERN 2 & 5: Context-Aware Footer Generation
func (i *YourIntent) getContextHelp() string {
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        return i.getListScreenFooter()
    case *YourDetailScreen:
        return i.getDetailScreenFooter()
    default:
        return ""
    }
}

// PATTERN 2: Themed Footer Building
func (i *YourIntent) getListScreenFooter() string {
    badges := []components.KeyBadge{
        components.NavigateBadge(),
        components.NewKeyBadge("Enter", "View"),
        components.NewKeyBadge("a", "Add"),
        components.NewKeyBadge("e", "Edit"),
        components.NewKeyBadge("d", "Delete"),
        components.NewKeyBadge("f", "Filter"),
        components.BackBadge(),
        components.QuitBadge(),
    }
    return components.RenderHelpFooter(i.theme, badges...)
}

// PATTERN 7: Screen Transition Helper
func (i *YourIntent) transitionToScreen(screen screens.Screen) tea.Cmd {
    screen.SetTerminalInfo(i.terminal)
    screen.SetTheme(i.theme)
    screen.SetLogo(i.logo)
    
    i.currentScreen = screen
    
    return screen.Init()
}

// PATTERN 8: Screen Result Handling
func (i *YourIntent) handleScreenResult(msg screens.ScreenResultMsg) tea.Cmd {
    switch result := msg.Result.(type) {
    case screens.ListActionResult:
        return i.handleListAction(result)
    case screens.NavigateBackResult:
        return i.handleNavigateBack()
    case screens.SelectItemResult:
        return i.handleItemSelected(result)
    case screens.DeleteConfirmedResult:
        return i.handleDeleteConfirmed(result)
    default:
        return nil
    }
}

// PATTERN 10: Action Routing
func (i *YourIntent) handleListAction(result screens.ListActionResult) tea.Cmd {
    switch result.Action {
    case "add":
        return i.handleAddNew()
    case "edit":
        return i.handleEdit(result.ItemID)
    case "delete":
        return i.handleDelete(result.ItemID)
    case "filter":
        return i.handleOpenFilter()
    default:
        return nil
    }
}

// PATTERN 12: Form Modal with Immediate Init
func (i *YourIntent) handleOpenFilter() tea.Cmd {
    i.filterModal = components.NewFilterModal(
        i.theme,
        i.activeFilters,
        i.terminal.Width,
        i.terminal.Height,
    )
    
    // CRITICAL: Call Init() for immediate rendering
    return i.filterModal.Init()
}

// PATTERN 11: Delete Confirmation Flow
func (i *YourIntent) handleDelete(itemID string) tea.Cmd {
    item, _ := i.findItemByID(itemID)
    
    confirmScreen := screens.NewDeleteConfirmScreen(
        "Delete Item",
        fmt.Sprintf("Are you sure you want to delete '%s'?", item.Name),
        itemID,
    )
    
    return i.transitionToScreen(confirmScreen)
}

func (i *YourIntent) handleDeleteConfirmed(result screens.DeleteConfirmedResult) tea.Cmd {
    if !result.Confirmed {
        return i.handleNavigateBack()
    }
    
    if err := i.service.Delete(result.ItemID); err != nil {
        return i.showError(err)
    }
    
    return i.refreshList()
}

// PATTERN 9: Filter/Sort Application
func (i *YourIntent) applyFilters(items []YourType) []YourType {
    if i.activeFilters == nil {
        return items
    }
    
    filtered := items
    
    // Apply filter criteria
    if len(i.activeFilters.Categories) > 0 {
        filtered = i.filterByCategories(filtered, i.activeFilters.Categories)
    }
    
    // Apply sorting
    if i.activeSort != nil {
        filtered = i.applySorting(filtered, i.activeSort)
    }
    
    return filtered
}

// Helper methods
func (i *YourIntent) loadInitialData() tea.Cmd {
    i.state = StateLoading
    
    return func() tea.Msg {
        items, err := i.service.GetAll()
        if err != nil {
            return errorMsg{err}
        }
        return dataLoadedMsg{items}
    }
}

func (i *YourIntent) refreshList() tea.Cmd {
    items, err := i.service.GetAll()
    if err != nil {
        return i.showError(err)
    }
    
    filtered := i.applyFilters(items)
    listScreen := NewYourListScreen(filtered)
    return i.transitionToScreen(listScreen)
}

// Result implements Intent.Result
func (i *YourIntent) Result() *IntentResult[interface{}] {
    // Return intent result
    return nil
}
```

---

## Migration Checklist

When refactoring an existing intent to use these patterns:

### Pre-Migration
- [ ] Read MODAL_OVERLAY_PATTERN.md
- [ ] Read THEMED_FOOTER_GUIDE.md
- [ ] Review BrowseTimeline as reference implementation
- [ ] Identify all modals used in intent
- [ ] List all screen types in intent
- [ ] Note existing keyboard shortcuts

### Pattern 1-3: View Rendering
- [ ] Move modal overlay to AFTER StandardView.Render()
- [ ] Separate rendering per screen type (one method per type)
- [ ] Ensure StandardView gets complete content (not placeholders)
- [ ] Add `renderModalOverlay()` helper methods

### Pattern 4: Key Handling
- [ ] Implement three-tier key handling (modal → global → screen)
- [ ] Verify global keys (q, ?, m) always work
- [ ] Test modal key consumption (tab, enter, esc)

### Pattern 2 & 5: Footers
- [ ] Convert all plain text footers to KeyBadge components
- [ ] Implement `getContextHelp()` method
- [ ] Create per-screen footer methods
- [ ] Use theme parameter in RenderHelpFooter

### Pattern 6-7: State Management
- [ ] Implement `GetState()` for breadcrumbs
- [ ] Create `transitionToScreen()` helper
- [ ] Ensure screens get terminal/theme/logo before Init()

### Pattern 8 & 10: Result Handling
- [ ] Implement `handleScreenResult()` router
- [ ] Create `handleListAction()` for action routing
- [ ] Add dedicated handlers per action type

### Pattern 12: Form Modals
- [ ] Call `Init()` when creating form modals
- [ ] Verify immediate rendering (no blank forms)

### Pattern 9 & 11: Data Management
- [ ] Extract `applyFilters()` method
- [ ] Implement delete confirmation flow
- [ ] Add filter/sort helpers

### Testing
- [ ] Update tests for new View() structure
- [ ] Test modal visibility and overlay
- [ ] Test global key interception
- [ ] Test footer rendering
- [ ] Visual test at multiple terminal sizes

### Documentation
- [ ] Update intent-specific documentation
- [ ] Add examples to guides if novel patterns used
- [ ] Note any deviations from standard patterns

---

## Related Documentation

- **[MODAL_OVERLAY_PATTERN.md](MODAL_OVERLAY_PATTERN.md)** - Critical modal rendering pattern
- **[THEMED_FOOTER_GUIDE.md](THEMED_FOOTER_GUIDE.md)** - Footer and KeyBadge usage
- **[TUI_STANDARDS.md](../TUI_STANDARDS.md)** - Overall TUI standards
- **[TUI_DEVELOPER_GUIDE.md](../TUI_DEVELOPER_GUIDE.md)** - TUI development guide

---

**Status**: ✅ **PRODUCTION STANDARDS - REQUIRED FOR ALL INTENTS**

*Reference implementation: `intents/browse_timeline/` subdirectory (2026-01-24)*  
*Migration guide: `docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md`*
