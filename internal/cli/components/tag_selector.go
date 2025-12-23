package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// TagSelector manages tag selection for career events
type TagSelector struct {
	selected map[string]bool
}

// NewTagSelector creates a new tag selector
func NewTagSelector() *TagSelector {
	return &TagSelector{
		selected: make(map[string]bool),
	}
}

// SelectedTags returns a list of currently selected tags
func (ts *TagSelector) SelectedTags() []string {
	tags := make([]string, 0, len(ts.selected))
	for tag := range ts.selected {
		tags = append(tags, tag)
	}
	return tags
}

// AvailableTags returns all available tags from the domain
func (ts *TagSelector) AvailableTags() []string {
	tags := make([]string, 0, len(career.AllowedTags))
	for tag := range career.AllowedTags {
		tags = append(tags, tag)
	}
	return tags
}

// SelectTag adds a tag to the selected list
func (ts *TagSelector) SelectTag(tag string) error {
	// Validate tag is allowed
	if !career.AllowedTags[tag] {
		return fmt.Errorf("%s is not a valid tag", tag)
	}

	// Check if already selected
	if ts.selected[tag] {
		return fmt.Errorf("tag %s is already selected", tag)
	}

	// Check max limit
	if len(ts.selected) >= 8 {
		return fmt.Errorf("maximum of 8 tags allowed")
	}

	ts.selected[tag] = true
	return nil
}

// DeselectTag removes a tag from the selected list
func (ts *TagSelector) DeselectTag(tag string) error {
	if !ts.selected[tag] {
		return fmt.Errorf("tag %s is not selected", tag)
	}

	delete(ts.selected, tag)
	return nil
}

// FilterTags returns tags that match the given prefix
func (ts *TagSelector) FilterTags(prefix string) []string {
	if prefix == "" {
		return ts.AvailableTags()
	}

	lowerPrefix := strings.ToLower(prefix)
	filtered := make([]string, 0)

	for tag := range career.AllowedTags {
		if strings.HasPrefix(strings.ToLower(tag), lowerPrefix) {
			filtered = append(filtered, tag)
		}
	}

	return filtered
}

// ToggleTag selects the tag if not selected, deselects if already selected
func (ts *TagSelector) ToggleTag(tag string) error {
	if ts.selected[tag] {
		return ts.DeselectTag(tag)
	}
	return ts.SelectTag(tag)
}

// IsSelected returns true if the tag is currently selected
func (ts *TagSelector) IsSelected(tag string) bool {
	return ts.selected[tag]
}

// Reset clears all selected tags
func (ts *TagSelector) Reset() {
	ts.selected = make(map[string]bool)
}

