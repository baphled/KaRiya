# UIKit Component Guide

**Version**: 1.1  
**Created**: 2026-01-19  
**Updated**: 2026-01-20  
**Status**: Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Quick Start](#quick-start)
4. [Primitives](#primitives)
5. [Containers](#containers)
6. [Feedback Components](#feedback-components) (Modal system)
7. [Layout Components](#layout-components) (ScreenLayout)
8. [Navigation Components](#navigation-components) (BreadcrumbBar)
9. [Display Components](#display-components) (Logo)
10. [Theme Integration](#theme-integration)
11. [Migration Guide](#migration-guide)
12. [Mandatory Patterns](#mandatory-patterns)
13. [Component Checklist](#component-checklist)
14. [API Reference](#api-reference)

---

## Overview

The UIKit is KaRiya's **standardized component library** for building consistent, theme-aware TUI components. It provides:

- **Primitives**: Foundational components (Text, Button, Badge, Input)
- **Containers**: Layout components (Box, Overlay)
- **Theme System**: Consistent styling via `theme.Aware` embedding

### Why UIKit?

| Before (Direct Lipgloss) | After (UIKit) |
|--------------------------|---------------|
| Scattered styling code | Centralized, reusable components |
| Manual theme handling | Automatic theme integration |
| Inconsistent colors | Semantic color system |
| No nil safety | Built-in nil guards |
| Verbose code | Fluent builder API |

### Package Structure

```
internal/cli/uikit/
├── primitives/           # Text, Button, ButtonGroup, Badge, Input
│   ├── text.go          # Semantic text styles
│   ├── button.go        # Interactive buttons
│   ├── button_group.go  # Keyboard-navigable button groups
│   ├── badge.go         # Status badges + help footer helpers
│   └── input.go         # Text input wrapper
├── containers/          # Box, Overlay
│   ├── box.go           # Bordered containers for modals/cards
│   └── overlay.go       # Centered modal overlays
├── feedback/            # Modal system (NEW)
│   ├── modal.go         # Modal types (Error, Loading, Progress, Success, Warning)
│   ├── modal_container.go # Custom modal dialog builder
│   └── help_modal.go    # Context-sensitive keyboard shortcuts
├── layout/              # Screen layouts (NEW)
│   ├── screen_layout.go # StandardView replacement - main layout component
│   ├── header.go        # Screen headers
│   └── footer.go        # Screen footers with status/mode
├── navigation/          # Navigation components (NEW)
│   └── breadcrumb.go    # BreadcrumbBar with icons
├── display/             # Display components (NEW)
│   └── logo.go          # Animated ASCII logo
└── theme/               # Theme infrastructure
    ├── theme.go         # Theme interface (re-exports from themes package)
    └── aware.go         # Embeddable theme awareness
```

---

## Architecture

### Component Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                         Theme System                             │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    theme.Aware (embedded)                   ││
│  │  • SetTheme(theme)  • Theme()  • PrimaryColor()  etc.      ││
│  └─────────────────────────────────────────────────────────────┘│
│                              │                                   │
│         ┌────────────────────┼────────────────────┐              │
│         ▼                    ▼                    ▼              │
│   ┌───────────┐        ┌───────────┐        ┌───────────┐       │
│   │ Primitives│        │ Containers│        │  Screens  │       │
│   │ Text      │        │ Box       │        │ (Intents) │       │
│   │ Button    │        │ Overlay   │        │           │       │
│   │ Badge     │        └───────────┘        └───────────┘       │
│   │ Input     │                                                  │
│   └───────────┘                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Fluent Builder Pattern

All UIKit components use a fluent builder pattern:

```go
// Create → Configure → Render
result := primitives.NewText("Hello", theme).
    Bold().
    MarginBottom(1).
    Render()
```

---

## Quick Start

### Using Text Primitives

```go
import "github.com/baphled/kariya/internal/cli/uikit/primitives"

// Semantic text styles
title := primitives.Title("Welcome", theme).Render()
subtitle := primitives.Subtitle("Settings", theme).Render()
body := primitives.Body("Content here", theme).Render()
muted := primitives.Muted("Hint text", theme).Italic().Render()
errorMsg := primitives.ErrorText("Failed!", theme).Bold().Render()
success := primitives.SuccessText("Saved!", theme).Render()
info := primitives.InfoText("Note: ...", theme).Render()

// Custom styling
custom := primitives.NewText("Custom", theme).
    Bold().
    Italic().
    Foreground(lipgloss.Color("#FF0000")).
    MarginTop(1).
    MarginBottom(2).
    Width(40).
    Render()
```

### Using Box Containers

```go
import "github.com/baphled/kariya/internal/cli/uikit/containers"

// Basic box
box := containers.NewBox(theme).
    Content("Hello World").
    Render()

// Modal box with solid background (REQUIRED for overlays)
modalBox := containers.NewBox(theme).
    Title("Confirm Delete").
    Content("Are you sure?").
    Variant(containers.BoxDestructive).
    Width(50).
    Padding(2).
    Background(theme.BackgroundColor()).  // REQUIRED for modals
    Render()

// Box variants (7 total)
containers.BoxDefault      // Standard border
containers.BoxEmphasized   // Thick border, primary color
containers.BoxDestructive  // Error/warning color
containers.BoxSubtle       // Muted color
containers.BoxSuccess      // Success color (green)
containers.BoxWarning      // Warning color (yellow)
containers.BoxInfo         // Info color (blue)
```

---

## Primitives

### Text Component

The `Text` component provides semantic text styling with automatic theme integration.

#### Semantic Styles

| Style | Constructor | Use Case |
|-------|-------------|----------|
| `TextTitle` | `Title(text, theme)` | Page/section titles |
| `TextSubtitle` | `Subtitle(text, theme)` | Section subtitles |
| `TextBody` | `Body(text, theme)` | Regular content |
| `TextMuted` | `Muted(text, theme)` | Hints, disabled text |
| `TextError` | `ErrorText(text, theme)` | Error messages |
| `TextSuccess` | `SuccessText(text, theme)` | Success messages |
| `TextWarning` | `WarningText(text, theme)` | Warning messages |
| `TextInfo` | `InfoText(text, theme)` | Informational text |

#### Text Methods

```go
// Formatting
.Bold()                  // Make text bold
.Italic()                // Make text italic
.Width(n)                // Set fixed width

// Margins
.MarginTop(n)            // Top margin in lines
.MarginBottom(n)         // Bottom margin in lines
.MarginLeft(n)           // Left margin in chars
.MarginRight(n)          // Right margin in chars

// Color override
.Foreground(color)       // Override default color

// Render
.Render() string         // Generate styled output
```

### Badge Component

Badges display small status indicators or labels.

```go
// Key badge - compact format "[Esc → Back]"
badge := primitives.KeyBadge("Esc", "Back", theme).Render()

// Help key badge - two-part format "[Esc] Back" (for help footers)
helpBadge := primitives.HelpKeyBadge("Esc", "Back", theme).Render()

// Status badge - with success styling
status := primitives.StatusBadge("Active", theme).Render()

// Tag badge - pill-style
tag := primitives.TagBadge("golang", theme).Render()
```

#### Badge Variants

| Variant | Format | Use Case |
|---------|--------|----------|
| `BadgeDefault` | Plain text | General labels |
| `BadgeKey` | `[key → action]` | Inline keyboard hints |
| `BadgeHelpKey` | `[key] action` | Help footer shortcuts |
| `BadgeStatus` | Bordered text | Status indicators |
| `BadgeTag` | Pill-style | Tags and categories |

#### Pre-built Help Key Badges

> **⚠️ MANDATORY**: Always use pre-built badge constructors instead of `HelpKeyBadge()` directly.
> This ensures consistency across the application. If a badge doesn't exist, add it to
> `primitives/badge.go` first, then use it.

For common keyboard shortcuts:

```go
// Navigation
primitives.NavigateBadge(theme)          // "[↑↓/jk] Navigate"
primitives.NavigateHorizontalBadge(theme) // "[←/→] Navigate"
primitives.NavigateVimBadge(theme)       // "[j/k] Navigate"
primitives.ToggleBadge(theme)            // "[←→/hl] Toggle"

// Actions
primitives.SelectBadge(theme)            // "[Enter] Select"
primitives.ConfirmBadge(theme)           // "[Enter] Confirm"
primitives.SubmitBadge(theme)            // "[Enter] Submit"
primitives.ContinueBadge(theme)          // "[Enter] Continue"
primitives.ApplyBadge(theme)             // "[Enter] Apply"
primitives.CancelBadge(theme)            // "[Esc] Cancel"
primitives.BackBadge(theme)              // "[Esc] Back"

// Yes/No
primitives.YesBadge(theme)               // "[y] Yes"
primitives.NoBadge(theme)                // "[n] No"

// CRUD Operations
primitives.AddBadge(theme)               // "[a] Add"
primitives.EditBadge(theme)              // "[e] Edit"
primitives.DeleteBadge(theme)            // "[d] Delete"
primitives.SaveBadge(theme)              // "[Ctrl+S] Save"

// Form Navigation
primitives.NextBadge(theme)              // "[Tab] Next"
primitives.NextFieldBadge(theme)         // "[Tab] Next field"
primitives.PrevBadge(theme)              // "[Shift+Tab] Previous"

// Search/Filter
primitives.SearchBadge(theme)            // "[/] Search"
primitives.FilterBadge(theme)            // "[f] Filter"

// Global
primitives.QuitBadge(theme)              // "[q] Quit"
primitives.HelpBadge(theme)              // "[?] Help"
primitives.MenuBadge(theme)              // "[m] Main Menu"

// Error Recovery
primitives.RetryBadge(theme)             // "[r] Retry"
primitives.RetryEnterBadge(theme)        // "[Enter/r] Retry"
```

**❌ DON'T do this:**
```go
// BAD: Using HelpKeyBadge directly
primitives.HelpKeyBadge("m", "Menu", theme)
```

**✅ DO this:**
```go
// GOOD: Using predefined badge constructor
primitives.MenuBadge(theme)
```

#### Help Footer Rendering

Render multiple badges as a help footer:

```go
footer := primitives.RenderHelpFooter(theme,
    primitives.NavigateBadge(theme),
    primitives.SelectBadge(theme),
    primitives.BackBadge(theme),
)
```

#### Standard Footer Presets

Pre-configured footers for common contexts:

```go
primitives.RenderMenuFooter(theme)     // Navigate, Select, Help, Quit
primitives.RenderListFooter(theme)     // Navigate, Select, Back
primitives.RenderFormFooter(theme)     // Confirm, Cancel
primitives.RenderEditFooter(theme)     // Save, Cancel
primitives.RenderConfirmFooter(theme)  // Confirm, Cancel
primitives.RenderBrowseFooter(theme)   // Navigate, Select, Edit, Delete, Back
primitives.RenderExportFooter(theme)   // Navigate, Select, Confirm, Back
```

### Button Component

Interactive buttons with multiple variants.

```go
// Button variants
btn := primitives.NewButton("Submit", theme).
    Variant(primitives.ButtonPrimary).  // Primary, Secondary, Danger
    Focused(true).                       // Show focus state
    Render()

// Convenience constructors
save := primitives.PrimaryButton("Save", theme)
cancel := primitives.SecondaryButton("Cancel", theme)
delete := primitives.DangerButton("Delete", theme).Disabled(true)
```

### ButtonGroup Component

Keyboard-navigable button groups.

```go
group := primitives.NewButtonGroup(theme).
    AddButton("Cancel", false).
    AddButton("Confirm", true).  // focused
    Render()
```

---

## Containers

### Box Component

The `Box` component provides bordered containers for modals, cards, and panels.

#### Box Methods

```go
// Content
.Content(string)         // Set box content
.Title(string)           // Set box title (rendered bold at top)

// Styling
.Variant(BoxVariant)     // BoxDefault, BoxEmphasized, BoxDestructive, BoxSubtle
.Padding(int)            // Internal padding

// Dimensions
.Width(int)              // Fixed width (0 = auto)
.Height(int)             // Fixed height (0 = auto)

// Background (REQUIRED for modals)
.Background(color)       // Set solid background color

// Effects
.WithShadow()            // Add drop shadow

// Render
.Render() string         // Generate output
```

#### Box Variants

| Variant | Border | Color | Use Case |
|---------|--------|-------|----------|
| `BoxDefault` | Rounded | Border color | Standard containers |
| `BoxEmphasized` | Thick | Primary color | Important content |
| `BoxDestructive` | Rounded | Error color | Delete confirmations |
| `BoxSubtle` | Rounded | Muted color | Background panels |
| `BoxSuccess` | Rounded | Success color | Success confirmations |
| `BoxWarning` | Rounded | Warning color | Caution messages |
| `BoxInfo` | Rounded | Info/Accent color | Informational content |

### Overlay Component

Centered modal overlay with background dimming.

```go
// NewOverlay takes width and height parameters (NOT theme)
overlay := containers.NewOverlay(80, 24).  // width, height
    Content(modalContent).
    Dimmed().                               // Enable background dimming
    Render()

// With custom dim character
overlay := containers.NewOverlay(80, 24).
    Content(modalContent).
    DimmedWith('░').                        // Custom dim character
    Render()
```

**Note**: The Overlay constructor takes `(width, height int)` parameters for terminal dimensions, not a theme. It does not have fluent `.Width()` or `.Height()` methods - dimensions are set in the constructor.

---

## Feedback Components

The `feedback` package provides modal and user feedback components.

### Modal Component

The Modal system provides 5 modal types for different feedback scenarios.

```go
import "github.com/baphled/kariya/internal/cli/uikit/feedback"

// Error modal - for error messages
modal := feedback.NewErrorModal("Operation Failed", "Could not save file: permission denied")

// Loading modal - with spinner
modal := feedback.NewLoadingModal("Processing...", true)  // true = cancellable

// Progress modal - with progress bar
modal := feedback.NewProgressModal("Uploading", "Uploading file...", 0.65)  // 65% complete

// Success modal - for success confirmation
modal := feedback.NewSuccessModal("File saved successfully!")

// Warning modal - for warnings
modal := feedback.NewWarningModal("Caution", "This action cannot be undone")

// Render with terminal dimensions
rendered := modal.Render(terminalWidth, terminalHeight)
```

#### Modal Types

| Type | Icon | Use Case |
|------|------|----------|
| `ModalError` | Warning icon | Error messages, failures |
| `ModalLoading` | Spinner | Async operations, loading states |
| `ModalProgress` | Progress bar | Multi-step operations |
| `ModalSuccess` | Checkmark | Success confirmations |
| `ModalWarning` | Warning icon | Caution messages |

#### Loading Message Rotator

For long operations, rotate through messages:

```go
rotator := feedback.NewLoadingMessageRotator([]string{
    "Analyzing data...",
    "Processing events...",
    "Building timeline...",
})
modal := feedback.NewLoadingModal(rotator.GetCurrent(), true)

// In Update() to rotate messages
rotator.Rotate()
```

#### Simple Spinner

For custom loading indicators:

```go
spinner := feedback.NewSimpleSpinner()
frame := spinner.GetFrame()  // Returns current frame (e.g., "⠋")
spinner.Advance()            // Move to next frame
```

---

## Layout Components

The `layout` package provides screen layout components.

### ScreenLayout Component

`ScreenLayout` is the primary layout component (replacement for the old `StandardView`).

```go
import (
    "github.com/baphled/kariya/internal/cli/uikit/layout"
    "github.com/baphled/kariya/internal/cli/uikit/display"
)

// Create layout
logo := display.NewLogo(false, termInfo.Width)
view := layout.NewScreenLayout(termInfo).
    WithTheme(theme).
    WithLogo(logo, 2).                              // Logo with 2 lines spacing
    WithBreadcrumbs("Main Menu", "Settings").       // Navigation trail
    WithContent("Your content here").
    WithHelp("↑/k Up  ↓/j Down  Enter Select").     // Help footer
    WithFooterSeparator(true)                       // Show separator line

// Alternative: Use title instead of breadcrumbs
view := layout.NewScreenLayout(termInfo).
    WithTitle("Settings", "Configure your preferences").
    WithContent("Content here")

// Render
output := view.Render()

// With modal overlay
view.ShowModalOverlay(modal)
output := view.Render()
```

#### ScreenLayout Methods

```go
// Content
.WithContent(content string) *ScreenLayout
.WithContentStyle(style lipgloss.Style) *ScreenLayout

// Header options
.WithLogo(logo LogoRenderer, spacing int) *ScreenLayout
.WithBreadcrumbs(crumbs ...string) *ScreenLayout
.WithTitle(title, subtitle string) *ScreenLayout
.WithTheme(theme themes.Theme) *ScreenLayout

// Footer options
.WithHelp(helpText string) *ScreenLayout
.WithFooterSeparator(show bool) *ScreenLayout

// Layout options
.SetUseFullWidth(full bool) *ScreenLayout

// Modal overlay
.ShowModalOverlay(modal ModalRenderer) *ScreenLayout

// Render
.Render() string
```

---

## Navigation Components

The `navigation` package provides navigation UI components.

### BreadcrumbBar Component

```go
import "github.com/baphled/kariya/internal/cli/uikit/navigation"

bar := navigation.NewBreadcrumbBar(80, false).  // width, boxed
    WithTheme(theme).
    AddCrumb(navigation.Breadcrumb{Label: "Home", Icon: navigation.IconHome}).
    AddCrumb(navigation.Breadcrumb{Label: "Settings", Icon: navigation.IconConfigure})

rendered := bar.View()
```

#### Icon Constants

| Icon | Constant | Use Case |
|------|----------|----------|
| 🏠 | `IconHome` | Main menu |
| ✏️ | `IconCaptureEvent` | Event capture |
| 📅 | `IconTimeline` | Timeline views |
| 📄 | `IconGenerateCV` | CV generation |
| 💾 | `IconExport` / `IconSave` | Export / Save |
| ⚙️ | `IconConfigure` | Settings |
| 📊 | `IconBursts` | Burst management |
| 💡 | `IconFacts` | Fact management |
| 📥 | `IconImport` | Import data |
| 🏷️ | `IconMetadata` | Metadata editor |
| 🔍 | `IconSearch` | Search |
| 🔎 | `IconFilter` | Filter |
| ↕️ | `IconSort` | Sort |

---

## Display Components

The `display` package provides visual display components.

### Logo Component

```go
import "github.com/baphled/kariya/internal/cli/uikit/display"

// Static logo
logo := display.NewLogo(false, 80).  // animated=false, width=80
    WithTheme(theme).
    WithTagline("Career Event Management System").
    WithVersion("v1.0.0")

rendered := logo.ViewStatic()

// Animated logo (for startup screens)
logo := display.NewLogo(true, 80)  // animated=true
cmd := logo.Init()
// In Update(): logo.Update(msg)
// In View(): logo.View()
```

---

## Theme Integration

### Automatic Theme Access

All UIKit components embed `theme.Aware`, providing automatic theme access:

```go
type MyComponent struct {
    theme.Aware  // Embed this
    // ... other fields
}

// Now you have access to:
comp.Theme()           // Get current theme
comp.PrimaryColor()    // Get primary color
comp.ErrorColor()      // Get error color
comp.BackgroundColor() // Get background color
// ... all color helpers
```

### Nil-Safe Theme Access

The `theme.Aware` type provides **automatic nil safety**:

```go
// Theme() returns default theme if none set
func (a *Aware) Theme() Theme {
    if a.theme == nil {
        return Default()  // Returns default theme
    }
    return a.theme
}
```

### Available Color Helpers

```go
theme.PrimaryColor()       // Main accent color
theme.SecondaryColor()     // Secondary accent
theme.AccentColor()        // Tertiary accent
theme.BackgroundColor()    // Background color
theme.ForegroundColor()    // Default text color
theme.MutedColor()         // Disabled/hint text
theme.SuccessColor()       // Success status
theme.WarningColor()       // Warning status
theme.ErrorColor()         // Error status
theme.InfoColor()          // Info status
theme.BorderColor()        // Default border
theme.BorderActiveColor()  // Active/focused border
```

---

## Migration Guide

### From Direct Lipgloss to UIKit

#### Step 1: Update Imports

```go
// BEFORE
import (
    "github.com/baphled/kariya/internal/cli/styles"
    "github.com/charmbracelet/lipgloss"
)

// AFTER
import (
    "github.com/baphled/kariya/internal/cli/themes"
    "github.com/baphled/kariya/internal/cli/uikit/containers"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
)
```

#### Step 2: Replace Box/Container Styling

```go
// BEFORE (direct lipgloss)
lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(styles.ColorBorder).
    Background(styles.ColorBackground).
    Padding(1, 2).
    Render(content)

// AFTER (UIKit)
theme := themes.NewDefaultTheme()
containers.NewBox(theme).
    Content(content).
    Padding(2).
    Background(theme.BackgroundColor()).
    Render()
```

#### Step 3: Replace Text Styling

```go
// BEFORE (direct lipgloss)
lipgloss.NewStyle().
    Bold(true).
    Foreground(styles.ColorTextPrimary).
    Render("Title")

// AFTER (UIKit)
primitives.Title("Title", theme).Render()

// BEFORE (muted text)
lipgloss.NewStyle().
    Foreground(styles.ColorTextMuted).
    Render("hint")

// AFTER (UIKit)
primitives.Muted("hint", theme).Render()
```

#### Step 4: Remove Unused Imports

After migration, remove:
- `"github.com/baphled/kariya/internal/cli/styles"` (if no longer used)
- `"github.com/charmbracelet/lipgloss"` (if no longer used)

---

## Mandatory Patterns

### Pattern 1: Nil Theme Guard in View()

**ALL components that use theme in View() MUST have a nil guard:**

```go
func (m *MyModal) View() string {
    if !m.visible {
        return ""
    }

    // REQUIRED: Nil theme guard
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }

    // Now safe to use theme
    return containers.NewBox(theme).
        Content(m.content).
        Background(theme.BackgroundColor()).
        Render()
}
```

**Why?** Components may be constructed with nil themes in tests or during initialization. Without this guard, `View()` will panic.

### Pattern 2: Solid Background for Modals

**ALL modal overlays MUST have a solid background:**

```go
// REQUIRED for modal boxes
containers.NewBox(theme).
    Content(content).
    Background(theme.BackgroundColor()).  // REQUIRED - prevents transparency
    Render()
```

**Why?** Without a solid background, the modal shows the underlying view through it, causing visual artifacts.

### Pattern 3: Theme in Component Fields

**Components that use theme in multiple methods should store it:**

```go
type MyModal struct {
    theme   themes.Theme  // Store theme
    visible bool
    // ...
}

func NewMyModal(theme themes.Theme) *MyModal {
    return &MyModal{
        theme: theme,  // Store for later use
    }
}
```

### Pattern 4: Default Theme Fallback

**When theme may be nil, use default:**

```go
// In constructors
func NewMyComponent(theme themes.Theme) *MyComponent {
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    return &MyComponent{theme: theme}
}

// Or inline in View()
theme := themes.NewDefaultTheme()  // Always safe
```

### Pattern 5: Predefined Badge Constructors

**ALWAYS use predefined badge constructors for help footers:**

```go
// ❌ BAD: Using HelpKeyBadge directly
helpFooter := primitives.RenderHelpFooter(theme,
    primitives.HelpKeyBadge("m", "Menu", theme),      // DON'T
    primitives.HelpKeyBadge("Enter", "Confirm", theme), // DON'T
)

// ✅ GOOD: Using predefined badge constructors
helpFooter := primitives.RenderHelpFooter(theme,
    primitives.MenuBadge(theme),    // DO
    primitives.ConfirmBadge(theme), // DO
)
```

**Why?**
- Ensures consistent key labels across the application (e.g., "Main Menu" not "Menu")
- Single source of truth for keyboard shortcuts
- Easier to update globally if key bindings change
- Prevents typos and inconsistencies

**If a badge doesn't exist:**
1. Add it to `internal/cli/uikit/primitives/badge.go`
2. Follow existing naming convention: `<Action>Badge(th theme.Theme)`
3. Use the new constructor in your code

---

## Component Checklist

Use this checklist when creating or modifying components:

### New Component Checklist

- [ ] Uses UIKit primitives for text styling (not direct lipgloss)
- [ ] Uses UIKit containers for boxes/borders (not direct lipgloss)
- [ ] Uses predefined badge constructors (not `HelpKeyBadge()` directly)
- [ ] Stores theme as field if used in multiple methods
- [ ] Has nil theme guard in View() method
- [ ] Uses `themes.NewDefaultTheme()` as fallback
- [ ] Modal boxes have solid background via `.Background()`
- [ ] Removes unused lipgloss/styles imports
- [ ] Tests pass with nil theme

### Migration Checklist

- [ ] Replace `lipgloss.NewStyle().Border(...)` with `containers.NewBox()`
- [ ] Replace text styling with `primitives.Title/Body/Muted/etc.`
- [ ] Replace `HelpKeyBadge()` calls with predefined badge constructors
- [ ] Add nil theme guard to View()
- [ ] Add `.Background(theme.BackgroundColor())` to modal boxes
- [ ] Update imports (add uikit, potentially remove lipgloss/styles)
- [ ] Run tests to verify no regressions
- [ ] Visual test to verify appearance

---

## API Reference

### primitives.Text

```go
// Constructors
NewText(content string, theme Theme) *Text
Title(content string, theme Theme) *Text
Subtitle(content string, theme Theme) *Text
Body(content string, theme Theme) *Text
Muted(content string, theme Theme) *Text
ErrorText(content string, theme Theme) *Text
SuccessText(content string, theme Theme) *Text
WarningText(content string, theme Theme) *Text
InfoText(content string, theme Theme) *Text

// Methods
(t *Text) Style(style TextStyle) *Text
(t *Text) Bold() *Text
(t *Text) Italic() *Text
(t *Text) Width(width int) *Text
(t *Text) MarginTop(lines int) *Text
(t *Text) MarginBottom(lines int) *Text
(t *Text) MarginLeft(chars int) *Text
(t *Text) MarginRight(chars int) *Text
(t *Text) Foreground(color lipgloss.Color) *Text
(t *Text) Render() string
```

### containers.Box

```go
// Constructor
NewBox(theme Theme) *Box

// Methods
(b *Box) Content(content string) *Box
(b *Box) Title(title string) *Box
(b *Box) Variant(variant BoxVariant) *Box
(b *Box) Width(width int) *Box
(b *Box) Height(height int) *Box
(b *Box) Padding(padding int) *Box
(b *Box) Background(color lipgloss.Color) *Box
(b *Box) WithShadow() *Box
(b *Box) Render() string
```

### theme.Aware

```go
// Embed in your struct
type MyComponent struct {
    theme.Aware
}

// Methods available after embedding
(a *Aware) SetTheme(t Theme)
(a *Aware) Theme() Theme
(a *Aware) PrimaryColor() lipgloss.Color
(a *Aware) SecondaryColor() lipgloss.Color
(a *Aware) AccentColor() lipgloss.Color
(a *Aware) ErrorColor() lipgloss.Color
(a *Aware) SuccessColor() lipgloss.Color
(a *Aware) WarningColor() lipgloss.Color
(a *Aware) InfoColor() lipgloss.Color
(a *Aware) BorderColor() lipgloss.Color
(a *Aware) BackgroundColor() lipgloss.Color
(a *Aware) MutedColor() lipgloss.Color
```

---

## Related Documentation

- [Theme Customization Guide](./THEME_CUSTOMIZATION_GUIDE.md) - Complete theme system docs
- [Modal Patterns](./MODAL_PATTERNS.md) - Modal usage patterns
- [TUI Developer Guide](./TUI_DEVELOPER_GUIDE.md) - General TUI development
- [TUI Standards](./TUI_STANDARDS.md) - Design principles and guidelines

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.1 | 2026-01-20 | Added feedback, layout, navigation, display packages; fixed Overlay constructor; added Box variants (Success, Warning, Info); documented Badge help footer features |
| 1.0 | 2026-01-19 | Initial release with primitives, containers, migration guide |
