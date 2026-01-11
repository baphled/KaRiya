package forms

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// FactFormData holds the form data for fact editing.
type FactFormData struct {
	Text                 string
	CompetencyCategories []string // Multi-select
	RoleFit              string
	AudienceRelevance    []string // Multi-select
	SubmitConfirmed      bool
}

// RoleFitOptions returns the available role fit options.
func RoleFitOptions() []SelectOption {
	return []SelectOption{
		{Key: string(career.RoleFitPrincipal), Value: "Principal"},
		{Key: string(career.RoleFitEM), Value: "Engineering Manager"},
		{Key: string(career.RoleFitStaff), Value: "Staff"},
		{Key: string(career.RoleFitSeniorIC), Value: "Senior IC"},
	}
}

// CompetencyCategoryOptions returns all available competency categories.
func CompetencyCategoryOptions() []huh.Option[string] {
	categories := []string{"technical", "leadership", "product", "consulting", "research", "mentoring"}
	options := make([]huh.Option[string], len(categories))
	for i, cat := range categories {
		// Capitalize for display
		displayName := cat
		if len(cat) > 0 {
			displayName = string(cat[0]-32) + cat[1:] // Simple capitalize
		}
		options[i] = huh.NewOption(displayName, cat)
	}
	return options
}

// AudienceRelevanceOptions returns all available audience relevance types.
func AudienceRelevanceOptions() []huh.Option[string] {
	audiences := []struct {
		key   string
		label string
	}{
		{"hiring_manager", "Hiring Manager"},
		{"recruiter", "Recruiter"},
		{"peer", "Peer"},
	}
	options := make([]huh.Option[string], len(audiences))
	for i, aud := range audiences {
		options[i] = huh.NewOption(aud.label, aud.key)
	}
	return options
}

// NewFactEditorForm creates a form for editing a fact.
// The form has 4 fields: Text, Competency Categories, Role Fit, Audience Relevance,
// plus a Submit confirmation button.
// Note: StrengthSignal is auto-generated and not user-editable.
func NewFactEditorForm(fact *career.Fact) *huh.Form {
	data := GetFactFormData(fact)
	return NewFactEditorFormWithData(data)
}

// NewFactEditorFormWithData creates a form for editing a fact with initial form data.
func NewFactEditorFormWithData(data *FactFormData) *huh.Form {
	return NewFactEditorFormWithDataAndHeight(data, 0)
}

// NewFactEditorFormWithDataAndHeight creates a form for editing a fact with initial form data and height.
// When height > 0, the form becomes scrollable if content exceeds the height.
func NewFactEditorFormWithDataAndHeight(data *FactFormData, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	// Build role fit options
	roleFitOpts := RoleFitOptions()
	huhRoleFitOpts := make([]huh.Option[string], len(roleFitOpts))
	for i, opt := range roleFitOpts {
		huhRoleFitOpts[i] = huh.NewOption(opt.Value, opt.Key)
	}

	group := huh.NewGroup(
		NewText(FieldConfig{
			Key:         "text",
			Title:       "Fact Text",
			Description: "Factual achievement or accomplishment (10-2000 characters)",
			Placeholder: "Enter fact text...",
			CharLimit:   2000,
			Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
		}).Value(&data.Text),

		huh.NewMultiSelect[string]().
			Key("competency_categories").
			Title("Competency Categories").
			Description("Select relevant competency categories").
			Options(CompetencyCategoryOptions()...).
			Value(&data.CompetencyCategories).
			Limit(6),

		huh.NewSelect[string]().
			Key("role_fit").
			Title("Role Fit").
			Description("Best fit role level for this fact").
			Options(huhRoleFitOpts...).
			Value(&data.RoleFit),

		huh.NewMultiSelect[string]().
			Key("audience_relevance").
			Title("Audience Relevance").
			Description("Select target audience types").
			Options(AudienceRelevanceOptions()...).
			Value(&data.AudienceRelevance).
			Limit(3),

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

// ApplyFactFormData applies the form data to a fact domain object.
// Note: StrengthSignal is preserved from the original fact as it's auto-generated.
func ApplyFactFormData(fact *career.Fact, data *FactFormData) error {
	fact.Text = data.Text
	fact.CompetencyCategories = data.CompetencyCategories
	fact.RoleFit = career.RoleFit(data.RoleFit)
	fact.AudienceRelevance = data.AudienceRelevance
	// StrengthSignal is NOT updated - it's auto-generated

	return nil
}

// GetFactFormData extracts form data from a fact domain object.
func GetFactFormData(fact *career.Fact) *FactFormData {
	// Ensure slices are not nil for proper multi-select binding
	categories := fact.CompetencyCategories
	if categories == nil {
		categories = []string{}
	}
	audiences := fact.AudienceRelevance
	if audiences == nil {
		audiences = []string{}
	}

	return &FactFormData{
		Text:                 fact.Text,
		CompetencyCategories: categories,
		RoleFit:              string(fact.RoleFit),
		AudienceRelevance:    audiences,
	}
}
