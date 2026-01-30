package forms

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// MetadataFormData holds the form data for event metadata editing.
type MetadataFormData struct {
	Date            string
	Company         string
	Project         string
	Tags            []string
	Categories      []string
	Skills          []string
	SubmitConfirmed bool
}

// NewMetadataEditorForm creates a form for editing event metadata.
// The form has 6 fields: Date, Company, Project, Tags, Categories, Skills,
// plus a Submit confirmation button.
//
// This form has no height constraint and renders all fields at once.
// For modal overlays where space is limited, use NewMetadataEditorFormWithHeight.
func NewMetadataEditorForm(
	event *career.Event, availableTags, availableCategories []string, availableSkills []*career.Skill,
) *huh.Form {
	data := GetMetadataFormData(event)
	return buildMetadataForm(data, availableTags, availableCategories, availableSkills, 0)
}

// NewMetadataEditorFormWithHeight creates a height-constrained metadata form.
// When height > 0, the huh library enables built-in scrolling so the form
// fits within modal overlays without overflowing.
func NewMetadataEditorFormWithHeight(
	event *career.Event, availableTags, availableCategories []string, availableSkills []*career.Skill,
	height int,
) *huh.Form {
	data := GetMetadataFormData(event)
	return buildMetadataForm(data, availableTags, availableCategories, availableSkills, height)
}

// NewMetadataEditorFormWithData creates a form for editing event metadata with initial form data.
//
// This form has no height constraint and renders all fields at once.
// For modal overlays, use NewMetadataEditorFormWithDataAndHeight.
func NewMetadataEditorFormWithData(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
) *huh.Form {
	return buildMetadataForm(data, availableTags, availableCategories, availableSkills, 0)
}

// NewMetadataEditorFormWithDataAndHeight creates a height-constrained metadata form with initial data.
// When height > 0, the huh library enables built-in scrolling so the form
// fits within modal overlays without overflowing.
func NewMetadataEditorFormWithDataAndHeight(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
	height int,
) *huh.Form {
	return buildMetadataForm(data, availableTags, availableCategories, availableSkills, height)
}

// buildMetadataForm constructs the metadata editor form fields and groups.
//
// When height > 0, the form is created with NewFormWithHeight to enable
// scrolling within constrained containers like modal overlays.
// When height <= 0, the form is created with NewForm (no height constraint).
func buildMetadataForm(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
	height int,
) *huh.Form {
	// Initialize submit confirmation to false.
	data.SubmitConfirmed = false

	// Convert tags and categories to options.
	tagOptions := make([]huh.Option[string], len(availableTags))
	for i, tag := range availableTags {
		tagOptions[i] = huh.NewOption(tag, tag)
	}

	categoryOptions := make([]huh.Option[string], len(availableCategories))
	for i, cat := range availableCategories {
		categoryOptions[i] = huh.NewOption(cat, cat)
	}

	// Convert skills to options (skill ID as value, skill name as label).
	skillOptions := make([]huh.Option[string], len(availableSkills))
	for i, skill := range availableSkills {
		skillOptions[i] = huh.NewOption(skill.Name, skill.ID)
	}

	group := huh.NewGroup(
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

		huh.NewMultiSelect[string]().
			Key("skills").
			Title("Skills").
			Description("Select relevant skills (optional)").
			Options(skillOptions...).
			Value(&data.Skills).
			Limit(10),

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

// ApplyMetadataFormData applies the form data to an event domain object.
func ApplyMetadataFormData(event *career.Event, data *MetadataFormData) error {
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
	event.Skills = data.Skills

	return nil
}

// GetMetadataFormData extracts form data from an event domain object.
func GetMetadataFormData(event *career.Event) *MetadataFormData {
	return &MetadataFormData{
		Date:       event.Date.Format("2006-01-02"),
		Company:    event.Company,
		Project:    event.Project,
		Tags:       event.Tags,
		Categories: event.Categories,
		Skills:     event.Skills,
	}
}
