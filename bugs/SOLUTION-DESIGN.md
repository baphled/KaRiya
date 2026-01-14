# Bug 001: Escape Key Navigation - Solution Design

**Date**: 2026-01-12  
**Status**: Design Complete, Implementation Ready  
**Approach**: Middleware Pattern with MessageInterceptor

---

## Problem Summary

Forms and modals consume escape key events before `HandleGlobalKeys()` can process them, trapping users in screens.

**Root Cause**: Message delegation happens BEFORE global key checking.

---

## Solution: MessageInterceptor Middleware

Create a reusable middleware layer that intercepts messages and handles global keys BEFORE delegation to sub-components.

### Architecture

```
User presses key
  → BubbleTea delivers tea.KeyMsg
  → Intent.Update(msg)
  → MessageInterceptor.InterceptOr(msg, fallback)
      → Check for global keys (esc, q, ?)
      → IF global key matched:
          → Call appropriate handler (OnBack, OnQuit, OnHelp)
          → Return command
      → IF no match:
          → Call fallback function (delegates to form/modal)
          → Return command
```

### Key Benefits

✅ **Consistent**: Same pattern across all intents  
✅ **Declarative**: Clear, readable intent code  
✅ **Type-Safe**: No manual switch statements  
✅ **Testable**: Easy to test interception logic  
✅ **Maintainable**: One place to update global key handling  
✅ **Reusable**: Works for forms, modals, any sub-component  

---

## Implementation

### Step 1: MessageInterceptor Helper (✅ DONE)

**File**: `internal/cli/intents/view_helpers.go` (after line 444)

```go
// MessageInterceptor provides a middleware layer for handling global keys before delegation.
// This ensures escape, quit, and other global keys are always processed first,
// preventing sub-components (forms, modals) from consuming them.
type MessageInterceptor struct {
	backHandler GlobalKeyHandler
	quitHandler GlobalKeyHandler
	helpHandler GlobalKeyHandler
}

// GlobalKeyHandler is a function that handles a global key event.
type GlobalKeyHandler func() tea.Cmd

// NewMessageInterceptor creates a new message interceptor with no handlers.
func NewMessageInterceptor() *MessageInterceptor {
	return &MessageInterceptor{}
}

// OnBack sets the handler for escape key (back navigation).
func (m *MessageInterceptor) OnBack(handler GlobalKeyHandler) *MessageInterceptor {
	m.backHandler = handler
	return m
}

// OnQuit sets the handler for quit key (q or Ctrl+C).
func (m *MessageInterceptor) OnQuit(handler GlobalKeyHandler) *MessageInterceptor {
	m.quitHandler = handler
	return m
}

// OnHelp sets the handler for help key (?).
func (m *MessageInterceptor) OnHelp(handler GlobalKeyHandler) *MessageInterceptor {
	m.helpHandler = handler
	return m
}

// InterceptOr checks for global keys and calls the appropriate handler.
// If no global key is matched, it calls the fallback function.
func (m *MessageInterceptor) InterceptOr(msg tea.Msg, fallback func() tea.Cmd) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return fallback()
	}

	result := HandleGlobalKeys(keyMsg)

	switch result {
	case KeyBack:
		if m.backHandler != nil {
			return m.backHandler()
		}
	case KeyQuit:
		if m.quitHandler != nil {
			return m.quitHandler()
		}
	case KeyHelp:
		if m.helpHandler != nil {
			return m.helpHandler()
		}
	}

	return fallback()
}

// Intercept is similar to InterceptOr but returns nil if no fallback is needed.
func (m *MessageInterceptor) Intercept(msg tea.Msg) tea.Cmd {
	return m.InterceptOr(msg, func() tea.Cmd { return nil })
}
```

**Status**: ✅ Implemented in `view_helpers.go`

---

### Step 2: Apply to CaptureEvent Intent

#### Fix 1: updateCaptureForm() - Form Delegation

