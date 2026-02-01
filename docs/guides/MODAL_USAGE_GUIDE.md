# Modal Usage Guide

This guide documents when to use each feedback modal type to ensure consistent user experience across the application.

## Feedback Modal Types

The `feedback` package (`internal/cli/uikit/feedback/modal.go`) provides five modal types via `*feedback.Modal`:

| Type | Constructor | Icon | Border | Bell | Dismissal | Purpose |
|------|-------------|------|--------|------|-----------|---------|
| `ModalError` | `NewErrorModal(title, msg)` | `⚠️` | Red (`ErrorColor`) | Yes | Esc to dismiss | Operation failures |
| `ModalWarning` | `NewWarningModal(title, msg)` | `⚠️` | Amber (`WarningColor`) | Yes | Esc to close | Validation guards, empty results |
| `ModalSuccess` | `NewSuccessModal(msg)` | `✅` | Green (`SuccessColor`) | No | Auto-dismiss 3s | Successful completions |
| `ModalLoading` | `NewLoadingModal(msg, cancel)` | `⏳` | Blue (`InfoColor`) | No | Esc to cancel | Async operations in progress |
| `ModalProgress` | `NewProgressModal(title, msg, val)` | `📊` | Blue (`InfoColor`) | No | None | Operations with known progress |

There is also a separate `InfoModal` struct (`info_modal.go`) with three variants:

| Variant | Constructor | Border | Purpose |
|---------|-------------|--------|---------|
| `InfoModalInfo` | `NewInfoModal(title, msg)` | Teal/Blue | Informational messages |
| `InfoModalWarning` | `NewWarningInfoModal(title, msg)` | Amber | Warning information |
| `InfoModalSuccess` | `NewSuccessInfoModal(title, msg)` | Green | Success information |

**Note:** `InfoModal` is a different struct from `Modal` with a different API (`Update(msg) bool`). Use `*feedback.Modal` types for intent-level feedback modals stored in `feedbackModal`.

---

## When to Use Each Type

### ModalError — Operation Failures

Use when an operation was attempted and **failed due to a system error**.

```go
// Database write failed
i.ShowErrorModal("Update Failed", err.Error())

// Service unavailable
i.ShowErrorModal("Inference Failed", "Skill inference service not available")

// Delete operation failed
i.ShowErrorModal("Delete Failed", err.Error())
```

**Characteristics:** The user attempted an action, the system tried to perform it, and it failed. The user needs to know the operation did not complete.

### ModalWarning — Validation Guards and Empty Results

Use when **no failure occurred** but the user needs to be informed of a condition.

```go
// Precondition not met (no failure, just a guard)
i.ShowWarningModal("Burst Not Confirmed", "Please confirm this burst before inferring skills.")

// Empty results from a successful operation
i.ShowWarningModal("No Skills Detected",
    "No skills were detected from the burst events. "+
        "The events may not contain enough technical details.")

// No data to process
i.ShowWarningModal("No Suggestions Found",
    "No suggestions were generated from your events. "+
        "Try adding more events or adjusting detection settings.")
```

**Characteristics:** The system operated correctly but the outcome is empty or the user's request cannot proceed. The user is not at fault and no error occurred.

### ModalSuccess — Successful Completions

Use when an operation **completed successfully** and the user needs confirmation.

```go
// Skills created
i.ShowSuccessModal("Skills Created", fmt.Sprintf("Successfully created %d skill(s)", len(skills)))

// Configuration saved
i.feedbackModal = feedback.NewSuccessModal("Configuration saved!")

// Event captured
i.feedbackModal = feedback.NewSuccessModal("Event saved!")
```

**Characteristics:** An operation the user initiated completed without error. The modal auto-dismisses after 3 seconds.

### ModalLoading — Async Operations

Use when an **asynchronous operation is in progress** and the user should wait.

```go
i.loadingModal = feedback.NewLoadingModal("Detecting burst patterns...", true).WithTheme(i.Theme())
return tea.Batch(i.loadingModal.Init(), i.startDetection())
```

**Characteristics:** Always pair with `.Init()` to start the spinner. Set `cancellable: true` for user-cancellable operations. Store in `loadingModal`, not `feedbackModal`.

---

## Decision Flowchart

