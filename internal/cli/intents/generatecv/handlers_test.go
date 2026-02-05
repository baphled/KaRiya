package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Handlers", func() {
	var (
		intent *generatecv.Intent
	)

	BeforeEach(func() {
		ctx := &generatecv.IntentContext{
			AvailableProfiles: []*generatecv.CVProfile{
				{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
			},
			Events: []*career.Event{fixtures.Event("e1")},
		}
		var err error
		intent, err = generatecv.NewIntent(ctx)
		Expect(err).ToNot(HaveOccurred())

		termInfo := intent.GetTerminalInfo()
		termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	})

	Describe("handleExportCompleteMsg", func() {
		It("should transition to export complete on success", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})

			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should transition to export complete on error", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			intent.Update(generatecv.ExportCompleteMsg{
				Path:  "",
				Error: errors.New("export failed"),
			})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("handleCVGenerated", func() {
		It("should store generated CV and transition to review", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})

			cv := fixtures.CVView("cv-1")
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: cv})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("handleWizardComplete", func() {
		It("should store wizard configuration and start extraction", func() {
			intent.Init()
			cmd := intent.Update(generatecv.WizardCompleteMsg{
				ProfileID:    "p1",
				Audience:     "recruiter",
				TechFocus:    "all",
				Technologies: []string{"Go"},
				FocusArea:    "backend",
				SkillsFormat: "grouped",
				SkillsLimit:  10,
				CVLength:     "standard",
			})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("handleTechExtracted", func() {
		It("should start CV generation after tech extraction", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})

			cmd := intent.Update(generatecv.TechnologiesExtractedMsg{
				Technologies: []*generatecv.ExtractedTechnology{},
				Suggestion:   &generatecv.FocusAreaSuggestion{},
			})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Ctrl+C cancellation", func() {
		It("should cancel and quit from any state", func() {
			intent.Init()

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("ProfileConfig Integration", func() {
		var intentWithProfile *generatecv.Intent

		BeforeEach(func() {
			ctx := &generatecv.IntentContext{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events: []*career.Event{fixtures.Event("e1")},
				ProfileConfig: &config.ProfileConfig{
					Name:     "Test User",
					Email:    "test@example.com",
					Title:    "Staff Engineer",
					Location: "Test City, TC",
					GitHub:   "https://github.com/testuser",
				},
			}
			var err error
			intentWithProfile, err = generatecv.NewIntent(ctx)
			Expect(err).ToNot(HaveOccurred())

			termInfo := intentWithProfile.GetTerminalInfo()
			termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
		})

		Context("PreviewScreen", func() {
			BeforeEach(func() {
				intentWithProfile.Init()
				intentWithProfile.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intentWithProfile.Update(generatecv.TechnologiesExtractedMsg{})
				intentWithProfile.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intentWithProfile.Update(tea.KeyMsg{Type: tea.KeyEnter})
			})

			It("should display user name from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Test User"))
			})

			It("should display user email from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("test@example.com"))
			})

			It("should display user location from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Test City, TC"))
			})

			It("should display user GitHub from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("github.com/testuser"))
			})
		})

		Context("ReviewScreen", func() {
			BeforeEach(func() {
				intentWithProfile.Init()
				intentWithProfile.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intentWithProfile.Update(generatecv.TechnologiesExtractedMsg{})
				intentWithProfile.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			})

			It("should display user name from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Test User"))
			})

			It("should display user email from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("test@example.com"))
			})

			It("should display user location from ProfileConfig", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Test City, TC"))
			})
		})
	})
})
