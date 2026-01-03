# Lipgloss & Bubbles Terminal UI Guide

This guide demonstrates best practices for building beautiful, interactive terminal UIs in Go using **lipgloss** for styling and **bubbles** for interactive components.

## Overview

- **lipgloss**: A library for styling and formatting text in terminal applications
- **bubbles**: A collection of ready-to-use interactive components (textinput, list, table, etc.)

## Architecture Principles

1. **Separation of Concerns**: Styling logic separate from component logic
2. **Reusable Styles**: Define styles once, apply everywhere
3. **Component Composition**: Build complex UIs from simple components
4. **Responsive Design**: Adapt to terminal width/height
5. **Consistent Theme**: Use a unified color scheme throughout

## Project Structure

```
internal/cli/
├── styles/
│   ├── styles.go           # Central style definitions
│   └── constants_export.go # Exported color constants
├── components/
│   ├── card_container.go   # Reusable card component
│   ├── form_container.go   # Reusable form component
│   └── ...other components
└── intents/
    ├── capture_event.go    # Intent implementation
    └── modals.go           # Modal sub-flows
```

## Key Concepts

### 1. Centralized Styling (styles/styles.go)

Define all styles in one place for consistency:

```go
var (
    // Color palette
    ColorBackground = lipgloss.Color("#1a1f2e")
    ColorAccentTeal = lipgloss.Color("#5fb3b3")

    // Reusable styles
    ButtonPrimary = lipgloss.NewStyle().
        Background(ColorAccentTeal).
        Foreground(ColorTextPrimary).
        Padding(0, 3)

    CardHeader = lipgloss.NewStyle().
        Bold(true).
        Foreground(ColorTextPrimary)
)
```

**Benefits:**
- Single source of truth for styling
- Easy to maintain and update themes
- Consistent across entire application
- Simple to support light/dark modes

### 2. Reusable Container Components

Build containers that wrap common UI patterns:

```go
type CardContainer struct {
    header string
    body   string
    footer string
}

func (cc *CardContainer) Render() string {
    // Combine header, body, footer with consistent styling
}
```

**Benefits:**
- Eliminates code duplication
- Enforces consistent spacing and alignment
- Easy to modify styling globally
- Reduces cognitive load when building UIs

### 3. Bubbles Components for Interactivity

Use bubbles components for user input:

```go
import "github.com/charmbracelet/bubbles/textinput"

type Model struct {
    input textinput.Model  // Bubbles component
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle key input
        var cmd tea.Cmd
        m.input, cmd = m.input.Update(msg)
        return m, cmd
    }
}
```

**Benefits:**
- Pre-built, well-tested interactive components
- Consistent keyboard handling
- Proper focus management
- Built-in validation support

### 4. Responsive Layout

Adapt UI to terminal dimensions:

```go
func (m Model) View() string {
    width := m.width
    height := m.height

    // Adjust layout based on available space
    if width < 80 {
        return m.renderCompactLayout()
    }
    return m.renderFullLayout()
}
```

## Complete Example: Enhanced Capture Event Intent

Below is a complete example showing how to build an interactive form using both lipgloss and bubbles:

