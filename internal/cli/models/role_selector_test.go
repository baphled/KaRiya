package models_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
)

var _ = ginkgo.Describe("RoleSelector Model", func() {
	var model *models.RoleSelectorModel

	ginkgo.BeforeEach(func() {
		baseModel := models.NewBaseStandardModel()
		model = models.NewRoleSelectorModel(baseModel)
	})

	ginkgo.Describe("initialization", func() {
		ginkgo.It("should create a new role selector model", func() {
			gomega.Expect(model).NotTo(gomega.BeNil())
		})

		ginkgo.It("should initialize with first role selected", func() {
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal("principal"))
		})

		ginkgo.It("should have empty error initially", func() {
			gomega.Expect(model.GetLastError()).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("navigation", func() {
		ginkgo.It("should move down with j key", func() {
			initialRole := model.GetSelectedRole()
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			gomega.Expect(model.GetSelectedRole()).NotTo(gomega.Equal(initialRole))
		})

		ginkgo.It("should move down with down arrow", func() {
			initialRole := model.GetSelectedRole()
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			gomega.Expect(model.GetSelectedRole()).NotTo(gomega.Equal(initialRole))
		})

		ginkgo.It("should move up with k key", func() {
			// Move down first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			gomega.Expect(model.GetSelectedRole()).NotTo(gomega.Equal("principal"))

			// Now move up
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal("principal"))
		})

		ginkgo.It("should move up with up arrow", func() {
			// Move down first
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			gomega.Expect(model.GetSelectedRole()).NotTo(gomega.Equal("principal"))

			// Now move up
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal("principal"))
		})

		ginkgo.It("should not move above first role", func() {
			// Should be at principal
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal("principal"))

			// Try to move up
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			// Should still be at principal
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal("principal"))
		})

		ginkgo.It("should not move below last role", func() {
			// Move to last role
			for i := 0; i < 5; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}

			lastRole := model.GetSelectedRole()

			// Try to move down more
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Should still be at last role
			gomega.Expect(model.GetSelectedRole()).To(gomega.Equal(lastRole))
		})
	})

	ginkgo.Describe("role selection", func() {
		ginkgo.It("should trigger callback on enter", func() {
			var selectedRole string
			model.SetOnRoleSelected(func(role string) {
				selectedRole = role
			})

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})
			gomega.Expect(selectedRole).To(gomega.Equal("principal"))
		})

		ginkgo.It("should return correct role after navigation", func() {
			var selectedRole string
			model.SetOnRoleSelected(func(role string) {
				selectedRole = role
			})

			// Move to staff
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\r'}})

			gomega.Expect(selectedRole).To(gomega.Equal("staff"))
		})
	})

	ginkgo.Describe("rendering", func() {
		ginkgo.It("should render without error", func() {
			view := model.View()
			gomega.Expect(view).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should include header text", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Select Target Role"))
		})

		ginkgo.It("should display all roles", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Principal"))
			gomega.Expect(view).To(gomega.ContainSubstring("Staff"))
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
