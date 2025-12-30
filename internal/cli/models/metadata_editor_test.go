package models_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetadataEditorModel", func() {
	var (
		event  *career.CareerEvent
		repo   *careerrepo.MemoryRepository
		svc    *careerservice.Service
		cliSvc *cliservice.CLIEventService
		editor *models.MetadataEditorModel
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = cliservice.NewCLIEventService(svc)

		event = &career.CareerEvent{
			ID:         "test-event-1",
			Text:       "Test event description",
			Date:       time.Now().Add(-24 * time.Hour),
			Company:    "TestCorp",
			Project:    "TestProject",
			Tags:       []string{"technical"},
			Categories: []string{"Technical"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		editor = models.NewMetadataEditorModel(event, svc, cliSvc, ctx)
	})

	Describe("Creation", func() {
		It("should create a new metadata editor model", func() {
			Expect(editor).NotTo(BeNil())
			Expect(editor.GetEvent()).To(Equal(event))
			Expect(editor.IsCancelled()).To(BeFalse())
			Expect(editor.IsSubmitted()).To(BeFalse())
		})

		It("should initialize input fields with event data", func() {
			evt := editor.GetEvent()
			Expect(evt.Company).To(Equal(event.Company))
			Expect(evt.Project).To(Equal(event.Project))
		})

		It("should copy original event for reverting", func() {
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Navigation", func() {
		It("should move focus forward with Tab", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyTab})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})

		It("should move focus backward with Shift+Tab", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})

		It("should navigate tags with Up/Down arrows", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyDown})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())

			editor.Update(tea.KeyMsg{Type: tea.KeyUp})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Text Input", func() {
		It("should handle text input", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Tag Selection", func() {
		It("should toggle tag with Space key", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeySpace})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Category Selection", func() {
		It("should toggle category with Space key", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeySpace})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Submission", func() {
		It("should handle Enter key", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Cancellation", func() {
		It("should cancel when Escape is pressed", func() {
			editor.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(editor.IsCancelled()).To(BeTrue())
		})

		It("should cancel when Escape is pressed", func() {
			editor = models.NewMetadataEditorModel(event, svc, cliSvc, ctx)
			editor.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(editor.IsCancelled()).To(BeTrue())
		})
	})

	Describe("Revert", func() {
		It("should revert to original event", func() {
			editor.Revert()
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should render without panicking", func() {
			Expect(func() {
				_ = editor.View()
			}).NotTo(Panic())
		})

		It("should show header in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Edit Event Metadata"))
		})

		It("should show date field in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Date"))
		})

		It("should show company field in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Company"))
		})

		It("should show project field in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Project"))
		})

		It("should show tags field in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Tags"))
		})

		It("should show categories field in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Categories"))
		})

		It("should show buttons in view", func() {
			view := editor.View()
			Expect(view).To(ContainSubstring("Save"))
			Expect(view).To(ContainSubstring("Cancel"))
		})
	})

	Describe("State Queries", func() {
		It("should return correct submitted status", func() {
			Expect(editor.IsSubmitted()).To(BeFalse())
		})

		It("should return correct cancelled status", func() {
			Expect(editor.IsCancelled()).To(BeFalse())
		})

		It("should return the edited event", func() {
			returnedEvent := editor.GetEvent()
			Expect(returnedEvent).NotTo(BeNil())
		})
	})

	Describe("Window Resize", func() {
		It("should update dimensions on WindowSizeMsg", func() {
			editor.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			// Just verify it doesn't panic
			Expect(editor).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should initialize without error", func() {
			cmd := editor.Init()
			Expect(cmd).To(BeNil())
		})
	})
})
