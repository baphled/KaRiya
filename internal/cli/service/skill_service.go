package service

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// CLISkillCreator wraps the skill repository to provide CLI-specific skill operations.
type CLISkillCreator struct {
	skillRepo careerrepo.SkillRepository
}

// NewCLISkillCreator creates a new CLI skill service.
//
// Expected:
//   - skillRepo must be a valid SkillRepository instance.
//
// Returns:
//   - A fully initialized CLISkillCreator ready for use.
//
// Side effects:
//   - None.
func NewCLISkillCreator(skillRepo careerrepo.SkillRepository) *CLISkillCreator {
	return &CLISkillCreator{
		skillRepo: skillRepo,
	}
}

// Create creates a new skill in the repository.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - skill must be a valid skill object.
//
// Returns:
//   - An error value if creation failed.
//
// Side effects:
//   - Creates a new skill in the database.
func (c *CLISkillCreator) Create(ctx context.Context, skill *career.Skill) error {
	return c.skillRepo.Create(ctx, skill)
}
