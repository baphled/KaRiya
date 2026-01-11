package career

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("SQLiteSkillRepository", func() {
	var (
		repository SkillRepository
		ctx        context.Context
		tempDir    string
		dbPath     string
	)

	BeforeEach(func() {
		var err error
		// Create a temporary directory for test database
		tempDir, err = os.MkdirTemp("", "kariya-skill-test-")
		Expect(err).NotTo(HaveOccurred())

		// Create database path
		dbPath = filepath.Join(tempDir, "test_skills.db")

		// Create main repository to run migrations
		mainRepo, err := NewSQLiteRepository(dbPath)
		Expect(err).NotTo(HaveOccurred())

		// Create skill repository with the same database connection
		repository = NewSQLiteSkillRepositoryWithDB(mainRepo.GetDB())

		ctx = context.Background()
	})

	AfterEach(func() {
		// Clean up temporary files
		os.RemoveAll(tempDir)
	})

	Describe("Create", func() {
		It("should create a new skill", func() {
			skill := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
			}

			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).NotTo(BeEmpty())
			Expect(skill.CreatedAt).NotTo(BeZero())
			Expect(skill.UpdatedAt).NotTo(BeZero())
		})

		It("should return error for duplicate skill name", func() {
			skill1 := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
			}

			err := repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())

			// Try to create another skill with the same name
			skill2 := &career.Skill{
				Name:     "Ruby",
				Category: "frontend",
			}

			err = repository.Create(ctx, skill2)
			Expect(err).To(MatchError(ErrDuplicateSkill))
		})

		It("should generate unique ID if not provided", func() {
			skill1 := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
			}

			skill2 := &career.Skill{
				Name:     "Go",
				Category: "backend",
			}

			err1 := repository.Create(ctx, skill1)
			err2 := repository.Create(ctx, skill2)

			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(skill1.ID).NotTo(BeEmpty())
			Expect(skill2.ID).NotTo(BeEmpty())
			Expect(skill1.ID).NotTo(Equal(skill2.ID))
		})
	})

	Describe("GetByID", func() {
		It("should retrieve an existing skill", func() {
			skill := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
				Level:    "advanced",
			}

			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := repository.GetByID(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(skill.ID))
			Expect(retrieved.Name).To(Equal(skill.Name))
			Expect(retrieved.Category).To(Equal(skill.Category))
			Expect(retrieved.Level).To(Equal(skill.Level))
		})

		It("should return error for non-existent skill", func() {
			_, err := repository.GetByID(ctx, "non-existent-id")
			Expect(err).To(MatchError(ErrSkillNotFound))
		})
	})

	Describe("GetByName", func() {
		It("should retrieve a skill by name", func() {
			skill := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
			}

			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := repository.GetByName(ctx, "Ruby")
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(skill.ID))
			Expect(retrieved.Name).To(Equal(skill.Name))
		})

		It("should return error for non-existent skill name", func() {
			_, err := repository.GetByName(ctx, "NonExistent")
			Expect(err).To(MatchError(ErrSkillNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create test skills
			skills := []*career.Skill{
				{Name: "Ruby", Category: "backend", Level: "advanced"},
				{Name: "Go", Category: "backend", Level: "intermediate"},
				{Name: "React", Category: "frontend", Level: "advanced"},
				{Name: "Vue", Category: "frontend", Level: "beginner"},
			}

			for _, skill := range skills {
				err := repository.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should list all skills with no filters", func() {
			skills, err := repository.List(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(4))
		})

		It("should filter skills by category", func() {
			filters := &SkillFilters{
				Category: "backend",
			}

			skills, err := repository.List(ctx, filters)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, skill := range skills {
				Expect(skill.Category).To(Equal("backend"))
			}
		})

		It("should filter skills by level", func() {
			filters := &SkillFilters{
				Level: "advanced",
			}

			skills, err := repository.List(ctx, filters)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, skill := range skills {
				Expect(skill.Level).To(Equal("advanced"))
			}
		})

		It("should apply pagination with limit", func() {
			filters := &SkillFilters{
				Limit: 2,
			}

			skills, err := repository.List(ctx, filters)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("should apply pagination with offset", func() {
			filters := &SkillFilters{
				Offset: 2,
			}

			skills, err := repository.List(ctx, filters)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})
	})

	Describe("Update", func() {
		It("should update an existing skill", func() {
			skill := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
				Level:    "intermediate",
			}

			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			// Update skill
			skill.Level = "advanced"
			years := 5
			skill.YearsUsed = &years

			err = repository.Update(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			// Verify update
			retrieved, err := repository.GetByID(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Level).To(Equal("advanced"))
			Expect(*retrieved.YearsUsed).To(Equal(5))
			Expect(retrieved.UpdatedAt.After(retrieved.CreatedAt)).To(BeTrue())
		})

		It("should return error for non-existent skill", func() {
			skill := &career.Skill{
				ID:       "non-existent-id",
				Name:     "Ruby",
				Category: "backend",
			}

			err := repository.Update(ctx, skill)
			Expect(err).To(MatchError(ErrSkillNotFound))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing skill", func() {
			skill := &career.Skill{
				Name:     "Ruby",
				Category: "backend",
			}

			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			// Delete skill
			err = repository.Delete(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify deletion
			_, err = repository.GetByID(ctx, skill.ID)
			Expect(err).To(MatchError(ErrSkillNotFound))
		})

		It("should return error for non-existent skill", func() {
			err := repository.Delete(ctx, "non-existent-id")
			Expect(err).To(MatchError(ErrSkillNotFound))
		})
	})

	Describe("GetByCategory", func() {
		BeforeEach(func() {
			skills := []*career.Skill{
				{Name: "Ruby", Category: "backend"},
				{Name: "Go", Category: "backend"},
				{Name: "React", Category: "frontend"},
			}

			for _, skill := range skills {
				err := repository.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should retrieve all skills in a category", func() {
			skills, err := repository.GetByCategory(ctx, "backend")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, skill := range skills {
				Expect(skill.Category).To(Equal("backend"))
			}
		})

		It("should return empty list for non-existent category", func() {
			skills, err := repository.GetByCategory(ctx, "non-existent")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetSkillsForEvent", func() {
		It("should retrieve skills associated with an event", func() {
			// This will be tested after we implement event-skill associations
			Skip("Requires event-skill associations - will be implemented in Phase 2")
		})
	})

	Describe("GetEventCountsForSkills", func() {
		It("should return event counts for all skills", func() {
			// This will be tested after we implement event-skill associations
			Skip("Requires event-skill associations - will be implemented in Phase 2")
		})
	})
})
