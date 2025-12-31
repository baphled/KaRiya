package components

import "fmt"

// AudienceRelevanceSelector manages selection of audience types for facts
type AudienceRelevanceSelector struct {
	selected map[string]bool
	options  []string
}

// NewAudienceRelevanceSelector creates a new audience relevance selector
func NewAudienceRelevanceSelector() *AudienceRelevanceSelector {
	return &AudienceRelevanceSelector{
		selected: make(map[string]bool),
		options:  []string{"hiring_manager", "recruiter", "peer"},
	}
}

// SetSelected sets the selected audience types
func (a *AudienceRelevanceSelector) SetSelected(audiences []string) {
	a.selected = make(map[string]bool)
	for _, aud := range audiences {
		a.selected[aud] = true
	}
}

// GetSelected returns the list of selected audience types
func (a *AudienceRelevanceSelector) GetSelected() []string {
	var result []string
	for _, opt := range a.options {
		if a.selected[opt] {
			result = append(result, opt)
		}
	}
	return result
}

// ToggleSelected toggles the current selected option
func (a *AudienceRelevanceSelector) ToggleSelected() {
	// This is a placeholder - implementation would depend on current focus
}

// IsSelected checks if an audience type is selected
func (a *AudienceRelevanceSelector) IsSelected(audience string) bool {
	return a.selected[audience]
}

// Render renders the selector
func (a *AudienceRelevanceSelector) Render() string {
	var result string
	for _, opt := range a.options {
		if a.selected[opt] {
			result += fmt.Sprintf("[✓ %s] ", opt)
		} else {
			result += fmt.Sprintf("[ %s ] ", opt)
		}
	}
	return result
}
