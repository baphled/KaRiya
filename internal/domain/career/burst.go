package career

import (
	"errors"
	"strings"
	"time"
)

// Burst represents a grouping of related CareerEvents.
type Burst struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	EventIDs    []string   `json:"event_ids"`
	Confirmed   bool       `json:"confirmed"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Validate checks if the Burst meets all defined criteria.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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

	return nil
}

// validateID ensures ID is not empty.
func (b *Burst) validateID() error {
	if strings.TrimSpace(b.ID) == "" {
		return errors.New("ID cannot be empty")
	}
	return nil
}

// validateName ensures name is not empty and within length constraints.
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

// validateEventIDs ensures at least 2 event IDs with no duplicates.
func (b *Burst) validateEventIDs() error {
	if len(b.EventIDs) < 2 {
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

// validateDescription ensures description is within length constraints (optional field).
func (b *Burst) validateDescription() error {
	if b.Description == "" {
		return nil
	}
	if len(b.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}
