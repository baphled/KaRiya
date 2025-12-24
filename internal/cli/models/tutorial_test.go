package models_test

import (
	"github.com/baphled/kariya/internal/cli/models"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TutorialModel", func() {
	var (
		model *models.TutorialModel
	)

	BeforeEach(func() {
		model = models.NewTutorialModel()
	})

	Context("when creating a new TutorialModel", func() {
		It("should initialize at welcome step", func() {
			Expect(model).NotTo(BeNil())
			Expect(model.CurrentStep()).To(Equal(0))
		})

		It("should not be completed initially", func() {
			Expect(model.IsCompleted()).To(BeFalse())
			Expect(model.IsSkipped()).To(BeFalse())
		})
	})

	Context("when navigating through tutorial", func() {
		It("should advance to next step with space", func() {
			msg := tea.KeyMsg{Type: tea.KeySpace}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.CurrentStep()).To(Equal(1))
		})

		It("should advance to next step with enter", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.CurrentStep()).To(Equal(1))
		})

		It("should advance to next step with right arrow", func() {
			msg := tea.KeyMsg{Type: tea.KeyRight}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.CurrentStep()).To(Equal(1))
		})

		It("should go back with left arrow", func() {
			// First advance
			msg1 := tea.KeyMsg{Type: tea.KeySpace}
			model, _ = model.Update(msg1).(*models.TutorialModel, nil)
			Expect(model.CurrentStep()).To(Equal(1))

			// Then go back
			msg2 := tea.KeyMsg{Type: tea.KeyLeft}
			updatedModel, _ := model.Update(msg2)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.CurrentStep()).To(Equal(0))
		})

		It("should not go before first step", func() {
			msg := tea.KeyMsg{Type: tea.KeyLeft}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.CurrentStep()).To(Equal(0))
		})

		It("should mark completed after last step", func() {
			// Advance through all steps
			for i := 0; i < 7; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			Expect(model.IsCompleted()).To(BeTrue())
		})
	})

	Context("when skipping tutorial", func() {
		It("should mark as skipped with escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.IsSkipped()).To(BeTrue())
			Expect(updatedTutorial.IsCompleted()).To(BeTrue())
		})

		It("should mark as skipped with 'q'", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial.IsSkipped()).To(BeTrue())
		})
	})

	Context("when rendering tutorial steps", func() {
		It("should render welcome step", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Welcome to KaRiya"))
		})

		It("should render modes step", func() {
			msg := tea.KeyMsg{Type: tea.KeySpace}
			model, _ = model.Update(msg).(*models.TutorialModel, nil)
			view := model.View()
			Expect(view).To(ContainSubstring("Capture Modes"))
		})

		It("should render capture step", func() {
			// Advance twice
			for i := 0; i < 2; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(ContainSubstring("Capturing Events"))
		})

		It("should render tags step", func() {
			// Advance 3 times
			for i := 0; i < 3; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(ContainSubstring("Tags & Categories"))
		})

		It("should render filtering step", func() {
			// Advance 4 times
			for i := 0; i < 4; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(ContainSubstring("Filtering & Searching"))
		})

		It("should render exporting step", func() {
			// Advance 5 times
			for i := 0; i < 5; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(ContainSubstring("Exporting Events"))
		})

		It("should render getting help step", func() {
			// Advance 6 times
			for i := 0; i < 6; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(ContainSubstring("Getting Help"))
		})
	})

	Context("when handling window size", func() {
		It("should update dimensions on window size change", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 30}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial).NotTo(BeNil())
		})
	})

	Context("when tutorial is completed", func() {
		It("should return empty view after completion", func() {
			// Advance through all steps
			for i := 0; i < 7; i++ {
				msg := tea.KeyMsg{Type: tea.KeySpace}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*models.TutorialModel)
			}
			view := model.View()
			Expect(view).To(BeEmpty())
		})

		It("should return empty view after skipping", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			model, _ = model.Update(msg).(*models.TutorialModel, nil)
			view := model.View()
			Expect(view).To(BeEmpty())
		})
	})

	Context("when handling other inputs", func() {
		It("should have Init method", func() {
			Expect(model.Init).NotTo(BeNil())
		})

		It("should handle other key inputs gracefully", func() {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			updatedModel, _ := model.Update(msg)
			updatedTutorial := updatedModel.(*models.TutorialModel)
			Expect(updatedTutorial).NotTo(BeNil())
		})
	})
})

