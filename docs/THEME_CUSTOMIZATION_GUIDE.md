# Theme Customization Guide

**Version**: 1.0  
**Created**: 2026-01-10  
**Status**: Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Theme Architecture](#theme-architecture)
4. [Using Themes in Intents](#using-themes-in-intents)
5. [Creating Custom Themes](#creating-custom-themes)
6. [Color Palette Reference](#color-palette-reference)
7. [StyleSet Reference](#styleset-reference)
8. [Bubbles Integration](#bubbles-integration)
9. [Best Practices](#best-practices)
10. [API Reference](#api-reference)

---

## Overview

KaRiya's theme system provides:

- **Runtime theme switching** - Change themes without restarting
- **Consistent styling** - All UI components use themed colors
- **Type-safe API** - Compile-time safety for theme access
- **Backwards compatibility** - Default theme matches original colors
- **Pre-composed styles** - Ready-to-use Lipgloss styles

### Package Locations

**Primary Theme Package** (`internal/cli/themes/`):
```
internal/cli/themes/
├── theme.go       # Theme interface and ColorPalette
├── styles.go      # StyleSet and GenerateStyles
├── default.go     # Default theme (KaRiya colors)
├── manager.go     # ThemeManager for registration/switching
├── detector.go    # Terminal detection utilities
├── bubbles.go     # Bubbles component integration
├── glamour.go     # Markdown rendering with themes
└── loading.go     # Themed loading views
```

**UIKit Theme Package** (`internal/cli/uikit/theme/`):
```
internal/cli/uikit/theme/
├── theme.go       # Re-exports themes.Theme + Default() helper
├── aware.go       # ThemeAware embeddable struct
└── theme_test.go  # Theme package tests
```

The UIKit theme package provides a convenience layer for UIKit components:
- Re-exports `themes.Theme` interface for easy imports
- Provides `Default()` function to get the default theme
- Provides `Aware` embeddable struct for theme-aware components

---

## Quick Start

### Accessing Theme in an Intent

```go
func (i *MyIntent) View() string {
    // Get active theme (nil-safe)
    theme := i.Theme()
    if theme == nil {
        // Fallback to default styling
        return i.renderWithDefaults()
    }
    
    // Use themed styles
    card := theme.Styles().CardBase.Render("Content")
    return card
}
```

### Using Semantic Colors

```go
func (i *MyIntent) viewError(msg string) string {
    theme := i.Theme()
    
    // Using color helpers (nil-safe with fallbacks)
    errorColor := i.getErrorColor()
    
    style := lipgloss.NewStyle().
        Foreground(errorColor).
        Bold(true)
    
    return style.Render(msg)
}

// Helper method pattern
func (i *MyIntent) getErrorColor() lipgloss.Color {
    if theme := i.Theme(); theme != nil {
        return theme.ErrorColor()
    }
    return styles.ColorError // fallback
}
```

---

## Theme Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                      ThemeManager                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Default   │  │   Custom    │  │   Custom    │  ...     │
│  │   Theme     │  │   Theme A   │  │   Theme B   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│         │                │                │                  │
│         └────────────────┼────────────────┘                  │
│                          ▼                                   │
│                   Active Theme                               │
│                          │                                   │
│         ┌────────────────┼────────────────┐                  │
│         ▼                ▼                ▼                  │
│   ColorPalette      StyleSet       Metadata                  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
                    TUI Components
```

### Theme Interface

```go
type Theme interface {
    // Metadata
    Name() string
    Description() string
    Author() string
    IsDark() bool

    // Color Access
    Palette() *ColorPalette

    // Pre-composed Styles
    Styles() *StyleSet

    // Semantic Color Helpers
    PrimaryColor() lipgloss.Color
    SecondaryColor() lipgloss.Color
    AccentColor() lipgloss.Color
    BackgroundColor() lipgloss.Color
    ForegroundColor() lipgloss.Color
    MutedColor() lipgloss.Color
    SuccessColor() lipgloss.Color
    WarningColor() lipgloss.Color
    ErrorColor() lipgloss.Color
    InfoColor() lipgloss.Color
    BorderColor() lipgloss.Color
    BorderActiveColor() lipgloss.Color
}
```

### Theme Propagation

Themes flow through the application via:

1. **ThemeManager** - Created in app initialization
2. **IntentRouter** - Propagates theme to intents on activation
3. **BaseIntent** - Provides `Theme()` method to all intents

```go
// Router propagates theme to new intents
func (r *DefaultIntentRouter) ActivateIntent(name string, ...) (tea.Cmd, error) {
    intent := factory()
    
    // Propagate theme manager
    if r.themeManager != nil {
        if setter, ok := intent.(interface{ SetThemeManager(*themes.ThemeManager) }); ok {
            setter.SetThemeManager(r.themeManager)
        }
    }
    
    return intent.Init(), nil
}
```

---

## Using Themes in Intents

### Pattern 0: ThemeAware Embedding (UIKit Components)

For UIKit components, use the `theme.Aware` embeddable struct for automatic theme access:

```go
import "github.com/baphled/kariya/internal/cli/uikit/theme"

type MyButton struct {
    theme.Aware  // Embed for theme access
    label string
}

func NewMyButton(label string) *MyButton {
    btn := &MyButton{label: label}
    // Theme defaults to Default() if not set
    return btn
}

func (b *MyButton) Render() string {
    // Use inherited color getters (nil-safe, auto-fallback)
    style := lipgloss.NewStyle().
        Foreground(b.PrimaryColor()).
        Background(b.BackgroundColor())
    return style.Render(b.label)
}

// Optionally override theme
btn := NewMyButton("Click Me")
btn.SetTheme(customTheme)
```

**Available Color Getters via `theme.Aware`**:
- `PrimaryColor()`, `SecondaryColor()`, `AccentColor()`
- `ErrorColor()`, `SuccessColor()`, `WarningColor()`, `InfoColor()`
- `BorderColor()`, `BackgroundColor()`, `MutedColor()`

**Benefits**:
- ✅ Automatic nil-safety (falls back to default theme)
- ✅ Consistent color access across all UIKit components
- ✅ No need for manual nil checks
- ✅ Theme propagation via `SetTheme()`

### Pattern 1: Helper Methods (Recommended)

Create helper methods in your intent for commonly used colors:

```go
type MyIntent struct {
    *BaseIntent
    // ... other fields
}

// getCardStyle returns themed card style with fallback
func (i *MyIntent) getCardStyle() lipgloss.Style {
    if theme := i.Theme(); theme != nil {
        return theme.Styles().CardBase
    }
    return lipgloss.NewStyle().
        Padding(1, 2).
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder)
}

// getPrimaryColor returns primary text color with fallback
func (i *MyIntent) getPrimaryColor() lipgloss.Color {
    if theme := i.Theme(); theme != nil {
        return theme.ForegroundColor()
    }
    return styles.ColorTextPrimary
}

// getAccentColor returns accent color with fallback
func (i *MyIntent) getAccentColor() lipgloss.Color {
    if theme := i.Theme(); theme != nil {
        return theme.PrimaryColor()
    }
    return styles.ColorAccentTeal
}
```

### Pattern 2: Direct Theme Access

For one-off styling or complex views:

```go
func (i *MyIntent) viewDetails() string {
    var content strings.Builder
    
    if theme := i.Theme(); theme != nil {
        // Use pre-composed styles
        header := theme.Styles().HeaderSection.Render("Details")
        content.WriteString(header)
        
        // Use semantic colors
        style := lipgloss.NewStyle().
            Foreground(theme.ForegroundColor()).
            Background(theme.Palette().BackgroundCard)
        content.WriteString(style.Render("Body text"))
    }
    
    return content.String()
}
```

### Pattern 3: Themed Tables

For tables using bubbles/table:

```go
func (i *MyIntent) Init() tea.Cmd {
    // Apply themed table styles
    if theme := i.Theme(); theme != nil {
        i.table.SetStyles(themes.NewThemedTableStyles(theme))
    }
    return nil
}
```

---

## Creating Custom Themes

### Step 1: Define Color Palette

```go
// internal/cli/themes/mytheme.go
package themes

import "github.com/charmbracelet/lipgloss"

var MyPalette = ColorPalette{
    // Background colors
    Background:     lipgloss.Color("#1e1e2e"),
    BackgroundAlt:  lipgloss.Color("#313244"),
    BackgroundCard: lipgloss.Color("#45475a"),

    // Foreground colors
    Foreground:      lipgloss.Color("#cdd6f4"),
    ForegroundDim:   lipgloss.Color("#a6adc8"),
    ForegroundMuted: lipgloss.Color("#6c7086"),

    // Accent colors
    Primary:   lipgloss.Color("#89b4fa"),
    Secondary: lipgloss.Color("#a6e3a1"),
    Tertiary:  lipgloss.Color("#cba6f7"),

    // Status colors
    Success: lipgloss.Color("#a6e3a1"),
    Warning: lipgloss.Color("#f9e2af"),
    Error:   lipgloss.Color("#f38ba8"),
    Info:    lipgloss.Color("#89dceb"),

    // Border colors
    Border:       lipgloss.Color("#45475a"),
    BorderActive: lipgloss.Color("#89b4fa"),
    BorderError:  lipgloss.Color("#f38ba8"),

    // Special colors
    Selection: lipgloss.Color("#45475a"),
    Highlight: lipgloss.Color("#585b70"),
    Link:      lipgloss.Color("#89dceb"),
}
```

### Step 2: Create Theme

```go
func NewMyTheme() Theme {
    return NewBaseTheme(
        "my-theme",
        "My Custom Theme - A beautiful color scheme",
        "Your Name",
        true, // isDark
        &MyPalette,
    )
}
```

### Step 3: Register Theme

```go
// In your app initialization
func initThemes(tm *themes.ThemeManager) {
    // Register custom theme
    if err := tm.Register(themes.NewMyTheme()); err != nil {
        log.Printf("Failed to register theme: %v", err)
    }
}
```

### Step 4: Switch Theme

```go
// Switch to custom theme
if err := themeManager.SetActive("my-theme"); err != nil {
    log.Printf("Failed to switch theme: %v", err)
}
```

---

## Color Palette Reference

### Background Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Background` | Primary background | `#1a1f2e` |
| `BackgroundAlt` | Alternate/contrast background | `#242936` |
| `BackgroundCard` | Card/panel background | `#2d3346` |

### Foreground Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Foreground` | Primary text | `#c7ccd1` |
| `ForegroundDim` | Secondary/dimmed text | `#8b92a0` |
| `ForegroundMuted` | Muted/disabled text | `#5e6673` |

### Accent Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Primary` | Primary action/brand color | `#5fb3b3` (teal) |
| `Secondary` | Secondary actions | `#6cb56c` (green) |
| `Tertiary` | Tertiary/highlight | `#a99bd1` (purple) |

### Status Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Success` | Success state | `#6cb56c` |
| `Warning` | Warning state | `#d9a66c` |
| `Error` | Error state | `#d76e6e` |
| `Info` | Info state | `#6ab0d3` |

### Border Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Border` | Default borders | `#3d4454` |
| `BorderActive` | Focused/active borders | `#5fb3b3` |
| `BorderError` | Error state borders | `#d76e6e` |

### Special Colors

| Field | Purpose | Default Value |
|-------|---------|---------------|
| `Selection` | Selected item background | `#3d4454` |
| `Highlight` | Highlighted text background | `#4d5566` |
| `Link` | Hyperlink color | `#6ab0d3` |

---

## StyleSet Reference

StyleSet provides pre-composed Lipgloss styles. All styles are automatically generated from the color palette.

### Card Styles

```go
theme.Styles().CardBase     // Standard card with border and padding
theme.Styles().CardHeader   // Bold header text
theme.Styles().CardContent  // Normal content text
theme.Styles().CardFooter   // Dimmed footer text
```

### Button Styles

```go
theme.Styles().ButtonBase      // Base button style
theme.Styles().ButtonPrimary   // Primary action button
theme.Styles().ButtonSecondary // Secondary action button
theme.Styles().ButtonFocused   // Focused button with thick border
theme.Styles().ButtonDisabled  // Disabled/faint button
```

### Input Styles

```go
theme.Styles().InputBase    // Base input style
theme.Styles().InputFocused // Focused input with active border
theme.Styles().InputError   // Input with error border
theme.Styles().InputLabel   // Input label (bold, dimmed)
theme.Styles().InputHint    // Input hint (italic, muted)
```

### List Styles

```go
theme.Styles().ListItem         // Normal list item
theme.Styles().ListItemSelected // Selected item with background
theme.Styles().ListItemFocused  // Focused item with left border
```

### Message Boxes

```go
theme.Styles().ErrorBox   // Error message container
theme.Styles().WarningBox // Warning message container
theme.Styles().SuccessBox // Success message container
theme.Styles().InfoBox    // Info message container
```

### Text Styles

```go
theme.Styles().ErrorText   // Bold error text
theme.Styles().WarningText // Bold warning text
theme.Styles().SuccessText // Bold success text
theme.Styles().InfoText    // Bold info text
theme.Styles().MutedText   // Muted/disabled text
```

### Header Styles

```go
theme.Styles().HeaderMain       // Main header (primary, bold)
theme.Styles().HeaderSection    // Section header (bold)
theme.Styles().HeaderSubsection // Subsection header (dimmed, bold)
```

### Key Badge Styles

```go
theme.Styles().KeyBadge     // Keyboard shortcut badge
theme.Styles().KeyBadgeHint // Key hint text
```

---

## Bubbles Integration

The theme system integrates with Charmbracelet bubbles components.

### Themed Table

```go
import "github.com/baphled/kariya/internal/cli/themes"

// Create themed table styles
tableStyles := themes.NewThemedTableStyles(theme)
myTable.SetStyles(tableStyles)

// Or apply to existing table
myTable = themes.ApplyThemeToTable(myTable, theme)
```

### Themed List

```go
// Create themed list styles
listStyles := themes.NewThemedListStyles(theme)
myList.Styles = listStyles

// Create themed delegate for item styling
delegate := themes.NewThemedListDelegate(theme)
myList.SetDelegate(delegate)

// Or apply to existing list
myList = themes.ApplyThemeToList(myList, theme)
```

### Themed Progress

```go
// Create themed progress bar with gradient
progressBar := themes.NewThemedProgress(theme)
```

### Themed Spinner

```go
// Create themed spinner
spinner := themes.NewThemedSpinner(theme)
```

### Themed Help

```go
// Create themed help styles
helpStyles := themes.NewThemedHelpStyles(theme)
myHelp.Styles = helpStyles
```

### Themed Huh Forms

The theme system integrates with Charmbracelet's huh library for form handling.

```go
import (
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/cli/themes"
    "github.com/charmbracelet/huh"
)

// Create a themed form using the forms package
form := forms.NewThemedForm(theme,
    huh.NewGroup(
        huh.NewInput().
            Key("name").
            Title("Full Name").
            Description("Enter your full name"),
        huh.NewSelect[string]().
            Key("role").
            Title("Role").
            Options(
                huh.NewOption("Developer", "dev"),
                huh.NewOption("Designer", "design"),
            ),
    ),
)

// Or generate a huh.Theme directly
huhTheme := themes.GenerateHuhTheme(theme)
form := huh.NewForm(groups...).WithTheme(huhTheme)
```

The themed form will have:
- **Focused fields**: Active border color, bold titles, visible selection indicators
- **Blurred fields**: Muted colors, less prominent borders
- **Error states**: Error colors for validation messages
- **Buttons**: Primary color for focused, muted for blurred
- **Help text**: Consistent with rest of TUI styling

---

## Best Practices

### 1. Always Provide Fallbacks

```go
// Good - handles nil theme
func (i *MyIntent) getErrorColor() lipgloss.Color {
    if theme := i.Theme(); theme != nil {
        return theme.ErrorColor()
    }
    return styles.ColorError
}

// Bad - will panic if theme is nil
func (i *MyIntent) getErrorColor() lipgloss.Color {
    return i.Theme().ErrorColor() // PANIC if nil!
}
```

### 2. Use Helper Methods for Repeated Colors

```go
// Good - reusable helper
func (i *MyIntent) viewList() string {
    accentColor := i.getAccentColor()
    // ... use accentColor multiple times
}

// Avoid - repeated nil checks
func (i *MyIntent) viewList() string {
    if theme := i.Theme(); theme != nil {
        // ...
    }
    // Same check repeated in another method
}
```

### 3. Use Pre-composed Styles

```go
// Good - use pre-composed styles
card := theme.Styles().CardBase.Render(content)

// Avoid - recreating styles manually
card := lipgloss.NewStyle().
    Padding(1, 2).
    BorderStyle(lipgloss.RoundedBorder()).
    // ... many more properties
    Render(content)
```

### 4. Theme Tables in Init()

```go
func (i *MyIntent) Init() tea.Cmd {
    // Apply theme once during initialization
    if theme := i.Theme(); theme != nil {
        i.table.SetStyles(themes.NewThemedTableStyles(theme))
    }
    return nil
}
```

### 5. Avoid Style.Copy()

```go
// Good - use direct assignment (styles are immutable)
cardStyle := theme.Styles().CardBase.BorderForeground(theme.ErrorColor())

// Deprecated - don't use Copy()
cardStyle := theme.Styles().CardBase.Copy().BorderForeground(theme.ErrorColor())
```

---

## API Reference

### ThemeManager

```go
// Create manager with default theme
tm := themes.NewThemeManager()

// Register a new theme
err := tm.Register(myTheme)

// Set active theme
err := tm.SetActive("theme-name")

// Get active theme
theme := tm.Active()

// Get specific theme
theme, err := tm.Get("theme-name")

// List all themes
names := tm.List()

// Subscribe to theme changes
id := tm.OnChange(func(t Theme) {
    // Handle theme change
})

// Unsubscribe
tm.RemoveChangeCallback(id)

// Quick access to active styles/palette
styles := tm.Styles()
palette := tm.Palette()
```

### BaseIntent Theme Methods

```go
// Get theme manager
tm := intent.GetThemeManager()

// Set theme manager
intent.SetThemeManager(tm)

// Get active theme (nil if no manager)
theme := intent.Theme()
```

### Terminal Detection

```go
// Detect terminal color depth
depth := themes.DetectColorDepth()
// Returns: ColorDepth16, ColorDepth256, or ColorDepthTrue

// Detect dark/light mode
isDark := themes.DetectDarkMode()
```

---

## UIKit Theme Integration

The UIKit package provides theme-aware components through the `uikit/theme` package.

### Importing UIKit Theme

```go
import (
    "github.com/baphled/kariya/internal/cli/uikit/theme"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
)
```

### Using ThemeAware in Custom Components

```go
type MyComponent struct {
    theme.Aware  // Embed for automatic theme access
    // ... other fields
}

func (c *MyComponent) View() string {
    // Color getters are nil-safe and auto-fallback to default
    return lipgloss.NewStyle().
        Foreground(c.PrimaryColor()).
        Render("My content")
}
```

### Theme Propagation

```go
// Set theme on component
component.SetTheme(myTheme)

// Get current theme (never nil - returns default if unset)
currentTheme := component.Theme()
```

### UIKit Primitives with Themes

All UIKit primitives accept a theme parameter:

```go
// Text with theme
text := primitives.NewText(theme).
    Title("Header").
    Render()

// Button with theme
btn := primitives.NewButton(theme).
    Label("Submit").
    Variant(primitives.ButtonPrimary).
    Render()

// Badge with theme
badge := primitives.NewBadge(theme).
    Text("NEW").
    Variant(primitives.BadgeStatus).
    Render()
```

---

## See Also

- [UIKit Guide](UIKIT_GUIDE.md) - Complete UIKit component documentation
- [TUI Standards](TUI_STANDARDS.md) - Overall TUI design guidelines
- [TUI Developer Guide](TUI_DEVELOPER_GUIDE.md) - Creating TUI components
- [TUI Visual Overhaul Spec](TUI_VISUAL_OVERHAUL_SPEC.md) - Original design specification
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss) - Styling library
