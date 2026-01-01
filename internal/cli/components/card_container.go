package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

// CardContainer renders a boxed content container with optional header, body, and footer sections.
// It provides a consistent card-like appearance with borders, padding, and background styling.
// CardContainer is a stateless rendering component.
type CardContainer struct {
	header          string
	body            string
	footer          string
	backgroundColor lipgloss.Color
	borderColor     lipgloss.Color
	hasHeader       bool
	hasBody         bool
	hasFooter       bool
}

// NewCardContainer creates a new CardContainer with default styling.
func NewCardContainer() *CardContainer {
	return &CardContainer{
		backgroundColor: styles.ColorBackgroundCard,
		borderColor:     styles.ColorBorder,
		hasHeader:       false,
		hasBody:         false,
		hasFooter:       false,
	}
}

// SetHeader sets the header content for the card.
// This method uses the builder pattern to allow method chaining.
func (cc *CardContainer) SetHeader(header string) *CardContainer {
	cc.header = header
	cc.hasHeader = true
	return cc
}

// SetBody sets the body content for the card.
// This method uses the builder pattern to allow method chaining.
func (cc *CardContainer) SetBody(body string) *CardContainer {
	cc.body = body
	cc.hasBody = true
	return cc
}

// SetFooter sets the footer content for the card.
// This method uses the builder pattern to allow method chaining.
func (cc *CardContainer) SetFooter(footer string) *CardContainer {
	cc.footer = footer
	cc.hasFooter = true
	return cc
}

// WithBackgroundColor sets a custom background color for the card.
// This method uses the builder pattern to allow method chaining.
func (cc *CardContainer) WithBackgroundColor(color lipgloss.Color) *CardContainer {
	cc.backgroundColor = color
	return cc
}

// WithBorderColor sets a custom border color for the card.
// This method uses the builder pattern to allow method chaining.
func (cc *CardContainer) WithBorderColor(color lipgloss.Color) *CardContainer {
	cc.borderColor = color
	return cc
}

// Render returns the styled card container as a string.
// It combines header, body, and footer sections with appropriate styling and spacing.
func (cc *CardContainer) Render() string {
	var sections []string

	// Render header if present
	if cc.hasHeader {
		headerStyle := styles.CardHeader.Copy().
			Foreground(styles.ColorTextPrimary)
		sections = append(sections, headerStyle.Render(cc.header))
	}

	// Render body if present
	if cc.hasBody {
		bodyStyle := styles.CardContent.Copy().
			Foreground(styles.ColorTextPrimary)
		sections = append(sections, bodyStyle.Render(cc.body))
	}

	// Render footer if present
	if cc.hasFooter {
		footerStyle := styles.CardFooter.Copy().
			Foreground(styles.ColorTextSecondary)
		sections = append(sections, footerStyle.Render(cc.footer))
	}

	// Combine all sections
	content := strings.Join(sections, "\n")

	// Apply card styling with border and background
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(cc.borderColor).
		Background(cc.backgroundColor).
		Foreground(styles.ColorTextPrimary)

	return cardStyle.Render(content)
}
