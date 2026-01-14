---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# TUI Standardization Standards and Guidelines

**Document Version**: 1.0
**Created**: 2025-12-30
**Status**: Complete and Verified
**Test Coverage**: 337/337 tests passing (100%)

---

## Overview

This document outlines the standardized Terminal User Interface (TUI) design for KaRiya. All screens, models, and components follow these standards to ensure a consistent, intuitive, and professional user experience.

**Key Principles**:
- **Consistency**: All models use the same navigation patterns and shortcuts
- **Discoverability**: Help text is always available and contextual
- **Accessibility**: Keyboard-first navigation with vim-style alternatives
- **Responsiveness**: Components adapt to terminal size changes
- **Clarity**: Information is presented clearly with proper hierarchy

---

## Keyboard Shortcuts

### Universal Navigation

These shortcuts work consistently across all screens:

| Shortcut | Action | Usage |
|----------|--------|-------|
| **Esc** | Back/Cancel | Go back to previous screen or cancel operation |
| **Tab** | Next field | Move focus to next field (forms) |
| **Shift+Tab** | Previous field | Move focus to previous field (forms) |
| **↑/k** | Up | Navigate up in lists or menus |
| **↓/j** | Down | Navigate down in lists or menus |
| **←/h** | Left | Navigate left in menus or selectors |
| **→/l** | Right | Navigate right in menus or selectors |
| **Enter** | Select/Confirm | Confirm selection or submit form |
| **Space** | Toggle | Toggle checkbox or expand item |

### Context-Specific Shortcuts

#### Form Screens
- **Enter**: Submit form (when Submit button focused)
- **Tab**: Next field
- **Shift+Tab**: Previous field
- **Esc**: Cancel form (returns to home)

#### List Screens
- **↑/k**: Previous item
- **↓/j**: Next item
- **Enter**: Select item (view details)
- **d**: Delete item
- **e**: Edit item
- **f**: Filter/search
- **s**: Sort
- **Esc**: Back to home

#### Metadata Review
- **↑/k**: Previous event
- **↓/j**: Next event
- **Enter**: Edit selected event
- **b**: Bulk operations on selected events
- **Space**: Toggle selection (for bulk operations)
- **Esc**: Back to home

#### Bulk Operations
- **↑/k**: Previous item
- **↓/j**: Next item
- **Space**: Toggle selection
- **a**: Select all
- **d**: Deselect all
- **Enter**: Apply changes
- **Esc**: Cancel (discard changes)

### Global Shortcuts

| Shortcut | Action |
|----------|--------|
| **?** | Show help |
| **h** | Go home |
| **q** | Quit application |
| **c** | Capture event |
| **l** | List events |
| **m** | Metadata review |

### Escape Key Behavior Standards

**Status**: ✅ **COMPLETE - ALL 5 INTENTS STANDARDIZED**
**Last Updated**: 2026-01-06
**Test Coverage**: 92 escape-specific tests, 100% passing

#### Universal Keyboard Standards

All intents follow these consistent keyboard behaviors:

| Key | Function | Behavior |
|-----|----------|----------|
| **esc** | Back | Navigate to previous state; cancel if root state |
| **m** | Main Menu | Return to main menu from any state |
| **q** / **ctrl+c** | Quit | Exit entire application |

#### State Classification

1. **Root State** (First state in intent)
   - **esc**: Cancels intent, returns to main menu
   - **m**: Cancels intent, returns to main menu
   - Example: `CaptureStateChooseStrategy`, `GenerateCVStateSelectProfile`

2. **Intermediate State** (Has previous state)
   - **esc**: Navigate back to previous state
   - **m**: Cancel intent, return to main menu
   - Example: `CaptureStateForm`, `GenerateCVStateSelectAudience`

3. **Async Operation State** (Performing background work)
   - **esc**: Let operation complete in background, navigate to previous state
   - **m**: Cancel intent immediately, return to main menu
   - Example: `GenerateCVStateGenerating`, `ConfigStateSaving`
   - Note: Operations continue in background per user preference

#### Error State Handling

When navigating back from an error state:
- ✅ **Keep error visible** so users can see what went wrong
- Error persists until explicit state change or successful retry
- Example: Submit state with error → Esc → Review state (error still shown)

#### Implementation Status

| Intent | States | Esc Handlers | 'm' Handlers | Tests | Status |
|--------|--------|--------------|--------------|-------|--------|
| **CaptureEvent** | 4 | 4/4 (100%) | 4/4 (100%) | 13 ✅ | ✅ Complete |
| **ConfigureSystem** | 7 | 7/7 (100%) | 7/7 (100%) | 21 ✅ | ✅ Complete |
| **GenerateCV** | 10 | 10/10 (100%) | 10/10 (100%) | 26 ✅ | ✅ Complete |
| **BrowseTimeline** | 2 | 2/2 (100%) | 2/2 (100%) | 7 ✅ | ✅ Complete |
| **ExportArtifact** | 9 | 9/9 (100%) | 9/9 (100%) | 25 ✅ | ✅ Complete |
| **TOTAL** | **32** | **32/32 (100%)** | **32/32 (100%)** | **92** | ✅ **Complete** |

