# StandardView Developer Guide

**Last Updated**: 2026-01-20  
**Author**: KaRiya Development Team  
**Status**: LEGACY - Migrating to UIKit

> **DEPRECATION NOTICE**: `StandardView` is being replaced by `layout.NewScreenLayout()`.
> For new code, use `layout.NewScreenLayout()` from `internal/cli/uikit/layout/`.
> See [UIKIT_GUIDE.md](./UIKIT_GUIDE.md) for the current component library.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Basic Usage](#basic-usage)
4. [Breadcrumbs](#breadcrumbs)
5. [Content](#content)
6. [Modals](#modals)
7. [Help Footer](#help-footer)
8. [Terminal Awareness](#terminal-awareness)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

---

## Overview

### What is StandardView?

StandardView is a component that provides a **standardized, consistent layout** for all TUI screens in KaRiya. It ensures every screen has:

- **Logo at the top** with configurable spacing
- **Breadcrumbs** for navigation context
- **Centered content** area
- **Modal overlays** for feedback (errors, loading, progress, success)
- **Help footer** with visual separator
- **Full terminal width** utilization

### Why Use StandardView?

**Consistency**: Every screen looks and feels the same, providing a professional user experience.

**Productivity**: Reduces boilerplate - no need to manually layout logo, headers, footers.

**Maintainability**: Changes to the standard layout automatically apply everywhere.

**Terminal Awareness**: Automatically handles different terminal sizes and resizing.

---

## Architecture

### Component Structure

```
┌────────────────────────────────────────┐
│         (Logo Spacing)                 │
│                                        │
├────────────────────────────────────────┤
│                                        │
│          ██╗  ██╗ █████╗               │
│          ██║ ██╔╝██╔══██╗              │
│          █████╔╝ ███████║              │  ← Logo
│          ██╔═██╗ ██╔══██║              │
│          ██║  ██╗██║  ██║              │
│          ╚═╝  ╚═╝╚═╝  ╚═╝              │
│                                        │
├────────────────────────────────────────┤
│  Home > Settings > Profile             │  ← Breadcrumbs (optional)
│                                        │
├────────────────────────────────────────┤
│                                        │
│         Your Content Here              │  ← Content Area
│                                        │
│                                        │
├────────────────────────────────────────┤
│ ────────────────────────────────────── │  ← Footer Separator (optional)
│ ↑/k Up  ↓/j Down  Enter Select  q Quit│  ← Help Footer
└────────────────────────────────────────┘

       ┌─────────────────┐
       │  Modal Overlay  │  ← Modal (optional, overlays content)
       └─────────────────┘
```

### Key Components

| Component | Purpose | Required |
|-----------|---------|----------|
| Logo | Branding, visual anchor | Optional |
| Breadcrumbs | Navigation context | Optional |
| Title/Subtitle | Page heading | Optional |
| Content | Main content area | Yes |
| Footer Separator | Visual divider | Optional |
| Help Text | User guidance | Recommended |
| Modal | Overlays for feedback | Optional |

---

## Basic Usage

### Creating a StandardView

```go
import (
    "github.com/baphled/kariya/internal/cli/components"
    "github.com/baphled/kariya/internal/cli/terminal"
)

func (i *YourIntent) View() string {
    termInfo := &terminal.Info{Width: 120, Height: 40}
    logo := components.NewASCIILogo(false, termInfo.Width)
    
    view := components.NewStandardView(termInfo).
        WithLogo(logo, 2).
        WithContent("Your content here").
        WithHelp("q Quit  h Help")
    
    return view.Render()
}
```

### Builder Pattern

StandardView uses a **fluent builder pattern** - all methods return `*StandardView` for chaining:

```go
view := components.NewStandardView(termInfo).
    WithLogo(logo, 2).                      // Add logo with 2-line spacing
    WithBreadcrumbs("Home", "Settings").    // Add breadcrumbs
    WithContent("Content").                 // Set content
    WithHelp("q Quit").                     // Add help text
    WithFooterSeparator(true).              // Show separator
    SetUseFullWidth(true)                   // Use full terminal width
```

### Minimal Example

The simplest possible StandardView:

```go
view := components.NewStandardView(termInfo).
    WithContent("Hello, World!").
    WithHelp("q Quit")

output := view.Render()
```

---

## Breadcrumbs

Breadcrumbs show the user's navigation context.

### Adding Breadcrumbs

```go
view := components.NewStandardView(termInfo).
    WithBreadcrumbs("Main Menu", "Settings", "Display")
```

**Rendered as**: `Main Menu > Settings > Display`

### Dynamic Breadcrumbs

```go
func (i *YourIntent) getBreadcrumbs() []string {
    switch i.state {
    case StateMenu:
        return []string{"Main Menu"}
    case StateSettings:
        return []string{"Main Menu", "Settings"}
    case StateProfile:
        return []string{"Main Menu", "Settings", "Profile"}
    default:
        return []string{"Main Menu"}
    }
}

view := components.NewStandardView(termInfo).
    WithBreadcrumbs(i.getBreadcrumbs()...)
```

### Best Practices

- ✅ Keep breadcrumbs concise (3-5 levels max)
- ✅ Use present tense ("Settings" not "Configure Settings")
- ✅ Match menu item names for consistency
- ❌ Don't include the current page twice (breadcrumbs + title)
- ❌ Don't use very long names (they may truncate)

---

## Content

The content area is where your main UI lives.

### Setting Content

```go
content := "Line 1\nLine 2\nLine 3"
view := components.NewStandardView(termInfo).
    WithContent(content)
```

### Styling Content

Apply Lipgloss styles to content:

```go
import "github.com/charmbracelet/lipgloss"

style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).
    Bold(true)

view := components.NewStandardView(termInfo).
    WithContent("Styled content").
    WithContentStyle(style)
```

### Content from Components

Use other components to build content:

```go
// Using a list
listContent := myListComponent.View()

// Using a form
formContent := myFormComponent.View()

// Using a table
tableContent := myTableComponent.View()

view := components.NewStandardView(termInfo).
    WithContent(listContent)  // or formContent, or tableContent
```

### Long Content

StandardView handles long content gracefully:

```go
longContent := strings.Repeat("Line of text\n", 100)
view := components.NewStandardView(termInfo).
    WithContent(longContent)
// Content will be rendered; scrolling is handled by your terminal
```

---

## Modals

Modals overlay the main content to provide feedback.

### Modal Types

1. **Error** - Show errors with dismiss option
2. **Loading** - Indeterminate progress with optional cancellation
3. **Progress** - Determinate progress (0.0 to 1.0)
4. **Success** - Success feedback with auto-dismiss
5. **Warning** - Warnings that need acknowledgment

### Error Modal

```go
modal := components.NewErrorModal("Error", "Failed to save file: permission denied")

view := components.NewStandardView(termInfo).
    WithContent("Background content").
    ShowModalOverlay(modal)
```

Features:
- Terminal bell alert (accessibility)
- Escape key to dismiss
- No auto-dismiss (requires acknowledgment)

### Loading Modal

```go
modal := components.NewLoadingModal("Processing your request...", false)

view := components.NewStandardView(termInfo).
    WithContent("Background content").
    ShowModalOverlay(modal)
```

With cancellation:

```go
modal := components.NewLoadingModal("Processing...", true)  // cancellable = true
```

### Progress Modal

```go
// progress: 0.0 (0%) to 1.0 (100%)
modal := components.NewProgressModal("Processing", "Analyzing data...", 0.75)

view := components.NewStandardView(termInfo).
    WithContent("Background content").
    ShowModalOverlay(modal)
```

Updating progress:

```go
modal.UpdateProgress(0.25)  // 25%
modal.UpdateProgress(0.50)  // 50%
modal.UpdateProgress(1.00)  // 100%
```

### Success Modal

```go
modal := components.NewSuccessModal("File saved successfully!")

view := components.NewStandardView(termInfo).
    WithContent("Background content").
    ShowModalOverlay(modal)
```

Features:
- Auto-dismisses after 3 seconds
- Green success icon
- Positive feedback

### Warning Modal

```go
modal := components.NewWarningModal("Warning", "This action cannot be undone")

view := components.NewStandardView(termInfo).
    WithContent("Background content").
    ShowModalOverlay(modal)
```

### Modal with Loading Messages

For long operations, rotate messages:

```go
messages := []string{
    "Analyzing events...",
    "Calculating metrics...",
    "Generating bullets...",
    "Formatting document...",
}
rotator := components.NewLoadingMessageRotator(messages, 2*time.Second)

modal := components.NewLoadingModal("Generating CV", false)
modal.SetMessageRotator(rotator)

view := components.NewStandardView(termInfo).
    ShowModalOverlay(modal)
```

See [MODAL_PATTERNS.md](MODAL_PATTERNS.md) for more modal examples.

---

## Help Footer

The help footer provides context-aware user guidance.

### Basic Help

```go
view := components.NewStandardView(termInfo).
    WithContent("Content").
    WithHelp("q Quit  h Help  m Menu")
```

### State-Specific Help

```go
func (i *YourIntent) getHelp() string {
    switch i.state {
    case StateList:
        return "↑/k Up  ↓/j Down  Enter Select  q Quit"
    case StateForm:
        return "Tab Next  Shift+Tab Prev  Enter Submit  Esc Cancel"
    case StateConfirm:
        return "y/Enter Yes  n/Esc No  q Quit"
    default:
        return "q Quit  h Help"
    }
}

view := components.NewStandardView(termInfo).
    WithHelp(i.getHelp())
```

### Footer Separator

Add a visual separator above the help:

```go
view := components.NewStandardView(termInfo).
    WithContent("Content").
    WithHelp("q Quit").
    WithFooterSeparator(true)
```

---

## Terminal Awareness

StandardView adapts to different terminal sizes.

### Setting Terminal Info

```go
termInfo := &terminal.Info{Width: 120, Height: 40}
view := components.NewStandardView(termInfo)
```

### Handling Nil Terminal Info

StandardView gracefully handles nil with defaults:

```go
view := components.NewStandardView(nil)  // Uses 120x40 default
```

### Responsive Rendering

StandardView automatically adjusts to terminal size:

```go
// Small terminal (80x24)
smallTerm := &terminal.Info{Width: 80, Height: 24}
view1 := components.NewStandardView(smallTerm)

// Large terminal (200x60)
largeTerm := &terminal.Info{Width: 200, Height: 60}
view2 := components.NewStandardView(largeTerm)

// Both render appropriately for their size
```

---

## Best Practices

### Do's ✅

1. **Always show help text** - users need guidance
   ```go
   .WithHelp("q Quit  h Help")
   ```

2. **Use breadcrumbs for deep navigation**
   ```go
   .WithBreadcrumbs("Main Menu", "Settings", "Advanced")
   ```

3. **Show the logo consistently** - builds brand recognition
   ```go
   .WithLogo(logo, 2)
   ```

4. **Use modals for important feedback**
   ```go
   .ShowModalOverlay(errorModal)
   ```

5. **Make help context-aware** - show relevant shortcuts
   ```go
   .WithHelp(i.getContextHelp())
   ```

### Don'ts ❌

1. **Don't render without terminal info** (use defaults if needed)
   ```go
   // Bad
   view := NewStandardView(nil)
   
   // Good
   termInfo := &terminal.Info{Width: 120, Height: 40}
   view := NewStandardView(termInfo)
   ```

2. **Don't forget footer separators** - they improve readability
   ```go
   .WithFooterSeparator(true)
   ```

3. **Don't use very long breadcrumbs** - they may truncate
   ```go
   // Bad
   .WithBreadcrumbs("A", "B", "C", "D", "E", "F", "G")
   
   // Good
   .WithBreadcrumbs("Menu", "Settings", "Display")
   ```

4. **Don't show multiple modals simultaneously** - confusing UX
   ```go
   // Bad - only one modal at a time
   .ShowModalOverlay(modal1).ShowModalOverlay(modal2)
   ```

---

## Troubleshooting

### Issue: View is empty

**Cause**: Content not set or terminal info is nil without defaults

**Solution**: Always set content and terminal info
```go
view := components.NewStandardView(termInfo).
    WithContent("Some content")
```

### Issue: Logo not appearing

**Cause**: Logo not set or `WithLogo()` not called

**Solution**: Create and attach logo
```go
logo := components.NewASCIILogo(false, termInfo.Width)
view.WithLogo(logo, 2)
```

### Issue: Modal not visible

**Cause**: Modal not added or fade-in incomplete

**Solution**: Use `ShowModalOverlay()` and wait for fade-in
```go
view.ShowModalOverlay(modal)
```

### Issue: Content cut off

**Cause**: Terminal too small or content exceeds height

**Solution**: StandardView handles this gracefully - ensure terminal is adequate size (minimum 80x24)

### Issue: Footer separator not showing

**Cause**: `WithFooterSeparator()` not called or set to false

**Solution**: Enable footer separator
```go
view.WithFooterSeparator(true)
```

### Issue: Performance is slow

**Cause**: Excessive re-rendering or very large content

**Solution**: 
- Render only when needed
- Use benchmarks to identify bottlenecks
- See performance targets: StandardView < 50ms, Modal < 20ms

---

## Examples

### Complete Intent View

```go
func (i *MyIntent) View() string {
    // Get terminal info from intent
    termInfo := i.GetTerminalInfo()
    
    // Create logo
    logo := i.GetLogo()
    
    // Build view
    view := components.NewStandardView(termInfo).
        WithLogo(logo, 2).
        WithBreadcrumbs(i.getBreadcrumbs()...).
        WithContent(i.getContent()).
        WithHelp(i.getHelp()).
        WithFooterSeparator(true)
    
    // Add modal if needed
    if i.showError {
        view.ShowModalOverlay(i.errorModal)
    }
    
    return view.Render()
}
```

### Multi-State View

```go
func (i *MyIntent) getContent() string {
    switch i.state {
    case StateMenu:
        return i.renderMenu()
    case StateForm:
        return i.renderForm()
    case StateConfirm:
        return i.renderConfirmation()
    default:
        return "Loading..."
    }
}
```

---

## Performance

Based on benchmarks (Phase 5, Task 16):

- **StandardView render**: 0.376ms (target: <50ms) ✓
- **Modal render**: 0.113ms (target: <20ms) ✓
- **Full view render**: 0.727ms (target: <100ms) ✓

All performance targets exceeded by 100x+. See `internal/cli/components/performance_test.go` for details.

---

## Modal Overlays with StandardView

StandardView works seamlessly with modal overlays created using `bubbletea-overlay`.

### Pattern

```go
func (i *YourIntent) View() string {
    // 1. Render base view with StandardView
    baseView := i.CreateViewWithBreadcrumbs("Main", "Section")
    baseView.WithContent(content)
    baseView.WithHelp(footer)
    renderedBase := baseView.Render()
    
    // 2. Overlay modal if visible
    if i.yourModal != nil && i.yourModal.IsVisible() {
        return i.renderModalOverlay(renderedBase)
    }
    
    return renderedBase
}
```

### Key Points

- ✅ StandardView renders the complete base layout
- ✅ Modal overlay is applied to the fully-rendered view
- ✅ Logo, breadcrumbs, and footer remain visible in background
- ✅ Modal appears centered over content area
- ✅ Y offset of -2 prevents footer overlap

**See**: 
- [MODAL_PATTERNS.md](MODAL_PATTERNS.md#modal-overlays-with-bubbletea-overlay) - Complete modal patterns
- [BUBBLETEA_OVERLAY_GUIDE.md](BUBBLETEA_OVERLAY_GUIDE.md) - Library usage guide
- `internal/cli/intents/browse_timeline_intent.go` - Real-world example (5 modals)

---

## Related Documentation

- [MODAL_PATTERNS.md](MODAL_PATTERNS.md) - Modal usage patterns
- [BUBBLETEA_OVERLAY_GUIDE.md](BUBBLETEA_OVERLAY_GUIDE.md) - bubbletea-overlay library guide **NEW!**
- [TUI_DEVELOPER_GUIDE.md](TUI_DEVELOPER_GUIDE.md) - General TUI development
- [TUI_STANDARDS.md](TUI_STANDARDS.md) - TUI design standards
- [LIPGLOSS_BUBBLES_GUIDE.md](LIPGLOSS_BUBBLES_GUIDE.md) - Styling guide

---

**Questions? Issues?** Open an issue on GitHub or check the troubleshooting section above.
