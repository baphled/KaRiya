package career

import (
	"errors"
	"strings"
	"time"
)

// Skill represents a professional skill or technology competency.
// Skills can be associated with career events to demonstrate technology expertise.
type Skill struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`       // "Ruby", "Kubernetes", "React"
	Category  string     `json:"category"`   // "backend", "frontend", "devops", etc.
	Level     string     `json:"level"`      // "beginner", "intermediate", "advanced", "expert" (optional)
	YearsUsed *int       `json:"years_used"` // Optional years of experience
	LastUsed  *time.Time `json:"last_used"`  // Optional last usage date (can be derived from events)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CommonSkillCategories provides suggested categories for skills.
// Users can define custom categories; these are suggestions, not constraints.
var CommonSkillCategories = []string{
	"backend",
	"frontend",
	"devops",
	"database",
	"cloud",
	"mobile",
	"tooling",
	"other",
}

// AllowedSkillLevels defines the valid skill level values.
var AllowedSkillLevels = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
	"expert":       true,
}

// Validate checks if the Skill meets all defined criteria.
func (s *Skill) Validate() error {
	// Validate name
	if err := s.validateName(); err != nil {
		return err
	}

	// Validate category
	if err := s.validateCategory(); err != nil {
		return err
	}

	// Validate level (optional)
	if err := s.validateLevel(); err != nil {
		return err
	}

	// Validate years used (optional)
	if err := s.validateYearsUsed(); err != nil {
		return err
	}

	// Validate last used (optional)
	if err := s.validateLastUsed(); err != nil {
		return err
	}

	return nil
}

// validateName ensures name is not empty and within length constraints.
func (s *Skill) validateName() error {
	trimmedName := strings.TrimSpace(s.Name)
	if trimmedName == "" {
		return errors.New("name cannot be empty")
	}
	if len(trimmedName) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}

// validateCategory ensures category is not empty and within length constraints.
func (s *Skill) validateCategory() error {
	trimmedCategory := strings.TrimSpace(s.Category)
	if trimmedCategory == "" {
		return errors.New("category cannot be empty")
	}
	if len(trimmedCategory) > 50 {
		return errors.New("category cannot exceed 50 characters")
	}
	return nil
}

// validateLevel ensures level, if provided, is one of the allowed values.
func (s *Skill) validateLevel() error {
	// Level is optional - empty string is allowed
	if s.Level == "" {
		return nil
	}

	if !AllowedSkillLevels[s.Level] {
		return errors.New("level must be one of: beginner, intermediate, advanced, expert")
	}
	return nil
}

// validateYearsUsed ensures years used, if provided, is within valid range.
func (s *Skill) validateYearsUsed() error {
	// YearsUsed is optional
	if s.YearsUsed == nil {
		return nil
	}

	if *s.YearsUsed < 0 || *s.YearsUsed > 50 {
		return errors.New("years_used must be between 0 and 50")
	}
	return nil
}

// validateLastUsed ensures last used date, if provided, is not in the future.
func (s *Skill) validateLastUsed() error {
	// LastUsed is optional
	if s.LastUsed == nil {
		return nil
	}

	now := time.Now()
	if s.LastUsed.After(now) {
		return errors.New("last_used cannot be in the future")
	}
	return nil
}