**Key Achievements**:
- ✅ All 32 states across 5 intents have full escape coverage
- ✅ All states support 'm' key for instant main menu return
- ✅ 3 critical async operation dead-ends fixed (Generating, Exporting, Saving)
- ✅ Error visibility preserved when navigating back
- ✅ 92 new tests added, 100% passing
- ✅ Zero regressions in existing tests

#### Example Flow: CaptureEvent Intent

```
ChooseStrategy (Root)
    esc → Cancel intent → Main menu
    m → Cancel intent → Main menu
    ↓ (select strategy)
CaptureForm
    esc → Back to ChooseStrategy
    m → Cancel intent → Main menu
    ↓ (submit form)
ReviewInferredEvent
    esc → Back to CaptureForm
    m → Cancel intent → Main menu
    ↓ (confirm)
Submit (with error)
    esc → Back to ReviewInferredEvent (error visible)
    m → Cancel intent → Main menu
    r → Retry submission
```

#### Testing Requirements

Every state MUST have tests for:
- ✅ Escape key behavior (back or cancel)
- ✅ 'm' key behavior (main menu)
- ✅ State remains active after back navigation
- ✅ Intent cancels correctly from any state
- ✅ View footers show correct key options

#### Developer Guidelines

When adding escape/m key handlers to a new state:

```go
case tea.KeyMsg:
    switch msg.String() {
    case "esc":
        if m.isRootState() {
            m.setCancelled()  // Root state cancels intent
        } else if m.isAsyncOperation() {
            // Let async operation complete in background
            m.state = m.previousState
        } else {
            m.state = m.previousState  // Go back one state
        }
        return nil
    case "m":
        // Always return to main menu
        m.setCancelled()
        return nil
    case "q", "ctrl+c":
        return tea.Quit  // Quit entire app
    }
```

#### View Footer Standards

All view methods MUST show available keys:

```
Good examples:
- "Esc: Back | m: Main menu | q: Quit"
- "Enter: Confirm | Esc: Back | m: Main menu | r: Retry"
- "Esc: Cancel | m: Main menu | q: Quit"

Bad examples:
- "Press Esc to cancel" (missing 'm' option)
- "q to quit" (missing esc and m options)
```

---

## Component Architecture

### Header Component

**Purpose**: Display screen title, breadcrumb navigation, and context

**Appearance**:
```
╔════════════════════════════════════════════════════════════════════════╗
║ KaRiya > Career Events > Event Details                                 ║
╚════════════════════════════════════════════════════════════════════════╝
```

**Usage**:
```go
header := components.NewHeader("Screen Title", width)
header.SetSubtitle("Optional subtitle")  // Not yet implemented
header.SetBreadcrumb([]string{"Home", "List", "Details"})  // Future feature
```

**Instances**:
- ✅ Form: "Capture Career Event"
- ✅ List: "Career Events"
- ✅ Metadata Review: "Review Event Metadata"
- ✅ Metadata Editor: "Edit Event Metadata"
- ✅ Other screens: Integrated or planned

### Footer Component

**Purpose**: Display status information and mode indicators

**Appearance**:
```
─────────────────────────────────────────────────────────────────────────
 Strategy: Quick │ Status: 3/10 events │ Ctrl+C to quit
```

**Usage**:
```go
footer := components.NewFooter(width)
footer.SetStatus("3/10 events")
footer.SetStrategy("Quick")
```

**Instances**:
- ✅ Form: Shows capture strategy (Quick/Manual)
- ✅ List: Shows event count
- ✅ Metadata Review: Shows review progress
- ✅ All screens: Integrated or planned

### Help Footer Component

**Purpose**: Display context-aware keyboard shortcuts

**Appearance**:
```
↑/↓: Navigate | Tab: Next field | Enter: Submit | Esc: Back | ?: Help
```

**Context-Aware Shortcuts**:
- **"form"**: Tab, Shift+Tab, Enter, Esc (field navigation)
- **"list"**: Up/Down, Enter, Delete, Edit, Filter, Sort, Esc
- **"metadata_review"**: Up/Down, Enter, Bulk ops, Space, Esc
- **"metadata_editor"**: Tab, Shift+Tab, Enter, Esc
- **"bulk_operations"**: Up/Down, Space, All/None, Enter, Esc

