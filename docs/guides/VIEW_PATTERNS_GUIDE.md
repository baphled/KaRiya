# View Patterns Guide

**Date**: 2026-01-01
**Version**: 1.0
**Status**: Complete

## Overview

This guide documents the standardized View() method patterns used across the KaRiya CLI application. Following these patterns ensures consistent UI/UX, maintainability, and code quality.

---

## Quick Reference

### Pattern Types

1. **Screen-Based Views** (ScreenContainer pattern)
2. **List-Based Views** (ListContainer pattern)
3. **Form-Based Views** (FormFieldContainer pattern)
4. **Modal/Dialog Views** (Custom pattern)

---

## 1. Screen-Based Views (ScreenContainer Pattern)

### When to Use
- Detail/review screens
- Single item viewing
- Complex content displays
- Views with headers, content, and footers

### Structure

```
┌─────────────────────────────────────────┐
│ Header (Title, Breadcrumbs)             │
├─────────────────────────────────────────┤
│                                         │
│ Content (ScreenContainer-wrapped)       │
│                                         │
├─────────────────────────────────────────┤
│ Footer (Help text, Navigation hints)    │
└─────────────────────────────────────────┘
```

### Implementation Template

```go
// View renders using the ScreenContainer pattern
func (m *YourModel) View() string {
    headerContent := m.renderHeader()
    contentArea := m.renderContent()
    footerContent := m.renderFooter()

    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerContent,
        "",
        contentArea,
        "",
        footerContent,
    )
}

// renderHeader renders the header section
func (m *YourModel) renderHeader() string {
    m.header.SetWidth(m.width)
    return m.header.View()
}

// renderContent renders the main content
func (m *YourModel) renderContent() string {
    var content []string

    // Build your content here
    content = append(content, "Your content")

    // Wrap in ScreenContainer for consistent padding
    screenContainer := components.NewScreenContainer(strings.Join(content, "\n")).
        WithPaddingMode(components.PaddingNormal)

    return screenContainer.Render()
}

// renderFooter renders the footer section
func (m *YourModel) renderFooter() string {
    m.helpFooter.SetWidth(m.width)
    return m.helpFooter.View()
}
```

### Reference Implementations

- **ViewEventModel**: `internal/cli/models/view_event.go`
- **ViewEventWithFactsModel**: `internal/cli/models/view_event_with_facts.go`
- **FactsResultsModel**: `internal/cli/models/facts_results.go`
- **ActionMenuModel**: `internal/cli/models/action_menu.go`
- **BurstSuggestionModel**: `internal/cli/models/burst_suggestion.go`

### Key Points

1. **Header Management**
   - Create header in constructor: `components.NewHeader("Title", width)`
   - Update width in WindowSizeMsg handler
   - Render via helper method

2. **Content Area**
   - Build content as string slice
   - Wrap in `ScreenContainer` for padding
   - Use `WithPaddingMode()` for layout control

3. **Footer Management**
   - Create footer in constructor: `components.NewHelpFooter("context_key", width)`
   - Update width in WindowSizeMsg handler
   - Render via helper method

4. **Spacing**
   - Use empty strings to add vertical spacing
   - Consistent spacing: header, empty line, content, empty line, footer

---

## 2. List-Based Views (ListContainer Pattern)

### When to Use
- Multiple items display
- Scrollable lists
- Selectable item lists
- Views with pagination

### Structure

```
┌─────────────────────────────────────────┐
│ Header (Title, Status)                  │
├─────────────────────────────────────────┤
│ Status Bar (Pagination, Filters)        │
│                                         │
│ List Items (with selection indicator)   │
│ - Item 1 (selected)                     │
│ - Item 2                                │
│ - Item 3                                │
│                                         │
├─────────────────────────────────────────┤
│ Footer (Help text, Navigation hints)    │
└─────────────────────────────────────────┘
```

### Implementation Template

