package service

import (
	"context"
	"time"

	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Event Service", func() {
	var (
		cliEventService *CLIEventService
		careerSvc       *careerservice.Service
		repo            careerrepo.EventRepository
	)

	BeforeEach(func() {
		// Set up in-memory repository for testing
		repo = careermemory.NewEventRepository()
		careerSvc = careerservice.NewService(repo)
		cliEventService = NewCLIEventService(careerSvc)
	})

	Context("Capturing Events", func() {
		It("should capture an event with default options", func() {
			ctx := context.Background()
			eventText := "Completed a significant project"
			eventDate := time.Now()
			mode := careerservice.TimelineJournaling

			// Capture the event
			err := cliEventService.CaptureEvent(ctx, eventText, eventDate, mode)

			// Assertions
			Expect(err).ToNot(HaveOccurred())
		})

		It("should capture an event with optional configurations", func() {
			ctx := context.Background()
			eventText := "Led a cross-functional team"
			eventDate := time.Now()
			mode := careerservice.TimelineJournaling

			// Capture the event with optional configurations
			err := cliEventService.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				mode,
				WithCompany("TechCorp"),
				WithProject("Platform Migration"),
				WithTags([]string{"leadership", "technical"}),
			)

			// Assertions
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Listing Events", func() {
		It("should support listing events", func() {
			ctx := context.Background()

			// Create a test event directly via repository
			event := fixtures.Event("test-list-1")
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// List events using CLI service
			filters := fixtures.EventListFilters()
			events, err := cliEventService.ListEvents(ctx, filters)

			// Assertions
			Expect(err).ToNot(HaveOccurred())
			Expect(events).ToNot(BeNil())
		})
	})

	Context("Get Event By ID", func() {
		It("should retrieve a specific event", func() {
			ctx := context.Background()

			// Create a test event
			event := fixtures.EventWith("", "Completed important milestone", "", "")
			err := careerSvc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).ToNot(HaveOccurred())
			eventID := event.ID

			// Get event by ID using CLI service
			retrievedEvent, err := cliEventService.GetEventByID(ctx, eventID)

			// Assertions
			Expect(err).ToNot(HaveOccurred())
			Expect(retrievedEvent).ToNot(BeNil())
			Expect(retrievedEvent.ID).To(Equal(eventID))
			Expect(retrievedEvent.Text).To(Equal("Completed important milestone"))
		})
	})
})

