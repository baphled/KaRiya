package intents

import (
	"context"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/technology"
)

var _ = Describe("GenerateCV Technology Extraction", func() {
	var (
		intent      *GenerateCVIntent
		ctx         *GenerateCVContext
		skillRepo   careerRepo.SkillRepository
		eventRepo   careerRepo.Repository
		testContext context.Context
	)

	BeforeEach(func() {
		testContext = context.Background()

		// Create in-memory repositories
		skillRepo = careerRepo.NewMemorySkillRepository()
		eventRepo = careerRepo.NewMemoryRepository()

		// Create test context with repositories
		ctx = &GenerateCVContext{
			AvailableProfiles: []*CVProfile{
				{
					ID:          "profile-1",
					Name:        "Senior Backend Engineer",
					TargetRole:  "senior_ic",
					Description: "Test profile",
				},
			},
			Events: []*career.CareerEvent{
				{
					ID:        "event-1",
					Text:      "Implemented API",
					Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Company:   "Test Co",
					Skills:    []string{"skill-ruby", "skill-postgres"},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "event-2",
					Text:      "Migrated Database",
					Date:      time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					Company:   "Test Co",
					Skills:    []string{"skill-postgres"},
					CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			Facts:           []*career.Fact{},
			SkillRepository: skillRepo,
			EventRepository: eventRepo,
			AppContext:      testContext,
		}

		// Add test skills to repository
		_ = skillRepo.Create(testContext, &career.Skill{
			ID:       "skill-ruby",
			Name:     "Ruby",
			Category: "backend",
		})
		_ = skillRepo.Create(testContext, &career.Skill{
			ID:       "skill-postgres",
			Name:     "PostgreSQL",
			Category: "database",
		})
		_ = skillRepo.Create(testContext, &career.Skill{
			ID:       "skill-react",
			Name:     "React",
			Category: "frontend",
		})

		// Add events to repository
		for _, event := range ctx.Events {
			_ = eventRepo.Create(testContext, event)
		}

		// Create intent
		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		_ = intent.Init()
	})

	Describe("Technology Extraction Flow", func() {
		Context("when audience is selected", func() {
			It("should transition to extracting technologies state", func() {
				// Start at audience selection
				intent.state.currentState = GenerateCVStateSelectAudience

				// Select audience (should trigger extraction)
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())

				// Should transition to extracting state
				Expect(intent.state.currentState).To(Equal(GenerateCVStateExtractingTechnologies))
			})

			It("should extract technologies with 3+ events", func() {
				// Start at audience selection
				intent.state.currentState = GenerateCVStateSelectAudience

				// Select audience
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())

				// Execute the extraction command
				msg := cmd()

				// Should be TechnologiesExtractedMsg
				extractedMsg, ok := msg.(TechnologiesExtractedMsg)
				Expect(ok).To(BeTrue(), "Expected TechnologiesExtractedMsg")
				Expect(extractedMsg.Error).To(BeNil())

				// Should have extracted technologies
				// PostgreSQL has 2 events (event-1, event-2) - below threshold
				// Ruby has 1 event (event-1) - below threshold
				// React has 0 events - below threshold
				// So filtered list should be empty
				Expect(extractedMsg.Technologies).To(HaveLen(0))
			})

			It("should provide focus area suggestion", func() {
				// Add more events to meet threshold
				for i := 0; i < 5; i++ {
					event := &career.CareerEvent{
						ID:        "event-extra-" + string(rune(i+'0')),
						Text:      "Extra Event",
						Date:      time.Date(2025, 1, 3+i, 0, 0, 0, 0, time.UTC),
						Company:   "Test Co",
						Skills:    []string{"skill-ruby", "skill-postgres"},
						CreatedAt: time.Date(2025, 1, 3+i, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 1, 3+i, 0, 0, 0, 0, time.UTC),
					}
					_ = eventRepo.Create(testContext, event)
					ctx.Events = append(ctx.Events, event)
				}

				// Start at audience selection
				intent.state.currentState = GenerateCVStateSelectAudience

				// Select audience
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())

				// Execute extraction
				msg := cmd()

				// Should have suggestion
				extractedMsg, ok := msg.(TechnologiesExtractedMsg)
				Expect(ok).To(BeTrue())
				Expect(extractedMsg.Suggestion).NotTo(BeNil())

				// Should suggest backend (majority of skills are backend/database)
				Expect(string(extractedMsg.Suggestion.Area)).To(Equal("backend"))
			})
		})

		Context("when extraction completes", func() {
			It("should store extracted technologies in state", func() {
				// Create extraction message with technologies
				msg := TechnologiesExtractedMsg{
					Technologies: []*ExtractedTechnology{
						{
							ID:         "skill-ruby",
							Name:       "Ruby",
							Category:   "backend",
							EventCount: 5,
							EventIDs:   []string{"e1", "e2", "e3", "e4", "e5"},
						},
					},
					Suggestion: &FocusAreaSuggestion{
						Area:       technology.FocusAreaBackend,
						Confidence: 0.85,
						Evidence:   map[string]int{"backend": 5},
					},
					Error: nil,
				}

				// Start at extracting state
				intent.state.currentState = GenerateCVStateExtractingTechnologies

				// Process extraction message
				_ = intent.Update(msg)

				// Should store technologies in state
				Expect(intent.state.extractedTechnologies).To(HaveLen(1))
				Expect(intent.state.extractedTechnologies[0].Name).To(Equal("Ruby"))
				Expect(intent.state.focusAreaSuggestion).NotTo(BeNil())
				Expect(string(intent.state.focusAreaSuggestion.Area)).To(Equal("backend"))
			})

			It("should transition to technology focus selection", func() {
				msg := TechnologiesExtractedMsg{
					Technologies: []*ExtractedTechnology{
						{ID: "skill-1", Name: "Ruby", Category: "backend", EventCount: 5},
						{ID: "skill-2", Name: "PostgreSQL", Category: "database", EventCount: 4},
						{ID: "skill-3", Name: "React", Category: "frontend", EventCount: 3},
					},
					Suggestion: &FocusAreaSuggestion{
						Area:       technology.FocusAreaBackend,
						Confidence: 0.70,
						Evidence:   map[string]int{"backend": 5, "database": 4, "frontend": 3},
					},
					Error: nil,
				}

				intent.state.currentState = GenerateCVStateExtractingTechnologies
				_ = intent.Update(msg)

				// Should transition to technology focus selection
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologyFocus))
				Expect(intent.state.technologiesAvailable).To(BeTrue())
			})

			It("should mark technologies as unavailable if < 3 technologies", func() {
				msg := TechnologiesExtractedMsg{
					Technologies: []*ExtractedTechnology{
						{ID: "skill-1", Name: "Ruby", Category: "backend", EventCount: 5},
					},
					Suggestion: &FocusAreaSuggestion{
						Area:       technology.FocusAreaBackend,
						Confidence: 0.90,
						Evidence:   map[string]int{"backend": 5},
					},
					Error: nil,
				}

				intent.state.currentState = GenerateCVStateExtractingTechnologies
				_ = intent.Update(msg)

				// Should mark as unavailable (only Language Agnostic allowed)
				Expect(intent.state.technologiesAvailable).To(BeFalse())
			})
		})

		Context("when extraction fails", func() {
			It("should go back to audience selection", func() {
				msg := TechnologiesExtractedMsg{
					Technologies: nil,
					Suggestion:   nil,
					Error:        errors.New("extraction failed"),
				}

				intent.state.currentState = GenerateCVStateExtractingTechnologies
				_ = intent.Update(msg)

				// Should return to audience selection
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
			})
		})

		Context("when escape is pressed during extraction", func() {
			It("should navigate back to audience selection", func() {
				intent.state.currentState = GenerateCVStateExtractingTechnologies

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should go back (extraction continues in background)
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
			})
		})
	})
})
