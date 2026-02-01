package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddEditModal", func() {
	var (
		modal *modals.AddEditModal
		skill *career.Skill
	)

	BeforeEach(func() {
		skill = fixtures.SkillWithYears("skill-1", "Go Programming", "backend", 5)
		skill.Level = "Expert"
	})

	Describe("NewAddEditModal (Add mode)", func() {
		It("creates a modal for adding a new skill", func() {
			modal = modals.NewAddEditModal(nil, 120, 40)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsEditMode()).To(BeFalse())
		})
	})

	Describe("NewAddEditModal (Edit mode)", func() {
		It("creates a modal for editing an existing skill", func() {
			modal = modals.NewAddEditModal(skill, 120, 40)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsEditMode()).To(BeTrue())
		})
	})

	Describe("Init", func() {
		It("returns a command", func() {
			modal = modals.NewAddEditModal(nil, 120, 40)
			cmd := modal.Init()
			// Form Init returns a command
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewAddEditModal(skill, 120, 40)
		})

		Context("when visible", func() {
			It("closes on Escape key", func() {
				Expect(modal.IsVisible()).To(BeTrue())
				_, completed, _ := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsVisible()).To(BeFalse())
				Expect(completed).To(BeFalse())
			})

			It("handles window size messages", func() {
				_, _, _ = modal.Update(tea.WindowSizeMsg{Width: 150, Height: 60})
				// Modal should still be visible
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when not visible", func() {
			BeforeEach(func() {
				modal.Hide()
			})

			It("ignores messages", func() {
				Expect(modal.IsVisible()).To(BeFalse())
				_, completed, skillData := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(completed).To(BeFalse())
				Expect(skillData).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		Context("when visible", func() {
			BeforeEach(func() {
				modal = modals.NewAddEditModal(skill, 120, 40)
			})

			It("renders the form", func() {
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("contains skill name field", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Skill Name"))
			})

			It("displays keyboard shortcuts in footer (KeyBadge pattern)", func() {
				view := modal.View()
				// Verify the footer contains expected keyboard hints.
				Expect(view).To(ContainSubstring("Tab"))
				Expect(view).To(ContainSubstring("Enter"))
				Expect(view).To(ContainSubstring("Esc"))
			})
		})

		Context("when not visible", func() {
			BeforeEach(func() {
				modal = modals.NewAddEditModal(nil, 120, 40)
				modal.Hide()
			})

			It("returns empty string", func() {
				view := modal.View()
				Expect(view).To(BeEmpty())
			})
		})
	})

	Describe("Show/Hide", func() {
		BeforeEach(func() {
			modal = modals.NewAddEditModal(nil, 120, 40)
		})

		It("hides the modal", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("shows the modal", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("GetOriginalSkill", func() {
		Context("in edit mode", func() {
			BeforeEach(func() {
				modal = modals.NewAddEditModal(skill, 120, 40)
			})

			It("returns the original skill", func() {
				Expect(modal.GetOriginalSkill()).To(Equal(skill))
			})
		})

		Context("in add mode", func() {
			BeforeEach(func() {
				modal = modals.NewAddEditModal(nil, 120, 40)
			})

			It("returns nil", func() {
				Expect(modal.GetOriginalSkill()).To(BeNil())
			})
		})
	})

	Describe("SkillEditData", func() {
		Describe("ToSkill", func() {
			It("converts form data to a new skill", func() {
				data := &modals.SkillEditData{
					Name:      "TypeScript",
					Category:  "frontend",
					Level:     "intermediate",
					YearsUsed: "3",
				}
				skill := data.ToSkill("")
				Expect(skill.Name).To(Equal("TypeScript"))
				Expect(skill.Category).To(Equal("frontend"))
				Expect(skill.Level).To(Equal("intermediate"))
				Expect(skill.YearsUsed).NotTo(BeNil())
				Expect(*skill.YearsUsed).To(Equal(3))
			})

			It("preserves skill ID when updating", func() {
				data := &modals.SkillEditData{
					Name:     "Go",
					Category: "backend",
				}
				skill := data.ToSkill("existing-id")
				Expect(skill.ID).To(Equal("existing-id"))
			})

			It("handles empty years", func() {
				data := &modals.SkillEditData{
					Name:      "Rust",
					Category:  "systems",
					YearsUsed: "",
				}
				skill := data.ToSkill("")
				Expect(skill.YearsUsed).To(BeNil())
			})

			It("handles invalid years", func() {
				data := &modals.SkillEditData{
					Name:      "Haskell",
					Category:  "functional",
					YearsUsed: "abc",
				}
				skill := data.ToSkill("")
				Expect(skill.YearsUsed).To(BeNil())
			})
		})
	})
})
