package components_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ViewEventSkillsModal", func() {
	var (
		modal  *components.ViewEventSkillsModal
		theme  themes.Theme
		skills []*career.Skill
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		yearsUsed := 5
		skills = []*career.Skill{
			{
				ID:        "skill-1",
				Name:      "Go",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: &yearsUsed,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "skill-2",
				Name:      "React",
				Category:  "frontend",
				Level:     "intermediate",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "skill-3",
				Name:      "Kubernetes",
				Category:  "devops",
				Level:     "expert",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		modal = components.NewViewEventSkillsModal("event-123", skills, theme)
	})

	Describe("NewViewEventSkillsModal", func() {
		It("should create a new modal with skills", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should handle nil skills slice", func() {
			nilModal := components.NewViewEventSkillsModal("event-123", nil, theme)
			Expect(nilModal).NotTo(BeNil())
			nilModal.Show()
			view := nilModal.View()
			Expect(view).To(ContainSubstring("No skills"))
		})

		It("should handle empty skills slice", func() {
			emptyModal := components.NewViewEventSkillsModal("event-123", []*career.Skill{}, theme)
			Expect(emptyModal).NotTo(BeNil())
			emptyModal.Show()
			view := emptyModal.View()
			Expect(view).To(ContainSubstring("No skills"))
		})
	})

	Describe("Show/Hide", func() {
		It("should show the modal", func() {
			Expect(modal.IsVisible()).To(BeFalse())
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide the modal", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("should update modal dimensions", func() {
			modal.SetDimensions(120, 40)
			modal.Show()
			// The view should render without error
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal.SetDimensions(100, 30)
			modal.Show()
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})

		It("should display skill names", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("React"))
			Expect(view).To(ContainSubstring("Kubernetes"))
		})

		It("should display skill categories", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("backend"))
			Expect(view).To(ContainSubstring("frontend"))
			Expect(view).To(ContainSubstring("devops"))
		})

		It("should display skill levels", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("advanced"))
			Expect(view).To(ContainSubstring("intermediate"))
			Expect(view).To(ContainSubstring("expert"))
		})

		It("should display years used when available", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("5"))
		})

		It("should show close instructions", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should show title with skill count", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).To(ContainSubstring("3"))
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal.SetDimensions(100, 30)
			modal.Show()
		})

		It("should close on escape key", func() {
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close on enter key", func() {
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close on 'q' key", func() {
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should close on backspace key", func() {
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should not process keys when hidden", func() {
			modal.Hide()
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			// Should not error, just return
		})

		It("should handle window size messages", func() {
			_, _ = modal.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetSkills", func() {
		It("should update the skills being displayed", func() {
			newYears := 10
			newSkills := []*career.Skill{
				{
					ID:        "skill-new",
					Name:      "Python",
					Category:  "backend",
					Level:     "expert",
					YearsUsed: &newYears,
				},
			}
			modal.SetSkills(newSkills)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Python"))
			Expect(view).NotTo(ContainSubstring("Go"))
		})
	})

	Describe("GetEventID", func() {
		It("should return the event ID", func() {
			Expect(modal.GetEventID()).To(Equal("event-123"))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Scrolling with many skills", func() {
		BeforeEach(func() {
			// Create many skills to test scrolling
			manySkills := make([]*career.Skill, 20)
			for i := 0; i < 20; i++ {
				years := i + 1
				manySkills[i] = &career.Skill{
					ID:        "skill-" + string(rune('a'+i)),
					Name:      "Skill " + string(rune('A'+i)),
					Category:  "category",
					Level:     "intermediate",
					YearsUsed: &years,
				}
			}
			modal = components.NewViewEventSkillsModal("event-123", manySkills, theme)
			modal.SetDimensions(100, 20) // Small height to force scrolling
			modal.Show()
		})

		It("should handle up/down scrolling", func() {
			// Initial view
			_ = modal.View()

			// Scroll down
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			viewAfterDown := modal.View()
			Expect(viewAfterDown).NotTo(BeEmpty())

			// Scroll up
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			viewAfterUp := modal.View()
			Expect(viewAfterUp).NotTo(BeEmpty())
		})

		It("should handle j/k vim-style scrolling", func() {
			_ = modal.View()

			// Scroll with j
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			viewAfterJ := modal.View()
			Expect(viewAfterJ).NotTo(BeEmpty())

			// Scroll with k
			_, _ = modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			viewAfterK := modal.View()
			Expect(viewAfterK).NotTo(BeEmpty())
		})
	})
})
