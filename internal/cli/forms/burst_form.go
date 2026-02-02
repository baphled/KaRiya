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
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - burstformdata must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstEditorFormWithData(data *BurstFormData) *huh.Form {
	return NewBurstEditorFormWithDataAndHeight(data, 0)
}

// NewBurstEditorFormWithDataAndHeight creates a form for editing a burst with initial form data and height.
//
// Expected:
//   - burstformdata must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstEditorFormWithDataAndHeight(data *BurstFormData, height int) *huh.Form {
	return NewBurstEditorFormWithDataAndDimensions(data, 0, height)
}

// NewBurstEditorFormWithDataAndDimensions creates a form for editing a burst with initial form data and dimensions.
//
// Expected:
//   - burstformdata must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstEditorFormWithDataAndDimensions(data *BurstFormData, width, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	// Create fields group (scrollable)
	fieldsGroup := huh.NewGroup(
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
	)

	return NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
}

// ApplyBurstFormData applies the form data to a burst domain object.
//
// Expected:
//   - burst must be valid.
//   - burstformdata must be valid.
//
// Side effects:
//   - None.
func ApplyBurstFormData(burst *career.Burst, data *BurstFormData) {
	burst.Name = data.Name
	burst.Description = data.Description
}

// GetBurstFormData extracts form data from a burst domain object.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A fully initialized BurstFormData ready for use.
//
// Side effects:
//   - None.
func GetBurstFormData(burst *career.Burst) *BurstFormData {
	return &BurstFormData{
		Name:        burst.Name,
		Description: burst.Description,
	}
}
