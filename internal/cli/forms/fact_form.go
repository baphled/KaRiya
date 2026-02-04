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
//
// Returns:
//   - A []SelectOption value.
//
// Side effects:
//   - None.
func RoleFitOptions() []SelectOption {
	return []SelectOption{
		{Key: string(career.RoleFitPrincipal), Value: "Principal"},
		{Key: string(career.RoleFitEM), Value: "Engineering Manager"},
		{Key: string(career.RoleFitStaff), Value: "Staff"},
		{Key: string(career.RoleFitSeniorIC), Value: "Senior IC"},
	}
}

// CompetencyCategoryOptions returns all available competency categories.
//
// Returns:
//   - A []huh.Option[string] value.
//
// Side effects:
//   - None.
func CompetencyCategoryOptions() []huh.Option[string] {
	categories := []string{"technical", "leadership", "product", "consulting", "research", "mentoring"}
	options := make([]huh.Option[string], len(categories))
	for i, cat := range categories {
		// Capitalize for display
		displayName := cat
		if cat != "" {
			displayName = string(cat[0]-32) + cat[1:] // Simple capitalize
		}
		options[i] = huh.NewOption(displayName, cat)
	}
	return options
}

// AudienceRelevanceOptions returns all available audience relevance types.
//
// Returns:
//   - A []huh.Option[string] value.
//
// Side effects:
//   - None.
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

// NewFactForm creates a fact edit form bound to the given data.
//
// Expected:
//   - data must be a valid FactFormData pointer.
//   - width and height control form dimensions (0 means unconstrained).
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewFactForm(data *FactFormData, width, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	// Build role fit options
	roleFitOpts := RoleFitOptions()
	huhRoleFitOpts := make([]huh.Option[string], len(roleFitOpts))
	for i, opt := range roleFitOpts {
		huhRoleFitOpts[i] = huh.NewOption(opt.Value, opt.Key)
	}

	fields := []huh.Field{
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
	}

	return newScrollableForm(fields, &data.SubmitConfirmed, width, height)
}

// ApplyFactFormData applies the form data to a fact domain object.
//
// Expected:
//   - fact must be valid.
//   - factformdata must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func ApplyFactFormData(fact *career.Fact, data *FactFormData) error {
	fact.Text = data.Text
	fact.CompetencyCategories = data.CompetencyCategories
	fact.RoleFit = career.RoleFit(data.RoleFit)
	fact.AudienceRelevance = data.AudienceRelevance
	// StrengthSignal is NOT updated - it's auto-generated

	return nil
}

// GetFactFormData extracts form data from a fact domain object.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized FactFormData ready for use.
//
// Side effects:
//   - None.
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
