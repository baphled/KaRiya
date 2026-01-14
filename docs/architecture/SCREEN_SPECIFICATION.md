---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya Screen Interface Specification

**Version**: 1.0
**Last Updated**: 2026-01-13
**Status**: Proposed

---

## Overview

Screens are reusable, composable view components that sit between Intents and Components. Each screen owns a single view within a workflow and communicates with its parent Intent via typed results.

---

## Core Interface

```go
package screens

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/baphled/kariya/internal/cli/terminal"
    "github.com/baphled/kariya/internal/cli/themes"
)

// Screen represents a single view within an intent workflow.
// Screens are reusable, composable, and testable in isolation.
type Screen interface {
    // Init initializes the screen (e.g., start spinner, load data)
    Init() tea.Cmd
    
    // Update handles input and returns:
    // - tea.Cmd for async operations
    // - ScreenResult when screen completes (nil if still active)
    Update(msg tea.Msg) (tea.Cmd, ScreenResult)
    
    // View renders the screen's current state
    View() string
    
    // SetTerminalInfo updates terminal dimensions
    SetTerminalInfo(info *terminal.Info)
    
    // SetTheme sets the active theme
    SetTheme(theme themes.Theme)
}
```

---

## Screen Results

```go
// ScreenResult communicates screen outcomes to the parent intent.
// Results are type-safe and exhaustive.
type ScreenResult interface {
    isScreenResult()  // marker method
}

// NavigateResult - user made a selection, ready to proceed
type NavigateResult struct {
    Data interface{}  // Selected value (type depends on screen)
}

func (r *NavigateResult) isScreenResult() {}

// CancelResult - user wants to go back
type CancelResult struct {
    Reason string  // "back", "escape", "cancel", "timeout"
}

func (r *CancelResult) isScreenResult() {}

// SubmitResult - form was submitted
type SubmitResult struct {
    Data   interface{}
    Errors []ValidationError
}

func (r *SubmitResult) isScreenResult() {}

// ErrorResult - an error occurred
type ErrorResult struct {
    Err error
}

func (r *ErrorResult) isScreenResult() {}

// ValidationError represents a field-level validation error
type ValidationError struct {
    Field   string
    Message string
}
```

---

## BaseScreen

All screens should embed `BaseScreen` for common functionality:

```go
// BaseScreen provides common screen functionality.
// Screens should embed this and call its methods.
type BaseScreen struct {
    // Terminal awareness
    terminalInfo *terminal.Info
    width        int
    height       int
    
    // Theme support
    theme themes.Theme
    
    // Modal overlay (for help, errors)
    activeModal ModalContent
    
    // Common view elements
    breadcrumbs []string
    help        string
    
    // Help modal
    helpVisible bool
}

// NewBaseScreen creates a new BaseScreen with defaults
func NewBaseScreen() *BaseScreen {
    return &BaseScreen{
        terminalInfo: terminal.NewInfo(),
        width:        80,
        height:       24,
    }
}

// SetTerminalInfo implements Screen interface
func (b *BaseScreen) SetTerminalInfo(info *terminal.Info) {
    b.terminalInfo = info
    if info != nil && info.IsValid {
        b.width = info.Width
        b.height = info.Height
    }
}

// SetTheme implements Screen interface
func (b *BaseScreen) SetTheme(theme themes.Theme) {
    b.theme = theme
}

// Theme returns the active theme
func (b *BaseScreen) Theme() themes.Theme {
    return b.theme
}

// SetBreadcrumbs sets navigation breadcrumbs
func (b *BaseScreen) SetBreadcrumbs(crumbs ...string) {
    b.breadcrumbs = crumbs
}

// SetHelp sets the help text
func (b *BaseScreen) SetHelp(help string) {
    b.help = help
}

// ShowModal displays a modal overlay
func (b *BaseScreen) ShowModal(modal ModalContent) {
    b.activeModal = modal
}

// HideModal hides the modal overlay
func (b *BaseScreen) HideModal() {
    b.activeModal = nil
}

// ToggleHelp toggles help modal visibility
func (b *BaseScreen) ToggleHelp() {
    b.helpVisible = !b.helpVisible
}

// CreateView returns a StandardView configured with common elements
func (b *BaseScreen) CreateView() *components.StandardView {
    view := components.NewStandardView(b.terminalInfo)
    view.WithBreadcrumbs(b.breadcrumbs...)
    view.WithHelp(b.help)
    view.SetUseFullWidth(true)
    
    if b.activeModal != nil {
        view.ShowModalOverlay(b.activeModal)
    }
    
    return view
}

// HandleGlobalKey handles global keys (?, q, esc)
// Returns ScreenResult if handled, nil if not
func (b *BaseScreen) HandleGlobalKey(msg tea.KeyMsg) ScreenResult {
    switch msg.String() {
    case "?":
        b.ToggleHelp()
        return nil  // Handled but no result
    case "q", "ctrl+c":
        return &CancelResult{Reason: "quit"}
    case "esc":
        if b.helpVisible {
            b.helpVisible = false
            return nil
        }
        return &CancelResult{Reason: "back"}
    }
    return nil
}
```

