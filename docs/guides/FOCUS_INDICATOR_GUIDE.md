# Focus Indicator Guide

## Overview

This guide documents the standardized focus indicator patterns used across the KaRiya CLI application. Focus indicators provide visual feedback to users about which element is currently active or selected.

**Current Status**: ✅ All 8 refactored models use consistent focus indicators

---

## Focus Indicator Types

### 1. Form Field Focus Indicators

**Used In**: form.go, fact_editor.go, metadata_editor.go

**Visual Appearance**:
- Focused field: Border style changes to primary color, input text becomes highlighted
- Non-focused field: Standard border with secondary color
- Error state: Border color changes to error red (overrides focus state)

**Implementation Pattern**:

```go
// Track which field is focused
type FormModel struct {
    focusIndex int  // 0 = TextField, 1 = DateField, etc.
    // ... other fields
}

// Render field with focus state
func (m *FormModel) renderTextFieldWithContainer() string {
    focused := m.focusIndex == int(TextField)

    fieldContent := components.NewFormFieldContainer().
        SetLabel("Event Text (required):").
        SetInput(m.inputs[0].View()).
        SetHint(charInfo).
        SetError(fieldErr).
        SetFocused(focused).  // ← Pass focus state to container
        Render()

    return fieldContent
}

// Update focus on navigation
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusIndex = (m.focusIndex + 1) % numFields
        case "shift+tab":
            m.focusIndex = (m.focusIndex - 1 + numFields) % numFields
        }
    }
    return m, nil
}
```

**Style Application** (in FormFieldContainer):

```go
// Render input with focus styling
if ffc.hasInput {
    var inputStyle lipgloss.Style
    if ffc.hasError {
        inputStyle = styles.InputError  // Error takes precedence
    } else if ffc.isFocused {
        inputStyle = styles.InputFocused  // Blue border, highlighted text
    } else {
        inputStyle = styles.InputBase  // Standard border
    }
    parts = append(parts, inputStyle.Render(ffc.input))
}
```

**Style Definitions**:

```go
// From styles.go
InputFocused = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(ColorPrimary).      // Blue border
    Padding(1, 2)

InputBase = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorder).       // Gray border
    Padding(1, 2)

InputError = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(ColorError).        // Red border
    Padding(1, 2)
```

**Navigation Keys**:
- `Tab` / `Shift+Tab`: Move between fields
- `j` / `k`: Move between fields (vim-style)
- `Enter`: Submit form
- `Esc`: Cancel form

**Visual Example**:

```
┌─ Event Text (required): ─────────────────┐  ← Focused (blue border)
│ My career milestone                      │
└──────────────────────────────────────────┘

┌─ Date (required): ───────────────────────┐  ← Not focused (gray border)
│ 2025-12-31                               │
└──────────────────────────────────────────┘

Error Example:
┌─ Company (required): ────────────────────┐  ← Error state (red border)
│ [empty]                                  │
└──────────────────────────────────────────┘
Error: Company name is required
```

---

### 2. List Item Focus Indicators

**Used In**: list.go, fact_list.go, burst_list.go

**Visual Appearance**:
- Focused item: Marker character `▶` + bold, primary color text
- Non-focused items: Two spaces + standard color text
- Empty state: "No items found" message

**Implementation Pattern**:

```go
// Track which item is focused/selected
type ListModel struct {
    selectedIdx int  // Index of selected item
    events      []*career.Event
    // ... other fields
}

// Render items with focus indicator
func (m *ListModel) renderListItems() []string {
    var items []string

    for i, event := range m.events {
        // Marker for selected item
        marker := "  "
        if i == m.selectedIdx {
            marker = "▶ "  // ← Selection indicator
        }

        // Truncate text to fit
        text := event.Text
        if len(text) > 100 {
            text = text[:97] + "..."
        }

        // Apply styling based on selection
        itemStyle := styles.ListItem
        if i == m.selectedIdx {
            itemStyle = styles.ListItemSelected  // ← Bold, primary color
        }

        items = append(items, itemStyle.Render(marker + text))
    }

    return items
}

// Update focus on navigation
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j":
            m.moveNext()  // ← Move to next item
        case "k":
            m.movePrev()  // ← Move to previous item
        case "g":
            m.goToFirstItem()
        case "G":
            m.goToLastItem()
        }
    }
    return m, nil
}
```

