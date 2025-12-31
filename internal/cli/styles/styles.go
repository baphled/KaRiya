package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color Scheme
// Professional dark theme with muted accents for focused, distraction-free interface
var (
	// Background colors
	ColorBackground     = lipgloss.Color("#1a1f2e") // Dark blue-gray background
	ColorBackgroundAlt  = lipgloss.Color("#242936") // Slightly lighter background for contrast
	ColorBackgroundCard = lipgloss.Color("#2d3346") // Card/panel background

	// Accent colors (muted and professional)
	ColorAccentTeal   = lipgloss.Color("#5fb3b3") // Muted teal for primary actions
	ColorAccentGreen  = lipgloss.Color("#6cb56c") // Muted green for success/confirmation
	ColorAccentPurple = lipgloss.Color("#a99bd1") // Muted purple for selected items

	// Text colors
	ColorTextPrimary   = lipgloss.Color("#c7ccd1") // Primary text (light gray)
	ColorTextSecondary = lipgloss.Color("#8b92a0") // Secondary text (medium gray)
	ColorTextMuted     = lipgloss.Color("#5e6673") // Muted text (dark gray)

	// Status colors
	ColorError   = lipgloss.Color("#d76e6e") // Error red (muted)
	ColorWarning = lipgloss.Color("#d9a66c") // Warning amber (muted)
	ColorSuccess = lipgloss.Color("#6cb56c") // Success green (same as accent)
	ColorInfo    = lipgloss.Color("#6ab0d3") // Info blue (muted)

	// Border colors
	ColorBorder       = lipgloss.Color("#3d4454") // Default border
	ColorBorderActive = lipgloss.Color("#5fb3b3") // Active/focused border (teal)
	ColorBorderError  = lipgloss.Color("#d76e6e") // Error border
)

// Base Styles
// Reusable style definitions for common UI elements

// Button styles
var (
	ButtonBase = lipgloss.NewStyle().
			Padding(0, 3).
			MarginRight(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	ButtonPrimary = ButtonBase.
			Foreground(ColorTextPrimary).
			BorderForeground(ColorAccentTeal).
			Background(ColorBackgroundCard)

	ButtonSecondary = ButtonBase.
			Foreground(ColorTextSecondary).
			BorderForeground(ColorBorder).
			Background(ColorBackground)

	ButtonFocused = ButtonPrimary.
			BorderForeground(ColorAccentTeal).
			Bold(true)

	ButtonDisabled = ButtonBase.
			Foreground(ColorTextMuted).
			BorderForeground(ColorBorder).
			Faint(true)
)

// Input field styles
var (
	InputBase = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorBackgroundCard)

	InputFocused = InputBase.Copy().
			BorderForeground(ColorBorderActive).
			BorderStyle(lipgloss.ThickBorder())

	InputError = InputBase.Copy().
			BorderForeground(ColorBorderError).
			BorderStyle(lipgloss.ThickBorder())

	InputLabel = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Bold(true).
			MarginBottom(1)

	InputHint = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Italic(true).
			MarginTop(1)
)

// Card styles
var (
	CardBase = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorBackgroundCard)

	CardHeader = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Bold(true).
			MarginBottom(1)

	CardContent = lipgloss.NewStyle().
			Foreground(ColorTextPrimary)

	CardFooter = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			MarginTop(1)
)

// Modal/Dialog styles
var (
	ModalBase = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorBackgroundCard)

	ModalTitle = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Bold(true).
			MarginBottom(1)

	ModalMessage = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			MarginBottom(2)

	ModalButtonContainer = lipgloss.NewStyle().
				MarginTop(2).
				MarginBottom(1)

	ModalInstructions = lipgloss.NewStyle().
				Foreground(ColorTextMuted).
				MarginTop(1)

	// Destructive modal styles (for delete confirmations)
	ModalDestructive = lipgloss.NewStyle().
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorError).
				Background(ColorBackgroundCard)

	ModalDestructiveTitle = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true).
				MarginBottom(1)
)

// Header styles
var (
	HeaderMain = lipgloss.NewStyle().
			Foreground(ColorAccentTeal).
			Bold(true).
			Padding(1, 0).
			MarginBottom(1)

	HeaderSection = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Bold(true).
			Padding(0, 0).
			MarginTop(1).
			MarginBottom(1)

	HeaderSubsection = lipgloss.NewStyle().
				Foreground(ColorTextSecondary).
				Bold(true).
				MarginBottom(1)
)

// Error message styles
var (
	ErrorBox = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderError).
			Background(ColorBackgroundCard).
			Foreground(ColorError)

	ErrorText = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	ErrorHint = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Italic(true).
			MarginTop(1)
)

// Warning message styles
var (
	WarningBox = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorWarning).
			Background(ColorBackgroundCard).
			Foreground(ColorWarning)

	WarningText = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	WarningHint = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Italic(true).
			MarginTop(1)
)

// Success message styles
var (
	SuccessBox = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Background(ColorBackgroundCard).
			Foreground(ColorSuccess)

	SuccessText = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	SuccessHint = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Italic(true).
			MarginTop(1)
)

// Info message styles
var (
	InfoBox = lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorInfo).
		Background(ColorBackgroundCard).
		Foreground(ColorInfo)

	InfoText = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Bold(true)

	InfoHint = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Italic(true).
			MarginTop(1)
)

// List styles
var (
	ListItem = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(ColorTextPrimary)

	ListItemSelected = lipgloss.NewStyle().
				Padding(0, 1).
				Foreground(ColorTextPrimary).
				Background(ColorAccentPurple).
				Bold(true)

	ListItemFocused = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(ColorTextPrimary).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(ColorAccentTeal)
)

