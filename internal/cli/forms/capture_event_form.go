package forms

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// CaptureEventFormData holds the form field values for capturing a career event.
type CaptureEventFormData struct {
	Text            string
	Date            string
	Company         string
	Project         string
	Tags            []string
	Categories      []string
	SubmitConfirmed bool
}

// NewCaptureEventFormData creates a new CaptureEventFormData with default values.
func NewCaptureEventFormData() *CaptureEventFormData {
	return &CaptureEventFormData{
		Tags:       []string{},
		Categories: []string{},
	}
}

// TagOptions returns huh options for all allowed tags.
func TagOptions() []huh.Option[string] {
	tags := []struct {
		key   string
		label string
	}{
		{"project", "Project"},
		{"achievement", "Achievement"},
		{"leadership", "Leadership"},
		{"technical", "Technical"},
		{"consulting", "Consulting"},
		{"research", "Research"},
		{"product", "Product"},
		{"mentoring", "Mentoring"},
	}
	options := make([]huh.Option[string], len(tags))
	for i, t := range tags {
		options[i] = huh.NewOption(t.label, t.key)
	}
	return options
}

// CategoryOptions returns huh options for all allowed categories.
func CategoryOptions() []huh.Option[string] {
	categories := []struct {
		key   string
		label string
	}{
		{"technical", "Technical"},
		{"leadership", "Leadership"},
		{"product", "Product"},
		{"consulting", "Consulting"},
		{"research", "Research"},
		{"mentoring", "Mentoring"},
	}
	options := make([]huh.Option[string], len(categories))
	for i, c := range categories {
		options[i] = huh.NewOption(c.label, c.key)
	}
	return options
}

// NewCaptureEventForm creates a capture event form.
// strategy: "quick" shows only text field, "manual" shows all fields.
// width, height: dimensions for the form (0 for auto).
func NewCaptureEventForm(data *CaptureEventFormData, strategy string, width, height int) *huh.Form {
	isQuickMode := strategy == "quick"

	// Build fields based on mode
	var fields []huh.Field

	// Text field - always visible
	fields = append(fields,
		NewText(FieldConfig{
			Key:         "text",
			Title:       "Event Description",
			Description: "Describe what you accomplished (required, 10-2000 characters)",
			Placeholder: "Enter event description...",
			CharLimit:   2000,
			Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
		}).Value(&data.Text).Lines(5),
	)

	// Optional fields - only in manual mode
	if !isQuickMode {
		fields = append(fields,
			NewInput(FieldConfig{
				Key:         "date",
				Title:       "Date",
				Description: "YYYY-MM-DD, 'today', or relative like '1 week ago'",
				Placeholder: "Defaults to today",
				Validate:    DateFormat,
			}).Value(&data.Date),

			NewInput(FieldConfig{
				Key:         "company",
				Title:       "Company",
				Description: "Company where this happened",
				Placeholder: "Enter company name (optional)",
				CharLimit:   100,
				Validate:    CompanyName,
			}).Value(&data.Company),

			NewInput(FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project this relates to",
				Placeholder: "Enter project name (optional)",
				CharLimit:   100,
			}).Value(&data.Project),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("Select relevant tags (optional)").
				Options(TagOptions()...).
				Value(&data.Tags).
				Limit(len(career.AllowedTags)),

			huh.NewMultiSelect[string]().
				Key("categories").
				Title("Categories").
				Description("Select competency categories (optional)").
				Options(CategoryOptions()...).
				Value(&data.Categories).
				Limit(len(career.AllowedCategories)),
		)
	}

	// Create scrollable fields group
	fieldsGroup := huh.NewGroup(fields...)

	// Use NewFormWithFixedConfirm for scrollable fields + fixed confirm button
	return NewFormWithFixedConfirmCustom(fieldsGroup, &data.SubmitConfirmed, width, height, "Submit Event", "Submit this career event for review")
}

// NewFormWithFixedConfirmCustom creates a form with scrollable fields and a custom confirm button.
// This allows customizing the confirm button title and description.
// The form is created with proper height for scrolling - the fields group scrolls while
// the confirm button stays fixed at the bottom.
func NewFormWithFixedConfirmCustom(fieldsGroup *huh.Group, confirmValue *bool, width, height int, confirmTitle, confirmDescription string) *huh.Form {
	// Calculate height for fields group (reserve space for confirm)
	fieldsHeight := height - ConfirmButtonHeight
	if fieldsHeight < 5 {
		fieldsHeight = 5
	}

	// Create confirm group (fixed at bottom)
	confirmGroup := huh.NewGroup(
		huh.NewConfirm().
			Key("submit").
			Title(confirmTitle).
			Description(confirmDescription).
			Affirmative("Submit").
			Negative("Cancel").
			Value(confirmValue),
	)

	// Apply height to fields group to make it scrollable
	fieldsGroup = fieldsGroup.WithHeight(fieldsHeight)

	form := huh.NewForm(fieldsGroup, confirmGroup).
		WithTheme(Theme()).
		WithLayout(huh.LayoutStack).
		WithHeight(height) // Constrain form output to fit in available space

	if width > 0 {
		form = form.WithWidth(width)
	}

	return form
}

// NewCaptureEventFormWithDimensions creates a capture event form with specified dimensions.
func NewCaptureEventFormWithDimensions(data *CaptureEventFormData, strategy string, width, height int) *huh.Form {
	return NewCaptureEventForm(data, strategy, width, height)
}

// ApplyCaptureEventFormData extracts a CareerEvent from form data.
// Note: ID, CreatedAt, and UpdatedAt should be set by the caller.
func ApplyCaptureEventFormData(data *CaptureEventFormData) *career.CareerEvent {
	return &career.CareerEvent{
		Text:       data.Text,
		Company:    data.Company,
		Project:    data.Project,
		Tags:       data.Tags,
		Categories: data.Categories,
	}
}