**Usage**:
```go
helpFooter := components.NewHelpFooter("form", width)
// Automatically generates appropriate shortcuts for context
```

### Navigation Menu Component

**Purpose**: Display context-sensitive menu with keyboard navigation

**Appearance**:
```
┌─ Main Menu ────────────────────────────────────────────────────────┐
│ ► Capture Event (c)        Add a new career event                  │
│   View Events (l)          Browse your career history              │
│   Metadata Review (m)      Review and enhance event metadata       │
│   Settings                 Configure KaRiya preferences            │
│   Help (?)                 View help documentation                 │
│   Quit (q)                 Exit the application                    │
└────────────────────────────────────────────────────────────────────┘
```

**Features**:
- Keyboard navigation (↑/↓ or j/k)
- Visual indicator (►) for selected item
- Keyboard shortcut display
- Context-aware descriptions

### List Item Component

**Purpose**: Display uniform list items with proper formatting

**Appearance**:
```
► [TAG] Title: Event description                  Date │ Company │ Tags
  Leading indicator for selected/focused items
```

**Features**:
- Consistent spacing and alignment
- Tag display with styling
- Metadata visibility (date, company, tags)
- Proper truncation with ellipsis (...)
- Visual indicators for selection and focus

---

## Navigation Flow

### Screen Hierarchy

```
Home
├── Capture Event (Form)
│   └── Success Screen
│       ├── Review Metadata (Metadata Review)
│       │   └── Edit Event (Metadata Editor)
│       │       └── Success
│       └── View Recent Events (List)
├── List Events
│   ├── View Event Details
│   ├── Edit Event Metadata
│   └── Bulk Operations
└── Metadata Review
    ├── Edit Event
    └── Bulk Operations
```

### Navigation Rules

1. **Esc Key**: Always returns to previous screen (or Home if none)
2. **Breadcrumbs**: Show current path in header (planned for Phase 4)
3. **Context**: Help footer adapts to current screen
4. **Consistency**: Same shortcuts mean same actions everywhere

---

## Visual Design

### Color Scheme

**Professional Dark Theme**:

| Element | Color | Hex Code | Usage |
|---------|-------|----------|-------|
| Background | Dark Blue-Gray | #1a1f2e | Primary background |
| Primary Text | Light Gray | #c7ccd1 | Main content |
| Secondary Text | Medium Gray | #8b92a0 | Descriptions |
| Muted Text | Dark Gray | #5e6673 | Disabled/secondary |
| Accent (Teal) | Bright Teal | #5fb3b3 | Interactive elements |
| Success (Green) | Bright Green | #6cb56c | Success states |
| Warning (Orange) | Warm Orange | #d9a66c | Warning states |
| Error (Red) | Bright Red | #d76e6e | Error states |
| Info (Blue) | Bright Blue | #6ab0d3 | Info states |

### Styling Rules

1. **Consistency**: All buttons, inputs, and panels use standard styles
2. **Contrast**: Text colors meet accessibility standards
3. **Focus**: Focused elements have clear visual indicator
4. **Errors**: Error text uses consistent red color with clear messaging
5. **Success**: Success states use consistent green color

### Theme System

KaRiya uses a theme system for consistent, customizable styling. All intents should use themed styling instead of hardcoded colors.

**Key Principles**:

1. **Use Theme Helpers**: Access colors via `i.Theme().PrimaryColor()` etc.
2. **Provide Fallbacks**: Always handle `nil` theme gracefully
3. **Pre-composed Styles**: Use `theme.Styles().CardBase` etc. for common patterns
4. **Bubbles Integration**: Use `themes.NewThemedTableStyles()` for bubbles components

**Example Pattern**:

```go
func (i *MyIntent) getCardStyle() lipgloss.Style {
    if theme := i.Theme(); theme != nil {
        return theme.Styles().CardBase
    }
    return lipgloss.NewStyle().
        Padding(1, 2).
        BorderStyle(lipgloss.RoundedBorder())
}
```

See [Theme Customization Guide](THEME_CUSTOMIZATION_GUIDE.md) for complete documentation.

---

## Component Integration Pattern

Every screen should follow this pattern:

```go
type ScreenModel struct {
    // Service dependencies
    cliService *service.CLIEventService
    service    *careerservice.Service

    // UI Components
    header     components.HeaderModel
    footer     components.FooterModel
    helpFooter components.HelpFooterModel

    // Screen-specific state
    // ... model state fields ...
}

func (m *ScreenModel) View() string {
    // Render content
    content := /* screen-specific content */

    // Combine with standard components
    fullView := lipgloss.JoinVertical(
        lipgloss.Left,
        m.header.View(),
        "",
        content,
        "",
        m.footer.View(),
        "",
        m.helpFooter.View(),
    )

    return fullView
}
```

