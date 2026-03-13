package skills

import (
	"strconv"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// SkillInput holds the raw data needed to create a skill.
// This struct decouples skill creation from any form library.
type SkillInput struct {
	Name      string
	Category  string
	Level     string
	YearsUsed string
}

// NewSkillFromInput creates a career.Skill from raw input data.
//
// Expected: Input fields may contain whitespace (trimmed automatically).
// YearsUsed may be empty (nil) or a valid integer string 0-50.
// Returns: A fully initialised career.Skill, or an error if validation fails.
// Side effects: None.
func NewSkillFromInput(input SkillInput, skillID string) (*career.Skill, error) {
	name := strings.TrimSpace(input.Name)
	category := strings.TrimSpace(input.Category)
	level := strings.TrimSpace(input.Level)
	yearsStr := strings.TrimSpace(input.YearsUsed)

	var yearsUsed *int
	if yearsStr != "" {
		years, err := strconv.Atoi(yearsStr)
		if err != nil {
			return nil, &YearsParseError{Value: yearsStr, Err: err}
		}
		yearsUsed = &years
	}

	skill := &career.Skill{
		ID:        skillID,
		Name:      name,
		Category:  category,
		Level:     level,
		YearsUsed: yearsUsed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := skill.Validate(); err != nil {
		return nil, err
	}

	return skill, nil
}

// YearsParseError is returned when a years string cannot be parsed.
type YearsParseError struct {
	Value string
	Err   error
}

// Error returns the error message for the unparseable years input.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (e *YearsParseError) Error() string {
	return "unable to parse years: " + e.Value
}

// Unwrap returns the underlying error.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (e *YearsParseError) Unwrap() error {
	return e.Err
}
