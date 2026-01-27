package skills_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillFormScreen", func() {
	var (
		screen *skills.SkillFormScreen
		skill  *career.Skill
	)

	BeforeEach(func() {
		skill = &career.Skill{
			ID:       "skill-1",
			Name:     "Kubernetes",
			Category: "devops",
			Level:    "advanced",
		}
	})

	Describe("Construction", func() {
		Context("creating new skill", func() {
			It("should create form screen for new skill", func() {
				screen = skills.NewSkillFormScreen(nil)
				Expect(screen).NotTo(BeNil())
			})

			It("should have mostly empty form data with default category", func() {
				screen = skills.NewSkillFormScreen(nil)
				data := screen.GetFormData()
				Expect(data.Name).To(BeEmpty())
				// Category gets default value from form (first in list: "backend")
				Expect(data.Category).NotTo(BeEmpty())
			})
		})

		Context("editing existing skill", func() {
			It("should create form screen for existing skill", func() {
				screen = skills.NewSkillFormScreen(skill)
				Expect(screen).NotTo(BeNil())
			})

			It("should pre-populate form data from skill", func() {
				screen = skills.NewSkillFormScreen(skill)
				data := screen.GetFormData()
				Expect(data.Name).To(Equal("Kubernetes"))
				Expect(data.Category).To(Equal("devops"))
				Expect(data.Level).To(Equal("advanced"))
			})

			It("should not set submit confirmed", func() {
				screen = skills.NewSkillFormScreen(skill)
				data := screen.GetFormData()
				Expect(data.SubmitConfirmed).To(BeFalse())
			})
		})
	})

	Describe("Update Handling", func() {
		BeforeEach(func() {
			screen = skills.NewSkillFormScreen(skill)
			screen.SetTerminalInfo(120, 40)
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})

		It("should return CancelResult on escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should delegate key messages to form", func() {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			_, result := screen.Update(msg)

			// Should not return a result (form is still active)
			Expect(result).To(BeNil())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = skills.NewSkillFormScreen(skill)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Manage Skills"))
		})

		It("should include form fields", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Skill Name"))
			Expect(view).To(ContainSubstring("Category"))
		})

		It("should show help text in footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Tab"))
		})
	})

	Describe("Form Data Management", func() {
		BeforeEach(func() {
			screen = skills.NewSkillFormScreen(skill)
		})

		It("should return form data", func() {
			data := screen.GetFormData()
			Expect(data).NotTo(BeNil())
			Expect(data.Name).To(Equal("Kubernetes"))
		})

		It("should allow form data modifications", func() {
			data := screen.GetFormData()
			data.Name = "Docker"
			data.Category = "devops"
			data.SubmitConfirmed = true

			// Verify data was modified
			updatedData := screen.GetFormData()
			Expect(updatedData.Name).To(Equal("Docker"))
			Expect(updatedData.SubmitConfirmed).To(BeTrue())
		})
	})

	Describe("Terminal Handling", func() {
		BeforeEach(func() {
			screen = skills.NewSkillFormScreen(skill)
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should rebuild form on dimension change", func() {
			screen.SetTerminalInfo(80, 24)
			view1 := screen.View()

			screen.SetTerminalInfo(120, 40)
			view2 := screen.View()

			// Views should be different after resize
			Expect(view1).NotTo(Equal(view2))
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			screen = skills.NewSkillFormScreen(skill)
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("Integration with forms package", func() {
		It("should use SkillFormData from forms package", func() {
			screen = skills.NewSkillFormScreen(skill)
			data := screen.GetFormData()

			// Verify data is of correct type
			var _ *forms.SkillFormData = data
		})

		It("should apply form data to domain object", func() {
			screen = skills.NewSkillFormScreen(nil)
			data := screen.GetFormData()

			// Simulate form submission
			data.Name = "Go"
			data.Category = "backend"
			data.Level = "expert"
			data.YearsUsed = "5"
			data.SubmitConfirmed = true

			// Create new skill from form data
			newSkill := &career.Skill{}
			forms.ApplySkillFormData(newSkill, data)

			Expect(newSkill.Name).To(Equal("Go"))
			Expect(newSkill.Category).To(Equal("backend"))
			Expect(newSkill.Level).To(Equal("expert"))
			Expect(newSkill.YearsUsed).NotTo(BeNil())
			Expect(*newSkill.YearsUsed).To(Equal(5))
		})
	})
})