---

## Base Screen Types

### BaseSelectScreen[T]

For single-item selection from a list:

```go
// BaseSelectScreen provides single-selection list functionality
type BaseSelectScreen[T any] struct {
    *BaseScreen
    
    items        []T
    selected     int
    itemRenderer func(T, bool) string  // Render item, bool=selected
}

type SelectConfig[T any] struct {
    Items        []T
    Default      int
    ItemRenderer func(T, bool) string
    Breadcrumbs  []string
    Help         string
}

func NewSelectScreen[T any](cfg SelectConfig[T]) *BaseSelectScreen[T] {
    s := &BaseSelectScreen[T]{
        BaseScreen:   NewBaseScreen(),
        items:        cfg.Items,
        selected:     cfg.Default,
        itemRenderer: cfg.ItemRenderer,
    }
    s.SetBreadcrumbs(cfg.Breadcrumbs...)
    s.SetHelp(cfg.Help)
    if s.help == "" {
        s.SetHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back")
    }
    return s
}

func (s *BaseSelectScreen[T]) Init() tea.Cmd {
    return nil
}

func (s *BaseSelectScreen[T]) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        // Check global keys first
        if result := s.HandleGlobalKey(keyMsg); result != nil {
            return nil, result
        }
        
        switch keyMsg.String() {
        case "up", "k":
            if s.selected > 0 {
                s.selected--
            }
        case "down", "j":
            if s.selected < len(s.items)-1 {
                s.selected++
            }
        case "enter":
            return nil, &NavigateResult{Data: s.items[s.selected]}
        case "g", "home":
            s.selected = 0
        case "G", "end":
            s.selected = len(s.items) - 1
        }
    }
    return nil, nil
}

func (s *BaseSelectScreen[T]) View() string {
    view := s.CreateView()
    
    var content strings.Builder
    for i, item := range s.items {
        isSelected := i == s.selected
        content.WriteString(s.itemRenderer(item, isSelected))
        content.WriteString("\n")
    }
    
    view.WithContent(content.String())
    return view.Render()
}

// Selected returns the currently selected item
func (s *BaseSelectScreen[T]) Selected() T {
    return s.items[s.selected]
}

// SelectedIndex returns the current selection index
func (s *BaseSelectScreen[T]) SelectedIndex() int {
    return s.selected
}
```

---

### BaseListScreen[T]

For browsable lists with actions:

