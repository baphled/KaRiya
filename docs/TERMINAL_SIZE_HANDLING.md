---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Terminal Size Handling Guide

**Created**: 2026-01-06
**Status**: Complete and Production-Ready
**Test Coverage**: 47 tests (25 terminal + 22 layout)

---

## Overview

This guide explains how KaRiya handles terminal dimensions and provides responsive, centered layouts across different screen sizes. The system ensures a consistent user experience whether the terminal is in fullscreen, split-screen, or compact mode.

### Key Features

- **Automatic size detection**: Terminal dimensions tracked automatically via BubbleTea
- **Responsive layouts**: Content adapts to 5 different size categories
- **Smart centering**: Intelligent content centering based on available space
- **Intent-aware**: Terminal size propagated to all intents automatically
- **Graceful degradation**: Minimal fallback for very small terminals

---

## Architecture

### Core Components

1. **`terminal.Info`** - Tracks terminal dimensions and provides size utilities
2. **`layout.Manager`** - Calculates responsive layouts and margins
3. **`TerminalAwareIntent`** - Interface for intents that need terminal dimensions
4. **`SmartContainer`** - Intelligent centering component
5. **Intent Router** - Propagates terminal size to active intents

### Data Flow

```
Terminal Resize Event (OS)
    ↓
BubbleTea WindowSizeMsg
    ↓
App Model (app.go)
    ↓
Intent Router
    ↓
Active Intent (if TerminalAwareIntent)
    ↓
Layout/Rendering Components
```

---

## Terminal Size Categories

The system recognizes 5 size categories with different layout strategies:

| Category | Width Range | Use Case | Margins |
|----------|------------|----------|---------|
| **SizeTiny** | < 60 cols | Very small/mobile | 0-1 chars |
| **SizeCompact** | 60-79 cols | Split screen | 1-2 chars |
| **SizeNormal** | 80-119 cols | Standard terminal | 2-4 chars |
| **SizeLarge** | 120-159 cols | Large terminal | 3-8 chars |
| **SizeXLarge** | ≥ 160 cols | Ultra-wide | 3-8 chars |

### Minimum Requirements

- **Minimum width**: 40 columns
- **Minimum height**: 15 rows
- **Default fallback**: 80x24 (when size unknown)

---

## Using Terminal Info

### In Intents

Intents can implement `TerminalAwareIntent` to receive terminal dimensions:

```go
type MyIntent struct {
    intents.BaseIntent  // Embeds terminal handling
    // ... other fields
}

// Implement TerminalAwareIntent interface
func (i *MyIntent) UpdateTerminalInfo(info *terminal.Info) {
    i.BaseIntent.UpdateTerminalInfo(info)
    // Optionally do custom handling
}

func (i *MyIntent) GetMinimumSize() (width, height int) {
    return 60, 20  // Custom minimum size for this intent
}

// Use in View()
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    
    // Check if terminal is large enough
    if !info.CanRender(terminal.DefaultConfig) {
        return "Terminal too small (min: 40x15)"
    }
    
    // Get size category for responsive rendering
    switch info.GetCategory() {
    case terminal.SizeTiny:
        return i.renderCompact()
    default:
        return i.renderNormal()
    }
}
```

### Automatic Propagation

The Intent Router automatically propagates terminal size to all `TerminalAwareIntent` implementations:

- On intent activation
- On terminal resize (WindowSizeMsg)
- When navigating back to a previous intent

No manual wiring required!

---

## Using Layout Manager

The `layout.Manager` provides responsive layout calculations:

```go
import (
    "github.com/baphled/kariya/internal/cli/layout"
    "github.com/baphled/kariya/internal/cli/terminal"
)

func renderMyView(info *terminal.Info) string {
    manager := layout.NewManager(info)
    
    // Get content area (terminal size minus margins)
    content := manager.GetContentArea()
    // content.Width, content.Height, content.X, content.Y
    
    // Get responsive margins
    margins := manager.GetMargins()
    
    // Calculate column layout
    columns := manager.CalculateColumns(3, 2)  // 3 columns, 2 char gutter
    
    // Check if should use compact layout
    if manager.ShouldUseCompactLayout() {
        return renderCompact()
    }
    
    return renderNormal()
}
```

### Layout Helpers

```go
// Check if multi-column layouts will work
if !manager.ShouldUseListLayout() {
    // Render as grid
} else {
    // Render as single column list
}

// Calculate equal columns with gutters
columns := manager.CalculateColumns(count, gutterWidth)
// Returns []int of column widths that sum to content width

// Update terminal info dynamically
manager.UpdateTerminalInfo(newInfo)
```

---

## Using SmartContainer

The `SmartContainer` provides intelligent centering and responsive rendering:

```go
import "github.com/baphled/kariya/internal/cli/components"

func renderCentered(info *terminal.Info, content string) string {
    container := components.NewSmartContainer(info)
    
    return container.
        SetContent(content).
        SetCenteringMode(components.CenterBoth).
        SetMinContentSize(40, 10).
        SetMaxContentSize(120, 40).
        Render()
}
```

### Centering Modes

```go
// No centering
container.SetCenteringMode(components.CenterNone)

// Horizontal only (default for many views)
container.SetCenteringMode(components.CenterHorizontal)

// Vertical only
container.SetCenteringMode(components.CenterVertical)

// Both directions (default for menus)
container.SetCenteringMode(components.CenterBoth)
```

### Advanced Options

```go
// Custom margins (overrides responsive defaults)
container.SetCustomMargins(components.Margins{
    Top: 3, Right: 6, Bottom: 3, Left: 6,
})

// Enforce size constraints
container.
    SetMinContentSize(50, 15).  // Minimum content dimensions
    SetMaxContentSize(100, 40)  // Maximum content dimensions

// Overflow handling (future enhancement)
container.SetOverflowMode(components.OverflowWrap)
```

---

## Best Practices

### 1. Always Check Terminal Validity

```go
info := i.GetTerminalInfo()

if !info.IsValid {
    // Use defaults, terminal size not yet received
    width, height := info.GetSafeDimensions(terminal.DefaultConfig)
}

if !info.CanRender(terminal.DefaultConfig) {
    // Terminal too small, show minimal message
    return "Terminal too small"
}
```

### 2. Provide Responsive Layouts

```go
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    
    switch info.GetCategory() {
    case terminal.SizeTiny, terminal.SizeCompact:
        return i.renderCompact()
    default:
        return i.renderNormal()
    }
}
```

### 3. Use SmartContainer for Centering

```go
// DON'T manually calculate centering
func badCentering(content string, width int) string {
    padding := (width - len(content)) / 2  // Breaks with ANSI codes!
    return strings.Repeat(" ", padding) + content
}

// DO use SmartContainer
func goodCentering(info *terminal.Info, content string) string {
    return components.NewSmartContainer(info).
        SetContent(content).
        SetCenteringMode(components.CenterHorizontal).
        Render()
}
```

### 4. Test Multiple Sizes

```go
var _ = Describe("MyIntent Responsive Rendering", func() {
    It("should render at tiny size", func() {
        info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
        view := intent.View()
        Expect(view).NotTo(BeEmpty())
    })
    
    It("should render at normal size", func() {
        info.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
        view := intent.View()
        Expect(view).NotTo(BeEmpty())
    })
    
    It("should render at large size", func() {
        info.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
        view := intent.View()
        Expect(view).NotTo(BeEmpty())
    })
})
```

---

## Migration Guide

### Converting Existing Intents

**Step 1**: Embed `BaseIntent`

```go
type MyIntent struct {
    intents.BaseIntent  // Add this
    // existing fields...
}
```

**Step 2**: Use terminal info in rendering

```go
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    manager := layout.NewManager(info)
    
    content := manager.GetContentArea()
    // Use content.Width, content.Height for layout
}
```

**Step 3**: Replace manual centering with SmartContainer

```go
// Before
func (i *MyIntent) View() string {
    content := i.renderContent()
    return components.CenterBlock(content, 80, 24, true)
}

// After
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    content := i.renderContent()
    
    return components.NewSmartContainer(info).
        SetContent(content).
        SetCenteringMode(components.CenterBoth).
        Render()
}
```

### Converting Components

Components that need terminal awareness can accept `*terminal.Info` as a parameter:

```go
// Before
func NewMyComponent(width, height int) *MyComponent {
    return &MyComponent{width: width, height: height}
}

// After
func NewMyComponent(info *terminal.Info) *MyComponent {
    width, height := info.GetSafeDimensions(terminal.DefaultConfig)
    return &MyComponent{
        terminalInfo: info,
        width:  width,
        height: height,
    }
}

// Update on resize
func (c *MyComponent) Update(msg tea.Msg) (*MyComponent, tea.Cmd) {
    if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
        c.terminalInfo.Update(wsMsg)
        c.width, c.height = c.terminalInfo.GetSafeDimensions(terminal.DefaultConfig)
    }
    // ...
}
```

---

## Common Patterns

### Pattern 1: Responsive Margins

```go
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    margins := i.getResponsiveMargins(info)
    
    width, height := info.ContentArea(margins)
    // Render with calculated dimensions
}

func (i *MyIntent) getResponsiveMargins(info *terminal.Info) terminal.Margins {
    switch info.GetCategory() {
    case terminal.SizeTiny:
        return terminal.Margins{0, 1, 0, 1}
    case terminal.SizeCompact:
        return terminal.Margins{1, 2, 1, 2}
    default:
        return terminal.Margins{2, 4, 2, 4}
    }
}
```