**Style Definitions**:

```go
// From styles.go
ListItem = lipgloss.NewStyle().
    Foreground(ColorTextPrimary).      // Standard text color
    PaddingLeft(1)

ListItemSelected = lipgloss.NewStyle().
    Foreground(ColorPrimary).          // Primary color (blue)
    Bold(true).                        // Bold text
    PaddingLeft(1)
```

**Navigation Keys**:
- `j`: Move down (vim-style)
- `k`: Move up (vim-style)
- `g`: Go to first item
- `G`: Go to last item
- `PageUp`: Page up
- `PageDown`: Page down
- `Enter`: Select current item
- `q` / `Esc`: Exit list

**Visual Example**:

```
📋 Career Events
▶ My career milestone                    2025-12-31 | Company Inc.
  Another important event                2025-12-20 | Other Corp
  Third event in the list                2025-12-10 | Third Co.

Showing 1-3 of 15 events
(Press 'j/k' to navigate, 'Enter' to select, 'q' to quit)
```

---

### 3. Dialog Button Focus Indicators

**Used In**: confirmation_dialog.go

**Visual Appearance**:
- Focused button: Border style with primary color
- Non-focused button: Standard secondary style
- Destructive action: Confirm button shows in error red color

**Implementation Pattern**:

```go
// Track which button is focused
type ConfirmationDialog struct {
    focused bool  // true = confirm button, false = cancel button
    // ... other fields
}

// Render buttons with focus state
func (d *ConfirmationDialog) renderButtons() []string {
    cancelStyle := styles.ButtonSecondary
    confirmStyle := styles.ButtonSecondary

    if d.focused {
        confirmStyle = styles.ButtonFocused  // ← Confirm button focused
    } else {
        cancelStyle = styles.ButtonFocused  // ← Cancel button focused
    }

    cancelButton := cancelStyle.Render(d.cancelText)
    confirmButton := confirmStyle.
        Foreground(styles.ColorError).     // ← Error color for destructive
        Render(d.confirmText)

    return []string{cancelButton, confirmButton}
}

// Update focus on navigation
func (d *ConfirmationDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab", "right", "l":
            d.focused = !d.focused  // ← Toggle focus
        case "shift+tab", "left", "h":
            d.focused = !d.focused  // ← Toggle focus
        case "enter":
            if d.focused {
                // Confirm action
            } else {
                // Cancel action
            }
        }
    }
    return d, nil
}
```

**Style Definitions**:

```go
// From styles.go
ButtonFocused = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(ColorPrimary).    // Blue border
    Foreground(ColorPrimary).          // Blue text
    Padding(0, 2)

ButtonSecondary = lipgloss.NewStyle().
    Foreground(ColorTextSecondary).    // Gray text
    Padding(0, 2)
```

**Navigation Keys**:
- `Tab` / `Shift+Tab`: Switch between buttons
- `Left` / `Right`: Switch between buttons
- `h` / `l`: Switch between buttons (vim-style)
- `Enter`: Confirm selected button
- `Esc`: Cancel dialog

**Visual Example**:

```
Are you sure you want to delete this event?

┌─────────────────────────────┐
│  Cancel        Delete Event │  ← Delete button focused (blue border)
└─────────────────────────────┘
(Tab: Switch | Enter: Confirm | Esc: Cancel)
```

---

## Best Practices

### ✅ DO

1. **Track focus state in model fields**
   ```go
   type MyModel struct {
       focusIndex int  // For forms
       selectedIdx int // For lists
       focused bool    // For dialogs
   }
   ```

2. **Use FormFieldContainer for form fields**
   ```go
   components.NewFormFieldContainer().
       SetFocused(isFocused).
       Render()
   ```

3. **Use marker character for list items**
   ```go
   marker := "  "
   if i == m.selectedIdx {
       marker = "▶ "
   }
   ```

4. **Apply styles from styles/constants_export.go**
   ```go
   itemStyle := styles.ListItem
   if i == m.selectedIdx {
       itemStyle = styles.ListItemSelected
   }
   ```

