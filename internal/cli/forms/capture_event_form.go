package forms

import (
	"sort"

	"github.com/baphled/kariya/internal/constants"
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

// ConfirmSubmit marks the form data as submitted.
// Implements base.QuickSubmittable for Ctrl+S quick submit support.
//
// Side effects:
//   - Sets SubmitConfirmed to true.
func (d *CaptureEventFormData) ConfirmSubmit() {
	d.SubmitConfirmed = true
}

// NewCaptureEventFormData creates a new CaptureEventFormData with default values.
//
// Returns:
//   - A fully initialized CaptureEventFormData ready for use.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - captureeventformdata must be valid.
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewCaptureEventForm(data *CaptureEventFormData, strategy string, width, height int) *huh.Form {
	var fields []huh.Field

	if strategy == "quick" {
		fields = []huh.Field{
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
		}
	} else {
		allTags := constants.AllEventTags()
		tagList := make([]string, 0, len(allTags))
		for _, tag := range allTags {
			tagList = append(tagList, string(tag))
		}
		sort.Strings(tagList)

		tagOptions := make([]huh.Option[string], 0, len(tagList))
		for _, tag := range tagList {
			tagOptions = append(tagOptions, huh.NewOption(tag, tag))
		}

		allCats := constants.AllCompetencyCategories()
		catList := make([]string, 0, len(allCats))
		for _, cat := range allCats {
			catList = append(catList, string(cat))
		}
		sort.Strings(catList)

		categoryOptions := make([]huh.Option[string], 0, len(catList))
		for _, cat := range catList {
			categoryOptions = append(categoryOptions, huh.NewOption(cat, cat))
		}

		fields = []huh.Field{
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
		}
	}

	return newScrollableForm(fields, &data.SubmitConfirmed, width, height)
}

// NewCaptureEventFormForModal creates a huh-based form for modal overlays (WITHOUT confirm button).
//
// Expected:
//   - captureeventformdata must be valid.
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
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
		// Prepare options for tags and categories from constants
		allTags := constants.AllEventTags()
		tagList := make([]string, 0, len(allTags))
		for _, tag := range allTags {
			tagList = append(tagList, string(tag))
		}
		sort.Strings(tagList)

		tagOptions := make([]huh.Option[string], 0, len(tagList))
		for _, tag := range tagList {
			tagOptions = append(tagOptions, huh.NewOption(tag, tag))
		}

		allCats := constants.AllCompetencyCategories()
		catList := make([]string, 0, len(allCats))
		for _, cat := range allCats {
			catList = append(catList, string(cat))
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

// GetCaptureEventFormData extracts form data from an existing event.
//
// Expected:
//   - event must not be nil.
//
// Returns:
//   - A fully initialized CaptureEventFormData ready for use.
//
// Side effects:
//   - None.
func GetCaptureEventFormData(event *career.Event) *CaptureEventFormData {
	tags := event.Tags
	if tags == nil {
		tags = []string{}
	}
	categories := event.Categories
	if categories == nil {
		categories = []string{}
	}

	return &CaptureEventFormData{
		Text:       event.Text,
		Date:       event.Date.Format("2006-01-02"),
		Company:    event.Company,
		Project:    event.Project,
		Tags:       tags,
		Categories: categories,
	}
}
