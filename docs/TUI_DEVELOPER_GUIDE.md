# KaRiya TUI Developer Guide

**A comprehensive guide for developers creating and maintaining TUI components**

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Architecture Overview](#architecture-overview)
3. [Creating Components](#creating-components)
4. [Styling and Layout](#styling-and-layout)
5. [Testing Components](#testing-components)
6. [Debugging TUI Issues](#debugging-tui-issues)
7. [Performance Optimization](#performance-optimization)
8. [Advanced Patterns](#advanced-patterns)
9. [Common Pitfalls](#common-pitfalls)

---

## Getting Started

### Prerequisites

- Go 1.24.0 or later
- Understanding of BubbleTea framework
- Familiarity with Lipgloss for styling
- Basic understanding of TUI concepts

### Key Dependencies

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"

    "github.com/baphled/kariya/internal/cli/navigation"
    "github.com/baphled/kariya/internal/cli/styles"
    "github.com/baphled/kariya/internal/cli/components"
)
```

### Recommended Reading

1. [TUI Standards](./TUI_STANDARDS.md) - Design principles and guidelines
2. [Keyboard Reference](./KEYBOARD_REFERENCE.md) - All keyboard shortcuts
3. [Theme Customization Guide](./THEME_CUSTOMIZATION_GUIDE.md) - Theme system documentation
4. [BubbleTea Docs](https://github.com/charmbracelet/bubbletea/tree/master/examples) - Framework examples
5. [Lipgloss Styling](https://github.com/charmbracelet/lipgloss/examples) - Styling examples

---

## Architecture Overview

### Directory Structure

```
internal/cli/
├── navigation/              # Keyboard shortcuts system
│   ├── constants.go        # Standardized keys
│   └── help.go            # Help text generation
├── components/            # Reusable UI components
│   ├── help_footer.go     # Help text footer
│   ├── tag_selector.go    # Tag selection
│   └── [future components]
├── models/                # Screen models
│   ├── form.go           # Event capture form
│   ├── list.go           # Event list view
│   ├── metadata_editor.go # Metadata editor
│   └── [other screens]
├── styles/                # Visual styling
│   └── styles.go         # Color and layout styles
├── app/                   # Application main model
│   └── app.go            # Screen navigation
└── service/               # Business logic adapter
    └── event_service.go   # Service wrapper
```

### Message Flow

```
BubbleTea Input
    ↓
App.Update() - routes to current screen
    ↓
Screen Model.Update() - handles input
    ↓
Returns Model, Cmd
    ↓
Cmd() generates Message
    ↓
App processes Message (like BackMsg, QuitMsg)
    ↓
App.View() renders screen
```

---

## Creating Components

### Step 1: Define the Component Structure

```go
package components

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// MyComponentModel represents your component
type MyComponentModel struct {
    width       int
    height      int
    state       string
    data        interface{}
    // Add your fields
}

// NewMyComponent creates a new component instance
func NewMyComponent(width int) MyComponentModel {
    return MyComponentModel{
        width:  width,
        height: 10, // Default height
        state:  "idle",
    }
}
```

### Step 2: Implement BubbleTea Interface

```go
// Init is called when the component is first created
func (m MyComponentModel) Init() tea.Cmd {
    // Return any initial commands (API calls, file reads, etc.)
    return nil
}

// Update handles incoming messages
func (m MyComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    // Handle custom messages
    case MyCustomMsg:
        m.state = "custom"
        return m, nil
    }

    return m, nil
}

// View renders the component to a string
func (m MyComponentModel) View() string {
    return m.render()
}
```

### Step 3: Implement Key Handling

```go
func (m MyComponentModel) handleKeyPress(
    msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "esc":
        // Always support Escape for back/cancel
        return m, func() tea.Msg {
            return BackMsg{}
        }

    case "q", "ctrl+c":
        // Support quit shortcuts
        return m, func() tea.Msg {
            return QuitMsg{}
        }

    case "up", "k":
        // Vim-style alternative
        m.moveUp()
        return m, nil

    case "down", "j":
        m.moveDown()
        return m, nil

    case "enter":
        // Confirm/select action
        return m, m.handleSelect()
    }

    return m, nil
}
```

### Step 4: Implement Rendering

```go
func (m MyComponentModel) render() string {
    // Create content
    content := m.renderContent()

    // Add help footer
    helpFooter := components.NewHelpFooter("context_name", m.width)
    help := helpFooter.View()

    // Combine with layout
    view := lipgloss.JoinVertical(
        lipgloss.Top,
        content,
        help,
    )

    // Handle responsive sizing
    if m.width < 40 {
        return m.renderCompact()
    }

    return view
}

func (m MyComponentModel) renderContent() string {
    // Use styles from styles package
    return styles.CardContent.Render("Your content")
}

func (m MyComponentModel) renderCompact() string {
    // Simplified layout for narrow terminals
    return "Compact view for " + m.state
}
```

### Step 5: Add Helper Methods

```go
// GetHeight returns the component height
func (m MyComponentModel) GetHeight() int {
    return m.height
}

// SetWidth updates component width (for responsive design)
func (m *MyComponentModel) SetWidth(width int) {
    m.width = width
}

// GetState returns current state (useful for parent models)
func (m MyComponentModel) GetState() string {
    return m.state
}

// RenderForWidth renders at specific width without changing state
func (m MyComponentModel) RenderForWidth(width int) string {
    original := m.width
    m.width = width
    view := m.View()
    m.width = original
    return view
}

// WithBorder adds a border to rendered output
func (m MyComponentModel) WithBorder() string {
    return styles.BorderStyle.Render(m.View())
}

// WithPadding adds padding to rendered output
func (m MyComponentModel) WithPadding(v, h int) string {
    return lipgloss.NewStyle().
        Padding(v, h).
        Render(m.View())
}
```

---

## Styling and Layout

### Using the Styles Package

```go
import "github.com/baphled/kariya/internal/cli/styles"

// Apply predefined styles
styledText := styles.CardHeader.Render("Header Text")
styledContent := styles.CardContent.Render("Content")
styledError := styles.ErrorBox.Render("Error message")

// Combine styles
highlighted := styles.ButtonPrimary.Render("Click me")
focused := styles.ButtonFocused.Render("Focused button")
disabled := styles.ButtonDisabled.Render("Disabled")
```

### Creating Custom Styles

```go
// Define custom style
myCustomStyle := lipgloss.NewStyle().
    Foreground(styles.ColorTextPrimary).
    Background(styles.ColorBackgroundCard).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(styles.ColorAccentTeal).
    Padding(1, 2).
    MarginTop(1).
    MarginBottom(1)

// Use custom style
styledContent := myCustomStyle.Render("Content here")
```

### Responsive Layout

```go
// Check terminal width and adapt
func (m MyComponentModel) render() string {
    switch {
    case m.width < 40:
        return m.renderVeryNarrow()
    case m.width < 80:
        return m.renderNarrow()
    case m.width < 120:
        return m.renderWide()
    default:
        return m.renderVeryWide()
    }
}

// Use MaxWidth helper for responsive centering
content := lipgloss.NewStyle().
    MaxWidth(120).
    Render(m.content)
```

### Layout Patterns

**Vertical Stack:**
```go
view := lipgloss.JoinVertical(
    lipgloss.Top,
    header,
    content,
    footer,
)
```

**Horizontal Layout:**
```go
view := lipgloss.JoinHorizontal(
    lipgloss.Left,
    sidebar,
    mainContent,
)
```

**Grid Layout:**
```go
view := lipgloss.NewStyle().
    Width(m.width).
    Height(m.height).
    Render(content)
```

### Theme System Integration

KaRiya uses a theme system for consistent styling across all UI components. See [Theme Customization Guide](THEME_CUSTOMIZATION_GUIDE.md) for full details.

**Quick Start:**

```go
import "github.com/baphled/kariya/internal/cli/themes"

// Access theme in an intent (via BaseIntent)
func (i *MyIntent) View() string {
    theme := i.Theme()
    if theme == nil {
        return i.renderWithDefaults()
    }
    
    // Use pre-composed styles
    card := theme.Styles().CardBase.Render("Content")
    return card
}
```

**Helper Method Pattern (Recommended):**

```go
// Create helpers for commonly used colors with fallbacks
func (i *MyIntent) getCardStyle() lipgloss.Style {
    if theme := i.Theme(); theme != nil {
        return theme.Styles().CardBase
    }
    return lipgloss.NewStyle().
        Padding(1, 2).
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder)
}

func (i *MyIntent) getAccentColor() lipgloss.Color {
    if theme := i.Theme(); theme != nil {
        return theme.PrimaryColor()
    }
    return styles.ColorAccentTeal
}
```

**Themed Bubbles Components:**

```go
// Apply theme to bubbles/table
if theme := i.Theme(); theme != nil {
    i.table.SetStyles(themes.NewThemedTableStyles(theme))
}

// Apply theme to bubbles/list
listStyles := themes.NewThemedListStyles(theme)
myList.Styles = listStyles
```

---

## Testing Components

### Unit Test Structure

```go
package components

import (
    "testing"

    tea "github.com/charmbracelet/bubbletea"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestMyComponent(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "MyComponent Suite")
}

var _ = Describe("MyComponentModel", func() {
    var component MyComponentModel

    BeforeEach(func() {
        component = NewMyComponent(80)
    })

    Describe("Initialization", func() {
        It("should initialize with correct width", func() {
            Expect(component.width).To(Equal(80))
        })

        It("should start in idle state", func() {
            Expect(component.state).To(Equal("idle"))
        })
    })

    Describe("Keyboard Navigation", func() {
        It("should handle Escape key", func() {
            msg := tea.KeyMsg{Type: tea.KeyEsc}
            _, cmd := component.Update(msg)
            Expect(cmd).NotTo(BeNil())
        })

        It("should support vim navigation", func() {
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
            updated, _ := component.Update(msg)
            Expect(updated).NotTo(BeNil())
        })
    })

    Describe("Rendering", func() {
        It("should render content", func() {
            view := component.View()
            Expect(view).NotTo(BeEmpty())
        })

        It("should handle various widths", func() {
            for width := 40; width <= 200; width += 20 {
                component.SetWidth(width)
                view := component.View()
                Expect(view).NotTo(BeEmpty())
            }
        })
    })
})
```

### Testing Best Practices

1. **Test keyboard shortcuts**
   ```go
   It("should navigate with arrow keys", func() {
       msg := tea.KeyMsg{Type: tea.KeyUp}
       _, _ = component.Update(msg)
       Expect(component.position).To(Equal(0))
   })
   ```

2. **Test responsive rendering**
   ```go
   It("should adapt to narrow terminals", func() {
       component.SetWidth(40)
       view := component.View()
       Expect(len(view)).To(BeLessThan(100))
   })
   ```

3. **Test state transitions**
   ```go
   It("should transition to selected state", func() {
       component.Update(tea.KeyMsg{Type: tea.KeyEnter})
       Expect(component.state).To(Equal("selected"))
   })
   ```

4. **Test message handling**
   ```go
   It("should handle custom messages", func() {
       msg := MyCustomMsg{data: "test"}
       _, _ = component.Update(msg)
       Expect(component.data).To(Equal("test"))
   })
   ```

---

## Debugging TUI Issues

### Common Issues and Solutions

**Problem: Keys not responding**
```go
// Debug: Log all key messages
case tea.KeyMsg:
    logger.Info("Key pressed: %s", msg.String())
    return m.handleKeyPress(msg)
```

**Problem: Content not rendering**
```go
// Debug: Check if View() returns empty string
func (m MyComponent) View() string {
    content := m.render()
    if content == "" {
        return "DEBUG: Empty content"
    }
    return content
}
```

**Problem: Terminal size issues**
```go
// Debug: Check width/height
func (m MyComponent) View() string {
    header := fmt.Sprintf("Width: %d, Height: %d\n", m.width, m.height)
    return header + m.render()
}
```

### Debugging Tools

**Logger usage:**
```go
import "github.com/baphled/kariya/internal/logger"

logger := logger.New(os.Stdout, logger.DebugLevel)
logger.WithFields(map[string]string{
    "component": "MyComponent",
    "width": fmt.Sprint(m.width),
}).Debug("Component rendered")
```

**Race condition detection:**
```bash
# Run tests with race detector
go test -race ./internal/cli/components/
```

**Test coverage:**
```bash
# Generate coverage report
go test -cover ./internal/cli/components/
go tool cover -html=coverage.out
```

---

## Performance Optimization

### Rendering Optimization

1. **Cache computed values**
   ```go
   type MyComponent struct {
       cachedView string
       dirty      bool
   }

   func (m *MyComponent) View() string {
       if !m.dirty {
           return m.cachedView
       }
       m.cachedView = m.render()
       m.dirty = false
       return m.cachedView
   }
   ```

2. **Minimize string operations**
   ```go
   // Bad: Creates many intermediate strings
   s := "a" + "b" + "c" + "d"

   // Good: Use builder or join
   parts := []string{"a", "b", "c", "d"}
   s := strings.Join(parts, "")
   ```

3. **Reuse expensive computations**
   ```go
   // Cache color parsing
   var colorCache = make(map[string]lipgloss.Color)

   func getColor(name string) lipgloss.Color {
       if color, ok := colorCache[name]; ok {
           return color
       }
       color := lipgloss.Color(name)
       colorCache[name] = color
       return color
   }
   ```

### Memory Optimization

1. **Avoid large allocations in Update()**
   ```go
   // Bad: Creates new slice each update
   func (m *Component) Update(...) {
       m.items = make([]string, 1000)
   }

   // Good: Reuse or pre-allocate
   func (m *Component) Init() {
       m.items = make([]string, 0, 1000)
   }
   ```

2. **Minimal message sizes**
   ```go
   // Bad: Large message with full object
   type UpdateMsg struct {
       Event *CareerEvent
   }

   // Good: Minimal message with just ID
   type UpdateMsg struct {
       EventID string
   }
   ```

---

## Advanced Patterns

### Custom Message Types

```go
// Define custom messages for your component
type MyComponentMsg struct {
    data string
}

type SelectMsg struct {
    index int
}

// Use in Update()
func (m MyComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case MyComponentMsg:
        m.data = msg.data
        return m, nil

    case SelectMsg:
        m.selectedIndex = msg.index
        return m, nil
    }
    return m, nil
}
```

### Delegating to Sub-models

```go
// Component can include sub-models
type FormComponent struct {
    textInput textinput.Model
    tagSel    TagSelector
}

// Delegate updates to sub-models
func (m *FormComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.textInput, cmd = m.textInput.Update(msg)
    if cmd != nil {
        return m, cmd
    }

    // Update other sub-components
    return m, nil
}

// Combine renders
func (m FormComponent) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Top,
        m.textInput.View(),
        m.tagSel.View(),
    )
}
```

### Concurrent Updates

```go
// Use channels to update from other goroutines
type AsyncComponent struct {
    updates chan string
}

func (m AsyncComponent) Init() tea.Cmd {
    return func() tea.Msg {
        return <-m.updates
    }
}

func (m AsyncComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case string:
        m.content = msg
        // Fetch next update
        return m, func() tea.Msg {
            return <-m.updates
        }
    }
    return m, nil
}
```

---

## Common Pitfalls

### 1. Forgetting to Handle Escape

```go
// ❌ Bad: No Escape handler
switch msg.String() {
case "enter":
    return m, m.handleSelect()
}