5. **Update focus on key navigation**
   ```go
   case "tab":
       m.focusIndex = (m.focusIndex + 1) % numFields
   ```

6. **Test focus state changes**
   ```go
   It("should move focus to next field on Tab", func() {
       m, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
       Expect(m.focusIndex).To(Equal(1))
   })
   ```

### ❌ DON'T

1. **Don't use inline colors for focus indicators**
   ```go
   // ❌ WRONG
   inputStyle = lipgloss.NewStyle().
       BorderForeground(lipgloss.Color("#0080ff"))  // Inline color

   // ✅ CORRECT
   inputStyle = styles.InputFocused  // Use constant
   ```

2. **Don't mix focus indicator styles**
   ```go
   // ❌ WRONG - Mixing border and background colors
   inputStyle = styles.InputFocused.
       Background(lipgloss.Color("#ffffff"))

   // ✅ CORRECT - Use consistent styling
   inputStyle = styles.InputFocused
   ```

3. **Don't forget to update focus on navigation**
   ```go
   // ❌ WRONG - Navigation doesn't update focus
   case "j":
       m.moveNext()  // But focusIndex stays the same

   // ✅ CORRECT - Update focus with navigation
   case "j":
       m.focusIndex++
       m.moveNext()
   ```

4. **Don't render focus state without updating it**
   ```go
   // ❌ WRONG - Focus never changes
   focused := true  // Hardcoded

   // ✅ CORRECT - Focus based on state
   focused := m.focusIndex == int(currentField)
   ```

5. **Don't use different focus indicators for similar elements**
   ```go
   // ❌ WRONG - Different markers for similar lists
   // list.go uses "▶"
   // fact_list.go uses "►"

   // ✅ CORRECT - All lists use "▶"
   marker := "▶ "
   ```

---

## Common Patterns

### Pattern 1: Form with Multiple Fields

```go
type MyFormModel struct {
    fields   []*textinput.Model
    focusIndex int
}

func (m *MyFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusIndex = (m.focusIndex + 1) % len(m.fields)
        case "shift+tab":
            m.focusIndex = (m.focusIndex - 1 + len(m.fields)) % len(m.fields)
        }
    }
    return m, nil
}

func (m *MyFormModel) View() string {
    var parts []string

    for i, field := range m.fields {
        focused := i == m.focusIndex

        container := components.NewFormFieldContainer().
            SetLabel(m.labels[i]).
            SetInput(field.View()).
            SetFocused(focused).
            Render()

        parts = append(parts, container)
    }

    return strings.Join(parts, "\n")
}
```

### Pattern 2: List with Selection

```go
type MyListModel struct {
    items       []string
    selectedIdx int
}

func (m *MyListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j":
            if m.selectedIdx < len(m.items)-1 {
                m.selectedIdx++
            }
        case "k":
            if m.selectedIdx > 0 {
                m.selectedIdx--
            }
        }
    }
    return m, nil
}

func (m *MyListModel) View() string {
    var items []string

    for i, item := range m.items {
        marker := "  "
        if i == m.selectedIdx {
            marker = "▶ "
        }

        style := styles.ListItem
        if i == m.selectedIdx {
            style = styles.ListItemSelected
        }

        items = append(items, style.Render(marker + item))
    }

    return strings.Join(items, "\n")
}
```

### Pattern 3: Dialog with Button Focus

```go
type MyDialogModel struct {
    buttons []string
    focused bool  // true = first button, false = second button
}

func (m *MyDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab", "right":
            m.focused = !m.focused
        case "shift+tab", "left":
            m.focused = !m.focused
        }
    }
    return m, nil
}

func (m *MyDialogModel) View() string {
    var buttons []string

    for i, label := range m.buttons {
        isFocused := (i == 0 && m.focused) || (i == 1 && !m.focused)

        style := styles.ButtonSecondary
        if isFocused {
            style = styles.ButtonFocused
        }

        buttons = append(buttons, style.Render(label))
    }

    return strings.Join(buttons, "  ")
}
```

---

## Testing Focus Indicators

### Unit Test Example

