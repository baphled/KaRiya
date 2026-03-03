package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/cli/screens"
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

	Describe("HandleSubmit", func() {
		It("should mark intent as completed", func() {
			intent.Init()

			intent.HandleSubmit(&screens.SubmitResult{})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Completed))
		})

		It("should return nil command", func() {
			intent.Init()

			cmd := intent.HandleSubmit(&screens.SubmitResult{})

			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleError", func() {
		It("should store the error", func() {
			intent.Init()

			testErr := errors.New("test error")
			cmd := intent.HandleError(&screens.ErrorResult{Err: testErr})

			Expect(cmd).To(BeNil())
		})

		It("should return nil command", func() {
			intent.Init()

			cmd := intent.HandleError(&screens.ErrorResult{Err: errors.New("error")})

			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleCancel", func() {
		Context("from StateConfiguring", func() {
			It("should mark intent as cancelled", func() {
				intent.Init()

				intent.HandleCancel(&screens.CancelResult{})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from StateReview", func() {
			BeforeEach(func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			})

			It("should return to wizard state", func() {
				intent.HandleCancel(&screens.CancelResult{})

				Expect(intent.GetState()).To(Equal(generatecv.StateConfiguring))
			})
		})

		Context("from StatePreview", func() {
			BeforeEach(func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			})

			It("should return to review state", func() {
				cmd := intent.HandleCancel(&screens.CancelResult{})

				Expect(cmd).To(BeNil())
				Expect(intent.GetState()).To(Equal(generatecv.StateReview))
			})
		})

		Context("from default state", func() {
			BeforeEach(func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			})

			It("should mark intent as cancelled", func() {
				intent.HandleCancel(&screens.CancelResult{})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})
	})

	Describe("HandleNavigate", func() {
		Context("with CVView data", func() {
			It("should mark intent as completed", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				intent.HandleNavigate(&screens.NavigateResult{ResultData: fixtures.CVView("cv-1")})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})
		})

		Context("with string data", func() {
			BeforeEach(func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			})

			It("should handle preview navigation", func() {
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "preview"})

				Expect(intent.GetState()).To(Equal(generatecv.StatePreview))
			})

			It("should handle export navigation", func() {
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})

				Expect(intent.GetState()).To(Equal(generatecv.StateExporting))
			})

			It("should handle edit navigation", func() {
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "edit"})

				Expect(intent.GetState()).To(Equal(generatecv.StateConfiguring))
			})

			It("should handle unknown string navigation", func() {
				cmd := intent.HandleNavigate(&screens.NavigateResult{ResultData: "unknown"})

				Expect(cmd).To(BeNil())
			})
		})

		Context("with unknown data type", func() {
			It("should return nil", func() {
				intent.Init()

				cmd := intent.HandleNavigate(&screens.NavigateResult{ResultData: 42})

				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("Test Getters", func() {
		Describe("GetTestContext", func() {
			It("should return the intent context", func() {
				ctx := intent.GetTestContext()

				Expect(ctx).NotTo(BeNil())
				Expect(ctx.AvailableProfiles).To(HaveLen(1))
				Expect(ctx.Events).To(HaveLen(1))
			})
		})

		Describe("GetIsExporting", func() {
			It("should return false initially", func() {
				intent.Init()

				Expect(intent.GetIsExporting()).To(BeFalse())
			})

			It("should return true when set via test setter", func() {
				intent.Init()
				intent.SetIsExportingForTest(true)

				Expect(intent.GetIsExporting()).To(BeTrue())
			})
		})

		Describe("GetSelectedSkillsLimit", func() {
			It("should return the configured skills limit", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{
					ProfileID:   "p1",
					Audience:    "hiring_manager",
					SkillsLimit: 15,
				})

				Expect(intent.GetSelectedSkillsLimit()).To(Equal(15))
			})
		})
	})

	Describe("View Rendering", func() {
		Context("in StateConfiguring", func() {
			It("should render wizard view", func() {
				intent.Init()

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StateExtracting", func() {
			It("should render progress modal", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StateGenerating", func() {
			It("should render progress modal", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StateReview", func() {
			It("should render review screen", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StatePreview", func() {
			It("should render preview screen", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StateExporting", func() {
			It("should render export modal over review screen", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in StateExportComplete", func() {
			It("should render export complete view with file path", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})
				intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("Export Complete"))
			})

			It("should render export complete view with error", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})
				intent.Update(generatecv.ExportCompleteMsg{Path: "", Error: errors.New("export failed")})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})

			It("should render export complete view for clipboard", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})
				intent.SetSelectedExportOptionForTest(generatecv.ExportOptionClipboard)
				intent.Update(generatecv.ExportCompleteMsg{Path: "", Error: nil})

				view := intent.View()

				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("handleExportCompleteKeypress", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})
			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})
		})

		It("should complete intent on enter key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Completed))
		})

		It("should return to export selection on esc key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(generatecv.StateExportSelectLocation))
		})

		It("should not change state on other keys", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			Expect(intent.GetState()).To(Equal(generatecv.StateExportComplete))
		})
	})

	Describe("Export Format Display", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})
		})

		It("should display Markdown format correctly", func() {
			intent.SetSelectedExportFormatForTest(generatecv.ExportFormatMarkdown)
			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})

			view := intent.View()
			Expect(view).To(ContainSubstring("Markdown"))
		})

		It("should display YAML format correctly", func() {
			intent.SetSelectedExportFormatForTest(generatecv.ExportFormatYAML)
			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.yaml", Error: nil})

			view := intent.View()
			Expect(view).To(ContainSubstring("YAML"))
		})

		It("should display Text format as default", func() {
			intent.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.txt", Error: nil})

			view := intent.View()
			Expect(view).To(ContainSubstring("Text"))
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

			It("should display profile name from wizard selection", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Profile: Staff Engineer"))
			})

			It("should display audience from wizard selection", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("Audience: hiring_manager"))
			})

			It("should display generation settings section header", func() {
				view := intentWithProfile.View()
				Expect(view).To(ContainSubstring("🎯 Generation Settings"))
			})
		})
	})
})
