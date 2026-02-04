package skillinference_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
		ctx = context.Background() //nolint:fatcontext // test setup

		// Simulate a confirmed burst with events
		testBurst = fixtures.Burst("burst-1", "event-1", "event-2", "event-3")

		testEvents = []*career.Event{
			fixtures.EventWith("event-1", "Built API gateway with Go and integrated with PostgreSQL database", "", ""),
			fixtures.EventWith("event-2", "Deployed services to Kubernetes using Docker containers", "", ""),
			fixtures.EventWith("event-3", "Implemented service mesh with Istio for traffic management", "", ""),
		}
		testEvents[0].Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		testEvents[1].Date = time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
		testEvents[2].Date = time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)

		// Add events to mock repo
		for _, event := range testEvents {
			eventRepo.events[event.ID] = event
		}
	})

	Describe("Complete Workflow", func() {
		It("should infer skills from burst events and create them", func() {
			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).NotTo(BeEmpty())

			skillNames := extractSkillNames(result.Suggestions)
			Expect(skillNames).To(ContainElement("Go"))
			Expect(skillNames).To(ContainElement("PostgreSQL"))
			Expect(skillNames).To(ContainElement("Kubernetes"))
			Expect(skillNames).To(ContainElement("Docker"))
			Expect(skillNames).To(ContainElement("Istio"))

			for _, suggestion := range result.Suggestions {
				Expect(suggestion.Confidence).To(BeNumerically(">", 0))
				Expect(suggestion.Confidence).To(BeNumerically("<=", 1.0))
			}

			acceptedSuggestions := result.Suggestions

			skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(len(acceptedSuggestions)))
			Expect(skillRepo.skills).To(HaveLen(len(acceptedSuggestions)))

			for _, eventID := range testBurst.EventIDs {
				event := eventRepo.events[eventID]
				Expect(event.Skills).NotTo(BeEmpty(), "Event %s should have linked skills", eventID)
			}

			for _, skill := range skills {
				Expect(skill.LastUsed).NotTo(BeNil())
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
			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			acceptedSuggestions := []skillinference.SkillSuggestion{}
			for _, suggestion := range result.Suggestions {
				if suggestion.Name == "Go" || suggestion.Name == "Kubernetes" {
					acceptedSuggestions = append(acceptedSuggestions, suggestion)
				}
			}

			skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))
			Expect(skillRepo.skills).To(HaveLen(2))

			skillNames := make([]string, len(skills))
			for i, skill := range skills {
				skillNames[i] = skill.Name
			}
			Expect(skillNames).To(ConsistOf("Go", "Kubernetes"))
		})

		It("should include existing skills in suggestions and report them in ExistingSkillNames", func() {
			existingTime := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
			existingSkill := fixtures.SkillWith("skill-existing", "Go", "Backend", "intermediate")
			existingSkill.LastUsed = &existingTime
			skillRepo.skills["go"] = existingSkill
			skillRepo.skillsByID["skill-existing"] = existingSkill

			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			suggestionNames := extractSkillNames(result.Suggestions)
			Expect(suggestionNames).To(ContainElement("Go"),
				"Go should still appear in suggestions even though it exists")
			Expect(suggestionNames).To(ContainElement("PostgreSQL"))
			Expect(suggestionNames).To(ContainElement("Docker"))

			Expect(result.ExistingSkillNames).To(ContainElement("Go"))
		})

		It("should report existing skill names while keeping them in suggestions", func() {
			existingTime := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)

			goSkill := fixtures.SkillWith("skill-go", "Go", "Backend", "intermediate")
			goSkill.LastUsed = &existingTime
			skillRepo.skills["go"] = goSkill
			skillRepo.skillsByID["skill-go"] = goSkill

			pgSkill := fixtures.SkillWith("skill-pg", "PostgreSQL", "Database", "intermediate")
			pgSkill.LastUsed = &existingTime
			skillRepo.skills["postgresql"] = pgSkill
			skillRepo.skillsByID["skill-pg"] = pgSkill

			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.ExistingSkillNames).To(ContainElement("Go"))
			Expect(result.ExistingSkillNames).To(ContainElement("PostgreSQL"))

			suggestionNames := extractSkillNames(result.Suggestions)
			Expect(suggestionNames).To(ContainElement("Go"),
				"Go should remain in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("PostgreSQL"),
				"PostgreSQL should remain in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("Docker"))
			Expect(suggestionNames).To(ContainElement("Kubernetes"))
		})

		It("should return empty ExistingSkillNames when no skills pre-exist", func() {
			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.ExistingSkillNames).To(BeEmpty())
			Expect(result.Suggestions).NotTo(BeEmpty())
		})
	})

	Describe("UI Integration Points", func() {
		It("should provide context for skill suggestion modal", func() {
			result, err := service.InferSkillsFromBurst(ctx, testBurst, testEvents)
			Expect(err).NotTo(HaveOccurred())

			for _, suggestion := range result.Suggestions {
				Expect(suggestion.Name).NotTo(BeEmpty(), "Skill name required for display")
				Expect(suggestion.Category).NotTo(BeEmpty(), "Category required for grouping")
				Expect(suggestion.Confidence).To(BeNumerically(">", 0), "Confidence required for sorting")
				Expect(suggestion.EventIDs).NotTo(BeEmpty(), "Event IDs required for context")

				if len(suggestion.Contexts) > 0 {
					Expect(suggestion.Contexts[0]).NotTo(BeEmpty())
				}
			}
		})
	})
})

// Helper to extract skill names from suggestions.
func extractSkillNames(suggestions []skillinference.SkillSuggestion) []string {
	names := make([]string, len(suggestions))
	for i, s := range suggestions {
		names[i] = s.Name
	}
	return names
}