```go
Describe("FormModel Focus", func() {
    It("should set focus on first field initially", func() {
        m := NewFormModel()
        Expect(m.focusIndex).To(Equal(0))
    })

    It("should move focus to next field on Tab", func() {
        m := NewFormModel()
        m.Update(tea.KeyMsg{Type: tea.KeyTab})
        Expect(m.focusIndex).To(Equal(1))
    })

    It("should show focus indicator in rendered output", func() {
        m := NewFormModel()
        view := m.View()

        // First field should have focused styling
        Expect(view).To(ContainSubstring(
            styles.InputFocused.Render("test"),
        ))
    })

    It("should wrap focus from last to first field", func() {
        m := NewFormModel()
        m.focusIndex = numFields - 1
        m.Update(tea.KeyMsg{Type: tea.KeyTab})
        Expect(m.focusIndex).To(Equal(0))
    })
})
```

### Integration Test Example

```go
Describe("List Focus Navigation", func() {
    It("should highlight selected item with marker", func() {
        m := NewListModel()
        m.selectedIdx = 1
        view := m.View()

        // Second item should have marker
        Expect(view).To(MatchRegexp(`▶\s+second item`))
    })

    It("should update selection on j key", func() {
        m := NewListModel()
        m.Update(tea.KeyMsg{Runes: []rune{'j'}})
        Expect(m.selectedIdx).To(Equal(1))
    })

    It("should apply selected style to focused item", func() {
        m := NewListModel()
        m.selectedIdx = 0
        view := m.View()

        // First item should have ListItemSelected style (bold)
        Expect(view).To(ContainSubstring(
            styles.ListItemSelected.Render("▶ first item"),
        ))
    })
})
```

---

## Troubleshooting

### Issue: Focus indicator not showing

**Cause**: Focus state not being passed to component

**Solution**:
```go
// ❌ WRONG - focused always false
focused := false
container.SetFocused(focused)

// ✅ CORRECT - focused based on state
focused := m.focusIndex == int(currentField)
container.SetFocused(focused)
```

### Issue: Focus not changing on navigation

**Cause**: Navigation event not updating focus index

**Solution**:
```go
// ❌ WRONG - focusIndex not updated
case "tab":
    m.moveNext()  // Only moves cursor, not focus

// ✅ CORRECT - Update focus index
case "tab":
    m.focusIndex = (m.focusIndex + 1) % numFields
    m.moveNext()
```

### Issue: Multiple items showing focus indicator

**Cause**: Focus condition checking wrong index

**Solution**:
```go
// ❌ WRONG - Multiple items match condition
if i < m.focusIndex {
    marker = "▶ "
}

// ✅ CORRECT - Only one item matches
if i == m.focusIndex {
    marker = "▶ "
}
```

### Issue: Focus indicator looks different in different models

**Cause**: Using different styles for same element type

**Solution**:
```go
// ❌ WRONG - Different styles
// form.go uses styles.InputFocused
// fact_editor.go uses custom blue border

// ✅ CORRECT - All use same style constant
container.SetFocused(focused)  // Uses styles.InputFocused
```

---

## Related Documentation

- [Error Handling Guide](./ERROR_HANDLING_GUIDE.md) - Error display alongside focus
- [Style Usage Guide](./STYLE_USAGE_GUIDE.md) - All available styles
- [Component Usage Guide](./COMPONENT_USAGE_GUIDE.md) - Container components
- [Task 4.7 Audit](../audits/TASK_4.7_FOCUS_INDICATOR_AUDIT.md) - Focus indicator audit results

---

## Summary

Focus indicators provide essential visual feedback in the CLI application:

- **Form Fields**: Border color change + text highlighting via FormFieldContainer
- **List Items**: Marker character `▶` + bold text via ListItemSelected style
- **Dialog Buttons**: Border color change via ButtonFocused style

All implementations:
- ✅ Use exported color constants from styles package
- ✅ Follow consistent patterns across similar element types
- ✅ Include proper navigation key handling
- ✅ Are thoroughly tested with unit and integration tests

When adding new models, follow these patterns to maintain consistency across the application.

---

**Last Updated**: 2026-01-01
**Status**: Production Ready
**Compliance**: ✅ All 8 refactored models follow these patterns

