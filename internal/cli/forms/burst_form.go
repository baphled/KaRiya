package forms

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// BurstFormData holds the form data for burst editing.
type BurstFormData struct {
	Name            string
	Description     string
	SubmitConfirmed bool
}

// NewBurstEditorForm creates a form for editing a burst.
// The form has two fields: Name (required) and Description (optional),
// plus a Submit confirmation button.
func NewBurstEditorForm(burst *career.Burst) *huh.Form {
	data := &BurstFormData{
		Name:            burst.Name,
		Description:     burst.Description,
		SubmitConfirmed: false,
	}

	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "name",
				Title:       "Burst Name",
				Description: "A descriptive name for this burst (required)",
				Placeholder: "Enter burst name...",
				CharLimit:   200,
				Validate:    Title,
			}).Value(&data.Name),

			NewText(FieldConfig{
				Key:         "description",
				Title:       "Description",
				Description: "Optional description of the burst",
				Placeholder: "Enter description...",
				CharLimit:   1000,
				Validate:    Description,
			}).Value(&data.Description),

			huh.NewConfirm().
				Key("submit").
				Title("Save Changes").
				Description("Submit the form to save your changes").
				Affirmative("Submit").
				Negative("Cancel").
				Value(&data.SubmitConfirmed),
		),
	)
}

// NewBurstEditorFormWithData creates a form for editing a burst with initial form data.
// This variant allows external data binding for more control.
func NewBurstEditorFormWithData(data *BurstFormData) *huh.Form {
	return NewBurstEditorFormWithDataAndHeight(data, 0)
}

// NewBurstEditorFormWithDataAndHeight creates a form for editing a burst with initial form data and height.
// When height > 0, the form becomes scrollable if content exceeds the height.
func NewBurstEditorFormWithDataAndHeight(data *BurstFormData, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	group := huh.NewGroup(
		NewInput(FieldConfig{
			Key:         "name",
			Title:       "Burst Name",
			Description: "A descriptive name for this burst (required)",
			Placeholder: "Enter burst name...",
			CharLimit:   200,
			Validate:    Title,
		}).Value(&data.Name),

		NewText(FieldConfig{
			Key:         "description",
			Title:       "Description",
			Description: "Optional description of the burst",
			Placeholder: "Enter description...",
			CharLimit:   1000,
			Validate:    Description,
		}).Value(&data.Description),

		huh.NewConfirm().
			Key("submit").
			Title("Save Changes").
			Description("Submit the form to save your changes").
			Affirmative("Submit").
			Negative("Cancel").
			Value(&data.SubmitConfirmed),
	)

	if height > 0 {
		return NewFormWithHeight(height, group)
	}
	return NewForm(group)
}

// ApplyBurstFormData applies the form data to a burst domain object.
func ApplyBurstFormData(burst *career.Burst, data *BurstFormData) {
	burst.Name = data.Name
	burst.Description = data.Description
}

// GetBurstFormData extracts form data from a burst domain object.
func GetBurstFormData(burst *career.Burst) *BurstFormData {
	return &BurstFormData{
		Name:        burst.Name,
		Description: burst.Description,
	}
}
