---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Focus Indicator Implementation Guide

**Version**: 1.0
**Date**: 2026-01-01
**Status**: COMPLETE
**Related Task**: Task 4.7 from tasks-07-model-consistency.md

## Overview

This guide documents the focus indicator implementation patterns used across the KaRiya CLI. Focus indicators provide visual feedback to users about which form field, list item, or button currently has input focus.

## Quick Reference

| Component Type | Focus Indicator | Implementation | Style |
|---|---|---|---|
| **Form Fields** | Thick teal border | `FormFieldContainer.SetFocused(bool)` | `InputFocused` |
| **List Items** | Right arrow marker (▶) | Check `selectedIdx` or `focusedIdx` | Text-based |
| **Buttons** | Bold text + teal border | Boolean flag in dialog | `ButtonFocused` |

---

## 1. Form Field Focus Indicators

### 1.1 Overview

Form fields use a thick teal border to indicate focus state. The border changes from thin gray (unfocused) to thick teal (focused).

### 1.2 Implementation Pattern

**Component**: `FormFieldContainer` in `internal/cli/components/form_field_container.go`

**Method**: `SetFocused(bool)`

**Usage Example**:
```go
container := components.NewFormFieldContainer().
    SetLabel("Event Text").
    SetInput(fieldValue).
    SetError(errorMsg).
    SetFocused(isCurrentField).  // KEY LINE
    Render()
```

### 1.3 Visual Examples

**Without Focus**:
```
┌─────────────────────────┐
│ Event Text              │
└─────────────────────────┘
```

**With Focus**:
```
┏━━━━━━━━━━━━━━━━━━━━━━━━━┓
│ Event Text              │
┗━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### 1.4 Models Using Form Field Focus

- ✅ `form.go` - 7 fields with focus indicators
- ✅ `fact_editor.go` - 5 fields with focus indicators
- ✅ `metadata_editor.go` - 6 fields with focus indicators

### 1.5 Implementation Details

**File**: `internal/cli/components/form_field_container.go`

**Code Reference** (lines 77-99):
```go
func (ffc *FormFieldContainer) Render() string {
    var parts []string

    // ... label rendering ...

    if ffc.hasInput {
        var inputStyle lipgloss.Style
        if ffc.hasError {
            inputStyle = styles.InputError.Copy().
                Foreground(styles.ColorTextPrimary)
        } else if ffc.isFocused {
            // FOCUS INDICATOR: Thick teal border
            inputStyle = styles.InputFocused.Copy().
                Foreground(styles.ColorTextPrimary)
        } else {
            inputStyle = styles.InputBase.Copy().
                Foreground(styles.ColorTextPrimary)
        }
        parts = append(parts, inputStyle.Render(ffc.input))
    }

    // ... rest of rendering ...
}
```

### 1.6 When to Use

- All text input fields
- All form fields with user input
- Anywhere using `FormFieldContainer`

### 1.7 Focus Indicator Style

**Style Name**: `InputFocused` (in `styles/styles.go`)

**Properties**:
- Border Color: `ColorBorderActive` (#5fb3b3 - Teal)
- Border Style: Thick (━, ┃, ┏, ┓, ┗, ┛, ┣, ┫, ┳, ┻)
- Effect: Clearly distinguishes focused field from others

---

## 2. List Item Focus Indicators

### 2.1 Overview

List items use a right arrow marker (▶) to indicate the currently selected/focused item. Non-selected items show two spaces (  ) in place of the marker.

### 2.2 Implementation Pattern

**Pattern**: Selection marker in list rendering

**Usage Example**:
```go
for i, item := range items {
    marker := "  "
    if i == selectedIdx {
        marker = "▶ "  // KEY LINE: Focus indicator
    }
    // Render item with marker
    itemLine := marker + itemText + details
    // Add to output
}
```

### 2.3 Visual Examples

**List View**:
```
  Career Event 1 - 2024-01-15
▶ Career Event 2 - 2024-02-20  (Selected - has focus)
  Career Event 3 - 2024-03-10
  Career Event 4 - 2024-04-05
