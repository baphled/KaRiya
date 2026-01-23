package selectors

import (
	"fmt"
	"sort"
	"strings"

	domain "github.com/baphled/kariya/internal/domain/career"
)

// SkillSelector manages skill selection for career events
type SkillSelector struct {
	availableSkills []*domain.Skill
	selected        map[string]bool          // skill ID -> selected
	skillsMap       map[string]*domain.Skill // skill ID -> skill
}

// NewSkillSelector creates a new skill selector with available skills
func NewSkillSelector(availableSkills []*domain.Skill) *SkillSelector {
	skillsMap := make(map[string]*domain.Skill)
	for _, skill := range availableSkills {
		skillsMap[skill.ID] = skill
	}

	return &SkillSelector{
		availableSkills: availableSkills,
		selected:        make(map[string]bool),
		skillsMap:       skillsMap,
	}
}

// SelectedSkillIDs returns a sorted list of currently selected skill IDs
func (ss *SkillSelector) SelectedSkillIDs() []string {
	ids := make([]string, 0, len(ss.selected))
	for id := range ss.selected {
		ids = append(ids, id)
	}
	// Sort for consistent output
	sort.Strings(ids)
	return ids
}

// AvailableSkills returns all available skills sorted by name
func (ss *SkillSelector) AvailableSkills() []*domain.Skill {
	skills := make([]*domain.Skill, len(ss.availableSkills))
	copy(skills, ss.availableSkills)

	// Sort by name
	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
	})

	return skills
}

// SelectSkill adds a skill to the selected list by ID
func (ss *SkillSelector) SelectSkill(skillID string) error {
	// Validate skill ID exists
	if _, exists := ss.skillsMap[skillID]; !exists {
		return fmt.Errorf("%s is not a valid skill ID", skillID)
	}

	// Check if already selected
	if ss.selected[skillID] {
		return fmt.Errorf("skill %s is already selected", skillID)
	}

	ss.selected[skillID] = true
	return nil
}

// DeselectSkill removes a skill from the selected list
func (ss *SkillSelector) DeselectSkill(skillID string) error {
	if !ss.selected[skillID] {
		return fmt.Errorf("skill %s is not selected", skillID)
	}

	delete(ss.selected, skillID)
	return nil
}

// ToggleSkill selects the skill if not selected, deselects if already selected
func (ss *SkillSelector) ToggleSkill(skillID string) error {
	if ss.selected[skillID] {
		return ss.DeselectSkill(skillID)
	}
	return ss.SelectSkill(skillID)
}

// IsSelected returns true if the skill is currently selected
func (ss *SkillSelector) IsSelected(skillID string) bool {
	return ss.selected[skillID]
}

// FilterSkills returns skills that match the given prefix (by name) in alphabetical order
func (ss *SkillSelector) FilterSkills(prefix string) []*domain.Skill {
	if prefix == "" {
		return ss.AvailableSkills()
	}

	lowerPrefix := strings.ToLower(prefix)
	filtered := make([]*domain.Skill, 0)

	for _, skill := range ss.availableSkills {
		if strings.HasPrefix(strings.ToLower(skill.Name), lowerPrefix) {
			filtered = append(filtered, skill)
		}
	}

	// Sort by name
	sort.Slice(filtered, func(i, j int) bool {
		return strings.ToLower(filtered[i].Name) < strings.ToLower(filtered[j].Name)
	})

	return filtered
}

// GetSkillByID returns the skill with the given ID
func (ss *SkillSelector) GetSkillByID(skillID string) (*domain.Skill, error) {
	skill, exists := ss.skillsMap[skillID]
	if !exists {
		return nil, fmt.Errorf("skill with ID %s not found", skillID)
	}
	return skill, nil
}

// Reset clears all selected skills
func (ss *SkillSelector) Reset() {
	ss.selected = make(map[string]bool)
}

// SetSelectedSkills sets the selected skills directly from a list of IDs (useful for editing)
func (ss *SkillSelector) SetSelectedSkills(skillIDs []string) {
	ss.selected = make(map[string]bool)
	for _, id := range skillIDs {
		if _, exists := ss.skillsMap[id]; exists {
			ss.selected[id] = true
		}
	}
}
