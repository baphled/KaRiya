---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Style Usage Guide

## Overview

This guide explains how to use the centralized style system in KaRiya. All colors, styles, and spacing constants are defined in `internal/cli/styles/` and exported through `constants_export.go` for consistent use throughout the application.

## Color Palette

### Background Colors

Use these colors for background surfaces:

```go
// Primary background - use for main content area
ColorBackground = "#1a1f2e"

// Alternate background - use for contrast/sections
ColorBackgroundAlt = "#242936"

// Card/panel background - use for cards, modals, containers
ColorBackgroundCard = "#2d3346"
```

**When to use:**
- `ColorBackground`: Main content area, list backgrounds
- `ColorBackgroundAlt`: Section backgrounds, alternating rows
- `ColorBackgroundCard`: Cards, panels, modals, containers

### Accent Colors (Muted & Professional)

Use these colors for interactive elements and status indication:

```go
// Primary action color - use for buttons, links, focus states
ColorAccentTeal = "#5fb3b3"

// Success/confirmation - use for positive actions/status
ColorAccentGreen = "#6cb56c"

// Selection/highlight - use for selected items, highlights
ColorAccentPurple = "#a99bd1"
```

**When to use:**
- `ColorAccentTeal`: Primary buttons, focused borders, active states
- `ColorAccentGreen`: Success status, confirmation actions
- `ColorAccentPurple`: Selected list items, highlighted text

### Text Colors

Use these colors for text rendering:

```go
// Primary text - use for main content
ColorTextPrimary = "#c7ccd1"

// Secondary text - use for labels, descriptions
ColorTextSecondary = "#8b92a0"

// Muted text - use for hints, disabled text
ColorTextMuted = "#5e6673"
```

**When to use:**
- `ColorTextPrimary`: Main content, body text, primary information
- `ColorTextSecondary`: Labels, descriptions, secondary information
- `ColorTextMuted`: Hints, disabled text, tertiary information

### Status Colors

Use these colors to indicate status or messages:

```go
// Error - use for errors, failures, destructive actions
ColorError = "#d76e6e"

// Warning - use for warnings, cautions
ColorWarning = "#d9a66c"

// Success - use for success messages, confirmations
ColorSuccess = "#6cb56c"

// Info - use for information messages
ColorInfo = "#6ab0d3"
```

**When to use:**
- `ColorError`: Error messages, error borders, delete confirmations
- `ColorWarning`: Warning messages, caution indicators
- `ColorSuccess`: Success messages, completion status
- `ColorInfo`: Information messages, help text

### Border Colors

Use these colors for borders and outlines:

```go
// Default border - use for normal borders
ColorBorder = "#3d4454"

// Active/focused border - use for focused elements
ColorBorderActive = "#5fb3b3"

// Error border - use for error states
ColorBorderError = "#d76e6e"
```

**When to use:**
- `ColorBorder`: Normal borders, unfocused elements
- `ColorBorderActive`: Focused inputs, active elements
- `ColorBorderError`: Error inputs, validation failures

## Using Colors in Code

### Correct Usage

```go
import "github.com/baphled/kariya/internal/cli/styles"

// Use getter functions
color := styles.GetColorAccentTeal()

// Or access constants directly
style := lipgloss.NewStyle().
    Foreground(styles.ColorTextPrimary).
    Background(styles.ColorBackgroundCard)

// Use exported styles
buttonStyle := styles.GetButtonPrimary()
inputStyle := styles.GetInputFocused()
```

### Incorrect Usage ❌

```go
// DON'T: Inline hex colors
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#5fb3b3"))

// DON'T: Create colors outside styles package
myColor := lipgloss.Color("#c7ccd1")

// DON'T: Use unexported style variables from other packages
// (Always go through styles package)
```

## Pre-built Styles

### Button Styles

```go
// Base button - foundation for all buttons
styles.GetButtonBase()

// Primary button - for main actions
styles.GetButtonPrimary()

// Secondary button - for secondary actions
styles.GetButtonSecondary()

// Focused button - for focused state
styles.GetButtonFocused()

// Disabled button - for disabled state
styles.GetButtonDisabled()
```

