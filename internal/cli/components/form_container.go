package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents a single field in the form with its label and content
type FormField struct {
	Label      string
	Input      string
	Error      string
	Hint       string
	IsFocused  bool
	IsRequired bool
	Width      int  // Optional: specify preferred width, 0 means auto
	FullWidth  bool // If true, field spans full width regardless of layout
}

// FormLayout defines how fields should be arranged
type FormLayout int

const (
	// SingleColumn layout: all fields stack vertically
	SingleColumn FormLayout = iota
	// TwoColumn layout: fields arranged in two columns when space allows
	TwoColumn
	// Responsive layout: automatically choose based on available width
	Responsive
)

// FormContainer is a smart form layout component that adapts to screen size and orientation.
// It intelligently arranges form fields based on available terminal width and height,
// supporting single-column, two-column, and responsive layouts.
type FormContainer struct {
	fields          []FormField
	width           int        // Available width in characters
	height          int        // Available height in characters
	layout          FormLayout // Layout strategy
	padding         int        // Horizontal padding
	verticalSpacing int        // Vertical spacing between fields
	columnGap       int        // Gap between columns in multi-column layouts
	minFieldWidth   int        // Minimum width for a field before wrapping
	maxFieldWidth   int        // Maximum width for a field
}

// NewFormContainer creates a new FormContainer with sensible defaults
func NewFormContainer() *FormContainer {
	return &FormContainer{
		fields:          []FormField{},
		width:           80, // Default terminal width
		height:          24, // Default terminal height
		layout:          Responsive,
		padding:         2,
		verticalSpacing: 2,
		columnGap:       4,
		minFieldWidth:   20,
		maxFieldWidth:   0, // 0 means no max
	}
}

// SetWidth sets the available width for the form container
func (fc *FormContainer) SetWidth(width int) *FormContainer {
	fc.width = width
	return fc
}

// SetHeight sets the available height for the form container
func (fc *FormContainer) SetHeight(height int) *FormContainer {
	fc.height = height
	return fc
}

// SetLayout sets the layout strategy for the form
func (fc *FormContainer) SetLayout(layout FormLayout) *FormContainer {
	fc.layout = layout
	return fc
}

// SetPadding sets the horizontal padding around the form
func (fc *FormContainer) SetPadding(padding int) *FormContainer {
	fc.padding = padding
	return fc
}

// SetVerticalSpacing sets the vertical spacing between fields
func (fc *FormContainer) SetVerticalSpacing(spacing int) *FormContainer {
	fc.verticalSpacing = spacing
	return fc
}

// SetColumnGap sets the gap between columns in multi-column layouts
func (fc *FormContainer) SetColumnGap(gap int) *FormContainer {
	fc.columnGap = gap
	return fc
}

// SetMinFieldWidth sets the minimum width for a field
func (fc *FormContainer) SetMinFieldWidth(width int) *FormContainer {
	fc.minFieldWidth = width
	return fc
}

// SetMaxFieldWidth sets the maximum width for a field (0 = no limit)
func (fc *FormContainer) SetMaxFieldWidth(width int) *FormContainer {
	fc.maxFieldWidth = width
	return fc
}

// AddField adds a form field to the container
func (fc *FormContainer) AddField(field FormField) *FormContainer {
	fc.fields = append(fc.fields, field)
	return fc
}

// AddFields adds multiple form fields to the container
func (fc *FormContainer) AddFields(fields ...FormField) *FormContainer {
	fc.fields = append(fc.fields, fields...)
	return fc
}

// Render returns the rendered form as a string
func (fc *FormContainer) Render() string {
	if len(fc.fields) == 0 {
		return ""
	}

	// Determine layout based on available width
	layout := fc.layout
	if layout == Responsive {
		layout = fc.determineOptimalLayout()
	}

	// Render based on chosen layout
	switch layout {
	case SingleColumn:
		return fc.renderSingleColumn()
	case TwoColumn:
		return fc.renderTwoColumn()
	default:
		return fc.renderSingleColumn()
	}
}