```go
// BaseListScreen provides browsable list with actions
type BaseListScreen[T any] struct {
    *BaseScreen
    
    items         []T
    selected      int
    itemRenderer  func(T, bool) string
    
    // Scrolling
    offset        int
    visibleRows   int
    
    // Actions
    actions       map[string]func(T) tea.Cmd
}

type ListConfig[T any] struct {
    Items        []T
    ItemRenderer func(T, bool) string
    VisibleRows  int
    Breadcrumbs  []string
    Help         string
}

func NewListScreen[T any](cfg ListConfig[T]) *BaseListScreen[T] {
    l := &BaseListScreen[T]{
        BaseScreen:   NewBaseScreen(),
        items:        cfg.Items,
        itemRenderer: cfg.ItemRenderer,
        visibleRows:  cfg.VisibleRows,
        actions:      make(map[string]func(T) tea.Cmd),
    }
    l.SetBreadcrumbs(cfg.Breadcrumbs...)
    l.SetHelp(cfg.Help)
    if l.help == "" {
        l.SetHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back")
    }
    return l
}

func (l *BaseListScreen[T]) RegisterAction(key string, action func(T) tea.Cmd) {
    l.actions[key] = action
}

func (l *BaseListScreen[T]) Init() tea.Cmd {
    return nil
}

func (l *BaseListScreen[T]) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if result := l.HandleGlobalKey(keyMsg); result != nil {
            return nil, result
        }
        
        switch keyMsg.String() {
        case "up", "k":
            if l.selected > 0 {
                l.selected--
                l.updateScroll()
            }
        case "down", "j":
            if l.selected < len(l.items)-1 {
                l.selected++
                l.updateScroll()
            }
        case "pgup":
            l.selected = max(0, l.selected-l.visibleRows)
            l.updateScroll()
        case "pgdown":
            l.selected = min(len(l.items)-1, l.selected+l.visibleRows)
            l.updateScroll()
        case "g", "home":
            l.selected = 0
            l.offset = 0
        case "G", "end":
            l.selected = len(l.items) - 1
            l.updateScroll()
        case "enter":
            return nil, &NavigateResult{Data: l.items[l.selected]}
        default:
            // Check for registered actions
            if action, ok := l.actions[keyMsg.String()]; ok {
                return action(l.items[l.selected]), nil
            }
        }
    }
    return nil, nil
}

func (l *BaseListScreen[T]) updateScroll() {
    if l.selected < l.offset {
        l.offset = l.selected
    } else if l.selected >= l.offset+l.visibleRows {
        l.offset = l.selected - l.visibleRows + 1
    }
}

func (l *BaseListScreen[T]) View() string {
    view := l.CreateView()
    
    var content strings.Builder
    end := min(len(l.items), l.offset+l.visibleRows)
    for i := l.offset; i < end; i++ {
        isSelected := i == l.selected
        content.WriteString(l.itemRenderer(l.items[i], isSelected))
        content.WriteString("\n")
    }
    
    view.WithContent(content.String())
    return view.Render()
}

// Selected returns the currently selected item
func (l *BaseListScreen[T]) Selected() T {
    return l.items[l.selected]
}

// SelectedIndex returns the current selection index
func (l *BaseListScreen[T]) SelectedIndex() int {
    return l.selected
}
```

---

### BaseFormScreen

For huh-based forms:

```go
// BaseFormScreen wraps huh forms with proper sizing
type BaseFormScreen struct {
    *BaseScreen
    
    form     *huh.Form
    formData interface{}
}

type FormConfig struct {
    Form        *huh.Form
    FormData    interface{}
    Breadcrumbs []string
    Help        string
}

func NewFormScreen(cfg FormConfig) *BaseFormScreen {
    f := &BaseFormScreen{
        BaseScreen: NewBaseScreen(),
        form:       cfg.Form,
        formData:   cfg.FormData,
    }
    f.SetBreadcrumbs(cfg.Breadcrumbs...)
    f.SetHelp(cfg.Help)
    if f.help == "" {
        f.SetHelp("Tab Next  Enter Submit  Esc Cancel")
    }
    return f
}

func (f *BaseFormScreen) Init() tea.Cmd {
    return f.form.Init()
}

func (f *BaseFormScreen) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    // Handle window resize
    if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
        f.SetTerminalInfo(&terminal.Info{
            Width:   wsMsg.Width,
            Height:  wsMsg.Height,
            IsValid: true,
        })
        f.form = f.form.WithHeight(wsMsg.Height - 10).WithWidth(wsMsg.Width - 4)
    }
    
    // Check for escape (form might consume it)
    if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
        return nil, &CancelResult{Reason: "cancel"}
    }
    
    // Delegate to form
    model, cmd := f.form.Update(msg)
    f.form = model.(*huh.Form)
    
    // Check if form completed
    if f.form.State == huh.StateCompleted {
        return nil, &SubmitResult{Data: f.formData}
    }
    
    return cmd, nil
}

func (f *BaseFormScreen) View() string {
    view := f.CreateView()
    view.WithContent(f.form.View())
    return view.Render()
}
```

