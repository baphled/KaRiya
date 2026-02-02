package selectors

import (
	"fmt"
	"strings"
)

// AudienceRelevanceSelector manages selection of audience types for facts.
type AudienceRelevanceSelector struct {
	selected map[string]bool
	options  []string
}

// NewAudienceRelevanceSelector creates a new audience relevance selector.
//
// Returns:
//   - A fully initialized AudienceRelevanceSelector ready for use.
//
// Side effects:
//   - None.
func NewAudienceRelevanceSelector() *AudienceRelevanceSelector {
	return &AudienceRelevanceSelector{
		selected: make(map[string]bool),
		options:  []string{"hiring_manager", "recruiter", "peer"},
	}
}

// SetSelected sets the selected audience types.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (a *AudienceRelevanceSelector) SetSelected(audiences []string) {
	a.selected = make(map[string]bool)
	for _, aud := range audiences {
		a.selected[aud] = true
	}
}

// GetSelected returns the list of selected audience types.
//
// Returns:
//   - A []string value.
//
// Side effects:
//   - None.
func (a *AudienceRelevanceSelector) GetSelected() []string {
	var result []string
	for _, opt := range a.options {
		if a.selected[opt] {
			result = append(result, opt)
		}
	}
	return result
}

// IsSelected checks if an audience type is selected.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (a *AudienceRelevanceSelector) IsSelected(audience string) bool {
	return a.selected[audience]
}

// Render renders the selector.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (a *AudienceRelevanceSelector) Render() string {
	var b strings.Builder
	for _, opt := range a.options {
		if a.selected[opt] {
			b.WriteString(fmt.Sprintf("[✓ %s] ", opt))
		} else {
			b.WriteString(fmt.Sprintf("[ %s ] ", opt))
		}
	}
	return b.String()
}
