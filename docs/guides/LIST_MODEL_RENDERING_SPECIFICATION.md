---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# List Model Rendering Specification

**Version**: 1.0
**Date**: 2026-01-01
**Status**: STANDARDIZED
**Audience**: Developers implementing or modifying list models

## Overview

This specification defines the standard rendering and behavior requirements for all list models in KaRiya. It ensures consistency across different list implementations (Career Events, Bursts, Facts) and provides guidelines for future list model development.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Rendering Pipeline](#rendering-pipeline)
3. [Component Structure](#component-structure)
4. [Visual Specification](#visual-specification)
5. [Navigation Specification](#navigation-specification)
6. [Styling Specification](#styling-specification)
7. [Implementation Checklist](#implementation-checklist)
8. [Code Examples](#code-examples)
9. [Common Patterns](#common-patterns)
10. [Troubleshooting](#troubleshooting)

## Architecture Overview

### Design Principles

1. **Composition Over Inheritance**: Models compose Container components rather than extending base classes
2. **Stateless Rendering**: Container components are pure rendering functions with no state
3. **Progressive Enhancement**: Shared navigation features with model-specific interactions
4. **Consistent UX**: Users see familiar patterns regardless of list type

### Component Hierarchy

```
Screen
├── Header Component
│   └── Title (e.g., "💥 Bursts", "📋 Facts", "Career Events")
├── List Content
│   ├── Optional: Custom Header (column names, filters, etc.)
│   ├── ListContainer
│   │   ├── Items (rendered via model-specific logic)
│   │   └── Pagination Info
│   ├── Empty State (if no items)
│   └── Status Messages
├── ScreenContainer (padding/layout wrapper)
└── Footer Component
    └── Status bar/help text

```

## Rendering Pipeline

### Standard View() Method Structure

All list models should follow this rendering pattern:

```go
func (m *ListModel) View() string {
    // 1. Get displayed items (with pagination calculation)
    displayedItems := m.getDisplayedItems()

    // 2. Build item strings
    var items []string
    for _, item := range displayedItems {
        items = append(items, m.renderItem(item))
    }

    // 3. Create pagination info
    paginationInfo := m.getPaginationInfo()

    // 4. Setup ListContainer with items
    listContainer := components.NewListContainer().
        SetItems(items).
        SetEmptyStateMessage(m.getEmptyStateMessage()).
        SetPaginationInfo(paginationInfo)

    // 5. Create screen content
    screenContent := lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderOptionalHeader(),  // If needed
        listContainer.Render(),
    )

    // 6. Wrap with ScreenContainer
    screenContainer := components.NewScreenContainer(screenContent).
        WithPaddingMode(components.PaddingNormal)

    // 7. Add header and footer
    headerView := components.NewHeader(m.getTitle(), m.width).View()
    footerView := components.NewFooter(m.width).View()

    // 8. Join all sections
    fullContent := lipgloss.JoinVertical(
        lipgloss.Left,
        headerView,
        "",
        screenContainer.Render(),
        "",
        footerView,
    )

    return fullContent
}
```

### Key Rendering Considerations

1. **Empty Lines for Spacing**: Use `""` between major sections for visual breathing room
2. **Consistent Ordering**: Header → Content → Footer (always in this order)
3. **Pagination Visibility**: Always include pagination info via `SetPaginationInfo()`
4. **Error Messages**: Handle gracefully with error styling if needed

## Component Structure

### Required Components

Every list model must include:

1. **Header Component** (via `components.NewHeader()`)
   - Shows list title/name
   - Displays at top of view
   - Title should include emoji for visual identity

2. **ListContainer Component** (via `components.NewListContainer()`)
   - Renders individual items
   - Handles empty state display
   - Displays pagination information

3. **ScreenContainer Component** (via `components.NewScreenContainer()`)
   - Wraps content with padding
   - Provides consistent margins
   - Uses `PaddingNormal` by default

4. **Footer Component** (via `components.NewFooter()`)
   - Shows status information
   - Displays at bottom of view
   - Optional: help text or action hints

### Optional Components

Models may optionally include:

1. **Custom Header Row** (for column labels or filters)
   - Example: Burst list shows column headers (Name, Events, Competency, Created)
   - Must use `styles.HeaderSection` for styling

2. **Help Footer** (additional key binding display)
   - Example: List model includes `helpFooter` for showing available shortcuts

## Visual Specification

### Title Display

**Pattern**: `emoji Title[Optional Filters]`

**Examples**:
- `Career Events` - Career events list
- `💥 Bursts` - Bursts list with emoji
- `📋 Facts (Competency: Technical)` - Facts list with filter indicator

**Implementation**:
```go
func (m *ListModel) getTitle() string {
    title := "📋 Facts"
    if m.competencyFilter != "" {
        title += fmt.Sprintf(" (Competency: %s)", m.competencyFilter)
    }
    return title
}
```

### Item Display Format

**General Pattern**: `[Marker] content [metadata]`

**Examples by Type**:

| Type | Format | Marker |
|------|--------|--------|
| Event | `▶ Event text [date] [company]` | `▶` selected, `  ` normal |
| Burst | `Name \| Events \| Competency \| Date` | Highlight (style) |
| Fact | `[☐/☑] Fact text [metadata]` | Checkbox for selection |

### Pagination Display

**Pattern**: `Showing X-Y of Z items` (STANDARDIZED as of 2026-01-01)

**Examples**:
- `Showing 1-10 of 50 events` - Events pagination
- `Showing 1-3 of 5 bursts` - Bursts pagination
- `Showing 1-8 of 12 facts` - Facts pagination

**Implementation**:
```go
func (m *ListModel) getPaginationInfo() string {
    maxIdx := m.height - 5
    endIdx := m.scrollOffset + maxIdx
    if endIdx > len(m.filtered) {
        endIdx = len(m.filtered)
    }
    startIdx := m.scrollOffset + 1

    info := fmt.Sprintf("Showing %d-%d of %d items", startIdx, endIdx, len(m.items))
    return info
}
```

**Note**: This format is MANDATORY for all list models as of 2026-01-01. The previous format `X/Y items` is DEPRECATED.

### Empty State Display

**Pattern**: `No <items> found` (STANDARDIZED as of 2026-01-01)

**Examples**:
- `"No events found"` - Career events list
- `"No bursts found"` - Bursts list
- `"No facts found"` - Facts list

**Guidelines**:
- Use ONLY the standardized format `"No <items> found"`
- Do NOT use conditional messages based on filter state
- Do NOT use alternative wording like "No matching <items>"
- Use proper capitalization and punctuation
- Be consistent across all list types

**Prohibited Patterns** (as of 2026-01-01):
- ❌ `"No matching bursts"` - Conditional based on filter
- ❌ `"No facts match the current filters"` - Filter-specific messaging
- ❌ `"No events found. Start capturing your career journey!"` - Additional guidance text

### Color Usage

**Exported Constants** (all colors must use these):

| Constant | Usage |
|----------|-------|
| `ColorTextPrimary` | Primary text content |
| `ColorTextSecondary` | Secondary information (dates, companies) |
| `ColorTextMuted` | Less important metadata |
| `ColorAccentTeal` | Focus indicator, primary action |
| `ColorError` | Error messages, destructive actions |
| `ColorSuccess` | Success states, positive feedback |
| `ColorInfo` | Informational hints and help text |
| `ColorBorder` | Normal borders, dividers |
| `ColorBorderActive` | Active element borders |

**Implementation Pattern**:
```go
styles.ListItem.
    Foreground(styles.ColorTextSecondary).
    Render(dateStr)
```

### Style Objects

**Exported Style Objects** (use these for consistent appearance):

| Object | Usage |
|--------|-------|
| `ListItem` | Normal list item rendering |
| `ListItemSelected` | Highlighting for selected item |
| `HeaderSection` | Column headers, section titles |
| `CardBase` | Container card backgrounds |
| `ErrorText` | Error message styling |
| `InfoHint` | Hint text and secondary info |
| `InfoText` | Informational text display |

## Navigation Specification

### Required Navigation Keys

All list models MUST support these navigation keys:

| Key | Alt | Action | Behavior |
|-----|-----|--------|----------|
| up | k | Previous item | Move cursor up, wrap around |
| down | j | Next item | Move cursor down, wrap around |
| Page Up | ctrl+b | Previous page | Scroll up by page height |
| Page Down | ctrl+f | Next page | Scroll down by page height |
| Home | g | First item | Jump to first item in list |
| End | G | Last item | Jump to last item in list |
| Escape | - | Go back | Return to parent screen |
| q | ctrl+c | Quit | Exit application |

### Recommended Navigation Keys

Models MAY implement these for specific functionality:

| Key | Action | Requirement |
|-----|--------|-------------|
| Enter | Select/View | Type-specific (view event, select fact, etc.) |
| Space | Toggle/Expand | Type-specific (expand burst, select fact, etc.) |
| ? | Show help | Optional, recommended |

### Key Handling Implementation Pattern

```go
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            m.prevItem()
        case "down", "j":
            m.nextItem()
        case "pgup", "ctrl+b":
            m.prevPage()
        case "pgdn", "ctrl+f":
            m.nextPage()
        case "home", "g":
            m.goToFirstItem()
        case "end", "G":
            m.goToLastItem()
        case "esc":
            return m, func() tea.Msg { return BackMsg{} }
        case "q", "ctrl+c":
            return m, func() tea.Msg { return QuitMsg{} }
        // Type-specific keys
        case "enter":
            return m, m.handleSelect()
        case " ", "space":
            return m, m.handleToggle()
        }
    }
    return m, nil
}
```

### Navigation Helper Methods

All models MUST implement these methods:

```go
func (m *ListModel) nextItem()          // Move to next item with wrap
func (m *ListModel) prevItem()          // Move to previous item with wrap
func (m *ListModel) nextPage()          // Move forward by page size
func (m *ListModel) prevPage()          // Move backward by page size
func (m *ListModel) goToFirstItem()     // Jump to first item
func (m *ListModel) goToLastItem()      // Jump to last item
```

### Navigation Response Time

- All navigation should be instant (< 1ms)
- No loading spinners for navigation
- Immediate visual feedback on key press
- No pagination loading (calculate in advance)

## Styling Specification

### Color Palette

**Must Use Exported Constants from `styles` package**:

All color assignments must use:
- `styles.Color*` constants for colors
- `styles.ListItem`, `styles.CardBase`, etc. for styles

**Prohibited**:
- Direct hex color values (e.g., `#FF00FF`)
- `lipgloss.Color()` calls outside styles package
- Magic number color codes

### Style Composition

**Correct Pattern**:
```go
styles.ListItem.
    Foreground(styles.ColorTextSecondary).
    Render(dateStr)
```

**Incorrect Pattern**:
```go
lipgloss.NewStyle().
    Foreground("#FF00FF").
    Render(text)  // WRONG: inline hex color
```

### Spacing Standards

**Padding**:
- Models should rely on ScreenContainer for external padding
- Use `PaddingNormal` mode for standard padding
- No additional custom padding on model output

**Margins**:
- Empty lines between sections: Use `""` in JoinVertical
- Item spacing: Handled by ListContainer
- No hardcoded margin values

**Line Heights**:
- Default line height from lipgloss
- No custom line spacing needed
- Ensure readability with proper margins

### Container Usage

**Standard Pattern**:
```go
screenContainer := components.NewScreenContainer(screenContent).
    WithPaddingMode(components.PaddingNormal)
```

**Alternatives**:
- Use `CardBase` for wrapping if additional background/border needed
- Use `components.CardContainer` for complex layouts
- Never hardcode padding values

## Implementation Checklist

When implementing a new list model, ensure:

### Rendering (Required)
- [ ] View() method follows standard structure
- [ ] Uses Header component with appropriate title
- [ ] Uses ListContainer for item rendering
- [ ] Includes pagination info via SetPaginationInfo()
- [ ] Handles empty state with clear message
- [ ] Uses ScreenContainer for wrapping
- [ ] Includes Footer component
- [ ] All sections joined in correct order

### Navigation (Required)
- [ ] Implements all 8 required navigation keys
- [ ] All navigation methods implemented
- [ ] Selection wraps at boundaries
- [ ] Keyboard response is instant
- [ ] Visual feedback shows selection clearly

### Styling (Required)
- [ ] All colors use styles.Color* constants
- [ ] All styles use exported style objects
- [ ] No inline hex colors
- [ ] No hardcoded magic numbers
- [ ] Uses appropriate colors for content type
- [ ] Proper text styling for emphasis/secondary info

### Data Handling (Required)
- [ ] Handles empty list gracefully
- [ ] Filters applied consistently
- [ ] Sorting works correctly
- [ ] Pagination info accurate
- [ ] Selected items tracked if applicable

### Testing (Required)
- [ ] Unit tests for rendering output
- [ ] Navigation tests for key handling
- [ ] Edge case tests (empty, single item, etc.)
- [ ] Style consistency verified
- [ ] No regressions in existing tests

### Documentation (Recommended)
- [ ] Code comments explaining complex logic
- [ ] Method documentation for public APIs
- [ ] Navigation shortcuts documented
- [ ] Any special features noted

## Code Examples

### Example 1: Simple List Model

```go
type SimpleListModel struct {
    items     []Item
    selectedIdx int
    width     int
    height    int
}

func (m *SimpleListModel) View() string {
    var items []string
    for i, item := range m.items {
        marker := "  "
        if i == m.selectedIdx {
            marker = "▶ "
        }
        items = append(items, marker + item.Text)
    }

    paginationInfo := fmt.Sprintf("%d/%d items",
        len(m.items), len(m.items))

    listContainer := components.NewListContainer().
        SetItems(items).
        SetEmptyStateMessage("No items found").
        SetPaginationInfo(paginationInfo)

    screenContent := lipgloss.JoinVertical(
        lipgloss.Left,
        listContainer.Render(),
    )

    screenContainer := components.NewScreenContainer(screenContent).
        WithPaddingMode(components.PaddingNormal)

    headerView := components.NewHeader("Items", m.width).View()
    footerView := components.NewFooter(m.width).View()

    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerView,
        "",
        screenContainer.Render(),
        "",
        footerView,
    )
}

func (m *SimpleListModel) nextItem() {
    if m.selectedIdx < len(m.items)-1 {
        m.selectedIdx++
    }
}

func (m *SimpleListModel) prevItem() {
    if m.selectedIdx > 0 {
        m.selectedIdx--
    }
}

// ... other navigation methods
```

### Example 2: List Model with Filtering

```go
func (m *FilteredListModel) getDisplayedItems() []Item {
    var displayed []Item
    for _, item := range m.items {
        if m.matchesFilter(item) {
            displayed = append(displayed, item)
        }
    }
    return displayed
}

func (m *FilteredListModel) getPaginationInfo() string {
    displayed := len(m.getDisplayedItems())
    total := len(m.items)
    info := fmt.Sprintf("%d/%d items", displayed, total)

    if displayed < total {
        info += fmt.Sprintf(" (filtered from %d)", total)
    }
    return info
}
```

## Common Patterns

### Pattern 1: Multi-Select List

Use checkboxes for selection indicators:
```go
checkbox := "☐"
if isSelected {
    checkbox = "☑"
}
item := fmt.Sprintf("%s %s", checkbox, text)
```

### Pattern 2: Expandable Items

Track expanded state and render nested items:
```go
if m.expandedIndices[i] {
    items = append(items, m.renderExpandedContent(item))
}
```

### Pattern 3: Filtered List

Always show filter status in pagination:
```go
info := fmt.Sprintf("%d/%d items", displayed, total)
if m.filterActive {
    info += " (filtered)"
}
```

### Pattern 4: Sortable List

Maintain sort order and indicate in header:
```go
title := "Items"
if m.sortBy != "default" {
    title += fmt.Sprintf(" (sorted by %s)", m.sortBy)
}
```

## Troubleshooting

### Issue: Rendering Inconsistency Between Models

**Symptom**: One list model looks different from others
**Solution**:
1. Check View() method structure matches specification
2. Verify ListContainer usage
3. Check all colors use exported constants
4. Ensure same Header/Footer components used

### Issue: Navigation Not Working

**Symptom**: Keyboard navigation doesn't respond
**Solution**:
1. Verify Update() method handles tea.KeyMsg
2. Check switch statement covers all required keys
3. Verify navigation methods (nextItem, prevItem, etc.) implemented
4. Test with automated tests

### Issue: Pagination Info Not Showing

**Symptom**: No pagination displayed in list
**Solution**:
1. Verify SetPaginationInfo() called on ListContainer
2. Check pagination string is not empty
3. Verify correct data passed to SetPaginationInfo()
4. Check ListContainer component renders pagination

### Issue: Styling Inconsistency

**Symptom**: Colors or styles don't match other lists
**Solution**:
1. Check all colors use styles.Color* constants
2. Verify styles use exported style objects
3. Search for inline color definitions
4. Run style audit: `grep -r "lipgloss.Color"`

### Issue: Empty State Not Showing

**Symptom**: List shows nothing when empty
**Solution**:
1. Verify SetEmptyStateMessage() called
2. Check message string is not empty
3. Verify getDisplayedItems() returns empty array
4. Check ListContainer handles empty state (should)

## References

- Component Documentation: `/docs/guides/COMPONENT_USAGE_GUIDE.md`
- Style Guide: `/docs/guides/STYLE_USAGE_GUIDE.md`
- Implementation Assessment: `/docs/features/06-model-consistency-implementation-assessment.md`
- Related Tests: `/internal/cli/models/list_rendering_consistency_test.go`

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-01-01 | Initial specification, covers core requirements |

## Contact & Support

For questions about list model implementation:
1. Review this specification and code examples
2. Check existing implementations (list.go, burst_list.go, fact_list.go)
3. Run regression tests to verify consistency
4. Refer to related documentation and guides

---

**Specification Status**: ACTIVE AND ENFORCED
**Last Review**: 2026-01-01
**Next Review**: Upon new list model implementation or major redesign

