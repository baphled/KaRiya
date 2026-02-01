package skillinference_test

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillSuggestion", func() {
	Describe("Data Structure", func() {
		It("should have required fields", func() {
			suggestion := skillinference.SkillSuggestion{
				Name:       "Go",
				Category:   "backend",
				Confidence: 0.95,
				EventIDs:   []string{"event-1", "event-2"},
				Contexts:   []string{"...built API using Go..."},
			}

			Expect(suggestion.Name).To(Equal("Go"))
			Expect(suggestion.Category).To(Equal("backend"))
			Expect(suggestion.Confidence).To(Equal(0.95))
			Expect(suggestion.EventIDs).To(HaveLen(2))
			Expect(suggestion.Contexts).To(HaveLen(1))
		})

		It("should support multiple contexts", func() {
			suggestion := skillinference.SkillSuggestion{
				Name:     "PostgreSQL",
				Category: "database",
				Contexts: []string{
					"...migrated to PostgreSQL...",
					"...optimized PostgreSQL queries...",
					"...set up PostgreSQL replication...",
				},
			}

			Expect(suggestion.Contexts).To(HaveLen(3))
		})

		It("should support empty contexts", func() {
			suggestion := skillinference.SkillSuggestion{
				Name:     "Docker",
				Category: "devops",
				Contexts: []string{},
			}

			Expect(suggestion.Contexts).To(BeEmpty())
		})
	})
})

var _ = Describe("SkillInferenceService Interface", func() {
	var (
		ctx    context.Context
		events []*career.Event
	)

	BeforeEach(func() {
		ctx = context.Background()
		events = []*career.Event{
			{
				ID:      "event-1",
				Text:    "Built API using Go and PostgreSQL",
				Company: "TechCorp",
				Project: "Backend Migration",
			},
			{
				ID:      "event-2",
				Text:    "Deployed services to Kubernetes cluster",
				Company: "TechCorp",
				Project: "DevOps Platform",
			},
		}
	})

	Describe("Interface Definition", func() {
		It("should define InferSkillsFromEvents method", func() {
			// This test validates interface compilation
			// Actual implementation tested in Phase 3
			var _ skillinference.SkillInferenceService = (*mockSkillInferenceService)(nil)
		})
	})

	Describe("Mock Implementation", func() {
		var mockService *mockSkillInferenceService

		BeforeEach(func() {
			mockService = &mockSkillInferenceService{}
		})

		It("should implement InferSkillsFromEvents", func() {
			result, err := mockService.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Suggestions).NotTo(BeNil())
		})

		It("should implement InferSkillsFromBurst", func() {
			burst := &career.Burst{
				ID:       "burst-1",
				Name:     "Backend Migration",
				EventIDs: []string{"event-1"},
			}

			result, err := mockService.InferSkillsFromBurst(ctx, burst, events)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Suggestions).NotTo(BeNil())
		})

		It("should implement CreateSkillsFromSuggestions", func() {
			suggestions := []skillinference.SkillSuggestion{
				{
					Name:       "Go",
					Category:   "backend",
					Confidence: 0.95,
					EventIDs:   []string{"event-1"},
					Contexts:   []string{"...built API using Go..."},
				},
			}

			skills, err := mockService.CreateSkillsFromSuggestions(ctx, suggestions)

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).NotTo(BeNil())
		})
	})

	Describe("Method Signatures", func() {
		var mockService *mockSkillInferenceService

		BeforeEach(func() {
			mockService = &mockSkillInferenceService{}
		})

		It("should accept context as first parameter", func() {
			result, err := mockService.InferSkillsFromEvents(ctx, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			burst := &career.Burst{ID: "burst-1"}
			result, err = mockService.InferSkillsFromBurst(ctx, burst, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			_, err = mockService.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error as second return value", func() {
			result, err := mockService.InferSkillsFromEvents(ctx, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})
	})
})

// mockSkillInferenceService is a test implementation of SkillInferenceService.
type mockSkillInferenceService struct{}

func (m *mockSkillInferenceService) InferSkillsFromEvents(
	_ context.Context,
	_ []*career.Event,
) (*skillinference.InferenceResult, error) {
	return &skillinference.InferenceResult{Suggestions: []skillinference.SkillSuggestion{}}, nil
}

func (m *mockSkillInferenceService) InferSkillsFromBurst(
	_ context.Context,
	_ *career.Burst,
	_ []*career.Event,
) (*skillinference.InferenceResult, error) {
	return &skillinference.InferenceResult{Suggestions: []skillinference.SkillSuggestion{}}, nil
}

func (m *mockSkillInferenceService) CreateSkillsFromSuggestions(
	_ context.Context,
	_ []skillinference.SkillSuggestion,
) ([]*career.Skill, error) {
	return []*career.Skill{}, nil
}
