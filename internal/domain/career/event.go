package career

import (
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/constants"
)

// CareerEvent represents a professional event or milestone
type CareerEvent struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	Date       time.Time `json:"date"`
	Company    string    `json:"company,omitempty"`
	Project    string    `json:"project,omitempty"`
	Tags       []string  `json:"tags"`
	Categories []string  `json:"categories,omitempty"`
	Skills     []string  `json:"skills,omitempty"` // Skill IDs associated with this event
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate checks if the CareerEvent meets all defined criteria
func (ce *CareerEvent) Validate() error {
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

// validateText ensures text is not empty and within length constraints
func (ce *CareerEvent) validateText() error {
	trimmedText := strings.TrimSpace(ce.Text)
	if trimmedText == "" {
		return errors.New("text cannot be empty")
	}
	if len(trimmedText) > 2000 {
		return errors.New("text cannot exceed 2000 characters")
	}
	return nil
}

// validateDate ensures date is not in the future
func (ce *CareerEvent) validateDate() error {
	now := time.Now()
	if ce.Date.After(now) {
		return errors.New("date cannot be in the future")
	}
	return nil
}

// validateTags checks that all tags are from the allowed set
func (ce *CareerEvent) validateTags() error {
	for _, tag := range ce.Tags {
		if !constants.IsValidEventTag(tag) {
			return errors.New("invalid tag: " + tag)
		}
	}
	return nil
}

// validateCategories checks that all categories are from the allowed set
func (ce *CareerEvent) validateCategories() error {
	for _, category := range ce.Categories {
		if !constants.IsValidCompetencyCategory(strings.ToLower(category)) {
			return errors.New("invalid category: " + category)
		}
	}
	return nil
}