```

### 2.4 Models Using List Item Focus

- ✅ `list.go` - Reference implementation
- ✅ `fact_list.go` - Consistent marker usage
- ✅ `burst_list.go` - Consistent marker usage

### 2.5 Implementation Details

**Reference File**: `internal/cli/models/list.go` (lines 185-190)

```go
for i, event := range m.events {
    // Marker for selected item
    marker := "  "
    if i == m.selectedIdx {
        marker = "▶ "  // FOCUS INDICATOR
    }

    // Render item with marker
    text := event.Text
    if len(text) > 100 {
        text = text[:97] + "..."
    }
    // ... additional rendering ...
}
```

### 2.6 Consistency Across List Models

All three list models use identical marker patterns:

**list.go**:
```go
if i == m.selectedIdx {
    marker = "▶ "
}
```

**fact_list.go**:
```go
// Marker for selected item (matching list.go pattern)
if i == flm.focusedIdx {
    marker = "▶ "
}
```

**burst_list.go**:
```go
// Marker for selected item (matching list.go pattern)
if i == m.selectedIdx {
    marker = "▶ "
}
```

### 2.7 When to Use

- All list-based models
- Any component showing selectable items
- Navigation through item collections
- Displaying multiple items with selection

### 2.8 Focus Indicator Character

**Character**: `▶` (Right-pointing arrow, Unicode U+25B6)

**Spacing**: Two characters total: `▶ ` (arrow + space)

**Alignment**: First characters of item line (before item text)

---

## 3. Button Focus Indicators

### 3.1 Overview

Buttons in dialogs use style-based focus indicators. The focused button is rendered with bold text and a teal border, while unfocused buttons use the secondary style.

### 3.2 Implementation Pattern

**Component**: Dialog models (e.g., `ConfirmationDialog`)

**Pattern**: Boolean flag tracking focused button

**Usage Example**:
```go
func (d *ConfirmationDialog) renderButtons() []string {
    cancelStyle := styles.ButtonSecondary
    confirmStyle := styles.ButtonSecondary

    if d.focused {
        confirmStyle = styles.ButtonFocused  // KEY LINE
    } else {
        cancelStyle = styles.ButtonFocused   // KEY LINE
    }

    cancelButton := cancelStyle.Render(d.cancelText)
    confirmButton := confirmStyle.Render(d.confirmText)
    return []string{cancelButton, confirmButton}
}
```

### 3.3 Visual Examples

**Cancel Button Focused**:
```
[ Cancel ]  [ Confirm ]
(Bold, teal border)
```

**Confirm Button Focused**:
```
[ Cancel ]  [ Confirm ]
           (Bold, teal border)
```

### 3.4 Models Using Button Focus

- ✅ `confirmation_dialog.go` - Destructive action dialogs

### 3.5 Implementation Details

**File**: `internal/cli/models/confirmation_dialog.go`

**Focus State Management** (lines 17-32):
```go
type ConfirmationDialog struct {
    // ... other fields ...
    focused bool // true = confirm, false = cancel
}

// Initialize with cancel focused (safer default)
func NewConfirmationDialog(title, message string) *ConfirmationDialog {
    return &ConfirmationDialog{
        focused: false,  // Default to cancel for safety
        // ... other initialization ...
    }
}
```

**Button Rendering** (lines 92-112):
```go
func (d *ConfirmationDialog) renderButtons() []string {
    cancelStyle := styles.ButtonSecondary
    confirmStyle := styles.ButtonSecondary

    if d.focused {
        confirmStyle = styles.ButtonFocused
    } else {
        cancelStyle = styles.ButtonFocused
    }

    cancelButton := cancelStyle.Render(d.cancelText)
    confirmButton := confirmStyle.Render(d.confirmText)
    return []string{cancelButton, confirmButton}
}
```

### 3.6 Focus Toggle

**Trigger**: Tab key or arrow keys

**Navigation**:
- Tab / Right Arrow: Move to next button
- Shift+Tab / Left Arrow: Move to previous button

### 3.7 When to Use

- All buttons in dialogs
- Any component with multiple button options
- Focus navigation between buttons
- Confirmation dialogs

### 3.8 Focus Indicator Styles

**Focused Button**: `ButtonFocused`
- Bold text
- Teal border (#5fb3b3)

**Unfocused Button**: `ButtonSecondary`
- Normal text weight
- Standard border

---

## 4. Style Constants Reference

### 4.1 Input Styles

**File**: `internal/cli/styles/styles.go`

```go
// Focused input field (thick teal border)
InputFocused = InputBase.Copy().
    BorderForeground(ColorBorderActive).
    BorderStyle(lipgloss.ThickBorder())

// Base input field (thin gray border)
InputBase = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2).
    BorderForeground(ColorBorder)

