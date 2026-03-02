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

func setupEventTestDB() *gorm.DB {
	sqlDB, err := stdsql.Open("sqlite", ":memory:")
	Expect(err).NotTo(HaveOccurred())

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	err = db.AutoMigrate(&models.Event{}, &models.Skill{})
	Expect(err).NotTo(HaveOccurred())

	err = db.Exec(`CREATE TABLE IF NOT EXISTS event_skills (
		event_id TEXT NOT NULL,
		skill_id TEXT NOT NULL,
		PRIMARY KEY (event_id, skill_id)
	)`).Error
	Expect(err).NotTo(HaveOccurred())

	return db
}

var _ = Describe("Event Repository", func() {
	var (
		repo *EventRepository
		db   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		db = setupEventTestDB()
		repo = NewEventRepository(db)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates an event with generated ID", func() {
			event := fixtures.EventWith("", "Implemented feature X", "TechCo", "")

			err := repo.Create(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).NotTo(BeEmpty())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})

		It("creates an event with provided ID", func() {
			event := fixtures.EventWith("custom-id", "Fixed bug Y", "TechCo", "")

			err := repo.Create(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).To(Equal("custom-id"))
		})

		It("stores tags and categories", func() {
			event := fixtures.EventWith("", "Led project", "", "")
			event.Tags = []string{"leadership", "project-management"}
			event.Categories = []string{"management", "technical"}

			Expect(repo.Create(ctx, event)).To(Succeed())

			found, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.Tags).To(ConsistOf("leadership", "project-management"))
			Expect(found.Categories).To(ConsistOf("management", "technical"))
		})

		It("saves skill associations", func() {
			skillRepo := NewSkillRepository(db)
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Docker", "devops", "")
			Expect(skillRepo.Create(ctx, skill1)).To(Succeed())
			Expect(skillRepo.Create(ctx, skill2)).To(Succeed())

			event := fixtures.EventWith("", "Built microservice", "", "")
			event.Skills = []string{skill1.ID, skill2.ID}

			Expect(repo.Create(ctx, event)).To(Succeed())

			found, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.Skills).To(ConsistOf(skill1.ID, skill2.ID))
		})
	})

	Describe("GetByID", func() {
		It("returns the event", func() {
			event := fixtures.EventWith("", "Did something", "Corp", "Alpha")
			Expect(repo.Create(ctx, event)).To(Succeed())

			found, err := repo.GetByID(ctx, event.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Text).To(Equal("Did something"))
			Expect(found.Company).To(Equal("Corp"))
			Expect(found.Project).To(Equal("Alpha"))
		})

		It("returns ErrEventNotFound for missing event", func() {
			_, err := repo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("Update", func() {
		It("updates the event", func() {
			event := fixtures.EventWith("", "Original", "", "")
			Expect(repo.Create(ctx, event)).To(Succeed())
			originalUpdatedAt := event.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			event.Text = "Updated"
			err := repo.Update(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(event.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))

			found, _ := repo.GetByID(ctx, event.ID)
			Expect(found.Text).To(Equal("Updated"))
		})

		It("updates skill associations", func() {
			skillRepo := NewSkillRepository(db)
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Python", "backend", "")
			Expect(skillRepo.Create(ctx, skill1)).To(Succeed())
			Expect(skillRepo.Create(ctx, skill2)).To(Succeed())

			event := fixtures.EventWith("", "Work", "", "")
			event.Skills = []string{skill1.ID}
			Expect(repo.Create(ctx, event)).To(Succeed())

			event.Skills = []string{skill2.ID}
			Expect(repo.Update(ctx, event)).To(Succeed())

			found, _ := repo.GetByID(ctx, event.ID)
			Expect(found.Skills).To(ConsistOf(skill2.ID))
		})

		It("returns ErrEventNotFound for missing event", func() {
			event := fixtures.EventWith("nonexistent", "Test", "", "")

			err := repo.Update(ctx, event)

			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the event", func() {
			event := fixtures.EventWith("", "To delete", "", "")
			Expect(repo.Create(ctx, event)).To(Succeed())

			err := repo.Delete(ctx, event.ID)

			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, event.ID)
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})

		It("returns ErrEventNotFound for missing event", func() {
			err := repo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			now := time.Now()
			e1 := fixtures.EventWith("", "Event 1", "", "")
			e1.Date = now.AddDate(0, 0, -2)
			e1.Tags = []string{"tag1"}
			e2 := fixtures.EventWith("", "Event 2", "", "")
			e2.Date = now.AddDate(0, 0, -1)
			e2.Tags = []string{"tag2"}
			e3 := fixtures.EventWith("", "Event 3", "", "")
			e3.Date = now
			e3.Tags = []string{"tag1", "tag2"}
			for _, e := range []*career.Event{e1, e2, e3} {
				Expect(repo.Create(ctx, e)).To(Succeed())
			}
		})

		It("returns all events without filters", func() {
			events, err := repo.List(ctx, *fixtures.EventListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(3))
		})

		It("filters by tag", func() {
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithTags([]string{"tag1"}))

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})

		It("filters by date range", func() {
			start := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithDateRange(&start, nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})

		It("sorts by date ascending", func() {
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithSort("date", "asc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(events[0].Text).To(Equal("Event 1"))
			Expect(events[2].Text).To(Equal("Event 3"))
		})

		It("sorts by date descending", func() {
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithSort("date", "desc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(events[0].Text).To(Equal("Event 3"))
			Expect(events[2].Text).To(Equal("Event 1"))
		})

		It("applies pagination", func() {
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithLimit(1, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			now := time.Now()
			e1 := fixtures.EventWith("", "Event 1", "", "")
			e1.Date = now
			e1.Tags = []string{"tag1"}
			e2 := fixtures.EventWith("", "Event 2", "", "")
			e2.Date = now
			e2.Tags = []string{"tag2"}
			e3 := fixtures.EventWith("", "Event 3", "", "")
			e3.Date = now
			e3.Tags = []string{"tag1"}
			for _, e := range []*career.Event{e1, e2, e3} {
				Expect(repo.Create(ctx, e)).To(Succeed())
			}
		})

		It("counts all events", func() {
			count, err := repo.Count(ctx, *fixtures.EventListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("counts with tag filter", func() {
			count, err := repo.Count(ctx, *fixtures.EventListFiltersWithTags([]string{"tag1"}))

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})

	Describe("LinkSkill", func() {
		It("should link a skill to an event", func() {
			event := fixtures.EventWith("event-1", "Test event", "", "")
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Link a skill to the event
			err = repo.LinkSkill(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			// Verify the skill was linked
			retrieved, err := repo.GetByID(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Skills).To(ContainElement("skill-1"))
		})

		It("should return error if event does not exist", func() {
			err := repo.LinkSkill(ctx, "nonexistent", "skill-1")
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})

		It("should not duplicate skill if already linked", func() {
			event := fixtures.EventWith("event-1", "Test event", "", "")
			event.Skills = []string{"skill-1"}
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Try to link the same skill again
			err = repo.LinkSkill(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			// Verify skill appears only once
			retrieved, err := repo.GetByID(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Skills).To(HaveLen(1))
			Expect(retrieved.Skills).To(ContainElement("skill-1"))
		})
	})

	Describe("UnlinkSkill", func() {
		It("should unlink a skill from an event", func() {
			event := fixtures.EventWith("event-1", "Test event", "", "")
			event.Skills = []string{"skill-1", "skill-2"}
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			err = repo.UnlinkSkill(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			retrieved, err := repo.GetByID(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Skills).ToNot(ContainElement("skill-1"))
			Expect(retrieved.Skills).To(ContainElement("skill-2"))
		})

		It("should return error if event does not exist", func() {
			err := repo.UnlinkSkill(ctx, "nonexistent", "skill-1")
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})

		It("should handle unlinking non-linked skill gracefully", func() {
			event := fixtures.EventWith("event-1", "Test event", "", "")
			event.Skills = []string{"skill-2"}
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			err = repo.UnlinkSkill(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			retrieved, err := repo.GetByID(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Skills).To(HaveLen(1))
			Expect(retrieved.Skills).To(ContainElement("skill-2"))
		})
	})

	Describe("List with skill associations", func() {
		It("loads skill IDs for listed events", func() {
			event := fixtures.EventWith("", "Built microservice", "", "")
			event.Skills = []string{"skill-1", "skill-2"}
			Expect(repo.Create(ctx, event)).To(Succeed())

			events, err := repo.List(ctx, *fixtures.EventListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Skills).To(ConsistOf("skill-1", "skill-2"))
		})

		It("handles events with no skills in list", func() {
			event := fixtures.EventWith("", "No skills event", "", "")
			Expect(repo.Create(ctx, event)).To(Succeed())

			events, err := repo.List(ctx, *fixtures.EventListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Skills).To(BeEmpty())
		})
	})

	Describe("Create with empty skills", func() {
		It("creates event with no skill associations", func() {
			event := fixtures.EventWith("", "No skills", "", "")
			event.Skills = []string{}

			err := repo.Create(ctx, event)

			Expect(err).NotTo(HaveOccurred())

			found, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.Skills).To(BeEmpty())
		})
	})

	Describe("End date filter", func() {
		It("filters by end date", func() {
			now := time.Now()
			e1 := fixtures.EventWith("", "Old event", "", "")
			e1.Date = now.Add(-72 * time.Hour)
			e2 := fixtures.EventWith("", "New event", "", "")
			e2.Date = now
			Expect(repo.Create(ctx, e1)).To(Succeed())
			Expect(repo.Create(ctx, e2)).To(Succeed())

			end := now.Add(-24 * time.Hour)
			events, err := repo.List(ctx, *fixtures.EventListFiltersWithDateRange(nil, &end))

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Text).To(Equal("Old event"))
		})
	})

	Describe("Error handling with closed DB", func() {
		var closedRepo *EventRepository

		BeforeEach(func() {
			closedDB, err := stdsql.Open("sqlite", ":memory:")
			Expect(err).NotTo(HaveOccurred())
			gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: closedDB}), &gorm.Config{})
			Expect(err).NotTo(HaveOccurred())
			closedDB.Close()
			closedRepo = NewEventRepository(gormDB)
		})

		It("returns error on Create with closed DB", func() {
			event := fixtures.EventWith("", "Test", "", "")
			err := closedRepo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
		})

		It("returns error on GetByID with closed DB", func() {
			_, err := closedRepo.GetByID(ctx, "id")
			Expect(err).To(HaveOccurred())
		})

		It("returns error on Update with closed DB", func() {
			event := fixtures.EventWith("id", "Test", "", "")
			err := closedRepo.Update(ctx, event)
			Expect(err).To(HaveOccurred())
		})

		It("returns error on List with closed DB", func() {
			_, err := closedRepo.List(ctx, *fixtures.EventListFilters())
			Expect(err).To(HaveOccurred())
		})

		It("returns error on Count with closed DB", func() {
			_, err := closedRepo.Count(ctx, *fixtures.EventListFilters())
			Expect(err).To(HaveOccurred())
		})

		It("returns error on LinkSkill with closed DB", func() {
			err := closedRepo.LinkSkill(ctx, "e1", "s1")
			Expect(err).To(HaveOccurred())
		})

		It("returns error on UnlinkSkill with closed DB", func() {
			err := closedRepo.UnlinkSkill(ctx, "e1", "s1")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("saveSkillAssociations error path", func() {
		It("returns error when event_skills table is missing", func() {
			Expect(db.Exec("DROP TABLE event_skills").Error).NotTo(HaveOccurred())

			event := fixtures.EventWith("", "Test event", "", "")
			event.Skills = []string{"skill-1"}

			err := repo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
		})
	})
})
