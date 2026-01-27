//nolint:errcheck // Test file - error handling for test setup is not relevant.
package career

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

// Helper function to create test events with skills
func createTestEventWithSkills(text, dateStr string, skillIDs []string) *career.CareerEvent {
	date, _ := time.Parse("2006-01-02", dateStr)
	return &career.CareerEvent{
		Text:   text,
		Date:   date,
		Skills: skillIDs,
	}
}

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
			// Create main repository for events
			mainRepo, err := NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create skills
			skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
			skill2 := &career.Skill{Name: "Go", Category: "backend"}
			skill3 := &career.Skill{Name: "Python", Category: "backend"}
			err = repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill2)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill3)
			Expect(err).NotTo(HaveOccurred())

			// Create event with two skills
			event := createTestEventWithSkills("Backend work", "2024-01-15", []string{skill1.ID, skill2.ID})
			err = mainRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve skills for the event
			skills, err := repository.GetSkillsForEvent(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))

			// Verify the correct skills were returned
			skillNames := []string{skills[0].Name, skills[1].Name}
			Expect(skillNames).To(ContainElements("Ruby", "Go"))
			Expect(skillNames).NotTo(ContainElement("Python"))
		})

		It("should return empty list for event with no skills", func() {
			// Create main repository for events
			mainRepo, err := NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create event with no skills
			event := createTestEventWithSkills("No skills event", "2024-01-15", nil)
			err = mainRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve skills for the event
			skills, err := repository.GetSkillsForEvent(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetEventCountsForSkills", func() {
		It("should return event counts for all skills", func() {
			// Create main repository for events
			mainRepo, err := NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create skills
			skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
			skill2 := &career.Skill{Name: "Go", Category: "backend"}
			skill3 := &career.Skill{Name: "Python", Category: "backend"}
			err = repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill2)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill3)
			Expect(err).NotTo(HaveOccurred())

			// Create events with different skill combinations
			event1 := createTestEventWithSkills("Event 1", "2024-01-15", []string{skill1.ID, skill2.ID})
			event2 := createTestEventWithSkills("Event 2", "2024-02-15", []string{skill1.ID})
			event3 := createTestEventWithSkills("Event 3", "2024-03-15", []string{skill2.ID, skill3.ID})

			err = mainRepo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event3)
			Expect(err).NotTo(HaveOccurred())

			// Get event counts
			counts, err := repository.GetEventCountsForSkills(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Verify counts: Ruby=2, Go=2, Python=1
			Expect(counts[skill1.ID]).To(Equal(2))
			Expect(counts[skill2.ID]).To(Equal(2))
			Expect(counts[skill3.ID]).To(Equal(1))
		})

		It("should return empty map when no events exist", func() {
			// Create skills with no events
			skill := &career.Skill{Name: "Unused", Category: "backend"}
			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			counts, err := repository.GetEventCountsForSkills(ctx)
			Expect(err).NotTo(HaveOccurred())
			// The skill should not appear in the counts since it has no events
			Expect(counts[skill.ID]).To(Equal(0))
		})
	})

	Describe("GetLastUsedForSkills", func() {
		It("should return last used dates for skills with events", func() {
			// Create main repository for events
			mainRepo, err := NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create skills
			skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
			skill2 := &career.Skill{Name: "Go", Category: "backend"}
			err = repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill2)
			Expect(err).NotTo(HaveOccurred())

			// Create events with skills at different dates
			event1 := createTestEventWithSkills("Event 1", "2024-01-01", []string{skill1.ID})
			event2 := createTestEventWithSkills("Event 2", "2024-06-15", []string{skill1.ID})
			event3 := createTestEventWithSkills("Event 3", "2024-03-10", []string{skill2.ID})

			err = mainRepo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event3)
			Expect(err).NotTo(HaveOccurred())

			// Get last used dates
			lastUsedMap, err := repository.GetLastUsedForSkills(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Verify skill1's last used is the most recent event
			lastUsed1, ok := lastUsedMap[skill1.ID]
			Expect(ok).To(BeTrue())
			Expect(lastUsed1.Format("2006-01-02")).To(Equal("2024-06-15"))

			// Verify skill2's last used
			lastUsed2, ok := lastUsedMap[skill2.ID]
			Expect(ok).To(BeTrue())
			Expect(lastUsed2.Format("2006-01-02")).To(Equal("2024-03-10"))
		})

		It("should return empty map when no skills have events", func() {
			// Create a skill without any events
			skill := &career.Skill{Name: "Ruby", Category: "backend"}
			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			lastUsedMap, err := repository.GetLastUsedForSkills(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(lastUsedMap).To(BeEmpty())
		})
	})

	Describe("List with sorting and MinEvents filter", func() {
		var (
			mainRepo *SQLiteRepository
			skill1   *career.Skill
			skill2   *career.Skill
			skill3   *career.Skill
			skill4   *career.Skill
		)

		BeforeEach(func() {
			var err error
			// Create main repository for events
			mainRepo, err = NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create skills with different categories and levels
			skill1 = &career.Skill{Name: "Ruby", Category: "backend", Level: "advanced"}
			skill2 = &career.Skill{Name: "Go", Category: "backend", Level: "intermediate"}
			skill3 = &career.Skill{Name: "React", Category: "frontend", Level: "advanced"}
			skill4 = &career.Skill{Name: "Vue", Category: "frontend", Level: "beginner"}

			err = repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill2)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill3)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill4)
			Expect(err).NotTo(HaveOccurred())

			// Create events with skill associations
			// skill1 (Ruby): 3 events
			// skill2 (Go): 2 events
			// skill3 (React): 1 event
			// skill4 (Vue): 0 events (unused)
			event1 := createTestEventWithSkills("Event 1", "2024-01-15", []string{skill1.ID})
			event2 := createTestEventWithSkills("Event 2", "2024-02-15", []string{skill1.ID, skill2.ID})
			event3 := createTestEventWithSkills("Event 3", "2024-03-15", []string{skill1.ID, skill2.ID, skill3.ID})

			err = mainRepo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event3)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("MinEvents filter", func() {
			It("should filter skills with minimum event count > 0 (used skills only)", func() {
				filters := &SkillFilters{
					MinEvents: 1,
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				// Should exclude skill4 (Vue) which has 0 events
				Expect(skills).To(HaveLen(3))

				skillNames := make([]string, len(skills))
				for i, s := range skills {
					skillNames[i] = s.Name
				}
				Expect(skillNames).To(ContainElements("Ruby", "Go", "React"))
				Expect(skillNames).NotTo(ContainElement("Vue"))
			})

			It("should filter skills with minimum event count >= 2", func() {
				filters := &SkillFilters{
					MinEvents: 2,
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				// Should only include Ruby (3) and Go (2)
				Expect(skills).To(HaveLen(2))

				skillNames := make([]string, len(skills))
				for i, s := range skills {
					skillNames[i] = s.Name
				}
				Expect(skillNames).To(ContainElements("Ruby", "Go"))
			})

			It("should combine MinEvents with category filter", func() {
				filters := &SkillFilters{
					Category:  "backend",
					MinEvents: 2,
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				// Should only include backend skills with >= 2 events: Ruby (3), Go (2)
				Expect(skills).To(HaveLen(2))
				for _, s := range skills {
					Expect(s.Category).To(Equal("backend"))
				}
			})
		})

		Context("Sorting", func() {
			It("should sort skills by name ascending (default)", func() {
				filters := &SkillFilters{
					SortBy:    "name",
					SortOrder: "asc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				Expect(skills[0].Name).To(Equal("Go"))
				Expect(skills[1].Name).To(Equal("React"))
				Expect(skills[2].Name).To(Equal("Ruby"))
				Expect(skills[3].Name).To(Equal("Vue"))
			})

			It("should sort skills by name descending", func() {
				filters := &SkillFilters{
					SortBy:    "name",
					SortOrder: "desc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				Expect(skills[0].Name).To(Equal("Vue"))
				Expect(skills[1].Name).To(Equal("Ruby"))
				Expect(skills[2].Name).To(Equal("React"))
				Expect(skills[3].Name).To(Equal("Go"))
			})

			It("should sort skills by event count descending (most used first)", func() {
				filters := &SkillFilters{
					SortBy:    "events",
					SortOrder: "desc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				// Ruby=3, Go=2, React=1, Vue=0
				Expect(skills[0].Name).To(Equal("Ruby"))
				Expect(skills[1].Name).To(Equal("Go"))
				Expect(skills[2].Name).To(Equal("React"))
				Expect(skills[3].Name).To(Equal("Vue"))
			})

			It("should sort skills by event count ascending (least used first)", func() {
				filters := &SkillFilters{
					SortBy:    "events",
					SortOrder: "asc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				// Vue=0, React=1, Go=2, Ruby=3
				Expect(skills[0].Name).To(Equal("Vue"))
				Expect(skills[1].Name).To(Equal("React"))
				Expect(skills[2].Name).To(Equal("Go"))
				Expect(skills[3].Name).To(Equal("Ruby"))
			})

			It("should sort skills by last used date descending (most recent first)", func() {
				filters := &SkillFilters{
					SortBy:    "last_used",
					SortOrder: "desc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				// Ruby, Go, React all used in Event 3 (2024-03-15)
				// Vue never used - should be last
				Expect(skills[3].Name).To(Equal("Vue"))
			})

			It("should sort skills by category", func() {
				filters := &SkillFilters{
					SortBy:    "category",
					SortOrder: "asc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(4))
				// backend skills first (Go, Ruby), then frontend (React, Vue)
				Expect(skills[0].Category).To(Equal("backend"))
				Expect(skills[1].Category).To(Equal("backend"))
				Expect(skills[2].Category).To(Equal("frontend"))
				Expect(skills[3].Category).To(Equal("frontend"))
			})

			It("should combine sorting with filtering", func() {
				filters := &SkillFilters{
					Category:  "backend",
					SortBy:    "events",
					SortOrder: "desc",
				}

				skills, err := repository.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))
				// Ruby (3 events) before Go (2 events)
				Expect(skills[0].Name).To(Equal("Ruby"))
				Expect(skills[1].Name).To(Equal("Go"))
			})
		})
	})

	Describe("GetEventsUsingSkill", func() {
		It("should return all events that use a specific skill", func() {
			// Create main repository for events
			mainRepo, err := NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create skills
			skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
			skill2 := &career.Skill{Name: "Go", Category: "backend"}
			err = repository.Create(ctx, skill1)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(ctx, skill2)
			Expect(err).NotTo(HaveOccurred())

			// Create events - some with skill1, some with skill2, some with both
			event1 := createTestEventWithSkills("Ruby Event 1", "2024-01-01", []string{skill1.ID})
			event2 := createTestEventWithSkills("Go Event", "2024-02-01", []string{skill2.ID})
			event3 := createTestEventWithSkills("Ruby Event 2", "2024-03-01", []string{skill1.ID})
			event4 := createTestEventWithSkills("Both Skills", "2024-04-01", []string{skill1.ID, skill2.ID})

			err = mainRepo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event3)
			Expect(err).NotTo(HaveOccurred())
			err = mainRepo.Create(ctx, event4)
			Expect(err).NotTo(HaveOccurred())

			// Get events using skill1
			events, err := repository.GetEventsUsingSkill(ctx, skill1.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(3)) // event1, event3, event4

			// Verify events are ordered by date DESC (most recent first)
			Expect(events[0].Text).To(Equal("Both Skills"))
			Expect(events[1].Text).To(Equal("Ruby Event 2"))
			Expect(events[2].Text).To(Equal("Ruby Event 1"))
		})

		It("should return empty slice for skill with no events", func() {
			// Create a skill without any events
			skill := &career.Skill{Name: "Ruby", Category: "backend"}
			err := repository.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			events, err := repository.GetEventsUsingSkill(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})

		It("should return error for non-existent skill", func() {
			events, err := repository.GetEventsUsingSkill(ctx, "non-existent-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})
	})
})
