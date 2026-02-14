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

	Describe("GetByCategory", func() {
		BeforeEach(func() {
			Expect(repo.Create(ctx, fixtures.SkillWith("", "Go", "backend", ""))).To(Succeed())
			Expect(repo.Create(ctx, fixtures.SkillWith("", "Ruby", "backend", ""))).To(Succeed())
			Expect(repo.Create(ctx, fixtures.SkillWith("", "React", "frontend", ""))).To(Succeed())
		})

		It("returns skills in the specified category", func() {
			skills, err := repo.GetByCategory(ctx, "backend")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, s := range skills {
				Expect(s.Category).To(Equal("backend"))
			}
		})

		It("returns empty for unknown category", func() {
			skills, err := repo.GetByCategory(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetSkillsForEvent", func() {
		It("returns skills linked to an event", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			event := &models.Event{ID: "e1", Text: "Test", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			Expect(db.Create(event).Error).NotTo(HaveOccurred())
			Expect(repo.LinkToEvent(ctx, skill.ID, "e1")).To(Succeed())

			skills, err := repo.GetSkillsForEvent(ctx, "e1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})

		It("returns empty for event with no skills", func() {
			skills, err := repo.GetSkillsForEvent(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetEventCountsForSkills", func() {
		It("returns event counts per skill", func() {
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			now := time.Now()
			for _, eid := range []string{"e1", "e2", "e3"} {
				event := &models.Event{ID: eid, Text: "Test", Date: now, CreatedAt: now, UpdatedAt: now}
				Expect(db.Create(event).Error).NotTo(HaveOccurred())
			}
			Expect(repo.LinkToEvent(ctx, skill1.ID, "e1")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill1.ID, "e2")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, "e3")).To(Succeed())

			counts, err := repo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(counts[skill1.ID]).To(Equal(2))
			Expect(counts[skill2.ID]).To(Equal(1))
		})

		It("returns empty map when no links exist", func() {
			counts, err := repo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(counts).To(BeEmpty())
		})
	})

	Describe("GetLastUsedForSkills", func() {
		It("returns empty map when no links exist", func() {
			lastUsed, err := repo.GetLastUsedForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(lastUsed).To(BeEmpty())
		})

		It("queries the database for last used dates", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			now := time.Now()
			e1 := &models.Event{ID: "e1", Text: "Test", Date: now, CreatedAt: now, UpdatedAt: now}
			Expect(db.Create(e1).Error).NotTo(HaveOccurred())
			Expect(repo.LinkToEvent(ctx, skill.ID, "e1")).To(Succeed())

			_, err := repo.GetLastUsedForSkills(ctx)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetEventsUsingSkill", func() {
		It("returns events ordered by date descending", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			now := time.Now()
			earlier := now.Add(-48 * time.Hour)
			e1 := &models.Event{ID: "e1", Text: "Old", Date: earlier, CreatedAt: now, UpdatedAt: now}
			e2 := &models.Event{ID: "e2", Text: "New", Date: now, CreatedAt: now, UpdatedAt: now}
			Expect(db.Create(e1).Error).NotTo(HaveOccurred())
			Expect(db.Create(e2).Error).NotTo(HaveOccurred())
			Expect(repo.LinkToEvent(ctx, skill.ID, "e1")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill.ID, "e2")).To(Succeed())

			events, err := repo.GetEventsUsingSkill(ctx, skill.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
			Expect(events[0].ID).To(Equal("e2"))
			Expect(events[1].ID).To(Equal("e1"))
		})

		It("returns empty for skill with no events", func() {
			events, err := repo.GetEventsUsingSkill(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})
	})

	Describe("Sort by events and last_used", func() {
		It("sorts by events ascending", func() {
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			now := time.Now()
			for _, eid := range []string{"e1", "e2"} {
				event := &models.Event{ID: eid, Text: "Test", Date: now, CreatedAt: now, UpdatedAt: now}
				Expect(db.Create(event).Error).NotTo(HaveOccurred())
			}
			Expect(repo.LinkToEvent(ctx, skill1.ID, "e1")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill1.ID, "e2")).To(Succeed())

			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("events", "desc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Name).To(Equal("Go"))
		})

		It("sorts by last_used descending", func() {
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			now := time.Now()
			earlier := now.Add(-48 * time.Hour)
			e1 := &models.Event{ID: "e1", Text: "Old", Date: earlier, CreatedAt: now, UpdatedAt: now}
			e2 := &models.Event{ID: "e2", Text: "New", Date: now, CreatedAt: now, UpdatedAt: now}
			Expect(db.Create(e1).Error).NotTo(HaveOccurred())
			Expect(db.Create(e2).Error).NotTo(HaveOccurred())
			Expect(repo.LinkToEvent(ctx, skill1.ID, "e1")).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, "e2")).To(Succeed())

			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("last_used", "desc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills[0].Name).To(Equal("Ruby"))
		})
	})
})

var _ = Describe("Skill Repository Error Handling", func() {
	var (
		closedRepo *SkillRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		closedDB, err := stdsql.Open("sqlite", ":memory:")
		Expect(err).NotTo(HaveOccurred())
		gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: closedDB}), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		closedDB.Close()
		closedRepo = NewSkillRepository(gormDB)
		ctx = context.Background()
	})

	It("returns error on Create with closed DB", func() {
		skill := fixtures.SkillWith("", "Go", "backend", "")
		err := closedRepo.Create(ctx, skill)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetByID with closed DB", func() {
		_, err := closedRepo.GetByID(ctx, "id")
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetByName with closed DB", func() {
		_, err := closedRepo.GetByName(ctx, "Go")
		Expect(err).To(HaveOccurred())
	})

	It("returns error on Update with closed DB", func() {
		skill := fixtures.SkillWith("id", "Go", "backend", "")
		err := closedRepo.Update(ctx, skill)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on List with closed DB", func() {
		_, err := closedRepo.List(ctx, nil)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on Count with closed DB", func() {
		_, err := closedRepo.Count(ctx, nil)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetSkillsForEvent with closed DB", func() {
		_, err := closedRepo.GetSkillsForEvent(ctx, "e1")
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetEventCountsForSkills with closed DB", func() {
		_, err := closedRepo.GetEventCountsForSkills(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetLastUsedForSkills with closed DB", func() {
		_, err := closedRepo.GetLastUsedForSkills(ctx)
		Expect(err).To(HaveOccurred())
	})

	It("returns error on GetEventsUsingSkill with closed DB", func() {
		_, err := closedRepo.GetEventsUsingSkill(ctx, "s1")
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("SQL Repositories", func() {
	It("creates repositories from GORM DB", func() {
		db := setupTestDB()
		repos := NewRepositoriesFromDB(db)

		Expect(repos.Event).NotTo(BeNil())
		Expect(repos.Skill).NotTo(BeNil())
		Expect(repos.Fact).NotTo(BeNil())
		Expect(repos.Burst).NotTo(BeNil())
	})

	It("creates repositories from sql.DB via NewRepositories", func() {
		sqlDB, err := stdsql.Open("sqlite", ":memory:")
		Expect(err).NotTo(HaveOccurred())

		repos, err := NewRepositories(sqlDB)

		Expect(err).NotTo(HaveOccurred())
		Expect(repos).NotTo(BeNil())
		Expect(repos.Event).NotTo(BeNil())
	})

	It("creates GORM DB from sql.DB via NewGormDB", func() {
		sqlDB, err := stdsql.Open("sqlite", ":memory:")
		Expect(err).NotTo(HaveOccurred())

		gormDB, err := NewGormDB(sqlDB)

		Expect(err).NotTo(HaveOccurred())
		Expect(gormDB).NotTo(BeNil())
	})

	It("opens a database via OpenDB", func() {
		tmpDB := GinkgoT().TempDir() + "/test.db"
		sqlDB, err := OpenDB(tmpDB)

		Expect(err).NotTo(HaveOccurred())
		Expect(sqlDB).NotTo(BeNil())
		sqlDB.Close()
	})

	It("creates repositories from path via NewRepositoriesFromPath", func() {
		tmpDB := GinkgoT().TempDir() + "/test.db"
		repos, err := NewRepositoriesFromPath(tmpDB)

		Expect(err).NotTo(HaveOccurred())
		Expect(repos).NotTo(BeNil())
		Expect(repos.Event).NotTo(BeNil())
	})
})
