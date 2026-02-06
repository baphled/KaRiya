package career

import (
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/constants"
)

// Event represents a professional event or milestone.
type Event struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	Date       time.Time `json:"date"`
	Company    string    `json:"company,omitempty"`
	Project    string    `json:"project,omitempty"`
	Tags       []string  `json:"tags"`
	Categories []string  `json:"categories,omitempty"`
	Skills     []string  `json:"skills,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate checks if the Event meets all defined criteria.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (ce *Event) Validate() error {
	// Validate text
	if err := ce.validateText(); err != nil {
		return err
	}

	// Validate date
	if err := ce.validateDate(); err != nil {
		return err
	}

	// Validate tags
	if err := ce.validateTags(); err != nil {
		return err
	}

	// Validate categories
	if err := ce.validateCategories(); err != nil {
		return err
	}

	return nil
}

// validateText ensures text is not empty and within length constraints.
// Text must be between 10 and 2000 characters.
func (ce *Event) validateText() error {
	trimmedText := strings.TrimSpace(ce.Text)
	if trimmedText == "" {
		return errors.New("text cannot be empty")
	}
	if len(trimmedText) < 10 {
		return errors.New("text must be at least 10 characters")
	}
	if len(trimmedText) > 2000 {
		return errors.New("text cannot exceed 2000 characters")
	}
	return nil
}

// validateDate ensures date is not in the future.
func (ce *Event) validateDate() error {
	now := time.Now()
	if ce.Date.After(now) {
		return errors.New("date cannot be in the future")
	}
	return nil
}

// validateTags checks that all tags are from the allowed set.
func (ce *Event) validateTags() error {
	for _, tag := range ce.Tags {
		if !constants.IsValidEventTag(tag) {
			return errors.New("invalid tag: " + tag)
		}
	}
	return nil
}

// validateCategories checks that all categories are from the allowed set.
func (ce *Event) validateCategories() error {
	for _, category := range ce.Categories {
		if !constants.IsValidCompetencyCategory(strings.ToLower(category)) {
			return errors.New("invalid category: " + category)
		}
	}
	return nil
}
