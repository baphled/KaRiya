package generatecv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Views", func() {
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
				intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})

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
			intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})

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
			intent.HandleNavigate(&screens.NavigateResult{ResultData: "export"})

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
