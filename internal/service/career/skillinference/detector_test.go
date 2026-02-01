package skillinference_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultSkillInferenceService", func() {
	var (
		ctx     context.Context
		service skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		ctx = context.Background()
		service = skillinference.NewSkillInferenceService(nil, nil) // No persistence needed for detection tests
	})

	Describe("Word Boundary Detection", func() {
		It("should detect 'go' as a keyword", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API using Go",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).To(HaveLen(1))
			Expect(result.Suggestions[0].Name).To(Equal("Go"))
		})

		It("should NOT match 'go' in 'goal'", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Our goal was to improve performance",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			// Should not detect "Go" from "goal"
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(BeEmpty())
		})

		It("should NOT match 'go' in 'going'", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "We are going to deploy tomorrow",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(BeEmpty())
		})

		It("should match 'go' with punctuation", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built with Go, PostgreSQL, and Docker",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(result.Suggestions)).To(BeNumerically(">=", 3))

			names := extractNames(result.Suggestions)
			Expect(names).To(ContainElement("Go"))
			Expect(names).To(ContainElement("PostgreSQL"))
			Expect(names).To(ContainElement("Docker"))
		})
	})

	Describe("Case-Insensitive Matching", func() {
		It("should match 'Go' (uppercase)", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).NotTo(BeEmpty())
			Expect(result.Suggestions[0].Name).To(Equal("Go"))
		})

		It("should match 'go' (lowercase)", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Wrote microservices in go",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).NotTo(BeEmpty())
			Expect(result.Suggestions[0].Name).To(Equal("Go"))
		})

		It("should match 'GO' (all caps)", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Migrated to GO from Java",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			names := extractNames(result.Suggestions)
			Expect(names).To(ContainElement("Go"))
		})
	})

	Describe("Context Extraction", func() {
		It("should extract context around keyword", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built a scalable API using Go and gRPC for microservices architecture",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Contexts).NotTo(BeEmpty())
			Expect(goSkill.Contexts[0]).To(ContainSubstring("using Go"))
		})

		It("should limit context to ~80 chars", func() {
			longText := "This is a very long event description that goes on and on with lots of details about how we built an amazing API using Go and gRPC for a distributed microservices architecture that handles millions of requests per day"
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: longText,
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Contexts).NotTo(BeEmpty())

			// Context should be truncated with ellipsis
			context := goSkill.Contexts[0]
			Expect(len(context)).To(BeNumerically("<=", 120)) // ~80 + ellipsis
		})

		It("should add ellipsis when truncated", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "First part of text before keyword and then we used Go for the implementation and then lots more text after",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Contexts).NotTo(BeEmpty())

			context := goSkill.Contexts[0]
			Expect(context).To(MatchRegexp(`\.\.\..*\.\.\.`)) // Has ellipsis on both sides
		})
	})

	Describe("Deduplication", func() {
		It("should deduplicate same skill across events", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go",
					Date: time.Now(),
				},
				{
					ID:   "event-2",
					Text: "Wrote Go microservices",
					Date: time.Now(),
				},
				{
					ID:   "event-3",
					Text: "Deployed Go services",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())

			// Should have ONE "Go" suggestion (not three)
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))

			// Should reference all 3 events
			Expect(goSkills[0].EventIDs).To(HaveLen(3))
			Expect(goSkills[0].EventIDs).To(ContainElement("event-1"))
			Expect(goSkills[0].EventIDs).To(ContainElement("event-2"))
			Expect(goSkills[0].EventIDs).To(ContainElement("event-3"))
		})

		It("should merge contexts from multiple events", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go",
					Date: time.Now(),
				},
				{
					ID:   "event-2",
					Text: "Wrote Go microservices",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Contexts).To(HaveLen(2))
		})

		It("should limit contexts to 3 max", func() {
			events := []*career.Event{
				{ID: "event-1", Text: "First Go usage", Date: time.Now()},
				{ID: "event-2", Text: "Second Go usage", Date: time.Now()},
				{ID: "event-3", Text: "Third Go usage", Date: time.Now()},
				{ID: "event-4", Text: "Fourth Go usage", Date: time.Now()},
				{ID: "event-5", Text: "Fifth Go usage", Date: time.Now()},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Contexts).To(HaveLen(3)) // Max 3 contexts
		})
	})

	Describe("Alias Support", func() {
		It("should detect 'golang' and return 'Go'", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built services with golang",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).NotTo(BeEmpty())
			// Should return canonical name "Go" not "golang"
			Expect(result.Suggestions[0].Name).To(Equal("Go"))
		})

		It("should deduplicate 'go' and 'golang' as same skill", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go",
					Date: time.Now(),
				},
				{
					ID:   "event-2",
					Text: "Wrote golang microservices",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())

			// Should have ONE "Go" suggestion (not two)
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))
			Expect(goSkills[0].EventIDs).To(HaveLen(2))
		})

		It("should detect 'k8s' and return 'Kubernetes'", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Deployed to k8s cluster",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(result.Suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
		})
	})

	Describe("Empty Results", func() {
		It("should return empty slice when no keywords found", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Had a meeting about the project",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).To(BeEmpty())
		})

		It("should return empty slice for empty events", func() {
			events := []*career.Event{}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).To(BeEmpty())
		})

		It("should return empty slice for nil events", func() {
			result, err := service.InferSkillsFromEvents(ctx, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Suggestions).To(BeEmpty())
		})
	})

	Describe("Multiple Technologies", func() {
		It("should detect multiple technologies in one event", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go, PostgreSQL, and deployed to Kubernetes",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(result.Suggestions)).To(BeNumerically(">=", 3))

			names := extractNames(result.Suggestions)
			Expect(names).To(ContainElement("Go"))
			Expect(names).To(ContainElement("PostgreSQL"))
			Expect(names).To(ContainElement("Kubernetes"))
		})

		It("should assign correct categories", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go, PostgreSQL, and React",
					Date: time.Now(),
				},
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())

			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Category).To(Equal("backend"))

			pgSkill := findByName(result.Suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Category).To(Equal("database"))

			reactSkill := findByName(result.Suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Category).To(Equal("frontend"))
		})
	})
})

// Helper functions for tests

func filterByName(suggestions []skillinference.SkillSuggestion, name string) []skillinference.SkillSuggestion { //nolint:unparam // test helper intentionally called with fixed values
	var result []skillinference.SkillSuggestion
	for _, s := range suggestions {
		if s.Name == name {
			result = append(result, s)
		}
	}
	return result
}

func findByName(suggestions []skillinference.SkillSuggestion, name string) *skillinference.SkillSuggestion {
	for i, s := range suggestions {
		if s.Name == name {
			return &suggestions[i]
		}
	}
	return nil
}

func extractNames(suggestions []skillinference.SkillSuggestion) []string {
	var names []string
	for _, s := range suggestions {
		names = append(names, s.Name)
	}
	return names
}
