package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/baphled/kariya/internal/cli/styles"
)

// SectionContainer groups content with a title and consistent spacing.
// It provides visual separation and organization of related content areas.
// SectionContainer is a stateless rendering component.
type SectionContainer struct {
	title       string
	content     string
	hasTitle    bool
	spacing     int
}

// NewSectionContainer creates a new SectionContainer with default spacing.
func NewSectionContainer(content string) *SectionContainer {
	return &SectionContainer{
		content:  content,
		hasTitle: false,
		spacing:  1, // Default spacing between title and content
	}
}

// SetTitle sets the title for the section.
// This method uses the builder pattern to allow method chaining.
func (sc *SectionContainer) SetTitle(title string) *SectionContainer {
	sc.title = title
	sc.hasTitle = true
	return sc
}

// WithSpacing sets the spacing (in lines) between title and content.
// This method uses the builder pattern to allow method chaining.
func (sc *SectionContainer) WithSpacing(spacing int) *SectionContainer {
	if spacing < 0 {
		spacing = 0
	}
	sc.spacing = spacing
	return sc
}

// Render returns the styled section container as a string.
// It combines the title and content with appropriate styling and spacing.
func (sc *SectionContainer) Render() string {
	var parts []string

	// Render title if present
	if sc.hasTitle {
		titleStyle := styles.HeaderSection.Copy().
			Foreground(styles.ColorTextPrimary)
		parts = append(parts, titleStyle.Render(sc.title))

		// Add spacing between title and content
		if sc.spacing > 0 {
			parts = append(parts, strings.Repeat("\n", sc.spacing))
		}
	}

	// Render content
	contentStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary)
	parts = append(parts, contentStyle.Render(sc.content))

	return strings.Join(parts, "")
}

