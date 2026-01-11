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
	Name      string
	Category  string
	Level     string
	YearsUsed string
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
// categoryOptions provides suggestions for category field (can be nil).
func NewSkillForm(skill *career.Skill, categoryOptions []string) *huh.Form {
	var data *SkillFormData
	if skill != nil {
		data = GetSkillFormData(skill)
	} else {
		data = &SkillFormData{}
	}

	return NewSkillFormWithData(data, categoryOptions)
}

// NewSkillFormWithData creates a form for editing a skill with initial form data.
// This variant allows external data binding for more control.
// categoryOptions provides suggestions for category field (can be nil).
func NewSkillFormWithData(data *SkillFormData, categoryOptions []string) *huh.Form {
	// Use provided categories or fall back to common ones
	categories := categoryOptions
	if categories == nil {
		categories = CommonSkillCategories
	}

	// Build category select options
	categorySelectOptions := make([]SelectOption, len(categories))
	for i, cat := range categories {
		categorySelectOptions[i] = SelectOption{Key: cat, Value: cat}
	}

	// Build level select options (with empty option for "not specified")
	levelOptions := []SelectOption{
		{Key: "", Value: "(not specified)"},
	}
	for _, level := range ValidSkillLevels {
		levelOptions = append(levelOptions, SelectOption{Key: level, Value: level})
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

			NewSelect("category", "Category", "Skill category or domain", categorySelectOptions).
				Value(&data.Category),

			NewSelect("level", "Proficiency Level", "Your proficiency level (optional)", levelOptions).
				Value(&data.Level),

			NewInput(FieldConfig{
				Key:         "years",
				Title:       "Years of Experience",
				Description: "Number of years using this skill (optional, 0-50)",
				Placeholder: "e.g., 3",
				CharLimit:   2,
				Validate:    SkillYearsUsed,
			}).Value(&data.YearsUsed),
		),
	)
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
// Validates data and returns error if invalid.
func ApplySkillFormData(skill *career.Skill, data *SkillFormData) error {
	// Trim whitespace
	name := strings.TrimSpace(data.Name)
	category := strings.TrimSpace(data.Category)
	level := strings.TrimSpace(data.Level)
	yearsStr := strings.TrimSpace(data.YearsUsed)

	// Apply required fields
	skill.Name = name
	skill.Category = category
	skill.Level = level

	// Parse and validate years if provided
	if yearsStr != "" {
		years, err := strconv.Atoi(yearsStr)
		if err != nil {
			return fmt.Errorf("invalid years format: must be a number")
		}

		if years < 0 || years > 50 {
			return fmt.Errorf("years must be between 0 and 50")
		}

		skill.YearsUsed = &years
	} else {
		skill.YearsUsed = nil
	}

	return nil
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
