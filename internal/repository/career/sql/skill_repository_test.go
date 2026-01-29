package sql

import (
	"context"
	stdsql "database/sql"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDB() *gorm.DB {
	// Use existing modernc.org/sqlite driver with GORM.
	sqlDB, err := stdsql.Open("sqlite", ":memory:")
	Expect(err).NotTo(HaveOccurred())

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	err = db.AutoMigrate(&models.Skill{}, &models.Event{})
	Expect(err).NotTo(HaveOccurred())

	err = db.Exec(`CREATE TABLE IF NOT EXISTS event_skills (
		event_id TEXT NOT NULL,
		skill_id TEXT NOT NULL,
		PRIMARY KEY (event_id, skill_id)
	)`).Error
	Expect(err).NotTo(HaveOccurred())

	return db
}

var _ = Describe("Skill Repository", func() {
	var (
		repo *SkillRepository
		db   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		db = setupTestDB()
		repo = NewSkillRepository(db)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates a skill with generated ID", func() {
			skill := &career.Skill{
				Name:     "Go",
				Category: "backend",
			}

			err := repo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).NotTo(BeEmpty())
			Expect(skill.CreatedAt).NotTo(BeZero())
			Expect(skill.UpdatedAt).NotTo(BeZero())
		})

		It("creates a skill with provided ID", func() {
			skill := &career.Skill{
				ID:       "custom-id",
				Name:     "Ruby",
				Category: "backend",
			}

			err := repo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).To(Equal("custom-id"))
		})

		It("fails on duplicate name", func() {
			skill1 := &career.Skill{Name: "Go", Category: "backend"}
			skill2 := &career.Skill{Name: "Go", Category: "frontend"}

			Expect(repo.Create(ctx, skill1)).To(Succeed())

			err := repo.Create(ctx, skill2)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetByID", func() {
		It("returns the skill", func() {
			skill := &career.Skill{Name: "Go", Category: "backend", Level: "expert"}
			Expect(repo.Create(ctx, skill)).To(Succeed())

			found, err := repo.GetByID(ctx, skill.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Name).To(Equal("Go"))
			Expect(found.Category).To(Equal("backend"))
			Expect(found.Level).To(Equal("expert"))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			_, err := repo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("GetByName", func() {
		It("returns the skill", func() {
			skill := &career.Skill{Name: "Go", Category: "backend"}
			Expect(repo.Create(ctx, skill)).To(Succeed())

			found, err := repo.GetByName(ctx, "Go")

			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(skill.ID))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			_, err := repo.GetByName(ctx, "Nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("Update", func() {
		It("updates the skill", func() {
			skill := &career.Skill{Name: "Go", Category: "backend", Level: "beginner"}
			Expect(repo.Create(ctx, skill)).To(Succeed())
			originalUpdatedAt := skill.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			skill.Level = "expert"
			err := repo.Update(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))

			found, _ := repo.GetByID(ctx, skill.ID)
			Expect(found.Level).To(Equal("expert"))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			skill := &career.Skill{ID: "nonexistent", Name: "Go", Category: "backend"}

			err := repo.Update(ctx, skill)

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the skill", func() {
			skill := &career.Skill{Name: "Go", Category: "backend"}
			Expect(repo.Create(ctx, skill)).To(Succeed())

			err := repo.Delete(ctx, skill.ID)

			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, skill.ID)
			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			err := repo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			skills := []*career.Skill{
				{Name: "Go", Category: "backend", Level: "expert"},
				{Name: "Ruby", Category: "backend", Level: "intermediate"},
				{Name: "React", Category: "frontend", Level: "advanced"},
				{Name: "Vue", Category: "frontend", Level: "beginner"},
				{Name: "Kubernetes", Category: "devops", Level: "advanced"},
			}
			for _, s := range skills {
				Expect(repo.Create(ctx, s)).To(Succeed())
			}
		})

		It("returns all skills without filters", func() {
			skills, err := repo.List(ctx, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(5))
		})

		It("filters by category", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{Category: "backend"})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, s := range skills {
				Expect(s.Category).To(Equal("backend"))
			}
		})

		It("filters by level", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{Level: "advanced"})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, s := range skills {
				Expect(s.Level).To(Equal("advanced"))
			}
		})

		It("sorts by name ascending (default)", func() {
			skills, err := repo.List(ctx, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Name).To(Equal("Go"))
			Expect(skills[4].Name).To(Equal("Vue"))
		})

		It("sorts by name descending", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{SortBy: "name", SortOrder: "desc"})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Name).To(Equal("Vue"))
			Expect(skills[4].Name).To(Equal("Go"))
		})

		It("sorts by category", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{SortBy: "category"})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Category).To(Equal("backend"))
			Expect(skills[len(skills)-1].Category).To(Equal("frontend"))
		})

		It("applies pagination with limit", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{Limit: 2})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("applies pagination with offset", func() {
			allSkills, _ := repo.List(ctx, nil)
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{Limit: 2, Offset: 2})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			Expect(skills[0].ID).To(Equal(allSkills[2].ID))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			skills := []*career.Skill{
				{Name: "Go", Category: "backend"},
				{Name: "Ruby", Category: "backend"},
				{Name: "React", Category: "frontend"},
			}
			for _, s := range skills {
				Expect(repo.Create(ctx, s)).To(Succeed())
			}
		})

		It("counts all skills", func() {
			count, err := repo.Count(ctx, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("counts with filters", func() {
			count, err := repo.Count(ctx, &career_repo.SkillListFilters{Category: "backend"})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})

	Describe("Event Links", func() {
		var skill *career.Skill

		BeforeEach(func() {
			skill = &career.Skill{Name: "Go", Category: "backend"}
			Expect(repo.Create(ctx, skill)).To(Succeed())

			event := &models.Event{ID: "event-1", Text: "Test", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			Expect(db.Create(event).Error).NotTo(HaveOccurred())
		})

		It("links skill to event", func() {
			err := repo.LinkToEvent(ctx, skill.ID, "event-1")
			Expect(err).NotTo(HaveOccurred())

			ids, err := repo.GetEventIDs(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(ContainElement("event-1"))
		})

		It("unlinks skill from event", func() {
			Expect(repo.LinkToEvent(ctx, skill.ID, "event-1")).To(Succeed())

			err := repo.UnlinkFromEvent(ctx, skill.ID, "event-1")
			Expect(err).NotTo(HaveOccurred())

			ids, _ := repo.GetEventIDs(ctx, skill.ID)
			Expect(ids).NotTo(ContainElement("event-1"))
		})

		It("ignores duplicate links", func() {
			Expect(repo.LinkToEvent(ctx, skill.ID, "event-1")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill.ID, "event-1")).To(Succeed())

			ids, _ := repo.GetEventIDs(ctx, skill.ID)
			Expect(ids).To(HaveLen(1))
		})
	})

	Describe("Filter by MinEvents", func() {
		BeforeEach(func() {
			skill1 := &career.Skill{Name: "Go", Category: "backend"}
			skill2 := &career.Skill{Name: "Ruby", Category: "backend"}
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			for i := 0; i < 3; i++ {
				event := &models.Event{
					ID:        "event-" + string(rune('a'+i)),
					Text:      "Test",
					Date:      time.Now(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				Expect(db.Create(event).Error).NotTo(HaveOccurred())
				Expect(repo.LinkToEvent(ctx, skill1.ID, event.ID)).To(Succeed())
			}
		})

		It("filters by minimum event count", func() {
			skills, err := repo.List(ctx, &career_repo.SkillListFilters{MinEvents: 2})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})
	})
})