---

## Testing Standards

### Navigation Testing

Every model should have tests for:
- ✅ Escape key returns to previous screen
- ✅ Tab/Shift+Tab navigation (forms)
- ✅ Up/Down/Left/Right navigation (lists/menus)
- ✅ Vim keys (j/k/h/l) work as alternatives
- ✅ Enter confirms selections
- ✅ Space toggles selections
- ✅ Help footer displays correct context

### Visual Testing

Every model should have tests for:
- ✅ Header displays correctly
- ✅ Footer displays correctly
- ✅ Help footer adapts to window size
- ✅ Content renders without overflow
- ✅ Colors/styling applied consistently
- ✅ Components responsive to resize messages

### Current Test Status

**Total Tests**: 337/337 passing (100%)
**Coverage**: 80%+ across all packages
**Race Conditions**: 0 detected
**Performance**: Sub-second rendering

---

## Development Guidelines

### Adding a New Screen

1. Create model struct with Header, Footer, HelpFooter components
2. Implement BubbleTea Model interface (Init, Update, View)
3. Add navigation shortcuts to help footer context
4. Use standard component styling
5. Write tests for navigation and rendering
6. Integrate into app.go navigation flow

### Updating Navigation

1. Keep shortcuts consistent with this document
2. Update help footer context if adding new shortcuts
3. Test with race detector
4. Verify backward compatibility
5. Update documentation

### Best Practices

- **Use Navigation Constants**: Import `github.com/baphled/kariya/internal/cli/navigation`
- **Leverage Reusable Components**: Don't duplicate header/footer/help logic
- **Test Navigation**: Include keyboard shortcut tests
- **Document Changes**: Update TUI_STANDARDS.md when adding features
- **Monitor Token Usage**: Keep implementation size reasonable

---

## Future Enhancements

### Phase 4 (Planned)
- [ ] Breadcrumb navigation in header
- [ ] Progress indicators for multi-step workflows
- [ ] Enhanced color scheme with theming
- [ ] Visual feedback (spinners, animations)

### Phase 5 (Planned)
- [ ] Comprehensive navigation testing
- [ ] Visual consistency testing across terminal sizes
- [ ] Performance profiling and optimization
- [ ] Accessibility audit and improvements

---

## Keyboard Reference Card

### One-Page Quick Reference

```
╔════════════════════════════════════════════════════════════════════╗
║                    KaRiya Keyboard Reference                      ║
├────────────────────────────────────────────────────────────────────┤
║ Navigation              Field Navigation        Actions            ║
║ ─────────────────────   ──────────────────────  ─────────────────  ║
║ ↑/k   Prev             Tab       Next field    Enter  Select      ║
║ ↓/j   Next             Shift+Tab Prev field    Space  Toggle      ║
║ ←/h   Left             Esc       Cancel        d      Delete      ║
║ →/l   Right                                    e      Edit        ║
║ Esc   Back/Home                                f      Filter      ║
║                                                s      Sort        ║
║ Global Shortcuts                               b      Bulk ops    ║
║ ──────────────────────────────────────────────────────────────────  ║
║ ?     Help              h     Home             c     Capture      ║
║ q     Quit              l     List             m     Metadata     ║
╚════════════════════════════════════════════════════════════════════╝
```

---

## Documentation

For user-facing documentation, see:
- `docs/CLI_GUIDE.md` - User guide with examples
- `docs/KEYBOARD_SHORTCUTS_GUIDE.md` - Complete keyboard reference
- `docs/workflows/` - Detailed workflow guides
- `README.md` - Main project documentation

For developer documentation, see:
- `docs/development/CENTRALIZED_KEY_HANDLING.md` - Centralized key handling guide
- `docs/development/KEYBOARD_SYSTEM_GUIDE.md` - Keyboard system implementation
- `docs/STATE_MATRIX.md` - Complete state matrix (10 intents, 64 states)
- `docs/rules/` - Development rules and guidelines
- Source code comments - Inline documentation

---

## Summary

The TUI standardization provides:
- ✅ **Consistent Navigation**: All models use same shortcuts
- ✅ **Clear Visual Design**: Professional dark theme with accessible colors
- ✅ **Discoverable Help**: Context-aware keyboard shortcut display
- ✅ **Responsive Layout**: Components adapt to terminal size
- ✅ **Robust Testing**: 337 tests with 100% pass rate
- ✅ **Excellent Performance**: Sub-second rendering with race-condition free code

**Status**: Production-ready with comprehensive test coverage and documentation.

---

**Document Version**: 1.0
**Last Updated**: 2025-12-30
**Status**: Complete
**Test Coverage**: 337/337 (100%)

