package technology

import (
	"context"

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
	eventRepo careerRepo.Repository
}

// NewExtractor creates a new technology extractor.
func NewExtractor(skillRepo careerRepo.SkillRepository, eventRepo careerRepo.Repository) *Extractor {
	return &Extractor{
		skillRepo: skillRepo,
		eventRepo: eventRepo,
	}
}

// ExtractFromUser aggregates user's defined skills with event associations.
// TODO: Implement extraction logic (RED phase stub)
func (e *Extractor) ExtractFromUser(ctx context.Context) ([]*ExtractedTechnology, error) {
	return nil, nil
}

// FilterByThreshold removes skills with < minEvents events.
// TODO: Implement filtering logic (RED phase stub)
func (e *Extractor) FilterByThreshold(techs []*ExtractedTechnology, minEvents int) []*ExtractedTechnology {
	return nil
}
