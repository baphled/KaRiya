package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// Color Palette Export
// ============================================================================
// All color constants are centralized and exported for use throughout the
// application. This ensures consistency and makes it easy to update colors
// globally.

// GetColorBackground returns the primary background color
func GetColorBackground() lipgloss.Color {
	return ColorBackground
}

// GetColorBackgroundAlt returns the alternate background color for contrast
func GetColorBackgroundAlt() lipgloss.Color {
	return ColorBackgroundAlt
}

// GetColorBackgroundCard returns the card/panel background color
func GetColorBackgroundCard() lipgloss.Color {
	return ColorBackgroundCard
}

// GetColorAccentTeal returns the primary accent color (teal)
func GetColorAccentTeal() lipgloss.Color {
	return ColorAccentTeal
}

// GetColorAccentGreen returns the success accent color (green)
func GetColorAccentGreen() lipgloss.Color {
	return ColorAccentGreen
}

// GetColorAccentPurple returns the selection/highlight color (purple)
func GetColorAccentPurple() lipgloss.Color {
	return ColorAccentPurple
}

// GetColorTextPrimary returns the primary text color
func GetColorTextPrimary() lipgloss.Color {
	return ColorTextPrimary
}

// GetColorTextSecondary returns the secondary text color
func GetColorTextSecondary() lipgloss.Color {
	return ColorTextSecondary
}

// GetColorTextMuted returns the muted text color
func GetColorTextMuted() lipgloss.Color {
	return ColorTextMuted
}

// GetColorError returns the error/failure color
func GetColorError() lipgloss.Color {
	return ColorError
}

// GetColorWarning returns the warning color
func GetColorWarning() lipgloss.Color {
	return ColorWarning
}

// GetColorSuccess returns the success color
func GetColorSuccess() lipgloss.Color {
	return ColorSuccess
}

// GetColorInfo returns the info color
func GetColorInfo() lipgloss.Color {
	return ColorInfo
}

// GetColorBorder returns the default border color
func GetColorBorder() lipgloss.Color {
	return ColorBorder
}

// GetColorBorderActive returns the active/focused border color
func GetColorBorderActive() lipgloss.Color {
	return ColorBorderActive
}

// GetColorBorderError returns the error border color
func GetColorBorderError() lipgloss.Color {
	return ColorBorderError
}

// ============================================================================
// Button Styles Export
// ============================================================================

// GetButtonBase returns the base button style
func GetButtonBase() lipgloss.Style {
	return ButtonBase
}

// GetButtonPrimary returns the primary button style
func GetButtonPrimary() lipgloss.Style {
	return ButtonPrimary
}

// GetButtonSecondary returns the secondary button style
func GetButtonSecondary() lipgloss.Style {
	return ButtonSecondary
}

// GetButtonFocused returns the focused button style
func GetButtonFocused() lipgloss.Style {
	return ButtonFocused
}

// GetButtonDisabled returns the disabled button style
func GetButtonDisabled() lipgloss.Style {
	return ButtonDisabled
}

// GetButtonPrimaryFocused returns the focused primary button style
func GetButtonPrimaryFocused() lipgloss.Style {
	return ButtonPrimaryFocused
}

// GetButtonSecondaryFocused returns the focused secondary button style
func GetButtonSecondaryFocused() lipgloss.Style {
	return ButtonSecondaryFocused
}

// ============================================================================
// Input Field Styles Export
// ============================================================================

// GetInputBase returns the base input field style
func GetInputBase() lipgloss.Style {
	return InputBase
}

// GetInputFocused returns the focused input field style
func GetInputFocused() lipgloss.Style {
	return InputFocused
}

// GetInputError returns the error input field style
func GetInputError() lipgloss.Style {
	return InputError
}

// GetInputLabel returns the input label style
func GetInputLabel() lipgloss.Style {
	return InputLabel
}

// GetInputHint returns the input hint style
func GetInputHint() lipgloss.Style {
	return InputHint
}

// ============================================================================
// Card Styles Export
// ============================================================================

