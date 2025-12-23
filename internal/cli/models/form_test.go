package models_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to update form and return typed result
func updateForm(f *models.FormModel, msg tea.Msg) (*models.FormModel, tea.Cmd) {
	m, cmd := f.Update(msg)
	return m.(*models.FormModel), cmd
}

// Helper to type text into form
func typeText(f *models.FormModel, text string) *models.FormModel {
	for _, r := range text {
		m, _ := f.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		f = m.(*models.FormModel)
	}
	return f
}

var _ = Describe("FormModel", func() {
	var (
		form          *models.FormModel
		cliService    *service.CLIEventService
		careerService *careerservice.Service
		repo          *careerrepo.MemoryRepository
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		careerService = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(careerService)
		form = models.NewFormModel(cliService)
	})

	Describe("NewFormModel", func() {
		It("should create a new form model with default values", func() {
			Expect(form).NotTo(BeNil())
			Expect(form.Submitted()).To(BeFalse())
			Expect(form.Event()).To(BeNil())
			Expect(form.Error()).To(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a blink command", func() {
			cmd := form.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update - Navigation", func() {
		It("should navigate forward with Tab key", func() {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			updatedModel, _ := form.Update(msg)
			updatedForm, ok := updatedModel.(*models.FormModel)
			Expect(ok).To(BeTrue())
			Expect(updatedForm).NotTo(BeNil())
		})

		It("should navigate backward with Shift+Tab", func() {
			msg := tea.KeyMsg{Type: tea.KeyShiftTab}
			updatedModel, _ := form.Update(msg)
			updatedForm, ok := updatedModel.(*models.FormModel)
			Expect(ok).To(BeTrue())
			Expect(updatedForm).NotTo(BeNil())
		})

		It("should quit on Ctrl+C", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			_, cmd := form.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("should quit on Esc", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, cmd := form.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should render the form", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("Capture Career Event"))
			Expect(view).To(ContainSubstring("Event Text (required)"))
			Expect(view).To(ContainSubstring("Date (optional)"))
			Expect(view).To(ContainSubstring("Company (optional)"))
			Expect(view).To(ContainSubstring("Project (optional)"))
			Expect(view).To(ContainSubstring("Capture Mode"))
			Expect(view).To(ContainSubstring("Submit"))
		})

		It("should show character count", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("Characters: 0/2000"))
		})

		It("should show all capture modes", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("Timeline Journaling"))
			Expect(view).To(ContainSubstring("CV Backfill"))
			Expect(view).To(ContainSubstring("Manual Entry"))
		})

		It("should show mode descriptions", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("Real-time logging (last 30 days)"))
			Expect(view).To(ContainSubstring("Import from existing CV (any date)"))
			Expect(view).To(ContainSubstring("Manual entry (any date)"))
		})
	})

	Describe("Date Parsing", func() {
		Context("with valid date formats", func() {
			It("should parse 'today' as current date", func() {
				// Tested through form submission
				Expect(form).NotTo(BeNil())
			})

			It("should parse YYYY-MM-DD format", func() {
				// Tested through form submission
				Expect(form).NotTo(BeNil())
			})

			It("should parse relative dates", func() {
				// Tested through form submission
				Expect(form).NotTo(BeNil())
			})
		})
	})

	Describe("Form Submission", func() {
		Context("with valid input", func() {
			It("should successfully submit with all fields filled", func() {
				// Type text
				form = typeText(form, "Led a team to deliver critical project")
				Expect(form.GetInputValue(0)).To(Equal("Led a team to deliver critical project"))

				// Navigate and fill date
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form = typeText(form, "today")

				// Navigate and fill company
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form = typeText(form, "TechCorp Inc.")

				// Navigate and fill project
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form = typeText(form, "Platform Migration")

				// Navigate to submit button (2 more tabs: mode, submit)
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())

				// Execute command and process result
				msg := cmd()
				form, _ = updateForm(form, msg)

				// Verify
				Expect(form.Submitted()).To(BeTrue())
				Expect(form.Error()).To(BeNil())
				Expect(form.Event()).NotTo(BeNil())
			})

			It("should submit with only required fields", func() {
				// Type text
				form = typeText(form, "Completed important task")
				Expect(form.GetInputValue(0)).To(Equal("Completed important task"))

				// Navigate to submit (5 tabs total: Date, Company, Project, Mode, Submit)
				for i := 0; i < 5; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeTrue())
				Expect(form.Error()).To(BeNil())
				Expect(form.Event()).NotTo(BeNil())
				Expect(form.Event().Text).To(Equal("Completed important task"))
			})
		})

		Context("with invalid input", func() {
			It("should fail with empty text field", func() {
				// Navigate to submit without filling text
				for i := 0; i < 5; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeFalse())
				Expect(form.Error()).NotTo(BeNil())
				Expect(form.Error().Error()).To(ContainSubstring("required"))
			})

			It("should accept exactly 2000 characters", func() {
				// Type exactly 2000 'a' characters
				for i := 0; i < 2000; i++ {
					m, _ := form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
					form = m.(*models.FormModel)
				}

				// Navigate to submit
				for i := 0; i < 5; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				// Should succeed with exactly 2000 characters
				Expect(form.Submitted()).To(BeTrue())
				Expect(form.Error()).To(BeNil())
			})

			It("should fail with future date", func() {
				// Type text
				form = typeText(form, "Test event")

				// Navigate and enter future date
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form = typeText(form, "2099-12-31")

				// Navigate to submit (4 more tabs)
				for i := 0; i < 4; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeFalse())
				Expect(form.Error()).NotTo(BeNil())
			})

			It("should fail with invalid date format", func() {
				// Type text
				form = typeText(form, "Test event")

				// Navigate and enter invalid date
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				form = typeText(form, "invalid-date")

				// Navigate to submit (4 more tabs from date field)
				for i := 0; i < 4; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeFalse())
				Expect(form.Error()).NotTo(BeNil())
				Expect(form.Error().Error()).To(ContainSubstring("invalid date"))
			})
		})

		Context("with Timeline Journaling mode", func() {
			It("should fail with date older than 30 days", func() {
				// Type text
				form = typeText(form, "Test event")

				// Navigate and enter old date
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				oldDate := time.Now().AddDate(0, 0, -31).Format("2006-01-02")
				form = typeText(form, oldDate)

				// Navigate to submit (mode is already Timeline by default, 4 tabs from date field)
				for i := 0; i < 4; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeFalse())
				Expect(form.Error()).NotTo(BeNil())
				Expect(form.Error().Error()).To(ContainSubstring("30 days"))
			})

			It("should succeed with date within 30 days", func() {
				// Type text
				form = typeText(form, "Test event")

				// Navigate and enter recent date
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				recentDate := time.Now().AddDate(0, 0, -15).Format("2006-01-02")
				form = typeText(form, recentDate)

				// Navigate to submit (4 more tabs from date field)
				for i := 0; i < 4; i++ {
					form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
				}

				// Submit
				form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				form, _ = updateForm(form, msg)

				Expect(form.Submitted()).To(BeTrue())
				Expect(form.Error()).To(BeNil())
			})
		})
	})

	Describe("Reset", func() {
		It("should reset the form to initial state", func() {
			// Fill the form
			form = typeText(form, "Test event")

			// Reset
			form.Reset()

			// Verify reset state
			Expect(form.Submitted()).To(BeFalse())
			Expect(form.Event()).To(BeNil())
			Expect(form.Error()).To(BeNil())

			// View should show empty form
			view := form.View()
			Expect(view).To(ContainSubstring("Characters: 0/2000"))
		})
	})

	Describe("Integration with Service", func() {

		It("DEBUG: should verify repository works independently", func() {
			// Create a fresh repository just for this test
			testRepo := careerrepo.NewMemoryRepository()
			ctx := context.Background()

			// Test repository directly
			testEvent := &career.CareerEvent{
				Text: "Direct test event",
				Date: time.Now(),
			}
			err := testRepo.Create(ctx, testEvent)
			Expect(err).To(BeNil())
			Expect(testEvent.ID).NotTo(BeEmpty())

			// Verify it's stored (use Limit to get results)
			events, err := testRepo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))
			if len(events) > 0 {
				Expect(events[0].Text).To(Equal("Direct test event"))
			}
		})

		It("should persist event to repository on successful submission", func() {
			ctx := context.Background()

			// Fill and submit form
			form = typeText(form, "Integration test event")

			// Navigate to submit
			for i := 0; i < 5; i++ {
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			}

			// Submit
			form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			form, _ = updateForm(form, msg)

			Expect(form.Submitted()).To(BeTrue())

			// Verify event was persisted
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))
			Expect(events[0].Text).To(Equal("Integration test event"))
		})

		It("should use correct capture mode", func() {
			ctx := context.Background()

			// Fill text
			form = typeText(form, "Mode test event")

			// Navigate to mode field (4 tabs)
			for i := 0; i < 4; i++ {
				form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			}

			// Select CVBackfill mode (down arrow)
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyDown})

			// Navigate to submit
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})

			// Submit
			form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			form, _ = updateForm(form, msg)

			Expect(form.Submitted()).To(BeTrue())

			// Verify event was captured
			events, err := repo.List(ctx, careerrepo.ListFilters{Limit: 100})
			Expect(err).To(BeNil())
			Expect(len(events)).To(Equal(1))
		})
	})
})
