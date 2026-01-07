package models_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// TestBurstEditor removed - tests are run by the main models_test.go suite

var _ = Describe("BurstEditorModel", func() {
	var (
		model   *models.BurstEditorModel
		burst   *career.Burst
		service *careerservice.Service
		ctx     context.Context
	)

	BeforeEach(func() {
		burst = &career.Burst{
			ID:          "burst-123",
			Name:        "Original Name",
			Description: "Original Description",
		}
		ctx = context.Background()
		// Note: service can be nil for basic tests since we're just testing form behavior
		service = nil
		model = models.NewBurstEditorModel(burst, service, ctx)
	})

	Describe("NewBurstEditorModel", func() {
		It("should create a model with burst data", func() {
			Expect(model).NotTo(BeNil())
			Expect(model.GetBurst()).To(Equal(burst))
		})

		It("should not be submitted initially", func() {
			Expect(model.IsSubmitted()).To(BeFalse())
		})

		It("should not be cancelled initially", func() {
			Expect(model.IsCancelled()).To(BeFalse())
		})

		It("should have no error initially", func() {
			Expect(model.GetError()).To(BeNil())
		})
	})

	Describe("Init", func() {
		It("should return a command", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle window size messages", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 50}
			updatedModel, cmd := model.Update(msg)

			Expect(updatedModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should handle quit on 'q'", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
			_, cmd := model.Update(msg)

			Expect(cmd).NotTo(BeNil())
		})

		It("should handle quit on 'ctrl+c'", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			_, cmd := model.Update(msg)

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Revert", func() {
		It("should revert changes to original burst", func() {
			// Simulate editing (in real usage, form would update formData)
			burst.Name = "Modified Name"
			burst.Description = "Modified Description"

			model.Revert()

			// After revert, the burst is reverted to original
			// Note: In the new huh implementation, Revert() updates the burst object
			Expect(model.GetBurst().Name).To(Equal("Original Name"))
			Expect(model.GetBurst().Description).To(Equal("Original Description"))
		})
	})

	Describe("View", func() {
		It("should render a view", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include header", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Burst Editor"))
		})
	})

	Describe("State management", func() {
		It("should track submitted state", func() {
			Expect(model.IsSubmitted()).To(BeFalse())
			// Note: Can't easily test submission without a real service and repository
		})

		It("should track cancelled state", func() {
			Expect(model.IsCancelled()).To(BeFalse())

			// Simulate escape key to cancel (huh forms handle Esc differently)
			// The form needs to be in a state where Esc will trigger abort
			// For now, just verify initial state
			// In real usage, pressing Esc in the form will trigger IsAborted() which sets cancelled = true
		})
	})
})
