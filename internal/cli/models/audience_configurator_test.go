package models_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
)

var _ = ginkgo.Describe("AudienceConfigurator Model", func() {
	var model *models.AudienceConfiguratorModel

	ginkgo.BeforeEach(func() {
		baseModel := models.NewBaseStandardModel()
		model = models.NewAudienceConfiguratorModel(baseModel, []string{})
	})

	ginkgo.Describe("initialization", func() {
		ginkgo.It("should create a new audience configurator model", func() {
			gomega.Expect(model).NotTo(gomega.BeNil())
		})

		ginkgo.It("should initialize with no selections", func() {
			gomega.Expect(model.GetSelectedAudiences()).To(gomega.HaveLen(0))
		})

		ginkgo.It("should initialize with first item focused", func() {
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal("hiring_manager"))
		})

		ginkgo.It("should support pre-selected audiences", func() {
			baseModel := models.NewBaseStandardModel()
			model = models.NewAudienceConfiguratorModel(baseModel, []string{"hiring_manager", "peer"})
			gomega.Expect(model.GetSelectedAudiences()).To(gomega.HaveLen(2))
		})

		ginkgo.It("should have empty error initially", func() {
			gomega.Expect(model.GetLastError()).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("navigation", func() {
		ginkgo.It("should move down with j key", func() {
			initialFocus := model.GetFocusedAudience()
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			gomega.Expect(model.GetFocusedAudience()).NotTo(gomega.Equal(initialFocus))
		})

		ginkgo.It("should move down with down arrow", func() {
			initialFocus := model.GetFocusedAudience()
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			gomega.Expect(model.GetFocusedAudience()).NotTo(gomega.Equal(initialFocus))
		})

		ginkgo.It("should move up with k key", func() {
			// Move down first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			gomega.Expect(model.GetFocusedAudience()).NotTo(gomega.Equal("hiring_manager"))

			// Now move up
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal("hiring_manager"))
		})

		ginkgo.It("should move up with up arrow", func() {
			// Move down first
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			gomega.Expect(model.GetFocusedAudience()).NotTo(gomega.Equal("hiring_manager"))

			// Now move up
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal("hiring_manager"))
		})

		ginkgo.It("should not move above first audience", func() {
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal("hiring_manager"))

			// Try to move up
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			// Should still be at hiring_manager
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal("hiring_manager"))
		})

		ginkgo.It("should not move below last audience", func() {
			// Move to last
			for i := 0; i < 5; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}

			lastAudience := model.GetFocusedAudience()

			// Try to move down more
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Should still be at last
			gomega.Expect(model.GetFocusedAudience()).To(gomega.Equal(lastAudience))
		})
	})

	ginkgo.Describe("selection with space", func() {
		ginkgo.It("should toggle selection with space", func() {
			gomega.Expect(model.IsSelected("hiring_manager")).To(gomega.BeFalse())
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			gomega.Expect(model.IsSelected("hiring_manager")).To(gomega.BeTrue())
		})

		ginkgo.It("should toggle off with space", func() {
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			gomega.Expect(model.IsSelected("hiring_manager")).To(gomega.BeTrue())
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			gomega.Expect(model.IsSelected("hiring_manager")).To(gomega.BeFalse())
		})

		ginkgo.It("should allow multiple selections", func() {
			// Select first
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			gomega.Expect(model.IsSelected("hiring_manager")).To(gomega.BeTrue())

			// Move to second
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Select second
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			gomega.Expect(model.IsSelected("recruiter")).To(gomega.BeTrue())

			// Verify both are selected
			selected := model.GetSelectedAudiences()
			gomega.Expect(selected).To(gomega.ContainElements("hiring_manager", "recruiter"))
		})
	})

	ginkgo.Describe("confirmation", func() {
		ginkgo.It("should require at least one selection", func() {
			var confirmedAudiences []string
			model.SetOnAudiencesSelected(func(audiences []string) {
				confirmedAudiences = audiences
			})

			// Try to confirm with no selections
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

			// Should have error
			gomega.Expect(model.GetLastError()).NotTo(gomega.BeNil())
			gomega.Expect(confirmedAudiences).To(gomega.HaveLen(0))
		})

		ginkgo.It("should confirm with one selection", func() {
			model.Update(tea.KeyMsg{Type: tea.KeySpace})

			var confirmedAudiences []string
			model.SetOnAudiencesSelected(func(audiences []string) {
				confirmedAudiences = audiences
			})

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

			gomega.Expect(confirmedAudiences).To(gomega.ContainElement("hiring_manager"))
		})

		ginkgo.It("should confirm with multiple selections", func() {
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			model.Update(tea.KeyMsg{Type: tea.KeySpace})

			var confirmedAudiences []string
			model.SetOnAudiencesSelected(func(audiences []string) {
				confirmedAudiences = audiences
			})

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

			gomega.Expect(confirmedAudiences).To(gomega.HaveLen(2))
			gomega.Expect(confirmedAudiences).To(gomega.ContainElements("hiring_manager", "recruiter"))
		})
	})

	ginkgo.Describe("rendering", func() {
		ginkgo.It("should render without error", func() {
			view := model.View()
			gomega.Expect(view).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should include header text", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Select Target Audiences"))
		})

		ginkgo.It("should display all audiences", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Hiring Manager"))
			gomega.Expect(view).To(gomega.ContainSubstring("Recruiter"))
			gomega.Expect(view).To(gomega.ContainSubstring("Peer"))
		})

		ginkgo.It("should show checkboxes", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.MatchRegexp(`\[\s*\]|\[.*\]`))
		})

		ginkgo.It("should handle window size updates", func() {
			_, cmd := model.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
			gomega.Expect(cmd).To(gomega.BeNil())
			gomega.Expect(model.GetWidth()).To(gomega.Equal(120))
			gomega.Expect(model.GetHeight()).To(gomega.Equal(30))
		})
	})

	ginkgo.Describe("model interface", func() {
		ginkgo.It("should implement Model interface", func() {
			var _ tea.Model = model
		})

		ginkgo.It("should implement StandardModel interface", func() {
			var _ models.StandardModel = model
		})

		ginkgo.It("should have Init method", func() {
			cmd := model.Init()
			gomega.Expect(cmd).To(gomega.BeNil())
		})

		ginkgo.It("should return BaseStandardModel", func() {
			gomega.Expect(model.BaseStandardModel).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("escape handling", func() {
		ginkgo.It("should handle escape key", func() {
			var backPressed bool
			model.SetOnBack(func() {
				backPressed = true
			})

			model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			gomega.Expect(backPressed).To(gomega.BeTrue())
		})
	})
})
