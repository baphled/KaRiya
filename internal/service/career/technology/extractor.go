package technology

import (
	"context"
	"sort"

	careerRepo "github.com/baphled/kariya/internal/repository/career"
)

// ExtractedTechnology represents a technology/skill with event associations.
type ExtractedTechnology struct {
	ID         string   // Skill ID
	Name       string   // "Ruby", "PostgreSQL"
	Category   string   // "backend", "database"
	EventCount int      // How many events use this skill
	EventIDs   []string // Which events
}

// Extractor aggregates and filters user skills with event associations.
type Extractor struct {
	skillRepo careerRepo.SkillRepository
	eventRepo careerRepo.EventRepository
}

// NewExtractor creates a new technology extractor.
func NewExtractor(skillRepo careerRepo.SkillRepository, eventRepo careerRepo.EventRepository) *Extractor {
	return &Extractor{
		skillRepo: skillRepo,
		eventRepo: eventRepo,
	}
}

// ExtractFromUser aggregates user's defined skills with event associations.
func (e *Extractor) ExtractFromUser(ctx context.Context) ([]*ExtractedTechnology, error) {
	// 1. Load all user-defined skills
	skills, err := e.skillRepo.List(ctx, &careerRepo.SkillListFilters{})
	if err != nil {
		return nil, err
	}

	if len(skills) == 0 {
		return []*ExtractedTechnology{}, nil
	}

	// 2. Load all events
	events, err := e.eventRepo.List(ctx, careerRepo.EventListFilters{})
	if err != nil {
		return nil, err
	}

	// 3. Build skill ID -> ExtractedTechnology map
	techMap := make(map[string]*ExtractedTechnology)
	for _, skill := range skills {
		techMap[skill.ID] = &ExtractedTechnology{
			ID:         skill.ID,
			Name:       skill.Name,
			Category:   skill.Category,
			EventCount: 0,
			EventIDs:   []string{},
		}
	}

	// 4. Count events per skill
	for _, event := range events {
		for _, skillID := range event.Skills {
			if tech, exists := techMap[skillID]; exists {
				tech.EventCount++
				tech.EventIDs = append(tech.EventIDs, event.ID)
			}
		}
	}

	// 5. Convert map to slice
	var result []*ExtractedTechnology
	for _, tech := range techMap {
		result = append(result, tech)
	}

	// 6. Sort by event count (descending)
	sort.Slice(result, func(i, j int) bool {
		return result[i].EventCount > result[j].EventCount
	})

	return result, nil
}

// FilterByThreshold removes skills with < minEvents events.
func (e *Extractor) FilterByThreshold(techs []*ExtractedTechnology, minEvents int) []*ExtractedTechnology {
	if len(techs) == 0 {
		return []*ExtractedTechnology{}
	}

	var filtered []*ExtractedTechnology
	for _, tech := range techs {
		if tech.EventCount >= minEvents {
			filtered = append(filtered, tech)
		}
	}

	return filtered
}
