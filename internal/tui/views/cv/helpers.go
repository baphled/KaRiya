// Package cv provides shared helpers for CV views.
package cv

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/config"
)

// NarrativeProfile contains presentation-ready profile details for CV rendering.
type NarrativeProfile struct {
	Name     string
	Role     string
	Location string
	Email    string
	GitHub   string
}

// NarrativeProfileFromConfig builds presentation-ready profile details from config.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A NarrativeProfile value.
//
// Side effects:
//   - None.
func NarrativeProfileFromConfig(cfg *config.ProfileConfig) NarrativeProfile {
	if cfg == nil {
		return NarrativeProfile{}
	}

	return NarrativeProfile{
		Name:     buildNarrativeProfileName(cfg.FirstName, cfg.LastName, cfg.Name),
		Role:     cfg.Title,
		Location: cfg.Location,
		Email:    cfg.Email,
		GitHub:   cfg.GitHub,
	}
}

func buildNarrativeProfileName(firstName, lastName, legacyName string) string {
	if firstName != "" {
		if lastName != "" {
			return strings.TrimSpace(firstName + " " + lastName)
		}
		return firstName
	}

	return legacyName
}

// WordWrap wraps text to a specified width at word boundaries.
//
// Expected:
//   - text is the string to wrap.
//   - width is the maximum line width (must be positive).
//
// Returns:
//   - A string with newlines inserted at word boundaries.
//
// Side effects:
//   - None.
func WordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	paragraphs := strings.Split(text, "\n")
	var wrappedParagraphs []string
	for _, para := range paragraphs {
		if para == "" {
			wrappedParagraphs = append(wrappedParagraphs, "")
			continue
		}
		wrappedParagraphs = append(wrappedParagraphs, WrapParagraph(para, width))
	}
	return strings.Join(wrappedParagraphs, "\n")
}

// WrapParagraph wraps a paragraph of text, preserving existing line breaks.
//
// Expected:
//   - text is the paragraph to wrap.
//   - width is the maximum line width (must be positive).
//
// Returns:
//   - A string with each line wrapped to the specified width.
//
// Side effects:
//   - None.
func WrapParagraph(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	var currentLine strings.Builder
	currentLen := 0

	words := strings.Fields(text)
	for i, word := range words {
		wordLen := len(word)
		if currentLen > 0 && currentLen+1+wordLen > width {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
			currentLine.Reset()
			currentLen = 0
		}
		if currentLen > 0 {
			currentLine.WriteString(" ")
			currentLen++
		}
		currentLine.WriteString(word)
		currentLen += wordLen
		if wordLen > width && i < len(words)-1 {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
			currentLine.Reset()
			currentLen = 0
		}
	}
	if currentLine.Len() > 0 {
		result.WriteString(currentLine.String())
	}
	return result.String()
}

// MinInt returns the smaller of two integers.
//
// Expected:
//   - a and b are integers to compare.
//
// Returns:
//   - The smaller of the two values.
//
// Side effects:
//   - None.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// handleViewportWindowResize processes a window size change for viewport-based views.
// It updates terminal info, stores dimensions, and marks the viewport as needing
// reinitialisation.
func handleViewportWindowResize(view interface {
	SetTerminalInfo(width, height int)
}, msg tea.WindowSizeMsg, width *int, height *int, ready *bool) {
	view.SetTerminalInfo(msg.Width, msg.Height)
	*width = msg.Width
	*height = msg.Height
	*ready = false
}
