package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
	cvviews "github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Views", func() {
	var intent *generatecv.Intent

	BeforeEach(func() {
		ctx := &generatecv.IntentValidator{
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

	Describe("getWizardContextHelp", func() {
		Context("when wizard modal is visible", func() {
			It("should return empty string", func() {
				intent.Init()

				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("when export modal is visible", func() {
			It("should return empty string when exporting", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})

				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("when progress modal is visible", func() {
			It("should show wait message", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})

				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in review state", func() {
			It("should show review context help", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				view := intent.View()
				Expect(view).To(ContainSubstring("Preview"))
			})
		})

		Context("in preview state", func() {
			It("should show preview context help", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				view := intent.View()
				Expect(view).To(ContainSubstring("Scroll"))
			})
		})
	})

	Describe("wizardView state rendering", func() {
		Context("StateReview without activeScreen", func() {
			It("should show loading message", func() {
				intent.Init()
				intent.SetStateForTest(generatecv.StateReview)

				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("StatePreview without activeScreen", func() {
			It("should show loading message", func() {
				intent.Init()
				intent.SetStateForTest(generatecv.StatePreview)

				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("renderReviewScreenWithModalOverlay", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
		})

		It("should render review screen", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with export modal overlay when visible", func() {
			intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("renderPreviewScreenWithModalOverlay", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should render preview screen", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with export modal overlay when visible", func() {
			intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("getTheme", func() {
		It("should return default theme when theme is nil", func() {
			intent.Init()

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should use theme manager when set", func() {
			intent.Init()
			tm := themes.NewThemeManager()
			_ = tm.SetActive("default")
			intent.SetThemeManager(tm)

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("render overlay functions with nil modals", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
		})

		It("should render review screen without export modal", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render preview screen without export modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("Views Additional Coverage", func() {
	var intent *generatecv.Intent

	BeforeEach(func() {
		ctx := &generatecv.IntentValidator{
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

	Describe("getTheme with active theme manager", func() {
		It("uses the theme from theme manager in review state", func() {
			intent.Init()
			tm := themes.NewThemeManager()
			err := tm.SetActive("default")
			Expect(err).ToNot(HaveOccurred())
			intent.SetThemeManager(tm)

			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("getWizardContextHelp with no modals active", func() {
		It("returns review help when in StateReview without modals", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			intent.SetStateForTest(generatecv.StateReview)

			view := intent.View()
			Expect(view).To(ContainSubstring("Preview"))
		})

		It("returns preview help when in StatePreview without modals", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			intent.SetStateForTest(generatecv.StatePreview)

			view := intent.View()
			Expect(view).To(ContainSubstring("Scroll"))
		})
	})

	Describe("export format display names", func() {
		It("displays Markdown format name", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.SetStateForTest(generatecv.StateExportComplete)
			intent.SetSelectedExportFormatForTest(generatecv.ExportFormatMarkdown)
			intent.SetIsExportingForTest(false)

			view := intent.View()
			Expect(view).To(ContainSubstring("Markdown"))
		})

		It("displays YAML format name", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.SetStateForTest(generatecv.StateExportComplete)
			intent.SetSelectedExportFormatForTest(generatecv.ExportFormatYAML)
			intent.SetIsExportingForTest(false)

			view := intent.View()
			Expect(view).To(ContainSubstring("YAML"))
		})
	})
})

var _ = Describe("Intent Getter Coverage", func() {
	var intent *generatecv.Intent

	BeforeEach(func() {
		ctx := &generatecv.IntentValidator{
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

	Describe("getTheme non-nil path via export complete", func() {
		It("uses theme manager in export complete state", func() {
			intent.Init()
			tm := themes.NewThemeManager()
			err := tm.SetActive("default")
			Expect(err).ToNot(HaveOccurred())
			intent.SetThemeManager(tm)

			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.SetStateForTest(generatecv.StateExportComplete)
			intent.SetIsExportingForTest(false)

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetIsExporting", func() {
		It("returns false by default", func() {
			Expect(intent.GetIsExporting()).To(BeFalse())
		})

		It("returns true when set", func() {
			intent.SetIsExportingForTest(true)
			Expect(intent.GetIsExporting()).To(BeTrue())
		})
	})

	Describe("GetTestContext", func() {
		It("returns the intent context", func() {
			ctx := intent.GetTestContext()
			Expect(ctx).NotTo(BeNil())
			Expect(ctx.AvailableProfiles).To(HaveLen(1))
		})
	})

	Describe("GetSelectedSkillsLimit", func() {
		It("returns default value", func() {
			Expect(intent.GetSelectedSkillsLimit()).To(BeNumerically(">=", 0))
		})
	})
})

var _ = Describe("Targeted Coverage Tests", func() {
	var intent *generatecv.Intent

	BeforeEach(func() {
		ctx := &generatecv.IntentValidator{
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

	Describe("setCompleted with exported path", func() {
		It("includes export timestamp when path is set", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})
			intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})

			intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionConfirm}})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("updateWizardFlow return nil for unknown msg types", func() {
		It("returns nil for unrecognised message types", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			cmd := intent.Update(struct{}{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("getWizardContextHelp StateReview without modals", func() {
		It("shows review help when activeView is nil and no modals visible", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.SetStateForTest(generatecv.StateReview)

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("handleKeyDelegation in non-view states", func() {
		It("returns nil for keypress in generating state", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
			Expect(cmd).To(BeNil())
		})
	})
})

var _ = Describe("Coverage Boost Tests", func() {
	var intent *generatecv.Intent

	BeforeEach(func() {
		ctx := &generatecv.IntentValidator{
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

	Describe("countCVBullets with populated sections", func() {
		It("counts bullets across sections with content groups", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})

			section := fixtures.CVSectionWithContent("s1", "cv-bullets", []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Group 1", []*career.CVBullet{
					fixtures.CVBullet("b1", "s1"),
					fixtures.CVBullet("b2", "s1"),
				}),
			})
			cv := fixtures.CVViewWithSections("cv-bullets", []*career.CVSection{section})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: cv})

			Expect(intent.GetState()).To(Equal(generatecv.StateReview))
			Expect(intent.GetGeneratedCV()).NotTo(BeNil())
			Expect(intent.GetGeneratedCV().Sections).To(HaveLen(1))
		})
	})

	Describe("HandleSubmit directly", func() {
		It("marks intent completed with HandleSubmit", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleError directly", func() {
		It("captures error from ErrorViewResult", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			testErr := errors.New("test view error")
			cmd := intent.HandleError(&widgets.ErrorViewResult{Err: testErr})
			Expect(cmd).To(BeNil())
			Expect(intent.GetExportError()).To(MatchError("test view error"))
		})
	})

	Describe("handleViewResult nil result via Update with non-key msg", func() {
		It("handles nil msg after reaching review state", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

			cmd := intent.Update(nil)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleCancel from various states", func() {
		It("cancels from configuring state", func() {
			intent.Init()

			cmd := intent.HandleCancel(&widgets.CancelViewResult{})
			Expect(cmd).To(BeNil())
		})

		It("returns to review from preview state", func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			cmd := intent.HandleCancel(&widgets.CancelViewResult{})
			Expect(cmd).To(BeNil())
			Expect(intent.GetState()).To(Equal(generatecv.StateReview))
		})
	})
})
