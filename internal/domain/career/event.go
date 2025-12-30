package career

import (
	"errors"
	"strings"
	"time"
)

// AllowedTags defines the set of valid tags for a CareerEvent
var AllowedTags = map[string]bool{
	"project":     true,
	"achievement": true,
	"leadership":  true,
	"technical":   true,
	"consulting":  true,
	"research":    true,
	"product":     true,
	"mentoring":   true,
}

// AllowedCategories defines the set of valid competency categories for a CareerEvent
var AllowedCategories = map[string]bool{
	"technical":  true,
	"leadership": true,
	"product":    true,
	"consulting": true,
	"research":   true,
	"mentoring":  true,
}

// CareerEvent represents a professional event or milestone
type CareerEvent struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	Date       time.Time `json:"date"`
	Company    string    `json:"company,omitempty"`
	Project    string    `json:"project,omitempty"`
	Tags       []string  `json:"tags"`
	Categories []string  `json:"categories,omitempty"`
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
		if !AllowedTags[tag] {
			return errors.New("invalid tag: " + tag)
		}
	}
	return nil
}

// validateCategories checks that all categories are from the allowed set
func (ce *CareerEvent) validateCategories() error {
	for _, category := range ce.Categories {
		if !AllowedCategories[strings.ToLower(category)] {
			return errors.New("invalid category: " + category)
		}
	}
	return nil
}
