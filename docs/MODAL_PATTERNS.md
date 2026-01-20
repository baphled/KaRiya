# Modal Patterns Guide

**Last Updated**: 2026-01-20  
**Purpose**: Common modal usage patterns for KaRiya TUI  
**Status**: Production Ready

> **UIKit Migration**: Use `feedback.NewErrorModal()`, `feedback.NewLoadingModal()`, etc.
> from `internal/cli/uikit/feedback/`. See [UIKIT_GUIDE.md](./UIKIT_GUIDE.md) for details.

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
12. [Modal Overlays with bubbletea-overlay](#modal-overlays-with-bubbletea-overlay) **NEW!**

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

## Modal Overlays with bubbletea-overlay

### Overview

KaRiya uses the `bubbletea-overlay` library (v0.6.3) for compositing modal dialogs over background content. This provides reliable, flicker-free overlay rendering without manual ANSI manipulation.

**Library**: https://github.com/rmhubbert/bubbletea-overlay  
**License**: MIT  
**Stars**: 100+

### Why bubbletea-overlay?

✅ **Automatic Positioning**: Centers modals automatically  
✅ **Reliable Compositing**: No ANSI code conflicts  
✅ **Type-Safe**: Works with any `tea.Model`  
✅ **Tested**: Battle-tested in production applications  
✅ **Simple API**: Only 5 parameters needed

### When to Use Overlay Modals

**Use overlay modals when:**
- Showing content over the main view (timeline, list, etc.)
- Preserving context (background remains visible)
- Capturing user input without changing screens
- Displaying details, confirmations, or forms

**Don't use overlay modals when:**
- Transitioning to a completely different view
- Full-screen content is more appropriate
- No background context needs to be preserved

### Implementation Pattern

#### 1. Create Modal Component

Modal must implement `tea.Model` interface and use UIKit components:

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    
    // UIKit components (REQUIRED)
    "github.com/baphled/kariya/internal/cli/themes"
    "github.com/baphled/kariya/internal/cli/uikit/containers"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
)

type YourModal struct {
    theme   themes.Theme  // REQUIRED: Store theme for View()
    visible bool
    width   int
    height  int
    // ... your fields
}

func (m *YourModal) Init() tea.Cmd { return nil }

func (m *YourModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if !m.visible { return m, nil }
    
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    case tea.KeyMsg:
        // Handle keys
    }
    return m, nil
}

func (m *YourModal) View() string {
    if !m.visible { return "" }
    
    // REQUIRED: Nil theme guard
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    
    // REQUIRED: Use UIKit Box with solid background
    return containers.NewBox(theme).
        Content(content).
        Padding(2).
        Background(theme.BackgroundColor()). // REQUIRED: Prevents transparency!
        Render()
}
```

> **CRITICAL: UIKit Requirements for Modal View()**
>
> 1. **Nil theme guard**: Always check for nil theme and use default
> 2. **Use UIKit containers.Box**: Not direct lipgloss styling
> 3. **Solid background**: Always call `.Background(theme.BackgroundColor())`
>
> See [UIKit Guide](./UIKIT_GUIDE.md) for complete patterns.

#### 2. Add Modal to Intent

```go
type YourIntent struct {
    *BaseIntent
    yourModal *components.YourModal
    // ... other fields
}
```

#### 3. Show Modal

```go
func (i *YourIntent) handleShowModal() tea.Cmd {
    termInfo := i.GetTerminalInfo()
    width := 120
    height := 40
    if termInfo != nil {
        width = termInfo.Width
        height = termInfo.Height
    }
    
    i.yourModal = components.NewYourModal(width, height)
    i.yourModal.Show()
    return i.yourModal.Init()
}
```

#### 4. Handle Modal Updates

```go
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    // Handle modal BEFORE other logic
    if i.yourModal != nil && i.yourModal.IsVisible() {
        _, cmd := i.yourModal.Update(msg)
        if !i.yourModal.IsVisible() {
            // Modal closed - handle result
            i.yourModal = nil
        }
        return cmd
    }
    
    // ... rest of update logic
}
```

#### 5. Create staticViewModel Helper

```go
type staticViewModel struct {
    content string
}

