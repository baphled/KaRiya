# Form Submission Issue - Detailed Analysis

## Problem Statement
When users try to submit the form, the data is not being saved to the database.

## Root Cause Analysis

### Issue 1: Form Submission Not Triggered by Ctrl+S
**Location**: `internal/cli/intents/capture_event_intent.go` - `updateCaptureForm()` method

**Problem**:
- The form submission is only triggered when:
  1. User navigates to the "Submit" button (using Tab)
  2. User presses Enter on the Submit button
- The form does NOT handle Ctrl+S to trigger submission
- The intent's `updateCaptureForm()` method does NOT handle Ctrl+S either

**Current Code** (lines 214-267):
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	// Delegate all messages to the form model to handle input and state
	_, formCmd := i.state.captureForm.Update(msg)

	// Check for special messages that indicate form completion or navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// User cancelled
			i.setCancelled()
			return nil

		case "esc":
			// Go back to strategy selection
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}
		// ← NO HANDLING FOR "ctrl+s"!
	}
	// ... rest of method
}
```

**Impact**:
- Users who press Ctrl+S to submit get no response
- The form appears to be unresponsive
- Users must navigate to the Submit button with Tab and press Enter (not intuitive)

### Issue 2: Form Submission Workflow Unclear
**Location**: Multiple files

**Problem**:
The submission workflow is complex and not well-documented:

1. **Form Submission Flow**:
   - User presses Enter on Submit button
   - Form.Update() calls m.submitForm()
   - submitForm() returns a tea.Cmd
   - The Cmd emits a SubmitMsg asynchronously
   - Intent receives SubmitMsg in next Update() call
   - Intent validates and transitions to review state

2. **Issues with this flow**:
   - Not obvious to users that they need to press Tab to reach the button
   - No visual indication of which button is focused
   - No keyboard shortcut (Ctrl+S) as users would expect
   - Form doesn't validate/show errors until Submit is pressed

### Issue 3: Data Not Being Persisted
**Location**: `internal/cli/intents/capture_event_intent.go` - `performSubmit()` method

**Potential Issues**:
1. **Event validation might be failing**:
   ```go
   if err := i.state.reviewState.Event.Validate(); err != nil {
       i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
       return nil
   }
   ```
   - If event validation fails, submission stops
   - User might not see error message

2. **Service might not be initialized**:
   ```go
   if i.eventService == nil {
       return SubmitErrorMsg{
           Code:    "SERVICE_ERROR",
           Message: "Event service not initialized",
           Cause:   nil,
       }
   }
   ```
   - If service is nil, submission fails silently

3. **Service call might be failing**:
   ```go
   err := i.eventService.CaptureEvent(
       ctx,
       event.Text,
       event.Date,
       mode,
   )
   ```
   - If this returns an error, it's captured but might not be displayed to user

## Solution

### Fix 1: Add Ctrl+S Support to Form Submission
**Action**: Modify `updateCaptureForm()` to handle Ctrl+S

**Implementation**:
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	_, formCmd := i.state.captureForm.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			// Submit the form (Ctrl+S or focused submit button + Enter)
			// Trigger form submission by calling submitForm directly
			return i.triggerFormSubmission()

		case "q", "ctrl+c":
			i.setCancelled()
			return nil

		case "esc":
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}
	}

	return formCmd
}

func (i *CaptureEventIntent) triggerFormSubmission() tea.Cmd {
	// Call the form's submitForm method
	return i.state.captureForm.SubmitForm()
}
```

### Fix 2: Add Explicit SubmitForm Method to FormModel
**Action**: Expose submitForm as a public method in FormModel

**Implementation**:
```go
// SubmitForm triggers form submission
// This is called by the intent when user presses Ctrl+S or Submit button
func (m *FormModel) SubmitForm() tea.Cmd {
	return m.submitForm()
}
```

### Fix 3: Add Better Error Feedback
**Action**: Ensure errors are displayed to the user

**Implementation**:
- When submission fails, show error state
- Display error message in a visible way
- Allow user to retry or go back

### Fix 4: Add Logging for Debugging
**Action**: Add structured logging to trace submission

**Implementation**:
```go
func (i *CaptureEventIntent) performSubmit() tea.Cmd {
	return func() tea.Msg {
		// Log submission attempt
		log.Info("Submitting event", "text", event.Text, "date", event.Date)

		// ... validation and service call ...

		if err != nil {
			log.Error("Event submission failed", "error", err)
			return SubmitErrorMsg{...}
		}

		log.Info("Event submitted successfully")
		return SubmitCompleteMsg{}
	}
}
```

## Testing Strategy

### Unit Tests
1. Test that Ctrl+S triggers form submission
2. Test that form validation works correctly
3. Test that submission command is returned properly

### Integration Tests
1. Test complete workflow: Input → Submit → Review → Confirm → Save
2. Test error scenarios: Invalid data, service failures
3. Test user interactions: Ctrl+S, Tab, Enter

### Manual Testing
1. Type event details
2. Press Ctrl+S to submit
3. Verify form transitions to review state
4. Confirm review
5. Verify event is saved to database

## Files to Modify

1. **internal/cli/intents/capture_event_intent.go**
   - Add Ctrl+S handling to updateCaptureForm()
   - Add triggerFormSubmission() method
   - Add logging to performSubmit()

2. **internal/cli/models/form.go**
   - Expose submitForm as public SubmitForm() method
   - Add logging for debugging

3. **internal/cli/intents/capture_event.go** (if needed)
   - Add logging configuration

## Expected Outcome

After implementing these fixes:
1. ✅ Users can press Ctrl+S to submit the form
2. ✅ Form transitions to review state on successful submission
3. ✅ Event data is validated before submission
4. ✅ Errors are displayed to the user
5. ✅ Event is saved to database
6. ✅ Users can see what went wrong if submission fails

## Verification Checklist

- [ ] Form submission triggered by Ctrl+S
- [ ] Form submission triggered by Tab + Enter
- [ ] Form transitions to review state
- [ ] Event data is saved to database
- [ ] Validation errors are displayed
- [ ] Service errors are displayed
- [ ] User can retry failed submissions
- [ ] User can go back to form to edit
- [ ] All tests pass
- [ ] No race conditions detected