// determineOptimalLayout chooses the best layout based on available width
func (fc *FormContainer) determineOptimalLayout() FormLayout {
	availableWidth := fc.width - (fc.padding * 2)

	// Very narrow terminals: single column only
	if availableWidth < 60 {
		return SingleColumn
	}

	// Wide terminals: try two-column if we have enough space
	// Each column needs at least minFieldWidth + columnGap/2
	columnWidth := (availableWidth - fc.columnGap) / 2
	if columnWidth >= fc.minFieldWidth {
		return TwoColumn
	}

	return SingleColumn
}

// renderSingleColumn renders all fields in a single column
func (fc *FormContainer) renderSingleColumn() string {
	availableWidth := fc.width - (fc.padding * 2)
	fieldWidth := fc.calculateFieldWidth(availableWidth)

	var rendered []string
	for _, field := range fc.fields {
		rendered = append(rendered, fc.renderField(field, fieldWidth))
	}

	content := strings.Join(rendered, strings.Repeat("\n", fc.verticalSpacing))
	return fc.applyFormPadding(content)
}

// renderTwoColumn renders fields in two columns where possible
func (fc *FormContainer) renderTwoColumn() string {
	availableWidth := fc.width - (fc.padding * 2)
	columnWidth := (availableWidth - fc.columnGap) / 2

	var lines []string

	// Separate fields into full-width and regular fields
	var fullWidthFields []FormField
	var regularFields []FormField

	for _, field := range fc.fields {
		if field.FullWidth {
			fullWidthFields = append(fullWidthFields, field)
		} else {
			regularFields = append(regularFields, field)
		}
	}

	// Process regular fields in pairs for two-column layout
	for i := 0; i < len(regularFields); i += 2 {
		if i+1 < len(regularFields) {
			// We have a pair - render side by side
			leftRendered := fc.renderField(regularFields[i], columnWidth)
			rightRendered := fc.renderField(regularFields[i+1], columnWidth)

			// Split by lines and combine horizontally
			leftLines := strings.Split(leftRendered, "\n")
			rightLines := strings.Split(rightRendered, "\n")

			// Pad to same height
			maxLines := len(leftLines)
			if len(rightLines) > maxLines {
				maxLines = len(rightLines)
			}

			for j := 0; j < maxLines; j++ {
				leftLine := ""
				rightLine := ""

				if j < len(leftLines) {
					leftLine = leftLines[j]
				}
				if j < len(rightLines) {
					rightLine = rightLines[j]
				}

				// Pad left line to column width
				leftLine = padRight(leftLine, columnWidth)

				// Combine with gap
				combined := leftLine + strings.Repeat(" ", fc.columnGap) + rightLine
				lines = append(lines, combined)
			}

			// Add vertical spacing between field pairs
			if i+2 < len(regularFields) {
				for k := 0; k < fc.verticalSpacing; k++ {
					lines = append(lines, "")
				}
			}
		} else {
			// Odd field out - render full width
			rendered := fc.renderField(regularFields[i], availableWidth)
			renderedLines := strings.Split(rendered, "\n")
			lines = append(lines, renderedLines...)

			if i+1 < len(regularFields) {
				for k := 0; k < fc.verticalSpacing; k++ {
					lines = append(lines, "")
				}
			}
		}
	}

	// Add full-width fields at the end
	for i, field := range fullWidthFields {
		if len(lines) > 0 && i == 0 {
			for k := 0; k < fc.verticalSpacing; k++ {
				lines = append(lines, "")
			}
		}

		rendered := fc.renderField(field, availableWidth)
		renderedLines := strings.Split(rendered, "\n")
		lines = append(lines, renderedLines...)

		if i+1 < len(fullWidthFields) {
			for k := 0; k < fc.verticalSpacing; k++ {
				lines = append(lines, "")
			}
		}
	}

	content := strings.Join(lines, "\n")
	return fc.applyFormPadding(content)
}

