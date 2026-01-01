# Task 4.6 - Error Display Standardization Audit

**Date**: 2026-01-01
**Status**: Audit Complete
**Priority**: HIGH

## Executive Summary

This audit examines error display patterns across all 8 refactored models in the KaRiya project to identify inconsistencies and establish a standardized error handling approach.

**Key Findings:**
- ✅ All models have error handling infrastructure
- ⚠️ Inconsistent error display patterns across models
- ⚠️ Inconsistent use of error styling (ErrorText vs ErrorBox)
- ⚠️ Field-level errors not consistently using FormFieldContainer
- ✅ Error style constants (ColorError, ErrorText, ErrorBox) are properly defined

## Detailed Audit Results

### 1. Form Model (`internal/cli/models/form.go`)

**Error Handling:**
- ✅ Field-level errors: `fieldErrors map[FormField]string`
- ✅ Model-level errors: `err error`
- ✅ Uses FormFieldContainer for field error display

**Error Display Pattern:**
```go
// Field-level errors (CORRECT)
if err, ok := m.fieldErrors[TextField]; ok {
    container.SetError(err)  // ✅ Using FormFieldContainer
}

// Model-level errors (INCONSISTENT)
if m.err != nil {
    content = append(content, styles.ErrorText.Render(m.err.Error()))  // ⚠️ Using ErrorText
}
```

**Issues:**
- Field errors: ✅ Correct (using FormFieldContainer)
- Model errors: ⚠️ Should use ErrorBox for consistency
- Error messages: Clear and descriptive

**Recommendation:**
- Change model-level error rendering from `ErrorText` to `ErrorBox`
- Wrap error in ScreenContainer for consistent padding

---

### 2. List Model (`internal/cli/models/list.go`)

**Error Handling:**
- ✅ Model-level errors: `err error`
- ❌ No field-level errors (not applicable)

**Error Display Pattern:**
```go
// Model-level error (INCONSISTENT)
if m.err != nil {
    content = append(content, styles.ErrorText.Render(fmt.Sprintf("Error loading events: %v", m.err)))
}
```

**Issues:**
- ⚠️ Uses ErrorText instead of ErrorBox
- Error message: "Error loading events: %v" - generic format
- No error recovery guidance

**Recommendation:**
- Change to ErrorBox for consistency with other models
- Add helpful recovery message (e.g., "Press 'r' to retry")
- Wrap in ScreenContainer for consistent padding

---

### 3. Details Model (`internal/cli/models/details.go`)

**Error Handling:**
- ✅ Uses ErrorBox for state errors

**Error Display Pattern:**
```go
// State error (CORRECT)
return styles.ErrorBox.Render("No event selected\n\nPress 'esc' to return")
```

**Issues:**
- ✅ Correct error style (ErrorBox)
- ✅ Includes recovery guidance

**Recommendation:**
- Already compliant - use as reference model

---

### 4. Confirmation Dialog Model (`internal/cli/models/confirmation_dialog.go`)

**Error Handling:**
- ❌ Minimal error handling
- ⚠️ Uses ColorError for styling but not in error display

**Error Display Pattern:**
```go
// Only found in button styling
Foreground(styles.ColorError)  // ⚠️ For destructive button, not errors
```

**Issues:**
- ❌ No error display mechanism for dialog operations
- ❌ No field-level error handling

**Recommendation:**
- Add error field to model: `err error`
- Implement error display in View() using ErrorBox
- Add error handling in Update() for dialog operations

---

### 5. Fact Editor Model (`internal/cli/models/fact_editor.go`)

**Error Handling:**
- ✅ Field-level errors: `fieldErrors map[int]string`
- ✅ Model-level errors: `err error`
- ✅ Uses FormFieldContainer for field error display

**Error Display Pattern:**
```go
// Field-level errors (CORRECT)
if err, ok := m.fieldErrors[FactTextFieldIdx]; ok {
    container.SetError(fieldErr)  // ✅ Using FormFieldContainer
}

// Model-level errors (INCONSISTENT)
if m.err != nil {
    content = append(content, styles.ErrorBox.Render(m.err.Error()))  // ✅ Using ErrorBox
}
```

