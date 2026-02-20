package capture

import (
	"errors"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// BurstEditInput holds the data needed to edit a burst.
type BurstEditInput struct {
	Name        string
	Description string
}

// BurstSuggestionInput holds the data from a burst suggestion.
type BurstSuggestionInput struct {
	Name        string
	Description string
	EventIDs    []string
}

// ApplyBurstEdit applies edit changes to an existing burst, modifying it in place.
//
// Expected: Burst must not be nil. Input.Name must not be empty.
// Returns: The modified burst, or an error if validation fails.
// Side effects: Modifies burst in place (Name, Description, UpdatedAt).
func ApplyBurstEdit(burst *career.Burst, input BurstEditInput) (*career.Burst, error) {
	if strings.TrimSpace(input.Name) == "" {
		return burst, errors.New("burst name cannot be empty")
	}

	burst.Name = input.Name
	burst.Description = input.Description
	burst.UpdatedAt = time.Now()

	return burst, nil
}

// CreateBurstFromSuggestion creates a new confirmed burst from a suggestion.
//
// Expected: Input.EventIDs should contain at least one event ID.
// Returns: A new career.Burst with Confirmed=true and no ID (assigned by repo).
// Side effects: None.
func CreateBurstFromSuggestion(input BurstSuggestionInput) *career.Burst {
	name := input.Name
	if name == "" {
		name = "Untitled burst"
	}

	return &career.Burst{
		Name:        name,
		Description: input.Description,
		EventIDs:    input.EventIDs,
		Confirmed:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ConfirmBurstTimestamps sets Confirmed=true and updates timestamps on a burst in place.
//
// Expected: Burst must not be nil.
// Returns: The modified burst with confirmation timestamps set.
// Side effects: Modifies burst in place.
func ConfirmBurstTimestamps(burst *career.Burst) *career.Burst {
	now := time.Now()
	burst.Confirmed = true
	burst.ConfirmedAt = &now
	burst.UpdatedAt = now
	return burst
}