// renderField renders a single form field with proper styling and wrapping
func (fc *FormContainer) renderField(field FormField, maxWidth int) string {
	var parts []string

	// Render label
	if field.Label != "" {
		labelStyle := styles.InputLabel.Copy().
			Foreground(styles.ColorTextSecondary)
		labelText := field.Label
		// Add focus indicator for focused fields
		if field.IsFocused {
			labelText = "► " + labelText
		}
		if field.IsRequired {
			labelText += " *"
		}
		parts = append(parts, labelStyle.Render(labelText))
	}

	// Render input
	if field.Input != "" {
		var inputStyle lipgloss.Style
		if field.Error != "" {
			inputStyle = styles.InputError.Copy().
				Foreground(styles.ColorTextPrimary)
		} else if field.IsFocused {
			inputStyle = styles.InputFocused.Copy().
				Foreground(styles.ColorTextPrimary)
		} else {
			inputStyle = styles.InputBase.Copy().
				Foreground(styles.ColorTextPrimary)
		}

		inputText := field.Input
		// Wrap input text if needed
		if maxWidth > 0 && len(inputText) > maxWidth && !strings.Contains(inputText, "\n") {
			inputText = fc.wrapText(inputText, maxWidth)
		}

		parts = append(parts, inputStyle.Render(inputText))
	}

	// Render error
	if field.Error != "" {
		errorStyle := styles.ErrorText.Copy().
			Foreground(styles.ColorError).
			MarginTop(1)
		errorText := field.Error
		if maxWidth > 0 && len(errorText) > maxWidth {
			errorText = fc.wrapText(errorText, maxWidth)
		}
		parts = append(parts, errorStyle.Render(errorText))
	}

	// Render hint
	if field.Hint != "" {
		hintStyle := styles.InputHint.Copy().
			Foreground(styles.ColorTextMuted)
		hintText := field.Hint
		if maxWidth > 0 && len(hintText) > maxWidth {
			hintText = fc.wrapText(hintText, maxWidth)
		}
		parts = append(parts, hintStyle.Render(hintText))
	}

	return strings.Join(parts, "\n")
}

// calculateFieldWidth calculates the optimal width for a field
func (fc *FormContainer) calculateFieldWidth(availableWidth int) int {
	if availableWidth < fc.minFieldWidth {
		return fc.minFieldWidth
	}

	width := availableWidth
	if fc.maxFieldWidth > 0 && width > fc.maxFieldWidth {
		width = fc.maxFieldWidth
	}

	return width
}

// applyFormPadding applies padding around the form content
func (fc *FormContainer) applyFormPadding(content string) string {
	if fc.padding == 0 {
		return content
	}

	padding := strings.Repeat(" ", fc.padding)
	lines := strings.Split(content, "\n")

	var paddedLines []string
	for _, line := range lines {
		paddedLines = append(paddedLines, padding+line)
	}

	return strings.Join(paddedLines, "\n")
}

// wrapText wraps text to fit within maxWidth characters
func (fc *FormContainer) wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= maxWidth {
			currentLine += " " + word
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// padRight pads a string to the right to reach targetWidth
func padRight(s string, targetWidth int) string {
	// Account for ANSI escape sequences in the string
	visibleLength := lipgloss.Width(s)
	if visibleLength >= targetWidth {
		return s
	}
	return s + strings.Repeat(" ", targetWidth-visibleLength)
}

// IsWideLayout returns true if the current layout is two-column
func (fc *FormContainer) IsWideLayout() bool {
	layout := fc.layout
	if layout == Responsive {
		layout = fc.determineOptimalLayout()
	}
	return layout == TwoColumn
}

// GetOptimalLayout returns the currently determined layout
func (fc *FormContainer) GetOptimalLayout() FormLayout {
	layout := fc.layout
	if layout == Responsive {
		layout = fc.determineOptimalLayout()
	}
	return layout
}

// GetCalculatedFieldWidth returns the width that would be used for fields in the current layout
func (fc *FormContainer) GetCalculatedFieldWidth() int {
	availableWidth := fc.width - (fc.padding * 2)
	layout := fc.GetOptimalLayout()

	if layout == TwoColumn {
		return (availableWidth - fc.columnGap) / 2
	}

	return fc.calculateFieldWidth(availableWidth)
}

// DebugInfo returns debug information about the form layout (useful for testing)
func (fc *FormContainer) DebugInfo() string {
	layout := fc.GetOptimalLayout()
	layoutName := "SingleColumn"
	if layout == TwoColumn {
		layoutName = "TwoColumn"
	}

	return fmt.Sprintf(
		"FormContainer Debug: width=%d, height=%d, layout=%s, fieldCount=%d, fieldWidth=%d",
		fc.width,
		fc.height,
		layoutName,
		len(fc.fields),
		fc.GetCalculatedFieldWidth(),
	)
}
