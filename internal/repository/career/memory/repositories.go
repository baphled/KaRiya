package memory

import (
	career_repo "github.com/baphled/kariya/internal/repository/career"
)

// NewRepositories creates all in-memory repositories for testing.
//
// Returns:
//   - A fully initialized career_repo.Repositories ready for use.
//
// Side effects:
//   - None.
func NewRepositories() *career_repo.Repositories {
	eventRepo := NewEventRepository()
	skillRepo := NewSkillRepository()

	skillRepo.SetEventRepository(eventRepo)
	eventRepo.SetSkillRepository(skillRepo)

	return &career_repo.Repositories{
		Event: eventRepo,
		Skill: skillRepo,
		Fact:  NewFactRepository(),
		Burst: NewBurstRepository(),
	}
}