// GetCardBase returns the base card style
func GetCardBase() lipgloss.Style {
	return CardBase
}

// GetCardHeader returns the card header style
func GetCardHeader() lipgloss.Style {
	return CardHeader
}

// GetCardContent returns the card content style
func GetCardContent() lipgloss.Style {
	return CardContent
}

// GetCardFooter returns the card footer style
func GetCardFooter() lipgloss.Style {
	return CardFooter
}

// ============================================================================
// Modal/Dialog Styles Export
// ============================================================================

// GetModalBase returns the base modal style
func GetModalBase() lipgloss.Style {
	return ModalBase
}

// GetModalTitle returns the modal title style
func GetModalTitle() lipgloss.Style {
	return ModalTitle
}

// GetModalMessage returns the modal message style
func GetModalMessage() lipgloss.Style {
	return ModalMessage
}

// GetModalButtonContainer returns the modal button container style
func GetModalButtonContainer() lipgloss.Style {
	return ModalButtonContainer
}

// GetModalInstructions returns the modal instructions style
func GetModalInstructions() lipgloss.Style {
	return ModalInstructions
}

// GetModalDestructive returns the destructive modal style
func GetModalDestructive() lipgloss.Style {
	return ModalDestructive
}

// GetModalDestructiveTitle returns the destructive modal title style
func GetModalDestructiveTitle() lipgloss.Style {
	return ModalDestructiveTitle
}

// ============================================================================
// Header Styles Export
// ============================================================================

// GetHeaderMain returns the main header style
func GetHeaderMain() lipgloss.Style {
	return HeaderMain
}

// GetHeaderSection returns the section header style
func GetHeaderSection() lipgloss.Style {
	return HeaderSection
}

// GetHeaderSubsection returns the subsection header style
func GetHeaderSubsection() lipgloss.Style {
	return HeaderSubsection
}

// ============================================================================
// Error Message Styles Export
// ============================================================================

// GetErrorBox returns the error box style
func GetErrorBox() lipgloss.Style {
	return ErrorBox
}

// GetErrorText returns the error text style
func GetErrorText() lipgloss.Style {
	return ErrorText
}

// GetErrorHint returns the error hint style
func GetErrorHint() lipgloss.Style {
	return ErrorHint
}

// GetErrorMsg returns the error message style
func GetErrorMsg() lipgloss.Style {
	return ErrorMsg
}

// ============================================================================
// Warning Message Styles Export
// ============================================================================

// GetWarningBox returns the warning box style
func GetWarningBox() lipgloss.Style {
	return WarningBox
}

// GetWarningText returns the warning text style
func GetWarningText() lipgloss.Style {
	return WarningText
}

// GetWarningHint returns the warning hint style
func GetWarningHint() lipgloss.Style {
	return WarningHint
}

// GetWarning returns the warning style
func GetWarning() lipgloss.Style {
	return Warning
}

// ============================================================================
// Success Message Styles Export
// ============================================================================

// GetSuccessBox returns the success box style
func GetSuccessBox() lipgloss.Style {
	return SuccessBox
}

// GetSuccessText returns the success text style
func GetSuccessText() lipgloss.Style {
	return SuccessText
}

// GetSuccessHint returns the success hint style
func GetSuccessHint() lipgloss.Style {
	return SuccessHint
}

// ============================================================================
// Info Message Styles Export
// ============================================================================

// GetInfoBox returns the info box style
func GetInfoBox() lipgloss.Style {
	return InfoBox
}

// GetInfoText returns the info text style
func GetInfoText() lipgloss.Style {
	return InfoText
}

// GetInfoHint returns the info hint style
func GetInfoHint() lipgloss.Style {
	return InfoHint
}

// ============================================================================
// List Styles Export
// ============================================================================

// GetListItem returns the list item style
func GetListItem() lipgloss.Style {
	return ListItem
}

// GetListItemSelected returns the selected list item style
func GetListItemSelected() lipgloss.Style {
	return ListItemSelected
}

// GetListItemFocused returns the focused list item style
func GetListItemFocused() lipgloss.Style {
	return ListItemFocused
}

