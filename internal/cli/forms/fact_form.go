package forms

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// FactFormData holds the form data for fact editing.
type FactFormData struct {
	Text                 string
	CompetencyCategories string // Comma-separated
	RoleFit              string
	AudienceRelevance    string // Comma-separated
	StrengthSignal       string
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

// NewFactEditorForm creates a form for editing a fact.
// The form has 5 fields: Text, Competency Categories, Role Fit, Audience Relevance, Strength Signal,
// plus a Submit confirmation button.
func NewFactEditorForm(fact *career.Fact) *huh.Form {
	data := GetFactFormData(fact)
	data.SubmitConfirmed = false

	roleFitOpts := RoleFitOptions()
	huhRoleFitOpts := make([]huh.Option[string], len(roleFitOpts))
	for i, opt := range roleFitOpts {
		huhRoleFitOpts[i] = huh.NewOption(opt.Value, opt.Key)
	}

	return NewForm(
		huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Fact Text",
				Description: "Factual achievement or accomplishment (1-2000 characters)",
				Placeholder: "Enter fact text...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text),

			huh.NewInput().
				Key("competency_categories").
				Title("Competency Categories").
				Description("Comma-separated competency categories").
				Placeholder("technical, leadership, communication").
				CharLimit(256).
				Value(&data.CompetencyCategories),

			huh.NewSelect[string]().
				Key("role_fit").
				Title("Role Fit").
				Description("Best fit role level for this fact").
				Options(huhRoleFitOpts...).
				Value(&data.RoleFit),

			huh.NewInput().
				Key("audience_relevance").
				Title("Audience Relevance").
				Description("Comma-separated audience types").
				Placeholder("startup, enterprise, technical").
				CharLimit(256).
				Value(&data.AudienceRelevance),

			NewInput(FieldConfig{
				Key:         "strength_signal",
				Title:       "Strength Signal",
				Description: "Signal of strength (0.0-1.0)",
				Placeholder: "0.8",
				Validate: Custom(
					func(val string) bool {
						if val == "" {
							return true
						}
						var f float64
						_, err := fmt.Sscanf(val, "%f", &f)
						return err == nil && f >= 0.0 && f <= 1.0
					},
					"must be a number between 0.0 and 1.0",
				),
			}).Value(&data.StrengthSignal),

			huh.NewConfirm().
				Key("submit").
				Title("Save Changes").
				Description("Submit the form to save your changes").
				Affirmative("Submit").
				Negative("Cancel").
				Value(&data.SubmitConfirmed),
		),
	)
}

// NewFactEditorFormWithData creates a form for editing a fact with initial form data.
func NewFactEditorFormWithData(data *FactFormData) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	roleFitOpts := RoleFitOptions()
	huhRoleFitOpts := make([]huh.Option[string], len(roleFitOpts))
	for i, opt := range roleFitOpts {
		huhRoleFitOpts[i] = huh.NewOption(opt.Value, opt.Key)
	}

	return NewForm(
		huh.NewGroup(
			NewText(FieldConfig{
				Key:         "text",
				Title:       "Fact Text",
				Description: "Factual achievement or accomplishment (1-2000 characters)",
				Placeholder: "Enter fact text...",
				CharLimit:   2000,
				Validate:    Compose(Required, MinLength(10), MaxLength(2000)),
			}).Value(&data.Text),

			huh.NewInput().
				Key("competency_categories").
				Title("Competency Categories").
				Description("Comma-separated competency categories").
				Placeholder("technical, leadership, communication").
				CharLimit(256).
				Value(&data.CompetencyCategories),

			huh.NewSelect[string]().
				Key("role_fit").
				Title("Role Fit").
				Description("Best fit role level for this fact").
				Options(huhRoleFitOpts...).
				Value(&data.RoleFit),

			huh.NewInput().
				Key("audience_relevance").
				Title("Audience Relevance").
				Description("Comma-separated audience types").
				Placeholder("startup, enterprise, technical").
				CharLimit(256).
				Value(&data.AudienceRelevance),

			NewInput(FieldConfig{
				Key:         "strength_signal",
				Title:       "Strength Signal",
				Description: "Signal of strength (0.0-1.0)",
				Placeholder: "0.8",
				Validate: Custom(
					func(val string) bool {
						if val == "" {
							return true
						}
						var f float64
						_, err := fmt.Sscanf(val, "%f", &f)
						return err == nil && f >= 0.0 && f <= 1.0
					},
					"must be a number between 0.0 and 1.0",
				),
			}).Value(&data.StrengthSignal),

			huh.NewConfirm().
				Key("submit").
				Title("Save Changes").
				Description("Submit the form to save your changes").
				Affirmative("Submit").
				Negative("Cancel").
				Value(&data.SubmitConfirmed),
		),
	)
}

// ApplyFactFormData applies the form data to a fact domain object.
func ApplyFactFormData(fact *career.Fact, data *FactFormData) error {
	fact.Text = data.Text
	fact.CompetencyCategories = parseStringSlice(data.CompetencyCategories)
	fact.RoleFit = career.RoleFit(data.RoleFit)
	fact.AudienceRelevance = parseStringSlice(data.AudienceRelevance)
	fact.StrengthSignal = data.StrengthSignal

	return nil
}

// GetFactFormData extracts form data from a fact domain object.
func GetFactFormData(fact *career.Fact) *FactFormData {
	return &FactFormData{
		Text:                 fact.Text,
		CompetencyCategories: formatStringSlice(fact.CompetencyCategories),
		RoleFit:              string(fact.RoleFit),
		AudienceRelevance:    formatStringSlice(fact.AudienceRelevance),
		StrengthSignal:       fact.StrengthSignal,
	}
}

// Helper functions

func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func formatStringSlice(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, ", ")
}