**Issues:**
- ✅ Field errors: Correct (using FormFieldContainer)
- ✅ Model errors: Correct (using ErrorBox)
- Error messages: Clear and specific

**Recommendation:**
- Already compliant - use as reference model
- Ensure consistent error message format

---

### 6. Metadata Editor Model (`internal/cli/models/metadata_editor.go`)

**Error Handling:**
- ✅ Field-level errors: `fieldErrors map[int]string`
- ✅ Model-level errors: `err error`
- ✅ Uses FormFieldContainer for field error display

**Error Display Pattern:**
```go
// Field-level errors (CORRECT)
if err, ok := m.fieldErrors[MetadataDateFieldIdx]; ok {
    container.SetError(fieldErr)  // ✅ Using FormFieldContainer
}

// Model-level errors (INCONSISTENT - VERBOSE)
if m.err != nil {
    content = append(content, styles.ErrorBox.Render("Error: "+m.err.Error()))
}

// Field errors also rendered as ErrorBox (REDUNDANT)
for fieldIdx, errMsg := range m.fieldErrors {
    content = append(content, styles.ErrorBox.Render(fmt.Sprintf("Field %d: %s", fieldIdx, errMsg)))
}
```

**Issues:**
- ⚠️ Field errors rendered twice (once in FormFieldContainer, once as ErrorBox)
- ⚠️ Redundant error display causes visual clutter
- ⚠️ "Error: " prefix adds unnecessary verbosity

**Recommendation:**
- Remove the loop that renders fieldErrors as ErrorBox (lines 340-341)
- Keep only FormFieldContainer rendering for field errors
- Change "Error: " prefix to just the error message
- Use consistent error message format

---

### 7. Fact List Model (`internal/cli/models/fact_list.go`)

**Error Handling:**
- ✅ Model-level errors: `err error`
- ❌ No visible error display in View()
- ✅ Has GetError() method but error not rendered

**Error Display Pattern:**
```go
// Has GetError() method but no rendering in View()
func (flm *FactListModel) GetError() error {
    return flm.err
}
```

**Issues:**
- ❌ Errors are stored but not displayed to user
- ❌ No error feedback in UI
- User won't see error messages

**Recommendation:**
- Add error display in View() method using ErrorBox
- Display error message with recovery guidance
- Wrap in ScreenContainer for consistent padding

---

### 8. Burst List Model (`internal/cli/models/burst_list.go`)

**Error Handling:**
- ✅ Model-level errors in messages: `Err: fmt.Errorf(...)`
- ❌ No visible error display in View()
- ❌ No error field in model struct

**Error Display Pattern:**
```go
// Error only in message type, not rendered
return BurstsLoadedMsg{Bursts: []*career.Burst{}, Err: fmt.Errorf("...")}
```

**Issues:**
- ❌ Errors in messages but not stored/displayed in model
- ❌ No error display in View()
- ❌ No error recovery mechanism
- User won't see error messages

**Recommendation:**
- Add error field to model: `err error`
- Handle error in Update() and store in model
- Add error display in View() using ErrorBox
- Wrap in ScreenContainer for consistent padding

---

## Standards Established

### Field-Level Error Display (STANDARD)

**Pattern:** Use FormFieldContainer with SetError()

```go
// In View() method
if err, ok := m.fieldErrors[fieldIndex]; ok {
    container := components.NewFormFieldContainer()
    container.SetLabel("Field Name")
    container.SetInput(inputView)
    container.SetError(err)  // ✅ STANDARD
    content = append(content, container.Render())
}
```

**Style:** `styles.ErrorText` (applied within FormFieldContainer)
**Location:** Below input field
**Message Format:** Specific, actionable error message

---

### Model-Level Error Display (STANDARD)

**Pattern:** Use ErrorBox with recovery guidance

```go
// In View() method
if m.err != nil {
    errorMsg := fmt.Sprintf("%s\n\nPress 'r' to retry or 'esc' to cancel", m.err.Error())
    content = append(content, styles.ErrorBox.Render(errorMsg))
}
```

