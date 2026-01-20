# Intent Architecture Guide

**Purpose**: Comprehensive guide to KaRiya's intent-driven TUI architecture  
**Last Updated**: 2026-01-20  
**Status**: Production Ready

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Layer Hierarchy](#layer-hierarchy)
3. [Intent Layer](#intent-layer)
4. [Screen Layer](#screen-layer)
5. [Modal Layer](#modal-layer)
6. [UIKit Foundation](#uikit-foundation)
7. [Dependency Flow](#dependency-flow)
8. [Creating New Intents](#creating-new-intents)
9. [Migration Guide](#migration-guide)
10. [Best Practices](#best-practices)

---

## Architecture Overview

KaRiya uses a **layered, intent-driven architecture** for its TUI:

```
┌─────────────────────────────────────────────────────────────┐
│                      App (Router)                           │
│  Routes user actions to intents, manages global state       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Intent Layer                             │
│  State machine, business logic orchestration, results       │
│  Examples: BrowseTimelineIntent, CaptureEventIntent         │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  Screen Layer   │  │   Modal Layer   │  │  Direct Render  │
│  Reusable views │  │  Overlay dialogs│  │  (Legacy/Simple)│
└─────────────────┘  └─────────────────┘  └─────────────────┘
              │               │               │
              └───────────────┼───────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    UIKit Foundation                         │
│  Primitives, Containers, Layout, Feedback, Theme            │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

1. **Intents own state machines** - All state transitions happen in intents
2. **Screens are stateless views** - Screens render UI and return results
3. **Modals overlay screens** - Modals appear over screens without replacing them
4. **UIKit provides consistency** - All rendering uses UIKit components
5. **One-way data flow** - Data flows down, events flow up

---

## Layer Hierarchy

### Package Structure

```
internal/cli/
├── app/                    # Application shell and router
│   └── app.go             # Root BubbleTea model
│
├── intents/               # Intent implementations
│   ├── base_intent.go     # BaseIntent with common helpers
│   ├── contract.go        # Intent interface definitions
│   ├── result.go          # IntentResult types
│   ├── router.go          # Intent routing
│   └── {name}_intent.go   # Intent implementations
│
├── screens/               # Screen implementations
│   ├── base/              # Base screen types
│   │   ├── base_screen.go # BaseScreen with theme/terminal
│   │   ├── select_screen.go
│   │   ├── form_screen.go
│   │   └── detail_screen.go
│   ├── capture/           # CaptureEvent screens
│   ├── configure/         # ConfigureSystem screens
│   ├── cv/                # GenerateCV screens
│   ├── skills/            # ManageSkills screens
│   └── timeline/          # BrowseTimeline screens
│
├── components/            # Modal components
│   ├── delete_confirm_modal.go
│   ├── filter_modal.go
│   └── {domain}_modal.go
│
├── uikit/                 # UIKit component library
│   ├── primitives/        # Text, Button, Badge, Input
│   ├── containers/        # Box, Overlay
│   ├── layout/            # ScreenLayout, Header, Footer
│   ├── feedback/          # Modal types (Error, Loading, etc.)
│   ├── navigation/        # Breadcrumbs
│   ├── display/           # Logo
│   └── theme/             # Theme infrastructure
│
├── behaviors/             # Embeddable behaviors
│   ├── table.go           # TableBehavior[T]
│   └── modal_helpers.go   # RenderModalOverlay
│
└── themes/                # Theme definitions
    └── default.go         # Catppuccin Macchiato theme
```

### Dependency Rules

```
app → intents → screens → uikit
         │         │
         ├─────────┼──→ components (modals)
         │         │
         └─────────┴──→ behaviors
```

**NEVER**:
- Screens importing intents
- UIKit importing screens/intents
- Behaviors importing screens/intents
- Circular dependencies

---

## Intent Layer

### Intent Interface

```go
// Intent defines the contract for all intents
type Intent interface {
    Init() tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    View() string
    Result() *IntentResult[interface{}]
}
```

### BaseIntent

All intents embed `*BaseIntent` for common functionality:

```go
type BaseIntent struct {
    terminalInfo *TerminalInfo
    theme        themes.Theme
    logo         *display.Logo
    helpModal    *feedback.HelpModal
}

// Key methods:
func (b *BaseIntent) GetTerminalInfo() *TerminalInfo
func (b *BaseIntent) Theme() themes.Theme
func (b *BaseIntent) GetLogo() *display.Logo
func (b *BaseIntent) CreateViewWithBreadcrumbs(breadcrumbs, content, help string) string
```

### Intent States

Each intent defines its own state enum:

```go
type BrowseTimelineState string

const (
    StateList    BrowseTimelineState = "list"
    StateDetail  BrowseTimelineState = "detail"
    StateDeleting BrowseTimelineState = "deleting"
)
```

### Intent Model Structure

```go
type BrowseTimelineIntent struct {
    *BaseIntent
    
    // State machine
    state  BrowseTimelineState
    active bool
    result *IntentResult[*BrowseTimelineResult]
    
    // Screens (one per state that uses screens)
    listScreen   *timeline.TimelineEventListScreen
    detailScreen *timeline.EventDetailScreen
    
    // Modals (shared across states)
    filterModal     *components.FilterModalModel
    deleteModal     *components.DeleteConfirmModal
    quickAddModal   *components.QuickAddEventModal
    
    // Data
    events        []*career.CareerEvent
    selectedEvent *career.CareerEvent
}
```

### Screen Result Handling

Intents implement `ScreenResultHandler` to process screen results:

```go
type ScreenResultHandler interface {
    HandleCancel(*screens.CancelResult) tea.Cmd
    HandleNavigate(*screens.NavigateResult) tea.Cmd
    HandleSubmit(*screens.SubmitResult) tea.Cmd
    HandleError(*screens.ErrorResult) tea.Cmd
}
```

Example implementation:

```go
func (i *BrowseTimelineIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    switch result.Target {
    case "detail":
        i.state = StateDetail
        i.detailScreen = timeline.NewEventDetailScreen(i.selectedEvent)
        return nil
    case "filter":
        i.filterModal.Show()
        return nil
    }
    return nil
}

func (i *BrowseTimelineIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    i.setCancelled()
    return nil
}
```

---

## Screen Layer

### Screen Interface

```go
type Screen interface {
    Update(msg tea.Msg) (tea.Cmd, ScreenResult)
    View() string
    SetTerminalInfo(info *intents.TerminalInfo)
    SetTheme(theme themes.Theme)
    SetLogo(logo *display.Logo)
}
```

### BaseScreen

All screens embed `*base.BaseScreen`:

```go
type BaseScreen struct {
    terminalInfo *intents.TerminalInfo
    theme        themes.Theme
    logo         *display.Logo
    breadcrumbs  string
}

// Key methods:
func (b *BaseScreen) CreateView(breadcrumbs, content, footer string) string
func (b *BaseScreen) GetDimensions() (width, height int)
func (b *BaseScreen) Theme() themes.Theme
```

### Screen Result Types

Screens return typed results to intents:

```go
type ScreenResult interface {
    Type() ScreenResultType
}

type NavigateResult struct {
    Target   string                 // e.g., "detail", "edit"
    Data     interface{}            // Optional payload
    Metadata map[string]interface{} // Context preservation
}

type CancelResult struct {
    Reason string
}

type SubmitResult struct {
    Data interface{}
}

type ErrorResult struct {
    Error   error
    Code    string
    Message string
}
```

### Screen Implementation Pattern

```go
package timeline

type TimelineEventListScreen struct {
    *base.BaseScreen
    
    // Behaviors (embeddable)
    table *behaviors.TableBehavior[*career.CareerEvent]
    
    // Local state
    events   []*career.CareerEvent
    selected int
}

func NewTimelineEventListScreen(events []*career.CareerEvent) *TimelineEventListScreen {
    s := &TimelineEventListScreen{
        BaseScreen: base.NewBaseScreen(),
        events:     events,
    }
    
    // Configure table behavior
    s.table = behaviors.NewTableBehavior[*career.CareerEvent]().
        WithColumns([]behaviors.ColumnDef{
            {Key: "date", Title: "Date", Width: 12},
            {Key: "title", Title: "Title", Width: 40},
            {Key: "company", Title: "Company", Width: 20},
        }).
        WithItems(events)
    
    return s
}

func (s *TimelineEventListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            return nil, screens.NewCancelResult("user cancelled")
        case "enter":
            selected := s.table.Selected()
            return nil, screens.NewNavigateResult("detail", selected)
        case "f":
            return nil, screens.NewNavigateResult("filter", nil)
        }
    }
    
    // Delegate to table behavior
    cmd := s.table.Update(msg)
    return cmd, nil
}

func (s *TimelineEventListScreen) View() string {
    content := s.table.View()
    footer := "Enter: View | f: Filter | Esc: Back"
    
    return s.CreateView(s.breadcrumbs, content, footer)
}
```

---

## Modal Layer

### Modal Interface

Modals are overlays that appear over screens:

```go
type Modal interface {
    tea.Model
    IsVisible() bool
    Show()
    Hide()
}
```

### Modal Implementation Pattern

```go
package components

type DeleteConfirmModal struct {
    visible   bool
    theme     themes.Theme
    width     int
    height    int
    itemName  string
    confirmed bool
}

func NewDeleteConfirmModal(theme themes.Theme, itemName string) *DeleteConfirmModal {
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    return &DeleteConfirmModal{
        theme:    theme,
        itemName: itemName,
        width:    50,
    }
}

func (m *DeleteConfirmModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if !m.visible {
        return m, nil
    }
    
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = min(msg.Width-4, 60)
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "y", "Y":
            m.confirmed = true
            m.Hide()
        case "n", "N", "esc":
            m.confirmed = false
            m.Hide()
        }
    }
    return m, nil
}

func (m *DeleteConfirmModal) View() string {
    if !m.visible {
        return ""
    }
    
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    
    title := primitives.Title("Confirm Delete", theme).Render()
    message := primitives.Body(
        fmt.Sprintf("Delete %q? This cannot be undone.", m.itemName),
        theme,
    ).Render()
    hint := primitives.Muted("y: Yes | n: No | Esc: Cancel", theme).Render()
    
    content := lipgloss.JoinVertical(lipgloss.Left, title, "", message, "", hint)
    
    // CRITICAL: Solid background for overlay
    return containers.NewBox(theme).
        Content(content).
        Width(m.width).
        Padding(1).
        Background(theme.BackgroundColor()).
        Variant(containers.BoxDestructive).
        Render()
}

func (m *DeleteConfirmModal) IsVisible() bool { return m.visible }
func (m *DeleteConfirmModal) Show()           { m.visible = true }
func (m *DeleteConfirmModal) Hide()           { m.visible = false }
func (m *DeleteConfirmModal) WasConfirmed() bool { return m.confirmed }
```

### Modal Overlay Rendering

Use `behaviors.RenderModalOverlay()` for overlay compositing:

```go
func (i *BrowseTimelineIntent) View() string {
    // Render base view from current screen
    var baseView string
    switch i.state {
    case StateList:
        baseView = i.listScreen.View()
    case StateDetail:
        baseView = i.detailScreen.View()
    }
    
    // Overlay modal if visible
    if i.filterModal != nil && i.filterModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.filterModal, baseView)
    }
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

---

## UIKit Foundation

### Package Overview

| Package | Components | Purpose |
|---------|------------|---------|
| `primitives/` | Text, Button, ButtonGroup, Input, Badge | Semantic styling |
| `containers/` | Box (4 variants), Overlay | Bordered containers |
| `layout/` | ScreenLayout, Header, Footer | Page structure |
| `feedback/` | Modal (5 types), OverlayModal | Status feedback |
| `navigation/` | BreadcrumbBar, Breadcrumb | Navigation context |
| `display/` | Logo | ASCII logo with animation |
| `theme/` | Theme interface, ThemeAware | Theme infrastructure |

### ScreenLayout

The standard way to create consistent layouts:

```go
func (s *MyScreen) View() string {
    termInfo := s.GetTerminalInfo()
    
    view := layout.NewScreenLayout(termInfo).
        WithLogo(s.GetLogo()).
        WithBreadcrumbs(s.breadcrumbs).
        WithContent(s.renderContent()).
        WithHelp(s.renderFooter()).
        WithFooterSeparator(true)
    
    // Optional: overlay modal
    if s.modal != nil && s.modal.IsVisible() {
        view.ShowModalOverlay(s.modal)
    }
    
    return view.Render()
}
```

### Theme Usage

Always handle nil themes:

```go
func (m *MyModal) View() string {
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    
    // Use UIKit primitives with theme
    title := primitives.Title("My Title", theme).Render()
    body := primitives.Body("Content here", theme).Render()
    
    return containers.NewBox(theme).
        Content(lipgloss.JoinVertical(lipgloss.Left, title, body)).
        Render()
}
```

---

## Dependency Flow

### Intent → Screen Communication

```
Intent                          Screen
   │                              │
   │  SetTerminalInfo(info)       │
   │─────────────────────────────▶│
   │                              │
   │  SetTheme(theme)             │
   │─────────────────────────────▶│
   │                              │
   │  SetLogo(logo)               │
   │─────────────────────────────▶│
   │                              │
   │  Update(msg)                 │
   │─────────────────────────────▶│
   │                              │
   │  (cmd, ScreenResult)         │
   │◀─────────────────────────────│
   │                              │
   │  View()                      │
   │─────────────────────────────▶│
   │                              │
   │  (string)                    │
   │◀─────────────────────────────│
```

### Intent → Modal Communication

```
Intent                          Modal
   │                              │
   │  modal.Show()                │
   │─────────────────────────────▶│
   │                              │
   │  modal.Update(msg)           │
   │─────────────────────────────▶│
   │                              │
   │  modal.IsVisible()           │
   │─────────────────────────────▶│
   │                              │
   │  modal.View()                │
   │─────────────────────────────▶│
   │                              │
   │  modal.WasConfirmed()        │  (domain-specific)
   │─────────────────────────────▶│
   │                              │
   │  modal.Hide()                │
   │─────────────────────────────▶│
```

### Screen → Behavior Communication

```
Screen                          TableBehavior
   │                                  │
   │  table.WithColumns(cols)         │
   │─────────────────────────────────▶│
   │                                  │
   │  table.WithItems(items)          │
   │─────────────────────────────────▶│
   │                                  │
   │  table.Update(msg)               │
   │─────────────────────────────────▶│
   │                                  │
   │  table.Selected()                │
   │─────────────────────────────────▶│
   │                                  │
   │  table.View()                    │
   │─────────────────────────────────▶│
```

---

## Creating New Intents

### Decision Tree: Direct Rendering vs Screens

```
Is this a simple intent with 1-2 states?
├── YES → Direct rendering (use BaseIntent helpers)
└── NO → Does it have reusable UI patterns?
         ├── YES → Create screens
         └── NO → Direct rendering with ScreenLayout
```

### Step-by-Step: Screen-Based Intent

#### 1. Create Screen Package

```bash
mkdir -p internal/cli/screens/myfeature
```

#### 2. Create Screens

```go
// internal/cli/screens/myfeature/list_screen.go
package myfeature

type ListScreen struct {
    *base.BaseScreen
    items []Item
    table *behaviors.TableBehavior[Item]
}

func NewListScreen(items []Item) *ListScreen {
    s := &ListScreen{
        BaseScreen: base.NewBaseScreen(),
        items:      items,
    }
    s.table = behaviors.NewTableBehavior[Item]().
        WithItems(items)
    return s
}
```

#### 3. Create Intent

```go
// internal/cli/intents/myfeature_intent.go
package intents

type MyFeatureIntent struct {
    *BaseIntent
    
    state      MyFeatureState
    active     bool
    result     *IntentResult[*MyFeatureResult]
    
    // Screens
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals
    filterModal *components.FilterModal
    
    // Data
    items []Item
}

func NewMyFeatureIntent(ctx context.Context) (*MyFeatureIntent, error) {
    intent := &MyFeatureIntent{
        BaseIntent: NewBaseIntent(),
        state:      StateList,
        active:     true,
    }
    return intent, nil
}
```

#### 4. Wire Screen to Intent

```go
func (i *MyFeatureIntent) Init() tea.Cmd {
    // Load data
    i.items = loadItems()
    
    // Create initial screen
    i.listScreen = myfeature.NewListScreen(i.items)
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.listScreen.SetLogo(i.GetLogo())
    
    return nil
}

func (i *MyFeatureIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Handle terminal size for all screens
    if wsm, ok := msg.(tea.WindowSizeMsg); ok {
        i.updateTerminalInfo(wsm)
        i.propagateTerminalInfo()
    }
    
    // Route to current screen
    switch i.state {
    case StateList:
        return i.updateList(msg)
    case StateDetail:
        return i.updateDetail(msg)
    }
    return nil
}

func (i *MyFeatureIntent) updateList(msg tea.Msg) tea.Cmd {
    // Handle modal first if visible
    if i.filterModal != nil && i.filterModal.IsVisible() {
        model, cmd := i.filterModal.Update(msg)
        i.filterModal = model.(*components.FilterModal)
        if !i.filterModal.IsVisible() {
            // Modal closed, apply filter
            i.applyFilter()
        }
        return cmd
    }
    
    // Delegate to screen
    cmd, result := i.listScreen.Update(msg)
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}

func (i *MyFeatureIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    switch r := result.(type) {
    case *screens.NavigateResult:
        return i.HandleNavigate(r)
    case *screens.CancelResult:
        return i.HandleCancel(r)
    case *screens.SubmitResult:
        return i.HandleSubmit(r)
    case *screens.ErrorResult:
        return i.HandleError(r)
    }
    return nil
}
```

#### 5. Implement View

```go
func (i *MyFeatureIntent) View() string {
    var baseView string
    
    switch i.state {
    case StateList:
        baseView = i.listScreen.View()
    case StateDetail:
        baseView = i.detailScreen.View()
    default:
        baseView = ""
    }
    
    // Overlay modals
    if i.filterModal != nil && i.filterModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.filterModal, baseView)
    }
    
    return baseView
}
```

---

## Migration Guide

### From Direct Rendering to Screens

#### Before (Direct Rendering)

```go
func (i *OldIntent) viewList() string {
    var content strings.Builder
    
    for idx, item := range i.items {
        if idx == i.selected {
            content.WriteString("> ")
        } else {
            content.WriteString("  ")
        }
        content.WriteString(item.Title)
        content.WriteString("\n")
    }
    
    return i.CreateViewWithBreadcrumbs(
        "My Feature",
        content.String(),
        "Enter: Select | Esc: Back",
    )
}
```

#### After (Screen-Based)

```go
// 1. Create screen
// internal/cli/screens/myfeature/list_screen.go
type ListScreen struct {
    *base.BaseScreen
    items    []Item
    selected int
}

func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            return nil, screens.NewCancelResult("")
        case "enter":
            return nil, screens.NewNavigateResult("detail", s.items[s.selected])
        case "up", "k":
            if s.selected > 0 {
                s.selected--
            }
        case "down", "j":
            if s.selected < len(s.items)-1 {
                s.selected++
            }
        }
    }
    return nil, nil
}

func (s *ListScreen) View() string {
    content := s.renderList()
    footer := "Enter: Select | Esc: Back"
    return s.CreateView("My Feature", content, footer)
}

// 2. Update intent
func (i *NewIntent) updateList(msg tea.Msg) tea.Cmd {
    cmd, result := i.listScreen.Update(msg)
    if result != nil {
        switch r := result.(type) {
        case *screens.NavigateResult:
            if r.Target == "detail" {
                i.state = StateDetail
                i.selectedItem = r.Data.(Item)
                i.detailScreen = myfeature.NewDetailScreen(i.selectedItem)
            }
        case *screens.CancelResult:
            i.setCancelled()
        }
    }
    return cmd
}

func (i *NewIntent) viewList() string {
    return i.listScreen.View()
}
```

### From Old Modals to UIKit Modals

#### Before (styles.* based)

```go
func (m *OldModal) View() string {
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder).
        Padding(1, 2)
    
    title := styles.CardHeader.Render("Title")
    body := styles.CardContent.Render("Content")
    
    return style.Render(title + "\n" + body)
}
```

#### After (UIKit based)

```go
func (m *NewModal) View() string {
    if !m.visible {
        return ""
    }
    
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    
    title := primitives.Title("Title", theme).Render()
    body := primitives.Body("Content", theme).Render()
    
    content := lipgloss.JoinVertical(lipgloss.Left, title, "", body)
    
    return containers.NewBox(theme).
        Content(content).
        Width(m.width).
        Padding(1).
        Background(theme.BackgroundColor()).
        Render()
}
```

---

## Best Practices

### DO

1. **Always embed BaseScreen/BaseIntent** for consistent behavior
2. **Handle nil themes** in View() methods
3. **Use ScreenLayout** for all page layouts
4. **Return ScreenResult** from screens, not raw state changes
5. **Use behaviors** for common UI patterns (tables, lists)
6. **Set solid backgrounds** on modals to prevent transparency
7. **Propagate terminal info** to all screens on WindowSizeMsg

### DON'T

1. **Don't import intents from screens** (breaks dependency hierarchy)
2. **Don't mutate screen state from intent** (use Update/Result pattern)
3. **Don't use raw lipgloss** for styling (use UIKit primitives)
4. **Don't stack modals** (close one before opening another)
5. **Don't skip WindowSizeMsg handling** (breaks responsiveness)
6. **Don't hardcode colors** (use theme colors)

### Code Quality Checklist

Before committing a new intent:

- [ ] Embeds `*BaseIntent`
- [ ] All screens embed `*base.BaseScreen`
- [ ] All modals use UIKit containers with solid backgrounds
- [ ] No circular dependencies
- [ ] Terminal info propagated to all screens
- [ ] Theme propagated to all screens and modals
- [ ] All screens return `ScreenResult` types
- [ ] Modal overlay uses `behaviors.RenderModalOverlay()`
- [ ] Tests cover all state transitions
- [ ] Escape key behavior tested for each state

---

## Reference Implementations

### Screen-Based Intent (Recommended)
- **BrowseTimelineIntent** - Full screen/modal integration
- **CaptureEventIntent** - Form-based screens
- **ConfigureSystemIntent** - Settings with modal sequence

### Direct Rendering Intent (Simple cases)
- **ExportArtifactIntent** - Simple state machine
- **MetadataEditorIntent** - Form-focused

### Modal Examples
- **DeleteConfirmModal** - Simple confirmation
- **FilterModalModel** - Form-based with huh
- **ViewEventDetailModal** - Read-only display

---

**See Also**:
- [TUI Developer Guide](./TUI_DEVELOPER_GUIDE.md) - Component development
- [UIKit Guide](./UIKIT_GUIDE.md) - Component library reference
- [Intent Development Checklist](./INTENT_DEVELOPMENT_CHECKLIST.md) - State machine patterns
- [Modal Patterns](./MODAL_PATTERNS.md) - Modal implementation guide
