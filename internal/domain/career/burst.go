package career

import (
	"errors"
	"strings"
	"time"
)

// Burst represents a grouping of related CareerEvents
type Burst struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	EventIDs        []string  `json:"event_ids"`
	CompetencyFocus string    `json:"competency_focus,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Validate checks if the Burst meets all defined criteria
func (b *Burst) Validate() error {
	// Validate ID
	if err := b.validateID(); err != nil {
		return err
	}

	// Validate name
	if err := b.validateName(); err != nil {
		return err
	}

	// Validate event IDs
	if err := b.validateEventIDs(); err != nil {
		return err
	}

	// Validate description (optional)
	if err := b.validateDescription(); err != nil {
		return err
	}

	// Validate competency focus (optional)
	if err := b.validateCompetencyFocus(); err != nil {
		return err
	}

	return nil
}

// validateID ensures ID is not empty
func (b *Burst) validateID() error {
	if strings.TrimSpace(b.ID) == "" {
		return errors.New("ID cannot be empty")
	}
	return nil
}

// validateName ensures name is not empty and within length constraints
func (b *Burst) validateName() error {
	trimmedName := strings.TrimSpace(b.Name)
	if trimmedName == "" {
		return errors.New("name cannot be empty")
	}
	if len(trimmedName) > 200 {
		return errors.New("name cannot exceed 200 characters")
	}
	return nil
}

// validateEventIDs ensures at least 2 event IDs with no duplicates
func (b *Burst) validateEventIDs() error {
	if b.EventIDs == nil || len(b.EventIDs) < 2 {
		return errors.New("burst must contain at least 2 events required for a burst")
	}

	// Check for empty event IDs
	seen := make(map[string]bool)
	for _, eventID := range b.EventIDs {
		if strings.TrimSpace(eventID) == "" {
			return errors.New("event ID cannot be empty")
		}
		if seen[eventID] {
			return errors.New("burst contains duplicate event IDs")
		}
		seen[eventID] = true
	}

	return nil
}

// validateDescription ensures description is within length constraints (optional field)
func (b *Burst) validateDescription() error {
	if b.Description == "" {
		return nil // Description is optional
	}
	if len(b.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// validateCompetencyFocus ensures competency focus is valid if provided
func (b *Burst) validateCompetencyFocus() error {
	if b.CompetencyFocus == "" {
		return nil // CompetencyFocus is optional
	}
	if !AllowedCategories[b.CompetencyFocus] {
		return errors.New("invalid competency focus: must be one of the allowed categories")
	}
	return nil
}