**Style:** `styles.ErrorBox`
**Color:** `styles.ColorError` (#d76e6e)
**Location:** Top of content area (after header, before main content)
**Message Format:** Error + recovery guidance
**Wrapping:** Optional ScreenContainer for consistent padding

---

### Error Message Format (STANDARD)

**Field-Level Errors:**
- Be specific: "Text exceeds 2000 characters limit"
- Be actionable: "Date cannot be in the future"
- Use present tense: "Company field is required"
- No "Error:" prefix (FormFieldContainer handles styling)

**Model-Level Errors:**
- Start with action: "Failed to save event"
- Include reason: "Database connection lost"
- Add recovery: "Press 'r' to retry"
- Example: `"Failed to save event: database connection lost\n\nPress 'r' to retry or 'esc' to cancel"`

---

## Compliance Matrix

| Model | Field Errors | Model Errors | Style | FormFieldContainer | ErrorBox | Status |
|-------|:---:|:---:|:---:|:---:|:---:|:---:|
| form.go | ✅ | ⚠️ | Mixed | ✅ | ❌ | NEEDS FIX |
| list.go | N/A | ⚠️ | ErrorText | N/A | ❌ | NEEDS FIX |
| details.go | N/A | ✅ | ErrorBox | N/A | ✅ | COMPLIANT |
| confirmation_dialog.go | ❌ | ❌ | None | ❌ | ❌ | NEEDS IMPL |
| fact_editor.go | ✅ | ✅ | Correct | ✅ | ✅ | COMPLIANT |
| metadata_editor.go | ✅ | ✅ | Redundant | ✅ | ⚠️ | NEEDS FIX |
| fact_list.go | N/A | ⚠️ | Missing | N/A | ❌ | NEEDS IMPL |
| burst_list.go | N/A | ❌ | Missing | N/A | ❌ | NEEDS IMPL |

**Legend:**
- ✅ Compliant/Implemented correctly
- ⚠️ Implemented but needs improvement
- ❌ Not implemented/Incorrect

---

## Remediation Plan

### Phase 1: High Priority (2 models)
1. **fact_list.go** - Implement error display
2. **burst_list.go** - Implement error display and storage

### Phase 2: Medium Priority (3 models)
3. **form.go** - Change model-level errors to ErrorBox
4. **list.go** - Change to ErrorBox and add recovery guidance
5. **metadata_editor.go** - Remove redundant error rendering

### Phase 3: Low Priority (1 model)
6. **confirmation_dialog.go** - Add error handling infrastructure

### Already Compliant (2 models)
- ✅ **details.go** - Use as reference
- ✅ **fact_editor.go** - Use as reference

---

## Implementation Checklist

- [ ] 4.6.2 Update form.go model-level errors to use ErrorBox
- [ ] 4.6.2 Update list.go model-level errors to use ErrorBox with recovery guidance
- [ ] 4.6.2 Update metadata_editor.go to remove redundant error rendering
- [ ] 4.6.3 Implement error display in fact_list.go using ErrorBox
- [ ] 4.6.3 Implement error storage and display in burst_list.go
- [ ] 4.6.3 Add error handling infrastructure to confirmation_dialog.go
- [ ] 4.6.4 Create error handling guide with examples
- [ ] Verify all models follow standard error display patterns
- [ ] Run full test suite to ensure no regressions
- [ ] Update AGENTS.md with completion status

---

## Reference Models

### Fact Editor (REFERENCE for form models)
- ✅ Field errors: FormFieldContainer.SetError()
- ✅ Model errors: ErrorBox with clear message
- ✅ Error recovery: Implicit in validation

### Details (REFERENCE for state errors)
- ✅ Uses ErrorBox for all errors
- ✅ Includes recovery guidance ("Press 'esc' to return")
- ✅ Clear, user-friendly messaging

---

## Document Information

- **Version**: 1.0
- **Created**: 2026-01-01
- **Status**: Audit Complete - Ready for Implementation
- **Related Task**: tasks/tasks-07-model-consistency.md (Task 4.6)
- **Next Steps**: Execute remediation plan (Task 4.6.2-4.6.4)

