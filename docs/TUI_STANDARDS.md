# KaRiya TUI Standards and Design Guidelines

**Document Version**: 1.0
**Last Updated**: 2025-12-30
**Status**: Complete and Ready for Implementation

---

## Table of Contents

1. [Overview](#overview)
2. [Keyboard Navigation](#keyboard-navigation)
3. [Component Architecture](#component-architecture)
4. [Visual Design](#visual-design)
5. [Implementation Patterns](#implementation-patterns)
6. [Best Practices](#best-practices)
7. [Accessibility](#accessibility)

---

## Overview

KaRiya's Terminal User Interface (TUI) follows Domain-Driven Design principles with consistent patterns for navigation, styling, and component composition. This document provides standards for creating and maintaining UI components across the application.

### Design Philosophy

- **Consistency**: Uniform keyboard shortcuts and patterns across all screens
- **Responsiveness**: Components adapt to various terminal sizes (40-300+ characters)
- **Simplicity**: Clear, uncluttered interfaces with focused functionality
- **Discoverability**: Help text and visual indicators guide users
- **Professionalism**: Dark theme with muted colors for extended use

---

## Keyboard Navigation

### Primary Navigation Keys

All navigation uses standardized keys defined in `internal/cli/navigation/constants.go`:

| Key | Function | Used In |
|-----|----------|---------|
| `Esc` | Back/Cancel | All screens |
| `↑/k` | Move up | Lists, menus, forms |
| `↓/j` | Move down | Lists, menus, forms |
| `←/h` | Move left | Forms, menus |
| `→/l` | Move right | Forms, menus |
| `Enter` | Select/Confirm | Forms, lists, menus |
| `Space` | Toggle | Checkboxes, selections |
| `Tab/Shift+Tab` | Next/Previous field | Forms |

### Contextual Shortcuts

| Key | Function | Context |
|-----|----------|---------|
| `c` | Capture event | Home screen |
| `l` | List events | Home screen |
| `m` | Metadata review | Home screen |
| `?` | Help | Any screen |
| `q` | Quit | Any screen |
| `f` | Filter | List view |
| `s` | Sort | List view |
| `/` | Search | List view |
| `e` | Edit | List/detail view |
| `d` | Delete | List/detail view |
| `b` | Bulk operations | List view |

### Vim-Style Navigation

- **j/k** for vertical navigation (move up/down)
- **h/l** for horizontal navigation (move left/right)
- **esc** for back navigation (consistent with vim convention)

Supports power users familiar with vim-style editing while remaining accessible to non-vim users through arrow keys.

### Design Rationale

- **Escape for back**: Standard across all screens, no duplication with text input
- **vim-style options**: Provides alternative navigation without replacing arrow keys
- **Tab navigation**: Standard form field navigation compatible with all terminals
- **Consistent modifiers**: Same key combinations work similarly across screens

---

## Component Architecture

### Component Hierarchy

```
BubbleTea Application (main.go)
├── App Model (internal/cli/app/app.go)
│   ├── FormModel (internal/cli/models/form.go)
│   │   ├── TextInput components
│   │   ├── TagSelector
│   │   └── CategorySelector
│   ├── ListModel (internal/cli/models/list.go)
│   │   ├── ListItems
│   │   └── Filters/Sort
│   ├── MetadataReviewModel (internal/cli/models/metadata_review.go)
│   ├── MetadataEditorModel (internal/cli/models/metadata_editor.go)
│   ├── BulkOperationsModel (internal/cli/models/bulk_operations.go)
│   └── [Other screen models...]
│
├── Reusable Components
│   ├── HelpFooter (internal/cli/components/help_footer.go)
│   ├── Header (internal/cli/components/header.go) [future]
│   ├── Footer (internal/cli/components/footer.go) [future]
│   ├── ListItem (internal/cli/components/list_item.go) [future]
│   ├── NavigationMenu (internal/cli/components/navigation_menu.go) [future]
│   ├── TagSelector (internal/cli/components/tag_selector.go)
│   └── CategorySelector (internal/cli/components/category_selector.go)
│
└── Navigation System
    ├── constants.go (keyboard definitions)
    └── help.go (help text generation)
```

### Screen Types

#### 1. **Input Screens** (Forms)
- **Purpose**: Capture user input
- **Components**: Text inputs, selectors, buttons
- **Navigation**: Tab/Shift+Tab between fields, Escape to cancel
- **Validation**: Real-time field validation with error display
- **Examples**: Event capture form, metadata editor

#### 2. **List Screens**
- **Purpose**: Display multiple items
- **Components**: List items, filters, sort controls
- **Navigation**: Up/Down/j/k to move, Enter to select, Escape to back
- **Selection**: Single selection with visual highlight
- **Examples**: Event list, bulk operations

#### 3. **Detail Screens**
- **Purpose**: Show full information for one item
- **Components**: Header, content, footer with actions
- **Navigation**: Escape to back, arrow keys for actions
- **Actions**: Edit, delete, related actions
- **Examples**: Event details, metadata view

#### 4. **Menu Screens**
- **Purpose**: Select from predefined options
- **Components**: Menu items with descriptions
- **Navigation**: Up/Down to move, Enter to select, Escape to cancel
- **Visual Hierarchy**: Focused item highlighted
- **Examples**: Help menu, action menus

### Component Interface

All screen models implement the BubbleTea Model interface:

```go
type Model interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

**Additional methods for consistent behavior:**

- `GetHeight() int` - For layout calculations
- `SetWidth(width int)` - For responsive sizing
- `IsCancelled() bool` - For dialog-like components
- `IsSubmitted() bool` - For form-like components

---

## Visual Design

### Color Scheme

**Professional Dark Theme:**

| Element | Color | Hex | Usage |
|---------|-------|-----|-------|
| Background | Dark blue-gray | #1a1f2e | Main background |
| Background Alt | Slightly lighter | #242936 | Secondary backgrounds |
| Card Background | Card color | #2d3346 | Cards/panels |
| Primary Text | Light gray | #c7ccd1 | Main text |
| Secondary Text | Medium gray | #8b92a0 | Help text |
| Muted Text | Dark gray | #5e6673 | Disabled text |
| Accent (Teal) | Muted teal | #5fb3b3 | Primary actions, borders |
| Accent (Green) | Muted green | #6cb56c | Success states |
| Accent (Purple) | Muted purple | #a99bd1 | Selected items |
| Error | Muted red | #d76e6e | Error states |
| Warning | Muted amber | #d9a66c | Warning states |
| Success | Muted green | #6cb56c | Success states |
| Info | Muted blue | #6ab0d3 | Information states |

### Styling Principles

1. **Minimal Color Use**: Limit color to 3-4 accent colors per screen
2. **Contrast**: Ensure text is readable on all backgrounds
3. **Consistency**: Use same colors for same meanings across screens
4. **Professional**: Muted tones for extended terminal use (less eye strain)
5. **Focus Indicators**: Clear visual feedback for focused elements

### Layout Guidelines

- **Maximum width**: 120 characters (for readability)
- **Minimum width**: 40 characters (mobile-friendly)
- **Padding**: 1-2 characters inside borders, 1 line above/below sections
- **Spacing**: Consistent 1-character gaps between elements
- **Alignment**: Left-aligned text for readability, centered headers
- **Responsive**: Scale content appropriately for terminal size

---

## Implementation Patterns

### Creating a New Screen

#### 1. Define the Model

```go
type MyScreenModel struct {
    width   int
    height  int
    state   string
    // Add your fields
}

// Constructor
func NewMyScreenModel(width int) MyScreenModel {
    return MyScreenModel{
        width: width,
        height: 20,
    }
}
```

#### 2. Implement Model Interface

```go
func (m MyScreenModel) Init() tea.Cmd {
    return nil
}

func (m MyScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m MyScreenModel) View() string {
    return m.render()
}

func (m MyScreenModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "esc":
        return m, func() tea.Msg { return BackMsg{} }
    case "q", "ctrl+c":
        return m, func() tea.Msg { return QuitMsg{} }
    }
    return m, nil
}
```

#### 3. Add Navigation and Help

```go
func (m MyScreenModel) render() string {
    // Render content
    content := "Your content here"

    // Add help footer
    helpFooter := components.NewHelpFooter("context_name", m.width)
    help := helpFooter.View()

    // Combine with proper layout
    return lipgloss.JoinVertical(
        lipgloss.Top,
        content,
        help,
    )
}
```

#### 4. Add Tests

```go
var _ = Describe("MyScreenModel", func() {
    Describe("Navigation", func() {
        It("should go back on Escape", func() {
            model := NewMyScreenModel(80)
            _, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
            Expect(cmd).NotTo(BeNil())
            msg := cmd()
            Expect(msg).To(BeAssignableToTypeOf(BackMsg{}))
        })
    })
})
```

### Creating a Reusable Component

#### 1. Define Component Interface

```go
type MyComponent struct {
    width   int
    data    interface{}
}

func NewMyComponent(width int) MyComponent {
    return MyComponent{width: width}
}
```

#### 2. Implement Rendering

```go
func (c MyComponent) Render() string {
    // Render component content
    return styles.Card.Render("Component content")
}

func (c MyComponent) RenderForWidth(width int) string {
    originalWidth := c.width
    c.width = width
    view := c.Render()
    c.width = originalWidth
    return view
}
```

#### 3. Add Styling Methods

```go
func (c MyComponent) WithBorder() MyComponent {
    c.style = c.style.Border(lipgloss.RoundedBorder())
    return c
}

func (c MyComponent) WithPadding(v, h int) MyComponent {
    c.style = c.style.Padding(v, h)
    return c
}
```

### Using Standardized Navigation

```go
import "github.com/baphled/kariya/internal/cli/navigation"

// Get help text for a context
helpText := navigation.GetContextualHelp("form")

// Generate help for specific keys
keys := []navigation.NavigationKey{
    navigation.KeyUp,
    navigation.KeyDown,
    navigation.KeySelect,
}
helpText := navigation.GetHelpTextCompact(keys, true)

// Add help footer
footer := components.NewHelpFooter("form", screenWidth)
footer.SetKeys(customKeys) // Optional: override context
renderedFooter := footer.View()
```

---

## Best Practices

### Navigation

1. **Always provide back navigation**
   - Every screen except Home should support Escape
   - Maintain previousScreen state for navigation history

2. **Consistent keyboard shortcuts**
   - Use standardized keys from navigation/constants.go
   - Avoid conflicts between text input and navigation

3. **Clear visual feedback**
   - Show focused element with highlight/border
   - Use help footer to show available shortcuts
   - Display validation errors inline with fields

### Performance

1. **Efficient rendering**
   - Only re-render when state changes
   - Use responsive width checks to avoid unnecessary recalculation
   - Cache computed values where appropriate

2. **Memory management**
   - Release resources in component cleanup
   - Avoid circular references in model state
   - Test with `go test -race` for concurrency issues

3. **Responsive design**
   - Check window size on resize messages
   - Gracefully handle narrow terminals (40+ chars)
   - Truncate long content with ellipsis

### Testing

1. **Unit tests for each model**
   - Test all keyboard shortcuts
   - Test state transitions
   - Test view rendering at various widths

2. **Integration tests**
   - Test navigation between screens
   - Test message passing between models
   - Test complete user workflows

3. **Property-based tests**
   - Test rendering on random terminal sizes
   - Test with random input data

### Code Quality

1. **Follow Go conventions**
   - Use clear, descriptive names
   - Keep functions small and focused
   - Write comments for exported functions

2. **Use existing components**
   - Don't duplicate functionality
   - Extend components instead of copying
   - Share styling through styles package

3. **Maintain consistency**
   - Match existing code style
   - Use same message types for similar actions
   - Follow established patterns for common operations

---

## Accessibility

### Keyboard Navigation

- ✅ All functions accessible via keyboard (no mouse required)
- ✅ Clear keyboard shortcuts shown in help text
- ✅ Navigation options clearly labeled
- ✅ Tab order logical and predictable

### Visual Design

- ✅ Sufficient color contrast (WCAG AA standard)
- ✅ No critical information conveyed by color alone
- ✅ Clear focus indicators for keyboard navigation
- ✅ Text-based content, not image-dependent

### Content

- ✅ Help text describes all available actions
- ✅ Error messages are specific and actionable
- ✅ Labels for all input fields
- ✅ Status messages provide context

---

## Future Enhancements

### Planned Components

- [ ] Navigation menu component (3.1-3.5)
- [ ] Header component with breadcrumbs (4.1-4.3)
- [ ] Footer component for status (4.4-4.6)
- [ ] List item component (5.1-5.3)
- [ ] Pagination support
- [ ] Search/filter indicators

### Planned Features

- [ ] Mouse support for clicking
- [ ] Theme customization
- [ ] Custom key binding
- [ ] Screen recording/playback
- [ ] Accessibility profiles
- [ ] Performance profiling

---

## References

- **Navigation System**: `internal/cli/navigation/`
- **Components**: `internal/cli/components/`
- **Styles**: `internal/cli/styles/`
- **Models**: `internal/cli/models/`
- **BubbleTea Docs**: https://github.com/charmbracelet/bubbletea
- **Lipgloss Docs**: https://github.com/charmbracelet/lipgloss

---

**Status**: ✅ Complete and ready for implementation
**Last Review**: 2025-12-30
**Next Review**: Upon completion of Phase 2