### Pattern 2: Adaptive Column Layout

```go
func (i *MyIntent) View() string {
    manager := layout.NewManager(i.GetTerminalInfo())
    
    var columns []int
    if manager.ShouldUseListLayout() {
        // Single column for narrow terminals
        columns = []int{manager.GetContentArea().Width}
    } else {
        // Multi-column for wider terminals
        columns = manager.CalculateColumns(3, 2)
    }
    
    return i.renderColumns(columns)
}
```

### Pattern 3: Size-Based Feature Toggle

```go
func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    
    // Show simplified version on small screens
    if info.GetCategory() <= terminal.SizeCompact {
        return i.renderSimplified()
    }
    
    // Full version with all features
    return i.renderFull()
}
```

---

## Troubleshooting

### Issue: Terminal size is (0, 0)

**Cause**: `WindowSizeMsg` not yet received

**Solution**: Always check `IsValid` and use `GetSafeDimensions()`

```go
if !info.IsValid {
    width, height := info.GetSafeDimensions(terminal.DefaultConfig)
    // Use fallback dimensions
}
```

### Issue: Content not centered

**Cause**: Using manual centering with ANSI codes

**Solution**: Use `SmartContainer` which handles ANSI codes correctly

```go
// Use lipgloss.Width() for proper width calculation
visualWidth := lipgloss.Width(line)
```

### Issue: Layout breaks on resize

**Cause**: Not updating components on `WindowSizeMsg`

**Solution**: Ensure intent implements `TerminalAwareIntent` or manually handles resize

```go
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
        i.terminalInfo.Update(wsMsg)
        // Recalculate layout...
    }
    // ...
}
```

### Issue: Windows doesn't support resize

**Limitation**: Windows doesn't send SIGWINCH signals

**Workaround**: App must be restarted after resizing on Windows

---

## Testing

### Test Utilities

```go
import "github.com/baphled/kariya/internal/cli/terminal"

// Common test sizes
var (
    TinyTerminal    = tea.WindowSizeMsg{Width: 40, Height: 15}
    CompactTerminal = tea.WindowSizeMsg{Width: 60, Height: 20}
    NormalTerminal  = tea.WindowSizeMsg{Width: 80, Height: 24}
    LargeTerminal   = tea.WindowSizeMsg{Width: 120, Height: 40}
)

// Test at multiple sizes
func testResponsiveRendering() {
    sizes := []tea.WindowSizeMsg{
        TinyTerminal,
        CompactTerminal,
        NormalTerminal,
        LargeTerminal,
    }
    
    for _, size := range sizes {
        info.Update(size)
        view := intent.View()
        Expect(view).NotTo(BeEmpty())
    }
}
```

### Test Checklist

- [ ] Renders with default (uninitialized) terminal info
- [ ] Renders at tiny size (50x20)
- [ ] Renders at compact size (70x24)
- [ ] Renders at normal size (80x24)
- [ ] Renders at large size (120x40)
- [ ] Handles resize during operation
- [ ] Falls back gracefully when too small

---

## Performance Considerations

### Resize Throttling

Currently, all `WindowSizeMsg` events are processed immediately. For very frequent resizes (e.g., dragging window), consider throttling:

```go
// Future enhancement: throttle resize events
const resizeThrottleMs = 100
```

### Content Caching

Cache rendered content when terminal size hasn't changed:

```go
type MyIntent struct {
    lastWidth  int
    lastHeight int
    cachedView string
}

func (i *MyIntent) View() string {
    info := i.GetTerminalInfo()
    
    if info.Width == i.lastWidth && info.Height == i.lastHeight {
        return i.cachedView
    }
    
    i.cachedView = i.render()
    i.lastWidth = info.Width
    i.lastHeight = info.Height
    
    return i.cachedView
}
```

---

## Future Enhancements

### Planned Features

1. **Resize throttling**: Debounce rapid resize events
2. **Scroll support**: Handle content overflow with scrolling
3. **Breakpoint configuration**: Custom size categories per intent
4. **Layout templates**: Pre-defined responsive layouts
5. **Orientation detection**: Detect portrait vs landscape
6. **Multi-screen support**: Handle terminal spanning multiple monitors

### API Stability

The current API is considered stable for production use. Any breaking changes will be versioned and documented.

---

## Summary

The KaRiya terminal size handling system provides:

✅ **Automatic size detection and propagation**
✅ **Responsive layouts with 5 size categories**
✅ **Smart centering that handles ANSI codes**
✅ **Intent-aware architecture**
✅ **Comprehensive testing (47 tests)**
✅ **Production-ready with graceful degradation**

For questions or issues, please refer to the [TUI Developer Guide](TUI_DEVELOPER_GUIDE.md) or open an issue on GitHub.
