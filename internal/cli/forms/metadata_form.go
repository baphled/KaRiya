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
// For modal overlays where space is limited, use NewMetadataEditorFormWithDimensions.
func NewMetadataEditorForm(
	event *career.Event, availableTags, availableCategories []string, availableSkills []*career.Skill,
) *huh.Form {
	data := GetMetadataFormData(event)
	return buildMetadataForm(data, MetadataFormConfig{
		AvailableTags: availableTags, AvailableCategories: availableCategories, AvailableSkills: availableSkills,
	})
}

// NewMetadataEditorFormWithHeight creates a height-constrained metadata form.
// When height > 0, the form fields scroll and the submit button stays fixed.
func NewMetadataEditorFormWithHeight(
	event *career.Event, availableTags, availableCategories []string, availableSkills []*career.Skill,
	height int,
) *huh.Form {
	data := GetMetadataFormData(event)
	return buildMetadataForm(data, MetadataFormConfig{
		AvailableTags: availableTags, AvailableCategories: availableCategories, AvailableSkills: availableSkills,
		Height: height,
	})
}

// NewMetadataEditorFormWithData creates a form for editing event metadata with initial form data.
//
// This form has no height constraint and renders all fields at once.
// For modal overlays, use NewMetadataEditorFormWithDataAndDimensions.
func NewMetadataEditorFormWithData(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
) *huh.Form {
	return buildMetadataForm(data, MetadataFormConfig{
		AvailableTags: availableTags, AvailableCategories: availableCategories, AvailableSkills: availableSkills,
	})
}

// NewMetadataEditorFormWithDataAndHeight creates a height-constrained metadata form with initial data.
// When height > 0, the form fields scroll and the submit button stays fixed.
func NewMetadataEditorFormWithDataAndHeight(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
	height int,
) *huh.Form {
	return buildMetadataForm(data, MetadataFormConfig{
		AvailableTags: availableTags, AvailableCategories: availableCategories, AvailableSkills: availableSkills,
		Height: height,
	})
}

// MetadataFormConfig bundles the option lists and dimension constraints
// shared by all metadata form constructors.
type MetadataFormConfig struct {
	AvailableTags       []string
	AvailableCategories []string
	AvailableSkills     []*career.Skill
	Width               int
	Height              int
}

// NewMetadataEditorFormWithDataAndDimensions creates a dimension-constrained metadata form.
// When Height > 0, the form fields scroll with the submit button always visible.
// When Width > 0, the form content is constrained to that width.
func NewMetadataEditorFormWithDataAndDimensions(data *MetadataFormData, cfg MetadataFormConfig) *huh.Form {
	return buildMetadataForm(data, cfg)
}

// buildMetadataForm constructs the metadata editor form fields and groups.
//
// When cfg.Height > 0, the form uses NewFormWithFixedConfirm so the submit
// button remains visible while the fields group scrolls independently.
// When cfg.Height <= 0, all fields (including confirm) go in a single group.
func buildMetadataForm(data *MetadataFormData, cfg MetadataFormConfig) *huh.Form {
	data.SubmitConfirmed = false

	fields := metadataFields(data, cfg.AvailableTags, cfg.AvailableCategories, cfg.AvailableSkills)

	if cfg.Height > 0 {
		fieldsGroup := huh.NewGroup(fields...)
		return NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, cfg.Width, cfg.Height)
	}

	// No height constraint: put all fields including confirm in one group.
	confirmField := huh.NewConfirm().
		Key("submit").
		Title("Save Changes").
		Description("Submit the form to save your changes").
		Affirmative("Submit").
		Negative("Cancel").
		Value(&data.SubmitConfirmed)

	allFields := append(fields, confirmField)
	return NewForm(huh.NewGroup(allFields...))
}

// metadataFields builds the editable field list shared by both form layouts.
func metadataFields(
	data *MetadataFormData, availableTags, availableCategories []string, availableSkills []*career.Skill,
) []huh.Field {
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

	return []huh.Field{
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
	}
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
