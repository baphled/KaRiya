package skill_test

import (
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	skillview "github.com/baphled/kariya/internal/tui/views/skill"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddEdit", func() {
	var (
		modal    *skillview.AddEdit
		skill    = fixtures.SkillWithYears("skill-1", "Go Programming", "backend", 5)
		formData *forms.SkillFormData
	)

	BeforeEach(func() {
		skill.Level = "Expert"
		formData = forms.GetSkillFormData(skill)
	})

	Describe("NewAddEdit", func() {
		It("creates a modal with nil form data", func() {
			modal = skillview.NewAddEdit(nil, 120, 40)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("creates a modal with form data", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			Expect(modal).NotTo(BeNil())
		})

		It("creates a modal in edit mode with skill ID", func() {
			modal = skillview.NewAddEdit(formData, 120, 40, "skill-1")
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsEditMode()).To(BeTrue())
		})

		It("creates a modal in add mode without skill ID", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsEditMode()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("returns a command", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil when form is nil", func() {
			modal = skillview.NewAddEdit(nil, 120, 40)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("initializes form with proper dimensions", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("handles small window dimensions", func() {
			modal = skillview.NewAddEdit(formData, 50, 20)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("handles large window dimensions", func() {
			modal = skillview.NewAddEdit(formData, 200, 100)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
		})

		It("closes modal on Escape key", func() {
			Expect(modal.IsVisible()).To(BeTrue())
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles window resize", func() {
			_, completed, data := modal.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("returns nil data when not visible", func() {
			modal.Hide()
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})
	})

	Describe("View", func() {
		It("returns empty string when not visible", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("renders form when visible", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Show/Hide", func() {
		BeforeEach(func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
		})

		It("shows the modal", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("hides the modal", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("GetOriginalSkillID", func() {
		It("returns empty string when not in edit mode", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			Expect(modal.GetOriginalSkillID()).To(Equal(""))
		})

		It("returns skill ID when in edit mode", func() {
			modal = skillview.NewAddEdit(formData, 120, 40, "skill-123")
			Expect(modal.GetOriginalSkillID()).To(Equal("skill-123"))
		})
	})

	Describe("IsEditMode", func() {
		It("returns false when adding new skill", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			Expect(modal.IsEditMode()).To(BeFalse())
		})

		It("returns true when editing existing skill", func() {
			modal = skillview.NewAddEdit(formData, 120, 40, "skill-456")
			Expect(modal.IsEditMode()).To(BeTrue())
		})
	})

	Describe("Update with various messages", func() {
		BeforeEach(func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
		})

		It("handles other key messages", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles window size message", func() {
			_, completed, data := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles very small window size", func() {
			_, completed, data := modal.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles very large window size", func() {
			_, completed, data := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 100})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles multiple window resize messages", func() {
			_, _, _ = modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			_, _, _ = modal.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
			_, _, _ = modal.Update(tea.WindowSizeMsg{Width: 80, Height: 25})
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("handles key message after window resize", func() {
			_, _, _ = modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles other key types", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("handles rune key messages", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})
	})

	Describe("View rendering", func() {
		It("renders with different widths", func() {
			modal = skillview.NewAddEdit(formData, 80, 40)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with small width", func() {
			modal = skillview.NewAddEdit(formData, 50, 40)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with large width", func() {
			modal = skillview.NewAddEdit(formData, 200, 40)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Edit mode initialization", func() {
		It("preserves skill ID in edit mode", func() {
			modal = skillview.NewAddEdit(formData, 120, 40, "skill-edit-123")
			Expect(modal.GetOriginalSkillID()).To(Equal("skill-edit-123"))
			Expect(modal.IsEditMode()).To(BeTrue())
		})

		It("handles empty skill ID in edit mode", func() {
			modal = skillview.NewAddEdit(formData, 120, 40, "")
			Expect(modal.GetOriginalSkillID()).To(Equal(""))
			Expect(modal.IsEditMode()).To(BeFalse())
		})
	})

	Describe("Form initialization", func() {
		It("initializes with nil form data", func() {
			modal = skillview.NewAddEdit(nil, 120, 40)
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("initializes with provided form data", func() {
			modal = skillview.NewAddEdit(formData, 120, 40)
			Expect(modal).NotTo(BeNil())
		})
	})
})