**Example:**
```go
primary := styles.GetButtonPrimary().Render("Submit")
secondary := styles.GetButtonSecondary().Render("Cancel")
```

### Input Field Styles

```go
// Base input - unfocused state
styles.GetInputBase()

// Focused input - when focused
styles.GetInputFocused()

// Error input - error state
styles.GetInputError()

// Input label - field labels
styles.GetInputLabel()

// Input hint - helper text
styles.GetInputHint()
```

**Example:**
```go
label := styles.GetInputLabel().Render("Email")
input := styles.GetInputFocused().Render("user@example.com")
hint := styles.GetInputHint().Render("Enter a valid email")
```

### Card Styles

```go
// Base card - foundation
styles.GetCardBase()

// Card header - header section
styles.GetCardHeader()

// Card content - content section
styles.GetCardContent()

// Card footer - footer section
styles.GetCardFooter()
```

**Example:**
```go
header := styles.GetCardHeader().Render("Event Details")
content := styles.GetCardContent().Render("Details here...")
card := styles.GetCardBase().Render(
    lipgloss.JoinVertical(lipgloss.Left, header, content),
)
```

### Modal/Dialog Styles

```go
// Base modal - foundation
styles.GetModalBase()

// Modal title - title section
styles.GetModalTitle()

// Modal message - message section
styles.GetModalMessage()

// Modal button container - buttons section
styles.GetModalButtonContainer()

// Modal instructions - helper text
styles.GetModalInstructions()

// Destructive modal - for delete/destructive actions
styles.GetModalDestructive()

// Destructive modal title - destructive title
styles.GetModalDestructiveTitle()
```

**Example:**
```go
title := styles.GetModalTitle().Render("Delete Event?")
message := styles.GetModalMessage().Render("This cannot be undone.")
modal := styles.GetModalBase().Render(
    lipgloss.JoinVertical(lipgloss.Center, title, message),
)
```

### Header Styles

```go
// Main header - page title
styles.GetHeaderMain()

// Section header - section title
styles.GetHeaderSection()

// Subsection header - subsection title
styles.GetHeaderSubsection()
```

**Example:**
```go
main := styles.GetHeaderMain().Render("Career Journal")
section := styles.GetHeaderSection().Render("Events")
```

### Message Styles

```go
// Error message
styles.GetErrorBox()
styles.GetErrorText()
styles.GetErrorHint()

// Warning message
styles.GetWarningBox()
styles.GetWarningText()
styles.GetWarningHint()

// Success message
styles.GetSuccessBox()
styles.GetSuccessText()
styles.GetSuccessHint()

// Info message
styles.GetInfoBox()
styles.GetInfoText()
styles.GetInfoHint()
```

**Example:**
```go
error := styles.GetErrorBox().Render(
    lipgloss.JoinVertical(lipgloss.Left,
        styles.GetErrorText().Render("Error"),
        styles.GetErrorHint().Render("Details"),
    ),
)
```

### List Styles

```go
// Normal list item
styles.GetListItem()

// Selected list item
styles.GetListItemSelected()

// Focused list item
styles.GetListItemFocused()
```

**Example:**
```go
item := styles.GetListItem().Render("Event 1")
selected := styles.GetListItemSelected().Render("Event 2")
focused := styles.GetListItemFocused().Render("Event 3")
```

### Tag/Badge Styles

```go
// Base tag
styles.GetTagBase()

// Selected tag
styles.GetTagSelected()

// Base badge
styles.GetBadge()

// Selected badge
styles.GetBadgeSelected()

// Focused badge
styles.GetBadgeFocused()
```

**Example:**
```go
tag := styles.GetTagBase().Render("Work")
selected := styles.GetTagSelected().Render("Career")
```

### Progress & Spinner Styles

```go
// Progress bar
styles.GetProgressBar()

// Progress text
styles.GetProgressText()

// Spinner
styles.GetSpinnerStyle()
```

**Example:**
```go
progress := styles.GetProgressBar().Render("████████░░")
text := styles.GetProgressText().Render("80%")
```

