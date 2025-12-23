package validation

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventValidator provides validation for career event inputs
type EventValidator struct{}

// NewEventValidator creates a new EventValidator instance
func NewEventValidator() *EventValidator {
	return &EventValidator{}
}

// ValidateText validates event text field
// Returns error if text is empty, whitespace-only, or exceeds 2000 characters
func (v *EventValidator) ValidateText(text string) error {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return fmt.Errorf("event text is required")
	}
	if len(text) > 2000 {
		return fmt.Errorf("event text cannot exceed 2000 characters")
	}
	return nil
}

// ValidateDate validates that the date is not in the future
func (v *EventValidator) ValidateDate(date time.Time) error {
	if date.After(time.Now()) {
		return fmt.Errorf("date cannot be in the future")
	}
	return nil
}

// ValidateDateForMode validates date based on capture mode constraints
func (v *EventValidator) ValidateDateForMode(date time.Time, mode careerservice.EventCaptureMode) error {
	// First validate date is not in future
	if err := v.ValidateDate(date); err != nil {
		return err
	}

	// Timeline journaling mode requires date within last 30 days
	if mode == careerservice.TimelineJournaling {
		thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
		if date.Before(thirtyDaysAgo) {
			return fmt.Errorf("timeline journaling events must be within the last 30 days")
		}
	}

	// CVBackfill and ManualEntry allow any past date
	return nil
}

// ValidateTags validates tag list for duplicates, count, and allowed values
func (v *EventValidator) ValidateTags(tags []string) error {
	// Check maximum tag count
	if len(tags) > 8 {
		return fmt.Errorf("cannot have more than 8 tags")
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, tag := range tags {
		if seen[tag] {
			return fmt.Errorf("duplicate tag: %s", tag)
		}
		seen[tag] = true
	}

	// Check against allowed tags
	for _, tag := range tags {
		if !isAllowedTag(tag) {
			return fmt.Errorf("invalid tag: %s", tag)
		}
	}

	return nil
}

// isAllowedTag checks if a tag is in the allowed set
func isAllowedTag(tag string) bool {
	return career.AllowedTags[tag]
}