// ============================================================================
// Tag Styles Export
// ============================================================================

// GetTagBase returns the base tag style
func GetTagBase() lipgloss.Style {
	return TagBase
}

// GetTagSelected returns the selected tag style
func GetTagSelected() lipgloss.Style {
	return TagSelected
}

// ============================================================================
// Progress Indicator Styles Export
// ============================================================================

// GetProgressBar returns the progress bar style
func GetProgressBar() lipgloss.Style {
	return ProgressBar
}

// GetProgressText returns the progress text style
func GetProgressText() lipgloss.Style {
	return ProgressText
}

// ============================================================================
// Spinner Styles Export
// ============================================================================

// GetSpinnerStyle returns the spinner style
func GetSpinnerStyle() lipgloss.Style {
	return SpinnerStyle
}

// ============================================================================
// Badge Styles Export
// ============================================================================

// GetBadge returns the base badge style
func GetBadge() lipgloss.Style {
	return Badge
}

// GetBadgeSelected returns the selected badge style
func GetBadgeSelected() lipgloss.Style {
	return BadgeSelected
}

// GetBadgeFocused returns the focused badge style
func GetBadgeFocused() lipgloss.Style {
	return BadgeFocused
}

// ============================================================================
// Label and Hint Styles Export
// ============================================================================

// GetLabel returns the label style
func GetLabel() lipgloss.Style {
	return Label
}

// GetLabelFocused returns the focused label style
func GetLabelFocused() lipgloss.Style {
	return LabelFocused
}

// GetHint returns the hint style
func GetHint() lipgloss.Style {
	return Hint
}

// ============================================================================
// Spacing Constants Export
// ============================================================================

// SpacingConstants holds all spacing values used throughout the application
type SpacingConstants struct {
	// Horizontal padding for buttons and form fields
	PaddingHorizontalSmall int
	PaddingHorizontalBase  int
	PaddingHorizontalLarge int

	// Vertical padding for cards and containers
	PaddingVerticalSmall int
	PaddingVerticalBase  int

	// Margins for spacing between elements
	MarginSmall int
	MarginBase  int
	MarginLarge int

	// Maximum content width for readability
	MaxContentWidth int

	// Grid gutter width
	GridGutterWidth int
}

// GetSpacingConstants returns the standard spacing constants
func GetSpacingConstants() SpacingConstants {
	return SpacingConstants{
		PaddingHorizontalSmall: 1,
		PaddingHorizontalBase:  2,
		PaddingHorizontalLarge: 3,
		PaddingVerticalSmall:   0,
		PaddingVerticalBase:    1,
		MarginSmall:            1,
		MarginBase:             2,
		MarginLarge:            3,
		MaxContentWidth:        120,
		GridGutterWidth:        2,
	}
}

// ============================================================================
// Color Palette Documentation
// ============================================================================
// Background Colors:
//   - ColorBackground: Primary dark background (#1a1f2e)
//   - ColorBackgroundAlt: Alternate background for contrast (#242936)
//   - ColorBackgroundCard: Card/panel background (#2d3346)
//
// Accent Colors (Muted & Professional):
//   - ColorAccentTeal: Primary action color (#5fb3b3)
//   - ColorAccentGreen: Success/confirmation color (#6cb56c)
//   - ColorAccentPurple: Selection/highlight color (#a99bd1)
//
// Text Colors:
//   - ColorTextPrimary: Main text (#c7ccd1)
//   - ColorTextSecondary: Secondary text (#8b92a0)
//   - ColorTextMuted: Muted text (#5e6673)
//
// Status Colors:
//   - ColorError: Error messages (#d76e6e)
//   - ColorWarning: Warning messages (#d9a66c)
//   - ColorSuccess: Success messages (#6cb56c)
//   - ColorInfo: Info messages (#6ab0d3)
//
// Border Colors:
//   - ColorBorder: Default border (#3d4454)
//   - ColorBorderActive: Active/focused border (#5fb3b3)
//   - ColorBorderError: Error border (#d76e6e)