## Spacing Constants

Use spacing constants for consistent layout:

```go
spacing := styles.GetSpacingConstants()

// Padding values
spacing.PaddingHorizontalSmall  // 1
spacing.PaddingHorizontalBase   // 2
spacing.PaddingHorizontalLarge  // 3
spacing.PaddingVerticalSmall    // 0
spacing.PaddingVerticalBase     // 1

// Margin values
spacing.MarginSmall             // 1
spacing.MarginBase              // 2
spacing.MarginLarge             // 3

// Layout values
spacing.MaxContentWidth         // 120
spacing.GridGutterWidth         // 2
```

**Example:**
```go
spacing := styles.GetSpacingConstants()
style := lipgloss.NewStyle().
    Padding(spacing.PaddingVerticalBase, spacing.PaddingHorizontalBase).
    Margin(spacing.MarginSmall, 0)
```

## Layout Helper Functions

### Responsive Layout

```go
// Get maximum width based on terminal width
maxWidth := styles.MaxWidth(terminalWidth)

// Get responsive style
style := styles.ResponsiveStyle(terminalWidth)

// Get responsive card
card := styles.ResponsiveCard(terminalWidth)

// Get responsive input
input := styles.ResponsiveInput(terminalWidth)
```

**Example:**
```go
// Get terminal width from BubbleTea
width := msg.Width

// Use responsive layout
maxWidth := styles.MaxWidth(width)
card := styles.ResponsiveCard(width)
```

### Column Layouts

```go
// Two column layout
leftWidth, rightWidth := styles.TwoColumn(terminalWidth)

// Three column layout
leftWidth, centerWidth, rightWidth := styles.ThreeColumn(terminalWidth)

// Grid layout
colWidth, gutter := styles.Grid(terminalWidth, numColumns)
```

**Example:**
```go
left, right := styles.TwoColumn(m.Width)
layout := lipgloss.JoinHorizontal(lipgloss.Top,
    lipgloss.NewStyle().Width(left).Render(leftContent),
    lipgloss.NewStyle().Width(right).Render(rightContent),
)
```

### Text Alignment

```go
// Center horizontally
centered := styles.CenterHorizontal(text, width)

// Center vertically
centered := styles.CenterVertical(text, height)

// Center both
centered := styles.Center(text, width, height)

// Align left
left := styles.AlignLeft(text, width)

// Align right
right := styles.AlignRight(text, width)
```

**Example:**
```go
title := styles.CenterHorizontal("Welcome", 40)
```

## Helper Functions

### Border Functions

```go
// Add default border
style := styles.WithBorder(baseStyle)

// Add focused border
style := styles.WithFocusedBorder(baseStyle)

// Add error border
style := styles.WithErrorBorder(baseStyle)
```

**Example:**
```go
inputStyle := lipgloss.NewStyle()
focusedStyle := styles.WithFocusedBorder(inputStyle)
```

### Padding & Margin Functions

```go
// Add padding
style := styles.WithPadding(baseStyle, vertical, horizontal)

// Add margin
style := styles.WithMargin(baseStyle, vertical, horizontal)
```

**Example:**
```go
style := styles.WithPadding(baseStyle, 1, 2)
```

## Best Practices

### ✅ Do

1. **Use getter functions** for colors and styles
   ```go
   color := styles.GetColorAccentTeal()
   ```

2. **Use spacing constants** for consistency
   ```go
   spacing := styles.GetSpacingConstants()
   ```

3. **Compose styles** using existing styles as base
   ```go
   myStyle := styles.GetCardBase().Copy().
       Foreground(styles.GetColorTextPrimary())
   ```

4. **Use responsive layouts** for terminal width
   ```go
   maxWidth := styles.MaxWidth(m.Width)
   ```

5. **Document style choices** in comments
   ```go
   // Use teal accent for primary action
   buttonStyle := styles.GetButtonPrimary()
   ```

### ❌ Don't

1. **Don't use inline hex colors**
   ```go
   // ❌ Bad
   Foreground(lipgloss.Color("#5fb3b3"))

   // ✅ Good
   Foreground(styles.GetColorAccentTeal())
   ```

