package sql

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Skill Repository", func() {
	var (
		repo *SkillRepository
		tx   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		Expect(sharedGormDB).NotTo(BeNil(), "shared DB not initialized - BeforeSuite not run")
		tx = sharedGormDB.Begin()
		DeferCleanup(func() { tx.Rollback() })
		repo = NewSkillRepository(tx)
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

		It("handles case-insensitive matching", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			found, err := repo.GetByName(ctx, "go")

			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(skill.ID))
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

		It("counts with MinEvents filter", func() {
			// Link one skill to events
			skill := fixtures.SkillWith("", "GoMinEvents", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			eventRepo := NewEventRepository(tx)
			for range 2 {
				event := fixtures.EventWith("", "Event", "", "")
				Expect(eventRepo.Create(ctx, event)).To(Succeed())
				Expect(repo.LinkToEvent(ctx, skill.ID, event.ID)).To(Succeed())
			}

			count, err := repo.Count(ctx, fixtures.SkillListFiltersWithMinEvents(2))

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("Event Links", func() {
		var skill *career.Skill

		BeforeEach(func() {
			skill = fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			event := &models.Event{ID: "event-1", Text: "Test", Date: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			Expect(tx.Create(event).Error).NotTo(HaveOccurred())
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
				Expect(tx.Create(event).Error).NotTo(HaveOccurred())
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
				err := tx.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
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
				err := tx.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
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
				err := tx.Exec(`INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
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
	})

	Describe("GetByCategory", func() {
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

		It("returns skills in the specified category", func() {
			skills, err := repo.GetByCategory(ctx, "backend")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			for _, s := range skills {
				Expect(s.Category).To(Equal("backend"))
			}
		})

		It("returns empty list for non-existent category", func() {
			skills, err := repo.GetByCategory(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetSkillsForEvent", func() {
		var event *career.Event

		BeforeEach(func() {
			event = fixtures.EventWith("test-event", "Test event", "", "")
			eventRepo := NewEventRepository(tx)
			Expect(eventRepo.Create(ctx, event)).To(Succeed())

			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Python", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill1.ID, event.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, event.ID)).To(Succeed())
		})

		It("returns skills linked to the event", func() {
			skills, err := repo.GetSkillsForEvent(ctx, event.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("returns empty list for event with no skills", func() {
			event2 := fixtures.EventWith("event-no-skills", "No skills event", "", "")
			eventRepo := NewEventRepository(tx)
			Expect(eventRepo.Create(ctx, event2)).To(Succeed())

			skills, err := repo.GetSkillsForEvent(ctx, event2.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetSkillsForEvents", func() {
		var event1, event2 *career.Event

		BeforeEach(func() {
			event1 = fixtures.EventWith("event-1", "Event 1", "", "")
			event2 = fixtures.EventWith("event-2", "Event 2", "", "")
			eventRepo := NewEventRepository(tx)
			Expect(eventRepo.Create(ctx, event1)).To(Succeed())
			Expect(eventRepo.Create(ctx, event2)).To(Succeed())

			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Python", "backend", "")
			skill3 := fixtures.SkillWith("", "React", "frontend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())
			Expect(repo.Create(ctx, skill3)).To(Succeed())

			// skill1 and skill2 linked to event1
			Expect(repo.LinkToEvent(ctx, skill1.ID, event1.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, event1.ID)).To(Succeed())
			// skill1 and skill3 linked to event2
			Expect(repo.LinkToEvent(ctx, skill1.ID, event2.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill3.ID, event2.ID)).To(Succeed())
		})

		It("returns unique skills for multiple events", func() {
			skills, err := repo.GetSkillsForEvents(ctx, []string{event1.ID, event2.ID})

			Expect(err).NotTo(HaveOccurred())
			// Should have 3 unique skills (Go, Python, React)
			Expect(skills).To(HaveLen(3))
		})

		It("returns empty list for empty event IDs", func() {
			skills, err := repo.GetSkillsForEvents(ctx, []string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("returns empty list for non-existent events", func() {
			skills, err := repo.GetSkillsForEvents(ctx, []string{"nonexistent"})

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetEventCountsForSkills", func() {
		var skill1, skill2 *career.Skill

		BeforeEach(func() {
			skill1 = fixtures.SkillWith("", "Go", "backend", "")
			skill2 = fixtures.SkillWith("", "Ruby", "backend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())

			eventRepo := NewEventRepository(tx)
			for range 3 {
				event := fixtures.EventWith("", "Event", "", "")
				Expect(eventRepo.Create(ctx, event)).To(Succeed())
				Expect(repo.LinkToEvent(ctx, skill1.ID, event.ID)).To(Succeed())
			}

			// skill2 has 1 event
			event := fixtures.EventWith("", "Event", "", "")
			Expect(eventRepo.Create(ctx, event)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, event.ID)).To(Succeed())
		})

		It("returns event counts for all skills", func() {
			counts, err := repo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(counts[skill1.ID]).To(Equal(3))
			Expect(counts[skill2.ID]).To(Equal(1))
		})

		It("returns only skills with events in the map", func() {
			// Add a skill with no events
			unusedSkill := fixtures.SkillWith("", "UnusedSkill", "backend", "")
			Expect(repo.Create(ctx, unusedSkill)).To(Succeed())

			counts, err := repo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			// Skills with events should be in the map
			Expect(counts).To(HaveKey(skill1.ID))
			Expect(counts).To(HaveKey(skill2.ID))
			// Unused skill should NOT be in the map
			Expect(counts).NotTo(HaveKey(unusedSkill.ID))
		})
	})

	Describe("GetEventsUsingSkill", func() {
		var skill *career.Skill

		BeforeEach(func() {
			skill = fixtures.SkillWith("", "Go", "backend", "")
			Expect(repo.Create(ctx, skill)).To(Succeed())

			eventRepo := NewEventRepository(tx)
			event1 := fixtures.EventWith("", "Old event", "", "")
			event1.Date = time.Now().AddDate(0, 0, -30)
			event2 := fixtures.EventWith("", "New event", "", "")
			event2.Date = time.Now()
			Expect(eventRepo.Create(ctx, event1)).To(Succeed())
			Expect(eventRepo.Create(ctx, event2)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill.ID, event1.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill.ID, event2.ID)).To(Succeed())
		})

		It("returns events using the skill, ordered by date DESC", func() {
			events, err := repo.GetEventsUsingSkill(ctx, skill.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
			// Should be sorted by date DESC (newest first)
			Expect(events[0].Date).To(BeTemporally(">", events[1].Date))
		})

		It("returns empty list for skill with no events", func() {
			unusedSkill := fixtures.SkillWith("", "Unused", "backend", "")
			Expect(repo.Create(ctx, unusedSkill)).To(Succeed())

			events, err := repo.GetEventsUsingSkill(ctx, unusedSkill.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})
	})

	Describe("Sorting edge cases", func() {
		BeforeEach(func() {
			skill1 := fixtures.SkillWith("", "Go", "backend", "")
			skill2 := fixtures.SkillWith("", "Ruby", "backend", "")
			skill3 := fixtures.SkillWith("", "Python", "frontend", "")
			Expect(repo.Create(ctx, skill1)).To(Succeed())
			Expect(repo.Create(ctx, skill2)).To(Succeed())
			Expect(repo.Create(ctx, skill3)).To(Succeed())

			// Link skill1 to 2 events, skill2 to 1 event
			eventRepo := NewEventRepository(tx)
			event1 := fixtures.EventWith("", "Event 1", "", "")
			event2 := fixtures.EventWith("", "Event 2", "", "")
			Expect(eventRepo.Create(ctx, event1)).To(Succeed())
			Expect(eventRepo.Create(ctx, event2)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill1.ID, event1.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill1.ID, event2.ID)).To(Succeed())
			Expect(repo.LinkToEvent(ctx, skill2.ID, event1.ID)).To(Succeed())
		})

		It("sorts by events count", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("events", "desc"))

			Expect(err).NotTo(HaveOccurred())
			// Go has 2 events, Ruby has 1, Python has 0
			Expect(skills[0].Name).To(Equal("Go"))
		})

		It("sorts by last used date", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("last_used", "desc"))

			Expect(err).NotTo(HaveOccurred())
			// Skills with events should come before those without
			Expect(skills).To(HaveLen(3))
		})

		It("uses default sorting for unknown sort field", func() {
			skills, err := repo.List(ctx, fixtures.SkillListFiltersWithSort("unknown_field", ""))

			Expect(err).NotTo(HaveOccurred())
			// Should fall back to name ASC
			Expect(skills).To(HaveLen(3))
		})
	})
})
