package forms

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/huh"
)

// SkillFormData holds the form data for skill editing.
type SkillFormData struct {
	Name            string
	Category        string
	Level           string
	YearsUsed       string
	SubmitConfirmed bool
}

// Common skill categories as suggestions
var CommonSkillCategories = []string{
	"backend",
	"frontend",
	"devops",
	"database",
	"cloud",
	"tooling",
	"other",
}

// Valid skill levels
var ValidSkillLevels = []string{
	"beginner",
	"intermediate",
	"advanced",
	"expert",
}

// NewSkillForm creates a form for adding or editing a skill.
// If skill is nil, creates form for new skill. Otherwise, pre-populates with existing data.
func NewSkillForm(skill *career.Skill) *huh.Form {
	data := &SkillFormData{}
	if skill != nil {
		data = GetSkillFormData(skill)
	}

	return NewForm(
		huh.NewGroup(
			NewInput(FieldConfig{
				Key:         "name",
				Title:       "Skill Name",
				Description: "Name of the skill or technology (required)",
				Placeholder: "e.g., Ruby, Kubernetes, React",
				CharLimit:   100,
				Validate:    SkillName,
			}).Value(&data.Name),

			NewSelect("category", "Category", "Skill category or domain", buildCategoryOptions()).
				Value(&data.Category),

			NewSelect("level", "Proficiency Level", "Your proficiency level (optional)", buildLevelOptions()).
				Value(&data.Level),

			NewInput(FieldConfig{
				Key:         "years",
				Title:       "Years of Experience",
				Description: "Number of years using this skill (optional, 0-50)",
				Placeholder: "e.g., 3",
				CharLimit:   2,
				Validate:    SkillYearsUsed,
			}).Value(&data.YearsUsed),

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

// NewSkillFormWithData creates a form for editing a skill with initial form data.
// This variant allows external data binding for more control.
func NewSkillFormWithData(data *SkillFormData) *huh.Form {
	return NewSkillFormWithDataAndHeight(data, 0)
}

// NewSkillFormWithDataAndHeight creates a form for editing a skill with initial form data and height.
// When height > 0, the form becomes scrollable if content exceeds the height.
func NewSkillFormWithDataAndHeight(data *SkillFormData, height int) *huh.Form {
	return NewSkillFormWithDataAndDimensions(data, 0, height)
}

// NewSkillFormWithDataAndDimensions creates a form for editing a skill with initial form data and dimensions.
// When height > 0, the form becomes scrollable if content exceeds the height.
// When width > 0, the form will be constrained to that width.
// The confirm button is fixed at the bottom, always visible.
func NewSkillFormWithDataAndDimensions(data *SkillFormData, width, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	// Create fields group (scrollable)
	fieldsGroup := huh.NewGroup(
		NewInput(FieldConfig{
			Key:         "name",
			Title:       "Skill Name",
			Description: "Name of the skill or technology (required)",
			Placeholder: "e.g., Ruby, Kubernetes, React",
			CharLimit:   100,
			Validate:    SkillName,
		}).Value(&data.Name),

		NewSelect("category", "Category", "Skill category or domain", buildCategoryOptions()).
			Value(&data.Category),

		NewSelect("level", "Proficiency Level", "Your proficiency level (optional)", buildLevelOptions()).
			Value(&data.Level),

		NewInput(FieldConfig{
			Key:         "years",
			Title:       "Years of Experience",
			Description: "Number of years using this skill (optional, 0-50)",
			Placeholder: "e.g., 3",
			CharLimit:   2,
			Validate:    SkillYearsUsed,
		}).Value(&data.YearsUsed),
	)

	return NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
}

// buildCategoryOptions builds the category select options.
func buildCategoryOptions() []SelectOption {
	options := make([]SelectOption, len(CommonSkillCategories))
	for i, cat := range CommonSkillCategories {
		options[i] = SelectOption{Key: cat, Value: cat}
	}
	return options
}

// buildLevelOptions builds the level select options.
func buildLevelOptions() []SelectOption {
	options := []SelectOption{
		{Key: "", Value: "(not specified)"},
	}
	for _, level := range ValidSkillLevels {
		options = append(options, SelectOption{Key: level, Value: level})
	}
	return options
}

// GetSkillFormData extracts form data from a skill domain object.
func GetSkillFormData(skill *career.Skill) *SkillFormData {
	yearsStr := ""
	if skill.YearsUsed != nil {
		yearsStr = strconv.Itoa(*skill.YearsUsed)
	}

	return &SkillFormData{
		Name:      skill.Name,
		Category:  skill.Category,
		Level:     skill.Level,
		YearsUsed: yearsStr,
	}
}

// ApplySkillFormData applies the form data to a skill domain object.
func ApplySkillFormData(skill *career.Skill, data *SkillFormData) {
	skill.Name = strings.TrimSpace(data.Name)
	skill.Category = strings.TrimSpace(data.Category)
	skill.Level = strings.TrimSpace(data.Level)

	yearsStr := strings.TrimSpace(data.YearsUsed)
	if yearsStr != "" {
		years, err := strconv.Atoi(yearsStr)
		if err == nil && years >= 0 && years <= 50 {
			skill.YearsUsed = &years
		}
	} else {
		skill.YearsUsed = nil
	}
}

// Skill-specific validators

// SkillName validates skill name (required, 1-100 characters).
func SkillName(value string) error {
	return Compose(
		Required,
		LengthRange(1, 100),
	)(value)
}

// SkillCategory validates skill category (required, 1-50 characters).
func SkillCategory(value string) error {
	return Compose(
		Required,
		LengthRange(1, 50),
	)(value)
}

// SkillLevel validates skill level (optional, must be one of valid levels if provided).
func SkillLevel(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil // Optional field
	}

	// Check if value is in valid levels
	for _, valid := range ValidSkillLevels {
		if trimmed == valid {
			return nil
		}
	}

	return fmt.Errorf("must be one of: %s", strings.Join(ValidSkillLevels, ", "))
}

// SkillYearsUsed validates years of experience (optional, 0-50 if provided).
func SkillYearsUsed(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil // Optional field
	}

	// Must be a valid integer
	years, err := strconv.Atoi(trimmed)
	if err != nil {
		return fmt.Errorf("must be a valid number")
	}

	// Range check
	if years < 0 {
		return fmt.Errorf("must be 0 or greater")
	}

	if years > 50 {
		return fmt.Errorf("must be 50 or less")
	}

	return nil
}