2. **Don't create colors outside styles package**
   ```go
   // ❌ Bad
   myColor := lipgloss.Color("#c7ccd1")

   // ✅ Good
   myColor := styles.GetColorTextPrimary()
   ```

3. **Don't use magic numbers for spacing**
   ```go
   // ❌ Bad
   Padding(1, 2).Margin(1, 0)

   // ✅ Good
   spacing := styles.GetSpacingConstants()
   Padding(spacing.PaddingVerticalBase, spacing.PaddingHorizontalBase)
   ```

4. **Don't copy styles unnecessarily**
   ```go
   // ❌ Bad
   myStyle := lipgloss.NewStyle().
       Foreground(styles.GetColorTextPrimary()).
       Background(styles.GetColorBackgroundCard())

   // ✅ Good (use CardBase)
   myStyle := styles.GetCardBase()
   ```

5. **Don't hardcode layout widths**
   ```go
   // ❌ Bad
   card := styles.GetCardBase().Width(100)

   // ✅ Good
   card := styles.ResponsiveCard(m.Width)
   ```

## Migration Guide

### From Inline Styles to Constants

**Before:**
```go
myStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#c7ccd1")).
    Background(lipgloss.Color("#2d3346")).
    Padding(1, 2).
    Margin(1, 0)
```

**After:**
```go
myStyle := styles.GetCardBase()
```

### From Manual Layout to Responsive

**Before:**
```go
width := m.Width
if width > 120 {
    width = 120
}
card := styles.GetCardBase().Width(width - 4)
```

**After:**
```go
card := styles.ResponsiveCard(m.Width)
```

## Common Patterns

### Form Field Layout

```go
label := styles.GetInputLabel().Render("Email")
input := styles.GetInputFocused().Render("user@example.com")
hint := styles.GetInputHint().Render("Enter a valid email")

field := lipgloss.JoinVertical(lipgloss.Left, label, input, hint)
```

### Card with Header and Content

```go
header := styles.GetCardHeader().Render("Title")
content := styles.GetCardContent().Render("Content here...")
footer := styles.GetCardFooter().Render("Footer")

card := styles.GetCardBase().Render(
    lipgloss.JoinVertical(lipgloss.Left, header, content, footer),
)
```

### Modal Dialog

```go
title := styles.GetModalTitle().Render("Confirm?")
message := styles.GetModalMessage().Render("Are you sure?")
buttons := styles.GetModalButtonContainer().Render("[ Yes ] [ No ]")

modal := styles.GetModalBase().Render(
    lipgloss.JoinVertical(lipgloss.Center, title, message, buttons),
)
```

### Status Message

```go
statusText := styles.GetSuccessText().Render("✓ Success")
hintText := styles.GetSuccessHint().Render("Operation completed")

message := styles.GetSuccessBox().Render(
    lipgloss.JoinVertical(lipgloss.Left, statusText, hintText),
)
```

## Troubleshooting

### Issue: Color doesn't match design

**Solution:** Check that you're using the correct getter function:
```go
// Verify the color name matches the design
color := styles.GetColorAccentTeal()  // vs GetColorAccentGreen()
```

### Issue: Spacing is inconsistent

**Solution:** Use spacing constants instead of magic numbers:
```go
// Use constants
spacing := styles.GetSpacingConstants()
Padding(spacing.PaddingVerticalBase, spacing.PaddingHorizontalBase)
```

### Issue: Layout breaks on small terminals

**Solution:** Use responsive layout functions:
```go
// Use responsive instead of fixed width
card := styles.ResponsiveCard(m.Width)
```

### Issue: Need custom style

**Solution:** Copy an existing style and modify:
```go
// Start with existing style
myStyle := styles.GetCardBase().Copy().
    Foreground(styles.GetColorError())
```

## References

- `internal/cli/styles/styles.go` - Style definitions
- `internal/cli/styles/constants_export.go` - Exported getters
- `internal/cli/components/` - Component examples
- Lipgloss documentation: https://github.com/charmbracelet/lipgloss