---

### BaseDetailScreen

For read-only detail views:

```go
// BaseDetailScreen displays read-only content with optional actions
type BaseDetailScreen struct {
    *BaseScreen
    
    title   string
    content string
    actions map[string]string  // key -> label
}

type DetailConfig struct {
    Title       string
    Content     string
    Actions     map[string]string
    Breadcrumbs []string
    Help        string
}

func NewDetailScreen(cfg DetailConfig) *BaseDetailScreen {
    d := &BaseDetailScreen{
        BaseScreen: NewBaseScreen(),
        title:      cfg.Title,
        content:    cfg.Content,
        actions:    cfg.Actions,
    }
    d.SetBreadcrumbs(cfg.Breadcrumbs...)
    d.SetHelp(cfg.Help)
    if d.help == "" {
        d.SetHelp("Esc Back")
    }
    return d
}

func (d *BaseDetailScreen) Init() tea.Cmd {
    return nil
}

func (d *BaseDetailScreen) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if result := d.HandleGlobalKey(keyMsg); result != nil {
            return nil, result
        }
        
        // Check for action keys
        for key := range d.actions {
            if keyMsg.String() == key {
                return nil, &NavigateResult{Data: key}
            }
        }
    }
    return nil, nil
}

func (d *BaseDetailScreen) View() string {
    view := d.CreateView()
    
    var content strings.Builder
    if d.title != "" {
        content.WriteString(d.title)
        content.WriteString("\n\n")
    }
    content.WriteString(d.content)
    
    if len(d.actions) > 0 {
        content.WriteString("\n\n")
        content.WriteString("Actions: ")
        for key, label := range d.actions {
            content.WriteString(fmt.Sprintf("%s %s  ", key, label))
        }
    }
    
    view.WithContent(content.String())
    return view.Render()
}
```

---

### BaseConfirmScreen

For yes/no confirmations:

```go
// BaseConfirmScreen provides yes/no confirmation
type BaseConfirmScreen struct {
    *BaseScreen
    
    title    string
    message  string
    selected bool  // true=yes, false=no
}

type ConfirmConfig struct {
    Title       string
    Message     string
    Default     bool  // Default selection
    Breadcrumbs []string
}

func NewConfirmScreen(cfg ConfirmConfig) *BaseConfirmScreen {
    c := &BaseConfirmScreen{
        BaseScreen: NewBaseScreen(),
        title:      cfg.Title,
        message:    cfg.Message,
        selected:   cfg.Default,
    }
    c.SetBreadcrumbs(cfg.Breadcrumbs...)
    c.SetHelp("←/→ Select  Enter Confirm  Esc Cancel")
    return c
}

func (c *BaseConfirmScreen) Init() tea.Cmd {
    return nil
}

func (c *BaseConfirmScreen) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if result := c.HandleGlobalKey(keyMsg); result != nil {
            return nil, result
        }
        
        switch keyMsg.String() {
        case "left", "h":
            c.selected = true
        case "right", "l":
            c.selected = false
        case "y":
            return nil, &NavigateResult{Data: true}
        case "n":
            return nil, &CancelResult{Reason: "declined"}
        case "enter":
            if c.selected {
                return nil, &NavigateResult{Data: true}
            }
            return nil, &CancelResult{Reason: "declined"}
        }
    }
    return nil, nil
}

func (c *BaseConfirmScreen) View() string {
    view := c.CreateView()
    
    var content strings.Builder
    content.WriteString(c.title)
    content.WriteString("\n\n")
    content.WriteString(c.message)
    content.WriteString("\n\n")
    
    // Render buttons
    yesStyle := lipgloss.NewStyle().Padding(0, 2)
    noStyle := lipgloss.NewStyle().Padding(0, 2)
    
    if c.selected {
        yesStyle = yesStyle.Background(lipgloss.Color("62")).Foreground(lipgloss.Color("0"))
    } else {
        noStyle = noStyle.Background(lipgloss.Color("62")).Foreground(lipgloss.Color("0"))
    }
    
    yes := yesStyle.Render("Yes")
    no := noStyle.Render("No")
    content.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, yes, "  ", no))
    
    view.WithContent(content.String())
    return view.Render()
}
```

