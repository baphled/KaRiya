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

// NewBurstForm creates a burst edit form bound to the given data.
//
// Expected:
//   - data must be a valid BurstFormData pointer.
//   - width and height control form dimensions (0 means unconstrained).
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewBurstForm(data *BurstFormData, width, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	fields := []huh.Field{
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
	}

	return newScrollableForm(fields, &data.SubmitConfirmed, width, height)
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
