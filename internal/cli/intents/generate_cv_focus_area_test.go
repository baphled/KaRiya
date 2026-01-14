package intents

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
)

var _ = Describe("GenerateCV Focus Area Selection", func() {
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
					Text:      "Test event",
					Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					Company:   "Test Co",
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			Facts:           []*career.Fact{},
			SkillRepository: skillRepo,
			EventRepository: eventRepo,
			AppContext:      testContext,
		}

		// Create intent
		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		_ = intent.Init()

		// Set up state for focus area selection
		intent.state.extractedTechnologies = []*ExtractedTechnology{
			{ID: "skill-ruby", Name: "Ruby", Category: "backend", EventCount: 7},
			{ID: "skill-postgres", Name: "PostgreSQL", Category: "database", EventCount: 5},
			{ID: "skill-react", Name: "React", Category: "frontend", EventCount: 2},
		}
		intent.state.focusAreaSuggestion = &FocusAreaSuggestion{
			Area:       technology.FocusAreaBackend,
			Confidence: 0.85,
			Evidence:   map[string]int{"backend": 7, "database": 5, "frontend": 2},
		}
		intent.state.currentState = GenerateCVStateSelectFocusArea
		intent.state.focusAreaCursor = 0
	})

	Describe("View rendering", func() {
		It("should show 4 focus area options", func() {
			view := intent.View()

			Expect(view).To(ContainSubstring("Backend"))
			Expect(view).To(ContainSubstring("Frontend"))
			Expect(view).To(ContainSubstring("Fullstack"))
			Expect(view).To(ContainSubstring("DevOps"))
		})

		It("should highlight suggested focus area", func() {
			view := intent.View()

			// Backend should be suggested (highest count in evidence)
			// Pattern: "Backend ⭐ (suggested)"
			Expect(view).To(MatchRegexp("(?i)backend.*⭐.*(suggested|recommended)"))
		})

		It("should show evidence counts", func() {
			view := intent.View()

			// Should show skill counts from evidence
			Expect(view).To(ContainSubstring("7")) // Backend count
			Expect(view).To(ContainSubstring("5")) // Database count
		})

		It("should show cursor indicator", func() {
			intent.state.focusAreaCursor = 1

			view := intent.View()

			Expect(view).To(ContainSubstring("▶"))
		})
	})

	Describe("Navigation", func() {
		It("should move cursor up with up/k", func() {
			intent.state.focusAreaCursor = 2

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

			Expect(intent.state.focusAreaCursor).To(Equal(1))
		})

		It("should move cursor down with down/j", func() {
			intent.state.focusAreaCursor = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.state.focusAreaCursor).To(Equal(1))
		})

		It("should not go below 0", func() {
			intent.state.focusAreaCursor = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

			Expect(intent.state.focusAreaCursor).To(Equal(0))
		})

		It("should not go above 3 (max index)", func() {
			intent.state.focusAreaCursor = 3

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.state.focusAreaCursor).To(Equal(3))
		})
	})

	Describe("Selection", func() {
		It("should store selected focus area on enter", func() {
			intent.state.focusAreaCursor = 0 // Backend

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(string(intent.state.selectedFocusArea)).To(Equal("backend"))
		})

		It("should transition to skills configuration selection", func() {
			intent.state.focusAreaCursor = 1 // Frontend

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectSkillsConfig))
			Expect(intent.state.selectedSkillsFormat).To(Equal("flat")) // Default
			Expect(intent.state.selectedSkillsLimit).To(Equal(0))       // Default (no limit)
		})

		It("should select Backend", func() {
			intent.state.focusAreaCursor = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedFocusArea).To(Equal(cv.FocusAreaBackend))
		})

		It("should select Frontend", func() {
			intent.state.focusAreaCursor = 1

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedFocusArea).To(Equal(cv.FocusAreaFrontend))
		})

		It("should select Fullstack", func() {
			intent.state.focusAreaCursor = 2

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedFocusArea).To(Equal(cv.FocusAreaFullstack))
		})

		It("should select DevOps", func() {
			intent.state.focusAreaCursor = 3

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedFocusArea).To(Equal(cv.FocusAreaDevOps))
		})
	})

	Describe("Escape behavior", func() {
		Context("when coming from Language Agnostic", func() {
			BeforeEach(func() {
				intent.state.selectedTechnologyFocus = cv.TechnologyFocusLanguageAgnostic
			})

			It("should go back to technology focus selection", func() {
				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologyFocus))
			})
		})

		Context("when coming from Generalist/Specialist", func() {
			BeforeEach(func() {
				intent.state.selectedTechnologyFocus = cv.TechnologyFocusGeneralist
			})

			It("should go back to technology selection", func() {
				_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectTechnologies))
			})
		})
	})
})