var _ = Describe("Skill Management Operations", func() {
	var (
		eventRepo   *careermemory.EventRepository
		skillRepo   *careermemory.SkillRepository
		careerSvc   *careerservice.Service
		cliEventSvc *CLIEventService
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		eventRepo = careermemory.NewEventRepository()
		skillRepo = careermemory.NewSkillRepository()
		eventRepo.SetSkillRepository(skillRepo)
		skillRepo.SetEventRepository(eventRepo)
		careerSvc = careerservice.NewService(eventRepo)
		careerSvc.SetSkillRepository(skillRepo)
		cliEventSvc = NewCLIEventService(careerSvc)
	})

	Describe("LinkSkillToEvent", func() {
		It("should link an existing skill to an event", func() {
			event := fixtures.EventWith("event-1", "Implemented feature", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			err = skillRepo.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			err = cliEventSvc.LinkSkillToEvent(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].ID).To(Equal("skill-1"))
		})

		It("should return error when event does not exist", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			err := skillRepo.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			err = cliEventSvc.LinkSkillToEvent(ctx, "non-existent", "skill-1")
			Expect(err).To(HaveOccurred())
		})

		It("should allow linking non-existent skill ID", func() {
			event := fixtures.EventWith("event-1", "Implemented feature", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			err = cliEventSvc.LinkSkillToEvent(ctx, "event-1", "non-existent")
			Expect(err).ToNot(HaveOccurred())

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("UnlinkSkillFromEvent", func() {
		It("should unlink an existing skill from an event", func() {
			event := fixtures.EventWith("event-1", "Implemented feature", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			err = skillRepo.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			err = eventRepo.LinkSkill(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			err = cliEventSvc.UnlinkSkillFromEvent(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("should return error when event does not exist", func() {
			err := cliEventSvc.UnlinkSkillFromEvent(ctx, "non-existent", "skill-1")
			Expect(err).To(HaveOccurred())
		})

		It("should handle gracefully when skill is not linked", func() {
			event := fixtures.EventWith("event-1", "Implemented feature", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			err = cliEventSvc.UnlinkSkillFromEvent(ctx, "event-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ListAllSkills", func() {
		It("should return all skills without filters", func() {
			skill1 := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			skill2 := fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate")
			skill3 := fixtures.SkillWith("skill-3", "React", "frontend", "advanced")

			err := skillRepo.Create(ctx, skill1)
			Expect(err).ToNot(HaveOccurred())
			err = skillRepo.Create(ctx, skill2)
			Expect(err).ToNot(HaveOccurred())
			err = skillRepo.Create(ctx, skill3)
			Expect(err).ToNot(HaveOccurred())

			skills, err := cliEventSvc.ListAllSkills(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(3))
		})

		It("should return empty slice when no skills exist", func() {
			skills, err := cliEventSvc.ListAllSkills(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})
})

var _ = Describe("UpdateEventMetadata", func() {
	var (
		repo   *careermemory.EventRepository
		svc    *careerservice.Service
		cliSvc *CLIEventService
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careermemory.NewEventRepository()
		svc = careerservice.NewService(repo)
		cliSvc = NewCLIEventService(svc)
	})

	It("should update event metadata without changing text or date", func() {
		// Create an event
		originalEvent := fixtures.EventWith("test-event-1", "Original text", "", "")
		originalEvent.Date = time.Now().Add(-24 * time.Hour)
		err := repo.Create(ctx, originalEvent)
		Expect(err).ToNot(HaveOccurred())

		// Update metadata
		updatedEvent := fixtures.EventWith("test-event-1", "This should be ignored", "NewCompany", "NewProject")
		updatedEvent.Tags = []string{"technical", "leadership"}
		updatedEvent.Categories = []string{"Technical"}

		err = cliSvc.UpdateEventMetadata(ctx, updatedEvent)
		Expect(err).ToNot(HaveOccurred())

		// Verify metadata was updated
		retrieved, err := svc.GetEventByID(ctx, "test-event-1")
		Expect(err).ToNot(HaveOccurred())
		Expect(retrieved.Company).To(Equal("NewCompany"))
		Expect(retrieved.Project).To(Equal("NewProject"))
		Expect(retrieved.Tags).To(ContainElements("technical", "leadership"))

		// Verify text and date were not changed
		Expect(retrieved.Text).To(Equal("Original text"))
		Expect(retrieved.Date).To(Equal(originalEvent.Date))
	})

	It("should return error when event is nil", func() {
		err := cliSvc.UpdateEventMetadata(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("event cannot be nil"))
	})

	It("should return error when event ID is empty", func() {
		event := fixtures.EventWith("", "Test", "Company", "")
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
	})

	It("should return error when event does not exist", func() {
		event := fixtures.EventWith("non-existent", "", "Company", "")
		err := cliSvc.UpdateEventMetadata(ctx, event)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("UpdateEvent", func() {
	var (
		repo        *careermemory.EventRepository
		careerSvc   *careerservice.Service
		cliEventSvc *CLIEventService
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careermemory.NewEventRepository()
		careerSvc = careerservice.NewService(repo)
		cliEventSvc = NewCLIEventService(careerSvc)
	})

	It("should update an existing event", func() {
		event := fixtures.EventWith("update-1", "Original career event text", "", "")
		err := repo.Create(ctx, event)
		Expect(err).ToNot(HaveOccurred())

		err = cliEventSvc.UpdateEvent(ctx, "update-1", "Updated career event text", time.Now())
		Expect(err).ToNot(HaveOccurred())

		retrieved, err := careerSvc.GetEventByID(ctx, "update-1")
		Expect(err).ToNot(HaveOccurred())
		Expect(retrieved.Text).To(Equal("Updated career event text"))
	})

	It("should apply optional configurations during update", func() {
		event := fixtures.EventWith("update-2", "Initial career event text", "", "")
		err := repo.Create(ctx, event)
		Expect(err).ToNot(HaveOccurred())

		err = cliEventSvc.UpdateEvent(
			ctx,
			"update-2",
			"Updated career event text",
			time.Now(),
			WithCompany("NewCorp"),
			WithProject("Migration"),
			WithCategories([]string{"technical", "leadership"}),
			WithTags([]string{"project", "achievement"}),
		)
		Expect(err).ToNot(HaveOccurred())

		retrieved, err := careerSvc.GetEventByID(ctx, "update-2")
		Expect(err).ToNot(HaveOccurred())
		Expect(retrieved.Company).To(Equal("NewCorp"))
		Expect(retrieved.Project).To(Equal("Migration"))
		Expect(retrieved.Categories).To(ContainElements("technical", "leadership"))
		Expect(retrieved.Tags).To(ContainElements("project", "achievement"))
	})

	It("should return error when event does not exist", func() {
		err := cliEventSvc.UpdateEvent(ctx, "non-existent", "Updated career event text", time.Now())
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("DeleteEvent", func() {
	var (
		repo        *careermemory.EventRepository
		careerSvc   *careerservice.Service
		cliEventSvc *CLIEventService
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careermemory.NewEventRepository()
		careerSvc = careerservice.NewService(repo)
		cliEventSvc = NewCLIEventService(careerSvc)
	})

	It("should delete an existing event", func() {
		event := fixtures.EventWith("delete-1", "Career event to be deleted", "", "")
		err := repo.Create(ctx, event)
		Expect(err).ToNot(HaveOccurred())

		err = cliEventSvc.DeleteEvent(ctx, "delete-1")
		Expect(err).ToNot(HaveOccurred())

		_, err = careerSvc.GetEventByID(ctx, "delete-1")
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Event Options", func() {
	Describe("WithCategories", func() {
		It("should set categories on event config", func() {
			config := defaultConfig()
			categories := []string{"technical", "leadership"}
			opt := WithCategories(categories)
			opt(config)
			Expect(config.Categories).To(Equal([]string{"technical", "leadership"}))
		})
	})
})

var _ = Describe("GetSkillsForEvent", func() {
	Context("when skill repository is configured", func() {
		var (
			eventRepo   *careermemory.EventRepository
			skillRepo   *careermemory.SkillRepository
			careerSvc   *careerservice.Service
			cliEventSvc *CLIEventService
			ctx         context.Context
		)

		BeforeEach(func() {
			ctx = context.Background()
			eventRepo = careermemory.NewEventRepository()
			skillRepo = careermemory.NewSkillRepository()
			eventRepo.SetSkillRepository(skillRepo)
			skillRepo.SetEventRepository(eventRepo)
			careerSvc = careerservice.NewService(eventRepo)
			careerSvc.SetSkillRepository(skillRepo)
			cliEventSvc = NewCLIEventService(careerSvc)
		})

		It("should return skills associated with an event", func() {
			event := fixtures.EventWith("event-skills-1", "Implemented feature for testing", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			err = skillRepo.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			err = eventRepo.LinkSkill(ctx, "event-skills-1", "skill-1")
			Expect(err).ToNot(HaveOccurred())

			skills, err := cliEventSvc.GetSkillsForEvent(ctx, "event-skills-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].ID).To(Equal("skill-1"))
		})

		It("should return empty slice when no skills are linked", func() {
			event := fixtures.EventWith("event-skills-2", "Implemented feature for testing", "", "")
			err := eventRepo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			skills, err := cliEventSvc.GetSkillsForEvent(ctx, "event-skills-2")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Context("when skill repository is not configured", func() {
		It("should return empty slice", func() {
			ctx := context.Background()
			eventRepo := careermemory.NewEventRepository()
			careerSvc := careerservice.NewService(eventRepo)
			cliEventSvc := NewCLIEventService(careerSvc)

			skills, err := cliEventSvc.GetSkillsForEvent(ctx, "any-event")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Context("when service is nil", func() {
		It("should return empty slice", func() {
			ctx := context.Background()
			cliEventSvc := &CLIEventService{service: nil}

			skills, err := cliEventSvc.GetSkillsForEvent(ctx, "any-event")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})
})

var _ = Describe("ListEvents edge cases", func() {
	It("should handle nil filters", func() {
		ctx := context.Background()
		repo := careermemory.NewEventRepository()
		careerSvc := careerservice.NewService(repo)
		cliEventSvc := NewCLIEventService(careerSvc)

		events, err := cliEventSvc.ListEvents(ctx, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(events).To(BeEmpty())
	})
})

var _ = Describe("ListAllSkills edge cases", func() {
	It("should return empty slice when skill repository is not configured", func() {
		ctx := context.Background()
		repo := careermemory.NewEventRepository()
		careerSvc := careerservice.NewService(repo)
		cliEventSvc := NewCLIEventService(careerSvc)

		skills, err := cliEventSvc.ListAllSkills(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(skills).To(BeEmpty())
	})
})