```go
// View renders using the ListContainer pattern
func (m *YourListModel) View() string {
    headerContent := m.renderHeader()
    statusContent := m.renderStatus()
    listContent := m.renderList()
    footerContent := m.renderFooter()

    // Handle errors
    var errorContent string
    if m.err != nil {
        errorContent = styles.ErrorBox.Render(fmt.Sprintf("Error: %v. Press 'r' to retry", m.err))
    }

    if errorContent != "" {
        return lipgloss.JoinVertical(
            lipgloss.Left,
            headerContent,
            "",
            errorContent,
            "",
            footerContent,
        )
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerContent,
        statusContent,
        "",
        listContent,
        "",
        footerContent,
    )
}

// renderHeader renders the header
func (m *YourListModel) renderHeader() string {
    m.header.SetWidth(m.width)
    return m.header.View()
}

// renderStatus renders pagination/status bar
func (m *YourListModel) renderStatus() string {
    statusText := fmt.Sprintf("Showing %d-%d of %d items",
        m.scrollOffset+1, m.endIdx, len(m.items))
    return styles.InputHint.Render(statusText)
}

// renderList renders the list items
func (m *YourListModel) renderList() string {
    var items []string
    for i, item := range m.items {
        if i == m.focusedIdx {
            items = append(items, m.renderSelectedItem(item))
        } else {
            items = append(items, m.renderItem(item))
        }
    }
    return strings.Join(items, "\n")
}

// renderFooter renders the footer
func (m *YourListModel) renderFooter() string {
    m.helpFooter.SetWidth(m.width)
    return m.helpFooter.View()
}
```

### Reference Implementations

- **ListModel**: `internal/cli/models/list.go`
- **FactListModel**: `internal/cli/models/fact_list.go`
- **BurstListModel**: `internal/cli/models/burst_list.go`
- **MetadataReviewModel**: `internal/cli/models/metadata_review.go`

### Key Points

1. **Status Bar**
   - Always show pagination: "Showing X-Y of Z items"
   - Include filter/sort status if applicable
   - Use `styles.InputHint` for consistent styling

2. **List Items**
   - Selection indicator: "▶ " for selected, "  " for unselected
   - Use `styles.ListItemSelected` for selected items
   - Use `styles.ListItem` for regular items

3. **Error Handling**
   - Use `styles.ErrorBox` for model-level errors
   - Include recovery guidance: "Press 'r' to retry"
   - Return early with error instead of showing list

4. **Empty State**
   - Check if list is empty
   - Show helpful message: "No items to display"
   - Suggest action: "Create one to get started"

---

## 3. Form-Based Views (FormFieldContainer Pattern)

### When to Use
- Input forms
- Event/fact editing
- Configuration screens
- Multi-field data entry

### Structure

```
┌─────────────────────────────────────────┐
│ Header (Title)                          │
├─────────────────────────────────────────┤
│ Form Fields:                            │
│ ┌─────────────────────────────────────┐ │
│ │ Label:                              │ │
│ │ [Input field with focus indicator]  │ │
│ │ Error message (if any)              │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ ┌─────────────────────────────────────┐ │
│ │ Label:                              │ │
│ │ [Input field]                       │ │
│ └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│ Footer (Help text, Navigation hints)    │
└─────────────────────────────────────────┘
```

### Implementation Template

```go
// View renders using the FormFieldContainer pattern
func (m *YourFormModel) View() string {
    headerContent := m.renderHeader()
    formContent := m.renderForm()
    footerContent := m.renderFooter()

    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerContent,
        "",
        formContent,
        "",
        footerContent,
    )
}

// renderHeader renders the header
func (m *YourFormModel) renderHeader() string {
    m.header.SetWidth(m.width)
    return m.header.View()
}

// renderForm renders all form fields
func (m *YourFormModel) renderForm() string {
    var fields []string

    // Add form fields
    for i, field := range m.fields {
        fieldContainer := components.NewFormFieldContainer()
        fieldContainer.SetLabel(field.Label)
        fieldContainer.SetValue(field.Value)
        if field.Focused {
            fieldContainer.SetFocused()
        }
        if field.Error != "" {
            fieldContainer.SetError(field.Error)
        }
        fields = append(fields, fieldContainer.View())
    }

    // Add model-level errors if any
    if m.err != nil {
        fields = append(fields, styles.ErrorBox.Render(m.err.Error()))
    }

    return strings.Join(fields, "\n\n")
}

// renderFooter renders the footer
func (m *YourFormModel) renderFooter() string {
    m.helpFooter.SetWidth(m.width)
    return m.helpFooter.View()
}
```

### Reference Implementations

- **FormModel**: `internal/cli/models/form.go`
- **FactEditorModel**: `internal/cli/models/fact_editor.go`
- **MetadataEditorModel**: `internal/cli/models/metadata_editor.go`