---

### BaseProgressScreen

For async operations with progress:

```go
// BaseProgressScreen shows progress for async operations
type BaseProgressScreen struct {
    *BaseScreen
    
    title     string
    message   string
    progress  float64  // 0.0 to 1.0
    spinner   spinner.Model
}

type ProgressConfig struct {
    Title       string
    Message     string
    Breadcrumbs []string
}

func NewProgressScreen(cfg ProgressConfig) *BaseProgressScreen {
    s := spinner.New()
    s.Spinner = spinner.Dot
    
    p := &BaseProgressScreen{
        BaseScreen: NewBaseScreen(),
        title:      cfg.Title,
        message:    cfg.Message,
        spinner:    s,
    }
    p.SetBreadcrumbs(cfg.Breadcrumbs...)
    p.SetHelp("Please wait...")
    return p
}

func (p *BaseProgressScreen) Init() tea.Cmd {
    return p.spinner.Tick
}

func (p *BaseProgressScreen) SetProgress(value float64, message string) {
    p.progress = value
    p.message = message
}

func (p *BaseProgressScreen) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    // Update spinner
    var cmd tea.Cmd
    p.spinner, cmd = p.spinner.Update(msg)
    return cmd, nil
}

func (p *BaseProgressScreen) View() string {
    view := p.CreateView()
    
    var content strings.Builder
    content.WriteString(p.title)
    content.WriteString("\n\n")
    content.WriteString(p.spinner.View())
    content.WriteString(" ")
    content.WriteString(p.message)
    
    if p.progress > 0 {
        content.WriteString("\n\n")
        content.WriteString(fmt.Sprintf("Progress: %.0f%%", p.progress*100))
    }
    
    view.WithContent(content.String())
    return view.Render()
}
```

---

## Domain-Specific Screen Example

Here's how to create a domain-specific screen using base screens:

```go
package skills

import (
    "github.com/baphled/kariya/internal/cli/screens/base"
    "github.com/baphled/kariya/internal/domain/career"
)

// SkillsList displays all skills grouped by category
type SkillsList struct {
    *base.BaseListScreen[*career.Skill]
    
    // Domain-specific state
    groupedByCategory map[string][]*career.Skill
    filters           *career.SkillFilters
}

type SkillsListConfig struct {
    Skills      []*career.Skill
    Filters     *career.SkillFilters
    Breadcrumbs []string
}

func NewSkillsList(cfg SkillsListConfig) *SkillsList {
    s := &SkillsList{
        BaseListScreen: base.NewListScreen(base.ListConfig[*career.Skill]{
            Items:        cfg.Skills,
            ItemRenderer: renderSkill,
            Breadcrumbs:  cfg.Breadcrumbs,
            Help:         "↑/k Up  ↓/j Down  n New  e Edit  d Delete  f Filter  Enter View",
        }),
        filters: cfg.Filters,
    }
    
    // Group skills by category
    s.groupedByCategory = groupByCategory(cfg.Skills)
    
    // Register domain-specific actions
    s.RegisterAction("n", func(_ *career.Skill) tea.Cmd { return nil }) // New
    s.RegisterAction("e", func(sk *career.Skill) tea.Cmd { return nil }) // Edit
    s.RegisterAction("d", func(sk *career.Skill) tea.Cmd { return nil }) // Delete
    
    return s
}

func renderSkill(skill *career.Skill, selected bool) string {
    marker := "  "
    if selected {
        marker = "▶ "
    }
    return fmt.Sprintf("%s%s (%s)", marker, skill.Name, skill.Category)
}

func groupByCategory(skills []*career.Skill) map[string][]*career.Skill {
    groups := make(map[string][]*career.Skill)
    for _, s := range skills {
        groups[s.Category] = append(groups[s.Category], s)
    }
    return groups
}
```