**Before** (BROKEN):
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
    // ❌ Delegate to form FIRST
    _, formCmd := i.state.captureForm.Update(msg)

    // ❌ Check global keys AFTER (too late!)
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyBack:
            i.state.currentState = CaptureStateChooseStrategy
            return nil
        }
    }
    // ...
}
```

**After** (FIXED):
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
    // ✅ Use MessageInterceptor to check global keys FIRST
    return NewMessageInterceptor().
        OnBack(func() tea.Cmd {
            i.state.currentState = CaptureStateChooseStrategy
            return nil
        }).
        OnQuit(func() tea.Cmd {
            return tea.Quit
        }).
        OnHelp(func() tea.Cmd {
            i.ToggleHelp()
            return nil
        }).
        InterceptOr(msg, func() tea.Cmd {
            // ✅ Only called if no global key matched
            return i.handleFormMessages(msg)
        })
}

// handleFormMessages processes form-specific messages after global keys are handled.
func (i *CaptureEventIntent) handleFormMessages(msg tea.Msg) tea.Cmd {
    // Check for Ctrl+S
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if keyMsg.String() == "ctrl+s" {
            return i.state.captureForm.SubmitForm()
        }
    }

    // ✅ NOW delegate to form
    _, formCmd := i.state.captureForm.Update(msg)

    // Check for completion messages
    switch msg := msg.(type) {
    case models.SubmitMsg:
        // ... handle submission
    case FormSubmittedMsg:
        // ... handle legacy message
    }

    return formCmd
}
```

---

#### Fix 2-4: updateReviewInferredEvent() - Modal Delegation

**Before** (BROKEN):
```go
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
    // ❌ Route to modals immediately
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        if i.state.reviewState.metadataModal != nil {
            // ❌ Delegate to modal first (wrong!)
            modal, cmd := i.state.reviewState.metadataModal.Update(msg)
            // ...
            return cmd
        }
    // ... other modals
    }

    // Global keys checked LATER (too late!)
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
            // ...
        }
    }
}
```

**After** (FIXED):
```go
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
    // ✅ Use MessageInterceptor to check global keys FIRST
    return NewMessageInterceptor().
        OnBack(func() tea.Cmd {
            // If modal is active, close it
            if i.state.reviewState.EditingMode != EditingModeNone {
                i.state.reviewState.metadataModal = nil
                i.state.reviewState.burstModal = nil
                i.state.reviewState.factModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
                return nil
            }
            // Otherwise go back to form
            i.state.currentState = CaptureStateForm
            return nil
        }).
        OnQuit(func() tea.Cmd {
            return tea.Quit
        }).
        OnHelp(func() tea.Cmd {
            i.ToggleHelp()
            return nil
        }).
        InterceptOr(msg, func() tea.Cmd {
            // ✅ Only called if no global key matched
            return i.handleReviewMessages(msg)
        })
}

// handleReviewMessages processes review-specific messages after global keys are handled.
func (i *CaptureEventIntent) handleReviewMessages(msg tea.Msg) tea.Cmd {
    // ✅ NOW route to modals if active
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        if i.state.reviewState.metadataModal != nil {
            modal, cmd := i.state.reviewState.metadataModal.Update(msg)
            i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

            // Check completion
            if i.state.reviewState.metadataModal.IsSubmitted() {
                i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
                i.state.reviewState.metadataModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
            } else if i.state.reviewState.metadataModal.IsCancelled() {
                i.state.reviewState.metadataModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
            }
            return cmd
        }

    case EditingModeBursts:
        if i.state.reviewState.burstModal != nil {
            modal, cmd := i.state.reviewState.burstModal.Update(msg)
            i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)
            // Modal closed by interceptor (no manual escape check needed!)
            return cmd
        }

    case EditingModeFacts:
        if i.state.reviewState.factModal != nil {
            modal, cmd := i.state.reviewState.factModal.Update(msg)
            i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

            // Check completion
            if i.state.reviewState.factModal.IsSubmitted() {
                i.state.reviewState.factModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
            } else if i.state.reviewState.factModal.IsCancelled() {
                i.state.reviewState.factModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
            }
            return cmd
        }
    }

    // Normal review handling (no modal active)
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+s", "enter":
            i.state.currentState = CaptureStateSubmit
            return i.performSubmit()
        case "e":
            i.state.reviewState.EditingMode = EditingModeMetadata
            return nil
        case "b":
            i.state.reviewState.EditingMode = EditingModeBursts
            return nil
        case "f":
            i.state.reviewState.EditingMode = EditingModeFacts
            return nil
        // ... other keys
        }
    }

    return nil
}
```

---

## Benefits Over Manual Approach

### Before (Manual Pattern - VERBOSE)
```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyQuit:
            return tea.Quit
        case KeyHelp:
            i.ToggleHelp()
            return nil
        case KeyBack:
            i.state = previousState
            return nil
        }
    }

    _, cmd := i.subComponent.Update(msg)
    return cmd
}
```