### Key Points

1. **Field Containers**
   - Use `FormFieldContainer` for consistent field rendering
   - Set label, value, and error via methods
   - Mark focused field with `SetFocused()`

2. **Error Handling**
   - Field-level errors: Use `FormFieldContainer.SetError()`
   - Model-level errors: Use `styles.ErrorBox`
   - Always include recovery guidance

3. **Validation Feedback**
   - Show field-level errors immediately
   - Use red color for errors (styles.ColorError)
   - Clear errors when field is corrected

---

## 4. Modal/Dialog Views

### When to Use
- Confirmation dialogs
- Success messages
- Quick menus
- Temporary overlays

### Structure

```
┌──────────────────────────────┐
│ Title / Header               │
├──────────────────────────────┤
│ Message / Content            │
├──────────────────────────────┤
│ Action Buttons or Options    │
│ (Yes/No, Confirm, etc.)      │
└──────────────────────────────┘
```

### Implementation Template

```go
// View renders a modal dialog
func (m *YourModalModel) View() string {
    var content []string

    // Title
    title := styles.HeaderMain.Render("Dialog Title")
    content = append(content, title)
    content = append(content, "")

    // Message
    message := "Your message here"
    content = append(content, message)
    content = append(content, "")

    // Options/Buttons
    options := "[Y]es  [N]o  [Esc] Cancel"
    content = append(content, styles.InputHint.Render(options))

    // Wrap in card
    card := styles.CardBase.
        Width(styles.MaxWidth(m.width) - 4).
        Render(strings.Join(content, "\n"))

    return card
}
```

### Reference Implementations

- **ConfirmationDialogModel**: `internal/cli/models/confirmation_dialog.go`
- **SuccessModel**: `internal/cli/models/success.go`

### Key Points

1. **Sizing**
   - Width: MaxWidth - 4 (for padding)
   - Height: Minimal, content-driven
   - Centered on screen (handled by parent)

2. **Styling**
   - Use `styles.CardBase` for wrapper
   - Use `styles.HeaderMain` for title
   - Use `styles.InputHint` for options

3. **Content**
   - Keep messages concise
   - Clear action buttons
   - Keyboard shortcuts in help text

---

## Error Handling Standards

### Field-Level Errors

```go
fieldContainer := components.NewFormFieldContainer()
fieldContainer.SetError("This field is required")
fieldContainer.View() // Displays error below field
```

**Styling**: `styles.ColorError` (red)
**Usage**: Form fields, input validation

### Model-Level Errors

```go
if m.err != nil {
    errorBox := styles.ErrorBox.Render(fmt.Sprintf("Error: %v. Press 'r' to retry", m.err))
    return errorBox
}
```

**Styling**: `styles.ErrorBox` (error box with background)
**Usage**: Loading errors, data fetch failures
**Guidance**: Always include recovery action

### Error Messages

- **Specific**: "Failed to load facts: database connection timeout"
- **Actionable**: "Press 'r' to retry" or "Check your connection and try again"
- **Consistent**: All errors use same format and style

---

## Color & Style Usage

### Text Styles

```go
styles.HeaderMain       // Main headers
styles.HeaderSecondary  // Sub-headers
styles.InputLabel       // Form labels
styles.InputHint        // Help text, hints
styles.ListItem         // Regular list items
styles.ListItemSelected // Selected list items
styles.ListItemFocused  // Focused list items
styles.ErrorBox         // Error messages
styles.ErrorText        // Field-level errors
styles.InfoText         // Information messages
styles.InfoBox          // Information boxes
```

### Colors

```go
styles.ColorText           // Primary text
styles.ColorTextSecondary  // Secondary text
styles.ColorError          // Error text (#d76e6e)
styles.ColorSuccess        // Success text (#6ee7b7)
styles.ColorInfo           // Info text (#a8dadc)
styles.ColorWarning        // Warning text (#f4a261)
```

---

## Best Practices

### 1. Separation of Concerns

✅ **Good**:
```go
func (m *Model) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        "",
        m.renderContent(),
        "",
        m.renderFooter(),
    )
}
```

❌ **Bad**:
```go
func (m *Model) View() string {
    // 100 lines of rendering logic mixed together
}
```

### 2. Width Management

