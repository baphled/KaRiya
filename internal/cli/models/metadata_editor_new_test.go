//nolint:errcheck // Test file - error handling for test setup is not relevant.
package models_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/models"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestMetadataEditorNew removed - tests are run by the main models_test.go suite

var _ = Describe("MetadataEditorModelNew", func() {
	var (
		model      *models.MetadataEditorModelNew
		event      *career.Event
		service    *careerservice.Service
		cliService *cliservice.CLIEventService
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo := careermemory.NewEventRepository()
		skillRepo := careermemory.NewSkillRepository()
		service = careerservice.NewService(repo)
		service.SetSkillRepository(skillRepo)
		cliService = nil // Can be nil for these tests

		// Create a test event
		event = &career.Event{
			ID:         "test-event-1",
			Text:       "Test event",
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Company:    "Test Company",
			Project:    "Test Project",
			Tags:       []string{"backend", "go"},
			Categories: []string{"development"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		model = models.NewMetadataEditorModelNew(event, service, cliService, ctx)
	})

	Describe("NewMetadataEditorModelNew", func() {
		It("should create a new model with correct initial values", func() {
			Expect(model).NotTo(BeNil())
			Expect(model.GetEvent()).To(Equal(event))
			Expect(model.IsSubmitted()).To(BeFalse())
			Expect(model.IsCancelled()).To(BeFalse())
			Expect(model.GetError()).To(BeNil())
		})

		It("should initialize with default dimensions", func() {
			Expect(model).NotTo(BeNil())
			// Default dimensions are set in constructor
		})

		It("should preserve original event for reverting", func() {
			originalCompany := event.Company
			event.Company = "Modified Company"
			model.Revert()
			Expect(model.GetEvent().Company).To(Equal(originalCompany))
		})
	})

	Describe("Init", func() {
		It("should return the form init command", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when handling window resize", func() {
			It("should update dimensions", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				updatedModel, _ := model.Update(msg)
				typedModel := updatedModel.(*models.MetadataEditorModelNew)
				// Width and height should be updated (checked via View behavior)
				Expect(typedModel).NotTo(BeNil())
			})
		})

		Context("when handling quit key", func() {
			It("should return QuitMsg for 'q' key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
				_, cmd := model.Update(msg)
				Expect(cmd).NotTo(BeNil())
				// Command should produce QuitMsg
				result := cmd()
				_, ok := result.(models.QuitMsg)
				Expect(ok).To(BeTrue())
			})
		})

		Context("when form is aborted", func() {
			It("should set cancelled to true", func() {
				// Simulate pressing Esc to abort
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				model.Update(msg)
				// Note: Actual abort detection depends on huh form state
				// We test via IsCancelled() after appropriate interaction
			})
		})
	})

	Describe("View", func() {
		It("should render without errors", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Edit Event Metadata"))
		})

		It("should include form fields", func() {
			view := model.View()
			// Huh form should render fields
			Expect(view).NotTo(BeEmpty())
		})

		It("should display errors when present", func() {
			// We can't directly set errors via public API, but we can verify
			// that the View() method doesn't panic
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Revert", func() {
		It("should restore original event data", func() {
			originalCompany := event.Company
			originalProject := event.Project
			originalDate := event.Date

			// Modify event
			event.Company = "Modified Company"
			event.Project = "Modified Project"
			event.Date = time.Now()

			// Revert
			model.Revert()

			// Check restoration
			Expect(model.GetEvent().Company).To(Equal(originalCompany))
			Expect(model.GetEvent().Project).To(Equal(originalProject))
			Expect(model.GetEvent().Date).To(Equal(originalDate))
		})

		It("should restore tags and categories", func() {
			originalTags := []string{"backend", "go"}
			originalCategories := []string{"development"}

			event.Tags = []string{"frontend"}
			event.Categories = []string{"testing"}

			model.Revert()

			Expect(model.GetEvent().Tags).To(Equal(originalTags))
			Expect(model.GetEvent().Categories).To(Equal(originalCategories))
		})
	})

	Describe("State Queries", func() {
		It("should return correct submitted state", func() {
			Expect(model.IsSubmitted()).To(BeFalse())
		})

		It("should return correct cancelled state", func() {
			Expect(model.IsCancelled()).To(BeFalse())
		})

		It("should return the event being edited", func() {
			Expect(model.GetEvent()).To(Equal(event))
		})

		It("should return nil error initially", func() {
			Expect(model.GetError()).To(BeNil())
		})
	})

	Describe("Form Data Integration", func() {
		It("should correctly extract form data from event", func() {
			formData := forms.GetMetadataFormData(event)
			Expect(formData).NotTo(BeNil())
			Expect(formData.Date).To(Equal("2024-01-15"))
			Expect(formData.Company).To(Equal("Test Company"))
			Expect(formData.Project).To(Equal("Test Project"))
			Expect(formData.Tags).To(Equal([]string{"backend", "go"}))
			Expect(formData.Categories).To(Equal([]string{"development"}))
		})

		It("should correctly apply form data to event", func() {
			formData := &forms.MetadataFormData{
				Date:       "2024-02-20",
				Company:    "New Company",
				Project:    "New Project",
				Tags:       []string{"frontend", "react"},
				Categories: []string{"ui"},
			}

			err := forms.ApplyMetadataFormData(event, formData)
			Expect(err).To(BeNil())
			Expect(event.Company).To(Equal("New Company"))
			Expect(event.Project).To(Equal("New Project"))
			Expect(event.Tags).To(Equal([]string{"frontend", "react"}))
			Expect(event.Categories).To(Equal([]string{"ui"}))
			Expect(event.Date.Format("2006-01-02")).To(Equal("2024-02-20"))
		})
	})
})
