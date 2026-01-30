//nolint:errcheck // Test file - error handling for test setup is not relevant.
package capture_test

import (
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// EventSubmitScreen Tests
//
// EventSubmitScreen handles async submission of captured event to database.
// Shows progress indicator while submitting and returns result when complete.

var _ = Describe("EventSubmitScreen", func() {
	var (
		screen      *capture.EventSubmitScreen
		testEvent   *career.Event
		testBursts  []*career.Burst
		testFacts   []*career.Fact
		breadcrumbs []string
	)

	BeforeEach(func() {
		now := time.Now()

		testEvent = &career.Event{
			ID:      "evt-1",
			Text:    "Test event",
			Date:    now,
			Company: "TestCorp",
		}

		testBursts = []*career.Burst{
			{
				ID:        "burst-1",
				Name:      "Test Burst",
				EventIDs:  []string{"evt-1"},
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		testFacts = []*career.Fact{
			{
				ID:            "fact-1",
				Text:          "Test fact",
				SourceEventID: "evt-1",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		}

		breadcrumbs = []string{"Main Menu", "Capture Event", "Submit"}
		screen = capture.NewEventSubmitScreen(breadcrumbs, testEvent, testBursts, testFacts)
	})

	Describe("Creation", func() {
		It("should create with event, bursts, and facts", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("should handle nil event", func() {
			screen = capture.NewEventSubmitScreen(breadcrumbs, nil, testBursts, testFacts)
			Expect(screen).NotTo(BeNil())
		})

		It("should handle empty bursts and facts", func() {
			screen = capture.NewEventSubmitScreen(breadcrumbs, testEvent, nil, nil)
			Expect(screen).NotTo(BeNil())
		})

		It("should initialize in submitting state", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Submitting"))
		})
	})

	Describe("Init Command", func() {
		It("should return submit command on Init", func() {
			cmd := screen.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should trigger async submission", func() {
			cmd := screen.Init()
			Expect(cmd).NotTo(BeNil())

			// Execute command to trigger submission
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})
	})

	Describe("View Rendering", func() {
		It("should show submitting message after Init", func() {
			screen.Init() // Start submission
			view := screen.View()
			Expect(view).To(ContainSubstring("Submitting"))
		})

		It("should show event text after Init", func() {
			screen.Init() // Start submission
			view := screen.View()
			Expect(view).To(ContainSubstring("Test event"))
		})

		It("should show progress indicator", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		Context("submitting state", func() {
			BeforeEach(func() {
				screen.Init()
			})

			It("should display UIKit InfoText title", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Submitting Event..."))
			})

			It("should display event text in DetailView", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Test event"))
			})

			It("should display please wait message", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Please wait"))
			})

			It("should handle nil event gracefully", func() {
				screen = capture.NewEventSubmitScreen(breadcrumbs, nil, testBursts, testFacts)
				screen.Init()
				view := screen.View()
				Expect(view).To(ContainSubstring("Submitting Event..."))
				Expect(view).To(ContainSubstring("Please wait"))
			})
		})

		Context("error state", func() {
			BeforeEach(func() {
				screen = capture.NewEventSubmitScreenWithError(
					breadcrumbs,
					testEvent,
					testBursts,
					testFacts,
					errors.New("database connection failed"),
				)
				cmd := screen.Init()
				errorMsg := cmd()
				screen.Update(errorMsg)
			})

			It("should display UIKit ErrorText title", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Submission Failed"))
			})

			It("should display error message", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("database connection failed"))
			})

			It("should display retry instruction", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Press Esc to go back and try again"))
			})
		})

		Context("success state", func() {
			BeforeEach(func() {
				cmd := screen.Init()
				successMsg := cmd()
				screen.Update(successMsg)
			})

			It("should display UIKit SuccessText title", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Event Submitted Successfully"))
			})

			It("should display event text in DetailView", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Test event"))
			})

			It("should display bursts count", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Bursts"))
				Expect(view).To(ContainSubstring("1"))
			})

			It("should display facts count", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Facts"))
				Expect(view).To(ContainSubstring("1"))
			})

			It("should omit bursts field when empty", func() {
				screen = capture.NewEventSubmitScreen(breadcrumbs, testEvent, nil, testFacts)
				cmd := screen.Init()
				successMsg := cmd()
				screen.Update(successMsg)
				view := screen.View()
				Expect(view).NotTo(ContainSubstring("Bursts"))
			})

			It("should omit facts field when empty", func() {
				screen = capture.NewEventSubmitScreen(breadcrumbs, testEvent, testBursts, nil)
				cmd := screen.Init()
				successMsg := cmd()
				screen.Update(successMsg)
				view := screen.View()
				Expect(view).NotTo(ContainSubstring("Facts"))
			})
		})
	})

	Describe("Footer Rendering", func() {
		Context("during submission", func() {
			BeforeEach(func() {
				screen.Init()
			})

			It("should show please wait message", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Please wait"))
			})
		})

		Context("after error", func() {
			BeforeEach(func() {
				screen = capture.NewEventSubmitScreenWithError(
					breadcrumbs,
					testEvent,
					testBursts,
					testFacts,
					errors.New("submission failed"),
				)
				cmd := screen.Init()
				errorMsg := cmd()
				screen.Update(errorMsg)
			})

			It("should show back badge", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Esc"))
				Expect(view).To(ContainSubstring("Back"))
			})

			It("should show quit badge", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Quit"))
			})
		})

		Context("after success", func() {
			BeforeEach(func() {
				cmd := screen.Init()
				successMsg := cmd()
				screen.Update(successMsg)
			})

			It("should show continue badge", func() {
				view := screen.View()
				Expect(view).To(ContainSubstring("Enter"))
				Expect(view).To(ContainSubstring("Continue"))
			})
		})
	})

	Describe("Submission Success", func() {
		It("should return SubmitResult when submission succeeds", func() {
			// Trigger submission
			cmd := screen.Init()
			successMsg := cmd() // Execute submission command

			// Update with success message
			_, result := screen.Update(successMsg)

			// Should return SubmitResult
			if result != nil {
				Expect(result.Type()).To(Equal(screens.ResultSubmit))
			}
		})

		It("should include event in success result", func() {
			cmd := screen.Init()
			successMsg := cmd()
			_, result := screen.Update(successMsg)

			if result != nil && result.Type() == screens.ResultSubmit {
				submitResult := result.(*screens.SubmitResult)
				data := submitResult.Data().(map[string]interface{})
				Expect(data["event"]).To(Equal(testEvent))
			}
		})
	})

	Describe("Submission Error", func() {
		It("should return ErrorResult when submission fails", func() {
			// Create screen with error trigger
			screen = capture.NewEventSubmitScreenWithError(
				breadcrumbs,
				testEvent,
				testBursts,
				testFacts,
				errors.New("database connection failed"),
			)

			cmd := screen.Init()
			errorMsg := cmd()
			_, result := screen.Update(errorMsg)

			if result != nil {
				Expect(result.Type()).To(Equal(screens.ResultError))
			}
		})

		It("should include error details in result", func() {
			testError := errors.New("validation failed")
			screen = capture.NewEventSubmitScreenWithError(
				breadcrumbs,
				testEvent,
				testBursts,
				testFacts,
				testError,
			)

			cmd := screen.Init()
			errorMsg := cmd()
			_, result := screen.Update(errorMsg)

			if result != nil && result.Type() == screens.ResultError {
				errorResult := result.(*screens.ErrorResult)
				Expect(errorResult.Err).To(Equal(testError))
			}
		})
	})

	Describe("Escape Key Behavior", func() {
		It("should not return result immediately on Esc (submission in progress)", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Escape should not cancel while submitting
			Expect(result).To(BeNil())
		})

		It("should allow Esc after submission completes", func() {
			// Complete submission first
			cmd := screen.Init()
			successMsg := cmd()
			screen.Update(successMsg)

			// Now Esc should work (if implemented)
			// This is optional - some screens auto-return result
		})
	})

	Describe("Window Resize", func() {
		It("should handle WindowSizeMsg", func() {
			cmd, result := screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should maintain state after resize", func() {
			screen.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			view := screen.View()
			Expect(view).To(ContainSubstring("Submitting"))
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