### After (Interceptor Pattern - CLEAN)
```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    return NewMessageInterceptor().
        OnBack(func() tea.Cmd { i.state = previousState; return nil }).
        OnQuit(func() tea.Cmd { return tea.Quit }).
        OnHelp(func() tea.Cmd { i.ToggleHelp(); return nil }).
        InterceptOr(msg, func() tea.Cmd {
            _, cmd := i.subComponent.Update(msg)
            return cmd
        })
}
```

**Improvements**:
- ✅ 12 lines → 6 lines (50% reduction)
- ✅ No nested switches
- ✅ Declarative, readable
- ✅ Impossible to get order wrong
- ✅ Consistent across intents

---

## Testing Strategy

### Unit Tests for MessageInterceptor

**New File**: `internal/cli/intents/view_helpers_test.go`

```go
var _ = Describe("MessageInterceptor", func() {
    var interceptor *MessageInterceptor

    BeforeEach(func() {
        interceptor = NewMessageInterceptor()
    })

    Describe("OnBack", func() {
        It("should call back handler when escape pressed", func() {
            called := false
            interceptor.OnBack(func() tea.Cmd {
                called = true
                return nil
            })

            interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
            Expect(called).To(BeTrue())
        })
    })

    Describe("OnQuit", func() {
        It("should call quit handler when q pressed", func() {
            called := false
            interceptor.OnQuit(func() tea.Cmd {
                called = true
                return tea.Quit
            })

            cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
            Expect(called).To(BeTrue())
            Expect(cmd).To(Equal(tea.Quit()))
        })
    })

    Describe("InterceptOr", func() {
        It("should call fallback when no global key matched", func() {
            fallbackCalled := false
            interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, func() tea.Cmd {
                fallbackCalled = true
                return nil
            })

            Expect(fallbackCalled).To(BeTrue())
        })

        It("should NOT call fallback when global key matched", func() {
            fallbackCalled := false
            interceptor.
                OnBack(func() tea.Cmd { return nil }).
                InterceptOr(tea.KeyMsg{Type: tea.KeyEsc}, func() tea.Cmd {
                    fallbackCalled = true
                    return nil
                })

            Expect(fallbackCalled).To(BeFalse())
        })
    })
})
```

### Integration Tests for CaptureEvent

Update existing escape tests to verify interceptor works:

```go
Describe("CaptureEvent - Escape with MessageInterceptor", func() {
    It("should handle escape from form using interceptor", func() {
        intent.state.currentState = CaptureStateForm
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

        Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
        Expect(intent.active).To(BeTrue())
    })

    It("should handle escape from modal using interceptor", func() {
        intent.state.currentState = CaptureStateReview
        intent.state.reviewState.EditingMode = EditingModeMetadata
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

        Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeNone))
    })
})
```

---

## Rollout Plan

### Phase 1: Core Infrastructure ✅ DONE
- [x] Create MessageInterceptor in view_helpers.go
- [x] Add tests for interceptor
- [x] Document pattern

### Phase 2: Fix CaptureEvent Intent
- [ ] Apply to updateCaptureForm()
- [ ] Apply to updateReviewInferredEvent()
- [ ] Remove manual escape checks (line 420)
- [ ] Add handleFormMessages() helper
- [ ] Add handleReviewMessages() helper

### Phase 3: Testing
- [ ] Unit tests pass
- [ ] E2E tests pass
- [ ] Manual testing complete

### Phase 4: Documentation
- [ ] Create ESCAPE_KEY_HANDLING.md guide
- [ ] Update TUI_STANDARDS.md with interceptor pattern
- [ ] Add examples to view_helpers.go docs

### Phase 5: Adoption (Optional)
- [ ] Offer to other intents (not required, but beneficial)
- [ ] Create migration guide
- [ ] Add to intent template

---

## Timeline

- Phase 1: ✅ Complete (1 hour)
- Phase 2: 1-2 hours (implementation)
- Phase 3: 1 hour (testing)
- Phase 4: 30 min (documentation)
- **Total**: 3.5-4.5 hours

---

## Success Criteria

✅ **MessageInterceptor created and tested**  
⬜ Escape works from CaptureEvent form  
⬜ Escape works from all 3 modals  
⬜ All tests passing  
⬜ No regressions  
⬜ Code cleaner and more maintainable  
⬜ Pattern documented for future use  

---

**Status**: Phase 1 Complete, Ready for Phase 2 Implementation  
**Next Action**: Apply interceptor to CaptureEvent intent functions  
**Estimated Completion**: 3-4 hours remaining
