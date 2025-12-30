package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/service/career/classification"
)

// AllowedCategories defines the set of valid competency categories
var AllowedCategories = map[string]bool{
	"technical":  true,
	"leadership": true,
	"product":    true,
	"consulting": true,
	"research":   true,
	"mentoring":  true,
}

// CategorySelector manages category selection for career events
type CategorySelector struct {
	selected map[string]bool
}

// NewCategorySelector creates a new category selector
func NewCategorySelector() *CategorySelector {
	return &CategorySelector{
		selected: make(map[string]bool),
	}
}

// SelectedCategories returns a sorted list of currently selected categories
func (cs *CategorySelector) SelectedCategories() []string {
	categories := make([]string, 0, len(cs.selected))
	for category := range cs.selected {
		categories = append(categories, category)
	}
	// Sort for consistent output
	sort.Strings(categories)
	return categories
}

// AvailableCategories returns all available categories
func (cs *CategorySelector) AvailableCategories() []string {
	categories := make([]string, 0, len(AllowedCategories))
	for category := range AllowedCategories {
		categories = append(categories, category)
	}
	// Sort for consistent output
	sort.Strings(categories)
	return categories
}

// SelectCategory adds a category to the selected list
func (cs *CategorySelector) SelectCategory(category string) error {
	// Validate category is allowed
	if !AllowedCategories[strings.ToLower(category)] {
		return fmt.Errorf("%s is not a valid category", category)
	}

	categoryLower := strings.ToLower(category)
	// Check if already selected
	if cs.selected[categoryLower] {
		return fmt.Errorf("category %s is already selected", category)
	}

	// No hard limit on categories (unlike tags which have max 8)
	cs.selected[categoryLower] = true
	return nil
}

// DeselectCategory removes a category from the selected list
func (cs *CategorySelector) DeselectCategory(category string) error {
	categoryLower := strings.ToLower(category)
	if !cs.selected[categoryLower] {
		return fmt.Errorf("category %s is not selected", category)
	}

	delete(cs.selected, categoryLower)
	return nil
}

// FilterCategories returns categories that match the given prefix
func (cs *CategorySelector) FilterCategories(prefix string) []string {
	if prefix == "" {
		return cs.AvailableCategories()
	}

	lowerPrefix := strings.ToLower(prefix)
	filtered := make([]string, 0)

	for category := range AllowedCategories {
		if strings.HasPrefix(strings.ToLower(category), lowerPrefix) {
			filtered = append(filtered, category)
		}
	}

	// Sort for consistent output
	sort.Strings(filtered)
	return filtered
}

// ToggleCategory selects the category if not selected, deselects if already selected
func (cs *CategorySelector) ToggleCategory(category string) error {
	categoryLower := strings.ToLower(category)
	if cs.selected[categoryLower] {
		return cs.DeselectCategory(category)
	}
	return cs.SelectCategory(category)
}

// IsSelected returns true if the category is currently selected
func (cs *CategorySelector) IsSelected(category string) bool {
	return cs.selected[strings.ToLower(category)]
}

// Clear removes all selected categories
func (cs *CategorySelector) Clear() {
	cs.selected = make(map[string]bool)
}

// SetSelected sets the selected categories from a list
func (cs *CategorySelector) SetSelected(categories []string) error {
	cs.Clear()
	for _, category := range categories {
		if err := cs.SelectCategory(category); err != nil {
			return err
		}
	}
	return nil
}

// GetCategoryDescription returns a human-friendly description of a category
func GetCategoryDescription(category string) string {
	descriptions := map[string]string{
		"technical":  "Technical skills and engineering work",
		"leadership": "Leadership and management experience",
		"product":    "Product management and strategy",
		"consulting": "Consulting and advisory work",
		"research":   "Research and investigation",
		"mentoring":  "Mentoring and coaching others",
	}
	if desc, ok := descriptions[strings.ToLower(category)]; ok {
		return desc
	}
	return category
}

// MapToClassificationCategories converts string categories to classification.CompetencyCategory
// This is useful for compatibility with the classification service
func MapToClassificationCategories(categories []string) []classification.CompetencyCategory {
	result := make([]classification.CompetencyCategory, 0, len(categories))
	for _, category := range categories {
		switch strings.ToLower(category) {
		case "technical":
			result = append(result, classification.TechnicalCompetency)
		case "leadership":
			result = append(result, classification.LeadershipCompetency)
		case "product":
			result = append(result, classification.ProductCompetency)
		case "consulting":
			result = append(result, classification.ConsultingCompetency)
		case "research":
			result = append(result, classification.ResearchCompetency)
		case "mentoring":
			result = append(result, classification.MentoringCompetency)
		}
	}
	return result
}
