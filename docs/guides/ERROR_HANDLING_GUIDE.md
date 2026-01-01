# Error Handling Guide for KaRiya Models

**Document Version**: 1.0
**Created**: 2026-01-01
**Status**: Complete
**Audience**: Developers working on KaRiya models

## Table of Contents

1. [Overview](#overview)
2. [Error Display Standards](#error-display-standards)
3. [Field-Level Errors](#field-level-errors)
4. [Model-Level Errors](#model-level-errors)
5. [Implementation Examples](#implementation-examples)
6. [Best Practices](#best-practices)
7. [Anti-Patterns](#anti-patterns)
8. [Testing Error Display](#testing-error-display)

---

## Overview

KaRiya follows a consistent error handling pattern across all models to ensure:
- **Consistent User Experience**: Users see errors in the same style and location
- **Clear Error Messages**: Errors are specific and actionable
- **Proper Recovery Guidance**: Users know how to resolve errors
- **Visual Hierarchy**: Errors are visually distinct but not overwhelming

### Error Types in KaRiya

1. **Field-Level Errors**: Validation errors for specific form fields
2. **Model-Level Errors**: Errors affecting the entire screen/operation
3. **Recovery Errors**: Errors during data operations (save, load, etc.)

---

## Error Display Standards

### Color System

All errors use the standardized error color from the style system:

```go
// From styles/styles.go
ColorError = lipgloss.Color("#d76e6e") // Error red (muted)
```

**Never use inline hex colors for errors.** Always use `styles.ColorError`.

### Error Styles

KaRiya provides two main error styling options:

#### 1. ErrorText (for field-level errors)

```go
// From styles/styles.go
ErrorText = lipgloss.NewStyle().
    Foreground(ColorError).
    Bold(false)
```

**Usage**: Field validation errors within FormFieldContainer
**Location**: Below the input field
**Example**: "Text exceeds 2000 characters limit"

#### 2. ErrorBox (for model-level errors)

```go
// From styles/styles.go
ErrorBox = lipgloss.NewStyle().
    Background(ColorErrorBg).
    Foreground(ColorError).
    Padding(1)
```

**Usage**: Errors affecting the entire model/screen
**Location**: Top of content area (after header, before main content)
**Example**: "Error loading events: database connection lost"

---

## Field-Level Errors

### Standard Implementation

Field-level errors are displayed using `FormFieldContainer.SetError()`:

```go
// In your model's View() method
func (m *FormModel) renderTextField() string {
    focused := m.focusIndex == TextFieldIdx
    fieldErr := ""

    // Check for validation error
    if err, ok := m.fieldErrors[TextFieldIdx]; ok {
        fieldErr = err
    }

    // Create container with error
    container := components.NewFormFieldContainer()
    container.SetLabel("Event Description")
    container.SetInput(m.inputs[TextFieldIdx].View())

    // SetError handles styling automatically
    if fieldErr != "" {
        container.SetError(fieldErr)
    }

    return container.Render()
}
```

### Field Error Message Format

Field errors should be:
- **Specific**: Describe the exact problem
- **Actionable**: Tell the user how to fix it
- **Concise**: Keep to one line when possible
- **Lowercase**: Start with lowercase (no "Error:" prefix)

**Good Examples:**
```
✓ "Text exceeds 2000 characters limit"
✓ "Date cannot be in the future"
✓ "Company field is required"
✓ "Invalid date format (expected YYYY-MM-DD)"
```

**Bad Examples:**
```
✗ "Error: text exceeds limit"        (verbose prefix)
✗ "Something went wrong"              (not specific)
✗ "INVALID INPUT"                     (too aggressive)
✗ "Error in field 0"                  (not helpful)
```

### Storing Field Errors

Store field errors in a map:

```go
type YourModel struct {
    fieldErrors map[int]string  // Map field index to error message
}

// In Update() method
func (m *YourModel) validateFields() {
    // Clear previous errors
    m.fieldErrors = make(map[int]string)

    // Validate each field
    if len(m.text) == 0 {
        m.fieldErrors[TextFieldIdx] = "Text is required"
    }

    if len(m.text) > 2000 {
        m.fieldErrors[TextFieldIdx] = "Text exceeds 2000 characters"
    }
}
```

### Clearing Field Errors

Clear field errors when:
1. User edits the field (in Update())
2. Form is submitted successfully
3. Form is reset/cancelled

```go
func (m *YourModel) clearFieldErrors() {
    m.fieldErrors = make(map[int]string)
}
```

---

## Model-Level Errors

### Standard Implementation

Model-level errors are displayed using `ErrorBox`:

```go
// In your model struct
type YourModel struct {
    err error  // Model-level error
}

// In View() method
func (m *YourModel) View() string {
    var content []string

    // Display error if present
    if m.err != nil {
        errorMsg := fmt.Sprintf(
            "Error loading data: %v\n\nPress 'r' to retry or 'esc' to cancel",
            m.err,
        )
        content = append(content, styles.ErrorBox.Render(errorMsg))
    }

    // ... rest of content
    return strings.Join(content, "\n\n")
}
```

### Model Error Message Format

Model errors should:
- **Start with action**: "Failed to save", "Error loading", etc.
- **Include reason**: "database connection lost", "invalid response"
- **Add recovery**: "Press 'r' to retry or 'esc' to cancel"
- **Use present tense**: "Connection lost" not "Connection was lost"

**Good Examples:**
```
✓ "Failed to save event: database connection lost\n\nPress 'r' to retry or 'esc' to cancel"
✓ "Error loading facts: invalid response from server\n\nPress 'r' to retry"
✓ "Burst repository not configured\n\nContact administrator"
```

**Bad Examples:**
```
✗ "Error"                           (not specific)
✗ "Something went wrong"            (not actionable)
✗ "ERROR: Database Error!"          (too aggressive)
✗ "Failed to save"                  (no recovery guidance)
```

### Storing Model Errors

Store model-level errors in an error field:

```go
// In Update() method
case SubmitMsg:
    if msg.Err != nil {
        m.err = msg.Err
        return m, nil
    }
    m.err = nil  // Clear error on success
    // ... handle success
```

### Clearing Model Errors

Clear model errors when:
1. User retries the operation
2. Form is submitted successfully
3. User navigates away
4. New operation starts

```go
func (m *YourModel) clearError() {
    m.err = nil
}
```

---

## Implementation Examples

### Example 1: Form with Field Errors

```go
type ContactFormModel struct {
    nameInput   textinput.Model
    emailInput  textinput.Model
    fieldErrors map[int]string  // Field index to error
    err         error            // Model-level error
    focusIndex  int
}

func (m *ContactFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            return m, m.submitForm()
        }

    case SubmitMsg:
        if msg.Err != nil {
            m.err = msg.Err
            return m, nil
        }
        // Success
        m.err = nil
        return m, nil
    }
    return m, nil
}

func (m *ContactFormModel) submitForm() tea.Cmd {
    return func() tea.Msg {
        // Validate fields
        m.fieldErrors = make(map[int]string)

        if len(strings.TrimSpace(m.nameInput.Value())) == 0 {
            m.fieldErrors[0] = "Name is required"
        }

        if len(strings.TrimSpace(m.emailInput.Value())) == 0 {
            m.fieldErrors[1] = "Email is required"
        }

        if len(m.fieldErrors) > 0 {
            return nil  // Don't submit if validation fails
        }

        // Attempt to save
        err := m.service.SaveContact(m.nameInput.Value(), m.emailInput.Value())
        return SubmitMsg{Err: err}
    }
}

func (m *ContactFormModel) View() string {
    var content []string

    // Display model-level error
    if m.err != nil {
        errorMsg := fmt.Sprintf("Failed to save: %v\n\nPress 'r' to retry", m.err)
        content = append(content, styles.ErrorBox.Render(errorMsg))
    }

    // Render name field with error
    nameContainer := components.NewFormFieldContainer()
    nameContainer.SetLabel("Name")
    nameContainer.SetInput(m.nameInput.View())
    if err, ok := m.fieldErrors[0]; ok {
        nameContainer.SetError(err)
    }
    content = append(content, nameContainer.Render())

    // Render email field with error
    emailContainer := components.NewFormFieldContainer()
    emailContainer.SetLabel("Email")
    emailContainer.SetInput(m.emailInput.View())
    if err, ok := m.fieldErrors[1]; ok {
        emailContainer.SetError(err)
    }
    content = append(content, emailContainer.Render())

    return strings.Join(content, "\n\n")
}
```

### Example 2: List with Load Error

```go
type EventListModel struct {
    events []*career.CareerEvent
    err    error
    width  int
    height int
}

func (m *EventListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "r" && m.err != nil {
            // Retry loading
            m.err = nil
            return m, m.loadEvents()
        }

    case EventsLoadedMsg:
        if msg.Err != nil {
            m.err = msg.Err
            return m, nil
        }
        m.events = msg.Events
        m.err = nil
        return m, nil
    }
    return m, nil
}

func (m *EventListModel) View() string {
    var content []string

    // Display error if present
    if m.err != nil {
        errorMsg := fmt.Sprintf(
            "Error loading events: %v\n\nPress 'r' to retry or 'esc' to cancel",
            m.err,
        )
        content = append(content, styles.ErrorBox.Render(errorMsg))
    } else if len(m.events) == 0 {
        content = append(content, "No events found")
    } else {
        // Render events
        for _, event := range m.events {
            content = append(content, fmt.Sprintf("▶ %s", event.Text))
        }
    }

    return strings.Join(content, "\n\n")
}
```

---

## Best Practices

### 1. Always Provide Recovery Guidance

Users should know how to recover from errors:

```go
// Good: Provides recovery path
errorMsg := fmt.Sprintf(
    "Failed to save: %v\n\nPress 'r' to retry or 'esc' to cancel",
    err,
)

// Bad: No recovery guidance
errorMsg := fmt.Sprintf("Error: %v", err)
```

### 2. Validate Early, Display Clearly

Validate fields as soon as possible and display errors immediately:

```go
// In Update() when field value changes
case tea.KeyMsg:
    // Update field
    m.inputs[idx].SetValue(newValue)

    // Validate immediately
    if err := m.validateField(idx); err != nil {
        m.fieldErrors[idx] = err.Error()
    } else {
        delete(m.fieldErrors, idx)
    }
```

### 3. Use Consistent Error Messages

Reuse error messages across the application:

```go
// Define error messages as constants
const (
    ErrRequired    = "This field is required"
    ErrTooLong     = "Text exceeds maximum length"
    ErrInvalidDate = "Invalid date format (expected YYYY-MM-DD)"
)

// Use in validation
if len(value) == 0 {
    m.fieldErrors[idx] = ErrRequired
}
```

### 4. Clear Errors Appropriately

Clear errors when they're no longer relevant:

```go
// Clear field error when user edits the field
case tea.KeyMsg:
    m.inputs[idx].SetValue(newValue)
    delete(m.fieldErrors, idx)  // Clear error

    // Re-validate
    if err := m.validateField(idx); err != nil {
        m.fieldErrors[idx] = err.Error()
    }
```

### 5. Distinguish Error Types

Use different styling for different error severities:

```go
// Critical error - use ErrorBox
if criticalError {
    return styles.ErrorBox.Render(fmt.Sprintf(
        "Critical error: %s\n\nPlease restart the application",
        err,
    ))
}

// Field error - use FormFieldContainer
if fieldError {
    container.SetError(err.Error())
}

// Warning - use different style
if warning {
    return styles.WarningBox.Render(fmt.Sprintf(
        "Warning: %s",
        msg,
    ))
}
```

---

## Anti-Patterns

### ❌ Don't: Use Inline Hex Colors

```go
// BAD
errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d76e6e"))
```

### ✅ Do: Use Style Constants

```go
// GOOD
content = append(content, styles.ErrorBox.Render(errorMsg))
```

---

### ❌ Don't: Render Field Errors Twice

```go
// BAD - Renders error twice
if err, ok := m.fieldErrors[idx]; ok {
    container.SetError(err)  // First render
}
// ... later ...
if err, ok := m.fieldErrors[idx]; ok {
    content = append(content, styles.ErrorBox.Render(err))  // Second render
}
```

### ✅ Do: Use FormFieldContainer for Field Errors

```go
// GOOD - Single render via container
container := components.NewFormFieldContainer()
container.SetLabel("Field")
container.SetInput(input)
if err, ok := m.fieldErrors[idx]; ok {
    container.SetError(err)
}
content = append(content, container.Render())
```

---

### ❌ Don't: Use Generic Error Messages

```go
// BAD
m.err = fmt.Errorf("error")
m.err = fmt.Errorf("something went wrong")
m.err = fmt.Errorf("failed")
```

### ✅ Do: Be Specific

```go
// GOOD
m.err = fmt.Errorf("failed to connect to database: %w", err)
m.err = fmt.Errorf("invalid response from API: status %d", status)
m.err = fmt.Errorf("file not found: %s", filename)
```

---

### ❌ Don't: Forget Recovery Guidance

```go
// BAD - User doesn't know what to do
errorMsg := fmt.Sprintf("Error loading data: %v", m.err)
content = append(content, styles.ErrorBox.Render(errorMsg))
```

### ✅ Do: Provide Recovery Steps

```go
// GOOD - User knows how to recover
errorMsg := fmt.Sprintf(
    "Error loading data: %v\n\nPress 'r' to retry or 'esc' to cancel",
    m.err,
)
content = append(content, styles.ErrorBox.Render(errorMsg))
```

---

### ❌ Don't: Mix Error Styles

```go
// BAD - Inconsistent styling
if m.err != nil {
    content = append(content, styles.ErrorText.Render(m.err.Error()))  // Field style for model error
}
```

### ✅ Do: Use Correct Style for Error Type

```go
// GOOD - Correct styling
if m.err != nil {
    errorMsg := fmt.Sprintf("Error: %v\n\nPress 'r' to retry", m.err)
    content = append(content, styles.ErrorBox.Render(errorMsg))  // Model style
}
```

---

## Testing Error Display

### Unit Tests for Field Errors

```go
ginkgo.Describe("ContactForm", func() {
    ginkgo.It("should display field error for missing name", func() {
        m := NewContactForm()
        m.fieldErrors[0] = "Name is required"

        view := m.View()
        gomega.Expect(view).To(gomega.ContainSubstring("Name is required"))
    })

    ginkgo.It("should clear field error after user edits", func() {
        m := NewContactForm()
        m.fieldErrors[0] = "Name is required"

        // Simulate user edit
        delete(m.fieldErrors, 0)

        view := m.View()
        gomega.Expect(view).NotTo(gomega.ContainSubstring("Name is required"))
    })
})
```

### Unit Tests for Model Errors

```go
ginkgo.Describe("EventList", func() {
    ginkgo.It("should display model error when loading fails", func() {
        m := NewEventList()
        m.err = fmt.Errorf("database connection lost")

        view := m.View()
        gomega.Expect(view).To(gomega.ContainSubstring("Error loading"))
        gomega.Expect(view).To(gomega.ContainSubstring("database connection lost"))
        gomega.Expect(view).To(gomega.ContainSubstring("Press 'r' to retry"))
    })

    ginkgo.It("should clear model error after successful retry", func() {
        m := NewEventList()
        m.err = fmt.Errorf("database connection lost")

        // Simulate successful retry
        m.err = nil
        m.events = []*Event{...}

        view := m.View()
        gomega.Expect(view).NotTo(gomega.ContainSubstring("Error loading"))
    })
})
```

### Integration Tests

```go
ginkgo.Describe("EventListIntegration", func() {
    ginkgo.It("should display error when service fails", func() {
        // Create model with failing service
        service := &MockService{
            ListErr: fmt.Errorf("network timeout"),
        }
        m := NewEventList(service)

        // Trigger load
        _, _ = m.Update(tea.KeyMsg{})

        // Verify error display
        view := m.View()
        gomega.Expect(view).To(gomega.ContainSubstring("Error loading"))
        gomega.Expect(view).To(gomega.ContainSubstring("network timeout"))
    })
})
```

---

## Reference Models

### Form Model (Reference Implementation)

**File**: `internal/cli/models/form.go`
**Features**:
- ✅ Field-level errors using FormFieldContainer
- ✅ Model-level errors using ErrorBox
- ✅ Proper error clearing
- ✅ Validation error messages

### Fact Editor Model (Reference Implementation)

**File**: `internal/cli/models/fact_editor.go`
**Features**:
- ✅ Field-level errors with specific messages
- ✅ Model-level errors with recovery guidance
- ✅ Error clearing on successful submission

### Details Model (Reference Implementation)

**File**: `internal/cli/models/details.go`
**Features**:
- ✅ State error display using ErrorBox
- ✅ Recovery guidance included
- ✅ Clear, user-friendly messages

---

## Compliance Checklist

When implementing error display in a new model, verify:

- [ ] Field errors use `FormFieldContainer.SetError()`
- [ ] Model errors use `styles.ErrorBox.Render()`
- [ ] Error messages are specific and actionable
- [ ] Recovery guidance is provided
- [ ] No inline hex colors used
- [ ] All error styles use `ColorError` constant
- [ ] Field errors cleared appropriately
- [ ] Model errors cleared on success
- [ ] Error messages are in present tense
- [ ] Tests verify error display
- [ ] Tests verify error clearing

---

## Related Documentation

- [Style System Guide](./STYLE_USAGE_GUIDE.md)
- [Component Usage Guide](./COMPONENT_USAGE_GUIDE.md)
- [Model Development Guide](./MODEL_DEVELOPMENT_GUIDE.md)
- [Task 4.6 Audit Report](../audits/TASK_4.6_ERROR_DISPLAY_AUDIT.md)

---

## Questions?

For questions about error handling:
1. Check this guide first
2. Review reference models (form.go, fact_editor.go)
3. Check existing tests for patterns
4. Ask in the development team

---

**Document Information**

- **Version**: 1.0
- **Created**: 2026-01-01
- **Last Updated**: 2026-01-01
- **Status**: Complete and Ready for Use
- **Related Task**: tasks/tasks-07-model-consistency.md (Task 4.6)

