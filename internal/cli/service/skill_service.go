package service

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// CLISkillService wraps the skill repository to provide CLI-specific skill operations.
type CLISkillService struct {
	skillRepo careerrepo.SkillRepository
}

// NewCLISkillService creates a new CLI skill service.
//
// Expected:
//   - skillRepo must be a valid SkillRepository instance.
//
// Returns:
//   - A fully initialized CLISkillService ready for use.
//
// Side effects:
//   - None.
func NewCLISkillService(skillRepo careerrepo.SkillRepository) *CLISkillService {
	return &CLISkillService{
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
func (c *CLISkillService) Create(ctx context.Context, skill *career.Skill) error {
	return c.skillRepo.Create(ctx, skill)
}
