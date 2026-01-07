# Modal Patterns Guide

**Last Updated**: 2026-01-07  
**Purpose**: Common modal usage patterns for KaRiya TUI  
**Status**: Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Modal Types](#modal-types)
3. [When to Use Each Type](#when-to-use-each-type)
4. [Error Modals](#error-modals)
5. [Loading Modals](#loading-modals)
6. [Progress Modals](#progress-modals)
7. [Success Modals](#success-modals)
8. [Warning Modals](#warning-modals)
9. [Modal Timing](#modal-timing)
10. [Accessibility](#accessibility)
11. [Common Patterns](#common-patterns)

---

## Overview

Modals are **overlay dialogs** that appear on top of the main content to provide immediate feedback or request user attention.

### Key Principles

- **One modal at a time** - Never stack modals
- **Clear purpose** - User should immediately understand what happened
- **Easy dismissal** - Escape key or explicit action
- **Accessibility** - Bell alerts, clear visual indicators
- **Non-blocking when possible** - Don't freeze the UI unnecessarily

---

## Modal Types

| Type | Icon | Purpose | Auto-Dismiss | Cancellable |
|------|------|---------|--------------|-------------|
| Error | ✖ | Show errors | No | Yes (Esc) |
| Loading | ⏳ | Indeterminate progress | No | Optional |
| Progress | ▓▓▓░░ | Determinate progress | No | No |
| Success | ✓ | Success feedback | Yes (3s) | No |
| Warning | ⚠ | Warnings | No | Yes (Esc) |

---

## When to Use Each Type

### Error Modals

**Use when:**
- An operation fails
- User input is invalid
- System encounters an error
- Need to alert user to a problem

**Don't use when:**
- Error is expected/recoverable inline (use inline validation instead)
- Multiple errors occur (combine into one message)

### Loading Modals

**Use when:**
- Operation takes > 500ms
- Progress is indeterminate
- User needs to wait

**Don't use when:**
- Operation is instant (< 500ms)
- Progress can be tracked (use Progress modal instead)

### Progress Modals

**Use when:**
- Progress is measurable (0% to 100%)
- Multi-step operations
- User needs to see progress

**Don't use when:**
- Progress is indeterminate (use Loading modal)
- Operation is instant

### Success Modals

**Use when:**
- Operation completes successfully
- User needs confirmation
- Want to provide positive feedback

**Don't use when:**
- Success is obvious from context
- Would interrupt workflow

### Warning Modals

**Use when:**
- Action is potentially dangerous
- Need user acknowledgment
- Want to prevent accidental actions

**Don't use when:**
- Could use inline warning instead
- Warning is not critical

---

## Error Modals

### Basic Error

```go
modal := components.NewErrorModal(
    "Save Failed",
    "Could not save file: permission denied"
)

view.ShowModalOverlay(modal)
```

**Features:**
- Terminal bell alert (beep)
- Red border and error icon (✖)
- Escape key to dismiss
- No auto-dismiss (requires acknowledgment)

### Validation Error

```go
modal := components.NewErrorModal(
    "Invalid Input",
    "Please fill in all required fields:\n• Name\n• Email\n• Date"
)
```

### Network Error

```go
modal := components.NewErrorModal(
    "Connection Failed",
    "Could not connect to server. Please check your network connection and try again."
)
```

### Error with Recovery Hint

```go
modal := components.NewErrorModal(
    "File Not Found",
    "Could not find config.yaml\n\nTry running 'kariya init' to create a default configuration."
)
```

### Critical Error

```go
modal := components.NewErrorModal(
    "Critical Error",
    "Database corrupted. Application will exit.\n\nPlease contact support with error code: ERR_DB_001"
)
modal.Bell = true  // Ensure bell is enabled for critical errors
```

---

## Loading Modals

### Basic Loading

```go
modal := components.NewLoadingModal("Loading data...", false)
view.ShowModalOverlay(modal)
```

### Loading with Spinner

The spinner animates automatically:

```go
modal := components.NewLoadingModal("Processing your request...", false)
// Spinner rotates through: ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
```

### Cancellable Loading

```go
modal := components.NewLoadingModal("Fetching data...", true)  // cancellable = true
// User can press Esc to cancel
```

### Loading with Message Rotation

For long operations, rotate through messages:

```go
messages := []string{
    "Connecting to server...",
    "Authenticating...",
    "Fetching data...",
    "Processing results...",
}

rotator := components.NewLoadingMessageRotator(messages, 2*time.Second)
modal := components.NewLoadingModal("Loading", false)
modal.SetMessageRotator(rotator)

// Messages automatically rotate every 2 seconds
```

### Pattern: Long Network Request

```go
// Show loading modal
modal := components.NewLoadingModal("Fetching data from API...", true)
view.ShowModalOverlay(modal)

// In your Update() method:
case types.DataFetchedMsg:
    // Hide modal, show data
    view.ShowModal = false
    
case tea.KeyMsg:
    if msg.String() == "esc" && modal.Cancellable {
        // Cancel operation
        return m, cancelFetch()
    }
```

---

## Progress Modals

### Basic Progress

```go
// progress: 0.0 to 1.0
modal := components.NewProgressModal("Processing", "Analyzing files...", 0.0)
view.ShowModalOverlay(modal)
```

### Updating Progress

```go
// Start at 0%
modal := components.NewProgressModal("Processing", "Starting...", 0.0)

// Update as work progresses
modal.UpdateProgress(0.25)  // 25%
modal.UpdateProgress(0.50)  // 50%
modal.UpdateProgress(0.75)  // 75%
modal.UpdateProgress(1.00)  // 100%
```

### Multi-Step Process

```go
type ProcessStep struct {
    Progress float64
    Title    string
    Message  string
}

steps := []ProcessStep{
    {0.00, "Starting", "Initializing process..."},
    {0.25, "Analyzing", "Analyzing career events..."},
    {0.50, "Calculating", "Calculating impact metrics..."},
    {0.75, "Generating", "Generating professional bullets..."},
    {1.00, "Complete", "CV generation complete!"},
}

for _, step := range steps {
    modal := components.NewProgressModal(step.Title, step.Message, step.Progress)
    view.ShowModalOverlay(modal)
    
    // Perform work...
    doWork(step)
}
```

### Pattern: CV Generation

```go
func (i *GenerateCVIntent) generateCV() tea.Cmd {
    return func() tea.Msg {
        // Step 1: Analyze (0-25%)
        i.progressModal = components.NewProgressModal(
            "Analyzing", 
            "Analyzing career events...", 
            0.0,
        )
        // ... do work ...
        i.progressModal.UpdateProgress(0.25)
        
        // Step 2: Calculate (25-50%)
        i.progressModal = components.NewProgressModal(
            "Calculating",
            "Calculating impact metrics...",
            0.25,
        )
        // ... do work ...
        i.progressModal.UpdateProgress(0.50)
        
        // Step 3: Generate (50-75%)
        i.progressModal = components.NewProgressModal(
            "Generating",
            "Generating professional bullets...",
            0.50,
        )
        // ... do work ...
        i.progressModal.UpdateProgress(0.75)
        
        // Step 4: Format (75-100%)
        i.progressModal = components.NewProgressModal(
            "Formatting",
            "Formatting final document...",
            0.75,
        )
        // ... do work ...
        i.progressModal.UpdateProgress(1.00)
        
        return CVGeneratedMsg{cv: result}
    }
}
```

---

## Success Modals

### Basic Success

```go
modal := components.NewSuccessModal("File saved successfully!")
view.ShowModalOverlay(modal)
// Auto-dismisses after 3 seconds
```

### Success with Details

```go
modal := components.NewSuccessModal(
    "CV generated successfully!\n\nSaved to: /path/to/cv.pdf\nPages: 2\nWords: 847"
)
```

### Success with Next Steps

```go
modal := components.NewSuccessModal(
    "Export complete!\n\n✓ File: output.yaml\n✓ Location: ~/kariya/exports/\n\nPress 'v' to view or 'q' to quit"
)
```

### Pattern: Save Confirmation

```go
// After successful save
modal := components.NewSuccessModal(
    fmt.Sprintf("Saved %d events to database", count)
)
view.ShowModalOverlay(modal)

// Auto-dismiss after 3 seconds
time.AfterFunc(3*time.Second, func() {
    view.ShowModal = false
})
```

---

## Warning Modals

### Basic Warning

```go
modal := components.NewWarningModal(
    "Unsaved Changes",
    "You have unsaved changes. Are you sure you want to quit?"
)
```

### Destructive Action Warning

```go
modal := components.NewWarningModal(
    "Delete Confirmation",
    "This will permanently delete 15 events. This action cannot be undone.\n\nPress 'y' to confirm or 'n' to cancel."
)
```

### Pattern: Confirmation Dialog

```go
// Show warning
modal := components.NewWarningModal(
    "Confirm Delete",
    "Delete this burst and all associated events?"
)
modal.Actions = []string{"y Yes", "n No"}
view.ShowModalOverlay(modal)

// In Update():
case tea.KeyMsg:
    if msg.String() == "y" {
        return m, deleteBurst()
    }
    if msg.String() == "n" || msg.String() == "esc" {
        view.ShowModal = false
    }
```

---

## Modal Timing

### Fade-In Animation

All modals fade in over 150ms for a smooth appearance:

```go
modal := components.NewErrorModal("Error", "Message")
modal.FadeInDuration = 150 * time.Millisecond  // Default
```

### Auto-Dismiss (Success Only)

Success modals auto-dismiss after 3 seconds:

```go
modal := components.NewSuccessModal("Success!")
modal.AutoDismiss = 3 * time.Second  // Default
```

### Custom Auto-Dismiss

```go
modal := components.NewSuccessModal("Quick notification")
modal.AutoDismiss = 1 * time.Second  // Dismiss after 1 second
```

### No Auto-Dismiss

```go
modal := components.NewSuccessModal("Important success")
modal.AutoDismiss = 0  // Never auto-dismiss
```

---

## Accessibility

### Terminal Bell

Error modals trigger a terminal bell for accessibility:

```go
modal := components.NewErrorModal("Error", "Message")
// modal.Bell = true  (automatic for errors)
```

### Keyboard Navigation

All modals support keyboard dismissal:

- **Escape** - Dismiss error/warning modals
- **Any key** - Acknowledge and dismiss (after fade-in)
- **y/n** - Confirm/cancel for warnings with actions

### Visual Indicators

Each modal type has distinct visual styling:

| Type | Border Color | Icon | Style |
|------|--------------|------|-------|
| Error | Red | ✖ | Bold, urgent |
| Loading | Blue | ⏳ | Animated spinner |
| Progress | Blue | ▓▓▓░░ | Progress bar |
| Success | Green | ✓ | Positive, light |
| Warning | Yellow | ⚠ | Attention-grabbing |

### Screen Reader Support

- Modals use clear, descriptive titles
- Icons are accompanied by text
- Important information appears first

---

## Common Patterns

### Pattern 1: Error → Retry

```go
// Show error
errorModal := components.NewErrorModal(
    "Connection Failed",
    "Could not connect to server.\n\nPress 'r' to retry or 'q' to quit."
)
view.ShowModalOverlay(errorModal)

// In Update():
case tea.KeyMsg:
    if msg.String() == "r" {
        view.ShowModal = false
        return m, retryConnection()
    }
```

### Pattern 2: Loading → Success

```go
// Start with loading
loadingModal := components.NewLoadingModal("Saving...", false)
view.ShowModalOverlay(loadingModal)

// On completion, show success
case SaveCompleteMsg:
    successModal := components.NewSuccessModal("Saved successfully!")
    view.ShowModalOverlay(successModal)
```

### Pattern 3: Progress → Success

```go
// Show progress
progressModal := components.NewProgressModal("Processing", "Working...", 0.5)
view.ShowModalOverlay(progressModal)

// Update progress
progressModal.UpdateProgress(1.0)

// Switch to success
case ProcessCompleteMsg:
    successModal := components.NewSuccessModal("Processing complete!")
    view.ShowModalOverlay(successModal)
```

### Pattern 4: Warning → Action

```go
// Show warning
warningModal := components.NewWarningModal(
    "Confirm Action",
    "This will modify 50 events. Continue?"
)
warningModal.Actions = []string{"y Yes", "n No"}
view.ShowModalOverlay(warningModal)

// Handle response
case tea.KeyMsg:
    if msg.String() == "y" {
        view.ShowModal = false
        return m, performAction()
    }
```

### Pattern 5: Loading with Timeout

```go
// Show loading with timeout
loadingModal := components.NewLoadingModal("Connecting...", true)
view.ShowModalOverlay(loadingModal)

// Set timeout
time.AfterFunc(30*time.Second, func() {
    errorModal := components.NewErrorModal(
        "Timeout",
        "Connection timed out after 30 seconds."
    )
    view.ShowModalOverlay(errorModal)
})
```

---

## Best Practices

### Do's ✅

1. **Use appropriate modal types**
   ```go
   // Error for failures
   NewErrorModal("Save Failed", "Permission denied")
   
   // Success for completions
   NewSuccessModal("Saved successfully!")
   ```

2. **Provide context in messages**
   ```go
   // Good
   "Could not save file: disk full (Error: ENOSPC)"
   
   // Bad
   "Save failed"
   ```

3. **Make loading operations cancellable when possible**
   ```go
   NewLoadingModal("Fetching data...", true)  // cancellable
   ```

4. **Use progress modals for multi-step operations**
   ```go
   // Show progress from 0% to 100%
   modal.UpdateProgress(0.25)
   modal.UpdateProgress(0.50)
   modal.UpdateProgress(0.75)
   modal.UpdateProgress(1.00)
   ```

### Don'ts ❌

1. **Don't show multiple modals simultaneously**
   ```go
   // Bad
   view.ShowModalOverlay(modal1)
   view.ShowModalOverlay(modal2)  // Replaces modal1, confusing
   ```

2. **Don't use success modals for obvious actions**
   ```go
   // Bad - success is obvious from context
   NewSuccessModal("Quit successful!")
   ```

3. **Don't forget to dismiss modals**
   ```go
   // Bad - modal stays forever
   view.ShowModalOverlay(errorModal)
   
   // Good - handle dismissal
   case tea.KeyMsg:
       if msg.String() == "esc" {
           view.ShowModal = false
       }
   ```

4. **Don't use loading modals for instant operations**
   ```go
   // Bad - operation takes < 100ms
   modal := NewLoadingModal("Loading...", false)
   doInstantOperation()
   
   // Good - no modal needed
   doInstantOperation()
   ```

---

## Testing

Modal behavior is tested in:
- `internal/cli/components/modal_test.go` - Unit tests
- `internal/cli/components/performance_test.go` - Performance benchmarks
- `cmd/test_all_views/main.go` - Visual test program

Run visual tests:
```bash
go build -o test_views ./cmd/test_all_views
./test_views
# Navigate to modal scenarios with arrow keys
```

---

## Related Documentation

- [STANDARDVIEW_GUIDE.md](STANDARDVIEW_GUIDE.md) - StandardView usage
- [TUI_DEVELOPER_GUIDE.md](TUI_DEVELOPER_GUIDE.md) - General TUI development
- [TUI_STANDARDS.md](TUI_STANDARDS.md) - Design standards

---

**Questions? Issues?** Open an issue on GitHub or check the examples above.
