package skillinference_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// Integration test showing the complete workflow:
// 1. User confirms a burst
// 2. System infers skills from burst events
// 3. User reviews suggestions
// 4. User accepts suggestions
// 5. System creates skills and links to events

var _ = Describe("Skill Inference Integration", func() {
	var (
		service    skillinference.SkillInferenceService
		skillRepo  *mockSkillRepository
		eventRepo  *mockEventRepository
		ctx        context.Context
		testBurst  *career.Burst
		testEvents []*career.Event
	)

	BeforeEach(func() {
		skillRepo = &mockSkillRepository{
			skills:     make(map[string]*career.Skill),
			skillsByID: make(map[string]*career.Skill),
		}
		eventRepo = &mockEventRepository{
			events: make(map[string]*career.Event),
		}

		service = skillinference.NewSkillInferenceService(skillRepo, eventRepo)
		ctx = context.Background()

		// Simulate a confirmed burst with events
		testBurst = &career.Burst{
			ID:          "burst-1",
			Name:        "Microservices Migration",
			Description: "Migrated monolith to microservices",
			EventIDs:    []string{"event-1", "event-2", "event-3"},
		}

		testEvents = []*career.Event{
			{
				ID:   "event-1",
				Text: "Built API gateway with Go and integrated with PostgreSQL database",
				Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:   "event-2",
				Text: "Deployed services to Kubernetes using Docker containers",
				Date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:   "event-3",
				Text: "Implemented service mesh with Istio for traffic management",
				Date: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			},
		}

		// Add events to mock repo
		for _, event := range testEvents {
			eventRepo.events[event.ID] = event
		}
	})

	Describe("Complete Workflow", func() {
		It("should infer skills from burst events and create them", func() {
			// Step 1: User confirms burst (outside this test, done in burst_management intent)
			// Burst is already confirmed in testBurst

			// Step 2: System infers skills from burst events
			suggestions, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)

			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).NotTo(BeEmpty())

			// Verify detected skills
			skillNames := extractSkillNames(suggestions)
			Expect(skillNames).To(ContainElement("Go"))
			Expect(skillNames).To(ContainElement("PostgreSQL"))
			Expect(skillNames).To(ContainElement("Kubernetes"))
			Expect(skillNames).To(ContainElement("Docker"))
			Expect(skillNames).To(ContainElement("Istio"))

			// Verify confidence scores are set
			for _, suggestion := range suggestions {
				Expect(suggestion.Confidence).To(BeNumerically(">", 0))
				Expect(suggestion.Confidence).To(BeNumerically("<=", 1.0))
			}

			// Step 3: User reviews suggestions in UI (simulated by selecting all)
			// In real UI, user would see modal with checkboxes to accept/reject

			// Step 4: User accepts all suggestions
			acceptedSuggestions := suggestions // In real UI, only checked suggestions

			// Step 5: System creates skills and links to events
			skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(len(acceptedSuggestions)))

			// Verify skills were created in repository
			Expect(skillRepo.skills).To(HaveLen(len(acceptedSuggestions)))

			// Verify skills are linked to events
			for _, eventID := range testBurst.EventIDs {
				event := eventRepo.events[eventID]
				Expect(event.Skills).NotTo(BeEmpty(), "Event %s should have linked skills", eventID)
			}

			// Verify LastUsed is set correctly
			for _, skill := range skills {
				Expect(skill.LastUsed).NotTo(BeNil())
				// LastUsed should be one of the event dates
				validDate := false
				for _, event := range testEvents {
					if skill.LastUsed.Equal(event.Date) {
						validDate = true
						break
					}
				}
				Expect(validDate).To(BeTrue(), "Skill %s has invalid LastUsed date", skill.Name)
			}
		})

		It("should handle partial acceptance of suggestions", func() {
			// Infer skills
			suggestions, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			// User accepts only Go and Kubernetes (simulating checkbox selection in UI)
			acceptedSuggestions := []skillinference.SkillSuggestion{}
			for _, suggestion := range suggestions {
				if suggestion.Name == "Go" || suggestion.Name == "Kubernetes" {
					acceptedSuggestions = append(acceptedSuggestions, suggestion)
				}
			}

			// Create only accepted skills
			skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			Expect(skillRepo.skills).To(HaveLen(2))

			// Verify only accepted skills were created
			skillNames := make([]string, len(skills))
			for i, skill := range skills {
				skillNames[i] = skill.Name
			}
			Expect(skillNames).To(ConsistOf("Go", "Kubernetes"))
		})

		It("should handle reuse of existing skills", func() {
			// Pre-create a skill
			existingTime := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
			existingSkill := &career.Skill{
				ID:       "skill-existing",
				Name:     "Go",
				Category: "Backend",
				LastUsed: &existingTime,
			}
			skillRepo.skills["go"] = existingSkill
			skillRepo.skillsByID["skill-existing"] = existingSkill

			// Infer skills (will detect Go again)
			suggestions, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			// Accept all suggestions
			skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)
			Expect(err).NotTo(HaveOccurred())

			// Verify Go skill was reused (not recreated)
			goSkill := findSkillByName(skills, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.ID).To(Equal("skill-existing"))

			// Verify LastUsed was updated to more recent date
			Expect(goSkill.LastUsed).NotTo(BeNil())
			Expect(goSkill.LastUsed.After(existingTime)).To(BeTrue())
		})
	})

	Describe("UI Integration Points", func() {
		It("should provide context for skill suggestion modal", func() {
			suggestions, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			// Verify suggestions contain all data needed for UI display
			for _, suggestion := range suggestions {
				Expect(suggestion.Name).NotTo(BeEmpty(), "Skill name required for display")
				Expect(suggestion.Category).NotTo(BeEmpty(), "Category required for grouping")
				Expect(suggestion.Confidence).To(BeNumerically(">", 0), "Confidence required for sorting")
				Expect(suggestion.EventIDs).NotTo(BeEmpty(), "Event IDs required for context")

				// Contexts should be available for showing usage examples
				if len(suggestion.Contexts) > 0 {
					Expect(suggestion.Contexts[0]).NotTo(BeEmpty())
				}
			}
		})
	})
})

// Helper to extract skill names from suggestions
func extractSkillNames(suggestions []skillinference.SkillSuggestion) []string {
	names := make([]string, len(suggestions))
	for i, s := range suggestions {
		names[i] = s.Name
	}
	return names
}
