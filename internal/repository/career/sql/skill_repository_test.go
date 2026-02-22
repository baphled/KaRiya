package sql

import (
	"context"
	stdsql "database/sql"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDB() *gorm.DB {
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
			skill := fixtures.SkillWith("", "Go", "backend", "")

			err := repo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).NotTo(BeEmpty())
			Expect(skill.CreatedAt).NotTo(BeZero())
			Expect(skill.UpdatedAt).NotTo(BeZero())
		})

		It("creates a skill with provided ID", func() {
			skill := fixtures.SkillWith("custom-id", "Ruby", "backend", "")

			err := repo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).To(Equal("custom-id"))
		})

		It("fails on duplicate name", func() {
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Go", "frontend", "")

			Expect(repo.Create(ctx, skill1)).To(Succeed())

			err := repo.Create(ctx, skill2)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetByID", func() {
		It("returns the skill", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "expert")
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
			skill := fixtures.SkillWith("", "Go", "backend", "")
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
			skill := fixtures.SkillWith("", "Go", "backend", "beginner")
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
			skill := fixtures.SkillWith("nonexistent", "Go", "backend", "")

			err := repo.Update(ctx, skill)

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the skill", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")
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
				fixtures.SkillWith("", "Go", "backend", "expert"),
				fixtures.SkillWith("", "Ruby", "backend", "intermediate"),
				fixtures.SkillWith("", "React", "frontend", "advanced"),
				fixtures.SkillWith("", "Vue", "frontend", "beginner"),
				fixtures.SkillWith("", "Kubernetes", "devops", "advanced"),
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
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithCategory("backend"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, s := range skills {
				Expect(s.Category).To(Equal("backend"))
			}
		})

		It("filters by level", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithLevel("advanced"))

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
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("name", "desc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Name).To(Equal("Vue"))
			Expect(skills[4].Name).To(Equal("Go"))
		})

		It("sorts by category", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("category", ""))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Category).To(Equal("backend"))
			Expect(skills[len(skills)-1].Category).To(Equal("frontend"))
		})

		It("applies pagination with limit", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithLimit(0, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("applies pagination with offset", func() {
			allSkills, _ := repo.List(ctx, nil)
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithLimit(2, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			Expect(skills[0].ID).To(Equal(allSkills[2].ID))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			skills := []*career.Skill{
				fixtures.SkillWith("", "Go", "backend", ""),
				fixtures.SkillWith("", "Ruby", "backend", ""),
				fixtures.SkillWith("", "React", "frontend", ""),
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
			count, err := repo.Count(ctx, fixtures.SkillListFiltersWithCategory("backend"))

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})

	Describe("Event Links", func() {
		var skill *career.Skill

		BeforeEach(func() {
			skill = fixtures.SkillWith("", "Go", "backend", "")
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
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			for i := range 3 {
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
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithMinEvents(2))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})
	})

	Describe("GetLastUsedForSkills", func() {
		var skill *career.Skill

		BeforeEach(func() {
			skill = fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())
		})

		Context("when event date is in RFC3339 format", func() {
			It("parses the date correctly", func() {
				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-rfc3339", "RFC3339 event", "2024-01-15T10:30:00Z", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-rfc3339")).To(Succeed())

				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveKey(skill.ID))
				Expect(result[skill.ID]).To(BeTemporally("~", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), time.Second))
			})
		})

		Context("when event date is in space-separated datetime format", func() {
			It("parses the date correctly", func() {
				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-space", "Space format event", "2024-06-20 14:00:00", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-space")).To(Succeed())

				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveKey(skill.ID))
				Expect(result[skill.ID]).To(BeTemporally("~", time.Date(2024, 6, 20, 14, 0, 0, 0, time.UTC), time.Second))
			})
		})

		Context("when event date is in date-only format", func() {
			It("parses as midnight UTC", func() {
				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-date", "Date only event", "2024-03-01", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-date")).To(Succeed())

				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveKey(skill.ID))
				Expect(result[skill.ID]).To(BeTemporally("~", time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), time.Second))
			})
		})

		Context("when skill has no events", func() {
			It("does not include the skill in the result", func() {
				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(HaveKey(skill.ID))
			})
		})

		Context("when multiple skills have events", func() {
			It("returns the correct date for each skill", func() {
				skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
				Expect(repo.Create(ctx, skill2)).To(Succeed())

				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-go", "Go event", "2024-01-15T10:30:00Z", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				err = db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-ruby", "Ruby event", "2024-06-20T14:00:00Z", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())

				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-go")).To(Succeed())
				Expect(repo.LinkToEvent(ctx, skill2.ID, "evt-ruby")).To(Succeed())

				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveLen(2))
				Expect(result[skill.ID]).To(BeTemporally("~", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), time.Second))
				Expect(result[skill2.ID]).To(BeTemporally("~", time.Date(2024, 6, 20, 14, 0, 0, 0, time.UTC), time.Second))
			})
		})

		Context("when a skill has multiple events", func() {
			It("returns the most recent date", func() {
				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-old", "Old event", "2023-01-01T00:00:00Z", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				err = db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-new", "New event", "2024-06-15T12:00:00Z", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())

				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-old")).To(Succeed())
				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-new")).To(Succeed())

				result, err := repo.GetLastUsedForSkills(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(result[skill.ID]).To(BeTemporally("~", time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), time.Second))
			})
		})

		Context("when event date is unparseable", func() {
			It("returns an error containing the skill ID", func() {
				err := db.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
					"evt-bad", "Bad date event", "not-a-date", time.Now(), time.Now()).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(repo.LinkToEvent(ctx, skill.ID, "evt-bad")).To(Succeed())

				_, err = repo.GetLastUsedForSkills(ctx)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(skill.ID))
			})
		})
	})
})