```go
package intents

import (
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/baphled/kariya/internal/cli/styles"
    "github.com/baphled/kariya/internal/cli/components"
)

type CaptureEventModel struct {
    // State
    state State

    // Interactive components (from bubbles)
    titleInput       textinput.Model
    descriptionInput textinput.Model
    companyInput     textinput.Model

    // UI state
    focusedField int
    width        int
    height       int

    // Data
    data *EventData
}

type State string

const (
    StateForm    State = "form"
    StateReview  State = "review"
    StateSubmit  State = "submit"
)

type EventData struct {
    Title       string
    Description string
    Company     string
    Tags        []string
}

// NewCaptureEventModel creates a new capture event model
func NewCaptureEventModel() *CaptureEventModel {
    m := &CaptureEventModel{
        state:        StateForm,
        focusedField: 0,
        data:         &EventData{},
    }

    // Initialize bubbles components
    m.titleInput = textinput.New()
    m.titleInput.Placeholder = "Enter event title..."
    m.titleInput.Focus()

    m.descriptionInput = textinput.New()
    m.descriptionInput.Placeholder = "Enter description..."

    m.companyInput = textinput.New()
    m.companyInput.Placeholder = "Enter company..."

    return m
}

// Init initializes the model
func (m *CaptureEventModel) Init() tea.Cmd {
    return textinput.Blink
}

// Update handles messages
func (m *CaptureEventModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit

        case "tab":
            m.focusedField = (m.focusedField + 1) % 3
            m.updateFocus()

        case "shift+tab":
            m.focusedField = (m.focusedField - 1 + 3) % 3
            m.updateFocus()

        case "enter":
            if m.state == StateForm {
                m.syncData()
                m.state = StateReview
            }
        }

        // Delegate to focused input
        return m.updateInput(msg)
    }

    return m, nil
}

// updateInput delegates keyboard input to the focused field
func (m *CaptureEventModel) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch m.focusedField {
    case 0:
        m.titleInput, cmd = m.titleInput.Update(msg)
    case 1:
        m.descriptionInput, cmd = m.descriptionInput.Update(msg)
    case 2:
        m.companyInput, cmd = m.companyInput.Update(msg)
    }

    return m, cmd
}

// updateFocus updates which field is focused
func (m *CaptureEventModel) updateFocus() {
    m.titleInput.Blur()
    m.descriptionInput.Blur()
    m.companyInput.Blur()

    switch m.focusedField {
    case 0:
        m.titleInput.Focus()
    case 1:
        m.descriptionInput.Focus()
    case 2:
        m.companyInput.Focus()
    }
}

// syncData copies input values to data struct
func (m *CaptureEventModel) syncData() {
    m.data.Title = m.titleInput.Value()
    m.data.Description = m.descriptionInput.Value()
    m.data.Company = m.companyInput.Value()
}

// View renders the model
func (m *CaptureEventModel) View() string {
    switch m.state {
    case StateForm:
        return m.viewForm()
    case StateReview:
        return m.viewReview()
    case StateSubmit:
        return m.viewSubmit()
    default:
        return ""
    }
}

// viewForm renders the input form
func (m *CaptureEventModel) viewForm() string {
    // Create header
    headerStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(styles.ColorTextPrimary).
        MarginBottom(1)
    header := headerStyle.Render("📝 Capture Career Event")

    // Create form fields with styling
    titleField := m.renderFormField("Title", m.titleInput.View(), m.focusedField == 0)
    descField := m.renderFormField("Description", m.descriptionInput.View(), m.focusedField == 1)
    companyField := m.renderFormField("Company", m.companyInput.View(), m.focusedField == 2)

    // Create footer with help text
    footerStyle := lipgloss.NewStyle().
        Foreground(styles.ColorTextSecondary).
        MarginTop(2)
    footer := footerStyle.Render("Tab: Next Field | Shift+Tab: Prev Field | Enter: Review")

    // Combine all sections
    form := lipgloss.JoinVertical(
        lipgloss.Left,
        titleField,
        descField,
        companyField,
    )

    // Use CardContainer for consistent styling
    card := components.NewCardContainer().
        SetHeader(header).
        SetBody(form).
        SetFooter(footer).
        Render()

    return card
}

// renderFormField renders a single form field with label and input
func (m *CaptureEventModel) renderFormField(label string, input string, focused bool) string {
    // Label styling
    labelStyle := lipgloss.NewStyle().
        Foreground(styles.ColorTextPrimary).
        Width(12).
        Bold(true)

    // Input wrapper styling
    borderColor := styles.ColorBorder
    if focused {
        borderColor = styles.ColorBorderActive
    }

    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(borderColor).
        Padding(0, 1)

    // Combine label and input
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        labelStyle.Render(label),
        inputStyle.Render(input),
    )
}

// viewReview renders the review screen
func (m *CaptureEventModel) viewReview() string {
    headerStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(styles.ColorTextPrimary).
        MarginBottom(1)
    header := headerStyle.Render("✓ Review Event")

    // Display captured data
    dataStyle := lipgloss.NewStyle().
        Foreground(styles.ColorTextPrimary).
        MarginBottom(1)

    data := lipgloss.JoinVertical(
        lipgloss.Left,
        dataStyle.Render("Title: "+m.data.Title),
        dataStyle.Render("Description: "+m.data.Description),
        dataStyle.Render("Company: "+m.data.Company),
    )

    // Action buttons
    confirmBtn := styles.ButtonPrimary.Render("✓ Confirm")
    editBtn := styles.ButtonSecondary.Render("✎ Edit")

    actions := lipgloss.JoinHorizontal(
        lipgloss.Top,
        confirmBtn,
        editBtn,
    )

    card := components.NewCardContainer().
        SetHeader(header).
        SetBody(data).
        SetFooter(actions).
        Render()

    return card
}

// viewSubmit renders the success screen
func (m *CaptureEventModel) viewSubmit() string {
    successStyle := lipgloss.NewStyle().
        Foreground(styles.ColorSuccess).
        Bold(true)

    message := successStyle.Render("✓ Event captured successfully!")

    card := components.NewCardContainer().
        SetHeader("Success").
        SetBody(message).
        Render()

    return card
}
```