// ✅ Good: Always handle Escape
switch msg.String() {
case "esc":
    return m, func() tea.Msg { return BackMsg{} }
case "enter":
    return m, m.handleSelect()
}
```

### 2. Ignoring Window Size Changes

```go
// ❌ Bad: Doesn't handle resize
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Ignores tea.WindowSizeMsg
    return m, nil
}

// ✅ Good: Responds to resize
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    return m, nil
}
```

### 3. Not Testing Edge Cases

```go
// ✅ Good: Test edge cases
It("should handle very narrow terminal", func() {
    m.width = 20
    view := m.View()
    Expect(view).NotTo(BeEmpty())
})

It("should handle very wide terminal", func() {
    m.width = 500
    view := m.View()
    Expect(view).NotTo(BeEmpty())
})
```

### 4. Hardcoding Colors

```go
// ❌ Bad: Hardcoded colors
style := lipgloss.NewStyle().Foreground("#FF0000")

// ✅ Good: Use centralized styles
style := styles.ErrorBox
```

### 5. Not Providing Help Text

```go
// ❌ Bad: No help for users
func (m Model) View() string {
    return m.renderContent()
}

// ✅ Good: Include help footer
func (m Model) View() string {
    footer := components.NewHelpFooter("form", m.width)
    return lipgloss.JoinVertical(
        lipgloss.Top,
        m.renderContent(),
        footer.View(),
    )
}
```

---

## Summary Checklist

When creating a new TUI component, ensure:

- ✅ Implements BubbleTea Model interface
- ✅ Handles `Esc` for back/cancel
- ✅ Handles window resize messages
- ✅ Handles navigation keys (↑↓ or j/k)
- ✅ Uses styles from `styles` package
- ✅ Includes help footer or help text
- ✅ Responsive to terminal width changes
- ✅ Has comprehensive unit tests
- ✅ Tests keyboard shortcuts
- ✅ Tests edge cases (narrow/wide terminals)
- ✅ No memory leaks or race conditions
- ✅ Clear, descriptive error messages

---

**Last Updated**: 2025-12-30
**Version**: 1.0
**Status**: Complete and ready for reference

For more information, see [TUI Standards](./TUI_STANDARDS.md) and [Keyboard Reference](./KEYBOARD_REFERENCE.md).

