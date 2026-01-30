package skillinference_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Confidence Scoring", func() {
	var (
		ctx     context.Context
		service skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		ctx = context.Background()
		service = skillinference.NewSkillInferenceService()
	})

	Describe("High Confidence Patterns (0.95)", func() {
		It("should score 'built with X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with Go for microservices",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'built using X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built microservices using Go and gRPC",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'developed in X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Developed backend services in Go",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'developed using X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Developed REST API using Go",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'implemented in X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Implemented authentication service in Go",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'implemented using X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Implemented caching layer using Redis",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			redisSkill := findByName(suggestions, "Redis")
			Expect(redisSkill).NotTo(BeNil())
			Expect(redisSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'wrote X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Wrote Go services for data processing",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'using X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Deployed services using Kubernetes and Helm",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'X developer' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Worked as a Go developer on backend team",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'X engineer' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Served as Python engineer for data platform",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pythonSkill := findByName(suggestions, "Python")
			Expect(pythonSkill).NotTo(BeNil())
			Expect(pythonSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'expert in X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Became expert in PostgreSQL optimization",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.95))
		})

		It("should score 'proficient in X' as high confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Became proficient in React development",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			reactSkill := findByName(suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Confidence).To(Equal(0.95))
		})
	})

	Describe("Medium Confidence Patterns (0.75)", func() {
		It("should score 'worked with X' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Worked with Docker for containerization",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			dockerSkill := findByName(suggestions, "Docker")
			Expect(dockerSkill).NotTo(BeNil())
			Expect(dockerSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'experience with X' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Gained experience with Kubernetes deployments",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X project' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Led the React project for frontend rewrite",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			reactSkill := findByName(suggestions, "React")
			Expect(reactSkill).NotTo(BeNil())
			Expect(reactSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X system' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Maintained PostgreSQL system for user data",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			pgSkill := findByName(suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X application' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Debugged Node.js application performance issues",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			nodeSkill := findByName(suggestions, "Node.js")
			Expect(nodeSkill).NotTo(BeNil())
			Expect(nodeSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'X service' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Monitored Redis service for caching layer",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			redisSkill := findByName(suggestions, "Redis")
			Expect(redisSkill).NotTo(BeNil())
			Expect(redisSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'migrated to X' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Migrated to Kubernetes from EC2 instances",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			k8sSkill := findByName(suggestions, "Kubernetes")
			Expect(k8sSkill).NotTo(BeNil())
			Expect(k8sSkill.Confidence).To(Equal(0.75))
		})

		It("should score 'integrated X' as medium confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Integrated GraphQL for API layer",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			graphqlSkill := findByName(suggestions, "GraphQL")
			Expect(graphqlSkill).NotTo(BeNil())
			Expect(graphqlSkill.Confidence).To(Equal(0.75))
		})
	})

	Describe("Low Confidence (0.5)", func() {
		It("should score simple keyword presence as low confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "The team discussed Docker at the meeting",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			dockerSkill := findByName(suggestions, "Docker")
			Expect(dockerSkill).NotTo(BeNil())
			Expect(dockerSkill.Confidence).To(Equal(0.5))
		})

		It("should score 'with X' (no action verb) as low confidence", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Attended meeting with Go team",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.5))
		})
	})

	Describe("Case-Insensitive Pattern Matching", func() {
		It("should match 'BUILT WITH' (uppercase)", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "BUILT API WITH GO",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})

		It("should match 'Built With' (mixed case)", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built With Go And PostgreSQL",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkill := findByName(suggestions, "Go")
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Confidence).To(Equal(0.95))
		})
	})

	Describe("Multiple Pattern Priority", func() {
		It("should use highest confidence when multiple patterns match", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Worked with Go and built microservices using Go",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkills := filterByName(suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))                // Deduplicated
			Expect(goSkills[0].Confidence).To(Equal(0.95)) // Highest confidence wins
		})

		It("should prefer 'built with' over 'worked with'", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Built API with PostgreSQL after working with MongoDB",
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())

			pgSkill := findByName(suggestions, "PostgreSQL")
			Expect(pgSkill).NotTo(BeNil())
			Expect(pgSkill.Confidence).To(Equal(0.95)) // High confidence

			mongoSkill := findByName(suggestions, "MongoDB")
			Expect(mongoSkill).NotTo(BeNil())
			Expect(mongoSkill.Confidence).To(Equal(0.75)) // Medium confidence
		})
	})

	Describe("Deduplication with Confidence Merging", func() {
		It("should keep highest confidence when merging same skill", func() {
			events := []*career.Event{
				{
					ID:   "event-1",
					Text: "Discussed Go at meeting", // Low confidence (0.5)
					Date: time.Now(),
				},
				{
					ID:   "event-2",
					Text: "Built API with Go", // High confidence (0.95)
					Date: time.Now(),
				},
			}

			suggestions, err := service.InferSkillsFromEvents(ctx, events)

			Expect(err).NotTo(HaveOccurred())
			goSkills := filterByName(suggestions, "Go")
			Expect(goSkills).To(HaveLen(1))
			Expect(goSkills[0].Confidence).To(Equal(0.95)) // Highest confidence kept
			Expect(goSkills[0].EventIDs).To(HaveLen(2))
		})
	})
})
