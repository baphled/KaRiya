package forms

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// MetadataFormData holds the form data for event metadata editing.
type MetadataFormData struct {
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// NewMetadataEditorForm creates a form for editing event metadata.
// The form has 5 fields: Date, Company, Project, Tags, Categories.
func NewMetadataEditorForm(event *career.CareerEvent, availableTags, availableCategories []string) *huh.Form {
	data := GetMetadataFormData(event)

	// Convert tags and categories to options
	tagOptions := make([]huh.Option[string], len(availableTags))
	for i, tag := range availableTags {
		tagOptions[i] = huh.NewOption(tag, tag)
	}

	categoryOptions := make([]huh.Option[string], len(availableCategories))
	for i, cat := range availableCategories {
		categoryOptions[i] = huh.NewOption(cat, cat)
	}

	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "date",
				Title:       "Date",
				Description: "Event date (YYYY-MM-DD, 'today', or '1 week ago')",
				Placeholder: "YYYY-MM-DD or 'today'",
				Validate:    DateFormat,
			}).Value(&data.Date),

			NewInput(FieldConfig{
				Key:         "company",
				Title:       "Company",
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   100,
				Validate:    CompanyName,
			}).Value(&data.Company),

			NewInput(FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project name (optional)",
				Placeholder: "Enter project name...",
				CharLimit:   100,
			}).Value(&data.Project),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("Select relevant tags").
				Options(tagOptions...).
				Value(&data.Tags).
				Limit(10),

			huh.NewMultiSelect[string]().
				Key("categories").
				Title("Categories").
				Description("Select relevant categories").
				Options(categoryOptions...).
				Value(&data.Categories).
				Limit(5),
		),
	)
}

// NewMetadataEditorFormWithData creates a form for editing event metadata with initial form data.
func NewMetadataEditorFormWithData(data *MetadataFormData, availableTags, availableCategories []string) *huh.Form {
	// Convert tags and categories to options
	tagOptions := make([]huh.Option[string], len(availableTags))
	for i, tag := range availableTags {
		tagOptions[i] = huh.NewOption(tag, tag)
	}

	categoryOptions := make([]huh.Option[string], len(availableCategories))
	for i, cat := range availableCategories {
		categoryOptions[i] = huh.NewOption(cat, cat)
	}

	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "date",
				Title:       "Date",
				Description: "Event date (YYYY-MM-DD, 'today', or '1 week ago')",
				Placeholder: "YYYY-MM-DD or 'today'",
				Validate:    DateFormat,
			}).Value(&data.Date),

			NewInput(FieldConfig{
				Key:         "company",
				Title:       "Company",
				Description: "Company name (optional)",
				Placeholder: "Enter company name...",
				CharLimit:   100,
				Validate:    CompanyName,
			}).Value(&data.Company),

			NewInput(FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project name (optional)",
				Placeholder: "Enter project name...",
				CharLimit:   100,
			}).Value(&data.Project),

			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Tags").
				Description("Select relevant tags").
				Options(tagOptions...).
				Value(&data.Tags).
				Limit(10),

			huh.NewMultiSelect[string]().
				Key("categories").
				Title("Categories").
				Description("Select relevant categories").
				Options(categoryOptions...).
				Value(&data.Categories).
				Limit(5),
		),
	)
}

// ApplyMetadataFormData applies the form data to an event domain object.
func ApplyMetadataFormData(event *career.CareerEvent, data *MetadataFormData) error {
	// Parse date string
	parsedDate, err := ParseDateString(data.Date)
	if err != nil {
		return err
	}

	event.Date = parsedDate
	event.Company = data.Company
	event.Project = data.Project
	event.Tags = data.Tags
	event.Categories = data.Categories

	return nil
}

// GetMetadataFormData extracts form data from an event domain object.
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
	return &MetadataFormData{
		Date:       event.Date.Format("2006-01-02"),
		Company:    event.Company,
		Project:    event.Project,
		Tags:       event.Tags,
		Categories: event.Categories,
	}
}