---

## Usage in Intent

```go
func (i *ManageSkillsIntent) transitionTo(state SkillsState) {
    i.state = state
    
    switch state {
    case SkillsList:
        screen := skills.NewSkillsList(skills.SkillsListConfig{
            Skills:      i.skills,
            Filters:     i.filters,
            Breadcrumbs: []string{"Main Menu", "Manage Skills"},
        })
        // Set terminal and theme from intent
        screen.SetTerminalInfo(i.terminalInfo)
        screen.SetTheme(i.theme)
        i.activeScreen = screen
        
    case SkillsDetail:
        // ... create detail screen
    }
}

func (i *ManageSkillsIntent) Update(msg tea.Msg) tea.Cmd {
    // Propagate terminal info to screen
    if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
        i.terminalInfo = &terminal.Info{
            Width:   wsMsg.Width,
            Height:  wsMsg.Height,
            IsValid: true,
        }
        i.activeScreen.SetTerminalInfo(i.terminalInfo)
    }
    
    // Delegate to screen
    cmd, result := i.activeScreen.Update(msg)
    
    // Handle screen result
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

func (i *ManageSkillsIntent) View() string {
    return i.activeScreen.View()
}
```

---

## Testing Screens

Screens can be tested independently:

```go
var _ = Describe("SkillsList", func() {
    var (
        screen *skills.SkillsList
        skills []*career.Skill
    )
    
    BeforeEach(func() {
        skills = []*career.Skill{
            {ID: "1", Name: "Go", Category: "backend"},
            {ID: "2", Name: "React", Category: "frontend"},
        }
        screen = skills.NewSkillsList(skills.SkillsListConfig{
            Skills: skills,
        })
    })
    
    Describe("Navigation", func() {
        It("moves down with j key", func() {
            _, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
            Expect(result).To(BeNil())
            Expect(screen.SelectedIndex()).To(Equal(1))
        })
        
        It("selects with enter key", func() {
            _, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
            Expect(result).To(BeAssignableToTypeOf(&NavigateResult{}))
            navResult := result.(*NavigateResult)
            Expect(navResult.Data).To(Equal(skills[0]))
        })
        
        It("cancels with escape key", func() {
            _, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
            Expect(result).To(BeAssignableToTypeOf(&CancelResult{}))
        })
    })
    
    Describe("View", func() {
        It("renders skill names", func() {
            view := screen.View()
            Expect(view).To(ContainSubstring("Go"))
            Expect(view).To(ContainSubstring("React"))
        })
        
        It("shows selection indicator", func() {
            view := screen.View()
            Expect(view).To(ContainSubstring("▶"))
        })
    })
})
```

---

## Best Practices

### 1. Screen Composition

Prefer composition over inheritance:

```go
// Good: Embed base screen, add domain logic
type SkillsList struct {
    *base.BaseListScreen[*career.Skill]
    groupedByCategory map[string][]*career.Skill
}

// Avoid: Duplicating base screen functionality
type SkillsList struct {
    items []career.Skill
    // ... reimplementing everything
}
```

### 2. Result Type Safety

Use specific result types for clarity:

```go
// Good: Clear what happened
case *screens.NavigateResult:
    skill := r.Data.(*career.Skill)
    
// Avoid: Generic data without context
case tea.Msg:
    // What type is this?
```

### 3. Terminal Awareness

Always handle window resize:

```go
func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, ScreenResult) {
    if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
        s.SetTerminalInfo(&terminal.Info{...})
    }
    // ...
}
```

### 4. Help Text

Provide clear, context-specific help:

```go
// Good: Specific to screen
"↑/k Up  ↓/j Down  n New  e Edit  d Delete  Enter View  Esc Back"

// Avoid: Generic or missing
"Use arrow keys"
```

---

## References

- **Architecture Doc**: `docs/architecture/TUI_ARCHITECTURE.md`
- **Naming Conventions**: `docs/rules/NAMING_CONVENTIONS.md`
- **Implementation Task**: `tasks/tasks-21-tui-architecture-refactor.md`
