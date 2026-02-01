package skillinference_test

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Confidence Scoring", func() {
	var (
		ctx     context.Context
		service skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		ctx = context.Background()                                  //nolint:fatcontext // test setup
		service = skillinference.NewSkillInferenceService(nil, nil) // No persistence needed for confidence tests
	})

	Describe("High Confidence Patterns (0.95)", func() {
		It("should score 'built with X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Built API with Go for microservices", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'built using X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Built microservices using Go and gRPC", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'developed in X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Developed backend services in Go", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'developed using X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Developed REST API using Go", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'implemented in X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Implemented authentication service in Go", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'implemented using X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Implemented caching layer using Redis", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			redisSkill := findByName(result.Suggestions, "Redis")
			Expect(redisSkill).NotTo(BeNil())
			Expect(redisSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'wrote X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Wrote Go services for data processing", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'using X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Deployed services using Kubernetes and Helm", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(result.Suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'X developer' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Worked as a Go developer on backend team", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'X engineer' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Served as Python engineer for data platform", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pythonSkill := findByName(result.Suggestions, "Python")
			Expect(pythonSkill).NotTo(BeNil())
			Expect(pythonSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'expert in X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Became expert in PostgreSQL optimization", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(result.Suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'proficient in X' as high confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Became proficient in React development", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			reactSkill := findByName(result.Suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'X project' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Led React project for frontend rewrite", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			reactSkill := findByName(result.Suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X system' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Maintained PostgreSQL system for user data", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(result.Suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X application' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Debugged Node.js application performance issues", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			nodeSkill := findByName(result.Suggestions, "Node.js")
			Expect(nodeSkill).NotTo(BeNil())
			Expect(nodeSkill.Confidence).To(Equal(0.75))
		})
	})

	Describe("Medium Confidence Patterns (0.75)", func() {
		It("should score 'worked with X' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Worked with Docker for containerization", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			dockerSkill := findByName(result.Suggestions, "Docker")
			Expect(dockerSkill).NotTo(BeNil())
			Expect(dockerSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'experience with X' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Gained experience with Kubernetes deployments", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(result.Suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X project' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Led React project for frontend rewrite", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			reactSkill := findByName(result.Suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X system' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Maintained PostgreSQL system for user data", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(result.Suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X application' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Debugged Node.js application performance issues", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			nodeSkill := findByName(result.Suggestions, "Node.js")
			Expect(nodeSkill).NotTo(BeNil())
			Expect(nodeSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X service' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Monitored Redis service for caching layer", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			redisSkill := findByName(result.Suggestions, "Redis")
			Expect(redisSkill).NotTo(BeNil())
			Expect(redisSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'migrated to X' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Migrated to Kubernetes from EC2 instances", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(result.Suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'integrated X' as medium confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Integrated GraphQL for API layer", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			graphqlSkill := findByName(result.Suggestions, "GraphQL")
			Expect(graphqlSkill).NotTo(BeNil())
			Expect(graphqlSkill.Confidence).To(Equal(0.75))
		})
	})

	Describe("Low Confidence (0.5)", func() {
		It("should score simple keyword presence as low confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "The team discussed Docker at meeting", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			dockerSkill := findByName(result.Suggestions, "Docker")
			Expect(dockerSkill).NotTo(BeNil())
			Expect(dockerSkill.Confidence).To(Equal(0.5))
		})

		It("should score 'with X' (no action verb) as low confidence", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Attended meeting with Go team", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.5))
		})
	})

	Describe("Case-Insensitive Pattern Matching", func() {
		It("should match 'BUILT WITH' (uppercase)", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "BUILT API WITH GO", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should match 'Built With' (mixed case)", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Built With Go And PostgreSQL", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(result.Suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})
	})

	Describe("Multiple Pattern Priority", func() {
		It("should use highest confidence when multiple patterns match", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Worked with Go and built microservices using Go", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))
			Expect(goSkills[0].Confidence).To(Equal(0.95))
		})

		It("should prefer 'built with' over 'worked with'", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Built API with PostgreSQL after working with MongoDB", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(result.Suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.95))
			mongoSkill := findByName(result.Suggestions, "MongoDB")
			Expect(mongoSkill).NotTo(BeNil())
			Expect(mongoSkill.Confidence).To(Equal(0.75))
		})
	})

	Describe("Deduplication with Confidence Merging", func() {
		It("should keep highest confidence when merging same skill", func() {
			events := []*career.Event{
				fixtures.EventWith("event-1", "Discussed Go at meeting", "", ""),
				fixtures.EventWith("event-2", "Built API with Go", "", ""),
			}

			result, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkills := filterByName(result.Suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))
			Expect(goSkills[0].Confidence).To(Equal(0.95))
			Expect(goSkills[0].EventIDs).To(HaveLen(2))
		})
	})
})

// Helper functions
func findByName(suggestions []skillinference.SkillSuggestion, name string) *skillinference.SkillSuggestion {
	for i := range suggestions {
		if suggestions[i].Name == name {
			return &suggestions[i]
		}
	}
	return nil
}

//nolint:unparam // name parameter allows reuse with different skills in future tests
func filterByName(suggestions []skillinference.SkillSuggestion, name string) []skillinference.SkillSuggestion {
	var filtered []skillinference.SkillSuggestion
	for _, s := range suggestions {
		if s.Name == name {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