func (s *staticViewModel) Init() tea.Cmd { return nil }
func (s *staticViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { 
    return s, nil 
}
func (s *staticViewModel) View() string { return s.content }
```

#### 6. Create Render Method

```go
func (i *YourIntent) renderYourModalOverlay(background string) string {
    // Wrap background in staticViewModel
    bgModel := &staticViewModel{content: background}
    
    // Use bubbletea-overlay to composite
    overlayModel := overlay.New(
        i.yourModal,    // Foreground (modal)
        bgModel,        // Background (main view)
        overlay.Center, // X position
        overlay.Center, // Y position
        0,              // X offset
        -2,             // Y offset (avoid footer)
    )
    
    return overlayModel.View()
}
```

#### 7. Integrate in View()

```go
func (i *YourIntent) View() string {
    // Render base view first
    baseView := i.renderMainContent()
    
    // If modal visible, overlay it
    if i.yourModal != nil && i.yourModal.IsVisible() {
        return i.renderYourModalOverlay(baseView)
    }
    
    return baseView
}
```

### Real-World Examples

#### Example 1: View Detail Modal (Read-Only)

**Component**: `internal/cli/components/view_event_detail_modal.go`  
**Usage**: Browse Timeline → Press Enter on event

```go
type ViewEventDetailModal struct {
    event   *career.CareerEvent
    theme   themes.Theme
    visible bool
    width   int
    height  int
    action  string // "", "edit", "delete"
}

// Key features:
// - Read-only display (no forms)
// - Action tracking (edit, delete, close)
// - Solid background to prevent transparency
// - Reuses existing RenderEventDetailCard component
```

**Intent Integration**:
```go
// Show modal on event selection
i.viewDetailModal = components.NewViewEventDetailModal(event, theme)
i.viewDetailModal.SetDimensions(width, height)
i.viewDetailModal.Show()

// Handle modal actions
action := i.viewDetailModal.GetAction()
switch action {
case "edit":
    i.editModal = components.NewEditEventModal(event, width, height)
case "delete":
    i.deleteModal = components.NewDeleteConfirmModal(...)
}
```

#### Example 2: Quick Add Modal (Form-Based)

**Component**: `internal/cli/components/quick_add_event_modal.go`  
**Usage**: Browse Timeline → Press 'a'

```go
type QuickAddEventModal struct {
    form      *huh.Form
    formData  *forms.EventFormData
    visible   bool
    completed bool
}

// Key features:
// - Huh form integration
// - Natural height for scrolling
// - Submit/cancel handling
// - Solid background wrapper
```

**Intent Integration**:
```go
// Show modal
i.quickAddModal = components.NewQuickAddEventModal(width, height)
return i.quickAddModal.Init()

// Handle completion
cmd, completed, eventData := i.quickAddModal.Update(msg)
if completed && eventData != nil {
    // Save event
    i.context.CLIEventService.CaptureEvent(...)
}
```

#### Example 3: Delete Confirmation Modal

**Component**: `internal/cli/components/delete_confirm_modal.go`  
**Usage**: Browse Timeline → Press 'd' on event

```go
type DeleteConfirmModal struct {
    entityName string
    title      string
    message    string
    visible    bool
    confirmed  bool
}

// Key features:
// - Simple yes/no confirmation
// - Red border for destructive actions
// - Clear button indicators
// - Solid background
```

**Intent Integration**:
```go
// Show modal
i.deleteModal = components.NewDeleteConfirmModal(
    event.Text,
    "Delete Event",
    fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
)
return i.deleteModal.Init()

// Handle confirmation
cmd, confirmed := i.deleteModal.Update(msg)
if confirmed {
    i.context.CLIEventService.DeleteEvent(ctx, event.ID)
}
```

### Common Patterns

#### Pattern 1: Modal with Actions

Modal returns action string instead of boolean:

```go
func (m *YourModal) GetAction() string {
    return m.action // "", "edit", "delete", "export", etc.
}

// In intent:
action := m.yourModal.GetAction()
switch action {
case "edit":
    // Handle edit
case "delete":
    // Handle delete
}
```

#### Pattern 2: Modal with Form Data

Modal returns structured data:

```go
func (m *YourModal) GetFormData() *YourFormData {
    return m.formData
}

// In intent:
if completed && eventData := m.yourModal.GetFormData(); eventData != nil {
    // Process form data
}
```

#### Pattern 3: Modal Chain

One modal triggers another:

```go
// Detail modal → Edit modal
if action := i.viewDetailModal.GetAction(); action == "edit" {
    i.viewDetailModal = nil
    i.editModal = components.NewEditEventModal(event, width, height)
    return i.editModal.Init()
}
```

### Modal Requirements Checklist

Use this checklist to ensure your modal implementation is complete and compliant with KaRiya standards.

#### 1. Modal Structure Requirements

- [ ] **Implements `tea.Model` interface** - Required for overlay compatibility
  - `Init() tea.Cmd`
  - `Update(msg tea.Msg) (tea.Model, tea.Cmd)` or custom signature
  - `View() string`
- [ ] **Has `visible bool` field** - Track modal visibility
- [ ] **Has `width int` and `height int` fields** - Responsive sizing
- [ ] **Has form/data fields** - Form model or data structure

#### 2. Update Method Requirements (CRITICAL)

- [ ] **Takes `tea.Msg` parameter** - NOT `tea.KeyMsg` (critical for huh forms)
  ```go
  func (m *YourModal) Update(msg tea.Msg) (tea.Cmd, bool, *FormData)
  ```
- [ ] **Handles `tea.WindowSizeMsg`** - Update width/height for responsive sizing
  ```go
  case tea.WindowSizeMsg:
      m.width = msg.Width
      m.height = msg.Height
  ```
- [ ] **Handles `tea.KeyMsg` for Esc** - Cancel/close modal
  ```go
  case tea.KeyMsg:
      switch msg.String() {
      case "esc":
          m.visible = false
          return nil, false, nil
      }
  ```
- [ ] **Returns early if not visible** - Performance optimization
  ```go
  if !m.visible { return nil, false, nil }
  ```
- [ ] **Forwards ALL messages to form** - Not just KeyMsg (critical!)
  ```go
  form, cmd := m.form.Update(msg)  // Pass full tea.Msg
  m.form = form.(*huh.Form)
  ```
- [ ] **Checks form completion** - Detect when user submits
  ```go
  if m.form.State == huh.StateCompleted {
      m.visible = false
      return cmd, true, m.formData
  }
  ```

#### 3. View Method Requirements

- [ ] **Returns empty string if not visible**
  ```go
  if !m.visible { return "" }
  ```
- [ ] **Has solid background** - CRITICAL to prevent transparency
  ```go
  Background(styles.ColorBackground)
  ```
- [ ] **Has rounded border** - Consistent styling
  ```go
  Border(lipgloss.RoundedBorder())
  ```
- [ ] **Has border color** - Theme-aware borders
  ```go
  BorderForeground(styles.ColorBorder)
  ```
- [ ] **Has padding** - Consistent spacing
  ```go
  Padding(1, 2)
  ```
- [ ] **Includes KeyBadge footer** - Show keyboard shortcuts (recommended)
  ```go
  footer := RenderHelpFooter(m.theme,
      NewKeyBadge("Tab", "Next field"),
      NewKeyBadge("Enter", "Submit"),
      NewKeyBadge("Esc", "Cancel"),
  )
  ```

#### 4. Form Building Requirements (for huh forms)

- [ ] **Uses 60% width calculation** - Consistent modal sizing
  ```go
  modalWidth := m.width * 60 / 100
  if modalWidth > 80 { modalWidth = 80 }
  if modalWidth < 40 { modalWidth = 40 }
  ```
- [ ] **Uses natural height (0)** - Let huh forms manage their own height
  ```go
  m.form = huh.NewForm(group).WithWidth(modalWidth)  // No WithHeight
  ```
- [ ] **Binds fields to formData** - Two-way data binding
  ```go
  Value(&m.formData.FieldName)
  ```
- [ ] **Has SubmitConfirmed field** - Track form submission
  ```go
  type FormData struct {
      Field1          string
      SubmitConfirmed bool
  }
  ```
- [ ] **Uses confirm button helper** - Consistent submission pattern
  ```go
  forms.NewFormWithFixedConfirm(group, &data.SubmitConfirmed, width, 0)
  ```

#### 5. Intent Integration Requirements

- [ ] **Modal field in intent struct**
  ```go
  type YourIntent struct {
      yourModal *components.YourModal
  }
  ```
- [ ] **Modal check BEFORE other logic** - Priority chain
  ```go
  // 1. Global keys (q, ?, m)
  // 2. Modal updates (if visible)
  // 3. Screen/state logic
  ```
- [ ] **Handler passes `tea.Msg`** - NOT `tea.KeyMsg` (critical!)
  ```go
  func (i *Intent) handleModalUpdate(msg tea.Msg) tea.Cmd {  // tea.Msg!
      cmd, applied, data := i.modal.Update(msg)
      // ...
  }
  ```
- [ ] **Creates modal with terminal dimensions**
  ```go
  termInfo := i.GetTerminalInfo()
  width, height := 120, 40
  if termInfo != nil {
      width = termInfo.Width
      height = termInfo.Height
  }
  i.modal = components.NewModal(width, height)
  ```
- [ ] **Calls Init() immediately** - Pattern 12: Immediate Init
  ```go
  return i.modal.Init()
  ```
- [ ] **Clears modal after use** - Prevent memory leaks
  ```go
  i.modal = nil
  ```

#### 6. Overlay Rendering Requirements

- [ ] **Has RenderOverlay method** - Composite modal over background
  ```go
  func (m *YourModal) RenderOverlay(baseView string) string
  ```
- [ ] **Uses staticViewModel wrapper** - For overlay compatibility
  ```go
  modalContent := staticViewModel{content: m.View()}
  bgModel := staticViewModel{content: baseView}
  ```
- [ ] **Uses overlay.New with Center positioning**
  ```go
  overlay.New(
      modalContent,   // Foreground
      bgModel,        // Background
      overlay.Center, // X position
      overlay.Center, // Y position
      0,              // X offset
      -2,             // Y offset (avoid footer)
  )
  ```
- [ ] **Intent checks IsVisible before rendering**
  ```go
  if i.modal != nil && i.modal.IsVisible() {
      return i.modal.RenderOverlay(baseView)
  }
  ```

#### 7. Keyboard Shortcuts Display Requirements

- [ ] **Modal footer shows shortcuts** - Tab/Enter/Esc for forms
- [ ] **List view shows navigation** - j/k/↑/↓ for lists
- [ ] **Uses RenderHelpFooter helper** - Consistent footer styling
- [ ] **Uses NewKeyBadge for each key** - Themed key badges
  ```go
  NewKeyBadge("j/k", "Navigate"),
  NewKeyBadge("Enter", "Select"),
  NewKeyBadge("Esc", "Back"),
  ```

#### 8. Testing Requirements

- [ ] **Component tests for toggle** - Show/Hide/IsVisible
- [ ] **Component tests for Tab navigation** - Field-to-field movement
- [ ] **Component tests for Enter submission** - Form completion
- [ ] **Component tests for Esc cancellation** - Close without saving
- [ ] **E2E tests for complete workflow** - Open → input → submit → verify
  - User opens modal
  - User types/navigates with Tab
  - User submits with Enter
  - System processes data correctly
- [ ] **Tests for WindowSizeMsg handling** - Responsive sizing

#### Compliance Quick Check

Run through this quick checklist:

1. ✅ **Update signature**: `Update(msg tea.Msg)` - NOT `tea.KeyMsg`
2. ✅ **Solid background**: `Background(styles.ColorBackground)`
3. ✅ **KeyBadge footer**: Shows Tab/Enter/Esc shortcuts
4. ✅ **E2E tests**: Proves complete workflow works
5. ✅ **Intent passes `tea.Msg`**: Handler doesn't cast to `tea.KeyMsg`

**If ANY of these fail, the modal is NOT compliant.**

---

### Best Practices

#### ✅ DO

1. **Always set solid background** to prevent transparency
   ```go
   Background(styles.ColorBackground)
   ```

2. **Handle WindowSizeMsg** for responsive modals
   ```go
   case tea.WindowSizeMsg:
       m.width = msg.Width
       m.height = msg.Height
   ```

3. **Check visibility before updating**
   ```go
   if !m.visible { return m, nil }
   ```

4. **Use Y offset of -2** to avoid footer overlap
   ```go
   overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
   ```

5. **Clear modal after use**
   ```go
   i.yourModal = nil
   ```

6. **Wrap background in staticViewModel**
   ```go
   bgModel := &staticViewModel{content: background}
   ```

#### ❌ DON'T

1. **Don't forget solid background** - causes transparency issues

2. **Don't manually center with ANSI** - use overlay.Center

3. **Don't stack modals** - hide one before showing another

4. **Don't constrain Huh form height** - use natural height (0)

5. **Don't forget to implement tea.Model** - required for overlay

6. **Don't use manual width calculations** - use MaxWidth/MaxHeight

### Troubleshooting

#### Issue: Background shows through modal (transparency)

**Solution**: Add solid background in modal's View():
```go
Background(styles.ColorBackground)
```

#### Issue: Modal not centered

**Solution**: Use overlay.Center for both X and Y:
```go
overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
```

#### Issue: Huh form doesn't scroll

**Solution**: Use natural height (0) instead of constrained height:
```go
form = forms.NewYourFormWithDataAndDimensions(data, width, 0) // 0 = natural
```

#### Issue: Modal overlaps footer

**Solution**: Use Y offset of -2:
```go
overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
```

#### Issue: Modal doesn't update on window resize

**Solution**: Handle WindowSizeMsg in modal Update():
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

### Complete Example: Browse Timeline Modals

All 5 Browse Timeline modals use this pattern:

| Modal | File | Lines | Features |
|-------|------|-------|----------|
| ViewEventDetailModal | `view_event_detail_modal.go` | 151 | Read-only, action tracking |
| QuickAddEventModal | `quick_add_event_modal.go` | 250 | Huh form, natural height |
| EditEventModal | `edit_event_modal.go` | 280 | Huh form, preserve original |
| DeleteConfirmModal | `delete_confirm_modal.go` | 200 | Yes/no, destructive style |
| FilterModalModel | `filter_modal.go` | 350 | Multi-field form |

**Intent Integration**: `internal/cli/intents/browse_timeline_intent.go`

- 5 modal fields in struct
- 5 render methods (one per modal)
- 5 update handlers (one per modal)
- 5 View() overlay checks

**See**: `VIEW_DETAIL_MODAL_SUMMARY.md` for complete implementation details.

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
- [VIEW_DETAIL_MODAL_SUMMARY.md](../VIEW_DETAIL_MODAL_SUMMARY.md) - Complete modal overlay implementation example
- [MODAL_REFACTOR_VERIFICATION.md](../MODAL_REFACTOR_VERIFICATION.md) - Modal refactoring verification guide

---

**Questions? Issues?** Open an issue on GitHub or check the examples above.
