package memory

import (
	career_repo "github.com/baphled/kariya/internal/repository/career"
)

// NewRepositories creates all in-memory repositories for testing.
func NewRepositories() *career_repo.Repositories {
	return &career_repo.Repositories{
		Event: NewEventRepository(),
		Skill: NewSkillRepository(),
		Fact:  NewFactRepository(),
		Burst: NewBurstRepository(),
	}
}