// Tag styles
var (
	TagBase = lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Background(ColorBackgroundAlt).
		Foreground(ColorTextSecondary)

	TagSelected = TagBase.Copy().
			BorderForeground(ColorAccentPurple).
			Background(ColorAccentPurple).
			Foreground(ColorTextPrimary).
			Bold(true)
)

// Progress indicator styles
var (
	ProgressBar = lipgloss.NewStyle().
			Foreground(ColorAccentTeal).
			Background(ColorBackgroundAlt)

	ProgressText = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			MarginLeft(2)
)

// Spinner styles
var (
	SpinnerStyle = lipgloss.NewStyle().
		Foreground(ColorAccentTeal)
)
// Badge and error message styles
var (
	Badge = lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Background(ColorBackgroundAlt).
		Foreground(ColorTextSecondary)

	BadgeSelected = Badge.Copy().
		BorderForeground(ColorAccentPurple).
		Background(ColorAccentPurple).
		Foreground(ColorTextPrimary).
		Bold(true)

	BadgeFocused = Badge.Copy().
		BorderForeground(ColorAccentTeal).
		BorderStyle(lipgloss.ThickBorder()).
		Background(ColorBackgroundCard)

	ErrorMsg = lipgloss.NewStyle().
		Foreground(ColorError).
		MarginTop(1)

	ButtonPrimaryFocused = ButtonPrimary.Copy().
		Bold(true).
		BorderStyle(lipgloss.ThickBorder())

	ButtonSecondaryFocused = ButtonSecondary.Copy().
		Bold(true).
		BorderStyle(lipgloss.ThickBorder())

	Label = lipgloss.NewStyle().
		Foreground(ColorTextSecondary).
		Bold(true).
		MarginBottom(1)

	LabelFocused = Label.Copy().
		Foreground(ColorAccentTeal).
		Bold(true)

	Hint = lipgloss.NewStyle().
		Foreground(ColorTextMuted).
		Italic(true).
		MarginTop(1)

	Warning = lipgloss.NewStyle().
		Foreground(ColorWarning).
		Bold(true)
)


// Helper functions for common layout patterns

// WithBorder adds a border to a style with the default border color
func WithBorder(style lipgloss.Style) lipgloss.Style {
	return style.Copy().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)
}

// WithFocusedBorder adds a focused border to a style
func WithFocusedBorder(style lipgloss.Style) lipgloss.Style {
	return style.Copy().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ColorBorderActive)
}

// WithErrorBorder adds an error border to a style
func WithErrorBorder(style lipgloss.Style) lipgloss.Style {
	return style.Copy().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ColorBorderError)
}

// WithPadding adds consistent padding to a style
func WithPadding(style lipgloss.Style, vertical, horizontal int) lipgloss.Style {
	return style.Copy().Padding(vertical, horizontal)
}

// WithMargin adds consistent margin to a style
func WithMargin(style lipgloss.Style, vertical, horizontal int) lipgloss.Style {
	return style.Copy().Margin(vertical, horizontal)
}

// Responsive Layout Helpers

// MaxWidth returns the maximum width for content based on terminal width
// Ensures content doesn't exceed reasonable line lengths for readability
func MaxWidth(terminalWidth int) int {
	const maxContentWidth = 120
	if terminalWidth < maxContentWidth {
		return terminalWidth - 4 // Leave margin
	}
	return maxContentWidth
}

// CenterHorizontal centers text horizontally within a given width
func CenterHorizontal(text string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, text)
}

// CenterVertical centers text vertically within a given height
func CenterVertical(text string, height int) string {
	return lipgloss.PlaceVertical(height, lipgloss.Center, text)
}

// Center centers text both horizontally and vertically
func Center(text string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, text)
}

// AlignLeft aligns text to the left within a given width
func AlignLeft(text string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, text)
}

// AlignRight aligns text to the right within a given width
func AlignRight(text string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Right, text)
}

// ResponsiveStyle returns a style with width adjusted to terminal size
func ResponsiveStyle(terminalWidth int) lipgloss.Style {
	width := MaxWidth(terminalWidth)
	return lipgloss.NewStyle().Width(width)
}

// ResponsiveCard returns a card style with responsive width
func ResponsiveCard(terminalWidth int) lipgloss.Style {
	width := MaxWidth(terminalWidth)
	return CardBase.Copy().Width(width - 4) // Account for padding
}

// ResponsiveInput returns an input style with responsive width
func ResponsiveInput(terminalWidth int) lipgloss.Style {
	width := MaxWidth(terminalWidth)
	return InputBase.Copy().Width(width - 8) // Account for padding and borders
}

// TwoColumn splits content into two columns with responsive widths
func TwoColumn(terminalWidth int) (leftWidth, rightWidth int) {
	maxWidth := MaxWidth(terminalWidth)
	leftWidth = maxWidth / 2
	rightWidth = maxWidth - leftWidth
	return leftWidth, rightWidth
}

// ThreeColumn splits content into three columns with responsive widths
func ThreeColumn(terminalWidth int) (leftWidth, centerWidth, rightWidth int) {
	maxWidth := MaxWidth(terminalWidth)
	leftWidth = maxWidth / 3
	centerWidth = maxWidth / 3
	rightWidth = maxWidth - leftWidth - centerWidth
	return leftWidth, centerWidth, rightWidth
}

// Grid returns dimensions for a grid layout with the given number of columns
func Grid(terminalWidth, columns int) (columnWidth, gutter int) {
	maxWidth := MaxWidth(terminalWidth)
	gutter = 2
	totalGutter := gutter * (columns - 1)
	columnWidth = (maxWidth - totalGutter) / columns
	return columnWidth, gutter
}
