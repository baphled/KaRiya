package forms

import (
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/charmbracelet/huh"
)

// BurstSuggestionFormData holds the form data for editing a burst suggestion.
// This is used in the BurstSuggestionModel when the user presses 'e' to edit
// the name and description of a suggested burst before confirming it.
type BurstSuggestionFormData struct {
	Name        string
	Description string
}

// NewBurstSuggestionEditForm creates a form for editing a burst suggestion's name and description.
//
// Expected:
//   - burstsuggestion must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstSuggestionEditForm(suggestion burstfact.BurstSuggestion) *huh.Form {
	data := GetBurstSuggestionFormData(suggestion)

	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "name",
				Title:       "Burst Name",
				Description: "Name for this burst (optional, will be auto-generated if empty)",
				Placeholder: "Enter burst name...",
				CharLimit:   100,
			}).Value(&data.Name),

			NewInput(FieldConfig{
				Key:         "description",
				Title:       "Burst Description",
				Description: "Description for this burst (optional)",
				Placeholder: "Enter burst description...",
				CharLimit:   500,
			}).Value(&data.Description),
		),
	)
}

// NewBurstSuggestionEditFormWithData creates a form with pre-populated data.
//
// Expected:
//   - burstsuggestionformdata must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstSuggestionEditFormWithData(data *BurstSuggestionFormData) *huh.Form {
	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "name",
				Title:       "Burst Name",
				Description: "Name for this burst (optional, will be auto-generated if empty)",
				Placeholder: "Enter burst name...",
				CharLimit:   100,
			}).Value(&data.Name),

			NewInput(FieldConfig{
				Key:         "description",
				Title:       "Burst Description",
				Description: "Description for this burst (optional)",
				Placeholder: "Enter burst description...",
				CharLimit:   500,
			}).Value(&data.Description),
		),
	)
}

// ApplyBurstSuggestionFormData applies the form data to a burst suggestion.
//
// Expected:
//   - burstsuggestion must be valid.
//   - burstsuggestionformdata must be valid.
//
// Side effects:
//   - None.
func ApplyBurstSuggestionFormData(suggestion *burstfact.BurstSuggestion, data *BurstSuggestionFormData) {
	suggestion.Name = data.Name
	suggestion.Description = data.Description
}

// GetBurstSuggestionFormData extracts form data from a burst suggestion.
//
// Expected:
//   - burstsuggestion must be valid.
//
// Returns:
//   - A fully initialized BurstSuggestionFormData ready for use.
//
// Side effects:
//   - None.
func GetBurstSuggestionFormData(suggestion burstfact.BurstSuggestion) *BurstSuggestionFormData {
	return &BurstSuggestionFormData{
		Name:        suggestion.Name,
		Description: suggestion.Description,
	}
}
