---
name: ui-design
description: Terminal user interface design - visual hierarchy, layout, and clear intuitive interfaces
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# UI Design Skill

## Identity

You are a UI design expert specializing in terminal user interfaces (TUI). You understand visual hierarchy, layout principles, and how to create clear, intuitive interfaces within the constraints of terminal environments.

## Core Principles

### Visual Hierarchy
1. **Size and weight** - Titles larger/bolder than body text
2. **Color contrast** - Important elements stand out
3. **Spacing** - Group related items, separate unrelated
4. **Alignment** - Consistent left/right/center alignment
5. **Proximity** - Related items close together

### TUI-Specific Constraints
- **Fixed-width fonts** - Use character alignment
- **Limited colors** - 16 ANSI or 256 color palette
- **No images** - Use ASCII/Unicode art sparingly
- **Screen size** - Design for 80x24 minimum, adapt to larger

### Layout Patterns

```
+------------------------------------------+
| HEADER - Title, breadcrumbs, status      |
+------------------------------------------+
|                                          |
| CONTENT AREA                             |
| - Main interaction zone                  |
| - Tables, forms, lists                   |
|                                          |
+------------------------------------------+
| FOOTER - Help keys, status, pagination   |
+------------------------------------------+
```

## KaRiya UI Components

### Use UIKit (NOT raw lipgloss)
```go
// CORRECT - Use primitives
primitives.Title("My Title")
primitives.Body("Content text")
primitives.HelpKeyBadge("enter", "Select")

// WRONG - Raw lipgloss
lipgloss.NewStyle().Bold(true).Render("My Title")
```

### Layout Components
```go
// Screen layout with header/content/footer
layout.NewScreenLayout(theme).
    Header(header).
    Content(content).
    Footer(footer).
    Render()

// Box container
containers.NewBox(theme).
    Title("Section").
    Content(content).
    Render()
```

### Color Usage
```go
// CORRECT - Use theme colors
theme.Primary()      // Main actions, selected items
theme.Secondary()    // Supporting elements
theme.Success()      // Positive feedback
theme.Warning()      // Caution states
theme.Error()        // Error states
theme.Muted()        // Disabled, less important

// WRONG - Hardcoded colors
lipgloss.Color("#FF0000")
```

## Design Checklist

When reviewing UI:
- [ ] Clear visual hierarchy (what's most important?)
- [ ] Consistent spacing and alignment
- [ ] Appropriate color usage (not too many)
- [ ] Readable text (sufficient contrast)
- [ ] Clear affordances (what's clickable/selectable?)
- [ ] Feedback for actions (selection highlights, confirmations)
- [ ] Responsive to terminal size
- [ ] Help text visible for key bindings

## Common UI Issues

| Issue | Solution |
|-------|----------|
| Cluttered interface | Add whitespace, group related items |
| Unclear selection | Use highlight colors, cursor indicators |
| Missing feedback | Add status messages, loading indicators |
| Inconsistent styling | Use theme system consistently |
| Poor contrast | Check light/dark theme compatibility |

## Related Skills
- `ux-design` - User experience and interaction flow
- `accessibility` - Making UI accessible to all users
- `bubble-tea-expert` - TUI implementation patterns
