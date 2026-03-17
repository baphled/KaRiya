package forms

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/constants"
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

// NewSkillForm creates a skill add/edit form bound to the given data.
//
// Expected:
//   - data must be a valid SkillFormData pointer.
//   - width and height control form dimensions (0 means unconstrained).
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewSkillForm(data *SkillFormData, width, height int) *huh.Form {
	// Initialize submit confirmation to false
	data.SubmitConfirmed = false

	fields := []huh.Field{
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
	}

	return newScrollableForm(fields, &data.SubmitConfirmed, width, height)
}

// buildCategoryOptions builds the category select options from centralized constants.
func buildCategoryOptions() []SelectOption {
	categories := constants.SuggestedSkillCategories()
	options := make([]SelectOption, len(categories))
	for i, cat := range categories {
		options[i] = SelectOption{Key: string(cat), Value: string(cat)}
	}
	return options
}

// buildLevelOptions builds the level select options from centralized constants.
func buildLevelOptions() []SelectOption {
	options := []SelectOption{
		{Key: "", Value: "(not specified)"},
	}
	for _, level := range constants.AllSkillLevels() {
		options = append(options, SelectOption{Key: string(level), Value: string(level)})
	}
	return options
}

// GetSkillFormData extracts form data from a skill domain object.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized SkillFormData ready for use.
//
// Side effects:
//   - None.
func GetSkillFormData(skill *career.Skill) *SkillFormData {
	if skill == nil {
		return &SkillFormData{}
	}

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
//
// Expected:
//   - skill must be valid.
//   - skillformdata must be valid.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SkillName(value string) error {
	return Compose(
		Required,
		LengthRange(1, 100),
	)(value)
}

// SkillCategory validates skill category (required, 1-50 characters).
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SkillCategory(value string) error {
	return Compose(
		Required,
		LengthRange(1, 50),
	)(value)
}

// SkillLevel validates skill level (optional, must be one of valid levels if provided).
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SkillLevel(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil // Optional field
	}

	// Use centralized skill level validation
	if constants.IsValidSkillLevel(trimmed) {
		return nil
	}

	// Build error message from all valid levels
	levels := constants.AllSkillLevels()
	levelStrs := make([]string, len(levels))
	for i, l := range levels {
		levelStrs[i] = string(l)
	}
	return fmt.Errorf("must be one of: %s", strings.Join(levelStrs, ", "))
}

// SkillYearsUsed validates years of experience (optional, 0-50 if provided).
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SkillYearsUsed(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil // Optional field
	}

	// Must be a valid integer
	years, err := strconv.Atoi(trimmed)
	if err != nil {
		return errors.New("must be a valid number")
	}

	// Range check
	if years < 0 {
		return errors.New("must be 0 or greater")
	}

	if years > 50 {
		return errors.New("must be 50 or less")
	}

	return nil
}
