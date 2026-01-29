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
			event := &career.Event{
				Text:    "Implemented feature X",
				Date:    time.Now(),
				Company: "TechCo",
			}

			err := repo.Create(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).NotTo(BeEmpty())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})

		It("creates an event with provided ID", func() {
			event := &career.Event{
				ID:      "custom-id",
				Text:    "Fixed bug Y",
				Date:    time.Now(),
				Company: "TechCo",
			}

			err := repo.Create(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).To(Equal("custom-id"))
		})

		It("stores tags and categories", func() {
			event := &career.Event{
				Text:       "Led project",
				Date:       time.Now(),
				Tags:       []string{"leadership", "project-management"},
				Categories: []string{"management", "technical"},
			}

			Expect(repo.Create(ctx, event)).To(Succeed())

			found, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.Tags).To(ConsistOf("leadership", "project-management"))
			Expect(found.Categories).To(ConsistOf("management", "technical"))
		})

		It("saves skill associations", func() {
			// Create skills first.
			skillRepo := NewSkillRepository(db)
			skill1 := &career.Skill{Name: "Go", Category: "backend"}
			skill2 := &career.Skill{Name: "Docker", Category: "devops"}
			Expect(skillRepo.Create(ctx, skill1)).To(Succeed())
			Expect(skillRepo.Create(ctx, skill2)).To(Succeed())

			event := &career.Event{
				Text:   "Built microservice",
				Date:   time.Now(),
				Skills: []string{skill1.ID, skill2.ID},
			}

			Expect(repo.Create(ctx, event)).To(Succeed())

			found, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.Skills).To(ConsistOf(skill1.ID, skill2.ID))
		})
	})

	Describe("GetByID", func() {
		It("returns the event", func() {
			event := &career.Event{
				Text:    "Did something",
				Date:    time.Now(),
				Company: "Corp",
				Project: "Alpha",
			}
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
			event := &career.Event{Text: "Original", Date: time.Now()}
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
			skill1 := &career.Skill{Name: "Go", Category: "backend"}
			skill2 := &career.Skill{Name: "Python", Category: "backend"}
			Expect(skillRepo.Create(ctx, skill1)).To(Succeed())
			Expect(skillRepo.Create(ctx, skill2)).To(Succeed())

			event := &career.Event{
				Text:   "Work",
				Date:   time.Now(),
				Skills: []string{skill1.ID},
			}
			Expect(repo.Create(ctx, event)).To(Succeed())

			// Update to different skill.
			event.Skills = []string{skill2.ID}
			Expect(repo.Update(ctx, event)).To(Succeed())

			found, _ := repo.GetByID(ctx, event.ID)
			Expect(found.Skills).To(ConsistOf(skill2.ID))
		})

		It("returns ErrEventNotFound for missing event", func() {
			event := &career.Event{ID: "nonexistent", Text: "Test", Date: time.Now()}

			err := repo.Update(ctx, event)

			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the event", func() {
			event := &career.Event{Text: "To delete", Date: time.Now()}
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
			events := []*career.Event{
				{Text: "Event 1", Date: now.AddDate(0, 0, -2), Tags: []string{"tag1"}},
				{Text: "Event 2", Date: now.AddDate(0, 0, -1), Tags: []string{"tag2"}},
				{Text: "Event 3", Date: now, Tags: []string{"tag1", "tag2"}},
			}
			for _, e := range events {
				Expect(repo.Create(ctx, e)).To(Succeed())
			}
		})

		It("returns all events without filters", func() {
			events, err := repo.List(ctx, career_repo.EventListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(3))
		})

		It("filters by tag", func() {
			events, err := repo.List(ctx, career_repo.EventListFilters{Tags: []string{"tag1"}})

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})

		It("filters by date range", func() {
			// Use start of yesterday to ensure Event 2 and Event 3 are included.
			start := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
			events, err := repo.List(ctx, career_repo.EventListFilters{StartDate: &start})

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})

		It("sorts by date ascending", func() {
			events, err := repo.List(ctx, career_repo.EventListFilters{SortBy: "date", SortOrder: "asc"})

			Expect(err).NotTo(HaveOccurred())
			Expect(events[0].Text).To(Equal("Event 1"))
			Expect(events[2].Text).To(Equal("Event 3"))
		})

		It("sorts by date descending", func() {
			events, err := repo.List(ctx, career_repo.EventListFilters{SortBy: "date", SortOrder: "desc"})

			Expect(err).NotTo(HaveOccurred())
			Expect(events[0].Text).To(Equal("Event 3"))
			Expect(events[2].Text).To(Equal("Event 1"))
		})

		It("applies pagination", func() {
			events, err := repo.List(ctx, career_repo.EventListFilters{Limit: 2, Offset: 1})

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			now := time.Now()
			events := []*career.Event{
				{Text: "Event 1", Date: now, Tags: []string{"tag1"}},
				{Text: "Event 2", Date: now, Tags: []string{"tag2"}},
				{Text: "Event 3", Date: now, Tags: []string{"tag1"}},
			}
			for _, e := range events {
				Expect(repo.Create(ctx, e)).To(Succeed())
			}
		})

		It("counts all events", func() {
			count, err := repo.Count(ctx, career_repo.EventListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("counts with tag filter", func() {
			count, err := repo.Count(ctx, career_repo.EventListFilters{Tags: []string{"tag1"}})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})
})