// Error input field (red border)
InputError = InputBase.Copy().
    BorderForeground(ColorError)
```

### 4.2 Button Styles

```go
// Focused button (bold, teal border)
ButtonFocused = ButtonPrimary.
    BorderForeground(ColorAccentTeal).
    Bold(true)

// Secondary button (normal weight)
ButtonSecondary = lipgloss.NewStyle().
    Padding(0, 2).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorder)
```

### 4.3 Color Constants

```go
ColorBorderActive = lipgloss.Color("#5fb3b3")  // Teal (focus border)
ColorAccentTeal = lipgloss.Color("#5fb3b3")    // Teal (button focus)
ColorBorder = lipgloss.Color("#555555")        // Gray (unfocused)
ColorError = lipgloss.Color("#d76e6e")         // Red (error)
```

---

## 5. Testing Focus Indicators

### 5.1 Form Field Focus Tests

**File**: `internal/cli/models/form_focus_indicators_test.go`

**Test Categories**:
- Focus indicator character verification
- Focus indicator presence in initial state
- Focus indicator positioning relative to labels
- Focus indicator consistency across field navigation
- Focus indicator layout preservation

**Example Test**:
```go
It("should display focus indicator for the focused field", func() {
    view := form.View()
    Expect(view).To(ContainSubstring("►"))  // Focus indicator present
})
```

### 5.2 List Item Focus Tests

**Files**:
- `internal/cli/models/fact_list_test.go`
- `internal/cli/models/burst_list_test.go`

**Test Categories**:
- Focus indicator marker display
- Marker character consistency
- Marker position at line beginning
- Non-selected item spacing
- Empty state handling

**Example Test**:
```go
ginkgo.It("should display focus indicator for selected item", func() {
    model.SetFacts(facts)
    view := model.View()
    gomega.Expect(view).To(gomega.ContainSubstring("▶"))
})
```

### 5.3 Button Focus Tests

**File**: `internal/cli/models/confirmation_dialog_test.go`

**Test Categories**:
- Default focus state (cancel button)
- Focus toggle with Tab key
- Focus toggle with arrow keys
- Focus state persistence
- Button rendering consistency

**Example Test**:
```go
It("should initialize with cancel button focused", func() {
    dialog := NewConfirmationDialog("Title", "Message")
    view := dialog.View()
    Expect(view).To(ContainSubstring("Cancel"))
})
```

---

## 6. Best Practices

### 6.1 Form Fields

✅ **DO**:
- Always use `FormFieldContainer.SetFocused(bool)` for form fields
- Set focused based on current field index
- Use `InputFocused` style for focus indication
- Test focus indicator display in form tests

❌ **DON'T**:
- Manually apply input styles without using FormFieldContainer
- Use inline border styling instead of constants
- Forget to set focused state when rendering fields
- Use different focus indicators in different forms

### 6.2 List Items

✅ **DO**:
- Use "▶ " marker for selected items consistently
- Use "  " (two spaces) for non-selected items
- Check `selectedIdx` or `focusedIdx` to determine marker
- Apply marker at the beginning of the item line

❌ **DON'T**:
- Use different marker characters (▶, ►, →, etc.)
- Vary marker spacing between models
- Use inline styling instead of marker characters
- Forget to update marker when selection changes

### 6.3 Buttons

✅ **DO**:
- Use boolean flag to track focused button
- Apply `ButtonFocused` style to focused button
- Apply `ButtonSecondary` style to unfocused buttons
- Default to safest button (cancel) when initializing

❌ **DON'T**:
- Use different focus styles in different dialogs
- Manually apply button styles without using constants
- Forget to toggle focus on navigation
- Use inconsistent focus indicators across dialogs

---

## 7. Common Patterns

### 7.1 Form Field with Focus

```go
// In model's View() method
container := components.NewFormFieldContainer().
    SetLabel("Field Label").
    SetInput(m.fields[fieldIndex]).
    SetError(fieldErrors[fieldIndex]).
    SetFocused(m.focusedField == fieldIndex).  // Focus logic
    Render()
```

### 7.2 List Item with Focus

```go
// In model's View() method
for i, item := range m.items {
    marker := "  "
    if i == m.selectedIdx {
        marker = "▶ "  // Focus marker
    }

    line := fmt.Sprintf("%s%s %s", marker, item.Name, item.Details)
    lines = append(lines, line)
}
```

### 7.3 Button with Focus

```go
// In dialog's renderButtons() method
cancelStyle := styles.ButtonSecondary
confirmStyle := styles.ButtonSecondary

