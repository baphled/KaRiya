package selectors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/constants"
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

// SelectedTags returns a sorted list of currently selected tags
func (ts *TagSelector) SelectedTags() []string {
	tags := make([]string, 0, len(ts.selected))
	for tag := range ts.selected {
		tags = append(tags, tag)
	}
	// Sort for consistent output
	sort.Strings(tags)
	return tags
}

// AvailableTags returns all available tags from the domain in alphabetical order
func (ts *TagSelector) AvailableTags() []string {
	allTags := constants.AllEventTags()
	tags := make([]string, 0, len(allTags))
	for _, tag := range allTags {
		tags = append(tags, string(tag))
	}
	// Sort for consistent output
	sort.Strings(tags)
	return tags
}

// SelectTag adds a tag to the selected list
func (ts *TagSelector) SelectTag(tag string) error {
	// Validate tag is allowed
	if !constants.IsValidEventTag(tag) {
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

// FilterTags returns tags that match the given prefix in alphabetical order
func (ts *TagSelector) FilterTags(prefix string) []string {
	if prefix == "" {
		return ts.AvailableTags()
	}

	lowerPrefix := strings.ToLower(prefix)
	filtered := make([]string, 0)

	for _, tag := range constants.AllEventTags() {
		tagStr := string(tag)
		if strings.HasPrefix(strings.ToLower(tagStr), lowerPrefix) {
			filtered = append(filtered, tagStr)
		}
	}

	// Sort for consistent output
	sort.Strings(filtered)
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

// SetSelectedTags sets the selected tags directly (useful for editing)
func (ts *TagSelector) SetSelectedTags(tags []string) {
	ts.selected = make(map[string]bool)
	for _, tag := range tags {
		if constants.IsValidEventTag(tag) && len(ts.selected) < 8 {
			ts.selected[tag] = true
		}
	}
}
