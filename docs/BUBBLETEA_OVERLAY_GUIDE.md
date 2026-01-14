---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# bubbletea-overlay Library Guide

**Library**: https://github.com/rmhubbert/bubbletea-overlay  
**Version**: v0.6.3  
**License**: MIT  
**Stars**: 100+  
**Purpose**: Reliable modal overlay compositing for Bubble Tea applications

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Basic Usage](#basic-usage)
4. [KaRiya Integration Pattern](#kariya-integration-pattern)
5. [Complete Example](#complete-example)
6. [API Reference](#api-reference)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)
9. [Real-World Examples](#real-world-examples)

---

## Overview

`bubbletea-overlay` is a library for compositing Bubble Tea models on top of each other. It provides:

✅ **Automatic Positioning**: Centers content without manual calculations  
✅ **Reliable Compositing**: No ANSI code conflicts or rendering artifacts  
✅ **Type-Safe**: Works with any `tea.Model`  
✅ **Simple API**: Only 5 parameters needed  
✅ **Battle-Tested**: Used in production applications

### Why Use It?

**Before** (manual compositing):
- Complex ANSI escape code manipulation
- Manual centering calculations
- Flicker and rendering artifacts
- Z-index conflicts
- Hard to debug

**After** (bubbletea-overlay):
- Simple 5-parameter API
- Automatic positioning
- Clean compositing
- No artifacts
- Easy to maintain

---

## Installation

```bash
go get github.com/rmhubbert/bubbletea-overlay@v0.6.3
```

**go.mod**:
```go
require (
    github.com/rmhubbert/bubbletea-overlay v0.6.3
)
```

---

## Basic Usage

### Step 1: Create Foreground Model

Modal must implement `tea.Model`:

```go
type MyModal struct {
    content string
}

func (m MyModal) Init() tea.Cmd {
    return nil
}

func (m MyModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil
}

func (m MyModal) View() string {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Render(m.content)
}
```

### Step 2: Create Background Model

Wrap existing view content:

```go
type staticViewModel struct {
    content string
}

func (s staticViewModel) Init() tea.Cmd { return nil }
func (s staticViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { 
    return s, nil 
}
func (s staticViewModel) View() string { 
    return s.content 
}
```

### Step 3: Composite with overlay.New()

```go
import "github.com/rmhubbert/bubbletea-overlay"

func compositeView(background string, modal *MyModal) string {
    bgModel := &staticViewModel{content: background}
    
    overlayModel := overlay.New(
        modal,          // Foreground (modal)
        bgModel,        // Background (main view)
        overlay.Center, // X position
        overlay.Center, // Y position
        0,              // X offset
        0,              // Y offset
    )
    
    return overlayModel.View()
}
```

---

## KaRiya Integration Pattern

KaRiya uses a standardized pattern for all modal overlays:

### 1. Modal Component Structure

```go
type YourModal struct {
    visible bool    // Control visibility
    width   int     // Terminal width
    height  int     // Terminal height
    // ... your fields
}

func (m *YourModal) View() string {
    if !m.visible { return "" }
    
    // CRITICAL: Always include solid background
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder).
        Background(styles.ColorBackground). // Prevents transparency!
        Padding(1, 2).
        MaxWidth(m.width - 8).
        MaxHeight(m.height - 8).
        Render(content)
}
```

### 2. staticViewModel Helper

**Location**: Each intent defines this locally

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

### 3. Render Method in Intent

```go
func (i *YourIntent) renderYourModalOverlay(background string) string {
    // Wrap background
    bgModel := &staticViewModel{content: background}
    
    // Composite with bubbletea-overlay
    overlayModel := overlay.New(
        i.yourModal,    // Foreground
        bgModel,        // Background
        overlay.Center, // X position
        overlay.Center, // Y position
        0,              // X offset
        -2,             // Y offset (avoids footer)
    )
    
    return overlayModel.View()
}
```

### 4. View Integration

```go
func (i *YourIntent) View() string {
    // Render base view
    baseView := i.renderMainContent()
    
    // Overlay modal if visible
    if i.yourModal != nil && i.yourModal.IsVisible() {
        return i.renderYourModalOverlay(baseView)
    }
    
    return baseView
}
```

---

## Complete Example

### Modal Component

```go
// internal/cli/components/example_modal.go
package components

import (
    "github.com/baphled/kariya/internal/cli/styles"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type ExampleModal struct {
    title   string
    message string
    visible bool
    width   int
    height  int
}

func NewExampleModal(title, message string, width, height int) *ExampleModal {
    return &ExampleModal{
        title:   title,
        message: message,
        visible: false,
        width:   width,
        height:  height,
    }
}

func (m *ExampleModal) Init() tea.Cmd {
    return nil
}

func (m *ExampleModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if !m.visible {
        return m, nil
    }
    
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
        
    case tea.KeyMsg:
        switch msg.String() {
        case "esc", "enter", "q":
            m.Hide()
            return m, nil
        }
    }
    
    return m, nil
}

func (m *ExampleModal) View() string {
    if !m.visible {
        return ""
    }
    
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        lipgloss.NewStyle().Bold(true).Render(m.title),
        "",
        m.message,
        "",
        lipgloss.NewStyle().
            Foreground(styles.ColorTextSecondary).
            Render("Press Enter or Esc to close"),
    )
    
    // Solid background prevents transparency
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder).
        Background(styles.ColorBackground).
        Padding(1, 2).
        MaxWidth(m.width - 8).
        Render(content)
}

func (m *ExampleModal) Show() { m.visible = true }
func (m *ExampleModal) Hide() { m.visible = false }
func (m *ExampleModal) IsVisible() bool { return m.visible }
```

### Intent Integration

```go
// internal/cli/intents/example_intent.go
package intents

import (
    "github.com/baphled/kariya/internal/cli/components"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/rmhubbert/bubbletea-overlay"
)

type staticViewModel struct {
    content string
}

func (s *staticViewModel) Init() tea.Cmd { return nil }
func (s *staticViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return s, nil
}
func (s *staticViewModel) View() string { return s.content }

type ExampleIntent struct {
    *BaseIntent
    exampleModal *components.ExampleModal
}

func (i *ExampleIntent) Update(msg tea.Msg) tea.Cmd {
    // Handle modal updates first
    if i.exampleModal != nil && i.exampleModal.IsVisible() {
        _, cmd := i.exampleModal.Update(msg)
        if !i.exampleModal.IsVisible() {
            // Modal closed
            i.exampleModal = nil
        }
        return cmd
    }
    
    // Handle key presses
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        switch keyMsg.String() {
        case "m":
            // Show modal
            termInfo := i.GetTerminalInfo()
            width, height := 120, 40
            if termInfo != nil {
                width = termInfo.Width
                height = termInfo.Height
            }
            i.exampleModal = components.NewExampleModal(
                "Example Modal",
                "This is an example modal overlay!",
                width,
                height,
            )
            i.exampleModal.Show()
            return i.exampleModal.Init()
        }
    }
    
    return nil
}

func (i *ExampleIntent) View() string {
    // Render main content
    baseView := "Main content here\n\nPress 'm' to show modal"
    
    // Overlay modal if visible
    if i.exampleModal != nil && i.exampleModal.IsVisible() {
        return i.renderExampleModalOverlay(baseView)
    }
    
    return baseView
}

func (i *ExampleIntent) renderExampleModalOverlay(background string) string {
    bgModel := &staticViewModel{content: background}
    
    overlayModel := overlay.New(
        i.exampleModal, // Foreground
        bgModel,        // Background
        overlay.Center, // X position
        overlay.Center, // Y position
        0,              // X offset
        -2,             // Y offset
    )
    
    return overlayModel.View()
}
```

---

## API Reference

### overlay.New()

```go
func New(
    fg tea.Model,  // Foreground model (modal)
    bg tea.Model,  // Background model (main view)
    x Position,    // X position (Left, Center, Right)
    y Position,    // Y position (Top, Center, Bottom)
    xOffset int,   // X offset in characters
    yOffset int,   // Y offset in lines
) tea.Model
```

**Parameters**:
- `fg`: The modal to display on top (must implement tea.Model)
- `bg`: The background view (must implement tea.Model)
- `x`: Horizontal position (overlay.Left, overlay.Center, overlay.Right)
- `y`: Vertical position (overlay.Top, overlay.Center, overlay.Bottom)
- `xOffset`: Additional horizontal offset (positive = right, negative = left)
- `yOffset`: Additional vertical offset (positive = down, negative = up)

**Returns**: A new tea.Model that composites fg on top of bg

### Position Constants

```go
overlay.Left    // Align to left edge
overlay.Center  // Center horizontally/vertically
overlay.Right   // Align to right edge
overlay.Top     // Align to top edge
overlay.Bottom  // Align to bottom edge
```

---

## Best Practices

### ✅ DO

1. **Always set solid background** in modal View():
   ```go
   Background(styles.ColorBackground)
   ```

2. **Use overlay.Center for both X and Y** (most common):
   ```go
   overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
   ```

3. **Use Y offset of -2** to avoid footer overlap:
   ```go
   overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
   ```

4. **Wrap background in staticViewModel**:
   ```go
   bgModel := &staticViewModel{content: background}
   ```

5. **Return empty string when modal not visible**:
   ```go
   func (m *Modal) View() string {
       if !m.visible { return "" }
       // ...
   }
   ```

6. **Handle WindowSizeMsg for responsiveness**:
   ```go
   case tea.WindowSizeMsg:
       m.width = msg.Width
       m.height = msg.Height
   ```

### ❌ DON'T

1. **Don't forget solid background** - causes transparency issues
2. **Don't manually center with ANSI** - defeats the purpose of overlay
3. **Don't stack overlays** - hide one before showing another
4. **Don't use large X/Y offsets** - breaks centering
5. **Don't skip WindowSizeMsg** - breaks on terminal resize
6. **Don't constrain Huh form height** - use natural height (0)

---

## Troubleshooting

### Issue: Background shows through modal (transparency)

**Cause**: Modal View() doesn't set background color  
**Solution**: Add `Background(styles.ColorBackground)`:
```go
return lipgloss.NewStyle().
    Background(styles.ColorBackground). // Add this!
    Border(lipgloss.RoundedBorder()).
    Render(content)
```

### Issue: Modal not centered correctly

**Cause**: Wrong position parameters  
**Solution**: Use `overlay.Center` for both X and Y:
```go
overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
```

### Issue: Modal overlaps footer

**Cause**: No Y offset  
**Solution**: Use Y offset of -2:
```go
overlay.New(modal, bg, overlay.Center, overlay.Center, 0, -2)
//                                                         ^^ This!
```

### Issue: Modal doesn't update on terminal resize

**Cause**: Not handling WindowSizeMsg  
**Solution**: Handle it in Update():
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    return m, nil
```

### Issue: Huh form doesn't scroll

**Cause**: Constrained form height  
**Solution**: Use natural height (0):
```go
form := forms.NewYourFormWithDimensions(data, width, 0) // 0 = natural
```

### Issue: Overlay flickers or has artifacts

**Cause**: Usually caused by stacking updates or improper refresh  
**Solution**: 
1. Check that modal returns empty string when not visible
2. Ensure only one modal visible at a time
3. Clear modal reference after hiding

---

## Real-World Examples

### Example 1: View Event Detail Modal

**File**: `internal/cli/components/view_event_detail_modal.go`  
**Lines**: 151  
**Use Case**: Display event details over timeline

**Features**:
- Read-only display
- Action tracking (edit, delete, close)
- Reuses existing RenderEventDetailCard
- Solid background
- Responsive sizing

**Integration**: `internal/cli/intents/browse_timeline_intent.go`
- Shows modal when user presses Enter on event
- Handles edit/delete actions
- Preserves timeline context

### Example 2: Quick Add Event Modal

**File**: `internal/cli/components/quick_add_event_modal.go`  
**Lines**: 250  
**Use Case**: Quick event entry over timeline

**Features**:
- Huh form integration
- Natural height for scrolling
- Submit/cancel handling
- Form data extraction

**Integration**: `internal/cli/intents/browse_timeline_intent.go`
- Shows modal when user presses 'a'
- Saves event on submit
- Refreshes timeline after save

### Example 3: Delete Confirmation Modal

**File**: `internal/cli/components/delete_confirm_modal.go`  
**Lines**: 200  
**Use Case**: Confirm destructive actions

**Features**:
- Simple yes/no confirmation
- Red border for warnings
- Clear button indicators
- Terminal bell alert

**Integration**: `internal/cli/intents/browse_timeline_intent.go`
- Shows modal when user presses 'd' on event
- Deletes event if confirmed
- Updates timeline after deletion

### All 5 Browse Timeline Modals

| Modal | Trigger | Purpose | Lines |
|-------|---------|---------|-------|
| ViewEventDetailModal | Enter | View details | 151 |
| QuickAddEventModal | 'a' | Add event | 250 |
| EditEventModal | 'e' | Edit event | 280 |
| DeleteConfirmModal | 'd' | Confirm delete | 200 |
| FilterModalModel | 'f' | Filter/sort | 350 |

**All use the same pattern**:
- Implement tea.Model
- Solid background in View()
- Wrapped with bubbletea-overlay
- Y offset of -2
- Handle WindowSizeMsg

---

## Related Documentation

- [MODAL_PATTERNS.md](./MODAL_PATTERNS.md) - Modal usage patterns
- [TUI_DEVELOPER_GUIDE.md](./TUI_DEVELOPER_GUIDE.md) - General TUI development
- [VIEW_DETAIL_MODAL_SUMMARY.md](../VIEW_DETAIL_MODAL_SUMMARY.md) - Complete implementation example
- [bubbletea-overlay GitHub](https://github.com/rmhubbert/bubbletea-overlay) - Library source

---

## Summary

`bubbletea-overlay` is the foundation of KaRiya's modal system:

✅ **Reliable**: No flicker, artifacts, or rendering issues  
✅ **Simple**: 5-parameter API, automatic positioning  
✅ **Type-Safe**: Works with any tea.Model  
✅ **Production-Ready**: Used in all 5 Browse Timeline modals  
✅ **Well-Documented**: Complete examples and troubleshooting

**Key Takeaway**: Always use solid backgrounds, handle WindowSizeMsg, and use Y offset of -2.

---

**Questions? Issues?** Check [MODAL_PATTERNS.md](./MODAL_PATTERNS.md) or open an issue on GitHub.
