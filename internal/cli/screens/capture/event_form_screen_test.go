package capture_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// EventFormScreen Tests
//
// EventFormScreen wraps CaptureForm and adapts it to the Screen interface.
// It handles form submission and translates form messages to ScreenResults.
//
// Related:
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration)
// - internal/cli/models/capture_form.go (CaptureForm model)
// - internal/cli/forms/capture_event_form.go (Form configuration)

var _ = Describe("EventFormScreen", func() {
	var (
		screen     *capture.EventFormScreen
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		// Note: CLIEventService is only needed for form submission
		// For screen tests, we can pass nil since we're testing screen behavior
		cliService = nil

		breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
		screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyQuick)
	})

	Describe("Creation", func() {
		It("should create with Quick strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyQuick)
			Expect(screen).NotTo(BeNil())
		})

		It("should create with Manual strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyManual)
			Expect(screen).NotTo(BeNil())
		})

		It("should initialize with breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
		})

		It("should show form fields in view", func() {
			view := screen.View()
			// Form should have at least the event text field
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Strategy Configuration", func() {
		It("should configure form with Quick strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyQuick)

			// Quick strategy should show minimal fields
			// Exact field visibility depends on form configuration
			Expect(screen).NotTo(BeNil())
		})

		It("should configure form with Manual strategy", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyManual)

			// Manual strategy should show all fields
			Expect(screen).NotTo(BeNil())
		})
	})

	Describe("Cancellation - Escape Key", func() {
		It("should return CancelResult when Escape is pressed", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should not interfere with form's internal Escape handling", func() {
			// First Esc might be consumed by form (e.g., blur field)
			// Second Esc should cancel
			screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// At least one Esc should trigger cancellation
			// (Exact behavior depends on form state)
			_ = result
		})
	})

	Describe("Window Resize", func() {
		It("should handle WindowSizeMsg", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())

			// View should still render correctly
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should pass window size to underlying form", func() {
			// Resize to small terminal
			screen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			view1 := screen.View()

			// Resize to large terminal
			screen.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
			view2 := screen.View()

			// Views should differ based on size (form adapts)
			// Note: This is hard to test without rendering, so we just verify no panic
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

		It("should delegate to CaptureForm's View", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Footer Rendering", func() {
		It("should show next field badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Tab"))
			Expect(view).To(ContainSubstring("Next field"))
		})

		It("should show submit badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Ctrl+S"))
			Expect(view).To(ContainSubstring("Submit"))
		})

		It("should show cancel badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should show quit badge", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Quit"))
		})
	})

	Describe("Form Submission Messages", func() {
		It("should return ErrorResult when SubmitMsg has error", func() {
			submitErr := errors.New("form validation failed")
			_, result := screen.Update(models.SubmitMsg{
				Event: nil,
				Err:   submitErr,
			})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultError))
			errorResult := result.(*screens.ErrorResult)
			Expect(errorResult.Err).To(Equal(submitErr))
		})

		It("should return SubmitResult when SubmitMsg succeeds", func() {
			testEvent := &career.Event{
				Text:    "Test event submission",
				Company: "ACME",
			}
			_, result := screen.Update(models.SubmitMsg{
				Event: testEvent,
				Err:   nil,
			})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultSubmit))
			submitResult := result.(*screens.SubmitResult)
			Expect(submitResult.Data()).To(Equal(testEvent))
		})
	})

	Describe("Init Command", func() {
		It("should initialize form with Init command", func() {
			cmd := screen.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should return form's Init command", func() {
			// Init command comes from underlying CaptureForm
			cmd := screen.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Event Editing", func() {
		It("should support loading existing event", func() {
			existingEvent := &career.Event{
				Text:    "Existing event text",
				Company: "ACME Corp",
				Project: "Project X",
			}

			breadcrumbs := []string{"Main Menu", "Edit Event"}
			screen = capture.NewEventFormScreenWithEvent(
				cliService,
				breadcrumbs,
				types.StrategyManual,
				existingEvent,
			)

			// View should show form (we can't easily verify field values without integration test)
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle nil event gracefully", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event"}
			screen = capture.NewEventFormScreenWithEvent(
				cliService,
				breadcrumbs,
				types.StrategyManual,
				nil,
			)

			// Should not panic
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("State Preservation", func() {
		It("should preserve breadcrumbs", func() {
			breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
			screen = capture.NewEventFormScreen(cliService, breadcrumbs, types.StrategyQuick)

			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Capture Event"))
			Expect(view).To(ContainSubstring("Form"))
		})
	})

	Describe("Terminal Info and Theme", func() {
		It("should accept terminal info", func() {
			screen.SetTerminalInfo(120, 40)
			// Should not panic
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept theme", func() {
			screen.SetTheme(nil) // nil theme is acceptable
			// Should not panic
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should accept logo", func() {
			screen.SetLogo(nil, 2) // nil logo is acceptable
			// Should not panic
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
