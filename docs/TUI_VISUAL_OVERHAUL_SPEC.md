# TUI Visual Overhaul Specification

**Version**: 1.0  
**Created**: 2026-01-09  
**Status**: Draft  
**Related Task**: Task 38

---

## Table of Contents

1. [Vision & Goals](#vision--goals)
2. [Theme System Architecture](#theme-system-architecture)
3. [Color Palette Structure](#color-palette-structure)
4. [Style Set Structure](#style-set-structure)
5. [Component Specifications](#component-specifications)
6. [Animation Specifications](#animation-specifications)
7. [Responsive Design](#responsive-design)
8. [Terminal Detection](#terminal-detection)
9. [Integration Guide](#integration-guide)
10. [Acceptance Criteria](#acceptance-criteria)

---

## Vision & Goals

### Inspiration

The visual overhaul draws inspiration from **btop** - a modern, polished terminal resource monitor known for:

- Clean, professional aesthetic
- Excellent use of terminal colors
- Smooth animations and transitions
- Responsive layouts that adapt to terminal size
- Consistent visual language throughout

### Goals

1. **Professional Polish**: Elevate KaRiya's visual presentation to match best-in-class TUI applications
2. **Theme System**: Enable runtime theme switching with multiple built-in themes
3. **Consistency**: Eliminate hardcoded styles in favor of semantic, theme-aware styling
4. **Accessibility**: Support dark/light modes and various terminal color depths
5. **Performance**: Maintain sub-50ms render times for all views
6. **Backwards Compatibility**: Default theme matches current colors for seamless migration

### Non-Goals

- Custom user-defined themes (can be added later)
- Per-component theme overrides
- CSS-like stylesheets

---

## Theme System Architecture

### Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      ThemeManager                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Theme A   │  │   Theme B   │  │   Theme C   │  ...     │
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
// Theme defines the contract for all themes
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
    Primary() lipgloss.Color
    Secondary() lipgloss.Color
    Accent() lipgloss.Color
    Background() lipgloss.Color
    Foreground() lipgloss.Color
    Muted() lipgloss.Color
    Success() lipgloss.Color
    Warning() lipgloss.Color
    Error() lipgloss.Color
    Info() lipgloss.Color
    Border() lipgloss.Color
    BorderActive() lipgloss.Color
}
```

### ThemeManager

```go
// ThemeManager handles theme registration and switching
type ThemeManager struct {
    themes      map[string]Theme
    active      Theme
    defaultName string
    onChange    []func(Theme)
}

// Key Methods
func NewThemeManager() *ThemeManager
func (tm *ThemeManager) Register(theme Theme) error
func (tm *ThemeManager) SetActive(name string) error
func (tm *ThemeManager) Active() Theme
func (tm *ThemeManager) List() []string
func (tm *ThemeManager) OnChange(callback func(Theme))
```

### Usage in GlobalContext

```go
// GlobalContext addition
type GlobalContext struct {
    // ... existing fields
    ThemeManager *themes.ThemeManager
}

// Usage in components
theme := ctx.ThemeManager.Active()
style := lipgloss.NewStyle().
    Foreground(theme.Foreground()).
    Background(theme.Background())
```

---

## Color Palette Structure

### ColorPalette Struct

```go
type ColorPalette struct {
    // Background colors
    Background     lipgloss.Color // Primary background
    BackgroundAlt  lipgloss.Color // Alternate/contrast background
    BackgroundCard lipgloss.Color // Card/panel background
    
    // Foreground colors
    Foreground      lipgloss.Color // Primary text
    ForegroundDim   lipgloss.Color // Secondary/dimmed text
    ForegroundMuted lipgloss.Color // Muted/disabled text
    
    // Accent colors
    Primary   lipgloss.Color // Primary action/brand color
    Secondary lipgloss.Color // Secondary actions
    Tertiary  lipgloss.Color // Tertiary/highlight
    
    // Status colors
    Success lipgloss.Color
    Warning lipgloss.Color
    Error   lipgloss.Color
    Info    lipgloss.Color
    
    // Border colors
    Border       lipgloss.Color // Default borders
    BorderActive lipgloss.Color // Focused/active borders
    BorderError  lipgloss.Color // Error state borders
    
    // Special colors
    Selection   lipgloss.Color // Selected item background
    Highlight   lipgloss.Color // Highlighted text background
    Link        lipgloss.Color // Hyperlink color
}
```

### Example: Default Theme Palette

Based on current KaRiya colors (backwards compatible):

```go
var DefaultPalette = ColorPalette{
    // Backgrounds
    Background:     lipgloss.Color("#1a1f2e"),
    BackgroundAlt:  lipgloss.Color("#242936"),
    BackgroundCard: lipgloss.Color("#2d3346"),
    
    // Foregrounds
    Foreground:      lipgloss.Color("#c7ccd1"),
    ForegroundDim:   lipgloss.Color("#8b92a0"),
    ForegroundMuted: lipgloss.Color("#5e6673"),
    
    // Accents
    Primary:   lipgloss.Color("#5fb3b3"), // Teal
    Secondary: lipgloss.Color("#6cb56c"), // Green
    Tertiary:  lipgloss.Color("#a99bd1"), // Purple
    
    // Status
    Success: lipgloss.Color("#6cb56c"),
    Warning: lipgloss.Color("#d9a66c"),
    Error:   lipgloss.Color("#d76e6e"),
    Info:    lipgloss.Color("#6ab0d3"),
    
    // Borders
    Border:       lipgloss.Color("#3d4454"),
    BorderActive: lipgloss.Color("#5fb3b3"),
    BorderError:  lipgloss.Color("#d76e6e"),
    
    // Special
    Selection: lipgloss.Color("#3d4454"),
    Highlight: lipgloss.Color("#4d5566"),
    Link:      lipgloss.Color("#6ab0d3"),
}
```

---

## Style Set Structure

### StyleSet Struct

Pre-composed Lipgloss styles generated from the theme's palette:

```go
type StyleSet struct {
    // Buttons
    ButtonBase            lipgloss.Style
    ButtonPrimary         lipgloss.Style
    ButtonSecondary       lipgloss.Style
    ButtonFocused         lipgloss.Style
    ButtonDisabled        lipgloss.Style
    
    // Inputs
    InputBase    lipgloss.Style
    InputFocused lipgloss.Style
    InputError   lipgloss.Style
    InputLabel   lipgloss.Style
    InputHint    lipgloss.Style
    
    // Cards/Panels
    CardBase    lipgloss.Style
    CardHeader  lipgloss.Style
    CardContent lipgloss.Style
    CardFooter  lipgloss.Style
    
    // Lists
    ListItem         lipgloss.Style
    ListItemSelected lipgloss.Style
    ListItemFocused  lipgloss.Style
    
    // Messages
    ErrorBox   lipgloss.Style
    WarningBox lipgloss.Style
    SuccessBox lipgloss.Style
    InfoBox    lipgloss.Style
    
    // Headers
    HeaderMain       lipgloss.Style
    HeaderSection    lipgloss.Style
    HeaderSubsection lipgloss.Style
    
    // Progress
    ProgressBar  lipgloss.Style
    ProgressText lipgloss.Style
    
    // Badges/Tags
    Badge         lipgloss.Style
    BadgeSelected lipgloss.Style
    Tag           lipgloss.Style
    TagSelected   lipgloss.Style
    
    // Key Badges (new)
    KeyBadge     lipgloss.Style
    KeyBadgeHint lipgloss.Style
}
```

### Style Generation

Styles are generated from the palette using a factory function:

```go
func GenerateStyles(palette *ColorPalette) *StyleSet {
    return &StyleSet{
        ButtonBase: lipgloss.NewStyle().
            Padding(0, 2).
            Border(lipgloss.RoundedBorder()).
            BorderForeground(palette.Border),
            
        ButtonPrimary: lipgloss.NewStyle().
            Padding(0, 2).
            Background(palette.Primary).
            Foreground(palette.Background).
            Border(lipgloss.RoundedBorder()).
            BorderForeground(palette.Primary),
            
        // ... etc
    }
}
```

---

## Component Specifications

### 1. SelectionList

Enhanced list component with visual highlighting for selected items.

**Features:**
- Full-width background highlighting on selected item
- Thick left border indicator on focused item
- Arrow indicator (▶) before selected item text
- Smooth keyboard navigation (j/k, up/down)

**Visual Example:**
```
┌─ Select Profile ─────────────────────────────┐
│                                              │
│   Senior Individual Contributor              │
│ ▶ Staff Engineer                    ◀────────│← Selected (highlighted bg)
│   Engineering Manager                        │
│   Technical Lead                             │
│                                              │
│ ↑/↓ or j/k to move • Enter to select        │
└──────────────────────────────────────────────┘
```

**Implementation:**
```go
type SelectionList struct {
    items       []string
    selected    int
    focused     bool
    theme       themes.Theme
    width       int
    height      int
}

func (s *SelectionList) View() string {
    var rows []string
    for i, item := range s.items {
        style := s.theme.Styles().ListItem
        if i == s.selected {
            style = s.theme.Styles().ListItemSelected.
                Width(s.width)
            item = "▶ " + item
        }
        rows = append(rows, style.Render(item))
    }
    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
```

### 2. KeyBadge

Styled keyboard shortcut badges for help footers.

**Visual Example:**
```
 [↑/↓] Navigate   [Enter] Select   [Esc] Cancel   [?] Help
```

**Styling:**
- Key in contrasting background (e.g., primary on card background)
- Rounded corners
- Consistent padding
- Hint text in muted color

**Implementation:**
```go
type KeyBadge struct {
    key   string
    hint  string
    theme themes.Theme
}

func (k *KeyBadge) View() string {
    keyStyle := k.theme.Styles().KeyBadge
    hintStyle := k.theme.Styles().KeyBadgeHint
    
    return keyStyle.Render(k.key) + " " + hintStyle.Render(k.hint)
}

// Helper for building footers
func BuildHelpFooter(theme themes.Theme, badges ...KeyBadge) string {
    var parts []string
    for _, badge := range badges {
        parts = append(parts, badge.View())
    }
    return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
```

### 3. GradientProgress

Progress bar with color gradient from start to end.

**Visual Example:**
```
Generating CV... ████████████░░░░░░░░ 60%
                 ↑ cyan → blue gradient
```

**Implementation:**
```go
type GradientProgress struct {
    percent    float64
    width      int
    startColor lipgloss.Color
    endColor   lipgloss.Color
    theme      themes.Theme
}

func (g *GradientProgress) View() string {
    filled := int(g.percent * float64(g.width))
    empty := g.width - filled
    
    // Interpolate colors for gradient effect
    bar := ""
    for i := 0; i < filled; i++ {
        color := interpolateColor(g.startColor, g.endColor, float64(i)/float64(g.width))
        bar += lipgloss.NewStyle().Foreground(color).Render("█")
    }
    bar += lipgloss.NewStyle().Foreground(g.theme.Muted()).Render(strings.Repeat("░", empty))
    
    return bar + fmt.Sprintf(" %.0f%%", g.percent*100)
}
```

### 4. ResponsiveCard

Card component that adapts to terminal width.

**Layout Modes:**

| Mode | Width | Characteristics |
|------|-------|-----------------|
| Compact | < 60 | No borders, minimal padding |
| Standard | 60-100 | Rounded borders, normal padding |
| Wide | > 100 | Double padding, max-width constraint |

**Implementation:**
```go
type ResponsiveCard struct {
    title   string
    content string
    theme   themes.Theme
    width   int
}

func (r *ResponsiveCard) View() string {
    style := r.theme.Styles().CardBase
    
    if r.width < 60 {
        // Compact: no border, minimal padding
        style = style.UnsetBorder().Padding(0, 1)
    } else if r.width > 100 {
        // Wide: constrain width, extra padding
        style = style.MaxWidth(100).Padding(1, 3)
    }
    
    header := r.theme.Styles().CardHeader.Render(r.title)
    body := r.theme.Styles().CardContent.Render(r.content)
    
    return style.Render(lipgloss.JoinVertical(lipgloss.Left, header, body))
}
```

### 5. EnhancedTable

Table with alternating row colors and hover effects.

**Visual Example:**
```
┌─────────────────────────────────────────────┐
│ Date       │ Company      │ Event          │
├─────────────────────────────────────────────┤
│ 2026-01-01 │ TechCorp     │ Led project... │← Row 1 (alt bg)
│ 2025-12-15 │ StartupXYZ   │ Implemented... │← Row 2 (normal bg)
│ 2025-11-30 │ BigCompany   │ Designed...    │← Row 3 (alt bg, selected)
└─────────────────────────────────────────────┘
```

**Implementation:**
```go
type EnhancedTable struct {
    headers  []string
    rows     [][]string
    selected int
    theme    themes.Theme
}

func (t *EnhancedTable) View() string {
    // Render headers
    headerStyle := t.theme.Styles().CardHeader
    
    // Render rows with alternating backgrounds
    var renderedRows []string
    for i, row := range t.rows {
        rowStyle := t.theme.Styles().ListItem
        if i%2 == 0 {
            rowStyle = rowStyle.Background(t.theme.Palette().BackgroundAlt)
        }
        if i == t.selected {
            rowStyle = t.theme.Styles().ListItemSelected
        }
        renderedRows = append(renderedRows, renderRow(row, rowStyle))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}
```

---

## Animation Specifications

### Dependencies

- **Harmonica**: Spring-based animations (`github.com/charmbracelet/harmonica`)
- **Glamour**: Markdown rendering (`github.com/charmbracelet/glamour`)

### 1. Logo Fade-In

Subtle fade-in animation for the KaRiya logo on startup.

```go
type LogoAnimation struct {
    spring   harmonica.Spring
    opacity  float64
    done     bool
}

func (l *LogoAnimation) Update(msg tea.Msg) tea.Cmd {
    switch msg.(type) {
    case AnimationTickMsg:
        l.opacity, _, l.done = l.spring.Update(l.opacity, 1.0)
        if !l.done {
            return TickAnimation()
        }
    }
    return nil
}

func (l *LogoAnimation) View() string {
    // Apply opacity by adjusting color brightness
    return renderLogoWithOpacity(l.opacity)
}
```

### 2. Loading Spinner

Enhanced spinner with rotating messages.

```go
type LoadingSpinner struct {
    spinner  spinner.Model
    messages []string
    current  int
    rotator  *components.LoadingMessageRotator
}

// Messages rotate every 2 seconds
var loadingMessages = []string{
    "Processing...",
    "Almost there...",
    "Working on it...",
    "Hang tight...",
}
```

### 3. Smooth List Scrolling

Spring-based smooth scrolling for long lists.

```go
type SmoothScroller struct {
    spring       harmonica.Spring
    targetOffset float64
    currentOffset float64
}

func (s *SmoothScroller) ScrollTo(index int) {
    s.targetOffset = float64(index * itemHeight)
}

func (s *SmoothScroller) Update() bool {
    s.currentOffset, _, done := s.spring.Update(s.currentOffset, s.targetOffset)
    return !done
}
```

---

## Responsive Design

### Breakpoints

| Category | Width (cols) | Height (rows) | Characteristics |
|----------|-------------|---------------|-----------------|
| Tiny | < 60 | < 20 | Minimal chrome, no borders |
| Compact | 60-79 | 20-29 | Reduced margins, compact elements |
| Normal | 80-119 | 30-39 | Standard layout, full features |
| Large | 120-159 | 40-49 | Spacious, multi-column possible |
| XLarge | 160+ | 50+ | Maximum features, wide layouts |

### Adaptation Strategies

**Margins:**
```go
func (m *Manager) GetMargins() terminal.Margins {
    switch m.terminalInfo.GetCategory() {
    case terminal.SizeTiny:
        return terminal.Margins{Top: 0, Right: 1, Bottom: 0, Left: 1}
    case terminal.SizeCompact:
        return terminal.Margins{Top: 1, Right: 2, Bottom: 1, Left: 2}
    case terminal.SizeLarge, terminal.SizeXLarge:
        return terminal.Margins{Top: 2, Right: 6, Bottom: 2, Left: 6}
    default:
        return terminal.Margins{Top: 2, Right: 4, Bottom: 2, Left: 4}
    }
}
```

**Content Decisions:**
- **Tiny**: Hide logo, minimal help text, no decorative borders
- **Compact**: Abbreviated help text, single-column layouts
- **Normal**: Full features, standard layouts
- **Large+**: Multi-column layouts, expanded details

---

## Terminal Detection

### Color Depth Detection

```go
type ColorDepth int

const (
    ColorDepth16    ColorDepth = 16    // Basic 16 colors
    ColorDepth256   ColorDepth = 256   // Extended palette
    ColorDepthTrue  ColorDepth = 16777216 // 24-bit true color
)

func DetectColorDepth() ColorDepth {
    // Check COLORTERM environment variable
    if colorterm := os.Getenv("COLORTERM"); colorterm != "" {
        if colorterm == "truecolor" || colorterm == "24bit" {
            return ColorDepthTrue
        }
    }
    
    // Check TERM environment variable
    term := os.Getenv("TERM")
    if strings.Contains(term, "256color") {
        return ColorDepth256
    }
    
    return ColorDepth16
}
```

### Dark/Light Mode Detection

```go
func DetectDarkMode() bool {
    // Check common environment variables
    if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
        parts := strings.Split(colorfgbg, ";")
        if len(parts) >= 2 {
            // Light background typically has higher value
            if bg, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
                return bg < 8 // Dark backgrounds are typically 0-7
            }
        }
    }
    
    // Default to dark mode (most terminal users prefer dark)
    return true
}
```

### Auto Theme Selection

```go
func (tm *ThemeManager) AutoSelect() {
    depth := DetectColorDepth()
    isDark := DetectDarkMode()
    
    // Select appropriate theme based on terminal capabilities
    if depth < ColorDepth256 {
        // Use high-contrast theme for limited colors
        tm.SetActive("high-contrast")
    } else if isDark {
        tm.SetActive("default") // or "btop-dark"
    } else {
        tm.SetActive("light") // or "gruvbox-light"
    }
}
```

---

## Integration Guide

### Updating Existing Components

**Before (hardcoded):**
```go
func (i *SomeIntent) viewContent() string {
    titleStyle := lipgloss.NewStyle().
        Foreground(styles.ColorAccentTeal).
        Bold(true)
    
    cardStyle := lipgloss.NewStyle().
        Background(styles.ColorBackgroundCard).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(styles.ColorBorder)
    
    return cardStyle.Render(titleStyle.Render("Title"))
}
```

**After (theme-aware):**
```go
func (i *SomeIntent) viewContent() string {
    theme := i.ctx.ThemeManager.Active()
    
    title := theme.Styles().HeaderSection.Render("Title")
    
    return theme.Styles().CardBase.Render(title)
}
```

### Creating New Theme

```go
// internal/cli/themes/mytheme.go
package themes

func NewMyTheme() Theme {
    palette := &ColorPalette{
        Background:     lipgloss.Color("#1e1e2e"),
        BackgroundAlt:  lipgloss.Color("#313244"),
        // ... define all colors
    }
    
    return &BaseTheme{
        name:        "my-theme",
        description: "My custom theme",
        author:      "Your Name",
        isDark:      true,
        palette:     palette,
        styles:      GenerateStyles(palette),
    }
}

// Register in manager
func init() {
    DefaultManager.Register(NewMyTheme())
}
```

### Huh Forms Integration

```go
// Generate Huh theme from active theme
func GenerateHuhTheme(theme themes.Theme) *huh.Theme {
    p := theme.Palette()
    
    return &huh.Theme{
        Focused: huh.FieldStyles{
            Base: lipgloss.NewStyle().
                Border(lipgloss.RoundedBorder()).
                BorderForeground(p.BorderActive),
            Title: lipgloss.NewStyle().
                Foreground(p.Primary),
            // ... etc
        },
        Blurred: huh.FieldStyles{
            Base: lipgloss.NewStyle().
                Border(lipgloss.RoundedBorder()).
                BorderForeground(p.Border),
            // ... etc
        },
    }
}
```

---

## Acceptance Criteria

### Functional Requirements

- [ ] Theme system fully functional with runtime switching
- [ ] Terminal color depth auto-detection working
- [ ] Dark/light mode auto-detection working
- [ ] All 5 primary intents using theme system
- [ ] Huh forms unified with main theme
- [ ] All hardcoded `styles.Color*` references replaced

### Visual Requirements

- [ ] Consistent styling across all screens
- [ ] Selection lists have background highlighting
- [ ] Key badges used in all footers
- [ ] Responsive layouts at all breakpoints
- [ ] No visual regressions from current appearance

### Performance Requirements

- [ ] View render time < 50ms
- [ ] Theme switching < 10ms
- [ ] No memory leaks from theme changes

### Quality Requirements

- [ ] All existing tests pass (2000+)
- [ ] Code coverage maintained ≥ 80%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions
- [ ] Theme system tests > 90% coverage

### Documentation Requirements

- [ ] Theme Customization Guide created
- [ ] TUI Developer Guide updated
- [ ] TUI Standards updated
- [ ] All new components documented

---

## Appendix: File Structure

```
internal/cli/
├── themes/
│   ├── theme.go           # Theme interface and ColorPalette
│   ├── theme_test.go
│   ├── manager.go         # ThemeManager
│   ├── manager_test.go
│   ├── detector.go        # Terminal detection
│   ├── detector_test.go
│   ├── styles.go          # StyleSet and GenerateStyles
│   ├── styles_test.go
│   └── default.go         # Default theme (current colors)
├── components/
│   ├── selection_list.go
│   ├── selection_list_test.go
│   ├── key_badge.go
│   ├── key_badge_test.go
│   ├── gradient_progress.go
│   ├── gradient_progress_test.go
│   ├── responsive_card.go
│   ├── responsive_card_test.go
│   ├── enhanced_table.go
│   └── enhanced_table_test.go
└── context/
    └── global.go          # Add ThemeManager field
```

---

## References

- [btop](https://github.com/aristocratos/btop) - Visual inspiration
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library
- [Harmonica](https://github.com/charmbracelet/harmonica) - Animation library
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Catppuccin](https://github.com/catppuccin/catppuccin) - Theme reference
- [Nord](https://www.nordtheme.com/) - Theme reference
