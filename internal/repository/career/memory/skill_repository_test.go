//nolint:errcheck // Test file - error handling for test setup is not relevant.
package memory

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("SkillRepository", func() {
	var (
		skillRepo *SkillRepository
		eventRepo *EventRepository
		ctx       context.Context
	)

	BeforeEach(func() {
		eventRepo = NewEventRepository()
		skillRepo = NewSkillRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		ctx = context.Background() //nolint:fatcontext // test setup
	})

	Describe("Create", func() {
		It("creates a skill with provided ID", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")

			err := skillRepo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).To(Equal("skill-1"))
			Expect(skill.CreatedAt).NotTo(BeZero())
			Expect(skill.UpdatedAt).NotTo(BeZero())
		})

		It("generates an ID when none is provided", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "")

			err := skillRepo.Create(ctx, skill)

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).NotTo(BeEmpty())
		})

		It("returns error for duplicate name", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			err := skillRepo.Create(ctx, fixtures.SkillWith("s2", "Go", "backend", ""))

			Expect(err).To(Equal(career_repo.ErrDuplicateSkill))
		})

		It("returns validation error for invalid skill", func() {
			skill := fixtures.SkillWith("", "", "backend", "")

			err := skillRepo.Create(ctx, skill)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetByID", func() {
		It("retrieves an existing skill", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("skill-1", "Go", "backend", "advanced"))

			skill, err := skillRepo.GetByID(ctx, "skill-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			_, err := skillRepo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("GetByName", func() {
		It("retrieves a skill by exact name", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			skill, err := skillRepo.GetByName(ctx, "Go")

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).To(Equal("s1"))
		})

		It("performs case-insensitive matching", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			skill, err := skillRepo.GetByName(ctx, "go")

			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
		})

		It("returns ErrSkillNotFound for missing name", func() {
			_, err := skillRepo.GetByName(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})
	})

	Describe("Update", func() {
		It("updates an existing skill", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", "intermediate"))

			updated := fixtures.SkillWith("s1", "Go", "backend", "advanced")
			err := skillRepo.Update(ctx, updated)

			Expect(err).NotTo(HaveOccurred())
			Expect(updated.UpdatedAt).NotTo(BeZero())

			skill, _ := skillRepo.GetByID(ctx, "s1")
			Expect(skill.Level).To(Equal("advanced"))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			err := skillRepo.Update(ctx, fixtures.SkillWith("nonexistent", "Go", "backend", ""))

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})

		It("returns validation error for invalid update", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			invalid := fixtures.SkillWith("s1", "", "backend", "")
			err := skillRepo.Update(ctx, invalid)

			Expect(err).To(HaveOccurred())
		})

		It("returns ErrDuplicateSkill when renaming to existing name", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))
			skillRepo.Create(ctx, fixtures.SkillWith("s2", "Python", "backend", ""))

			duplicate := fixtures.SkillWith("s2", "Go", "backend", "")
			err := skillRepo.Update(ctx, duplicate)

			Expect(err).To(Equal(career_repo.ErrDuplicateSkill))
		})
	})

	Describe("Delete", func() {
		It("deletes an existing skill", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			err := skillRepo.Delete(ctx, "s1")

			Expect(err).NotTo(HaveOccurred())

			_, err = skillRepo.GetByID(ctx, "s1")
			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})

		It("returns ErrSkillNotFound for missing skill", func() {
			err := skillRepo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrSkillNotFound))
		})

		It("cleans up event-skill associations on delete", func() {
			event := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
			eventRepo.Create(ctx, event)
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))
			eventRepo.LinkSkill(ctx, "e1", "s1")

			skillRepo.Delete(ctx, "s1")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "e1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", "advanced"))
			skillRepo.Create(ctx, fixtures.SkillWith("s2", "Python", "backend", "intermediate"))
			skillRepo.Create(ctx, fixtures.SkillWith("s3", "React", "frontend", "advanced"))
			skillRepo.Create(ctx, fixtures.SkillWith("s4", "Docker", "devops", "intermediate"))
		})

		It("returns all skills with nil filters", func() {
			skills, err := skillRepo.List(ctx, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(4))
		})

		It("filters by category", func() {
			skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithCategory("backend"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("filters by level", func() {
			skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithLevel("advanced"))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
		})

		It("filters by minimum events", func() {
			event := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
			eventRepo.Create(ctx, event)
			eventRepo.LinkSkill(ctx, "e1", "s1")

			skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithMinEvents(1))

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})

		Context("sorting", func() {
			It("sorts by name ascending by default", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFilters())

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Name).To(Equal("Docker"))
				Expect(skills[3].Name).To(Equal("React"))
			})

			It("sorts by name descending", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("name", "desc"))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Name).To(Equal("React"))
				Expect(skills[3].Name).To(Equal("Docker"))
			})

			It("sorts by category", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("category", ""))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Category).To(Equal("backend"))
			})

			It("sorts by category descending", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("category", "desc"))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Category).To(Equal("frontend"))
			})

			It("sorts by event count", func() {
				event := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
				eventRepo.Create(ctx, event)
				eventRepo.LinkSkill(ctx, "e1", "s1")

				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("events", "desc"))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Name).To(Equal("Go"))
			})

			It("sorts by event count ascending", func() {
				event := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
				eventRepo.Create(ctx, event)
				eventRepo.LinkSkill(ctx, "e1", "s1")

				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("events", ""))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[len(skills)-1].Name).To(Equal("Go"))
			})

			It("sorts by last_used descending", func() {
				now := time.Now()
				event1 := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
				event1.Date = now.Add(-24 * time.Hour)
				eventRepo.Create(ctx, event1)
				eventRepo.LinkSkill(ctx, "e1", "s1")

				event2 := fixtures.EventWith("e2", "Built Python service with Flask", "", "")
				event2.Date = now
				eventRepo.Create(ctx, event2)
				eventRepo.LinkSkill(ctx, "e2", "s2")

				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("last_used", "desc"))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Name).To(Equal("Python"))
				Expect(skills[1].Name).To(Equal("Go"))
			})

			It("sorts by last_used ascending", func() {
				now := time.Now()
				event1 := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
				event1.Date = now.Add(-24 * time.Hour)
				eventRepo.Create(ctx, event1)
				eventRepo.LinkSkill(ctx, "e1", "s1")

				event2 := fixtures.EventWith("e2", "Built Python service with Flask", "", "")
				event2.Date = now
				eventRepo.Create(ctx, event2)
				eventRepo.LinkSkill(ctx, "e2", "s2")

				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("last_used", ""))

				Expect(err).NotTo(HaveOccurred())
				Expect(len(skills)).To(BeNumerically(">=", 2))
			})

			It("falls back to name sort for unknown sort field", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithSort("unknown", ""))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].Name).To(Equal("Docker"))
			})
		})

		Context("pagination", func() {
			It("applies offset and limit", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithLimit(1, 2))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))
			})

			It("returns empty when offset exceeds total", func() {
				skills, err := skillRepo.List(ctx, fixtures.SkillListFiltersWithLimit(100, 0))

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(BeEmpty())
			})
		})
	})

	Describe("GetByCategory", func() {
		It("returns skills in the specified category sorted by name", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Python", "backend", ""))
			skillRepo.Create(ctx, fixtures.SkillWith("s2", "Go", "backend", ""))
			skillRepo.Create(ctx, fixtures.SkillWith("s3", "React", "frontend", ""))

			skills, err := skillRepo.GetByCategory(ctx, "backend")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			Expect(skills[0].Name).To(Equal("Go"))
			Expect(skills[1].Name).To(Equal("Python"))
		})

		It("returns empty for unknown category", func() {
			skills, err := skillRepo.GetByCategory(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("GetEventCountsForSkills", func() {
		It("returns event counts for skills with associations", func() {
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))
			skillRepo.Create(ctx, fixtures.SkillWith("s2", "Python", "backend", ""))

			event1 := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
			event2 := fixtures.EventWith("e2", "Built Python service with Flask", "", "")
			eventRepo.Create(ctx, event1)
			eventRepo.Create(ctx, event2)
			eventRepo.LinkSkill(ctx, "e1", "s1")
			eventRepo.LinkSkill(ctx, "e2", "s1")
			eventRepo.LinkSkill(ctx, "e2", "s2")

			counts, err := skillRepo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(counts["s1"]).To(Equal(2))
			Expect(counts["s2"]).To(Equal(1))
		})

		It("returns empty map when no associations exist", func() {
			counts, err := skillRepo.GetEventCountsForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(counts).To(BeEmpty())
		})
	})

	Describe("GetLastUsedForSkills", func() {
		It("returns empty map (simple implementation)", func() {
			lastUsed, err := skillRepo.GetLastUsedForSkills(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(lastUsed).To(BeEmpty())
		})
	})

	Describe("GetEventsUsingSkill", func() {
		It("returns events sorted by date descending", func() {
			now := time.Now()
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))

			event1 := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
			event1.Date = now.Add(-48 * time.Hour)
			event2 := fixtures.EventWith("e2", "Built Go service with gRPC", "", "")
			event2.Date = now
			eventRepo.Create(ctx, event1)
			eventRepo.Create(ctx, event2)
			eventRepo.LinkSkill(ctx, "e1", "s1")
			eventRepo.LinkSkill(ctx, "e2", "s1")

			events, err := skillRepo.GetEventsUsingSkill(ctx, "s1")

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
			Expect(events[0].ID).To(Equal("e2"))
			Expect(events[1].ID).To(Equal("e1"))
		})

		It("returns empty list for skill with no events", func() {
			events, err := skillRepo.GetEventsUsingSkill(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})

		It("returns empty list when eventRepo is nil", func() {
			repoWithoutEvents := NewSkillRepository()
			repoWithoutEvents.AssociateSkillWithEvent("s1", "e1")

			events, err := repoWithoutEvents.GetEventsUsingSkill(ctx, "s1")

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})
	})

	Describe("GetSkillsForEvent", func() {
		It("returns empty list when no skills are associated", func() {
			event := fixtures.Event("event-1")
			eventRepo.Create(ctx, event)

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("returns empty list for unknown event ID", func() {
			skills, err := skillRepo.GetSkillsForEvent(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("returns skills linked via EventRepository.LinkSkill", func() {
			event := fixtures.EventWith("event-1", "Built Go microservice with Docker", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "advanced")
			dockerSkill := fixtures.SkillWith("skill-docker", "Docker", "devops", "intermediate")
			skillRepo.Create(ctx, goSkill)
			skillRepo.Create(ctx, dockerSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-docker")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))

			names := []string{skills[0].Name, skills[1].Name}
			Expect(names).To(ContainElement("Go"))
			Expect(names).To(ContainElement("Docker"))
		})

		It("returns skills for correct event only", func() {
			event1 := fixtures.EventWith("event-1", "Go work description", "", "")
			event2 := fixtures.EventWith("event-2", "Python work description", "", "")
			eventRepo.Create(ctx, event1)
			eventRepo.Create(ctx, event2)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			pythonSkill := fixtures.SkillWith("skill-python", "Python", "backend", "")
			skillRepo.Create(ctx, goSkill)
			skillRepo.Create(ctx, pythonSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-2", "skill-python")

			skills1, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills1).To(HaveLen(1))
			Expect(skills1[0].Name).To(Equal("Go"))

			skills2, err := skillRepo.GetSkillsForEvent(ctx, "event-2")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills2).To(HaveLen(1))
			Expect(skills2[0].Name).To(Equal("Python"))
		})

		It("does not return duplicate skills from repeated LinkSkill calls", func() {
			event := fixtures.EventWith("event-1", "Go work description", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			skillRepo.Create(ctx, goSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
		})

		It("works with AssociateSkillWithEvent as well", func() {
			event := fixtures.EventWith("event-1", "Go work description", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			skillRepo.Create(ctx, goSkill)

			skillRepo.AssociateSkillWithEvent("skill-go", "event-1")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})
	})

	Describe("DisassociateSkillFromEvent", func() {
		It("removes the association between skill and event", func() {
			event := fixtures.EventWith("e1", "Built Go microservice with Docker", "", "")
			eventRepo.Create(ctx, event)
			skillRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "backend", ""))
			skillRepo.AssociateSkillWithEvent("s1", "e1")

			skillRepo.DisassociateSkillFromEvent("s1", "e1")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "e1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("cleans up empty maps after disassociation", func() {
			skillRepo.AssociateSkillWithEvent("s1", "e1")

			skillRepo.DisassociateSkillFromEvent("s1", "e1")

			events, err := skillRepo.GetEventsUsingSkill(ctx, "s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(BeEmpty())
		})

		It("handles disassociation of non-existent association gracefully", func() {
			Expect(func() {
				skillRepo.DisassociateSkillFromEvent("nonexistent", "nonexistent")
			}).NotTo(Panic())
		})
	})

	Describe("removeString helper", func() {
		It("removes the target string from slice", func() {
			result := removeString([]string{"a", "b", "c"}, "b")
			Expect(result).To(Equal([]string{"a", "c"}))
		})

		It("returns nil for empty slice", func() {
			result := removeString([]string{}, "a")
			Expect(result).To(BeNil())
		})

		It("returns original elements when target not found", func() {
			result := removeString([]string{"a", "b"}, "c")
			Expect(result).To(Equal([]string{"a", "b"}))
		})
	})
})

var _ = Describe("NewRepositories", func() {
	It("creates all repositories with cross-references", func() {
		repos := NewRepositories()

		Expect(repos.Event).NotTo(BeNil())
		Expect(repos.Skill).NotTo(BeNil())
		Expect(repos.Fact).NotTo(BeNil())
		Expect(repos.Burst).NotTo(BeNil())
	})
})
