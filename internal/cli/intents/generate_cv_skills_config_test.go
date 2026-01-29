package intents

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
)

var _ = Describe("GenerateCV Skills Configuration Selection", func() {
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

		// Set up state for skills config selection
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
		intent.state.selectedTechnologyFocus = cv.TechnologyFocusGeneralist
		intent.state.selectedTechnologies = []string{"skill-ruby", "skill-postgres"}
		intent.state.selectedFocusArea = cv.FocusAreaBackend
		intent.state.currentState = GenerateCVStateSelectSkillsConfig
		intent.state.skillsConfigCursor = 0

		// Set defaults
		intent.state.selectedSkillsFormat = "flat"
		intent.state.selectedSkillsLimit = 0 // no limit
	})

	Describe("View rendering", func() {
		It("should show skills format options (flat and grouped)", func() {
			view := intent.View()

			Expect(view).To(ContainSubstring("Flat"))
			Expect(view).To(ContainSubstring("Grouped"))
		})

		It("should show skills limit option", func() {
			view := intent.View()

			Expect(view).To(MatchRegexp("(?i)limit|maximum|max"))
		})

		It("should show current format selection", func() {
			intent.state.selectedSkillsFormat = "grouped"

			view := intent.View()

			// Should indicate grouped is selected
			Expect(view).To(MatchRegexp("(?i)grouped.*✓|✓.*grouped"))
		})

		It("should show current limit value", func() {
			intent.state.selectedSkillsLimit = 10

			view := intent.View()

			Expect(view).To(ContainSubstring("10"))
		})

		It("should show cursor indicator", func() {
			intent.state.skillsConfigCursor = 0

			view := intent.View()

			Expect(view).To(ContainSubstring("▶"))
		})
	})

	Describe("Navigation", func() {
		It("should move cursor up with up/k", func() {
			intent.state.skillsConfigCursor = 1

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

			Expect(intent.state.skillsConfigCursor).To(Equal(0))
		})

		It("should move cursor down with down/j", func() {
			intent.state.skillsConfigCursor = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.state.skillsConfigCursor).To(Equal(1))
		})

		It("should not go below 0", func() {
			intent.state.skillsConfigCursor = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyUp})

			Expect(intent.state.skillsConfigCursor).To(Equal(0))
		})

		It("should not go above 1 (max index)", func() {
			intent.state.skillsConfigCursor = 1

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			Expect(intent.state.skillsConfigCursor).To(Equal(1))
		})
	})

	Describe("Format selection", func() {
		It("should toggle format between flat and grouped when cursor is on format (0)", func() {
			intent.state.skillsConfigCursor = 0
			intent.state.selectedSkillsFormat = "flat"

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

			Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
		})

		It("should toggle back to flat from grouped", func() {
			intent.state.skillsConfigCursor = 0
			intent.state.selectedSkillsFormat = "grouped"

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

			Expect(intent.state.selectedSkillsFormat).To(Equal("flat"))
		})
	})

	Describe("Limit adjustment", func() {
		It("should increase limit with right arrow when cursor is on limit (1)", func() {
			intent.state.skillsConfigCursor = 1
			intent.state.selectedSkillsLimit = 5

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRight})

			Expect(intent.state.selectedSkillsLimit).To(Equal(10))
		})

		It("should decrease limit with left arrow when cursor is on limit (1)", func() {
			intent.state.skillsConfigCursor = 1
			intent.state.selectedSkillsLimit = 10

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyLeft})

			Expect(intent.state.selectedSkillsLimit).To(Equal(5))
		})

		It("should not go below 0", func() {
			intent.state.skillsConfigCursor = 1
			intent.state.selectedSkillsLimit = 0

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyLeft})

			Expect(intent.state.selectedSkillsLimit).To(Equal(0))
		})

		It("should not go above 50", func() {
			intent.state.skillsConfigCursor = 1
			intent.state.selectedSkillsLimit = 50

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRight})

			Expect(intent.state.selectedSkillsLimit).To(Equal(50))
		})

		It("should allow setting to 0 (no limit)", func() {
			intent.state.skillsConfigCursor = 1
			intent.state.selectedSkillsLimit = 5

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyLeft})

			Expect(intent.state.selectedSkillsLimit).To(Equal(0))
		})
	})

	Describe("Selection confirmation", func() {
		It("should transition to generating state on enter", func() {
			intent.state.selectedSkillsFormat = "grouped"
			intent.state.selectedSkillsLimit = 15

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateGenerating))
		})

		It("should preserve selected format", func() {
			intent.state.selectedSkillsFormat = "grouped"

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
		})

		It("should preserve selected limit", func() {
			intent.state.selectedSkillsLimit = 20

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.state.selectedSkillsLimit).To(Equal(20))
		})
	})

	Describe("Back navigation", func() {
		It("should return to focus area selection on escape", func() {
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectFocusArea))
		})

		It("should preserve previous selections when going back", func() {
			intent.state.selectedSkillsFormat = "grouped"
			intent.state.selectedSkillsLimit = 10

			_ = intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
			Expect(intent.state.selectedSkillsLimit).To(Equal(10))
		})
	})
})
