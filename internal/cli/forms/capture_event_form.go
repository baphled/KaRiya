package forms

import (
	"sort"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// CaptureEventFormData holds the form data for capturing career events.
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
		Text:            "",
		Date:            "",
		Company:         "",
		Project:         "",
		Tags:            []string{},
		Categories:      []string{},
		SubmitConfirmed: false,
	}
}

// NewCaptureEventForm creates a huh-based form for capturing career events.
// strategy: "quick" (text+date only) or "manual" (all fields)
// width, height: terminal dimensions for proper sizing
func NewCaptureEventForm(data *CaptureEventFormData, strategy string, width, height int) *huh.Form {
	var fieldsGroup *huh.Group

	if strategy == "quick" {
		// Quick capture: essential fields (text, date, company)
		fieldsGroup = huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Event Description",
				Description: "Describe what you accomplished (required, 10-2000 characters)",
				Placeholder: "Enter event description...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text).Lines(5),

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
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Company),
		)
	} else {
		// Manual capture: all fields
		// Prepare options for tags and categories from domain
		tagList := make([]string, 0, len(career.AllowedTags))
		for tag := range career.AllowedTags {
			tagList = append(tagList, tag)
		}
		sort.Strings(tagList)

		tagOptions := make([]huh.Option[string], 0, len(tagList))
		for _, tag := range tagList {
			tagOptions = append(tagOptions, huh.NewOption(tag, tag))
		}

		catList := make([]string, 0, len(career.AllowedCategories))
		for cat := range career.AllowedCategories {
			catList = append(catList, cat)
		}
		sort.Strings(catList)

		categoryOptions := make([]huh.Option[string], 0, len(catList))
		for _, cat := range catList {
			categoryOptions = append(categoryOptions, huh.NewOption(cat, cat))
		}

		fieldsGroup = huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Event Description",
				Description: "Describe what you accomplished (required, 10-2000 characters)",
				Placeholder: "Enter event description...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text).Lines(5),

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
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Company),

			NewInput(FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project name (optional)",
				Placeholder: "Enter project name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Project),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("Select relevant tags (optional, use space to select)").
				Options(tagOptions...).
				Value(&data.Tags).
				Limit(6),

			huh.NewMultiSelect[string]().
				Key("categories").
				Title("Categories").
				Description("Select relevant categories (optional, use space to select)").
				Options(categoryOptions...).
				Value(&data.Categories).
				Limit(6),
		)
	}

	// Create form with fixed confirm button at bottom
	return NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
}

// NewCaptureEventFormForModal creates a huh-based form for modal overlays (WITHOUT confirm button).
// This creates a simple form where pressing Enter on the last field completes the form.
// Use this for quick modal interactions, NOT for full-screen forms.
//
// strategy: "quick" (text+date+company) or "manual" (all fields)
// width, height: terminal dimensions for proper sizing
//
// Pattern (like FilterModal):
//   - Just form fields, no Submit/Cancel button
//   - Press Enter on last field = save and close
//   - Press Esc = cancel and close
func NewCaptureEventFormForModal(data *CaptureEventFormData, strategy string, width, height int) *huh.Form {
	var fieldsGroup *huh.Group

	if strategy == "quick" {
		// Quick capture: essential fields (text, date, company)
		fieldsGroup = huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Event Description",
				Description: "Describe what you accomplished (required, 10-2000 characters)",
				Placeholder: "Enter event description...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text).Lines(5),

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
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Company),
		)
	} else {
		// Manual capture: all fields
		// Prepare options for tags and categories from domain
		tagList := make([]string, 0, len(career.AllowedTags))
		for tag := range career.AllowedTags {
			tagList = append(tagList, tag)
		}
		sort.Strings(tagList)

		tagOptions := make([]huh.Option[string], 0, len(tagList))
		for _, tag := range tagList {
			tagOptions = append(tagOptions, huh.NewOption(tag, tag))
		}

		catList := make([]string, 0, len(career.AllowedCategories))
		for cat := range career.AllowedCategories {
			catList = append(catList, cat)
		}
		sort.Strings(catList)

		categoryOptions := make([]huh.Option[string], 0, len(catList))
		for _, cat := range catList {
			categoryOptions = append(categoryOptions, huh.NewOption(cat, cat))
		}

		fieldsGroup = huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Event Description",
				Description: "Describe what you accomplished (required, 10-2000 characters)",
				Placeholder: "Enter event description...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text).Lines(5),

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
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Company),

			NewInput(FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project name (optional)",
				Placeholder: "Enter project name...",
				CharLimit:   200,
				Validate:    MaxLength(200),
			}).Value(&data.Project),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("Select relevant tags (optional, use space to select)").
				Options(tagOptions...).
				Value(&data.Tags).
				Limit(6),

			huh.NewMultiSelect[string]().
				Key("categories").
				Title("Categories").
				Description("Select relevant categories (optional, use space to select)").
				Options(categoryOptions...).
				Value(&data.Categories).
				Limit(6),
		)
	}

	// Create simple form WITHOUT confirm button (modal pattern)
	// Just like FilterModal - press Enter on last field to complete
	return huh.NewForm(fieldsGroup).
		WithWidth(width).
		WithHeight(height).
		WithTheme(Theme())
}