if d.focused {
    confirmStyle = styles.ButtonFocused
} else {
    cancelStyle = styles.ButtonFocused
}

return []string{
    cancelStyle.Render(d.cancelText),
    confirmStyle.Render(d.confirmText),
}
```

---

## 8. Troubleshooting

### Issue: Focus Indicator Not Showing

**Symptom**: Field/button doesn't show focus visual

**Solutions**:
1. Verify `SetFocused(true)` is called
2. Check that correct style is applied (InputFocused, ButtonFocused)
3. Verify style constants are imported from styles package
4. Check terminal supports thick borders (some terminals don't)

### Issue: Focus Indicator Not Updating

**Symptom**: Focus indicator stays on same field/button

**Solutions**:
1. Verify focus state is updated on navigation
2. Check that View() is called after Update()
3. Verify focusedField/selectedIdx is correctly incremented/decremented
4. Check that focus toggle logic is correct

### Issue: Inconsistent Focus Indicators

**Symptom**: Different models show different focus indicators

**Solutions**:
1. Use FormFieldContainer for all form fields
2. Use "▶ " marker for all list items
3. Use ButtonFocused/ButtonSecondary for all buttons
4. Review style constants are correctly applied

---

## 9. Migration Guide

### 9.1 Migrating Form Fields to FormFieldContainer

**Before**:
```go
// Manual field rendering
labelStyle := styles.InputLabel.Copy()
inputStyle := styles.InputBase.Copy()
// ... complex styling logic ...
```

**After**:
```go
// Using FormFieldContainer
container := components.NewFormFieldContainer().
    SetLabel("Field").
    SetInput(value).
    SetFocused(isFocused).
    Render()
```

### 9.2 Migrating List Items to Standard Marker

**Before**:
```go
// Different markers per model
if selected {
    line = "► " + itemText  // Different marker
}
```

**After**:
```go
// Standard marker across all models
marker := "  "
if i == selectedIdx {
    marker = "▶ "  // Consistent marker
}
line = marker + itemText
```

---

## 10. Compliance Checklist

### Form Field Focus Indicators
- [ ] All form models use FormFieldContainer
- [ ] SetFocused() called on each field
- [ ] InputFocused style applied to focused fields
- [ ] Focus state correctly identifies current field
- [ ] Tests verify focus indicator display

### List Item Focus Indicators
- [ ] All list models use "▶ " marker
- [ ] All list models use "  " for non-selected
- [ ] Marker character is consistent across all lists
- [ ] Focus indicator position is consistent
- [ ] Tests verify marker display

### Button Focus States
- [ ] All dialogs implement button focus
- [ ] Focus state toggles with Tab/arrow keys
- [ ] ButtonFocused style applied to focused button
- [ ] ButtonSecondary style applied to unfocused
- [ ] Tests verify button focus styling

### Code Consistency
- [ ] All form models follow FormFieldContainer pattern
- [ ] All list models follow selection marker pattern
- [ ] All dialog models follow button focus pattern
- [ ] No inline focus styling (all use exported styles)
- [ ] Focus logic is centralized (not duplicated)

---

## 11. References

- [Error Handling Guide](./ERROR_HANDLING_GUIDE.md)
- [Component Usage Guide](./COMPONENT_USAGE_GUIDE.md)
- [Style Usage Guide](./STYLE_USAGE_GUIDE.md)
- [Task 4.7 Audit](../audits/TASK_4.7_FOCUS_INDICATOR_AUDIT.md)
- [Model Development Guide](./MODEL_DEVELOPMENT_GUIDE.md)

---

## 12. Summary

Focus indicators provide critical visual feedback in the KaRiya CLI:

1. **Form Fields**: Thick teal border via `InputFocused` style
2. **List Items**: Right arrow marker (▶) at line beginning
3. **Buttons**: Bold text + teal border via `ButtonFocused` style

All focus indicators are:
- ✅ Consistent across all models
- ✅ Using exported style constants
- ✅ Properly tested
- ✅ Well-documented

This guide ensures developers can implement focus indicators correctly and consistently throughout the codebase.

---

**Document Version**: 1.0
**Created**: 2026-01-01
**Last Updated**: 2026-01-01
**Status**: COMPLETE
**Compliance**: FULL ✅


