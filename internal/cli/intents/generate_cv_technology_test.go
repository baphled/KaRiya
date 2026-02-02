package intents

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("GenerateCV Technology Extraction", func() {
	var (
		intent      *GenerateCVIntent
		ctx         *GenerateCVContext
		skillRepo   careerRepo.SkillRepository
		eventRepo   careerRepo.EventRepository
		testContext context.Context
	)

	BeforeEach(func() {
		testContext = context.Background()

		// Create in-memory repositories
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()

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
			Events: func() []*career.Event {
				e1 := fixtures.EventWith("event-1", "Implemented API", "Test Co", "")
				e1.Skills = []string{"skill-ruby", "skill-postgres"}
				e2 := fixtures.EventWith("event-2", "Migrated Database", "Test Co", "")
				e2.Skills = []string{"skill-postgres"}
				return []*career.Event{e1, e2}
			}(),
			Facts:           []*career.Fact{},
			SkillRepository: skillRepo,
			EventRepository: eventRepo,
			AppContext:      testContext,
		}

		//nolint:errcheck // Test setup - error handling not relevant.
		skillRepo.Create(testContext, fixtures.SkillWith("skill-ruby", "Ruby", "backend", ""))
		//nolint:errcheck // Test setup - error handling not relevant.
		skillRepo.Create(testContext, fixtures.SkillWith("skill-postgres", "PostgreSQL", "database", ""))
		//nolint:errcheck // Test setup - error handling not relevant.
		skillRepo.Create(testContext, fixtures.SkillWith("skill-react", "React", "frontend", ""))

		// Add events to repository.
		for _, event := range ctx.Events {
			//nolint:errcheck // Test setup - error handling not relevant.
			eventRepo.Create(testContext, event)
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
				Expect(extractedMsg.Error).ToNot(HaveOccurred())

				// Should have extracted technologies
				// PostgreSQL has 2 events (event-1, event-2) - below threshold
				// Ruby has 1 event (event-1) - below threshold
				// React has 0 events - below threshold
				// So filtered list should be empty
				Expect(extractedMsg.Technologies).To(BeEmpty())
			})

			It("should provide focus area suggestion", func() {
				for i := range 5 {
					event := fixtures.EventWith("event-extra-"+string(rune(i+'0')), "Extra Event", "Test Co", "")
					event.Skills = []string{"skill-ruby", "skill-postgres"}
					eventRepo.Create(testContext, event)
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

	Describe("Technology Focus Selection", func() {
		BeforeEach(func() {
			// Set up state with extracted technologies
			intent.state.extractedTechnologies = []*ExtractedTechnology{
				{ID: "skill-1", Name: "Ruby", Category: "backend", EventCount: 5},
				{ID: "skill-2", Name: "PostgreSQL", Category: "database", EventCount: 4},
				{ID: "skill-3", Name: "React", Category: "frontend", EventCount: 3},
			}
			intent.state.technologiesAvailable = true
			intent.state.currentState = GenerateCVStateSelectTechnologyFocus
			intent.state.technologyFocusIndex = 0
		})

		Context("view rendering", func() {
			It("should show 3 technology focus options", func() {
				view := intent.View()

				Expect(view).To(ContainSubstring("Language Agnostic"))
				Expect(view).To(ContainSubstring("Generalist"))
				Expect(view).To(ContainSubstring("Specialist"))
			})

			It("should show technology count", func() {
				view := intent.View()

				// Should show how many technologies were found
				Expect(view).To(MatchRegexp("(?i)found.*3.*technolog"))
			})

			It("should disable Generalist and Specialist if < 3 technologies", func() {
				intent.state.extractedTechnologies = []*ExtractedTechnology{
					{ID: "skill-1", Name: "Ruby", Category: "backend", EventCount: 5},
				}
				intent.state.technologiesAvailable = false

				view := intent.View()

				// Should indicate only Language Agnostic is available
				Expect(view).To(MatchRegexp("(?i)(disabled|unavailable|only.*language agnostic)"))
			})

			It("should highlight current selection", func() {
				intent.state.technologyFocusIndex = 1 // Generalist

				view := intent.View()

				// Should have some indicator for current selection (▶ or similar)
				Expect(view).To(ContainSubstring("▶"))
			})
		})

		Context("navigation", func() {
			It("should move selection up with up/k", func() {
				intent.state.technologyFocusIndex = 1

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(intent.state.technologyFocusIndex).To(Equal(0))
			})

			It("should move selection down with down/j", func() {
				intent.state.technologyFocusIndex = 0

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(intent.state.technologyFocusIndex).To(Equal(1))
			})

			It("should not go below 0", func() {
				intent.state.technologyFocusIndex = 0

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(intent.state.technologyFocusIndex).To(Equal(0))
			})

			It("should not go above max index", func() {
				intent.state.technologyFocusIndex = 2

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(intent.state.technologyFocusIndex).To(Equal(2))
			})
		})

		Context("selection - Language Agnostic", func() {
			It("should transition to focus area selection", func() {
				intent.state.technologyFocusIndex = 0 // Language Agnostic

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should skip technology selection and go to focus area
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectFocusArea))
			})

			It("should set technology focus to language agnostic", func() {
				intent.state.technologyFocusIndex = 0

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("language_agnostic"))
			})
		})

		Context("selection - Generalist", func() {
			It("should transition to technology selection", func() {
				intent.state.technologyFocusIndex = 1 // Generalist

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should go to technology multi-select
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
			})

			It("should set technology focus to generalist", func() {
				intent.state.technologyFocusIndex = 1

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("generalist"))
			})
		})

		Context("selection - Specialist", func() {
			It("should transition to technology selection", func() {
				intent.state.technologyFocusIndex = 2 // Specialist

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should go to technology single-select
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
			})

			It("should set technology focus to specialist", func() {
				intent.state.technologyFocusIndex = 2

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("specialist"))
			})
		})

		Context("escape behavior", func() {
			It("should go back to extracting technologies state", func() {
				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should go back to previous state
				Expect(intent.state.currentState).To(Equal(GenerateCVStateExtractingTechnologies))
			})
		})
	})

	Describe("Technology Selection (Multi/Single)", func() {
		BeforeEach(func() {
			// Set up state with extracted technologies
			intent.state.extractedTechnologies = []*ExtractedTechnology{
				{ID: "skill-ruby", Name: "Ruby", Category: "backend", EventCount: 7},
				{ID: "skill-postgres", Name: "PostgreSQL", Category: "database", EventCount: 5},
				{ID: "skill-react", Name: "React", Category: "frontend", EventCount: 4},
				{ID: "skill-docker", Name: "Docker", Category: "devops", EventCount: 3},
			}
			intent.state.technologiesAvailable = true
			intent.state.currentState = GenerateCVStateSelectTechnologies
			intent.state.technologyCursor = 0
			intent.state.technologySelected = make(map[int]bool)
			intent.state.selectedTechnologies = []string{}
		})

		Context("Generalist mode - multi-select", func() {
			BeforeEach(func() {
				intent.state.selectedTechnologyFocus = cv.TechnologyFocusGeneralist
			})

			It("should show technology list with event counts", func() {
				view := intent.View()

				Expect(view).To(ContainSubstring("Ruby"))
				Expect(view).To(ContainSubstring("PostgreSQL"))
				Expect(view).To(ContainSubstring("React"))
				Expect(view).To(ContainSubstring("Docker"))
				Expect(view).To(MatchRegexp("7.*event")) // Event count for Ruby
			})

			It("should show instruction to select 2-5 technologies", func() {
				view := intent.View()

				Expect(view).To(MatchRegexp("(?i)(select|choose).*2.*5"))
			})

			It("should toggle selection with space", func() {
				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace})

				Expect(intent.state.technologySelected[0]).To(BeTrue())

				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace})

				Expect(intent.state.technologySelected[0]).To(BeFalse())
			})

			It("should move cursor with up/down", func() {
				intent.state.technologyCursor = 0

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(intent.state.technologyCursor).To(Equal(1))

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(intent.state.technologyCursor).To(Equal(0))
			})

			It("should not allow confirmation with < 2 selections", func() {
				intent.state.technologySelected[0] = true
				intent.state.selectedTechnologies = []string{"skill-ruby"}

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should not transition (still in select technologies state)
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
				Expect(cmd).To(BeNil())
			})

			It("should not allow confirmation with > 5 selections", func() {
				// Add more technologies to test > 5 selection
				intent.state.extractedTechnologies = append(intent.state.extractedTechnologies,
					&ExtractedTechnology{ID: "skill-5", Name: "Tech5", Category: "test", EventCount: 3},
					&ExtractedTechnology{ID: "skill-6", Name: "Tech6", Category: "test", EventCount: 3},
				)

				for i := range 6 {
					intent.state.technologySelected[i] = true
				}

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should not transition (too many selected)
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
				Expect(cmd).To(BeNil())
			})

			It("should transition to focus area with valid selection (2-5)", func() {
				intent.state.technologySelected[0] = true
				intent.state.technologySelected[1] = true
				intent.state.selectedTechnologies = []string{"skill-ruby", "skill-postgres"}

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectFocusArea))
				Expect(intent.state.selectedTechnologies).To(HaveLen(2))
			})

			It("should store selected technology IDs", func() {
				// Select Ruby and React
				intent.state.technologyCursor = 0
				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace}) // Select Ruby
				intent.state.technologyCursor = 2
				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace}) // Select React

				intent.state.selectedTechnologies = []string{"skill-ruby", "skill-react"}
				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.state.selectedTechnologies).To(ContainElement("skill-ruby"))
				Expect(intent.state.selectedTechnologies).To(ContainElement("skill-react"))
			})
		})

		Context("Specialist mode - single-select", func() {
			BeforeEach(func() {
				intent.state.selectedTechnologyFocus = cv.TechnologyFocusSpecialist
			})

			It("should show instruction to select 1 technology", func() {
				view := intent.View()

				Expect(view).To(MatchRegexp("(?i)(select|choose).*1.*technology"))
			})

			It("should select technology on space (replacing previous)", func() {
				intent.state.technologyCursor = 0
				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace})

				Expect(intent.state.technologySelected[0]).To(BeTrue())

				// Select different one
				intent.state.technologyCursor = 1
				_ = intent.Update(tea.KeyMsg{Type: tea.KeySpace})

				// Old selection should be cleared
				Expect(intent.state.technologySelected[0]).To(BeFalse())
				Expect(intent.state.technologySelected[1]).To(BeTrue())
			})

			It("should not allow confirmation without selection", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
				Expect(cmd).To(BeNil())
			})

			It("should transition to focus area with exactly 1 selection", func() {
				intent.state.technologySelected[0] = true
				intent.state.selectedTechnologies = []string{"skill-ruby"}

				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectFocusArea))
				Expect(intent.state.selectedTechnologies).To(HaveLen(1))
				Expect(intent.state.selectedTechnologies[0]).To(Equal("skill-ruby"))
			})
		})

		Context("escape behavior", func() {
			It("should go back to technology focus selection", func() {
				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologyFocus))
			})
		})
	})
})