✅ **Good**:
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.header.SetWidth(m.width)
    m.helpFooter.SetWidth(m.width)
```

❌ **Bad**:
```go
// Hardcoded widths
m.header.SetWidth(80)
m.helpFooter.SetWidth(80)
```

### 3. Error Recovery

✅ **Good**:
```go
if m.err != nil {
    return styles.ErrorBox.Render("Failed to load: " + m.err.Error() + ". Press 'r' to retry")
}
```

❌ **Bad**:
```go
if m.err != nil {
    return styles.ErrorBox.Render("Error: " + m.err.Error())
}
```

### 4. Helper Methods

✅ **Good**: Multiple small helper methods
- `renderHeader()`
- `renderContent()`
- `renderFooter()`
- `renderItem()`
- `renderError()`

❌ **Bad**: Single monolithic View() method

---

## Testing View Methods

### Basic View Test

```go
ginkgo.It("should render without errors", func() {
    model := NewYourModel()
    view := model.View()
    gomega.Expect(view).NotTo(gomega.BeEmpty())
})
```

### Structure Test

```go
ginkgo.It("should include header, content, and footer", func() {
    model := NewYourModel()
    view := model.View()
    gomega.Expect(view).To(gomega.ContainSubstring("Title"))
    gomega.Expect(view).To(gomega.ContainSubstring("Content"))
})
```

### Error Display Test

```go
ginkgo.It("should display errors correctly", func() {
    model := NewYourModel()
    model.err = fmt.Errorf("test error")
    view := model.View()
    gomega.Expect(view).To(gomega.ContainSubstring("Error"))
    gomega.Expect(view).To(gomega.ContainSubstring("retry"))
})
```

---

## Checklist for New Views

- [ ] View() method is < 20 lines
- [ ] Header is rendered via renderHeader()
- [ ] Content is rendered via renderContent()
- [ ] Footer is rendered via renderFooter()
- [ ] WindowSizeMsg updates width for all components
- [ ] Errors use ErrorBox with recovery guidance
- [ ] Colors use styles.* constants (no hex)
- [ ] Spacing is consistent (empty strings between sections)
- [ ] ScreenContainer used for padding
- [ ] Tests verify structure and error display
- [ ] Documentation in comments explains purpose

---

## Migration Guide

### From Manual Rendering to ScreenContainer

**Before**:
```go
func (m *Model) View() string {
    var content []string
    content = append(content, m.header.View())
    content = append(content, m.renderContent())
    content = append(content, m.footer.View())
    return strings.Join(content, "\n")
}
```

**After**:
```go
func (m *Model) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        "",
        m.renderContent(),
        "",
        m.renderFooter(),
    )
}

func (m *Model) renderContent() string {
    screenContainer := components.NewScreenContainer(content).
        WithPaddingMode(components.PaddingNormal)
    return screenContainer.Render()
}
```

---

## Common Patterns

### Pagination Format

```go
statusText := fmt.Sprintf("Showing %d-%d of %d items",
    startIdx, endIdx, totalItems)
content = append(content, styles.InputHint.Render(statusText))
```

### Selection Indicator

```go
if i == m.selectedIdx {
    item = "▶ " + styles.ListItemSelected.Render(item)
} else {
    item = "  " + styles.ListItem.Render(item)
}
```

### Empty State

```go
if len(m.items) == 0 {
    return styles.InfoText.Render("No items to display. Create one to get started.")
}
```

### Loading State

```go
if m.loading {
    return styles.InfoText.Render("Loading... Please wait.")
}
```

---

## Summary

The KaRiya CLI uses 4 main View() patterns:

1. **ScreenContainer** - Detail/review screens
2. **ListContainer** - Multiple item displays
3. **FormFieldContainer** - Input forms
4. **Modal** - Dialogs and overlays

All patterns follow:
- Separation of concerns (helper methods)
- Consistent error handling (ErrorBox with guidance)
- Responsive width management
- Style constant usage (no inline colors)
- Clear documentation and testing

---

## Related Documents

- [ERROR_HANDLING_GUIDE.md](./ERROR_HANDLING_GUIDE.md)
- [STYLE_USAGE_GUIDE.md](./STYLE_USAGE_GUIDE.md)
- [FOCUS_INDICATOR_GUIDE.md](./FOCUS_INDICATOR_GUIDE.md)


