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

// MetadataFormConfig bundles the option lists and dimension constraints
// for the metadata form.
type MetadataFormConfig struct {
	AvailableTags       []string
	AvailableCategories []string
	AvailableSkills     []*career.Skill
	Width               int
	Height              int
}

// NewMetadataForm creates a metadata edit form bound to the given data and config.
//
// Expected:
//   - data must be a valid MetadataFormData pointer.
//   - cfg provides available options and dimension constraints.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewMetadataForm(data *MetadataFormData, cfg MetadataFormConfig) *huh.Form {
	return buildMetadataForm(data, cfg)
}

// buildMetadataForm constructs the metadata editor form fields and groups.
//
// When cfg.Height > 0, the form uses newScrollableForm so the submit
// button remains visible while the fields group scrolls independently.
// When cfg.Height <= 0, all fields (including confirm) go in a single group.
func buildMetadataForm(data *MetadataFormData, cfg MetadataFormConfig) *huh.Form {
	data.SubmitConfirmed = false

	fields := metadataFields(data, cfg.AvailableTags, cfg.AvailableCategories, cfg.AvailableSkills)

	if cfg.Height > 0 {
		return newScrollableForm(fields, &data.SubmitConfirmed, cfg.Width, cfg.Height)
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
	return newForm(huh.NewGroup(allFields...))
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
//
// Expected:
//   - event must be valid.
//   - metadataformdata must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized MetadataFormData ready for use.
//
// Side effects:
//   - None.
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