## Best Practices

### 1. Style Organization

```go
// ❌ Bad: Styles defined inline
func (m Model) View() string {
    style := lipgloss.NewStyle().
        Background(lipgloss.Color("#fff")).
        Foreground(lipgloss.Color("#000"))
    return style.Render("text")
}

// ✅ Good: Styles defined centrally
func (m Model) View() string {
    return styles.ButtonPrimary.Render("text")
}
```

### 2. Component Reuse

```go
// ❌ Bad: Duplicating card rendering
func (m Model) View() string {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Render(content1)
}

// ✅ Good: Using CardContainer
func (m Model) View() string {
    return components.NewCardContainer().
        SetBody(content1).
        Render()
}
```

### 3. Focus Management

```go
// ✅ Good: Clear focus state
func (m *Model) updateFocus() {
    m.input1.Blur()
    m.input2.Blur()

    if m.focused == 0 {
        m.input1.Focus()
    } else {
        m.input2.Focus()
    }
}
```

### 4. Responsive Layouts

```go
// ✅ Good: Adapt to terminal width
func (m Model) View() string {
    if m.width < 80 {
        return m.renderCompact()
    }
    return m.renderFull()
}
```

## Common Patterns

### Form with Multiple Fields

```go
type FormModel struct {
    fields []textinput.Model
    focused int
}

func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focused = (m.focused + 1) % len(m.fields)
            m.updateFocus()
        }
    }

    m.fields[m.focused], _ = m.fields[m.focused].Update(msg)
    return m, nil
}
```

### List Selection

```go
import "github.com/charmbracelet/bubbles/list"

type ListModel struct {
    items list.Model
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.items, cmd = m.items.Update(msg)
    return m, cmd
}
```

### Modal Dialogs

```go
type ModalModel struct {
    visible bool
    content string
}

func (m *ModalModel) View() string {
    if !m.visible {
        return m.renderBackground()
    }

    modal := components.NewCardContainer().
        SetHeader("Confirm").
        SetBody(m.content).
        Render()

    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        modal,
    )
}
```

## Testing Styles

```go
func TestCardRendering(t *testing.T) {
    card := components.NewCardContainer().
        SetHeader("Test").
        SetBody("Content").
        Render()

    // Verify card contains expected text
    assert.Contains(t, card, "Test")
    assert.Contains(t, card, "Content")
}
```

## Performance Tips

1. **Cache Styles**: Define styles once, reuse everywhere
2. **Minimize Re-renders**: Only update when state changes
3. **Use Responsive Widths**: Set widths based on terminal size
4. **Lazy Render**: Only render visible components

## Resources

- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [KaRiya Styles Reference](../cli/styles/styles.go)
- [KaRiya Components](../cli/components/)

## Summary

By combining lipgloss for styling and bubbles for interactivity, you can create:

- ✅ Beautiful, professional terminal UIs
- ✅ Consistent styling across the application
- ✅ Reusable, maintainable components
- ✅ Responsive layouts that adapt to terminal size
- ✅ Interactive forms and dialogs
- ✅ Accessible keyboard navigation

The key is to **centralize styling**, **reuse components**, and **keep concerns separated**.

