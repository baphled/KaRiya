package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data Persistence", func() {
	Context("when capturing events with SQLite", func() {
		It("should persist events to SQLite database", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "test-events.db")

			// Create SQLite repository
			repo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo.Close()
			})

			// Create service
			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			// Capture an event
			ctx := context.Background()
			eventText := "Led implementation of critical feature"
			eventDate := time.Now().Add(-24 * time.Hour)

			err = cliSvc.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				careerservice.ManualEntry,
				service.WithCompany("TechCorp"),
				service.WithTags([]string{"technical", "leadership"}),
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify event was persisted to database
			events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1), "Event should be persisted to database")
			Expect(events[0].Text).To(Equal(eventText))
			Expect(events[0].Company).To(Equal("TechCorp"))
			Expect(events[0].Tags).To(ContainElements("technical", "leadership"))
		})

		It("should survive application restart with SQLite", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "persistent-events.db")

			// Create first repository instance and add event
			repo1, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo1.Close()
			})

			svc1 := careerservice.NewService(repo1)
			cliSvc1 := service.NewCLIEventService(svc1)

			ctx := context.Background()
			eventText := "Architected microservices platform"
			eventDate := time.Now().Add(-7 * 24 * time.Hour)

			err = cliSvc1.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				careerservice.CVBackfill,
				service.WithCompany("StartupCo"),
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify event was saved in first instance
			events1, err := svc1.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events1).To(HaveLen(1))

			// Create second repository instance pointing to same database
			repo2, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo2.Close()
			})

			svc2 := careerservice.NewService(repo2)

			// Verify event persists across instances
			events2, err := svc2.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events2).To(HaveLen(1), "Event should persist across database connections")
			Expect(events2[0].Text).To(Equal(eventText))
			Expect(events2[0].Company).To(Equal("StartupCo"))
		})

		It("should create database at custom path", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "custom-location.db")

			// Create SQLite repository at custom path
			repo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo.Close()
			})

			// Verify database file was created
			Expect(dbPath).To(BeAnExistingFile())

			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			ctx := context.Background()
			err = cliSvc.CaptureEvent(
				ctx,
				"Test event at custom path",
				time.Now().Add(-1*time.Hour),
				careerservice.ManualEntry,
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify persistence
			events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1))
		})
	})
})

var _ = Describe("Form Submission Persistence", func() {
	Context("when user submits form", func() {
		It("should return event data from form submission", func() {
			// Create temporary repository
			repo := careerrepo.NewMemoryRepository()
			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			// Create form model
			form := models.NewFormModel(cliSvc)

			// Optional fields are visible by default in manual mode (showOptionalFields=true)
			// No need to toggle them

			// Helper to type text
			typeText := func(f *models.FormModel, text string) *models.FormModel {
				for _, r := range text {
					m, _ := f.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
					f = m.(*models.FormModel)
				}
				return f
			}

			// Helper to update form
			updateForm := func(f *models.FormModel, msg tea.Msg) (*models.FormModel, tea.Cmd) {
				m, cmd := f.Update(msg)
				return m.(*models.FormModel), cmd
			}

			// Simulate user input
			form = typeText(form, "Led critical infrastructure upgrade")
			Expect(form.GetInputValue(0)).To(Equal("Led critical infrastructure upgrade"))

			// Navigate to date field
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "today")

			// Navigate to company field
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "TechCorp")

			// Navigate to project field
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "Infrastructure")

			// Navigate through Tags, Categories fields to reach Submit button
			// Current position: ProjectField (3)
			// Need to reach: SubmitButton (6)
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // TagsField (4)
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // CategoriesField (5)
			_, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // SubmitButton (6)

			// Submit form
			_, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			// Execute the command and get the result
			msg := cmd()

			// Verify the form returns a SubmitMsg with the event data
			Expect(msg).NotTo(BeNil())
			submitMsg, ok := msg.(models.SubmitMsg)
			Expect(ok).To(BeTrue(), "Command should return a SubmitMsg")
			Expect(submitMsg.Err).To(BeNil(), "SubmitMsg should not have errors")
			Expect(submitMsg.Event).NotTo(BeNil(), "SubmitMsg should contain event data")
			Expect(submitMsg.Event.Text).To(Equal("Led critical infrastructure upgrade"))
			Expect(submitMsg.Event.Company).To(Equal("TechCorp"))
			Expect(submitMsg.Event.Project).To(Equal("Infrastructure"))

			// NOTE: The form no longer persists events. The intent is responsible for persistence.
			// This test verifies that the form correctly collects and returns the data.
		})
	})
})

var _ = Describe("Form Submission Debug", func() {
	Context("debugging form persistence", func() {
		It("should show what mode is being used in form submission", func() {
			// Create temporary repository
			repo := careerrepo.NewMemoryRepository()
			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			// Create form model
			form := models.NewFormModel(cliSvc)

			// Optional fields are visible by default in manual mode (showOptionalFields=true)
			// No need to toggle them

			// Check initial mode
			Expect(form.GetInputValue(0)).To(Equal(""))

			// Helper to type text
			typeText := func(f *models.FormModel, text string) *models.FormModel {
				for _, r := range text {
					m, _ := f.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
					f = m.(*models.FormModel)
				}
				return f
			}

			// Helper to update form
			updateForm := func(f *models.FormModel, msg tea.Msg) (*models.FormModel, tea.Cmd) {
				m, cmd := f.Update(msg)
				return m.(*models.FormModel), cmd
			}

			// Simulate user input
			form = typeText(form, "Test event for debugging")
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "today")
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "TestCorp")
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})
			form = typeText(form, "TestProject")
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab})

			// Navigate through Tags, Categories, Mode fields to reach Submit button
			// Current position: ProjectField (3)
			// Need to reach: SubmitButton (7)
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // TagsField (4)
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // CategoriesField (5)
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // ModeField (6)
			form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // SubmitButton (7)

			// Submit form
			form, cmd := updateForm(form, tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			// Execute the command and process result
			msg := cmd()
			_, _ = updateForm(form, msg)

			// Check all events in repository (regardless of mode)
			ctx := context.Background()
			allEvents, err := repo.List(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())

			GinkgoT().Logf("Total events in repo: %d", len(allEvents))
			for i, e := range allEvents {
				GinkgoT().Logf("Event %d: %s (mode should be visible in capture logs)", i, e.Text)
			}
		})
	})
})