```
Did an operation fail with an error?
├── YES → ModalError (ShowErrorModal)
└── NO
    ├── Did an operation succeed?
    │   ├── YES → ModalSuccess (ShowSuccessModal)
    │   └── NO
    │       ├── Is a condition preventing the user's action?
    │       │   ├── YES → ModalWarning (ShowWarningModal)
    │       │   └── NO
    │       │       ├── Is an async operation in progress?
    │       │       │   ├── YES → ModalLoading (NewLoadingModal)
    │       │       │   └── NO → No modal needed
    │       └── Were results empty from a successful query?
    │           └── YES → ModalWarning (ShowWarningModal)
```

---

## Intent Field Naming Convention

### feedbackModal (not errorModal)

Intents that display feedback modals should name the field `feedbackModal` since it holds error, warning, and success modals:

```go
type Intent struct {
    *intents.BaseIntent

    // feedbackModal holds the feedback modal (shown for errors, warnings, and success messages).
    feedbackModal *feedback.Modal

    // loadingModal holds the loading modal (shown during async operations).
    loadingModal *feedback.Modal
}
```

**Rationale:** The field stores `*feedback.Modal` which can be any type (error, warning, success). Naming it `errorModal` is misleading when it also holds success and warning modals.

### Helper Methods

Intents should provide convenience helpers for each modal type:

```go
func (i *Intent) ShowErrorModal(title, message string) {
    i.feedbackModal = feedback.NewErrorModal(title, message)
}

func (i *Intent) ShowWarningModal(title, message string) {
    i.feedbackModal = feedback.NewWarningModal(title, message)
}

func (i *Intent) ShowSuccessModal(title, message string) {
    modal := feedback.NewSuccessModal(message)
    modal.Title = title
    i.feedbackModal = modal
}
```

### Handler Naming

The modal update handler should be named `handleFeedbackModalUpdate`:

```go
func (i *Intent) handleFeedbackModalUpdate(msg tea.Msg) tea.Cmd {
    if i.feedbackModal == nil {
        return nil
    }
    if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
        i.feedbackModal = nil
    }
    return noopCmd
}
```

---

## Alternative Patterns in the Codebase

### BaseIntent Transient Success (view_helpers.go)

Intents can use `BaseIntent.SetSuccess(message)` for transient success messages that auto-expire after 3 seconds. These are rendered automatically by `applyStateModals()` in `view_helpers.go`:

```go
i.SetSuccess("Operation completed")
```

This approach creates modals transiently during rendering and does not require a stored field.

### resultModal / submitModal (Legacy)

Older intents (`ConfigureSystemIntent`, `CaptureEventIntent`) use `resultModal` or `submitModal` as dual-purpose fields. New intents should use `feedbackModal` with the helper methods described above.

---

## Common Mistakes

| Mistake | Correct Approach |
|---------|------------------|
| Using `ShowErrorModal` for empty results | Use `ShowWarningModal` — no failure occurred |
| Using `ShowErrorModal` for validation guards | Use `ShowWarningModal` — the user needs guidance, not an error |
| Using `ShowErrorModal` for success messages | Use `ShowSuccessModal` — renders with green styling |
| Naming the field `errorModal` | Name it `feedbackModal` — it holds multiple modal types |
| Direct `i.feedbackModal = feedback.New*Modal(...)` | Use helpers: `ShowErrorModal`, `ShowWarningModal`, `ShowSuccessModal` |
| Storing loading modals in `feedbackModal` | Use a separate `loadingModal` field — loading has different lifecycle |

---

## Testing Modal Types

Use `GetFeedbackModal()` to verify the correct modal type in tests:

```go
modal := intent.GetFeedbackModal()
Expect(modal).NotTo(BeNil())
Expect(modal.Type).To(Equal(feedback.ModalWarning))  // or ModalError, ModalSuccess
```

For visibility checks:

```go
Expect(intent.HasVisibleFeedbackModal()).To(BeTrue())
```

---

## Reference

- Modal implementation: `internal/cli/uikit/feedback/modal.go`
- InfoModal implementation: `internal/cli/uikit/feedback/info_modal.go`
- ConfirmModal implementation: `internal/cli/uikit/feedback/confirm_modal.go`
- Burst management helpers: `internal/cli/intents/burst_management/helpers.go`
- UIKit guide: `docs/UIKIT_GUIDE.md`
- Modal patterns: `docs/MODAL_PATTERNS.md`
