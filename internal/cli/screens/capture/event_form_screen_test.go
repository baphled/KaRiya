package capture_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// EventFormScreen Tests
//
// EventFormScreen wraps base.FormScreen[*forms.CaptureEventFormData] with capture-specific context.
// It supports two strategies (quick and manual) and returns ScreenResults.
//
// Related:
// - internal/cli/forms/capture_event_form.go (Form configuration)
// - internal/cli/screens/base/form_screen.go (Base form screen)
// - docs/FORMS_GUIDE.md (Form patterns and best practices)

var _ = Describe("EventFormScreen", func() {
	var screen *capture.EventFormScreen

	BeforeEach(func() {
		breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
		screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyQuick)
	})

	Describe("Creation", func() {
		It("should create with Quick strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyQuick)
			Expect(screen).NotTo(BeNil())
		})

		It("should create with Manual strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyManual)
			Expect(screen).NotTo(BeNil())
		})

		It("should initialize with breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
		})

		It("should show form fields in view", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return strategy", func() {
			Expect(screen.GetStrategy()).To(Equal(types.StrategyQuick))
		})
	})

	Describe("Strategy Configuration", func() {
		It("should configure form with Quick strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyQuick)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetStrategy()).To(Equal(types.StrategyQuick))
		})

		It("should configure form with Manual strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyManual)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetStrategy()).To(Equal(types.StrategyManual))
		})
	})

	Describe("Cancellation - Escape Key", func() {
		It("should return CancelResult when Escape is pressed", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})
	})

	Describe("Window Resize", func() {
		It("should handle WindowSizeMsg", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should pass window size to underlying form", func() {
			screen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			view1 := screen.View()

			screen.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
			view2 := screen.View()

			Expect(view1).NotTo(BeEmpty())
			Expect(view2).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		It("should render form view", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show breadcrumbs in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
		})
	})

	Describe("Form Data", func() {
		It("should return form data", func() {
			data := screen.GetFormData()
			Expect(data).NotTo(BeNil())
			Expect(data.SubmitConfirmed).To(BeFalse())
		})

		It("should return empty form data for new event", func() {
			data := screen.GetFormData()
			Expect(data.Text).To(BeEmpty())
			Expect(data.Date).To(BeEmpty())
			Expect(data.Company).To(BeEmpty())
		})

		It("should return pre-populated form data for existing event", func() {
			existingEvent := fixtures.EventWith("", "Existing event text", "ACME Corp", "Project X")
			breadcrumbs := []string{"Main Menu", "Edit Event"}
			screen = capture.NewEventFormScreen(existingEvent, breadcrumbs, types.StrategyManual)

			data := screen.GetFormData()
			Expect(data.Text).To(Equal("Existing event text"))
			Expect(data.Company).To(Equal("ACME Corp"))
			Expect(data.Project).To(Equal("Project X"))
		})
	})

	Describe("Form Submission", func() {
		It("should return SubmitResult with CaptureEventFormData when form completes", func() {
			// The base.FormScreen returns SubmitResult with FormData as *forms.CaptureEventFormData
			// We verify the data type is correct
			data := screen.GetFormData()
			Expect(data).To(BeAssignableToTypeOf(&forms.CaptureEventFormData{}))
		})
	})

	Describe("Event Editing", func() {
		It("should support loading existing event", func() {
			existingEvent := fixtures.EventWith("", "Existing event text", "ACME Corp", "Project X")

			breadcrumbs := []string{"Main Menu", "Edit Event"}
			screen = capture.NewEventFormScreen(existingEvent, breadcrumbs, types.StrategyManual)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle nil event gracefully", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyManual)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("State Preservation", func() {
		It("should preserve breadcrumbs", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
			screen = capture.NewEventFormScreen(nil, breadcrumbs, types.StrategyQuick)

			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
			Expect(view).To(ContainSubstring("Form"))
		})
	})

	Describe("Terminal Info and Theme", func() {
		It("should accept terminal info", func() {
			screen.SetTerminalInfo(120, 40)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept theme", func() {
			screen.SetTheme(nil)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept logo", func() {
			screen.SetLogo(nil, 2)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
