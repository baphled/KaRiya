# Terminal UI Styling Reference

This document provides a comprehensive reference for styling terminal UIs in KaRiya using lipgloss and bubbles.

## Table of Contents

1. [Color Palette](#color-palette)
2. [Typography Styles](#typography-styles)
3. [Component Styles](#component-styles)
4. [Layout Patterns](#layout-patterns)
5. [Interactive Components](#interactive-components)
6. [Responsive Design](#responsive-design)
7. [Accessibility](#accessibility)
8. [Common Patterns](#common-patterns)

---

## Color Palette

### Primary Colors

All colors are defined in `internal/cli/styles/styles.go`:

```go
// Background colors
ColorBackground     = lipgloss.Color("#1a1f2e") // Dark blue-gray
ColorBackgroundAlt  = lipgloss.Color("#242936") // Slightly lighter
ColorBackgroundCard = lipgloss.Color("#2d3346") // Card/panel background

// Accent colors
ColorAccentTeal   = lipgloss.Color("#5fb3b3") // Primary action
ColorAccentGreen  = lipgloss.Color("#6cb56c") // Success/confirmation
ColorAccentPurple = lipgloss.Color("#a99bd1") // Selected items

// Text colors
ColorTextPrimary   = lipgloss.Color("#c7ccd1") // Primary text
ColorTextSecondary = lipgloss.Color("#8b92a0") // Secondary text
ColorTextMuted     = lipgloss.Color("#5e6673") // Muted text

// Status colors
ColorError   = lipgloss.Color("#d76e6e") // Error
ColorWarning = lipgloss.Color("#d9a66c") // Warning
ColorSuccess = lipgloss.Color("#6cb56c") // Success
ColorInfo    = lipgloss.Color("#6ab0d3") // Info
```

### Usage Guidelines

- **ColorAccentTeal**: Primary buttons, focused borders, important actions
- **ColorAccentGreen**: Success messages, confirmations, completed items
- **ColorError**: Validation errors, warnings, failed operations
- **ColorSuccess**: Success messages, completed tasks
- **ColorTextPrimary**: Main content, labels, headers
- **ColorTextSecondary**: Helper text, descriptions, meta information
- **ColorTextMuted**: Disabled items, placeholders, hints

---

## Typography Styles

### Pre-defined Styles

```go
// Header styles
CardHeader = lipgloss.NewStyle().
    Bold(true).
    Foreground(ColorTextPrimary)

// Content styles
CardContent = lipgloss.NewStyle().
    Foreground(ColorTextPrimary).
    Padding(1)

// Footer styles
CardFooter = lipgloss.NewStyle().
    Foreground(ColorTextSecondary).
    Italic(true)

// Button styles
ButtonPrimary = lipgloss.NewStyle().
    Background(ColorAccentTeal).
    Foreground(ColorBackground).
    Padding(0, 3).
    Bold(true)

ButtonSecondary = lipgloss.NewStyle().
    Background(ColorBackgroundAlt).
    Foreground(ColorAccentTeal).
    Padding(0, 3).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorAccentTeal)
```

### Creating Custom Styles

```go
// Bold text
boldStyle := lipgloss.NewStyle().Bold(true)

// Italic text
italicStyle := lipgloss.NewStyle().Italic(true)

// Underlined text
underlineStyle := lipgloss.NewStyle().Underline(true)

// Dimmed text
dimStyle := lipgloss.NewStyle().Faint(true)

// Combination
customStyle := lipgloss.NewStyle().
    Bold(true).
    Foreground(ColorTextPrimary).
    Background(ColorBackgroundAlt).
    Padding(1, 2)
```

---

## Component Styles

### Cards

```go
// Using CardContainer component
card := components.NewCardContainer().
    SetHeader("Card Title").
    SetBody("Card content here").
    SetFooter("Card footer").
    WithBorderColor(ColorBorderActive).
    Render()
```

### Forms

```go
// Using FormContainer component
form := components.NewFormContainer().
    AddField("Name", textInputModel).
    AddField("Email", emailInputModel).
    Render()
```

### Lists

```go
// Using ListContainer component
list := components.NewListContainer().
    SetItems(items).
    SetSelected(0).
    Render()
```

### Modal Dialogs

```go
// Using ModalContainer component
modal := components.NewModalContainer().
    SetTitle("Confirm Action").
    SetContent("Are you sure?").
    SetButtons("Yes", "No").
    Render()
```

---

## Layout Patterns

### Vertical Layout

```go
// Stack items vertically
layout := lipgloss.JoinVertical(
    lipgloss.Left,  // Alignment: Left, Center, Right
    item1,
    item2,
    item3,
)
```

### Horizontal Layout

```go
// Arrange items horizontally
layout := lipgloss.JoinHorizontal(
    lipgloss.Top,   // Alignment: Top, Center, Bottom
    item1,
    item2,
    item3,
)
```

### Centered Content

```go
// Center content in available space
centered := lipgloss.Place(
    width, height,
    lipgloss.Center, lipgloss.Center,
    content,
)
```

### Padding and Margins

```go
// Add padding (inside)
style := lipgloss.NewStyle().
    Padding(1, 2)  // Vertical, Horizontal

// Add margin (outside)
style := lipgloss.NewStyle().
    Margin(1, 2)   // Vertical, Horizontal

// Individual padding
style := lipgloss.NewStyle().
    PaddingTop(1).
    PaddingRight(2).
    PaddingBottom(1).
    PaddingLeft(2)
```

---

## Interactive Components

### Text Input

```go
import "github.com/charmbracelet/bubbles/textinput"

input := textinput.New()
input.Placeholder = "Enter text..."
input.Focus()
input.PromptStyle = lipgloss.NewStyle().Foreground(ColorAccentTeal)
input.TextStyle = lipgloss.NewStyle().Foreground(ColorTextPrimary)

// In Update()
input, cmd := input.Update(msg)

// In View()
return input.View()
```

### List Selection

```go
import "github.com/charmbracelet/bubbles/list"

items := []list.Item{
    item1, item2, item3,
}

l := list.New(items, list.NewDefaultDelegate(), width, height)
l.Title = "Select an item"
l.Styles.Title = lipgloss.NewStyle().
    Foreground(ColorTextPrimary).
    Bold(true)

// In Update()
l, cmd := l.Update(msg)

// In View()
return l.View()
```

### Spinner

```go
import "github.com/charmbracelet/bubbles/spinner"

s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(ColorAccentTeal)

// In Update()
s, cmd := s.Update(msg)

// In View()
return s.View() + " Loading..."
```

### Progress Bar

```go
import "github.com/charmbracelet/bubbles/progress"

p := progress.New(
    progress.WithDefaultGradient(),
    progress.WithoutPercentage(),
)

// In View()
return p.ViewAs(0.35)  // 35% progress
```

---

## Responsive Design

### Terminal Size Detection

```go
func (m Model) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return nil
}
```

### Adaptive Layouts

```go
func (m Model) View() string {
    // Full layout for wide terminals
    if m.width >= 120 {
        return m.renderWideLayout()
    }

    // Compact layout for narrow terminals
    if m.width < 80 {
        return m.renderCompactLayout()
    }

    // Standard layout
    return m.renderStandardLayout()
}
```

### Responsive Widths

```go
// Dynamic width based on terminal size
contentWidth := m.width - 4  // Account for padding

// Max width constraint
if contentWidth > 100 {
    contentWidth = 100
}

// Min width constraint
if contentWidth < 40 {
    contentWidth = 40
}

style := lipgloss.NewStyle().Width(contentWidth)
```

---

## Accessibility

### Keyboard Navigation

```go
// Clear focus indicators
func (m *Model) updateFocus() {
    m.input1.Blur()
    m.input2.Blur()

    if m.focused == 0 {
        m.input1.Focus()  // Visual feedback
    } else {
        m.input2.Focus()
    }
}

// Help text for keyboard shortcuts
helpText := "Tab: Next • Shift+Tab: Prev • Enter: Confirm • Esc: Cancel"
```

### Color Contrast

- Primary text on dark background: ✓ Good contrast
- Secondary text on dark background: ✓ Acceptable contrast
- Error text (red) on dark background: ✓ Good contrast
- Avoid color-only differentiation; use text labels

### Focus Indicators

```go
// Highlight focused element with border color
focusedStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorderActive)

// Use consistent focus styling across all components
```

---

## Common Patterns

### Form with Validation

```go
type FormModel struct {
    inputs []textinput.Model
    focused int
    errors map[int]string
}

func (m *FormModel) renderField(index int) string {
    input := m.inputs[index]

    // Highlight focused field
    borderColor := ColorBorder
    if m.focused == index {
        borderColor = ColorBorderActive
    }

    // Show error if present
    var errorText string
    if err, ok := m.errors[index]; ok {
        errorText = lipgloss.NewStyle().
            Foreground(ColorError).
            Render("✗ " + err)
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        input.View(),
        errorText,
    )
}
```

### Modal Overlay

```go
func (m Model) renderModal() string {
    modal := components.NewCardContainer().
        SetHeader("Confirm").
        SetBody("Are you sure?").
        SetFooter("Yes (y) / No (n)").
        Render()

    // Center modal on screen
    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        modal,
    )
}
```

### Status Bar

```go
func (m Model) renderStatusBar() string {
    left := lipgloss.NewStyle().
        Foreground(ColorTextPrimary).
        Render("Status: Ready")

    right := lipgloss.NewStyle().
        Foreground(ColorTextSecondary).
        Render("Ctrl+C to quit")

    gap := strings.Repeat(" ", m.width-len("Status: Ready")-len("Ctrl+C to quit"))

    return left + gap + right
}
```

### Loading Indicator

```go
func (m Model) renderLoading() string {
    spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"[m.spinner%10]

    return lipgloss.NewStyle().
        Foreground(ColorAccentTeal).
        Render(string(spinner) + " Loading...")
}
```

---

## Testing Styles

### Unit Tests

```go
func TestCardStyling(t *testing.T) {
    card := components.NewCardContainer().
        SetHeader("Test").
        SetBody("Content").
        Render()

    // Verify styling is applied
    assert.NotEmpty(t, card)
    assert.Contains(t, card, "Test")
    assert.Contains(t, card, "Content")
}
```

### Visual Tests

```go
// Test rendering at different terminal sizes
func TestResponsiveLayout(t *testing.T) {
    for _, width := range []int{40, 80, 120} {
        m := NewModel()
        m.width = width
        view := m.View()
        assert.NotEmpty(t, view)
    }
}
```

---

## Best Practices Summary

✅ **DO:**
- Use centralized style definitions
- Reuse component containers
- Test at different terminal widths
- Provide clear focus indicators
- Use consistent color meanings
- Document color/style choices
- Test keyboard navigation
- Support terminal resizing

❌ **DON'T:**
- Define styles inline
- Duplicate component rendering
- Assume fixed terminal size
- Use colors as only indicator
- Mix style definitions across files
- Ignore accessibility
- Hardcode dimensions
- Ignore responsive design

---

## Quick Reference

| Element | Style | Color | Notes |
|---------|-------|-------|-------|
| Header | Bold | Primary | Use CardHeader |
| Body | Normal | Primary | Use CardContent |
| Footer | Italic | Secondary | Use CardFooter |
| Primary Button | Bold | Teal BG | Use ButtonPrimary |
| Secondary Button | Normal | Teal FG | Use ButtonSecondary |
| Error | Normal | Error Red | Use ColorError |
| Success | Normal | Green | Use ColorSuccess |
| Focused Border | Normal | Teal | Use ColorBorderActive |
| Disabled | Faint | Muted | Use Faint(true) |

---

## Related Documentation

- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [KaRiya Styles Package](../internal/cli/styles/styles.go)
- [KaRiya Components](../internal/cli/components/)
- [Enhanced Example](../internal/cli/intents/enhanced_capture_example.go)

