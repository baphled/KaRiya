# Modal Overlay Pattern

**Last Updated**: 2026-01-13
**Status**: ✅ **PRODUCTION STANDARD - REQUIRED PATTERN**

---

## Table of Contents

1. [Overview](#overview)
2. [The Critical Rule](#the-critical-rule)
3. [Why This Matters](#why-this-matters)
4. [The Correct Pattern](#the-correct-pattern)
5. [Common Anti-Patterns](#common-anti-patterns)
6. [Implementation Guide](#implementation-guide)
7. [Testing Modal Overlays](#testing-modal-overlays)
8. [Examples](#examples)
9. [Troubleshooting](#troubleshooting)

---

## Overview

Modal overlays in KaRiya TUI **MUST** follow a specific rendering order to avoid alignment issues, content shifting, and visual corruption. This document defines the mandatory pattern for all intent implementations.

**Applies To**:
- All intent View() methods with modals
- Filter modals, edit modals, confirmation modals
- Any UI overlay that appears on top of existing content

**Key Principle**: Render the complete base view FIRST, then overlay the modal as the FINAL step.

---

## The Critical Rule

> **ALWAYS render StandardView FIRST (complete view), then overlay modal LAST**

```go
// ✅ CORRECT PATTERN
func (i *YourIntent) View() string {
    // 1. Create StandardView with ALL content
    view := i.CreateViewWithBreadcrumbs(...)
    view.WithContent(screen.RenderContent())
    view.WithHelp(i.getContextHelp())
    
    // 2. Render COMPLETE base view
    baseView := view.Render()
    
    // 3. Overlay modal as FINAL step (if visible)
    if i.modal != nil && i.modal.IsVisible() {
        return i.renderModalOverlay(baseView)
    }
    
    return baseView
}
```

**Why**: StandardView uses `lipgloss.Place()` to center content in the terminal. If you overlay the modal on raw content BEFORE StandardView rendering, you get double-centering and misalignment.

---

## Why This Matters

### Problem: Rendering Modal Before StandardView

```go
// ❌ WRONG - Leads to misalignment
func (i *YourIntent) View() string {
    content := screen.RenderContent()
    
    // Overlay modal on raw content
    if i.modal != nil && i.modal.IsVisible() {
        content = i.renderModalOverlay(content)
    }
    
    // Then pass to StandardView
    view := i.CreateViewWithBreadcrumbs(...)
    view.WithContent(content)  // Double-centering happens here!
    return view.Render()
}
```

**What Happens**:
1. Modal is centered on raw content (first centering)
2. StandardView centers the entire result again (second centering)
3. Modal appears off-center or shifted
4. Background content may be corrupted or duplicated
5. Terminal resizing causes visual glitches

### Solution: Rendering Modal After StandardView

```go
// ✅ CORRECT - Perfect alignment
func (i *YourIntent) View() string {
    // 1. Create complete view
    view := i.CreateViewWithBreadcrumbs(...)
    view.WithContent(screen.RenderContent())
    view.WithHelp(i.getContextHelp())
    
    // 2. Render once (with proper centering)
    baseView := view.Render()
    
    // 3. Overlay modal on complete, centered view
    if i.modal != nil && i.modal.IsVisible() {
        return i.renderModalOverlay(baseView)
    }
    
    return baseView
}
```

**What Happens**:
1. StandardView renders complete, properly centered view
2. Modal overlay replaces specific lines where modal appears
3. Single centering pass (StandardView only)
4. Modal appears exactly where intended
5. Background content remains intact

---

## The Correct Pattern

### Step-by-Step Implementation

#### Step 1: Create StandardView with Complete Content

```go
view := i.CreateViewWithBreadcrumbs(
    i.GetState(),
    i.terminal.Width,
    i.terminal.Height,
)

// Add ALL content that would appear without modal
view.WithContent(screen.RenderContent())
view.WithHelp(i.getContextHelp())
```

**Key Points**:
- Use actual screen content, not placeholder
- Include full help footer
- Set all view properties before rendering

#### Step 2: Render Complete Base View

```go
baseView := view.Render()
```

**Result**: A fully-rendered, properly-centered view that fills the terminal. This is what the user would see if NO modal were displayed.

#### Step 3: Conditionally Overlay Modal

```go
if i.modal != nil && i.modal.IsVisible() {
    return i.renderModalOverlay(baseView)
}

return baseView
```

**Key Points**:
- Check both modal existence AND visibility
- Only modify the already-rendered view
- Return the overlaid result directly

### Complete Template

```go
func (i *YourIntent) View() string {
    // Handle different screen types
    switch screen := i.currentScreen.(type) {
    case *YourListScreen:
        // 1. Create StandardView
        view := i.CreateViewWithBreadcrumbs(
            i.GetState(),
            i.terminal.Width,
            i.terminal.Height,
        )
        view.WithContent(screen.RenderContent())
        view.WithHelp(i.getContextHelp())
        
        // 2. Render COMPLETE view
        baseView := view.Render()
        
        // 3. Overlay modal (if visible)
        if i.filterModal != nil && i.filterModal.IsVisible() {
            return i.renderFilterModalOverlay(baseView)
        }
        
        return baseView
        
    case *YourDetailScreen:
        // Same pattern for other screens
        view := i.CreateViewWithBreadcrumbs(...)
        view.WithContent(screen.RenderContent())
        view.WithHelp(i.getContextHelp())
        baseView := view.Render()
        
        if i.editModal != nil && i.editModal.IsVisible() {
            return i.renderEditModalOverlay(baseView)
        }
        
        return baseView
        
    default:
        return "Unknown screen"
    }
}
```

---

## Common Anti-Patterns

### Anti-Pattern 1: Overlaying Before StandardView

```go
// ❌ WRONG
content := screen.RenderContent()
if modal.IsVisible() {
    content = overlayModal(content)  // Modal on raw content
}
view.WithContent(content)  // Then StandardView centers it
return view.Render()
```

**Problem**: Double-centering, misalignment

### Anti-Pattern 2: Conditional View Creation

```go
// ❌ WRONG
if modal.IsVisible() {
    view.WithContent(modalContent)
} else {
    view.WithContent(screenContent)
}
return view.Render()
```

**Problem**: Screen content not preserved, modal doesn't overlay

### Anti-Pattern 3: Splitting View Rendering

```go
// ❌ WRONG
header := view.RenderHeader()
content := screen.RenderContent()
if modal.IsVisible() {
    content = overlayModal(content)
}
footer := view.RenderFooter()
return header + content + footer
```

**Problem**: StandardView centering not applied, manual concatenation error-prone

### Anti-Pattern 4: Modifying Screen Content Directly

```go
// ❌ WRONG
if modal.IsVisible() {
    screen.SetContent(modal.Render())
}
view.WithContent(screen.RenderContent())
return view.Render()
```

**Problem**: Screen state corrupted, not a true overlay

---

## Implementation Guide

### Creating the Modal Overlay Helper

```go
func (i *YourIntent) renderModalOverlay(baseView string) string {
    if i.modal == nil || !i.modal.IsVisible() {
        return baseView
    }
    
    // Get modal content
    modalContent := i.modal.View()
    
    // Use ModalContent's RenderOverlay for proper positioning
    return components.RenderOverlay(
        baseView,
        modalContent,
        i.terminal.Width,
        i.terminal.Height,
    )
}
```

### The RenderOverlay Implementation

Located in `internal/cli/components/modal.go`:

```go
func RenderOverlay(baseView, modalContent string, termWidth, termHeight int) string {
    baseLines := strings.Split(baseView, "\n")
    
    // Determine modal dimensions
    modalLines := strings.Split(strings.TrimSpace(modalContent), "\n")
    modalHeight := len(modalLines)
    modalWidth := 0
    for _, line := range modalLines {
        w := lipgloss.Width(line)
        if w > modalWidth {
            modalWidth = w
        }
    }
    
    // Calculate centered position
    startLine := (termHeight - modalHeight) / 2
    if startLine < 0 {
        startLine = 0
    }
    
    // Replace lines where modal appears
    for i := 0; i < modalHeight && (startLine+i) < len(baseLines); i++ {
        modalLine := ""
        if i < len(modalLines) {
            modalLine = modalLines[i]
        }
        
        // Center the modal line horizontally
        centeredLine := lipgloss.PlaceHorizontal(
            termWidth,
            lipgloss.Center,
            modalLine,
        )
        
        baseLines[startLine+i] = centeredLine
    }
    
    return strings.Join(baseLines, "\n")
}
```

**Key Algorithm**:
1. Split base view into lines
2. Calculate modal dimensions and center position
3. Replace lines where modal should appear
4. Use `lipgloss.PlaceHorizontal` to center each modal line
5. Rejoin lines into final view

**Why This Works**:
- Base view is already properly rendered and centered by StandardView
- We only replace the specific lines where the modal appears
- Each modal line is individually centered horizontally
- No complex splicing or background preservation needed

---

## Testing Modal Overlays

### Unit Test Template

```go
var _ = Describe("Modal Overlay", func() {
    var (
        intent *YourIntent
        modal  *components.FilterModal
    )
    
    BeforeEach(func() {
        intent = NewYourIntent(context)
        intent.terminal = &cli.TerminalInfo{Width: 120, Height: 40}
        modal = components.NewFilterModal(...)
        intent.filterModal = modal
    })
    
    It("should render base view without modal when not visible", func() {
        modal.Hide()
        
        view := intent.View()
        
        Expect(view).NotTo(ContainSubstring("Filter"))
        Expect(view).To(ContainSubstring("Expected Screen Content"))
    })
    
    It("should overlay modal on complete view when visible", func() {
        modal.Show()
        
        view := intent.View()
        
        // Should contain both modal content AND background
        Expect(view).To(ContainSubstring("Filter"))
        Expect(view).To(ContainSubstring("Expected Screen Content"))
    })
    
    It("should center modal horizontally", func() {
        modal.Show()
        intent.terminal.Width = 120
        
        view := intent.View()
        lines := strings.Split(view, "\n")
        
        // Find modal line
        for _, line := range lines {
            if strings.Contains(line, "Filter") {
                // Count leading/trailing spaces
                trimmed := strings.TrimSpace(line)
                totalWidth := lipgloss.Width(line)
                contentWidth := lipgloss.Width(trimmed)
                leftPadding := strings.Index(line, trimmed)
                rightPadding := totalWidth - leftPadding - contentWidth
                
                // Should be roughly centered (allow ±1 for odd widths)
                Expect(leftPadding).To(BeNumerically("~", rightPadding, 1))
            }
        }
    })
    
    It("should preserve background content outside modal area", func() {
        modal.Show()
        
        view := intent.View()
        
        // Logo should still be visible (at top)
        Expect(view).To(ContainSubstring("KaRiya"))
        
        // Footer should still be visible (at bottom)
        Expect(view).To(ContainSubstring("Esc"))
    })
})
```

### Visual Test Checklist

When testing modal overlays manually:

- [ ] Modal appears centered horizontally
- [ ] Modal appears centered vertically
- [ ] Logo remains visible at top
- [ ] Footer remains visible at bottom
- [ ] Background content not duplicated
- [ ] No visual glitches on initial display
- [ ] Modal updates properly when form fields change
- [ ] Terminal resize doesn't break alignment
- [ ] Modal disappears cleanly when closed
- [ ] Base view returns to normal after modal closes

---

## Examples

### Example 1: BrowseTimeline with Filter Modal

**Reference Implementation**: `internal/cli/intents/browse_timeline_intent.go`

```go
func (i *BrowseTimelineIntent) View() string {
    switch screen := i.currentScreen.(type) {
    case *timeline.TimelineEventListScreen:
        // 1. Create StandardView with complete content
        view := i.CreateViewWithBreadcrumbs(
            i.GetState(),
            i.terminal.Width,
            i.terminal.Height,
        )
        view.WithContent(screen.RenderContent())
        view.WithHelp(i.getContextHelp())
        
        // 2. Render COMPLETE view
        baseView := view.Render()
        
        // 3. Overlay filter modal (if visible)
        if i.filterModal != nil && i.filterModal.IsVisible() {
            return i.renderFilterModalOverlay(baseView)
        }
        
        return baseView
        
    // ... other screen types
    }
}

func (i *BrowseTimelineIntent) renderFilterModalOverlay(baseView string) string {
    if i.filterModal == nil || !i.filterModal.IsVisible() {
        return baseView
    }
    
    modalContent := i.filterModal.View()
    
    return components.RenderOverlay(
        baseView,
        modalContent,
        i.terminal.Width,
        i.terminal.Height,
    )
}
```

### Example 2: ManageSkills with Edit Modal (TODO)

**Target Implementation**: `internal/cli/intents/manage_skills_intent.go`

```go
func (i *ManageSkillsIntent) View() string {
    switch screen := i.currentScreen.(type) {
    case *skills.SkillListScreen:
        // 1. Create StandardView
        view := i.CreateViewWithBreadcrumbs(
            i.GetState(),
            i.terminal.Width,
            i.terminal.Height,
        )
        view.WithContent(screen.RenderContent())
        view.WithHelp(i.getContextHelp())
        
        // 2. Render complete view
        baseView := view.Render()
        
        // 3. Overlay edit modal (if visible)
        if i.editModal != nil && i.editModal.IsVisible() {
            return i.renderEditModalOverlay(baseView)
        }
        
        return baseView
        
    // ... other screens
    }
}
```

---

## Troubleshooting

### Issue: Modal Appears Off-Center

**Symptom**: Modal is shifted left/right or up/down

**Likely Cause**: Overlaying modal before StandardView rendering

**Solution**: Move modal overlay AFTER `view.Render()`:
```go
baseView := view.Render()  // MUST render first
if modal.IsVisible() {
    return renderModalOverlay(baseView)  // Then overlay
}
return baseView
```

### Issue: Background Content Shifts When Modal Opens

**Symptom**: List items or other content moves when modal appears

**Likely Cause**: Modifying screen content instead of overlaying

**Solution**: Ensure screen.RenderContent() is called BEFORE checking modal visibility:
```go
// ✅ CORRECT - Screen renders normally
view.WithContent(screen.RenderContent())
baseView := view.Render()

// ❌ WRONG - Content changes based on modal state
if modal.IsVisible() {
    view.WithContent(modal.View())
} else {
    view.WithContent(screen.RenderContent())
}
```

### Issue: Modal Content Blank Until User Input

**Symptom**: Modal shows but fields are empty until user presses a key

**Likely Cause**: Form's `Init()` method not called

**Solution**: Always call `Init()` when creating form-based modals:
```go
modal := components.NewFilterModal(...)
return modal.Init()  // CRITICAL for immediate rendering
```

### Issue: Background Content Duplicated or Corrupted

**Symptom**: Screen content appears twice or garbled around modal

**Likely Cause**: Complex background preservation in RenderOverlay

**Solution**: Use simple line replacement, not fragment preservation:
```go
// ✅ CORRECT - Replace entire lines
baseLines[startLine+i] = centeredModalLine

// ❌ WRONG - Try to preserve background fragments
prefix := baseLines[startLine+i][:modalStart]
suffix := baseLines[startLine+i][modalEnd:]
baseLines[startLine+i] = prefix + modalLine + suffix
```

### Issue: Terminal Resize Breaks Modal Alignment

**Symptom**: Modal position is wrong after terminal resize

**Likely Cause**: Using cached dimensions instead of current terminal info

**Solution**: Always use `i.terminal.Width` and `i.terminal.Height` from current state:
```go
// ✅ CORRECT - Use current terminal info
return components.RenderOverlay(
    baseView,
    modalContent,
    i.terminal.Width,   // Current width
    i.terminal.Height,  // Current height
)

// ❌ WRONG - Use cached/stale dimensions
return components.RenderOverlay(
    baseView,
    modalContent,
    i.lastWidth,   // May be outdated
    i.lastHeight,
)
```

---

## Related Documentation

- **[THEMED_FOOTER_GUIDE.md](THEMED_FOOTER_GUIDE.md)** - Building consistent modal footers
- **[INTENT_PATTERNS_LIBRARY.md](INTENT_PATTERNS_LIBRARY.md)** - Complete intent patterns catalog
- **[MODAL_PATTERNS.md](../MODAL_PATTERNS.md)** - Modal types and usage patterns
- **[STANDARDVIEW_GUIDE.md](../STANDARDVIEW_GUIDE.md)** - StandardView component usage

---

## Checklist for Modal Overlay Implementation

When implementing or refactoring an intent with modals, verify:

- [ ] StandardView is created with complete content (no placeholders)
- [ ] `view.Render()` is called to get base view BEFORE modal check
- [ ] Modal overlay is the LAST step in View() method
- [ ] `renderModalOverlay()` helper uses `components.RenderOverlay()`
- [ ] Modal visibility is checked with `modal != nil && modal.IsVisible()`
- [ ] Form-based modals call `Init()` when created
- [ ] Current terminal dimensions are used (not cached values)
- [ ] Tests verify modal centering and background preservation
- [ ] Visual testing confirms proper alignment at multiple terminal sizes

---

**Status**: ✅ **MANDATORY PATTERN - REQUIRED FOR ALL INTENTS WITH MODALS**

*Last verified in BrowseTimeline intent (2026-01-13) - Reference implementation*
